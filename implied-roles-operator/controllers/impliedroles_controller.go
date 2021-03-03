/*
Copyright 2021.

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

	"github.com/go-logr/logr"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	rolev1alpha1 "github.com/CS6620-S21/Implied-Role-Assignment-for-Kubernetes/api/v1alpha1"
)

// ImpliedRolesReconciler reconciles a ImpliedRoles object
type ImpliedRolesReconciler struct {
	client.Client
	Log    logr.Logger
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=role.k8s.com,resources=impliedroles,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=role.k8s.com,resources=impliedroles/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=role.k8s.com,resources=impliedroles/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the ImpliedRoles object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.7.0/pkg/reconcile
func (r *ImpliedRolesReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log = r.Log.WithValues("impliedroles", req.NamespacedName)

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *ImpliedRolesReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&rolev1alpha1.ImpliedRoles{}).
		Owns(&appsv1.{}).
		Complete(r)
}
