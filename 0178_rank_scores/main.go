package main

import (
	"fmt"
	"math"
	"sort"
	"strings"
)

// Score models one row of the LeetCode `Scores` table:
//
//	+-------------+---------+
//	| Column Name | Type    |
//	+-------------+---------+
//	| id          | int     |
//	| score       | decimal |
//	+-------------+---------+
//
// This is a Database problem; the repo solves it in Go by loading the table
// into a slice of rows and implementing DENSE_RANK() by hand.
type Score struct {
	ID    int
	Score float64
}

// RankedScore is one row of the result table: the score plus its dense rank.
type RankedScore struct {
	Score float64
	Rank  int
}

// toCents converts a two-decimal score to an exact integer (3.85 → 385) so
// that equality and ordering never hit float64 precision traps.
func toCents(score float64) int {
	return int(math.Round(score * 100)) // round kills representation noise
}

// sortedByScoreDesc returns a copy of the rows ordered by score descending —
// the `ORDER BY score DESC` every approach needs for its output.
func sortedByScoreDesc(scores []Score) []Score {
	rows := make([]Score, len(scores))
	copy(rows, scores) // never mutate the caller's table
	sort.SliceStable(rows, func(i, j int) bool {
		return toCents(rows[i].Score) > toCents(rows[j].Score) // higher first
	})
	return rows
}

// formatRankings renders the result table on one line: [(score, rank) ...].
func formatRankings(rows []RankedScore) string {
	parts := make([]string, len(rows))
	for i, r := range rows {
		parts[i] = fmt.Sprintf("(%.2f, %d)", r.Score, r.Rank)
	}
	return "[" + strings.Join(parts, " ") + "]"
}

// ── Approach 1: Brute Force (Count Greater Distinct) ─────────────────────────
//
// bruteForce solves Rank Scores by computing each row's rank from first
// principles: one plus the number of DISTINCT scores strictly greater.
//
// Intuition:
//
//	The dense rank of a score is, by definition, how many distinct higher
//	scores exist above it, plus one: nothing higher → rank 1, one distinct
//	higher value → rank 2, and so on. Ties share a rank automatically
//	(identical scores see the identical set of higher values) and no rank is
//	ever skipped (adding one more distinct higher value raises the count by
//	exactly one). So just count, per row, with a nested scan.
//
// Algorithm:
//  1. Copy the rows and sort them by score descending (the output order).
//  2. For every row, scan the whole table collecting the distinct scores
//     strictly greater than it into a small hash set.
//  3. The row's rank is len(set) + 1; emit (score, rank).
//
// Time:  O(n²) — a full O(n) scan for each of the n rows (plus the sort).
// Space: O(n) — the sorted copy and the per-row set of greater values.
func bruteForce(scores []Score) []RankedScore {
	rows := sortedByScoreDesc(scores) // output must be ORDER BY score DESC
	result := make([]RankedScore, 0, len(rows))
	for _, row := range rows {
		myCents := toCents(row.Score)
		greater := map[int]bool{} // DISTINCT scores strictly above this row
		for _, other := range scores {
			if c := toCents(other.Score); c > myCents {
				greater[c] = true // set membership dedupes ties among highers
			}
		}
		// Dense rank = number of distinct better scores + 1.
		result = append(result, RankedScore{Score: row.Score, Rank: len(greater) + 1})
	}
	return result
}

// ── Approach 2: Hash Map Dense Rank (Coordinate Compression) ─────────────────
//
// hashMapDenseRank solves Rank Scores by precomputing every distinct score's
// rank once, then labelling rows via O(1) map lookups.
//
// Intuition:
//
//	All rows with the same score share one rank, so compute ranks per VALUE,
//	not per row: deduplicate the scores, sort the distinct values descending,
//	and the rank of each value is simply its position + 1 in that list. This
//	is coordinate compression (exactly #1331 Rank Transform of an Array) and
//	is how DENSE_RANK() is usually reasoned about: rank = index among the
//	sorted distinct values.
//
// Algorithm:
//  1. Collect the distinct score values with a hash set.
//  2. Sort the distinct values descending; build rankOf[value] = index + 1.
//  3. Sort the rows by score descending (output order).
//  4. Emit every row with rankOf[its score] — a constant-time lookup.
//
// Time:  O(n log n) — sorting rows dominates (d ≤ n distinct values).
// Space: O(n) — the distinct slice plus the value → rank map.
func hashMapDenseRank(scores []Score) []RankedScore {
	// Step 1: DISTINCT score values (as exact cents).
	seen := map[int]bool{}
	distinct := []int{}
	for _, s := range scores {
		if c := toCents(s.Score); !seen[c] {
			seen[c] = true
			distinct = append(distinct, c)
		}
	}
	// Step 2: descending order → position i holds the (i+1)-th best value.
	sort.Sort(sort.Reverse(sort.IntSlice(distinct)))
	rankOf := make(map[int]int, len(distinct)) // score cents → dense rank
	for i, c := range distinct {
		rankOf[c] = i + 1 // best value → 1, next distinct → 2, ...
	}
	// Steps 3–4: output rows in score order with their precomputed rank.
	rows := sortedByScoreDesc(scores)
	result := make([]RankedScore, 0, len(rows))
	for _, row := range rows {
		result = append(result, RankedScore{Score: row.Score, Rank: rankOf[toCents(row.Score)]})
	}
	return result
}

