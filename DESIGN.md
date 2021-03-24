# Implied Role Assignment for Kubernetes

## 1. Controller Design

## _Reconciler function:_

**Trigger:**

- Role-binding is created. &#8594; Role assignment (Through role binding yaml)\
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;Kubectl apply (Only to current User)
- Role-binding is updated. &#8594; Role assignment (Through role binding yaml)\ 
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;Kubectl apply (Only to current User)
- Role un assignment using yaml, deleted role binding - Need to look how
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;kubectl delete (Only to current User)
- Create the CRD that implies the implied role.\
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;Create implied roles and create role bindings (Across all users)
- Update the CRD that implies the implied role.\
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;Update implied roles and create role bindings (Across all users)
- Delete the CRD that implies the implied role.\
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;Delete all implied roles and implied role bindings (Across all users)
- If you delete the role, it needs to be deleted from role binding and CRD 
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;NOT IN SCOPE (Look into this if kubernetes automatically handles this)

**Logic:**

- Role-binding is created &#8594;\
        1. We get the user from the context and the role assigned.\ 
        2. We check if the role is an implied role.\
            if yes,\
            &nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;-   We go and retrieve all the implied roles for this role.\
            &nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;-   We assign these roles to the user via creating role bindings. (Also making sure we maintain labels for implied role)

- Role-binding is updated. &#8594;\
        1. We get the user from the context and the role assigned. \
        2. We check if the role is an implied role.\
            if yes, \
            &nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;-   We go and retrieve all the implied roles for this role.\
            &nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;-   We assign these roles to the user via creating role bindings. (Also making sure we maintain labels for implied role)

- Role-binding is deleted. &#8594;\
        1. We get the user from the context and the role assigned. \
        2. We check if the role is an implied role.\
            if yes, \
            &nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;-   We go and retrieve all the implied roles for this role.\
            &nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;-   We un-assign these roles to the user via deleting role bindings. (Also making sure we delete only those implied by this role and maintaining labels)

- Create the CRD that implies the implied role. &#8594; \
        1. For all users.\
        2. We check if there exists an implied role  role binding for created CRD.\
            if yes, \
            &nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;-   We go and retrieve all the implied roles from this role.\
            &nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;-   We assign these roles to the user via creating role bindings. (Also making sure we maintain labels for implied role)

- Update the CRD that implies the implied role. &#8594;\
        1. For all users.\
        2. We check if there exists an implied role  role binding for created CRD. \
            if yes, \
            &nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;-   We go and retrieve all the implied roles from this role.\
            &nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;-   We check what has changed since last and assign/unassign all those roles to the user via deleting/Creating role bindings. (Also making sure we maintain labels for implied role)
    
- Delete the CRD that implies the implied role. &#8594;\
        1. For all users.\
        2. We check if there exists an implied role  role binding for deleted CRD.\
            if yes, \
            &nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;-   We go and retrieve all the implied roles from this role.\
            &nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;-   We un-assign these roles to the user via creating role bindings. (Also making sure we maintain labels for implied role)

---

## 2. Spec and Status

**================= SPEC =================**

```
type Rule struct {
	   ParentRole	string	`json:"parent_role,omitempty"`
	   ChildRole	string	`json:"child_role,omitempty"`
}

type ImpliedRolesSpec struct {	
    // An object of the form prior_role => implied_role_1
    InferenceRule Rule `json:"inference_rule,omitempty"`
}
```
**========================================**

**================= STATUS =================**

    ```
    type RuleInferences struct {
	    ParentRole	string	`json:"parent_role,omitempty"`
	    ChildRoles	[]string	`json:"child_roles,omitempty"`
    }

    type ImpliedRolesStatus struct {
    Inferences RuleInferences `json:"inferences,omitempty"`
    }
    ```

**==========================================**

---

## 3. Custon resource definition

**================= CRD =================**
```
    inference_rule:
	    parent_role: admin
	    child_role: writer
```
**=======================================**

---

## 4. Example

*What the operator does in updating the status* -

*Adding newer role inferences* 

Inferences: \
admin -> developer\
admin -> reviewer\
developer -> writer\
writer -> pro\
writer -> noob\

RuleInferences:\
admin -> [ developer, reviewer, writer, pro, noob ]\
developer -> [ writer, pro, noob ]\
writer -> [ pro, noob ]\


*Deleting role inferences* 

Infernece: 
developer -> writer\

RuleInferences:\
admin -> [ developer, reviewer]\
writer -> [ pro, noob ]\

---
