package main

import "fmt"

// ── Approach 1: Plain Backtracking (Brute Force) ─────────────────────────────
//
// backtracking solves Flip Game II by trying every legal "++"→"--" move and
// recursing: the current player wins if ANY move leaves the opponent unable
// to win.
//
// Intuition:
//
//	This is a two-player game with perfect information. The player to move wins
//	if there is at least one move after which the opponent (now to move on the
//	new string) LOSES. That is a direct minimax recursion: canWin(s) is true
//	iff some flippable pair leads to a state where canWin(next) is false.
//
// Algorithm:
//
//  1. Scan for "++"; for each, flip to "--" in a byte copy.
//  2. If canWin(next) is false, the opponent loses → return true.
//  3. Restore is implicit (each branch uses a fresh copy).
//  4. If no move wins, return false.
//
// Time:  O(n!!) roughly (game-tree size) — exponential, each state re-explored.
// Space: O(n·depth) — a new string per move along a path of depth up to n/2.
func backtracking(s string) bool {
	b := []byte(s) // mutable working copy of the board
	for i := 0; i+1 < len(b); i++ {
		if b[i] == '+' && b[i+1] == '+' { // a legal move exists here
			b[i], b[i+1] = '-', '-'                 // make the move
			opponentWins := backtracking(string(b)) // opponent plays on the new board
			b[i], b[i+1] = '+', '+'                 // undo the move (restore for next candidate)
			if !opponentWins {                      // opponent can't win → we win by playing here
				return true
			}
		}
	}
	return false // no move leaves the opponent losing → current player loses
}

// ── Approach 2: Backtracking + Memoization ───────────────────────────────────
//
// memoBacktracking solves Flip Game II like Approach 1 but caches each board
// string's win/lose result so identical positions reachable by different move
// orders are computed once.
//
// Intuition:
//
//	Many move sequences reach the same board (order of far-apart flips doesn't
//	matter). Memoizing on the exact string collapses that redundancy. It is a
//	huge constant/exponent improvement while keeping the same minimax logic.
//
// Algorithm:
//
//  1. If s is cached, return the cached result.
//  2. Otherwise run the minimax scan; store and return the result.
//
// Time:  O(2^n) worst case bounded by number of distinct reachable boards.
// Space: O(2^n) memo in the worst case + recursion stack O(n).
func memoBacktracking(s string) bool {
	memo := map[string]bool{}
	var solve func(cur string) bool
	solve = func(cur string) bool {
		if v, ok := memo[cur]; ok { // already decided this exact board
			return v
		}
		b := []byte(cur)
		result := false
		for i := 0; i+1 < len(b); i++ {
			if b[i] == '+' && b[i+1] == '+' {
				b[i], b[i+1] = '-', '-'
				win := !solve(string(b)) // we win if opponent loses on the new board
				b[i], b[i+1] = '+', '+'
				if win {
					result = true
					break // one winning move suffices
				}
			}
		}
		memo[cur] = result // remember the verdict for this board
		return result
	}
	return solve(s)
}

// ── Approach 3: Sprague-Grundy Theorem (Optimal) ─────────────────────────────
//
// spragueGrundy solves Flip Game II by game theory: the board splits into
// independent maximal runs of '+', each an impartial game; the whole game is
// their XOR (nim-sum). The first player wins iff the nim-sum is nonzero.
//
// Intuition:
//
//	A move flips two ADJACENT pluses inside some run of consecutive '+'. A run
//	of length k, when you flip positions splitting it, becomes two shorter
//	independent runs. By Sprague-Grundy, the Grundy number of the whole game is
//	the XOR of the Grundy numbers of its independent runs. G(k) is the mex of
//	G(i-2) XOR G(k-i) over all internal flip positions. Nonzero total ⇒ win.
//
// Algorithm:
//
//  1. Precompute Grundy g[0..maxRun] via mex over splits.
//  2. Scan s, measure each maximal '+' run length; XOR g[len] into a running
//     nim-sum.
//  3. Return nimSum != 0.
//
// Time:  O(m² + n) — m = longest run for Grundy table, n to scan the string.
// Space: O(m) — the Grundy table.
func spragueGrundy(s string) bool {
	n := len(s)
	// Collect the lengths of maximal consecutive '+' runs.
	runs := []int{}
	count := 0
	for i := 0; i < n; i++ {
		if s[i] == '+' {
			count++ // extend the current run
		} else if count > 0 {
			runs = append(runs, count) // run ended, record it
			count = 0
		}
	}
	if count > 0 {
		runs = append(runs, count) // trailing run
	}
	maxRun := 0
	for _, r := range runs {
		if r > maxRun {
			maxRun = r
		}
	}
	// Grundy numbers: g[k] for a plus-run of length k.
	g := make([]int, maxRun+1)
	for k := 2; k <= maxRun; k++ {
		seen := map[int]bool{}
		// A flip at internal boundary i (1-based between pos i-1,i) turns a run of
		// length k into two runs of lengths (i-1)-1 and k-(i+1)... modeled as
		// splitting into g[left] XOR g[right] for every legal flip position.
		for left := 0; left <= k-2; left++ {
			right := k - 2 - left         // remaining pluses after removing the flipped pair
			seen[g[left]^g[right]] = true // XOR of the two resulting independent runs
		}
		// mex = smallest non-negative integer not in seen.
		mex := 0
		for seen[mex] {
			mex++
		}
		g[k] = mex
	}
	nim := 0
	for _, r := range runs {
		nim ^= g[r] // XOR the Grundy numbers of all runs
	}
	return nim != 0 // nonzero nim-sum ⇒ first player wins
}

func main() {
	fmt.Println("=== Approach 1: Backtracking ===")
	fmt.Println(backtracking("++++")) // expected true
	fmt.Println(backtracking("+"))    // expected false

	fmt.Println("=== Approach 2: Backtracking + Memoization ===")
	fmt.Println(memoBacktracking("++++")) // expected true
	fmt.Println(memoBacktracking("+"))    // expected false

	fmt.Println("=== Approach 3: Sprague-Grundy (Optimal) ===")
	fmt.Println(spragueGrundy("++++")) // expected true
	fmt.Println(spragueGrundy("+"))    // expected false
}
