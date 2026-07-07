package main

import "fmt"

// ── Approach 1: Brute Force ───────────────────────────────────────────────────
//
// maxProfitBrute solves Best Time to Buy and Sell Stock naively.
//
// Intuition:
//   Try every pair (i, j) with i < j and take the max difference.
//
// Time:  O(n^2)
// Space: O(1)
func maxProfitBrute(prices []int) int {
	maxP := 0
	for i := 0; i < len(prices); i++ {
		for j := i + 1; j < len(prices); j++ {
			if prices[j]-prices[i] > maxP {
				maxP = prices[j] - prices[i]
			}
		}
	}
	return maxP
}

// ── Approach 2: One Pass (Optimal) ───────────────────────────────────────────
//
// maxProfit solves Best Time to Buy and Sell Stock in one pass.
//
// Intuition:
//   Track the minimum price seen so far (best buy day).
//   At each price, compute profit if we sell today: price - minSoFar.
//   Update maxProfit accordingly.
//
// Time:  O(n)
// Space: O(1)
func maxProfit(prices []int) int {
	if len(prices) == 0 {
		return 0
	}
	minPrice := prices[0]
	maxP := 0

	for _, price := range prices {
		if price < minPrice {
			minPrice = price // found a cheaper buy day
		} else if price-minPrice > maxP {
			maxP = price - minPrice // found a better profit
		}
	}
	return maxP
}

func main() {
	fmt.Println("=== Approach 1: Brute Force ===")
	fmt.Printf("prices=[7,1,5,3,6,4]  got=%d  expected 5\n", maxProfitBrute([]int{7, 1, 5, 3, 6, 4}))
	fmt.Printf("prices=[7,6,4,3,1]  got=%d  expected 0\n", maxProfitBrute([]int{7, 6, 4, 3, 1}))

	fmt.Println("=== Approach 2: One Pass ===")
	fmt.Printf("prices=[7,1,5,3,6,4]  got=%d  expected 5\n", maxProfit([]int{7, 1, 5, 3, 6, 4}))
	fmt.Printf("prices=[7,6,4,3,1]  got=%d  expected 0\n", maxProfit([]int{7, 6, 4, 3, 1}))
}
