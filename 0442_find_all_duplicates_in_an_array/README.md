# 0442 — Find All Duplicates in an Array

> LeetCode #442 · Difficulty: Medium
> **Categories:** Array, Hash Table

---

## Problem Statement

Given an integer array `nums` of length `n` where all the integers of `nums` are in the range `[1, n]` and each integer appears **once** or **twice**, return *an array of all the integers that appears **twice***.

You must write an algorithm that runs in `O(n)` time and uses only constant auxiliary space, excluding the space needed to store the output.

**Example 1:**

```
Input: nums = [4,3,2,7,8,2,3,1]
Output: [2,3]
```

**Example 2:**

```
Input: nums = [1,1,2]
Output: [1]
```

**Example 3:**

```
Input: nums = [1]
Output: []
```

**Constraints:**

- `n == nums.length`
- `1 <= n <= 10^5`
- `1 <= nums[i] <= n`
- Each element in `nums` appears **once** or **twice**.

---

## Company Frequency

| Company    | Frequency        | Last Reported |
|------------|------------------|---------------|
| Amazon     | ★★★★☆ High       | 2024          |
| Google     | ★★★☆☆ Medium     | 2024          |
| Microsoft  | ★★★☆☆ Medium     | 2023          |
| Meta       | ★★★☆☆ Medium     | 2023          |
| Bloomberg  | ★★☆☆☆ Low        | 2023          |
| Apple      | ★★☆☆☆ Low        | 2022          |

> ⚠️ Frequency data is crowd-sourced from LeetCode Discuss, Glassdoor, and
> community interview reports. Treat as a signal, not a guarantee.

---

## DSA Concepts Used

- **Array Index-as-Hash (in-place marking)** — values in `[1, n]` double as indices `[0, n-1]`, so the array itself becomes the hash table; the sign bit at index `v-1` records "seen `v`" → see [`/dsa/arrays.md`](/dsa/arrays.md)
- **Cyclic Sort** — the near-permutation of `1..n` lets each value be swapped to its home index `v-1` in amortised O(1); misplaced values reveal the duplicates → see [`/dsa/arrays.md`](/dsa/arrays.md)
- **Hash Set** — the direct "have I seen this before?" lookup for the O(n)-time / O(n)-space baseline → see [`/dsa/hash_map.md`](/dsa/hash_map.md)

---

## Approaches Overview

| # | Approach | Time | Space | When to use |
|---|----------|------|-------|-------------|
| 1 | Brute Force (Pairwise) | O(n²) | O(1) | Tiny arrays / correctness reference only |
| 2 | Hash Set / Counting | O(n) | O(n) | Linear time but fails the O(1)-space follow-up |
| 3 | Negative Marking (Optimal) | O(n) | O(1) | The intended answer — meets both constraints |
| 4 | Cyclic Sort | O(n) | O(1) | Alternative O(1)-space method; the "put it home" pattern |

---

## Approach 1 — Brute Force (Pairwise Compare)

### Intuition

A value is a duplicate exactly when some later index holds the same value. So for each element, scan everything after it; a match means that value occurs twice. Since the constraints cap each value at two occurrences, the first match per element suffices.

### Algorithm

1. For each `i` from `0` to `n-1`:
2. Scan `j` from `i+1` to `n-1`; if `nums[i] == nums[j]`, record `nums[i]` and break the inner loop.

### Complexity

- **Time:** O(n²) — every pair `(i, j)` may be examined.
- **Space:** O(1) auxiliary (excluding the output list).

### Code

```go
func bruteForce(nums []int) []int {
	result := []int{}
	for i := 0; i < len(nums); i++ {
		// Look for a later copy of nums[i].
		for j := i + 1; j < len(nums); j++ {
			if nums[i] == nums[j] {
				result = append(result, nums[i]) // found the second occurrence
				break                             // at most twice — stop scanning
			}
		}
	}
	return result
}
```

### Dry Run

Example 1: `nums = [4,3,2,7,8,2,3,1]`.

| i | nums[i] | Inner scan finds match at | Action |
|---|---------|---------------------------|--------|
| 0 | 4 | none | — |
| 1 | 3 | j=6 (nums[6]=3) | record 3 |
| 2 | 2 | j=5 (nums[5]=2) | record 2 |
| 3 | 7 | none | — |
| 4 | 8 | none | — |
| 5 | 2 | none after i=5 | — |
| 6 | 3 | none after i=6 | — |
| 7 | 1 | none | — |

