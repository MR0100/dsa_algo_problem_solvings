# 0424 — Longest Repeating Character Replacement

> LeetCode #424 · Difficulty: Medium
> **Categories:** Hash Table, String, Sliding Window

---

## Problem Statement

You are given a string `s` and an integer `k`. You can choose any character of the string and change it to any other uppercase English character. You can perform this operation at most `k` times.

Return *the length of the longest substring containing the same letter you can get after performing the above operations*.

**Example 1:**

```
Input: s = "ABAB", k = 2
Output: 4
Explanation: Replace the two 'A's with two 'B's or vice versa.
```

**Example 2:**

```
Input: s = "AABABBA", k = 1
Output: 4
Explanation: Replace the one 'A' in the middle with 'B' and form "AABBBBA".
The substring "BBBB" has the longest repeating letters, which is 4.
There may exists other ways to achieve this answer too.
```

**Constraints:**

- `1 <= s.length <= 10^5`
- `s` consists of only uppercase English letters.
- `0 <= k <= s.length`

---

## Company Frequency

| Company    | Frequency        | Last Reported |
|------------|------------------|---------------|
| Amazon     | ★★★★☆ High       | 2024          |
| Google     | ★★★★☆ High       | 2024          |
| Meta       | ★★★☆☆ Medium     | 2023          |
| Microsoft  | ★★★☆☆ Medium     | 2023          |
| Bloomberg  | ★★☆☆☆ Low        | 2022          |
| TikTok     | ★★☆☆☆ Low        | 2023          |

> ⚠️ Frequency data is crowd-sourced from LeetCode Discuss, Glassdoor, and
> community interview reports. Treat as a signal, not a guarantee.

---

## DSA Concepts Used

- **Sliding Window (variable size)** — the answer is the widest window where the non-majority characters number at most `k`; a right pointer grows the window and a left pointer maintains the invariant → see [`/dsa/sliding_window.md`](/dsa/sliding_window.md)
- **Frequency counting** — a 26-entry array tracks each letter's count in the current window; `maxFreq` (the count of the majority letter) drives the validity test → see [`/dsa/hash_map.md`](/dsa/hash_map.md)

---

## Approaches Overview

| # | Approach | Time | Space | When to use |
|---|----------|------|-------|-------------|
| 1 | Brute Force | O(n²) | O(1) | Baseline / oracle; `n = 10⁵` gives 10¹⁰ substrings — TLE |
| 2 | Sliding Window with Recount | O(26·n) = O(n) | O(1) | Clear, easy-to-justify linear solution |
| 3 | Sliding Window, Non-Decreasing maxFreq (Optimal) | O(n) | O(1) | The slick one-pass answer; no per-step 26-scan |

---

## Approach 1 — Brute Force

### Intuition

A substring is convertible to all-same-letter within `k` operations exactly when the count of characters that are **not** its most frequent letter is `≤ k`, i.e. `length − maxFrequency ≤ k` (those are the ones you must overwrite). Enumerate every substring, compute its majority-letter count incrementally, test the condition, and keep the longest that passes.

### Algorithm

1. For each start `i`, extend end `j`, maintaining `count[26]` for `s[i..j]`.
2. Track `maxFreq`, the largest count in the window (update as each `s[j]` is added).
3. If `(j − i + 1) − maxFreq ≤ k`, the window is convertible — update `best`.
4. Return `best`.

### Complexity

- **Time:** O(n²) — there are ~n²/2 substrings; each extension is O(1) plus an O(1) majority update.
- **Space:** O(26) = O(1) — one fixed-size letter-count array per start.

### Code

```go
func bruteForce(s string, k int) int {
	best := 0
	for i := 0; i < len(s); i++ {
		var count [26]int // letter frequencies for windows starting at i
		maxFreq := 0       // most frequent letter's count in s[i..j]
		for j := i; j < len(s); j++ {
			count[s[j]-'A']++            // extend window to include s[j]
			if count[s[j]-'A'] > maxFreq {
				maxFreq = count[s[j]-'A'] // s[j] may be the new most frequent letter
			}
			windowLen := j - i + 1 // current substring length
			// Chars other than the majority must be replaced; need at most k.
			if windowLen-maxFreq <= k && windowLen > best {
				best = windowLen // this substring is convertible and longer
			}
		}
	}
	return best
}
```

