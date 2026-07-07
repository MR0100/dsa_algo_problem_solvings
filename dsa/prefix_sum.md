# Prefix Sum

## What it is

A **prefix sum** (also called a cumulative sum or running total) is an
auxiliary array `prefix` where each element stores the sum of all elements of
the original array up to that position:

```
prefix[i] = nums[0] + nums[1] + ... + nums[i-1]        (0-based, "offset" style)
```

Built once in O(n), it answers **any range-sum query in O(1)**:

```
sum(nums[l..r]) = prefix[r+1] - prefix[l]
```

That single identity is the whole trick. Anything that can be phrased as
"the sum (or count / XOR / product) of a contiguous range" stops being a loop
and becomes a subtraction.

Prefix sums generalise beyond `+`:

| Variant | prefix definition | range query |
|---------|-------------------|-------------|
| Sum | `p[i+1] = p[i] + a[i]` | `p[r+1] - p[l]` |
| Count (of a property) | `p[i+1] = p[i] + boolToInt(a[i] has property)` | `p[r+1] - p[l]` |
| XOR | `p[i+1] = p[i] ^ a[i]` | `p[r+1] ^ p[l]` |
| Product (no zeros) | `p[i+1] = p[i] * a[i]` | `p[r+1] / p[l]` |
| 2-D sum | `P[i+1][j+1] = A[i][j] + P[i][j+1] + P[i+1][j] - P[i][j]` | inclusion–exclusion (below) |

The inverse idea is the **difference array**: apply many range *updates* in
O(1) each, then one prefix-sum pass materialises the final values.

---

## How to recognise it — signals in the problem statement

Reach for prefix sums when you see:

1. **"Sum of a subarray / range" asked repeatedly.**
   "Given many queries `(l, r)`, return the sum of elements between `l` and
   `r`" — the textbook case (LeetCode #303, #304).
2. **"Count / find subarrays whose sum equals K"** (or is divisible by K, or
   has some target property). Every subarray sum is a *difference of two
   prefix values*, so "subarray with sum K" ⇔ "two prefix values that differ
   by K" — pair prefix sums with a **hash map** (LeetCode #560, #523, #525,
   #974). This is exactly the Two Sum complement trick applied to prefixes.
3. **"Immutable array" + "many queries".** The word *immutable* is a giant
   arrow pointing at precomputation.
