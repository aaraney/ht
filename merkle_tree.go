package main

import (
	"crypto/sha256"
	"fmt"
	"io"
	"log"
	"path/filepath"
	"slices"
	"strings"
)

// A Hasher takes in multiple strings and produces a single hash from the input.
// If a single input is provided it should be the output.
type Hasher interface {
	Hash(...string) string
}

type DefaultHasher struct{}

// Sorts input strings then hashes the sorted input using Sha256.
func (DefaultHasher) Hash(strs ...string) string {
	slices.Sort(strs)

	r := new(strings.Reader)
	h := sha256.New()
	for i := range strs {
		r.Reset(strs[i])
		_, err := io.Copy(h, r)
		if err != nil {
			//  TODO: improve error handling
			log.Fatal(err)
		}
	}
	hash := fmt.Sprintf("%x", h.Sum(nil))
	return hash
}

// Node in a Merkle Tree.
type Node struct {
	name     string
	children map[string]*Node
	val      string
}

type HashTree struct {
	hasher Hasher
	root   **Node
}

// Create a new HashTree that uses a given Hasher implementation or the DefaultHasher if a Hasher is
// not provided.
func NewHashTree(root string, h Hasher) HashTree {
	if h == nil {
		h = DefaultHasher{}
	}
	tree := NewTree(root)
	return HashTree{
		hasher: h,
		root:   &tree,
	}
}

// Uses the default hasher if an internal hasher is not present
func (h HashTree) BuildTree() {
	(*h.root).buildTree(h.hasher)
}

func (h HashTree) Add(filename, hashValue string) {
	(*h.root).Add(filename, hashValue)
}

func (h HashTree) String() string {
	return (*h.root).String()
}

// NewTree creates a empty root.
func NewTree(root string) *Node {
	return &Node{
		name:     filepath.Clean(root),
		children: make(map[string]*Node),
	}
}

// Returns if Node is a leaf node.
func (n *Node) Leaf() bool {
	return len(n.children) == 0
}

// Uses the default hasher if one is not provided.
func (n *Node) buildTree(h Hasher) string {
	if n == nil {
		panic("Cannot build tree from nil Node")
	} else if n.Leaf() {
		return n.val
	}

	if h == nil {
		h = DefaultHasher{}
	}

	var children []string
	for _, child := range n.children {
		children = append(children, child.buildTree(h))
	}

	var val string
	if len(children) == 1 {
		val = children[0]
	} else {
		val = h.Hash(children...)
	}
	n.val = val
	return val
}

type prefixAndNode struct {
	prefix string
	node   *Node
}

// TODO: I dont really like the name of this method.
func (n *Node) buildString(delm string) []string {
	if n == nil {
		panic("Cannot build string from nil Node")
	} else if n.Leaf() {
		return []string{fmt.Sprintf("%s %s", n.val, n.name)}
	}

	var current []prefixAndNode
	horizon := []prefixAndNode{{"", n}}
	var levelOuput []string
	var output []string

	for len(horizon) > 0 || len(current) > 0 {
		for i := range current {
			n := current[i]

			prefix := n.prefix
			node := n.node

			var hashAndName string

			if node.Leaf() {
				name := fmt.Sprintf("%s%s", prefix, node.name)
				hashAndName = fmt.Sprintf("%s %s", node.val, name)
			} else {
				name := fmt.Sprintf("%s%s%s", prefix, node.name, delm)
				hashAndName = fmt.Sprintf("%s %s", node.val, name)

				for _, child := range node.children {
					horizon = append(horizon, prefixAndNode{name, child})
				}
			}

			// TODO: sort before appending to output
			levelOuput = append(levelOuput, hashAndName)

		}
		slices.Sort(levelOuput)
		output = append(output, levelOuput...)
		levelOuput = nil

		current = horizon
		horizon = nil
	}

	return output
}

func (n *Node) String() string {
	out := n.buildString("/")
	return strings.Join(out, "\n")
}

type MerkleError string

func (m MerkleError) Error() string {
	return string(m)
}

const (
	EmptyFilePath       MerkleError = "empty filepath"
	NonRelativeFilePath MerkleError = "filepath not relative to root"
)

func (n *Node) Add(filePath string, hashValue string) error {
	if n == nil {
		panic("Cannot Add from nil Node")
	}

	if len(filePath) == 0 {
		return EmptyFilePath
	}

	filePath = filepath.Clean(filePath)

	relPath, err := filepath.Rel(n.name, filePath)
	if strings.HasPrefix(relPath, "../") || err != nil {
		return NonRelativeFilePath
	}

	split := strings.Split(relPath, "/")

	children := n.children
	var node *Node
	var exists bool
	for _, substr := range split {
		node, exists = children[substr]
		if !exists {
			// NOTE: may need to make adjustments here
			node = &Node{
				name:     substr,
				children: make(map[string]*Node),
			}
			children[substr] = node
		}
		children = node.children
	}
	// add metadata
	node.val = hashValue
	return nil
}
