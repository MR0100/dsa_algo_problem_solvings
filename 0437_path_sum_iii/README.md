# 0437 — Path Sum III

> LeetCode #437 · Difficulty: Medium
> **Categories:** Tree, Depth-First Search, Prefix Sum, Hash Table, Binary Tree

---

## Problem Statement

Given the `root` of a binary tree and an integer `targetSum`, return *the number of paths where the sum of the values along the path equals* `targetSum`.

The path does not need to start or end at the root or a leaf, but it must go **downwards** (i.e., traveling only from parent nodes to child nodes).

**Example 1:**

```
Input: root = [10,5,-3,3,2,null,11,3,-2,null,1], targetSum = 8
Output: 3
Explanation: The paths that sum to 8 are shown.
```

The three paths are `5 → 3`, `5 → 2 → 1`, and `-3 → 11`.

**Example 2:**

```
Input: root = [5,4,8,11,null,13,4,7,2,null,null,5,1], targetSum = 22
Output: 3
```

**Constraints:**

- The number of nodes in the tree is in the range `[0, 1000]`.
- `-10^9 <= Node.val <= 10^9`
- `-1000 <= targetSum <= 1000`

---

## Company Frequency

| Company    | Frequency        | Last Reported |
|------------|------------------|---------------|
| Amazon     | ★★★★☆ High       | 2024          |
| Facebook   | ★★★★☆ High       | 2024          |
| Google     | ★★★☆☆ Medium     | 2023          |
| Microsoft  | ★★★☆☆ Medium     | 2023          |
| Bloomberg  | ★★☆☆☆ Low        | 2022          |
| Apple      | ★★☆☆☆ Low        | 2022          |

> ⚠️ Frequency data is crowd-sourced from LeetCode Discuss, Glassdoor, and
> community interview reports. Treat as a signal, not a guarantee.

---

## DSA Concepts Used

- **Prefix Sum** — the values along a root-to-node chain form a prefix sum; a sub-path sums to target exactly when two prefixes differ by target — the tree version of "subarray sum equals k" → see [`/dsa/prefix_sum.md`](/dsa/prefix_sum.md)
- **Hash Map** — a frequency map of prefix sums seen on the current path turns "how many ancestors close a target path?" into an O(1) lookup → see [`/dsa/hash_map.md`](/dsa/hash_map.md)
- **Tree Traversal (DFS)** — both solutions are depth-first walks; the optimal one threads a running sum down and backtracks the map on the way up → see [`/dsa/tree_traversal.md`](/dsa/tree_traversal.md)

---

## Approaches Overview

| # | Approach | Time | Space | When to use |
|---|----------|------|-------|-------------|
| 1 | Brute Force (Double DFS) | O(n²) worst / O(n log n) balanced | O(h) | Small trees, or to build intuition; anchors a path start at every node |
| 2 | Prefix Sum + Hash Map (Optimal) | O(n) | O(n) | The interview answer; single pass with a backtracked prefix-sum map |

---

## Approach 1 — Brute Force (Double DFS)

### Intuition

Any valid path is a downward chain. If we fix its **top** node, counting paths reduces to "from this node, how many downward chains sum to target?" — a plain recursion that subtracts each node's value as it descends. Anchoring the top at every node in turn covers all paths. It repeats work on shared suffixes (a deep node is re-walked once per ancestor anchor), which is the source of the quadratic worst case.

### Algorithm

1. `countFrom(node, remaining)`: if `node == nil` return 0; let `rem = remaining - node.Val`; add 1 if `rem == 0`; recurse into both children with `rem` and sum.
2. `bruteForce(root, target)`: return `countFrom(root, target)` plus `bruteForce` on each child (which re-anchors the path start deeper).
3. Do **not** stop when `rem == 0` — later negative values can still form longer valid paths.

### Complexity

- **Time:** O(n²) for a skewed tree (n anchors × O(n) suffix each); O(n log n) when balanced.
- **Space:** O(h) — recursion stack, `h` = height.

### Code

```go
func bruteForce(root *TreeNode, targetSum int) int {
	if root == nil {
		return 0
	}
	// Paths starting exactly at root, plus all paths that start deeper (handled
	// by recursing the OUTER walk into each child with the full targetSum).
	return countFrom(root, targetSum) +
		bruteForce(root.Left, targetSum) +
		bruteForce(root.Right, targetSum)
}

func countFrom(node *TreeNode, remaining int) int {
	if node == nil {
		return 0
	}
	rem := remaining - node.Val // consume this node's value along the path
	count := 0
	if rem == 0 {
		count = 1 // a path ending here (root-of-this-call → node) hits the target
	}
	count += countFrom(node.Left, rem)
	count += countFrom(node.Right, rem)
	return count
}
```

### Dry Run

Example 1, `targetSum = 8`. Tree (partial):

```
        10
       /  \
      5    -3
     / \     \
    3   2     11
   / \    \
  3  -2    1
```

