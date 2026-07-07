# 0268 — Missing Number

> LeetCode #268 · Difficulty: Easy
> **Categories:** Array, Hash Table, Math, Binary Search, Bit Manipulation, Sorting

---

## Problem Statement

Given an array `nums` containing `n` distinct numbers in the range `[0, n]`,
return the only number in the range that is missing from the array.

**Example 1:**
```
Input: nums = [3,0,1]
Output: 2
```
Explanation: n = 3 since there are 3 numbers, so all numbers are in the range
[0,3]. 2 is the missing number in the range since it does not appear in nums.

**Example 2:**
```
Input: nums = [0,1]
Output: 2
```
Explanation: n = 2 since there are 2 numbers, so all numbers are in the range
[0,2]. 2 is the missing number in the range since it does not appear in nums.

**Example 3:**
```
Input: nums = [9,6,4,2,3,5,7,0,1]
Output: 8
```
Explanation: n = 9 since there are 9 numbers, so all numbers are in the range
[0,9]. 8 is the missing number in the range since it does not appear in nums.

**Constraints:**
- `n == nums.length`
- `1 <= n <= 10^4`
- `0 <= nums[i] <= n`
- All the numbers of `nums` are unique.

**Follow up:** Could you implement a solution using only O(1) extra space
complexity and O(n) runtime complexity?

---

## Company Frequency

| Company    | Frequency        | Last Reported |
|------------|------------------|---------------|
| Amazon     | ★★★★☆ High       | 2024          |
| Microsoft  | ★★★★☆ High       | 2024          |
| Google     | ★★★☆☆ Medium     | 2023          |
| Bloomberg  | ★★★☆☆ Medium     | 2023          |
| Apple      | ★★☆☆☆ Low        | 2022          |

> ⚠️ Frequency data is crowd-sourced from LeetCode Discuss, Glassdoor, and
> community interview reports. Treat as a signal, not a guarantee.

---

## DSA Concepts Used
- **Bit Manipulation (XOR)** — pairing indices with values cancels everything but the answer → see [`/dsa/bit_manipulation.md`](/dsa/bit_manipulation.md)
- **Math / Number Theory** — Gauss's sum `n(n+1)/2` gives the expected total → see [`/dsa/math_number_theory.md`](/dsa/math_number_theory.md)
- **Sorting** — a sorted range exposes the first index/value mismatch → see [`/dsa/sorting.md`](/dsa/sorting.md)
- **Hash Map / Set** — membership probe over `0..n` → see [`/dsa/hash_map.md`](/dsa/hash_map.md)

---

## Approaches Overview

| # | Approach | Time | Space | When to use |
|---|----------|------|-------|-------------|
| 1 | Sorting | O(n log n) | O(1) | Simple, no extra structures |
| 2 | Hash Set | O(n) | O(n) | Clearest to reason about |
| 3 | Gauss Sum | O(n) | O(1) | Meets O(n)/O(1) follow-up |
| 4 | XOR | O(n) | O(1) | Follow-up, overflow-safe |

---

## Approach 1 — Sorting

### Intuition
If nothing were missing, the sorted array would be `[0,1,2,...,n]` with
`nums[i] == i` everywhere. The first index where that breaks is the missing
value. If every index matches, the gap is at the end → `n`.

### Algorithm
1. Sort `nums`.
2. Scan; return the first `i` where `nums[i] != i`.
3. If none, return `n`.

### Complexity
- **Time:** O(n log n) — dominated by the sort.
- **Space:** O(1) — in-place sort.

### Code
```go
func sortScan(nums []int) int {
	sort.Ints(nums)       // arrange 0..n (with one gap) in order
	for i, v := range nums {
		if v != i { // first index whose value slipped
			return i
		}
	}
	return len(nums) // gap is at the very end -> n is missing
}
```

### Dry Run
Input `[3,0,1]` → sorted `[0,1,3]`:

| i | nums[i] | nums[i]==i? |
|---|---------|-------------|
| 0 | 0       | yes         |
| 1 | 1       | yes         |
| 2 | 3       | no → return 2 |

