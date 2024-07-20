//go:build !ignore_autogenerated
// +build !ignore_autogenerated

/*
Copyright The Kubernetes Authors.

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
func (in *Fischer) DeepCopyInto(out *Fischer) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	if in.DisallowedFlunders != nil {
		in, out := &in.DisallowedFlunders, &out.DisallowedFlunders
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
	in.Primary.DeepCopyInto(&out.Primary)
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Fischer.
func (in *Fischer) DeepCopy() *Fischer {
	if in == nil {
		return nil
	}
	out := new(Fischer)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *Fischer) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *FischerList) DeepCopyInto(out *FischerList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]Fischer, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new FischerList.
func (in *FischerList) DeepCopy() *FischerList {
	if in == nil {
		return nil
	}
	out := new(FischerList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *FischerList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Flunder) DeepCopyInto(out *Flunder) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
	out.Status = in.Status
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Flunder.
func (in *Flunder) DeepCopy() *Flunder {
	if in == nil {
		return nil
	}
	out := new(Flunder)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *Flunder) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *FlunderList) DeepCopyInto(out *FlunderList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]Flunder, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new FlunderList.
func (in *FlunderList) DeepCopy() *FlunderList {
	if in == nil {
		return nil
	}
	out := new(FlunderList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *FlunderList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *FlunderSpec) DeepCopyInto(out *FlunderSpec) {
	*out = *in
	if in.ReferenceType != nil {
		in, out := &in.ReferenceType, &out.ReferenceType
		*out = new(ReferenceType)
		**out = **in
	}
	in.Primary.DeepCopyInto(&out.Primary)
	if in.Extras != nil {
		in, out := &in.Extras, &out.Extras
		*out = make([]Widget, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	if in.More != nil {
		in, out := &in.More, &out.More
		*out = make(map[string]Widget, len(*in))
		for key, val := range *in {
			(*out)[key] = *val.DeepCopy()
		}
	}
	if in.Layer != nil {
		in, out := &in.Layer, &out.Layer
		*out = new(Layer)
		(*in).DeepCopyInto(*out)
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new FlunderSpec.
func (in *FlunderSpec) DeepCopy() *FlunderSpec {
	if in == nil {
		return nil
	}
	out := new(FlunderSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *FlunderStatus) DeepCopyInto(out *FlunderStatus) {
	*out = *in
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new FlunderStatus.
func (in *FlunderStatus) DeepCopy() *FlunderStatus {
	if in == nil {
		return nil
	}
	out := new(FlunderStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Layer) DeepCopyInto(out *Layer) {
	*out = *in
	if in.Extras != nil {
		in, out := &in.Extras, &out.Extras
		*out = make([]Widget, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Layer.
func (in *Layer) DeepCopy() *Layer {
	if in == nil {
		return nil
	}
	out := new(Layer)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Something) DeepCopyInto(out *Something) {
	*out = *in
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Something.
func (in *Something) DeepCopy() *Something {
	if in == nil {
		return nil
	}
	out := new(Something)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Widget) DeepCopyInto(out *Widget) {
	*out = *in
	if in.Something != nil {
		in, out := &in.Something, &out.Something
		*out = make([]Something, len(*in))
		copy(*out, *in)
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Widget.
func (in *Widget) DeepCopy() *Widget {
	if in == nil {
		return nil
	}
	out := new(Widget)
	in.DeepCopyInto(out)
	return out
}
