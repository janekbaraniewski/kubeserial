package gateway

import (
	"testing"

	"github.com/janekbaraniewski/kubeserial/pkg/apis/kubeserial/v1alpha1"
	"github.com/stretchr/testify/assert"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestCreateDeployment(t *testing.T) {
	// TODO: improve this test
	deployment := CreateDeployment(&v1alpha1.Device{
		ObjectMeta: v1.ObjectMeta{
			Name: "test-device",
		},
		Status: v1alpha1.DeviceStatus{
			NodeName: "test-node",
		},
	})

	assert.Equal(t, map[string]string{
		"kubernetes.io/hostname": "test-node",
	}, deployment.Spec.Template.Spec.NodeSelector)
}
