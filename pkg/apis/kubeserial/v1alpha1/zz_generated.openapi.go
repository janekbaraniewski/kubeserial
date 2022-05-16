//go:build !ignore_autogenerated
// +build !ignore_autogenerated

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
// Code generated by openapi-gen. DO NOT EDIT.

// This file was autogenerated by openapi-gen. Do not edit it manually!

package v1alpha1

import (
	common "k8s.io/kube-openapi/pkg/common"
	spec "k8s.io/kube-openapi/pkg/validation/spec"
)

func GetOpenAPIDefinitions(ref common.ReferenceCallback) map[string]common.OpenAPIDefinition {
	return map[string]common.OpenAPIDefinition{
		"github.com/janekbaraniewski/kubeserial/pkg/apis/kubeserial/v1alpha1.DeviceSpec":                   schema_pkg_apis_kubeserial_v1alpha1_DeviceSpec(ref),
		"github.com/janekbaraniewski/kubeserial/pkg/apis/kubeserial/v1alpha1.DeviceStatus":                 schema_pkg_apis_kubeserial_v1alpha1_DeviceStatus(ref),
		"github.com/janekbaraniewski/kubeserial/pkg/apis/kubeserial/v1alpha1.Device_2":                     schema_pkg_apis_kubeserial_v1alpha1_Device_2(ref),
		"github.com/janekbaraniewski/kubeserial/pkg/apis/kubeserial/v1alpha1.IngressSpec":                  schema_pkg_apis_kubeserial_v1alpha1_IngressSpec(ref),
		"github.com/janekbaraniewski/kubeserial/pkg/apis/kubeserial/v1alpha1.KubeSerial":                   schema_pkg_apis_kubeserial_v1alpha1_KubeSerial(ref),
		"github.com/janekbaraniewski/kubeserial/pkg/apis/kubeserial/v1alpha1.KubeSerialSpec":               schema_pkg_apis_kubeserial_v1alpha1_KubeSerialSpec(ref),
		"github.com/janekbaraniewski/kubeserial/pkg/apis/kubeserial/v1alpha1.KubeSerialStatus":             schema_pkg_apis_kubeserial_v1alpha1_KubeSerialStatus(ref),
		"github.com/janekbaraniewski/kubeserial/pkg/apis/kubeserial/v1alpha1.ManagerScheduleRequestSpec":   schema_pkg_apis_kubeserial_v1alpha1_ManagerScheduleRequestSpec(ref),
		"github.com/janekbaraniewski/kubeserial/pkg/apis/kubeserial/v1alpha1.ManagerScheduleRequestStatus": schema_pkg_apis_kubeserial_v1alpha1_ManagerScheduleRequestStatus(ref),
		"github.com/janekbaraniewski/kubeserial/pkg/apis/kubeserial/v1alpha1.ManagerSpec":                  schema_pkg_apis_kubeserial_v1alpha1_ManagerSpec(ref),
		"github.com/janekbaraniewski/kubeserial/pkg/apis/kubeserial/v1alpha1.ManagerStatus":                schema_pkg_apis_kubeserial_v1alpha1_ManagerStatus(ref),
	}
}

func schema_pkg_apis_kubeserial_v1alpha1_DeviceSpec(ref common.ReferenceCallback) common.OpenAPIDefinition {
	return common.OpenAPIDefinition{
		Schema: spec.Schema{
			SchemaProps: spec.SchemaProps{
				Description: "DeviceSpec defines the desired state of Device",
				Type:        []string{"object"},
				Properties: map[string]spec.Schema{
					"name": {
						SchemaProps: spec.SchemaProps{
							Default: "",
							Type:    []string{"string"},
							Format:  "",
						},
					},
					"idVendor": {
						SchemaProps: spec.SchemaProps{
							Default: "",
							Type:    []string{"string"},
							Format:  "",
						},
					},
					"idProduct": {
						SchemaProps: spec.SchemaProps{
							Default: "",
							Type:    []string{"string"},
							Format:  "",
						},
					},
					"manager": {
						SchemaProps: spec.SchemaProps{
							Type:   []string{"string"},
							Format: "",
						},
					},
				},
				Required: []string{"name", "idVendor", "idProduct"},
			},
		},
	}
}

func schema_pkg_apis_kubeserial_v1alpha1_DeviceStatus(ref common.ReferenceCallback) common.OpenAPIDefinition {
	return common.OpenAPIDefinition{
		Schema: spec.Schema{
			SchemaProps: spec.SchemaProps{
				Description: "DeviceStatus defines the observed state of Device",
				Type:        []string{"object"},
				Properties: map[string]spec.Schema{
					"conditions": {
						SchemaProps: spec.SchemaProps{
							Type: []string{"array"},
							Items: &spec.SchemaOrArray{
								Schema: &spec.Schema{
									SchemaProps: spec.SchemaProps{
										Default: map[string]interface{}{},
										Ref:     ref("github.com/janekbaraniewski/kubeserial/pkg/apis/kubeserial/v1alpha1.DeviceCondition"),
									},
								},
							},
						},
					},
					"nodeName": {
						SchemaProps: spec.SchemaProps{
							Type:   []string{"string"},
							Format: "",
						},
					},
				},
				Required: []string{"conditions"},
			},
		},
		Dependencies: []string{
			"github.com/janekbaraniewski/kubeserial/pkg/apis/kubeserial/v1alpha1.DeviceCondition"},
	}
}

