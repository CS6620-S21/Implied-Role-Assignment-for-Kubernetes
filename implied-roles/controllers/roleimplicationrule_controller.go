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
	"container/list"
	"context"
	"fmt"

	"github.com/go-logr/logr"
	rbacv1 "k8s.io/api/rbac/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	rolev1 "github.com/CS6620-S21/Implied-Role-Assignment-for-Kubernetes/api/v1"
	ktypes "k8s.io/apimachinery/pkg/types"
)

// RoleImplicationRuleReconciler reconciles a RoleImplicationRule object
type RoleImplicationRuleReconciler struct {
	client.Client
	Log    logr.Logger
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=role.neu.edu,resources=roleimplicationrules,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=role.neu.edu,resources=roleimplicationrules/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=role.neu.edu,resources=roleimplicationrules/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the RoleImplicationRule object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.7.2/pkg/reconcile
func (r *RoleImplicationRuleReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	logger := r.Log

	// Get and clear all implied role bindings.
	if err := r.DeleteExistingImpliedRoleBindings(ctx); err != nil {
		return ctrl.Result{}, err
	}

	// Get role bindings.
	// TODO: Switch this to use a pure function somehow.
	roleBindings := rbacv1.RoleBindingList{}
	if err := r.GetRoleBindings(ctx, req, &roleBindings); err != nil {
		return ctrl.Result{}, err
	}

	// Get user role mappings.
	// TODO: Handle errors inside this.
	_, err := GetUserRoleMappings(roleBindings)
	if err != nil {
		return ctrl.Result{}, err
	}

	// Get all role implication rules.
	roleImplicationRules := rolev1.RoleImplicationRuleList{}
	if err := r.GetRoleImplicationRules(ctx, &roleImplicationRules); err != nil {
		return ctrl.Result{}, err
	}

	// Get the graph from implication rules.
	// TODO: Handle errors
	roleImplicationGraph, err := GetRoleImplicationGraph(roleImplicationRules)
	if err != nil {
		return ctrl.Result{}, err
	}

	allRoleImplications, err := GetAllRoleImplicationsForRoles(roleImplicationGraph)
	if err != nil {
		return ctrl.Result{}, err
	}

	// Get rid of this.
	logger.Info(fmt.Sprintf("allRoleImplications => %v", allRoleImplications))

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *RoleImplicationRuleReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&rolev1.RoleImplicationRule{}).
		Owns(&rbacv1.RoleBinding{}).
		Complete(r)
}

func (r *RoleImplicationRuleReconciler) GetRoleBindings(ctx context.Context, req ctrl.Request, roleBindings *rbacv1.RoleBindingList) error {
	logger := r.Log

	opts := []client.ListOption{
		client.InNamespace(req.NamespacedName.Namespace),
	}

	if err := r.List(ctx, roleBindings, opts...); err != nil {
		logger.Error(err, "Error fetching role bindings")
		return err
	}
	return nil
}

func (r *RoleImplicationRuleReconciler) GetRoleImplicationRules(ctx context.Context, roleImplicationRules *rolev1.RoleImplicationRuleList) error {
	logger := r.Log

	if err := r.List(ctx, roleImplicationRules); err != nil {
		logger.Error(err, "Error fetching role implication rules")
		return err
	}
	return nil
}

func GetUserRoleMappings(roleBindings rbacv1.RoleBindingList) (map[string][]string, error) {
	userRoleMappings := make(map[string][]string)

	for _, roleBinding := range roleBindings.Items {
		for _, subject := range roleBinding.Subjects {
			userRoleMappings[subject.Name] = append(userRoleMappings[subject.Name], roleBinding.RoleRef.Name)
		}
	}

	return userRoleMappings, nil
}

func GetRoleImplicationGraph(roleImplicationRules rolev1.RoleImplicationRuleList) (map[string][]string, error) {
	roleImplicationGraph := make(map[string][]string)

	for _, implicationRule := range roleImplicationRules.Items {
		roleImplicationGraph[implicationRule.Spec.ImplicationRule.ParentRole] = append(roleImplicationGraph[implicationRule.Spec.ImplicationRule.ParentRole], implicationRule.Spec.ImplicationRule.ChildRole)
	}

	return roleImplicationGraph, nil
}

