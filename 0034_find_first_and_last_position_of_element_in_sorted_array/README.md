# 0034 — Find First and Last Position of Element in Sorted Array

> LeetCode #34 · Difficulty: Medium
> **Categories:** Array, Binary Search

---

## Problem Statement

Given an array of integers `nums` sorted in non-decreasing order, find the starting and ending position of a given `target` value.

If `target` is not found in the array, return `[-1, -1]`.

You must write an algorithm with `O(log n)` runtime complexity.

**Example 1**
```
Input:  nums = [5,7,7,8,8,10], target = 8
Output: [3,4]
```

**Example 2**
```
Input:  nums = [5,7,7,8,8,10], target = 6
Output: [-1,-1]
```

**Example 3**
```
Input:  nums = [], target = 0
Output: [-1,-1]
```

**Constraints**
- `0 <= nums.length <= 10⁵`
- `-10⁹ <= nums[i] <= 10⁹`
- `nums` is a non-decreasing integer array.
- `-10⁹ <= target <= 10⁹`

---

## Company Frequency

| Company   | Frequency       | Last Reported |
|-----------|-----------------|---------------|
| Amazon    | ★★★★★ Very High | 2024          |
| Google    | ★★★★☆ High      | 2024          |
| Microsoft | ★★★★☆ High      | 2024          |
| Meta      | ★★★☆☆ Medium    | 2023          |
| Bloomberg | ★★★☆☆ Medium    | 2023          |
| Apple     | ★★★☆☆ Medium    | 2023          |
| Adobe     | ★★☆☆☆ Low       | 2022          |

> ⚠️ Frequency data is crowd-sourced from LeetCode Discuss, Glassdoor, and community interview reports. Treat as a signal, not a guarantee.

---

## DSA Concepts Used

- **Binary Search** — two binary searches: one for the **lower bound** (first occurrence) and one for the **upper bound** (last occurrence). Each biases in opposite directions when `nums[mid] == target`.

---

## Approaches Overview

| # | Approach | Time | Space | When to use |
|---|----------|------|-------|-------------|
| 1 | Linear Scan | O(n) | O(1) | Trivial; doesn't satisfy O(log n) requirement |
| 2 | Two Binary Searches ✅ | O(log n) | O(1) | Required by problem; clean two-function implementation |

---

## Approach 1 — Linear Scan (Brute Force)

### Intuition
Walk the array; record the first and last indices where `nums[i] == target`.

### Complexity
- **Time:** O(n).
- **Space:** O(1).

---

## Approach 2 — Two Binary Searches (Recommended ✅)

### Intuition
Standard binary search finds *an* occurrence; we need the *first* and *last*. Two directional biases achieve this:

**Find first (lower bound):**
When `nums[mid] == target`, record `mid` as a candidate, then set `hi = mid - 1` to search further left.

**Find last (upper bound):**
When `nums[mid] == target`, record `mid` as a candidate, then set `lo = mid + 1` to search further right.

Both searches run independently in O(log n).

### Algorithm
```
findFirst(nums, target):
  lo=0, hi=n-1, result=-1
  while lo <= hi:
    mid = (lo+hi)/2
    if nums[mid] == target: result=mid; hi=mid-1    // bias left
    elif nums[mid] < target: lo=mid+1
    else: hi=mid-1
  return result

findLast(nums, target):
  lo=0, hi=n-1, result=-1
  while lo <= hi:
    mid = (lo+hi)/2
    if nums[mid] == target: result=mid; lo=mid+1    // bias right
    elif nums[mid] < target: lo=mid+1
    else: hi=mid-1
  return result
```

### Complexity
- **Time:** O(log n) — two independent binary searches.
- **Space:** O(1).

### Code
```go
func findFirst(nums []int, target int) int {
    lo, hi, result := 0, len(nums)-1, -1
    for lo <= hi {
        mid := (lo + hi) / 2
        if nums[mid] == target { result = mid; hi = mid - 1 } else
        if nums[mid] < target  { lo = mid + 1 } else { hi = mid - 1 }
    }
    return result
}
func findLast(nums []int, target int) int {
    lo, hi, result := 0, len(nums)-1, -1
    for lo <= hi {
        mid := (lo + hi) / 2
        if nums[mid] == target { result = mid; lo = mid + 1 } else
        if nums[mid] < target  { lo = mid + 1 } else { hi = mid - 1 }
    }
    return result
}
```

### Dry Run — `nums = [5,7,7,8,8,10]`, `target = 8`

**findFirst:**
```
lo=0, hi=5
mid=2: nums[2]=7<8 → lo=3
lo=3, hi=5
mid=4: nums[4]=8=8 → result=4; hi=3
lo=3, hi=3
mid=3: nums[3]=8=8 → result=3; hi=2
lo=3 > hi=2 → stop. first=3 ✓
```

**findLast:**
```
lo=0, hi=5
mid=2: nums[2]=7<8 → lo=3
lo=3, hi=5
mid=4: nums[4]=8=8 → result=4; lo=5
lo=5, hi=5
mid=5: nums[5]=10>8 → hi=4
lo=5 > hi=4 → stop. last=4 ✓
```

Result: `[3, 4]` ✓

---

## Key Takeaways

- **Lower/upper bound pattern** — memorise: `hi=mid-1` on match → lower bound; `lo=mid+1` on match → upper bound. This two-direction bias is the core of all "find range in sorted array" problems.
- **`result` variable instead of returning immediately** — returning immediately gives *an* occurrence; recording and continuing gives the *first* or *last*.
- **Two searches are better than one** — trying to find both bounds in a single search adds complexity without saving time; two clean O(log n) searches are clearer.
- **Empty array is handled by `hi = n-1 = -1`** — the loop never executes and `result = -1` is returned.

---

## Related Problems

- LeetCode #33 — Search in Rotated Sorted Array (binary search variant)
- LeetCode #35 — Search Insert Position (lower bound / first ≥ target)
- LeetCode #278 — First Bad Version (lower bound on a boolean predicate)
- LeetCode #744 — Find Smallest Letter Greater Than Target (upper bound variant)
