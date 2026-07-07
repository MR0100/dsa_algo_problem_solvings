package main

import "fmt"

// LeetCode 265 — Paint House II.
//
// There are n houses in a row, each to be painted one of k colors. costs[i][j]
// is the cost of painting house i with color j. No two ADJACENT houses may
// share a color. Return the minimum total painting cost.

// ── Approach 1: Dynamic Programming (Full Table) (Brute-ish DP) ───────────────
//
// dpFullTable solves Paint House II by computing, for every house and color,
// the min cost to paint up to that house ending in that color, scanning all
// previous colors for the best allowed predecessor.
//
// Intuition:
//
//	dp[i][j] = min cost to paint houses 0..i with house i colored j. To paint
//	house i color j, the previous house may be any color p != j, so
//	dp[i][j] = costs[i][j] + min over p!=j of dp[i-1][p]. The answer is the min
//	of the last row.
//
// Algorithm:
//  1. Initialise dp[0] = costs[0].
//  2. For each house i from 1: for each color j, find the min of dp[i-1][p]
//     over all p != j and add costs[i][j].
//  3. Return min of the final row.
//
// Time:  O(n·k²) — for each of n houses and k colors we scan k predecessors.
// Space: O(n·k) — the full dp table (can be reduced to two rows).
func dpFullTable(costs [][]int) int {
	if len(costs) == 0 {
		return 0
	}
	n, k := len(costs), len(costs[0])
	dp := make([][]int, n) // dp[i][j] = best cost for houses 0..i ending color j
	for i := range dp {
		dp[i] = make([]int, k)
	}
	copy(dp[0], costs[0]) // base row: cost of painting only the first house
	for i := 1; i < n; i++ {
		for j := 0; j < k; j++ {
			best := -1 // min dp[i-1][p] over p != j
			for p := 0; p < k; p++ {
				if p == j { // adjacent houses cannot share a color
					continue
				}
				if best == -1 || dp[i-1][p] < best {
					best = dp[i-1][p]
				}
			}
			dp[i][j] = costs[i][j] + best // extend the cheapest allowed previous color
		}
	}
	ans := dp[n-1][0] // minimum over the last house's colors
	for j := 1; j < k; j++ {
		if dp[n-1][j] < ans {
			ans = dp[n-1][j]
		}
	}
	return ans
}

// ── Approach 2: DP with Min1/Min2 Tracking (Optimal) ─────────────────────────
//
// dpMinTwo solves Paint House II in O(n·k) by noticing that the min over
// p != j of the previous row is the overall smallest previous cost, UNLESS the
// smallest one is at color j itself — in which case we use the second smallest.
//
// Intuition:
//
//	When extending to color j, we need the cheapest previous cost among colors
//	other than j. Precompute the previous row's smallest value (min1, at index
//	idx1) and its second smallest (min2). Then for each j: use min1 if j != idx1,
//	else min2. That removes the inner O(k) scan.
//
// Algorithm:
//  1. prev = costs[0]. Compute (min1, idx1, min2) of prev.
//  2. For each house i: for each color j, cur[j] = costs[i][j] +
//     (min1 if j != idx1 else min2). Recompute (min1, idx1, min2) from cur.
//  3. Answer is min1 after the last house.
//
// Time:  O(n·k) — one linear pass per house plus O(k) min-tracking.
// Space: O(k) — two rolling rows (prev / cur).
func dpMinTwo(costs [][]int) int {
	if len(costs) == 0 {
		return 0
	}
	n, k := len(costs), len(costs[0])
	if k == 1 { // single color: houses can't be adjacent-distinct unless n==1
		if n == 1 {
			return costs[0][0]
		}
		return 0 // per constraints k>=2 when n>1; guard anyway
	}
	prev := make([]int, k)
	copy(prev, costs[0]) // base row = first house costs
	// min1 = smallest in prev, idx1 its index, min2 = second smallest.
	min1, idx1, min2 := minTwo(prev)
	for i := 1; i < n; i++ {
		cur := make([]int, k)
		for j := 0; j < k; j++ {
			best := min1 // cheapest previous cost avoiding color j
			if j == idx1 {
				best = min2 // smallest is at j itself ⇒ take second smallest
			}
			cur[j] = costs[i][j] + best
		}
		min1, idx1, min2 = minTwo(cur) // refresh trackers for next house
		prev = cur
	}
	return min1 // min over the last row is the answer
}

// minTwo returns the smallest value, its index, and the second-smallest value.
func minTwo(row []int) (min1 int, idx1 int, min2 int) {
	min1, min2 = 1<<62, 1<<62 // "infinity" sentinels
	idx1 = -1
	for j, v := range row {
		if v < min1 { // new overall minimum; old min1 becomes min2
			min2 = min1
			min1 = v
			idx1 = j
		} else if v < min2 { // new runner-up
			min2 = v
		}
	}
	return
}

func main() {
	// Example 1: costs = [[1,5,3],[2,9,4]] ⇒ 5
	//   (house0 color0 =1 + house1 color2 =4 = 5, or color2=3 + color0=2 = 5).
	costs1 := [][]int{{1, 5, 3}, {2, 9, 4}}

	// Example 2: costs = [[1,3],[2,4]] ⇒ 5
	//   (house0 color0=1 + house1 color1=4 = 5).
	costs2 := [][]int{{1, 3}, {2, 4}}

	fmt.Println("=== Approach 1: DP Full Table ===")
	fmt.Println(dpFullTable(costs1)) // expected 5
	fmt.Println(dpFullTable(costs2)) // expected 5

	fmt.Println("=== Approach 2: DP Min1/Min2 (Optimal) ===")
	fmt.Println(dpMinTwo(costs1)) // expected 5
	fmt.Println(dpMinTwo(costs2)) // expected 5
}
