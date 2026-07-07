# 0203 — Remove Linked List Elements

> LeetCode #203 · Difficulty: Easy
> **Categories:** Linked List, Recursion

---

## Problem Statement

Given the `head` of a linked list and an integer `val`, remove all the nodes of the linked list that has `Node.val == val`, and return *the new head*.

**Example 1:**

```
Input: head = [1,2,6,3,4,5,6], val = 6
Output: [1,2,3,4,5]
```

**Example 2:**

```
Input: head = [], val = 1
Output: []
```

**Example 3:**

```
Input: head = [7,7,7,7], val = 7
Output: []
```

**Constraints:**

- The number of nodes in the list is in the range `[0, 10^4]`.
- `1 <= Node.val <= 50`
- `0 <= val <= 50`

---

## Company Frequency

| Company    | Frequency        | Last Reported |
|------------|------------------|---------------|
| Amazon     | ★★★★☆ High       | 2024          |
| Microsoft  | ★★★☆☆ Medium     | 2024          |
| Google     | ★★★☆☆ Medium     | 2023          |
| Apple      | ★★★☆☆ Medium     | 2023          |
| Bloomberg  | ★★☆☆☆ Low        | 2023          |
| Adobe      | ★★☆☆☆ Low        | 2022          |

> ⚠️ Frequency data is crowd-sourced from LeetCode Discuss, Glassdoor, and
> community interview reports. Treat as a signal, not a guarantee.

---

## DSA Concepts Used

- **Linked List — dummy (sentinel) head** — deletions need the node *before* the victim; a sentinel gives every real node a predecessor and erases the head special-case → see [`/dsa/linked_list.md`](/dsa/linked_list.md)
- **Linked List — in-place pointer surgery** — `prev.Next = prev.Next.Next` unlinks a node in O(1) without moving any data → see [`/dsa/linked_list.md`](/dsa/linked_list.md)

---

## Approaches Overview

| # | Approach | Time | Space | When to use |
|---|----------|------|-------|-------------|
| 1 | Iterative Without Dummy | O(n) | O(1) | Shows the raw problem: the head is the only special case |
| 2 | Recursion | O(n) | O(n) stack | Elegant 5-liner; fine for short lists, risky at 10⁴+ depth |
| 3 | Dummy Node / Sentinel (Optimal) | O(n) | O(1) | Always — one uniform loop, no special cases |

---

## Approach 1 — Iterative Without Dummy

### Intuition

To delete a node from a singly-linked list you must stand on the node **before** it — that's the only place from which its incoming pointer can be rewired. Every node has a predecessor except the head. So solve the problem in two phases: first slide `head` forward over any leading nodes with the target value (they simply stop being reachable), then do the standard walk in which `cur` always sits on a *kept* node and inspects `cur.Next`, splicing it out when it matches. This makes the structural pain of the problem explicit — the head is different — which is exactly what the sentinel of Approach 3 later abstracts away.

### Algorithm

1. While `head != nil` **and** `head.Val == val`: advance `head = head.Next` (drops all leading matches, possibly the entire list).
2. Set `cur = head`.
3. While `cur != nil` and `cur.Next != nil`:
   1. If `cur.Next.Val == val`: bypass it with `cur.Next = cur.Next.Next` (do **not** advance — the new `cur.Next` is unchecked).
   2. Else: advance `cur = cur.Next`.
4. Return `head`.

### Complexity

- **Time:** O(n) — each node is looked at exactly once, either skipped over or stepped onto.
- **Space:** O(1) — only the `cur` pointer; all relinking is in place.

### Code

```go
func iterativeWithoutDummy(head *ListNode, val int) *ListNode {
	// Phase 1: the head has no predecessor — peel off matching heads first.
	for head != nil && head.Val == val {
		head = head.Next // old head becomes garbage; next node is new head
	}
	// Phase 2: cur always sits on a KEPT node, so cur.Next is deletable.
	cur := head
	for cur != nil && cur.Next != nil {
		if cur.Next.Val == val {
			cur.Next = cur.Next.Next // splice the matching node out
		} else {
			cur = cur.Next // next node survives; step onto it
		}
	}
	return head
}
```

