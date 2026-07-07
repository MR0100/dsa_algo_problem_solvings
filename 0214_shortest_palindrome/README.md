# 0214 — Shortest Palindrome

> LeetCode #214 · Difficulty: Hard
> **Categories:** String, KMP, Rolling Hash, Two Pointers

---

## Problem Statement

You are given a string `s`. You can convert `s` to a palindrome by adding characters in front of it.

Return *the shortest palindrome you can find by performing this transformation*.

**Example 1:**
```
Input: s = "aacecaaa"
Output: "aaacecaaa"
```

**Example 2:**
```
Input: s = "abcd"
Output: "dcbabcd"
```

**Constraints:**
- `0 <= s.length <= 5 * 10⁴`
- `s` consists of lowercase English letters only.

---

## Company Frequency

| Company    | Frequency        | Last Reported |
|------------|------------------|---------------|
| Amazon     | ★★★☆☆ Medium     | 2024          |
| Google     | ★★★☆☆ Medium     | 2023          |
| Microsoft  | ★★☆☆☆ Low        | 2023          |
| Apple      | ★★☆☆☆ Low        | 2022          |

> ⚠️ Frequency data is crowd-sourced from LeetCode Discuss, Glassdoor, and
> community interview reports. Treat as a signal, not a guarantee.

---

## DSA Concepts Used

- **String Algorithms / KMP failure function** — the longest palindromic prefix equals the longest prefix of `s` that is a suffix of `reverse(s)`, computed by one KMP `lps` pass over `s + '#' + reverse(s)` → see [`/dsa/string_algorithms.md`](/dsa/string_algorithms.md)
- **Two Pointers** — palindrome and reverse checks scan from both ends inward → see [`/dsa/two_pointers.md`](/dsa/two_pointers.md)

---

## Approaches Overview

Let n = length of `s`.

| # | Approach | Time | Space | When to use |
|---|----------|------|-------|-------------|
| 1 | Brute Force (longest palindromic prefix) | O(n²) | O(n) | Small inputs; the clearest correctness argument |
| 2 | KMP failure function (Optimal) | O(n) | O(n) | Large inputs (n up to 5·10⁴) — the intended answer |

---

## Approach 1 — Brute Force (Longest Palindromic Prefix)

### Intuition
Characters may only be added to the **front**, so the palindrome we build has `s` as its suffix. The characters we prepend must mirror the tail of `s` that is not already symmetric. Concretely: find the **longest prefix** `s[0..k)` that is itself a palindrome — it needs nothing added in front of it. The remaining suffix `s[k:]` must be mirrored and placed at the very front: answer = `reverse(s[k:]) + s`. Fewer added characters ⇔ longer palindromic prefix, so maximise `k`.

### Algorithm
1. For `j` from `n` down to `0`, test whether `s[0..j)` is a palindrome.
2. The first (largest) `j` that passes is the longest palindromic prefix length `k`.
3. Return `reverse(s[k:]) + s`.

### Complexity
- **Time:** O(n²) — up to n prefix tests, each an O(n) palindrome check.
- **Space:** O(n) for the reversed suffix and the result string.

### Code
```go
func bruteForce(s string) string {
	n := len(s)
	for j := n; j >= 0; j-- {
		if isPalindrome(s[:j]) {
			return reverse(s[j:]) + s
		}
	}
	return s
}

func isPalindrome(t string) bool {
	for i, j := 0, len(t)-1; i < j; i, j = i+1, j-1 {
		if t[i] != t[j] {
			return false
		}
	}
	return true
}

func reverse(t string) string {
	b := []byte(t)
	for i, j := 0, len(b)-1; i < j; i, j = i+1, j-1 {
		b[i], b[j] = b[j], b[i]
	}
	return string(b)
}
```

### Dry Run (Example 1: `s = "aacecaaa"`, n = 8)

| j | prefix `s[:j]` | palindrome? |
|---|----------------|-------------|
| 8 | `aacecaaa` | `a…a` ok, but `a c e c a` centre mismatch → no |
| 7 | `aacecaa` | position 0 `a` vs 6 `a` ok, 1 `a` vs 5 `a` ok, 2 `c` vs 4 `c` ok, centre `e` → **yes** |

