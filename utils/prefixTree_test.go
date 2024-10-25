package utils

import (
	"fmt"
	"testing"
)

func Test(t *testing.T) {
	trie := NewTrie()
	trie.Insert("api/user/create", []string{"admin", "write"})
	trie.Insert("api/user/delete", []string{"admin", "write"})
	trie.Insert("api/user/view", []string{"user", "view"})

	testCases := []struct {
		path                string
		requiredPermissions []string
		expectedResult      bool
	}{
		{"api/user/create", []string{"admin", "write"}, true},
		{"api/user/delete", []string{"admin", "delete"}, true},
		{"api/user/view", []string{"user", "view"}, true},
		{"api/user/view", []string{"admin", "view"}, false},
		{"api/user/create", []string{"write"}, false},
	}
	for _, tc := range testCases {
		result := trie.ChrckPermissions(tc.path, tc.requiredPermissions)
		fmt.Printf("Path: %s,required: %v,Result: %v (Expected: %v)\n",
			tc.path, tc.requiredPermissions, result, tc.expectedResult)
	}
}
