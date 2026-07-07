# 0197 — Rising Temperature

> LeetCode #197 · Difficulty: Easy
> **Categories:** Database, Hash Map, Sorting

---

## Problem Statement

Table: `Weather`

```
+---------------+---------+
| Column Name   | Type    |
+---------------+---------+
| id            | int     |
| recordDate    | date    |
| temperature   | int     |
+---------------+---------+
```

`id` is the column with unique values for this table.
There are no different rows with the same `recordDate`.
This table contains information about the temperature on a certain day.

Write a solution to find all dates' `id` with higher temperatures compared to its previous dates (yesterday).

Return the result table in **any order**.

The result format is in the following example.

**Example 1**
```
Input:
Weather table:
+----+------------+-------------+
| id | recordDate | temperature |
+----+------------+-------------+
| 1  | 2015-01-01 | 10          |
| 2  | 2015-01-02 | 25          |
| 3  | 2015-01-03 | 20          |
| 4  | 2015-01-04 | 30          |
+----+------------+-------------+
Output:
+----+
| id |
+----+
| 2  |
| 4  |
+----+
Explanation:
In 2015-01-02, the temperature was higher than the previous day (10 -> 25).
In 2015-01-04, the temperature was higher than the previous day (20 -> 30).
```

**Constraints**
- Every `id` is unique.
- No two rows share the same `recordDate` (at most one reading per day).
- Dates are valid calendar dates; **they are not guaranteed to be consecutive** — a row only qualifies when a row dated exactly one day earlier exists.

