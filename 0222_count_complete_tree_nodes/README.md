# 0222 — Count Complete Tree Nodes

> LeetCode #222 · Difficulty: Easy
> **Categories:** Binary Tree, Binary Search, Depth-First Search, Bit Manipulation

---

## Problem Statement

Given the `root` of a **complete** binary tree, return the number of the nodes in the tree.

According to **Wikipedia**, every level, except possibly the last, is completely filled in a complete binary tree, and all nodes in the last level are as far left as possible. It can have between `1` and `2^h` nodes inclusive at the last level `h`.

Design an algorithm that runs in less than `O(n)` time complexity.

**Example 1:**

```
Input: root = [1,2,3,4,5,6]
Output: 6
```

**Example 2:**

```
Input: root = []
Output: 0
```

**Example 3:**

```
Input: root = [1]
Output: 1
```

**Constraints:**

- The number of nodes in the tree is in the range `[0, 5 * 10^4]`.
- `0 <= Node.val <= 5 * 10^4`
- The tree is guaranteed to be **complete**.

---

## Company Frequency

| Company    | Frequency        | Last Reported |
|------------|------------------|---------------|
| Amazon     | ★★★☆☆ Medium     | 2024          |
| Google     | ★★★☆☆ Medium     | 2024          |
| Microsoft  | ★★★☆☆ Medium     | 2023          |
| Apple      | ★★☆☆☆ Low        | 2023          |
| Bloomberg  | ★★☆☆☆ Low        | 2022          |
| Meta       | ★★☆☆☆ Low        | 2022          |

> ⚠️ Frequency data is crowd-sourced from LeetCode Discuss, Glassdoor, and
> community interview reports. Treat as a signal, not a guarantee.

---

## DSA Concepts Used

- **Tree Traversal (DFS)** — the O(n) baseline is a straight recursive node count → see [`/dsa/tree_traversal.md`](/dsa/tree_traversal.md)
- **Binary Search** — the sub-O(n) solutions binary-search either the perfect-subtree structure or the filled prefix of the last level → see [`/dsa/binary_search.md`](/dsa/binary_search.md)
- **Bit Manipulation** — heights give sizes as `2^h − 1` via `(1<<h)-1`, and last-level slot indices are decoded bit-by-bit to steer left/right → see [`/dsa/bit_manipulation.md`](/dsa/bit_manipulation.md)

---

## Approaches Overview

| # | Approach | Time | Space | When to use |
|---|----------|------|-------|-------------|
| 1 | Brute Force (Full Traversal) | O(n) | O(h) | Any tree; ignores the completeness hint (misses the follow-up target) |
| 2 | Perfect-Subtree Detection (Optimal) | O(log²n) | O(log n) | The canonical sub-O(n) answer; cleanest to reason about |
| 3 | Binary Search on Last Level | O(log²n) | O(1) | Same complexity, iterative, O(1) space; nice bit-indexing insight |

---

## Approach 1 — Brute Force (Full Traversal)

### Intuition
A tree's node count is `1 + count(left) + count(right)`. This is correct for any binary tree, so it works here too — it just doesn't use the "complete" guarantee, and therefore visits all `n` nodes.

### Algorithm
1. If `root == nil`, return 0.
2. Otherwise return `1 + bruteForce(root.Left) + bruteForce(root.Right)`.

### Complexity
- **Time:** O(n) — each node visited once.
- **Space:** O(h) — recursion stack; `h ≈ log n` for a complete tree.

### Code
```go
func bruteForce(root *TreeNode) int {
	if root == nil {
		return 0 // empty subtree contributes nothing
	}
	return 1 + bruteForce(root.Left) + bruteForce(root.Right)
}
```

### Dry Run
Example 1: `[1,2,3,4,5,6]` (node 2 has children 4,5; node 3 has child 6).

| Call | Returns |
|------|---------|
| count(4) | 1 + count(nil) + count(nil) = 1 |
| count(5) | 1 |
| count(2) | 1 + count(4) + count(5) = 3 |
| count(6) | 1 |
| count(3) | 1 + count(6) + count(nil) = 2 |
| count(1) | 1 + count(2) + count(3) = **6** |

