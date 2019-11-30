package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type Device struct {
	Name 		string 	`json:"name"`
	IDVendor	string 	`json:"idvendor"`
	IDProduct	string 	`json:"idproduct"`
	Manager		string 	`json:"manager"`
}

// KubeSerialSpec defines the desired state of KubeSerial
// +k8s:openapi-gen=true
type KubeSerialSpec struct {
	Devices []Device `json:"devices"`
}

// KubeSerialStatus defines the observed state of KubeSerial
// +k8s:openapi-gen=true
type KubeSerialStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "operator-sdk generate k8s" to regenerate code after modifying this file
	// Add custom validation using kubebuilder tags: https://book-v1.book.kubebuilder.io/beyond_basics/generating_crd.html
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// KubeSerial is the Schema for the kubeserials API
// +k8s:openapi-gen=true
// +kubebuilder:subresource:status
// +kubebuilder:resource:path=kubeserials,scope=Namespaced
type KubeSerial struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   KubeSerialSpec   `json:"spec,omitempty"`
	Status KubeSerialStatus `json:"status,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// KubeSerialList contains a list of KubeSerial
type KubeSerialList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []KubeSerial `json:"items"`
}

func init() {
	SchemeBuilder.Register(&KubeSerial{}, &KubeSerialList{})
}
