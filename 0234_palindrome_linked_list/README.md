# 0234 — Palindrome Linked List

> LeetCode #234 · Difficulty: Easy
> **Categories:** Linked List, Two Pointers, Stack, Recursion

---

## Problem Statement

Given the `head` of a singly linked list, return `true` if it is a palindrome or `false` otherwise.

**Example 1:**
```
Input: head = [1,2,2,1]
Output: true
```

**Example 2:**
```
Input: head = [1,2]
Output: false
```

**Constraints:**
- The number of nodes in the list is in the range `[1, 10⁵]`.
- `0 <= Node.val <= 9`

**Follow up:** Could you do it in `O(n)` time and `O(1)` space?

---

## Company Frequency

| Company    | Frequency        | Last Reported |
|------------|------------------|---------------|
| Amazon     | ★★★★☆ High       | 2024          |
| Microsoft  | ★★★☆☆ Medium     | 2024          |
| Meta       | ★★★☆☆ Medium     | 2024          |
| Google     | ★★★☆☆ Medium     | 2023          |
| Apple      | ★★☆☆☆ Low        | 2023          |
| Bloomberg  | ★★☆☆☆ Low        | 2023          |

> ⚠️ Frequency data is crowd-sourced from LeetCode Discuss, Glassdoor, and
> community interview reports. Treat as a signal, not a guarantee.

---

## DSA Concepts Used
- **Linked List** — traversal, in-place reversal, fast/slow midpoint → see [`/dsa/linked_list.md`](/dsa/linked_list.md)
- **Two Pointers** — array end-to-middle comparison and the fast/slow midpoint scan → see [`/dsa/two_pointers.md`](/dsa/two_pointers.md)
- **Stack** — the recursive approach uses the implicit call stack as a reverse traversal → see [`/dsa/stack.md`](/dsa/stack.md)

---

## Approaches Overview
| # | Approach | Time | Space | When to use |
|---|----------|------|-------|-------------|
| 1 | Copy to Array + Two Pointers | O(n) | O(n) | Simplest, most readable |
| 2 | Reverse Second Half In Place (Optimal) | O(n) | O(1) | Answers the O(1)-space follow-up |
| 3 | Recursion (Front Pointer) | O(n) | O(n) | Elegant; but O(n) stack risks overflow |

---

## Approach 1 — Copy to Array + Two Pointers

### Intuition
A singly-linked list can't be walked backwards, which is exactly what a palindrome check wants. The simplest fix is to materialise the values into a slice where random access is free, then compare the ends moving inward. Clear and hard to get wrong, at the cost of O(n) extra memory.

### Algorithm
1. Traverse the list once, appending each `Val` to a slice.
2. Set `i = 0`, `j = len-1`.
3. While `i < j`: if `slice[i] != slice[j]` return `false`; else `i++`, `j--`.
4. Return `true`.

### Complexity
- **Time:** O(n) — one traversal plus one two-pointer pass.
- **Space:** O(n) — the values slice.

### Code
```go
func arrayTwoPointers(head *ListNode) bool {
	vals := []int{}
	for node := head; node != nil; node = node.Next { // dump values into a slice
		vals = append(vals, node.Val)
	}
	i, j := 0, len(vals)-1
	for i < j { // compare symmetric ends moving toward the middle
		if vals[i] != vals[j] {
			return false // mismatch → not a palindrome
		}
		i++
		j--
	}
	return true
}
```

### Dry Run
Trace `head = [1,2,2,1]` → `vals = [1,2,2,1]`:

| Step | i | j | vals[i] | vals[j] | match? |
|------|---|---|---------|---------|--------|
| 1    | 0 | 3 | 1       | 1       | yes    |
| 2    | 1 | 2 | 2       | 2       | yes    |
| 3    | 2 | 1 | — (i ≥ j, loop ends) | | |

Return `true`.

---

## Approach 2 — Reverse Second Half In Place (Optimal)

### Intuition
We only need to compare the first half against the reversed second half. A fast/slow pointer finds the midpoint in one pass; reversing the second half in place then lets us walk both halves inward with no auxiliary array. This is the canonical O(n)-time / O(1)-space solution and the answer to the follow-up.

### Algorithm
1. Fast/slow walk: when `fast` reaches the end, `slow` is at the start of the second half.
2. Reverse the sublist starting at `slow`.
3. Walk `first` from head and `second` from the reversed head in lockstep, comparing values until `second` runs out.
4. Return whether all compared values matched.

