# 0259 — 3Sum Smaller

> LeetCode #259 · Difficulty: Medium
> **Categories:** Array, Two Pointers, Sorting, Binary Search

---

## Problem Statement

Given an array of `n` integers `nums` and an integer `target`, find the number of index triplets `i`, `j`, `k` with `0 <= i < j < k < n` that satisfy the condition `nums[i] + nums[j] + nums[k] < target`.

**Example 1:**

```
Input: nums = [-2,0,1,3], target = 2
Output: 2
Explanation: Because there are two triplets which sums are less than 2:
[-2,0,1]
[-2,0,3]
```

**Example 2:**

```
Input: nums = [], target = 0
Output: 0
```

**Example 3:**

```
Input: nums = [0], target = 0
Output: 0
```

**Constraints:**

- `n == nums.length`
- `0 <= n <= 3500`
- `-100 <= nums[i] <= 100`
- `-100 <= target <= 100`

**Follow up:** Could you solve it in `O(n^2)` runtime?

---

## Company Frequency

| Company    | Frequency        | Last Reported |
|------------|------------------|---------------|
| Google     | ★★★☆☆ Medium     | 2024          |
| Facebook   | ★★★☆☆ Medium     | 2023          |
| Amazon     | ★★★☆☆ Medium     | 2023          |
| Microsoft  | ★★☆☆☆ Low        | 2022          |
| Bloomberg  | ★★☆☆☆ Low        | 2022          |

> ⚠️ Frequency data is crowd-sourced from LeetCode Discuss, Glassdoor, and
> community interview reports. Treat as a signal, not a guarantee.

---

## DSA Concepts Used

- **Two Pointers** — after sorting, a `lo`/`hi` sweep counts all valid partners for a fixed first element in one linear pass, answering the O(n²) follow-up → see [`/dsa/two_pointers.md`](/dsa/two_pointers.md)
- **Sorting** — ordering the array is what makes the "count `hi - lo` at once" shortcut valid → see [`/dsa/sorting.md`](/dsa/sorting.md)

---

## Approaches Overview

| # | Approach | Time | Space | When to use |
|---|----------|------|-------|-------------|
| 1 | Brute Force (Triple Loop) | O(n³) | O(1) | Small n; baseline to verify against |
| 2 | Sort + Two Pointers (Optimal) | O(n²) | O(1) | The follow-up answer; counts a whole range of pairs per step |

---

## Approach 1 — Brute Force (Triple Loop)

### Intuition

The problem literally counts index triples `i < j < k` whose values sum to less than `target`. Enumerate every such triple with three nested loops and count the ones that qualify.

### Algorithm

1. For every `i < j < k`, test `nums[i] + nums[j] + nums[k] < target`.
2. Increment a counter on success.
3. Return the counter.

### Complexity

- **Time:** O(n³) — three nested loops.
- **Space:** O(1) — a single counter.

### Code

```go
func bruteForce(nums []int, target int) int {
	count := 0
	n := len(nums)
	for i := 0; i < n; i++ {
		for j := i + 1; j < n; j++ {
			for k := j + 1; k < n; k++ {
				if nums[i]+nums[j]+nums[k] < target { // valid triple
					count++
				}
			}
		}
	}
	return count
}
```

### Dry Run

Example 1: `nums = [-2,0,1,3], target = 2`.

| (i,j,k) | values | sum | sum < 2? | count |
|---------|--------|-----|----------|-------|
| (0,1,2) | -2,0,1 | -1 | yes | 1 |
| (0,1,3) | -2,0,3 | 1 | yes | 2 |
| (0,2,3) | -2,1,3 | 2 | no | 2 |
| (1,2,3) | 0,1,3 | 4 | no | 2 |

Result: `2` ✔

---

## Approach 2 — Sort + Two Pointers (Optimal)

### Intuition

Sort the array first. Fix the smallest member of the triple at index `i`, and count pairs `(lo, hi)` with `i < lo < hi` whose three-way sum is below `target`. Start `lo = i+1`, `hi = n-1`. The key shortcut: if `nums[i] + nums[lo] + nums[hi] < target`, then because the array is sorted, **every** `hi'` from `lo+1` up to `hi` also satisfies the inequality (those values are ≤ `nums[hi]`). That is `hi - lo` valid pairs added in one step; then advance `lo`. If the sum is not below target, the largest partner is too big, so decrement `hi`.

### Algorithm

1. Sort `nums`.
2. For `i = 0 .. n-3`: set `lo = i+1`, `hi = n-1`.
3. While `lo < hi`:
   - If `nums[i] + nums[lo] + nums[hi] < target`: add `hi - lo` to the count, then `lo++`.
   - Else `hi--`.
4. Return the count.

### Complexity

- **Time:** O(n²) — sort is O(n log n); the outer loop runs `n` times, each with an O(n) two-pointer sweep. Answers the follow-up.
- **Space:** O(1) extra (in-place sort).

### Code

```go
func twoPointers(nums []int, target int) int {
	sort.Ints(nums) // sorting lets us count many pairs in one comparison
	count := 0
	n := len(nums)
	for i := 0; i < n-2; i++ { // fix the smallest of the triple
		lo, hi := i+1, n-1
		for lo < hi {
			if nums[i]+nums[lo]+nums[hi] < target {
				// nums[lo] with ANY hi' in (lo, hi] is also < target because
				// those values are ≤ nums[hi]; count all of them at once.
				count += hi - lo
				lo++ // move to the next (larger) middle element
			} else {
				hi-- // sum too large; the largest partner must shrink
			}
		}
	}
	return count
}
```

### Dry Run

Example 1: `nums = [-2,0,1,3]` (already sorted), `target = 2`.

| i | nums[i] | lo | hi | sum | sum < 2? | action | count |
|---|---------|----|----|-----|----------|--------|-------|
| 0 | -2 | 1 | 3 | -2+0+3 = 1 | yes | count += hi-lo = 2; lo→2 | 2 |
| 0 | -2 | 2 | 3 | -2+1+3 = 2 | no | hi→2 | 2 |
| 0 | -2 | 2 | 2 | — | — | lo == hi, exit inner | 2 |
| 1 | 0 | 2 | 3 | 0+1+3 = 4 | no | hi→2 | 2 |
| 1 | 0 | 2 | 2 | — | — | exit inner | 2 |

Result: `2` ✔ — the two pairs counted at i=0 are `(-2,0,1)` and `(-2,0,3)`.

---

## Key Takeaways

- **Sort + two pointers for "count of tuples under a threshold":** when the target is an inequality (`< target`) rather than equality, a sorted two-pointer sweep counts a *contiguous block* of partners in one move (`hi - lo`), giving O(n²) overall.
- **The `count += hi - lo` shortcut** is the distinguishing trick versus classic 3Sum: don't count one triple at a time when sortedness lets you count a whole range.
- **Fix one element, reduce to 2Sum:** the standard reduction for k-Sum problems — pin the first index and solve the two-pointer subproblem on the suffix.

---

## Related Problems

- LeetCode #15 — 3Sum (equality version, dedup with two pointers)
- LeetCode #16 — 3Sum Closest (minimize distance to target)
- LeetCode #18 — 4Sum (one more fixed index)
- LeetCode #611 — Valid Triangle Number (same "count `hi - lo`" two-pointer counting)
