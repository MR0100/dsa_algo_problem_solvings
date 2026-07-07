# 0169 — Majority Element

> LeetCode #169 · Difficulty: Easy
> **Categories:** Array, Hash Table, Divide and Conquer, Sorting, Counting

---

## Problem Statement

Given an array `nums` of size `n`, return *the majority element*.

The majority element is the element that appears more than `⌊n / 2⌋` times. You may assume that the majority element **always exists** in the array.

**Example 1:**

```
Input: nums = [3,2,3]
Output: 3
```

**Example 2:**

```
Input: nums = [2,2,1,1,1,2,2]
Output: 2
```

**Constraints:**

- `n == nums.length`
- `1 <= n <= 5 * 10^4`
- `-10^9 <= nums[i] <= 10^9`

**Follow-up:** Could you solve the problem in linear time and in `O(1)` space?

---

## Company Frequency

| Company    | Frequency        | Last Reported |
|------------|------------------|---------------|
| Amazon     | ★★★★★ Very High  | 2024          |
| Google     | ★★★★☆ High       | 2024          |
| Microsoft  | ★★★★☆ High       | 2024          |
| Apple      | ★★★☆☆ Medium     | 2023          |
| Adobe      | ★★★☆☆ Medium     | 2023          |
| Bloomberg  | ★★☆☆☆ Low        | 2023          |

> ⚠️ Frequency data is crowd-sourced from LeetCode Discuss, Glassdoor, and
> community interview reports. Treat as a signal, not a guarantee.

---

## DSA Concepts Used

- **Hash Map** — value → frequency counting; the straightforward O(n) time / O(n) space answer → see [`/dsa/hash_map.md`](/dsa/hash_map.md)
- **Sorting** — a block longer than n/2 must cover index ⌊n/2⌋ after sorting → see [`/dsa/sorting.md`](/dsa/sorting.md)
- **Divide and Conquer** — the global majority must be the majority of at least one half; merge by counting → see [`/dsa/divide_and_conquer.md`](/dsa/divide_and_conquer.md)
- **Bit Manipulation** — the majority element dictates the majority bit at every one of the 32 positions → see [`/dsa/bit_manipulation.md`](/dsa/bit_manipulation.md)
- **Greedy / Streaming (Boyer–Moore Voting)** — pair off majority vs non-majority occurrences with a single counter; the survivor is the answer → see [`/dsa/greedy.md`](/dsa/greedy.md)

---

## Approaches Overview

| # | Approach | Time | Space | When to use |
|---|----------|------|-------|-------------|
| 1 | Brute Force | O(n²) | O(1) | Baseline; n = 5·10⁴ makes ~2.5·10⁹ comparisons — too slow |
| 2 | Hash Map Counting | O(n) | O(n) | Default quick answer when memory is free |
| 3 | Sorting | O(n log n) | O(n) copy (O(1) in-place) | One-liner when mutating/sorting is acceptable |
| 4 | Divide and Conquer | O(n log n) | O(log n) | To demonstrate the D&C decomposition pattern |
| 5 | Bit Manipulation | O(32·n) = O(n) | O(1) | Alternative O(1)-space answer; generalises to #137 Single Number II |
| 6 | Boyer–Moore Voting (Optimal) | O(n) | O(1) | Always — answers the follow-up exactly |

---

## Approach 1 — Brute Force

### Intuition

Count each candidate's occurrences directly. The element whose count exceeds ⌊n/2⌋ is by definition the majority. No cleverness, only nested loops.

### Algorithm

1. Compute the threshold `majorityCount = n / 2`.
2. For each `candidate` in `nums`, scan the entire array and count how many elements equal it.
3. Return the first candidate whose count strictly exceeds the threshold.

### Complexity

- **Time:** O(n²) — each of the n candidates triggers a full O(n) counting scan.
- **Space:** O(1) — two integer counters.

### Code