Answer `6`. ✔

---

## Approach 2 — Perfect-Subtree Detection (Optimal)

### Intuition
A **perfect** subtree (every level full) of height `h` has exactly `2^h − 1` nodes — no counting needed. Detect it by comparing the leftmost-path depth and rightmost-path depth: in a complete tree, equal depths ⇒ perfect. When they differ, only one child subtree is imperfect, so recurse into both; the perfect side resolves instantly, giving the O(log²n) bound.

### Algorithm
1. If `root == nil`, return 0.
2. `lh =` steps down the all-left path; `rh =` steps down the all-right path.
3. If `lh == rh`, the subtree is perfect → return `(1 << lh) - 1`.
4. Otherwise return `1 + perfectSubtree(root.Left) + perfectSubtree(root.Right)`.

### Complexity
- **Time:** O(log²n) — O(log n) levels of recursion, each computing two O(log n) height walks.
- **Space:** O(log n) — recursion depth equals tree height.

### Code
```go
func perfectSubtree(root *TreeNode) int {
	if root == nil {
		return 0 // empty subtree
	}
	lh := leftHeight(root)  // depth following only left children
	rh := rightHeight(root) // depth following only right children
	if lh == rh {
		return (1 << lh) - 1 // perfect subtree of height lh → 2^lh − 1 nodes
	}
	return 1 + perfectSubtree(root.Left) + perfectSubtree(root.Right)
}

func leftHeight(node *TreeNode) int {
	h := 0
	for node != nil {
		h++
		node = node.Left
	}
	return h
}

func rightHeight(node *TreeNode) int {
	h := 0
	for node != nil {
		h++
		node = node.Right
	}
	return h
}
```

### Dry Run
Example 1: `[1,2,3,4,5,6]`.

| Node | leftHeight | rightHeight | Perfect? | Result |
|------|-----------|-------------|----------|--------|
| 1 | 1→2→4 = 3 | 1→3→6 = 3 | **yes** (lh=rh=3) | `(1<<3)-1 = 7`? |

Wait — the root's left path `1→2→4` and right path `1→3→6` are both length 3, but the tree only has 6 nodes, not 7. The completeness guarantee ensures `lh==rh` ⇒ perfect **only when the tree truly is perfect**; here node 3 is missing its right child, so the two heights are *not* both 3 for the whole tree. Recomputing carefully: rightHeight walks `1→3→6→nil`, that is 3 nodes deep, and leftHeight walks `1→2→4→nil`, also 3. They match, which would wrongly claim 7.

The subtlety: this shortcut is applied per-subtree and the equal-height test is safe for a *complete* tree because a complete tree with matching left/right spine heights **is** perfect. Node 3's subtree is `[3,6]` — leftHeight `3→6` = 2, rightHeight `3→nil` = 1, they differ, so it recurses. Let's trace properly:

| Call | lh | rh | Branch | Value |
|------|----|----|--------|-------|
| node 1 | 3 (1,2,4) | 3 (1,3,6) | equal → but is it perfect? |  |

For node 1 the spines are `1,2,4` and `1,3,6`; both depth 3. A complete tree whose two spines are equal length is perfect, so this returns `2^3 - 1 = 7`. But the tree has 6 nodes... The resolution: **node 6 must be node 3's left child** in a complete tree (last level filled left-to-right), so the rightmost path is `1→3→6` only if 6 is a right child. In LeetCode's complete tree `[1,2,3,4,5,6]`, node 3's single child (6) is its **left** child. Therefore rightHeight walks `1→3→nil` after node 3 (node 3 has no right child) = depth 2, while leftHeight walks `1→2→4` = depth 3. `lh=3 ≠ rh=2`, so we recurse:

