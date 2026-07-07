# 0184 — Department Highest Salary

> LeetCode #184 · Difficulty: Medium
> **Categories:** Database, Hash Table, Sorting, Aggregation (Group By)

---

## Problem Statement

Table: `Employee`

```
+--------------+---------+
| Column Name  | Type    |
+--------------+---------+
| id           | int     |
| name         | varchar |
| salary       | int     |
| departmentId | int     |
+--------------+---------+
id is the primary key (column with unique values) for this table.
departmentId is a foreign key (reference columns) of the ID from the
Department table.
Each row of this table indicates the ID, name, and salary of an employee.
It also contains the ID of their department.
```

Table: `Department`

```
+-------------+---------+
| Column Name | Type    |
+-------------+---------+
| id          | int     |
| name        | varchar |
+-------------+---------+
id is the primary key (column with unique values) for this table.
It is guaranteed that department name is not NULL.
Each row of this table indicates the ID of a department and its name.
```

Write a solution to find employees who have the **highest salary in each of
the departments**.

Return the result table in **any order**.

**Example 1:**

```
Input:
Employee table:
+----+-------+--------+--------------+
| id | name  | salary | departmentId |
+----+-------+--------+--------------+
| 1  | Joe   | 70000  | 1            |
| 2  | Jim   | 90000  | 1            |
| 3  | Henry | 80000  | 2            |
| 4  | Sam   | 60000  | 2            |
| 5  | Max   | 90000  | 1            |
+----+-------+--------+--------------+
Department table:
+----+-------+
| id | name  |
+----+-------+
| 1  | IT    |
| 2  | Sales |
+----+-------+
Output:
+------------+----------+--------+
| Department | Employee | Salary |
+------------+----------+--------+
| IT         | Jim      | 90000  |
| Sales      | Henry    | 80000  |
| IT         | Max      | 90000  |
+------------+----------+--------+
Explanation: Max and Jim both have the highest salary in the IT department
and Henry has the highest salary in the Sales department.
```

**Constraints:**

- `id` columns are primary keys; `departmentId` always references a valid department (foreign key).
- Department names are never `NULL`.
- **Ties must all be reported** — every employee whose salary equals the department maximum appears in the result.

> **Repo note:** this is a LeetCode *Database* problem. In this Go repo the
> tables are modelled as `[]Employee{ID, Name, Salary, DepartmentID}` and
> `[]Department{ID, Name}`, results as `[]ResultRow{Department, Employee,
> Salary}`. Each approach re-implements `GROUP BY departmentId` + `MAX(salary)`
> plus the join back to department names.

---

## Company Frequency

| Company    | Frequency        | Last Reported |
|------------|------------------|---------------|
| Amazon     | ★★★★★ Very High  | 2024          |
| Microsoft  | ★★★★☆ High       | 2024          |
| Google     | ★★★☆☆ Medium     | 2024          |
| Meta       | ★★★☆☆ Medium     | 2023          |
| Bloomberg  | ★★★☆☆ Medium     | 2023          |
| Oracle     | ★★☆☆☆ Low        | 2023          |

> ⚠️ Frequency data is crowd-sourced from LeetCode Discuss, Glassdoor, and
> community interview reports. Treat as a signal, not a guarantee.

---

## DSA Concepts Used

- **Hash Map** — a departmentId → rolling-max map implements `GROUP BY` + `MAX` in one pass; a second map joins department ids to names → see [`/dsa/hash_map.md`](/dsa/hash_map.md)
- **Sorting** — ordering by (departmentId asc, salary desc) makes groups physical blocks whose first row is the group maximum → see [`/dsa/sorting.md`](/dsa/sorting.md)

---

## Approaches Overview

| # | Approach | Time | Space | When to use |
|---|----------|------|-------|-------------|
| 1 | Brute Force | O(n² + n·d) | O(1) | Baseline; direct "nobody strictly higher" definition, no aggregation |
| 2 | Sort and Group | O(n log n + d) | O(n + d) | When data must end up sorted anyway / GROUP BY intuition builder |
| 3 | Hash Map Two-Pass (Optimal) | O(n + d) | O(d) | Always — the streaming aggregate + filter pattern |

