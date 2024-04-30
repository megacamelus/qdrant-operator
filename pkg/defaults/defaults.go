package defaults

import "time"

const (
	SyncInterval                   = 5 * time.Second
	RetryInterval                  = 10 * time.Second
	ConflictInterval               = 1 * time.Second
	ReadTimeout      time.Duration = 10 * time.Second
	WriteTimeout     time.Duration = 10 * time.Second
	ContextTimeout   time.Duration = 5 * time.Second
	RequeueDuration  time.Duration = 5 * time.Second

	ProbeInitialDelay     int32 = 5
	ProbePeriod           int32 = 10
	ProbeTimeout          int32 = 10
	ProbeFailureThreshold int32 = 10
	ProbeSuccessThreshold int32 = 1

	FinalizerName = "qdrant.megacamelus.github.io/finalizer"
)

var (
	QdrantImage = "qdrant/qdrant:v1.9.0-unprivileged"
)
