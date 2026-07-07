# 0028 — Find the Index of the First Occurrence in a String

> LeetCode #28 · Difficulty: Easy
> **Categories:** String, Two Pointers, String Matching

---

## Problem Statement

Given two strings `haystack` and `needle`, return the index of the first occurrence of `needle` in `haystack`, or `-1` if `needle` is not part of `haystack`.

**Example 1**
```
Input:  haystack = "sadbutsad", needle = "sad"
Output: 0
Explanation: "sad" occurs at index 0 and 6. The first occurrence is at index 0.
```

**Example 2**
```
Input:  haystack = "leetcode", needle = "leeto"
Output: -1
Explanation: "leeto" did not occur in "leetcode", so return -1.
```

**Constraints**
- `1 <= haystack.length, needle.length <= 10⁴`
- `haystack` and `needle` consist of only lowercase English letters.

---

## Company Frequency

| Company   | Frequency       | Last Reported |
|-----------|-----------------|---------------|
| Amazon    | ★★★★☆ High      | 2024          |
| Google    | ★★★★☆ High      | 2024          |
| Microsoft | ★★★★☆ High      | 2024          |
| Meta      | ★★★☆☆ Medium    | 2023          |
| Bloomberg | ★★★☆☆ Medium    | 2023          |
| Apple     | ★★☆☆☆ Low       | 2022          |

> ⚠️ Frequency data is crowd-sourced from LeetCode Discuss, Glassdoor, and community interview reports. Treat as a signal, not a guarantee.

---

## DSA Concepts Used

- **String Matching** — the fundamental substring search problem; brute force is O(n·m); KMP achieves O(n+m).
- **KMP Failure Function (LPS array)** — the key to O(n+m): the LPS (Longest Proper Prefix which is also Suffix) array lets us skip re-comparing characters after a mismatch.

---

## Approaches Overview

| # | Approach | Time | Space | When to use |
|---|----------|------|-------|-------------|
| 1 | Brute Force | O(n·m) | O(1) | Short strings or interview warm-up |
| 2 | Standard Library | O(n) avg | O(1) | When stdlib use is acceptable |
| 3 | KMP ✅ | O(n+m) | O(m) | Interview-optimal; shows deep string knowledge |

n = len(haystack), m = len(needle).

---

## Approach 1 — Brute Force (Sliding Window)

### Intuition
Try every position `i` in haystack as a potential start. At each position, compare haystack[i..i+m-1] character by character with needle. Return `i` on the first full match.

### Algorithm
```
for i = 0 to n-m:
  j = 0
  while j < m and haystack[i+j] == needle[j]: j++
  if j == m: return i
return -1
```

### Complexity
- **Time:** O(n·m) — worst case: "aaaa...ab" / "aab" causes m comparisons at each of n-m positions.
- **Space:** O(1).

### Code
```go
func bruteForce(haystack, needle string) int {
    n, m := len(haystack), len(needle)
    if m == 0 {
        return 0
    }
    for i := 0; i <= n-m; i++ {
        j := 0
        for j < m && haystack[i+j] == needle[j] { // compare character by character
            j++
        }
        if j == m { // all m characters matched
            return i
        }
    }
    return -1
}
```

### Dry Run — `haystack="sadbutsad"`, `needle="sad"`
```
i=0: s==s, a==a, d==d → j=3=m → return 0 ✓
```

---

## Approach 2 — Standard Library (`strings.Index`)

### Intuition
Go's `strings.Index` implements an optimised string search. Useful when showing awareness of stdlib without needing to implement the algorithm.

### Complexity
- **Time:** O(n) average.
- **Space:** O(1).

### Code
```go
func useStdlib(haystack, needle string) int {
    return strings.Index(haystack, needle)
}
```

### Dry Run — `haystack="sadbutsad"`, `needle="sad"`
`strings.Index` scans internally; here we trace the equivalent search it performs.

