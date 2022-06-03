# Development

<!-- toc -->

```zsh 
➜  make help

Usage:
  make <target>

General
  help             Display this help.

Development
  fmt              Run go fmt against code.
  vet              Run go vet against code.
  test             Run tests.
  test-fswatch     Use fswatch to watch source files and run tests on chamnge

Run
  run              Run codegen and start controller from your host.

Docker
  kubeserial-docker-local        Build image for local development, tag local, supports only builder platform
  kubeserial-docker-all          Build and push image for all target platforms
  device-monitor-docker-local    Build image for local development, tag local, supports only builder platform
  device-monitor-docker-all      Build and push image for all target platforms
  injector-webhook-docker-local  Build image for local development, tag local, supports only builder platform
  injector-webhook-docker-all    Build and push image for all target platforms

Helm
  update-kubeserial-chart-version       Update version used in chart. Requires VERSION var to be set
  update-kubeserial-crds-chart-version  Update version used in chart. Requires VERSION var to be set
  helm-lint                             Run chart-testing to lint kubeserial chart.

Kind

Minikube
  minikube              Start local cluster, build image and deploy
  minikube-start        Start minikube cluster
  minikube-set-context  Set context to use minikube cluster
  minikube-deploy       Deploy the app to local minikube

Deployment
  uninstall        Uninstall release.
  deploy-dev       Install dev release in current context/namespace.

Docs
  docs-deps        Install mdbook (requires rust and cargo) + plugins
  docs-serve       Build docs, start server and open in browser

Build
  kubeserial        Build manager binary.
  device-monitor    Build device monitor binary
  injector-webhook  Build sidecar injector webhook binary binary
  all               Run codegen and build all components.
```