package main

import "fmt"

// ── Approach 1: Dynamic Programming (Bottom-Up) ──────────────────────────────
//
// dpBottomUp solves Nim Game by computing, for every heap size 1..n, whether
// the player to move can force a win, using the results for smaller heaps.
//
// Intuition:
//
//	A player wins from a heap of size k if there EXISTS a move (take 1, 2, or 3)
//	that leaves the opponent in a LOSING position. So win[k] is true iff at
//	least one of win[k-1], win[k-2], win[k-3] is false (that "false" is the
//	trap we hand to the opponent). Base: win[0] = false (no stones left → the
//	player to move already lost, since the person who took the last stone won).
//
// Algorithm:
//
//  1. win[0] = false.
//  2. For k = 1..n: win[k] = !win[k-1] || !win[k-2] || !win[k-3]
//     (treating win[<0] as true so it never contributes a winning move).
//  3. Answer is win[n].
//
// Time:  O(n)   — fill one boolean per heap size.
// Space: O(n)   — the win table (n grows to 2^31, so this is illustrative only).
func dpBottomUp(n int) bool {
	if n <= 0 {
		return false // no stones: the player to move cannot take the last stone
	}
	win := make([]bool, n+1) // win[k] = current player wins with k stones
	win[0] = false           // 0 stones → current player has no move → loses
	for k := 1; k <= n; k++ {
		// Try each legal removal of 1, 2, or 3 stones. If any of them leaves the
		// opponent in a losing state (win[k-take] == false), we win.
		for take := 1; take <= 3 && take <= k; take++ {
			if !win[k-take] { // opponent would be stuck in a losing position
				win[k] = true
				break // one winning move is enough
			}
		}
	}
	return win[n]
}

// ── Approach 2: Pattern Recognition / Math (Optimal) ─────────────────────────
//
// mathModulo solves Nim Game with the closed-form observation that the first
// player loses exactly when n is a multiple of 4.
//
// Intuition:
//
//	Run the DP by hand: sizes 1,2,3 are wins (take everything), size 4 is a
//	loss (whatever you take — 1,2,3 — leaves 3,2,1, all wins for the opponent).
//	From 5,6,7 you can always drop the opponent to exactly 4 (a loss for them),
//	so those are wins; 8 is a loss again. The losing positions are exactly the
//	multiples of 4. Strategy: always move to make the pile a multiple of 4 for
//	the opponent; you can only do that if n is NOT already a multiple of 4.
//
// Algorithm:
//
//	Return n % 4 != 0.
//
// Time:  O(1) — a single modulo.
// Space: O(1) — no extra memory.
func mathModulo(n int) bool {
	return n%4 != 0 // first player loses iff n is divisible by 4
}

func main() {
	fmt.Println("=== Approach 1: DP Bottom-Up ===")
	fmt.Println(dpBottomUp(4)) // expected false
	fmt.Println(dpBottomUp(1)) // expected true
	fmt.Println(dpBottomUp(2)) // expected true

	fmt.Println("=== Approach 2: Math Modulo (Optimal) ===")
	fmt.Println(mathModulo(4)) // expected false
	fmt.Println(mathModulo(1)) // expected true
	fmt.Println(mathModulo(2)) // expected true
}
