# 0112 — Path Sum

> LeetCode #112 · Difficulty: Easy
> **Categories:** Tree, Depth-First Search, Breadth-First Search, Binary Tree

---

## Problem Statement

Given the `root` of a binary tree and an integer `targetSum`, return `true` if the tree has a root-to-leaf path such that adding up all the values along the path equals `targetSum`.

A **leaf** is a node with no children.

**Example 1:**
```
Input: root = [5,4,8,11,null,13,4,7,2,null,null,null,1], targetSum = 22
Output: true
Explanation: The path 5 → 4 → 11 → 2 sums to 22.
```

**Example 2:**
```
Input: root = [1,2,3], targetSum = 5
Output: false
```

**Example 3:**
```
Input: root = [], targetSum = 0
Output: false
```

**Constraints:**
- The number of nodes in the tree is in the range `[0, 5000]`.
- `-1000 <= Node.val <= 1000`
- `-1000 <= targetSum <= 1000`

---

## Company Frequency

| Company   | Frequency       | Last Reported |
|-----------|-----------------|---------------|
| Amazon    | ★★★★★ Very High | 2024          |
| Microsoft | ★★★★☆ High      | 2024          |
| Google    | ★★★☆☆ Medium    | 2023          |
| Bloomberg | ★★★☆☆ Medium    | 2023          |

> ⚠️ Frequency data is crowd-sourced from LeetCode Discuss, Glassdoor, and community interview reports. Treat as a signal, not a guarantee.

---

## DSA Concepts Used

- **Tree DFS** — subtract current node value, check at leaf

---

## Approaches Overview

| # | Approach        | Time | Space | When to use           |
|---|-----------------|------|-------|-----------------------|
| 1 | Recursive DFS   | O(n) | O(h)  | Clean and concise     |
| 2 | Iterative DFS   | O(n) | O(h)  | Avoids recursion      |

---

## Approach 1 — Recursive DFS

### Intuition
Subtract the current node's value from `targetSum`. At a leaf, check if the remainder is 0. Recurse into left and right subtrees with the updated target.

### Algorithm
1. nil → false.
2. `remaining = targetSum - root.Val`.
3. If leaf: return `remaining == 0`.
4. Return `hasPathSum(left, remaining) || hasPathSum(right, remaining)`.

### Complexity
- **Time:** O(n)
- **Space:** O(h)

### Code
```go
func hasPathSum(root *TreeNode, targetSum int) bool {
    if root == nil { return false }
    remaining := targetSum - root.Val
    if root.Left == nil && root.Right == nil { return remaining == 0 }
    return hasPathSum(root.Left, remaining) || hasPathSum(root.Right, remaining)
}
```

### Dry Run
`targetSum=22`, path 5→4→11→2:

| Node | remaining before | remaining after |
|------|-----------------|-----------------|
| 5    | 22              | 17              |
| 4    | 17              | 13              |
| 11   | 13              | 2               |
| 2    | 2               | 0 → leaf → true |

---

## Approach 2 — Iterative DFS

### Intuition
Use a stack of `(node, remaining)` pairs. Pop each entry, compute new remaining, check at leaf.

### Complexity
- **Time:** O(n)
- **Space:** O(h)

### Code
```go
func hasPathSumIterative(root *TreeNode, targetSum int) bool {
    if root == nil { return false }
    type item struct { node *TreeNode; remaining int }
    stack := []item{{root, targetSum}}
    for len(stack) > 0 {
        curr := stack[len(stack)-1]; stack = stack[:len(stack)-1]
        rem := curr.remaining - curr.node.Val
        if curr.node.Left == nil && curr.node.Right == nil && rem == 0 { return true }
        if curr.node.Left  != nil { stack = append(stack, item{curr.node.Left,  rem}) }
        if curr.node.Right != nil { stack = append(stack, item{curr.node.Right, rem}) }
    }
    return false
}
```

### Dry Run
`targetSum=22`, tree `[1,2,3]`:

| Stack               | Pop     | rem | leaf? |
|--------------------|---------|-----|-------|
| [(1,22)]           | (1,22)  | 21  | no    |
| [(2,21),(3,21)]    | (3,21)  | 18  | leaf, 18≠0 |
| [(2,21)]           | (2,21)  | 19  | leaf, 19≠0 |

Return false ✓

---

## Key Takeaways
- Subtract-as-you-go instead of accumulating a sum avoids needing to carry the path sum down and check at the end.
- Short-circuit OR (`||`) is important — stop as soon as left returns true.
- Empty tree (`nil`) must return false even if `targetSum == 0`.

---

## Related Problems
- LeetCode #113 — Path Sum II (collect all qualifying paths)
- LeetCode #437 — Path Sum III (any node to any node, not root-to-leaf)
- LeetCode #666 — Path Sum IV
