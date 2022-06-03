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
