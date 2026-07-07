# 0177 — Nth Highest Salary

> LeetCode #177 · Difficulty: Medium
> **Categories:** Database, Sorting, Heap, Quickselect

---

## Problem Statement

Table: `Employee`

```
+-------------+------+
| Column Name | Type |
+-------------+------+
| id          | int  |
| salary      | int  |
+-------------+------+
```

`id` is the primary key (column with unique values) for this table.
Each row of this table contains information about the salary of an employee.

Write a solution to find the **nth highest distinct salary** from the `Employee` table. If there are less than `n` distinct salaries, return `null`.

The result format is in the following example.

**Example 1:**

```
Input:
Employee table:
+----+--------+
| id | salary |
+----+--------+
| 1  | 100    |
| 2  | 200    |
| 3  | 300    |
+----+--------+
n = 2
Output:
+------------------------+
| getNthHighestSalary(2) |
+------------------------+
| 200                    |
+------------------------+
```

**Example 2:**

```
Input:
Employee table:
+----+--------+
| id | salary |
+----+--------+
| 1  | 100    |
+----+--------+
n = 2
Output:
+------------------------+
| getNthHighestSalary(2) |
+------------------------+
| null                   |
+------------------------+
```

> **Note:** This is a Database problem (on LeetCode you write a SQL function
> `getNthHighestSalary(N INT)`). Per this repo's convention it is solved in
> **Go**: the table becomes `[]Employee{ID, Salary}`, each approach implements
> `getNthHighestSalary(n)` by hand, and `*int` models SQL `NULL` via `nil`.
> The canonical SQL appears in Key Takeaways.

---

## Company Frequency

| Company    | Frequency        | Last Reported |
|------------|------------------|---------------|
| Amazon     | ★★★★☆ High       | 2024          |
| Microsoft  | ★★★★☆ High       | 2024          |
| Google     | ★★★☆☆ Medium     | 2023          |
| Adobe      | ★★★☆☆ Medium     | 2023          |
| Uber       | ★★☆☆☆ Low        | 2023          |
| Flipkart   | ★★☆☆☆ Low        | 2023          |

> ⚠️ Frequency data is crowd-sourced from LeetCode Discuss, Glassdoor, and
> community interview reports. Treat as a signal, not a guarantee.

---

## DSA Concepts Used

- **Hash Map / Hash Set** — deduplication implements SQL `DISTINCT`; every approach must rank *distinct* values → see [`/dsa/hash_map.md`](/dsa/hash_map.md)
- **Sorting** — descending sort turns "n-th highest" into plain indexing (`LIMIT 1 OFFSET n-1`) → see [`/dsa/sorting.md`](/dsa/sorting.md)
- **Heap / Priority Queue** — a min-heap capped at n keeps exactly the top-n distinct values in a stream; its root is the n-th highest → see [`/dsa/heap_priority_queue.md`](/dsa/heap_priority_queue.md)
- **Quickselect** — one order statistic doesn't need a full sort; partition-and-discard finds it in average O(m) → see [`/dsa/quickselect.md`](/dsa/quickselect.md)

---

## Approaches Overview

| # | Approach | Time | Space | When to use |
|---|----------|------|-------|-------------|
| 1 | Brute Force (Repeated Max Stripping) | O(n·m) | O(1) | Tiny n; zero extra memory; no sort available |
| 2 | Sort Distinct Salaries | O(m log m) | O(m) | Default readable answer; mirrors the SQL plan exactly |
| 3 | Min-Heap of Size n | O(m log n) | O(m) set + O(n) heap | Streaming data / huge m with small n |
| 4 | Quickselect (Optimal Average) | O(m) avg, O(m²) worst | O(m) | One-shot query on in-memory data; the asymptotic winner |

*(m = number of rows, n = requested rank, d = number of distinct salaries ≤ m.)*

---

## Approach 1 — Brute Force (Repeated Max Stripping)

### Intuition

The 1st highest distinct salary is `MAX(salary)`. The 2nd highest is the maximum among salaries **strictly below** that. Iterate the idea: keep a falling `ceiling`, and n times find the largest salary strictly under it. Because the ceiling drops strictly each round, duplicates of already-peeled values can never win again — DISTINCT semantics fall out of the strict `<` for free. If some round finds nothing, there are fewer than n distinct salaries and the answer is `NULL`.

