# 0168 — Excel Sheet Column Title

> LeetCode #168 · Difficulty: Easy
> **Categories:** Math, String

---

## Problem Statement

Given an integer `columnNumber`, return *its corresponding column title as it appears in an Excel sheet*.

For example:

```
A -> 1
B -> 2
C -> 3
...
Z -> 26
AA -> 27
AB -> 28
...
```

**Example 1:**

```
Input: columnNumber = 1
Output: "A"
```

**Example 2:**

```
Input: columnNumber = 28
Output: "AB"
```

**Example 3:**

```
Input: columnNumber = 701
Output: "ZY"
```

**Constraints:**

- `1 <= columnNumber <= 2^31 - 1`

---

## Company Frequency

| Company    | Frequency        | Last Reported |
|------------|------------------|---------------|
| Microsoft  | ★★★★☆ High       | 2024          |
| Amazon     | ★★★☆☆ Medium     | 2024          |
| Google     | ★★★☆☆ Medium     | 2023          |
| Apple      | ★★☆☆☆ Low        | 2023          |
| Zenefits   | ★★☆☆☆ Low        | 2022          |

> ⚠️ Frequency data is crowd-sourced from LeetCode Discuss, Glassdoor, and
> community interview reports. Treat as a signal, not a guarantee.

---

## DSA Concepts Used

- **Math / Number Theory (base conversion)** — this is *bijective* base-26: digits run 1..26 with no zero, fixed by subtracting 1 before every `mod`/`div` → see [`/dsa/math_number_theory.md`](/dsa/math_number_theory.md)
- **String Algorithms** — digits emerge least-significant first, so the answer is built backwards and reversed (or prepended via recursion) → see [`/dsa/string_algorithms.md`](/dsa/string_algorithms.md)

---

## Approaches Overview

| # | Approach | Time | Space | When to use |
|---|----------|------|-------|-------------|
| 1 | Brute Force (length-block search) | O(k²), k = log₂₆(n) ≤ 7 | O(k) | To *derive* the numbering system from first principles |
| 2 | Iterative Base-26 with Offset (Optimal) | O(k) | O(k) | Always — the canonical subtract-one trick |
| 3 | Recursive Base-26 | O(k) | O(k) stack | Same math, cleaner code; shows recursion replacing the reverse step |

---

## Approach 1 — Brute Force (Digit-by-Digit Search)

### Intuition

Instead of spotting the base-26 trick, reason mechanically about **how many titles exist of each length**: 26 titles of length 1 (`A`..`Z`), 26² of length 2 (`AA`..`ZZ`), 26³ of length 3, and so on. Subtracting these block sizes locates the answer's length `k` and its 0-based offset inside the length-`k` block. Within that block titles are ordered lexicographically, so the first letter splits the block into 26 equal chunks of size 26^(k−1) — plain division picks the chunk (letter), and we recurse into it with `offset % chunk`.

### Algorithm

1. Set `offset = columnNumber`, `blockSize = 26`, `k = 1`.
2. While `offset > blockSize`: `offset -= blockSize`, `blockSize *= 26`, `k++`. Now the title has `k` letters.
3. Decrement `offset` to make it 0-based within the length-`k` block.
4. Compute `chunk = 26^(k-1)`.
5. For each of the `k` positions left to right: letter = `'A' + offset/chunk`; then `offset %= chunk`, `chunk /= 26`.

### Complexity

- **Time:** O(k²) with k = ⌈log₂₆(n)⌉ ≤ 7 — the power computation plus k divisions; effectively constant.
- **Space:** O(k) — the output buffer.

### Code

