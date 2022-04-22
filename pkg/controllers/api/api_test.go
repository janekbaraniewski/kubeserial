package api

import (
	"context"
	"testing"

	kubeserialv1alpha1 "github.com/janekbaraniewski/kubeserial/pkg/apis/kubeserial/v1alpha1"
	"github.com/stretchr/testify/assert"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/api/core/v1"
	networkingv1 "k8s.io/api/networking/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	"sigs.k8s.io/controller-runtime/pkg/client"
	runtimefake "sigs.k8s.io/controller-runtime/pkg/client/fake"
)

func GetFakeApiAndScheme() (*runtime.Scheme, client.Client) {
	scheme := runtime.NewScheme()
	utilruntime.Must(clientgoscheme.AddToScheme(scheme))
	utilruntime.Must(kubeserialv1alpha1.AddToScheme(scheme))
	fakeClient := runtimefake.NewClientBuilder().WithScheme(scheme).Build()
	return scheme, fakeClient
}

func GetApi(fakeClient client.Client, scheme *runtime.Scheme) API {
	return &ApiClient{
		Client: fakeClient,
		Scheme: scheme,
	}
}

func TestEnsureConfigMap(t *testing.T) {
	scheme, fakeClient := GetFakeApiAndScheme()
	api := GetApi(fakeClient, scheme)

	err := api.EnsureConfigMap(context.TODO(), &kubeserialv1alpha1.KubeSerial{}, &v1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "test-cm",
			Namespace: "test-namespace",
		},
		Data: map[string]string{
			"works": "true",
		},
	})

	assert.Equal(t, nil, err)
	found := &corev1.ConfigMap{}
	fakeClient.Get(
		context.TODO(),
		types.NamespacedName{Name: "test-cm", Namespace: "test-namespace"},
		found,
	)

	assert.Equal(t, "true", found.Data["works"])
}

func TestEnsureConfigMapDoesntOverwriteExisting(t *testing.T) {
	scheme, fakeClient := GetFakeApiAndScheme()
	api := GetApi(fakeClient, scheme)

	fakeClient.Create(context.TODO(), &v1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "test-cm",
			Namespace: "test-namespace",
		},
		Data: map[string]string{
			"data": "not-overwritten",
		},
	})

	err := api.EnsureConfigMap(context.TODO(), &kubeserialv1alpha1.KubeSerial{}, &v1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "test-cm",
			Namespace: "test-namespace",
		},
		Data: map[string]string{
			"data": "overwritten",
		},
	})

	assert.Equal(t, nil, err)
	found := &corev1.ConfigMap{}
	fakeClient.Get(
		context.TODO(),
		types.NamespacedName{Name: "test-cm", Namespace: "test-namespace"},
		found,
	)

	assert.Equal(t, "not-overwritten", found.Data["data"])

}

func TestEnsureService(t *testing.T) {
	scheme, fakeClient := GetFakeApiAndScheme()
	api := GetApi(fakeClient, scheme)

	err := api.EnsureService(context.TODO(), &kubeserialv1alpha1.KubeSerial{}, &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "test-service",
			Namespace: "test-namespace",
		},
	})

	assert.Equal(t, nil, err)

	found := &corev1.Service{}
	fakeClient.Get(context.TODO(), types.NamespacedName{Name: "test-service", Namespace: "test-namespace"}, found)
	assert.Equal(t, "test-service", found.ObjectMeta.Name)
}

func TestEnsureServiceDoesntOverwriteExisting(t *testing.T) {
	scheme, fakeClient := GetFakeApiAndScheme()
	api := GetApi(fakeClient, scheme)

	fakeClient.Create(context.TODO(), &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "test-service",
			Namespace: "test-namespace",
		},
		Spec: corev1.ServiceSpec{
			Ports: []corev1.ServicePort{
				{
					Port: 8080,
					Name: "original-port",
				},
			},
		},
	})

	err := api.EnsureService(context.TODO(), &kubeserialv1alpha1.KubeSerial{}, &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "test-service",
			Namespace: "test-namespace",
		},
		Spec: corev1.ServiceSpec{
			Ports: []corev1.ServicePort{
				{
					Port: 8080,
					Name: "overwritten-port",
				},
			},
		},
	})

	assert.Equal(t, nil, err)

	found := &corev1.Service{}
	fakeClient.Get(context.TODO(), types.NamespacedName{Name: "test-service", Namespace: "test-namespace"}, found)
	assert.Equal(t, "original-port", found.Spec.Ports[0].Name)
}

