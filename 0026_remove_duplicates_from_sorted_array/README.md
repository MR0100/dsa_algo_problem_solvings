# 0026 — Remove Duplicates from Sorted Array

> LeetCode #26 · Difficulty: Easy
> **Categories:** Array, Two Pointers

---

## Problem Statement

Given an integer array `nums` sorted in **non-decreasing order**, remove the duplicates **in-place** such that each unique element appears only **once**. The **relative order** of the elements should be kept the **same**. Then return the number of unique elements in `nums`.

Consider the number of unique elements of `nums` to be `k`. To get accepted, you need to do the following things:

- Change the array `nums` such that the first `k` elements of `nums` contain the unique elements in the order they were present in `nums` initially.
- The remaining elements of `nums` are not important as well as the size of `nums`.
- Return `k`.

**Example 1**
```
Input:  nums = [1,1,2]
Output: 2, nums = [1,2,_]
Explanation: Your function should return k = 2, with the first two elements of nums being 1 and 2 respectively.
```

**Example 2**
```
Input:  nums = [0,0,1,1,1,2,2,3,3,4]
Output: 5, nums = [0,1,2,3,4,_,_,_,_,_]
```

**Constraints**
- `1 <= nums.length <= 3 * 10⁴`
- `-100 <= nums[i] <= 100`
- `nums` is sorted in **non-decreasing** order.

---

## Company Frequency

| Company   | Frequency       | Last Reported |
|-----------|-----------------|---------------|
| Amazon    | ★★★★☆ High      | 2024          |
| Google    | ★★★★☆ High      | 2024          |
| Microsoft | ★★★☆☆ Medium    | 2024          |
| Meta      | ★★★☆☆ Medium    | 2023          |
| Apple     | ★★☆☆☆ Low       | 2022          |
| Bloomberg | ★★☆☆☆ Low       | 2022          |

> ⚠️ Frequency data is crowd-sourced from LeetCode Discuss, Glassdoor, and community interview reports. Treat as a signal, not a guarantee.

---

## DSA Concepts Used

- **Two Pointers** — a write pointer `k` and a read pointer `i` allow in-place deduplication in a single pass.
- **Sorted Array Property** — duplicates are always adjacent, so we only need to compare adjacent elements.

---

## Approaches Overview

| # | Approach | Time | Space | When to use |
|---|----------|------|-------|-------------|
| 1 | Brute Force (extra slice) | O(n) | O(n) | Never; violates in-place constraint |
| 2 | Two Pointers / Write Pointer ✅ | O(n) | O(1) | The canonical O(1)-space answer |

---

## Approach 1 — Brute Force (Extra Slice)

### Intuition
Copy unique elements into a separate slice, then copy them back into `nums`. Easy to reason about but uses O(n) extra space.

### Algorithm
1. `unique = [nums[0]]`.
2. For each element: append to `unique` if it differs from `unique`'s last element.
3. Copy `unique` back into `nums[0..k-1]`. Return `k = len(unique)`.

### Complexity
- **Time:** O(n).
- **Space:** O(n) — the extra slice.

---

## Approach 2 — Two Pointers / Write Pointer (Recommended ✅)

### Intuition
Use two pointers on the same array:
- `k` (write pointer) — position where the next unique element should go.
- `i` (read pointer) — scans forward through the array.

When `nums[i]` differs from `nums[k-1]` (the last unique element written), copy it to `nums[k]` and increment `k`. Because the array is sorted, duplicates are always adjacent — this one comparison per step is sufficient.

### Algorithm
```
k = 1  // nums[0] is trivially unique
for i = 1 to n-1:
  if nums[i] != nums[k-1]:
    nums[k] = nums[i]
    k++
return k
```

### Complexity
- **Time:** O(n) — single pass.
- **Space:** O(1) — only two index variables.

### Code
```go
func twoPointers(nums []int) int {
    k := 1
    for i := 1; i < len(nums); i++ {
        if nums[i] != nums[k-1] {
            nums[k] = nums[i]
            k++
        }
    }
    return k
}
```

### Dry Run — `nums = [0,0,1,1,1,2,2,3,3,4]`
```
k=1
i=1: nums[1]=0 == nums[0]=0 → skip
i=2: nums[2]=1 != nums[0]=0 → nums[1]=1; k=2
i=3: nums[3]=1 == nums[1]=1 → skip
i=4: nums[4]=1 == nums[1]=1 → skip
i=5: nums[5]=2 != nums[1]=1 → nums[2]=2; k=3
i=6: nums[6]=2 == nums[2]=2 → skip
i=7: nums[7]=3 != nums[2]=2 → nums[3]=3; k=4
i=8: nums[8]=3 == nums[3]=3 → skip
i=9: nums[9]=4 != nums[3]=3 → nums[4]=4; k=5

Result: k=5, nums=[0,1,2,3,4,...] ✓
```

---

## Key Takeaways

- **Sorted = adjacent duplicates** — because duplicates are always next to each other in a sorted array, a single comparison `nums[i] != nums[k-1]` catches all duplicates without a hash set.
- **Write pointer pattern** — `k` tracks "what has been written so far." The read pointer `i` laps ahead freely. This is the same core pattern as LeetCode #27 (Remove Element) and #80 (Remove Duplicates II).
- **Compare to write pointer, not read pointer** — compare `nums[i]` with `nums[k-1]` (last written), not `nums[i-1]` (last read). On a run of duplicates the read pointer advances but k does not; after the run, `k-1 < i-1`.

---

## Related Problems

- LeetCode #27 — Remove Element (same write-pointer pattern, filter by value)
- LeetCode #80 — Remove Duplicates from Sorted Array II (allow up to 2 duplicates)
- LeetCode #83 — Remove Duplicates from Sorted List (same idea on a linked list)
