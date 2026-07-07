# 0142 — Linked List Cycle II

> LeetCode #142 · Difficulty: Medium
> **Categories:** Hash Table, Linked List, Two Pointers

---

## Problem Statement

Given the `head` of a linked list, return *the node where the cycle begins. If there is no cycle, return* `null`.

There is a cycle in a linked list if there is some node in the list that can be reached again by continuously following the `next` pointer. Internally, `pos` is used to denote the index of the node that tail's `next` pointer is connected to (**0-indexed**). It is `-1` if there is no cycle. **Note that `pos` is not passed as a parameter.**

**Do not modify** the linked list.

**Example 1:**
```
Input: head = [3,2,0,-4], pos = 1
Output: tail connects to node index 1
Explanation: There is a cycle in the linked list, where tail connects to the second node.
```

**Example 2:**
```
Input: head = [1,2], pos = 0
Output: tail connects to node index 0
Explanation: There is a cycle in the linked list, where tail connects to the first node.
```

**Example 3:**
```
Input: head = [1], pos = -1
Output: no cycle
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
| Amazon     | ★★★★☆ High       | 2024          |
| Microsoft  | ★★★★☆ High       | 2024          |
| Google     | ★★★☆☆ Medium     | 2024          |
| Bloomberg  | ★★★☆☆ Medium     | 2024          |
| Apple      | ★★☆☆☆ Low        | 2023          |
| Adobe      | ★★☆☆☆ Low        | 2023          |

> ⚠️ Frequency data is crowd-sourced from LeetCode Discuss, Glassdoor, and
> community interview reports. Treat as a signal, not a guarantee.

---

## DSA Concepts Used

- **Linked List** — locating a structural feature (the cycle entrance) by pointer manipulation → see [`/dsa/linked_list.md`](/dsa/linked_list.md)
- **Two Pointers (Fast & Slow)** — Floyd's collision plus the `F = m·C − k` phase-2 walk pinpoints the entrance → see [`/dsa/two_pointers.md`](/dsa/two_pointers.md)
- **Hash Map / Hash Set** — the first node already in the visited set is the entrance → see [`/dsa/hash_map.md`](/dsa/hash_map.md)

---

## Approaches Overview

| # | Approach | Time | Space | When to use |
|---|----------|------|-------|-------------|
| 1 | Brute Force (Nested Scan) | O(n²) | O(1) | Baseline; proves "first revisited node = entrance" |
| 2 | Hash Set | O(n) | O(n) | Fastest to write correctly; fine when memory is free |
| 3 | Sentinel Rewiring | O(n) | O(1) | Teaching trick — violates the "do not modify" rule |
| 4 | Floyd's Cycle Detection (Optimal) | O(n) | O(1) | The follow-up answer; O(1) space without mutation |

---

## Approach 1 — Brute Force (Nested Scan)

### Intuition
Walking from `head`, every node strictly before the cycle entrance is passed exactly once — the entrance is the first node the walk ever returns to (the tail's back-edge lands there). So "find the entrance" reduces to "find the first revisited node", which we can detect with pointer-equality rescans of the prefix, no extra memory needed.

### Algorithm
1. Walk `cur` from `head` with counter `i` = number of nodes visited before `cur`.
2. For each `cur`, walk `probe` from `head` across the first `i` nodes.
3. If `probe == cur` → `cur` is the first revisited node → return `cur`.
4. If `cur` hits `nil` → the list terminates → return `nil`.

### Complexity
- **Time:** O(n²) — the first revisit happens within `n+1` steps, each step rescans up to `n` nodes.
- **Space:** O(1) — two pointers, two counters.

### Code
```go
func bruteForce(head *ListNode) *ListNode {
	i := 0 // number of nodes already visited before cur
	for cur := head; cur != nil; cur = cur.Next {
		// Rescan the first i nodes; a pointer match means cur is revisited.
		probe := head
		for j := 0; j < i; j++ {
			if probe == cur {
				return cur // first node seen twice = cycle entrance
			}
			probe = probe.Next
		}
		i++ // cur confirmed as a newly visited node
	}
	return nil // reached the tail → acyclic
}
```

### Dry Run (Example 1: head = [3,2,0,-4], pos = 1)

List: `N3 → N2 → N0 → N-4 → (back to N2)`.

| Step | `cur` | `i` | Prefix scanned | Match? | Action |
|------|-------|-----|----------------|--------|--------|
| 1 | N3 | 0 | (none) | — | i=1, advance |
| 2 | N2 | 1 | N3 | no | i=2, advance |
| 3 | N0 | 2 | N3, N2 | no | i=3, advance |
| 4 | N-4 | 3 | N3, N2, N0 | no | i=4, advance |
| 5 | **N2** | 4 | N3, **N2** | **yes** | **return N2 — index 1** ✓ |

---

## Approach 2 — Hash Set

### Intuition
Same "first revisited node = entrance" observation, but membership testing becomes O(1) by storing every visited `*ListNode` in a hash set. The first node found already in the set is exactly where the tail's `next` reconnects.

### Algorithm
1. `seen` = empty `map[*ListNode]bool`.
2. Walk `cur` from `head`.
3. If `cur` ∈ `seen` → return `cur` (the entrance).
4. Otherwise insert `cur`, advance to `cur.Next`.
5. `nil` reached → return `nil`.

### Complexity
- **Time:** O(n) — single pass, O(1) average hash operations.
- **Space:** O(n) — the set can hold all nodes.

### Code
```go
func hashSet(head *ListNode) *ListNode {
	seen := make(map[*ListNode]bool) // visited set keyed by node pointer
	for cur := head; cur != nil; cur = cur.Next {
		if seen[cur] {
			return cur // first repeat = where the tail's back-edge lands
		}
		seen[cur] = true // mark visited
	}
	return nil // no repeat before nil → no cycle
}
```

### Dry Run (Example 1: head = [3,2,0,-4], pos = 1)

| Step | `cur` | In `seen`? | `seen` after |
|------|-------|-----------|--------------|
| 1 | N3 | no | {N3} |
| 2 | N2 | no | {N3, N2} |
| 3 | N0 | no | {N3, N2, N0} |
| 4 | N-4 | no | {N3, N2, N0, N-4} |
| 5 | N2 | **yes** | **return N2 — index 1** ✓ |

---

## Approach 3 — Sentinel Rewiring (Destructive Marking)

### Intuition
If mutation were allowed, the list could store its own visited flags: after leaving a node, aim its `Next` at one shared sentinel node. Stepping onto a node whose `Next` is the sentinel means we arrived via the back-edge — and the *first* such node is the entrance. This achieves O(1) space, but **violates this problem's "Do not modify the linked list" rule**, so treat it as a thought experiment that motivates Floyd's.

### Algorithm
1. Create a sentinel node `S` (no real node ever points to it initially).
2. Walk `cur` from `head`.
3. If `cur.Next == S` → `cur` was visited before → return `cur`.
4. Otherwise save `next := cur.Next`, set `cur.Next = S`, advance to `next`.
5. `nil` reached → return `nil`.

### Complexity
- **Time:** O(n) — each node is entered at most twice.
- **Space:** O(1) — a single sentinel node.

### Code
```go
func sentinelRewiring(head *ListNode) *ListNode {
	sentinel := &ListNode{} // unique marker no real node ever points to
	cur := head
	for cur != nil {
		if cur.Next == sentinel {
			return cur // already marked → first revisited node → entrance
		}
		next := cur.Next    // save onward pointer before overwriting
		cur.Next = sentinel // mark cur as visited
		cur = next          // continue on the saved pointer
	}
	return nil // fell off the tail → no cycle
}
```

### Dry Run (Example 1: head = [3,2,0,-4], pos = 1)

| Step | `cur` | `cur.Next` before | Points to S? | Action |
|------|-------|-------------------|--------------|--------|
| 1 | N3 | N2 | no | N3.Next = S; cur = N2 |
| 2 | N2 | N0 | no | N2.Next = S; cur = N0 |
| 3 | N0 | N-4 | no | N0.Next = S; cur = N-4 |
| 4 | N-4 | N2 | no | N-4.Next = S; cur = N2 |
| 5 | N2 | S (set in step 2) | **yes** | **return N2 — index 1** ✓ |

---

## Approach 4 — Floyd's Cycle Detection (Optimal)

### Intuition
Phase 1 is classic tortoise-and-hare: collide inside the cycle. Phase 2 is the beautiful part. Let `F` = distance head → entrance, `C` = cycle length, and suppose the pointers collide `k` steps past the entrance. Then:

- slow travelled `F + k`
- fast travelled `F + k + m·C` for some laps `m ≥ 1`
- fast is twice as fast: `2(F + k) = F + k + m·C` ⟹ `F + k = m·C` ⟹ **`F = m·C − k`**

So from the collision point, walking `F` steps forward lands on the entrance (because `k + F ≡ 0 (mod C)`). We don't know `F` numerically — but a pointer starting at `head` reaches the entrance in exactly `F` steps too. March both at speed 1; their first meeting is the entrance.

### Algorithm
1. **Phase 1 — detect:** `slow = fast = head`. While `fast != nil && fast.Next != nil`: `slow` moves 1, `fast` moves 2. If `slow == fast` → collision found; else if the loop ends → return `nil`.
2. **Phase 2 — locate:** `ptr1 = head`, `ptr2 = collision node`.
3. While `ptr1 != ptr2`: advance both by exactly 1 step.
4. Return `ptr1` — the cycle entrance.

### Complexity
- **Time:** O(n) — phase 1 collides within `F + C ≤ 2n` slow-steps; phase 2 takes `F ≤ n` steps.
- **Space:** O(1) — four pointers, no mutation. Satisfies both the follow-up and the "do not modify" rule.

### Code
```go
func floydCycle(head *ListNode) *ListNode {
	slow, fast := head, head
	// Phase 1: advance until collision or end of list.
	for fast != nil && fast.Next != nil {
		slow = slow.Next      // tortoise: 1 step
		fast = fast.Next.Next // hare: 2 steps
		if slow == fast {
			// Phase 2: F = m·C − k ⇒ head and collision point are the same
			// distance (mod C) from the entrance.
			ptr1, ptr2 := head, slow
			for ptr1 != ptr2 {
				ptr1 = ptr1.Next // both move at speed 1
				ptr2 = ptr2.Next
			}
			return ptr1 // meeting point = cycle entrance
		}
	}
	return nil // fast fell off the end → acyclic
}
```

### Dry Run (Example 1: head = [3,2,0,-4], pos = 1)

List: `N3 → N2 → N0 → N-4 → N2 → …` — here `F = 1` (head → N2), `C = 3`.

**Phase 1** (both start at N3):

| Turn | `slow` | `fast` (2 hops) | Collision? |
|------|--------|-----------------|------------|
| 1 | N2 | N0 | no |
| 2 | N0 | N2 | no |
| 3 | N-4 | N-4 | **yes — collide at N-4** |

Check the math: collision is `k = 2` steps past entrance N2; `F = m·C − k = 1·3 − 2 = 1` ✓.

**Phase 2** (`ptr1 = head = N3`, `ptr2 = N-4`):

| Step | `ptr1` | `ptr2` | Equal? |
|------|--------|--------|--------|
| start | N3 | N-4 | no |
| 1 | N2 | N2 (N-4.Next) | **yes → return N2 — index 1** ✓ |

---

## Key Takeaways

- **First revisited node = cycle entrance** — this single observation powers the brute force, hash set, and sentinel approaches identically.
- **The Floyd phase-2 identity `F = m·C − k`** is the derivation interviewers ask for; be able to reproduce the two-line algebra (`2(F+k) = F+k+m·C`).
- Phase 2 is symmetric: one pointer from `head`, one from the collision, both at speed 1 — first meeting is the entrance, no counting required.
- Read the constraints: this problem *forbids* mutation, which rules out marking tricks and makes Floyd's the only O(1)-space answer.
- Pattern transfer: #287 (Find the Duplicate Number) is exactly this algorithm on the implicit list `i → nums[i]`.

---

## Related Problems

- LeetCode #141 — Linked List Cycle (phase 1 alone)
- LeetCode #287 — Find the Duplicate Number (Floyd's on an implicit linked list)
- LeetCode #202 — Happy Number (cycle detection on a digit-square sequence)
- LeetCode #160 — Intersection of Two Linked Lists (pointer meeting-point trick)
- LeetCode #457 — Circular Array Loop (fast/slow pointers in an array)
