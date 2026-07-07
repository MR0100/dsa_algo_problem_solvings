package main

import (
	"fmt"
	"sort"
)

// LeetCode 182 — Duplicate Emails.
//
// The original problem is a SQL one: given the Person table
//
//	+----+---------+
//	| id | email   |
//	+----+---------+
//
// report every email that appears more than once (each duplicate email exactly
// once, in any order). Here the table is modelled as a slice of Person rows,
// and each approach re-implements the GROUP BY / HAVING COUNT(*) > 1 query.

// Person models one row of the Person table.
type Person struct {
	ID    int
	Email string
}

// ── Approach 1: Brute Force ──────────────────────────────────────────────────
//
// bruteForce solves Duplicate Emails by comparing every row against every
// other row with nested loops.
//
// Intuition:
//
//	An email is a duplicate iff it appears in at least two rows. To report it
//	exactly once without extra memory, let only its FIRST occurrence do the
//	reporting: row i reports its email iff no earlier row holds the same email
//	(i is the first occurrence) and some later row repeats it.
//
// Algorithm:
//  1. For each row i, scan rows 0..i-1; if any holds the same email, skip i
//     (an earlier occurrence owns the report).
//  2. Otherwise scan rows i+1..n-1; on the first equal email, report it once.
//
// Time:  O(n^2) — up to two linear scans for each of the n rows.
// Space: O(1) — no auxiliary structures beyond the output slice.
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

// ── Approach 2: Sorting ──────────────────────────────────────────────────────
//
// sorting solves Duplicate Emails by sorting the emails so duplicates become
// adjacent, then scanning for runs of length >= 2.
//
// Intuition:
//
//	After sorting, all copies of an email sit next to each other. A duplicate
//	is any run of length >= 2 — and we can report it exactly once by reporting
//	only at the FIRST adjacency of each run (the boundary where a run starts).
//
// Algorithm:
//  1. Extract the email column into a slice and sort it.
//  2. Scan adjacent pairs: report emails[i] when emails[i] == emails[i-1] and
//     that pair starts the run (i == 1 or emails[i-2] differs).
//
// Time:  O(n log n) — sorting dominates the single linear scan.
// Space: O(n) — the extracted email column (input rows stay untouched).
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

// ── Approach 3: Hash Map Counting (Optimal) ──────────────────────────────────
//
// hashMap solves Duplicate Emails with a single counting pass over an
// email → occurrences map.
//
// Intuition:
//
//	SQL's `GROUP BY email HAVING COUNT(email) > 1` is frequency counting.
//	One pass builds the counts; emitting an email at the exact moment its
//	count reaches 2 reports every duplicate exactly once with no second pass.
//
// Algorithm:
//  1. Walk the rows incrementing counts[email].
//  2. When counts[email] becomes exactly 2, append the email (== 2, not >= 2,
//     so an email seen 3+ times is still reported only once).
//
// Time:  O(n) — one pass with O(1) average map updates.
// Space: O(n) — at most one map entry per distinct email.
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

// sorted returns an alphabetically sorted copy of emails so that the
// "return the result table in any order" outputs print deterministically.
func sorted(emails []string) []string {
	out := make([]string, len(emails))
	copy(out, emails)
	sort.Strings(out)
	return out
}

func main() {
	// Example 1 — Person table:
	//   | 1 | a@b.com |
	//   | 2 | c@d.com |
	//   | 3 | a@b.com |
	people := []Person{
		{ID: 1, Email: "a@b.com"},
		{ID: 2, Email: "c@d.com"},
		{ID: 3, Email: "a@b.com"},
	}

	fmt.Println("=== Approach 1: Brute Force ===")
	fmt.Println(sorted(bruteForce(people))) // [a@b.com]

	fmt.Println("=== Approach 2: Sorting ===")
	fmt.Println(sorted(sorting(people))) // [a@b.com]

	fmt.Println("=== Approach 3: Hash Map Counting (Optimal) ===")
	fmt.Println(sorted(hashMap(people))) // [a@b.com]
}
