package injector

import (
	"testing"

	"github.com/stretchr/testify/assert"
	corev1 "k8s.io/api/core/v1"
)

func TestAddDeviceInjector(t *testing.T) {
	testPod := &corev1.Pod{}

	AddDeviceInjector(&testPod.Spec, "test-device")

	assert.Equal(t, 1, len(testPod.Spec.Containers))
	assert.Equal(t, corev1.Container{
		Name:  "device-mounter",
		Image: "alpine/socat:1.7.4.3-r0",
		Command: []string{
			"/bin/sh",
		},
		Args: []string{
			"-c",
			"sleep 5 && socat -d -d pty,raw,echo=0,b115200,link=/dev/devices/test-device,perm=0660,group=tty tcp:test-device-gateway:3333",
		},
		VolumeMounts: []corev1.VolumeMount{
			{
				Name:      "devices",
				ReadOnly:  false,
				MountPath: "/dev/devices",
			},
		},
	}, testPod.Spec.Containers[0])
	assert.Equal(t, 1, len(testPod.Spec.Volumes))
	assert.Equal(t, corev1.Volume{
		Name: "devices",
		VolumeSource: corev1.VolumeSource{
			EmptyDir: &corev1.EmptyDirVolumeSource{},
		},
	}, testPod.Spec.Volumes[0])
}
