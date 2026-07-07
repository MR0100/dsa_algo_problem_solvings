package main

import "fmt"

// max returns the larger of two ints (Go's builtin min/max exist since 1.21,
// but we keep a helper for clarity in the DP recurrences).
func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

// ── Approach 1: Top-Down DP (Memoized State Machine) ─────────────────────────
//
// dpTopDown solves Buy/Sell Stock with Cooldown by recursing over (day, holding)
// and memoizing, where `holding` says whether we currently own a share.
//
// Intuition:
//
//	At each day, in each state (holding a share or not), we choose the better of
//	acting or waiting. If not holding: skip, or buy (which leaves us holding). If
//	holding: skip, or sell (which forces a cooldown, so the next available buy is
//	day+2). Overlapping subproblems on (day, holding) → memoize.
//
// Algorithm:
//
//	solve(day, holding):
//	  if day >= n: return 0
//	  skip = solve(day+1, holding)
//	  if holding: act = prices[day] + solve(day+2, false)  // sell, cooldown
//	  else:       act = -prices[day] + solve(day+1, true)   // buy
//	  return max(skip, act)
//	Answer = solve(0, false).
//
// Time:  O(n) — 2n states, each O(1).
// Space: O(n) — memo table + recursion stack.
func dpTopDown(prices []int) int {
	n := len(prices)
	if n == 0 {
		return 0
	}
	// memo[day][holding]; holding 0 = not holding, 1 = holding. -1 = uncomputed.
	memo := make([][2]int, n)
	seen := make([][2]bool, n)
	var solve func(day, holding int) int
	solve = func(day, holding int) int {
		if day >= n {
			return 0 // no days left → no more profit
		}
		if seen[day][holding] {
			return memo[day][holding] // reuse computed state
		}
		skip := solve(day+1, holding) // do nothing today
		var act int
		if holding == 1 {
			// Sell today: gain price, then cooldown skips day+1 → jump to day+2.
			act = prices[day] + solve(day+2, 0)
		} else {
			// Buy today: pay price, now holding on the next day.
			act = -prices[day] + solve(day+1, 1)
		}
		best := max(skip, act)
		memo[day][holding] = best
		seen[day][holding] = true
		return best
	}
	return solve(0, 0) // start not holding, on day 0
}

// ── Approach 2: Bottom-Up DP with Three State Arrays ─────────────────────────
//
// dpBottomUp solves the problem with three running states per day:
//
//	hold[i] = max profit on day i while holding a share,
//	sold[i] = max profit on day i having just sold (cooldown starts),
//	rest[i] = max profit on day i idle (not holding, free to buy).
//
// Intuition:
//
//	Model the days as a state machine. Transitions:
//	  hold[i] = max(hold[i-1], rest[i-1] - price)  // keep holding, or buy after rest
//	  sold[i] = hold[i-1] + price                  // sell what we held
//	  rest[i] = max(rest[i-1], sold[i-1])          // stay idle, or emerge from cooldown
//	The cooldown is baked in: you can only buy from `rest`, and you reach `rest`
//	only one day after `sold`. Answer = max(sold[n-1], rest[n-1]).
//
// Algorithm:
//
//	Initialize hold = -prices[0], sold = 0 (impossible → 0 works as floor via
//	transitions), rest = 0. Iterate days updating the three values.
//
// Time:  O(n).
// Space: O(n) if arrays kept; O(1) rolling (see approach 3).
func dpBottomUp(prices []int) int {
	n := len(prices)
	if n == 0 {
		return 0
	}
	hold := make([]int, n)
	sold := make([]int, n)
	rest := make([]int, n)
	hold[0] = -prices[0] // buy on day 0 to be holding
	sold[0] = 0          // can't have sold yet meaningfully
	rest[0] = 0          // idle with no transactions
	for i := 1; i < n; i++ {
		// Keep holding, or buy today coming from a rest day.
		hold[i] = max(hold[i-1], rest[i-1]-prices[i])
		// Sell today the share we were holding yesterday.
		sold[i] = hold[i-1] + prices[i]
		// Stay idle, or finish yesterday's cooldown (sold[i-1]).
		rest[i] = max(rest[i-1], sold[i-1])
	}
	// End either just-sold or resting; holding at the end is never optimal.
	return max(sold[n-1], rest[n-1])
}

// ── Approach 3: O(1) Space Rolling State Machine (Optimal) ────────────────────
//
// dpConstantSpace solves the problem with the same three-state recurrence but
// keeps only the previous day's values, dropping memory to O(1).
//
// Intuition:
//
//	Each transition depends only on day i-1, so three scalars replace the arrays.
//	This is the tightest form of the classic "cooldown state machine" and the
//	answer to give once you have derived the transitions.
//
// Algorithm:
//
//	hold = -prices[0]; sold = 0; rest = 0.
//	For each subsequent price:
//	  prevSold = sold
//	  sold = hold + price
//	  hold = max(hold, rest - price)
//	  rest = max(rest, prevSold)
//	Answer = max(sold, rest).
//
// Time:  O(n).
// Space: O(1).
func dpConstantSpace(prices []int) int {
	n := len(prices)
	if n == 0 {
		return 0
	}
	hold := -prices[0] // holding after buying day 0
	sold := 0          // just sold
	rest := 0          // idle
	for i := 1; i < n; i++ {
		prevSold := sold                 // remember yesterday's sold before overwriting
		sold = hold + prices[i]          // sell today what we held yesterday
		hold = max(hold, rest-prices[i]) // hold on, or buy from a rest day
		rest = max(rest, prevSold)       // stay idle, or exit yesterday's cooldown
	}
	return max(sold, rest) // best ending idle or just-sold
}

func main() {
	fmt.Println("=== Approach 1: Top-Down DP (Memoized) ===")
	fmt.Println(dpTopDown([]int{1, 2, 3, 0, 2})) // expected 3
	fmt.Println(dpTopDown([]int{1}))             // expected 0

	fmt.Println("=== Approach 2: Bottom-Up DP (State Arrays) ===")
	fmt.Println(dpBottomUp([]int{1, 2, 3, 0, 2})) // expected 3
	fmt.Println(dpBottomUp([]int{1}))             // expected 0

	fmt.Println("=== Approach 3: O(1) Space State Machine (Optimal) ===")
	fmt.Println(dpConstantSpace([]int{1, 2, 3, 0, 2})) // expected 3
	fmt.Println(dpConstantSpace([]int{1}))             // expected 0
}
