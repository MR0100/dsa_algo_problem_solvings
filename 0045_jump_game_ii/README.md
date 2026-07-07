# 0045 — Jump Game II

> LeetCode #45 · Difficulty: Medium
> **Categories:** Array, Dynamic Programming, Greedy

---

## Problem Statement

You are given a **0-indexed** array of integers `nums` of length `n`. You are initially positioned at `nums[0]`.

Each element `nums[i]` represents the maximum length of a forward jump from index `i`. In other words, if you are at `nums[i]`, you can jump to any `nums[i + j]` where:
- `0 <= j <= nums[i]` and
- `i + j < n`

Return the **minimum number of jumps** to reach `nums[n - 1]`. The test cases are generated such that you can always reach `nums[n - 1]`.

**Example 1**
```
Input:  nums = [2,3,1,1,4]
Output: 2
Explanation: Jump 1 step from index 0 to 1, then 3 steps to the last index.
```

**Example 2**
```
Input:  nums = [2,3,0,1,4]
Output: 2
```

**Constraints**
- `1 <= nums.length <= 10⁴`
- `0 <= nums[i] <= 1000`
- The test cases are generated such that you can always reach `nums[n-1]`.

---

## Company Frequency

| Company   | Frequency       | Last Reported |
|-----------|-----------------|---------------|
| Amazon    | ★★★★★ Very High | 2024          |
| Google    | ★★★★★ Very High | 2024          |
| Meta      | ★★★★★ Very High | 2024          |
| Microsoft | ★★★★☆ High      | 2024          |
| Bloomberg | ★★★★☆ High      | 2023          |
| Apple     | ★★★☆☆ Medium    | 2023          |

> ⚠️ Frequency data is crowd-sourced from LeetCode Discuss, Glassdoor, and community interview reports. Treat as a signal, not a guarantee.

---

## DSA Concepts Used

- **Greedy / BFS Level Expansion** — the key insight: treat each "jump" as a BFS level. All positions reachable in exactly k jumps form one level. The answer is the level at which we first reach (or pass) the last index.
- **Dynamic Programming** — Approach 3 uses a dp array; O(n²) but good for understanding.

---

## Approaches Overview

| # | Approach | Time | Space | When to use |
|---|----------|------|-------|-------------|
| 1 | BFS / Level Expansion ✅ | O(n) | O(1) | Intuitive; equivalent to greedy |
| 2 | Greedy ✅ | O(n) | O(1) | Cleaner code; same algorithm |
| 3 | DP | O(n²) | O(n) | Educational; clearly shows subproblem structure |

---

## Approach 1 — BFS / Level Expansion

### Intuition
Imagine the array as a graph: index 0 is the source, and from index `i` you can reach `i+1` to `i+nums[i]`. BFS finds the shortest path (in jumps). Each "level" in BFS is one jump.

Track:
- `curEnd` — the farthest index reachable with the current number of jumps.
- `farthest` — the farthest index reachable with one more jump from any index in the current level.

When we reach `curEnd` (end of current level), increment jumps and advance `curEnd = farthest`.

### Algorithm
```
jumps=0, curEnd=0, farthest=0
for i = 0 to n-2:
  farthest = max(farthest, i + nums[i])
  if i == curEnd:        // reached end of current BFS level
    jumps++
    curEnd = farthest
    if curEnd >= n-1: break
return jumps
```

### Complexity
- **Time:** O(n).
- **Space:** O(1).

### Dry Run — `nums = [2,3,1,1,4]`
```
jumps=0, curEnd=0, farthest=0

i=0: farthest=max(0,0+2)=2. i==curEnd(0): jumps=1, curEnd=2
i=1: farthest=max(2,1+3)=4. i<curEnd(2)
i=2: farthest=max(4,2+1)=4. i==curEnd(2): jumps=2, curEnd=4. 4>=4→break

return 2 ✓
```

---

## Approach 2 — Greedy (Recommended ✅)

### Intuition
At each position, greedily track the farthest index reachable so far. When forced to use a jump (when we've exhausted the current range), increment the jump count and extend the range.

This is the same algorithm as Approach 1, framed as greedy.

### Complexity
- **Time:** O(n).
- **Space:** O(1).

---

## Approach 3 — Dynamic Programming

### Intuition
`dp[i]` = minimum jumps to reach index `i`. For each `i`, look back at all `j < i` where `j + nums[j] >= i`, and take `dp[j] + 1`.

### Complexity
- **Time:** O(n²).
- **Space:** O(n).

---

## Key Takeaways

- **Loop to `n-2`, not `n-1`** — we process positions as sources; the last index `n-1` is the destination, not a source. Processing it would incorrectly increment `jumps`.
- **`farthest` is updated before the `if i==curEnd` check** — so when we increment `jumps` at `i==curEnd`, `farthest` already reflects the best jump from the entire current level.
- **Greedy choice is optimal** — at each "forced jump" point, taking the maximum reach is always best (no future choice can be hurt by reaching farther now, since reaching farther gives more options).
- **Difference from #55 (Jump Game I)** — #55 asks CAN you reach the end; this asks what is the MINIMUM number of jumps.

---

## Related Problems

- LeetCode #55 — Jump Game (can you reach the end; no count needed)
- LeetCode #1306 — Jump Game III (bidirectional jumps; BFS)
- LeetCode #1345 — Jump Game IV (different connections; BFS)
