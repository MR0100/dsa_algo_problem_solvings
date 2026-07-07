# 0030 — Substring with Concatenation of All Words

> LeetCode #30 · Difficulty: Hard
> **Categories:** Hash Table, String, Sliding Window

---

## Problem Statement

You are given a string `s` and an array of strings `words`. All the strings of `words` are of **the same length**.

A **concatenated string** is a string that exactly contains all the strings of any permutation of `words` concatenated.

Return an array of the **starting indices** of all the concatenated substrings in `s`. You can return the answer in **any order**.

**Example 1**
```
Input:  s = "barfoothefoobarman", words = ["foo","bar"]
Output: [0,9]
Explanation:
  The substring starting at 0 is "barfoo" = "bar"+"foo".
  The substring starting at 9 is "foobar" = "foo"+"bar".
```

**Example 2**
```
Input:  s = "wordgoodgoodgoodbestword", words = ["word","good","best","word"]
Output: []
```

**Example 3**
```
Input:  s = "barfoofoobarthefoobarman", words = ["bar","foo","the"]
Output: [6,9,12]
```

**Constraints**
- `1 <= s.length <= 10⁴`
- `1 <= words.length <= 5000`
- `1 <= words[i].length <= 30`
- `s` and `words[i]` consist of lowercase English letters.

---

## Company Frequency

| Company   | Frequency       | Last Reported |
|-----------|-----------------|---------------|
| Amazon    | ★★★★☆ High      | 2024          |
| Google    | ★★★★☆ High      | 2024          |
| Microsoft | ★★★☆☆ Medium    | 2024          |
| Bloomberg | ★★★☆☆ Medium    | 2023          |
| Meta      | ★★★☆☆ Medium    | 2023          |
| Apple     | ★★☆☆☆ Low       | 2022          |

> ⚠️ Frequency data is crowd-sourced from LeetCode Discuss, Glassdoor, and community interview reports. Treat as a signal, not a guarantee.

---

## DSA Concepts Used

- **Sliding Window** — the optimal approach uses a word-granularity sliding window, shrinking from the left when a word is over-counted. → see [`/dsa/sliding_window.md`](/dsa/sliding_window.md)
- **Hash Map (Frequency Map)** — both approaches maintain a `need` map of required word frequencies and a `have` map for the current window.
- **Fixed-Length Word Chunking** — since all words have the same length, the problem can be reduced to a fixed-width character sliding window at word granularity.

---

## Approaches Overview

| # | Approach | Time | Space | When to use |
|---|----------|------|-------|-------------|
| 1 | Brute Force | O(n·k·wordLen) | O(k) | Clear to reason about; passes within constraints |
| 2 | Sliding Window ✅ | O(n) | O(k) | Optimal; standard interview answer for Hard problems |

n = len(s), k = len(words), wordLen = len(words[0]).

---

## Approach 1 — Brute Force

### Intuition
A valid window has exactly `k * wordLen` characters. Check every possible starting position `i` in `s`: extract `k` consecutive words of length `wordLen` from `s[i..]` and verify they form a permutation of `words` using a frequency map.

### Algorithm
1. Build `need` = word frequency map from `words`.
2. `total = k * wordLen`.
3. For `i = 0` to `len(s) - total`:
   - Build `have` by extracting words `s[i+j*wordLen .. i+j*wordLen+wordLen-1]` for `j = 0..k-1`.
   - If any word is not in `need` or appears too many times, break.
   - If all k words matched: append `i` to result.

### Complexity
- **Time:** O(n·k·wordLen) — n windows × k word extractions × wordLen for map lookup.
- **Space:** O(k) — the frequency maps.

### Code
```go
func bruteForce(s string, words []string) []int {
    if len(s) == 0 || len(words) == 0 {
        return nil
    }
    wordLen := len(words[0])
    k := len(words)
    total := wordLen * k

    // build required frequency map
    need := make(map[string]int)
    for _, w := range words {
        need[w]++
    }

    var result []int
    for i := 0; i <= len(s)-total; i++ {
        have := make(map[string]int)
        j := 0
        for j < k {
            word := s[i+j*wordLen : i+j*wordLen+wordLen] // extract j-th word in window
            if need[word] == 0 {                           // word not in words list
                break
            }
            have[word]++
            if have[word] > need[word] { // word appears too many times
                break
            }
            j++
        }
        if j == k { // all k words matched
            result = append(result, i)
        }
    }
    return result
}
```

### Dry Run — `s="barfoothefoobarman"`, `words=["foo","bar"]`

`wordLen=3, k=2, total=6, need={foo:1,bar:1}`. Windows of length 6, start `i` in `[0 .. 12]`.

