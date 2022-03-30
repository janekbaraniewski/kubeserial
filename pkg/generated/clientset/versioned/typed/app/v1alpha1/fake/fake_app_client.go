// Code generated by client-gen. DO NOT EDIT.

package fake

import (
	v1alpha1 "github.com/janekbaraniewski/kubeserial/pkg/generated/clientset/versioned/typed/app/v1alpha1"
	rest "k8s.io/client-go/rest"
	testing "k8s.io/client-go/testing"
)

type FakeAppV1alpha1 struct {
	*testing.Fake
}

func (c *FakeAppV1alpha1) KubeSerials(namespace string) v1alpha1.KubeSerialInterface {
	return &FakeKubeSerials{c, namespace}
}

// RESTClient returns a RESTClient that is used to communicate
// with API server by this client implementation.
func (c *FakeAppV1alpha1) RESTClient() rest.Interface {
	var ret *rest.RESTClient
	return ret
}