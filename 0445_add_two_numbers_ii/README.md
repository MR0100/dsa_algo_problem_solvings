# 0445 — Add Two Numbers II

> LeetCode #445 · Difficulty: Medium
> **Categories:** Linked List, Math, Stack

---

## Problem Statement

You are given two **non-empty** linked lists representing two non-negative integers. The most significant digit comes first and each of their nodes contains a single digit. Add the two numbers and return the sum as a linked list.

You may assume the two numbers do not contain any leading zero, except the number 0 itself.

**Example 1:**

```
Input: l1 = [7,2,4,3], l2 = [5,6,4]
Output: [7,8,0,7]
```

(7243 + 564 = 7807.)

**Example 2:**

```
Input: l1 = [2,4,3], l2 = [5,6,4]
Output: [8,0,7]
```

**Example 3:**

```
Input: l1 = [0], l2 = [0]
Output: [0]
```

**Constraints:**

- The number of nodes in each linked list is in the range `[1, 100]`.
- `0 <= Node.val <= 9`
- It is guaranteed that the list represents a number that does not have leading zeros.

**Follow up:** Could you solve it without reversing the input lists?

---

## Company Frequency

| Company    | Frequency        | Last Reported |
|------------|------------------|---------------|
| Amazon     | ★★★★☆ High       | 2024          |
| Microsoft  | ★★★★☆ High       | 2024          |
| Meta       | ★★★☆☆ Medium     | 2024          |
| Google     | ★★★☆☆ Medium     | 2023          |
| Bloomberg  | ★★★☆☆ Medium     | 2023          |
| Adobe      | ★★☆☆☆ Low        | 2022          |

> ⚠️ Frequency data is crowd-sourced from LeetCode Discuss, Glassdoor, and
> community interview reports. Treat as a signal, not a guarantee.

---

## DSA Concepts Used

- **Linked List traversal & construction** — building the result by *prepending* nodes so the most-significant digit stays at the head, plus in-place reversal → see [`/dsa/linked_list.md`](/dsa/linked_list.md)
- **Stack (LIFO to reverse order without mutation)** — pushing digits and popping them gives least-significant-first access while leaving the inputs untouched, answering the follow-up → see [`/dsa/stack.md`](/dsa/stack.md)

---

## Approaches Overview

| # | Approach | Time | Space | When to use |
|---|----------|------|-------|-------------|
| 1 | Reverse Both, Add with Prepend | O(m + n) | O(1) extra | Cleanest if mutating inputs is allowed |
| 2 | Two Stacks (No Reversal, Optimal) | O(m + n) | O(m + n) | Answers the follow-up — inputs stay intact |

---

## Approach 1 — Reverse Both, Add with Prepend

### Intuition

Addition wants the least-significant digit first, but the lists store the most-significant digit first. Reverse both lists so the units digit sits at the head — now it is exactly LeetCode #2. Add digit-by-digit with a carry, and **prepend** each new digit to the result. Prepending an LSB-first computation lands the digits MSB-first automatically: the units digit is inserted first and ends up deepest, while the final carry is inserted last and becomes the head — so no separate reversal of the answer is needed.

### Algorithm

1. Reverse `l1` and `l2` (LSB now at each head).
2. Loop while `l1 != nil` or `l2 != nil` or `carry != 0`: `sum = carry + (l1.Val?) + (l2.Val?)`; advance the consumed lists.
3. `carry = sum / 10`; prepend a node with `sum % 10` to `result`.
4. Return `result`.

### Complexity

- **Time:** O(m + n) — two reversals plus one addition pass, each linear.
- **Space:** O(1) extra beyond the output — reversal is in place (mutates the inputs).

### Code

```go
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
```

### Dry Run

Example 1: `l1 = 7→2→4→3` (7243), `l2 = 5→6→4` (564).

After reversal: `l1 = 3→4→2→7`, `l2 = 4→6→5`.

