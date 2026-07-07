# 0078 — Subsets

> LeetCode #78 · Difficulty: Medium
> **Categories:** Array, Backtracking, Bit Manipulation

---

## Problem Statement

Given an integer array `nums` of **unique** elements, return all possible subsets (the power set).

The solution set **must not** contain duplicate subsets. Return the solution in **any order**.

**Example 1:**
```
Input: nums = [1,2,3]
Output: [[],[1],[2],[1,2],[3],[1,3],[2,3],[1,2,3]]
```

**Example 2:**
```
Input: nums = [0]
Output: [[],[0]]
```

**Constraints:**
- `1 <= nums.length <= 10`
- `-10 <= nums[i] <= 10`
- All the numbers of `nums` are **unique**.

---

## Company Frequency

| Company   | Frequency       | Last Reported |
|-----------|-----------------|---------------|
| Amazon    | ★★★★☆ High      | 2024          |
| Google    | ★★★☆☆ Medium    | 2024          |
| Facebook  | ★★★☆☆ Medium    | 2023          |
| Microsoft | ★★☆☆☆ Low       | 2023          |

> ⚠️ Frequency data is crowd-sourced from LeetCode Discuss, Glassdoor, and community interview reports. Treat as a signal, not a guarantee.

---

## DSA Concepts Used

- **Backtracking** — DFS where we record the path at every node, not just the leaves. See [`/dsa/backtracking.md`](/dsa/backtracking.md)
- **Bit Manipulation** — enumerate all 2^n bitmasks to represent subset membership.
- **Cascading** — iterative doubling: grow the result set by adding each new element to all existing subsets.

---

## Approaches Overview

| # | Approach | Time | Space | When to use |
|---|----------|------|-------|-------------|
| 1 | Backtracking | O(2^n × n) | O(n) | Most natural; good for extension to duplicates |
| 2 | Bit Manipulation | O(2^n × n) | O(n) | When n is small; elegant with bitmask tricks |
| 3 | Cascading (Iterative) | O(2^n × n) | O(2^n × n) | No recursion; easy to reason about |

---

## Approach 1 — Backtracking

### Intuition
Walk through each element and decide: include or skip. Record the current path at every level of the DFS tree — this naturally generates all 2^n subsets.

Unlike combination problems, we don't stop at a fixed depth; every node in the recursion tree is a valid subset.

### Algorithm
1. `bt(start, path)`:
   - Record a copy of `path` (this is a valid subset at any length).
   - For `i = start` to `n-1`: recurse with `bt(i+1, path + [nums[i]])`.
2. Start with `bt(0, [])`.

### Complexity
- **Time:** O(2^n × n) — 2^n subsets, each copied in O(n).
- **Space:** O(n) — recursion depth at most n.

### Code
```go
func backtracking(nums []int) [][]int {
    var result [][]int
    var bt func(start int, path []int)
    bt = func(start int, path []int) {
        tmp := make([]int, len(path))
        copy(tmp, path)
        result = append(result, tmp)
        for i := start; i < len(nums); i++ {
            bt(i+1, append(path, nums[i]))
        }
    }
    bt(0, nil)
    return result
}
```

### Dry Run (nums=[1,2,3])

```
bt(0, [])          → record []
  bt(1, [1])       → record [1]
    bt(2, [1,2])   → record [1,2]
      bt(3, [1,2,3]) → record [1,2,3]
    bt(3, [1,3])   → record [1,3]
  bt(2, [2])       → record [2]
    bt(3, [2,3])   → record [2,3]
  bt(3, [3])       → record [3]
```

Output: `[[], [1], [1,2], [1,2,3], [1,3], [2], [2,3], [3]]` — 8 subsets ✓

---

## Approach 2 — Bit Manipulation

### Intuition
For `n` elements, there are 2^n possible subsets, each uniquely identified by a bitmask of `n` bits. Bit `j` set in mask `i` means `nums[j]` is in the `i`-th subset.

Mask `0b000` = `{}`, `0b001` = `{nums[0]}`, `0b011` = `{nums[0], nums[1]}`, etc.

### Algorithm
1. For `mask = 0` to `2^n - 1`:
   - For `j = 0` to `n-1`: if bit `j` is set in `mask`, include `nums[j]`.
   - Append the subset to result.

### Complexity
- **Time:** O(2^n × n)
- **Space:** O(n) — current subset.

### Code
```go
func bitManipulation(nums []int) [][]int {
    n := len(nums)
    total := 1 << n
    result := make([][]int, 0, total)
    for mask := 0; mask < total; mask++ {
        subset := []int{}
        for j := 0; j < n; j++ {
            if mask&(1<<j) != 0 {
                subset = append(subset, nums[j])
            }
        }
        result = append(result, subset)
    }
    return result
}
```

### Dry Run (nums=[1,2,3], n=3)

| mask | binary | subset |
|------|--------|--------|
| 0 | 000 | [] |
| 1 | 001 | [1] |
| 2 | 010 | [2] |
| 3 | 011 | [1,2] |
| 4 | 100 | [3] |
| 5 | 101 | [1,3] |
| 6 | 110 | [2,3] |
| 7 | 111 | [1,2,3] |

8 subsets ✓

---

## Approach 3 — Cascading (Iterative)

### Intuition
Start with `result = [[]]`. For each element, add it to every existing subset and append the new subsets. This doubles the result size with each element.

Round 0 (before any element): `[[]]` (1 subset)
Round 1 (add 1): `[[], [1]]` (2 subsets)
Round 2 (add 2): `[[], [1], [2], [1,2]]` (4 subsets)
Round 3 (add 3): `[[], [1], [2], [1,2], [3], [1,3], [2,3], [1,2,3]]` (8 subsets)

### Algorithm
1. `result = [[]]`.
2. For each `num` in `nums`:
   - For each existing subset `sub` in `result`: append `sub + [num]` to result.

### Complexity
- **Time:** O(2^n × n) — total work across all rounds.
- **Space:** O(2^n × n) — output.

### Code
```go
func cascading(nums []int) [][]int {
    result := [][]int{{}}
    for _, num := range nums {
        n := len(result)
        for i := 0; i < n; i++ {
            newSub := make([]int, len(result[i])+1)
            copy(newSub, result[i])
            newSub[len(result[i])] = num
            result = append(result, newSub)
        }
    }
    return result
}
```

### Dry Run (nums=[1,2,3])

| After element | result |
|---------------|--------|
| (initial) | [[]] |
| 1 | [[], [1]] |
| 2 | [[], [1], [2], [1,2]] |
| 3 | [[], [1], [2], [1,2], [3], [1,3], [2,3], [1,2,3]] |

---

## Key Takeaways
- Record path at every DFS node (not just leaves) to generate all subsets, not just k-sized ones.
- Bitmask approach is clean and O(2^n × n) but only works for small n ≤ 20 due to integer size.
- Cascading is the most visual: the result literally doubles with each element.
- All three approaches produce 2^n subsets — same time complexity, different constants.

---

## Related Problems
- LeetCode #90 — Subsets II (with duplicates — sort + skip-dup guard)
- LeetCode #77 — Combinations (k-sized subsets only)
- LeetCode #784 — Letter Case Permutation (similar 2-branch DFS at each index)
