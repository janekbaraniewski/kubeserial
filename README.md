# KubeSerial

Manage your serial devices like 3D printers or CNC machines on your k8s cluster.

---

# Components

### controller

Manages operator components by observing state of each of the devices.

### monitor

Monitors cluster nodes waiting for specified serial devices to be connected and updates their state.

### gateway

Reacts to state changes and exposes your device in cluster network over TCP.

### manager

Creates deployment with management software, mounts your device over the network and gives you access through ingress rule.

# Requirements

- k8s cluster - ATM only ARM clusters are supported.
- Ingress controller installed in the cluster for ingress rules to work.

# Install

### Install manually 

Create the Deployment, CRDs, ServiceAccount etc.

```
kubectl create -f deploy/kubeserial.yaml
```

Configure

> Example configuration for Ender3 3D printer:

```yaml
# my-kubeserial.yaml
apiVersion: app.kubeserial.com/v1alpha1
kind: KubeSerial
metadata:
  name: kubeserial
  namespace: kubeserial
spec:
  devices:
    - name:       "ender3"
      idvendor:   "0403"
      idproduct:  "6001"
      manager:    "octoprint"
```

```
kubectl create -f my-kubeserial.yaml
```

