package main

import "fmt"

// ── Approach 1: Stack of (result, sign) Contexts ─────────────────────────────
//
// stackCalculator evaluates an expression containing non-negative integers,
// '+', '-', '(' , ')' and spaces by keeping a running result and sign, and
// pushing/popping the "outer" context each time a parenthesis group starts or
// ends. There is no '*' or '/', so no precedence handling is needed beyond
// parentheses.
//
// Intuition:
//
//	Scan left to right accumulating `result += sign * number`. A '(' begins a
//	fresh sub-expression whose value must later be multiplied by the sign in
//	front of it and added to whatever we had before the '(' — so push the
//	current (result, sign) and reset. A ')' finishes the group: the inner
//	result is folded back as innerResult*savedSign + savedResult.
//
// Algorithm:
//
//	Maintain result=0, sign=+1, and a stack. For each char:
//	  digit  → build the current number.
//	  '+'    → flush number into result; sign = +1.
//	  '-'    → flush number into result; sign = −1.
//	  '('    → flush; push result then sign; reset result=0, sign=+1.
//	  ')'    → flush; result = result*poppedSign + poppedResult.
//	Flush the trailing number at the end and return result.
//
// Time:  O(n) — one pass over the string.
// Space: O(n) — stack depth up to the nesting depth of parentheses.
func stackCalculator(s string) int {
	result := 0       // running total of the current (sub)expression
	number := 0       // integer currently being parsed
	sign := 1         // sign applied to the next number (+1 or −1)
	stack := []int{}  // holds (result, sign) pairs of enclosing contexts
	hasDigit := false // whether `number` currently holds parsed digits

	for i := 0; i < len(s); i++ {
		c := s[i]
		switch {
		case c >= '0' && c <= '9':
			number = number*10 + int(c-'0') // extend the multi-digit number
			hasDigit = true
		case c == '+':
			result += sign * number // commit the number just parsed
			number, hasDigit = 0, false
			sign = 1 // next number is added
		case c == '-':
			result += sign * number
			number, hasDigit = 0, false
			sign = -1 // next number is subtracted
		case c == '(':
			// entering a group: save context, then start fresh inside
			stack = append(stack, result, sign)
			result, sign = 0, 1
			number, hasDigit = 0, false
		case c == ')':
			result += sign * number // commit the last number inside the group
			number, hasDigit = 0, false
			// pop saved sign then saved result (pushed in that order)
			savedSign := stack[len(stack)-1]
			savedResult := stack[len(stack)-2]
			stack = stack[:len(stack)-2]
			// fold the inner value back with the sign that preceded '('
			result = savedResult + savedSign*result
		}
		// spaces fall through and are ignored
	}
	if hasDigit {
		result += sign * number // commit any trailing number
	}
	return result
}

// ── Approach 2: Recursive Descent (Parse Parentheses by Recursion) ───────────
//
// recursiveCalculator evaluates the same grammar with a recursive helper that
// consumes one parenthesised group per call, returning both the group's value
// and the index just past its ')'. This mirrors the stack approach but uses
// the call stack instead of an explicit one.
//
// Intuition:
//
//	Evaluating "…(sub)…" means: whenever we meet '(', recursively evaluate the
//	inside as its own expression, get its value, then continue as if that
//	value were a single number. Recursion naturally saves and restores the
//	surrounding (result, sign) via stack frames.
//
// Algorithm:
//
//	eval(i) scans from i keeping local result/number/sign:
//	  digit → build number; '+'/'-' → commit + set sign;
//	  '(' → (value, next) = eval(i+1); treat value as the number, jump to next;
//	  ')' → commit number, return (result, i) to the caller.
//	Top-level call returns the final result.
//
// Time:  O(n) — each character consumed once across all frames.
// Space: O(d) — recursion depth = parenthesis nesting depth.
func recursiveCalculator(s string) int {
	val, _ := eval(s, 0)
	return val
}

// eval evaluates the expression starting at index i until it hits a matching
// ')' or the end of the string, returning the value and the index of the
// character just after the ')' (or len(s) at top level).
func eval(s string, i int) (int, int) {
	result := 0 // running total for this parenthesis level
	number := 0 // integer being parsed
	sign := 1   // sign for the next number
	for i < len(s) {
		c := s[i]
		switch {
		case c >= '0' && c <= '9':
			number = number*10 + int(c-'0')
			i++
		case c == '+':
			result += sign * number
			number, sign = 0, 1
			i++
		case c == '-':
			result += sign * number
			number, sign = 0, -1
			i++
		case c == '(':
			// evaluate the nested group; its value plays the role of a number
			var inner int
			inner, i = eval(s, i+1)
			number = inner
		case c == ')':
			result += sign * number // commit last number of this level
			return result, i + 1    // hand back position after ')'
		default: // space
			i++
		}
	}
	result += sign * number // commit trailing number at top level
	return result, i
}

func main() {
	fmt.Println("=== Approach 1: Stack of (result, sign) Contexts ===")
	fmt.Println(stackCalculator("1 + 1"))               // expected 2
	fmt.Println(stackCalculator(" 2-1 + 2 "))           // expected 3
	fmt.Println(stackCalculator("(1+(4+5+2)-3)+(6+8)")) // expected 23

	fmt.Println("=== Approach 2: Recursive Descent ===")
	fmt.Println(recursiveCalculator("1 + 1"))               // expected 2
	fmt.Println(recursiveCalculator(" 2-1 + 2 "))           // expected 3
	fmt.Println(recursiveCalculator("(1+(4+5+2)-3)+(6+8)")) // expected 23
}
