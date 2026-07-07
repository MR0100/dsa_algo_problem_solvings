# Segment Tree & Fenwick Tree (Binary Indexed Tree)

> **Category:** Range-query data structures
> **Core superpower:** answer *range queries* (sum, min, max, count, …) **and**
> support *point/range updates* on the same array, both in **O(log n)** — where
> a plain array or prefix-sum array forces one of the two operations to O(n).

---

## 1. What the concept is

Both structures solve the same fundamental tension:

| Structure | Point update | Range query | Notes |
|---|---|---|---|
| Plain array | O(1) | O(n) | fast writes, slow reads |
| Prefix-sum array | O(n) | O(1) | fast reads, slow writes |
| **Fenwick tree (BIT)** | **O(log n)** | **O(log n)** | sums / invertible ops only, tiny code |
| **Segment tree** | **O(log n)** | **O(log n)** | any associative op; supports lazy range updates |

### Fenwick tree (Binary Indexed Tree)

A 1-indexed array `tree[1..n]` where `tree[i]` stores the sum of the block of
length `lowbit(i) = i & (-i)` ending at `i`. Prefix sums are assembled by
stripping the lowest set bit repeatedly; updates walk the other direction by
adding it. It only works for operations with an **inverse** (sum, XOR, count)
because a range query `[l, r]` is computed as `prefix(r) − prefix(l−1)`.

### Segment tree

A binary tree (usually stored in a flat array of size `4n`) where each node
covers an interval `[lo, hi]` and stores an aggregate (sum, min, max, gcd, …)
of that interval. The root covers `[0, n−1]`; each internal node splits its
interval in half. A query on `[l, r]` decomposes into at most `O(log n)`
disjoint node intervals. Works for **any associative operation** — no inverse
needed — and extends to **range updates** via lazy propagation.

