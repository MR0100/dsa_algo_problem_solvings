package main

import "fmt"

// ── Approach 1: Greedy (Peak-Valley) ─────────────────────────────────────────
//
// maxProfit solves Best Time to Buy and Sell Stock II greedily.
//
// Intuition:
//   Since we can hold multiple transactions (but only one at a time), we should
//   capture every upward slope. If prices[i] > prices[i-1], add the difference.
//   This is equivalent to finding all ascending consecutive pairs.
//
// Time:  O(n)
// Space: O(1)
func maxProfit(prices []int) int {
	profit := 0
	for i := 1; i < len(prices); i++ {
		if prices[i] > prices[i-1] {
			profit += prices[i] - prices[i-1] // capture every gain
		}
	}
	return profit
}

// ── Approach 2: Peak-Valley Explicit ─────────────────────────────────────────
//
// maxProfitPeakValley solves Best Time to Buy and Sell Stock II by explicitly
// finding local minima (valleys) and maxima (peaks).
//
// Intuition:
//   Buy at every valley, sell at every peak.
//   profit = sum of (peak - valley) for every ascending segment.
//
// Time:  O(n)
// Space: O(1)
func maxProfitPeakValley(prices []int) int {
	n := len(prices)
	profit := 0
	i := 0
	for i < n-1 {
		// find valley
		for i < n-1 && prices[i] >= prices[i+1] {
			i++
		}
		valley := prices[i]
		// find peak
		for i < n-1 && prices[i] <= prices[i+1] {
			i++
		}
		peak := prices[i]
		profit += peak - valley
	}
	return profit
}

func main() {
	fmt.Println("=== Approach 1: Greedy ===")
	fmt.Printf("prices=[7,1,5,3,6,4]  got=%d  expected 7\n", maxProfit([]int{7, 1, 5, 3, 6, 4}))
	fmt.Printf("prices=[1,2,3,4,5]  got=%d  expected 4\n", maxProfit([]int{1, 2, 3, 4, 5}))
	fmt.Printf("prices=[7,6,4,3,1]  got=%d  expected 0\n", maxProfit([]int{7, 6, 4, 3, 1}))

	fmt.Println("=== Approach 2: Peak-Valley ===")
	fmt.Printf("prices=[7,1,5,3,6,4]  got=%d  expected 7\n", maxProfitPeakValley([]int{7, 1, 5, 3, 6, 4}))
	fmt.Printf("prices=[1,2,3,4,5]  got=%d  expected 4\n", maxProfitPeakValley([]int{1, 2, 3, 4, 5}))
	fmt.Printf("prices=[7,6,4,3,1]  got=%d  expected 0\n", maxProfitPeakValley([]int{7, 6, 4, 3, 1}))
}
