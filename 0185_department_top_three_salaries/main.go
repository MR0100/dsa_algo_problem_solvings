package main

import (
	"container/heap"
	"fmt"
	"math"
	"sort"
)

// LeetCode 185 — Department Top Three Salaries.
//
// The original problem is a SQL one: given the tables
//
//	Employee: +----+------+--------+--------------+   Department: +----+------+
//	          | id | name | salary | departmentId |               | id | name |
//	          +----+------+--------+--------------+               +----+------+
//
// a "high earner" in a department is an employee whose salary is among the
// TOP THREE DISTINCT salaries of that department. Report every high earner as
// a row (Department, Employee, Salary), in any order. Here the tables are
// modelled as slices of rows, and every approach re-implements the dense-rank
// <= 3 query as an in-memory top-K-distinct strategy.

// Employee models one row of the Employee table.
// DepartmentID is a foreign key referencing Department.ID.
type Employee struct {
	ID           int
	Name         string
	Salary       int
	DepartmentID int
}

// Department models one row of the Department table.
type Department struct {
	ID   int
	Name string
}

// ResultRow is one row of the result table (Department, Employee, Salary).
type ResultRow struct {
	Department string
	Employee   string
	Salary     int
}

// nameMap builds the departmentId → name lookup shared by all approaches.
func nameMap(departments []Department) map[int]string {
	m := make(map[int]string, len(departments))
	for _, d := range departments {
		m[d.ID] = d.Name // ids are unique — one name per department
	}
	return m
}

// ── Approach 1: Brute Force (Correlated Distinct Count) ──────────────────────
//
// bruteForce solves Department Top Three Salaries by counting, for every
// employee, the distinct salaries in their department that are strictly
// higher.
//
// Intuition:
//
//	"Salary is among the top three distinct salaries" is equivalent to
//	"fewer than three DISTINCT salaries in the department are strictly
//	higher". That is exactly the classic correlated SQL subquery
//	`3 > COUNT(DISTINCT e2.salary WHERE e2.salary > e1.salary)` — evaluate it
//	literally with a nested scan per employee.
//
// Algorithm:
//  1. For each employee e, scan all employees of the same department and
//     collect the DISTINCT salaries strictly greater than e.Salary in a set.
//  2. If the set holds fewer than 3 salaries, e is a high earner: look up the
//     department name and emit (deptName, e.Name, e.Salary).
//
// Time:  O(n^2) — a full O(n) scan (with O(1) set inserts) per employee.
// Space: O(n) — the per-employee set of higher distinct salaries (reused
//        scratch; the department-name map adds O(d)).
func bruteForce(employees []Employee, departments []Department) []ResultRow {
	nameByDept := nameMap(departments)

	result := []ResultRow{}
	for _, e := range employees {
		higher := map[int]bool{} // DISTINCT salaries in e's dept above e.Salary
		for _, other := range employees {
			if other.DepartmentID == e.DepartmentID && other.Salary > e.Salary {
				higher[other.Salary] = true // set collapses duplicates for free
			}
		}
		if len(higher) < 3 { // fewer than 3 distinct salaries beat e → top three
			result = append(result, ResultRow{nameByDept[e.DepartmentID], e.Name, e.Salary})
		}
	}
	return result
}

// ── Approach 2: Sort and Group ───────────────────────────────────────────────
//
// sortAndGroup solves Department Top Three Salaries by collecting each
// department's DISTINCT salaries, sorting them descending, and cutting at
// the third one.
//
// Intuition:
//
//	Per department, the high earners are decided by a single number: the
//	third-highest DISTINCT salary. Everyone at or above that threshold is in;
//	everyone below is out. Departments with fewer than three distinct
//	salaries have no cut-off at all — every employee qualifies.
//
// Algorithm:
//  1. Group DISTINCT salaries per department into a set (dept → set).
//  2. For each department, dump the set into a slice, sort descending, and
//     record threshold[dept] = 3rd value (or math.MinInt when < 3 values).
//  3. Filter pass: emit every employee with e.Salary >= threshold[dept].
//
// Time:  O(n log n) — grouping is O(n); sorting all distinct salaries is at
//        worst O(n log n) across departments; the filter pass is O(n).
// Space: O(n) — the per-department distinct-salary sets.
func sortAndGroup(employees []Employee, departments []Department) []ResultRow {
	distinct := map[int]map[int]bool{} // departmentId → set of distinct salaries
	for _, e := range employees {
		if distinct[e.DepartmentID] == nil {
			distinct[e.DepartmentID] = map[int]bool{}
		}
		distinct[e.DepartmentID][e.Salary] = true // duplicates collapse here
	}

	threshold := map[int]int{} // departmentId → 3rd-highest distinct salary
	for dept, set := range distinct {
		salaries := make([]int, 0, len(set))
		for s := range set {
			salaries = append(salaries, s)
		}
		sort.Sort(sort.Reverse(sort.IntSlice(salaries))) // descending order
		if len(salaries) >= 3 {
			threshold[dept] = salaries[2] // the cut-off: 3rd distinct salary
		} else {
			threshold[dept] = math.MinInt // < 3 distinct → everyone qualifies
		}
	}

	nameByDept := nameMap(departments)
	result := []ResultRow{}
	for _, e := range employees { // filter pass against the per-dept cut-off
		if e.Salary >= threshold[e.DepartmentID] {
			result = append(result, ResultRow{nameByDept[e.DepartmentID], e.Name, e.Salary})
		}
	}
	return result
}

