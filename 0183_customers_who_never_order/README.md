# 0183 — Customers Who Never Order

> LeetCode #183 · Difficulty: Easy
> **Categories:** Database, Hash Table, Sorting, Binary Search

---

## Problem Statement

Table: `Customers`

```
+-------------+---------+
| Column Name | Type    |
+-------------+---------+
| id          | int     |
| name        | varchar |
+-------------+---------+
id is the primary key (column with unique values) for this table.
Each row of this table indicates the ID and name of a customer.
```

Table: `Orders`

```
+-------------+------+
| Column Name | Type |
+-------------+------+
| id          | int  |
| customerId  | int  |
+-------------+------+
id is the primary key (column with unique values) for this table.
customerId is a foreign key (reference columns) of the ID from the
Customers table.
Each row of this table indicates the ID of an order and the ID of the
customer who ordered it.
```

Write a solution to find all customers who **never order anything**.

Return the result table in **any order**.

**Example 1:**

```
Input:
Customers table:
+----+-------+
| id | name  |
+----+-------+
| 1  | Joe   |
| 2  | Henry |
| 3  | Sam   |
| 4  | Max   |
+----+-------+
Orders table:
+----+------------+
| id | customerId |
+----+------------+
| 1  | 3          |
| 2  | 1          |
+----+------------+
Output:
+-----------+
| Customers |
+-----------+
| Henry     |
| Max       |
+-----------+
Explanation: Henry and Max never placed any order; Joe and Sam each appear
in the Orders table, so they are excluded.
```

**Constraints:**

- `Customers.id` and `Orders.id` are primary keys — unique within their tables.
- `Orders.customerId` always references a valid `Customers.id` (foreign key).
- A customer may appear in `Orders` any number of times; one appearance is enough to exclude them.

> **Repo note:** this is a LeetCode *Database* problem. In this Go repo the
> tables are modelled as `[]Customer{ID, Name}` and `[]Order{ID, CustomerID}`,
> and each approach re-implements the SQL anti-join
> `WHERE id NOT IN (SELECT customerId FROM Orders)` in memory.

---

## Company Frequency

| Company    | Frequency        | Last Reported |
|------------|------------------|---------------|
| Amazon     | ★★★★☆ High       | 2024          |
| Bloomberg  | ★★★☆☆ Medium     | 2024          |
| Google     | ★★★☆☆ Medium     | 2023          |
| Microsoft  | ★★★☆☆ Medium     | 2023          |
| Apple      | ★★☆☆☆ Low        | 2023          |

> ⚠️ Frequency data is crowd-sourced from LeetCode Discuss, Glassdoor, and
> community interview reports. Treat as a signal, not a guarantee.

---

## DSA Concepts Used

- **Hash Map (as a Hash Set)** — materialising the subquery's ids into a set gives O(1) "has this customer ever ordered?" membership tests — a hash anti-join → see [`/dsa/hash_map.md`](/dsa/hash_map.md)
- **Sorting** — sorting the foreign-key column once simulates an index on `Orders.customerId` → see [`/dsa/sorting.md`](/dsa/sorting.md)
- **Binary Search** — membership testing inside the sorted id column in O(log o) → see [`/dsa/binary_search.md`](/dsa/binary_search.md)

---

## Approaches Overview

| # | Approach | Time | Space | When to use |
|---|----------|------|-------|-------------|
| 1 | Brute Force (Nested-Loop Anti-Join) | O(c·o) | O(1) | Baseline; acceptable only for tiny tables |
| 2 | Sort + Binary Search | O((c+o) log o) | O(o) | When the orders column is already sorted / hash memory is unwelcome |
| 3 | Hash Set (Hash Anti-Join, Optimal) | O(c+o) | O(o) | Always — the canonical `NOT IN` execution strategy |

---

## Approach 1 — Brute Force (Nested-Loop Anti-Join)

### Intuition

"Never ordered" is a **negative** condition: a customer belongs in the answer
only if *no* row of `Orders` references them. The naive evaluation of
`NOT IN` checks each customer against every single order — a nested-loop
anti-join. One match anywhere disqualifies the customer; only a completely
match-free scan keeps them.

### Algorithm

1. For every customer `c`, scan the whole `Orders` slice.
2. If any order has `o.CustomerID == c.ID`, mark `c` as having ordered and
   stop the scan early.
3. After the scan, if no order matched, append `c.Name` to the result.
4. Return the collected names.

### Complexity

- **Time:** O(c·o) — each of the c customers may scan all o orders (customers in the answer always scan all of them).
- **Space:** O(1) — a boolean flag and loop variables only.

### Code

```go
func bruteForce(customers []Customer, orders []Order) []string {
	result := []string{}
	for _, c := range customers {
		ordered := false
		for _, o := range orders { // inner scan: does any order reference c?
			if o.CustomerID == c.ID {
				ordered = true // found one — c is not in the answer
				break          // no need to scan the remaining orders
			}
		}
		if !ordered {
			result = append(result, c.Name) // survived the anti-join
		}
	}
	return result
}
```

### Dry Run

Example 1: `customers = [Joe(1), Henry(2), Sam(3), Max(4)]`,
`orders = [(1, cust 3), (2, cust 1)]`

| Step | Customer | Inner scan of orders | `ordered` | `result` |
|------|----------|----------------------|-----------|----------|
| 1 | Joe (id 1) | (1,3): 3≠1 → no; (2,1): 1==1 → match, break | true | `[]` |
| 2 | Henry (id 2) | (1,3): 3≠2; (2,1): 1≠2 → no match | false | `[Henry]` |
| 3 | Sam (id 3) | (1,3): 3==3 → match, break | true | `[Henry]` |
| 4 | Max (id 4) | (1,3): 3≠4; (2,1): 1≠4 → no match | false | `[Henry, Max]` |

