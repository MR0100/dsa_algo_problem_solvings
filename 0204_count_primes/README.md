# 0204 — Count Primes

> LeetCode #204 · Difficulty: Medium
> **Categories:** Math, Array, Enumeration, Number Theory

---

## Problem Statement

Given an integer `n`, return *the number of prime numbers that are strictly less than* `n`.

**Example 1:**

```
Input: n = 10
Output: 4
Explanation: There are 4 prime numbers less than 10, they are 2, 3, 5, 7.
```

**Example 2:**

```
Input: n = 0
Output: 0
```

**Example 3:**

```
Input: n = 1
Output: 0
```

**Constraints:**

- `0 <= n <= 5 * 10^6`

---

## Company Frequency

| Company       | Frequency        | Last Reported |
|---------------|------------------|---------------|
| Amazon        | ★★★★☆ High       | 2024          |
| Microsoft     | ★★★★☆ High       | 2024          |
| Google        | ★★★☆☆ Medium     | 2024          |
| Apple         | ★★★☆☆ Medium     | 2023          |
| Goldman Sachs | ★★★☆☆ Medium     | 2023          |
| Capital One   | ★★☆☆☆ Low        | 2023          |

> ⚠️ Frequency data is crowd-sourced from LeetCode Discuss, Glassdoor, and
> community interview reports. Treat as a signal, not a guarantee.

---

## DSA Concepts Used

- **Math / Number Theory — primes & trial division** — a number is prime iff no divisor ≤ its square root exists; divisors pair up around √x → see [`/dsa/math_number_theory.md`](/dsa/math_number_theory.md)
- **Math / Number Theory — Sieve of Eratosthenes** — invert the work: primes cross out their multiples in bulk, O(n log log n) for all primes below n at once → see [`/dsa/math_number_theory.md`](/dsa/math_number_theory.md)

---

## Approaches Overview