```go
func bruteForce(nums []int) int {
	majorityCount := len(nums) / 2 // need strictly more than this many
	for _, candidate := range nums {
		count := 0
		for _, v := range nums {
			if v == candidate {
				count++ // tally occurrences of this candidate
			}
		}
		if count > majorityCount {
			return candidate // first value crossing the threshold wins
		}
	}
	return -1 // unreachable: a majority element always exists
}
```

### Dry Run

Example 1: `nums = [3,2,3]`, threshold `majorityCount = 3/2 = 1`.

| Step | candidate | Inner scan counts | count | count > 1? | Action |
|------|-----------|-------------------|-------|------------|--------|
| 1 | 3 | 3 ✓, 2 ✗, 3 ✓ | 2 | yes | return 3 |

Result: `3` ✔

---

## Approach 2 — Hash Map Counting

### Intuition

Build all frequencies in one pass with a map. Because the majority element's final count exceeds n/2, we can return the instant *any* running count crosses the threshold — the answer cannot be beaten later.

### Algorithm

1. `counts := map[int]int{}`, threshold `n/2`.
2. For each `v`: increment `counts[v]`; if `counts[v] > n/2`, return `v`.

### Complexity

- **Time:** O(n) — single pass, O(1) average per map update.
- **Space:** O(n) — the map can hold up to ⌈n/2⌉+ distinct values before the threshold trips.

### Code

```go
func hashMap(nums []int) int {
	counts := map[int]int{} // value → occurrences seen so far
	majorityCount := len(nums) / 2
	for _, v := range nums {
		counts[v]++
		if counts[v] > majorityCount {
			return v // threshold crossed — must be the majority
		}
	}
	return -1 // unreachable: a majority element always exists
}
```

### Dry Run

Example 1: `nums = [3,2,3]`, threshold `1`.

| Step | v | counts after | counts[v] | counts[v] > 1? | Action |
|------|---|--------------|-----------|----------------|--------|
| 1 | 3 | {3:1} | 1 | no | continue |
| 2 | 2 | {3:1, 2:1} | 1 | no | continue |
| 3 | 3 | {3:2, 2:1} | 2 | yes | return 3 |

Result: `3` ✔

---

## Approach 3 — Sorting

### Intuition

Sort the array: equal values become one contiguous block. A block of length > n/2 cannot avoid the middle index — wherever it starts, it spans position ⌊n/2⌋. So `sorted[n/2]` *is* the majority element, unconditionally.

### Algorithm

1. Copy `nums` (to leave the input untouched) and sort the copy ascending.
2. Return the element at index `n/2`.

### Complexity

- **Time:** O(n log n) — the comparison sort dominates.
- **Space:** O(n) for the defensive copy — O(1) (ignoring sort internals) if in-place mutation is acceptable.

### Code

```go
func sorting(nums []int) int {
	sorted := make([]int, len(nums))
	copy(sorted, nums) // don't mutate the caller's slice
	sort.Ints(sorted)
	// A run longer than n/2 must straddle the midpoint.
	return sorted[len(sorted)/2]
}
```

### Dry Run

Example 1: `nums = [3,2,3]`.

| Step | Action | sorted | n/2 | sorted[n/2] |
|------|--------|--------|-----|-------------|
| 1 | copy input | [3,2,3] | — | — |
| 2 | sort ascending | [2,3,3] | — | — |
| 3 | index the middle | [2,3,3] | 1 | 3 |

Result: `3` ✔ (the 3-block occupies indices 1..2 and covers the midpoint 1).

---

## Approach 4 — Divide and Conquer

### Intuition

If `x` is the majority of the whole array, it must be the majority of the left half or the right half (or both). Proof by contraposition: if `x` is ≤ half in *both* halves, its total is ≤ n/2 — contradiction. So the global answer is among the two halves' winners; when they disagree, an O(range) counting pass decides.

### Algorithm

1. `majorityInRange(lo, hi)`: if `lo == hi`, return `nums[lo]`.
2. Recurse on `[lo, mid]` and `[mid+1, hi]`.
3. If both winners agree → return it.
4. Else count each winner's occurrences in `[lo, hi]` and return the more frequent one.

### Complexity

