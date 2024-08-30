/*
Copyright 2024 The Kubernetes Authors.

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

package validate

import (
	"k8s.io/apimachinery/pkg/api/operation"
	"k8s.io/apimachinery/pkg/api/validate/content"
	"k8s.io/apimachinery/pkg/util/validation/field"
)

// MaxLength verifies that the specified value is not longer than max
// characters.
func MaxLength(_ operation.Context, fldPath *field.Path, value, _ *string, max int) field.ErrorList {
	if value == nil {
		return nil
	}
	if len(*value) > max {
		return field.ErrorList{field.Invalid(fldPath, *value, content.MaxLenError(max))}
	}
	return nil
}

// Required verifies that the specified value is not the zero-value for its
// type.
func Required[T comparable](_ operation.Context, fldPath *field.Path, value, _ *T) field.ErrorList {
	if value != nil {
		var zero T
		if *value != zero {
			return nil
		}
	}
	return field.ErrorList{field.Required(fldPath, "")}
}

// Optional verifies that the specified value is not the zero-value for its
// type. This is identical to Required, but the caller should treat an error
// here as an indication that the optional value was not specified.
func Optional[T comparable](_ operation.Context, fldPath *field.Path, value, _ *T) field.ErrorList {
	if value != nil {
		var zero T
		if *value != zero {
			return nil
		}
	}
	return field.ErrorList{field.Required(fldPath, "optional value was not specified")}
}
