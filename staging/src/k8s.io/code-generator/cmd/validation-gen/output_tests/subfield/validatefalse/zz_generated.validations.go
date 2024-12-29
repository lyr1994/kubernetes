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

package validatefalse

import (
	fmt "fmt"

	operation "k8s.io/apimachinery/pkg/api/operation"
	safe "k8s.io/apimachinery/pkg/api/safe"
	validate "k8s.io/apimachinery/pkg/api/validate"
	field "k8s.io/apimachinery/pkg/util/validation/field"
	testscheme "k8s.io/code-generator/cmd/validation-gen/testscheme"
)

func init() { localSchemeBuilder.Register(RegisterValidations) }

// RegisterValidations adds validation functions to the given scheme.
// Public to allow building arbitrary schemes.
func RegisterValidations(scheme *testscheme.Scheme) error {
	scheme.AddValidationFunc((*T1)(nil), func(opCtx operation.Context, obj, oldObj interface{}, subresources ...string) field.ErrorList {
		if len(subresources) == 0 {
			return Validate_T1(opCtx, obj.(*T1), safe.Cast[*T1](oldObj), nil)
		}
		return field.ErrorList{field.InternalError(nil, fmt.Errorf("no validation found for %T, subresources: %v", obj, subresources))}
	})
	return nil
}

