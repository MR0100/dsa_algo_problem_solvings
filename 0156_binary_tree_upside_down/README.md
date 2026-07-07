# 0156 — Binary Tree Upside Down

> LeetCode #156 · Difficulty: Medium (Premium)
> **Categories:** Tree, Linked List, Depth-First Search, Binary Tree

---

## Problem Statement

Given the `root` of a binary tree, turn the tree upside down and return *the new root*.

You can turn a binary tree upside down with the following steps:

1. The original left child becomes the new root.
2. The original root becomes the new right child.
3. The original right child becomes the new left child.

The mentioned steps are done level by level. It is **guaranteed** that every right node has a sibling (a left node that shares the same parent) and has no children.

**Example 1:**
```
Input:  root = [1,2,3,4,5]
Output: [4,5,2,null,null,3,1]

        1                4
       / \              / \
      2   3    ──►     5   2
     / \                  / \
    4   5                3   1
```

**Example 2:**
```
Input:  root = []
Output: []
```

**Example 3:**
```
Input:  root = [1]
Output: [1]
```

**Constraints:**
- The number of nodes in the tree will be in the range `[0, 10]`.
- `1 <= Node.val <= 10`
- Every right node in the tree has a sibling (a left node that shares the same parent).
- Every right node in the tree has no children.

---

## Company Frequency

| Company   | Frequency       | Last Reported |
|-----------|-----------------|---------------|
| LinkedIn  | ★★★★★ Very High | 2024          |
| Google    | ★★★☆☆ Medium    | 2023          |
| Amazon    | ★★☆☆☆ Low       | 2023          |
| Facebook  | ★★☆☆☆ Low       | 2022          |

> ⚠️ Frequency data is crowd-sourced from LeetCode Discuss, Glassdoor, and
> community interview reports. Treat as a signal, not a guarantee.

---

## DSA Concepts Used

- **Tree Traversal (DFS along the left spine)** — the whole transformation happens on the path root → root.Left → root.Left.Left … → see [`/dsa/tree_traversal.md`](/dsa/tree_traversal.md)
- **Linked List In-Place Reversal** — the left spine behaves exactly like a singly linked list being reversed; the optimal solution is `reverse list` with one extra carried pointer → see [`/dsa/linked_list.md`](/dsa/linked_list.md)

---

## Approaches Overview

| # | Approach | Time | Space | When to use |
|---|----------|------|-------|-------------|
| 1 | Explicit Stack (Brute Force) | O(n) | O(n) | Easiest to visualise: capture the spine, then rebuild |
| 2 | Recursion (Bottom-Up) | O(n) | O(n) recursion stack | Cleanest code; classic interview answer |
| 3 | Iterative Pointer Rewiring (Optimal) | O(n) | O(1) | Follow-up answer: no stack, no recursion |

---

## Approach 1 — Explicit Stack (Brute Force)

### Intuition
The constraints ("every right node has a sibling and no children") force the tree into a very specific shape: everything hangs off the **left spine**, and right children are leaf decorations on spine nodes. So the flip only reorders the spine — bottom becomes top. Capture the spine in a stack, then pop it back out: each spine node adopts its original parent as right child and its parent's right sibling as left child.

### Algorithm
1. If `root` is `nil`, return `nil`.
2. Walk `root → root.Left → …`, pushing every spine node onto a stack.
3. The top of the stack (the leftmost, deepest node) is the **new root**.
4. Pop pairs from the top: for spine node `curr` with original parent `parent`, set `curr.Left = parent.Right` and `curr.Right = parent`.
5. The original root is now a leaf — set its `Left` and `Right` to `nil`.
6. Return the new root.

### Complexity
- **Time:** O(n) — each node is pushed once and rewired once.
- **Space:** O(n) — in the worst case (a pure left chain) the stack holds all n nodes.

### Code
```go
func bruteForceStack(root *TreeNode) *TreeNode {
    if root == nil {
        return nil
    }
    stack := []*TreeNode{}
    for node := root; node != nil; node = node.Left {
        stack = append(stack, node)
    }
    newRoot := stack[len(stack)-1]
    for i := len(stack) - 1; i > 0; i-- {
        curr := stack[i]
        parent := stack[i-1]
        curr.Left = parent.Right
        curr.Right = parent
    }
    root.Left, root.Right = nil, nil
    return newRoot
}
```

### Dry Run
Example 1: `root = [1,2,3,4,5]` (spine is 1 → 2 → 4).

| Step | Action | stack | Pointer changes |
|------|--------|-------|-----------------|
| 1 | push spine | `[1, 2, 4]` | — |
| 2 | newRoot = top | `[1, 2, 4]` | `newRoot = 4` |
| 3 | i=2: curr=4, parent=2 | | `4.Left = 2.Right = 5`, `4.Right = 2` |
| 4 | i=1: curr=2, parent=1 | | `2.Left = 1.Right = 3`, `2.Right = 1` |
| 5 | clear old root | | `1.Left = nil`, `1.Right = nil` |

Result tree: `4 → (5, 2)`, `2 → (3, 1)` = `[4,5,2,null,null,3,1]` ✓

---

## Approach 2 — Recursion (Bottom-Up)

