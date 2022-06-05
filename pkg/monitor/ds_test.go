package monitor

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/janekbaraniewski/kubeserial/pkg/utils"
	"github.com/stretchr/testify/assert"
)

func TestCreateDaemonSet(t *testing.T) {
	fs := utils.NewInMemoryFS()
	file, err := fs.Create("/config/monitor-daemonset.yaml")

	assert.Equal(t, nil, err)

	absPath, _ := filepath.Abs("../assets/monitor-daemonset.yaml")
	content, err := os.ReadFile(absPath)
	if err != nil {
		t.Fatalf("Failed to read yaml resource: %v", err)
	}

	file.Write(content)

	result, err := CreateDaemonSet(fs)

	assert.Equal(t, nil, err)
	assert.Equal(t, "kubeserial-monitor", result.ObjectMeta.Name)
	imageAndTag := strings.Split(result.Spec.Template.Spec.Containers[0].Image, ":")
	assert.Equal(t, "janekbaraniewski/kubeserial-device-monitor", imageAndTag[0])
}