### Algorithm

1. If `n < 1`, return `nil` (the rank is undefined).
2. Set `ceiling = +∞` (`math.MaxInt`).
3. Repeat n times:
   1. Scan all rows for the largest salary strictly below `ceiling` (track a `found` flag).
   2. If none exists, return `nil`.
   3. Otherwise assign it to `ceiling`.
4. After n rounds, `ceiling` is the n-th highest distinct salary — return it.

### Complexity

- **Time:** O(n·m) — n full scans over the m rows.
- **Space:** O(1) — one ceiling, one candidate, one flag.

### Code

```go
func bruteForce(employees []Employee, n int) *int {
	if n < 1 {
		return nil // "n-th highest" is undefined for n < 1
	}
	ceiling := math.MaxInt // +∞ sentinel: round 1 accepts every salary
	for round := 0; round < n; round++ {
		best := 0      // largest salary strictly below the ceiling this round
		found := false // whether any such salary exists
		for _, e := range employees {
			// Strict < keeps duplicates of previous winners out (DISTINCT).
			if e.Salary < ceiling && (!found || e.Salary > best) {
				best, found = e.Salary, true
			}
		}
		if !found {
			return nil // ran out of distinct salaries before reaching rank n
		}
		ceiling = best // lower the ceiling and peel the next rank
	}
	return &ceiling // the value peeled on round n
}
```

### Dry Run

Example 1: `salaries = [100, 200, 300]`, `n = 2`

| Step | round | ceiling before | scan result (max < ceiling) | found | ceiling after |
|------|-------|----------------|------------------------------|-------|---------------|
| 1 | 0 | +∞ | 300 (all qualify; 300 largest) | true | 300 |
| 2 | 1 | 300 | 200 (100 and 200 qualify) | true | 200 |
| 3 | — | loop done (2 rounds) | — | — | return **200** |

Result: `200` ✔ (Example 2: round 0 peels 100, round 1 finds nothing < 100 → `null` ✔)

---

## Approach 2 — Sort Distinct Salaries

### Intuition

"n-th highest DISTINCT salary" is the element at offset n−1 of the deduplicated salaries in descending order — the literal Go translation of the SQL solution `SELECT DISTINCT salary FROM Employee ORDER BY salary DESC LIMIT 1 OFFSET n-1`. Deduplicate with a hash set, sort descending, bounds-check, index.

### Algorithm

1. If `n < 1`, return `nil`.
2. Build the distinct salary slice via a hash set (`DISTINCT`).
3. Sort it descending (`ORDER BY salary DESC`).
4. If fewer than n distinct values exist, return `nil` (SQL returns `NULL`).
5. Return the element at index `n-1` (`LIMIT 1 OFFSET n-1`).

### Complexity

- **Time:** O(m log m) — O(m) dedup + sorting d ≤ m distinct values.
- **Space:** O(m) — the hash set and the distinct slice.

### Code

```go
func sortDistinct(employees []Employee, n int) *int {
	if n < 1 {
		return nil // "n-th highest" is undefined for n < 1
	}
	distinct := distinctSalaries(employees)
	// Descending order: rank 1 sits at index 0, rank n at index n-1.
	sort.Sort(sort.Reverse(sort.IntSlice(distinct)))
	if len(distinct) < n {
		return nil // fewer than n distinct salaries → SQL NULL
	}
	return &distinct[n-1]
}

// distinctSalaries extracts the unique salary values from the table —
// the Go equivalent of SQL's DISTINCT.
func distinctSalaries(employees []Employee) []int {
	seen := map[int]bool{} // salary → already collected?
	distinct := []int{}    // unique salaries in first-seen order
	for _, e := range employees {
		if !seen[e.Salary] {
			seen[e.Salary] = true // remember so duplicates are skipped
			distinct = append(distinct, e.Salary)
		}
	}
	return distinct
}
```

### Dry Run

Example 1: `salaries = [100, 200, 300]`, `n = 2`

| Step | Action | distinct | Check | Result |
|------|--------|----------|-------|--------|
| 1 | dedup 100, 200, 300 | [100, 200, 300] | — | — |
| 2 | sort descending | [300, 200, 100] | — | — |
| 3 | bounds check | — | 3 ≥ 2 | index n−1 = 1 valid |
| 4 | return `distinct[1]` | — | — | **200** |

