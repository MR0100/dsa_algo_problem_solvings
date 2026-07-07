# 0160 — Intersection of Two Linked Lists

> LeetCode #160 · Difficulty: Easy
> **Categories:** Hash Table, Linked List, Two Pointers

---

## Problem Statement

Given the heads of two singly linked-lists `headA` and `headB`, return *the node at which the two lists intersect*. If the two linked lists have no intersection at all, return `null`.

For example, the following two linked lists begin to intersect at node `c1`:

```
A:      a1 → a2
                ↘
                  c1 → c2 → c3
                ↗
B: b1 → b2 → b3
```

The test cases are generated such that there are no cycles anywhere in the entire linked structure.

**Note** that the linked lists must **retain their original structure** after the function returns.

**Custom Judge:**

The inputs to the **judge** are given as follows (your program is **not** given these inputs):

- `intersectVal` — The value of the node where the intersection occurs. This is `0` if there is no intersected node.
- `listA` — The first linked list.
- `listB` — The second linked list.
- `skipA` — The number of nodes to skip ahead in `listA` (starting from the head) to get to the intersected node.
- `skipB` — The number of nodes to skip ahead in `listB` (starting from the head) to get to the intersected node.

The judge will then create the linked structure based on these inputs and pass the two heads, `headA` and `headB` to your program. If you correctly return the intersected node, then your solution will be **accepted**.

**Example 1:**
```
Input: intersectVal = 8, listA = [4,1,8,4,5], listB = [5,6,1,8,4,5], skipA = 2, skipB = 3
Output: Intersected at '8'
Explanation: The intersected node's value is 8 (note that this must not be 0 if the two lists intersect).
From the head of A, it reads as [4,1,8,4,5]. From the head of B, it reads as [5,6,1,8,4,5]. There are 2 nodes before the intersected node in A; There are 3 nodes before the intersected node in B.
- Note that the intersected node's value is not 1 because the nodes with value 1 in A and B (2nd node in A and 3rd node in B) are different node references. In other words, they point to two different locations in memory, while the nodes with value 8 in A and B (3rd node in A and 4th node in B) point to the same location in memory.
```

**Example 2:**
```
Input: intersectVal = 2, listA = [1,9,1,2,4], listB = [3,2,4], skipA = 3, skipB = 1
Output: Intersected at '2'
Explanation: The intersected node's value is 2 (note that this must not be 0 if the two lists intersect).
From the head of A, it reads as [1,9,1,2,4]. From the head of B, it reads as [3,2,4]. There are 3 nodes before the intersected node in A; There are 1 node before the intersected node in B.
```

**Example 3:**
```
Input: intersectVal = 0, listA = [2,6,4], listB = [1,5], skipA = 3, skipB = 2
Output: No intersection
Explanation: From the head of A, it reads as [2,6,4]. From the head of B, it reads as [1,5]. Since the two lists do not intersect, intersectVal must be 0, while skipA and skipB can be arbitrary values.
Explanation: The two lists do not intersect, so return null.
```

**Constraints:**
- The number of nodes of `listA` is in the `m`.
- The number of nodes of `listB` is in the `n`.
- `1 <= m, n <= 3 * 10^4`
- `1 <= Node.val <= 10^5`
- `0 <= skipA <= m`
- `0 <= skipB <= n`
- `intersectVal` is `0` if `listA` and `listB` do not intersect.
- `intersectVal == listA[skipA] == listB[skipB]` if `listA` and `listB` intersect.

**Follow up:** Could you write a solution that runs in `O(m + n)` time and use only `O(1)` memory?

---

## Company Frequency

| Company    | Frequency       | Last Reported |
|------------|-----------------|---------------|
| Amazon     | ★★★★★ Very High | 2024          |
| Microsoft  | ★★★★☆ High      | 2024          |
| Google     | ★★★★☆ High      | 2024          |
| Facebook   | ★★★☆☆ Medium    | 2023          |
| Bloomberg  | ★★★☆☆ Medium    | 2023          |
| Adobe      | ★★☆☆☆ Low       | 2022          |

> ⚠️ Frequency data is crowd-sourced from LeetCode Discuss, Glassdoor, and
> community interview reports. Treat as a signal, not a guarantee.

---

## DSA Concepts Used

- **Linked List** — Y-shaped structure: two prefixes merging into one shared tail; intersection is pointer identity, never value equality → see [`/dsa/linked_list.md`](/dsa/linked_list.md)
- **Two Pointers** — the optimal solution walks two cursors over swapped paths so they travel equal total distance → see [`/dsa/two_pointers.md`](/dsa/two_pointers.md)
- **Hash Map / Set** — node pointers as keys give an easy O(m+n) solution at O(m) space → see [`/dsa/hash_map.md`](/dsa/hash_map.md)

---

## Approaches Overview