- **Time:** O(n log n) — recurrence T(n) = 2T(n/2) + O(n) (the disagreement counting), by the Master Theorem.
- **Space:** O(log n) — recursion depth; no auxiliary arrays.

### Code

```go
func divideAndConquer(nums []int) int {
	return majorityInRange(nums, 0, len(nums)-1)
}

// majorityInRange returns the majority element of nums[lo..hi] (inclusive).
func majorityInRange(nums []int, lo, hi int) int {
	// Base case: a single element is trivially the majority of its range.
	if lo == hi {
		return nums[lo]
	}
	mid := lo + (hi-lo)/2 // overflow-safe midpoint
	left := majorityInRange(nums, lo, mid)
	right := majorityInRange(nums, mid+1, hi)
	// Both halves elected the same value → it wins the merged range too.
	if left == right {
		return left
	}
	// Halves disagree → count both candidates over the whole range.
	if countInRange(nums, left, lo, hi) > countInRange(nums, right, lo, hi) {
		return left
	}
	return right
}

// countInRange counts occurrences of target in nums[lo..hi].
func countInRange(nums []int, target, lo, hi int) int {
	count := 0
	for i := lo; i <= hi; i++ {
		if nums[i] == target {
			count++
		}
	}
	return count
}
```

### Dry Run

Example 1: `nums = [3,2,3]`.

| Step | Call | lo | hi | mid | left result | right result | Action |
|------|------|----|----|-----|-------------|--------------|--------|
| 1 | `majorityInRange(0,2)` | 0 | 2 | 1 | pending | pending | recurse both halves |
| 2 | `majorityInRange(0,1)` | 0 | 1 | 0 | pending | pending | recurse |
| 3 | `majorityInRange(0,0)` | 0 | 0 | — | — | — | base → 3 |
| 4 | `majorityInRange(1,1)` | 1 | 1 | — | — | — | base → 2 |
| 5 | back in (0,1): left=3, right=2 disagree | 0 | 1 | — | 3 | 2 | count 3→1, count 2→1; tie → return right = 2 |
| 6 | `majorityInRange(2,2)` | 2 | 2 | — | — | — | base → 3 |
| 7 | back in (0,2): left=2, right=3 disagree | 0 | 2 | — | 2 | 3 | count 2→1, count 3→2; 3 wins → return 3 |

Result: `3` ✔ (note step 5: a *local* tie is harmless — the true majority always wins the *top-level* count).

---

## Approach 5 — Bit Manipulation

### Intuition

Decide the answer one bit at a time. At any bit position, more than n/2 of the array's elements are copies of the majority element — so whatever bit value the majority element has there is automatically the *majority bit* at that position. Count ones per position; if ones > n/2 the answer has a 1 there, else a 0. Doing the extraction on `int32` makes bit 31 the sign bit, so negative majorities reconstruct correctly via two's complement.

### Algorithm

1. For each `bit` in 0..31:
   1. Count elements whose `int32` form has that bit set.
   2. If the count exceeds n/2, OR `1 << bit` into an `int32` accumulator.
2. Widen the accumulator back to `int` (sign-extends for negatives) and return.

### Complexity

- **Time:** O(32·n) = O(n) — 32 full passes.
- **Space:** O(1) — one accumulator and one counter.

### Code

```go
func bitManipulation(nums []int) int {
	n := len(nums)
	var answer int32 // assemble in int32 so bit 31 doubles as the sign bit
	for bit := 0; bit < 32; bit++ {
		ones := 0
		for _, v := range nums {
			// Extract this bit from the two's-complement int32 form.
			if (int32(v)>>bit)&1 == 1 {
				ones++
			}
		}
		// The majority element dictates the majority bit at every position.
		if ones > n/2 {
			answer |= 1 << bit
		}
	}
	return int(answer) // widen back; sign extends automatically for negatives
}
```

### Dry Run

Example 1: `nums = [3,2,3]` (binary: 3 = `11`, 2 = `10`), n/2 = 1.

