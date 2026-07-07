# 0181 — Employees Earning More Than Their Managers

> LeetCode #181 · Difficulty: Easy
> **Categories:** Database, Hash Table, Sorting, Binary Search

---

## Problem Statement

Table: `Employee`

```
+-------------+---------+
| Column Name | Type    |
+-------------+---------+
| id          | int     |
| name        | varchar |
| salary      | int     |
| managerId   | int     |
+-------------+---------+
id is the primary key (column with unique values) for this table.
Each row of this table indicates the ID of an employee, their name,
salary, and the ID of their manager.
```

Write a solution to find the employees who earn **more than their managers**.

Return the result table in **any order**.

**Example 1:**

```
Input:
Employee table:
+----+-------+--------+-----------+
| id | name  | salary | managerId |
+----+-------+--------+-----------+
| 1  | Joe   | 70000  | 3         |
| 2  | Henry | 80000  | 4         |
| 3  | Sam   | 60000  | null      |
| 4  | Max   | 90000  | null      |
+----+-------+--------+-----------+
Output:
+----------+
| Employee |
+----------+
| Joe      |
+----------+
Explanation: Joe is the only employee who earns more than his manager.
```

**Constraints:**

- `id` is the primary key — every id is unique.
- `managerId` may be `null` (the employee has no manager) and otherwise references an `id` in the same table.
- "More than" means **strictly** greater salary.

> **Repo note:** this is a LeetCode *Database* problem. In this Go repo the
> table is modelled as `[]Employee{ID, Name, Salary, ManagerID}`, with
> `ManagerID == 0` standing in for SQL `NULL` (real ids start at 1, so 0 joins
> to nothing — exactly like a NULL join key). Each approach re-implements the
> SQL query as an in-memory join strategy.

---

## Company Frequency

| Company    | Frequency        | Last Reported |
|------------|------------------|---------------|
| Amazon     | ★★★★☆ High       | 2024          |
| Microsoft  | ★★★☆☆ Medium     | 2024          |
| Google     | ★★★☆☆ Medium     | 2023          |
| Bloomberg  | ★★★☆☆ Medium     | 2023          |
| Meta       | ★★☆☆☆ Low        | 2023          |
| Adobe      | ★★☆☆☆ Low        | 2023          |

> ⚠️ Frequency data is crowd-sourced from LeetCode Discuss, Glassdoor, and
> community interview reports. Treat as a signal, not a guarantee.

---

## DSA Concepts Used

- **Hash Map** — the id → salary map is the build side of a hash join; probing it by `managerId` makes each manager lookup O(1) → see [`/dsa/hash_map.md`](/dsa/hash_map.md)
- **Sorting** — sorting the table by primary key simulates a database index, enabling logarithmic lookups → see [`/dsa/sorting.md`](/dsa/sorting.md)
- **Binary Search** — locating a manager row by id inside the sorted copy → see [`/dsa/binary_search.md`](/dsa/binary_search.md)

---

## Approaches Overview

| # | Approach | Time | Space | When to use |
|---|----------|------|-------|-------------|
| 1 | Brute Force (Nested-Loop Join) | O(n²) | O(1) | Baseline; fine for tiny tables, mirrors naive SQL join evaluation |
| 2 | Sort + Binary Search | O(n log n) | O(n) | When memory for a hash table is tight but a sorted copy/index is acceptable |
| 3 | Hash Map (Hash Join, Optimal) | O(n) | O(n) | Always — the classic hash-join answer interviewers expect |

---

## Approach 1 — Brute Force (Nested-Loop Join)

### Intuition

`SELECT a.name FROM Employee a JOIN Employee b ON a.managerId = b.id WHERE a.salary > b.salary`
evaluated the naive way is a nested-loop self-join: pair every row with every
other row, keep the pairs where the join condition and the salary filter hold.
An employee with a `NULL` manager joins to nothing, so they can never appear —
our `ManagerID == 0` sentinel reproduces that for free because no row has id 0.

### Algorithm

1. Loop over every employee `e` (the outer side of the join).
2. For each `e`, scan the entire table for the row `m` with `m.ID == e.ManagerID`.
3. If no row matches (NULL manager), `e` is skipped naturally.
4. If a match is found and `e.Salary > m.Salary`, append `e.Name` to the result.
5. Because ids are unique, `break` out of the inner scan after the first match.

### Complexity

- **Time:** O(n²) — each of the n rows triggers a full O(n) scan for its manager.
- **Space:** O(1) — only loop variables; the output slice is not counted as auxiliary space.

### Code

```go
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
```

### Dry Run

Example 1: `employees = [Joe(1, 70000, mgr 3), Henry(2, 80000, mgr 4), Sam(3, 60000, mgr NULL→0), Max(4, 90000, mgr NULL→0)]`

| Step | Outer `e` | Inner scan finds `m` | `e.Salary > m.Salary`? | `result` |
|------|-----------|----------------------|------------------------|----------|
| 1 | Joe (70000, mgr 3) | Sam (id 3, 60000) | 70000 > 60000 → yes | `[Joe]` |
| 2 | Henry (80000, mgr 4) | Max (id 4, 90000) | 80000 > 90000 → no | `[Joe]` |
| 3 | Sam (60000, mgr 0) | no row has id 0 → no match | — (skipped) | `[Joe]` |
| 4 | Max (90000, mgr 0) | no row has id 0 → no match | — (skipped) | `[Joe]` |