| # | Approach | Time | Space | When to use |
|---|----------|------|-------|-------------|
| 1 | Brute Force (Nested Scan) | O(m·n) | O(1) | Baseline only; too slow at 3·10⁴ nodes |
| 2 | Hash Set | O(m+n) | O(m) | Quick to write; fails the O(1)-memory follow-up |
| 3 | Length Difference | O(m+n) | O(1) | Meets the follow-up; explicit and easy to explain |
| 4 | Two Pointers Switching Heads (Optimal) | O(m+n) | O(1) | Meets the follow-up with the least code; the classic |

---

## Approach 1 — Brute Force (Nested Scan)

### Intuition
Intersection means *the same node in memory*, so compare pointers, not values. For every node `a` of list A, scan all of list B looking for a node `b` with `a == b`. Scanning A front-to-back guarantees the first match is the earliest shared node — the intersection start.

### Algorithm
1. For each node `a` in A (from `headA`):
   1. For each node `b` in B (from `headB`):
      1. If `a == b` (pointer identity), return `a`.
2. If the loops finish without a match, return `nil`.

### Complexity
- **Time:** O(m·n) — up to every pair of nodes is compared (≈ 9·10⁸ at the max constraints — too slow in practice).
- **Space:** O(1) — two roaming cursors only.

### Code
```go
func bruteForce(headA, headB *ListNode) *ListNode {
    for a := headA; a != nil; a = a.Next {
        for b := headB; b != nil; b = b.Next {
            if a == b {
                return a
            }
        }
    }
    return nil
}
```

### Dry Run
Example 1: `A = 4→1→8→4→5`, `B = 5→6→1→8→4→5` (nodes `8→4→5` shared). Write A's nodes as `a1(4) a2(1) c1(8) c2(4) c3(5)` and B's prefix as `b1(5) b2(6) b3(1)`.

| Outer a | Inner scan of B | Match? |
|---------|-----------------|--------|
| a1 (4) | b1,b2,b3,c1,c2,c3 | no — a1 is not in B (note b1 has value 5, c2 value 4: values irrelevant, pointers differ) |
| a2 (1) | b1,b2,b3,c1,c2,c3 | no — b3 also holds value 1 but is a *different* node |
| c1 (8) | b1,b2,b3,**c1** | **yes** — same pointer reached via B after 3 steps |

Return node `c1` → "Intersected at '8'" ✓

---

## Approach 2 — Hash Set

### Intuition
Trade memory for a single pass each: record every node pointer of list A in a set, then walk B — the first node already in the set is where B merges into A. Everything after the merge point is shared, so the first hit is necessarily the intersection start.

### Algorithm
1. Create an empty set keyed by node pointer.
2. Walk A, inserting every node.
3. Walk B; the first node found in the set → return it.
4. B exhausted → return `nil`.

### Complexity
- **Time:** O(m+n) — one pass to fill the set, one pass to probe it (O(1) average per op).
- **Space:** O(m) — the set stores all of A's nodes.

### Code
```go
func hashSet(headA, headB *ListNode) *ListNode {
    seen := map[*ListNode]bool{}
    for a := headA; a != nil; a = a.Next {
        seen[a] = true
    }
    for b := headB; b != nil; b = b.Next {
        if seen[b] {
            return b
        }
    }
    return nil
}
```

### Dry Run
Example 1: `A = a1(4) a2(1) c1(8) c2(4) c3(5)`, `B = b1(5) b2(6) b3(1) → c1 c2 c3`.

| Phase | Step | seen / probe result |
|-------|------|---------------------|
| fill | insert a1, a2, c1, c2, c3 | seen = {a1, a2, c1, c2, c3} |
| probe | b1 (5) | not in seen |
| probe | b2 (6) | not in seen |
| probe | b3 (1) | not in seen — value 1 equals a2's value, but the pointer differs |
| probe | c1 (8) | **in seen → return c1** |

Return `c1` → "Intersected at '8'" ✓

---

## Approach 3 — Length Difference

### Intuition
Intersecting lists share a common **tail**, so measured *from the end* the intersection sits at the same distance in both lists. The only reason simple lockstep walking fails is the unequal prefix lengths. Fix that first: advance the longer list's cursor by `|m−n|` nodes. Now both cursors are the same distance from the end, and walking them together must land them on the intersection node simultaneously — or on `nil` simultaneously when disjoint.

### Algorithm
1. Count `lenA` (pass over A) and `lenB` (pass over B).
2. Set cursors `a = headA`, `b = headB`; advance the cursor of the longer list by `|lenA − lenB|` steps.
3. While `a != b`: advance both one step.
4. Return `a` (the shared node, or `nil` if both fell off the end together).

### Complexity
- **Time:** O(m+n) — two measuring passes plus one aligned pass, all linear.
- **Space:** O(1) — two counters and two cursors.

