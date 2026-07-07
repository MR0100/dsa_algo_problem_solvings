# 0272 — Closest Binary Search Tree Value II

> LeetCode #272 · Difficulty: Hard
> **Categories:** Tree, Binary Search Tree, Two Pointers, Heap (Priority Queue), Stack, Depth-First Search

---

## Problem Statement

Given the `root` of a binary search tree, a `target` value, and an integer `k`, return *the `k` values in the BST that are closest to the `target`*. You may return the answer in **any order**.

You are **guaranteed** to have only one unique set of `k` values in the BST that are closest to the `target`.

**Example 1:**

```
Input: root = [4,2,5,1,3], target = 3.714286, k = 2
Output: [4,3]
```

**Example 2:**

```
Input: root = [1], target = 0.000000, k = 1
Output: [1]
```

**Constraints:**

- The number of nodes in the tree is `n`.
- `1 <= k <= n <= 10^4`.
- `0 <= Node.val <= 10^9`
- `-10^9 <= target <= 10^9`

**Follow-up:** Assume that the BST is balanced. Could you solve it in less than `O(n)` runtime (where `n = total nodes`)?

---

## Company Frequency

| Company    | Frequency        | Last Reported |
|------------|------------------|---------------|
| Google     | ★★★★☆ High       | 2024          |
| Amazon     | ★★★☆☆ Medium     | 2023          |
| Meta       | ★★★☆☆ Medium     | 2023          |
| Microsoft  | ★★☆☆☆ Low        | 2023          |
| Uber       | ★★☆☆☆ Low        | 2022          |

> ⚠️ Frequency data is crowd-sourced from LeetCode Discuss, Glassdoor, and
> community interview reports. Treat as a signal, not a guarantee.

---

## DSA Concepts Used

- **Binary Search Tree** — inorder traversal of a BST yields the values in sorted order, which is what unlocks the sliding-window approach → see [`/dsa/binary_search_tree.md`](/dsa/binary_search_tree.md)
- **Tree Traversal (Inorder DFS)** — every approach walks the tree once to enumerate values → see [`/dsa/tree_traversal.md`](/dsa/tree_traversal.md)
- **Two Pointers** — on the sorted inorder array, two shrinking pointers isolate the k-wide closest window → see [`/dsa/two_pointers.md`](/dsa/two_pointers.md)
- **Heap / Priority Queue** — a size-k max-heap keyed by distance keeps only the best k candidates in one pass → see [`/dsa/heap_priority_queue.md`](/dsa/heap_priority_queue.md)

---

## Approaches Overview

| # | Approach | Time | Space | When to use |
|---|----------|------|-------|-------------|
| 1 | Inorder + Sort by Distance | O(n log n) | O(n) | Simplest correct baseline |
| 2 | Sliding Window on Sorted Array | O(n) | O(n) | Exploits BST → sorted inorder; the standard optimal |
| 3 | Max-Heap of Size k | O(n log k) | O(k) | When k ≪ n and you want O(k) extra space |

---

## Approach 1 — Inorder + Sort by Distance

### Intuition

"k closest to target" is a selection over all values. The most direct route ignores the tree shape entirely: gather every value, sort them by `|v - target|`, and take the first k. Correctness-first, not efficiency-first.

### Algorithm

1. Inorder-traverse the BST, collecting all values into a slice.
2. Sort the slice by ascending `|v - target|`.
3. Return the first k entries.

### Complexity

- **Time:** O(n log n) — the comparison sort dominates the linear traversal.
- **Space:** O(n) — all values are stored.

### Code

```go
func bruteForce(root *TreeNode, target float64, k int) []int {
	var vals []int
	var inorder func(*TreeNode)
	inorder = func(n *TreeNode) {
		if n == nil {
			return
		}
		inorder(n.Left)             // left subtree first
		vals = append(vals, n.Val)  // visit node
		inorder(n.Right)            // then right subtree
	}
	inorder(root)
	// Sort so that the closest-to-target values come first.
	sort.Slice(vals, func(i, j int) bool {
		return math.Abs(float64(vals[i])-target) < math.Abs(float64(vals[j])-target)
	})
	return vals[:k] // the k nearest
}
```

### Dry Run

Example 1: tree `[4,2,5,1,3]`, `target = 3.714286`, `k = 2`.

| Step | action | state |
|------|--------|-------|
| 1 | inorder collect | `vals = [1,2,3,4,5]` |
| 2 | distances to 3.714286 | `1→2.71, 2→1.71, 3→0.71, 4→0.29, 5→1.29` |
| 3 | sort by distance | `[4, 3, 5, 2, 1]` |
| 4 | take first k=2 | `[4, 3]` |

Result: `[4,3]` ✔

---

## Approach 2 — Sliding Window on Sorted Array

### Intuition

Inorder traversal of a BST is **sorted**. In a sorted array, the k closest values to any target always form a **contiguous** window of length k. So start with the full array and shrink it from whichever end is farther from target, one element at a time, until exactly k values remain — those are the answer.

### Algorithm

1. Inorder-traverse to get a sorted slice `vals`.
2. `lo = 0`, `hi = len(vals) - 1`. While `hi - lo >= k`: if the left end is farther from target than the right end, `lo++`; otherwise `hi--`.
3. Return `vals[lo : hi+1]`.

### Complexity

