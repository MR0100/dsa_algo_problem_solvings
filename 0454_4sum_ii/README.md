# 0454 — 4Sum II

> LeetCode #454 · Difficulty: Medium
> **Categories:** Array, Hash Table, Meet in the Middle

---

## Problem Statement

Given four integer arrays `nums1`, `nums2`, `nums3`, and `nums4` all of length `n`, return the number of tuples `(i, j, k, l)` such that:

- `0 <= i, j, k, l < n`
- `nums1[i] + nums2[j] + nums3[k] + nums4[l] == 0`

**Example 1:**

```
Input: nums1 = [1,2], nums2 = [-2,-1], nums3 = [-1,2], nums4 = [0,2]
Output: 2
Explanation:
The two tuples are:
1. (0, 0, 0, 1) -> nums1[0] + nums2[0] + nums3[0] + nums4[1] = 1 + (-2) + (-1) + 2 = 0
2. (1, 1, 0, 0) -> nums1[1] + nums2[1] + nums3[0] + nums4[0] = 2 + (-1) + (-1) + 0 = 0
```

**Example 2:**

```
Input: nums1 = [0], nums2 = [0], nums3 = [0], nums4 = [0]
Output: 1
```

**Constraints:**

- `n == nums1.length == nums2.length == nums3.length == nums4.length`
- `1 <= n <= 200`
- `-2^28 <= nums1[i], nums2[i], nums3[i], nums4[i] <= 2^28`

---

## Company Frequency

| Company    | Frequency        | Last Reported |
|------------|------------------|---------------|
| Amazon     | ★★★★☆ High       | 2024          |
| Google     | ★★★☆☆ Medium     | 2023          |
| Meta       | ★★★☆☆ Medium     | 2023          |
| Microsoft  | ★★☆☆☆ Low        | 2023          |
| Bloomberg  | ★★☆☆☆ Low        | 2022          |

> ⚠️ Frequency data is crowd-sourced from LeetCode Discuss, Glassdoor, and
> community interview reports. Treat as a signal, not a guarantee.

---

## DSA Concepts Used

- **Hash Map (complement counting)** — the engine of the fast solution: store how many index pairs produce each first-half sum, then look up the negated second-half sum. Counting complements in a map is the same idea that powers Two Sum, generalised to *pair* sums → see [`/dsa/hash_map.md`](/dsa/hash_map.md)
- **Two Pointers / Split-and-combine** — the "count pairs on each side, then match" strategy is the counting cousin of the two-pointer 4Sum split; here we split `4 = 2 + 2` and combine the halves rather than fix-and-scan → see [`/dsa/two_pointers.md`](/dsa/two_pointers.md)

---

## Approaches Overview

| # | Approach | Time | Space | When to use |
|---|----------|------|-------|-------------|
| 1 | Brute Force (four nested loops) | O(n⁴) | O(1) | Ground-truth only; `n=200` ⇒ 1.6×10⁹ iterations, TLE |
| 2 | Two Hash Maps, Split 2+2 (Meet in the Middle, Optimal) | O(n²) | O(n²) | The intended answer; two O(n²) passes with complement lookup |

---

## Approach 1 — Brute Force (Four Nested Loops)

### Intuition

The problem *is* a counting statement over four indices, so enumerate all `n⁴` index tuples and tally the ones summing to zero. Faithful and obviously correct — but at `n = 200` that is `1.6 × 10⁹` iterations, far past the time limit. It exists to validate the optimal solution's output.

### Algorithm

1. Set `count = 0`.
2. For every `a ∈ nums1`, `b ∈ nums2`, `c ∈ nums3`, `d ∈ nums4`: if `a + b + c + d == 0`, `count++`.
3. Return `count`.

### Complexity

- **Time:** O(n⁴) — four nested loops.
- **Space:** O(1) — only a counter.

### Code

```go
func bruteForce(nums1, nums2, nums3, nums4 []int) int {
	count := 0
	for _, a := range nums1 {
		for _, b := range nums2 {
			for _, c := range nums3 {
				for _, d := range nums4 {
					if a+b+c+d == 0 { // this quadruple sums to zero
						count++ // every index tuple is counted separately
					}
				}
			}
		}
	}
	return count
}
```

### Dry Run

Example 1: `nums1=[1,2]`, `nums2=[-2,-1]`, `nums3=[-1,2]`, `nums4=[0,2]`. Only the tuples that hit 0 are shown (16 total are examined).

