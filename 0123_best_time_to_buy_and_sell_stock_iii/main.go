package main

import "fmt"

// ── Approach 1: DP — Four States ─────────────────────────────────────────────
//
// maxProfit solves Best Time to Buy and Sell Stock III (at most 2 transactions)
// using four state variables.
//
// Intuition:
//   Track four states across all prices:
//   - buy1:  max cash after first buy  = max(buy1,  -price)
//   - sell1: max cash after first sell = max(sell1,  buy1 + price)
//   - buy2:  max cash after second buy = max(buy2,  sell1 - price)
//   - sell2: max cash after second sell= max(sell2,  buy2 + price)
//   Initialise buy1=buy2=-inf, sell1=sell2=0.
//
// Time:  O(n)
// Space: O(1)
func maxProfit(prices []int) int {
	buy1, sell1 := -1<<31, 0
	buy2, sell2 := -1<<31, 0

	for _, price := range prices {
		buy1 = max(buy1, -price)         // best we can do after first buy
		sell1 = max(sell1, buy1+price)   // best after first sell
		buy2 = max(buy2, sell1-price)    // best after second buy
		sell2 = max(sell2, buy2+price)   // best after second sell
	}
	return sell2
}

func max(a, b int) int {
	if a > b { return a }
	return b
}

// ── Approach 2: Left-Right DP ─────────────────────────────────────────────────
//
// maxProfitLeftRight solves Best Time to Buy and Sell Stock III by precomputing:
//   leftMax[i]  = max profit from a single transaction in prices[0..i]
//   rightMax[i] = max profit from a single transaction in prices[i..n-1]
// Then answer = max over all splits: leftMax[i] + rightMax[i+1].
//
// Time:  O(n)
// Space: O(n)
func maxProfitLeftRight(prices []int) int {
	n := len(prices)
	if n == 0 {
		return 0
	}
	leftMax := make([]int, n)
	minLeft := prices[0]
	for i := 1; i < n; i++ {
		if prices[i]-minLeft > leftMax[i-1] {
			leftMax[i] = prices[i] - minLeft
		} else {
			leftMax[i] = leftMax[i-1]
		}
		if prices[i] < minLeft {
			minLeft = prices[i]
		}
	}

	rightMax := make([]int, n)
	maxRight := prices[n-1]
	for i := n - 2; i >= 0; i-- {
		if maxRight-prices[i] > rightMax[i+1] {
			rightMax[i] = maxRight - prices[i]
		} else {
			rightMax[i] = rightMax[i+1]
		}
		if prices[i] > maxRight {
			maxRight = prices[i]
		}
	}

	ans := leftMax[n-1] // only first transaction
	for i := 0; i < n-1; i++ {
		if leftMax[i]+rightMax[i+1] > ans {
			ans = leftMax[i] + rightMax[i+1]
		}
	}
	return ans
}

func main() {
	fmt.Println("=== Approach 1: Four States DP ===")
	fmt.Printf("prices=[3,3,5,0,0,3,1,4]  got=%d  expected 6\n", maxProfit([]int{3, 3, 5, 0, 0, 3, 1, 4}))
	fmt.Printf("prices=[1,2,3,4,5]  got=%d  expected 4\n", maxProfit([]int{1, 2, 3, 4, 5}))
	fmt.Printf("prices=[7,6,4,3,1]  got=%d  expected 0\n", maxProfit([]int{7, 6, 4, 3, 1}))

	fmt.Println("=== Approach 2: Left-Right DP ===")
	fmt.Printf("prices=[3,3,5,0,0,3,1,4]  got=%d  expected 6\n", maxProfitLeftRight([]int{3, 3, 5, 0, 0, 3, 1, 4}))
	fmt.Printf("prices=[1,2,3,4,5]  got=%d  expected 4\n", maxProfitLeftRight([]int{1, 2, 3, 4, 5}))
	fmt.Printf("prices=[7,6,4,3,1]  got=%d  expected 0\n", maxProfitLeftRight([]int{7, 6, 4, 3, 1}))
}
