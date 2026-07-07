# 0035 — Search Insert Position

> LeetCode #35 · Difficulty: Easy
> **Categories:** Array, Binary Search

---

## Problem Statement

Given a sorted array of distinct integers and a target value, return the index if the target is found. If not, return the index where it would be if it were inserted in order.

You must write an algorithm with `O(log n)` runtime complexity.

**Example 1**
```
Input:  nums = [1,3,5,6], target = 5
Output: 2
```

**Example 2**
```
Input:  nums = [1,3,5,6], target = 2
Output: 1
```

**Example 3**
```
Input:  nums = [1,3,5,6], target = 7
Output: 4
```

**Constraints**
- `1 <= nums.length <= 10⁴`
- `-10⁴ <= nums[i] <= 10⁴`
- `nums` contains **distinct** values sorted in ascending order.
- `-10⁴ <= target <= 10⁴`

---

## Company Frequency

| Company   | Frequency       | Last Reported |
|-----------|-----------------|---------------|
| Amazon    | ★★★★☆ High      | 2024          |
| Google    | ★★★☆☆ Medium    | 2024          |
| Microsoft | ★★★☆☆ Medium    | 2024          |
| Meta      | ★★★☆☆ Medium    | 2023          |
| Bloomberg | ★★☆☆☆ Low       | 2022          |
| Apple     | ★★☆☆☆ Low       | 2022          |

> ⚠️ Frequency data is crowd-sourced from LeetCode Discuss, Glassdoor, and community interview reports. Treat as a signal, not a guarantee.

---

## DSA Concepts Used

- **Binary Search (Lower Bound)** — find the leftmost index where `nums[i] >= target`. This is exactly `std::lower_bound` in C++ / `sort.SearchInts` in Go.

---

## Approaches Overview

| # | Approach | Time | Space | When to use |
|---|----------|------|-------|-------------|
| 1 | Linear Scan | O(n) | O(1) | Doesn't satisfy O(log n) requirement |
| 2 | Binary Search ✅ | O(log n) | O(1) | Required; also the universal lower-bound template |

---

## Approach 1 — Linear Scan (Brute Force)

### Intuition
Scan forward; return the first index where `nums[i] >= target`, or `len(nums)` if not found.

### Complexity
- **Time:** O(n).
- **Space:** O(1).

### Code
```go
func linearScan(nums []int, target int) int {
    for i, v := range nums {
        if v >= target {
            return i
        }
    }
    return len(nums) // target is larger than all elements
}
```

### Dry Run — `nums = [1,3,5,6]`, `target = 2`

Scan left to right, return the first index whose value `≥ target`:

| `i` | `v` | `v >= 2`? | Action |
|-----|-----|-----------|--------|
| 0 | 1 | no | continue |
| 1 | 3 | yes | `return 1` |

Result: **1** (2 would be inserted between 1 and 3).

---

## Approach 2 — Binary Search / Lower Bound (Recommended ✅)

### Intuition
This problem is exactly the **lower bound** query: find the leftmost index where `nums[i] >= target`.

Binary search convergence:
- If `nums[mid] < target` → insertion point is strictly to the right → `lo = mid + 1`.
- If `nums[mid] >= target` → `mid` could be the answer, but there might be an earlier valid position → `hi = mid - 1`.

When the loop exits (`lo > hi`), `lo` points to the insertion position:
- If `target` is in the array, `lo` points to the first occurrence.
- If `target` is not in the array, `lo` is where it would be inserted.

### Algorithm
```
lo=0, hi=n-1
while lo <= hi:
  mid=(lo+hi)/2
  if nums[mid] < target: lo=mid+1
  else: hi=mid-1
return lo
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
        if nums[mid] < target { lo = mid + 1 } else { hi = mid - 1 }
    }
    return lo
}
```

### Dry Run — `nums = [1,3,5,6]`, `target = 2`
```
lo=0, hi=3
mid=1: nums[1]=3 >= 2 → hi=0

lo=0, hi=0
mid=0: nums[0]=1 < 2 → lo=1

lo=1 > hi=0 → stop.
return lo=1 ✓  (target 2 would be inserted between 1 and 3)
```

### Dry Run — `nums = [1,3,5,6]`, `target = 7`
```
lo=0, hi=3
mid=1: nums[1]=3 < 7 → lo=2
lo=2, hi=3
mid=2: nums[2]=5 < 7 → lo=3
lo=3, hi=3
mid=3: nums[3]=6 < 7 → lo=4
lo=4 > hi=3 → stop.
return lo=4 ✓  (insert at end)
```

---

## Key Takeaways

- **`lo` at loop exit is always the answer** — this is the universal property of the lower-bound binary search. It requires no special-casing for "not found".
- **`lo` converges to the insertion point** — every iteration either pushes `lo` right (nums[mid] too small) or pushes `hi` left (nums[mid] large enough). When they cross, `lo` is the smallest index where `nums[lo] >= target`.
- **Template to memorise** — `while lo<=hi: if nums[mid]<target: lo=mid+1 else: hi=mid-1; return lo`. This handles search, insert, lower bound all at once.
- **`sort.SearchInts(nums, target)` in Go** — the standard library implements exactly this; good to know for production code.

---

## Related Problems

- LeetCode #34 — Find First and Last Position (lower bound + upper bound)
- LeetCode #278 — First Bad Version (lower bound on a boolean predicate)
- LeetCode #367 — Valid Perfect Square (binary search on integers)
- LeetCode #744 — Find Smallest Letter Greater Than Target (upper bound)
