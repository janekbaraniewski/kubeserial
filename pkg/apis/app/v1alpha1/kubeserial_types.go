package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// Device defines monitored device
// +k8s:openapi-gen=true
type Device struct {
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
	// +listType=set
	Devices []Device    `json:"devices"`
	Ingress IngressSpec `json:"ingress"`
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
