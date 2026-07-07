# 0410 — Split Array Largest Sum

> LeetCode #410 · Difficulty: Hard
> **Categories:** Array, Binary Search, Dynamic Programming, Greedy, Prefix Sum

---

## Problem Statement

Given an integer array `nums` and an integer `k`, split `nums` into `k` non-empty subarrays such that the largest sum of any subarray is **minimized**.

Return *the minimized largest sum of the split*.

A **subarray** is a contiguous part of the array.

**Example 1:**

```
Input: nums = [7,2,5,10,8], k = 2
Output: 18
Explanation: There are four ways to split nums into two subarrays.
The best way is to split it into [7,2,5] and [10,8], where the largest sum among the two subarrays is only 18.
```

**Example 2:**

```
Input: nums = [1,2,3,4,5], k = 2
Output: 9
Explanation: There are four ways to split nums into two subarrays.
The best way is to split it into [1,2,3] and [4,5], where the largest sum among the two subarrays is only 9.
```

**Constraints:**

- `1 <= nums.length <= 1000`
- `0 <= nums[i] <= 10^6`
- `1 <= k <= min(50, nums.length)`

---

## Company Frequency

| Company    | Frequency        | Last Reported |
|------------|------------------|---------------|
| Google     | ★★★★☆ High       | 2024          |
| Amazon     | ★★★★☆ High       | 2024          |
| Meta       | ★★★☆☆ Medium     | 2023          |
| Microsoft  | ★★★☆☆ Medium     | 2023          |
| Apple      | ★★☆☆☆ Low        | 2023          |
| ByteDance  | ★★☆☆☆ Low        | 2023          |

> ⚠️ Frequency data is crowd-sourced from LeetCode Discuss, Glassdoor, and
> community interview reports. Treat as a signal, not a guarantee.

---

## DSA Concepts Used

- **Binary Search on the Answer** — the optimal method searches over the *value* of the largest allowed subarray sum, exploiting that "is a cap feasible with ≤ k parts?" is monotone → see [`/dsa/binary_search.md`](/dsa/binary_search.md)
- **Greedy (feasibility check)** — for a fixed cap, greedily fill each part until it would overflow, then cut; this counts the minimum number of parts needed → see [`/dsa/greedy.md`](/dsa/greedy.md)
- **Partition Dynamic Programming** — the exact DP `dp[i][j] = min over p of max(dp[p][j-1], sum(p..i))` splits prefixes into a fixed number of parts → see [`/dsa/dynamic_programming_2d.md`](/dsa/dynamic_programming_2d.md)
- **Prefix Sum** — O(1) subarray-sum lookups power both the greedy check and the DP transition → see [`/dsa/prefix_sum.md`](/dsa/prefix_sum.md)

---

## Approaches Overview

| # | Approach | Time | Space | When to use |
|---|----------|------|-------|-------------|
| 1 | Binary Search on the Answer (Optimal) | O(n · log(sum − max)) | O(1) | The expected answer; fast and short |
| 2 | Dynamic Programming (Partition DP) | O(n² · k) | O(n · k) | Exact, no monotonicity insight needed; teaches the partition-DP pattern |

`n = len(nums)`.

---

## Approach 1 — Binary Search on the Answer (Optimal)

### Intuition

Don't search for the *split* — search for the *value* of the answer. The largest subarray sum must lie in `[max(nums), sum(nums)]`:

- it can't be less than the biggest single element (that element occupies at least a whole part by itself);
- it can't exceed the total (taking `k = 1` dumps everything into one part).

Now the key monotonicity: **"can we split so that no part exceeds `cap`?"** gets *easier* as `cap` grows. If a cap works, every larger cap works too. So the set of feasible caps is a suffix of the range, and we binary-search for its smallest element.

Feasibility for a fixed `cap` is a one-pass greedy: keep adding numbers to the current part; the moment the next number would push it over `cap`, start a new part. The number of parts this produces is the *minimum* possible for that cap — if it is ≤ `k`, the cap is feasible.

### Algorithm

