# 0338 — Counting Bits

> LeetCode #338 · Difficulty: Easy
> **Categories:** Dynamic Programming, Bit Manipulation

---

## Problem Statement

Given an integer `n`, return *an array `ans` of length `n + 1` such that for each `i` (`0 <= i <= n`), `ans[i]` is the **number of `1`'s** in the binary representation of `i`*.

**Example 1:**

```
Input: n = 2
Output: [0,1,1]
Explanation:
0 --> 0
1 --> 1
2 --> 10
```

**Example 2:**

```
Input: n = 5
Output: [0,1,1,2,1,2]
Explanation:
0 --> 0
1 --> 1
2 --> 10
3 --> 11
4 --> 100
5 --> 101
```

**Constraints:**

- `0 <= n <= 10^5`

**Follow up:**

- It is very easy to come up with a solution with a runtime of `O(n log n)`. Can you do it in linear time `O(n)` and possibly in a single pass?
- Can you do it without using any built-in function (i.e., like `__builtin_popcount` in C++)?

---

## Company Frequency

| Company    | Frequency        | Last Reported |
|------------|------------------|---------------|
| Amazon     | ★★★★☆ High       | 2024          |
| Google     | ★★★☆☆ Medium     | 2024          |
| Microsoft  | ★★★☆☆ Medium     | 2023          |
| Apple      | ★★☆☆☆ Low        | 2023          |
| Meta       | ★★☆☆☆ Low        | 2022          |

> ⚠️ Frequency data is crowd-sourced from LeetCode Discuss, Glassdoor, and
> community interview reports. Treat as a signal, not a guarantee.

---

## DSA Concepts Used

- **Bit Manipulation** — every solution leans on a bit identity: `x & (x-1)` clears the lowest set bit, `x >> 1` drops it, and the highest power of two carries one bit → see [`/dsa/bit_manipulation.md`](/dsa/bit_manipulation.md)
- **Dynamic Programming (1D)** — the linear solutions express `ans[i]` in terms of a smaller, already-computed `ans[j]` with `j < i` → see [`/dsa/dynamic_programming_1d.md`](/dsa/dynamic_programming_1d.md)

---

## Approaches Overview

| # | Approach | Time | Space | When to use |
|---|----------|------|-------|-------------|
| 1 | Brute Force (Kernighan per number) | O(n log n) | O(1) extra | Baseline; satisfies the easy `O(n log n)` bound |
| 2 | DP with Highest Power of Two | O(n) | O(1) extra | Linear; intuitive "strip the high bit" recurrence |
| 3 | DP with Right Shift | O(n) | O(1) extra | Linear; `ans[i>>1] + (i&1)` — very clean |
| 4 | DP with Kernighan (Optimal) | O(n) | O(1) extra | Linear; `ans[i&(i-1)] + 1` — the tightest recurrence |
| 5 | Library popcount (reference) | O(n) | O(1) extra | Verification only; uses a built-in |

---

## Approach 1 — Brute Force (Kernighan per number)

### Intuition

The problem literally asks for `popcount(i)` for every `i` in `[0, n]`. The most direct route is to count the set bits of each number independently. Brian Kernighan's identity `x &= x-1` removes exactly the lowest set bit, so the loop runs once per set bit — cheaper than scanning all 32 positions.

### Algorithm

1. Allocate `ans` of length `n+1`.
2. For each `i` in `[0, n]`, set `x = i, count = 0`.
3. While `x > 0`: do `x &= x-1` and `count++`.
4. Store `ans[i] = count`.

### Complexity

- **Time:** O(n log n) — for each of the `n+1` numbers the inner loop runs `popcount(i)` times (≤ 32).
- **Space:** O(1) extra — only counters beyond the required output array.

### Code

```go
func bruteForce(n int) []int {
	ans := make([]int, n+1) // ans[i] will hold popcount(i)
	for i := 0; i <= n; i++ {
		count := 0 // number of set bits found in i so far
		x := i     // work on a copy so i stays intact
		for x > 0 {
			x &= x - 1 // Kernighan: erase the lowest set bit
			count++    // each erase corresponds to exactly one set bit
		}
		ans[i] = count // record the popcount of i
	}
	return ans
}
```

### Dry Run

Example 1: `n = 2`.

| i | x start | inner steps (x &= x-1) | count | ans[i] |
|---|---------|------------------------|-------|--------|
| 0 | 0 | (loop skipped)                | 0 | 0 |
| 1 | 1 | `1 & 0 = 0` → count 1          | 1 | 1 |
| 2 | 2 (`10`) | `10 & 01 = 0` → count 1     | 1 | 1 |

Result: `[0, 1, 1]` ✔

---

## Approach 2 — DP with Highest Power of Two

### Intuition

Let `offset` be the largest power of two that is `≤ i`. Then `i` differs from `i - offset` by exactly that one high bit, so `ans[i] = ans[i-offset] + 1`. `offset` doubles precisely when `i` hits the next power of two.

### Algorithm

1. Set `offset = 1`, `ans[0] = 0`.
2. For `i` from `1` to `n`: if `offset*2 == i`, set `offset *= 2`.
3. Set `ans[i] = ans[i-offset] + 1`.

### Complexity

- **Time:** O(n) — one O(1) transition per number, single pass.
- **Space:** O(1) extra beyond the output array.

### Code

