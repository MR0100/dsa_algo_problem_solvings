# 0039 — Combination Sum

> LeetCode #39 · Difficulty: Medium
> **Categories:** Array, Backtracking

---

## Problem Statement

Given an array of **distinct** integers `candidates` and a target integer `target`, return a list of all **unique combinations** of `candidates` where the chosen numbers sum to `target`. You may return the combinations in **any order**.

The **same** number may be chosen from `candidates` an **unlimited number of times**. Two combinations are unique if the frequency of at least one of the chosen numbers is different.

The test cases are generated such that the number of unique combinations that sum up to `target` is less than `150` combinations for the given input.

**Example 1**
```
Input:  candidates = [2,3,6,7], target = 7
Output: [[2,2,3],[7]]
```

**Example 2**
```
Input:  candidates = [2,3,5], target = 8
Output: [[2,2,2,2],[2,3,3],[3,5]]
```

**Example 3**
```
Input:  candidates = [2], target = 1
Output: []
```

**Constraints**
- `1 <= candidates.length <= 30`
- `2 <= candidates[i] <= 40`
- All elements of `candidates` are **distinct**.
- `1 <= target <= 40`

---

## Company Frequency

| Company   | Frequency       | Last Reported |
|-----------|-----------------|---------------|
| Amazon    | ★★★★★ Very High | 2024          |
| Google    | ★★★★★ Very High | 2024          |
| Microsoft | ★★★★☆ High      | 2024          |
| Meta      | ★★★★☆ High      | 2024          |
| Bloomberg | ★★★★☆ High      | 2023          |
| Apple     | ★★★☆☆ Medium    | 2023          |
| Adobe     | ★★★☆☆ Medium    | 2023          |

> ⚠️ Frequency data is crowd-sourced from LeetCode Discuss, Glassdoor, and community interview reports. Treat as a signal, not a guarantee.

---

## DSA Concepts Used

- **Backtracking** — explore combinations by deciding at each step to include the current candidate again (with the same `i`) or move to the next (`i+1`).
- **Pruning** — sorting candidates and breaking early when `candidates[i] > remaining` eliminates all impossible branches.

---

## Approaches Overview

| # | Approach | Time | Space | When to use |
|---|----------|------|-------|-------------|
| 1 | Backtracking (no sort) | O(N^(T/M) * N) | O(T/M) | Without sorting; may do redundant work |
| 2 | Backtracking with Sort + Pruning ✅ | O(N^(T/M)) | O(T/M) | Standard interview answer; cleaner pruning |

N = len(candidates), T = target, M = min(candidates).

---

## Approach 1 — Backtracking (Brute Force)

### Intuition
Try each candidate at every position. Pass `start` to avoid generating permutations of the same multiset (e.g., `[2,3]` and `[3,2]` are the same combination). Allow `i` instead of `i+1` to reuse the same candidate.

### Complexity
- **Time:** O(N^(T/M)) — branching factor N, depth T/M.
- **Space:** O(T/M) — recursion depth.

---

## Approach 2 — Backtracking with Sort + Pruning (Recommended ✅)

### Intuition
Sort candidates first. During backtracking, when `candidates[i] > remaining`, all subsequent candidates are also too large (sorted order), so `break` immediately. This prunes branches before exploring them.

- Use `start` to avoid counting the same combination in different orders.
- Recurse with same index `i` (not `i+1`) to allow unlimited reuse.

### Algorithm
```
sort(candidates)
bt(start, remaining, path):
  if remaining == 0: record path; return
  for i = start to len-1:
    if candidates[i] > remaining: break  // pruning
    bt(i, remaining-candidates[i], path+[candidates[i]])
```

### Complexity
- **Time:** O(N^(T/M)) — same asymptotic, but with far fewer explorations in practice.
- **Space:** O(T/M) — maximum recursion depth.

### Code
```go
func backtracking(candidates []int, target int) [][]int {
    sort.Ints(candidates)
    var result [][]int
    var bt func(start, remaining int, path []int)
    bt = func(start, remaining int, path []int) {
        if remaining == 0 {
            tmp := make([]int, len(path)); copy(tmp, path)
            result = append(result, tmp); return
        }
        for i := start; i < len(candidates); i++ {
            if candidates[i] > remaining { break }
            bt(i, remaining-candidates[i], append(path, candidates[i]))
        }
    }
    bt(0, target, nil)
    return result
}
```

### Dry Run — `candidates = [2,3,6,7]`, `target = 7`
```
bt(0, 7, []):
  i=0 (2): bt(0, 5, [2]):
    i=0 (2): bt(0, 3, [2,2]):
      i=0 (2): bt(0, 1, [2,2,2]):
        i=0 (2): 2>1 → break
      i=1 (3): bt(1, 0, [2,2,3]): remaining=0 → record [2,2,3]
      i=2 (6): 6>3 → break
    i=1 (3): bt(1, 2, [2,3]): 3>2 → break
    i=2 (6): 6>5 → break
  i=1 (3): bt(1, 4, [3]):
    i=1 (3): bt(1, 1, [3,3]): 3>1 → break
    i=2 (6): 6>4 → break
  i=2 (6): bt(2, 1, [6]): 6>1 → break
  i=3 (7): bt(3, 0, [7]): remaining=0 → record [7]

Result: [[2,2,3],[7]] ✓
```

---

## Key Takeaways

- **`i` not `i+1` for unlimited reuse** — in #39 (unlimited reuse) recurse with `i`; in #40 (each used once) recurse with `i+1`. This one-index difference is the entire distinction.
- **`start` prevents permutations** — without `start`, `[2,3]` and `[3,2]` would both appear. By only allowing candidates at index ≥ `start`, we enforce lexicographic (non-decreasing) order within each combination.
- **Sort before break-pruning** — the `break` pruning only works because candidates are sorted. Without sorting, a large candidate early in the list wouldn't prevent exploring smaller candidates after it.
- **Copy the path before recording** — `append(path, ...)` may share underlying memory; always `copy` before appending to `result`.

---

## Related Problems

- LeetCode #40 — Combination Sum II (each element used at most once; duplicates in input)
- LeetCode #216 — Combination Sum III (exactly k numbers, values 1–9)
- LeetCode #377 — Combination Sum IV (count permutations summing to target)
