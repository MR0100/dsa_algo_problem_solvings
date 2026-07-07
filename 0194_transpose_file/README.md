# 0194 — Transpose File

> LeetCode #194 · Difficulty: Medium
> **Categories:** Shell, Matrix, String

---

## Problem Statement

Given a text file `file.txt`, transpose its content.

You may assume that each row has the same number of columns, and each field is separated by the `' '` character.

**Example:**

If `file.txt` has the following content:

```
name age
alice 21
ryan 30
```

Output the following:

```
name alice ryan
age 21 30
```

> **Repo note:** #194 is one of LeetCode's four Shell problems (the official
> ask is an awk/bash script). Following this repo's Go-only convention, every
> approach below re-implements the identical transpose in Go, treating the
> file's content as an in-memory string (`file.txt` → the `fileTxt` constant).
> The canonical awk answer appears in Key Takeaways for interview completeness.

---

## Company Frequency

| Company    | Frequency        | Last Reported |
|------------|------------------|---------------|
| Amazon     | ★★☆☆☆ Low        | 2023          |
| Google     | ★★☆☆☆ Low        | 2023          |
| Microsoft  | ★★☆☆☆ Low        | 2022          |
| Bloomberg  | ★☆☆☆☆ Rare       | 2022          |
| Yandex     | ★☆☆☆☆ Rare       | 2022          |

> ⚠️ Frequency data is crowd-sourced from LeetCode Discuss, Glassdoor, and
> community interview reports. Treat as a signal, not a guarantee.

---

## DSA Concepts Used

- **Matrix Traversal** — transposition is the index swap `output[j][i] = input[i][j]`; iterating a rectangular grid column-major is the heart of the problem → see [`/dsa/matrix_traversal.md`](/dsa/matrix_traversal.md)
- **String Algorithms (tokenisation + builders)** — rows must be split into fields and output lines assembled efficiently; `strings.Builder` avoids quadratic concatenation → see [`/dsa/string_algorithms.md`](/dsa/string_algorithms.md)

---

## Approaches Overview

| # | Approach | Time | Space | When to use |
|---|----------|------|-------|-------------|
| 1 | Brute Force (re-scan the file per column) | O(C·N) | O(R + L) | Baseline; no stored matrix, but re-parses everything C times |
| 2 | Matrix Transpose (parse once, swap indices) | O(R·C) | O(R·C) | The standard answer; simplest correct linear solution |
| 3 | Streaming Column Builders (Optimal) | O(R·C) | O(output) | Huge files read as a stream; holds only the output, never the input |

*R = rows, C = columns, N = total characters in the file, L = length of one row.*

---

## Approach 1 — Brute Force (Re-Scan the File per Column)

### Intuition
Output line `j` is "the j-th field of every row, in row order". The most naive way to obtain it is to re-read the whole file once per column — no stored matrix, just repeated extraction. It mirrors the naive shell loop `for j in 1..C: cut -d' ' -f$j file.txt | paste -s -d' '`, and its wastefulness (every row re-tokenised C times) is exactly what the better approaches remove.

### Algorithm
1. Split the file into rows; the first row's field count fixes C (all rows are guaranteed equally wide).
2. For each column `j` in `0..C-1`:
   a. Walk every row, re-split it into fields, take field `j`.
   b. Join the collected fields with single spaces → output line `j`.

### Complexity
- **Time:** O(C·N) — the entire input of N characters is re-tokenised once for each of the C output lines.
- **Space:** O(R + L) per pass — one column's R fields plus one row's transient tokens (no full matrix).

### Code
```go
func bruteForce(fileContent string) []string {
	rows := strings.Split(fileContent, "\n") // file → rows
	if len(rows) == 0 || fileContent == "" {
		return nil // empty file → nothing to transpose
	}
	cols := len(strings.Fields(rows[0])) // every row is guaranteed this wide

	out := make([]string, 0, cols)
	for j := 0; j < cols; j++ { // one full file pass per output line
		parts := make([]string, 0, len(rows))
		for _, row := range rows {
			fields := strings.Fields(row) // wasteful: re-split this row on every pass
			parts = append(parts, fields[j]) // pluck only the j-th field this pass
		}
		out = append(out, strings.Join(parts, " ")) // column j → output line j
	}
	return out
}
```

