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

type FakeAPIClient struct {
	Operations []Operation
}

func NewFakeAPIClient() *FakeAPIClient {
	return &FakeAPIClient{
		Operations: []Operation{},
	}
}

func (a *FakeAPIClient) EnsureObject(_ context.Context, _ metav1.Object, obj client.Object) error {
	a.Operations = append(a.Operations, Operation{"EnsureObject", obj.GetObjectKind().GroupVersionKind()})
	return nil
}

func (a *FakeAPIClient) DeleteObject(_ context.Context, obj client.Object) error {
	a.Operations = append(a.Operations, Operation{"DeleteObject", obj.GetObjectKind().GroupVersionKind()})
	return nil
}
