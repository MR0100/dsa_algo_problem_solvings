# 0083 — Remove Duplicates from Sorted List

> LeetCode #83 · Difficulty: Easy
> **Categories:** Linked List

---

## Problem Statement

Given the `head` of a sorted linked list, delete all duplicates such that each element appears only once. Return the linked list **sorted** as well.

**Example 1:**
```
Input: head = [1,1,2]
Output: [1,2]
```

**Example 2:**
```
Input: head = [1,1,2,3,3]
Output: [1,2,3]
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
| Microsoft | ★★★☆☆ Medium   | 2023          |
| Google    | ★★☆☆☆ Low      | 2023          |

> ⚠️ Frequency data is crowd-sourced from LeetCode Discuss, Glassdoor, and community interview reports. Treat as a signal, not a guarantee.

---

## DSA Concepts Used

- **Linked List Traversal** — in-place pointer manipulation. See [`/dsa/linked_list.md`](/dsa/linked_list.md)

---

## Approaches Overview

| # | Approach | Time | Space | When to use |
|---|----------|------|-------|-------------|
| 1 | Iterative | O(n) | O(1) | Standard; simplest |
| 2 | Recursive | O(n) | O(n) | Alternative; uses call stack |

---

## Approach 1 — Iterative

### Intuition
Walk through the list with one pointer `cur`. If `cur.Val == cur.Next.Val`, skip `cur.Next` (it's a duplicate). Otherwise advance `cur`. The head is never removed (we always keep one copy), so no dummy head is needed.

### Algorithm
1. `cur = head`.
2. While `cur != nil && cur.Next != nil`:
   - If `cur.Val == cur.Next.Val`: `cur.Next = cur.Next.Next`.
   - Else: `cur = cur.Next`.
3. Return `head`.

### Complexity
- **Time:** O(n)
- **Space:** O(1)

### Code
```go
func deleteDuplicates(head *ListNode) *ListNode {
    cur := head
    for cur != nil && cur.Next != nil {
        if cur.Val == cur.Next.Val {
            cur.Next = cur.Next.Next
        } else {
            cur = cur.Next
        }
    }
    return head
}
```

### Dry Run (head=[1,1,2,3,3])

| cur | cur.Val | cur.Next.Val | action |
|-----|---------|--------------|--------|
| n1(1) | 1 | 1 | skip: cur.Next = n3(2) |
| n1(1) | 1 | 2 | advance: cur = n3(2) |
| n3(2) | 2 | 3 | advance: cur = n4(3) |
| n4(3) | 3 | 3 | skip: cur.Next = nil |
| n4(3) | 3 | nil | loop ends |

Result: `1→2→3` ✓

---

## Approach 2 — Recursive

### Intuition
Recursively deduplicate the rest of the list first. Then, if the current node has the same value as its (already-deduped) next, skip the next node.

### Algorithm
1. If `head == nil || head.Next == nil`: return `head`.
2. `head.Next = deleteDuplicatesRecursive(head.Next)`.
3. If `head.Val == head.Next.Val`: return `head.Next`.
4. Return `head`.

### Complexity
- **Time:** O(n)
- **Space:** O(n) — call stack.

### Code
```go
func deleteDuplicatesRecursive(head *ListNode) *ListNode {
    if head == nil || head.Next == nil {
        return head
    }
    head.Next = deleteDuplicatesRecursive(head.Next)
    if head.Val == head.Next.Val {
        return head.Next
    }
    return head
}
```

### Dry Run (head=[1,1,2])
```
rec([1,1,2]):
  head.Next = rec([1,2]):
    head.Next = rec([2]) = [2]
    1 != 2 → return [1,2]
  head.Next = [1,2]
  head.Val(1) == head.Next.Val(1) → return head.Next = [1,2]... 
  wait, head.Next is now [1,2] so head.Next.Val = 1 → 1==1 → return [1,2]
```
Wait — after `head.Next = rec([1,2])` we get `[1,2]`. Now `head.Val=1 == head.Next.Val=1` → skip head → return `[1,2]` which is `1→2`. ✓

---

## Key Takeaways
- This problem is simpler than #82 because we keep one copy and the head is never removed.
- The key is: only advance `cur` when the next node has a different value; when skipping, stay at `cur` (it may need to skip multiple duplicates).
- Contrast with #82 (delete ALL copies including the first).

---

## Related Problems
- LeetCode #82 — Remove Duplicates from Sorted List II (delete all copies)
- LeetCode #26 — Remove Duplicates from Sorted Array (same idea, array version)
