# 0187 — Repeated DNA Sequences

> LeetCode #187 · Difficulty: Medium
> **Categories:** Hash Table, String, Bit Manipulation, Sliding Window, Rolling Hash

---

## Problem Statement

The **DNA sequence** is composed of a series of nucleotides abbreviated as `'A'`, `'C'`, `'G'`, and `'T'`.

- For example, `"ACGAATTCCG"` is a **DNA sequence**.

When studying **DNA**, it is useful to identify repeated sequences within the DNA.

Given a string `s` that represents a **DNA sequence**, return all the **`10`-letter-long** sequences (substrings) that occur more than once in a DNA molecule. You may return the answer in **any order**.

**Example 1:**
```
Input: s = "AAAAACCCCCAAAAACCCCCCAAAAAGGGTTT"
Output: ["AAAAACCCCC","CCCCCAAAAA"]
```

**Example 2:**
```
Input: s = "AAAAAAAAAAAAA"
Output: ["AAAAAAAAAA"]
```

**Constraints:**
- `1 <= s.length <= 10^5`
- `s[i]` is either `'A'`, `'C'`, `'G'`, or `'T'`.

---

## Company Frequency

| Company    | Frequency        | Last Reported |
|------------|------------------|---------------|
| Meta       | ★★★★☆ High       | 2024          |
| Google     | ★★★☆☆ Medium     | 2024          |
| Amazon     | ★★★☆☆ Medium     | 2024          |
| Microsoft  | ★★★☆☆ Medium     | 2023          |
| LinkedIn   | ★★☆☆☆ Low        | 2023          |
| Bloomberg  | ★★☆☆☆ Low        | 2023          |

> ⚠️ Frequency data is crowd-sourced from LeetCode Discuss, Glassdoor, and
> community interview reports. Treat as a signal, not a guarantee.

---

## DSA Concepts Used

- **Hash Map / Hash Set** — counting window occurrences (or membership testing on fingerprints) turns "does this substring repeat?" into an O(1) lookup → see [`/dsa/hash_map.md`](/dsa/hash_map.md)
- **Sliding Window (fixed size)** — the candidates are exactly the n−9 windows of length 10; each step slides the window by one character → see [`/dsa/sliding_window.md`](/dsa/sliding_window.md)
- **Bit Manipulation** — a 4-letter alphabet means 2 bits per letter, so a 10-letter window packs losslessly into a 20-bit integer key → see [`/dsa/bit_manipulation.md`](/dsa/bit_manipulation.md)
- **String Algorithms (rolling hash / Rabin–Karp)** — updating the window fingerprint in O(1) per slide instead of re-hashing all 10 characters → see [`/dsa/string_algorithms.md`](/dsa/string_algorithms.md)

---

## Approaches Overview

| # | Approach | Time | Space | When to use |
|---|----------|------|-------|-------------|
| 1 | Brute Force | O(n² · L) | O(1) extra | Only to establish correctness; too slow for n = 10⁵ |
| 2 | Hash Map Counting | O(n · L) | O(n · L) | The standard interview answer; simple and fast enough |
| 3 | Rolling Hash / Bit Manipulation (Optimal) | O(n) | O(n) | Follow-up flex: drops the ×L factor and stores ints, not strings |

*(L = 10, the fixed sequence length.)*

---

## Approach 1 — Brute Force

### Intuition
A 10-letter sequence belongs in the answer iff it starts at two different indices. The most literal check: for every window, scan the whole string for another occurrence. The only subtlety is reporting each repeated *text* once — solved by letting only the **first** occurrence of a text do the reporting: if the same window text already appeared at an earlier index, skip it, because the earlier pass already decided it.

### Algorithm
1. For each start index `i` from `0` to `n-10`:
   1. Let `window = s[i:i+10]`.
   2. **Dedup check:** scan `j < i`; if `s[j:j+10] == window`, skip this `i` (already reported by its first occurrence).
   3. **Repeat check:** scan `j > i`; if `s[j:j+10] == window`, append `window` to the result and stop scanning.
2. Return the result (ordered by first occurrence).

### Complexity
- **Time:** O(n² · L) — each of the n windows is compared against up to n others, and each string comparison costs up to L = 10 character checks. At n = 10⁵ that is ~10¹¹ operations — hopeless in production, fine on the tiny examples.
- **Space:** O(1) extra — string slices in Go are views (pointer + length), so comparisons allocate nothing; only the output slice grows.

### Code
```go
func bruteForce(s string) []string {
	result := []string{} // answers, ordered by first occurrence
	n := len(s)
	for i := 0; i+seqLen <= n; i++ {
		window := s[i : i+seqLen] // candidate 10-letter sequence
		// (a) reported already? — check all earlier windows for the same text
		seenBefore := false
		for j := 0; j < i; j++ {
			if s[j:j+seqLen] == window {
				seenBefore = true // its first occurrence handled the reporting
				break
			}
		}
		if seenBefore {
			continue
		}
		// (b) does it repeat later? — one later match is enough to qualify
		for j := i + 1; j+seqLen <= n; j++ {
			if s[j:j+seqLen] == window {
				result = append(result, window) // first occurrence reports it
				break
			}
		}
	}
	return result
}
```

