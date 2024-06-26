package monitor

import (
	"strings"
	"testing"

	kubeserial "github.com/janekbaraniewski/kubeserial/pkg"
	appv1alpha1 "github.com/janekbaraniewski/kubeserial/pkg/apis/v1alpha1"
	"github.com/janekbaraniewski/kubeserial/pkg/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCreateDaemonSet(t *testing.T) {
	t.Parallel()
	fs := utils.NewInMemoryFS()
	if err := fs.AddFileFromHostPath(string(kubeserial.MonitorDSSpecPath)); err != nil {
		t.Fatalf("Failed to load test asset: %v", err)
	}

	result, err := CreateDaemonSet(fs)
	require.NoError(t, err)
	assert.Equal(t, "kubeserial-monitor", result.ObjectMeta.Name)
	imageAndTag := strings.Split(result.Spec.Template.Spec.Containers[0].Image, ":")
	assert.Equal(t, "ghcr.io/janekbaraniewski/kubeserial-device-monitor", imageAndTag[0])
}

func TestCreateConfigMap(t *testing.T) {
	t.Parallel()
	fs := utils.NewInMemoryFS()
	if err := fs.AddFileFromHostPath(string(kubeserial.MonitorCMSpecPath)); err != nil {
		t.Fatalf("Failed to load test asset: %v", err)
	}

	devices := []appv1alpha1.SerialDevice2{
		{
			Name:      "testdevice",
			IDVendor:  "123",
			IDProduct: "456",
		},
	}

	result, err := CreateConfigMap(fs, devices)

	require.NoError(t, err)

	expectedUdevConfig := "SUBSYSTEM==\"tty\", ATTRS{idVendor}==\"123\", ATTRS{idProduct}==\"456\", SYMLINK+=\"testdevice\"\n"

	assert.Equal(t, expectedUdevConfig, result.Data["98-devices.rules"])
}
