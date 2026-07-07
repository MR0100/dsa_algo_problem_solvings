# 0107 — Binary Tree Level Order Traversal II

> LeetCode #107 · Difficulty: Medium
> **Categories:** Tree, Breadth-First Search, Binary Tree

---

## Problem Statement

Given the `root` of a binary tree, return the bottom-up level order traversal of its nodes' values (i.e., from left to right, level by level from leaf to root).

**Example 1:**
```
Input: root = [3,9,20,null,null,15,7]
Output: [[15,7],[9,20],[3]]
```

**Example 2:**
```
Input: root = [1]
Output: [[1]]
```

**Example 3:**
```
Input: root = []
Output: []
```

**Constraints:**
- The number of nodes in the tree is in the range `[0, 2000]`.
- `-1000 <= Node.val <= 1000`

---

## Company Frequency

| Company   | Frequency    | Last Reported |
|-----------|--------------|---------------|
| Amazon    | ★★★★☆ High  | 2024          |
| Microsoft | ★★★☆☆ Medium | 2023          |
| Bloomberg | ★★★☆☆ Medium | 2023          |

> ⚠️ Frequency data is crowd-sourced from LeetCode Discuss, Glassdoor, and community interview reports. Treat as a signal, not a guarantee.

---

## DSA Concepts Used

- **BFS** — standard level-order, then reverse → see [`/dsa/tree_traversal.md`](/dsa/tree_traversal.md)

---

## Approaches Overview

| # | Approach        | Time | Space | When to use             |
|---|-----------------|------|-------|-------------------------|
| 1 | BFS + Reverse   | O(n) | O(w)  | Direct extension of #102|
| 2 | DFS + Reverse   | O(n) | O(h)  | DFS-preferred contexts  |

---

## Approach 1 — BFS + Reverse

### Intuition
Identical to #102 BFS level order, then reverse the outer slice so leaf levels appear first.

### Algorithm
1. BFS level order → collect `result` top-down.
2. Reverse `result` in-place.

### Complexity
- **Time:** O(n)
- **Space:** O(w) — queue; O(n) for the result slice.

### Code
```go
func levelOrderBottom(root *TreeNode) [][]int {
    if root == nil { return nil }
    var result [][]int
    queue := []*TreeNode{root}
    for len(queue) > 0 {
        levelSize := len(queue)
        level := make([]int, 0, levelSize)
        for i := 0; i < levelSize; i++ {
            node := queue[0]; queue = queue[1:]
            level = append(level, node.Val)
            if node.Left  != nil { queue = append(queue, node.Left)  }
            if node.Right != nil { queue = append(queue, node.Right) }
        }
        result = append(result, level)
    }
    for i, j := 0, len(result)-1; i < j; i, j = i+1, j-1 {
        result[i], result[j] = result[j], result[i]
    }
    return result
}
```

### Dry Run
Input: `[3,9,20,null,null,15,7]`

BFS collects: `[[3],[9,20],[15,7]]`. After reverse: `[[15,7],[9,20],[3]]`.

---

## Approach 2 — DFS + Reverse

### Intuition
DFS with depth tracking (same as #102 DFS approach), then reverse at the end.

### Complexity
- **Time:** O(n)
- **Space:** O(h)

### Code
```go
func levelOrderBottomDFS(root *TreeNode) [][]int {
    var result [][]int
    var dfs func(node *TreeNode, depth int)
    dfs = func(node *TreeNode, depth int) {
        if node == nil { return }
        if depth == len(result) { result = append(result, []int{}) }
        result[depth] = append(result[depth], node.Val)
        dfs(node.Left, depth+1); dfs(node.Right, depth+1)
    }
    dfs(root, 0)
    for i, j := 0, len(result)-1; i < j; i, j = i+1, j-1 {
        result[i], result[j] = result[j], result[i]
    }
    return result
}
```

### Dry Run
DFS collects `[[3],[9,20],[15,7]]`. Reverse → `[[15,7],[9,20],[3]]`.

---

## Key Takeaways
- Bottom-up level order = top-down BFS + one reverse pass. No special data structure needed.
- The reverse is O(L) where L = number of levels, not O(n).

---

## Related Problems
- LeetCode #102 — Binary Tree Level Order Traversal (top-down version)
- LeetCode #199 — Binary Tree Right Side View
