# 0163 — Missing Ranges

> LeetCode #163 · Difficulty: Easy (Premium)
> **Categories:** Array, Intervals

---

## Problem Statement

You are given an inclusive range `[lower, upper]` and a **sorted unique** integer array `nums`, where all elements are within the inclusive range.

A number `x` is considered **missing** if `x` is in the range `[lower, upper]` and `x` is not in `nums`.

Return the **shortest sorted** list of ranges that **exactly covers all the missing numbers**. That is, no element of `nums` is included in any of the ranges, and each missing number is covered by one of the ranges.

**Example 1:**
```
Input: nums = [0,1,3,50,75], lower = 0, upper = 99
Output: [[2,2],[4,49],[51,74],[76,99]]
Explanation: The ranges are:
[2,2]
[4,49]
[51,74]
[76,99]
```

**Example 2:**
```
Input: nums = [-1], lower = -1, upper = -1
Output: []
Explanation: There are no missing ranges since there are no missing numbers.
```

**Constraints:**
- `-10^9 <= lower <= upper <= 10^9`
- `0 <= nums.length <= 100`
- `lower <= nums[i] <= upper`
- All the values of `nums` are **unique**.

---

## Company Frequency

| Company   | Frequency       | Last Reported |
|-----------|-----------------|---------------|
| Google    | ★★★★★ Very High | 2024          |
| Meta      | ★★★★☆ High      | 2024          |
| Amazon    | ★★★☆☆ Medium    | 2023          |
| Microsoft | ★★★☆☆ Medium    | 2023          |
| Apple     | ★★☆☆☆ Low       | 2023          |
| Oracle    | ★★☆☆☆ Low       | 2023          |

> ⚠️ Frequency data is crowd-sourced from LeetCode Discuss, Glassdoor, and community interview reports. Treat as a signal, not a guarantee.

---

## DSA Concepts Used

- **Intervals** — the answer is a set of disjoint inclusive ranges built from the gaps between covered points → see [`/dsa/intervals.md`](/dsa/intervals.md)
- **Hash Map (Set)** — the brute force uses a set for O(1) "is x present?" membership tests during the sweep → see [`/dsa/hash_map.md`](/dsa/hash_map.md)
- **Sorting (sorted-input exploitation)** — the optimal scan works only because `nums` is sorted and unique, so gaps appear in order and exactly once → see [`/dsa/sorting.md`](/dsa/sorting.md)

---

## Approaches Overview

| # | Approach                              | Time     | Space | When to use                                                        |
|---|---------------------------------------|----------|-------|--------------------------------------------------------------------|
| 1 | Brute Force (Hash Set + Range Sweep)  | O(n + U) | O(n)  | Only when `U = upper − lower + 1` is small; span may be ~2·10⁹     |
| 2 | Linear Scan Over nums (Optimal)       | O(n)     | O(1)  | Always — work depends on the array length, not the numeric span    |

---

## Approach 1 — Brute Force (Hash Set + Full Range Sweep)

### Intuition
Apply the definition of "missing" literally. Put every element of `nums` into a hash set, then walk `x` through every value from `lower` to `upper`, asking "is `x` present?". Consecutive missing values must be merged into one range to satisfy the "shortest list" requirement, so keep a flag for "am I currently inside a missing stretch?" and remember where the stretch began. The fatal flaw: the loop length is the *numeric span* `upper − lower + 1`, which the constraints allow to reach about 2·10⁹ even when `nums` has ≤ 100 elements.

### Algorithm
1. Insert every element of `nums` into set `present`.
2. Initialise `inRange = false` (no missing stretch open).
3. For `x` from `lower` to `upper`:
   1. If `x` is **missing** and no stretch is open: open one with `start = x`, `inRange = true`.
   2. If `x` is **present** and a stretch is open: close it by appending `[start, x−1]`, set `inRange = false`.
4. After the sweep, if a stretch is still open, append the final `[start, upper]`.
5. Return the collected ranges.

### Complexity
- **Time:** O(n + U), where `U = upper − lower + 1` — building the set is O(n) and the sweep touches every value in the span once. Up to ~2·10⁹ iterations in the worst case, i.e. impractical.
- **Space:** O(n) — the hash set of present values (output excluded).

### Code
```go
func bruteForce(nums []int, lower, upper int) [][]int {
	present := make(map[int]bool, len(nums)) // O(1) membership tests
	for _, v := range nums {
		present[v] = true
	}
	result := [][]int{}
	start := 0       // first number of the currently open missing stretch
	inRange := false // whether a missing stretch is currently open
	for x := lower; x <= upper; x++ {
		if !present[x] { // x is missing
			if !inRange { // a new missing stretch begins here
				start = x
				inRange = true
			}
		} else if inRange { // x is present → the open stretch ended at x−1
			result = append(result, []int{start, x - 1})
			inRange = false
		}
	}
	if inRange { // the range ended while a stretch was still open
		result = append(result, []int{start, upper})
	}
	return result
}
```

### Dry Run
`nums = [0,1,3,50,75]`, `lower = 0`, `upper = 99` (Example 1). Set = `{0,1,3,50,75}`; only the sweep's state-changing steps shown:

