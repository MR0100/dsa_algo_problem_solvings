# 0257 — Binary Tree Paths

> LeetCode #257 · Difficulty: Easy
> **Categories:** Tree, Depth-First Search, String, Backtracking, Binary Tree

---

## Problem Statement

Given the `root` of a binary tree, return *all root-to-leaf paths in **any order***.

A **leaf** is a node with no children.

**Example 1:**

```
Input: root = [1,2,3,null,5]
Output: ["1->2->5","1->3"]
```

**Example 2:**

```
Input: root = [1]
Output: ["1"]
```

**Constraints:**

- The number of nodes in the tree is in the range `[1, 100]`.
- `-100 <= Node.val <= 100`

---

## Company Frequency

| Company    | Frequency        | Last Reported |
|------------|------------------|---------------|
| Amazon     | ★★★★☆ High       | 2024          |
| Google     | ★★★☆☆ Medium     | 2024          |
| Facebook   | ★★★☆☆ Medium     | 2023          |
| Microsoft  | ★★☆☆☆ Low        | 2023          |
| Apple      | ★★☆☆☆ Low        | 2022          |

> ⚠️ Frequency data is crowd-sourced from LeetCode Discuss, Glassdoor, and
> community interview reports. Treat as a signal, not a guarantee.

---

## DSA Concepts Used

- **Tree Traversal (DFS)** — a root-to-leaf path is a full descent of the tree; DFS naturally builds each path as it walks down → see [`/dsa/tree_traversal.md`](/dsa/tree_traversal.md)
- **Backtracking** — Approach 2 keeps a single shared path slice, pushing on entry and popping on exit to reuse memory across branches → see [`/dsa/backtracking.md`](/dsa/backtracking.md)
- **Stack** — Approach 3 replaces the recursion with an explicit stack of (node, path) frames → see [`/dsa/stack.md`](/dsa/stack.md)

---

## Approaches Overview

| # | Approach | Time | Space | When to use |
|---|----------|------|-------|-------------|
| 1 | DFS String Accumulation | O(n·h) | O(n·h) | Shortest to write; immutable path strings, no undo needed |
| 2 | DFS Backtracking Slice | O(n·h) | O(h) aux | Reuses one path buffer; the canonical backtracking template |
| 3 | Iterative DFS Explicit Stack (Optimal) | O(n·h) | O(n·h) | When recursion is unwanted / stack depth is a concern |

*(h = tree height; each of up to O(n) leaves produces a path of length O(h).)*

---

## Approach 1 — DFS String Accumulation

### Intuition

A root-to-leaf path is exactly the sequence of node values from the root down to a node with no children. Carry the partial path (as a string) down the recursion; whenever you reach a leaf, the accumulated string is a finished answer. Because strings are immutable in Go, each recursive call gets its own copy — no cleanup needed.

### Algorithm

1. `dfs(node, path)`: append `node.Val` to `path`.
2. If `node` is a leaf (both children nil), append `path` to results and return.
3. Otherwise recurse into each non-nil child with `path + "->"`.

### Complexity

- **Time:** O(n·h) — `n` nodes visited; assembling up to O(n) paths of length O(h).
- **Space:** O(n·h) — total output size; recursion stack is O(h).

### Code

```go
func dfsStringAccum(root *TreeNode) []string {
	var result []string
	if root == nil {
		return result // empty tree → no paths
	}
	var dfs func(node *TreeNode, path string)
	dfs = func(node *TreeNode, path string) {
		// Append this node's value to the path built so far.
		path += strconv.Itoa(node.Val)
		if node.Left == nil && node.Right == nil { // leaf → path is complete
			result = append(result, path)
			return
		}
		if node.Left != nil {
			dfs(node.Left, path+"->") // extend with arrow then recurse
		}
		if node.Right != nil {
			dfs(node.Right, path+"->")
		}
	}
	dfs(root, "")
	return result
}
```

### Dry Run

Example 1: tree `[1,2,3,null,5]` (node 2 has right child 5).

| Call | path in | leaf? | action |
|------|---------|-------|--------|
| dfs(1, "") | "1" | no | recurse left(2), right(3) |
| dfs(2, "1->") | "1->2" | no | recurse right(5) |
| dfs(5, "1->2->") | "1->2->5" | yes | append "1->2->5" |
| dfs(3, "1->") | "1->3" | yes | append "1->3" |

Result: `["1->2->5", "1->3"]` ✔

---

## Approach 2 — DFS Backtracking Slice

### Intuition

Rather than copy a growing string at every level, keep one slice of the current path's values. Push the current node's value on entry, pop it on exit — so the same buffer is reused across all branches. At a leaf, join the slice with `"->"` to snapshot the path.

### Algorithm

1. `dfs(node)`: push `node.Val` onto the shared `path` slice.
2. If `node` is a leaf, join `path` with `"->"` and record it.
3. Else recurse into non-nil children.
4. Pop the value before returning (restore state for the sibling/parent).

