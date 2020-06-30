DOCKERHUB=janekbaraniewski/kubeserial
VERSION=stable

PHONY: .compile
compile: compilearm

PHONY: .compilearm
compilearm: build/_output/bin/kubeserial

compilearm: export GOOS=linux
compilearm: export GOARCH=arm
compilearm: export GOARM=5

build/_output/bin/kubeserial:
	@mkdir -p build/_output/bin/
	@go build -o build/_output/bin/kubeserial cmd/manager/main.go

PHONY: .build
build: compilearm
	cd build/ && docker build . -t $(DOCKERHUB):$(VERSION) 

PHONY: .clean
clean:
	@rm -rf build/_output