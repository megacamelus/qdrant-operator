package defaults

import "time"

const (
	SyncInterval       = 5 * time.Second
	RetryInterval      = 10 * time.Second
	ConflictInterval   = 1 * time.Second
	KaotoFinalizerName = "qdrant.lburgazzoli.github.io/finalizer"
)

var (
	QdrantImage = "qdrant/qdrant:v1.7.4-unprivileged"
)
