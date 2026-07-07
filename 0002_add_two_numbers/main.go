package main

import "fmt"

// ListNode is a singly-linked list node as defined by LeetCode.
type ListNode struct {
	Val  int
	Next *ListNode
}

// ── Approach 1: Convert to Numbers, Add, Convert Back ────────────────────────
//
// convertAndAdd reconstructs the integer each list represents, adds them,
// then builds a result linked list from the sum.
//
// Intuition:
//   The lists are just digits of a number stored in reverse. Reconstruct each
//   number, use built-in addition, then decompose the result back into a list.
//
// Why this is limited:
//   Go's int64 holds up to ~9.2 × 10¹⁸, which means lists longer than 18-19
//   digits silently overflow. For the LeetCode constraints (up to 100 nodes)
//   this approach is mathematically wrong in the general case.
//
// Algorithm:
//   1. Walk l1 accumulating: num1 += digit * 10^position.
//   2. Walk l2 accumulating: num2 += digit * 10^position.
//   3. sum = num1 + num2.
//   4. If sum == 0, return node(0). Otherwise build list digit by digit with sum % 10, sum /= 10.
//
// Time:  O(max(m,n)) — three linear passes (two to decode, one to encode).
// Space: O(max(m,n)) — the result list.
//
// ⚠️  INCORRECT for lists longer than ~18 nodes due to int64 overflow.
//     Included only for educational comparison.
func convertAndAdd(l1 *ListNode, l2 *ListNode) *ListNode {
	// Decode l1 into an integer.
	num1, mul := 0, 1
	for l1 != nil {
		num1 += l1.Val * mul
		mul *= 10
		l1 = l1.Next
	}
	// Decode l2 into an integer.
	num2, mul := 0, 1
	for l2 != nil {
		num2 += l2.Val * mul
		mul *= 10
		l2 = l2.Next
	}

	sum := num1 + num2
	if sum == 0 {
		return &ListNode{Val: 0}
	}

	// Encode sum back into a reversed-digit linked list.
	dummy := &ListNode{}
	cur := dummy
	for sum > 0 {
		cur.Next = &ListNode{Val: sum % 10}
		cur = cur.Next
		sum /= 10
	}
	return dummy.Next
}

// ── Approach 2: Iterative Simulation with Carry ───────────────────────────────
//
// iterativeCarry simulates column-by-column addition exactly as you would do
// it by hand, using a carry variable. This is the canonical solution.
//
// Intuition:
//   Because the digits are already stored in reverse (LSB first), the head of
//   each list is the ones digit — exactly where manual addition starts. Walk
//   both lists simultaneously, summing digit + digit + carry at each step.
//   Write the ones digit of that sum as a new node, propagate the carry.
//   Continue until both lists are exhausted AND there is no remaining carry.
//
// Algorithm:
//   1. Use a dummy head node to simplify list building (no nil-check on head).
//   2. carry = 0, cur = dummy.
//   3. While l1 != nil OR l2 != nil OR carry != 0:
//        sum = carry + (l1.Val if l1 != nil else 0) + (l2.Val if l2 != nil else 0)
//        carry = sum / 10
//        cur.Next = new node(sum % 10)
//        advance cur, l1, l2.
//   4. Return dummy.Next.
//
// Time:  O(max(m,n)) — one pass over both lists; at most one extra node for carry.
// Space: O(max(m,n)) — the result list (O(1) auxiliary beyond output).
func iterativeCarry(l1 *ListNode, l2 *ListNode) *ListNode {
	dummy := &ListNode{} // sentinel; dummy.Next will be the real head
	cur := dummy
	carry := 0

	// Continue while either list has digits OR there is a carry to flush.
	for l1 != nil || l2 != nil || carry != 0 {
		sum := carry // start with incoming carry

		if l1 != nil {
			sum += l1.Val
			l1 = l1.Next
		}
		if l2 != nil {
			sum += l2.Val
			l2 = l2.Next
		}

		carry = sum / 10          // 0 or 1 (max digit sum is 9+9+1=19)
		cur.Next = &ListNode{Val: sum % 10} // write the ones digit
		cur = cur.Next
	}

	return dummy.Next
}