### Dry Run

Input `s = "AABABBA"`, `k = 1`. Showing the start `i = 0` row (the one that yields the answer). Window is `s[0..j]`.

| j | char | count (A,B) | maxFreq | len | len−maxFreq | ≤ k? | best |
|---|------|-------------|---------|-----|-------------|------|------|
| 0 | A    | (1,0)       | 1       | 1   | 0           | yes  | 1    |
| 1 | A    | (2,0)       | 2       | 2   | 0           | yes  | 2    |
| 2 | B    | (2,1)       | 2       | 3   | 1           | yes  | 3    |
| 3 | A    | (3,1)       | 3       | 4   | 1           | yes  | **4** |
| 4 | B    | (3,2)       | 3       | 5   | 2           | no   | 4    |
| 5 | B    | (3,3)       | 3       | 6   | 3           | no   | 4    |
| 6 | A    | (4,3)       | 4       | 7   | 3           | no   | 4    |

Other starts never beat 4, so the answer is `4`.

---

## Approach 2 — Sliding Window with Recount

### Intuition

Instead of restarting for every substring, keep one window `[left, right]`. It is valid while `len − maxFreq ≤ k`. Move `right` to grow it; whenever adding `s[right]` breaks validity, move `left` forward until it holds again. The widest width ever seen is the answer. Here we recompute `maxFreq` by scanning the 26 counts each time — transparent, and still linear overall.

### Algorithm

1. `left = 0`; for each `right`, do `count[s[right]]++`.
2. Compute `maxFreq = max` over the 26 counts.
3. While `(right − left + 1) − maxFreq > k`: `count[s[left]]--`, `left++`, recompute `maxFreq`.
4. Track the maximum `right − left + 1`.

### Complexity

- **Time:** O(26·n) = O(n) — each index enters and leaves the window once; each validity check scans 26 counters (constant).
- **Space:** O(26) = O(1).

### Code

```go
func slidingWindowRecount(s string, k int) int {
	var count [26]int
	best := 0
	left := 0
	for right := 0; right < len(s); right++ {
		count[s[right]-'A']++ // include the new right-hand character

		// Recompute the window's most frequent letter count.
		maxFreq := 0
		for _, c := range count {
			if c > maxFreq {
				maxFreq = c
			}
		}

		// Shrink while too many chars would need replacing.
		for (right-left+1)-maxFreq > k {
			count[s[left]-'A']-- // drop the leftmost character
			left++
			maxFreq = 0 // window changed → recompute the majority count
			for _, c := range count {
				if c > maxFreq {
					maxFreq = c
				}
			}
		}

		if right-left+1 > best {
			best = right - left + 1 // widest valid window so far
		}
	}
	return best
}
```

### Dry Run

Input `s = "AABABBA"`, `k = 1`. `w = right−left+1`.

| right | char | window | count (A,B) | maxFreq | w−maxFreq>k? shrink | best |
|-------|------|--------|-------------|---------|---------------------|------|
| 0 | A | `A`     | (1,0) | 1 | 0>1? no  | 1 |
| 1 | A | `AA`    | (2,0) | 2 | 0>1? no  | 2 |
| 2 | B | `AAB`   | (2,1) | 2 | 1>1? no  | 3 |
| 3 | A | `AABA`  | (3,1) | 3 | 1>1? no  | **4** |
| 4 | B | `AABAB` | (3,2) | 3 | 2>1? yes → drop `A`, left=1 (`ABAB`, w=4, maxFreq 2, 2>1 yes) → drop `A`, left=2 (`BAB`, w=3, maxFreq 2, 1>1 no) | 4 |
| 5 | B | `BABB` | (1,3) | 3 | 1>1? no  | 4 |
| 6 | A | `BABBA` | (2,3) | 3 | 2>1? yes → drop `B`, left=3 (`ABBA`, maxFreq 2, 2>1 yes) → drop `A`, left=4 (`BBA`, maxFreq 2, 1>1 no) | 4 |

Answer `4`.

---

## Approach 3 — Sliding Window, Non-Decreasing maxFreq (Optimal)

