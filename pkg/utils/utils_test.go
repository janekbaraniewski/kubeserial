package utils

import (
	"context"
	"testing"

	"github.com/janekbaraniewski/kubeserial/pkg/apis/kubeserial/v1alpha1"
	"github.com/janekbaraniewski/kubeserial/pkg/generated/clientset/versioned/fake"
	"github.com/stretchr/testify/assert"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestSetDeviceCondition(t *testing.T) {
	device := &v1alpha1.Device{
		ObjectMeta: v1.ObjectMeta{
			Name:      "test-device",
			Namespace: "test-ns",
		},
	}

	client := fake.NewSimpleClientset()

	client.AppV1alpha1().Devices("test-ns").Create(context.TODO(), device, v1.CreateOptions{})

	SetDeviceCondition(&device.Status.Conditions, v1alpha1.DeviceCondition{
		Status: v1.ConditionTrue,
		Type:   v1alpha1.DeviceReady,
		Reason: "TestPassed",
	})

	assert.Equal(t, "TestPassed", device.Status.Conditions[0].Reason)

	client.AppV1alpha1().Devices("test-ns").Update(context.TODO(), device, v1.UpdateOptions{})

	foundDevice, err := client.AppV1alpha1().Devices("test-ns").Get(context.TODO(), "test-device", v1.GetOptions{})

	assert.Equal(t, nil, err)

	assert.Equal(t, "TestPassed", foundDevice.Status.Conditions[0].Reason)
}
