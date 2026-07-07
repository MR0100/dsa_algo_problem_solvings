# 0260 — Single Number III

> LeetCode #260 · Difficulty: Medium
> **Categories:** Array, Bit Manipulation, Hash Table

---

## Problem Statement

Given an integer array `nums`, in which exactly two elements appear only once and all the other elements appear exactly twice. Find the two elements that appear only once. You can return the answer in **any order**.

You must write an algorithm that runs in linear runtime complexity and uses only constant extra space.

**Example 1:**

```
Input: nums = [1,2,1,3,2,5]
Output: [3,5]
Explanation: [5, 3] is also a valid answer.
```

**Example 2:**

```
Input: nums = [-1,0]
Output: [-1,0]
```

**Example 3:**

```
Input: nums = [0,1]
Output: [1,0]
```

**Constraints:**

- `2 <= nums.length <= 3 * 10^4`
- `-2^31 <= nums[i] <= 2^31 - 1`
- Each integer in `nums` will appear twice, only two integers will appear once.

---

## Company Frequency

| Company    | Frequency        | Last Reported |
|------------|------------------|---------------|
| Amazon     | ★★★☆☆ Medium     | 2024          |
| Google     | ★★★☆☆ Medium     | 2023          |
| Microsoft  | ★★★☆☆ Medium     | 2023          |
| Facebook   | ★★☆☆☆ Low        | 2023          |
| Bloomberg  | ★★☆☆☆ Low        | 2022          |

> ⚠️ Frequency data is crowd-sourced from LeetCode Discuss, Glassdoor, and
> community interview reports. Treat as a signal, not a guarantee.

---

## DSA Concepts Used

- **Bit Manipulation (XOR)** — XORing all elements cancels the paired numbers and leaves `a ^ b`; a differing bit then partitions the array so each unique can be isolated in constant space → see [`/dsa/bit_manipulation.md`](/dsa/bit_manipulation.md)
- **Hash Map** — the O(n)-space baseline counts frequencies and returns the count-1 keys → see [`/dsa/hash_map.md`](/dsa/hash_map.md)

---

## Approaches Overview

| # | Approach | Time | Space | When to use |
|---|----------|------|-------|-------------|
| 1 | Hash Map Counting | O(n) | O(n) | Simplest; violates the O(1)-space requirement but easy to reason about |
| 2 | XOR Partition (Optimal) | O(n) | O(1) | The required linear-time, constant-space solution |

---

## Approach 1 — Hash Map Counting

### Intuition

Every value appears twice except the two we want, which appear once. Count frequencies in a map and return the two keys with count 1.

### Algorithm

1. Tally each value's frequency into a map.
2. Collect keys whose count is exactly 1 (there are exactly two).
3. Return them.

### Complexity

- **Time:** O(n) — one counting pass plus one pass over the map.
- **Space:** O(n) — the frequency map (does **not** meet the O(1) constraint).

### Code

```go
func hashMap(nums []int) []int {
	count := make(map[int]int)
	for _, x := range nums { // tally frequencies
		count[x]++
	}
	res := make([]int, 0, 2)
	for x, c := range count {
		if c == 1 { // appears exactly once
			res = append(res, x)
		}
	}
	sort.Ints(res) // stable output for testing (map order is random)
	return res
}
```

### Dry Run

Example 1: `nums = [1,2,1,3,2,5]`.

| step | value | count map |
|------|-------|-----------|
| 1 | 1 | {1:1} |
| 2 | 2 | {1:1, 2:1} |
| 3 | 1 | {1:2, 2:1} |
| 4 | 3 | {1:2, 2:1, 3:1} |
| 5 | 2 | {1:2, 2:2, 3:1} |
| 6 | 5 | {1:2, 2:2, 3:1, 5:1} |

Keys with count 1: `3` and `5`. Result (sorted): `[3, 5]` ✔

---

## Approach 2 — XOR Partition by Differing Bit (Optimal)

### Intuition

XOR of the whole array cancels every duplicated pair (`x ^ x = 0`), leaving `xorAll = a ^ b`, where `a` and `b` are the two uniques. Since `a != b`, `xorAll` has at least one set bit. Take its **lowest** set bit via `xorAll & -xorAll`. That bit is `1` in exactly one of `a`, `b`. Use it to split *all* numbers into two groups; duplicates stay together (they share every bit), so each group contains one unique plus some pairs. XOR each group independently — the pairs cancel and the two uniques fall out.

### Algorithm

1. `xorAll = XOR of all nums` (equals `a ^ b`).
2. `diff = xorAll & -xorAll` — isolate the lowest differing bit.
3. Walk the array: if `x & diff != 0`, XOR into `a`; else XOR into `b`.
4. Return `{a, b}`.

### Complexity

- **Time:** O(n) — two linear passes.
- **Space:** O(1) — a handful of integer accumulators; meets the constraint.

### Code

```go
func xorPartition(nums []int) []int {
	xorAll := 0
	for _, x := range nums { // pairs cancel; left with a ^ b
		xorAll ^= x
	}
	// Isolate the lowest bit where a and b differ (two's-complement trick).
	diff := xorAll & (-xorAll)
	a, b := 0, 0
	for _, x := range nums {
		if x&diff != 0 { // this number has the differing bit set
			a ^= x // group A: pairs cancel, leaves one unique
		} else {
			b ^= x // group B: leaves the other unique
		}
	}
	res := []int{a, b}
	sort.Ints(res) // stable output for testing
	return res
}
```

### Dry Run

Example 1: `nums = [1,2,1,3,2,5]`. Uniques are 3 and 5.

Pass 1 — `xorAll`:
`1^2^1^3^2^5` = `(1^1)^(2^2)^3^5` = `0^0^3^5` = `3 ^ 5` = `011 ^ 101` = `110` (6).

`diff = 6 & -6 = 0b110 & 0b...010 = 0b010` (2) — the lowest set bit of `xorAll`.

Pass 2 — partition on bit `2`:

| x | x (bin) | x & 2 ? | group | a | b |
|---|---------|---------|-------|---|---|
| 1 | 001 | 0 | B | 0 | 1 |
| 2 | 010 | 2 | A | 2 | 1 |
| 1 | 001 | 0 | B | 2 | 0 |
| 3 | 011 | 2 | A | 1 | 0 |
| 2 | 010 | 2 | A | 3 | 0 |
| 5 | 101 | 0 | B | 3 | 5 |

Final `a = 3`, `b = 5`. Result (sorted): `[3, 5]` ✔

---

## Key Takeaways

- **XOR is the "cancel the pairs" hammer:** any problem where duplicates appear an even number of times collapses under `^`. Single Number I is just `XOR everything`; this problem adds one separation step.
- **`x & -x` isolates the lowest set bit** — the standard two's-complement idiom. Here it gives a bit guaranteed to differ between the two answers, which is all we need to partition.
- **Partition-then-reduce:** when one XOR leaves a *combination* of two unknowns, find a bit that separates them and rerun the reduction on each half.
- The hash-map approach is O(n) time but O(n) space — mention it, then upgrade to XOR to satisfy the constant-space constraint.

---

## Related Problems

- LeetCode #136 — Single Number (one unique; plain XOR)
- LeetCode #137 — Single Number II (every element thrice except one; bit-count mod 3)
- LeetCode #268 — Missing Number (XOR indices with values)
- LeetCode #389 — Find the Difference (XOR two strings)
