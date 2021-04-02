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

package v1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type Rule struct {
	// TODO: make these fields required with kubebuilder
	ParentRole string `json:"parent_role"`
	ChildRole  string `json:"child_role"`
}

// RoleImplicationRuleSpec defines the desired state of RoleImplicationRule
type RoleImplicationRuleSpec struct {
	ImplicationRule Rule `json:"implication_rule"`
}

type GeneratedImplications struct {
	ParentRole string   `json:"parent_role"`
	ChildRoles []string `json:"child_roles"`
}

// RoleImplicationRuleStatus defines the observed state of RoleImplicationRule
type RoleImplicationRuleStatus struct {
	RoleImplications GeneratedImplications `json:"role_implications"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// RoleImplicationRule is the Schema for the roleimplicationrules API
type RoleImplicationRule struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   RoleImplicationRuleSpec   `json:"spec,omitempty"`
	Status RoleImplicationRuleStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// RoleImplicationRuleList contains a list of RoleImplicationRule
type RoleImplicationRuleList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []RoleImplicationRule `json:"items"`
}

func init() {
	SchemeBuilder.Register(&RoleImplicationRule{}, &RoleImplicationRuleList{})
}
