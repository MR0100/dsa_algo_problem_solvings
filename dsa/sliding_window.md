# Sliding Window

## What it is

Sliding Window is a technique for problems that ask about **contiguous subarrays or substrings**. Instead of re-examining overlapping regions from scratch, you maintain a window `[left, right]` and slide it across the input — expanding or contracting based on a condition.

It converts many O(n²) or O(n³) brute-force solutions into O(n).

---

## Two variants

### 1. Fixed-size window
The window always has exactly `k` elements. Slide one step at a time: add the new right element, remove the old left element.

```go
// Sum of every subarray of size k
windowSum := 0
for i := 0; i < k; i++ { windowSum += nums[i] }
maxSum := windowSum
for i := k; i < len(nums); i++ {
    windowSum += nums[i] - nums[i-k]  // slide: add right, remove left
    if windowSum > maxSum { maxSum = windowSum }
}
```

### 2. Variable-size window
The window grows and shrinks based on a validity condition. `right` expands freely; `left` advances to restore validity when the window becomes invalid.

```go
left := 0
for right := 0; right < len(s); right++ {
    // include s[right] into the window
    for !windowIsValid() {
        // exclude s[left], shrink
        left++
    }
    // window [left, right] is now valid — update answer
}
```

---

## When to recognise it

| Problem signals | Window type |
|-----------------|-------------|
| "Longest/shortest subarray/substring satisfying X" | Variable |
| "Maximum/minimum sum of subarray of size k" | Fixed |
| "Number of subarrays satisfying X" | Variable |
| "Find an anagram / permutation of pattern in string" | Fixed (size = pattern length) |
| "At most K distinct characters / elements" | Variable |

Key phrase: **contiguous** — if the problem allows non-contiguous elements, it is likely not a sliding window.

---

## Common patterns

### Longest window (maximise)
```go
// Expand right freely; shrink left only when invalid
left := 0
for right := 0; right < n; right++ {
    add(s[right])
    for !valid() { remove(s[left]); left++ }
    maxLen = max(maxLen, right-left+1)
}
```

### Shortest window (minimise)
```go
// Expand right until valid; shrink left as long as still valid
left := 0
for right := 0; right < n; right++ {
    add(s[right])
    for valid() {
        minLen = min(minLen, right-left+1)
        remove(s[left]); left++
    }
}
```

### Fixed-size window
```go
// Initialise first window, then slide
for i := k; i < n; i++ {
    add(nums[i]); remove(nums[i-k])
    update answer
}
```

---

## Complexity

| | Time | Space |
|-|------|-------|
| Fixed window | O(n) | O(1) or O(k) |
| Variable window | O(n) | O(window size) |

`left` and `right` together advance at most `2n` steps, so the inner loop does not make the total complexity O(n²).

---

## Common pitfalls

1. **Moving left backwards** — if you store last-seen indices in a map, always guard `lastSeen[ch] >= left` before jumping left; otherwise a stale entry before the window causes left to retreat.
2. **Off-by-one in window length** — window length is `right - left + 1` (both endpoints inclusive).
3. **Shrinking condition** — be precise about what "valid" means. For "at most K distinct", shrink when `distinct > K`; for "no duplicates", shrink when the incoming char is already in the set.
4. **Updating the answer** — for maximum-window problems, update after restoring validity. For minimum-window problems, update before shrinking.

---

## Problems in this repo that use Sliding Window

- [0003 — Longest Substring Without Repeating Characters](/0003_longest_substring_without_repeating_characters/README.md)
