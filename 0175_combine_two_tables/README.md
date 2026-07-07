# 0175 — Combine Two Tables

> LeetCode #175 · Difficulty: Easy
> **Categories:** Database, SQL, LEFT JOIN

---

## Problem Statement

**Table: `Person`**

```
+-------------+---------+
| Column Name | Type    |
+-------------+---------+
| personId    | int     |
| lastName    | varchar |
| firstName   | varchar |
+-------------+---------+
personId is the primary key (column with unique values) for this table.
This table contains information about the ID of some persons and their first and last names.
```

**Table: `Address`**

```
+-------------+---------+
| Column Name | Type    |
+-------------+---------+
| addressId   | int     |
| personId    | int     |
| city        | varchar |
| state       | varchar |
+-------------+---------+
addressId is the primary key (column with unique values) for this table.
Each row of this table contains information about the city and state of one person with ID = PersonId.
```

Write a solution to report the first name, last name, city, and state of each person in the `Person` table. If the address of a `personId` is not present in the `Address` table, report `null` instead.

Return the result table in **any order**.

The result format is in the following example.

**Example 1:**

```
Input:
Person table:
+----------+----------+-----------+
| personId | lastName | firstName |
+----------+----------+-----------+
| 1        | Wang     | Allen     |
| 2        | Alice    | Bob       |
+----------+----------+-----------+
Address table:
+-----------+----------+---------------+------------+
| addressId | personId | city          | state      |
+-----------+----------+---------------+------------+
| 1         | 2        | New York City | New York   |
| 2         | 3        | Leetcode      | California |
+-----------+----------+---------------+------------+
Output:
+-----------+----------+---------------+----------+
| firstName | lastName | city          | state    |
+-----------+----------+---------------+----------+
| Allen     | Wang     | Null          | Null     |
| Bob       | Alice    | New York City | New York |
+-----------+----------+---------------+----------+
Explanation:
There is no address in the address table for the personId = 1 so we return null in their city and state.
addressId = 1 contains information about the address of personId = 2.
```

---

## Company Frequency

| Company    | Frequency        | Last Reported |
|------------|------------------|---------------|
| Amazon     | ★★★★☆ High       | 2024          |
| Meta       | ★★★☆☆ Medium     | 2023          |
| Google     | ★★★☆☆ Medium     | 2023          |
| Microsoft  | ★★★☆☆ Medium     | 2023          |
| Bloomberg  | ★★☆☆☆ Low        | 2023          |
| Adobe      | ★★☆☆☆ Low        | 2022          |

> ⚠️ Frequency data is crowd-sourced from LeetCode Discuss, Glassdoor, and
> community interview reports. Treat as a signal, not a guarantee.

---

## DSA Concepts Used

This problem is officially a **SQL** exercise — the accepted answer is a single `LEFT JOIN`:

```sql
SELECT p.firstName, p.lastName, a.city, a.state
FROM Person p LEFT JOIN Address a ON p.personId = a.personId;
```

The Go implementation reproduces what a database engine actually *does* to execute that join, using the three classic join algorithms. The relevant DSA concepts are therefore the data structures/patterns each algorithm relies on:

- **Hash Map** — the hash join builds `map[personId] → []Address` to index the right table for O(1)-average probing, exactly what a database index on `Address.personId` provides → see [`/dsa/hash_map.md`](/dsa/hash_map.md)
- **Sorting** — the sort-merge join sorts both tables on the join key so that matching rows line up during a single sweep → see [`/dsa/sorting.md`](/dsa/sorting.md)
- **Two Pointers** — after sorting, the merge phase advances two indices (one per table) in lockstep, never rewinding → see [`/dsa/two_pointers.md`](/dsa/two_pointers.md)

---

## Approaches Overview

| # | Approach | Time | Space | When to use |
|---|----------|------|-------|-------------|
| 1 | Brute Force (Nested Loop Join) | O(P·A) | O(1) extra | No index, no memory for a hash table; the fallback a DB uses on tiny inputs |
| 2 | Sort-Merge Join | O(P log P + A log A) | O(P + A) | Inputs already sorted, or must be sorted anyway (e.g. `ORDER BY` on the join key) |
| 3 | Hash Join (Optimal) | O(P + A) | O(A) | Unsorted, unindexed tables — the query planner's default equi-join |

*(P = number of persons, A = number of addresses.)*

---

## Approach 1 — Brute Force (Nested Loop Join)

### Intuition

A `LEFT JOIN` keeps **every** row of the left table (`Person`). For each person, scan the whole `Address` table for rows with the same `personId`: emit one output row per match; if the scan finds nothing, still emit the person once, padding `city`/`state` with `null`. This is literally how a database executes a join when it has no index and no memory for a hash table — the driving (outer) table is the left one so that unmatched left rows are never dropped.

### Algorithm

1. For each person `p` (preserving `Person`-table order):
2. Scan all addresses; for each with `a.personId == p.personId`, emit `(p.firstName, p.lastName, a.city, a.state)` and set a `matched` flag.
3. If no address matched, emit `(p.firstName, p.lastName, null, null)` — the LEFT JOIN survivor row.

