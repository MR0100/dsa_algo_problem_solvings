package main

import "fmt"

// TreeNode is a standard binary-tree node.
type TreeNode struct {
	Val   int
	Left  *TreeNode
	Right *TreeNode
}

// ── Approach 1: Path-to-Node + Divergence ────────────────────────────────────
//
// pathBased solves Lowest Common Ancestor of a Binary Tree by computing the
// root→p and root→q paths explicitly, then walking both from the top until they
// diverge; the last shared node is the LCA.
//
// Intuition:
//
//	The LCA is the deepest node that appears on BOTH the root→p path and the
//	root→q path. If we record each full path as a list of nodes, the answer is
//	simply the last position where the two lists agree.
//
// Algorithm:
//  1. Depth-first search to build the root→p path (list of *TreeNode).
//  2. Same for root→q.
//  3. Walk both lists in lockstep; keep the last node where they are equal.
//
// Time:  O(n) — two DFS traversals, each visiting at most every node.
// Space: O(n) — the two paths plus recursion stack (skewed tree).
func pathBased(root, p, q *TreeNode) *TreeNode {
	var pathTo func(node, target *TreeNode, path *[]*TreeNode) bool
	// pathTo appends nodes to *path while searching for target; returns true
	// once target is located so the caller keeps the built prefix.
	pathTo = func(node, target *TreeNode, path *[]*TreeNode) bool {
		if node == nil {
			return false // dead end, nothing added
		}
		*path = append(*path, node) // tentatively include this node on the path
		if node == target {
			return true // found — keep the path as-is
		}
		// Explore left then right; either subtree finding target keeps node.
		if pathTo(node.Left, target, path) || pathTo(node.Right, target, path) {
			return true
		}
		*path = (*path)[:len(*path)-1] // backtrack: node is not on the path
		return false
	}

	var pPath, qPath []*TreeNode
	pathTo(root, p, &pPath) // build root→p
	pathTo(root, q, &qPath) // build root→q

	var lca *TreeNode
	// Compare positions until one path ends or they diverge; keep last match.
	for i := 0; i < len(pPath) && i < len(qPath); i++ {
		if pPath[i] == qPath[i] {
			lca = pPath[i] // still on the shared prefix
		} else {
			break // paths diverged; previous match is the LCA
		}
	}
	return lca
}

// ── Approach 2: Single-Pass Recursion (Optimal) ──────────────────────────────
//
// recursive solves Lowest Common Ancestor of a Binary Tree with one post-order
// traversal that returns, for each subtree, whether p and/or q were found below.
//
// Intuition:
//
//	Ask each node: "how many of {p, q} live in my subtree?" A node is the LCA
//	exactly when the two targets are found on different sides (one left, one
//	right) or when the node itself is one target and the other lies below it.
//	The first node (deepest, since post-order bubbles up) satisfying this is the
//	answer.
//
// Algorithm:
//  1. Recurse. If node is nil / p / q, return node (a "found" signal).
//  2. left  = recurse(node.Left); right = recurse(node.Right).
//  3. If both non-nil, node is the split point → return node.
//  4. Otherwise return whichever side is non-nil (propagates the found target up).
//
// Time:  O(n) — each node visited once.
// Space: O(h) — recursion stack, h = tree height (O(n) worst, O(log n) balanced).
func recursive(root, p, q *TreeNode) *TreeNode {
	if root == nil || root == p || root == q {
		return root // base case: hit an empty branch or one of the targets
	}
	left := recursive(root.Left, p, q)   // search left subtree
	right := recursive(root.Right, p, q) // search right subtree
	if left != nil && right != nil {
		return root // p and q split here → this node is the LCA
	}
	if left != nil {
		return left // both targets (or the only found one) are on the left
	}
	return right // otherwise everything relevant is on the right (or nil)
}

