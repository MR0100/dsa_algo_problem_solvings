# 0226 — Invert Binary Tree

> LeetCode #226 · Difficulty: Easy
> **Categories:** Tree, Depth-First Search, Breadth-First Search, Binary Tree

---

## Problem Statement

Given the `root` of a binary tree, invert the tree, and return *its root*.

**Example 1:**
```
Input: root = [4,2,7,1,3,6,9]
Output: [4,7,2,9,6,3,1]
```

**Example 2:**
```
Input: root = [2,1,3]
Output: [2,3,1]
```

**Example 3:**
```
Input: root = []
Output: []
```

**Constraints:**
- The number of nodes in the tree is in the range `[0, 100]`.
- `-100 <= Node.val <= 100`

---

## Company Frequency

| Company    | Frequency        | Last Reported |
|------------|------------------|---------------|
| Google     | ★★★★★ Very High  | 2024          |
| Amazon     | ★★★★☆ High       | 2024          |
| Microsoft  | ★★★★☆ High       | 2024          |
| Meta       | ★★★☆☆ Medium     | 2023          |
| Apple      | ★★★☆☆ Medium     | 2023          |
| Bloomberg  | ★★☆☆☆ Low        | 2023          |

> ⚠️ Frequency data is crowd-sourced from LeetCode Discuss, Glassdoor, and
> community interview reports. Treat as a signal, not a guarantee.

---

## DSA Concepts Used

- **Tree Traversal** — inverting requires visiting every node exactly once, the defining property of a full traversal → see [`/dsa/tree_traversal.md`](/dsa/tree_traversal.md)
- **Depth-First Search** — the recursive and stack-based solutions both walk the tree depth-first → see [`/dsa/graph_bfs_dfs.md`](/dsa/graph_bfs_dfs.md)
- **Breadth-First Search** — the queue-based solution mirrors the tree level by level → see [`/dsa/graph_bfs_dfs.md`](/dsa/graph_bfs_dfs.md)
- **Queue / Stack** — the iterative variants use an explicit FIFO/LIFO to replace the call stack → see [`/dsa/queue_deque.md`](/dsa/queue_deque.md)

---

## Approaches Overview

| # | Approach | Time | Space | When to use |
|---|----------|------|-------|-------------|
| 1 | Recursive DFS | O(n) | O(h) | Cleanest, default choice; fine for the tiny constraint |
| 2 | Iterative BFS (queue) | O(n) | O(w) | Avoid recursion; process level by level |
| 3 | Iterative DFS (stack) | O(n) | O(h) | Avoid recursion while keeping depth-first order |

---

## Approach 1 — Recursive DFS (Swap Children)

### Intuition
Inverting a tree is producing its mirror image. At every node, the entire left
subtree and the entire right subtree simply trade places. If we perform that
single swap at every node, the whole tree ends up mirrored. Recursion states
this naturally: invert the children first, then swap the two inverted subtrees.

### Algorithm
1. If the current node is `nil`, return `nil` (an empty subtree is its own mirror).
2. Recursively invert the left subtree.
3. Recursively invert the right subtree.
4. Swap the node's `Left` and `Right` pointers.
5. Return the node.

### Complexity
- **Time:** O(n) — each of the `n` nodes is visited once and does O(1) work.
- **Space:** O(h) — recursion stack depth equals the tree height `h` (O(n) for a degenerate tree, O(log n) when balanced).

### Code
```go
func invertRecursive(root *TreeNode) *TreeNode {
	if root == nil { // empty subtree — nothing to mirror
		return nil
	}
	// Invert both sides first, then swap the returned (inverted) subtrees.
	left := invertRecursive(root.Left)   // fully-inverted left subtree
	right := invertRecursive(root.Right) // fully-inverted right subtree
	root.Left = right                    // right subtree now hangs on the left
	root.Right = left                    // left subtree now hangs on the right
	return root
}
```

### Dry Run
Tree `[4,2,7,1,3,6,9]` — recursion returns from the leaves upward:

| Call on node | Left after invert | Right after invert | Node after swap |
|--------------|-------------------|--------------------|-----------------|
| 1 (leaf)     | nil               | nil                | 1               |
| 3 (leaf)     | nil               | nil                | 3               |
| 2            | inverted 1        | inverted 3         | Left=3, Right=1 |
| 6 (leaf)     | nil               | nil                | 6               |
| 9 (leaf)     | nil               | nil                | 9               |
| 7            | inverted 6        | inverted 9         | Left=9, Right=6 |
| 4 (root)     | inverted 2        | inverted 7         | Left=7, Right=2 |

Final level-order: `[4,7,2,9,6,3,1]`. ✅

---

## Approach 2 — Iterative BFS (Queue)

