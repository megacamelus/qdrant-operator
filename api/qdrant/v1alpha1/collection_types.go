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

package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// CollectionSpec defines the desired state of Collection.
type CollectionSpec struct {
	// +kubebuilder:validation:Required
	Cluster string `json:"cluster"`

	// +kubebuilder:validation:Optional
	Name string `json:"name,omitempty"`

	// +kubebuilder:validation:Required
	VectorParams *VectorParams `json:"vectorParams,omitempty"`
}

type VectorParams struct {
	// +kubebuilder:validation:Required
	Size uint64 `json:"size"`

	// +kubebuilder:validation:Required
	// +kubebuilder:validation:Enum:=Cosine;Euclid;Dot;Manhattan
	Distance string `json:"distance"`
}

// CollectionStatus defines the observed state of Collection.
type CollectionStatus struct {
	Phase              string             `json:"phase"`
	Conditions         []metav1.Condition `json:"conditions,omitempty"`
	ObservedGeneration int64              `json:"observedGeneration,omitempty"`

	Status       string `json:"status"`
	VectorsCount uint64 `json:"vectorsCount"`
	PointsCount  uint64 `json:"pointsCount"`
}

// +genclient
// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:storageversion
// +kubebuilder:printcolumn:name="Cluster",type=string,JSONPath=`.spec.cluster`,description="The Cluster"

// Collection is the Schema for the collections API.
type Collection struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   CollectionSpec   `json:"spec,omitempty"`
	Status CollectionStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// CollectionList contains a list of Collection.
type CollectionList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Collection `json:"items"`
}
