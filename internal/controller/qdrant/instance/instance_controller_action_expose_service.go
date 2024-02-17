package instance

import (
	"context"

	"github.com/lburgazzoli/qdrant-operator/internal/controller/qdrant"
	"github.com/lburgazzoli/qdrant-operator/pkg/apply"
	"github.com/lburgazzoli/qdrant-operator/pkg/controller/client"

	"k8s.io/apimachinery/pkg/api/meta"
	"k8s.io/apimachinery/pkg/util/intstr"

	"sigs.k8s.io/controller-runtime/pkg/builder"
	"sigs.k8s.io/controller-runtime/pkg/predicate"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	corev1ac "k8s.io/client-go/applyconfigurations/core/v1"
)

func NewServiceAction() Action {
	return &serviceAction{}
}

type serviceAction struct {
}

func (a *serviceAction) Configure(_ context.Context, _ *client.Client, b *builder.Builder) (*builder.Builder, error) {
	b = b.Owns(&corev1.Service{}, builder.WithPredicates(
		predicate.Or(
			predicate.ResourceVersionChangedPredicate{},
		)))

	return b, nil
}

func (a *serviceAction) Cleanup(context.Context, *ReconciliationRequest) error {
	return nil
}

func (a *serviceAction) Apply(ctx context.Context, rr *ReconciliationRequest) error {
	serviceCondition := metav1.Condition{
		Type:               "Service",
		Status:             metav1.ConditionTrue,
		Reason:             "Deployed",
		Message:            "Deployed",
		ObservedGeneration: rr.Instance.Generation,
	}

	err := a.service(ctx, rr)
	if err != nil {
		serviceCondition.Status = metav1.ConditionFalse
		serviceCondition.Reason = "Failure"
		serviceCondition.Message = err.Error()

		return err
	}

	meta.SetStatusCondition(&rr.Instance.Status.Conditions, serviceCondition)

	return nil
}

func (a *serviceAction) service(ctx context.Context, rr *ReconciliationRequest) error {
	service := corev1ac.Service(rr.Instance.Name, rr.Instance.Namespace).
		WithOwnerReferences(apply.WithOwnerReference(rr.Instance)).
		WithLabels(Labels(rr.Instance)).
		WithSpec(corev1ac.ServiceSpec().
			WithPorts(
				corev1ac.ServicePort().
					WithName(qdrant.QdrantHttpPortType).
					WithProtocol(corev1.ProtocolTCP).
					WithPort(qdrant.QdrantHttpPort).
					WithTargetPort(intstr.FromString(qdrant.QdrantHttpPortType)),
				corev1ac.ServicePort().
					WithName(qdrant.QdrantGrpcPortType).
					WithProtocol(corev1.ProtocolTCP).
					WithPort(qdrant.QdrantGrpcPort).
					WithTargetPort(intstr.FromString(qdrant.QdrantGrpcPortType))).
			WithSelector(LabelsForSelector(rr.Instance)).
			WithSessionAffinity(corev1.ServiceAffinityNone).
			WithPublishNotReadyAddresses(true))

	_, err := rr.Client.CoreV1().Services(rr.Instance.Namespace).Apply(
		ctx,
		service,
		metav1.ApplyOptions{
			FieldManager: qdrant.QdrantOperatorFieldManager,
			Force:        true,
		},
	)

	return err
}
