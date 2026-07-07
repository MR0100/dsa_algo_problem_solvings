# 0481 — Magical String

> LeetCode #481 · Difficulty: Medium
> **Categories:** Two Pointers, String, Simulation

---

## Problem Statement

A magical string `s` consists of only `'1'` and `'2'` and obeys the following rules:

- The string `s` is **magical** because concatenating the number of contiguous occurrences of characters `'1'` and `'2'` generates the string `s` itself.

The first few elements of `s` is `s = "1221121221221121122……"`. If we group the consecutive `'1'`s and `'2'`s in `s`, it will be `"1 22 11 2 1 22 1 22 11 2 11 22 ......"` and the occurrences of `'1'`s or `'2'`s in each group are `"1 2 2 1 1 2 1 2 2 1 2 2 ......"`. You can see that the occurrence sequence is `s` itself.

Given an integer `n`, return the number of `'1'`s in the first `n` number in the magical string `s`.

**Example 1:**

```
Input: n = 6
Output: 3
Explanation: The first 6 elements of magical string s is "122112" and it contains three 1's, so return 3.
```

**Example 2:**

```
Input: n = 1
Output: 1
```

**Constraints:**

- `1 <= n <= 10^5`

---

## Company Frequency

| Company    | Frequency        | Last Reported |
|------------|------------------|---------------|
| Google     | ★★★☆☆ Medium     | 2023          |
| Amazon     | ★★☆☆☆ Low        | 2023          |
| Microsoft  | ★★☆☆☆ Low        | 2022          |
| Bloomberg  | ★★☆☆☆ Low        | 2022          |
| Adobe      | ★☆☆☆☆ Rare       | 2021          |

> ⚠️ Frequency data is crowd-sourced from LeetCode Discuss, Glassdoor, and
> community interview reports. Treat as a signal, not a guarantee.

---

## DSA Concepts Used

- **Two Pointers (slow read / fast write)** — a slow pointer walks the prefix already built and reads each value as a *run length*, while the string keeps growing at the tail. The read pointer never overtakes the write pointer, so a single forward scan generates the whole sequence → see [`/dsa/two_pointers.md`](/dsa/two_pointers.md)
- **String / self-referential simulation** — the sequence is defined by decoding its own characters as run-length instructions; recognising this "the data is its own program" structure is the key insight → see [`/dsa/string_algorithms.md`](/dsa/string_algorithms.md)

---

## Approaches Overview

| # | Approach | Time | Space | When to use |
|---|----------|------|-------|-------------|
| 1 | Full Generation with a Group Queue | O(n) | O(n) | Clearest mental model: build the whole string, then count |
| 2 | In-Place Two-Pointer, Count While Building (Optimal) | O(n) | O(n) | Fewest passes; count folded into generation, stops exactly at n |

> Both are O(n) time and space — the string itself is unavoidable. Approach 2 shaves the constant factor (one pass instead of two) and never writes past index `n-1`.

---

## Approach 1 — Full Generation with a Group Queue

### Intuition

The magical string is *self-describing*: reading its own digits left to right tells you the length of each successive run. `s[0]=1` ⇒ the first run (the `"1"`) has length 1. `s[1]=2` ⇒ the second run (the `"22"`) has length 2. `s[2]=2` ⇒ the third run (the `"11"`) has length 2, and so on. The digit of each run **strictly alternates** `1,2,1,2,…`, so we never have to decide *what* to append, only *how many*. A "read head" index walking the already-built prefix supplies those counts. Because every run is at most length 2, the read head can never catch up to the write head, so the string keeps growing until it is long enough. Then we simply count the 1's in the first `n` characters.

### Algorithm

