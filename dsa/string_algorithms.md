# String Algorithms

> Pattern matching, string transformation, parsing, and canonical-form techniques
> for problems whose input is text.

---

## What the concept is

"String algorithms" is the umbrella for techniques that exploit the structure of
character sequences rather than treating them as generic arrays:

1. **Pattern matching** — find occurrences of a pattern `p` (length m) inside a
   text `s` (length n): brute force O(n·m), **KMP** O(n+m), **Rabin-Karp**
   (rolling hash) O(n+m) average, **Z-algorithm** O(n+m).
2. **Palindrome techniques** — expand-around-center O(n²)/O(1),
   **Manacher's algorithm** O(n).
3. **Canonical form + hash map** — map each string to a normalized key
   (sorted characters, character-count signature) so "equivalent" strings
   collide, e.g. grouping anagrams.
4. **String building & transformation** — constructing output efficiently with
   `strings.Builder` (Go strings are immutable; `+=` in a loop is O(L²)),
   run-length encoding, digit-by-digit arithmetic on numeric strings.
5. **Parsing / scanning** — single-pass state machines for atoi, valid-number
   validation, path simplification, tokenizing with `strings.Split` /
   `strings.Fields`.

Related-but-separate families (covered in their own files): sliding window over
a string (`sliding_window.md`), two pointers (`two_pointers.md`), string DP
like edit distance / regex matching (`dynamic_programming_2d.md`), and prefix
trees (`trie.md`).

### How to recognise it — signals in the problem statement

- "Find the first occurrence of `needle` in `haystack`" / "does `s` contain
  `t`?" → **substring search** (brute force → KMP/Rabin-Karp).
- "Longest palindromic …", "is it a palindrome after …" → **palindrome
  expansion / Manacher / two pointers**.
- "Group / count anagrams", "are these permutations of each other" →
  **canonical form + hash map**.
- "Add / multiply two numbers given as strings" (too big for int64) →
  **digit-by-digit simulation with carry**.
- "Parse / validate / convert this string" (atoi, Roman numerals, valid
  number, simplify path) → **linear scan, often a small state machine**.
- "Build / print the string that results from …" (zigzag, count-and-say,
  text justification) → **simulation + `strings.Builder`**.
- Constraints like `n ≤ 10⁵` on a substring-search problem → the O(n·m)
  brute force may TLE; reach for **KMP or rolling hash**.
- ASCII/lowercase-only alphabet mentioned → a fixed-size `[26]int` /
  `[128]int` count array beats a `map[byte]int`.

---

## General templates (Go)

### 1. KMP substring search — O(n+m) time, O(m) space

The failure function `lps[i]` = length of the longest proper prefix of
`p[0..i]` that is also a suffix. On a mismatch, the pattern pointer falls back
to `lps[j-1]` instead of restarting — the text pointer never moves backwards.

```go
// buildLPS computes the longest-proper-prefix-that-is-also-suffix table.
func buildLPS(p string) []int {
    lps := make([]int, len(p)) // lps[0] is always 0
    length := 0                // length of the current matched prefix
    for i := 1; i < len(p); {
        if p[i] == p[length] {
            length++          // extend the prefix-suffix match
            lps[i] = length
            i++
        } else if length > 0 {
            length = lps[length-1] // fall back to the next shorter border
        } else {
            lps[i] = 0        // no border at all
            i++
        }
    }
    return lps
}

// kmpSearch returns the index of the first occurrence of p in s, or -1.
func kmpSearch(s, p string) int {
    if len(p) == 0 {
        return 0
    }
    lps := buildLPS(p)
    j := 0 // chars of p matched so far
    for i := 0; i < len(s); i++ {
        for j > 0 && s[i] != p[j] {
            j = lps[j-1]      // mismatch: reuse the longest border
        }
        if s[i] == p[j] {
            j++               // match: advance in the pattern
        }
        if j == len(p) {
            return i - len(p) + 1 // full pattern matched ending at i
        }
    }
    return -1
}
```

### 2. Rabin-Karp rolling hash — O(n+m) average

