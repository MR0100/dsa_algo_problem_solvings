# 0145 — Binary Tree Postorder Traversal

> LeetCode #145 · Difficulty: Easy
> **Categories:** Stack, Tree, Depth-First Search, Binary Tree

---

## Problem Statement

Given the `root` of a binary tree, return *the postorder traversal of its nodes' values*.

**Example 1:**
```
Input: root = [1,null,2,3]
Output: [3,2,1]
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
Output: [4,6,7,5,2,9,8,3,1]
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
- The number of the nodes in the tree is in the range `[0, 100]`.
- `-100 <= Node.val <= 100`

**Follow-up:** Recursive solution is trivial, could you do it iteratively?

---

## Company Frequency

| Company    | Frequency        | Last Reported |
|------------|------------------|---------------|
| Amazon     | ★★★★☆ High       | 2024          |
| Microsoft  | ★★★☆☆ Medium     | 2024          |
| Google     | ★★★☆☆ Medium     | 2024          |
| Meta       | ★★★☆☆ Medium     | 2023          |
| Apple      | ★★☆☆☆ Low        | 2023          |
| LinkedIn   | ★★☆☆☆ Low        | 2023          |

> ⚠️ Frequency data is crowd-sourced from LeetCode Discuss, Glassdoor, and
> community interview reports. Treat as a signal, not a guarantee.

---

## DSA Concepts Used

- **Tree Traversal (DFS)** — postorder = left, right, root; children strictly before parents → see [`/dsa/tree_traversal.md`](/dsa/tree_traversal.md)
- **Stack** — both iterative variants replace the call stack; the last-visited trick is postorder-specific → see [`/dsa/stack.md`](/dsa/stack.md)
- **Morris Threading** — O(1)-space traversal via temporary predecessor threads (see the traversal notes in [`/dsa/tree_traversal.md`](/dsa/tree_traversal.md))
- **Linked List Reversal** — Morris postorder reverses right-pointer spines in place, exactly like reversing a linked list → see [`/dsa/linked_list.md`](/dsa/linked_list.md)

---

## Approaches Overview

| # | Approach | Time | Space | When to use |
|---|----------|------|-------|-------------|
| 1 | Recursive DFS | O(n) | O(h) | Default; clearest code |
| 2 | Reversed Modified Preorder | O(n) | O(n) | Fastest iterative version to write correctly |
| 3 | One Stack + Last Visited | O(n) | O(h) | True streaming postorder; the "real" follow-up answer |
| 4 | Morris Traversal | O(n) | O(1) | Space-critical; hardest to write, temporarily mutates the tree |

---

## Approach 1 — Recursive DFS

### Intuition
Postorder emits a node only after **both** subtrees are completely finished — the definition (left, right, root) is recursive, so the code is a direct transcription. This is the order used to safely delete a tree or compute child-dependent values (heights, subtree sums).

### Algorithm
1. Define `dfs(node)`:
   1. If `node == nil` → return.
   2. `dfs(node.Left)`.
   3. `dfs(node.Right)`.
   4. Append `node.Val`.
2. Call `dfs(root)`; return the result.

### Complexity
- **Time:** O(n) — each node is processed exactly once.
- **Space:** O(h) — recursion depth equals tree height: O(log n) balanced, O(n) skewed.

### Code
```go
func postorderRecursive(root *TreeNode) []int {
	result := []int{} // non-nil so an empty tree prints as [] not nil
	var dfs func(node *TreeNode)
	dfs = func(node *TreeNode) {
		if node == nil {
			return // empty subtree contributes nothing
		}
		dfs(node.Left)                    // finish the LEFT subtree first
		dfs(node.Right)                   // then the RIGHT subtree
		result = append(result, node.Val) // ROOT last
	}
	dfs(root)
	return result
}
```

### Dry Run (Example 1: root = [1,null,2,3])

Tree: `1` (no left, right = `2`); `2` (left = `3`, no right).

| Step | Call | Action | `result` |
|------|------|--------|----------|
| 1 | dfs(1) | recurse left (nil) then right (2) | [] |
| 2 | dfs(2) | recurse left (3) | [] |
| 3 | dfs(3) | left nil, right nil → emit 3 | [3] |
| 4 | back in dfs(2) | right nil → emit 2 | [3, 2] |
| 5 | back in dfs(1) | both subtrees done → emit 1 | [3, 2, 1] |

Output: `[3,2,1]` ✓

---

## Approach 2 — Reversed Modified Preorder

### Intuition
Postorder is `left, right, root`. Read it backwards: `root, right, left` — that is preorder with the children mirrored. Preorder is the easy stack traversal, so run the mirrored preorder (push left before right, so right pops first) and reverse the output at the end.

### Algorithm
1. If `root == nil` → return `[]`.
2. `stack = [root]`.
3. While the stack is non-empty:
   1. Pop `node`; append `node.Val` (building root-right-left).
   2. Push `node.Left` if non-nil (popped last).
   3. Push `node.Right` if non-nil (popped next).
4. Reverse the result slice in place; return it.

### Complexity
- **Time:** O(n) — n pushes + n pops + O(n) reversal.
- **Space:** O(n) — the stack can hold up to O(n) nodes for bushy trees (reversal itself is in-place).

### Code
```go
func reversedPreorder(root *TreeNode) []int {
	result := []int{}
	if root == nil {
		return result // nothing to traverse
	}
	stack := []*TreeNode{root}
	for len(stack) > 0 {
		node := stack[len(stack)-1]       // peek top
		stack = stack[:len(stack)-1]      // pop it
		result = append(result, node.Val) // building root-right-left order
		if node.Left != nil {
			stack = append(stack, node.Left) // pushed FIRST → popped LAST
		}
		if node.Right != nil {
			stack = append(stack, node.Right) // pushed LAST → popped NEXT (right first)
		}
	}
	// Reverse root-right-left into left-right-root (= postorder).
	for i, j := 0, len(result)-1; i < j; i, j = i+1, j-1 {
		result[i], result[j] = result[j], result[i]
	}
	return result
}
```

### Dry Run (Example 1: root = [1,null,2,3])

| Step | Popped | Pushed | Stack after | `result` (pre-reverse) |
|------|--------|--------|-------------|------------------------|
| init | — | 1 | [1] | [] |
| 1 | 1 | right=2 (no left) | [2] | [1] |
| 2 | 2 | left=3 (no right) | [3] | [1, 2] |
| 3 | 3 | (leaf) | [] | [1, 2, 3] |

Reverse `[1,2,3]` → `[3,2,1]` ✓

---

## Approach 3 — One Stack + Last Visited

### Intuition
A node may only be emitted after its right subtree is done — but a plain stack cannot tell "arriving at this node for the first time" from "returning after finishing its right child". Track the **last emitted node**: if the top's right child is exactly that node (or is nil), the right subtree is finished and the top may be emitted; otherwise dive into the right subtree first. This yields postorder in correct order *as a stream* — no reversal step.

### Algorithm
1. `curr = root`, empty stack, `lastVisited = nil`.
2. While `curr != nil` **or** the stack is non-empty:
   1. Slide left: while `curr != nil`, push `curr`, `curr = curr.Left`.
   2. `peek` = top of stack (do not pop yet).
   3. If `peek.Right != nil` and `peek.Right != lastVisited` → `curr = peek.Right` (right subtree still pending).
   4. Else pop and emit `peek`, set `lastVisited = peek`, leave `curr = nil`.
3. Return the result.

### Complexity
- **Time:** O(n) — each node is pushed once and popped once; each peek is O(1).
- **Space:** O(h) — the stack holds at most one root-to-leaf path.

### Code
```go
func oneStack(root *TreeNode) []int {
	result := []int{}
	stack := []*TreeNode{}
	var lastVisited *TreeNode // most recently emitted node
	curr := root

	for curr != nil || len(stack) > 0 {
		// Phase a: push the whole left spine of the current subtree.
		for curr != nil {
			stack = append(stack, curr)
			curr = curr.Left
		}
		peek := stack[len(stack)-1] // candidate for emission
		if peek.Right != nil && peek.Right != lastVisited {
			// Right subtree exists and is NOT finished yet → traverse it first.
			curr = peek.Right
		} else {
			// Right subtree is absent or already emitted → safe to emit peek.
			stack = stack[:len(stack)-1] // pop
			result = append(result, peek.Val)
			lastVisited = peek // remember so the parent knows right is done
			// curr stays nil: next loop iteration re-examines the new top.
		}
	}
	return result
}
```

### Dry Run (Example 1: root = [1,null,2,3])

| Step | `curr` at entry | Stack after left-slide | `peek` | `peek.Right` vs `lastVisited` | Action | `result` |
|------|-----------------|------------------------|--------|-------------------------------|--------|----------|
| 1 | 1 | [1] (1 has no left child) | 1 | Right = 2, last = nil → right pending | curr = 2 | [] |
| 2 | 2 | [1, 2, 3] (push 2, slide to its left child 3, 3 has no left) | 3 | Right = nil | pop + emit 3; last = 3 | [3] |
| 3 | nil | [1, 2] | 2 | Right = nil | pop + emit 2; last = 2 | [3, 2] |
| 4 | nil | [1] | 1 | Right = 2 **==** last → right finished | pop + emit 1; last = 1 | [3, 2, 1] |

Stack empty, `curr` nil → output `[3,2,1]` ✓

---

## Approach 4 — Morris Traversal (Optimal Space)

### Intuition
Postorder of a tree = for each node reached by inorder-Morris threading, when the thread brings us back (left subtree exhausted), emit the **right spine of its left child in reverse**. Hanging the whole tree as the **left child of a dummy node** makes the final right spine (root → rightmost path) come out through the same rule. Reversing a spine costs no memory: flip its `Right` pointers like reversing a linked list, walk it, flip them back.

### Algorithm
1. `dummy = &TreeNode{Left: root}`, `curr = dummy`.
2. While `curr != nil`:
   1. If `curr.Left == nil` → `curr = curr.Right`.
   2. Else find `pred` = rightmost node in `curr.Left` (stop at an existing thread to `curr`).
   3. If `pred.Right == nil` (first arrival) → thread `pred.Right = curr`; `curr = curr.Left`.
   4. Else (second arrival) → remove the thread; **emit the chain `curr.Left … pred` in reverse** (flip pointers → walk → flip back); `curr = curr.Right`.
3. Return the result.

### Complexity
- **Time:** O(n) — every edge is touched a constant number of times: threading, unthreading, and two pointer-flips during spine emission.
- **Space:** O(1) — beyond the output, only a few pointers; all threads and flips are undone, so the tree is fully restored.

### Code
```go
func morrisPostorder(root *TreeNode) []int {
	result := []int{}
	dummy := &TreeNode{Left: root} // ensures the rightmost spine of the real
	curr := dummy                  // root is also emitted via the same rule

	for curr != nil {
		if curr.Left == nil {
			curr = curr.Right // nothing to emit here; keep moving
			continue
		}
		// Find inorder predecessor of curr within its left subtree.
		pred := curr.Left
		for pred.Right != nil && pred.Right != curr {
			pred = pred.Right
		}
		if pred.Right == nil {
			pred.Right = curr // first arrival: lay the return thread
			curr = curr.Left  // and descend left
		} else {
			pred.Right = nil                           // second arrival: remove thread
			emitReverseSpine(curr.Left, pred, &result) // emit spine curr.Left…pred backwards
			curr = curr.Right                          // move on past the finished left subtree
		}
	}
	return result
}

