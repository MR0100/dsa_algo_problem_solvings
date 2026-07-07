# 0230 — Kth Smallest Element in a BST

> LeetCode #230 · Difficulty: Medium
> **Categories:** Tree, Depth-First Search, Binary Search Tree, Binary Tree

---

## Problem Statement

Given the `root` of a binary search tree, and an integer `k`, return *the* `kᵗʰ` *smallest value (**1-indexed**) of all the values of the nodes in the tree*.

**Example 1:**
```
Input: root = [3,1,4,null,2], k = 1
Output: 1
```

**Example 2:**
```
Input: root = [5,3,6,2,4,null,null,1], k = 3
Output: 3
```

**Constraints:**
- The number of nodes in the tree is `n`.
- `1 <= k <= n <= 10⁴`
- `0 <= Node.val <= 10⁴`

**Follow-up:** If the BST is modified often (i.e., we can do insert and delete operations) and you need to find the kth smallest frequently, how would you optimize?

---

## Company Frequency

| Company    | Frequency        | Last Reported |
|------------|------------------|---------------|
| Amazon     | ★★★★★ Very High  | 2024          |
| Google     | ★★★★☆ High       | 2024          |
| Meta       | ★★★★☆ High       | 2024          |
| Microsoft  | ★★★☆☆ Medium     | 2023          |
| Bloomberg  | ★★★☆☆ Medium     | 2024          |
| Uber       | ★★☆☆☆ Low        | 2023          |

> ⚠️ Frequency data is crowd-sourced from LeetCode Discuss, Glassdoor, and
> community interview reports. Treat as a signal, not a guarantee.

---

## DSA Concepts Used

- **Binary Search Tree** — an in-order traversal of a BST yields values in ascending order, which is the whole trick → see [`/dsa/binary_search_tree.md`](/dsa/binary_search_tree.md)
- **In-Order Tree Traversal** — left → node → right visits nodes smallest-first → see [`/dsa/tree_traversal.md`](/dsa/tree_traversal.md)
- **Stack (iterative traversal)** — the optimal version replaces recursion with an explicit stack, enabling early exit and a pausable iterator → see [`/dsa/stack.md`](/dsa/stack.md)

---

## Approaches Overview

| # | Approach | Time | Space | When to use |
|---|----------|------|-------|-------------|
| 1 | In-order to full list | O(n) | O(n) | Simplest; when memory is not a concern |
| 2 | In-order early stop (recursive) | O(h + k) | O(h) | Avoid visiting nodes past the k-th |
| 3 | Iterative in-order (stack) | O(h + k) | O(h) | Optimal; no list, supports pausable iterator |

---

## Approach 1 — In-Order Traversal to Full List

### Intuition
An in-order traversal of a BST visits nodes in ascending value order. So the
k-th smallest is exactly the k-th element that traversal produces. The simplest
form materializes the entire sorted list and indexes into it at position `k-1`.

### Algorithm
1. Recursively traverse in-order (left, node, right), appending each value to a list.
2. Return `list[k-1]` (1-indexed `k` maps to 0-indexed slice).

### Complexity
- **Time:** O(n) — every node is visited once.
- **Space:** O(n) — the full value list, plus O(h) recursion stack.

### Code
```go
func inorderFullList(root *TreeNode, k int) int {
	var vals []int
	var inorder func(node *TreeNode)
	inorder = func(node *TreeNode) {
		if node == nil {
			return
		}
		inorder(node.Left)            // all smaller values first
		vals = append(vals, node.Val) // then this node (ascending order)
		inorder(node.Right)           // then all larger values
	}
	inorder(root)
	return vals[k-1] // k is 1-indexed; slice is 0-indexed
}
```

### Dry Run
Example 1, tree `[3,1,4,null,2]`, `k = 1`. In-order visitation:

| Step | Node visited | vals appended |
|------|--------------|---------------|
| 1    | 1            | [1]           |
| 2    | 2            | [1,2]         |
| 3    | 3            | [1,2,3]       |
| 4    | 4            | [1,2,3,4]     |

`vals[k-1] = vals[0] = 1`. ✅

---

## Approach 2 — In-Order with Early Stop (Counter)

### Intuition
We don't need the whole sorted list — only its k-th element. Count nodes as the
in-order walk emits them; when the count reaches `k`, that node's value is the
answer and everything to its right is irrelevant. This trims work to O(h + k).

