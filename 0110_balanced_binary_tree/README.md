# 0110 — Balanced Binary Tree

> LeetCode #110 · Difficulty: Easy
> **Categories:** Tree, Depth-First Search, Binary Tree

---

## Problem Statement

Given a binary tree, determine if it is **height-balanced**.

A height-balanced binary tree is a binary tree in which the depth of the two subtrees of every node never differs by more than one.

**Example 1:**
```
Input: root = [3,9,20,null,null,15,7]
Output: true
```

**Example 2:**
```
Input: root = [1,2,2,3,3,null,null,4,4]
Output: false
```

**Example 3:**
```
Input: root = []
Output: true
```

**Constraints:**
- The number of nodes in the tree is in the range `[0, 5000]`.
- `-10^4 <= Node.val <= 10^4`

---

## Company Frequency

| Company   | Frequency       | Last Reported |
|-----------|-----------------|---------------|
| Amazon    | ★★★★★ Very High | 2024          |
| Microsoft | ★★★★☆ High      | 2024          |
| Google    | ★★★★☆ High      | 2024          |
| Bloomberg | ★★★☆☆ Medium    | 2023          |

> ⚠️ Frequency data is crowd-sourced from LeetCode Discuss, Glassdoor, and community interview reports. Treat as a signal, not a guarantee.

---

## DSA Concepts Used

- **Tree DFS post-order** — compute height bottom-up → see [`/dsa/tree_traversal.md`](/dsa/tree_traversal.md)
- **Sentinel Value** — return -1 from `check()` to propagate failure without extra state

---

## Approaches Overview

| # | Approach                    | Time       | Space | When to use                    |
|---|-----------------------------|------------|-------|--------------------------------|
| 1 | Top-Down (naive)            | O(n log n) | O(h)  | Simple to understand           |
| 2 | Bottom-Up Optimal (Optimal) | O(n)       | O(h)  | Always; single-pass            |

---

## Approach 1 — Top-Down Recursive (Naive)

### Intuition
At each node, compute heights of both subtrees and check the difference. If balanced, recurse into children. Problem: `height()` is called repeatedly on the same nodes → O(n log n).

### Complexity
- **Time:** O(n log n) — height computation at each node takes O(n).
- **Space:** O(h)

### Code
```go
func isBalanced(root *TreeNode) bool {
    var height func(node *TreeNode) int
    height = func(node *TreeNode) int {
        if node == nil { return 0 }
        l, r := height(node.Left), height(node.Right)
        if l > r { return l + 1 }
        return r + 1
    }
    if root == nil { return true }
    l, r := height(root.Left), height(root.Right)
    diff := l - r; if diff < 0 { diff = -diff }
    return diff <= 1 && isBalanced(root.Left) && isBalanced(root.Right)
}
```

### Dry Run
Input: `[3,9,20,null,null,15,7]`

| Node | left height | right height | |left-right| | Balanced? |
|------|-------------|--------------|------------|-----------|
| 9    | 0           | 0            | 0          | ✓         |
| 15   | 0           | 0            | 0          | ✓         |
| 7    | 0           | 0            | 0          | ✓         |
| 20   | 1           | 1            | 0          | ✓         |
| 3    | 1           | 2            | 1          | ✓         |

---

## Approach 2 — Bottom-Up Optimal

### Intuition
Single post-order traversal. Return the node's height if balanced, or -1 as a sentinel if any subtree is unbalanced. Once -1 is detected, propagate it upward immediately without further work.

### Algorithm
1. `check(node)`:
   - nil → 0.
   - `l = check(left)`;  if l==-1 → return -1.
   - `r = check(right)`; if r==-1 → return -1.
   - `|l-r| > 1` → return -1.
   - Return `1 + max(l, r)`.
2. `isBalanced` = `check(root) != -1`.

### Complexity
- **Time:** O(n) — each node visited exactly once.
- **Space:** O(h)

### Code
```go
func isBalancedOptimal(root *TreeNode) bool {
    var check func(node *TreeNode) int
    check = func(node *TreeNode) int {
        if node == nil { return 0 }
        l := check(node.Left);  if l == -1 { return -1 }
        r := check(node.Right); if r == -1 { return -1 }
        diff := l - r; if diff < 0 { diff = -diff }
        if diff > 1 { return -1 }
        if l > r { return l + 1 }
        return r + 1
    }
    return check(root) != -1
}
```

### Dry Run
Input: `[1,2,2,3,3,null,null,4,4]`

Post-order traversal:
- check(4)=1, check(4)=1.
- check(3)=2 (|1-1|=0 ok).
- check(3)=1 (leaf).
- check(2): l=2, r=1, |2-1|=1 ✓ → 3.
- check(2): l=0, r=0 → 1.
- check(1): l=3, r=1, |3-1|=2 > 1 → **-1**.

Return: -1 ≠ -1 is false → `false`. ✓

---

## Key Takeaways
- Top-down is O(n log n); bottom-up is O(n). Always use bottom-up in interviews.
- The -1 sentinel lets a single DFS detect and propagate unbalanced state early.
- Pattern: "compute a property and detect failure in one DFS" → return sentinel on failure.

---

## Related Problems
- LeetCode #104 — Maximum Depth of Binary Tree (height computation)
- LeetCode #543 — Diameter of Binary Tree (similar bottom-up DFS)
- LeetCode #1382 — Balance a Binary Search Tree
