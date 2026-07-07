# 0438 — Find All Anagrams in a String

> LeetCode #438 · Difficulty: Medium
> **Categories:** String, Sliding Window, Hash Table

---

## Problem Statement

Given two strings `s` and `p`, return *an array of all the start indices of* `p`*'s anagrams in* `s`. You may return the answer in **any order**.

**Example 1:**

```
Input: s = "cbaebabacd", p = "abc"
Output: [0,6]
Explanation:
The substring with start index = 0 is "cba", which is an anagram of "abc".
The substring with start index = 6 is "bac", which is an anagram of "abc".
```

**Example 2:**

```
Input: s = "abab", p = "ab"
Output: [0,1,2]
Explanation:
The substring with start index = 0 is "ab", which is an anagram of "ab".
The substring with start index = 1 is "ba", which is an anagram of "ab".
The substring with start index = 2 is "ab", which is an anagram of "ab".
```

**Constraints:**

- `1 <= s.length, p.length <= 3 * 10^4`
- `s` and `p` consist of lowercase English letters.

---

## Company Frequency

| Company    | Frequency        | Last Reported |
|------------|------------------|---------------|
| Amazon     | ★★★★★ Very High  | 2024          |
| Facebook   | ★★★★☆ High       | 2024          |
| Google     | ★★★★☆ High       | 2024          |
| Microsoft  | ★★★☆☆ Medium     | 2023          |
| Bloomberg  | ★★★☆☆ Medium     | 2023          |
| Uber       | ★★☆☆☆ Low        | 2022          |

> ⚠️ Frequency data is crowd-sourced from LeetCode Discuss, Glassdoor, and
> community interview reports. Treat as a signal, not a guarantee.

---

## DSA Concepts Used

- **Sliding Window (fixed size)** — the answer is every length-`|p|` window of `s` that is an anagram of `p`; the window slides one character at a time → see [`/dsa/sliding_window.md`](/dsa/sliding_window.md)
- **Hash Table / Frequency Count** — anagram equivalence is equality of letter counts; a 26-slot array is the compact hash of a lowercase string → see [`/dsa/hash_map.md`](/dsa/hash_map.md)
- **Counting Sort idea** — building a fixed 26-bucket tally per string is the counting-sort histogram, which is what makes the comparison O(26) instead of O(m log m) → see [`/dsa/counting_sort.md`](/dsa/counting_sort.md)

---

## Approaches Overview

| # | Approach | Time | Space | When to use |
|---|----------|------|-------|-------------|
| 1 | Brute Force (Sort Every Window) | O((n−m)·m log m) | O(m) | Baseline; clear but re-sorts each window |
| 2 | Sliding Window + Count Compare | O(n·26) = O(n) | O(1) | Clean and fast; compares two 26-arrays per slide |
| 3 | Sliding Window + Match Counter (Optimal) | O(n) | O(1) | Tightest constant; O(1) per slide, no 26-loop |

Here `n = len(s)`, `m = len(p)`.

---

## Approach 1 — Brute Force (Sort Every Window)

### Intuition

Two strings are anagrams iff their sorted characters are identical. So slide a window of width `m` over `s`, sort the window, and compare to a pre-sorted `p`. It is the most direct reading of "is this window an anagram?", at the price of an O(m log m) sort per window.

### Algorithm

1. Compute `sortedP` = sorted characters of `p`.
2. For each start `i` from `0` to `n − m`: copy `s[i:i+m]`, sort it, compare to `sortedP`.
3. Append `i` to the result whenever they match.

### Complexity

- **Time:** O((n − m) · m log m) — one sort per window, up to `n − m + 1` windows.
- **Space:** O(m) — the per-window buffer.

### Code

```go
func bruteForce(s string, p string) []int {
	n, m := len(s), len(p)
	res := []int{}
	if m > n {
		return res // p can't fit in s → no anagrams
	}
	sortedP := []byte(p)
	sort.Slice(sortedP, func(a, b int) bool { return sortedP[a] < sortedP[b] }) // canonical form of p

	for i := 0; i+m <= n; i++ {
		window := []byte(s[i : i+m])                                             // copy the current window
		sort.Slice(window, func(a, b int) bool { return window[a] < window[b] }) // canonicalise it
		if string(window) == string(sortedP) {                                  // same multiset of letters?
			res = append(res, i)
		}
	}
	return res
}
```

### Dry Run

