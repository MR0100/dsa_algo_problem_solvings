# 0337 — House Robber III

> LeetCode #337 · Difficulty: Medium
> **Categories:** Tree, Depth-First Search, Dynamic Programming, Binary Tree

---

## Problem Statement

The thief has found himself a new place for his thievery again. There is only one entrance to this area, called `root`.

Besides the `root`, each house has one and only one parent house. After a tour, the smart thief realized that all houses in this place form a binary tree. It will automatically contact the police if **two directly-linked houses were broken into on the same night**.

Given the `root` of the binary tree, return *the maximum amount of money the thief can rob **without alerting the police***.

**Example 1:**

```
Input: root = [3,2,3,null,3,null,1]
Output: 7
Explanation: Maximum amount of money the thief can rob = 3 + 3 + 1 = 7.
```

**Example 2:**

```
Input: root = [3,4,5,1,3,null,1]
Output: 9
Explanation: Maximum amount of money the thief can rob = 4 + 5 = 9.
```

**Constraints:**

- The number of nodes in the tree is in the range `[1, 10^4]`.
- `0 <= Node.val <= 10^4`

---

## Company Frequency

| Company    | Frequency        | Last Reported |
|------------|------------------|---------------|
| Amazon     | ★★★★☆ High       | 2024          |
| Google     | ★★★☆☆ Medium     | 2024          |
| Microsoft  | ★★★☆☆ Medium     | 2023          |
| Meta       | ★★★☆☆ Medium     | 2023          |
| Uber       | ★★☆☆☆ Low        | 2022          |

> ⚠️ Frequency data is crowd-sourced from LeetCode Discuss, Glassdoor, and
> community interview reports. Treat as a signal, not a guarantee.

---

## DSA Concepts Used

- **Tree Traversal (post-order DFS)** — the optimal solution solves children before parents so a node can combine their results → see [`/dsa/tree_traversal.md`](/dsa/tree_traversal.md)
- **Dynamic Programming (1D state per node)** — each node keeps two states (robbed / not robbed); this is tree DP, the closest existing file is → see [`/dsa/dynamic_programming_1d.md`](/dsa/dynamic_programming_1d.md)
- **Graph BFS/DFS** — depth-first recursion over the tree structure → see [`/dsa/graph_bfs_dfs.md`](/dsa/graph_bfs_dfs.md)

---

## Approaches Overview

| # | Approach | Time | Space | When to use |
|---|----------|------|-------|-------------|
| 1 | Naive Recursion | O(2^n) | O(h) | Explains the rob/skip choice; too slow for large trees |
| 2 | Memoized Recursion (Top-Down DP) | O(n) | O(n) | Fixes the recomputation while keeping the same recurrence |
| 3 | DFS Returning a Pair (Optimal) | O(n) | O(h) | The clean interview answer; one pass, O(h) space |

---

## Approach 1 — Naive Recursion

### Intuition

At any node the thief either **robs it** — then its children are forbidden, so the next legal houses are the four grandchildren — or **skips it**, leaving the children free to rob. Take the max of the two. Correct, but robbing recurses into grandchildren while skipping recurses into children which again recurse into those same grandchildren, so subtrees are recomputed exponentially.

### Algorithm

1. If `root` is `nil`, return `0`.
2. `robThis = root.Val + rob(all four grandchildren)`.
3. `skipThis = rob(root.Left) + rob(root.Right)`.
4. Return `max(robThis, skipThis)`.

### Complexity

- **Time:** O(2^n) — the same subtrees are re-solved through both the rob and skip branches.
- **Space:** O(h) — recursion stack, `h` = tree height.

### Code

```go
func naiveRecursion(root *TreeNode) int {
	if root == nil {
		return 0 // empty subtree yields nothing
	}
	robThis := root.Val // we rob this house...
	if root.Left != nil {
		// ...so its children are off-limits; jump to grandchildren.
		robThis += naiveRecursion(root.Left.Left) + naiveRecursion(root.Left.Right)
	}
	if root.Right != nil {
		robThis += naiveRecursion(root.Right.Left) + naiveRecursion(root.Right.Right)
	}
	// Or skip this house and rob its children freely.
	skipThis := naiveRecursion(root.Left) + naiveRecursion(root.Right)
	return max(robThis, skipThis) // best of the two choices
}
```

### Dry Run

Example 1: `root = [3,2,3,null,3,null,1]`. Tree:

```
        3
       / \
      2   3
       \    \
        3    1
```

