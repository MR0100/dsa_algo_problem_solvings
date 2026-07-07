# 0340 — Longest Substring with At Most K Distinct Characters

> LeetCode #340 · Difficulty: Medium
> **Categories:** Hash Table, String, Sliding Window

---

## Problem Statement

Given a string `s` and an integer `k`, return *the length of the longest substring of `s` that contains at most `k` distinct characters*.

**Example 1:**

```
Input: s = "eceba", k = 2
Output: 3
Explanation: The substring is "ece" with length 3.
```

**Example 2:**

```
Input: s = "aa", k = 1
Output: 2
Explanation: The substring is "aa" with length 2.
```

**Constraints:**

- `1 <= s.length <= 5 * 10^4`
- `0 <= k <= 50`

---

## Company Frequency

| Company    | Frequency        | Last Reported |
|------------|------------------|---------------|
| Google     | ★★★★☆ High       | 2024          |
| Amazon     | ★★★★☆ High       | 2024          |
| Facebook   | ★★★☆☆ Medium     | 2023          |
| Microsoft  | ★★☆☆☆ Low        | 2023          |
| Bloomberg  | ★★☆☆☆ Low        | 2022          |

> ⚠️ Frequency data is crowd-sourced from LeetCode Discuss, Glassdoor, and
> community interview reports. Treat as a signal, not a guarantee.

---

## DSA Concepts Used

- **Sliding Window** — the optimal answer maintains a variable-size window that never holds more than `k` distinct characters, expanding right and shrinking left → see [`/dsa/sliding_window.md`](/dsa/sliding_window.md)
- **Hash Map** — a character→count map tracks the distinct-character set inside the window in O(1) per update → see [`/dsa/hash_map.md`](/dsa/hash_map.md)
- **String Algorithms** — scanning and windowing over the characters of `s` → see [`/dsa/string_algorithms.md`](/dsa/string_algorithms.md)

---

## Approaches Overview

| # | Approach | Time | Space | When to use |
|---|----------|------|-------|-------------|
| 1 | Brute Force (all substrings) | O(n²) | O(k) | Baseline; clear but too slow for n = 5·10⁴ |
| 2 | Sliding Window + Hash Map (Optimal) | O(n) | O(k) | General alphabet; the standard interview answer |
| 3 | Sliding Window + Fixed Array (Optimal) | O(n) | O(1) | ASCII input; avoids hashing/allocation |

---

## Approach 1 — Brute Force (all substrings)

### Intuition

Enumerate every substring by its start `i` and extend the end `j`, keeping a frequency set of the characters seen. As soon as the set grows past `k`, no longer substring from that start is valid, so stop extending. Record the longest valid length across all starts.

### Algorithm

1. If `k == 0`, return `0`.
2. For each start `i`, reset a distinct-set and extend `j` from `i`.
3. Add `s[j]`; if distinct count `> k`, break.
4. Otherwise update `best = max(best, j-i+1)`.

### Complexity

- **Time:** O(n²) — `n` starts, each extending up to `n` characters with O(1) set updates.
- **Space:** O(k) — the distinct-set holds at most `k+1` characters before breaking.

### Code

```go
func bruteForce(s string, k int) int {
	if k == 0 {
		return 0 // no characters allowed → empty substring only
	}
	best := 0
	for i := 0; i < len(s); i++ {
		freq := map[byte]int{} // distinct chars in the current window s[i..j]
		for j := i; j < len(s); j++ {
			freq[s[j]]++ // extend the window by one character
			if len(freq) > k {
				break // too many distinct chars; no longer window from this i works
			}
			if j-i+1 > best {
				best = j - i + 1 // record a longer valid window
			}
		}
	}
	return best
}
```

### Dry Run

Example 1: `s = "eceba", k = 2`. Only start `i = 0` (which yields the answer) shown:

| i | j | s[j] | freq (set) | distinct | valid? | best |
|---|---|------|------------|----------|--------|------|
| 0 | 0 | e | {e} | 1 | yes | 1 |
| 0 | 1 | c | {e,c} | 2 | yes | 2 |
| 0 | 2 | e | {e,c} | 2 | yes | 3 |
| 0 | 3 | b | {e,c,b} | 3 | no → break | 3 |

Later starts never beat 3. Result: `3` ✔

---

## Approach 2 — Sliding Window with Hash Map (Optimal)

### Intuition

Keep a window `[left, right]` whose distinct-character count is always `≤ k`. Push `right` forward one character at a time. If adding it makes the map hold more than `k` distinct characters, advance `left` — decrementing counts and deleting a key when its count reaches 0 — until the window is valid again. Each index enters and exits the window at most once, so the total work is linear.

### Algorithm

1. If `k == 0`, return `0`.
2. Expand: `freq[s[right]]++`.
3. While `len(freq) > k`: `freq[s[left]]--`, delete if 0, `left++`.
4. `best = max(best, right-left+1)`.

