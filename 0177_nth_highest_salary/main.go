package main

import (
	"container/heap"
	"fmt"
	"math"
	"sort"
	"strconv"
)

// Employee models one row of the LeetCode `Employee` table:
//
//	+-------------+------+
//	| Column Name | Type |
//	+-------------+------+
//	| id          | int  |
//	| salary      | int  |
//	+-------------+------+
//
// This is a Database problem; the repo solves it in Go by loading the table
// into a slice of rows and implementing getNthHighestSalary(n) by hand.
type Employee struct {
	ID     int
	Salary int
}

// formatNullable renders a *int the way LeetCode renders a SQL result cell:
// the number itself, or "null" when the value does not exist.
func formatNullable(v *int) string {
	if v == nil {
		return "null" // SQL NULL — fewer than n distinct salaries exist
	}
	return strconv.Itoa(*v)
}

// distinctSalaries extracts the unique salary values from the table —
// the Go equivalent of SQL's DISTINCT. Shared by approaches 2–4.
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

// ── Approach 1: Brute Force (Repeated Max Stripping) ─────────────────────────
//
// bruteForce solves Nth Highest Salary by peeling off the maximum n times.
//
// Intuition:
//
//	The 1st highest distinct salary is MAX(salary). The 2nd is the max among
//	salaries strictly below that. Repeat "find the max strictly below the
//	previous ceiling" n times and the last value peeled is the n-th highest
//	distinct salary — duplicates never bother us because the ceiling drops
//	strictly each round. No sorting, no extra memory, n full scans.
//
// Algorithm:
//  1. If n < 1, return nil (rank is undefined).
//  2. ceiling = +∞. Repeat n times:
//     a. Scan all rows for the largest salary strictly below ceiling.
//     b. If none exists, there are fewer than n distinct salaries → nil.
//     c. Otherwise set ceiling to that value.
//  3. After n rounds, ceiling holds the n-th highest distinct salary.
//
// Time:  O(n·m) — n scans over m rows.
// Space: O(1) — only the ceiling and a found flag.
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

// ── Approach 2: Sort Distinct Salaries ───────────────────────────────────────
//
// sortDistinct solves Nth Highest Salary by sorting the distinct salaries
// descending and indexing position n-1.
//
// Intuition:
//
//	"n-th highest DISTINCT salary" is the element at offset n-1 of the
//	deduplicated salaries in descending order — a direct translation of the
//	SQL `SELECT DISTINCT salary ORDER BY salary DESC LIMIT 1 OFFSET n-1`.
//
// Algorithm:
//  1. If n < 1, return nil.
//  2. Deduplicate the salaries with a hash set.
//  3. Sort the distinct values in descending order.
//  4. If fewer than n distinct values exist, return nil; otherwise return
//     the element at index n-1.
//
// Time:  O(m log m) — sorting the (≤ m) distinct salaries dominates.
// Space: O(m) — the hash set plus the distinct slice.
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

// ── Approach 3: Min-Heap of Size n ───────────────────────────────────────────
//
// minHeapTopN solves Nth Highest Salary by streaming the distinct salaries
// through a min-heap capped at n elements.
//
// Intuition:
//
//	Keep only the n largest distinct salaries seen so far in a min-heap.
//	The heap root is the smallest of those n — i.e. the current candidate
//	for "n-th highest". Push each distinct salary; whenever the heap grows
//	past n, pop the root (something too small to matter). After the stream,
//	a full heap's root IS the n-th highest. Classic top-k: great when m is
//	huge and n is small, because memory stays O(n).
//
// Algorithm:
//  1. If n < 1, return nil.
//  2. Deduplicate the salaries (DISTINCT semantics).
//  3. Push each distinct salary onto a min-heap; if the size exceeds n,
//     pop the minimum.
//  4. If the heap holds fewer than n values, return nil; otherwise the
//     root is the answer.
//
// Time:  O(m log n) — m pushes/pops against a heap capped at n.
// Space: O(m) — the dedup set (the heap itself is only O(n)).
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