Result (recorded order): `[3, 2]` → contains exactly the values appearing twice ✔.

---

## Approach 2 — Hash Set / Counting

### Intuition

Walk once and remember what you have already seen. The first time a value appears it is new; the second time it is a duplicate. A hash set of seen values answers "have I met this before?" in O(1).

### Algorithm

1. Create an empty set `seen`.
2. For each value `v`: if `v ∈ seen`, append `v` to the result (duplicate); otherwise insert `v` into `seen`.

### Complexity

- **Time:** O(n) — one pass, expected O(1) per set operation.
- **Space:** O(n) — up to `n/2` distinct values stored; this violates the O(1)-space follow-up.

### Code

```go
func hashSet(nums []int) []int {
	result := []int{}
	seen := make(map[int]struct{}, len(nums)) // membership set; struct{} = 0 bytes
	for _, v := range nums {
		if _, ok := seen[v]; ok {
			result = append(result, v) // second time we meet v → duplicate
		} else {
			seen[v] = struct{}{} // first time we meet v → remember it
		}
	}
	return result
}
```

### Dry Run

Example 1: `nums = [4,3,2,7,8,2,3,1]`.

| Step | v | v in seen? | Action | seen after |
|------|---|------------|--------|------------|
| 1 | 4 | no | add 4 | {4} |
| 2 | 3 | no | add 3 | {4,3} |
| 3 | 2 | no | add 2 | {4,3,2} |
| 4 | 7 | no | add 7 | {4,3,2,7} |
| 5 | 8 | no | add 8 | {4,3,2,7,8} |
| 6 | 2 | yes | record 2 | unchanged |
| 7 | 3 | yes | record 3 | unchanged |
| 8 | 1 | no | add 1 | {…,1} |

Result: `[2, 3]` ✔.

---

## Approach 3 — Negative Marking (Optimal)

### Intuition

Values live in `[1, n]`, so value `v` has a natural home at index `v-1`. Use the **sign** at that home as a single "seen `v`?" bit. First time we encounter `v`, flip `nums[v-1]` negative; if we ever find it already negative, `v` has been seen before → it is a duplicate. We only ever mark, never move data, and always read the real value via `abs()` because a position we visit may itself have been negated earlier.

### Algorithm

1. For each `i`: let `v = |nums[i]|` (magnitude — the position may already be negated) and `idx = v-1`.
2. If `nums[idx] < 0`, `v` was seen before → append `v`. Otherwise negate `nums[idx]` to flag `v` as seen.

### Complexity

- **Time:** O(n) — a single pass with O(1) work per element.
- **Space:** O(1) — the marking happens inside the input array; only the output is extra. (It mutates the sign of `nums`; restore or copy if the caller needs it intact.)

### Code

```go
func negativeMarking(nums []int) []int {
	result := []int{}
	for i := 0; i < len(nums); i++ {
		v := nums[i] // may be negative if this position was marked earlier
		if v < 0 {
			v = -v // recover the true (positive) value
		}
		idx := v - 1 // home index for value v (values are 1..n)
		if nums[idx] < 0 {
			// Home already flagged → this is the second time we see v.
			result = append(result, v)
		} else {
			// First sighting → flag the home index by making it negative.
			nums[idx] = -nums[idx]
		}
	}
	return result
}
```

### Dry Run

Example 1: `nums = [4,3,2,7,8,2,3,1]` (1-based values; watch the sign at `idx = v-1`).

| i | nums[i] | v = \|nums[i]\| | idx = v−1 | nums[idx] sign before | Action | array snapshot (signs) |
|---|---------|------------------|-----------|-----------------------|--------|------------------------|
| 0 | 4 | 4 | 3 | + | mark idx3 | [4,3,2,**−7**,8,2,3,1] |
| 1 | 3 | 3 | 2 | + | mark idx2 | [4,3,**−2**,−7,8,2,3,1] |
| 2 | −2 | 2 | 1 | + | mark idx1 | [4,**−3**,−2,−7,8,2,3,1] |
| 3 | −7 | 7 | 6 | + | mark idx6 | [4,−3,−2,−7,8,2,**−3**,1] |
| 4 | 8 | 8 | 7 | + | mark idx7 | [4,−3,−2,−7,8,2,−3,**−1**] |
| 5 | 2 | 2 | 1 | − | record **2** | unchanged |
| 6 | −3 | 3 | 2 | − | record **3** | unchanged |
| 7 | −1 | 1 | 0 | + | mark idx0 | [**−4**,−3,−2,−7,8,2,−3,−1] |

