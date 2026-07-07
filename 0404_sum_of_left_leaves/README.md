# 0404 — Sum of Left Leaves

> LeetCode #404 · Difficulty: Easy
> **Categories:** Tree, Depth-First Search, Breadth-First Search, Binary Tree

---

## Problem Statement

Given the `root` of a binary tree, return *the sum of all left leaves.*

A **leaf** is a node with no children. A **left leaf** is a leaf that is the left child of another node.

**Example 1:**

```
Input: root = [3,9,20,null,null,15,7]
Output: 24
Explanation: There are two left leaves in the binary tree, with values 9 and 15 respectively.
```

**Example 2:**

```
Input: root = [1]
Output: 0
```

**Constraints:**

- The number of nodes in the tree is in the range `[1, 1000]`.
- `-1000 <= Node.val <= 1000`

---

## Company Frequency

| Company    | Frequency        | Last Reported |
|------------|------------------|---------------|
| Amazon     | ★★★☆☆ Medium     | 2024          |
| Adobe      | ★★★☆☆ Medium     | 2023          |
| Microsoft  | ★★☆☆☆ Low        | 2023          |
| Facebook   | ★★☆☆☆ Low        | 2022          |
| Google     | ★★☆☆☆ Low        | 2022          |

> ⚠️ Frequency data is crowd-sourced from LeetCode Discuss, Glassdoor, and
> community interview reports. Treat as a signal, not a guarantee.

---

## DSA Concepts Used

- **Tree Traversal (DFS/BFS)** — the answer is a single full walk of the tree accumulating qualifying leaf values; the traversal skeleton is the whole solution → see [`/dsa/tree_traversal.md`](/dsa/tree_traversal.md)
- **Stack** — the iterative version replaces the recursion stack with an explicit `(node, isLeft)` stack, the standard way to de-recurse a DFS → see [`/dsa/stack.md`](/dsa/stack.md)

---

## Approaches Overview

| # | Approach | Time | Space | When to use |
|---|----------|------|-------|-------------|
| 1 | Recursive DFS (isLeft flag) | O(n) | O(h) | Cleanest expression; pass down "am I a left child?" |
| 2 | Parent-Check DFS | O(n) | O(h) | Avoids the flag by inspecting the left child from the parent |
| 3 | Iterative DFS (explicit stack) | O(n) | O(h) | When recursion is disallowed or stack depth is a concern |

---

## Approach 1 — Recursive DFS

### Intuition

The trap in this problem: a node cannot decide on its own whether it is a "left leaf", because that depends on **which edge its parent used** to reach it. So thread a boolean `isLeft` down the recursion. At each node, if it is a leaf (no children) *and* it was reached via a left edge, contribute its value; otherwise recurse into the left child with `isLeft = true` and the right child with `isLeft = false`. The root is entered with `isLeft = false`, since it has no parent.

### Algorithm

1. `helper(node, isLeft)`: if `node` is `nil`, return `0`.
2. If `node` is a leaf and `isLeft`, return `node.Val`; if a leaf but not a left child, return `0`.
3. Otherwise return `helper(node.Left, true) + helper(node.Right, false)`.
4. Answer = `helper(root, false)`.

### Complexity

- **Time:** O(n) — every node is visited exactly once.
- **Space:** O(h) — recursion depth equals the tree height `h` (up to `n` if degenerate).

### Code

```go
func recursiveDFS(root *TreeNode) int {
	var helper func(node *TreeNode, isLeft bool) int
	helper = func(node *TreeNode, isLeft bool) int {
		if node == nil {
			return 0 // empty subtree adds nothing
		}
		// A leaf reached through a LEFT edge is exactly a "left leaf".
		if node.Left == nil && node.Right == nil {
			if isLeft {
				return node.Val
			}
			return 0 // a leaf, but it was a right child — ignore
		}
		// Recurse: left child is a left edge, right child is not.
		return helper(node.Left, true) + helper(node.Right, false)
	}
	// The root has no parent, so it is treated as "not a left child".
	return helper(root, false)
}
```

### Dry Run

Example 1: `root = [3,9,20,null,null,15,7]`.

| call | node | isLeft | leaf? | contribution / recurse |
|------|------|--------|-------|------------------------|
| helper(3,false) | 3 | false | no | recurse L(9,true) + R(20,false) |
| helper(9,true) | 9 | true | yes | **+9** (left leaf) |
| helper(20,false) | 20 | false | no | recurse L(15,true) + R(7,false) |
| helper(15,true) | 15 | true | yes | **+15** (left leaf) |
| helper(7,false) | 7 | false | yes | +0 (leaf but right child) |

Total = `9 + 15 = 24` ✔

---

## Approach 2 — Parent-Check DFS

### Intuition

Rather than carrying a flag downward, look one level **down** from each node. Standing at a parent, its left child is a "left leaf" exactly when that left child exists and is itself a leaf. Add it on the spot, then continue the traversal into both subtrees. This trades the flag parameter for a one-step lookahead.

### Algorithm

