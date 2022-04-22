package gateway

import (
	"testing"

	"github.com/janekbaraniewski/kubeserial/pkg/apis/kubeserial/v1alpha1"
	"github.com/stretchr/testify/assert"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestCreateConfigMap(t *testing.T) {
	// TODO: improve this test
	cm := CreateConfigMap(
		&v1alpha1.KubeSerial{
			ObjectMeta: v1.ObjectMeta{
				Name: "test-cr",
			},
		},
		&v1alpha1.Device{
			Name: "test-device",
		},
	)

	assert.Equal(t, "test-cr-test-device-gateway", cm.ObjectMeta.Name)
}
