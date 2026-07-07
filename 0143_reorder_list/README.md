# 0143 — Reorder List

> LeetCode #143 · Difficulty: Medium
> **Categories:** Linked List, Two Pointers, Stack, Recursion

---

## Problem Statement

You are given the head of a singly linked-list. The list can be represented as:

```
L0 → L1 → … → Ln-1 → Ln
```

*Reorder the list to be on the following form:*

```
L0 → Ln → L1 → Ln-1 → L2 → Ln-2 → …
```

You may not modify the values in the list's nodes. Only nodes themselves may be changed.

**Example 1:**
```
Input: head = [1,2,3,4]
Output: [1,4,2,3]
```

**Example 2:**
```
Input: head = [1,2,3,4,5]
Output: [1,5,2,4,3]
```

**Constraints:**
- The number of nodes in the list is in the range `[1, 5 * 10⁴]`.
- `1 <= Node.val <= 1000`

---

## Company Frequency

| Company    | Frequency        | Last Reported |
|------------|------------------|---------------|
| Meta       | ★★★★★ Very High  | 2024          |
| Amazon     | ★★★★☆ High       | 2024          |
| Microsoft  | ★★★★☆ High       | 2024          |
| Google     | ★★★☆☆ Medium     | 2024          |
| Bloomberg  | ★★★☆☆ Medium     | 2023          |
| Adobe      | ★★☆☆☆ Low        | 2023          |

> ⚠️ Frequency data is crowd-sourced from LeetCode Discuss, Glassdoor, and
> community interview reports. Treat as a signal, not a guarantee.

---

## DSA Concepts Used

- **Linked List** — the whole problem is pointer surgery: split, reverse, splice → see [`/dsa/linked_list.md`](/dsa/linked_list.md)
- **Two Pointers (Fast & Slow)** — finding the middle in one pass; also the index pair in the array approach → see [`/dsa/two_pointers.md`](/dsa/two_pointers.md)
- **In-place Reversal** — reversing the second half is LeetCode #206 as a subroutine (see the reversal template in [`/dsa/linked_list.md`](/dsa/linked_list.md))

---

## Approaches Overview

| # | Approach | Time | Space | When to use |
|---|----------|------|-------|-------------|
| 1 | Brute Force (Tail Splice) | O(n²) | O(1) | Baseline; no extra memory, tiny lists |
| 2 | Array + Two Pointers | O(n) | O(n) | Quick to write; random access removes all pointer pain |
| 3 | Middle + Reverse + Merge (Optimal) | O(n) | O(1) | The interview answer; combines three classic primitives |

---

## Approach 1 — Brute Force (Tail Splice)

### Intuition
The target order alternates "next node from the front" with "next node from the back". A singly-linked list cannot walk backwards, but it can always *find* the back: walk to the end. So for each front node, walk to the tail, detach it, and splice it right after the front node. The next front is the node right after the freshly spliced tail.

### Algorithm
1. `cur = head`.
2. While `cur != nil && cur.Next != nil`:
   1. Walk `(prev, last)` to the end so `last` is the tail and `prev` its predecessor.
   2. If `last == cur.Next`, the tail is already adjacent to `cur` — the remainder is fully reordered → stop.
   3. Detach: `prev.Next = nil`.
   4. Splice: `last.Next = cur.Next`, then `cur.Next = last`.
   5. Advance past the placed pair: `cur = last.Next`.

### Complexity
- **Time:** O(n²) — about n/2 splices, each rescanning up to n nodes to find the current tail.
- **Space:** O(1) — three pointer variables.

### Code
```go
func bruteForce(head *ListNode) {
	cur := head
	for cur != nil && cur.Next != nil {
		// Walk to the tail, tracking its predecessor for detachment.
		prev, last := cur, cur.Next
		for last.Next != nil {
			prev = last
			last = last.Next
		}
		if last == cur.Next {
			break // tail is directly after cur → nothing left to interleave
		}
		prev.Next = nil      // detach the tail from the list
		last.Next = cur.Next // tail now points at the old front-successor
		cur.Next = last      // front now points at the tail
		cur = last.Next      // skip over the pair (front, tail) just fixed
	}
}
```