### Algorithm
1. Keep a running `count` and an `ans` holder plus a `found` flag.
2. Recurse in-order; on visiting a node, increment `count`.
3. If `count == k`, record `ans = node.Val`, set `found`, and short-circuit.
4. Return `ans`.

### Complexity
- **Time:** O(h + k) — descend to the smallest (O(h)) then emit k nodes.
- **Space:** O(h) — recursion stack depth.

### Code
```go
func inorderEarlyStop(root *TreeNode, k int) int {
	count := 0 // how many nodes emitted so far
	ans := 0   // the k-th smallest value once found
	found := false

	var inorder func(node *TreeNode)
	inorder = func(node *TreeNode) {
		if node == nil || found {
			return // stop descending once the answer is fixed
		}
		inorder(node.Left) // smaller values first
		if found {
			return // answer found in the left subtree — unwind
		}
		count++ // this node is the next value in ascending order
		if count == k {
			ans = node.Val // k-th smallest reached
			found = true
			return
		}
		inorder(node.Right) // otherwise continue into larger values
	}
	inorder(root)
	return ans
}
```

### Dry Run
Example 1, tree `[3,1,4,null,2]`, `k = 1`:

| Node entered | Left done? | count after visit | count == k? | Action |
|--------------|-----------|-------------------|-------------|--------|
| 3            | descend left first | — | — | recurse into 1 |
| 1            | left is nil | 1 | yes (1==1) | ans=1, found=true |

Recursion unwinds immediately. Returns `1`. ✅

---

## Approach 3 — Iterative In-Order (Stack, Optimal)

### Intuition
The iterative in-order pattern pushes left spines onto a stack, then pops to emit
nodes in ascending order. Emitting one node at a time lets us stop exactly at the
k-th — no recursion, no list. This control flow also underlies the follow-up: an
iterator you can pause and resume, ideal when the tree is modified frequently.

### Algorithm
1. Start with an empty stack and `curr = root`.
2. Push `curr` and all its left descendants (go as far left as possible).
3. Pop a node — it is the next in ascending order. Decrement `k`; if `k == 0`,
   return its value.
4. Set `curr = popped.Right` and repeat.

### Complexity
- **Time:** O(h + k) — O(h) to reach the smallest, then k pops.
- **Space:** O(h) — the stack holds at most one root-to-leaf path.

### Code
```go
func iterativeInorder(root *TreeNode, k int) int {
	stack := []*TreeNode{}
	curr := root
	for curr != nil || len(stack) > 0 {
		for curr != nil { // dive to the leftmost unvisited node
			stack = append(stack, curr)
			curr = curr.Left
		}
		// Pop: this node is the next smallest not yet emitted.
		curr = stack[len(stack)-1]
		stack = stack[:len(stack)-1]
		k--
		if k == 0 {
			return curr.Val // exactly the k-th smallest
		}
		curr = curr.Right // explore the right subtree (larger values)
	}
	return -1 // unreachable given valid input (1 ≤ k ≤ n)
}
```

### Dry Run
Example 1, tree `[3,1,4,null,2]`, `k = 1`:

| Step | Action | stack | curr | k |
|------|--------|-------|------|---|
| 1    | push 3, then push 1 (left of 3) | [3,1] | nil | 1 |
| 2    | pop 1 (leftmost) | [3] | 1 | 1 |
| 3    | k-- → 0, k==0 → return 1 | [3] | 1 | 0 |

Returns `1`. ✅

---

## Key Takeaways
- In-order traversal of a BST = sorted order; the k-th emitted node is the answer.
- Early-stopping (recursive counter or iterative stack) cuts O(n) down to
  O(h + k) by never touching nodes past the k-th.
- The iterative stack version is the interview-favorite because it exits early,
  uses O(h) space, and generalizes to a resumable iterator.
- **Follow-up:** for frequent insert/delete + kth queries, augment each node with
  a subtree-size count. Then kth-smallest is an O(h) descent (compare k against
  the left subtree size), and inserts/deletes maintain counts in O(h).

---

## Related Problems
- LeetCode #94 — Binary Tree Inorder Traversal (the core traversal)
- LeetCode #173 — Binary Search Tree Iterator (pausable in-order iterator)
- LeetCode #108 — Convert Sorted Array to BST (BST ↔ sorted order)
- LeetCode #235 — Lowest Common Ancestor of a BST (BST ordering property)
