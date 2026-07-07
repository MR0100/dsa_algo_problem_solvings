package main

import (
	"fmt"
	"strings"
)

// Given a non-negative integer as a decimal string `num` and an integer `k`,
// remove exactly k digits so the remaining string forms the SMALLEST possible
// number. Return it without leading zeros; the empty result is "0".

// ── Approach 1: Brute Force (remove one costly digit at a time) ──────────────
//
// bruteForce performs k passes; in each pass it deletes the first digit that is
// larger than the digit right after it (or the last digit if the string is
// non-decreasing).
//
// Intuition:
//
//	To shrink a number by deleting one digit, delete the first "descent" — the
//	first position i where num[i] > num[i+1]. Removing that peak lowers a more
//	significant place, which dominates any lower-place change. If no descent
//	exists the digits are non-decreasing, so the last digit is the least
//	valuable to keep — drop it. Repeat k times.
//
// Algorithm:
//  1. Repeat k times:
//     a. Scan for the first index i with num[i] > num[i+1].
//     b. Remove num[i]; if no such i, remove the final character.
//  2. Strip leading zeros; return "0" if nothing remains.
//
// Time:  O(k · n) — each of the k passes scans the up-to-n-length string.
// Space: O(n) — the working string copy.
func bruteForce(num string, k int) string {
	b := []byte(num) // mutable copy so we can splice out characters
	for ; k > 0; k-- {
		i := 0
		// Walk to the first spot where a digit exceeds its successor.
		for i < len(b)-1 && b[i] <= b[i+1] {
			i++
		}
		// Delete b[i]: if a descent was found, i is that peak; otherwise the
		// string is non-decreasing and i == len(b)-1 (the last, largest digit).
		b = append(b[:i], b[i+1:]...)
	}
	return normalize(string(b))
}

// ── Approach 2: Monotonic Increasing Stack (Optimal) ─────────────────────────
//
// monotonicStack builds the answer left to right, keeping a stack of kept
// digits in non-decreasing order and popping any bigger digit while removals
// remain.
//
// Intuition:
//
//	We want the kept digits to increase from most-significant to least, because
//	a smaller digit in an earlier place beats anything later. So as each new
//	digit arrives, pop previously-kept digits that are LARGER than it (spending
//	one removal each) — the new digit is a better, smaller occupant of that
//	place. Stop popping when we run out of removals. Any removals still unused
//	at the end are spent trimming from the tail (the largest kept digits).
//
// Algorithm:
//  1. For each digit d in num:
//     - While k>0 and stack non-empty and top>d: pop, k--.
//     - Push d.
//  2. If k>0 remain, drop the last k digits (tail is largest).
//  3. Strip leading zeros; return "0" if empty.
//
// Time:  O(n) — every digit is pushed once and popped at most once.
// Space: O(n) — the stack of kept digits.
func monotonicStack(num string, k int) string {
	stack := make([]byte, 0, len(num)) // kept digits, maintained non-decreasing
	for i := 0; i < len(num); i++ {
		d := num[i]
		// Pop any kept digit strictly greater than the incoming one, as long as
		// we still have removals: the new smaller digit improves that place.
		for k > 0 && len(stack) > 0 && stack[len(stack)-1] > d {
			stack = stack[:len(stack)-1]
			k--
		}
		stack = append(stack, d) // the incoming digit is now kept
	}
	// Any leftover removals: the remaining digits are non-decreasing, so the
	// biggest ones are at the end — chop them off.
	stack = stack[:len(stack)-k]
	return normalize(string(stack))
}

// ── Approach 3: Preallocated Char Array as an Explicit Stack ─────────────────
//
// arrayStack is the same monotonic-stack algorithm but written over a fixed
// []byte with a manual top pointer, avoiding slice re-slicing — a common
// interview-friendly form that makes the O(n) push/pop bookkeeping obvious.
//
// Intuition:
//
//	Identical greedy as Approach 2: keep a non-decreasing sequence of digits,
//	removing a larger predecessor whenever a smaller digit appears and removals
//	remain. Using an index `top` into a preallocated buffer just makes the stack
//	operations explicit and allocation-free.
//
// Algorithm:
//  1. Allocate result buffer of length n, top = 0.
//  2. For each digit d: while top>0 and k>0 and buf[top-1]>d: top--, k--. Then
//     buf[top] = d, top++.
//  3. Final kept length = top - k (spend any leftover removals off the tail).
//  4. Strip leading zeros over buf[:keptLen]; "0" if empty.
//
// Time:  O(n) — one push and at most one pop per digit.
// Space: O(n) — the result buffer.
func arrayStack(num string, k int) string {
	buf := make([]byte, len(num)) // preallocated stack storage
	top := 0                      // index one past the last kept digit
	for i := 0; i < len(num); i++ {
		d := num[i]
		// Discard larger kept digits while budget allows.
		for top > 0 && k > 0 && buf[top-1] > d {
			top--
			k--
		}
		buf[top] = d // push the incoming digit
		top++
	}
	keptLen := top - k // leftover removals trim the (largest) tail digits
	return normalize(string(buf[:keptLen]))
}

// normalize strips leading zeros from a digit string and returns "0" if the
// result would be empty. Shared by all three approaches.
func normalize(s string) string {
	s = strings.TrimLeft(s, "0") // "0200" -> "200", "000" -> ""
	if s == "" {
		return "0" // an all-removed or all-zero result represents the value 0
	}
	return s
}

func main() {
	fmt.Println("=== Approach 1: Brute Force ===")
	fmt.Printf("num=1432219 k=3 -> %q\n", bruteForce("1432219", 3)) // expected "1219"
	fmt.Printf("num=10200 k=1   -> %q\n", bruteForce("10200", 1))   // expected "200"
	fmt.Printf("num=10 k=2      -> %q\n", bruteForce("10", 2))      // expected "0"

	fmt.Println("=== Approach 2: Monotonic Stack ===")
	fmt.Printf("num=1432219 k=3 -> %q\n", monotonicStack("1432219", 3)) // expected "1219"
	fmt.Printf("num=10200 k=1   -> %q\n", monotonicStack("10200", 1))   // expected "200"
	fmt.Printf("num=10 k=2      -> %q\n", monotonicStack("10", 2))      // expected "0"

	fmt.Println("=== Approach 3: Array Stack ===")
	fmt.Printf("num=1432219 k=3 -> %q\n", arrayStack("1432219", 3)) // expected "1219"
	fmt.Printf("num=10200 k=1   -> %q\n", arrayStack("10200", 1))   // expected "200"
	fmt.Printf("num=10 k=2      -> %q\n", arrayStack("10", 2))      // expected "0"
}
