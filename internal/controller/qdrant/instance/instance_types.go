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

package instance

import (
	"context"
	"fmt"

	qdrantv1alpha1 "github.com/lburgazzoli/qdrant-operator/api/qdrant/v1alpha1"

	"k8s.io/apimachinery/pkg/api/resource"
	"k8s.io/apimachinery/pkg/types"

	"sigs.k8s.io/controller-runtime/pkg/builder"

	"github.com/lburgazzoli/qdrant-operator/pkg/controller/client"
)

var (
	QdrantInstanceDefaultMemory = resource.MustParse("600Mi")
	QdrantInstanceDefaultCPU    = resource.MustParse("500m")
	QdrantInstanceStorage       = resource.MustParse("1Gi")
)

type ReconciliationRequest struct {
	*client.Client

	Instance *qdrantv1alpha1.Instance
}

func (rr *ReconciliationRequest) Key() types.NamespacedName {
	return types.NamespacedName{
		Namespace: rr.Instance.Namespace,
		Name:      rr.Instance.Name,
	}
}

func (rr *ReconciliationRequest) String() string {
	return fmt.Sprintf("%s/%s", rr.Instance.Namespace, rr.Instance.Name)
}

type Action interface {
	Configure(context.Context, *client.Client, *builder.Builder) (*builder.Builder, error)
	Apply(context.Context, *ReconciliationRequest) error
	Cleanup(context.Context, *ReconciliationRequest) error
}
