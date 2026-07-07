# 0058 — Length of Last Word

> LeetCode #58 · Difficulty: Easy
> **Categories:** String

---

## Problem Statement

Given a string `s` consisting of words and spaces, return the **length of the last word** in the string.

A **word** is a maximal substring consisting of non-space characters only.

**Example 1**
```
Input:  s = "Hello World"
Output: 5
```

**Example 2**
```
Input:  s = "   fly me   to   the moon  "
Output: 4
```

**Example 3**
```
Input:  s = "luffy is still joyboy"
Output: 6
```

**Constraints**
- `1 <= s.length <= 10⁴`
- `s` consists of only English letters and spaces `' '`.
- There will be at least one word in `s`.

---

## Company Frequency

| Company   | Frequency       | Last Reported |
|-----------|-----------------|---------------|
| Amazon    | ★★★★☆ High      | 2024          |
| Google    | ★★★★☆ High      | 2024          |
| Bloomberg | ★★★☆☆ Medium    | 2023          |
| Microsoft | ★★★☆☆ Medium    | 2023          |
| Apple     | ★★☆☆☆ Low       | 2022          |
| Meta      | ★★☆☆☆ Low       | 2022          |

> ⚠️ Frequency data is crowd-sourced from LeetCode Discuss, Glassdoor, and community interview reports. Treat as a signal, not a guarantee.

---

## DSA Concepts Used

- **String Scanning** — trivial; scan from the end to skip trailing spaces, then count the last word.

---

## Approaches Overview

| # | Approach | Time | Space | When to use |
|---|----------|------|-------|-------------|
| 1 | Split + Take Last | O(n) | O(n) | Simplest; allocates words slice |
| 2 | Reverse Scan ✅ | O(n) worst / O(k) typical | O(1) | Optimal; no allocation |

---

## Approach 1 — Split + Take Last

### Intuition
`strings.Fields` splits on any whitespace and ignores leading/trailing spaces. The last element is the last word.

### Complexity
- **Time:** O(n).
- **Space:** O(n) — allocates the words slice.

### Code
```go
func splitCount(s string) int {
    words := strings.Fields(s)
    if len(words) == 0 {
        return 0
    }
    return len(words[len(words)-1])
}
```

### Dry Run — `s = "   fly me   to   the moon  "`

`strings.Fields(s)` splits on runs of whitespace and drops leading/trailing spaces, then we take the last element's length.

| step | value |
|------|-------|
| input `s` | `"   fly me   to   the moon  "` |
| `words = strings.Fields(s)` | `["fly", "me", "to", "the", "moon"]` |
| `len(words)` | 5 (not 0, so skip the guard) |
| last word `words[4]` | `"moon"` |
| `len(words[4])` | 4 |

Return 4 ✓

---

## Approach 2 — Reverse Scan (Recommended ✅)

### Intuition
Scan from the right:
1. Skip trailing spaces.
2. Count characters until a space or the beginning is hit.

No allocation needed.

### Algorithm
```
i = len(s) - 1
while i >= 0 and s[i] == ' ': i--
count = 0
while i >= 0 and s[i] != ' ': count++; i--
return count
```

### Complexity
- **Time:** O(n) worst case (all spaces except one leading word). O(k) typical where k = trailing_spaces + last_word_length.
- **Space:** O(1).

### Code
```go
func reverseScan(s string) int {
    i := len(s) - 1
    for i >= 0 && s[i] == ' ' { i-- }
    count := 0
    for i >= 0 && s[i] != ' ' { count++; i-- }
    return count
}
```

### Dry Run — `s = "   fly me   to   the moon  "`
```
i starts at len-1 = 27 (space)
Skip spaces: i=25 ('n' in "moon")

Count: 'n','o','o','m' → count=4, i=21 (space)

Return 4 ✓
```

---

## Key Takeaways

- **Skip trailing spaces first** — the tricky part of this problem is trailing spaces (`"Hello World  "`). Without skipping them, you'd count 0.
- **Reverse scan is always faster in practice** — avoids scanning the entire string if the last word is short and at the end.
- **`strings.Fields`** in Go handles multiple spaces and trim automatically — useful for quick solutions in non-interview contexts.

---

## Related Problems

- LeetCode #151 — Reverse Words in a String (reverse order of words; similar scanning)
- LeetCode #434 — Number of Segments in a String (count word segments)
