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

// Code generated by validation-gen. DO NOT EDIT.

package primitive_pointers

import (
	fmt "fmt"

	operation "k8s.io/apimachinery/pkg/api/operation"
	safe "k8s.io/apimachinery/pkg/api/safe"
	validate "k8s.io/apimachinery/pkg/api/validate"
	runtime "k8s.io/apimachinery/pkg/runtime"
	field "k8s.io/apimachinery/pkg/util/validation/field"
)

func init() { localSchemeBuilder.Register(RegisterValidations) }

// RegisterValidations adds validation functions to the given scheme.
// Public to allow building arbitrary schemes.
func RegisterValidations(scheme *runtime.Scheme) error {
	scheme.AddValidationFunc((*T1)(nil), func(opCtx operation.Context, obj, oldObj interface{}, subresources ...string) field.ErrorList {
		if len(subresources) == 0 {
			return Validate_T1(opCtx, obj.(*T1), safe.Cast[T1](oldObj), nil)
		}
		return field.ErrorList{field.InternalError(nil, fmt.Errorf("no validation found for %T, subresources: %v", obj, subresources))}
	})
	return nil
}

func Validate_T1(opCtx operation.Context, obj, oldObj *T1, fldPath *field.Path) (errs field.ErrorList) {
	// field T1.SP
	errs = append(errs,
		func(obj *string, oldObj *string, fldPath *field.Path) (errs field.ErrorList) {
			if obj != nil {
				if vContext.Operation == k8s.io/apimachinery/pkg/api/operation.Update && oldObj != nil {
					errs = append(errs, validate.FixedResultUpdate(fldPath, *obj, *oldObj, true, "T1.SP, UpdateOnly")...)
				}
			}
			return
		}(obj.SP, safe.Field(oldObj, func(oldObj T1) *string { return oldObj.SP }), fldPath.Child("sp"))...)

	// field T1.IP
	errs = append(errs,
		func(obj *int, oldObj *int, fldPath *field.Path) (errs field.ErrorList) {
			if obj != nil {
				if vContext.Operation == k8s.io/apimachinery/pkg/api/operation.Update && oldObj != nil {
					errs = append(errs, validate.FixedResultUpdate(fldPath, *obj, *oldObj, true, "T1.IP, UpdateOnly")...)
				}
			}
			return
		}(obj.IP, safe.Field(oldObj, func(oldObj T1) *int { return oldObj.IP }), fldPath.Child("ip"))...)

	// field T1.BP
	errs = append(errs,
		func(obj *bool, oldObj *bool, fldPath *field.Path) (errs field.ErrorList) {
			if obj != nil {
				if vContext.Operation == k8s.io/apimachinery/pkg/api/operation.Update && oldObj != nil {
					errs = append(errs, validate.FixedResultUpdate(fldPath, *obj, *oldObj, true, "T1.BP, UpdateOnly")...)
				}
			}
			return
		}(obj.BP, safe.Field(oldObj, func(oldObj T1) *bool { return oldObj.BP }), fldPath.Child("bp"))...)

	// field T1.FP
	errs = append(errs,
		func(obj *float64, oldObj *float64, fldPath *field.Path) (errs field.ErrorList) {
			if obj != nil {
				if vContext.Operation == k8s.io/apimachinery/pkg/api/operation.Update && oldObj != nil {
					errs = append(errs, validate.FixedResultUpdate(fldPath, *obj, *oldObj, true, "T1.FP, UpdateOnly")...)
				}
			}
			return
		}(obj.FP, safe.Field(oldObj, func(oldObj T1) *float64 { return oldObj.FP }), fldPath.Child("fp"))...)

	return errs
}
