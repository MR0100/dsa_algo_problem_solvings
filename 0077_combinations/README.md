# 0077 — Combinations

> LeetCode #77 · Difficulty: Medium
> **Categories:** Backtracking, Combinatorics

---

## Problem Statement

Given two integers `n` and `k`, return all possible combinations of `k` numbers chosen from the range `[1, n]`.

You may return the answer in **any order**.

**Example 1:**
```
Input: n = 4, k = 2
Output: [[1,2],[1,3],[1,4],[2,3],[2,4],[3,4]]
```

**Example 2:**
```
Input: n = 1, k = 1
Output: [[1]]
```

**Constraints:**
- `1 <= n <= 20`
- `1 <= k <= n`

---

## Company Frequency

| Company   | Frequency      | Last Reported |
|-----------|----------------|---------------|
| Amazon    | ★★★☆☆ Medium   | 2024          |
| Google    | ★★★☆☆ Medium   | 2024          |
| Microsoft | ★★☆☆☆ Low      | 2023          |

> ⚠️ Frequency data is crowd-sourced from LeetCode Discuss, Glassdoor, and community interview reports. Treat as a signal, not a guarantee.

---

## DSA Concepts Used

- **Backtracking** — recursively build combinations; prune impossible branches. See [`/dsa/backtracking.md`](/dsa/backtracking.md)
- **Combinatorics** — generating C(n,k) combinations systematically.

---

## Approaches Overview

| # | Approach | Time | Space | When to use |
|---|----------|------|-------|-------------|
| 1 | Backtracking (with pruning) | O(C(n,k) × k) | O(k) | Standard interview answer |
| 2 | Iterative (lexicographic) | O(C(n,k) × k) | O(k) | When recursion overhead matters |

---

## Approach 1 — Backtracking

### Intuition
Choose numbers in strictly increasing order (from `start` to `n`) to avoid generating duplicate combinations. Record each path once it has exactly `k` elements.

**Pruning:** if we're at position `start` and need `(k - len(path))` more elements, but only `(n - start + 1)` numbers remain, there's no point iterating past `n - (k - len(path)) + 1`. This cuts many branches.

### Algorithm
1. `bt(start, path)`:
   - If `len(path) == k`: record and return.
   - `limit = n - (k - len(path)) + 1` — pruned upper bound.
   - For `i = start` to `limit`: recurse with `i+1` and `path + [i]`.
2. Start with `bt(1, [])`.

### Complexity
- **Time:** O(C(n,k) × k) — C(n,k) combinations, copying each costs O(k).
- **Space:** O(k) — max recursion depth is `k`.

### Code
```go
func backtracking(n, k int) [][]int {
    var result [][]int
    var bt func(start int, path []int)
    bt = func(start int, path []int) {
        if len(path) == k {
            tmp := make([]int, k)
            copy(tmp, path)
            result = append(result, tmp)
            return
        }
        limit := n - (k - len(path)) + 1
        for i := start; i <= limit; i++ {
            bt(i+1, append(path, i))
        }
    }
    bt(1, nil)
    return result
}
```

### Dry Run (n=4, k=2)

```
bt(1, [])
  i=1 → bt(2, [1])
    i=2 → bt(3, [1,2]) → len==k → record [1,2]
    i=3 → bt(4, [1,3]) → record [1,3]
    i=4 → bt(5, [1,4]) → record [1,4]
  i=2 → bt(3, [2])
    i=3 → record [2,3]
    i=4 → record [2,4]
  i=3 → bt(4, [3])
    i=4 → record [3,4]
  i=4 → bt(5, [4])  [limit=4-(2-1)+1=4, so allowed]
```

Pruning example: when `path=[]` and `start=4`, `limit = 4-(2-0)+1 = 3`, so `i=4` is never tried — there aren't 2 elements left starting at 4.

Result: `[[1,2],[1,3],[1,4],[2,3],[2,4],[3,4]]` ✓

---

## Approach 2 — Iterative (Lexicographic)

### Intuition
Start with the lexicographically smallest combination `[1, 2, ..., k]`. Repeatedly find the rightmost element that can still be incremented (i.e., `nums[i] < n - k + i + 1`), increment it, and fill the rest in order.

This mimics what backtracking does but without the call stack.

### Algorithm
1. Initialize `nums = [1, 2, ..., k]`.
2. Loop:
   - Record current combination.
   - Find rightmost `i` where `nums[i] < n - k + i + 1`.
   - If none found: done.
   - `nums[i]++`; set `nums[j] = nums[j-1] + 1` for `j > i`.

### Complexity
- **Time:** O(C(n,k) × k)
- **Space:** O(k) — one current-combination array.

### Code
```go
func iterative(n, k int) [][]int {
    var result [][]int
    nums := make([]int, k)
    for i := range nums { nums[i] = i + 1 }
    for {
        tmp := make([]int, k)
        copy(tmp, nums)
        result = append(result, tmp)
        i := k - 1
        for i >= 0 && nums[i] == n-k+i+1 { i-- }
        if i < 0 { break }
        nums[i]++
        for j := i + 1; j < k; j++ { nums[j] = nums[j-1] + 1 }
    }
    return result
}
```

### Dry Run (n=4, k=2)

| nums | rightmost i can increment? | action |
|------|---------------------------|--------|
| [1,2] | i=1 (2 < 4-2+1+1=4? yes) | record; increment i=1: [1,3] |
| [1,3] | i=1 (3 < 4? yes) | record; [1,4] |
| [1,4] | i=1 (4 < 4? no); i=0 (1 < 3? yes) | record; i=0 → 2; fill j=1: [2,3] |
| [2,3] | i=1 (3 < 4? yes) | record; [2,4] |
| [2,4] | i=1 no; i=0 (2 < 3? yes) | record; [3,4] |
| [3,4] | i=1 no; i=0 (3 < 3? no); i=-1 | record; done |

---

## Key Takeaways
- The pruning bound `limit = n - (k - len(path)) + 1` is critical — it prunes all branches that can't produce a full k-element combination.
- "Choose in increasing order" = no need to track visited; naturally avoids duplicates.
- The iterative approach is O(1) extra space (excluding output) and avoids call-stack overhead.

---

## Related Problems
- LeetCode #39 — Combination Sum (combinations with repetition allowed)
- LeetCode #40 — Combination Sum II (combinations, no repetition, skip duplicates)
- LeetCode #78 — Subsets (all subsets, not just k-sized)
- LeetCode #216 — Combination Sum III (k numbers that sum to n)
