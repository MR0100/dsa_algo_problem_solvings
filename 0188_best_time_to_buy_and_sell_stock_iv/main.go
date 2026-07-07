package main

import "fmt"

// ── Approach 1: Brute Force (Exhaustive Recursion) ───────────────────────────
//
// bruteForce solves Best Time to Buy and Sell Stock IV by trying every
// possible sequence of buy/sell/skip decisions day by day.
//
// Intuition:
//
//	On any day there are at most two sensible moves: do nothing, or trade
//	(buy if we hold nothing, sell if we hold a share). Exploring both choices
//	at every day enumerates every legal trading schedule with at most k
//	completed transactions; the best leaf value is the answer. A transaction
//	is counted when it completes (on the sell), so `remaining` only drops then.
//
// Algorithm:
//  1. Recurse on state (day, remaining, holding).
//  2. Base case: past the last day or no transactions remaining → profit 0.
//  3. Otherwise take the max of: skip today; sell today (+prices[day],
//     remaining-1) if holding; buy today (-prices[day]) if not holding.
//
// Time:  O(2^n) — two branches per day with no state reuse.
// Space: O(n) — recursion depth is one frame per day.
func bruteForce(k int, prices []int) int {
	var solve func(day, remaining int, holding bool) int
	solve = func(day, remaining int, holding bool) int {
		// no days left, or no transactions left to earn with → nothing more to gain
		if day == len(prices) || remaining == 0 {
			return 0
		}
		best := solve(day+1, remaining, holding) // choice 1: do nothing today
		if holding {
			// choice 2a: sell today — realise the price, one transaction completes
			if p := prices[day] + solve(day+1, remaining-1, false); p > best {
				best = p
			}
		} else {
			// choice 2b: buy today — pay the price, transaction counted on the future sell
			if p := -prices[day] + solve(day+1, remaining, true); p > best {
				best = p
			}
		}
		return best
	}
	return solve(0, k, false) // start on day 0, k transactions available, holding nothing
}

// ── Approach 2: DP Top-Down (Memoized Recursion) ─────────────────────────────
//
// dpTopDown solves Best Time to Buy and Sell Stock IV with the same recursion
// as the brute force, but caches every (day, remaining, holding) state.
//
// Intuition:
//
//	The brute force recomputes identical futures: how we *reached* day d with
//	t transactions left and a given holding flag never changes what the best
//	continuation is worth. There are only n·(k+1)·2 distinct states, so
//	memoizing collapses the exponential tree to one visit per state.
//
// Algorithm:
//  1. memo[day][remaining][holding] caches the best future profit; -1 marks
//     "not computed" (safe because doing nothing forever guarantees >= 0).
//  2. Recurse exactly like the brute force, consulting/filling the cache.
//
// Time:  O(n·k) — n·(k+1)·2 states, O(1) transition work each.
// Space: O(n·k) — the memo table plus O(n) recursion stack.
func dpTopDown(k int, prices []int) int {
	n := len(prices)
	if n == 0 || k == 0 {
		return 0 // no days or no allowed transactions → no profit possible
	}
	// memo[day][remaining][holding]; -1 = unset (all real answers are >= 0)
	memo := make([][][2]int, n)
	for d := range memo {
		memo[d] = make([][2]int, k+1)
		for t := range memo[d] {
			memo[d][t] = [2]int{-1, -1}
		}
	}
	var solve func(day, remaining, holding int) int
	solve = func(day, remaining, holding int) int {
		if day == n || remaining == 0 {
			return 0 // out of days or out of transactions
		}
		if memo[day][remaining][holding] != -1 {
			return memo[day][remaining][holding] // state already solved once
		}
		best := solve(day+1, remaining, holding) // skip today
		if holding == 1 {
			// sell: pocket today's price, the transaction completes
			if p := prices[day] + solve(day+1, remaining-1, 0); p > best {
				best = p
			}
		} else {
			// buy: spend today's price, stay on the same transaction budget
			if p := -prices[day] + solve(day+1, remaining, 1); p > best {
				best = p
			}
		}
		memo[day][remaining][holding] = best // cache before returning
		return best
	}
	return solve(0, k, 0)
}

