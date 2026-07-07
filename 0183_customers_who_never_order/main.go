package main

import (
	"fmt"
	"sort"
)

// LeetCode 183 — Customers Who Never Order.
//
// The original problem is a SQL one: given the tables
//
//	Customers: +----+------+        Orders: +----+------------+
//	           | id | name |                | id | customerId |
//	           +----+------+                +----+------------+
//
// report the names of customers who never placed any order (any order).
// Here the tables are modelled as slices of rows, and every approach
// re-implements the anti-join `WHERE id NOT IN (SELECT customerId ...)`.

// Customer models one row of the Customers table.
type Customer struct {
	ID   int
	Name string
}

// Order models one row of the Orders table.
// CustomerID is a foreign key referencing Customers.ID.
type Order struct {
	ID         int
	CustomerID int
}

// ── Approach 1: Brute Force (Nested-Loop Anti-Join) ──────────────────────────
//
// bruteForce solves Customers Who Never Order by scanning the whole Orders
// table once per customer.
//
// Intuition:
//
//	"Never ordered" means no row in Orders references the customer. The naive
//	evaluation of `NOT IN` checks each customer against every order — a
//	nested-loop anti-join: keep the customer only when the inner scan finds
//	no match at all.
//
// Algorithm:
//  1. For every customer c, scan all orders looking for o.CustomerID == c.ID.
//  2. If any order matches, c has ordered — discard and stop the scan early.
//  3. If the scan finishes with no match, append c.Name to the result.
//
// Time:  O(c*o) — a full scan of the o orders for each of the c customers.
// Space: O(1) — no auxiliary structures beyond the output slice.
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

// ── Approach 2: Sort + Binary Search ─────────────────────────────────────────
//
// sortAndBinarySearch solves Customers Who Never Order by sorting the
// customerId column of Orders and binary-searching it per customer.
//
// Intuition:
//
//	The inner scan only asks "does customer id X appear in Orders?" — a
//	membership test. Sorting the customerId column once lets every membership
//	test run in O(log o), like an index on the foreign-key column.
//
// Algorithm:
//  1. Project Orders down to its CustomerID column and sort it.
//  2. For each customer, binary search the sorted ids (sort.SearchInts).
//  3. If the id is absent (search lands past the end or on a different
//     value), the customer never ordered — keep their name.
//
// Time:  O((c+o) log o) — sort the o ids, then c binary searches.
// Space: O(o) — the projected, sorted copy of the customerId column.
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

// ── Approach 3: Hash Set (Hash Anti-Join, Optimal) ───────────────────────────
//
// hashSet solves Customers Who Never Order by materialising the set of
// customer ids that appear in Orders, then filtering customers against it.
//
// Intuition:
//
//	This is how a database executes `NOT IN` efficiently: build a hash set of
//	the subquery's values, then keep the rows whose key misses the set —
//	a hash anti-join. Two linear passes, O(1) membership tests.
//
// Algorithm:
//  1. Pass 1 (build): insert every o.CustomerID into a hash set.
//  2. Pass 2 (probe): for each customer, keep the name iff c.ID is NOT in
//     the set.
//
// Time:  O(c+o) — one pass over each table with O(1) average set operations.
// Space: O(o) — the set holds at most one entry per order row.
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

// sorted returns an alphabetically sorted copy of names so that the
// "return the result table in any order" outputs print deterministically.
func sorted(names []string) []string {
	out := make([]string, len(names))
	copy(out, names)
	sort.Strings(out)
	return out
}

func main() {
	// Example 1 — Customers table:        Orders table:
	//   | 1 | Joe   |                       | 1 | 3 |
	//   | 2 | Henry |                       | 2 | 1 |
	//   | 3 | Sam   |
	//   | 4 | Max   |
	customers := []Customer{
		{ID: 1, Name: "Joe"},
		{ID: 2, Name: "Henry"},
		{ID: 3, Name: "Sam"},
		{ID: 4, Name: "Max"},
	}
	orders := []Order{
		{ID: 1, CustomerID: 3},
		{ID: 2, CustomerID: 1},
	}

	fmt.Println("=== Approach 1: Brute Force (Nested-Loop Anti-Join) ===")
	fmt.Println(sorted(bruteForce(customers, orders))) // [Henry Max]

	fmt.Println("=== Approach 2: Sort + Binary Search ===")
	fmt.Println(sorted(sortAndBinarySearch(customers, orders))) // [Henry Max]

	fmt.Println("=== Approach 3: Hash Set (Hash Anti-Join, Optimal) ===")
	fmt.Println(sorted(hashSet(customers, orders))) // [Henry Max]
}