### Complexity

- **Time:** O(n·h) — up to O(n) leaves, each joining an O(h)-length slice.
- **Space:** O(h) auxiliary — the path slice + recursion depth (output excluded).

### Code

```go
func dfsBacktrack(root *TreeNode) []string {
	var result []string
	if root == nil {
		return result
	}
	path := []string{} // current root-to-node values as strings
	var dfs func(node *TreeNode)
	dfs = func(node *TreeNode) {
		path = append(path, strconv.Itoa(node.Val)) // push current value
		if node.Left == nil && node.Right == nil {   // leaf
			result = append(result, strings.Join(path, "->")) // snapshot the path
		} else {
			if node.Left != nil {
				dfs(node.Left)
			}
			if node.Right != nil {
				dfs(node.Right)
			}
		}
		path = path[:len(path)-1] // pop: undo this node before returning
	}
	dfs(root)
	return result
}
```

### Dry Run

Example 1: tree `[1,2,3,null,5]`.

| Step | action | path slice after |
|------|--------|------------------|
| 1 | enter 1, push | ["1"] |
| 2 | enter 2, push | ["1","2"] |
| 3 | enter 5, push (leaf) → record "1->2->5" | ["1","2","5"] |
| 4 | pop 5 | ["1","2"] |
| 5 | pop 2 | ["1"] |
| 6 | enter 3, push (leaf) → record "1->3" | ["1","3"] |
| 7 | pop 3, pop 1 | [] |

Result: `["1->2->5", "1->3"]` ✔

---

## Approach 3 — Iterative DFS with Explicit Stack (Optimal)

### Intuition

Any recursive DFS becomes iterative by carrying the "call state" on a stack. The state each node needs is the path string ending at it. Pop a node with its path; if it is a leaf, emit; otherwise push its children with the extended path. Pushing the right child before the left makes the left pop first, matching a left-to-right recursive order.

### Algorithm

1. Push `(root, "root value")` onto the stack.
2. While the stack is non-empty, pop `(node, path)`.
3. If `node` is a leaf, add `path` to results.
4. Else push children (right then left) with `path + "->" + childVal`.

### Complexity

- **Time:** O(n·h) — same total work as the recursive versions.
- **Space:** O(n·h) — the stack can hold O(n) node/path frames.

### Code

```go
func iterativeStack(root *TreeNode) []string {
	var result []string
	if root == nil {
		return result
	}
	type frame struct {
		node *TreeNode
		path string
	}
	// Seed the stack with the root and its own value as the starting path.
	stack := []frame{{root, strconv.Itoa(root.Val)}}
	for len(stack) > 0 {
		top := stack[len(stack)-1] // peek
		stack = stack[:len(stack)-1] // pop
		node, path := top.node, top.path
		if node.Left == nil && node.Right == nil { // leaf → complete path
			result = append(result, path)
			continue
		}
		// Push RIGHT first so LEFT is popped first (stack is LIFO); this makes
		// the output order match a left-to-right recursive DFS.
		if node.Right != nil {
			stack = append(stack, frame{node.Right, path + "->" + strconv.Itoa(node.Right.Val)})
		}
		if node.Left != nil {
			stack = append(stack, frame{node.Left, path + "->" + strconv.Itoa(node.Left.Val)})
		}
	}
	return result
}
```

### Dry Run

Example 1: tree `[1,2,3,null,5]`.

| Step | popped (node, path) | leaf? | stack after |
|------|---------------------|-------|-------------|
| 0 | — (seed) | — | [(1,"1")] |
| 1 | (1,"1") | no | [(3,"1->3"), (2,"1->2")] |
| 2 | (2,"1->2") | no | [(3,"1->3"), (5,"1->2->5")] |
| 3 | (5,"1->2->5") | yes → emit | [(3,"1->3")] |
| 4 | (3,"1->3") | yes → emit | [] |

Result: `["1->2->5", "1->3"]` ✔

---

## Key Takeaways

- **Root-to-leaf enumeration is a DFS where the "answer" is recorded at leaves**, not at the root. The trigger condition is `left == nil && right == nil`.
- **Immutable strings vs. backtracking slice**: passing a fresh string per call is simplest and needs no undo; a shared slice with push/pop saves allocation and is the reusable backtracking template.
- **Recursion ↔ explicit stack**: bundle whatever the recursive call carried (here the path) into the stack frame. Push children in reverse of the desired visit order because a stack is LIFO.
- LeetCode accepts paths in **any order**, so worrying about order is only for matching a specific expected string — not correctness.

---

## Related Problems

- LeetCode #112 — Path Sum (root-to-leaf DFS with a target)
- LeetCode #113 — Path Sum II (collect all root-to-leaf paths hitting a sum)
- LeetCode #124 — Binary Tree Maximum Path Sum (arbitrary path, not root-to-leaf)
- LeetCode #988 — Smallest String Starting From Leaf (path building leaf-to-root)
