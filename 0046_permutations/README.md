# 0046 — Permutations

> LeetCode #46 · Difficulty: Medium
> **Categories:** Array, Backtracking

---

## Problem Statement

Given an array `nums` of distinct integers, return all the possible permutations. You can return the answer in **any order**.

**Example 1**
```
Input:  nums = [1,2,3]
Output: [[1,2,3],[1,3,2],[2,1,3],[2,3,1],[3,1,2],[3,2,1]]
```

**Example 2**
```
Input:  nums = [0,1]
Output: [[0,1],[1,0]]
```

**Example 3**
```
Input:  nums = [1]
Output: [[1]]
```

**Constraints**
- `1 <= nums.length <= 6`
- `-10 <= nums[i] <= 10`
- All the integers of `nums` are **unique**.

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
| Adobe     | ★★★☆☆ Medium    | 2023          |

> ⚠️ Frequency data is crowd-sourced from LeetCode Discuss, Glassdoor, and community interview reports. Treat as a signal, not a guarantee.

---

## DSA Concepts Used

- **Backtracking** — choose, recurse, unchoose. At each step we pick one unused element, recurse until the path is full, then unchoose.

---

## Approaches Overview

| # | Approach | Time | Space | When to use |
|---|----------|------|-------|-------------|
| 1 | Backtracking + Visited Array ✅ | O(n! × n) | O(n) | Cleanest to explain; explicit choice tracking |
| 2 | Swap Backtracking (In-Place) ✅ | O(n! × n) | O(n) | No visited array; slightly less space |
| 3 | Iterative Insertion | O(n! × n) | O(n! × n) | Good for understanding; less memory-efficient |

---

## Approach 1 — Backtracking with Visited Array (Recommended ✅)

### Intuition
At each recursion level, loop through all indices. If `nums[i]` hasn't been used, add it to the path, recurse, then remove it (backtrack). When the path reaches length n, record it.

### Algorithm
```
bt(path, visited):
  if len(path)==n: record path; return
  for i=0 to n-1:
    if not visited[i]:
      visited[i]=true; bt(path+[nums[i]]); visited[i]=false
```

### Complexity
- **Time:** O(n! × n) — n! permutations × O(n) copy each.
- **Space:** O(n) — visited array + recursion stack of depth n.

### Code
```go
// backtracking solves Permutations by building each permutation position by
// position, using a boolean visited array to track which elements are used.
//
// Time:  O(n! * n) — n! permutations, each costs O(n) to copy
// Space: O(n) — recursion depth n; visited array n
func backtracking(nums []int) [][]int {
    n := len(nums)
    var result [][]int
    visited := make([]bool, n)

    var bt func(path []int)
    bt = func(path []int) {
        if len(path) == n {
            tmp := make([]int, n)
            copy(tmp, path)
            result = append(result, tmp)
            return
        }
        for i := 0; i < n; i++ {
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

### Dry Run — `nums = [1,2,3]`
```
bt([], vis=[F,F,F]):
  i=0: vis=[T,F,F]; bt([1]):
    i=1: bt([1,2]):
      i=2: bt([1,2,3]) → record [1,2,3]
    i=2: bt([1,3]):
      i=1: bt([1,3,2]) → record [1,3,2]
  i=1: vis=[F,T,F]; bt([2]):
    ... → [2,1,3], [2,3,1]
  i=2: ... → [3,1,2], [3,2,1]
