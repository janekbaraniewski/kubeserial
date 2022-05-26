#!/usr/bin/env bash

set -o errexit
set -o nounset
set -o pipefail

sed="sed"

unameOut="$(uname -s)"
case "${unameOut}" in
    Linux*)     sed=sed;;
    Darwin*)    sed=gsed;;
    *)          sed="UNKNOWN:${unameOut}"
esac

cd "$(git rev-parse --show-toplevel)"

: ${GO=go}

export GOPATH=$("$GO" env GOPATH | awk -F ':' '{print $1}')
export PATH=$PATH:$GOPATH/bin

source hack/utils.sh

GO111MODULE=on "$GO" install \
    sigs.k8s.io/controller-tools/cmd/controller-gen

printf "controller-gen CRD... "

controller-gen rbac:roleName=manager-role crd paths=./pkg/apis/... output:crd:dir=/tmp/deploy/crds
find /tmp/deploy/crds -name "*.yaml" | xargs -I % python3 ./hack/update-crd-metadata.py % ./hack/crd_metadata_template.yaml
find /tmp/deploy/crds -name "*.yaml" | xargs ${sed} -i 's/\x27{{/{{/g'  # change '{{ -> {{
find /tmp/deploy/crds -name "*.yaml" | xargs ${sed} -i 's/}}\x27/}}/g' # change }}' -> }}
find /tmp/deploy/crds -name "*.yaml" | xargs ${sed} -i 's/\\\x27//g' # change \' -> '
replace_or_compare /tmp/deploy/crds/ deploy/chart/kubeserial-crds/templates/
rm -r /tmp/deploy/crds

printf "Done!\n"

printf "controller-gen webhook... "

controller-gen rbac:roleName=manager-role webhook paths=./pkg/webhooks/... output:webhook:dir=/tmp/webhooks
find /tmp/webhooks -name "*.yaml" | xargs -I % python3 ./hack/update-webhook-template.py % ./hack/webhook_template.yaml
find /tmp/webhooks -name "*.yaml" | xargs ${sed} -i 's/\x27{{/{{/g'  # change '{{ -> {{
find /tmp/webhooks -name "*.yaml" | xargs ${sed} -i 's/}}\x27/}}/g' # change }}' -> }}
find /tmp/webhooks -name "*.yaml" | xargs ${sed} -i 's/\\\x27//g' # change \' -> '
replace_or_compare /tmp/webhooks/manifests.yaml deploy/chart/kubeserial/templates/device_injector_webhook_webhooks.yaml
rm -r /tmp/webhooks

printf "Done!\n"
