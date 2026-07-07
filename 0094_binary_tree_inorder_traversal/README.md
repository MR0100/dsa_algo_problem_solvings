# 0094 — Binary Tree Inorder Traversal

> LeetCode #94 · Difficulty: Easy
> **Categories:** Stack, Tree, Depth-First Search, Binary Tree

---

## Problem Statement

Given the `root` of a binary tree, return the **inorder traversal** of its nodes' values.

**Example 1:**
```
Input: root = [1,null,2,3]
Output: [1,3,2]
```

**Example 2:**
```
Input: root = []
Output: []
```

**Example 3:**
```
Input: root = [1]
Output: [1]
```

**Constraints:**
- The number of nodes in the tree is in the range `[0, 100]`.
- `-100 <= Node.val <= 100`

**Follow-up:** Recursive solution is trivial. Could you do it iteratively?

---

## Company Frequency

| Company   | Frequency       | Last Reported |
|-----------|-----------------|---------------|
| Amazon    | ★★★★★ Very High | 2024          |
| Microsoft | ★★★★☆ High      | 2024          |
| Google    | ★★★☆☆ Medium    | 2024          |
| Facebook  | ★★★☆☆ Medium    | 2023          |

> ⚠️ Frequency data is crowd-sourced from LeetCode Discuss, Glassdoor, and community interview reports. Treat as a signal, not a guarantee.

---

## DSA Concepts Used

- **Tree Traversal (DFS)** — inorder = left, root, right. See [`/dsa/tree_traversal.md`](/dsa/tree_traversal.md)
- **Stack** — iterative simulation of the recursive call stack.
- **Morris Threading** — O(1) space traversal by temporarily modifying tree pointers.

---

## Approaches Overview

| # | Approach | Time | Space | When to use |
|---|----------|------|-------|-------------|
| 1 | Recursive | O(n) | O(h) | Simplest; always valid |
| 2 | Iterative (Stack) | O(n) | O(h) | Follow-up; no recursion |
| 3 | Morris Traversal | O(n) | O(1) | Space-critical; modifies tree temporarily |

---

## Approach 1 — Recursive

### Intuition
Inorder traversal visits: left subtree, root, right subtree. This matches the natural recursive structure.

### Algorithm
1. If `node == nil`: return.
2. `inorder(node.Left)`.
3. Append `node.Val`.
4. `inorder(node.Right)`.

### Complexity
- **Time:** O(n) — each node visited once.
- **Space:** O(h) — call stack, where h = height. O(log n) balanced, O(n) skewed.

### Code
```go
func inorderRecursive(root *TreeNode) []int {
    var result []int
    var inorder func(node *TreeNode)
    inorder = func(node *TreeNode) {
        if node == nil { return }
        inorder(node.Left)
        result = append(result, node.Val)
        inorder(node.Right)
    }
    inorder(root)
    return result
}
```

### Dry Run (root=[1,null,2,3])

```
inorder(1):
  inorder(nil) [left of 1] → return
  append 1
  inorder(2):
    inorder(3):
      inorder(nil) → return
      append 3
      inorder(nil) → return
    append 2
    inorder(nil) → return
Result: [1, 3, 2]
```

---

## Approach 2 — Iterative (Stack)

### Intuition
Simulate the recursive call stack. Push left nodes aggressively. When `curr` is nil, the top of the stack is the next node to visit (it has no more left children). After visiting it, go right.

### Algorithm
1. `stack = [], curr = root`.
2. While `curr != nil || stack not empty`:
   - While `curr != nil`: push `curr`, `curr = curr.Left`.
   - `curr = pop`. Visit (append). `curr = curr.Right`.

### Complexity
- **Time:** O(n)
- **Space:** O(h)

### Code
```go
func inorderIterative(root *TreeNode) []int {
    var result []int
    stack := []*TreeNode{}
    curr := root
    for curr != nil || len(stack) > 0 {
        for curr != nil { stack = append(stack, curr); curr = curr.Left }
        curr = stack[len(stack)-1]; stack = stack[:len(stack)-1]
        result = append(result, curr.Val)
        curr = curr.Right
    }
    return result
}
```

### Dry Run (root=[1,null,2,3])

| curr | stack | action |
|------|-------|--------|
| 1 | [] | push 1, curr=nil |
| nil | [1] | pop 1, visit 1, curr=1.Right=2 |
| 2 | [] | push 2, curr=2.Left=3 |
| 3 | [2] | push 3, curr=3.Left=nil |
| nil | [2,3] | pop 3, visit 3, curr=3.Right=nil |
| nil | [2] | pop 2, visit 2, curr=2.Right=nil |
| nil | [] | done |

Result: [1,3,2] ✓

---

## Approach 3 — Morris Traversal (O(1) Space)

### Intuition
Instead of a stack, temporarily thread the inorder predecessor's right pointer to the current node. This lets us return to the current node after processing its left subtree, without a stack.

For each node `curr`:
1. If no left child: visit, go right.
2. Else: find inorder predecessor (`pred` = rightmost node in `curr.Left`).
   - If `pred.Right == nil`: thread (`pred.Right = curr`), go left.
   - If `pred.Right == curr`: unthread (`pred.Right = nil`), visit `curr`, go right.

### Complexity
- **Time:** O(n) — each node's predecessor is found at most twice.
- **Space:** O(1) — no stack; modifies tree pointers in-place (temporarily).

### Code
```go
func inorderMorris(root *TreeNode) []int {
    var result []int
    curr := root
    for curr != nil {
        if curr.Left == nil {
            result = append(result, curr.Val); curr = curr.Right
        } else {
            pred := curr.Left
            for pred.Right != nil && pred.Right != curr { pred = pred.Right }
            if pred.Right == nil {
                pred.Right = curr; curr = curr.Left
            } else {
                pred.Right = nil; result = append(result, curr.Val); curr = curr.Right
            }
        }
    }
    return result
}
```

### Dry Run (root=[1,null,2,3])

| curr | pred | action |
|------|------|--------|
| 1 | nil (no left) | visit 1, curr=2 |
| 2 | pred=3 (rightmost in 2.left=3), 3.Right=nil | thread 3→2, curr=3 |
| 3 | nil (no left) | visit 3, curr=3.Right=2 |
| 2 | pred=3, 3.Right==2 (threaded) | unthread 3.Right=nil, visit 2, curr=2.Right=nil |
| nil | done |

Result: [1,3,2] ✓

---

## Key Takeaways
- Inorder traversal of a BST gives elements in sorted order — critical property.
- Morris threading: when `pred.Right == curr`, we're returning from the left subtree. Unthread and visit.
- The iterative stack approach is the standard follow-up answer; Morris is a bonus.
- Same structure applies to preorder and postorder (with minor modifications).

---

## Related Problems
- LeetCode #144 — Binary Tree Preorder Traversal
- LeetCode #145 — Binary Tree Postorder Traversal
- LeetCode #98 — Validate Binary Search Tree (inorder gives sorted sequence)
- LeetCode #230 — Kth Smallest Element in a BST (inorder, stop at k)