// ── Approach 3: Min-Heap Top-K ───────────────────────────────────────────────

// salaryMinHeap is a min-heap of salaries (container/heap interface); the
// smallest of the tracked top salaries sits at index 0, ready to be evicted.
type salaryMinHeap []int

func (h salaryMinHeap) Len() int           { return len(h) }
func (h salaryMinHeap) Less(i, j int) bool { return h[i] < h[j] }
func (h salaryMinHeap) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }
func (h *salaryMinHeap) Push(x any)        { *h = append(*h, x.(int)) }
func (h *salaryMinHeap) Pop() any {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[:n-1]
	return x
}

// minHeapTopK solves Department Top Three Salaries by keeping, per
// department, a size-3 min-heap of its distinct salaries.
//
// Intuition:
//
//	"Top three distinct salaries" is a top-K problem with K = 3. The standard
//	top-K tool is a min-heap of size K: its root is the weakest member of the
//	current top K, so any new (distinct) salary bigger than the root evicts
//	it. A companion set per department enforces distinctness. Unlike full
//	sorting, this generalises to any K at O(n log K).
//
// Algorithm:
//  1. Stream the employees; for each department keep (heap, set).
//  2. Skip salaries already in the set (distinct only). Push when the heap
//     holds < 3; otherwise, if salary > root, pop the root and push.
//  3. Filter pass: employee qualifies iff their department's heap holds
//     fewer than 3 salaries (no cut-off yet) or e.Salary >= heap root.
//
// Time:  O(n log K + d) = O(n) for K = 3 — each row does O(log 3) = O(1) heap work.
// Space: O(d) — at most 3 heap slots and 3 set entries per department.
func minHeapTopK(employees []Employee, departments []Department) []ResultRow {
	const k = 3
	heaps := map[int]*salaryMinHeap{}  // departmentId → min-heap of top salaries
	inHeap := map[int]map[int]bool{}   // departmentId → salaries currently held
	for _, e := range employees {
		h, ok := heaps[e.DepartmentID]
		if !ok { // first employee of this department: create its containers
			h = &salaryMinHeap{}
			heaps[e.DepartmentID] = h
			inHeap[e.DepartmentID] = map[int]bool{}
		}
		set := inHeap[e.DepartmentID]
		if set[e.Salary] {
			continue // this distinct salary is already tracked — skip duplicates
		}
		if h.Len() < k {
			heap.Push(h, e.Salary) // room left: admit unconditionally
			set[e.Salary] = true
		} else if e.Salary > (*h)[0] { // beats the weakest of the current top 3
			delete(set, (*h)[0])   // evict the root from the distinct set
			heap.Pop(h)            // ...and from the heap
			heap.Push(h, e.Salary) // admit the stronger salary
			set[e.Salary] = true
		}
	}

	nameByDept := nameMap(departments)
	result := []ResultRow{}
	for _, e := range employees {
		h := heaps[e.DepartmentID]
		// < 3 distinct salaries → no cut-off, everyone in the dept qualifies;
		// otherwise the heap root IS the 3rd-highest distinct salary.
		if h.Len() < k || e.Salary >= (*h)[0] {
			result = append(result, ResultRow{nameByDept[e.DepartmentID], e.Name, e.Salary})
		}
	}
	return result
}

