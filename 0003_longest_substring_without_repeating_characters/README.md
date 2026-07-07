# 0003 — Longest Substring Without Repeating Characters

> LeetCode #3 · Difficulty: Medium
> **Categories:** String, Sliding Window, Hash Map, Two Pointers

---

## Problem Statement

Given a string `s`, find the length of the **longest substring** without duplicate characters.

**Example 1**
```
Input:  s = "abcabcbb"
Output: 3
Explanation: The answer is "abc", with the length of 3.
```

**Example 2**
```
Input:  s = "bbbbb"
Output: 1
Explanation: The answer is "b", with the length of 1.
```

**Example 3**
```
Input:  s = "pwwkew"
Output: 3
Explanation: The answer is "wke", with the length of 3.
Notice that the answer must be a substring — "pwke" is a subsequence, not a substring.
```

**Constraints**
- `0 <= s.length <= 5 × 10⁴`
- `s` consists of English letters, digits, symbols and spaces.

---

## Company Frequency

| Company   | Frequency       | Last Reported |
|-----------|-----------------|---------------|
| Amazon    | ★★★★★ Very High | 2024          |
| Google    | ★★★★★ Very High | 2024          |
| Meta      | ★★★★★ Very High | 2024          |
| Microsoft | ★★★★☆ High      | 2024          |
| Bloomberg | ★★★★☆ High      | 2024          |
| Apple     | ★★★★☆ High      | 2023          |
| Uber      | ★★★☆☆ Medium    | 2023          |
| Adobe     | ★★★☆☆ Medium    | 2023          |
| Netflix   | ★★★☆☆ Medium    | 2023          |
| LinkedIn  | ★★☆☆☆ Low       | 2022          |
| Salesforce| ★★☆☆☆ Low       | 2022          |

> ⚠️ Frequency data is crowd-sourced from LeetCode Discuss, Glassdoor, and community
> interview reports. Treat as a signal, not a guarantee.

---

## DSA Concepts Used

- **Sliding Window** — a variable-width window `[left, right]` expands right and contracts left to maintain the invariant that all characters inside are unique. → see [`/dsa/sliding_window.md`](/dsa/sliding_window.md)
- **Hash Map** — maps each character to its most-recently-seen index, enabling O(1) detection of duplicates and O(1) jump of the left pointer. → see [`/dsa/hash_map.md`](/dsa/hash_map.md)
- **Two Pointers** — `left` and `right` both move forward (same direction, different speeds); together they define the current candidate window. → see [`/dsa/two_pointers.md`](/dsa/two_pointers.md)

---

## Approaches Overview

| # | Approach | Time | Space | When to use |
|---|----------|------|-------|-------------|
| 1 | Brute Force | O(n³) | O(min(n,a)) | Never in practice; baseline understanding |
| 2 | Sliding Window + HashSet | O(n) | O(min(n,a)) | Clear to write; left shrinks one step at a time |
| 3 | Sliding Window + HashMap ✅ | O(n) | O(min(n,a)) | General optimal — O(1) left jump |
| 4 | Sliding Window + Array | O(n) | O(1) | ASCII input; fastest constant factor |

---

## Approach 1 — Brute Force

### Intuition
Generate every possible substring and check each one for uniqueness. Keep the length of the longest valid substring seen.

### Algorithm
1. For each start index `i` from `0` to `n-1`:
2. Maintain a `seen` set; for each `j` from `i` to `n-1`:
   - If `s[j]` is already in `seen` → break (all longer substrings starting at `i` are also invalid).
   - Otherwise add `s[j]` to `seen` and update `maxLen`.

### Complexity
- **Time:** O(n³) — O(n²) pairs × O(n) for building/checking the set.
- **Space:** O(min(n, a)) — the set, where `a` is the alphabet size.

### Code
```go
func bruteForce(s string) int {
    n := len(s)
    maxLen := 0
    for i := 0; i < n; i++ {
        seen := make(map[byte]bool)
        for j := i; j < n; j++ {
            if seen[s[j]] {
                break
            }
            seen[s[j]] = true
            if j-i+1 > maxLen {
                maxLen = j - i + 1
            }
        }
    }
    return maxLen
}
```

### Dry Run — `s = "abcabcbb"`