Example 1: `s = "cbaebabacd", p = "abc"`, `m = 3`, `sortedP = "abc"`.

| start i | window | sorted window | == "abc"? | record |
|---------|--------|---------------|-----------|--------|
| 0 | `cba` | `abc` | yes | **0** |
| 1 | `bae` | `abe` | no | — |
| 2 | `aeb` | `abe` | no | — |
| 3 | `eba` | `abe` | no | — |
| 4 | `bab` | `abb` | no | — |
| 5 | `aba` | `aab` | no | — |
| 6 | `bac` | `abc` | yes | **6** |
| 7 | `acd` | `acd` | no | — |

Result: `[0, 6]` ✔

---

## Approach 2 — Sliding Window with Count Comparison

### Intuition

Anagram equivalence is equality of letter frequencies, so represent each string by a 26-length count array. The window's counts change by only **two entries** when it slides (one letter enters, one leaves), so we never rebuild or re-sort — we just increment the incoming letter, decrement the outgoing one, and compare the window array to `p`'s array. In Go, `[26]int == [26]int` is a full element-wise comparison, giving a tidy O(26) test per position.

### Algorithm

1. Build `need[26]` from `p` and `win[26]` from the first `m` characters of `s`.
2. If `win == need`, record start `0`.
3. For `i` from `m` to `n − 1`: `win[s[i]]++` (enter), `win[s[i-m]]--` (leave); if `win == need`, record start `i − m + 1`.

### Complexity

- **Time:** O(n · 26) = O(n) — a constant 26-length array compare at each of `n` positions.
- **Space:** O(1) — two fixed 26-length arrays.

### Code

```go
func slidingWindowCompare(s string, p string) []int {
	n, m := len(s), len(p)
	res := []int{}
	if m > n {
		return res
	}
	var need, win [26]int // fixed alphabet: 'a'..'z'
	for i := 0; i < m; i++ {
		need[p[i]-'a']++ // target letter frequencies
		win[s[i]-'a']++  // first window's letter frequencies
	}
	if win == need { // array value comparison in Go is a full element-wise check
		res = append(res, 0)
	}
	// Slide the window one character at a time.
	for i := m; i < n; i++ {
		win[s[i]-'a']++   // the new right-hand character enters the window
		win[s[i-m]-'a']-- // the leftmost character leaves the window
		if win == need {
			res = append(res, i-m+1) // window now covers s[i-m+1 .. i]
		}
	}
	return res
}
```

### Dry Run

Example 1: `s = "cbaebabacd", p = "abc"`. `need = {a:1, b:1, c:1}`.

Initial window `s[0:3] = "cba"` → `win = {a:1, b:1, c:1}` → equals `need` → record **0**.

| step i | char in (s[i]) | char out (s[i-m]) | window covers | win == need? | record |
|--------|----------------|-------------------|---------------|--------------|--------|
| 3 | `e` | `c` | `bae` | {a:1,b:1,e:1} ≠ | — |
| 4 | `b` | `b` | `aeb` | {a:1,b:1,e:1} ≠ | — |
| 5 | `a` | `a` | `eba` | {a:1,b:1,e:1} ≠ | — |
| 6 | `b` | `e` | `bab` | {a:1,b:2} ≠ | — |
| 7 | `a` | `b` | `aba` | {a:2,b:1} ≠ | — |
| 8 | `c` | `a` | `bac` | {a:1,b:1,c:1} = | **6** |
| 9 | `d` | `b` | `acd` | {a:1,c:1,d:1} ≠ | — |

Result: `[0, 6]` ✔

---

## Approach 3 — Sliding Window with Match Counter (Optimal)

### Intuition

Approach 2 re-checks all 26 slots per slide — cheap, but avoidable. Keep a single integer `matches` = the number of letters whose current window count **exactly equals** the needed count. A window is an anagram precisely when `matches == 26`. Since each slide changes only two letters' counts, we can fix `matches` by looking at each affected letter's count as it crosses into or out of equality — O(1) per step, no inner loop.

### Algorithm

1. Build `need[26]` from `p`. Start `matches` = number of letters `c` with `need[c] == 0` (already satisfied at window count 0).
2. Define `add(c)` / `remove(c)`: before changing `win[c]`, if it was equal to `need[c]` do `matches--`; change it; if now equal, do `matches++`.
3. For each `r` in `[0, n)`: `add(s[r])`; if `r >= m`, `remove(s[r-m])`; if `r >= m-1` and `matches == 26`, record `r − m + 1`.

