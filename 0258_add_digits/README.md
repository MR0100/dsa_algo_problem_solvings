# 0258 — Add Digits

> LeetCode #258 · Difficulty: Easy
> **Categories:** Math, Simulation, Number Theory

---

## Problem Statement

Given an integer `num`, repeatedly add all its digits until the result has only one digit, and return it.

**Example 1:**

```
Input: num = 38
Output: 2
Explanation: The process is
38 --> 3 + 8 --> 11
11 --> 1 + 1 --> 2
Since 2 has only one digit, return it.
```

**Example 2:**

```
Input: num = 0
Output: 0
```

**Constraints:**

- `0 <= num <= 2^31 - 1`

**Follow up:** Could you do it without any loop/recursion in `O(1)` runtime?

---

## Company Frequency

| Company    | Frequency        | Last Reported |
|------------|------------------|---------------|
| Amazon     | ★★★☆☆ Medium     | 2023          |
| Microsoft  | ★★★☆☆ Medium     | 2023          |
| Google     | ★★☆☆☆ Low        | 2023          |
| Apple      | ★★☆☆☆ Low        | 2022          |
| Adobe      | ★★☆☆☆ Low        | 2022          |

> ⚠️ Frequency data is crowd-sourced from LeetCode Discuss, Glassdoor, and
> community interview reports. Treat as a signal, not a guarantee.

---

## DSA Concepts Used

- **Math / Number Theory (digital root)** — repeatedly summing digits preserves the value mod 9, so the answer is the digital root `1 + (num-1) mod 9`; this also answers the O(1) follow-up → see [`/dsa/math_number_theory.md`](/dsa/math_number_theory.md)

---

## Approaches Overview

| # | Approach | Time | Space | When to use |
|---|----------|------|-------|-------------|
| 1 | Iterative Digit Sum | O(log num) | O(1) | Direct simulation of the definition |
| 2 | Single-Pass Recursion | O(log num) | O(log* num) | Same logic expressed recursively |
| 3 | Digital Root Formula (Optimal) | O(1) | O(1) | The follow-up answer: no loop/recursion |

---

## Approach 1 — Iterative Digit Sum

### Intuition

The problem is stated as a loop: replace `num` with the sum of its digits, and repeat while the result still has more than one digit. Implement exactly that.

### Algorithm

1. While `num >= 10` (more than one digit):
   1. `sum = 0`; while `num > 0`: `sum += num % 10`, `num /= 10`.
   2. `num = sum`.
2. Return `num`.

### Complexity

- **Time:** O(log num) — each pass costs O(digits); the value collapses toward its digit sum extremely fast, so only a couple of passes run.
- **Space:** O(1) — two scalars.

### Code

```go
func iterativeDigitSum(num int) int {
	for num >= 10 { // keep going while more than one digit
		sum := 0
		for num > 0 { // add up the decimal digits
			sum += num % 10 // take the lowest digit
			num /= 10       // drop it
		}
		num = sum // replace num with its digit sum, repeat
	}
	return num
}
```

### Dry Run

Example 1: `num = 38`.

| Outer pass | inner sum computation | num after |
|------------|-----------------------|-----------|
| 1 | 8 + 3 = 11 | 11 |
| 2 | 1 + 1 = 2 | 2 |
| — | 2 < 10 → exit | 2 |

Result: `2` ✔

---

## Approach 2 — Single-Pass Recursion

### Intuition

The same repeated digit-sum, expressed recursively. If `num` is already a single digit it is the answer; otherwise compute its digit sum once and recurse on that.

### Algorithm

1. If `num < 10`, return `num` (base case).
2. Compute the digit sum `s` of `num`.
3. Return `recursiveDigitSum(s)`.

### Complexity

- **Time:** O(log num) amortized — same shrinking behavior as Approach 1.
- **Space:** O(log* num) recursion depth — effectively constant.

### Code

```go
func recursiveDigitSum(num int) int {
	if num < 10 { // base case: already a single digit
		return num
	}
	sum := 0
	for n := num; n > 0; n /= 10 { // sum the digits once
		sum += n % 10
	}
	return recursiveDigitSum(sum) // recurse on the reduced value
}
```

### Dry Run

Example 1: `num = 38`.

| Call | num | digit sum | recurse into |
|------|-----|-----------|--------------|
| recursiveDigitSum(38) | 38 | 3+8 = 11 | recursiveDigitSum(11) |
| recursiveDigitSum(11) | 11 | 1+1 = 2 | recursiveDigitSum(2) |
| recursiveDigitSum(2) | 2 | — (< 10) | return 2 |

Result: `2` ✔

---

## Approach 3 — Digital Root O(1) Formula (Optimal)

### Intuition

Because `10 ≡ 1 (mod 9)`, every power of ten is congruent to 1 mod 9, so a number is congruent to the sum of its digits mod 9. Repeated digit-summing therefore never changes the value mod 9 and converges to the **digital root**. The clean closed form is: `0` when `num == 0`, otherwise `1 + (num-1) mod 9`. The `(num-1)` shift maps multiples of 9 to `9` instead of `0`.

### Algorithm

1. If `num == 0`, return `0`.
2. Return `1 + (num-1) % 9`.

### Complexity

- **Time:** O(1) — pure arithmetic, satisfies the follow-up (no loop/recursion).
- **Space:** O(1).

### Code

```go
func digitalRoot(num int) int {
	if num == 0 {
		return 0 // 0 stays 0 (the formula below would also give 0, guarded for clarity)
	}
	// 1 + (num-1)%9 collapses to 9 for multiples of 9, else num%9.
	return 1 + (num-1)%9
}
```

### Dry Run

Example 1: `num = 38`.

| Step | computation | value |
|------|-------------|-------|
| 1 | num != 0 | — |
| 2 | (num - 1) = 37 | 37 |
| 3 | 37 % 9 = 1 | 1 |
| 4 | 1 + 1 = 2 | 2 |

Result: `2` ✔. Check `num = 9`: `1 + (8 % 9) = 1 + 8 = 9` (multiple of 9 → 9, not 0).

---

## Key Takeaways

- **Digital root = `1 + (n-1) mod 9`** for `n > 0`, and `0` for `n == 0`. Memorize this — it turns a loop into O(1).
- **Why mod 9 works:** `10^k ≡ 1 (mod 9)`, so a number ≡ its digit sum (mod 9). The same fact powers the "casting out nines" checksum.
- **The `-1 … +1` shift** is the standard trick to fold a `0` residue up to the top of a 1-based range (here, mapping multiples of 9 to 9).
- Simulation is fine when values shrink fast; but recognizing the invariant (value mod 9) is what unlocks the constant-time answer.

---

## Related Problems

- LeetCode #202 — Happy Number (repeated digit transformation, cycle detection)
- LeetCode #1837 — Sum of Digits in Base K (digit-sum in another base)
- LeetCode #1945 — Sum of Digits After Convert (iterated digit sums)
- LeetCode #2544 — Alternating Digit Sum (digit-by-digit processing)