```go
func bruteForce(columnNumber int) string {
	// Step 1: find the title length k and the 0-based offset inside that block.
	offset := columnNumber
	blockSize := 26 // number of titles having the current length (26^k)
	k := 1
	for offset > blockSize {
		offset -= blockSize // skip all titles of length k
		blockSize *= 26     // next block: titles of length k+1
		k++
	}
	offset-- // make the offset 0-based within the length-k block

	// Step 2: pick letters left to right by dividing into 26 equal chunks.
	out := make([]byte, k)
	// chunk = 26^(k-1): how many titles share the same first letter.
	chunk := 1
	for i := 0; i < k-1; i++ {
		chunk *= 26
	}
	for i := 0; i < k; i++ {
		out[i] = byte('A' + offset/chunk) // which of the 26 chunks we are in
		offset %= chunk                   // descend into that chunk
		if chunk > 1 {
			chunk /= 26 // next position partitions 26x finer
		}
	}
	return string(out)
}
```

### Dry Run

Example 1: `columnNumber = 1`.

| Step | Action | offset | blockSize | k | chunk | out |
|------|--------|--------|-----------|---|-------|-----|
| 1 | init | 1 | 26 | 1 | — | [] |
| 2 | `offset (1) > blockSize (26)`? no → length found | 1 | 26 | 1 | — | [] |
| 3 | 0-base: `offset--` | 0 | 26 | 1 | — | [] |
| 4 | `chunk = 26^0 = 1` | 0 | — | 1 | 1 | [] |
| 5 | position 0: letter `'A' + 0/1 = 'A'`; `offset = 0%1 = 0` | 0 | — | 1 | 1 | ['A'] |

Result: `"A"` ✔

Bonus trace of Example 2 (`columnNumber = 28`):

| Step | Action | offset | blockSize | k | chunk | out |
|------|--------|--------|-----------|---|-------|-----|
| 1 | init | 28 | 26 | 1 | — | [] |
| 2 | `28 > 26` → skip length-1 block | 2 | 676 | 2 | — | [] |
| 3 | `2 > 676`? no; 0-base: `offset--` | 1 | — | 2 | — | [] |
| 4 | `chunk = 26^1 = 26` | 1 | — | 2 | 26 | [] |
| 5 | position 0: `'A' + 1/26 = 'A'`; `offset = 1%26 = 1`; `chunk = 1` | 1 | — | 2 | 1 | ['A'] |
| 6 | position 1: `'A' + 1/1 = 'B'`; `offset = 0` | 0 | — | 2 | 1 | ['A','B'] |

Result: `"AB"` ✔

---

## Approach 2 — Iterative Base-26 with Offset (Optimal)

### Intuition

The column system is **bijective base-26**: the "digits" are A=1 … Z=26 and there is *no zero*. Ordinary base conversion (`n % 26`, `n / 26`) breaks because a remainder of 0 should mean `Z` (26), not a new digit. The classic fix: **subtract 1 before every mod/div step**. That maps the digit range 1..26 onto 0..25, where `% 26` and `/ 26` behave. Each step yields the *last* letter of the remaining title, so collect the letters and reverse at the end.

### Algorithm

1. While `columnNumber > 0`:
   1. `columnNumber--` — shift the 1..26 digit down to 0..25.
   2. Append `'A' + columnNumber % 26` to the buffer (this is the current least-significant letter).
   3. `columnNumber /= 26` — strip the emitted digit.
2. Reverse the buffer and return it.

### Complexity

- **Time:** O(log₂₆(n)) — one loop iteration per output letter; at most 7 for `2^31 − 1` (`"FXSHRXW"`).
- **Space:** O(log₂₆(n)) — the byte buffer holding the answer.

### Code

```go
func iterativeBase26(columnNumber int) string {
	out := []byte{}
	for columnNumber > 0 {
		columnNumber-- // shift digits from 1..26 to 0..25 (no zero in this system)
		// Least-significant "digit" is the last letter of the title.
		out = append(out, byte('A'+columnNumber%26))
		columnNumber /= 26 // drop the digit we just emitted
	}
	// Digits were produced right-to-left → reverse in place.
	for i, j := 0, len(out)-1; i < j; i, j = i+1, j-1 {
		out[i], out[j] = out[j], out[i]
	}
	return string(out)
}
```

