# 0255 — Verify Preorder Sequence in Binary Search Tree

> LeetCode #255 · Difficulty: Medium
> **Categories:** Stack, Tree, Binary Search Tree, Recursion, Monotonic Stack, Divide and Conquer

---

## Problem Statement

Given an array of **unique** integers `preorder`, return `true` if it is the correct preorder traversal sequence of a binary search tree.

**Example 1:**

```
Input: preorder = [5,2,1,3,6]
Output: true
```

**Example 2:**

```
Input: preorder = [5,2,6,1,3]
Output: false
```

**Constraints:**

- `1 <= preorder.length <= 10^4`
- `1 <= preorder[i] <= 10^4`
- All the elements of `preorder` are **unique**.

**Follow-up:** Could you do it using only constant space complexity?

---

## Company Frequency

| Company   | Frequency        | Last Reported |
|-----------|------------------|---------------|
| Google    | ★★★★☆ High       | 2023          |
| Amazon    | ★★★☆☆ Medium     | 2022          |
| Microsoft | ★★★☆☆ Medium     | 2022          |
| Bloomberg | ★★☆☆☆ Low        | 2021          |

> ⚠️ Frequency data is crowd-sourced from LeetCode Discuss, Glassdoor, and community interview reports. Treat as a signal, not a guarantee.

---

## DSA Concepts Used

- **Binary search tree** — the ordering invariant (left < root < right) is what we validate → see [`/dsa/binary_search_tree.md`](/dsa/binary_search_tree.md)
- **Monotonic stack** — the optimal pass keeps a decreasing stack and a running lower bound → see [`/dsa/monotonic_stack.md`](/dsa/monotonic_stack.md)
- **Divide & conquer** — the recursive approach splits each range into left/right subtrees → see [`/dsa/divide_and_conquer.md`](/dsa/divide_and_conquer.md)
- **Stack** — the core auxiliary structure of the optimal solution → see [`/dsa/stack.md`](/dsa/stack.md)

---

## Approaches Overview

| # | Approach | Time | Space | When to use |
|---|----------|------|-------|-------------|
| 1 | Recursive Range Partition | O(n²) | O(n) | Clear divide-and-conquer intuition |
| 2 | Monotonic Stack + Lower Bound (Optimal) | O(n) | O(n) / O(1) in place | Best time; answers the follow-up |

---

## Approach 1 — Recursive Range Partition

### Intuition
In preorder, the first element of any (sub)sequence is the subtree's root. Everything after it that is smaller forms a contiguous prefix — the left subtree — and everything after that must all be greater than the root (the right subtree). If any element after the left-run is not greater than the root, the sequence is invalid. Recurse into both parts with tightened bounds.

### Algorithm
1. `verify(lo, hi, min, max)`: the first element `preorder[lo]` is the root; it must satisfy `min < root < max`.
2. Scan forward from `lo+1` while values `< root` — that marks the end of the left subtree.
3. Every remaining value up to `hi` (the right subtree) must be `> root`; otherwise return `false`.
4. Recurse on left with bounds `(min, root)` and right with `(root, max)`.

### Complexity
- **Time:** O(n²) worst case — skewed inputs re-scan overlapping ranges.
- **Space:** O(n) — recursion depth on a degenerate tree.

### Code
```go
func recursivePartition(preorder []int) bool {
	var verify func(lo, hi, min, max int) bool
	verify = func(lo, hi, min, max int) bool {
		if lo > hi {
			return true
		}
		root := preorder[lo]
		if root <= min || root >= max {
			return false
		}
		i := lo + 1
		for i <= hi && preorder[i] < root {
			i++
		}
		for j := i; j <= hi; j++ {
			if preorder[j] <= root {
				return false
			}
		}
		return verify(lo+1, i-1, min, root) && verify(i, hi, root, max)
	}
	return verify(0, len(preorder)-1, -1<<63, 1<<63-1)
}
```

