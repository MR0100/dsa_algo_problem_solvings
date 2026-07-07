# 0235 — Lowest Common Ancestor of a Binary Search Tree

> LeetCode #235 · Difficulty: Medium
> **Categories:** Binary Search Tree, Tree, Depth-First Search

---

## Problem Statement

Given a binary search tree (BST), find the lowest common ancestor (LCA) node of two given nodes in the BST.

According to the definition of LCA on Wikipedia: "The lowest common ancestor is defined between two nodes `p` and `q` as the lowest node in `T` that has both `p` and `q` as descendants (where we allow **a node to be a descendant of itself**)."

**Example 1:**
```
Input: root = [6,2,8,0,4,7,9,null,null,3,5], p = 2, q = 8
Output: 6
Explanation: The LCA of nodes 2 and 8 is 6.
```

**Example 2:**
```
Input: root = [6,2,8,0,4,7,9,null,null,3,5], p = 2, q = 4
Output: 2
Explanation: The LCA of nodes 2 and 4 is 2, since a node can be a descendant of itself according to the LCA definition.
```

**Example 3:**
```
Input: root = [2,1], p = 2, q = 1
Output: 2
```

**Constraints:**
- The number of nodes in the tree is in the range `[2, 10⁵]`.
- `-10⁹ <= Node.val <= 10⁹`
- All `Node.val` are **unique**.
- `p != q`
- `p` and `q` will exist in the BST.

---

## Company Frequency

| Company    | Frequency        | Last Reported |
|------------|------------------|---------------|
| Amazon     | ★★★★☆ High       | 2024          |
| Meta       | ★★★★☆ High       | 2024          |
| Microsoft  | ★★★☆☆ Medium     | 2024          |
| Google     | ★★★☆☆ Medium     | 2023          |
| Bloomberg  | ★★☆☆☆ Low        | 2023          |

> ⚠️ Frequency data is crowd-sourced from LeetCode Discuss, Glassdoor, and
> community interview reports. Treat as a signal, not a guarantee.

---

## DSA Concepts Used
- **Binary Search Tree** — the ordering property is what lets us pick a direction at each node → see [`/dsa/binary_search_tree.md`](/dsa/binary_search_tree.md)
- **Tree Traversal** — descending root-to-node, and building root-to-node paths → see [`/dsa/tree_traversal.md`](/dsa/tree_traversal.md)

---

