# 0196 — Delete Duplicate Emails

> LeetCode #196 · Difficulty: Easy
> **Categories:** Database, Hash Map, Sorting

---

## Problem Statement

Table: `Person`

```
+-------------+---------+
| Column Name | Type    |
+-------------+---------+
| id          | int     |
| email       | varchar |
+-------------+---------+
```

`id` is the primary key (column with unique values) for this table.
Each row of this table contains an email. The emails will not contain uppercase letters.

Write a solution to **delete** all duplicate emails, keeping only one unique email with the smallest `id`.

For SQL users, please note that you are supposed to write a `DELETE` statement and not a `SELECT` one.

For Pandas users, please note that you are supposed to modify `Person` in place.

After running your script, the answer shown is the `Person` table. The driver will first compile and run your piece of code and then show the `Person` table. The final order of the `Person` table **does not matter**.

The result format is in the following example.

**Example 1**
```
Input:
Person table:
+----+------------------+
| id | email            |
+----+------------------+
| 1  | john@example.com |
| 2  | bob@example.com  |
| 3  | john@example.com |
+----+------------------+
Output:
+----+------------------+
| id | email            |
+----+------------------+
| 1  | john@example.com |
| 2  | bob@example.com  |
+----+------------------+
Explanation: john@example.com is repeated two times. We keep the row with the smallest Id = 1.
```

**Constraints**
- `id` is the primary key — every id is unique.
- Emails are non-NULL and contain no uppercase letters.
- The final order of the surviving rows does not matter.

