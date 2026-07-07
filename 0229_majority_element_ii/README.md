# 0229 — Majority Element II

> LeetCode #229 · Difficulty: Medium
> **Categories:** Array, Hash Table, Counting, Sorting

---

## Problem Statement

Given an integer array of size `n`, find all elements that appear more than `⌊ n/3 ⌋` times.

**Example 1:**
```
Input: nums = [3,2,3]
Output: [3]
```

**Example 2:**
```
Input: nums = [1]
Output: [1]
```

**Example 3:**
```
Input: nums = [1,2]
Output: [1,2]
```

**Constraints:**
- `1 <= nums.length <= 5 * 10⁴`
- `-10⁹ <= nums[i] <= 10⁹`

**Follow-up:** Could you solve the problem in linear time and in `O(1)` space?

---

## Company Frequency

| Company    | Frequency        | Last Reported |
|------------|------------------|---------------|
| Amazon     | ★★★★☆ High       | 2024          |
| Google     | ★★★★☆ High       | 2024          |
| Meta       | ★★★☆☆ Medium     | 2023          |
| Microsoft  | ★★★☆☆ Medium     | 2023          |
| Bloomberg  | ★★☆☆☆ Low        | 2023          |

> ⚠️ Frequency data is crowd-sourced from LeetCode Discuss, Glassdoor, and
> community interview reports. Treat as a signal, not a guarantee.

---

## DSA Concepts Used

- **Hash Map (frequency counting)** — the direct approach tallies each value's occurrences → see [`/dsa/hash_map.md`](/dsa/hash_map.md)
- **Boyer–Moore Voting (generalized)** — the >n/k majority pattern needs k-1 candidate/counter pairs; here k=3 gives O(1) space → see [`/dsa/greedy.md`](/dsa/greedy.md)
- **Sorting** — used only to produce deterministic output order → see [`/dsa/sorting.md`](/dsa/sorting.md)

---

## Approaches Overview

| # | Approach | Time | Space | When to use |
|---|----------|------|-------|-------------|
| 1 | Hash map count | O(n) | O(n) | Simplest; when extra space is fine |
| 2 | Boyer–Moore voting | O(n) | O(1) | Meets the follow-up: linear time, constant space |

---

## Approach 1 — Hash Map (Count and Filter)

### Intuition
"More than ⌊n/3⌋ times" is a direct frequency question. Count how many times
each value appears, then keep the values whose count clears the threshold. There
can be at most two such values, since three distinct values each exceeding n/3
would sum to more than n.

### Algorithm
1. Build a map `value → count` in one pass.
2. Compute `threshold = n / 3` (integer division gives ⌊n/3⌋).
3. Collect every value with `count > threshold`.
4. Sort the (≤2 element) result for deterministic output.

### Complexity
- **Time:** O(n) — one counting pass; sorting ≤2 elements is O(1).
- **Space:** O(n) — the map holds up to `n` distinct keys.

### Code
```go
func hashMapCount(nums []int) []int {
	counts := make(map[int]int) // value → number of occurrences
	for _, v := range nums {
		counts[v]++ // tally this value
	}
	threshold := len(nums) / 3 // must appear strictly MORE than ⌊n/3⌋ times
	res := []int{}
	for v, c := range counts {
		if c > threshold { // clears the majority-of-thirds bar
			res = append(res, v)
		}
	}
	sort.Ints(res) // map iteration order is random; sort for a stable answer
	return res
}
```

### Dry Run
`nums = [3,2,3]`, `threshold = 3/3 = 1`:

| Step | Element | counts map | 
|------|---------|-----------|
| 1    | 3       | {3:1}     |
| 2    | 2       | {3:1, 2:1}|
| 3    | 3       | {3:2, 2:1}|

Filter `count > 1`: only `3` (count 2) qualifies; `2` (count 1) does not.
Result: `[3]`. ✅

---

## Approach 2 — Boyer–Moore Voting, Generalized (Optimal)

