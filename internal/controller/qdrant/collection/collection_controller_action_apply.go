package collection

import (
	"context"
	"fmt"
	"strings"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/lburgazzoli/qdrant-operator/pkg/controller/client"
	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/builder"

	pb "github.com/qdrant/go-client/qdrant"
)

func NewApplyAction() Action {
	return &applyAction{}
}

type applyAction struct {
}

func (a *applyAction) Configure(_ context.Context, _ *client.Client, b *builder.Builder) (*builder.Builder, error) {
	//
	// TODO: there should be a timer event here to reconcile the cluster status with the
	//       local status, but the current implementation is just enough for a POC which
	//       assumes that collections are only managed through this operator
	//

	return b, nil
}

func (a *applyAction) Cleanup(ctx context.Context, rr *ReconciliationRequest) error {
	return WithCollectionsClient(ctx, rr, a.delete)
}

func (a *applyAction) Apply(ctx context.Context, rr *ReconciliationRequest) error {

	applyCondition := metav1.Condition{
		Type:               "Apply",
		Status:             metav1.ConditionTrue,
		Reason:             "Applied",
		Message:            "Applied",
		ObservedGeneration: rr.Collection.Generation,
	}

	// TODO: create --> upsert
	if err := WithCollectionsClient(ctx, rr, a.create); err != nil {
		applyCondition.Status = metav1.ConditionFalse
		applyCondition.Reason = "Failure"
		applyCondition.Message = err.Error()
	}

	meta.SetStatusCondition(&rr.Collection.Status.Conditions, applyCondition)

	return nil
}

func (a *applyAction) create(ctx context.Context, rr *ReconciliationRequest, cc pb.CollectionsClient) error {
	var distance pb.Distance

	switch rr.Collection.Spec.VectorParams.Distance {
	case "Cosine":
		distance = pb.Distance_Cosine
	case "Manhattan":
		distance = pb.Distance_Manhattan
	case "Euclid":
		distance = pb.Distance_Euclid
	case "Dot":
		distance = pb.Distance_Dot
	default:
		return fmt.Errorf("unsupported vector distance type: %s", rr.Collection.Spec.VectorParams.Distance)
	}

	name := rr.Collection.Spec.Name
	if name == "" {
		name = rr.Collection.Name
	}

	_, err := cc.Create(ctx, &pb.CreateCollection{
		CollectionName: name,
		VectorsConfig: &pb.VectorsConfig{
			Config: &pb.VectorsConfig_Params{
				Params: &pb.VectorParams{
					Size:     rr.Collection.Spec.VectorParams.Size,
					Distance: distance,
				},
			},
		},
	})

	if err != nil {
		ge, ok := status.FromError(err)
		if !ok {
			return fmt.Errorf("could not create collections %s, %w", rr.Collection.Name, err)
		}

		switch {
		case ge.Code() == codes.AlreadyExists:
			return nil
		case ge.Code() == codes.InvalidArgument && strings.Contains(ge.Message(), "already exists"):
			// for some reason qdrant returns InvalidArgument also to signal that a collection
			// already exists
			return nil
		default:
			return fmt.Errorf("could not create collections %s, %w", rr.Collection.Name, err)
		}

	}

	return nil
}

func (a *applyAction) delete(ctx context.Context, rr *ReconciliationRequest, cc pb.CollectionsClient) error {
	_, err := cc.Delete(ctx, &pb.DeleteCollection{
		CollectionName: rr.Collection.Name,
	})

	if err != nil {
		ge, ok := status.FromError(err)
		if !ok {
			return fmt.Errorf("could not delete collections %s, %w", rr.Collection.Name, err)
		}

		if ge.Code() == codes.InvalidArgument {
			// for some reason qdrant returns InvalidArgument also to signal that a collection
			// des not exist
			return nil
		} else {
			return fmt.Errorf("could not delete collections %s, %w", rr.Collection.Name, err)
		}

	}

	return nil
}
