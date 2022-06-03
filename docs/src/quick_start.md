# Quick Start

<!-- toc -->

## Requirements

- k8s cluster
- Ingress controller installed in the cluster for ingress rules to work.

## Install with Helm

### Add help repo

```bash
$ helm repo add baraniewski https://baraniewski.com/charts/
```

### Install CRDs

Due to way in which helm handles CRDs, they are managed using separate chart.

```bash
$ helm upgrade --install kubeserial-crds baraniewski/kubeserial-crds
```

### Install Controller

```bash
$ helm upgrade --install kubeserial baraniewski/kubeserial
```

### Create your config

Create your configuration file based on example above. To find out values of `idVendor` and `idProduct` for your device, connect it to your computer, locate where it is (let's say `/dev/ttyUSB0`) and run:

```bash
udevadm info -q all -n /dev/ttyUSB0 --attribute-walk
```

Look for them from the top. Once you've got your configuration, run
