package main

import (
	"fmt"
	"strconv"
	"strings"
)

// NestedInteger mirrors LeetCode's interface: each instance holds EITHER a
// single integer OR a list of NestedInteger. isInt distinguishes the two.
type NestedInteger struct {
	isInt bool
	num   int
	list  []*NestedInteger
}

// NewInt builds a NestedInteger holding a single integer.
func NewInt(v int) *NestedInteger { return &NestedInteger{isInt: true, num: v} }

// NewList builds an (initially empty) NestedInteger holding a list.
func NewList() *NestedInteger { return &NestedInteger{isInt: false, list: []*NestedInteger{}} }

// IsInteger reports whether this NestedInteger holds a single integer (true)
// or a nested list (false).
func (n *NestedInteger) IsInteger() bool { return n.isInt }

// SetInteger sets this NestedInteger to hold a single integer.
func (n *NestedInteger) SetInteger(value int) { n.isInt = true; n.num = value }

// Add nests another NestedInteger inside this one (this becomes a list).
func (n *NestedInteger) Add(elem *NestedInteger) {
	n.isInt = false
	n.list = append(n.list, elem)
}

// GetInteger returns the held integer (valid only when IsInteger()).
func (n *NestedInteger) GetInteger() int { return n.num }

// GetList returns the held list (valid only when !IsInteger()).
func (n *NestedInteger) GetList() []*NestedInteger { return n.list }

// String renders a NestedInteger back into LeetCode's bracket notation,
// so we can print and compare parser output against the input string.
func (n *NestedInteger) String() string {
	if n.IsInteger() {
		return strconv.Itoa(n.GetInteger())
	}
	parts := make([]string, 0, len(n.list))
	for _, e := range n.list {
		parts = append(parts, e.String())
	}
	return "[" + strings.Join(parts, ",") + "]"
}

// ── Approach 1: Explicit Stack ───────────────────────────────────────────────
//
// stackParse deserializes the string using an explicit stack of open lists,
// scanning left to right one character at a time.
//
// Intuition:
//
//	Brackets nest, so the natural tool is a stack. If the string does not
//	start with '[', it is a bare integer — return it directly. Otherwise,
//	each '[' opens a new list pushed on the stack; a number token becomes a
//	NestedInteger added to the list on top; each ']' closes the top list and
//	nests it into the (new) top below. The lone remaining stack item at the
//	end is the answer.
//
// Algorithm:
//  1. If s[0] != '[', return NewInt(atoi(s)).
//  2. Walk s. Maintain a stack of *NestedInteger lists and a number buffer.
//  3. On '[': push a fresh list.
//  4. On ',' or ']': if a number was buffered, flush it as an int into the
//     top list; on ']' additionally pop the top list and Add it into the new
//     top (unless it was the root).
//  5. On a digit or '-': accumulate into the number buffer.
//  6. Return the single remaining list.
//
// Time:  O(L) — one pass over the L-character string (atoi is amortized O(len)).
// Space: O(D) — stack depth D = maximum nesting; O(L) worst case.
func stackParse(s string) *NestedInteger {
	if s[0] != '[' { // a bare integer like "324" or "-7"
		v, _ := strconv.Atoi(s)
		return NewInt(v)
	}

	stack := []*NestedInteger{} // open lists, innermost on top
	numStart := -1              // start index of the current number token, -1 = none

	// flushNumber turns any pending digit run into an int added to the top list.
	flushNumber := func(end int) {
		if numStart != -1 {
			v, _ := strconv.Atoi(s[numStart:end])
			top := stack[len(stack)-1]
			top.Add(NewInt(v)) // nest the integer into the current list
			numStart = -1
		}
	}

	for i := 0; i < len(s); i++ {
		switch c := s[i]; {
		case c == '[':
			stack = append(stack, NewList()) // open a new list
		case c == ',':
			flushNumber(i) // a number token (if any) just ended
		case c == ']':
			flushNumber(i) // close out any trailing number in this list
			if len(stack) > 1 {
				top := stack[len(stack)-1]
				stack = stack[:len(stack)-1] // pop the finished list
				stack[len(stack)-1].Add(top) // nest it into its parent
			}
		default: // digit or leading '-'
			if numStart == -1 {
				numStart = i // begin a new number token
			}
		}
	}
	return stack[0] // the root list
}

// ── Approach 2: Recursive Descent (Optimal / cleanest) ───────────────────────
//
// recursiveParse deserializes the string with a recursive-descent parser that
// mirrors the grammar element = int | '[' (element (',' element)*)? ']'.
//
// Intuition:
//
//	The data is a grammar, so parse it recursively. A shared cursor `i` walks
//	the string. parseValue looks at s[i]: if it is not '[', read an integer
//	token; otherwise consume '[', then repeatedly parseValue for each comma-
//	separated child until ']'. Recursion naturally matches the nesting depth,
//	so no explicit stack is needed — the call stack IS the stack.
//
// Algorithm:
//  1. Keep an index i (closure variable).
//  2. parseValue: if s[i] != '[', scan a signed integer and wrap in NewInt.
//  3. Else consume '['; make an empty list; while s[i] != ']', parseValue a
//     child, Add it, and skip a ',' if present; consume ']'; return the list.
//  4. Call parseValue once from index 0.
//
// Time:  O(L) — each character consumed once across the recursion.
// Space: O(D) — recursion depth = maximum nesting; O(L) worst case.
func recursiveParse(s string) *NestedInteger {
	i := 0 // shared cursor into s

	var parseValue func() *NestedInteger
	parseValue = func() *NestedInteger {
		if s[i] != '[' { // integer token: optional '-' then digits
			start := i
			if s[i] == '-' {
				i++
			}
			for i < len(s) && s[i] >= '0' && s[i] <= '9' {
				i++
			}
			v, _ := strconv.Atoi(s[start:i])
			return NewInt(v)
		}

		i++ // consume '['
		lst := NewList()
		for s[i] != ']' { // parse children until the matching ']'
			lst.Add(parseValue()) // recurse for each element
			if s[i] == ',' {
				i++ // skip the separator
			}
		}
		i++ // consume ']'
		return lst
	}

	return parseValue()
}

func main() {
	// Example 1: a bare integer.
	fmt.Println("=== Approach 1: Explicit Stack ===")
	fmt.Println(stackParse("324"))               // expected 324
	fmt.Println(stackParse("[123,[456,[789]]]")) // expected [123,[456,[789]]]
	fmt.Println(stackParse("[]"))                // expected []
	fmt.Println(stackParse("[-1,2,[3,-4]]"))     // expected [-1,2,[3,-4]]

	fmt.Println("=== Approach 2: Recursive Descent (Optimal) ===")
	fmt.Println(recursiveParse("324"))               // expected 324
	fmt.Println(recursiveParse("[123,[456,[789]]]")) // expected [123,[456,[789]]]
	fmt.Println(recursiveParse("[]"))                // expected []
	fmt.Println(recursiveParse("[-1,2,[3,-4]]"))     // expected [-1,2,[3,-4]]
}