### Intuition
Post-order thinking: first flip everything **below** my left child; whatever root that produces is the global answer. Then fix my own little triangle — my left child becomes my parent (`left.Right = root`), and my right child slides under it (`left.Left = right`). Finally I detach my own children because I have become a leaf.

The key trick is saving `left` and `right` **before** recursing, because the recursion rewires pointers underneath us.

### Algorithm
1. Base case: if `root == nil` or `root.Left == nil`, this node is already the new root — return it.
2. Save `left = root.Left`, `right = root.Right`.
3. `newRoot = recursion(left)` — flip the deeper levels first.
4. Rewire the triangle: `left.Left = right` (rule 3), `left.Right = root` (rule 2).
5. `root.Left = nil`, `root.Right = nil` — the old root is now a leaf.
6. Return `newRoot` unchanged all the way up.

### Complexity
- **Time:** O(n) — every node is visited exactly once.
- **Space:** O(n) — recursion depth equals the spine length (O(h), worst case n).

### Code
```go
func recursion(root *TreeNode) *TreeNode {
    if root == nil || root.Left == nil {
        return root
    }
    left := root.Left
    right := root.Right
    newRoot := recursion(left)
    left.Left = right
    left.Right = root
    root.Left, root.Right = nil, nil
    return newRoot
}
```

### Dry Run
Example 1: `root = [1,2,3,4,5]`.

| Call frame | left | right | Recursive result | Rewiring done in this frame |
|------------|------|-------|------------------|------------------------------|
| `recursion(4)` | — | — | returns `4` (base case: no left child) | none |
| `recursion(2)` | 4 | 5 | `newRoot = 4` | `4.Left = 5`, `4.Right = 2`, `2.Left = 2.Right = nil` |
| `recursion(1)` | 2 | 3 | `newRoot = 4` | `2.Left = 3`, `2.Right = 1`, `1.Left = 1.Right = nil` |

Final: root `4`, children `(5, 2)`; node `2`, children `(3, 1)` → `[4,5,2,null,null,3,1]` ✓

---

## Approach 3 — Iterative Pointer Rewiring (Optimal)

### Intuition
The left spine **is** a singly linked list (via `Left` pointers). Flipping the tree is exactly reversing that list, with one extra piece of luggage: as we move down, each node must hand its **right sibling** down to the next spine node, which will adopt it as its new left child. Carry two "previous" pointers — `prev` (parent already flipped) and `prevRight` (that parent's original right child) — and rewire in a single pass. No stack, no recursion.

### Algorithm
1. Initialise `prev = nil`, `prevRight = nil`, `curr = root`.
2. While `curr != nil`:
   1. `next = curr.Left` — save where we walk next before destroying the pointer.
   2. `curr.Left = prevRight` — the parent's right sibling becomes my new left child.
   3. `prevRight = curr.Right` — stash my sibling for the next node down.
   4. `curr.Right = prev` — my parent becomes my new right child.
   5. `prev = curr`, `curr = next` — advance down the original spine.
3. Return `prev` (the last spine node processed = original leftmost = new root).

### Complexity
- **Time:** O(n) — one pass down the left spine; every node rewired once.
- **Space:** O(1) — only three auxiliary pointers regardless of tree size.

### Code
```go
func iterative(root *TreeNode) *TreeNode {
    var prev, prevRight *TreeNode
    curr := root
    for curr != nil {
        next := curr.Left
        curr.Left = prevRight
        prevRight = curr.Right
        curr.Right = prev
        prev = curr
        curr = next
    }
    return prev
}
```

### Dry Run
Example 1: `root = [1,2,3,4,5]` (spine 1 → 2 → 4).

| Iter | curr | next=curr.Left | curr.Left ← prevRight | prevRight ← curr.Right | curr.Right ← prev | prev after | curr after |
|------|------|----------------|-----------------------|------------------------|-------------------|------------|------------|
| 1 | 1 | 2 | `1.Left = nil` | `prevRight = 3` | `1.Right = nil` | 1 | 2 |
| 2 | 2 | 4 | `2.Left = 3` | `prevRight = 5` | `2.Right = 1` | 2 | 4 |
| 3 | 4 | nil | `4.Left = 5` | `prevRight = nil` | `4.Right = 2` | 4 | nil |

Loop exits, return `prev = 4`. Tree: `4 → (5, 2)`, `2 → (3, 1)` = `[4,5,2,null,null,3,1]` ✓

---

## Key Takeaways
- **Read the constraints for structure**: "every right node has a sibling and no children" collapses a general tree problem into a left-spine (linked-list) problem.
- The optimal solution is literally **linked-list reversal with one extra carried pointer** (`prevRight`). Recognising the reversal pattern instantly gives the O(1)-space answer.
- In the recursive version, **save child pointers before recursing** — the recursion mutates the structure underneath you.
- Bottom-up recursion pattern: the answer is produced at the deepest level and bubbled up unchanged; each frame only fixes its local pointers.

---

## Related Problems
- LeetCode #206 — Reverse Linked List (identical pointer-reversal pattern)
- LeetCode #92 — Reverse Linked List II (partial in-place reversal)
- LeetCode #114 — Flatten Binary Tree to Linked List (tree restructuring in place)
- LeetCode #226 — Invert Binary Tree (simpler tree transformation)
