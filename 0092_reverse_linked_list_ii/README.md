# 0092 — Reverse Linked List II

> LeetCode #92 · Difficulty: Medium
> **Categories:** Linked List

---

## Problem Statement

Given the `head` of a singly linked list and two integers `left` and `right` where `left <= right`, reverse the nodes of the list from position `left` to position `right`, and return the reversed list.

**Example 1:**
```
Input: head = [1,2,3,4,5], left = 2, right = 4
Output: [1,4,3,2,5]
```

**Example 2:**
```
Input: head = [5], left = 1, right = 1
Output: [5]
```

**Constraints:**
- The number of nodes in the list is `n`.
- `1 <= n <= 500`
- `-500 <= Node.val <= 500`
- `1 <= left <= right <= n`

**Follow-up:** Could you do it in one pass?

---

## Company Frequency

| Company   | Frequency      | Last Reported |
|-----------|----------------|---------------|
| Amazon    | ★★★★☆ High     | 2024          |
| Facebook  | ★★★★☆ High     | 2024          |
| Microsoft | ★★★☆☆ Medium   | 2023          |
| Google    | ★★★☆☆ Medium   | 2023          |

> ⚠️ Frequency data is crowd-sourced from LeetCode Discuss, Glassdoor, and community interview reports. Treat as a signal, not a guarantee.

---

## DSA Concepts Used

- **Linked List — In-Place Reversal** — "insert at front" technique avoids needing to re-scan. See [`/dsa/linked_list.md`](/dsa/linked_list.md)
- **Dummy Head Node** — handles the case where `left = 1` (reversing from the head).

---

## Approaches Overview

| # | Approach | Time | Space | When to use |
|---|----------|------|-------|-------------|
| 1 | Collect Values + Rebuild | O(n) | O(n) | Simple but uses extra space |
| 2 | One-Pass In-Place | O(n) | O(1) | Optimal — single pass, O(1) space |

---

## Approach 1 — Collect Values + Rebuild

### Intuition
Extract all values into a slice, reverse the sub-slice `[left-1..right-1]`, then write values back into the nodes.

### Complexity
- **Time:** O(n)
- **Space:** O(n)

### Code
```go
func reverseBetweenSimple(head *ListNode, left int, right int) *ListNode {
    vals := []int{}
    for cur := head; cur != nil; cur = cur.Next { vals = append(vals, cur.Val) }
    l, r := left-1, right-1
    for l < r { vals[l], vals[r] = vals[r], vals[l]; l++; r-- }
    cur := head
    for i := range vals { cur.Val = vals[i]; cur = cur.Next }
    return head
}
```

### Dry Run (head=[1,2,3,4,5], left=2, right=4)
vals=[1,2,3,4,5]. Reverse indices 1..3: vals=[1,4,3,2,5]. Write back: [1,4,3,2,5] ✓

---

## Approach 2 — One-Pass In-Place ("Insert at Front")

### Intuition
After advancing `prev` to position `left-1` and setting `curr` to the first node of the reversal range, repeatedly move the node *after* `curr` to the front of the reversed segment (just after `prev`).

In each iteration, node `next = curr.Next` is detached from its current position and inserted right after `prev`. After `right - left` iterations, the segment is fully reversed.

**No need to move `curr`** — it naturally ends up at the tail of the reversed segment.

### Algorithm
1. `dummy.Next = head; prev = dummy`.
2. Advance `prev` to position `left-1`.
3. `curr = prev.Next`.
4. Repeat `right-left` times:
   - `next = curr.Next`.
   - `curr.Next = next.Next`.
   - `next.Next = prev.Next`.
   - `prev.Next = next`.

### Complexity
- **Time:** O(n)
- **Space:** O(1)

### Code
```go
func reverseBetween(head *ListNode, left int, right int) *ListNode {
    dummy := &ListNode{Next: head}
    prev := dummy
    for i := 1; i < left; i++ { prev = prev.Next }
    curr := prev.Next
    for i := 0; i < right-left; i++ {
        next := curr.Next
        curr.Next = next.Next
        next.Next = prev.Next
        prev.Next = next
    }
    return dummy.Next
}
```

### Dry Run (head=[1,2,3,4,5], left=2, right=4)

Initial: `dummy→1→2→3→4→5`, `prev=node(1)`, `curr=node(2)`.

**Iteration 1** (move node 3 to front):
- `next=node(3)`.
- `curr.Next = node(4)` → `2→4`.
- `next.Next = node(2)` → `3→2`.
- `prev.Next = node(3)` → `1→3`.
- List: `dummy→1→3→2→4→5`.

**Iteration 2** (move node 4 to front):
- `next=node(4)`.
- `curr.Next = node(5)` → `2→5`.
- `next.Next = node(3)` → `4→3`.
- `prev.Next = node(4)` → `1→4`.
- List: `dummy→1→4→3→2→5`.

Result: `[1,4,3,2,5]` ✓

---

## Key Takeaways
- The "insert at front" technique: take the node just after `curr` and move it to just after `prev`. `curr` doesn't move — it's the tail anchor.
- Dummy head is essential when `left = 1` (head itself might be the new tail).
- After `right-left` iterations, the segment `[left..right]` is reversed.
- This exact technique appears in other in-place reversal problems (e.g., reverse in groups of k).

---

## Related Problems
- LeetCode #206 — Reverse Linked List (entire list; same idea without boundaries)
- LeetCode #25 — Reverse Nodes in k-Group (reverse multiple segments)
- LeetCode #24 — Swap Nodes in Pairs (reverse pairs; special case of k=2)
