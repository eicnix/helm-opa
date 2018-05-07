# Helm  [OPA](https://www.openpolicyagent.org/) Plugin
[![CircleCI](https://circleci.com/gh/eicnix/helm-opa.svg?style=svg)](https://circleci.com/gh/eicnix/helm-opa)

This plugin enables you to check your rendered templates files again [Open Policy Agent](https://www.openpolicyagent.org/) policies to ensure that they match your policies.

<a href="https://asciinema.org/a/180130" target="_blank"><img src="https://asciinema.org/a/180130.png" /></a>

## Usage

Define policies inside the `policies` folder in your Chart. e.g. Pods must run as nonRoot:
```
deny[msg] {
  input.kind = "Deployment"
  not input.spec.template.spec.securityContext.runAsNonRoot = true
  msg = "Containers must not run as root"
}
```

Run the policy check:
```sh
# helm opa CHART
Processing file deployment.yaml
Violations:
- Containers must not run as root
Processing file ingress.yaml
Processing file service.yaml
===
Result: Chart is not compliant
```

## Install

```sh
$ helm plugin install https://github.com/eicnix/helm-opa
```

## License

Copyright 2018 Lukas Eichler

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
