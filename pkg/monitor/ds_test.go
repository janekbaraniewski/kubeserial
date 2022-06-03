package monitor

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/janekbaraniewski/kubeserial/pkg/utils"
	"github.com/stretchr/testify/assert"
)

func TestCreateDaemonSet(t *testing.T) {
	fs := utils.NewInMemoryFS()
	file, err := fs.Create("/config/monitor-spec.yaml")

	assert.Equal(t, nil, err)

	absPath, _ := filepath.Abs("../assets/monitor-spec.yaml")
	content, err := os.ReadFile(absPath)
	if err != nil {
		t.Fatalf("Failed to read yaml resource: %v", err)
	}

	file.Write(content)

	result, err := CreateDaemonSet(fs)

	assert.Equal(t, nil, err)
	assert.Equal(t, "test-monitor", result.ObjectMeta.Name)
	assert.Equal(t, "sample-image:dev", result.Spec.Template.Spec.Containers[0].Image)
}