| i | j | window | seen | valid? | maxLen |
|---|---|--------|------|--------|--------|
| 0 | 0 | "a" | {a} | ✅ | 1 |
| 0 | 1 | "ab" | {a,b} | ✅ | 2 |
| 0 | 2 | "abc" | {a,b,c} | ✅ | 3 |
| 0 | 3 | "abca" | — | ❌ 'a' dup → break | 3 |
| 1 | 1 | "b" | {b} | ✅ | 3 |
| 1 | 2 | "bc" | {b,c} | ✅ | 3 |
| 1 | 3 | "bca" | {b,c,a} | ✅ | 3 |
| 1 | 4 | "bcab" | — | ❌ 'b' dup → break | 3 |
| ... | ... | ... | ... | ... | 3 |

Final answer: **3**

---

## Approach 2 — Sliding Window with HashSet

### Intuition
Instead of restarting from scratch at each `i`, maintain a sliding window using a set that holds exactly the characters in `[left, right]`. Expand `right` freely; when a duplicate enters, shrink `left` one step at a time until the duplicate is removed.

### Algorithm
1. `left = 0`, `charSet = {}`, `maxLen = 0`.
2. For `right = 0` to `n-1`:
   - While `s[right]` is in `charSet`: remove `s[left]` from set, `left++`.
   - Add `s[right]` to set.
   - `maxLen = max(maxLen, right - left + 1)`.
3. Return `maxLen`.

### Complexity
- **Time:** O(n) — `left` moves right at most n times total; `right` moves n times. Total = 2n = O(n).
- **Space:** O(min(n, a)) — set holds at most one entry per unique character in the window.

### Code
```go
func slidingWindowSet(s string) int {
    charSet := make(map[byte]bool)
    left := 0
    maxLen := 0
    for right := 0; right < len(s); right++ {
        for charSet[s[right]] {
            delete(charSet, s[left])
            left++
        }
        charSet[s[right]] = true
        if right-left+1 > maxLen {
            maxLen = right - left + 1
        }
    }
    return maxLen
}
```

### Dry Run — `s = "abcabcbb"`

| right | s[right] | action | window | maxLen |
|-------|----------|--------|--------|--------|
| 0 | 'a' | add | [0,0]="a" | 1 |
| 1 | 'b' | add | [0,1]="ab" | 2 |
| 2 | 'c' | add | [0,2]="abc" | 3 |
| 3 | 'a' | 'a'∈set → remove s[0]='a', left=1; add 'a' | [1,3]="bca" | 3 |
| 4 | 'b' | 'b'∈set → remove s[1]='b', left=2; add 'b' | [2,4]="cab" | 3 |
| 5 | 'c' | 'c'∈set → remove s[2]='c', left=3; add 'c' | [3,5]="abc" | 3 |
| 6 | 'b' | 'b'∈set → remove s[3]='a',left=4; 'b'∈set → remove s[4]='b',left=5; add 'b' | [5,6]="cb" | 3 |
| 7 | 'b' | 'b'∈set → remove s[5]='c',left=6; 'b'∈set → remove s[6]='b',left=7; add 'b' | [7,7]="b" | 3 |

Final: **3**

---

## Approach 3 — Sliding Window with HashMap, Jump Left (Optimal)

### Intuition
The inner while-loop in Approach 2 can move `left` many times in one step. Instead, store the **last-seen index** of each character. When a duplicate `s[right]` is found inside the window, jump `left` directly to `lastSeen[s[right]] + 1` in a single O(1) operation.

The critical guard: only jump if `lastSeen[ch] >= left`. A character seen before the window started is irrelevant — jumping to its position would move `left` backward, which is wrong.

### Algorithm
1. `lastSeen = {}`, `left = 0`, `maxLen = 0`.
2. For `right = 0` to `n-1`:
   - `ch = s[right]`
   - If `ch` in `lastSeen` AND `lastSeen[ch] >= left`: `left = lastSeen[ch] + 1`.
   - `lastSeen[ch] = right`.
   - `maxLen = max(maxLen, right - left + 1)`.
3. Return `maxLen`.

### Complexity
- **Time:** O(n) — single pass; each character visited exactly once.
- **Space:** O(min(n, a)) — map has at most one entry per distinct character.

### Code
```go
func slidingWindowMap(s string) int {
    lastSeen := make(map[byte]int)
    left := 0
    maxLen := 0
    for right := 0; right < len(s); right++ {
        ch := s[right]
        if idx, ok := lastSeen[ch]; ok && idx >= left {
            left = idx + 1
        }
        lastSeen[ch] = right
        if right-left+1 > maxLen {
            maxLen = right - left + 1
        }
    }
    return maxLen
}
```

