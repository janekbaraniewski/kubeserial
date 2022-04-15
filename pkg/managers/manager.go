package managers

import (
	"context"

	appv1alpha1 "github.com/janekbaraniewski/kubeserial/pkg/apis/kubeserial/v1alpha1"
	"github.com/janekbaraniewski/kubeserial/pkg/controllers/api"
)

type Manager struct {
	Image      string
	RunCmnd    string
	Config     string
	ConfigPath string
}

func (m *Manager) Schedule(ctx context.Context, cr *appv1alpha1.KubeSerial, device *appv1alpha1.Device, api *api.ApiClient) error {
	cm := m.CreateConfigMap(cr, device)
	deploy := m.CreateDeployment(cr, device)
	svc := m.CreateService(cr, device)

	if err := api.EnsureConfigMap(ctx, cr, cm); err != nil {
		return err
	}

	if err := api.EnsureDeployment(ctx, cr, deploy); err != nil {
		return err
	}

	if err := api.EnsureService(ctx, cr, svc); err != nil {
		return err
	}

	if cr.Spec.Ingress.Enabled {
		ingress := m.CreateIngress(cr, device, cr.Spec.Ingress.Domain)
		if err := api.EnsureIngress(ctx, cr, ingress); err != nil {
			return err
		}
	}

	return nil
}

func (m *Manager) Delete(ctx context.Context, cr *appv1alpha1.KubeSerial, device *appv1alpha1.Device, api *api.ApiClient) error {
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