Result: `200` ✔ (Example 2: distinct = [100], 1 < 2 → `null` ✔)

---

## Approach 3 — Min-Heap of Size n

### Intuition

To know the n-th highest you only ever need to remember the **n largest distinct values seen so far** — and the interesting one among them is the *smallest of those n*, which is exactly what a min-heap root gives in O(1). Stream the distinct salaries through a min-heap capped at n: push each value, and when the size exceeds n pop the root (a value provably not in the top n). When the stream ends, a full heap's root **is** the n-th highest distinct salary. This is the classic top-k pattern: ideal when m is enormous and n is small, since working memory beyond the dedup set is only O(n).

### Algorithm

1. If `n < 1`, return `nil`.
2. Deduplicate salaries (DISTINCT semantics).
3. For each distinct salary: `heap.Push`; if `Len() > n`, `heap.Pop` the minimum.
4. If the heap holds fewer than n values, return `nil`.
5. Otherwise return the root `(*h)[0]`.

### Complexity

- **Time:** O(m log n) — each of ≤ m distinct values costs one push and at most one pop against a heap of size ≤ n+1.
- **Space:** O(m) for the dedup set; the heap itself is only O(n).

### Code

```go
func minHeapTopN(employees []Employee, n int) *int {
	if n < 1 {
		return nil // "n-th highest" is undefined for n < 1
	}
	distinct := distinctSalaries(employees)
	h := &intMinHeap{}
	heap.Init(h)
	for _, s := range distinct {
		heap.Push(h, s) // consider this distinct salary for the top n
		if h.Len() > n {
			heap.Pop(h) // evict the smallest — it can't be in the top n
		}
	}
	if h.Len() < n {
		return nil // fewer than n distinct salaries → SQL NULL
	}
	root := (*h)[0] // smallest of the n largest = the n-th highest
	return &root
}

// intMinHeap is a min-heap of ints for container/heap — the root is always
// the smallest element, so it can evict the smallest of the "top n so far".
type intMinHeap []int

func (h intMinHeap) Len() int           { return len(h) }
func (h intMinHeap) Less(i, j int) bool { return h[i] < h[j] }
func (h intMinHeap) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }
func (h *intMinHeap) Push(x any)        { *h = append(*h, x.(int)) }
func (h *intMinHeap) Pop() any {
	old := *h
	last := len(old) - 1
	v := old[last]  // heap.Pop has already swapped the min to the end
	*h = old[:last] // shrink the slice
	return v
}
```

### Dry Run

Example 1: `salaries = [100, 200, 300]`, `n = 2` (heap shown smallest-first)

| Step | s | heap after push | Len > 2? | pop | heap after |
|------|-----|-----------------|----------|-----|------------|
| 1 | 100 | [100] | no | — | [100] |
| 2 | 200 | [100, 200] | no | — | [100, 200] |
| 3 | 300 | [100, 200, 300] | yes | 100 | [200, 300] |
| 4 | — | stream done, Len = 2 = n | — | — | root = **200** |

Result: `200` ✔ (Example 2: heap ends as [100], Len 1 < 2 → `null` ✔)

---

## Approach 4 — Quickselect (Optimal Average)

### Intuition

Sorting computes *every* rank; we need exactly *one*. Quickselect partitions the distinct values around a pivot (like quicksort) but then recurses into only the side containing the target index and discards the other side wholesale — the work per level shrinks geometrically, giving average linear time. Rank conversion: the n-th largest of d distinct values sits at ascending index `d − n`. A middle-element pivot avoids the notorious O(d²) degeneration on already-sorted inputs.

### Algorithm

