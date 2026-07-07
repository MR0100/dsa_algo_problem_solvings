# 0136 — Single Number

> LeetCode #136 · Difficulty: Easy
> **Categories:** Array, Bit Manipulation

---

## Problem Statement

Given a **non-empty** array of integers `nums`, every element appears **twice** except for one. Find that single one.

You must implement a solution with a linear runtime complexity and use only constant extra space.

**Example 1:**
```
Input: nums = [2,2,1]
Output: 1
```

**Example 2:**
```
Input: nums = [4,1,2,1,2]
Output: 4
```

**Example 3:**
```
Input: nums = [1]
Output: 1
```

**Constraints:**
- `1 <= nums.length <= 3 * 10^4`
- `-3 * 10^4 <= nums[i] <= 3 * 10^4`
- Each element in the array appears twice except for one element which appears only once.

---

## Company Frequency

| Company   | Frequency       | Last Reported |
|-----------|-----------------|---------------|
| Amazon    | ★★★★★ Very High | 2024          |
| Google    | ★★★★☆ High      | 2024          |
| Apple     | ★★★★☆ High      | 2024          |
| Microsoft | ★★★★☆ High      | 2024          |
| Adobe     | ★★★☆☆ Medium    | 2023          |
| Bloomberg | ★★★☆☆ Medium    | 2023          |

> ⚠️ Frequency data is crowd-sourced from LeetCode Discuss, Glassdoor, and community interview reports. Treat as a signal, not a guarantee.

---

## DSA Concepts Used

- **Bit Manipulation** — XOR's self-cancelling property (`a ⊕ a = 0`) makes pairs vanish, leaving the loner → see [`/dsa/bit_manipulation.md`](/dsa/bit_manipulation.md)
- **Hash Map** — frequency counting is the generic "find the odd one out" tool → see [`/dsa/hash_map.md`](/dsa/hash_map.md)
- **Sorting** — sorting groups duplicates adjacently, enabling a pair-wise scan → see [`/dsa/sorting.md`](/dsa/sorting.md)

---

## Approaches Overview

| # | Approach            | Time       | Space | When to use                                          |
|---|---------------------|------------|-------|------------------------------------------------------|
| 1 | Brute Force         | O(n²)      | O(1)  | Never in practice; baseline to explain the problem   |
| 2 | Hash Map            | O(n)       | O(n)  | When elements may repeat any number of times         |
| 3 | Sorting             | O(n log n) | O(n)  | When mutation is allowed and memory is tight-ish     |
| 4 | Math (Set Sum)      | O(n)       | O(n)  | Neat trick; generalizes to "appears k times"         |
| 5 | Bitwise XOR (Optimal) | O(n)     | O(1)  | Always — meets the required O(n)/O(1) follow-up      |

---

## Approach 1 — Brute Force

### Intuition
The most direct reading of the problem: an element is "the single one" if its total occurrence count in the array is exactly 1. So for every element, count how often it appears.

### Algorithm
1. For each index `i` from `0` to `n-1`:
   1. Scan the whole array and count how many `j` satisfy `nums[j] == nums[i]`.
   2. If the count is `1`, return `nums[i]`.

### Complexity
- **Time:** O(n²) — each of the n candidates triggers a full O(n) counting scan.
- **Space:** O(1) — only loop counters and one count variable.

### Code
```go
func bruteForce(nums []int) int {
	for i := 0; i < len(nums); i++ { // candidate element
		count := 0
		for j := 0; j < len(nums); j++ { // count its occurrences everywhere
			if nums[j] == nums[i] {
				count++
			}
		}
		if count == 1 { // appears exactly once → it is the single number
			return nums[i]
		}
	}
	return -1 // unreachable per problem guarantee
}
```

### Dry Run
`nums = [2,2,1]` (Example 1):

| i | nums[i] | occurrences found | count | action        |
|---|---------|-------------------|-------|---------------|
| 0 | 2       | j=0 (2), j=1 (2)  | 2     | keep looking  |
| 1 | 2       | j=0 (2), j=1 (2)  | 2     | keep looking  |
| 2 | 1       | j=2 (1)           | 1     | return **1** ✅ |

---

## Approach 2 — Hash Map

### Intuition
Counting is what hash maps do best. One pass to tally frequencies, one pass over the table to find the key whose count is 1.

### Algorithm
1. Create an empty map `freq` from value → count.
2. For every `num` in `nums`, do `freq[num]++`.
3. Iterate the map; return the key whose count equals `1`.

### Complexity
- **Time:** O(n) — one pass to build the map plus one pass over at most n keys, each map op amortized O(1).
- **Space:** O(n) — the map stores up to `(n+1)/2` distinct values.

### Code
```go
func hashMap(nums []int) int {
	freq := make(map[int]int, len(nums)) // value → occurrence count
	for _, num := range nums {
		freq[num]++ // tally every element
	}
	for num, count := range freq {
		if count == 1 { // the unique element
			return num
		}
	}
	return -1 // unreachable per problem guarantee
}
```

### Dry Run
`nums = [2,2,1]` (Example 1):

| step | num read | freq map after step |
|------|----------|---------------------|
| 1    | 2        | {2:1}               |
| 2    | 2        | {2:2}               |
| 3    | 1        | {2:2, 1:1}          |
| scan | —        | key 2 has count 2 (skip); key 1 has count 1 → return **1** ✅ |

