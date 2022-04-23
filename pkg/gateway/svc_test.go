package gateway

import (
	"testing"

	"github.com/janekbaraniewski/kubeserial/pkg/apis/kubeserial/v1alpha1"
	"github.com/stretchr/testify/assert"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestCreateService(t *testing.T) {
	// TODO: improve this test
	svc := CreateService(
		&v1alpha1.KubeSerial{
			ObjectMeta: v1.ObjectMeta{
				Name: "test-cr",
			},
		},
		&v1alpha1.Device_2{
			Name: "test-device",
		},
	)

	assert.Equal(t, "test-cr-test-device-gateway", svc.ObjectMeta.Name)
}
