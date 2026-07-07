package main

import (
	"fmt"
	"sort"
)

// LeetCode 180 — Consecutive Numbers.
//
// The original problem is a SQL one: given the table
//
//	Logs: +----+-----+
//	      | id | num |
//	      +----+-----+
//
// where `id` is an autoincrement primary key (sequential here), find every
// value of `num` that appears in AT LEAST THREE CONSECUTIVE rows (by ascending
// id). Report each such value once as `ConsecutiveNums`, in any order.
//
// The accepted SQL answer (classic three-way self-join on adjacent ids):
//
//	SELECT DISTINCT l1.num AS ConsecutiveNums
//	FROM Logs l1, Logs l2, Logs l3
//	WHERE l1.id = l2.id - 1
//	  AND l2.id = l3.id - 1
//	  AND l1.num = l2.num
//	  AND l2.num = l3.num;
//
// Here the table is modelled as a slice of Log rows and every approach
// re-implements that query as an in-memory "three consecutive equal values"
// scan over the id-ordered rows.

// Log models one row of the Logs table (id is the primary key, num the value).
type Log struct {
	ID  int
	Num int
}

// ── Approach 1: Self-Join Simulation (Brute Force) ───────────────────────────
//
// selfJoin solves Consecutive Numbers by mirroring the SQL three-way
// self-join `Logs l1, l2, l3` with a triple-nested scan.
//
// Intuition:
//
//	The accepted SQL pairs three table aliases whose ids form a run
//	(id, id+1, id+2) and whose nums are all equal. Reproduce it literally:
//	for every triple of rows (i, j, k) check the join predicates
//	l1.id = l2.id-1, l2.id = l3.id-1 and l1.num = l2.num = l3.num, collecting
//	the shared num into a DISTINCT set. This is the join a database runs with
//	no index — correct but O(n^3).
//
// Algorithm:
//  1. For every ordered triple (a, b, c) of rows:
//  2. Check the id chain: b.ID == a.ID+1 and c.ID == b.ID+1.
//  3. Check the value match: a.Num == b.Num == c.Num.
//  4. When all hold, add a.Num to a DISTINCT set.
//  5. Return the set's members, sorted ascending.
//
// Time:  O(n^3) — three nested passes over the rows (mirrors the 3-alias join).
// Space: O(d) — the DISTINCT set of qualifying values (d distinct answers).
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

// ── Approach 2: Adjacent-Triple Scan (Sorted Window) ─────────────────────────
//
// adjacentTriple solves Consecutive Numbers by sorting rows on id and testing
// each contiguous window of three neighbours.
//
// Intuition:
//
//	The self-join only ever pairs rows whose ids differ by exactly 1, so the
//	triple must be three PHYSICALLY adjacent rows once ordered by id. Sort by
//	id, then slide a size-3 window: a window qualifies when its ids are
//	consecutive (i, i+1, i+2) and its three nums are equal. This drops the
//	O(n^3) join to a single O(n) sweep after the sort.
//
// Algorithm:
//  1. Copy and sort the rows ascending by id (don't mutate the caller's slice).
//  2. For each start index t from 0 to n-3, take rows t, t+1, t+2.
//  3. Require consecutive ids (guards id gaps) AND all three nums equal.
//  4. Collect the shared num into a DISTINCT set; return it sorted.
//
// Time:  O(n log n) — the id sort dominates; the window sweep is O(n).
// Space: O(n) — the sorted copy of the rows (plus the O(d) answer set).
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

// ── Approach 3: Single-Pass Running-Streak Counter (Optimal) ─────────────────
//
// runningStreak solves Consecutive Numbers in one pass by counting how long
// the current value has been repeating.
//
// Intuition:
//
//	Walk the id-ordered rows keeping a streak length for the current value.
//	Each row either extends the streak (same value AND id is exactly one more
//	than the previous id) or restarts it at 1. The moment a streak first hits
//	3, that value is a valid answer — record it once. No self-join, no window
//	array: just two counters. This is the query-planner-free O(n) solution.
//
// Algorithm:
//  1. Copy and sort the rows ascending by id.
//  2. Track prevNum, prevID and a streak counter.
//  3. For each row: if num == prevNum AND id == prevID+1, streak++, else
//     streak = 1. On streak reaching exactly 3, mark num in the answer set.
//  4. Update prevNum/prevID; return the set sorted ascending.
//
// Time:  O(n log n) — the id sort; the streak scan itself is O(n) / O(n) if
//
//	the input is already id-sorted (as LeetCode guarantees).
//
// Space: O(1) extra beyond the sorted copy (plus the O(d) answer set).
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

// sortedKeys returns the keys of a set as an ascending slice, so every
// approach prints its DISTINCT result deterministically.
func sortedKeys(set map[int]bool) []int {
	out := make([]int, 0, len(set))
	for k := range set {
		out = append(out, k)
	}
	sort.Ints(out)
	return out
}

func main() {
	// Example 1 — Logs table:
	//   | id | num |
	//   |  1 |  1  |
	//   |  2 |  1  |
	//   |  3 |  1  |   ← 1 appears three times consecutively (ids 1,2,3)
	//   |  4 |  2  |
	//   |  5 |  1  |
	//   |  6 |  2  |
	//   |  7 |  2  |
	logs := []Log{
		{ID: 1, Num: 1},
		{ID: 2, Num: 1},
		{ID: 3, Num: 1},
		{ID: 4, Num: 2},
		{ID: 5, Num: 1},
		{ID: 6, Num: 2},
		{ID: 7, Num: 2},
	}

	fmt.Println("=== Approach 1: Self-Join Simulation (Brute Force) ===")
	fmt.Println(selfJoin(logs)) // [1]

	fmt.Println("=== Approach 2: Adjacent-Triple Scan (Sorted Window) ===")
	fmt.Println(adjacentTriple(logs)) // [1]

	fmt.Println("=== Approach 3: Single-Pass Running-Streak Counter (Optimal) ===")
	fmt.Println(runningStreak(logs)) // [1]
}
