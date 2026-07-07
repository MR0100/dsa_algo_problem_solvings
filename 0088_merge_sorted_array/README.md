# 0088 — Merge Sorted Array

> LeetCode #88 · Difficulty: Easy
> **Categories:** Array, Two Pointers, Sorting

---

## Problem Statement

You are given two integer arrays `nums1` and `nums2`, sorted in **non-decreasing order**, and two integers `m` and `n`, representing the number of elements in `nums1` and `nums2` respectively.

**Merge** `nums1` and `nums2` into a single array sorted in **non-decreasing order**.

The final sorted array should not be returned by the function, but instead be **stored inside the array** `nums1`. To accommodate this, `nums1` has a length of `m + n`, where the first `m` elements denote the elements that should be merged, and the last `n` elements are set to 0 and should be ignored. `nums2` has a length of `n`.

**Example 1:**
```
Input: nums1 = [1,2,3,0,0,0], m = 3, nums2 = [2,5,6], n = 3
Output: [1,2,2,3,5,6]
```

**Example 2:**
```
Input: nums1 = [1], m = 1, nums2 = [], n = 0
Output: [1]
```

**Example 3:**
```
Input: nums1 = [0], m = 0, nums2 = [1], n = 1
Output: [1]
```

**Constraints:**
- `nums1.length == m + n`
- `nums2.length == n`
- `0 <= m, n <= 200`
- `1 <= m + n`
- `-10^9 <= nums1[i], nums2[j] <= 10^9`

**Follow-up:** Can you come up with an algorithm that runs in `O(m + n)` time?

---

## Company Frequency

| Company   | Frequency       | Last Reported |
|-----------|-----------------|---------------|
| Facebook  | ★★★★★ Very High | 2024          |
| Amazon    | ★★★★☆ High      | 2024          |
| Google    | ★★★☆☆ Medium    | 2024          |
| Microsoft | ★★★☆☆ Medium    | 2023          |

> ⚠️ Frequency data is crowd-sourced from LeetCode Discuss, Glassdoor, and community interview reports. Treat as a signal, not a guarantee.

---

## DSA Concepts Used

- **Two Pointers (merge from back)** — writing from the end avoids overwriting unread elements. See [`/dsa/two_pointers.md`](/dsa/two_pointers.md)

---

## Approaches Overview

| # | Approach | Time | Space | When to use |
|---|----------|------|-------|-------------|
| 1 | Copy + Sort | O((m+n)log(m+n)) | O(1) | Never in interview; conceptual baseline |
| 2 | Three Pointers from End | O(m+n) | O(1) | Always — optimal and expected |

---

## Approach 1 — Copy + Sort

### Intuition
Copy `nums2` into `nums1[m:]`, then sort the entire array. Doesn't exploit the sorted property.

### Algorithm
1. `copy(nums1[m:], nums2[:n])`.
2. Sort `nums1[0:m+n]`.

### Complexity
- **Time:** O((m+n) log(m+n))
- **Space:** O(1) — in-place sort.

### Code
```go
func mergeSimple(nums1 []int, m int, nums2 []int, n int) {
    copy(nums1[m:], nums2[:n])
    // insertion sort (or use sort.Ints(nums1))
    for i := 1; i < m+n; i++ {
        key := nums1[i]
        j := i - 1
        for j >= 0 && nums1[j] > key { nums1[j+1] = nums1[j]; j-- }
        nums1[j+1] = key
    }
}
```

### Dry Run (nums1=[1,2,3,0,0,0], nums2=[2,5,6])
After copy: `[1,2,3,2,5,6]`. After sort: `[1,2,2,3,5,6]` ✓

---

## Approach 2 — Three Pointers from End (Optimal)

### Intuition
If we write from the *front*, we'd overwrite `nums1` elements before reading them. Writing from the *back* solves this: the tail of `nums1` (positions `m` to `m+n-1`) is empty (zeros), so we can safely write there.

Compare the largest unmerged elements from both arrays (`p1 = m-1`, `p2 = n-1`), write the larger to position `p = m+n-1`, and decrement the appropriate pointer.

**Key insight:** if `p2 < 0` (nums2 exhausted), the remaining `nums1` elements are already in their correct positions.

### Algorithm
1. `p1=m-1, p2=n-1, p=m+n-1`.
2. While `p2 >= 0`:
   - If `p1 >= 0 && nums1[p1] > nums2[p2]`: `nums1[p] = nums1[p1]; p1--`.
   - Else: `nums1[p] = nums2[p2]; p2--`.
   - `p--`.

### Complexity
- **Time:** O(m+n) — single pass.
- **Space:** O(1)

### Code
```go
func merge(nums1 []int, m int, nums2 []int, n int) {
    p1 := m - 1
    p2 := n - 1
    p := m + n - 1
    for p2 >= 0 {
        if p1 >= 0 && nums1[p1] > nums2[p2] {
            nums1[p] = nums1[p1]; p1--
        } else {
            nums1[p] = nums2[p2]; p2--
        }
        p--
    }
}
```

### Dry Run (nums1=[1,2,3,0,0,0] m=3, nums2=[2,5,6] n=3)

| p1 | p2 | p | nums1[p1] | nums2[p2] | write | nums1 |
|----|----|----|-----------|-----------|-------|-------|
| 2 | 2 | 5 | 3 | 6 | 6 | [...,0,6] |
| 2 | 1 | 4 | 3 | 5 | 5 | [...,5,6] |
| 2 | 0 | 3 | 3 | 2 | 3 | [1,2,3,3,5,6] wait: write 3 at p=3 |
| 1 | 0 | 2 | 2 | 2 | 2(p2) | [1,2,2,3,5,6]? Let me retrace |

Corrected trace:

| p1 | p2 | p | compare | write | nums1 |
|----|----|----|---------|-------|-------|
| 2 | 2 | 5 | 3 vs 6 → 6 wins | nums1[5]=6 | [1,2,3,0,0,6] |
| 2 | 1 | 4 | 3 vs 5 → 5 wins | nums1[4]=5 | [1,2,3,0,5,6] |
| 2 | 0 | 3 | 3 vs 2 → 3 wins | nums1[3]=3 | [1,2,3,3,5,6] |
| 1 | 0 | 2 | 2 vs 2 → tie → p2 wins | nums1[2]=2 | [1,2,2,3,5,6] |
| 1 | -1 | 1 | p2<0, done | | [1,2,2,3,5,6] ✓ |

---

## Key Takeaways
- Writing from the back exploits the `m+n` buffer at the end of `nums1`.
- When `p2 < 0`, stop — remaining `nums1` elements are already in place.
- When elements are equal, prefer `p2` (write from `nums2`); this is correct since both values are equal.
- This O(m+n) in-place merge without extra space is a classic interview pattern.

---

## Related Problems
- LeetCode #21 — Merge Two Sorted Lists (same merge logic on linked lists)
- LeetCode #23 — Merge k Sorted Lists (generalization)
- LeetCode #912 — Sort an Array (merge sort uses this merge step)