### Dry Run
Example: rows = `[name age | alice 21 | ryan 30]`, so R = 3, C = 2.

| pass j | row scanned | fields (re-split) | field j taken | parts after | emitted line |
|--------|-------------|-------------------|---------------|-------------|--------------|
| 0 | `name age` | `[name age]` | `name` | `[name]` | |
| 0 | `alice 21` | `[alice 21]` | `alice` | `[name alice]` | |
| 0 | `ryan 30` | `[ryan 30]` | `ryan` | `[name alice ryan]` | `name alice ryan` |
| 1 | `name age` | `[name age]` (split **again**) | `age` | `[age]` | |
| 1 | `alice 21` | `[alice 21]` (again) | `21` | `[age 21]` | |
| 1 | `ryan 30` | `[ryan 30]` (again) | `30` | `[age 21 30]` | `age 21 30` |

Output: `name alice ryan`, `age 21 30`. ✓ (Note each row was tokenised C = 2 times.)

---

## Approach 2 — Matrix Transpose (Parse Once, Swap Indices)

### Intuition
Transposition is nothing but an index swap: `output[j][i] = input[i][j]`. If we materialise the file as a 2-D slice, producing the answer is just iterating with the loops exchanged — each cell is parsed exactly once, fixing Approach 1's redundant re-splitting at the cost of holding the whole matrix in memory.

### Algorithm
1. Split the file into rows; split each row into fields → `matrix[i][j]`.
2. For each column `j` in `0..C-1`, collect `matrix[0][j], matrix[1][j], …, matrix[R-1][j]`.
3. Join each collection with single spaces → output line `j`.

### Complexity
- **Time:** O(R·C) — every cell is parsed once and emitted once; joins are linear in output size.
- **Space:** O(R·C) — the full matrix lives in memory (fine unless the file is enormous).

### Code
```go
func matrixTranspose(fileContent string) []string {
	rows := strings.Split(fileContent, "\n")
	if len(rows) == 0 || fileContent == "" {
		return nil
	}

	// Parse once: matrix[i] = fields of row i.
	matrix := make([][]string, len(rows))
	for i, row := range rows {
		matrix[i] = strings.Fields(row)
	}

	rCount, cCount := len(matrix), len(matrix[0]) // R rows, C columns (rectangular by constraint)
	out := make([]string, cCount)
	for j := 0; j < cCount; j++ { // walk column-major: output[j][i] = matrix[i][j]
		parts := make([]string, rCount)
		for i := 0; i < rCount; i++ {
			parts[i] = matrix[i][j] // the index swap that IS the transpose
		}
		out[j] = strings.Join(parts, " ")
	}
	return out
}
```

### Dry Run
Example 1. Parse phase produces:

```
matrix = [ [name  age]
           [alice 21 ]
           [ryan  30 ] ]        R = 3, C = 2
```

Emit phase:

| j | i | matrix[i][j] | parts state | out[j] |
|---|---|--------------|-------------|--------|
| 0 | 0 | `name` | `[name _ _]` | |
| 0 | 1 | `alice` | `[name alice _]` | |
| 0 | 2 | `ryan` | `[name alice ryan]` | `name alice ryan` |
| 1 | 0 | `age` | `[age _ _]` | |
| 1 | 1 | `21` | `[age 21 _]` | |
| 1 | 2 | `30` | `[age 21 30]` | `age 21 30` |

Output: `name alice ryan`, `age 21 30`. ✓

---

## Approach 3 — Streaming Column Builders (Optimal)

