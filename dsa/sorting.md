# Sorting

> **Category:** Algorithmic technique / preprocessing tool
> **Typical complexity:** O(n log n) comparison sort · O(n + k) counting sort · O(1)–O(n) extra space

---

## What it is

Sorting rearranges a collection into a defined order (ascending, descending, or
by a custom key). On LeetCode, sorting is rarely the *answer* itself — it is a
**preprocessing step that buys structure**. Once data is sorted you gain three
superpowers:

1. **Adjacency** — equal or related elements sit next to each other
   (duplicate detection, grouping, interval overlap).
2. **Monotonicity** — values only grow as you scan, enabling two pointers,
   binary search, and greedy decisions.
3. **Canonical form** — a sorted version of an object is a unique fingerprint
   for its equivalence class (anagrams, multisets).

Comparison-based sorting has a proven lower bound of **Ω(n log n)**. Beating it
requires exploiting value constraints (counting sort, bucket sort, radix sort,
cyclic sort / index-as-hash) — a recurring trick when a problem demands O(n).

---

## How to recognise a sorting problem

Signals in the problem statement:

| Signal | Why sorting helps |
|--------|-------------------|
| "Find pairs/triplets that sum to…" | Sort → two pointers converging from both ends |
| "Return results **without duplicates**" | Sort → duplicates become adjacent → skip with `nums[i] == nums[i-1]` |
| "Merge / overlap / intervals" | Sort by start → only adjacent intervals can overlap |
| "Group items that are equivalent" (anagrams, multisets) | Sorted form = canonical hash key |
| "K-th largest / smallest", "top k", "closest" | Sort (or partial sort / heap / quickselect) |
| "Array contains values in range [1, n]" or a small fixed alphabet | Counting sort / cyclic sort / Dutch national flag → O(n) |
| "Input is **already sorted**" | Don't sort — the problem wants binary search or two pointers |
| "Must run in O(n)" but the natural solution sorts | Look for a non-comparison sort or a hash-based shortcut |
| "Sort by custom criterion" (meeting times, tasks, envelopes) | `sort.Slice` with a custom `less` + greedy/DP on the sorted order |

