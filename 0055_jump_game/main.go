package main

import "fmt"

// ── Approach 1: Brute Force (Recursive with Memoization) ─────────────────────
//
// memoization solves Jump Game using top-down DP.
//
// Intuition:
//   From each index, try all reachable next indices. A position is "good" if
//   we can reach the last index from it. Memoize to avoid recomputation.
//   0=unknown, 1=good, 2=bad.
//
// Algorithm:
//   memo[n-1] = good
//   canJump(i):
//     if memo[i] != unknown: return memo[i] == good
//     max_j = min(i + nums[i], n-1)
//     for j = i+1 to max_j:
//       if canJump(j): memo[i]=good; return true
//     memo[i]=bad; return false
//
// Time:  O(n²) — each index computed once, inner loop up to n.
// Space: O(n)  — memo array + recursion stack.
func memoization(nums []int) bool {
	n := len(nums)
	memo := make([]int, n) // 0=unknown, 1=good, 2=bad
	memo[n-1] = 1          // last index is always good

	var canJump func(i int) bool
	canJump = func(i int) bool {
		if memo[i] != 0 {
			return memo[i] == 1
		}
		maxReach := i + nums[i]
		if maxReach >= n-1 {
			memo[i] = 1
			return true
		}
		for j := i + 1; j <= maxReach; j++ {
			if canJump(j) {
				memo[i] = 1
				return true
			}
		}
		memo[i] = 2
		return false
	}

	return canJump(0)
}

// ── Approach 2: DP Bottom-Up ──────────────────────────────────────────────────
//
// dpBottomUp solves Jump Game iterating right to left.
//
// Intuition:
//   dp[i] = true if we can reach the last index from index i.
//   Start from i = n-2 down to 0. dp[i] = true if any j in [i+1, i+nums[i]]
//   has dp[j] = true.
//
// Time:  O(n²)
// Space: O(n)
func dpBottomUp(nums []int) bool {
	n := len(nums)
	dp := make([]bool, n)
	dp[n-1] = true // last index is reachable from itself

	for i := n - 2; i >= 0; i-- {
		maxReach := i + nums[i]
		if maxReach >= n-1 {
			dp[i] = true // can jump directly to or past the end
			continue
		}
		for j := i + 1; j <= maxReach; j++ {
			if dp[j] {
				dp[i] = true
				break
			}
		}
	}

	return dp[0]
}

// ── Approach 3: Greedy (Optimal) ─────────────────────────────────────────────
//
// greedy solves Jump Game in O(n) time and O(1) space.
//
// Intuition:
//   Track the farthest index reachable so far. At each index i, update
//   farthest = max(farthest, i + nums[i]). If i ever exceeds farthest,
//   we can never reach index i — return false. If farthest >= n-1, return true.
//
// Algorithm:
//   farthest = 0
//   for i = 0 to n-1:
//     if i > farthest: return false   // i is unreachable
//     farthest = max(farthest, i + nums[i])
//   return true
//
// Time:  O(n) — single pass.
// Space: O(1)
func greedy(nums []int) bool {
	farthest := 0
	for i, num := range nums {
		if i > farthest {
			return false // can't reach index i
		}
		if i+num > farthest {
			farthest = i + num // extend reach
		}
		if farthest >= len(nums)-1 {
			return true // can reach or pass the last index
		}
	}
	return true
}

func main() {
	fmt.Println("=== Approach 1: Memoization ===")
	fmt.Printf("nums=[2,3,1,1,4]  got=%v  expected true\n", memoization([]int{2, 3, 1, 1, 4}))
	fmt.Printf("nums=[3,2,1,0,4]  got=%v  expected false\n", memoization([]int{3, 2, 1, 0, 4}))
	fmt.Printf("nums=[0]          got=%v  expected true\n", memoization([]int{0}))
	fmt.Printf("nums=[2,0,0]      got=%v  expected true\n", memoization([]int{2, 0, 0}))

	fmt.Println("=== Approach 2: DP Bottom-Up ===")
	fmt.Printf("nums=[2,3,1,1,4]  got=%v  expected true\n", dpBottomUp([]int{2, 3, 1, 1, 4}))
	fmt.Printf("nums=[3,2,1,0,4]  got=%v  expected false\n", dpBottomUp([]int{3, 2, 1, 0, 4}))
	fmt.Printf("nums=[0]          got=%v  expected true\n", dpBottomUp([]int{0}))
	fmt.Printf("nums=[2,0,0]      got=%v  expected true\n", dpBottomUp([]int{2, 0, 0}))

	fmt.Println("=== Approach 3: Greedy (Optimal) ===")
	fmt.Printf("nums=[2,3,1,1,4]  got=%v  expected true\n", greedy([]int{2, 3, 1, 1, 4}))
	fmt.Printf("nums=[3,2,1,0,4]  got=%v  expected false\n", greedy([]int{3, 2, 1, 0, 4}))
	fmt.Printf("nums=[0]          got=%v  expected true\n", greedy([]int{0}))
	fmt.Printf("nums=[2,0,0]      got=%v  expected true\n", greedy([]int{2, 0, 0}))
}
