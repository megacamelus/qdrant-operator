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

package collection

import (
	"strings"

	"github.com/lburgazzoli/qdrant-operator/internal/controller/qdrant"

	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/selection"

	ctrl "sigs.k8s.io/controller-runtime/pkg/client"
)

func Labels(ref ctrl.Object) map[string]string {
	return map[string]string{
		qdrant.KubernetesLabelAppName:      AppName,
		qdrant.KubernetesLabelAppInstance:  AppName + "-" + strings.ToLower(ref.GetName()),
		qdrant.KubernetesLabelAppComponent: strings.ToLower(ref.GetObjectKind().GroupVersionKind().Kind),
		qdrant.KubernetesLabelAppPartOf:    qdrant.QdrantAppName,
		qdrant.KubernetesLabelAppManagedBy: qdrant.QdrantOperatorFieldManager,
	}
}

func LabelsForSelector(ref ctrl.Object) map[string]string {
	return map[string]string{
		qdrant.KubernetesLabelAppName:     AppName,
		qdrant.KubernetesLabelAppInstance: AppName + "-" + strings.ToLower(ref.GetName()),
	}
}

func AppSelector() (labels.Selector, error) {
	appName, err := labels.NewRequirement(qdrant.KubernetesLabelAppName, selection.Equals, []string{AppName})
	if err != nil {
		return nil, err
	}

	selector := labels.NewSelector().
		Add(*appName)

	return selector, nil
}
