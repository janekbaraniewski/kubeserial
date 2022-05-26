package monitor

import (
	"testing"

	appv1alpha1 "github.com/janekbaraniewski/kubeserial/pkg/apis/v1alpha1"
	"github.com/stretchr/testify/assert"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestCreateDaemonSet(t *testing.T) {
	// TODO: improve this test
	cr := &appv1alpha1.KubeSerial{
		ObjectMeta: v1.ObjectMeta{
			Name: "test-config",
		},
	}

	result := CreateDaemonSet(cr, "0.0.1")

	assert.Equal(t, "test-config-monitor", result.ObjectMeta.Name)
	assert.Equal(t, "janekbaraniewski/kubeserial-device-monitor:0.0.1", result.Spec.Template.Spec.Containers[0].Image)
}
