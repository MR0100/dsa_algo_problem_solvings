# 0327 — Count of Range Sum

> LeetCode #327 · Difficulty: Hard
> **Categories:** Array, Binary Search, Divide and Conquer, Binary Indexed Tree, Segment Tree, Merge Sort, Ordered Set

---

## Problem Statement

Given an integer array `nums` and two integers `lower` and `upper`, return *the number of range sums that lie in `[lower, upper]` inclusive*.

Range sum `S(i, j)` is defined as the sum of the elements in `nums` between indices `i` and `j` inclusive, where `i <= j`.

**Example 1:**

```
Input: nums = [-2,5,-1], lower = -2, upper = 2
Output: 3
Explanation: The three ranges are [0,0], [2,2], and [0,2] and their respective sums are -2, -1, 2.
```

**Example 2:**

```
Input: nums = [0], lower = 0, upper = 0
Output: 1
```

**Constraints:**

- `1 <= nums.length <= 10^5`
- `-2^31 <= nums[i] <= 2^31 - 1`
- `-10^5 <= lower <= upper <= 10^5`
- The answer is guaranteed to fit in a 32-bit integer.

---

## Company Frequency

| Company    | Frequency        | Last Reported |
|------------|------------------|---------------|
| Google     | ★★★★☆ High       | 2024          |
| Amazon     | ★★★☆☆ Medium     | 2024          |
| Microsoft  | ★★★☆☆ Medium     | 2023          |
| Meta       | ★★★☆☆ Medium     | 2023          |
| Bloomberg  | ★★★☆☆ Medium     | 2023          |
| Apple      | ★★☆☆☆ Low        | 2022          |

> ⚠️ Frequency data is crowd-sourced from LeetCode Discuss, Glassdoor, and
> community interview reports. Treat as a signal, not a guarantee.

---

## DSA Concepts Used

- **Prefix Sums** — the whole problem is reduced to prefix sums: with `P[0]=0` and `P[k]=nums[0]+...+nums[k-1]`, the range sum `S(i,j) = P[j+1] - P[i]`, so we count pairs `a < b` with `lower <= P[b]-P[a] <= upper` → see [`/dsa/prefix_sum.md`](/dsa/prefix_sum.md)
- **Divide and Conquer / Merge Sort** — the canonical `O(n log n)` solution sorts the prefix-sum array and counts qualifying cross-half pairs during the merge, in the same family as counting inversions → see [`/dsa/divide_and_conquer.md`](/dsa/divide_and_conquer.md) and [`/dsa/sorting.md`](/dsa/sorting.md)
- **Binary Indexed Tree (Fenwick)** — streaming the prefix sums into a Fenwick tree over compressed coordinates lets each "how many earlier prefix sums fall in `[P[b]-upper, P[b]-lower]`" query run in `O(log n)` → see [`/dsa/segment_tree_fenwick.md`](/dsa/segment_tree_fenwick.md)
- **Coordinate Compression + Binary Search** — prefix sums and their query bounds span a huge int64 range, so we sort/dedup them to dense ranks and use binary search (`sort.Search`) to map a value to its rank → see [`/dsa/sorting.md`](/dsa/sorting.md)

---

## Approaches Overview

| # | Approach | Time | Space | When to use |
|---|----------|------|-------|-------------|
| 1 | Brute Force | O(n²) | O(n) | Baseline / correctness oracle; fine for n ≤ a few thousand, TLE at n = 10⁵ |
| 2 | Merge Sort Count (Divide and Conquer) | O(n log n) | O(n) | The standard interview answer; clean and cache-friendly |
| 3 | Binary Indexed Tree (Fenwick) (Optimal) | O(n log n) | O(n) | When you already think in Fenwick/BIT terms or need an online-style sweep |

---

## Approach 1 — Brute Force

### Intuition

Once you build the prefix-sum array, every range sum is a single subtraction: `S(i, j) = P[j+1] - P[i]`. So instead of re-summing sub-arrays, iterate over all pairs of prefix indices `a < b` and test whether `P[b] - P[a]` lands in `[lower, upper]`. This turns an `O(n³)` naive re-summation into `O(n²)` — still too slow for `n = 10⁵`, but it is the ground truth we validate the fast solutions against.

