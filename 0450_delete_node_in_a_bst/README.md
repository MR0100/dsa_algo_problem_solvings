# 0450 — Delete Node in a BST

> LeetCode #450 · Difficulty: Medium
> **Categories:** Tree, Binary Search Tree, Recursion

---

## Problem Statement

Given a `root` node reference of a BST and a `key`, delete the node with the given `key` in the BST. Return *the root node reference (possibly updated) of the BST*.

Basically, the deletion can be divided into two stages:

1. Search for a node to remove.
2. If the node is found, delete the node.

**Example 1:**

```
Input: root = [5,3,6,2,4,null,7], key = 3
Output: [5,4,6,2,null,null,7]
Explanation: Given key to delete is 3. So we find the node with value 3 and delete it.
One valid answer is [5,4,6,2,null,null,7], shown in the above BST.
Please notice that another valid answer is [5,2,6,null,4,null,7] and it's also accepted.
```

**Example 2:**

```
Input: root = [5,3,6,2,4,null,7], key = 0
Output: [5,3,6,2,4,null,7]
Explanation: The tree does not contain a node with value = 0.
```

**Example 3:**

```
Input: root = [], key = 0
Output: []
```

**Constraints:**

- The number of nodes in the tree is in the range `[0, 10^4]`.
- `-10^5 <= Node.val <= 10^5`
- Each node has a **unique** value.
- `root` is a valid binary search tree.
- `-10^5 <= key <= 10^5`

**Follow up:** Could you solve it with time complexity `O(height of tree)`?

---

## Company Frequency

| Company    | Frequency        | Last Reported |
|------------|------------------|---------------|
| Amazon     | ★★★★☆ High       | 2024          |
| Microsoft  | ★★★★☆ High       | 2024          |
| Google     | ★★★☆☆ Medium     | 2023          |
| Meta       | ★★★☆☆ Medium     | 2023          |
| Bloomberg  | ★★☆☆☆ Low        | 2022          |
| Apple      | ★★☆☆☆ Low        | 2022          |

> ⚠️ Frequency data is crowd-sourced from LeetCode Discuss, Glassdoor, and
> community interview reports. Treat as a signal, not a guarantee.

---

## DSA Concepts Used

- **Binary Search Tree** — both search and repair rely on the BST ordering; the two-children case is resolved with the **inorder successor** (smallest node in the right subtree), the next-larger value that keeps the ordering intact → see [`/dsa/binary_search_tree.md`](/dsa/binary_search_tree.md)
- **Tree Traversal** — reaching the target and locating the successor are downward walks; the recursive form reattaches subtrees via return values on the way back up → see [`/dsa/tree_traversal.md`](/dsa/tree_traversal.md)

---

## Approaches Overview

| # | Approach | Time | Space | When to use |
|---|----------|------|-------|-------------|
| 1 | Recursive Delete (Inorder Successor) | O(h) | O(h) | Cleanest to write; subtrees reattach through return values |
| 2 | Iterative Delete (Parent Pointer) | O(h) | O(1) | O(1) space follow-up; explicit pointer rewiring |

*(h = height of the tree: O(log n) balanced, O(n) skewed.)*

---

## Approach 1 — Recursive Delete (Inorder Successor)

### Intuition

Removing a BST node splits into three cases once you find it:

- **Leaf or single child** — return the (possibly `nil`) child so the parent reattaches it directly.
- **Two children** — you cannot just drop it, so replace its *value* with its **inorder successor**: the smallest value in the right subtree (that node has no left child). The successor is the next-larger value, so overwriting keeps the BST ordering valid. Then recursively delete that successor from the right subtree — and since the successor has no left child, that recursive delete lands in an easy single-child case.

Recursion reattaches subtrees automatically: each call returns the new root of the subtree it was given, and the parent stores that return value into the correct child pointer.

### Algorithm

1. If `root == nil`, return `nil` (key not present).
2. If `key < root.Val`, set `root.Left = delete(root.Left, key)`; if `key > root.Val`, recurse right.
3. Else `root` is the target:
   - if `root.Left == nil`, return `root.Right`;
   - if `root.Right == nil`, return `root.Left`;
   - otherwise find the leftmost node of `root.Right` (the successor), copy its value into `root`, and `root.Right = delete(root.Right, successorValue)`.
4. Return `root`.

### Complexity

- **Time:** O(h) — one root-to-target path plus a successor descent inside the right subtree.
- **Space:** O(h) — recursion stack depth.

### Code

```go
func recursiveDelete(root *TreeNode, key int) *TreeNode {
	if root == nil {
		return nil // key not found; nothing to remove
	}
	if key < root.Val {
		root.Left = recursiveDelete(root.Left, key) // target is in the left subtree
	} else if key > root.Val {
		root.Right = recursiveDelete(root.Right, key) // target is in the right subtree
	} else {
		// Found the node to delete — handle by number of children.
		if root.Left == nil {
			return root.Right // 0 or 1 child: promote the right child
		}
		if root.Right == nil {
			return root.Left // exactly one (left) child: promote it
		}
		// Two children: find inorder successor = leftmost node of right subtree.
		succ := root.Right
		for succ.Left != nil {
			succ = succ.Left // walk to the smallest value greater than root
		}
		root.Val = succ.Val // overwrite value with successor's (ordering kept)
		// Remove the successor (a node with no left child) from the right subtree.
		root.Right = recursiveDelete(root.Right, succ.Val)
	}
	return root
}
```

### Dry Run

Example 1: `root = [5,3,6,2,4,null,7]`, `key = 3`.