### Intuition
The swap at each node is completely independent of every other node's swap, so
the *order* in which we visit nodes is irrelevant — we only need to visit them
all. A queue gives a level-order sweep with no recursion, which removes any
deep-stack concern.

### Algorithm
1. If `root` is `nil`, return `nil`.
2. Push `root` into a queue.
3. While the queue is non-empty:
   a. Dequeue a node.
   b. Swap its `Left` and `Right` children.
   c. Enqueue each non-nil child.
4. Return `root`.

### Complexity
- **Time:** O(n) — each node enters and leaves the queue exactly once.
- **Space:** O(w) — the queue holds at most one full level, up to O(n) in the worst case.

### Code
```go
func invertBFS(root *TreeNode) *TreeNode {
	if root == nil {
		return nil
	}
	queue := []*TreeNode{root} // FIFO of nodes still needing their swap
	for len(queue) > 0 {
		node := queue[0]  // dequeue the front node
		queue = queue[1:] // pop it off the queue
		// Swap this node's children — the core mirror operation.
		node.Left, node.Right = node.Right, node.Left
		if node.Left != nil { // enqueue children (post-swap) for later processing
			queue = append(queue, node.Left)
		}
		if node.Right != nil {
			queue = append(queue, node.Right)
		}
	}
	return root
}
```

### Dry Run
Tree `[4,2,7,1,3,6,9]`:

| Step | Dequeued | After swap (Left,Right) | Queue after enqueue |
|------|----------|-------------------------|---------------------|
| 1    | 4        | (7,2)                   | [7, 2]              |
| 2    | 7        | (9,6)                   | [2, 9, 6]           |
| 3    | 2        | (3,1)                   | [9, 6, 3, 1]        |
| 4    | 9        | (nil,nil)               | [6, 3, 1]           |
| 5    | 6        | (nil,nil)               | [3, 1]              |
| 6    | 3        | (nil,nil)               | [1]                 |
| 7    | 1        | (nil,nil)               | []                  |

Final level-order: `[4,7,2,9,6,3,1]`. ✅

---

## Approach 3 — Iterative DFS (Stack)

### Intuition
Same independent-swap observation as BFS, but a LIFO stack reproduces the
depth-first visiting order of the recursive version *without* recursion —
handy when you want to avoid growing the call stack or hitting recursion limits.

### Algorithm
1. If `root` is `nil`, return `nil`.
2. Push `root` onto a stack.
3. While the stack is non-empty:
   a. Pop a node.
   b. Swap its `Left` and `Right` children.
   c. Push each non-nil child.
4. Return `root`.

### Complexity
- **Time:** O(n) — every node is pushed and popped once.
- **Space:** O(h) — the stack holds at most a root-to-leaf path plus siblings.

### Code
```go
func invertDFSStack(root *TreeNode) *TreeNode {
	if root == nil {
		return nil
	}
	stack := []*TreeNode{root} // LIFO of nodes awaiting their swap
	for len(stack) > 0 {
		node := stack[len(stack)-1]  // peek the top
		stack = stack[:len(stack)-1] // pop it
		// Swap children — the mirror operation.
		node.Left, node.Right = node.Right, node.Left
		if node.Left != nil { // push children to keep going deeper
			stack = append(stack, node.Left)
		}
		if node.Right != nil {
			stack = append(stack, node.Right)
		}
	}
	return root
}
```

### Dry Run
Tree `[4,2,7,1,3,6,9]`:

| Step | Popped | After swap (Left,Right) | Stack after push |
|------|--------|-------------------------|------------------|
| 1    | 4      | (7,2)                   | [7, 2]           |
| 2    | 2      | (3,1)                   | [7, 3, 1]        |
| 3    | 1      | (nil,nil)               | [7, 3]           |
| 4    | 3      | (nil,nil)               | [7]              |
| 5    | 7      | (9,6)                   | [9, 6]           |
| 6    | 6      | (nil,nil)               | [9]              |
| 7    | 9      | (nil,nil)               | []               |

Final level-order: `[4,7,2,9,6,3,1]`. ✅

---

## Key Takeaways
- Inverting a binary tree = swapping every node's two children; the swap is
  local and order-independent, so any traversal (pre/post-order, BFS, DFS) works.
- The recursive one-liner is the canonical answer, but knowing the iterative
  queue/stack versions demonstrates you can avoid recursion on demand.
- "Order-independent per-node operation" is a recurring signal that BFS and DFS
  are interchangeable — pick whichever fits the constraints.

---

## Related Problems
- LeetCode #101 — Symmetric Tree (mirror comparison instead of mutation)
- LeetCode #100 — Same Tree (structural tree recursion)
- LeetCode #104 — Maximum Depth of Binary Tree (same traversal skeleton)
- LeetCode #543 — Diameter of Binary Tree (post-order aggregation)
