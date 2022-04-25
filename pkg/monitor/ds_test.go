package monitor

import (
	"testing"

	appv1alpha1 "github.com/janekbaraniewski/kubeserial/pkg/apis/kubeserial/v1alpha1"
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

	result := CreateDaemonSet(cr)

	assert.Equal(t, "test-config-monitor", result.ObjectMeta.Name)
}
