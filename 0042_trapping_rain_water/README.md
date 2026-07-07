# 0042 — Trapping Rain Water

> LeetCode #42 · Difficulty: Hard
> **Categories:** Array, Two Pointers, Dynamic Programming, Stack, Monotonic Stack

---

## Problem Statement

Given `n` non-negative integers representing an elevation map where the width of each bar is `1`, compute how much water it can trap after raining.

**Example 1**
```
Input:  height = [0,1,0,2,1,0,1,3,2,1,2,1]
Output: 6
```

**Example 2**
```
Input:  height = [4,2,0,3,2,5]
Output: 9
```

**Constraints**
- `n == height.length`
- `1 <= n <= 2 * 10⁴`
- `0 <= height[i] <= 10⁵`

---

## Company Frequency

| Company   | Frequency       | Last Reported |
|-----------|-----------------|---------------|
| Amazon    | ★★★★★ Very High | 2024          |
| Google    | ★★★★★ Very High | 2024          |
| Meta      | ★★★★★ Very High | 2024          |
| Microsoft | ★★★★★ Very High | 2024          |
| Bloomberg | ★★★★☆ High      | 2024          |
| Apple     | ★★★★☆ High      | 2023          |
| Adobe     | ★★★☆☆ Medium    | 2023          |

> ⚠️ Frequency data is crowd-sourced from LeetCode Discuss, Glassdoor, and community interview reports. Treat as a signal, not a guarantee.

---

## DSA Concepts Used

- **Two Pointers** — Approach 3 maintains left/right pointers and maxL/maxR trackers, processing whichever side has the shorter current boundary.
- **Precomputed Max Arrays** — Approach 2 uses O(n) space to precompute leftMax and rightMax.
- **Monotonic Stack** — Approach 4 uses a decreasing stack to detect valleys and compute water horizontally.

---

## Approaches Overview

| # | Approach | Time | Space | When to use |
|---|----------|------|-------|-------------|
| 1 | Brute Force | O(n²) | O(1) | TLE; never in practice |
| 2 | Precomputed Arrays | O(n) | O(n) | Easy to understand and implement |
| 3 | Two Pointers ✅ | O(n) | O(1) | Optimal; O(1) space follow-up answer |
| 4 | Monotonic Stack | O(n) | O(n) | Good alternative; horizontal water computation |

---

## Approach 1 — Brute Force

### Intuition
Water at column `i` = `min(maxLeft[i], maxRight[i]) - height[i]`. Compute maxLeft and maxRight naively with two inner loops per position.

### Complexity
- **Time:** O(n²).
- **Space:** O(1).

---

## Approach 2 — Precomputed Left/Right Max Arrays

### Intuition
Precompute `maxL[i]` = max height in `height[0..i]` and `maxR[i]` = max height in `height[i..n-1]`. Then water at i = `min(maxL[i], maxR[i]) - height[i]`.

### Complexity
- **Time:** O(n).
- **Space:** O(n).

---

## Approach 3 — Two Pointers (Recommended ✅)

### Intuition
Maintain `left` and `right` pointers, and `maxL`/`maxR` (max heights seen so far from each side).

**Key insight:** if `height[left] <= height[right]`, the water at `left` is determined by `maxL` (not by the right side, which is at least as tall). So we can safely compute water at `left` and advance `left`. Symmetrically for the right side.

This eliminates the need for the precomputed arrays by processing whichever side has the tighter constraint.

### Algorithm
```
left=0, right=n-1, maxL=0, maxR=0, total=0
while left < right:
  if height[left] <= height[right]:
    if height[left] >= maxL: maxL = height[left]
    else: total += maxL - height[left]
    left++
  else:
    if height[right] >= maxR: maxR = height[right]
    else: total += maxR - height[right]
    right--
```

### Complexity
- **Time:** O(n).
- **Space:** O(1).

### Code
```go
func twoPointers(height []int) int {
    left, right := 0, len(height)-1
    maxL, maxR, total := 0, 0, 0
    for left < right {
        if height[left] <= height[right] {
            if height[left] >= maxL { maxL = height[left] } else { total += maxL - height[left] }
            left++
        } else {
            if height[right] >= maxR { maxR = height[right] } else { total += maxR - height[right] }
            right--
        }
    }
    return total
}
```

### Dry Run — `height = [4,2,0,3,2,5]`
```
left=0,right=5, maxL=0,maxR=0
h[0]=4 <= h[5]=5: 4>=0→maxL=4; left=1
h[1]=2 <= h[5]=5: 2<4→total+=4-2=2; left=2
h[2]=0 <= h[5]=5: 0<4→total+=4-0=4; left=3   [total=6]
h[3]=3 <= h[5]=5: 3<4→total+=4-3=1; left=4   [total=7]
h[4]=2 <= h[5]=5: 2<4→total+=4-2=2; left=5   [total=9]
left=5==right=5 → stop
Result: 9 ✓
```

---

## Approach 4 — Monotonic Stack

### Intuition
Maintain a stack of indices in **decreasing height** order. When a taller bar `i` appears, pop the valley floor (the top of the stack), then compute water horizontally between the remaining stack top (left wall) and `i` (right wall).

This computes water **layer by layer horizontally**, while Approach 3 computes it **column by column vertically**.

### Complexity
- **Time:** O(n).
- **Space:** O(n).

---

## Key Takeaways

- **Two core formulas:**
  - Vertical view: `water[i] = min(maxLeft[i], maxRight[i]) - height[i]`
  - Horizontal view: `water = (minWall - valleyFloor) × width`
- **Why two-pointer works** — when `height[left] <= height[right]`, `maxL` is the binding constraint for `left` regardless of what's on the right (we know at least `height[right] >= height[left]` limits water from spilling left). We can commit the computation.
- **This is a Top-10 interview question** — trapping rain water appears at virtually every FAANG loop. Know all four approaches; emphasise two-pointer for O(1) space, monotonic stack for horizontal-layer intuition.

---

## Related Problems

- LeetCode #11 — Container With Most Water (two pointers; one container not trappi
- LeetCode #84 — Largest Rectangle in Histogram (monotonic stack, similar pattern)
- LeetCode #407 — Trapping Rain Water II (3D extension with a min-heap)
