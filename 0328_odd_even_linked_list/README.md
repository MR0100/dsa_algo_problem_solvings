# 0328 — Odd Even Linked List

> LeetCode #328 · Difficulty: Medium
> **Categories:** Linked List

---

## Problem Statement

Given the `head` of a singly linked list, group all the nodes with odd indices
together followed by the nodes with even indices, and return the reordered list.

The **first** node is considered **odd**, and the **second** node is **even**,
and so on.

Note that the relative order inside both the even and odd groups should remain
as it was in the input.

You must solve the problem in `O(1)` extra space complexity and `O(n)` time
complexity.

> ⚠️ "Odd" and "even" refer to the **position** (1-based index) of a node in the
> list, **not** the node's value.

**Example 1:**

```
Input:  head = [1,2,3,4,5]
Output: [1,3,5,2,4]
```

Explanation: The odd-position nodes are 1 (pos 1), 3 (pos 3), 5 (pos 5); the
even-position nodes are 2 (pos 2), 4 (pos 4). Grouping odd positions first, then
even positions, gives `1 -> 3 -> 5 -> 2 -> 4`.

**Example 2:**

```
Input:  head = [2,1,3,5,6,4,7]
Output: [2,3,6,7,1,5,4]
```

Explanation: Odd-position nodes are 2 (pos 1), 3 (pos 3), 6 (pos 5), 7 (pos 7);
even-position nodes are 1 (pos 2), 5 (pos 4), 4 (pos 6). Concatenating gives
`2 -> 3 -> 6 -> 7 -> 1 -> 5 -> 4`.

**Constraints:**

- The number of nodes in the linked list is in the range `[0, 10^4]`.
- `-10^6 <= Node.val <= 10^6`.

---

## Company Frequency

| Company    | Frequency        | Last Reported |
|------------|------------------|---------------|
| Amazon     | ★★★★☆ High       | 2024          |
| Microsoft  | ★★★★☆ High       | 2024          |
| Meta       | ★★★☆☆ Medium     | 2023          |
| Google     | ★★★☆☆ Medium     | 2023          |
| Bloomberg  | ★★★☆☆ Medium     | 2024          |
| Apple      | ★★☆☆☆ Low        | 2023          |

> ⚠️ Frequency data is crowd-sourced from LeetCode Discuss, Glassdoor, and
> community interview reports. Treat as a signal, not a guarantee.

---

## DSA Concepts Used

- **Singly Linked List traversal & pointer rewiring** — the whole task is to
  relink existing nodes into a new order using only next-pointer surgery (no
  value copying, no auxiliary array), which is the defining skill of linked-list
  problems → see [`/dsa/linked_list.md`](/dsa/linked_list.md).
- **Two pointers (parallel walkers)** — the optimal solution advances an `odd`
  pointer and an `even` pointer through the list in lockstep, each skipping over
  the other's nodes to unzip one list into two interleaved chains → see
  [`/dsa/two_pointers.md`](/dsa/two_pointers.md).

---

## Approaches Overview

| # | Approach | Time | Space | When to use |
|---|----------|------|-------|-------------|
| 1 | Two Sublists via Builder Tails | O(n) | O(1) | Clarity-first: explicit odd/even sublists make the grouping obvious. |
| 2 | In-Place Two-Pointer Weave (Optimal) | O(n) | O(1) | The canonical answer — fewest pointers, cleanest interview solution. |

Both meet the required `O(n)` time and `O(1)` extra space; neither allocates a
new node — they only rewire existing ones.

---

## Approach 1 — Two Sublists via Builder Tails

### Intuition

Because relative order must be preserved inside each group, a **stable split**
is all we need. Sweep the list once and, for each node, ask "is my 1-based
position odd or even?" — then append the node to the matching sublist. At the
end the answer is simply the odd sublist followed by the even sublist. Two dummy
heads plus two tail pointers let us append in O(1) each, and because we move the
*existing* nodes (not copies) the space stays O(1).

### Algorithm

1. If `head` is `nil` or `head.Next` is `nil`, return `head` unchanged (0 or 1
   node — nothing to reorder).
