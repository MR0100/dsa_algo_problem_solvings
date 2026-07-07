# Linked List

> **Category:** Data Structure / Pointer Manipulation
> **Difficulty to master:** Easy to learn, medium to master (pointer bugs are subtle)

---

## What is a linked list?

A **linked list** is a linear data structure where elements (**nodes**) are not
stored in contiguous memory. Each node holds a value and one or more pointers
to neighbouring nodes:

- **Singly linked list** — each node points to the `next` node only.
- **Doubly linked list** — each node points to both `next` and `prev`.
- **Circular linked list** — the tail's `next` points back to the head.

```go
// Singly linked list node (the shape LeetCode uses for almost all problems).
type ListNode struct {
    Val  int
    Next *ListNode
}

// Doubly linked list node (LRU cache, browser history, etc.).
type DListNode struct {
    Val        int
    Prev, Next *DListNode
}
```

### Linked list vs array — the fundamental trade-off

| Operation                      | Array   | Linked List |
|--------------------------------|---------|-------------|
| Access i-th element            | O(1)    | O(n) — must walk from head |
| Insert/delete at known position| O(n) — shift elements | O(1) — relink pointers |
| Insert/delete at front         | O(n)    | O(1)        |
| Memory layout                  | contiguous, cache-friendly | scattered, pointer overhead |
| Binary search                  | O(log n)| impossible directly (no random access) |

The single most important consequence: **linked list algorithms are about
pointer surgery, not index arithmetic.** You cannot jump to the middle; you
can only walk. Every classic linked-list technique (slow/fast pointers, dummy
head, in-place reversal) exists to work around the lack of random access.

---

## How to recognise a linked-list problem

Signals in the problem statement:

1. **Explicit:** "Given the `head` of a linked list…" — the input type is
   `*ListNode`. Nearly every LeetCode linked-list problem says this verbatim.
2. **"Do it in-place" / "O(1) extra space"** on a sequence — often means
   relink existing nodes instead of copying values into an array.
