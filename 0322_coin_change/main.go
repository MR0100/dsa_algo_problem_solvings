package main

import (
	"fmt"
	"sort"
)

// ── Approach 1: Top-Down DP (Memoized Recursion) ─────────────────────────────
//
// dpTopDown solves Coin Change by recursively asking "fewest coins to make
// amount r", caching each sub-answer so every amount is solved once.
//
// Intuition:
//
//	minCoins(r) = 1 + min over each coin c of minCoins(r-c). The base case
//	minCoins(0) = 0. Because many recursion paths revisit the same remaining
//	amount, memoization collapses the exponential tree into O(amount) distinct
//	states.
//
// Algorithm:
//  1. memo[r] caches the answer for amount r; -2 means "not computed yet".
//  2. solve(r): if r == 0 return 0; if r < 0 return -1 (impossible).
//  3. Try every coin, recurse on r-c, keep the smallest (sub+1).
//  4. Store and return; -1 if no coin leads to a solution.
//
// Time:  O(amount * len(coins)) — each of amount+1 states tries every coin once.
// Space: O(amount) — memo table plus recursion depth.
func dpTopDown(coins []int, amount int) int {
	memo := make([]int, amount+1) // memo[r] = fewest coins for amount r
	for i := range memo {
		memo[i] = -2 // sentinel: not yet computed
	}
	var solve func(r int) int
	solve = func(r int) int {
		if r == 0 {
			return 0 // no coins needed for amount 0
		}
		if r < 0 {
			return -1 // overshot: this path is impossible
		}
		if memo[r] != -2 {
			return memo[r] // reuse cached result
		}
		best := -1 // -1 = still no valid way found
		for _, c := range coins {
			sub := solve(r - c) // fewest coins for the remainder
			if sub >= 0 && (best == -1 || sub+1 < best) {
				best = sub + 1 // take this coin plus the sub-solution
			}
		}
		memo[r] = best // cache before returning
		return best
	}
	return solve(amount)
}

// ── Approach 2: Bottom-Up DP (Unbounded Knapsack) (Optimal) ──────────────────
//
// dpBottomUp solves Coin Change by filling a table dp[0..amount] where dp[a] is
// the fewest coins to make amount a, building small amounts before large ones.
//
// Intuition:
//
//	To make amount a we choose a last coin c <= a, leaving a-c already solved.
//	So dp[a] = 1 + min over coins c<=a of dp[a-c]. Initialise dp[0]=0 and every
//	other entry to "infinity"; any entry still infinite at the end is
//	unreachable.
//
// Algorithm:
//  1. dp[0] = 0; dp[a] = amount+1 (infinity sentinel) for a >= 1.
//  2. For a from 1..amount, for each coin c <= a: dp[a] = min(dp[a], dp[a-c]+1).
//  3. Return dp[amount] if it stayed finite, else -1.
//
// Time:  O(amount * len(coins)) — nested loops over amounts and coins.
// Space: O(amount) — the dp table.
func dpBottomUp(coins []int, amount int) int {
	inf := amount + 1              // larger than any real coin count
	dp := make([]int, amount+1)    // dp[a] = fewest coins for amount a
	for a := 1; a <= amount; a++ { // dp[0] stays 0 (Go zero value)
		dp[a] = inf // start unreachable
	}
	for a := 1; a <= amount; a++ {
		for _, c := range coins {
			if c <= a && dp[a-c]+1 < dp[a] { // coin fits and improves the count
				dp[a] = dp[a-c] + 1 // use one coin c on top of dp[a-c]
			}
		}
	}
	if dp[amount] == inf {
		return -1 // amount never became reachable
	}
	return dp[amount]
}

