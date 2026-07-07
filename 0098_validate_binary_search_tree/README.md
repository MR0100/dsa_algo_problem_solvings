# 0098 — Validate Binary Search Tree

> LeetCode #98 · Difficulty: Medium
> **Categories:** Tree, Depth-First Search, Binary Search Tree

---

## Problem Statement

Given the `root` of a binary tree, determine if it is a valid binary search tree (BST).

A **valid BST** is defined as follows:
- The left subtree of a node contains only nodes with keys **less than** the node's key.
- The right subtree of a node contains only nodes with keys **greater than** the node's key.
- Both the left and right subtrees must also be binary search trees.

**Example 1:**
```
Input: root = [2,1,3]
Output: true
```

**Example 2:**
```
Input: root = [5,1,4,null,null,3,6]
Output: false
Explanation: The root node's value is 5 but its right child's value is 4.
```

**Constraints:**
- The number of nodes in the tree is in the range `[1, 10^4]`.
- `-2^31 <= Node.val <= 2^31 - 1`

---

## Company Frequency

| Company   | Frequency       | Last Reported |
|-----------|-----------------|---------------|
| Amazon    | ★★★★★ Very High | 2024          |
| Facebook  | ★★★★☆ High      | 2024          |
| Microsoft | ★★★★☆ High      | 2024          |
| Google    | ★★★☆☆ Medium    | 2023          |
| Bloomberg | ★★★☆☆ Medium    | 2023          |

> ⚠️ Frequency data is crowd-sourced from LeetCode Discuss, Glassdoor, and community interview reports. Treat as a signal, not a guarantee.

---

## DSA Concepts Used

- **BST Property** — left < root < right, recursively for all nodes (not just direct children). See [`/dsa/tree_traversal.md`](/dsa/tree_traversal.md)
- **Inorder Traversal** — inorder of a valid BST is strictly increasing.
- **Min/Max Bounds** — pass a valid range `(min, max)` down the recursion.

---

## Approaches Overview

| # | Approach | Time | Space | When to use |
|---|----------|------|-------|-------------|
| 1 | Inorder Traversal | O(n) | O(h) | Natural; leverages BST's sorted inorder |
| 2 | Min/Max Range Check | O(n) | O(h) | Most explicit; passes bounds down |
| 3 | Iterative Inorder | O(n) | O(h) | No recursion |

---

## Approach 1 — Inorder Traversal

### Intuition
A BST's inorder traversal yields a **strictly increasing** sequence. Walk inorder; if any value ≤ previous value, the tree is invalid.

**Common mistake:** checking only direct parent-child pairs (e.g., node 10's left child is 5, node 5's right child is 7 — valid parent-child, but 7 > root 10, making it invalid). Inorder comparison catches this.

### Algorithm
1. Walk inorder (left, root, right).
2. Track previous value `prev`. If `current.Val <= prev`, return invalid.
3. Update `prev = current.Val`.

### Complexity
- **Time:** O(n)
- **Space:** O(h) — recursion stack.

### Code
```go
func isValidBSTInorder(root *TreeNode) bool {
    prev := math.MinInt64; valid := true
    var inorder func(node *TreeNode)
    inorder = func(node *TreeNode) {
        if node == nil || !valid { return }
        inorder(node.Left)
        if node.Val <= prev { valid = false; return }
        prev = node.Val
        inorder(node.Right)
    }
    inorder(root); return valid
}
```

### Dry Run (root=[5,1,4,null,null,3,6])

Inorder: 1, 5, 3, 6. At step `node=3`: `3 <= prev(5)` → invalid ✓.

---

## Approach 2 — Min/Max Range Check

### Intuition
For each node, check that its value is within a valid range `(min, max)`. Initially `(-∞, +∞)`. When going left, the upper bound becomes the parent's value. When going right, the lower bound becomes the parent's value.

This correctly handles the global invariant (not just parent-child).

### Algorithm
1. `check(node, min, max)`:
   - If `node == nil`: true.
   - If `node.Val <= min || node.Val >= max`: false.
   - Return `check(node.Left, min, node.Val) && check(node.Right, node.Val, max)`.
2. Start: `check(root, MinInt64, MaxInt64)`.

### Complexity
- **Time:** O(n)
- **Space:** O(h)

### Code
```go
func isValidBST(root *TreeNode) bool {
    var check func(node *TreeNode, min, max int) bool
    check = func(node *TreeNode, min, max int) bool {
        if node == nil { return true }
        if node.Val <= min || node.Val >= max { return false }
        return check(node.Left, min, node.Val) && check(node.Right, node.Val, max)
    }
    return check(root, math.MinInt64, math.MaxInt64)
}
```

### Dry Run (root=[5,1,4,null,null,3,6])

```
check(5, -∞, +∞): 5 in (-∞,+∞) ✓
  check(1, -∞, 5): 1 in (-∞,5) ✓ → check(nil,..)|check(nil,..) → true
  check(4, 5, +∞): 4 <= 5 → FALSE
```

Returns false ✓.

---

## Approach 3 — Iterative Inorder

### Code
```go
func isValidBSTIterative(root *TreeNode) bool {
    stack := []*TreeNode{}; prev := math.MinInt64; curr := root
    for curr != nil || len(stack) > 0 {
        for curr != nil { stack = append(stack, curr); curr = curr.Left }
        curr = stack[len(stack)-1]; stack = stack[:len(stack)-1]
        if curr.Val <= prev { return false }
        prev = curr.Val; curr = curr.Right
    }
    return true
}
```

---

## Key Takeaways
- Common mistake: only comparing with direct parent. Use min/max bounds or inorder comparison to catch cross-subtree violations.
- Use `math.MinInt64` and `math.MaxInt64` as initial bounds — but note node values can be exactly `INT_MIN` or `INT_MAX`. The strict inequality (`<=` for min, `>=` for max) handles this correctly since the bounds are exclusive.
- Inorder approach: stop early if any value ≤ previous. Min/max approach: prune subtrees that can't possibly be valid.

---

## Related Problems
- LeetCode #94 — Binary Tree Inorder Traversal
- LeetCode #99 — Recover Binary Search Tree (find two swapped nodes)
- LeetCode #230 — Kth Smallest Element in a BST
- LeetCode #701 — Insert into a Binary Search Tree
