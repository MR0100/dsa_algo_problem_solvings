# 0189 — Rotate Array

> LeetCode #189 · Difficulty: Medium
> **Categories:** Array, Math, Two Pointers

---

## Problem Statement

Given an integer array `nums`, rotate the array to the right by `k` steps, where `k` is non-negative.

**Example 1:**
```
Input: nums = [1,2,3,4,5,6,7], k = 3
Output: [5,6,7,1,2,3,4]
Explanation:
rotate 1 steps to the right: [7,1,2,3,4,5,6]
rotate 2 steps to the right: [6,7,1,2,3,4,5]
rotate 3 steps to the right: [5,6,7,1,2,3,4]
```

**Example 2:**
```
Input: nums = [-1,-100,3,99], k = 2
Output: [3,99,-1,-100]
Explanation:
rotate 1 steps to the right: [99,-1,-100,3]
rotate 2 steps to the right: [3,99,-1,-100]
```

**Constraints:**
- `1 <= nums.length <= 10^5`
- `-2^31 <= nums[i] <= 2^31 - 1`
- `0 <= k <= 10^5`

**Follow-up:**
- Try to come up with as many solutions as you can. There are at least **three** different ways to solve this problem.
- Could you do it in-place with `O(1)` extra space?

---

## Company Frequency

| Company    | Frequency        | Last Reported |
|------------|------------------|---------------|
| Microsoft  | ★★★★★ Very High  | 2024          |
| Amazon     | ★★★★☆ High       | 2024          |
| Google     | ★★★★☆ High       | 2024          |
| Apple      | ★★★☆☆ Medium     | 2024          |
| Meta       | ★★★☆☆ Medium     | 2023          |
| Bloomberg  | ★★★☆☆ Medium     | 2023          |

> ⚠️ Frequency data is crowd-sourced from LeetCode Discuss, Glassdoor, and
> community interview reports. Treat as a signal, not a guarantee.

---

## DSA Concepts Used

- **Two Pointers (in-place reversal)** — every reversal is a converging swap loop, and "reverse all, then reverse each block" is the O(1)-space workhorse → see [`/dsa/two_pointers.md`](/dsa/two_pointers.md)
- **Math / Modular Arithmetic** — the destination of index `i` is `(i+k) mod n`; the cycle structure of that permutation (there are exactly `gcd(n,k)` cycles) drives the cyclic-replacement proof → see [`/dsa/math_number_theory.md`](/dsa/math_number_theory.md)

---

## Approaches Overview

| # | Approach | Time | Space | When to use |
|---|----------|------|-------|-------------|
| 1 | Brute Force (rotate 1 step, k times) | O(n·k) | O(1) | Conceptual baseline only; TLEs at 10⁵ × 10⁵ |
| 2 | Extra Array | O(n) | O(n) | Simplest correct fast answer when memory is free |
| 3 | Cyclic Replacements | O(n) | O(1) | In-place with the minimum possible writes (one per element) |
| 4 | Reversal (Optimal) | O(n) | O(1) | The interview answer: in-place, short, impossible to get wrong |

---

## Approach 1 — Brute Force (Rotate One Step, k Times)

### Intuition
A rotation by `k` is literally `k` rotations by one, and a single right rotation is trivial: remember the last element, slide everything one slot right, and drop the remembered element at the front. First reduce `k` modulo `n` — rotating by exactly `n` returns the array to itself, so only `k mod n` matters.

### Algorithm
1. `k %= n`.
2. Repeat `k` times:
   1. `last = nums[n-1]`.
   2. For `i` from `n-1` down to `1`: `nums[i] = nums[i-1]` (right-to-left so nothing is overwritten before it is moved).
   3. `nums[0] = last`.

### Complexity
- **Time:** O(n·k) — each of the k single-step rotations shifts all n elements; worst case ~10¹⁰ operations at the constraint limits, so this TLEs on LeetCode.
- **Space:** O(1) — only the single saved element.

### Code
```go
func bruteForce(nums []int, k int) {
	n := len(nums)
	if n == 0 {
		return
	}
	k %= n // rotating by n is a no-op, so only the remainder matters
	for step := 0; step < k; step++ {
		last := nums[n-1] // the element that wraps around to the front
		// shift every element one slot to the right (right to left to avoid clobbering)
		for i := n - 1; i > 0; i-- {
			nums[i] = nums[i-1]
		}
		nums[0] = last // wrapped element lands at the front
	}
}
```

