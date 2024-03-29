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
	"context"
	"fmt"

	qdrantv1alpha1 "github.com/megacamelus/qdrant-operator/api/qdrant/v1alpha1"

	"k8s.io/apimachinery/pkg/types"

	"sigs.k8s.io/controller-runtime/pkg/builder"

	"github.com/megacamelus/qdrant-operator/pkg/controller/client"
)

const (
	AppName string = "collection"
)

type ReconciliationRequest struct {
	*client.Client

	Collection *qdrantv1alpha1.Collection
}

func (rr *ReconciliationRequest) Key() types.NamespacedName {
	return types.NamespacedName{
		Namespace: rr.Collection.Namespace,
		Name:      rr.Collection.Name,
	}
}

func (rr *ReconciliationRequest) String() string {
	return fmt.Sprintf("%s/%s", rr.Collection.Namespace, rr.Collection.Name)
}

type Action interface {
	Configure(context.Context, *client.Client, *builder.Builder) (*builder.Builder, error)
	Apply(context.Context, *ReconciliationRequest) error
	Cleanup(context.Context, *ReconciliationRequest) error
}