| i | window `s[i..i+5]` | words extracted | have valid? | record? |
|---|--------------------|-----------------|-------------|---------|
| 0 | `barfoo` | bar, foo | both in need, counts ok → j=k | **append 0** |
| 1 | `arfoot` | arf → not in need | break | no |
| 2 | `rfooth` | rfo → not in need | break | no |
| 3 | `foothe` | foo(ok), the → not in need | break | no |
| 4 | `oothef` | oot → not in need | break | no |
| 5 | `othefo` | oth → not in need | break | no |
| 6 | `thefoo` | the → not in need | break | no |
| 7 | `hefoob` | hef → not in need | break | no |
| 8 | `efooba` | efo → not in need | break | no |
| 9 | `foobar` | foo, bar | both in need, counts ok → j=k | **append 9** |
| 10 | `oobarm` | oob → not in need | break | no |
| 11 | `obarma` | oba → not in need | break | no |
| 12 | `barman` | bar(ok), man → not in need | break | no |

Result: `[0, 9]` ✓

---

## Approach 2 — Sliding Window per Offset (Recommended ✅)

### Intuition
Words are fixed-length, so the search space splits into `wordLen` independent "grids":
- Grid 0: positions 0, wordLen, 2·wordLen, …
- Grid 1: positions 1, 1+wordLen, 1+2·wordLen, …
- …

Within each grid, run a single sliding window at word granularity. Maintain:
- `have`: current word counts in the window.
- `matched`: number of words that have been correctly counted (not over/under).
- `left`: left boundary of the window (word-aligned within the grid).

Advance the right pointer word by word:
- **Unknown word:** reset the window.
- **Word goes over-count:** shrink from the left until count is correct.
- **matched == k:** record `left`; shrink left by one word.

### Algorithm
```
for offset = 0 to wordLen-1:
  reset have, matched, left=offset
  for i = offset to n-wordLen step wordLen:
    word = s[i..i+wordLen-1]
    if word not in need: reset; continue
    have[word]++
    if have[word] <= need[word]: matched++
    while have[word] > need[word]: shrink from left
    if matched == k: record left; shrink left by one word
```

### Complexity
- **Time:** O(n) — each character is examined at most twice per offset grid; there are wordLen grids; total = O(wordLen × n/wordLen) = O(n).
- **Space:** O(k).

### Code
```go
func slidingWindow(s string, words []string) []int {
    wordLen, k, n := len(words[0]), len(words), len(s)
    need := make(map[string]int)
    for _, w := range words { need[w]++ }
    var result []int
    for offset := 0; offset < wordLen; offset++ {
        have := make(map[string]int); matched, left := 0, offset
        for i := offset; i <= n-wordLen; i += wordLen {
            word := s[i : i+wordLen]
            if _, ok := need[word]; !ok {
                have = make(map[string]int); matched = 0; left = i+wordLen; continue
            }
            have[word]++
            if have[word] <= need[word] { matched++ }
            for have[word] > need[word] {
                lw := s[left : left+wordLen]
                have[lw]--
                if have[lw] < need[lw] { matched-- }
                left += wordLen
            }
            if matched == k {
                result = append(result, left)
                lw := s[left : left+wordLen]; have[lw]--
                if have[lw] < need[lw] { matched-- }
                left += wordLen
            }
        }
    }
    return result
}
```

### Dry Run — `s="barfoothefoobarman"`, `words=["foo","bar"]`

`wordLen=3, k=2, total=6, need={foo:1,bar:1}`

**Grid offset=0** (positions 0,3,6,9,12,15):
```
left=0, matched=0
i=0: "bar" ∈ need, have={bar:1}, matched=1
i=3: "foo" ∈ need, have={bar:1,foo:1}, matched=2=k → record 0. shrink: left=3
i=6: "the" ∉ need → reset, left=9
i=9: "foo" ∈ need, have={foo:1}, matched=1
i=12:"bar" ∈ need, have={foo:1,bar:1}, matched=2=k → record 9. shrink: left=12
i=15:"man" ∉ need → reset
```
Grid offset=1, offset=2: no matches.

Result: [0, 9] ✓

---

## Key Takeaways

- **wordLen grids = the core insight** — because words have fixed length, the search grid is determined by `offset % wordLen`. Running `wordLen` separate sliding windows gives an O(n) total.
- **Three window events**: (1) unknown word → hard reset; (2) over-count word → shrink from left; (3) matched == k → record and shrink by one word.
- **Don't re-make the `need` map inside the loop** — `need` is constant; only `have` changes. Re-creating `have` on reset is necessary but cheap.
- **This is LeetCode's hardest sliding window** — the fixed-word-length structure is the key observation. Without it, the problem would require a suffix automaton.

---

## Related Problems

- LeetCode #3 — Longest Substring Without Repeating Characters (variable sliding window)
- LeetCode #76 — Minimum Window Substring (sliding window with frequency maps)
- LeetCode #567 — Permutation in String (fixed-size window, same multiset check)