Return `[Henry, Max]` ✓ (matches the expected output).

---

## Approach 2 — Sort + Binary Search

### Intuition

The inner scan only ever asks one question: *does customer id X appear
anywhere in Orders?* That is a membership test on the `customerId` column.
Sorting that column once (like building an index on the foreign key) lets
every membership test run as an O(log o) binary search instead of an O(o)
scan.

### Algorithm

1. Project `Orders` down to its `CustomerID` column into a fresh int slice.
2. Sort the slice ascending.
3. For each customer, run `sort.SearchInts(ids, c.ID)` — it returns the first
   index whose value is `>= c.ID`.
4. If that index is past the end, or holds a different value, `c.ID` is absent
   → the customer never ordered → keep `c.Name`.

### Complexity

- **Time:** O((c+o) log o) — O(o log o) to sort, then c binary searches of O(log o) each.
- **Space:** O(o) — the projected, sorted copy of the customerId column.

### Code

```go
func sortAndBinarySearch(customers []Customer, orders []Order) []string {
	ids := make([]int, len(orders))
	for i, o := range orders {
		ids[i] = o.CustomerID // project the foreign-key column
	}
	sort.Ints(ids) // sorted column ≈ index on Orders.customerId

	result := []string{}
	for _, c := range customers {
		// SearchInts returns the first index with ids[idx] >= c.ID.
		idx := sort.SearchInts(ids, c.ID)
		if idx == len(ids) || ids[idx] != c.ID { // id absent → never ordered
			result = append(result, c.Name)
		}
	}
	return result
}
```

### Dry Run

Example 1: projected column `[3, 1]` → after sorting `ids = [1, 3]`.

| Step | Customer | `SearchInts(ids, c.ID)` → `idx` | `idx == len(ids)`? | `ids[idx] != c.ID`? | Keep? | `result` |
|------|----------|--------------------------------|--------------------|--------------------|-------|----------|
| 1 | Joe (id 1) | 0 | no | `ids[0] = 1` == 1 → no | no | `[]` |
| 2 | Henry (id 2) | 1 | no | `ids[1] = 3` ≠ 2 → yes | yes | `[Henry]` |
| 3 | Sam (id 3) | 1 | no | `ids[1] = 3` == 3 → no | no | `[Henry]` |
| 4 | Max (id 4) | 2 | yes (len = 2) | — | yes | `[Henry, Max]` |

Return `[Henry, Max]` ✓.

---

## Approach 3 — Hash Set (Hash Anti-Join, Optimal)

### Intuition

This is how a database actually executes `NOT IN (subquery)`: materialise the
subquery's values into a hash set, then stream the outer table through it and
keep the rows whose key **misses** — a *hash anti-join*. Two linear passes,
each membership test O(1) on average.

### Algorithm

1. **Build pass:** insert every `o.CustomerID` into the `hasOrdered` set
   (duplicates collapse for free).
2. **Probe pass:** for each customer `c`, keep `c.Name` iff `c.ID` is *not*
   in the set.
3. Return the collected names.

### Complexity

- **Time:** O(c+o) — one pass over each table with O(1) average set operations.
- **Space:** O(o) — the set holds at most one entry per distinct ordering customer.

### Code

```go
func hashSet(customers []Customer, orders []Order) []string {
	hasOrdered := make(map[int]bool, len(orders)) // ids that appear in Orders
	for _, o := range orders {
		hasOrdered[o.CustomerID] = true // build the subquery's value set
	}

	result := []string{}
	for _, c := range customers { // probe side of the anti-join
		if !hasOrdered[c.ID] { // miss → this customer never ordered
			result = append(result, c.Name)
		}
	}
	return result
}
```

### Dry Run

Example 1 — build pass over `orders = [(1,3), (2,1)]` produces
`hasOrdered = {3: true, 1: true}`.

| Step | Probe customer | `hasOrdered[c.ID]` | Keep (`!hasOrdered`)? | `result` |
|------|----------------|---------------------|------------------------|----------|
| 1 | Joe (id 1) | true | no | `[]` |
| 2 | Henry (id 2) | false (miss) | yes | `[Henry]` |
| 3 | Sam (id 3) | true | no | `[Henry]` |
| 4 | Max (id 4) | false (miss) | yes | `[Henry, Max]` |

Return `[Henry, Max]` ✓.

---

## Key Takeaways

- **`NOT IN` / `NOT EXISTS` ≡ hash anti-join:** build a set from the inner
  query, keep outer rows whose key misses. The mirror image of the hash join
  from LeetCode #181 (there a probe *hit* kept the row; here a *miss* does).
- **Set membership is the primitive** — pick linear scan O(o), sorted +
  binary search O(log o), or hash set O(1) per test depending on constraints;
  the surrounding logic never changes.
- **Duplicates in the probe data cost nothing** with a set: inserting the same
  customerId twice collapses into one entry.
- Canonical SQL solutions this maps to:

```sql
-- Subquery anti-join
SELECT name AS Customers
FROM Customers
WHERE id NOT IN (SELECT customerId FROM Orders);

-- LEFT JOIN + NULL filter
SELECT c.name AS Customers
FROM Customers c LEFT JOIN Orders o ON c.id = o.customerId
WHERE o.id IS NULL;
```

---

## Related Problems

- LeetCode #181 — Employees Earning More Than Their Managers (hash join; this problem is the anti-join twin)
- LeetCode #182 — Duplicate Emails (same table-modelling + hash-map pattern)
- LeetCode #196 — Delete Duplicate Emails (set-based row filtering)
- LeetCode #262 — Trips and Users (join + filtering on membership)
- LeetCode #577 — Employee Bonus (LEFT JOIN with NULL filter — the SQL twin of this pattern)
