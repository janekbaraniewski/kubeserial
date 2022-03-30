DOCKERHUB=janekbaraniewski/kubeserial
TARGET_PLATFORMS=$(shell cat TARGET_PLATFORMS)
VERSION ?= $(shell git rev-parse --short HEAD)
DOCKERBUILD_EXTRA_OPTS ?=
DOCKERBUILD_PLATFORM_OPT=--platform
GO_BUILD_OUTPUT_PATH ?= build/_output/bin/kubeserial

PHONY: .all
all: kubeserial

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
kubeserial-docker-all: PLATFORMS=${TARGET_PLATFORMS}
kubeserial-docker-all: DOCKERBUILD_ACTION=--push
kubeserial-docker-all: kubeserial-docker

kubeserial-docker: DOCKERFILE=Dockerfile
kubeserial-docker: docker-build

docker-build:
	docker buildx build . -f ${DOCKERFILE} ${DOCKERBUILD_EXTRA_OPTS} ${DOCKERBUILD_PLATFORM_OPT} ${PLATFORMS} -t $(DOCKERHUB):$(VERSION) ${DOCKERBUILD_ACTION}

PHONY: .codegen
codegen: COPY_OR_DIFF=copy
codegen: code-gen.sh

PHONY: .check-codegen
check-codegen: COPY_OR_DIFF=diff
check-codegen: code-gen.sh

code-gen.sh:
	@COPY_OR_DIFF=${COPY_OR_DIFF} ./scripts/code-gen.sh

PHONY: .test
test:
