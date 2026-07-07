package main

import "fmt"

// ── Approach 1: Right-to-Left Stack ──────────────────────────────────────────
//
// stackRightToLeft solves Ternary Expression Parser by scanning the string from
// right to left, using a stack so that whenever a '?' appears its two already-
// resolved operands sit on top, ready to be collapsed by the preceding
// condition.
//
// Intuition:
//
//	Ternary groups right-to-left, so the RIGHTMOST '?' is the innermost complete
//	expression. Walking from the end, we push characters; the moment we hit a
//	'?', the top of the stack holds the true-branch and (below a ':') the
//	false-branch of that '?'. The character just before the '?' (the next one to
//	the left) is the condition — pick the branch and push the winner back. By
//	the time we reach index 0, the stack holds the single final value.
//
// Algorithm:
//  1. Iterate i from len-1 down to 0.
//  2. If s[i] == '?': pop the true-branch, pop the ':', pop the false-branch;
//     the condition is s[i-1]; push (cond=='T' ? trueBranch : falseBranch) and
//     skip i by one extra (consume the condition).
//  3. Otherwise push s[i].
//  4. The lone remaining stack element is the answer.
//
// Time:  O(n) — each character is pushed and popped O(1) times.
// Space: O(n) — the stack can hold up to the whole string in the worst case.
func stackRightToLeft(expression string) string {
	stack := []byte{} // holds resolved single characters (values / ':' markers)
	for i := len(expression) - 1; i >= 0; i-- {
		c := expression[i]
		if c == '?' {
			trueBranch := stack[len(stack)-1]  // top = value if condition true
			stack = stack[:len(stack)-1]       // pop true-branch
			stack = stack[:len(stack)-1]       // pop the ':' separator
			falseBranch := stack[len(stack)-1] // next = value if condition false
			stack = stack[:len(stack)-1]       // pop false-branch
			cond := expression[i-1]            // the condition sits just left of '?'
			if cond == 'T' {
				stack = append(stack, trueBranch) // condition true → keep true-branch
			} else {
				stack = append(stack, falseBranch) // condition false → keep false-branch
			}
			i-- // we consumed the condition character too; skip it
		} else {
			stack = append(stack, c) // digit, 'T', 'F', or ':' — defer resolution
		}
	}
	return string(stack[0]) // exactly one value remains
}

// ── Approach 2: Recursive Descent ────────────────────────────────────────────
//
// recursiveDescent solves Ternary Expression Parser by parsing with a moving
// index: read one atom (condition), and if it is followed by '?', recursively
// parse the true-branch and false-branch, skipping the branch not taken.
//
// Intuition:
//
//	Grammar: expr = value | value '?' expr ':' expr. Read the leading value; if
//	the next char is '?', this is a ternary. Evaluate the taken branch by
//	recursing; but we must still ADVANCE the index past the untaken branch so
//	the caller resumes at the right place. A depth counter skips a full nested
//	branch: increment on '?', decrement on ':', and the branch ends when depth
//	returns to its starting balance at a ':' or end.
//
// Algorithm:
//  1. Keep a package-style index via closure; parse() reads expression[pos].
//  2. cond = current char; advance. If next is not '?', return cond.
//  3. Otherwise skip '?'; parse the true-branch recursively; skip ':'; parse
//     the false-branch recursively; return whichever matches cond.
//     (Both branches are parsed to move the cursor, but only one value is kept.)
//
// Time:  O(n) — every character is read once by the recursive walk.
// Space: O(n) — recursion depth up to the nesting depth of the expression.
func recursiveDescent(expression string) string {
	pos := 0 // shared cursor into expression

	var parse func() string
	parse = func() string {
		cond := expression[pos] // first atom: a value or a condition (T/F)
		pos++                   // move past it
		// A bare value (no trailing '?') is a complete expression.
		if pos >= len(expression) || expression[pos] != '?' {
			return string(cond)
		}
		pos++               // skip '?'
		trueVal := parse()  // recursively evaluate the true-branch
		pos++               // skip ':'
		falseVal := parse() // recursively evaluate the false-branch
		if cond == 'T' {    // choose based on the condition
			return trueVal
		}
		return falseVal
	}

	return parse()
}

// ── Approach 3: Iterative Forward Skip (Optimal) ─────────────────────────────
//
// iterativeForwardSkip solves Ternary Expression Parser with a single forward
// pointer and NO recursion and NO stack of operands: at each ternary it keeps
// the taken branch and fast-forwards the pointer past the whole untaken branch
// using a nesting-depth counter.
//
// Intuition:
//
//	Read the leading condition. If it is 'T' we want the true-branch, which
//	starts right after '?', so just step in and keep going. If it is 'F' we must
//	discard the true-branch and jump to the false-branch, which begins after the
//	':' that matches THIS '?'. Finding that ':' means skipping any nested
//	ternaries: a depth counter rises on '?' and falls on ':', and the matching
//	separator is the ':' seen while depth is balanced. Symmetric logic skips the
//	trailing false-branch once we have finished a taken true-branch — but with a
//	forward "keep taken, skip untaken" rule the pointer always lands on the head
//	of the next value we care about, so the final character reached is the
//	answer.
//
// Algorithm:
//  1. i = 0.
//  2. Loop: if the char after expression[i] is not '?', expression[i] is the
//     final value — return it.
//  3. Otherwise expression[i] is a condition. If 'T', step i past "cond?" into
//     the true-branch. If 'F', skip past the true-branch (using depth) to just
//     after the matching ':' and continue from the false-branch.
//
// Time:  O(n) — the pointer only moves forward across the string.
// Space: O(1) — a couple of integer counters, no stack or recursion.
func iterativeForwardSkip(expression string) string {
	i := 0
	for {
		// A value with no following '?' is a complete (sub-)expression's result.
		if i+1 >= len(expression) || expression[i+1] != '?' {
			return string(expression[i])
		}
		if expression[i] == 'T' {
			// Condition true → descend into the true-branch (right after "T?").
			i += 2
		} else {
			// Condition false → skip the ENTIRE true-branch, land on false-branch.
			// Start just past "F?" and walk to the ':' that matches this '?'.
			i += 2     // move past "F?" to the first char of the true-branch
			depth := 0 // nesting level relative to this branch
			for depth > 0 || expression[i] != ':' {
				switch expression[i] {
				case '?':
					depth++ // entering a nested ternary
				case ':':
					depth-- // leaving a nested ternary
				}
				i++
			}
			i++ // step over the matching ':' onto the false-branch head
		}
	}
}

func main() {
	fmt.Println("=== Approach 1: Right-to-Left Stack ===")
	fmt.Println(stackRightToLeft("T?2:3"))     // expected 2
	fmt.Println(stackRightToLeft("F?1:T?4:5")) // expected 4
	fmt.Println(stackRightToLeft("T?T?F:5:3")) // expected F

	fmt.Println("=== Approach 2: Recursive Descent ===")
	fmt.Println(recursiveDescent("T?2:3"))     // expected 2
	fmt.Println(recursiveDescent("F?1:T?4:5")) // expected 4
	fmt.Println(recursiveDescent("T?T?F:5:3")) // expected F

	fmt.Println("=== Approach 3: Iterative Forward Skip (Optimal) ===")
	fmt.Println(iterativeForwardSkip("T?2:3"))     // expected 2
	fmt.Println(iterativeForwardSkip("F?1:T?4:5")) // expected 4
	fmt.Println(iterativeForwardSkip("T?T?F:5:3")) // expected F
}