// ── Approach 3: Recursive Simulation with Carry ───────────────────────────────
//
// recursiveCarry solves the problem using recursion. Each call handles one
// digit column and recurses on the remaining tails.
//
// Intuition:
//   The recursive structure mirrors the iterative one: at each level, consume
//   one node from each list (if available), compute sum + carry, create a
//   result node for the current digit, and recurse with the remaining tails
//   and the new carry. Base case: both lists nil and carry == 0.
//
// Algorithm:
//   helper(l1, l2, carry):
//     if l1 == nil AND l2 == nil AND carry == 0 → return nil
//     sum = carry + l1.Val (if l1 != nil) + l2.Val (if l2 != nil)
//     node = new ListNode(sum % 10)
//     node.Next = helper(l1.Next, l2.Next, sum/10)
//     return node
//
// Time:  O(max(m,n)) — one recursive call per digit column.
// Space: O(max(m,n)) — call stack depth equals the number of digit columns.
func recursiveCarry(l1 *ListNode, l2 *ListNode) *ListNode {
	return addHelper(l1, l2, 0)
}

// addHelper is the recursive worker for recursiveCarry.
func addHelper(l1, l2 *ListNode, carry int) *ListNode {
	// Base case: nothing left to process.
	if l1 == nil && l2 == nil && carry == 0 {
		return nil
	}

	sum := carry
	if l1 != nil {
		sum += l1.Val
		l1 = l1.Next
	}
	if l2 != nil {
		sum += l2.Val
		l2 = l2.Next
	}

	// Build this column's node and attach the rest recursively.
	node := &ListNode{Val: sum % 10}
	node.Next = addHelper(l1, l2, sum/10)
	return node
}

// ── Helpers ───────────────────────────────────────────────────────────────────

// makeList converts a slice of ints into a linked list (first element = head = LSB).
func makeList(vals []int) *ListNode {
	dummy := &ListNode{}
	cur := dummy
	for _, v := range vals {
		cur.Next = &ListNode{Val: v}
		cur = cur.Next
	}
	return dummy.Next
}

// listToSlice converts a linked list to a slice for easy printing.
func listToSlice(head *ListNode) []int {
	var result []int
	for head != nil {
		result = append(result, head.Val)
		head = head.Next
	}
	return result
}

func main() {
	type example struct {
		l1, l2 []int
		expect []int
	}
	examples := []example{
		{[]int{2, 4, 3}, []int{5, 6, 4}, []int{7, 0, 8}},       // 342+465=807
		{[]int{0}, []int{0}, []int{0}},                           // 0+0=0
		{[]int{9, 9, 9, 9, 9, 9, 9}, []int{9, 9, 9, 9}, []int{8, 9, 9, 9, 0, 0, 0, 1}}, // 9999999+9999
	}

	approaches := []struct {
		name string
		fn   func(*ListNode, *ListNode) *ListNode
	}{
		{"Approach 1: Convert & Add (⚠️ overflow risk) O(max(m,n)) T | O(max(m,n)) S", convertAndAdd},
		{"Approach 2: Iterative Carry (Optimal)        O(max(m,n)) T | O(max(m,n)) S", iterativeCarry},
		{"Approach 3: Recursive Carry                  O(max(m,n)) T | O(max(m,n)) S", recursiveCarry},
	}

	for _, ex := range examples {
		fmt.Printf("l1=%v  l2=%v  expect=%v\n", ex.l1, ex.l2, ex.expect)
		for _, ap := range approaches {
			result := listToSlice(ap.fn(makeList(ex.l1), makeList(ex.l2)))
			fmt.Printf("  %-75s → %v\n", ap.name, result)
		}
		fmt.Println()
	}
}