### Algorithm

1. Build `P[0..n]` with `P[0] = 0` and `P[k] = P[k-1] + nums[k-1]` (in `int64`).
2. For every pair `a < b`, compute `diff = P[b] - P[a]`.
3. If `lower <= diff <= upper`, increment the counter.
4. Return the counter.

### Complexity

- **Time:** O(n²) — the double loop touches every pair of the `n+1` prefix indices.
- **Space:** O(n) — the prefix-sum array; the counter is O(1).

### Code

```go
func bruteForce(nums []int, lower int, upper int) int {
	n := len(nums)
	prefix := make([]int64, n+1) // prefix[0] = 0 already (zero value)
	for i := 0; i < n; i++ {
		prefix[i+1] = prefix[i] + int64(nums[i]) // running prefix sum in int64
	}

	lo, hi := int64(lower), int64(upper) // widen bounds once to compare cleanly
	count := 0
	for a := 0; a <= n; a++ { // a is the "start" prefix index
		for b := a + 1; b <= n; b++ { // b > a is the "end" prefix index
			diff := prefix[b] - prefix[a] // this equals range sum S(a, b-1)
			if diff >= lo && diff <= hi { // inside the inclusive window?
				count++ // valid range sum found
			}
		}
	}
	return count
}
```

### Dry Run

Example 1: `nums = [-2, 5, -1]`, `lower = -2`, `upper = 2`.

Prefix sums: `P = [0, -2, 3, 2]` (`P[0]=0`, `P[1]=-2`, `P[2]=-2+5=3`, `P[3]=3-1=2`).

| Pair (a,b) | `P[b] - P[a]` | value | in `[-2, 2]`? | running count |
|------------|---------------|-------|---------------|---------------|
| (0,1) | `-2 - 0`  | -2 | yes | 1 |
| (0,2) | `3 - 0`   | 3  | no  | 1 |
| (0,3) | `2 - 0`   | 2  | yes | 2 |
| (1,2) | `3 - (-2)`| 5  | no  | 2 |
| (1,3) | `2 - (-2)`| 4  | no  | 2 |
| (2,3) | `2 - 3`   | -1 | yes | 3 |

Result: `3` ✔ — the three qualifying pairs are `(0,1) → S(0,0) = -2`, `(0,3) → S(0,2) = 2`, and `(2,3) → S(2,2) = -1`, matching ranges `[0,0]`, `[0,2]`, `[2,2]`.

---

## Approach 2 — Merge Sort Count (Divide and Conquer)

### Intuition

We need to count pairs `a < b` with `lower <= P[b] - P[a] <= upper`. This is the same shape as counting inversions, and merge sort solves it the same way. Split the prefix array in half. Any qualifying pair is either fully inside the left half, fully inside the right half (handled by recursion), or straddles the split with `a` on the left and `b` on the right. Crucially, `a < b` holds automatically for straddling pairs because every left index is smaller than every right index. Once both halves are **sorted by value**, for a fixed left value `P[a]` the right values `P[b]` satisfying `P[a]+lower <= P[b] <= P[a]+upper` form a contiguous window; as `P[a]` increases, that window only slides right, so two monotone pointers sweep all of it in linear time. Then merge the halves so the parent recursion sees a sorted block.

### Algorithm

1. Build `int64` prefix sums `P[0..n]`.
2. `sortCount(left, right)`: if the block has ≤ 1 element, return 0 (already sorted). Otherwise pick `mid` and recurse into `[left, mid]` and `[mid+1, right]`, summing their counts.
3. Count cross pairs: keep two pointers `low` and `high` into the right half. For each `a` from `left` to `mid`, advance `low` to the first right index with `P[low] - P[a] >= lower`, and `high` to the first right index with `P[high] - P[a] > upper`. Add `high - low` (the size of the valid window) to the count.
4. Merge the two sorted halves into a scratch buffer and copy back.
5. Return the total count from `sortCount(0, n)`.

### Complexity

