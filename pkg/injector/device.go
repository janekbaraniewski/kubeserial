package injector

import (
	"fmt"
	"strings"

	corev1 "k8s.io/api/core/v1"
)

func AddDeviceInjector(spec *corev1.PodSpec, device string) {
	spec.Containers = append(spec.Containers, sidecarContainerSpec(device))
	spec.Volumes = append(spec.Volumes, podVolumeSpec())
	containers := []corev1.Container{}
	for _, container := range spec.Containers {
		container.VolumeMounts = append(container.VolumeMounts, volumeMountSpec())
		containers = append(containers, container)
	}
	spec.Containers = containers
}

func sidecarContainerSpec(device string) corev1.Container {
	return corev1.Container{
		Name:    "device-mounter",
		Image:   "alpine/socat:1.7.4.3-r0",
		Command: []string{"/bin/sh"},
		Args: []string{
			"-c",
			fmt.Sprintf("sleep 5 && socat -d -d pty,raw,echo=0,b115200,link=/dev/devices/%v,perm=0660,group=tty tcp:%v:3333", device, strings.ToLower(device+"-gateway")),
		},
	}
}

func podVolumeSpec() corev1.Volume {
	return corev1.Volume{
		Name: "devices",
		VolumeSource: corev1.VolumeSource{
			EmptyDir: &corev1.EmptyDirVolumeSource{},
		},
	}
}

func volumeMountSpec() corev1.VolumeMount {
	return corev1.VolumeMount{
		Name:      "devices",
		ReadOnly:  false,
		MountPath: "/dev/devices",
	}
}
