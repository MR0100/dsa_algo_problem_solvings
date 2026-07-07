# 0080 — Remove Duplicates from Sorted Array II

> LeetCode #80 · Difficulty: Medium
> **Categories:** Array, Two Pointers

---

## Problem Statement

Given an integer array `nums` sorted in **non-decreasing order**, remove some duplicates **in-place** such that each unique element appears **at most twice**. The **relative order** of the elements should be kept the same.

Since it is impossible to change the length of the array in some languages, you must instead have the result be placed in the **first part** of the array `nums`. More formally, if there are `k` elements after removing the duplicates, then the first `k` elements of `nums` should hold the final result. It does not matter what you leave beyond the first `k` elements.

Return `k` after placing the final result in the first `k` slots of `nums`.

**Do not** allocate extra space for another array. You must do this by **modifying the input array in-place** with O(1) extra memory.

**Example 1:**
```
Input: nums = [1,1,1,2,2,3]
Output: 5, nums = [1,1,2,2,3,_]
Explanation: Your function should return k = 5, with the first five elements of nums being 1, 1, 2, 2 and 3.
```

**Example 2:**
```
Input: nums = [0,0,1,1,1,1,2,3,3]
Output: 7, nums = [0,0,1,1,2,3,3,_,_]
Explanation: Your function should return k = 7, with the first seven elements of nums being 0, 0, 1, 1, 2, 3, and 3.
```

**Constraints:**
- `1 <= nums.length <= 3 * 10^4`
- `-10^4 <= nums[i] <= 10^4`
- `nums` is sorted in **non-decreasing order**.

---

## Company Frequency

| Company   | Frequency      | Last Reported |
|-----------|----------------|---------------|
| Amazon    | ★★★☆☆ Medium   | 2024          |
| Facebook  | ★★★☆☆ Medium   | 2023          |
| Microsoft | ★★☆☆☆ Low      | 2023          |

> ⚠️ Frequency data is crowd-sourced from LeetCode Discuss, Glassdoor, and community interview reports. Treat as a signal, not a guarantee.

---

## DSA Concepts Used

- **Two Pointers (Write Pointer)** — one pointer reads, one writes; write pointer tracks the "valid" prefix. See [`/dsa/two_pointers.md`](/dsa/two_pointers.md)
- **In-Place Modification** — modify the front of the array without allocating extra space.

---

## Approaches Overview

| # | Approach | Time | Space | When to use |
|---|----------|------|-------|-------------|
| 1 | Count-Based Write Pointer | O(n) | O(1) | Intuitive; easy to explain |
| 2 | Compare with k-2 | O(n) | O(1) | Elegant; generalises to "at most m duplicates" |

---

## Approach 1 — Count-Based Write Pointer

### Intuition
Walk through `nums` with a read pointer. Track how many times the current element has been seen consecutively (`count`). Only write it to position `k` if `count <= 2`.

### Algorithm
1. `k=0, count=0, prev=nil`.
2. For each `num`:
   - If `num == prev`: `count++`.
   - Else: `count=1; prev=num`.
   - If `count <= 2`: `nums[k]=num; k++`.
3. Return `k`.

### Complexity
- **Time:** O(n)
- **Space:** O(1)

### Code
```go
func removeDuplicatesCount(nums []int) int {
    k := 0
    count := 0
    var prev *int

    for _, num := range nums {
        if prev != nil && num == *prev {
            count++
        } else {
            count = 1
            v := num
            prev = &v
        }
        if count <= 2 {
            nums[k] = num
            k++
        }
    }
    return k
}
```

### Dry Run (nums=[1,1,1,2,2,3])

| num | prev | count | count<=2? | nums (first k) | k |
|-----|------|-------|-----------|----------------|---|
| 1 | nil | 1 | yes | [1] | 1 |
| 1 | 1 | 2 | yes | [1,1] | 2 |
| 1 | 1 | 3 | no | [1,1] | 2 |
| 2 | 1 | 1 | yes | [1,1,2] | 3 |
| 2 | 2 | 2 | yes | [1,1,2,2] | 4 |
| 3 | 2 | 1 | yes | [1,1,2,2,3] | 5 |

Return `k=5`, nums prefix = `[1,1,2,2,3]` ✓

---

## Approach 2 — Compare with k-2 (Elegant)

### Intuition
Use write pointer `k`. For each element `num`, write it if:
- `k < 2` (the first two elements always go in — there can't be 3 duplicates if only 2 exist yet), OR
- `num != nums[k-2]` (the element two positions back in the *output* is different, so adding `num` won't create 3+ consecutive duplicates).

If `nums[i] == nums[k-2]`, the last two written positions already hold this value — writing a third copy would violate the constraint.

**Why `k-2`?** Because `nums[k-1]` and `nums[k-2]` are the last two written values. If `num == nums[k-2]`, the sequence `nums[k-2], nums[k-1], num` would have 3 equal values.

### Algorithm
1. `k = 0`.
2. For each `num`: if `k < 2 || num != nums[k-2]`: `nums[k] = num; k++`.
3. Return `k`.

### Complexity
- **Time:** O(n)
- **Space:** O(1)

### Code
```go
func removeDuplicates(nums []int) int {
    k := 0
    for _, num := range nums {
        if k < 2 || num != nums[k-2] {
            nums[k] = num
            k++
        }
    }
    return k
}
```

### Dry Run (nums=[0,0,1,1,1,1,2,3,3])

| num | k | nums[k-2] | write? | k after | nums prefix |
|-----|---|-----------|--------|---------|-------------|
| 0 | 0 | — (k<2) | yes | 1 | [0] |
| 0 | 1 | — (k<2) | yes | 2 | [0,0] |
| 1 | 2 | 0 | 1≠0: yes | 3 | [0,0,1] |
| 1 | 3 | 0 | 1≠0: yes | 4 | [0,0,1,1] |
| 1 | 4 | 1 | 1==1: no | 4 | [0,0,1,1] |
| 1 | 4 | 1 | 1==1: no | 4 | [0,0,1,1] |
| 2 | 4 | 1 | 2≠1: yes | 5 | [0,0,1,1,2] |
| 3 | 5 | 1 | 3≠1: yes | 6 | [0,0,1,1,2,3] |
| 3 | 6 | 2 | 3≠2: yes | 7 | [0,0,1,1,2,3,3] |

Return `k=7` ✓

---

## Key Takeaways
- The `k-2` trick is the "canonical" solution because it generalises: to allow at most `m` duplicates, change `k-2` to `k-m` and `k < 2` to `k < m`.
- Works correctly for `k < 2` because there's no risk of creating a triple with fewer than 2 elements written.
- This is the same two-pointer "write head" pattern as LeetCode #26 (Remove Duplicates I), extended by one position.

---

## Related Problems
- LeetCode #26 — Remove Duplicates from Sorted Array (at most 1 occurrence — compare with `k-1`)
- LeetCode #27 — Remove Element (remove all occurrences of a value)
- LeetCode #283 — Move Zeroes (in-place rearrangement with write pointer)
