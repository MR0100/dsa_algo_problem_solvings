package main

import "fmt"

// ── Approach 1: Minimax Backtracking (No Memo) ───────────────────────────────
//
// bruteForceMinimax solves Can I Win by simulating the game: the current player
// tries every unused number and wins if any choice reaches the total OR forces
// the opponent into a losing position.
//
// Intuition:
//
//	This is a two-player, perfect-information, zero-sum game, so it is decided
//	by minimax. On your turn you win if EITHER (a) some unused number x makes
//	the running remaining ≤ 0 immediately, OR (b) some x leaves the opponent
//	in a position from which THEY cannot win. Model the "which numbers are
//	still available" state as a bitmask (bit i set ⇒ number i+1 is used).
//	Recurse; the first player controls the root.
//
// Algorithm:
//  1. Handle trivialities: if desiredTotal ≤ 0 the mover has already won (true);
//     if the sum 1..max < desiredTotal, nobody can ever reach it (false).
//  2. canWin(used, remaining): for each unused number x, if x ≥ remaining OR
//     the opponent loses on canWin(used|bit, remaining−x), return true.
//  3. If no choice wins, return false.
//
// Time:  O(n!) worst case — every ordering of picks explored without caching.
// Space: O(n) recursion depth.
func bruteForceMinimax(maxChoosableInteger int, desiredTotal int) bool {
	if desiredTotal <= 0 {
		return true // target already met before anyone moves
	}
	// If even grabbing every number cannot reach the target, first player can't win.
	if maxChoosableInteger*(maxChoosableInteger+1)/2 < desiredTotal {
		return false
	}

	var canWin func(used, remaining int) bool
	canWin = func(used, remaining int) bool {
		for x := 1; x <= maxChoosableInteger; x++ {
			bit := 1 << uint(x-1) // bit representing "number x is taken"
			if used&bit != 0 {
				continue // x already used this game
			}
			// Win now (x meets/exceeds what's left) OR opponent loses after we take x.
			if x >= remaining || !canWin(used|bit, remaining-x) {
				return true
			}
		}
		return false // no move leads to a win → current player loses
	}
	return canWin(0, desiredTotal)
}

// ── Approach 2: Minimax + Bitmask Memoization (Optimal) ──────────────────────
//
// bitmaskMemo solves Can I Win identically to Approach 1 but caches each
// game state so it is evaluated once — the standard efficient solution.
//
// Intuition:
//
//	The outcome of a position depends ONLY on which numbers remain, not on the
//	order they were taken, because `remaining = desiredTotal − Σ(used)` is
//	determined by the used-set. So the state is fully captured by the bitmask
//	`used` alone. There are 2^max distinct masks; memoize each. That collapses
//	the factorial search into O(2^max · max).
//
// Algorithm:
//  1. Same trivial checks as Approach 1.
//  2. memo[used] ∈ {unknown, win, lose}. In canWin(used, remaining): if cached,
//     return it; else try every unused x (win if x ≥ remaining or opponent
//     loses); store and return the result.
//
// Time:  O(2^max · max) — each of 2^max masks does ≤ max work once.
// Space: O(2^max) — the memo table plus O(max) recursion depth.
func bitmaskMemo(maxChoosableInteger int, desiredTotal int) bool {
	if desiredTotal <= 0 {
		return true
	}
	if maxChoosableInteger*(maxChoosableInteger+1)/2 < desiredTotal {
		return false
	}

	// memo maps a used-mask to the current mover's outcome: 0 unknown, 1 win, -1 lose.
	memo := make([]int8, 1<<uint(maxChoosableInteger))

	var canWin func(used, remaining int) bool
	canWin = func(used, remaining int) bool {
		if memo[used] != 0 {
			return memo[used] == 1 // already solved this position
		}
		result := false
		for x := 1; x <= maxChoosableInteger; x++ {
			bit := 1 << uint(x-1)
			if used&bit != 0 {
				continue // x taken
			}
			// Immediate win, or push opponent into a proven loss.
			if x >= remaining || !canWin(used|bit, remaining-x) {
				result = true
				break // one winning move is enough
			}
		}
		if result { // cache the verdict for this mask
			memo[used] = 1
		} else {
			memo[used] = -1
		}
		return result
	}
	return canWin(0, desiredTotal)
}

func main() {
	fmt.Println("=== Approach 1: Minimax Backtracking (No Memo) ===")
	fmt.Printf("max=10, total=11  got=%t  expected false\n", bruteForceMinimax(10, 11))
	fmt.Printf("max=10, total=0   got=%t  expected true\n", bruteForceMinimax(10, 0))
	fmt.Printf("max=10, total=1   got=%t  expected true\n", bruteForceMinimax(10, 1))
	fmt.Printf("max=4,  total=6   got=%t  expected true\n", bruteForceMinimax(4, 6))   // classic: first player picks 1
	fmt.Printf("max=18, total=79  got=%t  expected true\n", bruteForceMinimax(18, 79)) // sum 171 ≥ 79, first player wins

	fmt.Println("=== Approach 2: Minimax + Bitmask Memoization (Optimal) ===")
	fmt.Printf("max=10, total=11  got=%t  expected false\n", bitmaskMemo(10, 11))
	fmt.Printf("max=10, total=0   got=%t  expected true\n", bitmaskMemo(10, 0))
	fmt.Printf("max=10, total=1   got=%t  expected true\n", bitmaskMemo(10, 1))
	fmt.Printf("max=4,  total=6   got=%t  expected true\n", bitmaskMemo(4, 6))
	fmt.Printf("max=18, total=79  got=%t  expected true\n", bitmaskMemo(18, 79))
	fmt.Printf("max=20, total=210 got=%t  expected false\n", bitmaskMemo(20, 210)) // sum=210 forces all 20 picks; the 20th (winning) move is player 2's
}
