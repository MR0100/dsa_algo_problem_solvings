# 0108 — Convert Sorted Array to Binary Search Tree

> LeetCode #108 · Difficulty: Easy
> **Categories:** Array, Divide and Conquer, Tree, Binary Search Tree, Binary Tree

---

## Problem Statement

Given an integer array `nums` where the elements are sorted in **ascending order**, convert it to a **height-balanced** binary search tree.

**Example 1:**
```
Input: nums = [-10,-3,0,5,9]
Output: [0,-3,9,-10,null,5]
Explanation: [0,-10,5,null,-3,null,9] is also accepted.
```

**Example 2:**
```
Input: nums = [1,3]
Output: [3,1]
Explanation: [1,null,3] and [3,1] are both height-balanced BSTs.
```

**Constraints:**
- `1 <= nums.length <= 10^4`
- `-10^4 <= nums[i] <= 10^4`
- `nums` is sorted in strictly ascending order.

---

## Company Frequency

| Company   | Frequency       | Last Reported |
|-----------|-----------------|---------------|
| Amazon    | ★★★★★ Very High | 2024          |
| Google    | ★★★★☆ High      | 2024          |
| Microsoft | ★★★☆☆ Medium    | 2023          |
| Apple     | ★★★☆☆ Medium    | 2023          |

> ⚠️ Frequency data is crowd-sourced from LeetCode Discuss, Glassdoor, and community interview reports. Treat as a signal, not a guarantee.

---

## DSA Concepts Used

- **Divide and Conquer** — always pick the middle element as root
- **Binary Search Tree Property** — left < root < right guaranteed by sorted input

---

## Approaches Overview

| # | Approach                 | Time | Space    | When to use             |
|---|--------------------------|------|----------|-------------------------|
| 1 | Recursive D&C            | O(n) | O(log n) | Cleanest; standard      |
| 2 | Iterative (Queue)        | O(n) | O(n)     | Avoids recursion        |

---

## Approach 1 — Recursive Divide and Conquer

### Intuition
A height-balanced BST requires roughly equal left/right subtree sizes. The middle element of the sorted array perfectly splits it in half — make it the root, recurse on left half and right half.

### Algorithm
1. `build(lo, hi)`:
   - `lo > hi` → return nil.
   - `mid = (lo+hi)/2`.
   - `root = nums[mid]`.
   - `root.Left  = build(lo, mid-1)`.
   - `root.Right = build(mid+1, hi)`.
2. Return `build(0, n-1)`.

### Complexity
- **Time:** O(n) — each element used as a node exactly once.
- **Space:** O(log n) — recursion stack (balanced tree depth).

### Code
```go
func sortedArrayToBST(nums []int) *TreeNode {
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
`nums = [-10,-3,0,5,9]`

| Call          | lo | hi | mid | nums[mid] |
|---------------|----|----|-----|-----------|
| build(0,4)    | 0  | 4  | 2   | 0 (root)  |
| build(0,1)    | 0  | 1  | 0   | -10       |
| build(1,1)    | 1  | 1  | 1   | -3        |
| build(3,4)    | 3  | 4  | 3   | 5         |
| build(4,4)    | 4  | 4  | 4   | 9         |

Tree: 0(left=-10(right=-3), right=5(right=9)) — height-balanced ✓

---

## Approach 2 — Iterative Queue

### Intuition
Same D&C logic but BFS-order using a queue of `(parent, lo, hi, isLeft)` tasks. Avoid the recursion stack entirely.

### Complexity
- **Time:** O(n)
- **Space:** O(n) — queue stores O(n) tasks.

### Code
```go
func sortedArrayToBSTIterative(nums []int) *TreeNode {
    if len(nums) == 0 { return nil }
    type task struct { parent *TreeNode; lo, hi int; isLeft bool }
    mid := len(nums) / 2
    root := &TreeNode{Val: nums[mid]}
    queue := []task{{root,0,mid-1,true},{root,mid+1,len(nums)-1,false}}
    for len(queue) > 0 {
        curr := queue[0]; queue = queue[1:]
        if curr.lo > curr.hi { continue }
        m := (curr.lo + curr.hi) / 2
        node := &TreeNode{Val: nums[m]}
        if curr.isLeft { curr.parent.Left = node } else { curr.parent.Right = node }
        queue = append(queue, task{node,curr.lo,m-1,true}, task{node,m+1,curr.hi,false})
    }
    return root
}
```

### Dry Run
`nums=[-10,-3,0,5,9]`, mid=2 → root=0.

Queue: [(root,0,1,L),(root,3,4,R)]
- (root,0,1,L): m=0, node=-10. root.Left=-10. Enqueue (-10,0,-1,L),(-10,1,1,R).
- (root,3,4,R): m=3, node=5. root.Right=5. Enqueue (5,3,2,L),(5,4,4,R).
- (-10,0,-1,L): skip.
- (-10,1,1,R): m=1, node=-3. -10.Right=-3.
- (5,3,2,L): skip.
- (5,4,4,R): m=4, node=9. 5.Right=9.

---

## Key Takeaways
- Pick middle element as root → guarantees height balance.
- `mid = (lo+hi)/2` uses left-middle for even-length, giving one valid answer (right-middle `(lo+hi+1)/2` also accepted).
- Same pattern applies to sorted linked list (#109) but without random access.

---

## Related Problems
- LeetCode #109 — Convert Sorted List to Binary Search Tree (linked list version)
- LeetCode #1382 — Balance a Binary Search Tree
