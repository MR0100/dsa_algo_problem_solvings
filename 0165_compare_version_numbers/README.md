# 0165 — Compare Version Numbers

> LeetCode #165 · Difficulty: Medium
> **Categories:** Two Pointers, String

---

## Problem Statement

Given two **version strings**, `version1` and `version2`, compare them. A version string consists of **revisions** separated by dots `'.'`. The **value of the revision** is its **integer conversion** ignoring leading zeros.

To compare version strings, compare their revision values in **left-to-right order**. If one of the version strings has fewer revisions, treat the missing revision values as `0`.

Return the following:

- If `version1 < version2`, return `-1`.
- If `version1 > version2`, return `1`.
- Otherwise, return `0`.

**Example 1:**
```
Input: version1 = "1.2", version2 = "1.10"
Output: -1
Explanation:
version1's second revision is "2" and version2's second revision is "10": 2 < 10, so version1 < version2.
```

**Example 2:**
```
Input: version1 = "1.01", version2 = "1.001"
Output: 0
Explanation:
Ignoring leading zeroes, both "01" and "001" represent the same integer "1".
```

**Example 3:**
```
Input: version1 = "1.0", version2 = "1.0.0.0"
Output: 0
Explanation:
version1 has less revisions, which means every missing revision are treated as "0".
```

**Constraints:**
- `1 <= version1.length, version2.length <= 500`
- `version1` and `version2` only contain digits and `'.'`.
- `version1` and `version2` are valid version numbers.
- All the given revisions in `version1` and `version2` can be stored in a **32-bit integer**.

---

## Company Frequency

| Company   | Frequency       | Last Reported |
|-----------|-----------------|---------------|
| Apple     | ★★★★★ Very High | 2024          |
| Amazon    | ★★★★☆ High      | 2024          |
| Microsoft | ★★★★☆ High      | 2024          |
| Meta      | ★★★☆☆ Medium    | 2023          |
| Bloomberg | ★★★☆☆ Medium    | 2023          |
| Zoom      | ★★★☆☆ Medium    | 2023          |

> ⚠️ Frequency data is crowd-sourced from LeetCode Discuss, Glassdoor, and community interview reports. Treat as a signal, not a guarantee.

---

## DSA Concepts Used

- **Two Pointers** — one independent cursor per string parses revision-by-revision in lock-step, padding the exhausted side with zeros → see [`/dsa/two_pointers.md`](/dsa/two_pointers.md)
- **String Algorithms** — tokenising on a delimiter and manual digit-accumulation parsing (the `r = r*10 + digit` idiom from atoi) → see [`/dsa/string_algorithms.md`](/dsa/string_algorithms.md)

---

## Approaches Overview

| # | Approach                              | Time     | Space    | When to use                                                     |
|---|---------------------------------------|----------|----------|------------------------------------------------------------------|
| 1 | Brute Force (Split + Parse + Compare) | O(m + n) | O(m + n) | Production code — clearest to read and review                   |
| 2 | Two Pointers In-Place Parse (Optimal) | O(m + n) | O(1)     | Interviews / memory-tight paths — no substring allocations       |

---

## Approach 1 — Brute Force (Split + Parse + Compare)

### Intuition
The statement *is* the algorithm: revisions are the dot-separated chunks, a revision's value is its integer conversion (which throws away leading zeros automatically), and the shorter version behaves as if padded with zero revisions. So split both strings on `'.'`, convert chunk by chunk, and compare position by position up to the *longer* list — the first inequality decides everything.

### Algorithm
1. `parts1 = Split(version1, ".")`, `parts2 = Split(version2, ".")`.
2. `longer = max(len(parts1), len(parts2))`.
3. For `i` from `0` to `longer − 1`:
   1. `r1 = Atoi(parts1[i])` if `i < len(parts1)`, else `0`.
   2. `r2 = Atoi(parts2[i])` if `i < len(parts2)`, else `0`.
   3. If `r1 < r2`, return `-1`; if `r1 > r2`, return `1`.
4. Every position tied → return `0`.

### Complexity
- **Time:** O(m + n) — splitting, parsing, and comparing visit each character of both strings a constant number of times.
- **Space:** O(m + n) — `Split` materialises every revision as a substring slice.

### Code
```go
func splitAndCompare(version1, version2 string) int {
	parts1 := strings.Split(version1, ".") // revision chunks of version1
	parts2 := strings.Split(version2, ".") // revision chunks of version2
	longer := len(parts1)
	if len(parts2) > longer {
		longer = len(parts2) // compare up to the longer revision list
	}
	for i := 0; i < longer; i++ {
		r1, r2 := 0, 0 // missing revisions count as 0
		if i < len(parts1) {
			r1, _ = strconv.Atoi(parts1[i]) // Atoi drops leading zeros ("001" → 1)
		}
		if i < len(parts2) {
			r2, _ = strconv.Atoi(parts2[i])
		}
		if r1 < r2 { // first differing revision decides the order
			return -1
		}
		if r1 > r2 {
			return 1
		}
	}
	return 0 // every revision matched → the versions are equal
}
```

