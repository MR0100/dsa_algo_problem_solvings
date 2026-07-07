# 0369 — Plus One Linked List

> LeetCode #369 · Difficulty: Medium
> **Categories:** Linked List, Math

---

## Problem Statement

Given a non-negative integer represented as a linked list of digits, *plus one* to the integer.

The digits are stored such that the most significant digit is at the `head` of the list.

**Example 1:**

```
Input: head = [1,2,3]
Output: [1,2,4]
```

**Example 2:**

```
Input: head = [0]
Output: [1]
```

**Constraints:**

- The number of nodes in the linked list is in the range `[1, 100]`.
- `0 <= Node.val <= 9`
- The number represented by the linked list does not contain leading zeros except for the number `0` itself.

---

## Company Frequency

| Company    | Frequency        | Last Reported |
|------------|------------------|---------------|
| Google     | ★★★☆☆ Medium     | 2023          |
| Amazon     | ★★★☆☆ Medium     | 2023          |
| Microsoft  | ★★☆☆☆ Low        | 2022          |
| Facebook   | ★★☆☆☆ Low        | 2022          |

> ⚠️ Frequency data is crowd-sourced from LeetCode Discuss, Glassdoor, and
> community interview reports. Treat as a signal, not a guarantee.

---

## DSA Concepts Used

- **Linked List traversal & pointer surgery** — reversing, re-linking, and prepending nodes to model digit arithmetic → see [`/dsa/linked_list.md`](/dsa/linked_list.md)
- **Elementary number theory (carry propagation)** — adding one only ripples a carry through a trailing run of 9s → see [`/dsa/math_number_theory.md`](/dsa/math_number_theory.md)

---

## Approaches Overview

| # | Approach | Time | Space | When to use |
|---|----------|------|-------|-------------|
| 1 | Reverse, Add, Reverse Back (Brute Force) | O(n) | O(1) | Intuitive: make LSB reachable first |
| 2 | Rightmost Non-Nine (Optimal, No Reverse) | O(n) | O(1) | Cleanest single-forward-pass answer |
| 3 | Recursion Carry-Back (Optimal) | O(n) | O(n) | Elegant; carry rides the call stack |

---

## Approach 1 — Reverse, Add, Reverse Back (Brute Force)

### Intuition

Addition starts at the least-significant digit, but the list stores the most-significant digit first. Reverse the list to put the ones digit at the head, propagate the `+1` carry left-to-right, then reverse back. If a carry survives past the front (e.g. `99 → 100`), append a new node.

### Algorithm

1. Reverse the list.
2. Walk it with `carry` starting at 1: `sum = val + carry`; store `sum % 10`; `carry = sum / 10`.
3. If `carry` remains after the last node, append a node holding it.
4. Reverse back and return the head.

### Complexity

- **Time:** O(n) — three linear passes.
- **Space:** O(1) — in-place pointer surgery.

### Code

```go
func reverseAddReverse(head *ListNode) *ListNode {
	reverse := func(node *ListNode) *ListNode {
		var prev *ListNode
		for node != nil {
			node.Next, prev, node = prev, node, node.Next // classic 3-way swap
		}
		return prev
	}

	head = reverse(head) // least-significant digit now at the head
	carry := 1           // the "+1" we must add
	cur := head
	var tail *ListNode // remember the last node to append an overflow digit
	for cur != nil {
		sum := cur.Val + carry
		cur.Val = sum % 10  // keep the low digit
		carry = sum / 10    // propagate the carry (0 or 1)
		tail = cur
		cur = cur.Next
	}
	if carry > 0 { // overflow past the most-significant digit (e.g. 99 → 100)
		tail.Next = &ListNode{Val: carry}
	}
	return reverse(head) // restore most-significant-first order
}
```

### Dry Run

Example 1: `head = [1,2,3]`.

| Step | Action | List state |
|------|--------|------------|
| 1 | reverse | `3 → 2 → 1` |
| 2 | 3+1=4, carry 0 | `4 → 2 → 1` |
| 3 | 2+0=2, carry 0 | `4 → 2 → 1` |
| 4 | 1+0=1, carry 0 | `4 → 2 → 1` |
| 5 | carry 0, no append | `4 → 2 → 1` |
| 6 | reverse back | `1 → 2 → 4` |

Result: `[1,2,4]` ✔

---

## Approach 2 — Rightmost Non-Nine (Optimal, No Reverse)

### Intuition

Adding 1 to a number only ripples a carry through a trailing run of 9s. The **last non-9 digit** is where the carry stops: bump it by one and turn every 9 to its right into 0. If *every* digit is 9, the whole number rolls over — so we plant a sentinel (dummy) node in front whose 0 acts as "the last non-nine": incrementing it produces the leading 1 (e.g. `999 → 1000`), and we return the dummy in that case.

