package main

import (
	"fmt"
	"math"
)

// ── Approach 1: DP Bottom-Up (Coin-Change Style) ─────────────────────────────
//
// dpBottomUp solves Perfect Squares by treating perfect squares as coin
// denominations and finding the fewest "coins" summing to n.
//
// Intuition:
//
//	Let dp[i] = minimum number of perfect squares summing to i. To form i we
//	pick some square j*j (j*j <= i) as the last term, leaving i - j*j to be
//	formed optimally: dp[i] = min over j of dp[i - j*j] + 1. dp[0] = 0.
//
// Algorithm:
//  1. dp[0] = 0; dp[i>0] = +∞.
//  2. For i = 1..n: for each j with j*j <= i: dp[i] = min(dp[i], dp[i-j*j]+1).
//  3. Return dp[n].
//
// Time:  O(n * sqrt(n)) — for each i we try up to sqrt(i) squares.
// Space: O(n) — the dp table.
func dpBottomUp(n int) int {
	dp := make([]int, n+1) // dp[i] = fewest squares that sum to i
	for i := 1; i <= n; i++ {
		dp[i] = math.MaxInt32 // start each as "unreachable / infinity"
		for j := 1; j*j <= i; j++ {
			// take square j*j as the last term; add 1 to the best for the remainder
			if dp[i-j*j]+1 < dp[i] {
				dp[i] = dp[i-j*j] + 1
			}
		}
	}
	return dp[n]
}

// ── Approach 2: BFS (Shortest Path in a Graph of Sums) ────────────────────────
//
// bfs solves Perfect Squares by viewing each integer 0..n as a node and each
// perfect square as an edge; the fewest squares is the shortest path 0→n.
//
// Intuition:
//
//	From value v we can jump to v + j*j for every square j*j <= n - v. Each jump
//	costs one square. BFS from 0 explores by increasing number of jumps, so the
//	first time we reach n is the minimum count.
//
// Algorithm:
//  1. Precompute squares <= n.
//  2. BFS level-by-level from 0; level number = squares used so far.
//  3. For each frontier value, jump by every square; first to reach n wins.
//
// Time:  O(n * sqrt(n)) — each of ≤ n nodes expands ≤ sqrt(n) edges.
// Space: O(n) — visited set + queue.
func bfs(n int) int {
	if n == 0 {
		return 0
	}
	// build the list of usable square values
	squares := []int{}
	for j := 1; j*j <= n; j++ {
		squares = append(squares, j*j)
	}
	visited := make([]bool, n+1) // avoid revisiting a value
	queue := []int{0}            // start BFS at sum 0
	visited[0] = true
	level := 0 // number of squares used to reach the current frontier
	for len(queue) > 0 {
		level++         // one more square added at this BFS layer
		next := []int{} // frontier for the next layer
		for _, v := range queue {
			for _, s := range squares {
				nv := v + s
				if nv == n {
					return level // reached target with `level` squares
				}
				if nv < n && !visited[nv] {
					visited[nv] = true      // mark to prevent duplicates
					next = append(next, nv) // explore later
				}
				if nv > n {
					break // squares are ascending; further ones overshoot too
				}
			}
		}
		queue = next
	}
	return level
}

// ── Approach 3: Lagrange's Four-Square Theorem (Optimal, Math) ────────────────
//
// mathFourSquare solves Perfect Squares in O(sqrt(n)) using number theory.
//
// Intuition:
//
//	Lagrange: every natural number is the sum of at most four squares. So the
//	answer is 1, 2, 3, or 4. We classify:
//	  - 1 if n itself is a perfect square.
//	  - 4 iff n = 4^a * (8b + 7)  (Legendre's three-square theorem: these are
//	    exactly the numbers NOT expressible as three squares).
//	  - 2 if n = a² + b² for some a (test all a with a² <= n).
//	  - otherwise 3.
//
// Algorithm:
//  1. If n is a perfect square → 1.
//  2. Strip factors of 4; if remainder ≡ 7 (mod 8) → 4.
//  3. Try a from 1 while a² <= n: if n - a² is a perfect square → 2.
//  4. Otherwise → 3.
//
// Time:  O(sqrt(n)) — the two-square check dominates.
// Space: O(1).
func mathFourSquare(n int) int {
	isSquare := func(x int) bool {
		r := int(math.Sqrt(float64(x))) // candidate integer root
		return r*r == x                 // exact when r*r reproduces x
	}
	if isSquare(n) {
		return 1 // n is itself a perfect square
	}
	// Legendre's three-square theorem: answer is 4 iff n = 4^a*(8b+7).
	m := n
	for m%4 == 0 { // strip out factors of 4
		m /= 4
	}
	if m%8 == 7 {
		return 4
	}
	// Try to write n as a sum of two squares.
	for a := 1; a*a <= n; a++ {
		if isSquare(n - a*a) {
			return 2
		}
	}
	return 3 // not 1, 2, or 4 → must be 3 by Lagrange
}

func main() {
	fmt.Println("=== Approach 1: DP Bottom-Up ===")
	fmt.Println(dpBottomUp(12)) // expected 3
	fmt.Println(dpBottomUp(13)) // expected 2

	fmt.Println("=== Approach 2: BFS ===")
	fmt.Println(bfs(12)) // expected 3
	fmt.Println(bfs(13)) // expected 2

	fmt.Println("=== Approach 3: Lagrange's Four-Square (Optimal) ===")
	fmt.Println(mathFourSquare(12)) // expected 3
	fmt.Println(mathFourSquare(13)) // expected 2
}
