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
|image.repository|`ghcr.io/janekbaraniewski/kubeserial`||true|
|image.pullPolicy|`IfNotPresent`||true|
|image.tag|`APP_VERSION`||true|
|monitor.image.repository|`ghcr.io/janekbaraniewski/kubeserial-device-monitor`||true|
|monitor.image.pullPolicy|`IfNotPresent`||true|
|monitor.image.tag|`APP_VERSION`||true|
|monitor.resources|`{}`||true|
|webhook.image.repository|`ghcr.io/janekbaraniewski/kubeserial-injector-webhook`||true|
|webhook.image.pullPolicy|`IfNotPresent`||true|
|webhook.image.tag|`APP_VERSION`||true|
|webhook.failurePolicy|`Ignore`|What happens to a pod creation when the injector is unreachable. `Ignore` skips injection; `Fail` rejects the pod.|false|
|webhook.namespaceSelector|excludes `kube-system`, `kube-node-lease`, `kube-public`, `cert-manager`|Namespaces the injector is consulted for.|false|
|webhook.objectSelector|`{}`|Extra pod-label narrowing. Cannot match annotations, so it cannot key off `app.kubeserial.com/inject-device`.|false|
|webhook.serviceAccount.create|`true`|Create a dedicated service account for the injector, holding only read access to SerialDevices.|false|
|webhook.serviceAccount.name|`""`|Name of the injector service account. Defaults to `<release>-device-injector`.|false|
|webhook.serviceAccount.annotations|`{}`||false|
|monitoring.prometheusMonitors.enabled|`true`||true|


[comment]: # (Links)
[kubeserial]: https://github.com/janekbaraniewski/kubeserial