Longest palindromic prefix `k = 7` (`aacecaa`). Suffix `s[7:] = "a"`, `reverse("a") = "a"`.
Answer = `"a" + "aacecaaa"` = **`"aaacecaaa"`** ✓

---

## Approach 2 — KMP Failure Function (Optimal)

### Intuition
The longest palindromic prefix of `s` equals the longest prefix of `s` that is also a suffix of `reverse(s)`. Form `combined = s + sep + reverse(s)` with a separator (`'#'`) that never occurs in `s`, so no matched span can straddle it. Running KMP's failure function (`lps`) over `combined`, the **last** `lps` value is exactly the length `k` of the longest prefix-of-`s` matching a suffix-of-`reverse(s)` — i.e. the longest palindromic prefix. Then prepend `reverse(s[k:])`.

### Algorithm
1. `combined = s + "#" + reverse(s)`.
2. Build `lps[]` where `lps[i]` = length of the longest proper prefix of `combined[:i+1]` that is also a suffix of it.
3. `k = lps[last]` = longest palindromic prefix length.
4. Return `reverse(s[k:]) + s` (implemented as `rev[:n-k] + s`, since `reverse(s[k:]) == reverse(s)[:n-k]`).

### Complexity
- **Time:** O(n) — one KMP failure-array build over a length-`~2n` string.
- **Space:** O(n) for the `lps` array and the result.

### Code
```go
func kmp(s string) string {
	if len(s) == 0 {
		return ""
	}
	rev := reverse(s)
	combined := s + "#" + rev
	lps := buildLPS(combined)
	k := lps[len(lps)-1]
	return rev[:len(s)-k] + s
}

func buildLPS(t string) []int {
	lps := make([]int, len(t))
	length := 0
	for i := 1; i < len(t); i++ {
		for length > 0 && t[i] != t[length] {
			length = lps[length-1] // fall back along failure links
		}
		if t[i] == t[length] {
			length++
		}
		lps[i] = length
	}
	return lps
}
```

### Dry Run (Example 1: `s = "aacecaaa"`)

`rev = "aaacecaa"`, `combined = "aacecaaa#aaacecaa"` (length 17).

Key `lps` progression (index : char : lps):

| i | char | lps[i] | note |
|---|------|--------|------|
| 0 | a | 0 | base |
| 1 | a | 1 | `a` matches prefix `a` |
| 2 | c | 0 | reset |
| … | … | … | separator `#` forces lps back to 0 |
| 8 | # | 0 | boundary — matches never cross it |
| 9..16 | a a a c e c a a | … | matches climb along the reversed half |
| 16 | a | **7** | final value = longest palindromic prefix length |

`k = lps[16] = 7`. `rev[:n-k] = rev[:1] = "a"`. Answer = `"a" + "aacecaaa"` = **`"aaacecaaa"`** ✓

---

## Key Takeaways

- **"Add to front to make a palindrome" ⇒ find the longest palindromic prefix.** Only the non-palindromic tail must be mirrored and prepended; that reversed tail is the minimal addition.
- **KMP turns an O(n²) prefix search into O(n).** The trick — concatenate `s + '#' + reverse(s)` and read off the final `lps` — reappears in many "prefix-equals-suffix" string problems.
- **The separator is mandatory.** Without a character absent from `s`, the failure function could match across the join and report a `k` larger than `n`, breaking the slice `rev[:n-k]`.
- **`reverse(s[k:]) == reverse(s)[:n-k]`** — reversing once and slicing avoids reversing a substring separately.
- Rolling hashes (double hashing to dodge collisions) are an alternative O(n) route with the same longest-palindromic-prefix idea.

---

## Related Problems

- LeetCode #5 — Longest Palindromic Substring (palindrome structure)
- LeetCode #28 — Find the Index of the First Occurrence in a String (KMP itself)
- LeetCode #459 — Repeated Substring Pattern (KMP failure-function trick)
- LeetCode #125 — Valid Palindrome (two-pointer palindrome check)
- LeetCode #336 — Palindrome Pairs (prefix/suffix palindrome matching)