// ── Approach 4: One-Pass Top-3 Tracking (Optimal) ────────────────────────────
//
// onePassTopThree solves Department Top Three Salaries by maintaining, per
// department, its three highest DISTINCT salaries in a fixed [3]int array.
//
// Intuition:
//
//	Since K is a constant 3, the heap is overkill: three ordered slots per
//	department, updated with constant-time compare-and-shift inserts, track
//	the top three distinct salaries exactly. After one pass, slot t[2] is the
//	department's cut-off (or MinInt when fewer than 3 distinct exist), and a
//	second pass keeps every employee at or above it.
//
// Algorithm:
//  1. Stream the employees; per department keep t = [3]int, descending, with
//     math.MinInt marking empty slots.
//  2. Insert each salary: ignore if it equals any slot (distinct only);
//     otherwise shift-and-place it into rank order.
//  3. Filter pass: emit every employee with e.Salary >= t[2] (MinInt makes
//     departments with < 3 distinct salaries accept everyone).
//
// Time:  O(n + d) — O(1) insert work per row, then a linear filter pass.
// Space: O(d) — exactly three tracked salaries per department.
func onePassTopThree(employees []Employee, departments []Department) []ResultRow {
	top := map[int]*[3]int{} // departmentId → its top-3 distinct salaries (desc)
	for _, e := range employees {
		t, ok := top[e.DepartmentID]
		if !ok { // first employee of this department: all slots empty
			t = &[3]int{math.MinInt, math.MinInt, math.MinInt}
			top[e.DepartmentID] = t
		}
		insertDistinctTop3(t, e.Salary) // O(1) compare-and-shift insert
	}

	nameByDept := nameMap(departments)
	result := []ResultRow{}
	for _, e := range employees {
		t := top[e.DepartmentID]
		// t[2] is the 3rd-highest distinct salary — the qualification cut-off.
		// When the dept has < 3 distinct salaries t[2] is MinInt, so every
		// employee passes, matching the problem's definition.
		if e.Salary >= t[2] {
			result = append(result, ResultRow{nameByDept[e.DepartmentID], e.Name, e.Salary})
		}
	}
	return result
}

// insertDistinctTop3 merges salary s into the descending top-3 array t,
// keeping only DISTINCT values. Empty slots hold math.MinInt.
func insertDistinctTop3(t *[3]int, s int) {
	switch {
	case s == t[0] || s == t[1] || s == t[2]:
		// already tracked — duplicates of a top salary change nothing
	case s > t[0]:
		t[0], t[1], t[2] = s, t[0], t[1] // new #1: shift old #1 and #2 down
	case s > t[1]:
		t[1], t[2] = s, t[1] // new #2: shift old #2 down
	case s > t[2]:
		t[2] = s // new #3: replace the cut-off
	}
}

// sortedRows orders result rows by (Department asc, Salary desc, Employee asc)
// so that the "return the result table in any order" outputs print
// deterministically across all approaches.
func sortedRows(rows []ResultRow) []ResultRow {
	out := make([]ResultRow, len(rows))
	copy(out, rows)
	sort.Slice(out, func(i, j int) bool {
		if out[i].Department != out[j].Department {
			return out[i].Department < out[j].Department
		}
		if out[i].Salary != out[j].Salary {
			return out[i].Salary > out[j].Salary
		}
		return out[i].Employee < out[j].Employee
	})
	return out
}

func main() {
	// Example 1 — Employee table:                Department table:
	//   | 1 | Joe   | 85000 | 1 |                  | 1 | IT    |
	//   | 2 | Henry | 80000 | 2 |                  | 2 | Sales |
	//   | 3 | Sam   | 60000 | 2 |
	//   | 4 | Max   | 90000 | 1 |
	//   | 5 | Janet | 69000 | 1 |
	//   | 6 | Randy | 85000 | 1 |
	//   | 7 | Will  | 70000 | 1 |
	employees := []Employee{
		{ID: 1, Name: "Joe", Salary: 85000, DepartmentID: 1},
		{ID: 2, Name: "Henry", Salary: 80000, DepartmentID: 2},
		{ID: 3, Name: "Sam", Salary: 60000, DepartmentID: 2},
		{ID: 4, Name: "Max", Salary: 90000, DepartmentID: 1},
		{ID: 5, Name: "Janet", Salary: 69000, DepartmentID: 1},
		{ID: 6, Name: "Randy", Salary: 85000, DepartmentID: 1},
		{ID: 7, Name: "Will", Salary: 70000, DepartmentID: 1},
	}
	departments := []Department{
		{ID: 1, Name: "IT"},
		{ID: 2, Name: "Sales"},
	}

	fmt.Println("=== Approach 1: Brute Force (Correlated Distinct Count) ===")
	fmt.Println(sortedRows(bruteForce(employees, departments)))
	// [{IT Max 90000} {IT Joe 85000} {IT Randy 85000} {IT Will 70000} {Sales Henry 80000} {Sales Sam 60000}]

	fmt.Println("=== Approach 2: Sort and Group ===")
	fmt.Println(sortedRows(sortAndGroup(employees, departments)))
	// [{IT Max 90000} {IT Joe 85000} {IT Randy 85000} {IT Will 70000} {Sales Henry 80000} {Sales Sam 60000}]

	fmt.Println("=== Approach 3: Min-Heap Top-K ===")
	fmt.Println(sortedRows(minHeapTopK(employees, departments)))
	// [{IT Max 90000} {IT Joe 85000} {IT Randy 85000} {IT Will 70000} {Sales Henry 80000} {Sales Sam 60000}]

	fmt.Println("=== Approach 4: One-Pass Top-3 Tracking (Optimal) ===")
	fmt.Println(sortedRows(onePassTopThree(employees, departments)))
	// [{IT Max 90000} {IT Joe 85000} {IT Randy 85000} {IT Will 70000} {Sales Henry 80000} {Sales Sam 60000}]
}
