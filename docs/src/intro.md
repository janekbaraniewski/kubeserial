# KubeSerial 
<b>Work with serial devices in Kubernetes</b>

---

KubeSerial is a set of Kubernetes controllers and resources that make working with serial devices in Kubernetes clusters easy. It decouples node that device is connected to from workload that's using it.

## How it works

KubeSerial uses [Device Monitors][DM] to monitor all cluster nodes for [Devices][D]. Once [Device][D] is detected, it will schedule [Device Gateway][DG] to specific cluster node. [Device Gateway][DG] exposes this device as TCP server. Once [Device Gateway][DG] is available, it can be used by [Manager][M]. There are 2 supported modes for [Managers][M]:
- Use predefined Manager that will be scheduled when Device is available
- Use annotation to inject device to any workload


<!-- Links  -->
[D]:  /configuration/devices.md            "Device"
[DM]: /components/monitor.md               "Device Monitor"
[DG]: /components/gateway.md               "Device Gateway"
[M]:  /configuration/managers/SUMMARY.md   "Managers"
