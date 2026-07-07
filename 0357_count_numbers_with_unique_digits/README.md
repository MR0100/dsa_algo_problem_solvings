# 0357 — Count Numbers with Unique Digits

> LeetCode #357 · Difficulty: Medium
> **Categories:** Math, Dynamic Programming, Backtracking, Combinatorics

---

## Problem Statement

Given an integer `n`, return the count of all numbers with unique digits, `x`, where `0 <= x < 10^n`.

**Example 1:**

```
Input: n = 2
Output: 91
Explanation: The answer should be the total numbers in the range of 0 ≤ x < 100,
excluding 11, 22, 33, 44, 55, 66, 77, 88, 99
```

**Example 2:**

```
Input: n = 0
Output: 1
```

**Constraints:**

- `0 <= n <= 8`

---

## Company Frequency

| Company    | Frequency        | Last Reported |
|------------|------------------|---------------|
| Google     | ★★★☆☆ Medium     | 2023          |
| Amazon     | ★★☆☆☆ Low        | 2022          |
| Microsoft  | ★★☆☆☆ Low        | 2022          |
| Bloomberg  | ★★☆☆☆ Low        | 2021          |

> ⚠️ Frequency data is crowd-sourced from LeetCode Discuss, Glassdoor, and
> community interview reports. Treat as a signal, not a guarantee.

---

## DSA Concepts Used

- **Combinatorics / Counting** — the answer is a product of shrinking digit choices (9 · 9 · 8 · …), counted per digit length → see [`/dsa/math_number_theory.md`](/dsa/math_number_theory.md)
- **Dynamic Programming (build-up)** — each length's count derives from the previous by one multiplication → see [`/dsa/dynamic_programming_1d.md`](/dsa/dynamic_programming_1d.md)
- **Backtracking** — DFS constructs every unique-digit number with a `used[]` mask, counting valid prefixes → see [`/dsa/backtracking.md`](/dsa/backtracking.md)

---

## Approaches Overview

| # | Approach | Time | Space | When to use |
|---|----------|------|-------|-------------|
| 1 | Brute Force | O(10ⁿ · n) | O(1) | Baseline / verification only |
| 2 | Combinatorics (Optimal) | O(min(n, 10)) | O(1) | Best; closed-form counting |
| 3 | Backtracking (DFS) | ≤ O(10!) bounded | O(n) | When you must also *generate* the numbers |

---

## Approach 1 — Brute Force

### Intuition

The problem is a literal count of `x` in `[0, 10^n)` with distinct digits, so enumerate every `x` and check it. For each `x`, walk its decimal digits marking a 10-slot `seen` array; the first repeat disqualifies it. `0` counts as it is a single digit. Correct but exponential — fine as a reference oracle.

### Algorithm

1. Compute `limit = 10^n`.
2. For each `x` in `[0, limit)`, extract digits and test uniqueness with a `seen[10]` mask.
3. Increment a counter for each unique `x`.

### Complexity

- **Time:** O(10ⁿ · n) — enumerates the whole range, digit-checking each.
- **Space:** O(1) — fixed 10-slot array.

### Code

```go
func bruteForce(n int) int {
	limit := 1 // will become 10^n
	for i := 0; i < n; i++ {
		limit *= 10
	}
	count := 0
	for x := 0; x < limit; x++ {
		if hasUniqueDigits(x) {
			count++ // x qualifies
		}
	}
	return count
}

func hasUniqueDigits(x int) bool {
	var seen [10]bool // seen[d] = have we already used digit d?
	if x == 0 {
		return true // "0" is a single unique digit
	}
	for x > 0 {
		d := x % 10 // lowest digit
		if seen[d] {
			return false // repeat found
		}
		seen[d] = true
		x /= 10 // drop the digit we just consumed
	}
	return true
}
```

### Dry Run

Example 1: `n = 2`, `limit = 100`. (Showing the disqualifications.)

| x | digits | unique? |
|---|--------|---------|
| 0..10 | e.g. 0,1,…,10 | all unique |
| 11 | 1,1 | ✗ excluded |
| 12..21 | unique | ✓ |
| 22 | 2,2 | ✗ excluded |
| … | … | … |
| 99 | 9,9 | ✗ excluded |

Excluded exactly `{11,22,…,99}` = 9 numbers ⇒ `100 - 9 = 91`. Result: **91** ✔

---

## Approach 2 — Combinatorics (Optimal)

### Intuition

Count by number of digits `k` and sum. One-digit numbers `0..9` give 10. For `k`-digit numbers (`k ≥ 2`): the leading digit has **9** choices (`1..9`, no leading zero), the next has **9** choices (`0..9` minus the one used), then **8**, **7**, …. So the count of `k`-digit unique numbers is `9 · 9 · 8 · … · (10 - k + 1)`. Sum over `k = 1..n`. Since only 10 distinct digits exist, any `n > 10` adds nothing beyond `k = 10`.