// ── Approach 4: Quickselect (Optimal Average) ────────────────────────────────
//
// quickSelect solves Nth Highest Salary with Hoare's selection algorithm on
// the distinct salaries — average linear time, no full sort.
//
// Intuition:
//
//	Sorting all distinct salaries does more work than needed: we want one
//	order statistic, not the whole order. Quickselect partitions around a
//	pivot (as quicksort does) but then recurses into only the side that
//	contains the target index, discarding the other side entirely. The
//	n-th highest of d distinct values is the (d-n)-th smallest (0-indexed),
//	so one selection call answers the query in O(d) expected time.
//
// Algorithm:
//  1. If n < 1, return nil. Deduplicate salaries into `vals` (d values).
//  2. If d < n, return nil.
//  3. target = d - n (index of the answer in ascending order).
//  4. Loop: pick the middle element as pivot (guards sorted inputs), move
//     it to the end, Lomuto-partition [lo..hi]; the pivot lands at index p
//     with smaller values left, larger-or-equal right.
//  5. If p == target, vals[p] is the answer; if p < target, recurse right
//     (lo = p+1); else recurse left (hi = p-1).
//
// Time:  O(m) average — partitions shrink geometrically (O(m²) worst case).
// Space: O(m) — the distinct slice; the selection itself is in-place, O(1).
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

func main() {
	// Example 1: Employee = [(1,100), (2,200), (3,300)], n = 2 → 200
	example1 := []Employee{{ID: 1, Salary: 100}, {ID: 2, Salary: 200}, {ID: 3, Salary: 300}}
	// Example 2: Employee = [(1,100)], n = 2 → null
	example2 := []Employee{{ID: 1, Salary: 100}}
	// Edge: duplicates + n beyond the distinct count.
	edge := []Employee{{ID: 1, Salary: 300}, {ID: 2, Salary: 300}, {ID: 3, Salary: 100}}

	fmt.Println("=== Approach 1: Brute Force (Repeated Max Stripping) ===")
	fmt.Printf("Example 1: salaries=[100 200 300] n=2  got=%s  expected 200\n", formatNullable(bruteForce(example1, 2)))
	fmt.Printf("Example 2: salaries=[100]         n=2  got=%s  expected null\n", formatNullable(bruteForce(example2, 2)))
	fmt.Printf("Edge:      salaries=[300 300 100] n=2  got=%s  expected 100\n", formatNullable(bruteForce(edge, 2)))
	fmt.Printf("Edge:      salaries=[300 300 100] n=3  got=%s  expected null\n", formatNullable(bruteForce(edge, 3)))

	fmt.Println("=== Approach 2: Sort Distinct Salaries ===")
	fmt.Printf("Example 1: salaries=[100 200 300] n=2  got=%s  expected 200\n", formatNullable(sortDistinct(example1, 2)))
	fmt.Printf("Example 2: salaries=[100]         n=2  got=%s  expected null\n", formatNullable(sortDistinct(example2, 2)))
	fmt.Printf("Edge:      salaries=[300 300 100] n=2  got=%s  expected 100\n", formatNullable(sortDistinct(edge, 2)))
	fmt.Printf("Edge:      salaries=[300 300 100] n=3  got=%s  expected null\n", formatNullable(sortDistinct(edge, 3)))

	fmt.Println("=== Approach 3: Min-Heap of Size n ===")
	fmt.Printf("Example 1: salaries=[100 200 300] n=2  got=%s  expected 200\n", formatNullable(minHeapTopN(example1, 2)))
	fmt.Printf("Example 2: salaries=[100]         n=2  got=%s  expected null\n", formatNullable(minHeapTopN(example2, 2)))
	fmt.Printf("Edge:      salaries=[300 300 100] n=2  got=%s  expected 100\n", formatNullable(minHeapTopN(edge, 2)))
	fmt.Printf("Edge:      salaries=[300 300 100] n=3  got=%s  expected null\n", formatNullable(minHeapTopN(edge, 3)))

	fmt.Println("=== Approach 4: Quickselect (Optimal Average) ===")
	fmt.Printf("Example 1: salaries=[100 200 300] n=2  got=%s  expected 200\n", formatNullable(quickSelect(example1, 2)))
	fmt.Printf("Example 2: salaries=[100]         n=2  got=%s  expected null\n", formatNullable(quickSelect(example2, 2)))
	fmt.Printf("Edge:      salaries=[300 300 100] n=2  got=%s  expected 100\n", formatNullable(quickSelect(edge, 2)))
	fmt.Printf("Edge:      salaries=[300 300 100] n=3  got=%s  expected null\n", formatNullable(quickSelect(edge, 3)))
}
