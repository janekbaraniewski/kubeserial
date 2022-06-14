#! /bin/bash

mkdir -p /tmp/kubeserial/test-assets
helm template --name-template kubeserial --show-only templates/config_map.yaml charts/kubeserial > /tmp/kubeserial/test-assets/config-map.yaml
rm -r test-assets && mkdir test-assets
ls charts/kubeserial/specs | xargs -I{} bash -c "yq eval '.data[\"{}\"]' /tmp/kubeserial/test-assets/config-map.yaml > test-assets/{}"
rm -r /tmp/kubeserial/test-assets