### Dry Run (Example 1: head = [1,2,3,4])

| Iteration | `cur` | Tail found (`prev`, `last`) | `last == cur.Next`? | Splice performed | List after |
|-----------|-------|-----------------------------|---------------------|------------------|------------|
| 1 | 1 | prev=3, last=4 | no (cur.Next=2) | 3.Next=nil; 4.Next=2; 1.Next=4 | 1 → 4 → 2 → 3 |
| — | cur = 4.Next = 2 | | | | |
| 2 | 2 | prev=2, last=3 | **yes** (cur.Next=3) | none — break | 1 → 4 → 2 → 3 |

Result: `[1,4,2,3]` ✓

---

## Approach 2 — Array + Two Pointers

### Intuition
Everything hard about this problem is "reach the k-th node from the back in O(1)". Copy the node *pointers* into a slice and you get random access. Then two indices — `i` from the front, `j` from the back — rewire `Next` pointers in exactly the interleaved order, no searching required.

### Algorithm
1. Traverse once, appending every `*ListNode` to `nodes`.
2. `i = 0`, `j = len(nodes) − 1`.
3. While `i < j`:
   1. `nodes[i].Next = nodes[j]` (front links to back), `i++`.
   2. If `i == j`, break (all nodes placed).
   3. `nodes[j].Next = nodes[i]` (back links to next front), `j--`.
4. `nodes[i].Next = nil` — the node where the indices met is the new tail.

### Complexity
- **Time:** O(n) — one pass to collect + one pass of rewiring.
- **Space:** O(n) — the slice holds all n node pointers.

### Code
```go
func arrayTwoPointers(head *ListNode) {
	// Collect pointers for O(1) indexed access.
	var nodes []*ListNode
	for cur := head; cur != nil; cur = cur.Next {
		nodes = append(nodes, cur)
	}

	i, j := 0, len(nodes)-1
	for i < j {
		nodes[i].Next = nodes[j] // front node links to current back node
		i++                      // next front
		if i == j {
			break // pointers met → all nodes placed
		}
		nodes[j].Next = nodes[i] // back node links to next front node
		j--                      // next back
	}
	nodes[i].Next = nil // whoever the pointers met on becomes the new tail
}
```

### Dry Run (Example 1: head = [1,2,3,4])

`nodes = [1, 2, 3, 4]` (by value; entries are pointers).

| Step | `i` | `j` | Action | Links so far |
|------|-----|-----|--------|--------------|
| 1 | 0 | 3 | `nodes[0].Next = nodes[3]` → 1→4; i=1 | 1→4 |
| 2 | 1 | 3 | i ≠ j → `nodes[3].Next = nodes[1]` → 4→2; j=2 | 1→4→2 |
| 3 | 1 | 2 | i < j → `nodes[1].Next = nodes[2]` → 2→3; i=2 | 1→4→2→3 |
| 4 | 2 | 2 | i == j → break | |
| 5 | 2 | 2 | `nodes[2].Next = nil` → 3 is the tail | 1→4→2→3→nil |

Result: `[1,4,2,3]` ✓

---

## Approach 3 — Middle + Reverse + Merge (Optimal)

