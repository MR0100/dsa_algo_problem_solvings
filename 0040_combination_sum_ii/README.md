# 0040 — Combination Sum II

> LeetCode #40 · Difficulty: Medium
> **Categories:** Array, Backtracking

---

## Problem Statement

Given a collection of candidate numbers (`candidates`) and a target number (`target`), find all unique combinations in `candidates` where the candidate numbers sum to `target`.

Each number in `candidates` may only be used **once** in the combination.

**Note:** The solution set must not contain duplicate combinations.

**Example 1**
```
Input:  candidates = [10,1,2,7,6,1,5], target = 8
Output: [[1,1,6],[1,2,5],[1,7],[2,6]]
```

**Example 2**
```
Input:  candidates = [2,5,2,1,2], target = 5
Output: [[1,2,2],[5]]
```

**Constraints**
- `1 <= candidates.length <= 100`
- `1 <= candidates[i] <= 50`
- `1 <= target <= 30`

---

## Company Frequency

| Company   | Frequency       | Last Reported |
|-----------|-----------------|---------------|
| Amazon    | ★★★★☆ High      | 2024          |
| Google    | ★★★★☆ High      | 2024          |
| Microsoft | ★★★★☆ High      | 2024          |
| Meta      | ★★★☆☆ Medium    | 2023          |
| Bloomberg | ★★★☆☆ Medium    | 2023          |
| Apple     | ★★★☆☆ Medium    | 2023          |
| Adobe     | ★★☆☆☆ Low       | 2022          |

> ⚠️ Frequency data is crowd-sourced from LeetCode Discuss, Glassdoor, and community interview reports. Treat as a signal, not a guarantee.

---

## DSA Concepts Used

- **Backtracking** — same recursive structure as #39, but advance `i+1` since each element is used at most once.
- **Skip-Duplicate Pruning** — after sorting, skip candidates at the same recursion level that equal the previous candidate (`i > start && candidates[i] == candidates[i-1]`). This eliminates duplicate combinations without extra data structures.

---

## Approaches Overview

| # | Approach | Time | Space | When to use |
|---|----------|------|-------|-------------|
| 1 | Backtracking + Map Dedup | O(2^N × N) | O(2^N × N) | Simple to write; wastes space on dedup map |
| 2 | Backtracking + Skip-Dup Pruning ✅ | O(2^N) | O(N) | Optimal; no extra data structures |

---

## Approach 1 — Backtracking + Map Dedup (Brute Force)

### Intuition
Run standard backtracking (each element used at most once → recurse with `i+1`). Store results in a map keyed by the sorted combination string to deduplicate. Correct but uses O(2^N × N) extra space.

### Complexity
- **Time:** O(2^N × N) — 2^N subsets, each takes O(N) to stringify for the map.
- **Space:** O(2^N × N) — the dedup map.

### Code
```go
func bruteForce(candidates []int, target int) [][]int {
    sort.Ints(candidates)
    seen := map[string]bool{}
    var result [][]int
    var bt func(start, remaining int, path []int)
    bt = func(start, remaining int, path []int) {
        if remaining == 0 {
            // encode the combo as a string key for dedup
            key := fmt.Sprint(path)
            if !seen[key] {
                seen[key] = true
                tmp := make([]int, len(path))
                copy(tmp, path)
                result = append(result, tmp)
            }
            return
        }
        for i := start; i < len(candidates); i++ {
            if candidates[i] > remaining {
                break
            }
            bt(i+1, remaining-candidates[i], append(path, candidates[i]))
        }
    }
    bt(0, target, nil)
    return result
}
```

### Dry Run — `candidates = [10,1,2,7,6,1,5]` sorted → `[1,1,2,5,6,7,10]`, `target = 8`

No skip-dup guard here — every branch is explored (advancing `i+1`), and the map `seen` filters out combinations that stringify to a key already recorded.

| Step | Call | Action | `seen` after | `result` after |
|------|------|--------|--------------|----------------|
| 1 | reach `remaining=0` via `[1,1,6]` | key `[1 1 6]` new → record | `{[1 1 6]}` | `[[1 1 6]]` |
| 2 | reach `remaining=0` via `[1,2,5]` | key `[1 2 5]` new → record | `+[1 2 5]` | `[[1 1 6] [1 2 5]]` |
| 3 | reach `remaining=0` via `[1,7]` | key `[1 7]` new → record | `+[1 7]` | `[[1 1 6] [1 2 5] [1 7]]` |
| 4 | second `1` branch reaches `[1,2,5]` again (starting from the 2nd `1`) | key `[1 2 5]` already in `seen` → skip | unchanged | unchanged |
| 5 | second `1` branch reaches `[1,7]` again | key `[1 7]` already in `seen` → skip | unchanged | unchanged |
| 6 | reach `remaining=0` via `[2,6]` | key `[2 6]` new → record | `+[2 6]` | `[[1 1 6] [1 2 5] [1 7] [2 6]]` |

