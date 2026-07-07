package main

import (
	"fmt"
	"sort"
	"strconv"
)

// ── Approach 1: Plain Divide and Conquer ─────────────────────────────────────
//
// divideAndConquer solves Different Ways to Add Parentheses by splitting the
// expression at each operator and combining every left result with every right
// result.
//
// Intuition:
//
//	Every full parenthesization corresponds to choosing a "last" operator that
//	is evaluated last. If we fix that operator, the operands on its left form a
//	sub-expression and the operands on its right form another sub-expression,
//	each of which can itself be parenthesized in every possible way. So we
//	recurse on both sides and take the Cartesian product of their result lists,
//	applying the fixed operator to every pair.
//
// Algorithm:
//  1. If the string is a pure number (no operator), return [number].
//  2. For each position i holding an operator:
//     a. Recurse on expr[:i]   → list of left values.
//     b. Recurse on expr[i+1:] → list of right values.
//     c. For every (l, r) pair, apply the operator and append to results.
//  3. Return the accumulated results.
//
// Time:  O(Catalan(n)) — the number of distinct parenthesizations grows like
//
//	the Catalan number in n (number of operators); every combination is
//	produced once.
//
// Space: O(Catalan(n)) for the output plus O(n) recursion depth.
func divideAndConquer(expression string) []int {
	// Base case: no operator present means the whole string is one integer.
	if isNumber(expression) {
		n, _ := strconv.Atoi(expression) // parse the standalone number
		return []int{n}
	}

	var results []int
	// Try every operator position as the last-evaluated operator.
	for i := 0; i < len(expression); i++ {
		c := expression[i]
		if c == '+' || c == '-' || c == '*' {
			// Solve both sides independently.
			left := divideAndConquer(expression[:i])
			right := divideAndConquer(expression[i+1:])
			// Combine every left value with every right value.
			for _, l := range left {
				for _, r := range right {
					results = append(results, apply(l, r, c))
				}
			}
		}
	}
	return results
}

// ── Approach 2: Divide and Conquer + Memoization (Optimal) ───────────────────
//
// memoized solves Different Ways to Add Parentheses like the plain recursion
// but caches results for each substring so identical sub-expressions are
// computed only once.
//
// Intuition:
//
//	The same substring (e.g. "2*3") can be reached through many different top
//	level splits, and each time it is re-solved from scratch. Because the set
//	of results for a substring depends only on the substring itself, we can
//	store it in a map keyed by the substring and reuse it.
//
// Algorithm:
//  1. If the substring is cached, return the cached slice.
//  2. Otherwise run the same split-and-combine logic as Approach 1.
//  3. Store the result under the substring key before returning.
//
// Time:  O(Catalan(n)) results still, but each distinct substring is expanded
//
//	once; the memo removes the exponential re-computation of shared
//	sub-expressions.
//
// Space: O(number of distinct substrings × results) for the cache.
func memoized(expression string) []int {
	memo := make(map[string][]int) // substring → its list of possible values
	var solve func(expr string) []int
	solve = func(expr string) []int {
		if v, ok := memo[expr]; ok {
			return v // reuse previously computed results
		}
		if isNumber(expr) {
			n, _ := strconv.Atoi(expr)
			memo[expr] = []int{n}
			return memo[expr]
		}
		var results []int
		for i := 0; i < len(expr); i++ {
			c := expr[i]
			if c == '+' || c == '-' || c == '*' {
				left := solve(expr[:i])    // cached recursion
				right := solve(expr[i+1:]) // cached recursion
				for _, l := range left {
					for _, r := range right {
						results = append(results, apply(l, r, c))
					}
				}
			}
		}
		memo[expr] = results // cache before returning
		return results
	}
	return solve(expression)
}

// isNumber reports whether s contains no operator (i.e. is a plain integer).
func isNumber(s string) bool {
	for i := 0; i < len(s); i++ {
		if s[i] == '+' || s[i] == '-' || s[i] == '*' {
			return false
		}
	}
	return true
}

// apply computes l op r for a single operator byte.
func apply(l, r int, op byte) int {
	switch op {
	case '+':
		return l + r
	case '-':
		return l - r
	default: // '*'
		return l * r
	}
}

// sorted returns a sorted copy so outputs are comparable regardless of the
// order in which parenthesizations were generated.
func sorted(v []int) []int {
	out := append([]int(nil), v...)
	sort.Ints(out)
	return out
}

func main() {
	fmt.Println("=== Approach 1: Divide and Conquer ===")
	fmt.Println(sorted(divideAndConquer("2-1-1")))   // expected [0 2]
	fmt.Println(sorted(divideAndConquer("2*3-4*5"))) // expected [-34 -14 -10 -10 10]

	fmt.Println("=== Approach 2: Divide and Conquer + Memoization (Optimal) ===")
	fmt.Println(sorted(memoized("2-1-1")))   // expected [0 2]
	fmt.Println(sorted(memoized("2*3-4*5"))) // expected [-34 -14 -10 -10 10]
}