| Call | lh | rh | Branch | Value |
|------|----|----|--------|-------|
| node 1 | 3 | 2 | differ → recurse | 1 + L + R |
| node 2 | 2 (2,4) | 2 (2,5) | equal → perfect | `2^2−1 = 3` |
| node 3 | 2 (3,6) | 1 (3) | differ → recurse | 1 + L + R |
| node 6 | 1 | 1 | equal → perfect | `2^1−1 = 1` |
| node 3 (nil right) | — | — | 1 + 1 + 0 | 2 |
| node 1 total | | | 1 + 3 + 2 | **6** |

Answer `6`. ✔ (Key correctness fact: last level fills left-to-right, so node 3's child is a *left* child, breaking the spine equality at the root.)

---

## Approach 3 — Binary Search on the Last Level

### Intuition
The top `h` levels of a height-`h` complete tree are always full: `2^h − 1` nodes. The last level holds a left-filled prefix of its `2^h` slots. "Is slot `i` present?" is monotone (all present slots come before all absent ones), so binary-search the boundary. To test slot `i`, read its `h` bits from most- to least-significant: bit 1 ⇒ go right, bit 0 ⇒ go left; if you never hit `nil`, the slot exists.

### Algorithm
1. `h =` edges on the leftmost path (tree height). If `h == 0`, return 1.
2. Binary-search `lo=0, hi=2^h−1` over last-level slot indices; `exists(i)` walks `h` steps guided by the bits of `i`.
3. Return `(2^h − 1)` upper nodes `+ lo` present last-level leaves.

### Complexity
- **Time:** O(log²n) — O(log n) search steps × O(log n) per existence walk.
- **Space:** O(1) — iterative, no recursion.

### Code
```go
func binarySearchLastLevel(root *TreeNode) int {
	if root == nil {
		return 0
	}
	h := 0
	for n := root.Left; n != nil; n = n.Left {
		h++ // height = edges down the leftmost path
	}
	if h == 0 {
		return 1 // only the root exists
	}
	lo, hi := 0, (1<<h)-1
	for lo <= hi {
		mid := (lo + hi) / 2
		if exists(root, mid, h) {
			lo = mid + 1 // slot present → search right
		} else {
			hi = mid - 1 // slot absent → search left
		}
	}
	return (1<<h - 1) + lo
}

func exists(root *TreeNode, idx, h int) bool {
	node := root
	for bit := h - 1; bit >= 0; bit-- {
		if idx&(1<<bit) != 0 {
			node = node.Right // bit set → right
		} else {
			node = node.Left // bit clear → left
		}
		if node == nil {
			return false
		}
	}
	return true
}
```

### Dry Run
Example 1: `[1,2,3,4,5,6]`. Leftmost path `1→2→4` → `h = 2` edges. Upper nodes `= 2^2 − 1 = 3` (nodes 1,2,3). Last level has slots `0..3` (`2^2 = 4`).

| mid | bits (h=2) | walk | exists? | lo,hi |
|-----|-----------|------|---------|-------|
| start | | | | lo=0, hi=3 |
| 1 | `01` → L,R | 1→2→5 | yes | lo=2 |
| 2 | `10` → R,L | 1→3→6 | yes | lo=3 |
| 3 | `11` → R,R | 1→3→nil | no | hi=2 |

Loop ends (`lo=3 > hi=2`). Present leaves `= lo = 3` (slots 0,1,2 = nodes 4,5,6). Total `= 3 (upper) + 3 = 6`. ✔

---

## Key Takeaways
- The "complete tree" guarantee is what unlocks sub-O(n): a matching left/right spine means a **perfect** subtree of size `2^h − 1`, computed in O(1).
- Last-level slots fill left-to-right, so "is slot i present?" is **monotone** — a textbook binary-search-on-answer setup, where the path to a slot is decoded from the bits of its index.
- Sizes and slot navigation come straight from bit tricks: `(1<<h)-1` for a perfect subtree's node count, `idx & (1<<bit)` to steer the descent.

---

## Related Problems
- LeetCode #104 — Maximum Depth of Binary Tree (height walk building block)
- LeetCode #110 — Balanced Binary Tree (per-node height comparisons)
- LeetCode #958 — Check Completeness of a Binary Tree (verifies the property this problem assumes)
- LeetCode #199 — Binary Tree Right Side View (rightmost-path reasoning)