### Intuition

We only care about the **longest** valid window, so let `maxFreq` be a high-water mark that never decreases. When the window is valid (`width − maxFreq ≤ k`), extend it. When it isn't, advance `left` by exactly one so the width does not grow — but do **not** lower `maxFreq`. This is safe: a shorter window whose majority count is smaller than a value we already recorded can never produce a longer answer, so we don't need an accurate `maxFreq` once it has peaked. The width `right − left + 1` therefore never shrinks and equals the answer at the end.

### Algorithm

1. `left = 0`, `maxFreq = 0`.
2. For each `right`: `count[s[right]]++`; `maxFreq = max(maxFreq, count[s[right]])`.
3. If `(right − left + 1) − maxFreq > k`: `count[s[left]]--`, `left++` (slide, keeping `maxFreq`).
4. Return the final width `len(s) − left`.

### Complexity

- **Time:** O(n) — one pass, O(1) work per character (no 26-scan).
- **Space:** O(26) = O(1).

### Code

```go
func slidingWindowOptimal(s string, k int) int {
	var count [26]int
	maxFreq := 0 // high-water mark of any letter's count in the window; never decreases
	left := 0
	for right := 0; right < len(s); right++ {
		count[s[right]-'A']++ // add s[right] to the window
		if count[s[right]-'A'] > maxFreq {
			maxFreq = count[s[right]-'A'] // update the running majority count
		}

		// If more than k chars would need replacing, slide the window right by
		// one (keep its size fixed). We do not decrease maxFreq — a shorter
		// window can never improve on a length already achieved.
		if (right-left+1)-maxFreq > k {
			count[s[left]-'A']-- // remove the leftmost character
			left++
		}
	}
	// The window never shrank, so its final width is the longest valid length.
	return len(s) - left
}
```

### Dry Run

Input `s = "AABABBA"`, `k = 1`. `w = right−left+1`; note `left` only ever moves by at most one per step, and `maxFreq` never falls.

| right | char | count (A,B) | maxFreq | w before | w−maxFreq>k? | left after | window width |
|-------|------|-------------|---------|----------|--------------|------------|--------------|
| 0 | A | (1,0) | 1 | 1 | 0>1? no  | 0 | 1 |
| 1 | A | (2,0) | 2 | 2 | 0>1? no  | 0 | 2 |
| 2 | B | (2,1) | 2 | 3 | 1>1? no  | 0 | 3 |
| 3 | A | (3,1) | 3 | 4 | 1>1? no  | 0 | 4 |
| 4 | B | (3,2) | 3 | 5 | 2>1? yes → drop `s[0]=A` | 1 | 4 |
| 5 | B | (2,3) | 3 | 4 | 1>1? no  | 1 | 4 |
| 6 | A | (3,3) | 3 | 5 | 2>1? yes → drop `s[1]=A` | 2 | 4 |

Final `len(s) − left = 7 − 2 = 4`. Answer `4`. Notice `maxFreq` stayed at 3 even though no window from `left` onward literally contained 3 of a letter — that's the deliberate high-water-mark trick, and it never overstates the result.

---

## Key Takeaways

- **Window validity test:** `windowLength − maxFrequency ≤ k`. The characters you must replace are exactly the non-majority ones; this single inequality is the whole problem.
- **The high-water-mark trick.** For "longest window" problems you often don't need an exact running maximum — only a value that never underestimates the best-so-far. Letting `maxFreq` only increase drops an O(alphabet) factor and gives a clean single `if` (no inner shrink loop).
- **Never-shrinking window.** Because the optimal window only ever slides (never contracts), its final width *is* the answer — no separate `best` variable needed.
- Recognise this template: it recurs in "longest substring with at most k of something" problems.

---

## Related Problems

- LeetCode #340 — Longest Substring with At Most K Distinct Characters (same window family)
- LeetCode #1004 — Max Consecutive Ones III (flip at most k zeros — identical shape)
- LeetCode #3 — Longest Substring Without Repeating Characters (variable window)
- LeetCode #1208 — Get Equal Substrings Within Budget (window with a cost budget)
- LeetCode #76 — Minimum Window Substring (shrinking-window counterpart)