```go
// rabinKarp returns the first index of p in s, or -1.
func rabinKarp(s, p string) int {
    n, m := len(s), len(p)
    if m > n {
        return -1
    }
    const base, mod = 256, 1_000_000_007
    // pow = base^(m-1) % mod, used to remove the leading char when sliding.
    pow := 1
    for i := 0; i < m-1; i++ {
        pow = pow * base % mod
    }
    hp, hs := 0, 0 // hash of pattern, hash of current window
    for i := 0; i < m; i++ {
        hp = (hp*base + int(p[i])) % mod
        hs = (hs*base + int(s[i])) % mod
    }
    for i := 0; ; i++ {
        // Hashes equal → verify char-by-char to rule out collisions.
        if hs == hp && s[i:i+m] == p {
            return i
        }
        if i+m == n {
            return -1
        }
        // Slide: drop s[i], append s[i+m].
        hs = ((hs-int(s[i])*pow%mod+mod)%mod*base + int(s[i+m])) % mod
    }
}
```

### 3. Expand around center (palindromes) — O(n²) time, O(1) space

```go
// expand grows outward from (l, r) while the characters match,
// returning the length of the palindrome found.
func expand(s string, l, r int) int {
    for l >= 0 && r < len(s) && s[l] == s[r] {
        l--
        r++
    }
    return r - l - 1 // window overshot by one on each side
}

// longestPalindrome tries every center: 2n-1 of them (chars and gaps).
func longestPalindrome(s string) string {
    start, maxLen := 0, 0
    for i := range s {
        if l := expand(s, i, i); l > maxLen { // odd-length center
            maxLen, start = l, i-(l-1)/2
        }
        if l := expand(s, i, i+1); l > maxLen { // even-length center
            maxLen, start = l, i-(l-2)/2
        }
    }
    return s[start : start+maxLen]
}
```

### 4. Canonical form + hash map (anagram grouping)

```go
// key builds a character-count signature — identical for all anagrams.
func key(s string) [26]int {
    var cnt [26]int
    for i := 0; i < len(s); i++ {
        cnt[s[i]-'a']++ // fixed alphabet → array beats map
    }
    return cnt // arrays are comparable in Go, so usable as a map key
}

groups := map[[26]int][]string{}
for _, w := range words {
    k := key(w)
    groups[k] = append(groups[k], w)
}
```

### 5. Efficient string building

```go
var b strings.Builder
b.Grow(n)            // optional: pre-allocate when the size is known
for _, ch := range parts {
    b.WriteByte(ch)  // O(1) amortised — never `result += string(ch)` in a loop
}
return b.String()
```

---

## Worked example — KMP on LeetCode #28

Search for pattern **`p = "ababa"`** in text **`s = "ababcababa"`** — a case
where the naive search wastes work re-scanning and KMP's fallback shines.

**Step 1 — build `lps` for `p = "ababa"`:**

| i | p[i] | compare p[i] vs p[length] | action | lps |
|---|------|---------------------------|--------|-----|
| 1 | b | b vs a (length=0) — mismatch, length=0 | lps[1]=0, i=2 | [0,0,·,·,·] |
| 2 | a | a vs a — match | length=1, lps[2]=1, i=3 | [0,0,1,·,·] |
| 3 | b | b vs b (length=1) — match | length=2, lps[3]=2, i=4 | [0,0,1,2,·] |
| 4 | a | a vs a (length=2) — match | length=3, lps[4]=3 | [0,0,1,2,3] |

**Step 2 — scan `s = "ababcababa"` with `j` = chars matched:**

| i | s[i] | before: j | comparison | after: j | note |
|---|------|-----------|------------|----------|------|
| 0 | a | 0 | a==p[0] ✓ | 1 | |
| 1 | b | 1 | b==p[1] ✓ | 2 | |
| 2 | a | 2 | a==p[2] ✓ | 3 | |
| 3 | b | 3 | b==p[3] ✓ | 4 | |
| 4 | c | 4 | c≠p[4] → j=lps[3]=2; c≠p[2] → j=lps[1]=0; c≠p[0] | 0 | fell back **without moving i** |
| 5 | a | 0 | a==p[0] ✓ | 1 | |
| 6 | b | 1 | b==p[1] ✓ | 2 | |
| 7 | a | 2 | a==p[2] ✓ | 3 | |
| 8 | b | 3 | b==p[3] ✓ | 4 | |
| 9 | a | 4 | a==p[4] ✓ | 5 = m → **match at i-m+1 = 5** | |

Result: first occurrence at index **5**. Total work: each character of `s` is
read once and `j` only decreases via `lps` (bounded by total increases), so the
scan is O(n); the table build is O(m).

---

## Common pitfalls

- **O(L²) concatenation** — `result += s` in a loop reallocates every time.
  Use `strings.Builder` (or a `[]byte` you convert once at the end).
