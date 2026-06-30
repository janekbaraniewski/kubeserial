# Development

<!-- toc -->

## Intro

Most things you'd need for local development are covered in the Makefile. Just run `make help` to see available commands:

```zsh
➜  make help

Usage:
  make <target>

General
  help             Display this help.

Development
  generate         Generate files
  check-generated  Check that generated files are up to date
  fmt              Run go fmt against code.
  vet              Run go vet against code.
  test             Run tests.
  test-fswatch     Use fswatch to watch source files and run tests on change
  lint             Run golangci-lint
  check            Run linters and check code gen

Run
  run              Run codegen and start controller from your host.

Docker
  docker-local                   Build images for local dev
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
  update-version                        Update charts version.

Kind
  kind                  Create kind cluster, install certmanager, build and load images, install dev release.
  install-certmanager   Install cert-manager from jetstack/cert-manager

Minikube
  minikube              Start local cluster, build image and deploy
  minikube-start        Start minikube cluster
  minikube-set-context  Set context to use minikube cluster

Deployment
  uninstall        Uninstall release.
  deploy-dev       Install dev release in current context/namespace.
  install-dev      Install dev release in current context/namespace.

Docs
  docs-deps        Install mdbook (requires rust and cargo) + plugins
  docs-serve       Build docs, start server and open in browser
```

## Code generation

KubeSerial uses two kinds of generated files: the typed client / deepcopy / openapi code under `pkg/generated` and `pkg/apis`, and the CRD manifests rendered into `charts/kubeserial-crds`. Both are produced by scripts in `hack/`:

```bash
$ make generate          # regenerate client code (code-gen) and CRD manifests (manifests-gen)
$ make check-generated   # verify the committed output is up to date (used in CI)
```

`make run` and `make test` run `generate` for you, so after editing types under `pkg/apis/v1alpha1` you normally don't need to invoke it by hand.

## Running tests

After any change run

```bash
$ make test
```

to run the test suite. This formats and vets the code, downloads `setup-envtest`, renders the CRDs and starts an `envtest` control plane (the `ENVTEST_K8S_VERSION` set in the Makefile) so the controller integration tests under `pkg/controllers/integration_tests` can run against a real API server. Coverage is written to `coverage.txt`.

You can also find it helpful to just run tests every time there is a change in project source files, for this run

```bash
$ make test-fswatch
```

which will run the `test` target every time it detects a change using `fswatch` (`fswatch` must be installed).

## Linting

```bash
$ make lint    # golangci-lint
$ make check   # check-generated + lint, the combined CI gate
```

## Building images

There are sets of 2 targets for each of the three images in this project (`kubeserial`, `device-monitor`, `injector-webhook`):

> *-docker-local

which builds the docker image for local development, and

> *-docker-all

which builds and pushes the image for all target architectures (the platforms listed in the `TARGET_PLATFORMS` file).

For local development you'll only need to build local images.

### Building images locally

<mark>A Docker buildx builder with support for the platforms listed in TARGET_PLATFORMS is required</mark>

To build all three images at once using your local Docker, run

```zsh
$ make docker-local
```

This runs every `*-docker-local` target and loads the images into your local Docker.

## Running on a local cluster

There are two ready-made flows for spinning up a full local environment: kind and minikube. Both build the images, install cert-manager (the webhook needs it for its serving certificate) and deploy a dev release of the CRDs and the controllers.

### kind

```bash
$ make kind
```

This creates a kind cluster named `kubeserial`, installs cert-manager, builds the local images, loads them into the cluster and installs the dev release.

### minikube

```bash
$ make minikube
```

This starts a minikube profile named `kubeserial`, builds the controller and monitor images directly into minikube's Docker daemon, installs cert-manager and deploys the dev release.

Use `make minikube-set-context` to point your kubectl context at the minikube cluster.

### Deploying into an existing cluster

If you already have a cluster and just want to (re)install the dev release into it:

```bash
$ make install-dev
```

This bumps the chart versions and installs both the `kubeserial-crds` and `kubeserial` charts into the `kubeserial` namespace, using `charts/kubeserial/values-local.yaml`. Run `make uninstall` to remove the release.

## Working on the docs

The docs are built with [mdBook](https://rust-lang.github.io/mdBook/) and use the `mermaid`, `toc` and `open-on-gh` preprocessors.

```bash
$ make docs-deps    # install mdbook and the required plugins (needs rust + cargo)
$ make docs-serve   # build, serve and open the book in your browser
```
