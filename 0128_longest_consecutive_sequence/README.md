# 0128 — Longest Consecutive Sequence

> LeetCode #128 · Difficulty: Medium
> **Categories:** Array, Hash Table, Union Find

---

## Problem Statement

Given an unsorted array of integers `nums`, return the length of the longest consecutive elements sequence.

You must write an algorithm that runs in O(n) time.

**Example 1:**
```
Input: nums = [100,4,200,1,3,2]
Output: 4
Explanation: The longest consecutive sequence is [1,2,3,4], length 4.
```

**Example 2:**
```
Input: nums = [0,3,7,2,5,8,4,6,0,1]
Output: 9
```

**Constraints:**
- `0 <= nums.length <= 10^5`
- `-10^9 <= nums[i] <= 10^9`

---

## Company Frequency

| Company   | Frequency       | Last Reported |
|-----------|-----------------|---------------|
| Amazon    | ★★★★★ Very High | 2024          |
| Google    | ★★★★★ Very High | 2024          |
| Facebook  | ★★★★☆ High      | 2024          |
| Microsoft | ★★★★☆ High      | 2024          |
| Bloomberg | ★★★★☆ High      | 2024          |

> ⚠️ Frequency data is crowd-sourced from LeetCode Discuss, Glassdoor, and community interview reports. Treat as a signal, not a guarantee.

---

## DSA Concepts Used

- **HashSet** — O(1) membership test enables O(n) algorithm
- **Sequence start detection** — only start counting from sequence beginnings

---

## Approaches Overview

| # | Approach    | Time       | Space | When to use                   |
|---|-------------|------------|-------|-------------------------------|
| 1 | Sort        | O(n log n) | O(1)  | Simple but doesn't meet O(n)  |
| 2 | HashSet     | O(n)       | O(n)  | Required; always use          |

---

## Approach 1 — Sort

### Intuition
Sort the array. Scan linearly — if `nums[i] == nums[i-1]+1`, extend current sequence. Skip duplicates.

### Complexity
- **Time:** O(n log n)
- **Space:** O(1)

### Code
```go
func longestConsecutiveBrute(nums []int) int {
    if len(nums) == 0 { return 0 }
    sort.Ints(nums)
    maxLen, curr := 1, 1
    for i := 1; i < len(nums); i++ {
        if nums[i] == nums[i-1] { continue }
        if nums[i] == nums[i-1]+1 { curr++ } else { curr = 1 }
        if curr > maxLen { maxLen = curr }
    }
    return maxLen
}
```

### Dry Run
`[100,4,200,1,3,2]` sorted: `[1,2,3,4,100,200]`. Scan: 1→2→3→4 (len 4), reset, 100, 200. maxLen=4.

---

## Approach 2 — HashSet (Optimal)

### Intuition
Put all numbers in a set. For each number `n`, only start counting if `n-1` is NOT in the set (meaning `n` is the start of a sequence). Then count how long `n, n+1, n+2, ...` continues.

Each number is visited at most twice — once as a start candidate and once while extending a sequence. Total: O(n).

### Algorithm
1. Build `numSet`.
2. For each `n` in numSet:
   - If `n-1` not in numSet: it's a sequence start.
   - Count length: extend `n+1, n+2, ...` until not in set.
   - Update maxLen.

### Complexity
- **Time:** O(n) — each element visited at most twice.
- **Space:** O(n)

### Code
```go
func longestConsecutive(nums []int) int {
    numSet := make(map[int]bool)
    for _, n := range nums { numSet[n] = true }
    maxLen := 0
    for n := range numSet {
        if !numSet[n-1] { // sequence start
            curr := n; length := 1
            for numSet[curr+1] { curr++; length++ }
            if length > maxLen { maxLen = length }
        }
    }
    return maxLen
}
```

### Dry Run
`[100,4,200,1,3,2]`, numSet={1,2,3,4,100,200}:

| n   | n-1 in set? | length |
|-----|------------|--------|
| 100 | no         | 1      |
| 1   | no         | count 1,2,3,4 → 4 |
| 200 | no         | 1      |
| 2,3,4 | yes (2→1, 3→2, 4→3) | skip |

maxLen = 4 ✓

---

## Key Takeaways
- Key insight: **only count from sequence starts** (where `n-1` is absent).
- This avoids restarting counts mid-sequence → each element counted at most twice → O(n).
- Duplicates are handled automatically by the HashSet (they just don't re-trigger counting).

---

## Related Problems
- LeetCode #298 — Binary Tree Longest Consecutive Sequence
- LeetCode #549 — Binary Tree Longest Consecutive Sequence II
