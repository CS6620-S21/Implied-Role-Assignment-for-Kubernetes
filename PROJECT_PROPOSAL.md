# Implied Role Assignment for Kubernetes

## 1. Vision and Goals of the Project

The RBAC (role-based access control) model needs to support improved delegation in order to scale. The goal is to transfer from coarse grained access control to finer-grained model, to be able to minimize the burden on the administrators, and grant a user all the privileges they need for a specific namespace with a single role assignment. High-level goals of RBAC include:

- A user should be able to delegate a subset of their capabilities
- Casual users should be able to provide access to a subset of their resources without providing full access to all of them
- Allow relationships between roles to be defined where one (larger more powerful role) can imply that the user also obtains a set of smaller (less powerful) roles
- Both explicitly assigned and implied roles will be included when listing roles assigned to a given user
- The model developed should has a generalizable pattern, which is reusable, and scales as our organization scales without having to redefine the rules for every project within our organization

---

## 2. Users/Personas of the Project

The aim of the project is to design and implement a system for implied roles within Kubernetes. Consequently, the following are the personas of the project.

The feature in its completed state will be provided as a Kubernetes operator that will be published on [OperatorHub.io](https://operatorhub.io/).

As a result, the end users are expected to be users 
- that administer or have permissions to administer Kuberentes.
- use the operator as published on OperatorHub for the Implied Role feature.


---

## 3. Scope and Features of the Project:

Users are assigned to roles within a namespace to perform operations. In order to better model the typical hierarchical authority model of a large organization, Users are assigned to roles within a namespace to perform operations. In order to better model the typical hierarchical authority model of a large organization, we will allow relationships between roles to be defined where one (larger more powerful role) can imply that the user also obtains a set of smaller (less powerful) roles.

The roles can be viewed as a hierarchy:
- Larger roles inherit the permissions assigned to smaller roles

For example, if a rule states that *admin* implies a *member*, any user assigned the *admin* role will automatically receive the *member* role as well.

The implementation avoids a strict hierarchy in favor of generating a directed-acyclic-graph (DAG): the same role may be implied by multiple prior roles. At enforcement time the required abstraction is a set of role assignments, not a tree or a graph.

The role relationships are illustrated in this ASCII DAG diagram. The prior roles are above implied roles, with the arrows showing the direction of implication. The table below that also explicitly shows these relationships:

![implied-role-dag](assets/dag.png)

For this project, we want to have the feature:
- Assign user new roles by creating implication rules so that the hierarchy of roles shown in the above pircture is created by creating/applying new implication rules, instead of directly assigning new roles to the user.

---

## 4. Solution Concept

A high-level orientation of the solution that is envisaged in order to meet the objectives of the architecture engagement -- 

Prerequisites - \
Workspace \
SQL backend (implicit roles assignment)

Steps:
1. Meet prerequisites
2. Create/Configure Groups
3. Configure cluster to use these groups
4. Define and assign permission to these groups using ClusterRole / RoleBinding
5. Troubleshooting and Debugging
6. Limitations

Design Implications :

\<Wireframes\>

key objectives \
requirements \
constraints for the engagement \
highlight work areas to be investigated in more detail


Discussion:

Design decisions made during the global architecture design.

Limitations:

Discussions and decisions taken, why?

---

## 5. Acceptance Criteria

- After creating/updating implied rules, user will be automatically assigned to implied roles according to the implied rules assigned. For example: if an user has developer role, and we create an implied rule that specifies that developer implies writer, and after applying this implied rule, now the user will also have writer role.
- When an implication rule is deleted, the corresponding implied role binding will be automatically deleted, without affecting user's original role binding. For example: if an user has developer role and implied writer role from the implication rule that states that developer implies writer. When the implication rule is deleted, the user will not have writer role anymore, but still keeps the developer role. 
---

## 6.  Release Planning:

### Project Preparation #1 (due by week 1)
- Understand what containers are
- Understand what Kubernetes is
- Make sure how to run Minikube
- Understand how the role access controls works under minikube in general

### Project Preparation #2 (due by week 2)
- Have a Minikube system up and running, with a CRD – customer resource descriptor
- A role and two different users – one who has that role can get access to it, and one who - doesn’t have that role doesn’t have access to it
- Get a minimal operator working on Minikube
