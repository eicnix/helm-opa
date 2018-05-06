# Helm OPA Plugin
[![CircleCI](https://circleci.com/gh/eicnix/helm-opa.svg?style=svg)](https://circleci.com/gh/eicnix/helm-opa)

This plugin enables you to check your rendered templates files again OPA policies to ensure that they match your policies.
<a href="hhttps://asciinema.org/a/179963" target="_blank"><img src="https://asciinema.org/a/179963" width="589"/></a>

yaml```
helm opa CHART
```

## Install

$ helm plugin install https://github.com/eicnix/helm-opa

## License
Copyright [2018] [Lukas Eichler]

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