### Intuition
You do not need the whole input matrix at once — while reading row `i`, each field already knows its destination: field `j` belongs at the end of output line `j`. Appending it immediately to a per-column `strings.Builder` means memory holds only the OUTPUT under construction, and the input is consumed in a single streaming pass. This is exactly how the canonical awk one-liner builds `s[j] = s[j] " " $j` while scanning stdin once — the right shape when the file is too big to keep around.

### Algorithm
1. Walk the rows once, splitting each into fields.
2. On the first row, create one builder per column — as **pointers**, because a `strings.Builder` must never be copied after its first write (slice growth copies elements, and a copied Builder panics on use).
3. Append field `j` to builder `j`, writing a `' '` separator first unless the builder is still empty.
4. After the single pass, materialise each builder into output line `j`.

### Complexity
- **Time:** O(R·C) — each cell is touched exactly once; `Builder` appends are amortised O(1) (no quadratic string concatenation).
- **Space:** O(output) — only the transposed lines being built; the raw matrix is never stored.

### Code
```go
func streamBuilders(fileContent string) []string {
	if fileContent == "" {
		return nil
	}

	var builders []*strings.Builder // builders[j] accumulates output line j
	for _, row := range strings.Split(fileContent, "\n") { // single pass over rows
		for j, field := range strings.Fields(row) {
			if j == len(builders) { // first row discovers each new column
				builders = append(builders, &strings.Builder{}) // pointer: Builders must not be copied after use
			}
			if builders[j].Len() > 0 {
				builders[j].WriteByte(' ') // separator before every field except the first
			}
			builders[j].WriteString(field) // append this cell to its column's line
		}
	}

	out := make([]string, len(builders))
	for j, b := range builders {
		out[j] = b.String() // materialise each accumulated line
	}
	return out
}
```

### Dry Run
Example 1, single pass:

| row | field (j) | action | builder 0 state | builder 1 state |
|-----|-----------|--------|-----------------|-----------------|
| `name age` | `name` (0) | new builder 0; write `name` | `name` | — |
| `name age` | `age` (1) | new builder 1; write `age` | `name` | `age` |
| `alice 21` | `alice` (0) | space + `alice` | `name alice` | `age` |
| `alice 21` | `21` (1) | space + `21` | `name alice` | `age 21` |
| `ryan 30` | `ryan` (0) | space + `ryan` | `name alice ryan` | `age 21` |
| `ryan 30` | `30` (1) | space + `30` | `name alice ryan` | `age 21 30` |

Materialise: `["name alice ryan", "age 21 30"]`. ✓

---

## Key Takeaways

- **Transpose = swap the loop indices**: `out[j][i] = in[i][j]` — the same skeleton as LC 867 *Transpose Matrix*, just with string fields instead of ints.
- **Stream when you only need the output**: per-destination accumulators (builders, buckets) let you consume huge inputs in one pass without materialising them — a pattern that generalises to log processing and MapReduce-style shuffles.
- **`strings.Builder` must not be copied after use** — store pointers (`[]*strings.Builder`) when builders live in a growable slice, or the slice's reallocation copy will panic at the next write.
- **Join, don't concatenate in a loop** — `strings.Join`/`Builder` are O(n); `s += field` in a loop is O(n²).
- The canonical awk one-liner, worth knowing verbatim:
  ```bash
  awk '{ for (i = 1; i <= NF; i++) { s[i] = (NR == 1) ? $i : s[i] " " $i } }
       END { for (i = 1; i <= NF; i++) print s[i] }' file.txt
  ```
  `NF` = fields per row, `NR` = current row number; `s[i]` accumulates output line `i` exactly like Approach 3's builders.

---

## Related Problems

- LeetCode #867 — Transpose Matrix (the same index swap on an int matrix)
- LeetCode #48 — Rotate Image (transpose + reverse = 90° rotation)
- LeetCode #54 — Spiral Matrix (non-trivial matrix traversal order)
- LeetCode #192 — Word Frequency (Shell problem set)
- LeetCode #193 — Valid Phone Numbers (Shell problem set)
- LeetCode #195 — Tenth Line (Shell problem set)
