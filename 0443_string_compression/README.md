# 0443 — String Compression

> LeetCode #443 · Difficulty: Medium
> **Categories:** Two Pointers, String

---

## Problem Statement

Given an array of characters `chars`, compress it using the following algorithm:

Begin with an empty string `s`. For each group of **consecutive repeating characters** in `chars`:

- If the group's length is `1`, append the character to `s`.
- Otherwise, append the character followed by the group's length.

The compressed string `s` **should not be returned separately**, but instead, be stored **in the input character array `chars`**. Note that group lengths that are `10` or longer will be split into multiple characters in `chars`.

After you are done **modifying the input array**, return *the new length of the array*.

You must write an algorithm that uses only constant extra space.

**Example 1:**

```
Input: chars = ["a","a","b","b","c","c","c"]
Output: Return 6, and the first 6 characters of the input array should be: ["a","2","b","2","c","3"]
Explanation: The groups are "aa", "bb", and "ccc". This compresses to "a2b2c3".
```

**Example 2:**

```
Input: chars = ["a"]
Output: Return 1, and the first character of the input array should be: ["a"]
Explanation: The only group is "a", which remains uncompressed since it's a single character.
```

**Example 3:**

```
Input: chars = ["a","b","b","b","b","b","b","b","b","b","b","b","b"]
Output: Return 4, and the first 4 characters of the input array should be: ["a","b","1","2"].
Explanation: The groups are "a" and "bbbbbbbbbbbb". This compresses to "ab12".
```

**Constraints:**

- `1 <= chars.length <= 2000`
- `chars[i]` is a lowercase English letter, uppercase English letter, digit, or symbol.

---

## Company Frequency

| Company    | Frequency        | Last Reported |
|------------|------------------|---------------|
| Amazon     | ★★★★★ Very High  | 2024          |
| Microsoft  | ★★★★☆ High       | 2024          |
| Google     | ★★★☆☆ Medium     | 2023          |
| Meta       | ★★★☆☆ Medium     | 2023          |
| Apple      | ★★☆☆☆ Low        | 2023          |
| Bloomberg  | ★★☆☆☆ Low        | 2022          |

> ⚠️ Frequency data is crowd-sourced from LeetCode Discuss, Glassdoor, and
> community interview reports. Treat as a signal, not a guarantee.

---

## DSA Concepts Used

- **Two Pointers (read/write in place)** — a `read` pointer groups each run while a lagging `write` pointer emits the compressed output into the same array; the write head never overtakes the read head → see [`/dsa/two_pointers.md`](/dsa/two_pointers.md)
- **Run-Length Encoding on strings** — grouping consecutive equal characters and emitting `char + count` is the core string-scan pattern here → see [`/dsa/string_algorithms.md`](/dsa/string_algorithms.md)

---

## Approaches Overview

| # | Approach | Time | Space | When to use |
|---|----------|------|-------|-------------|
| 1 | Extra Buffer | O(n) | O(n) | Easiest to reason about; fails the O(1)-space rule |
| 2 | Two Pointers In-Place (Optimal) | O(n) | O(1) | The intended answer — compresses inside `chars` |

---

## Approach 1 — Extra Buffer (Build Then Copy Back)

### Intuition

Compression is run-length encoding. Build the encoded form in a scratch buffer — for each maximal run, push the character, and if the run is longer than one, push the digits of its length. Then copy the buffer back into `chars` so the result lives in the input as required. Simple and correct, at the cost of O(n) extra space.

### Algorithm

1. Scan runs: for a run starting at `i`, advance `j` while `chars[j] == chars[i]`.
2. Append `chars[i]` to the buffer; if `count = j - i > 1`, append each decimal digit of `count`.
3. Set `i = j` and continue. Finally `copy` the buffer into `chars` and return its length.

### Complexity

- **Time:** O(n) — one linear scan of the runs plus an O(n) copy back.
- **Space:** O(n) — the temporary buffer; this does **not** satisfy the constant-space follow-up.

### Code

```go
func extraBuffer(chars []byte) int {
	out := make([]byte, 0, len(chars)) // compressed result
	i := 0
	for i < len(chars) {
		ch := chars[i] // the character of this run
		j := i
		// Extend the run while the same character repeats.
		for j < len(chars) && chars[j] == ch {
			j++
		}
		count := j - i         // length of the run
		out = append(out, ch)  // always emit the character
		if count > 1 {
			// Emit the count's decimal digits (may be multiple, e.g. "12").
			out = append(out, []byte(strconv.Itoa(count))...)
		}
		i = j // jump to the next run
	}
	copy(chars, out) // write compressed content back into chars in place
	return len(out)  // new logical length
}
```

