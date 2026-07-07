package main

import "fmt"

// ── Approach 1: Recursion (Brute Force) ──────────────────────────────────────
//
// recursion solves Paint Fence by defining a recurrence on the number of ways
// to paint the last post.
//
// Intuition:
//
//	When we paint post i we may either use a NEW color (different from post
//	i-1) or REUSE the color of post i-1. Reusing is only legal if post i-1
//	and i-2 differ (otherwise three in a row share a color). So the count of
//	ways to paint i posts splits into:
//	  - "same as previous": we must have differed one step earlier → ways(i-2)
//	    times 1 legal same-color choice.
//	  - "different from previous": (k-1) new-color choices → ways(i-1)*(k-1).
//	This yields total(i) = (k-1) * (total(i-1) + total(i-2)).
//
// Algorithm:
//  1. Base cases: 0 posts → 0, 1 post → k, 2 posts → k*k.
//  2. Otherwise return (k-1)*(recursion(n-1)+recursion(n-2)).
//
// Time:  O(2^n) — the tree branches into two recursive calls without memoization.
// Space: O(n) — recursion stack depth.
func recursion(n, k int) int {
	if n == 0 { // no posts → no way to paint
		return 0
	}
	if n == 1 { // a single post can be any of the k colors
		return k
	}
	if n == 2 { // two posts: k choices for first, k for second (may match)
		return k * k
	}
	// (k-1) new-color choices multiply BOTH the "reuse" (n-2) and "new" (n-1) cases
	return (k - 1) * (recursion(n-1, k) + recursion(n-2, k))
}

// ── Approach 2: DP Bottom-Up (Table) ─────────────────────────────────────────
//
// dpBottomUp solves Paint Fence by filling a table of subproblem answers so
// each is computed once.
//
// Intuition:
//
//	Same recurrence total(i) = (k-1)*(total(i-1)+total(i-2)), but iterate from
//	the base cases upward, storing results so we never recompute.
//
// Algorithm:
//  1. Handle n==0, n==1, n==2 directly.
//  2. dp[1]=k, dp[2]=k*k.
//  3. For i=3..n: dp[i] = (k-1)*(dp[i-1]+dp[i-2]).
//  4. Return dp[n].
//
// Time:  O(n) — one pass filling the table.
// Space: O(n) — the dp array.
func dpBottomUp(n, k int) int {
	if n == 0 {
		return 0
	}
	if n == 1 {
		return k
	}
	dp := make([]int, n+1) // dp[i] = ways to paint the first i posts
	dp[1] = k              // 1 post → k ways
	dp[2] = k * k          // 2 posts → k*k ways
	for i := 3; i <= n; i++ {
		dp[i] = (k - 1) * (dp[i-1] + dp[i-2]) // fold the recurrence
	}
	return dp[n]
}

// ── Approach 3: DP O(1) Space (Optimal) ──────────────────────────────────────
//
// dpConstantSpace solves Paint Fence keeping only the two most recent states.
//
// Intuition:
//
//	total(i) depends only on total(i-1) and total(i-2), so two rolling
//	variables replace the whole table.
//
// Algorithm:
//  1. prev2 = k (i=1), prev1 = k*k (i=2).
//  2. For i=3..n: cur = (k-1)*(prev1+prev2); slide prev2=prev1, prev1=cur.
//  3. Return prev1.
//
// Time:  O(n) — single loop.
// Space: O(1) — two scalars.
func dpConstantSpace(n, k int) int {
	if n == 0 {
		return 0
	}
	if n == 1 {
		return k
	}
	prev2 := k     // ways for i-2 (starts at i=1)
	prev1 := k * k // ways for i-1 (starts at i=2)
	for i := 3; i <= n; i++ {
		cur := (k - 1) * (prev1 + prev2) // recurrence
		prev2 = prev1                    // slide the window forward
		prev1 = cur
	}
	return prev1 // holds ways for i=n
}

func main() {
	fmt.Println("=== Approach 1: Recursion ===")
	fmt.Println(recursion(3, 2)) // expected 6
	fmt.Println(recursion(1, 1)) // expected 1
	fmt.Println(recursion(7, 2)) // expected 42

	fmt.Println("=== Approach 2: DP Bottom-Up ===")
	fmt.Println(dpBottomUp(3, 2)) // expected 6
	fmt.Println(dpBottomUp(1, 1)) // expected 1
	fmt.Println(dpBottomUp(7, 2)) // expected 42

	fmt.Println("=== Approach 3: DP O(1) Space (Optimal) ===")
	fmt.Println(dpConstantSpace(3, 2)) // expected 6
	fmt.Println(dpConstantSpace(1, 1)) // expected 1
	fmt.Println(dpConstantSpace(7, 2)) // expected 42
}