| x        | present[x]? | inRange before | action                                  | result so far                       |
|----------|-------------|----------------|-----------------------------------------|-------------------------------------|
| 0        | yes         | false          | nothing to close                        | []                                  |
| 1        | yes         | false          | nothing to close                        | []                                  |
| 2        | no          | false          | open stretch, `start=2`                 | []                                  |
| 3        | yes         | true           | close `[2, 2]`                          | [[2,2]]                             |
| 4        | no          | false          | open stretch, `start=4`                 | [[2,2]]                             |
| 5…49     | no          | true           | stretch stays open                      | [[2,2]]                             |
| 50       | yes         | true           | close `[4, 49]`                         | [[2,2],[4,49]]                      |
| 51       | no          | false          | open stretch, `start=51`                | [[2,2],[4,49]]                      |
| 75       | yes         | true           | close `[51, 74]`                        | [[2,2],[4,49],[51,74]]              |
| 76       | no          | false          | open stretch, `start=76`                | [[2,2],[4,49],[51,74]]              |
| 77…99    | no          | true           | stretch stays open; sweep ends          | [[2,2],[4,49],[51,74]]              |
| post-loop| —           | true           | close final `[76, 99]`                  | **[[2,2],[4,49],[51,74],[76,99]]** ✅ |

---

## Approach 2 — Linear Scan Over nums (Optimal)

### Intuition
Flip the perspective: instead of asking "which numbers are missing?", ask "where are the gaps around the numbers that are *present*?". Since `nums` is sorted, unique, and fully inside `[lower, upper]`, the missing numbers form at most `n + 1` contiguous blocks: before the first element, between each consecutive pair, and after the last element. Track `next` — the smallest value of the range not yet accounted for. Each element `v` either sits exactly at `next` (no gap) or leaves the missing block `[next, v−1]` before it. Every gap is emitted in O(1) no matter how many billions of numbers it spans.

### Algorithm
1. Set `next = lower` — the smallest value still unaccounted for.
2. For each `v` in `nums` (in order):
   1. If `v > next`, the whole block `[next, v−1]` is missing → append it.
   2. Set `next = v + 1` — everything up to and including `v` is now covered.
3. After the loop, if `next <= upper`, append the tail block `[next, upper]`.
4. Return the collected ranges. (Empty `nums` naturally yields the single range `[lower, upper]`.)

### Complexity
- **Time:** O(n) — one pass; each of the ≤ 100 elements does constant work, regardless of the 2·10⁹ numeric span.
- **Space:** O(1) — one integer of state (`next`), output excluded.

### Code
```go
func linearScan(nums []int, lower, upper int) [][]int {
	result := [][]int{}
	next := lower // smallest number in [lower, upper] not yet covered
	for _, v := range nums {
		if v > next { // gap [next, v-1] is entirely missing
			result = append(result, []int{next, v - 1})
		}
		next = v + 1 // v itself is present → coverage advances past it
	}
	if next <= upper { // tail gap after the last element of nums
		result = append(result, []int{next, upper})
	}
	return result
}
```

### Dry Run
`nums = [0,1,3,50,75]`, `lower = 0`, `upper = 99` (Example 1):

| step      | v  | next (before) | v > next? | emitted range | next (after) | result so far                          |
|-----------|----|---------------|-----------|---------------|--------------|----------------------------------------|
| 1         | 0  | 0             | no        | —             | 1            | []                                     |
| 2         | 1  | 1             | no        | —             | 2            | []                                     |
| 3         | 3  | 2             | yes       | [2, 2]        | 4            | [[2,2]]                                |
| 4         | 50 | 4             | yes       | [4, 49]       | 51           | [[2,2],[4,49]]                         |
| 5         | 75 | 51            | yes       | [51, 74]      | 76           | [[2,2],[4,49],[51,74]]                 |
| post-loop | —  | 76            | 76 ≤ 99   | [76, 99]      | —            | **[[2,2],[4,49],[51,74],[76,99]]** ✅   |

For Example 2 (`nums = [-1]`, `lower = -1`, `upper = -1`): `next = -1`; `v = -1` is not `> next`, so nothing is emitted and `next` becomes `0`; post-loop `0 ≤ -1` is false → **[]** ✅.

---

## Key Takeaways

- **Iterate over the data, not the domain.** When the numeric span (10⁹-scale) dwarfs the input size (≤ 100), any per-value sweep is doomed; per-element gap emission is the scalable formulation.
- The `next` (a.k.a. `prev + 1`) cursor is a reusable micro-pattern for sorted-coverage problems: it makes the three gap cases — head gap, middle gaps, tail gap — fall out of one uniform loop plus one post-loop check.
- Sorted + unique input is what allows O(1) state; if duplicates were allowed you would first dedupe or guard against `v == next − 1` re-emission.
- Watch arithmetic at the extremes: with `lower = -2^31` or `upper = 2^31 − 1` (older 32-bit statement), `prev = lower − 1` or `v + 1` can overflow 32-bit ints. Go's 64-bit `int` absorbs this; in other languages widen the type. Formulating with `next = lower` (instead of `prev = lower − 1`) sidesteps the underflow entirely.
- The premium problem has an older variant that returns strings (`"2"`, `"4->49"`); the interval logic is identical — only the formatting of each emitted `[start, end]` pair changes.

---

## Related Problems

- LeetCode #228 — Summary Ranges (the mirror problem: summarise the *present* numbers instead of the missing ones)
- LeetCode #268 — Missing Number (single missing value in `[0, n]`)
- LeetCode #41 — First Missing Positive (find the smallest absent value)
- LeetCode #57 — Insert Interval (maintaining disjoint sorted intervals)
- LeetCode #352 — Data Stream as Disjoint Intervals (building coverage intervals online)
