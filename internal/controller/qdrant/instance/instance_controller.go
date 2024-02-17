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

package instance

import (
	"context"
	"sort"

	"github.com/go-logr/logr"
	"github.com/lburgazzoli/qdrant-operator/pkg/defaults"
	"github.com/pkg/errors"
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

	"github.com/lburgazzoli/qdrant-operator/pkg/controller/client"

	qdrantApi "github.com/lburgazzoli/qdrant-operator/api/qdrant/v1alpha1"
)

func NewInstanceReconciler(manager ctrl.Manager) (*InstanceReconciler, error) {
	c, err := client.NewClient(manager.GetConfig(), manager.GetScheme(), manager.GetClient())
	if err != nil {
		return nil, err
	}

	rec := InstanceReconciler{}
	rec.l = ctrl.Log.WithName("instance")
	rec.Client = c
	rec.Scheme = manager.GetScheme()

	rec.actions = make([]Action, 0)
	rec.actions = append(rec.actions, NewServiceAction())
	rec.actions = append(rec.actions, NewPersistentVolumeClaimAction())
	rec.actions = append(rec.actions, NewDeployAction())

	return &rec, nil
}

// InstanceReconciler reconciles a Instance object
type InstanceReconciler struct {
	*client.Client

	Scheme  *runtime.Scheme
	actions []Action
	l       logr.Logger
}

// +kubebuilder:rbac:groups=qdrant.lburgazzoli.github.io,resources=instances,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=qdrant.lburgazzoli.github.io,resources=instances/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=qdrant.lburgazzoli.github.io,resources=instances/finalizers,verbs=update
// +kubebuilder:rbac:groups="",resources=services,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups="",resources=persistentvolumeclaims,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=apps,resources=deployments,verbs=get;list;watch;create;update;patch;delete

// SetupWithManager sets up the controller with the Manager.
func (r *InstanceReconciler) SetupWithManager(ctx context.Context, mgr ctrl.Manager) error {
	c := ctrl.NewControllerManagedBy(mgr)

	c = c.For(&qdrantApi.Instance{}, builder.WithPredicates(
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

	return c.Complete(reconcile.AsReconciler[*qdrantApi.Instance](mgr.GetClient(), r))
}

func (r *InstanceReconciler) Reconcile(ctx context.Context, res *qdrantApi.Instance) (ctrl.Result, error) {
	l := log.FromContext(ctx)

	rr := ReconciliationRequest{
		Client:   r.Client,
		Instance: res.DeepCopy(),
	}

	l.Info("Reconciling", "resource", rr.String())

	if rr.Instance.ObjectMeta.DeletionTimestamp.IsZero() {

		//
		// Add finalizer
		//

		if controllerutil.AddFinalizer(rr.Instance, defaults.KaotoFinalizerName) {
			if err := r.Update(ctx, rr.Instance); err != nil {
				if k8serrors.IsConflict(err) {
					return ctrl.Result{}, err
				}

				return ctrl.Result{}, errors.Wrapf(err, "failure adding finalizer to connector cluster %s", rr.String())
			}
		}
	} else {

		//
		// Cleanup leftovers if needed
		//

		for i := len(r.actions) - 1; i >= 0; i-- {
			if err := r.actions[i].Cleanup(ctx, &rr); err != nil {
				return ctrl.Result{}, err
			}
		}

		//
		// Handle finalizer
		//

		if controllerutil.RemoveFinalizer(rr.Instance, defaults.KaotoFinalizerName) {
			if err := r.Update(ctx, rr.Instance); err != nil {
				if k8serrors.IsConflict(err) {
					return ctrl.Result{}, err
				}

				return ctrl.Result{}, errors.Wrapf(err, "failure removing finalizer from connector cluster %s", rr.String())
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
		ObservedGeneration: rr.Instance.Generation,
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

		rr.Instance.Status.Phase = "Error"
	} else {
		rr.Instance.Status.ObservedGeneration = rr.Instance.Generation
		rr.Instance.Status.Phase = "Ready"
	}

	meta.SetStatusCondition(&rr.Instance.Status.Conditions, reconcileCondition)

	sort.SliceStable(rr.Instance.Status.Conditions, func(i, j int) bool {
		return rr.Instance.Status.Conditions[i].Type < rr.Instance.Status.Conditions[j].Type
	})

	//
	// Update status
	//

	err := r.Status().Update(ctx, rr.Instance)
	if err != nil && k8serrors.IsConflict(err) {
		l.Info(err.Error())
		return ctrl.Result{Requeue: true}, nil
	} else if err != nil {
		allErrors = multierr.Append(allErrors, err)
	}

	return ctrl.Result{}, allErrors
}
