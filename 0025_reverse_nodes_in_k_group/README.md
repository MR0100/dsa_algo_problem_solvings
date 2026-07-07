# 0025 — Reverse Nodes in k-Group

> LeetCode #25 · Difficulty: Hard
> **Categories:** Linked List, Recursion

---

## Problem Statement

Given the `head` of a linked list, reverse the nodes of the list `k` at a time, and return *the modified list*.

`k` is a positive integer and is less than or equal to the length of the linked list. If the number of nodes is not a multiple of `k` then left-out nodes, in the end, should remain as is.

You may not alter the values in the list's nodes, only the nodes themselves may be changed.

**Example 1**
```
Input:  head = [1,2,3,4,5], k = 2
Output: [2,1,4,3,5]
```

**Example 2**
```
Input:  head = [1,2,3,4,5], k = 3
Output: [3,2,1,4,5]
```

**Constraints**
- The number of nodes in the list is `n`.
- `1 <= k <= n <= 5000`
- `0 <= Node.val <= 1000`

**Follow-up:** Can you solve the problem in `O(1)` extra memory space?

---

## Company Frequency

| Company   | Frequency       | Last Reported |
|-----------|-----------------|---------------|
| Amazon    | ★★★★★ Very High | 2024          |
| Google    | ★★★★★ Very High | 2024          |
| Meta      | ★★★★☆ High      | 2024          |
| Microsoft | ★★★★☆ High      | 2024          |
| Bloomberg | ★★★★☆ High      | 2024          |
| Apple     | ★★★☆☆ Medium    | 2023          |
| Adobe     | ★★★☆☆ Medium    | 2023          |

> ⚠️ Frequency data is crowd-sourced from LeetCode Discuss, Glassdoor, and community interview reports. Treat as a signal, not a guarantee.

---

## DSA Concepts Used

- **Linked List** — requires pointer manipulation to reverse a sublist without allocating new nodes.
- **Dummy Head + `prev` group tail pointer** — the dummy head absorbs the "special first group" edge case; `prevGroupTail` links each newly reversed group to the previous.
- **Recursion** — Approach 2 reverses the first k nodes and recurses for the rest; elegant but uses O(n/k) stack space.

---

## Approaches Overview

| # | Approach | Time | Space | When to use |
|---|----------|------|-------|-------------|
| 1 | Iterative ✅ | O(n) | O(1) | The follow-up answer; O(1) extra space |
| 2 | Recursive | O(n) | O(n/k) | Cleaner to derive; fine for small lists |

---

## Approach 1 — Iterative (Recommended ✅)

### Intuition
Use a dummy head. Maintain `prevGroupTail` — the tail of the most recently completed group (initially `dummy`). For each group:
1. **Check:** use `getKth` to verify k nodes exist; if not, done.
2. **Isolate:** note `groupHead = prevGroupTail.Next` and `nextGroupHead = kthNode.Next`.
3. **Reverse:** standard in-place reversal stopping at `nextGroupHead`.
4. **Relink:** `prevGroupTail.Next = kthNode` (new group head after reversal).
5. **Advance:** `prevGroupTail = groupHead` (original head is now the group tail).

### Algorithm
```
dummy.Next = head; prevGroupTail = dummy
loop:
  kth = getKth(prevGroupTail, k)
  if kth == nil: break
  groupHead = prevGroupTail.Next
  nextGroupHead = kth.Next
  reverse(groupHead → kth), stop before nextGroupHead
  prevGroupTail.Next = kth   // kth is now the group's new head
  prevGroupTail = groupHead  // groupHead is now the group's new tail
```

### Complexity
- **Time:** O(n) — each node reversed once.
- **Space:** O(1) — only pointer variables.

### Code
```go
func iterative(head *ListNode, k int) *ListNode {
    dummy := &ListNode{Next: head}
    prevGroupTail := dummy
    for {
        kthNode := getKth(prevGroupTail, k)
        if kthNode == nil { break }
        groupHead := prevGroupTail.Next
        nextGroupHead := kthNode.Next
        prev, cur := nextGroupHead, groupHead
        for cur != nextGroupHead {
            nxt := cur.Next; cur.Next = prev; prev = cur; cur = nxt
        }
        prevGroupTail.Next = kthNode
        prevGroupTail = groupHead
    }
    return dummy.Next
}
func getKth(node *ListNode, k int) *ListNode {
    for k > 0 && node != nil { node = node.Next; k-- }
    return node
}
```

### Dry Run — `head = [1,2,3,4,5]`, `k = 2`
```
dummy→1→2→3→4→5; prevGroupTail=dummy

Group 1: kth=getKth(dummy,2)=node(2). groupHead=1, nextGroupHead=3
  Reverse 1→2 (stop at 3): 2→1→3
  prevGroupTail.Next=2; prevGroupTail=1
  List: dummy→2→1→3→4→5

Group 2: kth=getKth(node(1),2)=node(4). groupHead=3, nextGroupHead=5
  Reverse 3→4 (stop at 5): 4→3→5
  prevGroupTail.Next=4; prevGroupTail=3
  List: dummy→2→1→4→3→5

Group 3: kth=getKth(node(3),2)=nil. Break.

Return dummy.Next = 2→1→4→3→5 ✓
```

---

## Approach 2 — Recursive

### Intuition
Count k nodes. If fewer than k exist, return head unchanged. Otherwise reverse k nodes in-place using the standard reversal loop (stop after k steps). Connect the original group head (now tail) to `recursive(curr, k)` where `curr` is the first node of the next group.

### Complexity
- **Time:** O(n).
- **Space:** O(n/k) — one stack frame per group.

---

## Key Takeaways

- **`getKth` is the key helper** — checking whether k nodes exist before reversing avoids processing the leftover tail. Writing it cleanly as a separate function makes the main loop simple.
- **After reversal, original `groupHead` is the group's tail** — so `prevGroupTail = groupHead`. This is the only slightly non-obvious line.
- **`nextGroupHead` is the reversal sentinel** — the in-place reversal uses `prev = nextGroupHead` as its starting "previous" value, so when `cur` reaches `nextGroupHead`, the loop stops automatically and the group tail is correctly linked to the next group.
- **Generalises #24** — Swap Nodes in Pairs is exactly this problem with `k = 2`.

---

## Related Problems

- LeetCode #24 — Swap Nodes in Pairs (`k = 2` special case)
- LeetCode #206 — Reverse Linked List (full list reversal, the subroutine)
- LeetCode #92 — Reverse Linked List II (reverse a subrange by indices)