func schema_pkg_apis_kubeserial_v1alpha1_Device_2(ref common.ReferenceCallback) common.OpenAPIDefinition {
	return common.OpenAPIDefinition{
		Schema: spec.Schema{
			SchemaProps: spec.SchemaProps{
				Description: "Device defines monitored device",
				Type:        []string{"object"},
				Properties: map[string]spec.Schema{
					"name": {
						SchemaProps: spec.SchemaProps{
							Default: "",
							Type:    []string{"string"},
							Format:  "",
						},
					},
					"idVendor": {
						SchemaProps: spec.SchemaProps{
							Default: "",
							Type:    []string{"string"},
							Format:  "",
						},
					},
					"idProduct": {
						SchemaProps: spec.SchemaProps{
							Default: "",
							Type:    []string{"string"},
							Format:  "",
						},
					},
					"manager": {
						SchemaProps: spec.SchemaProps{
							Default: "",
							Type:    []string{"string"},
							Format:  "",
						},
					},
					"subsystem": {
						SchemaProps: spec.SchemaProps{
							Default: "",
							Type:    []string{"string"},
							Format:  "",
						},
					},
				},
				Required: []string{"name", "idVendor", "idProduct", "subsystem"},
			},
		},
	}
}

func schema_pkg_apis_kubeserial_v1alpha1_IngressSpec(ref common.ReferenceCallback) common.OpenAPIDefinition {
	return common.OpenAPIDefinition{
		Schema: spec.Schema{
			SchemaProps: spec.SchemaProps{
				Description: "IngressSpec defines the desired Ingress configuration",
				Type:        []string{"object"},
				Properties: map[string]spec.Schema{
					"enabled": {
						SchemaProps: spec.SchemaProps{
							Default: false,
							Type:    []string{"boolean"},
							Format:  "",
						},
					},
					"domain": {
						SchemaProps: spec.SchemaProps{
							Type:   []string{"string"},
							Format: "",
						},
					},
					"annotations": {
						SchemaProps: spec.SchemaProps{
							Type: []string{"object"},
							AdditionalProperties: &spec.SchemaOrBool{
								Allows: true,
								Schema: &spec.Schema{
									SchemaProps: spec.SchemaProps{
										Default: "",
										Type:    []string{"string"},
										Format:  "",
									},
								},
							},
						},
					},
				},
				Required: []string{"enabled"},
			},
		},
	}
}

func schema_pkg_apis_kubeserial_v1alpha1_KubeSerial(ref common.ReferenceCallback) common.OpenAPIDefinition {
	return common.OpenAPIDefinition{
		Schema: spec.Schema{
			SchemaProps: spec.SchemaProps{
				Description: "KubeSerial is the Schema for the kubeserials API",
				Type:        []string{"object"},
				Properties: map[string]spec.Schema{
					"kind": {
						SchemaProps: spec.SchemaProps{
							Description: "Kind is a string value representing the REST resource this object represents. Servers may infer this from the endpoint the client submits requests to. Cannot be updated. In CamelCase. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#types-kinds",
							Type:        []string{"string"},
							Format:      "",
						},
					},
					"apiVersion": {
						SchemaProps: spec.SchemaProps{
							Description: "APIVersion defines the versioned schema of this representation of an object. Servers should convert recognized schemas to the latest internal value, and may reject unrecognized values. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#resources",
							Type:        []string{"string"},
							Format:      "",
						},
					},
					"metadata": {
						SchemaProps: spec.SchemaProps{
							Default: map[string]interface{}{},
							Ref:     ref("k8s.io/apimachinery/pkg/apis/meta/v1.ObjectMeta"),
						},
					},
					"spec": {
						SchemaProps: spec.SchemaProps{
							Default: map[string]interface{}{},
							Ref:     ref("github.com/janekbaraniewski/kubeserial/pkg/apis/kubeserial/v1alpha1.KubeSerialSpec"),
						},
					},
					"status": {
						SchemaProps: spec.SchemaProps{
							Default: map[string]interface{}{},
							Ref:     ref("github.com/janekbaraniewski/kubeserial/pkg/apis/kubeserial/v1alpha1.KubeSerialStatus"),
						},
					},
				},
				Required: []string{"spec"},
			},
		},
		Dependencies: []string{
			"github.com/janekbaraniewski/kubeserial/pkg/apis/kubeserial/v1alpha1.KubeSerialSpec", "github.com/janekbaraniewski/kubeserial/pkg/apis/kubeserial/v1alpha1.KubeSerialStatus", "k8s.io/apimachinery/pkg/apis/meta/v1.ObjectMeta"},
	}
}

