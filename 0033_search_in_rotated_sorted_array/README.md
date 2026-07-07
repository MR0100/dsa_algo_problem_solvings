# 0033 — Search in Rotated Sorted Array

> LeetCode #33 · Difficulty: Medium
> **Categories:** Array, Binary Search

---

## Problem Statement

There is an integer array `nums` sorted in ascending order (with **distinct** values).

Prior to being passed to your function, `nums` is possibly rotated at an unknown pivot index `k` (`1 <= k < nums.length`) such that the resulting array is `[nums[k], nums[k+1], ..., nums[n-1], nums[0], nums[1], ..., nums[k-1]]`.

Given the array `nums` after the possible rotation and an integer `target`, return the index of `target` if it is in `nums`, or `-1` if it is not in `nums`.

You must write an algorithm with `O(log n)` runtime complexity.

**Example 1**
```
Input:  nums = [4,5,6,7,0,1,2], target = 0
Output: 4
```

**Example 2**
```
Input:  nums = [4,5,6,7,0,1,2], target = 3
Output: -1
```

**Example 3**
```
Input:  nums = [1], target = 0
Output: -1
```

**Constraints**
- `1 <= nums.length <= 5000`
- `-10⁴ <= nums[i] <= 10⁴`
- All values of `nums` are **unique**.
- `nums` is an ascending array that is possibly rotated.
- `-10⁴ <= target <= 10⁴`

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
| Adobe     | ★★★☆☆ Medium    | 2023          |

> ⚠️ Frequency data is crowd-sourced from LeetCode Discuss, Glassdoor, and community interview reports. Treat as a signal, not a guarantee.

---

## DSA Concepts Used

- **Binary Search** — modified to work on a partially sorted (rotated) array by determining at each step which half is sorted.

---

## Approaches Overview

| # | Approach | Time | Space | When to use |
|---|----------|------|-------|-------------|
| 1 | Linear Scan | O(n) | O(1) | Quick sanity check; not the intended solution |
| 2 | Binary Search ✅ | O(log n) | O(1) | Required by the problem; standard interview answer |

---

## Approach 1 — Linear Scan (Brute Force)

### Intuition
Ignore the rotation; scan for the target.

### Complexity
- **Time:** O(n).
- **Space:** O(1).

### Code
```go
func bruteForce(nums []int, target int) int {
    for i, v := range nums {
        if v == target {
            return i
        }
    }
    return -1
}
```

### Dry Run — `nums = [4,5,6,7,0,1,2]`, `target = 0`
Linear scan visits each index until `nums[i] == target`.

| i | nums[i] | nums[i] == 0? | action |
|---|---------|---------------|--------|
| 0 | 4       | no            | continue |
| 1 | 5       | no            | continue |
| 2 | 6       | no            | continue |
| 3 | 7       | no            | continue |
| 4 | 0       | **yes**       | return 4 ✓ |

---

## Approach 2 — Modified Binary Search (Recommended ✅)

### Intuition
In a rotated sorted array, at least one of the two halves `[lo..mid]` or `[mid..hi]` is always **fully sorted**. We can tell which by comparing `nums[lo]` with `nums[mid]`:

- If `nums[lo] <= nums[mid]`: **left half is sorted**.
  - Is target in `[nums[lo], nums[mid])`? Search left.
  - Else? Search right.
- Else: **right half is sorted**.
  - Is target in `(nums[mid], nums[hi]]`? Search right.
  - Else? Search left.

This allows us to eliminate half the search space at each step, preserving O(log n).

### Algorithm
```
lo=0, hi=n-1
while lo <= hi:
  mid = (lo+hi)/2
  if nums[mid] == target: return mid
  if nums[lo] <= nums[mid]:  // left is sorted
    if nums[lo] <= target < nums[mid]: hi=mid-1
    else: lo=mid+1
  else:                      // right is sorted
    if nums[mid] < target <= nums[hi]: lo=mid+1
    else: hi=mid-1
return -1
```

### Complexity
- **Time:** O(log n).
- **Space:** O(1).

### Code
```go
func binarySearch(nums []int, target int) int {
    lo, hi := 0, len(nums)-1
    for lo <= hi {
        mid := (lo + hi) / 2
        if nums[mid] == target { return mid }
        if nums[lo] <= nums[mid] {
            if nums[lo] <= target && target < nums[mid] { hi = mid - 1 } else { lo = mid + 1 }
        } else {
            if nums[mid] < target && target <= nums[hi] { lo = mid + 1 } else { hi = mid - 1 }
        }
    }
    return -1
}
```

### Dry Run — `nums = [4,5,6,7,0,1,2]`, `target = 0`
```
lo=0, hi=6
mid=3: nums[3]=7 ≠ 0.
  nums[0]=4 <= nums[3]=7 → left sorted.
  target=0: 4 <= 0? No → search right. lo=4.

lo=4, hi=6
mid=5: nums[5]=1 ≠ 0.
  nums[4]=0 <= nums[5]=1 → left sorted.
  target=0: 0 <= 0 < 1? Yes → search left. hi=4.

lo=4, hi=4
mid=4: nums[4]=0 == 0 → return 4 ✓
```

---

## Key Takeaways

- **Always one sorted half** — this is the invariant that makes the modified binary search work. With distinct elements, `nums[lo] <= nums[mid]` unambiguously identifies the sorted half.
- **Strict vs inclusive inequalities matter** — `target < nums[mid]` (strict) because `nums[mid]` itself was already checked. `target <= nums[hi]` (inclusive) because `nums[hi]` hasn't been checked yet.
- **Variant with duplicates** — LeetCode #81 adds duplicates; when `nums[lo] == nums[mid]`, we can't determine which half is sorted → must increment `lo` (degrades to O(n) worst case).

---

## Related Problems

- LeetCode #81 — Search in Rotated Sorted Array II (allows duplicates)
- LeetCode #153 — Find Minimum in Rotated Sorted Array (find the pivot)
- LeetCode #154 — Find Minimum in Rotated Sorted Array II (with duplicates)
