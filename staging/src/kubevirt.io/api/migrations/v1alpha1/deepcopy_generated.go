//go:build !ignore_autogenerated
// +build !ignore_autogenerated

/*
This file is part of the KubeVirt project

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.

Copyright The KubeVirt Authors.
*/

// Code generated by deepcopy-gen. DO NOT EDIT.

package v1alpha1

import (
	runtime "k8s.io/apimachinery/pkg/runtime"
)

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in LabelSelector) DeepCopyInto(out *LabelSelector) {
	{
		in := &in
		*out = make(LabelSelector, len(*in))
		for key, val := range *in {
			(*out)[key] = val
		}
		return
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new LabelSelector.
func (in LabelSelector) DeepCopy() LabelSelector {
	if in == nil {
		return nil
	}
	out := new(LabelSelector)
	in.DeepCopyInto(out)
	return *out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *MigrationPolicy) DeepCopyInto(out *MigrationPolicy) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
	out.Status = in.Status
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new MigrationPolicy.
func (in *MigrationPolicy) DeepCopy() *MigrationPolicy {
	if in == nil {
		return nil
	}
	out := new(MigrationPolicy)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *MigrationPolicy) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *MigrationPolicyList) DeepCopyInto(out *MigrationPolicyList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]MigrationPolicy, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new MigrationPolicyList.
func (in *MigrationPolicyList) DeepCopy() *MigrationPolicyList {
	if in == nil {
		return nil
	}
	out := new(MigrationPolicyList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *MigrationPolicyList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *MigrationPolicySpec) DeepCopyInto(out *MigrationPolicySpec) {
	*out = *in
	if in.Selectors != nil {
		in, out := &in.Selectors, &out.Selectors
		*out = new(Selectors)
		(*in).DeepCopyInto(*out)
	}
	if in.AllowAutoConverge != nil {
		in, out := &in.AllowAutoConverge, &out.AllowAutoConverge
		*out = new(bool)
		**out = **in
	}
	if in.BandwidthPerMigration != nil {
		in, out := &in.BandwidthPerMigration, &out.BandwidthPerMigration
		x := (*in).DeepCopy()
		*out = &x
	}
	if in.CompletionTimeoutPerGiB != nil {
		in, out := &in.CompletionTimeoutPerGiB, &out.CompletionTimeoutPerGiB
		*out = new(int64)
		**out = **in
	}
	if in.AllowPostCopy != nil {
		in, out := &in.AllowPostCopy, &out.AllowPostCopy
		*out = new(bool)
		**out = **in
	}
	if in.AllowWorkloadDisruption != nil {
		in, out := &in.AllowWorkloadDisruption, &out.AllowWorkloadDisruption
		*out = new(bool)
		**out = **in
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new MigrationPolicySpec.
func (in *MigrationPolicySpec) DeepCopy() *MigrationPolicySpec {
	if in == nil {
		return nil
	}
	out := new(MigrationPolicySpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *MigrationPolicyStatus) DeepCopyInto(out *MigrationPolicyStatus) {
	*out = *in
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new MigrationPolicyStatus.
func (in *MigrationPolicyStatus) DeepCopy() *MigrationPolicyStatus {
	if in == nil {
		return nil
	}
	out := new(MigrationPolicyStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Selectors) DeepCopyInto(out *Selectors) {
	*out = *in
	if in.NamespaceSelector != nil {
		in, out := &in.NamespaceSelector, &out.NamespaceSelector
		*out = make(LabelSelector, len(*in))
		for key, val := range *in {
			(*out)[key] = val
		}
	}
	if in.VirtualMachineInstanceSelector != nil {
		in, out := &in.VirtualMachineInstanceSelector, &out.VirtualMachineInstanceSelector
		*out = make(LabelSelector, len(*in))
		for key, val := range *in {
			(*out)[key] = val
		}
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Selectors.
func (in *Selectors) DeepCopy() *Selectors {
	if in == nil {
		return nil
	}
	out := new(Selectors)
	in.DeepCopyInto(out)
	return out
}
