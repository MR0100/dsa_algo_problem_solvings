# 0426 — Convert Binary Search Tree to Sorted Doubly Linked List

> LeetCode #426 · Difficulty: Medium
> **Categories:** Tree, Binary Search Tree, Depth-First Search, Linked List, Stack, Divide and Conquer

---

## Problem Statement

Convert a **Binary Search Tree** to a sorted **Circular Doubly-Linked List** in place.

You can think of the left and right pointers as synonymous to the predecessor and successor pointers in a doubly-linked list. For a circular doubly linked list, the predecessor of the first element is the last element, and the successor of the last element is the first element.

We want to do the transformation **in place**. After the transformation, the left pointer of the tree node should point to its predecessor, and the right pointer should point to its successor. You should return the pointer to the smallest element of the linked list.

**Example 1:**
```
Input: root = [4,2,5,1,3]
Output: [1,2,3,4,5]
```
Explanation: The figure below shows the transformed BST. The solid line indicates the successor relationship, while the dashed line means the predecessor relationship.
```
BST:                 Circular DLL (return head = node 1):

      4               ┌──────────────────────────────┐
     / \              ▼                              │
    2   5      1 ⇄ 2 ⇄ 3 ⇄ 4 ⇄ 5  (and 5.right = 1, 1.left = 5)
   / \
  1   3
```

**Example 2:**
```
Input: root = [2,1,3]
Output: [1,2,3]
```

**Example 3:**
```
Input: root = []
Output: []
```
Explanation: Input is an empty tree. Output is also an empty Linked List.

**Constraints:**
- The number of nodes in the tree is in the range `[0, 2000]`.
- `-1000 <= Node.val <= 1000`
- All the values of the tree are **unique**.

---

## Company Frequency

| Company    | Frequency        | Last Reported |
|------------|------------------|---------------|
| Facebook   | ★★★★☆ High       | 2024          |
| Amazon     | ★★★☆☆ Medium     | 2024          |
| Microsoft  | ★★★☆☆ Medium     | 2023          |
| Google     | ★★★☆☆ Medium     | 2023          |
| Bloomberg  | ★★☆☆☆ Low        | 2023          |
| Oracle     | ★★☆☆☆ Low        | 2023          |

> ⚠️ Frequency data is crowd-sourced from LeetCode Discuss, Glassdoor, and
> community interview reports. Treat as a signal, not a guarantee.

---

## DSA Concepts Used

- **Binary Search Tree in-order = sorted order** — the key property: an in-order walk of a BST visits values ascending, which is precisely the list order required → see [`/dsa/binary_search_tree.md`](/dsa/binary_search_tree.md)
- **Tree Traversal (in-order)** — recursive or stack-based left-node-right visiting drives every approach → see [`/dsa/tree_traversal.md`](/dsa/tree_traversal.md)
- **Doubly Linked List construction** — wiring `Left`=predecessor / `Right`=successor and closing the ring is in-place pointer surgery → see [`/dsa/linked_list.md`](/dsa/linked_list.md)
- **Stack (iterative traversal)** — Approach 3 replaces the call stack with an explicit stack to produce the same in-order sequence → see [`/dsa/stack.md`](/dsa/stack.md)

---

## Approaches Overview

Let n = number of nodes, h = tree height.

| # | Approach | Time | Space | When to use |
|---|----------|------|-------|-------------|
| 1 | In-order into a slice, then link | O(n) | O(n) | Clearest to reason about; when clarity beats memory |
| 2 | In-order with a running `prev` pointer (Optimal) | O(n) | O(h) | Interview-standard: single pass, no auxiliary array |
| 3 | Iterative in-order with an explicit stack | O(n) | O(h) | When recursion must be avoided (deep/unbalanced trees) |

---

## Approach 1 — In-order into a Slice, then Link

### Intuition
An in-order traversal of a BST produces its nodes in ascending value order — exactly the order the circular list must present. So decouple the two concerns: first collect all node **pointers** into a slice in sorted order (correctness is self-evident), then make a second pass linking neighbour `i` to neighbour `i+1` via `Right`/`Left`. Using modular indices to pick neighbours makes closing the circle (last↔first) fall out with no special-casing.

### Algorithm
1. In-order traverse, appending each node pointer to `nodes`.
2. If the tree was empty, return `nil`.
3. For every `i`: `nodes[i].Right = nodes[(i+1) mod n]` and `nodes[i].Left = nodes[(i-1+n) mod n]`.
4. Return `nodes[0]` — the smallest, i.e. the head.