func Validate_T1(opCtx operation.Context, obj, oldObj *T1, fldPath *field.Path) (errs field.ErrorList) {
	// type T1
	errs = append(errs, validate.FixedResult(opCtx, fldPath, obj, oldObj, false, "type T1")...)

	// field T1.TypeMeta has no validation

	// field T1.T2
	errs = append(errs,
		func(obj, oldObj *T2, fldPath *field.Path) (errs field.ErrorList) {
			// field T2.MapField
			errs = append(errs,
				func(obj, oldObj map[string]string, fldPath *field.Path) (errs field.ErrorList) {
					errs = append(errs, validate.FixedResult(opCtx, fldPath, obj, oldObj, false, "subfield T1.T2.MapField")...)
					return
				}(obj.MapField, safe.Field(oldObj, func(oldObj *T2) map[string]string { return oldObj.MapField }), fldPath.Child("mapField"))...)

			// field T2.PointerField
			errs = append(errs,
				func(obj, oldObj *string, fldPath *field.Path) (errs field.ErrorList) {
					errs = append(errs, validate.FixedResult(opCtx, fldPath, obj, oldObj, false, "subfield T1.T2.PointerField")...)
					return
				}(obj.PointerField, safe.Field(oldObj, func(oldObj *T2) *string { return oldObj.PointerField }), fldPath.Child("pointerField"))...)

			// field T2.SliceField
			errs = append(errs,
				func(obj, oldObj []string, fldPath *field.Path) (errs field.ErrorList) {
					errs = append(errs, validate.FixedResult(opCtx, fldPath, obj, oldObj, false, "subfield T1.T2.SliceField")...)
					return
				}(obj.SliceField, safe.Field(oldObj, func(oldObj *T2) []string { return oldObj.SliceField }), fldPath.Child("sliceField"))...)

			// field T2.StringField
			errs = append(errs,
				func(obj, oldObj *string, fldPath *field.Path) (errs field.ErrorList) {
					errs = append(errs, validate.FixedResult(opCtx, fldPath, obj, oldObj, false, "subfield T1.T2.StringField")...)
					return
				}(&obj.StringField, safe.Field(oldObj, func(oldObj *T2) *string { return &oldObj.StringField }), fldPath.Child("stringField"))...)

			// field T2.StringFieldWithValidation
			errs = append(errs,
				func(obj, oldObj *string, fldPath *field.Path) (errs field.ErrorList) {
					errs = append(errs, validate.FixedResult(opCtx, fldPath, obj, oldObj, false, "subfield T1.T2.StringFieldWithValidation")...)
					return
				}(&obj.StringFieldWithValidation, safe.Field(oldObj, func(oldObj *T2) *string { return &oldObj.StringFieldWithValidation }), fldPath.Child("stringFieldWithValidation"))...)

			// field T2.StructField
			errs = append(errs,
				func(obj, oldObj *StructField, fldPath *field.Path) (errs field.ErrorList) {
					errs = append(errs, validate.FixedResult(opCtx, fldPath, obj, oldObj, false, "subfield T1.T2.StructField")...)
					return
				}(&obj.StructField, safe.Field(oldObj, func(oldObj *T2) *StructField { return &oldObj.StructField }), fldPath.Child("structField"))...)

			errs = append(errs, Validate_T2(opCtx, obj, oldObj, fldPath)...)
			return
		}(&obj.T2, safe.Field(oldObj, func(oldObj *T1) *T2 { return &oldObj.T2 }), fldPath.Child("t2"))...)

	// field T1.PT2
	errs = append(errs,
		func(obj, oldObj *T2, fldPath *field.Path) (errs field.ErrorList) {
			// field T2.MapField
			errs = append(errs,
				func(obj, oldObj map[string]string, fldPath *field.Path) (errs field.ErrorList) {
					errs = append(errs, validate.FixedResult(opCtx, fldPath, obj, oldObj, false, "subfield T1.PT2.MapField")...)
					return
				}(obj.MapField, safe.Field(oldObj, func(oldObj *T2) map[string]string { return oldObj.MapField }), fldPath.Child("mapField"))...)

			// field T2.PointerField
			errs = append(errs,
				func(obj, oldObj *string, fldPath *field.Path) (errs field.ErrorList) {
					errs = append(errs, validate.FixedResult(opCtx, fldPath, obj, oldObj, false, "subfield T1.PT2.PointerField")...)
					return
				}(obj.PointerField, safe.Field(oldObj, func(oldObj *T2) *string { return oldObj.PointerField }), fldPath.Child("pointerField"))...)

			// field T2.SliceField
			errs = append(errs,
				func(obj, oldObj []string, fldPath *field.Path) (errs field.ErrorList) {
					errs = append(errs, validate.FixedResult(opCtx, fldPath, obj, oldObj, false, "subfield T1.PT2.SliceField")...)
					return
				}(obj.SliceField, safe.Field(oldObj, func(oldObj *T2) []string { return oldObj.SliceField }), fldPath.Child("sliceField"))...)

			// field T2.StringField
			errs = append(errs,
				func(obj, oldObj *string, fldPath *field.Path) (errs field.ErrorList) {
					errs = append(errs, validate.FixedResult(opCtx, fldPath, obj, oldObj, false, "subfield T1.PT2.StringField")...)
					return
				}(&obj.StringField, safe.Field(oldObj, func(oldObj *T2) *string { return &oldObj.StringField }), fldPath.Child("stringField"))...)

			// field T2.StringFieldWithValidation
			errs = append(errs,
				func(obj, oldObj *string, fldPath *field.Path) (errs field.ErrorList) {
					errs = append(errs, validate.FixedResult(opCtx, fldPath, obj, oldObj, false, "subfield T1.PT2.StringFieldWithValidation")...)
					return
				}(&obj.StringFieldWithValidation, safe.Field(oldObj, func(oldObj *T2) *string { return &oldObj.StringFieldWithValidation }), fldPath.Child("stringFieldWithValidation"))...)

			// field T2.StructField
			errs = append(errs,
				func(obj, oldObj *StructField, fldPath *field.Path) (errs field.ErrorList) {
					errs = append(errs, validate.FixedResult(opCtx, fldPath, obj, oldObj, false, "subfield T1.PT2.StructField")...)
					return
				}(&obj.StructField, safe.Field(oldObj, func(oldObj *T2) *StructField { return &oldObj.StructField }), fldPath.Child("structField"))...)

			errs = append(errs, Validate_T2(opCtx, obj, oldObj, fldPath)...)
			return
		}(obj.PT2, safe.Field(oldObj, func(oldObj *T1) *T2 { return oldObj.PT2 }), fldPath.Child("pt2"))...)

	return errs
}

func Validate_T2(opCtx operation.Context, obj, oldObj *T2, fldPath *field.Path) (errs field.ErrorList) {
	// field T2.MapField has no validation
	// field T2.PointerField has no validation
	// field T2.SliceField has no validation
	// field T2.StringField has no validation

	// field T2.StringFieldWithValidation
	errs = append(errs,
		func(obj, oldObj *string, fldPath *field.Path) (errs field.ErrorList) {
			errs = append(errs, validate.FixedResult(opCtx, fldPath, obj, oldObj, false, "field T2.StringFieldWithValidation")...)
			return
		}(&obj.StringFieldWithValidation, safe.Field(oldObj, func(oldObj *T2) *string { return &oldObj.StringFieldWithValidation }), fldPath.Child("stringFieldWithValidation"))...)

	// field T2.StructField has no validation
	return errs
}
