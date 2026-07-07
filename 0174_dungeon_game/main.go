package main

import (
	"fmt"
	"math"
)

// ── Approach 1: Brute Force (Recursive Path Exploration) ─────────────────────
//
// bruteForce solves Dungeon Game by recursively trying both moves from every
// room and combining results with the max(1, need-room) rule.
//
// Intuition:
//
//	Think backwards from any room: if I know the minimum health required to
//	survive from the room to the right and from the room below, the cheaper
//	of the two is what I must still have AFTER this room's effect. So before
//	the room I need min(right, down) - dungeon[i][j] — but never less than 1,
//	because health must stay positive at every moment (a huge orb cannot
//	"bank" health below 1). Recursing on this rule explores every monotone
//	path.
//
// Algorithm:
//  1. need(i, j) = max(1, 1 - dungeon[i][j]) at the princess cell.
//  2. Otherwise need(i, j) = max(1, min(need(i+1,j), need(i,j+1)) - dungeon[i][j]),
//     taking only in-bounds moves.
//  3. Answer = need(0, 0).
//
// Time:  O(2^(m+n)) — every cell branches into two recursive calls; the same
//
//	subproblems are recomputed exponentially many times.
//
// Space: O(m+n) — recursion depth equals the path length.
func bruteForce(dungeon [][]int) int {
	m, n := len(dungeon), len(dungeon[0])
	var need func(i, j int) int
	need = func(i, j int) int {
		// Base case: the princess room — leave it with at least 1 HP.
		if i == m-1 && j == n-1 {
			return max(1, 1-dungeon[i][j])
		}
		onward := math.MaxInt32 // cheapest requirement among legal next rooms
		if i+1 < m {
			onward = min(onward, need(i+1, j)) // option: move down
		}
		if j+1 < n {
			onward = min(onward, need(i, j+1)) // option: move right
		}
		// Must hold `onward` after this room's effect; clamp at 1 because
		// health may never be 0 or below, even inside a healing room.
		return max(1, onward-dungeon[i][j])
	}
	return need(0, 0)
}

// ── Approach 2: DP Top-Down (Memoized Recursion) ─────────────────────────────
//
// dpTopDown solves Dungeon Game with the same recursion as Approach 1 plus a
// memo table, collapsing the exponential tree to one evaluation per cell.
//
// Intuition:
//
//	need(i, j) depends only on (i, j) — not on how we got there — so the
//	recursion has just m×n distinct subproblems. Cache each answer the first
//	time it is computed; every later visit is a table lookup.
//
// Algorithm:
//  1. memo[i][j] = 0 means "not computed yet" (real answers are always ≥ 1).
//  2. Run the Approach-1 recursion, checking/filling memo around it.
//  3. Answer = need(0, 0).
//
// Time:  O(m·n) — each cell computed once, O(1) work apiece.
// Space: O(m·n) — the memo table, plus O(m+n) recursion stack.
func dpTopDown(dungeon [][]int) int {
	m, n := len(dungeon), len(dungeon[0])
	memo := make([][]int, m)
	for i := range memo {
		memo[i] = make([]int, n) // zero value 0 = "uncomputed" sentinel
	}
	var need func(i, j int) int
	need = func(i, j int) int {
		if i == m-1 && j == n-1 {
			return max(1, 1-dungeon[i][j]) // princess room base case
		}
		if memo[i][j] != 0 {
			return memo[i][j] // already solved — reuse
		}
		onward := math.MaxInt32
		if i+1 < m {
			onward = min(onward, need(i+1, j)) // requirement if we go down
		}
		if j+1 < n {
			onward = min(onward, need(i, j+1)) // requirement if we go right
		}
		memo[i][j] = max(1, onward-dungeon[i][j]) // clamp: health ≥ 1 always
		return memo[i][j]
	}
	return need(0, 0)
}

