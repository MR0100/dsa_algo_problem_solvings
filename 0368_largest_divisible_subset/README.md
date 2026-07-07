# 0368 — Largest Divisible Subset

> LeetCode #368 · Difficulty: Medium
> **Categories:** Array, Math, Dynamic Programming, Sorting

---

## Problem Statement

Given a set of **distinct** positive integers `nums`, return the largest subset `answer` such that every pair `(answer[i], answer[j])` of elements in this subset satisfies:

- `answer[i] % answer[j] == 0`, or
- `answer[j] % answer[i] == 0`.

If there are multiple solutions, return any of them.

**Example 1:**

```
Input: nums = [1,2,3]
Output: [1,2]
Explanation: [1,3] is also accepted.
```

**Example 2:**

```
Input: nums = [1,2,4,8]
Output: [1,2,4,8]
```

**Constraints:**

- `1 <= nums.length <= 1000`
- `1 <= nums[i] <= 2 * 10^9`
- All the integers in `nums` are **unique**.

---

## Company Frequency

| Company    | Frequency        | Last Reported |
|------------|------------------|---------------|
| Google     | ★★★☆☆ Medium     | 2023          |
| Amazon     | ★★★☆☆ Medium     | 2023          |
| Microsoft  | ★★☆☆☆ Low        | 2022          |
| Adobe      | ★★☆☆☆ Low        | 2022          |

> ⚠️ Frequency data is crowd-sourced from LeetCode Discuss, Glassdoor, and
> community interview reports. Treat as a signal, not a guarantee.

---

## DSA Concepts Used

- **Dynamic Programming (1D, LIS-shaped)** — `dp[i]` = longest divisible chain ending at `i`; the recurrence is the Longest Increasing Subsequence pattern with "divides" replacing "<" → see [`/dsa/longest_increasing_subsequence.md`](/dsa/longest_increasing_subsequence.md)
- **Sorting** — sorting ascending makes divisibility transitive along the array, so only consecutive-in-subset pairs must divide → see [`/dsa/sorting.md`](/dsa/sorting.md)

---

## Approaches Overview

| # | Approach | Time | Space | When to use |
|---|----------|------|-------|-------------|
| 1 | DP with Parent Pointers (Optimal) | O(n²) | O(n) | The interview answer: O(n) memory, reconstruct via parents |
| 2 | DP Storing Full Subsets | O(n²)–O(n³) | O(n²) | Easier to read; stores each subset directly at the cost of memory |

---

## Approach 1 — DP with Parent Pointers (Optimal)

### Intuition

Divisibility on a **sorted** array behaves like `≤`: if `a | b` and `b | c` then `a | c`. So once we sort ascending, any subset in which each element divides the *next larger chosen* element is automatically pairwise-divisible — we only need consecutive pairs to divide. That reduces the task to "longest chain where `nums[j]` divides `nums[i]` for `j < i`", a Longest-Increasing-Subsequence-shaped DP. `dp[i]` is the length of the best chain ending at `i`, and `parent[i]` remembers which earlier index we extended so we can rebuild the actual subset.

### Algorithm

1. Sort `nums` ascending.
2. Initialise `dp[i] = 1`, `parent[i] = -1`.
3. For each `i`, for each `j < i`: if `nums[i] % nums[j] == 0` and `dp[j]+1 > dp[i]`, set `dp[i] = dp[j]+1`, `parent[i] = j`.
4. Track `best` = index of the maximum `dp`.
5. Follow `parent` pointers from `best` back to the start, then reverse to ascending.

### Complexity

- **Time:** O(n²) — the double loop dominates the O(n log n) sort.
- **Space:** O(n) — the `dp` and `parent` arrays.

### Code

```go
func dpWithParents(nums []int) []int {
	n := len(nums)
	if n == 0 {
		return []int{}
	}
	sort.Ints(nums) // divisibility acts like ≤ once sorted ascending

	dp := make([]int, n)     // dp[i] = longest divisible chain ending at i
	parent := make([]int, n) // parent[i] = previous index in that chain
	best := 0                // index where the overall-longest chain ends

	for i := 0; i < n; i++ {
		dp[i] = 1       // a single element is always a valid chain
		parent[i] = -1  // no predecessor yet
		for j := 0; j < i; j++ {
			// nums[j] < nums[i] (sorted); if it divides nums[i] we may extend
			// the chain that ended at j.
			if nums[i]%nums[j] == 0 && dp[j]+1 > dp[i] {
				dp[i] = dp[j] + 1
				parent[i] = j
			}
		}
		if dp[i] > dp[best] { // remember the global best end-point
			best = i
		}
	}

	// Reconstruct by following parent pointers from best back to the start.
	var result []int
	for i := best; i != -1; i = parent[i] {
		result = append(result, nums[i])
	}
	// We collected largest→smallest; reverse to ascending for a tidy answer.
	for l, r := 0, len(result)-1; l < r; l, r = l+1, r-1 {
		result[l], result[r] = result[r], result[l]
	}
	return result
}
```

