# 0414 — Third Maximum Number

> LeetCode #414 · Difficulty: Easy
> **Categories:** Array, Sorting, Math

---

## Problem Statement

Given an integer array `nums`, return *the **third distinct maximum** number in this array. If the third maximum does not exist, return the **maximum** number*.

**Example 1:**

```
Input: nums = [3,2,1]
Output: 1
Explanation:
The first distinct maximum is 3.
The second distinct maximum is 2.
The third distinct maximum is 1.
```

**Example 2:**

```
Input: nums = [1,2]
Output: 2
Explanation:
The first distinct maximum is 2.
The second distinct maximum does not exist, so the maximum (2) is returned instead.
```

**Example 3:**

```
Input: nums = [2,2,3,1]
Output: 1
Explanation:
The first distinct maximum is 3.
The second distinct maximum is 2 (both 2's are counted together since they have the same value).
The third distinct maximum is 1.
```

**Constraints:**

- `1 <= nums.length <= 10^4`
- `-2^31 <= nums[i] <= 2^31 - 1`

**Follow up:** Can you find an `O(n)` solution?

---

## Company Frequency

| Company    | Frequency        | Last Reported |
|------------|------------------|---------------|
| Amazon     | ★★★☆☆ Medium     | 2024          |
| Google     | ★★☆☆☆ Low        | 2023          |
| Microsoft  | ★★☆☆☆ Low        | 2023          |
| Bloomberg  | ★★☆☆☆ Low        | 2022          |
| Apple      | ★☆☆☆☆ Rare       | 2022          |

> ⚠️ Frequency data is crowd-sourced from LeetCode Discuss, Glassdoor, and
> community interview reports. Treat as a signal, not a guarantee.

---

## DSA Concepts Used

- **Array scanning / selection** — the answer is a rank-3 order statistic over the *distinct* values, computable in one linear pass without full sorting → see [`/dsa/arrays.md`](/dsa/arrays.md)
- **Sorting** — the baseline sorts the distinct values descending and indexes position 3 → see [`/dsa/sorting.md`](/dsa/sorting.md)
- **Hash set (deduplication)** — "distinct" maxima require collapsing equal values before counting → see [`/dsa/hash_map.md`](/dsa/hash_map.md)

---

## Approaches Overview

| # | Approach | Time | Space | When to use |
|---|----------|------|-------|-------------|
| 1 | Brute Force (Sort + Deduplicate) | O(n log n) | O(n) | Shortest to write; fine when n is small |
| 2 | Three-Variable Scan (Optimal) | O(n) | O(1) | The follow-up answer; single pass, constant space |
| 3 | Bounded Min-Set (Ordered Top-3) | O(n) | O(1) | Generalises to top-k; mirrors size-limited heap/TreeSet |

---

## Approach 1 — Brute Force (Sort + Deduplicate)

### Intuition

"Third distinct maximum" is by definition the element at position 3 of the distinct values listed in descending order. So deduplicate, sort big→small, and read index 2. If fewer than three distinct values exist, the rule says return the overall maximum, which is index 0.

### Algorithm

1. Insert every `nums[i]` into a set to drop duplicates.
2. Copy the set into a slice and sort it descending.
3. If the slice has `>= 3` elements, return element at index 2; otherwise return element at index 0.

### Complexity

- **Time:** O(n log n) — sorting up to `n` distinct values dominates.
- **Space:** O(n) — the set plus the deduplicated slice.

### Code

```go
func bruteForce(nums []int) int {
	seen := make(map[int]struct{}, len(nums)) // set of distinct values
	for _, v := range nums {
		seen[v] = struct{}{} // insert; duplicates collapse automatically
	}
	distinct := make([]int, 0, len(seen))
	for v := range seen {
		distinct = append(distinct, v) // materialise the distinct values
	}
	sort.Sort(sort.Reverse(sort.IntSlice(distinct))) // descending order
	if len(distinct) >= 3 {
		return distinct[2] // the 3rd distinct maximum
	}
	return distinct[0] // fewer than 3 distinct → return the maximum
}
```

### Dry Run

Example 1: `nums = [3,2,1]`.

| Step | Action | State |
|------|--------|-------|
| 1 | build set | `{3,2,1}` |
| 2 | to slice | `[3,2,1]` (any order) |
| 3 | sort descending | `[3,2,1]` |
| 4 | `len == 3 >= 3` → index 2 | `1` |

Result: `1` ✔

---

## Approach 2 — Three-Variable Scan (Optimal)

### Intuition

We only need the top three distinct values, so track a three-slot podium `first > second > third`. For each number: if it equals a slot already on the podium, ignore it (the maxima must be **distinct**); otherwise insert it into the correct place, cascading smaller slots down. The subtle trap is sentinels: because `nums[i]` can be as small as `-2^31`, a `MinInt` sentinel for "empty" is unsafe — so we use `*int` pointers where `nil` unambiguously means "slot not yet filled".

### Algorithm

1. Start with `first = second = third = nil`.
2. For each `v`: if `v` equals any *filled* slot, skip it.
3. Otherwise:
   - `v > first` (or `first` empty) → shift down and set `first = v`.
   - else `v > second` (or `second` empty) → shift `third = second`, set `second = v`.
   - else `v > third` (or `third` empty) → set `third = v`.
