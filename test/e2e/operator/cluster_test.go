package operator

import (
	"encoding/json"
	"github.com/lburgazzoli/gomega-matchers/pkg/matchers/jq"
	"testing"

	"github.com/rs/xid"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"

	. "github.com/megacamelus/qdrant-operator/test/support"
	. "github.com/onsi/gomega"

	qdrantAc "github.com/megacamelus/qdrant-operator/pkg/client/applyconfiguration/qdrant/v1alpha1"

	pb "github.com/qdrant/go-client/qdrant"
)

func TestClusterDeploy(t *testing.T) {
	test := With(t)

	ns := test.NewTestNamespace()
	name := xid.New().String()

	cluster := DeployCluster(
		test,
		qdrantAc.Cluster(name, ns.Name),
	)

	test.Eventually(ConfigMap(test, cluster.Name, cluster.Namespace), TestTimeoutLong).Should(
		Not(BeNil()))
	test.Eventually(Deployment(test, cluster.Name, cluster.Namespace), TestTimeoutLong).Should(
		WithTransform(ConditionStatus(appsv1.DeploymentAvailable), Equal(corev1.ConditionTrue)))

	test.Eventually(Cluster(test, cluster), TestTimeoutLong).Should(
		WithTransform(json.Marshal, And(
			jq.Match(`.status.phase == "Ready"`),
		)),
	)
}

func TestCollectionDeploy(t *testing.T) {
	test := With(t)

	ns := test.NewTestNamespace()
	name := xid.New().String()

	cluster := DeployCluster(
		test,
		qdrantAc.Cluster(name, ns.Name),
	)

	test.Eventually(ConfigMap(test, cluster.Name, cluster.Namespace), TestTimeoutLong).Should(
		Not(BeNil()))
	test.Eventually(Deployment(test, cluster.Name, cluster.Namespace), TestTimeoutLong).Should(
		WithTransform(ConditionStatus(appsv1.DeploymentAvailable), Equal(corev1.ConditionTrue)))

	test.Eventually(Cluster(test, cluster), TestTimeoutLong).Should(
		WithTransform(json.Marshal, And(
			jq.Match(`.status.phase == "Ready"`),
		)),
	)

	collection := DeployCollection(
		test,
		qdrantAc.Collection(name, ns.Name).
			WithSpec(
				qdrantAc.CollectionSpec().
					WithCluster(cluster.Name).
					WithVectorParams(
						qdrantAc.VectorParams().
							WithDistance(pb.Distance_Cosine.String()).
							WithSize(384))),
	)

	test.Eventually(Collection(test, collection), TestTimeoutLong).Should(
		WithTransform(json.Marshal, And(
			jq.Match(`.status.phase == "Ready"`),
			jq.Match(`.status.collection.status == "Green"`),
		)),
	)
}
