# 0004 — Median of Two Sorted Arrays

> LeetCode #4 · Difficulty: Hard
> **Categories:** Array, Binary Search, Divide and Conquer

---

## Problem Statement

Given two sorted arrays `nums1` and `nums2` of size `m` and `n` respectively, return the **median** of the two sorted arrays.

The overall run time complexity should be `O(log(m+n))`.

**Example 1**
```
Input:  nums1 = [1,3], nums2 = [2]
Output: 2.00000
Explanation: merged array = [1,2,3] and median is 2.
```

**Example 2**
```
Input:  nums1 = [1,2], nums2 = [3,4]
Output: 2.50000
Explanation: merged array = [1,2,3,4] and median is (2+3)/2 = 2.5.
```

**Constraints**
- `nums1.length == m`
- `nums2.length == n`
- `0 <= m <= 1000`
- `0 <= n <= 1000`
- `1 <= m + n <= 2000`
- `-10⁶ <= nums1[i], nums2[i] <= 10⁶`

---

## Company Frequency

| Company   | Frequency       | Last Reported |
|-----------|-----------------|---------------|
| Google    | ★★★★★ Very High | 2024          |
| Amazon    | ★★★★☆ High      | 2024          |
| Microsoft | ★★★★☆ High      | 2024          |
| Meta      | ★★★★☆ High      | 2024          |
| Apple     | ★★★☆☆ Medium    | 2023          |
| Bloomberg | ★★★☆☆ Medium    | 2023          |
| Adobe     | ★★☆☆☆ Low       | 2022          |
| Uber      | ★★☆☆☆ Low       | 2022          |

> ⚠️ Frequency data is crowd-sourced from LeetCode Discuss, Glassdoor, and community
> interview reports. Treat as a signal, not a guarantee.

---

## DSA Concepts Used

- **Binary Search** — we binary-search the partition index of the shorter array; each step halves the search space, giving O(log(min(m,n))). → see [`/dsa/binary_search.md`](/dsa/binary_search.md)
- **Divide and Conquer** — finding the median is equivalent to finding the correct split of both arrays into a combined left half and right half; this split can be computed without merging. → see [`/dsa/divide_and_conquer.md`](/dsa/divide_and_conquer.md)
- **Array** — the problem requires indexed access into two sorted arrays; the binary search relies on reading boundary values by index.

---

## Approaches Overview

| # | Approach | Time | Space | When to use |
|---|----------|------|-------|-------------|
| 1 | Merge & Find | O(m+n) | O(m+n) | Small inputs; easy to implement |
| 2 | Concat & Sort | O((m+n)log(m+n)) | O(m+n) | Worst of all; only as a baseline |
| 3 | Two-Pointer Walk | O(m+n) | O(1) | When O(m+n) time is acceptable but O(m+n) space is not |
| 4 | Binary Search on Partition ✅ | O(log(min(m,n))) | O(1) | Required approach; meets the O(log) constraint |

---

## Approach 1 — Merge and Find Median

### Intuition
Build the fully merged sorted array, then read the median from the middle position. The "sorted" property of both inputs makes the merge O(m+n) instead of O((m+n) log(m+n)).

### Algorithm
1. Use two pointers (standard merge-sort merge step) to build `merged[]` in sorted order.
2. `total = m + n`.
3. If `total` is odd → return `merged[total/2]`.
4. If `total` is even → return `(merged[total/2-1] + merged[total/2]) / 2.0`.

### Complexity
- **Time:** O(m+n) — one linear merge pass.
- **Space:** O(m+n) — the merged array.

### Code
```go
func mergeAndFind(nums1 []int, nums2 []int) float64 {
    merged := make([]int, 0, len(nums1)+len(nums2))
    i, j := 0, 0
    for i < len(nums1) && j < len(nums2) {
        if nums1[i] <= nums2[j] { merged = append(merged, nums1[i]); i++ } else { merged = append(merged, nums2[j]); j++ }
    }
    merged = append(merged, nums1[i:]...)
    merged = append(merged, nums2[j:]...)
    total := len(merged); mid := total / 2
    if total%2 == 1 { return float64(merged[mid]) }
    return float64(merged[mid-1]+merged[mid]) / 2.0
}
```

