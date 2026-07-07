# 0412 — Fizz Buzz

> LeetCode #412 · Difficulty: Easy
> **Categories:** Math, String, Simulation

---

## Problem Statement

Given an integer `n`, return a string array `answer` (**1-indexed**) where:

- `answer[i] == "FizzBuzz"` if `i` is divisible by `3` and `5`.
- `answer[i] == "Fizz"` if `i` is divisible by `3`.
- `answer[i] == "Buzz"` if `i` is divisible by `5`.
- `answer[i] == i` (as a string) if none of the above conditions are true.

**Example 1:**

```
Input: n = 3
Output: ["1","2","Fizz"]
```

**Example 2:**

```
Input: n = 5
Output: ["1","2","Fizz","4","Buzz"]
```

**Example 3:**

```
Input: n = 15
Output: ["1","2","Fizz","4","Buzz","Fizz","7","8","Fizz","Buzz","11","Fizz","13","14","FizzBuzz"]
```

**Constraints:**

- `1 <= n <= 10^4`

---

## Company Frequency

| Company    | Frequency        | Last Reported |
|------------|------------------|---------------|
| Amazon     | ★★★☆☆ Medium     | 2024          |
| Microsoft  | ★★★☆☆ Medium     | 2023          |
| Google     | ★★☆☆☆ Low        | 2023          |
| Apple      | ★★☆☆☆ Low        | 2023          |
| Bloomberg  | ★★☆☆☆ Low        | 2022          |
| Adobe      | ★☆☆☆☆ Rare       | 2022          |

> ⚠️ Frequency data is crowd-sourced from LeetCode Discuss, Glassdoor, and
> community interview reports. Treat as a signal, not a guarantee.

---

## DSA Concepts Used

- **Number Theory (divisibility & LCM)** — the whole task is a divisibility table; the "both 3 and 5" case is exactly "divisible by `lcm(3,5) = 15`", which is why testing `i%15` (or the two flags together) is correct → see [`/dsa/math_number_theory.md`](/dsa/math_number_theory.md)
- **String building** — each cell is a tiny string assembled from labels or the number's decimal text → see [`/dsa/string_algorithms.md`](/dsa/string_algorithms.md)

---

## Approaches Overview

| # | Approach | Time | Space | When to use |
|---|----------|------|-------|-------------|
| 1 | Brute Force (Modulo Checks) | O(n) | O(1) extra | The canonical answer; order the `%15` branch first |
| 2 | String Concatenation | O(n) | O(1) extra | Scales to extra divisors (7→"Bazz") with no combined-modulo explosion |
| 3 | Counter Increments (No Modulo) (Optimal) | O(n) | O(1) extra | When `%`/division is expensive or disallowed — additions only |

> All three are O(n) and optimal in the big-O sense (you must emit `n` cells). They differ in *which operations* they use, which is the real interview discussion.

---

## Approach 1 — Brute Force (Modulo Checks)

### Intuition

The statement is a rule table. Iterate every value `i` from `1` to `n` and choose the label by testing divisibility **in the correct order**. The `FizzBuzz` case (divisible by both 3 and 5, i.e. by 15) must be checked *before* the individual `Fizz`/`Buzz` cases; otherwise a multiple of 15 would match `i%3 == 0` first and be mislabelled "Fizz".

### Algorithm

1. Allocate `answer` of length `n` (index `i-1` holds the label for value `i`).
2. For `i` from `1` to `n`, pick the first matching case:
   1. `i % 15 == 0` → `"FizzBuzz"`.
   2. `i % 3 == 0` → `"Fizz"`.
   3. `i % 5 == 0` → `"Buzz"`.
   4. otherwise → `strconv.Itoa(i)`.
3. Return `answer`.

### Complexity

- **Time:** O(n) — one pass, constant work per element.
- **Space:** O(1) auxiliary beyond the required `n`-length output slice.

### Code

```go
func bruteForce(n int) []string {
	answer := make([]string, n) // 0-indexed slice; answer[i-1] holds the label for value i
	for i := 1; i <= n; i++ {   // problem is 1-indexed, so iterate values 1..n
		switch {
		case i%15 == 0: // divisible by 3 AND 5 → must be checked FIRST
			answer[i-1] = "FizzBuzz"
		case i%3 == 0: // divisible by 3 only
			answer[i-1] = "Fizz"
		case i%5 == 0: // divisible by 5 only
			answer[i-1] = "Buzz"
		default: // divisible by neither → the number as text
			answer[i-1] = strconv.Itoa(i)
		}
	}
	return answer
}
```

### Dry Run

Example 1: `n = 3`.

| Step | i | i%15 | i%3 | i%5 | Matched case | answer[i-1] |
|------|---|------|-----|-----|--------------|-------------|
| 1 | 1 | 1 | 1 | 1 | default | `"1"` |
| 2 | 2 | 2 | 2 | 2 | default | `"2"` |
| 3 | 3 | 3 | 0 | 3 | `i%3==0` | `"Fizz"` |

Result: `["1","2","Fizz"]` ✔

---

## Approach 2 — String Concatenation

### Intuition

Drop the special "divisible by 15" branch. Append `"Fizz"` whenever 3 divides `i` and `"Buzz"` whenever 5 divides `i`. A multiple of 15 passes **both** tests and naturally becomes `"Fizz"+"Buzz" = "FizzBuzz"`. If neither test fired, the accumulator is still empty, so substitute the number. The pattern extends to more rules (add `"Bazz"` for 7, etc.) without a combinatorial explosion of combined-modulo cases.

