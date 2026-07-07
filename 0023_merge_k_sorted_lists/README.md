# 0023 — Merge k Sorted Lists

> LeetCode #23 · Difficulty: Hard
> **Categories:** Linked List, Divide and Conquer, Heap (Priority Queue), Merge Sort

---

## Problem Statement

You are given an array of `k` linked-lists `lists`, each linked-list is sorted in ascending order.

*Merge all the linked-lists into one sorted linked-list and return it.*

**Example 1**
```
Input:  lists = [[1,4,5],[1,3,4],[2,6]]
Output: [1,1,2,3,4,4,5,6]
Explanation: The linked-lists are:
  1->4->5, 1->3->4, 2->6
Merging them gives: 1->1->2->3->4->4->5->6
```

**Example 2**
```
Input:  lists = []
Output: []
```

**Example 3**
```
Input:  lists = [[]]
Output: []
```

**Constraints**
- `k == lists.length`
- `0 <= k <= 10⁴`
- `0 <= lists[i].length <= 500`
- `-10⁴ <= lists[i][j] <= 10⁴`
- `lists[i]` is sorted in **ascending** order.
- The sum of `lists[i].length` will not exceed `10⁴`.

---

## Company Frequency

| Company   | Frequency       | Last Reported |
|-----------|-----------------|---------------|
| Amazon    | ★★★★★ Very High | 2024          |
| Google    | ★★★★★ Very High | 2024          |
| Meta      | ★★★★★ Very High | 2024          |
| Microsoft | ★★★★☆ High      | 2024          |
| Bloomberg | ★★★★☆ High      | 2024          |
| Apple     | ★★★★☆ High      | 2023          |
| Adobe     | ★★★☆☆ Medium    | 2023          |
| Uber      | ★★★☆☆ Medium    | 2023          |

> ⚠️ Frequency data is crowd-sourced from LeetCode Discuss, Glassdoor, and community interview reports. Treat as a signal, not a guarantee.

---

## DSA Concepts Used

- **Linked List** — all solutions relink existing nodes; no new nodes are allocated.
- **Divide and Conquer** — Approach 3 pairs lists and merges in O(log k) rounds, making it O(N log k).
- **Min-Heap / Priority Queue** — Approach 4 keeps the k current list heads in a min-heap, popping the global minimum each time. → see [`/dsa/heap.md`](/dsa/heap.md)

---

## Approaches Overview

| # | Approach | Time | Space | When to use |
|---|----------|------|-------|-------------|
| 1 | Brute Force (sort all values) | O(N log N) | O(N) | Quick to write; ignores sorted structure |
| 2 | Sequential Merge | O(kN) | O(1) | Fine for small k; degrades badly for large k |
| 3 | Divide and Conquer ✅ | O(N log k) | O(log k) | Optimal time; uses the merge-two-lists subroutine |
| 4 | Min-Heap ✅ | O(N log k) | O(k) | Optimal time; natural for streaming / online scenarios |

N = total nodes across all lists.

---

## Approach 1 — Brute Force (Collect and Sort)

### Intuition
Collect every node's value into a slice, sort it, then build a new linked list. Ignores the fact that each list is already sorted.

### Complexity
- **Time:** O(N log N) — dominated by the sort.
- **Space:** O(N) — values slice.

---

## Approach 2 — Sequential Merge

### Intuition
Merge `lists[0]` with `lists[1]`, then merge the result with `lists[2]`, etc. Like insertion sort: always work with the accumulated result.

### Why it's suboptimal
If each list has N/k nodes, the i-th merge handles a list of size `i·(N/k)` and a new list of size `N/k`. Total work: Σᵢ₌₁ᵏ `(i+1)·(N/k)` ≈ O(k·N). For large k, this is much worse than O(N log k).

### Complexity
- **Time:** O(kN).
- **Space:** O(1) — in-place relinking.

---

## Approach 3 — Divide and Conquer (Recommended ✅)

