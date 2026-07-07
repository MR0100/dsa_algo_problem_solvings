# 0100 — Same Tree

> LeetCode #100 · Difficulty: Easy
> **Categories:** Tree, Depth-First Search, Breadth-First Search, Binary Tree

---

## Problem Statement

Given the roots of two binary trees `p` and `q`, write a function to check if they are the same or not.

Two binary trees are considered the same if they are structurally identical, and the nodes have the same value.

**Example 1:**
```
Input: p = [1,2,3], q = [1,2,3]
Output: true
```

**Example 2:**
```
Input: p = [1,2], q = [1,null,2]
Output: false
```

**Example 3:**
```
Input: p = [1,2,1], q = [1,1,2]
Output: false
```

**Constraints:**
- The number of nodes in both trees is in the range `[0, 100]`.
- `-10^4 <= Node.val <= 10^4`

---

## Company Frequency

| Company   | Frequency       | Last Reported |
|-----------|-----------------|---------------|
| Amazon    | ★★★★☆ High      | 2024          |
| Facebook  | ★★★☆☆ Medium    | 2024          |
| Microsoft | ★★★☆☆ Medium    | 2023          |
| Google    | ★★☆☆☆ Low       | 2023          |

> ⚠️ Frequency data is crowd-sourced from LeetCode Discuss, Glassdoor, and community interview reports. Treat as a signal, not a guarantee.

---

## DSA Concepts Used

- **Tree Traversal (DFS/BFS)** — simultaneously traverse both trees, comparing at each step. See [`/dsa/tree_traversal.md`](/dsa/tree_traversal.md)

---

## Approaches Overview

| # | Approach | Time | Space | When to use |
|---|----------|------|-------|-------------|
| 1 | Recursive DFS | O(n) | O(h) | Cleanest; 3 lines of logic |
| 2 | Iterative BFS | O(n) | O(w) | When recursion depth is a concern |

---

## Approach 1 — Recursive DFS

### Intuition
Two trees are the same iff:
1. Both are `nil` (both empty trees are the same), OR
2. Both are non-nil, have equal root values, and their left and right subtrees are recursively the same.

### Algorithm
1. If `p == nil && q == nil`: return `true`.
2. If `p == nil || q == nil`: return `false` (one exists, other doesn't).
3. If `p.Val != q.Val`: return `false`.
4. Return `isSameTree(p.Left, q.Left) && isSameTree(p.Right, q.Right)`.

### Complexity
- **Time:** O(n) — n = min(|p|, |q|); stops at first mismatch.
- **Space:** O(h) — recursion stack depth.

### Code
```go
func isSameTree(p *TreeNode, q *TreeNode) bool {
    if p == nil && q == nil { return true }
    if p == nil || q == nil { return false }
    if p.Val != q.Val { return false }
    return isSameTree(p.Left, q.Left) && isSameTree(p.Right, q.Right)
}
```

### Dry Run (p=[1,2,3], q=[1,2,3])

```
isSameTree(1,1): vals match → isSameTree(2,2) && isSameTree(3,3)
  isSameTree(2,2): match → isSameTree(nil,nil) && isSameTree(nil,nil) → true
  isSameTree(3,3): match → true && true → true
→ true ✓
```

### Dry Run (p=[1,2], q=[1,null,2])

```
isSameTree(1,1): match →
  isSameTree(2,nil): p≠nil, q=nil → false
→ false ✓
```

---

## Approach 2 — Iterative BFS

### Intuition
Use a queue of `(p, q)` node pairs. For each pair:
- Both nil: skip (same structure here).
- One nil: return false.
- Values differ: return false.
- Otherwise: enqueue their children pairs.

### Algorithm
1. `queue = [(p, q)]`.
2. While queue not empty:
   - Dequeue `(curr_p, curr_q)`.
   - If both nil: continue.
   - If one nil or values differ: return false.
   - Enqueue `(curr_p.Left, curr_q.Left)` and `(curr_p.Right, curr_q.Right)`.
3. Return true.

### Complexity
- **Time:** O(n)
- **Space:** O(w) — w = max width (O(n) worst case for complete tree).

### Code
```go
func isSameTreeBFS(p *TreeNode, q *TreeNode) bool {
    type pair struct{ p, q *TreeNode }
    queue := []pair{{p, q}}
    for len(queue) > 0 {
        curr := queue[0]; queue = queue[1:]
        if curr.p == nil && curr.q == nil { continue }
        if curr.p == nil || curr.q == nil { return false }
        if curr.p.Val != curr.q.Val { return false }
        queue = append(queue, pair{curr.p.Left, curr.q.Left})
        queue = append(queue, pair{curr.p.Right, curr.q.Right})
    }
    return true
}
```

### Dry Run (p=[1,2,3], q=[1,2,3])

| Queue | Action |
|-------|--------|
| [(1,1)] | dequeue; 1==1; enqueue (2,2),(3,3) |
| [(2,2),(3,3)] | dequeue (2,2); 2==2; enqueue (nil,nil),(nil,nil) |
| [(3,3),(nil,nil),(nil,nil)] | dequeue (3,3); 3==3; enqueue (nil,nil),(nil,nil) |
| [(nil,nil)×4] | all nil pairs → continue |
| [] | done → true ✓ |

---

## Key Takeaways
- The recursive solution is extremely concise — short-circuit evaluation (`&&`) naturally handles the case where left subtrees differ (don't even check right).
- `nil == nil` is the "structural match" base case; one nil + one non-nil is the "structural mismatch" case.
- BFS processes level by level; useful when the trees might be very deep (to avoid stack overflow).
- This same pattern (simultaneous DFS on two trees) appears in #101 (Symmetric Tree) and #572 (Subtree of Another Tree).

---

## Related Problems
- LeetCode #101 — Symmetric Tree (same recursive structure, compare left.left with right.right)
- LeetCode #572 — Subtree of Another Tree (check if p matches any subtree of q)
- LeetCode #617 — Merge Two Binary Trees (simultaneous traversal, merge values)
