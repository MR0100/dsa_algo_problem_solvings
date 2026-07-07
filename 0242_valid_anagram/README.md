# 0242 — Valid Anagram

> LeetCode #242 · Difficulty: Easy
> **Categories:** Hash Table, String, Sorting, Counting

---

## Problem Statement

Given two strings `s` and `t`, return `true` *if* `t` *is an anagram of* `s`*, and* `false` *otherwise*.

An **anagram** is a word or phrase formed by rearranging the letters of a different word or phrase, typically using all the original letters exactly once.

**Example 1:**

```
Input: s = "anagram", t = "nagaram"
Output: true
```

**Example 2:**

```
Input: s = "rat", t = "car"
Output: false
```

**Constraints:**

- `1 <= s.length, t.length <= 5 * 10^4`
- `s` and `t` consist of lowercase English letters.

**Follow up:** What if the inputs contain Unicode characters? How would you adapt your solution to such a case?

---

## Company Frequency

| Company    | Frequency        | Last Reported |
|------------|------------------|---------------|
| Amazon     | ★★★★★ Very High  | 2024          |
| Google     | ★★★★☆ High       | 2024          |
| Microsoft  | ★★★★☆ High       | 2024          |
| Bloomberg  | ★★★☆☆ Medium     | 2024          |
| Meta       | ★★★☆☆ Medium     | 2023          |
| Apple      | ★★★☆☆ Medium     | 2023          |

> ⚠️ Frequency data is crowd-sourced from LeetCode Discuss, Glassdoor, and
> community interview reports. Treat as a signal, not a guarantee.

---

## DSA Concepts Used

- **Hash Map / Frequency Counting** — an anagram is a character multiset match; count occurrences and check they cancel → see [`/dsa/hash_map.md`](/dsa/hash_map.md)
- **Sorting** — sorting canonicalizes a multiset, so anagrams sort to identical strings → see [`/dsa/sorting.md`](/dsa/sorting.md)
- **String Processing** — iterating characters and mapping letters to array indices → see [`/dsa/string_algorithms.md`](/dsa/string_algorithms.md)

---

## Approaches Overview

| # | Approach | Time | Space | When to use |
|---|----------|------|-------|-------------|
| 1 | Sorting | O(n log n) | O(n) | Shortest to write; no assumptions on alphabet |
| 2 | Hash Map Count | O(n) | O(k) | Unicode / arbitrary alphabet (follow-up) |
| 3 | Fixed Array Counter (Optimal) | O(n) | O(1) | Known small alphabet (a–z) — fastest |

(n = string length, k = number of distinct characters.)

---

## Approach 1 — Sorting

### Intuition
Two strings are anagrams exactly when they contain the same characters with the same multiplicities. Sorting turns a multiset into one canonical ordering, so anagrams become byte-for-byte equal after sorting.

### Algorithm
1. If lengths differ, return `false` immediately.
2. Sort the character slice of `s`.
3. Sort the character slice of `t`.
4. Return whether the two sorted strings are equal.

### Complexity
- **Time:** O(n log n) — the two sorts dominate.
- **Space:** O(n) — mutable rune slices to sort (Go strings are immutable).

### Code
```go
func sortingApproach(s string, t string) bool {
	if len(s) != len(t) {
		return false
	}
	a := []rune(s)
	b := []rune(t)
	sort.Slice(a, func(i, j int) bool { return a[i] < a[j] })
	sort.Slice(b, func(i, j int) bool { return b[i] < b[j] })
	return string(a) == string(b)
}
```

### Dry Run
Trace `s = "anagram"`, `t = "nagaram"`:

| Step | Variable | Value |
|------|----------|-------|
| 1 | lengths | 7 == 7 → continue |
| 2 | sorted `a` | `a a a g m n r` |
| 3 | sorted `b` | `a a a g m n r` |
| 4 | equal? | `"aaagmnr" == "aaagmnr"` → **true** |

Result: `true`. ✓

