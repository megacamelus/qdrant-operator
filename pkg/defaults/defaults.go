package defaults

import "time"

const (
	SyncInterval       = 5 * time.Second
	RetryInterval      = 10 * time.Second
	ConflictInterval   = 1 * time.Second
	KaotoFinalizerName = "qdrant.megacamelus.github.io/finalizer"
)

var (
	QdrantImage = "qdrant/qdrant:v1.8.3-unprivileged"
)
