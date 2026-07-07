package main

import "fmt"

// ── Approach 1: Stack ────────────────────────────────────────────────────────
//
// stackEval evaluates the expression by pushing terms onto a stack, applying
// '*' and '/' immediately and deferring '+'/'-' as signed numbers, then summing.
//
// Intuition:
//
//	Multiplication and division bind tighter than addition and subtraction.
//	If we scan left to right and remember the operator that PRECEDED the
//	current number, we can resolve precedence with a stack: for '+' push the
//	number, for '-' push its negation, and for '*' or '/' pop the top and push
//	the combined result. Whatever is left on the stack are independent additive
//	terms — summing them gives the answer.
//
// Algorithm:
//  1. Track prevOp (the operator before the current number), starting as '+'.
//  2. Walk each character, building a multi-digit number `num`.
//  3. When we hit an operator or the end of string, act on prevOp:
//     '+' → push num; '-' → push -num; '*' → push pop*num; '/' → push pop/num.
//  4. Reset num, set prevOp to the current operator.
//  5. Sum the stack.
//
// Time:  O(n) — one pass over the string, each element pushed/popped O(1) times.
// Space: O(n) — the stack can hold up to one entry per additive term.
func stackEval(s string) int {
	stack := []int{} // holds resolved additive terms (already signed / multiplied)
	num := 0         // the integer currently being parsed
	prevOp := byte('+')

	for i := 0; i < len(s); i++ {
		c := s[i]
		if c >= '0' && c <= '9' {
			num = num*10 + int(c-'0') // extend the multi-digit number
		}
		// Act at an operator OR at the final character (flush the last number).
		// Spaces are skipped by the digit check but still trigger the flush only
		// at end-of-string, so we guard on "not a space" for operators.
		if (c != ' ' && c < '0') || i == len(s)-1 {
			switch prevOp {
			case '+':
				stack = append(stack, num) // additive term, keep sign positive
			case '-':
				stack = append(stack, -num) // subtraction = adding a negative
			case '*':
				top := stack[len(stack)-1]      // multiply into the pending term
				stack[len(stack)-1] = top * num // replace top with product
			case '/':
				top := stack[len(stack)-1]      // integer-divide the pending term
				stack[len(stack)-1] = top / num // Go truncates toward zero, as required
			}
			prevOp = c // this operator governs the NEXT number
			num = 0    // reset the number accumulator
		}
	}

	sum := 0
	for _, v := range stack { // remaining entries are independent additive terms
		sum += v
	}
	return sum
}

// ── Approach 2: O(1) Space (Running Accumulators) ────────────────────────────
//
// constantSpaceEval evaluates the expression without a stack by keeping a
// running total and the "last term" so '*'/'/' can rewrite it before it is
// committed to the total.
//
// Intuition:
//
//	The stack only ever needs its TOP element for '*' and '/', and everything
//	below the top is already just being summed. So replace the stack with two
//	ints: `result` (the committed sum of all finished terms) and `lastNum`
//	(the most recent term, still open to being multiplied/divided). When a '+'
//	or '-' arrives we know the previous term is finalized, so fold lastNum into
//	result. For '*' or '/' we combine lastNum with the new number in place.
//
// Algorithm:
//  1. Keep result = 0, lastNum = 0, prevOp = '+'.
//  2. Parse each number; at each operator or end:
//     '+'/'-' → add lastNum to result, set lastNum = ±num;
//     '*'     → lastNum *= num;  '/' → lastNum /= num.
//  3. After the loop, add lastNum to result and return it.
//
// Time:  O(n) — single pass.
// Space: O(1) — three integer accumulators, no stack.
func constantSpaceEval(s string) int {
	result := 0  // sum of all fully-committed additive terms
	lastNum := 0 // the current term, still open to * or /
	num := 0     // integer being parsed
	prevOp := byte('+')

	for i := 0; i < len(s); i++ {
		c := s[i]
		if c >= '0' && c <= '9' {
			num = num*10 + int(c-'0') // build the multi-digit number
		}
		if (c != ' ' && c < '0') || i == len(s)-1 {
			switch prevOp {
			case '+':
				result += lastNum // commit previous term
				lastNum = num     // new open term is +num
			case '-':
				result += lastNum // commit previous term
				lastNum = -num    // new open term is -num
			case '*':
				lastNum *= num // fold multiplication into the open term
			case '/':
				lastNum /= num // fold division into the open term (truncates toward zero)
			}
			prevOp = c
			num = 0
		}
	}
	result += lastNum // commit the final open term
	return result
}

func main() {
	fmt.Println("=== Approach 1: Stack ===")
	fmt.Println(stackEval("3+2*2"))     // 7
	fmt.Println(stackEval(" 3/2 "))     // 1
	fmt.Println(stackEval(" 3+5 / 2 ")) // 5

	fmt.Println("=== Approach 2: O(1) Space (Running Accumulators) ===")
	fmt.Println(constantSpaceEval("3+2*2"))     // 7
	fmt.Println(constantSpaceEval(" 3/2 "))     // 1
	fmt.Println(constantSpaceEval(" 3+5 / 2 ")) // 5
}
