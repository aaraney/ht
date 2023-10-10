package main

import (
	"fmt"
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

	for _, child := range n.children {
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
