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

import (
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/selection"
)

const (
	QdrantAppName              string = "qdrant"
	QdrantOperatorFieldManager string = "qdrant-operator"
	QdrantHTTPPort             int32  = 6333
	QdrantHTTPPortType         string = "http"
	QdrantGrpcPort             int32  = 6334
	QdrantGrpcPortType         string = "grpc"
	QdrantLivenessProbePath    string = "/livez"
	QdrantReadinessProbePath   string = "/readyz"

	KubernetesLabelAppName      = "app.kubernetes.io/name"
	KubernetesLabelAppInstance  = "app.kubernetes.io/instance"
	KubernetesLabelAppComponent = "app.kubernetes.io/component"
	KubernetesLabelAppPartOf    = "app.kubernetes.io/part-of"
	KubernetesLabelAppManagedBy = "app.kubernetes.io/managed-by"
)

func AppSelector() (labels.Selector, error) {
	appName, err := labels.NewRequirement(KubernetesLabelAppPartOf, selection.Equals, []string{QdrantAppName})
	if err != nil {
		return nil, err
	}

	selector := labels.NewSelector().
		Add(*appName)

	return selector, nil
}