### Dry Run — `s = "abcabcbb"`

| right | ch | lastSeen[ch] | left before | jump? | left after | window len | maxLen |
|-------|----|--------------|-------------|-------|------------|------------|--------|
| 0 | a | — | 0 | ❌ | 0 | 1 | 1 |
| 1 | b | — | 0 | ❌ | 0 | 2 | 2 |
| 2 | c | — | 0 | ❌ | 0 | 3 | 3 |
| 3 | a | 0 ≥ 0 | 0 | ✅ | 1 | 3 | 3 |
| 4 | b | 1 ≥ 1 | 1 | ✅ | 2 | 3 | 3 |
| 5 | c | 2 ≥ 2 | 2 | ✅ | 3 | 3 | 3 |
| 6 | b | 4 ≥ 3 | 3 | ✅ | 5 | 2 | 3 |
| 7 | b | 6 ≥ 5 | 5 | ✅ | 7 | 1 | 3 |

Final: **3** — one clean pass, no inner loop.

---

## Approach 4 — Sliding Window with Fixed Array

### Intuition
Identical to Approach 3 but replaces the hash map with a 128-element `int` array (one slot per ASCII character value). Array indexing (`lastSeen[ch]`) has a smaller constant than map lookup (`map[byte]int`): no hashing, no collision handling, everything in contiguous memory.

Initialize every slot to `-1` to mean "not seen". No need for an `ok` check — `lastSeen[ch] >= left` (with `left >= 0`) safely returns false when the character hasn't been seen (slot = -1 < 0 ≤ left).

### Algorithm
Same as Approach 3, replacing `map[byte]int` with `[128]int` initialized to `-1`.

### Complexity
- **Time:** O(n) — same as Approach 3.
- **Space:** O(1) — the 128-slot array is constant-size regardless of input.

### Code
```go
func slidingWindowArray(s string) int {
    var lastSeen [128]int
    for i := range lastSeen { lastSeen[i] = -1 }
    left := 0
    maxLen := 0
    for right := 0; right < len(s); right++ {
        ch := s[right]
        if lastSeen[ch] >= left {
            left = lastSeen[ch] + 1
        }
        lastSeen[ch] = right
        if right-left+1 > maxLen {
            maxLen = right - left + 1
        }
    }
    return maxLen
}
```

### Dry Run — `s = "bbbbb"` (all same character)

| right | ch | lastSeen['b'] | left | jump to | window | maxLen |
|-------|----|---------------|------|---------|--------|--------|
| 0 | b | -1 < 0 | 0 | — | 1 | 1 |
| 1 | b | 0 ≥ 0 | 0 | 1 | 1 | 1 |
| 2 | b | 1 ≥ 1 | 1 | 2 | 1 | 1 |
| 3 | b | 2 ≥ 2 | 2 | 3 | 1 | 1 |
| 4 | b | 3 ≥ 3 | 3 | 4 | 1 | 1 |

Final: **1** ✓

---

## Key Takeaways

- **Sliding window pattern** — any problem asking for the longest/shortest contiguous subarray/substring satisfying a constraint is a sliding window candidate. Expand right, shrink left, maintain the invariant.
- **Jump vs shrink** — replacing the inner while-loop (shrink one step) with a direct jump (using a hash map for last-seen index) is a classic O(n) optimisation. The guard `lastSeen[ch] >= left` is easy to forget but critical.
- **Array vs map** — for fixed-size character sets (ASCII 128, extended 256), a fixed array is strictly faster than a hash map due to cache locality and no hashing overhead. For Unicode, a map is necessary.
- **The `>= left` guard** — without this, `left` could move backwards if a character appears both before and inside the current window. Always check that the cached index is actually inside the current window.

---

## Related Problems

- LeetCode #159 — Longest Substring with At Most Two Distinct Characters (window with count ≤ 2)
- LeetCode #340 — Longest Substring with At Most K Distinct Characters (generalised)
- LeetCode #424 — Longest Repeating Character Replacement (sliding window with a twist)
- LeetCode #76 — Minimum Window Substring (minimum window containing all chars of t)
- LeetCode #567 — Permutation in String (fixed-size window anagram check)
