# 0462 — Minimum Moves to Equal Array Elements II

> LeetCode #462 · Difficulty: Medium
> **Categories:** Array, Math, Sorting, Quickselect

---

## Problem Statement

Given an integer array `nums` of size `n`, return *the minimum number of moves required to make all array elements equal*.

In one move, you can increment or decrement an element of the array by `1`.

Test cases are designed so that the answer will fit in a **32-bit** integer.

**Example 1:**

```
Input: nums = [1,2,3]
Output: 2
Explanation:
Only two moves are needed (remember each move increments or decrements one element):
[1,2,3]  =>  [2,2,3]  =>  [2,2,2]
```

**Example 2:**

```
Input: nums = [1,10,2,9]
Output: 16
```

**Constraints:**

- `n == nums.length`
- `1 <= nums.length <= 10^5`
- `-10^9 <= nums[i] <= 10^9`

---

## Company Frequency

| Company    | Frequency        | Last Reported |
|------------|------------------|---------------|
| Amazon     | ★★★☆☆ Medium     | 2024          |
| Google     | ★★★☆☆ Medium     | 2023          |
| Microsoft  | ★★☆☆☆ Low        | 2023          |
| Bloomberg  | ★★☆☆☆ Low        | 2022          |
| Adobe      | ★★☆☆☆ Low        | 2022          |

> ⚠️ Frequency data is crowd-sourced from LeetCode Discuss, Glassdoor, and
> community interview reports. Treat as a signal, not a guarantee.

---

## DSA Concepts Used

- **Median minimises Σ|xᵢ − t|** — the key math fact: the sum of absolute deviations from a point is minimised at the median (a 1-D geometric-median / L¹ result), so the optimal common target is the median value → see [`/dsa/math_number_theory.md`](/dsa/math_number_theory.md)
- **Sorting** — the simplest way to reach the median (and to pair outer elements) is to sort ascending → see [`/dsa/sorting.md`](/dsa/sorting.md)
- **Quickselect** — the median is a k-th order statistic, findable in expected O(n) without a full sort → see [`/dsa/quickselect.md`](/dsa/quickselect.md)
- **Two Pointers** — pairing the smallest and largest inward accumulates each pair's span, summing distances to the median without computing it → see [`/dsa/two_pointers.md`](/dsa/two_pointers.md)

---

## Approaches Overview

| # | Approach | Time | Space | When to use |
|---|----------|------|-------|-------------|
| 1 | Brute Force (try every target) | O(n · range) | O(1) | Baseline; TLE when value range is large |
| 2 | Sort + Median (two pointers) | O(n log n) | O(1) extra | The standard interview answer; clean and short |
| 3 | Quickselect Median (Optimal) | O(n) expected | O(n) | When O(n log n) is not enough; shows median = order statistic |

---

## Approach 1 — Brute Force (Try Every Target Value)

### Intuition

Making every element equal to a target `t` costs `Σ |nums[i] - t|` moves, since each move shifts one element by 1. The best target can never lie outside `[min, max]` (moving past the extremes only increases every distance). So enumerate each integer target in that range, total its cost, and keep the minimum. Correct and obvious, but the value range can be up to `2·10⁹`, so it exists only to anchor the idea.

### Algorithm

1. Scan once to find `lo = min(nums)` and `hi = max(nums)`.
2. For every target `t` in `[lo, hi]`, compute `cost = Σ |nums[i] - t|`.
3. Return the smallest cost.

### Complexity

- **Time:** O(n · (hi − lo)) — one array pass per candidate target.
- **Space:** O(1) — a few counters.

### Code

```go
func bruteForce(nums []int) int {
	lo, hi := nums[0], nums[0] // seed min and max with the first element
	for _, v := range nums {
		if v < lo {
			lo = v
		}
		if v > hi {
			hi = v
		}
	}
	best := -1 // sentinel: no cost computed yet
	// Every integer target in the value range is a candidate.
	for t := lo; t <= hi; t++ {
		cost := 0
		for _, v := range nums {
			cost += abs(v - t) // moves to drag v onto target t
		}
		if best == -1 || cost < best {
			best = cost // remember the cheapest target so far
		}
	}
	return best
}
```

### Dry Run

Example 1: `nums = [1,2,3]`, so `lo = 1`, `hi = 3`. Try each target:

| target t | \|1−t\| | \|2−t\| | \|3−t\| | cost | best after |
|----------|---------|---------|---------|------|------------|
| 1 | 0 | 1 | 2 | 3 | 3 |
| 2 | 1 | 0 | 1 | 2 | 2 |
| 3 | 2 | 1 | 0 | 3 | 2 |

Minimum cost is at `t = 2` (the median). Result: `2` ✔

---

## Approach 2 — Sort + Median (Two Pointers)

### Intuition

`Σ |nums[i] - t|` is minimised at the **median**. Rather than compute the median and then the sum, sort and pair the outermost elements: bringing the current smallest `nums[i]` and current largest `nums[j]` to *any* meeting point between them costs exactly `nums[j] - nums[i]` — the meeting point cancels out. All optimal meeting points lie at the median, so summing `nums[j] - nums[i]` over every nested pair yields the total distance to the median directly, and it sidesteps the odd/even-length median subtlety.

### Algorithm