### Dry Run
Example 1: `nums = [1,2,3,4,5,6,7], k = 3`

| step | saved `last` | array after shift + placement |
|------|--------------|-------------------------------|
| init | — | `[1,2,3,4,5,6,7]` |
| 1 | 7 | `[7,1,2,3,4,5,6]` |
| 2 | 6 | `[6,7,1,2,3,4,5]` |
| 3 | 5 | `[5,6,7,1,2,3,4]` ✓ |

Matches the explanation in the problem statement step for step.

---

## Approach 2 — Extra Array

### Intuition
Why move elements repeatedly when every element's final home is known in closed form? A right rotation by `k` sends index `i` to `(i + k) mod n`. Allocate a scratch array, teleport each element straight to its destination, then copy the scratch back over `nums` (the judge inspects `nums` itself, so the copy-back matters).

### Algorithm
1. Allocate `rotated` of length `n`.
2. For every index `i`: `rotated[(i+k) % n] = nums[i]`.
3. `copy(nums, rotated)`.

### Complexity
- **Time:** O(n) — one placement pass and one copy pass.
- **Space:** O(n) — the auxiliary array; this is exactly what the follow-up asks us to eliminate.

### Code
```go
func extraArray(nums []int, k int) {
	n := len(nums)
	if n == 0 {
		return
	}
	rotated := make([]int, n) // scratch array in final order
	for i, v := range nums {
		rotated[(i+k)%n] = v // closed-form destination of index i
	}
	copy(nums, rotated) // the problem wants nums itself mutated
}
```

### Dry Run
Example 1: `nums = [1,2,3,4,5,6,7], k = 3` (n = 7)

| i | nums[i] | destination `(i+3) % 7` | rotated so far |
|---|---------|--------------------------|----------------|
| 0 | 1 | 3 | `[_,_,_,1,_,_,_]` |
| 1 | 2 | 4 | `[_,_,_,1,2,_,_]` |
| 2 | 3 | 5 | `[_,_,_,1,2,3,_]` |
| 3 | 4 | 6 | `[_,_,_,1,2,3,4]` |
| 4 | 5 | 0 | `[5,_,_,1,2,3,4]` |
| 5 | 6 | 1 | `[5,6,_,1,2,3,4]` |
| 6 | 7 | 2 | `[5,6,7,1,2,3,4]` |

Copy back → `nums = [5,6,7,1,2,3,4]` ✓

---

## Approach 3 — Cyclic Replacements

### Intuition
The mapping `i → (i+k) mod n` is a permutation, and every permutation decomposes into disjoint cycles — here exactly `gcd(n, k)` of them. Walk a cycle carrying one value "in hand": drop it in its destination slot, pick up the evicted occupant, hop `k` forward, and repeat; the walk returns to its start precisely when the cycle closes. Each element is written exactly once — the theoretical minimum. When `gcd(n,k) > 1` a single cycle misses some indices (e.g. n = 6, k = 2 only touches even slots from start 0), so keep a placement counter and start a new cycle at the next index until all `n` placements are done.

### Algorithm
1. `k %= n`; if `k == 0`, return (identity).
2. `count = 0`; for `start = 0, 1, 2, …` while `count < n`:
   1. `current = start`, `carried = nums[start]`.
   2. Loop: `next = (current + k) % n`; swap `carried` with `nums[next]`; `current = next`; `count++`; stop when `current == start`.

### Complexity
- **Time:** O(n) — the counter guarantees exactly n placements in total across all cycles.
- **Space:** O(1) — one carried value and three integers.

### Code
```go
func cyclicReplacements(nums []int, k int) {
	n := len(nums)
	if n == 0 {
		return
	}
	k %= n
	if k == 0 {
		return // identity rotation: touching nothing is correct
	}
	count := 0 // how many elements have reached their final slot
	for start := 0; count < n; start++ {
		current := start
		carried := nums[start] // value in hand, waiting to be placed k ahead
		for {
			next := (current + k) % n // final slot of the carried value
			// place the carried value and pick up the displaced one in a single swap
			nums[next], carried = carried, nums[next]
			current = next
			count++
			if current == start { // cycle closed — everything on it is placed
				break
			}
		}
	}
}
```

### Dry Run
Example 1: `nums = [1,2,3,4,5,6,7], k = 3` — gcd(7,3) = 1, so one cycle covers everything (start = 0, carrying 1):

