# 0373 — Find K Pairs with Smallest Sums

> LeetCode #373 · Difficulty: Medium
> **Categories:** Array, Heap (Priority Queue)

---

## Problem Statement

You are given two integer arrays `nums1` and `nums2` sorted in **non-decreasing order** and an integer `k`.

Define a pair `(u, v)` which consists of one element from the first array and one element from the second array.

Return *the* `k` *pairs* `(u1, v1), (u2, v2), ..., (uk, vk)` *with the smallest sums.*

**Example 1:**

```
Input: nums1 = [1,7,11], nums2 = [2,4,6], k = 3
Output: [[1,2],[1,4],[1,6]]
Explanation: The first 3 pairs are returned from the sequence:
             [1,2],[1,4],[1,6],[7,2],[7,4],[11,2],[7,6],[11,4],[11,6]
```

**Example 2:**

```
Input: nums1 = [1,1,2], nums2 = [1,2,3], k = 2
Output: [[1,1],[1,1]]
Explanation: The first 2 pairs are returned from the sequence:
             [1,1],[1,1],[1,2],[2,1],[1,2],[2,2],[1,3],[1,3],[2,3]
```

**Example 3:**

```
Input: nums1 = [1,2], nums2 = [3], k = 3
Output: [[1,3],[2,3]]
Explanation: All possible pairs are returned from the sequence: [1,3],[2,3]
```

**Constraints:**

- `1 <= nums1.length, nums2.length <= 10^5`
- `-10^9 <= nums1[i], nums2[i] <= 10^9`
- `nums1` and `nums2` both are sorted in non-decreasing order.
- `1 <= k <= 10^4`

---

## Company Frequency

| Company    | Frequency        | Last Reported |
|------------|------------------|---------------|
| Amazon     | ★★★★☆ High       | 2024          |
| Google     | ★★★☆☆ Medium     | 2023          |
| Microsoft  | ★★★☆☆ Medium     | 2023          |
| Uber       | ★★☆☆☆ Low        | 2022          |
| Bloomberg  | ★★☆☆☆ Low        | 2023          |

> ⚠️ Frequency data is crowd-sourced from LeetCode Discuss, Glassdoor, and community interview reports. Treat as a signal, not a guarantee.

---

## DSA Concepts Used

- **Min-Heap / Priority Queue** — always pop the globally smallest candidate sum next → see [`/dsa/heap_priority_queue.md`](/dsa/heap_priority_queue.md)
- **K-way merge / sorted-frontier expansion** — because both arrays are sorted, sums increase along rows and columns, so we expand a staircase boundary → see [`/dsa/k_way_merge.md`](/dsa/k_way_merge.md)
- **Sorting** — the brute-force baseline sorts all pairs by sum → see [`/dsa/sorting.md`](/dsa/sorting.md)

---

## Approaches Overview

| # | Approach | Time | Space | When to use |
|---|----------|------|-------|-------------|
| 1 | Brute force (all pairs + sort) | O(m·n log(m·n)) | O(m·n) | Small arrays; simplest correctness baseline |
| 2 | Min-heap over sorted frontier | O(k log k) | O(k) | Large arrays, small k (optimal) |

---

## Approach 1 — Brute Force

### Intuition
The literal reading: enumerate the full `m × n` grid of pairs, sort by sum, take the `k` smallest. Correct but wasteful when the grid is huge and `k` is small.

### Algorithm
1. Build every pair `[u, v]` for `u ∈ nums1`, `v ∈ nums2`.
2. Stable-sort by `u + v`.
3. Return the first `min(k, len)` pairs.

### Complexity
- **Time:** O(m·n log(m·n)) — building and sorting the entire grid.
- **Space:** O(m·n) — the list of all pairs.

### Code
```go
func bruteForce(nums1, nums2 []int, k int) [][]int {
	pairs := make([][]int, 0, len(nums1)*len(nums2))
	for _, u := range nums1 { // every element of the first array
		for _, v := range nums2 { // paired with every element of the second
			pairs = append(pairs, []int{u, v})
		}
	}
	// Sort ascending by the pair sum; ties keep their relative order.
	sort.SliceStable(pairs, func(i, j int) bool {
		return pairs[i][0]+pairs[i][1] < pairs[j][0]+pairs[j][1]
	})
	if k > len(pairs) { // never ask for more pairs than exist
		k = len(pairs)
	}
	return pairs[:k] // the k smallest-sum pairs
}
```

### Dry Run
`nums1 = [1,7,11]`, `nums2 = [2,4,6]`, `k = 3`:

