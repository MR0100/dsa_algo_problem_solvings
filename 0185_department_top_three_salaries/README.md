# 0185 — Department Top Three Salaries

> LeetCode #185 · Difficulty: Hard
> **Categories:** Database (SQL), Hash Map, Sorting, Heap / Priority Queue, Top-K

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
```

`id` is the primary key (column with unique values) for this table.
`departmentId` is a foreign key (reference column) of the `id` from the `Department` table.
Each row of this table indicates the ID, name, and salary of an employee. It also contains the ID of their department.

Table: `Department`

```
+-------------+---------+
| Column Name | Type    |
+-------------+---------+
| id          | int     |
| name        | varchar |
+-------------+---------+
```

`id` is the primary key (column with unique values) for this table.
Each row of this table indicates the ID of a department and its name.

A company's executives are interested in seeing who earns the most money in each of the company's departments. A **high earner** in a department is an employee who has a salary in the **top three unique salaries** for that department.

Write a solution to find the employees who are high earners in each of the departments.

Return the result table **in any order**.

The result format is in the following example.

**Example 1:**

```
Input:
Employee table:
+----+-------+--------+--------------+
| id | name  | salary | departmentId |
+----+-------+--------+--------------+
| 1  | Joe   | 85000  | 1            |
| 2  | Henry | 80000  | 2            |
| 3  | Sam   | 60000  | 2            |
| 4  | Max   | 90000  | 1            |
| 5  | Janet | 69000  | 1            |
| 6  | Randy | 85000  | 1            |
| 7  | Will  | 70000  | 1            |
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
| IT         | Max      | 90000  |
| IT         | Joe      | 85000  |
| IT         | Randy    | 85000  |
| IT         | Will     | 70000  |
| Sales      | Henry    | 80000  |
| Sales      | Sam      | 60000  |
+------------+----------+--------+

Explanation:
In the IT department:
- Max earns the highest unique salary
- Both Randy and Joe earn the second-highest unique salary
- Will earns the third-highest unique salary

