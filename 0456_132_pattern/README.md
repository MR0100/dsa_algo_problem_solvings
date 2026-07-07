# 0456 — 132 Pattern

> LeetCode #456 · Difficulty: Medium
> **Categories:** Array, Binary Search, Stack, Monotonic Stack, Ordered Set

---

## Problem Statement

Given an array of `n` integers `nums`, a **132 pattern** is a subsequence of three integers `nums[i]`, `nums[j]` and `nums[k]` such that `i < j < k` and `nums[i] < nums[k] < nums[j]`.

Return `true` *if there is a **132 pattern** in* `nums`*, otherwise, return* `false`*.*

**Example 1:**

```
Input: nums = [1,2,3,4]
Output: false
Explanation: There is no 132 pattern in the sequence.
```

**Example 2:**

```
Input: nums = [3,1,4,2]
Output: true
Explanation: There is a 132 pattern in the sequence: [1, 4, 2].
```

**Example 3:**

```
Input: nums = [-1,3,2,0]
Output: true
Explanation: There are three 132 patterns in the sequence: [-1, 3, 2], [-1, 3, 0] and [-1, 2, 0].
```

**Constraints:**

- `n == nums.length`
- `1 <= n <= 2 * 10^5`
- `-10^9 <= nums[i] <= 10^9`

**Follow up:** Can you implement the solution in `O(n)` time complexity?

---

## Company Frequency

| Company    | Frequency        | Last Reported |
|------------|------------------|---------------|
| Amazon     | ★★★★☆ High       | 2024          |
| Google     | ★★★☆☆ Medium     | 2024          |
| Bloomberg  | ★★★☆☆ Medium     | 2023          |
| Microsoft  | ★★☆☆☆ Low        | 2023          |
| Adobe      | ★★☆☆☆ Low        | 2022          |
| ByteDance  | ★★☆☆☆ Low        | 2023          |

> ⚠️ Frequency data is crowd-sourced from LeetCode Discuss, Glassdoor, and
> community interview reports. Treat as a signal, not a guarantee.

---

## DSA Concepts Used

- **Monotonic Stack** — the optimal solution keeps a *decreasing* stack of candidate "3" values while scanning right→left; popping a smaller value promotes it to the best "2" that already has a larger "3" to its right → see [`/dsa/monotonic_stack.md`](/dsa/monotonic_stack.md)
- **Stack** — the candidate structure is a plain LIFO stack (push each element, pop while it is dominated) → see [`/dsa/stack.md`](/dsa/stack.md)
- **Arrays** — the whole problem is a subsequence-relation query over a 1-D array, and Approach 2 leans on the classic *running prefix minimum* array trick → see [`/dsa/arrays.md`](/dsa/arrays.md)

---

## Approaches Overview

| # | Approach | Time | Space | When to use |
|---|----------|------|-------|-------------|
| 1 | Brute Force (three loops) | O(n³) | O(1) | Tiny arrays / correctness oracle; TLEs at n = 2·10⁵ |
| 2 | Fix j + prefix minimum | O(n²) | O(1) | Removes the i-loop; still too slow for the largest inputs |
| 3 | Monotonic Stack, right→left (Optimal) | O(n) | O(n) | The intended answer; meets the O(n) follow-up |

---

## Approach 1 — Brute Force

### Intuition

The pattern is literally a triple of indices `i < j < k` with `nums[i] < nums[k] < nums[j]`. Enumerate every such ordered triple and test the inequality. No cleverness — just a direct transcription of the definition, useful as a correctness reference for the faster approaches.

### Algorithm

1. Loop `i` over all indices (the "1").
2. Loop `j` over indices after `i` (the "3").
3. Loop `k` over indices after `j` (the "2").
4. If `nums[i] < nums[k]` and `nums[k] < nums[j]`, return `true`.
5. If no triple qualifies, return `false`.

### Complexity

- **Time:** O(n³) — three nested loops, each up to length `n`.
- **Space:** O(1) — only loop counters.

### Code

