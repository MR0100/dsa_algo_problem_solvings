# Non-Comparison Sorts — Counting, Bucket, Radix

> **The big idea:** comparison sorts (quicksort, mergesort, heapsort) are provably bounded by **Ω(n log n)** because they only learn about the data through pairwise `<` comparisons. If instead the keys are **bounded integers** (or map cleanly onto array indices / small ranges), you can *index directly* and sort in **O(n + k)** — no comparisons at all.

---

## What it is

A **comparison sort** decides order solely by asking "is a < b?". Any such algorithm can be modelled as a decision tree; sorting `n` items means distinguishing `n!` permutations, and a binary tree tall enough to have `n!` leaves needs height `≥ log₂(n!) = Ω(n log n)`. That is a hard lower bound — you *cannot* beat `n log n` while only comparing.

**Non-comparison (distribution) sorts** sidestep the bound by using the keys themselves as information:

- **Counting sort** — for integer keys in a known range `[0, k)`: count how many times each key occurs, then reconstruct the sorted output from the counts. `O(n + k)`.
- **Bucket sort** — scatter elements into `m` buckets by some key-derived index (e.g. value range, or a hash), sort inside each bucket, then concatenate. `O(n + m)` when the elements distribute evenly and buckets stay small.
- **Radix sort** — sort multi-digit keys one digit at a time, least-significant-digit first, using a **stable** counting sort as the per-digit subroutine. `O(d·(n + b))` for `d` digits in base `b`.

They all trade generality for speed: you must be able to turn a key into a bucket index cheaply, and the range/spread of keys must be controlled, or the space cost blows up.

### Where they sit

| | Comparison sorts | Distribution sorts |
|--|------------------|--------------------|
| Information used | pairwise `<` | the key value itself → an index |
| Lower bound | Ω(n log n) | none from comparisons |
| Typical time | O(n log n) | O(n + k), O(d·(n+b)) |
| Needs | any `Less` ordering | bounded integer-like keys |
| Stable? | mergesort yes; quicksort no | counting/radix yes (if careful); bucket depends on inner sort |

---

## When to recognise it

| Signal in the problem | Which sort | Why |
|-----------------------|-----------|-----|
| Keys are integers in a **small, known range** (0..k) | Counting sort | index by value directly, O(n+k) |
| Only a handful of distinct values (e.g. 0/1/2) | Counting sort | k is tiny → effectively O(n) |
| "sort in O(n)" hint with bounded values | Counting / radix | comparison sort can't hit O(n) |
| **Frequencies** matter and you want "top-k frequent" | Bucket sort by frequency | frequency ∈ [1, n] → n+1 buckets, no heap needed |
| Values (or their transform) spread **uniformly** over a range | Bucket sort | even spread ⇒ O(1) per bucket average |
| **Fixed-width** integer keys / large integers | Radix (LSD) | avoid k = huge range; process digit by digit |
| "maximum gap between consecutive sorted values in O(n)" | Bucket / radix | pigeonhole buckets give the gap without full sort |
| Counting-based ranking (e.g. H-Index over a bounded score) | Counting sort | cap the range at n, tally, sweep |

**When *not* to use them:** keys are arbitrary comparables (strings by locale, floats with huge dynamic range, custom objects with only a `Less`), the value range `k` dwarfs `n` (space explodes), or you genuinely need an in-place, comparison-based sort. Then fall back to `sort.Slice` / mergesort / quickselect.

---

## General template / pseudocode

### Counting sort (stable, general integer keys)

```go
// countingSort returns nums sorted ascending, assuming every value is in [0, k).
// Stable: equal keys keep their original relative order (needed for radix).
// Time: O(n + k)   Space: O(n + k)
func countingSort(nums []int, k int) []int {
    count := make([]int, k) // count[v] = how many times v appears
    for _, v := range nums {
        count[v]++
    }
    // Prefix-sum the counts so count[v] becomes the number of elements ≤ v.
    // After this, count[v] is the *end index (exclusive)* for value v.
    for v := 1; v < k; v++ {
        count[v] += count[v-1]
    }
    out := make([]int, len(nums))
    // Walk right-to-left to keep the sort STABLE: later equal elements land
    // at higher indices, preserving input order.
    for i := len(nums) - 1; i >= 0; i-- {
        v := nums[i]
        count[v]--           // the slot just before the current end
        out[count[v]] = v
    }
    return out
}
```

