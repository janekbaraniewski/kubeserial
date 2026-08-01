# Manager

A Manager is the management software that talks to a device (for example OctoPrint for a 3D printer, or Zigbee2MQTT for a Zigbee dongle). KubeSerial can schedule a Manager for you automatically whenever its device is connected, and tear it down again when the device goes away.

<!-- toc -->

## How it works

You declare a [`Manager`](../configuration/managers/internal.md) resource describing the image, the run command and optional config, and reference it from a `SerialDevice` (`spec.manager`). When that device becomes available the [controllers](controllers.md) create a `ManagerScheduleRequest`, and the ManagerScheduleRequest controller schedules the workload:

- a **Deployment** running the manager image,
- a **`ClusterIP` Service** on port `80`,
- a **ConfigMap** with the rendered config, mounted at the manager's `configPath` (only when `spec.config` is set).

The deployment's entrypoint is wrapped so that, before the manager starts, a `socat` bridge connects the [Device Gateway](gateway.md) to a local PTY:

```sh
socat -d -d pty,raw,echo=0,b115200,link=/dev/device,perm=0660,group=tty tcp:<device>-gateway:3333 & <runCmd>
```

The manager then talks to `/dev/device` as if the hardware were attached locally, while the bytes actually travel over the network to the node holding the device.

> Ingress for managers is not wired up in the current release; reach managers through their Service.

## Two ways to attach a device

Running a predefined Manager is one of two supported modes. The other is the [device-injection webhook](../configuration/managers/external.md), which applies the same `socat` bridge to *any* pod via an annotation, without you defining a Manager resource. See the [Managers](../configuration/managers/SUMMARY.md) overview for when to use each.
