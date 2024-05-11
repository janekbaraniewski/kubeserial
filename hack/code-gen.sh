#!/usr/bin/env bash

set -o errexit
set -o nounset
set -o pipefail

cd "$(git rev-parse --show-toplevel)"

: ${GO=go}

export GOPATH=$("$GO" env GOPATH | awk -F ':' '{print $1}')
export PATH=$PATH:$GOPATH/bin

source hack/utils.sh

printf "Start code-gen script \n"

# Install the required binaries with modules enabled
go install k8s.io/code-generator
chmod +x "$GOPATH"/pkg/mod/k8s.io/code-generator@v0.30.0/kube_codegen.sh

printf "Got all dependencies \n"

# Ensure proper package tags are in place (See: https://pkg.go.dev/k8s.io/code-generator/cmd/deepcopy-gen)
printf "Running code generators...\n"

# Use generate-groups.sh helper script to run all code generators
"$GOPATH"/pkg/mod/k8s.io/code-generator@v0.30.0/kube_codegen.sh all \
    github.com/janekbaraniewski/kubeserial/pkg/generated \
    github.com/janekbaraniewski/kubeserial/pkg/apis \
    v1alpha1

printf "Finished kube_codegen.sh, updating source files...\n"

# Manual copy might be required if output paths are incorrectly set by the tools, adjust paths as needed
if [[ "${COPY_OR_DIFF}" == "copy" ]]; then
    printf "Removing old generated files...\n"
    rm -rf ./pkg/generated
    printf "Populating with new generated files...\n"
    mkdir -p ./pkg/generated
    cp -r "$GOPATH/src/github.com/janekbaraniewski/kubeserial/pkg/generated/"* ./pkg/generated/
    printf "Populating DONE\n"
    printf "New files populated with: \n"
    ls -la ./pkg/generated/
fi

printf "Final checks...\n"
replace_or_compare "$GOPATH/src/github.com/janekbaraniewski/kubeserial/pkg/generated/" ./pkg/generated/

printf "All generators have completed.\n"