Anti-signal: if the problem needs **original indices** after sorting, either
sort `(value, index)` pairs or switch to a hash map (see LeetCode #1).

---

## Go templates

### 1. The standard library (what you write 95% of the time)

```go
import "sort"

// Sort a slice of ints ascending.
sort.Ints(nums)

// Sort with a custom comparator (e.g., intervals by start).
sort.Slice(intervals, func(i, j int) bool {
    return intervals[i][0] < intervals[j][0] // less() — true if i must come before j
})

// Stable variant: preserves relative order of equal elements
// (needed when a secondary implicit order matters).
sort.SliceStable(items, func(i, j int) bool { return items[i].key < items[j].key })

// Sort a string: convert to byte/rune slice first (strings are immutable in Go).
b := []byte(s)
sort.Slice(b, func(i, j int) bool { return b[i] < b[j] })
canonical := string(b)

// Sort (value, originalIndex) pairs when indices must survive the sort.
type pair struct{ val, idx int }
pairs := make([]pair, len(nums))
for i, v := range nums { pairs[i] = pair{v, i} }
sort.Slice(pairs, func(i, j int) bool { return pairs[i].val < pairs[j].val })
```

Go's `sort.Slice` uses pattern-defeating quicksort (pdqsort since Go 1.19):
O(n log n) worst case, in-place, **not stable**. `sort.SliceStable` is stable
at the cost of O(log n) extra merge passes.

### 2. Merge sort (know it cold — it's the merge step that interviews test)

```go
// mergeSort sorts nums using divide and conquer.
//
// Time:  O(n log n) — log n levels, O(n) merge work per level.
// Space: O(n) — temporary buffer for merging.
func mergeSort(nums []int) []int {
    if len(nums) <= 1 {          // base case: 0 or 1 element is already sorted
        return nums
    }
    mid := len(nums) / 2
    left := mergeSort(nums[:mid])  // sort left half
    right := mergeSort(nums[mid:]) // sort right half
    return merge(left, right)      // combine two sorted halves
}

// merge combines two sorted slices into one sorted slice.
// This exact pattern reappears in: Merge Two Sorted Lists (#21),
// Merge Sorted Array (#88), Merge k Sorted Lists (#23), count-inversions.
func merge(a, b []int) []int {
    out := make([]int, 0, len(a)+len(b))
    i, j := 0, 0
    for i < len(a) && j < len(b) {
        if a[i] <= b[j] {          // <= keeps the sort stable
            out = append(out, a[i]); i++
        } else {
            out = append(out, b[j]); j++
        }
    }
    out = append(out, a[i:]...)    // drain leftovers (only one of these is non-empty)
    out = append(out, b[j:]...)
    return out
}
```

### 3. Quicksort + partition (basis of quickselect for "k-th element")

```go
// quickSort sorts nums in place.
//
// Time:  O(n log n) average, O(n²) worst (mitigate with random pivot).
// Space: O(log n) recursion stack.
func quickSort(nums []int, lo, hi int) {
    if lo >= hi { return }
    p := partition(nums, lo, hi)  // pivot lands at its final sorted position p
    quickSort(nums, lo, p-1)      // recurse left of pivot
    quickSort(nums, p+1, hi)      // recurse right of pivot
}

// partition (Lomuto): everything < pivot moves left of it.
func partition(nums []int, lo, hi int) int {
    pivot := nums[hi]             // pick last element as pivot (randomise in practice)
    i := lo                       // i = next slot for a "small" element
    for j := lo; j < hi; j++ {
        if nums[j] < pivot {      // small element found
            nums[i], nums[j] = nums[j], nums[i]
            i++
        }
    }
    nums[i], nums[hi] = nums[hi], nums[i] // place pivot into its slot
    return i
}
```

### 4. Counting sort (O(n + k) when values live in a small range)

```go
// countingSort sorts values known to be in [0, k].
//
// Time:  O(n + k)   Space: O(k)
func countingSort(nums []int, k int) {
    count := make([]int, k+1)     // count[v] = occurrences of value v
    for _, v := range nums { count[v]++ }
    i := 0
    for v := 0; v <= k; v++ {     // rewrite nums in order
        for c := 0; c < count[v]; c++ {
            nums[i] = v; i++
        }
    }
}
```

### 5. Dutch national flag — one-pass 3-way partition (LeetCode #75)

```go
// low  = boundary of 0s, high = boundary of 2s, mid = scanner.
low, mid, high := 0, 0, len(nums)-1
for mid <= high {
    switch nums[mid] {
    case 0: nums[low], nums[mid] = nums[mid], nums[low]; low++; mid++
    case 1: mid++                       // 1s stay in the middle
    case 2: nums[mid], nums[high] = nums[high], nums[mid]; high--
        // note: mid does NOT advance — the swapped-in value is unexamined
    }
}
```

### 6. Cyclic sort / index-as-hash (values in [1, n] → O(n), O(1) space)

```go
// Put value v at index v-1; anything out of place after this pass is the answer.
// Used in First Missing Positive (#41), Find Duplicate, Find Missing.
for i := 0; i < len(nums); i++ {
    for nums[i] >= 1 && nums[i] <= len(nums) && nums[nums[i]-1] != nums[i] {
        nums[nums[i]-1], nums[i] = nums[i], nums[nums[i]-1] // swap v home
    }
}
```

### Algorithm cheat sheet

| Algorithm | Time (avg / worst) | Space | Stable | When |
|-----------|--------------------|-------|--------|------|
| Merge sort | O(n log n) / O(n log n) | O(n) | Yes | Linked lists, external sort, count inversions |
| Quicksort | O(n log n) / O(n²) | O(log n) | No | In-place general sort; partition → quickselect |
| Heap sort | O(n log n) / O(n log n) | O(1) | No | Guaranteed bound + O(1) space; top-k via heap |
| Insertion sort | O(n²) / O(n²) | O(1) | Yes | Tiny or nearly-sorted input; sorted-list insert |
| Counting sort | O(n + k) | O(k) | Yes* | Small value range (colors, chars, ratings) |
| Bucket sort | O(n) avg | O(n) | — | Uniformly distributed floats/frequencies |
| Radix sort | O(d·(n + k)) | O(n + k) | Yes | Fixed-width ints / strings |
| Cyclic sort | O(n) | O(1) | No | Values in [1, n]; missing/duplicate numbers |

---

## Worked example — Merge Intervals (LeetCode #56)

Problem: given `intervals = [[1,3],[2,6],[8,10],[15,18]]`, merge all
overlapping intervals.

**Key insight:** unsorted, *any* pair might overlap (O(n²) checks). Sorted by
start, an interval can only overlap its immediate predecessor — one O(n) pass.

```go
sort.Slice(intervals, func(i, j int) bool { return intervals[i][0] < intervals[j][0] })

merged := [][]int{intervals[0]}
for _, cur := range intervals[1:] {
    last := merged[len(merged)-1]
    if cur[0] <= last[1] {                 // overlap: cur starts before last ends
        if cur[1] > last[1] { last[1] = cur[1] } // extend the merged interval
    } else {
        merged = append(merged, cur)       // gap: start a new interval
    }
}
```

Step-by-step trace (input already happens to be sorted here):

| Step | `cur` | `last` (top of merged) | `cur[0] <= last[1]`? | Action | `merged` after |
|------|-------|------------------------|----------------------|--------|----------------|
| init | —     | —                      | —                    | seed with first | `[[1,3]]` |
| 1    | [2,6] | [1,3]                  | 2 ≤ 3 → yes          | extend end to max(3,6)=6 | `[[1,6]]` |
| 2    | [8,10]| [1,6]                  | 8 ≤ 6 → no           | append | `[[1,6],[8,10]]` |
| 3    | [15,18]| [8,10]                | 15 ≤ 10 → no         | append | `[[1,6],[8,10],[15,18]]` |

Result: `[[1,6],[8,10],[15,18]]`. Total cost O(n log n) for the sort +
O(n) for the sweep.

The same "sort, then one linear pass" shape solves meeting rooms, non-overlapping
intervals, insert interval, and most greedy scheduling problems.

---

## Common pitfalls

1. **Destroying original indices.** Sorting reorders the array; if the answer
   must be reported in terms of original positions (Two Sum #1), sort
   `(value, index)` pairs or avoid sorting entirely.
2. **Sorting when the input is already sorted.** "Sorted array" in the
   statement means the intended solution is binary search or two pointers —
   re-sorting wastes the constraint (and O(n log n) when O(log n) was expected).
3. **Forgetting stability.** `sort.Slice` is NOT stable. If equal elements
   must keep their relative order (secondary criteria, ties broken by
   appearance order), use `sort.SliceStable` or add a tiebreaker to `less`.
4. **Duplicate-skip without sorting first.** The classic
   `if i > start && nums[i] == nums[i-1] { continue }` pruning
   (3Sum #15, Combination Sum II #40, Permutations II #47) only works because
   sorting made duplicates adjacent. Skipping the sort silently produces
   duplicate results.
5. **Comparator bugs.** In Go, `less(i, j)` must be a strict weak ordering —
   returning `true` for equal elements can panic or corrupt the sort. Use `<`,
   never `<=`. For descending order flip to `>`, don't negate `<=`.
6. **Integer overflow in comparators.** `return a[i]-a[j] < 0` overflows for
   extreme ints; compare directly with `<`.
7. **Assuming O(n log n) is optimal.** If constraints say values ∈ [1, n], a
   small alphabet, or "must be O(n)", reach for counting / cyclic / bucket
   sort or the Dutch national flag instead of a comparison sort (#41, #75).
8. **Sorting strings the wrong way.** `sort.Slice` on a `string` directly
   won't compile — convert to `[]byte`/`[]rune` first; and rune-sort if the
   input can contain multi-byte characters.
9. **Quicksort worst case in adversarial inputs.** Sorted or all-equal input
   with a fixed pivot is O(n²); randomise the pivot or use median-of-three.
   (Go's pdqsort already handles this — a reason to prefer the stdlib.)
10. **Merging in place from the front.** Merge Sorted Array (#88) overwrites
    unread elements if you merge left-to-right; merge **backwards from the
    largest** to use the free space at the end.

---

## Problems in this repo

Problems whose solutions use sorting as a core technique (0131+ to be added in
a later pass):

- [0001 — Two Sum](../0001_two_sum/README.md) — sort `(value, index)` pairs, then two pointers; shows how to preserve original indices through a sort.
- [0014 — Longest Common Prefix](../0014_longest_common_prefix/README.md) — lexicographic sort reduces the LCP of n strings to LCP(first, last).
- [0015 — 3Sum](../0015_3sum/README.md) — sort to enable converging two pointers and adjacent-duplicate skipping.
- [0016 — 3Sum Closest](../0016_3sum_closest/README.md) — sort, fix one element, two-pointer sweep tracking the closest sum.
- [0018 — 4Sum](../0018_4sum/README.md) — sort, two nested fixed indices + two pointers, duplicate skipping at every level.
- [0023 — Merge k Sorted Lists](../0023_merge_k_sorted_lists/README.md) — the k-way merge step of (external) merge sort, via divide-and-conquer or a min-heap.
- [0039 — Combination Sum](../0039_combination_sum/README.md) — sort candidates so `break`-pruning cuts impossible branches early.
- [0040 — Combination Sum II](../0040_combination_sum_ii/README.md) — sort so duplicate candidates are adjacent, enabling skip-duplicate pruning.
- [0041 — First Missing Positive](../0041_first_missing_positive/README.md) — sorting baseline vs. O(n) cyclic sort (index-as-hash).
- [0047 — Permutations II](../0047_permutations_ii/README.md) — sort first so the `nums[i]==nums[i-1]` duplicate-skip rule is reliable.
- [0049 — Group Anagrams](../0049_group_anagrams/README.md) — sorted string as canonical form / hash key for each anagram class.
- [0056 — Merge Intervals](../0056_merge_intervals/README.md) — sort by start so only adjacent intervals can overlap (worked example above).
- [0075 — Sort Colors](../0075_sort_colors/README.md) — counting sort and one-pass Dutch national flag 3-way partition.
- [0086 — Partition List](../0086_partition_list/README.md) — stable two-list partition of a linked list by a pivot value (quicksort's partition idea, stable form).
- [0088 — Merge Sorted Array](../0088_merge_sorted_array/README.md) — in-place backwards merge, the merge step of merge sort under a space constraint.

Related references: [`two_pointers.md`](./two_pointers.md) (what sorting most
often enables), [`hash_map.md`](./hash_map.md) (the usual O(n) alternative to
sort-based approaches).
