/*
Copyright 2023.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package collection

import (
	"context"
	"fmt"
	"sort"
	"time"

	"github.com/go-logr/logr"
	"github.com/megacamelus/qdrant-operator/pkg/defaults"
	"go.uber.org/multierr"

	"k8s.io/apimachinery/pkg/api/meta"
	"k8s.io/apimachinery/pkg/runtime"

	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	ctrl "sigs.k8s.io/controller-runtime"

	"sigs.k8s.io/controller-runtime/pkg/builder"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	"github.com/megacamelus/qdrant-operator/pkg/controller/client"

	qdrantApi "github.com/megacamelus/qdrant-operator/api/qdrant/v1alpha1"
)

func NewCollectionReconciler(manager ctrl.Manager) (*Reconciler, error) {
	c, err := client.NewClient(manager.GetConfig(), manager.GetScheme(), manager.GetClient())
	if err != nil {
		return nil, err
	}

	rec := Reconciler{}
	rec.l = ctrl.Log.WithName("cluster")
	rec.Client = c
	rec.Scheme = manager.GetScheme()

	rec.actions = make([]Action, 0)
	rec.actions = append(rec.actions, NewApplyAction())
	rec.actions = append(rec.actions, NewStatusAction())

	return &rec, nil
}

func Setup(ctx context.Context, mgr ctrl.Manager) error {
	rec, err := NewCollectionReconciler(mgr)
	if err != nil {
		return fmt.Errorf("unable to create Collection reconciler %w", err)
	}

	if err := rec.SetupWithManager(ctx, mgr); err != nil {
		return fmt.Errorf("unable to setup Collection reconciler %w", err)
	}

	return nil
}

// Reconciler reconciles a Cluster object.
type Reconciler struct {
	*client.Client

	Scheme  *runtime.Scheme
	actions []Action
	l       logr.Logger
}

// +kubebuilder:rbac:groups=qdrant.megacamelus.github.io,resources=collections,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=qdrant.megacamelus.github.io,resources=collections/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=qdrant.megacamelus.github.io,resources=collections/finalizers,verbs=update
// +kubebuilder:rbac:groups=qdrant.megacamelus.github.io,resources=clusters,verbs=get;list;watch
// +kubebuilder:rbac:groups=qdrant.megacamelus.github.io,resources=clusters/status,verbs=get

// SetupWithManager sets up the controller with the Manager.
func (r *Reconciler) SetupWithManager(ctx context.Context, mgr ctrl.Manager) error {
	c := ctrl.NewControllerManagedBy(mgr)

	c = c.For(&qdrantApi.Collection{}, builder.WithPredicates(
		predicate.Or(
			predicate.GenerationChangedPredicate{},
		)))

	for i := range r.actions {
		b, err := r.actions[i].Configure(ctx, r.Client, c)
		if err != nil {
			return err
		}

		c = b
	}

	return c.Complete(reconcile.AsReconciler[*qdrantApi.Collection](mgr.GetClient(), r))
}

func (r *Reconciler) Reconcile(ctx context.Context, res *qdrantApi.Collection) (ctrl.Result, error) {
	l := log.FromContext(ctx)

	rr := ReconciliationRequest{
		Client:     r.Client,
		Collection: res.DeepCopy(),
	}

	l.Info("Reconciling", "resource", rr.String())

	if rr.Collection.ObjectMeta.DeletionTimestamp.IsZero() {

		//
		// Add finalizer
		//

		if controllerutil.AddFinalizer(rr.Collection, defaults.KaotoFinalizerName) {
			if err := r.Update(ctx, rr.Collection); err != nil {
				if k8serrors.IsConflict(err) {
					return ctrl.Result{}, err
				}

				return ctrl.Result{}, fmt.Errorf("failure adding finalizer to cluster %s %w", rr.String(), err)
			}
		}
	} else {

		//
		// Cleanup leftovers if needed
		//

		for i := len(r.actions) - 1; i >= 0; i-- {
			if err := r.actions[i].Cleanup(ctx, &rr); err != nil {
				l.Info("Failure while performing cleanup, retrying in 5 sec")
				return ctrl.Result{RequeueAfter: 5 * time.Second}, err
			}
		}

		//
		// Handle finalizer
		//

		if controllerutil.RemoveFinalizer(rr.Collection, defaults.KaotoFinalizerName) {
			if err := r.Update(ctx, rr.Collection); err != nil {
				if k8serrors.IsConflict(err) {
					return ctrl.Result{}, err
				}

				return ctrl.Result{}, fmt.Errorf("failure removing finalizer to cluster %s %w", rr.String(), err)
			}
		}

		return ctrl.Result{}, nil
	}

	//
	// Reconcile
	//

	reconcileCondition := metav1.Condition{
		Type:               "Reconcile",
		Status:             metav1.ConditionTrue,
		Reason:             "Reconciled",
		Message:            "Reconciled",
		ObservedGeneration: rr.Collection.Generation,
	}

	var allErrors error

	for i := range r.actions {
		if err := r.actions[i].Apply(ctx, &rr); err != nil {
			allErrors = multierr.Append(allErrors, err)
		}
	}

	if allErrors != nil {
		reconcileCondition.Status = metav1.ConditionFalse
		reconcileCondition.Reason = "Failure"
		reconcileCondition.Message = "Failure"

		rr.Collection.Status.Phase = "Error"
	} else {
		rr.Collection.Status.ObservedGeneration = rr.Collection.Generation
		rr.Collection.Status.Phase = "Ready"
	}

	meta.SetStatusCondition(&rr.Collection.Status.Conditions, reconcileCondition)

	sort.SliceStable(rr.Collection.Status.Conditions, func(i, j int) bool {
		return rr.Collection.Status.Conditions[i].Type < rr.Collection.Status.Conditions[j].Type
	})

	//
	// Update status
	//

	err := r.Status().Update(ctx, rr.Collection)
	if err != nil && k8serrors.IsConflict(err) {
		l.Info(err.Error())
		return ctrl.Result{Requeue: true}, nil
	} else if err != nil {
		allErrors = multierr.Append(allErrors, err)
	}

	return ctrl.Result{}, allErrors
}