### Dry Run — Example 2: `nums1=[1,2]`, `nums2=[3,4]`
```
Merge step:
  i=0,j=0: 1<3 → merged=[1], i=1
  i=1,j=0: 2<3 → merged=[1,2], i=2
  i exhausted → append [3,4]
  merged = [1,2,3,4]
total=4 (even), mid=2
return (merged[1]+merged[2])/2 = (2+3)/2 = 2.5 ✓
```

---

## Approach 2 — Concatenate and Sort

### Intuition
The crudest approach: ignore the sorted property of both inputs, concatenate them into one array, sort it, then read the median.

### Algorithm
1. `combined = nums1 + nums2`.
2. `sort.Ints(combined)`.
3. Read median from middle.

### Complexity
- **Time:** O((m+n) log(m+n)) — dominated by the sort.
- **Space:** O(m+n) — the combined array.

### Why included
This is strictly worse than Approach 1 — it wastes the sorted property. Included to show the progression from "completely naive" to "uses sorted property" to "optimal".

---

## Approach 3 — Two-Pointer Walk to Median Position

### Intuition
The same as Approach 1 but without storing the merged array. We walk two pointers forward, always advancing the one with the smaller current value, and count until we reach the median position(s). We only need to remember the last two values seen.

### Algorithm
1. `targetHigh = total / 2` (the median position in 0-indexed merged sequence).
2. Walk pointers: at each step advance whichever pointer has the smaller current value.
3. After `targetHigh + 1` steps, `prev` holds position `targetHigh-1` and `cur` holds `targetHigh`.
4. Odd total: return `cur`. Even total: return `(prev + cur) / 2.0`.

### Complexity
- **Time:** O(m+n) — walk up to `total/2 + 1` steps.
- **Space:** O(1) — no merged array; only two stored values.

### Dry Run — Example 1: `nums1=[1,3]`, `nums2=[2]`
```
total=3, targetHigh=1
step 0: min(1,2)=1 → cur=1, i=1
step 1: min(3,2)=2 → cur=2, j=1 (exhausted)
total odd → return cur = 2.0 ✓
```

---

## Approach 4 — Binary Search on Partition (Optimal)

### Intuition
The key insight: **finding the median is equivalent to finding the correct partition of both arrays into a left half and a right half.**

We want partition indices `i` (in `nums1`) and `j` (in `nums2`) such that:

```
Left half:  nums1[0..i-1]  ∪  nums2[0..j-1]   ← has ⌈(m+n)/2⌉ elements
Right half: nums1[i..m-1]  ∪  nums2[j..n-1]
```

The partition is **correct** when:
```
nums1[i-1] ≤ nums2[j]   AND   nums2[j-1] ≤ nums1[i]
```
(i.e., every element on the left is ≤ every element on the right)

The median then is:
- **Odd total:** `max(nums1[i-1], nums2[j-1])` — the largest element in the left half.
- **Even total:** `(max(left) + min(right)) / 2.0`.

Since `j = half - i` (determined by `i`), we only need to binary-search `i` over `[0, m]`. Always binary-search on the shorter array to get O(log(min(m,n))).

**Steering:**
- If `nums1[i-1] > nums2[j]` → `i` is too large → `hi = i - 1`.
- If `nums2[j-1] > nums1[i]` → `i` is too small → `lo = i + 1`.

**Sentinel values:** use `-∞` when `i=0` or `j=0` (nothing on the left) and `+∞` when `i=m` or `j=n` (nothing on the right), so the partition condition is trivially true for edge cuts.