2. Create two sentinel dummies `oddDummy`, `evenDummy` with tails `oddTail`,
   `evenTail` initialized to those dummies.
3. Walk the original list with a 1-based position counter `pos`:
   - If `pos` is odd, append the current node to `oddTail` and advance `oddTail`.
   - If `pos` is even, append the current node to `evenTail` and advance
     `evenTail`.
   (Save `cur.Next` before rewiring, since appending overwrites it later.)
4. Terminate the even sublist with `evenTail.Next = nil` (its last node may still
   dangle into the original list).
5. Join the groups: `oddTail.Next = evenDummy.Next`.
6. Return `oddDummy.Next` — the head of the odd group.

### Complexity

- **Time:** O(n) — a single pass with O(1) work per node.
- **Space:** O(1) — four pointers (two dummies, two tails); no node allocation.

### Code

```go
func twoListsExtraNodes(head *ListNode) *ListNode {
	// Edge cases: an empty list or a single node is already grouped correctly.
	if head == nil || head.Next == nil {
		return head
	}

	oddDummy := &ListNode{}  // sentinel before the odd-position sublist
	evenDummy := &ListNode{} // sentinel before the even-position sublist
	oddTail := oddDummy      // last node currently in the odd sublist
	evenTail := evenDummy    // last node currently in the even sublist

	pos := 1    // 1-based position of `cur` in the original list
	cur := head // node we are currently classifying
	for cur != nil {
		next := cur.Next // remember the successor before we rewire cur.Next
		if pos%2 == 1 {
			oddTail.Next = cur // odd position → extend the odd sublist
			oddTail = cur      // advance the odd tail onto this node
		} else {
			evenTail.Next = cur // even position → extend the even sublist
			evenTail = cur      // advance the even tail onto this node
		}
		cur = next // move to the successor we saved
		pos++      // its position is one greater
	}

	evenTail.Next = nil           // cap the even sublist so it can't loop back
	oddTail.Next = evenDummy.Next // stitch odd group in front of even group
	return oddDummy.Next          // head of the odd group is the new head
}
```

### Dry Run

Tracing Example 1, `head = [1,2,3,4,5]`. Nodes are written by value; `oddTail`
and `evenTail` show the value each tail currently sits on (starting on their
dummies, shown as `Od`/`Ev`).

| Step | pos | cur | odd/even? | Action | oddTail | evenTail |
|------|-----|-----|-----------|--------|---------|----------|
| init | 1   | 1   | —         | dummies created | Od | Ev |
| 1    | 1   | 1   | odd       | append 1 to odd sublist | 1  | Ev |
| 2    | 2   | 2   | even      | append 2 to even sublist | 1 | 2 |
| 3    | 3   | 3   | odd       | append 3 to odd sublist | 3 | 2 |
| 4    | 4   | 4   | even      | append 4 to even sublist | 3 | 4 |
| 5    | 5   | 5   | odd       | append 5 to odd sublist | 5 | 4 |
| loop end | 6 | nil | —      | stop walking | 5 | 4 |

Finalize: `evenTail.Next = nil` caps even group `2 -> 4`; odd group is
`1 -> 3 -> 5`. Then `oddTail.Next = evenDummy.Next` links `5 -> 2`, giving
`1 -> 3 -> 5 -> 2 -> 4`. Return `oddDummy.Next = 1`. ✅ Output `[1 3 5 2 4]`.

---

## Approach 2 — In-Place Two-Pointer Weave (Optimal)

### Intuition

Odd-position and even-position nodes **already alternate** in the input: odd,
even, odd, even, .... So from any odd node, the *next* odd node is exactly two
hops away (`odd.Next.Next`), and the same holds for evens. We keep an `odd`
pointer and an `even` pointer and repeatedly splice each one past its neighbour,
"unzipping" the single list into two interleaved chains — all by pointer
surgery, never moving a value. A saved `evenHead` lets us reattach the even
chain after the odd chain once the evens run out.

### Algorithm

1. If `head` is `nil` or `head.Next` is `nil`, return `head` unchanged (0 or 1
   node).
