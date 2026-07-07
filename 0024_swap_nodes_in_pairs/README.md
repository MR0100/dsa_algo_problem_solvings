# 0024 — Swap Nodes in Pairs

> LeetCode #24 · Difficulty: Medium
> **Categories:** Linked List, Recursion

---

## Problem Statement

Given a linked list, swap every two adjacent nodes and return its head. You must solve the problem without modifying the values in the list's nodes (i.e., only nodes themselves may be changed.)

**Example 1**
```
Input:  head = [1,2,3,4]
Output: [2,1,4,3]
```

**Example 2**
```
Input:  head = []
Output: []
```

**Example 3**
```
Input:  head = [1]
Output: [1]
```

**Constraints**
- The number of nodes in the list is in the range `[0, 100]`.
- `0 <= Node.val <= 100`

---

## Company Frequency

| Company   | Frequency       | Last Reported |
|-----------|-----------------|---------------|
| Amazon    | ★★★★☆ High      | 2024          |
| Google    | ★★★★☆ High      | 2024          |
| Microsoft | ★★★☆☆ Medium    | 2024          |
| Bloomberg | ★★★☆☆ Medium    | 2023          |
| Meta      | ★★★☆☆ Medium    | 2023          |
| Apple     | ★★☆☆☆ Low       | 2022          |
| Adobe     | ★★☆☆☆ Low       | 2022          |

> ⚠️ Frequency data is crowd-sourced from LeetCode Discuss, Glassdoor, and community interview reports. Treat as a signal, not a guarantee.

---

## DSA Concepts Used

- **Linked List** — the problem requires relinking nodes (not swapping values), so pointer manipulation is the core skill.
- **Dummy Head Node** — a sentinel node before the list simplifies the swap of the very first pair.
- **Recursion** — Approach 2 reduces the problem to: swap the first pair, then recurse on the remaining list.

---

## Approaches Overview

| # | Approach | Time | Space | When to use |
|---|----------|------|-------|-------------|
| 1 | Iterative ✅ | O(n) | O(1) | Preferred; constant stack space |
| 2 | Recursive | O(n) | O(n) | Elegant; easy to reason about |

---

## Approach 1 — Iterative with Dummy Head (Recommended ✅)

### Intuition
Use a `prev` pointer that sits just before the current pair. For each pair `(first, second)`:
1. `prev.Next = second` — link prev to the new head of the pair.
2. `first.Next = second.Next` — first skips over second.
3. `second.Next = first` — second points back to first.
4. `prev = first` — first is now the pair's tail; advance prev.

### Algorithm
```
dummy → 1 → 2 → 3 → 4
prev = dummy

Pair (1, 2):
  dummy → 2 → 1 → 3 → 4
  prev = node(1)

Pair (3, 4):
  ... → 1 → 4 → 3
  prev = node(3)

Result: dummy → 2 → 1 → 4 → 3
```

### Complexity
- **Time:** O(n) — one pass; n/2 swaps.
- **Space:** O(1) — only pointer variables.

### Code
```go
func iterative(head *ListNode) *ListNode {
    dummy := &ListNode{Next: head}
    prev := dummy
    for prev.Next != nil && prev.Next.Next != nil {
        first, second := prev.Next, prev.Next.Next
        prev.Next = second
        first.Next = second.Next
        second.Next = first
        prev = first
    }
    return dummy.Next
}
```

### Dry Run — `head = [1,2,3,4]`
```
Initial:  dummy→1→2→3→4,  prev=dummy

Iteration 1: first=1, second=2
  dummy→2, 1→3, 2→1   →  dummy→2→1→3→4
  prev=1

Iteration 2: first=3, second=4
  1→4, 3→nil, 4→3     →  dummy→2→1→4→3
  prev=3

prev.Next=nil → stop.
Return: 2→1→4→3 ✓
```

---

## Approach 2 — Recursive

### Intuition
If the list has fewer than 2 nodes, return it unchanged (base case). Otherwise:
1. `second = head.Next`
2. `head.Next = recursive(second.Next)` — head's next is the swapped tail.
3. `second.Next = head` — second precedes head.
4. Return `second`.

### Complexity
- **Time:** O(n) — n/2 recursive calls.
- **Space:** O(n) — call stack depth n/2.

### Code
```go
func recursive(head *ListNode) *ListNode {
    if head == nil || head.Next == nil { return head }
    second := head.Next
    head.Next = recursive(second.Next)
    second.Next = head
    return second
}
```

### Dry Run — `head = [1,2,3,4]`
```
recursive(1):
  second=2, head.Next=recursive(3)
    recursive(3):
      second=4, head.Next=recursive(nil)=nil
      4.Next=3; return 4
  head.Next=4→3; 2.Next=1; return 2

Result: 2→1→4→3 ✓
```

---

## Key Takeaways

- **Three pointer moves per pair** — `prev→second`, `first→second.Next`, `second→first`. Miss one and the list is corrupted or you lose nodes.
- **Dummy head makes the first pair identical to all others** — without it, the first pair needs a special case to update `head`.
- **`prev = first` after swap** — after swapping, `first` is the tail of the pair (the right node). The next pair starts at `first.Next`, which is why `prev` should point to `first`.
- **Iterative preferred** — for 100 nodes, recursion depth is 50 (fine here), but as a habit prefer iterative for linked list problems to avoid stack issues.

---

## Related Problems

- LeetCode #25 — Reverse Nodes in k-Group (generalisation: reverse k nodes at a time)
- LeetCode #206 — Reverse Linked List (reverse the entire list)
- LeetCode #92 — Reverse Linked List II (reverse a subrange)
