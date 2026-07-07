# 0018 — 4Sum

> LeetCode #18 · Difficulty: Medium
> **Categories:** Array, Two Pointers, Sorting

---

## Problem Statement

Given an array `nums` of `n` integers, return an array of all the **unique** quadruplets `[nums[a], nums[b], nums[c], nums[d]]` such that:
- `0 <= a, b, c, d < n`
- `a`, `b`, `c`, and `d` are **distinct**.
- `nums[a] + nums[b] + nums[c] + nums[d] == target`

You may return the answer in **any order**.

**Example 1**
```
Input:  nums = [1,0,-1,0,-2,2], target = 0
Output: [[-2,-1,1,2],[-2,0,0,2],[-1,0,0,1]]
```

**Example 2**
```
Input:  nums = [2,2,2,2,2], target = 8
Output: [[2,2,2,2]]
```

**Constraints**
- `1 <= nums.length <= 200`
- `-10⁹ <= nums[i] <= 10⁹`
- `-10⁹ <= target <= 10⁹`

---

## Company Frequency

| Company   | Frequency       | Last Reported |
|-----------|-----------------|---------------|
| Amazon    | ★★★★☆ High      | 2024          |
| Google    | ★★★★☆ High      | 2024          |
| Microsoft | ★★★☆☆ Medium    | 2024          |
| Bloomberg | ★★★☆☆ Medium    | 2023          |
| Meta      | ★★★☆☆ Medium    | 2023          |
| Apple     | ★★☆☆☆ Low       | 2022          |

> ⚠️ Frequency data is crowd-sourced from LeetCode Discuss, Glassdoor, and community interview reports. Treat as a signal, not a guarantee.

---

## DSA Concepts Used

- **Sorting** — enables duplicate skipping and the two-pointer inner loop.
- **Two Pointers** — the innermost two-pointer pair reduces O(n²) pairs to O(n), giving O(n³) overall. → see [`/dsa/two_pointers.md`](/dsa/two_pointers.md)

---

## Approaches Overview

| # | Approach | Time | Space | When to use |
|---|----------|------|-------|-------------|
| 1 | Brute Force | O(n⁴) | O(n) | Never in practice |
| 2 | Sort + Two Pointers ✅ | O(n³) | O(1) | The standard answer; direct extension of 3Sum |

---

## Approach 1 — Brute Force

### Intuition
Try all O(n⁴) quadruplets. Deduplicate using a sorted 4-tuple as a map key.

### Complexity
- **Time:** O(n⁴).
- **Space:** O(n) — the dedup map.

---

## Approach 2 — Sort + Two Pointers (Recommended ✅)

### Intuition
4Sum = 3Sum with one extra outer loop. Sort the array. Fix two elements `nums[i]` and `nums[j]`, then use converging two pointers to find pairs in `nums[j+1:]` summing to `target - nums[i] - nums[j]`.

Duplicate handling:
- Skip `i` if `nums[i] == nums[i-1]` (same first element).
- Skip `j` if `nums[j] == nums[j-1]` and `j > i+1` (same second element within the same `i`).
- Skip `l` and `r` after recording a match (same as 3Sum).

**Overflow guard:** values can be up to ±10⁹, so four values summed can reach ±4×10⁹, which overflows `int32`. Cast to `int64` before comparing.

### Algorithm
1. Sort `nums`.
2. For `i` from 0 to n-4:
   - Skip dup `i`.
   - For `j` from `i+1` to n-3:
     - Skip dup `j` (guard: `j > i+1`).
     - Two-pointer `l=j+1`, `r=n-1` for `int64(sum) == int64(target)`.
     - Record, skip dups, advance on exact match; steer pointers otherwise.

### Complexity
- **Time:** O(n³) — O(n log n) + O(n²) outer × O(n) inner.
- **Space:** O(1) extra.

### Code
```go
func twoPointers(nums []int, target int) [][]int {
    sort.Ints(nums)
    n := len(nums)
    var result [][]int
    for i := 0; i < n-3; i++ {
        if i > 0 && nums[i] == nums[i-1] { continue }
        for j := i + 1; j < n-2; j++ {
            if j > i+1 && nums[j] == nums[j-1] { continue }
            l, r := j+1, n-1
            for l < r {
                sum := int64(nums[i]) + int64(nums[j]) + int64(nums[l]) + int64(nums[r])
                if sum == int64(target) {
                    result = append(result, []int{nums[i], nums[j], nums[l], nums[r]})
                    for l < r && nums[l] == nums[l+1] { l++ }
                    for l < r && nums[r] == nums[r-1] { r-- }
                    l++; r--
                } else if sum < int64(target) {
                    l++
                } else {
                    r--
                }
            }
        }
    }
    return result
}
```

### Dry Run — `nums = [1,0,-1,0,-2,2]`, `target = 0`
```
After sort: [-2,-1,0,0,1,2]

i=0 (nums[0]=-2):
  j=1 (nums[1]=-1), need sum=3:
    l=2,r=5: 0+2=2 < 3 → l++
    l=3,r=5: 0+2=2 < 3 → l++
    l=4,r=5: 1+2=3 == 3 → record [-2,-1,1,2]. l=5,r=4 → stop.
  j=2 (nums[2]=0), need sum=2:
    l=3,r=5: 0+2=2 == 2 → record [-2,0,0,2]. skip dup l(nums[3]=0=nums[4]? no), skip dup r. l=4,r=4 → stop.
  j=3 (nums[3]=0): skip (nums[3]==nums[2])
  j=4 (nums[4]=1), need sum=1:
    l=5,r=5 → stop immediately.

i=1 (nums[1]=-1):
  j=2 (nums[2]=0), need sum=1:
    l=3,r=5: 0+2=2 > 1 → r--
    l=3,r=4: 0+1=1 == 1 → record [-1,0,0,1]. l=4,r=3 → stop.
  j=3 (nums[3]=0): skip
  ...

Result: [[-2,-1,1,2],[-2,0,0,2],[-1,0,0,1]] ✓
```

---

## Key Takeaways

- **4Sum = 3Sum + one loop** — the kSum pattern generalises: add one more outer loop and recurse/reduce to (k-1)Sum. Base case is 2Sum with two pointers on a sorted array.
- **`int64` overflow guard** — values can be ±10⁹; four of them can exceed `int32` range (±2.1×10⁹). Always use `int64` for the sum comparison.
- **Duplicate-skipping guard for `j`** — `j > i+1` is critical. Without it, we'd skip valid quadruplets where `nums[j] == nums[i]` (the second element equals the first).
- **Generalise to kSum** — for k ≥ 3: sort once, fix (k-2) elements with nested loops, two-pointer for the innermost pair. Time: O(n^(k-1)).

---

## Related Problems

- LeetCode #15 — 3Sum (same two-pointer pattern, one fewer outer loop)
- LeetCode #16 — 3Sum Closest (closest sum variant)
- LeetCode #1 — Two Sum (2Sum with hash map)
- LeetCode #259 — 3Sum Smaller (count triplets < target)
