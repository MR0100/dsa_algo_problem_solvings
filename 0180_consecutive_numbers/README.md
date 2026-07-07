# 0180 — Consecutive Numbers

> LeetCode #180 · Difficulty: Medium
> **Categories:** Database, SQL, Self-Join, Window Scan, Sliding Window

---

## Problem Statement

Table: `Logs`

```
+-------------+---------+
| Column Name | Type    |
+-------------+---------+
| id          | int     |
| num         | varchar |
+-------------+---------+
```

In SQL, `id` is the primary key for this table. `id` is an autoincrement
column starting from `1`.

Find all numbers that appear **at least three times consecutively**.

Return the result table in **any order**.

The result format is in the following example.

### Example 1

**Input:**

`Logs` table:

```
+----+-----+
| id | num |
+----+-----+
| 1  | 1   |
| 2  | 1   |
| 3  | 1   |
| 4  | 2   |
| 5  | 1   |
| 6  | 2   |
| 7  | 2   |
+----+-----+
```

**Output:**

```
+-----------------+
| ConsecutiveNums |
+-----------------+
| 1               |
+-----------------+
```

**Explanation:** `1` is the only number that appears consecutively for at least
three times (rows with `id` 1, 2 and 3).

### Constraints

- The rows of the table are ordered by `id` ascending, which is a sequential
  autoincrement starting from `1` (may have gaps in general, but the run must be
  over consecutive `id`s to count).

---

## Company Frequency

| Company    | Frequency        | Last Reported |
|------------|------------------|---------------|
| Amazon     | ★★★★☆ High       | 2024          |
| Meta       | ★★★★☆ High       | 2023          |
| Google     | ★★★☆☆ Medium     | 2023          |
| Microsoft  | ★★★☆☆ Medium     | 2023          |
| Apple      | ★★☆☆☆ Low        | 2022          |
| Bloomberg  | ★★☆☆☆ Low        | 2022          |

> ⚠️ Frequency data is crowd-sourced from LeetCode Discuss, Glassdoor, and
> community interview reports. Treat as a signal, not a guarantee.

---

## DSA Concepts Used

- **Hash Map / Set** — a `map[int]bool` deduplicates the qualifying values,
  exactly reproducing SQL's `SELECT DISTINCT` → see [`/dsa/hash_map.md`](/dsa/hash_map.md)
- **Sliding Window** — the optimal scan slides a fixed size-3 window (or a
  running streak counter) across the id-ordered rows → see [`/dsa/sliding_window.md`](/dsa/sliding_window.md)
- **Two Pointers** — the adjacent-triple check compares neighbouring rows,
  the same "walk in lockstep over ordered data" pattern → see [`/dsa/two_pointers.md`](/dsa/two_pointers.md)
- **Sorting** — rows are ordered by `id` so that a "consecutive run" becomes a
  block of physically adjacent entries → see [`/dsa/sorting.md`](/dsa/sorting.md)

---

## Approaches Overview

| # | Approach | Time | Space | When to use |
|---|----------|------|-------|-------------|
| 1 | Self-Join Simulation (Brute Force) | O(n³) | O(d) | Literal translation of the SQL self-join; clarity over speed |
| 2 | Adjacent-Triple Scan (Sorted Window) | O(n log n) | O(n) | Rows not guaranteed id-sorted; still want a simple window test |
| 3 | Single-Pass Running-Streak Counter (Optimal) | O(n log n) → O(n) if pre-sorted | O(1) extra | Best in practice; one linear pass over id-ordered rows |

*(n = number of rows, d = number of distinct qualifying values.)*

---

## Approach 1 — Self-Join Simulation (Brute Force)

### Intuition
The accepted SQL joins three aliases of the same table, `Logs l1, l2, l3`,
whose ids form a run `(id, id+1, id+2)` and whose `num`s are all equal. We
reproduce that join literally: examine every triple of rows and keep the shared
value whenever the id-chain and value-match predicates both hold. A `DISTINCT`
set collapses repeated hits.

### Algorithm
1. For every ordered triple of rows `(a, b, c)`:
2. Check the id chain: `b.ID == a.ID+1` and `c.ID == b.ID+1`.
3. Check the value match: `a.Num == b.Num == c.Num`.
4. When all predicates hold, add `a.Num` to a DISTINCT set.
5. Return the set's members, sorted ascending.

