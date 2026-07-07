# 0303 — Range Sum Query - Immutable

> LeetCode #303 · Difficulty: Easy
> **Categories:** Array, Design, Prefix Sum

---

## Problem Statement

Given an integer array `nums`, handle multiple queries of the following type:

1. Calculate the **sum** of the elements of `nums` between indices `left` and `right` **inclusive** where `left <= right`.

Implement the `NumArray` class:

- `NumArray(int[] nums)` Initializes the object with the integer array `nums`.
- `int sumRange(int left, int right)` Returns the sum of the elements of `nums` between indices `left` and `right` **inclusive** (i.e. `nums[left] + nums[left + 1] + ... + nums[right]`).

**Example 1:**

```
Input
["NumArray", "sumRange", "sumRange", "sumRange"]
[[[-2, 0, 3, -5, 2, -1]], [0, 2], [2, 5], [0, 5]]
Output
[null, 1, -1, -3]

Explanation
NumArray numArray = new NumArray([-2, 0, 3, -5, 2, -1]);
numArray.sumRange(0, 2); // return (-2) + 0 + 3 = 1
numArray.sumRange(2, 5); // return 3 + (-5) + 2 + (-1) = -1
numArray.sumRange(0, 5); // return (-2) + 0 + 3 + (-5) + 2 + (-1) = -3
```

**Constraints:**

- `1 <= nums.length <= 10^4`
- `-10^5 <= nums[i] <= 10^5`
- `0 <= left <= right < nums.length`
- At most `10^4` calls will be made to `sumRange`.

---

## Company Frequency

| Company   | Frequency        | Last Reported |
|-----------|------------------|---------------|
| Amazon    | ★★★★☆ High       | 2024          |
| Facebook  | ★★★☆☆ Medium     | 2023          |
| Google    | ★★★☆☆ Medium     | 2023          |
| Microsoft | ★★☆☆☆ Low        | 2022          |

> ⚠️ Frequency data is crowd-sourced from LeetCode Discuss, Glassdoor, and community interview reports. Treat as a signal, not a guarantee.

---

## DSA Concepts Used

- **Prefix Sum** — precompute cumulative sums so any range answers in O(1) → see [`/dsa/prefix_sum.md`](/dsa/prefix_sum.md)
- **Design (immutable structure)** — preprocess once in the constructor, answer many queries cheaply → see [`/dsa/design_data_structures.md`](/dsa/design_data_structures.md)

---

## Approaches Overview

| # | Approach | Constructor | Query | Space | When to use |
|---|----------|-------------|-------|-------|-------------|
| 1 | Brute force (re-sum) | O(1) | O(n) | O(n) | Few queries, or as a correctness baseline |
| 2 | Prefix sum (Optimal) | O(n) | O(1) | O(n) | Many queries on an immutable array — the intended solution |

---

## Approach 1 — Brute Force (Re-sum Each Query)

### Intuition
The most literal reading: store the numbers, and when asked for a range, loop and add. Correct, but every query costs O(n); with up to 10⁴ queries this is wasteful because the array never changes.

### Algorithm
1. Store the array reference (it is immutable, so no copy needed).
2. For `sumRange(left, right)`, iterate `i` from `left` to `right` accumulating a sum.
3. Return the sum.

### Complexity
- **Time:** constructor O(1); each `SumRange` is O(right − left + 1) = O(n) worst case.
- **Space:** O(n) to hold the array.

### Code
```go
type NumArrayBrute struct {
	nums []int // the original, immutable numbers
}

func NewNumArrayBrute(nums []int) NumArrayBrute {
	return NumArrayBrute{nums: nums}
}

func (a NumArrayBrute) SumRange(left, right int) int {
	sum := 0
	for i := left; i <= right; i++ { // walk the requested window
		sum += a.nums[i]
	}
	return sum
}
```

### Dry Run
Array `[-2, 0, 3, -5, 2, -1]`, query `sumRange(0, 2)`.

| i | nums[i] | sum |
|---|---|---|
| 0 | -2 | -2 |
| 1 | 0  | -2 |
| 2 | 3  | 1  |

Returns **1**.

---

## Approach 2 — Prefix Sum (Optimal)

### Intuition
Define `prefix[i]` = sum of the first `i` elements (`prefix[0] = 0`). Then the sum of `nums[left..right]` telescopes to `prefix[right+1] − prefix[left]`: the subtracted term removes everything strictly before `left`, leaving exactly the window. A single O(n) preprocessing pass makes every subsequent query O(1).

### Algorithm
1. Allocate `prefix` of length `n+1`, with `prefix[0] = 0`.
2. Fill `prefix[i+1] = prefix[i] + nums[i]`.
3. Answer `sumRange(left, right)` as `prefix[right+1] − prefix[left]`.

### Complexity
- **Time:** constructor O(n); each `SumRange` O(1).
- **Space:** O(n) for the prefix array.

### Code
```go
type NumArray struct {
	prefix []int // prefix[i] = sum of the first i elements
}

func NewNumArray(nums []int) NumArray {
	prefix := make([]int, len(nums)+1) // one extra slot; prefix[0] stays 0
	for i, v := range nums {
		prefix[i+1] = prefix[i] + v // running cumulative sum
	}
	return NumArray{prefix: prefix}
}

func (a NumArray) SumRange(left, right int) int {
	// Everything up to right, minus everything before left, = the window.
	return a.prefix[right+1] - a.prefix[left]
}
```

### Dry Run
Array `[-2, 0, 3, -5, 2, -1]`.

Build prefix:

| i | nums[i] | prefix[i+1] |
|---|---|---|
| — | —  | prefix[0] = 0 |
| 0 | -2 | -2 |
| 1 | 0  | -2 |
| 2 | 3  | 1  |
| 3 | -5 | -4 |
| 4 | 2  | -2 |
| 5 | -1 | -3 |

Queries:
- `sumRange(0,2)` = `prefix[3] − prefix[0]` = `1 − 0` = **1**
- `sumRange(2,5)` = `prefix[6] − prefix[2]` = `-3 − (-2)` = **-1**
- `sumRange(0,5)` = `prefix[6] − prefix[0]` = `-3 − 0` = **-3**

---

## Key Takeaways

- **Immutable + many queries ⇒ preprocess once.** The prefix-sum table converts O(n) queries into O(1) lookups.
- The `n+1`-sized prefix with a leading zero removes all off-by-one pain: `sum(l..r) = prefix[r+1] − prefix[l]` works even when `l = 0`.
- This is the foundational pattern behind range-sum problems; the mutable variant (#307) upgrades it to a Fenwick/segment tree.

---

## Related Problems

- LeetCode #304 — Range Sum Query 2D - Immutable (2D prefix sums)
- LeetCode #307 — Range Sum Query - Mutable (Fenwick / segment tree)
- LeetCode #560 — Subarray Sum Equals K (prefix sums + hash map)
- LeetCode #1480 — Running Sum of 1d Array (the prefix array itself)
