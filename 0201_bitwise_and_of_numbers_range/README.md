# 0201 — Bitwise AND of Numbers Range

> LeetCode #201 · Difficulty: Medium
> **Categories:** Bit Manipulation

---

## Problem Statement

Given two integers `left` and `right` that represent the range `[left, right]`, return *the bitwise AND of all numbers in this range, inclusive*.

**Example 1:**

```
Input: left = 5, right = 7
Output: 4
```

**Example 2:**

```
Input: left = 0, right = 0
Output: 0
```

**Example 3:**

```
Input: left = 1, right = 2147483647
Output: 0
```

**Constraints:**

- `0 <= left <= right <= 2^31 - 1`

---

## Company Frequency

| Company    | Frequency        | Last Reported |
|------------|------------------|---------------|
| Google     | ★★★★☆ High       | 2024          |
| Amazon     | ★★★☆☆ Medium     | 2024          |
| Microsoft  | ★★★☆☆ Medium     | 2023          |
| Meta       | ★★★☆☆ Medium     | 2023          |
| Apple      | ★★☆☆☆ Low        | 2023          |
| Bloomberg  | ★★☆☆☆ Low        | 2022          |

> ⚠️ Frequency data is crowd-sourced from LeetCode Discuss, Glassdoor, and
> community interview reports. Treat as a signal, not a guarantee.

---

## DSA Concepts Used

- **Bit Manipulation** — the whole problem is a statement about which bits survive an AND across a consecutive range; the answer is the *common binary prefix* of `left` and `right`, and Kernighan's `x & (x-1)` lowest-set-bit trick extracts it → see [`/dsa/bit_manipulation.md`](/dsa/bit_manipulation.md)

---

## Approaches Overview

| # | Approach | Time | Space | When to use |
|---|----------|------|-------|-------------|
| 1 | Brute Force | O(right − left) | O(1) | Baseline only; a range like [2³⁰, 2³¹−1] runs ~10⁹ iterations — TLE |
| 2 | Common Prefix (Bit Shifting) | O(log right) ≤ 31 steps | O(1) | The standard interview answer; easiest to prove |
| 3 | Brian Kernighan's Trick (Optimal) | O(popcount(right)) ≤ O(log right) | O(1) | Fewest iterations in practice; shows off the `x & (x−1)` idiom |

---

## Approach 1 — Brute Force

### Intuition

The problem statement *is* an algorithm: fold every number of the range into a running AND. The only intelligence added is an early exit — AND can only turn bits off, so a running result of `0` is final and the loop can stop. That early exit is what lets Example 3 (`right = 2^31 − 1`) finish instantly: `1 & 2 = 0` after a single step. It does **not** save adversarial ranges such as `[2³⁰, 2³¹−1]`, where bit 30 stays set for a billion iterations.

### Algorithm

1. Initialise `result = left`.
2. For `num` from `left + 1` up to `right`: set `result &= num`.
3. If `result` ever becomes `0`, return `0` immediately (it can never come back).
4. After the loop, return `result`.

### Complexity

- **Time:** O(right − left) — worst case touches every number in the range, up to ~2³¹ iterations; the early exit only helps ranges that straddle a power of two.
- **Space:** O(1) — one integer accumulator.

### Code

```go
func bruteForce(left int, right int) int {
	result := left // accumulate the AND starting from the low end
	for num := left + 1; num <= right; num++ {
		result &= num // fold the next number into the running AND
		if result == 0 {
			return 0 // bits only turn off — once zero, always zero
		}
	}
	return result
}
```

### Dry Run

Example 1: `left = 5, right = 7` (binary: 5 = `101`, 6 = `110`, 7 = `111`).

| Step | num | result before | AND operation | result after | result == 0? |
|------|-----|---------------|---------------|--------------|--------------|
| 0 | — | — | initialise | 5 (`101`) | no |
| 1 | 6 | 5 (`101`) | `101 & 110` | 4 (`100`) | no |
| 2 | 7 | 4 (`100`) | `100 & 111` | 4 (`100`) | no |

Loop ends (`num > 7`). Result: `4` ✔

---

## Approach 2 — Common Prefix (Bit Shifting)

### Intuition

Write `left` and `right` in binary and split each into "common prefix + differing tail". Key claim: **every bit below the common prefix ends up 0 in the AND.** Why: counting from `left` up to `right` must at some point flip the highest differing bit; on the way, each lower bit position cycles through both 0 and 1 (in particular, a carry ripples a 0 into it). So some number in the range has a 0 at every tail position, killing that bit in the AND. What survives is exactly the shared prefix. To find it, chop the lowest bit off both numbers until they become equal, then shift the surviving prefix back into position, padding the chopped positions with zeros.