1. If `n < 1`, return `nil`. Deduplicate salaries into `vals` (d values).
2. If `d < n`, return `nil`.
3. `target = d − n` (the answer's index in ascending order).
4. Loop on window `[lo, hi]`:
   1. Swap the middle element to `hi` as pivot; Lomuto-partition: values `< pivot` shift left, pivot lands at index `p`.
   2. If `p == target`, return `vals[p]`.
   3. If `p < target`, set `lo = p + 1` (answer is in the larger side).
   4. Else set `hi = p − 1` (answer is in the smaller side).

### Complexity

- **Time:** O(m) average — d + d/2 + d/4 + … ≈ 2d comparisons after the O(m) dedup; O(m²) worst case with adversarial pivots.
- **Space:** O(m) — the distinct slice; the selection itself is in-place, O(1) auxiliary.

### Code

```go
func quickSelect(employees []Employee, n int) *int {
	if n < 1 {
		return nil // "n-th highest" is undefined for n < 1
	}
	vals := distinctSalaries(employees)
	if len(vals) < n {
		return nil // fewer than n distinct salaries → SQL NULL
	}
	// n-th largest (1-indexed) == element at ascending index d-n (0-indexed).
	target := len(vals) - n
	lo, hi := 0, len(vals)-1
	for {
		// Middle-element pivot avoids the O(d²) trap on already-sorted data.
		mid := lo + (hi-lo)/2
		vals[mid], vals[hi] = vals[hi], vals[mid] // stash pivot at the end
		p := partition(vals, lo, hi)
		switch {
		case p == target:
			return &vals[p] // pivot landed exactly on the target rank
		case p < target:
			lo = p + 1 // answer lies in the right (larger) side
		default:
			hi = p - 1 // answer lies in the left (smaller) side
		}
	}
}

// partition performs a Lomuto partition of vals[lo..hi] around the pivot
// stored at vals[hi]; returns the pivot's final resting index.
func partition(vals []int, lo, hi int) int {
	pivot := vals[hi] // pivot value (previously swapped to the end)
	i := lo           // boundary: vals[lo..i-1] are all < pivot
	for j := lo; j < hi; j++ {
		if vals[j] < pivot {
			vals[i], vals[j] = vals[j], vals[i] // grow the "< pivot" prefix
			i++
		}
	}
	vals[i], vals[hi] = vals[hi], vals[i] // drop the pivot into its slot
	return i
}
```

### Dry Run

Example 1: `salaries = [100, 200, 300]`, `n = 2` → `vals = [100, 200, 300]`, `d = 3`, `target = 3 − 2 = 1`

| Step | lo | hi | mid | pivot | vals after partition | p | p vs target | Action |
|------|----|----|-----|-------|----------------------|---|-------------|--------|
| 1 | 0 | 2 | 1 | 200 (swapped to hi: [100, 300, 200]) | [100, 200, 300] | 1 | p == 1 == target | return `vals[1]` |

Partition detail (step 1): pivot = 200; j=0: 100 < 200 → keep left, i→1; j=1: 300 ≥ 200 → stays; final swap puts 200 at index 1.

Result: `200` ✔ (Example 2: d = 1 < n = 2 → `null` before any partitioning ✔)

---

## Key Takeaways

- **One rank ≠ full sort.** The ladder of tools for "k-th largest": repeated max O(k·m) → sort O(m log m) → heap O(m log k) → quickselect O(m) average. Pick by the relative sizes of k and m and by whether data streams.
- **DISTINCT first, rank second** — dedup with a hash set (or via strict `<` ceilings) before any ranking logic; forgetting it is *the* classic bug in this family (`[300,300,100]`, n=2 must give 100, not 300).
- **Rank conversion formula:** n-th largest of d values = ascending index `d − n`; off-by-ones here are silent and deadly.
- **Guard degenerate n** (`n ≤ 0`) and `d < n` explicitly — in SQL a negative `OFFSET N-1` is a runtime error unless handled (`SET N = N - 1` first).
- Canonical SQL (MySQL function form):
  ```sql
  CREATE FUNCTION getNthHighestSalary(N INT) RETURNS INT
  BEGIN
    SET N = N - 1;
    RETURN (
      SELECT DISTINCT salary FROM Employee
      ORDER BY salary DESC LIMIT 1 OFFSET N
    );
  END
  ```

---

## Related Problems

- LeetCode #176 — Second Highest Salary (this problem with n fixed to 2)
- LeetCode #185 — Department Top Three Salaries (top-k distinct per group)
- LeetCode #215 — Kth Largest Element in an Array (same quickselect/heap ladder, no distinct twist)
- LeetCode #703 — Kth Largest Element in a Stream (the size-k min-heap as a live data structure)
- LeetCode #414 — Third Maximum Number (O(1)-space distinct top-3 tracking)
