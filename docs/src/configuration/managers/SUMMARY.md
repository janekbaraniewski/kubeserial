# Managers

> To learn more about the concept of a Manager take a look at the [Manager docs][MD].

Configuring a manager is not required, but it is the easiest way to actually use a device once KubeSerial has exposed it. There are two ways to attach something to a device. You can use both at the same time for different devices, but not for the same device.

## Manager scheduled by KubeSerial

You create a `Manager` resource holding the spec of the management software you want to run and bind it to a `SerialDevice` (via `spec.manager`). When KubeSerial detects that the device is connected it schedules the management software for you and bridges the device into the container. When the device is disconnected everything is cleaned up. Learn how to configure it by [reading the docs][MIC].

## Manager scheduled externally

You add an annotation to your Pod and the KubeSerial mutating webhook rewrites the pod when it is created to bridge in the device. This works with any image and needs no `Manager` resource. Learn how to configure it by [reading the docs][MEC].

<!-- Links  -->
[MD]:  ../../components/manager.md   "Manager"
[MIC]: internal.md                  "Manager internal"
[MEC]: external.md                  "Manager external"