// ── Approach 3: DP Bottom-Up (Transactions × Days + maxDiff Trick) ───────────
//
// dpBottomUp solves Best Time to Buy and Sell Stock IV with a 2D table
// dp[t][d] = best profit using at most t transactions within days 0..d.
//
// Intuition:
//
//	On day d with t transactions we either don't trade (dp[t][d-1]) or we
//	sell on day d after buying on some earlier day m, which is worth
//	prices[d] - prices[m] + dp[t-1][m]. Scanning all m each time costs
//	O(n^2·k); instead carry maxDiff = max over m<d of (dp[t-1][m] - prices[m])
//	forward incrementally, so the best "buy day" for any future sell day is
//	always available in O(1).
//
// Algorithm:
//  1. dp has k+1 rows and n columns, row 0 and column 0 stay 0.
//  2. For each t = 1..k: start maxDiff = dp[t-1][0] - prices[0]; for each
//     d = 1..n-1: dp[t][d] = max(dp[t][d-1], prices[d]+maxDiff), then fold in
//     maxDiff = max(maxDiff, dp[t-1][d] - prices[d]).
//  3. Answer is dp[k][n-1].
//
// Time:  O(n·k) — one O(n) sweep per transaction row thanks to maxDiff.
// Space: O(n·k) — the full DP table (kept 2D for clarity; see Approach 4 for O(k)).
func dpBottomUp(k int, prices []int) int {
	n := len(prices)
	if n < 2 || k == 0 {
		return 0 // need at least two days (a buy and a sell) to profit
	}
	// dp[t][d] = best profit with at most t transactions over days 0..d
	dp := make([][]int, k+1)
	for t := range dp {
		dp[t] = make([]int, n) // row t=0 and column d=0 remain 0 (no trades possible)
	}
	for t := 1; t <= k; t++ {
		// maxDiff = best (dp[t-1][m] - prices[m]) over buy days m seen so far
		maxDiff := dp[t-1][0] - prices[0]
		for d := 1; d < n; d++ {
			// don't trade today, or sell today on top of the best earlier buy
			dp[t][d] = max(dp[t][d-1], prices[d]+maxDiff)
			// let day d compete as the buy day for future sell days
			maxDiff = max(maxDiff, dp[t-1][d]-prices[d])
		}
	}
	return dp[k][n-1] // at most k transactions over all n days
}

// ── Approach 4: State Machine, Space Optimized (Optimal) ─────────────────────
//
// stateMachine solves Best Time to Buy and Sell Stock IV by tracking, for
// each transaction slot j, the best cash after its buy and after its sell.
//
// Intuition:
//
//	A schedule with at most k transactions walks through 2k+1 states:
//	start → bought#1 → sold#1 → ... → bought#k → sold#k. For each state keep
//	the best cash balance achievable so far. Each new price relaxes every
//	state: buy[j] improves by buying today out of sell[j-1]; sell[j] improves
//	by selling today out of buy[j]. Extra insight: one transaction consumes
//	at least two days, so when k >= n/2 the cap never binds and the problem
//	degenerates to "unlimited transactions" — solved greedily by pocketing
//	every upward price step (this also guards the DP against huge k).
//
// Algorithm:
//  1. If k >= n/2: return the sum of all positive day-to-day price rises.
//  2. Otherwise keep arrays buy[1..k] (init -inf) and sell[0..k] (init 0).
//  3. For every price p and every slot j: buy[j] = max(buy[j], sell[j-1]-p),
//     then sell[j] = max(sell[j], buy[j]+p).
//  4. Return sell[k].
//
// Time:  O(n·k) — k relaxations per price (or O(n) via greedy when k >= n/2).
// Space: O(k) — two arrays of k+1 running states; the n dimension is gone.
func stateMachine(k int, prices []int) int {
	n := len(prices)
	if n < 2 || k == 0 {
		return 0 // can't complete any transaction
	}
	// A completed transaction needs a buy day and a later sell day (>= 2 days),
	// so at most n/2 transactions ever fit: for k >= n/2 the cap is irrelevant.
	if k >= n/2 {
		profit := 0
		for d := 1; d < n; d++ {
			if prices[d] > prices[d-1] {
				profit += prices[d] - prices[d-1] // harvest every ascent (unlimited trades)
			}
		}
		return profit
	}
	const negInf = -1 << 60 // effectively -infinity: "state not reachable yet"
	buy := make([]int, k+1)  // buy[j]  = best cash right after the j-th buy
	sell := make([]int, k+1) // sell[j] = best cash right after the j-th sell (sell[0]=0: no trades)
	for j := 1; j <= k; j++ {
		buy[j] = negInf // no buy has happened yet
	}
	for _, p := range prices {
		for j := 1; j <= k; j++ {
			// start transaction j today: cash after (j-1)-th sell minus today's price
			buy[j] = max(buy[j], sell[j-1]-p)
			// close transaction j today: cash after its buy plus today's price
			sell[j] = max(sell[j], buy[j]+p)
		}
	}
	return sell[k] // best cash after at most k completed transactions
}

func main() {
	// Example 1: k = 2, prices = [2,4,1]        → 2
	// Example 2: k = 2, prices = [3,2,6,5,0,3]  → 7
	k1, p1 := 2, []int{2, 4, 1}
	k2, p2 := 2, []int{3, 2, 6, 5, 0, 3}

	fmt.Println("=== Approach 1: Brute Force (Exhaustive Recursion) ===")
	fmt.Println(bruteForce(k1, p1)) // expected: 2
	fmt.Println(bruteForce(k2, p2)) // expected: 7

	fmt.Println("=== Approach 2: DP Top-Down (Memoized Recursion) ===")
	fmt.Println(dpTopDown(k1, p1)) // expected: 2
	fmt.Println(dpTopDown(k2, p2)) // expected: 7

	fmt.Println("=== Approach 3: DP Bottom-Up (Transactions x Days + maxDiff) ===")
	fmt.Println(dpBottomUp(k1, p1)) // expected: 2
	fmt.Println(dpBottomUp(k2, p2)) // expected: 7

	fmt.Println("=== Approach 4: State Machine, Space Optimized (Optimal) ===")
	fmt.Println(stateMachine(k1, p1)) // expected: 2
	fmt.Println(stateMachine(k2, p2)) // expected: 7
}