In the Sales department:
- Henry earns the highest salary
- Sam earns the second-highest salary
- There is no third-highest salary as there are only two employees
```

**Constraints:**

- There are no employees with the exact same name, salary, and department.

---

## Company Frequency

| Company    | Frequency        | Last Reported |
|------------|------------------|---------------|
| Amazon     | ★★★★☆ High       | 2024          |
| Meta       | ★★★★☆ High       | 2023          |
| Google     | ★★★☆☆ Medium     | 2023          |
| Bloomberg  | ★★★☆☆ Medium     | 2023          |
| Microsoft  | ★★☆☆☆ Low        | 2022          |
| Apple      | ★★☆☆☆ Low        | 2022          |

> ⚠️ Frequency data is crowd-sourced from LeetCode Discuss, Glassdoor, and
> community interview reports. Treat as a signal, not a guarantee.

---

## DSA Concepts Used

- **Hash Map** — group employees by `departmentId`, deduplicate salaries with a set, and build the `departmentId → name` lookup, all in O(1) amortized per operation → see [`/dsa/hash_map.md`](/dsa/hash_map.md)
- **Sorting** — Approach 2 collects each department's distinct salaries and sorts them descending to read off the third-highest (the cut-off) → see [`/dsa/sorting.md`](/dsa/sorting.md)
- **Heap / Priority Queue** — Approach 3 frames "top three distinct salaries" as a classic top-K problem, maintaining a size-3 min-heap whose root is the current cut-off → see [`/dsa/heap_priority_queue.md`](/dsa/heap_priority_queue.md)

---

## Approaches Overview

| # | Approach | Time | Space | When to use |
|---|----------|------|-------|-------------|
| 1 | Brute Force (Correlated Distinct Count) | O(n²) | O(n) | Direct translation of the SQL correlated subquery; fine for small tables, TLE on large ones |
| 2 | Sort and Group | O(n log n) | O(n) | Clean and intuitive; sort each department's distinct salaries and cut at the third |
| 3 | Min-Heap Top-K | O(n log K) = O(n) | O(d) | The general top-K pattern; scales to any K without re-sorting |
| 4 | One-Pass Top-3 Tracking (Optimal) | O(n) | O(d) | K is the constant 3, so three ordered slots per department beat the heap |

Where `n` = number of employees and `d` = number of departments.

---

## Approach 1 — Brute Force (Correlated Distinct Count)

### Intuition

"Salary is among the top three distinct salaries" is exactly equivalent to "fewer than three **distinct** salaries in the same department are strictly higher." That is the textbook SQL correlated subquery `3 > (SELECT COUNT(DISTINCT e2.salary) FROM Employee e2 WHERE e2.departmentId = e1.departmentId AND e2.salary > e1.salary)`. We evaluate it literally: for each employee, scan the whole table and count the distinct higher salaries in their department.

### Algorithm

1. Build the `departmentId → name` lookup once.
2. For each employee `e`, scan all employees of the same department and collect into a set the **distinct** salaries strictly greater than `e.Salary`.
3. If that set holds fewer than 3 salaries, `e` is a high earner — emit `(deptName, e.Name, e.Salary)`.

### Complexity

- **Time:** O(n²) — an O(n) inner scan (with O(1) set inserts) for each of the n employees.
- **Space:** O(n) — the per-employee set of higher distinct salaries (reused scratch), plus O(d) for the department-name map.

### Code

```go
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
```

### Dry Run

Example 1. For each employee we list the distinct higher salaries **within their department**, then decide.

| Employee | Dept | Salary | Distinct higher salaries in dept | count | count < 3? | High earner? |
|----------|------|--------|----------------------------------|-------|------------|--------------|
| Joe (IT) | 1 | 85000 | {90000} | 1 | yes | ✅ |
| Henry (Sales) | 2 | 80000 | {} | 0 | yes | ✅ |
| Sam (Sales) | 2 | 60000 | {80000} | 1 | yes | ✅ |
| Max (IT) | 1 | 90000 | {} | 0 | yes | ✅ |
| Janet (IT) | 1 | 69000 | {70000, 85000, 90000} | 3 | no | ❌ |
| Randy (IT) | 1 | 85000 | {90000} | 1 | yes | ✅ |
| Will (IT) | 1 | 70000 | {85000, 90000} | 2 | yes | ✅ |

Note Janet sees three distinct higher salaries (70000, 85000, 90000), so she is excluded — matching the expected output.

---

## Approach 2 — Sort and Group

### Intuition

Per department, the entire membership question hinges on one number: the **third-highest distinct salary**. Everyone at or above that threshold is in; everyone below is out. Departments with fewer than three distinct salaries have no cut-off — every employee qualifies. So collect each department's distinct salaries, sort them descending, and read off the third.

### Algorithm

1. Group **distinct** salaries per department into a set (`departmentId → set`).
2. For each department, dump the set into a slice, sort descending, and record `threshold[dept]` = the 3rd value (or `math.MinInt` when the department has fewer than 3 distinct salaries).
3. Filter pass: emit every employee with `e.Salary >= threshold[dept]`.

### Complexity

- **Time:** O(n log n) — grouping is O(n); sorting all distinct salaries is at worst O(n log n) across departments; the filter pass is O(n).
- **Space:** O(n) — the per-department distinct-salary sets.

### Code

```go
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
```

### Dry Run

Example 1.

Step 1 — distinct salaries per department:

| Dept | Distinct salaries collected |
|------|-----------------------------|
| 1 (IT) | {85000, 90000, 69000, 85000→dup, 70000} = {90000, 85000, 70000, 69000} |
| 2 (Sales) | {80000, 60000} |

Step 2 — sort descending and pick the 3rd:

| Dept | Sorted descending | len ≥ 3? | threshold |
|------|-------------------|----------|-----------|
| 1 (IT) | [90000, 85000, 70000, 69000] | yes | 70000 (index 2) |
| 2 (Sales) | [80000, 60000] | no | math.MinInt |

Step 3 — filter every employee against `threshold[dept]`:

| Employee | Dept | Salary | threshold | Salary ≥ threshold? |
|----------|------|--------|-----------|---------------------|
| Joe | 1 | 85000 | 70000 | ✅ |
| Henry | 2 | 80000 | MinInt | ✅ |
| Sam | 2 | 60000 | MinInt | ✅ |
| Max | 1 | 90000 | 70000 | ✅ |
| Janet | 1 | 69000 | 70000 | ❌ |
| Randy | 1 | 85000 | 70000 | ✅ |
| Will | 1 | 70000 | 70000 | ✅ |

---

## Approach 3 — Min-Heap Top-K

### Intuition

"Top three distinct salaries" is a top-K problem with K = 3. The standard tool is a min-heap of size K: its root is the **weakest** member of the current top K, so any new distinct salary bigger than the root evicts it. A companion set per department enforces distinctness. Unlike full sorting, this generalizes to any K at O(n log K) and never re-sorts.

### Algorithm

1. Stream the employees; for each department keep a `(heap, set)` pair.
2. Skip salaries already in the set (distinct only). If the heap holds fewer than 3, push unconditionally; otherwise, if `salary > root`, pop the root (and drop it from the set) and push the new salary.
3. Filter pass: an employee qualifies iff their department's heap holds fewer than 3 salaries (no cut-off yet) or `e.Salary >= heap root`.

### Complexity

- **Time:** O(n log K + d) = O(n) for K = 3 — each row does O(log 3) = O(1) heap work.
- **Space:** O(d) — at most 3 heap slots and 3 set entries per department.

### Code

```go
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
```

The heap itself is a standard `container/heap` min-heap over `int`:

```go
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
```

### Dry Run

Example 1 — build pass, tracing the IT department's size-3 min-heap. (Sales only ever sees two distinct salaries, so its heap stays under-full: `[60000, 80000]`.)

| Row | Salary | In set already? | Heap action | Heap (root first) | Set |
|-----|--------|-----------------|-------------|-------------------|-----|
| Joe | 85000 | no | len 0 < 3 → push | [85000] | {85000} |
| Max | 90000 | no | len 1 < 3 → push | [85000, 90000] | {85000, 90000} |
| Janet | 69000 | no | len 2 < 3 → push | [69000, 90000, 85000] | {69000, 85000, 90000} |
| Randy | 85000 | yes | duplicate → skip | [69000, 90000, 85000] | {69000, 85000, 90000} |
| Will | 70000 | no | full; 70000 > root 69000 → evict 69000, push 70000 | [70000, 90000, 85000] | {70000, 85000, 90000} |

Final IT heap root = **70000** (the 3rd-highest distinct salary). Filter pass:

| Employee | Dept | Salary | root / under-full | Qualifies? |
|----------|------|--------|-------------------|------------|
| Joe | IT | 85000 | ≥ 70000 | ✅ |
| Max | IT | 90000 | ≥ 70000 | ✅ |
| Janet | IT | 69000 | < 70000 | ❌ |
| Randy | IT | 85000 | ≥ 70000 | ✅ |
| Will | IT | 70000 | ≥ 70000 | ✅ |
| Henry | Sales | 80000 | heap under-full (len 2 < 3) | ✅ |
| Sam | Sales | 60000 | heap under-full (len 2 < 3) | ✅ |

---

## Approach 4 — One-Pass Top-3 Tracking (Optimal)

### Intuition

Since K is the constant 3, the heap is overkill. Three ordered slots per department — held in a `[3]int` sorted descending, with `math.MinInt` marking empty slots — track the top three distinct salaries exactly, updated with constant-time compare-and-shift inserts. After one build pass, slot `t[2]` is the department's cut-off (or `MinInt` when fewer than 3 distinct salaries exist), and a second pass keeps every employee at or above it.

### Algorithm

1. Stream the employees; per department keep `t = [3]int`, descending, with `math.MinInt` in empty slots.
2. Insert each salary with `insertDistinctTop3`: ignore it if it equals any slot (distinct only); otherwise shift-and-place it into rank order.
3. Filter pass: emit every employee with `e.Salary >= t[2]` — `MinInt` makes departments with fewer than 3 distinct salaries accept everyone.

### Complexity

- **Time:** O(n + d) — O(1) insert work per row, then a linear filter pass.
- **Space:** O(d) — exactly three tracked salaries per department.

### Code

```go
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
```

### Dry Run

Example 1 — build pass for the IT department. `t` starts as `[MinInt, MinInt, MinInt]`.

| Row | Salary | Case matched | `t` after (desc) |
|-----|--------|--------------|------------------|
| Joe | 85000 | s > t[0] (85000 > MinInt) | [85000, MinInt, MinInt] |
| Max | 90000 | s > t[0] (90000 > 85000) | [90000, 85000, MinInt] |
| Janet | 69000 | s > t[2] (69000 > MinInt) | [90000, 85000, 69000] |
| Randy | 85000 | s == t[1] → duplicate, no change | [90000, 85000, 69000] |
| Will | 70000 | s > t[2] (70000 > 69000) | [90000, 85000, 70000] |

Final IT slots = `[90000, 85000, 70000]`, so cut-off `t[2] = 70000`. Sales sees only `[80000, 60000, MinInt]`, cut-off `MinInt`. Filter pass:

| Employee | Dept | Salary | t[2] (cut-off) | Salary ≥ t[2]? |
|----------|------|--------|----------------|----------------|
| Joe | IT | 85000 | 70000 | ✅ |
| Max | IT | 90000 | 70000 | ✅ |
| Janet | IT | 69000 | 70000 | ❌ |
| Randy | IT | 85000 | 70000 | ✅ |
| Will | IT | 70000 | 70000 | ✅ |
| Henry | Sales | 80000 | MinInt | ✅ |
| Sam | Sales | 60000 | MinInt | ✅ |

All four approaches produce the same six high earners.

---

## Key Takeaways

- **"Top K distinct per group" reduces to one number per group: the K-th-highest distinct value (the cut-off).** Compute it once, then a single linear filter (`value >= cutoff`) answers membership for every row. The SQL `DENSE_RANK() OVER (PARTITION BY dept ORDER BY salary DESC) <= 3` says exactly this.
- **The correlated-subquery framing — "count how many DISTINCT values strictly exceed mine; qualify iff that count < K" — is the direct brute-force translation** of the SQL, and worth recognizing on sight even when you plan to optimize past its O(n²).
- **A size-K min-heap is the general top-K workhorse:** its root is the weakest survivor and the exact cut-off. But when **K is a small constant**, fixed ordered slots with compare-and-shift inserts drop the log factor and the allocation overhead entirely.
- **Distinctness is a set concern, kept orthogonal to ranking.** Every approach layers a `map[int]bool` (or an equality check against tracked slots) on top of its ranking logic so duplicate salaries collapse before they compete.
- **Sentinel `math.MinInt` for "empty slot" / "no cut-off"** lets the same `salary >= cutoff` comparison handle both full departments and departments with fewer than three distinct salaries — no special-casing in the filter loop.

---

## Related Problems

- LeetCode #184 — Department Highest Salary (same partition-by-department pattern, top-1 instead of top-3)
- LeetCode #176 — Second Highest Salary (K-th distinct salary, single group)
- LeetCode #177 — Nth Highest Salary (generalized K-th distinct salary)
- LeetCode #215 — Kth Largest Element in an Array (min-heap / quickselect top-K core)
- LeetCode #347 — Top K Frequent Elements (top-K by a per-group aggregate)
- LeetCode #692 — Top K Frequent Words (top-K with tie-breaking)
