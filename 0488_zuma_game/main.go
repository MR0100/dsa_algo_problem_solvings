package main

import (
	"fmt"
	"sort"
	"strings"
)

// Zuma Game — find the MINIMUM number of balls we must insert from the hand to
// clear the whole board. Inserting a ball can trigger a cascade: any group of
// 3+ same-coloured balls that forms is removed, which may cause further groups
// to touch and be removed, and so on. Return -1 if the board can never be cleared.
//
// Both solutions share two helpers:
//   - collapse: repeatedly delete 3+ runs until the board is stable.
//   - sortHand: a canonical hand string, so equal multisets share a memo key.

// collapse removes every maximal run of length ≥ 3 from the board, repeating
// until no such run remains, and returns the stabilised board.
//
// Why a loop: deleting one run can make its two neighbours (previously
// separated) become adjacent and form a new ≥ 3 run — the classic Zuma cascade.
func collapse(board string) string {
	// Keep scanning until a full pass removes nothing.
	for {
		n := len(board)
		removed := false
		i := 0
		var b strings.Builder
		for i < n {
			j := i
			// Extend j over the maximal run of the same colour starting at i.
			for j < n && board[j] == board[i] {
				j++
			}
			if j-i >= 3 {
				removed = true // this run of length ≥ 3 disappears
				// (append nothing — the run is deleted)
			} else {
				b.WriteString(board[i:j]) // keep short runs verbatim
			}
			i = j // jump to the next colour
		}
		board = b.String()
		if !removed {
			return board // stable: a whole pass deleted nothing
		}
	}
}

// sortHand returns the hand's letters in sorted order, so that hands with the
// same multiset of colours (e.g. "RB" and "BR") produce the same memo key.
func sortHand(hand string) string {
	bs := []byte(hand)
	sort.Slice(bs, func(a, b int) bool { return bs[a] < bs[b] })
	return string(bs)
}

// ── Approach 1: Brute-Force DFS (try every insertion) ─────────────────────────
//
// bruteForceDFS solves Zuma Game by recursively trying to insert each distinct
// hand ball at every gap of the board, collapsing after each insertion, and
// taking the minimum insertion count that clears the board.
//
// Intuition:
//
//	At any state (board, hand) we must eventually place some hand ball
//	somewhere. So branch on "which colour do we insert, and at which position",
//	recurse on the resulting collapsed board with that ball removed from the
//	hand, and keep the cheapest solution (+1 for the ball we just used). An
//	empty board costs 0; running out of hand without clearing costs ∞ (-1).
//
// Optimisations that keep it correct AND fast enough (board ≤ 16, hand ≤ 5):
//   - Only insert BEFORE index i when board[i] starts a new colour, and only a
//     colour equal to board[i] (matching an existing run) OR — as a fallback —
//     we still allow inserting the colour at that seam; matching insor­tions are
//     the only ones that can ever trigger a removal, so non-matching inserts are
//     never better and are skipped.
//
// Time:  O(exponential) — bounded by hand permutations × board positions;
//
//	tractable only because |hand| ≤ 5.
//
// Space: O(|hand|) recursion depth (plus transient board strings).
func bruteForceDFS(board, hand string) int {
	const INF = 1 << 30
	var dfs func(board, hand string) int
	dfs = func(board, hand string) int {
		if len(board) == 0 {
			return 0 // board cleared — no more insertions needed
		}
		if len(hand) == 0 {
			return INF // out of balls but board not empty — impossible from here
		}
		best := INF
		// Try inserting one hand ball; consider each distinct colour once.
		for h := 0; h < len(hand); h++ {
			if h > 0 && hand[h] == hand[h-1] {
				continue // skip duplicate colours (hand is sorted) — same result
			}
			ball := hand[h]
			// Remaining hand with this one ball removed.
			rest := hand[:h] + hand[h+1:]
			// Try every gap 0..len(board); dedupe positions inside equal runs.
			for pos := 0; pos <= len(board); pos++ {
				// Skip a gap that sits strictly inside a run of the same colour:
				// inserting there is equivalent to inserting at the run's start.
				if pos > 0 && pos < len(board) && board[pos] == board[pos-1] {
					continue
				}
				// Heuristic prune: only insert next to a matching colour (that is
				// the only way a removal can ever start). Inserting where the ball
				// touches its own colour, i.e. board[pos]==ball or board[pos-1]==ball.
				matchRight := pos < len(board) && board[pos] == ball
				matchLeft := pos > 0 && board[pos-1] == ball
				if !matchLeft && !matchRight {
					continue // a lone ball with no like neighbour can never help
				}
				// Build the new board and collapse any cascade.
				next := collapse(board[:pos] + string(ball) + board[pos:])
				if r := dfs(next, rest); r != INF && r+1 < best {
					best = r + 1 // +1 for the ball we just inserted
				}
			}
		}
		return best
	}
	res := dfs(collapse(board), sortHand(hand)) // pre-collapse in case input had a run
	if res >= INF {
		return -1 // never cleared
	}
	return res
}