### Algorithm

1. If `n == 0`, the range is `[0, 1)` ⇒ only `0` ⇒ return `1`.
2. Start `total = 10` (all one-digit numbers).
3. `uniqueDigits = 9`, `available = 9`. For each length `k = 2..n`: `uniqueDigits *= available`, `total += uniqueDigits`, `available--`.
4. Return `total`.

### Complexity

- **Time:** O(min(n, 10)) — at most ~10 multiplications.
- **Space:** O(1).

### Code

```go
func combinatorics(n int) int {
	if n == 0 {
		return 1 // only 0 lies in [0, 1)
	}
	total := 10             // all one-digit numbers 0..9
	uniqueDigits := 9       // count of k-digit unique numbers, starts at k=1 (9)
	availableDigits := 9    // choices for the next position (0..9 minus used)
	for k := 2; k <= n && availableDigits > 0; k++ {
		uniqueDigits *= availableDigits // extend numbers by one more distinct digit
		total += uniqueDigits           // add all k-digit unique numbers
		availableDigits--               // one fewer digit remains for the next slot
	}
	return total
}
```

### Dry Run

Example 1: `n = 2`.

| Step | k | uniqueDigits | available | total |
|------|---|--------------|-----------|-------|
| init | — | 9 | 9 | 10 |
| loop | 2 | 9 · 9 = 81 | 9 → 8 | 10 + 81 = 91 |

Loop ends (`k > n`). Result: **91** ✔

---

## Approach 3 — Backtracking (DFS)

### Intuition

Every unique-digit number is a root-to-node path in a tree: choose a leading digit `1..9`, then any unused digit for each further position, up to `n` digits. Each node of length `1..n` is exactly one number. DFS with a `used[]` mask, counting every node visited. This constructs the numbers explicitly — handy for the "generate them all" variant — and its total matches the combinatorial count.

### Algorithm

1. Start `count = 1` to pre-count the number `0`.
2. For each leading digit `1..9`: mark it used, count it, DFS to extend.
3. In DFS at depth `d < n`: for each unused digit `0..9`, mark it, count the longer prefix, recurse, then unmark (backtrack).
4. Return `count`.

### Complexity

- **Time:** bounded by the number of unique-digit numbers (≈ 8.9M for n = 10) ≤ O(10!).
- **Space:** O(n) recursion depth + O(10) mask.

### Code

```go
func backtracking(n int) int {
	if n == 0 {
		return 1
	}
	var used [10]bool // digits already placed on the current path
	count := 1        // pre-count the number 0

	var dfs func(depth int)
	dfs = func(depth int) {
		if depth == n {
			return // cannot append more digits
		}
		for d := 0; d <= 9; d++ {
			if used[d] {
				continue // digit already used on this path
			}
			used[d] = true // place digit d
			count++        // this longer prefix is itself a valid number
			dfs(depth + 1) // extend further
			used[d] = false // backtrack: free digit d
		}
	}

	for d := 1; d <= 9; d++ {
		used[d] = true
		count++        // the 1-digit number d
		dfs(1)         // extend to 2..n digits
		used[d] = false
	}
	return count
}
```

### Dry Run

Example 1: `n = 2`. Start `count = 1` (the number 0).

| Phase | Detail | count |
|-------|--------|-------|
| leading digits 1..9 | each counted as a 1-digit number (+9) | 1 + 9 = 10 |
| extend each (depth 1 < 2) | for leading `d`, append any of the 9 remaining digits (+9 each, ×9 leaders) | 10 + 81 = 91 |

DFS at depth 2 returns immediately (`depth == n`). Result: **91** ✔

---

## Key Takeaways

- **Count by structure, not by enumeration.** Bucketing by digit length turns an exponential scan into a tiny product `9 · 9 · 8 · 7 · …`.
- **"No leading zero" costs one choice**: the first slot has 9 options (`1..9`), every later slot also has 9 initially (`0..9` minus one used), then decreasing.
- **Pigeonhole cap**: with only 10 distinct digits, no unique-digit number exceeds 10 digits, so the sum saturates at `n = 10` — larger `n` adds nothing.
- Backtracking and combinatorics agree because DFS visits exactly one node per valid number; use DFS only when you must emit the numbers, else the closed form is O(1)-ish.

---

## Related Problems

- LeetCode #46 — Permutations (backtracking with a used mask)
- LeetCode #1012 — Numbers With Repeated Digits (complement counting, digit DP)
- LeetCode #233 — Number of Digit One (digit-position combinatorics)
- LeetCode #902 — Numbers At Most N Given Digit Set (digit DP counting)