| Step | l1.Val | l2.Val | carry in | sum | digit = sum%10 | carry out | result (head → …) |
|------|--------|--------|----------|-----|----------------|-----------|--------------------|
| 1 | 3 | 4 | 0 | 7 | 7 | 0 | 7 |
| 2 | 4 | 6 | 0 | 10 | 0 | 1 | 0 → 7 |
| 3 | 2 | 5 | 1 | 8 | 8 | 0 | 8 → 0 → 7 |
| 4 | 7 | — | 0 | 7 | 7 | 0 | 7 → 8 → 0 → 7 |

Loop ends (both lists empty, carry 0). Result: `7 → 8 → 0 → 7` ✔.

---

## Approach 2 — Two Stacks (No Reversal, Optimal)

### Intuition

A stack hands digits back in reverse order — so pushing every digit and popping gives the same LSB-first traversal as reversal, but **without mutating the inputs** (the follow-up). Pop from both stacks to add units-first; build the result by **prepending** so the most-significant digit stays at the head.

### Algorithm

1. Push all of `l1`'s values onto stack `s1`, all of `l2`'s onto `s2` (tops = units digits).
2. While `s1` or `s2` is non-empty or `carry != 0`: pop a value from each (0 if empty), `sum = popped + popped + carry`; prepend a node with `sum % 10`; `carry = sum / 10`.
3. Return the head (the most-significant digit, prepended last).

### Complexity

- **Time:** O(m + n) — one pass to fill the stacks, one to drain them.
- **Space:** O(m + n) — the two stacks hold every digit; the inputs are left unchanged.

### Code

```go
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

func listToStack(head *ListNode) []int {
	stack := []int{}
	for head != nil {
		stack = append(stack, head.Val)
		head = head.Next
	}
	return stack
}
```

### Dry Run

Example 1: `l1 = 7→2→4→3`, `l2 = 5→6→4`.

Stacks (top on the right): `s1 = [7,2,4,3]`, `s2 = [5,6,4]`.

| Step | pop s1 | pop s2 | carry in | sum | digit = sum%10 | carry out | result (head → …) |
|------|--------|--------|----------|-----|----------------|-----------|--------------------|
| 1 | 3 | 4 | 0 | 7 | 7 | 0 | 7 |
| 2 | 4 | 6 | 0 | 10 | 0 | 1 | 0 → 7 |
| 3 | 2 | 5 | 1 | 8 | 8 | 0 | 8 → 0 → 7 |
| 4 | 7 | (empty→0) | 0 | 7 | 7 | 0 | 7 → 8 → 0 → 7 |

Both stacks empty, carry 0 → stop. Result: `7 → 8 → 0 → 7` ✔ — and `l1`, `l2` were never modified.

---

## Key Takeaways

- **"MSB-first, but add LSB-first" is the whole puzzle.** Two ways to bridge the mismatch: physically reverse the lists, or use a stack to *virtually* reverse the traversal. The stack version is preferred because it leaves the inputs intact (the follow-up).
- **Prepend, don't append.** Building the result with `result = &ListNode{Val: d, Next: result}` places each newly computed (less-significant-later) digit ahead of the previous ones, giving MSB-first order for free — no final reverse pass.
- **Carry outlives the lists.** The loop condition must include `carry != 0`, or `99 + 1 = 100` loses its leading `1`. Always keep summing while any input remains *or* a carry is pending.
- **Empty-stack / empty-list guards default the missing digit to 0**, cleanly handling unequal lengths without padding.

---

## Related Problems

- LeetCode #2 — Add Two Numbers (digits stored LSB-first; the reversed twin of this one)
- LeetCode #206 — Reverse Linked List (the reversal primitive used in Approach 1)
- LeetCode #369 — Plus One Linked List (single-operand carry propagation, MSB-first)
- LeetCode #67 — Add Binary (same carry loop over strings/bits)
- LeetCode #43 — Multiply Strings (grade-school arithmetic with carries)