- **Time:** O(n log n) — merge-sort recursion of depth `log n`; per level the counting sweep and the merge are each linear (both pointers only move forward).
- **Space:** O(n) — the prefix array plus one scratch buffer, plus O(log n) recursion stack.

### Code

```go
func mergeSortCount(nums []int, lower int, upper int) int {
	n := len(nums)
	prefix := make([]int64, n+1)
	for i := 0; i < n; i++ {
		prefix[i+1] = prefix[i] + int64(nums[i]) // int64 prefix sums
	}

	lo, hi := int64(lower), int64(upper)
	scratch := make([]int64, len(prefix)) // reusable merge buffer

	// sortCount sorts prefix[left..right] in place and returns the number of
	// valid cross pairs discovered while merging its two halves.
	var sortCount func(left, right int) int
	sortCount = func(left, right int) int {
		if left >= right {
			return 0 // 0 or 1 element: nothing to pair, already sorted
		}
		mid := left + (right-left)/2 // split point (overflow-safe midpoint)
		// Count pairs fully inside each half first (both halves get sorted here).
		count := sortCount(left, mid) + sortCount(mid+1, right)

		// Count cross pairs: a in [left, mid], b in [mid+1, right].
		low, high := mid+1, mid+1 // moving window bounds into the right half
		for a := left; a <= mid; a++ {
			// low = first right index with prefix[low] - prefix[a] >= lower.
			for low <= right && prefix[low]-prefix[a] < lo {
				low++
			}
			// high = first right index with prefix[high] - prefix[a] > upper.
			for high <= right && prefix[high]-prefix[a] <= hi {
				high++
			}
			// [low, high) is the window of valid right partners for this a.
			count += high - low
		}

		// Standard merge of the two sorted halves into scratch, then copy back.
		i, j, k := left, mid+1, left
		for i <= mid && j <= right {
			if prefix[i] <= prefix[j] { // take the smaller front element
				scratch[k] = prefix[i]
				i++
			} else {
				scratch[k] = prefix[j]
				j++
			}
			k++
		}
		for i <= mid { // drain any remaining left elements
			scratch[k] = prefix[i]
			i++
			k++
		}
		for j <= right { // drain any remaining right elements
			scratch[k] = prefix[j]
			j++
			k++
		}
		copy(prefix[left:right+1], scratch[left:right+1]) // write sorted block back
		return count
	}

	return sortCount(0, n) // sort/count over all n+1 prefix indices
}
```

### Dry Run

Example 1: `nums = [-2, 5, -1]`, `lower = -2`, `upper = 2`. Prefix sums `P = [0, -2, 3, 2]`, indices `0..3`. Call `sortCount(0, 3)`, `mid = 1`.

**Recurse left `sortCount(0, 1)`** on values `[P0=0, P1=-2]`, `mid = 0`:
- Left `[0,0]` and right `[1,1]` are single elements → 0 each.
- Cross count, `a = 0` (`P[0]=0`): advance `low` while `P[low] - 0 < -2` — `P[1] = -2`, not `< -2`, so `low` stays at 1. Advance `high` while `P[high] - 0 <= 2` — `P[1] = -2 <= 2`, so `high → 2` (past `right=1`). Window `high - low = 2 - 1 = 1`. **Count += 1** (pair `P[1]-P[0] = -2`, in range).
- Merge sorts this block to `[-2, 0]`. So now `P = [-2, 0, 3, 2]`.

**Recurse right `sortCount(2, 3)`** on values `[P2=3, P3=2]`, `mid = 2`:
- Singletons → 0 each.
- Cross count, `a = 2` (`P[2]=3`): advance `low` while `P[low] - 3 < -2` — `P[3] = 2`, `2 - 3 = -1`, not `< -2`, `low` stays 3. Advance `high` while `P[high] - 3 <= 2` — `2 - 3 = -1 <= 2`, `high → 4`. Window `4 - 3 = 1`. **Count += 1** (pair `P[3]-P[2] = -1`, in range).
- Merge sorts to `[2, 3]`. Now `P = [-2, 0, 2, 3]`.

