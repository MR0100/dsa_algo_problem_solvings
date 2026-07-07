# 0090 — Subsets II

> LeetCode #90 · Difficulty: Medium
> **Categories:** Array, Backtracking, Bit Manipulation

---

## Problem Statement

Given an integer array `nums` that **may contain duplicates**, return all possible subsets (the power set).

The solution set **must not** contain duplicate subsets. Return the solution in **any order**.

**Example 1:**
```
Input: nums = [1,2,2]
Output: [[],[1],[1,2],[1,2,2],[2],[2,2]]
```

**Example 2:**
```
Input: nums = [0]
Output: [[],[0]]
```

**Constraints:**
- `1 <= nums.length <= 10`
- `-10 <= nums[i] <= 10`

---

## Company Frequency

| Company   | Frequency      | Last Reported |
|-----------|----------------|---------------|
| Amazon    | ★★★☆☆ Medium   | 2024          |
| Facebook  | ★★★☆☆ Medium   | 2023          |
| Google    | ★★☆☆☆ Low      | 2023          |

> ⚠️ Frequency data is crowd-sourced from LeetCode Discuss, Glassdoor, and community interview reports. Treat as a signal, not a guarantee.

---

## DSA Concepts Used

- **Backtracking + Sort + Skip-Duplicate Guard** — sort first; skip `nums[i]` if `nums[i]==nums[i-1]` at the same recursion level. See [`/dsa/backtracking.md`](/dsa/backtracking.md)
- **Cascading** — iterative doubling with controlled extension range for duplicates.

---

## Approaches Overview

| # | Approach | Time | Space | When to use |
|---|----------|------|-------|-------------|
| 1 | Backtracking + skip-dup guard | O(2^n × n) | O(n) | Standard; generalises to all dup-aware backtracking |
| 2 | Cascading with dup detection | O(2^n × n) | O(2^n × n) | Iterative; no recursion |

---

## Approach 1 — Backtracking with Skip-Duplicate Guard

### Intuition
Sort first so duplicates are adjacent. During backtracking, at each recursion level, if we're about to pick `nums[i]` and `nums[i] == nums[i-1]` and `i > start` (the previous element at this level was already tried), skip `nums[i]`. This avoids generating the same subset twice from different positions.

**Critical: `i > start`, not `i > 0`.** The condition `i > 0 && nums[i]==nums[i-1]` would also skip the *first* use of a duplicate, but we want to allow it as the starting element of this recursive call. `i > start` allows the first element at each level but skips duplicates at the same level.

### Algorithm
1. `sort(nums)`.
2. `bt(start, path)`:
   - Record `path`.
   - For `i = start` to `n-1`:
     - If `i > start && nums[i] == nums[i-1]`: `continue` (skip duplicate).
     - `bt(i+1, path+[nums[i]])`.

### Complexity
- **Time:** O(2^n × n) — 2^n subsets in worst case (no dups), each copied O(n).
- **Space:** O(n) — recursion depth.

### Code
```go
func subsetsWithDup(nums []int) [][]int {
    sort.Ints(nums)
    var result [][]int
    var bt func(start int, path []int)
    bt = func(start int, path []int) {
        tmp := make([]int, len(path))
        copy(tmp, path)
        result = append(result, tmp)
        for i := start; i < len(nums); i++ {
            if i > start && nums[i] == nums[i-1] { continue }
            bt(i+1, append(path, nums[i]))
        }
    }
    bt(0, nil)
    return result
}
```

### Dry Run (nums=[1,2,2] after sort)

```
bt(0, [])       → record []
  i=0 (1): bt(1, [1])  → record [1]
    i=1 (2): bt(2, [1,2])  → record [1,2]
      i=2 (2): bt(3, [1,2,2])  → record [1,2,2]
    i=2 (2): i>start(1) && nums[2]==nums[1]=2 → SKIP
  i=1 (2): bt(2, [2])  → record [2]
    i=2 (2): bt(3, [2,2])  → record [2,2]
  i=2 (2): i>start(0) && nums[2]==nums[1]=2 → SKIP
```

Output: `[[], [1], [1,2], [1,2,2], [2], [2,2]]` — 6 subsets ✓

---

## Approach 2 — Cascading with Duplicate Detection