### Complexity
- **Time:** O(n³) — three nested passes over the rows, mirroring the three
  table aliases of the self-join.
- **Space:** O(d) — the DISTINCT set of qualifying values.

### Code
```go
func selfJoin(logs []Log) []int {
	found := map[int]bool{} // DISTINCT l1.num that satisfy the join
	for i := range logs {   // l1
		for j := range logs { // l2
			for k := range logs { // l3
				a, b, c := logs[i], logs[j], logs[k]
				// id chain: consecutive ids l1.id = l2.id-1, l2.id = l3.id-1
				if b.ID != a.ID+1 || c.ID != b.ID+1 {
					continue
				}
				// value match: l1.num = l2.num = l3.num
				if a.Num == b.Num && b.Num == c.Num {
					found[a.Num] = true // DISTINCT collapses repeats for free
				}
			}
		}
	}
	return sortedKeys(found)
}
```

### Dry Run
Input rows (id, num): (1,1) (2,1) (3,1) (4,2) (5,1) (6,2) (7,2). Only triples
whose ids chain as `x, x+1, x+2` are tested — one per starting id.

| Triple ids | nums     | id chain ok? | nums equal? | Action           |
|------------|----------|--------------|-------------|------------------|
| 1,2,3      | 1,1,1    | yes          | yes         | add `1` to set   |
| 2,3,4      | 1,1,2    | yes          | no          | skip             |
| 3,4,5      | 1,2,1    | yes          | no          | skip             |
| 4,5,6      | 2,1,2    | yes          | no          | skip             |
| 5,6,7      | 1,2,2    | yes          | no          | skip             |

Set = `{1}` → sorted result **`[1]`**.

---

## Approach 2 — Adjacent-Triple Scan (Sorted Window)

### Intuition
The self-join only ever pairs rows whose ids differ by exactly 1, so a valid
triple must be three **physically adjacent** rows once the table is ordered by
id. Sort by id, then slide a size-3 window: a window qualifies when its ids are
consecutive and its three `num`s are equal. This collapses the O(n³) join into
a single linear sweep after the sort.

### Algorithm
1. Copy the rows and sort them ascending by id (do not mutate the input).
2. For each start index `t` from `0` to `n-3`, take rows `t, t+1, t+2`.
3. Require consecutive ids (which guards against id gaps) **and** all three
   `num`s equal.
4. Collect the shared `num` into a DISTINCT set; return it sorted.

### Complexity
- **Time:** O(n log n) — the id sort dominates; the window sweep is O(n).
- **Space:** O(n) — the sorted copy of the rows (plus the O(d) answer set).

### Code
```go
func adjacentTriple(logs []Log) []int {
	rows := append([]Log(nil), logs...) // sort a copy: leave input untouched
	sort.Slice(rows, func(i, j int) bool { return rows[i].ID < rows[j].ID })

	found := map[int]bool{}
	for t := 0; t+2 < len(rows); t++ { // window [t, t+1, t+2]
		a, b, c := rows[t], rows[t+1], rows[t+2]
		// consecutive ids: rejects any run interrupted by an id gap
		if b.ID != a.ID+1 || c.ID != b.ID+1 {
			continue
		}
		// three equal values in a row → a consecutive-3 number
		if a.Num == b.Num && b.Num == c.Num {
			found[a.Num] = true
		}
	}
	return sortedKeys(found)
}
```

### Dry Run
Already id-sorted: (1,1) (2,1) (3,1) (4,2) (5,1) (6,2) (7,2). Windows of 3:

| t | window (id:num)       | ids consecutive? | nums equal? | Action         |
|---|-----------------------|------------------|-------------|----------------|
| 0 | 1:1, 2:1, 3:1         | yes              | yes         | add `1`        |
| 1 | 2:1, 3:1, 4:2         | yes              | no          | skip           |
| 2 | 3:1, 4:2, 5:1         | yes              | no          | skip           |
| 3 | 4:2, 5:1, 6:2         | yes              | no          | skip           |
| 4 | 5:1, 6:2, 7:2         | yes              | no          | skip           |

Set = `{1}` → sorted result **`[1]`**.

---

