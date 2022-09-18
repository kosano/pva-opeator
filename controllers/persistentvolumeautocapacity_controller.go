/*
Copyright 2022.

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

package controllers

import (
	"context"

	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	"github.com/go-logr/logr"
	pvav1 "github.com/kosano/pva-operator/api/v1"
	"github.com/kosano/pva-operator/pkg/handlers"
	"github.com/kosano/pva-operator/pkg/stats"
	"github.com/kosano/pva-operator/pkg/utils"
)

// PersistentVolumeAutocapacityReconciler reconciles a PersistentVolumeAutocapacity object
type PersistentVolumeAutocapacityReconciler struct {
	client.Client
	Scheme  *runtime.Scheme
	Summary *stats.Summary
	Log     logr.Logger
}

//+kubebuilder:rbac:groups=pva.kosano.io,resources=persistentvolumeautocapacities,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=pva.kosano.io,resources=persistentvolumeautocapacities/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=pva.kosano.io,resources=persistentvolumeautocapacities/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the PersistentVolumeAutocapacity object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.11.0/pkg/reconcile
func (r *PersistentVolumeAutocapacityReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	_ = log.FromContext(ctx)
	log := r.Log.WithValues("PersistentVolumeAutocapacity", req.NamespacedName)
	// TODO(user): your logic here
	log.Info("fetching persistentvolumeautocapacity resource.")
	pva := pvav1.PersistentVolumeAutocapacity{}
	r.Get(ctx, req.NamespacedName, &pva)
	pvcHander := handlers.NewPVCHander(r.Client)
	go func() {
		for _, vs := range r.Summary.VolumeStats {
			var err error
			switch {
			case pva.Spec.Namespaces[0] == "*" && pva.Spec.PVCNames[0] == "*":
				err = pvcHander.IncreaseCapacity(ctx, vs, pva)
			case pva.Spec.Namespaces[0] != "*" && pva.Spec.PVCNames[0] == "*":
				if utils.IsInclude(vs.PVCRef.Namespace, pva.Spec.Namespaces) {
					err = pvcHander.IncreaseCapacity(ctx, vs, pva)
				}
			case pva.Spec.Namespaces[0] != "*" && pva.Spec.PVCNames[0] != "*":
				if utils.IsInclude(vs.PVCRef.Namespace, pva.Spec.Namespaces) && utils.IsInclude(vs.PVCRef.Name, pva.Spec.PVCNames) {
					err = pvcHander.IncreaseCapacity(ctx, vs, pva)
				}
			}
			if err != nil {
				log.Error(err, "pvc increase capactity failed")
				// utils.SendMail()
			}
		}
	}()

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *PersistentVolumeAutocapacityReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&pvav1.PersistentVolumeAutocapacity{}).
		Complete(r)
}