| call | node.Val | branch taken | effect |
|------|----------|--------------|--------|
| delete(5, 3) | 5 | 3 < 5 → recurse left | `root.Left = delete(3, 3)` |
| delete(3, 3) | 3 | match; has both children (2 and 4) | successor = leftmost of right subtree = **4** |
| — | — | copy 4 into this node | node value 3 → 4 |
| delete(4, 4) on right subtree | 4 | match; `Left == nil` | return `Right` (nil) → node 4 removed from right subtree |

Node "3" is now value 4 with left child 2 and right child nil; subtree root 5 keeps right child 6 (with right child 7). Level-order: `[5,4,6,2,null,null,7]`. ✔

---

## Approach 2 — Iterative Delete (Parent Pointer)

### Intuition

The same three-case deletion, done with a loop to spend only O(1) extra space. Descend to the target while remembering its **parent**, so you know which of the parent's child pointers to rewire. When the target is found:

- **At most one child** → the replacement is that child (or `nil`).
- **Two children** → detach the inorder successor from the right subtree, tracking *its* parent so you can unlink it cleanly, then splice the successor into the target's slot, adopting the target's left (and, if the successor was deeper, its right) subtree.

Finally, hang the computed replacement off the parent — or, if the target was the root, return the replacement as the new root.

### Algorithm

1. Search from `root`, keeping `parent`, until `cur.Val == key` or `cur == nil`. If not found, return `root`.
2. Compute `replacement`:
   - if `cur.Left == nil` → `cur.Right`; else if `cur.Right == nil` → `cur.Left`;
   - else find successor `succ` (leftmost of `cur.Right`) with its parent `succParent`. If `succParent != cur`, set `succParent.Left = succ.Right` and `succ.Right = cur.Right`. Always `succ.Left = cur.Left`; `replacement = succ`.
3. Reattach: if `parent == nil`, return `replacement`; else set `parent.Left` or `parent.Right` (whichever equalled `cur`) to `replacement`. Return `root`.

### Complexity

- **Time:** O(h) — a search descent plus a successor descent.
- **Space:** O(1) — a fixed set of pointers, no recursion.

### Code

```go
func iterativeDelete(root *TreeNode, key int) *TreeNode {
	// Step 1: locate the target node and remember its parent.
	var parent *TreeNode
	cur := root
	for cur != nil && cur.Val != key {
		parent = cur
		if key < cur.Val {
			cur = cur.Left // go left for smaller keys
		} else {
			cur = cur.Right // go right for larger keys
		}
	}
	if cur == nil {
		return root // key absent → tree unchanged
	}

	// Step 2: compute the subtree that will replace cur.
	var replacement *TreeNode
	if cur.Left == nil {
		replacement = cur.Right // 0/1 child: right child (or nil) takes over
	} else if cur.Right == nil {
		replacement = cur.Left // only a left child
	} else {
		// Two children: detach inorder successor (leftmost of right subtree).
		succParent := cur
		succ := cur.Right
		for succ.Left != nil {
			succParent = succ
			succ = succ.Left
		}
		if succParent != cur {
			// Successor is deeper: unlink it, then let it adopt cur.Right.
			succParent.Left = succ.Right // successor's right child fills its slot
			succ.Right = cur.Right       // successor adopts the whole right subtree
		}
		succ.Left = cur.Left // successor adopts cur's left subtree
		replacement = succ   // successor now stands in for cur
	}

	// Step 3: attach the replacement where cur used to hang.
	if parent == nil {
		return replacement // cur was the root
	}
	if parent.Left == cur {
		parent.Left = replacement // cur was a left child
	} else {
		parent.Right = replacement // cur was a right child
	}
	return root
}
```

### Dry Run

Example 1: `root = [5,3,6,2,4,null,7]`, `key = 3`.

| step | state | detail |
|------|-------|--------|
| search | cur=5, key 3 < 5 → left; parent=5 | cur=3, `cur.Val == key` → stop |
| classify | cur=3 has both children (2, 4) | go to two-children branch |
| successor | succ = cur.Right = 4; `4.Left == nil` | loop body never runs → succParent = cur |
| adopt | `succParent == cur`, so skip unlink; `succ.Left = cur.Left` (= 2) | node 4 now has left child 2, right child nil |
| replacement | replacement = succ (node 4) | — |
| reattach | parent = 5, `parent.Left == cur` (node 3) | `parent.Left = node 4` |

Resulting tree level-order: `[5,4,6,2,null,null,7]`. ✔

---

## Key Takeaways

- **BST deletion has exactly three cases**: no child, one child, two children. Only the two-children case is subtle.
- **Two children ⇒ inorder successor (or predecessor).** Copy the successor's value up, then delete the successor — which is guaranteed to have no left child, collapsing to an easy case. Using the predecessor (max of left subtree) is equally valid and explains why LeetCode accepts multiple answers.
- **Recursion reattaches via return values**; iteration reattaches via an explicit **parent pointer**. Knowing both is the difference between O(h) stack space and O(1) space.
- **Watch the "successor is the direct right child" edge case.** When `cur.Right` itself has no left child, the successor *is* `cur.Right`; do not accidentally null out its right subtree — the guard `if succParent != cur` handles this.
- Search + repair is O(height), meeting the follow-up bound; a balanced BST makes this O(log n).

---

## Related Problems

- LeetCode #701 — Insert into a Binary Search Tree (the mirror operation)
- LeetCode #700 — Search in a Binary Search Tree (the search half)
- LeetCode #235 — Lowest Common Ancestor of a BST (BST-ordered descent)
- LeetCode #449 — Serialize and Deserialize BST (BST structure round-trip)
- LeetCode #776 — Split BST (splitting on a value using the ordering)
