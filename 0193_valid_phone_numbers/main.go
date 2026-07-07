package main

import (
	"fmt"
	"regexp"
	"strings"
)

// LeetCode #193 is one of the four Shell problems ("write a one-line bash
// script..."). Per this repo's Go-only rule, each approach re-implements the
// same filter in Go: file.txt is simulated as an in-memory string and every
// function returns only the lines that are valid phone numbers, i.e. exactly
// (xxx) xxx-xxxx or xxx-xxx-xxxx where x is a digit.

// fileTxt mirrors the official example content of file.txt.
const fileTxt = `987-123-4567
123 456 7890
(123) 456-7890`

// isDigit reports whether b is an ASCII digit '0'..'9'.
func isDigit(b byte) bool {
	return b >= '0' && b <= '9'
}

// matchesDashed checks the fixed template xxx-xxx-xxxx (length 12,
// dashes at indices 3 and 7, digits everywhere else).
func matchesDashed(s string) bool {
	if len(s) != 12 { // template is exactly 12 characters
		return false
	}
	for i := 0; i < 12; i++ {
		if i == 3 || i == 7 { // separator slots
			if s[i] != '-' {
				return false
			}
		} else if !isDigit(s[i]) { // every other slot must be a digit
			return false
		}
	}
	return true
}

// matchesParen checks the fixed template (xxx) xxx-xxxx (length 14, literal
// '(' ')' ' ' '-' at indices 0, 4, 5, 9, digits everywhere else).
func matchesParen(s string) bool {
	if len(s) != 14 { // template is exactly 14 characters
		return false
	}
	if s[0] != '(' || s[4] != ')' || s[5] != ' ' || s[9] != '-' {
		return false // a literal separator is out of place
	}
	for _, i := range []int{1, 2, 3, 6, 7, 8, 10, 11, 12, 13} { // digit slots
		if !isDigit(s[i]) {
			return false
		}
	}
	return true
}

// ── Approach 1: Brute Force (Position-by-Position Template Check) ────────────
//
// bruteForce solves Valid Phone Numbers by comparing every line, character by
// character, against the two allowed fixed-width templates.
//
// Intuition:
//
//	Both valid formats are rigid: every character position is either "must be
//	a digit" or "must be this exact separator". So a line is valid iff its
//	length matches one template and every index passes that template's check.
//	This is what a regex engine would do, written out by hand.
//
// Algorithm:
//  1. Split the file into lines.
//  2. For each line, run matchesDashed (xxx-xxx-xxxx) and, failing that,
//     matchesParen ((xxx) xxx-xxxx).
//  3. Keep the line if either template accepts it.
//
// Time:  O(N) — N total characters; each line is scanned a constant number of times.
// Space: O(1) extra beyond the output slice — only index variables.
func bruteForce(fileContent string) []string {
	valid := []string{}
	for _, line := range strings.Split(fileContent, "\n") { // one candidate per line
		if matchesDashed(line) || matchesParen(line) { // accept if either template fits
			valid = append(valid, line)
		}
	}
	return valid
}

// allDigits reports whether s is non-empty and consists solely of digits.
func allDigits(s string) bool {
	if len(s) == 0 {
		return false // an empty group can never be a digit block
	}
	for i := 0; i < len(s); i++ {
		if !isDigit(s[i]) {
			return false
		}
	}
	return true
}

// ── Approach 2: Split and Validate Tokens ────────────────────────────────────
//
// splitAndValidate solves Valid Phone Numbers by decomposing each line into
// its separator-delimited digit groups and validating group lengths.
//
// Intuition:
//
//	Instead of walking indices, think in tokens: a dashed number is exactly
//	three '-'-separated groups of sizes 3/3/4; a parenthesised number is
//	"(ddd)" + space + "ddd-dddd". Splitting expresses the format rules at a
//	higher level and localises each rule to one small check.
//
// Algorithm:
//  1. If the line starts with '(', split once on the space: the head must be
//     "(ddd)" and the tail must split on '-' into digit groups of 3 and 4.
//  2. Otherwise split on '-': accept exactly three all-digit groups with
//     lengths 3, 3, 4.
//  3. Collect the lines that pass.
//
// Time:  O(N) — each line is split and scanned a constant number of times.
// Space: O(L) per line for the transient token slices (L = line length).
func splitAndValidate(fileContent string) []string {
	valid := []string{}
	for _, line := range strings.Split(fileContent, "\n") {
		if validBySplit(line) {
			valid = append(valid, line)
		}
	}
	return valid
}

