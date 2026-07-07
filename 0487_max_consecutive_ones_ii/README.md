# 0487 — Max Consecutive Ones II

> LeetCode #487 · Difficulty: Medium
> **Categories:** Array, Dynamic Programming, Sliding Window

---

## Problem Statement

Given a binary array `nums`, return *the maximum number of consecutive* `1`*'s in the array if you can flip at most one* `0`.

**Example 1:**

```
Input: nums = [1,0,1,1,0]
Output: 4
Explanation:
- If we flip the first zero, nums becomes [1,1,1,1,0] and we have 3 consecutive ones.
- If we flip the second zero, nums becomes [1,0,1,1,1] and we have 3 consecutive ones.
The max number of consecutive ones is 4.
```

**Example 2:**

```
Input: nums = [1,0,1,1,0,1]
Output: 4
```

**Constraints:**

- `1 <= nums.length <= 10^5`
- `nums[i]` is either `0` or `1`.

**Follow-up:** What if the input numbers come in one by one as an infinite stream? In other words, you can't store all numbers coming from the stream as it's too large to hold in memory. Could you solve it efficiently?

---

## Company Frequency

| Company    | Frequency        | Last Reported |
|------------|------------------|---------------|
| Facebook   | ★★★★☆ High       | 2024          |
| Amazon     | ★★★☆☆ Medium     | 2023          |
| Google     | ★★★☆☆ Medium     | 2023          |
| Microsoft  | ★★☆☆☆ Low        | 2022          |
| Adobe      | ★★☆☆☆ Low        | 2022          |

> ⚠️ Frequency data is crowd-sourced from LeetCode Discuss, Glassdoor, and
> community interview reports. Treat as a signal, not a guarantee.

---

## DSA Concepts Used

- **Sliding Window** — the answer is the widest contiguous window containing at most one `0`; a right-growing / left-shrinking window solves it in one pass and directly answers the streaming follow-up → see [`/dsa/sliding_window.md`](/dsa/sliding_window.md)
- **Array traversal** — the brute force and the previous-count scan are pure left-to-right array passes over a binary array → see [`/dsa/arrays.md`](/dsa/arrays.md)
- **1D Dynamic Programming** — the previous-count method carries two running run-lengths (flip-used / flip-unused) as constant-size DP state that transitions element by element → see [`/dsa/dynamic_programming_1d.md`](/dsa/dynamic_programming_1d.md)

---

## Approaches Overview

| # | Approach | Time | Space | When to use |
|---|----------|------|-------|-------------|
| 1 | Brute Force (flip each zero) | O(n²) | O(1) | Explains the "join left run + right run" idea; too slow for n = 10⁵ |
| 2 | Previous-Count DP | O(n) | O(1) | One pass; carries flipped/unflipped run lengths |
| 3 | Sliding Window (Optimal) | O(n) | O(1) | Best general answer; window with ≤ 1 zero, streaming-friendly |

---

## Approach 1 — Brute Force

### Intuition

You may flip at most one `0`. So the answer is either the longest existing run of `1`s (flip nothing), or — for some specific zero — the length you get by turning **that** zero into a `1` and gluing the run of `1`s on its left to the run of `1`s on its right. Try every zero as the flip target, measure `left + 1 + right`, and keep the maximum (also comparing against the best flip-free run so all-ones inputs work).

### Algorithm

1. First pass: find the longest run of `1`s with no flip (handles the all-ones case).
2. For each index `i` where `nums[i] == 0`:
   - count consecutive `1`s extending **left** from `i` (stop at the first `0`).
   - count consecutive `1`s extending **right** from `i` (stop at the first `0`).
   - candidate `= left + 1 + right` (the `+1` is the flipped zero itself).
3. Return the maximum over all candidates and the flip-free run.

### Complexity

- **Time:** O(n²) — for each zero we scan left and right, up to O(n) work per index.
- **Space:** O(1) — a handful of counters.

### Code

```go
func bruteForce(nums []int) int {
	n := len(nums)
	best := 0
	// First, the longest run with no flip (covers arrays of all 1s).
	run := 0
	for _, v := range nums {
		if v == 1 {
			run++ // extend the current run of ones
			if run > best {
				best = run
			}
		} else {
			run = 0 // a zero ends the flip-free run
		}
	}
	// Now try spending the single flip on each zero.
	for i := 0; i < n; i++ {
		if nums[i] != 0 {
			continue // only zeros are worth flipping
		}
		left := 0
		for j := i - 1; j >= 0 && nums[j] == 1; j-- {
			left++ // ones immediately to the left of i
		}
		right := 0
		for j := i + 1; j < n && nums[j] == 1; j++ {
			right++ // ones immediately to the right of i
		}
		if cand := left + 1 + right; cand > best { // +1 for the flipped zero itself
			best = cand
		}
	}
	return best
}
```

### Dry Run

Example 1: `nums = [1, 0, 1, 1, 0]`. Flip-free longest run = 2 (indices 2..3).

| Zero at i | left run | right run | candidate = left+1+right | best |
|-----------|----------|-----------|--------------------------|------|
| i = 1 | ones at index 0 → 1 | ones at indices 2,3 → 2 | 1 + 1 + 2 = 4 | 4 |
| i = 4 | ones at indices 2,3 → 2 | none (end) → 0 | 2 + 1 + 0 = 3 | 4 |

Maximum = **4** ✔.

---

## Approach 2 — Previous-Count DP

### Intuition

Scan left to right carrying two running lengths of `1`-runs **ending at the current index**:

- `cur` = length using **no** flip,
- `prev` = length using **exactly one** flip.

