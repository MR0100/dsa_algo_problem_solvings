# 0019 — Remove Nth Node From End of List

> LeetCode #19 · Difficulty: Medium
> **Categories:** Linked List, Two Pointers

---

## Problem Statement

Given the `head` of a linked list, remove the `n`-th node from the end of the list and return its head.

**Example 1**
```
Input:  head = [1,2,3,4,5], n = 2
Output: [1,2,3,5]
```

**Example 2**
```
Input:  head = [1], n = 1
Output: []
```

**Example 3**
```
Input:  head = [1,2], n = 1
Output: [1]
```

**Constraints**
- The number of nodes in the list is `sz`.
- `1 <= sz <= 30`
- `0 <= Node.val <= 100`
- `1 <= n <= sz`

**Follow-up:** Could you do this in one pass?

---

## Company Frequency

| Company   | Frequency       | Last Reported |
|-----------|-----------------|---------------|
| Amazon    | ★★★★★ Very High | 2024          |
| Google    | ★★★★☆ High      | 2024          |
| Microsoft | ★★★★☆ High      | 2024          |
| Bloomberg | ★★★★☆ High      | 2024          |
| Meta      | ★★★☆☆ Medium    | 2023          |
| Apple     | ★★★☆☆ Medium    | 2023          |
| Adobe     | ★★☆☆☆ Low       | 2022          |
| Uber      | ★★☆☆☆ Low       | 2022          |

> ⚠️ Frequency data is crowd-sourced from LeetCode Discuss, Glassdoor, and community interview reports. Treat as a signal, not a guarantee.

---

## DSA Concepts Used

- **Linked List** — the fundamental traversal and unlinking operations.
- **Two Pointers (fast/slow)** — Approach 2 uses a gap of `n+1` nodes between a fast and slow pointer to land slow at the predecessor of the target in one pass. → see [`/dsa/two_pointers.md`](/dsa/two_pointers.md)
- **Dummy Head Node** — a sentinel node before `head` simplifies edge cases where the head itself must be removed.

---

## Approaches Overview

| # | Approach | Time | Space | When to use |
|---|----------|------|-------|-------------|
| 1 | Two-Pass (count length) | O(L) | O(1) | Simple; two traversals |
| 2 | One-Pass (two pointers) ✅ | O(L) | O(1) | Interview follow-up; single traversal |

---

## Approach 1 — Two-Pass (Count Length)

### Intuition
First pass: count the list length `L`. The n-th node from the end is the `(L-n+1)`-th from the start (1-indexed), or equivalently the `(L-n)`-th from a dummy head (0-indexed). Second pass: walk to that predecessor and unlink `predecessor.Next`.

### Algorithm
1. Count `L` (first pass).
2. Walk `dummy` pointer `L-n` steps (second pass).
3. `cur.Next = cur.Next.Next`.

### Complexity
- **Time:** O(L) — two passes.
- **Space:** O(1).

---

## Approach 2 — One-Pass with Two Pointers (Recommended ✅)

### Intuition
Maintain a gap of `n+1` nodes between `fast` and `slow`, both starting at a dummy head. First advance `fast` by `n+1` steps. Then advance both together until `fast` is `nil`. At that point `slow` is at the predecessor of the target.

Why `n+1` (not `n`)? We want `slow` to stop **before** the target (to be able to set `slow.Next = slow.Next.Next`). If `fast` is `n+1` ahead of `slow`, when `fast` hits `nil`, `slow` is `n+1` before the end — which is the predecessor of the n-th from the end.

### Algorithm
1. `dummy.Next = head`. `fast = slow = dummy`.
2. Advance `fast` `n+1` times.
3. While `fast != nil`: advance both `fast` and `slow`.
4. `slow.Next = slow.Next.Next`.
5. Return `dummy.Next`.

### Complexity
- **Time:** O(L) — single traversal.
- **Space:** O(1).

### Code
```go
func onePass(head *ListNode, n int) *ListNode {
    dummy := &ListNode{Next: head}
    fast, slow := dummy, dummy
    for i := 0; i <= n; i++ { fast = fast.Next }
    for fast != nil { fast = fast.Next; slow = slow.Next }
    slow.Next = slow.Next.Next
    return dummy.Next
}
```

### Dry Run — `head = [1,2,3,4,5]`, `n = 2`
```
dummy → 1 → 2 → 3 → 4 → 5 → nil

After advancing fast n+1=3 steps from dummy:
  fast → 3 (dummy, 1, 2, 3)
  slow → dummy

Walk both until fast = nil:
  step 1: fast→4, slow→1
  step 2: fast→5, slow→2
  step 3: fast→nil, slow→3

slow=3, slow.Next=4 (the 2nd from end).
slow.Next = slow.Next.Next = 5.

List: 1→2→3→5 ✓
```

### Dry Run — `head = [1]`, `n = 1`
```
dummy → 1 → nil

Advance fast n+1=2 steps: fast = nil (went past end)
fast is already nil → inner loop doesn't run.
slow = dummy.
slow.Next = slow.Next.Next = nil.
Return dummy.Next = nil → [] ✓
```

---

## Key Takeaways

- **Dummy head** — makes the head-removal case identical to all other cases. Without it, you need a special check for `n == L`.
- **Gap of n+1, not n** — the gap between slow and fast is `n+1` so that slow lands at the **predecessor**, not the target. Off-by-one here is the most common mistake.
- **Fast/slow pointer pattern** — this is the standard "find n-th from end in one pass" technique. The same idea appears in "find middle of list" (gap = L/2) and "detect cycle" (fast moves 2 steps).
- **Guaranteed valid input** — the constraints say `1 ≤ n ≤ sz`, so we never need to handle invalid `n`.

---

## Related Problems

- LeetCode #876 — Middle of the Linked List (fast pointer moves 2x — finds mid in one pass)
- LeetCode #141 — Linked List Cycle (fast/slow pointer cycle detection)
- LeetCode #206 — Reverse Linked List (fundamental list manipulation)
- LeetCode #21 — Merge Two Sorted Lists (linked list merge)