3. **"Without modifying node values"** — forces true pointer manipulation
   (e.g. #24 Swap Nodes in Pairs bans the easy value-swap).
4. **Digits stored one per node** (e.g. #2 Add Two Numbers) — simulation of
   arithmetic where you cannot see all digits at once.
5. **"k-th node from the end"** — you cannot index from the back, so this
   screams two-pointer gap technique.
6. **Cycle detection** — "does the list loop?" → Floyd's tortoise & hare.
7. **Design problems** needing O(1) insert/delete at both ends with known
   node handles — LRU cache (#146) = hash map + doubly linked list.
8. **Merging sorted streams** (#21, #23) — linked lists let you splice nodes
   in O(1) without allocating a result array.

If the problem gives you an array and asks for insert-heavy behaviour, or
gives you a list and asks "find the middle / reverse / detect cycle", you are
in linked-list-pattern territory.

---

## Core templates (Go)

### 1. Traversal

```go
// Walk the list once. O(n) time, O(1) space.
for curr := head; curr != nil; curr = curr.Next {
    // process curr.Val
}
```

### 2. Dummy (sentinel) head — the #1 trick

Whenever the head itself might be removed, replaced, or unknown until the end,
put a fake node in front. All insert/delete logic then treats the real head
like any other node — no special cases.

```go
// Pseudocode:
//   dummy -> head -> ... (dummy.Next is always the answer's head)
//   keep a `prev` pointer starting at dummy
//   mutate prev.Next freely; return dummy.Next at the end
dummy := &ListNode{Next: head} // sentinel sits before the real head
prev := dummy                  // prev is always the last node we are keeping
for prev.Next != nil {
    if shouldDelete(prev.Next) {
        prev.Next = prev.Next.Next // unlink: skip over the doomed node
    } else {
        prev = prev.Next           // keep it: advance prev
    }
}
return dummy.Next // real head, even if the original head was deleted
```

### 3. In-place reversal — the second most important template

```go
// reverseList reverses a singly linked list iteratively.
//
// Invariant at the top of every iteration:
//   prev = head of the already-reversed prefix
//   curr = head of the not-yet-reversed suffix
//
// Time:  O(n) — each node visited once
// Space: O(1) — three pointers only
func reverseList(head *ListNode) *ListNode {
    var prev *ListNode  // nil: the reversed prefix is empty at first
    curr := head
    for curr != nil {
        next := curr.Next // 1. save the rest of the list BEFORE breaking the link
        curr.Next = prev  // 2. flip the arrow to point backwards
        prev = curr       // 3. reversed prefix grows by one node
        curr = next       // 4. step into the saved suffix
    }
    return prev // curr is nil; prev is the new head
}
```

Memorise the four-line body `next / flip / prev / curr` — it appears inside
#25 Reverse Nodes in k-Group, #92 Reverse Linked List II, #234 Palindrome
Linked List, and many more.

### 4. Slow/fast pointers (find middle, detect cycle, k-th from end)

```go
// Middle of the list: fast moves 2 steps per 1 of slow.
// When fast falls off the end, slow is at the middle.
slow, fast := head, head
for fast != nil && fast.Next != nil { // guard BOTH before fast.Next.Next
    slow = slow.Next
    fast = fast.Next.Next
}
// slow == middle (second middle for even length)

// Cycle detection (Floyd): same loop; if slow == fast inside the loop → cycle.

// n-th from the end: give fast a head start of n, then move both;
// when fast hits nil, slow is n nodes from the end.
```

### 5. Merge two sorted lists (splicing)

```go
// Splice nodes from l1 and l2 in sorted order onto a dummy tail.
dummy := &ListNode{}
tail := dummy
for l1 != nil && l2 != nil {
    if l1.Val <= l2.Val {
        tail.Next = l1 // link the smaller node — no allocation
        l1 = l1.Next
    } else {
        tail.Next = l2
        l2 = l2.Next
    }
    tail = tail.Next
}
if l1 != nil { tail.Next = l1 } else { tail.Next = l2 } // append leftovers
return dummy.Next
```

### 6. Recursion on lists

A list is either `nil` or `node + smaller list` — a natural recursive shape.

```go
// Recursive reversal. Time O(n), Space O(n) recursion stack.
func reverse(head *ListNode) *ListNode {
    if head == nil || head.Next == nil { // base: empty or single node
        return head
    }
    newHead := reverse(head.Next) // reverse the tail; newHead is the last node
    head.Next.Next = head        // make my successor point back at me
    head.Next = nil              // I become the (temporary) tail
    return newHead
}
```

Prefer iteration in interviews unless recursion is clearly cleaner — the O(n)
stack matters for n up to 10^5 and interviewers ask about it.

---

## Worked example — reverse `1 → 2 → 3 → 4`, traced step by step

Using the iterative template (#3 above). `∅` = nil.

| Step | Action                              | prev        | curr | next | List state (arrows)        |
|------|-------------------------------------|-------------|------|------|----------------------------|
| init | —                                   | ∅           | 1    | —    | 1→2→3→4→∅                  |
| 1a   | `next = curr.Next`                  | ∅           | 1    | 2    | 1→2→3→4→∅                  |
| 1b   | `curr.Next = prev` (1 now → ∅)      | ∅           | 1    | 2    | ∅←1  2→3→4→∅               |
| 1c   | `prev = curr; curr = next`          | 1           | 2    | 2    | ∅←1  2→3→4→∅               |
| 2    | save 3; flip 2→1; advance           | 2           | 3    | 3    | ∅←1←2  3→4→∅               |
| 3    | save 4; flip 3→2; advance           | 3           | 4    | 4    | ∅←1←2←3  4→∅               |
| 4    | save ∅; flip 4→3; advance           | 4           | ∅    | ∅    | ∅←1←2←3←4                  |
| exit | `curr == nil` → return `prev`       | **4** (new head) | | | 4→3→2→1→∅              |

Key observation: at every step the list is split into a fully reversed prefix
(headed by `prev`) and an untouched suffix (headed by `curr`), and `next` is
the lifeline that keeps the suffix reachable while we break the link.

---

## Common pitfalls and how to avoid them

1. **Losing the rest of the list.** Writing `curr.Next = prev` before saving
   `curr.Next` orphans the suffix forever. *Rule: save before you sever.*
2. **Forgetting the dummy head.** Deleting/inserting at the head without a
   sentinel forces ugly `if node == head` special cases and is the source of
   most off-by-one-node bugs. When in doubt, use a dummy — it costs one node.
3. **Nil-pointer dereference in fast-pointer loops.** `fast.Next.Next` panics
   if `fast.Next` is nil. Always guard `fast != nil && fast.Next != nil`, in
   that order (Go short-circuits `&&`).
4. **Not cutting the tail.** After reversal/reordering, the old tail may still
   point into the list, creating an accidental cycle. Explicitly set
   `node.Next = nil` where the new list should end (see the recursive reverse:
   `head.Next = nil`).
5. **Advancing the wrong pointer after deletion.** In the dummy-head delete
   loop, after `prev.Next = prev.Next.Next` you must **not** also advance
   `prev` — the new `prev.Next` might need deleting too.
6. **Confusing "second middle" vs "first middle"** for even-length lists.
   `for fast != nil && fast.Next != nil` lands slow on the **second** middle;
   `for fast.Next != nil && fast.Next.Next != nil` lands on the **first**.
   Pick deliberately — it changes split-in-half logic (#109, #234).
7. **Value-swapping when the problem forbids it.** Some problems (#24)
   explicitly require relinking nodes, not swapping `Val`s. Read constraints.
8. **Off-by-one in gap two-pointers.** For "n-th from end", decide whether
   fast starts `n` or `n+1` ahead depending on whether slow should land *on*
   the target or *before* it (before it, if you need to delete it — #19).
9. **Recursion depth.** n ≤ 10^5 with recursion = ~10^5 stack frames. Fine in
   Go usually, but state the O(n) space cost and know the iterative version.
10. **Modifying the input when asked not to.** Some problems require restoring
    the list or building a new one; note whether destruction is allowed.

---

## Decision cheat-sheet

| Problem smell | Technique |
|---|---|
| Head might be deleted / built from scratch | Dummy head |
| Middle, cycle, k-th from end | Slow/fast pointers |
| Reverse all / part / groups | In-place reversal (prev/curr/next) |
| Merge sorted lists | Splice with dummy + tail |
| Stable partition / split into sublists | Two dummy heads, join at end |
| O(1) insert+delete with lookup | Doubly linked list + hash map |
| k lists at once | Heap of heads, or divide & conquer pairwise merge |

---

## Problems in this repo

Problems whose solutions centrally use linked lists (0131+ will be added in a
later pass):

- [0002 — Add Two Numbers](../0002_add_two_numbers/README.md) — dummy head + carry simulation over two lists
- [0019 — Remove Nth Node From End of List](../0019_remove_nth_node_from_end_of_list/README.md) — two-pointer gap + dummy head
- [0021 — Merge Two Sorted Lists](../0021_merge_two_sorted_lists/README.md) — canonical splice-merge with dummy
- [0023 — Merge k Sorted Lists](../0023_merge_k_sorted_lists/README.md) — heap / divide & conquer over list heads
- [0024 — Swap Nodes in Pairs](../0024_swap_nodes_in_pairs/README.md) — pure pointer relinking (value swap forbidden)
- [0025 — Reverse Nodes in k-Group](../0025_reverse_nodes_in_k_group/README.md) — grouped in-place reversal
- [0061 — Rotate List](../0061_rotate_list/README.md) — close into a ring, then cut
- [0082 — Remove Duplicates from Sorted List II](../0082_remove_duplicates_from_sorted_list_ii/README.md) — dummy head, delete whole duplicate runs
- [0083 — Remove Duplicates from Sorted List](../0083_remove_duplicates_from_sorted_list/README.md) — in-place skip of adjacent duplicates
- [0086 — Partition List](../0086_partition_list/README.md) — two dummy heads, stable split + join
- [0092 — Reverse Linked List II](../0092_reverse_linked_list_ii/README.md) — reverse a sublist between positions
- [0109 — Convert Sorted List to Binary Search Tree](../0109_convert_sorted_list_to_binary_search_tree/README.md) — slow/fast middle-finding for balanced build
- [0114 — Flatten Binary Tree to Linked List](../0114_flatten_binary_tree_to_linked_list/README.md) — tree flattened into a right-pointer "list"
- [0116 — Populating Next Right Pointers in Each Node](../0116_populating_next_right_pointers_in_each_node/README.md) — treat each tree level as a linked list
- [0117 — Populating Next Right Pointers in Each Node II](../0117_populating_next_right_pointers_in_each_node_ii/README.md) — dummy head per level while wiring `next`

Related concepts: [`two_pointers.md`](two_pointers.md) (slow/fast is a linked-list
specialisation), [`hash_map.md`](hash_map.md) (node lookup in O(1), LRU cache).