1. `dfs(node)`: if `node` is `nil`, return `0`. Initialise `sum = 0`.
2. If `node.Left != nil`:
   - If `node.Left` is a leaf, `sum += node.Left.Val`.
   - Else `sum += dfs(node.Left)`.
3. `sum += dfs(node.Right)`.
4. Return `sum`; answer = `dfs(root)`.

### Complexity

- **Time:** O(n) — each node is visited once.
- **Space:** O(h) — recursion stack proportional to height.

### Code

```go
func parentCheckDFS(root *TreeNode) int {
	var dfs func(node *TreeNode) int
	dfs = func(node *TreeNode) int {
		if node == nil {
			return 0
		}
		sum := 0
		if node.Left != nil {
			// Is the left child itself a leaf? Then it's a left leaf — count it.
			if node.Left.Left == nil && node.Left.Right == nil {
				sum += node.Left.Val
			} else {
				sum += dfs(node.Left) // otherwise descend into it
			}
		}
		sum += dfs(node.Right) // right subtree may contain its own left leaves
		return sum
	}
	return dfs(root)
}
```

### Dry Run

Example 1: `root = [3,9,20,null,null,15,7]`.

| call | node | left child | left child is leaf? | action | right recurse |
|------|------|-----------|---------------------|--------|---------------|
| dfs(3) | 3 | 9 | yes | **+9** | dfs(20) |
| dfs(20) | 20 | 15 | yes | **+15** | dfs(7) |
| dfs(7) | 7 | nil | — | nothing | dfs(nil)=0 |

Total = `9 + 15 = 24` ✔

---

## Approach 3 — Iterative DFS (Explicit Stack)

### Intuition

Any recursive DFS can be rewritten with an explicit stack, which removes the risk of deep recursion and is often requested in interviews. Store each node together with the same `isLeft` flag the recursive version carried. Pop nodes one at a time; when a popped node is a leaf that arrived via a left edge, add its value.

### Algorithm

1. If `root` is `nil`, return `0`. Push `(root, false)`.
2. While the stack is non-empty, pop `(node, isLeft)`:
   - If `node` is a leaf and `isLeft`, add `node.Val`.
   - Otherwise push `(node.Left, true)` and `(node.Right, false)` when non-nil.
3. Return the accumulated sum.

### Complexity

- **Time:** O(n) — each node is pushed and popped exactly once.
- **Space:** O(h) — the stack holds at most one root-to-leaf path plus branching siblings.

### Code

```go
func iterativeStack(root *TreeNode) int {
	if root == nil {
		return 0
	}
	type frame struct {
		node   *TreeNode
		isLeft bool // whether this node is the left child of its parent
	}
	stack := []frame{{root, false}} // root has no parent -> not a left child
	sum := 0
	for len(stack) > 0 {
		// Pop the top frame.
		top := stack[len(stack)-1]
		stack = stack[:len(stack)-1]
		node, isLeft := top.node, top.isLeft

		if node.Left == nil && node.Right == nil {
			if isLeft {
				sum += node.Val // leaf reached via a left edge = left leaf
			}
			continue // leaves have no children to push
		}
		// Push children with the correct edge label.
		if node.Left != nil {
			stack = append(stack, frame{node.Left, true})
		}
		if node.Right != nil {
			stack = append(stack, frame{node.Right, false})
		}
	}
	return sum
}
```

### Dry Run

Example 1: `root = [3,9,20,null,null,15,7]`. Stack shown top-last.

| step | popped (node, isLeft) | leaf & left? | sum | stack after |
|------|-----------------------|--------------|-----|-------------|
| init | — | — | 0 | `(3,F)` |
| 1 | (3, F) | no | 0 | `(9,T) (20,F)` |
| 2 | (20, F) | no | 0 | `(9,T) (15,T) (7,F)` |
| 3 | (7, F) | leaf, not left | 0 | `(9,T) (15,T)` |
| 4 | (15, T) | leaf & left | **15** | `(9,T)` |
| 5 | (9, T) | leaf & left | **24** | `` (empty) |

Result: `24` ✔

---

## Key Takeaways

- **"Left leaf" is a property of the edge, not the node.** You must know whether you arrived through a left link — carry an `isLeft` flag, or peek at the left child from the parent.
- **Two clean formulations:** push context *down* (flag) or look one step *ahead* (check `node.Left` is a leaf). Both are O(n)/O(h); pick whichever reads clearer to you.
- **De-recursion pattern:** replace the call stack with a stack of `(node, extra-state)` frames — here the extra state is the single boolean. This generalises to any DFS that threads context.
- **Guard the root:** the root can itself be a leaf (Example 2), but it is never a *left* leaf, so it must enter with `isLeft = false`.

---

## Related Problems

- LeetCode #112 — Path Sum (DFS carrying a running total)
- LeetCode #543 — Diameter of Binary Tree (post-order DFS aggregation)
- LeetCode #257 — Binary Tree Paths (DFS carrying path context)
- LeetCode #1022 — Sum of Root To Leaf Binary Numbers (DFS accumulating leaf values)
- LeetCode #104 — Maximum Depth of Binary Tree (DFS/BFS traversal skeleton)
