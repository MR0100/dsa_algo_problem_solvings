# 0114 — Flatten Binary Tree to Linked List

> LeetCode #114 · Difficulty: Medium
> **Categories:** Linked List, Stack, Tree, Depth-First Search, Binary Tree

---

## Problem Statement

Given the `root` of a binary tree, flatten the tree into a "linked list":

- The "linked list" should use the same `TreeNode` class where the `right` child pointer points to the next node in the list and the `left` child pointer is always `null`.
- The "linked list" should be in the same order as a **preorder traversal** of the binary tree.

**Example 1:**
```
Input: root = [1,2,5,3,4,null,6]
Output: [1,null,2,null,3,null,4,null,5,null,6]
```

**Example 2:**
```
Input: root = []
Output: []
```

**Example 3:**
```
Input: root = [0]
Output: [0]
```

**Constraints:**
- The number of nodes in the tree is in the range `[0, 2000]`.
- `-100 <= Node.val <= 100`

**Follow up:** Can you flatten the tree in-place (with O(1) extra space)?

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

- **Preorder Traversal** — flattened order matches preorder
- **Morris Traversal (variant)** — find predecessor to rewire without extra space

---

## Approaches Overview

| # | Approach                     | Time | Space | When to use              |
|---|------------------------------|------|-------|--------------------------|
| 1 | Collect Preorder + Rewire    | O(n) | O(n)  | Simplest to understand   |
| 2 | Morris-like In-Place         | O(n) | O(1)  | Follow-up requirement    |
| 3 | Reverse Preorder (post-order)| O(n) | O(h)  | Elegant recursive trick  |

---

## Approach 1 — Collect Preorder then Rewire

### Intuition
Preorder traversal visits nodes in exactly the order the flattened list needs them. Collect node pointers into a slice, then wire `node[i].Right = node[i+1]`, `node[i].Left = nil`.

### Complexity
- **Time:** O(n)
- **Space:** O(n)

### Code
```go
func flatten(root *TreeNode) {
    var nodes []*TreeNode
    var preorder func(node *TreeNode)
    preorder = func(node *TreeNode) {
        if node == nil { return }
        nodes = append(nodes, node)
        preorder(node.Left); preorder(node.Right)
    }
    preorder(root)
    for i := 0; i < len(nodes)-1; i++ {
        nodes[i].Left = nil; nodes[i].Right = nodes[i+1]
    }
}
```

### Dry Run
`[1,2,5,3,4,null,6]` preorder: `[1,2,3,4,5,6]`.
Wire: 1→2→3→4→5→6 (all left=nil). ✓

---

## Approach 2 — Morris-Like In-Place (O(1) Space)

### Intuition
For each node with a left child:
1. Find the **rightmost** node of the left subtree.
2. Attach the current right subtree to that rightmost node's right.
3. Move the left subtree to become the right child; clear left.

Advance to `curr.Right` and repeat.

### Algorithm
```
curr = root
while curr != nil:
    if curr.Left != nil:
        rightmost = find_rightmost(curr.Left)
        rightmost.Right = curr.Right
        curr.Right = curr.Left
        curr.Left = nil
    curr = curr.Right
```

### Complexity
- **Time:** O(n) — each node visited at most twice.
- **Space:** O(1)

### Code
```go
func flattenInPlace(root *TreeNode) {
    curr := root
    for curr != nil {
        if curr.Left != nil {
            rightmost := curr.Left
            for rightmost.Right != nil { rightmost = rightmost.Right }
            rightmost.Right = curr.Right
            curr.Right = curr.Left
            curr.Left = nil
        }
        curr = curr.Right
    }
}
```

### Dry Run
`[1,2,5,3,4,null,6]`:

| curr | left | rightmost of left | action                          |
|------|------|-------------------|---------------------------------|
| 1    | 2    | 4                 | 4.Right=5, 1.Right=2, 1.Left=nil|
| 2    | 3    | 3                 | 3.Right=4, 2.Right=3, 2.Left=nil|
| 3    | nil  | —                 | advance                         |
| 4    | nil  | —                 | advance                         |
| 5    | nil  | —                 | advance                         |
| 6    | nil  | —                 | done                            |

Result: 1→2→3→4→5→6 ✓

---

## Approach 3 — Reverse Preorder (Right→Left→Root)

### Intuition
Process nodes in **reverse preorder**: right subtree, then left, then root. Maintain a `prev` pointer. Wire `node.Right = prev`, `node.Left = nil`. When we finish, the list is correctly ordered.

### Complexity
- **Time:** O(n)
- **Space:** O(h) — recursion stack.

### Code
```go
func flattenReverse(root *TreeNode) {
    var prev *TreeNode
    var dfs func(node *TreeNode)
    dfs = func(node *TreeNode) {
        if node == nil { return }
        dfs(node.Right); dfs(node.Left)
        node.Right = prev; node.Left = nil; prev = node
    }
    dfs(root)
}
```

### Dry Run
Reverse preorder of `[1,2,5,3,4,null,6]`: 6, 5, 4, 3, 2, 1.

| Process | prev before | node.Right set to | prev after |
|---------|-------------|-------------------|------------|
| 6       | nil         | nil               | 6          |
| 5       | 6           | 6                 | 5          |
| 4       | 5           | 5                 | 4          |
| 3       | 4           | 4                 | 3          |
| 2       | 3           | 3                 | 2          |
| 1       | 2           | 2                 | 1          |

Result: 1→2→3→4→5→6 ✓

---

## Key Takeaways
- Preorder = flattened order. Any O(1)-space method must simulate preorder without a stack.
- Morris-like: find rightmost of left subtree → stitch right subtree there → move left to right.
- Reverse-preorder trick: process right→left→root, wiring `prev` backward — reads forward.

---

## Related Problems
- LeetCode #116 — Populating Next Right Pointers in Each Node
- LeetCode #430 — Flatten a Multilevel Doubly Linked List