1. `lo = max(nums)`, `hi = sum(nums)`.
2. While `lo < hi`:
   - `mid = lo + (hi − lo)/2`;
   - if `canSplit(mid)` (needs ≤ `k` parts), set `hi = mid`;
   - else set `lo = mid + 1`.
3. Return `lo`.

`canSplit(cap)`: walk `nums`, greedily filling a part; cut a new part when it would overflow `cap`; return `parts ≤ k`.

### Complexity

- **Time:** O(n · log(sum − max)) — each feasibility scan is `O(n)`; the search runs `O(log(range))` times.
- **Space:** O(1) — just a few counters.

### Code

```go
func binarySearch(nums []int, k int) int {
	lo, hi := 0, 0
	for _, v := range nums {
		if v > lo {
			lo = v // lower bound: no part can be smaller than the largest element
		}
		hi += v // upper bound: everything in one part
	}

	// canSplit reports whether nums can be cut into <= k parts each summing <= cap.
	canSplit := func(cap int) bool {
		parts := 1    // we always have at least one part
		current := 0  // running sum of the part we are currently filling
		for _, v := range nums {
			if current+v > cap {
				// v doesn't fit in the current part → start a new part with v.
				parts++
				current = v
				if parts > k {
					return false // needed more than k parts ⇒ cap too small
				}
			} else {
				current += v // v fits, keep filling this part
			}
		}
		return true
	}

	// Standard "find smallest feasible value" binary search.
	for lo < hi {
		mid := lo + (hi-lo)/2 // candidate cap (avoids overflow)
		if canSplit(mid) {
			hi = mid // mid works; try to do even smaller
		} else {
			lo = mid + 1 // mid too small; need a bigger cap
		}
	}
	return lo // smallest cap that is feasible = minimized largest sum
}
```

### Dry Run

Example 1: `nums = [7,2,5,10,8]`, `k = 2`. `lo = max = 10`, `hi = sum = 32`.

| lo | hi | mid | canSplit(mid)? (greedy parts) | move |
|----|----|-----|-------------------------------|------|
| 10 | 32 | 21 | [7,2,5]=14, [10,8]=18 → 2 parts ≤ 2 ✓ | hi = 21 |
| 10 | 21 | 15 | [7,2,5]=14, [10]=10, [8]=8 → 3 parts > 2 ✗ | lo = 16 |
| 16 | 21 | 18 | [7,2,5]=14, [10,8]=18 → 2 parts ≤ 2 ✓ | hi = 18 |
| 16 | 18 | 17 | [7,2,5]=14, [10]=10, [8] → 3 parts > 2 ✗ | lo = 18 |
| 18 | 18 | — | loop ends (lo == hi) | — |

Return `lo = 18`.

Result: **18** ✔ (split `[7,2,5] | [10,8]`).

---

## Approach 2 — Dynamic Programming (Partition DP)

### Intuition

Solve it exactly with a partition DP. Let `dp[i][j]` = the minimized largest-subarray-sum when splitting the **first `i` elements** into exactly `j` parts. To build `dp[i][j]`, decide where the **last** part begins: say it covers elements `(p .. i-1]`. Then:

- the last part's sum is `prefix[i] − prefix[p]`;
- the first `p` elements are split into `j − 1` parts optimally, costing `dp[p][j-1]`.

The largest part of this whole arrangement is the **max** of those two values. We pick the split point `p` that minimizes that max:

```
dp[i][j] = min over p in [j-1, i-1] of  max( dp[p][j-1], prefix[i] - prefix[p] )
```

Base case `dp[0][0] = 0`. The answer is `dp[n][k]`.

### Algorithm

1. Build `prefix[i] = nums[0] + … + nums[i-1]`.
2. Init `dp[i][j] = +inf`, `dp[0][0] = 0`.
3. For `i` from `1..n`, for `j` from `1..min(i,k)`, for split point `p` from `j-1..i-1`:
   `dp[i][j] = min(dp[i][j], max(dp[p][j-1], prefix[i] − prefix[p]))`.
4. Return `dp[n][k]`.

### Complexity

- **Time:** O(n² · k) — `O(n·k)` states, each scanning up to `n` split points.
- **Space:** O(n · k) — the DP table.

### Code

