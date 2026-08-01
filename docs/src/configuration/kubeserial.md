# KubeSerial

The `KubeSerial` resource is the single place where you declare which devices the cluster should expose. The Helm chart creates one for you from the `kubeserial.serialDevices` values, and the [KubeSerial controller](../components/controllers.md) turns its `serialDevices` list into the [Device Monitor](../components/monitor.md) and the per-device [`SerialDevice`](devices.md) objects.

```yaml
apiVersion: app.kubeserial.com/v1alpha1
kind: KubeSerial
metadata:
  name: kubeserial
  namespace: kubeserial
spec:
  ingress:
    enabled: false
  serialDevices:
  - idProduct: "6001"
    idVendor: "0403"
    manager: octoprint
    name: ender3
  - idProduct: ea60
    idVendor: 10c4
    name: sonoff-zigbee
```

## Spec

| Field | Required | Description |
| --- | --- | --- |
| `serialDevices` | yes | List of devices to expose. |
| `serialDevices[].name` | yes | Device name. Becomes the `/dev/<name>` udev symlink and the prefix for the generated gateway/manager resources. |
| `serialDevices[].idVendor` | yes | USB vendor id used to match the device. |
| `serialDevices[].idProduct` | yes | USB product id used to match the device. |
| `serialDevices[].manager` | no | Name of a [`Manager`](managers/internal.md) to schedule automatically when this device is connected. Omit it to drive the device with the [injection webhook](managers/external.md) instead. |
| `ingress.enabled` | yes | Reserved for manager ingress. Ingress is not wired up in the current release, leave it `false`. |
| `ingress.domain` | no | Domain used when ingress is enabled. |
| `ingress.annotations` | no | Annotations added to generated ingress objects. |

> Quote `idVendor` / `idProduct` values that are all digits (for example `"6001"`) so they are treated as strings.
