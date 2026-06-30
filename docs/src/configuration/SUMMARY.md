# Configuration

KubeSerial is configured through a handful of cluster-scoped custom resources. In a typical Helm install you describe everything in one [`KubeSerial`](kubeserial.md) resource and the controllers derive the rest, but you can also manage the individual objects directly.

- [KubeSerial](kubeserial.md) - the top-level resource listing every device you want to expose.
- [Devices](devices.md) - the `SerialDevice` objects that represent a single device and carry its status.
- [Managers](managers/SUMMARY.md) - how to attach management software or arbitrary workloads to a device.
