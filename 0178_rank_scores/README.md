# 0178 — Rank Scores

> LeetCode #178 · Difficulty: Medium
> **Categories:** Database, Sorting, Coordinate Compression

---

## Problem Statement

Table: `Scores`

```
+-------------+---------+
| Column Name | Type    |
+-------------+---------+
| id          | int     |
| score       | decimal |
+-------------+---------+
```

`id` is the primary key (column with unique values) for this table.
Each row of this table contains the score of a game. Score is a floating point value with two decimal places.

Write a solution to find the **rank** of the scores. The ranking should be calculated according to the following rules:

- The scores should be ranked from the highest to the lowest.
- If there is a tie between two scores, both should have the same ranking.
- After a tie, the next ranking number should be the next consecutive integer value. In other words, there should be **no holes** between ranks.

Return the result table ordered by `score` in descending order.

The result format is in the following example.

**Example 1:**

```
Input:
Scores table:
+----+-------+
| id | score |
+----+-------+
| 1  | 3.50  |
| 2  | 3.65  |
| 3  | 4.00  |
| 4  | 3.85  |
| 5  | 4.00  |
| 6  | 3.65  |
+----+-------+
Output:
+-------+------+
| score | rank |
+-------+------+
| 4.00  | 1    |
| 4.00  | 1    |
| 3.85  | 2    |
| 3.65  | 3    |
| 3.65  | 3    |
| 3.50  | 4    |
+-------+------+
```

> **Note:** This is a Database problem (the SQL one-liner is `DENSE_RANK()`).
> Per this repo's convention it is solved in **Go**: the table becomes
> `[]Score{ID, Score}` and each approach implements the dense-rank semantics by
> hand, returning `[]RankedScore{Score, Rank}` ordered by score descending.
> Scores are compared as exact integer *cents* (`3.85 → 385`) to dodge float64
> precision traps. The canonical SQL appears in Key Takeaways.

---

## Company Frequency

| Company    | Frequency        | Last Reported |
|------------|------------------|---------------|
| Amazon     | ★★★★☆ High       | 2024          |
| Google     | ★★★★☆ High       | 2024          |
| Meta       | ★★★☆☆ Medium     | 2024          |
| Microsoft  | ★★★☆☆ Medium     | 2023          |
| Bloomberg  | ★★☆☆☆ Low        | 2023          |
| Adobe      | ★★☆☆☆ Low        | 2023          |

> ⚠️ Frequency data is crowd-sourced from LeetCode Discuss, Glassdoor, and
> community interview reports. Treat as a signal, not a guarantee.

---

## DSA Concepts Used

- **Sorting** — the output must be `ORDER BY score DESC`, and sorted order is also what makes ranks computable in one sweep → see [`/dsa/sorting.md`](/dsa/sorting.md)
- **Hash Map / Hash Set** — deduplicating score values (rank belongs to a *value*, not a row) and O(1) value → rank lookups → see [`/dsa/hash_map.md`](/dsa/hash_map.md)
- **Coordinate Compression (Dense Rank)** — mapping raw values to their index in the sorted distinct list is exactly `DENSE_RANK()` and exactly #1331 Rank Transform of an Array → see [`/dsa/sorting.md`](/dsa/sorting.md)

---

## Approaches Overview

| # | Approach | Time | Space | When to use |
|---|----------|------|-------|-------------|
| 1 | Brute Force (Count Greater Distinct) | O(n²) | O(n) | Definition-first baseline; fine for tiny tables |
| 2 | Hash Map Dense Rank (Coordinate Compression) | O(n log n) | O(n) | When ranks must be *reusable* (labelling many queries/rows by value) |
| 3 | Sort and Sweep (Optimal) | O(n log n) | O(1) extra | Always — single pass after the mandatory sort, no auxiliary map |

---

## Approach 1 — Brute Force (Count Greater Distinct)

### Intuition

