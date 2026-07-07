# 0237 — Delete Node in a Linked List

> LeetCode #237 · Difficulty: Medium
> **Categories:** Linked List

---

## Problem Statement

There is a singly-linked list `head` and we want to delete a node `node` in it.

You are given the node to be deleted `node`. You will **not be given access** to the first node of `head`.

All the values of the linked list are **unique**, and it is guaranteed that the given node `node` is **not the last node** in the linked list.

Delete the given node. Note that by deleting the node, we do not mean removing it from memory. We mean:

- The value of the given node should not exist in the linked list.
- The number of nodes in the linked list should decrease by one.
- All the values before `node` should be in the same order.
- All the values after `node` should be in the same order.

**Custom testing:**

- For the input, you should provide the entire linked list `head` and the node to be given `node`. `node` should not be the last node of the list and should be an actual node in the list.
- We will build the linked list and pass the node to your function.
- The output will be the entire list after calling your function.

**Example 1:**

```
Input: head = [4,5,1,9], node = 5
Output: [4,1,9]
Explanation: You are given the second node with value 5, the linked list should become 4 -> 1 -> 9 after calling your function.
```

**Example 2:**

```
Input: head = [4,5,1,9], node = 1
Output: [4,5,9]
Explanation: You are given the third node with value 1, the linked list should become 4 -> 5 -> 9 after calling your function.
```

**Constraints:**

- The number of the nodes in the given list is in the range `[2, 1000]`.
- `-1000 <= Node.val <= 1000`
- The value of each node in the list is **unique**.
- The `node` to be deleted is **in the list** and is **not a tail** node.

---

## Company Frequency

| Company    | Frequency        | Last Reported |
|------------|------------------|---------------|
| Amazon     | ★★★☆☆ Medium     | 2024          |
| Microsoft  | ★★★☆☆ Medium     | 2023          |
| Adobe      | ★★★☆☆ Medium     | 2023          |
| Apple      | ★★☆☆☆ Low        | 2023          |
| Bloomberg  | ★★☆☆☆ Low        | 2022          |

> ⚠️ Frequency data is crowd-sourced from LeetCode Discuss, Glassdoor, and
> community interview reports. Treat as a signal, not a guarantee.

---

## DSA Concepts Used

- **Linked List** — the whole trick is a pointer/value manipulation on a singly-linked list where you cannot reach the predecessor → see [`/dsa/linked_list.md`](/dsa/linked_list.md)

---

## Approaches Overview

| # | Approach | Time | Space | When to use |
|---|----------|------|-------|-------------|
| 1 | Copy-Next-Value then Skip (Optimal) | O(1) | O(1) | Always — the intended answer |
| 2 | Cascade-Shift Values | O(k) | O(1) | Only to illustrate the "overwrite + trim tail" mental model |

---

## Approach 1 — Copy-Next-Value then Skip (Optimal)

### Intuition

The catch: we have no reference to the node **before** `node`, so we cannot re-link the list to route around `node`. But nothing says the *physical* node must be the one removed. We make `node` **impersonate its successor**: copy the successor's value into `node`, then bypass the successor. The list now reads exactly as if `node` had been deleted — it is really the successor that got unlinked, but it was carrying the value we wanted gone. This is legal precisely because `node` is guaranteed not to be the tail, so a successor always exists.

### Algorithm

1. Copy `node.Next.Val` into `node.Val` — `node` now looks like its successor.
2. Set `node.Next = node.Next.Next` — unlink the original successor.

### Complexity

- **Time:** O(1) — two assignments, no traversal.
- **Space:** O(1) — nothing allocated.

### Code

```go
func copyNextAndSkip(node *ListNode) {
	node.Val = node.Next.Val   // steal the successor's value into this node
	node.Next = node.Next.Next // splice the (now duplicated) successor out
}
```

### Dry Run

Example 1: `head = [4,5,1,9]`, `node` = the node holding `5`.

| Step | Operation | List state |
|------|-----------|-----------|
| 0 | initial | `4 → 5 → 1 → 9` (`node` points at `5`) |
| 1 | `node.Val = node.Next.Val` (5 ← 1) | `4 → 1 → 1 → 9` |
| 2 | `node.Next = node.Next.Next` | `4 → 1 → 9` |

Result: `[4, 1, 9]` ✔ — the value `5` is gone and the length dropped by one.

---

## Approach 2 — Cascade-Shift Values

### Intuition

Deleting `node` is *equivalent* to shifting every following value one slot toward the head and then trimming the final node. This walks all the way to the tail, so it is strictly worse than Approach 1 — but it makes the "delete = overwrite forward + drop the last node" model explicit without depending on a single successor hop.

### Algorithm

1. Set `curr = node`.
2. While `curr.Next.Next != nil`: overwrite `curr.Val = curr.Next.Val`, then advance `curr = curr.Next`.
3. Now `curr.Next` is the last node: set `curr.Val = curr.Next.Val` and `curr.Next = nil` to trim the duplicate tail.

### Complexity

- **Time:** O(k) — `k` nodes from `node` to the tail are visited.
- **Space:** O(1) — a single roving pointer.

### Code

```go
func cascadeShift(node *ListNode) {
	curr := node // will slide down the tail copying values forward
	for curr.Next.Next != nil {
		curr.Val = curr.Next.Val // overwrite with the next value
		curr = curr.Next         // advance
	}
	curr.Val = curr.Next.Val // absorb the last node's value
	curr.Next = nil          // trim the now-duplicate tail
}
```

### Dry Run

Example 1: `head = [4,5,1,9]`, `node` = the node holding `5`.

| Step | `curr.Val` | Action | List state |
|------|-----------|--------|-----------|
| 0 | 5 | start at `node` | `4 → 5 → 1 → 9` |
| 1 | 5→1 | `curr.Val = 1`, advance | `4 → 1 → 1 → 9` |
| 2 | 1→9 | `curr.Next.Next == nil` → `curr.Val = 9`, `curr.Next = nil` | `4 → 1 → 9` |

Result: `[4, 1, 9]` ✔

---

## Key Takeaways

- **When you cannot reach the predecessor, delete the successor instead.** Copying `node.Next`'s value into `node` and skipping `node.Next` deletes the *value* in O(1) even though a different physical node is unlinked.
- This only works because the node is **guaranteed not to be the tail** — the successor copy would be impossible otherwise. Always confirm that guarantee before using the trick.
- The technique generalizes: "impersonate the neighbor" is useful anytime you have a node but not its context in a singly-linked structure.
- Contrast with normal deletion (#203, #19), which requires the predecessor pointer and therefore either a dummy head or a traversal.

---

## Related Problems

- LeetCode #203 — Remove Linked List Elements (needs predecessor / dummy head)
- LeetCode #19 — Remove Nth Node From End of List (two-pointer to find predecessor)
- LeetCode #83 — Remove Duplicates from Sorted List
- LeetCode #1836 — Remove Duplicates From an Unsorted Linked List