1. Handle tiny `n`: the bootstrap prefix `"122"` covers `n = 1,2,3` and always contains exactly one `1`.
2. Seed the string as `s = [1,2,2]` (the first digit can't be derived from an empty string, so we hard-code the self-consistent start).
3. Set `head = 2` (the first run length not yet consumed) and `next = 1` (digit of the next run).
4. While `len(s) < n`: append `s[head]` copies of `next`, flip `next` between 1 and 2, and increment `head`.
5. Count and return the number of 1's among `s[0..n-1]`.

### Complexity

- **Time:** O(n) — each character is produced exactly once (a run appends ≤ 2), and the final count is one O(n) sweep.
- **Space:** O(n) — the generated slice grows to length ≈ `n`.

### Code

```go
func generateWithQueue(n int) int {
	if n == 0 {
		return 0 // no characters → no 1's
	}
	if n <= 3 {
		// The bootstrap prefix "122" already covers n = 1,2,3. Its 1's:
		// "1"→1, "12"→1, "122"→1. So the count is always 1 here.
		return 1
	}

	s := []int{1, 2, 2} // known self-consistent seed of the magical string
	head := 2           // s[2]=2 is the next run length we have not applied yet
	next := 1           // the digit that the next run will consist of (alternates)

	// Grow until we have at least n characters; a run of length s[head]
	// appends that many copies of `next`.
	for len(s) < n {
		runLen := s[head] // how many identical digits the next group holds (1 or 2)
		for i := 0; i < runLen; i++ {
			s = append(s, next) // emit one digit of the current run
		}
		next ^= 3 // flip 1<->2 (1^3=2, 2^3=1) so runs alternate digit
		head++    // consume the next run-length instruction
	}

	ones := 0
	for i := 0; i < n; i++ { // count 1's within the first n characters only
		if s[i] == 1 {
			ones++
		}
	}
	return ones
}
```

### Dry Run

Example 1: `n = 6`. Seed `s = [1,2,2]`, `head = 2`, `next = 1`, target length 6.

| Iteration | len(s) < 6? | s[head] (run len) | next (digit appended) | s after append | next after flip | head after |
|-----------|-------------|-------------------|-----------------------|----------------|-----------------|------------|
| start | — | — | — | `[1,2,2]` | 1 | 2 |
| 1 | yes (3) | s[2]=2 | 1 | `[1,2,2,1,1]` | 2 | 3 |
| 2 | yes (5) | s[3]=1 | 2 | `[1,2,2,1,1,2]` | 1 | 4 |
| 3 | no (6) | — | — | stop | — | — |

Count 1's in `s[0..5] = [1,2,2,1,1,2]` → indices 0, 3, 4 are 1 ⇒ **3** ✔

---

## Approach 2 — In-Place Two-Pointer, Count While Building (Optimal)

### Intuition

Exactly the same self-referential generation, but we fold the answer into the build loop. A slow read pointer `i` supplies run lengths; a fast write appends the alternating digit; and the instant we append a `1` we bump a running `ones` counter. We also stop the moment the string reaches length `n`, so we never produce or count a character past position `n-1`. This eliminates the separate counting pass and any need to remember state beyond `ones`.

### Algorithm

1. Return 0 for `n = 0`; return 1 for `n ≤ 3` (prefix `"122"` always has one `1`).
2. Seed `s = [1,2,2]` and `ones = 1` (the seed's single `1` at index 0).
3. Read pointer `i = 2`; current run digit `digit = 1`.
4. While `len(s) < n`: append `s[i]` copies of `digit`, breaking early if we hit length `n`; for each appended `1`, increment `ones`. Then flip `digit` and advance `i`.
5. Return `ones`.

### Complexity

- **Time:** O(n) — one unit of work per produced character, with counting inlined.
- **Space:** O(n) — the string buffer; no auxiliary structures.

### Code

```go
func twoPointers(n int) int {
	if n == 0 {
		return 0 // empty prefix has zero 1's
	}
	if n <= 3 {
		return 1 // "1","12","122" each contain exactly one 1
	}

	s := []int{1, 2, 2} // magical-string seed
	ones := 1           // the seed contributes exactly one '1' (at index 0)
	i := 2              // read pointer: s[2] is the next run length to apply
	digit := 1          // digit of the run being appended right now (alternates)

	for len(s) < n {
		for c := 0; c < s[i]; c++ { // append s[i] copies of `digit`
			if len(s) >= n {
				break // never write past position n-1
			}
			s = append(s, digit)
			if digit == 1 { // count 1's inline, only for indices < n
				ones++
			}
		}
		digit ^= 3 // 1<->2 alternation for the next run
		i++        // advance the read pointer to the next run length
	}
	return ones
}
```

### Dry Run

Example 1: `n = 6`. Seed `s = [1,2,2]`, `ones = 1`, `i = 2`, `digit = 1`.

| Iteration | len(s) < 6? | s[i] (run len) | digit | appends (respecting n) | ones after | s after | digit flip | i after |
|-----------|-------------|----------------|-------|------------------------|------------|---------|-----------|---------|
| start | — | — | — | — | 1 | `[1,2,2]` | 1 | 2 |
| 1 | yes (3) | s[2]=2 | 1 | append 1,1 | 3 | `[1,2,2,1,1]` | 2 | 3 |
| 2 | yes (5) | s[3]=1 | 2 | append 2 (len hits 6) | 3 | `[1,2,2,1,1,2]` | 1 | 4 |
| 3 | no (6) | — | — | stop | 3 | — | — | — |

Return `ones = 3` ✔ — identical result, computed without a second scan.

---

## Key Takeaways

- **Self-generating sequences**: when a sequence *is* its own run-length encoding, a single slow/fast two-pointer pass builds it — the read pointer decodes counts, the write pointer emits digits, and a bounded run length (here ≤ 2) guarantees the reader never starves.
- **Alternating digit trick**: `x ^= 3` flips 1↔2 (`1^3=2`, `2^3=1`) in one op, cleaner than an `if`. Generalises to toggling between any two values `a`, `b` via `x ^= a^b`.
- **Count during construction** when the final scan would just re-read what you already produced — one pass, O(1) extra state.
- Hard-coding a tiny **self-consistent seed** (`"122"`) sidesteps the chicken-and-egg bootstrap of a self-referential definition.

---

## Related Problems

- LeetCode #38 — Count and Say (another self-describing string built by decoding runs)
- LeetCode #443 — String Compression (run-length encoding with two pointers)
- LeetCode #779 — K-th Symbol in Grammar (self-referential binary sequence)
- LeetCode #900 — RLE Iterator (decoding run-length pairs on demand)
