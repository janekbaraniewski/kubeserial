package managers

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/janekbaraniewski/kubeserial/pkg/apis/v1alpha1"
	"github.com/janekbaraniewski/kubeserial/pkg/kubeapi"
	"github.com/janekbaraniewski/kubeserial/pkg/utils"
	"github.com/stretchr/testify/assert"
)

func TestSchedule(t *testing.T) {
	fs := utils.NewInMemoryFS()

	AddSpecFilesToFilesystem(t, fs)

	req := &v1alpha1.ManagerScheduleRequest{}
	manager := &v1alpha1.Manager{}
	api := kubeapi.NewFakeApiClient()

	err := Schedule(context.TODO(), fs, req, manager, "kubeserial", api)

	assert.Equal(t, nil, err)
}

func AddSpecFilesToFilesystem(t *testing.T, fs *utils.InMemoryFS) {
	AddFileToFileSystem(t, fs, "manager-configmap.yaml")
	AddFileToFileSystem(t, fs, "manager-deployment.yaml")
	AddFileToFileSystem(t, fs, "manager-service.yaml")
}

func AddFileToFileSystem(t *testing.T, fs *utils.InMemoryFS, path string) {
	file, err := fs.Create(fmt.Sprintf("/config/%v", path))

	assert.Equal(t, nil, err)

	absPath, _ := filepath.Abs(fmt.Sprintf("../../test-assets/%v", path))
	content, err := os.ReadFile(absPath)
	if err != nil {
		t.Fatalf("Failed to read yaml resource: %v", err)
	}

	file.Write(content)
	file.Close()
}
