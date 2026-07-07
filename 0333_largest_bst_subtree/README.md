# 0333 — Largest BST Subtree

> LeetCode #333 · Difficulty: Medium (Premium)
> **Categories:** Tree, Depth-First Search, Binary Search Tree, Dynamic Programming, Binary Tree

---

## Problem Statement

Given the root of a binary tree, find the largest subtree, which is also a Binary Search Tree (BST), where the largest means subtree has the largest number of nodes.

A **Binary Search Tree (BST)** is a tree in which all the nodes follow the below-mentioned properties:

- The left subtree values are less than the value of their parent (root) node's value.
- The right subtree values are greater than the value of their parent (root) node's value.

**Note:** A subtree must include all of its descendants.

**Example 1:**

```
Input: root = [10,5,15,1,8,null,7]
Output: 3
Explanation: The Largest BST Subtree in this case is the highlighted one. The
return value is the subtree's size, which is 3.
```

```
        10
       /  \
      5    15
     / \     \
    1   8     7
```

The subtree rooted at 5 (nodes 5, 1, 8) is a valid BST with 3 nodes.

**Example 2:**

```
Input: root = [4,2,7,2,3,5,null,2,null,null,null,null,null,1]
Output: 2
```

**Constraints:**

- The number of nodes in the tree is in the range `[0, 10^4]`.
- `-10^4 <= Node.val <= 10^4`

**Follow-up:** Can you figure out ways to solve it with `O(n)` time complexity?

---

## Company Frequency

| Company    | Frequency        | Last Reported |
|------------|------------------|---------------|
| Amazon     | ★★★☆☆ Medium     | 2023          |
| Microsoft  | ★★★☆☆ Medium     | 2023          |
| Google     | ★★☆☆☆ Low        | 2023          |
| Bloomberg  | ★★☆☆☆ Low        | 2022          |
| ByteDance  | ★★☆☆☆ Low        | 2022          |

> ⚠️ Frequency data is crowd-sourced from LeetCode Discuss, Glassdoor, and
> community interview reports. Treat as a signal, not a guarantee.

---

## DSA Concepts Used

- **Binary Search Tree** — the core validity rule (left < node < right, applied to the whole subtree) is what each approach checks → see [`/dsa/binary_search_tree.md`](/dsa/binary_search_tree.md)
- **Tree Traversal (Post-Order)** — the optimal solution collects `(isBST, size, min, max)` from children before deciding the parent, i.e. bottom-up DFS → see [`/dsa/tree_traversal.md`](/dsa/tree_traversal.md)
- **Dynamic Programming on Trees** — post-order returns a small tuple of subproblem answers per node so the parent decides in O(1), the classic tree-DP shape → see [`/dsa/tree_dp.md`](/dsa/tree_dp.md)

---

## Approaches Overview

| # | Approach | Time | Space | When to use |
|---|----------|------|-------|-------------|
| 1 | Brute Force (Validate Every Subtree) | O(n²) | O(h) | Baseline; clear but re-validates descendants repeatedly |
| 2 | Post-Order Bottom-Up (Optimal) | O(n) | O(h) | The follow-up answer; one pass, O(1) work per node |

*(n = node count, h = tree height.)*

---

## Approach 1 — Brute Force (Validate Every Subtree)

### Intuition

Read the problem literally: for each node, ask two questions — "is the subtree rooted here a *valid* BST?" and "how many nodes does it contain?" If it is a BST, its node count is a candidate answer; take the maximum over all nodes. Validity is checked the standard way, propagating an allowed value range `(lo, hi)` downward. The waste is obvious: validating a big subtree re-touches every descendant, and we do that at every node.

### Algorithm

1. Define `bruteForce(node)`:
   - If `node == nil`, return 0.
   - If `isBST(node, -inf, +inf)`, the whole subtree is a BST → return `countNodes(node)`.
   - Otherwise the best BST is entirely inside a child → return `max(bruteForce(left), bruteForce(right))`.
2. `isBST(node, lo, hi)` verifies `lo < node.Val < hi` and recurses with tightened bounds.
3. `countNodes` returns the subtree size.

### Complexity

- **Time:** O(n²) — a skewed tree makes each `isBST`/`countNodes` O(n), invoked at O(n) nodes.
- **Space:** O(h) — recursion depth.

### Code

```go
func bruteForce(root *TreeNode) int {
	if root == nil {
		return 0
	}
	// If this whole subtree is a BST, its node count is a candidate.
	if isBST(root, -1<<62, 1<<62) {
		return countNodes(root)
	}
	// Otherwise the best BST lies within one of the children's subtrees.
	l := bruteForce(root.Left)
	r := bruteForce(root.Right)
	if l > r {
		return l
	}
	return r
}

func isBST(node *TreeNode, lo, hi int) bool {
	if node == nil {
		return true // empty subtree is trivially a BST
	}
	if node.Val <= lo || node.Val >= hi {
		return false // value violates the inherited bound
	}
	return isBST(node.Left, lo, node.Val) && isBST(node.Right, node.Val, hi)
}

func countNodes(node *TreeNode) int {
	if node == nil {
		return 0
	}
	return 1 + countNodes(node.Left) + countNodes(node.Right)
}
```

