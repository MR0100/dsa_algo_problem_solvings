# 0206 — Reverse Linked List

> LeetCode #206 · Difficulty: Easy
> **Categories:** Linked List, Recursion, Two Pointers, Stack

---

## Problem Statement

Given the `head` of a singly linked list, reverse the list, and return *the reversed list*.

**Example 1:**
```
Input: head = [1,2,3,4,5]
Output: [5,4,3,2,1]
```

**Example 2:**
```
Input: head = [1,2]
Output: [2,1]
```

**Example 3:**
```
Input: head = []
Output: []
```

**Constraints:**
- The number of nodes in the list is the range `[0, 5000]`.
- `-5000 <= Node.val <= 5000`

**Follow-up:** A linked list can be reversed either iteratively or recursively. Could you implement both?

---

## Company Frequency

| Company    | Frequency        | Last Reported |
|------------|------------------|---------------|
| Amazon     | ★★★★★ Very High  | 2024          |
| Microsoft  | ★★★★★ Very High  | 2024          |
| Google     | ★★★★☆ High       | 2024          |
| Meta       | ★★★★☆ High       | 2024          |
| Apple      | ★★★☆☆ Medium     | 2024          |
| Bloomberg  | ★★★☆☆ Medium     | 2023          |
| Adobe      | ★★★☆☆ Medium     | 2023          |
| Oracle     | ★★☆☆☆ Low        | 2023          |

> ⚠️ Frequency data is crowd-sourced from LeetCode Discuss, Glassdoor, and
> community interview reports. Treat as a signal, not a guarantee.

---

## DSA Concepts Used

- **Linked List** — the whole problem is pointer manipulation on a singly linked list; reversal is *the* foundational list operation reused by dozens of harder problems → see [`/dsa/linked_list.md`](/dsa/linked_list.md)
- **Two Pointers** — the optimal solution walks a `prev`/`cur` pointer pair down the list, flipping one arrow per step → see [`/dsa/two_pointers.md`](/dsa/two_pointers.md)
- **Stack** — LIFO order is reversal by definition; pushing all nodes and popping them re-links the list backwards (and the recursive solution is this same stack, hidden in the call frames) → see [`/dsa/stack.md`](/dsa/stack.md)

---

## Approaches Overview

| # | Approach | Time | Space | When to use |
|---|----------|------|-------|-------------|
| 1 | Brute Force (Copy Values & Rebuild) | O(n) | O(n) | Never in an interview — but proves reversal needs no pointer surgery; useful when input must stay immutable |
| 2 | Stack of Nodes | O(n) | O(n) | Stepping stone: makes the "reverse = LIFO" insight explicit before optimising |
| 3 | Recursion | O(n) | O(n) | The follow-up's second half; elegant, and the sub-structure trick powers #92/#25 |
| 4 | Iterative Pointer Reversal (Optimal) | O(n) | O(1) | Default interview answer; constant space, one pass |

---

## Approach 1 — Brute Force (Copy Values & Rebuild)

### Intuition
If pointer surgery feels risky, sidestep it entirely. A linked list is just a sequence of values, so record the sequence in a slice, then manufacture a brand-new list that emits those values back to front. No existing pointer is ever modified, which makes correctness trivial — at the cost of O(n) extra memory and all-new nodes.

### Algorithm
1. Walk the list once, appending each node's value to a slice `vals`.
2. Create a sentinel (`dummy`) node and a `tail` pointer at the sentinel.
3. Iterate `vals` from index `len(vals)-1` down to `0`; for each value append a new node at `tail.Next` and advance `tail`.
4. Return `dummy.Next` — the head of the freshly built reversed copy.

### Complexity
- **Time:** O(n) — one pass to collect n values plus one pass to build n nodes.
- **Space:** O(n) — the value slice and the n newly allocated nodes.

### Code
```go
func bruteForce(head *ListNode) *ListNode {
	// Collect the values in original order.
	vals := []int{}
	for cur := head; cur != nil; cur = cur.Next {
		vals = append(vals, cur.Val) // remember each value as we pass it
	}
	// Rebuild a completely new list, reading the slice backwards.
	dummy := &ListNode{} // sentinel so the first append needs no special case
	tail := dummy        // tail always points at the last node built so far
	for i := len(vals) - 1; i >= 0; i-- {
		tail.Next = &ListNode{Val: vals[i]} // append a fresh node holding vals[i]
		tail = tail.Next                    // advance tail onto the new node
	}
	return dummy.Next // real head sits right after the sentinel
}
```

### Dry Run (Example 1: head = [1,2,3,4,5])

Phase 1 — collect values:

| Step | `cur` | `vals` after append |
|------|-------|---------------------|
| 1 | node(1) | [1] |
| 2 | node(2) | [1 2] |
| 3 | node(3) | [1 2 3] |
| 4 | node(4) | [1 2 3 4] |
| 5 | node(5) | [1 2 3 4 5] |
| 6 | nil — loop ends | [1 2 3 4 5] |