func schema_pkg_apis_kubeserial_v1alpha1_KubeSerialSpec(ref common.ReferenceCallback) common.OpenAPIDefinition {
	return common.OpenAPIDefinition{
		Schema: spec.Schema{
			SchemaProps: spec.SchemaProps{
				Description: "KubeSerialSpec defines the desired state of KubeSerial",
				Type:        []string{"object"},
				Properties: map[string]spec.Schema{
					"devices": {
						SchemaProps: spec.SchemaProps{
							Type: []string{"array"},
							Items: &spec.SchemaOrArray{
								Schema: &spec.Schema{
									SchemaProps: spec.SchemaProps{
										Default: map[string]interface{}{},
										Ref:     ref("github.com/janekbaraniewski/kubeserial/pkg/apis/kubeserial/v1alpha1.Device_2"),
									},
								},
							},
						},
					},
					"ingress": {
						SchemaProps: spec.SchemaProps{
							Default: map[string]interface{}{},
							Ref:     ref("github.com/janekbaraniewski/kubeserial/pkg/apis/kubeserial/v1alpha1.IngressSpec"),
						},
					},
				},
				Required: []string{"devices", "ingress"},
			},
		},
		Dependencies: []string{
			"github.com/janekbaraniewski/kubeserial/pkg/apis/kubeserial/v1alpha1.Device_2", "github.com/janekbaraniewski/kubeserial/pkg/apis/kubeserial/v1alpha1.IngressSpec"},
	}
}

func schema_pkg_apis_kubeserial_v1alpha1_KubeSerialStatus(ref common.ReferenceCallback) common.OpenAPIDefinition {
	return common.OpenAPIDefinition{
		Schema: spec.Schema{
			SchemaProps: spec.SchemaProps{
				Description: "KubeSerialStatus defines the observed state of KubeSerial",
				Type:        []string{"object"},
			},
		},
	}
}

func schema_pkg_apis_kubeserial_v1alpha1_ManagerScheduleRequestSpec(ref common.ReferenceCallback) common.OpenAPIDefinition {
	return common.OpenAPIDefinition{
		Schema: spec.Schema{
			SchemaProps: spec.SchemaProps{
				Description: "ManagerScheduleRequestSpec defines the desired state of ManagerScheduleRequest",
				Type:        []string{"object"},
				Properties: map[string]spec.Schema{
					"device": {
						SchemaProps: spec.SchemaProps{
							Default: "",
							Type:    []string{"string"},
							Format:  "",
						},
					},
					"manager": {
						SchemaProps: spec.SchemaProps{
							Default: "",
							Type:    []string{"string"},
							Format:  "",
						},
					},
				},
				Required: []string{"device", "manager"},
			},
		},
	}
}

func schema_pkg_apis_kubeserial_v1alpha1_ManagerScheduleRequestStatus(ref common.ReferenceCallback) common.OpenAPIDefinition {
	return common.OpenAPIDefinition{
		Schema: spec.Schema{
			SchemaProps: spec.SchemaProps{
				Description: "ManagerScheduleRequestStatus defines the observed state of ManagerScheduleRequest",
				Type:        []string{"object"},
				Properties: map[string]spec.Schema{
					"fulfilled": {
						SchemaProps: spec.SchemaProps{
							Default: false,
							Type:    []string{"boolean"},
							Format:  "",
						},
					},
				},
				Required: []string{"fulfilled"},
			},
		},
	}
}

func schema_pkg_apis_kubeserial_v1alpha1_ManagerSpec(ref common.ReferenceCallback) common.OpenAPIDefinition {
	return common.OpenAPIDefinition{
		Schema: spec.Schema{
			SchemaProps: spec.SchemaProps{
				Description: "ManagerSpec defines the desired state of Manager",
				Type:        []string{"object"},
				Properties: map[string]spec.Schema{
					"image": {
						SchemaProps: spec.SchemaProps{
							Default: map[string]interface{}{},
							Ref:     ref("github.com/janekbaraniewski/kubeserial/pkg/apis/kubeserial/v1alpha1.Image"),
						},
					},
					"runCmd": {
						SchemaProps: spec.SchemaProps{
							Default: "",
							Type:    []string{"string"},
							Format:  "",
						},
					},
					"config": {
						SchemaProps: spec.SchemaProps{
							Default: "",
							Type:    []string{"string"},
							Format:  "",
						},
					},
					"configPath": {
						SchemaProps: spec.SchemaProps{
							Default: "",
							Type:    []string{"string"},
							Format:  "",
						},
					},
				},
				Required: []string{"image", "runCmd"},
			},
		},
		Dependencies: []string{
			"github.com/janekbaraniewski/kubeserial/pkg/apis/kubeserial/v1alpha1.Image"},
	}
}

func schema_pkg_apis_kubeserial_v1alpha1_ManagerStatus(ref common.ReferenceCallback) common.OpenAPIDefinition {
	return common.OpenAPIDefinition{
		Schema: spec.Schema{
			SchemaProps: spec.SchemaProps{
				Description: "ManagerStatus defines the observed state of Manager",
				Type:        []string{"object"},
			},
		},
	}
}
