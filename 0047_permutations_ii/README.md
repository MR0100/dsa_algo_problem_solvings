# 0047 — Permutations II

> LeetCode #47 · Difficulty: Medium
> **Categories:** Array, Backtracking

---

## Problem Statement

Given a collection of numbers, `nums`, that might contain duplicates, return all possible unique permutations in any order.

**Example 1**
```
Input:  nums = [1,1,2]
Output: [[1,1,2],[1,2,1],[2,1,1]]
```

**Example 2**
```
Input:  nums = [1,2,3]
Output: [[1,2,3],[1,3,2],[2,1,3],[2,3,1],[3,1,2],[3,2,1]]
```

**Constraints**
- `1 <= nums.length <= 8`
- `-10 <= nums[i] <= 10`

---

## Company Frequency

| Company   | Frequency       | Last Reported |
|-----------|-----------------|---------------|
| Amazon    | ★★★★☆ High      | 2024          |
| Google    | ★★★★☆ High      | 2024          |
| Microsoft | ★★★★☆ High      | 2024          |
| Meta      | ★★★☆☆ Medium    | 2023          |
| Bloomberg | ★★★☆☆ Medium    | 2023          |
| Apple     | ★★☆☆☆ Low       | 2022          |

> ⚠️ Frequency data is crowd-sourced from LeetCode Discuss, Glassdoor, and community interview reports. Treat as a signal, not a guarantee.

---

## DSA Concepts Used

- **Backtracking** — same as #46.
- **Skip-Duplicate Pruning** — sort first; at each recursion level, skip a candidate if it has the same value as the previous candidate AND the previous candidate was already backtracked at this level (`!visited[i-1]`).

---

## Approaches Overview

| # | Approach | Time | Space | When to use |
|---|----------|------|-------|-------------|
| 1 | Backtracking + Map Dedup | O(n! × n) | O(n! × n) | Simple; wastes space |
| 2 | Backtracking + Skip-Dup Pruning ✅ | O(n! × n) best case | O(n) | Standard optimal answer |

---

## Approach 1 — Backtracking + Map Dedup (Brute Force)

### Intuition
Generate all permutations (same as #46) and store unique ones in a string-keyed map. Correct but uses O(n! × n) extra space.

### Complexity
- **Time:** O(n! × n).
- **Space:** O(n! × n).

### Code
```go
func bruteForce(nums []int) [][]int {
    seen := map[string]bool{}
    var result [][]int
    visited := make([]bool, len(nums))
    var bt func(path []int)
    bt = func(path []int) {
        if len(path) == len(nums) {
            key := fmt.Sprint(path)
            if !seen[key] {
                seen[key] = true
                tmp := make([]int, len(nums))
                copy(tmp, path)
                result = append(result, tmp)
            }
            return
        }
        for i := 0; i < len(nums); i++ {
            if !visited[i] {
                visited[i] = true
                bt(append(path, nums[i]))
                visited[i] = false
            }
        }
    }
    bt(nil)
    return result
}
```

### Dry Run — `nums = [1,1,2]` (no sort; dedup via `seen` map)
| Full path reached | `key` | In `seen`? | Action |
|-------------------|-------|------------|--------|
| [1,1,2] (idx 0,1,2) | "[1 1 2]" | no | add → record [1,1,2] |
| [1,2,1] (idx 0,2,1) | "[1 2 1]" | no | add → record [1,2,1] |
| [1,1,2] (idx 1,0,2) | "[1 1 2]" | yes | skip (duplicate) |
| [1,2,1] (idx 1,2,0) | "[1 2 1]" | yes | skip |
| [2,1,1] (idx 2,0,1) | "[2 1 1]" | no | add → record [2,1,1] |
| [2,1,1] (idx 2,1,0) | "[2 1 1]" | yes | skip |

All 3! = 6 leaves are visited; the map keeps only the 3 unique ones.
Result: [[1,1,2],[1,2,1],[2,1,1]] count=3 ✓

---

## Approach 2 — Backtracking with Skip-Dup Pruning (Recommended ✅)

### Intuition
Sort `nums`. At each recursion level, if we encounter `nums[i] == nums[i-1]` and `!visited[i-1]`, skip `nums[i]`.

**Why `!visited[i-1]`?**
- `visited[i-1] == true` means `nums[i-1]` is currently in our path (we're building `[..., nums[i-1], ...]`). In this case, choosing `nums[i]` (same value) at the next level gives a different position — allowed.
- `visited[i-1] == false` means `nums[i-1]` was chosen and then backtracked at this level. Choosing `nums[i]` (same value) here would start exactly the same subtree as `nums[i-1]` did — duplicate! Skip.

### Algorithm
```
sort(nums)
bt(path, visited):
  if len(path)==n: record; return
  for i=0 to n-1:
    if visited[i]: continue
    if i>0 and nums[i]==nums[i-1] and NOT visited[i-1]: continue  // skip dup
    visited[i]=true; bt(path+[nums[i]]); visited[i]=false
```

### Complexity
- **Time:** O(n! × n) worst (no dups); significantly reduced with dups.
- **Space:** O(n) — visited + recursion stack.

### Code
```go
func backtracking(nums []int) [][]int {
    sort.Ints(nums); n := len(nums)
    var result [][]int; visited := make([]bool, n)
    var bt func(path []int)
    bt = func(path []int) {
        if len(path) == n {
            tmp := make([]int, n); copy(tmp, path)
            result = append(result, tmp); return
        }
        for i := 0; i < n; i++ {
            if visited[i] { continue }
            if i > 0 && nums[i] == nums[i-1] && !visited[i-1] { continue }
            visited[i] = true; bt(append(path, nums[i])); visited[i] = false
        }
    }
    bt(nil); return result
}
```

### Dry Run — `nums = [1,1,2]` sorted
```
bt([], vis=[F,F,F]):
  i=0 (1): vis=[T,F,F]; bt([1]):
    i=0: visited[0]=true → skip
    i=1 (1): vis=[T,T,F]; bt([1,1]):
      i=0,1: visited → skip
      i=2 (2): bt([1,1,2]) → record ✓
    i=2 (2): bt([1,2]):
      i=0: visited → skip
      i=1 (1): vis=[T,F,T]; bt([1,2,1]) → record ✓
  i=1 (1): i>0, nums[1]==nums[0]=1, !visited[0]=true → SKIP
  i=2 (2): vis=[F,F,T]; bt([2]):
    i=0 (1): bt([2,1]):
      i=1 (1): bt([2,1,1]) → record ✓
    i=1 (1): i>0, nums[1]==nums[0], !visited[0] → SKIP
Result: [[1,1,2],[1,2,1],[2,1,1]] count=3 ✓
```

---

## Key Takeaways

- **`!visited[i-1]` is the entire dedup logic** — a subtle one-line guard eliminates all duplicates without any extra data structure.
- **Must sort first** — `nums[i]==nums[i-1]` is the condition; this is only reliable if duplicates are adjacent (i.e., after sorting).
- **Comparison to #40 (Combination Sum II)** — in #40 the skip condition is `i > start && nums[i] == nums[i-1]` (no visited array). In #47 we use a visited array and skip when `!visited[i-1]`. Both avoid revisiting the same value at the same decision level.

---

## Related Problems

- LeetCode #46 — Permutations (distinct elements; no dedup needed)
- LeetCode #40 — Combination Sum II (skip-dup in subset problems)
- LeetCode #90 — Subsets II (same skip-dup idea for subsets)
