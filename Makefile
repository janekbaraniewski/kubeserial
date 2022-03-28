DOCKERHUB=janekbaraniewski/kubeserial
VERSION=stable

PHONY: .compile
compile: compile-arm

PHONY: .compile-arm
compile-arm: export GOOS=linux
compile-arm: export GOARCH=arm
compile-arm: export GOARM=5
compile-arm: build/_output/bin/kubeserial

build/_output/bin/kubeserial:
	@mkdir -p build/_output/bin/
	go build -o build/_output/bin/kubeserial cmd/manager/main.go

PHONY: .docker-arm
docker-arm: compile-arm
	cd build/ && docker build . -t $(DOCKERHUB):$(VERSION)

PHONY: .clean
clean:
	@rm -rf build/_output