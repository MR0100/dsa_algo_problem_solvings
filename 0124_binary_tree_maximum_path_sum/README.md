# 0124 — Binary Tree Maximum Path Sum

> LeetCode #124 · Difficulty: Hard
> **Categories:** Dynamic Programming, Tree, Depth-First Search, Binary Tree

---

## Problem Statement

A **path** in a binary tree is a sequence of nodes where each pair of adjacent nodes in the sequence has an edge connecting them. A node can only appear in the sequence **at most once**. Note that the path does not need to pass through the root.

The **path sum** of a path is the sum of the node's values in the path.

Given the `root` of a binary tree, return the **maximum path sum** of any **non-empty** path.

**Example 1:**
```
Input: root = [1,2,3]
Output: 6
Explanation: The path is 2 → 1 → 3 with sum 6.
```

**Example 2:**
```
Input: root = [-10,9,20,null,null,15,7]
Output: 42
Explanation: The path is 15 → 20 → 7 with sum 42.
```

**Constraints:**
- The number of nodes in the tree is in the range `[1, 3 * 10^4]`.
- `-1000 <= Node.val <= 1000`

---

## Company Frequency

| Company   | Frequency       | Last Reported |
|-----------|-----------------|---------------|
| Amazon    | ★★★★★ Very High | 2024          |
| Microsoft | ★★★★★ Very High | 2024          |
| Google    | ★★★★☆ High      | 2024          |
| Facebook  | ★★★★☆ High      | 2024          |
| Bloomberg | ★★★★☆ High      | 2024          |

> ⚠️ Frequency data is crowd-sourced from LeetCode Discuss, Glassdoor, and community interview reports. Treat as a signal, not a guarantee.

---

## DSA Concepts Used

- **Post-order DFS** — compute subtree contribution bottom-up
- **Global maximum with local return** — update global max at each node but return only the single-arm extension

---

## Approaches Overview

| # | Approach        | Time | Space | When to use |
|---|-----------------|------|-------|-------------|
| 1 | Post-Order DFS  | O(n) | O(h)  | Always      |

---

## Approach 1 — Post-Order DFS

### Intuition
For each node, the maximum path **through that node** bends at it: left arm + node + right arm. Both arms are optional (set to 0 if negative).

But when reporting up to the parent, we can only extend **one** arm (the path can't branch), so we return `node.Val + max(leftGain, rightGain)`.

Maintain a global max updated at each node.

### Algorithm
1. `maxPath = -∞`.
2. `gain(node)`:
   - nil → 0.
   - `leftGain  = max(0, gain(node.Left))`.
   - `rightGain = max(0, gain(node.Right))`.
   - `maxPath = max(maxPath, node.Val + leftGain + rightGain)`.
   - Return `node.Val + max(leftGain, rightGain)`.
3. Call `gain(root)`, return `maxPath`.

### Complexity
- **Time:** O(n) — each node visited once.
- **Space:** O(h) — recursion stack.

### Code
```go
func maxPathSum(root *TreeNode) int {
    maxPath := -1<<31
    var gain func(node *TreeNode) int
    gain = func(node *TreeNode) int {
        if node == nil { return 0 }
        leftGain  := gain(node.Left);  if leftGain  < 0 { leftGain  = 0 }
        rightGain := gain(node.Right); if rightGain < 0 { rightGain = 0 }
        if node.Val+leftGain+rightGain > maxPath { maxPath = node.Val+leftGain+rightGain }
        if leftGain > rightGain { return node.Val+leftGain }
        return node.Val+rightGain
    }
    gain(root)
    return maxPath
}
```

### Dry Run
`[-10, 9, 20, null, null, 15, 7]`:

| Node | leftGain | rightGain | pathThrough | return |
|------|----------|-----------|-------------|--------|
| 9    | 0        | 0         | 9           | 9      |
| 15   | 0        | 0         | 15          | 15     |
| 7    | 0        | 0         | 7           | 7      |
| 20   | 15       | 7         | 20+15+7=42  | 20+15=35 |
| -10  | 9→0      | 35→35     | -10+0+35=25 | -10+35=25 |

maxPath = 42 ✓

---

## Key Takeaways
- Key insight: `gain()` returns a single-arm extension (not the full path) because a path can't branch when continuing upward.
- Clamp gains to 0 — never extend with a negative subtree.
- Update global max `before` returning from each node.
- Works even for all-negative trees because we initialize `maxPath = -∞` and a single node is a valid path.

---

## Related Problems
- LeetCode #543 — Diameter of Binary Tree (similar "gain from both sides" trick)
- LeetCode #687 — Longest Univalue Path
- LeetCode #112 — Path Sum (root-to-leaf)