**Back in `sortCount(0, 3)`**, cross count over sorted halves left `[-2, 0]` (indices 0..1), right `[2, 3]` (indices 2..3). `low = high = 2`:
- `a = 0` (`P[0] = -2`): `low` while `P[low] - (-2) < -2` → `P[2] = 2`, `2+2 = 4`, not `< -2`, `low` stays 2. `high` while `P[high] + 2 <= 2` → `P[2]=2`, `2+2 = 4 <= 2`? no, `high` stays 2. Window 0.
- `a = 1` (`P[1] = 0`): `low` while `P[low] - 0 < -2` → `P[2] = 2`, not `< -2`, `low` stays 2. `high` while `P[high] - 0 <= 2` → `P[2] = 2 <= 2`, `high → 3`; `P[3] = 3 <= 2`? no, stop. Window `3 - 2 = 1`. **Count += 1** (a straddling pair: original `P[3] - P[1] = 2 - (-2) = ... ` in sorted terms `2 - 0 = 2`, corresponds to `S(0,2) = 2`, in range).

Total: `1 (left) + 1 (right) + 1 (cross) = 3` ✔

---

## Approach 3 — Binary Indexed Tree (Fenwick) (Optimal)

### Intuition

Sweep the prefix sums left to right. At step `b`, everything before it (`P[a]` for `a < b`) is already recorded, so the count of valid pairs ending at `b` is "how many recorded values lie in `[P[b] - upper, P[b] - lower]`" — this rearranges `lower <= P[b] - P[a] <= upper` for the unknown `P[a]`. A Fenwick tree answers "how many recorded values ≤ X" in `O(log n)`; subtract two such prefix counts to get a range count. Because the values span a huge `int64` range, first **coordinate-compress**: gather every value that could ever be inserted (`P[v]`) or used as a query bound (`P[v]-lower`, `P[v]-upper`), sort and dedup them to dense ranks, and binary-search to map a value to its rank. Then insert `P[b]` so future ends can use it.

### Algorithm

1. Build `int64` prefix sums `P[0..n]`.
2. Collect all values `P[v]`, `P[v]-lower`, `P[v]-upper`; sort and dedup into `uniq`. Define `rank(x)` = 1-based index of the first `uniq` value `>= x` (via `sort.Search`).
3. Maintain a Fenwick tree of counts over ranks. For each `p = P[b]` in order: let `L = p - upper`, `R = p - lower`. Add `query(rank(R+1)-1) - query(rank(L)-1)` — the number of already-inserted values in `[L, R]` — to the answer, then `update(rank(p))`.
4. Return the answer.

### Complexity

- **Time:** O(n log n) — building/sorting the coordinate list is `O(n log n)`; each of the `n+1` steps does `O(log n)` Fenwick work plus `O(log n)` binary searches.
- **Space:** O(n) — the coordinate list and the Fenwick tree.

### Code

```go
func binaryIndexedTree(nums []int, lower int, upper int) int {
	n := len(nums)
	prefix := make([]int64, n+1)
	for i := 0; i < n; i++ {
		prefix[i+1] = prefix[i] + int64(nums[i]) // int64 prefix sums
	}
	lo, hi := int64(lower), int64(upper)

	// Gather every value that will ever be inserted or used as a query bound.
	all := make([]int64, 0, 3*(n+1))
	for _, p := range prefix {
		all = append(all, p)    // the value we insert
		all = append(all, p-lo) // R bound: P[b] - lower
		all = append(all, p-hi) // L bound: P[b] - upper
	}
	sort.Slice(all, func(i, j int) bool { return all[i] < all[j] }) // sort ascending
	// Dedup in place so equal values share one rank.
	uniq := all[:0:0] // fresh slice, cap 0, so appends don't alias `all`
	for i, v := range all {
		if i == 0 || v != all[i-1] {
			uniq = append(uniq, v)
		}
	}

	// rank returns the 1-based index of the first compressed value >= x.
	rank := func(x int64) int {
		return sort.Search(len(uniq), func(i int) bool { return uniq[i] >= x }) + 1
	}

	// Fenwick tree (1-indexed) storing counts of inserted compressed values.
	tree := make([]int, len(uniq)+1)
	update := func(i int) { // add 1 at position i
		for ; i < len(tree); i += i & (-i) { // climb via lowest-set-bit jumps
			tree[i]++
		}
	}
	query := func(i int) int { // prefix count over positions [1, i]
		s := 0
		for ; i > 0; i -= i & (-i) { // descend via lowest-set-bit jumps
			s += tree[i]
		}
		return s
	}

	count := 0
	for _, p := range prefix { // p plays the role of P[b], an "end" prefix sum
		L := p - hi // lower edge of the valid P[a] interval: P[b] - upper
		R := p - lo // upper edge of the valid P[a] interval: P[b] - lower
		// Count already-inserted values in [L, R].
		// rank(L) is the first value >= L; rank(R+1)-1 is the last value <= R.
		rightRank := rank(R+1) - 1 // last compressed index with value <= R
		leftRank := rank(L)        // first compressed index with value >= L
		count += query(rightRank) - query(leftRank-1)
		update(rank(p)) // now P[b] is available to future ends as a P[a]
	}
	return count
}
```

