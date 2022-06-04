# kubeserial

App chart for [KubeSerial][kubeserial]

## Install

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



## Configuration

|name|default|description|required|
|---|---|---|---|
|kubeserial.devices|`[]`|List of device configuration to create when installing|false|
|certManagerIssuer.name|||true|
|certManagerIssuer.kind|`Issuer`||true|
|image.repository|`janekbaraniewski/kubeserial`||true|
|image.pullPolicy|`IfNotPresent`||true|
|image.tag|`APP_VERSION`||true|
|monitor.image.repository|`janekbaraniewski/kubeserial-device-monitor`||true|
|monitor.image.pullPolicy|`IfNotPresent`||true|
|monitor.image.tag|`APP_VERSION`||true|
|monitor.resources|`{}`||true|
|webhook.image.repository|`janekbaraniewski/kubeserial-injector-webhook`||true|
|webhook.image.pullPolicy|`IfNotPresent`||true|
|webhook.image.tag|`APP_VERSION`||true|
|monitoring.prometheusMonitors.enabled|`true`||true|


[comment]: # (Links)
[kubeserial]: https://github.com/janekbaraniewski/kubeserial
