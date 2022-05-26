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

// Image is spec for image to be used by manager
type Image struct {
	// +required
	// +kubebuilder:validation:Required
	Repository string `json:"repository"`
	// +required
	// +kubebuilder:validation:Required
	Tag string `json:"tag"`
}

// +k8s:openapi-gen=true
// ManagerSpec defines the desired state of Manager
type ManagerSpec struct {
	// +required
	// +kubebuilder:validation:Required
	Image Image `json:"image"`
	// +required
	// +kubebuilder:validation:Required
	RunCmd string `json:"runCmd"`
	// +optional
	// +kubebuilder:validation:Optional
	Config string `json:"config"`
	// +optional
	// +kubebuilder:validation:Optional
	ConfigPath string `json:"configPath"`
}

// +k8s:openapi-gen=true
// ManagerStatus defines the observed state of Manager
type ManagerStatus struct{}

// +genclient
// +k8s:openapi-gen=true
// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:resource:path=managers,scope=Namespaced
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// Manager is the Schema for the managers API
type Manager struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ManagerSpec   `json:"spec,omitempty"`
	Status ManagerStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// ManagerList contains a list of Manager
type ManagerList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Manager `json:"items"`
}
