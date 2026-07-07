# Binary Search

> **Category:** Divide and Conquer / Search
> **Typical complexity:** O(log n) time, O(1) space (iterative)

---

## What it is

Binary search finds a target — or a *boundary* — in a **monotonic** search space
by repeatedly halving the range. At every step you inspect the middle element,
decide which half can be discarded, and recurse/iterate on the remaining half.
Because the search space shrinks by half each step, only ⌈log₂ n⌉ + 1 probes are
ever needed.

The key insight that generalises far beyond "find x in a sorted array":

> Binary search works on **any predicate `P(i)` that is monotonic** over the
> search space — i.e. once it flips from `false` to `true` it never flips back
> (`F F F F T T T`). Binary search finds the flip point in O(log n) checks.

"Sorted array contains target" is just the special case where
`P(i) = (nums[i] >= target)`.

---

## How to recognise it — signals in the problem statement

| Signal | Example phrasing |
|--------|------------------|
| Input is **sorted** (or "rotated sorted", "sorted rows") | "given a sorted array…" — #0033, #0034, #0035, #0074 |
| Required complexity is **O(log n)** | "must run in O(log n) time" — #0004, #0033, #0034 |
| **Find first / last / count** of something in ordered data | "first and last position of target" — #0034 |
| **Insert position** / floor / ceiling queries | "index where it would be inserted" — #0035 |
| **"Minimise the maximum" / "maximise the minimum"** | binary search on the answer (e.g. LC #410, #875, #1011) |
| Answer lies in a **numeric range and feasibility is monotonic** | "largest k such that k² ≤ x" — #0069 |
| A **monotonic yes/no property over a length or value** | "if a prefix of length L is common, every shorter one is too" — #0014 |
| Need to avoid multiplication/division but can **halve/double** | #0029 uses the same halving idea via doubling |

Rule of thumb: if you can phrase the problem as *"find the smallest/largest X
such that condition(X) holds"* and `condition` is monotonic, it is binary
search — even if no array is sorted.

---

## The three flavours

1. **Exact-match search** — return index of `target`, or -1. (#0033, #0074)
2. **Boundary search (lower/upper bound)** — find first index where a
   monotonic predicate is true. Subsumes flavour 1 and is the least
   bug-prone template. (#0034, #0035)
3. **Binary search on the answer space** — the "array" is an implicit range
   of candidate answers `[lo, hi]`; a feasibility check plays the role of
   comparison. (#0069, and partition search in #0004)

---

## Templates in Go

### Template 1 — classic exact match

```go
// binarySearch returns the index of target in sorted nums, or -1.
//
// Invariant: if target exists, it is always inside [lo, hi].
//
// Time:  O(log n)
// Space: O(1)
func binarySearch(nums []int, target int) int {
    lo, hi := 0, len(nums)-1        // inclusive search range [lo, hi]
    for lo <= hi {                  // range still non-empty
        mid := lo + (hi-lo)/2       // overflow-safe midpoint
        switch {
        case nums[mid] == target:
            return mid              // found it
        case nums[mid] < target:
            lo = mid + 1            // target is strictly right of mid
        default:
            hi = mid - 1            // target is strictly left of mid
        }
    }
    return -1                       // range empty → not present
}
```

### Template 2 — lower bound (first index where predicate is true) — **learn this one**

```go
// lowerBound returns the smallest index i in [0, n] such that P(i) is true,
// where P is monotonic (false...false, true...true). Returns n if P is never
// true. With P(i) = nums[i] >= target this is exactly sort.SearchInts.
//
// Invariant: P is false for every index < lo, true for every index >= hi.
// When lo == hi, that common value is the boundary.
//
// Time:  O(log n)
// Space: O(1)
func lowerBound(n int, P func(i int) bool) int {
    lo, hi := 0, n                  // half-open answer range [lo, hi]
    for lo < hi {                   // stops when lo == hi — no infinite loop
        mid := lo + (hi-lo)/2       // mid < hi always, so hi = mid shrinks
        if P(mid) {
            hi = mid                // mid might BE the answer → keep it
        } else {
            lo = mid + 1            // mid is definitely not → discard it
        }
    }
    return lo                       // == hi == first true index
}
```

Everything reduces to `lowerBound`:

```go
// first occurrence of target (or -1)
i := lowerBound(len(nums), func(i int) bool { return nums[i] >= target })
if i == len(nums) || nums[i] != target { i = -1 }

// last occurrence of target = (first index with nums[i] > target) - 1
j := lowerBound(len(nums), func(i int) bool { return nums[i] > target }) - 1

// insert position (LC #0035) is lowerBound directly
// count of target = j - i + 1  (LC #0034)
```

### Template 3 — binary search on the answer space

```go
// maxFeasible finds the LARGEST k in [lo, hi] with feasible(k) true,
// where feasible is monotonically decreasing (true...true, false...false).
// Example (LC #0069 Sqrt): feasible(k) = k*k <= x.
//
// Time:  O(log(hi-lo) * cost(feasible))
// Space: O(1)
func maxFeasible(lo, hi int, feasible func(k int) bool) int {
    ans := lo - 1                   // sentinel: "nothing feasible yet"
    for lo <= hi {
        mid := lo + (hi-lo)/2
        if feasible(mid) {
            ans = mid               // record candidate...
            lo = mid + 1            // ...and try to push higher
        } else {
            hi = mid - 1            // too big → shrink from the right
        }
    }
    return ans
}
```

### Go standard library

```go
sort.SearchInts(nums, target)          // lower bound over []int
sort.Search(n, func(i int) bool {...}) // generic lower bound (Template 2)
```

---

## Worked example — lower bound, traced step by step

Find the **first occurrence of `target = 8`** in
`nums = [5, 7, 7, 8, 8, 10]` (LC #0034, Example 1). Expected answer: index 3.

Predicate: `P(i) = nums[i] >= 8` → over indices 0..5: `F F F T T T`.
Run `lowerBound(6, P)`:

| Step | lo | hi | mid | nums[mid] | P(mid) = nums[mid] ≥ 8 | Action |
|------|----|----|-----|-----------|------------------------|--------|
| 1 | 0 | 6 | 3 | 8 | true | mid could be the answer → `hi = 3` |
| 2 | 0 | 3 | 1 | 7 | false | mid ruled out → `lo = 2` |
| 3 | 2 | 3 | 2 | 7 | false | mid ruled out → `lo = 3` |
| 4 | 3 | 3 | — | — | — | `lo == hi` → loop exits |

Return `lo = 3`. ✓ Three probes for six elements — and `nums[3] == 8`, so the
target exists. For the *last* occurrence, run again with
`P(i) = nums[i] > 8` (→ returns 5), subtract 1 → index 4. Result `[3, 4]`.

---

## Common pitfalls and how to avoid them

1. **Infinite loops from a bad mid/update pairing.**
   With `lo < hi` and `hi = mid`, you MUST use the left-leaning mid
   `mid = lo + (hi-lo)/2` (never `+1`), otherwise `mid` can equal `hi` and the
   range never shrinks. Symmetrically, if you write `lo = mid` (searching for
   the *last* true), you need the right-leaning mid `mid = lo + (hi-lo+1)/2`.
   Rule: *the half that keeps `mid` must be the half that `mid` rounds away
   from.*

2. **Integer overflow in `(lo + hi) / 2`.**
   In fixed-width languages `lo + hi` can overflow. Always write
   `mid := lo + (hi-lo)/2`. (Go ints are 64-bit on modern platforms, but the
   habit matters — and it *does* overflow with `int32` or in Java/C++.)

3. **Off-by-one from mixing inclusive and half-open conventions.**
   Pick one and be consistent:
   - Inclusive `[lo, hi]` → loop `lo <= hi`, updates `mid±1` (Template 1/3).
   - Half-open `[lo, hi)` → loop `lo < hi`, updates `hi = mid` / `lo = mid+1`
     (Template 2).
   Mixing them (`lo <= hi` with `hi = mid`) is the classic infinite loop.

4. **Forgetting the "not found" check after a boundary search.**
   `lowerBound` returns an *insert position*, not proof of existence. Always
   verify `i < len(nums) && nums[i] == target` before claiming a hit.

5. **Returning the wrong side at termination.**
   When the loop ends with `lo == hi` (half-open), `lo` is the first true
   index. When it ends with `lo > hi` (inclusive), `lo` is the insert
   position and `hi` is the last false index. Know which one your invariant
   promises — derive it from the invariant, don't guess.

6. **Assuming the predicate is monotonic when it isn't.**
   Binary search silently returns garbage on non-monotonic data. Rotated
   arrays (#0033/#0081) are *not* globally sorted — you must first decide
   which half is sorted at each step. Duplicates (#0081) can make the
   sorted-half test ambiguous (`nums[lo] == nums[mid] == nums[hi]`) and
   degrade the worst case to O(n) — shrink both ends by one in that case.

7. **Searching indices when you should search values (or vice versa).**
   "Minimise the max page load", "smallest divisor", "sqrt" — the array of
   *candidate answers* is implicit. Define `feasible(k)`, prove it is
   monotonic in `k`, set correct `lo`/`hi` bounds, then apply Template 3.

8. **Bad initial bounds on answer-space searches.**
   `hi` must provably contain the answer (e.g. for sqrt, `hi = x/2 + 1` or
   simply `x`), and `lo` must be a valid extreme. Off-by-one at the bounds is
   the #1 bug in "binary search on answer" solutions.

---

## Problems in this repo

| # | Problem | How binary search is used |
|---|---------|---------------------------|
| 0004 | [Median of Two Sorted Arrays](../0004_median_of_two_sorted_arrays/README.md) | Binary search the **partition index** of the shorter array — O(log(min(m,n))) |
| 0014 | [Longest Common Prefix](../0014_longest_common_prefix/README.md) | Binary search on **prefix length** — monotonic "is length-L prefix common?" predicate |
| 0033 | [Search in Rotated Sorted Array](../0033_search_in_rotated_sorted_array/README.md) | Modified binary search — decide which half is sorted at each step |
| 0034 | [Find First and Last Position of Element in Sorted Array](../0034_find_first_and_last_position_of_element_in_sorted_array/README.md) | Two boundary searches — **lower bound** and **upper bound** (Template 2) |
| 0035 | [Search Insert Position](../0035_search_insert_position/README.md) | The canonical **lower bound** / insert-position problem |
| 0069 | [Sqrt(x)](../0069_sqrt_x/README.md) | Binary search on the **answer space** — largest k with k² ≤ x (Template 3) |
| 0074 | [Search a 2D Matrix](../0074_search_a_2d_matrix/README.md) | Treat the matrix as a flat sorted array via index mapping `mid → (mid/n, mid%n)` |
| 0081 | [Search in Rotated Sorted Array II](../0081_search_in_rotated_sorted_array_ii/README.md) | Rotated search **with duplicates** — ambiguous halves degrade worst case to O(n) |

*(Problems 0131+ are being added; a later pass will extend this table.)*