Result: `[[1 1 6] [1 2 5] [1 7] [2 6]]` ✓ — the map absorbs the duplicate `[1 2 5]` and `[1 7]` that the two identical `1`s would otherwise produce.

---

## Approach 2 — Backtracking with Skip-Duplicate Pruning (Recommended ✅)

### Intuition
Sort candidates. During backtracking, at each recursion level (same `start`), skip a candidate if it equals the previous candidate at the same level. This prevents exploring two branches that would produce the same combination.

**Why `i > start` matters:**
- `i == start`: this is the first choice at this level; we must consider it.
- `i > start && candidates[i] == candidates[i-1]`: we already explored a branch starting with `candidates[i-1]`; starting with `candidates[i]` (same value) would produce an identical sub-tree. Skip.

### Algorithm
```
sort(candidates)
bt(start, remaining, path):
  if remaining == 0: record path; return
  for i = start to len-1:
    if candidates[i] > remaining: break
    if i > start and candidates[i] == candidates[i-1]: continue  // skip dup
    bt(i+1, remaining-candidates[i], path+[candidates[i]])
```

### Complexity
- **Time:** O(2^N) — each element is either included or excluded.
- **Space:** O(N) — recursion depth (at most N deep since each element used once).

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
            if i > start && candidates[i] == candidates[i-1] { continue }
            bt(i+1, remaining-candidates[i], append(path, candidates[i]))
        }
    }
    bt(0, target, nil)
    return result
}
```

### Dry Run — `candidates = [10,1,2,7,6,1,5]` sorted → `[1,1,2,5,6,7,10]`, `target = 8`
```
bt(0, 8, []):
  i=0 (1): bt(1, 7, [1]):
    i=1 (1): bt(2, 6, [1,1]):
      i=2 (2): bt(3, 4, [1,1,2]):
        i=3 (5): 5>4 → break
      i=3 (5): bt(4, 1, [1,1,5]): 6>1 → break
      i=4 (6): bt(5, 0, [1,1,6]): record [1,1,6] ✓
      i=5 (7): 7>6 → break
    i=2 (2): bt(3, 5, [1,2]):
      i=3 (5): bt(4, 0, [1,2,5]): record [1,2,5] ✓
      i=4 (6): 6>5 → break
    i=3 (5): bt(4, 2, [1,5]): 6>2 → break
    i=4 (6): bt(5, 1, [1,6]): 7>1 → break
    i=5 (7): bt(6, 0, [1,7]): record [1,7] ✓
    i=6 (10): 10>7 → break
  i=1 (1): skip (i>0 and candidates[1]==candidates[0])
  i=2 (2): bt(3, 6, [2]):
    i=3 (5): bt(4, 1, [2,5]): 6>1 → break
    i=4 (6): bt(5, 0, [2,6]): record [2,6] ✓
    i=5 (7): 7>6 → break
  ...

Result: [[1,1,6],[1,2,5],[1,7],[2,6]] ✓
```

---

## Key Takeaways

- **`i+1` vs `i`** — in #39 each candidate can be reused (recurse with `i`); here each is used at most once (recurse with `i+1`). This single difference is the #39 vs #40 distinction.
- **`i > start` is the critical guard** — skipping `i` when `candidates[i] == candidates[i-1]` without this guard would incorrectly skip the first occurrence of a duplicate, preventing combinations like `[1,1,6]` from being found (which requires both 1s).
- **Sort first** — the skip-dup rule only works because duplicates are adjacent after sorting. Without sorting, you'd need a hash set.
- **Generalises to kSum** — the "skip same value at same level" rule applies to 3Sum, 4Sum, and all kSum problems to eliminate duplicate triplets/quadruplets.

---

## Related Problems

- LeetCode #39 — Combination Sum (unlimited reuse; no duplicates in input)
- LeetCode #15 — 3Sum (same skip-duplicate pruning in sorted array)
- LeetCode #18 — 4Sum (same pattern, two outer loops)
- LeetCode #216 — Combination Sum III (k numbers from 1–9)
