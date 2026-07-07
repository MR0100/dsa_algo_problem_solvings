# 0111 — Minimum Depth of Binary Tree

> LeetCode #111 · Difficulty: Easy
> **Categories:** Tree, Depth-First Search, Breadth-First Search, Binary Tree

---

## Problem Statement

Given a binary tree, find its minimum depth.

The minimum depth is the number of nodes along the shortest path from the root node down to the nearest leaf node.

**Note:** A leaf is a node with no children.

**Example 1:**
```
Input: root = [3,9,20,null,null,15,7]
Output: 2
```

**Example 2:**
```
Input: root = [2,null,3,null,4,null,5,null,6]
Output: 5
```

**Constraints:**
- The number of nodes in the tree is in the range `[0, 10^5]`.
- `-1000 <= Node.val <= 1000`

---

## Company Frequency

| Company   | Frequency    | Last Reported |
|-----------|--------------|---------------|
| Amazon    | ★★★★☆ High  | 2024          |
| Microsoft | ★★★☆☆ Medium | 2023          |
| Google    | ★★★☆☆ Medium | 2023          |

> ⚠️ Frequency data is crowd-sourced from LeetCode Discuss, Glassdoor, and community interview reports. Treat as a signal, not a guarantee.

---

## DSA Concepts Used

- **Tree DFS** — must handle null-child case carefully
- **BFS early termination** — first leaf found = minimum depth

---

## Approaches Overview

| # | Approach         | Time | Space | When to use                         |
|---|------------------|------|-------|-------------------------------------|
| 1 | Recursive DFS    | O(n) | O(h)  | Simple, handles null children       |
| 2 | BFS (early exit) | O(n) | O(w)  | Optimal for wide/deep trees (stops early) |

---

## Approach 1 — Recursive DFS

### Intuition
The critical subtlety: if a node has only one child, we **cannot** take the path through nil — nil is not a leaf. We must go through the existing child.

### Algorithm
1. nil → 0.
2. Both children nil → leaf → 1.
3. Only right child → `1 + minDepth(right)`.
4. Only left child → `1 + minDepth(left)`.
5. Both children → `1 + min(minDepth(left), minDepth(right))`.

### Complexity
- **Time:** O(n)
- **Space:** O(h)

### Code
```go
func minDepth(root *TreeNode) int {
    if root == nil { return 0 }
    if root.Left == nil && root.Right == nil { return 1 }
    if root.Left  == nil { return 1 + minDepth(root.Right) }
    if root.Right == nil { return 1 + minDepth(root.Left)  }
    l, r := minDepth(root.Left), minDepth(root.Right)
    if l < r { return l + 1 }
    return r + 1
}
```

### Dry Run
Input: `[2,null,3,null,4,null,5,null,6]` (right-spine only)

| Node | left | right | path chosen |
|------|------|-------|-------------|
| 6    | nil  | nil   | leaf → 1    |
| 5    | nil  | 1     | nil-left → 2 |
| 4    | nil  | 2     | nil-left → 3 |
| 3    | nil  | 3     | nil-left → 4 |
| 2    | nil  | 4     | nil-left → 5 |

Result: 5 ✓

---

## Approach 2 — BFS (Early Termination)

### Intuition
BFS processes nodes level by level. The first leaf node we encounter is at the shallowest level — return immediately.

### Algorithm
1. Queue with root, `depth=0`.
2. Each level: increment depth, process all nodes.
3. When a node with no children is found, return `depth`.

### Complexity
- **Time:** O(n) worst case (skewed tree), but often much faster.
- **Space:** O(w)

### Code
```go
func minDepthBFS(root *TreeNode) int {
    if root == nil { return 0 }
    queue := []*TreeNode{root}; depth := 0
    for len(queue) > 0 {
        levelSize := len(queue); depth++
        for i := 0; i < levelSize; i++ {
            node := queue[0]; queue = queue[1:]
            if node.Left == nil && node.Right == nil { return depth }
            if node.Left  != nil { queue = append(queue, node.Left)  }
            if node.Right != nil { queue = append(queue, node.Right) }
        }
    }
    return depth
}
```

### Dry Run
Input: `[3,9,20,null,null,15,7]`

| Level | Queue       | depth | Leaf? |
|-------|-------------|-------|-------|
| 1     | [3]         | 1     | no    |
| 2     | [9,20]      | 2     | node 9 is leaf → return 2 |

---

## Key Takeaways
- Common mistake: treating nil as a leaf. Only nodes with no children are leaves.
- BFS terminates at first leaf — often faster than DFS for wide trees.
- Contrast with #104 maxDepth: maxDepth takes `max`, minDepth takes `min` — but must guard against nil children.

---

## Related Problems
- LeetCode #104 — Maximum Depth of Binary Tree
- LeetCode #112 — Path Sum (root-to-leaf path)
