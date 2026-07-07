# 0061 — Rotate List

> LeetCode #61 · Difficulty: Medium
> **Categories:** Linked List, Two Pointers

---

## Problem Statement

Given the `head` of a linked list, rotate the list to the right by `k` places.

**Example 1**
```
Input:  head = [1,2,3,4,5], k = 2
Output: [4,5,1,2,3]
```

**Example 2**
```
Input:  head = [0,1,2], k = 4
Output: [2,0,1]
```

**Constraints**
- The number of nodes in the list is in the range `[0, 500]`.
- `-100 <= Node.val <= 100`
- `0 <= k <= 2 * 10⁹`

---

## Company Frequency

| Company   | Frequency       | Last Reported |
|-----------|-----------------|---------------|
| Amazon    | ★★★★☆ High      | 2024          |
| Google    | ★★★★☆ High      | 2024          |
| Microsoft | ★★★☆☆ Medium    | 2023          |
| Meta      | ★★★☆☆ Medium    | 2023          |
| Bloomberg | ★★★☆☆ Medium    | 2023          |
| Apple     | ★★☆☆☆ Low       | 2022          |

> ⚠️ Frequency data is crowd-sourced from LeetCode Discuss, Glassdoor, and community interview reports. Treat as a signal, not a guarantee.

---

## DSA Concepts Used

- **Linked List Manipulation** — find tail, make circular, find new tail, break link.
- **Modular Arithmetic** — `k % n` normalises large k values.

---

## Approaches Overview

| # | Approach | Time | Space | When to use |
|---|----------|------|-------|-------------|
| 1 | Array Copy | O(n) | O(n) | Simple; but allocates extra space |
| 2 | Find Tail and Reconnect ✅ | O(n) | O(1) | Optimal; pure linked list manipulation |

---

## Approach 1 — Array Copy

### Intuition
Collect all values into a slice, rotate the slice by `k % n` positions (using slice rearrangement), write back into the original nodes.

### Complexity
- **Time:** O(n).
- **Space:** O(n) — extra slice.

---

## Approach 2 — Find Tail and Reconnect (Recommended ✅)

### Intuition
Rotating right by `k` is equivalent to taking the last `k % n` nodes and prepending them. With a linked list, we can do this in O(n) by:
1. Walking to find the length `n` and the tail.
2. Computing effective rotation `k = k % n`. If 0, return unchanged.
3. Making the list circular (tail → head).
4. Walking `n - k - 1` steps from head to reach the **new tail**.
5. New head = new tail's next; break the circle.

### Algorithm
```
find n and tail; tail.Next = head
k = k % n; if k == 0: return head
newTail = head; advance newTail (n-k-1) steps
newHead = newTail.Next; newTail.Next = nil
return newHead
```

### Complexity
- **Time:** O(n) — two passes (length + advance).
- **Space:** O(1).

### Code
```go
func reconnect(head *ListNode, k int) *ListNode {
    n := 1; tail := head
    for tail.Next != nil { tail = tail.Next; n++ }
    k = k % n
    if k == 0 { return head }
    tail.Next = head   // make circular
    newTail := head
    for i := 0; i < n-k-1; i++ { newTail = newTail.Next }
    newHead := newTail.Next; newTail.Next = nil
    return newHead
}
```

### Dry Run — `head = [1,2,3,4,5]`, `k = 2`
```
n=5, tail=node(5), k=2%5=2
Make circular: 5→1
Advance n-k-1 = 2 steps: head→1→2→3 (newTail = node(3))
newHead = node(4)
Break: node(3).Next = nil
Result: 4→5→1→2→3 ✓
```

---

## Key Takeaways

- **`k % n` is essential** — k can be up to 2×10⁹ but the effective rotation is always in [0, n-1].
- **Make circular, then break** — this is a clean pattern for list rotation; no need for multiple passes or complex pointer manipulation.
- **New tail position** — `n - k - 1` from head (0-indexed). The `-1` accounts for 0-indexing: if k=2 and n=5, new tail is at index 2 (node 3), new head is at index 3 (node 4).

---

## Related Problems

- LeetCode #189 — Rotate Array (same concept for arrays; three-reversal trick)
- LeetCode #24 — Swap Nodes in Pairs (linked list pointer manipulation)
- LeetCode #25 — Reverse Nodes in k-Group (more complex linked list ops)
