package main

import (
	"bufio"
	"fmt"
	"strings"
)

// LeetCode #195 is one of the four Shell problems ("Given a text file
// file.txt, print just the 10th line of the file"). Per this repo's Go-only
// rule, each approach re-implements the logic in Go: file.txt is simulated as
// an in-memory string and every function returns the 10th line, or "" when
// the file has fewer than 10 lines (the note's edge case — print nothing).

// fileTxt mirrors the official example content of file.txt (exactly 10 lines).
const fileTxt = `Line 1
Line 2
Line 3
Line 4
Line 5
Line 6
Line 7
Line 8
Line 9
Line 10`

// shortTxt exercises the note's edge case: fewer than 10 lines → no output.
const shortTxt = `Line 1
Line 2
Line 3`

// ── Approach 1: Brute Force (Materialise Every Line) ─────────────────────────
//
// bruteForce solves Tenth Line by splitting the entire file into a slice of
// lines and indexing the 10th one.
//
// Intuition:
//
//	The blunt instrument: load everything, then array-index. lines[9] is the
//	10th line (0-based indexing), and a length check answers the follow-up
//	"what if the file has fewer than 10 lines?" — return nothing. This is the
//	spirit of the bash `head -n 10 file.txt | tail -n +10` family.
//
// Algorithm:
//  1. Split the whole content on '\n' into a slice.
//  2. If the slice has fewer than 10 entries, return "" (print nothing).
//  3. Otherwise return lines[9].
//
// Time:  O(N) — the split walks all N characters, even those after line 10.
// Space: O(N) — every line is materialised, even though only one is needed.
func bruteForce(fileContent string) string {
	lines := strings.Split(fileContent, "\n") // materialise ALL lines up front
	if len(lines) < 10 {                      // fewer than 10 lines → nothing to print
		return ""
	}
	return lines[9] // 10th line lives at 0-based index 9
}

// ── Approach 2: Line Stream with Early Exit ──────────────────────────────────
//
// streamEarlyExit solves Tenth Line by scanning line-by-line and stopping the
// moment line 10 is produced — later lines are never even read.
//
// Intuition:
//
//	Only one line matters, so reading past it is pure waste. A streaming
//	scanner keeps just the current line in memory and can abandon the input
//	as soon as the target is reached — exactly how `sed -n '10p;10q'` beats
//	plain `sed -n '10p'` on a gigantic file.
//
// Algorithm:
//  1. Wrap the content in a bufio.Scanner (line-splitting mode).
//  2. Scan lines, incrementing a counter.
//  3. When the counter hits 10, return that line immediately.
//  4. If input ends first, return "".
//
// Time:  O(P) — P = characters in the first 10 lines only (early exit);
//
//	O(N) worst case when the file is short.
//
// Space: O(L) — the scanner buffers a single line of length L at a time.
func streamEarlyExit(fileContent string) string {
	scanner := bufio.NewScanner(strings.NewReader(fileContent)) // stream, don't slurp
	lineNo := 0
	for scanner.Scan() { // pulls ONE line per iteration
		lineNo++
		if lineNo == 10 {
			return scanner.Text() // found it — stop reading, skip the rest of the file
		}
	}
	return "" // ran out of lines before reaching 10
}

// ── Approach 3: Raw Byte Scan, Zero Allocation (Optimal) ─────────────────────
//
// byteScan solves Tenth Line by counting '\n' bytes directly: the 10th line is
// the substring between the 9th and 10th newline, so no line slice, scanner,
// or copy is ever allocated.
//
// Intuition:
//
//	Lines are just the gaps between newline bytes. The 10th line starts right
//	after the 9th '\n' and ends right before the 10th '\n' (or at EOF). Track
//	two things — how many newlines seen, where the 10th line starts — and
//	slice the original string once. Go string slicing shares the underlying
//	bytes, so this does zero allocations.
//
// Algorithm:
//  1. Walk the bytes with a newline counter and a start marker.
//  2. On the 9th '\n', record start = i+1 (the 10th line begins here).
//  3. On the 10th '\n', return content[start:i] and stop — early exit.
//  4. If EOF arrives with exactly 9 newlines seen, the 10th line is the
//     final unterminated line: return content[start:].
//  5. Otherwise (fewer than 9 newlines) return "".
//
// Time:  O(P) — bytes up to the 10th newline only; O(N) worst case.
// Space: O(1) — two integers; the returned slice shares the input's memory.
func byteScan(fileContent string) string {
	newlines := 0 // how many '\n' bytes seen so far
	start := 0    // index just past the 9th newline = start of the 10th line
	for i := 0; i < len(fileContent); i++ {
		if fileContent[i] != '\n' {
			continue // only newline bytes drive the state machine
		}
		newlines++
		if newlines == 9 {
			start = i + 1 // the 10th line begins immediately after the 9th '\n'
		}
		if newlines == 10 {
			return fileContent[start:i] // ends just before the 10th '\n'; early exit
		}
	}
	if newlines == 9 {
		return fileContent[start:] // 10th line is the last line, no trailing '\n'
	}
	return "" // fewer than 10 lines in the file
}

func main() {
	fmt.Println("=== Approach 1: Brute Force (Materialise Every Line) ===")
	fmt.Println(bruteForce(fileTxt)) // Line 10

	fmt.Println("=== Approach 2: Line Stream with Early Exit ===")
	fmt.Println(streamEarlyExit(fileTxt)) // Line 10

	fmt.Println("=== Approach 3: Raw Byte Scan, Zero Allocation (Optimal) ===")
	fmt.Println(byteScan(fileTxt)) // Line 10

	// Note edge case: fewer than 10 lines → every approach prints nothing ("").
	fmt.Println("=== Edge Case: 3-Line File (all approaches) ===")
	fmt.Printf("%q\n", bruteForce(shortTxt))      // ""
	fmt.Printf("%q\n", streamEarlyExit(shortTxt)) // ""
	fmt.Printf("%q\n", byteScan(shortTxt))        // ""
}
