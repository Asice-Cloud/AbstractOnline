package utils

// prefix tree node
type TrieNode struct {
	children    map[rune]*TrieNode
	isEnd       bool
	permissions []string
}

// prefix tree
type Trie struct {
	root *TrieNode
}

func NewTrie() *Trie {
	return &Trie{root: &TrieNode{children: make(map[rune]*TrieNode)}}
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

func (t *Trie) Search(path string) (*TrieNode, bool) {
	node := t.root
	for _, char := range path {
		if _, found := node.children[char]; !found {
			return nil, false
		}
	}
	return node, node.isEnd
}

func (t *Trie) ChrckPermissions(path string, requiredPermissions []string) bool {
	node := t.root
	var maxMatchNode *TrieNode
	for _, char := range path {
		if nextNode, found := node.children[char]; found {
			node = nextNode
			if node.isEnd {
				maxMatchNode = node
			} else {
				break
			}
		}
	}
	if maxMatchNode == nil {
		return false
	}

	return hasRequiredPermissions(maxMatchNode.permissions, requiredPermissions)
}

func hasRequiredPermissions(nodePermissions, requiredPermissions []string) bool {
	permSet := make(map[string]bool)
	for _, perm := range nodePermissions {
		permSet[perm] = true
	}
	for _, reqPerm := range requiredPermissions {
		if !permSet[reqPerm] {
			return false
		}
	}
	return true
}