### Algorithm

1. For `i` from `1` to `n`, start with an empty builder `sb`.
2. If `i % 3 == 0`, write `"Fizz"`.
3. If `i % 5 == 0`, write `"Buzz"`.
4. If `sb` is still empty, write `strconv.Itoa(i)`.
5. Store `sb.String()`.

### Complexity

- **Time:** O(n) — constant, bounded string building per element.
- **Space:** O(1) auxiliary beyond the output (each builder holds ≤ 8 bytes).

### Code

```go
func stringConcat(n int) []string {
	answer := make([]string, n)
	for i := 1; i <= n; i++ {
		var sb strings.Builder // accumulate the label without a combined 15-check
		if i%3 == 0 {
			sb.WriteString("Fizz") // 3 contributes "Fizz"
		}
		if i%5 == 0 {
			sb.WriteString("Buzz") // 5 contributes "Buzz"; both → "FizzBuzz"
		}
		if sb.Len() == 0 { // no divisor matched → use the number itself
			sb.WriteString(strconv.Itoa(i))
		}
		answer[i-1] = sb.String()
	}
	return answer
}
```

### Dry Run

Example 1: `n = 3`.

| Step | i | i%3==0? | i%5==0? | sb after 3-check | sb after 5-check | empty? → number | final |
|------|---|---------|---------|------------------|------------------|-----------------|-------|
| 1 | 1 | no | no | `""` | `""` | yes → `"1"` | `"1"` |
| 2 | 2 | no | no | `""` | `""` | yes → `"2"` | `"2"` |
| 3 | 3 | yes | no | `"Fizz"` | `"Fizz"` | no | `"Fizz"` |

Result: `["1","2","Fizz"]` ✔

---

## Approach 3 — Counter Increments (No Modulo) (Optimal)

### Intuition

Divisibility by 3 recurs every 3 steps and by 5 every 5 steps. Keep two counters `fizz` and `buzz` that increment each iteration; when `fizz` reaches 3 we are on a multiple of 3 (reset it to 0), and when `buzz` reaches 5 we are on a multiple of 5 (reset it to 0). Combine the two "just reset?" flags exactly like Approach 2. This uses **no modulo or division** — only additions and comparisons — which matters on hardware where `%` is costly or in interview variants that forbid it.

### Algorithm

1. Initialise `fizz = 0`, `buzz = 0`.
2. For `i` from `1` to `n`: `fizz++`, `buzz++`.
3. `isFizz = (fizz == 3)`; if so reset `fizz = 0`. `isBuzz = (buzz == 5)`; if so reset `buzz = 0`.
4. Emit `"FizzBuzz"` / `"Fizz"` / `"Buzz"` / `strconv.Itoa(i)` from the two flags.

### Complexity

- **Time:** O(n) — constant work per element; only `++` and `==`.
- **Space:** O(1) auxiliary beyond the output (two integer counters).

### Code

```go
func counterNoModulo(n int) []string {
	answer := make([]string, n)
	fizz, buzz := 0, 0 // steps since the last multiple of 3 / of 5
	for i := 1; i <= n; i++ {
		fizz++ // advance both cyclic counters
		buzz++
		isFizz := fizz == 3 // reached a multiple of 3 this step?
		isBuzz := buzz == 5 // reached a multiple of 5 this step?
		if isFizz {
			fizz = 0 // restart the 3-cycle
		}
		if isBuzz {
			buzz = 0 // restart the 5-cycle
		}
		switch {
		case isFizz && isBuzz: // both cycles landed → multiple of 15
			answer[i-1] = "FizzBuzz"
		case isFizz:
			answer[i-1] = "Fizz"
		case isBuzz:
			answer[i-1] = "Buzz"
		default:
			answer[i-1] = strconv.Itoa(i)
		}
	}
	return answer
}
```

### Dry Run

Example 1: `n = 3`.

| Step | i | fizz→ | buzz→ | isFizz | isBuzz | reset fizz? | reset buzz? | answer[i-1] |
|------|---|-------|-------|--------|--------|-------------|-------------|-------------|
| 1 | 1 | 1 | 1 | no | no | — | — | `"1"` |
| 2 | 2 | 2 | 2 | no | no | — | — | `"2"` |
| 3 | 3 | 3 | 3 | **yes** | no | fizz=0 | — | `"Fizz"` |

Result: `["1","2","Fizz"]` ✔

---

## Key Takeaways

- **Order the most-specific case first.** Testing `%15` (or "both flags") before `%3`/`%5` is the whole correctness trick; the classic bug is checking `Fizz` first and never reaching `FizzBuzz`.
- **"Divisible by both a and b" = "divisible by lcm(a,b)".** Here `lcm(3,5) = 15` because 3 and 5 are coprime.
- **Flag-accumulation generalises.** The concat approach adds new rules without a combinatorial blowup of combined-modulo branches — the go-to design when the rule set may grow.
- **Modulo is not mandatory.** Cyclic counters replace `%` with additions, a handy pattern when division is expensive or banned.
- All solutions are Θ(n): you must produce `n` outputs, so no algorithm beats linear here.

---

## Related Problems

- LeetCode #1195 — Fizz Buzz Multithreaded (same rules, concurrency twist)
- LeetCode #1290 — Convert Binary Number in a Linked List to Integer (digit/number formatting)
- LeetCode #66 — Plus One (per-digit simulation)
- LeetCode #7 — Reverse Integer (integer ↔ string handling)
