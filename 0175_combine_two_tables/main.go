package main

import (
	"fmt"
	"sort"
)

// Person mirrors the LeetCode `Person` table (personId is the primary key).
type Person struct {
	PersonID  int
	LastName  string
	FirstName string
}

// Address mirrors the LeetCode `Address` table (addressId is the primary key;
// personId here is a foreign key that may or may not match a Person row).
type Address struct {
	AddressID int
	PersonID  int
	City      string
	State     string
}

// Row is one line of the result set: firstName, lastName, city, state
// ("null" marks a person with no address, exactly like SQL's NULL).
type Row struct {
	FirstName, LastName, City, State string
}

// This problem is officially a SQL exercise; the accepted answer is:
//
//	SELECT p.firstName, p.lastName, a.city, a.state
//	FROM Person p LEFT JOIN Address a ON p.personId = a.personId;
//
// Below we implement what the database engine actually does for that LEFT
// JOIN, as three classic join algorithms: nested-loop, sort-merge, and hash.

// ── Approach 1: Brute Force (Nested Loop Join) ───────────────────────────────
//
// nestedLoopJoin solves Combine Two Tables by comparing every person against
// every address — the textbook O(P×A) nested-loop LEFT JOIN.
//
// Intuition:
//
//	A LEFT JOIN keeps EVERY row of the left table (Person). For each person,
//	scan the whole Address table for rows with the same personId: emit one
//	output row per match; if the scan finds nothing, still emit the person
//	once, padding city/state with null. This is literally how a database
//	executes a join when it has no index and no memory for a hash table.
//
// Algorithm:
//  1. For each person p (preserving Person-table order):
//  2. Scan all addresses; for each with a.personId == p.personId, emit
//     (p.firstName, p.lastName, a.city, a.state) and mark matched.
//  3. If no address matched, emit (p.firstName, p.lastName, null, null).
//
// Time:  O(P·A) — every person scans the entire address list.
// Space: O(1) extra — output aside, only a boolean flag per person scan.
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

// ── Approach 2: Sort-Merge Join ──────────────────────────────────────────────
//
// sortMergeJoin solves Combine Two Tables by sorting both tables on personId
// and sweeping them together with two pointers.
//
// Intuition:
//
//	Once both tables are ordered by the join key, matching rows must line up:
//	advance whichever pointer lags. When keys are equal, emit the pairing;
//	when the address key is smaller, that address matches nobody — skip it;
//	when the person key is smaller, that person has no address — emit a null
//	row. This is the join databases pick when inputs are already sorted or
//	must be sorted anyway (e.g. for ORDER BY on the join key).
//
// Algorithm:
//  1. Copy and sort persons by personId, addresses by personId.
//  2. Two pointers i (persons), j (addresses):
//     a.PersonID < p.PersonID → j++ (orphan address);
//     equal → emit one row per equal-key address, then i++;
//     otherwise → emit null row, i++.
//  3. Remaining persons after addresses run out get null rows.
//
// Time:  O(P log P + A log A) — sorting dominates; the merge sweep is O(P+A).
// Space: O(P + A) — sorted copies (we must not mutate the caller's tables).
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

// ── Approach 3: Hash Join (Optimal) ──────────────────────────────────────────
//
// hashJoin solves Combine Two Tables by indexing the Address table in a hash
// map keyed on personId, then probing it once per person.
//
// Intuition:
//
//	The nested loop wastes time re-scanning addresses. Build the "index"
//	ourselves: one pass over Address fills map[personId] → addresses; then
//	each person finds its matches in O(1) average. Build + probe = one pass
//	over each table. This is the hash join a query planner chooses for
//	unsorted, unindexed tables — and the map is exactly what a database
//	index on Address.personId would give us.
//
// Algorithm:
//  1. byPerson = map[int][]Address; append every address under its personId
//     (a slice value keeps LEFT JOIN correct if a person has several addresses).
//  2. For each person (preserving Person-table order): look up byPerson[id];
//     emit one row per hit, or a null row when the bucket is empty.
//
// Time:  O(P + A) — one pass to build the map, one pass to probe it.
// Space: O(A) — the hash map holds every address once.
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

// printRows renders a result set in the LeetCode output-table format.
func printRows(rows []Row) {
	for _, r := range rows {
		fmt.Printf("| %-9s | %-8s | %-13s | %-8s |\n", r.FirstName, r.LastName, r.City, r.State)
	}
}

func main() {
	// Example 1 input tables.
	persons := []Person{
		{PersonID: 1, LastName: "Wang", FirstName: "Allen"},
		{PersonID: 2, LastName: "Alice", FirstName: "Bob"},
	}
	addresses := []Address{
		{AddressID: 1, PersonID: 2, City: "New York City", State: "New York"},
		{AddressID: 2, PersonID: 3, City: "Leetcode", State: "California"},
	}

	// Expected (any row order is accepted; we emit Person-table order):
	// | Allen     | Wang     | null          | null     |
	// | Bob       | Alice    | New York City | New York |

	fmt.Println("=== Approach 1: Brute Force (Nested Loop Join) ===")
	printRows(nestedLoopJoin(persons, addresses)) // Allen/Wang/null/null then Bob/Alice/New York City/New York

	fmt.Println("=== Approach 2: Sort-Merge Join ===")
	printRows(sortMergeJoin(persons, addresses)) // same two rows (sorted by personId: 1 then 2)

	fmt.Println("=== Approach 3: Hash Join (Optimal) ===")
	printRows(hashJoin(persons, addresses)) // same two rows
}
