# 0262 — Trips and Users

> LeetCode #262 · Difficulty: Hard
> **Categories:** Database, SQL, Join, Group By, Aggregation

---

## Problem Statement

Table: `Trips`

```
+-------------+----------+
| Column Name | Type     |
+-------------+----------+
| id          | int      |
| client_id   | int      |
| driver_id   | int      |
| city_id     | int      |
| status      | enum     |
| request_at  | varchar  |
+-------------+----------+
id is the primary key (column with unique values) for this table.
The table holds all taxi trips. Each trip has a unique id, while client_id and
driver_id are foreign keys to the users_id at the Users table.
status is an ENUM (category) type of ('completed', 'cancelled_by_driver',
'cancelled_by_client').
```

Table: `Users`

```
+-------------+----------+
| Column Name | Type     |
+-------------+----------+
| users_id    | int      |
| banned      | enum     |
| role        | enum     |
+-------------+----------+
users_id is the primary key (column with unique values) for this table.
The table holds all users. Each user has a unique users_id, and role is an ENUM
type of ('client', 'driver', 'partner').
banned is an ENUM (category) type of ('Yes', 'No').
```

The **cancellation rate** is computed by dividing the number of canceled (by client or driver) requests with unbanned users by the total number of requests with unbanned users on that day.

Write a solution to find the cancellation rate of requests with unbanned users (both client and driver must not be banned) each day between `"2013-10-01"` and `"2013-10-03"`. Round `Cancellation Rate` to two decimal points.

Return the result table in **any order**.

**Example 1:**

```
Input:
Trips table:
+----+-----------+-----------+---------+---------------------+------------+
| id | client_id | driver_id | city_id | status              | request_at |
+----+-----------+-----------+---------+---------------------+------------+
| 1  | 1         | 10        | 1       | completed           | 2013-10-01 |
| 2  | 2         | 11        | 1       | cancelled_by_driver | 2013-10-01 |
| 3  | 3         | 12        | 6       | completed           | 2013-10-01 |
| 4  | 4         | 13        | 6       | cancelled_by_client | 2013-10-01 |
| 5  | 1         | 10        | 1       | completed           | 2013-10-02 |
| 6  | 2         | 11        | 6       | completed           | 2013-10-02 |
| 7  | 3         | 12        | 6       | completed           | 2013-10-02 |
| 8  | 2         | 12        | 12      | completed           | 2013-10-03 |
| 9  | 3         | 10        | 12      | completed           | 2013-10-03 |
| 10 | 4         | 13        | 12      | cancelled_by_driver | 2013-10-03 |
+----+-----------+-----------+---------+---------------------+------------+

Users table:
+----------+--------+--------+
| users_id | banned | role   |
+----------+--------+--------+
| 1        | No     | client |
| 2        | Yes    | client |
| 3        | No     | client |
| 4        | No     | client |
| 10       | No     | driver |
| 11       | No     | driver |
| 12       | No     | driver |
| 13       | No     | driver |
+----------+--------+--------+

Output:
+------------+-------------------+
| Day        | Cancellation Rate |
+------------+-------------------+
| 2013-10-01 | 0.33              |
| 2013-10-02 | 0.00              |
| 2013-10-03 | 0.50              |
+------------+-------------------+

Explanation:
On 2013-10-01:
  - There were 4 requests in total, 2 of which were canceled.
  - However, the request with id=2 was made by a banned client (user_id=2), so
    it is ignored in the calculation.
  - Hence there are 3 unbanned requests in total, 1 of which was canceled.
  - The Cancellation Rate is (1 / 3) = 0.33
On 2013-10-02:
  - There were 3 requests in total, 0 of which were canceled.
  - The request with id=6 was made by a banned client, so it is ignored.
  - Hence there are 2 unbanned requests in total, 0 of which were canceled.
  - The Cancellation Rate is (0 / 2) = 0.00
On 2013-10-03:
  - There were 3 requests in total, 1 of which was canceled.
  - The request with id=8 was made by a banned client, so it is ignored.
  - Hence there are 2 unbanned requests in total, 1 of which was canceled.
  - The Cancellation Rate is (1 / 2) = 0.50
```

