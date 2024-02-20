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

// CollectionInfoApplyConfiguration represents an declarative configuration of the CollectionInfo type for use
// with apply.
type CollectionInfoApplyConfiguration struct {
	Name         *string `json:"name,omitempty"`
	Status       *string `json:"status,omitempty"`
	VectorsCount *uint64 `json:"vectorsCount,omitempty"`
	PointsCount  *uint64 `json:"pointsCount,omitempty"`
}

// CollectionInfoApplyConfiguration constructs an declarative configuration of the CollectionInfo type for use with
// apply.
func CollectionInfo() *CollectionInfoApplyConfiguration {
	return &CollectionInfoApplyConfiguration{}
}

// WithName sets the Name field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the Name field is set to the value of the last call.
func (b *CollectionInfoApplyConfiguration) WithName(value string) *CollectionInfoApplyConfiguration {
	b.Name = &value
	return b
}

// WithStatus sets the Status field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the Status field is set to the value of the last call.
func (b *CollectionInfoApplyConfiguration) WithStatus(value string) *CollectionInfoApplyConfiguration {
	b.Status = &value
	return b
}

// WithVectorsCount sets the VectorsCount field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the VectorsCount field is set to the value of the last call.
func (b *CollectionInfoApplyConfiguration) WithVectorsCount(value uint64) *CollectionInfoApplyConfiguration {
	b.VectorsCount = &value
	return b
}

// WithPointsCount sets the PointsCount field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the PointsCount field is set to the value of the last call.
func (b *CollectionInfoApplyConfiguration) WithPointsCount(value uint64) *CollectionInfoApplyConfiguration {
	b.PointsCount = &value
	return b
}
