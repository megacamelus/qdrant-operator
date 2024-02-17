package instance

import (
	"context"

	"github.com/lburgazzoli/qdrant-operator/internal/controller/qdrant"
	"github.com/lburgazzoli/qdrant-operator/pkg/apply"
	"github.com/lburgazzoli/qdrant-operator/pkg/controller/client"
	"github.com/lburgazzoli/qdrant-operator/pkg/defaults"

	"k8s.io/apimachinery/pkg/api/meta"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	appsv1ac "k8s.io/client-go/applyconfigurations/apps/v1"
	corev1ac "k8s.io/client-go/applyconfigurations/core/v1"
	metav1ac "k8s.io/client-go/applyconfigurations/meta/v1"

	"sigs.k8s.io/controller-runtime/pkg/builder"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
)

func NewDeployAction() Action {
	return &deployAction{}
}

type deployAction struct {
}

func (a *deployAction) Configure(_ context.Context, _ *client.Client, b *builder.Builder) (*builder.Builder, error) {
	b = b.Owns(&appsv1.Deployment{}, builder.WithPredicates(
		predicate.Or(
			predicate.ResourceVersionChangedPredicate{},
		)))

	return b, nil
}

func (a *deployAction) Cleanup(context.Context, *ReconciliationRequest) error {
	return nil
}

func (a *deployAction) Apply(ctx context.Context, rr *ReconciliationRequest) error {
	deploymentCondition := metav1.Condition{
		Type:               "Deployment",
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

func (a *deployAction) deploy(ctx context.Context, rr *ReconciliationRequest) error {

	d := a.deployment(rr)

	_, err := rr.Client.AppsV1().Deployments(rr.Instance.Namespace).Apply(
		ctx,
		d,
		metav1.ApplyOptions{
			FieldManager: qdrant.QdrantOperatorFieldManager,
			Force:        true,
		},
	)

	return err
}

func (a *deployAction) deployment(rr *ReconciliationRequest) *appsv1ac.DeploymentApplyConfiguration {
	image := rr.Instance.Spec.Image
	if image == "" {
		image = defaults.QdrantImage
	}

	labels := Labels(rr.Instance)

	envs := make([]*corev1ac.EnvVarApplyConfiguration, 0)
	envs = append(envs, apply.WithEnvFromField("NAMESPACE", "metadata.namespace"))

	return appsv1ac.Deployment(rr.Instance.Name, rr.Instance.Namespace).
		WithOwnerReferences(apply.WithOwnerReference(rr.Instance)).
		WithLabels(Labels(rr.Instance)).
		WithSpec(appsv1ac.DeploymentSpec().
			WithReplicas(1).
			WithSelector(metav1ac.LabelSelector().WithMatchLabels(labels)).
			WithTemplate(corev1ac.PodTemplateSpec().
				WithLabels(labels).
				WithSpec(corev1ac.PodSpec().
					WithSecurityContext(corev1ac.PodSecurityContext().
						WithRunAsNonRoot(true).
						WithSeccompProfile(corev1ac.SeccompProfile().WithType(corev1.SeccompProfileTypeRuntimeDefault))).
					WithContainers(corev1ac.Container().
						WithImage(image).
						WithImagePullPolicy(corev1.PullAlways).
						WithName(qdrant.QdrantAppName).
						WithPorts(
							apply.WithPort(qdrant.QdrantHttpPortType, qdrant.QdrantHttpPort),
							apply.WithPort(qdrant.QdrantGrpcPortType, qdrant.QdrantGrpcPort)).
						WithReadinessProbe(apply.WithHTTPProbe(qdrant.QdrantReadinessProbePath, qdrant.QdrantHttpPort)).
						WithLivenessProbe(apply.WithHTTPProbe(qdrant.QdrantLivenessProbePath, qdrant.QdrantHttpPort)).
						WithEnv(envs...).
						WithResources(corev1ac.ResourceRequirements().WithRequests(corev1.ResourceList{
							corev1.ResourceMemory: QdrantInstanceDefaultMemory,
							corev1.ResourceCPU:    QdrantInstanceDefaultCPU,
						})).
						WithSecurityContext(corev1ac.SecurityContext().
							WithAllowPrivilegeEscalation(false).
							WithRunAsNonRoot(true).
							WithSeccompProfile(corev1ac.SeccompProfile().WithType(corev1.SeccompProfileTypeRuntimeDefault)))))))
}