| Step | bit | bits of [3,2,3] at position | ones | ones > 1? | answer (binary) |
|------|-----|------------------------------|------|-----------|-----------------|
| 1 | 0 | 1, 0, 1 | 2 | yes | `...01` |
| 2 | 1 | 1, 1, 1 | 3 | yes | `...11` |
| 3 | 2..31 | 0, 0, 0 | 0 | no | `...11` |

Result: binary `11` = `3` ✔

---

## Approach 6 — Boyer–Moore Voting (Optimal)

### Intuition

Think of it as an election with mutual annihilation: keep one `candidate` and a `count`. Matching elements vote +1; non-matching elements vote −1; at zero the next element takes over as candidate. Every −1 pairs one candidate occurrence against one different element. Since the majority element has more than n/2 copies, it can absorb *all* pair-offs (at most n/2 of them) and still have votes left — so it must be the candidate standing at the end. This is the follow-up's linear-time, constant-space answer.

### Algorithm

1. `count = 0`, `candidate` undefined.
2. For each `v` in `nums`:
   1. If `count == 0`, set `candidate = v`.
   2. If `v == candidate`, `count++`; else `count--`.
3. Return `candidate`. (If a majority were **not** guaranteed, a second verification pass counting `candidate` would be mandatory.)

### Complexity

- **Time:** O(n) — exactly one pass, O(1) work per element.
- **Space:** O(1) — one candidate and one counter, regardless of n.

### Code

```go
func boyerMooreVoting(nums []int) int {
	count := 0
	candidate := 0
	for _, v := range nums {
		// Counter exhausted → previous candidate fully paired off; adopt v.
		if count == 0 {
			candidate = v
		}
		if v == candidate {
			count++ // a vote for the current candidate
		} else {
			count-- // a vote against — cancels one supporter
		}
	}
	// Majority is guaranteed to exist, so no verification pass is needed.
	return candidate
}
```

### Dry Run

Example 1: `nums = [3,2,3]`.

| Step | v | count == 0? | candidate | v == candidate? | count after |
|------|---|-------------|-----------|-----------------|-------------|
| 1 | 3 | yes → candidate = 3 | 3 | yes | 1 |
| 2 | 2 | no | 3 | no | 0 |
| 3 | 3 | yes → candidate = 3 | 3 | yes | 1 |

Loop ends with `candidate = 3`. Result: `3` ✔

Bonus trace of Example 2 (`nums = [2,2,1,1,1,2,2]`):

| Step | v | count == 0? | candidate | count after |
|------|---|-------------|-----------|-------------|
| 1 | 2 | yes → candidate = 2 | 2 | 1 |
| 2 | 2 | no | 2 | 2 |
| 3 | 1 | no | 2 | 1 |
| 4 | 1 | no | 2 | 0 |
| 5 | 1 | yes → candidate = 1 | 1 | 1 |
| 6 | 2 | no | 1 | 0 |
| 7 | 2 | yes → candidate = 2 | 2 | 1 |

Result: `2` ✔ (the lead changes twice, but the true majority survives the final takeover).

---

## Key Takeaways

- **Boyer–Moore Voting** is the canonical "> n/2 occurrences" tool: one candidate + one counter, O(n)/O(1). Generalisation: for elements appearing > n/k times, keep k−1 candidate/counter pairs (LeetCode #229 uses k = 3).
- **Remember the verification pass.** The single-pass vote is only trusted here because the majority is *guaranteed*; without that guarantee, always re-count the surviving candidate.
- **"Majority" transfers to sub-structures:** the majority must dominate ≥ 1 half (divide & conquer) and dominates every bit position (bit manipulation) — the same counting argument reused in two different domains.
- Sorted-array shortcut: any element with > n/2 copies must occupy index ⌊n/2⌋. Useful as a quick correctness check in tests.

---

## Related Problems

- LeetCode #229 — Majority Element II (elements appearing > n/3 times; generalised Boyer–Moore with two candidates)
- LeetCode #1150 — Check If a Number Is Majority Element in a Sorted Array (binary-search variant)
- LeetCode #2404 — Most Frequent Even Element (frequency counting pattern)
- LeetCode #137 — Single Number II (the same per-bit counting trick)
