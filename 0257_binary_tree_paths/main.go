package main

import (
	"fmt"
	"strconv"
	"strings"
)

// TreeNode is a binary tree node.
type TreeNode struct {
	Val   int
	Left  *TreeNode
	Right *TreeNode
}

// ── Approach 1: DFS with String Accumulation ─────────────────────────────────
//
// dfsStringAccum solves Binary Tree Paths by carrying the path built so far as
// a string down the recursion.
//
// Intuition:
//
//	A root-to-leaf path is exactly the sequence of node values from root to a
//	node with no children. Walk down; append each value to the running path
//	string; when we hit a leaf, the running string is a complete answer.
//
// Algorithm:
//  1. dfs(node, path): append node.Val to path.
//  2. If node is a leaf (no children), append path to results.
//  3. Else recurse into non-nil children with path + "->".
//
// Time:  O(n * h) — n nodes visited; building each of up to O(n) path strings
//
//	costs up to O(h) where h is height (path length).
//
// Space: O(n * h) — output strings; recursion stack O(h).
func dfsStringAccum(root *TreeNode) []string {
	var result []string
	if root == nil {
		return result // empty tree → no paths
	}
	var dfs func(node *TreeNode, path string)
	dfs = func(node *TreeNode, path string) {
		// Append this node's value to the path built so far.
		path += strconv.Itoa(node.Val)
		if node.Left == nil && node.Right == nil { // leaf → path is complete
			result = append(result, path)
			return
		}
		if node.Left != nil {
			dfs(node.Left, path+"->") // extend with arrow then recurse
		}
		if node.Right != nil {
			dfs(node.Right, path+"->")
		}
	}
	dfs(root, "")
	return result
}

// ── Approach 2: DFS with Backtracking Slice ──────────────────────────────────
//
// dfsBacktrack solves Binary Tree Paths using a shared slice of node values,
// pushing on entry and popping on exit (classic backtracking).
//
// Intuition:
//
//	Instead of copying a growing string at every level, keep one slice of the
//	current path's values. Push the value when entering a node, pop it when
//	leaving. At a leaf, join the slice into the "a->b->c" form.
//
// Algorithm:
//  1. dfs(node): push node.Val onto path slice.
//  2. If leaf, join path with "->" and record.
//  3. Recurse into non-nil children.
//  4. Pop node.Val before returning (restore state).
//
// Time:  O(n * h) — each of up to O(n) leaves joins an O(h)-length slice.
// Space: O(h) auxiliary (the path slice + recursion), plus O(n*h) output.
func dfsBacktrack(root *TreeNode) []string {
	var result []string
	if root == nil {
		return result
	}
	path := []string{} // current root-to-node values as strings
	var dfs func(node *TreeNode)
	dfs = func(node *TreeNode) {
		path = append(path, strconv.Itoa(node.Val)) // push current value
		if node.Left == nil && node.Right == nil {  // leaf
			result = append(result, strings.Join(path, "->")) // snapshot the path
		} else {
			if node.Left != nil {
				dfs(node.Left)
			}
			if node.Right != nil {
				dfs(node.Right)
			}
		}
		path = path[:len(path)-1] // pop: undo this node before returning
	}
	dfs(root)
	return result
}

// ── Approach 3: Iterative DFS with Explicit Stack (Optimal) ──────────────────
//
// iterativeStack solves Binary Tree Paths without recursion, pairing each node
// with the path string that leads to it on an explicit stack.
//
// Intuition:
//
//	Any recursive DFS can be made iterative by carrying the "call state" on a
//	stack. Here the state per node is the path string ending at it. Pop a node,
//	extend its path; if leaf, emit; else push children with the extended path.
//
// Algorithm:
//  1. Push (root, "root's value") onto the stack.
//  2. While stack non-empty: pop (node, path).
//  3. If node is a leaf, add path to results.
//  4. Else push children (right then left, so left pops first) with
//     path + "->" + childVal.
//
// Time:  O(n * h) — same work as recursion.
// Space: O(n * h) — stack can hold O(n) node/path pairs.
func iterativeStack(root *TreeNode) []string {
	var result []string
	if root == nil {
		return result
	}
	type frame struct {
		node *TreeNode
		path string
	}
	// Seed the stack with the root and its own value as the starting path.
	stack := []frame{{root, strconv.Itoa(root.Val)}}
	for len(stack) > 0 {
		top := stack[len(stack)-1]   // peek
		stack = stack[:len(stack)-1] // pop
		node, path := top.node, top.path
		if node.Left == nil && node.Right == nil { // leaf → complete path
			result = append(result, path)
			continue
		}
		// Push RIGHT first so LEFT is popped first (stack is LIFO); this makes
		// the output order match a left-to-right recursive DFS.
		if node.Right != nil {
			stack = append(stack, frame{node.Right, path + "->" + strconv.Itoa(node.Right.Val)})
		}
		if node.Left != nil {
			stack = append(stack, frame{node.Left, path + "->" + strconv.Itoa(node.Left.Val)})
		}
	}
	return result
}

func main() {
	// Example 1:      1
	//               /   \
	//              2     3
	//               \
	//                5
	// Expected paths: ["1->2->5","1->3"]
	ex1 := &TreeNode{
		Val:   1,
		Left:  &TreeNode{Val: 2, Right: &TreeNode{Val: 5}},
		Right: &TreeNode{Val: 3},
	}
	// Example 2: single node [1] → ["1"]
	ex2 := &TreeNode{Val: 1}

	fmt.Println("=== Approach 1: DFS String Accumulation ===")
	fmt.Println(dfsStringAccum(ex1)) // expected [1->2->5 1->3]
	fmt.Println(dfsStringAccum(ex2)) // expected [1]

	fmt.Println("=== Approach 2: DFS Backtracking Slice ===")
	fmt.Println(dfsBacktrack(ex1)) // expected [1->2->5 1->3]
	fmt.Println(dfsBacktrack(ex2)) // expected [1]

	fmt.Println("=== Approach 3: Iterative DFS with Explicit Stack (Optimal) ===")
	fmt.Println(iterativeStack(ex1)) // expected [1->2->5 1->3]
	fmt.Println(iterativeStack(ex2)) // expected [1]
}