### Intuition
The reordered list is exactly the first half `[L0..Lmid]` zipped with the *reversed* second half `[Ln..Lmid+1]`. Each piece is a classic O(n)/O(1) primitive:
1. **Find the middle** — fast & slow pointers (LeetCode #876).
2. **Reverse the second half** — iterative reversal (LeetCode #206).
3. **Merge alternately** — zip two lists (like #21's splicing, but alternating unconditionally).

Composing known primitives is the intended "aha" of this problem.

### Algorithm
1. Guard: 0 or 1 node → return.
2. `slow, fast = head, head`; while `fast.Next != nil && fast.Next.Next != nil`: `slow = slow.Next`, `fast = fast.Next.Next`. Now `slow` is the last node of the first half (the first half is never shorter than the second — needed for the zip).
3. Split: `second = slow.Next`; `slow.Next = nil`.
4. Reverse `second` with the `prev/cur/next` rotation; `second = prev` afterwards.
5. Zip: while `second != nil`: save `n1 = first.Next`, `n2 = second.Next`; link `first.Next = second`, `second.Next = n1`; advance `first = n1`, `second = n2`.

### Complexity
- **Time:** O(n) — ~n/2 steps to find the middle + n/2 to reverse + n/2 to merge.
- **Space:** O(1) — all three phases are in-place pointer surgery.

### Code
```go
func middleReverseMerge(head *ListNode) {
	if head == nil || head.Next == nil {
		return // 0 or 1 node → already reordered
	}

	// Step 1: find the end of the first half.
	// Using fast.Next/fast.Next.Next stops slow at ⌈n/2⌉-th node, so the
	// first half is never shorter than the second — required for the zip.
	slow, fast := head, head
	for fast.Next != nil && fast.Next.Next != nil {
		slow = slow.Next
		fast = fast.Next.Next
	}

	// Step 2: split the list after slow.
	second := slow.Next
	slow.Next = nil // terminate the first half

	// Step 3: reverse the second half in place.
	var prev *ListNode
	for cur := second; cur != nil; {
		next := cur.Next // save onward pointer
		cur.Next = prev  // flip the link backwards
		prev = cur       // advance prev
		cur = next       // advance cur
	}
	second = prev // prev is the new head of the reversed half

	// Step 4: zip-merge the two halves, alternating first → second.
	first := head
	for second != nil {
		n1, n2 := first.Next, second.Next // save both onward pointers
		first.Next = second               // front node → back node
		second.Next = n1                  // back node → next front node
		first, second = n1, n2            // advance both halves
	}
}
```

### Dry Run (Example 1: head = [1,2,3,4])

**Step 1 — find middle** (`slow=1`, `fast=1`):

| Turn | Guard (`fast.Next`, `fast.Next.Next`) | `slow` | `fast` |
|------|----------------------------------------|--------|--------|
| 1 | 2, 3 — both non-nil → move | 2 | 3 |
| 2 | 4, nil — stop | 2 | 3 |

**Step 2 — split:** `second = 3`, `2.Next = nil`. Halves: `1→2` and `3→4`.

**Step 3 — reverse second half:**

| Iter | `cur` | `next` | Link set | `prev` after |
|------|-------|--------|----------|--------------|
| 1 | 3 | 4 | 3.Next = nil | 3 |
| 2 | 4 | nil | 4.Next = 3 | 4 |

`second = 4 → 3`.

**Step 4 — zip** (`first = 1`, `second = 4`):

| Iter | `n1` | `n2` | Links set | List so far | `first` | `second` |
|------|------|------|-----------|-------------|---------|----------|
| 1 | 2 | 3 | 1.Next=4; 4.Next=2 | 1→4→2 | 2 | 3 |
| 2 | nil | nil | 2.Next=3; 3.Next=nil | 1→4→2→3 | nil | nil |

`second == nil` → done. Result: `[1,4,2,3]` ✓

---

## Key Takeaways

- **Decompose into known primitives**: middle-of-list + reverse-list + zip-merge. Recognising a problem as a composition is a reusable interview superpower.
- The middle-finding guard `fast.Next != nil && fast.Next.Next != nil` makes `slow` stop at ⌈n/2⌉, guaranteeing the first half is ≥ the second — the zip loop can then simply run until `second` is exhausted.
- When zipping, **save both onward pointers first** (`n1`, `n2`) before overwriting any `Next` — the #1 source of bugs in pointer-splicing code.
- A slice of node pointers converts any "k-th from the back" linked-list problem into trivial array indexing — the standard time/space trade-off.
- "You may not modify values, only nodes" is the problem telling you value-swapping is cheating; expect the same clause in other pointer-surgery questions.

---

## Related Problems

- LeetCode #876 — Middle of the Linked List (step 1 as a standalone problem)
- LeetCode #206 — Reverse Linked List (step 3 as a standalone problem)
- LeetCode #21 — Merge Two Sorted Lists (splice-style merging)
- LeetCode #234 — Palindrome Linked List (identical middle+reverse skeleton, compare instead of merge)
- LeetCode #61 — Rotate List (tail-to-front pointer surgery)
- LeetCode #24 — Swap Nodes in Pairs (local pointer rewiring drills)
