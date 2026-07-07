package main

import "fmt"

// ── Approach 1: Brute Force (Recursive Subset Sum) ───────────────────────────
//
// bruteForce solves Partition Equal Subset Sum by asking the equivalent
// question "is there a subset that sums to total/2?" and answering it with
// naive include/exclude recursion over every element.
//
// Intuition:
//
//	Splitting nums into two equal-sum halves means each half must sum to
//	total/2. So the whole problem reduces to the classic subset-sum decision:
//	can we pick a subset of nums whose sum is exactly target = total/2? If
//	total is odd no such split exists. Otherwise try every element two ways —
//	take it (subtract from the remaining target) or skip it — and see if any
//	combination lands on exactly 0.
//
// Algorithm:
//  1. sum all numbers; if odd, return false immediately.
//  2. target = sum/2.
//  3. Recurse from index 0 with remaining = target:
//     - remaining == 0 → found a valid subset → true.
//     - index past end or remaining < 0 → dead end → false.
//     - otherwise: take nums[index] OR skip it.
//
// Time:  O(2^n) — each element branches into take/skip.
// Space: O(n) — recursion depth equals the number of elements.
func bruteForce(nums []int) bool {
	sum := 0
	for _, v := range nums { // total of all elements
		sum += v
	}
	if sum%2 != 0 { // an odd total can never split into two equal integer halves
		return false
	}
	target := sum / 2 // each side must reach exactly this

	var canReach func(index, remaining int) bool
	canReach = func(index, remaining int) bool {
		if remaining == 0 {
			return true // exact subset found
		}
		if index >= len(nums) || remaining < 0 {
			return false // ran out of items or overshot the target
		}
		// Branch: use nums[index] toward the target, OR leave it out.
		return canReach(index+1, remaining-nums[index]) ||
			canReach(index+1, remaining)
	}
	return canReach(0, target)
}

// ── Approach 2: Top-Down DP (Memoised Subset Sum) ────────────────────────────
//
// dpTopDown solves Partition Equal Subset Sum by memoising the brute-force
// recursion on (index, remaining), collapsing the exponential tree into a
// polynomial number of distinct states.
//
// Intuition:
//
//	The recursion revisits the same (index, remaining) pair through many
//	different take/skip paths. Since the answer for a state never changes,
//	cache it. There are only n × (target+1) distinct states, so filling each
//	once turns O(2^n) into O(n·target).
//
// Algorithm:
//  1. Same reduction to target = sum/2 (odd sum → false).
//  2. memo[index][remaining] ∈ {unset, false, true}; recurse as before but
//     read/write the cache for each (index, remaining).
//
// Time:  O(n·target) — every (index, remaining) state solved once.
// Space: O(n·target) memo + O(n) recursion stack.
func dpTopDown(nums []int) bool {
	sum := 0
	for _, v := range nums {
		sum += v
	}
	if sum%2 != 0 {
		return false
	}
	target := sum / 2

	// memo[i][r]: 0 = unvisited, 1 = false, 2 = true (avoids a separate seen map).
	memo := make([][]int8, len(nums))
	for i := range memo {
		memo[i] = make([]int8, target+1)
	}

	var canReach func(index, remaining int) bool
	canReach = func(index, remaining int) bool {
		if remaining == 0 {
			return true
		}
		if index >= len(nums) || remaining < 0 {
			return false
		}
		if memo[index][remaining] != 0 { // seen this exact state before
			return memo[index][remaining] == 2
		}
		res := canReach(index+1, remaining-nums[index]) || // take
			canReach(index+1, remaining) // skip
		if res {
			memo[index][remaining] = 2 // cache true
		} else {
			memo[index][remaining] = 1 // cache false
		}
		return res
	}
	return canReach(0, target)
}

