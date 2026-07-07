# 0457 — Circular Array Loop

> LeetCode #457 · Difficulty: Medium
> **Categories:** Array, Hash Table, Two Pointers

---

## Problem Statement

You are playing a game involving a **circular** array of non-zero integers `nums`. Each `nums[i]` denotes the number of indices forward/backward you must move if you are located at index `i`:

- If `nums[i]` is positive, move `nums[i]` steps **forward**, and
- If `nums[i]` is negative, move `nums[i]` steps **backward**.

Since the array is **circular**, you may assume that moving forward from the last element puts you on the first element, and moving backwards from the first element puts you on the last element.

A **cycle** in the array consists of a sequence of indices `seq` of length `k` where:

- Following the movement rules above results in the repeating index sequence `seq[0] -> seq[1] -> ... -> seq[k - 1] -> seq[0] -> ...`
- Every `nums[seq[j]]` is **either all positive or all negative**.
- `k > 1`

Return `true` *if there is a **cycle** in* `nums`*, or* `false` *otherwise*.

**Example 1:**

```
Input: nums = [2,-1,1,2,2]
Output: true
Explanation: There is a cycle from index 0 -> 2 -> 3 -> 0 -> ...
The cycle's length is 3.
```

**Example 2:**

```
Input: nums = [-1,-2,-3,-4,-5,6]
Output: false
Explanation: The only cycle is of length 1, so it is not a cycle.
```

**Example 3:**

```
Input: nums = [1,-1,5,1,4]
Output: true
Explanation: There is a cycle from index 0 -> 1 -> 0 -> ...
The cycle's length is 2.
Note: The sequence of indices [3, 4] is not a cycle because [4, 1] is not all positive or all negative.
```

**Constraints:**

- `1 <= nums.length <= 5000`
- `-1000 <= nums[i] <= 1000`
- `nums[i] != 0`

**Follow up:** Could you solve it in `O(n)` time complexity and `O(1)` extra space complexity?

---

## Company Frequency

| Company    | Frequency        | Last Reported |
|------------|------------------|---------------|
| Google     | ★★★☆☆ Medium     | 2024          |
| Amazon     | ★★★☆☆ Medium     | 2023          |
| Microsoft  | ★★☆☆☆ Low        | 2023          |
| Bloomberg  | ★★☆☆☆ Low        | 2022          |
| Uber       | ★★☆☆☆ Low        | 2022          |

> ⚠️ Frequency data is crowd-sourced from LeetCode Discuss, Glassdoor, and
> community interview reports. Treat as a signal, not a guarantee.

---

## DSA Concepts Used