```

---

## Approach 2 — Swap Backtracking (In-Place)

### Intuition
At position `start`, try each index `i >= start` by swapping `nums[start]` with `nums[i]`, recursing for `start+1`, then swapping back.

No visited array needed — the elements in `nums[0..start-1]` are "fixed" and those in `nums[start..n-1]` are candidates.

### Code
```go
func swapBacktracking(nums []int) [][]int {
    var result [][]int
    var bt func(start int)
    bt = func(start int) {
        if start == len(nums) {
            tmp := make([]int, len(nums)); copy(tmp, nums)
            result = append(result, tmp); return
        }
        for i := start; i < len(nums); i++ {
            nums[start], nums[i] = nums[i], nums[start]
            bt(start+1)
            nums[start], nums[i] = nums[i], nums[start]
        }
    }
    bt(0); return result
}
```

### Complexity
- **Time:** O(n! × n) — n! leaves reached; each records a length-n copy in O(n).
- **Space:** O(n) — recursion stack depth n; swaps mutate `nums` in place, no visited array.

### Dry Run — `nums = [1,2,3]`
| Call | `nums` on entry | Loop action | Recurse | On record |
|------|-----------------|-------------|---------|-----------|
| bt(0) | [1,2,3] | i=0 swap(0,0) → [1,2,3] | bt(1) | |
| bt(1) | [1,2,3] | i=1 swap(1,1) → [1,2,3] | bt(2) | |
| bt(2) | [1,2,3] | i=2 swap(2,2) → [1,2,3] | bt(3) | |
| bt(3) | [1,2,3] | start==n | — | record [1,2,3] |
| bt(2) | [1,2,3] | i=2 swap(1,2) → [1,3,2] | bt(3) | record [1,3,2]; restore [1,2,3] |
| bt(1) | [1,2,3] | i=2 swap(0,2) → [3,2,1] | bt(1)… | eventually [3,2,1],[3,1,2]; restore |
| bt(0) | [1,2,3] | i=1 swap(0,1) → [2,1,3] | bt(1)… | eventually [2,1,3],[2,3,1]; restore |

Result (order of discovery): [1,2,3],[1,3,2],[2,1,3],[2,3,1],[3,2,1],[3,1,2] — count 6 ✓

---

## Approach 3 — Iterative Insertion

### Intuition
Build all permutations incrementally. Start with `[[]]`. For each number, insert it at every possible position (0 to len(perm)) of every existing permutation.

### Complexity
- **Time:** O(n! × n).
- **Space:** O(n! × n) — all permutations stored simultaneously.

### Code
```go
// iterative solves Permutations by iteratively inserting each number into every
// possible position of existing permutations.
//
// Time:  O(n! * n)
// Space: O(n! * n)
func iterative(nums []int) [][]int {
    result := [][]int{{}} // start with one empty permutation
    for _, num := range nums {
        var next [][]int
        for _, perm := range result {
            // insert num at every position in perm
            for pos := 0; pos <= len(perm); pos++ {
                newPerm := make([]int, len(perm)+1)
                copy(newPerm[:pos], perm[:pos])
                newPerm[pos] = num
                copy(newPerm[pos+1:], perm[pos:])
                next = append(next, newPerm)
            }
        }
        result = next
    }
    return result
}
```

### Dry Run — `nums = [1,2,3]`
| Step | `num` | `result` before | Insert positions per perm | `result` after |
|------|-------|-----------------|---------------------------|----------------|
| init | — | — | — | [[]] |
| 1 | 1 | [[]] | insert 1 into [] at pos 0 | [[1]] |
| 2 | 2 | [[1]] | into [1] at pos 0,1 | [[2,1],[1,2]] |
| 3 | 3 | [[2,1],[1,2]] | into [2,1] at 0,1,2; into [1,2] at 0,1,2 | [[3,2,1],[2,3,1],[2,1,3],[3,1,2],[1,3,2],[1,2,3]] |

Result: 6 permutations ✓

---

## Key Takeaways

- **Visited array vs swap** — visited array is O(n) extra; swap is O(1) extra but mutates the input (must restore). Both work in interviews.
- **Always copy before recording** — `append(path, x)` in Go may share the underlying slice. Always `copy` into a new slice before appending to `result`.
- **n! grows fast** — n=6: 720, n=8: 40,320, n=10: 3.6M. For n > 8, outputting all permutations becomes impractical.
- **This is the base for #47** — Permutations II adds duplicates; the skip-dup rule is the only additional mechanism needed.

---

## Related Problems

- LeetCode #47 — Permutations II (duplicates in input; skip-dup pruning)
- LeetCode #31 — Next Permutation (generate one specific permutation)
- LeetCode #60 — Permutation Sequence (k-th permutation directly)
