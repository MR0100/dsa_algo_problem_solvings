# 0099 — Recover Binary Search Tree

> LeetCode #99 · Difficulty: Medium
> **Categories:** Tree, Depth-First Search, Binary Search Tree

---

## Problem Statement

You are given the `root` of a binary search tree (BST), where the values of **exactly two** nodes of the tree were swapped by mistake. Recover the tree without changing its structure.

**Example 1:**
```
Input: root = [1,3,null,null,2]
Output: [3,1,null,null,2]
Explanation: 3 cannot be a left child of 1 because 3 > 1. Swapping 1 and 3 makes the BST valid.
```

**Example 2:**
```
Input: root = [3,1,4,null,null,2]
Output: [2,1,4,null,null,3]
Explanation: 2 cannot be a left child of 4 because 2 < 4. Swapping 2 and 3 makes the BST valid.
```

**Constraints:**
- The number of nodes in the tree is in the range `[2, 1000]`.
- `-2^31 <= Node.val <= 2^31 - 1`

**Follow-up:** A solution using `O(n)` space is fairly straightforward. Could you devise a constant space solution?

---

## Company Frequency

| Company   | Frequency      | Last Reported |
|-----------|----------------|---------------|
| Amazon    | ★★★☆☆ Medium   | 2024          |
| Facebook  | ★★★☆☆ Medium   | 2023          |
| Microsoft | ★★☆☆☆ Low      | 2023          |

> ⚠️ Frequency data is crowd-sourced from LeetCode Discuss, Glassdoor, and community interview reports. Treat as a signal, not a guarantee.

---

## DSA Concepts Used

- **Inorder Traversal + Inversion Detection** — two swapped nodes create at most 2 "inversions" in the inorder sequence.
- **Morris Traversal** — O(1) space inorder traversal for the follow-up.

---

## Approaches Overview

| # | Approach | Time | Space | When to use |
|---|----------|------|-------|-------------|
| 1 | Inorder Traversal | O(n) | O(h) | Standard; easy to reason about |
| 2 | Morris Traversal | O(n) | O(1) | Follow-up; O(1) extra space |

---

## Approach 1 — Inorder Traversal (Find Two Swapped Nodes)

### Intuition
A valid BST's inorder traversal is strictly increasing. When two nodes are swapped, the inorder sequence has inversions (places where a value is greater than the next).

Two cases:
- **Adjacent swap** (e.g., [..., 3, 2, ...]): only 1 inversion point. `first = prev(3)`, `second = curr(2)`.
- **Non-adjacent swap** (e.g., [..., 5, ..., 3, ...]): 2 inversion points.
  - First inversion: `first = prev(5)`.
  - Second inversion: `second = curr(3)`.

By always updating `second = curr` at every inversion and only setting `first` once, both cases are handled correctly.

### Algorithm
1. Walk inorder. Track `prev`, `first`, `second`.
2. At each node: if `prev != nil && prev.Val > node.Val`:
   - If `first == nil`: `first = prev`.
   - `second = node` (always update).
3. After traversal: `swap(first.Val, second.Val)`.

### Complexity
- **Time:** O(n)
- **Space:** O(h) — recursion stack.

### Code
```go
func recoverTree(root *TreeNode) {
    var first, second, prev *TreeNode
    var inorder func(node *TreeNode)
    inorder = func(node *TreeNode) {
        if node == nil { return }
        inorder(node.Left)
        if prev != nil && prev.Val > node.Val {
            if first == nil { first = prev }
            second = node
        }
        prev = node
        inorder(node.Right)
    }
    inorder(root)
    first.Val, second.Val = second.Val, first.Val
}
```

### Dry Run (root=[3,1,4,null,null,2])

Inorder: 1, 3, 2, 4.
- Visit 1: prev=nil → prev=1.
- Visit 3: 3>1, no inversion → prev=3.
- Visit 2: 3>2 → inversion! first=node(3), second=node(2). prev=2.
- Visit 4: 2<4 → no inversion.

Only 1 inversion point → adjacent swap. Swap 3 and 2 → inorder becomes [1,2,3,4] ✓.

### Dry Run (root=[5,3,8,1,4,null,9,null,null,null,null,6,7])

Inorder: 1, 3, 4, 5, 6, 7, 8, 9 → normal. Consider swapping 5 and 3:

Inorder of corrupted tree: 1, **5**, 4, **3**, 6, 7, 8, 9.
- 5>4 → first=node(5), second=node(4). (1st inversion)
- 3<4? No: 4>3 → second=node(3). (2nd inversion, update second)

Swap 5 and 3 → restored ✓.

---

## Approach 2 — Morris Traversal (O(1) Space)

### Intuition
Apply the same inversion-detection logic but using Morris inorder traversal (no stack, O(1) space). The same `first`/`second` tracking applies; only the traversal mechanism changes.

### Code
```go
func recoverTreeMorris(root *TreeNode) {
    var first, second, prev *TreeNode
    curr := root
    for curr != nil {
        if curr.Left == nil {
            if prev != nil && prev.Val > curr.Val {
                if first == nil { first = prev }
                second = curr
            }
            prev = curr; curr = curr.Right
        } else {
            pred := curr.Left
            for pred.Right != nil && pred.Right != curr { pred = pred.Right }
            if pred.Right == nil {
                pred.Right = curr; curr = curr.Left
            } else {
                pred.Right = nil
                if prev != nil && prev.Val > curr.Val {
                    if first == nil { first = prev }
                    second = curr
                }
                prev = curr; curr = curr.Right
            }
        }
    }
    first.Val, second.Val = second.Val, first.Val
}
```

---

## Key Takeaways
- `second` is always updated at every inversion; `first` is only set on the first inversion. This elegantly handles both adjacent and non-adjacent swaps with a single pass.
- Only values are swapped (`first.Val, second.Val = second.Val, first.Val`), not nodes — the tree structure is preserved.
- Morris traversal allows O(1) space at the cost of temporarily modifying right pointers.

---

## Related Problems
- LeetCode #94 — Binary Tree Inorder Traversal
- LeetCode #98 — Validate Binary Search Tree
- LeetCode #501 — Find Mode in Binary Search Tree (inorder + frequency count)
