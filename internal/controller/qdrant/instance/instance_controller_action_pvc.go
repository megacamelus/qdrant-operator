package instance

import (
	"context"

	"github.com/lburgazzoli/qdrant-operator/internal/controller/qdrant"
	"github.com/lburgazzoli/qdrant-operator/pkg/apply"
	"github.com/lburgazzoli/qdrant-operator/pkg/controller/client"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	corev1ac "k8s.io/client-go/applyconfigurations/core/v1"

	"k8s.io/apimachinery/pkg/api/meta"

	"sigs.k8s.io/controller-runtime/pkg/builder"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
)

func NewPersistentVolumeClaimAction() Action {
	return &pvcAction{}
}

type pvcAction struct {
}

func (a *pvcAction) Configure(_ context.Context, _ *client.Client, b *builder.Builder) (*builder.Builder, error) {
	b = b.Owns(&appsv1.Deployment{}, builder.WithPredicates(
		predicate.Or(
			predicate.ResourceVersionChangedPredicate{},
		)))

	return b, nil
}

func (a *pvcAction) Cleanup(context.Context, *ReconciliationRequest) error {
	return nil
}

func (a *pvcAction) Apply(ctx context.Context, rr *ReconciliationRequest) error {
	deploymentCondition := metav1.Condition{
		Type:               "PersistentVolumeClaim",
		Status:             metav1.ConditionTrue,
		Reason:             "Deployed",
		Message:            "Deployed",
		ObservedGeneration: rr.Instance.Generation,
	}

	err := a.deploy(ctx, rr)
	if err != nil {
		deploymentCondition.Status = metav1.ConditionFalse
		deploymentCondition.Reason = "Failure"
		deploymentCondition.Message = err.Error()
	}

	meta.SetStatusCondition(&rr.Instance.Status.Conditions, deploymentCondition)

	return err
}

func (a *pvcAction) deploy(ctx context.Context, rr *ReconciliationRequest) error {
	pvc := corev1ac.PersistentVolumeClaim(rr.Instance.Name, rr.Instance.Namespace).
		WithOwnerReferences(apply.WithOwnerReference(rr.Instance)).
		WithLabels(Labels(rr.Instance)).
		WithSpec(corev1ac.PersistentVolumeClaimSpec().
			WithAccessModes(corev1.ReadWriteOnce).
			WithResources(corev1ac.VolumeResourceRequirements().
				WithRequests(corev1.ResourceList{
					"storage": QdrantInstanceStorage,
				})))

	_, err := rr.Client.CoreV1().PersistentVolumeClaims(rr.Instance.Namespace).Apply(
		ctx,
		pvc,
		metav1.ApplyOptions{
			FieldManager: qdrant.QdrantOperatorFieldManager,
			Force:        true,
		},
	)

	return err
}
