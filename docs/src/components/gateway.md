# Device Gateway

Exposes a specific device on the cluster network over TCP.

<!-- toc -->

## How it works

When a [`SerialDevice`](../configuration/devices.md) becomes `Available`, the [SerialDevice controller](controllers.md) creates a Device Gateway for it. A gateway is three resources, all named `<device>-gateway`:

- a **ConfigMap** holding a `ser2net.conf`,
- a **Deployment** running `ser2net`,
- a **`ClusterIP` Service** publishing port `3333`.

The `ser2net` container is privileged and mounts the host `/dev`. Its config maps the device to a raw TCP port:

```
3333:raw:600:/dev/ender3:115200 8DATABITS NONE 1STOPBIT -XONXOFF LOCAL -RTSCTS HANGUP_WHEN_DONE
```

The deployment is pinned to the node the device is attached to using a `kubernetes.io/hostname` node selector set from the device's `status.nodeName` (written by the [Device Monitor](monitor.md)). This is why the device only needs to be physically present on one node: the gateway always lands on that node, and everything else reaches it through the service.

Consumers connect to `<device>-gateway:3333` and turn that TCP stream back into a local serial port with `socat` (see the [Manager](manager.md) and the [device-injection webhook](../configuration/managers/external.md)). When the device is unplugged, the controller deletes the ConfigMap, Deployment and Service again.