### Complexity
- **Time:** O(n) — one traversal plus one linear linking pass.
- **Space:** O(n) — the slice of `n` node pointers, plus O(h) recursion stack.

### Code
```go
func inorderSlice(root *Node) *Node {
	if root == nil {
		return nil
	}
	nodes := []*Node{} // node pointers in ascending value order
	var inorder func(*Node)
	inorder = func(n *Node) {
		if n == nil {
			return
		}
		inorder(n.Left)          // left subtree: smaller values first
		nodes = append(nodes, n) // visit: record this node
		inorder(n.Right)         // right subtree: larger values
	}
	inorder(root)

	for i := 0; i < len(nodes); i++ {
		nodes[i].Right = nodes[(i+1)%len(nodes)]           // successor (wraps to head at the end)
		nodes[i].Left = nodes[(i-1+len(nodes))%len(nodes)] // predecessor (wraps to tail at start)
	}
	return nodes[0] // smallest element = head of the circular DLL
}
```

### Dry Run (Example 1)

BST `[4,2,5,1,3]`. In-order collects `nodes = [1,2,3,4,5]` (n=5).

| i | node | `Right = nodes[(i+1)%5]` | `Left = nodes[(i-1+5)%5]` |
|---|------|--------------------------|---------------------------|
| 0 | 1 | nodes[1]=2 | nodes[4]=5 |
| 1 | 2 | nodes[2]=3 | nodes[0]=1 |
| 2 | 3 | nodes[3]=4 | nodes[1]=2 |
| 3 | 4 | nodes[4]=5 | nodes[2]=3 |
| 4 | 5 | nodes[0]=1 | nodes[3]=4 |

Return `nodes[0]` = node 1. Forward (`Right`): `1→2→3→4→5→1`; backward (`Left`): `5→4→3→2→1→5`. Output `[1,2,3,4,5]` ✓

---

## Approach 2 — In-order with a Running Previous Pointer (Optimal)

### Intuition
The slice in Approach 1 exists only to know each node's predecessor. But during an in-order walk we *already* visit nodes in that exact order — so at the instant we visit `cur`, the previously-visited node `prev` **is** its predecessor. Link them on the spot (`prev.Right = cur`, `cur.Left = prev`) and no array is needed. Remember the first visited node as `head`; when the walk finishes, `prev` is the tail, and we close the ring.

### Algorithm
1. Keep `head` (first visited) and `prev` (last visited), both initially `nil`.
2. In-order: recurse left; **on visit**, if `prev != nil` set `prev.Right = cur`, `cur.Left = prev`; else `head = cur` (this is the minimum). Then `prev = cur`; recurse right.
3. After the walk, close the circle: `head.Left = prev`, `prev.Right = head`.
4. Return `head`.

### Complexity
- **Time:** O(n) — each node visited exactly once.
- **Space:** O(h) — only the recursion stack; no auxiliary array.

### Code
```go
func inorderInPlace(root *Node) *Node {
	if root == nil {
		return nil
	}
	var head, prev *Node // head = smallest node; prev = last node linked so far
	var inorder func(*Node)
	inorder = func(cur *Node) {
		if cur == nil {
			return
		}
		inorder(cur.Left) // recurse left first (smaller values)
		if prev != nil {
			prev.Right = cur // previous node's successor is this node
			cur.Left = prev  // this node's predecessor is the previous node
		} else {
			head = cur // no predecessor yet ⇒ this is the smallest node
		}
		prev = cur         // this node becomes the "previous" for the next visit
		inorder(cur.Right) // then recurse right (larger values)
	}
	inorder(root)

	// Close the circle: smallest's predecessor is the largest and vice versa.
	head.Left = prev
	prev.Right = head
	return head
}
```

### Dry Run (Example 1)

In-order visit order: 1, 2, 3, 4, 5.

| Visit `cur` | `prev` before | Action | `head` | `prev` after |
|-------------|---------------|--------|--------|--------------|
| 1 | nil | `head = 1` | 1 | 1 |
| 2 | 1 | `1.Right=2`, `2.Left=1` | 1 | 2 |
| 3 | 2 | `2.Right=3`, `3.Left=2` | 1 | 3 |
| 4 | 3 | `3.Right=4`, `4.Left=3` | 1 | 4 |
| 5 | 4 | `4.Right=5`, `5.Left=4` | 1 | 5 |
| — (end) | — | `head.Left=5` (`1.Left=5`), `prev.Right=head` (`5.Right=1`) | 1 | 5 |