```go
func bruteForce(nums []int) bool {
	n := len(nums)
	for i := 0; i < n; i++ { // choose the "1" (smallest, leftmost)
		for j := i + 1; j < n; j++ { // choose the "3" (largest, middle index)
			for k := j + 1; k < n; k++ { // choose the "2" (middle value, rightmost)
				// 132 means nums[i] < nums[k] < nums[j].
				if nums[i] < nums[k] && nums[k] < nums[j] {
					return true // found one qualifying triple
				}
			}
		}
	}
	return false // exhausted all triples, none matched
}
```

### Dry Run

Example 1: `nums = [1, 2, 3, 4]` (strictly increasing, so no "3 then smaller 2" ever appears).

| i | nums[i] | j | nums[j] | k | nums[k] | nums[i] < nums[k] < nums[j]? |
|---|---------|---|---------|---|---------|------------------------------|
| 0 | 1 | 1 | 2 | 2 | 3 | 1<3 but 3<2 false |
| 0 | 1 | 1 | 2 | 3 | 4 | 1<4 but 4<2 false |
| 0 | 1 | 2 | 3 | 3 | 4 | 1<4 but 4<3 false |
| 1 | 2 | 2 | 3 | 3 | 4 | 2<4 but 4<3 false |
| … | … | … | … | … | … | every case fails |

All triples fail → return `false` ✔

---

## Approach 2 — Fix j, Track Min-i So Far

### Intuition

Fix the middle element `nums[j]` as the "3". The best possible "1" for it is the *smallest value strictly to its left* — anything smaller only makes the `nums[i] < nums[k]` test easier to pass. So maintain a running prefix minimum `minLeft` as `j` sweeps left→right, and for each `j` look to the right for a "2": some `nums[k]` with `minLeft < nums[k] < nums[j]`. This collapses the `i`-loop into a single precomputed value.

### Algorithm

1. If `n < 3`, return `false`.
2. Initialise `minLeft = nums[0]` (running minimum of everything left of `j`).
3. For `j` from `1` to `n-1`:
   1. For `k` from `j+1` to `n-1`: if `minLeft < nums[k] < nums[j]`, return `true`.
   2. Update `minLeft = min(minLeft, nums[j])`.
4. Return `false`.

### Complexity

- **Time:** O(n²) — the outer `j` loop times the inner `k` scan; the `i` search is now free.
- **Space:** O(1) — one running minimum.

### Code

```go
func minPrefix(nums []int) bool {
	n := len(nums)
	if n < 3 {
		return false // need at least three elements for a triple
	}
	minLeft := nums[0] // smallest value seen strictly left of j (best "1")
	for j := 1; j < n; j++ {
		// Look for a "2": some k > j whose value sits strictly between the best
		// "1" (minLeft) and the current "3" (nums[j]).
		for k := j + 1; k < n; k++ {
			if minLeft < nums[k] && nums[k] < nums[j] {
				return true // minLeft (i) < nums[k] (k) < nums[j] (j): a 132 pattern
			}
		}
		if nums[j] < minLeft {
			minLeft = nums[j] // extend the prefix minimum for future j's
		}
	}
	return false
}
```

### Dry Run

Example 1: `nums = [1, 2, 3, 4]`, `minLeft` starts at `1`.

| j | nums[j] | minLeft (before) | k scan for minLeft < nums[k] < nums[j] | found? | minLeft (after) |
|---|---------|------------------|----------------------------------------|--------|-----------------|
| 1 | 2 | 1 | k=2→3 (3<2 no), k=3→4 (4<2 no) | no | min(1,2)=1 |
| 2 | 3 | 1 | k=3→4 (4<3 no) | no | min(1,3)=1 |
| 3 | 4 | 1 | (no k) | no | min(1,4)=1 |

No `j` finds a valid "2" → return `false` ✔

---

## Approach 3 — Monotonic Stack, Scan Right→Left (Optimal)

### Intuition

