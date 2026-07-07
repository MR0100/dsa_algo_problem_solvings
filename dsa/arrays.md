# Array

> **The** foundational data structure. Almost every other structure (stacks, heaps,
> hash tables, adjacency lists, DP tables) is built on top of a contiguous array.
> **Core power:** O(1) random access by index. **Core weakness:** O(n) insert/delete
> in the middle, and a fixed capacity you must grow explicitly.

---

## What it is

An **array** is a block of **contiguous memory** holding `n` elements of the same
type, addressed by an integer **index** `0 .. n-1`. Because the elements sit next
to each other and every element has the same size, the address of element `i` is
just `base + i*elementSize` — a single multiply-and-add. That is why indexing is
**O(1)** and why arrays are cache-friendly: scanning `a[0], a[1], a[2], …` walks
memory in order, so the CPU prefetcher keeps the pipeline full.

Where it sits among data structures:

- **vs. linked list** — array gives O(1) index access but O(n) middle insert/delete;
  a linked list is the mirror image (O(1) splice if you hold the node, O(n) to *find* it).
- **vs. hash map** — a map is "an array indexed by a hash of an arbitrary key". When
  your keys are already small integers `0..k`, skip the hash and use a plain array
  (a *bucket* / *direct-address table*) — it is faster and simpler (see #41, #268).
- **as a substrate** — heaps, Fenwick/segment trees, DP tables, ring buffers, and
  counting sort are all "an array plus an indexing rule".

### Array vs. slice in Go — the one thing to get right

Go distinguishes a fixed-size **array** from a growable **slice**, and 99% of
LeetCode Go code uses **slices**.

```go
var arr [5]int          // ARRAY: fixed length 5, length is part of the TYPE,
                        //        passed BY VALUE (copied) when handed to a func.

s := []int{1, 2, 3}     // SLICE: a 3-word header {ptr, len, cap} that VIEWS an
                        //        underlying array. Passed by value, but the header
                        //        points at shared backing storage.

s = make([]int, n)      // slice of n zeroed ints
s = make([]int, 0, n)   // len 0, cap n — pre-size to avoid re-allocations
s = append(s, 4)        // amortized O(1); grows (usually doubles) cap when full
```

Consequences you must internalise:

- **Slices share backing storage.** `b := a[1:3]` does **not** copy — writing `b[0]`
  mutates `a[1]`. To get an independent copy use `copy(dst, src)` or
  `append([]int{}, src...)`.
- **`append` may or may not reallocate.** If `cap` has room it writes in place and
  the original slice sees the change; if not, it allocates a new array and the two
  slices diverge. Never rely on which happened — either pre-size with `make`, or
  treat the returned slice as the source of truth: `s = append(s, x)`.
- **Passing a slice to a function** lets the function mutate existing elements (shared
  storage) but not change the caller's length — `append` inside the callee is invisible
  to the caller unless you return the slice.
- **2-D "arrays"** are slices of slices; each row must be allocated:
  `grid := make([][]int, rows); for i := range grid { grid[i] = make([]int, cols) }`.

---

## When to recognise it

"Use an array" is almost never the *whole* answer — the skill is recognising which
**array-scan idiom** the problem wants. Signals:

| Signal in the problem | Array technique |
|-----------------------|-----------------|
| "modify the array **in place**, O(1) extra space" | read/write **two-pointer** compaction (#26, #27, #283) |
| "**rotate** / reverse a portion" | index arithmetic or the **reversal trick** (#48, #189) |
| "running max/min/sum as you scan" | **single-pass state machine** / Kadane (#53, #121, #152) |
| "answer[i] depends on stuff **left and right** of i" | **two-pass** prefix + suffix scan (#238, #42) |
| "each value is a small integer in a known range" | use the array itself as a **bucket / sign-marker** (#41, #448) |
| "find pair/triple summing to target" in a **sorted** array | **opposite-end two pointers** (#11, #167) |
| "merge / compare two arrays from the back" | fill **from the end** to avoid overwrites (#88) |
| "count / group by index distance" | single scan tracking **last-seen index** (#243, #245) |
| "collapse consecutive runs into ranges" | scan tracking **start of current run** (#228, #163) |

Two meta-patterns dominate and are worth naming explicitly:

**Single-pass state machine.** Walk the array once, carrying a tiny bit of state
(a running best, a running sum, a last-seen index, a phase flag). O(n) time, O(1)
space. This subsumes Kadane's algorithm, "best time to buy/sell stock", and most
"last occurrence" scans.

**Two-pass scan.** When element `i`'s answer needs information from *both sides*, one
left-to-right pass computes prefixes, one right-to-left pass folds in suffixes.
"Product of array except self" and the DP formulation of "trapping rain water" are
the archetypes. This is the array face of the **prefix-sum** technique — see
[`/dsa/prefix_sum.md`](/dsa/prefix_sum.md) for the running-aggregate depth.

---

## General templates / idioms

### 1. In-place compaction (read/write two pointers)

The single most reused array idiom: `write` marks where the next *kept* element
goes; `read` scans everything. Keep the element ⇒ copy it to `write`, advance `write`.

```go
// removeElement deletes every occurrence of val in place and returns the new length.
// All kept elements end up in nums[:newLen]; the tail is garbage we ignore.
func removeElement(nums []int, val int) int {
    write := 0                     // next slot for a kept element
    for read := 0; read < len(nums); read++ {
        if nums[read] != val {     // keep this one
            nums[write] = nums[read]
            write++                // slot consumed
        }
        // else: skip — read advances, write stays put
    }
    return write                   // count of kept elements
}
```

The same skeleton, with the keep-condition changed, solves "remove duplicates from a
sorted array" (`nums[read] != nums[write-1]`), "move zeroes" (keep non-zero, then
zero-fill the tail), and "remove duplicates II" (allow up to two — compare against
`nums[write-2]`).

### 2. Running state machine (single pass)

```go
// maxSubArray — Kadane's algorithm. State = best sum of a subarray ENDING at i.
func maxSubArray(nums []int) int {
    best := nums[0]       // best answer seen anywhere
    cur := nums[0]        // best subarray ending at the current index
    for i := 1; i < len(nums); i++ {
        // extend the previous run, or start fresh at nums[i] — whichever is larger
        cur = max(cur+nums[i], nums[i])
        best = max(best, cur)
    }
    return best
}
```

### 3. Two-pass prefix/suffix (answer needs both sides)

```go
// productExceptSelf — out[i] = product of everything except nums[i], no division.
// Pass 1 fills out[i] with the product of the PREFIX (everything left of i).
// Pass 2 multiplies in the product of the SUFFIX (everything right of i) on the fly.
func productExceptSelf(nums []int) []int {
    n := len(nums)
    out := make([]int, n)
    out[0] = 1
    for i := 1; i < n; i++ {
        out[i] = out[i-1] * nums[i-1] // product of nums[0..i-1]
    }
    suffix := 1
    for i := n - 1; i >= 0; i-- {
        out[i] *= suffix              // fold in product of nums[i+1..n-1]
        suffix *= nums[i]             // extend the suffix leftward
    }
    return out
}
```

### 4. The array as its own hash table (index = key)

When values are integers in `[1..n]` (or `[0..n-1]`), you don't need a separate map —
the array's own indices *are* the keys. Two flavours:

```go
// Cyclic sort: put each value v at index v-1 by swapping until it's home.
// Used by #41 First Missing Positive, #448 Find All Numbers Disappeared.
for i := 0; i < len(nums); i++ {
    for nums[i] >= 1 && nums[i] <= len(nums) && nums[nums[i]-1] != nums[i] {
        nums[i], nums[nums[i]-1] = nums[nums[i]-1], nums[i] // send nums[i] home
    }
}
// Afterwards, the first index i where nums[i] != i+1 exposes the missing value.
```

```go
// Sign marking: encode "I have seen value v" by flipping the sign of nums[v-1].
// O(1) extra space because the "seen" set lives inside the array itself.
for _, x := range nums {
    idx := abs(x) - 1
    if nums[idx] > 0 {
        nums[idx] = -nums[idx] // mark presence of value (idx+1)
    }
}
```

### 5. The reversal trick (rotate without extra array)

```go
// rotate shifts nums right by k using three reversals — O(1) extra space.
// [1,2,3,4,5,6,7], k=3  ->  reverse all -> reverse first k -> reverse rest.
func rotate(nums []int, k int) {
    n := len(nums)
    k %= n                       // k may exceed n
    reverse(nums, 0, n-1)        // 7 6 5 4 3 2 1
    reverse(nums, 0, k-1)        // 5 6 7 | 4 3 2 1
    reverse(nums, k, n-1)        // 5 6 7 | 1 2 3 4
}
func reverse(a []int, lo, hi int) {
    for lo < hi {
        a[lo], a[hi] = a[hi], a[lo]
        lo, hi = lo+1, hi-1
    }
}
```

### 6. Fill-from-the-back (merge without a scratch buffer)

```go
// merge writes nums1 (which has trailing space) from the RIGHT so we never
// overwrite an nums1 element we still need to read. Classic #88.
func merge(nums1 []int, m int, nums2 []int, n int) {
    i, j, w := m-1, n-1, m+n-1    // read heads for each input; write head at the end
    for j >= 0 {                 // nums2 not exhausted
        if i >= 0 && nums1[i] > nums2[j] {
            nums1[w] = nums1[i]; i--
        } else {
            nums1[w] = nums2[j]; j--
        }
        w--
    }
}
```

---

## Worked example — two-pointer in-place removal (#27)

Input: `nums = [3, 2, 2, 3]`, `val = 3`. Expected: length `2`, `nums[:2] == [2, 2]`.

Trace of `removeElement` (`w` = write, values marked `*` are freshly written):

| read | nums[read] | keep? | action | write after | array state |
|------|-----------|-------|--------|-------------|-------------|
| —    | —          | —     | init   | 0 | `[3, 2, 2, 3]` |
| 0    | 3          | no    | skip   | 0 | `[3, 2, 2, 3]` |
| 1    | 2          | yes   | `nums[0]=2` | 1 | `[2*, 2, 2, 3]` |
| 2    | 2          | yes   | `nums[1]=2` | 2 | `[2, 2*, 2, 3]` |
| 3    | 3          | no    | skip   | 2 | `[2, 2, 2, 3]` |

Return `write = 2`. The judge reads only `nums[:2] = [2, 2]`; the tail `[2, 3]` is
leftover garbage and correctly ignored. **Key insight:** `write` never outruns `read`,
so we always copy *forward into space we have already scanned* — no element is lost.

### Worked example — two-pass product (#238)

Input `nums = [1, 2, 3, 4]`, expected `[24, 12, 8, 6]`.

Pass 1 (prefixes) fills `out[i] = product of nums[0..i-1]`:

| i | out[i] = out[i-1]*nums[i-1] | out so far |
|---|------------------------------|------------|
| 0 | seed = 1                     | `[1, _, _, _]` |
| 1 | 1 * nums[0]=1                | `[1, 1, _, _]` |
| 2 | 1 * nums[1]=2                | `[1, 1, 2, _]` |
| 3 | 2 * nums[2]=3                | `[1, 1, 2, 6]` |

Pass 2 (suffixes), `suffix` starts at 1:

| i | out[i] *= suffix | suffix *= nums[i] | out so far |
|---|------------------|-------------------|------------|
| 3 | 6 * 1  = 6       | 1 * 4 = 4         | `[1, 1, 2, 6]` |
| 2 | 2 * 4  = 8       | 4 * 3 = 12        | `[1, 1, 8, 6]` |
| 1 | 1 * 12 = 12      | 12 * 2 = 24       | `[1, 12, 8, 6]` |
| 0 | 1 * 24 = 24      | 24 * 1 = 24       | `[24, 12, 8, 6]` |

Result `[24, 12, 8, 6]`. Each element saw the product of everything to its left
(pass 1) times everything to its right (pass 2) — and never itself.

---

## Complexity

| Operation | Array / slice | Why |
|-----------|---------------|-----|
| Index read/write `a[i]` | **O(1)** | address = base + i·size |
| Scan / traverse | **O(n)** | touch each element once |
| Search (unsorted) | **O(n)** | must inspect every element |
| Search (sorted) | **O(log n)** | binary search — see [`/dsa/binary_search.md`](/dsa/binary_search.md) |
| Insert / delete at **end** (slice) | **O(1)** amortized | `append`; occasional O(n) regrow |
| Insert / delete in **middle** | **O(n)** | shift every following element |
| `append` regrow | **O(n)** one-off, **O(1)** amortized | doubling ⇒ total copies ≤ 2n |

The idioms above are all **O(n) time**. Their whole point is **O(1) extra space** —
they mutate the input array in place instead of allocating a second one. When a
problem says "do it without extra space" or "O(1) space", it is asking for one of the
in-place idioms (compaction, sign-marking, cyclic sort, reversal).

---

## Common pitfalls

1. **Off-by-one on bounds.** The valid index range is `0 .. len-1`. `for i := 0; i <= len(a); i++` reads one past the end and panics. Prefer `range` when you don't need the index arithmetic.

2. **Aliasing via slices.** `b := a[i:j]` shares storage with `a`. Mutating `b` mutates `a`, and vice-versa. When you need an independent copy, `copy()` or `append([]T{}, a...)` — never assume slicing copies.

3. **`append` reallocation surprises.** After `b := append(a, x)`, `a` and `b` may or may not share backing memory depending on capacity. Writing through one may or may not be seen by the other. Always reassign: `a = append(a, x)`.

4. **Modifying a slice while ranging over it.** `for i, v := range nums` evaluates `nums` once; growing it inside the loop won't extend the iteration, and shrinking it can leave `v` stale. For in-place compaction use an explicit `read`/`write` index loop, not `range`.

5. **Overwriting data you still need.** In-place merges and shifts can clobber an element before it's read. The fix is direction: fill **from the back** (#88) or move the write pointer so it never passes the read pointer (#26/#27).

6. **Nil vs. empty slice.** A `nil` slice has `len 0` and is safe to `range` and `append`, but `s == nil` is true only for the nil case, not for `[]int{}`. Don't branch on emptiness with `== nil`; use `len(s) == 0`.

7. **2-D slice rows sharing a backing array.** `row := make([]int, n); grid := [][]int{row, row}` makes both rows the *same* slice — editing one edits the other. Allocate each row separately.

8. **Integer overflow in index math / products.** `(lo+hi)/2` can overflow for huge indices (use `lo + (hi-lo)/2`); running products (#238, #152) can overflow — mind the value range.

9. **Assuming sorted when it isn't.** Two-pointer-from-both-ends and binary search require a sorted array. If the problem doesn't guarantee order, you must sort first (adds O(n log n)) or pick a different technique.

---

## Problems in this repo that use it

In-place two-pointer compaction:
- [0026 — Remove Duplicates from Sorted Array](/0026_remove_duplicates_from_sorted_array/README.md) — keep when `nums[read] != nums[write-1]`
- [0027 — Remove Element](/0027_remove_element/README.md) — the canonical read/write skeleton (traced above)
- [0080 — Remove Duplicates from Sorted Array II](/0080_remove_duplicates_from_sorted_array_ii/README.md) — allow two copies via `nums[write-2]`
- [0283 — Move Zeroes](/0283_move_zeroes/README.md) — compact non-zeros, then zero-fill the tail
- [0075 — Sort Colors](/0075_sort_colors/README.md) — Dutch-national-flag three-way partition, one pass

Single-pass running state:
- [0053 — Maximum Subarray](/0053_maximum_subarray/README.md) — Kadane's running best (traced above)
- [0121 — Best Time to Buy and Sell Stock](/0121_best_time_to_buy_and_sell_stock/README.md) — track running min price, best profit
- [0152 — Maximum Product Subarray](/0152_maximum_product_subarray/README.md) — carry running max *and* min (signs flip)
- [0169 — Majority Element](/0169_majority_element/README.md) — Boyer–Moore vote counter, O(1) space

Two-pass / prefix-suffix:
- [0238 — Product of Array Except Self](/0238_product_of_array_except_self/README.md) — prefix pass then suffix pass (traced above)
- [0042 — Trapping Rain Water](/0042_trapping_rain_water/README.md) — max-to-the-left and max-to-the-right per column
- [0011 — Container With Most Water](/0011_container_with_most_water/README.md) — opposite-end two pointers

Array-as-hash / index=key:
- [0041 — First Missing Positive](/0041_first_missing_positive/README.md) — cyclic sort into `nums[v-1]`, O(1) space
- [0073 — Set Matrix Zeroes](/0073_set_matrix_zeroes/README.md) — use first row/column as the marker store

Index arithmetic / rotation / merge:
- [0048 — Rotate Image](/0048_rotate_image/README.md) — transpose then reverse each row, in place
- [0189 — Rotate Array](/0189_rotate_array/README.md) — three-reversal trick (see template 5)
- [0066 — Plus One](/0066_plus_one/README.md) — carry propagation scanning from the least-significant digit
- [0088 — Merge Sorted Array](/0088_merge_sorted_array/README.md) — fill from the back (template 6)

Scan tracking last-seen / run start:
- [0243 — Shortest Word Distance](/0243_shortest_word_distance/README.md) — one pass, remember last index of each word
- [0245 — Shortest Word Distance III](/0245_shortest_word_distance_iii/README.md) — same-word variant of the last-index scan
- [0228 — Summary Ranges](/0228_summary_ranges/README.md) — track the start of each consecutive run
- [0163 — Missing Ranges](/0163_missing_ranges/README.md) — scan gaps between consecutive elements

### Foundational — see also

- [`/dsa/two_pointers.md`](/dsa/two_pointers.md) — the read/write and opposite-end pointer patterns in depth
- [`/dsa/prefix_sum.md`](/dsa/prefix_sum.md) — running aggregates as a first-class technique
- [`/dsa/sliding_window.md`](/dsa/sliding_window.md) — contiguous-subarray scans with a moving window
- [`/dsa/sorting.md`](/dsa/sorting.md) and [`/dsa/binary_search.md`](/dsa/binary_search.md) — what a *sorted* array unlocks
- [`/dsa/matrix_traversal.md`](/dsa/matrix_traversal.md) — the 2-D generalisation