---

## Approach 1 — Brute Force

### Intuition

Flip the aggregation around: an employee has the department's highest salary
**iff no colleague in the same department earns strictly more**. That
predicate needs no MAX at all — test it directly for every employee with a
nested scan. Ties survive automatically, because for two tied maxima neither
is *strictly* higher than the other, so both pass.

### Algorithm

1. For each employee `e`, scan all employees.
2. If any `other` has `other.DepartmentID == e.DepartmentID` and
   `other.Salary > e.Salary`, mark `e` as beaten and stop the scan.
3. If `e` survives, scan `Department` for the row with `d.ID == e.DepartmentID`
   and emit `(d.Name, e.Name, e.Salary)`.
4. Return all emitted rows.

### Complexity

- **Time:** O(n² + n·d) — an O(n) colleague scan per employee, plus an O(d) department-name scan per winner.
- **Space:** O(1) — flags and loop variables only.

### Code

```go
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
```

### Dry Run

Example 1: `employees = [Joe 70000 d1, Jim 90000 d1, Henry 80000 d2, Sam 60000 d2, Max 90000 d1]`

| Step | `e` | Same-dept colleague strictly higher? | `isMax` | Dept name | `result` |
|------|-----|--------------------------------------|---------|-----------|----------|
| 1 | Joe (70000, d1) | Jim 90000 > 70000 → yes | false | — | `[]` |
| 2 | Jim (90000, d1) | Joe 70000 ✗, Max 90000 ✗ (not strict) | true | IT | `[{IT Jim 90000}]` |
| 3 | Henry (80000, d2) | Sam 60000 ✗ | true | Sales | `+ {Sales Henry 80000}` |
| 4 | Sam (60000, d2) | Henry 80000 > 60000 → yes | false | — | unchanged |
| 5 | Max (90000, d1) | Joe ✗, Jim 90000 ✗ (not strict) | true | IT | `+ {IT Max 90000}` |

Return `{IT Jim 90000}, {Sales Henry 80000}, {IT Max 90000}` ✓ — exactly the
three expected rows (order is irrelevant).

---

## Approach 2 — Sort and Group

### Intuition

`GROUP BY` made physical: sort by `(departmentId asc, salary desc)` and each
department becomes a contiguous block whose **first row carries the block's
maximum salary**. Walk the sorted rows once — on entering a block, remember
its leader's salary; emit every row tied with it. As soon as a lower salary
appears, the equality test fails for the rest of the block.

### Algorithm

1. Copy `Employee` and sort the copy by departmentId ascending, salary descending.
2. Build `nameByDept` (departmentId → name) once for O(1) name joins.
3. Scan the sorted rows; when `i == 0` or the departmentId changed, set
   `blockMax` to the current row's salary (the new block's maximum).
4. Emit `(deptName, e.Name, e.Salary)` for every row with `e.Salary == blockMax`.

### Complexity

- **Time:** O(n log n + d) — the sort dominates; the block scan is O(n) and the name map costs O(d).
- **Space:** O(n + d) — the sorted copy plus the department-name map.

### Code

```go
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
```

### Dry Run

Example 1: after sorting,
`rows = [Jim 90000 d1, Max 90000 d1, Joe 70000 d1, Henry 80000 d2, Sam 60000 d2]`
and `nameByDept = {1: IT, 2: Sales}`.

| Step | `i` | Row | New block? | `blockMax` | `Salary == blockMax`? | `result` |
|------|-----|-----|------------|------------|------------------------|----------|
| 1 | 0 | Jim 90000 d1 | yes (i = 0) | 90000 | yes | `[{IT Jim 90000}]` |
| 2 | 1 | Max 90000 d1 | no | 90000 | yes (tie) | `+ {IT Max 90000}` |
| 3 | 2 | Joe 70000 d1 | no | 90000 | no | unchanged |
| 4 | 3 | Henry 80000 d2 | yes (d1→d2) | 80000 | yes | `+ {Sales Henry 80000}` |
| 5 | 4 | Sam 60000 d2 | no | 80000 | no | unchanged |