### Algorithm
1. Swap so `nums1` is the shorter (search space = `[0, m]`).
2. `half = (m + n + 1) / 2`.
3. Binary search `lo=0, hi=m`:
   - `i = (lo+hi)/2`, `j = half - i`.
   - Compute the four boundary values with ±∞ sentinels.
   - If `nums1LeftMax ≤ nums2RightMin` AND `nums2LeftMax ≤ nums1RightMin` → found; return median.
   - Else if `nums1LeftMax > nums2RightMin` → `hi = i-1`.
   - Else → `lo = i+1`.

### Complexity
- **Time:** O(log(min(m,n))) — binary search on the shorter array.
- **Space:** O(1) — only index variables and four boundary values.

### Code
```go
func binarySearchPartition(nums1 []int, nums2 []int) float64 {
    if len(nums1) > len(nums2) { nums1, nums2 = nums2, nums1 }
    m, n := len(nums1), len(nums2)
    half := (m + n + 1) / 2
    lo, hi := 0, m
    for lo <= hi {
        i := (lo + hi) / 2
        j := half - i
        n1L := math.MinInt64; if i > 0 { n1L = nums1[i-1] }
        n1R := math.MaxInt64; if i < m { n1R = nums1[i] }
        n2L := math.MinInt64; if j > 0 { n2L = nums2[j-1] }
        n2R := math.MaxInt64; if j < n { n2R = nums2[j] }
        if n1L <= n2R && n2L <= n1R {
            leftMax := n1L; if n2L > leftMax { leftMax = n2L }
            if (m+n)%2 == 1 { return float64(leftMax) }
            rightMin := n1R; if n2R < rightMin { rightMin = n2R }
            return float64(leftMax+rightMin) / 2.0
        } else if n1L > n2R {
            hi = i - 1
        } else {
            lo = i + 1
        }
    }
    return 0.0
}
```

### Dry Run — Example 2: `nums1=[1,2]`, `nums2=[3,4]`
```
m=2, n=2, half=2, lo=0, hi=2

Iteration 1: i=1, j=1
  n1L=nums1[0]=1, n1R=nums1[1]=2
  n2L=nums2[0]=3, n2R=nums2[1]=4
  Check: 1 ≤ 4 ✅  BUT  3 ≤ 2 ❌ → n2L > n1R → lo = 2

Iteration 2: i=2, j=0
  n1L=nums1[1]=2, n1R=+∞ (i=m)
  n2L=-∞ (j=0),  n2R=nums2[0]=3
  Check: 2 ≤ 3 ✅  AND  -∞ ≤ +∞ ✅ → partition correct!
  leftMax = max(2, -∞) = 2
  total=4 (even) → rightMin = min(+∞, 3) = 3
  return (2+3)/2 = 2.5 ✓
```

---

## Key Takeaways

- **Partition = median** — the fundamental insight of this problem. You don't need to physically merge; you just need to find the right cut. Once you find it, the median is just `max(left)` or `(max(left)+min(right))/2`.
- **Binary search on a derived quantity** — we don't binary-search on values, we binary-search on a partition index. This is a common advanced binary search pattern.
- **Always search on the shorter array** — ensures `j = half - i` is always valid (j stays within [0,n]) and keeps time O(log(min(m,n))).
- **Sentinel values** — using `-∞` and `+∞` for edge cases (partition at the start/end of an array) eliminates all four out-of-bounds guards, making the code much cleaner.
- **`(m+n+1)/2` vs `(m+n)/2`** — adding `1` before dividing biases the left half to have the extra element when the total is odd. This ensures the median is `leftMax` (not `rightMin`) for odd totals.

---

## Related Problems

- LeetCode #215 — Kth Largest Element in an Array (similar "find kth element" theme)
- LeetCode #295 — Find Median from Data Stream (dynamic median with two heaps)
- LeetCode #378 — Kth Smallest Element in a Sorted Matrix (binary search on value range)
