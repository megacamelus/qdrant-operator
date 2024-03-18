package cluster

import (
	"context"
	"fmt"

	"github.com/megacamelus/qdrant-operator/internal/controller/qdrant"
	"github.com/megacamelus/qdrant-operator/pkg/apply"
	"github.com/megacamelus/qdrant-operator/pkg/controller/client"
	"github.com/megacamelus/qdrant-operator/pkg/defaults"

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

func (a *deployAction) Cleanup(ctx context.Context, rr *ReconciliationRequest) error {
	collections, err := rr.Client.Qdrant.QdrantV1alpha1().Collections(rr.Cluster.Namespace).List(ctx, metav1.ListOptions{})
	if err != nil {
		return err
	}

	if len(collections.Items) == 0 {
		return nil
	}

	refs := make([]string, 0)

	for i := range collections.Items {
		if collections.Items[i].Spec.Cluster != rr.Cluster.Name {
			continue
		}

		refs = append(refs, collections.Items[i].Name)
	}

	if len(collections.Items) > 0 {
		return fmt.Errorf("cannot delete cluster with name %s, in namespace %s as it is referenced by %d collections (%s)",
			rr.Cluster.Name,
			rr.Cluster.Namespace,
			len(collections.Items),
			refs)
	}

	return nil
}

func (a *deployAction) Apply(ctx context.Context, rr *ReconciliationRequest) error {
	deploymentCondition := metav1.Condition{
		Type:               "Deployment",
		Status:             metav1.ConditionTrue,
		Reason:             "Deployed",
		Message:            "Deployed",
		ObservedGeneration: rr.Cluster.Generation,
	}

	err := a.deploy(ctx, rr)
	if err != nil {
		deploymentCondition.Status = metav1.ConditionFalse
		deploymentCondition.Reason = "Failure"
		deploymentCondition.Message = err.Error()
	}

	meta.SetStatusCondition(&rr.Cluster.Status.Conditions, deploymentCondition)

	return err
}

func (a *deployAction) deploy(ctx context.Context, rr *ReconciliationRequest) error {

	d := a.deployment(rr)

	_, err := rr.Client.AppsV1().Deployments(rr.Cluster.Namespace).Apply(
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
	image := rr.Cluster.Spec.Image
	if image == "" {
		image = defaults.QdrantImage
	}

	labels := Labels(rr.Cluster)

	envs := make([]*corev1ac.EnvVarApplyConfiguration, 0)
	envs = append(envs, apply.WithEnvFromField("NAMESPACE", "metadata.namespace"))

	return appsv1ac.Deployment(rr.Cluster.Name, rr.Cluster.Namespace).
		WithOwnerReferences(apply.WithOwnerReference(rr.Cluster)).
		WithLabels(labels).
		WithSpec(appsv1ac.DeploymentSpec().
			WithReplicas(1).
			WithSelector(metav1ac.LabelSelector().WithMatchLabels(labels)).
			WithTemplate(corev1ac.PodTemplateSpec().
				WithLabels(labels).
				WithSpec(corev1ac.PodSpec().
					WithSecurityContext(corev1ac.PodSecurityContext().
						WithRunAsNonRoot(true).
						WithSeccompProfile(corev1ac.SeccompProfile().WithType(corev1.SeccompProfileTypeRuntimeDefault))).
					WithVolumes(
						corev1ac.Volume().
							WithName(rr.Cluster.Name+"-storage").
							WithPersistentVolumeClaim(corev1ac.PersistentVolumeClaimVolumeSource().
								WithClaimName(rr.Cluster.Name)),
						corev1ac.Volume().
							WithName(rr.Cluster.Name+"-snapshot").
							WithEmptyDir(corev1ac.EmptyDirVolumeSource()),
						corev1ac.Volume().
							WithName(rr.Cluster.Name+"-init").
							WithEmptyDir(corev1ac.EmptyDirVolumeSource())).
					WithContainers(corev1ac.Container().
						WithImage(image).
						WithImagePullPolicy(corev1.PullAlways).
						WithName(qdrant.QdrantAppName).
						WithPorts(
							apply.WithPort(qdrant.QdrantHTTPPortType, qdrant.QdrantHTTPPort),
							apply.WithPort(qdrant.QdrantGrpcPortType, qdrant.QdrantGrpcPort)).
						WithReadinessProbe(apply.WithHTTPProbe(qdrant.QdrantReadinessProbePath, qdrant.QdrantHTTPPort)).
						WithLivenessProbe(apply.WithHTTPProbe(qdrant.QdrantLivenessProbePath, qdrant.QdrantHTTPPort)).
						WithEnv(envs...).
						WithResources(corev1ac.ResourceRequirements().WithRequests(corev1.ResourceList{
							corev1.ResourceMemory: QdrantClusterDefaultMemory,
							corev1.ResourceCPU:    QdrantClusterDefaultCPU,
						})).
						WithSecurityContext(corev1ac.SecurityContext().
							WithAllowPrivilegeEscalation(false).
							WithRunAsNonRoot(true).
							WithSeccompProfile(corev1ac.SeccompProfile().WithType(corev1.SeccompProfileTypeRuntimeDefault))).
						WithVolumeMounts(
							corev1ac.VolumeMount().
								WithName(rr.Cluster.Name+"-storage").
								WithMountPath("/qdrant/storage"),
							corev1ac.VolumeMount().
								WithName(rr.Cluster.Name+"-snapshot").
								WithMountPath("/qdrant/snapshot"),
							corev1ac.VolumeMount().
								WithName(rr.Cluster.Name+"-init").
								WithMountPath("/qdrant/init"))))))
}
