# Longest Increasing Subsequence (LIS)

> **Problem shape:** given a sequence, find the length (and optionally the elements) of the longest **subsequence** (order preserved, elements need not be contiguous) that is strictly increasing.
> **Two workhorses:** O(n²) DP (`dp[i]` = LIS ending at `i`) and O(n log n) **patience sorting** (a `tails` array + binary search).

---

## What it is

A **subsequence** keeps the original left-to-right order but may skip elements (unlike a *substring/subarray*, which must be contiguous). The LIS of `[10, 9, 2, 5, 3, 7, 101, 18]` is `[2, 3, 7, 18]` (or `[2, 3, 7, 101]`) — length **4**.

LIS is the canonical example of a problem with two classic solutions at different complexities, and a beautiful bridge (**patience sorting**) between a card game and binary search. It is also a *pattern generator*: many problems reduce to "sort/transform, then run LIS" — Russian Doll Envelopes, Largest Divisible Subset, Longest Chain of pairs, box stacking, and more.

Key facts to internalise:

- **Strict vs non-strict** matters. "Strictly increasing" uses `<`; "non-decreasing" uses `≤`. In the O(n log n) method this is the difference between `LowerBound` (strict) and `UpperBound` (non-decreasing) — a one-line change that is a frequent source of off-by-one bugs.
- The `tails` array in the fast method is **not** an actual subsequence — its *length* is the answer, but its contents are only a set of "smallest possible tail for each length".
- LIS ↔ **Longest Common Subsequence**: the LIS of `a` equals the LCS of `a` and `sorted(unique(a))`. Rarely the fastest route, but a useful mental link.

---

## When to recognise it

