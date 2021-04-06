package main

import (
	"fmt"

	"container/list"
)

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
func transform(allRoleImplications map[string][]string, x map[string][]string, role string, irole string) (map[string][]string, map[string][]string) {

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
			_, found := Find(allRoleImplications[role], enew)
			if !found {
				allRoleImplications[role] = append(allRoleImplications[role], enew)
			}
		}

	}

	return allRoleImplications, x

}

// func (r *RoleImplicationRuleReconciler) GetAllRoleImplicationsForRoles(roleImplicationGraph map[string][]string) (map[string][]string, error) {

// 	allRoleImplications := make(map[string][]string)

// 	for role, irole := range roleImplicationGraph {
// 		for i := 0; i < len(irole); i++ {
// 			transform(role, irole[i])
// 		}
// 	}
// 	allRoleImplications = res
// 	return allRoleImplications, nil
// }

func main() {

	// Temp Map to store and compare the initial roles and implications.
	var Implicationgraph = make(map[string][]string)

	var allRoleImplications = make(map[string][]string)

	var roleImplicationGraph = map[string][]string{
		"admin":     {"developer", "reviewer"},
		"writer":    {"pro", "noob"},
		"developer": {"writer"},
	}

	for role, irole := range roleImplicationGraph {
		for i := 0; i < len(irole); i++ {
			allRoleImplications, Implicationgraph = transform(allRoleImplications, Implicationgraph, role, irole[i])
		}
	}
	fmt.Print(allRoleImplications)
}