> ℹ️ This is a **Database** problem. Per this repo's Go-only convention, `main.go`
> models the `Person` table as a slice of structs and implements the delete logic
> with classic DSA techniques; the canonical SQL answers appear in
> [Key Takeaways](#key-takeaways).

---

## Company Frequency

| Company    | Frequency        | Last Reported |
|------------|------------------|---------------|
| Amazon     | ★★★★☆ High       | 2024          |
| Google     | ★★★☆☆ Medium     | 2024          |
| Microsoft  | ★★★☆☆ Medium     | 2023          |
| Bloomberg  | ★★☆☆☆ Low        | 2023          |
| Meta       | ★★☆☆☆ Low        | 2023          |
| Adobe      | ★★☆☆☆ Low        | 2023          |

> ⚠️ Frequency data is crowd-sourced from LeetCode Discuss, Glassdoor, and
> community interview reports. Treat as a signal, not a guarantee.

---

## DSA Concepts Used

- **Hash Map** — one pass builds `email → min(id)`; a second pass keeps exactly the rows matching that minimum. Turns the O(n²) self-join into O(n). → see [`/dsa/hash_map.md`](/dsa/hash_map.md)
- **Sorting** — sorting by `(email, id)` makes duplicates adjacent with the keeper first, reducing dedup to a linear sweep. → see [`/dsa/sorting.md`](/dsa/sorting.md)

---

## Approaches Overview

| # | Approach | Time | Space | When to use |
|---|----------|------|-------|-------------|
| 1 | Brute Force (self-join scan) | O(n²) | O(n) | Tiny tables; mirrors the SQL self-join `DELETE` exactly |
| 2 | Sort + Sweep | O(n log n) | O(n) | No hash map available; want duplicates grouped for inspection |
| 3 | Hash Map (Optimal) ✅ | O(n) | O(n) | General case — single-pass group-by in linear time |

---

## Approach 1 — Brute Force

### Intuition
A row dies exactly when the table contains another row with the **same email and a smaller id**. So test each row against every other row and keep it only if no such witness exists. This is a literal re-enactment of the classic SQL self-join delete:

```sql
DELETE p1 FROM Person p1, Person p2
WHERE p1.email = p2.email AND p1.id > p2.id;
```

### Algorithm
1. For every row `i`, scan every row `j` of the table.
2. If `person[j].Email == person[i].Email` and `person[j].ID < person[i].ID`, row `i` is a duplicate — mark it and stop scanning (one witness suffices; the strict `<` also rules out `j == i`).
3. Append every unmarked row to the result.
4. Sort the survivors by id for deterministic output (any order is accepted).

### Complexity
- **Time:** O(n²) — each of the n rows may compare against all n rows.
- **Space:** O(n) — only the output slice of kept rows.

### Code
```go
func bruteForce(person []Person) []Person {
	kept := []Person{} // rows that survive the DELETE
	for i := range person {
		duplicate := false
		for j := range person {
			// A row with the same email but a strictly smaller id proves row i
			// is a duplicate (the id comparison also rules out j == i).
			if person[j].Email == person[i].Email && person[j].ID < person[i].ID {
				duplicate = true
				break // one witness is enough
			}
		}
		if !duplicate {
			kept = append(kept, person[i])
		}
	}
	return sortByID(kept)
}
```

### Dry Run — Example 1: `Person = [(1, john@), (2, bob@), (3, john@)]`

| Row i checked | Witness search (same email, smaller id) | duplicate? | kept after step |
|---------------|------------------------------------------|------------|-----------------|
| (1, john@)    | (2, bob@) ✗ email · (3, john@) ✗ id 3 > 1 | no  | [(1, john@)] |
| (2, bob@)     | no row shares email `bob@`                | no  | [(1, john@), (2, bob@)] |
| (3, john@)    | (1, john@) ✓ same email, id 1 < 3         | yes | unchanged |

Survivors sorted by id → `[(1, john@example.com), (2, bob@example.com)]` ✓

---

## Approach 2 — Sort + Sweep

### Intuition
Sort a copy of the table by `(email, id)`. Now every group of duplicate emails is contiguous **and its keeper (smallest id) stands first in the group**. A single sweep that keeps only "first row of each group" performs the whole delete.

### Algorithm
1. Copy the table (never reorder the caller's data).
2. Sort the copy by email ascending, breaking ties by id ascending.
3. Sweep left to right: keep row `i` iff `i == 0` or `rows[i-1].Email != rows[i].Email`.
4. Sort the survivors by id for deterministic output.

### Complexity
- **Time:** O(n log n) — the sort dominates; the sweep is O(n).
- **Space:** O(n) — the sorted working copy.

### Code
```go
func sortAndSweep(person []Person) []Person {
	rows := make([]Person, len(person))
	copy(rows, person) // work on a copy; leave the input table untouched

	// Sort by (email, id): duplicates adjacent, smallest id first in a group.
	sort.Slice(rows, func(i, j int) bool {
		if rows[i].Email != rows[j].Email {
			return rows[i].Email < rows[j].Email
		}
		return rows[i].ID < rows[j].ID
	})

	kept := []Person{}
	for i, row := range rows {
		// The first row of every email group is the survivor (smallest id).
		if i == 0 || rows[i-1].Email != row.Email {
			kept = append(kept, row)
		}
	}
	return sortByID(kept)
}
```

### Dry Run — Example 1: `Person = [(1, john@), (2, bob@), (3, john@)]`

**After sorting by (email, id):** `[(2, bob@), (1, john@), (3, john@)]`

| i | row        | prev email | first of its group? | kept after step |
|---|------------|------------|---------------------|-----------------|
| 0 | (2, bob@)  | —          | yes (i == 0)        | [(2, bob@)] |
| 1 | (1, john@) | bob@       | yes (email changed) | [(2, bob@), (1, john@)] |
| 2 | (3, john@) | john@      | no (same email)     | unchanged |

Survivors sorted by id → `[(1, john@example.com), (2, bob@example.com)]` ✓

---

## Approach 3 — Hash Map (Optimal)

### Intuition
The keeper of each email group is completely described by one fact: `email → minimum id`. A hash map computes that in a single pass; a second pass keeps exactly the rows whose id equals their email's minimum. This is the Go version of the `GROUP BY` subquery delete:

```sql
DELETE FROM Person
WHERE id NOT IN (
    SELECT * FROM (SELECT MIN(id) FROM Person GROUP BY email) AS keep
);
```

### Algorithm
1. **Pass 1:** for each row, if its email is unseen or its id is smaller than the recorded one, store `minID[email] = id`.
2. **Pass 2:** keep row iff `minID[row.Email] == row.ID` — exactly one row per email can satisfy this because ids are unique.
3. Sort survivors by id for deterministic output.

### Complexity
- **Time:** O(n) — two linear passes with O(1) map operations each.
- **Space:** O(n) — one map entry per distinct email.

### Code
```go
func hashMap(person []Person) []Person {
	minID := make(map[string]int, len(person)) // email → smallest id seen so far
	for _, row := range person {
		// First sighting of an email, or a smaller id than recorded, wins.
		if best, ok := minID[row.Email]; !ok || row.ID < best {
			minID[row.Email] = row.ID
		}
	}

	kept := []Person{}
	for _, row := range person {
		// Exactly one row per email satisfies this: the one with the min id.
		if minID[row.Email] == row.ID {
			kept = append(kept, row)
		}
	}
	return sortByID(kept)
}
```

### Dry Run — Example 1: `Person = [(1, john@), (2, bob@), (3, john@)]`

**Pass 1 — build `minID`:**

| row        | email seen before? | minID after |
|------------|--------------------|-------------|
| (1, john@) | no                 | {john@: 1} |
| (2, bob@)  | no                 | {john@: 1, bob@: 2} |
| (3, john@) | yes, min is 1; 3 > 1 → no change | {john@: 1, bob@: 2} |

**Pass 2 — keep rows whose id is their email's minimum:**

| row        | minID[email] | id == min? | kept after step |
|------------|--------------|------------|-----------------|
| (1, john@) | 1            | ✅         | [(1, john@)] |
| (2, bob@)  | 2            | ✅         | [(1, john@), (2, bob@)] |
| (3, john@) | 1            | ❌ (3 ≠ 1) | unchanged |

Result → `[(1, john@example.com), (2, bob@example.com)]` ✓

---

## Key Takeaways

- **"Keep one representative per group" = group-by-min pattern.** Whether in SQL (`GROUP BY email, MIN(id)`) or Go (`map[email]minID`), dedup-with-a-tiebreak is always: compute each group's champion, then filter everything else out.
- **A duplicate is defined by a witness.** Row *x* dies iff a row with the same key and a smaller id exists — phrasing deletes as an existence test is what makes the brute force (and the SQL self-join) correct even with 3+ copies of the same email.
- **The canonical SQL answers**, both accepted on LeetCode:
  ```sql
  -- Self-join delete (most common interview answer)
  DELETE p1 FROM Person p1, Person p2
  WHERE p1.email = p2.email AND p1.id > p2.id;

  -- GROUP BY subquery delete (needs the derived-table wrapper in MySQL,
  -- because MySQL cannot DELETE from a table it is directly SELECTing from)
  DELETE FROM Person
  WHERE id NOT IN (
      SELECT * FROM (SELECT MIN(id) FROM Person GROUP BY email) AS keep
  );
  ```
- **Sort makes duplicates adjacent.** Sorting by `(key, tiebreak)` turns any grouping problem into a linear sweep — the same trick used in Remove Duplicates from Sorted Array (#26) and Group Anagrams (#49).
- **Unique-id tiebreaks are exact filters.** Because ids are unique, `id == minID[email]` selects *exactly one* row per email — no extra "already kept" bookkeeping needed.

---

## Related Problems

- LeetCode #182 — Duplicate Emails (the SELECT half of this problem: find the duplicated emails)
- LeetCode #175 — Combine Two Tables (basic SQL table manipulation)
- LeetCode #181 — Employees Earning More Than Their Managers (self-join comparison pattern)
- LeetCode #183 — Customers Who Never Order (anti-join / NOT IN pattern)
- LeetCode #26 — Remove Duplicates from Sorted Array (keep-first-occurrence sweep over sorted data)
- LeetCode #217 — Contains Duplicate (hash-set duplicate detection, the array analogue)
