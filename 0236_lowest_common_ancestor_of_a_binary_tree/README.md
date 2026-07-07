# 0236 — Lowest Common Ancestor of a Binary Tree

> LeetCode #236 · Difficulty: Medium
> **Categories:** Tree, Depth-First Search, Binary Tree

---

## Problem Statement

Given a binary tree, find the lowest common ancestor (LCA) of two given nodes in the tree.

According to the [definition of LCA on Wikipedia](https://en.wikipedia.org/wiki/Lowest_common_ancestor): "The lowest common ancestor is defined between two nodes `p` and `q` as the lowest node in `T` that has both `p` and `q` as descendants (where we allow **a node to be a descendant of itself**)."

**Example 1:**

```
Input: root = [3,5,1,6,2,0,8,null,null,7,4], p = 5, q = 1
Output: 3
Explanation: The LCA of nodes 5 and 1 is 3.
```

**Example 2:**

```
Input: root = [3,5,1,6,2,0,8,null,null,7,4], p = 5, q = 4
Output: 5
Explanation: The LCA of nodes 5 and 4 is 5, since a node can be a descendant of itself according to the LCA definition.
```

**Example 3:**

```
Input: root = [1,2], p = 1, q = 2
Output: 1
```

**Constraints:**

- The number of nodes in the tree is in the range `[2, 10^5]`.
- `-10^9 <= Node.val <= 10^9`
- All `Node.val` are **unique**.
- `p != q`
- `p` and `q` will exist in the tree.

---

## Company Frequency

| Company    | Frequency        | Last Reported |
|------------|------------------|---------------|
| Amazon     | ★★★★★ Very High  | 2024          |
| Facebook   | ★★★★★ Very High  | 2024          |
| Microsoft  | ★★★★☆ High       | 2024          |
| Google     | ★★★★☆ High       | 2024          |
| Apple      | ★★★☆☆ Medium     | 2023          |
| Bloomberg  | ★★★☆☆ Medium     | 2023          |
| LinkedIn   | ★★★☆☆ Medium     | 2023          |

> ⚠️ Frequency data is crowd-sourced from LeetCode Discuss, Glassdoor, and
> community interview reports. Treat as a signal, not a guarantee.

---

## DSA Concepts Used

- **Tree Traversal (post-order DFS)** — the optimal solution is a single post-order recursion that lets each subtree report whether it contains `p` or `q` → see [`/dsa/tree_traversal.md`](/dsa/tree_traversal.md)
- **Depth-First Search** — every approach walks the tree with DFS, either to build paths or to bubble a "found" signal upward → see [`/dsa/graph_bfs_dfs.md`](/dsa/graph_bfs_dfs.md)
- **Hash Map** — the parent-pointer variant stores `node → parent` and an ancestor set for O(1) membership tests → see [`/dsa/hash_map.md`](/dsa/hash_map.md)

---

## Approaches Overview

| # | Approach | Time | Space | When to use |
|---|----------|------|-------|-------------|
| 1 | Path-to-Node + Divergence | O(n) | O(n) | Most intuitive; also gives you the actual paths if you need them |
| 2 | Single-Pass Recursion (Optimal) | O(n) | O(h) | The canonical interview answer — one traversal, tiny code |
| 3 | Parent Pointers + Ancestor Set | O(n) | O(n) | When repeated LCA queries or "climb up" access is handy |

---

## Approach 1 — Path-to-Node + Divergence

### Intuition

The LCA is the deepest node lying on **both** the root→`p` path and the root→`q` path. If we record each full path top-down, the answer is simply the last index where the two lists still agree — after that they branch into different subtrees.

### Algorithm

1. DFS from the root, building the ordered list of nodes on the path to `p` (backtrack when a branch fails).
2. Do the same for `q`.
3. Walk both lists together from the front; remember the last node where they are equal. Stop at the first mismatch.
4. That remembered node is the LCA.

### Complexity

- **Time:** O(n) — two DFS traversals, each touching at most every node once.
- **Space:** O(n) — the two path lists plus recursion stack for a skewed tree.

### Code

```go
func pathBased(root, p, q *TreeNode) *TreeNode {
	var pathTo func(node, target *TreeNode, path *[]*TreeNode) bool
	pathTo = func(node, target *TreeNode, path *[]*TreeNode) bool {
		if node == nil {
			return false // dead end, nothing added
		}
		*path = append(*path, node) // tentatively include this node on the path
		if node == target {
			return true // found — keep the path as-is
		}
		if pathTo(node.Left, target, path) || pathTo(node.Right, target, path) {
			return true
		}
		*path = (*path)[:len(*path)-1] // backtrack: node is not on the path
		return false
	}

	var pPath, qPath []*TreeNode
	pathTo(root, p, &pPath) // build root→p
	pathTo(root, q, &qPath) // build root→q

	var lca *TreeNode
	for i := 0; i < len(pPath) && i < len(qPath); i++ {
		if pPath[i] == qPath[i] {
			lca = pPath[i] // still on the shared prefix
		} else {
			break // paths diverged; previous match is the LCA
		}
	}
	return lca
}
```

### Dry Run

Example 1: `p = 5, q = 1`. Tree root is `3`.

| Step | Action | State |
|------|--------|-------|
| 1 | Build path to `5` | `[3, 5]` |
| 2 | Build path to `1` | `[3, 1]` |
| 3 | Compare index 0 | `3 == 3` → `lca = 3` |
| 4 | Compare index 1 | `5 != 1` → break |

Result: `lca = 3` ✔

---

## Approach 2 — Single-Pass Recursion (Optimal)

### Intuition

Ask each node one question: *"how many of `{p, q}` are in my subtree?"* A node is the LCA precisely when the two targets are found on **different sides** (one in the left subtree, one in the right), or when the node **is** one of the targets and the other lies below it. Because post-order returns bottom-up, the first node that "sees" both is automatically the deepest such node.

### Algorithm

1. If `root` is `nil`, `p`, or `q`, return `root` (a "found here" signal, or nil for empty).
2. Recurse into the left subtree → `left`; into the right subtree → `right`.
3. If **both** `left` and `right` are non-nil, the targets split at `root` → return `root`.
4. Otherwise return whichever of `left`/`right` is non-nil (propagate the found target upward), or nil.

### Complexity

- **Time:** O(n) — each node is visited exactly once.
- **Space:** O(h) — recursion stack, `h` = tree height (O(n) skewed, O(log n) balanced).

### Code

```go
func recursive(root, p, q *TreeNode) *TreeNode {
	if root == nil || root == p || root == q {
		return root // base case: hit an empty branch or one of the targets
	}
	left := recursive(root.Left, p, q)   // search left subtree
	right := recursive(root.Right, p, q) // search right subtree
	if left != nil && right != nil {
		return root // p and q split here → this node is the LCA
	}
	if left != nil {
		return left // both targets (or the only found one) are on the left
	}
	return right // otherwise everything relevant is on the right (or nil)
}
```

### Dry Run

Example 1: `p = 5, q = 1`, root `3`.

| Node visited (post-order) | left | right | Returns |
|---------------------------|------|-------|---------|
| `6` | nil | nil | nil |
| `7` | nil | nil | nil |
| `4` | nil | nil | nil |
| `2` | nil (7) | nil (4) | nil |
| `5` (== p) | — | — | `5` (base case) |
| `0` | nil | nil | nil |
| `8` | nil | nil | nil |
| `1` (== q) | — | — | `1` (base case) |
| `3` | `5` | `1` | both non-nil → **`3`** |

Result: `3` ✔

---

## Approach 3 — Parent Pointers + Ancestor Set

### Intuition

If every node knows its parent, we can climb from any node up to the root. Collect all of `p`'s ancestors (including `p` itself) into a set; then climb from `q` — the first ancestor of `q` that also belongs to `p`'s set is the lowest common one.

### Algorithm

1. DFS/BFS from the root filling `parent[node]` for each node, stopping once both `p` and `q` have parents recorded.
2. Starting at `p`, climb to the root, inserting every node into a set.
3. Starting at `q`, climb upward; the first node found in the set is the LCA.

### Complexity

- **Time:** O(n) — one traversal to build parents, then O(h) to climb.
- **Space:** O(n) — the parent map and the ancestor set.

### Code

```go
func parentPointers(root, p, q *TreeNode) *TreeNode {
	parent := map[*TreeNode]*TreeNode{root: nil} // root has no parent
	stack := []*TreeNode{root}                   // explicit DFS stack
	for parent[p] == nil || parent[q] == nil {
		node := stack[len(stack)-1] // pop
		stack = stack[:len(stack)-1]
		if node.Left != nil {
			parent[node.Left] = node // record edge node→node.Left
			stack = append(stack, node.Left)
		}
		if node.Right != nil {
			parent[node.Right] = node // record edge node→node.Right
			stack = append(stack, node.Right)
		}
	}

	ancestors := map[*TreeNode]bool{} // p and all of p's ancestors
	for n := p; n != nil; n = parent[n] {
		ancestors[n] = true // climb p to root
	}
	for n := q; ; n = parent[n] {
		if ancestors[n] {
			return n // first common ancestor = lowest common ancestor
		}
	}
}
```

### Dry Run

Example 1: `p = 5, q = 1`.

| Step | Action | State |
|------|--------|-------|
| 1 | Build parents until 5 and 1 seen | `parent[5]=3, parent[1]=3, …` |
| 2 | Climb from `p=5` | `ancestors = {5, 3}` |
| 3 | Climb from `q=1`: is `1` in set? | no |
| 4 | Next `parent[1]=3`: is `3` in set? | yes → return `3` |

Result: `3` ✔

---

## Key Takeaways

- **The LCA is the split point.** In one post-order pass, a node returns non-nil from both children exactly when it is the first node whose subtree contains both targets — that is the LCA. This is the pattern to reach for first.
- **"A node can be its own descendant"** is baked into the base case `root == p || root == q`: hitting a target immediately returns it, so an ancestor-of relationship (Example 2) is handled for free.
- **Parent pointers turn a tree into a climbable structure**, converting LCA into the "intersection of two upward chains" problem — the same idea as finding where two linked lists merge (#160).
- The single-pass recursion needs no extra data structures beyond the O(h) stack, which is why it is the preferred answer when the tree fits in memory.

---

## Related Problems

- LeetCode #235 — Lowest Common Ancestor of a Binary Search Tree (BST version: use the ordering to descend in one path)
- LeetCode #1650 — Lowest Common Ancestor of a Binary Tree III (nodes carry parent pointers — the two-chain intersection)
- LeetCode #160 — Intersection of Two Linked Lists (same "merge point of two chains" idea)
- LeetCode #1123 — Lowest Common Ancestor of Deepest Leaves
- LeetCode #865 — Smallest Subtree with all the Deepest Nodes
