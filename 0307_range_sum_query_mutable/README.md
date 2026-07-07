# 0307 — Range Sum Query - Mutable

> LeetCode #307 · Difficulty: Medium
> **Categories:** Array, Design, Binary Indexed Tree, Segment Tree

---

## Problem Statement

Given an integer array `nums`, handle multiple queries of the following types:

1. **Update** the value of an element in `nums`.
2. Calculate the **sum** of the elements of `nums` between indices `left` and `right` **inclusive** where `left <= right`.

Implement the `NumArray` class:

- `NumArray(int[] nums)` Initializes the object with the integer array `nums`.
- `void update(int index, int val)` **Updates** the value of `nums[index]` to be `val`.
- `int sumRange(int left, int right)` Returns the **sum** of the elements of `nums` between indices `left` and `right` **inclusive** (i.e. `nums[left] + nums[left + 1] + ... + nums[right]`).

**Example 1:**

```
Input:
["NumArray", "sumRange", "update", "sumRange"]
[[[1, 3, 5]], [0, 2], [1, 2], [0, 2]]
Output:
[null, 9, null, 8]

Explanation:
NumArray numArray = new NumArray([1, 3, 5]);
numArray.sumRange(0, 2); // return 1 + 3 + 5 = 9
numArray.update(1, 2);   // nums = [1, 2, 5]
numArray.sumRange(0, 2); // return 1 + 2 + 5 = 8
```

**Constraints:**

- `1 <= nums.length <= 3 * 10^4`
- `-100 <= nums[i] <= 100`
- `0 <= index < nums.length`
- `-100 <= val <= 100`
- `0 <= left <= right < nums.length`
- At most `3 * 10^4` calls will be made to `update` and `sumRange`.

---

## Company Frequency

| Company   | Frequency        | Last Reported |
|-----------|------------------|---------------|
| Google    | ★★★★☆ High       | 2024          |
| Amazon    | ★★★★☆ High       | 2024          |
| Microsoft | ★★★☆☆ Medium     | 2023          |
| Meta      | ★★★☆☆ Medium     | 2023          |
| Bloomberg | ★★★☆☆ Medium     | 2023          |

> ⚠️ Frequency data is crowd-sourced from LeetCode Discuss, Glassdoor, and community interview reports. Treat as a signal, not a guarantee.

---

## DSA Concepts Used

- **Binary Indexed Tree (Fenwick)** — point update + prefix sum in O(log n) via lowbit jumps → see [`/dsa/segment_tree_fenwick.md`](/dsa/segment_tree_fenwick.md)
- **Segment Tree** — array-backed range-sum tree, generalizes to other aggregates → see [`/dsa/segment_tree_fenwick.md`](/dsa/segment_tree_fenwick.md)
- **Prefix Sum** — range sum expressed as prefix(r) − prefix(l−1); the Fenwick tree makes prefixes updatable → see [`/dsa/prefix_sum.md`](/dsa/prefix_sum.md)
- **Design Data Structures** — implementing a class with a fixed operation contract → see [`/dsa/design_data_structures.md`](/dsa/design_data_structures.md)

---

## Approaches Overview

| # | Approach | Update | SumRange | Space | When to use |
|---|----------|--------|----------|-------|-------------|
| 1 | Brute Force (array) | O(1) | O(n) | O(n) | Few queries, many updates |
| 2 | Fenwick Tree (Optimal) | O(log n) | O(log n) | O(n) | Interleaved update + sum; cleanest |
| 3 | Segment Tree | O(log n) | O(log n) | O(n) | When query later generalizes (min/max/lazy) |

---

## Approach 1 — Brute Force

### Intuition

Keep the raw array. `update` is one assignment; `sumRange` walks the window. This is the baseline the smarter structures must beat when updates and queries are interleaved.

### Algorithm

1. Store a copy of `nums`.
2. `Update(index, val)`: `nums[index] = val`.
3. `SumRange(left, right)`: loop `i` from `left` to `right` adding `nums[i]`.

### Complexity

- **Time:** Update O(1); SumRange O(n) — walks the whole window.
- **Space:** O(n) — the stored array.

### Code

```go
type BruteForce struct {
	nums []int
}

func NewBruteForce(nums []int) *BruteForce {
	cp := make([]int, len(nums))
	copy(cp, nums)
	return &BruteForce{nums: cp}
}

func (b *BruteForce) Update(index, val int) {
	b.nums[index] = val
}

func (b *BruteForce) SumRange(left, right int) int {
	sum := 0
	for i := left; i <= right; i++ {
		sum += b.nums[i]
	}
	return sum
}
```

### Dry Run

Ops on `nums = [1,3,5]`:

| Op | State | Computation | Output |
|----|-------|-------------|--------|
| `sumRange(0,2)` | [1,3,5] | 1+3+5 | 9 |
| `update(1,2)`   | [1,2,5] | nums[1]=2 | null |
| `sumRange(0,2)` | [1,2,5] | 1+2+5 | 8 |

---

## Approach 2 — Fenwick Tree (Optimal)

### Intuition

A Fenwick tree (Binary Indexed Tree) stores partial sums so that any prefix sum is the sum of O(log n) cells and a point update touches O(log n) cells. Cell `i` covers the range `(i − lowbit(i), i]` where `lowbit(i) = i & -i`. Range sum `[l, r] = prefix(r) − prefix(l−1)`. Since updates are "set to `val`", we store the current values and push the difference `val − old`.

### Algorithm

1. Build a 1-indexed tree of size `n+1`; `add` each initial value.
2. `add(i, d)`: while `i <= n`, `tree[i] += d`, `i += lowbit(i)`.
3. `prefix(i)`: while `i > 0`, `s += tree[i]`, `i -= lowbit(i)`.
4. `Update(index, val)`: `d = val − nums[index]`; store `val`; `add(index+1, d)`.
5. `SumRange(left, right)`: `prefix(right+1) − prefix(left)`.

