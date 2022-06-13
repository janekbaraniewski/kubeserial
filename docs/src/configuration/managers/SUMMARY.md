# Managers

> In order to learn more about concept of Manager take a look at [Manager docs][MD]

Manager configuration is not required, although it's the easiest way to work with serial devices using KubeSerial.

There are 2 ways to work with managers. You can use them both at the same time for different devices, but not for the same one.

## Manager scheduled by KubeSerial

In this approach, you'll need to create `Manager` type object which will hold spec of manager software you want to run and bind it with `SerialDevice`. Then, when KubeSerial detects that device is connected, it will schedule management software for you and create link to the device inside. Once device is disconnected, everything is cleaned up. Learn how to configure it by [reading the docs][MIC]
## Manager scheduled externaly

In this approach, you add annotation to your Pod and KubeSerial mutating webhook will update pod spec when it's created to inject connection to device. Learn how to configure it by [reading the docs][MEC]

<!-- Links  -->
[MD]:  /components/manager.md            "Manager"
[MIC]: /configuration/managers/internal.md "Manager internal"
[MEC]: /configuration/managers/external.md "Manager external"