// ── Approach 3: BFS on Amounts ───────────────────────────────────────────────
//
// bfs solves Coin Change by treating amounts as graph nodes: from amount a there
// is an edge to a-c for each coin c. The shortest path (fewest edges) from
// amount to 0 is the fewest coins.
//
// Intuition:
//
//	BFS explores states in order of distance, so the first time we reach 0 we
//	have used the minimum number of coins. A visited set stops us revisiting an
//	amount already reached at an equal-or-smaller depth.
//
// Algorithm:
//  1. Queue starts with amount at level 0; mark it visited.
//  2. Pop a level: for each node subtract every coin; if we hit 0, return
//     level+1; push each new positive amount not yet visited.
//  3. If the queue drains without reaching 0, return -1.
//
// Time:  O(amount * len(coins)) — each amount is enqueued at most once.
// Space: O(amount) — visited set and queue.
func bfs(coins []int, amount int) int {
	if amount == 0 {
		return 0 // nothing to make
	}
	visited := make([]bool, amount+1) // visited[a] = amount a already reached
	visited[amount] = true
	queue := []int{amount} // BFS frontier of remaining amounts
	level := 0             // coins used to reach the current frontier
	for len(queue) > 0 {
		level++         // we are about to spend one more coin
		next := []int{} // amounts reachable after this coin
		for _, cur := range queue {
			for _, c := range coins {
				rem := cur - c // remaining amount after using coin c
				if rem == 0 {
					return level // reached exactly zero: shortest path
				}
				if rem > 0 && !visited[rem] {
					visited[rem] = true      // do not revisit this amount
					next = append(next, rem) // explore it next level
				}
			}
		}
		queue = next // advance to the next BFS layer
	}
	return -1 // 0 unreachable
}

// ── Approach 4: Greedy + DFS with Pruning ────────────────────────────────────
//
// greedyDFS solves Coin Change by sorting coins descending and trying to use as
// many of the largest coin as possible first, backtracking with strong pruning.
// Greedy alone is NOT correct for arbitrary coin sets, so it is paired with DFS
// that explores fewer large coins when needed — the sort makes good answers
// appear early, letting the bound prune the rest.
//
// Intuition:
//
//	If we already know a solution using `best` coins, any partial path that has
//	spent `count` coins and still cannot possibly beat `best` (even using the
//	largest remaining coin greedily) is abandoned. Descending order finds a tight
//	upper bound fast, so pruning is aggressive.
//
// Algorithm:
//  1. Sort coins descending.
//  2. dfs(idx, remaining, count): if remaining == 0 update best.
//  3. Take k = remaining/coins[idx] copies of the current coin, then recurse on
//     the next coin with fewer copies; prune when count + (ideal remaining
//     coins) >= best.
//
// Time:  O(exponential) worst case, but pruning is fast in practice.
// Space: O(len(coins)) — recursion depth.
func greedyDFS(coins []int, amount int) int {
	sorted := make([]int, len(coins))
	copy(sorted, coins)
	sort.Sort(sort.Reverse(sort.IntSlice(sorted))) // largest coin first
	best := -1                                     // fewest coins found; -1 = none yet
	var dfs func(idx, remaining, count int)
	dfs = func(idx, remaining, count int) {
		if remaining == 0 {
			if best == -1 || count < best {
				best = count // record an improved solution
			}
			return
		}
		if idx == len(sorted) {
			return // ran out of coin types
		}
		coin := sorted[idx]
		maxK := remaining / coin // most copies of this coin that fit
		// Try k from many to few; break early once the optimistic bound
		// (count + k, the fewest coins we could add) can't beat best.
		for k := maxK; k >= 0; k-- {
			if best != -1 && count+k >= best {
				continue // even k coins here can't improve; skip
			}
			dfs(idx+1, remaining-k*coin, count+k) // move to next coin type
		}
	}
	dfs(0, amount, 0)
	return best
}

func main() {
	fmt.Println("=== Approach 1: Top-Down DP (Memoized) ===")
	fmt.Println(dpTopDown([]int{1, 2, 5}, 11)) // expected 3
	fmt.Println(dpTopDown([]int{2}, 3))        // expected -1
	fmt.Println(dpTopDown([]int{1}, 0))        // expected 0

	fmt.Println("=== Approach 2: Bottom-Up DP (Optimal) ===")
	fmt.Println(dpBottomUp([]int{1, 2, 5}, 11)) // expected 3
	fmt.Println(dpBottomUp([]int{2}, 3))        // expected -1
	fmt.Println(dpBottomUp([]int{1}, 0))        // expected 0

	fmt.Println("=== Approach 3: BFS on Amounts ===")
	fmt.Println(bfs([]int{1, 2, 5}, 11)) // expected 3
	fmt.Println(bfs([]int{2}, 3))        // expected -1
	fmt.Println(bfs([]int{1}, 0))        // expected 0

	fmt.Println("=== Approach 4: Greedy + DFS with Pruning ===")
	fmt.Println(greedyDFS([]int{1, 2, 5}, 11)) // expected 3
	fmt.Println(greedyDFS([]int{2}, 3))        // expected -1
	fmt.Println(greedyDFS([]int{1}, 0))        // expected 0
}