### Algorithm

1. Set `shift = 0`.
2. While `left != right`: do `left >>= 1`, `right >>= 1`, `shift++`.
3. When they meet, `left` holds the common prefix; return `left << shift`.

### Complexity

- **Time:** O(log right) — each iteration discards one bit; at most 31 iterations for 32-bit inputs, independent of the range width.
- **Space:** O(1) — a shift counter and the two parameters mutated in place.

### Code

```go
func commonPrefixShift(left int, right int) int {
	shift := 0 // how many low bits we discarded
	// Keep chopping the lowest bit until the remaining prefixes agree.
	for left != right {
		left >>= 1  // drop lowest bit of left
		right >>= 1 // drop lowest bit of right
		shift++     // remember how far we shifted
	}
	// left == right == common prefix; move it back to its true position,
	// filling the discarded (differing) bits with zeros.
	return left << shift
}
```

### Dry Run

Example 1: `left = 5 (101), right = 7 (111)`.

| Step | left (bin) | right (bin) | left == right? | Action | shift after |
|------|------------|-------------|----------------|--------|-------------|
| 1 | `101` (5) | `111` (7) | no | shift both right by 1 | 1 |
| 2 | `10` (2)  | `11` (3)  | no | shift both right by 1 | 2 |
| 3 | `1` (1)   | `1` (1)   | yes | exit loop | 2 |

Reconstruct: `1 << 2 = 100₂ = 4`. Result: `4` ✔ — the common prefix of `101` and `111` is the single bit `1__`, padded with two zeros.

---

## Approach 3 — Brian Kernighan's Trick (Optimal)

### Intuition

Same prefix insight, attacked from the other side. `right & (right − 1)` clears the **lowest set bit** of `right` in one operation (Kernighan's trick). Every set bit of `right` that lies below the common prefix is doomed anyway (Approach 2's argument), so erase them: while `right > left`, keep clearing `right`'s lowest set bit. The loop stops at the first value of `right` that is ≤ `left`. That value is a prefix of the original `right` with trailing zeros — and since it is ≤ `left` but shares `left`'s high bits, it is exactly the common prefix, i.e. the answer. Each iteration removes a whole set bit rather than a single position, so it never runs more than `popcount(right)` times.

### Algorithm

1. While `left < right`: set `right &= right − 1` (erase the lowest set bit of `right`).
2. Return `right`.

### Complexity

- **Time:** O(popcount(right)) ≤ O(log right) — at most one iteration per set bit of `right` (≤ 31), usually far fewer because it stops as soon as `right ≤ left`.
- **Space:** O(1) — everything happens in the two parameters.

### Code

```go
func brianKernighan(left int, right int) int {
	// Erase right's lowest set bit until right sinks down to (or below) left.
	for left < right {
		right &= right - 1 // Kernighan: clears exactly the lowest set bit
	}
	return right // now the shared prefix of the original [left, right]
}
```

### Dry Run

Example 1: `left = 5 (101), right = 7 (111)`.

| Step | left | right (bin) | left < right? | Operation | right after |
|------|------|-------------|----------------|-----------|-------------|
| 1 | 5 | `111` (7) | yes | `111 & 110` | `110` (6) |
| 2 | 5 | `110` (6) | yes | `110 & 101` | `100` (4) |
| 3 | 5 | `100` (4) | no (4 < 5) | exit loop | `100` (4) |

Result: `4` ✔ — two bit-clears found the common prefix `100`.

---

## Key Takeaways

- **AND over a consecutive range = common binary prefix** of the endpoints, right-padded with zeros. Any bit position that differs between `left` and `right` — or sits below such a position — passes through 0 somewhere in the range.
- **AND accumulation is monotone:** bits only turn off. This gives both the early-exit in brute force and the guarantee that the prefix answer is final.
- **`x & (x − 1)` clears the lowest set bit** — the same idiom that powers LeetCode #191 (counting bits) and power-of-two checks (`x & (x-1) == 0`). Keep it in the toolbox.
- Range width is irrelevant to the optimal solutions: `[1, 2^31 − 1]` costs the same ~30 operations as `[5, 7]`. When a "range aggregate" collapses to a property of the endpoints, look for the invariant instead of iterating.

---

## Related Problems

- LeetCode #191 — Number of 1 Bits (Kernighan's lowest-set-bit loop)
- LeetCode #190 — Reverse Bits (bit-by-bit reconstruction)
- LeetCode #338 — Counting Bits (`dp[x] = dp[x & (x-1)] + 1`)
- LeetCode #371 — Sum of Two Integers (pure bit-manipulation arithmetic)
- LeetCode #2411 — Smallest Subarrays With Maximum Bitwise OR (range bit aggregates)
