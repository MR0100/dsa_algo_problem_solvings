# 0300 — Longest Increasing Subsequence

> LeetCode #300 · Difficulty: Medium
> **Categories:** Array, Binary Search, Dynamic Programming

---

## Problem Statement

Given an integer array `nums`, return *the length of the longest strictly increasing subsequence*.

**Example 1:**

```
Input: nums = [10,9,2,5,3,7,101,18]
Output: 4
Explanation: The longest increasing subsequence is [2,3,7,101], therefore the length is 4.
```

**Example 2:**

```
Input: nums = [0,1,0,3,2,3]
Output: 4
```

**Example 3:**

```
Input: nums = [7,7,7,7,7,7,7]
Output: 1
```

**Constraints:**

- `1 <= nums.length <= 2500`
- `-10^4 <= nums[i] <= 10^4`

**Follow up:** Can you come up with an algorithm that runs in `O(n log n)` time complexity?

---

## Company Frequency

| Company    | Frequency        | Last Reported |
|------------|------------------|---------------|
| Amazon     | ★★★★★ Very High  | 2024          |
| Google     | ★★★★☆ High       | 2024          |
| Microsoft  | ★★★★☆ High       | 2023          |
| Meta       | ★★★☆☆ Medium     | 2023          |
| Bloomberg  | ★★★☆☆ Medium     | 2023          |

> ⚠️ Frequency data is crowd-sourced from LeetCode Discuss, Glassdoor, and community interview reports. Treat as a signal, not a guarantee.

---

## DSA Concepts Used

- **1-D Dynamic Programming** — `dp[i]` = LIS length ending at `i` → see [`/dsa/dynamic_programming_1d.md`](/dsa/dynamic_programming_1d.md)
- **Binary Search (lower bound)** — patience-sorting placement in the `tails` array → see [`/dsa/binary_search.md`](/dsa/binary_search.md)
- **Greedy** — keeping each length's tail as small as possible → see [`/dsa/greedy.md`](/dsa/greedy.md)

---

## Approaches Overview

| # | Approach | Time | Space | When to use |
|---|----------|------|-------|-------------|
| 1 | DP (ending-at-i) | O(n²) | O(n) | n ≤ few thousand; also recovers actual subsequence easily |
| 2 | Patience + Binary Search | O(n log n) | O(n) | Large n; the follow-up answer |

---

## Approach 1 — Dynamic Programming O(n²)

### Intuition
Let `dp[i]` be the length of the longest strictly increasing subsequence that **ends at** index `i`. Any such subsequence is formed by appending `nums[i]` to a shorter one ending at some earlier, smaller element. So `dp[i] = 1 + max(dp[j])` over `j < i` with `nums[j] < nums[i]`, defaulting to `1`.

### Algorithm
1. Initialise `dp[i] = 1` for all `i`.
2. For `i` from `1..n-1`, for each `j < i`: if `nums[j] < nums[i]`, set `dp[i] = max(dp[i], dp[j]+1)`.
3. Return the maximum entry in `dp`.

### Complexity
- **Time:** O(n²) — every pair `(j, i)` inspected.
- **Space:** O(n) — the `dp` array.

### Code
```go
func dpQuadratic(nums []int) int {
	n := len(nums)
	if n == 0 {
		return 0
	}
	dp := make([]int, n)
	best := 1
	for i := 0; i < n; i++ {
		dp[i] = 1
		for j := 0; j < i; j++ {
			if nums[j] < nums[i] && dp[j]+1 > dp[i] {
				dp[i] = dp[j] + 1
			}
		}
		if dp[i] > best {
			best = dp[i]
		}
	}
	return best
}
```

### Dry Run
Example 1: `nums = [10, 9, 2, 5, 3, 7, 101, 18]`.

| i | nums[i] | j's with nums[j] < nums[i] | dp[i] | best |
|---|---------|----------------------------|-------|------|
| 0 | 10      | —                          | 1     | 1    |
| 1 | 9       | —                          | 1     | 1    |
| 2 | 2       | —                          | 1     | 1    |
| 3 | 5       | 2 (dp=1)                   | 2     | 2    |
| 4 | 3       | 2 (dp=1)                   | 2     | 2    |
| 5 | 7       | 2,5,3 (max dp 2)           | 3     | 3    |
| 6 | 101     | all smaller (max dp 3)     | 4     | 4    |
| 7 | 18      | 10,9,2,5,3,7 (max dp 3)    | 4     | 4    |

Answer = **4**.

---

## Approach 2 — Patience Sorting + Binary Search O(n log n) (Optimal)

### Intuition
Maintain `tails`, where `tails[k]` is the **smallest possible tail** of any increasing subsequence of length `k+1`. Keeping tails as small as possible leaves maximum room to extend later. For each number, binary-search the first tail `>= x`: overwrite it (keeps the tail minimal) or, if none exists, append `x` (the LIS grew). The final `len(tails)` is the answer. Note: `tails` is not itself a valid LIS, but its length is exact.

### Algorithm
1. `tails = []`.
2. For each `x`: find the leftmost index `i` with `tails[i] >= x` (lower bound).
   - if `i == len(tails)`, append `x` (extends the longest run).
   - else set `tails[i] = x` (smaller tail for that length).
3. Return `len(tails)`.

### Complexity
- **Time:** O(n log n) — a binary search per element.
- **Space:** O(n) — the `tails` array.

### Code
```go
func patienceBinarySearch(nums []int) int {
	tails := []int{}
	for _, x := range nums {
		i := sort.SearchInts(tails, x)
		if i == len(tails) {
			tails = append(tails, x)
		} else {
			tails[i] = x
		}
	}
	return len(tails)
}
```

### Dry Run
Example 1: `nums = [10, 9, 2, 5, 3, 7, 101, 18]`. `sort.SearchInts` finds the leftmost tail `>= x`.

| x   | search index in tails | action            | tails after |
|-----|-----------------------|-------------------|-------------|
| 10  | 0 (empty)             | append            | [10] |
| 9   | 0 (10>=9)             | overwrite tails[0]| [9] |
| 2   | 0 (9>=2)              | overwrite tails[0]| [2] |
| 5   | 1 (none >=5)          | append            | [2,5] |
| 3   | 1 (5>=3)              | overwrite tails[1]| [2,3] |
| 7   | 2 (none >=7)          | append            | [2,3,7] |
| 101 | 3 (none >=101)        | append            | [2,3,7,101] |
| 18  | 3 (101>=18)           | overwrite tails[3]| [2,3,7,18] |

`len(tails) = 4` → answer **4**. (The final `tails` `[2,3,7,18]` is a valid LIS here, but in general only its length is guaranteed correct.)

---

## Key Takeaways
- The DP definition **"LIS ending at i"** is the standard framing; the answer is the max over all end positions, not `dp[n-1]`.
- The O(n log n) speedup replaces the inner "find the best predecessor" scan with a **binary search over tails** — a greedy invariant: smallest tail per length.
- Use **lower bound** (`SearchInts`, first index `>= x`) for *strictly* increasing; use **upper bound** (first index `> x`) for *non-decreasing* subsequences.
- `tails` gives the correct length but not necessarily an actual subsequence — recovering the sequence needs parent pointers.

---

## Related Problems
- LeetCode #354 — Russian Doll Envelopes (2-D LIS via sort + LIS)
- LeetCode #673 — Number of Longest Increasing Subsequences
- LeetCode #674 — Longest Continuous Increasing Subsequence
- LeetCode #646 — Maximum Length of Pair Chain
- LeetCode #1143 — Longest Common Subsequence