| # | Approach | Time | Space | When to use |
|---|----------|------|-------|-------------|
| 1 | Brute Force (Trial Division) | O(n·√n) | O(1) | Only for a handful of queries or tiny n; too slow at 5·10⁶ scale |
| 2 | Sieve of Eratosthenes | O(n log log n) | O(n) | The expected interview answer; simple and near-linear |
| 3 | Odd-Only Sieve | O(n log log n) | O(n/2) | Same asymptotics, ~half the memory/work — a strong constant-factor upgrade |
| 4 | Linear (Euler's) Sieve (Optimal) | O(n) | O(n) | True linear time; also yields smallest-prime-factor info for free |

---

## Approach 1 — Brute Force (Trial Division)

### Intuition

Apply the definition to every candidate. `x` is prime iff no `d ∈ [2, √x]` divides it — checking beyond √x is pointless because divisors come in pairs `(d, x/d)` and the smaller member of every pair is ≤ √x. Counting primes below `n` then costs one trial-division test per candidate. Composites usually fail fast (half of everything falls to `d = 2`), but each *prime* pays the full √x loop, and with ~348 513 primes below 5·10⁶ this multiplies out to hundreds of millions of divisions — the reason this approach exists only as a baseline.

### Algorithm

1. Set `count = 0`.
2. For every `x` from 2 to `n − 1`:
   1. Test primality: for `d = 2` while `d·d ≤ x`, if `x % d == 0` then `x` is composite — stop.
   2. If no divisor was found, increment `count`.
3. Return `count`.

### Complexity

- **Time:** O(n·√n) — n candidates × up to √n trial divisions each (primes pay the full price; the constraint max n = 5·10⁶ makes this impractical).
- **Space:** O(1) — no arrays, just loop variables.

### Code

```go
func bruteForce(n int) int {
	count := 0
	for x := 2; x < n; x++ { // "strictly less than n" per the statement
		if isPrime(x) {
			count++
		}
	}
	return count
}

// isPrime reports whether x is prime by trial division up to √x.
func isPrime(x int) bool {
	if x < 2 {
		return false // 0 and 1 are not prime by definition
	}
	for d := 2; d*d <= x; d++ { // divisors pair up around √x
		if x%d == 0 {
			return false // found a nontrivial divisor → composite
		}
	}
	return true
}
```

### Dry Run

Example 1: `n = 10`.

| Step | x | Trial divisions (d with d² ≤ x) | Verdict | count after |
|------|---|--------------------------------|---------|-------------|
| 1 | 2 | none (2² > 2) | prime | 1 |
| 2 | 3 | none (2² > 3) | prime | 2 |
| 3 | 4 | 4 % 2 == 0 | composite | 2 |
| 4 | 5 | 5 % 2 ≠ 0 (3² > 5, stop) | prime | 3 |
| 5 | 6 | 6 % 2 == 0 | composite | 3 |
| 6 | 7 | 7 % 2 ≠ 0 (3² > 7, stop) | prime | 4 |
| 7 | 8 | 8 % 2 == 0 | composite | 4 |
| 8 | 9 | 9 % 2 ≠ 0, 9 % 3 == 0 | composite | 4 |

Result: `4` ✔ (primes 2, 3, 5, 7)

---

## Approach 2 — Sieve of Eratosthenes

### Intuition

Trial division asks, for each number, "who divides me?" — and mostly gets "nobody". The sieve flips the direction of work: each discovered prime `p` *announces* its multiples as composite, in bulk, with pure addition. Two classic refinements make it fast: (1) start marking at `p²`, because any smaller multiple `k·p` with `k < p` has a prime factor smaller than `p` and was already crossed out by it; (2) stop seeding once `p² ≥ n`, because such a prime has nothing left to announce. Whatever survives all crossings is prime.

### Algorithm

1. If `n < 3`, return 0 (no primes strictly below 0, 1, or 2).
2. Allocate `isComposite[0..n-1]`, all `false`.
3. For `i = 2` while `i·i < n`:
   1. If `isComposite[i]`, skip — a composite seed's prime factors already ran.
   2. Otherwise `i` is prime: mark `i², i²+i, i²+2i, … < n` as composite.
4. Count indices `x ∈ [2, n)` with `isComposite[x] == false`; return the count.

### Complexity

- **Time:** O(n log log n) — total marking work is `n·Σ(1/p)` over primes `p < √n`, and the sum of reciprocal primes grows as log log n (Mertens' theorem).
- **Space:** O(n) — one boolean per number below n (5 MB at the constraint max).

### Code

```go
func sieveOfEratosthenes(n int) int {
	if n < 3 {
		return 0 // no primes strictly below 0, 1, or 2
	}
	isComposite := make([]bool, n) // index = number; false means "maybe prime"
	for i := 2; i*i < n; i++ {     // only seeds up to √n can start new crossings
		if isComposite[i] {
			continue // composite seeds add nothing new — their primes already ran
		}
		for multiple := i * i; multiple < n; multiple += i {
			isComposite[multiple] = true // i is prime; kill its multiples from i²
		}
	}
	count := 0
	for x := 2; x < n; x++ {
		if !isComposite[x] {
			count++ // survivor of all crossings → prime
		}
	}
	return count
}
```

### Dry Run

Example 1: `n = 10`. Seeds run while `i² < 10`, i.e. `i ∈ {2, 3}`.

| Step | i | i prime? | Marks (from i², step i) | isComposite view (2..9) |
|------|---|----------|--------------------------|--------------------------|
| 1 | 2 | yes (unmarked) | 4, 6, 8 | 2 3 ~~4~~ 5 ~~6~~ 7 ~~8~~ 9 |
| 2 | 3 | yes (unmarked) | 9 | 2 3 ~~4~~ 5 ~~6~~ 7 ~~8~~ ~~9~~ |
| 3 | 4 | — | loop ends (4² = 16 ≥ 10) | — |

Counting pass: unmarked = {2, 3, 5, 7} → `4` ✔

---

## Approach 3 — Odd-Only Sieve

### Intuition

Half of the classic sieve's table — every even index — is dead weight: 2 is the *only* even prime. So count 2 once, then store only odd numbers, mapping odd `x` to index `x/2` (3→1, 5→2, 7→3, …). Marking also stays odd-only: an odd prime `p` crosses out `p², p²+2p, p²+4p, …` — stepping by `2p` skips the even multiples that the table doesn't even represent. Same O(n log log n) asymptotics, but half the memory, half the marking work, and better cache behaviour — the standard competitive-programming upgrade.

### Algorithm

1. If `n < 3`, return 0.
2. Allocate `isComposite` of length `n/2`; index `i` represents the odd number `2i+1`.
3. For `i = 1` while `(2i+1)² < n`:
   1. If `isComposite[i]`, skip.
   2. Else `p = 2i+1` is prime: for `m = p²; m < n; m += 2p`, set `isComposite[m/2] = true`.
4. Return `1` (the prime 2) + the number of unmarked indices `i ≥ 1`.

### Complexity

- **Time:** O(n log log n) — identical sum as Approach 2 restricted to odd multiples, roughly halving the constant factor.
- **Space:** O(n/2) — one boolean per odd number below n (2.5 MB at the constraint max).

### Code

```go
func oddOnlySieve(n int) int {
	if n < 3 {
		return 0 // primes strictly below n require n ≥ 3 (first prime is 2)
	}
	// Index i represents the odd number 2i+1 (i=0 ↔ 1, i=1 ↔ 3, ...).
	// 2i+1 < n  ⇔  i < n/2 for the sizes we need, so length n/2 covers all.
	isComposite := make([]bool, n/2)
	for i := 1; (2*i+1)*(2*i+1) < n; i++ { // p = 2i+1, seed while p² < n
		if isComposite[i] {
			continue // p already known composite — skip its multiples
		}
		p := 2*i + 1
		// Odd multiples only: p² is odd, and += 2p preserves oddness.
		for multiple := p * p; multiple < n; multiple += 2 * p {
			isComposite[multiple/2] = true // odd m maps to index m/2
		}
	}
	count := 1 // the prime 2, which the odd table cannot represent
	for i := 1; i < len(isComposite); i++ {
		if !isComposite[i] {
			count++ // odd survivor 2i+1 is prime
		}
	}
	return count
}
```

### Dry Run

Example 1: `n = 10`. Table length `10/2 = 5`; index i ↔ odd number 2i+1: `[1, 3, 5, 7, 9]`.

| Step | i | p = 2i+1 | p² < 10? | Marks (m, step 2p) | Table (index: number, ~~struck~~) |
|------|---|----------|----------|---------------------|-----------------------------------|
| 1 | 1 | 3 | yes (9 < 10) | m = 9 → index 4 | 0:1, 1:3, 2:5, 3:7, 4:~~9~~ |
| 2 | 2 | 5 | no (25 ≥ 10) | seed loop ends | unchanged |

Count: start at 1 (prime 2); unmarked indices 1, 2, 3 → numbers 3, 5, 7 → count = 1 + 3 = `4` ✔

---

## Approach 4 — Linear (Euler's) Sieve (Optimal)

### Intuition

The classic sieve re-marks numbers: 45 is crossed out by 3 (3·15) *and* by 5 (5·9). Euler's sieve gives every composite exactly one "certificate of compositeness": `c = p · i`, where `p` is the **smallest prime factor** (SPF) of `c`. Enumerate `i` upward, and for each `i` mark `i·p` for the known primes `p` in increasing order — but **stop as soon as `p` divides `i`**. Why stop there: if `p | i`, then for any larger prime `q`, the number `q·i` contains the factor `p` (inside `i`), so its SPF is `p`, not `q` — it will be generated later as `p · (q·i/p)`, and marking it now would be the duplicate work we are eliminating. Each composite is therefore marked exactly once → true O(n).

### Algorithm

1. If `n < 3`, return 0.
2. Allocate `isComposite[0..n-1]`; create an empty `primes` list.
3. For `i = 2 … n−1`:
   1. If `isComposite[i]` is false, `i` is prime — append to `primes`.
   2. For each `p` in `primes` (ascending): if `i·p ≥ n`, break; set `isComposite[i·p] = true`; if `i % p == 0`, break (p is i's SPF).
4. Return `len(primes)`.

### Complexity

- **Time:** O(n) — the outer loop is n steps and every inner-loop iteration marks a distinct composite exactly once, so total inner work equals the number of composites < n.
- **Space:** O(n) — the boolean table plus the primes list (π(n) ≈ n/ln n entries).

### Code

```go
func linearSieve(n int) int {
	if n < 3 {
		return 0 // nothing strictly below n can be prime for n ≤ 2
	}
	isComposite := make([]bool, n)
	primes := []int{} // discovered primes in increasing order
	for i := 2; i < n; i++ {
		if !isComposite[i] {
			primes = append(primes, i) // never marked → i is prime
		}
		// Mark i·p for each known prime p ≤ smallest prime factor of i.
		for _, p := range primes {
			if i*p >= n {
				break // product out of range — and grows with p, so stop
			}
			isComposite[i*p] = true // p is the smallest prime factor of i·p
			if i%p == 0 {
				// p divides i ⇒ for any larger prime q, q·i's smallest
				// factor is still p (via i), not q — stop to avoid re-marks.
				break
			}
		}
	}
	return len(primes)
}
```

### Dry Run

Example 1: `n = 10`.

| Step | i | i prime? | primes | Inner marks (i·p, stop rule) |
|------|---|----------|--------------------|------------------------------------------------|
| 1 | 2 | yes | [2] | 2·2=4 marked; 2 % 2 == 0 → break |
| 2 | 3 | yes | [2, 3] | 3·2=6 marked; 3·3=9 marked; 3 % 3 == 0 → break |
| 3 | 4 | no (marked) | [2, 3] | 4·2=8 marked; 4 % 2 == 0 → break (8's SPF is 2) |
| 4 | 5 | yes | [2, 3, 5] | 5·2=10 ≥ 10 → break immediately |
| 5 | 6 | no | [2, 3, 5] | 6·2=12 ≥ 10 → break |
| 6 | 7 | yes | [2, 3, 5, 7] | 7·2=14 ≥ 10 → break |
| 7 | 8 | no | [2, 3, 5, 7] | 8·2=16 ≥ 10 → break |
| 8 | 9 | no | [2, 3, 5, 7] | 9·2=18 ≥ 10 → break |

Every composite {4, 6, 8, 9} was marked exactly once. Return `len(primes)` = `4` ✔

---

## Key Takeaways

- **Trial division only to √x** — divisors pair up as `(d, x/d)`; one member of each pair is ≤ √x. This halves-the-exponent trick appears everywhere (factor counting, perfect squares).
- **Sieve = inverted work direction:** instead of testing each number against all primes, let each prime bulk-mark its multiples. Start at `p²` and seed only up to √n — both follow from "smaller factors already did the job".
- **Odd-only compression** (index `i ↔ 2i+1`, step `2p`) is the cheapest 2× win on any sieve; the same wheel idea extends to skipping multiples of 3 and 5.
- **Euler's linear sieve** marks each composite once via its smallest prime factor; the `if i % p == 0 break` line is the whole theorem. Bonus: recording `spf[i·p] = p` gives O(log c) factorisation of any c < n afterwards — a frequent contest sub-routine.
- Precompute-once, answer-many: for repeated "how many primes / is x prime" queries, build the sieve a single time and answer from the table in O(1).

---

## Related Problems

- LeetCode #263 — Ugly Number (divisibility reasoning on a single number)
- LeetCode #264 — Ugly Number II (generate numbers from prime factors)
- LeetCode #313 — Super Ugly Number (generalised multi-prime generation)
- LeetCode #2523 — Closest Prime Numbers in Range (sieve as a subroutine)
- LeetCode #952 — Largest Component Size by Common Factor (smallest-prime-factor factorisation)
