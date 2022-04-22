package api

import (
	"context"
	"testing"

	kubeserialv1alpha1 "github.com/janekbaraniewski/kubeserial/pkg/apis/kubeserial/v1alpha1"
	"github.com/stretchr/testify/assert"
	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/api/core/v1"
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