### Complexity

- **Time:** O(n) — `right` and `left` each traverse the string once.
- **Space:** O(k) — the map holds at most `k+1` distinct keys.

### Code

```go
func slidingWindow(s string, k int) int {
	if k == 0 {
		return 0 // vacuously no valid non-empty window
	}
	freq := map[byte]int{} // char → count within [left, right]
	best, left := 0, 0
	for right := 0; right < len(s); right++ {
		freq[s[right]]++ // include the new right character
		// Shrink until at most k distinct characters remain.
		for len(freq) > k {
			freq[s[left]]-- // one fewer occurrence of the leftmost char
			if freq[s[left]] == 0 {
				delete(freq, s[left]) // its last copy left the window
			}
			left++ // move the left edge inward
		}
		if right-left+1 > best {
			best = right - left + 1 // widest valid window so far
		}
	}
	return best
}
```

### Dry Run

Example 1: `s = "eceba", k = 2`.

| right | s[right] | freq after add | len>k? shrink | left | window | best |
|-------|----------|----------------|---------------|------|--------|------|
| 0 | e | {e:1} | no | 0 | "e" | 1 |
| 1 | c | {e:1,c:1} | no | 0 | "ec" | 2 |
| 2 | e | {e:2,c:1} | no | 0 | "ece" | 3 |
| 3 | b | {e:2,c:1,b:1} → 3>2 | drop e (2→1), left=1 | 1 | still {c,e,b}=3 → drop c(→0,del), left=2 | 2 |
| 3 | (after shrink) | {e:1,b:1} | ok | 2 | "eb" | 3 |
| 4 | a | {e:1,b:1,a:1} → 3>2 | drop e(→0,del), left=3 | 3 | {b,a}=2 | 3 |

Result: `3` ✔

---

## Approach 3 — Sliding Window with Fixed Array (Optimal, ASCII)

### Intuition

For an ASCII alphabet the map is overkill: a fixed `[128]int` counter plus a running `distinct` integer reproduces `len(freq)` exactly. Increment `distinct` when a slot goes `0→1`; decrement it when a slot goes `1→0`. All updates are O(1) with zero allocation.

### Algorithm

1. If `k == 0`, return `0`.
2. Add `s[right]`: if its count was 0, `distinct++`, then `count++`.
3. While `distinct > k`: decrement `count[s[left]]`, if it hits 0 `distinct--`, `left++`.
4. `best = max(best, right-left+1)`.

### Complexity

- **Time:** O(n) — one pass, O(1) per step.
- **Space:** O(1) — a fixed 128-entry array regardless of `n` or `k`.

### Code

```go
func slidingWindowArray(s string, k int) int {
	if k == 0 {
		return 0
	}
	var count [128]int // ASCII code → occurrences within the window
	distinct := 0      // number of characters currently present (count > 0)
	best, left := 0, 0
	for right := 0; right < len(s); right++ {
		if count[s[right]] == 0 {
			distinct++ // a brand-new character entered the window
		}
		count[s[right]]++
		for distinct > k { // too many distinct → shrink from the left
			count[s[left]]--
			if count[s[left]] == 0 {
				distinct-- // the leftmost char's last copy left
			}
			left++
		}
		if right-left+1 > best {
			best = right - left + 1
		}
	}
	return best
}
```

### Dry Run

Example 1: `s = "eceba", k = 2`.

| right | char | distinct after add | shrink? | left | window | best |
|-------|------|--------------------|---------|------|--------|------|
| 0 | e | 1 | no | 0 | "e" | 1 |
| 1 | c | 2 | no | 0 | "ec" | 2 |
| 2 | e | 2 | no | 0 | "ece" | 3 |
| 3 | b | 3 | drop e,e→still, drop... left→2, distinct=2 | 2 | "eb" | 3 |
| 4 | a | 3 | drop e→distinct 2, left=3 | 3 | "ba" | 3 |

Result: `3` ✔

---

## Key Takeaways

- **"At most K" ⇒ variable-size sliding window.** Grow right greedily; shrink left only when the constraint (distinct ≤ k) breaks. Amortized O(n) because each index is added/removed once.
- A **count map + `len(map)`** is the general distinct-tracker; for a small fixed alphabet, a **counter array + a `distinct` integer** removes hashing and allocation.
- The sibling **"exactly K distinct"** reduces to `atMost(K) - atMost(K-1)` — a common follow-up (see #992).
- Delete the key (or decrement `distinct`) the moment a count hits 0, or `len(freq)` overcounts and the window shrinks too aggressively.

---

## Related Problems

- LeetCode #159 — Longest Substring with At Most Two Distinct Characters (k = 2 special case)
- LeetCode #3 — Longest Substring Without Repeating Characters (distinct = window length)
- LeetCode #992 — Subarrays with K Different Integers (exactly-K via atMost trick)
- LeetCode #76 — Minimum Window Substring (shrinking-window on a coverage constraint)