### Complexity

- **Time:** O(P·A) — every person scans the entire address list; the two nested loops multiply.
- **Space:** O(1) extra — output aside, only a boolean flag per person scan.

### Code

```go
func nestedLoopJoin(persons []Person, addresses []Address) []Row {
	rows := []Row{}
	for _, p := range persons { // outer loop: LEFT table drives a LEFT JOIN
		matched := false
		for _, a := range addresses { // inner loop: full scan per person
			if a.PersonID == p.PersonID { // join predicate ON p.personId = a.personId
				rows = append(rows, Row{p.FirstName, p.LastName, a.City, a.State})
				matched = true // remember: this person must not get a null row
			}
		}
		if !matched { // LEFT JOIN semantics: unmatched left rows survive with NULLs
			rows = append(rows, Row{p.FirstName, p.LastName, "null", "null"})
		}
	}
	return rows
}
```

### Dry Run

Example 1: `persons = [{1, Wang, Allen}, {2, Alice, Bob}]`, `addresses = [{1, personId 2, New York City, New York}, {2, personId 3, Leetcode, California}]`.

| Step | Outer person `p` | Inner address `a` | `a.personId == p.personId`? | Action | `matched` | rows emitted so far |
|------|------------------|-------------------|-----------------------------|--------|-----------|---------------------|
| 1 | {1, Allen Wang} | {personId 2} | 2 == 1? no | skip | false | — |
| 2 | {1, Allen Wang} | {personId 3} | 3 == 1? no | skip | false | — |
| 3 | {1, Allen Wang} | (scan done) | — | not matched → emit null row | — | `Allen, Wang, null, null` |
| 4 | {2, Bob Alice} | {personId 2} | 2 == 2? **yes** | emit match row | true | + `Bob, Alice, New York City, New York` |
| 5 | {2, Bob Alice} | {personId 3} | 3 == 2? no | skip | true | — |
| 6 | {2, Bob Alice} | (scan done) | — | already matched → no null row | — | — |

Result:
```
| Allen     | Wang     | null          | null     |
| Bob       | Alice    | New York City | New York |
```
✔

---

## Approach 2 — Sort-Merge Join

### Intuition

Once both tables are ordered by the join key, matching rows must line up — so instead of re-scanning, we *sweep* them together with two pointers. Advance whichever pointer lags: when the address key is smaller, that address matches nobody, skip it; when keys are equal, emit the pairing; when the person key is smaller, that person has no address, emit a `null` row. This is the join a database picks when the inputs are already sorted or must be sorted anyway (e.g. for an `ORDER BY` on the join key).

### Algorithm

