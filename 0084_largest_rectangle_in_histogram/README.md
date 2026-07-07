# 0084 — Largest Rectangle in Histogram

> LeetCode #84 · Difficulty: Hard
> **Categories:** Array, Stack, Monotonic Stack

---

## Problem Statement

Given an array of integers `heights` representing the histogram's bar height where the width of each bar is `1`, return the area of the largest rectangle in the histogram.

**Example 1:**
```
Input: heights = [2,1,5,6,2,3]
Output: 10
Explanation: The largest rectangle is shown in the shaded area with width = 2 and height = 5.
```

**Example 2:**
```
Input: heights = [2,4]
Output: 4
```

**Constraints:**
- `1 <= heights.length <= 10^5`
- `0 <= heights[i] <= 10^4`

---

## Company Frequency

| Company   | Frequency       | Last Reported |
|-----------|-----------------|---------------|
| Amazon    | ★★★★★ Very High | 2024          |
| Google    | ★★★★☆ High      | 2024          |
| Facebook  | ★★★★☆ High      | 2024          |
| Microsoft | ★★★☆☆ Medium    | 2024          |
| Apple     | ★★★☆☆ Medium    | 2023          |

> ⚠️ Frequency data is crowd-sourced from LeetCode Discuss, Glassdoor, and community interview reports. Treat as a signal, not a guarantee.

---

## DSA Concepts Used

- **Monotonic Stack** — maintain indices in increasing height order; pop when a shorter bar is seen to compute widths. See [`/dsa/stack.md`](/dsa/stack.md)
- **Sentinel Value** — append height 0 to force all remaining bars to be popped at the end.

---

## Approaches Overview

| # | Approach | Time | Space | When to use |
|---|----------|------|-------|-------------|
| 1 | Brute Force | O(n²) | O(1) | Understanding only |
| 2 | Monotonic Stack | O(n) | O(n) | Always — standard optimal solution |

---

## Approach 1 — Brute Force

### Intuition
For each pair of left and right boundaries (i, j), find the minimum height in `heights[i..j]`. The rectangle area is `minHeight × (j - i + 1)`. Track the maximum.

### Algorithm
1. For `i = 0` to `n-1`:
   - `minH = heights[i]`.
   - For `j = i` to `n-1`:
     - `minH = min(minH, heights[j])`.
     - `maxArea = max(maxArea, minH × (j-i+1))`.

### Complexity
- **Time:** O(n²)
- **Space:** O(1)

### Code
```go
func bruteForce(heights []int) int {
    n := len(heights)
    maxArea := 0
    for i := 0; i < n; i++ {
        minH := heights[i]
        for j := i; j < len(heights); j++ {
            if heights[j] < minH { minH = heights[j] }
            area := minH * (j - i + 1)
            if area > maxArea { maxArea = area }
        }
    }
    return maxArea
}
```

### Dry Run (heights=[2,1,5,6,2,3])

| i | j=i | j=i+1 | ... | maxArea |
|---|-----|-------|-----|---------|
| 0 | 2×1=2 | min(2,1)×2=2 | ... | 2 |
| 1 | 1×1=1 | min(1,5)×2=2 | min(1,5,6)×3=3 | min(1,5,6,2)×4=4 | min×5=5 | min×6=6 → 6 |
| 2 | 5×1=5 | min(5,6)×2=10 ← MAX | ... |
Result: 10 ✓

---

## Approach 2 — Monotonic Stack (Optimal)

### Intuition
For each bar, the maximum rectangle using that bar as the *shortest* bar extends left until a shorter bar is found and right until a shorter bar is found.

A **monotonic increasing stack** (storing indices) tracks potential left boundaries. When we encounter a bar `h` shorter than `heights[stack.top]`, the top bar has found its right boundary (current index `i`). Its left boundary is `stack[top-1] + 1` (or 0 if empty). So we can compute the area immediately.

**Sentinel:** Appending height 0 at the end ensures all bars still in the stack get processed.

### Algorithm
1. Append `0` to `heights` as a sentinel.
2. `stack = []`, `maxArea = 0`.
3. For each index `i` with height `h`:
   - While `heights[stack.top] > h`:
     - Pop `topIdx`. `height = heights[topIdx]`.
     - `width = i` if stack empty, else `i - stack.top - 1`.
     - `maxArea = max(maxArea, height × width)`.
   - Push `i`.
4. Return `maxArea`.

### Complexity
- **Time:** O(n) — each index pushed and popped at most once.
- **Space:** O(n) — stack.

### Code
```go
func monoStack(heights []int) int {
    heights = append(heights, 0)
    stack := []int{}
    maxArea := 0
    for i, h := range heights {
        for len(stack) > 0 && heights[stack[len(stack)-1]] > h {
            topIdx := stack[len(stack)-1]
            stack = stack[:len(stack)-1]
            height := heights[topIdx]
            var width int
            if len(stack) == 0 {
                width = i
            } else {
                width = i - stack[len(stack)-1] - 1
            }
            if area := height * width; area > maxArea {
                maxArea = area
            }
        }
        stack = append(stack, i)
    }
    return maxArea
}
```

### Dry Run (heights=[2,1,5,6,2,3,0])

| i | h | stack before | action | maxArea |
|---|---|-------------|--------|---------|
| 0 | 2 | [] | push 0 | 0 |
| 1 | 1 | [0] | h(1)<h(0)=2 → pop 0: height=2, width=1(stack empty), area=2 | 2 |
| | | [] | push 1 | |
| 2 | 5 | [1] | push 2 | 2 |
| 3 | 6 | [1,2] | push 3 | 2 |
| 4 | 2 | [1,2,3] | h(2)<h(3)=6 → pop 3: height=6, width=4-2-1=1, area=6 | 6 |
| | | [1,2] | h(2)<h(2)=5 → pop 2: height=5, width=4-1-1=2, area=10 | 10 |
| | | [1] | h(2)>h(1)=1 → push 4 | |
| 5 | 3 | [1,4] | push 5 | 10 |
| 6 | 0 | [1,4,5] | pop 5: height=3, width=6-4-1=1, area=3 | 10 |
| | | [1,4] | pop 4: height=2, width=6-1-1=4, area=8 | 10 |
| | | [1] | pop 1: height=1, width=6(stack empty), area=6 | 10 |
| | | [] | push 6 | 10 |

Final: 10 ✓

---

## Key Takeaways
- The stack stores indices, not heights — needed to compute widths.
- When popping top with index `t`, the right boundary is current `i` and the left boundary is `stack.top + 1` (or 0 if empty).
- The sentinel 0 at the end flushes all remaining bars without a special post-loop.
- This same helper is used in #85 (Maximal Rectangle) to solve each row's histogram.

---

## Related Problems
- LeetCode #85 — Maximal Rectangle (extend #84 to 2D matrix using row histograms)
- LeetCode #42 — Trapping Rain Water (monotonic stack on histogram, complementary problem)
- LeetCode #316 — Remove Duplicate Letters (monotonic stack for lexicographic order)
