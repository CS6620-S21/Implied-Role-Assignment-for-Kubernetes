package auth

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httputil"
	"reflect"
	"strings"
	"testing"
	"time"

	rbacapi "k8s.io/api/rbac/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/apiserver/pkg/authentication/request/bearertoken"
	"k8s.io/apiserver/pkg/authentication/token/tokenfile"
	"k8s.io/apiserver/pkg/authentication/user"
	"k8s.io/apiserver/pkg/authorization/authorizer"

	"k8s.io/apiserver/pkg/registry/generic"
	clientset "k8s.io/client-go/kubernetes"
	"k8s.io/client-go/transport"
	"k8s.io/klog/v2"
	"k8s.io/kubernetes/pkg/api/legacyscheme"

	"k8s.io/kubernetes/pkg/controlplane"
	"k8s.io/kubernetes/pkg/registry/rbac/clusterrole"

	"k8s.io/kubernetes/pkg/registry/rbac/clusterrolebinding"

	"k8s.io/kubernetes/pkg/registry/rbac/role"

	"k8s.io/kubernetes/pkg/registry/rbac/rolebinding"
	"k8s.io/kubernetes/plugin/pkg/auth/authorizer/rbac"
	"k8s.io/kubernetes/test/integration/framework"
)



// bootstrapRoles are a set of RBAC roles which will be populated before the test.
type bootstrapRoles struct {
	roles               []rbacapi.Role
	roleBindings        []rbacapi.RoleBinding
	clusterRoles        []rbacapi.ClusterRole
	clusterRoleBindings []rbacapi.ClusterRoleBinding
}


// bootstrap uses the provided client to create the bootstrap roles and role bindings.
// client should be authenticated as the RBAC super user.
func (b bootstrapRoles) bootstrap(client clientset.Interface) error {

	for _, r := range b.roles {
		_, err := client.RbacV1().Roles(r.Namespace).Create(context.TODO(), &r, metav1.CreateOptions{})
		if err != nil {
			return fmt.Errorf("failed to make request: %v", err)
		}
	}

	for _, r := range b.roleBindings {
		_, err := client.RbacV1().RoleBindings(r.Namespace).Create(context.TODO(), &r, metav1.CreateOptions{})
		if err != nil {
			return fmt.Errorf("failed to make request: %v", err)
		}
	}

	return nil
}



tests := []struct {
		bootstrapRoles bootstrapRoles
	} { {    
    bootstrapRoles: bootstrapRoles{
				roleBindings: []rbacapi.RoleBinding{
					{
						ObjectMeta: metav1.ObjectMeta{Name: "write-jobs", Namespace: "job-namespace"},
						Subjects:   []rbacapi.Subject{{Kind: "User", Name: "job-writer-namespace"}},
						RoleRef:    rbacapi.RoleRef{Kind: "ClusterRole", Name: "write-jobs"},
					},
					{
						ObjectMeta: metav1.ObjectMeta{Name: "create-rolebindings", Namespace: "job-namespace"},
						Subjects: []rbacapi.Subject{
							{Kind: "User", Name: "job-writer-namespace"},
							{Kind: "User", Name: "any-rolebinding-writer-namespace"},
						},
						RoleRef: rbacapi.RoleRef{Kind: "ClusterRole", Name: "create-rolebindings"},
					},
					{
						ObjectMeta: metav1.ObjectMeta{Name: "bind-any-clusterrole", Namespace: "job-namespace"},
						Subjects:   []rbacapi.Subject{{Kind: "User", Name: "any-rolebinding-writer-namespace"}},
						RoleRef:    rbacapi.RoleRef{Kind: "ClusterRole", Name: "bind-any-clusterrole"},
					},
				},
			}
	    }
	}