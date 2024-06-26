package monitor

import (
	"context"
	"testing"

	"github.com/janekbaraniewski/kubeserial/pkg/apis/v1alpha1"
	"github.com/janekbaraniewski/kubeserial/pkg/generated/clientset/versioned/fake"
	"github.com/janekbaraniewski/kubeserial/pkg/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	testclient "k8s.io/client-go/kubernetes/fake"
)

func TestUpdateDeviceState_ConfigMap(t *testing.T) {
	t.Parallel()
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
	fs := utils.NewInMemoryFS()
	monitor := NewMonitor(
		fakeClientset,
		fakeClientsetKubeserial,
		"test-ns",
		"test-node",
		fs,
	)
	monitor.UpdateDeviceState(ctx)
}

func TestUpdateDeviceState_Device(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	getDevice := func(ready, available v1.ConditionStatus, node string) *v1alpha1.SerialDevice {
		return &v1alpha1.SerialDevice{
			ObjectMeta: v1.ObjectMeta{
				Name: "test-device",
			},
			Spec: v1alpha1.SerialDeviceSpec{
				Name:      "test-device",
				IDVendor:  "123",
				IDProduct: "456",
				Manager:   "test-manager",
			},
			Status: v1alpha1.SerialDeviceStatus{
				Conditions: []v1alpha1.SerialDeviceCondition{
					{
						Type:   v1alpha1.SerialDeviceReady,
						Status: ready,
					},
					{
						Type:   v1alpha1.SerialDeviceAvailable,
						Status: available,
						Reason: "testing",
					},
				},
				NodeName: node,
			},
		}
	}
	testCases := []struct {
		Name            string
		InitReady       v1.ConditionStatus
		InitAvailable   v1.ConditionStatus
		InitNode        string
		ResultAvailable v1.ConditionStatus
		ResultNode      string
		CreateDevice    bool
	}{
		{
			Name:            "test-update-available-condition-when-available",
			InitReady:       v1.ConditionTrue,
			InitAvailable:   v1.ConditionFalse,
			InitNode:        "",
			ResultAvailable: v1.ConditionTrue,
			ResultNode:      "test-node",
			CreateDevice:    true,
		},
		{
			Name:            "test-update-available-condition-when-unavailable",
			InitReady:       v1.ConditionTrue,
			InitAvailable:   v1.ConditionTrue,
			InitNode:        "test-node",
			ResultAvailable: v1.ConditionFalse,
			ResultNode:      "",
			CreateDevice:    false,
		},
		{
			Name:            "test-dont-update-not-ready-device",
			InitReady:       v1.ConditionFalse,
			InitAvailable:   v1.ConditionUnknown,
			InitNode:        "",
			ResultAvailable: v1.ConditionUnknown,
			ResultNode:      "",
		},
	}
	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			t.Parallel()
			device := getDevice(testCase.InitReady, testCase.InitAvailable, testCase.InitNode)
			fakeClientset := testclient.NewSimpleClientset()
			fakeClientsetKubeserial := fake.NewSimpleClientset(device)
			fs := utils.NewInMemoryFS()
			if testCase.CreateDevice {
				//nolint:errcheck
				fs.Create("/dev/" + device.Name)
			}
			monitor := NewMonitor(
				fakeClientset,
				fakeClientsetKubeserial,
				"test-ns",
				"test-node",
				fs,
			)

			monitor.UpdateDeviceState(ctx)

			foundDevice, err := fakeClientsetKubeserial.AppV1alpha1().SerialDevices().Get(
				ctx, device.Name, v1.GetOptions{})

			require.NoError(t, err)
			availableCondition := foundDevice.GetCondition(v1alpha1.SerialDeviceAvailable)
			assert.Equal(t, testCase.ResultAvailable, availableCondition.Status)
			assert.Equal(t, testCase.ResultNode, foundDevice.Status.NodeName)
		})
	}
}