| a | b | c | d | a+b+c+d | zero? |
|---|---|---|---|---------|-------|
| 1 | -2 | -1 | 2 | 0 | ✔ count=1 |
| 2 | -1 | -1 | 0 | 0 | ✔ count=2 |
| ... | ... | ... | ... | ≠ 0 | the other 14 tuples miss |

Result: `2` ✔.

---

## Approach 2 — Two Hash Maps, Split 2+2 (Meet in the Middle, Optimal)

### Intuition

Rewrite the target as a balance between two halves:

```
a + b + c + d == 0   ⇔   (a + b) == -(c + d).
```

So precompute every possible left-half sum `a + b` (a from `nums1`, b from `nums2`) and remember, in a hash map, *how many index pairs* produce each sum. Then walk every right-half sum `c + d` and add the stored count of `-(c + d)`: each recorded left pair combines with the current `(c, d)` to make one valid zero-sum quadruple. Two `O(n²)` passes replace the `O(n⁴)` enumeration — meet-in-the-middle, and the map counts *multiplicities* so repeated sums are handled automatically.

### Algorithm

1. Build `sumAB`: for each `a ∈ nums1`, `b ∈ nums2`, increment `sumAB[a+b]`.
2. Set `count = 0`.
3. For each `c ∈ nums3`, `d ∈ nums4`: add `sumAB[-(c+d)]` to `count` (absent key contributes 0).
4. Return `count`.

### Complexity

- **Time:** O(n²) — one double loop to fill the map, one to query it; map operations are O(1) average.
- **Space:** O(n²) — up to `n²` distinct pair sums stored.

### Code

```go
func meetInTheMiddle(nums1, nums2, nums3, nums4 []int) int {
	// sumAB[s] = number of (i, j) index pairs with nums1[i] + nums2[j] == s.
	sumAB := make(map[int]int, len(nums1)*len(nums2))
	for _, a := range nums1 {
		for _, b := range nums2 {
			sumAB[a+b]++ // record one more pair achieving this sum
		}
	}

	count := 0
	for _, c := range nums3 {
		for _, d := range nums4 {
			// We need a+b == -(c+d) to reach a total of 0. Every stored pair
			// with that sum pairs with the current (c,d) to form one quadruple.
			count += sumAB[-(c + d)] // missing key yields 0, which is correct
		}
	}
	return count
}
```

### Dry Run

Example 1: `nums1=[1,2]`, `nums2=[-2,-1]`, `nums3=[-1,2]`, `nums4=[0,2]`.

Build `sumAB` from all `a+b`:

| a | b | a+b | sumAB after |
|---|---|-----|-------------|
| 1 | -2 | -1 | `{-1:1}` |
| 1 | -1 | 0 | `{-1:1, 0:1}` |
| 2 | -2 | 0 | `{-1:1, 0:2}` |
| 2 | -1 | 1 | `{-1:1, 0:2, 1:1}` |

Now scan `c+d` and add `sumAB[-(c+d)]`:

| c | d | c+d | need -(c+d) | sumAB[need] | count |
|---|---|-----|-------------|-------------|-------|
| -1 | 0 | -1 | 1 | 1 | 1 |
| -1 | 2 | 1 | -1 | 1 | 2 |
| 2 | 0 | 2 | -2 | 0 | 2 |
| 2 | 2 | 4 | -4 | 0 | 2 |

Result: `2` ✔ — matches the brute-force count.

---

## Key Takeaways

- **`k`-Sum with independent arrays ⇒ split into halves and meet in the middle.** For four arrays, split `4 = 2 + 2`: hash all `n²` left sums, then look up complements for the `n²` right sums, turning `O(n^k)` into `O(n^{k/2})`.
- **Store counts, not booleans.** The map value is a *multiplicity* — how many index pairs reach a sum — because the problem counts index tuples, and duplicate values/sums must each be counted.
- **`map[int]int` zero-value is your friend in Go:** a missing key returns `0`, so `count += sumAB[-(c+d)]` needs no "exists" check.
- This is the same complement-lookup idea as **Two Sum** (#1), lifted one level: instead of matching a single number's complement, match a *pair sum's* complement.

---

## Related Problems

- LeetCode #1 — Two Sum (complement lookup, the seed idea)
- LeetCode #18 — 4Sum (one array, distinct quadruplets, two-pointer split)
- LeetCode #15 — 3Sum (fix one, two-pointer the rest)
- LeetCode #16 — 3Sum Closest (nearest sum to a target)
- LeetCode #653 — Two Sum IV (complement lookup over a BST)