Return `[Joe]` ✓ (matches the expected output).

---

## Approach 2 — Sort + Binary Search

### Intuition

The inner scan of Approach 1 is nothing but a lookup by primary key. If the
table is sorted by `id`, that lookup becomes an O(log n) binary search — this
is literally what a database index on the primary key does. Sort once, then
answer n lookups fast.

### Algorithm

1. Copy the table (never mutate the input) and sort the copy by `ID` ascending.
2. For each employee `e`, binary search the sorted copy for `e.ManagerID`
   using `sort.Search` (first index with `ID >= ManagerID`).
3. If the found slot holds exactly `ManagerID` (probe hit) **and**
   `e.Salary > manager.Salary`, append `e.Name`.
4. A NULL manager (`0`) or missing id lands on a non-matching slot → skipped.

### Complexity

- **Time:** O(n log n) — O(n log n) for the sort, plus n binary searches at O(log n) each.
- **Space:** O(n) — the sorted copy of the table.

### Code

```go
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
```

### Dry Run

Example 1: after sorting, `byID = [Joe(1), Henry(2), Sam(3), Max(4)]` (already ordered by id).

| Step | `e` | Search target | `idx` | `byID[idx].ID == target`? | Salary check | `result` |
|------|-----|---------------|-------|---------------------------|--------------|----------|
| 1 | Joe (70000) | 3 | 2 | yes → Sam (60000) | 70000 > 60000 ✓ | `[Joe]` |
| 2 | Henry (80000) | 4 | 3 | yes → Max (90000) | 80000 > 90000 ✗ | `[Joe]` |
| 3 | Sam (60000) | 0 | 0 | `byID[0].ID = 1 ≠ 0` → miss | skipped | `[Joe]` |
| 4 | Max (90000) | 0 | 0 | `byID[0].ID = 1 ≠ 0` → miss | skipped | `[Joe]` |

Return `[Joe]` ✓.

---

## Approach 3 — Hash Map (Hash Join, Optimal)

### Intuition

This is exactly how a real query planner executes the self-join efficiently: a
**hash join**. Build a hash table keyed on the join column (`id`) in one pass,
then probe it once per row in a second pass. Each probe is O(1) on average, so
the whole query collapses to linear time.

### Algorithm

1. **Build pass:** fill `salaryByID[e.ID] = e.Salary` for every row.
2. **Probe pass:** for each employee `e`, look up `salaryByID[e.ManagerID]`.
3. The probe misses for NULL (`0`) managers — mirroring SQL's NULL join semantics.
4. When the probe hits and `e.Salary > managerSalary`, append `e.Name`.

### Complexity

- **Time:** O(n) — one build pass plus one probe pass, O(1) average per map operation.
- **Space:** O(n) — the map stores one id → salary entry per row.

### Code

```go
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
```

### Dry Run

Example 1 — build pass produces `salaryByID = {1: 70000, 2: 80000, 3: 60000, 4: 90000}`.

| Step | Probe `e` | Lookup `salaryByID[e.ManagerID]` | `ok`? | `e.Salary > managerSalary`? | `result` |
|------|-----------|----------------------------------|-------|------------------------------|----------|
| 1 | Joe (70000, mgr 3) | `salaryByID[3] = 60000` | yes | 70000 > 60000 ✓ | `[Joe]` |
| 2 | Henry (80000, mgr 4) | `salaryByID[4] = 90000` | yes | 80000 > 90000 ✗ | `[Joe]` |
| 3 | Sam (60000, mgr 0) | `salaryByID[0]` → miss | no | skipped | `[Joe]` |
| 4 | Max (90000, mgr 0) | `salaryByID[0]` → miss | no | skipped | `[Joe]` |

Return `[Joe]` ✓.

---

## Key Takeaways

- **A self-join is just a lookup problem.** "Find each row's manager" = "look up a row by primary key" — the three classic lookup strategies (linear scan, sorted + binary search, hash table) map 1:1 onto the three approaches here.
- **Hash join pattern:** build a map on the join key in one pass, probe it in a second pass — O(n²) → O(n). This is the single most reusable trick for turning SQL-style problems into linear-time Go.
- **Model SQL NULL as a value that can never match** (here `0` when ids start at 1). NULL join keys silently join to nothing; a failed map probe (`ok == false`) reproduces that behaviour exactly.
- Canonical SQL solutions this maps to:

```sql
-- Self-join
SELECT a.name AS Employee
FROM Employee a JOIN Employee b ON a.managerId = b.id
WHERE a.salary > b.salary;

-- Correlated subquery
SELECT name AS Employee FROM Employee e
WHERE salary > (SELECT salary FROM Employee WHERE id = e.managerId);
```

---

## Related Problems

- LeetCode #182 — Duplicate Emails (same table-modelling + hash-map pattern)
- LeetCode #183 — Customers Who Never Order (hash anti-join instead of hash join)
- LeetCode #184 — Department Highest Salary (join + per-group aggregation)
- LeetCode #570 — Managers with at Least 5 Direct Reports (self-join + grouping)
- LeetCode #577 — Employee Bonus (left join with NULL filtering)
