# Development

<!-- toc -->

## Intro

Most things you'd need for local development are covered in Makefile. Just run `make help` to see available commands:

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

## Running tests

After any change run

```bash
$ make test
```

to run tests suite

You can also find it helpful to just run tests every time there is a change in project source files, for this run

```bash
$ make test-fswatch
```

which will run `test` target every time it detects change using `fswatch` (`fswatch` must be installed)

## Building images

There are sets of 2 targets for each image that is a part of this project:

> *-docker-local

which builds docker image for local development

and 

> *-docker-all

which builds docker image for all target architectures.

For local development you'll only need to build local images

### Building images localy

<mark>Docker buildx builder with support for platforms listed in TARGET_PLATFORMS is required</mark>

If you want to build all images using your local Docker, run

```zsh
$ make docker-local 
```

This will execute all `*-docker-local` targets and build images using your local Docker.

## Running minikube

TODO
<!-- TODO -->
