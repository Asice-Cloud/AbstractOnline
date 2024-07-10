package utils

type TrieNode struct {
	children    map[rune]*TrieNode
	isEnd       bool
	permissions []string
}

type Trie struct {
	root *TrieNode
}

func NewTrie() *Trie {
	return &Trie{
		root: &TrieNode{children: make(map[rune]*TrieNode)},
	}
}

func (t *Trie) Insert(path string, permissions []string) {
	node := t.root
	for _, char := range path {
		if _, found := node.children[char]; !found {
			node.children[char] = &TrieNode{children: make(map[rune]*TrieNode)}
		}
		node = node.children[char]
	}
	node.isEnd = true
	node.permissions = permissions
}

func (t *Trie) CheckPermissions(path string, require []string) bool {
	node := t.root
	var maxMatchNode *TrieNode

	for _, char := range path {
		if nextnode, found := node.children[char]; found {
			node = nextnode
			if node.isEnd {
				maxMatchNode = node
			}
		} else {
			break
		}
	}
	if maxMatchNode == nil {
		return false
	}
	return hasRequire(maxMatchNode.permissions, require)
}

func hasRequire(nodePermissions, require []string) bool {
	permissionSet := make(map[string]bool)
	for _, perm := range nodePermissions {
		permissionSet[perm] = true
	}
	for _, req := range require {
		if !permissionSet[req] {
			return false
		}
	}
	return true
}

/*func main() {
	trie := NewTrie()
	trie.Insert("api/user/create", []string{"admin", "read", "write"})
	trie.Insert("api/user/delete", []string{"admin", "read", "delete"})
	trie.Insert("api/user/view", []string{"user", "read"})

	test := []struct {
		path    string
		require []string
		expect  bool
	}{
		{"api/user/create", []string{"admin", "read", "write"}, true},
		{"api/user/delete", []string{"admin", "delete"}, true},
		{"api/user/view", []string{"user", "read"}, true},
		{"api/user/view", []string{"admin", "read"}, true},
		{"api/user/create", []string{"admin", "read"}, false},
		{"api/user/create", []string{"write"}, false},
	}
	for _, tr := range test {
		result := trie.CheckPermissions(tr.path, tr.require)
		fmt.Printf("Path: %s, Require: %v, Result: %v (Expect: %v)\n", tr.path, tr.require, result, tr.expect)
	}
}
*/