1. Copy and sort `persons` by `personId`, and `addresses` by `personId` (copies, so the caller's tables stay untouched).
2. Sweep with pointer `j` into the sorted addresses; for each person `p` in order:
   - Skip addresses whose `personId < p.personId` (`j++`) — orphan addresses that join to nobody.
   - Emit one row for every address whose `personId == p.personId` (handles a person with several addresses), setting `matched`.
   - If none matched, emit a `null` row.
3. Because `personId` is unique in `Person`, each subsequent person's id is strictly larger, so leaving `j` at the block start re-skips correctly.

### Complexity

- **Time:** O(P log P + A log A) — sorting both tables dominates; the merge sweep itself is O(P + A).
- **Space:** O(P + A) — sorted copies of both tables (we must not mutate the caller's input).

### Code

```go
func sortMergeJoin(persons []Person, addresses []Address) []Row {
	// Sort copies so the input tables stay untouched.
	ps := append([]Person(nil), persons...)
	as := append([]Address(nil), addresses...)
	sort.Slice(ps, func(x, y int) bool { return ps[x].PersonID < ps[y].PersonID })
	sort.Slice(as, func(x, y int) bool { return as[x].PersonID < as[y].PersonID })

	rows := []Row{}
	j := 0 // sweep pointer into the sorted addresses
	for i := 0; i < len(ps); i++ {
		p := ps[i]
		// Skip addresses whose personId is smaller — they join to nobody here.
		for j < len(as) && as[j].PersonID < p.PersonID {
			j++
		}
		// Emit every address sharing this person's id (handles duplicates).
		matched := false
		for k := j; k < len(as) && as[k].PersonID == p.PersonID; k++ {
			rows = append(rows, Row{p.FirstName, p.LastName, as[k].City, as[k].State})
			matched = true
		}
		if !matched { // no address block for this person → LEFT JOIN null row
			rows = append(rows, Row{p.FirstName, p.LastName, "null", "null"})
		}
		// NOTE: j is left at the block start; personId is unique in Person,
		// so the next person's id is strictly larger and re-skips correctly.
	}
	return rows // ordered by personId (a valid order: "return in any order")
}
```

### Dry Run

Example 1. After sorting: `ps = [{1, Allen Wang}, {2, Bob Alice}]` (already sorted), `as = [{personId 2, New York City}, {personId 3, Leetcode}]` (already sorted).

| Step | `i` / person `p` | `j` (skip while `as[j].personId < p.personId`) | equal-key block scan | `matched` | rows emitted |
|------|------------------|-----------------------------------------------|----------------------|-----------|--------------|
| 1 | i=0, p={1, Allen} | `as[0].personId = 2`, not `< 1` → no skip, `j=0` | `as[0].personId = 2 == 1`? no → block empty | false | emit `Allen, Wang, null, null` |
| 2 | i=1, p={2, Bob} | `as[0].personId = 2`, not `< 2` → no skip, `j=0` | `as[0].personId = 2 == 2`? **yes** → emit; `as[1].personId = 3 == 2`? no → stop | true | emit `Bob, Alice, New York City, New York` |

Result (ordered by personId 1 then 2):
```
| Allen     | Wang     | null          | null     |
| Bob       | Alice    | New York City | New York |
```
✔

---

## Approach 3 — Hash Join (Optimal)

### Intuition

The nested loop wastes time re-scanning addresses. Build the "index" ourselves: one pass over `Address` fills `map[personId] → []Address`; then each person finds its matches in O(1) average. Build + probe = one pass over each table. This is the hash join a query planner chooses for unsorted, unindexed tables — and the map is exactly what a database index on `Address.personId` would give us. Storing a **slice** per key keeps LEFT JOIN correct when a person has several addresses.

### Algorithm

1. **Build phase:** create `byPerson = map[int][]Address`; append every address under its `personId`.
2. **Probe phase:** for each person (preserving `Person`-table order), look up `byPerson[personId]`:
   - Empty bucket → emit `(firstName, lastName, null, null)`.
   - Otherwise → emit one row per address in the bucket.

### Complexity

- **Time:** O(P + A) — one pass to build the map, one pass to probe it; map operations are O(1) average.
- **Space:** O(A) — the hash map holds every address once.

### Code

```go
func hashJoin(persons []Person, addresses []Address) []Row {
	// Build phase: index the smaller/right table by the join key.
	byPerson := make(map[int][]Address, len(addresses))
	for _, a := range addresses {
		byPerson[a.PersonID] = append(byPerson[a.PersonID], a)
	}
	// Probe phase: each person fetches its address bucket in O(1) average.
	rows := []Row{}
	for _, p := range persons {
		bucket := byPerson[p.PersonID] // nil slice when the person has no address
		if len(bucket) == 0 {
			rows = append(rows, Row{p.FirstName, p.LastName, "null", "null"})
			continue
		}
		for _, a := range bucket { // one output row per matching address
			rows = append(rows, Row{p.FirstName, p.LastName, a.City, a.State})
		}
	}
	return rows
}
```

### Dry Run

Example 1.

**Build phase** — scan `addresses`:

| Step | address `a` | map after insert |
|------|-------------|------------------|
| 1 | {personId 2, New York City, New York} | `{2: [New York City/New York]}` |
| 2 | {personId 3, Leetcode, California} | `{2: [New York City/New York], 3: [Leetcode/California]}` |

**Probe phase** — scan `persons` in order:

| Step | person `p` | `byPerson[p.personId]` | bucket empty? | rows emitted |
|------|------------|------------------------|---------------|--------------|
| 1 | {1, Allen Wang} | `byPerson[1]` → nil | yes | `Allen, Wang, null, null` |
| 2 | {2, Bob Alice} | `byPerson[2]` → [New York City/New York] | no | `Bob, Alice, New York City, New York` |

Result:
```
| Allen     | Wang     | null          | null     |
| Bob       | Alice    | New York City | New York |
```
✔

---

## Key Takeaways

- **A `LEFT JOIN` keeps every left-table row.** The unmatched left rows survive with `NULL`s on the right-table columns — the driving/outer loop (or probe loop) must run over the left table, and each left row that finds no match still emits exactly one output row.
- **There are three canonical ways an engine executes an equi-join**, and knowing them explains query-plan choices: nested-loop (simple, O(P·A), good for tiny inputs), sort-merge (O(n log n), great when data is already sorted or ordering is needed anyway), and hash join (O(P + A), the default for unsorted unindexed tables).
- **A hash map *is* an index.** Building `map[key] → []rows` and probing it is exactly what a database index on the join column buys you: O(1)-average lookups turning an O(P·A) scan into O(P + A).
- **Store a slice, not a single value, per key** when a key can map to many rows (one person, many addresses) — otherwise multi-match joins silently drop rows.
- **"Return in any order"** is a license to pick whatever order is cheapest: nested-loop and hash join keep Person-table order for free; sort-merge naturally yields join-key order.

---

## Related Problems

- LeetCode #181 — Employees Earning More Than Their Managers (self-join on the same table)
- LeetCode #183 — Customers Who Never Order (LEFT JOIN + `IS NULL` anti-join)
- LeetCode #197 — Rising Temperature (self-join on adjacent rows)
- LeetCode #577 — Employee Bonus (LEFT JOIN keeping unmatched left rows)
- LeetCode #584 — Find Customer Referee (NULL-aware filtering after a join)
