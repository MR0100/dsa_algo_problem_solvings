# 0354 — Russian Doll Envelopes

> LeetCode #354 · Difficulty: Hard
> **Categories:** Array, Binary Search, Dynamic Programming, Sorting, LIS

---

## Problem Statement

You are given a 2D array of integers `envelopes` where
`envelopes[i] = [wi, hi]` represents the width and the height of an envelope.

One envelope can fit into another if and only if both the width and height of one
envelope are greater than the other envelope's width and height.

Return *the maximum number of envelopes you can Russian doll (i.e., put one
inside the other)*.

**Note:** You cannot rotate an envelope.

**Example 1:**

```
Input: envelopes = [[5,4],[6,4],[6,7],[2,3]]
Output: 3
Explanation: The maximum number of envelopes you can Russian doll is 3
([2,3] => [5,4] => [6,7]).
```

**Example 2:**

```
Input: envelopes = [[1,1],[1,1],[1,1]]
Output: 1
```

**Constraints:**

- `1 <= envelopes.length <= 10^5`
- `envelopes[i].length == 2`
- `1 <= wi, hi <= 10^5`

---

## Company Frequency

| Company    | Frequency        | Last Reported |
|------------|------------------|---------------|
| Google     | ★★★★☆ High       | 2024          |
| Amazon     | ★★★☆☆ Medium     | 2023          |
| Microsoft  | ★★★☆☆ Medium     | 2023          |
| Apple      | ★★☆☆☆ Low        | 2022          |

> ⚠️ Frequency data is crowd-sourced from LeetCode Discuss, Glassdoor, and
> community interview reports. Treat as a signal, not a guarantee.

---

## DSA Concepts Used

- **Longest Increasing Subsequence (LIS)** — the core of the problem once one
  dimension is fixed by sorting → see [`/dsa/dynamic_programming_1d.md`](/dsa/dynamic_programming_1d.md)
- **Sorting with a tie-break trick** — width asc, height desc removes the
  same-width trap → see [`/dsa/sorting.md`](/dsa/sorting.md)
- **Binary search (patience piles)** — the O(n log n) LIS uses `sort.Search`
  over the tails array → see [`/dsa/binary_search.md`](/dsa/binary_search.md)

---

## Approaches Overview

| # | Approach | Time | Space | When to use |
|---|----------|------|-------|-------------|
| 1 | Sort + O(n²) DP (LIS) | O(n²) | O(n) | Small n; easy to reason about |
| 2 | Sort (w asc, h desc) + O(n log n) LIS (Optimal) | O(n log n) | O(n) | Large n (up to 10⁵) |

---

## Approach 1 — Sort + O(n²) DP (LIS on height)

### Intuition
Sort by width ascending. Any nesting chain then reads left-to-right as a
subsequence. Since equal widths cannot nest, an envelope `j` fits inside `i` only
when **both** `w[j] < w[i]` and `h[j] < h[i]`. That is the classic LIS DP with a
2-D strict-increase comparison.

### Algorithm
1. Sort envelopes by width asc; on tie, by height asc.
2. `dp[i]` = length of the longest chain ending at envelope `i`, initialized 1.
3. For each `i`, for each `j < i`: if `w[j] < w[i]` and `h[j] < h[i]`, set
   `dp[i] = max(dp[i], dp[j] + 1)`.
4. Return `max(dp)`.

### Complexity
- **Time:** O(n²) — the nested loop over all pairs.
- **Space:** O(n) for the `dp` array.

