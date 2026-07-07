package main

import (
	"fmt"
	"strings"
)

// ── Approach 1: Stack (Fold Complete Subtrees) ───────────────────────────────
//
// stackFold verifies a preorder serialization by repeatedly collapsing a
// completed subtree "x,#,#" back into a single "#" placeholder.
//
// Intuition:
//
//	A valid preorder is built from the pattern: a real node followed by its
//	(eventually) two complete children. Whenever we see the pattern
//	number,#,# it means a node whose both children are already resolved to
//	null — so that whole subtree behaves like a leaf's null slot. Replace it
//	with a single # and keep folding. A fully valid tree folds down to exactly
//	one # (the null-slot the root itself occupies from its parent's view).
//
// Algorithm:
//  1. Push tokens left→right onto a stack.
//  2. After each push, while the top three are "#","#",number (in stack
//     order that means number,#,# in reading order), pop all three and push #.
//  3. At the end the stack must be exactly ["#"].
//
// Time:  O(n) — every token is pushed once and popped at most once.
// Space: O(n) — the stack.
func stackFold(preorder string) bool {
	tokens := strings.Split(preorder, ",") // split the CSV serialization
	stack := []string{}                    // holds unresolved tokens
	for _, t := range tokens {
		stack = append(stack, t) // read one token
		// Collapse "number,#,#" (top-of-stack order: #, #, number) into "#".
		for len(stack) >= 3 &&
			stack[len(stack)-1] == "#" &&
			stack[len(stack)-2] == "#" &&
			stack[len(stack)-3] != "#" {
			stack = stack[:len(stack)-3] // pop the two nulls and the node
			stack = append(stack, "#")   // the completed subtree acts as a null
		}
	}
	// A valid full tree reduces to a single null placeholder.
	return len(stack) == 1 && stack[0] == "#"
}

// ── Approach 2: Slot Counting (Optimal) ──────────────────────────────────────
//
// slotCounting verifies the serialization by tracking available "child slots".
//
// Intuition:
//
//	Think of the tree as a set of edge slots waiting to be filled. We begin
//	with one slot (the root's own incoming edge). Every token consumes one
//	slot. A non-null node then *creates two new slots* for its two children;
//	a null node creates none. If we ever run out of slots before finishing,
//	the string is invalid. At the very end exactly zero slots may remain.
//
// Algorithm:
//  1. slots = 1.
//  2. For each token: if slots == 0 mid-stream → invalid (no place to attach).
//     Consume one slot (slots--). If the token is a number, add two slots.
//  3. Valid iff slots == 0 after all tokens.
//
// Time:  O(n) — one pass over the tokens.
// Space: O(n) for the split (O(1) extra if scanned in place).
func slotCounting(preorder string) bool {
	tokens := strings.Split(preorder, ",") // the sequence of nodes
	slots := 1                             // the root occupies one incoming slot
	for _, t := range tokens {
		if slots == 0 {
			return false // a token arrived with nowhere to attach
		}
		slots-- // this token fills one open slot
		if t != "#" {
			slots += 2 // a real node opens two child slots
		}
	}
	return slots == 0 // every slot must be exactly filled
}

func main() {
	fmt.Println("=== Approach 1: Stack (Fold Complete Subtrees) ===")
	fmt.Printf("%q -> got=%t  expected true\n", "9,3,4,#,#,1,#,#,2,#,6,#,#", stackFold("9,3,4,#,#,1,#,#,2,#,6,#,#"))
	fmt.Printf("%q -> got=%t  expected false\n", "1,#", stackFold("1,#"))
	fmt.Printf("%q -> got=%t  expected false\n", "9,#,#,1", stackFold("9,#,#,1"))
	fmt.Printf("%q -> got=%t  expected true\n", "#", stackFold("#"))

	fmt.Println("=== Approach 2: Slot Counting (Optimal) ===")
	fmt.Printf("%q -> got=%t  expected true\n", "9,3,4,#,#,1,#,#,2,#,6,#,#", slotCounting("9,3,4,#,#,1,#,#,2,#,6,#,#"))
	fmt.Printf("%q -> got=%t  expected false\n", "1,#", slotCounting("1,#"))
	fmt.Printf("%q -> got=%t  expected false\n", "9,#,#,1", slotCounting("9,#,#,1"))
	fmt.Printf("%q -> got=%t  expected true\n", "#", slotCounting("#"))
}
