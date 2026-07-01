# Device Monitor

Monitors cluster nodes waiting for the configured serial devices to be connected and keeps each [`SerialDevice`](../configuration/devices.md) status in sync with what is actually plugged in.

<!-- toc -->

## How it works

The Device Monitor is deployed by the [KubeSerial controller](controllers.md) as a `DaemonSet`, so one instance runs on every node. Each pod has two containers:

- **`udev-monitor`** mounts the host `/dev` and a generated udev rules file (`98-devices.rules`). For each device in the [KubeSerial](../configuration/kubeserial.md) spec it gets a rule of the form:

  ```
  SUBSYSTEM=="tty", ATTRS{idVendor}=="0403", ATTRS{idProduct}=="6001", SYMLINK+="ender3"
  ```

  This makes udev create a stable `/dev/<device-name>` symlink whenever the matching USB device is connected, regardless of the kernel-assigned `ttyUSB*` name.

- **`device-monitor`** runs the reconcile loop. Once per second it lists all `SerialDevice` resources whose `Ready` condition is `True` and, for each, `stat`s `/dev/<device-name>` on its node:
  - if the device file exists and the device is not yet marked available, it sets `Available=True` and `Free=True` and records its own node name in `status.nodeName`;
  - if the device disappears from the node that owned it, it sets `Available=False`, `Free=Unknown` and clears `status.nodeName`.

Because the monitor writes `status.nodeName`, the rest of the system always knows which node a device is attached to. That is what lets the controller pin the [Device Gateway](gateway.md) to the correct node.

Both containers run privileged and mount the host `/dev` so they can see the physical devices. The `device-monitor` container reads its namespace and node name from the `OPERATOR_NAMESPACE` and `NODE_NAME` environment variables.
