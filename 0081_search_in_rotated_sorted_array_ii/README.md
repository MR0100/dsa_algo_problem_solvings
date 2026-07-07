# 0081 — Search in Rotated Sorted Array II

> LeetCode #81 · Difficulty: Medium
> **Categories:** Array, Binary Search

---

## Problem Statement

There is an integer array `nums` sorted in non-decreasing order (not necessarily with distinct values). Before being passed to your function, `nums` is possibly rotated at an unknown pivot index `k` such that the resulting array is `[nums[k], nums[k+1], ..., nums[n-1], nums[0], nums[1], ..., nums[k-1]]`.

Given the array `nums` after the possible rotation and an integer `target`, return `true` if `target` is in `nums`, or `false` if it is not.

You must decrease the overall operation steps as much as possible.

**Example 1:**
```
Input: nums = [2,5,6,0,0,1,2], target = 0
Output: true
```

**Example 2:**
```
Input: nums = [2,5,6,0,0,1,2], target = 3
Output: false
```

**Constraints:**
- `1 <= nums.length <= 5000`
- `-10^4 <= nums[i] <= 10^4`
- `nums` is an ascending array that is possibly rotated.
- `-10^4 <= target <= 10^4`

**Follow-up:** This problem is similar to Search in Rotated Sorted Array (#33), but `nums` may contain duplicates. Would this affect the runtime complexity? How and why?

---

## Company Frequency

| Company   | Frequency      | Last Reported |
|-----------|----------------|---------------|
| Amazon    | ★★★☆☆ Medium   | 2024          |
| LinkedIn  | ★★★☆☆ Medium   | 2023          |
| Microsoft | ★★☆☆☆ Low      | 2023          |

> ⚠️ Frequency data is crowd-sourced from LeetCode Discuss, Glassdoor, and community interview reports. Treat as a signal, not a guarantee.

---

## DSA Concepts Used

- **Binary Search** — with special handling for the ambiguous case when `nums[lo] == nums[mid]`. See [`/dsa/binary_search.md`](/dsa/binary_search.md)
- **Rotated Array** — one half is always sorted (unless duplicates prevent determination).

---

## Approaches Overview

| # | Approach | Time | Space | When to use |
|---|----------|------|-------|-------------|
| 1 | Linear Scan | O(n) | O(1) | Simple fallback; always correct |
| 2 | Binary Search | O(log n) avg, O(n) worst | O(1) | Preferred; fast except all-same arrays |

---

## Approach 1 — Linear Scan

### Intuition
When duplicates are present, the worst-case complexity of binary search degrades to O(n) anyway. A simple linear scan is always O(n) and correct.

### Algorithm
Scan every element. Return `true` if any equals `target`.

### Complexity
- **Time:** O(n)
- **Space:** O(1)

### Code
```go
func linearScan(nums []int, target int) bool {
    for _, v := range nums {
        if v == target {
            return true
        }
    }
    return false
}
```

### Dry Run (nums=[2,5,6,0,0,1,2], target=0)
Scan: 2 → 5 → 6 → 0 → match! Return true.

---

## Approach 2 — Binary Search (with duplicate skip)

### Intuition
The key difference from #33 is: when `nums[lo] == nums[mid]`, we cannot determine which half is sorted. Example: `[1,1,1,0,1]` — is the 0 on the left or right of mid?

Fix: when `nums[lo] == nums[mid]`, increment `lo` by 1 (skip one duplicate). This allows us to eventually reach a state where we can determine the sorted half. Worst case (all same) degrades to O(n).

When `nums[lo] != nums[mid]`, exactly one half is sorted (standard #33 logic applies).

### Algorithm
1. `lo=0, hi=n-1`.
2. While `lo <= hi`:
   - `mid = lo + (hi-lo)/2`.
   - If `nums[mid] == target`: return `true`.
   - If `nums[lo] == nums[mid]`: `lo++` (skip; can't determine sorted half).
   - Else if left half sorted (`nums[lo] <= nums[mid]`):
     - If `target` in `[nums[lo], nums[mid])`: `hi = mid-1`.
     - Else: `lo = mid+1`.
   - Else (right half sorted):
     - If `target` in `(nums[mid], nums[hi]]`: `lo = mid+1`.
     - Else: `hi = mid-1`.
3. Return `false`.

### Complexity
- **Time:** O(log n) average; O(n) worst case (e.g., `[1,1,1,1,1]`).
- **Space:** O(1)

### Code
```go
func binarySearch(nums []int, target int) bool {
    lo, hi := 0, len(nums)-1
    for lo <= hi {
        mid := lo + (hi-lo)/2
        if nums[mid] == target {
            return true
        }
        if nums[lo] == nums[mid] {
            lo++
            continue
        }
        if nums[lo] <= nums[mid] {
            if nums[lo] <= target && target < nums[mid] {
                hi = mid - 1
            } else {
                lo = mid + 1
            }
        } else {
            if nums[mid] < target && target <= nums[hi] {
                lo = mid + 1
            } else {
                hi = mid - 1
            }
        }
    }
    return false
}
```

### Dry Run (nums=[2,5,6,0,0,1,2], target=0)

| lo | hi | mid | nums[mid] | action |
|----|----|-----|-----------|--------|
| 0 | 6 | 3 | 0 | nums[mid]==target → return true |

(Found immediately at mid=3.)

### Dry Run (nums=[1,0,1,1,1], target=0)

| lo | hi | mid | nums[mid] | nums[lo] | action |
|----|----|-----|-----------|----------|--------|
| 0 | 4 | 2 | 1 | 1 | nums[lo]==nums[mid]: lo++ |
| 1 | 4 | 2 | 1 | 0 | Left NOT sorted (nums[lo]=0 > nums[mid]=1)... wait, 0<=1 so left IS sorted |
| | | | | | Left sorted: target=0 in [nums[1]=0, nums[2]=1)? 0 in [0,1) yes → hi=1 |
| 1 | 1 | 1 | 0 | | nums[mid]==target → true |

---

## Key Takeaways
- Duplicates destroy the ability to determine which half is sorted when `nums[lo] == nums[mid]`. The fix is conservative: skip one element.
- This is why the problem says "would this affect runtime complexity?" — yes, worst case becomes O(n) instead of O(log n).
- In practice for interviews, this O(n) edge case is acceptable and the binary search is the expected answer.

---

## Related Problems
- LeetCode #33 — Search in Rotated Sorted Array (no duplicates; true O(log n))
- LeetCode #153 — Find Minimum in Rotated Sorted Array
- LeetCode #154 — Find Minimum in Rotated Sorted Array II (with duplicates)