---

## Approach 3 — Sorting

### Intuition
After sorting, both copies of every duplicate are adjacent. Walking the array in steps of two, every well-formed pair satisfies `arr[i] == arr[i+1]`. The first index where that fails is where the single number slid in and shifted the pairing.

### Algorithm
1. Copy the array (to leave the caller's data untouched) and sort it.
2. For `i = 0, 2, 4, ...` while `i+1 < n`: if `arr[i] != arr[i+1]`, return `arr[i]`.
3. If the loop finishes, every pair matched, so the single number is the last element — return `arr[n-1]`.

### Complexity
- **Time:** O(n log n) — the sort dominates the O(n) scan.
- **Space:** O(n) — for the defensive copy (O(1) extra if sorting in place is acceptable).

### Code
```go
func sortAndScan(nums []int) int {
	arr := make([]int, len(nums)) // copy so we don't mutate the input
	copy(arr, nums)
	sort.Ints(arr) // duplicates become adjacent
	for i := 0; i+1 < len(arr); i += 2 {
		if arr[i] != arr[i+1] { // pairing broken → arr[i] is alone
			return arr[i]
		}
	}
	return arr[len(arr)-1] // all pairs matched → the last element is the single one
}
```

### Dry Run
`nums = [2,2,1]` (Example 1), sorted copy `arr = [1,2,2]`:

| i | arr[i] | arr[i+1] | equal? | action |
|---|--------|----------|--------|--------|
| 0 | 1      | 2        | no     | return **1** ✅ |

---

## Approach 4 — Math (Set Sum)

### Intuition
If every value appeared exactly twice, the array total would be exactly twice the sum of the distinct values. Our single number contributes once instead of twice, so it is short by exactly one copy of itself:

`2 · sum(distinct) − sum(all) = single number`

### Algorithm
1. Walk the array once, accumulating `sumAll` over every element.
2. Simultaneously insert values into a set; the first time a value is seen, add it to `sumSet`.
3. Return `2*sumSet − sumAll`.

### Complexity
- **Time:** O(n) — a single pass with O(1) set operations.
- **Space:** O(n) — the set of distinct values.

### Code
```go
func mathSum(nums []int) int {
	seen := make(map[int]bool, len(nums)) // distinct values
	sumAll, sumSet := 0, 0
	for _, num := range nums {
		sumAll += num // sum of everything, duplicates included
		if !seen[num] {
			seen[num] = true
			sumSet += num // sum of each distinct value once
		}
	}
	return 2*sumSet - sumAll // duplicates cancel, the single number remains
}
```

### Dry Run
`nums = [2,2,1]` (Example 1):

| step | num | seen after      | sumAll | sumSet |
|------|-----|-----------------|--------|--------|
| 1    | 2   | {2}             | 2      | 2      |
| 2    | 2   | {2}             | 4      | 2      |
| 3    | 1   | {2,1}           | 5      | 3      |

Result: `2*3 − 5 = 1` → **1** ✅

---

## Approach 5 — Bitwise XOR (Optimal)

### Intuition
XOR has three magic properties: it is commutative, associative, `a ⊕ a = 0`, and `a ⊕ 0 = a`. So XOR-ing the entire array lets us mentally reorder it into pairs — each pair annihilates to 0 — and the single number XOR-ed with 0 is itself. This meets the problem's O(n) time / O(1) space requirement exactly.

### Algorithm
1. Initialize `result = 0` (the XOR identity).
2. For every `num` in `nums`, set `result ^= num`.
3. Return `result`.

### Complexity
- **Time:** O(n) — one pass, one XOR per element.
- **Space:** O(1) — a single integer accumulator.

### Code
```go
func bitwiseXOR(nums []int) int {
	result := 0 // XOR identity
	for _, num := range nums {
		result ^= num // pairs cancel to 0; the loner survives
	}
	return result
}
```

### Dry Run
`nums = [2,2,1]` (Example 1). In binary: 2 = `10`, 1 = `01`.

| step | num | result before | XOR operation   | result after |
|------|-----|---------------|-----------------|--------------|
| 1    | 2   | 00            | 00 ⊕ 10 = 10    | 2            |
| 2    | 2   | 10            | 10 ⊕ 10 = 00    | 0            |
| 3    | 1   | 00            | 00 ⊕ 01 = 01    | 1            |

Return **1** ✅

---

## Key Takeaways

- **XOR cancels pairs**: `a ⊕ a = 0`, `a ⊕ 0 = a`, plus commutativity — the canonical trick for "everything appears twice except one". Memorize it.
- The same XOR idea extends to *Single Number III* (two loners): XOR everything, then split the array by any set bit of the combined XOR.
- The **set-sum identity** `k·sum(set) − sum(all)` generalizes to "every element appears k times except one" — that is exactly one accepted solution for LeetCode #137.
- Sorting is the universal fallback for duplicate-grouping problems, at an O(n log n) cost.
- When a problem explicitly says "linear time, constant space", it is usually hinting at bit tricks or in-place index games.

---

## Related Problems

- LeetCode #137 — Single Number II (every element appears three times except one)
- LeetCode #260 — Single Number III (two elements appear once)
- LeetCode #268 — Missing Number (XOR over indices and values)
- LeetCode #389 — Find the Difference (XOR over two strings)
- LeetCode #287 — Find the Duplicate Number (constant-space duplicate finding)