### Dry Run
Example 1: `s = "AAAAACCCCCAAAAACCCCCCAAAAAGGGTTT"` (n = 32, windows start at i = 0..22)

| i | window | earlier match? | later match? | result |
|---|--------|----------------|--------------|--------|
| 0 | `AAAAACCCCC` | — (nothing before) | yes, at j = 10 | `[AAAAACCCCC]` |
| 1 | `AAAACCCCCA` | no | no | unchanged |
| 2–4 | `AAACCCCCAA`, `AACCCCCAAA`, `ACCCCCAAAA` | no | no | unchanged |
| 5 | `CCCCCAAAAA` | no | yes, at j = 16 | `[AAAAACCCCC, CCCCCAAAAA]` |
| 6–9 | various | no | no | unchanged |
| 10 | `AAAAACCCCC` | **yes, at j = 0** → skip | — | unchanged (no duplicate) |
| 11–15 | various | no | no | unchanged |
| 16 | `CCCCCAAAAA` | **yes, at j = 5** → skip | — | unchanged |
| 17–22 | various (drift into `G`/`T` region) | no | no | unchanged |

Final: `["AAAAACCCCC", "CCCCCAAAAA"]` ✓

---

## Approach 2 — Hash Map Counting

### Intuition
One pass with a `substring → count` map answers "have I seen this window before?" in O(1). The neat trick is to append a window **exactly when its count becomes 2**: earlier than 2 it is not yet repeated, later than 2 it was already reported. This kills two birds — dedup and deterministic left-to-right detection order — without a second pass or an extra set.

### Algorithm
1. Slide `i` over every window start (`0` to `n-10`); let `sub = s[i:i+10]`.
2. Increment `counts[sub]`.
3. If `counts[sub] == 2`, append `sub` to the result.
4. Return the result.

### Complexity
- **Time:** O(n · L) — n window slides; hashing each 10-byte key costs L = 10. Effectively linear with a constant factor of 10.
- **Space:** O(n · L) — worst case (all windows distinct) the map stores one 10-byte key per window.

### Code
```go
func hashMap(s string) []string {
	result := []string{}          // sequences confirmed repeated, in detection order
	counts := map[string]int{}    // 10-letter window → occurrences so far
	for i := 0; i+seqLen <= len(s); i++ {
		sub := s[i : i+seqLen] // current window snapshot
		counts[sub]++
		if counts[sub] == 2 { // exactly the moment it becomes repeated
			result = append(result, sub) // == 2 (not >= 2) guarantees a single report
		}
	}
	return result
}
```

### Dry Run
Example 1: `s = "AAAAACCCCCAAAAACCCCCCAAAAAGGGTTT"`

| i | sub | counts[sub] after ++ | action | result |
|---|-----|----------------------|--------|--------|
| 0 | `AAAAACCCCC` | 1 | — | `[]` |
| 1..9 | 9 distinct windows | 1 each | — | `[]` |
| 10 | `AAAAACCCCC` | **2** | append | `[AAAAACCCCC]` |
| 11..15 | distinct windows | 1 each | — | unchanged |
| 16 | `CCCCCAAAAA` | **2** (first seen at i = 5) | append | `[AAAAACCCCC, CCCCCAAAAA]` |
| 17..22 | remaining windows | 1 each | — | unchanged |

Final: `["AAAAACCCCC", "CCCCCAAAAA"]` ✓

