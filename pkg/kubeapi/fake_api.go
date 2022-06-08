package kubeapi

import (
	"context"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type Operation struct {
	Action string
	GVK    schema.GroupVersionKind
}

type FakeApiClient struct {
	Operations []Operation
}

func NewFakeApiClient() *FakeApiClient {
	return &FakeApiClient{
		Operations: []Operation{},
	}
}

func (a *FakeApiClient) EnsureObject(ctx context.Context, cr metav1.Object, obj client.Object) error {
	a.Operations = append(a.Operations, Operation{"EnsureObject", obj.GetObjectKind().GroupVersionKind()})
	return nil
}
func (a *FakeApiClient) DeleteObject(ctx context.Context, obj client.Object) error {
	a.Operations = append(a.Operations, Operation{"DeleteObject", obj.GetObjectKind().GroupVersionKind()})
	return nil
}
