package main

import "fmt"

// ── Approach 1: Brute Force (Try Every Burst Order) ───────────────────────────
//
// bruteForce solves Burst Balloons by recursively trying every possible next
// balloon to burst and taking the maximum total coins.
//
// Intuition:
//
//	A naive reading of the problem is a sequence of decisions: which balloon to
//	pop next. At each step, popping balloon i earns nums[left]*nums[i]*nums[right]
//	where left/right are the CURRENT surviving neighbours. Try all n choices,
//	recurse on the shrunken array, and keep the best. This explores n! orders.
//
// Algorithm:
//
//  1. Pad the array with 1 on both ends so boundary balloons have neighbours.
//  2. remaining is the list of still-alive original balloons (their padded
//     indices). For each alive balloon, its neighbours are the alive balloons
//     immediately before/after it in `remaining`.
//  3. Pick each balloon as "burst now", add its coins, recurse on the rest.
//
// Time:  O(n!) — every ordering of bursts is tried. Only feasible for tiny n.
// Space: O(n) recursion depth (plus copies of the remaining list).
func bruteForce(nums []int) int {
	// Pad so index 0 and last are the virtual "1" boundary balloons.
	padded := make([]int, len(nums)+2)
	padded[0], padded[len(padded)-1] = 1, 1
	copy(padded[1:], nums)

	// remaining lists padded indices of balloons not yet burst (1..n).
	remaining := make([]int, len(nums))
	for i := range remaining {
		remaining[i] = i + 1
	}

	var solve func(rem []int) int
	solve = func(rem []int) int {
		if len(rem) == 0 {
			return 0 // nothing left to burst
		}
		best := 0
		for idx := 0; idx < len(rem); idx++ {
			// Current neighbours of rem[idx]: previous alive balloon (or padded
			// boundary 0) and next alive balloon (or padded boundary len-1).
			left := 1
			if idx > 0 {
				left = padded[rem[idx-1]]
			} else {
				left = padded[0]
			}
			right := 1
			if idx < len(rem)-1 {
				right = padded[rem[idx+1]]
			} else {
				right = padded[len(padded)-1]
			}
			coins := left * padded[rem[idx]] * right // coins for bursting rem[idx] now

			// Build the remaining list without index idx.
			next := make([]int, 0, len(rem)-1)
			next = append(next, rem[:idx]...)
			next = append(next, rem[idx+1:]...)

			if got := coins + solve(next); got > best {
				best = got
			}
		}
		return best
	}
	return solve(remaining)
}

// ── Approach 2: Interval DP Top-Down (Memoised "Last Balloon") ────────────────
//
// dpTopDown solves Burst Balloons with the key reframing: instead of asking
// which balloon to burst FIRST, ask which balloon in an open interval is burst
// LAST. Then its neighbours are the fixed interval boundaries, decoupling the
// two subproblems.
//
// Intuition:
//
//	If balloon k is the LAST to burst inside open interval (l, r), then when it
//	pops both boundaries l and r are still present, earning nums[l]*nums[k]*nums[r].
//	Everything strictly between l and k, and between k and r, was already burst
//	independently. So dp(l,r) = max over k in (l,r) of
//	   dp(l,k) + nums[l]*nums[k]*nums[r] + dp(k,r).
//	"Last" makes the boundaries fixed; "first" would leave them changing.
//
// Algorithm:
//
//  1. Pad with 1 on both ends -> array of size n+2, valid balloons at 1..n.
//  2. dp(l,r) = best coins burstable strictly between l and r (exclusive).
//  3. Base case: no balloon between l and r (r == l+1) -> 0.
//  4. Memoise dp(l,r) in a 2-D table.
//
// Time:  O(n^3) — O(n^2) intervals, each scanning O(n) split points.
// Space: O(n^2) memo + O(n) recursion depth.
func dpTopDown(nums []int) int {
	n := len(nums)
	padded := make([]int, n+2)
	padded[0], padded[n+1] = 1, 1
	copy(padded[1:], nums)

	// memo[l][r] caches dp(l,r); -1 marks "not computed yet".
	memo := make([][]int, n+2)
	for i := range memo {
		memo[i] = make([]int, n+2)
		for j := range memo[i] {
			memo[i][j] = -1
		}
	}

	var dp func(l, r int) int
	dp = func(l, r int) int {
		if r-l <= 1 {
			return 0 // no balloon strictly between l and r
		}
		if memo[l][r] != -1 {
			return memo[l][r]
		}
		best := 0
		for k := l + 1; k < r; k++ { // k = last balloon burst in (l, r)
			// When k pops last, boundaries l and r are still alive.
			coins := padded[l]*padded[k]*padded[r] + dp(l, k) + dp(k, r)
			if coins > best {
				best = coins
			}
		}
		memo[l][r] = best
		return best
	}
	return dp(0, n+1)
}

// ── Approach 3: Interval DP Bottom-Up (Optimal, Tabulated) ────────────────────
//
// dpBottomUp solves Burst Balloons with the same "last balloon" recurrence but
// fills the DP table iteratively by increasing interval length.
//
// Intuition:
//
//	dp[l][r] depends only on strictly smaller intervals dp[l][k] and dp[k][r].
//	So if we process intervals from shortest to longest, every dependency is
//	already computed. This removes recursion overhead entirely.
//
// Algorithm:
//
//  1. Pad with 1 on both ends.
//  2. For interval length len from 2..n+1 (distance r-l), for each left l with
//     r = l+len: dp[l][r] = max over k in (l,r) of
//     padded[l]*padded[k]*padded[r] + dp[l][k] + dp[k][r].
//  3. Answer is dp[0][n+1].
//
// Time:  O(n^3) — intervals x split points.
// Space: O(n^2) — the DP table.
func dpBottomUp(nums []int) int {
	n := len(nums)
	padded := make([]int, n+2)
	padded[0], padded[n+1] = 1, 1
	copy(padded[1:], nums)

	// dp[l][r] = max coins from bursting all balloons strictly between l and r.
	dp := make([][]int, n+2)
	for i := range dp {
		dp[i] = make([]int, n+2)
	}

	// length = r - l = distance between the two boundaries.
	for length := 2; length <= n+1; length++ {
		for l := 0; l+length <= n+1; l++ {
			r := l + length
			for k := l + 1; k < r; k++ { // k burst last inside (l, r)
				coins := padded[l]*padded[k]*padded[r] + dp[l][k] + dp[k][r]
				if coins > dp[l][r] {
					dp[l][r] = coins
				}
			}
		}
	}
	return dp[0][n+1]
}

func main() {
	// Official Example 1
	e1 := []int{3, 1, 5, 8}
	fmt.Println("=== Approach 1: Brute Force ===")
	fmt.Println(bruteForce(e1)) // expected 167
	fmt.Println("=== Approach 2: Interval DP Top-Down ===")
	fmt.Println(dpTopDown(e1)) // expected 167
	fmt.Println("=== Approach 3: Interval DP Bottom-Up (Optimal) ===")
	fmt.Println(dpBottomUp(e1)) // expected 167

	// Official Example 2
	e2 := []int{1, 5}
	fmt.Println("=== Approach 1: Brute Force (Example 2) ===")
	fmt.Println(bruteForce(e2)) // expected 10
	fmt.Println("=== Approach 2: Interval DP Top-Down (Example 2) ===")
	fmt.Println(dpTopDown(e2)) // expected 10
	fmt.Println("=== Approach 3: Interval DP Bottom-Up (Example 2) ===")
	fmt.Println(dpBottomUp(e2)) // expected 10
}
