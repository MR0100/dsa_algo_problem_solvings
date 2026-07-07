# 0250 — Count Univalue Subtrees

> LeetCode #250 · Difficulty: Medium
> **Categories:** Tree, Depth-First Search, Binary Tree

---

## Problem Statement

Given the `root` of a binary tree, return *the number of **uni-value** subtrees*.

A **uni-value subtree** means all nodes of the subtree have the same value.

**Example 1:**

```
Input: root = [5,1,5,5,5,null,5]
Output: 4
```

**Example 2:**

```
Input: root = []
Output: 0
```

**Example 3:**

```
Input: root = [5,5,5,5,5,null,5]
Output: 6
```

**Constraints:**

- The number of the node in the tree will be in the range `[0, 1000]`.
- `-1000 <= Node.val <= 1000`

---

## Company Frequency

| Company    | Frequency        | Last Reported |
|------------|------------------|---------------|
| Google     | ★★★☆☆ Medium     | 2024          |
| Amazon     | ★★★☆☆ Medium     | 2023          |
| Facebook   | ★★☆☆☆ Low        | 2023          |
| Microsoft  | ★★☆☆☆ Low        | 2022          |

> ⚠️ Frequency data is crowd-sourced from LeetCode Discuss, Glassdoor, and
> community interview reports. Treat as a signal, not a guarantee.

---

## DSA Concepts Used

- **Tree Traversal (Post-Order DFS)** — a node's univalue status depends on its children, so resolve children first, then decide the node → see [`/dsa/tree_traversal.md`](/dsa/tree_traversal.md)
- **Graph BFS/DFS** — recursion over a tree is depth-first traversal; the counter is accumulated on the way back up → see [`/dsa/graph_bfs_dfs.md`](/dsa/graph_bfs_dfs.md)

---

## Approaches Overview

| # | Approach | Time | Space | When to use |
|---|----------|------|-------|-------------|
| 1 | Post-Order DFS with Return Flag (Optimal) | O(n) | O(h) | Canonical: return bool up, bump a shared counter |
| 2 | Post-Order Returning (isUnival, count) Pair | O(n) | O(h) | Functional style; no captured mutable counter |

---

## Approach 1 — Post-Order DFS with Return Flag (Optimal)

### Intuition

A subtree is univalue iff every node in it shares one value. That is a bottom-up property: a leaf is trivially univalue; an internal node is univalue iff **both** child subtrees are univalue **and** every present child's value equals the node's own value. Post-order (children before parent) lets each node combine its already-computed children's results. We return "is this subtree univalue" up the stack and increment a shared counter whenever the answer is yes.

### Algorithm

1. `dfs(node)` returns whether `node`'s subtree is univalue.
2. `nil → true` (a missing child imposes no constraint).
3. Recurse both children first (post-order): `leftUni`, `rightUni`.
4. If either child subtree isn't univalue, return `false`.
5. If any present child's value differs from `node.Val`, return `false`.
6. Otherwise `count++` and return `true`.

### Complexity

- **Time:** O(n) — every node visited exactly once.
- **Space:** O(h) — recursion stack, `h` = height (O(n) skewed, O(log n) balanced).

### Code

```go
func countUnivalPostOrder(root *TreeNode) int {
	count := 0

	var dfs func(node *TreeNode) bool
	dfs = func(node *TreeNode) bool {
		if node == nil {
			return true
		}

		leftUni := dfs(node.Left)
		rightUni := dfs(node.Right)

		if !leftUni || !rightUni {
			return false
		}
		if node.Left != nil && node.Left.Val != node.Val {
			return false
		}
		if node.Right != nil && node.Right.Val != node.Val {
			return false
		}

		count++
		return true
	}

	dfs(root)
	return count
}
```

### Dry Run

Tree `[5,1,5,5,5,null,5]`:

```
        5(a)
       /    \
    1(b)     5(c)
    /  \        \
 5(d)  5(e)     5(f)
```

Post-order visit order: d, e, b, f, c, a.

| Node | children univalue? | value match? | univalue? | count |
|------|--------------------|--------------|-----------|-------|
| d (5, leaf) | — | — | yes | 1 |
| e (5, leaf) | — | — | yes | 2 |
| b (1) | d,e univalue | d.Val=5 ≠ b.Val=1 | **no** | 2 |
| f (5, leaf) | — | — | yes | 3 |
| c (5) | f univalue | f.Val=5 == c.Val=5 | yes | 4 |
| a (5) | b **not** univalue | — | **no** | 4 |

Result = `4`. ✓

---

## Approach 2 — Post-Order Returning (isUnival, count) Pair

### Intuition

Same post-order logic, but instead of a shared closure counter, each call returns a `(isUnival, count)` pair: whether the subtree is univalue, and the total number of univalue subtrees within it. The parent sums both children's counts, decides its own univalue status, and adds 1 if applicable. No mutable captured state.

### Algorithm

1. `dfs(node)` returns `(isUnival bool, count int)`.
2. `nil → (true, 0)`.
3. Recurse both children; `total = leftCount + rightCount`.
4. `isUnival = leftUni && rightUni`, then false if any present child's value differs.
5. If `isUnival`, `total++`.
6. Return `(isUnival, total)`.

### Complexity

- **Time:** O(n).
- **Space:** O(h) recursion depth.

### Code

```go
func countUnivalPair(root *TreeNode) int {
	var dfs func(node *TreeNode) (bool, int)
	dfs = func(node *TreeNode) (bool, int) {
		if node == nil {
			return true, 0
		}

		leftUni, leftCount := dfs(node.Left)
		rightUni, rightCount := dfs(node.Right)
		total := leftCount + rightCount

		isUni := leftUni && rightUni
		if isUni && node.Left != nil && node.Left.Val != node.Val {
			isUni = false
		}
		if isUni && node.Right != nil && node.Right.Val != node.Val {
			isUni = false
		}

		if isUni {
			total++
		}
		return isUni, total
	}

	_, count := dfs(root)
	return count
}
```

### Dry Run

Same tree `[5,1,5,5,5,null,5]`, returning `(isUni, count)`:

| Node | left ret | right ret | total before | isUni | total after |
|------|----------|-----------|--------------|-------|-------------|
| d (5) | (true,0) | (true,0) | 0 | true | 1 |
| e (5) | (true,0) | (true,0) | 0 | true | 1 |
| b (1) | (true,1) from d | (true,1) from e | 2 | false (d.Val≠1) | 2 |
| f (5) | (true,0) | (true,0) | 0 | true | 1 |
| c (5) | (true,0) | (true,1) from f | 1 | true | 2 |
| a (5) | (false,2) from b | (true,2) from c | 4 | false | 4 |

Returns `(false, 4)` → answer `4`. ✓

---

## Key Takeaways

- Whenever a node's answer depends on its children, reach for **post-order DFS**: recurse first, combine on the way back up.
- Returning a small tuple `(property, aggregate)` avoids a mutable closure counter and keeps the recursion pure — a clean interview pattern.
- Treat `nil` children as the identity (univalue = true, count = 0) so leaf and one-child cases fall out of the general rule.

---

## Related Problems

- LeetCode #965 — Univalued Binary Tree (whole-tree version)
- LeetCode #543 — Diameter of Binary Tree (post-order returning a value)
- LeetCode #687 — Longest Univalue Path (same-value path via post-order)
- LeetCode #366 — Find Leaves of Binary Tree (bottom-up aggregation)