Rule of thumb: **reach for Fenwick when you need sums/counts with point
updates (it's 10 lines); reach for a segment tree when you need min/max/gcd,
non-invertible merges, or lazy range updates.**

---

## 2. How to recognise a problem needs it — signals in the statement

- **"queries" plural, mixed with "updates"** — e.g. *"Implement `update(i, val)`
  and `sumRange(l, r)`"* (LeetCode #307 Range Sum Query — Mutable). The mix of
  read + write is the giveaway; static arrays would just use prefix sums.
- **"count of smaller/greater elements to the left/right"** — count-inversions
  flavour (LC #315 Count of Smaller Numbers After Self, #493 Reverse Pairs).
  Process elements in one order, query "how many already inserted are < x"
  → Fenwick over the *value domain* (after coordinate compression).
- **"number of inversions" / "important reverse pairs"** — same pattern.
- **Online processing / streaming**: elements arrive one at a time and each
  step asks a question about everything seen so far.
- **"maximum/minimum in a range" with updates** — segment tree (min/max has no
  inverse, so Fenwick's prefix trick doesn't apply directly).
- **"range update, range query"** (add `v` to all of `[l, r]`, then ask sums)
  — segment tree with lazy propagation, or two Fenwick trees.
- **k-th smallest among remaining / order statistics with deletions** — Fenwick
  of 0/1 presence flags + binary search on prefix sums. (LC #60 Permutation
  Sequence can use this to make "remove the k-th unused digit" O(log n).)
- **Skyline / interval max height, calendar booking with counts** (LC #218,
  #732) — segment tree, often with coordinate compression or dynamic nodes.
- Constraints hint: `n` and `q` both up to ~10⁵ and an O(n·q) brute force is
  ~10¹⁰ ops → you need O((n+q) log n).

**Anti-signals:** array is static (prefix sums suffice); only suffix/prefix
questions with no updates; small n (≤ ~1000) where O(n²) passes.

---

## 3. General templates in Go

### 3a. Fenwick tree (point update, prefix/range sum)

```go
// Fenwick (Binary Indexed Tree), 1-indexed internally.
//
// Pseudocode idea:
//   update(i, delta): while i <= n { tree[i] += delta; i += i & (-i) }
//   prefix(i):        while i >= 1 { sum += tree[i];  i -= i & (-i) }
//
// Time:  O(log n) per op  ·  Space: O(n)
type Fenwick struct {
	n    int
	tree []int // tree[i] holds the sum of the block of length i&(-i) ending at i
}

func NewFenwick(n int) *Fenwick {
	return &Fenwick{n: n, tree: make([]int, n+1)} // index 0 unused
}

// Update adds delta at position i (1-indexed).
func (f *Fenwick) Update(i, delta int) {
	for ; i <= f.n; i += i & (-i) { // climb: add lowest set bit
		f.tree[i] += delta
	}
}

// Prefix returns sum of [1..i].
func (f *Fenwick) Prefix(i int) int {
	s := 0
	for ; i > 0; i -= i & (-i) { // descend: strip lowest set bit
		s += f.tree[i]
	}
	return s
}

// RangeSum returns sum of [l..r] (1-indexed, inclusive) using the inverse.
func (f *Fenwick) RangeSum(l, r int) int {
	return f.Prefix(r) - f.Prefix(l-1)
}
```

### 3b. Coordinate compression (companion to Fenwick counting)

```go
// Map arbitrary values to ranks 1..m so the Fenwick is sized by the number
// of DISTINCT values, not the value range (which may be up to ±10^9).
func compress(nums []int) map[int]int {
	sorted := append([]int(nil), nums...)
	sort.Ints(sorted)
	rank := make(map[int]int, len(nums))
	for _, v := range sorted {
		if _, seen := rank[v]; !seen {
			rank[v] = len(rank) + 1 // ranks start at 1 (Fenwick is 1-indexed)
		}
	}
	return rank
}
```

### 3c. Segment tree (point update, range query — sum shown, swap the merge for min/max)

```go
// Segment tree over nums[0..n-1], stored in a flat array of size 4n.
// Node `node` covers [lo, hi]; children are 2*node+1 and 2*node+2.
//
// Time:  build O(n), query/update O(log n)  ·  Space: O(n)
type SegTree struct {
	n    int
	tree []int
}

func NewSegTree(nums []int) *SegTree {
	st := &SegTree{n: len(nums), tree: make([]int, 4*len(nums))}
	st.build(nums, 0, 0, st.n-1)
	return st
}

func (st *SegTree) build(nums []int, node, lo, hi int) {
	if lo == hi { // leaf: covers a single element
		st.tree[node] = nums[lo]
		return
	}
	mid := lo + (hi-lo)/2
	st.build(nums, 2*node+1, lo, mid)    // build left half
	st.build(nums, 2*node+2, mid+1, hi)  // build right half
	st.tree[node] = st.tree[2*node+1] + st.tree[2*node+2] // merge (change for min/max)
}

// Update sets nums[i] = val (point assignment).
func (st *SegTree) Update(i, val int) { st.update(0, 0, st.n-1, i, val) }

func (st *SegTree) update(node, lo, hi, i, val int) {
	if lo == hi { // reached the leaf for index i
		st.tree[node] = val
		return
	}
	mid := lo + (hi-lo)/2
	if i <= mid {
		st.update(2*node+1, lo, mid, i, val)
	} else {
		st.update(2*node+2, mid+1, hi, i, val)
	}
	st.tree[node] = st.tree[2*node+1] + st.tree[2*node+2] // re-merge on the way up
}

// Query returns the aggregate over [l, r] (0-indexed, inclusive).
func (st *SegTree) Query(l, r int) int { return st.query(0, 0, st.n-1, l, r) }

func (st *SegTree) query(node, lo, hi, l, r int) int {
	if r < lo || hi < l { // no overlap → identity element (0 for sum, +inf for min)
		return 0
	}
	if l <= lo && hi <= r { // total overlap → this node's aggregate is usable as-is
		return st.tree[node]
	}
	mid := lo + (hi-lo)/2 // partial overlap → recurse both sides and merge
	return st.query(2*node+1, lo, mid, l, r) + st.query(2*node+2, mid+1, hi, l, r)
}
```

### 3d. Lazy propagation sketch (range update, range query)

```go
// For "add v to every element of [l, r]" keep a parallel lazy[] array.
// Pseudocode:
//   push(node, lo, hi):                 // flush pending update to children
//     if lazy[node] != 0 and lo != hi:
//       for each child c: lazy[c] += lazy[node]; tree[c] += lazy[node] * lenOf(c)
//       lazy[node] = 0
//   rangeAdd(node, lo, hi, l, r, v):
//     no overlap  -> return
//     total       -> tree[node] += v * (hi-lo+1); lazy[node] += v; return
//     partial     -> push(node); recurse children; re-merge
//   query is the same shape, calling push() before recursing.
```

---

## 4. Worked example, traced step by step

**Problem (LC #315 pattern):** for `nums = [5, 2, 6, 1]`, count for each
element how many numbers **after** it are **smaller**. Expected: `[2, 1, 1, 0]`.

**Plan:** compress values → ranks `{1:1, 2:2, 5:3, 6:4}`. Scan **right to
left**; before inserting `nums[i]`, query `Prefix(rank−1)` = how many
already-inserted (i.e. to the right) values are strictly smaller. Then
`Update(rank, +1)`.

Fenwick of size 4. `tree[i]` covers a block of length `i&(-i)`:
`tree[1]`→[1], `tree[2]`→[1..2], `tree[3]`→[3], `tree[4]`→[1..4].

| Step | i | nums[i] | rank | Query `Prefix(rank−1)` | Result | Update path (`+1`) | tree after `[t1,t2,t3,t4]` |
|---|---|---|---|---|---|---|---|
| 1 | 3 | 1 | 1 | `Prefix(0)` = 0 | ans[3]=0 | 1 → 2 → 4 | [1, 1, 0, 1] |
| 2 | 2 | 6 | 4 | `Prefix(3)` = t3 + t2 = 0+1 | ans[2]=1 | 4 | [1, 1, 0, 2] |
| 3 | 1 | 2 | 2 | `Prefix(1)` = t1 = 1 | ans[1]=1 | 2 → 4 | [1, 2, 0, 3] |
| 4 | 0 | 5 | 3 | `Prefix(2)` = t2 = 2 | ans[0]=2 | 3 → 4 | [1, 2, 1, 4] |

Query decomposition detail for step 2: `Prefix(3)` starts at `i=3`
(`3&-3 = 1`, take `tree[3]`, then `i = 3−1 = 2`), then `tree[2]`
(`i = 2−2 = 0`, stop). Update path detail for step 1: insert at `i=1`, then
`i = 1 + (1&-1) = 2`, then `i = 2 + (2&-2) = 4`, then `4+4 = 8 > n`, stop.

Final answer read left to right: **[2, 1, 1, 0]**. Total work:
n·2·O(log n) = O(n log n) versus O(n²) brute force.

---

## 5. Common pitfalls and how to avoid them

1. **0-indexing a Fenwick tree.** `lowbit(0) = 0` → infinite loop. Always keep
   the tree 1-indexed internally; convert at the API boundary.
2. **Using Fenwick for min/max.** `Prefix(r) − Prefix(l−1)` needs an inverse;
   min/max have none. Use a segment tree (or a specialised BIT that only
   supports prefix-min with restricted updates).
3. **Forgetting coordinate compression.** Values up to ±10⁹ do not fit in a
   Fenwick sized by value. Compress to ranks 1..distinct first; remember
   duplicates share a rank.
4. **`update` with a delta vs. an assignment.** Fenwick's natural op is
   `+= delta`. For "set index i to v" you must add `v − current[i]` and keep a
   shadow copy of current values. Segment tree assignment is direct.
5. **Segment tree array too small.** Use `4n`, not `2n` — the flat-array layout
   of an unbalanced recursion can index up to ~4n. (Or build a bottom-up
   iterative tree of exactly `2n`.)
6. **Wrong identity element in query.** No-overlap must return the identity of
   the merge: `0` for sum, `+inf` for min, `−inf` for max, `0` for gcd/xor.
   Returning `0` for a min-tree silently corrupts answers.
7. **Forgetting to re-merge on the way up** after a point update
   (`tree[node] = merge(children)`), or forgetting to `push` lazy values
   *before* recursing in a lazy tree — both give stale aggregates.
8. **Off-by-one between inclusive/exclusive ranges.** Pick a convention
   (inclusive `[l, r]` here) and enforce it in every function signature.
9. **Query direction for "count smaller after self".** Scan right→left and
   query strictly-less (`rank−1`); scanning the wrong way or using `rank`
   counts equal elements — a classic wrong-answer on duplicates.
10. **Overflow.** Sums of 10⁵ elements up to 10⁹ exceed int32; in Go use
    `int` (64-bit on modern platforms) or `int64` explicitly.
11. **Reverse pairs (`nums[i] > 2*nums[j]`) compression trap:** you must
    compress the *union* of `v` and `2v` values (or binary-search ranks per
    query), otherwise `2*nums[j]` has no rank.

---

## 6. Problems in this repo

No solved problem in this repo (0001–0130 at the time of writing) uses a
segment tree or Fenwick tree as a primary approach. One problem mentions it
as an alternative optimisation:

- [0060 — Permutation Sequence](../0060_permutation_sequence/README.md) — the
  "remove the k-th unused digit" step can be done in O(log n) with a Fenwick
  tree of presence flags + binary search on prefix sums (irrelevant for n ≤ 9,
  noted in its Key Takeaways).

Classic LeetCode problems to link here once solved: #307 Range Sum Query —
Mutable, #315 Count of Smaller Numbers After Self, #327 Count of Range Sum,
#493 Reverse Pairs, #218 The Skyline Problem, #732 My Calendar III.

*(Problems 0131–0400 are being written concurrently; a later pass will add
links from that range.)*
