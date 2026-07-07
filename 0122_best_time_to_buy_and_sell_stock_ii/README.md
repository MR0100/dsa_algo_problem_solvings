# 0122 — Best Time to Buy and Sell Stock II

> LeetCode #122 · Difficulty: Medium
> **Categories:** Array, Dynamic Programming, Greedy

---

## Problem Statement

You are given an integer array `prices` where `prices[i]` is the price of a given stock on the `i`th day.

On each day, you may decide to buy and/or sell the stock. You can only hold **at most one share** at a time. However, you can buy it then immediately sell it on the same day.

Find and return the **maximum profit** you can achieve.

**Example 1:**
```
Input: prices = [7,1,5,3,6,4]
Output: 7
Explanation: Buy day 2 (1), sell day 3 (5), profit 4. Buy day 4 (3), sell day 5 (6), profit 3. Total 7.
```

**Example 2:**
```
Input: prices = [1,2,3,4,5]
Output: 4
Explanation: Buy day 1, sell day 5. Profit 4.
```

**Example 3:**
```
Input: prices = [7,6,4,3,1]
Output: 0
```

**Constraints:**
- `1 <= prices.length <= 3 * 10^4`
- `0 <= prices[i] <= 10^4`

---

## Company Frequency

| Company   | Frequency       | Last Reported |
|-----------|-----------------|---------------|
| Amazon    | ★★★★★ Very High | 2024          |
| Microsoft | ★★★★☆ High      | 2024          |
| Google    | ★★★★☆ High      | 2024          |
| Bloomberg | ★★★★☆ High      | 2024          |

> ⚠️ Frequency data is crowd-sourced from LeetCode Discuss, Glassdoor, and community interview reports. Treat as a signal, not a guarantee.

---

## DSA Concepts Used

- **Greedy** — capture every upward move

---

## Approaches Overview

| # | Approach            | Time | Space | When to use |
|---|---------------------|------|-------|-------------|
| 1 | Greedy (sum gains)  | O(n) | O(1)  | Always      |
| 2 | Peak-Valley Explicit| O(n) | O(1)  | Intuitive   |

---

## Approach 1 — Greedy

### Intuition
Sum up all positive consecutive differences. `prices[i]-prices[i-1]` contributes to profit whenever positive. This is optimal because any buying-selling strategy is decomposable into consecutive segments.

### Complexity
- **Time:** O(n)
- **Space:** O(1)

### Code
```go
func maxProfit(prices []int) int {
    profit := 0
    for i := 1; i < len(prices); i++ {
        if prices[i] > prices[i-1] { profit += prices[i] - prices[i-1] }
    }
    return profit
}
```

### Dry Run
`[7,1,5,3,6,4]`:

| i | diff | add? | profit |
|---|------|------|--------|
| 1 | -6   | no   | 0      |
| 2 | +4   | yes  | 4      |
| 3 | -2   | no   | 4      |
| 4 | +3   | yes  | 7      |
| 5 | -2   | no   | 7      |

---

## Approach 2 — Peak-Valley Explicit

### Intuition
Find each valley (buy) and peak (sell) pair. Sum `(peak - valley)` for each ascending segment.

### Complexity
- **Time:** O(n)
- **Space:** O(1)

### Code
```go
func maxProfitPeakValley(prices []int) int {
    n := len(prices); profit := 0; i := 0
    for i < n-1 {
        for i < n-1 && prices[i] >= prices[i+1] { i++ }
        valley := prices[i]
        for i < n-1 && prices[i] <= prices[i+1] { i++ }
        peak := prices[i]
        profit += peak - valley
    }
    return profit
}
```

### Dry Run
`[7,1,5,3,6,4]`:
- Find valley: i=1 (price=1). Find peak: i=2 (price=5). profit=4.
- Find valley: i=3 (price=3). Find peak: i=4 (price=6). profit=7.

---

## Key Takeaways
- Unlimited transactions = sum of all upward moves.
- Greedy approach 1 is cleanest — one loop, constant space.
- Contrast with #121 (1 transaction) and #123 (2 transactions).

---

## Related Problems
- LeetCode #121 — Best Time to Buy and Sell Stock (1 transaction)
- LeetCode #123 — Best Time to Buy and Sell Stock III (2 transactions)
- LeetCode #309 — Best Time to Buy and Sell Stock with Cooldown