Phase 2 — rebuild backwards (`i` from 4 down to 0):

| `i` | `vals[i]` | List built so far (after dummy) |
|-----|-----------|---------------------------------|
| 4 | 5 | 5 |
| 3 | 4 | 5 → 4 |
| 2 | 3 | 5 → 4 → 3 |
| 1 | 2 | 5 → 4 → 3 → 2 |
| 0 | 1 | 5 → 4 → 3 → 2 → 1 |

Return `dummy.Next` = **[5,4,3,2,1]** ✓

---

## Approach 2 — Stack of Nodes

### Intuition
A stack reverses order by nature: the last node pushed is the first node popped. Push every node pointer during one traversal, then pop them one by one, wiring each popped node's `Next` to the following pop. Unlike Approach 1 this reuses the original nodes — only their `Next` pointers change — and it makes the "reversal = LIFO" insight explicit. The recursive solution (Approach 3) is exactly this stack, hidden inside the call frames.

### Algorithm
1. Traverse the list, pushing each `*ListNode` onto a slice-backed stack.
2. If the stack is empty, return `nil` (empty input).
3. Pop the top node (the original tail) — it becomes `newHead`.
4. While nodes remain, wire the previously popped node's `Next` to the next pop and advance.
5. Set the last popped node's `Next` (the original head) to `nil` to terminate the list.

### Complexity
- **Time:** O(n) — every node is pushed exactly once and popped exactly once.
- **Space:** O(n) — the stack holds all n node pointers at its peak.

### Code
```go
func stackBased(head *ListNode) *ListNode {
	// Push every node onto the stack in traversal order.
	stack := []*ListNode{}
	for cur := head; cur != nil; cur = cur.Next {
		stack = append(stack, cur) // node pointers, not values — we reuse nodes
	}
	if len(stack) == 0 {
		return nil // empty list reverses to the empty list
	}
	// The last node pushed (original tail) is the new head.
	newHead := stack[len(stack)-1]
	cur := newHead
	// Pop the remaining nodes, chaining each onto the reversed list.
	for i := len(stack) - 2; i >= 0; i-- {
		cur.Next = stack[i] // wire current node to the next-popped node
		cur = cur.Next      // advance along the growing reversed list
	}
	cur.Next = nil // the original head is now the tail — terminate it
	return newHead
}
```

### Dry Run (Example 1: head = [1,2,3,4,5])

Nodes named by value: N1…N5. After the push phase: `stack = [N1 N2 N3 N4 N5]`, `newHead = N5`.

| `i` | `stack[i]` | Wire made | Reversed chain so far |
|-----|-----------|-----------|------------------------|
| — | — | — | N5 |
| 3 | N4 | N5.Next = N4 | N5 → N4 |
| 2 | N3 | N4.Next = N3 | N5 → N4 → N3 |
| 1 | N2 | N3.Next = N2 | N5 → N4 → N3 → N2 |
| 0 | N1 | N2.Next = N1 | N5 → N4 → N3 → N2 → N1 |

Finally `N1.Next = nil`. Return `N5` = **[5,4,3,2,1]** ✓

---

## Approach 3 — Recursion

### Intuition
Trust the recursion to reverse everything *after* `head`. Once `recursive(head.Next)` returns, the picture is `head → (tail … head.Next)` — i.e. `head` still points forward at what is now the **last** node of the reversed sublist. So `head.Next.Next = head` attaches `head` to the end, and `head.Next = nil` terminates it. The deepest call returns the original tail, and that pointer is passed up unchanged — it is the new head.

### Algorithm
1. Base case: if `head` is `nil` or `head.Next` is `nil`, the list is its own reversal — return `head`.
2. Recurse: `newHead := recursive(head.Next)` — reverses the whole sublist after `head`.
3. `head.Next` is now the reversed sublist's tail, so set `head.Next.Next = head` to append `head` behind it.
4. Set `head.Next = nil` — `head` is the current tail (a caller one level up may re-point it).
5. Return `newHead` unchanged.

### Complexity
- **Time:** O(n) — exactly one call per node, O(1) work per call.
- **Space:** O(n) — the call stack reaches n frames before unwinding.

### Code
```go
func recursive(head *ListNode) *ListNode {
	// Base case: nothing to reverse for an empty or single-node list.
	if head == nil || head.Next == nil {
		return head
	}
	// Reverse everything after head; newHead is the original tail.
	newHead := recursive(head.Next)
	// head.Next is the last node of the reversed sublist — hook head behind it.
	head.Next.Next = head
	// Break head's old forward pointer; it now ends the list (or gets
	// re-pointed by the caller one level up).
	head.Next = nil
	return newHead
}
```

### Dry Run (Example 1: head = [1,2,3,4,5])

Descent: calls stack up as `recursive(N1) → recursive(N2) → recursive(N3) → recursive(N4) → recursive(N5)`; `recursive(N5)` hits the base case and returns N5. Unwinding:

