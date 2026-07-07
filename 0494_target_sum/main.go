package main

import "fmt"

// Assign '+' or '-' before each nums[i]; count the sign assignments whose signed
// sum equals target. Key transform: if P = sum of '+' elements and N = sum of
// '-' elements, then P - N = target and P + N = total, so P = (total+target)/2.
// The problem reduces to "count subsets summing to (total+target)/2".

// ── Approach 1: Brute Force DFS (Try Both Signs) ─────────────────────────────
//
// bruteForceDFS explores a binary decision tree: at each index add or subtract
// the current number, counting leaves whose running sum hits target.
//
// Intuition:
//
//	Every expression is a choice of + or - at each of the n positions — 2^n
//	leaves. Walk the tree, carrying the running signed sum. At the end (all
//	numbers consumed) it is a hit iff the sum equals target. No memo, so shared
//	subproblems are recomputed.
//
// Algorithm:
//  1. dfs(i, sum): if i == n, return 1 if sum == target else 0.
//  2. Otherwise return dfs(i+1, sum + nums[i]) + dfs(i+1, sum - nums[i]).
//
// Time:  O(2^n) — two branches per index, n up to 20 → ~10^6, fine here.
// Space: O(n) — recursion depth.
func bruteForceDFS(nums []int, target int) int {
	var dfs func(i, sum int) int
	dfs = func(i, sum int) int {
		if i == len(nums) {
			if sum == target {
				return 1 // one valid signed expression reached target
			}
			return 0
		}
		// Branch on the sign of nums[i]: '+' then '-'.
		return dfs(i+1, sum+nums[i]) + dfs(i+1, sum-nums[i])
	}
	return dfs(0, 0)
}

// ── Approach 2: Top-Down DP with Memoization ─────────────────────────────────
//
// dpTopDown is the DFS above plus a cache keyed on (index, running sum). Many
// different sign prefixes land on the same (i, sum) state, so caching collapses
// the exponential tree into a polynomial number of distinct states.
//
// Intuition:
//
//	The only thing that matters about the choices so far is the current index i
//	and the running sum. Sums range within [-total, +total], so there are at
//	most n * (2*total + 1) distinct states. Memoize each state's answer. We
//	offset the sum by `total` to use it as a non-negative array/map index.
//
// Algorithm:
//  1. memo[i][sum] caches the number of ways from state (i, sum).
//  2. dfs(i, sum): base case as before; otherwise combine both branches and
//     store the result before returning.
//
// Time:  O(n * total) — number of distinct (i, sum) states, each O(1) work.
// Space: O(n * total) — the memо table plus O(n) recursion.
func dpTopDown(nums []int, target int) int {
	total := 0
	for _, v := range nums {
		total += v // maximum absolute reachable sum
	}
	// If target is unreachable in magnitude, no assignment can work.
	if target > total || target < -total {
		return 0
	}
	// memo[i] maps a running sum -> ways; use a map to keep it sparse and to
	// sidestep negative indexing cleanly.
	memo := make([]map[int]int, len(nums))
	for i := range memo {
		memo[i] = map[int]int{}
	}
	var dfs func(i, sum int) int
	dfs = func(i, sum int) int {
		if i == len(nums) {
			if sum == target {
				return 1
			}
			return 0
		}
		if v, ok := memo[i][sum]; ok {
			return v // state already solved
		}
		ways := dfs(i+1, sum+nums[i]) + dfs(i+1, sum-nums[i])
		memo[i][sum] = ways // cache before returning
		return ways
	}
	return dfs(0, 0)
}

// ── Approach 3: Subset-Sum 0/1 Knapsack, 1-D DP (Optimal) ────────────────────
//
// subsetSumDP transforms the problem into "count subsets with sum s", where
// s = (total + target) / 2, then counts them with a classic 1-D knapsack DP.
//
// Intuition:
//
//	Let P be the sum of the '+' group and N the sum of the '-' group. Then
//	P - N = target and P + N = total, so P = (total + target) / 2. The answer is
//	the number of subsets of nums that sum to P. That is the count version of
//	subset-sum: dp[s] = number of ways to reach sum s. Process each number,
//	iterating s DOWNWARD so each item is used at most once (0/1 knapsack).
//
// Algorithm:
//  1. Compute total. If (total + target) is odd or |target| > total, return 0
//     (no integer subset target exists).
//  2. s = (total + target) / 2. Initialize dp[0] = 1 (one way to make sum 0:
//     pick nothing).
//  3. For each num: for j from s down to num: dp[j] += dp[j - num].
//  4. Return dp[s].
//
// Time:  O(n * s) — n numbers times the subset-sum capacity s.
// Space: O(s) — a single 1-D DP array.
func subsetSumDP(nums []int, target int) int {
	total := 0
	for _, v := range nums {
		total += v
	}
	// Need P = (total+target)/2 to be a non-negative integer.
	if target > total || target < -total || (total+target)%2 != 0 {
		return 0
	}
	s := (total + target) / 2 // required '+'-group subset sum
	dp := make([]int, s+1)    // dp[j] = # of subsets summing to j
	dp[0] = 1                 // empty subset makes sum 0 in exactly one way
	for _, num := range nums {
		// Iterate downward so num contributes to each j at most once (0/1).
		for j := s; j >= num; j-- {
			dp[j] += dp[j-num] // ways to reach j using this num
		}
	}
	return dp[s]
}

func main() {
	fmt.Println("=== Approach 1: Brute Force DFS (Try Both Signs) ===")
	fmt.Println(bruteForceDFS([]int{1, 1, 1, 1, 1}, 3)) // expected 5
	fmt.Println(bruteForceDFS([]int{1}, 1))             // expected 1

	fmt.Println("=== Approach 2: Top-Down DP with Memoization ===")
	fmt.Println(dpTopDown([]int{1, 1, 1, 1, 1}, 3)) // expected 5
	fmt.Println(dpTopDown([]int{1}, 1))             // expected 1

	fmt.Println("=== Approach 3: Subset-Sum 0/1 Knapsack, 1-D DP (Optimal) ===")
	fmt.Println(subsetSumDP([]int{1, 1, 1, 1, 1}, 3)) // expected 5
	fmt.Println(subsetSumDP([]int{1}, 1))             // expected 1
}