- **Bytes vs runes in Go** — `s[i]` is a `byte`; `for _, r := range s` yields
  `rune`s and the *index jumps by UTF-8 width*. For LeetCode (ASCII inputs)
  index with bytes; only use runes when Unicode is explicitly in scope.
- **KMP fallback must loop** — after a mismatch you may need to fall back
  *repeatedly* (`for j > 0 && mismatch { j = lps[j-1] }`), not just once.
- **Forgetting the empty-pattern convention** — searching for `""` returns 0
  (matches `strings.Index` and C's `strstr`).
- **Rolling hash without verification** — hash equality can be a collision;
  always confirm with a direct comparison before reporting a match (or use
  double hashing).
- **Negative values in modular rolling hash** — `(hs - c*pow) % mod` can go
  negative in Go; add `mod` before taking `%`.
- **Even-length palindrome centers** — expanding only from single characters
  misses `"abba"`; check all `2n-1` centers (`(i,i)` and `(i,i+1)`).
- **Sorting as anagram key is O(k log k) per word** — a `[26]int` count array
  is O(k) and, being a comparable array, works directly as a Go map key.
- **Off-by-one when recovering the start index** — after `expand` returns
  length `l` centered near `i`, the start is `i-(l-1)/2` (odd) vs `i-(l-2)/2`
  (even). Derive it once, comment it, and reuse.
- **Integer overflow when parsing numeric strings** — atoi-style problems
  require checking the overflow *before* multiplying by 10 and adding a digit.

---

## Problems in this repo

Core string-algorithm problems:

- [0005 — Longest Palindromic Substring](../0005_longest_palindromic_substring/README.md) — expand around center, Manacher's algorithm
- [0028 — Find the Index of the First Occurrence in a String](../0028_find_the_index_of_the_first_occurrence_in_a_string/README.md) — brute force vs **KMP**
- [0049 — Group Anagrams](../0049_group_anagrams/README.md) — canonical form (sorted key / count array) + hash map
- [0125 — Valid Palindrome](../0125_valid_palindrome/README.md) — two-pointer palindrome check with filtering

Parsing / scanning and conversion:

- [0008 — String to Integer (atoi)](../0008_string_to_integer_atoi/README.md) — linear-scan parser with overflow checks
- [0012 — Integer to Roman](../0012_integer_to_roman/README.md) / [0013 — Roman to Integer](../0013_roman_to_integer/README.md) — symbol-table conversion
- [0014 — Longest Common Prefix](../0014_longest_common_prefix/README.md) — horizontal / vertical scan
- [0058 — Length of Last Word](../0058_length_of_last_word/README.md) — reverse scan
- [0065 — Valid Number](../0065_valid_number/README.md) — state-machine validation
- [0071 — Simplify Path](../0071_simplify_path/README.md) — split + stack of path components
- [0093 — Restore IP Addresses](../0093_restore_ip_addresses/README.md) — segment validation + backtracking

String building / simulation / string arithmetic:

- [0006 — Zigzag Conversion](../0006_zigzag_conversion/README.md) — row simulation with `strings.Builder`
- [0038 — Count and Say](../0038_count_and_say/README.md) — run-length encoding
- [0043 — Multiply Strings](../0043_multiply_strings/README.md) — digit-by-digit multiplication with carry
- [0067 — Add Binary](../0067_add_binary/README.md) — digit-by-digit addition with carry
- [0068 — Text Justification](../0068_text_justification/README.md) — greedy line packing + space distribution

Adjacent patterns on strings (see their own concept files):

- [0003 — Longest Substring Without Repeating Characters](../0003_longest_substring_without_repeating_characters/README.md) — sliding window
- [0030 — Substring with Concatenation of All Words](../0030_substring_with_concatenation_of_all_words/README.md) — fixed-stride sliding window over words
- [0076 — Minimum Window Substring](../0076_minimum_window_substring/README.md) — sliding window with need/have counts
- [0010 — Regular Expression Matching](../0010_regular_expression_matching/README.md) / [0044 — Wildcard Matching](../0044_wildcard_matching/README.md) — 2-D string DP
- [0072 — Edit Distance](../0072_edit_distance/README.md), [0097 — Interleaving String](../0097_interleaving_string/README.md), [0115 — Distinct Subsequences](../0115_distinct_subsequences/README.md) — 2-D string DP
- [0127 — Word Ladder](../0127_word_ladder/README.md) / [0126 — Word Ladder II](../0126_word_ladder_ii/README.md) — BFS over word transformations

*(Problems 0131+ will be linked in a later pass.)*
