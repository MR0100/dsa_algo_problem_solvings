# 0095 — Unique Binary Search Trees II

> LeetCode #95 · Difficulty: Medium
> **Categories:** Dynamic Programming, Backtracking, Tree, Binary Search Tree

---

## Problem Statement

Given an integer `n`, return **all structurally unique BST's** (binary search trees), which has exactly `n` nodes of unique values from `1` to `n`. Return the answer in **any order**.

**Example 1:**
```
Input: n = 3
Output: [[1,null,2,null,3],[1,null,3,2],[2,1,3],[3,1,null,null,2],[3,2,null,1]]
(5 unique BSTs)
```

**Example 2:**
```
Input: n = 1
Output: [[1]]
```

**Constraints:**
- `1 <= n <= 8`

---

## Company Frequency

| Company   | Frequency      | Last Reported |
|-----------|----------------|---------------|
| Amazon    | ★★★☆☆ Medium   | 2024          |
| Google    | ★★★☆☆ Medium   | 2023          |
| Microsoft | ★★☆☆☆ Low      | 2023          |

> ⚠️ Frequency data is crowd-sourced from LeetCode Discuss, Glassdoor, and community interview reports. Treat as a signal, not a guarantee.

---

## DSA Concepts Used

- **Divide and Conquer on BST** — choose root from [start..end]; left subtree uses [start..root-1], right subtree uses [root+1..end]. See [`/dsa/tree_traversal.md`](/dsa/tree_traversal.md)
- **Catalan Numbers** — the number of unique BSTs with n nodes is the n-th Catalan number: C(n) = C(2n,n)/(n+1).
- **Memoization** — cache results for each (start, end) pair.

---

## Approaches Overview

| # | Approach | Time | Space | When to use |
|---|----------|------|-------|-------------|
| 1 | Recursion (Divide and Conquer) | O(Catalan(n) × n) | O(Catalan(n) × n) | Natural; clean |
| 2 | Memoized Recursion | O(Catalan(n) × n) | O(n² + Catalan(n) × n) | Avoids recomputing same (start,end) pairs |

---

## Approach 1 — Recursion (Divide and Conquer)

### Intuition
For a range `[start..end]`, every value `i` in `[start..end]` can be the root. The left subtree contains all BSTs over `[start..i-1]` and the right subtree contains all BSTs over `[i+1..end]`. Combine each left tree with each right tree to form a complete BST with root `i`.

**Base case:** if `start > end`, return `[nil]` (one empty tree). This is essential — it allows ranges like `[start..root-1]` when root = start to return `[nil]` rather than nothing.

### Algorithm
1. `generate(start, end)`:
   - If `start > end`: return `[nil]`.
   - For each `i = start` to `end`:
     - `leftTrees = generate(start, i-1)`.
     - `rightTrees = generate(i+1, end)`.
     - For each `(left, right)` pair: create `TreeNode{i, left, right}` and append.

### Complexity
- **Time:** O(Catalan(n) × n) — the number of trees grown is Catalan(n); for each tree, O(n) nodes created.
- **Space:** O(Catalan(n) × n) — storing all trees.

### Code
```go
func generateTrees(n int) []*TreeNode {
    if n == 0 { return nil }
    var generate func(start, end int) []*TreeNode
    generate = func(start, end int) []*TreeNode {
        if start > end { return []*TreeNode{nil} }
        var allTrees []*TreeNode
        for i := start; i <= end; i++ {
            leftTrees := generate(start, i-1)
            rightTrees := generate(i+1, end)
            for _, left := range leftTrees {
                for _, right := range rightTrees {
                    allTrees = append(allTrees, &TreeNode{Val: i, Left: left, Right: right})
                }
            }
        }
        return allTrees
    }
    return generate(1, n)
}
```

### Dry Run (n=3)

`generate(1,3)`:
- `i=1`: left=`gen(1,0)`=[nil], right=`gen(2,3)`:
  - `i=2`: left=[nil], right=`gen(3,3)`:
    - `i=3`: left=[nil], right=[nil] → tree `{3,nil,nil}`.
  - tree `{2,nil,{3}}`.
  - `i=3`: left=`gen(2,2)`=[{2,nil,nil}], right=[nil] → tree `{3,{2},nil}`.
  - returns `[{2,nil,{3}}, {3,{2},nil}]`.
  - Creates: `{1,nil,{2,nil,{3}}}` and `{1,nil,{3,{2},nil}}`.
