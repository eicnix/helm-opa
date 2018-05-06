package main


deny[msg] {
  input.kind = "Deployment"
  not input.spec.template.spec.securityContext.runAsNonRoot = true
  msg = "Containers must not run as root"
}

deny[msg] {
  input.kind = "Deployment"
  input.spec.selectors.matchLabels.app = null
  msg = "Containers must provide app/release labls for pod selectors"
}