### Intuition
At most two values can appear more than n/3 times. The classic Boyer–Moore
majority vote (for the >n/2 case) tracks *one* candidate; generalizing to >n/k
needs `k-1` candidates. Here k=3, so keep **two** candidates with two independent
counters. Each number either votes for a candidate (increment), fills an empty
slot, or cancels one vote from *both* candidates. The two survivors are only
*possible* answers — a second pass verifies each truly exceeds n/3.

### Algorithm
1. Initialize `cand1, cand2` to distinct placeholders, `count1 = count2 = 0`.
2. For each `num`:
   - if `num == cand1` → `count1++`
   - else if `num == cand2` → `count2++`
   - else if `count1 == 0` → adopt `num` as `cand1`, `count1 = 1`
   - else if `count2 == 0` → adopt `num` as `cand2`, `count2 = 1`
   - else → `count1--` and `count2--` (three-way cancellation)
3. Recount `cand1` and `cand2` over the array.
4. Output whichever confirmed count exceeds `n/3`.

### Complexity
- **Time:** O(n) — one voting pass plus one verification pass.
- **Space:** O(1) — a fixed set of scalars.

### Code
```go
func boyerMooreVoting(nums []int) []int {
	// Two candidate slots. Use distinct initial values so an all-same edge case
	// can't accidentally match an unset candidate; counts of 0 make them "empty".
	cand1, cand2 := 0, 1
	count1, count2 := 0, 0

	for _, v := range nums {
		switch {
		case v == cand1: // vote for candidate 1
			count1++
		case v == cand2: // vote for candidate 2
			count2++
		case count1 == 0: // slot 1 is empty → adopt v
			cand1, count1 = v, 1
		case count2 == 0: // slot 2 is empty → adopt v
			cand2, count2 = v, 1
		default: // v differs from both live candidates → cancel one vote each
			count1--
			count2--
		}
	}

	// The two survivors are only CANDIDATES — verify their true frequencies.
	count1, count2 = 0, 0
	for _, v := range nums {
		if v == cand1 {
			count1++
		} else if v == cand2 {
			count2++
		}
	}

	res := []int{}
	n := len(nums)
	if count1 > n/3 { // confirmed to exceed ⌊n/3⌋
		res = append(res, cand1)
	}
	if count2 > n/3 {
		res = append(res, cand2)
	}
	sort.Ints(res) // deterministic ordering
	return res
}
```

### Dry Run
`nums = [3,2,3]`, start `cand1=0,cand2=1,count1=0,count2=0`:

| Element | Branch taken            | cand1,count1 | cand2,count2 |
|---------|-------------------------|--------------|--------------|
| 3       | count1==0 → adopt 3     | 3, 1         | 1, 0         |
| 2       | count2==0 → adopt 2     | 3, 1         | 2, 1         |
| 3       | v==cand1 → count1++     | 3, 2         | 2, 1         |

Verification pass over `[3,2,3]`: cand1=3 appears 2×, cand2=2 appears 1×.
Threshold `n/3 = 1`. `2 > 1` → keep 3; `1 > 1` is false → drop 2.
Result: `[3]`. ✅

---

## Key Takeaways
- The count of elements exceeding n/k is at most `k-1` — a pigeonhole bound that
  caps the number of candidates you must track.
- Boyer–Moore generalizes: `>n/2` → 1 candidate, `>n/3` → 2 candidates, `>n/k`
  → k-1 candidates, each with its own counter, plus a mandatory verification pass.
- The verification pass is essential — surviving candidates are not guaranteed to
  actually meet the threshold (e.g., when no such element exists).
- Initialize the two candidate slots to *different* placeholder values so a
  single distinct input value can't occupy both slots.

---

## Related Problems
- LeetCode #169 — Majority Element (the >n/2 base case, single candidate)
- LeetCode #1150 — Check If a Number Is Majority Element in a Sorted Array
- LeetCode #17.10 — Find Majority Element (LCCI variant)