func TestEnsureIngress(t *testing.T) {
	scheme, fakeClient := GetFakeApiAndScheme()
	api := GetApi(fakeClient, scheme)

	err := api.EnsureIngress(context.TODO(), &kubeserialv1alpha1.KubeSerial{}, &networkingv1.Ingress{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "test-ingress",
			Namespace: "test-namespace",
		},
		Spec: networkingv1.IngressSpec{
			Rules: []networkingv1.IngressRule{
				{
					Host: "test-host.com",
				},
			},
		},
	})

	assert.Equal(t, nil, err)

	found := &networkingv1.Ingress{}
	fakeClient.Get(context.TODO(), types.NamespacedName{Name: "test-ingress", Namespace: "test-namespace"}, found)
	assert.Equal(t, "test-host.com", found.Spec.Rules[0].Host)
}

func TestEnsureIngressDoesntOverwriteExisting(t *testing.T) {
	scheme, fakeClient := GetFakeApiAndScheme()
	api := GetApi(fakeClient, scheme)

	fakeClient.Create(context.TODO(), &networkingv1.Ingress{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "test-ingress",
			Namespace: "test-namespace",
		},
		Spec: networkingv1.IngressSpec{
			Rules: []networkingv1.IngressRule{
				{
					Host: "original-host.com",
				},
			},
		},
	})

	err := api.EnsureIngress(context.TODO(), &kubeserialv1alpha1.KubeSerial{}, &networkingv1.Ingress{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "test-ingress",
			Namespace: "test-namespace",
		},
		Spec: networkingv1.IngressSpec{
			Rules: []networkingv1.IngressRule{
				{
					Host: "new-host.com",
				},
			},
		},
	})

	assert.Equal(t, nil, err)

	found := &networkingv1.Ingress{}
	fakeClient.Get(context.TODO(), types.NamespacedName{Name: "test-ingress", Namespace: "test-namespace"}, found)
	assert.Equal(t, "original-host.com", found.Spec.Rules[0].Host)
}

func TestEnsureDeployment(t *testing.T) {
	scheme, fakeClient := GetFakeApiAndScheme()
	api := GetApi(fakeClient, scheme)

	err := api.EnsureDeployment(context.TODO(), &kubeserialv1alpha1.KubeSerial{}, &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "test-deployment",
			Namespace: "test-namespace",
		},
		Spec: appsv1.DeploymentSpec{
			Template: v1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Name: "test-deployment",
				},
			},
		},
	})

	assert.Equal(t, nil, err)

	found := &appsv1.Deployment{}
	fakeClient.Get(context.TODO(), types.NamespacedName{Name: "test-deployment", Namespace: "test-namespace"}, found)
	assert.Equal(t, "test-deployment", found.Spec.Template.ObjectMeta.Name)
}

func TestEnsureDeploymentDoesntOverwriteExisting(t *testing.T) {
	scheme, fakeClient := GetFakeApiAndScheme()
	api := GetApi(fakeClient, scheme)

	fakeClient.Create(context.TODO(), &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "test-deployment",
			Namespace: "test-namespace",
		},
		Spec: appsv1.DeploymentSpec{
			Template: v1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Name: "original-deployment",
				},
			},
		},
	})

	err := api.EnsureDeployment(context.TODO(), &kubeserialv1alpha1.KubeSerial{}, &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "test-deployment",
			Namespace: "test-namespace",
		},
		Spec: appsv1.DeploymentSpec{
			Template: v1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Name: "new-deployment",
				},
			},
		},
	})

	assert.Equal(t, nil, err)

	found := &appsv1.Deployment{}
	fakeClient.Get(context.TODO(), types.NamespacedName{Name: "test-deployment", Namespace: "test-namespace"}, found)
	assert.Equal(t, "original-deployment", found.Spec.Template.ObjectMeta.Name)
}
