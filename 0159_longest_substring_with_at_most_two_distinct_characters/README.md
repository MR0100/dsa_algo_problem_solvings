# 0159 — Longest Substring with At Most Two Distinct Characters

> LeetCode #159 · Difficulty: Medium (Premium)
> **Categories:** Hash Table, String, Sliding Window

---

## Problem Statement

Given a string `s`, return *the length of the longest substring that contains **at most two distinct characters***.

**Example 1:**
```
Input: s = "eceba"
Output: 3
Explanation: The substring is "ece" which its length is 3.
```

**Example 2:**
```
Input: s = "ccaabbb"
Output: 5
Explanation: The substring is "aabbb" which its length is 5.
```

**Constraints:**
- `1 <= s.length <= 10^5`
- `s` consists of English letters.

---

## Company Frequency

| Company   | Frequency       | Last Reported |
|-----------|-----------------|---------------|
| Google    | ★★★★★ Very High | 2024          |
| Facebook  | ★★★★☆ High      | 2024          |
| Amazon    | ★★★☆☆ Medium    | 2023          |
| Microsoft | ★★★☆☆ Medium    | 2023          |
| Lyft      | ★★☆☆☆ Low       | 2022          |

> ⚠️ Frequency data is crowd-sourced from LeetCode Discuss, Glassdoor, and
> community interview reports. Treat as a signal, not a guarantee.

---

## DSA Concepts Used

- **Sliding Window** — "longest substring satisfying an invariant" is the canonical grow-right / shrink-left window problem → see [`/dsa/sliding_window.md`](/dsa/sliding_window.md)
- **Hash Map** — tracks either per-character counts (Approach 2) or last-occurrence indices (Approach 3) inside the window → see [`/dsa/hash_map.md`](/dsa/hash_map.md)
- **Two Pointers** — `left` and `right` delimit the window; both only move forward, giving amortised O(n) → see [`/dsa/two_pointers.md`](/dsa/two_pointers.md)

---

## Approaches Overview

| # | Approach | Time | Space | When to use |
|---|----------|------|-------|-------------|
| 1 | Brute Force | O(n²) | O(1) | Baseline; fine for tiny inputs only |
| 2 | Sliding Window + Frequency Map | O(n) | O(1) — ≤ 3 keys | Standard answer; generalises to "at most k" |
| 3 | Sliding Window + Last-Occurrence Map (Optimal) | O(n) | O(1) — ≤ 3 keys | Same asymptotics; left pointer jumps in one step |

---

## Approach 1 — Brute Force

### Intuition
Enumerate every starting index `i` and stretch the end `j` rightwards, collecting the distinct characters of `s[i..j]` in a small set. As soon as a third distinct character shows up, every longer substring starting at `i` is invalid too, so break early and try the next start.

### Algorithm
1. `best = 0`.
2. For each `i` from `0` to `n-1`:
   1. Reset `distinct` to an empty set.
   2. For each `j` from `i` to `n-1`:
      1. Insert `s[j]` into `distinct`.
      2. If `len(distinct) > 2`, break the inner loop.
      3. Else update `best = max(best, j-i+1)`.
3. Return `best`.

### Complexity
- **Time:** O(n²) — n choices of start, each scanning up to n characters (early break helps but not asymptotically, e.g. `"ababab…"`).
- **Space:** O(1) — the set never exceeds 3 entries.

### Code
```go
func bruteForce(s string) int {
    best := 0
    for i := 0; i < len(s); i++ {
        distinct := map[byte]bool{}
        for j := i; j < len(s); j++ {
            distinct[s[j]] = true
            if len(distinct) > 2 {
                break
            }
            if j-i+1 > best {
                best = j - i + 1
            }
        }
    }
    return best
}
```

### Dry Run
Example 1: `s = "eceba"`.

| i | j sweep | distinct set evolution | valid window lengths | best |
|---|---------|------------------------|----------------------|------|
| 0 | e,c,e,b | {e} → {e,c} → {e,c} → {e,c,b} break | 1, 2, 3 | 3 |
| 1 | c,e,b | {c} → {c,e} → {c,e,b} break | 1, 2 | 3 |
| 2 | e,b,a | {e} → {e,b} → {e,b,a} break | 1, 2 | 3 |
| 3 | b,a | {b} → {b,a} | 1, 2 | 3 |
| 4 | a | {a} | 1 | 3 |

Return **3** (`"ece"`) ✓

---

## Approach 2 — Sliding Window + Frequency Map

### Intuition
Keep a window `[left..right]` whose invariant is *≤ 2 distinct characters*. Slide `right` forward absorbing one character at a time; whenever the invariant breaks (3 distinct), advance `left`, decrementing counts and removing keys that reach zero, until the invariant is restored. Because both pointers only move forward, total work is linear — the classic "grow right, shrink left" template.

### Algorithm
1. `count = {}` (char → occurrences inside the window), `left = 0`, `best = 0`.
2. For `right` from `0` to `n-1`:
   1. `count[s[right]]++`.
   2. While `len(count) > 2`:
      1. `count[s[left]]--`.
      2. If `count[s[left]] == 0`, delete the key (a distinct char left the window).
      3. `left++`.
   3. `best = max(best, right-left+1)`.
