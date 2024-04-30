package support

import (
	qdrantApi "github.com/megacamelus/qdrant-operator/api/qdrant/v1alpha1"
	qdrantAc "github.com/megacamelus/qdrant-operator/pkg/client/applyconfiguration/qdrant/v1alpha1"

	"k8s.io/apimachinery/pkg/runtime"

	"github.com/onsi/gomega"

	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func Cluster(t Test, instance *qdrantApi.Cluster) func(g gomega.Gomega) (*qdrantApi.Cluster, error) {
	return func(g gomega.Gomega) (*qdrantApi.Cluster, error) {
		answer, err := t.Client().Qdrant().QdrantV1alpha1().Clusters(instance.Namespace).Get(
			t.Ctx(),
			instance.Name,
			metav1.GetOptions{},
		)

		if k8serrors.IsNotFound(err) {
			return nil, nil
		}

		return answer, err
	}
}

func DeployCluster(
	t Test,
	in *qdrantAc.ClusterApplyConfiguration,
) *qdrantApi.Cluster {
	t.T().Helper()

	cp := t.Client().Qdrant().QdrantV1alpha1().Clusters(*in.Namespace)

	res, err := cp.Apply(
		t.Ctx(),
		in,
		metav1.ApplyOptions{
			FieldManager: "qdrant-e2e-" + t.ID(),
			Force:        true,
		})

	t.Expect(err).
		ToNot(gomega.HaveOccurred())

	t.T().Logf("Cluster %s in namespace %s created", res.Name, res.Namespace)

	t.Cleanup(func() []runtime.Object {
		return []runtime.Object{res}
	})

	return res
}

func Collection(t Test, instance *qdrantApi.Collection) func(g gomega.Gomega) (*qdrantApi.Collection, error) {
	return func(g gomega.Gomega) (*qdrantApi.Collection, error) {
		answer, err := t.Client().Qdrant().QdrantV1alpha1().Collections(instance.Namespace).Get(
			t.Ctx(),
			instance.Name,
			metav1.GetOptions{},
		)

		if k8serrors.IsNotFound(err) {
			return nil, nil
		}

		return answer, err
	}
}

func DeployCollection(
	t Test,
	in *qdrantAc.CollectionApplyConfiguration,
) *qdrantApi.Collection {
	t.T().Helper()

	cp := t.Client().Qdrant().QdrantV1alpha1().Collections(*in.Namespace)

	res, err := cp.Apply(
		t.Ctx(),
		in,
		metav1.ApplyOptions{
			FieldManager: "qdrant-e2e-" + t.ID(),
			Force:        true,
		})

	t.Expect(err).
		ToNot(gomega.HaveOccurred())

	t.T().Logf("Collectin %s in namespace %s created", res.Name, res.Namespace)

	t.Cleanup(func() []runtime.Object {
		return []runtime.Object{res}
	})

	return res
}
