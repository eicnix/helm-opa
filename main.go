package main

import (
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/ghodss/yaml"
	"github.com/open-policy-agent/opa/ast"
	"github.com/open-policy-agent/opa/rego"
	"github.com/spf13/cobra"
	"k8s.io/helm/pkg/chartutil"
	"k8s.io/helm/pkg/engine"
	"k8s.io/helm/pkg/proto/hapi/chart"
	"k8s.io/helm/pkg/timeconv"
)

const globalUsage = `
Checks the rendered chart templates against Open Policy Agent policies.
All policies under policies/ will be evaluated. 
`

var (
	flagVerbose bool
	showNotes   bool
)

var version = "DEV"

func main() {
	cmd := &cobra.Command{
		Use:   "template [flags] CHART",
		Short: fmt.Sprintf("locally render templates (helm-template %s)", version),
		RunE:  run,
	}

	f := cmd.Flags()
	f.BoolVarP(&flagVerbose, "verbose", "v", false, "show the computed YAML values as well.")
	f.BoolVar(&showNotes, "notes", false, "show the computed NOTES.txt file as well.")

	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func run(cmd *cobra.Command, args []string) error {
	if len(args) < 1 {
		return errors.New("chart is required")
	}
	c, err := chartutil.Load(args[0])
	if err != nil {
		return err
	}

	vv, err := yaml.Marshal(map[string]interface{}{})
	config := &chart.Config{Raw: string(vv), Values: map[string]*chart.Value{}}

	options := chartutil.ReleaseOptions{
		Name:      "RELEASE",
		Time:      timeconv.Now(),
		Namespace: "NAMESPACE",
	}

	renderer := engine.New()

	vals, err := chartutil.ToRenderValues(c, config, options)
	if err != nil {
		return err
	}

	out, err := renderer.Render(c, vals)
	if err != nil {
		return err
	}

	sortedKeys := make([]string, 0, len(out))
	for key := range out {
		sortedKeys = append(sortedKeys, key)
	}
	sort.Strings(sortedKeys)

	foundViolations := false

	compiler, err := buildCopmiler(args[0] + "/" + "policies")

	for _, name := range sortedKeys {
		data := out[name]
		fileName := filepath.Base(name)

		if strings.HasSuffix(fileName, ".yaml") {
			r, _ := processFile(fileName, data, compiler)
			if !foundViolations && r {
				foundViolations = r
			}
		}

	}

	fmt.Println("===")
	if foundViolations {
		fmt.Println("Result: Chart is compliant")
	} else {
		fmt.Println("Result: Chart is not compliant")
	}

	return nil
}

func processFile(fileName string, data string, compiler *ast.Compiler) (bool, error) {
	fmt.Printf("Processing file %v\n", fileName)

	ctx := context.Background()
	var input interface{}
	err := yaml.Unmarshal([]byte(data), &input)

	rego := rego.New(
		rego.Query("data.main.deny"),
		rego.Compiler(compiler),
		rego.Input(input))

	rs, err := rego.Eval(ctx)

	hasResults := func(expression interface{}) bool {
		if v, ok := expression.([]interface{}); ok {
			return len(v) > 0
		}
		return false
	}

	foundViolations := false
	for _, r := range rs {
		for _, e := range r.Expressions {
			value := e.Value
			if hasResults(value) {
				foundViolations = true
				fmt.Println("Violations:")
				for _, v := range value.([]interface{}) {
					fmt.Printf("- %v\n", v)
				}
			}
		}
	}

	return foundViolations, err
}

func buildCopmiler(path string) (*ast.Compiler, error) {
	files, err := ioutil.ReadDir(path)
	if err != nil {
		return nil, err
	}

	modules := map[string]*ast.Module{}

	for _, file := range files {
		if !strings.HasSuffix(file.Name(), ".rego") {
			continue
		}

		out, err := ioutil.ReadFile(path + "/" + file.Name())
		if err != nil {
			return nil, err
		}

		parsed, err := ast.ParseModule(file.Name(), string(out[:]))
		if err != nil {
			return nil, err
		}
		modules[file.Name()] = parsed

	}

	compiler := ast.NewCompiler()
	compiler.Compile(modules)

	if compiler.Failed() {
		panic(compiler.Errors)
	}

	return compiler, nil
}
