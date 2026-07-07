# 0493 — Reverse Pairs

> LeetCode #493 · Difficulty: Hard
> **Categories:** Divide and Conquer, Binary Indexed Tree, Segment Tree, Merge Sort, Array, Ordered Set

---

## Problem Statement

Given an integer array `nums`, return *the number of **reverse pairs** in the array*.

A **reverse pair** is a pair `(i, j)` where:

- `0 <= i < j < nums.length` and
- `nums[i] > 2 * nums[j]`.

**Example 1:**

```
Input: nums = [1,3,2,3,1]
Output: 2
Explanation: The reverse pairs are:
(1, 4) --> nums[1] = 3, nums[4] = 1, 3 > 2 * 1
(3, 4) --> nums[3] = 3, nums[4] = 1, 3 > 2 * 1
```

**Example 2:**

```
Input: nums = [2,4,3,5,1]
Output: 3
Explanation: The reverse pairs are:
(1, 4) --> nums[1] = 4, nums[4] = 1, 4 > 2 * 1
(2, 4) --> nums[2] = 3, nums[4] = 1, 3 > 2 * 1
(3, 4) --> nums[3] = 5, nums[4] = 1, 5 > 2 * 1
```

**Constraints:**

- `1 <= nums.length <= 5 * 10^4`
- `-2^31 <= nums[i] <= 2^31 - 1`

---

## Company Frequency

| Company    | Frequency        | Last Reported |
|------------|------------------|---------------|
| Google     | ★★★★☆ High       | 2024          |
| Amazon     | ★★★☆☆ Medium     | 2023          |
| Microsoft  | ★★★☆☆ Medium     | 2023          |
| Meta       | ★★★☆☆ Medium     | 2023          |
| Bloomberg  | ★★☆☆☆ Low        | 2022          |
| Uber       | ★★☆☆☆ Low        | 2022          |

> ⚠️ Frequency data is crowd-sourced from LeetCode Discuss, Glassdoor, and
> community interview reports. Treat as a signal, not a guarantee.

---

## DSA Concepts Used

- **Divide and Conquer (modified merge sort)** — count cross-pairs during the merge step where both halves are already sorted, mirroring the classic "count inversions" algorithm generalized to the `nums[i] > 2·nums[j]` condition → see [`/dsa/divide_and_conquer.md`](/dsa/divide_and_conquer.md)
- **Binary Indexed Tree (Fenwick) + coordinate compression** — sweep `j` left→right and use prefix-frequency queries to count how many already-seen values exceed `2·nums[j]` → see [`/dsa/segment_tree_fenwick.md`](/dsa/segment_tree_fenwick.md)

---

## Approaches Overview

| # | Approach | Time | Space | When to use |
|---|----------|------|-------|-------------|
| 1 | Brute Force (All Pairs) | O(n²) | O(1) | Tiny inputs / correctness oracle; TLE at n = 5·10⁴ |
| 2 | Merge Sort (Divide & Conquer, Optimal) | O(n log n) | O(n) | The canonical answer; no compression bookkeeping |
| 3 | BIT + Coordinate Compression | O(n log n) | O(n) | Same complexity; reusable "count values > x seen so far" pattern |

---

## Approach 1 — Brute Force (All Pairs)

### Intuition