- `i=2`: left=`gen(1,1)`=[{1}], right=`gen(3,3)`=[{3}] → `{2,{1},{3}}`.
- `i=3`: left=`gen(1,2)`=[{1,nil,{2}}, {2,{1},nil}], right=[nil]:
  - Creates `{3,{1,nil,{2}},nil}` and `{3,{2,{1},nil},nil}`.

Total: 5 trees ✓

**Catalan(3) = 5**: for n=1:1, n=2:2, n=3:5, n=4:14, n=5:42, n=6:132.

---

## Approach 2 — Memoized Recursion

### Intuition
The `generate(start, end)` call recomputes the same ranges multiple times for larger `n`. Cache results by `(start, end)` key.

**Note:** the trees for `generate(1,3)` and `generate(2,4)` have different node values even though they have the same structure, so caching by `(start, end)` is correct — each unique `(start, end)` pair has unique node values.

### Complexity
- Same asymptotic as approach 1 but avoids recomputing overlapping subproblems.

### Code
```go
func generateTreesMemo(n int) []*TreeNode {
    memo := make(map[[2]int][]*TreeNode)
    var generate func(start, end int) []*TreeNode
    generate = func(start, end int) []*TreeNode {
        if start > end { return []*TreeNode{nil} }
        key := [2]int{start, end}
        if trees, ok := memo[key]; ok { return trees }
        var allTrees []*TreeNode
        for i := start; i <= end; i++ {
            for _, left := range generate(start, i-1) {
                for _, right := range generate(i+1, end) {
                    allTrees = append(allTrees, &TreeNode{Val: i, Left: left, Right: right})
                }
            }
        }
        memo[key] = allTrees
        return allTrees
    }
    return generate(1, n)
}
```

### Dry Run (n=3)

Same divide-and-conquer as Approach 1, but each `(start,end)` result is stored in `memo` and reused on a repeat call instead of being rebuilt.

| call `generate(start,end)` | in memo? | action | result stored | count |
|----------------------------|----------|--------|---------------|-------|
| `(1,3)`                     | no       | compute; try roots 1,2,3 | `[5 trees]` | 5 |
| ↳ `(1,0)` (left of root 1)  | no       | start>end base case | `[nil]` | 1 |
| ↳ `(2,3)` (right of root 1) | no       | compute; try roots 2,3 | `[2 trees]` | 2 |
| ↳↳ `(3,3)`                  | no       | compute → `{3}` | `[{3}]` | 1 |
| ↳↳ `(2,2)`                  | no       | compute → `{2}` | `[{2}]` | 1 |
| ↳ `(1,1)` (left of root 2)  | no       | compute → `{1}` | `[{1}]` | 1 |
| ↳ `(3,3)` (right of root 2) | **yes**  | reuse cached `[{3}]` | — | 1 |
| ↳ `(1,2)` (left of root 3)  | no       | compute; try roots 1,2 | `[2 trees]` | 2 |
| ↳↳ `(2,2)`                  | **yes**  | reuse cached `[{2}]` | — | 1 |
| ↳↳ `(1,1)`                  | **yes**  | reuse cached `[{1}]` | — | 1 |

`generate(1,3)` returns 5 trees ✓ — identical output to Approach 1, but `(3,3)`, `(2,2)`, `(1,1)` are each computed once and served from `memo` thereafter.

---

## Key Takeaways
- Return `[nil]` (not `[]`) as the base case for empty ranges — this enables the cross-product logic to produce trees even when one subtree is empty.
- The number of unique BSTs with n nodes = Catalan(n): 1, 1, 2, 5, 14, 42, 132, 429...
- Catalan number appears in: BST counting, balanced parentheses, polygon triangulation, mountain ranges.
- Different from #96 (Unique BSTs) which only counts — here we enumerate all trees.

---

## Related Problems
- LeetCode #96 — Unique Binary Search Trees (count only, not enumerate)
- LeetCode #241 — Different Ways to Add Parentheses (same divide-and-conquer on intervals)
- LeetCode #894 — All Possible Full Binary Trees (enumerate trees with different structure)
