# 0141 — Linked List Cycle

> LeetCode #141 · Difficulty: Easy
> **Categories:** Hash Table, Linked List, Two Pointers

---

## Problem Statement

Given `head`, the head of a linked list, determine if the linked list has a cycle in it.

There is a cycle in a linked list if there is some node in the list that can be reached again by continuously following the `next` pointer. Internally, `pos` is used to denote the index of the node that tail's `next` pointer is connected to. **Note that `pos` is not passed as a parameter.**

Return `true` *if there is a cycle in the linked list*. Otherwise, return `false`.

**Example 1:**
```
Input: head = [3,2,0,-4], pos = 1
Output: true
Explanation: There is a cycle in the linked list, where the tail connects to the 1st node (0-indexed).
```

**Example 2:**
```
Input: head = [1,2], pos = 0
Output: true
Explanation: There is a cycle in the linked list, where the tail connects to the 0th node.
```

**Example 3:**
```
Input: head = [1], pos = -1
Output: false
Explanation: There is no cycle in the linked list.
```

**Constraints:**
- The number of the nodes in the list is in the range `[0, 10⁴]`.
- `-10⁵ <= Node.val <= 10⁵`
- `pos` is `-1` or a **valid index** in the linked-list.

**Follow-up:** Can you solve it using `O(1)` (i.e. constant) memory?

---

## Company Frequency

| Company    | Frequency        | Last Reported |
|------------|------------------|---------------|
| Amazon     | ★★★★★ Very High  | 2024          |
| Microsoft  | ★★★★☆ High       | 2024          |
| Bloomberg  | ★★★★☆ High       | 2024          |
| Google     | ★★★☆☆ Medium     | 2024          |
| Apple      | ★★★☆☆ Medium     | 2023          |
| Oracle     | ★★☆☆☆ Low        | 2023          |

> ⚠️ Frequency data is crowd-sourced from LeetCode Discuss, Glassdoor, and
> community interview reports. Treat as a signal, not a guarantee.

---

## DSA Concepts Used

- **Linked List** — cycle detection is the canonical linked-list structural question → see [`/dsa/linked_list.md`](/dsa/linked_list.md)
- **Two Pointers (Fast & Slow)** — Floyd's tortoise-and-hare collides iff a cycle exists → see [`/dsa/two_pointers.md`](/dsa/two_pointers.md)
- **Hash Map / Hash Set** — remembering visited node *pointers* detects the first revisit → see [`/dsa/hash_map.md`](/dsa/hash_map.md)

---

## Approaches Overview

| # | Approach | Time | Space | When to use |
|---|----------|------|-------|-------------|
| 1 | Brute Force (Nested Scan) | O(n²) | O(1) | Baseline; shows what "revisit" means without extra memory |
| 2 | Hash Set | O(n) | O(n) | First idea in an interview; simplest correct linear solution |
| 3 | Pointer Rewiring | O(n) | O(1) | Party trick when mutation is allowed; destroys the list |
| 4 | Floyd's Tortoise & Hare (Optimal) | O(n) | O(1) | The follow-up answer; expected in interviews |

---

## Approach 1 — Brute Force (Nested Scan)

### Intuition
A cycle means the traversal revisits a node it has already passed. Without any extra memory, we can detect a revisit by rescanning: the node reached at step `k` is a revisit iff it is pointer-equal to one of the first `k-1` nodes. Values cannot be used for identity (duplicates are allowed) — only pointer comparison is sound.

### Algorithm
1. Walk `cur` from `head`, keeping counter `i` = number of nodes visited before `cur`.
2. For each `cur`, walk `probe` from `head` across the first `i` nodes.
3. If `probe == cur` (same pointer), `cur` was already visited → return `true`.
4. If `cur` reaches `nil`, the list has a real tail → return `false`.

### Complexity
- **Time:** O(n²) — the first revisit must occur within `n+1` steps, and each step rescans up to `n` prefix nodes.
- **Space:** O(1) — only two pointers and two counters.