> ℹ️ This is a **Database** problem. Per this repo's Go-only convention, `main.go`
> models the `Weather` table as a slice of structs and implements the query logic
> with classic DSA techniques; the canonical SQL answers appear in
> [Key Takeaways](#key-takeaways).

---

## Company Frequency

| Company    | Frequency        | Last Reported |
|------------|------------------|---------------|
| Amazon     | ★★★★☆ High       | 2024          |
| Google     | ★★★☆☆ Medium     | 2024          |
| Microsoft  | ★★★☆☆ Medium     | 2023          |
| Apple      | ★★☆☆☆ Low        | 2023          |
| Bloomberg  | ★★☆☆☆ Low        | 2023          |

> ⚠️ Frequency data is crowd-sourced from LeetCode Discuss, Glassdoor, and
> community interview reports. Treat as a signal, not a guarantee.

---

## DSA Concepts Used

- **Hash Map** — indexing rows by date turns "find the row dated yesterday" from an O(n) scan into an O(1) lookup; it is exactly the hash join a database engine would run. → see [`/dsa/hash_map.md`](/dsa/hash_map.md)
- **Sorting** — chronological order makes each row's only possible "yesterday" its immediate predecessor, mirroring the SQL `LAG()` window function. → see [`/dsa/sorting.md`](/dsa/sorting.md)

---

## Approaches Overview

| # | Approach | Time | Space | When to use |
|---|----------|------|-------|-------------|
| 1 | Brute Force (pairwise scan) | O(n²) | O(1) aux | Tiny tables; mirrors the SQL self-join literally |
| 2 | Sort + Scan | O(n log n) | O(n) | Data naturally processed in date order; LAG-style reasoning |
| 3 | Hash Map (Optimal) ✅ | O(n) | O(n) | General case — constant-time yesterday lookups |

---

## Approach 1 — Brute Force

### Intuition
For each day we need the reading taken **exactly one calendar day earlier**. Without any index, just scan the whole table for it. This is the nested-loop execution of the SQL self-join:

```sql
SELECT w1.id
FROM Weather w1, Weather w2
WHERE DATEDIFF(w1.recordDate, w2.recordDate) = 1
  AND w1.temperature > w2.temperature;
```

The one subtlety of this whole problem: "previous date" means *calendar* yesterday, not the previous row — ids and dates can have gaps, and month/year boundaries (`2015-02-01` after `2015-01-31`) must be handled by real date arithmetic, never by `id - 1` or day-number arithmetic.

### Algorithm
1. For every row `today`, scan all rows for one whose date is exactly one day before `today` (checked with `AddDate(0, 0, 1)` — calendar-safe).
2. If that yesterday row exists **and** `today.Temperature > yesterday.Temperature`, collect `today.ID`. Dates are unique, so at most one match exists — stop scanning at the first.
3. Sort the collected ids (any order is accepted; sorting makes output deterministic).

### Complexity
- **Time:** O(n²) — every row scans the entire table once.
- **Space:** O(1) auxiliary — nothing allocated beyond the output slice.

### Code
```go
func bruteForce(weather []Weather) []int {
	ids := []int{}
	for _, today := range weather {
		for _, other := range weather {
			// other qualifies only when dated exactly one day before today.
			if isNextDay(other.RecordDate, today.RecordDate) {
				if today.Temperature > other.Temperature {
					ids = append(ids, today.ID)
				}
				break // dates are unique → at most one "yesterday" exists
			}
		}
	}
	sort.Ints(ids) // deterministic output
	return ids
}
```

### Dry Run — Example 1

| today (id, date, temp) | yesterday row found?          | temp rise?    | ids after step |
|------------------------|-------------------------------|---------------|----------------|
| (1, 2015-01-01, 10)    | none — no 2014-12-31 row      | —             | [] |
| (2, 2015-01-02, 25)    | (1, 2015-01-01, 10)           | 25 > 10 ✅    | [2] |
| (3, 2015-01-03, 20)    | (2, 2015-01-02, 25)           | 20 > 25 ❌    | [2] |
| (4, 2015-01-04, 30)    | (3, 2015-01-03, 20)           | 30 > 20 ✅    | [2, 4] |

Result → `[2 4]` ✓

---

## Approach 2 — Sort + Scan

### Intuition
After sorting chronologically, a row's only possible "yesterday" is the row **directly before it** — any earlier row is at least two days away. This is precisely how the window-function solution thinks:

```sql
SELECT id FROM (
    SELECT id, temperature, recordDate,
           LAG(temperature) OVER (ORDER BY recordDate) AS prevTemp,
           LAG(recordDate)  OVER (ORDER BY recordDate) AS prevDate
    FROM Weather
) t
WHERE temperature > prevTemp AND DATEDIFF(recordDate, prevDate) = 1;
```

…including its classic pitfall: adjacent rows in date order may still be **more than one day apart**, so the calendar gap must be verified explicitly.

### Algorithm
1. Copy the table (never reorder the caller's data) and sort it by `recordDate` ascending — `"YYYY-MM-DD"` strings sort lexicographically in chronological order, so plain string comparison is correct.
2. For each `i ≥ 1`, check **both** conditions: row `i-1` is exactly one calendar day earlier, and `rows[i].Temperature > rows[i-1].Temperature`. Collect `rows[i].ID` when both hold.
3. Sort the collected ids.

### Complexity
- **Time:** O(n log n) — the sort dominates; the adjacent scan is O(n).
- **Space:** O(n) — the sorted working copy.

### Code
```go
func sortAndScan(weather []Weather) []int {
	rows := make([]Weather, len(weather))
	copy(rows, weather) // never reorder the caller's table

	// "YYYY-MM-DD" strings sort lexicographically in chronological order,
	// so a plain string comparison is a correct date comparison here.
	sort.Slice(rows, func(i, j int) bool { return rows[i].RecordDate < rows[j].RecordDate })

	ids := []int{}
	for i := 1; i < len(rows); i++ {
		// Both conditions required: consecutive calendar days AND a rise.
		if isNextDay(rows[i-1].RecordDate, rows[i].RecordDate) &&
			rows[i].Temperature > rows[i-1].Temperature {
			ids = append(ids, rows[i].ID)
		}
	}
	sort.Ints(ids) // deterministic output
	return ids
}
```

### Dry Run — Example 1

**After sorting by date** (input already chronological): `[(1, 01-01, 10), (2, 01-02, 25), (3, 01-03, 20), (4, 01-04, 30)]`

| i | rows[i-1] → rows[i]            | one day apart? | temp rise? | ids after step |
|---|--------------------------------|----------------|------------|----------------|
| 1 | (01-01, 10) → (01-02, 25)      | ✅             | 25 > 10 ✅ | [2] |
| 2 | (01-02, 25) → (01-03, 20)      | ✅             | 20 > 25 ❌ | [2] |
| 3 | (01-03, 20) → (01-04, 30)      | ✅             | 30 > 20 ✅ | [2, 4] |

Result → `[2 4]` ✓

---

## Approach 3 — Hash Map (Optimal)

### Intuition
The brute force wastes its inner loop searching for yesterday. Index every row by its date string once; then each row answers "who was yesterday?" with a single O(1) map lookup. This is exactly the **hash join** a database engine picks for the self-join above.

### Algorithm
1. **Pass 1:** build `byDate: recordDate → row` for every row.
2. **Pass 2:** for each row, compute yesterday's date string with `AddDate(0, 0, -1)` (calendar-safe across months, years, and leap days) and look it up. If present and today's temperature is strictly higher, collect the id.
3. Sort the collected ids.

### Complexity
- **Time:** O(n) — two linear passes, O(1) work per row.
- **Space:** O(n) — one map entry per row.

### Code
```go
func hashMap(weather []Weather) []int {
	// Index every row by its date for O(1) "who was yesterday?" lookups.
	byDate := make(map[string]Weather, len(weather))
	for _, row := range weather {
		byDate[row.RecordDate] = row
	}

	ids := []int{}
	for _, today := range weather {
		// Calendar-correct yesterday (handles month/year edges, leap days).
		yesterday := mustDate(today.RecordDate).AddDate(0, 0, -1).Format(dateLayout)
		if prev, ok := byDate[yesterday]; ok && today.Temperature > prev.Temperature {
			ids = append(ids, today.ID)
		}
	}
	sort.Ints(ids) // deterministic output
	return ids
}
```

### Dry Run — Example 1

**Pass 1 — index by date:** `byDate = {01-01: (1,10), 01-02: (2,25), 01-03: (3,20), 01-04: (4,30)}`

**Pass 2 — lookup yesterday:**

| today (id, date, temp) | yesterday key | in map?        | temp rise? | ids after step |
|------------------------|---------------|----------------|------------|----------------|
| (1, 01-01, 10)         | 2014-12-31    | ❌             | —          | [] |
| (2, 01-02, 25)         | 2015-01-01    | ✅ (1, 10)     | 25 > 10 ✅ | [2] |
| (3, 01-03, 20)         | 2015-01-02    | ✅ (2, 25)     | 20 > 25 ❌ | [2] |
| (4, 01-04, 30)         | 2015-01-03    | ✅ (3, 20)     | 30 > 20 ✅ | [2, 4] |

Result → `[2 4]` ✓

---

## Key Takeaways

- **"Compare with yesterday" = self-join on a calendar condition.** Whenever a row must be matched with an offset version of its own table (yesterday, previous event, manager of employee), think self-join — and in code, think hash map keyed by the join column.
- **Never do date math with `id - 1` or raw subtraction.** Ids can have gaps and dates cross month/year boundaries; only calendar-aware arithmetic (`DATEDIFF` in SQL, `AddDate` in Go) is correct. `recordDate - 1` style arithmetic on date types is a classic wrong answer here.
- **The canonical SQL answers**, both accepted on LeetCode:
  ```sql
  -- Self-join with DATEDIFF (most common interview answer)
  SELECT w1.id
  FROM Weather w1
  JOIN Weather w2 ON DATEDIFF(w1.recordDate, w2.recordDate) = 1
  WHERE w1.temperature > w2.temperature;

  -- Window-function version: LAG must also verify the 1-day gap!
  SELECT id FROM (
      SELECT id, temperature, recordDate,
             LAG(temperature) OVER (ORDER BY recordDate) AS prevTemp,
             LAG(recordDate)  OVER (ORDER BY recordDate) AS prevDate
      FROM Weather
  ) t
  WHERE temperature > prevTemp AND DATEDIFF(recordDate, prevDate) = 1;
  ```
- **ISO-8601 dates (`YYYY-MM-DD`) sort lexicographically** — fixed-width, most-significant-first. Plain string comparison is a valid chronological sort; a handy trick for logs, filenames, and keys.
- **Hash join beats nested loop.** The O(n²) → O(n) jump here is the same "index the lookup side" move as Two Sum: build a map on the join key, then probe it.

---

## Related Problems

- LeetCode #196 — Delete Duplicate Emails (database table modelled with the same in-memory technique)
- LeetCode #181 — Employees Earning More Than Their Managers (self-join comparing a row with its related row)
- LeetCode #182 — Duplicate Emails (GROUP BY / hash-map aggregation)
- LeetCode #1661 — Average Time of Process per Machine (pairing related rows of one table)
- LeetCode #1 — Two Sum (the same "index one side, probe with the other" hash pattern in array form)
