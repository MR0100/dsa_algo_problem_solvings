# 0125 — Valid Palindrome

> LeetCode #125 · Difficulty: Easy
> **Categories:** Two Pointers, String

---

## Problem Statement

A phrase is a **palindrome** if, after converting all uppercase letters into lowercase letters and removing all non-alphanumeric characters, it reads the same forward and backward. Alphanumeric characters include letters and numbers.

Given a string `s`, return `true` if it is a palindrome, or `false` otherwise.

**Example 1:**
```
Input: s = "A man, a plan, a canal: Panama"
Output: true
Explanation: "amanaplanacanalpanama" is a palindrome.
```

**Example 2:**
```
Input: s = "race a car"
Output: false
Explanation: "raceacar" is not a palindrome.
```

**Example 3:**
```
Input: s = " "
Output: true
Explanation: After removing non-alphanumeric chars, the string is empty, which is a palindrome.
```

**Constraints:**
- `1 <= s.length <= 2 * 10^5`
- `s` consists only of printable ASCII characters.

---

## Company Frequency

| Company   | Frequency       | Last Reported |
|-----------|-----------------|---------------|
| Amazon    | ★★★★★ Very High | 2024          |
| Microsoft | ★★★★☆ High      | 2024          |
| Facebook  | ★★★★☆ High      | 2024          |
| Bloomberg | ★★★★☆ High      | 2024          |
| Google    | ★★★☆☆ Medium    | 2023          |

> ⚠️ Frequency data is crowd-sourced from LeetCode Discuss, Glassdoor, and community interview reports. Treat as a signal, not a guarantee.

---

## DSA Concepts Used

- **Two Pointers** — converge from both ends, skip non-alphanumeric characters

---

## Approaches Overview

| # | Approach                   | Time | Space | When to use         |
|---|----------------------------|------|-------|---------------------|
| 1 | Filter String + Two Ptrs   | O(n) | O(n)  | Easy to understand  |
| 2 | In-Place Two Pointers      | O(n) | O(1)  | Follow-up / optimal |

---

## Approach 1 — Filter String then Two Pointers

### Intuition
Build a filtered lowercase alphanumeric string, then use two pointers to check palindrome.

### Complexity
- **Time:** O(n)
- **Space:** O(n) — filtered string.

### Code
```go
func isPalindrome(s string) bool {
    var sb strings.Builder
    for _, ch := range s {
        if unicode.IsLetter(ch) || unicode.IsDigit(ch) { sb.WriteRune(unicode.ToLower(ch)) }
    }
    filtered := sb.String()
    for i, j := 0, len(filtered)-1; i < j; i, j = i+1, j-1 {
        if filtered[i] != filtered[j] { return false }
    }
    return true
}
```

### Dry Run
`"A man, a plan, a canal: Panama"` → `"amanaplanacanalpanama"`.

Two pointers: a==a, m==m, a==a, n==n, a==a, p==p, l==l, a==a, n==n, a==a → true ✓

---

## Approach 2 — In-Place Two Pointers (O(1) Space)

### Intuition
Use `l, r` pointers on the original string. Skip non-alphanumeric from each side, then compare lowercased characters.

### Algorithm
1. `l=0, r=len-1`.
2. While `l < r`:
   - Skip non-alphanumeric at `l`, `r`.
   - Compare `toLower(s[l])` with `toLower(s[r])`.
   - If different: return false.
   - `l++`, `r--`.
3. Return true.

### Complexity
- **Time:** O(n)
- **Space:** O(1)

### Code
```go
func isPalindromeO1(s string) bool {
    l, r := 0, len(s)-1
    for l < r {
        for l < r && !isAlphanumeric(s[l]) { l++ }
        for l < r && !isAlphanumeric(s[r]) { r-- }
        if toLower(s[l]) != toLower(s[r]) { return false }
        l++; r--
    }
    return true
}
```

### Dry Run
`"race a car"` → filtered equivalent `"raceacar"`:

- r==r ✓, a==a ✓, c==c ✓, e≠a → return false ✓

---

## Key Takeaways
- Skip non-alphanumeric at both ends before comparing.
- Empty string (all non-alphanumeric removed) is a palindrome.
- Case-insensitive: compare lowercase versions only.

---

## Related Problems
- LeetCode #680 — Valid Palindrome II (at most one deletion)
- LeetCode #234 — Palindrome Linked List
- LeetCode #9 — Palindrome Number
