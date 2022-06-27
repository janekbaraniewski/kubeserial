# KubeSerial

![Release](https://badgen.net/github/release/janekbaraniewski/kubeserial)
[![License](https://img.shields.io/github/license/janekbaraniewski/kubeserial.svg)](LICENSE)
![GH workflow](https://github.com/janekbaraniewski/kubeserial/actions/workflows/test.yml/badge.svg)
[![Go Report Card](https://goreportcard.com/badge/github.com/janekbaraniewski/kubeserial)](https://goreportcard.com/report/github.com/janekbaraniewski/kubeserial)
[![codecov](https://codecov.io/gh/janekbaraniewski/kubeserial/branch/master/graph/badge.svg?token=Y95FB6H188)](https://codecov.io/gh/janekbaraniewski/kubeserial)
[![Go Reference](https://pkg.go.dev/badge/github.com/janekbaraniewski/kubeserial.svg)](https://pkg.go.dev/github.com/janekbaraniewski/kubeserial)
[![CodeQL](https://github.com/janekbaraniewski/kubeserial/workflows/CodeQL/badge.svg)](https://github.com/janekbaraniewski/kubeserial/actions?query=workflow%3ACodeQL)


KubeSerial monitors your cluster nodes for physical devices specified in spec. Once the device is connected, it creates gateway service that exposes it over the network and manager service with specified management software. When the device gets disconnected everything is cleaned up.

![Example usage 1](docs/demo1.gif)

## Quick start

You can find quick start guide [here](https://baraniewski.com/kubeserial/quick_start.html)

## Docs

Docs can be found [here](https://baraniewski.com/kubeserial/)
