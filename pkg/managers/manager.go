package managers

import (
	"context"

	appv1alpha1 "github.com/janekbaraniewski/kubeserial/pkg/apis/kubeserial/v1alpha1"
	"github.com/janekbaraniewski/kubeserial/pkg/controllers/api"
	"k8s.io/apimachinery/pkg/types"
)

type Manager struct {
	Image      string
	RunCmnd    string
	Config     string
	ConfigPath string
}

func Schedule(ctx context.Context, request *appv1alpha1.ManagerScheduleRequest, mgr *appv1alpha1.Manager, api api.API) error {
	manager := &Manager{
		Image:      mgr.Spec.Image.Repository + ":" + mgr.Spec.Image.Tag,
		RunCmnd:    mgr.Spec.RunCmd,
		Config:     mgr.Spec.Config,
		ConfigPath: mgr.Spec.ConfigPath,
	}
	cr := types.NamespacedName{
		Name:      request.Name,
		Namespace: request.Namespace,
	}
	device := types.NamespacedName{
		Name:      request.Spec.Device,
		Namespace: mgr.Namespace,
	}
	if mgr.Spec.Config != "" {
		cm := manager.CreateConfigMap(cr, device)
		if err := api.EnsureConfigMap(ctx, request, cm); err != nil {
			return err
		}
	}
	deploy := manager.CreateDeployment(cr, device, mgr.Spec.Config != "")
	svc := manager.CreateService(cr, device)

	if err := api.EnsureDeployment(ctx, request, deploy); err != nil {
		return err
	}

	if err := api.EnsureService(ctx, request, svc); err != nil {
		return err
	}

	// if cr.Spec.Ingress.Enabled {
	// 	ingress := manager.CreateIngress(cr, device, cr.Spec.Ingress.Domain)
	// 	if err := api.EnsureIngress(ctx, cr, ingress); err != nil {
	// 		return err
	// 	}
	// }
	return nil
}

func (m *Manager) Delete(ctx context.Context, cr *appv1alpha1.KubeSerial, device *appv1alpha1.SerialDevice_2, api api.API) error {
	name := m.GetName(cr.Name, device.Name)

	if err := api.DeleteDeployment(ctx, cr, name); err != nil {
		return err
	}
	if err := api.DeleteConfigMap(ctx, cr, name); err != nil {
		return err
	}
	if err := api.DeleteService(ctx, cr, name); err != nil {
		return err
	}
	if err := api.DeleteIngress(ctx, cr, name); err != nil {
		return err
	}

	return nil
}