- **Time:** O(n) — traversal is O(n); shrinking discards `n - k` elements with O(1) work each.
- **Space:** O(n) — the sorted array (only O(k) beyond it for the result).

### Code

```go
func slidingWindow(root *TreeNode, target float64, k int) []int {
	var vals []int
	var inorder func(*TreeNode)
	inorder = func(n *TreeNode) {
		if n == nil {
			return
		}
		inorder(n.Left)
		vals = append(vals, n.Val)
		inorder(n.Right)
	}
	inorder(root) // vals is now sorted ascending

	lo, hi := 0, len(vals)-1
	// Repeatedly discard the endpoint farther from target until exactly k remain.
	for hi-lo >= k {
		if math.Abs(float64(vals[lo])-target) > math.Abs(float64(vals[hi])-target) {
			lo++ // left end is farther → drop it
		} else {
			hi-- // right end is farther (or tie) → drop it
		}
	}
	return vals[lo : hi+1] // the k contiguous nearest values
}
```

### Dry Run

Example 1: `vals = [1,2,3,4,5]`, `target = 3.714286`, `k = 2`.

| Step | lo | hi | hi-lo≥k? | left dist \|vals[lo]-t\| | right dist \|vals[hi]-t\| | action |
|------|----|----|----------|--------------------------|---------------------------|--------|
| 1 | 0 | 4 | 4≥2 yes | \|1-3.71\|=2.71 | \|5-3.71\|=1.29 | left farther → lo++ |
| 2 | 1 | 4 | 3≥2 yes | \|2-3.71\|=1.71 | \|5-3.71\|=1.29 | left farther → lo++ |
| 3 | 2 | 4 | 2≥2 yes | \|3-3.71\|=0.71 | \|5-3.71\|=1.29 | right farther → hi-- |
| 4 | 2 | 3 | 1≥2 no | — | — | stop |

Return `vals[2:4] = [3,4]` (any order) ✔

---

## Approach 3 — Max-Heap of Size k

### Intuition

When k is far smaller than n, sorting all n values is wasteful. Keep only the best k seen so far in a **max-heap keyed by distance**: the root is the *worst* of the current best k. For each new value, push it; if the heap overflows past k, pop the root (the farthest). After the full pass, the heap holds exactly the k closest.

### Algorithm

1. Inorder-traverse. For each value `v`: push `(|v - target|, v)` onto the max-heap; if heap size `> k`, pop the max.
2. Drain the heap into the result.

### Complexity

- **Time:** O(n log k) — n pushes/pops, each O(log k).
- **Space:** O(k) — the heap never exceeds k+1 elements.

### Code

```go
func maxHeapK(root *TreeNode, target float64, k int) []int {
	h := &distHeap{}
	heap.Init(h)
	var inorder func(*TreeNode)
	inorder = func(n *TreeNode) {
		if n == nil {
			return
		}
		inorder(n.Left)
		// Push this value keyed by its distance to target.
		heap.Push(h, item{dist: math.Abs(float64(n.Val) - target), val: n.Val})
		// If we now hold more than k, evict the farthest (heap root).
		if h.Len() > k {
			heap.Pop(h)
		}
		inorder(n.Right)
	}
	inorder(root)
	// Whatever remains in the heap are the k closest values.
	res := make([]int, 0, k)
	for h.Len() > 0 {
		res = append(res, heap.Pop(h).(item).val)
	}
	return res
}
```

### Dry Run

Example 1: inorder gives `1,2,3,4,5`, `target = 3.714286`, `k = 2`. Heap is a max-heap by distance (root = farthest).

| Visit v | dist | push then trim to k=2 | heap (val:dist) | evicted |
|---------|------|-----------------------|-----------------|---------|
| 1 | 2.71 | push | {1:2.71} | — |
| 2 | 1.71 | push | {1:2.71, 2:1.71} | — |
| 3 | 0.71 | push → size 3 → pop max (1:2.71) | {2:1.71, 3:0.71} | 1 |
| 4 | 0.29 | push → size 3 → pop max (2:1.71) | {3:0.71, 4:0.29} | 2 |
| 5 | 1.29 | push → size 3 → pop max (5:1.29) | {3:0.71, 4:0.29} | 5 |

Drain heap → `[4, 3]` (order may vary) ✔

---

## Key Takeaways

- **BST inorder = sorted array.** Whenever a BST problem asks about order, ranges, or nearest values, reach for inorder first.
- **k closest in a sorted array is a contiguous window** — shrink from the farther end (two pointers) for O(n), no sort needed.
- **Size-k heap** is the go-to when you want the "top/closest k" of a stream and k ≪ n, trading O(k) space for avoiding a full O(n log n) sort.
- The **follow-up** (sub-O(n) on a balanced BST) is solved with two stacks — a predecessor iterator and a successor iterator around target — merging the two nearest streams until k values are pulled; that yields O(k + log n).

---

## Related Problems

- LeetCode #270 — Closest Binary Search Tree Value (single closest value)
- LeetCode #658 — Find K Closest Elements (same window logic on a plain sorted array)
- LeetCode #94 — Binary Tree Inorder Traversal (the traversal primitive)
- LeetCode #230 — Kth Smallest Element in a BST (inorder + counting)
- LeetCode #973 — K Closest Points to Origin (size-k max-heap by distance)
