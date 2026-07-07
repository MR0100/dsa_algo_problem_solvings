package main

import "fmt"

// ── Approach 1: DP Top-Down (Memoized Recursion) ─────────────────────────────
//
// dpTopDown solves Paint House by recursion + memoization.
//
// Intuition:
//
//	Paint houses left to right. The only thing that constrains house i is the
//	color chosen for house i-1 (adjacent houses must differ). So define
//	solve(i, prevColor) = minimum cost to paint houses i..n-1 given house i-1
//	was painted prevColor. At house i we try each of the 3 colors that is not
//	prevColor, add its cost, and recurse. Memoize on (i, prevColor).
//
// Algorithm:
//  1. solve(i, prev): if i == n return 0.
//  2. For color c in {0,1,2} with c != prev: candidate = costs[i][c] + solve(i+1, c).
//  3. Return the minimum candidate; cache it.
//  4. Answer = solve(0, -1) (no previous color for the first house).
//
// Time:  O(n * 3 * 3) = O(n) — n*4 states, each tries 3 colors.
// Space: O(n) — memo table + recursion stack.
func dpTopDown(costs [][]int) int {
	n := len(costs)
	if n == 0 {
		return 0
	}
	// memo[i][prev+1]: prev ranges -1..2, shift by +1 to index 0..3.
	memo := make([][4]int, n)
	seen := make([][4]bool, n)

	var solve func(i, prev int) int
	solve = func(i, prev int) int {
		if i == n { // painted every house
			return 0
		}
		if seen[i][prev+1] { // already computed this state
			return memo[i][prev+1]
		}
		best := 1 << 30 // large sentinel (no valid cost can reach this)
		for c := 0; c < 3; c++ {
			if c == prev { // adjacent houses cannot share a color
				continue
			}
			cost := costs[i][c] + solve(i+1, c) // paint i with c, recurse
			if cost < best {
				best = cost
			}
		}
		seen[i][prev+1] = true // memoize the result for this state
		memo[i][prev+1] = best
		return best
	}
	return solve(0, -1) // house 0 has no left neighbor
}

// ── Approach 2: DP Bottom-Up (Full Table) ────────────────────────────────────
//
// dpBottomUp solves Paint House with an explicit DP table filled row by row.
//
// Intuition:
//
//	dp[i][c] = minimum total cost to paint houses 0..i with house i painted
//	color c. To paint house i color c, house i-1 must be one of the other two
//	colors, so dp[i][c] = costs[i][c] + min(dp[i-1][other two]). Answer is the
//	minimum over the last row.
//
// Algorithm:
//  1. dp[0] = costs[0].
//  2. For i = 1..n-1, c = 0..2: dp[i][c] = costs[i][c] + min(dp[i-1] over c'!=c).
//  3. Return min(dp[n-1]).
//
// Time:  O(n) — n rows, constant work per cell.
// Space: O(n) — the full dp table (kept for clarity).
func dpBottomUp(costs [][]int) int {
	n := len(costs)
	if n == 0 {
		return 0
	}
	dp := make([][3]int, n)
	dp[0] = [3]int{costs[0][0], costs[0][1], costs[0][2]} // base row = first house's costs
	for i := 1; i < n; i++ {
		// Each color depends on the best of the OTHER two colors above it.
		dp[i][0] = costs[i][0] + min2(dp[i-1][1], dp[i-1][2])
		dp[i][1] = costs[i][1] + min2(dp[i-1][0], dp[i-1][2])
		dp[i][2] = costs[i][2] + min2(dp[i-1][0], dp[i-1][1])
	}
	last := dp[n-1]
	return min2(last[0], min2(last[1], last[2])) // cheapest way to finish
}

// ── Approach 3: DP Bottom-Up O(1) Space (Optimal) ────────────────────────────
//
// dpRolling solves Paint House keeping only the previous row's three values.
//
// Intuition:
//
//	dp[i] only ever reads dp[i-1]. So we do not need the whole table — three
//	rolling variables (red, blue, green) for the previous house suffice. Update
//	them in place per house.
//
// Algorithm:
//  1. r,b,g = costs[0].
//  2. For each next house: newR = costs[i][0]+min(b,g), etc.; then r,b,g = new*.
//  3. Return min(r,b,g).
//
// Time:  O(n) — one pass.
// Space: O(1) — three scalars.
func dpRolling(costs [][]int) int {
	n := len(costs)
	if n == 0 {
		return 0
	}
	r, b, g := costs[0][0], costs[0][1], costs[0][2] // prev house's 3 totals
	for i := 1; i < n; i++ {
		// Compute this house's totals from the previous three, then roll over.
		nr := costs[i][0] + min2(b, g) // red now: cheaper of prev blue/green
		nb := costs[i][1] + min2(r, g) // blue now: cheaper of prev red/green
		ng := costs[i][2] + min2(r, b) // green now: cheaper of prev red/blue
		r, b, g = nr, nb, ng
	}
	return min2(r, min2(b, g)) // best final choice
}

// min2 returns the smaller of two ints.
func min2(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func main() {
	ex1 := [][]int{{17, 2, 17}, {16, 16, 5}, {14, 3, 19}}
	ex2 := [][]int{{7, 6, 2}}

	fmt.Println("=== Approach 1: DP Top-Down (Memoized) ===")
	fmt.Println(dpTopDown(ex1)) // expected 10
	fmt.Println(dpTopDown(ex2)) // expected 2

	fmt.Println("=== Approach 2: DP Bottom-Up (Full Table) ===")
	fmt.Println(dpBottomUp(ex1)) // expected 10
	fmt.Println(dpBottomUp(ex2)) // expected 2

	fmt.Println("=== Approach 3: DP Bottom-Up O(1) Space (Optimal) ===")
	fmt.Println(dpRolling(ex1)) // expected 10
	fmt.Println(dpRolling(ex2)) // expected 2
}
