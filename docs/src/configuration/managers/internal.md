# Manager scheduled by KubeSerial

A `Manager` describes a piece of management software to run for a device. Reference it from a `SerialDevice` (`spec.manager`) and KubeSerial schedules it automatically while the device is connected. See the [Manager component docs](../../components/manager.md) for how the workload is built and bridged to the device.

<!-- toc -->

## Spec

| Field | Required | Description |
| --- | --- | --- |
| `image.repository` | yes | Manager image repository. |
| `image.tag` | yes | Manager image tag. |
| `runCmd` | yes | Command used to start the software. It is run after the `socat` bridge is established, so the device is already present at `/dev/device`. |
| `config` | no | Inline configuration file content. When set, it is written to a ConfigMap and mounted at `configPath`. |
| `configPath` | no | Path inside the container where `config` is mounted (for example `/data/config.yaml`). |

`Manager` is cluster-scoped.

## Example

This is an OctoPrint manager for the `ender3` printer used in the [Quick Start](../../quick_start.md):

```yaml
apiVersion: app.kubeserial.com/v1alpha1
kind: Manager
metadata:
  name: octoprint
  namespace: kubeserial
spec:
  image:
    repository: janekbaraniewski/octoprint
    tag: 1.3.10
  configPath: /data/config.yaml
  config: |
    accessControl:
      enabled: false
    serial:
      additionalPorts:
      - /dev/device
      autoconnect: true
      baudrate: 0
      port: /dev/device
    server:
      firstRun: false
  runCmd: mkdir /root/.octoprint && cp /data/config.yaml /root/.octoprint/config.yaml && /OctoPrint-1.3.10/run --iknowwhatimdoing --port 80
```

The device shows up inside the manager container as `/dev/device`, so the manager's own configuration should point at that path.
