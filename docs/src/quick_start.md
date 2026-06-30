# Quick Start

<!-- toc -->

## Requirements

- a Kubernetes cluster
- [cert-manager](https://cert-manager.io/) installed, with an Issuer or ClusterIssuer configured (the device-injection webhook needs it for its serving certificate)

## Install with Helm

### Add the Helm repo

First add the Helm repository that hosts the KubeSerial charts:

```bash
> helm repo add baraniewski https://baraniewski.com/charts/
```

### Install CRDs

Because of the way Helm handles CRDs, they are shipped as a separate chart:

```bash
> helm upgrade --install kubeserial-crds baraniewski/kubeserial-crds
```

### Create a minimal values file

For the webhook to work you must tell it which issuer to use for its TLS certificate. Create a values file like this (adjust to your setup):

```yaml
certManagerIssuer:
  name: selfsigned-issuer
  kind: ClusterIssuer
```

<mark>Set the proper name and kind of your Issuer or ClusterIssuer</mark>

### Install the controllers

```bash
> helm upgrade --install kubeserial baraniewski/kubeserial -f my-values.yaml
```

## Get device attributes

To find the `idVendor` and `idProduct` for your device, connect it to a machine, find its node (say `/dev/ttyUSB0`) and run:

```bash
> udevadm info -q all -n /dev/ttyUSB0 --attribute-walk
```

Look for them near the top of the output.

## Add your devices

Now declare your devices. The example below has two: one is bound to a predefined manager, one is not. To learn about the difference, see the [manager configuration docs](configuration/managers/SUMMARY.md).

Add this to your values file:

```yaml
kubeserial:
  serialDevices:
  - idProduct: "6001"
    idVendor: "0403"
    manager: octoprint
    name: ender3
  - idProduct: ea60
    idVendor: 10c4
    name: sonoff-zigbee
```

> Quote values that are all digits (for example `"6001"`) so they stay strings.

## Update your release with the new values

```bash
> helm upgrade kubeserial baraniewski/kubeserial -f my-values.yaml
```

## Validate that everything is working

You should see three workloads in your cluster - the controller manager, the device-injection webhook and the device monitor:

```bash
$ kubectl get pods
NAME                                         READY   STATUS    RESTARTS   AGE
kubeserial-7d58555d4c-97nqk                  1/1     Running   0          5m
kubeserial-device-injector-cc7696b59-mfdn5   1/1     Running   0          5m
kubeserial-monitor-rjs82                     2/2     Running   0          5m
```

The number of `kubeserial-monitor` pods should match the number of nodes in your cluster.

You should also see your devices:

```bash
$ kubectl get serialdevice
NAME            READY   AVAILABLE   NODE
ender3          True    True        worker-1
sonoff-zigbee   True    True        worker-2
```

A device shows `AVAILABLE=True` and a `NODE` once the monitor detects it plugged in. From here you can attach workloads to it: either through a predefined [Manager](configuration/managers/internal.md), or by annotating any pod for the [device-injection webhook](configuration/managers/external.md).

You're all set with the basic configuration of KubeSerial. Read through the rest of the docs to configure it to your needs.