| Signal in the problem | Why it's LIS |
|-----------------------|--------------|
| "longest **subsequence** that is increasing / strictly increasing" | literal LIS |
| "longest **chain** / can nest / can stack" on pairs `(w, h)` | sort by one coordinate, LIS on the other (Russian Doll, envelopes, boxes) |
| "largest subset where every pair is divisible" | sort, then LIS where the "increasing" relation is `b % a == 0` (#368) |
| "increasing **triplet**" — just asks *does* a length-3 increasing subsequence exist | LIS with early exit; O(n)/O(1) two-variable trick (#334) |
| "minimum number of ... to cover / partition" | often **Dilworth**: min chains to cover = length of longest antichain, and vice-versa — many such problems become LIS on a transformed key |
| "longest **bitonic** / longest **wiggle**" | LIS run forwards *and* backwards, or a two-state DP |
| Patience / "smallest tail" phrasing, or you need better than O(n²) for n up to 10⁵ | the `tails` + binary-search method |

**When it's *not* LIS:** if the subsequence must be **contiguous**, that's a subarray problem (Kadane / sliding window / two pointers), not LIS. If you need the *count* of increasing subsequences of all lengths, that's a different (often Fenwick-tree) DP.

---

## General template / pseudocode

### 1. O(n²) DP — `dp[i]` = length of the LIS **ending exactly at index i**

The most intuitive version, and the one to reach for when you also need to **reconstruct** the subsequence.

```go
// lengthOfLISDP returns the LIS length in O(n^2) time, O(n) space.
//
// dp[i] = length of the longest strictly-increasing subsequence that ENDS at i.
// Every such subsequence's second-to-last element is some j < i with nums[j] < nums[i].
// Base: dp[i] = 1 (the element alone). Answer: max over all dp[i].
func lengthOfLISDP(nums []int) int {
    n := len(nums)
    if n == 0 {
        return 0
    }
    dp := make([]int, n)
    best := 1
    for i := 0; i < n; i++ {
        dp[i] = 1 // the single element nums[i] is always a valid length-1 subsequence
        for j := 0; j < i; j++ {
            if nums[j] < nums[i] && dp[j]+1 > dp[i] {
                // nums[i] can extend the best subsequence ending at j
                dp[i] = dp[j] + 1
            }
        }
        if dp[i] > best {
            best = dp[i]
        }
    }
    return best
}
```

### 2. O(n log n) — patience sorting with a `tails` array + binary search

Maintain `tails`, where `tails[k]` = **the smallest possible tail value of any increasing subsequence of length `k+1`** seen so far. `tails` is always sorted, so we can binary-search it.

For each `x`: find the **first** tail `≥ x` (strict LIS) and overwrite it with `x`; if none exists, `x` extends the longest run, so append it. The final `len(tails)` is the answer.

```go
// lengthOfLIS returns the strictly-increasing LIS length in O(n log n).
//
// tails[k] = smallest tail of an increasing subsequence of length k+1.
// Invariant: tails is strictly increasing, so binary search is valid.
// For strict LIS we replace the first element >= x (lower bound).
func lengthOfLIS(nums []int) int {
    tails := make([]int, 0, len(nums))
    for _, x := range nums {
        // Binary search for the leftmost index whose tail is >= x.
        lo, hi := 0, len(tails)
        for lo < hi {
            mid := lo + (hi-lo)/2
            if tails[mid] < x {
                lo = mid + 1 // tails[mid] can still precede x; go right
            } else {
                hi = mid // tails[mid] >= x is a candidate; go left
            }
        }
        if lo == len(tails) {
            tails = append(tails, x) // x is larger than all tails → extends LIS
        } else {
            tails[lo] = x // x becomes a smaller/equal tail for length lo+1
        }
    }
    return len(tails)
}
```

**Non-decreasing variant** (allow equal, `≤`): change the comparison from `tails[mid] < x` to `tails[mid] <= x` (i.e. use an **upper bound** — find the first tail *strictly greater* than `x`). That way an equal value extends rather than replaces.

Go's standard library gives this directly via `sort.SearchInts`:

```go
import "sort"

func lengthOfLIS(nums []int) int {
    tails := []int{}
    for _, x := range nums {
        i := sort.SearchInts(tails, x) // first index with tails[i] >= x  (strict LIS)
        if i == len(tails) {
            tails = append(tails, x)
        } else {
            tails[i] = x
        }
    }
    return len(tails)
}
```

### 3. Reconstructing the actual LIS (O(n²) DP with parent pointers)

`tails` alone can't be read back as a valid subsequence; to recover the elements, track a predecessor for each index.

```go
// reconstructLIS returns one longest strictly-increasing subsequence itself.
func reconstructLIS(nums []int) []int {
    n := len(nums)
    if n == 0 {
        return nil
    }
    dp := make([]int, n)    // dp[i] = LIS length ending at i
    prev := make([]int, n)  // prev[i] = index of the element before nums[i], or -1
    bestLen, bestEnd := 1, 0
    for i := 0; i < n; i++ {
        dp[i], prev[i] = 1, -1
        for j := 0; j < i; j++ {
            if nums[j] < nums[i] && dp[j]+1 > dp[i] {
                dp[i], prev[i] = dp[j]+1, j
            }
        }
        if dp[i] > bestLen {
            bestLen, bestEnd = dp[i], i
        }
    }
    // Walk the predecessor chain back from the best endpoint, then reverse.
    seq := make([]int, 0, bestLen)
    for i := bestEnd; i != -1; i = prev[i] {
        seq = append(seq, nums[i])
    }
    for l, r := 0, len(seq)-1; l < r; l, r = l+1, r-1 {
        seq[l], seq[r] = seq[r], seq[l]
    }
    return seq
}
```

(Reconstruction in O(n log n) is also possible by remembering, for each `x`, which `tails` position it landed in, plus a parent index — more bookkeeping, same idea.)

### 4. Counting the number of LISes (O(n²))

Track both a length and a count per index.

```go
func findNumberOfLIS(nums []int) int {
    n := len(nums)
    length := make([]int, n) // longest chain ending at i
    count := make([]int, n)  // how many such longest chains end at i
    maxLen, ans := 0, 0
    for i := 0; i < n; i++ {
        length[i], count[i] = 1, 1
        for j := 0; j < i; j++ {
            if nums[j] < nums[i] {
                if length[j]+1 > length[i] {
                    length[i] = length[j] + 1
                    count[i] = count[j] // new best length: inherit j's count
                } else if length[j]+1 == length[i] {
                    count[i] += count[j] // tie: accumulate ways
                }
            }
        }
        if length[i] > maxLen {
            maxLen, ans = length[i], count[i]
        } else if length[i] == maxLen {
            ans += count[i]
        }
    }
    return ans
}
```

### 5. LIS on pairs — the "sort then LIS" reduction (Russian Doll Envelopes #354)

Nest/stack problems become 1-D LIS *after* sorting, with a crucial tie-break trick:

```go
import "sort"

// maxEnvelopes: an envelope (w,h) nests in (W,H) iff w<W and h<H (strict both).
// Sort by width ascending; for EQUAL widths sort height DESCENDING so that two
// envelopes of the same width can never both appear in an increasing-by-height
// subsequence (they can't nest). Then the answer is the LIS over heights.
func maxEnvelopes(envelopes [][]int) int {
    sort.Slice(envelopes, func(i, j int) bool {
        if envelopes[i][0] == envelopes[j][0] {
            return envelopes[i][1] > envelopes[j][1] // equal width → height DESC
        }
        return envelopes[i][0] < envelopes[j][0] // width ASC
    })
    // LIS (strict) over the heights column.
    tails := []int{}
    for _, e := range envelopes {
        h := e[1]
        i := sort.SearchInts(tails, h) // strict: first tail >= h
        if i == len(tails) {
            tails = append(tails, h)
        } else {
            tails[i] = h
        }
    }
    return len(tails)
}
```

The height-descending tie-break is the whole trick: it collapses the 2-D strict-nesting constraint into a clean 1-D strict-LIS on heights.

---

## Worked example

### O(n²) DP on `nums = [10, 9, 2, 5, 3, 7, 101, 18]`

`dp[i]` = LIS length ending at index `i`. For each `i` we look at all earlier smaller elements.

| i | nums[i] | j with nums[j] < nums[i] (and their dp) | dp[i] | Running best |
|---|---------|------------------------------------------|-------|--------------|
| 0 | 10 | — | 1 | 1 |
| 1 | 9 | none (10 ≥ 9) | 1 | 1 |
| 2 | 2 | none | 1 | 1 |
| 3 | 5 | j=2 (nums 2, dp 1) | 1+1 = **2** | 2 |
| 4 | 3 | j=2 (nums 2, dp 1) | 1+1 = **2** | 2 |
| 5 | 7 | j=2(dp1), j=3(dp2), j=4(dp2) → best 2 | 2+1 = **3** | 3 |
| 6 | 101 | j=0,1,2,3,4,5 → best dp is dp[5]=3 | 3+1 = **4** | 4 |
| 7 | 18 | j=2,3,4,5 → best dp is dp[5]=3 | 3+1 = **4** | 4 |

Answer: **4**. Reconstructing from index 6 (`101`) via parents: `101 ← 7 ← 3 ← 2`, i.e. `[2, 3, 7, 101]`.

### O(n log n) patience sorting on the same input

We build `tails`, overwriting the first element `≥ x` (strict).

| x | Action (binary search for first tail ≥ x) | `tails` after |
|-----|--------------------------------------------|---------------|
| 10 | empty → append | `[10]` |
| 9 | first ≥ 9 is index 0 → replace | `[9]` |
| 2 | first ≥ 2 is index 0 → replace | `[2]` |
| 5 | none ≥ 5 → append | `[2, 5]` |
| 3 | first ≥ 3 is index 1 (value 5) → replace | `[2, 3]` |
| 7 | none ≥ 7 → append | `[2, 3, 7]` |
| 101 | none ≥ 101 → append | `[2, 3, 7, 101]` |
| 18 | first ≥ 18 is index 3 (value 101) → replace | `[2, 3, 7, 18]` |

`len(tails) = 4` → answer **4**. Note the final `tails = [2,3,7,18]` *happens* to be a real LIS here, but in general its contents are just "best tails", not a guaranteed subsequence — only its length is meaningful.

### Increasing Triplet (#334) — LIS specialised to length 3, O(n)/O(1)

Keep the two smallest "tails" for lengths 1 and 2; if any element beats both, a triplet exists.

```go
func increasingTriplet(nums []int) bool {
    first, second := math.MaxInt, math.MaxInt
    for _, x := range nums {
        switch {
        case x <= first:
            first = x // smallest so far (tail of a length-1 chain)
        case x <= second:
            second = x // smallest value that has some smaller element before it
        default:
            return true // x > second > (something before second) ⇒ triplet
        }
    }
    return false
}
```

This is exactly the `tails` idea frozen at size 2.

---

## Complexity

| Method | Time | Space | Notes |
|--------|------|-------|-------|
| O(n²) DP (`dp[i]`) | **O(n²)** | O(n) | simplest; needed for reconstruction & counting variants |
| Patience sorting + binary search | **O(n log n)** | O(n) | `n` elements, each an O(log n) binary search into `tails` |
| Increasing triplet (#334) | **O(n)** | **O(1)** | LIS pinned to length 3 → two scalars |
| Russian Doll (#354) | **O(n log n)** | O(n) | O(n log n) sort + O(n log n) LIS |
| Largest Divisible Subset (#368) | **O(n²)** | O(n) | divisibility isn't a total order compatible with the tails trick, so it stays O(n²) with reconstruction |

Why the fast method is O(n log n): `tails` stays sorted by construction (each write either appends a new max or lowers an existing slot without breaking order), so every element does one binary search — n × log n.

---

## Common pitfalls

1. **Strict vs non-decreasing mix-up.** Strict LIS replaces the first tail `≥ x` (`sort.SearchInts` / lower bound). Non-decreasing LIS replaces the first tail `> x` (upper bound). Using the wrong bound gives an answer off by the number of duplicate runs — the classic "why is my LIS one too long/short" bug.

2. **Treating `tails` as the actual subsequence.** Its *length* is the LIS length, but its contents are only "smallest possible tail per length" and often do **not** form a real subsequence of the input. To output the elements, use the O(n²) parent-pointer method (or the heavier O(n log n) reconstruction).

3. **Russian Doll: forgetting the height-descending tie-break (#354).** Sorting equal widths ascending by height lets two same-width envelopes both enter the height-LIS, falsely "nesting" them (widths equal, not strictly less). Sort equal widths by height **descending**.

4. **Off-by-one in the hand-rolled binary search.** Use the half-open `[lo, hi)` convention with `hi = len(tails)` and `mid = lo + (hi-lo)/2`; the loop condition is `lo < hi`. Mixing inclusive/exclusive bounds is the usual source of infinite loops or index-out-of-range.

5. **Divisibility LIS is not O(n log n)-able (#368).** The `tails` trick requires the relation to be a *total order* consistent with numeric comparison. Divisibility is only a partial order, so #368 stays O(n²) DP (sort first so any valid predecessor lies to the left, then `dp[i]` over divisors), and you must track parents to output the subset.

6. **Empty input / all-equal input.** LIS of `[]` is 0; LIS of `[7,7,7]` is 1 (strict) or 3 (non-decreasing). Make sure the base cases and the strict/non-strict choice agree with the problem.

7. **Increasing Triplet: the `<=` vs `<` on updates (#334).** Updating `first`/`second` with `<=` (not `<`) is what correctly handles duplicates so that `[2,2,2]` returns false. Getting this comparison wrong is the standard trap.

8. **Reconstruction tie-break ambiguity.** When multiple LISes exist, the parent-pointer walk returns *one* of them; if the problem expects the lexicographically smallest (or a specific one), add a tie-break rule when updating `prev`.

---

## Problems in this repo that use it

- [0300 — Longest Increasing Subsequence](/0300_longest_increasing_subsequence/README.md) — the base problem; both the O(n²) DP and the O(n log n) patience-sorting method.
- [0334 — Increasing Triplet Subsequence](/0334_increasing_triplet_subsequence/README.md) — LIS collapsed to length 3, solved in O(n) time / O(1) space with two scalars.
- [0354 — Russian Doll Envelopes](/0354_russian_doll_envelopes/README.md) — sort (width asc, height desc), then O(n log n) LIS on heights; the canonical "sort then LIS" reduction.
- [0368 — Largest Divisible Subset](/0368_largest_divisible_subset/README.md) — LIS where "increasing" is the divisibility relation; O(n²) DP with parent-pointer reconstruction.

### Related classics to know (may not be in repo)

- [0376 — Wiggle Subsequence](/0376_wiggle_subsequence/README.md) — a subsequence-length DP cousin (alternating up/down), solvable greedily in O(n).
- LeetCode #673 — Number of Longest Increasing Subsequences (the counting variant, template 4).
- LeetCode #646 — Maximum Length of Pair Chain (sort by second coordinate; greedy or LIS).
- LeetCode #1671 — Minimum Removals to Make Mountain Array (LIS forwards + backwards).
