# 0285 — Inorder Successor in BST

> LeetCode #285 · Difficulty: Medium
> **Categories:** Tree, Binary Search Tree, Depth-First Search, Binary Tree

---

## Problem Statement

Given the `root` of a binary search tree and a node `p` in it, return the **in-order successor** of that node in the BST. If the given node has no in-order successor in the tree, return `null`.

The successor of a node `p` is the node with the smallest key greater than `p.val`.

**Example 1:**
```
Input: root = [2,1,3], p = 1
Output: 2
Explanation: 1's in-order successor node is 2. Note that both p and the return value is of TreeNode type.
```

**Example 2:**
```
Input: root = [5,3,6,2,4,null,null,1], p = 6
Output: null
Explanation: There is no in-order successor of the current node, so the answer is null.
```

**Constraints:**
- The number of nodes in the tree is in the range `[1, 10⁴]`.
- `-10⁵ <= Node.val <= 10⁵`
- All Nodes will have unique values.

---

## Company Frequency

| Company    | Frequency        | Last Reported |
|------------|------------------|---------------|
| Meta       | ★★★★☆ High       | 2024          |
| Amazon     | ★★★★☆ High       | 2024          |
| Microsoft  | ★★★☆☆ Medium     | 2023          |
| Google     | ★★★☆☆ Medium     | 2023          |
| Bloomberg  | ★★☆☆☆ Low        | 2023          |

> ⚠️ Frequency data is crowd-sourced from LeetCode Discuss, Glassdoor, and
> community interview reports. Treat as a signal, not a guarantee.

---

## DSA Concepts Used
- **Binary search tree ordering** — left < node < right lets us find the successor without visiting both subtrees → see [`/dsa/binary_search_tree.md`](/dsa/binary_search_tree.md)
- **In-order traversal** — the successor is literally the next node in ascending order → see [`/dsa/tree_traversal.md`](/dsa/tree_traversal.md)

---

## Approaches Overview
| # | Approach | Time | Space | When to use |
|---|----------|------|-------|-------------|
| 1 | In-Order Traversal (Brute Force) | O(n) | O(n) | Baseline; also works on non-BSTs |
| 2 | BST Property, Iterative (Optimal) | O(h) | O(1) | Best; single descent, no `parent` pointers |
| 3 | Right-Subtree + Ancestor (explicit cases) | O(h) | O(1) | Clarifies the two-case reasoning |

---

## Approach 1 — In-Order Traversal

### Intuition
The in-order traversal of a BST lists nodes in ascending value. The successor of `p` is simply the node right after `p` in that list. Flatten the tree to a sorted list of node pointers, find `p`, and return the next one (or `nil` if `p` is last).

### Algorithm
1. In-order traverse, appending each node pointer to a slice.
2. Scan the slice for `p`; return the element immediately after it.
3. If `p` is the last element, return `nil`.

### Complexity
- **Time:** O(n) — visits every node.
- **Space:** O(n) — the flattened list plus O(h) recursion.

### Code
```go
func inorderTraversal(root, p *TreeNode) *TreeNode {
	order := []*TreeNode{}
	var walk func(*TreeNode)
	walk = func(n *TreeNode) {
		if n == nil {
			return
		}
		walk(n.Left)
		order = append(order, n)
		walk(n.Right)
	}
	walk(root)
	for i, n := range order {
		if n == p && i+1 < len(order) {
			return order[i+1]
		}
	}
	return nil
}
```

### Dry Run
Example 1: `root = [2,1,3]`, `p = node(1)`.

| step | in-order list built |
|------|---------------------|
| visit left of 2 → 1 | [1] |
| visit 2 | [1,2] |
| visit right of 2 → 3 | [1,2,3] |

Find `p = 1` at index 0; next element is `2` ⇒ return node `2`.

---

## Approach 2 — BST Property, Iterative (Optimal)

