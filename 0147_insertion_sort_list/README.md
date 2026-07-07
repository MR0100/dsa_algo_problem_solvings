# 0147 — Insertion Sort List

> LeetCode #147 · Difficulty: Medium
> **Categories:** Linked List, Sorting

---

## Problem Statement

Given the `head` of a singly linked list, sort the list using **insertion sort**, and return the sorted list's head.

The steps of the **insertion sort** algorithm:

1. Insertion sort iterates, consuming one input element each repetition and growing a sorted output list.
2. At each iteration, insertion sort removes one element from the input data, finds the location it belongs within the sorted list and inserts it there.
3. It repeats until no input elements remain.

The following is a graphical example of the insertion sort algorithm. The partially sorted list (black) initially contains only the first element in the list. One element (red) is removed from the input data and inserted in-place into the sorted list with each iteration.

**Example 1:**
```
Input: head = [4,2,1,3]
Output: [1,2,3,4]
```

**Example 2:**
```
Input: head = [-1,5,3,4,0]
Output: [-1,0,3,4,5]
```

**Constraints:**
- The number of nodes in the list is in the range `[1, 5000]`.
- `-5000 <= Node.val <= 5000`

---

## Company Frequency

| Company    | Frequency        | Last Reported |
|------------|------------------|---------------|
| Microsoft  | ★★★★☆ High       | 2024          |
| Amazon     | ★★★☆☆ Medium     | 2024          |
| Google     | ★★★☆☆ Medium     | 2023          |
| Bloomberg  | ★★★☆☆ Medium     | 2023          |
| Adobe      | ★★☆☆☆ Low        | 2023          |

> ⚠️ Frequency data is crowd-sourced from LeetCode Discuss, Glassdoor, and
> community interview reports. Treat as a signal, not a guarantee.

---

## DSA Concepts Used

- **Linked List splicing** — detach a node and re-link it elsewhere with pure pointer surgery → see [`/dsa/linked_list.md`](/dsa/linked_list.md)
- **Sorting (insertion sort)** — the O(n²) grow-a-sorted-prefix algorithm, and why it's list-friendly → see [`/dsa/sorting.md`](/dsa/sorting.md)
- **Dummy (sentinel) head** — makes "insert before the current head" a non-special case → see [`/dsa/linked_list.md`](/dsa/linked_list.md)

---

## Approaches Overview

| # | Approach | Time | Space | When to use |
|---|----------|------|-------|-------------|
| 1 | Brute Force (Copy Values + Array Sort) | O(n log n) | O(n) | Baseline sanity check; ignores the exercise |
| 2 | Insertion Sort (Classic) | O(n²) | O(1) | The asked-for algorithm; clean reference version |
| 3 | Insertion Sort + Tail Shortcut (Optimal) | O(n²) worst, O(n) best | O(1) | Same algorithm, fast on nearly-sorted input |

---

## Approach 1 — Brute Force (Copy Values + Array Sort)

### Intuition
If we ignore the "use insertion sort" instruction, the path of least resistance is to move the values into an array — where sorting is a stdlib one-liner — and then write them back over the same nodes in order. The links never change; only the values do. It's a useful baseline and a reminder that "sort a linked list" is trivial once you allow O(n) extra space.

### Algorithm
1. Traverse the list, appending each `Val` to a slice.
2. `sort.Ints(vals)` — O(n log n).
3. Traverse the list a second time, overwriting `n.Val` with `vals[i]` in order.
4. Return the original head.

### Complexity
- **Time:** O(n log n) — two O(n) walks plus the array sort, which dominates.
- **Space:** O(n) — the values slice; this defeats the point of an in-place list sort.

### Code
```go
func bruteForceArraySort(head *ListNode) *ListNode {
	var vals []int
	for n := head; n != nil; n = n.Next {
		vals = append(vals, n.Val)
	}
	sort.Ints(vals)
	i := 0
	for n := head; n != nil; n = n.Next {
		n.Val = vals[i]
		i++
	}
	return head
}
```