### Complexity
- **Time:** O(n) — midpoint scan + reversal + comparison, all linear.
- **Space:** O(1) — only a few pointers.

### Code
```go
func reverseHalf(head *ListNode) bool {
	if head == nil || head.Next == nil {
		return true // empty or single node is trivially a palindrome
	}

	slow, fast := head, head
	for fast != nil && fast.Next != nil {
		slow = slow.Next      // advances one node
		fast = fast.Next.Next // advances two nodes
	}

	var prev *ListNode
	cur := slow
	for cur != nil {
		next := cur.Next // remember the rest of the list
		cur.Next = prev  // flip the pointer backward
		prev = cur       // advance prev
		cur = next       // advance cur
	}

	first, second := head, prev
	for second != nil { // second half is the shorter/equal one → drives the loop
		if first.Val != second.Val {
			return false // symmetric values differ
		}
		first = first.Next
		second = second.Next
	}
	return true
}
```

### Dry Run
Trace `head = [1,2,2,1]`:

| Phase | State |
|-------|-------|
| Midpoint | slow starts at head. fast=1→(2's 2nd)→nil. After loop: slow at the 3rd node (value 2, second half = `[2,1]`). |
| Reverse | Reverse `2→1` into `1→2`; `prev` (reversed head) points at value 1. |
| Compare | first=`[1,2,2,1]`, second=`[1,2]` (reversed). |

Comparison loop:

| Step | first.Val | second.Val | match? |
|------|-----------|------------|--------|
| 1    | 1         | 1          | yes    |
| 2    | 2         | 2          | yes    |
| 3    | second == nil → loop ends | | |

Return `true`.

---

## Approach 3 — Recursion (Front Pointer)

### Intuition
Recursion gives us a free reverse traversal: the deepest call sees the last node, and as the stack unwinds we visit nodes from the back. Pairing each unwinding node with a `front` pointer that advances from the head compares position `i` against position `n-1-i`. Elegant, but it uses O(n) stack space, which risks a stack overflow on the 10⁵-node upper bound.

### Algorithm
1. Keep a shared `front` pointer starting at head.
2. `recurse(node)`: if `node` is nil, return `true`. Recurse on `node.Next`; if that sub-call failed, propagate `false`. Compare `node.Val` with `front.Val`, advance `front`, and return the equality.

### Complexity
- **Time:** O(n) — each node visited once.
- **Space:** O(n) — recursion depth equals the list length.

### Code
```go
func recursive(head *ListNode) bool {
	front := head // shared pointer sweeping from the front

	var recurse func(node *ListNode) bool
	recurse = func(node *ListNode) bool {
		if node == nil {
			return true // reached past the tail: base case
		}
		if !recurse(node.Next) { // dive to the end first
			return false // a deeper mismatch already failed
		}
		if node.Val != front.Val { // node unwinds from the back; front from the head
			return false
		}
		front = front.Next // advance the front pointer for the next comparison
		return true
	}
	return recurse(head)
}
```

### Dry Run
Trace `head = [1,2,2,1]`. Recursion descends to the tail, then unwinds. `front` starts at node 0 (value 1):

| Unwind step | node.Val (back) | front.Val (front) | match? | front advances to |
|-------------|-----------------|-------------------|--------|-------------------|
| 1 (last node) | 1             | 1                 | yes    | node 1 (value 2)  |
| 2           | 2               | 2                 | yes    | node 2 (value 2)  |
| 3           | 2               | 2                 | yes    | node 3 (value 1)  |
| 4 (first node) | 1            | 1                 | yes    | nil               |

All matched → return `true`.

---

## Key Takeaways
- **Fast/slow pointers** find the midpoint of a linked list in one pass — the workhorse for "middle of list", "cycle detection", and half-based comparisons.
- **In-place reversal** of the second half turns an O(n)-space problem into O(1) space. Polite implementations restore the list afterward.
- The recursive solution is a neat demonstration that the call stack *is* a stack — but O(n) depth can overflow; prefer the iterative reverse-half in production.
- Comparing a list to its reverse is the linked-list analogue of the array two-pointer palindrome check.

## Related Problems
- LeetCode #206 — Reverse Linked List (the reversal subroutine)
- LeetCode #876 — Middle of the Linked List (the fast/slow midpoint)
- LeetCode #143 — Reorder List (reverse-half + merge)
- LeetCode #9 — Palindrome Number (same idea on integer digits)
