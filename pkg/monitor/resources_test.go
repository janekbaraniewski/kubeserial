package monitor

import (
	"strings"
	"testing"

	kubeserial "github.com/janekbaraniewski/kubeserial/pkg"
	appv1alpha1 "github.com/janekbaraniewski/kubeserial/pkg/apis/v1alpha1"
	"github.com/janekbaraniewski/kubeserial/pkg/utils"
	"github.com/stretchr/testify/assert"
)

func TestCreateDaemonSet(t *testing.T) {
	fs := utils.NewInMemoryFS()
	if err := fs.AddFileFromHostPath(string(kubeserial.MonitorDSSpecPath)); err != nil {
		t.Fatalf("Failed to load test asset: %v", err)
	}

	result, err := CreateDaemonSet(fs)
	assert.Equal(t, nil, err)
	assert.Equal(t, "kubeserial-monitor", result.ObjectMeta.Name)
	imageAndTag := strings.Split(result.Spec.Template.Spec.Containers[0].Image, ":")
	assert.Equal(t, "janekbaraniewski/kubeserial-device-monitor", imageAndTag[0])
}

func TestCreateConfigMap(t *testing.T) {
	fs := utils.NewInMemoryFS()
	if err := fs.AddFileFromHostPath(string(kubeserial.MonitorCMSpecPath)); err != nil {
		t.Fatalf("Failed to load test asset: %v", err)
	}

	devices := []appv1alpha1.SerialDevice_2{
		{
			Name:      "testdevice",
			IdVendor:  "123",
			IdProduct: "456",
		},
	}

	result, err := CreateConfigMap(fs, devices)

	assert.Equal(t, nil, err)

	expectedUdevConfig := "SUBSYSTEM==\"tty\", ATTRS{idVendor}==\"123\", ATTRS{idProduct}==\"456\", SYMLINK+=\"testdevice\"\n"

	assert.Equal(t, expectedUdevConfig, result.Data["98-devices.rules"])
}
