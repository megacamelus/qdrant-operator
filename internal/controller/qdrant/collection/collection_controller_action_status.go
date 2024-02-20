package collection

import (
	"context"
	"fmt"

	qdrantApi "github.com/lburgazzoli/qdrant-operator/api/qdrant/v1alpha1"

	"github.com/lburgazzoli/qdrant-operator/pkg/controller/client"
	pb "github.com/qdrant/go-client/qdrant"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"sigs.k8s.io/controller-runtime/pkg/builder"
)

func NewStatusAction() Action {
	return &statusAction{}
}

type statusAction struct {
}

func (a *statusAction) Configure(_ context.Context, _ *client.Client, b *builder.Builder) (*builder.Builder, error) {
	return b, nil
}

func (a *statusAction) Cleanup(_ context.Context, _ *ReconciliationRequest) error {
	return nil
}

func (a *statusAction) Apply(ctx context.Context, rr *ReconciliationRequest) error {
	return WithCollectionsClient(ctx, rr, a.info)
}

func (a *statusAction) info(ctx context.Context, rr *ReconciliationRequest, cc pb.CollectionsClient) error {
	name := rr.Collection.Spec.Name
	if name == "" {
		name = rr.Collection.Name
	}

	r, err := cc.Get(ctx, &pb.GetCollectionInfoRequest{
		CollectionName: name,
	})

	if err != nil {
		ge, ok := status.FromError(err)
		if !ok {
			return fmt.Errorf("could not retrieve collections %s, %w", rr.Collection.Name, err)
		}

		switch {
		case ge.Code() == codes.NotFound:
			return nil
		case ge.Code() == codes.InvalidArgument:
			// for some reason qdrant returns InvalidArgument also to signal that a collection
			// already exists
			return nil
		default:
			return fmt.Errorf("could not retrieve collections %s, %w", rr.Collection.Name, err)
		}
	}

	rr.Collection.Status.CollectionInfo = &qdrantApi.CollectionInfo{
		Name: name,
	}

	if r.GetResult() != nil {
		switch r.GetResult().GetStatus() {
		case pb.CollectionStatus_Green:
			rr.Collection.Status.CollectionInfo.Status = "Green"
		case pb.CollectionStatus_Yellow:
			rr.Collection.Status.CollectionInfo.Status = "Yellow"
		case pb.CollectionStatus_Red:
			rr.Collection.Status.CollectionInfo.Status = "Red"
		default:
			rr.Collection.Status.CollectionInfo.Status = "Unknown"
		}

		rr.Collection.Status.CollectionInfo.PointsCount = r.GetResult().GetPointsCount()
		rr.Collection.Status.CollectionInfo.VectorsCount = r.GetResult().GetVectorsCount()
	}

	return nil
}
