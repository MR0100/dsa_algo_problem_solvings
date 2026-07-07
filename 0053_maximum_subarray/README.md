# 0053 — Maximum Subarray

> LeetCode #53 · Difficulty: Medium
> **Categories:** Array, Divide and Conquer, Dynamic Programming

---

## Problem Statement

Given an integer array `nums`, find the **subarray** with the largest sum, and return its sum.

**Example 1**
```
Input:  nums = [-2,1,-3,4,-1,2,1,-5,4]
Output: 6
Explanation: The subarray [4,-1,2,1] has the largest sum 6.
```

**Example 2**
```
Input:  nums = [1]
Output: 1
```

**Example 3**
```
Input:  nums = [5,4,-1,7,8]
Output: 23
```

**Constraints**
- `1 <= nums.length <= 10⁵`
- `-10⁴ <= nums[i] <= 10⁴`

**Follow-up:** If you have figured out the O(n) solution, try coding another solution using the **divide and conquer** approach, which is more subtle.

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

- **Kadane's Algorithm** — greedy single-pass O(n) solution.
- **Dynamic Programming** — `dp[i]` = max subarray sum ending at `i`.
- **Divide and Conquer** — split at mid; best subarray is in left, right, or crosses mid.

---

## Approaches Overview

| # | Approach | Time | Space | When to use |
|---|----------|------|-------|-------------|
| 1 | Brute Force (all pairs) | O(n²) | O(1) | Reference; TLE for n=10⁵ |
| 2 | Kadane's Algorithm ✅ | O(n) | O(1) | The canonical interview answer |
| 3 | Divide and Conquer | O(n log n) | O(log n) | Follow-up requirement; good for distributed systems |
| 4 | DP Bottom-Up | O(n) | O(n) | Textbook DP formulation; can be optimised to O(1) space |

---

## Approach 1 — Brute Force

### Intuition
Try every subarray `nums[i..j]`, accumulate sum incrementally, track the maximum.

### Complexity
- **Time:** O(n²).
- **Space:** O(1).

### Code
```go
func bruteForce(nums []int) int {
    best := math.MinInt64
    for i := 0; i < len(nums); i++ {
        sum := 0
        for j := i; j < len(nums); j++ {
            sum += nums[j]
            if sum > best {
                best = sum
            }
        }
    }
    return best
}
```

### Dry Run — `nums = [-2,1,-3,4,-1,2,1,-5,4]`

Each outer `i` restarts `sum=0` and extends `j` from `i` to the end, tracking the best sum seen. Showing the running `sum` for a few representative starting points and the global `best`:

| i (start) | sums as j advances (`nums[i..j]`)                              | max for this i | global best |
|-----------|----------------------------------------------------------------|----------------|-------------|
| 0         | -2, -1, -4, 0, -1, 1, 2, -3, 1                                  | 2              | 2           |
| 1         | 1, -2, 2, 1, 3, 4, -1, 3                                        | 4              | 4           |
| 3         | 4, 3, 5, **6**, 1, 5                                            | **6**          | **6**       |
| 4         | -1, 1, 2, -3, 1                                                 | 2              | 6           |
| 8         | 4                                                               | 4              | 6           |

The best over all starting points is 6, achieved by `i=3` extending to `j=6` (subarray `[4,-1,2,1]`). Result: 6 ✓

---

## Approach 2 — Kadane's Algorithm (Recommended ✅)

### Intuition
At each index, the maximum subarray ending here either:
1. Extends the previous best subarray: `curSum + nums[i]`.
2. Starts fresh from here: `nums[i]`.

So: if `curSum` drops below 0, any subarray starting before this point drags the future sum down. Reset `curSum = 0` (effectively starting a new subarray at the next element).

### Algorithm
```
curSum = 0; best = nums[0]
for each num in nums:
  curSum += num
  best = max(best, curSum)
  if curSum < 0: curSum = 0
```

### Complexity
- **Time:** O(n) — one pass.
- **Space:** O(1).

