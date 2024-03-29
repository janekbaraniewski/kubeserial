/*
Copyright 2022.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
// Code generated by client-gen. DO NOT EDIT.

package fake

import (
	"context"

	v1alpha1 "github.com/janekbaraniewski/kubeserial/pkg/apis/v1alpha1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	labels "k8s.io/apimachinery/pkg/labels"
	schema "k8s.io/apimachinery/pkg/runtime/schema"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	testing "k8s.io/client-go/testing"
)

// FakeManagerScheduleRequests implements ManagerScheduleRequestInterface
type FakeManagerScheduleRequests struct {
	Fake *FakeAppV1alpha1
}

var managerschedulerequestsResource = schema.GroupVersionResource{Group: "app.kubeserial.com", Version: "v1alpha1", Resource: "managerschedulerequests"}

var managerschedulerequestsKind = schema.GroupVersionKind{Group: "app.kubeserial.com", Version: "v1alpha1", Kind: "ManagerScheduleRequest"}

// Get takes name of the managerScheduleRequest, and returns the corresponding managerScheduleRequest object, and an error if there is any.
func (c *FakeManagerScheduleRequests) Get(ctx context.Context, name string, options v1.GetOptions) (result *v1alpha1.ManagerScheduleRequest, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootGetAction(managerschedulerequestsResource, name), &v1alpha1.ManagerScheduleRequest{})
	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.ManagerScheduleRequest), err
}

// List takes label and field selectors, and returns the list of ManagerScheduleRequests that match those selectors.
func (c *FakeManagerScheduleRequests) List(ctx context.Context, opts v1.ListOptions) (result *v1alpha1.ManagerScheduleRequestList, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootListAction(managerschedulerequestsResource, managerschedulerequestsKind, opts), &v1alpha1.ManagerScheduleRequestList{})
	if obj == nil {
		return nil, err
	}

	label, _, _ := testing.ExtractFromListOptions(opts)
	if label == nil {
		label = labels.Everything()
	}
	list := &v1alpha1.ManagerScheduleRequestList{ListMeta: obj.(*v1alpha1.ManagerScheduleRequestList).ListMeta}
	for _, item := range obj.(*v1alpha1.ManagerScheduleRequestList).Items {
		if label.Matches(labels.Set(item.Labels)) {
			list.Items = append(list.Items, item)
		}
	}
	return list, err
}

// Watch returns a watch.Interface that watches the requested managerScheduleRequests.
func (c *FakeManagerScheduleRequests) Watch(ctx context.Context, opts v1.ListOptions) (watch.Interface, error) {
	return c.Fake.
		InvokesWatch(testing.NewRootWatchAction(managerschedulerequestsResource, opts))
}

// Create takes the representation of a managerScheduleRequest and creates it.  Returns the server's representation of the managerScheduleRequest, and an error, if there is any.
func (c *FakeManagerScheduleRequests) Create(ctx context.Context, managerScheduleRequest *v1alpha1.ManagerScheduleRequest, opts v1.CreateOptions) (result *v1alpha1.ManagerScheduleRequest, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootCreateAction(managerschedulerequestsResource, managerScheduleRequest), &v1alpha1.ManagerScheduleRequest{})
	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.ManagerScheduleRequest), err
}

// Update takes the representation of a managerScheduleRequest and updates it. Returns the server's representation of the managerScheduleRequest, and an error, if there is any.
func (c *FakeManagerScheduleRequests) Update(ctx context.Context, managerScheduleRequest *v1alpha1.ManagerScheduleRequest, opts v1.UpdateOptions) (result *v1alpha1.ManagerScheduleRequest, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootUpdateAction(managerschedulerequestsResource, managerScheduleRequest), &v1alpha1.ManagerScheduleRequest{})
	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.ManagerScheduleRequest), err
}

// UpdateStatus was generated because the type contains a Status member.
// Add a +genclient:noStatus comment above the type to avoid generating UpdateStatus().
func (c *FakeManagerScheduleRequests) UpdateStatus(ctx context.Context, managerScheduleRequest *v1alpha1.ManagerScheduleRequest, opts v1.UpdateOptions) (*v1alpha1.ManagerScheduleRequest, error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootUpdateSubresourceAction(managerschedulerequestsResource, "status", managerScheduleRequest), &v1alpha1.ManagerScheduleRequest{})
	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.ManagerScheduleRequest), err
}

// Delete takes name of the managerScheduleRequest and deletes it. Returns an error if one occurs.
func (c *FakeManagerScheduleRequests) Delete(ctx context.Context, name string, opts v1.DeleteOptions) error {
	_, err := c.Fake.
		Invokes(testing.NewRootDeleteActionWithOptions(managerschedulerequestsResource, name, opts), &v1alpha1.ManagerScheduleRequest{})
	return err
}

// DeleteCollection deletes a collection of objects.
func (c *FakeManagerScheduleRequests) DeleteCollection(ctx context.Context, opts v1.DeleteOptions, listOpts v1.ListOptions) error {
	action := testing.NewRootDeleteCollectionAction(managerschedulerequestsResource, listOpts)

	_, err := c.Fake.Invokes(action, &v1alpha1.ManagerScheduleRequestList{})
	return err
}

// Patch applies the patch and returns the patched managerScheduleRequest.
func (c *FakeManagerScheduleRequests) Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts v1.PatchOptions, subresources ...string) (result *v1alpha1.ManagerScheduleRequest, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootPatchSubresourceAction(managerschedulerequestsResource, name, pt, data, subresources...), &v1alpha1.ManagerScheduleRequest{})
	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.ManagerScheduleRequest), err
}