---

## Approach 2 — Hash Map Frequency Count

### Intuition
Add `+1` for each character of `s` and `-1` for each character of `t` in one shared map. If they are anagrams, every character's net count is zero. A map covers any alphabet, including Unicode (the follow-up).

### Algorithm
1. If lengths differ, return `false`.
2. For each rune in `s`, increment its count.
3. For each rune in `t`, decrement its count.
4. If every value is `0`, return `true`; else `false`.

### Complexity
- **Time:** O(n) — three linear passes.
- **Space:** O(k) — one map entry per distinct character.

### Code
```go
func hashMap(s string, t string) bool {
	if len(s) != len(t) {
		return false
	}
	count := make(map[rune]int)
	for _, r := range s {
		count[r]++
	}
	for _, r := range t {
		count[r]--
	}
	for _, c := range count {
		if c != 0 {
			return false
		}
	}
	return true
}
```

### Dry Run
Trace `s = "anagram"`, `t = "nagaram"`:

| Phase | Map state |
|-------|-----------|
| after counting `s` | `a:3, n:1, g:1, r:1, m:1` |
| after subtracting `t` | `a:0, n:0, g:0, r:0, m:0` |
| final check | all zero → **true** |

Result: `true`. ✓

---

## Approach 3 — Fixed Array Counter (Optimal, a–z)

### Intuition
When the alphabet is the fixed set of lowercase letters, replace the hash map with a 26-element array indexed by `letter - 'a'`. No hashing, contiguous memory — the fastest form. Increment on `s`, decrement on `t` in one fused loop.

### Algorithm
1. If lengths differ, return `false`.
2. Create `counts[26]` initialized to zero.
3. For each index `i`: `counts[s[i]-'a']++` and `counts[t[i]-'a']--`.
4. If any slot is non-zero, return `false`; otherwise `true`.

### Complexity
- **Time:** O(n) — one pass touching both strings.
- **Space:** O(1) — a constant 26-element array.

### Code
```go
func arrayCount(s string, t string) bool {
	if len(s) != len(t) {
		return false
	}
	var counts [26]int
	for i := 0; i < len(s); i++ {
		counts[s[i]-'a']++
		counts[t[i]-'a']--
	}
	for _, c := range counts {
		if c != 0 {
			return false
		}
	}
	return true
}
```

### Dry Run
Trace `s = "anagram"`, `t = "nagaram"` (showing only touched letters):

| i | s[i] | t[i] | counts['a'] | counts['n'] | counts['g'] | counts['r'] | counts['m'] |
|---|------|------|-------------|-------------|-------------|-------------|-------------|
| 0 | a | n | +1 | -1 | 0 | 0 | 0 |
| 1 | n | a | 0 | 0 | 0 | 0 | 0 |
| 2 | a | g | +1 | 0 | -1 | 0 | 0 |
| 3 | g | a | 0 | 0 | 0 | 0 | 0 |
| 4 | r | r | 0 | 0 | 0 | 0 | 0 |
| 5 | a | a | 0 | 0 | 0 | 0 | 0 |
| 6 | m | m | 0 | 0 | 0 | 0 | 0 |

All slots end at 0 → **true**. ✓

---

## Key Takeaways
- **Anagram = equal character multiset.** Two canonical ways to check: sort both, or count and cancel.
- The **length check** is a free early exit — unequal lengths can never be anagrams.
- Prefer a **fixed-size array counter** when the alphabet is small and known; use a **map** for arbitrary/Unicode input (the stated follow-up).
- Counting in a **single fused loop** (`+` on one string, `-` on the other) halves passes versus two separate loops.

---

## Related Problems
- LeetCode #49 — Group Anagrams (bucket strings by anagram signature)
- LeetCode #438 — Find All Anagrams in a String (sliding-window frequency match)
- LeetCode #383 — Ransom Note (subset frequency check)
- LeetCode #567 — Permutation in String (window anagram check)