Example 2: `s = "AAAAAAAAAAAAA"` (13 A's, windows at i = 0..3, all `AAAAAAAAAA`): count hits 2 at i = 1 → appended once; counts 3 and 4 at i = 2, 3 do **not** re-append. Final: `["AAAAAAAAAA"]` ✓

---

## Approach 3 — Rolling Hash / Bit Manipulation (Optimal)

### Intuition
Two observations sharpen Approach 2:
1. **The alphabet has 4 letters** → each letter needs only 2 bits (`A=00, C=01, G=10, T=11`), so a 10-letter window fits in exactly 20 bits. That is a *perfect* hash — two windows share a fingerprint iff they are the same text, so collisions are impossible.
2. **Consecutive windows overlap in 9 letters** → instead of re-hashing 10 characters per step, shift the previous fingerprint left by 2 bits, OR in the new letter, and mask to 20 bits. Each slide is O(1) — the Rabin–Karp rolling-hash idea in its cleanest form.

Two integer sets replace the string-keyed counter: `seen` (fingerprints observed) and `added` (fingerprints already reported), giving exactly-once reporting.

### Algorithm
1. Map letters to 2-bit codes: `A→0, C→1, G→2, T→3`.
2. For each `i` from `0` to `n-1`: `window = ((window << 2) | code[s[i]]) & 0xFFFFF` (keep the low 20 bits = the newest 10 letters).
3. Once `i >= 9` the fingerprint covers a full window ending at `i`:
   - If `window ∈ seen` and `window ∉ added`: append `s[i-9:i+1]` and put `window` in `added`.
   - Put `window` in `seen`.
4. Return the result.

### Complexity
- **Time:** O(n) — a single pass with O(1) shift/mask/set work per character; the ×L re-hash factor of Approach 2 is gone.
- **Space:** O(n) — at most one 20-bit integer per window in each set (integer keys are far lighter than 10-byte string keys; at most 2²⁰ ≈ 10⁶ distinct fingerprints exist).

### Code
```go
func rollingHash(s string) []string {
	result := []string{}
	if len(s) < seqLen {
		return result // no 10-letter window even exists
	}
	// 2-bit code per nucleotide — 10 letters pack into 20 bits losslessly
	code := map[byte]uint32{'A': 0, 'C': 1, 'G': 2, 'T': 3}
	const mask = 1<<(2*seqLen) - 1 // 0xFFFFF: keep only the newest 10 letters
	var window uint32              // rolling 20-bit fingerprint of the last 10 letters
	seen := map[uint32]bool{}      // fingerprints observed at least once
	added := map[uint32]bool{}     // fingerprints already appended to result
	for i := 0; i < len(s); i++ {
		// push the new letter into the low bits, drop the letter that fell out
		window = (window<<2 | code[s[i]]) & mask
		if i >= seqLen-1 { // window spans a full 10 letters ending at i
			if seen[window] && !added[window] {
				result = append(result, s[i-seqLen+1:i+1]) // second sighting → report once
				added[window] = true                       // never report this text again
			}
			seen[window] = true // record this window for future sightings
		}
	}
	return result
}
```

### Dry Run
Example 1: `s = "AAAAACCCCCAAAAACCCCCCAAAAAGGGTTT"` (codes: A=00, C=01)

| i | letter | window text (= fingerprint decoded) | in `seen`? | action |
|---|--------|--------------------------------------|------------|--------|
| 0–8 | A×5, C×4 | still filling (< 10 letters) | — | build fingerprint only |
| 9 | C | letters 0..9 = `AAAAACCCCC` (bits `00000000000101010101`) | no | seen ← + key(`AAAAACCCCC`) |
| 10 | A | letters 1..10 = `AAAACCCCCA` | no | add to seen |
| 11–13 | A,A,A | `AAACCCCCAA`, `AACCCCCAAA`, `ACCCCCAAAA` | no | add to seen |
| 14 | A | letters 5..14 = `CCCCCAAAAA` | no | seen ← + key(`CCCCCAAAAA`) |
| 15–18 | C,C,C,C | `CCCCAAAAAC`, `CCCAAAAACC`, `CCAAAAACCC`, `CAAAAACCCC` | no | add to seen |
| 19 | C | letters 10..19 = `AAAAACCCCC` → same 20-bit key as i = 9 | **yes**, not added | append `AAAAACCCCC`; added ← + key |
| 20–24 | C,A,A,A,A | 5 new fingerprints | no | add to seen |
| 25 | A | letters 16..25 = `CCCCCAAAAA` → same key as i = 14 | **yes**, not added | append `CCCCCAAAAA`; added ← + key |
| 26–31 | G,G,G,T,T,T | new fingerprints containing G/T | no | add to seen |

Final: `["AAAAACCCCC", "CCCCCAAAAA"]` ✓ — each fingerprint's *second* sighting (and only the second, thanks to `added`) triggers a report.

---

## Key Takeaways

- **Small alphabet ⇒ bit-packed perfect hash.** With ≤ 4 symbols, 2 bits per symbol packs a length-10 window into 20 bits — membership tests become integer set lookups with zero collision risk.
- **Rolling update pattern:** `key = ((key << bits) | newCode) & mask` slides a fixed window in O(1); it is the integer form of the Rabin–Karp rolling hash.
- **"Append when count == 2"** is a one-liner that deduplicates *and* preserves deterministic detection order — no separate result set required.
- Fixed-length substring problems are sliding-window problems in disguise: the candidate space is exactly the n − L + 1 windows, so aim for O(1) work per slide.

---

## Related Problems

- LeetCode #28 — Find the Index of the First Occurrence in a String (Rabin–Karp rolling hash)
- LeetCode #1044 — Longest Duplicate Substring (rolling hash + binary search on length)
- LeetCode #3 — Longest Substring Without Repeating Characters (sliding window over a string)
- LeetCode #438 — Find All Anagrams in a String (fixed-size window + frequency fingerprint)
- LeetCode #2156 — Find Substring With Given Hash Value (explicit rolling-hash construction)
