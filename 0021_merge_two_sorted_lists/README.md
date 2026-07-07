# 0021 — Merge Two Sorted Lists

> LeetCode #21 · Difficulty: Easy
> **Categories:** Linked List, Recursion, Two Pointers

---

## Problem Statement

You are given the heads of two sorted linked lists `list1` and `list2`.

Merge the two lists into one **sorted** list. The list should be made by splicing together the nodes of the first two lists.

Return the head of the merged linked list.

**Example 1**
```
Input:  list1 = [1,2,4], list2 = [1,3,4]
Output: [1,1,2,3,4,4]
```

**Example 2**
```
Input:  list1 = [], list2 = []
Output: []
```

**Example 3**
```
Input:  list1 = [], list2 = [0]
Output: [0]
```

**Constraints**
- The number of nodes in both lists is in the range `[0, 50]`.
- `-100 <= Node.val <= 100`
- Both `list1` and `list2` are sorted in non-decreasing order.

---

## Company Frequency

| Company   | Frequency       | Last Reported |
|-----------|-----------------|---------------|
| Amazon    | ★★★★★ Very High | 2024          |
| Google    | ★★★★★ Very High | 2024          |
| Meta      | ★★★★★ Very High | 2024          |
| Microsoft | ★★★★☆ High      | 2024          |
| Bloomberg | ★★★★☆ High      | 2024          |
| Apple     | ★★★☆☆ Medium    | 2023          |
| Adobe     | ★★★☆☆ Medium    | 2023          |
| Uber      | ★★☆☆☆ Low       | 2022          |

> ⚠️ Frequency data is crowd-sourced from LeetCode Discuss, Glassdoor, and community interview reports. Treat as a signal, not a guarantee.

---

## DSA Concepts Used

- **Linked List** — the problem requires relinking existing nodes without allocating new ones.
- **Two Pointers** — both approaches maintain a pointer into each list and advance the smaller one.
- **Dummy Head Node** — a sentinel `&ListNode{}` before the result simplifies the "attach first node" edge case.

---

## Approaches Overview

| # | Approach | Time | Space | When to use |
|---|----------|------|-------|-------------|
| 1 | Iterative ✅ | O(m+n) | O(1) | Preferred; no stack overflow risk |
| 2 | Recursive | O(m+n) | O(m+n) | Elegant; risks stack overflow for very long lists |

---

## Approach 1 — Iterative with Dummy Head (Recommended ✅)

### Intuition
Use a dummy node as the head of the result list. A `cur` pointer advances through the result, always attaching the smaller of the two current list heads. When one list is exhausted, attach the remainder of the other (it is already sorted, so this is O(1)).

### Algorithm
1. `dummy = &ListNode{}`, `cur = dummy`.
2. While both lists are non-nil:
   - Attach the smaller head to `cur.Next`, advance that list.
   - `cur = cur.Next`.
3. Attach the non-nil remainder.
4. Return `dummy.Next`.

### Complexity
- **Time:** O(m+n) — one comparison per node.
- **Space:** O(1) — no new nodes; only the dummy head is extra.

### Code
```go
func iterative(list1, list2 *ListNode) *ListNode {
    dummy := &ListNode{}
    cur := dummy
    for list1 != nil && list2 != nil {
        if list1.Val <= list2.Val {
            cur.Next = list1; list1 = list1.Next
        } else {
            cur.Next = list2; list2 = list2.Next
        }
        cur = cur.Next
    }
    if list1 != nil { cur.Next = list1 } else { cur.Next = list2 }
    return dummy.Next
}
```

### Dry Run — `list1=[1,2,4]`, `list2=[1,3,4]`
```
cur=dummy
1<=1 → attach list1(1), list1→2
1≤1? no: 1<=1 → attach list1 first (<=) → attach 1, list1→2
  actually: 1(list1) <= 1(list2) → attach list1's 1
cur→1; list1=2
2>1  → attach list2's 1; cur→1; list2=3
2<=3 → attach list1's 2; cur→2; list1=4
4>3  → attach list2's 3; cur→3; list2=4
4<=4 → attach list1's 4; cur→4; list1=nil
list1 nil → attach list2 remainder [4]
Result: 1→1→2→3→4→4 ✓
```

---

## Approach 2 — Recursive

### Intuition
The smaller head becomes the merged list's head. Its `Next` is the merged result of the remaining nodes — a strictly smaller sub-problem.

### Algorithm
```
merge(l1, l2):
  if l1 == nil: return l2
  if l2 == nil: return l1
  if l1.Val <= l2.Val:
    l1.Next = merge(l1.Next, l2)
    return l1
  else:
    l2.Next = merge(l1, l2.Next)
    return l2
```

### Complexity
- **Time:** O(m+n).
- **Space:** O(m+n) — recursion stack depth equals total nodes.

---

## Key Takeaways

- **Dummy head eliminates the "which list starts first" special case** — without it, you'd need to determine the initial head separately before the loop.
- **Attach remainder in O(1)** — since both input lists are already sorted, whichever list has nodes left when the loop exits can be attached directly. No need to traverse it.
- **Iterative > recursive for linked lists** — recursion uses O(n) stack space. For long lists (millions of nodes), iterative is safe; recursive may stack-overflow.
- **This is the merge step of merge sort** — understanding this problem unlocks LeetCode #148 (Sort List).

---

## Related Problems

- LeetCode #23 — Merge k Sorted Lists (extension to k lists)
- LeetCode #148 — Sort List (merge sort on linked list uses this as a subroutine)
- LeetCode #88 — Merge Sorted Array (same idea but on arrays)
