package main

import (
	"fmt"
	"slices"

	"github.com/jedib0t/go-pretty/v6/list"
)

func buildTreeView(prefix string, n *Node, l *list.List) {
	if n == nil {
		panic("Cannot build tree view from nil Root")
	}
	if l == nil {
		panic("Cannot build tree view from nil list.List")
	}

	var hashAndName string
	var name string
	if n.Leaf() {
		name = fmt.Sprintf("%s%s", prefix, n.name)
		hashAndName = fmt.Sprintf("%s %s", n.val, name)
	} else {
		name = fmt.Sprintf("%s%s%s", prefix, n.name, "/")
		hashAndName = fmt.Sprintf("%s %s", n.val, name)
	}
	l.AppendItem(hashAndName)

	var children []*Node
	for _, child := range n.children {
		children = append(children, child)
	}
	// sort by hash val then by name if needed to make levels deterministic
	slices.SortFunc(children, func(a, b *Node) int {
        // sort by val then by name
		if a.val == b.val {
			if a.name == b.name {
				return 0
			} else if a.name < b.name {
				return -1
			} else {
				return 1
			}
		} else if a.val < b.val {
			return -1
		} else {
			return 1
		}
	})

	for _, child := range children {
		l.Indent()
		buildTreeView(name, child, l)
		l.UnIndent()
	}
}

func BuildTreeView(tree HashTree, l *list.List) {
	if tree.root == nil {
		panic("Cannot build tree view from nil Root")
	}
	buildTreeView("", *tree.root, l)
}
