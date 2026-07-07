# 0070 — Climbing Stairs

> LeetCode #70 · Difficulty: Easy
> **Categories:** Math, Dynamic Programming, Memoization

---

## Problem Statement

You are climbing a staircase. It takes `n` steps to reach the top.

Each time you can either climb `1` or `2` steps. In how many distinct ways can you climb to the top?

**Example 1**
```
Input:  n = 2
Output: 2
Explanation: There are two ways to climb to the top.
1. 1 step + 1 step
2. 2 steps
```

**Example 2**
```
Input:  n = 3
Output: 3
Explanation: There are three ways to climb to the top.
1. 1 step + 1 step + 1 step
2. 1 step + 2 steps
3. 2 steps + 1 step
```

**Constraints**
- `1 <= n <= 45`

---

## Company Frequency

| Company   | Frequency       | Last Reported |
|-----------|-----------------|---------------|
| Amazon    | ★★★★★ Very High | 2024          |
| Google    | ★★★★★ Very High | 2024          |
| Meta      | ★★★★★ Very High | 2024          |
| Microsoft | ★★★★★ Very High | 2024          |
| Bloomberg | ★★★★☆ High      | 2024          |
| Apple     | ★★★★☆ High      | 2023          |
| Adobe     | ★★★☆☆ Medium    | 2023          |

> ⚠️ Frequency data is crowd-sourced from LeetCode Discuss, Glassdoor, and community interview reports. Treat as a signal, not a guarantee.

---

## DSA Concepts Used

- **Dynamic Programming (Fibonacci)** — `ways(n) = ways(n-1) + ways(n-2)`.
- **Space Optimization** — only the last two values needed; reduce O(n) → O(1) space.

---

## Approaches Overview

| # | Approach | Time | Space | When to use |
|---|----------|------|-------|-------------|
| 1 | Memoization (Top-Down DP) | O(n) | O(n) | Good learning step; explicit recursion |
| 2 | DP Bottom-Up (Table) | O(n) | O(n) | Textbook DP |
| 3 | Two Variables ✅ | O(n) | O(1) | Optimal; the interview answer |

---

## Approach 1 — Memoization (Top-Down DP)

### Intuition
`ways(n)` = ways to reach step n = ways to reach n-1 (then take 1 step) + ways to reach n-2 (then take 2 steps). This is the Fibonacci recurrence.

Base: `ways(0) = 1` (one way to stand at bottom), `ways(1) = 1`.

### Complexity
- **Time:** O(n).
- **Space:** O(n) — memo + stack.

---

## Approach 2 — DP Bottom-Up

### Intuition
Fill `dp[0..n]` iteratively. `dp[i] = dp[i-1] + dp[i-2]`.

### Complexity
- **Time:** O(n).
- **Space:** O(n).

### Dry Run — `n = 5`
```
dp[0]=1, dp[1]=1
dp[2]=2, dp[3]=3, dp[4]=5, dp[5]=8
Return dp[5] = 8 ✓
```

---

## Approach 3 — Two Variables (Recommended ✅)

### Intuition
The DP table only needs the last two values. Keep `prev` and `curr`, update in a loop.

### Algorithm
```
prev=1, curr=1
for i=2 to n: prev, curr = curr, prev+curr
return curr
```

### Complexity
- **Time:** O(n).
- **Space:** O(1).

### Code
```go
func twoVars(n int) int {
    if n <= 1 { return 1 }
    prev, curr := 1, 1
    for i := 2; i <= n; i++ { prev, curr = curr, prev+curr }
    return curr
}
```

### Dry Run — `n = 5`
```
Start: prev=1, curr=1
i=2: prev=1, curr=1+1=2
i=3: prev=2, curr=1+2=3
i=4: prev=3, curr=2+3=5
i=5: prev=5, curr=3+5=8
Return 8 ✓
```

---

## Key Takeaways

- **This IS Fibonacci** — `ways(n) = Fib(n+1)` (with Fib(1)=1, Fib(2)=1). `n=1→1, n=2→2, n=3→3, n=4→5, ...`.
- **Follow-up: k steps** — if you can take 1 to k steps, `ways(n) = sum(ways(n-1..n-k))`. Use a sliding window sum for O(n) time.
- **n ≤ 45** — the answer fits in a 32-bit integer (Fib(46) = 1,836,311,903 < 2³¹-1).
- **Most famous intro DP problem** — the interviewer knows you know it; what they're testing is whether you can articulate the recurrence clearly and optimise space.

---

## Related Problems

- LeetCode #509 — Fibonacci Number (exact same pattern)
- LeetCode #746 — Min Cost Climbing Stairs (DP with cost; add min() to the recurrence)
- LeetCode #91 — Decode Ways (counting paths with conditions; DP extension)
- LeetCode #377 — Combination Sum IV (k-step climbing generalisation)
