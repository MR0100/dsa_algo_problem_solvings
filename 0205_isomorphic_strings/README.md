# 0205 — Isomorphic Strings

> LeetCode #205 · Difficulty: Easy
> **Categories:** Hash Table, String

---

## Problem Statement

Given two strings `s` and `t`, *determine if they are isomorphic*.

Two strings `s` and `t` are **isomorphic** if the characters in `s` can be replaced to get `t`.

All occurrences of a character must be replaced with another character while preserving the order of characters. No two characters may map to the same character, but a character may map to itself.

**Example 1:**

```
Input: s = "egg", t = "add"
Output: true
Explanation:
The strings s and t can be made identical by:
- Mapping 'e' to 'a'.
- Mapping 'g' to 'd'.
```

**Example 2:**

```
Input: s = "foo", t = "bar"
Output: false
Explanation:
The strings s and t can not be made identical as 'o' needs to be mapped to both 'a' and 'r'.
```

**Example 3:**

```
Input: s = "paper", t = "title"
Output: true
```

**Constraints:**

- `1 <= s.length <= 5 * 10^4`
- `t.length == s.length`
- `s` and `t` consist of any valid ascii characters.

---

## Company Frequency

| Company    | Frequency        | Last Reported |
|------------|------------------|---------------|
| LinkedIn   | ★★★★★ Very High  | 2024          |
| Google     | ★★★★☆ High       | 2024          |
| Amazon     | ★★★☆☆ Medium     | 2024          |
| Microsoft  | ★★★☆☆ Medium     | 2023          |
| Apple      | ★★☆☆☆ Low        | 2023          |
| Bloomberg  | ★★☆☆☆ Low        | 2022          |

> ⚠️ Frequency data is crowd-sourced from LeetCode Discuss, Glassdoor, and
> community interview reports. Treat as a signal, not a guarantee.

---

## DSA Concepts Used

