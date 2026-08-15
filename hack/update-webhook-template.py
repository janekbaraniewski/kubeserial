#!/usr/env python3

import sys
import yaml

if len(sys.argv) != 3:
    print("Wrong number of arguments")
    exit

crd_path = sys.argv[1]
metadata_template_path = sys.argv[2]

with open(crd_path) as f:
    crd_yaml = yaml.safe_load(f)

with open(metadata_template_path) as f:
    metadata_template = yaml.safe_load(f)

crd_yaml['metadata'] = metadata_template['metadata']
crd_yaml['webhooks'][0]['clientConfig'] = metadata_template['clientConfig']

# Scoping/failure behaviour are chart-configurable, so they come from the
# template rather than from the controller-gen marker.
for key in ('failurePolicy', 'namespaceSelector', 'objectSelector'):
    if key in metadata_template:
        crd_yaml['webhooks'][0][key] = metadata_template[key]

with open(crd_path, 'w') as f:
    yaml.dump(crd_yaml, f, default_style=None, default_flow_style=False, width=4096)