### Dry Run
Input `[5,2,1,3,6]`, initial `verify(0, 4, -inf, +inf)`.

| call                     | root | bounds ok?   | left-run (< root) | right check (> root) | recurse            |
|--------------------------|------|--------------|-------------------|----------------------|--------------------|
| verify(0,4,-inf,+inf)    | 5    | yes          | [2,1,3] (idx 1..3)| [6] > 5 ok           | left(1,3,-inf,5), right(4,4,5,+inf) |
| verify(1,3,-inf,5)       | 2    | yes          | [1] (idx 2)       | [3] > 2 ok           | left(2,2), right(3,3) |
| verify(2,2,-inf,2)       | 1    | yes          | none              | none                 | both empty → true  |
| verify(3,3,2,5)          | 3    | yes          | none              | none                 | true               |
| verify(4,4,5,+inf)       | 6    | yes (6>5)    | none              | none                 | true               |

All branches return true → answer **true**.

---

## Approach 2 — Monotonic Stack + Lower Bound (Optimal)

### Intuition
Walk the preorder left to right. While we keep descending left children, values decrease — push them onto a stack kept decreasing. When we hit a value larger than the stack top, we've turned into a right subtree: pop all smaller values, and the **last popped** value becomes a new lower bound, because everything that follows lives in that popped node's right subtree and must exceed it. If any later value is `<= lowerBound`, the sequence is invalid.

### Algorithm
1. `lowerBound = -inf`; empty stack.
2. For each value `v`: if `v <= lowerBound`, return `false`.
3. While the stack is non-empty and `v > stack.top()`: `lowerBound = pop()`.
4. Push `v`.
5. After processing all values, return `true`.

### Complexity
- **Time:** O(n) — each element is pushed and popped at most once.
- **Space:** O(n) for the stack, reducible to O(1) by reusing the input array as the stack — answering the follow-up.

### Code
```go
func monotonicStack(preorder []int) bool {
	lowerBound := -1 << 63
	stack := []int{}
	for _, v := range preorder {
		if v <= lowerBound {
			return false
		}
		for len(stack) > 0 && v > stack[len(stack)-1] {
			lowerBound = stack[len(stack)-1]
			stack = stack[:len(stack)-1]
		}
		stack = append(stack, v)
	}
	return true
}
```

### Dry Run
Input `[5,2,1,3,6]`, `lowerBound = -inf`, `stack = []`.

| v | v ≤ lb? | pops (while v > top)        | new lowerBound | stack after |
|---|---------|----------------------------|----------------|-------------|
| 5 | no      | none                       | -inf           | [5]         |
| 2 | no      | none (2 > 5 false)         | -inf           | [5,2]       |
| 1 | no      | none (1 > 2 false)         | -inf           | [5,2,1]     |
| 3 | no      | pop 1, pop 2 (3 > 5 false) | 2              | [5,3]       |
| 6 | no (6>2)| pop 3, pop 5               | 5              | [6]         |

No violation → answer **true**. (For `[5,2,6,1,3]`, after 6 pops everything the lowerBound becomes 5, then `1 <= 5` triggers `false`.)

---

## Key Takeaways
- Preorder validity of a BST reduces to: whenever you move into a right subtree, nothing afterward may drop back below the ancestor you turned right from.
- A decreasing monotonic stack naturally models the left-descent chain; popping on a larger value models turning right and locking in a new floor.
- The stack can be maintained in the input array itself with an index pointer, giving O(1) extra space — the classic follow-up answer.
- The recursive range-partition is easier to explain but O(n²); the stack pass is the linear optimal.

---

## Related Problems
- LeetCode #98 — Validate Binary Search Tree (BST invariant with bounds)
- LeetCode #1008 — Construct BST from Preorder Traversal (build instead of verify)
- LeetCode #331 — Verify Preorder Serialization of a Binary Tree (preorder validity, slot counting)
- LeetCode #105 — Construct Binary Tree from Preorder and Inorder Traversal