On a `1`, both runs extend. On a `0`, the no-flip run resets to `0`, and the one-flip run becomes `cur + 1` — we spend the flip on *this* zero, gluing onto the flip-free run that just ended plus the zero itself. Since `prev ≥ cur` always, the answer is the largest `prev` observed.

### Algorithm

1. `cur = prev = 0`, `best = 0`.
2. For each value `v`:
   - if `v == 1`: `cur++`, `prev++`.
   - else: `prev = cur + 1`, `cur = 0`.
   - `best = max(best, prev)`.
3. Return `best`.

### Complexity

- **Time:** O(n) — a single pass.
- **Space:** O(1) — two counters plus the best-so-far.

### Code

```go
func prevCountDP(nums []int) int {
	cur := 0  // run of 1s ending here with NO flip used
	prev := 0 // run of 1s ending here with exactly ONE flip used
	best := 0
	for _, v := range nums {
		if v == 1 {
			cur++  // extend the flip-free run
			prev++ // the flipped run also grows over a real 1
		} else {
			prev = cur + 1 // flip THIS zero: glue onto the flip-free run + itself
			cur = 0        // the flip-free run is broken by a genuine 0
		}
		if prev > best {
			best = prev // prev always ≥ cur, so it alone bounds the answer
		}
	}
	return best
}
```

### Dry Run

Example 1: `nums = [1, 0, 1, 1, 0]`.

| i | v | branch | cur after | prev after | best |
|---|---|--------|-----------|------------|------|
| 0 | 1 | ones | 1 | 1 | 1 |
| 1 | 0 | zero → prev = cur+1 = 2 | 0 | 2 | 2 |
| 2 | 1 | ones | 1 | 3 | 3 |
| 3 | 1 | ones | 2 | 4 | 4 |
| 4 | 0 | zero → prev = cur+1 = 3 | 0 | 3 | 4 |

Answer = **4** ✔.

---

## Approach 3 — Sliding Window (Optimal)

### Intuition

A valid answer is exactly a contiguous window you can turn into all-`1`s by flipping at most one `0` — i.e. a window that contains **at most one zero**. Grow the window to the right; the moment it would hold a second zero, shrink from the left just past the older zero. The widest window ever seen is the answer. Crucially, the left edge only moves forward and we never re-examine old elements, so the algorithm works on a stream: keep only `left`, the current index, and (implicitly, via the counter) the last zero — answering the follow-up.

### Algorithm

1. `left = 0`, `zeros = 0`, `best = 0`.
2. For `right` from `0` to `n − 1`:
   - if `nums[right] == 0`, increment `zeros`.
   - while `zeros > 1`: if `nums[left] == 0` decrement `zeros`; then `left++`.
   - `best = max(best, right − left + 1)`.
3. Return `best`.

### Complexity

- **Time:** O(n) — `right` and `left` each advance at most `n` times total (amortised O(1) per step).
- **Space:** O(1) — two indices and a zero counter; no need to retain the array (streaming-friendly).

### Code

```go
func slidingWindow(nums []int) int {
	left := 0  // left edge of the current window
	zeros := 0 // number of zeros currently inside the window
	best := 0
	for right := 0; right < len(nums); right++ {
		if nums[right] == 0 {
			zeros++ // the new element on the right is a zero
		}
		// If we now hold two zeros, advance left until only one remains.
		for zeros > 1 {
			if nums[left] == 0 {
				zeros-- // the zero leaving on the left frees our single flip
			}
			left++ // shrink the window from the left
		}
		if w := right - left + 1; w > best {
			best = w // widest window with ≤ 1 zero seen so far
		}
	}
	return best
}
```

### Dry Run

Example 1: `nums = [1, 0, 1, 1, 0]`.

| right | nums[right] | zeros | shrink? (zeros > 1) | left | window width | best |
|-------|-------------|-------|---------------------|------|--------------|------|
| 0 | 1 | 0 | no | 0 | 1 | 1 |
| 1 | 0 | 1 | no | 0 | 2 | 2 |
| 2 | 1 | 1 | no | 0 | 3 | 3 |
| 3 | 1 | 1 | no | 0 | 4 | 4 |
| 4 | 0 | 2 | yes → drop nums[0]=1 (left→1), nums[1]=0 (zeros→1, left→2) | 2 | 3 | 4 |

Answer = **4** ✔ (window `[1,1,0]` at indices 2..4 has one zero, width 3; the max width was 4).

---

## Key Takeaways

- **"At most one flip / one deletion" → window with a bounded bad-element count.** Max Consecutive Ones II is the special case "≤ 1 zero in the window"; the general version (LeetCode #1004) allows `k` flips by shrinking whenever `zeros > k`.
- **The sliding window auto-answers the streaming follow-up.** Because `left` is monotone and we only ever look at `nums[right]` and `nums[left]`, you can process an infinite stream keeping O(1) state (the last zero's position suffices).
- **Two-state run-length DP is an equally clean O(n) alternative** when you'd rather carry "run ending here with/without the flip" than manage window edges.
- **Brute force reveals the structure**: any single-flip answer is `left_run + 1 + right_run` around a chosen zero — a useful sanity check for the linear solutions.

---

## Related Problems

- LeetCode #485 — Max Consecutive Ones (no flip; the base problem)
- LeetCode #1004 — Max Consecutive Ones III (flip at most `k` zeros; same window, `zeros > k`)
- LeetCode #424 — Longest Repeating Character Replacement (window with ≤ k replacements)
- LeetCode #1493 — Longest Subarray of 1's After Deleting One Element (must delete exactly one)
- LeetCode #340 — Longest Substring with At Most K Distinct Characters (bounded-content window)
