# 0209 — Minimum Size Subarray Sum

> LeetCode #209 · Difficulty: Medium
> **Categories:** Array, Sliding Window, Two Pointers, Prefix Sum, Binary Search

---

## Problem Statement

Given an array of positive integers `nums` and a positive integer `target`, return *the **minimal length** of a subarray whose sum is greater than or equal to* `target`. If there is no such subarray, return `0` instead.

A **subarray** is a contiguous non-empty sequence of elements within an array.

**Example 1:**
```
Input: target = 7, nums = [2,3,1,2,4,3]
Output: 2
Explanation: The subarray [4,3] has the minimal length under the problem constraint.
```

**Example 2:**
```
Input: target = 4, nums = [1,4,4]
Output: 1
```

**Example 3:**
```
Input: target = 11, nums = [1,1,1,1,1,1,1,1]
Output: 0
```

**Constraints:**
- `1 <= target <= 10⁹`
- `1 <= nums.length <= 10⁵`
- `1 <= nums[i] <= 10⁴`

**Follow-up:** If you have figured out the `O(n)` solution, try coding another solution of which the time complexity is `O(n log(n))`.

---

## Company Frequency

| Company       | Frequency        | Last Reported |
|---------------|------------------|---------------|
| Meta          | ★★★★☆ High       | 2024          |
| Amazon        | ★★★★☆ High       | 2024          |
| Google        | ★★★☆☆ Medium     | 2024          |
| Microsoft     | ★★★☆☆ Medium     | 2023          |
| Goldman Sachs | ★★★☆☆ Medium     | 2023          |
| Apple         | ★★☆☆☆ Low        | 2023          |
| Bloomberg     | ★★☆☆☆ Low        | 2023          |

> ⚠️ Frequency data is crowd-sourced from LeetCode Discuss, Glassdoor, and
> community interview reports. Treat as a signal, not a guarantee.

---

## DSA Concepts Used

- **Sliding Window (variable size)** — all-positive elements make the window sum monotone under both edge moves, the exact precondition for the shrink-while-valid pattern → see [`/dsa/sliding_window.md`](/dsa/sliding_window.md)
- **Two Pointers** — `left` and `right` only ever move forward, giving amortised O(n) → see [`/dsa/two_pointers.md`](/dsa/two_pointers.md)
- **Prefix Sum** — `sum(i..j) = prefix[j+1] − prefix[i]` converts range-sum questions into pair-of-indices questions → see [`/dsa/prefix_sum.md`](/dsa/prefix_sum.md)
- **Binary Search (lower bound)** — positivity makes the prefix array strictly increasing, so the needed end index is found with `sort.SearchInts` → see [`/dsa/binary_search.md`](/dsa/binary_search.md)

---

## Approaches Overview

