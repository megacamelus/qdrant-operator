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

package applyconfiguration

import (
	v1alpha1 "github.com/megacamelus/qdrant-operator/api/qdrant/v1alpha1"
	qdrantv1alpha1 "github.com/megacamelus/qdrant-operator/pkg/client/qdrant/applyconfiguration/qdrant/v1alpha1"
	schema "k8s.io/apimachinery/pkg/runtime/schema"
)

// ForKind returns an apply configuration type for the given GroupVersionKind, or nil if no
// apply configuration type exists for the given GroupVersionKind.
func ForKind(kind schema.GroupVersionKind) interface{} {
	switch kind {
	// Group=qdrant, Version=v1alpha1
	case v1alpha1.SchemeGroupVersion.WithKind("Cluster"):
		return &qdrantv1alpha1.ClusterApplyConfiguration{}
	case v1alpha1.SchemeGroupVersion.WithKind("ClusterSpec"):
		return &qdrantv1alpha1.ClusterSpecApplyConfiguration{}
	case v1alpha1.SchemeGroupVersion.WithKind("ClusterStatus"):
		return &qdrantv1alpha1.ClusterStatusApplyConfiguration{}
	case v1alpha1.SchemeGroupVersion.WithKind("Collection"):
		return &qdrantv1alpha1.CollectionApplyConfiguration{}
	case v1alpha1.SchemeGroupVersion.WithKind("CollectionInfo"):
		return &qdrantv1alpha1.CollectionInfoApplyConfiguration{}
	case v1alpha1.SchemeGroupVersion.WithKind("CollectionSpec"):
		return &qdrantv1alpha1.CollectionSpecApplyConfiguration{}
	case v1alpha1.SchemeGroupVersion.WithKind("CollectionStatus"):
		return &qdrantv1alpha1.CollectionStatusApplyConfiguration{}
	case v1alpha1.SchemeGroupVersion.WithKind("VectorParams"):
		return &qdrantv1alpha1.VectorParamsApplyConfiguration{}

	}
	return nil
}
