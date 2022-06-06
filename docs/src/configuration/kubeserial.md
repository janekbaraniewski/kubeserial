# KubeSerial

```yaml
apiVersion: app.kubeserial.com/v1alpha1
kind: KubeSerial
metadata:
  annotations:
    ...
  labels:
    ...
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