| Call | robThis | skipThis | returns |
|------|---------|----------|---------|
| node 2 (left child, has right child 3) | 2 + 0 (grandkids of 2 = none) = 2 | rob(3)=3 | max(2,3)=3 |
| node 3 (right child, has right child 1) | 3 + 0 = 3 | rob(1)=1 | max(3,1)=3 |
| root 3 | 3 + [3's grandkids: 3 and 1] = 3+3+1 = 7 | rob(2)+rob(3) = 3+3 = 6 | max(7,6)=7 |

Result: `7` ✔

---

## Approach 2 — Memoized Recursion (Top-Down DP)

### Intuition

The naive recurrence is right; only the repetition hurts. Cache each node's best result keyed by its pointer. The first evaluation stores the answer; every later request for that subtree is an O(1) lookup, so each node is solved exactly once.

### Algorithm

1. Keep `memo: *TreeNode → int`.
2. `dfs(node)`: if `nil` return 0; if cached return it.
3. Otherwise compute `robThis`/`skipThis` as in Approach 1, store `max` in `memo`, return it.

### Complexity

- **Time:** O(n) — each node computed once with O(1) combining work.
- **Space:** O(n) memo + O(h) recursion stack.

### Code

```go
func memoized(root *TreeNode) int {
	memo := map[*TreeNode]int{} // node → best amount for its subtree
	var dfs func(node *TreeNode) int
	dfs = func(node *TreeNode) int {
		if node == nil {
			return 0
		}
		if v, ok := memo[node]; ok {
			return v // already solved this subtree
		}
		robThis := node.Val
		if node.Left != nil {
			robThis += dfs(node.Left.Left) + dfs(node.Left.Right)
		}
		if node.Right != nil {
			robThis += dfs(node.Right.Left) + dfs(node.Right.Right)
		}
		skipThis := dfs(node.Left) + dfs(node.Right)
		best := max(robThis, skipThis)
		memo[node] = best // cache before returning
		return best
	}
	return dfs(root)
}
```

### Dry Run

Example 1 (same tree as above). Nodes are memoized bottom-up:

| Node | robThis | skipThis | memo stored |
|------|---------|----------|-------------|
| leaf 3 (under node 2) | 3 | 0 | 3 |
| leaf 1 (under node 3) | 1 | 0 | 1 |
| node 2 | 2 | memo[3]=3 | 3 |
| node 3 (right) | 3 | memo[1]=1 | 3 |
| root 3 | 3 + memo[3] + memo[1] = 3+3+1 = 7 | memo[2]+memo[3] = 3+3 = 6 | 7 |

Result: `7` ✔ (each subtree computed once).

---

## Approach 3 — DFS Returning a Pair (Optimal)

### Intuition

Instead of peeking at grandchildren, have each node return **two** values to its parent: the best if this node is robbed, and the best if it is not. The parent then decides using only its children's pairs:

- rob node ⇒ children must be skipped ⇒ `node.Val + left.notRob + right.notRob`
- skip node ⇒ each child takes its own best ⇒ `max(left) + max(right)`

One post-order pass, no repeated work, O(h) space.

### Algorithm

1. `dfs(node)` returns `(rob, notRob)`.
2. Base: `nil → (0, 0)`.
3. Solve `left` and `right` first; then
   `rob = node.Val + lNot + rNot`,
   `notRob = max(lRob, lNot) + max(rRob, rNot)`.
4. Answer = `max(rootRob, rootNotRob)`.

### Complexity

- **Time:** O(n) — one visit per node, O(1) combining.
- **Space:** O(h) — recursion stack only; no memo table.

### Code

```go
func robTreeDP(root *TreeNode) int {
	// dfs returns {rob = best including node, notRob = best excluding node}.
	var dfs func(node *TreeNode) (int, int)
	dfs = func(node *TreeNode) (int, int) {
		if node == nil {
			return 0, 0 // nothing to rob, nothing to skip
		}
		lRob, lNot := dfs(node.Left)  // children solved first (post-order)
		rRob, rNot := dfs(node.Right) //
		// If we rob this node, both children must be skipped.
		rob := node.Val + lNot + rNot
		// If we skip this node, each child independently takes its own best.
		notRob := max(lRob, lNot) + max(rRob, rNot)
		return rob, notRob
	}
	rob, notRob := dfs(root)
	return max(rob, notRob) // whole-tree best
}
```

### Dry Run

Example 1 (same tree). Each node returns `(rob, notRob)`:

| Node | lRob,lNot | rRob,rNot | rob = val+lNot+rNot | notRob = max(l)+max(r) | pair |
|------|-----------|-----------|----------------------|-------------------------|------|
| leaf 3 (under 2) | 0,0 | 0,0 | 3+0+0 = 3 | 0+0 = 0 | (3,0) |
| leaf 1 (under 3) | 0,0 | 0,0 | 1+0+0 = 1 | 0 | (1,0) |
| node 2 | (nil)=0,0 | (3,0) | 2+0+0 = 2 | max(0,0)+max(3,0) = 3 | (2,3) |
| node 3 (right) | (nil)=0,0 | (1,0) | 3+0+0 = 3 | max(0,0)+max(1,0) = 1 | (3,1) |
| root 3 | (2,3) | (3,1) | 3+3+1 = 7 | max(2,3)+max(3,1) = 3+3 = 6 | (7,6) |

Answer = `max(7, 6) = 7` ✔

---

## Key Takeaways

- **Tree DP = return per-node states bottom-up.** When a node's optimum depends on whether it is "used", return a tuple `(used, notUsed)` so the parent combines in O(1).
- The **rob/skip dichotomy** is the same as linear House Robber (#198); the tree just replaces "previous index" with "children".
- Returning a pair beats memoization: it removes the grandchild double-recursion *and* the O(n) memo map, landing at O(n) time / O(h) space.
- Post-order traversal is the natural fit whenever a parent needs fully-computed child results.

---

## Related Problems

- LeetCode #198 — House Robber (linear version, same rob/skip DP)
- LeetCode #213 — House Robber II (circular array)
- LeetCode #124 — Binary Tree Maximum Path Sum (post-order returns per-node value)
- LeetCode #543 — Diameter of Binary Tree (post-order aggregation)
