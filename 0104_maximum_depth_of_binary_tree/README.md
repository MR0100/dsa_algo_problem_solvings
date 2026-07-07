# 0104 — Maximum Depth of Binary Tree

> LeetCode #104 · Difficulty: Easy
> **Categories:** Tree, Depth-First Search, Breadth-First Search, Binary Tree

---

## Problem Statement

Given the `root` of a binary tree, return its maximum depth.

A binary tree's **maximum depth** is the number of nodes along the longest path from the root node down to the farthest leaf node.

**Example 1:**
```
Input: root = [3,9,20,null,null,15,7]
Output: 3
```

**Example 2:**
```
Input: root = [1,null,2]
Output: 2
```

**Constraints:**
- The number of nodes in the tree is in the range `[0, 10^4]`.
- `-100 <= Node.val <= 100`

---

## Company Frequency

| Company   | Frequency       | Last Reported |
|-----------|-----------------|---------------|
| Amazon    | ★★★★★ Very High | 2024          |
| Google    | ★★★★☆ High      | 2024          |
| Microsoft | ★★★★☆ High      | 2024          |
| Apple     | ★★★☆☆ Medium    | 2023          |

> ⚠️ Frequency data is crowd-sourced from LeetCode Discuss, Glassdoor, and community interview reports. Treat as a signal, not a guarantee.

---

## DSA Concepts Used

- **Tree DFS** — `1 + max(left, right)` recurrence → see [`/dsa/tree_traversal.md`](/dsa/tree_traversal.md)
- **BFS** — count levels in queue-based traversal

---

## Approaches Overview

| # | Approach        | Time | Space | When to use                    |
|---|-----------------|------|-------|--------------------------------|
| 1 | Recursive DFS   | O(n) | O(h)  | Cleanest; always preferred     |
| 2 | Iterative BFS   | O(n) | O(w)  | Avoids recursion stack         |
| 3 | Iterative DFS   | O(n) | O(h)  | Stack-explicit, no recursion   |

---

## Approach 1 — Recursive DFS

### Intuition
The depth of a tree rooted at `node` is 1 plus the maximum of its children's depths. Base case: nil node has depth 0.

### Algorithm
1. If node is nil, return 0.
2. `left = maxDepth(node.Left)`.
3. `right = maxDepth(node.Right)`.
4. Return `1 + max(left, right)`.

### Complexity
- **Time:** O(n) — visits every node once.
- **Space:** O(h) — recursion stack; O(log n) balanced, O(n) skewed.

### Code
```go
func maxDepth(root *TreeNode) int {
    if root == nil { return 0 }
    left := maxDepth(root.Left)
    right := maxDepth(root.Right)
    if left > right { return left + 1 }
    return right + 1
}
```

### Dry Run
Input: `[3,9,20,null,null,15,7]`

| Call         | Returns |
|--------------|---------|
| maxDepth(9)  | 1       |
| maxDepth(15) | 1       |
| maxDepth(7)  | 1       |
| maxDepth(20) | 2       |
| maxDepth(3)  | 3       |

---

## Approach 2 — Iterative BFS

### Intuition
Level order BFS visits nodes level by level. Count the number of levels completed.

### Algorithm
1. Queue with root, `depth = 0`.
2. Each iteration: process all nodes in current queue (one full level), increment `depth`.
3. Return `depth`.

### Complexity
- **Time:** O(n)
- **Space:** O(w) — max width of tree.

### Code
```go
func maxDepthBFS(root *TreeNode) int {
    if root == nil { return 0 }
    depth := 0
    queue := []*TreeNode{root}
    for len(queue) > 0 {
        levelSize := len(queue); depth++
        for i := 0; i < levelSize; i++ {
            node := queue[0]; queue = queue[1:]
            if node.Left != nil  { queue = append(queue, node.Left)  }
            if node.Right != nil { queue = append(queue, node.Right) }
        }
    }
    return depth
}
```

### Dry Run
Input: `[3,9,20,null,null,15,7]`

| Iteration | Queue before | depth after |
|-----------|-------------|-------------|
| 1         | [3]         | 1           |
| 2         | [9, 20]     | 2           |
| 3         | [15, 7]     | 3           |

---

## Approach 3 — Iterative DFS (Stack)

### Intuition
DFS using an explicit stack storing `(node, depth)` pairs. Track the max depth seen.

### Algorithm
1. Push `(root, 1)` onto stack.
2. While stack non-empty: pop `(node, d)`, update `maxD = max(maxD, d)`, push children with `d+1`.
3. Return `maxD`.

### Complexity
- **Time:** O(n)
- **Space:** O(h)

### Code
```go
func maxDepthDFS(root *TreeNode) int {
    if root == nil { return 0 }
    type item struct { node *TreeNode; depth int }
    stack := []item{{root, 1}}
    maxD := 0
    for len(stack) > 0 {
        curr := stack[len(stack)-1]; stack = stack[:len(stack)-1]
        if curr.depth > maxD { maxD = curr.depth }
        if curr.node.Left  != nil { stack = append(stack, item{curr.node.Left,  curr.depth+1}) }
        if curr.node.Right != nil { stack = append(stack, item{curr.node.Right, curr.depth+1}) }
    }
    return maxD
}
```

### Dry Run
Input: `[3,9,20,null,null,15,7]`

| Stack (top→bottom)                | Pop           | maxD |
|-----------------------------------|---------------|------|
| [(3,1)]                           | (3,1)         | 1    |
| [(9,2),(20,2)]                    | (20,2)        | 2    |
| [(9,2),(15,3),(7,3)]              | (7,3)         | 3    |
| [(9,2),(15,3)]                    | (15,3)        | 3    |
| [(9,2)]                           | (9,2)         | 3    |

---

## Key Takeaways
- Recursive DFS depth formula: `1 + max(left, right)`.
- BFS alternative: count completed levels.
- For any tree depth/path problem, the `(node, depth)` pair pattern is the iterative DFS template.

---

## Related Problems
- LeetCode #111 — Minimum Depth of Binary Tree
- LeetCode #543 — Diameter of Binary Tree
- LeetCode #110 — Balanced Binary Tree
