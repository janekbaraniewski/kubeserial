DOCKERHUB=janekbaraniewski/kubeserial
SUPPORTED_PLATFORMS=linux/amd64,linux/386,linux/arm64,linux/arm/v7
VERSION ?= $(shell git rev-parse --short HEAD)
DOCKERBUILD_EXTRA_OPTS ?=
DOCKERBUILD_PLATFORM_OPT=--platform
GO_BUILD_OUTPUT_PATH ?= build/_output/bin/kubeserial

PHONY: .kubeserial
kubeserial:
	@mkdir -p build/_output/bin/
	go build -o ${GO_BUILD_OUTPUT_PATH} cmd/manager/main.go

PHONY: .kubeserial-docker-local
kubeserial-docker-local: PLATFORMS=
kubeserial-docker-local: DOCKERBUILD_PLATFORM_OPT=
kubeserial-docker-local: DOCKERBUILD_ACTION=--load
kubeserial-docker-local: VERSION=local
kubeserial-docker-local: kubeserial-docker

PHONY: .kubeserial-docker-all
kubeserial-docker-all: PLATFORMS=${SUPPORTED_PLATFORMS}
kubeserial-docker-all: DOCKERBUILD_ACTION=--push
kubeserial-docker-all: kubeserial-docker

kubeserial-docker: DOCKERFILE=Dockerfile
kubeserial-docker: docker-build

docker-build:
	docker buildx build . -f ${DOCKERFILE} ${DOCKERBUILD_EXTRA_OPTS} ${DOCKERBUILD_PLATFORM_OPT} ${PLATFORMS} -t $(DOCKERHUB):$(VERSION) ${DOCKERBUILD_ACTION}

PHONY: .clean
clean:
	@rm -rf build/_output
