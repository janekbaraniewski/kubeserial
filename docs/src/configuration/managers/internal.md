# Manager scheduled by KubeSerial

```yaml
apiVersion: app.kubeserial.com/v1alpha1
kind: Manager
metadata:
  annotations:
    meta.helm.sh/release-name: kubeserial
    meta.helm.sh/release-namespace: kubeserial
  creationTimestamp: "2022-06-05T23:48:43Z"
  generation: 1
  labels:
    app.kubernetes.io/instance: kubeserial
    app.kubernetes.io/managed-by: Helm
    app.kubernetes.io/name: kubeserial
    app.kubernetes.io/version: 0.0.1-8c68648
    helm.sh/chart: kubeserial-0.0.1-8c68648
  name: octoprint
  resourceVersion: "43158"
  uid: c908bad7-8716-4f50-aa3f-aeb91bebe71f
spec:
  config: |
    accessControl:
      enabled: false
    plugins:
      announcements:
        _config_version: 1
        channels:
          _blog:
            read_until: 1573642500
          _important:
            read_until: 1521111600
          _octopi:
            read_until: 1573722900
          _plugins:
            read_until: 1573862400
          _releases:
            read_until: 1574699400
      discovery:
        upnpUuid: ef35acc7-a859-4947-980d-d5edb10508e4
      softwareupdate:
        _config_version: 6
      tracking:
        enabled: false
    deviceProfiles:
      default: _default
    serial:
      additionalPorts:
      - /dev/devices/ender3
      autoconnect: true
      baudrate: 0
      port: /dev/device
    server:
      firstRun: false
      onlineCheck:
        enabled: true
      pluginBlacklist:
        enabled: false
      seenWizards:
        corewizard: 3
        cura: null
        tracking: null`
  configPath: /data/config.yaml
  image:
    repository: janekbaraniewski/octoprint
    tag: 1.3.10
  runCmd: mkdir /root/.octoprint && cp /data/config.yaml /root/.octoprint/config.yaml
    && /OctoPrint-1.3.10/run --iknowwhatimdoing --port 80
```
