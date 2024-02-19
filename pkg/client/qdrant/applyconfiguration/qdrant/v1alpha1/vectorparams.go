/*
Copyright 2023.

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
// Code generated by applyconfiguration-gen. DO NOT EDIT.

package v1alpha1

// VectorParamsApplyConfiguration represents an declarative configuration of the VectorParams type for use
// with apply.
type VectorParamsApplyConfiguration struct {
	Size     *uint64 `json:"size,omitempty"`
	Distance *string `json:"distance,omitempty"`
}

// VectorParamsApplyConfiguration constructs an declarative configuration of the VectorParams type for use with
// apply.
func VectorParams() *VectorParamsApplyConfiguration {
	return &VectorParamsApplyConfiguration{}
}

// WithSize sets the Size field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the Size field is set to the value of the last call.
func (b *VectorParamsApplyConfiguration) WithSize(value uint64) *VectorParamsApplyConfiguration {
	b.Size = &value
	return b
}

// WithDistance sets the Distance field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the Distance field is set to the value of the last call.
func (b *VectorParamsApplyConfiguration) WithDistance(value string) *VectorParamsApplyConfiguration {
	b.Distance = &value
	return b
}