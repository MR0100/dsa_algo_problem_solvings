# 0109 — Convert Sorted List to Binary Search Tree

> LeetCode #109 · Difficulty: Medium
> **Categories:** Linked List, Divide and Conquer, Tree, Binary Search Tree, Binary Tree

---

## Problem Statement

Given the `head` of a singly linked list where elements are sorted in **ascending order**, convert it to a **height-balanced** binary search tree.

**Example 1:**
```
Input: head = [-10,-3,0,5,9]
Output: [0,-3,9,-10,null,5]
```

**Example 2:**
```
Input: head = []
Output: []
```

**Constraints:**
- The number of nodes in the linked list is in the range `[0, 2 * 10^4]`.
- `-10^5 <= Node.val <= 10^5`

---

## Company Frequency

| Company   | Frequency    | Last Reported |
|-----------|--------------|---------------|
| Amazon    | ★★★★☆ High  | 2024          |
| Google    | ★★★☆☆ Medium | 2024          |
| Microsoft | ★★★☆☆ Medium | 2023          |

> ⚠️ Frequency data is crowd-sourced from LeetCode Discuss, Glassdoor, and community interview reports. Treat as a signal, not a guarantee.

---

## DSA Concepts Used

- **Divide and Conquer** — same middle-element-as-root idea as #108
- **Inorder Simulation** — advance a shared list pointer during recursion to avoid O(n) space

---

## Approaches Overview

| # | Approach                | Time | Space    | When to use                  |
|---|-------------------------|------|----------|------------------------------|
| 1 | Array Copy → D&C        | O(n) | O(n)     | Simplest to reason about     |
| 2 | Inorder Simulation      | O(n) | O(log n) | Optimal; advanced technique  |

---

## Approach 1 — Array Copy, then Divide and Conquer

### Intuition
Copy all values to a slice. Now it's identical to #108 — pick middle, recurse.

### Algorithm
1. Walk list, collect values into `nums`.
2. Apply `build(lo, hi)` from #108.

### Complexity
- **Time:** O(n)
- **Space:** O(n) — array copy.

### Code
```go
func sortedListToBST(head *ListNode) *TreeNode {
    var nums []int
    for cur := head; cur != nil; cur = cur.Next { nums = append(nums, cur.Val) }
    var build func(lo, hi int) *TreeNode
    build = func(lo, hi int) *TreeNode {
        if lo > hi { return nil }
        mid := (lo + hi) / 2
        root := &TreeNode{Val: nums[mid]}
        root.Left  = build(lo, mid-1)
        root.Right = build(mid+1, hi)
        return root
    }
    return build(0, len(nums)-1)
}
```

### Dry Run
`[-10,-3,0,5,9]` → `nums=[-10,-3,0,5,9]` → same as #108 dry run.

---

## Approach 2 — Inorder Simulation (O(log n) Space)

### Intuition
Key insight: **inorder traversal visits nodes in sorted order**. If we recurse in inorder (left → root → right) and advance the list pointer each time we "visit" a root, each root gets assigned the correct sorted value without random access.

We use the `(lo, hi)` range only to know how many elements are in each subtree, not to index into the list.

### Algorithm
1. Count list length `n`.
2. Keep a shared pointer `cur = head`.
3. `build(lo, hi)`:
   - `lo > hi` → return nil.
   - `mid = (lo+hi)/2`.
   - `left = build(lo, mid-1)` ← this advances `cur` to the mid-th element.
   - `root = &TreeNode{Val: cur.Val}`, then `cur = cur.Next`.
   - `root.Left = left`.
   - `root.Right = build(mid+1, hi)`.
4. Return root.

### Complexity
- **Time:** O(n) — each list node consumed exactly once.
- **Space:** O(log n) — only the recursion stack (height of balanced BST).

### Code
```go
func sortedListToBSTInOrder(head *ListNode) *TreeNode {
    length := 0
    for cur := head; cur != nil; cur = cur.Next { length++ }
    cur := head
    var build func(lo, hi int) *TreeNode
    build = func(lo, hi int) *TreeNode {
        if lo > hi { return nil }
        mid := (lo + hi) / 2
        left := build(lo, mid-1)    // builds left subtree, advances cur
        root := &TreeNode{Val: cur.Val}
        cur = cur.Next              // consume this node
        root.Left = left
        root.Right = build(mid+1, hi)
        return root
    }
    return build(0, length-1)
}
```

### Dry Run
`[-10,-3,0,5,9]`, length=5.

`build(0,4)`: mid=2.
- `left = build(0,1)`:
  - mid=0.
  - `left2 = build(0,-1)` = nil.
  - root2 = cur=-10, cur→-3.
  - right2 = build(1,1): mid=1, left=nil, root=-3 (cur→0), right=nil. Return -3.
  - root2=(-10, right=-3). Return.
- root = cur=0 (cur→5). root.Left = (-10,right=-3).
- `right = build(3,4)`:
  - mid=3. left=build(3,2)=nil. root=5 (cur→9). right=build(4,4)=9 (cur→nil).
  - Return 5(right=9).

Final: 0(left=-10(right=-3), right=5(right=9)) ✓

---

## Key Takeaways
- Array copy approach: easy but O(n) extra space.
- Inorder simulation: O(log n) space by leveraging that inorder visits elements in order.
- The `lo,hi` range encodes size, not array indices — `cur` advances naturally.

---

## Related Problems
- LeetCode #108 — Convert Sorted Array to Binary Search Tree (array version)
- LeetCode #876 — Middle of the Linked List (slow/fast pointer)
