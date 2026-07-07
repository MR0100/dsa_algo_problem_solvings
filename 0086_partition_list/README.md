# 0086 — Partition List

> LeetCode #86 · Difficulty: Medium
> **Categories:** Linked List, Two Pointers

---

## Problem Statement

Given the `head` of a linked list and a value `x`, partition it such that all nodes **less than** `x` come before nodes **greater than or equal to** `x`.

You should **preserve the original relative order** of the nodes in each of the two partitions.

**Example 1:**
```
Input: head = [1,4,3,2,5,2], x = 3
Output: [1,2,2,4,3,5]
```

**Example 2:**
```
Input: head = [2,1], x = 2
Output: [1,2]
```

**Constraints:**
- The number of nodes in the list is in the range `[0, 200]`.
- `-100 <= Node.val <= 100`
- `-200 <= x <= 200`

---

## Company Frequency

| Company   | Frequency      | Last Reported |
|-----------|----------------|---------------|
| Amazon    | ★★★☆☆ Medium   | 2024          |
| Facebook  | ★★☆☆☆ Low      | 2023          |
| Microsoft | ★★☆☆☆ Low      | 2023          |

> ⚠️ Frequency data is crowd-sourced from LeetCode Discuss, Glassdoor, and community interview reports. Treat as a signal, not a guarantee.

---

## DSA Concepts Used

- **Linked List — Dummy Heads** — two dummy heads create two independent sublists that can be joined. See [`/dsa/linked_list.md`](/dsa/linked_list.md)
- **Two Pointers** — `less` and `greater` are tail pointers for the two sublists.

---

## Approaches Overview

| # | Approach | Time | Space | When to use |
|---|----------|------|-------|-------------|
| 1 | Two Dummy Lists | O(n) | O(1) | Only practical approach; clean and O(n) |

---

## Approach 1 — Two Dummy Lists

### Intuition
Maintain two sublists:
- **less**: nodes with `Val < x`.
- **greater**: nodes with `Val >= x`.

Walk through the original list and append each node to the appropriate sublist. At the end, join `less.tail → greater.head`. Terminate `greater.tail.Next = nil` to avoid cycles (the last node may have a stale `.Next` pointer from the original list).

Two dummy head nodes avoid special-casing the insertion of the first element into each sublist.

### Algorithm
1. `lessHead`, `greaterHead` = new dummy nodes.
2. `less = lessHead`, `greater = greaterHead`.
3. For each node `head`:
   - If `head.Val < x`: `less.Next = head; less = less.Next`.
   - Else: `greater.Next = head; greater = greater.Next`.
   - `head = head.Next`.
4. `greater.Next = nil` — terminate the "greater" list.
5. `less.Next = greaterHead.Next` — connect the two lists.
6. Return `lessHead.Next`.

### Complexity
- **Time:** O(n) — single pass.
- **Space:** O(1) — only pointer reassignment, no new nodes allocated.

### Code
```go
func partition(head *ListNode, x int) *ListNode {
    lessHead := &ListNode{}
    greaterHead := &ListNode{}
    less := lessHead
    greater := greaterHead

    for head != nil {
        if head.Val < x {
            less.Next = head
            less = less.Next
        } else {
            greater.Next = head
            greater = greater.Next
        }
        head = head.Next
    }
    greater.Next = nil
    less.Next = greaterHead.Next
    return lessHead.Next
}
```

### Dry Run (head=[1,4,3,2,5,2], x=3)

| node | Val < 3? | less list | greater list |
|------|----------|-----------|--------------|
| 1 | yes | 1 | — |
| 4 | no | 1 | 4 |
| 3 | no | 1 | 4→3 |
| 2 | yes | 1→2 | 4→3 |
| 5 | no | 1→2 | 4→3→5 |
| 2 | yes | 1→2→2 | 4→3→5 |

Join: `1→2→2→4→3→5` ✓

---

## Key Takeaways
- Two dummy heads eliminate the "first element" edge case for both sublists.
- **Must** set `greater.Next = nil` before joining — otherwise the last "greater" node retains its old `.Next` pointer, creating a cycle.
- Preserving relative order comes naturally because we append in forward traversal order.
- The same two-list technique applies to sorting linked lists by any predicate.

---

## Related Problems
- LeetCode #82 — Remove Duplicates from Sorted List II (dummy head pattern)
- LeetCode #21 — Merge Two Sorted Lists (merging two lists)
- LeetCode #328 — Odd Even Linked List (same two-list separation pattern)