```go
func dpBottomUp(nums []int, k int) int {
	n := len(nums)
	const inf = int(1e18) // sentinel for "not reachable"

	// prefix[i] = nums[0] + ... + nums[i-1].
	prefix := make([]int, n+1)
	for i := 0; i < n; i++ {
		prefix[i+1] = prefix[i] + nums[i]
	}

	// dp[i][j] = min largest-part sum splitting first i elems into j parts.
	dp := make([][]int, n+1)
	for i := range dp {
		dp[i] = make([]int, k+1)
		for j := range dp[i] {
			dp[i][j] = inf // start unreachable
		}
	}
	dp[0][0] = 0 // zero elements, zero parts, zero cost

	for i := 1; i <= n; i++ {
		// Can't have more parts than elements, nor more than k.
		for j := 1; j <= k && j <= i; j++ {
			// Try every split point p: last part is (p..i-1], value prefix[i]-prefix[p].
			for p := j - 1; p < i; p++ {
				if dp[p][j-1] == inf {
					continue // first p elements can't be split into j-1 parts
				}
				lastPart := prefix[i] - prefix[p]   // sum of the last (j-th) part
				candidate := dp[p][j-1]             // largest part among the first j-1
				if lastPart > candidate {
					candidate = lastPart // the overall largest part is the max of the two
				}
				if candidate < dp[i][j] {
					dp[i][j] = candidate // keep the split that minimizes the largest part
				}
			}
		}
	}
	return dp[n][k]
}
```

### Dry Run

Example 1: `nums = [7,2,5,10,8]`, `k = 2`. `prefix = [0,7,9,14,24,32]`.

We want `dp[5][2]`. First the single-part row `dp[i][1] = prefix[i]` (whole prefix in one part):

`dp[1][1]=7, dp[2][1]=9, dp[3][1]=14, dp[4][1]=24, dp[5][1]=32`.

Now `dp[5][2]` = min over split point `p ∈ [1,4]` of `max(dp[p][1], prefix[5] − prefix[p])`:

| p | dp[p][1] (first part group) | last part = 32 − prefix[p] | max |
|---|-----------------------------|----------------------------|-----|
| 1 | 7 | 32 − 7 = 25 | 25 |
| 2 | 9 | 32 − 9 = 23 | 23 |
| 3 | 14 | 32 − 14 = 18 | **18** |
| 4 | 24 | 32 − 24 = 8 | 24 |

Minimum is `18` at `p = 3` (last part `[10,8]`, first group `[7,2,5]`).

Result: `dp[5][2] = 18` ✔

---

## Key Takeaways

- **"Minimize the maximum" (or maximize the minimum) ⇒ binary-search the answer.** Whenever the objective is a bottleneck value and "can we achieve bound X?" is monotone, search over X instead of over configurations. This pattern covers Koko eating bananas, ship packages in D days, and more.
- **The feasibility check is the real work.** Here it's a greedy one-pass part-counter; getting its direction right (feasible ⇒ shrink `hi`) is what makes the binary search converge to the *smallest* feasible cap.
- **Bounds matter:** `lo = max(nums)` (a part can't be smaller than its biggest element) and `hi = sum(nums)` (one part holds everything). Starting `lo` at `0` also works but wastes iterations.
- **Partition DP is the exact fallback:** `dp[i][j] = min_p max(dp[p][j-1], sum(p..i))` is the canonical "split a sequence into k pieces optimizing an aggregate" recurrence — memorise its shape; it recurs across many problems.
- Binary search is `O(n log range)` vs the DP's `O(n²k)` — for `n` up to 1000 the binary search is dramatically faster, which is why it's the expected solution.

---

## Related Problems

- LeetCode #1011 — Capacity To Ship Packages Within D Days (identical binary-search-on-answer)
- LeetCode #875 — Koko Eating Bananas (minimize max via feasibility search)
- LeetCode #1231 — Divide Chocolate (maximize the minimum piece — the dual)
- LeetCode #1335 — Minimum Difficulty of a Job Schedule (partition DP, same recurrence)
- LeetCode #813 — Largest Sum of Averages (partition DP variant)