- **Hash Map** — build the character → character replacement table on the fly; a second, reversed map enforces that no two characters share an image (injectivity) → see [`/dsa/hash_map.md`](/dsa/hash_map.md)
- **String Algorithms — canonical encoding** — replacing each character by its first/last-occurrence index normalises both strings; isomorphic ⇔ identical encodings (the same trick as pattern matching in #290/#890) → see [`/dsa/string_algorithms.md`](/dsa/string_algorithms.md)

---

## Approaches Overview

| # | Approach | Time | Space | When to use |
|---|----------|------|-------|-------------|
| 1 | Brute Force (Pairwise Consistency) | O(n²) | O(1) | Definition-level check; ~1.25·10⁹ pairs at n = 5·10⁴ — too slow |
| 2 | Two Hash Maps | O(n) | O(k), k = alphabet | The natural interview answer; generalises to any alphabet (runes) |
| 3 | First-Occurrence Encoding (Optimal) | O(n) | O(1) — two [256] arrays | Fastest constants; fixed byte alphabet makes maps unnecessary |

---

## Approach 1 — Brute Force (Pairwise Consistency)

### Intuition

Forget maps and look at what a valid replacement *implies*: the two strings must have identical **equality structure**. For any two positions `j < i`, if `s[j] == s[i]` then `t[j]` must equal `t[i]` (one character, one image), and if `s[j] != s[i]` then `t[j]` must differ from `t[i]` (two characters may not share an image). So compare the boolean `s[i] == s[j]` with the boolean `t[i] == t[j]` for every pair — any disagreement in either direction kills the isomorphism, and if all pairs agree, the position-wise pairing itself is a valid mapping.

### Algorithm

1. For `i` from 0 to n−1:
   1. For `j` from 0 to i−1:
      1. If `(s[i] == s[j]) != (t[i] == t[j])`, return `false`.
2. All pairs consistent → return `true`.

### Complexity

- **Time:** O(n²) — n(n−1)/2 position pairs, O(1) work each; at the constraint max (5·10⁴) that is ~1.25·10⁹ comparisons — TLE territory.
- **Space:** O(1) — literally nothing stored.

### Code

```go
func bruteForce(s string, t string) bool {
	for i := 0; i < len(s); i++ {
		for j := 0; j < i; j++ {
			// The equality patterns must agree at every pair of positions:
			// a mismatch in either direction breaks the bijection.
			if (s[i] == s[j]) != (t[i] == t[j]) {
				return false
			}
		}
	}
	return true // identical equality structure → an isomorphism exists
}
```

### Dry Run

Example 1: `s = "egg", t = "add"`.

| Step | i | j | s[i], s[j] equal? | t[i], t[j] equal? | Agree? |
|------|---|---|--------------------|--------------------|--------|
| 1 | 1 | 0 | g, e → no | d, a → no | yes ✓ |
| 2 | 2 | 0 | g, e → no | d, a → no | yes ✓ |
| 3 | 2 | 1 | g, g → yes | d, d → yes | yes ✓ |

No pair disagrees → `true` ✔ (For `s="foo", t="bar"`: pair i=2, j=1 gives `o == o` → yes but `r == a` → no — disagreement → `false`.)

---

## Approach 2 — Two Hash Maps

### Intuition

Simulate the replacement while reading the strings once. A forward map `sToT` pins down consistency: the first time a character `a` appears, its image is fixed forever; seeing `a` later with a different image is a contradiction. But the forward map alone misses the "no two characters may map to the same character" rule — in `s = "badc", t = "baba"`, forward gives `b→b, a→a, d→b` with no forward conflict, yet `d` steals the image `b` that already belongs to `b`. The backward map `tToS` catches exactly that: each t-character remembers which s-character owns it. Both maps must stay conflict-free for the strings to be isomorphic (a bijection is a consistent function whose inverse is also a function).

### Algorithm

1. Initialise empty maps `sToT` and `tToS`.
2. For each index `i`, let `a = s[i]`, `b = t[i]`:
   1. If `sToT[a]` exists and `sToT[a] != b` → return `false` (inconsistent forward mapping).
   2. If `tToS[b]` exists and `tToS[b] != a` → return `false` (image `b` already owned by a different character).
   3. Set `sToT[a] = b` and `tToS[b] = a`.
3. Return `true`.

### Complexity

- **Time:** O(n) — one pass; each step is two average-O(1) lookups and two inserts.
- **Space:** O(k) — at most one entry per distinct character in each map; k ≤ 256 for byte alphabets, so effectively constant.

### Code

```go
func twoHashMaps(s string, t string) bool {
	sToT := map[byte]byte{} // established forward mapping s-char → t-char
	tToS := map[byte]byte{} // established backward mapping t-char → s-char
	for i := 0; i < len(s); i++ {
		a, b := s[i], t[i]
		// Forward check: a must always map to the same image.
		if mapped, ok := sToT[a]; ok && mapped != b {
			return false // a already maps elsewhere — inconsistent
		}
		// Backward check: b must not already be claimed by another s-char.
		if mapped, ok := tToS[b]; ok && mapped != a {
			return false // two s-chars would share the image b — not injective
		}
		sToT[a] = b // (re-)recording an identical pair is harmless
		tToS[b] = a
	}
	return true
}
```

### Dry Run

Example 1: `s = "egg", t = "add"`.

| Step | i | a = s[i] | b = t[i] | sToT[a] conflict? | tToS[b] conflict? | sToT after | tToS after |
|------|---|----------|----------|-------------------|-------------------|------------------|------------------|
| 1 | 0 | e | a | absent — ok | absent — ok | {e→a} | {a→e} |
| 2 | 1 | g | d | absent — ok | absent — ok | {e→a, g→d} | {a→e, d→g} |
| 3 | 2 | g | d | g→d, matches | d→g, matches | unchanged | unchanged |

End of string → `true` ✔ (For `s="badc", t="baba"` step 3 has a=d, b=b: `sToT[d]` absent, but `tToS[b] = b ≠ d` → backward conflict → `false`.)

---

## Approach 3 — First-Occurrence Encoding (Optimal)

### Intuition

Normalise both strings into a form where the actual letters vanish and only the *repeat pattern* remains: replace each character by the position where it was last seen. `"egg"` and `"add"` both become "new, new, repeats-position-1", so they are isomorphic; `"foo"` and `"bar"` diverge at index 2 ("repeats-position-1" vs "new"). Implementation trick: instead of building the encodings, verify them on the fly — at each position `i`, `s[i]` and `t[i]` must have the same last-seen value. Two flat `[256]int` arrays (ASCII alphabet) hold the values; storing `i + 1` lets the zero value of the array mean "never seen". This replaces hash maps with two L1-cache-resident arrays — same O(n), much smaller constants.

### Algorithm

1. Declare `lastSeenS`, `lastSeenT` as `[256]int` (all zeros = never seen).
2. For each index `i`:
   1. If `lastSeenS[s[i]] != lastSeenT[t[i]]`, return `false` — one character is new while the other is not, or they last occurred at different places.
   2. Set both `lastSeenS[s[i]]` and `lastSeenT[t[i]]` to `i + 1`.
3. Return `true`.

### Complexity

- **Time:** O(n) — a single pass with two array reads and two array writes per character.
- **Space:** O(1) — two fixed 256-entry arrays (2 KB total), independent of input length.

### Code

```go
func firstOccurrenceEncoding(s string, t string) bool {
	var lastSeenS, lastSeenT [256]int // last position + 1; 0 = never seen
	for i := 0; i < len(s); i++ {
		// Both characters must have the same "history": either both brand
		// new (0 == 0) or both last seen at the same index.
		if lastSeenS[s[i]] != lastSeenT[t[i]] {
			return false
		}
		lastSeenS[s[i]] = i + 1 // +1 shift keeps index 0 distinguishable
		lastSeenT[t[i]] = i + 1
	}
	return true
}
```

### Dry Run

Example 1: `s = "egg", t = "add"`.

| Step | i | s[i] | t[i] | lastSeenS[s[i]] | lastSeenT[t[i]] | Equal? | Store i+1 |
|------|---|------|------|------------------|------------------|--------|-----------|
| 1 | 0 | e | a | 0 (new) | 0 (new) | yes ✓ | lastSeenS[e]=1, lastSeenT[a]=1 |
| 2 | 1 | g | d | 0 (new) | 0 (new) | yes ✓ | lastSeenS[g]=2, lastSeenT[d]=2 |
| 3 | 2 | g | d | 2 | 2 | yes ✓ | lastSeenS[g]=3, lastSeenT[d]=3 |

All positions agree → `true` ✔ (For `s="foo", t="bar"` at i=2: `lastSeenS[o] = 2` but `lastSeenT[r] = 0` → mismatch → `false`.)

---

## Key Takeaways

- **Bijection = two one-way checks.** A forward map catches "one character, two images"; only the reverse map catches "two characters, one image". Forgetting the second check (test case `s="badc", t="baba"`) is *the* classic bug in this problem.
- **Canonical-form thinking:** encoding each character as its first/last-occurrence index erases identities and keeps structure — the exact tool for "same pattern?" problems (#290 Word Pattern, #890 Find and Replace Pattern).
- **Fixed alphabet ⇒ arrays beat maps.** When keys are bytes, a `[256]int` array is an O(1)-space, cache-friendly hash map replacement; the `i + 1` shift is the standard trick to keep the zero value meaning "empty".
- The `lastSeen` comparison works because agreement at *every* prefix forces agreement of the full occurrence pattern — verifying incrementally avoids materialising the encoded arrays.
- Isomorphism here is an equivalence relation on strings; "is `s` isomorphic to `t`" is symmetric, which is why both directions must be policed equally.

---

## Related Problems

- LeetCode #290 — Word Pattern (identical bijection check, chars ↔ words)
- LeetCode #890 — Find and Replace Pattern (isomorphism test applied to many candidates)
- LeetCode #49 — Group Anagrams (canonical-form grouping, different normalisation)
- LeetCode #242 — Valid Anagram (frequency structure instead of positional structure)
- LeetCode #1153 — String Transforms Into Another String (one-way mapping without injectivity — instructive contrast)