### Complexity

- **Time:** Update O(log n), SumRange O(log n); Build O(n log n).
- **Space:** O(n) — the tree plus the value cache.

### Code

```go
type FenwickTree struct {
	n    int
	tree []int
	nums []int
}

func NewFenwickTree(nums []int) *FenwickTree {
	n := len(nums)
	f := &FenwickTree{n: n, tree: make([]int, n+1), nums: make([]int, n)}
	for i, v := range nums {
		f.nums[i] = v
		f.add(i+1, v)
	}
	return f
}

func (f *FenwickTree) add(i, delta int) {
	for ; i <= f.n; i += i & (-i) {
		f.tree[i] += delta
	}
}

func (f *FenwickTree) prefix(i int) int {
	s := 0
	for ; i > 0; i -= i & (-i) {
		s += f.tree[i]
	}
	return s
}

func (f *FenwickTree) Update(index, val int) {
	delta := val - f.nums[index]
	f.nums[index] = val
	f.add(index+1, delta)
}

func (f *FenwickTree) SumRange(left, right int) int {
	return f.prefix(right+1) - f.prefix(left)
}
```

### Dry Run

Build from `[1,3,5]` (1-indexed tree, n=3). After seeding: `tree[1]=1`, `tree[2]=1+3=4`, `tree[3]=5`.

| Op | Steps | Output |
|----|-------|--------|
| `sumRange(0,2)` | prefix(3) = tree[3]+tree[2] = 5+4 = 9; prefix(0) = 0; 9−0 | **9** |
| `update(1,2)` | delta = 2−3 = −1; add(2,−1): tree[2]=3 | null |
| `sumRange(0,2)` | prefix(3) = tree[3]+tree[2] = 5+3 = 8; prefix(0)=0; 8−0 | **8** |

---

## Approach 3 — Segment Tree

### Intuition

A segment tree stores range sums in a binary tree flattened into an array of size `2n`: leaves `n..2n−1` hold the elements and each internal node holds the sum of its two children. A point update walks leaf → root fixing O(log n) parents; a range query splits `[l, r]` into O(log n) canonical segments. Prefer this when the query later generalizes to min/max or lazy range updates.

### Algorithm

1. Build: place values at `[n, 2n)`; for `i = n−1` down to `1`, `tree[i] = tree[2i] + tree[2i+1]`.
2. `Update(index, val)`: set leaf `tree[index+n] = val`; climb halving the index and resumming parents.
3. `SumRange(left, right)`: `l = left+n`, `r = right+n`; while `l <= r`, take `tree[l]` when `l` is a right child, take `tree[r]` when `r` is a left child, then halve both.

### Complexity

- **Time:** Update O(log n), SumRange O(log n); Build O(n).
- **Space:** O(n) — the `2n` array.

### Code

```go
type SegmentTree struct {
	n    int
	tree []int
}

func NewSegmentTree(nums []int) *SegmentTree {
	n := len(nums)
	t := make([]int, 2*n)
	for i := 0; i < n; i++ {
		t[n+i] = nums[i]
	}
	for i := n - 1; i > 0; i-- {
		t[i] = t[2*i] + t[2*i+1]
	}
	return &SegmentTree{n: n, tree: t}
}

func (s *SegmentTree) Update(index, val int) {
	i := index + s.n
	s.tree[i] = val
	for i > 1 {
		i /= 2
		s.tree[i] = s.tree[2*i] + s.tree[2*i+1]
	}
}

func (s *SegmentTree) SumRange(left, right int) int {
	l, r := left+s.n, right+s.n
	sum := 0
	for l <= r {
		if l%2 == 1 {
			sum += s.tree[l]
			l++
		}
		if r%2 == 0 {
			sum += s.tree[r]
			r--
		}
		l /= 2
		r /= 2
	}
	return sum
}
```

### Dry Run

Build from `[1,3,5]` (n=3, array size 6). Leaves: `tree[3]=1, tree[4]=3, tree[5]=5`. Internals: `tree[2]=tree[4]+tree[5]=8`, `tree[1]=tree[2]+tree[3]=9`.

| Op | Steps | Output |
|----|-------|--------|
| `sumRange(0,2)` | l=3, r=5: l odd → +tree[3]=1, l=4; r odd (not even) skip; l=2,r=2: r even → +tree[2]=8, r=1; l=1,r=0 stop | **9** |
| `update(1,2)` | leaf tree[4]=2; tree[2]=tree[4]+tree[5]=7; tree[1]=tree[2]+tree[3]=8 | null |
| `sumRange(0,2)` | l=3,r=5: +tree[3]=1,l=4; l=2,r=2: +tree[2]=7,r=1; stop | **8** |

---

## Key Takeaways

- **Range sum + point update = Fenwick or segment tree.** The BIT is the shortest to code and hardest to get wrong; the segment tree is the more general hammer.
- **`lowbit(i) = i & -i`** isolates the lowest set bit — the whole Fenwick tree is built on climbing (`+= lowbit`) and descending (`−= lowbit`) by it.
- **Set vs. add:** when the API sets a value, cache the current values and feed the *delta* into an additive structure.
- **Iterative (non-recursive) segment tree** with leaves at `[n, 2n)` avoids recursion overhead and is easy to memorize.

---

## Related Problems

- LeetCode #303 — Range Sum Query - Immutable (prefix sums, no updates)
- LeetCode #308 — Range Sum Query 2D - Mutable (2D Fenwick tree)
- LeetCode #315 — Count of Smaller Numbers After Self (Fenwick tree on ranks)
- LeetCode #218 — The Skyline Problem (segment tree / heap)
