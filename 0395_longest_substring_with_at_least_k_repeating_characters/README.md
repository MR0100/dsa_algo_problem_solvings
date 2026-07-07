# 0395 — Longest Substring with At Least K Repeating Characters

> LeetCode #395 · Difficulty: Medium
> **Categories:** Hash Table, String, Divide and Conquer, Sliding Window

---

## Problem Statement

Given a string `s` and an integer `k`, return the length of the longest substring of `s`
such that the frequency of each character in this substring is greater than or equal to
`k`.

If no such substring exists, return `0`.

**Example 1:**

```
Input: s = "aaabb", k = 3
Output: 3
Explanation: The longest substring is "aaa", as 'a' is repeated 3 times.
```

**Example 2:**

```
Input: s = "ababbc", k = 2
Output: 5
Explanation: The longest substring is "ababb", as 'a' is repeated 2 times and 'b' is
repeated 3 times.
```

**Constraints:**

- `1 <= s.length <= 10^4`
- `s` consists of only lowercase English letters.
- `1 <= k <= 10^5`

---

## Company Frequency

| Company   | Frequency         | Last Reported |
|-----------|-------------------|---------------|
| Google    | ★★★★☆ High        | 2024          |
| Amazon    | ★★★★☆ High        | 2024          |
| Microsoft | ★★★☆☆ Medium      | 2023          |
| Bloomberg | ★★★☆☆ Medium      | 2023          |

> ⚠️ Frequency data is crowd-sourced from LeetCode Discuss, Glassdoor, and
> community interview reports. Treat as a signal, not a guarantee.

---

## DSA Concepts Used

- **Hash / frequency counting** — 26-letter tallies to test the "each ≥ k" condition → see [`/dsa/hash_map.md`](/dsa/hash_map.md)
- **Divide and Conquer** — split on characters that can never qualify → see [`/dsa/divide_and_conquer.md`](/dsa/divide_and_conquer.md)
- **Sliding Window** — make the condition monotone by fixing the distinct-char count → see [`/dsa/sliding_window.md`](/dsa/sliding_window.md)

---

## Approaches Overview

| # | Approach | Time | Space | When to use |
|---|----------|------|-------|-------------|
| 1 | Brute Force | O(n²·26) | O(1) | Baseline / correctness reference |
| 2 | Divide and Conquer | O(n·26) typ. | O(26·depth) | Elegant; the classic "splitter" insight |
| 3 | Sliding Window (fixed unique) | O(26·n) | O(1) | Optimal linear-ish; makes window monotone |

---

## Approach 1 — Brute Force

### Intuition

Encode the definition directly: for every substring, keep a frequency table and accept it
if every present character occurs at least `k` times; track the longest accepted length.

### Algorithm

1. For each start `i`, reset a 26-frequency array.
2. For each end `j ≥ i`: increment `count[s[j]]`, then check all present counts ≥ k;
   update the best length if valid.
3. Return the best.

### Complexity

- **Time:** O(n²·26) — O(n²) substrings × O(26) validity check.
- **Space:** O(1) — fixed 26-length array.

### Code

```go
func bruteForce(s string, k int) int {
	n := len(s)
	best := 0
	for i := 0; i < n; i++ {
		var count [26]int
		for j := i; j < n; j++ {
			count[s[j]-'a']++
			if allAtLeastK(count, k) && j-i+1 > best {
				best = j - i + 1
			}
		}
	}
	return best
}
```

### Dry Run

`s = "aaabb"`, `k = 3`. Windows starting at `i = 0`:

| j | window | counts | valid (all ≥3)? | best |
|---|--------|--------|-----------------|------|
| 0 | "a" | a:1 | no | 0 |
| 1 | "aa" | a:2 | no | 0 |
| 2 | "aaa" | a:3 | yes | 3 |
| 3 | "aaab" | a:3,b:1 | no | 3 |
| 4 | "aaabb" | a:3,b:2 | no | 3 |

Later starts yield nothing longer ⇒ **`3`**.

---

## Approach 2 — Divide and Conquer

### Intuition

Any character that appears fewer than `k` times in a segment can **never** belong to a
valid substring — so it acts as a splitter. If a segment has no such rare character, the
whole segment is valid. Otherwise pick a rare character, split the segment on every
occurrence of it, and recurse into each piece.

### Algorithm

1. Count frequencies in `s[lo:hi]`.
2. If no character is present-but-rare (all ≥ k), return `hi-lo`.
3. Else pick such a splitter; recurse on each maximal piece between its occurrences; return
   the max over pieces.

### Complexity

- **Time:** O(n·26) typical, O(n²) worst case — each recursion level eliminates at least
  one character class (≤ 26 levels), each level scans O(n).
