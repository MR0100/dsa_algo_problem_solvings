# 0477 — Total Hamming Distance

> LeetCode #477 · Difficulty: Medium
> **Categories:** Bit Manipulation, Array, Math

---

## Problem Statement

The **Hamming distance** between two integers is the number of positions at which the corresponding bits are different.

Given an integer array `nums`, return *the sum of **Hamming distances** between all the pairs of the integers in* `nums`.

**Example 1:**

```
Input: nums = [4,14,2]
Output: 6
Explanation: In binary representation, the 4 is 0100, 14 is 1110, and 2 is 0010 (just
showing the four bits relevant in this case).
The answer will be:
HammingDistance(4, 14) + HammingDistance(4, 2) + HammingDistance(14, 2) = 2 + 2 + 2 = 6.
```

**Example 2:**

```
Input: nums = [4,14,4]
Output: 4
```

**Constraints:**

- `1 <= nums.length <= 10^4`
- `0 <= nums[i] <= 10^9`
- The answer for the given input will fit in a **32-bit** integer.

---

## Company Frequency

| Company    | Frequency        | Last Reported |
|------------|------------------|---------------|
| Facebook   | ★★★★☆ High       | 2024          |
| Amazon     | ★★★☆☆ Medium     | 2023          |
| Google     | ★★★☆☆ Medium     | 2023          |
| Microsoft  | ★★☆☆☆ Low        | 2022          |
| Uber       | ★★☆☆☆ Low        | 2022          |

> ⚠️ Frequency data is crowd-sourced from LeetCode Discuss, Glassdoor, and
> community interview reports. Treat as a signal, not a guarantee.

---

## DSA Concepts Used

- **Bit Manipulation** — Hamming distance is `popcount(x ^ y)`, and the optimal solution decomposes the total into independent per-bit-column contributions → see [`/dsa/bit_manipulation.md`](/dsa/bit_manipulation.md)
- **Counting / Combinatorics** — the trick is *swap the order of summation*: for each bit column, the number of differing pairs is `ones × zeros`, a pure counting argument on the array values → see [`/dsa/arrays.md`](/dsa/arrays.md)

---

## Approaches Overview

| # | Approach | Time | Space | When to use |
|---|----------|------|-------|-------------|
| 1 | Brute Force (All Pairs) | O(n² · 32) | O(1) | Baseline; `n = 10⁴` → ~10⁸ pair-ops, TLE |
| 2 | Bit-Position Counting (Optimal) | O(n · 32) = O(n) | O(1) | The intended answer; per-column `ones × zeros` |
| 3 | Bit-Position Counting, Early Exit | O(n · maxBit) | O(1) | Same, but skips all-zero high columns |

---

## Approach 1 — Brute Force (All Pairs)

### Intuition

The problem is defined as a sum over all unordered pairs. Just enumerate each pair `(i, j)` once and add its Hamming distance, which is the popcount of the XOR (`x ^ y` is `1` exactly where the two numbers disagree).

### Algorithm

1. `total = 0`.
2. For `i` in `0..n-1`, for `j` in `i+1..n-1`: `total += popcount(nums[i] ^ nums[j])`.
3. Return `total`.

### Complexity

- **Time:** O(n² · w) with word width `w ≤ 32` — there are `n(n-1)/2` pairs, each a constant-width popcount. At `n = 10⁴` this is ~10⁸ and times out.
- **Space:** O(1).

### Code

```go
func hammingDistance(x, y int) int {
	return bits.OnesCount(uint(x ^ y)) // popcount of the disagreement mask
}

func bruteForce(nums []int) int {
	total := 0
	for i := 0; i < len(nums); i++ { // pick the first element of the pair
		for j := i + 1; j < len(nums); j++ { // pair it with every later element
			total += hammingDistance(nums[i], nums[j]) // add this pair's distance
		}
	}
	return total
}
```

### Dry Run

Example 1: `nums = [4, 14, 2]` → `4 = 0100`, `14 = 1110`, `2 = 0010`.

| Pair (i, j) | nums[i] ^ nums[j] | popcount | running total |
|-------------|-------------------|----------|---------------|
| (0, 1) 4,14 | `0100 ^ 1110 = 1010` | 2 | 2 |
| (0, 2) 4,2  | `0100 ^ 0010 = 0110` | 2 | 4 |
| (1, 2) 14,2 | `1110 ^ 0010 = 1100` | 2 | 6 |

Result: `6` ✔

---

## Approach 2 — Bit-Position Counting (Optimal)

### Intuition

Rewrite the double sum by **swapping the order of summation**. Instead of "for each pair, count differing bits", do "for each bit position, count differing pairs". Fix a bit column `b`. Suppose `ones` of the `n` numbers have a `1` there; the other `zeros = n - ones` have a `0`. A pair differs at column `b` **iff** one member is from the `ones` group and the other from the `zeros` group — so column `b` contributes exactly `ones × zeros` to the grand total (each such pair adds `1`). Summing `ones × zeros` over all 32 columns gives the answer in linear time.