Scan from the right and treat each element as the "1". We want to have already discovered, to the right, a valid ("3", "2") pair — a larger element followed by a strictly smaller one. Keep the best such "2" in a variable `third`: it is the largest value that we know had a strictly bigger element to *its* right. Maintain a **decreasing stack** of candidate "3"s. When the current element is bigger than the stack top, that top was a smaller value sitting to the right of a bigger one (the current element) — so it is a legitimate "2"; pop it and raise `third`. Now, if at any later (further-left) step the current element drops below `third`, that element is a valid "1" and we have `i < j < k` with `nums[i] < nums[k] < nums[j]`. The right-to-left order guarantees the index relations.

### Algorithm

1. If `n < 3`, return `false`.
2. Set `third = -∞` (best "2" that already has a larger "3" to its right).
3. Use an empty stack of candidate "3" values.
4. Scan `i` from `n-1` down to `0`:
   1. If `nums[i] < third`, return `true` (`nums[i]` is a valid "1").
   2. While the stack is non-empty and `nums[i] > top`, pop the top into `third` (it is a valid "2").
   3. Push `nums[i]` as a new candidate "3".
5. Return `false`.

### Complexity

- **Time:** O(n) — every element is pushed once and popped at most once; the inner while-loop is amortised O(1).
- **Space:** O(n) — the stack can grow to `n` in the worst case (a strictly increasing-from-the-right array).

### Code

```go
func monotonicStack(nums []int) bool {
	n := len(nums)
	if n < 3 {
		return false // impossible to form a triple
	}
	third := math.MinInt64 // best "2" that already has a larger "3" to its right
	stack := []int{}       // decreasing stack of candidate "3" values (right side)
	for i := n - 1; i >= 0; i-- {
		// nums[i] plays the "1": if it is below the best known "2", we win,
		// because that "2" had a strictly larger "3" between i and it.
		if nums[i] < third {
			return true
		}
		// Everything smaller than nums[i] on the stack is a "2" that pairs with
		// nums[i] as its larger "3"; raise `third` to the largest such value.
		for len(stack) > 0 && nums[i] > stack[len(stack)-1] {
			third = stack[len(stack)-1] // this popped value is a valid "2"
			stack = stack[:len(stack)-1]
		}
		stack = append(stack, nums[i]) // nums[i] is a fresh candidate "3"
	}
	return false
}
```

### Dry Run

Example 2: `nums = [3, 1, 4, 2]`, scanning right→left. `third = -∞`, `stack = []`.

| i | nums[i] | nums[i] < third? | pops (value → third) | stack after push | third after |
|---|---------|------------------|----------------------|------------------|-------------|
| 3 | 2 | 2 < -∞? no | none | [2] | -∞ |
| 2 | 4 | 4 < -∞? no | pop 2 → third=2 | [4] | 2 |
| 1 | 1 | 1 < 2? **yes** → return true | — | — | — |

At `i = 1`, `nums[1] = 1 < third = 2`, so a 132 pattern exists — here the "1" is `nums[1]=1`, the "3" is `nums[2]=4`, and the "2" is the popped `nums[3]=2` (`1 < 2 < 4`). Return `true` ✔

---

## Key Takeaways

- **Fix the middle, precompute the extreme.** Turning O(n³) into O(n²) came from realising that for a fixed "3", the only "1" worth trying is the prefix minimum — a recurring "fix one index, reduce the others to a running aggregate" move.
- **Scan direction is a design lever.** Going *right→left* lets us commit each element as the "1" only after its potential ("3","2") pair on the right is already summarised in `third`.
- **Monotonic stack = "nearest dominating element" machine.** Here the decreasing stack surfaces, for the current large value, all smaller values to its right — exactly the "2" candidates. The same skeleton solves Next Greater Element, Daily Temperatures, and Largest Rectangle.
- **A single scalar can replace a whole structure query.** `third` compresses "the best 2 with a larger 3 to its right" into one number, which is what makes the check O(1).

---

## Related Problems

- LeetCode #496 — Next Greater Element I (monotonic stack basics)
- LeetCode #503 — Next Greater Element II (circular monotonic stack)
- LeetCode #739 — Daily Temperatures (decreasing stack of indices)
- LeetCode #84 — Largest Rectangle in Histogram (monotonic stack, area)
- LeetCode #907 — Sum of Subarray Minimums (contribution via monotonic stack)
