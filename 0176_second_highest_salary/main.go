package main

import (
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
// into a slice of rows and implementing the query logic by hand.
type Employee struct {
	ID     int
	Salary int
}

// formatNullable renders a *int the way LeetCode renders a SQL result cell:
// the number itself, or "null" when the value does not exist.
func formatNullable(v *int) string {
	if v == nil {
		return "null" // SQL NULL — no second highest distinct salary exists
	}
	return strconv.Itoa(*v)
}

// ── Approach 1: Brute Force (Sort Distinct Salaries) ─────────────────────────
//
// bruteForce solves Second Highest Salary by materialising every distinct
// salary, sorting descending, and picking index 1.
//
// Intuition:
//
//	"Second highest DISTINCT salary" is literally "element at offset 1 of the
//	distinct salaries sorted descending". Deduplicate, sort, index — the exact
//	plan of the SQL `SELECT DISTINCT salary ORDER BY salary DESC LIMIT 1
//	OFFSET 1`.
//
// Algorithm:
//  1. Walk the table, inserting each salary into a hash set; collect the
//     first occurrence of every value into a slice.
//  2. Sort the distinct slice in descending order.
//  3. If fewer than 2 distinct salaries exist, return nil (SQL NULL);
//     otherwise return a pointer to the element at index 1.
//
// Time:  O(n log n) — the sort dominates the O(n) dedup pass.
// Space: O(n) — the hash set plus the distinct-values slice.
func bruteForce(employees []Employee) *int {
	seen := map[int]bool{} // salary → already collected?
	distinct := []int{}    // unique salaries in first-seen order
	for _, e := range employees {
		if !seen[e.Salary] {
			seen[e.Salary] = true // mark so duplicates are skipped
			distinct = append(distinct, e.Salary)
		}
	}
	// Descending sort puts the highest at index 0, second highest at index 1.
	sort.Sort(sort.Reverse(sort.IntSlice(distinct)))
	if len(distinct) < 2 {
		return nil // fewer than two distinct salaries → SQL NULL
	}
	return &distinct[1] // the second highest distinct salary
}

// ── Approach 2: Two-Pass Max ─────────────────────────────────────────────────
//
// twoPassMax solves Second Highest Salary by finding the maximum, then the
// maximum among salaries strictly below it.
//
// Intuition:
//
//	The second highest distinct salary is exactly MAX(salary) restricted to
//	salaries < MAX(salary). Two linear scans, no sorting, no extra memory —
//	the Go translation of `SELECT MAX(salary) WHERE salary < (SELECT
//	MAX(salary) ...)`. Duplicates of the maximum are excluded automatically
//	by the strict `<`.
//
// Algorithm:
//  1. Pass 1: scan all rows to find the highest salary `max1`.
//  2. Pass 2: scan again keeping the largest salary strictly less than
//     `max1` (track whether any such salary was found).
//  3. If pass 2 found nothing, every row equals `max1` → return nil (NULL);
//     otherwise return the runner-up value.
//
// Time:  O(n) — two full scans of the table.
// Space: O(1) — two scalars and a found flag.
func twoPassMax(employees []Employee) *int {
	if len(employees) == 0 {
		return nil // empty table → no salaries at all
	}
	// Pass 1: global maximum salary.
	max1 := employees[0].Salary
	for _, e := range employees {
		if e.Salary > max1 {
			max1 = e.Salary // found a higher salary — update the maximum
		}
	}
	// Pass 2: maximum among salaries strictly below max1.
	best := 0      // candidate runner-up value (valid only when found == true)
	found := false // whether any salary < max1 exists
	for _, e := range employees {
		// Strict < skips every duplicate of the maximum (distinct semantics).
		if e.Salary < max1 && (!found || e.Salary > best) {
			best, found = e.Salary, true
		}
	}
	if !found {
		return nil // all rows share one distinct salary → SQL NULL
	}
	return &best
}

// ── Approach 3: Single Pass Top-2 Tracking (Optimal) ─────────────────────────
//
// singlePass solves Second Highest Salary in one scan by maintaining the two
// largest distinct salaries seen so far.
//
// Intuition:
//
//	Keep two running slots: `first` (highest so far) and `second` (highest
//	strictly below `first`). A new salary either dethrones `first` (pushing
//	the old `first` down to `second`), slots between them, or is ignored —
//	duplicates of `first` fall through both conditions, preserving the
//	DISTINCT requirement. One pass, O(1) space: the streaming version of
//	Approach 2, and the same pattern as #414 Third Maximum Number.
//
// Algorithm:
//  1. Initialise first = second = -∞ (math.MinInt as "empty" sentinel).
//  2. For each salary s:
//     a. If s > first: second = first, first = s.
//     b. Else if s < first and s > second: second = s (strict < skips
//     duplicates of the current maximum).
//  3. If second is still the sentinel, no second distinct salary exists →
//     return nil; otherwise return second.
//
// Time:  O(n) — exactly one scan of the table.
// Space: O(1) — two integer slots.
func singlePass(employees []Employee) *int {
	first, second := math.MinInt, math.MinInt // -∞ sentinels mean "slot empty"
	for _, e := range employees {
		s := e.Salary
		switch {
		case s > first:
			second = first // old champion becomes the runner-up
			first = s      // new distinct maximum
		case s < first && s > second:
			second = s // fits strictly between the two slots
		}
		// s == first (duplicate of the max) matches neither case → ignored,
		// which is exactly the DISTINCT semantics the problem demands.
	}
	if second == math.MinInt {
		return nil // never filled → fewer than two distinct salaries
	}
	return &second
}

func main() {
	// Example 1: Employee = [(1,100), (2,200), (3,300)] → 200
	example1 := []Employee{{ID: 1, Salary: 100}, {ID: 2, Salary: 200}, {ID: 3, Salary: 300}}
	// Example 2: Employee = [(1,100)] → null
	example2 := []Employee{{ID: 1, Salary: 100}}
	// Edge: duplicated maximum — distinct semantics must return 100, not 200.
	edge := []Employee{{ID: 1, Salary: 200}, {ID: 2, Salary: 200}, {ID: 3, Salary: 100}}

	fmt.Println("=== Approach 1: Brute Force (Sort Distinct Salaries) ===")
	fmt.Printf("Example 1: salaries=[100 200 300]  got=%s  expected 200\n", formatNullable(bruteForce(example1)))
	fmt.Printf("Example 2: salaries=[100]          got=%s  expected null\n", formatNullable(bruteForce(example2)))
	fmt.Printf("Edge:      salaries=[200 200 100]  got=%s  expected 100\n", formatNullable(bruteForce(edge)))

	fmt.Println("=== Approach 2: Two-Pass Max ===")
	fmt.Printf("Example 1: salaries=[100 200 300]  got=%s  expected 200\n", formatNullable(twoPassMax(example1)))
	fmt.Printf("Example 2: salaries=[100]          got=%s  expected null\n", formatNullable(twoPassMax(example2)))
	fmt.Printf("Edge:      salaries=[200 200 100]  got=%s  expected 100\n", formatNullable(twoPassMax(edge)))

	fmt.Println("=== Approach 3: Single Pass Top-2 Tracking (Optimal) ===")
	fmt.Printf("Example 1: salaries=[100 200 300]  got=%s  expected 200\n", formatNullable(singlePass(example1)))
	fmt.Printf("Example 2: salaries=[100]          got=%s  expected null\n", formatNullable(singlePass(example2)))
	fmt.Printf("Edge:      salaries=[200 200 100]  got=%s  expected 100\n", formatNullable(singlePass(edge)))
}
