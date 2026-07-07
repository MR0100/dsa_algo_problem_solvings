# Divide and Conquer

> **Category:** Algorithm Design Paradigm
> **Typical complexity:** governed by the recurrence `T(n) = a·T(n/b) + f(n)` — most often O(n log n) time, O(log n) space (recursion stack)

---

## What it is

Divide and conquer (D&C) solves a problem by:

1. **Divide** — split the input into two or more **independent** subproblems,
   usually of roughly equal size.
2. **Conquer** — solve each subproblem recursively; a subproblem small enough
   (the *base case*) is solved directly.
3. **Combine** — merge the sub-solutions into the solution for the whole input.

The paradigm pays off when the divide + combine work is cheap relative to the
brute-force cost of the whole problem. Halving `n` repeatedly gives a recursion
tree of depth ⌈log₂ n⌉; if each level does O(n) combine work the total is
O(n log n) — this is exactly why merge sort beats the O(n²) quadratic sorts.

Classic members of the family: **merge sort**, **quick sort / quickselect**,
**binary search** (a degenerate D&C where one half is discarded, so there is no
combine step), **fast exponentiation**, **Karatsuba multiplication**,
**closest pair of points**, and most **recursive tree constructions** ("build
left subtree, build right subtree, attach to root").

### D&C vs. dynamic programming — the litmus test

Both split a problem into subproblems. The difference:

- **D&C:** subproblems are **disjoint / non-overlapping** — each element of the
  input belongs to exactly one subproblem. Plain recursion suffices.
- **DP:** subproblems **overlap** — the same subproblem is reached via many
  paths, so you memoise or tabulate.

If your recursion tree recomputes the same state, you have DP in disguise;
if every node of the tree works on a fresh slice of the input, it is D&C.

---

## How to recognise it — signals in the problem statement

| Signal | Example phrasing / problem |
|--------|----------------------------|
| Required complexity is **O(n log n)** or **O(log(m+n))** | "overall run time complexity should be O(log (m+n))" — #0004 |
| Input can be **split at a midpoint** and the halves solved independently | maximum subarray — #0053 (answer is in left half, right half, or crosses the middle) |
| **"Merge k …"** — pairwise merging halves the count each round | merge k sorted lists — #0023 |
| A **root/pivot naturally partitions** the input into left and right parts | build BST from sorted array/list — #0108, #0109; build tree from traversals — #0105, #0106 |
| **Sorted input → balanced structure** | "height-balanced BST" from sorted data — #0108, #0109 |
| **Generate all structures** over a range that a chosen root splits | unique BSTs II — #0095 |
| A constraint like **"at most two transactions"** invites a split point | best time to buy/sell stock III — #0123 (best split of the timeline) |
| Counting **cross-boundary pairs** (inversions, reverse pairs, range sums) | LC #493, #315, #327 — merge-sort variants where the combine step counts |
| **k-th / median / top-k** selection without full sorting | quickselect — LC #215; median of two sorted arrays — #0004 |

Rule of thumb: ask *"if I magically had the answers for the left half and the
right half, could I stitch the full answer together quickly?"* If yes — and the
halves don't share state — divide and conquer applies.

---

## General templates (Go)

### Template 1 — canonical array D&C (merge-sort shape)

```go
// solve returns the answer for nums[lo..hi] (inclusive).
func solve(nums []int, lo, hi int) Result {
    // ── Base case: a slice of size 0 or 1 is trivially solved. ──
    if lo >= hi {
        return baseResult(nums, lo)
    }

    // ── Divide: split at the midpoint. ──
    mid := lo + (hi-lo)/2 // overflow-safe midpoint

    // ── Conquer: recurse on the two independent halves. ──
    left := solve(nums, lo, mid)    // answer for nums[lo..mid]
    right := solve(nums, mid+1, hi) // answer for nums[mid+1..hi]

    // ── Combine: merge sub-answers, handling anything that CROSSES mid. ──
    // This is where each problem differs: merge sorted runs, take the max
    // of {left, right, crossing}, count cross-boundary pairs, etc.
    return combine(left, right, nums, lo, mid, hi)
}
```

Commented pseudocode of the shape:

```text
solve(range):
    if range is trivially small: return direct answer   // base case
    mid  = middle of range                               // divide
    L    = solve(left  half)                             // conquer
    R    = solve(right half)                             // conquer
    return combine(L, R, cross-boundary work)            // combine
```

### Template 2 — build-a-tree D&C (root splits the range)

```go
// build constructs the (sub)tree for the ordered range [lo, hi].
func build(vals []int, lo, hi int) *TreeNode {
    // ── Base case: empty range → nil subtree. ──
    if lo > hi {
        return nil
    }

    // ── Divide: pick the element that becomes the root.
    //    Sorted array → midpoint (balance); traversal problems → the
    //    preorder-first / postorder-last element locates the root. ──
    mid := lo + (hi-lo)/2
    root := &TreeNode{Val: vals[mid]}

    // ── Conquer: everything left of the root forms the left subtree,
    //    everything right forms the right subtree — fully independent. ──
    root.Left = build(vals, lo, mid-1)
    root.Right = build(vals, mid+1, hi)

    // ── Combine: attaching the children IS the combine step. ──
    return root
}
```

### Template 3 — pairwise reduction ("merge k" shape)

```go
// mergeKLists repeatedly merges lists in pairs: k → k/2 → k/4 → … → 1.
// Each element is touched once per round and there are ⌈log k⌉ rounds,
// so total time is O(N log k) instead of O(N·k) for one-by-one merging.
func mergeKLists(lists []*ListNode) *ListNode {
    if len(lists) == 0 {
        return nil
    }
    for len(lists) > 1 {
        var merged []*ListNode
        for i := 0; i < len(lists); i += 2 {
            if i+1 < len(lists) {
                merged = append(merged, mergeTwo(lists[i], lists[i+1]))
            } else {
                merged = append(merged, lists[i]) // odd one out survives the round
            }
        }
        lists = merged // half as many lists remain
    }
    return lists[0]
}
```

### Analysing the runtime — Master Theorem cheat sheet

For `T(n) = a·T(n/b) + O(n^d)` (a = subproblems, n/b = subproblem size,
n^d = divide+combine cost):

| Case | Condition | Result | Example |
|------|-----------|--------|---------|
| 1 | `d > log_b a` | O(n^d) | combine dominates |
| 2 | `d = log_b a` | O(n^d · log n) | merge sort: a=2, b=2, d=1 → O(n log n) |
| 3 | `d < log_b a` | O(n^(log_b a)) | recursion dominates (e.g. Karatsuba) |

Binary search: a=1, b=2, d=0 → case 2 → O(log n).

---

## Worked example — Maximum Subarray (#0053) traced step by step

Problem: find the contiguous subarray with the largest sum.
Key D&C insight: for any midpoint, the best subarray either lies **entirely in
the left half**, **entirely in the right half**, or **crosses the middle** —
and the crossing case is computable in O(n) (best suffix of the left half +
best prefix of the right half).

```go
func maxSubArray(nums []int, lo, hi int) int {
    if lo == hi {
        return nums[lo] // single element
    }
    mid := lo + (hi-lo)/2
    left := maxSubArray(nums, lo, mid)
    right := maxSubArray(nums, mid+1, hi)

    // best sum of a suffix ending at mid
    bestLeft, sum := math.MinInt, 0
    for i := mid; i >= lo; i-- {
        sum += nums[i]
        bestLeft = max(bestLeft, sum)
    }
    // best sum of a prefix starting at mid+1
    bestRight, sum2 := math.MinInt, 0
    for i := mid + 1; i <= hi; i++ {
        sum2 += nums[i]
        bestRight = max(bestRight, sum2)
    }
    cross := bestLeft + bestRight

    return max(left, max(right, cross))
}
```

Trace on `nums = [-2, 1, -3, 4]` (indices 0..3):

| Step | Call | mid | left result | right result | bestLeft (suffix) | bestRight (prefix) | cross | return |
|------|------|-----|-------------|--------------|-------------------|--------------------|-------|--------|
| 1 | `solve(0,3)` | 1 | → step 2 | → step 5 | — | — | — | — |
| 2 | `solve(0,1)` | 0 | → step 3 | → step 4 | — | — | — | — |
| 3 | `solve(0,0)` | — | — | — | — | — | — | **-2** (base) |
| 4 | `solve(1,1)` | — | — | — | — | — | — | **1** (base) |
| 5 | back in `solve(0,1)` | 0 | -2 | 1 | suffix of `[-2]` = -2 | prefix of `[1]` = 1 | -2+1 = **-1** | max(-2, 1, -1) = **1** |
| 6 | `solve(2,3)` | 2 | `solve(2,2)` = **-3** | `solve(3,3)` = **4** | suffix of `[-3]` = -3 | prefix of `[4]` = 4 | -3+4 = **1** | max(-3, 4, 1) = **4** |
| 7 | back in `solve(0,3)` | 1 | 1 (step 5) | 4 (step 6) | best suffix of `[-2,1]` ending at idx 1: `1` (just `[1]`) | best prefix of `[-3,4]`: `-3+4 = 1` | 1+1 = **2** | max(1, 4, 2) = **4** |

Answer: **4** (the subarray `[4]`). Recursion depth log₂ 4 = 2; each level does
O(n) crossing work → O(n log n) total, O(log n) stack space.

---

## Common pitfalls and how to avoid them

1. **Missing or wrong base case → infinite recursion / stack overflow.**
   Always handle the empty range (`lo > hi`) and the single element
   (`lo == hi`) explicitly, and make sure every recursive call strictly
   shrinks the range (`[lo, mid]` + `[mid+1, hi]`, never `[lo, mid]` +
   `[mid, hi]` when `mid` can equal `hi`).

2. **Midpoint overflow.** Use `mid := lo + (hi-lo)/2`, never `(lo+hi)/2`.
   Rarely bites in Go with `int` on 64-bit, but it is the interview-expected
   idiom and matters with 32-bit indices.

3. **Forgetting the crossing case in the combine step.** In problems like
   maximum subarray or counting inversions, the answer that *spans* the
   midpoint is the whole point of the combine step — omitting it silently
   returns wrong answers on inputs like `[1, -1, 2]`.

4. **Expensive combine that kills the complexity.** If combine is O(n²) the
   recursion gives O(n² log n) — worse than brute force. E.g. in #0105/#0106,
   searching for the root's index in the inorder slice linearly at every level
   gives O(n²) worst case; precompute a value→index hash map to make each
   lookup O(1) and the whole build O(n).

5. **Repeated slicing/copying instead of index ranges.** `nums[lo:mid]` in Go
   shares the backing array (cheap), but building *new* slices/lists per call
   (e.g. `append` copies) can add hidden O(n) work and O(n log n) memory.
   Prefer passing `(lo, hi)` indices over materialising sub-slices.

6. **Applying D&C to overlapping subproblems.** If both recursive calls can
   examine the same elements/states, you will recompute exponentially — that
   is a memoisation/DP problem (e.g. naive Fibonacci), not D&C. Check that
   the halves are disjoint before committing to the paradigm.

7. **Unbalanced splits degrade the depth.** Quickselect/quicksort with a bad
   pivot degrades to O(n) depth and O(n²) time; use a random pivot or
   median-of-three. Merge-sort-style fixed midpoints never have this issue.

8. **Recursion-stack space counts.** A "constant extra space" claim is wrong
   for recursive D&C — the stack is O(log n) (balanced) or O(n) (degenerate,
   e.g. building a BST from a skewed input). State it in the complexity
   analysis.

9. **In "merge k" problems, merging one-by-one instead of pairwise.**
   Sequential merging is O(N·k); pairwise (tournament) merging is O(N log k).
   The halving structure is what makes it divide and conquer — see #0023.

---

## Problems in this repo

| # | Problem | How D&C is used |
|---|---------|-----------------|
| 0004 | [Median of Two Sorted Arrays](../0004_median_of_two_sorted_arrays/README.md) | Binary-search partition of the smaller array — discard half the search space each step for O(log(min(m,n))) |
| 0023 | [Merge k Sorted Lists](../0023_merge_k_sorted_lists/README.md) | Pairwise merging halves the list count each round: O(N log k) |
| 0053 | [Maximum Subarray](../0053_maximum_subarray/README.md) | Best subarray is left, right, or crossing the midpoint (worked example above) |
| 0095 | [Unique Binary Search Trees II](../0095_unique_binary_search_trees_ii/README.md) | Each candidate root splits [1..n]; combine = cartesian product of left/right subtree sets |
| 0105 | [Construct Binary Tree from Preorder and Inorder Traversal](../0105_construct_binary_tree_from_preorder_and_inorder_traversal/README.md) | Preorder head is the root; its inorder index splits both traversals into left/right subtree ranges |
| 0106 | [Construct Binary Tree from Inorder and Postorder Traversal](../0106_construct_binary_tree_from_inorder_and_postorder_traversal/README.md) | Postorder tail is the root; same inorder split, building right-then-left |
| 0108 | [Convert Sorted Array to Binary Search Tree](../0108_convert_sorted_array_to_binary_search_tree/README.md) | Midpoint becomes the root → guaranteed height balance (Template 2 verbatim) |
| 0109 | [Convert Sorted List to Binary Search Tree](../0109_convert_sorted_list_to_binary_search_tree/README.md) | Same midpoint idea on a linked list — via slow/fast pointer or inorder simulation |
| 0123 | [Best Time to Buy and Sell Stock III](../0123_best_time_to_buy_and_sell_stock_iii/README.md) | Split the timeline: best single transaction in prefix + best in suffix, maximised over the split point |

> Problems #0131+ are being added concurrently; a later pass will extend this table.

## Related concepts in this library

- [`binary_search.md`](binary_search.md) — D&C with one half discarded and no combine step
- [`sorting.md`](sorting.md) — merge sort and quicksort are the archetypal D&C algorithms
- [`tree_traversal.md`](tree_traversal.md) — most recursive tree algorithms are structurally D&C
- [`heap_priority_queue.md`](heap_priority_queue.md) — the alternative to D&C for "merge k" problems
