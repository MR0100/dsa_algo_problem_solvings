package main

import (
	"fmt"
	"strings"
)

// ── Approach 1: Stack of Path Lengths ────────────────────────────────────────
//
// stackLengths solves Longest Absolute File Path by splitting the input on
// '\n' and maintaining a stack whose entry at depth d is the total path length
// of the current directory chain up to level d.
//
// Intuition:
//
//	The number of leading '\t' tabs gives an entry's depth. The absolute path
//	length of an entry at depth d is (length of its own name) + (path length of
//	its parent at depth d-1) + 1 for the '/' separator. A stack indexed by
//	depth lets us look up the parent's accumulated length in O(1). When we hit
//	an entry containing a '.', it is a file, so we update the best answer.
//
// Algorithm:
//  1. Split input on '\n' into lines.
//  2. For each line: depth = count of leading '\t'; name = line without tabs.
//  3. curLen = (depth==0 ? len(name) : stack[depth-1] + 1 + len(name)).
//  4. Store stack[depth] = curLen. If name has a '.', it is a file → update max.
//  5. Return the max file path length (0 if no file).
//
// Time:  O(N) — N = total characters; each line scanned a constant number of times.
// Space: O(D) — D = maximum directory depth (stack size).
func stackLengths(input string) int {
	lines := strings.Split(input, "\n") // each token is one dir/file entry
	// stack[d] = cumulative path length of the current chain at depth d.
	stack := make([]int, len(lines)+1)
	longest := 0

	for _, line := range lines {
		// Count leading tabs → this is the entry's nesting depth.
		depth := 0
		for depth < len(line) && line[depth] == '\t' {
			depth++
		}
		name := line[depth:] // strip the leading tabs to get the raw name
		nameLen := len(name) // characters in this dir/file name

		curLen := nameLen // path length if this were a top-level entry
		if depth > 0 {
			// parent chain length + '/' + this name
			curLen = stack[depth-1] + 1 + nameLen
		}
		stack[depth] = curLen // record cumulative length at this depth

		// A '.' in the name marks a file (e.g. "file.ext"); measure it.
		if strings.Contains(name, ".") {
			if curLen > longest {
				longest = curLen
			}
		}
	}
	return longest
}

// ── Approach 2: Map depth → running length (Optimal, single pass) ────────────
//
// mapDepthLength solves the same problem using a map from depth to the running
// path length, avoiding pre-splitting the whole string into a slice.
//
// Intuition:
//
//	Identical accounting to Approach 1, but we treat depth→length as a map so
//	the code reads as "the parent length at depth-1 plus my name plus a slash."
//	This makes the O(1) parent-lookup explicit and works even if the input were
//	streamed line by line.
//
// Algorithm:
//  1. Split on '\n'. lengths[0] = 0 (a virtual root so top-level uses +0).
//  2. For each line: depth = tabs + 1 (so root children live at depth 1);
//     name = line without tabs.
//  3. If it is a file (has '.'): answer = max(answer, lengths[depth-1] + len(name)).
//     Else: lengths[depth] = lengths[depth-1] + len(name) + 1 (name plus '/').
//  4. Return answer.
//
// Time:  O(N) — one pass over all characters.
// Space: O(D) — one map entry per active depth.
func mapDepthLength(input string) int {
	lines := strings.Split(input, "\n")
	// lengths[d] = path length (including trailing '/') of the current dir at depth d.
	// lengths[0] = 0 acts as the virtual filesystem root.
	lengths := map[int]int{0: 0}
	answer := 0

	for _, line := range lines {
		depth := 0
		for depth < len(line) && line[depth] == '\t' {
			depth++
		}
		name := line[depth:]
		level := depth + 1 // shift so top-level entries sit at level 1

		if strings.Contains(name, ".") {
			// File: full path = parent dir path + file name (no trailing slash).
			pathLen := lengths[level-1] + len(name)
			if pathLen > answer {
				answer = pathLen
			}
		} else {
			// Directory: store its path length WITH a trailing '/' for children.
			lengths[level] = lengths[level-1] + len(name) + 1
		}
	}
	return answer
}

func main() {
	// Example 1: "dir\n\tsubdir1\n\tsubdir2\n\t\tfile.ext" → "dir/subdir2/file.ext" = 20
	ex1 := "dir\n\tsubdir1\n\tsubdir2\n\t\tfile.ext"
	// Example 2: longest is "dir/subdir2/subsubdir2/file2.ext" = 32
	ex2 := "dir\n\tsubdir1\n\t\tfile1.ext\n\t\tsubsubdir1\n\tsubdir2\n\t\tsubsubdir2\n\t\t\tfile2.ext"
	// Example 3: "a" → no file → 0
	ex3 := "a"

	fmt.Println("=== Approach 1: Stack of Path Lengths ===")
	fmt.Println(stackLengths(ex1)) // expected 20
	fmt.Println(stackLengths(ex2)) // expected 32
	fmt.Println(stackLengths(ex3)) // expected 0

	fmt.Println("=== Approach 2: Map depth → running length (Optimal) ===")
	fmt.Println(mapDepthLength(ex1)) // expected 20
	fmt.Println(mapDepthLength(ex2)) // expected 32
	fmt.Println(mapDepthLength(ex3)) // expected 0
}