// ── Approach 3: DP Bottom-Up (2D Table from the Princess) ────────────────────
//
// dpBottomUp solves Dungeon Game iteratively, filling a table of minimum
// health requirements from the bottom-right corner back to the entrance.
//
// Intuition:
//
//	The recursion's dependency order is fixed: (i, j) needs (i+1, j) and
//	(i, j+1). Sweeping rows bottom→top and columns right→left guarantees
//	both are ready. A sentinel row/column of +∞ with two 1-cells beside the
//	princess lets one formula cover base case, edges, and interior alike.
//	Forward DP ("max health reaching each cell") fails because a path can
//	look rich early but die later — the state would need two numbers;
//	backward DP needs just one.
//
// Algorithm:
//  1. Build (m+1)×(n+1) table filled with +∞; set need[m][n-1] = need[m-1][n] = 1
//     (virtual rooms "after" the princess: exit alive with exactly 1 HP).
//  2. For i = m-1 … 0, j = n-1 … 0:
//     need[i][j] = max(1, min(need[i+1][j], need[i][j+1]) - dungeon[i][j]).
//  3. Answer = need[0][0].
//
// Time:  O(m·n) — one constant-time cell fill each.
// Space: O(m·n) — the (m+1)×(n+1) requirement table.
func dpBottomUp(dungeon [][]int) int {
	m, n := len(dungeon), len(dungeon[0])
	// One extra sentinel row and column so every cell uses the same formula.
	need := make([][]int, m+1)
	for i := range need {
		need[i] = make([]int, n+1)
		for j := range need[i] {
			need[i][j] = math.MaxInt32 // walls: never the min() winner
		}
	}
	// Virtual cells flanking the princess: surviving means exiting with 1 HP.
	need[m][n-1], need[m-1][n] = 1, 1
	for i := m - 1; i >= 0; i-- {
		for j := n - 1; j >= 0; j-- {
			// Cheapest requirement among the two onward rooms...
			req := min(need[i+1][j], need[i][j+1]) - dungeon[i][j]
			if req < 1 {
				req = 1 // ...clamped: health can never sit at 0 or below
			}
			need[i][j] = req
		}
	}
	return need[0][0]
}

// ── Approach 4: DP Bottom-Up, Space Optimized (1D Rolling Row) ───────────────
//
// dpSpaceOptimized solves Dungeon Game with the Approach-3 recurrence but
// keeps only one row of the table alive.
//
// Intuition:
//
//	Filling row i touches only row i+1 (below) and the already-updated part
//	of row i (right). So a single slice works: before writing index j,
//	dp[j] still holds row i+1's value ("below") and dp[j+1] already holds
//	row i's value ("right"). Seed the slice so the princess cell's formula
//	falls out of the same code path.
//
// Algorithm:
//  1. dp = length n+1, all +∞ except dp[n-1] = 1 (the virtual room below
//     the princess).
//  2. For i = m-1 … 0, j = n-1 … 0: dp[j] = max(1, min(dp[j], dp[j+1]) - dungeon[i][j]).
//  3. Answer = dp[0].
//
// Time:  O(m·n) — identical work to Approach 3.
// Space: O(n) — one rolling row (n+1 ints).
func dpSpaceOptimized(dungeon [][]int) int {
	m, n := len(dungeon), len(dungeon[0])
	dp := make([]int, n+1)
	for j := range dp {
		dp[j] = math.MaxInt32 // sentinel: off-grid moves are never chosen
	}
	dp[n-1] = 1 // virtual room below the princess: arrive there with 1 HP
	for i := m - 1; i >= 0; i-- {
		for j := n - 1; j >= 0; j-- {
			// dp[j] = requirement below (old row), dp[j+1] = right (new row).
			req := min(dp[j], dp[j+1]) - dungeon[i][j]
			if req < 1 {
				req = 1 // health floor: even big orbs can't push need below 1
			}
			dp[j] = req // overwrite in place: row i+1's slot becomes row i's
		}
	}
	return dp[0]
}