| Step | Detail |
|------|--------|
| Build | `[1,2],[1,4],[1,6],[7,2],[7,4],[7,6],[11,2],[11,4],[11,6]` |
| Sums | `3,5,7,9,11,13,13,15,17` |
| Sort | already ascending (stable): `[1,2](3),[1,4](5),[1,6](7),[7,2](9),…` |
| Take k=3 | `[[1,2],[1,4],[1,6]]` |

Return `[[1,2],[1,4],[1,6]]`. ✓

---

## Approach 2 — Min-Heap over the Sorted Frontier (Optimal)

### Intuition
Picture pairs as a grid where row `i` uses `nums1[i]` and column `j` uses `nums2[j]`. Because both arrays are sorted, sums **increase down each row and across each column**. The globally smallest unused pair always lies on the "staircase" boundary. Seed the heap with the first column `(i, 0)` for each row (only the first `k` rows can matter). Each time we pop `(i, j)`, the only new candidate that could now be minimal is its right neighbour `(i, j+1)`, so we push that. This is a k-way merge of the rows.

### Algorithm
1. Push `(i, 0)` for `i` in `0..min(len(nums1), k)-1` with sum `nums1[i]+nums2[0]`.
2. Repeat `k` times (while heap non-empty):
   1. Pop the smallest `(i, j)`; append `[nums1[i], nums2[j]]` to the answer.
   2. If `j+1 < len(nums2)`, push `(i, j+1)`.
3. Return the collected pairs.

### Complexity
- **Time:** O(k log k) — the heap holds at most ~k items; k pops, each with a push.
- **Space:** O(k) — heap plus output.

### Code
```go
func heapFrontier(nums1, nums2 []int, k int) [][]int {
	res := make([][]int, 0, k)
	if len(nums1) == 0 || len(nums2) == 0 || k == 0 {
		return res // nothing to pair
	}

	h := &minHeap{}
	heap.Init(h)
	// Seed with the first column: pairs (nums1[i], nums2[0]). Only the first k
	// rows can ever contribute to the k smallest sums.
	limit := len(nums1)
	if limit > k {
		limit = k
	}
	for i := 0; i < limit; i++ {
		heap.Push(h, pairItem{i: i, j: 0, sum: nums1[i] + nums2[0]})
	}

	// Pop k smallest sums, expanding the row's next column each time.
	for h.Len() > 0 && len(res) < k {
		it := heap.Pop(h).(pairItem)                       // current smallest-sum pair
		res = append(res, []int{nums1[it.i], nums2[it.j]}) // record it
		if it.j+1 < len(nums2) {                           // slide right within the same row
			heap.Push(h, pairItem{i: it.i, j: it.j + 1, sum: nums1[it.i] + nums2[it.j+1]})
		}
	}
	return res
}
```

### Dry Run
`nums1 = [1,7,11]`, `nums2 = [2,4,6]`, `k = 3`:

| Step | Heap (i,j,sum) after action | Popped | Output |
|------|------------------------------|--------|--------|
| Seed | (0,0,3),(1,0,9),(2,0,13) | — | [] |
| Pop 1 | pop (0,0,3); push (0,1,5) → heap {(0,1,5),(1,0,9),(2,0,13)} | [1,2] | [[1,2]] |
| Pop 2 | pop (0,1,5); push (0,2,7) → heap {(0,2,7),(1,0,9),(2,0,13)} | [1,4] | [[1,2],[1,4]] |
| Pop 3 | pop (0,2,7); j+1=3 out of range, no push | [1,6] | [[1,2],[1,4],[1,6]] |

`len(res) == k == 3` → stop. Return `[[1,2],[1,4],[1,6]]`. ✓

---

## Key Takeaways
- **Sorted inputs turn a global search into a frontier expansion.** Sums monotonically increase along rows/columns, so only staircase-boundary cells are candidates.
- **Seed the first column, expand right on each pop** — a clean k-way merge that never touches the full `m·n` grid.
- **Bound the seed by k** (`min(len(nums1), k)` rows): rows beyond the k-th can never appear in the first k smallest sums.
- Heap size stays O(k), giving O(k log k) regardless of how large the arrays are.

---

## Related Problems
- LeetCode #378 — Kth Smallest Element in a Sorted Matrix (same frontier idea)
- LeetCode #23 — Merge k Sorted Lists (k-way merge with a heap)
- LeetCode #264 — Ugly Number II (heap/pointer frontier)
- LeetCode #632 — Smallest Range Covering Elements from K Lists (heap over sorted lists)
