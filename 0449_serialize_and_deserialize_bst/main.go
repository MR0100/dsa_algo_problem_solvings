package main

import (
	"fmt"
	"strconv"
	"strings"
)

// TreeNode is the standard LeetCode binary-tree node.
type TreeNode struct {
	Val   int
	Left  *TreeNode
	Right *TreeNode
}

// codec is the interface every approach implements so main() can drive the same
// serialize → deserialize round-trip through all of them.
type codec interface {
	serialize(root *TreeNode) string
	deserialize(data string) *TreeNode
}

// ── Approach 1: Preorder with Null Markers ───────────────────────────────────
//
// PreorderNullCodec serializes ANY binary tree (BST property unused) via a
// preorder walk that writes an explicit sentinel for every missing child.
//
// Intuition:
//
//	A preorder listing alone is ambiguous for a general binary tree, but adding
//	a marker ("#") wherever a child is nil makes the shape unambiguous: on the
//	way back the same preorder consumption rebuilds exactly one tree. It ignores
//	the BST ordering entirely, so it is the most general (and most verbose)
//	codec — the baseline before we exploit the BST structure.
//
// Algorithm (serialize): preorder DFS, emit node value or "#" for nil, space-joined.
// Algorithm (deserialize): read tokens left-to-right; "#" → nil, else build node
//
//	then recursively build left, then right (preorder order).
//
// Time:  O(n) serialize and O(n) deserialize — each node/marker visited once.
// Space: O(n) for the string plus O(h) recursion stack.
type PreorderNullCodec struct{}

// serialize writes a preorder traversal with "#" sentinels for nil children.
func (PreorderNullCodec) serialize(root *TreeNode) string {
	var sb strings.Builder
	var pre func(node *TreeNode)
	pre = func(node *TreeNode) {
		if node == nil {
			sb.WriteString("# ") // sentinel marks an absent child
			return
		}
		sb.WriteString(strconv.Itoa(node.Val)) // node value first (preorder)
		sb.WriteByte(' ')
		pre(node.Left)  // then entire left subtree
		pre(node.Right) // then entire right subtree
	}
	pre(root)
	return strings.TrimSpace(sb.String())
}

// deserialize rebuilds the tree by consuming the preorder token stream.
func (PreorderNullCodec) deserialize(data string) *TreeNode {
	if data == "" {
		return nil
	}
	tokens := strings.Fields(data) // split on whitespace into value/"#" tokens
	pos := 0                       // cursor into tokens, advanced as we consume
	var build func() *TreeNode
	build = func() *TreeNode {
		tok := tokens[pos] // current token dictates node vs nil
		pos++
		if tok == "#" {
			return nil // sentinel → no node here
		}
		val, _ := strconv.Atoi(tok)
		node := &TreeNode{Val: val}
		node.Left = build()  // preorder: left subtree consumed next
		node.Right = build() // then right subtree
		return node
	}
	return build()
}

// ── Approach 2: Preorder-Only, BST-Bounded Rebuild (Optimal) ─────────────────
//
// BSTPreorderCodec serializes a BST as a bare preorder value list (NO null
// markers) and reconstructs it using the BST ordering to place each value.
//
// Intuition:
//
//	For a BST the preorder sequence alone determines the tree — no sentinels
//	needed — because the search-tree ordering tells us where each value belongs.
//	The first value is the root. Everything smaller than it (a contiguous prefix
//	of the remaining preorder) forms the left subtree; the rest forms the right.
//	Reconstruct recursively with an (lower, upper) value window: consume the
//	next value only while it fits the current node's allowed range. This drops
//	the markers, roughly halving the payload, and rebuilds in O(n) via a moving
//	cursor + bounds.
//
// Algorithm (serialize): plain preorder, values only, space-joined.
// Algorithm (deserialize): cursor over values; build(lower, upper) takes the
//
//	next value if it lies in (lower, upper), makes it a node, then recursively
//	builds left with upper = val and right with lower = val.
//
// Time:  O(n) serialize; O(n) deserialize — each value consumed once, bounds are O(1).
// Space: O(n) for the string plus O(h) recursion stack.
type BSTPreorderCodec struct{}

// serialize writes just the preorder values — BST ordering makes markers unnecessary.
func (BSTPreorderCodec) serialize(root *TreeNode) string {
	var vals []string
	var pre func(node *TreeNode)
	pre = func(node *TreeNode) {
		if node == nil {
			return // no marker: absence is inferred from value bounds later
		}
		vals = append(vals, strconv.Itoa(node.Val)) // preorder value
		pre(node.Left)
		pre(node.Right)
	}
	pre(root)
	return strings.Join(vals, " ")
}