## Approaches Overview
| # | Approach | Time | Space | When to use |
|---|----------|------|-------|-------------|
| 1 | Recursive BST Walk | O(h) | O(h) | Cleanest expression of the BST logic |
| 2 | Iterative BST Walk (Optimal) | O(h) | O(1) | Best space; the canonical answer |
| 3 | Root-to-Node Paths | O(h) | O(h) | Generalises to any binary tree (LC #236) |

(`h` = tree height; O(log n) if balanced, O(n) worst case.)

---

## Approach 1 — Recursive BST Walk

### Intuition
In a BST every node partitions keys: everything smaller is on the left, everything larger on the right. The LCA of `p` and `q` is the first node where they fall on **different** sides (or one equals the node). If both are smaller, their shared ancestor is left; if both larger, right; otherwise the current node is the split point — the LCA.

### Algorithm
1. At node `root`: if both `p.Val` and `q.Val < root.Val`, recurse left.
2. Else if both `> root.Val`, recurse right.
3. Else `root` is the LCA (the values straddle it, or one equals it).

### Complexity
- **Time:** O(h) — one node per level.
- **Space:** O(h) — recursion stack.

### Code
```go
func recursiveBST(root, p, q *TreeNode) *TreeNode {
	if p.Val < root.Val && q.Val < root.Val {
		return recursiveBST(root.Left, p, q) // both keys live in the left subtree
	}
	if p.Val > root.Val && q.Val > root.Val {
		return recursiveBST(root.Right, p, q) // both keys live in the right subtree
	}
	return root // p and q split here (or one equals root) → lowest common ancestor
}
```

### Dry Run
Tree `[6,2,8,0,4,7,9,null,null,3,5]`, `p = 2`, `q = 8`:

| Call | root.Val | both < root? | both > root? | action |
|------|----------|--------------|--------------|--------|
| 1    | 6        | 2<6 but 8>6 → no | 8>6 but 2<6 → no | values straddle → return **6** |

Return node `6`. ✓

---

## Approach 2 — Iterative BST Walk (Optimal)

### Intuition
The recursion is tail-shaped: each step just moves to one child. Replacing it with a `while` loop removes the call stack entirely, giving the optimal O(h) time and **O(1) space** solution.

### Algorithm
1. Start `cur = root`.
2. While `cur != nil`:
   - if both `p`,`q` `< cur.Val`: `cur = cur.Left`.
   - else if both `> cur.Val`: `cur = cur.Right`.
   - else return `cur` (split point).

### Complexity
- **Time:** O(h) — descends one level per iteration.
- **Space:** O(1) — a single pointer.

### Code
```go
func iterativeBST(root, p, q *TreeNode) *TreeNode {
	cur := root
	for cur != nil {
		switch {
		case p.Val < cur.Val && q.Val < cur.Val:
			cur = cur.Left // both smaller → go left
		case p.Val > cur.Val && q.Val > cur.Val:
			cur = cur.Right // both larger → go right
		default:
			return cur // values straddle cur (or one equals it) → LCA
		}
	}
	return nil // unreachable given p and q exist in the tree
}
```

### Dry Run
Same tree, `p = 2`, `q = 8`:

| Step | cur.Val | both < cur? | both > cur? | next cur |
|------|---------|-------------|-------------|----------|
| 1    | 6       | no (8 > 6)  | no (2 < 6)  | straddle → return **6** |

Return node `6`. ✓

For `p = 2`, `q = 4`: at `cur = 6`, both `2,4 < 6` → go left to `cur = 2`. Now `p.Val == cur.Val` (2), so neither "both smaller" nor "both larger" holds → return **2** (a node is its own descendant). ✓

---

## Approach 3 — General LCA via Root-to-Node Paths

### Intuition
Ignore the BST property for a moment: the LCA of two nodes is the **last common node** on their two root-to-node paths. Collect both paths as lists of nodes, then walk them in parallel; the deepest position where they still agree is the LCA. Slower and heavier than the BST-specific walk, but it generalises to arbitrary binary trees (see LC #236).

### Algorithm
1. Using the BST ordering, build `pathP = root..p` and `pathQ = root..q`.
2. Walk both from the root while `pathP[i] == pathQ[i]`.
3. The last matching node is the LCA.

### Complexity
- **Time:** O(h) — building each path is O(h); comparison is O(h).
- **Space:** O(h) — the two path slices.

### Code
```go
func pathIntersection(root, p, q *TreeNode) *TreeNode {
	pathTo := func(target *TreeNode) []*TreeNode {
		path := []*TreeNode{}
		cur := root
		for cur != nil {
			path = append(path, cur) // record every node on the way down
			if target.Val < cur.Val {
				cur = cur.Left
			} else if target.Val > cur.Val {
				cur = cur.Right
			} else {
				break // reached the target node itself
			}
		}
		return path
	}

	pathP := pathTo(p)
	pathQ := pathTo(q)

	var lca *TreeNode
	for i := 0; i < len(pathP) && i < len(pathQ); i++ {
		if pathP[i] == pathQ[i] {
			lca = pathP[i] // still on the common prefix → update the candidate
		} else {
			break // paths diverge here; earlier node stays as the LCA
		}
	}
	return lca
}
```

### Dry Run
Same tree, `p = 2`, `q = 8`:

| Path | nodes (by value) |
|------|------------------|
| pathP (to 2) | [6, 2] |
| pathQ (to 8) | [6, 8] |

Parallel walk:

| i | pathP[i] | pathQ[i] | equal? | lca |
|---|----------|----------|--------|-----|
| 0 | 6        | 6        | yes    | 6   |
| 1 | 2        | 8        | no → break | 6 |

Return node `6`. ✓

---

## Key Takeaways
- **Use the BST ordering.** When both targets are on the same side of a node, the LCA is deeper on that side; the moment they split, you've found it — no full tree search needed.
- Prefer the **iterative** walk: same O(h) time but O(1) space, and no recursion depth concern on a 10⁵-node tree.
- A node can be its **own** ancestor — the "straddle or equal" default branch handles the `p == cur` / `q == cur` case for free.
- The path-intersection method is the bridge to the **general** LCA problem (#236) where there's no ordering to exploit.

## Related Problems
- LeetCode #236 — Lowest Common Ancestor of a Binary Tree (no BST ordering)
- LeetCode #1650 — LCA of a Binary Tree III (nodes have parent pointers)
- LeetCode #700 — Search in a BST (same directional descent)
- LeetCode #98 — Validate Binary Search Tree (BST ordering property)
