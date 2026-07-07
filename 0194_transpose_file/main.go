package main

import (
	"fmt"
	"strings"
)

// LeetCode #194 is one of the four Shell problems ("Given a text file
// file.txt, transpose its content"). Per this repo's Go-only rule, each
// approach re-implements the transpose in Go: file.txt is simulated as an
// in-memory string, every row has the same number of space-separated fields,
// and each function returns the transposed lines (column j of the input
// becomes line j of the output).

// fileTxt mirrors the official example content of file.txt.
const fileTxt = `name age
alice 21
ryan 30`

// ── Approach 1: Brute Force (Re-Scan the File per Column) ────────────────────
//
// bruteForce solves Transpose File by making one full pass over the raw file
// for every output line: pass j re-splits every row and plucks its j-th field.
//
// Intuition:
//
//	Output line j is "the j-th field of every row, in row order". The most
//	naive way to get it is to re-read the whole file once per column — no
//	stored matrix, just repeated extraction. This mirrors the naive shell
//	loop `for j in 1..C: cut -d' ' -f$j file.txt | paste -s`.
//
// Algorithm:
//  1. Split the file into rows; the first row's field count fixes C.
//  2. For each column j in 0..C-1:
//     a. Walk every row, re-split it into fields, take field j.
//     b. Join the collected fields with single spaces → output line j.
//
// Time:  O(C·N) — the entire input (N characters, C columns per row) is
//
//	re-tokenised for each of the C output lines.
//
// Space: O(R + L) per pass — one column's fields (R rows) and a row's tokens.
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
			fields := strings.Fields(row)    // wasteful: re-split this row on every pass
			parts = append(parts, fields[j]) // pluck only the j-th field this pass
		}
		out = append(out, strings.Join(parts, " ")) // column j → output line j
	}
	return out
}

// ── Approach 2: Matrix Transpose (Parse Once, Swap Indices) ──────────────────
//
// matrixTranspose solves Transpose File the classic way: parse the whole file
// into a 2-D slice once, then read it out column-major.
//
// Intuition:
//
//	Transposition is an index swap: output[j][i] = input[i][j]. If we
//	materialise the file as a matrix, producing the answer is just iterating
//	with the loops exchanged — each cell is parsed exactly once (unlike
//	Approach 1, which re-parses every row C times).
//
// Algorithm:
//  1. Split the file into rows; split each row into fields → matrix[i][j].
//  2. For each column j, collect matrix[0][j], matrix[1][j], …, matrix[R-1][j].
//  3. Join each collection with spaces → output line j.
//
// Time:  O(R·C) — every cell is parsed once and emitted once.
// Space: O(R·C) — the full matrix lives in memory (fine unless the file is huge).
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

// ── Approach 3: Streaming Column Builders (Optimal) ──────────────────────────
//
// streamBuilders solves Transpose File in a single streaming pass: one
// strings.Builder per column accumulates its output line as rows fly by, so
// the raw file never needs to be retained or revisited.
//
// Intuition:
//
//	You do not need the whole matrix at once — while reading row i, each field
//	already knows its destination: field j belongs at the end of output line j.
//	Appending it immediately means memory holds only the OUTPUT under
//	construction, exactly how the canonical awk solution builds s[j]=s[j]" "$j
//	while scanning stdin once.
//
// Algorithm:
//  1. Walk the rows once, splitting each into fields.
//  2. On the first row, create one builder per column (as pointers — a
//     strings.Builder must not be copied once written to).
//  3. Append field j to builder j, prefixing a space unless it is the first
//     entry of that line.
//  4. After the pass, materialise each builder → output line j.
//
// Time:  O(R·C) — each cell is touched exactly once; builder appends are
//
//	amortised O(1).
//
// Space: O(output) — only the transposed lines themselves, no input matrix.
func streamBuilders(fileContent string) []string {
	if fileContent == "" {
		return nil
	}

	var builders []*strings.Builder                        // builders[j] accumulates output line j
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

// printLines prints each transposed line, matching the script output.
func printLines(lines []string) {
	for _, l := range lines {
		fmt.Println(l)
	}
}

func main() {
	fmt.Println("=== Approach 1: Brute Force (Re-Scan the File per Column) ===")
	printLines(bruteForce(fileTxt))
	// expected:
	// name alice ryan
	// age 21 30

	fmt.Println("=== Approach 2: Matrix Transpose (Parse Once, Swap Indices) ===")
	printLines(matrixTranspose(fileTxt))
	// expected:
	// name alice ryan
	// age 21 30

	fmt.Println("=== Approach 3: Streaming Column Builders (Optimal) ===")
	printLines(streamBuilders(fileTxt))
	// expected:
	// name alice ryan
	// age 21 30
}
