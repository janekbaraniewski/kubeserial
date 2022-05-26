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
// Code generated by deepcopy-gen. DO NOT EDIT.

package v1alpha1

import (
	runtime "k8s.io/apimachinery/pkg/runtime"
)

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Image) DeepCopyInto(out *Image) {
	*out = *in
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Image.
func (in *Image) DeepCopy() *Image {
	if in == nil {
		return nil
	}
	out := new(Image)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *IngressSpec) DeepCopyInto(out *IngressSpec) {
	*out = *in
	if in.Annotations != nil {
		in, out := &in.Annotations, &out.Annotations
		*out = make(map[string]string, len(*in))
		for key, val := range *in {
			(*out)[key] = val
		}
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new IngressSpec.
func (in *IngressSpec) DeepCopy() *IngressSpec {
	if in == nil {
		return nil
	}
	out := new(IngressSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *KubeSerial) DeepCopyInto(out *KubeSerial) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
	out.Status = in.Status
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new KubeSerial.
func (in *KubeSerial) DeepCopy() *KubeSerial {
	if in == nil {
		return nil
	}
	out := new(KubeSerial)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *KubeSerial) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *KubeSerialList) DeepCopyInto(out *KubeSerialList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]KubeSerial, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new KubeSerialList.
func (in *KubeSerialList) DeepCopy() *KubeSerialList {
	if in == nil {
		return nil
	}
	out := new(KubeSerialList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *KubeSerialList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *KubeSerialSpec) DeepCopyInto(out *KubeSerialSpec) {
	*out = *in
	if in.SerialDevices != nil {
		in, out := &in.SerialDevices, &out.SerialDevices
		*out = make([]SerialDevice_2, len(*in))
		copy(*out, *in)
	}
	in.Ingress.DeepCopyInto(&out.Ingress)
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new KubeSerialSpec.
func (in *KubeSerialSpec) DeepCopy() *KubeSerialSpec {
	if in == nil {
		return nil
	}
	out := new(KubeSerialSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *KubeSerialStatus) DeepCopyInto(out *KubeSerialStatus) {
	*out = *in
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new KubeSerialStatus.
func (in *KubeSerialStatus) DeepCopy() *KubeSerialStatus {
	if in == nil {
		return nil
	}
	out := new(KubeSerialStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Manager) DeepCopyInto(out *Manager) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	out.Spec = in.Spec
	out.Status = in.Status
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Manager.
func (in *Manager) DeepCopy() *Manager {
	if in == nil {
		return nil
	}
	out := new(Manager)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *Manager) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ManagerList) DeepCopyInto(out *ManagerList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]Manager, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ManagerList.
func (in *ManagerList) DeepCopy() *ManagerList {
	if in == nil {
		return nil
	}
	out := new(ManagerList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *ManagerList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ManagerScheduleRequest) DeepCopyInto(out *ManagerScheduleRequest) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	out.Spec = in.Spec
	out.Status = in.Status
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ManagerScheduleRequest.
func (in *ManagerScheduleRequest) DeepCopy() *ManagerScheduleRequest {
	if in == nil {
		return nil
	}
	out := new(ManagerScheduleRequest)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *ManagerScheduleRequest) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ManagerScheduleRequestList) DeepCopyInto(out *ManagerScheduleRequestList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]ManagerScheduleRequest, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ManagerScheduleRequestList.
func (in *ManagerScheduleRequestList) DeepCopy() *ManagerScheduleRequestList {
	if in == nil {
		return nil
	}
	out := new(ManagerScheduleRequestList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *ManagerScheduleRequestList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ManagerScheduleRequestSpec) DeepCopyInto(out *ManagerScheduleRequestSpec) {
	*out = *in
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ManagerScheduleRequestSpec.
func (in *ManagerScheduleRequestSpec) DeepCopy() *ManagerScheduleRequestSpec {
	if in == nil {
		return nil
	}
	out := new(ManagerScheduleRequestSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ManagerScheduleRequestStatus) DeepCopyInto(out *ManagerScheduleRequestStatus) {
	*out = *in
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ManagerScheduleRequestStatus.
func (in *ManagerScheduleRequestStatus) DeepCopy() *ManagerScheduleRequestStatus {
	if in == nil {
		return nil
	}
	out := new(ManagerScheduleRequestStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ManagerSpec) DeepCopyInto(out *ManagerSpec) {
	*out = *in
	out.Image = in.Image
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ManagerSpec.
func (in *ManagerSpec) DeepCopy() *ManagerSpec {
	if in == nil {
		return nil
	}
	out := new(ManagerSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ManagerStatus) DeepCopyInto(out *ManagerStatus) {
	*out = *in
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ManagerStatus.
func (in *ManagerStatus) DeepCopy() *ManagerStatus {
	if in == nil {
		return nil
	}
	out := new(ManagerStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *SerialDevice) DeepCopyInto(out *SerialDevice) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	out.Spec = in.Spec
	in.Status.DeepCopyInto(&out.Status)
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new SerialDevice.
func (in *SerialDevice) DeepCopy() *SerialDevice {
	if in == nil {
		return nil
	}
	out := new(SerialDevice)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *SerialDevice) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *SerialDeviceCondition) DeepCopyInto(out *SerialDeviceCondition) {
	*out = *in
	in.LastTransitionTime.DeepCopyInto(&out.LastTransitionTime)
	in.LastHeartbeatTime.DeepCopyInto(&out.LastHeartbeatTime)
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new SerialDeviceCondition.
func (in *SerialDeviceCondition) DeepCopy() *SerialDeviceCondition {
	if in == nil {
		return nil
	}
	out := new(SerialDeviceCondition)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *SerialDeviceList) DeepCopyInto(out *SerialDeviceList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]SerialDevice, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new SerialDeviceList.
func (in *SerialDeviceList) DeepCopy() *SerialDeviceList {
	if in == nil {
		return nil
	}
	out := new(SerialDeviceList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *SerialDeviceList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *SerialDeviceSpec) DeepCopyInto(out *SerialDeviceSpec) {
	*out = *in
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new SerialDeviceSpec.
func (in *SerialDeviceSpec) DeepCopy() *SerialDeviceSpec {
	if in == nil {
		return nil
	}
	out := new(SerialDeviceSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *SerialDeviceStatus) DeepCopyInto(out *SerialDeviceStatus) {
	*out = *in
	if in.Conditions != nil {
		in, out := &in.Conditions, &out.Conditions
		*out = make([]SerialDeviceCondition, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new SerialDeviceStatus.
func (in *SerialDeviceStatus) DeepCopy() *SerialDeviceStatus {
	if in == nil {
		return nil
	}
	out := new(SerialDeviceStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *SerialDevice_2) DeepCopyInto(out *SerialDevice_2) {
	*out = *in
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new SerialDevice_2.
func (in *SerialDevice_2) DeepCopy() *SerialDevice_2 {
	if in == nil {
		return nil
	}
	out := new(SerialDevice_2)
	in.DeepCopyInto(out)
	return out
}