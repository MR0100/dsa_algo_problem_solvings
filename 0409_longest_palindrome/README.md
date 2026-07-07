# 0409 — Longest Palindrome

> LeetCode #409 · Difficulty: Easy
> **Categories:** Hash Table, String, Greedy

---

## Problem Statement

Given a string `s` which consists of lowercase or uppercase letters, return *the length of the **longest palindrome*** that can be built with those letters.

Letters are **case sensitive**, for example, `"Aa"` is not considered a palindrome.

**Example 1:**

```
Input: s = "abccccdd"
Output: 7
Explanation: One longest palindrome that can be built is "dccaccd", whose length is 7.
```

**Example 2:**

```
Input: s = "a"
Output: 1
Explanation: The longest palindrome that can be built is "a", whose length is 1.
```

**Constraints:**

- `1 <= s.length <= 2000`
- `s` consists of lowercase **and/or** uppercase English letters only.

---

## Company Frequency

| Company    | Frequency        | Last Reported |
|------------|------------------|---------------|
| Amazon     | ★★★★☆ High       | 2024          |
| Google     | ★★★☆☆ Medium     | 2023          |
| Microsoft  | ★★★☆☆ Medium     | 2023          |
| Meta       | ★★☆☆☆ Low        | 2023          |
| Apple      | ★★☆☆☆ Low        | 2022          |

> ⚠️ Frequency data is crowd-sourced from LeetCode Discuss, Glassdoor, and
> community interview reports. Treat as a signal, not a guarantee.

---

## DSA Concepts Used

- **Hash Map / Frequency Counting** — the answer depends only on how many times each character appears, so we tally counts and reason about their parity → see [`/dsa/hash_map.md`](/dsa/hash_map.md)
- **Counting with a fixed array** — because the alphabet is a small fixed set (letters ⊂ ASCII), a flat `[128]int` count array replaces the map for O(1) space, the same index-directly-into-buckets idea behind counting sort → see [`/dsa/counting_sort.md`](/dsa/counting_sort.md)
- **Greedy** — we greedily take every full pair of each character and then grab one leftover single for the center; no character choice is ever regretted → see [`/dsa/greedy.md`](/dsa/greedy.md)

---

## Approaches Overview

| # | Approach | Time | Space | When to use |
|---|----------|------|-------|-------------|
| 1 | Hash Map Frequency | O(n) | O(k) | Clear and general; works for any alphabet |
| 2 | Fixed Count Array + Odd Tally (Optimal) | O(n) | O(1) | Same idea, O(1) space; track odd-count parity on the fly |

`k` = number of distinct characters (≤ 52 here); `n = len(s)`.

---

## Approach 1 — Hash Map Frequency

### Intuition

A palindrome reads the same forwards and backwards, so every character must appear in **mirrored pairs** around the center — an even number of times — with **at most one** exception: a single character may sit exactly in the middle. So the plan is greedy:

- For a character seen `c` times, use as many as form full pairs: `c − (c mod 2)` (i.e. `c` rounded down to even).
- If **any** character has an odd count, one unpaired character can be placed in the center, adding `1`.

### Algorithm

1. Count occurrences of each character in a map.
2. Initialise `length = 0`, `hasOdd = false`.
3. For each count `c`: add `c − (c mod 2)` to `length`; if `c` is odd, set `hasOdd = true`.
4. If `hasOdd`, add `1` for the center.
5. Return `length`.

### Complexity

- **Time:** O(n) — one pass to build the counts, one pass over at most `k` buckets.
- **Space:** O(k) — the frequency map, bounded by 52 distinct letters.

### Code

```go
func hashMap(s string) int {
	counts := make(map[byte]int) // character → how many times it appears
	for i := 0; i < len(s); i++ {
		counts[s[i]]++ // tally this character
	}

	length := 0        // running palindrome length
	hasOdd := false    // did we see any character with an odd count?
	for _, c := range counts {
		length += c - c%2 // use the largest even number of this char (full pairs)
		if c%2 == 1 {
			hasOdd = true // at least one odd-count char exists → a center is available
		}
	}
	if hasOdd {
		length++ // one leftover single character can occupy the exact middle
	}
	return length
}
```

### Dry Run

Example 1: `s = "abccccdd"`.

Counts: `a:1, b:1, c:4, d:2`.