4. If `third` is filled, return `*third`; else return `*first`.

### Complexity

- **Time:** O(n) — one pass, constant comparisons per element.
- **Space:** O(1) — three pointer slots, independent of `n`.

### Code

```go
func threeVariableScan(nums []int) int {
	var first, second, third *int // podium slots; nil means "not yet filled"
	for i := range nums {
		v := nums[i]
		// Skip duplicates of any already-placed podium value: the three
		// maxima must be DISTINCT.
		if (first != nil && v == *first) ||
			(second != nil && v == *second) ||
			(third != nil && v == *third) {
			continue
		}
		switch {
		case first == nil || v > *first: // new overall maximum
			third = second // everyone slides down one place
			second = first
			nv := v // take an address of a fresh copy (loop var reuse safety)
			first = &nv
		case second == nil || v > *second: // fits between 1st and 2nd
			third = second
			nv := v
			second = &nv
		case third == nil || v > *third: // fits into 3rd place
			nv := v
			third = &nv
		}
	}
	if third != nil { // a genuine third distinct maximum exists
		return *third
	}
	return *first // fewer than 3 distinct values → the maximum
}
```

### Dry Run

Example 1: `nums = [3,2,1]` (`•` = nil).

| Step | v | duplicate? | branch | first | second | third |
|------|---|-----------|--------|-------|--------|-------|
| 0 | — | — | init | • | • | • |
| 1 | 3 | no | `first` empty | 3 | • | • |
| 2 | 2 | no | `second` empty | 3 | 2 | • |
| 3 | 1 | no | `third` empty | 3 | 2 | 1 |

`third != nil` → return `*third = 1` ✔

---

## Approach 3 — Bounded Min-Set (Ordered Top-3)

### Intuition

Keep a set of distinct candidates capped at size 3. Insert each distinct value; whenever the set overflows to 4, evict its **minimum**. After the scan the set holds the three largest distinct values (or all of them if fewer than three exist). This is the "top-k with a size-limited min-heap / TreeSet" pattern; with `k = 3` a tiny set is enough and each operation is O(1).

### Algorithm

1. Maintain a set `top` of distinct values.
2. For each `v`: skip if already present; else insert; if `|top| > 3`, delete the minimum.
3. If `|top| >= 3`, return the **minimum** of `top` (the smallest of the top three = the 3rd maximum); else return the **maximum** of `top`.

### Complexity

- **Time:** O(n) — the set never exceeds 3 elements, so insert/evict/scan are all O(1).
- **Space:** O(1) — at most three tracked values.

### Code

```go
func boundedSet(nums []int) int {
	top := make(map[int]struct{}, 4) // distinct candidates, capped at 3
	for _, v := range nums {
		if _, ok := top[v]; ok {
			continue // already a candidate → keep distinctness
		}
		top[v] = struct{}{}
		if len(top) > 3 { // one too many → evict the smallest
			minV := math.MaxInt64
			for k := range top {
				if k < minV {
					minV = k
				}
			}
			delete(top, minV)
		}
	}
	// Read out either the min (3rd max) or the max (fewer than 3 distinct).
	if len(top) >= 3 {
		minV := math.MaxInt64
		for k := range top {
			if k < minV {
				minV = k // smallest of the top 3 = the 3rd maximum
			}
		}
		return minV
	}
	maxV := math.MinInt64
	for k := range top {
		if k > maxV {
			maxV = k // fewer than 3 distinct → overall maximum
		}
	}
	return maxV
}
```

### Dry Run

Example 1: `nums = [3,2,1]`.

| Step | v | present? | top after insert | size > 3? evict | top |
|------|---|----------|-------------------|-----------------|-----|
| 1 | 3 | no | `{3}` | no | `{3}` |
| 2 | 2 | no | `{3,2}` | no | `{3,2}` |
| 3 | 1 | no | `{3,2,1}` | no | `{3,2,1}` |

Final size `3 >= 3` → return min `= 1` ✔

---

## Key Takeaways

- **Never use `MinInt` as an "empty" sentinel when `MinInt` is a legal input.** `nums[i]` can equal `-2^31`, so use `nil` pointers (or a separate "filled" boolean) to mark unset slots — a very common source of wrong answers on this problem.
- **Top-k without full sorting:** for a fixed small `k`, a constant number of tracked slots (or a size-`k` set/heap) turns an O(n log n) sort into an O(n) scan — the intended follow-up.
- **"Distinct" means dedupe first.** Equal values collapse into one rank, so a set (or an equality check against current slots) is mandatory.
- The min of the "top-3 set" is exactly the 3rd maximum — a clean reformulation that generalises to "kth largest".

---

## Related Problems

- LeetCode #215 — Kth Largest Element in an Array (general top-k)
- LeetCode #347 — Top K Frequent Elements (bounded heap / bucket)
- LeetCode #628 — Maximum Product of Three Numbers (track top-3 / bottom-2)
- LeetCode #703 — Kth Largest Element in a Stream (size-k min-heap)
