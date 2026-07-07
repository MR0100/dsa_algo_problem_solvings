package main

import "fmt"

// ListNode is a singly-linked list node as defined by LeetCode.
type ListNode struct {
	Val  int
	Next *ListNode
}

// ── Approach 1: Reverse Both Lists, Add with Prepend ─────────────────────────
//
// reverseApproach solves Add Two Numbers II by reversing both inputs so the
// least-significant digit comes first (as in LeetCode #2), then adding with
// carry while PREPENDING each result digit — which yields MSB-first output with
// no separate final reversal.
//
// Intuition:
//
//	Addition is naturally done from the least-significant digit, but the digits
//	are stored most-significant-first. Reverse both lists to expose the units
//	digit at the head and add digit-by-digit propagating carry. If we build the
//	answer by prepending (new node points at the current head), the digits land
//	MSB-first automatically — the units digit is inserted first and ends up
//	deepest, the final carry is inserted last and becomes the head.
//
// Algorithm:
//  1. Reverse l1 and l2 (LSB now at the head).
//  2. Walk both, summing val1+val2+carry; new node val = sum%10, carry = sum/10.
//  3. Prepend each new node to the growing result.
//  4. Continue while either list remains or carry != 0; return the result head.
//
// Time:  O(m + n) — two reversals plus one addition pass, all linear.
// Space: O(1) extra (ignoring the output list) — reversal is in place.
//
// Note: mutates the input lists (they end up reversed). Acceptable unless the
// follow-up forbids it; main() rebuilds fresh lists per call.
func reverseApproach(l1 *ListNode, l2 *ListNode) *ListNode {
	l1 = reverseList(l1) // now LSB-first
	l2 = reverseList(l2)

	var result *ListNode // grows via prepend → ends up MSB-first
	carry := 0
	// Add while digits remain in either list or a carry is pending.
	for l1 != nil || l2 != nil || carry != 0 {
		sum := carry
		if l1 != nil {
			sum += l1.Val
			l1 = l1.Next
		}
		if l2 != nil {
			sum += l2.Val
			l2 = l2.Next
		}
		carry = sum / 10 // carry into the next (more significant) digit
		// Prepend the new digit; prepending an LSB-first computation yields
		// MSB-first output automatically.
		result = &ListNode{Val: sum % 10, Next: result}
	}
	return result
}

// reverseList reverses a singly linked list in place and returns the new head.
func reverseList(head *ListNode) *ListNode {
	var prev *ListNode
	for head != nil {
		next := head.Next // remember the rest
		head.Next = prev  // flip the pointer backward
		prev = head       // advance prev
		head = next       // advance head
	}
	return prev // new head = old tail
}

// ── Approach 2: Two Stacks (No Reversal, Optimal) ────────────────────────────
//
// twoStacks solves Add Two Numbers II by pushing every digit of each list onto a
// stack, then popping to add from least-significant to most-significant — without
// mutating the input lists (answers the follow-up).
//
// Intuition:
//
//	A stack gives us the digits in reverse (LSB-first) on pop, so we get the same
//	right-to-left addition as reversal but WITHOUT touching the inputs. Popping
//	both stacks yields the units digits first; we build the result by PREPENDING
//	each new digit, which keeps the final list MSB-first.
//
// Algorithm:
//  1. Push all of l1's values onto stack s1, all of l2's onto s2.
//  2. While either stack is non-empty or carry != 0: pop a value from each (0 if
//     empty), sum with carry, prepend a node with sum%10, set carry = sum/10.
//  3. Return the head (the most-significant digit, built last via prepend).
//
// Time:  O(m + n) — one pass to fill the stacks, one to drain them.
// Space: O(m + n) — the two stacks.
func twoStacks(l1 *ListNode, l2 *ListNode) *ListNode {
	s1 := listToStack(l1) // digits of number 1, LSB on top
	s2 := listToStack(l2) // digits of number 2, LSB on top

	var result *ListNode // grows via prepend → stays MSB-first
	carry := 0
	// Continue while any digits remain to add or a carry is outstanding.
	for len(s1) > 0 || len(s2) > 0 || carry != 0 {
		sum := carry
		if len(s1) > 0 {
			sum += s1[len(s1)-1] // top of stack = current least-significant digit
			s1 = s1[:len(s1)-1]  // pop
		}
		if len(s2) > 0 {
			sum += s2[len(s2)-1]
			s2 = s2[:len(s2)-1]
		}
		carry = sum / 10
		// Prepend keeps the most-significant digit at the head as we go.
		result = &ListNode{Val: sum % 10, Next: result}
	}
	return result
}

// listToStack copies a list's values into a slice used as a stack (index 0 = MSB,
// last index = LSB, so popping from the end walks LSB → MSB).
func listToStack(head *ListNode) []int {
	stack := []int{}
	for head != nil {
		stack = append(stack, head.Val)
		head = head.Next
	}
	return stack
}

// ── Helpers to build lists and print them for verification ───────────────────

// buildList makes a linked list from digits given most-significant-first.
func buildList(digits ...int) *ListNode {
	dummy := &ListNode{}
	cur := dummy
	for _, d := range digits {
		cur.Next = &ListNode{Val: d}
		cur = cur.Next
	}
	return dummy.Next
}

// listString renders a list as "7 -> 8 -> 0 -> 7".
func listString(head *ListNode) string {
	if head == nil {
		return "<nil>"
	}
	s := ""
	for head != nil {
		s += fmt.Sprintf("%d", head.Val)
		if head.Next != nil {
			s += " -> "
		}
		head = head.Next
	}
	return s
}

func main() {
	fmt.Println("=== Approach 1: Reverse Both, Add with Prepend ===")
	// 7243 + 564 = 7807
	fmt.Printf("[7 2 4 3] + [5 6 4]  got=%s  expected 7 -> 8 -> 0 -> 7\n",
		listString(reverseApproach(buildList(7, 2, 4, 3), buildList(5, 6, 4))))
	// 0 + 0 = 0
	fmt.Printf("[0] + [0]            got=%s  expected 0\n",
		listString(reverseApproach(buildList(0), buildList(0))))
	// 5 + 5 = 10 (carry creates a new most-significant digit)
	fmt.Printf("[5] + [5]            got=%s  expected 1 -> 0\n",
		listString(reverseApproach(buildList(5), buildList(5))))
	// 99 + 1 = 100 (cascading carry)
	fmt.Printf("[9 9] + [1]          got=%s  expected 1 -> 0 -> 0\n",
		listString(reverseApproach(buildList(9, 9), buildList(1))))

	fmt.Println("=== Approach 2: Two Stacks (No Reversal, Optimal) ===")
	fmt.Printf("[7 2 4 3] + [5 6 4]  got=%s  expected 7 -> 8 -> 0 -> 7\n",
		listString(twoStacks(buildList(7, 2, 4, 3), buildList(5, 6, 4))))
	fmt.Printf("[0] + [0]            got=%s  expected 0\n",
		listString(twoStacks(buildList(0), buildList(0))))
	fmt.Printf("[5] + [5]            got=%s  expected 1 -> 0\n",
		listString(twoStacks(buildList(5), buildList(5))))
	fmt.Printf("[9 9] + [1]          got=%s  expected 1 -> 0 -> 0\n",
		listString(twoStacks(buildList(9, 9), buildList(1))))
}
