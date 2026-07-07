package main

import "fmt"

// ── Approach 1: Brute Force (Check Every Number) ─────────────────────────────
//
// bruteForce solves Count Numbers with Unique Digits by testing each number in
// [0, 10^n) for digit uniqueness.
//
// Intuition:
//
//	The definition is direct: count x in [0, 10^n) whose decimal digits are
//	all distinct. So enumerate every such x and, for each, walk its digits
//	using a 10-slot "seen" boolean array; the first repeated digit disqualifies
//	it. 0 counts (it has the single digit 0). Correct but exponential.
//
// Algorithm:
//  1. limit = 10^n.
//  2. For x in [0, limit): reset a seen[10] mask, extract digits of x,
//     and mark them; a duplicate ⇒ not unique.
//  3. Count the unique ones.
//
// Time:  O(10^n · d) where d ≈ n is the digit count — enumerates the range.
// Space: O(1) — a fixed 10-slot seen array.
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

// hasUniqueDigits reports whether every decimal digit of x is distinct.
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

// ── Approach 2: Combinatorics (Optimal) ──────────────────────────────────────
//
// combinatorics solves Count Numbers with Unique Digits by counting, per digit
// length, how many numbers of that length have all-distinct digits.
//
// Intuition:
//
//	Count by number of digits k (1..n) and add the single number 0.
//	 - 1-digit numbers (1..9): 9 of them all trivially unique. Plus 0 ⇒ 10.
//	 - k-digit numbers (k >= 2): the leading digit has 9 choices (1..9, no
//	   leading zero); the next has 9 choices (0..9 minus the one used); then
//	   8, 7, ... So count_k = 9 · 9 · 8 · … · (10 - k + 1).
//	Sum over k = 1..n. Because only 10 distinct digits exist, no unique-digit
//	number has more than 10 digits, so for n > 10 the extra terms are 0.
//
// Algorithm:
//  1. If n == 0 ⇒ only the number 0 ⇒ return 1.
//  2. total = 10 (all 1-digit numbers 0..9).
//  3. uniqueDigits = 9, available = 9. For each extra digit length up to n:
//     uniqueDigits *= available; available--; total += uniqueDigits.
//  4. Return total.
//
// Time:  O(min(n, 10)) — a handful of multiplications.
// Space: O(1).
func combinatorics(n int) int {
	if n == 0 {
		return 1 // only 0 lies in [0, 1)
	}
	total := 10          // all one-digit numbers 0..9
	uniqueDigits := 9    // count of k-digit unique numbers, starts at k=1 (9)
	availableDigits := 9 // choices for the next position (0..9 minus used)
	for k := 2; k <= n && availableDigits > 0; k++ {
		uniqueDigits *= availableDigits // extend numbers by one more distinct digit
		total += uniqueDigits           // add all k-digit unique numbers
		availableDigits--               // one fewer digit remains for the next slot
	}
	return total
}

// ── Approach 3: Backtracking (DFS over digit choices) ────────────────────────
//
// backtracking solves Count Numbers with Unique Digits by DFS-building numbers
// digit by digit, never reusing a digit, and counting each valid prefix.
//
// Intuition:
//
//	Every unique-digit number is a path in a tree: pick a first digit (1..9),
//	then any unused digit for each following position, up to n digits. Each
//	such node (of length 1..n) is one number. Explore all paths with a used[]
//	mask and count nodes. This mirrors the combinatorics but constructs the
//	numbers explicitly, which is instructive for the "generate all" variant.
//
// Algorithm:
//  1. Answer starts at 1 to account for the number 0.
//  2. For each starting digit 1..9, DFS with depth 1: at each node count it,
//     then if depth < n try every unused digit 0..9 as the next digit.
//  3. Sum all counted nodes.
//
// Time:  O(sum of unique counts) ≤ O(10!) but bounded by ~8.9M for n=10.
// Space: O(n) recursion + O(10) used mask.
func backtracking(n int) int {
	if n == 0 {
		return 1
	}
	var used [10]bool // digits already placed on the current path
	count := 1        // pre-count the number 0

	// dfs counts every valid number reachable by extending the current prefix.
	var dfs func(depth int)
	dfs = func(depth int) {
		if depth == n {
			return // cannot append more digits
		}
		for d := 0; d <= 9; d++ {
			if used[d] {
				continue // digit already used on this path
			}
			used[d] = true  // place digit d
			count++         // this longer prefix is itself a valid number
			dfs(depth + 1)  // extend further
			used[d] = false // backtrack: free digit d
		}
	}

	// Leading digit must be 1..9 (no leading zeros) — each is a 1-digit number.
	for d := 1; d <= 9; d++ {
		used[d] = true
		count++ // the 1-digit number d
		dfs(1)  // extend to 2..n digits
		used[d] = false
	}
	return count
}

func main() {
	fmt.Println("=== Approach 1: Brute Force ===")
	fmt.Printf("n=2   got=%d  expected 91\n", bruteForce(2))
	fmt.Printf("n=0   got=%d  expected 1\n", bruteForce(0))
	fmt.Printf("n=1   got=%d  expected 10\n", bruteForce(1))
	fmt.Printf("n=3   got=%d  expected 739\n", bruteForce(3))

	fmt.Println("=== Approach 2: Combinatorics (Optimal) ===")
	fmt.Printf("n=2   got=%d  expected 91\n", combinatorics(2))
	fmt.Printf("n=0   got=%d  expected 1\n", combinatorics(0))
	fmt.Printf("n=1   got=%d  expected 10\n", combinatorics(1))
	fmt.Printf("n=3   got=%d  expected 739\n", combinatorics(3))
	fmt.Printf("n=11  got=%d  expected 8877691\n", combinatorics(11)) // caps at 10 distinct digits

	fmt.Println("=== Approach 3: Backtracking ===")
	fmt.Printf("n=2   got=%d  expected 91\n", backtracking(2))
	fmt.Printf("n=0   got=%d  expected 1\n", backtracking(0))
	fmt.Printf("n=1   got=%d  expected 10\n", backtracking(1))
	fmt.Printf("n=3   got=%d  expected 739\n", backtracking(3))
}
