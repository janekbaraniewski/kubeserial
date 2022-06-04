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

Now you're ready to create configuration for your devices. Here you can find example of 2 devices - one is using predefined manager, one doesn't. To learn about the difference, please refer to [manager configuration docs](configuration/managers/SUMMARY.md)

Add this config to your values file:


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

## Update your helm release with new values

```bash
$ helm upgrade kubeserial baraniewski/kubeserial -f my-values.yaml
```

## Validate that everything is working

You should see 3 workloads in your cluster - controler manager, webhook and device monitor:

```bash
$ ➜  kubectl get pods                                                                  
NAME                                         READY   STATUS    RESTARTS   AGE
kubeserial-7d58555d4c-97nqk                  1/1     Running   0          5m
kubeserial-device-injector-cc7696b59-mfdn5   1/1     Running   0          5m
kubeserial-monitor-rjs82                     2/2     Running   0          5m
```

Number of `kubeserial-monitor` pods should match number of your cluster nodes.

You should also see your devices:

```bash
$ ➜  kubectl get serialdevice
NAME            READY   AVAILABLE   NODE
ender3          True    False       
sonoff-zigbee   True    False       
```

You're all set with basic configuration of KubeSerial. Please read through rest of docs to configure it to your needs.
