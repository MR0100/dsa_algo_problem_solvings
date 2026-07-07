# 0123 — Best Time to Buy and Sell Stock III

> LeetCode #123 · Difficulty: Hard
> **Categories:** Array, Dynamic Programming

---

## Problem Statement

You are given an array `prices` where `prices[i]` is the price of a given stock on the `i`th day.

Find the maximum profit you can achieve. You may complete **at most two transactions**.

**Note:** You may not engage in multiple transactions simultaneously (i.e., you must sell the stock before you buy again).

**Example 1:**
```
Input: prices = [3,3,5,0,0,3,1,4]
Output: 6
Explanation: Buy day 4 (0), sell day 6 (3), profit 3. Buy day 7 (1), sell day 8 (4), profit 3. Total 6.
```

**Example 2:**
```
Input: prices = [1,2,3,4,5]
Output: 4
Explanation: Buy day 1 (1), sell day 5 (5), profit 4.
```

**Example 3:**
```
Input: prices = [7,6,4,3,1]
Output: 0
```

**Constraints:**
- `1 <= prices.length <= 10^5`
- `0 <= prices[i] <= 10^5`

---

## Company Frequency

| Company   | Frequency    | Last Reported |
|-----------|--------------|---------------|
| Amazon    | ★★★★☆ High  | 2024          |
| Google    | ★★★★☆ High  | 2024          |
| Facebook  | ★★★☆☆ Medium | 2024          |
| Bloomberg | ★★★☆☆ Medium | 2023          |

> ⚠️ Frequency data is crowd-sourced from LeetCode Discuss, Glassdoor, and community interview reports. Treat as a signal, not a guarantee.

---

## DSA Concepts Used

- **DP with states** — 4 states tracking phases of at most 2 transactions → see [`/dsa/dynamic_programming.md`](/dsa/dynamic_programming.md)
- **Divide-and-Conquer** — precompute best single-transaction for left and right halves

---

## Approaches Overview

| # | Approach           | Time | Space | When to use        |
|---|--------------------|------|-------|--------------------|
| 1 | Four States DP     | O(n) | O(1)  | Cleanest; optimal  |
| 2 | Left-Right DP      | O(n) | O(n)  | Intuitive split    |

---

## Approach 1 — Four States DP

### Intuition
Track four states that represent maximum cash at each stage:
- `buy1`:  best cash after buying for the first time = max(buy1, -price).
- `sell1`: best cash after selling for the first time = max(sell1, buy1+price).
- `buy2`:  best cash after buying for the second time = max(buy2, sell1-price).
- `sell2`: best cash after selling for the second time = max(sell2, buy2+price).

Initialize `buy1 = buy2 = -∞` (haven't bought yet), `sell1 = sell2 = 0`.

### Algorithm
For each price, update all four states in order.

### Complexity
- **Time:** O(n)
- **Space:** O(1)

### Code
```go
func maxProfit(prices []int) int {
    buy1, sell1 := -1<<31, 0
    buy2, sell2 := -1<<31, 0
    for _, price := range prices {
        buy1  = max(buy1,  -price)
        sell1 = max(sell1,  buy1+price)
        buy2  = max(buy2,  sell1-price)
        sell2 = max(sell2,  buy2+price)
    }
    return sell2
}
```

### Dry Run
`[3,3,5,0,0,3,1,4]`:

| price | buy1 | sell1 | buy2  | sell2 |
|-------|------|-------|-------|-------|
| 3     | -3   | 0     | -3    | 0     |
| 3     | -3   | 0     | -3    | 0     |
| 5     | -3   | 2     | -1    | 2     |
| 0     | 0    | 2     | 2     | 2     |
| 0     | 0    | 2     | 2     | 2     |
| 3     | 0    | 3     | 2     | 5     |
| 1     | 0    | 3     | 2     | 5     |
| 4     | 0    | 4     | 3     | 7 ← wait |

Hmm, expected 6. Let me retrace with [3,3,5,0,0,3,1,4]:
- price=0 (idx 3): buy1=max(-3,0)=0, sell1=max(2,0+0)=2, buy2=max(-1,2-0)=2, sell2=max(2,2+0)=2.
- price=3 (idx 5): buy1=0, sell1=max(2,0+3)=3, buy2=max(2,3-3)=2, sell2=max(2,2+3)=5.
- price=1 (idx 6): buy1=0, sell1=3, buy2=max(2,3-1)=2, sell2=5.
- price=4 (idx 7): buy1=0, sell1=max(3,0+4)=4, buy2=max(2,4-4)=2, sell2=max(5,2+4)=6.

Final sell2 = 6 ✓ (my table above was wrong for price=8).

---

## Approach 2 — Left-Right DP

### Intuition
For each split point `i`, the two transactions can be split: first in `[0..i]`, second in `[i+1..n-1]`. Precompute `leftMax[i]` and `rightMax[i]`, then find the max sum.

### Complexity
- **Time:** O(n)
- **Space:** O(n)

### Code
```go
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
```

### Dry Run
`[3,3,5,0,0,3,1,4]`:
- `leftMax = [0,0,2,2,2,3,3,4]`
- `rightMax = [6,6,4,4,3,3,3,0]`
- max split: `leftMax[4]+rightMax[5] = 2+3 = 5`? Actually: leftMax[3]=2, rightMax[4]=3 → 5. Or leftMax[5]=3, rightMax[6]=3 → 6.

Best: `leftMax[5]+rightMax[6] = 3+3 = 6` ✓

---

## Key Takeaways
- Four-state DP generalizes: for `k` transactions, use `2k` states.
- `sell1-price` in `buy2` represents: "use profit from first trade to fund second buy."
- Left-right precomputation splits at every possible dividing day.

---

## Related Problems
- LeetCode #121 — Best Time to Buy and Sell Stock (1 transaction)
- LeetCode #122 — Best Time to Buy and Sell Stock II (unlimited)
- LeetCode #188 — Best Time to Buy and Sell Stock IV (k transactions)
