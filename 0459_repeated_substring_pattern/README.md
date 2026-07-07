# 0459 — Repeated Substring Pattern

> LeetCode #459 · Difficulty: Easy
> **Categories:** String, String Matching

---

## Problem Statement

Given a string `s`, check if it can be constructed by taking a substring of it and appending multiple copies of the substring together.

**Example 1:**

```
Input: s = "abab"
Output: true
Explanation: It is the substring "ab" twice.
```

**Example 2:**

```
Input: s = "aba"
Output: false
```

**Example 3:**

```
Input: s = "abcabcabcabc"
Output: true
Explanation: It is the substring "abc" four times or the substring "abcabc" twice.
```

**Constraints:**

- `1 <= s.length <= 10^4`
- `s` consists of lowercase English letters.

---

## Company Frequency

| Company    | Frequency        | Last Reported |
|------------|------------------|---------------|
| Amazon     | ★★★☆☆ Medium     | 2024          |
| Google     | ★★★☆☆ Medium     | 2023          |
| Microsoft  | ★★☆☆☆ Low        | 2023          |
| Adobe      | ★★☆☆☆ Low        | 2022          |
| Apple      | ★★☆☆☆ Low        | 2022          |

> ⚠️ Frequency data is crowd-sourced from LeetCode Discuss, Glassdoor, and
> community interview reports. Treat as a signal, not a guarantee.

---

## DSA Concepts Used

- **String Algorithms (KMP prefix function)** — the optimal solution derives the smallest period from `lps[n-1]` (longest proper prefix that is also a suffix); `s` is periodic iff that period divides `n` → see [`/dsa/string_algorithms.md`](/dsa/string_algorithms.md)
- **Two Pointers** — the brute-force period test walks index `i` against `i-L`, a lockstep two-index scan over the string → see [`/dsa/two_pointers.md`](/dsa/two_pointers.md)

---

## Approaches Overview

| # | Approach | Time | Space | When to use |
|---|----------|------|-------|-------------|
| 1 | Try every divisor period | O(n·d(n)) ⊆ O(n²) | O(1) | Most intuitive; fine for n ≤ 10⁴ |
| 2 | Concatenation trick `(s+s)[1:-1]` | O(n) with KMP search / O(n²) naive | O(n) | Shortest to write; one substring search |
| 3 | KMP failure function (Optimal) | O(n) | O(n) | Cleanest linear answer, no doubling |

---

## Approach 1 — Try Every Divisor Period

### Intuition

If `s` is `k` copies of a block of length `L`, then `L` must be a proper divisor of `n = len(s)`, and every character must equal the one `L` positions earlier (`s[i] == s[i-L]`). So enumerate candidate block lengths `L` that divide `n` (only up to `n/2`, since the block must repeat at least twice) and verify the tiling condition for each. The first `L` that satisfies it proves periodicity.

### Algorithm

