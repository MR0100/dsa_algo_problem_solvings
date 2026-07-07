package main

import (
	"fmt"
	"sort"
	"strconv"
)

// ── Approach 1: Brute Force (Enumerate Operator Assignments) ─────────────────
//
// bruteForce solves Expression Add Operators by trying every assignment of
// {"", "+", "-", "*"} to the n-1 gaps between digits, then splitting the
// string accordingly, evaluating each full expression, and keeping those equal
// to target.
//
// Intuition:
//
//	Between consecutive digits there are three real operators plus "no
//	operator" (which glues digits into a multi-digit number). That's 4^(n-1)
//	strings to build. Build each candidate expression, evaluate it respecting
//	operator precedence (* before + and -), and collect the ones that hit the
//	target. Dead simple, but it re-parses and re-evaluates each full string.
//
// Algorithm:
//
//  1. For every gap mask in base-4 over n-1 gaps, insert the chosen operator
//     (or nothing) between digits to form an expression string.
//  2. Reject any expression whose numeric tokens have leading zeros.
//  3. Evaluate with precedence; if it equals target, keep it.
//
// Time:  O(4^n · n) — 4^(n-1) expressions, each parsed/evaluated in O(n).
// Space: O(n) recursion/parse plus the result list.
func bruteForce(num string, target int) []string {
	n := len(num)
	res := []string{}
	if n == 0 {
		return res
	}
	ops := []string{"", "+", "-", "*"} // "" means concatenate digits
	gaps := n - 1

	// total = 4^gaps combinations, enumerated by mixed-radix counting.
	total := 1
	for i := 0; i < gaps; i++ {
		total *= 4
	}

	for mask := 0; mask < total; mask++ {
		// Build the expression as a list of tokens (numbers and operators).
		expr := string(num[0])
		tokens := []string{string(num[0])} // running number/operator tokens
		valid := true
		m := mask
		for g := 0; g < gaps; g++ {
			op := ops[m%4]
			m /= 4
			expr += op + string(num[g+1])
			if op == "" {
				// Glue digit onto the last number token.
				tokens[len(tokens)-1] += string(num[g+1])
			} else {
				tokens = append(tokens, op, string(num[g+1]))
			}
		}
		// Reject leading-zero numbers like "05".
		for i := 0; i < len(tokens); i += 2 {
			t := tokens[i]
			if len(t) > 1 && t[0] == '0' {
				valid = false
				break
			}
		}
		if valid && evaluate(tokens) == target {
			res = append(res, expr)
		}
	}
	sort.Strings(res) // deterministic order for testing
	return res
}

// evaluate computes a token list [num, op, num, op, ...] honouring that * binds
// tighter than + and -, using the standard "running term" technique.
func evaluate(tokens []string) int {
	result := 0                        // sum of completed +/- terms
	term, _ := strconv.Atoi(tokens[0]) // current multiplicative term
	for i := 1; i < len(tokens); i += 2 {
		op := tokens[i]
		v, _ := strconv.Atoi(tokens[i+1])
		switch op {
		case "+":
			result += term // close the term, start a new positive one
			term = v
		case "-":
			result += term
			term = -v
		case "*":
			term *= v // extend the current term
		}
	}
	return result + term // add the final open term
}

// ── Approach 2: Backtracking (Optimal) ───────────────────────────────────────
//
// backtracking solves Expression Add Operators by a DFS that, at each position,
// tries extending the current number and inserting +, -, or * — while carrying
// the running evaluated value so no re-parsing is ever needed.
//
// Intuition:
//
//	The killer detail is multiplication precedence. We track two numbers as we
//	build the expression: `total` (the value so far) and `prev` (the last
//	operand's signed contribution). For "+v": total += v, prev = v. For "-v":
//	total -= v, prev = -v. For "*v": we must undo prev and re-apply it times v
//	→ total = total - prev + prev*v, prev = prev*v. This makes each candidate
//	evaluated incrementally in O(1) at the leaf, and prunes leading-zero
//	numbers immediately.
//
// Algorithm:
//
//	DFS(pos, expr, total, prev):
//	  if pos == n: if total == target, record expr; return.
//	  for end = pos..n-1:
//	    take num[pos..end]; break if it has a leading zero (num[pos]=='0' and end>pos).
//	    if pos == 0: seed the first operand (no operator) → DFS(end+1, str, cur, cur).
//	    else try:
//	      +cur → DFS(end+1, expr+"+"+str, total+cur, cur)
//	      -cur → DFS(end+1, expr+"-"+str, total-cur, -cur)
//	      *cur → DFS(end+1, expr+"*"+str, total-prev+prev*cur, prev*cur)
//
// Time:  O(4^n · n) worst case (branching 4 per gap, O(n) to append the string).
// Space: O(n) recursion depth plus output.
func backtracking(num string, target int) []string {
	res := []string{}
	n := len(num)
	if n == 0 {
		return res
	}

	var dfs func(pos int, expr string, total, prev int)
	dfs = func(pos int, expr string, total, prev int) {
		if pos == n { // consumed all digits
			if total == target {
				res = append(res, expr)
			}
			return
		}
		for end := pos; end < n; end++ {
			// Leading-zero guard: "0" alone is fine, "05" is not.
			if end > pos && num[pos] == '0' {
				break
			}
			cur, _ := strconv.Atoi(num[pos : end+1]) // operand num[pos..end]
			str := num[pos : end+1]
			if pos == 0 {
				dfs(end+1, str, cur, cur) // first operand: no leading operator
			} else {
				dfs(end+1, expr+"+"+str, total+cur, cur)                // addition
				dfs(end+1, expr+"-"+str, total-cur, -cur)               // subtraction
				dfs(end+1, expr+"*"+str, total-prev+prev*cur, prev*cur) // fix precedence
			}
		}
	}
	dfs(0, "", 0, 0)
	sort.Strings(res) // deterministic order for testing
	return res
}

func main() {
	fmt.Println("=== Approach 1: Brute Force ===")
	fmt.Println(bruteForce("123", 6))           // [1*2*3 1+2+3]
	fmt.Println(bruteForce("232", 8))           // [2*3+2 2+3*2]
	fmt.Println(bruteForce("105", 5))           // [1*0+5 10-5]
	fmt.Println(bruteForce("00", 0))            // [0*0 0+0 0-0]
	fmt.Println(bruteForce("3456237490", 9191)) // []

	fmt.Println("=== Approach 2: Backtracking (Optimal) ===")
	fmt.Println(backtracking("123", 6))           // [1*2*3 1+2+3]
	fmt.Println(backtracking("232", 8))           // [2*3+2 2+3*2]
	fmt.Println(backtracking("105", 5))           // [1*0+5 10-5]
	fmt.Println(backtracking("00", 0))            // [0*0 0+0 0-0]
	fmt.Println(backtracking("3456237490", 9191)) // []
}