Work straight from the definition. A score's dense rank is "how many **distinct** score values beat it, plus one": nothing higher → rank 1; one distinct higher value → rank 2; and so on. Ties automatically share a rank (equal scores see the identical set of higher values), and ranks can't have holes (each extra distinct higher value adds exactly +1). So for every row, scan the whole table and count the distinct greater values with a small set.

### Algorithm

1. Copy the rows and sort by score descending (the required output order).
2. For each row, scan **all** rows; insert every score strictly greater than the current row's score into a hash set (the set dedupes ties among the higher scores).
3. Emit `(score, len(set) + 1)`.

### Complexity

- **Time:** O(n²) — an O(n) counting scan per row dominates the O(n log n) sort.
- **Space:** O(n) — the sorted copy plus the transient per-row set.

### Code

```go
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

// toCents converts a two-decimal score to an exact integer (3.85 → 385) so
// that equality and ordering never hit float64 precision traps.
func toCents(score float64) int {
	return int(math.Round(score * 100)) // round kills representation noise
}
```

### Dry Run

Example 1, rows sorted descending: `4.00, 4.00, 3.85, 3.65, 3.65, 3.50`

| Step | row score | greater set collected | len(set) | rank emitted |
|------|-----------|------------------------|----------|--------------|
| 1 | 4.00 | {} | 0 | 1 |
| 2 | 4.00 | {} | 0 | 1 |
| 3 | 3.85 | {400} | 1 | 2 |
| 4 | 3.65 | {400, 385} | 2 | 3 |
| 5 | 3.65 | {400, 385} | 2 | 3 |
| 6 | 3.50 | {400, 385, 365} | 3 | 4 |

Result: `[(4.00,1) (4.00,1) (3.85,2) (3.65,3) (3.65,3) (3.50,4)]` ✔ — note step 4: the two 4.00 rows collapse to the single set entry `{400}`, which is what keeps the ranking *dense*.

---

## Approach 2 — Hash Map Dense Rank (Coordinate Compression)

### Intuition

Rank is a property of the **value**, not of the row — all rows scoring 3.65 share one rank. So compute each distinct value's rank exactly once: deduplicate, sort the distinct values descending, and value #i in that list has rank i+1. Then labelling the n rows is a mere O(1) map lookup each. This is textbook coordinate compression (identical to #1331 Rank Transform of an Array) and mirrors how you *reason* about `DENSE_RANK()`: rank = 1 + index among sorted distinct values.

### Algorithm

1. Collect distinct score values (as exact cents) with a hash set.
2. Sort the distinct values descending.
3. Build `rankOf[value] = position + 1`.
4. Sort the rows by score descending; emit each row with `rankOf[value]`.

### Complexity

- **Time:** O(n log n) — sorting the n rows dominates; the d ≤ n distinct values sort in O(d log d).
- **Space:** O(n) — the distinct slice and the value → rank map (O(d), worst case d = n).

### Code

```go
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
```

### Dry Run

Example 1: `scores = [3.50, 3.65, 4.00, 3.85, 4.00, 3.65]` (cents: 350, 365, 400, 385, 400, 365)

| Step | Action | State |
|------|--------|-------|
| 1 | dedup values | distinct = [350, 365, 400, 385] |
| 2 | sort descending | distinct = [400, 385, 365, 350] |
| 3 | build rank map | rankOf = {400:1, 385:2, 365:3, 350:4} |
| 4 | sort rows desc | [4.00, 4.00, 3.85, 3.65, 3.65, 3.50] |
| 5 | lookup 4.00 → 1, 4.00 → 1 | [(4.00,1) (4.00,1)] |
| 6 | lookup 3.85 → 2 | … (3.85,2) |
| 7 | lookup 3.65 → 3, 3.65 → 3 | … (3.65,3) (3.65,3) |
| 8 | lookup 3.50 → 4 | … (3.50,4) |

Result: `[(4.00,1) (4.00,1) (3.85,2) (3.65,3) (3.65,3) (3.50,4)]` ✔