// deserialize rebuilds the BST from the preorder values using value-range bounds.
func (BSTPreorderCodec) deserialize(data string) *TreeNode {
	if data == "" {
		return nil
	}
	fields := strings.Fields(data)
	vals := make([]int, len(fields))
	for i, f := range fields {
		vals[i], _ = strconv.Atoi(f) // parse preorder values once
	}
	pos := 0 // cursor into vals; only moves when a value is placed

	// build constructs the subtree whose node values must fall in (lower, upper).
	// Sentinels ±∞ bound the root.
	var build func(lower, upper int) *TreeNode
	build = func(lower, upper int) *TreeNode {
		if pos == len(vals) {
			return nil // stream exhausted
		}
		v := vals[pos] // peek the next preorder value
		if v < lower || v > upper {
			return nil // out of this node's allowed window → belongs elsewhere
		}
		pos++ // consume v as the current subtree's root
		node := &TreeNode{Val: v}
		node.Left = build(lower, v)  // left subtree values must be < v
		node.Right = build(v, upper) // right subtree values must be > v
		return node
	}
	// Use a wide integer window as ±∞ (values fit within it per constraints).
	return build(-1<<62, 1<<62)
}

// levelOrder renders a tree as a LeetCode-style level-order list (with "null"
// for missing children, trailing nulls trimmed) purely for verifying round-trips.
func levelOrder(root *TreeNode) string {
	if root == nil {
		return "[]"
	}
	var out []string
	queue := []*TreeNode{root}
	for len(queue) > 0 {
		node := queue[0]
		queue = queue[1:]
		if node == nil {
			out = append(out, "null")
			continue
		}
		out = append(out, strconv.Itoa(node.Val))
		queue = append(queue, node.Left, node.Right)
	}
	// Trim trailing "null"s to match LeetCode's compact display.
	for len(out) > 0 && out[len(out)-1] == "null" {
		out = out[:len(out)-1]
	}
	return "[" + strings.Join(out, ",") + "]"
}

// buildBST inserts values (in the given order) into a BST — a convenient way to
// construct the example trees for main().
func buildBST(vals ...int) *TreeNode {
	var root *TreeNode
	var insert func(node *TreeNode, v int) *TreeNode
	insert = func(node *TreeNode, v int) *TreeNode {
		if node == nil {
			return &TreeNode{Val: v}
		}
		if v < node.Val {
			node.Left = insert(node.Left, v)
		} else {
			node.Right = insert(node.Right, v)
		}
		return node
	}
	for _, v := range vals {
		root = insert(root, v)
	}
	return root
}

// roundTrip serializes then deserializes with the given codec and returns the
// resulting tree's level-order form.
func roundTrip(c codec, root *TreeNode) string {
	return levelOrder(c.deserialize(c.serialize(root)))
}

func main() {
	// Example 1: BST from insertions 2,1,3 → root 2, left 1, right 3 → [2,1,3].
	ex1 := buildBST(2, 1, 3)
	// A richer BST: insert 5,3,6,2,4,7 → [5,3,6,2,4,null,7].
	ex2 := buildBST(5, 3, 6, 2, 4, 7)
	// Empty tree.
	var ex3 *TreeNode

	fmt.Println("=== Approach 1: Preorder with Null Markers ===")
	c1 := PreorderNullCodec{}
	fmt.Printf("serialize([2,1,3])            = %q\n", c1.serialize(ex1)) // "2 1 # # 3 # #"
	fmt.Printf("round-trip [2,1,3]            → %s  expected [2,1,3]\n", roundTrip(c1, ex1))
	fmt.Printf("round-trip [5,3,6,2,4,null,7] → %s  expected [5,3,6,2,4,null,7]\n", roundTrip(c1, ex2))
	fmt.Printf("round-trip []                 → %s  expected []\n", roundTrip(c1, ex3))

	fmt.Println("=== Approach 2: Preorder-Only, BST-Bounded (Optimal) ===")
	c2 := BSTPreorderCodec{}
	fmt.Printf("serialize([2,1,3])            = %q\n", c2.serialize(ex1)) // "2 1 3"
	fmt.Printf("round-trip [2,1,3]            → %s  expected [2,1,3]\n", roundTrip(c2, ex1))
	fmt.Printf("round-trip [5,3,6,2,4,null,7] → %s  expected [5,3,6,2,4,null,7]\n", roundTrip(c2, ex2))
	fmt.Printf("round-trip []                 → %s  expected []\n", roundTrip(c2, ex3))
}
