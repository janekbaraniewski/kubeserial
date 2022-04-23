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
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// Device defines monitored device
// +k8s:openapi-gen=true
type Device_2 struct {
	Name      string `json:"name"`
	IdVendor  string `json:"idVendor"`
	IdProduct string `json:"idProduct"`
	Manager   string `json:"manager"`
	Subsystem string `json:"subsystem"`
}

// IngressSpec defines the desired Ingress configuration
// +k8s:openapi-gen=true
type IngressSpec struct {
	Enabled     bool              `json:"enabled"`
	Domain      string            `json:"domain,omitempty"`
	Annotations map[string]string `json:"annotations,omitempty"`
}

// KubeSerialSpec defines the desired state of KubeSerial
// +k8s:openapi-gen=true
type KubeSerialSpec struct {
	Devices []Device_2  `json:"devices"`
	Ingress IngressSpec `json:"ingress"`
}

// KubeSerialStatus defines the observed state of KubeSerial
// +k8s:openapi-gen=true
type KubeSerialStatus struct{}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// KubeSerial is the Schema for the kubeserials API
// +genclient
// +k8s:openapi-gen=true
// +kubebuilder:subresource:status
// +kubebuilder:resource:path=kubeserials,scope=Namespaced
type KubeSerial struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   KubeSerialSpec   `json:"spec"`
	Status KubeSerialStatus `json:"status,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// KubeSerialList contains a list of KubeSerial
type KubeSerialList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []KubeSerial `json:"items"`
}
