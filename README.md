# KubeSerial
[![License](https://img.shields.io/badge/License-Apache_2.0-blue.svg)](https://opensource.org/licenses/Apache-2.0)
![GH workflow](https://github.com/janekbaraniewski/kubeserial/actions/workflows/test.yml/badge.svg)
[![Go Report Card](https://goreportcard.com/badge/github.com/janekbaraniewski/kubeserial)](https://goreportcard.com/report/github.com/janekbaraniewski/kubeserial)
[![codecov](https://codecov.io/gh/janekbaraniewski/kubeserial/branch/master/graph/badge.svg?token=Y95FB6H188)](https://codecov.io/gh/janekbaraniewski/kubeserial)

---

KubeSerial monitors your cluster nodes for physical devices specified in spec. Once the device is connected, it creates gateway service that exposes it over the network and manager service with specified management software. When the device gets disconnected everything is cleaned up.

### Example:

This is an example configuration for 2 popular cheap printers - Ender 3 and Anet A8 - that handles their management and exposes Octoprint instance for each of them at `http://ender3.my.home` and `http://aneta8.my.home`.

```yaml
apiVersion: app.kubeserial.com/v1alpha1
kind: KubeSerial
metadata:
  name: kubeserial
  namespace: kubeserial
spec:
  devices:
    - name:       "ender3"
      idVendor:   "0403"
      idProduct:  "6001"
      manager:    "octoprint"
      subsystem:  "tty"
    - name:       "aneta8"
      idVendor:   "1a86"
      idProduct:  "7523"
      manager:    "octoprint"
      subsystem:  "tty"
  ingress:
    enabled: true
    domain: ".my.home"
    annotations:
      kubernetes.io/ingress.class: traefik
```

![Example usage 1](demo1.gif)


# Requirements

- k8s cluster
- Ingress controller installed in the cluster for ingress rules to work.

# Install with Helm

## Add help repo


```bash
$ helm repo add baraniewski https://baraniewski.com/charts/
```
## Install CRDs

Due to way in which helm handles CRDs, they are managed using separate chart.

```bash
$ helm upgrade --install kubeserial-crds baraniewski/kubeserial-crds
```
## Install Controller

```bash
$ helm upgrade --install kubeserial baraniewski/kubeserial
```

## Create your configuration

Create your configuration file based on example above. To find out values of `idVendor` and `idProduct` for your device, connect it to your computer, locate where it is (let's say `/dev/ttyUSB0`) and run:
```
udevadm info -q all -n /dev/ttyUSB0 --attribute-walk
```
Look for them from the top. Once you've got your configuration, run

```
kubectl create -f my-kubeserial.yaml
```

# Components

### Controller

Manages operator components by observing state of each of the devices.

### Monitor

Monitors cluster nodes waiting for specified serial devices to be connected and updates their state.

### Gateway

Exposed specific device in cluster network over TCP.

### Manager

Creates deployment with management software, mounts your device over the network and gives you access through ingress rule.

###### Supported managers

- Octoprint


# Build

## Multi-Arch Setup

If running from an x86 system, run the following to install necessary qemu emulation:

```sh
docker run --rm --privileged multiarch/qemu-user-static --reset -p yes
```

See <https://github.com/multiarch/qemu-user-static> for more details.

## To Compile Kubeserial Binary

```sh
make compile
```

## Build The Docker Comtainer

```sh
make build
```

