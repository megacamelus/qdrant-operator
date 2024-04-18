package cluster

import (
	"context"

	"sigs.k8s.io/yaml"

	"github.com/megacamelus/qdrant-operator/internal/controller/qdrant"
	"github.com/megacamelus/qdrant-operator/pkg/apply"
	"github.com/megacamelus/qdrant-operator/pkg/controller/client"
	"k8s.io/apimachinery/pkg/api/meta"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	corev1ac "k8s.io/client-go/applyconfigurations/core/v1"
	"sigs.k8s.io/controller-runtime/pkg/builder"
)

func NewDeployConfigAction() Action {
	return &deployConfigAction{}
}

type deployConfigAction struct {
}

func (a *deployConfigAction) Configure(_ context.Context, _ *client.Client, b *builder.Builder) (*builder.Builder, error) {
	b = b.Owns(&corev1.ConfigMap{})

	return b, nil
}

func (a *deployConfigAction) Cleanup(ctx context.Context, rr *ReconciliationRequest) error {
	return nil
}

func (a *deployConfigAction) Apply(ctx context.Context, rr *ReconciliationRequest) error {
	deploymentCondition := metav1.Condition{
		Type:               "Configuration",
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

func (a *deployConfigAction) deploy(ctx context.Context, rr *ReconciliationRequest) error {
	c, err := a.configMap(rr)
	if err != nil {
		return err
	}

	_, err = rr.Client.CoreV1().ConfigMaps(rr.Cluster.Namespace).Apply(
		ctx,
		c,
		metav1.ApplyOptions{
			FieldManager: qdrant.QdrantOperatorFieldManager,
			Force:        true,
		},
	)

	return err
}

func (a *deployConfigAction) configMap(rr *ReconciliationRequest) (*corev1ac.ConfigMapApplyConfiguration, error) {
	labels := Labels(rr.Cluster)

	config := ""
	if rr.Cluster.Spec.Configuration != nil {
		out, err := yaml.JSONToYAML(rr.Cluster.Spec.Configuration.RawMessage)
		if err != nil {
			return nil, err
		}

		config = string(out)
	}

	cm := corev1ac.ConfigMap(rr.Cluster.Name, rr.Cluster.Namespace).
		WithOwnerReferences(apply.WithOwnerReference(rr.Cluster)).
		WithLabels(labels).
		WithData(map[string]string{
			"config": config,
		})

	return cm, nil
}