### Algorithm

1. `n = len(nums)`, `total = 0`.
2. For each bit column `b` in `0..31`:
   - `ones = ` number of `nums` with bit `b` set.
   - `total += ones × (n - ones)`.
3. Return `total`.

### Complexity

- **Time:** O(n · w) with `w = 32`, i.e. O(n) — one pass over the array per column.
- **Space:** O(1).

### Code

```go
func bitColumnCount(nums []int) int {
	n := len(nums)
	total := 0
	for b := 0; b < 32; b++ { // examine each bit column independently
		ones := 0 // how many numbers have a 1 in this column
		for _, x := range nums {
			ones += (x >> b) & 1 // add 1 when bit b of x is set
		}
		zeros := n - ones     // the rest have a 0 in this column
		total += ones * zeros // each (one,zero) pair differs here → contributes 1
	}
	return total
}
```

### Dry Run

Example 1: `nums = [4, 14, 2]` = `[0100, 1110, 0010]`, `n = 3`. Only columns with any `1` contribute.

| Column b | bits of (4,14,2) | ones | zeros = 3 − ones | ones × zeros | running total |
|----------|------------------|------|------------------|--------------|---------------|
| 0 | (0,0,0) | 0 | 3 | 0 | 0 |
| 1 | (0,1,1) | 2 | 1 | 2 | 2 |
| 2 | (1,1,0) | 2 | 1 | 2 | 4 |
| 3 | (0,1,0) | 1 | 2 | 2 | 6 |
| 4..31 | all 0 | 0 | 3 | 0 | 6 |

Result: `6` ✔

---

## Approach 3 — Bit-Position Counting, Early Exit

### Intuition

Values are `≤ 10⁹ < 2³⁰`, so most of the 32 columns are all-zero and contribute `0 × n = 0`. OR every number together first: the highest set bit of that union is the topmost column worth scanning. Iterate columns only up to it. The asymptotics are unchanged, but the constant factor drops for arrays of small numbers.

### Algorithm

1. `orAll = OR of all nums`.
2. For `b = 0` while `(orAll >> b) != 0`: accumulate `ones × (n - ones)` for column `b`.
3. Return `total`.

### Complexity

- **Time:** O(n · maxBit) where `maxBit` is the index of the highest set bit across the array (`≤ 30` here).
- **Space:** O(1).

### Code

```go
func bitColumnEarlyExit(nums []int) int {
	n := len(nums)
	orAll := 0
	for _, x := range nums {
		orAll |= x // union of all bits present anywhere in the array
	}
	total := 0
	for b := 0; orAll>>b != 0; b++ { // stop once no number has a bit at/above b
		ones := 0
		for _, x := range nums {
			ones += (x >> b) & 1 // count 1s in column b
		}
		total += ones * (n - ones) // differing pairs contributed by column b
	}
	return total
}
```

### Dry Run

Example 1: `nums = [4, 14, 2]`. `orAll = 4 | 14 | 2 = 1110₂ = 14`, so the loop runs for `b = 0, 1, 2, 3` (since `14 >> 4 = 0` stops it).

| Column b | orAll >> b | ones | zeros | contribution | total |
|----------|-----------|------|-------|--------------|-------|
| 0 | `1110` ≠ 0 | 0 | 3 | 0 | 0 |
| 1 | `111`  ≠ 0 | 2 | 1 | 2 | 2 |
| 2 | `11`   ≠ 0 | 2 | 1 | 2 | 4 |
| 3 | `1`    ≠ 0 | 1 | 2 | 2 | 6 |
| 4 | `0` → stop | — | — | —  | 6 |

Result: `6` ✔

---

## Key Takeaways

- **Swap the order of summation.** "Sum over pairs of a per-bit quantity" becomes "sum over bits of a per-pair count" — the master move that turns O(n²) into O(n).
- At any bit column, differing pairs = `ones × zeros`. This *contribution counting* pattern recurs whenever a total decomposes into independent coordinates (bits, characters, colors).
- `popcount(x ^ y)` is the one-liner for Hamming distance between two integers.
- Iterating fixed 32 columns keeps the algorithm O(n) regardless of value magnitude; OR-ing first only trims constants.

---

## Related Problems

- LeetCode #461 — Hamming Distance (single pair, `popcount(x ^ y)`)
- LeetCode #191 — Number of 1 Bits (popcount building block)
- LeetCode #201 — Bitwise AND of Numbers Range (per-bit reasoning over a set)
- LeetCode #1863 — Sum of All Subset XOR Totals (per-bit contribution counting)
- LeetCode #2588 — Count the Number of Beautiful Subarrays (bit-column parity)
