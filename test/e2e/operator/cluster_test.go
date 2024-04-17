package operator

import (
	"testing"

	"github.com/rs/xid"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"

	. "github.com/megacamelus/qdrant-operator/test/support"
	. "github.com/megacamelus/qdrant-operator/test/support/matchers"
	. "github.com/onsi/gomega"

	qdrantAc "github.com/megacamelus/qdrant-operator/pkg/client/qdrant/applyconfiguration/qdrant/v1alpha1"
)

func TestClusterDeploy(t *testing.T) {
	test := With(t)

	ns := test.NewTestNamespace()
	name := xid.New().String()

	cluster := DeployCluster(
		test,
		qdrantAc.Cluster(name, ns.Name),
	)

	test.Eventually(Deployment(test, cluster.Name, cluster.Namespace), TestTimeoutLong).Should(
		WithTransform(ConditionStatus(appsv1.DeploymentAvailable), Equal(corev1.ConditionTrue)))

	test.Eventually(Cluster(test, cluster), TestTimeoutLong).Should(
		WithTransform(AsJSON(), And(
			MatchJQ(`.status.phase == "Ready"`),
		)),
	)
}
