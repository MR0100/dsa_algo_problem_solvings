# 0266 — Palindrome Permutation

> LeetCode #266 · Difficulty: Easy
> **Categories:** Hash Table, String, Bit Manipulation

---

## Problem Statement

Given a string `s`, return `true` if a permutation of the string could form a
palindrome and `false` otherwise.

**Example 1:**
```
Input: s = "code"
Output: false
```

**Example 2:**
```
Input: s = "aab"
Output: true
```
Explanation: "aba" is a palindrome permutation of "aab".

**Example 3:**
```
Input: s = "carerac"
Output: true
```
Explanation: "racecar" is a palindrome permutation of "carerac".

**Constraints:**
- `1 <= s.length <= 5000`
- `s` consists of only lowercase English letters.

---

## Company Frequency

| Company    | Frequency        | Last Reported |
|------------|------------------|---------------|
| Amazon     | ★★★★☆ High       | 2024          |
| Google     | ★★★☆☆ Medium     | 2023          |
| Bloomberg  | ★★★☆☆ Medium     | 2023          |
| Microsoft  | ★★☆☆☆ Low        | 2022          |

> ⚠️ Frequency data is crowd-sourced from LeetCode Discuss, Glassdoor, and
> community interview reports. Treat as a signal, not a guarantee.

---

## DSA Concepts Used
- **Hash Map / Frequency Count** — parity of character counts decides palindrome feasibility → see [`/dsa/hash_map.md`](/dsa/hash_map.md)
- **String Algorithms** — reasoning about palindrome structure → see [`/dsa/string_algorithms.md`](/dsa/string_algorithms.md)
- **Bit / Set Parity Trick** — toggling set membership tracks odd/even counts without storing them → see [`/dsa/bit_manipulation.md`](/dsa/bit_manipulation.md)

---

## Approaches Overview

| # | Approach | Time | Space | When to use |
|---|----------|------|-------|-------------|
| 1 | Hash Map Count | O(n) | O(k) | Clear, explicit counting |
| 2 | Single-Pass Set Toggle (Optimal) | O(n) | O(k) | Only parity matters — skip the counts |

---

## Approach 1 — Hash Map Count

### Intuition
A string can be rearranged into a palindrome iff **at most one** distinct
character has an odd count. Even counts pair symmetrically around the centre; a
single odd-count character can occupy the middle position.

### Algorithm
1. Count how many times each character appears.
2. Count how many characters have an odd frequency.
3. Return `true` if the number of odd-count characters is 0 or 1.

### Complexity
- **Time:** O(n) — one pass to count, one pass over the distinct chars.
- **Space:** O(k) — k distinct characters (≤ 26 here, so effectively O(1)).

### Code
```go
func hashMap(s string) bool {
	counts := make(map[rune]int) // char -> frequency
	for _, c := range s {
		counts[c]++ // tally every character
	}
	odd := 0 // number of characters seen an odd number of times
	for _, v := range counts {
		if v%2 == 1 { // odd frequency
			odd++
		}
	}
	return odd <= 1 // palindrome possible with at most one odd-count char
}
```

### Dry Run
Input `s = "code"`:

| Step | Char | counts state           |
|------|------|------------------------|
| 1    | c    | {c:1}                  |
| 2    | o    | {c:1, o:1}             |
| 3    | d    | {c:1, o:1, d:1}        |
| 4    | e    | {c:1, o:1, d:1, e:1}   |

Odd-count characters: c, o, d, e → `odd = 4`. `4 <= 1` is false → return **false**. ✅

---

## Approach 2 — Single-Pass Set Toggle (Optimal)

### Intuition
We do not need the exact counts, only their **parity**. Keep a set of characters
currently seen an odd number of times: insert on first sight, remove on the
second, insert again on the third, and so on. After the scan the set holds
exactly the odd-count characters; a palindrome is possible iff its size ≤ 1.

### Algorithm
1. For each char: if it is already in the set, delete it; otherwise insert it.
2. After the pass, return `len(set) <= 1`.

### Complexity
- **Time:** O(n) — single pass over the string.
- **Space:** O(k) — the set of distinct odd-parity characters.

### Code
```go
func setToggle(s string) bool {
	seen := make(map[rune]struct{}) // chars currently at odd parity
	for _, c := range s {
		if _, ok := seen[c]; ok { // second (even) occurrence
			delete(seen, c) // parity flips back to even
		} else {
			seen[c] = struct{}{} // odd occurrence
		}
	}
	return len(seen) <= 1 // at most one leftover odd char
}
```

### Dry Run
Input `s = "code"`:

| Step | Char | Action              | seen set        |
|------|------|---------------------|-----------------|
| 1    | c    | insert (new)        | {c}             |
| 2    | o    | insert (new)        | {c, o}          |
| 3    | d    | insert (new)        | {c, o, d}       |
| 4    | e    | insert (new)        | {c, o, d, e}    |

`len(seen) = 4`, `4 <= 1` is false → return **false**. ✅

---

## Key Takeaways
- Palindrome feasibility reduces to **counting odd-frequency characters** — at most one is allowed.
- When only parity matters, a **toggle set** replaces a full frequency map and removes the second scan.
- The same parity idea appears with XOR bitmasks (each of 26 letters = one bit) for O(1) extra space.

---

## Related Problems
- LeetCode #267 — Palindrome Permutation II (build the actual palindromes)
- LeetCode #409 — Longest Palindrome (same odd-count reasoning, returns a length)
- LeetCode #5 — Longest Palindromic Substring (palindrome structure)
- LeetCode #125 — Valid Palindrome (palindrome checking)