Anchoring `countFrom` at each node (only non-zero anchors shown):

| anchor node | downward chains found summing to 8 | contribution |
|-------------|-------------------------------------|--------------|
| 10 | 10→5→(-3)? no… no chain from 10 sums to 8 | 0 |
| 5  | `5→3` (=8 ✓); `5→2→1` (=8 ✓) | 2 |
| -3 | `-3→11` (=8 ✓) | 1 |
| others | none | 0 |

Total = 2 + 1 = **3** ✔

---

## Approach 2 — Prefix Sum + Hash Map (Optimal)

### Intuition

Walk once from the root, keeping `curr` = sum of values from the root down to the current node — a prefix sum. A sub-path from some ancestor `A` (exclusive) down to the current node sums to `target` precisely when `curr - prefix(A) == target`, i.e. `prefix(A) == curr - target`. So maintain a frequency map of prefix sums seen **on the current root-to-node chain**; the number of paths ending at this node is `freq[curr - target]`. Register `curr` before descending and **decrement it on the way back up** (backtrack), so the map never mixes in sibling branches. Seeding `freq[0] = 1` accounts for full paths that start at the root itself.

### Algorithm

1. Initialise `freq = {0: 1}` (empty prefix), `total = 0`.
2. `dfs(node, curr)`: if nil return; `curr += node.Val`; `total += freq[curr - target]`; `freq[curr]++`; recurse into both children; then `freq[curr]--`.
3. Return `total`.

### Complexity

- **Time:** O(n) — each node visited once with O(1) map operations.
- **Space:** O(n) — the prefix-sum map (up to one entry per path node) plus O(h) recursion.

### Code

```go
func prefixSumHashMap(root *TreeNode, targetSum int) int {
	freq := map[int]int{0: 1} // prefix sum 0 seen once: the empty prefix at the root
	total := 0

	var dfs func(node *TreeNode, curr int)
	dfs = func(node *TreeNode, curr int) {
		if node == nil {
			return
		}
		curr += node.Val // running sum from root down to (and including) node
		// Any ancestor prefix equal to (curr - target) closes a path summing to
		// target ending at this node; count all such ancestors.
		total += freq[curr-targetSum]
		freq[curr]++ // register this node's prefix for its descendants
		dfs(node.Left, curr)
		dfs(node.Right, curr)
		freq[curr]-- // backtrack: leave the map holding only the current chain
	}

	dfs(root, 0)
	return total
}
```

### Dry Run

Example 1, `targetSum = 8`. Trace the branch `10 → 5 → 3 → 3` then the relevant others. `freq` starts `{0:1}`.

| visit node | curr | look up `freq[curr-8]` | total | freq after `freq[curr]++` |
|------------|------|------------------------|-------|----------------------------|
| 10 | 10 | freq[2] = 0 | 0 | {0:1, 10:1} |
| 5 (left) | 15 | freq[7] = 0 | 0 | {0:1, 10:1, 15:1} |
| 3 (5.left) | 18 | freq[10] = **1** (ancestor 10) → path `5→3` | 1 | {…, 18:1} |
| 3 (3.left) | 21 | freq[13] = 0 | 1 | {…, 21:1} |
| backtrack 3, 3 … then 2 (5.right) | 17 | freq[9] = 0 | 1 | {…, 17:1} |
| 1 (2.right) | 18 | freq[10] = **1** (ancestor 10) → path `5→2→1` | 2 | {…} |
| … right subtree 11 via -3 | -3: curr 7; 11: curr 18 | at 11, freq[10] = **1** → path `-3→11` | 3 | — |

Total = **3** ✔ (backtracking removes each branch's prefixes so, e.g., the `5→3→3` prefixes do not leak into the `-3→11` branch).

---

## Key Takeaways

- **Path sum on a tree = subarray sum on a line.** The root-to-node chain is an array; "count sub-paths equal to k" is [LeetCode #560](https://leetcode.com/problems/subarray-sum-equals-k/) lifted onto a tree. Same `freq[curr - k]` lookup.
- **Backtrack the map, not just the recursion.** Decrementing `freq[curr]` when returning is what keeps a sibling branch from seeing the current branch's prefixes. Forgetting it is the #1 bug here.
- **Seed `freq[0] = 1`** so a path that begins at the root (no ancestor to subtract) is counted.
- **Don't early-exit on a zero remainder** in the brute force — negative node values mean a longer path can still be valid.

---

## Related Problems

- LeetCode #112 — Path Sum (does *a* root-to-leaf path equal target?)
- LeetCode #113 — Path Sum II (enumerate the root-to-leaf paths)
- LeetCode #560 — Subarray Sum Equals K (the 1-D prefix-sum-map original)
- LeetCode #124 — Binary Tree Maximum Path Sum (arbitrary path, DP variant)
- LeetCode #666 — Path Sum IV (path sums over an encoded tree)
