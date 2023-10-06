package main

import (
	"testing"
)

// TODO: look at how to group tests
func TestBuildTree(t *testing.T) {
	root := NewTree("")

	filename := "somefile"
	// "feedbeef" hashed with sha256
	hashValue := "32549bff6d8404c4d121b589f4d24ac6416ed48c25163e1f08d92d67ca0bb0b3"

	root.Add(filename, hashValue)

	var h DefaultHasher
	root.buildTree(h)

	if len(root.children) != 1 {
		t.Errorf("expect a single child node. got %v", root.children)
	}

	node, ok := root.children[filename]
	if !ok {
		t.Errorf("expect filename: %q to be in map. map is: %v", filename, root.children)
	}
	if filename != node.name {
		t.Errorf("expect not filename: %q, got %q", filename, node.name)
	}
	if hashValue != node.val {
		t.Errorf("expect node val: %q, got %q", hashValue, node.val)
	}
}

func TestBuildTreeMultiple(t *testing.T) {
	root := NewTree(".")

	// leaf values are not "hash values" in this test. however when we build the tree their values
	// will be hashed.

	// NOTE: what about adding just a directory or adding a file with the same name as a directory?
	root.Add("./a", "a")
	root.Add("./b/c", "c")
	root.Add("./b/d", "d")

	// c and d will be concatenated and then hashed (e.g. hash("c" + "d")), this will become the
	// value of b's node.

	BHashValue := "21e721c35a5823fdb452fa2f9f0a612c74fb952e06927489c6b27a43b817bed4"

	// hash("c" + "d") = 21e721c35a5823fdb452fa2f9f0a612c74fb952e06927489c6b27a43b817bed4
	// hash(hash("c" + "d") + "a")
	rootHashValue := "735e521309a9c209bc8effbb7f7c90716a5eece7925249a2bf60dc631d7d5b93"

	var h DefaultHasher
	root.buildTree(h)

	if rootHashValue != root.val {
		t.Errorf("expect node val: %q, got %q", rootHashValue, root.val)
	}

	bNodeVal := root.children["b"].val
	if BHashValue != bNodeVal {
		t.Errorf("expect node val: %q, got %q", BHashValue, bNodeVal)
	}

	horizon := []*Node{root}

	for len(horizon) > 0 {
		node := horizon[0]
		t.Log(node.name)
		horizon = horizon[1:]
		for _, n := range node.children {
			horizon = append(horizon, n)
		}
	}
}

func TestCannotAddNonRelativeNode(t *testing.T) {
	root := NewTree("/a")
	err := root.Add("/b", "")
	if err == nil {
		t.Error(err)
	}
}
