# 0448 — Find All Numbers Disappeared in an Array

> LeetCode #448 · Difficulty: Easy
> **Categories:** Array, Hash Table

---

## Problem Statement

Given an array `nums` of `n` integers where `nums[i]` is in the range `[1, n]`, return *an array of all the integers in the range* `[1, n]` *that do not appear in* `nums`.

**Example 1:**

```
Input: nums = [4,3,2,7,8,2,3,1]
Output: [5,6]
```

**Example 2:**

```
Input: nums = [1,1]
Output: [2]
```

**Constraints:**

- `n == nums.length`
- `1 <= n <= 10^5`
- `1 <= nums[i] <= n`

**Follow up:** Could you do it without extra space and in `O(n)` runtime? You may assume the returned list does not count as extra space.

---

## Company Frequency

| Company    | Frequency        | Last Reported |
|------------|------------------|---------------|
| Google     | ★★★★☆ High       | 2024          |
| Amazon     | ★★★★☆ High       | 2024          |
| Microsoft  | ★★★☆☆ Medium     | 2023          |
| Bloomberg  | ★★★☆☆ Medium     | 2023          |
| Apple      | ★★☆☆☆ Low        | 2022          |
| Adobe      | ★★☆☆☆ Low        | 2022          |

> ⚠️ Frequency data is crowd-sourced from LeetCode Discuss, Glassdoor, and
> community interview reports. Treat as a signal, not a guarantee.

---

## DSA Concepts Used

- **Array (index-as-hash)** — because values are confined to `[1, n]` and there are exactly `n` slots, value `v` can be addressed directly at index `v−1`, letting the array double as its own presence table → see [`/dsa/arrays.md`](/dsa/arrays.md)
- **Hash Table** — the brute-force baseline records seen values in a boolean/set structure for O(1) membership before collecting the absentees → see [`/dsa/hash_map.md`](/dsa/hash_map.md)

---

## Approaches Overview

| # | Approach | Time | Space | When to use |
|---|----------|------|-------|-------------|
| 1 | Brute Force (Seen Set) | O(n) | O(n) | Clear and safe when mutating the input is disallowed |
| 2 | In-Place Negation Marking (Optimal) | O(n) | O(1) | The follow-up answer; O(1) extra space by using the sign bit |

---

## Approach 1 — Brute Force (Seen Set)

### Intuition

Every value is guaranteed to lie in `[1, n]`. Keep a boolean table `seen` of size `n+1`; mark `seen[x] = true` for each `x` in `nums`. Afterwards, sweep `v` from `1` to `n`: any `v` whose flag is still `false` never appeared, so it is one of the missing numbers. Simple and non-destructive, at the cost of an O(n) side table.

### Algorithm

1. Allocate `seen := make([]bool, n+1)` (index 0 unused).
2. For each `x` in `nums`, set `seen[x] = true`.
3. For `v` from `1` to `n`: if `!seen[v]`, append `v` to the result.
4. Return the result.

### Complexity

- **Time:** O(n) — one pass to mark, one pass to collect.
- **Space:** O(n) — the boolean table (the output list is not counted).

### Code

```go
func bruteForce(nums []int) []int {
	n := len(nums)
	seen := make([]bool, n+1) // index 0 unused; values range over 1..n
	for _, x := range nums {
		seen[x] = true // this value is present
	}
	res := []int{}
	for v := 1; v <= n; v++ {
		if !seen[v] { // v in [1,n] but never marked → it is missing
			res = append(res, v)
		}
	}
	return res
}
```

### Dry Run

Example 1: `nums = [4,3,2,7,8,2,3,1]`, `n = 8`.

| Step | Action | seen (indices 1..8) marked true |
|------|--------|----------------------------------|
| mark | process all x | {1,2,3,4,7,8} |
| collect v=1 | seen[1]=T | — |
| collect v=2 | seen[2]=T | — |
| collect v=3 | seen[3]=T | — |
| collect v=4 | seen[4]=T | — |
| collect v=5 | seen[5]=F | append **5** |
| collect v=6 | seen[6]=F | append **6** |
| collect v=7 | seen[7]=T | — |
| collect v=8 | seen[8]=T | — |