Return **2**. ✅

---

## Approach 2 — Hash Set

### Intuition
Store every present number, then ask which of `0..n` is absent.

### Algorithm
1. Insert every value into a set.
2. For `i` in `0..n`: if `i` is not in the set, return `i`.

### Complexity
- **Time:** O(n) — build the set, then probe `0..n`.
- **Space:** O(n) — the set.

### Code
```go
func hashSet(nums []int) int {
	present := make(map[int]struct{}, len(nums)) // set of present values
	for _, v := range nums {
		present[v] = struct{}{}
	}
	for i := 0; i <= len(nums); i++ { // 0..n inclusive
		if _, ok := present[i]; !ok { // i never appeared
			return i
		}
	}
	return -1 // unreachable for valid input
}
```

### Dry Run
Input `[3,0,1]`, `n = 3`, set = {3,0,1}:

| i | in set? |
|---|---------|
| 0 | yes     |
| 1 | yes     |
| 2 | no → return 2 |

Return **2**. ✅

---

## Approach 3 — Gauss Sum

### Intuition
The complete set `0..n` sums to `n(n+1)/2`. The array is that set minus one
number, so `expected - actual = missing`.

### Algorithm
1. `expected = n(n+1)/2` where `n = len(nums)`.
2. `actual = sum(nums)`.
3. Return `expected - actual`.

### Complexity
- **Time:** O(n) — one pass to sum.
- **Space:** O(1).

### Code
```go
func gaussSum(nums []int) int {
	n := len(nums)
	expected := n * (n + 1) / 2 // sum of 0..n
	actual := 0
	for _, v := range nums {
		actual += v // sum of present values
	}
	return expected - actual // the gap
}
```

### Dry Run
Input `[3,0,1]`, `n = 3`:

| Quantity | Value              |
|----------|--------------------|
| expected | 3·4/2 = 6          |
| actual   | 3+0+1 = 4          |
| result   | 6 − 4 = **2**      |

Return **2**. ✅

---

## Approach 4 — XOR (Optimal)

### Intuition
XOR-ing a value with itself is 0. XOR every index `0..n` together with every
array value: each present number cancels its own index, leaving only the missing
number. Unlike summation, XOR never overflows.

### Algorithm
1. `result = n` (accounts for the top index that has no matching element).
2. For each `i`: `result ^= i ^ nums[i]`.
3. Return `result`.

### Complexity
- **Time:** O(n) — single pass.
- **Space:** O(1).

### Code
```go
func xorBits(nums []int) int {
	result := len(nums) // seed with n (index with no element)
	for i, v := range nums {
		result ^= i ^ v // cancel index against value
	}
	return result // only the missing number survives
}
```

### Dry Run
Input `[3,0,1]`, seed `result = 3`:

| i | nums[i] | result ^= i ^ nums[i] | result |
|---|---------|-----------------------|--------|
| — | —       | seed                  | 3      |
| 0 | 3       | 3 ^ 0 ^ 3 = 0         | 0      |
| 1 | 0       | 0 ^ 1 ^ 0 = 1         | 1      |
| 2 | 1       | 1 ^ 2 ^ 1 = 2         | 2      |

Return **2**. ✅

---

## Key Takeaways
- **XOR self-cancellation** turns "find the odd one out" problems into O(1)-space single passes.
- **Gauss's formula** gives the same result but risks integer overflow for large ranges; XOR does not.
- Seeding the accumulator with `n` neatly covers the highest index that has no array element.
- Sorting and hashing are fine fallbacks but miss the O(n)/O(1) follow-up target.

---

## Related Problems
- LeetCode #136 — Single Number (XOR cancellation)
- LeetCode #287 — Find the Duplicate Number (opposite: an extra rather than a gap)
- LeetCode #41 — First Missing Positive (missing element with cyclic placement)
- LeetCode #448 — Find All Numbers Disappeared in an Array (multiple gaps)