### Dry Run

Example 1: `head = [1,2,6,3,4,5,6], val = 6`.

Phase 1: `head.Val = 1 ≠ 6` → no leading matches, head stays at node `1`.

| Step | cur | cur.Next | cur.Next.Val == 6? | Action | List state |
|------|-----|----------|--------------------|--------|------------|
| 1 | 1 | 2 | no | advance cur → 2 | 1→2→6→3→4→5→6 |
| 2 | 2 | 6 | yes | splice: 2.Next = 3 | 1→2→3→4→5→6 |
| 3 | 2 | 3 | no | advance cur → 3 | 1→2→3→4→5→6 |
| 4 | 3 | 4 | no | advance cur → 4 | 1→2→3→4→5→6 |
| 5 | 4 | 5 | no | advance cur → 5 | 1→2→3→4→5→6 |
| 6 | 5 | 6 | yes | splice: 5.Next = nil | 1→2→3→4→5 |
| 7 | 5 | nil | — | exit loop | 1→2→3→4→5 |

Result: `[1,2,3,4,5]` ✔ (Example 3 `[7,7,7,7]` is consumed entirely by Phase 1: head slides off all four 7s to `nil`, Phase 2 never runs, return `nil`.)

---

## Approach 2 — Recursion

### Intuition

A list with all `val` nodes removed is: *the cleaned rest of the list*, preceded by *this node if it survives*. That self-similarity is a recursion. Let the recursive call fully clean everything after the current node, hook that cleaned suffix on with `head.Next = ...`, then decide the current node's fate: if it matches, the node deletes itself by returning the cleaned suffix instead of itself; otherwise it returns itself. Head deletions need no special handling — a deleted head simply hands back its suffix, and the caller (or `main`) receives the correct new head.

### Algorithm

1. Base case: `head == nil` → return `nil` (an empty list is already clean).
2. Recurse: `head.Next = recursiveApproach(head.Next, val)` — after this line, everything behind `head` is guaranteed clean.
3. If `head.Val == val`: return `head.Next` (self-delete by bypassing).
4. Else: return `head`.

### Complexity

- **Time:** O(n) — exactly one call per node, O(1) work each.
- **Space:** O(n) — the call stack is as deep as the list; at the constraint maximum of 10⁴ nodes this is safe in Go but the pattern fails on very long lists (stack overflow risk in other languages).

### Code

```go
func recursiveApproach(head *ListNode, val int) *ListNode {
	// Base case: nothing to remove in an empty list.
	if head == nil {
		return nil
	}
	// Recursively clean everything after this node first.
	head.Next = recursiveApproach(head.Next, val)
	// Now decide this node's fate: matching nodes vanish by returning
	// the cleaned suffix instead of themselves.
	if head.Val == val {
		return head.Next
	}
	return head
}
```

### Dry Run

Example 1: `head = [1,2,6,3,4,5,6], val = 6`. Calls descend to the end, then resolve back-to-front:

| Step (unwind order) | Node | Cleaned suffix received | Node.Val == 6? | Returns |
|---------------------|------|-------------------------|-----------------|---------|
| 1 | nil | — | — | nil (base case) |
| 2 | 6 (last) | nil | yes | nil — deletes itself |
| 3 | 5 | nil | no | 5 → nil |
| 4 | 4 | 5 | no | 4→5 |
| 5 | 3 | 4→5 | no | 3→4→5 |
| 6 | 6 (third) | 3→4→5 | yes | 3→4→5 — deletes itself |
| 7 | 2 | 3→4→5 | no | 2→3→4→5 |
| 8 | 1 | 2→3→4→5 | no | 1→2→3→4→5 |

Result: `[1,2,3,4,5]` ✔

---

## Approach 3 — Dummy Node / Sentinel (Optimal)

### Intuition