### Counting sort, in-place *overwrite* (when you don't need stability — e.g. Sort Colors)

```go
// countingSortInPlace overwrites nums in ascending order for values in [0, k).
// Not stable, but O(k) extra space and dead simple.
func countingSortInPlace(nums []int, k int) {
    count := make([]int, k)
    for _, v := range nums {
        count[v]++
    }
    idx := 0
    for v := 0; v < k; v++ {
        for count[v] > 0 { // write v as many times as it occurred
            nums[idx] = v
            idx++
            count[v]--
        }
    }
}
```

### Bucket-by-frequency (the "top-k frequent" pattern)

```go
// topKFrequent returns the k most frequent values using bucket sort by count.
// Frequencies range 1..len(nums), so index buckets by frequency directly.
// Time: O(n)   Space: O(n)
func topKFrequent(nums []int, k int) []int {
    freq := make(map[int]int)
    for _, v := range nums {
        freq[v]++
    }
    // buckets[f] holds all values that occur exactly f times.
    buckets := make([][]int, len(nums)+1)
    for v, f := range freq {
        buckets[f] = append(buckets[f], v)
    }
    res := make([]int, 0, k)
    // Sweep from highest frequency down; collect until we have k values.
    for f := len(buckets) - 1; f >= 1 && len(res) < k; f-- {
        res = append(res, buckets[f]...)
    }
    return res[:k]
}
```

### Radix sort (LSD, base 10, non-negative ints)

```go
// radixSortLSD sorts non-negative ints least-significant-digit first,
// using a stable counting sort per digit.
// Time: O(d·(n + 10))   Space: O(n)
func radixSortLSD(nums []int) []int {
    if len(nums) == 0 {
        return nums
    }
    max := nums[0]
    for _, v := range nums {
        if v > max {
            max = v
        }
    }
    out := append([]int(nil), nums...)
    for exp := 1; max/exp > 0; exp *= 10 { // one pass per digit position
        count := make([]int, 10)
        for _, v := range out {
            count[(v/exp)%10]++ // tally this digit
        }
        for d := 1; d < 10; d++ {
            count[d] += count[d-1] // prefix sums → end indices
        }
        tmp := make([]int, len(out))
        for i := len(out) - 1; i >= 0; i-- { // right-to-left ⇒ stable
            d := (out[i] / exp) % 10
            count[d]--
            tmp[count[d]] = out[i]
        }
        out = tmp
    }
    return out
}
```

---

## Worked example — step-by-step trace

### Counting sort on `nums = [2, 5, 3, 0, 2, 3, 0, 3]`, k = 6 (values 0..5)

**1. Tally** each value:

```
value:  0  1  2  3  4  5
count:  2  0  2  3  0  1        (two 0s, two 2s, three 3s, one 5)
```

**2. Prefix-sum** → `count[v]` = number of elements ≤ v = exclusive end index of value `v`:

```
value:  0  1  2  3  4  5
count:  2  2  4  7  7  8
```

**3. Place** (right-to-left for stability). Read input from the end:

| i | nums[i]=v | count[v] before | write out[count[v]-1] | count[v] after |
|---|-----------|-----------------|-----------------------|----------------|
| 7 | 3 | 7 | out[6]=3 | 6 |
| 6 | 0 | 2 | out[1]=0 | 1 |
| 5 | 3 | 6 | out[5]=3 | 5 |
| 4 | 2 | 4 | out[3]=2 | 3 |
| 3 | 0 | 1 | out[0]=0 | 0 |
| 2 | 3 | 5 | out[4]=3 | 4 |
| 1 | 5 | 8 | out[7]=5 | 7 |
| 0 | 2 | 3 | out[2]=2 | 2 |

Result: `out = [0, 0, 2, 2, 3, 3, 3, 5]` — sorted, in O(n + k).

### Bucket-by-frequency on `nums = [1,1,1,2,2,3]`, k = 2

Frequencies: `{1:3, 2:2, 3:1}`. Buckets indexed by frequency (size n+1 = 7):

```
freq:    0   1    2    3    4  5  6
bucket:  []  [3]  [2]  [1]  [] [] []
```

Sweep from `freq = 6` down, collecting until 2 values: `freq=3 → [1]`, `freq=2 → [2]`. Result `[1, 2]` — the two most frequent, no sorting or heap required.

