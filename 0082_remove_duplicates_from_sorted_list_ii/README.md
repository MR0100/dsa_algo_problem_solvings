# 0082 — Remove Duplicates from Sorted List II

> LeetCode #82 · Difficulty: Medium
> **Categories:** Linked List, Two Pointers

---

## Problem Statement

Given the `head` of a sorted linked list, delete all nodes that have duplicate numbers, leaving only distinct numbers from the original list. Return the linked list **sorted** as well.

**Example 1:**
```
Input: head = [1,2,3,3,4,4,5]
Output: [1,2,5]
```

**Example 2:**
```
Input: head = [1,1,1,2,3]
Output: [2,3]
```

**Constraints:**
- The number of nodes in the list is in the range `[0, 300]`.
- `-100 <= Node.val <= 100`
- The list is guaranteed to be **sorted** in ascending order.

---

## Company Frequency

| Company   | Frequency      | Last Reported |
|-----------|----------------|---------------|
| Amazon    | ★★★☆☆ Medium   | 2024          |
| Microsoft | ★★★☆☆ Medium   | 2024          |
| Facebook  | ★★☆☆☆ Low      | 2023          |

> ⚠️ Frequency data is crowd-sourced from LeetCode Discuss, Glassdoor, and community interview reports. Treat as a signal, not a guarantee.

---

## DSA Concepts Used

- **Dummy Head Node** — simplifies edge cases where the head itself might be deleted. See [`/dsa/linked_list.md`](/dsa/linked_list.md)
- **Two Pointers (prev/cur)** — `prev` is the last confirmed distinct node; `cur` scans forward.

---

## Approaches Overview

| # | Approach | Time | Space | When to use |
|---|----------|------|-------|-------------|
| 1 | Iterative (dummy + prev pointer) | O(n) | O(1) | Standard interview answer |
| 2 | Recursive | O(n) | O(n) | Elegant but uses call stack |

---

## Approach 1 — Iterative (Dummy Head)

### Intuition
Use a dummy head so the result list's first node can be removed without special casing. Walk with `prev` (the last confirmed distinct node). When `prev.Next` starts a duplicate run (i.e., `prev.Next.Next.Val == prev.Next.Val`), skip the *entire* run of that value by advancing `prev.Next` past all nodes with that value. Otherwise, advance `prev` normally.

### Algorithm
1. `dummy.Next = head; prev = dummy`.
2. While `prev.Next != nil`:
   - `cur = prev.Next`.
   - If `cur.Next != nil && cur.Next.Val == cur.Val`:
     - `val = cur.Val`.
     - While `prev.Next != nil && prev.Next.Val == val`: `prev.Next = prev.Next.Next`.
   - Else: `prev = prev.Next`.
3. Return `dummy.Next`.

### Complexity
- **Time:** O(n)
- **Space:** O(1)

### Code
```go
func deleteDuplicates(head *ListNode) *ListNode {
    dummy := &ListNode{Next: head}
    prev := dummy
    for prev.Next != nil {
        cur := prev.Next
        if cur.Next != nil && cur.Next.Val == cur.Val {
            val := cur.Val
            for prev.Next != nil && prev.Next.Val == val {
                prev.Next = prev.Next.Next
            }
        } else {
            prev = prev.Next
        }
    }
    return dummy.Next
}
```

### Dry Run (head=[1,2,3,3,4,4,5])

| Step | prev.Val | cur.Val | Action |
|------|----------|---------|--------|
| 1 | dummy | 1 | 1.Next=2≠1 → advance prev to 1 |
| 2 | 1 | 2 | 2.Next=3≠2 → advance prev to 2 |
| 3 | 2 | 3 | 3.Next=3 → skip all 3s: prev.Next=4 |
| 4 | 2 | 4 | 4.Next=4 → skip all 4s: prev.Next=5 |
| 5 | 2 | 5 | 5.Next=nil → advance prev to 5 |
| 6 | 5 | nil | done |

Result: `dummy→1→2→5` ✓

---

## Approach 2 — Recursive

### Intuition
At each call: if the current node begins a duplicate run, skip the entire run and recurse on the remainder. Otherwise, keep the current node and recurse on `head.Next`.

### Algorithm
1. If `head == nil`: return `nil`.
2. If `head.Next != nil && head.Next.Val == head.Val`:
   - `val = head.Val`.
   - Skip all nodes with this value.
   - Return `deleteDuplicatesRecursive(head)` — the first node after the run.
3. Else: `head.Next = deleteDuplicatesRecursive(head.Next); return head`.

### Complexity
- **Time:** O(n)
- **Space:** O(n) — call stack depth proportional to list length.

### Code
```go
func deleteDuplicatesRecursive(head *ListNode) *ListNode {
    if head == nil {
        return nil
    }
    if head.Next != nil && head.Next.Val == head.Val {
        val := head.Val
        for head != nil && head.Val == val {
            head = head.Next
        }
        return deleteDuplicatesRecursive(head)
    }
    head.Next = deleteDuplicatesRecursive(head.Next)
    return head
}
```

### Dry Run (head=[1,1,1,2,3])

```
rec([1,1,1,2,3]) → head.Next.Val==head.Val → skip all 1s → rec([2,3])
  rec([2,3]) → 3≠2 → head.Next = rec([3]) → ...
    rec([3]) → Next=nil → head.Next = rec(nil) = nil → return 3
  → return 2→3
→ return 2→3
```

---

## Key Takeaways
- The dummy head is essential here (unlike #83) because the head node itself might be deleted.
- The "skip entire run" pattern: save the value, then loop `while prev.Next.Val == val`.
- Contrast with #83 (keep one copy) vs #82 (delete all copies).

---

## Related Problems
- LeetCode #83 — Remove Duplicates from Sorted List (keep one copy)
- LeetCode #26 — Remove Duplicates from Sorted Array (array version, keep one)
- LeetCode #80 — Remove Duplicates from Sorted Array II (keep two copies)