### Dry Run
`version1 = "1.2"`, `version2 = "1.10"` (Example 1). After splitting: `parts1 = ["1","2"]`, `parts2 = ["1","10"]`, `longer = 2`:

| i | parts1[i] | parts2[i] | r1 | r2 | comparison    | action        |
|---|-----------|-----------|----|----|---------------|---------------|
| 0 | "1"       | "1"       | 1  | 1  | equal         | next revision |
| 1 | "2"       | "10"      | 2  | 10 | `2 < 10`      | return **-1** ✅ |

(Example 3 `"1.0"` vs `"1.0.0.0"`: positions 2 and 3 read `r1 = 0` from the exhausted `parts1` versus `r2 = Atoi("0") = 0` — all tie → **0** ✅.)

---

## Approach 2 — Two Pointers In-Place Parse (Optimal)

### Intuition
We never need both full revision lists in memory — only the *current pair* of revision values. Keep one cursor per string; each round, accumulate digits (`r = r*10 + digit`) until hitting a `'.'` or the end. A cursor that has run off its string simply contributes `0`, which implements the "missing revisions are 0" rule with zero extra code. Leading zeros vanish for free because accumulating `0` first multiplies nothing in (`0*10 + d = d`). No substrings, no allocations — O(1) auxiliary space.

### Algorithm
1. `i = 0` (cursor into `version1`), `j = 0` (cursor into `version2`).
2. While `i < len(version1)` **or** `j < len(version2)`:
   1. `r1 = 0`; while `i` points at a digit, `r1 = r1*10 + digit(version1[i])`, advance `i`.
   2. `r2 = 0`; while `j` points at a digit, `r2 = r2*10 + digit(version2[j])`, advance `j`.
   3. If `r1 < r2`, return `-1`; if `r1 > r2`, return `1`.
   4. Advance both cursors one step past their `'.'` (a past-the-end increment on an exhausted string is harmless — the cursor just stays out of range).
3. Both strings consumed with all revisions equal → return `0`.

### Complexity
- **Time:** O(m + n) — each cursor moves strictly forward, so every character is inspected exactly once.
- **Space:** O(1) — two cursors and two integer accumulators; no substring slices.

### Code
```go
func twoPointers(version1, version2 string) int {
	i, j := 0, 0 // cursors into version1 and version2
	for i < len(version1) || j < len(version2) {
		r1 := 0
		for i < len(version1) && version1[i] != '.' {
			r1 = r1*10 + int(version1[i]-'0') // accumulate digit (leading zeros vanish)
			i++
		}
		r2 := 0
		for j < len(version2) && version2[j] != '.' {
			r2 = r2*10 + int(version2[j]-'0')
			j++
		}
		if r1 < r2 { // first differing revision decides the order
			return -1
		}
		if r1 > r2 {
			return 1
		}
		i++ // skip the '.' (harmless past-the-end increment when exhausted)
		j++
	}
	return 0 // both strings fully consumed with equal revisions
}
```

### Dry Run
`version1 = "1.2"`, `version2 = "1.10"` (Example 1):

| round | i before | j before | r1 parse                    | r2 parse                          | r1 | r2 | result / cursors after     |
|-------|----------|----------|------------------------------|-----------------------------------|----|----|-----------------------------|
| 1     | 0        | 0        | reads '1', stops at '.' (i=1) | reads '1', stops at '.' (j=1)     | 1  | 1  | tie → skip dots: i=2, j=2   |
| 2     | 2        | 2        | reads '2', hits end (i=3)     | reads '1' then '0' → 10 (j=4)     | 2  | 10 | `2 < 10` → return **-1** ✅ |

(Example 2 `"1.01"` vs `"1.001"`: round 2 accumulates `0*10+0=0`, then `0*10+1=1` on one side and `0→0→1` on the other — both give `r=1`, so leading zeros never matter → **0** ✅.)

---

## Key Takeaways

- **Version strings are not decimals** — `"1.10" > "1.2"` even though 1.10 < 1.2 as floats. Never compare version strings lexicographically or numerically as a whole; compare revision by revision.
- The digit-accumulator idiom `r = r*10 + int(c − '0')` (from #8 atoi) parses integers in place with no allocations and neutralises leading zeros as a side effect.
- **Virtual zero-padding**: letting an exhausted cursor yield `0` handles unequal revision counts (`"1.0"` vs `"1.0.0.0"`) with no special-case code — a pattern that also appears when adding numbers stored as lists/strings (#2, #415).
- Loop `while i or j has input` (OR, not AND) whenever two sequences of different lengths must be consumed to the very end.
- The constraint "each revision fits in a 32-bit integer" is what makes plain `int` accumulation safe; if revisions could be arbitrarily long you would compare chunk lengths after stripping leading zeros, then compare lexicographically.

---

## Related Problems

- LeetCode #8 — String to Integer (atoi) (the same manual digit-parsing core)
- LeetCode #415 — Add Strings (parallel cursors with virtual padding over unequal lengths)
- LeetCode #2 — Add Two Numbers (unequal-length lock-step consumption with implicit zeros)
- LeetCode #43 — Multiply Strings (digit-by-digit string arithmetic)
- LeetCode #468 — Validate IP Address (dot-delimited token parsing and validation)
