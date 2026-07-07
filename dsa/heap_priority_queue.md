# Heap / Priority Queue

> **Category:** Data Structure · **Core operations:** push O(log n), pop O(log n), peek O(1), heapify O(n)

---

## What it is

A **heap** is a complete binary tree stored in a flat array that maintains the
**heap property**:

- **Min-heap:** every parent ≤ its children → the minimum is always at the root.
- **Max-heap:** every parent ≥ its children → the maximum is always at the root.

A **priority queue** is the abstract interface ("give me the highest-priority
item next"); a binary heap is its standard implementation.

Because the tree is *complete* (filled left-to-right, level by level), it needs
no pointers — for the node at index `i` (0-based):

```
parent(i) = (i - 1) / 2
left(i)   = 2*i + 1
right(i)  = 2*i + 2
```

The heap is **partially ordered**, not sorted. Only the root is guaranteed to be
the min/max; siblings have no defined order. That partial order is exactly what
makes push/pop cheap: you only fix one root-to-leaf path (`O(log n)`), never the
whole structure.

### Operation costs

| Operation | Cost | How |
|---|---|---|
| `Peek` (read min/max) | O(1) | it's `a[0]` |
| `Push` | O(log n) | append at end, **sift up** while smaller than parent |
| `Pop` | O(log n) | swap root with last, shrink, **sift down** the new root |
| `heap.Init` (heapify) | O(n) | sift down from the last internal node backwards |
| Update/remove at index `i` | O(log n) | `heap.Fix(h, i)` / `heap.Remove(h, i)` |
| Search for arbitrary value | O(n) | heaps don't support search — pair with a hash map if needed |

---

## How to recognise a heap problem — signals in the statement

Reach for a heap when the problem needs **repeated access to the current
min/max of a changing collection**, but never needs full sorted order.

Strong textual signals:

- **"k-th largest / k-th smallest / top k / k closest / k most frequent"**
  → size-k heap of the *opposite* polarity (min-heap for top-k largest).
- **"merge k sorted lists/arrays/streams"** → min-heap of the k current heads
  (k-way merge). See [0023](../0023_merge_k_sorted_lists/README.md).
- **"median of a data stream" / "running median"** → two heaps: max-heap for
  the lower half, min-heap for the upper half (LeetCode #295).
- **"schedule tasks / meeting rooms / CPU intervals / minimum platforms"**
  → sort by start, min-heap keyed on end time (LeetCode #253, #621, #1834).
- **"repeatedly take the two smallest/largest and combine"** → greedy +
  heap (Huffman coding, LeetCode #1046 Last Stone Weight, #1167 sticks).
- **"shortest path with weighted edges"** → Dijkstra = BFS with a min-heap
  keyed on distance (LeetCode #743, #787). Same idea for Prim's MST and
  A* search.
- **"maximise/minimise while greedily picking the best available option each
  step"** → heap is the "give me the best candidate" engine (LeetCode #502
  IPO, #871 Refueling Stops).
- **"stream / online / data arrives one element at a time"** — you can't sort
  what you haven't seen; a heap maintains order incrementally.

Rules of thumb:

- Need only the min **or** max repeatedly, with interleaved inserts → **heap**.
- Need the k smallest of n where k ≪ n → **size-k heap**, `O(n log k)` beats
  full sort `O(n log n)`.
- Need everything sorted once, no further inserts → just **sort**.
- Need min of a **sliding window** with FIFO eviction → a **monotonic deque**
  is O(n) and usually beats a heap (heaps evict the min, not the oldest).
- Need arbitrary lookup/predecessor queries too → balanced BST / ordered map,
  not a heap.

---

## Go templates

Go has no built-in priority queue type; you implement `heap.Interface` from
`container/heap` (five methods) and the package supplies the sift logic.

### Template 1 — min-heap of ints

```go
package main

import (
    "container/heap"
    "fmt"
)

// IntHeap is a min-heap of ints. For a MAX-heap, flip Less to h[i] > h[j].
type IntHeap []int

func (h IntHeap) Len() int           { return len(h) }            // number of elements
func (h IntHeap) Less(i, j int) bool { return h[i] < h[j] }       // "<" ⇒ min-heap
func (h IntHeap) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }  // used by sift up/down

// Push/Pop take pointer receivers because they resize the slice.
// NOTE: these are called BY heap.Push/heap.Pop — never call them directly.
func (h *IntHeap) Push(x any) { *h = append(*h, x.(int)) }         // append at the end
func (h *IntHeap) Pop() any {                                      // remove the LAST element
    old := *h
    n := len(old)
    x := old[n-1]     // heap.Pop has already swapped the root here
    *h = old[:n-1]
    return x
}

func main() {
    h := &IntHeap{5, 2, 8}
    heap.Init(h)                 // heapify in O(n) — required if starting non-empty
    heap.Push(h, 1)              // O(log n): append + sift up
    fmt.Println((*h)[0])         // peek min in O(1) → 1
    fmt.Println(heap.Pop(h))     // O(log n): pop min → 1
}
```

### Template 2 — heap of structs (priority queue with payload)

```go
// Item carries a payload plus the key we order by.
type Item struct {
    value    string // arbitrary payload
    priority int    // ordering key
}

type PQ []Item

func (pq PQ) Len() int            { return len(pq) }
func (pq PQ) Less(i, j int) bool  { return pq[i].priority < pq[j].priority } // min by priority
func (pq PQ) Swap(i, j int)       { pq[i], pq[j] = pq[j], pq[i] }
func (pq *PQ) Push(x any)         { *pq = append(*pq, x.(Item)) }
func (pq *PQ) Pop() any {
    old := *pq
    n := len(old)
    item := old[n-1]
    *pq = old[:n-1]
    return item
}

// Usage:
//   pq := &PQ{}
//   heap.Push(pq, Item{"task", 3})
//   top := heap.Pop(pq).(Item)   // type-assert the any back to Item
```

### Template 3 — top-k pattern (k largest elements)

```go
// kLargest returns the k largest values of nums using a size-k MIN-heap.
//
// Invariant: the heap always holds the k largest values seen so far,
// with the SMALLEST of those k at the root — the "bouncer at the door".
// Any newcomer must beat the root to get in.
//
// Time:  O(n log k)   Space: O(k)
func kLargest(nums []int, k int) []int {
    h := &IntHeap{}              // min-heap from Template 1
    heap.Init(h)
    for _, x := range nums {
        if h.Len() < k {
            heap.Push(h, x)      // heap not full yet: everyone gets in
        } else if x > (*h)[0] {
            (*h)[0] = x          // newcomer beats the weakest of the top k:
            heap.Fix(h, 0)       // replace root in place and sift down (O(log k))
        }                        // else: newcomer too small, discard
    }
    return *h                    // k largest, in arbitrary (heap) order
}
```

### Template 4 — k-way merge (commented pseudocode)

```text
build min-heap H with the first element of each of the k sources   // O(k)
while H is not empty:                                              // N total pops
    e ← pop(H)                       // global minimum among all current heads
    append e to output
    if e's source has a next element:
        push(H, next element)        // heap size never exceeds k
# Total: O(N log k) time, O(k) extra space.
```

Full Go implementation with `*ListNode`:
[`0023_merge_k_sorted_lists/main.go`](../0023_merge_k_sorted_lists/main.go).

### Manual heap (when asked to implement from scratch)

```go
// siftUp: after appending at index i, bubble it up while it beats its parent.
func siftUp(a []int, i int) {
    for i > 0 {
        p := (i - 1) / 2          // parent index
        if a[i] >= a[p] {         // heap property satisfied — stop
            break
        }
        a[i], a[p] = a[p], a[i]   // violated — swap with parent
        i = p                     // continue from the parent's position
    }
}

// siftDown: after placing a value at index i, push it down to its place.
func siftDown(a []int, i, n int) {
    for {
        l, r, smallest := 2*i+1, 2*i+2, i
        if l < n && a[l] < a[smallest] { smallest = l } // left child smaller?
        if r < n && a[r] < a[smallest] { smallest = r } // right child smaller?
        if smallest == i {                              // neither — done
            break
        }
        a[i], a[smallest] = a[smallest], a[i]           // swap with smaller child
        i = smallest
    }
}

// heapify: build a valid min-heap in O(n) — sift down every internal node,
// last-to-first. Leaves (indices ≥ n/2) are already 1-element heaps.
func heapify(a []int) {
    for i := len(a)/2 - 1; i >= 0; i-- {
        siftDown(a, i, len(a))
    }
}
```

---

## Worked example — LeetCode #23, Merge k Sorted Lists

Merge `lists = [[1,4,5], [1,3,4], [2,6]]` with a min-heap of the current heads
(Template 4). `N = 8` nodes, `k = 3` lists — heap never exceeds 3 elements.

| Step | Pop (list) | Output so far | Push next | Heap after |
|---|---|---|---|---|
| init | — | — | heads 1(L1), 1(L2), 2(L3) | {1(L1), 1(L2), 2(L3)} |
| 1 | 1 (L1) | 1 | 4 (L1) | {1(L2), 2(L3), 4(L1)} |
| 2 | 1 (L2) | 1→1 | 3 (L2) | {2(L3), 3(L2), 4(L1)} |
| 3 | 2 (L3) | 1→1→2 | 6 (L3) | {3(L2), 4(L1), 6(L3)} |
| 4 | 3 (L2) | 1→1→2→3 | 4 (L2) | {4(L1), 4(L2), 6(L3)} |
| 5 | 4 (L1) | 1→1→2→3→4 | 5 (L1) | {4(L2), 5(L1), 6(L3)} |
| 6 | 4 (L2) | …→4→4 | — (L2 done) | {5(L1), 6(L3)} |
| 7 | 5 (L1) | …→5 | — (L1 done) | {6(L3)} |
| 8 | 6 (L3) | …→6 | — (L3 done) | {} |

Result: `1→1→2→3→4→4→5→6`. Each of the 8 nodes was pushed and popped exactly
once at `O(log 3)` each → `O(N log k)`, versus `O(N log N)` for
concatenate-and-sort and `O(N·k)` for repeatedly scanning all k heads.

**Tracing step 1 inside the heap** (array form, min-heap):

1. Heap array `[1(L1), 1(L2), 2(L3)]`. Pop: swap root with last →
   `[2(L3), 1(L2) | 1(L1)]`, detach `1(L1)`.
2. Sift down `2(L3)` from index 0: child `1(L2)` at index 1 is smaller → swap →
   `[1(L2), 2(L3)]`. Property restored.
3. Push `4(L1)`: append → `[1(L2), 2(L3), 4(L1)]`; sift up: parent of index 2
   is index 0 holding `1 ≤ 4` → stays put.

---

## Common pitfalls (and how to avoid them)

1. **Calling `h.Push(x)` / `h.Pop()` instead of `heap.Push(h, x)` /
   `heap.Pop(h)`.** Your methods only append/truncate the slice; the *package*
   functions do the sifting. Direct calls silently corrupt the heap order.
   Rule: your five methods exist for the package, not for you.
2. **Forgetting `heap.Init(h)` on a non-empty slice.** A slice built with
   `append` is not a heap until heapified; `heap.Pop` will then return garbage.
   (Starting empty and only using `heap.Push` is fine.)
3. **Wrong polarity.** Min vs max is decided solely by `Less`. For top-k
   *largest* you want a *min*-heap (evict the smallest of the keepers), and
   vice versa — getting this backwards is the classic top-k bug.
4. **Assuming the backing slice is sorted.** Only `a[0]` is the extreme.
   Iterating the slice, or popping your Pop method directly, does not yield
   sorted order — pop via `heap.Pop` n times if you need it sorted.
5. **Mutating an element's priority in place without `heap.Fix`.** Changing
   `pq[i].priority` breaks the invariant; call `heap.Fix(pq, i)` (or
   `heap.Remove` + `heap.Push`). For Dijkstra, the standard trick is **lazy
   deletion**: push a duplicate with the better distance and skip stale
   entries when popped (check against the dist table).
6. **Pointer vs value receivers.** `Push`/`Pop` must have pointer receivers
   (they resize the slice); pass `&h` to the heap functions. Passing the value
   compiles in some arrangements but loses the resize.
7. **Type-assertion slips.** `heap.Pop` returns `any`; assert to the exact
   stored type (`.(Item)` vs `.(*Item)` mismatch panics at runtime).
8. **Comparator overflow.** `return a-b < 0`-style comparisons overflow for
   extreme ints; compare directly with `<`.
9. **Ties without a tiebreak.** If `Less` treats equal priorities as
   unordered, equal-priority pop order is arbitrary (heaps are not stable).
   When the problem demands FIFO among ties, add an insertion sequence number
   to the ordering key.
10. **Using a heap where a monotonic deque / quickselect fits better.**
    Sliding-window max is O(n) with a deque; a single k-th order statistic is
    O(n) average with quickselect. A heap is for *repeated, interleaved*
    min/max access.

---

## Problems in this repo

*(Problems 0131+ are being added concurrently; this list covers what exists now
and will be extended in a later pass.)*

- [0023 — Merge k Sorted Lists](../0023_merge_k_sorted_lists/README.md) —
  Approach 4 is the canonical **k-way merge with a min-heap** of the k list
  heads: `O(N log k)` time, `O(k)` space, and the merge step of external
  merge sort.

Heap-adjacent notes elsewhere in the repo (heap appears as a follow-up or
alternative, not the implemented solution):

- [0004 — Median of Two Sorted Arrays](../0004_median_of_two_sorted_arrays/README.md)
  — the streaming variant (LeetCode #295) is the classic **two-heaps** median.
- [0042 — Trapping Rain Water](../0042_trapping_rain_water/README.md) — the 3D
  extension (LeetCode #407) requires a min-heap boundary BFS.
- [0056 — Merge Intervals](../0056_merge_intervals/README.md) — the follow-up
  Meeting Rooms II (LeetCode #253) is **sort + min-heap on end times**.
