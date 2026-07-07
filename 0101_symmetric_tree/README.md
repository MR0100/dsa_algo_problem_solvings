# 0101 — Symmetric Tree

> LeetCode #101 · Difficulty: Easy
> **Categories:** Tree, Depth-First Search, Breadth-First Search, Binary Tree

---

## Problem Statement

Given the `root` of a binary tree, check whether it is a mirror of itself (i.e., symmetric around its center).

**Example 1:**
```
Input: root = [1,2,2,3,4,4,3]
Output: true
```

**Example 2:**
```
Input: root = [1,2,2,null,3,null,3]
Output: false
```

**Constraints:**
- The number of nodes in the tree is in the range `[1, 1000]`.
- `-100 <= Node.val <= 100`

**Follow up:** Could you solve it both recursively and iteratively?

---

## Company Frequency

| Company   | Frequency       | Last Reported |
|-----------|-----------------|---------------|
| Amazon    | ★★★★★ Very High | 2024          |
| Microsoft | ★★★★☆ High      | 2024          |
| Bloomberg | ★★★☆☆ Medium    | 2023          |
| Google    | ★★★☆☆ Medium    | 2024          |

> ⚠️ Frequency data is crowd-sourced from LeetCode Discuss, Glassdoor, and community interview reports. Treat as a signal, not a guarantee.

---

## DSA Concepts Used

- **Tree DFS** — recursive mirror comparison → see [`/dsa/tree_traversal.md`](/dsa/tree_traversal.md)
- **BFS / Queue** — iterative pair-comparison level by level

---

## Approaches Overview

| # | Approach       | Time | Space | When to use                |
|---|----------------|------|-------|----------------------------|
| 1 | Recursive DFS  | O(n) | O(h)  | Cleaner, most common       |
| 2 | Iterative BFS  | O(n) | O(w)  | When stack space is limited|

---

## Approach 1 — Recursive DFS

### Intuition
A tree is symmetric if its left subtree is a mirror image of its right subtree. Two subtrees are mirrors iff:
- Their roots have the same value.
- The left's left mirrors the right's right.
- The left's right mirrors the right's left.

### Algorithm
1. If root is nil, return true.
2. Define `mirror(left, right)`:
   - Both nil → true.
   - One nil or different values → false.
   - Otherwise → `mirror(left.Left, right.Right) && mirror(left.Right, right.Left)`.
3. Return `mirror(root.Left, root.Right)`.

### Complexity
- **Time:** O(n) — each node visited once.
- **Space:** O(h) — recursion stack height.

### Code
```go
func isSymmetric(root *TreeNode) bool {
    var mirror func(left, right *TreeNode) bool
    mirror = func(left, right *TreeNode) bool {
        if left == nil && right == nil {
            return true
        }
        if left == nil || right == nil {
            return false
        }
        return left.Val == right.Val &&
            mirror(left.Left, right.Right) &&
            mirror(left.Right, right.Left)
    }
    if root == nil {
        return true
    }
    return mirror(root.Left, root.Right)
}
```

### Dry Run
Input: `[1,2,2,3,4,4,3]`

| Call                    | left.Val | right.Val | Match? |
|-------------------------|----------|-----------|--------|
| mirror(2, 2)            | 2        | 2         | ✓      |
| mirror(3, 3)            | 3        | 3         | ✓      |
| mirror(nil, nil)        | —        | —         | ✓      |
| mirror(4, 4)            | 4        | 4         | ✓      |
| mirror(nil, nil)        | —        | —         | ✓      |

All comparisons pass → `true`.

---

## Approach 2 — Iterative BFS

### Intuition
Instead of recursion, use a queue of (left, right) pairs that should be mirrors. Start with `(root.Left, root.Right)`. For each pair, check equality and enqueue mirror children.

### Algorithm
1. Push `(root.Left, root.Right)` onto queue.
2. While queue not empty:
   - Dequeue `(l, r)`.
   - Both nil → continue.
   - One nil or `l.Val != r.Val` → return false.
   - Enqueue `(l.Left, r.Right)` and `(l.Right, r.Left)`.
3. Return true.

### Complexity
- **Time:** O(n)
- **Space:** O(w) — width of tree; O(n) worst case for complete tree.

### Code
```go
func isSymmetricIterative(root *TreeNode) bool {
    if root == nil {
        return true
    }
    type pair struct{ l, r *TreeNode }
    queue := []pair{{root.Left, root.Right}}
    for len(queue) > 0 {
        curr := queue[0]; queue = queue[1:]
        if curr.l == nil && curr.r == nil { continue }
        if curr.l == nil || curr.r == nil || curr.l.Val != curr.r.Val {
            return false
        }
        queue = append(queue, pair{curr.l.Left, curr.r.Right})
        queue = append(queue, pair{curr.l.Right, curr.r.Left})
    }
    return true
}
```

### Dry Run
Input: `[1,2,2,null,3,null,3]`

| Queue head  | l.Val | r.Val | Action          |
|-------------|-------|-------|-----------------|
| (2, 2)      | 2     | 2     | enqueue children|
| (nil, nil)  | —     | —     | continue        |
| (3, 3)  ← but wait: l=2.Right=3, r=2.Left=nil | 3 | nil | return false |

Result: `false`.

---

## Key Takeaways
- Mirror check = cross-compare children: `(l.Left, r.Right)` and `(l.Right, r.Left)`.
- Same pattern works for BFS (queue of pairs) or DFS (recursion with two pointers).
- Base cases: both nil = symmetric; one nil or value mismatch = not symmetric.

---

## Related Problems
- LeetCode #100 — Same Tree (compare two trees node-by-node, same recursive structure)
- LeetCode #226 — Invert Binary Tree (flip left/right children)
- LeetCode #572 — Subtree of Another Tree