2. Set `odd = head`, `even = head.Next`, and `evenHead = even` (remember where
   the even group begins).
3. While `even != nil` **and** `even.Next != nil`:
   - `odd.Next = even.Next; odd = odd.Next` — jump `odd` to the next odd node.
   - `even.Next = odd.Next; even = even.Next` — jump `even` to the next even node.
4. `odd.Next = evenHead` — append the whole even group after the odd group.
5. Return `head` — still the first odd node, now the head of the reordered list.

### Complexity

- **Time:** O(n) — each node is visited a constant number of times in one sweep.
- **Space:** O(1) — three pointers (`odd`, `even`, `evenHead`); pure pointer
  surgery, no allocation.

### Code

```go
func inPlacePointers(head *ListNode) *ListNode {
	// Edge cases: 0 or 1 node — already grouped, nothing to weave.
	if head == nil || head.Next == nil {
		return head
	}

	odd := head       // walks the odd-position (1,3,5,...) chain
	even := head.Next // walks the even-position (2,4,6,...) chain
	evenHead := even  // remember the even group's head for the final join

	// Continue while there is a further even node to relink. Testing both
	// `even` and `even.Next` guards the two length parities safely.
	for even != nil && even.Next != nil {
		odd.Next = even.Next // odd skips the even node to reach the next odd
		odd = odd.Next       // advance odd onto that next-odd node
		even.Next = odd.Next // even skips the (new) odd node to the next even
		even = even.Next     // advance even onto that next-even node
	}

	odd.Next = evenHead // splice the even group onto the end of the odd group
	return head         // the original first node is still the overall head
}
```

### Dry Run

Tracing Example 1, `head = [1,2,3,4,5]`. `evenHead` is fixed on node 2 the whole
time. Each row shows the pointer positions after that iteration's rewiring.

| Iter | Guard `even && even.Next` | Rewiring | odd | even | evenHead | Chain so far |
|------|---------------------------|----------|-----|------|----------|--------------|
| init | —                         | odd=1, even=2, evenHead=2 | 1 | 2 | 2 | 1→2→3→4→5 |
| 1    | even=2, even.Next=3 → true | odd.Next=3, odd=3; even.Next=4, even=4 | 3 | 4 | 2 | odd chain 1→3; even chain 2→4 |
| 2    | even=4, even.Next=5 → true | odd.Next=5, odd=5; even.Next=nil, even=nil | 5 | nil | 2 | odd chain 1→3→5; even chain 2→4 |
| end  | even=nil → false, stop    | — | 5 | nil | 2 | — |

Final join: `odd.Next = evenHead` links node 5 → node 2, producing
`1 -> 3 -> 5 -> 2 -> 4`. Return `head = 1`. ✅ Output `[1 3 5 2 4]`.

---

## Key Takeaways

- **"Odd/even" = position, not value.** The first trap is misreading the prompt;
  always confirm it's the 1-based index that matters.
- **Unzip with a two-hop stride.** Because odd and even nodes alternate, the next
  same-parity node is always `.Next.Next` away — a reusable pattern for splitting
  or de-interleaving a list in place.
- **Save `evenHead` early.** You destroy the original links as you weave, so
  capture the even group's head before the loop to reattach it at the end.
- **Loop guard `even != nil && even.Next != nil`.** This single condition
  correctly handles both odd-length and even-length lists — `odd` is guaranteed
  non-nil whenever `even.Next` is being read.
- **A sentinel-tail split is the intuitive alternative.** Building two explicit
  sublists (Approach 1) is easier to reason about and still O(1) space, since you
  relink existing nodes rather than copy them.

---

## Related Problems

- LeetCode #206 — Reverse Linked List (fundamental pointer-rewiring in place)
- LeetCode #86 — Partition List (stable split of a list into two groups by a
  predicate, then rejoin — same tail-builder pattern as Approach 1)
- LeetCode #24 — Swap Nodes in Pairs (local pointer surgery on adjacent nodes)
- LeetCode #725 — Split Linked List in Parts (partitioning a list by position)