### Code
```go
func kadane(nums []int) int {
    best, curSum := nums[0], 0
    for _, num := range nums {
        curSum += num
        if curSum > best { best = curSum }
        if curSum < 0   { curSum = 0 }
    }
    return best
}
```

### Dry Run — `nums = [-2,1,-3,4,-1,2,1,-5,4]`

| i | num | curSum before reset | best |
|---|-----|---------------------|------|
| 0 | -2  | -2 → reset → 0      | -2   |
| 1 | 1   | 1                   | 1    |
| 2 | -3  | -2 → reset → 0      | 1    |
| 3 | 4   | 4                   | 4    |
| 4 | -1  | 3                   | 4    |
| 5 | 2   | 5                   | 5    |
| 6 | 1   | 6                   | **6**|
| 7 | -5  | 1                   | 6    |
| 8 | 4   | 5                   | 6    |

Result: 6 ✓ (subarray [4,-1,2,1] = indices 3–6)

---

## Approach 3 — Divide and Conquer

### Intuition
Split at midpoint `mid`. The answer is the max of:
- Best subarray entirely in `left = nums[lo..mid]` (recurse).
- Best subarray entirely in `right = nums[mid+1..hi]` (recurse).
- Best subarray crossing `mid`: expand left from `mid` and right from `mid+1`, sum both arms.

### Algorithm
```
maxSubArray(lo, hi):
  if lo == hi: return nums[lo]
  mid = (lo + hi) / 2
  return max(maxSubArray(lo, mid), maxSubArray(mid+1, hi), maxCross(lo, mid, hi))

maxCross(lo, mid, hi):
  leftBest = -∞; sum = 0
  for i = mid downto lo: sum += nums[i]; leftBest = max(leftBest, sum)
  rightBest = -∞; sum = 0
  for i = mid+1 to hi: sum += nums[i]; rightBest = max(rightBest, sum)
  return leftBest + rightBest
```

### Complexity
- **Time:** O(n log n) — `T(n) = 2T(n/2) + O(n)`, Master theorem case 2.
- **Space:** O(log n) — recursion stack depth.

### Code
```go
func divideAndConquer(nums []int) int {
    return dac(nums, 0, len(nums)-1)
}

func dac(nums []int, lo, hi int) int {
    if lo == hi {
        return nums[lo]
    }
    mid := (lo + hi) / 2
    left := dac(nums, lo, mid)
    right := dac(nums, mid+1, hi)
    cross := maxCross(nums, lo, mid, hi)
    return max3(left, right, cross)
}

// maxCross computes the max subarray sum that crosses mid.
func maxCross(nums []int, lo, mid, hi int) int {
    leftSum := math.MinInt64
    sum := 0
    for i := mid; i >= lo; i-- { // expand left from mid
        sum += nums[i]
        if sum > leftSum {
            leftSum = sum
        }
    }
    rightSum := math.MinInt64
    sum = 0
    for i := mid + 1; i <= hi; i++ { // expand right from mid+1
        sum += nums[i]
        if sum > rightSum {
            rightSum = sum
        }
    }
    return leftSum + rightSum
}

func max3(a, b, c int) int {
    if a >= b && a >= c {
        return a
    }
    if b >= c {
        return b
    }
    return c
}
```

### Dry Run — `nums = [-2,1,-3,4,-1,2,1,-5,4]` (indices 0..8)

`dac(lo,hi)` splits at `mid=(lo+hi)/2` and returns `max(left, right, cross)`. Tracing the top-level call and its key recursions:

| Call `dac(lo,hi)` | mid | left | right | cross (`maxCross`) | returns |
|-------------------|-----|------|-------|--------------------|---------|
| dac(0,8)          | 4   | 4 (from 0..4) | 4 (from 5..8) | **6** (arms across mid=4) | **6** |
| ↳ dac(0,4)        | 2   | 1 (0..2) | 4 (3..4) | 4 | 4 |
| ↳ dac(5,8)        | 6   | 3 (5..6) | 4 (7..8) | 2 | 4 |
| ↳ dac(3,4)        | 3   | 4 (nums[3]) | -1 (nums[4]) | 3 | 4 |
| ↳ dac(5,6)        | 5   | 2 (nums[5]) | 1 (nums[6]) | 3 | 3 |

