package main

import (
	"fmt"
	"sort"
)

// LeetCode 181 — Employees Earning More Than Their Managers.
//
// The original problem is a SQL one: given the Employee table
//
//	+----+-------+--------+-----------+
//	| id | name  | salary | managerId |
//	+----+-------+--------+-----------+
//
// report the names of employees who earn strictly more than their managers.
// Here the table is modelled as a slice of Employee rows, and every approach
// re-implements the query as an in-memory join strategy.

// Employee models one row of the Employee table.
// ManagerID == 0 stands for SQL NULL (real ids start at 1, so 0 can never
// match any row — exactly like a NULL join key, which joins to nothing).
type Employee struct {
	ID        int
	Name      string
	Salary    int
	ManagerID int
}

// ── Approach 1: Brute Force (Nested-Loop Join) ───────────────────────────────
//
// bruteForce solves Employees Earning More Than Their Managers with a
// nested-loop self-join of the Employee table.
//
// Intuition:
//
//	SQL's `FROM Employee a JOIN Employee b ON a.managerId = b.id` evaluated
//	naively compares every row against every other row. Do exactly that: for
//	each employee, scan the whole table to find the row holding their manager,
//	then compare the two salaries.
//
// Algorithm:
//  1. For every employee e, scan all rows looking for m with m.ID == e.ManagerID.
//  2. A NULL manager (ManagerID 0) matches no row, so e is skipped naturally.
//  3. If e.Salary > m.Salary, append e.Name to the result.
//
// Time:  O(n^2) — a full O(n) manager scan for each of the n rows.
// Space: O(1) — no auxiliary structures beyond the output slice.
func bruteForce(employees []Employee) []string {
	result := []string{}
	for _, e := range employees { // outer side of the self-join
		for _, m := range employees { // inner scan: locate e's manager row
			if m.ID == e.ManagerID { // join condition: a.managerId = b.id
				if e.Salary > m.Salary { // filter: a.salary > b.salary
					result = append(result, e.Name)
				}
				break // ids are unique — manager found, stop the inner scan
			}
		}
	}
	return result
}

// ── Approach 2: Sort + Binary Search ─────────────────────────────────────────
//
// sortAndBinarySearch solves Employees Earning More Than Their Managers by
// sorting a copy of the table by id and binary-searching each manager row.
//
// Intuition:
//
//	The inner O(n) scan is just a lookup by primary key. Sorting the table by
//	id turns each lookup into an O(log n) binary search — the in-memory
//	analogue of a database index on the primary-key column.
//
// Algorithm:
//  1. Copy the table and sort the copy by ID ascending.
//  2. For each employee, binary search the copy for their ManagerID.
//  3. If the manager row exists and e.Salary > manager.Salary, keep e.Name.
//
// Time:  O(n log n) — sorting dominates; each of the n lookups is O(log n).
// Space: O(n) — the sorted copy of the table.
func sortAndBinarySearch(employees []Employee) []string {
	byID := make([]Employee, len(employees))
	copy(byID, employees) // never mutate the caller's table
	sort.Slice(byID, func(i, j int) bool { return byID[i].ID < byID[j].ID })

	result := []string{}
	for _, e := range employees {
		// First index whose ID is >= the manager id we are looking for.
		idx := sort.Search(len(byID), func(i int) bool { return byID[i].ID >= e.ManagerID })
		// A NULL manager (0) or a missing id lands on a non-matching slot → skip.
		if idx < len(byID) && byID[idx].ID == e.ManagerID && e.Salary > byID[idx].Salary {
			result = append(result, e.Name) // earns strictly more than the manager
		}
	}
	return result
}

// ── Approach 3: Hash Map (Hash Join, Optimal) ────────────────────────────────
//
// hashMap solves Employees Earning More Than Their Managers with an
// id → salary hash map, i.e. a classic hash join.
//
// Intuition:
//
//	A database executes this query efficiently as a hash join: build a hash
//	table keyed on the join column (id), then probe it once per row. Two
//	linear passes replace the quadratic nested scan.
//
// Algorithm:
//  1. Pass 1 (build): salaryByID[e.ID] = e.Salary for every row.
//  2. Pass 2 (probe): look up salaryByID[e.ManagerID]; the probe misses for
//     NULL (0) managers, mirroring SQL NULL join semantics.
//  3. Keep e.Name whenever the probe hits and e.Salary > managerSalary.
//
// Time:  O(n) — two passes with O(1) average map operations.
// Space: O(n) — the hash map holds one entry per row.
func hashMap(employees []Employee) []string {
	salaryByID := make(map[int]int, len(employees)) // build side: id → salary
	for _, e := range employees {
		salaryByID[e.ID] = e.Salary
	}

	result := []string{}
	for _, e := range employees { // probe side
		if managerSalary, ok := salaryByID[e.ManagerID]; ok && e.Salary > managerSalary {
			result = append(result, e.Name) // earns strictly more than the manager
		}
	}
	return result
}

// sorted returns an alphabetically sorted copy of names so that the
// "return the result table in any order" outputs print deterministically.
func sorted(names []string) []string {
	out := make([]string, len(names))
	copy(out, names)
	sort.Strings(out)
	return out
}

func main() {
	// Example 1 — Employee table (ManagerID 0 models SQL NULL):
	//   | 1 | Joe   | 70000 | 3    |
	//   | 2 | Henry | 80000 | 4    |
	//   | 3 | Sam   | 60000 | null |
	//   | 4 | Max   | 90000 | null |
	employees := []Employee{
		{ID: 1, Name: "Joe", Salary: 70000, ManagerID: 3},
		{ID: 2, Name: "Henry", Salary: 80000, ManagerID: 4},
		{ID: 3, Name: "Sam", Salary: 60000, ManagerID: 0},
		{ID: 4, Name: "Max", Salary: 90000, ManagerID: 0},
	}

	fmt.Println("=== Approach 1: Brute Force (Nested-Loop Join) ===")
	fmt.Println(sorted(bruteForce(employees))) // [Joe]

	fmt.Println("=== Approach 2: Sort + Binary Search ===")
	fmt.Println(sorted(sortAndBinarySearch(employees))) // [Joe]

	fmt.Println("=== Approach 3: Hash Map (Hash Join, Optimal) ===")
	fmt.Println(sorted(hashMap(employees))) // [Joe]
}