The definition is a nested loop: for every earlier index `i` and later index `j`, test `nums[i] > 2·nums[j]`. Count the ones that pass. It is the literal transcription of the problem and serves as the ground truth for the fast solutions. (Note: `2·nums[j]` can exceed `int32`, but Go's `int` is 64-bit, so plain multiplication is safe within the `[-2³¹, 2³¹−1]` value range.)

### Algorithm

1. `count = 0`.
2. For `i` in `0..n-1`, for `j` in `i+1..n-1`: if `nums[i] > 2·nums[j]`, `count++`.
3. Return `count`.

### Complexity

- **Time:** O(n²) — inspects all `~n²/2` ordered pairs; `2.5·10⁹` for `n = 5·10⁴` → TLE.
- **Space:** O(1).

### Code

```go
func bruteForce(nums []int) int {
	count := 0
	for i := 0; i < len(nums); i++ {
		for j := i + 1; j < len(nums); j++ {
			// int is 64-bit here, so 2*nums[j] cannot overflow the range.
			if nums[i] > 2*nums[j] {
				count++ // (i, j) is a reverse pair
			}
		}
	}
	return count
}
```

### Dry Run

Example 1: `nums = [1,3,2,3,1]`. Only pairs that satisfy `nums[i] > 2·nums[j]` are listed.

| i | nums[i] | j | nums[j] | nums[i] > 2·nums[j]? | count |
|---|---------|---|---------|----------------------|-------|
| 1 | 3 | 4 | 1 | 3 > 2 → yes | 1 |
| 3 | 3 | 4 | 1 | 3 > 2 → yes | 2 |

Every other pair fails (e.g. `i=0,j=…`: `1 > 2·x` never holds for these `x`). Total: `2` ✔

---

## Approach 2 — Merge Sort (Divide and Conquer, Optimal)

### Intuition

Split the array in half. Every reverse pair is (a) fully inside the left half, (b) fully inside the right half, or (c) **crossing** — `i` in the left, `j` in the right. Recursion handles (a) and (b). For (c), once both halves are sorted the crossing pairs collapse to a linear sweep: run a pointer `i` up the (sorted) left half and a pointer `j` up the (sorted) right half; for each `i`, advance `j` while `left[i] > 2·right[j]`. Since `left[i]` only grows, `j` never rewinds — so counting all crossing pairs is O(n) per level. After counting, merge the halves so the parent level can repeat the trick. **Count first, then merge** (the count needs the "raw" sorted halves before interleaving).

### Algorithm

1. If the segment has `< 2` elements, return `0`.
2. Recurse on `[lo..mid]` and `[mid+1..hi]`, summing their counts (this also sorts each half).
3. Crossing count: `j = mid+1`; for `i` from `lo` to `mid`, advance `j` while `j ≤ hi` and `arr[i] > 2·arr[j]`; add `j − (mid+1)` to the count.
4. Standard merge of the two sorted halves into place via a temp buffer.
5. Return the total count.

### Complexity

- **Time:** O(n log n) — `log n` levels of recursion, each doing O(n) counting plus O(n) merging.
- **Space:** O(n) — the reusable merge buffer (plus O(log n) recursion stack).

### Code

```go
func mergeSortCount(nums []int) int {
	if len(nums) < 2 {
		return 0
	}
	arr := make([]int, len(nums)) // work on a copy so the input stays intact
	copy(arr, nums)
	tmp := make([]int, len(nums)) // reusable merge buffer
	return mergeCount(arr, tmp, 0, len(arr)-1)
}

func mergeCount(arr, tmp []int, lo, hi int) int {
	if lo >= hi {
		return 0 // 0 or 1 element: no pairs
	}
	mid := lo + (hi-lo)/2
	// 1) Count pairs fully inside each half (and sort each half).
	count := mergeCount(arr, tmp, lo, mid) + mergeCount(arr, tmp, mid+1, hi)

	// 2) Count crossing pairs. Both arr[lo..mid] and arr[mid+1..hi] are sorted.
	//    For each i on the left, extend j on the right while arr[i] > 2*arr[j].
	j := mid + 1
	for i := lo; i <= mid; i++ {
		// arr grows with i, so j never needs to move backward (monotonic).
		for j <= hi && arr[i] > 2*arr[j] {
			j++
		}
		count += j - (mid + 1) // all right elements before j pair with arr[i]
	}

	// 3) Merge arr[lo..mid] and arr[mid+1..hi] into sorted order via tmp.
	i, k, r := lo, lo, mid+1
	for i <= mid && r <= hi {
		if arr[i] <= arr[r] {
			tmp[k] = arr[i]
			i++
		} else {
			tmp[k] = arr[r]
			r++
		}
		k++
	}
	for i <= mid { // drain leftovers from the left half
		tmp[k] = arr[i]
		i++
		k++
	}
	for r <= hi { // drain leftovers from the right half
		tmp[k] = arr[r]
		r++
		k++
	}
	copy(arr[lo:hi+1], tmp[lo:hi+1]) // write the merged run back in place
	return count
}
```

### Dry Run

Example 1: `nums = [1,3,2,3,1]`. The recursion splits into `[1,3]` | `[2,3,1]`. We show the **crossing-count** phase for the top-level merge, at which point the halves have already been recursively sorted to `left = [1,3]` and `right = [1,2,3]`.

Crossing count with sorted `left = [1,3]`, sorted `right = [1,2,3]` (`j` starts at first right index):

| i | left[i] | advance j while left[i] > 2·right[j] | pairs added (j − start) |
|---|---------|--------------------------------------|-------------------------|
| — | 1 | `1 > 2·1`? no → j stays | 0 |
| — | 3 | `3 > 2·1`? yes → j++; `3 > 2·2`? no → stop | 1 |

Crossing pairs at the top level: `1`. The remaining `1` comes from within the right subtree (`[2,3,1]` → its own merge finds the pair `(3, 1)` where `3 > 2·1`). Left subtree `[1,3]` contributes `0`. Total `= 0 (left) + 1 (right subtree) + 1 (crossing) = 2` ✔

---

## Approach 3 — Binary Indexed Tree + Coordinate Compression

### Intuition

Process indices left→right. When the sweep reaches `j`, every `i < j` has already been recorded in a value-frequency structure. A reverse pair needs `nums[i] > 2·nums[j]`, i.e. we want **"how many already-seen values are strictly greater than `2·nums[j]`"**. A Fenwick tree gives prefix-count queries in O(log n): count of seen values `> 2·nums[j]` = `insertedSoFar − (count of seen values ≤ 2·nums[j])`. The values (`nums[i]` and `2·nums[j]`) span a huge range, so we **coordinate-compress** the union of all `nums[i]` and all `2·nums[j]` down to ranks `1..m` first.

### Algorithm

1. Build a sorted, de-duplicated array of every `nums[i]` and every `2·nums[i]`; map each value to a 1-based rank.
2. For `j` from `0` to `n-1`:
   a. `r = query(rank(2·nums[j]))` = number of inserted values with rank ≤ that; add `inserted − r` to the answer (values strictly greater).
   b. `update(rank(nums[j]))`; increment `inserted` (now `nums[j]` is a valid future "i").
3. Return the answer.

### Complexity

- **Time:** O(n log n) — one `sort` for compression, then `n` updates and `n` queries, each O(log m) with `m ≤ 2n`.
- **Space:** O(n) — the compressed value array and the Fenwick tree.

### Code

```go
func bitCount(nums []int) int {
	n := len(nums)
	if n < 2 {
		return 0
	}
	// 1) Collect every value we will ever query or insert, then compress.
	vals := make([]int, 0, 2*n)
	for _, v := range nums {
		vals = append(vals, v)   // values we INSERT (the nums[i])
		vals = append(vals, 2*v) // values we QUERY against (the 2*nums[j])
	}
	sort.Ints(vals)
	uniq := vals[:0] // dedup in place
	for i, v := range vals {
		if i == 0 || v != vals[i-1] {
			uniq = append(uniq, v)
		}
	}
	// rank returns the 1-based position of x in the sorted unique list.
	rank := func(x int) int {
		return sort.SearchInts(uniq, x) + 1 // +1 → Fenwick indices start at 1
	}

	tree := make([]int, len(uniq)+1) // Fenwick tree of value frequencies
	// update adds 1 at position i.
	update := func(i int) {
		for ; i < len(tree); i += i & (-i) {
			tree[i]++
		}
	}
	// query returns the count of inserted values with rank in [1, i].
	query := func(i int) int {
		s := 0
		for ; i > 0; i -= i & (-i) {
			s += tree[i]
		}
		return s
	}

	answer := 0
	inserted := 0 // how many nums[i] are already in the tree (i < j)
	for j := 0; j < n; j++ {
		// Count already-inserted values strictly greater than 2*nums[j]:
		//   inserted - (# values with rank <= rank(2*nums[j])).
		r := query(rank(2 * nums[j]))
		answer += inserted - r
		update(rank(nums[j])) // now nums[j] becomes an eligible "i" for later j
		inserted++
	}
	return answer
}
```

### Dry Run

Example 1: `nums = [1,3,2,3,1]`. Compressed unique values (union of `nums` and `2·nums = {2,6,4,6,2}`): `[1,2,3,4,6]` → ranks `1:1, 2:2, 3:3, 4:4, 6:5`.

| j | nums[j] | 2·nums[j] | rank(2·nums[j]) | query = seen ≤ that | inserted | added = inserted − query | then insert nums[j] (rank) | answer |
|---|---------|-----------|-----------------|---------------------|----------|--------------------------|----------------------------|--------|
| 0 | 1 | 2 | 2 | 0 | 0 | 0 | rank(1)=1 | 0 |
| 1 | 3 | 6 | 5 | 1 (the `1`) | 1 | 0 | rank(3)=3 | 0 |
| 2 | 2 | 4 | 4 | 1 (the `1`; `3`→rank3 ≤4 also) → 2 | 2 | 0 | rank(2)=2 | 0 |
| 3 | 3 | 6 | 5 | 3 (all of 1,2,3 ≤ 6) | 3 | 0 | rank(3)=3 | 0 |
| 4 | 1 | 2 | 2 | 2 (values ≤ 2 seen: the `1` and `2`) | 4 | 4 − 2 = **2** | rank(1)=1 | 2 |

Final answer: `2` ✔ — at `j = 4` (`nums[j] = 1`), the two seen `3`s exceed `2·1 = 2`, contributing both reverse pairs.

---

## Key Takeaways

- **"Count pairs across the array with an order condition" → merge sort or a Fenwick/segment tree.** Reverse Pairs is the generalized inversion-count: replace `left[i] > right[j]` with `left[i] > 2·right[j]` in the merge sweep.
- In the merge-sort version, **count the crossing pairs before merging**, using a separate monotonic `j` pointer; the merge itself then interleaves the halves. Do not fuse the two loops — the count needs the halves un-interleaved.
- The Fenwick approach embodies a reusable idiom: *sweep once, and at each element ask "how many previously seen values satisfy a range predicate?"* — coordinate-compress when values are sparse/huge.
- **Overflow watch:** `2·nums[j]` can exceed `int32`. In Go the default `int` is 64-bit so it is safe; in languages with 32-bit ints, widen to 64-bit (or compare `nums[i]/2.0 > nums[j]`).

---

## Related Problems

- LeetCode #315 — Count of Smaller Numbers After Self (BIT / merge sort, same sweep)
- LeetCode #327 — Count of Range Sum (merge sort on prefix sums, same cross-count trick)
- LeetCode #912 — Sort an Array (the merge-sort skeleton)
- LeetCode #307 — Range Sum Query - Mutable (Fenwick tree fundamentals)
- Classic — Counting Inversions (`nums[i] > nums[j]`, the parent problem)
