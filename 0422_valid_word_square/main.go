package main

import (
	"fmt"
	"strings"
)

// ── Approach 1: Build Columns and Compare ────────────────────────────────────
//
// buildColumns solves Valid Word Square by explicitly constructing each column
// string and comparing it to the corresponding row.
//
// Intuition:
//
//	A word square is valid when, for every k, "row k" equals "column k". The
//	most literal reading is: build column k by reading the k-th character of
//	every row that is long enough, then string-compare it with row k. The only
//	real trap is that rows can have different lengths (a "ragged" grid), so a
//	row only contributes to column k when that row has a k-th character.
//
// Algorithm:
//  1. Let n = number of rows.
//  2. For each k in 0..n-1:
//     a. Build column k: for each row r, if k < len(words[r]) append words[r][k].
//     b. If that column string != words[k], return false.
//  3. If all columns matched their rows, return true.
//
// Time:  O(n·L) where L is the max word length — every character is visited
//
//	while building columns (n columns, each up to n chars).
//
// Space: O(L) — one column string built at a time (O(n·L) total transiently).
func buildColumns(words []string) bool {
	n := len(words) // number of rows; a square has at most n columns
	for k := 0; k < n; k++ {
		var col strings.Builder // column k, assembled character by character
		for r := 0; r < n; r++ {
			// Row r contributes to column k only if it is long enough to have
			// a k-th character; a ragged (short) row simply stops early.
			if k < len(words[r]) {
				col.WriteByte(words[r][k]) // the k-th char of row r sits in column k, row r
			} else {
				break // rows are read top-to-bottom; once one is too short, a full
				// square would require the column to end here — but we compare the
				// assembled prefix directly to row k below, which catches mismatches.
			}
		}
		if col.String() != words[k] { // column k must read exactly like row k
			return false
		}
	}
	return true
}

// ── Approach 2: Symmetric Index Check (Optimal) ──────────────────────────────
//
// symmetricCheck solves Valid Word Square without building any strings, by
// checking the single symmetry condition words[i][j] == words[j][i] for every
// filled cell, guarded by careful bounds tests.
//
// Intuition:
//
//	"Row k equals column k for all k" is exactly the statement that the grid is
//	symmetric across its main diagonal: cell (i, j) must equal cell (j, i).
//	Because rows are ragged, we must be strict about existence: if (i, j) is a
//	real character, then (j, i) must also be a real character AND hold the same
//	value. In particular, if row i has a j-th character, there must even *be* a
//	row j — otherwise the transpose cell is missing and the square is invalid.
//
// Algorithm:
//  1. For each row i and each column j within row i:
//     a. If j >= number of rows → row j does not exist → return false
//     (column i would be longer than row i allows).
//     b. If i >= len(words[j]) → the mirror cell (j, i) is missing → false.
//     c. If words[i][j] != words[j][i] → asymmetric → false.
//  2. If no violation is found, return true.
//
// Time:  O(n·L) — each existing character is checked once against its mirror.
// Space: O(1) — pure index arithmetic, no auxiliary strings.
func symmetricCheck(words []string) bool {
	n := len(words) // number of rows == max possible number of columns
	for i := 0; i < n; i++ {
		for j := 0; j < len(words[i]); j++ {
			// (i, j) exists. Its mirror is (j, i). For validity that mirror
			// must also exist and match.

			// There must be a j-th row at all; otherwise column i extends past
			// the number of rows while row i still has a character — invalid.
			if j >= n {
				return false
			}
			// Row j must have an i-th character, else the mirror cell is absent.
			if i >= len(words[j]) {
				return false
			}
			// The defining symmetry: cell (i,j) equals cell (j,i).
			if words[i][j] != words[j][i] {
				return false
			}
		}
	}
	return true
}

func main() {
	fmt.Println("=== Approach 1: Build Columns and Compare ===")
	fmt.Println(buildColumns([]string{"abcd", "bnrt", "crmy", "dtye"})) // expected true
	fmt.Println(buildColumns([]string{"abcd", "bnrt", "crm", "dt"}))    // expected true
	fmt.Println(buildColumns([]string{"ball", "area", "read", "lady"})) // expected false

	fmt.Println("=== Approach 2: Symmetric Index Check (Optimal) ===")
	fmt.Println(symmetricCheck([]string{"abcd", "bnrt", "crmy", "dtye"})) // expected true
	fmt.Println(symmetricCheck([]string{"abcd", "bnrt", "crm", "dt"}))    // expected true
	fmt.Println(symmetricCheck([]string{"ball", "area", "read", "lady"})) // expected false
}
