/*
Copyright 2022.

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

package qdrant

const (
	QdrantAppName              string = "qdrant"
	QdrantComponentInstance    string = "instance"
	QdrantOperatorFieldManager string = "qdrant-operator"
	QdrantHttpPort             int32  = 6333
	QdrantHttpPortType         string = "http"
	QdrantGrpcPort             int32  = 6334
	QdrantGrpcPortType         string = "grpc"
	QdrantLivenessProbePath    string = "/readyz"
	QdrantReadinessProbePath   string = "/readyz"

	KubernetesLabelAppName      = "app.kubernetes.io/name"
	KubernetesLabelAppInstance  = "app.kubernetes.io/instance"
	KubernetesLabelAppComponent = "app.kubernetes.io/component"
	KubernetesLabelAppPartOf    = "app.kubernetes.io/part-of"
	KubernetesLabelAppManagedBy = "app.kubernetes.io/managed-by"
)