| # | Approach | Time | Space | When to use |
|---|----------|------|-------|-------------|
| 1 | Brute Force (Try Every Start, Extend) | O(n²) | O(1) | Baseline; fine for n in the hundreds, times out at n = 10⁵ |
| 2 | Prefix Sums + Binary Search | O(n log n) | O(n) | The follow-up answer; also the template that survives when negatives appear (→ #862 needs a monotone deque instead) |
| 3 | Sliding Window (Optimal) | O(n) | O(1) | Default answer; requires all-positive elements |

---

## Approach 1 — Brute Force (Try Every Start, Extend)

### Intuition
A subarray is determined by its two endpoints — so enumerate them. For each start `i`, extend the end `j` rightwards while keeping a running sum. Because every element is positive, the sum can only grow as `j` advances, so the *first* `j` where `sum ≥ target` yields the shortest qualifying subarray that starts at `i`; extending further only lengthens it, so break immediately. (The truly naive version recomputes each window's sum from scratch — O(n³); carrying the running sum removes one factor of n for free.)

### Algorithm
1. Set `best = n + 1` (sentinel: "nothing found").
2. For every start `i` from `0` to `n-1`:
   1. Reset `sum = 0`.
   2. For `j` from `i` to `n-1`: add `nums[j]` to `sum`.
   3. The first time `sum >= target`, update `best = min(best, j-i+1)` and break to the next `i`.
3. Return `0` if `best` is still the sentinel, else `best`.

### Complexity
- **Time:** O(n²) — n choices of start, each scanning up to n elements (early break helps but not asymptotically).
- **Space:** O(1) — only `sum`, `best`, and loop indices.

### Code
```go
func bruteForce(target int, nums []int) int {
	n := len(nums)
	best := n + 1 // sentinel: longer than any real subarray means "not found yet"

	for i := 0; i < n; i++ { // try every start position
		sum := 0
		for j := i; j < n; j++ { // extend the end one element at a time
			sum += nums[j] // running sum of nums[i..j]
			if sum >= target {
				if j-i+1 < best {
					best = j - i + 1 // shorter qualifying window found
				}
				break // any longer window from this i is worse — next start
			}
		}
	}

	if best == n+1 {
		return 0 // no subarray ever reached target
	}
	return best
}
```

### Dry Run (Example 1: target = 7, nums = [2,3,1,2,4,3])

| `i` (start) | Extension trace (sum after each `j`) | First `sum ≥ 7` | Window / length | `best` after |
|-------------|--------------------------------------|------------------|-----------------|--------------|
| 0 | 2, 5, 6, 8 | j=3 (sum 8) | [2,3,1,2] len 4 | 4 |
| 1 | 3, 4, 6, 10 | j=4 (sum 10) | [3,1,2,4] len 4 | 4 |
| 2 | 1, 3, 7 | j=4 (sum 7) | [1,2,4] len 3 | 3 |
| 3 | 2, 6, 9 | j=5 (sum 9) | [2,4,3] len 3 | 3 |
| 4 | 4, 7 | j=5 (sum 7) | [4,3] len 2 | **2** |
| 5 | 3 | never | — | 2 |

`best = 2` ≠ sentinel → return **2** ✓

---

## Approach 2 — Prefix Sums + Binary Search

### Intuition
`sum(i..j) = prefix[j+1] − prefix[i]`, where `prefix[k]` is the sum of the first `k` elements. Because every `nums[k] ≥ 1`, `prefix` is **strictly increasing** — a sorted array, which invites binary search. Fixing the start `i`, the window `i..e-1` qualifies iff `prefix[e] ≥ prefix[i] + target`; the shortest such window uses the **smallest** such `e`, i.e. a lower-bound search. n starts × log n search = the O(n log n) solution the follow-up requests.

### Algorithm
1. Build `prefix[0..n]` with `prefix[0] = 0`, `prefix[k+1] = prefix[k] + nums[k]`.
2. For each start `i` in `0..n-1`:
   1. Compute `need = prefix[i] + target`.
   2. Lower-bound search: `e = sort.SearchInts(prefix, need)` — the smallest index with `prefix[e] ≥ need`.
   3. If `e ≤ n`, window `nums[i..e-1]` qualifies with length `e - i`; update the best.
3. Return `0` if no start produced a window, else the best length.

### Complexity
- **Time:** O(n log n) — O(n) to build prefixes, then n binary searches of O(log n) each.
- **Space:** O(n) — the prefix-sum array of n+1 entries.

### Code
```go
func prefixSumBinarySearch(target int, nums []int) int {
	n := len(nums)
	// prefix[k] = sum of the first k elements; strictly increasing since nums[i] ≥ 1.
	prefix := make([]int, n+1)
	for k := 0; k < n; k++ {
		prefix[k+1] = prefix[k] + nums[k]
	}

	best := n + 1 // sentinel meaning "no valid window found yet"
	for i := 0; i < n; i++ {
		need := prefix[i] + target // window i..e-1 works iff prefix[e] ≥ need
		// Lower bound: smallest e with prefix[e] ≥ need (searches the whole
		// array; results e ≤ i are impossible since that window would be empty
		// or negative-length, and prefix[e] < need there anyway).
		e := sort.SearchInts(prefix, need)
		if e <= n { // found a real prefix index → window nums[i..e-1] qualifies
			if e-i < best {
				best = e - i // record the shorter window length
			}
		}
	}

	if best == n+1 {
		return 0 // target unreachable from any start
	}
	return best
}
```

### Dry Run (Example 1: target = 7, nums = [2,3,1,2,4,3])

`prefix = [0, 2, 5, 6, 8, 12, 15]` (indices 0..6).

| `i` | `need = prefix[i]+7` | Lower-bound `e` (first `prefix[e] ≥ need`) | Window length `e−i` | `best` after |
|-----|----------------------|---------------------------------------------|----------------------|--------------|
| 0 | 7 | e=4 (prefix[4]=8) | 4 | 4 |
| 1 | 9 | e=5 (prefix[5]=12) | 4 | 4 |
| 2 | 12 | e=5 (prefix[5]=12) | 3 | 3 |
| 3 | 13 | e=6 (prefix[6]=15) | 3 | 3 |
| 4 | 15 | e=6 (prefix[6]=15) | 2 | **2** |
| 5 | 19 | e=7 > n → no window | — | 2 |

`best = 2` → return **2** ✓

---

## Approach 3 — Sliding Window (Optimal)

### Intuition
Positivity makes the window sum **monotone in both edges**: moving `right` rightwards can only increase the sum, moving `left` rightwards can only decrease it. So neither pointer ever needs to back up. Expand `right` until the window is valid (`sum ≥ target`); then shrink from `left` while validity survives, recording the length at each valid state — that finds the *tightest* window for every right edge. Each element enters the window once and leaves at most once, so the whole thing is amortised O(n) with O(1) memory. This monotonicity is precisely what breaks with negative numbers (then you need #862's monotone-deque technique).

### Algorithm
1. Initialise `left = 0`, `sum = 0`, `best = n + 1`.
2. For `right` from `0` to `n-1`:
   1. `sum += nums[right]` (expand).
   2. While `sum >= target`: update `best = min(best, right-left+1)`, then `sum -= nums[left]`, `left++` (shrink).
3. Return `0` if `best` is still the sentinel, else `best`.

### Complexity
- **Time:** O(n) — `right` advances n times; `left` advances at most n times across the entire run (amortised constant per step).
- **Space:** O(1) — two indices plus two accumulators, independent of n.

### Code
```go
func slidingWindow(target int, nums []int) int {
	n := len(nums)
	best := n + 1 // sentinel: "no valid window seen"
	sum := 0      // sum of the current window nums[left..right]
	left := 0     // left edge of the window

	for right := 0; right < n; right++ {
		sum += nums[right] // expand: pull nums[right] into the window
		// Shrink while still valid — finds the tightest window ending at right.
		for sum >= target {
			if right-left+1 < best {
				best = right - left + 1 // new shortest qualifying window
			}
			sum -= nums[left] // expel the leftmost element
			left++            // window's left edge moves right
		}
	}

	if best == n+1 {
		return 0 // total array sum < target — impossible
	}
	return best
}
```

### Dry Run (Example 1: target = 7, nums = [2,3,1,2,4,3])

| `right` | Added | `sum` after add | Shrink steps (`sum ≥ 7`)? | `left` after | Window after | `best` |
|---------|-------|-----------------|----------------------------|--------------|--------------|--------|
| 0 | 2 | 2 | no | 0 | [2] | — |
| 1 | 3 | 5 | no | 0 | [2,3] | — |
| 2 | 1 | 6 | no | 0 | [2,3,1] | — |
| 3 | 2 | 8 | 8≥7: record len 4 → best=4; drop 2 → sum 6; 6<7 stop | 1 | [3,1,2] | 4 |
| 4 | 4 | 10 | 10≥7: record len 4 (best stays 4); drop 3 → sum 7. 7≥7: record len 3 → best=**3**; drop 1 → sum 6; stop | 3 | [2,4] | 3 |
| 5 | 3 | 9 | 9≥7: record len 3 (best stays 3); drop 2 → sum 7. 7≥7: record len 2 → best=**2**; drop 4 → sum 3; stop | 5 | [3] | **2** |

Loop ends; `best = 2` → return **2** ✓

---

## Key Takeaways

- **Shrink-while-valid template:** for "shortest subarray satisfying a condition" with monotone validity — expand right; `while (valid) { record; shrink left }`. Its mirror ("longest subarray", e.g. #3) shrinks *while invalid* and records after the shrink loop. Knowing which loop the recording sits in is the whole pattern.
- **Positivity is the licence for the sliding window.** All `nums[i] ≥ 1` ⇒ window sum strictly monotone in both edges ⇒ pointers never back up. If negatives are allowed, the window breaks — reach for prefix sums + monotone deque (#862) instead.
- **Prefix sums turn subarray sums into index pairs**, and *increasing* prefix sums unlock binary search. `sort.SearchInts(a, x)` is Go's lower bound: the smallest index `i` with `a[i] ≥ x`.
- **Sentinel trick:** initialise the best length to `n+1` (impossible value); a final equality check doubles as the "no answer → 0" branch, avoiding a separate found-flag.
- Watch the boundary: the requirement is `sum ≥ target`, not `>` — off-by-one in that comparison silently fails Example 2.
- Amortised analysis in one line: each element is added once and removed at most once, so the two nested loops together do ≤ 2n work — nested loops do not automatically mean O(n²).

---

## Related Problems

- LeetCode #3 — Longest Substring Without Repeating Characters (mirror-image window: longest, shrink-while-invalid)
- LeetCode #76 — Minimum Window Substring (shortest window, validity via character counts)
- LeetCode #713 — Subarray Product Less Than K (same positive-array window, product instead of sum)
- LeetCode #862 — Shortest Subarray with Sum at Least K (this problem **with negatives** — window fails, monotone deque required)
- LeetCode #560 — Subarray Sum Equals K (prefix sums + hash map when equality, not ≥, is asked)
- LeetCode #325 — Maximum Size Subarray Sum Equals k (longest window via first-occurrence prefix map)