### Code
```go
func dpQuadratic(envelopes [][]int) int {
	n := len(envelopes)
	if n == 0 {
		return 0
	}
	sort.Slice(envelopes, func(i, j int) bool {
		if envelopes[i][0] != envelopes[j][0] {
			return envelopes[i][0] < envelopes[j][0]
		}
		return envelopes[i][1] < envelopes[j][1]
	})

	dp := make([]int, n) // dp[i] = longest nesting chain ending at i
	best := 0
	for i := 0; i < n; i++ {
		dp[i] = 1
		for j := 0; j < i; j++ {
			if envelopes[j][0] < envelopes[i][0] && envelopes[j][1] < envelopes[i][1] {
				if dp[j]+1 > dp[i] {
					dp[i] = dp[j] + 1
				}
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
Input `[[5,4],[6,4],[6,7],[2,3]]`. After sort by width asc, height asc:
`[[2,3],[5,4],[6,4],[6,7]]` (indices 0..3).

| i | envelope | j candidates (both dims smaller) | dp[i] | best |
|---|----------|----------------------------------|-------|------|
| 0 | [2,3] | none | 1 | 1 |
| 1 | [5,4] | j=0 [2,3]: 2<5 & 3<4 ✓ → dp=2 | 2 | 2 |
| 2 | [6,4] | j=0 ✓ dp=2; j=1 [5,4]: 4<4 ✗ | 2 | 2 |
| 3 | [6,7] | j=0 ✓; j=1 [5,4] ✓ dp=3; j=2 [6,4]: 6<6 ✗ | 3 | 3 |

Answer **3**, via `[2,3] → [5,4] → [6,7]`.

---

## Approach 2 — Sort (w asc, h desc) + O(n log n) LIS (Optimal)

### Intuition
Fix width by sorting it ascending. Naively running LIS on height would let two
envelopes of the **same width** (heights increasing) falsely nest. Sorting height
**descending** within equal widths makes those heights non-increasing, so a
strictly-increasing LIS can never pick two of them. The problem collapses to:
longest strictly-increasing subsequence of the height array — solved in
O(n log n) by patience sorting.

### Algorithm
1. Sort by width asc; on tie, height **descending**.
2. Maintain `tails`, where `tails[k]` is the smallest possible tail height of an
   increasing subsequence of length `k+1`.
3. For each height `h`: binary-search the first `tails[i] >= h`. If none, append
   `h` (extends the longest chain); else overwrite that slot with `h`.
4. Return `len(tails)`.

### Complexity
- **Time:** O(n log n) — sort plus one binary search per envelope.
- **Space:** O(n) for `tails`.

### Code
```go
func lisBinarySearch(envelopes [][]int) int {
	n := len(envelopes)
	if n == 0 {
		return 0
	}
	sort.Slice(envelopes, func(i, j int) bool {
		if envelopes[i][0] != envelopes[j][0] {
			return envelopes[i][0] < envelopes[j][0]
		}
		return envelopes[i][1] > envelopes[j][1] // height DESC on tie
	})

	tails := []int{}
	for _, e := range envelopes {
		h := e[1]
		lo := sort.Search(len(tails), func(i int) bool { return tails[i] >= h })
		if lo == len(tails) {
			tails = append(tails, h)
		} else {
			tails[lo] = h
		}
	}
	return len(tails)
}
```

### Dry Run
Input `[[5,4],[6,4],[6,7],[2,3]]`. Sort by width asc, height desc →
`[[2,3],[5,4],[6,7],[6,4]]`. Heights stream: `3, 4, 7, 4`.

| h | binary-search (first tail ≥ h) | action | tails |
|---|--------------------------------|--------|-------|
| 3 | len 0 → append | append 3 | [3] |
| 4 | no tail ≥ 4 → append | append 4 | [3,4] |
| 7 | no tail ≥ 7 → append | append 7 | [3,4,7] |
| 4 | tails[1]=4 ≥ 4 → overwrite idx 1 | tails[1]=4 | [3,4,7] |

`len(tails) = 3`. Note how sorting `6,7` before `6,4` prevented the two width-6
envelopes from both entering the chain.

---

## Key Takeaways

- **Reduce 2-D nesting to 1-D LIS** by sorting one dimension.
- **The tie-break is the whole trick.** Sorting height *descending* within equal
  widths silently enforces the strict-width rule, so a plain strictly-increasing
  LIS on heights is correct.
- **Patience sorting** turns LIS into O(n log n): `tails[k]` holds the minimal
  tail of a length-`k+1` increasing subsequence; `sort.Search` finds the slot.
- `len(tails)` is the LIS length even though `tails` itself is not a real
  subsequence.

---

## Related Problems

- LeetCode #300 — Longest Increasing Subsequence (the 1-D base case)
- LeetCode #646 — Maximum Length of Pair Chain (interval-chain LIS variant)
- LeetCode #1691 — Maximum Height by Stacking Cuboids (3-D generalization)
- LeetCode #368 — Largest Divisible Subset (chain DP)