func emitReverseSpine(from, to *TreeNode, result *[]int) {
	reverseRightChain(from) // flip: to → … → from
	// Walk from `to` back down to `from`, emitting values.
	for node := to; ; node = node.Right {
		*result = append(*result, node.Val)
		if node == from {
			break // reached the start of the original spine
		}
	}
	reverseRightChain(to) // flip back: from → … → to (tree restored)
}

func reverseRightChain(head *TreeNode) {
	var prev *TreeNode
	for node := head; node != nil; {
		next := node.Right // save onward pointer
		node.Right = prev  // flip the link
		prev = node        // advance prev
		node = next        // advance node
	}
}
```

### Dry Run (Example 1: root = [1,null,2,3])

Tree: `1 → right 2`, `2 → left 3`. Add `dummy.Left = 1`.

| Step | `curr` | `curr.Left` | Predecessor search | Action | `result` |
|------|--------|-------------|--------------------|--------|----------|
| 1 | dummy | 1 | rightmost of 1's subtree: 1 → 2 → (2.Right nil) ⇒ pred = 2 | thread 2.Right = dummy; curr = 1 | [] |
| 2 | 1 | nil | — | curr = 1.Right = 2 | [] |
| 3 | 2 | 3 | pred = 3 (3.Right nil) | thread 3.Right = 2; curr = 3 | [] |
| 4 | 3 | nil | — | curr = 3.Right = 2 (via thread) | [] |
| 5 | 2 | 3 | pred = 3, `3.Right == 2` → thread found | unthread; emit spine [3] reversed → 3; curr = 2.Right = dummy (via thread) | [3] |
| 6 | dummy | 1 | pred: 1 → 2, `2.Right == dummy` → thread found | unthread; spine 1 → 2 reversed → emit 2, 1; curr = dummy.Right = nil | [3, 2, 1] |
| 7 | nil | — | — | loop ends | [3, 2, 1] |

Output: `[3,2,1]` ✓ — all threads removed, tree restored.

---

## Key Takeaways

- **Postorder reversed = mirrored preorder** (`root, right, left`) — the single most useful traversal identity; it turns the hardest iterative traversal into the easiest one plus a reverse.
- The **last-visited pointer** is the canonical way to make one-stack postorder stream in order — the peek/pop distinction ("is the right subtree finished?") is the entire difficulty of iterative postorder.
- **Morris postorder = inorder threading + reversed right spines under a dummy root.** The dummy node trick (hang the tree as a *left* child) is what lets the root's own spine be emitted uniformly.
- Postorder is the "children before parent" order — reach for it whenever a node's answer depends on its subtrees: tree deletion, subtree sums, heights, DP on trees (#337, #124).
- Approach 2 is O(n) stack space (it doesn't mirror a root-to-leaf path), while Approach 3 is O(h) — worth mentioning when an interviewer pushes on space.

---

## Related Problems

- LeetCode #94 — Binary Tree Inorder Traversal (same trio of techniques)
- LeetCode #144 — Binary Tree Preorder Traversal (the mirrored building block of Approach 2)
- LeetCode #590 — N-ary Tree Postorder Traversal (same reverse trick, k children)
- LeetCode #104 — Maximum Depth of Binary Tree (postorder-style bottom-up computation)
- LeetCode #110 — Balanced Binary Tree (bottom-up heights = postorder)
- LeetCode #124 — Binary Tree Maximum Path Sum (postorder DP on trees)
