# 0233 — Number of Digit One

> LeetCode #233 · Difficulty: Hard
> **Categories:** Math, Dynamic Programming, Recursion

---

## Problem Statement

Given an integer `n`, count the total number of digit `1` appearing in all non-negative integers less than or equal to `n`.

**Example 1:**
```
Input: n = 13
Output: 6
Explanation: The digit 1 occurs in the following numbers: 1, 10, 11, 12, 13.
Counting each occurrence: 1 → 1, 10 → 1, 11 → 2, 12 → 1, 13 → 1. Total = 6.
```

**Example 2:**
```
Input: n = 0
Output: 0
```

**Constraints:**
- `0 <= n <= 10⁹`

---

## Company Frequency

| Company    | Frequency        | Last Reported |
|------------|------------------|---------------|
| Google     | ★★★☆☆ Medium     | 2023          |
| Amazon     | ★★☆☆☆ Low        | 2023          |
| Microsoft  | ★★☆☆☆ Low        | 2022          |
| Adobe      | ★★☆☆☆ Low        | 2022          |

> ⚠️ Frequency data is crowd-sourced from LeetCode Discuss, Glassdoor, and
> community interview reports. Treat as a signal, not a guarantee.

---

## DSA Concepts Used
- **Number Theory / Digit Counting** — count contributions per decimal place rather than per number → see [`/dsa/math_number_theory.md`](/dsa/math_number_theory.md)
- **Digit DP (combinatorial counting)** — the place-value split is the closed-form heart of digit dynamic programming → see [`/dsa/dynamic_programming_1d.md`](/dsa/dynamic_programming_1d.md)

---

## Approaches Overview
| # | Approach | Time | Space | When to use |
|---|----------|------|-------|-------------|
| 1 | Brute Force (count per number) | O(n log n) | O(1) | Small n, or as a correctness oracle |
| 2 | Digit-Position Counting (Optimal) | O(log n) | O(1) | Any n up to 10⁹; the intended solution |

---

## Approach 1 — Brute Force (Count Per Number)

### Intuition
The definition is literal: sum, over all integers `i` in `[1, n]`, the number of times digit `1` appears in `i`. So just do that — for each `i`, strip its decimal digits and tally the ones. Simple and obviously correct, useful as a reference oracle, but O(n log n) — too slow for `n` near 10⁹.

### Algorithm
1. For each `i` from 1 to `n`:
   1. While `i` has digits left, look at the lowest digit.
   2. If it is 1, increment the total.
   3. Drop the lowest digit (`i /= 10`).
2. Return the accumulated total.

### Complexity
- **Time:** O(n log n) — n numbers, each with O(log n) digits.
- **Space:** O(1).

### Code
```go
func bruteForce(n int) int {
	total := 0
	for i := 1; i <= n; i++ { // consider every integer in [1, n]
		x := i
		for x > 0 { // examine each decimal digit of x
			if x%10 == 1 { // lowest digit is a 1 → count it
				total++
			}
			x /= 10 // discard the lowest digit
		}
	}
	return total
}
```

### Dry Run
Trace `n = 13`:

| i  | digits scanned | ones in i | running total |
|----|----------------|-----------|---------------|
| 1  | 1              | 1         | 1             |
| 2–9| (no 1s)        | 0         | 1             |
| 10 | 0, 1           | 1         | 2             |
| 11 | 1, 1           | 2         | 4             |
| 12 | 2, 1           | 1         | 5             |
| 13 | 3, 1           | 1         | 6             |

Return `6`.

---

## Approach 2 — Digit-Position Counting (Optimal)

### Intuition
Instead of counting per number, count **per digit position**. Fix a place value `p` (1, 10, 100, …). Split `n` around that place into `high`, `cur`, and `low`. How many numbers in `[1, n]` carry a 1 at place `p` depends only on `cur`:

- `cur == 0` : `high * p` — only the completed high cycles contribute.
- `cur == 1` : `high * p + (low + 1)` — the in-progress cycle contributes ones for `low+1` numbers.
- `cur >= 2` : `(high + 1) * p` — the place has fully passed 1 once more.

Summing over all places gives the answer in O(log n).

### Algorithm
1. For place `p = 1, 10, 100, …` while `p <= n`:
   - `high = n / (p*10)`, `cur = (n / p) % 10`, `low = n % p`.
   - if `cur == 0`: add `high * p`.
   - if `cur == 1`: add `high * p + low + 1`.
   - if `cur >= 2`: add `(high + 1) * p`.
2. Return the running sum.

### Complexity
- **Time:** O(log n) — one iteration per decimal place of `n`.
- **Space:** O(1).

### Code
```go
func digitPosition(n int) int {
	count := 0
	for p := 1; p <= n; p *= 10 { // iterate over each place value 1,10,100,...
		high := n / (p * 10) // digits above the current place
		cur := (n / p) % 10  // the digit sitting at the current place
		low := n % p         // digits below the current place

		switch {
		case cur == 0:
			count += high * p
		case cur == 1:
			count += high*p + low + 1
		default: // cur >= 2
			count += (high + 1) * p
		}
	}
	return count
}
```

### Dry Run
Trace `n = 13`:

| p  | high = n/(p*10) | cur = (n/p)%10 | low = n%p | case      | contribution | count |
|----|-----------------|----------------|-----------|-----------|--------------|-------|
| 1  | 13/10 = 1       | 13%10 = 3      | 0         | cur ≥ 2   | (1+1)*1 = 2  | 2     |
| 10 | 13/100 = 0      | (13/10)%10 = 1 | 13%10 = 3 | cur == 1  | 0*10 + 3 + 1 = 4 | 6 |

`p = 100 > 13`, loop ends. Return `6`.

- The `p = 1` step counts ones in the **units** place: numbers 1 and 11 → 2 occurrences.
- The `p = 10` step counts ones in the **tens** place: numbers 10, 11, 12, 13 → 4 occurrences.
- Total `2 + 4 = 6`. ✓

---

## Key Takeaways
- **Count by position, not by value.** For "how many times does digit d appear up to n", decompose each place into high/cur/low and use a closed form — this is the standard digit-DP identity.
- The three-way split on `cur` (`0` / `1` / `≥2`) captures whether the partial cycle has reached, is at, or has passed the digit 1.
- Keep a brute-force oracle around: it makes it trivial to fuzz-test the O(log n) formula for off-by-one errors on the `low + 1` term.
- The same template generalises to counting any fixed digit (with a small tweak for digit 0, which can't lead).

## Related Problems
- LeetCode #172 — Factorial Trailing Zeroes (place-value / factor counting)
- LeetCode #357 — Count Numbers with Unique Digits (combinatorial digit counting)
- LeetCode #902 — Numbers At Most N Given Digit Set (digit DP)
- LeetCode #1067 — Digit Count in Range (direct generalisation)