Return `head` = 1. Forward `1,2,3,4,5`; backward `5,4,3,2,1` ✓

---

## Approach 3 — Iterative In-order with an Explicit Stack

### Intuition
Recursion contributed nothing but the in-order visit *sequence*; an explicit stack reproduces that sequence without risking a deep call stack on a skewed tree. Push the entire left spine, pop to visit the next-smallest node, then move into its right subtree and repeat. The `prev`/`head` linking is byte-for-byte the same as Approach 2 — only the traversal driver changes.

### Algorithm
1. `stack = []`, `cur = root`, `head = prev = nil`.
2. While `cur != nil` or the stack is non-empty:
   1. Push every left descendant: while `cur != nil`, push `cur`, `cur = cur.Left`.
   2. Pop `cur` — the next node in ascending order. Link to `prev` (or set `head`); `prev = cur`.
   3. `cur = cur.Right`.
3. Close the ring: `head.Left = prev`, `prev.Right = head`; return `head`.

### Complexity
- **Time:** O(n) — each node is pushed and popped exactly once.
- **Space:** O(h) — the explicit stack holds at most one root-to-leaf path.

### Code
```go
func iterativeInorder(root *Node) *Node {
	if root == nil {
		return nil
	}
	var head, prev *Node
	stack := []*Node{}
	cur := root
	for cur != nil || len(stack) > 0 {
		for cur != nil {
			stack = append(stack, cur) // remember the path down the left spine
			cur = cur.Left
		}
		cur = stack[len(stack)-1]    // top = next node in ascending order
		stack = stack[:len(stack)-1] // pop it

		if prev != nil {
			prev.Right = cur // link predecessor → current
			cur.Left = prev
		} else {
			head = cur // first popped node is the minimum
		}
		prev = cur

		cur = cur.Right // explore the right subtree next
	}
	head.Left = prev
	prev.Right = head
	return head
}
```

### Dry Run (Example 1)

BST `[4,2,5,1,3]`: root 4, 4.Left=2 (2.Left=1, 2.Right=3), 4.Right=5.

| Phase | `cur` in | Stack after left-push | Pop → `cur` | Link | `prev` | Move to `cur.Right` |
|-------|----------|-----------------------|-------------|------|--------|---------------------|
| a | 4 | [4,2,1] | 1 | head=1 | 1 | 1.Right = nil |
| b | nil | [4,2] | 2 | 1.Right=2, 2.Left=1 | 2 | 2.Right = 3 |
| c | 3 | [4,3] | 3 | 2.Right=3, 3.Left=2 | 3 | 3.Right = nil |
| d | nil | [4] | 4 | 3.Right=4, 4.Left=3 | 4 | 4.Right = 5 |
| e | 5 | [5] | 5 | 4.Right=5, 5.Left=4 | 5 | 5.Right = nil |
| end | nil | [] | — | 1.Left=5, 5.Right=1 | 5 | — |

Return `head` = 1 → `[1,2,3,4,5]` ✓

---

## Key Takeaways

- **In-order traversal of a BST is the sorted sequence** — the single fact this whole problem rests on. Any "turn a BST into something ordered" task starts here.
- **The running-`prev` trick** turns "collect then link" into "link while visiting", dropping O(n) auxiliary space to O(h). It generalises to #114 (flatten to a list), #94/#173 (in-order iterator), and #530 (min abs difference in BST).
- **Close the circle exactly once, after the walk**: `head.Left = prev; prev.Right = head`. Doing it inside the traversal is the classic off-by-one that breaks circularity or creates a self-loop on the head.
- The recursive and stack-based in-order produce an identical node order; choose the stack version only when call-stack depth (up to n on a degenerate tree) is a real risk.
- Return the **smallest** node (`head`), not the root — a subtle spec detail that trips up first attempts.

---

## Related Problems
- LeetCode #94 — Binary Tree Inorder Traversal (the traversal at the core)
- LeetCode #114 — Flatten Binary Tree to Linked List (same "linearise a tree in place" idea)
- LeetCode #173 — Binary Search Tree Iterator (stack-based in-order, on demand)
- LeetCode #109 — Convert Sorted List to Binary Search Tree (the inverse transformation)
- LeetCode #430 — Flatten a Multilevel Doubly Linked List (in-place DLL splicing)
- LeetCode #708 — Insert into a Sorted Circular Linked List (circular DLL pointer surgery)