### Dry Run

Example 1: `chars = ["a","a","b","b","c","c","c"]`.

| Run start i | ch | run end j | count | Appended to buffer | buffer so far |
|-------------|----|-----------|-------|--------------------|---------------|
| 0 | a | 2 | 2 | `a`, `2` | `a2` |
| 2 | b | 4 | 2 | `b`, `2` | `a2b2` |
| 4 | c | 7 | 3 | `c`, `3` | `a2b2c3` |

Copy `a2b2c3` back into `chars`; return length `6` ✔.

---

## Approach 2 — Two Pointers In-Place (Optimal)

### Intuition

The key safety observation: any run of length `L` compresses to at most `L` characters (1 for the letter, plus `≤ L-1` digits — and for `L ≥ 2`, `1 + ⌈log₁₀L⌉ ≤ L`). So the **write** pointer can never pass the **read** pointer; overwriting in place is always safe. `read` groups the current run and counts it; `write` emits the character and, when `count > 1`, the digits of the count.

### Algorithm

1. Initialise `read = 0`, `write = 0`.
2. While `read < n`: record `ch = chars[read]`; advance `read` across the whole run, incrementing `count`.
3. Write `chars[write++] = ch`. If `count > 1`, convert `count` to its digit string and write each digit at `chars[write++]`.
4. Return `write` — the compressed length.

### Complexity

- **Time:** O(n) — `read` touches each character once; `write` performs no more writes than `read` performs reads.
- **Space:** O(1) — two indices plus a digit buffer bounded by `⌈log₁₀ n⌉ ≤ 4` for `n ≤ 2000`.

### Code

```go
func twoPointers(chars []byte) int {
	write := 0 // next slot to write compressed output
	read := 0  // scanning position over the input runs
	for read < len(chars) {
		ch := chars[read] // character starting this run
		count := 0
		// Consume the entire run of ch, counting its length.
		for read < len(chars) && chars[read] == ch {
			read++
			count++
		}
		chars[write] = ch // write the run's character
		write++
		if count > 1 {
			// Write the count's digits in order (in-place, left to right).
			for _, d := range strconv.Itoa(count) {
				chars[write] = byte(d) // rune digit fits in a byte ('0'..'9')
				write++
			}
		}
	}
	return write // length of the compressed array
}
```

### Dry Run

Example 1: `chars = ["a","a","b","b","c","c","c"]`.

| read (before) | ch | run consumed → read after | count | Writes | write after | chars[0:write] |
|---------------|----|---------------------------|-------|--------|-------------|-----------------|
| 0 | a | 2 | 2 | `a`,`2` | 2 | `a2` |
| 2 | b | 4 | 2 | `b`,`2` | 4 | `a2b2` |
| 4 | c | 7 | 3 | `c`,`3` | 6 | `a2b2c3` |

Loop ends (`read = 7 = n`). Return `write = 6`; `chars[:6] = "a2b2c3"` ✔.

---

## Key Takeaways

- **In-place because the output can't outgrow the input.** For any run of length ≥ 2, `1 + digits(L) ≤ L`, so the compressed form is never longer than the source — this is *why* a single array with read/write pointers is safe. Always verify this "write never passes read" invariant before compressing in place.
- **Multi-digit counts are just digits appended left-to-right.** A run of 12 becomes `'1','2'` (two array cells), not the byte `12`. Convert the integer to its decimal string and stream the characters.
- **Group-then-emit** is the reusable RLE skeleton: outer loop finds a maximal run, inner action emits the encoding. It generalises to decoding, run-length image compression, and "collapse consecutive duplicates" problems.
- **Single-character runs emit no count** — the `count > 1` guard is the whole trick that keeps `"a"` as `"a"` rather than `"a1"`.

---

## Related Problems

- LeetCode #38 — Count and Say (run-length encoding, but building outward)
- LeetCode #271 — Encode and Decode Strings (length-prefixed serialization)
- LeetCode #26 — Remove Duplicates from Sorted Array (read/write two-pointer in place)
- LeetCode #604 — Design Compressed String Iterator (decode an RLE stream lazily)