Result: `[5, 6]`. ✔

---

## Approach 2 — In-Place Negation Marking (Optimal)

### Intuition

Turn the input array into its own presence table. Value `v` maps to index `v−1`. Walk the array; for each value `v`, negate the element at `nums[v−1]` to stamp "value `v` was seen." Read `v` as `abs(nums[i])` because an earlier stamp may already have flipped `nums[i]` negative — the sign is scratch space, but the *magnitude* still holds the real value. After the pass, any slot `i` that remains **positive** was never stamped, so value `i+1` is absent. This achieves the follow-up's O(1) extra space.

### Algorithm

1. **Pass 1 (stamp):** for each `i`, let `v = abs(nums[i])`; if `nums[v−1] > 0`, negate it.
2. **Pass 2 (collect):** for each index `i`, if `nums[i] > 0`, value `i+1` is missing → append it.
3. Return the collected values. (Signs can be restored by taking absolute values if the caller needs the array intact.)

### Complexity

- **Time:** O(n) — two linear passes.
- **Space:** O(1) — mutates `nums` in place; only the output list is allocated.

### Code

```go
func inPlaceMarking(nums []int) []int {
	n := len(nums)
	// Pass 1: stamp presence by making nums[value-1] negative.
	for i := 0; i < n; i++ {
		v := nums[i] // current value (may already be negated by a prior stamp)
		if v < 0 {
			v = -v // recover the true value 1..n
		}
		idx := v - 1 // the slot that represents value v
		if nums[idx] > 0 {
			nums[idx] = -nums[idx] // stamp: mark value v as seen
		}
	}
	// Pass 2: any still-positive slot i was never stamped → value i+1 missing.
	res := []int{}
	for i := 0; i < n; i++ {
		if nums[i] > 0 {
			res = append(res, i+1) // slot i encodes value i+1
		}
	}
	return res
}
```

### Dry Run

Example 1: `nums = [4,3,2,7,8,2,3,1]`, `n = 8`. Pass 1 stamps `nums[v−1]` negative for each value `v = abs(nums[i])`:

| i | v = abs(nums[i]) | idx = v−1 | nums after stamping (bold = flipped) |
|---|------------------|-----------|--------------------------------------|
| 0 | 4 | 3 | [4,3,2,**-7**,8,2,3,1] |
| 1 | 3 | 2 | [4,3,**-2**,-7,8,2,3,1] |
| 2 | 2 (abs of −2) | 1 | [4,**-3**,-2,-7,8,2,3,1] |
| 3 | 7 (abs of −7) | 6 | [4,-3,-2,-7,8,2,**-3**,1] |
| 4 | 8 | 7 | [4,-3,-2,-7,8,2,-3,**-1**] |
| 5 | 2 | 1 | nums[1]=−3<0 → no change |
| 6 | 3 (abs of −3) | 2 | nums[2]=−2<0 → no change |
| 7 | 1 (abs of −1) | 0 | [**-4**,-3,-2,-7,8,2,-3,-1] |

Pass 2 — positive slots: index 4 (`8>0` → value **5**) and index 5 (`2>0` → value **6**). Result: `[5, 6]`. ✔

---

## Key Takeaways

- **Values bounded by `[1, n]` in an `n`-length array ⇒ index-as-hash.** Map value `v` to slot `v−1` and let the array store its own membership info.
- **The sign bit is free scratch space.** Negating in place records presence without a second array; always read the *absolute value* to recover the datum you overwrote the sign of.
- Same pattern solves the sibling problems: find duplicates (LC #442), find the single duplicate (LC #287), first missing positive (LC #41). Recognise the family.
- If mutation is forbidden, fall back to the O(n)-space boolean/set version — both are O(n) time.

---

## Related Problems

- LeetCode #442 — Find All Duplicates in an Array (same negation trick, collect on the *second* stamp)
- LeetCode #287 — Find the Duplicate Number (values in `[1, n]`, cycle detection)
- LeetCode #41 — First Missing Positive (index-as-hash with cyclic swaps)
- LeetCode #268 — Missing Number (single missing value via XOR / sum)
- LeetCode #645 — Set Mismatch (one duplicate + one missing)
