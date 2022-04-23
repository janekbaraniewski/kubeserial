package gateway

import (
	"testing"

	"github.com/janekbaraniewski/kubeserial/pkg/apis/kubeserial/v1alpha1"
	"github.com/stretchr/testify/assert"
)

func TestCreateDeployment(t *testing.T) {
	// TODO: improve this test
	deployment := CreateDeployment(&v1alpha1.KubeSerial{}, &v1alpha1.Device_2{}, "test-node")

	assert.Equal(t, map[string]string{
		"kubernetes.io/hostname": "test-node",
	}, deployment.Spec.Template.Spec.NodeSelector)
}
