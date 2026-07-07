# 0144 — Binary Tree Preorder Traversal

> LeetCode #144 · Difficulty: Easy
> **Categories:** Stack, Tree, Depth-First Search, Binary Tree

---

## Problem Statement

Given the `root` of a binary tree, return *the preorder traversal of its nodes' values*.

**Example 1:**
```
Input: root = [1,null,2,3]
Output: [1,2,3]
Explanation:
  1
   \
    2
   /
  3
```

**Example 2:**
```
Input: root = [1,2,3,4,5,null,8,null,null,6,7,9]
Output: [1,2,4,5,6,7,3,8,9]
Explanation:
        1
      /   \
     2     3
    / \     \
   4   5     8
      / \   /
     6   7 9
```

**Example 3:**
```
Input: root = []
Output: []
```

**Example 4:**
```
Input: root = [1]
Output: [1]
```

**Constraints:**
- The number of nodes in the tree is in the range `[0, 100]`.
- `-100 <= Node.val <= 100`

**Follow-up:** Recursive solution is trivial, could you do it iteratively?

---

## Company Frequency

| Company    | Frequency        | Last Reported |
|------------|------------------|---------------|
| Amazon     | ★★★★☆ High       | 2024          |
| Microsoft  | ★★★★☆ High       | 2024          |
| Google     | ★★★☆☆ Medium     | 2024          |
| Meta       | ★★★☆☆ Medium     | 2023          |
| Apple      | ★★☆☆☆ Low        | 2023          |
| Bloomberg  | ★★☆☆☆ Low        | 2023          |

> ⚠️ Frequency data is crowd-sourced from LeetCode Discuss, Glassdoor, and
> community interview reports. Treat as a signal, not a guarantee.

---

## DSA Concepts Used

- **Tree Traversal (DFS)** — preorder = root, left, right; one of the three fundamental depth-first orders → see [`/dsa/tree_traversal.md`](/dsa/tree_traversal.md)
- **Stack** — the iterative version simulates the call stack explicitly → see [`/dsa/stack.md`](/dsa/stack.md)
- **Morris Threading** — O(1)-space traversal by temporarily re-pointing predecessor `Right` links (covered with tree traversal → see [`/dsa/tree_traversal.md`](/dsa/tree_traversal.md))

---

## Approaches Overview

| # | Approach | Time | Space | When to use |
|---|----------|------|-------|-------------|
| 1 | Recursive DFS | O(n) | O(h) | Default; clearest possible code |
| 2 | Iterative (Stack) | O(n) | O(h) | The follow-up; required when recursion depth is a risk |
| 3 | Morris Traversal | O(n) | O(1) | Space-critical settings; may temporarily mutate the tree |

---

## Approach 1 — Recursive DFS

### Intuition
Preorder is defined recursively — visit the **root**, then traverse the **left** subtree, then the **right** subtree. Translating the definition into a helper function is the entire solution; the call stack remembers where to resume after each subtree.

### Algorithm
1. Define `dfs(node)`:
   1. If `node == nil` → return.
   2. Append `node.Val` to the result (root first).
   3. `dfs(node.Left)`.
   4. `dfs(node.Right)`.
2. Call `dfs(root)`, return the result.

### Complexity
- **Time:** O(n) — each of the n nodes is visited exactly once.
- **Space:** O(h) — deepest recursion equals the tree height: O(log n) balanced, O(n) skewed.

### Code
```go
func preorderRecursive(root *TreeNode) []int {
	result := []int{} // non-nil so an empty tree prints as [] not nil
	var dfs func(node *TreeNode)
	dfs = func(node *TreeNode) {
		if node == nil {
			return // empty subtree contributes nothing
		}
		result = append(result, node.Val) // ROOT first
		dfs(node.Left)                    // then the whole LEFT subtree
		dfs(node.Right)                   // then the whole RIGHT subtree
	}
	dfs(root)
	return result
}
```

### Dry Run (Example 1: root = [1,null,2,3])

Tree: `1` has no left child, right child `2`; `2` has left child `3`.

| Step | Call | Action | `result` |
|------|------|--------|----------|
| 1 | dfs(1) | visit 1 | [1] |
| 2 | dfs(1.Left = nil) | return | [1] |
| 3 | dfs(1.Right = 2) | visit 2 | [1, 2] |
| 4 | dfs(2.Left = 3) | visit 3 | [1, 2, 3] |
| 5 | dfs(3.Left = nil) | return | [1, 2, 3] |
| 6 | dfs(3.Right = nil) | return | [1, 2, 3] |
| 7 | dfs(2.Right = nil) | return | [1, 2, 3] |

Output: `[1,2,3]` ✓

---

## Approach 2 — Iterative (Stack)

### Intuition
Replace the call stack with your own. A stack is LIFO, and we visit a node the moment we pop it; to make the left subtree come out before the right one, push the **right child first and the left child second** — LIFO order flips them back.

### Algorithm
1. If `root == nil` → return `[]`.
2. `stack = [root]`.
3. While the stack is non-empty:
   1. Pop `node`; append `node.Val`.
   2. If `node.Right != nil` → push it (popped last).
   3. If `node.Left != nil` → push it (popped next).
4. Return the result.

### Complexity
- **Time:** O(n) — every node is pushed once and popped once.
- **Space:** O(h) in the balanced case; up to O(n) worst case for the stored frontier.

