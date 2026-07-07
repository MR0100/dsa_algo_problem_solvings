# 0116 — Populating Next Right Pointers in Each Node

> LeetCode #116 · Difficulty: Medium
> **Categories:** Linked List, Tree, Depth-First Search, Breadth-First Search, Binary Tree

---

## Problem Statement

You are given a **perfect binary tree** where all leaves are on the same level, and every parent has two children. The tree has a `Node` struct with an additional `Next` pointer.

Populate each `next` pointer to point to its next right node. If there is no next right node, set it to `NULL`.

Initially, all `next` pointers are set to `NULL`.

**Example 1:**
```
Input: root = [1,2,3,4,5,6,7]
Output: [1,#,2,3,#,4,5,6,7,#]
Explanation: # denotes the end of a level.
```

**Example 2:**
```
Input: root = []
Output: []
```

**Constraints:**
- The number of nodes in the tree is in the range `[0, 2^12 - 1]`.
- `-1000 <= node.val <= 1000`

**Follow-up:** Can you solve this using only O(1) extra space?

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

- **BFS level order** — natural fit for wiring nodes within each level
- **Perfect binary tree property** — enables O(1) space: each level is traversable via Next pointers set in the previous pass

---

## Approaches Overview

| # | Approach             | Time | Space | When to use              |
|---|----------------------|------|-------|--------------------------|
| 1 | BFS Queue            | O(n) | O(w)  | Simple and general       |
| 2 | O(1) Space (perfect) | O(n) | O(1)  | Follow-up; exploits tree |

---

## Approach 1 — BFS Level Order

### Intuition
Standard BFS collects one level at a time. For each node at position `i` in its level, set `node.Next = queue[0]` (the next node) if `i < levelSize-1`.

### Algorithm
1. Queue with root.
2. Each level: process `levelSize` nodes. For `i < levelSize-1`, set `node.Next = queue[0]`.
3. Enqueue children.

### Complexity
- **Time:** O(n)
- **Space:** O(w) = O(n/2) for perfect binary tree at leaf level.

### Code
```go
func connect(root *Node) *Node {
    if root == nil { return nil }
    queue := []*Node{root}
    for len(queue) > 0 {
        levelSize := len(queue)
        for i := 0; i < levelSize; i++ {
            node := queue[0]; queue = queue[1:]
            if i < levelSize-1 { node.Next = queue[0] }
            if node.Left  != nil { queue = append(queue, node.Left)  }
            if node.Right != nil { queue = append(queue, node.Right) }
        }
    }
    return root
}
```

### Dry Run
Level `[2, 3]` (levelSize=2):
- i=0: node=2, `2.Next = queue[0] = 3`. Enqueue 4, 5.
- i=1: node=3, `i == levelSize-1`, no Next. Enqueue 6, 7.

---

## Approach 2 — O(1) Space (Perfect Tree)

### Intuition
Use the `Next` pointers already set on the current level to traverse it and wire the next level. Two wiring cases:
1. **Same parent**: `node.Left.Next = node.Right`.
2. **Different parents**: `node.Right.Next = node.Next.Left`.

Start at `leftmost = root`, descend one level at a time.

### Algorithm
1. `leftmost = root`.
2. While `leftmost.Left != nil` (not leaf level):
   - `curr = leftmost`.
   - While `curr != nil`:
     - Wire `curr.Left.Next = curr.Right`.
     - If `curr.Next`: wire `curr.Right.Next = curr.Next.Left`.
     - `curr = curr.Next`.
   - `leftmost = leftmost.Left`.

### Complexity
- **Time:** O(n)
- **Space:** O(1)

### Code
```go
func connectO1(root *Node) *Node {
    if root == nil { return nil }
    leftmost := root
    for leftmost.Left != nil {
        curr := leftmost
        for curr != nil {
            curr.Left.Next = curr.Right
            if curr.Next != nil { curr.Right.Next = curr.Next.Left }
            curr = curr.Next
        }
        leftmost = leftmost.Left
    }
    return root
}
```

### Dry Run
Root=1. Level [1]:
- Wire 2.Next=3 (same parent: 1.Left.Next=1.Right). No curr.Next.
- leftmost = 2.

Level [2, 3]:
- curr=2: 4.Next=5. 2.Next=3 → 5.Next=6.
- curr=3: 6.Next=7. 3.Next=nil → skip.

Level [4,5,6,7]: leaf level → stop.

---

## Key Takeaways
- BFS approach works for any binary tree (see #117).
- O(1) approach exploits perfect binary tree: both children always exist and Next pointers from the previous level are available for traversal.
- Two wiring cases: same-parent (Left→Right) and cross-parent (Right→Next.Left).

---

## Related Problems
- LeetCode #117 — Populating Next Right Pointers II (arbitrary tree)
- LeetCode #102 — Binary Tree Level Order Traversal
