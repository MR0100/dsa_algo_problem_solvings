# 0031 — Next Permutation

> LeetCode #31 · Difficulty: Medium
> **Categories:** Array, Two Pointers

---

## Problem Statement

A **permutation** of an array of integers is an arrangement of its members into a sequence or linear order.

The **next permutation** of an array of integers is the next lexicographically greater permutation of its integer. More formally, if all the permutations of the array are sorted in one container according to their lexicographical order, then the **next permutation** of that arrangement is the arrangement that follows it in that sorted order. If such an arrangement is not possible, the array must be rearranged as the lowest possible order (i.e., sorted in ascending order).

Modify the array **in-place** with only **constant extra memory**.

**Example 1**
```
Input:  nums = [1,2,3]
Output: [1,3,2]
```

**Example 2**
```
Input:  nums = [3,2,1]
Output: [1,2,3]
```

**Example 3**
```
Input:  nums = [1,1,5]
Output: [1,5,1]
```

**Constraints**
- `1 <= nums.length <= 100`
- `0 <= nums[i] <= 100`

---

## Company Frequency

| Company   | Frequency       | Last Reported |
|-----------|-----------------|---------------|
| Amazon    | ★★★★★ Very High | 2024          |
| Google    | ★★★★★ Very High | 2024          |
| Microsoft | ★★★★☆ High      | 2024          |
| Meta      | ★★★★☆ High      | 2024          |
| Bloomberg | ★★★★☆ High      | 2023          |
| Apple     | ★★★☆☆ Medium    | 2023          |
| Adobe     | ★★★☆☆ Medium    | 2023          |

> ⚠️ Frequency data is crowd-sourced from LeetCode Discuss, Glassdoor, and community interview reports. Treat as a signal, not a guarantee.

---

## DSA Concepts Used

- **Array In-Place Manipulation** — the solution modifies the array using only index variables (two index passes + a reverse).
- **Two Pointers** — the reverse step in pass 2 uses two converging pointers.
- **Greedy** — at each decision point we make the locally smallest change (swap the pivot with the smallest element to its right that is still larger).

---

## Approaches Overview

| # | Approach | Time | Space | When to use |
|---|----------|------|-------|-------------|
| 1 | Generate all permutations | O(n! · n) | O(n!) | Never; astronomically slow |
| 2 | Two-Pass In-Place ✅ | O(n) | O(1) | The only correct interview answer |

---

## Approach 1 — Brute Force (Generate All Permutations)

### Intuition
Generate every permutation in sorted order; find the current one; return the next. Impractical for n > 8.

### Complexity
- **Time:** O(n! · n).
- **Space:** O(n!).

---

## Approach 2 — Two-Pass In-Place (Recommended ✅)

### Intuition
The next permutation differs from the current one only in the **shortest possible suffix**. That means:

1. Find the rightmost **pivot** — the rightmost position `i` where `nums[i] < nums[i+1]`. Everything to the right of `i` is already in descending order (no "next" exists for that suffix alone).
2. Find the rightmost element `j` to the right of `i` that is strictly greater than `nums[i]` (the smallest such element, since the suffix is descending).
3. Swap `nums[i]` and `nums[j]`. Now `nums[i]` is slightly larger, but the suffix (right of `i`) is still descending.
4. Reverse the suffix `nums[i+1:]` to make it the smallest possible ascending order.

**Edge case:** If no pivot exists (entire array is non-increasing), it's already the last permutation → reverse the whole array to get the first.

### Algorithm
```
Step 1: i = n-2
        while i >= 0 and nums[i] >= nums[i+1]: i--

Step 2: if i >= 0:
          j = n-1
          while nums[j] <= nums[i]: j--
          swap(nums[i], nums[j])

Step 3: reverse nums[i+1 .. n-1]
```

### Complexity
- **Time:** O(n) — two left-to-right scans + one O(n) reverse.
- **Space:** O(1) — only index variables.

### Code
```go
func optimal(nums []int) {
    n := len(nums)
    i := n - 2
    for i >= 0 && nums[i] >= nums[i+1] { i-- }
    if i >= 0 {
        j := n - 1
        for nums[j] <= nums[i] { j-- }
        nums[i], nums[j] = nums[j], nums[i]
    }
    left, right := i+1, n-1
    for left < right {
        nums[left], nums[right] = nums[right], nums[left]
        left++; right--
    }
}
```

### Dry Run — `nums = [1,3,2]`
```
Step 1: Find pivot
  i=1: nums[1]=3 >= nums[2]=2 → i--
  i=0: nums[0]=1 <  nums[1]=3 → stop. Pivot i=0.

Step 2: Find swap partner
  j=2: nums[2]=2 > nums[0]=1 → stop. j=2.
  Swap: [1,3,2] → [2,3,1]

Step 3: Reverse suffix [i+1=1 .. n-1=2]:
  [2,3,1] → [2,1,3]

Result: [2,1,3]  ✓ (next after [1,3,2] in lex order)
```

---

## Key Takeaways

- **Suffix is descending when we stop** — the pivot-finding loop terminates when `nums[i] < nums[i+1]`, meaning `nums[i+1..n-1]` is entirely non-increasing. This is why reversing it gives the smallest possible suffix.
- **Swap with the *rightmost* element > pivot** — because the suffix is descending, the rightmost element greater than the pivot is also the smallest such element. This makes the new suffix as small as possible.
- **No pivot → reverse all** — `i` ends at -1, so `i+1 = 0` and reversing `nums[0..n-1]` gives the sorted (first) permutation.
- **Works with duplicates** — the `>=` in the pivot loop (`nums[i] >= nums[i+1]`) and `<=` in the j loop (`nums[j] <= nums[i]`) correctly handle equal elements.

---

## Related Problems

- LeetCode #46 — Permutations (generate all permutations)
- LeetCode #60 — Permutation Sequence (find the k-th permutation)
- LeetCode #556 — Next Greater Element III (same logic on digits of a number)
