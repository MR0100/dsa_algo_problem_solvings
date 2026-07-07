# K-Way Merge

> **The pattern:** you have **k already-sorted** sequences (lists, rows, generated streams) and want to consume their combined order — the global minimum first, then the next, and so on. The engine is a **min-heap holding one "frontier" element per sequence**: pop the smallest, then push that sequence's successor. This same *sorted-frontier expansion* generalises to "kth-smallest in a sorted structure" and to generating sequences in increasing order (ugly numbers, k smallest pairs).

---

## What it is

Two-way merge (the merge step of mergesort) walks two sorted lists with two pointers, always taking the smaller head. **K-way merge** does the same across `k` lists at once — but scanning all `k` heads on every step would cost `O(k)` per output element. A **min-heap of the k current heads** reduces the "who is smallest now?" question to `O(log k)`.

The mental model is a **frontier**: at any moment you hold exactly one candidate from each still-active sequence — the smallest unconsumed element of that sequence. The global next element is the minimum of the frontier. Pop it, and *advance only the sequence it came from* (push that sequence's next element into the heap). The frontier stays size ≤ k, and repeating `n` times streams out the fully merged order.

This generalises past "merge k lists" into a whole family:

- **Sorted-matrix / staircase expansion** — a matrix with sorted rows *and* columns is `k` sorted rows; the frontier walks a monotone "staircase" from the top-left. Finding the kth-smallest is `k` heap pops.
- **k smallest pairs / sums** — the "sequences" are implicit: pair `(a_i, b_j)` has neighbours `(a_{i+1}, b_j)` and `(a_i, b_{j+1})`. Seed the heap with the smallest, expand its successors — a k-way merge over a 2-D sorted grid.
- **Merged-multiple-generator** (ugly numbers, super ugly numbers) — generate the sorted sequence `{ prev × p : p in primes }` by merging k monotone streams, one per multiplier.

So k-way merge is really *"repeatedly extract the minimum of a set of monotone frontiers, then push each popped element's successors."* Recognising that unifies half a dozen "hard" heap problems.

---

## When to recognise it

| Signal in the problem | Why k-way merge fits |
|-----------------------|----------------------|
| "**merge k sorted** lists / arrays" | The definition; heap of k heads |
| "**kth smallest** in a sorted matrix / among sorted lists" | Do k min-extractions from the frontier |
| "**k pairs / k sums** with smallest total" from sorted arrays | Implicit sorted grid; expand successors of the popped pair |
| Generate the "**nth ugly / super-ugly**" number | Merge k monotone multiplier-streams in increasing order |
| "combine many sorted streams, keep global order" | Streaming merge; heap holds one frontier per stream |
| Rows **and** columns sorted (Young tableau / staircase) | Sorted frontier moves monotonically; heap or binary-search-on-value |
| "smallest / next in increasing order, built from previous ones" | Frontier expansion — push successors, pop the min |

**Complementary technique — binary search on the answer.** For *kth-smallest* questions where the value range is known (sorted matrix, k-th smallest pair distance), you can **binary-search the value** and count "how many elements ≤ mid" instead of popping k times from a heap. That trades the heap's `O(k log …)` for `O((n or rows)·log(range))` and O(1) space. Know both; pick by whether `k` or the value range is the smaller lever.

**When *not* to use it:** the inputs aren't sorted (sort or use a different structure first); there are only two sequences (a plain two-pointer merge is simpler and heap-free); or you need *all* pairs/products, not the smallest few (a heap gives you order but not a speedup over generating everything).

---

## General template / pseudocode

Go's heap lives in `container/heap`; you implement the five-method `heap.Interface`. A reusable min-heap of "items" (value + which sequence + position) is the backbone.

### Merge k sorted linked lists

```go
import "container/heap"

// ListNode is the standard singly-linked list node.
type ListNode struct {
    Val  int
    Next *ListNode
}

// nodeHeap is a min-heap of list heads, ordered by Val.
type nodeHeap []*ListNode

func (h nodeHeap) Len() int            { return len(h) }
func (h nodeHeap) Less(i, j int) bool  { return h[i].Val < h[j].Val } // MIN-heap
func (h nodeHeap) Swap(i, j int)       { h[i], h[j] = h[j], h[i] }
func (h *nodeHeap) Push(x any)         { *h = append(*h, x.(*ListNode)) }
func (h *nodeHeap) Pop() any {
    old := *h
    n := len(old)
    item := old[n-1]
    *h = old[:n-1]
    return item
}

// mergeKLists merges k sorted lists into one sorted list.
// Time: O(N log k) for N total nodes   Space: O(k) for the heap
func mergeKLists(lists []*ListNode) *ListNode {
    h := &nodeHeap{}
    heap.Init(h)
    for _, node := range lists {
        if node != nil {
            heap.Push(h, node) // seed the frontier: one head per list
        }
    }
    dummy := &ListNode{}
    tail := dummy
    for h.Len() > 0 {
        node := heap.Pop(h).(*ListNode) // global minimum among frontiers
        tail.Next = node
        tail = tail.Next
        if node.Next != nil {
            heap.Push(h, node.Next) // advance only the list we consumed
        }
    }
    return dummy.Next
}
```

### Kth smallest in a row+column sorted matrix (heap / staircase)

```go
type cell struct{ val, r, c int }
type cellHeap []cell

func (h cellHeap) Len() int           { return len(h) }
func (h cellHeap) Less(i, j int) bool { return h[i].val < h[j].val }
func (h cellHeap) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }
func (h *cellHeap) Push(x any)        { *h = append(*h, x.(cell)) }
func (h *cellHeap) Pop() any {
    old := *h
    n := len(old)
    it := old[n-1]
    *h = old[:n-1]
    return it
}

// kthSmallest returns the k-th smallest value (1-indexed) in an n×n matrix
// whose rows AND columns are sorted ascending. Treat each row as a sorted
// list; the heap holds the current frontier cell of each active row.
// Time: O(k log k)   Space: O(k)
func kthSmallest(matrix [][]int, k int) int {
    n := len(matrix)
    h := &cellHeap{}
    // Seed with the first column: the smallest element of each row.
    for r := 0; r < n && r < k; r++ {
        *h = append(*h, cell{matrix[r][0], r, 0})
    }
    heap.Init(h)
    var cur cell
    for i := 0; i < k; i++ { // pop k times → the k-th is the answer
        cur = heap.Pop(h).(cell)
        if cur.c+1 < n {
            // advance rightward within the same row (its next-larger element)
            heap.Push(h, cell{matrix[cur.r][cur.c+1], cur.r, cur.c + 1})
        }
    }
    return cur.val
}
```

### Frontier expansion for generated sequences (Ugly Number II, multi-pointer)

Not every k-way merge needs a heap. When there are only a *few* streams and each is a monotone multiple of the growing result, `k` integer pointers beat a heap:

```go
// nthUglyNumber returns the n-th positive integer whose only prime factors
// are 2, 3, and 5. It merges three monotone streams (×2, ×3, ×5) over the
// ugly numbers already built, using one pointer per stream.
// Time: O(n)   Space: O(n)
func nthUglyNumber(n int) int {
    ugly := make([]int, n)
    ugly[0] = 1
    i2, i3, i5 := 0, 0, 0 // pointers into `ugly` for each multiplier
    for i := 1; i < n; i++ {
        n2, n3, n5 := ugly[i2]*2, ugly[i3]*3, ugly[i5]*5
        next := min3(n2, n3, n5) // smallest frontier value
        ugly[i] = next
        // Advance EVERY pointer whose candidate equals next (dedupes, e.g. 6=2×3).
        if next == n2 {
            i2++
        }
        if next == n3 {
            i3++
        }
        if next == n5 {
            i5++
        }
    }
    return ugly[n-1]
}

func min3(a, b, c int) int {
    m := a
    if b < m {
        m = b
    }
    if c < m {
        m = c
    }
    return m
}
```

---

## Worked example — step-by-step trace

### Merge k lists: `[[1,4,5], [1,3,4], [2,6]]`

Seed the heap with the three heads. Heap shown as a multiset (min at left):

| step | heap (frontier) | pop | push successor | output so far |
|------|-----------------|-----|----------------|---------------|
| 0 | {1ₐ, 1_b, 2} | — | — | |
| 1 | {1_b, 2, 4ₐ} | 1ₐ | 4 (from list A) | 1 |
| 2 | {2, 3_b, 4ₐ} | 1_b | 3 (from list B) | 1,1 |
| 3 | {3_b, 4ₐ, 6} | 2 | 6 (from list C) | 1,1,2 |
| 4 | {4ₐ, 4_b, 6} | 3_b | 4 (from list B) | 1,1,2,3 |
| 5 | {4_b, 5ₐ, 6} | 4ₐ | 5 (from list A) | 1,1,2,3,4 |
| 6 | {5ₐ, 6} | 4_b | (B exhausted) | 1,1,2,3,4,4 |
| 7 | {6} | 5ₐ | (A exhausted) | …,5 |
| 8 | {} | 6 | (C exhausted) | …,6 |

Merged: `1,1,2,3,4,4,5,6`. The heap never exceeded size 3 = k; total work `O(N log k)`.

### Kth smallest in a sorted matrix, k = 5

```
matrix = [ 1  5  9 ]
         [10 11 13 ]
         [12 13 15 ]
```

Seed the first column (one frontier per row): heap = {(1,r0,c0), (10,r1,c0), (12,r2,c0)}.

| pop # | pop (val,r,c) | push (matrix[r][c+1]) | heap after |
|-------|----------------|-----------------------|------------|
| 1 | (1, 0,0) | (5, 0,1) | {5, 10, 12} |
| 2 | (5, 0,1) | (9, 0,2) | {9, 10, 12} |
| 3 | (9, 0,2) | row 0 done | {10, 12} |
| 4 | (10, 1,0) | (11, 1,1) | {11, 12} |
| 5 | (11, 1,1) | (13, 1,2) | {12, 13} |

The 5th pop is **11** → answer. Only the *frontier* was ever in the heap, and each pop advanced a single row rightward along the staircase.

---

## Complexity

| Problem shape | Time | Space | Note |
|---------------|------|-------|------|
| Merge k sorted lists/arrays (N total elements) | O(N log k) | O(k) | heap holds ≤ k frontiers |
| Kth smallest in n×n sorted matrix (heap) | O(k log k), capped O(k log n) | O(min(k, n)) | seed ≤ min(k,n) rows |
| Kth smallest in sorted matrix (binary search on value) | O(n · log(max−min)) | O(1) | count ≤ mid per candidate |
| K smallest pairs / sums | O(k log k) | O(k) | expand popped pair's neighbours |
| Nth ugly number (k fixed multipliers) | O(n·k) | O(n) | k pointers, no heap; k tiny |

- **Why `log k`, not `log N`:** the heap only ever contains one live element per sequence, so its size is bounded by `k` (or by how many frontiers you seed). Each of the `N` elements is pushed and popped once → `N` heap operations at `O(log k)` each.
- **Heap vs. binary-search-on-value for kth-smallest:** the heap streams the first `k` in order and is intuitive; binary search on the value avoids the heap entirely and uses O(1) space, at the cost of a per-guess counting pass. Choose by which of `k` and the value range is smaller.
- **Multi-pointer beats the heap when k is a small constant** (ugly numbers: k = 3): three comparisons per step is cheaper than heap bookkeeping, giving a clean `O(n·k)`.

---

## Common pitfalls

1. **Max-heap instead of min-heap.** K-way merge needs the *smallest* frontier each step. In Go, `Less` must be `<`. Reversing it silently produces descending garbage.
2. **Advancing the wrong sequence (or all of them).** After popping the minimum, push **only** the successor of *that* element's sequence. Pushing every sequence's next, or re-pushing the whole frontier, breaks correctness and blows up the heap.
3. **Pushing duplicate frontier cells (pairs/matrix).** When two neighbours can reach the same cell (e.g. `(i,j)` from both `(i-1,j)` and `(i,j-1)`), guard with a `visited` set (a `map[[2]int]bool` or a boolean grid), or you'll process the same element multiple times and miscount k.
4. **Forgetting to skip nil / empty inputs when seeding.** A nil list head or an empty row must not enter the heap; check before the initial push, and re-check `node.Next != nil` before every subsequent push.
5. **Not deduplicating in generated sequences.** Ugly Number II must advance *every* pointer whose candidate equals the chosen minimum (6 = 2×3 = 3×2). Advancing only one leaves 6 in two streams and emits it twice. Use independent `if`s, not `else if`.
6. **Seeding more than necessary for kth-smallest.** You never need more than `min(k, rows)` frontiers — seeding all n rows when k is small wastes space (and, for k < n, is simply unnecessary).
7. **Integer overflow in binary-search-on-value / sums.** `mid = lo + (hi-lo)/2` avoids overflow; summing large pair values may need care. Also ensure the value-count is monotone in `mid` for the binary search to be valid.
8. **Confusing "k smallest" with "kth smallest".** Some problems want the single kth element (stop after k pops); others want the whole list of k (collect each pop). Read which, and size the loop accordingly.

---

## Problems in this repo that use it

- [0023 — Merge k Sorted Lists](/0023_merge_k_sorted_lists/README.md) — the canonical k-way merge: a min-heap of the k list heads (also solvable by divide-and-conquer pairwise merging in the same O(N log k)).
- [0373 — Find K Pairs with Smallest Sums](/0373_find_k_pairs_with_smallest_sums/README.md) — implicit sorted 2-D grid; seed the heap with the smallest pairs and expand each popped pair's two neighbours, guarding against duplicates.
- [0378 — Kth Smallest Element in a Sorted Matrix](/0378_kth_smallest_element_in_a_sorted_matrix/README.md) — treat rows as sorted lists and pop k times from the frontier heap, or binary-search the value and count elements ≤ mid.
- [0264 — Ugly Number II](/0264_ugly_number_ii/README.md) — merge three monotone multiplier streams (×2, ×3, ×5) with three pointers (a heap-free k-way merge), advancing every pointer that ties for the minimum.
- [0313 — Super Ugly Number](/0313_super_ugly_number/README.md) — the same frontier-expansion generalised to an arbitrary set of `k` primes: one pointer (or one heap entry) per prime.

### Related classics to know

- LeetCode #21 — Merge Two Sorted Lists (the k = 2 base case; plain two-pointer merge, no heap).
- LeetCode #23 / #373 / #378 above are the standard k-way-merge interview set.
- LeetCode #295 — Find Median from a Data Stream (two heaps; a cousin of the frontier idea, maintaining order incrementally).