### Code
```go
func preorderIterative(root *TreeNode) []int {
	result := []int{}
	if root == nil {
		return result // nothing to traverse
	}
	stack := []*TreeNode{root}
	for len(stack) > 0 {
		node := stack[len(stack)-1]       // peek top
		stack = stack[:len(stack)-1]      // pop it
		result = append(result, node.Val) // visit on pop → root before children
		if node.Right != nil {
			stack = append(stack, node.Right) // pushed FIRST → popped LAST
		}
		if node.Left != nil {
			stack = append(stack, node.Left) // pushed LAST → popped NEXT (left first)
		}
	}
	return result
}
```

### Dry Run (Example 1: root = [1,null,2,3])

| Step | Popped | Pushed | Stack after | `result` |
|------|--------|--------|-------------|----------|
| init | — | 1 | [1] | [] |
| 1 | 1 | right=2 (no left) | [2] | [1] |
| 2 | 2 | left=3 (no right) | [3] | [1, 2] |
| 3 | 3 | (leaf — nothing) | [] | [1, 2, 3] |

Stack empty → output `[1,2,3]` ✓

---

## Approach 3 — Morris Traversal (Optimal Space)

### Intuition
The stack exists solely to remember "come back here after the left subtree". Morris threading stores that return address **inside the tree**: the inorder predecessor of `curr` (rightmost node of `curr.Left`) has a free `nil` Right pointer — temporarily aim it back at `curr`. The key preorder twist: emit the node on the **first** arrival (before descending left), whereas inorder Morris emits on the second arrival.

### Algorithm
1. `curr = root`.
2. While `curr != nil`:
   1. If `curr.Left == nil` → visit `curr`, move to `curr.Right`.
   2. Else find `pred` = rightmost node of `curr.Left`, stopping if `pred.Right == curr` (an existing thread).
   3. If `pred.Right == nil` (first arrival) → **visit `curr` now**, set `pred.Right = curr` (thread), descend `curr = curr.Left`.
   4. Else (second arrival via thread) → unset `pred.Right = nil` (restore), move `curr = curr.Right`.

### Complexity
- **Time:** O(n) — every edge is traversed at most twice (once laying the thread, once removing it), so the predecessor searches sum to O(n) overall.
- **Space:** O(1) — no stack or recursion; all threads are removed, so the tree is restored.

### Code
```go
func preorderMorris(root *TreeNode) []int {
	result := []int{}
	curr := root
	for curr != nil {
		if curr.Left == nil {
			result = append(result, curr.Val) // leaf-ward: visit and move on
			curr = curr.Right
			continue
		}
		// Find the inorder predecessor: rightmost node in the left subtree,
		// stopping early if we hit a thread pointing back to curr.
		pred := curr.Left
		for pred.Right != nil && pred.Right != curr {
			pred = pred.Right
		}
		if pred.Right == nil {
			result = append(result, curr.Val) // FIRST arrival → preorder visit
			pred.Right = curr                 // lay the return thread
			curr = curr.Left                  // descend into the left subtree
		} else {
			pred.Right = nil  // second arrival → remove thread (restore tree)
			curr = curr.Right // left subtree done; continue rightward
		}
	}
	return result
}
```

### Dry Run (Example 1: root = [1,null,2,3])

Tree: `1 → right 2`, `2 → left 3`.

| Step | `curr` | `curr.Left` | Predecessor logic | Action | `result` |
|------|--------|-------------|-------------------|--------|----------|
| 1 | 1 | nil | — | visit 1; go right | [1] |
| 2 | 2 | 3 | pred = 3, `3.Right == nil` | visit 2; thread 3.Right = 2; go left | [1, 2] |
| 3 | 3 | nil | — | visit 3; go right (follows thread → 2) | [1, 2, 3] |
| 4 | 2 | 3 | pred = 3, `3.Right == 2` (thread) | unthread 3.Right = nil; go right (nil) | [1, 2, 3] |
| 5 | nil | — | — | loop ends | [1, 2, 3] |

Output: `[1,2,3]` ✓ — and the tree is back to its original shape.

---

## Key Takeaways

- Preorder's iterative version is the **easiest of the three DFS orders**: visit-on-pop plus "push right before left" is the whole trick.
- **Morris preorder vs Morris inorder differ by one line**: preorder emits on the *first* arrival (before threading left), inorder emits on the *second* arrival. Same threading machinery.
- Preorder is the natural order for **copying/serialising** a tree (root appears before its subtrees — LeetCode #297 uses it).
- Return `[]int{}` rather than `nil` when the tree is empty if the caller expects `[]` — a small Go-specific detail for clean output.
- The visit-on-pop stack pattern generalises to n-ary trees (#589) by pushing children in reverse order.

---

## Related Problems

- LeetCode #94 — Binary Tree Inorder Traversal (same trio: recursive / stack / Morris)
- LeetCode #145 — Binary Tree Postorder Traversal (the hard sibling)
- LeetCode #589 — N-ary Tree Preorder Traversal (same pattern, k children)
- LeetCode #102 — Binary Tree Level Order Traversal (BFS counterpart)
- LeetCode #105 — Construct Binary Tree from Preorder and Inorder Traversal (uses preorder's root-first property)
- LeetCode #297 — Serialize and Deserialize Binary Tree (preorder serialisation)
