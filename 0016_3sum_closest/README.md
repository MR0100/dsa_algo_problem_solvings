# 0016 — 3Sum Closest

> LeetCode #16 · Difficulty: Medium
> **Categories:** Array, Two Pointers, Sorting

---

## Problem Statement

Given an integer array `nums` of length `n` and an integer `target`, find three integers in `nums` such that the sum is closest to `target`.

Return the *sum* of the three integers.

You may assume that each input would have exactly one solution.

**Example 1**
```
Input:  nums = [-1,2,1,-4], target = 1
Output: 2
Explanation: The sum that is closest to the target is 2. (-1 + 2 + 1 = 2).
```

**Example 2**
```
Input:  nums = [0,0,0], target = 1
Output: 0
```

**Constraints**
- `3 <= nums.length <= 500`
- `-1000 <= nums[i] <= 1000`
- `-10⁴ <= target <= 10⁴`

---

## Company Frequency

| Company   | Frequency       | Last Reported |
|-----------|-----------------|---------------|
| Amazon    | ★★★★★ Very High | 2024          |
| Google    | ★★★★☆ High      | 2024          |
| Microsoft | ★★★★☆ High      | 2024          |
| Bloomberg | ★★★☆☆ Medium    | 2023          |
| Meta      | ★★★☆☆ Medium    | 2023          |
| Apple     | ★★☆☆☆ Low       | 2022          |
| Adobe     | ★★☆☆☆ Low       | 2022          |

> ⚠️ Frequency data is crowd-sourced from LeetCode Discuss, Glassdoor, and community interview reports. Treat as a signal, not a guarantee.

---

## DSA Concepts Used

- **Sorting** — enables the two-pointer technique and early-exit optimisations.
- **Two Pointers** — same converging-pointer pattern as 3Sum, adapted for "closest" instead of "exact". → see [`/dsa/two_pointers.md`](/dsa/two_pointers.md)

---

## Approaches Overview

| # | Approach | Time | Space | When to use |
|---|----------|------|-------|-------------|
| 1 | Brute Force | O(n³) | O(1) | Never in practice; n can be 500 |
| 2 | Sort + Two Pointers ✅ | O(n²) | O(1) | The standard answer |

---

## Approach 1 — Brute Force

### Intuition
Try all O(n³) triplets. Track the one with the smallest `|sum - target|`.

### Complexity
- **Time:** O(n³).
- **Space:** O(1).

---

## Approach 2 — Sort + Two Pointers (Recommended ✅)

### Intuition
Identical setup to 3Sum but instead of checking for exact zero, we track `closest` — the sum with minimum absolute distance to `target`. After sorting, fix `i` and use converging `l`, `r`:
- If `sum < target`: `l++` (need a larger sum).
- If `sum > target`: `r--` (need a smaller sum).
- If `sum == target`: return immediately (can't do better).
- At each step: update `closest` if `|sum - target| < |closest - target|`.

Duplicate skipping for `i` is an optional optimisation (safe to include since same `i` value leads to same search space).

### Algorithm
1. `sort.Ints(nums)`. `closest = nums[0]+nums[1]+nums[2]`.
2. For `i` from 0 to n-3:
   - Skip duplicate `i`.
   - `l = i+1`, `r = n-1`.
   - While `l < r`: compute `sum`, update `closest`, steer pointers.
3. Return `closest`.

### Complexity
- **Time:** O(n log n) + O(n²) = O(n²).
- **Space:** O(1) extra.

### Code
```go
func twoPointers(nums []int, target int) int {
    sort.Ints(nums)
    n := len(nums)
    closest := math.MaxInt32
    for i := 0; i < n-2; i++ {
        if i > 0 && nums[i] == nums[i-1] { continue }
        l, r := i+1, n-1
        for l < r {
            sum := nums[i] + nums[l] + nums[r]
            if abs(sum-target) < abs(closest-target) { closest = sum }
            if sum == target { return sum }
            if sum < target { l++ } else { r-- }
        }
    }
    return closest
}
```

### Dry Run — `nums = [-1,2,1,-4]`, `target = 1`
```
After sort: [-4,-1,1,2]
closest = -4+-1+1 = -4 (initial, will be updated)

i=0, nums[0]=-4:
  l=1,r=3: sum=-4+-1+2=-3. |(-3)-1|=4 < |(-4)-1|=5 → closest=-3. -3<1 → l++
  l=2,r=3: sum=-4+1+2=-1.  |(-1)-1|=2 < 4 → closest=-1. -1<1 → l++
  l=3,r=3 → stop

i=1, nums[1]=-1:
  l=2,r=3: sum=-1+1+2=2. |2-1|=1 < |(-1)-1|=2 → closest=2. 2>1 → r--
  l=2,r=2 → stop

i=2 → i < n-2=2? No → stop.

Return closest=2 ✓
```

---

## Key Takeaways

- **Same template as 3Sum** — sort + fix one element + two-pointer inner. The only change is tracking `closest` instead of collecting exact-zero triplets.
- **Steer toward target** — `sum < target → l++` (increase), `sum > target → r--` (decrease). The sorted array guarantees this moves sum in the right direction.
- **Early exit on exact match** — if `sum == target`, no triplet can be closer, so return immediately.
- **Initialise `closest` carefully** — initialise to any valid triplet sum (e.g., `nums[0]+nums[1]+nums[2]`) rather than `math.MaxInt32` to avoid edge cases with the `abs` comparison.

---

## Related Problems

- LeetCode #15 — 3Sum (exact zero sum; return all triplets)
- LeetCode #18 — 4Sum (k=4 variant)
- LeetCode #259 — 3Sum Smaller (count triplets with sum < target)