4. **Left side vs right side.** "Pivot index where left sum equals right
   sum", "product of array except self" — anything comparing an
   aggregate *before* index `i` with an aggregate *after* it (LeetCode #724,
   #238 which uses prefix *products*).
5. **Binary array reframing.** "Longest subarray with equal 0s and 1s" —
   map 0 → −1 and the problem becomes "longest subarray with sum 0"
   (LeetCode #525).
6. **Many range *updates*, values read once at the end.** "Add `v` to every
   element in `[l, r]` for each of Q bookings/flights" — difference array
   (LeetCode #1109, #370).
7. **2-D region sums.** "Sum of elements inside a rectangle" (LeetCode #304).
8. **A brute force that recomputes overlapping sums.** If your O(n²) inner
   loop is `for j := i; ...; sum += nums[j]`, a prefix sum usually deletes
   one of the loops — or, combined with a hash map, both.

**When it does NOT apply:** if the array is frequently *mutated* between
queries, a static prefix array goes stale after every write — that is Fenwick
tree (BIT) / segment tree territory. Also, "max/min over a range" does not
telescope like sums do; use sparse tables or segment trees for those.

---

## General templates (Go)

### 1. Basic prefix array + O(1) range query

The `n+1`-length "offset" layout (`prefix[0] = 0`) avoids every off-by-one
and special case for `l = 0`:

```go
// buildPrefix returns p of length n+1 where p[i] = sum of nums[0..i-1].
//
// Pseudocode:
//   p[0] = 0                      // empty prefix sums to zero
//   for i in 0..n-1:
//       p[i+1] = p[i] + nums[i]   // extend the running total by one element
func buildPrefix(nums []int) []int {
    p := make([]int, len(nums)+1) // one extra slot for the empty prefix
    for i, v := range nums {
        p[i+1] = p[i] + v // each entry = previous entry + current element
    }
    return p
}

// rangeSum returns sum(nums[l..r]) inclusive, in O(1).
//
//   sum(l..r) = p[r+1] - p[l]    // total up to r, minus total before l
func rangeSum(p []int, l, r int) int {
    return p[r+1] - p[l]
}
```

### 2. Prefix sum + hash map — count subarrays with sum K

The single most interview-frequent form (LeetCode #560 pattern):

```go
// subarraySumK counts contiguous subarrays summing exactly to k.
//
// Key identity: sum(i..j) == k  ⇔  prefix[j+1] - prefix[i] == k
//                              ⇔  prefix[i] == prefix[j+1] - k
// So while scanning, ask: "how many earlier prefixes equal current - k?"
//
// Time:  O(n)   Space: O(n)
func subarraySumK(nums []int, k int) int {
    count := 0
    sum := 0                    // running prefix sum up to current index
    seen := map[int]int{0: 1}   // prefix value -> how many times seen;
                                // {0:1} = the empty prefix, so subarrays
                                // starting at index 0 are counted too
    for _, v := range nums {
        sum += v                // extend prefix to include v
        count += seen[sum-k]    // every earlier prefix equal to sum-k
                                // closes one subarray with sum k here
        seen[sum]++             // register current prefix for later indices
    }
    return count
}
```

Swap the map's meaning to solve the siblings:
- **longest** subarray with sum k → store *earliest index* per prefix value.
- sum **divisible by k** → key on `((sum % k) + k) % k` instead of `sum`.
- equal 0s/1s → map 0 → −1, then k = 0 with earliest-index map.

### 3. Difference array — batch range updates

```go
// applyRangeUpdates: for each update (l, r, v), add v to a[l..r]; then
// materialise. Q updates + one pass = O(n + Q) instead of O(n*Q).
//
// Pseudocode:
//   diff[l]   += v    // start adding v from index l onward
//   diff[r+1] -= v    // stop adding v after index r
//   answer = prefix sum of diff
func applyRangeUpdates(n int, updates [][3]int) []int {
    diff := make([]int, n+1) // extra slot so r+1 never overflows
    for _, u := range updates {
        l, r, v := u[0], u[1], u[2]
        diff[l] += v   // v takes effect at l
        diff[r+1] -= v // v is cancelled just past r
    }
    a := make([]int, n)
    running := 0
    for i := 0; i < n; i++ {
        running += diff[i] // prefix sum of deltas = final value at i
        a[i] = running
    }
    return a
}
```

### 4. 2-D prefix sum — O(1) rectangle sums

```go
// build2D: P[i][j] = sum of the rectangle A[0..i-1][0..j-1].
// Add the cell, the sum above, the sum to the left; the top-left block was
// added twice, subtract it once (inclusion–exclusion).
func build2D(A [][]int) [][]int {
    m, n := len(A), len(A[0])
    P := make([][]int, m+1)
    for i := range P {
        P[i] = make([]int, n+1) // row 0 / col 0 stay zero (empty prefixes)
    }
    for i := 0; i < m; i++ {
        for j := 0; j < n; j++ {
            P[i+1][j+1] = A[i][j] + P[i][j+1] + P[i+1][j] - P[i][j]
        }
    }
    return P
}

// sumRegion returns the sum of A[r1..r2][c1..c2] inclusive.
func sumRegion(P [][]int, r1, c1, r2, c2 int) int {
    return P[r2+1][c2+1] - P[r1][c2+1] - P[r2+1][c1] + P[r1][c1]
}
```

---

## Worked example — Subarray Sum Equals K

`nums = [1, 2, 3, -2, 2]`, `k = 3`.

First, the prefix array (offset layout): `p = [0, 1, 3, 6, 4, 6]`.
Every subarray with sum 3 is a pair `i < j` with `p[j] − p[i] = 3`:

- `p[2] − p[0] = 3` → `nums[0..1] = [1, 2]`
- `p[3] − p[2] = 3` → `nums[2..2] = [3]`
- `p[4] − p[1] = 3` → `nums[1..3] = [2, 3, −2]`
- `p[5] − p[2] = 3` → `nums[2..4] = [3, −2, 2]`

Expected answer: **4**.

Trace of `subarraySumK` (template 2). `seen` starts as `{0: 1}`.

| step | v | sum (prefix) | need = sum−k | seen[need] | count | subarray closed here | seen after |
|------|----|----|----|---|---|----------------------|---------------------------|
| 1 | 1  | 1 | −2 | 0 | 0 | — | {0:1, 1:1} |
| 2 | 2  | 3 | 0  | 1 | 1 | `[1,2]` (p2−p0) | {0:1, 1:1, 3:1} |
| 3 | 3  | 6 | 3  | 1 | 2 | `[3]` (p3−p2) | {0:1, 1:1, 3:1, 6:1} |
| 4 | −2 | 4 | 1  | 1 | 3 | `[2,3,−2]` (p4−p1) | {0:1, 1:1, 3:1, 6:1, 4:1} |
| 5 | 2  | 6 | 3  | 1 | 4 | `[3,−2,2]` (p5−p2) | {0:1, 1:1, 3:1, 6:2, 4:1} |

Final `count = 4`, matching the enumeration above.

**Lesson embedded in the example:** with negative numbers, subarrays hiding
"cancel-out" segments are easy to miss by eye — the prefix identity finds all
of them mechanically. (This is also why sliding window *cannot* replace
prefix+hashmap here: negatives break the window's monotonicity.)

---

## Common pitfalls (and fixes)

1. **Off-by-one in the range formula.** With the `n+1` offset layout the
   query is `p[r+1] − p[l]` — not `p[r] − p[l]` or `p[r] − p[l−1]`. Always
   use the offset layout; it also makes `l = 0` need no special case.
2. **Forgetting to seed the map with the empty prefix** (`seen[0] = 1` /
   `firstIdx[0] = -1`). Without it, subarrays that start at index 0 are
   silently dropped. This is the #1 bug in the sum-equals-K family.
3. **Update order inside the loop:** query the map (`count += seen[sum-k]`)
   *before* inserting the current prefix. Inserting first lets a subarray of
   length 0 match itself when `k == 0`.
4. **Reaching for sliding window when the array has negatives.** Sliding
   window needs "adding an element never shrinks the sum" (all non-negative).
   With negatives, prefix sum + hash map is the correct tool.
5. **Negative numbers and `%` in Go.** Go's `%` can return negatives
   (`-4 % 3 == -1`). For divisible-by-K problems normalise:
   `r := ((sum % k) + k) % k`.
6. **Overflow.** `n` up to 10⁵ with values up to 10⁹ overflows `int32`. Go's
   `int` is 64-bit on mainstream platforms, but be explicit (`int64`) when
   porting or when products are involved.
7. **Using a static prefix array on a mutable array.** One write invalidates
   O(n) prefix entries. If the problem interleaves updates and queries, use a
   Fenwick tree or segment tree instead.
8. **2-D inclusion–exclusion sign errors.** Rectangle sum is
   `P[r2+1][c2+1] − P[r1][c2+1] − P[r2+1][c1] + P[r1][c1]` — two minuses, one
   plus-back. Draw the four rectangles once; memorise the shape, not the
   indices.
9. **Difference array: forgetting `diff[r+1] -= v`** (or sizing `diff` as `n`
   so `r+1` panics at the last index). Allocate `n+1`.
10. **Recomputing the prefix inside a query loop.** The whole point is
    build-once O(n), query-many O(1). If your prefix build is inside the
    per-query loop you are back to brute force with extra steps.

---

## Problems in this repo

Direct prefix-sum problems from the 0131+ range are being added; links below
are what exists today.

- [0053 — Maximum Subarray](../0053_maximum_subarray/README.md) — Kadane's
  algorithm is prefix sums in disguise: `maxSubarray = max over j of
  (prefix[j] − min prefix before j)`; resetting `curSum` when it goes
  negative is exactly "discard the smaller prefix".
- [0001 — Two Sum](../0001_two_sum/README.md) — supplies the complement +
  hash-map trick that prefix sums reuse verbatim in the sum-equals-K family
  (see template 2 and LeetCode #560).

Classic LeetCode problems for this concept (link here as they land in the
repo): #303 Range Sum Query – Immutable, #304 Range Sum Query 2D, #238
Product of Array Except Self (prefix products), #525 Contiguous Array, #523
Continuous Subarray Sum, #560 Subarray Sum Equals K, #724 Find Pivot Index,
#974 Subarray Sums Divisible by K, #1109 Corporate Flight Bookings
(difference array).
