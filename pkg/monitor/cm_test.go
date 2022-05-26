package monitor

import (
	"testing"

	appv1alpha1 "github.com/janekbaraniewski/kubeserial/pkg/apis/v1alpha1"
	"github.com/stretchr/testify/assert"
)

func TestCreateConfigMap(t *testing.T) {
	cr := &appv1alpha1.KubeSerial{
		Spec: appv1alpha1.KubeSerialSpec{
			SerialDevices: []appv1alpha1.SerialDevice_2{
				{
					Name:      "testdevice",
					Subsystem: "tty",
					IdVendor:  "123",
					IdProduct: "456",
				},
			},
		},
	}

	result := CreateConfigMap(cr)

	expectedUdevConfig := "SUBSYSTEM==\"tty\", ATTRS{idVendor}==\"123\", ATTRS{idProduct}==\"456\", SYMLINK+=\"testdevice\"\n"

	assert.Equal(t, expectedUdevConfig, result.Data["98-devices.rules"])
}
