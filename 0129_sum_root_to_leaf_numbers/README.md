# 0129 — Sum Root to Leaf Numbers

> LeetCode #129 · Difficulty: Medium
> **Categories:** Tree, Depth-First Search, Binary Tree

---

## Problem Statement

You are given the `root` of a binary tree containing digits from `0` to `9` only.

Each root-to-leaf path in the tree represents a number. For example, the root-to-leaf path `1 → 2 → 3` represents the number `123`.

Return the **total sum** of all root-to-leaf numbers.

**Example 1:**
```
Input: root = [1,2,3]
Output: 25
Explanation: 1→2 = 12, 1→3 = 13. Total = 25.
```

**Example 2:**
```
Input: root = [4,9,0,5,1]
Output: 1026
Explanation: 4→9→5=495, 4→9→1=491, 4→0=40. Total=1026.
```

**Constraints:**
- The number of nodes in the tree is in the range `[1, 1000]`.
- `0 <= Node.val <= 9`
- The depth of the tree will not exceed `10`.

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

- **DFS with running value** — multiply by 10 and add digit at each step

---

## Approaches Overview

| # | Approach        | Time | Space | When to use        |
|---|-----------------|------|-------|--------------------|
| 1 | Recursive DFS   | O(n) | O(h)  | Clean and concise  |
| 2 | Iterative DFS   | O(n) | O(h)  | Avoids recursion   |

---

## Approach 1 — Recursive DFS

### Intuition
Pass `running` (the number formed so far) down the recursion. At each node: `running = running*10 + node.Val`. At a leaf, add `running` to the total.

### Algorithm
1. `dfs(node, running)`:
   - nil → 0.
   - `running = running*10 + node.Val`.
   - Leaf → return `running`.
   - Return `dfs(left, running) + dfs(right, running)`.

### Complexity
- **Time:** O(n)
- **Space:** O(h)

### Code
```go
func sumNumbers(root *TreeNode) int {
    var dfs func(node *TreeNode, running int) int
    dfs = func(node *TreeNode, running int) int {
        if node == nil { return 0 }
        running = running*10 + node.Val
        if node.Left == nil && node.Right == nil { return running }
        return dfs(node.Left, running) + dfs(node.Right, running)
    }
    return dfs(root, 0)
}
```

### Dry Run
`[4,9,0,5,1]`:

| Call            | running |
|-----------------|---------|
| dfs(4, 0)       | 4       |
| dfs(9, 4)       | 49      |
| dfs(5, 49)      | 495 → leaf |
| dfs(1, 49)      | 491 → leaf |
| dfs(0, 4)       | 40 → leaf  |

Total: 495+491+40 = 1026 ✓

---

## Approach 2 — Iterative DFS

### Intuition
Stack of `(node, running)` pairs. Same logic as recursive.

### Complexity
- **Time:** O(n)
- **Space:** O(h)

### Code
```go
func sumNumbersIterative(root *TreeNode) int {
    if root == nil { return 0 }
    type item struct { node *TreeNode; running int }
    stack := []item{{root, 0}}; total := 0
    for len(stack) > 0 {
        curr := stack[len(stack)-1]; stack = stack[:len(stack)-1]
        running := curr.running*10 + curr.node.Val
        if curr.node.Left == nil && curr.node.Right == nil { total += running; continue }
        if curr.node.Left  != nil { stack = append(stack, item{curr.node.Left,  running}) }
        if curr.node.Right != nil { stack = append(stack, item{curr.node.Right, running}) }
    }
    return total
}
```

### Dry Run
Same as recursive, stack-based. LIFO order processes leaves correctly.

---

## Key Takeaways
- `running = running*10 + node.Val` encodes the digit sequence as a number.
- Only leaves contribute to the sum — intermediate nodes merely accumulate digits.
- Same structure as Path Sum (#112): pass a running value down, check at leaf.

---

## Related Problems
- LeetCode #112 — Path Sum
- LeetCode #113 — Path Sum II
- LeetCode #257 — Binary Tree Paths
