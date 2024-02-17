package qdrant

type ClusterType string

const (
	QdrantAppName              string = "qdrant"
	QdrantComponentInstance    string = "instance"
	QdrantOperatorFieldManager string = "qdrant-operator"
	QdrantPort                 int32  = 8080
	QdrantPortType             string = "http"
	QdrantLivenessProbePath    string = "/"
	QdrantReadinessProbePath   string = "/"

	KubernetesLabelAppName      = "app.kubernetes.io/name"
	KubernetesLabelAppInstance  = "app.kubernetes.io/instance"
	KubernetesLabelAppComponent = "app.kubernetes.io/component"
	KubernetesLabelAppPartOf    = "app.kubernetes.io/part-of"
	KubernetesLabelAppManagedBy = "app.kubernetes.io/managed-by"
)
