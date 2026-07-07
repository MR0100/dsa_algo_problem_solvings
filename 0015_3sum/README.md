# 0015 — 3Sum

> LeetCode #15 · Difficulty: Medium
> **Categories:** Array, Two Pointers, Sorting

---

## Problem Statement

Given an integer array `nums`, return all the triplets `[nums[i], nums[j], nums[k]]` such that `i != j`, `i != k`, and `j != k`, and `nums[i] + nums[j] + nums[k] == 0`.

Notice that the solution set must not contain duplicate triplets.

**Example 1**
```
Input:  nums = [-1,0,1,2,-1,-4]
Output: [[-1,-1,2],[-1,0,1]]
Explanation:
  nums[0] + nums[1] + nums[2] = (-1) + 0 + 1 = 0.
  nums[1] + nums[2] + nums[4] = 0 + 1 + (-1) = 0.
  nums[0] + nums[3] + nums[4] = (-1) + 2 + (-1) = 0.
  The distinct triplets are [-1,0,1] and [-1,-1,2].
```

**Example 2**
```
Input:  nums = [0,1,1]
Output: []
```

**Example 3**
```
Input:  nums = [0,0,0]
Output: [[0,0,0]]
```

**Constraints**
- `3 <= nums.length <= 3000`
- `-10⁵ <= nums[i] <= 10⁵`

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
| Adobe     | ★★★☆☆ Medium    | 2023          |
| Uber      | ★★★☆☆ Medium    | 2023          |
| LinkedIn  | ★★☆☆☆ Low       | 2022          |

> ⚠️ Frequency data is crowd-sourced from LeetCode Discuss, Glassdoor, and community interview reports. Treat as a signal, not a guarantee.

---

## DSA Concepts Used

- **Sorting** — sorting enables duplicate-skipping and the two-pointer technique. → see [`/dsa/two_pointers.md`](/dsa/two_pointers.md)
- **Two Pointers** — after fixing the first element at `i`, converging `l` and `r` pointers find all valid pairs in `nums[i+1:]` in O(n).
- **Hash Set** — Approach 3 uses a set to detect the required complement, at the cost of O(n) extra space per outer iteration.

---

## Approaches Overview

| # | Approach | Time | Space | When to use |
|---|----------|------|-------|-------------|
| 1 | Brute Force | O(n³) | O(n) | Never in practice; TLE at n=3000 |
| 2 | Sort + Two Pointers ✅ | O(n²) | O(1) | The standard answer; sort once, find pairs in O(n) each |
| 3 | Hash Set per Pair | O(n²) | O(n) | When sorting is forbidden; harder to deduplicate |

---

## Approach 1 — Brute Force

### Intuition
Try every combination of three indices (i, j, k). Deduplicate using a sorted triple as a map key.

### Complexity
- **Time:** O(n³).
- **Space:** O(n) — the dedup map.

---

## Approach 2 — Sort + Two Pointers (Recommended ✅)

### Intuition
Sort `nums`. For each index `i` (the leftmost element of the triplet), find pairs in `nums[i+1:]` that sum to `-nums[i]` using converging two pointers. Since the array is sorted, we can:
- Skip duplicate `i` values (same first element → same triplets).
- Skip duplicate `l` and `r` values after recording a match.
- Early-exit if `nums[i] > 0` (sorted, so all remaining elements are also positive; no zero-sum possible).

### Algorithm
1. `sort.Ints(nums)`.
2. For `i` from 0 to n-3:
   - If `nums[i] > 0`: break.
   - If `i > 0 && nums[i] == nums[i-1]`: skip (duplicate first element).
   - Set `l = i+1`, `r = n-1`, `target = -nums[i]`.
   - While `l < r`:
     - `sum = nums[l] + nums[r]`.
     - `sum == target`: record, skip duplicate l and r, advance both.
     - `sum < target`: `l++`.
     - `sum > target`: `r--`.

### Complexity
- **Time:** O(n log n) + O(n²) = O(n²).
- **Space:** O(1) extra (sort in-place; output is O(k) where k = result count).

### Code
```go
func twoPointers(nums []int) [][]int {
    sort.Ints(nums)
    n := len(nums)
    var result [][]int
    for i := 0; i < n-2; i++ {
        if nums[i] > 0 { break }
        if i > 0 && nums[i] == nums[i-1] { continue }
        l, r := i+1, n-1
        for l < r {
            sum := nums[l] + nums[r]
            if sum == -nums[i] {
                result = append(result, []int{nums[i], nums[l], nums[r]})
                for l < r && nums[l] == nums[l+1] { l++ }
                for l < r && nums[r] == nums[r-1] { r-- }
                l++; r--
            } else if sum < -nums[i] {
                l++
            } else {
                r--
            }
        }
    }
    return result
}
```

### Dry Run — `nums = [-1,0,1,2,-1,-4]`
```
After sort: [-4,-1,-1,0,1,2]

i=0, nums[0]=-4, target=4:
  l=1,r=5: -1+2=1 < 4 → l++
  l=2,r=5: -1+2=1 < 4 → l++
  l=3,r=5: 0+2=2 < 4  → l++
  l=4,r=5: 1+2=3 < 4  → l++
  l=5, l≥r → stop

i=1, nums[1]=-1, target=1:
  l=2,r=5: -1+2=1 == 1 → record [-1,-1,2]
    skip dup l: nums[2]=-1, nums[3]=0 → no dup
    skip dup r: nums[5]=2, nums[4]=1  → no dup
    l=3, r=4
  l=3,r=4: 0+1=1 == 1  → record [-1,0,1]
    l=4, r=3 → stop

i=2, nums[2]=-1, skip (nums[2]==nums[1])

i=3, nums[3]=0, target=0:
  l=4,r=5: 1+2=3 > 0 → r--
  l=4,r=4 → stop

Result: [[-1,-1,2],[-1,0,1]] ✓
```

---

## Approach 3 — Hash Set per Pair

### Intuition
Fix index `i`. For each `j > i`, the required third element is `complement = -nums[i]-nums[j]`. Check if `complement` was already seen in the inner pass (nums[i+1..j-1]). If yes, we found a valid triplet. Deduplicate via a sorted-triple map.

### Complexity
- **Time:** O(n²).
- **Space:** O(n) per outer iteration.

---

## Key Takeaways

- **Sort → two-pointer is the canonical O(n²) technique for k-sum** — it reduces 3Sum to repeated 2Sum on sorted arrays.
- **Three duplicate-skipping rules** — (1) skip duplicate `i`, (2) skip duplicate `l` after a match, (3) skip duplicate `r` after a match. Forgetting any one produces duplicate triplets.
- **Early exit when `nums[i] > 0`** — in a sorted array, if the smallest element in the remaining range is positive, no zero-sum triplet can exist.
- **Generalises to kSum** — for 4Sum, 5Sum: fix one element and recurse into (k-1)Sum. The two-pointer base case is 2Sum on a sorted array.

---

## Related Problems

- LeetCode #1 — Two Sum (2Sum with hash map)
- LeetCode #16 — 3Sum Closest (track closest instead of exact zero)
- LeetCode #18 — 4Sum (fix two elements, two-pointer for the inner pair)
- LeetCode #259 — 3Sum Smaller (count pairs with sum < target)