// ── Approach 3: Parent Pointers + Ancestor Set ───────────────────────────────
//
// parentPointers solves Lowest Common Ancestor of a Binary Tree by recording
// every node's parent, then walking p's ancestors into a set and finding q's
// first ancestor already in that set.
//
// Intuition:
//
//	With a parent map we can climb from any node to the root. Collect all of
//	p's ancestors (including p); then climb from q — the first ancestor that is
//	also one of p's ancestors is the lowest common one.
//
// Algorithm:
//  1. BFS/DFS from root filling parent[node] for every node until both p and q
//     have been seen.
//  2. Climb from p to root, inserting each node into a set.
//  3. Climb from q; return the first node present in the set.
//
// Time:  O(n) — one traversal to build parents, then O(h) climbs.
// Space: O(n) — parent map and ancestor set.
func parentPointers(root, p, q *TreeNode) *TreeNode {
	parent := map[*TreeNode]*TreeNode{root: nil} // root has no parent
	stack := []*TreeNode{root}                   // explicit DFS stack
	// Traverse until we have recorded parents for both targets.
	for parent[p] == nil || parent[q] == nil {
		node := stack[len(stack)-1] // pop
		stack = stack[:len(stack)-1]
		if node.Left != nil {
			parent[node.Left] = node // record edge node→node.Left
			stack = append(stack, node.Left)
		}
		if node.Right != nil {
			parent[node.Right] = node // record edge node→node.Right
			stack = append(stack, node.Right)
		}
	}

	ancestors := map[*TreeNode]bool{} // p and all of p's ancestors
	for n := p; n != nil; n = parent[n] {
		ancestors[n] = true // climb p to root
	}
	// Climb q until we hit a node already on p's ancestor chain.
	for n := q; ; n = parent[n] {
		if ancestors[n] {
			return n // first common ancestor = lowest common ancestor
		}
	}
}

// buildSampleTree builds the tree from the LeetCode examples:
//
//	     3
//	   /   \
//	  5     1
//	 / \   / \
//	6   2 0   8
//	   / \
//	  7   4
//
// and returns root plus a lookup of node values → *TreeNode for convenience.
func buildSampleTree() (*TreeNode, map[int]*TreeNode) {
	n := func(v int) *TreeNode { return &TreeNode{Val: v} }
	nodes := map[int]*TreeNode{}
	for _, v := range []int{3, 5, 1, 6, 2, 0, 8, 7, 4} {
		nodes[v] = n(v)
	}
	nodes[3].Left, nodes[3].Right = nodes[5], nodes[1]
	nodes[5].Left, nodes[5].Right = nodes[6], nodes[2]
	nodes[1].Left, nodes[1].Right = nodes[0], nodes[8]
	nodes[2].Left, nodes[2].Right = nodes[7], nodes[4]
	return nodes[3], nodes
}

func main() {
	root, nodes := buildSampleTree()

	fmt.Println("=== Approach 1: Path-to-Node + Divergence ===")
	fmt.Println(pathBased(root, nodes[5], nodes[1]).Val) // expected 3
	fmt.Println(pathBased(root, nodes[5], nodes[4]).Val) // expected 5

	fmt.Println("=== Approach 2: Single-Pass Recursion (Optimal) ===")
	fmt.Println(recursive(root, nodes[5], nodes[1]).Val) // expected 3
	fmt.Println(recursive(root, nodes[5], nodes[4]).Val) // expected 5

	fmt.Println("=== Approach 3: Parent Pointers + Ancestor Set ===")
	fmt.Println(parentPointers(root, nodes[5], nodes[1]).Val) // expected 3
	fmt.Println(parentPointers(root, nodes[5], nodes[4]).Val) // expected 5

	// Example 3: root = [1,2], p = 1, q = 2 → LCA = 1
	root2 := &TreeNode{Val: 1, Left: &TreeNode{Val: 2}}
	fmt.Println("=== Example 3 (root=[1,2], p=1, q=2) ===")
	fmt.Println(recursive(root2, root2, root2.Left).Val) // expected 1
}