| Returning frame | `head` | `newHead` | After `head.Next.Next = head` and `head.Next = nil` | List state |
|-----------------|--------|-----------|------------------------------------------------------|------------|
| recursive(N4) | N4 | N5 | N5.Next = N4, N4.Next = nil | 5 → 4 (1 → 2 → 3 → 4 detached) |
| recursive(N3) | N3 | N5 | N4.Next = N3, N3.Next = nil | 5 → 4 → 3 |
| recursive(N2) | N2 | N5 | N3.Next = N2, N2.Next = nil | 5 → 4 → 3 → 2 |
| recursive(N1) | N1 | N5 | N2.Next = N1, N1.Next = nil | 5 → 4 → 3 → 2 → 1 |

Top-level call returns N5 = **[5,4,3,2,1]** ✓

---

## Approach 4 — Iterative Pointer Reversal (Optimal)

### Intuition
Reversing a list is literally "make every arrow point the other way". Sweep once with two pointers: `prev` is the head of the already-reversed prefix, `cur` is the first untouched node. Each step redirects `cur.Next` from its successor to `prev` — but that overwrites the only route forward, so a temporary `next` saves the successor first. When `cur` runs off the end, `prev` holds the original tail: the new head.

### Algorithm
1. Initialise `prev = nil` (reversed prefix is empty) and `cur = head`.
2. While `cur != nil`:
   1. `next := cur.Next` — save the untouched suffix before overwriting.
   2. `cur.Next = prev` — flip the arrow backwards.
   3. `prev = cur` — the reversed prefix grows by one node.
   4. `cur = next` — step into the suffix.
3. Return `prev`.

### Complexity
- **Time:** O(n) — each node is visited exactly once with constant work.
- **Space:** O(1) — three pointers (`prev`, `cur`, `next`) regardless of length.

### Code
```go
func iterative(head *ListNode) *ListNode {
	var prev *ListNode // head of the already-reversed prefix (starts empty)
	cur := head        // first node not yet reversed
	for cur != nil {
		next := cur.Next // save successor — cur.Next is about to be overwritten
		cur.Next = prev  // flip: cur now points backwards at the reversed prefix
		prev = cur       // reversed prefix grows to include cur
		cur = next       // advance into the untouched suffix
	}
	// cur is nil, so prev holds the original tail — the new head.
	return prev
}
```

### Dry Run (Example 1: head = [1,2,3,4,5])

| Step | `cur` (before) | `next` saved | Flip performed | `prev` (after) | `cur` (after) | Reversed prefix |
|------|----------------|--------------|----------------|----------------|----------------|-----------------|
| init | N1 | — | — | nil | N1 | (empty) |
| 1 | N1 | N2 | N1.Next = nil | N1 | N2 | 1 |
| 2 | N2 | N3 | N2.Next = N1 | N2 | N3 | 2 → 1 |
| 3 | N3 | N4 | N3.Next = N2 | N3 | N4 | 3 → 2 → 1 |
| 4 | N4 | N5 | N4.Next = N3 | N4 | N5 | 4 → 3 → 2 → 1 |
| 5 | N5 | nil | N5.Next = N4 | N5 | nil | 5 → 4 → 3 → 2 → 1 |

Loop exits (`cur == nil`); return `prev` = N5 = **[5,4,3,2,1]** ✓

---

## Key Takeaways

- **The three-pointer dance (`prev`, `cur`, `next`) is a core primitive.** Memorise it cold: save `next`, flip `cur.Next`, advance both. It reappears inside #92 (reverse a sublist), #25 (reverse k-groups), #234 (palindrome check), and #143 (reorder list).
- **Always save the successor before overwriting `Next`** — the single most common linked-list bug is losing the rest of the list mid-flip.
- **Recursion trick:** after `recursive(head.Next)` returns, `head.Next` points at the reversed sublist's *tail*, so `head.Next.Next = head` appends `head`. Understanding *why* that line works is the real test of the recursive version.
- **Stack ⇄ recursion equivalence:** the explicit stack version and the recursive version are the same algorithm; recursion just hides the stack in call frames. Converting between the two is a general skill.
- **`nil` handles itself:** with `prev` starting at `nil` (iterative) or `head == nil` as a base case (recursive), the empty list needs no special-casing — a hallmark of clean pointer code.
- Iterative wins in production: O(1) space and no risk of stack overflow on a 5000-node (or 5-million-node) list.

---

## Related Problems

- LeetCode #92 — Reverse Linked List II (reverse only positions left..right, same flip primitive)
- LeetCode #25 — Reverse Nodes in k-Group (repeated bounded reversal)
- LeetCode #24 — Swap Nodes in Pairs (k-group reversal with k = 2)
- LeetCode #234 — Palindrome Linked List (reverse the second half, then compare)
- LeetCode #143 — Reorder List (reverse the back half, then interleave)
- LeetCode #61 — Rotate List (pointer rewiring on a singly linked list)
