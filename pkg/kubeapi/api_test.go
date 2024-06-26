package kubeapi

import (
	"context"
	"testing"

	kubeserialv1alpha1 "github.com/janekbaraniewski/kubeserial/pkg/apis/v1alpha1"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	"sigs.k8s.io/controller-runtime/pkg/client"
	runtimefake "sigs.k8s.io/controller-runtime/pkg/client/fake"
)

func GetFakeAPIAndScheme() (*runtime.Scheme, client.Client) {
	scheme := runtime.NewScheme()
	utilruntime.Must(clientgoscheme.AddToScheme(scheme))
	utilruntime.Must(kubeserialv1alpha1.AddToScheme(scheme))
	fakeClient := runtimefake.NewClientBuilder().WithScheme(scheme).Build()
	return scheme, fakeClient
}

func GetAPI(fakeClient client.Client, scheme *runtime.Scheme) API {
	return NewAPIClient(fakeClient, scheme)
}

func TestEnsureConfigMap(t *testing.T) {
	t.Parallel()
	scheme, fakeClient := GetFakeAPIAndScheme()
	api := GetAPI(fakeClient, scheme)

	cm := &corev1.ConfigMap{
		TypeMeta: metav1.TypeMeta{
			Kind:       "ConfigMap",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      "test-cm",
			Namespace: "test-namespace",
		},
		Data: map[string]string{
			"works": "true",
		},
	}

	err := api.EnsureObject(context.TODO(), &kubeserialv1alpha1.KubeSerial{}, cm)

	require.NoError(t, err)
	found := &corev1.ConfigMap{}
	//nolint:errcheck
	fakeClient.Get(
		context.TODO(),
		types.NamespacedName{Name: "test-cm", Namespace: "test-namespace"},
		found,
	)

	assert.Equal(t, "true", found.Data["works"])
}

func TestEnsureConfigMapUpdatesExisting(t *testing.T) {
	t.Parallel()
	scheme, fakeClient := GetFakeAPIAndScheme()
	api := GetAPI(fakeClient, scheme)
	//nolint:errcheck
	fakeClient.Create(context.TODO(), &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "test-cm",
			Namespace: "test-namespace",
		},
		Data: map[string]string{
			"data": "not-overwritten",
		},
	})

	err := api.EnsureObject(context.TODO(), &kubeserialv1alpha1.KubeSerial{}, &corev1.ConfigMap{
		TypeMeta: metav1.TypeMeta{
			Kind:       "ConfigMap",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      "test-cm",
			Namespace: "test-namespace",
		},
		Data: map[string]string{
			"data": "overwritten",
		},
	})

	require.NoError(t, err)
	found := &corev1.ConfigMap{}
	//nolint:errcheck
	fakeClient.Get(
		context.TODO(),
		types.NamespacedName{Name: "test-cm", Namespace: "test-namespace"},
		found,
	)

	assert.Equal(t, "overwritten", found.Data["data"])
}

func TestDeleteObject(t *testing.T) {
	t.Parallel()
	scheme, fakeClient := GetFakeAPIAndScheme()
	api := GetAPI(fakeClient, scheme)

	obj := &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "test-cm",
			Namespace: "test-namespace",
		},
		Data: map[string]string{
			"data": "not-overwritten",
		},
	}

	//nolint:errcheck
	fakeClient.Create(context.TODO(), obj)

	err := api.DeleteObject(context.TODO(), obj)
	require.NoError(t, err)

	lookup := &corev1.ConfigMap{}

	err = fakeClient.Get(
		context.TODO(),
		types.NamespacedName{Name: "test-cm", Namespace: "test-namespace"},
		lookup,
	)

	assert.True(t, errors.IsNotFound(err))
}
