# 0430 — Flatten a Multilevel Doubly Linked List

> LeetCode #430 · Difficulty: Medium
> **Categories:** Linked List, Depth-First Search, Stack, Doubly Linked List

---

## Problem Statement

You are given a doubly linked list, which contains nodes that have a next pointer, a previous pointer, and an additional **child pointer**. This child pointer may or may not point to a separate doubly linked list, also containing these special nodes. These child lists may have one or more children of their own, and so on, to produce a **multilevel data structure**, as shown in the example below.

Given the `head` of the first level of the list, **flatten** the list so that all the nodes appear in a single-level, doubly linked list. Let `curr` be a node with a child list. The nodes in the child list should appear **after** `curr` and **before** `curr.next` in the flattened list.

Return the `head` of the flattened list. The nodes in the list must have **all** of their child pointers set to `null`.

**Example 1:**
```
Input: head = [1,2,3,4,5,6,null,null,null,7,8,9,10,null,null,11,12]
Output: [1,2,3,7,8,11,12,9,10,4,5,6]
```
Explanation: The multilevel linked list in the input is as follows:
```
1---2---3---4---5---6--NULL
        |
        7---8---9---10--NULL
            |
            11--12--NULL
```
After flattening the multilevel linked list it becomes:
```
1-2-3-7-8-11-12-9-10-4-5-6-NULL
```

**Example 2:**
```
Input: head = [1,2,null,3]
Output: [1,3,2]
```
Explanation: The multilevel linked list in the input is as follows:
```
1---2---NULL
|
3---NULL
```

**Example 3:**
```
Input: head = []
Output: []
```
Explanation: There could be empty list in the input.

**Constraints:**
- The number of Nodes will not exceed `1000`.
- `1 <= Node.val <= 10⁵`

**How the multilevel linked list is represented in test cases:**
We use the multilevel linked list from **Example 1** above:
```
 1---2---3---4---5---6--NULL
         |
         7---8---9---10--NULL
             |
             11--12--NULL
```
The serialization of each level is fixed. After serializing each level, we concatenate them, replacing empty levels' serialization with `null`, and trimming any trailing `null`s from the whole list.

---

## Company Frequency

| Company    | Frequency        | Last Reported |
|------------|------------------|---------------|
| Amazon     | ★★★★☆ High       | 2024          |
| Microsoft  | ★★★☆☆ Medium     | 2024          |
| Facebook   | ★★★☆☆ Medium     | 2023          |
| Bloomberg  | ★★★☆☆ Medium     | 2023          |
| Google     | ★★☆☆☆ Low        | 2023          |

> ⚠️ Frequency data is crowd-sourced from LeetCode Discuss, Glassdoor, and
> community interview reports. Treat as a signal, not a guarantee.

---

## DSA Concepts Used

- **Doubly Linked List surgery** — the whole task is careful re-splicing of `Prev`/`Next`/`Child` pointers in place → see [`/dsa/linked_list.md`](/dsa/linked_list.md)
- **Depth-First Search** — diving into a child before continuing the current level is a pre-order DFS; the child branch is fully consumed before the deferred successor → see [`/dsa/graph_bfs_dfs.md`](/dsa/graph_bfs_dfs.md)
- **Stack (explicit or call stack)** — remembering the "resume here" successor while a child branch is flattened is a LIFO discipline, whether via recursion or an explicit stack → see [`/dsa/stack.md`](/dsa/stack.md)

---

## Approaches Overview

Let n = total nodes, d = maximum nesting depth.

| # | Approach | Time | Space | When to use |
|---|----------|------|-------|-------------|
| 1 | Recursive flatten (return the tail) | O(n) | O(d) | Cleanest expression; the tail-return trick makes splices trivial |
| 2 | Iterative flatten with an explicit stack | O(n) | O(d) | When recursion depth is a concern, or an interviewer asks for O(1) recursion |

---

## Approach 1 — Recursive Flatten (Return the Tail)

### Intuition
When a node `cur` has a child, that child's **entire flattened sublist** must be inserted between `cur` and `cur.Next`. The splice is easy if a helper can flatten a sublist *and hand back its tail*: link `cur → childHead` on the left and `childTail → oldNext` on the right — three pointer updates. Recursion makes children-of-children free, because flattening a sublist recurses into *its* children first, so by the time we get its tail, everything below is already a single level.

### Algorithm
1. Walk the current level with a cursor `cur`, remembering `next = cur.Next` (the *original* successor) each step.
2. If `cur.Child != nil`:
   1. Recursively flatten the child, getting `childTail`.
   2. Link `cur ↔ cur.Child`; set `cur.Child = nil`.
   3. Link `childTail ↔ next` (if `next` exists).
3. Advance `cur = next`; keep the last non-nil node as this level's `tail`.
4. Return `tail` so a parent-level splice can chain onto it.

### Complexity
- **Time:** O(n) — every node is visited a constant number of times (once on its own level, plus the pointer fix-ups).
- **Space:** O(d) — recursion depth equals the nesting depth `d` (up to `n` in the fully-nested worst case).

### Code
```go
func recursiveFlatten(head *Node) *Node {
	flatten(head) // ignore the returned tail at the top level
	return head
}

// flatten flattens the sublist starting at `head` in place and returns its
// tail (last node). Returns nil for a nil head.
func flatten(head *Node) *Node {
	cur := head
	var tail *Node // last node seen so far on this level
	for cur != nil {
		next := cur.Next // the ORIGINAL successor, before any splice
		if cur.Child != nil {
			childTail := flatten(cur.Child) // recursively flatten the nested list

			cur.Next = cur.Child  // node now points into the child list
			cur.Child.Prev = cur  // back-link the child head to node
			cur.Child = nil       // child pointer must be cleared per the spec

			childTail.Next = next // stitch the child's tail to the old successor
			if next != nil {
				next.Prev = childTail // back-link if there was a successor
			}
		}
		tail = cur  // cur is (so far) the furthest node on this level
		cur = next  // continue with the original successor
	}
	return tail
}
```

