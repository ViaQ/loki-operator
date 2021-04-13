// +build !ignore_autogenerated

/*
Copyright 2021.

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

// Code generated by controller-gen. DO NOT EDIT.

package v1beta1

import (
	"k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	runtime "k8s.io/apimachinery/pkg/runtime"
)

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *IngestionLimitSpec) DeepCopyInto(out *IngestionLimitSpec) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new IngestionLimitSpec.
func (in *IngestionLimitSpec) DeepCopy() *IngestionLimitSpec {
	if in == nil {
		return nil
	}
	out := new(IngestionLimitSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *LimitsSpec) DeepCopyInto(out *LimitsSpec) {
	*out = *in
	out.Global = in.Global
	if in.Tenants != nil {
		in, out := &in.Tenants, &out.Tenants
		*out = make(map[string]LimitsTemplateSpec, len(*in))
		for key, val := range *in {
			(*out)[key] = val
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new LimitsSpec.
func (in *LimitsSpec) DeepCopy() *LimitsSpec {
	if in == nil {
		return nil
	}
	out := new(LimitsSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *LimitsTemplateSpec) DeepCopyInto(out *LimitsTemplateSpec) {
	*out = *in
	out.IngestionLimits = in.IngestionLimits
	out.QueryLimits = in.QueryLimits
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new LimitsTemplateSpec.
func (in *LimitsTemplateSpec) DeepCopy() *LimitsTemplateSpec {
	if in == nil {
		return nil
	}
	out := new(LimitsTemplateSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *LokiComponentSpec) DeepCopyInto(out *LokiComponentSpec) {
	*out = *in
	if in.NodeSelector != nil {
		in, out := &in.NodeSelector, &out.NodeSelector
		*out = make(map[string]string, len(*in))
		for key, val := range *in {
			(*out)[key] = val
		}
	}
	if in.Tolerations != nil {
		in, out := &in.Tolerations, &out.Tolerations
		*out = make([]v1.Toleration, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new LokiComponentSpec.
func (in *LokiComponentSpec) DeepCopy() *LokiComponentSpec {
	if in == nil {
		return nil
	}
	out := new(LokiComponentSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *LokiStack) DeepCopyInto(out *LokiStack) {
	*out = *in
	in.Spec.DeepCopyInto(&out.Spec)
	in.Status.DeepCopyInto(&out.Status)
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	out.TypeMeta = in.TypeMeta
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new LokiStack.
func (in *LokiStack) DeepCopy() *LokiStack {
	if in == nil {
		return nil
	}
	out := new(LokiStack)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *LokiStack) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *LokiStackList) DeepCopyInto(out *LokiStackList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]LokiStack, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new LokiStackList.
func (in *LokiStackList) DeepCopy() *LokiStackList {
	if in == nil {
		return nil
	}
	out := new(LokiStackList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *LokiStackList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *LokiStackSpec) DeepCopyInto(out *LokiStackSpec) {
	*out = *in
	out.Storage = in.Storage
	in.Limits.DeepCopyInto(&out.Limits)
	in.Template.DeepCopyInto(&out.Template)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new LokiStackSpec.
func (in *LokiStackSpec) DeepCopy() *LokiStackSpec {
	if in == nil {
		return nil
	}
	out := new(LokiStackSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *LokiStackStatus) DeepCopyInto(out *LokiStackStatus) {
	*out = *in
	if in.Conditions != nil {
		in, out := &in.Conditions, &out.Conditions
		*out = make([]metav1.Condition, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new LokiStackStatus.
func (in *LokiStackStatus) DeepCopy() *LokiStackStatus {
	if in == nil {
		return nil
	}
	out := new(LokiStackStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *LokiTemplateSpec) DeepCopyInto(out *LokiTemplateSpec) {
	*out = *in
	in.Compactor.DeepCopyInto(&out.Compactor)
	in.Distributor.DeepCopyInto(&out.Distributor)
	in.Ingester.DeepCopyInto(&out.Ingester)
	in.Querier.DeepCopyInto(&out.Querier)
	in.QueryFrontend.DeepCopyInto(&out.QueryFrontend)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new LokiTemplateSpec.
func (in *LokiTemplateSpec) DeepCopy() *LokiTemplateSpec {
	if in == nil {
		return nil
	}
	out := new(LokiTemplateSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ObjectStorageSecretSpec) DeepCopyInto(out *ObjectStorageSecretSpec) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ObjectStorageSecretSpec.
func (in *ObjectStorageSecretSpec) DeepCopy() *ObjectStorageSecretSpec {
	if in == nil {
		return nil
	}
	out := new(ObjectStorageSecretSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ObjectStorageSpec) DeepCopyInto(out *ObjectStorageSpec) {
	*out = *in
	out.Secret = in.Secret
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ObjectStorageSpec.
func (in *ObjectStorageSpec) DeepCopy() *ObjectStorageSpec {
	if in == nil {
		return nil
	}
	out := new(ObjectStorageSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *QueryLimitSpec) DeepCopyInto(out *QueryLimitSpec) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new QueryLimitSpec.
func (in *QueryLimitSpec) DeepCopy() *QueryLimitSpec {
	if in == nil {
		return nil
	}
	out := new(QueryLimitSpec)
	in.DeepCopyInto(out)
	return out
}