### Dry Run

Example 1: `root = [10,5,15,1,8,null,7]`.

| Node | isBST(node, −∞, +∞)? | Reason | Candidate size |
|------|----------------------|--------|----------------|
| 10 | no | right subtree has 7 < 10 but 7 sits under 15 (7 < 10 breaks the BST rule at root) | — |
| 5 | **yes** | 1 < 5 < 8, valid | `countNodes = 3` |
| 15 | no | 7 is a right child but 7 < 15 | — |
| 1 | yes | leaf | 1 |
| 8 | yes | leaf | 1 |
| 7 | yes | leaf | 1 |

Best candidate = 3 (subtree at 5). Result: `3` ✔

---

## Approach 2 — Post-Order Bottom-Up (Optimal)

### Intuition

Turn the O(n²) re-validation into O(n) by carrying results **upward**. A subtree at `node` is a BST **iff** its left and right subtrees are both BSTs *and* `leftMax < node.Val < rightMin`. So if every child hands back `(isBST, size, min, max)`, the parent decides in O(1) and combines the sizes. Track a global best size. The one subtlety is the `nil` base case: return an "empty BST" with an inverted range (`min = +∞, max = −∞`) so the checks `leftMax < val` and `val < rightMin` pass vacuously for a missing child.

### Algorithm

1. Define post-order `dfs(node)` returning `info{isBST, size, min, max}`:
   - `nil` → `{true, 0, +∞, −∞}` (empty, inverted range).
   - Compute `l = dfs(left)`, `r = dfs(right)`.
   - If `l.isBST && r.isBST && l.max < node.Val < r.min`: this subtree is a BST of `size = l.size + r.size + 1`; update `best`; return `{true, size, min(node.Val, l.min), max(node.Val, r.max)}`.
   - Else return `{isBST: false}`.
2. Return `best`.

### Complexity

- **Time:** O(n) — every node visited exactly once with O(1) combine work.
- **Space:** O(h) — recursion stack.

### Code

```go
func postOrder(root *TreeNode) int {
	best := 0
	type info struct {
		isBST    bool
		size     int
		min, max int
	}
	var dfs func(node *TreeNode) info
	dfs = func(node *TreeNode) info {
		if node == nil {
			return info{isBST: true, size: 0, min: 1 << 62, max: -1 << 62}
		}
		l := dfs(node.Left)  // info about left subtree
		r := dfs(node.Right) // info about right subtree
		if l.isBST && r.isBST && l.max < node.Val && node.Val < r.min {
			sz := l.size + r.size + 1 // combined size including this node
			if sz > best {
				best = sz // record a larger valid BST
			}
			return info{
				isBST: true,
				size:  sz,
				min:   min(node.Val, l.min), // smallest value in this BST
				max:   max(node.Val, r.max), // largest value in this BST
			}
		}
		return info{isBST: false}
	}
	dfs(root)
	return best
}
```

### Dry Run

Example 1: `root = [10,5,15,1,8,null,7]`. Post-order visits leaves first.

| Node | l.info | r.info | BST here? | size | returns |
|------|--------|--------|-----------|------|---------|
| 1 | {T,0,+∞,−∞} | {T,0,+∞,−∞} | yes | 1 | {T,1,1,1}, best=1 |
| 8 | {T,0,+∞,−∞} | {T,0,+∞,−∞} | yes | 1 | {T,1,8,8}, best=1 |
| 5 | {T,1,1,1} | {T,1,8,8} | 1<5<8 → yes | 3 | {T,3,1,8}, **best=3** |
| 7 | {T,0,+∞,−∞} | {T,0,+∞,−∞} | yes | 1 | {T,1,7,7} |
| 15 | {T,0,+∞,−∞} | {T,1,7,7} | need val<r.min: 15<7 false | — | {F} |
| 10 | {T,3,1,8} | {F} | r not BST | — | {F} |

`best = 3`. Result: `3` ✔

---

## Key Takeaways

- **Bottom-up beats top-down when validity depends on aggregated child facts.** Returning a `(isBST, size, min, max)` tuple lets each parent decide in O(1) — the recurring trick for tree problems that would otherwise re-scan subtrees (compare LeetCode #124, #543, #110).
- **The inverted sentinel range** (`min = +∞, max = −∞`) for the empty subtree is what makes a node with one missing child validate cleanly without special-casing.
- **BST validity is a *range* property, not a *local* one:** it is not enough that `left.val < node < right.val`; every value in the left subtree must be < node. That's exactly why the brute force propagates `(lo, hi)` and the optimal carries up `(min, max)`.
- "A subtree must include all descendants" — you cannot cherry-pick nodes; this is what separates the problem from generic subset searches.

---

## Related Problems

- LeetCode #98 — Validate Binary Search Tree (the validity check, standalone)
- LeetCode #124 — Binary Tree Maximum Path Sum (post-order returns a value + updates a global best)
- LeetCode #110 — Balanced Binary Tree (bottom-up height + validity)
- LeetCode #543 — Diameter of Binary Tree (post-order aggregate with global max)
- LeetCode #250 — Count Univalue Subtrees (same "is this subtree special + propagate" pattern)
