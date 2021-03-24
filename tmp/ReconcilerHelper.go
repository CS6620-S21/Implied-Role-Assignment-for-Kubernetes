package main

import (
	"fmt"
)

func Difference(a, b []string) (diff []string) {
	m := make(map[string]bool)

	for _, item := range b {
		m[item] = true
	}

	for _, item := range a {
		if _, ok := m[item]; !ok {
			diff = append(diff, item)
		}
	}
	return
}

func main() {
	var current_implied_role = []string{"create", "update", "patch", "delete", "watch"}
	var existing_roles = []string{"watch","patch"}
	fmt.Println(Difference(current_implied_role , existing_roles ))
}