- **Two Pointers (Floyd's Tortoise & Hare)** — the optimal O(1)-space solution runs a slow/fast pointer per start index to detect a cycle without a visited set → see [`/dsa/two_pointers.md`](/dsa/two_pointers.md)
- **Arrays** — indices are treated as nodes and `nums[i]` as circular edge weights; the modular `next` helper and in-place "dead node = 0" marking are pure array manipulation → see [`/dsa/arrays.md`](/dsa/arrays.md)
- **Hash Map** — the brute-force baseline records the current walk's visited indices in a hash set to spot a repeat → see [`/dsa/hash_map.md`](/dsa/hash_map.md)

---

## Approaches Overview

| # | Approach | Time | Space | When to use |
|---|----------|------|-------|-------------|
| 1 | Brute Force (per-start visited set) | O(n²) | O(n) | Clear correctness model; fine for n ≤ 5000 |
| 2 | Floyd Tortoise & Hare + 0-marking (Optimal) | O(n) | O(1) | Meets the follow-up: linear time, constant extra space |

---

## Approach 1 — Brute Force (Per-Start Visited Set)

### Intuition

Treat each index as a node with a single outgoing edge (`next(i)`). A valid cycle must (a) return to a start index, (b) keep one direction throughout, and (c) have length > 1. Start a fresh walk from every index, storing the indices visited **on that walk** in a set. Revisiting a set member means we closed a loop. Two guards enforce the extra rules: bail the instant a node's sign disagrees with the start's direction, and reject a one-hop self-loop (`next(cur) == cur`).

### Algorithm

1. For each start index `s`:
   1. `dir = nums[s] > 0` (required direction).
   2. Walk with an empty `seen` set, `cur = s`.
   3. If `sign(nums[cur]) != dir`, break (direction changed).
   4. Compute `nxt = next(cur)`; if `nxt == cur`, break (self-loop).
   5. If `cur` is already in `seen`, return `true` (cycle of length > 1).
   6. Add `cur` to `seen`, set `cur = nxt`.
2. If no start yields a cycle, return `false`.

### Complexity

- **Time:** O(n²) — up to `n` starts, each walking up to `n` nodes.
- **Space:** O(n) — the per-walk visited set.

### Code

```go
func bruteForce(nums []int) bool {
	n := len(nums)
	for s := 0; s < n; s++ {
		dir := nums[s] > 0    // required direction: true = forward, false = backward
		seen := map[int]bool{} // indices visited on this particular walk
		cur := s
		for {
			// Direction check: every element on the cycle must share the sign
			// of the start; a mixed-sign path is disqualified.
			if (nums[cur] > 0) != dir {
				break // direction flipped → this start cannot yield a valid cycle
			}
			nxt := next(nums, cur) // one circular jump from cur
			if nxt == cur {
				break // self-loop (length 1) is explicitly not allowed
			}
			if seen[cur] {
				return true // we have returned to a node already on this walk → cycle len > 1
			}
			seen[cur] = true // mark cur as visited on this walk
			cur = nxt        // advance
		}
	}
	return false // no start index produced a valid cycle
}
```

### Dry Run

Example 1: `nums = [2, -1, 1, 2, 2]`, `n = 5`. Start `s = 0`, `dir = true` (forward). Recall `next(i) = ((i+nums[i])%n + n) % n`.

| cur | nums[cur] | sign == dir? | next(cur) | self-loop? | cur in seen? | action |
|-----|-----------|--------------|-----------|------------|--------------|--------|
| 0 | 2 | yes | (0+2)%5=2 | no | no | add 0, cur=2 |
| 2 | 1 | yes | (2+1)%5=3 | no | no | add 2, cur=3 |
| 3 | 2 | yes | (3+2)%5=0 | no | no | add 3, cur=0 |
| 0 | 2 | yes | 2 | no | **yes (0∈seen)** | return true |

Cycle `0 → 2 → 3 → 0`, length 3, all forward → `true` ✔

---

## Approach 2 — Floyd's Tortoise & Hare, In-Place Marking (Optimal)

### Intuition

Cycle detection with O(1) space is Floyd's algorithm: a slow pointer (one hop) and a fast pointer (two hops); if they collide, a cycle exists. Two problem rules bolt on top. **Direction:** define a "valid step" that only proceeds if the destination keeps the start's sign and is not a self-loop; any invalid step kills the current attempt. **Length > 1:** the self-loop check handles it. The key to *linear total time* is cleanup: once a start attempt fails, every node it walked can never belong to a valid cycle, so overwrite those nodes with `0` (a dead marker). Future starts skip `0` nodes, so each element is charged O(1) amortised work.

### Algorithm

1. For each start `i` with `nums[i] != 0`:
   1. `dir = nums[i] > 0`; `slow = fast = i`.
   2. Loop: advance `slow` by one valid step and `fast` by two valid steps. A valid step requires same direction and no self-loop; an invalid step returns `-1` and breaks the loop.
   3. If `slow == fast`, return `true` (they met inside a valid cycle).
2. On failure, re-walk from `i` overwriting each node with `0` until a direction change or self-loop, so it is never re-explored.
3. If no start succeeds, return `false`.

### Complexity

- **Time:** O(n) — the `0`-marking guarantees every index participates in at most a constant number of hops across all starts.
- **Space:** O(1) — only pointers; the visited information is folded back into the input array.

### Code

```go
func floydFastSlow(nums []int) bool {
	n := len(nums)

	// sameDirection reports whether index i still moves in direction `dir`
	// (dir true = forward). A node whose sign disagrees ends the current path.
	sameDirection := func(i int, dir bool) bool {
		return (nums[i] > 0) == dir
	}

	for i := 0; i < n; i++ {
		if nums[i] == 0 {
			continue // already marked dead by a previous failed attempt
		}
		dir := nums[i] > 0 // required direction for this attempt
		slow, fast := i, i // both pointers start at i

		for {
			// Advance slow by one valid hop.
			slow = validNext(nums, slow, dir, sameDirection)
			if slow == -1 {
				break // slow hit a direction change or self-loop → path invalid
			}
			// Advance fast by two valid hops.
			fast = validNext(nums, fast, dir, sameDirection)
			if fast == -1 {
				break
			}
			fast = validNext(nums, fast, dir, sameDirection)
			if fast == -1 {
				break
			}
			if slow == fast {
				return true // pointers met inside a same-direction cycle of length > 1
			}
		}

		// This start failed: overwrite the whole path with 0 so it is never
		// re-explored, keeping the amortised cost linear.
		j := i
		for nums[j] != 0 && sameDirection(j, dir) {
			nxt := next(nums, j) // remember where this node pointed
			nums[j] = 0          // mark node j as dead
			if nxt == j {
				break // self-loop node; stop marking
			}
			j = nxt
		}
	}
	return false // no valid cycle from any start
}
```

Supporting helpers:

```go
func next(nums []int, i int) int {
	n := len(nums)
	return ((i+nums[i])%n + n) % n // circular index, safe for negative modulo
}

func validNext(nums []int, i int, dir bool, sameDirection func(int, bool) bool) int {
	if !sameDirection(i, dir) {
		return -1 // sign of nums[i] disagrees with the cycle's direction
	}
	nxt := next(nums, i)
	if nxt == i {
		return -1 // length-1 self-loop is not a valid cycle
	}
	return nxt
}
```

### Dry Run

Example 1: `nums = [2, -1, 1, 2, 2]`, start `i = 0`, `dir = true`. `slow, fast = 0, 0`. All jumps forward and never self-loop, so every `validNext` succeeds.

| iteration | slow (1 hop) | fast (2 hops) | slow == fast? |
|-----------|--------------|---------------|----------------|
| start | 0 | 0 | — |
| 1 | 0→2 | 0→2→3 | 2 vs 3, no |
| 2 | 2→3 | 3→0→2 | 3 vs 2, no |
| 3 | 3→0 | 2→3→0 | **0 vs 0, yes** |

`slow == fast == 0` inside the all-forward cycle `0→2→3→0` → return `true` ✔

---

## Key Takeaways

- **Functional graph = one edge per node.** Whenever every element points to exactly one successor (here `next(i)`), cycle-detection tools (Floyd, visited sets, union-find) apply directly.
- **Negative modulo needs `((x % n) + n) % n`.** Go's `%` keeps the sign of the dividend, so backward jumps require the double-mod normalisation.
- **Fold the "visited" bit into the data.** Overwriting dead nodes with a sentinel (`0`) turns a per-start O(n²) walk into O(n) total while staying O(1) space — a classic in-place amortisation trick.
- **Extra validity rules layer onto a known algorithm.** The direction and length-> 1 constraints are just guards inside each step; the skeleton is still plain tortoise-and-hare.

---

## Related Problems

- LeetCode #141 — Linked List Cycle (Floyd's on a linked list)
- LeetCode #142 — Linked List Cycle II (find the cycle entrance)
- LeetCode #202 — Happy Number (Floyd's on a number sequence)
- LeetCode #287 — Find the Duplicate Number (Floyd's on an implicit functional graph)
- LeetCode #565 — Array Nesting (longest cycle in a permutation)
