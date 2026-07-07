# 0102 — Binary Tree Level Order Traversal

> LeetCode #102 · Difficulty: Medium
> **Categories:** Tree, Breadth-First Search, Binary Tree

---

## Problem Statement

Given the `root` of a binary tree, return the level order traversal of its nodes' values (i.e., from left to right, level by level).

**Example 1:**
```
Input: root = [3,9,20,null,null,15,7]
Output: [[3],[9,20],[15,7]]
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

| Company   | Frequency       | Last Reported |
|-----------|-----------------|---------------|
| Amazon    | ★★★★★ Very High | 2024          |
| Microsoft | ★★★★★ Very High | 2024          |
| Google    | ★★★★☆ High      | 2024          |
| Facebook  | ★★★★☆ High      | 2024          |
| Bloomberg | ★★★☆☆ Medium    | 2023          |

> ⚠️ Frequency data is crowd-sourced from LeetCode Discuss, Glassdoor, and community interview reports. Treat as a signal, not a guarantee.

---

## DSA Concepts Used

- **BFS / Queue** — process nodes level by level using a queue → see [`/dsa/tree_traversal.md`](/dsa/tree_traversal.md)
- **DFS with Depth** — use recursion depth to determine which level slice to append to

---

## Approaches Overview

| # | Approach           | Time | Space | When to use              |
|---|--------------------|------|-------|--------------------------|
| 1 | BFS Queue          | O(n) | O(w)  | Most natural for BFS     |
| 2 | DFS + depth param  | O(n) | O(h)  | When DFS is preferred    |

---

## Approach 1 — BFS Queue

### Intuition
BFS naturally processes nodes level by level. At the start of each BFS iteration, the queue contains exactly all nodes of the current level — its size tells us how many to dequeue for this level.

### Algorithm
1. If root is nil, return nil.
2. Initialize queue with root.
3. While queue is non-empty:
   - `levelSize = len(queue)`.
   - Dequeue `levelSize` nodes, collect their values.
   - Enqueue each node's non-nil children.
   - Append collected values to result.
4. Return result.

### Complexity
- **Time:** O(n) — each node enqueued/dequeued once.
- **Space:** O(w) — queue holds at most one full level (O(n) for complete tree).

### Code
```go
func levelOrder(root *TreeNode) [][]int {
    if root == nil { return nil }
    var result [][]int
    queue := []*TreeNode{root}
    for len(queue) > 0 {
        levelSize := len(queue)
        level := make([]int, 0, levelSize)
        for i := 0; i < levelSize; i++ {
            node := queue[0]; queue = queue[1:]
            level = append(level, node.Val)
            if node.Left != nil  { queue = append(queue, node.Left)  }
            if node.Right != nil { queue = append(queue, node.Right) }
        }
        result = append(result, level)
    }
    return result
}
```

### Dry Run
Input: `[3,9,20,null,null,15,7]`

| Iteration | Queue before     | levelSize | Collected | Queue after    |
|-----------|-----------------|-----------|-----------|----------------|
| 1         | [3]             | 1         | [3]       | [9, 20]        |
| 2         | [9, 20]         | 2         | [9, 20]   | [15, 7]        |
| 3         | [15, 7]         | 2         | [15, 7]   | []             |

Result: `[[3],[9,20],[15,7]]`

---

## Approach 2 — DFS with Depth

### Intuition
DFS pre-order tracks the depth at each call. When `depth == len(result)`, we're visiting a new level for the first time — start a new slice. Append the node's value to `result[depth]`.

### Algorithm
1. Define `dfs(node, depth)`:
   - Return if node is nil.
   - If `depth == len(result)`, append new empty slice to result.
   - `result[depth] = append(result[depth], node.Val)`.
   - Recurse left with `depth+1`, then right.
2. Call `dfs(root, 0)`.

### Complexity
- **Time:** O(n)
- **Space:** O(h) — recursion stack depth.

### Code
```go
func levelOrderDFS(root *TreeNode) [][]int {
    var result [][]int
    var dfs func(node *TreeNode, depth int)
    dfs = func(node *TreeNode, depth int) {
        if node == nil { return }
        if depth == len(result) { result = append(result, []int{}) }
        result[depth] = append(result[depth], node.Val)
        dfs(node.Left, depth+1)
        dfs(node.Right, depth+1)
    }
    dfs(root, 0)
    return result
}
```

### Dry Run
Input: `[3,9,20,null,null,15,7]`

| Call           | depth | result after                  |
|----------------|-------|-------------------------------|
| dfs(3, 0)      | 0     | [[3]]                         |
| dfs(9, 1)      | 1     | [[3],[9]]                     |
| dfs(nil, 2)    | —     | (return)                      |
| dfs(nil, 2)    | —     | (return)                      |
| dfs(20, 1)     | 1     | [[3],[9,20]]                  |
| dfs(15, 2)     | 2     | [[3],[9,20],[15]]             |
| dfs(7, 2)      | 2     | [[3],[9,20],[15,7]]           |

---

## Key Takeaways
- BFS level order: use queue + snapshot `levelSize = len(queue)` at start of each iteration.
- DFS can produce level order: track depth, `result[depth]` tells which level slice to use.
- The `depth == len(result)` guard is the DFS equivalent of "starting a new level."

---

## Related Problems
- LeetCode #103 — Binary Tree Zigzag Level Order Traversal (alternate direction per level)
- LeetCode #107 — Binary Tree Level Order Traversal II (bottom-up)
- LeetCode #199 — Binary Tree Right Side View (rightmost node per level)
- LeetCode #637 — Average of Levels in Binary Tree
