/*
Copyright 2014 The Kubernetes Authors.

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
	"k8s.io/apimachinery/pkg/util/validation/field"
)

// FixedResult asserts a fixed boolean result.  This is mostly useful for
// testing.
func FixedResult[T any](fldPath *field.Path, value T, result bool, arg string) field.ErrorList {
	if result {
		return nil
	}
	return field.ErrorList{
		field.Invalid(fldPath, value, "forced failure: "+arg),
	}
}

// FixedResultUpdate asserts a fixed boolean result.  This is mostly useful for
// testing updates.
func FixedResultUpdate[T any](fldPath *field.Path, value, oldValue T, result bool, arg string) field.ErrorList {
	if result {
		return nil
	}
	return field.ErrorList{
		field.Invalid(fldPath, value, "forced failure: "+arg),
	}
}
