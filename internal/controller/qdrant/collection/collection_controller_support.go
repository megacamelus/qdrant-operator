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
	"errors"
	"fmt"
	"time"

	pb "github.com/qdrant/go-client/qdrant"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func endpoint(ctx context.Context, rr *ReconciliationRequest) (string, error) {
	c, err := rr.Client.Qdrant.QdrantV1alpha1().Clusters(rr.Collection.Namespace).Get(
		ctx,
		rr.Collection.Spec.Cluster,
		metav1.GetOptions{},
	)

	if k8serrors.IsNotFound(err) {
		return "", nil
	}
	if err != nil {
		return "", err
	}

	if c.Status.GrpcEndpoint == "" {
		return "", errors.New("unable to determine Grpc endpoint")
	}

	return c.Status.GrpcEndpoint, nil
}

func WithCollectionsClient(ctx context.Context, rr *ReconciliationRequest, f func(context.Context, *ReconciliationRequest, pb.CollectionsClient) error) error {
	endpoint, err := endpoint(ctx, rr)
	if err != nil {
		return fmt.Errorf("unable to determine Grpc endpoint for qdrant cluster %s, %w", rr.Collection.Spec.Cluster, err)
	}

	newCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	conn, err := grpc.DialContext(
		newCtx,
		endpoint,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)

	if err != nil {
		return fmt.Errorf("unable to connect to %s, %w", endpoint, err)
	}

	defer func() {
		_ = conn.Close()
	}()

	cc := pb.NewCollectionsClient(conn)

	return f(newCtx, rr, cc)
}