---

## Company Frequency

| Company   | Frequency        | Last Reported |
|-----------|------------------|---------------|
| Amazon    | ★★★☆☆ Medium     | 2023          |
| Meta      | ★★★☆☆ Medium     | 2023          |
| Uber      | ★★★★☆ High       | 2022          |
| Google    | ★★☆☆☆ Low        | 2022          |

> ⚠️ Frequency data is crowd-sourced from LeetCode Discuss, Glassdoor, and community interview reports. Treat as a signal, not a guarantee.

---

## DSA Concepts Used

- **Hash Join** — index the Users table by `users_id` for O(1) client/driver lookups → see [`/dsa/hash_map.md`](/dsa/hash_map.md)
- **Group-By Aggregation** — bucket qualifying trips per day and accumulate counts → see [`/dsa/hash_map.md`](/dsa/hash_map.md)
- **Sorting** — order the output rows by day → see [`/dsa/sorting.md`](/dsa/sorting.md)

---

## Approaches Overview

| # | Approach | Time | Space | When to use |
|---|----------|------|-------|-------------|
| 1 | Nested-Loop Join + Group-By | O(T·U) | O(D) | Tiny tables; the literal join simulation |
| 2 | Hash Join + Group-By (Optimal) | O(T + U) | O(U + D) | Any real size; O(1) user lookups |

*(T = trips, U = users, D = distinct days.)*

---

## Approach 1 — Nested-Loop Join + Group-By

### Intuition
The two SQL `JOIN`s attach the client and driver User rows to each trip. Reproduce them literally: for each in-window trip, linear-scan Users to find both parties and keep the trip only if **both** are unbanned. Then bucket kept trips by `request_at` and, per day, the rate is (non-completed)/(total).

### Algorithm
1. For each trip with `request_at` in `[2013-10-01, 2013-10-03]`:
2. Linear-scan Users to look up the client and the driver rows.
3. If either is missing or banned, skip the trip.
4. Else increment that day's `total`, and its `cancelled` if `status != "completed"`.
5. For each day compute `cancelled/total`, round to 2 dp, sort days ascending.

### Complexity
- **Time:** O(T·U) — each kept trip rescans the whole Users table.
- **Space:** O(D) — one counter bucket per distinct day.

### Code
```go
func nestedLoopJoin(trips []Trip, users []User) []Result {
	total := map[string]int{}     // day → count of qualifying trips
	cancelled := map[string]int{} // day → count of qualifying non-completed trips
	inRange := func(d string) bool {
		return d >= "2013-10-01" && d <= "2013-10-03" // lexical works for ISO dates
	}
	bannedOf := func(id int) (string, bool) { // linear scan lookup: (banned, found)
		for _, u := range users {
			if u.UsersID == id {
				return u.Banned, true
			}
		}
		return "", false
	}
	for _, t := range trips {
		if !inRange(t.RequestAt) { // WHERE request_at BETWEEN ...
			continue
		}
		cb, cok := bannedOf(t.ClientID) // JOIN Users c ON client_id
		db, dok := bannedOf(t.DriverID) // JOIN Users d ON driver_id
		if !cok || !dok || cb != "No" || db != "No" {
			continue // drop trips where either party is missing or banned
		}
		total[t.RequestAt]++
		if t.Status != "completed" { // status != 'completed' ⇒ a cancellation
			cancelled[t.RequestAt]++
		}
	}
	return finalize(total, cancelled)
}
```

### Dry Run
`2013-10-01` trips (ids 1–4). Users 1,3,4,10–13 unbanned; user **2 is banned**.