### Algorithm

1. Create `dummy → head`; set `lastNotNine = dummy`.
2. Walk the real nodes; whenever a value `!= 9`, update `lastNotNine`.
3. Increment `lastNotNine.Val`.
4. Set every node after `lastNotNine` to 0.
5. Return `dummy.Next` if `dummy` stayed 0, else `dummy` (a leading 1 was created).

### Complexity

- **Time:** O(n) — one pass to locate, a partial pass to zero out.
- **Space:** O(1).

### Code

```go
func rightmostNonNine(head *ListNode) *ListNode {
	dummy := &ListNode{Val: 0, Next: head} // sentinel absorbs an all-nines carry
	lastNotNine := dummy                   // rightmost node whose value != 9
	for node := head; node != nil; node = node.Next {
		if node.Val != 9 {
			lastNotNine = node
		}
	}
	lastNotNine.Val++ // the carry lands here and stops
	// Everything to the right of the increment point becomes 0.
	for node := lastNotNine.Next; node != nil; node = node.Next {
		node.Val = 0
	}
	if dummy.Val == 0 { // no rollover ⇒ dummy is unused
		return dummy.Next
	}
	return dummy // dummy became the new leading 1 (e.g. 999 → 1000)
}
```

### Dry Run

Example 1: `head = [1,2,3]`.

| Step | node | val != 9? | lastNotNine |
|------|------|-----------|-------------|
| scan | 1 | yes | node(1) |
| scan | 2 | yes | node(2) |
| scan | 3 | yes | node(3) |
| — | increment lastNotNine(3) → 4 | | list `1→2→4` |
| — | nodes after → 0 (none) | | list `1→2→4` |
| — | dummy.Val==0 → return dummy.Next | | `[1,2,4]` |

Result: `[1,2,4]` ✔ (Edge `[9,9]`: no non-nine among reals, so `lastNotNine` stays dummy → dummy becomes 1, both 9s zeroed → `[1,0,0]`.)

---

## Approach 3 — Recursion Carry-Back (Optimal)

### Intuition

Recursion reaches the least-significant digit last, so we add the carry as the stack unwinds — exactly the right-to-left order arithmetic needs, with no reversal. The base case (past the tail) returns carry 1; each node adds it, keeps the low digit, and returns the new carry upward. A carry surviving at the head means prepend a 1.

### Algorithm

1. `add(node)`: if `node == nil` return 1 (the injected `+1`).
2. `carry = add(node.Next)`; `sum = node.Val + carry`.
3. `node.Val = sum % 10`; return `sum / 10`.
4. In the caller: if `add(head) == 1`, prepend a node with value 1.

### Complexity

- **Time:** O(n) — one node per stack frame.
- **Space:** O(n) — recursion depth.

### Code

```go
func recursionCarry(head *ListNode) *ListNode {
	var add func(node *ListNode) int
	add = func(node *ListNode) int {
		if node == nil {
			return 1 // reached the end; this is the +1 to inject
		}
		carry := add(node.Next) // finish the lower digits first
		sum := node.Val + carry
		node.Val = sum % 10 // keep low digit
		return sum / 10     // pass carry up
	}
	if add(head) == 1 { // carry escaped the most-significant digit
		return &ListNode{Val: 1, Next: head}
	}
	return head
}
```

### Dry Run

Example 1: `head = [1,2,3]`. Frames unwind from the tail.

| Unwind | node | carry in | sum | node.Val | carry out |
|--------|------|----------|-----|----------|-----------|
| base | nil | — | — | — | 1 |
| 3 | 3 | 1 | 4 | 4 | 0 |
| 2 | 2 | 0 | 2 | 2 | 0 |
| 1 | 1 | 0 | 1 | 1 | 0 |

`add(head)` returns 0 → no prepend. List `1 → 2 → 4`. Result: `[1,2,4]` ✔

---

## Key Takeaways

- **Carry only ripples through trailing 9s.** Finding the rightmost non-nine, incrementing it, and zeroing the rest is the whole problem — a single forward pass, no reversal.
- **Sentinel/dummy node handles rollover uniformly.** Its 0 becomes the new leading 1 when every digit is 9, avoiding a special case.
- **Recursion supplies free right-to-left order** for LSB-first arithmetic on an MSB-first list.
- Reverse-add-reverse is the fallback when you can't or don't want to think in sentinels — same O(1) space, more passes.

---

## Related Problems

- LeetCode #66 — Plus One (array version of the exact same carry logic)
- LeetCode #2 — Add Two Numbers (LSB-first addition of two lists)
- LeetCode #445 — Add Two Numbers II (MSB-first; reverse or stack)
- LeetCode #206 — Reverse Linked List (the reversal primitive used here)