3. Return `best`.

### Complexity
- **Time:** O(n) — `right` advances n times; `left` advances at most n times over the whole run (amortised O(1) per step).
- **Space:** O(1) — the map momentarily holds at most 3 keys.

### Code
```go
func slidingWindow(s string) int {
    count := map[byte]int{}
    left, best := 0, 0
    for right := 0; right < len(s); right++ {
        count[s[right]]++
        for len(count) > 2 {
            count[s[left]]--
            if count[s[left]] == 0 {
                delete(count, s[left])
            }
            left++
        }
        if right-left+1 > best {
            best = right - left + 1
        }
    }
    return best
}
```

### Dry Run
Example 1: `s = "eceba"`.

| right | s[right] | count after add | shrink? | left | window | best |
|-------|----------|-----------------|---------|------|--------|------|
| 0 | e | {e:1} | no | 0 | `"e"` | 1 |
| 1 | c | {e:1, c:1} | no | 0 | `"ec"` | 2 |
| 2 | e | {e:2, c:1} | no | 0 | `"ece"` | **3** |
| 3 | b | {e:2, c:1, b:1} | yes: drop e (e:1), drop c (c:0 → delete) → {e:1, b:1} | 2 | `"eb"` | 3 |
| 4 | a | {e:1, b:1, a:1} | yes: drop e (0 → delete) → {b:1, a:1} | 3 | `"ba"` | 3 |

Return **3** ✓

---

## Approach 3 — Sliding Window + Last-Occurrence Map (Optimal)

### Intuition
Instead of counting occurrences, remember each window character's **last occurrence index**. When a third distinct character arrives, exactly one existing character must be evicted — necessarily the one whose last occurrence is furthest left (everything after it belongs to the other two chars). The window can then **jump** `left` directly to `thatIndex + 1` in one assignment instead of stepping. The map never exceeds 3 entries, so scanning it for the minimum is O(1). This formulation is the one that scales verbatim to "at most k distinct characters" (LeetCode #340) with an ordered map / careful bookkeeping.

### Algorithm
1. `last = {}` (char → index of most recent occurrence), `left = 0`, `best = 0`.
2. For `right` from `0` to `n-1`:
   1. `last[s[right]] = right` (insert or refresh).
   2. If `len(last) > 2`:
      1. Find the entry `(evict, minIdx)` with the smallest index among the ≤ 3 entries.
      2. `delete(last, evict)`; `left = minIdx + 1`.
   3. `best = max(best, right-left+1)`.
3. Return `best`.

### Complexity
- **Time:** O(n) — one pass; the eviction scan touches at most 3 map entries, and `left` only moves forward.
- **Space:** O(1) — at most 3 map entries at any time.

### Code
```go
func slidingWindowLastIndex(s string) int {
    last := map[byte]int{}
    left, best := 0, 0
    for right := 0; right < len(s); right++ {
        last[s[right]] = right
        if len(last) > 2 {
            minIdx := right
            var evict byte
            for c, idx := range last {
                if idx < minIdx {
                    minIdx = idx
                    evict = c
                }
            }
            delete(last, evict)
            left = minIdx + 1
        }
        if right-left+1 > best {
            best = right - left + 1
        }
    }
    return best
}
```

### Dry Run
Example 1: `s = "eceba"`.

| right | s[right] | last map after update | >2 distinct? | eviction | left | window | best |
|-------|----------|------------------------|--------------|----------|------|--------|------|
| 0 | e | {e:0} | no | — | 0 | `"e"` | 1 |
| 1 | c | {e:0, c:1} | no | — | 0 | `"ec"` | 2 |
| 2 | e | {e:2, c:1} | no | — | 0 | `"ece"` | **3** |
| 3 | b | {e:2, c:1, b:3} | yes | evict c (idx 1); left = 1+1 = 2 | 2 | `"eb"` | 3 |
| 4 | a | {e:2, b:3, a:4} | yes | evict e (idx 2); left = 2+1 = 3 | 3 | `"ba"` | 3 |

Return **3** ✓

---

## Key Takeaways
- "Longest substring/subarray with at most K of something" → **sliding window** is almost always the answer; the invariant lives in a small hash map.
- Two interchangeable bookkeeping styles: **frequency counts** (decrement to zero to detect a char leaving) vs **last-occurrence indices** (evict the leftmost-last char and jump). Both are O(n); the second moves `left` in one hop.
- Space is O(1) here only because the invariant bounds the map at k+1 = 3 entries — for general k it is O(k).
- Substring problems where a violation at `(i, j)` implies violation for all `j' > j` admit the brute-force early break — and that monotonicity is precisely why the sliding window works at all.
- This problem is the k = 2 special case of LeetCode #340; write the window code so that `2` appears once and it generalises for free.

---

## Related Problems
- LeetCode #3 — Longest Substring Without Repeating Characters (window with "all distinct" invariant)
- LeetCode #340 — Longest Substring with At Most K Distinct Characters (direct generalisation)
- LeetCode #904 — Fruit Into Baskets (identical problem in disguise, k = 2)
- LeetCode #76 — Minimum Window Substring (shrinking-window counterpart)
- LeetCode #424 — Longest Repeating Character Replacement (window with budgeted violations)