### Dry Run
Example 1: `head = [4,2,1,3]`.

| Step | Action                | State                                  |
|------|-----------------------|----------------------------------------|
| 1    | collect values        | `vals = [4,2,1,3]`, list `4→2→1→3`     |
| 2    | `sort.Ints`           | `vals = [1,2,3,4]`                     |
| 3    | write back `i=0`      | node1.Val = 1 → list `1→2→1→3`         |
| 4    | write back `i=1`      | node2.Val = 2 → list `1→2→1→3`         |
| 5    | write back `i=2`      | node3.Val = 3 → list `1→2→3→3`         |
| 6    | write back `i=3`      | node4.Val = 4 → list `1→2→3→4`         |

Output: `[1,2,3,4]` ✓

---

## Approach 2 — Insertion Sort (Classic)

### Intuition
Sorting cards in your hand: pick up the next card from the pile (the unsorted remainder of the list), scan your hand (the sorted prefix) left-to-right, and slot the card in front of the first bigger card. On a linked list this is *cheaper* than on an array — an array insertion shifts every later element, but a list insertion is O(1) pointer surgery once the position is found. The scan itself is what costs O(n).

A **dummy head** in front of the sorted result means "insert at the very beginning" needs no special case: every insertion is "splice after some node `scan`".

### Algorithm
1. Create `dummy` (sentinel). `dummy.Next` is the always-sorted output list, initially empty.
2. Let `curr = head`. While `curr != nil`:
   1. Save `next := curr.Next` (the rest of the unsorted input) — `curr` is about to be re-linked.
   2. Scan from `dummy`: advance `scan` while `scan.Next != nil && scan.Next.Val < curr.Val`. Now `scan` is the last sorted node smaller than `curr` (or `dummy`).
   3. Splice: `curr.Next = scan.Next; scan.Next = curr`.
   4. `curr = next`.
3. Return `dummy.Next`.

