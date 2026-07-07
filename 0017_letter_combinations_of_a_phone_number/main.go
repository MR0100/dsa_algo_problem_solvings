package main

import "fmt"

var phoneMap = map[byte]string{
	'2': "abc", '3': "def", '4': "ghi", '5': "jkl",
	'6': "mno", '7': "pqrs", '8': "tuv", '9': "wxyz",
}

// ── Approach 1: Backtracking (Recursive) ─────────────────────────────────────
//
// backtracking builds each combination character by character, branching
// over all letters for the current digit, and collecting complete strings.
//
// Intuition:
//   At each position i in digits, we have len(phoneMap[digits[i]]) choices.
//   Recurse with i+1 after appending each letter. When i == len(digits),
//   the current path is a complete combination — add it to results.
//
// Time:  O(4^n × n) — 4^n combinations (worst case: all '7' or '9' digits
//                      with 4 letters each), each of length n to copy.
// Space: O(n)       — recursion stack depth + current path buffer.
func backtracking(digits string) []string {
	if len(digits) == 0 {
		return []string{}
	}
	var result []string
	btHelper(digits, 0, []byte{}, &result)
	return result
}

func btHelper(digits string, idx int, path []byte, result *[]string) {
	if idx == len(digits) {
		*result = append(*result, string(path))
		return
	}
	for _, ch := range phoneMap[digits[idx]] {
		path = append(path, byte(ch))
		btHelper(digits, idx+1, path, result)
		path = path[:len(path)-1] // backtrack
	}
}

// ── Approach 2: Iterative BFS / Queue ────────────────────────────────────────
//
// iterativeBFS starts with a queue containing an empty string and expands
// each partial combination by all letters of the next digit.
//
// Intuition:
//   Treat partial combinations as states in a BFS. For each digit, expand
//   every existing partial by all its letters. The queue holds all partial
//   combinations of the same length at each step.
//
// Time:  O(4^n × n) — same as backtracking.
// Space: O(4^n × n) — the queue holds all partial combinations.
func iterativeBFS(digits string) []string {
	if len(digits) == 0 {
		return []string{}
	}
	queue := []string{""}

	for i := 0; i < len(digits); i++ {
		letters := phoneMap[digits[i]]
		next := make([]string, 0, len(queue)*len(letters))
		for _, partial := range queue {
			for _, ch := range letters {
				next = append(next, partial+string(ch))
			}
		}
		queue = next
	}
	return queue
}

func main() {
	examples := []struct {
		digits string
		expect []string
	}{
		{"23", []string{"ad", "ae", "af", "bd", "be", "bf", "cd", "ce", "cf"}},
		{"", []string{}},
		{"2", []string{"a", "b", "c"}},
	}

	approaches := []struct {
		name string
		fn   func(string) []string
	}{
		{"Approach 1: Backtracking (recursive) O(4^n·n) T | O(n)     S", backtracking},
		{"Approach 2: Iterative BFS          ✅ O(4^n·n) T | O(4^n·n) S", iterativeBFS},
	}

	for _, ex := range examples {
		fmt.Printf("digits=%q  expect=%v\n", ex.digits, ex.expect)
		for _, ap := range approaches {
			fmt.Printf("  %-60s → %v\n", ap.name, ap.fn(ex.digits))
		}
		fmt.Println()
	}
}
