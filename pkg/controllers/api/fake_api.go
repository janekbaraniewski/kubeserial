package api

import (
	"context"

	appv1alpha1 "github.com/janekbaraniewski/kubeserial/pkg/apis/kubeserial/v1alpha1"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	networkingv1 "k8s.io/api/networking/v1"
)

type FakeApiClient struct {
	Operations []string
}

func NewFakeApiClient() FakeApiClient {
	return FakeApiClient{
		Operations: []string{},
	}
}

func (a *FakeApiClient) EnsureConfigMap(ctx context.Context, cr *appv1alpha1.KubeSerial, cm *corev1.ConfigMap) error {
	a.Operations = append(a.Operations, "EnsureConfigMap")
	return nil
}
func (a *FakeApiClient) EnsureService(ctx context.Context, cr *appv1alpha1.KubeSerial, svc *corev1.Service) error {
	a.Operations = append(a.Operations, "EnsureService")
	return nil
}
func (a *FakeApiClient) EnsureIngress(ctx context.Context, cr *appv1alpha1.KubeSerial, ingress *networkingv1.Ingress) error {
	a.Operations = append(a.Operations, "EnsureIngress")
	return nil
}
func (a *FakeApiClient) EnsureDeployment(ctx context.Context, cr *appv1alpha1.KubeSerial, deployment *appsv1.Deployment) error {
	a.Operations = append(a.Operations, "EnsureDeployment")
	return nil
}
func (a *FakeApiClient) EnsureDaemonSet(ctx context.Context, cr *appv1alpha1.KubeSerial, ds *appsv1.DaemonSet) error {
	a.Operations = append(a.Operations, "EnsureDaemonSet")
	return nil
}
func (a *FakeApiClient) DeleteDeployment(ctx context.Context, cr *appv1alpha1.KubeSerial, name string) error {
	a.Operations = append(a.Operations, "DeleteDeployment")
	return nil
}
func (a *FakeApiClient) DeleteConfigMap(ctx context.Context, cr *appv1alpha1.KubeSerial, name string) error {
	a.Operations = append(a.Operations, "DeleteConfigMap")
	return nil
}
func (a *FakeApiClient) DeleteService(ctx context.Context, cr *appv1alpha1.KubeSerial, name string) error {
	a.Operations = append(a.Operations, "DeleteService")
	return nil
}
func (a *FakeApiClient) DeleteIngress(ctx context.Context, cr *appv1alpha1.KubeSerial, name string) error {
	a.Operations = append(a.Operations, "DeleteIngress")
	return nil
}
