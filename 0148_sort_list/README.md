# 0148 — Sort List

> LeetCode #148 · Difficulty: Medium
> **Categories:** Linked List, Two Pointers, Divide and Conquer, Sorting, Merge Sort

---

## Problem Statement

Given the `head` of a linked list, return the list after sorting it in **ascending order**.

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

**Example 3:**
```
Input: head = []
Output: []
```

**Constraints:**
- The number of nodes in the list is in the range `[0, 5 * 10^4]`.
- `-10^5 <= Node.val <= 10^5`

**Follow-up:** Can you sort the linked list in `O(n logn)` time and `O(1)` memory (i.e. constant space)?

---

## Company Frequency

| Company    | Frequency        | Last Reported |
|------------|------------------|---------------|
| Facebook   | ★★★★★ Very High  | 2024          |
| Amazon     | ★★★★☆ High       | 2024          |
| Microsoft  | ★★★★☆ High       | 2024          |
| Google     | ★★★★☆ High       | 2024          |
| Adobe      | ★★★☆☆ Medium     | 2023          |
| Bloomberg  | ★★★☆☆ Medium     | 2023          |

> ⚠️ Frequency data is crowd-sourced from LeetCode Discuss, Glassdoor, and
> community interview reports. Treat as a signal, not a guarantee.

---

## DSA Concepts Used

- **Linked List** — all sorting happens via pointer re-linking, never value copying (in approaches 2–3) → see [`/dsa/linked_list.md`](/dsa/linked_list.md)
- **Divide and Conquer** — split, sort halves, combine: the merge sort skeleton → see [`/dsa/divide_and_conquer.md`](/dsa/divide_and_conquer.md)
- **Two Pointers (slow/fast)** — find the middle of a list in one pass → see [`/dsa/two_pointers.md`](/dsa/two_pointers.md)
- **Sorting (merge sort)** — the list-friendly O(n log n) sort; merge step reused from Merge Two Sorted Lists → see [`/dsa/sorting.md`](/dsa/sorting.md)

---

## Approaches Overview

| # | Approach | Time | Space | When to use |
|---|----------|------|-------|-------------|
| 1 | Brute Force (Copy Values + Array Sort) | O(n log n) | O(n) | Quick baseline; fails the O(1)-space follow-up |
| 2 | Merge Sort Top-Down (Recursive) | O(n log n) | O(log n) stack | The standard interview answer; clean and readable |
| 3 | Merge Sort Bottom-Up (Optimal) | O(n log n) | O(1) | The follow-up answer; no recursion at all |

---

## Approach 1 — Brute Force (Copy Values + Array Sort)

### Intuition
Sorting is easy in arrays and awkward in lists, so convert: copy the node values into a slice, sort it with the standard library, and write the values back over the nodes in order. The list's structure never changes — only the values move. This meets the O(n log n) time bar but spends O(n) extra memory, which is exactly what the follow-up forbids.

### Algorithm
1. Traverse the list, appending each `Val` to a slice `vals`.
2. `sort.Ints(vals)`.
3. Traverse the list a second time, assigning `n.Val = vals[i]` in order.
4. Return `head`.

### Complexity
- **Time:** O(n log n) — two O(n) passes plus the dominant array sort.
- **Space:** O(n) — the values slice.

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

| Step | Action           | `vals`        | List        |
|------|------------------|---------------|-------------|
| 1    | collect values   | `[4,2,1,3]`   | `4→2→1→3`   |
| 2    | `sort.Ints`      | `[1,2,3,4]`   | `4→2→1→3`   |
| 3    | write back i=0..3| `[1,2,3,4]`   | `1→2→3→4`   |

Output: `[1,2,3,4]` ✓

---

## Approach 2 — Merge Sort Top-Down (Recursive)

### Intuition
Merge sort is the canonical linked-list sort for two reasons:

1. **No random access needed.** Quicksort and heapsort want O(1) indexing; merge sort only ever walks forward.
2. **Merging lists is free.** Merging two sorted arrays needs an O(n) buffer; merging two sorted *lists* is pure pointer re-linking — zero extra memory.

