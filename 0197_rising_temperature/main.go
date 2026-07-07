package main

import (
	"fmt"
	"sort"
	"time"
)

// Weather models one row of the Weather table.
// ID is unique, and no two rows share the same RecordDate.
type Weather struct {
	ID          int
	RecordDate  string // "YYYY-MM-DD"
	Temperature int
}

// dateLayout is Go's reference layout for the "YYYY-MM-DD" date format.
const dateLayout = "2006-01-02"

// mustDate parses a "YYYY-MM-DD" string into a time.Time (UTC midnight).
// The table schema guarantees valid dates, so a parse error is unreachable.
func mustDate(s string) time.Time {
	t, err := time.Parse(dateLayout, s)
	if err != nil {
		panic(err) // unreachable with schema-valid data
	}
	return t
}

// isNextDay reports whether curr falls exactly one calendar day after prev.
// AddDate handles month/year boundaries and leap days correctly, which is the
// whole subtlety of this problem — ids/dates may have gaps, so "previous row"
// is NOT automatically "yesterday".
func isNextDay(prev, curr string) bool {
	return mustDate(prev).AddDate(0, 0, 1).Equal(mustDate(curr))
}

// ── Approach 1: Brute Force ──────────────────────────────────────────────────
//
// bruteForce solves Rising Temperature by pairing every row with every other
// row, looking for the row dated exactly one day earlier.
//
// Intuition:
//   For each day we need its "yesterday" row. With no index at hand, simply
//   scan the whole table for a row whose date is exactly one day before.
//   This is the nested-loop version of the SQL self-join:
//     SELECT w1.id FROM Weather w1, Weather w2
//     WHERE DATEDIFF(w1.recordDate, w2.recordDate) = 1
//       AND w1.temperature > w2.temperature;
//
// Algorithm:
//   1. For every row "today", scan all rows for one dated today − 1 day.
//   2. If that yesterday row exists and today's temperature is strictly
//      higher, collect today's id (dates are unique → at most one match).
//   3. Sort the collected ids (any order is accepted; we normalise).
//
// Time:  O(n²) — every row scans the entire table once.
// Space: O(1) auxiliary — only the output slice is allocated.
func bruteForce(weather []Weather) []int {
	ids := []int{}
	for _, today := range weather {
		for _, other := range weather {
			// other qualifies only when dated exactly one day before today.
			if isNextDay(other.RecordDate, today.RecordDate) {
				if today.Temperature > other.Temperature {
					ids = append(ids, today.ID)
				}
				break // dates are unique → at most one "yesterday" exists
			}
		}
	}
	sort.Ints(ids) // deterministic output
	return ids
}

// ── Approach 2: Sort + Scan ──────────────────────────────────────────────────
//
// sortAndScan solves Rising Temperature by sorting a copy chronologically and
// comparing each row only with its immediate predecessor.
//
// Intuition:
//   After sorting by date, a row's potential "yesterday" can only be the row
//   directly before it — anything earlier is at least two days away. This is
//   the window-function idea (LAG(...) OVER (ORDER BY recordDate)), including
//   its classic pitfall: adjacent rows may still be more than one day apart,
//   so the calendar gap must be verified explicitly.
//
// Algorithm:
//   1. Copy the table and sort it by recordDate ascending.
//   2. For each i ≥ 1, check that row i-1 is exactly one calendar day earlier
//      AND that row i's temperature is strictly higher; collect row i's id.
//   3. Sort the collected ids.
//
// Time:  O(n log n) — dominated by the sort; the scan is O(n).
// Space: O(n) — the sorted working copy.
func sortAndScan(weather []Weather) []int {
	rows := make([]Weather, len(weather))
	copy(rows, weather) // never reorder the caller's table

	// "YYYY-MM-DD" strings sort lexicographically in chronological order,
	// so a plain string comparison is a correct date comparison here.
	sort.Slice(rows, func(i, j int) bool { return rows[i].RecordDate < rows[j].RecordDate })

	ids := []int{}
	for i := 1; i < len(rows); i++ {
		// Both conditions required: consecutive calendar days AND a rise.
		if isNextDay(rows[i-1].RecordDate, rows[i].RecordDate) &&
			rows[i].Temperature > rows[i-1].Temperature {
			ids = append(ids, rows[i].ID)
		}
	}
	sort.Ints(ids) // deterministic output
	return ids
}

// ── Approach 3: Hash Map (Optimal) ───────────────────────────────────────────
//
// hashMap solves Rising Temperature by indexing every row by its date, then
// looking each row's yesterday up in O(1).
//
// Intuition:
//   The brute force wastes time searching for yesterday; a map keyed by date
//   answers "who was yesterday?" instantly. This is exactly what the database
//   engine does when it executes the self-join with a hash join.
//
// Algorithm:
//   1. Pass 1: build map recordDate → row.
//   2. Pass 2: for each row, compute yesterday's date string with AddDate
//      (calendar-safe) and look it up; if present with a strictly lower
//      temperature, collect the id.
//   3. Sort the collected ids.
//
// Time:  O(n) — two linear passes, O(1) work per row.
// Space: O(n) — the date-indexed map.
func hashMap(weather []Weather) []int {
	// Index every row by its date for O(1) "who was yesterday?" lookups.
	byDate := make(map[string]Weather, len(weather))
	for _, row := range weather {
		byDate[row.RecordDate] = row
	}

	ids := []int{}
	for _, today := range weather {
		// Calendar-correct yesterday (handles month/year edges, leap days).
		yesterday := mustDate(today.RecordDate).AddDate(0, 0, -1).Format(dateLayout)
		if prev, ok := byDate[yesterday]; ok && today.Temperature > prev.Temperature {
			ids = append(ids, today.ID)
		}
	}
	sort.Ints(ids) // deterministic output
	return ids
}

func main() {
	// Example 1 (the only official example): Weather table.
	weather := []Weather{
		{1, "2015-01-01", 10},
		{2, "2015-01-02", 25},
		{3, "2015-01-03", 20},
		{4, "2015-01-04", 30},
	}

	fmt.Println("=== Approach 1: Brute Force ===")
	fmt.Println(bruteForce(weather)) // [2 4]

	fmt.Println("=== Approach 2: Sort + Scan ===")
	fmt.Println(sortAndScan(weather)) // [2 4]

	fmt.Println("=== Approach 3: Hash Map (Optimal) ===")
	fmt.Println(hashMap(weather)) // [2 4]
}
