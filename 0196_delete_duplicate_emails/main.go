package main

import (
	"fmt"
	"sort"
)

// Person models one row of the Person table.
// ID is the primary key (unique); Email is the value that may repeat across rows.
type Person struct {
	ID    int
	Email string
}

// sortByID orders a result table by id ascending so every approach prints the
// surviving rows in the same deterministic order (LeetCode accepts the final
// table in any order; we normalise purely for easy visual comparison).
func sortByID(rows []Person) []Person {
	sort.Slice(rows, func(i, j int) bool { return rows[i].ID < rows[j].ID })
	return rows
}

// ── Approach 1: Brute Force ──────────────────────────────────────────────────
//
// bruteForce solves Delete Duplicate Emails by comparing every row against
// every other row.
//
// Intuition:
//   A row must be deleted exactly when some other row carries the same email
//   with a smaller id. So for each row, scan the whole table looking for such
//   a "better" row and keep the row only if none exists. This is literally the
//   self-join the classic SQL answer performs:
//     DELETE p1 FROM Person p1, Person p2
//     WHERE p1.email = p2.email AND p1.id > p2.id;
//
// Algorithm:
//   1. For every row i, scan every row j of the table.
//   2. If row j has the same email as row i and a smaller id, row i is a
//      duplicate — mark it and stop scanning.
//   3. Collect every row never marked; sort by id for deterministic output.
//
// Time:  O(n²) — each of the n rows may be checked against all n rows.
// Space: O(n) — only the output slice holding the kept rows.
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

// ── Approach 2: Sort + Sweep ─────────────────────────────────────────────────
//
// sortAndSweep solves Delete Duplicate Emails by sorting a copy of the table
// by (email, id) so duplicates become adjacent, then keeping the first row of
// each email group.
//
// Intuition:
//   Sorting groups equal emails together and, inside each group, places the
//   smallest id first. The keeper of every group is therefore exactly the
//   first row of that group — one linear sweep collects them all.
//
// Algorithm:
//   1. Copy the table (do not disturb the caller's row order).
//   2. Sort the copy by email ascending, then by id ascending.
//   3. Sweep once: keep a row iff it is the very first row or its email
//      differs from the previous row's email.
//   4. Sort survivors by id for deterministic output.
//
// Time:  O(n log n) — dominated by the sort; the sweep is O(n).
// Space: O(n) — the sorted working copy.
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

// ── Approach 3: Hash Map (Optimal) ───────────────────────────────────────────
//
// hashMap solves Delete Duplicate Emails with a single map from email to the
// smallest id seen for that email.
//
// Intuition:
//   The survivor of each email group is fully described by "email → min id".
//   One pass computes that map; a second pass keeps exactly the rows whose id
//   equals their email's minimum. This mirrors the GROUP BY subquery form:
//     DELETE FROM Person WHERE id NOT IN (
//       SELECT * FROM (SELECT MIN(id) FROM Person GROUP BY email) AS keep
//     );
//
// Algorithm:
//   1. Pass 1: for each row, record the minimum id per email in a map.
//   2. Pass 2: keep a row iff minID[row.Email] == row.ID.
//   3. Sort survivors by id for deterministic output.
//
// Time:  O(n) — two linear passes with O(1) map operations each.
// Space: O(n) — the map holds one entry per distinct email.
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

func main() {
	// Example 1 (the only official example):
	// Person table before the delete.
	person := []Person{
		{1, "john@example.com"},
		{2, "bob@example.com"},
		{3, "john@example.com"},
	}

	fmt.Println("=== Approach 1: Brute Force ===")
	fmt.Println(bruteForce(person)) // [{1 john@example.com} {2 bob@example.com}]

	fmt.Println("=== Approach 2: Sort + Sweep ===")
	fmt.Println(sortAndSweep(person)) // [{1 john@example.com} {2 bob@example.com}]

	fmt.Println("=== Approach 3: Hash Map (Optimal) ===")
	fmt.Println(hashMap(person)) // [{1 john@example.com} {2 bob@example.com}]
}