### Complexity

- **Time:** O(n) — O(1) per character; the only 26-loop is the one-time setup.
- **Space:** O(1) — two fixed arrays and a counter.

### Code

```go
func slidingWindowMatchCounter(s string, p string) []int {
	n, m := len(s), len(p)
	res := []int{}
	if m > n {
		return res
	}
	var need, win [26]int
	for i := 0; i < m; i++ {
		need[p[i]-'a']++
	}
	matches := 0
	// Letters that p does NOT use start already satisfied (both counts 0).
	for c := 0; c < 26; c++ {
		if need[c] == 0 {
			matches++
		}
	}

	// add incorporates character c (index into window) and repairs `matches`.
	add := func(c int) {
		if win[c] == need[c] { // was equal → about to break equality
			matches--
		}
		win[c]++
		if win[c] == need[c] { // reached equality
			matches++
		}
	}
	// remove drops character c from the window and repairs `matches`.
	remove := func(c int) {
		if win[c] == need[c] { // was equal → about to break equality
			matches--
		}
		win[c]--
		if win[c] == need[c] { // reached equality
			matches++
		}
	}

	for r := 0; r < n; r++ {
		add(int(s[r] - 'a')) // extend window to the right
		if r >= m {
			remove(int(s[r-m] - 'a')) // shrink from the left to keep width m
		}
		if r >= m-1 && matches == 26 { // full-width window AND all 26 letters agree
			res = append(res, r-m+1)
		}
	}
	return res
}
```

### Dry Run

Example 1: `s = "cbaebabacd", p = "abc"`, `m = 3`. `need = {a:1,b:1,c:1}`; the other 23 letters have `need = 0`, so `matches` starts at **23**.

Fill the first window `c, b, a` (each moves its letter from count 0→1, i.e. from `need` mismatch to match):

| r | action | affected letter | matches after | window ready? | matches==26? |
|---|--------|-----------------|---------------|---------------|--------------|
| 0 | add `c` | c: 0→1 (=need) | 24 | no (r<2) | — |
| 1 | add `b` | b: 0→1 (=need) | 25 | no | — |
| 2 | add `a` | a: 0→1 (=need) | 26 | yes | **yes → record 0** |
| 3 | add `e` (e:0→1, breaks) ; remove `c` (c:1→0, breaks) | e−, c− | 24 | yes | no |
| 4 | add `b` (b:1→2, breaks) ; remove `b` (b:2→1, fixes) | b−, b+ | 24 | yes | no |
| 5 | add `a` (a:1→2, breaks) ; remove `a` (a:2→1, fixes) | a−, a+ | 24 | yes | no |
| 6 | add `b` (b:1→2, breaks) ; remove `e` (e:1→0, fixes) | b−, e+ | 24 | yes | no |
| 7 | add `a` (a:1→2, breaks) ; remove `b` (b:2→1, fixes) | a−, b+ | 24 | yes | no |
| 8 | add `c` (c:0→1, fixes) ; remove `a` (a:2→1, fixes) | c+, a+ | 26 | yes | **yes → record 6** |
| 9 | add `d` (d:0→1, breaks) ; remove `b` (b:1→0, breaks) | d−, b− | 24 | yes | no |

Result: `[0, 6]` ✔

---

## Key Takeaways

- **Anagram ⇔ equal frequency vectors.** For a lowercase alphabet that vector is a `[26]int`; never sort when you can count.
- **Fixed-size sliding window** = advance the right edge and retract the left edge in lockstep, so exactly one character enters and one leaves per step. Only the changed buckets need updating.
- **The `matches` counter trick** upgrades an O(Σ) per-step compare to O(1): track *how many* buckets are already correct and repair the count only around the two that change. Reusable for #567 (permutation in string) and #76 (minimum window substring).
- Go bonus: **arrays are comparable by value** (`[26]int == [26]int`), which makes Approach 2 a one-liner test — but the match counter is strictly cheaper for large alphabets.

---

## Related Problems

- LeetCode #567 — Permutation in String (does *any* window match? boolean variant)
- LeetCode #76 — Minimum Window Substring (variable window + match counter)
- LeetCode #242 — Valid Anagram (single frequency comparison)
- LeetCode #49 — Group Anagrams (canonical-form hashing)
- LeetCode #30 — Substring with Concatenation of All Words (window of word-multiset)
