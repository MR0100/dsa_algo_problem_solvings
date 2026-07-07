package main

import "fmt"

// ── Approach 1: BFS / Level-by-Level Expansion ───────────────────────────────
//
// bfs solves Jump Game II by treating each jump level as a BFS level.
//
// Intuition: Each position reachable with k jumps defines a "level". The
// farthest we can reach from each level determines what positions are reachable
// in k+1 jumps. We count levels until we reach or exceed the last index.
//
// Algorithm:
//  jumps=0, curEnd=0, farthest=0
//  for i = 0 to n-2:
//    farthest = max(farthest, i + nums[i])
//    if i == curEnd:  // finished current level
//      jumps++
//      curEnd = farthest
//      if curEnd >= n-1: break
//
// Time:  O(n)
// Space: O(1)
func bfs(nums []int) int {
	n := len(nums)
	if n <= 1 {
		return 0
	}
	jumps := 0
	curEnd := 0   // rightmost index reachable with current jump count
	farthest := 0 // rightmost index reachable with one more jump
	for i := 0; i < n-1; i++ {
		if i+nums[i] > farthest {
			farthest = i + nums[i] // update farthest reachable from i
		}
		if i == curEnd { // we've exhausted the current jump level
			jumps++
			curEnd = farthest
			if curEnd >= n-1 {
				break // already reached the end
			}
		}
	}
	return jumps
}

// ── Approach 2: Greedy (Optimal) ─────────────────────────────────────────────
//
// greedy solves Jump Game II identically to the BFS but framed as a greedy:
// at each position, track the farthest reachable index. Use a jump whenever we
// must advance to the next "frontier".
//
// This is the same algorithm as bfs above; it's presented separately to clarify
// the greedy intuition: always prefer the jump that gets you the farthest.
//
// Greedy correctness: We never need to jump before we have to (when we reach
// curEnd). At that point, the best choice is to jump to wherever farthest is.
//
// Time:  O(n)
// Space: O(1)
func greedy(nums []int) int {
	jumps, curEnd, farthest := 0, 0, 0
	for i := 0; i < len(nums)-1; i++ {
		if i+nums[i] > farthest {
			farthest = i + nums[i]
		}
		if i == curEnd {
			jumps++
			curEnd = farthest
		}
	}
	return jumps
}

// ── Approach 3: DP (for comparison) ──────────────────────────────────────────
//
// dpApproach solves Jump Game II with backward DP.
//
// dp[i] = minimum jumps needed to reach index i.
// For each i, look back at all j < i where j + nums[j] >= i:
//   dp[i] = min(dp[j] + 1).
//
// Time:  O(n²)
// Space: O(n)
func dpApproach(nums []int) int {
	n := len(nums)
	dp := make([]int, n)
	for i := range dp {
		dp[i] = 1<<31 - 1 // infinity
	}
	dp[0] = 0
	for i := 1; i < n; i++ {
		for j := 0; j < i; j++ {
			if j+nums[j] >= i && dp[j]+1 < dp[i] {
				dp[i] = dp[j] + 1
			}
		}
	}
	return dp[n-1]
}

func main() {
	cases := []struct {
		nums []int
		want int
	}{
		{[]int{2, 3, 1, 1, 4}, 2},
		{[]int{2, 3, 0, 1, 4}, 2},
		{[]int{1, 2, 3}, 2},
		{[]int{0}, 0},
		{[]int{1, 1, 1, 1}, 3},
	}

	fmt.Println("=== Approach 1: BFS / Level Expansion ===")
	for _, c := range cases {
		fmt.Printf("nums=%v => %d  expected %d\n", c.nums, bfs(c.nums), c.want)
	}

	fmt.Println("\n=== Approach 2: Greedy (Optimal) ===")
	for _, c := range cases {
		fmt.Printf("nums=%v => %d  expected %d\n", c.nums, greedy(c.nums), c.want)
	}

	fmt.Println("\n=== Approach 3: DP ===")
	for _, c := range cases {
		fmt.Printf("nums=%v => %d  expected %d\n", c.nums, dpApproach(c.nums), c.want)
	}
}
