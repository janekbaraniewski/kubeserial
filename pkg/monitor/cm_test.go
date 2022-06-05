package monitor

import (
	"os"
	"path/filepath"
	"testing"

	appv1alpha1 "github.com/janekbaraniewski/kubeserial/pkg/apis/v1alpha1"
	"github.com/janekbaraniewski/kubeserial/pkg/utils"
	"github.com/stretchr/testify/assert"
)

func TestCreateConfigMap(t *testing.T) {
	fs := utils.NewInMemoryFS()
	file, err := fs.Create("/config/monitor-configmap.yaml")

	assert.Equal(t, nil, err)

	absPath, _ := filepath.Abs("../assets/monitor-configmap.yaml")
	content, err := os.ReadFile(absPath)
	if err != nil {
		t.Fatalf("Failed to read yaml resource: %v", err)
	}

	file.Write(content)

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