The winner is the **crossing** subarray of the top call. `maxCross(0,4,8)` around `mid=4`: the left arm expands down from index 4 and its best is `nums[3..4]=[4,-1]` = 3; the right arm expands up from index 5 and its best is `nums[5..6]=[2,1]` = 3. Sum of arms = 3 + 3 = 6, the subarray `[4,-1,2,1]`. The top call returns `max(4, 4, 6) = 6`. Result: 6 ✓

---

## Approach 4 — DP Bottom-Up

### Intuition
`dp[i]` = maximum subarray sum ending at index `i`.
- `dp[0] = nums[0]`.
- `dp[i] = max(nums[i], dp[i-1] + nums[i])` — start fresh or extend.

Answer = `max(dp[0..n-1])`.

This is equivalent to Kadane's but stores all dp values explicitly (reducible to O(1) space by only keeping `dp[i-1]`).

### Complexity
- **Time:** O(n).
- **Space:** O(n) — reducible to O(1) by using a single variable.

### Code
```go
func dpBottomUp(nums []int) int {
    n := len(nums)
    dp := make([]int, n)
    dp[0] = nums[0]
    best := dp[0]
    for i := 1; i < n; i++ {
        if dp[i-1]+nums[i] > nums[i] {
            dp[i] = dp[i-1] + nums[i]
        } else {
            dp[i] = nums[i]
        }
        if dp[i] > best {
            best = dp[i]
        }
    }
    return best
}
```

### Dry Run — `nums = [-2,1,-3,4,-1,2,1,-5,4]`

`dp[i] = max(nums[i], dp[i-1] + nums[i])`; `best = max(dp)`.

| i | nums[i] | dp[i-1] | dp[i-1]+nums[i] | dp[i] = max(nums[i], dp[i-1]+nums[i]) | best |
|---|---------|---------|-----------------|---------------------------------------|------|
| 0 | -2      | —       | —               | -2                                    | -2   |
| 1 | 1       | -2      | -1              | 1                                     | 1    |
| 2 | -3      | 1       | -2              | -2                                    | 1    |
| 3 | 4       | -2      | 2               | 4                                     | 4    |
| 4 | -1      | 4       | 3               | 3                                     | 4    |
| 5 | 2       | 3       | 5               | 5                                     | 5    |
| 6 | 1       | 5       | 6               | 6                                     | **6**|
| 7 | -5      | 6       | 1               | 1                                     | 6    |
| 8 | 4       | 1       | 5               | 5                                     | 6    |

Max `dp[i]` is 6 at `i=6` (subarray `[4,-1,2,1]`). Result: 6 ✓

---

## Key Takeaways

- **Kadane's is the answer for interviews** — O(n) time, O(1) space, 5 lines.
- **Reset-to-zero vs reset-to-current** — resetting to 0 (not to `nums[i]`) is correct when we can start a new subarray at any index; the loop's next iteration will add `nums[i]`.
- **All-negative arrays** — the answer is the single largest element. Kadane's handles this correctly because `best` is initialised to `nums[0]` (not `0`) and we update `best` before potentially resetting `curSum`.
- **Divide and conquer insight** — useful if you need to report the subarray indices (not just the sum), or in a parallel/distributed context where each machine handles a sub-range.
- **Follow-up: return indices** — track `start`, `end`, and `tempStart` in Kadane's loop.

---

## Related Problems

- LeetCode #918 — Maximum Sum Circular Subarray (Kadane's + circular adaptation)
- LeetCode #152 — Maximum Product Subarray (similar greedy; track min/max)
- LeetCode #560 — Subarray Sum Equals K (count subarrays; prefix sum + hash map)
