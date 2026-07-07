# Quickselect

> **Category:** Selection algorithm · Divide and Conquer · Partitioning
> **Prerequisite reading:** [`/dsa/sorting.md`](/dsa/sorting.md) (quicksort partitioning), [`/dsa/heap_priority_queue.md`](/dsa/heap_priority_queue.md) (the main alternative)

---

## What it is

Quickselect (Hoare's selection algorithm, 1961 — same author as quicksort) finds
the **k-th smallest (or largest) element of an unsorted array in average O(n) time
without fully sorting the array**.

It reuses quicksort's core insight: after one **partition** pass around a pivot,
the pivot lands at its *final sorted position* `p`, with everything smaller on its
left and everything larger on its right — even though neither side is itself sorted.

- If `p == k` → the pivot **is** the answer. Done.
- If `p > k` → the answer lives entirely in the **left** part.
- If `p < k` → the answer lives entirely in the **right** part.

The crucial difference from quicksort: quicksort must recurse into **both** halves
(→ O(n log n)), while quickselect recurses into **only one** half. The work forms a
geometric series `n + n/2 + n/4 + … ≈ 2n`, giving **average O(n)**.

| | Quicksort | Quickselect |
|---|---|---|
| Goal | full sorted order | one order statistic |
| Recursion | both sides | one side only |
| Average time | O(n log n) | **O(n)** |
| Worst time | O(n²) | O(n²) (O(n) guaranteed with median-of-medians or shuffling in expectation) |
| Space | O(log n) stack | O(1) iterative / O(log n) recursive |

---

## How to recognise it — signals in the problem statement

Reach for quickselect when you see:

1. **"k-th largest / k-th smallest"** — the canonical trigger.
   *"Find the kth largest element in an unsorted array."*
2. **"top k" / "k closest" / "k most frequent"** where the **output order among the
   k doesn't matter**. Partitioning puts the best k into the first k slots without
   sorting them internally.
3. **"median"** of an unsorted array — the median is just the (n/2)-th order
   statistic. (E.g. median used as a pivot value in wiggle sort / minimize-moves
   problems.)
4. **A follow-up asking to beat O(n log n)** on a selection-flavoured problem —
   the interviewer is fishing for average-O(n) quickselect after you offer
   sort (O(n log n)) or heap (O(n log k)).
5. **All data fits in memory and can be mutated.** Quickselect reorders the array
   in place. If the input is a *stream*, or must not be modified, a heap is the
   right tool instead (see pitfalls).

Rule of thumb for interviews: mention **three** ladder rungs —
sort `O(n log n)` → heap `O(n log k)` → quickselect `O(n)` average — then implement
the one the interviewer steers you to.

---

## General template (Go)

### Template 1 — Lomuto partition, iterative (recommended default)

The Lomuto scheme is the easiest to write correctly under pressure: the pivot ends
at an exactly known index, so the comparison with `k` is direct.

```go
// quickSelect returns the element that would sit at index k (0-based)
// if nums were sorted ascending — i.e. the (k+1)-th smallest element.
// Average O(n) time, O(1) space. Mutates nums.
func quickSelect(nums []int, k int) int {
	lo, hi := 0, len(nums)-1 // current search window, inclusive both ends

	for {
		// Randomised pivot defeats the O(n²) adversarial/sorted-input case:
		// expected O(n) regardless of input order.
		p := lo + rand.Intn(hi-lo+1)
		// Move the chosen pivot to the end so Lomuto can scan [lo, hi-1].
		nums[p], nums[hi] = nums[hi], nums[p]

		// ── Lomuto partition ────────────────────────────────────────────
		// Invariant: nums[lo..store-1] < pivot,  nums[store..i-1] >= pivot.
		pivot := nums[hi]
		store := lo // next slot for an element smaller than the pivot
		for i := lo; i < hi; i++ {
			if nums[i] < pivot {
				nums[store], nums[i] = nums[i], nums[store]
				store++
			}
		}
		// Drop the pivot into its final sorted position.
		nums[store], nums[hi] = nums[hi], nums[store]
		// ────────────────────────────────────────────────────────────────

		switch {
		case store == k: // pivot IS the k-th smallest — done
			return nums[store]
		case store < k: // answer is strictly to the right
			lo = store + 1
		default: // store > k: answer is strictly to the left
			hi = store - 1
		}
	}
}
```

Commented pseudocode of the same idea:

```text
quickselect(A, k):
    lo, hi = 0, len(A)-1
    loop:
        p = partition(A, lo, hi)   # pivot lands at final index p
        if p == k: return A[p]     # found the order statistic
        if p <  k: lo = p + 1      # discard left half INCLUDING pivot
        else:      hi = p - 1      # discard right half INCLUDING pivot
```

### Adapting the template

- **k-th LARGEST** (LeetCode 215 phrasing): the k-th largest is the element at
  sorted index `n - k`. Call `quickSelect(nums, len(nums)-k)`.
  Alternatively flip the comparison to `>` and select index `k-1` directly.
- **Top-k / k-closest**: run quickselect for index `k-1`; afterwards
  `nums[0..k-1]` holds the k best (unordered). Return that slice.
- **Custom keys** (distance to origin, frequency): compare by `key(nums[i]) < key(pivot)`
  or quickselect over an index/pair slice.

### Template 2 — Hoare partition (fewer swaps, trickier bounds)

Hoare's scheme does ~3× fewer swaps but the returned index `j` is **not**
necessarily the pivot's final position — it only guarantees
`A[lo..j] <= A[j+1..hi]` element-wise between the two blocks. So the recursion
condition changes: recurse on `[lo, j]` when `k <= j`, else `[j+1, hi]`, and you
can never "return early on p == k" — the loop terminates when `lo == hi`.

```go
func quickSelectHoare(nums []int, k int) int {
	lo, hi := 0, len(nums)-1
	for lo < hi {
		pivot := nums[lo+rand.Intn(hi-lo+1)] // pivot by VALUE, not index
		i, j := lo-1, hi+1
		for {
			for i++; nums[i] < pivot; i++ { // advance to first elem >= pivot
			}
			for j--; nums[j] > pivot; j-- { // retreat to first elem <= pivot
			}
			if i >= j { // pointers crossed: [lo..j] <= [j+1..hi]
				break
			}
			nums[i], nums[j] = nums[j], nums[i]
		}
		if k <= j { // k-th smallest lies in the left block
			hi = j
		} else { // ... or the right block
			lo = j + 1
		}
	}
	return nums[k] // window shrunk to one element: the answer
}
```

Use Hoare when the input has **many duplicates** — Lomuto degrades to O(n²) on an
all-equal array (see pitfalls), Hoare handles it gracefully. Or use a
**three-way (Dutch national flag) partition** for a Lomuto-style scheme that is
also duplicate-proof.

### Worst-case O(n) guarantee — median of medians (know, don't code)

Picking the pivot with the *median-of-medians* rule (median of the medians of
groups of 5) guarantees the pivot discards ≥30% of elements each round →
**worst-case O(n)**, at the price of a large constant factor and fiddly code.
In interviews: *name it*, state the bound, and say that in practice a random
pivot's expected O(n) is what everyone ships. `C++`'s `nth_element` and Go's
lack of a stdlib equivalent are good trivia to mention.

---

## Worked example — step-by-step trace

Problem: **k-th largest** in `nums = [3, 2, 1, 5, 6, 4]`, `k = 2` (LeetCode 215,
Example 1; expected answer **5**).

Convert: k-th largest = sorted index `n - k = 6 - 2 = 4`. Run
`quickSelect(nums, 4)` with the Lomuto template. Assume the "random" pivot pick
happens to choose the last element each round (deterministic for the trace).

**Round 1** — window `[lo=0, hi=5]`, array `[3, 2, 1, 5, 6, 4]`, pivot `nums[5] = 4`.

| i | nums[i] | < 4 ? | action | array after | store |
|---|---------|-------|--------|-------------|-------|
| — | — | — | init | `[3, 2, 1, 5, 6, 4]` | 0 |
| 0 | 3 | yes | swap(0,0), store→1 | `[3, 2, 1, 5, 6, 4]` | 1 |
| 1 | 2 | yes | swap(1,1), store→2 | `[3, 2, 1, 5, 6, 4]` | 2 |
| 2 | 1 | yes | swap(2,2), store→3 | `[3, 2, 1, 5, 6, 4]` | 3 |
| 3 | 5 | no | skip | `[3, 2, 1, 5, 6, 4]` | 3 |
| 4 | 6 | no | skip | `[3, 2, 1, 5, 6, 4]` | 3 |
| end | — | — | swap pivot into slot 3 | `[3, 2, 1, 4, 6, 5]` | 3 |

Pivot's final index `p = 3`. We want `k = 4`. Since `3 < 4`, discard the left side
**and** the pivot: new window `[lo=4, hi=5]`.

**Round 2** — window `[4, 5]`, array `[3, 2, 1, 4, 6, 5]`, pivot `nums[5] = 5`.

| i | nums[i] | < 5 ? | action | array after | store |
|---|---------|-------|--------|-------------|-------|
| — | — | — | init | `[3, 2, 1, 4, 6, 5]` | 4 |
| 4 | 6 | no | skip | `[3, 2, 1, 4, 6, 5]` | 4 |
| end | — | — | swap pivot into slot 4 | `[3, 2, 1, 4, 5, 6]` | 4 |

Pivot's final index `p = 4 == k`. **Return `nums[4] = 5`.** ✓

Notice the array is *not* fully sorted at any point we relied on — only the
pivot positions 3 and 4 were guaranteed. Total elements scanned: 6 + 2 = 8 ≈ n,
illustrating the geometric-series O(n) behaviour.

---

## Complexity summary

- **Time:** average/expected **O(n)** (random pivot); worst **O(n²)**
  (unlucky pivots — impossible to force from outside once the pivot is random);
  worst **O(n)** with median-of-medians.
- **Space:** **O(1)** iterative (both templates above); O(log n) expected stack
  if written recursively.
- **Recurrence intuition:** `T(n) = T(n/2) + O(n)` → Master theorem case 3 → O(n).

---

## Common pitfalls and how to avoid them

1. **Off-by-one between "k-th largest" and array index.**
   k-th largest ⇔ sorted index `n - k`; k-th smallest ⇔ index `k - 1`.
   Fix the convention in a comment on line 1 (*"k is a 0-based sorted index"*)
   and convert once at the call site — never juggle both conventions inside
   the loop.
2. **Recursing/looping on the wrong side, or keeping the pivot in the window.**
   After Lomuto, the pivot at `p` is FINAL — the new window must *exclude* it
   (`lo = p+1` or `hi = p-1`). Keeping it risks an infinite loop.
3. **Forgetting to randomise the pivot.** Always-last-element pivot is O(n²) on
   already-sorted input — a classic hidden test case (LeetCode 215 has one).
   One `rand.Intn` line fixes it. (Shuffling the whole array first also works.)
4. **Duplicates + Lomuto = O(n²).** On `[7,7,7,…,7]`, strict `<` puts everything
   on one side every round. If constraints allow heavy duplication, use Hoare or a
   three-way partition (`< pivot | == pivot | > pivot`, then check which band
   contains `k`).
5. **Mixing up Hoare's return value with the pivot position.** Hoare's `j` is a
   *split point*, not the pivot's final index — you cannot early-return on
   `j == k`, and the left recursion must be `[lo, j]` (inclusive), not `[lo, j-1]`.
   When unsure, use Lomuto.
6. **Infinite loop in Hoare from pivot choice.** With Hoare, pick the pivot
   *value* from `nums[lo]`-side (or copy the value before swapping); choosing
   `nums[hi]` as pivot value with this do-while structure can loop forever when
   `j` never moves past `lo`.
7. **Mutating input that must stay intact.** Quickselect reorders the array. If
   the caller needs the original order, copy first (costs O(n) space — say so).
8. **Using quickselect on streaming / huge-N data.** If elements arrive one at a
   time or don't fit in memory, quickselect is inapplicable — use a size-k heap
   (O(n log k), O(k) space). Interviewers often ask this exact follow-up.
9. **`rand.Intn(0)` panic.** Guard the window: when `lo == hi` the answer is
   `nums[lo]`; templates above never call `Intn` with a zero span because
   `hi-lo+1 >= 1`, but a hand-rolled variant that computes `Intn(hi-lo)` will
   panic on a width-1 window.

---

## Decision cheat-sheet: quickselect vs heap vs sort

| Situation | Use |
|---|---|
| One-shot k-th statistic, array in memory, mutation OK | **Quickselect** — O(n) avg |
| Streaming input / n unbounded / can't hold all data | **Heap of size k** — O(n log k) |
| Need the top-k **in sorted order** | Heap, or quickselect then sort the k prefix: O(n + k log k) |
| Need many different order statistics of the same array | **Sort once** — O(n log n) |
| Hard worst-case guarantee demanded | Median-of-medians quickselect — O(n) worst |

---

## Problems in this repo

*None of the currently-written problems (0001–0016) use quickselect.*

Classic quickselect problems that belong here once solved (a later pass will
link them as their folders land, including the in-flight 0131–0400 batch):

- LeetCode #215 — Kth Largest Element in an Array (the canonical problem)
- LeetCode #347 — Top K Frequent Elements (quickselect on frequency pairs)
- LeetCode #373 / #378 — kth-smallest variants where quickselect is one option
- LeetCode #462 — Minimum Moves to Equal Array Elements II (median via quickselect)
- LeetCode #973 — K Closest Points to Origin (quickselect on squared distance)