### Code
```go
func bruteForce(head *ListNode) bool {
	i := 0 // number of nodes already visited before cur
	for cur := head; cur != nil; cur = cur.Next {
		// Rescan the first i nodes looking for cur by pointer identity.
		probe := head
		for j := 0; j < i; j++ {
			if probe == cur {
				return true // cur was seen before → we looped back → cycle
			}
			probe = probe.Next
		}
		i++ // one more distinct node confirmed visited
	}
	// Reached nil → the list has a tail → no cycle possible.
	return false
}
```

### Dry Run (Example 1: head = [3,2,0,-4], pos = 1)

List: `3 → 2 → 0 → -4 → (back to 2)`. Nodes named by value: N3, N2, N0, N-4.

| Step | `cur` | `i` | Prefix scanned (`probe` path) | Pointer match? | Action |
|------|-------|-----|-------------------------------|----------------|--------|
| 1 | N3 | 0 | (none) | — | i=1, advance |
| 2 | N2 | 1 | N3 | no | i=2, advance |
| 3 | N0 | 2 | N3, N2 | no | i=3, advance |
| 4 | N-4 | 3 | N3, N2, N0 | no | i=4, advance |
| 5 | **N2** | 4 | N3, **N2** | **yes** (probe == cur) | **return `true`** ✓ |

---

## Approach 2 — Hash Set

### Intuition
Trade memory for time: store every visited `*ListNode` in a set. The moment we step onto a node already in the set, the `next` pointers have looped back. Keying on the pointer (not the value) makes duplicates harmless.

### Algorithm
1. Create an empty `map[*ListNode]bool` used as a set.
2. Walk the list from `head`.
3. For each node: if it is already in the set → return `true`.
4. Otherwise insert it and move to `Next`.
5. Reaching `nil` → return `false`.

### Complexity
- **Time:** O(n) — each node is visited at most once; hash insert/lookup are O(1) average.
- **Space:** O(n) — the set may hold every node of an acyclic list.

### Code
```go
func hashSet(head *ListNode) bool {
	seen := make(map[*ListNode]bool) // set keyed by node pointer, not value
	for cur := head; cur != nil; cur = cur.Next {
		if seen[cur] {
			return true // this exact node was visited before → cycle
		}
		seen[cur] = true // mark node as visited
	}
	return false // fell off the end → no cycle
}
```

### Dry Run (Example 1: head = [3,2,0,-4], pos = 1)

| Step | `cur` | `cur` in `seen`? | `seen` after step |
|------|-------|------------------|-------------------|
| 1 | N3 | no | {N3} |
| 2 | N2 | no | {N3, N2} |
| 3 | N0 | no | {N3, N2, N0} |
| 4 | N-4 | no | {N3, N2, N0, N-4} |
| 5 | N2 | **yes** | **return `true`** ✓ |

---

## Approach 3 — Pointer Rewiring (Destructive Marking)

### Intuition
If mutation is allowed, the list itself can store the "visited" flag: after leaving a node, point its `Next` at itself. A fresh node never self-loops (we create self-loops only after visiting), so stepping onto a node whose `Next` is itself proves we arrived through a cycle's back-edge.

### Algorithm
1. Walk `cur` from `head`.
2. If `cur.Next == cur` → this node was visited before → return `true`.
3. Otherwise save `next := cur.Next`, set `cur.Next = cur` (mark visited), then `cur = next`.
4. Reaching `nil` → return `false`.

> Destroys the list — acceptable for a pure detection puzzle, but note that LeetCode #142 explicitly says *"Do not modify the linked list."*

### Complexity
- **Time:** O(n) — every node is entered at most twice (once fresh, once via the back-edge).
- **Space:** O(1) — the visited flag lives inside the list's own pointers.

### Code
```go
func markVisited(head *ListNode) bool {
	cur := head
	for cur != nil {
		if cur.Next == cur {
			return true // self-loop marker found → we re-entered a visited node
		}
		next := cur.Next // save the onward pointer before overwriting it
		cur.Next = cur   // mark this node as visited via a self-loop
		cur = next       // continue the walk on the saved pointer
	}
	return false // hit nil → list ends → no cycle
}
```