---

## Complexity

| Algorithm | Time | Space | Stable | Constraint |
|-----------|------|-------|--------|------------|
| Counting sort | O(n + k) | O(n + k) | yes (right-to-left placement) | integer keys in [0, k) |
| Counting sort (overwrite) | O(n + k) | O(k) | no | integer keys in [0, k) |
| Bucket sort | O(n + m) avg, O(n²) worst | O(n + m) | depends on inner sort | keys spread over m buckets |
| Radix sort (LSD) | O(d·(n + b)) | O(n + b) | yes | fixed-width keys, base b, d digits |

- **Counting** is linear only while `k = O(n)`. If `k ≫ n` (e.g. sorting eight 32-bit ints by raw value → k = 2³²) the `count` array dominates — that is exactly when you switch to **radix** to keep the per-pass range small (`b = 10` or `256`) at the cost of `d` passes.
- **Bucket** achieves O(n) *average* only under a uniform-spread assumption; adversarial input piling everything into one bucket degrades to the inner sort's worst case (O(n²) with insertion sort). Choosing `m ≈ n` buckets and a good index function is what keeps buckets O(1)-sized.
- **The Ω(n log n) barrier does not apply** to any of these because none of them decide order via comparisons — that is the entire point.

---

## Common pitfalls

1. **Range blow-up.** Counting/radix cost space proportional to the key range (or base). Values up to 10⁹ with counting sort allocate a billion-slot array — use radix, or a hash-map/bucket approach, instead. Always sanity-check `k` against `n`.
2. **Negative numbers.** Plain counting sort indexes by value and can't handle negatives. Offset by `-min` (shift all keys into `[0, max-min]`), or split into negative/non-negative passes. LSD radix as written assumes non-negative keys.
3. **Losing stability (and breaking radix).** LSD radix *requires* a stable per-digit sort — place elements **right-to-left** after prefix-summing. Placing left-to-right, or using an unstable subroutine, silently corrupts the final order.
4. **Off-by-one in the prefix sum.** After prefix-summing, `count[v]` is the *exclusive* end index for value `v`; you must **decrement before writing** (`count[v]--; out[count[v]] = v`). Writing then decrementing overshoots the array.
5. **Bucket sort on skewed data.** If keys aren't roughly uniform, one bucket swallows most elements and you're back to O(n²). Either transform the key to spread it, or don't use bucket sort. For "top-k *frequent*", note frequencies are naturally bounded by `[1, n]`, which is why the frequency-bucket trick is always safe.
6. **Forgetting the max scan in radix.** The number of digit passes `d` depends on `max`; compute it first. Iterating a fixed 10 passes wastes work or, if keys exceed 10 digits, under-sorts.
7. **Assuming counting sort is always faster.** For small `n` or huge `k`, the constant factors and the `O(k)` array make it lose to a good comparison sort. It wins specifically when `k = O(n)` and `n` is large.

---

## Problems in this repo that use it

- [0075 — Sort Colors](/0075_sort_colors/README.md) — three values {0,1,2}; a two-pass counting sort (tally then overwrite) is the textbook non-optimal solution, with the one-pass Dutch-National-Flag partition as the O(1)-space optimal.
- [0164 — Maximum Gap](/0164_maximum_gap/README.md) — the marquee example: achieve O(n) by **bucketing** values by range (pigeonhole) so the max gap must fall *between* buckets, or by **radix** sorting — both dodge the O(n log n) of a comparison sort.
- [0274 — H-Index](/0274_h_index/README.md) — cap citation counts at `n` and **counting-sort** them into `n+1` buckets, then sweep from high to low accumulating papers to find the h-index in O(n).
- [0347 — Top K Frequent Elements](/0347_top_k_frequent_elements/README.md) — count frequencies, then **bucket by frequency** (index `1..n`) and sweep from the top for an O(n) alternative to the heap / quickselect solutions.
- [0192 — Word Frequency](/0192_word_frequency/README.md) — a Bash/shell problem whose sort-by-count step is conceptually the bucket-by-frequency idea (`sort | uniq -c | sort -rn`).

### Related classics to know

- LeetCode #451 — Sort Characters By Frequency (the canonical bucket-by-frequency string problem; same idea as #347, out of this repo's current 1–400 range).
- LeetCode #164 / #274 / #347 above are the standard non-comparison-sort interview set.
