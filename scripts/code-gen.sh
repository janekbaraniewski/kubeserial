#!/usr/bin/env bash

set -o errexit
set -o nounset
set -o pipefail

cd "$(git rev-parse --show-toplevel)"

: ${GO=go}

export GOPATH=$("$GO" env GOPATH | awk -F ':' '{print $1}')
export PATH=$PATH:$GOPATH/bin

exit_on_error () {
    printf "\n${1}\n\n"
    exit 1
}

replace_or_compare () {
    echo "Path1 -> ${1}"
    echo "Path2 -> ${2}"
    echo "PWD -> ${PWD}"

    ls ${1}
    ls ${2}

    printf "\n\n\n file 1 \n\n\n"
    cat ${1}
    printf "\n\n\n file 2 \n\n\n"
    cat ${2}

    if [[ "${COPY_OR_DIFF}" == "copy" ]]; then
        cp -r $1 $2
    elif [[ "${COPY_OR_DIFF}" == "diff" ]]; then
        diff -r $1 $2 || exit_on_error "To fix run:\n    make codegen"
    fi
}

# Install the required binaries
GO111MODULE=on "$GO" install \
    k8s.io/code-generator/cmd/deepcopy-gen \
    k8s.io/code-generator/cmd/register-gen \
    k8s.io/code-generator/cmd/client-gen \
    k8s.io/code-generator/cmd/lister-gen \
    k8s.io/code-generator/cmd/informer-gen \
    k8s.io/code-generator/cmd/openapi-gen \
    sigs.k8s.io/controller-tools/cmd/controller-gen

printf "deepcopy-gen... "

deepcopy-gen \
  --go-header-file "scripts/boilerplate.go.txt" \
  --input-dirs="github.com/janekbaraniewski/kubeserial/pkg/apis/app/v1alpha1" \
  --output-package="github.com/janekbaraniewski/kubeserial/pkg/apis/app/v1alpha1" \
  --output-file-base=zz_generated.deepcopy -v 1

replace_or_compare $GOPATH/src/github.com/janekbaraniewski/kubeserial/pkg/apis/app/v1alpha1/zz_generated.deepcopy.go ./pkg/apis/app/v1alpha1/zz_generated.deepcopy.go

printf "Done!\n"

printf "register-gen... "

register-gen all \
  --go-header-file "scripts/boilerplate.go.txt" \
  --input-dirs="github.com/janekbaraniewski/kubeserial/pkg/apis/app/v1alpha1" \
  --output-package="github.com/janekbaraniewski/kubeserial/pkg/apis/app/v1alpha1" \
  --output-file-base=zz_generated.register -v 1

replace_or_compare $GOPATH/src/github.com/janekbaraniewski/kubeserial/pkg/apis/app/v1alpha1/zz_generated.register.go ./pkg/apis/app/v1alpha1/zz_generated.register.go

printf "Done!\n"

printf "openapi-gen... "

openapi-gen \
  --go-header-file "scripts/boilerplate.go.txt" \
  --input-dirs="github.com/janekbaraniewski/kubeserial/pkg/apis/app/v1alpha1" \
  --output-package="github.com/janekbaraniewski/kubeserial/pkg/apis/app/v1alpha1" \
  --output-file-base=zz_generated.openapi -v 1

replace_or_compare $GOPATH/src/github.com/janekbaraniewski/kubeserial/pkg/apis/app/v1alpha1/zz_generated.openapi.go ./pkg/apis/app/v1alpha1/zz_generated.openapi.go

printf "Done!\n"

printf "client-gen... "

client-gen \
    --go-header-file "scripts/boilerplate.go.txt" \
    --input-base="" \
    --input="github.com/janekbaraniewski/kubeserial/pkg/apis/app/v1alpha1" \
    --output-package=github.com/janekbaraniewski/kubeserial/pkg/generated/clientset \
    --clientset-name=versioned

printf "Done!\n"

printf "lister-gen... "

lister-gen \
    --go-header-file "scripts/boilerplate.go.txt" \
    --input-dirs="github.com/janekbaraniewski/kubeserial/pkg/apis/app/v1alpha1" \
    --output-package=github.com/janekbaraniewski/kubeserial/pkg/generated/listers
printf "Done!\n"

printf "informer-gen... "

informer-gen \
    --go-header-file "scripts/boilerplate.go.txt" \
    --input-dirs="github.com/janekbaraniewski/kubeserial/pkg/apis/app/v1alpha1" \
    --versioned-clientset-package=github.com/janekbaraniewski/kubeserial/pkg/generated/clientset/versioned \
    --listers-package=github.com/janekbaraniewski/kubeserial/pkg/generated/listers \
    --output-package=github.com/janekbaraniewski/kubeserial/pkg/generated/informers

rm -rf ./pkg/generated || true
replace_or_compare $GOPATH/src/github.com/janekbaraniewski/kubeserial/pkg/generated/ ./pkg/generated

printf "Done!\n"

printf "controller-gen... "

controller-gen crd paths=./pkg/apis/app/... output:crd:dir=/tmp/deploy/crds
replace_or_compare /tmp/deploy/crds deploy/

printf "Done!\n"
