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

with open(crd_path, 'w') as f:
    yaml.dump(crd_yaml, f, default_style=None, default_flow_style=False, width=4096)