### Intuition
Extend the cascading approach (#78): when a duplicate element is encountered, only extend the subsets added in the **previous** step (not all existing subsets). This mirrors what the skip-dup guard does in backtracking.

Track `startIdx` = the index where the previous round's new subsets begin. For a duplicate element, only extend from `startIdx`; for a new element, extend from 0 (all existing subsets).

### Algorithm
1. `result = [[]], startIdx = 0`.
2. For each `i, num` in `nums`:
   - `start = 0` if new element, `startIdx` if duplicate.
   - `startIdx = len(result)`.
   - For `j = start` to `startIdx-1`: append `result[j] + [num]`.

### Complexity
- **Time:** O(2^n × n)
- **Space:** O(2^n × n)

### Code
```go
func subsetsWithDupCascading(nums []int) [][]int {
    sort.Ints(nums)
    result := [][]int{{}}
    startIdx := 0
    for i, num := range nums {
        start := 0
        if i > 0 && nums[i] == nums[i-1] { start = startIdx }
        startIdx = len(result)
        for j := start; j < startIdx; j++ {
            newSub := make([]int, len(result[j])+1)
            copy(newSub, result[j]); newSub[len(result[j])] = num
            result = append(result, newSub)
        }
    }
    return result
}
```

### Dry Run (nums=[1,2,2] after sort)

| element | start | startIdx (before) | new subsets added | result |
|---------|-------|-------------------|-------------------|--------|
| 1 | 0 | 0→1 | [1] | [[], [1]] |
| 2 (first) | 0 | 1→2 | [2],[1,2] | [[], [1], [2], [1,2]] |
| 2 (dup) | 1 | 2→4 | only extend from idx≥1: [1]+[2]→[1,2,2]? No: from idx=1(startIdx before=2): extend result[1]=[1]→[1,2,2]? wait: startIdx=2 means "previous round added result[1] and result[2]=[2] and result[3]=[1,2]" |

Let me retrace: after element `2` (first), result = `[[], [1], [2], [1,2]]`, startIdx=2 (we set startIdx=len(result)=4 at this step... let me read the code carefully).

At i=1 (first 2): `start=0`, `startIdx=len(result)=1`, then extend j=0..0 → add `[2]`, `[1,2]`. Now result=`[[],[1],[2],[1,2]]`, startIdx=1 set before extension.

At i=2 (second 2): dup so `start=startIdx=1`, `startIdx=len(result)=4`, extend j=1..3 → add `[1]+[2]=[1,2]`? No: extend `result[1]=[1]`→`[1,2]`, `result[2]=[2]`→`[2,2]`, `result[3]=[1,2]`→`[1,2,2]`.

But that would add `[1,2]` again! Actually `start=startIdx=1` (the startIdx saved from the previous step, which was 1), so j=1,2,3 → new subsets from those: `[1,2]` (from [1]), `[2,2]` (from [2]), `[1,2,2]` (from [1,2]).

Hmm, `[1,2]` would be a duplicate. Let me re-check the code: `startIdx` before this iteration = 1 (saved at i=1). So j=1..3 adds `[1]+2=[1,2]` (dup!), `[2]+2=[2,2]`, `[1,2]+2=[1,2,2]`.

Actually looking again: the key insight is `startIdx` from i=1 round was 1, meaning "subsets that existed before we processed the first '2'". The subsets added by the first '2' are at indices 2 and 3. So for the second '2', we should extend only indices 2..3. This means `startIdx` at i=2 should be 2, not 1.

Looking at the code: at i=1, `startIdx = len(result)` is set BEFORE extending → `startIdx=1` (result has 2 elements: `[[],[1]]`). After extending, result has 4 elements. At i=2, we read `startIdx=1` → `start=1` → extend j=1..3. This adds `[1]+2,[2]+2,[1,2]+2 = [1,2],[2,2],[1,2,2]`. The `[1,2]` would be a duplicate.

The verified output shows count=6 which is correct. Let me re-examine: result starts as `[[]]`. After i=0 (num=1): startIdx=0, then startIdx=len=1, extend j=0: add `[1]`. result=`[[],[1]]`. After i=1 (num=2, new): start=0, startIdx=len=2, extend j=0,1: add `[2],[1,2]`. result=`[[],[1],[2],[1,2]]`. startIdx before set to 2 (not updated — the code sets startIdx before the extension loop, so startIdx=2 at end of i=1 step since `startIdx=len(result)=2` and result had 2 elements before extending).

Wait: `startIdx = len(result)` is set BEFORE the inner loop. At i=1 start of step: startIdx=1 (saved from prev), start=0, then `startIdx=len(result)=2` (result has [[],[1]]). Then extend j=0..1: add [2],[1,2]. result=[[],[1],[2],[1,2]].

At i=2 (dup): start=startIdx=2 (saved from i=1), `startIdx=len(result)=4`. Extend j=2..3: result[2]=[2]→[2,2], result[3]=[1,2]→[1,2,2]. result=`[[],[1],[2],[1,2],[2,2],[1,2,2]]` = 6 subsets ✓. No duplicates!

---

## Key Takeaways
- **Sort first** is mandatory for both approaches — duplicates must be adjacent.
- In backtracking: `i > start` (not `i > 0`) is the critical condition. `i > 0` would incorrectly skip the very first pick of a duplicate value.
- In cascading: only extend subsets from the previous round when a duplicate is encountered.
- This skip-dup pattern applies to Combinations (#40), Permutations (#47), Subsets (#90) — any backtracking with duplicates.

---

## Related Problems
- LeetCode #78 — Subsets (no duplicates)
- LeetCode #40 — Combination Sum II (combinations with duplicates; same skip-dup guard)
- LeetCode #47 — Permutations II (permutations with duplicates; different skip condition)
