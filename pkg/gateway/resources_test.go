package gateway

import (
	"testing"

	kubeserial "github.com/janekbaraniewski/kubeserial/pkg"
	"github.com/janekbaraniewski/kubeserial/pkg/apis/v1alpha1"
	"github.com/janekbaraniewski/kubeserial/pkg/utils"
	"github.com/stretchr/testify/assert"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestBuilder_Build(t *testing.T) {
	t.Parallel()
	fs := utils.NewInMemoryFS()

	AddSpecFilesToFilesystem(t, fs)

	device := &v1alpha1.SerialDevice{
		ObjectMeta: v1.ObjectMeta{
			Name: "test-device",
		},
	}

	objects := NewBuilder(device, fs).Build()

	assert.NotEqual(t, nil, objects)
	assert.Equal(t, 3, len(objects))
}

func TestCreateConfigMap(t *testing.T) {
	t.Parallel()
	fs := utils.NewInMemoryFS()

	AddSpecFilesToFilesystem(t, fs)

	cm, err := CreateConfigMap(
		&v1alpha1.SerialDevice{
			ObjectMeta: v1.ObjectMeta{
				Name: "test-device",
			},
		},
		fs,
	)

	assert.Equal(t, nil, err)
	assert.Equal(t, "test-device-gateway", cm.ObjectMeta.Name)
}

func TestCreateDeployment(t *testing.T) {
	t.Parallel()
	fs := utils.NewInMemoryFS()

	AddSpecFilesToFilesystem(t, fs)

	deployment, err := CreateDeployment(
		&v1alpha1.SerialDevice{
			ObjectMeta: v1.ObjectMeta{
				Name: "test-device",
			},
			Status: v1alpha1.SerialDeviceStatus{
				NodeName: "test-node",
			},
		},
		fs,
	)

	assert.Equal(t, nil, err)
	assert.Equal(t, map[string]string{
		"kubernetes.io/hostname": "test-node",
	}, deployment.Spec.Template.Spec.NodeSelector)
	for _, volume := range deployment.Spec.Template.Spec.Volumes {
		if volume.Name == "config" {
			assert.Equal(t, "test-device-gateway", volume.ConfigMap.Name)
		}
	}
}

func TestCreateService(t *testing.T) {
	t.Parallel()
	fs := utils.NewInMemoryFS()

	AddSpecFilesToFilesystem(t, fs)

	// TODO: improve this test
	svc, err := CreateService(
		&v1alpha1.SerialDevice{
			ObjectMeta: v1.ObjectMeta{
				Name: "test-device",
			},
		},
		fs,
	)

	assert.Equal(t, nil, err)
	assert.Equal(t, "test-device-gateway", svc.ObjectMeta.Name)
}

func AddSpecFilesToFilesystem(t *testing.T, fs *utils.InMemoryFS) {
	t.Helper()
	if err := fs.AddFileFromHostPath(string(kubeserial.GatewayCMSpecPath)); err != nil {
		t.Fatalf("Failed to load test asset: %v", err)
	}
	if err := fs.AddFileFromHostPath(string(kubeserial.GatewayDeploySpecPath)); err != nil {
		t.Fatalf("Failed to load test asset: %v", err)
	}
	if err := fs.AddFileFromHostPath(string(kubeserial.GatewaySvcSpecPath)); err != nil {
		t.Fatalf("Failed to load test asset: %v", err)
	}
}