| Trip id | client | driver | banned party? | kept? | status → cancel? | total(10-01) | cancelled(10-01) |
|---------|--------|--------|---------------|-------|------------------|--------------|------------------|
| 1       | 1      | 10     | none          | yes   | completed → no   | 1            | 0                |
| 2       | 2      | 11     | client 2      | **no**| —                | 1            | 0                |
| 3       | 3      | 12     | none          | yes   | completed → no   | 2            | 0                |
| 4       | 4      | 13     | none          | yes   | cancelled → yes  | 3            | 1                |

Rate for 2013-10-01 = `round2(1/3)` = **0.33**.

---

## Approach 2 — Hash Join + Group-By (Optimal)

### Intuition
The nested-loop version rescans Users for every trip. Index Users once into a map `users_id → banned`; then a single pass over the trips does two O(1) lookups each and accumulates the per-day totals, matching the accepted SQL's hash-join execution plan.

### Algorithm
1. Build `banned[users_id] = "Yes"/"No"` from the Users table.
2. For each in-window trip, look up client and driver in O(1); keep only if both present and `"No"`.
3. Accumulate `total` and `cancelled` per day.
4. Compute `cancelled/total` rounded to 2 dp, sorted by day.

### Complexity
- **Time:** O(T + U) — one pass to index users, one pass over trips.
- **Space:** O(U + D) — the user index plus per-day buckets.

### Code
```go
func hashJoin(trips []Trip, users []User) []Result {
	banned := make(map[int]string, len(users)) // users_id → "Yes"/"No"
	for _, u := range users {
		banned[u.UsersID] = u.Banned // index the whole Users table once
	}
	total := map[string]int{}
	cancelled := map[string]int{}
	for _, t := range trips {
		if t.RequestAt < "2013-10-01" || t.RequestAt > "2013-10-03" {
			continue // outside the reporting window
		}
		cb, cok := banned[t.ClientID] // O(1) client lookup
		db, dok := banned[t.DriverID] // O(1) driver lookup
		if !cok || !dok || cb != "No" || db != "No" {
			continue // either party unknown or banned ⇒ exclude
		}
		total[t.RequestAt]++
		if t.Status != "completed" {
			cancelled[t.RequestAt]++
		}
	}
	return finalize(total, cancelled)
}
```

### Dry Run
Index: `banned = {1:No, 2:Yes, 3:No, 4:No, 10:No, 11:No, 12:No, 13:No}`.
Processing `2013-10-03` trips (ids 8–10):

| Trip id | client | driver | lookup client / driver | kept? | status → cancel? | total | cancelled |
|---------|--------|--------|------------------------|-------|------------------|-------|-----------|
| 8       | 2      | 12     | Yes / No               | **no**| —                | 0     | 0         |
| 9       | 3      | 10     | No / No                | yes   | completed → no   | 1     | 0         |
| 10      | 4      | 13     | No / No                | yes   | cancelled → yes  | 2     | 1         |

Rate for 2013-10-03 = `round2(1/2)` = **0.50**.

---

## Key Takeaways

- **The banned filter is the trap:** a trip counts only if *both* its client and driver are unbanned. A trip made by a banned party is dropped from **both** numerator and denominator, not just the numerator.
- Cancellation = `status != 'completed'`, which is why `SUM(CASE WHEN status != 'completed' ...)` (or `status LIKE 'cancelled%'`) is cleaner than enumerating the two cancel enums.
- ISO date strings (`YYYY-MM-DD`) compare correctly with plain lexical `<`/`>`, so no date parsing is needed for the `BETWEEN` window.
- Modelling SQL as a hash-join + group-by makes the query's execution plan explicit and turns an O(T·U) nested loop into O(T + U).

---

## Related Problems

- LeetCode #1075 — Project Employees I (group-by + rounded average)
- LeetCode #1211 — Queries Quality and Percentage (conditional aggregation with rounding)
- LeetCode #1633 — Percentage of Users Attended a Contest (ratio over a filtered set)
- LeetCode #0175 — Combine Two Tables (join modelled in Go)