---

## Approach 3 — Sort and Sweep (Optimal)

### Intuition

After sorting descending, the ranks are completely forced — no lookup structure needed. Walk from the top: the first row is rank 1; each time the score **changes** the rank ticks up by exactly 1; equal neighbours inherit the current rank. Ties share (rule 2) because the counter doesn't move on equality; no holes exist (rule 3) because the counter only ever moves by +1. This is precisely how the `DENSE_RANK()` window function evaluates a sorted partition, using just two scalars of state.

### Algorithm

1. Copy the rows and sort by score descending.
2. Initialise `rank = 0`, `prev = sentinel`.
3. For each row top-down:
   1. If its score (in cents) differs from `prev`: `rank++`, `prev = cents`.
   2. Emit `(score, rank)`.

### Complexity

- **Time:** O(n log n) — the mandatory output sort; the sweep itself is O(n).
- **Space:** O(1) — beyond the sorted copy and the output slice, only `rank` and `prev`.

### Code

```go
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
```

### Dry Run

Example 1, rows sorted descending: `4.00, 4.00, 3.85, 3.65, 3.65, 3.50`

| Step | row score (cents) | cents != prev? | rank after | prev after | emit |
|------|-------------------|----------------|------------|------------|------|
| 0 | — | — | 0 | −∞ | — |
| 1 | 4.00 (400) | yes | 1 | 400 | (4.00, 1) |
| 2 | 4.00 (400) | no — tie | 1 | 400 | (4.00, 1) |
| 3 | 3.85 (385) | yes | 2 | 385 | (3.85, 2) |
| 4 | 3.65 (365) | yes | 3 | 365 | (3.65, 3) |
| 5 | 3.65 (365) | no — tie | 3 | 365 | (3.65, 3) |
| 6 | 3.50 (350) | yes | 4 | 350 | (3.50, 4) |

Result: `[(4.00,1) (4.00,1) (3.85,2) (3.65,3) (3.65,3) (3.50,4)]` ✔

---

## Key Takeaways

- **Dense rank = 1 + count of distinct greater values.** Three equivalent computations: count per row (O(n²)), index in sorted distinct values (coordinate compression), or a sorted sweep incrementing on change (window-function style). The sweep is the cheapest.
- **RANK() vs DENSE_RANK()** — the only difference is what the counter does after a tie: `RANK()` jumps by the tie-group size (holes: 1,1,3), `DENSE_RANK()` moves by +1 (no holes: 1,1,2). In the sweep, that's `rank++` on change vs `rank = i+1` on change.
- **Never compare money/two-decimal floats raw** — normalise to integer cents (`round(x*100)`) before `==`/`>`; float64 cannot represent 3.65 exactly.
- **Rank belongs to the value, not the row** — dedupe first whenever ranks must be shared by ties; this is the same insight as #176/#177's "distinct" requirement.
- Canonical SQL:
  ```sql
  SELECT score, DENSE_RANK() OVER (ORDER BY score DESC) AS `rank`
  FROM Scores;
  -- Pre-window-function form (MySQL 5.7 era):
  SELECT s.score,
         (SELECT COUNT(DISTINCT s2.score) FROM Scores s2
          WHERE s2.score > s.score) + 1 AS `rank`
  FROM Scores s ORDER BY s.score DESC;
  ```
  (`rank` needs backticks — it became a reserved word in MySQL 8.)

---

## Related Problems

- LeetCode #1331 — Rank Transform of an Array (dense rank on an array; identical coordinate compression)
- LeetCode #176 — Second Highest Salary (distinct-value ranking, single answer)
- LeetCode #177 — Nth Highest Salary (distinct-value ranking, parameterised)
- LeetCode #185 — Department Top Three Salaries (dense rank *per group*, `PARTITION BY`)
- LeetCode #506 — Relative Ranks (rank assignment via sorting, no ties)