### Intuition
Pair up lists and merge each pair simultaneously. After one round, k → k/2 lists. Repeat for O(log k) rounds, each touching all N nodes once.

This is the "combine" step of merge sort applied to k lists rather than k elements.

### Algorithm
```
dnc(lists, lo, hi):
  if lo == hi: return lists[lo]
  mid = (lo+hi)/2
  left  = dnc(lists, lo, mid)
  right = dnc(lists, mid+1, hi)
  return mergeTwoLists(left, right)
```

### Complexity
- **Time:** O(N log k) — O(log k) rounds × O(N) per round.
- **Space:** O(log k) — recursion depth.

### Dry Run — `lists = [L1, L2, L3, L4]` (4 lists)
```
Round 1: merge(L1,L2)→M12, merge(L3,L4)→M34
Round 2: merge(M12,M34)→result
Total: 2 rounds = log₂(4)
Each round touches all N nodes once → 2·N work = O(N log k) ✓
```

---

## Approach 4 — Min-Heap (Priority Queue)

### Intuition
Maintain a min-heap of size k holding one node per list (the current front of each non-exhausted list). At each step, pop the globally minimum node and attach it to the result. Push that node's successor (if it exists) onto the heap.

Each of the N nodes is pushed and popped exactly once, each O(log k).

### Code
```go
func minHeap(lists []*ListNode) *ListNode {
    h := &nodeHeap{}
    heap.Init(h)
    for _, node := range lists { if node != nil { heap.Push(h, node) } }
    dummy := &ListNode{}; cur := dummy
    for h.Len() > 0 {
        node := heap.Pop(h).(*ListNode)
        cur.Next = node; cur = cur.Next
        if node.Next != nil { heap.Push(h, node.Next) }
    }
    return dummy.Next
}
```

### Complexity
- **Time:** O(N log k) — N pop+push operations, each O(log k).
- **Space:** O(k) — the heap holds at most k nodes.

### Dry Run — `lists = [[1,4,5],[1,3,4],[2,6]]`
```
Initial heap: {1(L1), 1(L2), 2(L3)}

Pop 1(L1) → result: 1. Push 4(L1). Heap: {1(L2),2(L3),4(L1)}
Pop 1(L2) → result: 1→1. Push 3(L2). Heap: {2(L3),3(L2),4(L1)}
Pop 2(L3) → result: 1→1→2. Push 6(L3). Heap: {3(L2),4(L1),6(L3)}
Pop 3(L2) → result: 1→1→2→3. Push 4(L2). Heap: {4(L1),4(L2),6(L3)}
Pop 4(L1) → result: 1→1→2→3→4. Push 5(L1). Heap: {4(L2),5(L1),6(L3)}
Pop 4(L2) → result: ...→4. No next. Heap: {5(L1),6(L3)}
Pop 5(L1) → result: ...→5. No next. Heap: {6(L3)}
Pop 6(L3) → result: ...→6. No next. Heap: {}

Final: [1,1,2,3,4,4,5,6] ✓
```

---

## Key Takeaways

- **Sequential vs divide-and-conquer** — sequential merge is O(kN); divide-and-conquer is O(N log k). The log k comes from doing O(log k) rounds instead of k.
- **Heap and divide-and-conquer are both optimal at O(N log k)** — choose divide-and-conquer when code simplicity matters; choose the heap when the lists are streamed online (you don't have all k lists upfront).
- **Go's `container/heap`** — requires implementing `heap.Interface`: `Len`, `Less`, `Swap`, `Push`, `Pop`. The `Less` function determines min-heap vs max-heap.
- **This is the merge step of external merge sort** — when sorting data too large for RAM, you write sorted runs to disk and merge them with a k-way heap merge.

---

## Related Problems

- LeetCode #21 — Merge Two Sorted Lists (the subroutine)
- LeetCode #148 — Sort List (merge sort on linked list)
- LeetCode #264 — Ugly Number II (k-way merge from multiple virtual sequences)
- LeetCode #373 — Find K Pairs with Smallest Sums (min-heap, k-way merge)
