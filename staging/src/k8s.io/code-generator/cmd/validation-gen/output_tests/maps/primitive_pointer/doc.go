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

// +k8s:validation-gen=TypeMeta
// +k8s:validation-gen-scheme-registry=k8s.io/code-generator/cmd/validation-gen/testscheme.Scheme

// This is a test package.
package primitivepointer

import "k8s.io/code-generator/cmd/validation-gen/testscheme"

var localSchemeBuilder = testscheme.New()

// +validateTrue="type T1"
type T1 struct {
	TypeMeta int

	// +validateTrue="field T1.MSPS"
	// +eachKey=+validateTrue="T1.MSPS[keys]"
	// +eachVal=+validateTrue="T1.MSPS[vals]"
	MSPS map[string]*string `json:"msps"`

	// Duplicate with no validation.
	AnotherMSPS map[string]*string `json:"anothermsps"`
}
