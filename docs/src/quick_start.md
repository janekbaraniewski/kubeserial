# Quick Start

<!-- toc -->

## Requirements

- k8s cluster
- CertManager installed and Cert Issuer configured

## Install with Helm

### Add help repo

First you'll need to add helm repository that stores KubeSerial charts

```bash
$ helm repo add baraniewski https://baraniewski.com/charts/
```

### Install CRDs

Due to way in which helm handles CRDs, they are managed using separate chart.

```bash
$ helm upgrade --install kubeserial-crds baraniewski/kubeserial-crds
```

### Create minimal values file

In order to make webhook work, you'll need to specify which Cert Issuer should be used for SSL cert used by webhook. To do this, create values file with following structure (change values depending on your setup):

```yaml
certManagerIssuer:
  name: selfsigned-issuer
  kind: ClusterIssuer
```

<mark>Set proper name and kind of your Issuer or ClusterIssuer</mark>

### Install Controller

```bash
$ helm upgrade --install kubeserial baraniewski/kubeserial
```

## Get device attributes

To find out values of `idVendor` and `idProduct` for your device, connect it to your computer, locate where it is (let's say `/dev/ttyUSB0`) and run:

```bash
udevadm info -q all -n /dev/ttyUSB0 --attribute-walk
```

Look for them from the top.

## Update helm values file

TODO: more details

```yaml
kubeserial:
  serialDevices:
  - idProduct: "6001"
    idVendor: "0403"
    manager: octoprint
    name: ender3
    subsystem: tty
  - idProduct: ea60
    idVendor: 10c4
    name: sonoff-zigbee
    subsystem: tty
```