Return `{IT Jim 90000}, {IT Max 90000}, {Sales Henry 80000}` ✓.

---

## Approach 3 — Hash Map Two-Pass (Optimal)

### Intuition

`GROUP BY departmentId` + `MAX(salary)` is a **rolling maximum per key**: one
linear pass keeps, for each department, the best salary seen so far in a hash
map. But the query also needs *who* earns it — and there may be ties — so a
second pass re-reads the rows and emits every employee whose salary equals
their department's recorded maximum. Aggregate, then filter.

### Algorithm

1. **Aggregate pass:** for each employee, `maxSalary[dept] = max(maxSalary[dept], e.Salary)`
   (a missing key is treated as "new department").
2. Build `nameByDept` (departmentId → name) for the final join.
3. **Filter pass:** emit `(deptName, e.Name, e.Salary)` for every employee
   with `e.Salary == maxSalary[e.DepartmentID]`.

### Complexity

- **Time:** O(n + d) — two passes over Employee and one over Department, all with O(1) average map operations.
- **Space:** O(d) — one rolling-max entry and one name entry per department.

### Code

```go
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
```

### Dry Run

Example 1 — **aggregate pass** over the five employees:

| Step | Employee | `maxSalary` after update |
|------|----------|---------------------------|
| 1 | Joe 70000 d1 | `{1: 70000}` (new department) |
| 2 | Jim 90000 d1 | `{1: 90000}` (90000 > 70000) |
| 3 | Henry 80000 d2 | `{1: 90000, 2: 80000}` (new department) |
| 4 | Sam 60000 d2 | unchanged (60000 < 80000) |
| 5 | Max 90000 d1 | unchanged (90000 not > 90000) |

**Filter pass** with `maxSalary = {1: 90000, 2: 80000}`, `nameByDept = {1: IT, 2: Sales}`:

| Step | Employee | `Salary == maxSalary[dept]`? | `result` |
|------|----------|-------------------------------|----------|
| 1 | Joe 70000 d1 | 70000 == 90000 → no | `[]` |
| 2 | Jim 90000 d1 | 90000 == 90000 → yes | `[{IT Jim 90000}]` |
| 3 | Henry 80000 d2 | 80000 == 80000 → yes | `+ {Sales Henry 80000}` |
| 4 | Sam 60000 d2 | 60000 == 80000 → no | unchanged |
| 5 | Max 90000 d1 | 90000 == 90000 → yes | `+ {IT Max 90000}` |

Return `{IT Jim 90000}, {Sales Henry 80000}, {IT Max 90000}` ✓.

---

## Key Takeaways

- **Aggregate-then-filter is the universal "group max with ties" recipe:** pass 1 computes `MAX` per key into a map, pass 2 keeps rows equal to their key's max. It generalises to MIN/SUM/COUNT per group in one pass each.
- **You cannot pick winners during the aggregate pass** — a row that looks like the max early on may be beaten later, and ties would be missed. The second pass is what makes the answer complete and tie-safe.
- **"Nobody strictly higher" ≡ "is a maximum"** — the brute-force predicate — handles ties for free and is a handy correctness oracle for testing cleverer solutions.
- Sorting by `(group key, value desc)` turns groups into blocks led by their maximum — the physical picture behind `GROUP BY`.
- Canonical SQL solution this maps to:

```sql
SELECT d.name AS Department, e.name AS Employee, e.salary AS Salary
FROM Employee e
JOIN Department d ON e.departmentId = d.id
WHERE (e.departmentId, e.salary) IN (
    SELECT departmentId, MAX(salary)
    FROM Employee
    GROUP BY departmentId
);
```

---

## Related Problems

- LeetCode #185 — Department Top Three Salaries (same pattern generalised from top-1 to top-3 distinct)
- LeetCode #176 — Second Highest Salary (ranking salaries instead of taking the max)
- LeetCode #177 — Nth Highest Salary (parameterised ranking)
- LeetCode #178 — Rank Scores (dense ranking, the concept behind ties here)
- LeetCode #181 — Employees Earning More Than Their Managers (same join-via-hash-map modelling)
