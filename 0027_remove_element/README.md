# 0027 — Remove Element

> LeetCode #27 · Difficulty: Easy
> **Categories:** Array, Two Pointers

---

## Problem Statement

Given an integer array `nums` and an integer `val`, remove all occurrences of `val` in `nums` **in-place**. The order of the elements may be changed. Then return the number of elements in `nums` which are not equal to `val`.

Consider the number of elements in `nums` which are not equal to `val` be `k`. To get accepted, you need to do the following things:

- Change the array `nums` such that the first `k` elements of `nums` contain the elements which are not equal to `val`.
- The remaining elements of `nums` are not important as well as the size of `nums`.
- Return `k`.

**Example 1**
```
Input:  nums = [3,2,2,3], val = 3
Output: 2, nums = [2,2,_,_]
```

**Example 2**
```
Input:  nums = [0,1,2,2,3,0,4,2], val = 2
Output: 5, nums = [0,1,4,0,3,_,_,_]
```

**Constraints**
- `0 <= nums.length <= 100`
- `0 <= nums[i] <= 50`
- `0 <= val <= 100`

---

## Company Frequency

| Company   | Frequency       | Last Reported |
|-----------|-----------------|---------------|
| Amazon    | ★★★☆☆ Medium    | 2024          |
| Google    | ★★★☆☆ Medium    | 2024          |
| Microsoft | ★★☆☆☆ Low       | 2023          |
| Meta      | ★★☆☆☆ Low       | 2023          |
| Bloomberg | ★★☆☆☆ Low       | 2022          |

> ⚠️ Frequency data is crowd-sourced from LeetCode Discuss, Glassdoor, and community interview reports. Treat as a signal, not a guarantee.

---

## DSA Concepts Used

- **Two Pointers** — both a "write pointer" approach (preserves order, O(n) writes) and a "swap from end" approach (minimal writes, O(k) where k = occurrences of val).

---

## Approaches Overview

| # | Approach | Time | Space | When to use |
|---|----------|------|-------|-------------|
| 1 | Shift Left | O(n²) | O(1) | Never; O(n²) is too slow |
| 2 | Write Pointer ✅ | O(n) | O(1) | Order must be preserved |
| 3 | Swap from End ✅ | O(n) | O(1) | val is rare; minimise write count |

---

## Approach 1 — Shift Left (Brute Force)

### Intuition
When we find `nums[i] == val`, shift all later elements one position left to fill the gap. Decrement the array length. Repeat until no more matches.

### Algorithm
```
n = len(nums); i = 0
while i < n:
  if nums[i] == val:
    shift nums[i+1..n-1] left by 1; n--
  else:
    i++
return n
```

### Complexity
- **Time:** O(n²) — each match triggers an O(n) shift.
- **Space:** O(1).

---

## Approach 2 — Write Pointer (Recommended ✅, order preserved)

### Intuition
Walk `i` across the array. Every element that should be kept (`nums[i] != val`) is written to `nums[k]`, then `k` increments. Elements equal to `val` are simply skipped — the write pointer doesn't advance.

### Algorithm
```
k = 0
for each e in nums:
  if e != val:
    nums[k] = e
    k++
return k
```

### Complexity
- **Time:** O(n).
- **Space:** O(1).

### Code
```go
func twoPointers(nums []int, val int) int {
    k := 0
    for _, e := range nums {
        if e != val {
            nums[k] = e
            k++
        }
    }
    return k
}
```

### Dry Run — `nums = [0,1,2,2,3,0,4,2]`, `val = 2`
```
k=0
e=0: 0≠2 → nums[0]=0; k=1
e=1: 1≠2 → nums[1]=1; k=2
e=2: 2=2 → skip
e=2: 2=2 → skip
e=3: 3≠2 → nums[2]=3; k=3
e=0: 0≠2 → nums[3]=0; k=4
e=4: 4≠2 → nums[4]=4; k=5
e=2: 2=2 → skip

Result: k=5, nums=[0,1,3,0,4,...] ✓
```

---

## Approach 3 — Swap from End (Optimal for rare val)

### Intuition
Use two boundaries: `left` scans from the start, `right` marks the end of the valid region. When `nums[left] == val`, overwrite it with `nums[right-1]` and shrink `right`. No shift needed; `right` decrements instead. When `nums[left] != val`, advance `left`.

This performs at most `k` writes (one per occurrence of val) vs Approach 2's n writes.

### Algorithm
```
left=0, right=len(nums)
while left < right:
  if nums[left] == val:
    nums[left] = nums[right-1]
    right--
  else:
    left++
return right
```

### Complexity
- **Time:** O(n).
- **Space:** O(1).

### Dry Run — `nums = [3,2,2,3]`, `val = 3`
```
left=0, right=4
nums[0]=3=val: nums[0]=nums[3]=3; right=3  → [3,2,2,3] (same value)
nums[0]=3=val: nums[0]=nums[2]=2; right=2  → [2,2,2,3]
nums[0]=2≠val: left=1
nums[1]=2≠val: left=2
left=2 == right=2 → stop

Return 2. nums[:2]=[2,2] ✓
```

---

## Key Takeaways

- **Write pointer vs swap from end** — use write pointer when element order matters; use swap from end when order doesn't matter and `val` is rare (fewer writes = better cache performance).
- **In-place does not mean zero writes** — "in-place" means O(1) extra space, not zero modifications. Both O(n) solutions here write at most n times to the same array.
- **The judge only checks nums[:k]** — elements beyond k are irrelevant; you don't need to zero them out.

---

## Related Problems

- LeetCode #26 — Remove Duplicates from Sorted Array (same write-pointer pattern)
- LeetCode #80 — Remove Duplicates from Sorted Array II (allow up to 2 duplicates)
- LeetCode #283 — Move Zeroes (write-pointer pattern, keep order, zeros at end)