// ── Approach 5: Binary Search on Initial Health + Greedy Simulation ──────────
//
// binarySearchOnAnswer solves Dungeon Game by binary-searching the starting
// health and checking each guess with a forward max-health DP.
//
// Intuition:
//
//	Survivability is monotone: if H hit points suffice, H+1 certainly does
//	(the knight's health along any path just shifts up by 1). Monotone
//	yes/no ⇒ binary search the smallest yes. Checking a FIXED start H is a
//	easy forward DP: health changes additively, so at each cell the best we
//	can do is arrive with maximum health; cells where even the maximum
//	drops to 0 or below are dead ends.
//
// Algorithm:
//  1. feasible(H): dp over rows, best[j] = max health standing on (i, j)
//     having stayed ≥ 1 the whole way, or -∞ if unreachable alive;
//     transition best = max(from above, from left) + dungeon[i][j].
//  2. Binary search the smallest H in [1, 1000·(m+n)+1] with feasible(H).
//
// Time:  O(m·n · log(1000·(m+n))) — one grid DP per binary-search probe.
// Space: O(n) — rolling row for the feasibility check.
func binarySearchOnAnswer(dungeon [][]int) int {
	m, n := len(dungeon), len(dungeon[0])
	const dead = math.MinInt32 // marker: cannot reach this cell alive

	// feasible reports whether the knight can reach the princess alive
	// when starting with exactly `start` health points.
	feasible := func(start int) bool {
		best := make([]int, n+1) // best[j+1] = max health at (i, j); best[0] = left wall
		for j := range best {
			best[j] = dead
		}
		for i := 0; i < m; i++ {
			for j := 0; j < n; j++ {
				var h int
				if i == 0 && j == 0 {
					h = start + dungeon[0][0] // entering the first room applies its effect
				} else {
					// Best surviving predecessor: above (old best[j+1]) or left (new best[j]).
					prev := max(best[j+1], best[j])
					if prev == dead {
						best[j+1] = dead // no way to arrive alive
						continue
					}
					h = prev + dungeon[i][j] // health is purely additive
				}
				if h <= 0 {
					h = dead // died on entry — 0 or below is instant death
				}
				best[j+1] = h
			}
		}
		return best[n] != dead // alive at the princess cell?
	}

	// Any single path visits m+n-1 rooms, each draining at most 1000 HP,
	// so 1000·(m+n)+1 is always survivable → a valid search upper bound.
	lo, hi := 1, 1000*(m+n)+1
	for lo < hi {
		mid := lo + (hi-lo)/2 // avoid (lo+hi) overflow
		if feasible(mid) {
			hi = mid // mid works — try smaller
		} else {
			lo = mid + 1 // mid dies — need more health
		}
	}
	return lo // smallest health that survives
}

func main() {
	ex1 := func() [][]int { return [][]int{{-2, -3, 3}, {-5, -10, 1}, {10, 30, -5}} }
	ex2 := func() [][]int { return [][]int{{0}} }

	fmt.Println("=== Approach 1: Brute Force (Recursive Path Exploration) ===")
	fmt.Printf("dungeon=[[-2,-3,3],[-5,-10,1],[10,30,-5]] got=%-3d expected 7\n", bruteForce(ex1()))
	fmt.Printf("dungeon=[[0]]                             got=%-3d expected 1\n", bruteForce(ex2()))

	fmt.Println("=== Approach 2: DP Top-Down (Memoized Recursion) ===")
	fmt.Printf("dungeon=[[-2,-3,3],[-5,-10,1],[10,30,-5]] got=%-3d expected 7\n", dpTopDown(ex1()))
	fmt.Printf("dungeon=[[0]]                             got=%-3d expected 1\n", dpTopDown(ex2()))

	fmt.Println("=== Approach 3: DP Bottom-Up (2D Table from the Princess) ===")
	fmt.Printf("dungeon=[[-2,-3,3],[-5,-10,1],[10,30,-5]] got=%-3d expected 7\n", dpBottomUp(ex1()))
	fmt.Printf("dungeon=[[0]]                             got=%-3d expected 1\n", dpBottomUp(ex2()))

	fmt.Println("=== Approach 4: DP Bottom-Up, Space Optimized (1D Rolling Row) ===")
	fmt.Printf("dungeon=[[-2,-3,3],[-5,-10,1],[10,30,-5]] got=%-3d expected 7\n", dpSpaceOptimized(ex1()))
	fmt.Printf("dungeon=[[0]]                             got=%-3d expected 1\n", dpSpaceOptimized(ex2()))

	fmt.Println("=== Approach 5: Binary Search on Initial Health + Greedy Simulation ===")
	fmt.Printf("dungeon=[[-2,-3,3],[-5,-10,1],[10,30,-5]] got=%-3d expected 7\n", binarySearchOnAnswer(ex1()))
	fmt.Printf("dungeon=[[0]]                             got=%-3d expected 1\n", binarySearchOnAnswer(ex2()))
}
