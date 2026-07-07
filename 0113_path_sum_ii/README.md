# 0113 — Path Sum II

> LeetCode #113 · Difficulty: Medium
> **Categories:** Backtracking, Tree, Depth-First Search, Binary Tree

---

## Problem Statement

Given the `root` of a binary tree and an integer `targetSum`, return all **root-to-leaf** paths where the sum of the node values in the path equals `targetSum`. Each path should be returned as a list of node values, not node references.

A **leaf** is a node with no children.

**Example 1:**
```
Input: root = [5,4,8,11,null,13,4,7,2,null,null,5,1], targetSum = 22
Output: [[5,4,11,2],[5,8,4,5]]
```

**Example 2:**
```
Input: root = [1,2,3], targetSum = 5
Output: []
```

**Example 3:**
```
Input: root = [1,2], targetSum = 0
Output: []
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
| Google    | ★★★☆☆ Medium    | 2024          |
| Bloomberg | ★★★☆☆ Medium    | 2023          |

> ⚠️ Frequency data is crowd-sourced from LeetCode Discuss, Glassdoor, and community interview reports. Treat as a signal, not a guarantee.

---

## DSA Concepts Used

- **Backtracking** — append on descent, pop on return → see [`/dsa/backtracking.md`](/dsa/backtracking.md)
- **Tree DFS** — standard recursive traversal

---

## Approaches Overview

| # | Approach           | Time   | Space | When to use     |
|---|--------------------|--------|-------|-----------------|
| 1 | Backtracking DFS   | O(n²)  | O(h)  | Standard; always|

---

## Approach 1 — Backtracking DFS

### Intuition
DFS with a running path slice. At each node, append the value. At a leaf with `remaining==0`, record a copy. On return, pop (backtrack) to undo the append.

### Algorithm
1. `path = []`, `result = []`.
2. `dfs(node, remaining)`:
   - Append `node.Val` to path, subtract from remaining.
   - If leaf and `remaining==0`: append copy of path to result.
   - Recurse left, recurse right.
   - Pop last element from path (backtrack).

### Complexity
- **Time:** O(n²) — up to n/2 leaves, copying path of length O(n) for each.
- **Space:** O(h) — path + recursion stack.

### Code
```go
func pathSum(root *TreeNode, targetSum int) [][]int {
    var result [][]int
    var path []int
    var dfs func(node *TreeNode, remaining int)
    dfs = func(node *TreeNode, remaining int) {
        if node == nil { return }
        path = append(path, node.Val)
        remaining -= node.Val
        if node.Left == nil && node.Right == nil && remaining == 0 {
            cp := make([]int, len(path)); copy(cp, path)
            result = append(result, cp)
        }
        dfs(node.Left, remaining)
        dfs(node.Right, remaining)
        path = path[:len(path)-1] // backtrack
    }
    dfs(root, targetSum)
    return result
}
```

### Dry Run
`targetSum=22`, path 5→4→11→2:

| Action           | path         | remaining |
|-----------------|--------------|-----------|
| visit 5          | [5]          | 17        |
| visit 4          | [5,4]        | 13        |
| visit 11         | [5,4,11]     | 2         |
| visit 7          | [5,4,11,7]   | -5, leaf≠0 |
| backtrack        | [5,4,11]     | —         |
| visit 2          | [5,4,11,2]   | 0, leaf✓ → record |
| backtrack        | [5,4,11]     | —         |
| backtrack        | [5,4]        | —         |
| ...              | (right side) | ...       |

---

## Key Takeaways
- **Must copy** the path slice before appending to result — `append(result, path)` shares the backing array.
- Backtrack by `path = path[:len(path)-1]` — simple and efficient.
- Classic backtracking: choose (append) → explore → unchoose (pop).

---

## Related Problems
- LeetCode #112 — Path Sum (boolean existence only)
- LeetCode #437 — Path Sum III (any prefix, not just root-to-leaf)
- LeetCode #257 — Binary Tree Paths (return string paths)