- **Space:** O(26·depth) — recursion depth ≤ 26 plus per-frame count arrays.

### Code

```go
func divideAndConquer(s string, k int) int {
	var solve func(lo, hi int) int
	solve = func(lo, hi int) int {
		if hi-lo < k {
			return 0
		}
		var count [26]int
		for i := lo; i < hi; i++ {
			count[s[i]-'a']++
		}
		splitter := byte(0)
		found := false
		for c := 0; c < 26; c++ {
			if count[c] > 0 && count[c] < k {
				splitter = byte('a' + c)
				found = true
				break
			}
		}
		if !found {
			return hi - lo
		}
		best := 0
		start := lo
		for i := lo; i < hi; i++ {
			if s[i] == splitter {
				if r := solve(start, i); r > best {
					best = r
				}
				start = i + 1
			}
		}
		if r := solve(start, hi); r > best {
			best = r
		}
		return best
	}
	return solve(0, len(s))
}
```

### Dry Run

`s = "aaabb"`, `k = 3`. `solve(0,5)`: counts a:3, b:2. `b` is present-but-rare → splitter
`b`.

| piece | range | recurse result |
|-------|-------|----------------|
| "aaa" | [0,3) | counts a:3 → no rare char → return 3 |
| "" between the two b's | [4,4) | length 0 < k → 0 |
| "" after last b | [5,5) | 0 |

Max = **`3`**.

---

## Approach 3 — Sliding Window (Fixed Unique Count)

### Intuition

A plain sliding window fails because "valid" is not monotone as the window grows. Fix the
number of **distinct** characters allowed, `unique = 1..26`. For each fixed target, run a
window keeping exactly `unique` distinct characters and count how many of them already have
frequency ≥ k. When `distinct == unique` **and** all of them are ≥ k, the window is valid.
Fixing the distinct count restores monotonicity so the two-pointer window works.

### Algorithm

1. For `unique` from 1 to 26:
   - Expand `right`; on adding a char update `distinct` and the `atLeastK` tally.
   - While `distinct > unique`, shrink from `left`, updating the tallies.
   - When `distinct == unique` and `atLeastK == unique`, update the answer.
2. Return the best length.

### Complexity

- **Time:** O(26·n) — 26 linear passes.
- **Space:** O(1) — a 26-length frequency array per pass.

### Code

```go
func slidingWindow(s string, k int) int {
	n := len(s)
	best := 0
	for unique := 1; unique <= 26; unique++ {
		var count [26]int
		distinct := 0
		atLeastK := 0
		left := 0
		for right := 0; right < n; right++ {
			ri := s[right] - 'a'
			if count[ri] == 0 {
				distinct++
			}
			count[ri]++
			if count[ri] == k {
				atLeastK++
			}
			for distinct > unique {
				li := s[left] - 'a'
				if count[li] == k {
					atLeastK--
				}
				count[li]--
				if count[li] == 0 {
					distinct--
				}
				left++
			}
			if distinct == unique && atLeastK == unique {
				if right-left+1 > best {
					best = right - left + 1
				}
			}
		}
	}
	return best
}
```

### Dry Run

`s = "aaabb"`, `k = 3`, focus on `unique = 1`:

| right | char | count | distinct | atLeastK | window valid? | best |
|-------|------|-------|----------|----------|---------------|------|
| 0 | a | a:1 | 1 | 0 | 1==1 but atLeastK 0 → no | 0 |
| 1 | a | a:2 | 1 | 0 | no | 0 |
| 2 | a | a:3 | 1 | 1 | 1==1 & 1==1 → yes, len 3 | 3 |
| 3 | b | a:3,b:1 | 2>1 → shrink a's out until distinct=1 (window "b") | 1 | 0 | no | 3 |
| 4 | b | b:2 | 1 | 0 | no | 3 |

Other `unique` values never beat 3 ⇒ **`3`**.

---

## Key Takeaways

- The "each char ≥ k" condition is **not monotone**, so a naive one-pass sliding window
  does not work — a crucial trap.
- **Divide and conquer on a splitter character** (one that can never appear in the answer)
  is the signature trick; recognize it whenever a "forbidden element" partitions the space.
- **Fix the distinct-char count (1..26)** to restore monotonicity and reuse the standard
  sliding-window template — a powerful "add a bounded outer loop to make the inner problem
  monotone" pattern.
- Maintain an `atLeastK` counter incrementally so window validity is O(1) to test.

---

## Related Problems

- LeetCode #340 — Longest Substring with At Most K Distinct Characters (fixed-distinct window)
- LeetCode #003 — Longest Substring Without Repeating Characters (classic sliding window)
- LeetCode #159 — Longest Substring with At Most Two Distinct Characters
- LeetCode #076 — Minimum Window Substring (window with a coverage condition)
