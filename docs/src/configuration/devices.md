# Device

<!-- toc -->

## Spec

```yaml
apiVersion: app.kubeserial.com/v1alpha1
kind: SerialDevice
metadata:
  annotations:
    ...
  labels:
    ...
  name: ender3
spec:
  idProduct: "6001"
  idVendor: "0403"
  manager: octoprint
  name: ender3
```

## Status conditions

```yaml
apiVersion: app.kubeserial.com/v1alpha1
kind: SerialDevice
...
status:
  conditions:
  - lastHeartbeatTime: "2022-06-05T23:49:02Z"
    lastTransitionTime: "2022-06-05T23:49:02Z"
    message: ""
    reason: NotValidated
    status: "False"
    type: Available
  - lastHeartbeatTime: "2022-06-05T23:49:02Z"
    lastTransitionTime: "2022-06-05T23:49:02Z"
    message: ""
    reason: AllChecksPassed
    status: "True"
    type: Ready
  - lastHeartbeatTime: "2022-06-05T23:49:02Z"
    lastTransitionTime: "2022-06-05T23:49:02Z"
    message: ""
    reason: NotValidated
    status: "False"
    type: Free
```