### Dry Run

Example 1: `nums = [-2, 5, -1]`, `lower = -2`, `upper = 2`. Prefix sums `P = [0, -2, 3, 2]`.

We process each `p = P[b]` in order. For each, `L = p - upper = p - 2`, `R = p - lower = p + 2`, count inserted values in `[L, R]`, then insert `p`. (Compressed ranks are handled internally; here we reason with raw values — the tree only contains values inserted so far.)

| Step | `p = P[b]` | `[L, R] = [p-2, p+2]` | inserted values so far | count in `[L,R]` | running answer | then insert |
|------|-----------|-----------------------|------------------------|------------------|----------------|-------------|
| 1 | `0`  | `[-2, 2]` | {} (empty) | 0 | 0 | insert 0 |
| 2 | `-2` | `[-4, 0]` | {0} | 0 is in `[-4,0]` → 1 | 1 | insert -2 |
| 3 | `3`  | `[1, 5]`  | {0, -2} | none in `[1,5]` → 0 | 1 | insert 3 |
| 4 | `2`  | `[0, 4]`  | {0, -2, 3} | 0 and 3 are in `[0,4]` → 2 | 3 | insert 2 |

Result: `3` ✔ — step 2 catches `S` corresponding to `P[1]-P[0] = -2`, and step 4 catches `P[3]-P[0] = 2` and `P[3]-P[2] = -1`.

---

## Key Takeaways

- **Range-sum-in-window collapses to pair counting on prefix sums.** `S(i,j) ∈ [lower, upper]` ⟺ `P[b] - P[a] ∈ [lower, upper]` with `a < b`. Learn to see this reduction — it powers this problem and #560, #523, etc.
- **"Count pairs with a bounded difference" is a merge-sort / BIT template.** The same machinery counts inversions (#493 Reverse Pairs, #315 Count of Smaller). The merge version needs the halves sorted **by value** while the pair constraint `a < b` is preserved for free by index order.
- **Two monotone pointers turn each merge level linear.** Because both the left value `P[a]` and the window bounds `P[a]+lower`, `P[a]+upper` increase together, `low`/`high` never move backward — that is what keeps the count step `O(n)` per level.
- **Always use `int64` for prefix sums here.** Elements up to `2^31` over `10^5` terms reach `~2×10^{14}`, well past `int32`. Overflow is the single most common bug on this problem.
- **Fenwick trees need coordinate compression** when values are sparse over a huge range: gather every insert value **and** every query bound, sort, dedup, then index by rank. Forgetting to include the query bounds (`p-lower`, `p-upper`) is a classic mistake.

---

## Related Problems

- LeetCode #315 — Count of Smaller Numbers After Self (same merge-sort / BIT counting pattern)
- LeetCode #493 — Reverse Pairs (bounded-difference pair counting via merge sort)
- LeetCode #307 — Range Sum Query - Mutable (Fenwick / segment tree for prefix sums)
- LeetCode #53 — Maximum Subarray (prefix-sum reasoning over sub-array sums)