func Find(slice []string, val string) (int, bool) {
	// checks if irole exists in a list of role that is assignd to the key
	for i, item := range slice {
		if item == val {
			return i, true
		}
	}
	return -1, false
}

func transform(allRoleImplications map[string][]string, x map[string][]string, role string, irole string) (map[string][]string, map[string][]string) {
	// Creates the final Map for the reconciliation.
	// Need to pass the role and its implicated role.

	// add roles if they dont exist or append irole to existing role.
	x[role] = append(x[role], irole)

	// iterate throught the temp MAP
	for role := range x {
		q := list.New()
		q.PushBack(role)
		result := list.New()
		for q.Len() != 0 {
			st := q.Front().Value
			st_temp := q.Front()
			q.Remove(st_temp)
			if st != role {
				result.PushBack(st)
			}
			stnew := fmt.Sprintf("%v", st)
			if x[stnew] != nil {
				for _, s := range x[stnew] {
					q.PushBack(s)
				}
			}
		}

		// reiterate and appende  all the implied roles from Queue.
		for e := result.Front(); e != nil; e = e.Next() {
			enew := fmt.Sprintf("%v", e.Value)
			_, found := Find(allRoleImplications[role], enew)
			if !found {
				allRoleImplications[role] = append(allRoleImplications[role], enew)
			}
		}
	}

	return allRoleImplications, x

}

func GetAllRoleImplicationsForRoles(roleImplicationGraph map[string][]string) (map[string][]string, error) {
	allRoleImplications := make(map[string][]string)

	var Implicationgraph = make(map[string][]string)

	// var allRoleImplications = make(map[string][]string)

	// var roleImplicationGraph = map[string][]string{
	// 	"admin":     {"developer", "reviewer"},
	// 	"writer":    {"pro", "noob"},
	// 	"developer": {"writer"},
	// }

	for role, irole := range roleImplicationGraph {
		for i := 0; i < len(irole); i++ {
			allRoleImplications, Implicationgraph = transform(allRoleImplications, Implicationgraph, role, irole[i])
		}
	}
	//fmt.Print(allRoleImplications)

	return allRoleImplications, nil
}

func (r *RoleImplicationRuleReconciler) CreateRoleBindings(ctx context.Context) error {
	// TODO: Fix this to actually do stuff.
	logger := r.Log

	// Check if the rolebinding already exists. Don't create it if it already exists.
	namespacedName := ktypes.NamespacedName{Namespace: "default", Name: "test-rolebinding"}
	existingRoleBinding := &rbacv1.RoleBinding{}
	if err := r.Get(ctx, namespacedName, existingRoleBinding); err == nil {
		return nil
	}

	roleRef := rbacv1.RoleRef{
		Name: "test-role-ref",
		Kind: "Role",
	}

	testUser := rbacv1.Subject{
		Kind: rbacv1.UserKind,
		Name: "test-user",
	}

	subjects := make([]rbacv1.Subject, 0)
	subjects = append(subjects, testUser)

	metadata := v1.ObjectMeta{
		Name: "test-rolebinding",
		Labels: map[string]string{
			"source": "implied-roles",
		},
		Namespace: "default",
	}

	roleBinding := rbacv1.RoleBinding{
		ObjectMeta: metadata,
		Subjects:   subjects,
		RoleRef:    roleRef,
	}

	if err := r.Create(ctx, &roleBinding); err != nil {
		logger.Error(err, "Failed test rolebinding creation")
		return err
	}

	return nil
}

func (r *RoleImplicationRuleReconciler) DeleteExistingImpliedRoleBindings(ctx context.Context) error {
	/*
		Get all role bindings create by this operator and delete them.
	*/
	logger := r.Log

	opts := []client.ListOption{
		client.MatchingLabels{"type": "implied"},
	}

	impliedRoleBindings := rbacv1.RoleBindingList{}

	if err := r.List(ctx, &impliedRoleBindings, opts...); err != nil {
		logger.Error(err, "Error fetching implied role bindings")
		return err
	}

	// Delete all the roles we've got
	for _, roleBinding := range impliedRoleBindings.Items {
		if err := r.Delete(ctx, &roleBinding); err != nil {
			logger.Error(err, fmt.Sprintf("Failed to delete RoleBinding %s", roleBinding.Name))
		}
	}

	return nil
}
