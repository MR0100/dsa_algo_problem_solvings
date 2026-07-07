# 0182 — Duplicate Emails

> LeetCode #182 · Difficulty: Easy
> **Categories:** Database, Hash Table, Counting, Sorting

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
id is the primary key (column with unique values) for this table.
Each row of this table contains an email. The emails will not contain
uppercase letters.
```

Write a solution to report all the **duplicate emails**. Note that it's
guaranteed that the email field is not `NULL`.

Return the result table in **any order**.

**Example 1:**

```
Input:
Person table:
+----+---------+
| id | email   |
+----+---------+
| 1  | a@b.com |
| 2  | c@d.com |
| 3  | a@b.com |
+----+---------+
Output:
+---------+
| Email   |
+---------+
| a@b.com |
+---------+
Explanation: a@b.com is repeated two times.
```

**Constraints:**

- `id` is the primary key — every id is unique.
- Emails contain no uppercase letters and are never `NULL`.
- Each duplicate email must appear **exactly once** in the result, no matter how many times it repeats.

> **Repo note:** this is a LeetCode *Database* problem. In this Go repo the
> table is modelled as `[]Person{ID, Email}`, and each approach re-implements
> the SQL `GROUP BY email HAVING COUNT(email) > 1` query in memory.

---

## Company Frequency

| Company    | Frequency        | Last Reported |
|------------|------------------|---------------|
| Amazon     | ★★★★☆ High       | 2024          |
| Microsoft  | ★★★☆☆ Medium     | 2024          |
| Google     | ★★★☆☆ Medium     | 2023          |
| Bloomberg  | ★★☆☆☆ Low        | 2023          |
| Oracle     | ★★☆☆☆ Low        | 2023          |
| Salesforce | ★★☆☆☆ Low        | 2023          |

> ⚠️ Frequency data is crowd-sourced from LeetCode Discuss, Glassdoor, and
> community interview reports. Treat as a signal, not a guarantee.

---

## DSA Concepts Used

- **Hash Map** — email → frequency counting is the direct translation of `GROUP BY ... HAVING COUNT(*) > 1`; emitting at count == 2 dedupes for free → see [`/dsa/hash_map.md`](/dsa/hash_map.md)
- **Sorting** — sorting groups equal emails into adjacent runs, so duplicates are found by comparing neighbours → see [`/dsa/sorting.md`](/dsa/sorting.md)

---

## Approaches Overview

| # | Approach | Time | Space | When to use |
|---|----------|------|-------|-------------|
| 1 | Brute Force | O(n²) | O(1) | Baseline; only for tiny tables or when no extra memory is allowed |
| 2 | Sorting | O(n log n) | O(n) | When the data is already sorted or memory for a map is unwelcome |
| 3 | Hash Map Counting (Optimal) | O(n) | O(n) | Always — one pass, and the standard duplicate-detection idiom |

---

## Approach 1 — Brute Force

### Intuition

An email is a duplicate iff at least two rows carry it. The subtle part is
reporting it **once**: with no memory allowed, let only the *first occurrence*
of each email do the reporting. Row `i` reports its email iff no earlier row
holds the same email (so `i` is the first occurrence) **and** some later row
repeats it (so it really is a duplicate).

### Algorithm

1. For each row `i`, scan rows `0..i-1`; if any holds the same email, skip
   row `i` — an earlier occurrence owns the report.
2. Otherwise scan rows `i+1..n-1`; on the first row with an equal email,
   append the email to the result and stop scanning.
3. Return the collected emails.

### Complexity

- **Time:** O(n²) — each of the n rows may trigger two linear scans of the table.
- **Space:** O(1) — only loop indices; no auxiliary containers beyond the output.

### Code

```go
func bruteForce(people []Person) []string {
	result := []string{}
	for i, p := range people {
		isFirst := true
		for j := 0; j < i; j++ { // does an earlier row already hold this email?
			if people[j].Email == p.Email {
				isFirst = false // yes → that row owns the report, not us
				break
			}
		}
		if !isFirst {
			continue // avoid reporting the same email twice
		}
		for k := i + 1; k < len(people); k++ { // is the email repeated later?
			if people[k].Email == p.Email {
				result = append(result, p.Email) // first occurrence + a repeat = duplicate
				break                            // one report per email is enough
			}
		}
	}
	return result
}
```

### Dry Run

Example 1: `people = [(1, a@b.com), (2, c@d.com), (3, a@b.com)]`

| Step | `i` | Email | Earlier copy (`j < i`)? | Later copy (`k > i`)? | `result` |
|------|-----|---------|--------------------------|------------------------|----------|
| 1 | 0 | a@b.com | no → first occurrence | yes, at k = 2 | `[a@b.com]` |
| 2 | 1 | c@d.com | no → first occurrence | no | `[a@b.com]` |
| 3 | 2 | a@b.com | yes, at j = 0 → skip | — | `[a@b.com]` |

Return `[a@b.com]` ✓ (matches the expected output).

---

## Approach 2 — Sorting

### Intuition

Sorting pulls all copies of the same email next to each other, so a duplicate
is simply a run of length ≥ 2. To emit each duplicate exactly once — even when
an email appears 3+ times — report only at the **first adjacency** of each run
(the position where the run starts).

### Algorithm

1. Project the email column into a fresh slice (keep the input rows untouched).
2. Sort the slice lexicographically.
3. Scan indices `1..n-1`: report `emails[i]` when `emails[i] == emails[i-1]`
   **and** the pair starts a run (`i == 1` or `emails[i-2] != emails[i-1]`).

### Complexity

- **Time:** O(n log n) — the sort dominates; the run scan is a single O(n) pass.
- **Space:** O(n) — the projected copy of the email column.

### Code

```go
func sorting(people []Person) []string {
	emails := make([]string, len(people))
	for i, p := range people {
		emails[i] = p.Email // project the email column, keep input immutable
	}
	sort.Strings(emails) // duplicates become adjacent

	result := []string{}
	for i := 1; i < len(emails); i++ {
		// Report only at the first adjacency of a run so each duplicate email
		// is emitted exactly once, even when it appears 3+ times.
		if emails[i] == emails[i-1] && (i == 1 || emails[i-1] != emails[i-2]) {
			result = append(result, emails[i])
		}
	}
	return result
}
```

### Dry Run

Example 1: projected column `[a@b.com, c@d.com, a@b.com]` → after sorting
`emails = [a@b.com, a@b.com, c@d.com]`.

| Step | `i` | `emails[i-1]` vs `emails[i]` | Equal? | Run start (`i==1` or `emails[i-2]` differs)? | `result` |
|------|-----|------------------------------|--------|-----------------------------------------------|----------|
| 1 | 1 | a@b.com vs a@b.com | yes | yes (`i == 1`) | `[a@b.com]` |
| 2 | 2 | a@b.com vs c@d.com | no | — | `[a@b.com]` |

Return `[a@b.com]` ✓.

---

## Approach 3 — Hash Map Counting (Optimal)

### Intuition

`GROUP BY email HAVING COUNT(email) > 1` is frequency counting in disguise.
One pass over the rows builds the counts, and there is a neat trick to skip the
second pass entirely: append an email at the *exact moment* its count reaches
2. Earlier (count 1) it is not yet a duplicate; later (count 3, 4, …) it has
already been reported.

### Algorithm

1. Create `counts`, a map from email → occurrences seen so far.
2. For each row, increment `counts[email]`.
3. If the count just became exactly 2, append the email to the result.
4. Return the result after the single pass.

### Complexity

- **Time:** O(n) — one pass, O(1) average per map update.
- **Space:** O(n) — at most one map entry per distinct email.

### Code

```go
func hashMap(people []Person) []string {
	counts := make(map[string]int, len(people)) // email → occurrences so far
	result := []string{}
	for _, p := range people {
		counts[p.Email]++
		if counts[p.Email] == 2 { // the exact moment it becomes a duplicate
			result = append(result, p.Email) // report once, never again
		}
	}
	return result
}
```

### Dry Run

Example 1: `people = [(1, a@b.com), (2, c@d.com), (3, a@b.com)]`

| Step | Row | `counts` after increment | Count just hit 2? | `result` |
|------|-----|---------------------------|--------------------|----------|
| 1 | (1, a@b.com) | `{a@b.com: 1}` | no | `[]` |
| 2 | (2, c@d.com) | `{a@b.com: 1, c@d.com: 1}` | no | `[]` |
| 3 | (3, a@b.com) | `{a@b.com: 2, c@d.com: 1}` | yes | `[a@b.com]` |

Return `[a@b.com]` ✓.

---

## Key Takeaways

- **`GROUP BY` + `HAVING COUNT > 1` ≡ hash-map frequency counting** — the most common SQL-to-code translation there is.
- **Emit when the count hits exactly 2** — a one-pass dedup trick: `== 2` fires once per duplicate, whereas `>= 2` would fire on every repeat. The same trick powers streaming "first time X became a duplicate" problems.
- **Sorting turns global duplicate detection into a local neighbour check**; reporting only at run starts dedupes without any set.
- Canonical SQL solutions this maps to:

```sql
-- Aggregation
SELECT email AS Email
FROM Person
GROUP BY email
HAVING COUNT(email) > 1;

-- Self-join alternative
SELECT DISTINCT p1.email AS Email
FROM Person p1 JOIN Person p2
  ON p1.email = p2.email AND p1.id <> p2.id;
```

---

## Related Problems

- LeetCode #181 — Employees Earning More Than Their Managers (same table-modelling pattern)
- LeetCode #183 — Customers Who Never Order (hash set membership instead of counting)
- LeetCode #196 — Delete Duplicate Emails (same table, delete instead of report)
- LeetCode #217 — Contains Duplicate (identical duplicate-detection pattern on arrays)
- LeetCode #219 — Contains Duplicate II (duplicate detection with an index-distance twist)
