package main

import (
	"fmt"
	"strconv"
)

// isOperator reports whether a token is one of the four RPN operators.
// A negative number like "-11" has length > 1, so it never matches.
func isOperator(token string) bool {
	return token == "+" || token == "-" || token == "*" || token == "/"
}

// apply evaluates `a op b` for one operator token.
// Go's integer division already truncates toward zero, as the problem demands.
func apply(a, b int, op string) int {
	switch op {
	case "+":
		return a + b
	case "-":
		return a - b
	case "*":
		return a * b
	default: // "/"
		return a / b // truncates toward zero for mixed signs, e.g. 6/-132 = 0
	}
}

// ── Approach 1: Brute Force (Repeated Scan and Reduce) ───────────────────────
//
// bruteForceReduce solves Evaluate Reverse Polish Notation by repeatedly
// finding the leftmost "number number operator" triple and collapsing it.
//
// Intuition:
//
//	An RPN expression contains at least one operator that directly follows
//	its two (already-literal) operands. Rewrite that triple as its computed
//	value and the expression shrinks by two tokens while staying valid RPN.
//	Repeat until a single number remains — exactly how you'd simplify the
//	expression by hand on paper.
//
// Algorithm:
//  1. Copy the tokens (we mutate the working slice).
//  2. While more than one token remains:
//     a. Scan left→right for the first operator whose two predecessors are
//     both numbers.
//     b. Evaluate the triple, splice the result back in place of the three
//     tokens.
//  3. Parse and return the lone remaining token.
//
// Time:  O(n^2) — each of the ~n/2 reductions rescans and reshifts O(n) tokens.
// Space: O(n) — the mutable copy of the token list.
func bruteForceReduce(tokens []string) int {
	work := make([]string, len(tokens))
	copy(work, tokens) // never mutate the caller's slice

	for len(work) > 1 {
		for i := 2; i < len(work); i++ {
			// a reducible triple is: number at i-2, number at i-1, operator at i
			if isOperator(work[i]) && !isOperator(work[i-2]) && !isOperator(work[i-1]) {
				a, _ := strconv.Atoi(work[i-2]) // left operand
				b, _ := strconv.Atoi(work[i-1]) // right operand
				result := apply(a, b, work[i])
				// splice: replace work[i-2 : i+1] with the single result token
				work[i-2] = strconv.Itoa(result)
				work = append(work[:i-1], work[i+1:]...)
				break // restart the scan on the shortened expression
			}
		}
	}
	value, _ := strconv.Atoi(work[0]) // single token left = the answer
	return value
}

// ── Approach 2: Recursion from the Right ─────────────────────────────────────
//
// recursiveEval solves Evaluate Reverse Polish Notation by consuming tokens
// right-to-left as an implicit expression tree.
//
// Intuition:
//
//	The LAST token of an RPN expression is the root of its expression tree.
//	If it is an operator, the tokens before it split into the right subtree
//	(immediately before the operator) followed by the left subtree. Walking
//	an index backwards, evaluate the RIGHT operand first, then the left —
//	each recursive call consumes exactly the tokens of its subtree.
//
// Algorithm:
//  1. Keep an index starting at the last token.
//  2. eval(): take the token at idx, decrement idx.
//     a. Operator → right := eval(), left := eval(), return apply(left, right, op).
//     b. Number  → return its parsed value.
//
// Time:  O(n) — every token is visited exactly once.
// Space: O(n) — recursion depth equals expression-tree height (worst case n).
func recursiveEval(tokens []string) int {
	idx := len(tokens) - 1 // shared cursor, consumed right-to-left
	var eval func() int
	eval = func() int {
		token := tokens[idx]
		idx-- // consume this token
		if isOperator(token) {
			right := eval() // right operand sits immediately before the operator
			left := eval()  // then everything before that is the left operand
			return apply(left, right, token)
		}
		value, _ := strconv.Atoi(token) // literal number: a leaf
		return value
	}
	return eval()
}

// ── Approach 3: Stack (Optimal) ──────────────────────────────────────────────
//
// stackApproach solves Evaluate Reverse Polish Notation with the classic
// single left-to-right pass over a stack of operands.
//
// Intuition:
//
//	RPN was designed for stack machines: a number is pushed; an operator pops
//	the top two values (top = RIGHT operand!), applies itself, and pushes the
//	result. Because every operator appears after its operands, its operands
//	are guaranteed to be the two most recent stack values. One pass, done.
//
// Algorithm:
//  1. For each token left→right:
//     a. Number   → push.
//     b. Operator → b := pop (right), a := pop (left), push apply(a, b, op).
//  2. The single value left on the stack is the result.
//
// Time:  O(n) — one pass, O(1) work per token.
// Space: O(n) — the operand stack (worst case: all numbers first).
func stackApproach(tokens []string) int {
	stack := make([]int, 0, len(tokens)) // operand stack
	for _, token := range tokens {
		if isOperator(token) {
			b := stack[len(stack)-1] // top of stack = RIGHT operand
			a := stack[len(stack)-2] // beneath it   = LEFT operand
			stack = stack[:len(stack)-2]
			stack = append(stack, apply(a, b, token)) // push the folded value
		} else {
			value, _ := strconv.Atoi(token) // literal number
			stack = append(stack, value)
		}
	}
	return stack[0] // valid RPN guarantees exactly one value remains
}

func main() {
	// Official LeetCode examples.
	example1 := []string{"2", "1", "+", "3", "*"}                                             // ((2+1)*3) = 9
	example2 := []string{"4", "13", "5", "/", "+"}                                            // (4+(13/5)) = 6
	example3 := []string{"10", "6", "9", "3", "+", "-11", "*", "/", "*", "17", "+", "5", "+"} // 22

	fmt.Println("=== Approach 1: Brute Force (Repeated Scan and Reduce) ===")
	fmt.Println(bruteForceReduce(example1)) // 9
	fmt.Println(bruteForceReduce(example2)) // 6
	fmt.Println(bruteForceReduce(example3)) // 22

	fmt.Println("=== Approach 2: Recursion from the Right ===")
	fmt.Println(recursiveEval(example1)) // 9
	fmt.Println(recursiveEval(example2)) // 6
	fmt.Println(recursiveEval(example3)) // 22

	fmt.Println("=== Approach 3: Stack (Optimal) ===")
	fmt.Println(stackApproach(example1)) // 9
	fmt.Println(stackApproach(example2)) // 6
	fmt.Println(stackApproach(example3)) // 22
}
