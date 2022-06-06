package managers

import (
	"context"
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
	if err := fs.AddFileFromHostPath("manager-configmap.yaml"); err != nil {
		t.Fatalf("Failed to load test asset: %v", err)
	}
	if err := fs.AddFileFromHostPath("manager-deployment.yaml"); err != nil {
		t.Fatalf("Failed to load test asset: %v", err)
	}
	if err := fs.AddFileFromHostPath("manager-service.yaml"); err != nil {
		t.Fatalf("Failed to load test asset: %v", err)
	}
}
