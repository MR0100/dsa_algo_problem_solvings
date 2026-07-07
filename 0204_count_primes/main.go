package main

import "fmt"

// ── Approach 1: Brute Force (Trial Division) ─────────────────────────────────
//
// bruteForce solves Count Primes by testing every candidate below n with
// trial division.
//
// Intuition:
//
//	Apply the definition directly: x is prime iff no d in [2, √x] divides it.
//	Divisors come in pairs (d, x/d) with the smaller one ≤ √x, so trialing up
//	to √x suffices. Do this for every x < n and count the survivors.
//
// Algorithm:
//  1. count = 0.
//  2. For x = 2 .. n−1: run isPrime(x) — try every d with d·d ≤ x; if any
//     divides x, x is composite.
//  3. Increment count for each prime; return count.
//
// Time:  O(n·√n) — n candidates, each trial-divided by up to √n numbers
//
//	(composites usually fail fast, but primes pay the full √x).
//
// Space: O(1) — no tables, just counters.
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

// ── Approach 2: Sieve of Eratosthenes ────────────────────────────────────────
//
// sieveOfEratosthenes solves Count Primes by crossing out composites instead
// of testing candidates.
//
// Intuition:
//
//	Flip the direction of work: rather than asking "who divides x?" for every
//	x, let each discovered prime p announce its own multiples p², p²+p, ...
//	as composite. Marking starts at p² because any smaller multiple k·p
//	(k < p) has a factor smaller than p and was already crossed out by it.
//	Whatever is never crossed out is prime.
//
// Algorithm:
//  1. isComposite = bool slice of length n (all false).
//  2. For i = 2 while i·i < n: if i is still unmarked (prime), mark
//     i², i²+i, i²+2i, ... < n as composite.
//  3. Count the unmarked indices in [2, n).
//
// Time:  O(n log log n) — the classic bound: Σ over primes p<√n of n/p.
// Space: O(n) — one boolean per number below n.
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

// ── Approach 3: Odd-Only Sieve ───────────────────────────────────────────────
//
// oddOnlySieve solves Count Primes with a sieve that stores and marks only
// odd numbers, halving memory and work.
//
// Intuition:
//
//	2 is the only even prime — count it once and never store an even number
//	again. Map odd number x to index x/2 (3→1, 5→2, 7→3, ...). An odd prime
//	p only needs to cross out its ODD multiples (even ones aren't stored),
//	which are p², p²+2p, p²+4p, ... — stepping by 2p keeps the values odd.
//
// Algorithm:
//  1. If n ≤ 2 return 0 (and n ≤ 2 has no odd primes either).
//  2. isComposite[i] represents odd number 2i+1, slice length n/2.
//  3. For each odd p = 2i+1 with p² < n and p unmarked: mark p², p²+2p, ...
//     (indices m/2).
//  4. Answer = 1 (for the prime 2) + unmarked indices i ≥ 1.
//
// Time:  O(n log log n) — same asymptotics as Approach 2 with ~half the
//
//	constant factor.
//
// Space: O(n/2) — one boolean per odd number below n.
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

// ── Approach 4: Linear (Euler's) Sieve (Optimal) ─────────────────────────────
//
// linearSieve solves Count Primes in true O(n) by crossing out every
// composite exactly once — via its smallest prime factor.
//
// Intuition:
//
//	The classic sieve marks 45 three times (3·15, 5·9, ...). Euler's sieve
//	assigns each composite a unique certificate: c = p·i where p is the
//	SMALLEST prime factor of c. Iterating i upward and, for each i, marking
//	i·p for known primes p — but stopping as soon as p divides i — produces
//	exactly that factorisation once and never again, so total marking work
//	is O(n).
//
// Algorithm:
//  1. primes = empty list; isComposite = bool slice of length n.
//  2. For i = 2 .. n−1: if i unmarked, append i to primes.
//  3. For each known prime p (ascending) while i·p < n: mark i·p, and break
//     when p divides i (p is i's smallest prime factor — bigger p would
//     mark i·p by a non-smallest factor, duplicating work).
//  4. Return len(primes).
//
// Time:  O(n) — every composite is marked exactly once.
// Space: O(n) — the composite table plus the list of primes (~n/ln n).
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

func main() {
	fmt.Println("=== Approach 1: Brute Force (Trial Division) ===")
	fmt.Printf("n=10       got=%d  expected 4\n", bruteForce(10)) // primes: 2,3,5,7
	fmt.Printf("n=0        got=%d  expected 0\n", bruteForce(0))
	fmt.Printf("n=1        got=%d  expected 0\n", bruteForce(1))
	fmt.Printf("n=2        got=%d  expected 0\n", bruteForce(2))           // strictly less than 2
	fmt.Printf("n=1000000  got=%d  expected 78498\n", bruteForce(1000000)) // π(10⁶)

	fmt.Println("=== Approach 2: Sieve of Eratosthenes ===")
	fmt.Printf("n=10       got=%d  expected 4\n", sieveOfEratosthenes(10))
	fmt.Printf("n=0        got=%d  expected 0\n", sieveOfEratosthenes(0))
	fmt.Printf("n=1        got=%d  expected 0\n", sieveOfEratosthenes(1))
	fmt.Printf("n=2        got=%d  expected 0\n", sieveOfEratosthenes(2))
	fmt.Printf("n=1000000  got=%d  expected 78498\n", sieveOfEratosthenes(1000000))
	fmt.Printf("n=5000000  got=%d  expected 348513\n", sieveOfEratosthenes(5000000)) // constraint max

	fmt.Println("=== Approach 3: Odd-Only Sieve ===")
	fmt.Printf("n=10       got=%d  expected 4\n", oddOnlySieve(10))
	fmt.Printf("n=0        got=%d  expected 0\n", oddOnlySieve(0))
	fmt.Printf("n=1        got=%d  expected 0\n", oddOnlySieve(1))
	fmt.Printf("n=2        got=%d  expected 0\n", oddOnlySieve(2))
	fmt.Printf("n=3        got=%d  expected 1\n", oddOnlySieve(3)) // just the prime 2
	fmt.Printf("n=1000000  got=%d  expected 78498\n", oddOnlySieve(1000000))
	fmt.Printf("n=5000000  got=%d  expected 348513\n", oddOnlySieve(5000000))

	fmt.Println("=== Approach 4: Linear (Euler's) Sieve (Optimal) ===")
	fmt.Printf("n=10       got=%d  expected 4\n", linearSieve(10))
	fmt.Printf("n=0        got=%d  expected 0\n", linearSieve(0))
	fmt.Printf("n=1        got=%d  expected 0\n", linearSieve(1))
	fmt.Printf("n=2        got=%d  expected 0\n", linearSieve(2))
	fmt.Printf("n=1000000  got=%d  expected 78498\n", linearSieve(1000000))
	fmt.Printf("n=5000000  got=%d  expected 348513\n", linearSieve(5000000))
}