```go
func dpHighBit(n int) []int {
	ans := make([]int, n+1) // ans[0] = 0 by default
	offset := 1             // highest power of two seen so far (starts at 1 = 2^0)
	for i := 1; i <= n; i++ {
		if offset*2 == i {
			offset *= 2 // i just reached the next power of two → update the high bit
		}
		ans[i] = ans[i-offset] + 1 // i = (i-offset) plus one extra high bit
	}
	return ans
}
```

### Dry Run

Example 1: `n = 2`.

| i | offset check | offset | i-offset | ans[i-offset] | ans[i] |
|---|--------------|--------|----------|---------------|--------|
| 1 | 1*2==1? no  | 1 | 0 | 0 | 1 |
| 2 | 1*2==2? yes → offset=2 | 2 | 0 | 0 | 1 |

Result: `[0, 1, 1]` ✔

---

## Approach 3 — DP with Right Shift

### Intuition

`i >> 1` removes `i`'s lowest bit; that bit is `i & 1`. Therefore `popcount(i) = popcount(i >> 1) + (i & 1)`. Because `i>>1 < i`, its answer is already stored.

### Algorithm

1. Set `ans[0] = 0`.
2. For `i` from `1` to `n`: set `ans[i] = ans[i>>1] + (i & 1)`.

### Complexity

- **Time:** O(n) — one O(1) transition per number.
- **Space:** O(1) extra beyond the output array.

### Code

```go
func dpRightShift(n int) []int {
	ans := make([]int, n+1)
	for i := 1; i <= n; i++ {
		ans[i] = ans[i>>1] + (i & 1) // half's popcount plus the bit we shifted off
	}
	return ans
}
```

### Dry Run

Example 1: `n = 2`.

| i | i>>1 | ans[i>>1] | i & 1 | ans[i] |
|---|------|-----------|-------|--------|
| 1 | 0 | 0 | 1 | 1 |
| 2 | 1 | 1 | 0 | 1 |

Result: `[0, 1, 1]` ✔

---

## Approach 4 — DP with Kernighan (Optimal)

### Intuition

`i & (i-1)` equals `i` with its lowest set bit cleared — a strictly smaller, already-solved number that has exactly one fewer 1-bit. Hence `ans[i] = ans[i & (i-1)] + 1`. This is the tightest of the recurrences and needs no `offset` bookkeeping.

### Algorithm

1. Set `ans[0] = 0`.
2. For `i` from `1` to `n`: set `ans[i] = ans[i & (i-1)] + 1`.

### Complexity

- **Time:** O(n) — single pass, O(1) per number.
- **Space:** O(1) extra beyond the output array.

### Code

```go
func dpKernighan(n int) []int {
	ans := make([]int, n+1)
	for i := 1; i <= n; i++ {
		ans[i] = ans[i&(i-1)] + 1 // clear lowest set bit → one fewer 1-bit
	}
	return ans
}
```

### Dry Run

Example 1: `n = 2`.

| i | i-1 | i & (i-1) | ans[i&(i-1)] | ans[i] |
|---|-----|-----------|--------------|--------|
| 1 | 0 | `1 & 0 = 0` | 0 | 1 |
| 2 | 1 | `10 & 01 = 0` | 0 | 1 |

Result: `[0, 1, 1]` ✔

---

## Approach 5 — Library popcount (reference)

### Intuition

`math/bits.OnesCount` maps to a hardware POPCNT-style instruction, giving `popcount(i)` in O(1). Useful to cross-check the DP outputs, but it sidesteps the follow-up's "no built-in" constraint, so it is a reference rather than the intended answer.

### Algorithm

1. For `i` in `[0, n]`: set `ans[i] = bits.OnesCount(uint(i))`.

### Complexity

- **Time:** O(n) — each `OnesCount` is O(1).
- **Space:** O(1) extra beyond the output array.

### Code

```go
func libPopcount(n int) []int {
	ans := make([]int, n+1)
	for i := 0; i <= n; i++ {
		ans[i] = bits.OnesCount(uint(i)) // constant-time hardware popcount
	}
	return ans
}
```

### Dry Run

Example 1: `n = 2`.

| i | binary | bits.OnesCount | ans[i] |
|---|--------|----------------|--------|
| 0 | `0`  | 0 | 0 |
| 1 | `1`  | 1 | 1 |
| 2 | `10` | 1 | 1 |

Result: `[0, 1, 1]` ✔

---

## Key Takeaways

- **Reuse a smaller subproblem.** All three linear solutions build `ans[i]` from an already-computed `ans[j]` with `j < i` — classic 1D DP over integers.
- **Three bit identities, one shape:** strip the high bit (`i-offset`), strip the low bit via shift (`i>>1` + `i&1`), or clear the low set bit (`i & (i-1)`). All give `+1` transitions.
- **`x & (x-1)` clears the lowest set bit** — the same idiom that powers #191 and #201. It yields the cleanest recurrence here.
- The `O(n log n)` brute force already passes; the value of the DP is hitting the follow-up's **single-pass linear** bound with no built-ins.

---

## Related Problems

- LeetCode #191 — Number of 1 Bits (single-number popcount, Kernighan loop)
- LeetCode #201 — Bitwise AND of Numbers Range (`x & (x-1)` trick)
- LeetCode #231 — Power of Two (`x & (x-1) == 0`)
- LeetCode #260 — Single Number III (bit partitioning)