### Complexity
- **Time:** O(n²) — inserting the i-th node scans up to i sorted nodes; Σi = n(n−1)/2. (Always scanning from the head also forfeits insertion sort's O(n) best case — fixed in Approach 3.)
- **Space:** O(1) — nodes are re-linked in place; only three pointers used.

### Code
```go
func insertionSortList(head *ListNode) *ListNode {
	dummy := &ListNode{} // sentinel head of the sorted result
	curr := head
	for curr != nil {
		next := curr.Next // save the unsorted remainder

		scan := dummy // find last sorted node with value < curr.Val
		for scan.Next != nil && scan.Next.Val < curr.Val {
			scan = scan.Next
		}

		curr.Next = scan.Next // splice curr in after scan
		scan.Next = curr

		curr = next
	}
	return dummy.Next
}
```

### Dry Run
Example 1: `head = [4,2,1,3]`. Sorted list shown after each insertion (dummy omitted).

| Step | `curr` | Scan stops because…                          | Splice action              | Sorted list after | `next` (remaining) |
|------|--------|----------------------------------------------|----------------------------|-------------------|---------------------|
| 1    | 4      | sorted list empty (`scan.Next == nil`)       | insert 4 after dummy       | `4`               | `2→1→3`            |
| 2    | 2      | `scan.Next.Val = 4 ≥ 2` immediately          | insert 2 before 4          | `2→4`             | `1→3`              |
| 3    | 1      | `scan.Next.Val = 2 ≥ 1` immediately          | insert 1 before 2          | `1→2→4`           | `3`                |
| 4    | 3      | passed 1, 2; stopped at `scan.Next.Val = 4 ≥ 3` | insert 3 between 2 and 4 | `1→2→3→4`         | `nil`              |

Loop ends (`curr == nil`). Return `dummy.Next` → `[1,2,3,4]` ✓

---

## Approach 3 — Insertion Sort + Tail Shortcut (Optimal)

### Intuition
The classic version has an embarrassing flaw: even if the input is already sorted, every node triggers a full scan from the head — O(n²) on the *easiest* input, exactly where insertion sort should shine at O(n). Fix: remember the **tail** of the sorted prefix. If the incoming value is `>= tail.Val`, it belongs at the end — append in O(1), no scan. Only genuinely out-of-order nodes pay for a scan. Worst case is unchanged, but sorted / nearly-sorted inputs (insertion sort's real-world use case) drop to ~O(n).

### Algorithm
1. `dummy` sentinel + `tail` pointer (nil while the sorted list is empty).
2. For each node `curr` (saving `next := curr.Next` first):
   1. **Fast path:** if `tail != nil && curr.Val >= tail.Val` → `tail.Next = curr; curr.Next = nil; tail = curr`.
   2. **Slow path:** scan from `dummy` for the first sorted node `≥ curr.Val`; splice `curr` in before it. If `curr` ended up last (`curr.Next == nil`), update `tail = curr`.
3. Return `dummy.Next`.

### Complexity
- **Time:** O(n²) worst case (reverse-sorted input still scans every time); **O(n) best case** — sorted input takes the fast path on every node.
- **Space:** O(1) — in-place, a constant number of pointers.

### Code
```go
func insertionSortTailOptimized(head *ListNode) *ListNode {
	dummy := &ListNode{}
	var tail *ListNode // last node of the sorted prefix
	for curr := head; curr != nil; {
		next := curr.Next

		if tail != nil && curr.Val >= tail.Val {
			// fast path: belongs at the end, no scan
			tail.Next = curr
			curr.Next = nil
			tail = curr
		} else {
			// slow path: scan for the insertion point
			scan := dummy
			for scan.Next != nil && scan.Next.Val < curr.Val {
				scan = scan.Next
			}
			curr.Next = scan.Next
			scan.Next = curr
			if curr.Next == nil {
				tail = curr
			}
		}

		curr = next
	}
	return dummy.Next
}
```

### Dry Run
Example 1: `head = [4,2,1,3]`.

| Step | `curr` | `tail` before | Path taken                                   | Sorted list after | `tail` after |
|------|--------|---------------|----------------------------------------------|-------------------|--------------|
| 1    | 4      | nil           | slow (list empty); insert after dummy; curr is last → tail=4 | `4`   | 4            |
| 2    | 2      | 4             | `2 >= 4`? no → slow; insert before 4         | `2→4`             | 4            |
| 3    | 1      | 4             | `1 >= 4`? no → slow; insert before 2         | `1→2→4`           | 4            |
| 4    | 3      | 4             | `3 >= 4`? no → slow; passes 1,2, stops at 4; insert before 4 | `1→2→3→4` | 4  |

Return `dummy.Next` → `[1,2,3,4]` ✓
(On an already-sorted input like `[1,2,3,4]`, steps 2–4 would all hit the fast path: zero scans.)

---

## Key Takeaways

- **Dummy head = zero special cases.** Any algorithm that may insert before the current head of a list should start from a sentinel node.
- **Save `next` before re-linking.** The single most common linked-list bug: re-wiring `curr.Next` before storing where the unsorted remainder continues.
- **Insertion on a list is O(1), search is O(n)** — the mirror image of arrays (search O(log n) with binary search, insertion O(n) shifting). That's why "binary insertion sort" helps arrays but not lists.
- **Tail check is the classic insertion-sort optimization** for lists: it restores the O(n) best case on nearly-sorted data, which is precisely when you'd choose insertion sort in practice.
- Insertion sort is the right tool for **small n, nearly-sorted data, or online streams** (sort as elements arrive); for general list sorting in O(n log n), see merge sort in LeetCode #148.

---

## Related Problems

- LeetCode #148 — Sort List (same task, O(n log n) merge sort — the follow-up)
- LeetCode #21 — Merge Two Sorted Lists (the splice-into-sorted-list primitive)
- LeetCode #92 — Reverse Linked List II (pointer surgery with a dummy head)
- LeetCode #86 — Partition List (stable re-linking of nodes into ordered groups)
- LeetCode #1019 — Next Greater Node In Linked List (list traversal + ordering)
