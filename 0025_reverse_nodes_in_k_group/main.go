package main

import "fmt"

type ListNode struct {
	Val  int
	Next *ListNode
}

// ── Approach 1: Iterative ─────────────────────────────────────────────────────
//
// iterative reverses k nodes at a time using a prev/tail tracking approach.
//
// Intuition:
//   Use a dummy head. For each group of k nodes:
//     1. Check that k nodes exist ahead; if not, stop.
//     2. Reverse those k nodes in-place, keeping track of the group's tail
//        (which becomes its head after reversal).
//     3. Relink: connect the previous group's tail to the new head of this
//        reversed group.
//   After reversal, the original "head" of the group is now its tail →
//   advance prev to it for the next iteration.
//
// Time:  O(n) — each node reversed exactly once.
// Space: O(1) — in-place; only pointer variables.
func iterative(head *ListNode, k int) *ListNode {
	dummy := &ListNode{Next: head}
	prevGroupTail := dummy

	for {
		// Check if k nodes remain.
		kthNode := getKth(prevGroupTail, k)
		if kthNode == nil {
			break
		}

		groupHead := prevGroupTail.Next
		nextGroupHead := kthNode.Next

		// Reverse k nodes in [groupHead .. kthNode].
		prev := nextGroupHead
		cur := groupHead
		for cur != nextGroupHead {
			nxt := cur.Next
			cur.Next = prev
			prev = cur
			cur = nxt
		}

		// Relink: previous tail → new group head (was kthNode, now prev).
		prevGroupTail.Next = kthNode
		// groupHead is now the group's tail; advance prevGroupTail.
		prevGroupTail = groupHead
	}
	return dummy.Next
}

// getKth returns the k-th node after `node`, or nil if fewer than k exist.
func getKth(node *ListNode, k int) *ListNode {
	for k > 0 && node != nil {
		node = node.Next
		k--
	}
	return node
}

// ── Approach 2: Recursive ─────────────────────────────────────────────────────
//
// recursive reverses the first k nodes, then attaches the recursed tail.
//
// Intuition:
//   Count k nodes. If fewer than k exist, return head (no reversal).
//   Otherwise reverse k nodes: standard in-place reversal stopping after k.
//   Connect the original group head (now group tail) to the result of
//   reversing the rest of the list recursively.
//
// Time:  O(n) — each node visited once.
// Space: O(n/k) — recursion depth equals number of groups.
func recursive(head *ListNode, k int) *ListNode {
	// Count k nodes.
	count, cur := 0, head
	for cur != nil && count < k {
		cur = cur.Next
		count++
	}
	if count < k {
		return head // fewer than k nodes — don't reverse
	}

	// Reverse k nodes starting from head.
	var prev *ListNode
	curr := head
	for i := 0; i < k; i++ {
		nxt := curr.Next
		curr.Next = prev
		prev = curr
		curr = nxt
	}
	// head is now the tail of the reversed group.
	// curr is the start of the next group.
	head.Next = recursive(curr, k)
	return prev // prev is the new head of this reversed group
}

// ── Helpers ───────────────────────────────────────────────────────────────────

func makeList(vals []int) *ListNode {
	dummy := &ListNode{}
	cur := dummy
	for _, v := range vals {
		cur.Next = &ListNode{Val: v}
		cur = cur.Next
	}
	return dummy.Next
}

func listToSlice(head *ListNode) []int {
	var out []int
	for head != nil {
		out = append(out, head.Val)
		head = head.Next
	}
	return out
}

func main() {
	examples := []struct {
		vals   []int
		k      int
		expect []int
	}{
		{[]int{1, 2, 3, 4, 5}, 2, []int{2, 1, 4, 3, 5}},
		{[]int{1, 2, 3, 4, 5}, 3, []int{3, 2, 1, 4, 5}},
		{[]int{1, 2, 3, 4, 5}, 1, []int{1, 2, 3, 4, 5}},
		{[]int{1}, 1, []int{1}},
	}

	approaches := []struct {
		name string
		fn   func(*ListNode, int) *ListNode
	}{
		{"Approach 1: Iterative  ✅ O(n) T | O(1)   S", iterative},
		{"Approach 2: Recursive     O(n) T | O(n/k) S", recursive},
	}

	for _, ex := range examples {
		fmt.Printf("vals=%v  k=%d  expect=%v\n", ex.vals, ex.k, ex.expect)
		for _, ap := range approaches {
			fmt.Printf("  %-45s → %v\n", ap.name, listToSlice(ap.fn(makeList(ex.vals), ex.k)))
		}
		fmt.Println()
	}
}