So: find the middle with slow/fast pointers, **cut** the list into two independent halves, recursively sort each, and merge (exactly LeetCode #21). The only subtlety is starting `fast` at `head.Next` so that `slow` lands on the *last node of the first half* — letting us cut with `slow.Next = nil` and avoiding infinite recursion on 2-node lists.

### Algorithm
1. Base case: `head == nil || head.Next == nil` → return `head`.
2. `slow, fast := head, head.Next`; while `fast != nil && fast.Next != nil`: `slow = slow.Next`, `fast = fast.Next.Next`.
3. `mid := slow.Next`; `slow.Next = nil` — the cut.
4. `left := mergeSortTopDown(head)`, `right := mergeSortTopDown(mid)`.
5. Return `merge(left, right)`: dummy head; repeatedly attach the smaller front node (`<=` for stability); attach the leftover run.

### Complexity
- **Time:** O(n log n) — log n split levels; each level does O(n) total merge work.
- **Space:** O(log n) — recursion stack depth; the merges themselves re-link in place.

### Code
```go
func mergeSortTopDown(head *ListNode) *ListNode {
	if head == nil || head.Next == nil {
		return head
	}

	slow, fast := head, head.Next // fast one ahead → slow ends the first half
	for fast != nil && fast.Next != nil {
		slow = slow.Next
		fast = fast.Next.Next
	}
	mid := slow.Next
	slow.Next = nil // cut into two halves

	left := mergeSortTopDown(head)
	right := mergeSortTopDown(mid)
	return merge(left, right)
}

func merge(a, b *ListNode) *ListNode {
	dummy := &ListNode{}
	tail := dummy
	for a != nil && b != nil {
		if a.Val <= b.Val { // <= keeps it stable
			tail.Next = a
			a = a.Next
		} else {
			tail.Next = b
			b = b.Next
		}
		tail = tail.Next
	}
	if a != nil {
		tail.Next = a
	} else {
		tail.Next = b
	}
	return dummy.Next
}
```

### Dry Run
Example 1: `head = [4,2,1,3]`.

| Step | Call                     | slow/fast walk                          | Split into        | Returns          |
|------|--------------------------|------------------------------------------|-------------------|------------------|
| 1    | `sort(4→2→1→3)`          | slow: 4→2, fast: 2→3 → stop              | `4→2` and `1→3`   | merge of steps 2,5 |
| 2    | `sort(4→2)`              | slow: 4, fast: 2 → stop immediately      | `4` and `2`       | merge of steps 3,4 |
| 3    | `sort(4)`                | base case                                | —                 | `4`              |
| 4    | `sort(2)`                | base case                                | —                 | `2`              |
| 5    | `merge(4, 2)`            | 2 ≤ 4 → take 2; then take 4              | —                 | `2→4`            |
| 6    | `sort(1→3)`              | split into `1` and `3` (as steps 2–5)    | `1` and `3`       | `1→3`            |
| 7    | `merge(2→4, 1→3)`        | take 1 (1<2), take 2 (2≤3), take 3 (3<4), append 4 | —      | `1→2→3→4`        |

Output: `[1,2,3,4]` ✓

Merge table for step 7 in detail:

| Iteration | a    | b    | Compare      | Result list so far |
|-----------|------|------|--------------|---------------------|
| 1         | 2→4  | 1→3  | 2 > 1 → take b(1) | `1`            |
| 2         | 2→4  | 3    | 2 ≤ 3 → take a(2) | `1→2`          |
| 3         | 4    | 3    | 4 > 3 → take b(3) | `1→2→3`        |
| 4         | 4    | nil  | b empty → append a | `1→2→3→4`     |

---

## Approach 3 — Merge Sort Bottom-Up (Optimal, O(1) Space)

### Intuition
The recursion in Approach 2 only exists to *discover* the run boundaries. We can compute them arithmetically instead: at the bottom of merge sort's recursion tree, every single node is a sorted run of length 1. Merge adjacent runs pairwise → sorted runs of length 2. Repeat with width 4, 8, … After ⌈log₂ n⌉ passes the whole list is one run. No call stack, so space drops to O(1) — this is the answer to the follow-up.

Two helpers do the surgery:
- `split(head, n)` — cuts off the first `n` nodes and returns the head of the remainder.
- `mergeTail(a, b, tail)` — merges runs `a` and `b` directly onto the end of the output being rebuilt this pass, returning the new tail.

### Algorithm
1. Count the length `n` of the list.
2. For `width = 1; width < n; width *= 2` (one pass per width):
   1. `tail = dummy`, `curr = dummy.Next`.
   2. While `curr != nil`:
      - `left = curr`; `right = split(left, width)` — cut run A;
      - `curr = split(right, width)` — cut run B, remember the remainder;
      - `tail = mergeTail(left, right, tail)` — merge and append.
3. Return `dummy.Next`.

### Complexity
- **Time:** O(n log n) — ⌈log₂ n⌉ passes, each pass splits+merges all n nodes once.
- **Space:** O(1) — a fixed handful of pointers; no recursion, no buffers.

### Code
```go
func mergeSortBottomUp(head *ListNode) *ListNode {
	n := 0
	for node := head; node != nil; node = node.Next {
		n++
	}

	dummy := &ListNode{Next: head}
	for width := 1; width < n; width *= 2 {
		tail := dummy
		curr := dummy.Next
		for curr != nil {
			left := curr
			right := split(left, width)  // cut run A
			curr = split(right, width)   // cut run B; curr = remainder
			tail = mergeTail(left, right, tail)
		}
	}
	return dummy.Next
}

func split(head *ListNode, n int) *ListNode {
	for i := 1; head != nil && i < n; i++ {
		head = head.Next
	}
	if head == nil {
		return nil
	}
	rest := head.Next
	head.Next = nil // terminate the first run
	return rest
}

func mergeTail(a, b *ListNode, tail *ListNode) *ListNode {
	curr := tail
	for a != nil && b != nil {
		if a.Val <= b.Val {
			curr.Next = a
			a = a.Next
		} else {
			curr.Next = b
			b = b.Next
		}
		curr = curr.Next
	}
	if a != nil {
		curr.Next = a
	} else {
		curr.Next = b
	}
	for curr.Next != nil {
		curr = curr.Next // advance to the merged run's end
	}
	return curr
}
```

### Dry Run
Example 1: `head = [4,2,1,3]`, so `n = 4`.

**Pass 1 — width 1** (merge single nodes pairwise):

| Step | `left` | `right` | `curr` after splits | Merge result appended | List rebuilt so far |
|------|--------|---------|----------------------|------------------------|----------------------|
| 1    | `4`    | `2`     | `1→3`                | `2→4`                  | `2→4`                |
| 2    | `1`    | `3`     | `nil`                | `1→3`                  | `2→4→1→3`            |

**Pass 2 — width 2** (merge runs of two):

| Step | `left` | `right` | `curr` after splits | Merge result appended | List rebuilt so far |
|------|--------|---------|----------------------|------------------------|----------------------|
| 1    | `2→4`  | `1→3`   | `nil`                | `1→2→3→4`              | `1→2→3→4`            |

`width` doubles to 4, `4 < 4` is false → loop ends.
Output: `[1,2,3,4]` ✓

---

## Key Takeaways

- **Merge sort is the linked-list sort**: sequential access only, and the merge step re-links nodes with zero buffer. Quicksort/heapsort lose their appeal without random access.
- **slow/fast with `fast = head.Next`** makes `slow` stop at the end of the first half — the right spot to cut. Starting both at `head` overshoots and infinite-loops on 2-node lists.
- **Cut before you recurse** (`slow.Next = nil`): forgetting the cut is the #1 bug in list merge sort.
- **Top-down vs bottom-up** is a space trade: O(log n) stack (elegant, easy to write) vs O(1) iterative (the follow-up answer). Same time complexity, same merge primitive.
- The `merge` helper is literally LeetCode #21 — sorted-list problems compose.
- Bottom-up's `split`/`mergeTail` pattern ("cut fixed-width runs, merge onto a tail") reappears in external sorting and in real runtimes (e.g. list sort in the Linux kernel is bottom-up merge sort).

---

## Related Problems

- LeetCode #21 — Merge Two Sorted Lists (the merge primitive used here)
- LeetCode #23 — Merge k Sorted Lists (divide and conquer over k lists)
- LeetCode #147 — Insertion Sort List (the O(n²) predecessor to this problem)
- LeetCode #876 — Middle of the Linked List (the slow/fast split step in isolation)
- LeetCode #912 — Sort an Array (same divide-and-conquer, array flavour)
