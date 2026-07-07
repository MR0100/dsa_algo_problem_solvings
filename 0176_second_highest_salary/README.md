# 0176 — Second Highest Salary

> LeetCode #176 · Difficulty: Medium
> **Categories:** Database, Sorting, Top-K Tracking

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

Write a solution to find the **second highest distinct salary** from the `Employee` table. If there is no second highest salary, return `null` (return `None` in Pandas).

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
Output:
+---------------------+
| SecondHighestSalary |
+---------------------+
| 200                 |
+---------------------+
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
Output:
+---------------------+
| SecondHighestSalary |
+---------------------+
| null                |
+---------------------+
```

> **Note:** This is a Database problem. Per this repo's convention everything is
> solved in **Go**: the table is modeled as `[]Employee{ID, Salary}` and each
> approach implements the query plan by hand, returning `*int` so a missing
> answer maps to `nil` (SQL `NULL`). The canonical SQL appears in Key Takeaways.

---

## Company Frequency

| Company    | Frequency        | Last Reported |
|------------|------------------|---------------|
| Amazon     | ★★★★★ Very High  | 2024          |
| Microsoft  | ★★★★☆ High       | 2024          |
| Google     | ★★★☆☆ Medium     | 2024          |
| Adobe      | ★★★☆☆ Medium     | 2023          |
| Apple      | ★★☆☆☆ Low        | 2023          |
| Bloomberg  | ★★☆☆☆ Low        | 2023          |

> ⚠️ Frequency data is crowd-sourced from LeetCode Discuss, Glassdoor, and
> community interview reports. Treat as a signal, not a guarantee.

---

## DSA Concepts Used

- **Hash Map / Hash Set** — deduplicating salaries is the Go equivalent of SQL `DISTINCT` → see [`/dsa/hash_map.md`](/dsa/hash_map.md)
- **Sorting** — "second highest" is offset 1 of the distinct values sorted descending (`ORDER BY ... DESC LIMIT 1 OFFSET 1`) → see [`/dsa/sorting.md`](/dsa/sorting.md)
- **Greedy / Streaming Scan** — the optimal answer maintains the top-2 distinct values in one pass with O(1) state, the same running-extrema pattern as #414 Third Maximum Number → see [`/dsa/greedy.md`](/dsa/greedy.md)

---

## Approaches Overview

| # | Approach | Time | Space | When to use |
|---|----------|------|-------|-------------|
| 1 | Brute Force (Sort Distinct Salaries) | O(n log n) | O(n) | Mirrors the SQL `DISTINCT + ORDER BY + OFFSET` plan; trivially generalises to *n-th* highest |
| 2 | Two-Pass Max | O(n) | O(1) | Mirrors the SQL `MAX(salary) WHERE salary < MAX` sub-query; no sorting needed |
| 3 | Single Pass Top-2 Tracking (Optimal) | O(n) | O(1) | Streaming data / one scan allowed only; the interview-favourite form |

---

## Approach 1 — Brute Force (Sort Distinct Salaries)

### Intuition

"Second highest **distinct** salary" is literally "the element at offset 1 of the distinct salaries sorted descending". So do exactly that: deduplicate with a hash set, sort descending, index position 1. If fewer than two distinct salaries survive the dedup, the answer is `NULL`. This is a one-to-one translation of `SELECT DISTINCT salary FROM Employee ORDER BY salary DESC LIMIT 1 OFFSET 1`.

### Algorithm

1. Walk the table; insert each salary into a hash set, appending first occurrences to a `distinct` slice.
2. Sort `distinct` in descending order.
3. If `len(distinct) < 2`, return `nil` (SQL `NULL`).
4. Otherwise return `distinct[1]`.

### Complexity

- **Time:** O(n log n) — the descending sort dominates the O(n) dedup pass.
- **Space:** O(n) — the hash set and the distinct-values slice can hold all n salaries.

### Code

```go
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
```

### Dry Run

Example 1: `Employee = [(1,100), (2,200), (3,300)]`

| Step | Row (id, salary) | seen after | distinct after | Action |
|------|------------------|------------|----------------|--------|
| 1 | (1, 100) | {100} | [100] | new value → collect |
| 2 | (2, 200) | {100, 200} | [100, 200] | new value → collect |
| 3 | (3, 300) | {100, 200, 300} | [100, 200, 300] | new value → collect |
| 4 | sort descending | — | [300, 200, 100] | highest first |
| 5 | len check | — | 3 ≥ 2 | index 1 exists |
| 6 | return `distinct[1]` | — | — | **200** |

Result: `200` ✔ (Example 2: distinct = [100], len 1 < 2 → `null` ✔)

---

## Approach 2 — Two-Pass Max

### Intuition

The second highest distinct salary is exactly `MAX(salary)` over the rows whose salary is **strictly below** the global maximum. So: one pass to find the maximum, a second pass to find the best value under it. Duplicates of the maximum are excluded automatically by the strict `<` — no explicit dedup structure needed. This is the Go translation of the classic `SELECT MAX(salary) FROM Employee WHERE salary < (SELECT MAX(salary) FROM Employee)`, which handles the `NULL` case for free in SQL (an empty `MAX` is `NULL`).

### Algorithm

1. Pass 1: scan every row, tracking the maximum salary `max1`.
2. Pass 2: scan again, tracking the largest salary strictly less than `max1` plus a `found` flag.
3. If nothing was found (all rows share the same salary), return `nil`.
4. Otherwise return the runner-up value.

### Complexity

- **Time:** O(n) — two linear scans of the table.
- **Space:** O(1) — two scalars and a boolean flag.

### Code

```go
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
```

### Dry Run

Example 1: `Employee = [(1,100), (2,200), (3,300)]`

| Step | Pass | Row salary | max1 | best | found | Action |
|------|------|-----------|------|------|-------|--------|
| 1 | 1 | 100 | 100 | — | — | init / not greater |
| 2 | 1 | 200 | 200 | — | — | 200 > 100 → update |
| 3 | 1 | 300 | 300 | — | — | 300 > 200 → update |
| 4 | 2 | 100 | 300 | 100 | true | 100 < 300 and slot empty → take |
| 5 | 2 | 200 | 300 | 200 | true | 200 < 300 and 200 > 100 → improve |
| 6 | 2 | 300 | 300 | 200 | true | 300 < 300 fails → skip (the max itself) |
| 7 | — | return | — | 200 | — | **200** |

Result: `200` ✔ (Example 2: pass 2 finds nothing below 100 → `found=false` → `null` ✔)

---

## Approach 3 — Single Pass Top-2 Tracking (Optimal)

### Intuition

Fuse the two passes into one by carrying **two running slots**: `first` (highest distinct so far) and `second` (highest distinct strictly below `first`). Each incoming salary either dethrones `first` (the old `first` slides down into `second`), slots in between, or is ignored. A duplicate of `first` satisfies neither `s > first` nor `s < first`, so it falls through untouched — the DISTINCT requirement is enforced by the strict comparisons alone. This is the streaming form of Approach 2 and the exact pattern of #414 Third Maximum Number generalised to k = 2.

### Algorithm

1. Initialise `first = second = -∞` (`math.MinInt` as the "slot empty" sentinel).
2. For each salary `s`:
   1. If `s > first`: `second = first`, then `first = s`.
   2. Else if `s < first && s > second`: `second = s`.
   3. Else (`s == first` or `s ≤ second`): do nothing.
3. After the scan, if `second` still holds the sentinel there is no second distinct salary → return `nil`; otherwise return `second`.

### Complexity

- **Time:** O(n) — exactly one scan, O(1) work per row.
- **Space:** O(1) — two integer slots regardless of table size.

### Code

```go
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
```

### Dry Run

Example 1: `Employee = [(1,100), (2,200), (3,300)]` (−∞ = `math.MinInt`)

| Step | s | s > first? | s < first && s > second? | first | second |
|------|-----|-----------|--------------------------|-------|--------|
| 0 | — | — | — | −∞ | −∞ |
| 1 | 100 | yes | — | 100 | −∞ |
| 2 | 200 | yes | — | 200 | 100 |
| 3 | 300 | yes | — | 300 | 200 |
| 4 | return | — | second ≠ −∞ | — | **200** |

Result: `200` ✔ (Example 2: after 100, `second` stays −∞ → `null` ✔. Edge `[200,200,100]`: the second 200 matches neither branch, then 100 fills `second` → `100` ✔)

---

## Key Takeaways

- **"k-th highest DISTINCT" ⇒ dedupe first, then rank.** Strict inequalities (`<`, `>`) are the cheapest dedup tool: they make duplicates of a tracked value invisible without any set.
- **Top-k tracking with k slots** — for small fixed k, k running variables beat sorting: O(n) time, O(1) space, works on streams. k = 2 here, k = 3 in #414 Third Maximum Number.
- **Nullable results in Go** — return `*int` (or `(int, bool)`) to model SQL `NULL`; never overload a real value like `-1` when the domain doesn't forbid it.
- **Sentinel caveat** — `math.MinInt` as "empty" is safe only because salaries can't be `math.MinInt`; when the full domain is possible, carry an explicit `found` flag (as Approach 2 does).
- Canonical SQL, both forms:
  ```sql
  -- Form 1: OFFSET, wrapped so an empty result becomes NULL
  SELECT (SELECT DISTINCT salary FROM Employee
          ORDER BY salary DESC LIMIT 1 OFFSET 1) AS SecondHighestSalary;
  -- Form 2: MAX below MAX (empty MAX is already NULL)
  SELECT MAX(salary) AS SecondHighestSalary FROM Employee
  WHERE salary < (SELECT MAX(salary) FROM Employee);
  ```

---

## Related Problems

- LeetCode #177 — Nth Highest Salary (this problem generalised from 2 to n)
- LeetCode #184 — Department Highest Salary (MAX per group instead of global)
- LeetCode #185 — Department Top Three Salaries (top-k distinct per group)
- LeetCode #414 — Third Maximum Number (identical top-k distinct tracking on an array)
- LeetCode #215 — Kth Largest Element in an Array (rank selection without the distinct twist)