// ── Approach 3: Sort and Sweep (Optimal) ─────────────────────────────────────
//
// sortAndSweep solves Rank Scores in a single sweep over the sorted rows —
// the way the DENSE_RANK() window function actually evaluates.
//
// Intuition:
//
//	Once the rows are sorted descending, dense ranks are forced: the first
//	row is rank 1, and walking down, the rank increases by exactly 1 every
//	time the score CHANGES — equal neighbours inherit the same rank (ties),
//	and because the increment is always +1 there can be no holes. One
//	counter and one "previous value" variable replace all auxiliary maps.
//
// Algorithm:
//  1. Copy the rows and sort by score descending.
//  2. rank = 0, prev = sentinel (no previous value).
//  3. For each row top-down: if its score differs from prev, increment rank
//     and update prev; emit (score, rank).
//
// Time:  O(n log n) — the sort; the sweep itself is O(n).
// Space: O(1) — beyond the sorted copy and the output, just two scalars.
func sortAndSweep(scores []Score) []RankedScore {
	rows := sortedByScoreDesc(scores)
	result := make([]RankedScore, 0, len(rows))
	rank := 0           // dense rank of the previous distinct value
	prev := math.MinInt // sentinel: no score seen yet
	for _, row := range rows {
		c := toCents(row.Score)
		if c != prev {
			rank++   // new distinct value → next consecutive rank (no holes)
			prev = c // remember it so following ties reuse this rank
		}
		result = append(result, RankedScore{Score: row.Score, Rank: rank})
	}
	return result
}

func main() {
	// Example 1: Scores = [(1,3.50), (2,3.65), (3,4.00), (4,3.85), (5,4.00), (6,3.65)]
	example1 := []Score{
		{ID: 1, Score: 3.50},
		{ID: 2, Score: 3.65},
		{ID: 3, Score: 4.00},
		{ID: 4, Score: 3.85},
		{ID: 5, Score: 4.00},
		{ID: 6, Score: 3.65},
	}
	// Edge: single row — must simply be rank 1.
	edge := []Score{{ID: 1, Score: 3.50}}

	expected1 := "[(4.00, 1) (4.00, 1) (3.85, 2) (3.65, 3) (3.65, 3) (3.50, 4)]"
	expectedEdge := "[(3.50, 1)]"

	fmt.Println("=== Approach 1: Brute Force (Count Greater Distinct) ===")
	fmt.Printf("Example 1: got=%s\n", formatRankings(bruteForce(example1))) // expected [(4.00, 1) (4.00, 1) (3.85, 2) (3.65, 3) (3.65, 3) (3.50, 4)]
	fmt.Printf("           exp=%s\n", expected1)
	fmt.Printf("Edge:      got=%s  exp=%s\n", formatRankings(bruteForce(edge)), expectedEdge) // expected [(3.50, 1)]

	fmt.Println("=== Approach 2: Hash Map Dense Rank (Coordinate Compression) ===")
	fmt.Printf("Example 1: got=%s\n", formatRankings(hashMapDenseRank(example1))) // expected [(4.00, 1) (4.00, 1) (3.85, 2) (3.65, 3) (3.65, 3) (3.50, 4)]
	fmt.Printf("           exp=%s\n", expected1)
	fmt.Printf("Edge:      got=%s  exp=%s\n", formatRankings(hashMapDenseRank(edge)), expectedEdge) // expected [(3.50, 1)]

	fmt.Println("=== Approach 3: Sort and Sweep (Optimal) ===")
	fmt.Printf("Example 1: got=%s\n", formatRankings(sortAndSweep(example1))) // expected [(4.00, 1) (4.00, 1) (3.85, 2) (3.65, 3) (3.65, 3) (3.50, 4)]
	fmt.Printf("           exp=%s\n", expected1)
	fmt.Printf("Edge:      got=%s  exp=%s\n", formatRankings(sortAndSweep(edge)), expectedEdge) // expected [(3.50, 1)]
}
