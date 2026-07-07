# 0467 — Unique Substrings in Wraparound String

> LeetCode #467 · Difficulty: Medium
> **Categories:** String, Dynamic Programming

---

## Problem Statement

We define the string `base` to be the infinite wraparound string of `"abcdefghijklmnopqrstuvwxyz"`, so `base` will look like this:

- `"...zabcdefghijklmnopqrstuvwxyzabcdefghijklmnopqrstuvwxyzabcd..."`.

Given a string `s`, return *the number of **unique non-empty substrings** of* `s` *that are present in* `base`.

**Example 1:**

```
Input: s = "a"
Output: 1
Explanation: Only the substring "a" of s is in base.
```

**Example 2:**

```
Input: s = "cac"
Output: 2
Explanation: There are two substrings ("a", "c") of s in base.
```

**Example 3:**

```
Input: s = "zab"
Output: 6
Explanation: There are six substrings ("z", "a", "b", "za", "ab", and "zab") of s in base.
```

**Constraints:**

- `1 <= s.length <= 10^5`
- `s` consists of lowercase English letters.

---

## Company Frequency

| Company    | Frequency        | Last Reported |
|------------|------------------|---------------|
| Amazon     | ★★★☆☆ Medium     | 2023          |
| Google     | ★★☆☆☆ Low        | 2022          |
| Microsoft  | ★★☆☆☆ Low        | 2022          |
| Bloomberg  | ★★☆☆☆ Low        | 2021          |

> ⚠️ Frequency data is crowd-sourced from LeetCode Discuss, Glassdoor, and
> community interview reports. Treat as a signal, not a guarantee.

---

## DSA Concepts Used

- **1-D Dynamic Programming (state = ending character)** — the count of distinct valid substrings ending at each letter is the longest wraparound run ending there; keeping one number per letter is a 26-state DP → see [`/dsa/dynamic_programming_1d.md`](/dsa/dynamic_programming_1d.md)
- **String structure / consecutive-run detection** — validity of a substring reduces to "every adjacent pair steps forward by one letter cyclically", a contiguous-run property of the character sequence → see [`/dsa/string_algorithms.md`](/dsa/string_algorithms.md)

---

## Approaches Overview

| # | Approach | Time | Space | When to use |
|---|----------|------|-------|-------------|
| 1 | Brute Force (Enumerate + Set) | O(n²) substrings, up to O(n³) with hashing | O(n²) | Understanding / tiny inputs; TLE + MLE for n up to 10⁵ |
| 2 | DP on Longest Run Ending at Each Letter (Optimal) | O(n) | O(1) | Always; single pass, 26 counters, no set |

---

## Approach 1 — Brute Force (Enumerate + Set)

### Intuition

`base` is the alphabet repeated forever, so a substring occurs in `base` **iff** each adjacent pair `(prev, cur)` steps forward by exactly one letter cyclically: `(cur − prev + 26) % 26 == 1`, with `z → a` wrapping. Enumerate all O(n²) substrings; while extending a substring one character at a time, stop the instant the wraparound step breaks (no longer substring from that start can be valid). Insert each valid substring into a set so duplicates — like the two `"c"`s in `"cac"` — are counted once.

### Algorithm

1. For each start index `i`, walk an end index `j` from `i` forward.
2. When `j > i`, check the step between `s[j−1]` and `s[j]`; if it is not `+1` cyclically, `break` (this start is exhausted).
3. Insert the still-valid `s[i..j]` into a set `seen`.
4. Return `len(seen)`.

### Complexity

- **Time:** O(n²) substrings; hashing each substring of length up to `n` makes it up to O(n³) for a fully increasing string. Only feasible for small `s`.
- **Space:** O(number of distinct valid substrings), which can be O(n²).

### Code

```go
func bruteForce(s string) int {
	seen := make(map[string]struct{}) // distinct valid substrings
	n := len(s)
	for i := 0; i < n; i++ {
		// Grow the window [i..j]; stop as soon as the wraparound chain breaks.
		for j := i; j < n; j++ {
			if j > i {
				// step from previous char to current char, cyclically
				step := (int(s[j]) - int(s[j-1]) + 26) % 26
				if step != 1 { // not a +1 wraparound move → chain broken
					break
				}
			}
			seen[s[i:j+1]] = struct{}{} // record this contiguous wraparound run
		}
	}
	return len(seen) // each distinct substring counted exactly once
}
```

### Dry Run

Example 3: `s = "zab"`.

