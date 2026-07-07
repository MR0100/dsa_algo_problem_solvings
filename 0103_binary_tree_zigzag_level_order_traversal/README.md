# 0103 — Binary Tree Zigzag Level Order Traversal

> LeetCode #103 · Difficulty: Medium
> **Categories:** Tree, Breadth-First Search, Binary Tree

---

## Problem Statement

Given the `root` of a binary tree, return the zigzag level order traversal of its nodes' values (i.e., from left to right, then right to left for the next level and alternate between).

**Example 1:**
```
Input: root = [3,9,20,null,null,15,7]
Output: [[3],[20,9],[15,7]]
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
- `-100 <= Node.val <= 100`

---

## Company Frequency

| Company   | Frequency       | Last Reported |
|-----------|-----------------|---------------|
| Amazon    | ★★★★★ Very High | 2024          |
| Microsoft | ★★★★☆ High      | 2024          |
| Bloomberg | ★★★★☆ High      | 2023          |
| Google    | ★★★☆☆ Medium    | 2024          |

> ⚠️ Frequency data is crowd-sourced from LeetCode Discuss, Glassdoor, and community interview reports. Treat as a signal, not a guarantee.

---

## DSA Concepts Used

- **BFS Queue** — standard level order, with direction flag → see [`/dsa/tree_traversal.md`](/dsa/tree_traversal.md)
- **DFS with Depth** — odd depth = prepend value, even depth = append

---

## Approaches Overview

| # | Approach               | Time | Space | When to use            |
|---|------------------------|------|-------|------------------------|
| 1 | BFS + direction flag   | O(n) | O(w)  | Most intuitive         |
| 2 | DFS + depth parity     | O(n) | O(h)  | DFS-preferred contexts |

---

## Approach 1 — BFS with Direction Flag

### Intuition
Same as #102 BFS, but instead of always filling `level[i]`, we use a `leftToRight` boolean. When true, write `level[i]` left to right. When false, write `level[levelSize-1-i]` — fill from the right end without reversing after.

### Algorithm
1. Queue with root, `leftToRight = true`.
2. Each level: allocate `level := make([]int, levelSize)`.
3. Fill position:
   - leftToRight: `level[i] = node.Val`
   - rightToLeft: `level[levelSize-1-i] = node.Val`
4. Toggle `leftToRight`.

### Complexity
- **Time:** O(n)
- **Space:** O(w)

### Code
```go
func zigzagLevelOrder(root *TreeNode) [][]int {
    if root == nil { return nil }
    var result [][]int
    queue := []*TreeNode{root}
    leftToRight := true
    for len(queue) > 0 {
        levelSize := len(queue)
        level := make([]int, levelSize)
        for i := 0; i < levelSize; i++ {
            node := queue[0]; queue = queue[1:]
            if leftToRight {
                level[i] = node.Val
            } else {
                level[levelSize-1-i] = node.Val
            }
            if node.Left != nil  { queue = append(queue, node.Left)  }
            if node.Right != nil { queue = append(queue, node.Right) }
        }
        result = append(result, level)
        leftToRight = !leftToRight
    }
    return result
}
```

### Dry Run
Input: `[3,9,20,null,null,15,7]`, `leftToRight=true`

| Level | levelSize | Direction | Fill order      | level    |
|-------|-----------|-----------|-----------------|----------|
| 0     | 1         | L→R       | [0]=3           | [3]      |
| 1     | 2         | R→L       | [1]=9, [0]=20   | [20, 9]  |
| 2     | 2         | L→R       | [0]=15, [1]=7   | [15, 7]  |

Result: `[[3],[20,9],[15,7]]`

---

## Approach 2 — DFS with Depth Parity

### Intuition
DFS with depth. Even depths: append normally (left-to-right). Odd depths: prepend (right-to-left). Prepend is O(n) per node in worst case but acceptable for this problem size.

### Algorithm
1. `dfs(node, depth)`:
   - Start new slice at `result[depth]` if needed.
   - `depth%2 == 0`: append; `depth%2 == 1`: prepend.
   - Recurse left then right.

### Complexity
- **Time:** O(n²) worst case due to prepend (slice copy). O(n) in balanced trees.
- **Space:** O(h)

### Code
```go
func zigzagLevelOrderDFS(root *TreeNode) [][]int {
    var result [][]int
    var dfs func(node *TreeNode, depth int)
    dfs = func(node *TreeNode, depth int) {
        if node == nil { return }
        if depth == len(result) { result = append(result, []int{}) }
        if depth%2 == 0 {
            result[depth] = append(result[depth], node.Val)
        } else {
            result[depth] = append([]int{node.Val}, result[depth]...)
        }
        dfs(node.Left, depth+1)
        dfs(node.Right, depth+1)
    }
    dfs(root, 0)
    return result
}
```

### Dry Run
Input: `[3,9,20,null,null,15,7]`

| Call        | depth | parity | action           | result[depth] |
|-------------|-------|--------|------------------|---------------|
| dfs(3, 0)   | 0     | even   | append 3         | [3]           |
| dfs(9, 1)   | 1     | odd    | prepend 9        | [9]           |
| dfs(20, 1)  | 1     | odd    | prepend 20       | [20, 9]       |
| dfs(15, 2)  | 2     | even   | append 15        | [15]          |
| dfs(7, 2)   | 2     | even   | append 7         | [15, 7]       |

---

## Key Takeaways
- Pre-allocate level slice with `make([]int, levelSize)` and fill by index avoids reversing after collection.
- Zigzag = #102 BFS + alternating fill direction per level.
- DFS zigzag: even depth → append; odd depth → prepend.

---

## Related Problems
- LeetCode #102 — Binary Tree Level Order Traversal (straight BFS, no zigzag)
- LeetCode #107 — Binary Tree Level Order Traversal II (bottom-up)
