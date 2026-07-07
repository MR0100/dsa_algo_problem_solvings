package main

import (
	"fmt"
	"strconv"
	"strings"
)

// Serialize and Deserialize Binary Tree
//
// Design an algorithm to serialize a binary tree to a string and deserialize
// that string back into the identical tree structure. There is no restriction
// on the encoding format.

// TreeNode is the standard binary-tree node.
type TreeNode struct {
	Val   int
	Left  *TreeNode
	Right *TreeNode
}

// ── Approach 1: Preorder DFS Codec (Optimal) ─────────────────────────────────
//
// PreorderCodec serializes with a preorder DFS, writing "#" for nil children.
//
// Intuition:
//
//	A preorder walk that ALSO records nil children uniquely determines the tree:
//	you always know, while reading, whether the current position is a real node
//	or an absent child, so no second traversal is needed to reconstruct.
//
// Algorithm (serialize):
//  1. Visit root; if nil append "#", else append its value.
//  2. Recurse left, then right.
//  3. Join tokens with commas.
//
// Algorithm (deserialize):
//  1. Split into tokens; keep an index.
//  2. Read one token: if "#", return nil; else make a node.
//  3. Recursively build left then right (preorder consumes in the same order).
//
// Time:  O(n) serialize + O(n) deserialize — each node visited once.
// Space: O(n) — output string plus O(h) recursion stack.
type PreorderCodec struct{}

// serialize turns the tree into a comma-separated preorder string.
func (PreorderCodec) serialize(root *TreeNode) string {
	var sb strings.Builder
	var dfs func(node *TreeNode)
	dfs = func(node *TreeNode) {
		if node == nil {
			sb.WriteString("#,") // marker for an absent child
			return
		}
		sb.WriteString(strconv.Itoa(node.Val)) // record this node's value
		sb.WriteByte(',')
		dfs(node.Left)  // preorder: left subtree next
		dfs(node.Right) // then right subtree
	}
	dfs(root)
	return sb.String()
}

// deserialize rebuilds the tree from the preorder string.
func (PreorderCodec) deserialize(data string) *TreeNode {
	tokens := strings.Split(data, ",") // may end with an empty token from trailing comma
	idx := 0
	var build func() *TreeNode
	build = func() *TreeNode {
		tok := tokens[idx] // read the next token in preorder order
		idx++
		if tok == "#" {
			return nil // absent child
		}
		val, _ := strconv.Atoi(tok)
		node := &TreeNode{Val: val}
		node.Left = build()  // consume the left subtree's tokens
		node.Right = build() // then the right subtree's tokens
		return node
	}
	return build()
}

// ── Approach 2: BFS Level-Order Codec ────────────────────────────────────────
//
// BFSCodec serializes using a level-order (queue) traversal, appending "#" for
// nil positions, mirroring LeetCode's own display format.
//
// Intuition:
//
//	A breadth-first sweep records nodes level by level; writing "#" for each
//	missing child preserves shape. On the way back, a queue re-links children in
//	the exact order they were written.
//
// Algorithm (serialize):
//  1. Push root into a queue.
//  2. Pop a node: if nil append "#", else append value and push both children.
//  3. Repeat until the queue empties.
//
// Algorithm (deserialize):
//  1. First token is the root; enqueue it.
//  2. Pop a parent; the next two tokens are its left/right children — create
//     non-nil ones, enqueue them, and attach.
//
// Time:  O(n) serialize + O(n) deserialize.
// Space: O(n) — queue and output.
type BFSCodec struct{}

func (BFSCodec) serialize(root *TreeNode) string {
	if root == nil {
		return "" // empty tree serializes to the empty string
	}
	var out []string
	queue := []*TreeNode{root}
	for len(queue) > 0 {
		node := queue[0] // dequeue the front
		queue = queue[1:]
		if node == nil {
			out = append(out, "#") // absent slot
			continue
		}
		out = append(out, strconv.Itoa(node.Val))
		queue = append(queue, node.Left)  // enqueue children (possibly nil)
		queue = append(queue, node.Right) // to be resolved on later pops
	}
	return strings.Join(out, ",")
}

func (BFSCodec) deserialize(data string) *TreeNode {
	if data == "" {
		return nil // empty string -> empty tree
	}
	tokens := strings.Split(data, ",")
	root := &TreeNode{Val: mustAtoi(tokens[0])}
	queue := []*TreeNode{root}
	i := 1 // index of the next child token to consume
	for len(queue) > 0 && i < len(tokens) {
		parent := queue[0] // parent whose children we resolve now
		queue = queue[1:]

		if tokens[i] != "#" { // left child present
			parent.Left = &TreeNode{Val: mustAtoi(tokens[i])}
			queue = append(queue, parent.Left)
		}
		i++
		if i < len(tokens) && tokens[i] != "#" { // right child present
			parent.Right = &TreeNode{Val: mustAtoi(tokens[i])}
			queue = append(queue, parent.Right)
		}
		i++
	}
	return root
}

func mustAtoi(s string) int {
	n, _ := strconv.Atoi(s)
	return n
}

// ── Test helpers ─────────────────────────────────────────────────────────────

// levelOrder produces a compact level-order view (with nulls trimmed at the
// tail) so we can confirm the round-tripped tree matches the original.
func levelOrder(root *TreeNode) []string {
	if root == nil {
		return []string{}
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
	// trim trailing "null" tokens for a clean, canonical comparison
	for len(out) > 0 && out[len(out)-1] == "null" {
		out = out[:len(out)-1]
	}
	return out
}

func main() {
	// Example 1: root = [1,2,3,null,null,4,5]
	//        1
	//       / \
	//      2   3
	//         / \
	//        4   5
	root1 := &TreeNode{
		Val:  1,
		Left: &TreeNode{Val: 2},
		Right: &TreeNode{
			Val:   3,
			Left:  &TreeNode{Val: 4},
			Right: &TreeNode{Val: 5},
		},
	}
	// Example 2: root = [] (empty tree)
	var root2 *TreeNode

	fmt.Println("=== Approach 1: Preorder DFS Codec ===")
	pc := PreorderCodec{}
	s1 := pc.serialize(root1)
	fmt.Println(fmt.Sprint(levelOrder(pc.deserialize(s1)))) // expected [1 2 3 null null 4 5]
	s2 := pc.serialize(root2)
	fmt.Println(fmt.Sprint(levelOrder(pc.deserialize(s2)))) // expected []

	fmt.Println("=== Approach 2: BFS Level-Order Codec ===")
	bc := BFSCodec{}
	b1 := bc.serialize(root1)
	fmt.Println(b1)                                         // expected 1,2,3,#,#,4,5,#,#,#,#
	fmt.Println(fmt.Sprint(levelOrder(bc.deserialize(b1)))) // expected [1 2 3 null null 4 5]
	b2 := bc.serialize(root2)
	fmt.Println(fmt.Sprint(levelOrder(bc.deserialize(b2)))) // expected []
}