// ── Approach 2: Memoised DFS (Optimal) ────────────────────────────────────────
//
// memoDFS solves Zuma Game with the same insertion search as brute force, but
// caches the best result for each (board, sortedHand) state so identical
// sub-problems reached by different insertion orders are solved once.
//
// Intuition:
//
//	Different insertion sequences frequently converge on the same
//	(board, remaining-hand) pair — e.g. using two interchangeable balls in
//	either order. The minimum extra insertions from a state depends ONLY on
//	that state, so memoise on (board, sorted hand). This turns the branching
//	search into a lookup over the reachable state set.
//
// Algorithm:
//  1. Key each state by board + "|" + sorted(hand).
//  2. dfs(board, hand): if board empty → 0; if hand empty → INF; if cached → return.
//  3. Otherwise branch over (colour, position) exactly as brute force, collapse,
//     recurse, take min+1, store in the cache, and return.
//
// Time:  O(S · |board| · |hand|) where S = number of distinct reachable states;
//
//	each state does O(|board|·|hand|) branching work.
//
// Space: O(S) memo entries + O(|hand|) recursion depth.
func memoDFS(board, hand string) int {
	const INF = 1 << 30
	memo := map[string]int{} // (board|sortedHand) → min insertions to clear
	var dfs func(board, hand string) int
	dfs = func(board, hand string) int {
		if len(board) == 0 {
			return 0 // cleared
		}
		if len(hand) == 0 {
			return INF // stuck: balls exhausted, board remains
		}
		key := board + "|" + hand // hand is already sorted by the caller
		if v, ok := memo[key]; ok {
			return v // this exact sub-problem was already solved
		}
		best := INF
		for h := 0; h < len(hand); h++ {
			if h > 0 && hand[h] == hand[h-1] {
				continue // one representative per colour
			}
			ball := hand[h]
			rest := hand[:h] + hand[h+1:] // hand stays sorted after removal
			for pos := 0; pos <= len(board); pos++ {
				if pos > 0 && pos < len(board) && board[pos] == board[pos-1] {
					continue // interior of a same-colour run: redundant gap
				}
				matchRight := pos < len(board) && board[pos] == ball
				matchLeft := pos > 0 && board[pos-1] == ball
				if !matchLeft && !matchRight {
					continue // insertion cannot trigger any removal
				}
				next := collapse(board[:pos] + string(ball) + board[pos:])
				if r := dfs(next, rest); r != INF && r+1 < best {
					best = r + 1
				}
			}
		}
		memo[key] = best // record the optimum for this state
		return best
	}
	res := dfs(collapse(board), sortHand(hand))
	if res >= INF {
		return -1
	}
	return res
}

func main() {
	fmt.Println("=== Approach 1: Brute-Force DFS ===")
	fmt.Println(bruteForceDFS("WRRBBW", "RB"))       // expected -1
	fmt.Println(bruteForceDFS("WWRRBBWW", "WRBRW"))  // expected 2
	fmt.Println(bruteForceDFS("G", "GGGGG"))         // expected 2
	fmt.Println(bruteForceDFS("RBYYBBRRB", "YRBGB")) // expected 3

	fmt.Println("=== Approach 2: Memoised DFS (Optimal) ===")
	fmt.Println(memoDFS("WRRBBW", "RB"))       // expected -1
	fmt.Println(memoDFS("WWRRBBWW", "WRBRW"))  // expected 2
	fmt.Println(memoDFS("G", "GGGGG"))         // expected 2
	fmt.Println(memoDFS("RBYYBBRRB", "YRBGB")) // expected 3
}
