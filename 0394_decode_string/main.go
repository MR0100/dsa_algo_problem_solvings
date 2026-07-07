package main

import (
	"fmt"
	"strings"
)

// ── Approach 1: Two Stacks (Optimal, Iterative) ──────────────────────────────
//
// twoStacks decodes strings of the form k[encoded] where the bracketed part is
// repeated k times, supporting arbitrary nesting.
//
// Intuition:
//
//	Nesting is naturally a stack problem. Keep the "string built so far" and the
//	"repeat count pending" for the current level. When we hit '[', we are
//	descending one level: push the current partial string and the current count,
//	then reset them for the inner level. When we hit ']', we finish the inner
//	level: pop the outer string and the count k, and append the inner string
//	repeated k times to the outer one.
//
// Algorithm:
//  1. countStack (ints), stringStack (strings); cur = "" and num = 0.
//  2. Scan each char:
//     - digit → num = num*10 + digit (multi-digit counts).
//     - '[' → push num and cur; reset num = 0, cur = "".
//     - ']' → pop k and prev; cur = prev + cur repeated k times.
//     - letter → cur += letter.
//  3. Return cur.
//
// Time:  O(total output length) — each output character produced once.
// Space: O(output length) for the stacks/builders.
func twoStacks(s string) string {
	countStack := []int{}     // pending repeat counts, one per open bracket
	stringStack := []string{} // partial strings accumulated before each '['
	cur := ""                 // string being built at the current nesting level
	num := 0                  // repeat count currently being parsed (may be multi-digit)

	for i := 0; i < len(s); i++ {
		ch := s[i]
		switch {
		case ch >= '0' && ch <= '9':
			num = num*10 + int(ch-'0') // accumulate multi-digit number
		case ch == '[':
			countStack = append(countStack, num)   // remember how many times inner repeats
			stringStack = append(stringStack, cur) // remember outer string so far
			num = 0                                // reset for the inner level
			cur = ""                               // inner level starts empty
		case ch == ']':
			// Pop the repeat count for this bracket.
			k := countStack[len(countStack)-1]
			countStack = countStack[:len(countStack)-1]
			// Pop the outer string that preceded this bracket.
			prev := stringStack[len(stringStack)-1]
			stringStack = stringStack[:len(stringStack)-1]
			// Outer + (inner repeated k times).
			cur = prev + strings.Repeat(cur, k)
		default:
			cur += string(ch) // ordinary letter joins the current level
		}
	}
	return cur
}

// ── Approach 2: Recursive Descent (Index Pointer) ────────────────────────────
//
// recursiveDescent decodes by recursion: it parses one "chunk" at a time and
// recurses whenever it opens a bracket, mirroring the grammar directly.
//
// Intuition:
//
//	The grammar is: sequence of (letters | number '[' sequence ']'). A recursive
//	function that decodes from a shared position and returns when it hits ']'
//	naturally handles nesting: on a digit it reads the count, expects '[',
//	recurses to decode the inside, then repeats that result k times.
//
// Algorithm:
//  1. Use a shared index pointer i into s.
//  2. decode() loops while i < len and s[i] != ']':
//     - letter → append to result, i++.
//     - digit → read full number; skip '['; res := decode(); skip ']';
//     append number-times res.
//  3. Return the built string; the top-level call decodes everything.
//
// Time:  O(total output length).
// Space: O(output length + nesting depth) for recursion + builders.
func recursiveDescent(s string) string {
	i := 0 // shared cursor into s, advanced by the closure below
	var decode func() string
	decode = func() string {
		var sb strings.Builder
		// Stop at end of string or at the ']' that closes this level.
		for i < len(s) && s[i] != ']' {
			ch := s[i]
			if ch >= '0' && ch <= '9' {
				// Parse the (possibly multi-digit) repeat count.
				k := 0
				for i < len(s) && s[i] >= '0' && s[i] <= '9' {
					k = k*10 + int(s[i]-'0')
					i++
				}
				i++               // consume '['
				inner := decode() // decode the bracketed content
				i++               // consume ']'
				// Repeat the inner decoded string k times.
				for r := 0; r < k; r++ {
					sb.WriteString(inner)
				}
			} else {
				sb.WriteByte(ch) // plain letter
				i++
			}
		}
		return sb.String()
	}
	return decode()
}

func main() {
	fmt.Println("=== Approach 1: Two Stacks ===")
	fmt.Printf("s=\"3[a]2[bc]\":       got=%q  expected \"aaabcbc\"\n", twoStacks("3[a]2[bc]"))            // expected aaabcbc
	fmt.Printf("s=\"3[a2[c]]\":        got=%q  expected \"accaccacc\"\n", twoStacks("3[a2[c]]"))           // expected accaccacc
	fmt.Printf("s=\"2[abc]3[cd]ef\":   got=%q  expected \"abcabccdcdcdef\"\n", twoStacks("2[abc]3[cd]ef")) // expected abcabccdcdcdef

	fmt.Println("=== Approach 2: Recursive Descent ===")
	fmt.Printf("s=\"3[a]2[bc]\":       got=%q  expected \"aaabcbc\"\n", recursiveDescent("3[a]2[bc]"))            // expected aaabcbc
	fmt.Printf("s=\"3[a2[c]]\":        got=%q  expected \"accaccacc\"\n", recursiveDescent("3[a2[c]]"))           // expected accaccacc
	fmt.Printf("s=\"2[abc]3[cd]ef\":   got=%q  expected \"abcabccdcdcdef\"\n", recursiveDescent("2[abc]3[cd]ef")) // expected abcabccdcdcdef
}
