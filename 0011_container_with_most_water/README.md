# 0011 — Container With Most Water

> LeetCode #11 · Difficulty: Medium
> **Categories:** Array, Two Pointers, Greedy

---

## Problem Statement

You are given an integer array `height` of length `n`. There are `n` vertical lines drawn such that the two endpoints of the `i`-th line are `(i, 0)` and `(i, height[i])`.

Find two lines that together with the x-axis form a container such that the container contains the most water.

Return the maximum amount of water a container can store.

**Notice** that you may not slant the container.

**Example 1**
```
Input:  height = [1,8,6,2,5,4,8,3,7]
Output: 49
Explanation: The vertical lines at index 1 and 8 form the container.
             min(8,7) * (8-1) = 7*7 = 49.
```

**Example 2**
```
Input:  height = [1,1]
Output: 1
```

**Constraints**
- `n == height.length`
- `2 <= n <= 10⁵`
- `0 <= height[i] <= 10⁴`

---

## Company Frequency

| Company   | Frequency       | Last Reported |
|-----------|-----------------|---------------|
| Amazon    | ★★★★★ Very High | 2024          |
| Google    | ★★★★★ Very High | 2024          |
| Meta      | ★★★★★ Very High | 2024          |
| Microsoft | ★★★★☆ High      | 2024          |
| Bloomberg | ★★★★☆ High      | 2024          |
| Apple     | ★★★☆☆ Medium    | 2023          |
| Uber      | ★★★☆☆ Medium    | 2023          |
| Adobe     | ★★☆☆☆ Low       | 2022          |

> ⚠️ Frequency data is crowd-sourced from LeetCode Discuss, Glassdoor, and community interview reports. Treat as a signal, not a guarantee.

---

## DSA Concepts Used

- **Two Pointers** — the optimal solution starts with the widest container and greedily moves the shorter wall inward. → see [`/dsa/two_pointers.md`](/dsa/two_pointers.md)
- **Greedy** — at each step, we make the locally optimal choice (move the shorter wall) which is provably globally optimal.

---

## Approaches Overview

| # | Approach | Time | Space | When to use |
|---|----------|------|-------|-------------|
| 1 | Brute Force | O(n²) | O(1) | Never in interviews; TLE at n=10⁵ |
| 2 | Two Pointers ✅ | O(n) | O(1) | Always; this is the textbook answer |

---

## Approach 1 — Brute Force

### Intuition
Try every pair `(i, j)` and compute `min(height[i], height[j]) * (j-i)`. Track the maximum.

### Algorithm
1. For every `i` from 0 to n-2:
2. For every `j` from `i+1` to n-1:
   - `area = min(height[i], height[j]) * (j - i)`
   - Update `maxArea`.

### Complexity
- **Time:** O(n²).
- **Space:** O(1).

---

## Approach 2 — Two Pointers (Recommended ✅)

### Intuition
Start with the widest possible container (`l=0`, `r=n-1`). At each step, the area is `min(height[l], height[r]) * (r-l)`. The width can only decrease as we move inward, so we need height to increase to beat the current area. The only way to possibly get a taller bottle is to discard the shorter wall — moving the taller wall inward can only decrease (or maintain) the minimum height, which can never help.

**Formal proof:** When we advance `l` (because `height[l] <= height[r]`), we claim no pair `(l, j)` with `j < r` can beat `(l, r)`:
- For any `j < r`: `area(l,j) = min(height[l], height[j]) * (j-l) ≤ height[l] * (r-l)`.
- The current area `(l,r) = min(height[l], height[r]) * (r-l) ≥ height[l] * (r-l)` when `height[r] ≥ height[l]`.

So `(l,r)` is already the best pair involving `l`. We safely discard `l`.

### Algorithm
1. `l=0, r=n-1, maxArea=0`.
2. While `l < r`:
   - `area = min(height[l], height[r]) * (r-l)`.
   - Update `maxArea`.
   - If `height[l] <= height[r]`: `l++`, else `r--`.
3. Return `maxArea`.

### Complexity
- **Time:** O(n) — each pointer moves at most n steps total.
- **Space:** O(1).

### Code
```go
func twoPointers(height []int) int {
    l, r := 0, len(height)-1
    maxArea := 0
    for l < r {
        h := height[l]
        if height[r] < h { h = height[r] }
        if area := h * (r - l); area > maxArea { maxArea = area }
        if height[l] <= height[r] { l++ } else { r-- }
    }
    return maxArea
}
```

### Dry Run — `height = [1,8,6,2,5,4,8,3,7]`
```
l=0, r=8: area = min(1,7)*(8-0) = 1*8 = 8. h[0]=1 < h[8]=7 → l++
l=1, r=8: area = min(8,7)*(8-1) = 7*7 = 49. h[1]=8 > h[8]=7 → r--
l=1, r=7: area = min(8,3)*(7-1) = 3*6 = 18. h[7]=3 < h[1]=8 → r--
l=1, r=6: area = min(8,8)*(6-1) = 8*5 = 40. h[1]=8 ≥ h[6]=8 → l++
l=2, r=6: area = min(6,8)*(6-2) = 6*4 = 24. h[2]=6 < h[6]=8 → l++
l=3, r=6: area = min(2,8)*(6-3) = 2*3 = 6.  h[3]=2 < h[6]=8 → l++
l=4, r=6: area = min(5,8)*(6-4) = 5*2 = 10. h[4]=5 < h[6]=8 → l++
l=5, r=6: area = min(4,8)*(6-5) = 4*1 = 4.  h[5]=4 < h[6]=8 → l++
l=6 == r=6 → stop.
maxArea = 49 ✓
```

---

## Key Takeaways

- **Width decreases as we move inward** — the width can never recover. So to find a better area, we need a taller minimum height. Moving the shorter wall is the only action that can help.
- **The greedy choice is safe** — when we advance `l`, we provably discard all remaining pairs involving `l`. This is the key insight to justify the algorithm's correctness.
- **This is the canonical two-pointer problem** — the pattern "start from both ends, advance the smaller one" appears in many problems (e.g., Trapping Rain Water, 3Sum).
- **Difference from Trapping Rain Water (#42)** — this problem asks for the area of a single container (two walls); #42 asks for the total water trapped between all bars. Different DP/stack approaches apply there.

---

## Implementation (Go)

See [main.go](main.go).

### Verification
```
height=[1,8,6,2,5,4,8,3,7]  → 49 ✓
height=[1,1]                  → 1  ✓
height=[4,3,2,1,4]            → 16 ✓
height=[1,2,1]                → 2  ✓
```

---

## Related Problems

- LeetCode #42 — Trapping Rain Water (water trapped between all bars; harder)
- LeetCode #15 — 3Sum (two pointers on sorted array)
- LeetCode #16 — 3Sum Closest (two pointers variant)