| hop | next = (cur+3)%7 | place carried | pick up | array state | count |
|-----|-------------------|---------------|---------|-------------|-------|
| 1 | 3 | 1 → slot 3 | 4 | `[1,2,3,1,5,6,7]` | 1 |
| 2 | 6 | 4 → slot 6 | 7 | `[1,2,3,1,5,6,4]` | 2 |
| 3 | 2 | 7 → slot 2 | 3 | `[1,2,7,1,5,6,4]` | 3 |
| 4 | 5 | 3 → slot 5 | 6 | `[1,2,7,1,5,3,4]` | 4 |
| 5 | 1 | 6 → slot 1 | 2 | `[1,6,7,1,5,3,4]` | 5 |
| 6 | 4 | 2 → slot 4 | 5 | `[1,6,7,1,2,3,4]` | 6 |
| 7 | 0 | 5 → slot 0 | 1 | `[5,6,7,1,2,3,4]` ✓ | 7 = n → done |

The walk 0→3→6→2→5→1→4→0 closes exactly at the start with every element placed once.

---

## Approach 4 — Reversal (Optimal)

### Intuition
A right rotation by `k` splits the array into two blocks: the last `k` elements (which must move to the front, order preserved) and the first `n−k` (which follow, order preserved). Reversing the **whole** array brings the tail block to the front — but both blocks come out internally backwards. Reversing each block separately repairs their internal order. Each element is reversed exactly twice inside its block, cancelling out, while the block-order flip from the full reversal survives. Three loops, no cycle bookkeeping, O(1) space.

### Algorithm
1. `k %= n`.
2. Reverse `nums[0..n-1]` (whole array).
3. Reverse `nums[0..k-1]` (the arrived tail block).
4. Reverse `nums[k..n-1]` (the pushed-back head block).

### Complexity
- **Time:** O(n) — three linear passes; every element is swapped at most twice in total.
- **Space:** O(1) — all reversals swap within `nums`; only loop indices are allocated.

### Code
```go
func reversal(nums []int, k int) {
	n := len(nums)
	if n == 0 {
		return
	}
	k %= n
	reverseRange(nums, 0, n-1) // whole array: tail block arrives at the front (backwards)
	reverseRange(nums, 0, k-1) // fix internal order of the first k (the old tail)
	reverseRange(nums, k, n-1) // fix internal order of the rest (the old head)
}

// reverseRange reverses nums[lo..hi] in place with two converging pointers.
func reverseRange(nums []int, lo, hi int) {
	for lo < hi {
		nums[lo], nums[hi] = nums[hi], nums[lo] // swap the outermost pair
		lo++
		hi--
	}
}
```

### Dry Run
Example 1: `nums = [1,2,3,4,5,6,7], k = 3`

| Step | Action | Array state |
|------|--------|-------------|
| 0 | initial | `[1,2,3,4,5,6,7]` |
| 1 | reverse `[0..6]` (whole) | `[7,6,5,4,3,2,1]` |
| 2 | reverse `[0..2]` (first k = 3) | `[5,6,7,4,3,2,1]` |
| 3 | reverse `[3..6]` (rest) | `[5,6,7,1,2,3,4]` ✓ |

Edge check: if `k % n == 0`, step 2 reverses an empty prefix (no-op) and steps 1 + 3 reverse the full array twice — identity, as required.

---

## Key Takeaways

- **Always `k %= n` first** — rotation is periodic with period n; this both fixes correctness for `k > n` and protects brute-force costs.
- **Triple reversal** is the go-to O(1)-space block-reordering trick: reverse the whole, then each block. The identical pattern solves #186/#151 (word reversal) and string rotations.
- **Cyclic replacement** achieves the minimum possible writes (one per element) and its termination argument — `gcd(n,k)` disjoint cycles — is a classic number-theory-meets-arrays interview probe.
- In Go, remember the judge (and `main`) observe the *same slice*: in-place approaches mutate directly, and buffer-based approaches must `copy` back into `nums`, not reassign the local slice header.

---

## Related Problems

- LeetCode #61 — Rotate List (same rotation on a linked list)
- LeetCode #186 — Reverse Words in a String II (same triple-reversal pattern)
- LeetCode #151 — Reverse Words in a String (block reversal with parsing)
- LeetCode #48 — Rotate Image (2D rotation via reversals/transpose)
- LeetCode #396 — Rotate Function (math over all rotations without performing them)