### Dry Run (Example 1: head = [3,2,0,-4], pos = 1)

| Step | `cur` | `cur.Next` before | Self-loop? | Action |
|------|-------|-------------------|------------|--------|
| 1 | N3 | N2 | no | N3.Next = N3; cur = N2 |
| 2 | N2 | N0 | no | N2.Next = N2; cur = N0 |
| 3 | N0 | N-4 | no | N0.Next = N0; cur = N-4 |
| 4 | N-4 | N2 | no | N-4.Next = N-4; cur = N2 |
| 5 | N2 | N2 (marked in step 2) | **yes** | **return `true`** ✓ |

---

## Approach 4 — Floyd's Tortoise & Hare (Optimal)

### Intuition
Run two pointers at different speeds: `slow` moves 1 step per turn, `fast` moves 2. In an acyclic list, `fast` hits `nil`. In a cyclic list both pointers eventually enter the cycle, and inside the cycle the gap between them changes by exactly 1 each turn — so `fast` cannot jump over `slow`; it must land exactly on it. Collision ⟺ cycle.

### Algorithm
1. Initialise `slow = head`, `fast = head`.
2. While `fast != nil && fast.Next != nil`:
   1. `slow = slow.Next` (1 step).
   2. `fast = fast.Next.Next` (2 steps).
   3. If `slow == fast` → return `true`.
3. Loop exits only when `fast` runs off the end → return `false`.

### Complexity
- **Time:** O(n) — `slow` takes at most `n + C` steps (C = cycle length) before the collision; each turn is O(1).
- **Space:** O(1) — two pointers. This answers the follow-up.

### Code
```go
func floydTwoPointers(head *ListNode) bool {
	slow, fast := head, head
	// fast needs two valid hops per turn, so guard fast and fast.Next.
	for fast != nil && fast.Next != nil {
		slow = slow.Next      // tortoise: 1 step
		fast = fast.Next.Next // hare: 2 steps
		if slow == fast {
			return true // collision can only happen inside a cycle
		}
	}
	return false // fast fell off the end → acyclic
}
```

### Dry Run (Example 1: head = [3,2,0,-4], pos = 1)

List: `N3 → N2 → N0 → N-4 → N2 → …` (N-4.Next = N2 is the back-edge). Both pointers start at N3.

| Turn | `slow` move | `slow` | `fast` move (2 hops) | `fast` | `slow == fast`? |
|------|-------------|--------|----------------------|--------|-----------------|
| 1 | N3 → N2 | N2 | N3 → N2 → N0 | N0 | no |
| 2 | N2 → N0 | N0 | N0 → N-4 → N2 | N2 | no |
| 3 | N0 → N-4 | N-4 | N2 → N0 → N-4 | N-4 | **yes → return `true`** ✓ |

---

## Key Takeaways

- **Pointer identity, never value identity** — duplicate values are legal, so all cycle logic compares `*ListNode` pointers.
- **Fast & slow pointers close a gap of 1 per turn** inside a cycle, so a 2×-speed hare provably lands *exactly* on the tortoise — this is why Floyd's algorithm never misses.
- The hash-set solution is the "say it first" interview answer; Floyd's is the follow-up answer for O(1) space. Know both cold.
- Loop guard for the hare is always `fast != nil && fast.Next != nil` — memorise this exact condition; it prevents nil-pointer panics on even/odd length lists alike.
- The self-loop marking trick generalises: a mutable structure can encode "visited" in itself, giving O(1)-space marking when the interviewer permits mutation.

---

## Related Problems

- LeetCode #142 — Linked List Cycle II (same pattern: Floyd's, plus finding the cycle entrance)
- LeetCode #202 — Happy Number (cycle detection on an implicit sequence)
- LeetCode #287 — Find the Duplicate Number (Floyd's on an array treated as a linked list)
- LeetCode #876 — Middle of the Linked List (fast & slow pointers)
- LeetCode #234 — Palindrome Linked List (fast & slow to find the midpoint)