### Dry Run

Example 1: `nums = [1,2,3]` (already sorted).

| i | nums[i] | j scanned | divisible? | dp[i] | parent[i] | best |
|---|---------|-----------|------------|-------|-----------|------|
| 0 | 1 | — | — | 1 | -1 | 0 |
| 1 | 2 | j=0 (1): 2%1==0 | yes, dp=1+1=2 | 2 | 0 | 1 |
| 2 | 3 | j=0 (1): 3%1==0 → dp=2; j=1 (2): 3%2≠0 | partial | 2 | 0 | 1 |

`best = 1` (dp=2). Walk parents from index 1: `nums[1]=2`, parent→0 `nums[0]=1`, parent→-1. Collected `[2,1]`, reversed → `[1,2]` ✔

---

## Approach 2 — DP Storing Full Subsets

### Intuition

Same recurrence, but instead of parent pointers we store the entire best subset that ends at each index. `subsets[i]` holds the largest divisible subset ending at `nums[i]`; to extend, copy the best predecessor's subset and append `nums[i]`. Simpler to read, at the cost of O(n²) memory and slice copying.

### Algorithm

1. Sort ascending.
2. `subsets[i]` starts as `[nums[i]]`.
3. For each `i`, for each `j < i` with `nums[i] % nums[j] == 0`: if `len(subsets[j])+1 > len(subsets[i])`, set `subsets[i]` to a copy of `subsets[j]` plus `nums[i]`.
4. Return the longest `subsets[i]`.

### Complexity

- **Time:** O(n²) comparisons plus O(n) per copy — up to O(n³) copying in the worst case.
- **Space:** O(n²) — each `subsets[i]` may hold up to `n` elements.

### Code

```go
func dpFullSubsets(nums []int) []int {
	n := len(nums)
	if n == 0 {
		return []int{}
	}
	sort.Ints(nums)

	subsets := make([][]int, n) // subsets[i] = best divisible subset ending at i
	best := 0
	for i := 0; i < n; i++ {
		subsets[i] = []int{nums[i]} // at minimum, the element by itself
		for j := 0; j < i; j++ {
			if nums[i]%nums[j] == 0 && len(subsets[j])+1 > len(subsets[i]) {
				// Copy predecessor's subset so we don't alias/mutate it.
				cp := make([]int, len(subsets[j]))
				copy(cp, subsets[j])
				subsets[i] = append(cp, nums[i])
			}
		}
		if len(subsets[i]) > len(subsets[best]) {
			best = i
		}
	}
	return subsets[best]
}
```

### Dry Run

Example 1: `nums = [1,2,3]`.

| i | nums[i] | j check | subsets[i] after | best |
|---|---------|---------|------------------|------|
| 0 | 1 | — | `[1]` | 0 |
| 1 | 2 | j=0: 2%1==0, len 1+1>1 | `[1,2]` | 1 |
| 2 | 3 | j=0: 3%1==0, len 1+1>1 → `[1,3]`; j=1: 3%2≠0 | `[1,3]` | 1 (tie, keep first) |

Longest is `subsets[1] = [1,2]` ✔ (`[1,3]` would also be accepted.)

---

## Key Takeaways

- **Sort first to make divisibility transitive.** A pairwise-divisibility subset over a sorted array only needs each consecutive chosen pair to divide — turning an O(2ⁿ) subset search into an O(n²) LIS DP.
- **LIS template with a swapped predicate.** Replace `nums[j] < nums[i]` with `nums[i] % nums[j] == 0`; everything else (dp array, parent pointers, reconstruction) is standard LIS.
- **Parent pointers vs. stored subsets:** parents keep memory at O(n) and are the preferred reconstruction technique; storing full subsets is more readable but O(n²) memory.
- Reconstruction walks predecessors backward, so remember to reverse the collected list.

---

## Related Problems

- LeetCode #300 — Longest Increasing Subsequence (the parent DP this is built on)
- LeetCode #673 — Number of Longest Increasing Subsequence (LIS variant)
- LeetCode #1048 — Longest String Chain (chain DP with a different predicate)
- LeetCode #646 — Maximum Length of Pair Chain (sort-then-chain)
