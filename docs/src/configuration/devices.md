# Device

A `SerialDevice` represents one physical device. You normally don't create these by hand: the [KubeSerial controller](../components/controllers.md) generates one per entry in the [KubeSerial](kubeserial.md) `serialDevices` list. The [Device Monitor](../components/monitor.md) and the controllers keep its `status` in sync with reality.

<!-- toc -->

## Spec

```yaml
apiVersion: app.kubeserial.com/v1alpha1
kind: SerialDevice
metadata:
  name: ender3
spec:
  name: ender3
  idVendor: "0403"
  idProduct: "6001"
  manager: octoprint
```

| Field | Required | Description |
| --- | --- | --- |
| `name` | yes | Device name. Matches the `/dev/<name>` symlink created by the udev rule. |
| `idVendor` | yes | USB vendor id. |
| `idProduct` | yes | USB product id. |
| `manager` | no | Name of a [`Manager`](managers/internal.md) to schedule when the device is available. |

`SerialDevice` is cluster-scoped. `kubectl get serialdevice` prints the `Ready`, `Available` and `Node` columns.

## Status conditions

The status carries three conditions and, once the device is detected, the node it is attached to:

| Condition | Set by | Meaning |
| --- | --- | --- |
| `Ready` | SerialDevice controller | Configuration has been validated. If the device references a manager, that `Manager` object must exist; otherwise the device is ready immediately. |
| `Available` | Device Monitor | The device file is currently present on a node. `status.nodeName` records which one. |
| `Free` | Device Monitor / webhook | The device is available and not in use. The [injection webhook](managers/external.md) only injects into a pod when `Free` is `True`. |

```yaml
apiVersion: app.kubeserial.com/v1alpha1
kind: SerialDevice
...
status:
  nodeName: worker-1
  conditions:
  - type: Ready
    status: "True"
    reason: AllChecksPassed
    lastHeartbeatTime: "2022-06-05T23:49:02Z"
    lastTransitionTime: "2022-06-05T23:49:02Z"
    message: ""
  - type: Available
    status: "True"
    reason: DeviceAvailable
    lastHeartbeatTime: "2022-06-05T23:49:02Z"
    lastTransitionTime: "2022-06-05T23:49:02Z"
    message: ""
  - type: Free
    status: "True"
    reason: DeviceFree
    lastHeartbeatTime: "2022-06-05T23:49:02Z"
    lastTransitionTime: "2022-06-05T23:49:02Z"
    message: ""
```

When the device is unplugged the monitor sets `Available` to `False`, `Free` to `Unknown` and clears `status.nodeName`, and the controller removes the [gateway](../components/gateway.md) and any manager scheduled for it.
