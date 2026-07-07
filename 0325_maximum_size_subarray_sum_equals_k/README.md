# 0325 — Maximum Size Subarray Sum Equals k

> LeetCode #325 · Difficulty: Medium (Premium)
> **Categories:** Array, Hash Table, Prefix Sum

---

## Problem Statement

Given an integer array `nums` and an integer `k`, return the maximum length of a
subarray that sums to `k`. If there is not one, return `0` instead.

**Example 1:**

```
Input:  nums = [1,-1,5,-2,3], k = 3
Output: 4
Explanation: The subarray [1, -1, 5, -2] sums to 3 and is the longest.
```

**Example 2:**

```
Input:  nums = [-2,-1,2,1], k = 1
Output: 2
Explanation: The subarray [-1, 2] sums to 1 and is the longest.
```

**Constraints:**

- `1 <= nums.length <= 2 * 10^5`
- `-10^4 <= nums[i] <= 10^4`
- `-10^9 <= k <= 10^9`
- The sum of the entire `nums` array is guaranteed to fit within the 32-bit
  signed integer range.

---

## Company Frequency

| Company   | Frequency        | Last Reported |
|-----------|------------------|---------------|
| Facebook  | ★★★★☆ High       | 2024          |
| Amazon    | ★★★☆☆ Medium     | 2023          |
| Google    | ★★★☆☆ Medium     | 2023          |
| Palantir  | ★★☆☆☆ Low        | 2022          |

> ⚠️ Frequency data is crowd-sourced from LeetCode Discuss, Glassdoor, and
> community interview reports. Treat as a signal, not a guarantee.

---

## DSA Concepts Used

- **Prefix Sum** — a subarray sum is the difference of two prefix sums → see
  [`/dsa/prefix_sum.md`](/dsa/prefix_sum.md)
- **Hash Map** — map each prefix sum to its earliest index for O(1) lookup of the
  complement → see [`/dsa/hash_map.md`](/dsa/hash_map.md)

---

## Approaches Overview

| # | Approach | Time | Space | When to use |
|---|----------|------|-------|-------------|
| 1 | Brute Force (all subarrays) | O(n²) | O(1) | Tiny inputs / baseline |
| 2 | Prefix Sum + Hash Map (Optimal) | O(n) | O(n) | Any real input; negatives allowed |

---

## Approach 1 — Brute Force (All Subarrays)

### Intuition
Enumerate every contiguous subarray. Fix a start `i`, extend the end `j` while
keeping a running sum; whenever the sum hits `k`, its length competes for the max.

### Algorithm
1. `best = 0`.
2. For each start `i`: `sum = 0`; for each end `j >= i`: `sum += nums[j]`; if
   `sum == k`, `best = max(best, j-i+1)`.
3. Return `best`.

### Complexity
- **Time:** O(n²) — every (start, end) pair.
- **Space:** O(1) — a couple of scalars.

### Code
```go
func bruteForce(nums []int, k int) int {
	best := 0
	for i := 0; i < len(nums); i++ {
		sum := 0
		for j := i; j < len(nums); j++ {
			sum += nums[j]
			if sum == k && j-i+1 > best {
				best = j - i + 1
			}
		}
	}
	return best
}
```

### Dry Run
Example 1: `nums = [1,-1,5,-2,3]`, `k = 3`.

| i | j-range sums that equal 3          | candidate length | best |
|---|------------------------------------|------------------|------|
| 0 | j=3: 1-1+5-2 = 3                    | 4                | 4 |
| 0 | j=4: 1-1+5-2+3 = 6 (≠3)             | —                | 4 |
| 1 | none equal 3                       | —                | 4 |
| 2 | j=4: 5-2+3 = 6 (≠3)                 | —                | 4 |
| 4 | j=4: 3 = 3                          | 1                | 4 |

`best = 4`. Output `4`.

---

## Approach 2 — Prefix Sum + Hash Map (Optimal)

### Intuition
Let `P(i)` be the prefix sum through index `i`. A subarray `j+1..i` sums to
`P(i) - P(j)`. We want that equal to `k`, i.e. `P(j) = P(i) - k`. If we know the
**earliest** index where the prefix sum `P(i) - k` occurred, the subarray from
just after it to `i` is the longest one ending at `i`. Storing each prefix sum's
first index only (never overwriting) maximises length. Seed the map with
`{0: -1}` so subarrays starting at index 0 are handled uniformly. Works with
negatives (a sliding window would not).

### Algorithm
1. `seen = {0: -1}` — prefix sum 0 exists "before" index 0.
2. `sum = 0`, `best = 0`.
3. For `i, x` in `nums`: `sum += x`; if `sum-k` is in `seen`,
   `best = max(best, i - seen[sum-k])`; if `sum` is new, `seen[sum] = i`.
4. Return `best`.

### Complexity
- **Time:** O(n) — single pass, O(1) hash operations.
- **Space:** O(n) — up to n distinct prefix sums stored.

### Code
```go
func prefixSumHashMap(nums []int, k int) int {
	seen := map[int]int{0: -1}
	sum := 0
	best := 0
	for i, x := range nums {
		sum += x
		if j, ok := seen[sum-k]; ok && i-j > best {
			best = i - j
		}
		if _, ok := seen[sum]; !ok {
			seen[sum] = i
		}
	}
	return best
}
```

### Dry Run
Example 1: `nums = [1,-1,5,-2,3]`, `k = 3`. Start `seen = {0:-1}`, `sum = 0`,
`best = 0`.

| i | x  | sum | sum-k | seen has sum-k? | best (i - seen[sum-k]) | store seen[sum] |
|---|----|-----|-------|-----------------|------------------------|-----------------|
| 0 | 1  | 1   | -2    | no              | 0                      | {0:-1, 1:0} |
| 1 | -1 | 0   | -3    | no              | 0                      | (0 exists, keep -1) |
| 2 | 5  | 5   | 2     | no              | 0                      | {...,5:2} |
| 3 | -2 | 3   | 0     | yes → seen[0]=-1| 3-(-1)=**4**           | {...,3:3} |
| 4 | 3  | 6   | 3     | yes → seen[3]=3 | 4-3=1 (not > 4)        | {...,6:4} |

`best = 4`. Output `4`.

---

## Key Takeaways
- **Subarray-sum-equals-target ⇒ prefix sums in a hash map.** The complement is
  `currentPrefix - k`.
- For **longest** length, store each prefix sum's **first** occurrence (never
  overwrite); for **counting** subarrays (#560) you would instead accumulate
  frequencies.
- Seeding `{0: -1}` elegantly covers subarrays that begin at index 0.
- A sliding window does **not** work here because `nums` may contain negatives,
  so the running sum is not monotonic.

---

## Related Problems
- LeetCode #560 — Subarray Sum Equals K (count instead of longest length)
- LeetCode #523 — Continuous Subarray Sum (prefix sum mod k)
- LeetCode #974 — Subarray Sums Divisible by K (prefix sum mod)
- LeetCode #1074 — Number of Submatrices That Sum to Target (2D prefix + hashmap)
