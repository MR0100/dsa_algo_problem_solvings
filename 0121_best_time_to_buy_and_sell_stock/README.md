# 0121 — Best Time to Buy and Sell Stock

> LeetCode #121 · Difficulty: Easy
> **Categories:** Array, Dynamic Programming

---

## Problem Statement

You are given an array `prices` where `prices[i]` is the price of a given stock on the `i`th day.

You want to maximize your profit by choosing a **single day** to buy one stock and choosing a **different day in the future** to sell that stock.

Return the maximum profit you can achieve from this transaction. If you cannot achieve any profit, return `0`.

**Example 1:**
```
Input: prices = [7,1,5,3,6,4]
Output: 5
Explanation: Buy on day 2 (price=1), sell on day 5 (price=6), profit=6-1=5.
```

**Example 2:**
```
Input: prices = [7,6,4,3,1]
Output: 0
Explanation: No profitable transaction possible.
```

**Constraints:**
- `1 <= prices.length <= 10^5`
- `0 <= prices[i] <= 10^4`

---

## Company Frequency

| Company   | Frequency       | Last Reported |
|-----------|-----------------|---------------|
| Amazon    | ★★★★★ Very High | 2024          |
| Microsoft | ★★★★★ Very High | 2024          |
| Google    | ★★★★★ Very High | 2024          |
| Facebook  | ★★★★★ Very High | 2024          |
| Bloomberg | ★★★★★ Very High | 2024          |

> ⚠️ Frequency data is crowd-sourced from LeetCode Discuss, Glassdoor, and community interview reports. Treat as a signal, not a guarantee.

---

## DSA Concepts Used

- **Greedy / One-pass scan** — track minimum price seen so far
- **Dynamic Programming** — Kadane-like structure

---

## Approaches Overview

| # | Approach      | Time   | Space | When to use       |
|---|---------------|--------|-------|-------------------|
| 1 | Brute Force   | O(n²)  | O(1)  | Understanding only|
| 2 | One Pass      | O(n)   | O(1)  | Always            |

---

## Approach 1 — Brute Force

### Intuition
Try every buy-sell pair and take the maximum profit.

### Complexity
- **Time:** O(n²)
- **Space:** O(1)

### Code
```go
func maxProfitBrute(prices []int) int {
    maxP := 0
    for i := 0; i < len(prices); i++ {
        for j := i+1; j < len(prices); j++ {
            if prices[j]-prices[i] > maxP { maxP = prices[j]-prices[i] }
        }
    }
    return maxP
}
```

### Dry Run
`[7,1,5,3,6,4]`: best pair is (i=1,j=4): 6-1=5.

---

## Approach 2 — One Pass (Optimal)

### Intuition
Track `minPrice` seen so far. For each price, compute `price - minPrice`. Update maxProfit.
This works because: if we're at day `j`, the best buy day is the cheapest day before `j`.

### Algorithm
1. `minPrice = prices[0]`, `maxP = 0`.
2. For each price:
   - If `price < minPrice`: update minPrice (better buy day).
   - Else if `price - minPrice > maxP`: update maxP.
3. Return maxP.

### Complexity
- **Time:** O(n)
- **Space:** O(1)

### Code
```go
func maxProfit(prices []int) int {
    minPrice := prices[0]; maxP := 0
    for _, price := range prices {
        if price < minPrice { minPrice = price }
        else if price-minPrice > maxP { maxP = price-minPrice }
    }
    return maxP
}
```

### Dry Run
`[7,1,5,3,6,4]`:

| price | minPrice | profit | maxP |
|-------|----------|--------|------|
| 7     | 7        | 0      | 0    |
| 1     | 1        | 0      | 0    |
| 5     | 1        | 4      | 4    |
| 3     | 1        | 2      | 4    |
| 6     | 1        | 5      | 5    |
| 4     | 1        | 3      | 5    |

Result: 5 ✓

---

## Key Takeaways
- Single transaction: track running minimum and compute profit at each step.
- Equivalent to Kadane's algorithm on the difference array `prices[i]-prices[i-1]`.
- See #122 (unlimited transactions), #123 (at most 2), #188 (at most k).

---

## Related Problems
- LeetCode #122 — Best Time to Buy and Sell Stock II (multiple transactions)
- LeetCode #123 — Best Time to Buy and Sell Stock III (at most 2 transactions)
- LeetCode #188 — Best Time to Buy and Sell Stock IV (at most k transactions)