| Char | count c | c − (c mod 2) added | odd? | length so far | hasOdd |
|------|---------|---------------------|------|---------------|--------|
| a | 1 | 0 | yes | 0 | true |
| b | 1 | 0 | yes | 0 | true |
| c | 4 | 4 | no | 4 | true |
| d | 2 | 2 | no | 6 | true |

`hasOdd = true` → add 1 → `length = 7`.

Result: **7** ✔ (e.g. `"dccaccd"` — the leftover `a` is the center; `b` is discarded.)

---

## Approach 2 — Fixed Count Array + Odd Tally (Optimal)

### Intuition

Same parity insight, but avoid the map. Index a flat `[128]int` array by ASCII code (letters live in `0..127`), and maintain a live counter `oddCount` = how many characters currently have an odd count. Each character that ends up with an odd count wastes exactly **one** unpaired copy. We are allowed to keep **one** of those wasted singles as the palindrome's center; the rest must be dropped. Therefore:

- if `oddCount > 0`: answer = `n − oddCount + 1` (discard `oddCount − 1` singles, keep one center);
- else: answer = `n` (everything pairs perfectly).

Flipping `oddCount` incrementally (`+1` when a slot turns odd, `−1` when it turns even) means the whole thing is a single pass with no second loop.

### Algorithm

1. `count[128]`, `oddCount = 0`.
2. For each character: increment its slot; if the slot is now odd, `oddCount++`, else `oddCount--`.
3. If `oddCount > 0`, return `n − oddCount + 1`.
4. Else return `n`.

### Complexity

- **Time:** O(n) — one linear scan of the string; no post-processing loop.
- **Space:** O(1) — a fixed 128-entry array, independent of input size.

### Code

```go
func countArray(s string) int {
	var count [128]int // ASCII code → occurrence count (letters fit in 0..127)
	oddCount := 0      // how many characters currently have an odd count
	for i := 0; i < len(s); i++ {
		count[s[i]]++ // record this character
		if count[s[i]]%2 == 1 {
			oddCount++ // count for this char just turned odd
		} else {
			oddCount-- // it turned even again — this char is fully paired now
		}
	}

	n := len(s)
	if oddCount > 0 {
		// Each odd-count char contributes one unpaired single; we may keep exactly
		// one of them in the center, so we discard (oddCount - 1) characters.
		return n - oddCount + 1
	}
	return n // all counts even → the whole string is usable
}
```

### Dry Run

Example 1: `s = "abccccdd"`, `n = 8`.

| i | char | slot count after | parity now | oddCount after |
|---|------|------------------|------------|----------------|
| 0 | a | a=1 | odd | 1 |
| 1 | b | b=1 | odd | 2 |
| 2 | c | c=1 | odd | 3 |
| 3 | c | c=2 | even | 2 |
| 4 | c | c=3 | odd | 3 |
| 5 | c | c=4 | even | 2 |
| 6 | d | d=1 | odd | 3 |
| 7 | d | d=2 | even | 2 |

Final `oddCount = 2` (`a` and `b` remain odd). Since `oddCount > 0`: answer = `n − oddCount + 1 = 8 − 2 + 1 = 7`.

Result: **7** ✔

---

## Key Takeaways

- **A palindrome tolerates at most one odd-count character** (the center); everything else must pair. This single parity fact solves the whole problem — no need to actually build the palindrome.
- **`c − (c mod 2)`** is the clean way to "round a count down to the nearest even" when you can only use matched pairs.
- **Two equivalent framings:** sum the even-floored counts and add 1 if any odd exists; *or* start from `n` and subtract the wasted singles, adding one back for the center (`n − oddCount + 1`). They give the same answer — pick whichever you can derive fastest.
- **Small fixed alphabet ⇒ array, not map.** When keys are letters/ASCII, a flat count array is faster and O(1) space, and lets you maintain derived quantities (like `oddCount`) incrementally.
- **Watch case sensitivity:** `'A'` and `'a'` are different characters here, which is exactly why a 128-slot ASCII array (not 26) is the safe choice.

---

## Related Problems

- LeetCode #5 — Longest Palindromic Substring (contiguous, expand around center)
- LeetCode #125 — Valid Palindrome (two-pointer check)
- LeetCode #266 — Palindrome Permutation (can a palindrome be formed at all?)
- LeetCode #234 — Palindrome Linked List (palindrome verification)
- LeetCode #267 — Palindrome Permutation II (build all palindromes from the letters)