1. Sort `nums` ascending.
2. Set `i = 0`, `j = n - 1`, `moves = 0`.
3. While `i < j`: add `nums[j] - nums[i]` to `moves`, then `i++`, `j--`.
4. Return `moves`.

### Complexity

- **Time:** O(n log n) — the sort dominates; the pairing pass is O(n).
- **Space:** O(1) extra beyond the sort (here we copy to avoid mutating the caller).

### Code

```go
func sortMedian(nums []int) int {
	sorted := make([]int, len(nums)) // copy so the input is not mutated
	copy(sorted, nums)
	sort.Ints(sorted) // ascending order

	moves := 0
	i, j := 0, len(sorted)-1 // outermost pair of pointers
	// Each pair contributes its span; the meeting point (the median) cancels.
	for i < j {
		moves += sorted[j] - sorted[i] // cost to converge this outer pair
		i++                            // shrink toward the center
		j--
	}
	return moves
}
```

### Dry Run

Example 1: `nums = [1,2,3]` → sorted `[1,2,3]`, `i = 0`, `j = 2`.

| Step | i | j | nums[i] | nums[j] | nums[j]−nums[i] | moves after |
|------|---|---|---------|---------|-----------------|-------------|
| 1 | 0 | 2 | 1 | 3 | 2 | 2 |
| — | 1 | 1 | — | — | i == j → stop | 2 |

Result: `2` ✔ (the single central element `2` needs no move).

---

## Approach 3 — Quickselect Median (Optimal)

### Intuition

Only the median matters, not a fully sorted array. The median is the `n/2`-th order statistic, and **quickselect** finds a k-th order statistic in expected linear time: partition around a pivot, then recurse into just the side that contains index `n/2`. With the median `m` in hand, the answer is `Σ |nums[i] - m|`. This trades the sort's `O(n log n)` for expected `O(n)`.

### Algorithm

1. `m = quickselect(nums, n/2)` — the lower median (any value between the two central elements is optimal, so the lower one is fine).
2. `moves = Σ |nums[i] - m|`.
3. Return `moves`.

### Complexity

- **Time:** O(n) expected for quickselect + O(n) for the sum = **O(n) expected** (O(n²) worst case with adversarial pivots).
- **Space:** O(n) for the working copy so the caller's slice is untouched.

### Code

```go
func quickselectMedian(nums []int) int {
	work := make([]int, len(nums)) // partitioning reorders elements; copy first
	copy(work, nums)

	median := quickselect(work, 0, len(work)-1, len(work)/2) // k-th smallest, k=n/2
	moves := 0
	for _, v := range nums {
		moves += abs(v - median) // total distance to the median
	}
	return moves
}

func quickselect(a []int, lo, hi, k int) int {
	for lo < hi {
		p := partition(a, lo, hi) // pivot lands at its final sorted index p
		switch {
		case p == k:
			return a[k] // pivot is exactly the k-th element
		case p < k:
			lo = p + 1 // target is to the right of the pivot
		default:
			hi = p - 1 // target is to the left of the pivot
		}
	}
	return a[lo] // single-element window is the answer
}

func partition(a []int, lo, hi int) int {
	pivot := a[hi] // choose the last element as pivot
	i := lo        // boundary: a[lo..i-1] are < pivot
	for j := lo; j < hi; j++ {
		if a[j] < pivot {
			a[i], a[j] = a[j], a[i] // push a smaller element into the left region
			i++
		}
	}
	a[i], a[hi] = a[hi], a[i] // drop the pivot just past the smaller region
	return i
}
```

### Dry Run

Example 1: `nums = [1,2,3]`, `n = 3`, seek `k = n/2 = 1` (0-indexed median). `work = [1,2,3]`.

Quickselect on `work[0..2]`, `k = 1`:

| Step | window [lo,hi] | pivot (a[hi]) | partition index p | decision |
|------|----------------|---------------|-------------------|----------|
| 1 | [0,2] | `3` | after partition, `3` is largest → p = 2 | p (2) > k (1) → hi = 1 |
| 2 | [0,1] | `2` | `1 < 2`, so `1` left, pivot `2` at index 1 → p = 1 | p == k → return `a[1] = 2` |

Median `m = 2`. Sum of distances: `|1−2| + |2−2| + |3−2| = 1 + 0 + 1 = 2`. Result: `2` ✔

---

## Key Takeaways

- **Minimising Σ|xᵢ − t| ⇒ pick the median.** (Minimising Σ(xᵢ − t)² would pick the mean; minimising the max deviation would pick the midrange.) Knowing which central statistic each objective wants is a reusable interview fact.
- **The meeting-point cancels for outer pairs.** After sorting, `Σ (nums[j] − nums[i])` over nested pairs equals the total distance to the median — no need to special-case odd/even length.
- **A k-th order statistic does not require a sort.** Quickselect gets it in expected O(n); use it when a single positional value (median, k-th largest) is all you need.
- Contrast with #453 (Minimum Moves I), where a move increments `n−1` elements — that reframes to `Σ(nums[i] − min)` instead of a median problem.

---

## Related Problems

- LeetCode #453 — Minimum Moves to Equal Array Elements (increment n−1 elements)
- LeetCode #296 — Best Meeting Point (2-D median, same L¹ idea)
- LeetCode #215 — Kth Largest Element in an Array (quickselect)
- LeetCode #2033 — Minimum Operations to Make a Uni-Value Grid (median with a step k)