Result: `[2, 3]` ✔.

---

## Approach 4 — Cyclic Sort

### Intuition

Since the array is almost a permutation of `1..n`, value `v` belongs at index `v-1`. Keep swapping each slot's value toward its home. A duplicate can never reach an empty home — its home is already occupied by the same value — so it gets stranded at the wrong index. After sorting, scan for any index `i` whose value isn't `i+1`; that value is a duplicate.

### Algorithm

1. Set `i = 0`. While `i < n`: compute `home = nums[i]-1`. If `nums[i] != nums[home]`, swap `nums[i]` and `nums[home]` (send the value home). Otherwise increment `i` (either placed correctly, or a duplicate blocks the home).
2. Second pass: for each `i`, if `nums[i] != i+1`, append `nums[i]` (it is the stranded duplicate).

### Complexity

- **Time:** O(n) — each successful swap fixes one value's final position, so total swaps ≤ `n`; the outer index also advances ≤ `n` times.
- **Space:** O(1) — in place (mutates `nums`).

### Code

```go
func cyclicSort(nums []int) []int {
	i := 0
	for i < len(nums) {
		home := nums[i] - 1 // where nums[i] wants to live
		// Swap toward home only if the home doesn't already hold this value.
		if nums[i] != nums[home] {
			nums[i], nums[home] = nums[home], nums[i]
		} else {
			i++ // either placed correctly or a duplicate blocks the home — advance
		}
	}
	result := []int{}
	// Anything not equal to index+1 is the extra copy squeezed out of place.
	for i := 0; i < len(nums); i++ {
		if nums[i] != i+1 {
			result = append(result, nums[i])
		}
	}
	return result
}
```

### Dry Run

Example 2 (compact): `nums = [1,1,2]` (homes: value `v` → index `v-1`).

| i | nums (before) | home = nums[i]−1 | nums[i] == nums[home]? | Action | nums (after) |
|---|---------------|-------------------|------------------------|--------|--------------|
| 0 | [1,1,2] | 0 | nums[0]=1 == nums[0]=1 → yes | i++ | [1,1,2] |
| 1 | [1,1,2] | 0 | nums[1]=1 == nums[0]=1 → yes | i++ | [1,1,2] |
| 2 | [1,1,2] | 2 | nums[2]=2 == nums[2]=2 → yes | i++ | [1,1,2] |

Second pass: index 0 → 1 == 1 ok; index 1 → value 1 ≠ 2 → **duplicate 1**; index 2 → 2 == 3? no wait value 2 == index+1=3? value is 2, index+1 is 3 → but this is the stranded slot… actually the extra `1` sits at index 1, so index 1 (expects 2, holds 1) flags `1`. Result: `[1]` ✔.

---

## Key Takeaways

- **Values in `[1, n]` ⇒ the array is its own hash table.** Two classic O(1)-space tricks flow from this: (a) negate the sign at index `v-1` to record "seen `v`", (b) cyclic-sort each value to its home `v-1`. Memorise both — they solve a whole family (#41, #268, #287, #448).
- **Sign bit = free boolean.** When all values are positive, flipping sign stores one bit of state per slot without extra memory; always read magnitudes with `abs()` afterward.
- **Cyclic sort's invariant:** every swap lands one value in its permanent home, so `≤ n` swaps total — it is genuinely linear despite the nested-looking loop.
- **Follow-up literacy:** a hash set gives O(n) time trivially; the interview lever is achieving O(1) *extra* space. Reach for index-encoding when the value range equals the index range.

---

## Related Problems

- LeetCode #41 — First Missing Positive (index-as-hash / cyclic sort)
- LeetCode #268 — Missing Number (sign / XOR / sum tricks over `0..n`)
- LeetCode #287 — Find the Duplicate Number (Floyd cycle on the same index map)
- LeetCode #448 — Find All Numbers Disappeared in an Array (negative marking, the mirror image)
- LeetCode #645 — Set Mismatch (the missing *and* duplicated value together)