### Dry Run (Example 1)

Levels: top `1-2-3-4-5-6`, node 3's child `7-8-9-10`, node 8's child `11-12`.

| Event | Action | Resulting links (relevant) |
|-------|--------|----------------------------|
| `flatten(1..6)` reaches node 3 | `next = 4`; node 3 has child → recurse `flatten(7..10)` | pending: reconnect `4` after child tail |
| `flatten(7..10)` reaches node 8 | `next = 9`; node 8 has child → recurse `flatten(11,12)` | pending: reconnect `9` after `12` |
| `flatten(11,12)` returns | tail = `12` | `11-12` unchanged |
| back in `7..10` | splice: `8 → 11`, `12 → 9` | `7-8-11-12-9-10`, returns tail `10` |
| back in `1..6` | splice: `3 → 7`, `10 → 4` | `1-2-3-7-8-11-12-9-10-4-5-6` |

Walking `Next` from head: `1,2,3,7,8,11,12,9,10,4,5,6` ✓ — and every `Child` was set to `nil` during its splice.

---

## Approach 2 — Iterative Flatten with an Explicit Stack

### Intuition
Diving into a child is a depth-first **detour**: the node we *would* have visited next on the current level (`cur.Next`) is where we must resume once the child branch is exhausted. Push those deferred successors onto a stack. Then just walk `Next`: whenever the current node runs out of `Next` (`cur.Next == nil`) and the stack is non-empty, pop the most recently deferred successor and splice it on — that reconnects the branch we detoured away from, innermost first (LIFO), exactly matching the recursive order.

### Algorithm
1. `cur = head`. While `cur != nil`:
   1. If `cur.Child != nil`: push `cur.Next` (if any) onto the stack; set `cur.Next = cur.Child`, back-link, and clear `cur.Child`.
   2. If `cur.Next == nil` and the stack is non-empty: pop `top`, set `cur.Next = top`, `top.Prev = cur`.
   3. `cur = cur.Next`.

### Complexity
- **Time:** O(n) — each node is processed exactly once as `cur`.
- **Space:** O(d) — the stack holds at most one deferred successor per active nesting level.

### Code
```go
func stackFlatten(head *Node) *Node {
	if head == nil {
		return nil
	}
	stack := []*Node{} // deferred "resume here" successors, LIFO
	cur := head
	for cur != nil {
		if cur.Child != nil {
			if cur.Next != nil {
				stack = append(stack, cur.Next) // defer the current-level successor
			}
			cur.Next = cur.Child // dive into the child list
			cur.Child.Prev = cur
			cur.Child = nil // clear per the spec
		}
		if cur.Next == nil && len(stack) > 0 {
			top := stack[len(stack)-1]  // most recently deferred successor
			stack = stack[:len(stack)-1] // pop it
			cur.Next = top               // reconnect it after the finished branch
			top.Prev = cur
		}
		cur = cur.Next // move forward along the now-single-level list
	}
	return head
}
```

### Dry Run (Example 1)

| `cur` | Child? | Stack action | Next==nil? pop? | `cur.Next` becomes | Stack after |
|-------|--------|--------------|-----------------|--------------------|-------------|
| 1 | no | — | no | 2 | `[]` |
| 2 | no | — | no | 3 | `[]` |
| 3 | yes (7) | push 4 | no | 7 | `[4]` |
| 7 | no | — | no | 8 | `[4]` |
| 8 | yes (11) | push 9 | no | 11 | `[4,9]` |
| 11 | no | — | no | 12 | `[4,9]` |
| 12 | no | — | yes → pop 9 | 9 | `[4]` |
| 9 | no | — | no | 10 | `[4]` |
| 10 | no | — | yes → pop 4 | 4 | `[]` |
| 4 | no | — | no | 5 | `[]` |
| 5 | no | — | no | 6 | `[]` |
| 6 | no | — | yes, stack empty | nil (stop) | `[]` |

Order: `1,2,3,7,8,11,12,9,10,4,5,6` ✓

---

## Key Takeaways

- **"Insert child between `cur` and `cur.Next`" ⇒ you need the child's tail.** Whether you get it by returning it (recursion) or by stacking the deferred `Next` (iteration), the splice is the same three-pointer dance on both sides.
- **Recursion's call stack == the explicit stack.** The iterative version just makes the deferred-successor stack visible; both flatten innermost-first (LIFO / depth-first).
- **Always snapshot `cur.Next` before mutating**, and remember to fix `Prev` pointers *and* null out `Child` — the three easiest bugs (forgetting the back-link, forgetting to clear `Child`, or losing the original successor) each fail a hidden test.
- A single-child, single-node case (Example 2 → `1,3,2`) is the minimal test that still exercises the full splice on both ends.

---

## Related Problems
- LeetCode #114 — Flatten Binary Tree to Linked List (same "splice a subtree in-line" idea on a tree)
- LeetCode #21 — Merge Two Sorted Lists (linked-list pointer surgery)
- LeetCode #138 — Copy List with Random Pointer (extra-pointer linked list manipulation)
- LeetCode #708 — Insert into a Sorted Circular Linked List (careful in-place splicing)
- LeetCode #369 — Plus One Linked List (recursion returning info up the list)
