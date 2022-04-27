package monitor

import (
	"context"
	"os"
	"testing"

	"github.com/janekbaraniewski/kubeserial/pkg/apis/kubeserial/v1alpha1"
	"github.com/janekbaraniewski/kubeserial/pkg/generated/clientset/versioned/fake"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	testclient "k8s.io/client-go/kubernetes/fake"
)

func TestUpdateDeviceState_ConfigMap(t *testing.T) {
	ctx := context.Background()

	cm := &corev1.ConfigMap{
		ObjectMeta: v1.ObjectMeta{
			Name:      "test-device-state",
			Namespace: "test-ns",
			Labels: map[string]string{
				"type": "DeviceState",
			},
		},
		Data: map[string]string{
			"available": "false",
		},
	}

	fakeClientset := testclient.NewSimpleClientset(cm)
	fakeClientsetKubeserial := fake.NewSimpleClientset()
	fs := afero.NewMemMapFs()
	monitor := NewMonitor(fakeClientset, fakeClientsetKubeserial, "test-ns", fs.Stat)
	monitor.UpdateDeviceState(ctx)
}

func TestUpdateDeviceState_Device(t *testing.T) {
	ctx := context.Background()
	device := &v1alpha1.Device{
		ObjectMeta: v1.ObjectMeta{
			Name:      "test-device",
			Namespace: "test-ns",
		},
		Spec: v1alpha1.DeviceSpec{
			Name:      "test-device",
			IdVendor:  "123",
			IdProduct: "456",
			Subsystem: "tty",
			Manager:   "test-manager",
		},
		Status: v1alpha1.DeviceStatus{
			Conditions: []v1alpha1.DeviceCondition{
				{
					Type:   v1alpha1.DeviceAvailable,
					Status: v1.ConditionFalse,
				},
			},
		},
	}
	os.Setenv("NODE_NAME", "test-node")
	{
		t.Run("test-update-available-condition-when-available", func(t *testing.T) {
			fakeClientset := testclient.NewSimpleClientset()
			fakeClientsetKubeserial := fake.NewSimpleClientset(device)
			fs := afero.NewMemMapFs()
			fs.Create("/dev/tty" + device.Name)
			monitor := NewMonitor(fakeClientset, fakeClientsetKubeserial, "test-ns", fs.Stat)

			monitor.UpdateDeviceState(ctx)

			foundDevice, err := fakeClientsetKubeserial.AppV1alpha1().Devices("test-ns").Get(
				ctx, device.Name, v1.GetOptions{})

			assert.Equal(t, nil, err)
			assert.Equal(t, v1.ConditionTrue, foundDevice.Status.Conditions[0].Status)
			assert.Equal(t, "DeviceAvailable", foundDevice.Status.Conditions[0].Reason)
			assert.Equal(t, "test-node", foundDevice.Status.NodeName)
		})
	}
	{
		t.Run("test-update-available-condition-when-unavailable", func(t *testing.T) {
			device.Status.Conditions[0].Status = v1.ConditionTrue
			device.Status.NodeName = "test-node"
			fakeClientset := testclient.NewSimpleClientset()
			fakeClientsetKubeserial := fake.NewSimpleClientset(device)
			fs := afero.NewMemMapFs()
			monitor := NewMonitor(fakeClientset, fakeClientsetKubeserial, "test-ns", fs.Stat)

			monitor.UpdateDeviceState(ctx)

			foundDevice, err := fakeClientsetKubeserial.AppV1alpha1().Devices("test-ns").Get(
				ctx, device.Name, v1.GetOptions{})

			assert.Equal(t, nil, err)
			assert.Equal(t, v1.ConditionFalse, foundDevice.Status.Conditions[0].Status)
			assert.Equal(t, "DeviceUnavailable", foundDevice.Status.Conditions[0].Reason)
			assert.Equal(t, "", foundDevice.Status.NodeName)
		})
	}
	os.Unsetenv("NODE_NAME")
}
