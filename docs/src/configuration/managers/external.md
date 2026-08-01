# Manager scheduled externally

Instead of defining a `Manager` resource, you can attach a device to any pod you already run by adding a single annotation. The KubeSerial mutating webhook rewrites the pod at creation time so its container reaches the device over the network. This lets you use a device from an off-the-shelf image without rebuilding it or running a privileged pod yourself.

<!-- toc -->

## Usage

Add the annotation to the pod template, naming the [`SerialDevice`](../devices.md) you want:

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: zigbee2mqtt
spec:
  template:
    metadata:
      annotations:
        app.kubeserial.com/inject-device: sonoff-zigbee
    spec:
      containers:
      - name: zigbee2mqtt
        image: koenkk/zigbee2mqtt:latest
```

The device name must match a `SerialDevice` (`sonoff-zigbee` in the [Quick Start](../../quick_start.md) example). Inside the container the device is available as `/dev/device`, so point your software at that path.

## How it works

The webhook is a `MutatingWebhookConfiguration` that intercepts pod `create` and `update` calls. For each pod it:

1. Reads the `app.kubeserial.com/inject-device` annotation. If it is missing, the pod is left untouched.
2. Looks up the named `SerialDevice` and checks that its `Free` condition is `True`. If the device does not exist or is not free, the pod is admitted unchanged (no bridge is injected).
3. Rewrites the first container's `command`/`args` to start a `socat` bridge before the original process:

   ```sh
   socat -d -d pty,raw,echo=0,b115200,link=/dev/device,perm=0660,group=tty tcp:sonoff-zigbee-gateway:3333 & <original entrypoint>
   ```

   `socat` creates a local PTY at `/dev/device` and connects it to the device's [gateway](../../components/gateway.md) service. The original entrypoint then runs as usual and talks to `/dev/device`.

If the container does not set an explicit `command`, the webhook reads the image's OCI config to recover the image's entrypoint and command, so it can wrap them correctly. It also sets an `app.kubeserial.com/device-injected` annotation so the same pod is not rewritten twice.

## Requirements and limitations

- The device's gateway must be running, which means the `SerialDevice` has to be `Available` and `Free`. If it is not free at the moment the pod is created, the pod starts without the bridge.
- Only the first container in the pod is rewritten.
- The webhook needs its serving certificate; cert-manager and a configured issuer are part of the [install requirements](../../quick_start.md).