| step | window in haystack | compare vs `sad` | result |
|------|--------------------|------------------|--------|
| start 0 | `sad`butsad | `sad`==`sad` | full match |
| return | — | — | index `0` |

`strings.Index("sadbutsad", "sad")` returns `0` on the first matching window. ✓

---

## Approach 3 — KMP (Knuth-Morris-Pratt) — Recommended ✅

### Intuition
After a mismatch at `needle[j]`, we've already matched `needle[0..j-1]` against some suffix of the already-scanned haystack. Instead of restarting from scratch, we can jump `j` back to the longest prefix of `needle` that is also a suffix of `needle[0..j-1]`. This value is `lps[j-1]`.

Example: `needle = "AAACAAAA"`:
```
lps = [0, 1, 2, 0, 1, 2, 3, 3]
```
If we matched `needle[0..5]="AAACAA"` and hit a mismatch at j=6, lps[5]=2, so we jump j to 2 (already matched "AA") instead of restarting from j=0.

### Algorithm
**Build LPS:**
```
lps[0] = 0; length = 0; i = 1
while i < m:
  if needle[i] == needle[length]: length++; lps[i]=length; i++
  elif length != 0: length = lps[length-1]  // try shorter prefix
  else: lps[i]=0; i++
```

**Search:**
```
i=0, j=0
while i < n:
  if haystack[i] == needle[j]: i++; j++
  if j == m: return i-m
  elif i<n and mismatch:
    if j != 0: j = lps[j-1]
    else: i++
return -1
```

### Complexity
- **Time:** O(n+m) — O(m) to build LPS; O(n) to search (i never decrements).
- **Space:** O(m) — the LPS array.

### Code
```go
func kmp(haystack, needle string) int {
    n, m := len(haystack), len(needle)
    if m == 0 { return 0 }
    lps := make([]int, m)
    length := 0
    for i := 1; i < m; {
        if needle[i] == needle[length] { length++; lps[i] = length; i++ } else
        if length != 0 { length = lps[length-1] } else { lps[i] = 0; i++ }
    }
    i, j := 0, 0
    for i < n {
        if haystack[i] == needle[j] { i++; j++ }
        if j == m { return i - m }
        if i < n && haystack[i] != needle[j] {
            if j != 0 { j = lps[j-1] } else { i++ }
        }
    }
    return -1
}
```

### Dry Run — `haystack="aabaabaaf"`, `needle="aabaaf"`
```
needle LPS: a=0,a=1,b=0,a=1,a=2,f=0 → [0,1,0,1,2,0]

i=0,j=0: a==a → i=1,j=1
i=1,j=1: a==a → i=2,j=2
i=2,j=2: b==b → i=3,j=3
i=3,j=3: a==a → i=4,j=4
i=4,j=4: a==a → i=5,j=5
i=5,j=5: b≠f → j=lps[4]=2
i=5,j=2: b==b → i=6,j=3
i=6,j=3: a==a → i=7,j=4
i=7,j=4: a==a → i=8,j=5
i=8,j=5: f==f → i=9,j=6 = m → return 9-6 = 3 ✓
```

---

## Key Takeaways

- **KMP never moves `i` backward** — this is why it's O(n+m) instead of O(n·m). Each position in haystack is visited at most once.
- **LPS = skip table** — `lps[j-1]` after a mismatch at j means: "I've verified that this prefix of needle also appears as a suffix — restart j there."
- **When needle is empty, return 0** — this matches the convention for `strstr` in C and `strings.Index` in Go.
- **This is asked at FAANG** — KMP is the expected O(n+m) answer; brute force is acceptable for easy-level, but interviewers at Google/Meta may push you to explain KMP.

---

## Related Problems

- LeetCode #214 — Shortest Palindrome (uses KMP failure function)
- LeetCode #459 — Repeated Substring Pattern (uses KMP or period detection)
- LeetCode #686 — Repeated String Match (find repetitions needed to contain needle)
