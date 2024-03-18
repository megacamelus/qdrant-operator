package cluster

import (
	"context"
	"fmt"

	"github.com/megacamelus/qdrant-operator/internal/controller/qdrant"
	appsv1 "k8s.io/api/apps/v1"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/megacamelus/qdrant-operator/pkg/controller/client"

	"sigs.k8s.io/controller-runtime/pkg/builder"
)

func NewStatusAction() Action {
	return &statusAction{}
}

type statusAction struct {
}

func (a *statusAction) Configure(_ context.Context, _ *client.Client, b *builder.Builder) (*builder.Builder, error) {
	return b, nil
}

func (a *statusAction) Cleanup(_ context.Context, _ *ReconciliationRequest) error {
	return nil
}

func (a *statusAction) Apply(ctx context.Context, rr *ReconciliationRequest) error {
	d, err := rr.Client.AppsV1().Deployments(rr.Cluster.Namespace).Get(ctx, rr.Cluster.Name, metav1.GetOptions{})

	if k8serrors.IsNotFound(err) {
		return nil
	}
	if err != nil {
		return err
	}

	s, err := rr.Client.CoreV1().Services(rr.Cluster.Namespace).Get(ctx, rr.Cluster.Name, metav1.GetOptions{})

	if k8serrors.IsNotFound(err) {
		return nil
	}
	if err != nil {
		return err
	}

	var available *appsv1.DeploymentCondition

	for i := range d.Status.Conditions {
		if d.Status.Conditions[i].Type == appsv1.DeploymentAvailable {
			available = &d.Status.Conditions[i]
			break
		}
	}

	if available == nil {
		return nil
	}

	for i := range s.Spec.Ports {
		switch s.Spec.Ports[i].Name {
		case qdrant.QdrantHTTPPortType:
			rr.Cluster.Status.HTTPEndpoint = fmt.Sprintf("http://%s.%s.svc.cluster.local:%d", s.Name, s.Namespace, s.Spec.Ports[i].Port)
		case qdrant.QdrantGrpcPortType:
			rr.Cluster.Status.GrpcEndpoint = fmt.Sprintf("%s.%s.svc.cluster.local:%d", s.Name, s.Namespace, s.Spec.Ports[i].Port)
		}
	}

	return nil
}