### Intuition
The successor is the smallest node whose value is strictly greater than `p.val`. Walk down from the root: whenever `p.val < cur.val`, `cur` is a valid (larger) candidate — remember it and go left to search for an even smaller valid value; otherwise `cur <= p`, so the successor must be larger — go right. The last remembered candidate is the answer. This never visits both subtrees, giving O(h).

### Algorithm
1. `successor = nil`, `cur = root`.
2. While `cur != nil`:
   - If `p.val < cur.val`: `successor = cur`; `cur = cur.Left`.
   - Else: `cur = cur.Right`.
3. Return `successor`.

### Complexity
- **Time:** O(h) — one root-to-leaf descent (`h` = height).
- **Space:** O(1) — iterative, no recursion or parent pointers.

### Code
```go
func bstSearch(root, p *TreeNode) *TreeNode {
	var successor *TreeNode
	cur := root
	for cur != nil {
		if p.Val < cur.Val {
			successor = cur
			cur = cur.Left
		} else {
			cur = cur.Right
		}
	}
	return successor
}
```

### Dry Run
Example 1: `root = 2 (left 1, right 3)`, `p.Val = 1`.

| iter | cur | compare `1 < cur.Val` | successor | move |
|------|-----|-----------------------|-----------|------|
| 1 | 2 | 1 < 2 → true | 2 | go left |
| 2 | 1 | 1 < 1 → false | 2 | go right |
| 3 | nil | — | 2 | loop ends |

Return `successor = node(2)`.

(Example 2: `p.Val = 6` is the max; every comparison `6 < cur.Val` is false, so `successor` stays `nil`.)

---

## Approach 3 — Right-Subtree + Ancestor (Explicit Cases)

### Intuition
Same ordering facts, spelled out as two cases:
- If `p` **has a right subtree**, the successor is that subtree's leftmost node (smallest value still larger than `p`).
- If `p` has **no right subtree**, the successor is the nearest ancestor from which we turned left descending toward `p`.

### Algorithm
1. If `p.Right != nil`: `cur = p.Right`; while `cur.Left != nil` `cur = cur.Left`; return `cur`.
2. Else walk from root comparing values; each time `p.Val < cur.Val` record `successor = cur` and go left, else go right; stop at `p`; return `successor`.

### Complexity
- **Time:** O(h).
- **Space:** O(1).

### Code
```go
func rightSubtreeCase(root, p *TreeNode) *TreeNode {
	if p.Right != nil {
		cur := p.Right
		for cur.Left != nil {
			cur = cur.Left
		}
		return cur
	}
	var successor *TreeNode
	cur := root
	for cur != nil && cur.Val != p.Val {
		if p.Val < cur.Val {
			successor = cur
			cur = cur.Left
		} else {
			cur = cur.Right
		}
	}
	return successor
}
```

### Dry Run
Example 1: `p = node(1)`, `p.Right == nil` → Case 2. Walk from root:

| iter | cur | `1 < cur.Val`? | successor | move |
|------|-----|----------------|-----------|------|
| 1 | 2 | true | 2 | go left |
| 2 | 1 | cur.Val == p.Val → stop | 2 | — |

Return `successor = node(2)`.

---

## Key Takeaways
- **Successor via descent, not parent pointers:** track the last node you turned left from; it's the smallest value `> p.val` seen so far. This is the O(h)/O(1) go-to.
- The **two-case view** (right subtree present vs absent) is a useful mental model, but the single-loop "candidate" version handles both uniformly in less code.
- Because it's a BST, you never need to visit both children — that's what turns O(n) traversal into O(h) search.

---

## Related Problems
- LeetCode #510 — Inorder Successor in BST II (with parent pointers, no root)
- LeetCode #173 — Binary Search Tree Iterator (repeated successor calls)
- LeetCode #700 — Search in a Binary Search Tree (same descent skeleton)
- LeetCode #235 — Lowest Common Ancestor of a BST (descent using ordering)