1. `n = len(s)`. For `L` from `1` to `n/2`:
   1. Skip if `n % L != 0` (a length-`L` block can't tile `n` evenly).
   2. Check `s[i] == s[i-L]` for every `i` in `[L, n)`. If all match, return `true`.
2. If no divisor works, return `false`.

### Complexity

- **Time:** O(n · d(n)) ⊆ O(n²) — at most `n/2` candidate periods, each verified with an O(n) scan (`d(n)` = number of divisors).
- **Space:** O(1) — only index arithmetic.

### Code

```go
func divisorBruteForce(s string) bool {
	n := len(s)
	for l := 1; l <= n/2; l++ { // candidate block length; a real period is ≤ n/2
		if n%l != 0 {
			continue // block of length l can't tile n characters evenly
		}
		ok := true
		// Verify periodicity: each char must equal the one one block earlier.
		for i := l; i < n; i++ {
			if s[i] != s[i-l] {
				ok = false // mismatch → l is not a valid period
				break
			}
		}
		if ok {
			return true // s is the length-l block repeated n/l times
		}
	}
	return false // no proper divisor period rebuilds s
}
```

### Dry Run

Example 1: `s = "abab"`, `n = 4`. Candidate `L` in `1..2`.

| L | n % L == 0? | tiling check s[i]==s[i-L] | result |
|---|-------------|---------------------------|--------|
| 1 | 4%1=0 yes | i=1: s[1]='b' vs s[0]='a' ✗ | fail |
| 2 | 4%2=0 yes | i=2: s[2]='a' vs s[0]='a' ✓; i=3: s[3]='b' vs s[1]='b' ✓ | **all match → true** |

Block `"ab"` (length 2) tiles `"abab"` → return `true` ✔

---

## Approach 2 — Concatenation Trick (s+s Doubling)

### Intuition

Form `t = s + s`. If `s` is `k ≥ 2` copies of a block of length `L`, then a copy of `s` begins again at offset `L` inside `t`, where `0 < L < n`. Removing the first and last characters of `t` destroys the two *trivial* occurrences (the one starting at offset 0 and the one starting at offset `n`). So `s` still appears in `t[1 : 2n-1]` **iff** a genuine internal period exists. A single substring search answers the question. (Slick to remember: `s in (s + s)[1:-1]`.)

### Algorithm

1. `doubled = s + s`.
2. `middle = doubled[1 : len(doubled)-1]` (strip first and last char).
3. Return whether `middle` contains `s`.

### Complexity

- **Time:** O(n) when the underlying search is KMP/Z (Go's `strings.Contains` is a tuned algorithm and behaves linearly here); O(n²) with a naive character-by-character search.
- **Space:** O(n) — the doubled string.

### Code

```go
func concatTrick(s string) bool {
	doubled := s + s                        // two back-to-back copies
	middle := doubled[1 : len(doubled)-1]   // drop the trivial offset-0 and offset-n matches
	return strings.Contains(middle, s)      // a surviving match ⇒ a real internal period
}
```

### Dry Run

Example 1: `s = "abab"`.

| Step | value |
|------|-------|
| doubled = s + s | `"abababab"` (length 8) |
| middle = doubled[1:7] | `"bababa"` |
| Does `"bababa"` contain `"abab"`? | yes — `"bababa"` = `b·abab·a`, `s` sits at offset 1 |

Match found in the stripped double → return `true` ✔

For contrast, `s = "aba"`: doubled = `"abaaba"`, middle = `"baab"`, which does **not** contain `"aba"` → `false`.

---

## Approach 3 — KMP Failure Function Period (Optimal)

### Intuition

The KMP prefix function `lps[i]` is the length of the longest proper prefix of `s[0..i]` that is also a suffix of it. For the whole string, `border = lps[n-1]` is the longest prefix that reappears as a suffix, and `p = n - border` is the **smallest period** of `s`. The string is a repetition of a shorter block exactly when that period actually tiles the string: `border > 0` (a non-empty overlap exists) **and** `n % p == 0`. This needs a single linear prefix-function pass — no doubling, no divisor loop.

### Algorithm

1. If `n < 2`, return `false`.
2. Build `lps[]` via the standard KMP prefix-function construction in O(n).
3. Let `last = lps[n-1]`, `p = n - last`.
4. Return `last > 0 && n % p == 0`.

### Complexity

- **Time:** O(n) — the prefix function amortises to a single linear scan.
- **Space:** O(n) — the `lps` array.

### Code

```go
func kmpFailure(s string) bool {
	n := len(s)
	if n < 2 {
		return false // a single character can't be a repeat of a shorter block
	}
	lps := make([]int, n) // lps[i] = length of longest proper prefix==suffix of s[0..i]
	length := 0           // length of the current matching prefix
	// Standard KMP prefix-function build.
	for i := 1; i < n; i++ {
		// Fall back through shorter borders while characters disagree.
		for length > 0 && s[i] != s[length] {
			length = lps[length-1] // reuse the next-longest border
		}
		if s[i] == s[length] {
			length++ // extend the current border by one character
		}
		lps[i] = length // record the border length ending at i
	}
	last := lps[n-1] // longest border of the whole string
	period := n - last
	// Periodic iff a non-empty border exists and its induced period tiles n.
	return last > 0 && n%period == 0
}
```

### Dry Run

Example 1: `s = "abab"`, `n = 4`. Build `lps` (`length` starts at 0):

| i | s[i] | compare s[i] vs s[length] | length after | lps[i] |
|---|------|---------------------------|--------------|--------|
| 1 | 'b' | 'b' vs s[0]='a' mismatch, length stays 0 | 0 | 0 |
| 2 | 'a' | 'a' vs s[0]='a' match | 1 | 1 |
| 3 | 'b' | 'b' vs s[1]='b' match | 2 | 2 |

`last = lps[3] = 2`, `period = 4 - 2 = 2`. Check `last > 0` (yes) and `4 % 2 == 0` (yes) → return `true` ✔ — the smallest period is 2 (block `"ab"`), and it tiles the length-4 string.

---

## Key Takeaways

- **`period = n - lps[n-1]`** is the smallest period of a string; it tiles the string iff `n % period == 0`. This one line turns KMP's prefix function into a periodicity test.
- **The `(s+s)[1:-1]` trick** is the fastest to recall: a string is periodic iff it reappears in its own doubling with the two trivial ends stripped. Great for a 1-liner, but know *why* the strip removes exactly the offset-0 and offset-n matches.
- **A valid period must divide `n` and be ≤ n/2.** Both the brute force and the KMP check hinge on the divisibility condition — repetition means integer tiling.
- **Prefix function ≠ just substring search.** It also encodes borders/periods, which powers this problem, string compression, and shortest-repeating-unit questions.

---

## Related Problems

- LeetCode #28 — Find the Index of the First Occurrence in a String (KMP search itself)
- LeetCode #686 — Repeated String Match (how many copies until a pattern fits)
- LeetCode #1392 — Longest Happy Prefix (direct `lps[n-1]` output)
- LeetCode #1668 — Maximum Repeating Substring (repetition counting)
- LeetCode #214 — Shortest Palindrome (KMP prefix function on a transformed string)