### Dry Run

Example 1: `columnNumber = 1`.

| Step | columnNumber (loop entry) | after `--` | `% 26` → letter | after `/ 26` | out |
|------|---------------------------|------------|-----------------|--------------|-----|
| 1 | 1 | 0 | 0 → `'A'` | 0 | ['A'] |
| 2 | 0 → loop exits | — | — | — | ['A'] |
| 3 | reverse (single byte, unchanged) | — | — | — | `"A"` |

Result: `"A"` ✔

Bonus trace of Example 3 (`columnNumber = 701`):

| Step | columnNumber (loop entry) | after `--` | `% 26` → letter | after `/ 26` | out |
|------|---------------------------|------------|-----------------|--------------|-----|
| 1 | 701 | 700 | 700 % 26 = 24 → `'Y'` | 26 | ['Y'] |
| 2 | 26 | 25 | 25 % 26 = 25 → `'Z'` | 0 | ['Y','Z'] |
| 3 | 0 → exit; reverse | — | — | — | `"ZY"` |

Result: `"ZY"` ✔ (note step 2: without the `--`, 26 % 26 = 0 would wrongly emit `'A'`).

---

## Approach 3 — Recursive Base-26

### Intuition

Identical math, expressed as a recurrence: `title(n) = title((n-1)/26) + letter((n-1)%26)`. The recursion unwinds most-significant letter first, so the string comes out already in the right order — the call stack replaces the explicit reverse.

### Algorithm

1. Base case: `columnNumber == 0` → return `""`.
2. `columnNumber--` (the same bijective-base shift).
3. Return `recursiveBase26(columnNumber/26)` concatenated with the letter `'A' + columnNumber%26`.

### Complexity

- **Time:** O(log₂₆(n)) — one recursive call per letter (string concatenation cost is negligible at ≤ 7 letters).
- **Space:** O(log₂₆(n)) — recursion depth, at most 7 frames.

### Code

```go
func recursiveBase26(columnNumber int) string {
	// Base case: nothing left to convert.
	if columnNumber == 0 {
		return ""
	}
	columnNumber-- // bijective base-26: shift 1..26 → 0..25 before splitting
	// Prefix letters first (recursion), then this position's letter.
	return recursiveBase26(columnNumber/26) + string(byte('A'+columnNumber%26))
}
```

### Dry Run

Example 1: `columnNumber = 1`.

| Step | Call | after `--` | Recurse on | Letter appended | Returns |
|------|------|------------|-----------|-----------------|---------|
| 1 | `recursiveBase26(1)` | 0 | `recursiveBase26(0)` | `'A' + 0%26 = 'A'` | pending |
| 2 | `recursiveBase26(0)` | — (base case) | — | — | `""` |
| 3 | unwind step 1 | — | — | — | `"" + "A"` = `"A"` |

Result: `"A"` ✔

---

## Key Takeaways

- **Recognise bijective numeration:** digit set 1..k with no zero ⇒ subtract 1 before every `mod`/`div`. The single `columnNumber--` is the entire difficulty of this problem.
- The tell-tale failure mode: multiples of 26 (`Z`, `AZ`, `ZZ`…). If your code turns 26 into `"A0"`-ish garbage or 52 into `"BZ"` vs `"AZ"`, the missing `-1` is the bug.
- This problem is the exact **inverse of LeetCode #171** (Excel Sheet Column Number); implementing both back-to-back cements the mapping.
- Base-conversion loops naturally emit least-significant digits first — either reverse at the end (iterative) or let recursion order the output for you.

---

## Related Problems

- LeetCode #171 — Excel Sheet Column Number (the inverse mapping, title → number)
- LeetCode #7 — Reverse Integer (digit-by-digit mod/div processing)
- LeetCode #12 — Integer to Roman (value → positional string system)
- LeetCode #504 — Base 7 (ordinary base conversion for contrast)
