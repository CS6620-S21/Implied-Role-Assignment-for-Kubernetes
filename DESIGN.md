# Implied Role Assignment for Kubernetes

This is the design document for the implied roles operator.

## Overview

The operator is built as a combination of a **CRD** (Custom Resource Definition) and a **Controller** (written in Go). The design specifications of both components are outlined below.

## CRD

A **Custom Resource** (CR) contains a **Spec** and a **Status**. The **Spec** is effectively the schema for the desired state of implied roles in a cluster. The **Status** is a schema for the current state of implied roles in a cluster.

Each _implication_ (parent role implies child role) is a CR in the cluster.

### Spec

The spec is defined below. An implied roles CR contains a one to one mapping of `ParentRole` to `ChildRole`.

```golang
type Rule struct {
    ParentRole      string      `json:"parent_role,omitempty"`
    ChildRole       string      `json:"child_role,omitempty"`
}

type ImpliedRolesSpec struct {	
    InferenceRule   Rule        `json:"inference_rule,omitempty"`
}
```

### Status

The status is defined below. The `Inferences` contain a one-to-many mapping of `ParentRole` to multiple `ChildRoles`. The `ChildRoles` are all roles that a `ParentRole` would imply i.e. all roles that are descendants of `ParentRole` if one were to imagine all implication rules as a tree. 

```golang
type RuleInferences struct {
    ParentRole  string          `json:"parent_role,omitempty"`
    ChildRoles  []string        `json:"child_roles,omitempty"`
}

type ImpliedRolesStatus struct {
    Inferences  RuleInferences  `json:"inferences,omitempty"`
}
```

### CR Sample

This is an example of a implied roles CR defined in a `yaml`.

```yaml
kind: ImpliedRoles
metadata:
    name: impliedroles-sample
spec:
    inference_rule:
        parent_role: admin
        child_role: writer
```

## Controller

### Reconciler function

**Trigger:**

- Role-binding is created. &#8594; Role assignment (Through role binding yaml) \
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;Kubectl apply (Only to current User)
- Role-binding is updated. &#8594; Role assignment (Through role binding yaml) \ 
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;Kubectl apply (Only to current User)
- Role un assignment using yaml, deleted role binding - Need to look how\
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;kubectl delete (Only to current User)
- Create the CRD that implies the implied role.\
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;Create implied roles and create role bindings (Across all users)
- Update the CRD that implies the implied role.\
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;Update implied roles and create role bindings (Across all users)
- Delete the CRD that implies the implied role.\
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;Delete all implied roles and implied role bindings (Across all users)
- If you delete the role, it needs to be deleted from role binding and CRD \
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;NOT IN SCOPE (Look into this if kubernetes automatically handles this)

**Logic:**

- Role-binding is created &#8594;\
        1. We get the user from the context and the role assigned. \ 
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

## 4. Example

*What the operator does in updating the status* -

*Adding newer role inferences* 

Inferences: \
admin -> developer\
admin -> reviewer\
developer -> writer\
writer -> pro\
writer -> noob

RuleInferences:\
admin -> [ developer, reviewer, writer, pro, noob ]\
developer -> [ writer, pro, noob ]\
writer -> [ pro, noob ]


*Deleting role inferences* 

Inferences: 
developer -> writer

RuleInferences:\
admin -> [ developer, reviewer]\
writer -> [ pro, noob ]

---