| i | j | pair checked | step | valid? | substring inserted | set size |
|---|---|--------------|------|--------|--------------------|----------|
| 0 | 0 | — | — | yes | "z" | 1 |
| 0 | 1 | z→a | (0−25+26)%26 = 1 | yes | "za" | 2 |
| 0 | 2 | a→b | (1−0+26)%26 = 1 | yes | "zab" | 3 |
| 1 | 1 | — | — | yes | "a" | 4 |
| 1 | 2 | a→b | 1 | yes | "ab" | 5 |
| 2 | 2 | — | — | yes | "b" | 6 |

Set = { z, za, zab, a, ab, b }. Answer `= 6` ✔

---

## Approach 2 — DP on Longest Run Ending at Each Letter (Optimal)

### Intuition

Bucket every valid substring by its **last character**. Claim: the number of *distinct* valid substrings ending at letter `c` equals the **length of the longest valid wraparound run ending at `c`**. Why: any valid substring ending at `c` is a suffix of the longest run ending at `c`; if that longest run has length `L`, its suffixes have lengths `1..L` — exactly `L` distinct strings ending in `c`. Two different runs that both end at `c` share their shorter suffixes, so keeping only the **maximum** run length per ending letter counts each distinct substring exactly once. Summing `maxEnd[c]` over all 26 letters is the answer — no set required.

### Algorithm

1. `maxEnd[26]` — longest valid run ending at each letter, all 0.
2. Scan `s`, maintaining `curLen`, the run length ending at `s[i]`:
   - if `i > 0` and `(s[i] − s[i−1] + 26) % 26 == 1`, then `curLen++`;
   - else `curLen = 1` (run restarts).
3. `maxEnd[s[i]−'a'] = max(maxEnd[s[i]−'a'], curLen)`.
4. Return the sum of `maxEnd[0..25]`.

### Complexity

- **Time:** O(n) — one linear pass plus a constant 26-element sum.
- **Space:** O(1) — a fixed array of 26 integers.

### Code

```go
func dpMaxEndingAt(s string) int {
	var maxEnd [26]int // maxEnd[c] = longest wraparound run ending at letter c
	curLen := 0        // length of the run ending at the current position
	for i := 0; i < len(s); i++ {
		if i > 0 && (int(s[i])-int(s[i-1])+26)%26 == 1 {
			curLen++ // s[i] extends the wraparound run from s[i-1]
		} else {
			curLen = 1 // run restarts at s[i] (or first char)
		}
		c := s[i] - 'a'          // 0..25 index of the ending letter
		if curLen > maxEnd[c] {  // keep the longest run ever ending at c
			maxEnd[c] = curLen
		}
	}
	total := 0
	for _, v := range maxEnd { // sum the per-letter maxima
		total += v // v distinct substrings end at that letter
	}
	return total
}
```

### Dry Run

Example 3: `s = "zab"`.

| i | s[i] | continues run? (step==1) | curLen | letter | maxEnd update |
|---|------|--------------------------|--------|--------|----------------|
| 0 | z | n/a (first) | 1 | z | maxEnd[z] = 1 |
| 1 | a | z→a: (0−25+26)%26 = 1 ✓ | 2 | a | maxEnd[a] = 2 |
| 2 | b | a→b: 1 ✓ | 3 | b | maxEnd[b] = 3 |

Sum `= maxEnd[z] + maxEnd[a] + maxEnd[b] = 1 + 2 + 3 = 6` ✔ — the three letters account for substrings ending in `z` (`z`), ending in `a` (`a`, `za`), and ending in `b` (`b`, `ab`, `zab`).

---

## Key Takeaways

- **Count distinct substrings by bucketing on the last character.** When "valid" substrings form contiguous runs, the distinct ones ending at a letter are exactly the suffixes of the *longest* such run — so `sum of longest-run-per-ending-letter` counts every distinct substring once, replacing an O(n²) set with 26 integers.
- **Cyclic "next letter" test:** `(cur − prev + 26) % 26 == 1` handles the `z → a` wrap without a special case.
- **The +1 running-length recurrence** (`curLen++` on continuation, else reset to 1) is the same shape as *Longest Consecutive/Increasing Run*; here we additionally take a max per ending symbol to deduplicate.
- Deduplication via a max-per-key array (instead of a hash set) is a recurring trick whenever the objects form nested/suffix families.

---

## Related Problems

- LeetCode #300 — Longest Increasing Subsequence (running-length DP relatives)
- LeetCode #674 — Longest Continuous Increasing Subsequence (the reset-to-1 run recurrence)
- LeetCode #128 — Longest Consecutive Sequence (consecutive-run reasoning)
- LeetCode #940 — Distinct Subsequences II (counting distinct substrings/subsequences)
- LeetCode #1638 — Count Substrings That Differ by One Character (bucketed substring counting)