// validBySplit applies the token rules for the two formats to one line.
func validBySplit(line string) bool {
	if strings.HasPrefix(line, "(") { // candidate for the (xxx) xxx-xxxx form
		parts := strings.SplitN(line, " ", 2) // split into "(xxx)" and "xxx-xxxx"
		if len(parts) != 2 {
			return false // no space → cannot match the parenthesised form
		}
		head, rest := parts[0], parts[1]
		// head must be exactly "(" + 3 digits + ")".
		if len(head) != 5 || head[0] != '(' || head[4] != ')' || !allDigits(head[1:4]) {
			return false
		}
		// rest must be exactly 3 digits, '-', 4 digits.
		groups := strings.Split(rest, "-")
		return len(groups) == 2 &&
			len(groups[0]) == 3 && allDigits(groups[0]) &&
			len(groups[1]) == 4 && allDigits(groups[1])
	}

	// Candidate for the xxx-xxx-xxxx form: exactly three digit groups 3/3/4.
	groups := strings.Split(line, "-")
	return len(groups) == 3 &&
		len(groups[0]) == 3 && allDigits(groups[0]) &&
		len(groups[1]) == 3 && allDigits(groups[1]) &&
		len(groups[2]) == 4 && allDigits(groups[2])
}

// phoneRe encodes both formats in one anchored pattern:
//
//	^        start of line (nothing before the number)
//	\(\d{3}\) ␣  →  "(xxx) "  … or …  \d{3}-  →  "xxx-"
//	\d{3}-\d{4}  →  the shared "xxx-xxxx" tail
//	$        end of line (nothing after the number)
//
// Compiled once at package level — compiling inside the loop would be wasteful.
var phoneRe = regexp.MustCompile(`^(\(\d{3}\) |\d{3}-)\d{3}-\d{4}$`)

// ── Approach 3: Regular Expression (Optimal) ─────────────────────────────────
//
// regexMatch solves Valid Phone Numbers with a single anchored regular
// expression — the direct translation of the intended one-line bash answer
// `grep -P '^(\(\d{3}\) |\d{3}-)\d{3}-\d{4}$' file.txt`.
//
// Intuition:
//
//	Both formats share the tail "xxx-xxxx" and differ only in the prefix:
//	"(xxx) " versus "xxx-". One alternation captures exactly that, and the
//	^...$ anchors reject any extra leading/trailing characters — the classic
//	bug in naive solutions (an unanchored pattern would accept
//	"0(001) 345-0000").
//
// Algorithm:
//  1. Compile ^(\(\d{3}\) |\d{3}-)\d{3}-\d{4}$ once.
//  2. Keep every line the pattern matches.
//
// Time:  O(N) — Go's RE2 engine guarantees linear-time matching, no backtracking.
// Space: O(1) per line — the compiled automaton is fixed-size and shared.
func regexMatch(fileContent string) []string {
	valid := []string{}
	for _, line := range strings.Split(fileContent, "\n") {
		if phoneRe.MatchString(line) { // anchored full-line match
			valid = append(valid, line)
		}
	}
	return valid
}

// printLines prints each surviving line, matching the script output.
func printLines(lines []string) {
	for _, l := range lines {
		fmt.Println(l)
	}
}

func main() {
	fmt.Println("=== Approach 1: Brute Force (Position-by-Position Template Check) ===")
	printLines(bruteForce(fileTxt))
	// expected:
	// 987-123-4567
	// (123) 456-7890

	fmt.Println("=== Approach 2: Split and Validate Tokens ===")
	printLines(splitAndValidate(fileTxt))
	// expected:
	// 987-123-4567
	// (123) 456-7890

	fmt.Println("=== Approach 3: Regular Expression (Optimal) ===")
	printLines(regexMatch(fileTxt))
	// expected:
	// 987-123-4567
	// (123) 456-7890
}
