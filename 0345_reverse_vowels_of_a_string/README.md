# 0345 — Reverse Vowels of a String

> LeetCode #345 · Difficulty: Easy
> **Categories:** Two Pointers, String

---

## Problem Statement

Given a string `s`, reverse only all the vowels in the string and return it.

The vowels are `'a'`, `'e'`, `'i'`, `'o'`, and `'u'`, and they can appear in both lower and upper cases, more than once.

**Example 1:**

```
Input: s = "IceCreAm"
Output: "AceCreIm"
Explanation: The vowels in s are ['I', 'e', 'e', 'A'].
On reversing the vowels, s becomes "AceCreIm".
```

**Example 2:**

```
Input: s = "leetcode"
Output: "leotcede"
```

**Constraints:**

- `1 <= s.length <= 3 * 10^5`
- `s` consist of printable ASCII characters.

---

## Company Frequency

| Company    | Frequency        | Last Reported |
|------------|------------------|---------------|
| Amazon     | ★★★☆☆ Medium     | 2024          |
| Google     | ★★☆☆☆ Low        | 2023          |
| Microsoft  | ★★☆☆☆ Low        | 2023          |
| Facebook   | ★★☆☆☆ Low        | 2022          |
| Apple      | ★★☆☆☆ Low        | 2022          |

> ⚠️ Frequency data is crowd-sourced from LeetCode Discuss, Glassdoor, and
> community interview reports. Treat as a signal, not a guarantee.

---

## DSA Concepts Used

- **Two Pointers** — one pointer from each end, each skipping non-vowels, swapping when both land on vowels; the classic selective-reversal pattern → see [`/dsa/two_pointers.md`](/dsa/two_pointers.md)
- **String** — in-place character manipulation over a byte slice → see [`/dsa/string_algorithms.md`](/dsa/string_algorithms.md)

---

## Approaches Overview

| # | Approach | Time | Space | When to use |
|---|----------|------|-------|-------------|
| 1 | Two Pointers (Optimal) | O(n) | O(n) buffer | The intended answer; single pass, minimal extra memory |
| 2 | Collect Indices then Reverse | O(n) | O(V) | Clear two-pass alternative; easy to reason about |

---

## Approach 1 — Two Pointers (Optimal)

### Intuition
Reversing "only the vowels" means the vowel subsequence is reversed while consonants stay put. Use a pointer from each end: advance each past non-vowels; when both point at vowels, swap them and step inward. Consonants are simply skipped and never moved.

### Algorithm
1. Convert to `[]byte`; `left = 0`, `right = len-1`.
2. Advance `left` while `s[left]` is not a vowel.
3. Advance `right` (leftward) while `s[right]` is not a vowel.
4. If `left < right`, swap the two vowels; step both inward.
5. Repeat until `left >= right`.

### Complexity
- **Time:** O(n) — each index is examined at most once.
- **Space:** O(n) for the mutable byte buffer (Go strings are immutable); O(1) extra beyond it.

### Code
```go
func twoPointers(s string) string {
	b := []byte(s)
	left, right := 0, len(b)-1
	for left < right {
		for left < right && !isVowel(b[left]) {
			left++
		}
		for left < right && !isVowel(b[right]) {
			right--
		}
		if left < right {
			b[left], b[right] = b[right], b[left]
			left++
			right--
		}
	}
	return string(b)
}
```

### Dry Run
Input `s = "IceCreAm"` → bytes `I c e C r e A m` (indices 0..7):

| Step | left | right | b[left] | b[right] | Action | string |
|------|------|-------|---------|----------|--------|--------|
| 0 | 0 | 7 | I (vowel) | m (skip) → right=6 | A (vowel); swap I↔A | `A c e C r e I m` |
| 1 | 1 | 5 | c (skip)→2 | e (vowel) | b[2]=e vowel; swap e↔e | `A c e C r e I m` |
| 2 | 3 | 4 | C (skip)→4 | r (skip)→... | left=4,right=4 → stop | `A c e C r e I m` |

Output `"AceCreIm"`.

---

## Approach 2 — Collect Indices then Reverse

### Intuition
First find *where* the vowels are (their indices, left to right). The vowels must end up in reverse order at those same slots, so pair the k-th index from the front with the k-th index from the back and swap.

### Algorithm
1. Pass 1: collect all indices `idx` where `b[i]` is a vowel.
2. Pass 2: with `i` from the front of `idx` and `j` from the back, swap `b[idx[i]]` and `b[idx[j]]`, moving inward.

### Complexity
- **Time:** O(n) — two linear passes.
- **Space:** O(V) — the list of vowel indices (V = number of vowels).

### Code
```go
func collectIndices(s string) string {
	b := []byte(s)
	idx := []int{}
	for i := 0; i < len(b); i++ {
		if isVowel(b[i]) {
			idx = append(idx, i)
		}
	}
	for i, j := 0, len(idx)-1; i < j; i, j = i+1, j-1 {
		b[idx[i]], b[idx[j]] = b[idx[j]], b[idx[i]]
	}
	return string(b)
}
```

### Dry Run
Input `s = "IceCreAm"`:

| Phase | State |
|-------|-------|
| Pass 1 | vowels at indices `idx = [0, 2, 5, 6]` (I,e,e,A) |
| Pass 2, i=0 j=3 | swap b[0]↔b[6] (I↔A) → `A c e C r e I m` |
| Pass 2, i=1 j=2 | swap b[2]↔b[5] (e↔e) → unchanged |
| i=2 j=1 | i >= j → stop |

Output `"AceCreIm"`.

---

## Key Takeaways

- This is the two-pointer template plus a **skip condition**: only act when both pointers satisfy a predicate (here, "is a vowel").
- A helper like `isVowel` should be case-insensitive — vowels appear in both cases.
- The swap moves the exact characters (preserving their original case), so `I`/`A` keep their casing after swapping positions.
- The two-pass "collect indices then swap symmetric pairs" formulation generalises to reversing any filtered subsequence in place.

---

## Related Problems

- LeetCode #344 — Reverse String (unconditional two-pointer reversal)
- LeetCode #125 — Valid Palindrome (two pointers with skip condition)
- LeetCode #917 — Reverse Only Letters (same pattern, letters instead of vowels)
- LeetCode #541 — Reverse String II (block reversal)
