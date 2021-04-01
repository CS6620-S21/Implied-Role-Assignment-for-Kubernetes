package main

import (
	"fmt"

	"container/list"
)

// Temp Map to store and compare the initial roles and implications.
var x = make(map[string][]string)

// Final Map that holds all the roles an its implications.
var res = make(map[string][]string)

// checks if irole exists in a list of role that is assignd to the key
func Find(slice []string, val string) (int, bool) {
	for i, item := range slice {
		if item == val {
			return i, true
		}
	}
	return -1, false
}

// Creates the final Map for the reconciliation.
// Need to pass the role and its implicated role.
func transform(role string, irole string) {

	// add roles if they dont exist or append irole to existing role.
	x[role] = append(x[role], irole)

	// iterate throught the temp MAP
	for role, irole := range x {
		fmt.Println("k:", role, "v:", irole)
		q := list.New()
		q.PushBack(role)
		result := list.New()
		for q.Len() != 0 {
			st := q.Front().Value
			st_temp := q.Front()
			q.Remove(st_temp)
			fmt.Println(st)
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
			fmt.Println(e.Value)
			_, found := Find(res[role], enew)
			if !found {
				res[role] = append(res[role], enew)
			}
		}

	}

}

func (r *RoleImplicationRuleReconciler) GetAllRoleImplicationsForRoles(roleImplicationGraph map[string][]string) (map[string][]string, error) {
	allRoleImplications := make(map[string][]string)

	for role, irole := range roleImplicationGraph {
		for i := 0; i < len(irole); i++ {
			transform(role, irole[i])
		}
	}
	allRoleImplications = res
	return allRoleImplications, nil
}