// ── Approach 3: Bottom-Up 2D DP ──────────────────────────────────────────────
//
// dp2D solves Partition Equal Subset Sum with a classic 0/1-knapsack table:
// dp[i][s] = "using the first i numbers, can we make sum s?"
//
// Intuition:
//
//	Build subset-sum reachability iteratively. With the first i items you can
//	make sum s if either the first i-1 items already make s (skip item i), or
//	they make s - nums[i-1] (take item i). Sum s == 0 is always reachable
//	(empty subset). The answer is dp[n][target].
//
// Algorithm:
//  1. Reduce to target = sum/2 (odd → false).
//  2. dp[i][0] = true for all i (empty subset sums to 0).
//  3. For each item i and each sum s: dp[i][s] = dp[i-1][s] OR
//     (s >= nums[i-1] AND dp[i-1][s-nums[i-1]]).
//  4. Return dp[n][target].
//
// Time:  O(n·target) — fill the whole table once.
// Space: O(n·target) — the full 2D table.
func dp2D(nums []int) bool {
	sum := 0
	for _, v := range nums {
		sum += v
	}
	if sum%2 != 0 {
		return false
	}
	target := sum / 2
	n := len(nums)

	// dp[i][s] = can the first i numbers form subset-sum s.
	dp := make([][]bool, n+1)
	for i := range dp {
		dp[i] = make([]bool, target+1)
		dp[i][0] = true // sum 0 always achievable via the empty subset
	}

	for i := 1; i <= n; i++ {
		v := nums[i-1] // the i-th number (0-indexed)
		for s := 1; s <= target; s++ {
			dp[i][s] = dp[i-1][s] // case 1: don't use v
			if s >= v {
				// case 2: use v, provided the remainder s-v was reachable before
				dp[i][s] = dp[i][s] || dp[i-1][s-v]
			}
		}
	}
	return dp[n][target]
}

// ── Approach 4: Bottom-Up 1D DP (Optimal) ────────────────────────────────────
//
// dp1D solves Partition Equal Subset Sum by compressing the knapsack table to
// a single boolean row, iterating each sum from high to low.
//
// Intuition:
//
//	Row i of the 2D table only reads row i-1, so one array suffices. The catch:
//	each item may be used at most once, so when folding item v in we must scan
//	sums from target down to v. Descending guarantees dp[s-v] still refers to
//	the previous item's state (item not yet used at a lower sum this round);
//	ascending would let one item be counted multiple times.
//
// Algorithm:
//  1. Reduce to target = sum/2 (odd → false).
//  2. dp[0] = true.
//  3. For each v: for s from target down to v: dp[s] = dp[s] || dp[s-v].
//     Early exit once dp[target] is true.
//  4. Return dp[target].
//
// Time:  O(n·target) — same states, one row.
// Space: O(target) — a single boolean array.
func dp1D(nums []int) bool {
	sum := 0
	for _, v := range nums {
		sum += v
	}
	if sum%2 != 0 {
		return false
	}
	target := sum / 2

	dp := make([]bool, target+1) // dp[s] = some processed subset sums to s
	dp[0] = true                 // empty subset

	for _, v := range nums {
		// Descend so each v folds in at most once (0/1-knapsack requirement).
		for s := target; s >= v; s-- {
			if dp[s-v] { // s-v was reachable without v → s reachable with v
				dp[s] = true
			}
		}
		if dp[target] {
			return true // target already reachable — no need to process more items
		}
	}
	return dp[target]
}

func main() {
	fmt.Println("=== Approach 1: Brute Force (Recursive Subset Sum) ===")
	fmt.Println(bruteForce([]int{1, 5, 11, 5})) // expected true
	fmt.Println(bruteForce([]int{1, 2, 3, 5}))  // expected false

	fmt.Println("=== Approach 2: Top-Down DP (Memoised) ===")
	fmt.Println(dpTopDown([]int{1, 5, 11, 5})) // expected true
	fmt.Println(dpTopDown([]int{1, 2, 3, 5}))  // expected false

	fmt.Println("=== Approach 3: Bottom-Up 2D DP ===")
	fmt.Println(dp2D([]int{1, 5, 11, 5})) // expected true
	fmt.Println(dp2D([]int{1, 2, 3, 5}))  // expected false

	fmt.Println("=== Approach 4: Bottom-Up 1D DP (Optimal) ===")
	fmt.Println(dp1D([]int{1, 5, 11, 5})) // expected true
	fmt.Println(dp1D([]int{1, 2, 3, 5}))  // expected false
}
