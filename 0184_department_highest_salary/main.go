package main

import (
	"fmt"
	"sort"
)

// LeetCode 184 — Department Highest Salary.
//
// The original problem is a SQL one: given the tables
//
//	Employee: +----+------+--------+--------------+   Department: +----+------+
//	          | id | name | salary | departmentId |               | id | name |
//	          +----+------+--------+--------------+               +----+------+
//
// report every employee who has the highest salary in their department, as
// rows (Department, Employee, Salary) — ties included, any order. Here the
// tables are modelled as slices of rows, and every approach re-implements the
// GROUP BY departmentId + MAX(salary) query plus the join back to names.

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

// ── Approach 1: Brute Force ──────────────────────────────────────────────────
//
// bruteForce solves Department Highest Salary by checking, for every
// employee, whether anyone in the same department out-earns them.
//
// Intuition:
//
//	An employee has the department's highest salary iff NO colleague in the
//	same department earns strictly more. Testing that directly needs no
//	aggregation at all — just a nested scan per employee. Ties survive
//	naturally because "nobody strictly higher" holds for every tied maximum.
//
// Algorithm:
//  1. For each employee e, scan all employees; if any shares e's department
//     and earns strictly more, e is not a department maximum.
//  2. If e survives, linearly scan Department for e's department name and
//     emit the row (deptName, e.Name, e.Salary).
//
// Time:  O(n^2 + n*d) — an O(n) colleague scan per employee, plus an O(d)
//        name scan per winner.
// Space: O(1) — no auxiliary structures beyond the output slice.
func bruteForce(employees []Employee, departments []Department) []ResultRow {
	result := []ResultRow{}
	for _, e := range employees {
		isMax := true
		for _, other := range employees { // does any colleague out-earn e?
			if other.DepartmentID == e.DepartmentID && other.Salary > e.Salary {
				isMax = false // strictly higher salary found → e loses
				break
			}
		}
		if !isMax {
			continue
		}
		for _, d := range departments { // join: departmentId → department name
			if d.ID == e.DepartmentID {
				result = append(result, ResultRow{d.Name, e.Name, e.Salary})
				break // department ids are unique
			}
		}
	}
	return result
}

// ── Approach 2: Sort and Group ───────────────────────────────────────────────
//
// sortAndGroup solves Department Highest Salary by sorting rows so each
// department forms a contiguous block that starts with its maximum salary.
//
// Intuition:
//
//	Sorting by (departmentId asc, salary desc) is SQL's GROUP BY made
//	physical: each department becomes one block, and the block's first row
//	carries the department's maximum salary. Every leading row tied with
//	that maximum is a winner; the first lower salary ends the winning run.
//
// Algorithm:
//  1. Sort a copy of Employee by departmentId asc, then salary desc.
//  2. Build a departmentId → name map once for O(1) name joins.
//  3. Scan the sorted rows: on entering a new block, record blockMax from its
//     first row; emit every row whose salary equals blockMax.
//
// Time:  O(n log n + d) — the sort dominates; the scan and map build are linear.
// Space: O(n + d) — the sorted copy plus the department-name map.
func sortAndGroup(employees []Employee, departments []Department) []ResultRow {
	rows := make([]Employee, len(employees))
	copy(rows, employees) // sort a copy — never mutate the caller's table
	sort.Slice(rows, func(i, j int) bool {
		if rows[i].DepartmentID != rows[j].DepartmentID {
			return rows[i].DepartmentID < rows[j].DepartmentID // group departments
		}
		return rows[i].Salary > rows[j].Salary // highest salary first inside a block
	})

	nameByDept := make(map[int]string, len(departments)) // departmentId → name
	for _, d := range departments {
		nameByDept[d.ID] = d.Name
	}

	result := []ResultRow{}
	blockMax := 0 // maximum salary of the department block being scanned
	for i, e := range rows {
		if i == 0 || rows[i-1].DepartmentID != e.DepartmentID {
			blockMax = e.Salary // block leader = this department's max salary
		}
		if e.Salary == blockMax { // leader and everyone tied with it win
			result = append(result, ResultRow{nameByDept[e.DepartmentID], e.Name, e.Salary})
		}
	}
	return result
}

// ── Approach 3: Hash Map Two-Pass (Optimal) ──────────────────────────────────
//
// hashMap solves Department Highest Salary with a departmentId → max-salary
// map built in one pass and applied in a second.
//
// Intuition:
//
//	GROUP BY + MAX is a rolling maximum per key: one pass over the rows keeps
//	the best salary seen per department in a hash map. A second pass then
//	emits every employee whose salary equals their department's recorded
//	maximum — ties included by construction.
//
// Algorithm:
//  1. Pass 1 (aggregate): maxSalary[dept] = max(maxSalary[dept], e.Salary).
//  2. Build the departmentId → name map for the final join.
//  3. Pass 2 (filter): emit (deptName, e.Name, e.Salary) for every employee
//     with e.Salary == maxSalary[e.DepartmentID].
//
// Time:  O(n + d) — two linear passes over Employee and one over Department.
// Space: O(d) — one max entry and one name entry per department.
func hashMap(employees []Employee, departments []Department) []ResultRow {
	maxSalary := map[int]int{} // departmentId → highest salary seen so far
	for _, e := range employees {
		if cur, ok := maxSalary[e.DepartmentID]; !ok || e.Salary > cur {
			maxSalary[e.DepartmentID] = e.Salary // new department or new maximum
		}
	}

	nameByDept := make(map[int]string, len(departments)) // departmentId → name
	for _, d := range departments {
		nameByDept[d.ID] = d.Name
	}

	result := []ResultRow{}
	for _, e := range employees { // filter pass: keep the per-department maxima
		if e.Salary == maxSalary[e.DepartmentID] {
			result = append(result, ResultRow{nameByDept[e.DepartmentID], e.Name, e.Salary})
		}
	}
	return result
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
	//   | 1 | Joe   | 70000 | 1 |                  | 1 | IT    |
	//   | 2 | Jim   | 90000 | 1 |                  | 2 | Sales |
	//   | 3 | Henry | 80000 | 2 |
	//   | 4 | Sam   | 60000 | 2 |
	//   | 5 | Max   | 90000 | 1 |
	employees := []Employee{
		{ID: 1, Name: "Joe", Salary: 70000, DepartmentID: 1},
		{ID: 2, Name: "Jim", Salary: 90000, DepartmentID: 1},
		{ID: 3, Name: "Henry", Salary: 80000, DepartmentID: 2},
		{ID: 4, Name: "Sam", Salary: 60000, DepartmentID: 2},
		{ID: 5, Name: "Max", Salary: 90000, DepartmentID: 1},
	}
	departments := []Department{
		{ID: 1, Name: "IT"},
		{ID: 2, Name: "Sales"},
	}

	fmt.Println("=== Approach 1: Brute Force ===")
	fmt.Println(sortedRows(bruteForce(employees, departments))) // [{IT Jim 90000} {IT Max 90000} {Sales Henry 80000}]

	fmt.Println("=== Approach 2: Sort and Group ===")
	fmt.Println(sortedRows(sortAndGroup(employees, departments))) // [{IT Jim 90000} {IT Max 90000} {Sales Henry 80000}]

	fmt.Println("=== Approach 3: Hash Map Two-Pass (Optimal) ===")
	fmt.Println(sortedRows(hashMap(employees, departments))) // [{IT Jim 90000} {IT Max 90000} {Sales Henry 80000}]
}