Approach 1's two-phase structure exists solely because the head lacks a predecessor. So *manufacture a predecessor*: allocate one sentinel node `dummy` with `dummy.Next = head`. Now every real node — including the head — has a node before it, and a single uniform loop handles every deletion, including "delete every node in the list". The invariant is crisp: `prev` always points at the last node we have **decided to keep** (initially the dummy, which is kept by definition), and `prev.Next` is the next node awaiting judgment. At the end, `dummy.Next` is the true head, whatever happened to the original one. This is the #1 linked-list interview trick and the version to write by default.

### Algorithm

1. Create `dummy = &ListNode{Next: head}`; set `prev = dummy`.
2. While `prev.Next != nil`:
   1. If `prev.Next.Val == val`: unlink it with `prev.Next = prev.Next.Next` (stay on `prev` — the new `prev.Next` still needs judgment).
   2. Else: the node survives; advance `prev = prev.Next`.
3. Return `dummy.Next`.

### Complexity

- **Time:** O(n) — one pass; each node is judged exactly once.
- **Space:** O(1) — a single extra node regardless of input size.

### Code

```go
func dummyNode(head *ListNode, val int) *ListNode {
	dummy := &ListNode{Next: head} // sentinel sits before the real head
	prev := dummy                  // prev is always the last KEPT node
	for prev.Next != nil {
		if prev.Next.Val == val {
			prev.Next = prev.Next.Next // unlink the matching node
		} else {
			prev = prev.Next // keep it; move the frontier forward
		}
	}
	return dummy.Next // real head, even if the original head was deleted
}
```

### Dry Run

Example 1: `head = [1,2,6,3,4,5,6], val = 6`. List with sentinel: `D→1→2→6→3→4→5→6`.

| Step | prev | prev.Next | match 6? | Action | List after |
|------|------|-----------|----------|--------|------------|
| 1 | D | 1 | no | advance prev → 1 | D→1→2→6→3→4→5→6 |
| 2 | 1 | 2 | no | advance prev → 2 | D→1→2→6→3→4→5→6 |
| 3 | 2 | 6 | yes | unlink: 2.Next = 3 | D→1→2→3→4→5→6 |
| 4 | 2 | 3 | no | advance prev → 3 | D→1→2→3→4→5→6 |
| 5 | 3 | 4 | no | advance prev → 4 | D→1→2→3→4→5→6 |
| 6 | 4 | 5 | no | advance prev → 5 | D→1→2→3→4→5→6 |
| 7 | 5 | 6 | yes | unlink: 5.Next = nil | D→1→2→3→4→5 |
| 8 | 5 | nil | — | exit loop | D→1→2→3→4→5 |

Return `dummy.Next` = node `1`. Result: `[1,2,3,4,5]` ✔ (Example 3 `[7,7,7,7]`: prev never advances off the dummy — all four nodes are unlinked one by one and `dummy.Next` ends as `nil` → `[]`.)

---

## Key Takeaways

- **Sentinel/dummy head is the default weapon** whenever the head itself might be deleted or replaced (#19, #82, #83, #86, #92): allocate `dummy → head`, work with a `prev` pointer, return `dummy.Next`.
- **Deletion happens from the predecessor.** In a singly-linked list, "remove X" always means "make X's predecessor point past X" — position your pointer accordingly.
- **Don't advance after a splice.** After `prev.Next = prev.Next.Next`, the *new* `prev.Next` is unexamined (consecutive matches like `[7,7,7,7]` break if you advance unconditionally) — the classic bug in this problem.
- **Recursive list processing** = "clean the suffix, then decide self". Beautiful and uniform (no head case), but costs O(n) stack — mention the trade-off in interviews.
- The loop invariant "`prev` is the last kept node" makes correctness self-evident; state it out loud when whiteboarding.

---

## Related Problems

- LeetCode #83 — Remove Duplicates from Sorted List (same splice loop, different predicate)
- LeetCode #82 — Remove Duplicates from Sorted List II (dummy node is mandatory)
- LeetCode #27 — Remove Element (identical task on an array)
- LeetCode #19 — Remove Nth Node From End of List (dummy + two pointers)
- LeetCode #237 — Delete Node in a Linked List (deletion without the predecessor — the inverted trick)
