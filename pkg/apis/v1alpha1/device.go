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

package v1alpha1

import (
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type SerialDeviceConditionType string

const (
	SerialDeviceAvailable SerialDeviceConditionType = "Available"
	SerialDeviceFree      SerialDeviceConditionType = "Free"
	SerialDeviceReady     SerialDeviceConditionType = "Ready"
)

// +k8s:openapi-gen=true
// SerialDeviceSpec defines the desired state of SerialDevice
type SerialDeviceSpec struct {
	// +required
	// +kubebuilder:validation:Required
	Name string `json:"name"`
	// +required
	// +kubebuilder:validation:Required
	IdVendor string `json:"idVendor"`
	// +required
	// +kubebuilder:validation:Required
	IdProduct string `json:"idProduct"`
	// +optional
	// +kubebuilder:validation:Optional
	Manager string `json:"manager,omitempty"`
}

type SerialDeviceCondition struct {
	// +required
	Type SerialDeviceConditionType `json:"type" protobuf:"bytes,1,opt,name=type"`
	// +required
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:Enum=True;False;Unknown
	Status metav1.ConditionStatus `json:"status" protobuf:"bytes,2,opt,name=status"`
	// +optional
	// +kubebuilder:validation:Minimum=0
	ObservedGeneration int64 `json:"observedGeneration,omitempty" protobuf:"varint,3,opt,name=observedGeneration"`
	// +required
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:Type=string
	// +kubebuilder:validation:Format=date-time
	LastTransitionTime metav1.Time `json:"lastTransitionTime" protobuf:"bytes,4,opt,name=lastTransitionTime"`
	// +optional
	// +kubebuilder:validation:Optional
	// +kubebuilder:validation:Type=string
	// +kubebuilder:validation:Format=date-time
	LastHeartbeatTime metav1.Time `json:"lastHeartbeatTime" protobuf:"bytes,4,opt,name=lastHeartbeatTime"`
	// +required
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:MaxLength=1024
	// +kubebuilder:validation:MinLength=1
	// +kubebuilder:validation:Pattern=`^[A-Za-z]([A-Za-z0-9_,:]*[A-Za-z0-9_])?$`
	Reason string `json:"reason" protobuf:"bytes,5,opt,name=reason"`
	// message is a human readable message indicating details about the transition.
	// This may be an empty string.
	// +required
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:MaxLength=32768
	Message string `json:"message" protobuf:"bytes,6,opt,name=message"`
}

// +k8s:openapi-gen=true
// SerialDeviceStatus defines the observed state of SerialDevice
type SerialDeviceStatus struct {
	Conditions []SerialDeviceCondition `json:"conditions"`
	NodeName   string                  `json:"nodeName,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +genclient:nonNamespaced
// +genclient
// +k8s:openapi-gen=true
// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:resource:path=serialdevices,scope=Cluster
// +kubebuilder:printcolumn:name="Ready",type=string,JSONPath=`.status.conditions[?(@.type=="Ready")].status`
// +kubebuilder:printcolumn:name="Available",type=string,JSONPath=`.status.conditions[?(@.type=="Available")].status`
// +kubebuilder:printcolumn:name="SerialDevice Node",type=string,JSONPath=`.status.nodeName`

// SerialDevice is the Schema for the SerialDevices API
type SerialDevice struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   SerialDeviceSpec   `json:"spec,omitempty"`
	Status SerialDeviceStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// SerialDeviceList contains a list of SerialDevice
type SerialDeviceList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []SerialDevice `json:"items"`
}

// NeedsManager checks if SerialDevice needs Manager
func (d *SerialDevice) NeedsManager() bool {
	return d.Spec.Manager != ""
}

func (d *SerialDevice) IsAvailable() bool {
	availableCondition := d.GetCondition(SerialDeviceAvailable)
	return availableCondition.Status == metav1.ConditionTrue
}

func (d *SerialDevice) IsReady() bool {
	readyCondition := d.GetCondition(SerialDeviceReady)
	return readyCondition.Status == metav1.ConditionTrue
}

func (d *SerialDevice) IsFree() bool {
	freeCondition := d.GetCondition(SerialDeviceFree)
	return freeCondition.Status == metav1.ConditionTrue
}

func (d *SerialDevice) GetCondition(conditionType SerialDeviceConditionType) *SerialDeviceCondition {
	for i := range d.Status.Conditions {
		if d.Status.Conditions[i].Type == conditionType {
			return &d.Status.Conditions[i]
		}
	}
	return nil
}

func (d *SerialDevice) SetCondition(newCondition SerialDeviceCondition) {
	existing := d.GetCondition(newCondition.Type)

	if existing == nil {
		if newCondition.LastTransitionTime.IsZero() {
			newCondition.LastTransitionTime = v1.NewTime(time.Now())
		}
		newCondition.LastHeartbeatTime = v1.NewTime(time.Now())
		d.Status.Conditions = append(d.Status.Conditions, newCondition)
		return
	}

	if existing.Status != newCondition.Status {
		existing.Status = newCondition.Status
		if !newCondition.LastTransitionTime.IsZero() {
			existing.LastTransitionTime = newCondition.LastTransitionTime
		} else {
			existing.LastTransitionTime = v1.NewTime(time.Now())
		}
	}

	existing.Reason = newCondition.Reason
	existing.Message = newCondition.Message
	existing.ObservedGeneration = newCondition.ObservedGeneration
	existing.LastHeartbeatTime = v1.NewTime(time.Now())
}