### Code
```go
func lengthDifference(headA, headB *ListNode) *ListNode {
    lenA, lenB := 0, 0
    for n := headA; n != nil; n = n.Next {
        lenA++
    }
    for n := headB; n != nil; n = n.Next {
        lenB++
    }
    a, b := headA, headB
    for ; lenA > lenB; lenA-- {
        a = a.Next
    }
    for ; lenB > lenA; lenB-- {
        b = b.Next
    }
    for a != b {
        a = a.Next
        b = b.Next
    }
    return a
}
```

### Dry Run
Example 1: `A = a1(4) a2(1) c1(8) c2(4) c3(5)` (len 5), `B = b1(5) b2(6) b3(1) c1 c2 c3` (len 6).

| Step | a | b | Note |
|------|---|---|------|
| measure | — | — | lenA = 5, lenB = 6 |
| align | a1 | b2 | B longer by 1 → b skips b1 |
| walk 1 | a2 | b3 | a1≠b2 was checked; still different nodes |
| walk 2 | c1 | c1 | pointers now EQUAL → loop exits |

Return `c1` → "Intersected at '8'" ✓

---

## Approach 4 — Two Pointers Switching Heads (Optimal)

### Intuition
Decompose the lists as `A = a + c` and `B = b + c` (`a`, `b` = exclusive prefixes, `c` = shared tail, possibly empty). Let pointer `pa` traverse `A` then continue from `headB`; let `pb` traverse `B` then continue from `headA`. Each travels `a + c + b` total steps, so after exactly `a + b` steps both stand at the start of `c` — the intersection — at the same moment. If `c` is empty, they reach `nil` together after `a + b` steps and the loop exits returning `nil`. Length arithmetic happens implicitly; no counting, no extra memory.

### Algorithm
1. If either head is `nil`, return `nil`.
2. `pa = headA`, `pb = headB`.
3. While `pa != pb`:
   1. `pa = headB` if `pa == nil`, else `pa = pa.Next`.
   2. `pb = headA` if `pb == nil`, else `pb = pb.Next`.
4. Return `pa` — the intersection node, or `nil` when both exhausted their swapped paths.

### Complexity
- **Time:** O(m+n) — each pointer visits each list at most once (≤ m+n+1 steps each).
- **Space:** O(1) — exactly two pointers.

### Code
```go
func twoPointers(headA, headB *ListNode) *ListNode {
    if headA == nil || headB == nil {
        return nil
    }
    pa, pb := headA, headB
    for pa != pb {
        if pa == nil {
            pa = headB
        } else {
            pa = pa.Next
        }
        if pb == nil {
            pb = headA
        } else {
            pb = pb.Next
        }
    }
    return pa
}
```

### Dry Run
Example 1: `A = a1(4) a2(1) c1(8) c2(4) c3(5)`, `B = b1(5) b2(6) b3(1) c1 c2 c3`. Path lengths: a = 2, b = 3, c = 3.

| Step | pa | pb | pa == pb? |
|------|----|----|-----------|
| start | a1(4) | b1(5) | no |
| 1 | a2(1) | b2(6) | no |
| 2 | c1(8) | b3(1) | no |
| 3 | c2(4) | c1(8) | no |
| 4 | c3(5) | c2(4) | no |
| 5 | nil | c3(5) | no |
| 6 | b1(5) (switched to B) | nil | no |
| 7 | b2(6) | a1(4) (switched to A) | no |
| 8 | b3(1) | a2(1) | no — equal values (1), different pointers! |
| 9 | **c1(8)** | **c1(8)** | **yes — same node** |

Both pointers travelled a + c + b = 2 + 3 + 3 = 8 transitions past the start to reach `c1` together. Return `c1` → "Intersected at '8'" ✓

---

## Key Takeaways
- **Identity, not equality**: intersection questions on linked structures compare *pointers*; equal values (like the two `1` nodes in Example 1) are red herrings.
- Intersecting singly linked lists form a **Y, never an X**: once merged they can never diverge (each node has one `Next`), so the shared part is always a suffix — this is what makes tail-alignment reasoning valid.
- **Equalise path lengths** is the master trick: do it explicitly (count and skip, Approach 3) or implicitly (swap heads so both walk a+b+c, Approach 4).
- The head-switching loop terminates even with no intersection because both pointers hit `nil` at the same step — `nil == nil` acts as the natural sentinel; don't special-case it away.
- Same "align two traversals" idea reappears in cycle detection (#141/#142) and in finding the k-th node from the end (#19).

---

## Related Problems
- LeetCode #141 — Linked List Cycle (pointer-meeting argument, same family)
- LeetCode #142 — Linked List Cycle II (find the exact meeting node)
- LeetCode #19 — Remove Nth Node From End of List (gap-aligned two pointers)
- LeetCode #599 — Minimum Index Sum of Two Lists (hash-set intersection of two sequences)
