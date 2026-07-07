# 0075 — Sort Colors

> LeetCode #75 · Difficulty: Medium
> **Categories:** Array, Two Pointers, Sorting

---

## Problem Statement

Given an array `nums` with `n` objects colored red, white, or blue, sort them **in-place** so that objects of the same color are adjacent, with the colors in the order red, white, and blue.

We will use the integers `0`, `1`, and `2` to represent the color red, white, and blue, respectively.

You must solve this problem without using the library's sort function.

**Example 1**
```
Input:  nums = [2,0,2,1,1,0]
Output: [0,0,1,1,2,2]
```

**Example 2**
```
Input:  nums = [2,0,1]
Output: [0,1,2]
```

**Constraints**
- `n == nums.length`
- `1 <= n <= 300`
- `nums[i]` is either `0`, `1`, or `2`.

**Follow-up:** Could you come up with a one-pass algorithm using only constant extra space?

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

- **Dutch National Flag Algorithm** — three-way partition using `lo`, `mid`, `hi` pointers; single pass O(1) space.
- **Count Sort** — O(n) time O(1) space but two passes.

---

## Approaches Overview

| # | Approach | Time | Space | Passes | When to use |
|---|----------|------|-------|--------|-------------|
| 1 | Count Sort | O(n) | O(1) | 2 | Simple; limited to known value range |
| 2 | Dutch National Flag ✅ | O(n) | O(1) | 1 | Optimal one-pass; the classic answer |

---

## Approach 1 — Count Sort

### Intuition
Count the frequency of 0, 1, and 2. Overwrite the array: `count[0]` zeros, then `count[1]` ones, then `count[2]` twos.

### Complexity
- **Time:** O(n) — two passes.
- **Space:** O(1) — three-element count.

### Code
```go
func countSort(nums []int) {
	count := [3]int{}
	for _, v := range nums {
		count[v]++
	}
	i := 0
	for color := 0; color <= 2; color++ {
		for j := 0; j < count[color]; j++ {
			nums[i] = color
			i++
		}
	}
}
```

### Dry Run — `nums = [2,0,2,1,1,0]`

**Pass 1 — count each value:**

| v (scanned) | count[0] | count[1] | count[2] |
|-------------|----------|----------|----------|
| 2           | 0        | 0        | 1        |
| 0           | 1        | 0        | 1        |
| 2           | 1        | 0        | 2        |
| 1           | 1        | 1        | 2        |
| 1           | 1        | 2        | 2        |
| 0           | 2        | 2        | 2        |

Final counts: `count = [2, 2, 2]`.

**Pass 2 — overwrite** (`i` = write index):

| color | writes | nums after            |
|-------|--------|-----------------------|
| 0     | 2 → i=0,1 | `[0,0,_,_,_,_]`    |
| 1     | 2 → i=2,3 | `[0,0,1,1,_,_]`    |
| 2     | 2 → i=4,5 | `[0,0,1,1,2,2]`    |

Result: `[0,0,1,1,2,2]` ✓

---

## Approach 2 — Dutch National Flag (Recommended ✅)

### Intuition
Named after the Dutch flag's three horizontal stripes. Use three pointers:
- `lo`: next position for 0 (everything left of `lo` is 0).
- `mid`: current element being processed.
- `hi`: next position for 2 (everything right of `hi` is 2).

The **invariant** at all times:
- `nums[0..lo-1]` = all 0s.
- `nums[lo..mid-1]` = all 1s.
- `nums[mid..hi]` = unknown (being processed).
- `nums[hi+1..n-1]` = all 2s.

At each step:
- `nums[mid] == 0`: swap with `lo`, advance both `lo` and `mid`. The 0 joins the left section; the 1 that came from `lo` is safe to advance past.
- `nums[mid] == 1`: advance `mid` only (1 is already in the correct middle section).
- `nums[mid] == 2`: swap with `hi`, decrement `hi`. DON'T advance `mid` — the swapped element from `hi` is unknown and must be re-examined.

### Algorithm
```
lo=0, mid=0, hi=n-1
while mid <= hi:
  switch nums[mid]:
    0: swap(lo,mid); lo++; mid++
    1: mid++
    2: swap(mid,hi); hi--
```

### Complexity
- **Time:** O(n) — each element examined at most twice (once by `mid`, once after a swap).
- **Space:** O(1).

### Code
```go
func dutchFlag(nums []int) {
    lo, mid, hi := 0, 0, len(nums)-1
    for mid <= hi {
        switch nums[mid] {
        case 0: nums[lo], nums[mid] = nums[mid], nums[lo]; lo++; mid++
        case 1: mid++
        case 2: nums[mid], nums[hi] = nums[hi], nums[mid]; hi--
        }
    }
}
```

### Dry Run — `nums = [2,0,2,1,1,0]`
```
lo=0, mid=0, hi=5:
  nums[0]=2: swap(0,5). nums=[0,0,2,1,1,2]. hi=4.
  nums[0]=0: swap(0,0). nums=[0,...]. lo=1, mid=1.
  nums[1]=0: swap(1,1). lo=2, mid=2.
  nums[2]=2: swap(2,4). nums=[0,0,1,1,2,2]. hi=3.
  nums[2]=1: mid=3.
  nums[3]=1: mid=4. mid>hi=3 → stop.

Result: [0,0,1,1,2,2] ✓
```

---

## Key Takeaways

- **Don't advance `mid` after swapping with `hi`** — the element that came from `hi` is unknown (could be 0, 1, or 2) and must be re-checked. This is the most common implementation mistake.
- **After swapping with `lo`, we CAN advance `mid`** — because `nums[lo..mid-1]` is all 1s (invariant), so the element that came from `lo` must be a 1, which can be skipped.
- **Generalises to k colors** — for k distinct values, use a more complex variant (but in practice, DNF only appears as the 3-color problem).
- **This is the classic partitioning primitive** — QuickSort's 3-way partition is Dutch National Flag applied to pivot elements.

---

## Related Problems

- LeetCode #283 — Move Zeroes (partition 0s and non-zeros; 2-pointer variant)
- LeetCode #905 — Sort Array by Parity (partition even/odd; same 2-pointer idea)
- LeetCode #912 — Sort an Array (sort without library; merge/quick sort)