## Approach 3 — Single-Pass Running-Streak Counter (Optimal)

### Intuition
Instead of re-reading a 3-row window at every step, keep a single **streak**
counter for the current value. Each row either extends the streak (same value
*and* an id exactly one greater than the previous id) or restarts it at 1. The
first time a streak reaches 3, that value qualifies — record it once. Two
counters replace the whole self-join.

### Algorithm
1. Copy the rows and sort ascending by id.
2. Track `prevNum`, `prevID`, and a `streak` counter.
3. For each row: if `num == prevNum` **and** `id == prevID+1`, increment
   `streak`; otherwise reset `streak = 1`.
4. When `streak` reaches exactly 3, mark `num` in the answer set (using `== 3`
   rather than `>= 3` records each run's value only once).
5. Update `prevNum`/`prevID`; return the set sorted ascending.

### Complexity
- **Time:** O(n log n) for the id sort; the streak scan itself is O(n). If the
  input is already id-sorted (as LeetCode guarantees), the whole thing is O(n).
- **Space:** O(1) extra beyond the sorted copy (plus the O(d) answer set).

### Code
```go
func runningStreak(logs []Log) []int {
	rows := append([]Log(nil), logs...) // don't disturb the caller's ordering
	sort.Slice(rows, func(i, j int) bool { return rows[i].ID < rows[j].ID })

	found := map[int]bool{}
	streak := 0                  // length of the current same-value run
	prevNum, prevID := 0, -1<<62 // sentinels: first row always restarts
	for _, r := range rows {
		if r.Num == prevNum && r.ID == prevID+1 {
			streak++ // same value AND contiguous id → run continues
		} else {
			streak = 1 // value changed or id gap → new run of length 1
		}
		if streak == 3 { // first time a run reaches 3 → qualifying value
			found[r.Num] = true // == 3 (not >=) records each value once per run
		}
		prevNum, prevID = r.Num, r.ID
	}
	return sortedKeys(found)
}
```

### Dry Run
Rows in id order: (1,1) (2,1) (3,1) (4,2) (5,1) (6,2) (7,2).

| row (id:num) | continues run? (num==prev && id==prev+1) | streak | streak==3? | set    |
|--------------|------------------------------------------|--------|------------|--------|
| 1:1          | no (sentinel)                            | 1      | no         | {}     |
| 2:1          | yes                                      | 2      | no         | {}     |
| 3:1          | yes                                      | 3      | **yes**    | {1}    |
| 4:2          | no (num changed)                         | 1      | no         | {1}    |
| 5:1          | no (num changed)                         | 1      | no         | {1}    |
| 6:2          | no (num changed)                         | 1      | no         | {1}    |
| 7:2          | yes                                      | 2      | no         | {1}    |

Set = `{1}` → sorted result **`[1]`**.

---

## Key Takeaways

- **"Consecutive N times" ⇒ think self-join or streak scan.** The textbook SQL
  answer chains N aliases on `id = id ± 1`; the linear in-memory analogue is a
  running-streak counter — the general shape for any "K in a row" question.
- **Always tie the run to consecutive ids, not just equal values.** Ids can
  have gaps; a value repeating across a gap is *not* consecutive. Every approach
  checks `id == prev+1` in addition to `num == prev`.
- **`SELECT DISTINCT` ↔ a hash set.** Multiple overlapping runs of the same
  value must yield a single output row, exactly what a `map[int]bool` gives you.
- **Fire on `streak == 3`, not `>= 3`.** Using `==` records each qualifying run
  exactly once, avoiding redundant set writes as a long run keeps growing.
- **The N-alias self-join generalises poorly** (O(n^N)); the streak counter is
  O(n) for any N, which is why interviewers like the follow-up "what if it's
  K consecutive?"

---

## Related Problems

- LeetCode #181 — Employees Earning More Than Their Managers (SQL self-join)
- LeetCode #182 — Duplicate Emails (SQL grouping / `HAVING COUNT >= 2`)
- LeetCode #196 — Delete Duplicate Emails (SQL self-join on adjacency)
- LeetCode #1454 — Active Users (consecutive-days run detection, same streak idea)
- LeetCode #601 — Human Traffic of Stadium (≥ 3 consecutive high-traffic rows)
