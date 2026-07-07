# 0449 — Serialize and Deserialize BST

> LeetCode #449 · Difficulty: Medium
> **Categories:** Tree, Binary Search Tree, DFS, Design, String

---

## Problem Statement

Serialization is converting a data structure or object into a sequence of bits so that it can be stored in a file or memory buffer, or transmitted across a network connection link to be reconstructed later in the same or another computer environment.

Design an algorithm to serialize and deserialize a **binary search tree**. There is no restriction on how your serialization/deserialization algorithm should work. You need to ensure that a binary search tree can be serialized to a string, and this string can be deserialized to the original tree structure.

**The encoded string should be as compact as possible.**

**Example 1:**

```
Input: root = [2,1,3]
Output: [2,1,3]
```

**Example 2:**

```
Input: root = []
Output: []
```

**Constraints:**

- The number of nodes in the tree is in the range `[0, 10^4]`.
- `0 <= Node.val <= 10^4`
- The input tree is **guaranteed** to be a binary search tree.

---

## Company Frequency

| Company    | Frequency        | Last Reported |
|------------|------------------|---------------|
| Amazon     | ★★★★☆ High       | 2024          |
| Google     | ★★★★☆ High       | 2024          |
| Microsoft  | ★★★☆☆ Medium     | 2023          |
| Meta       | ★★★☆☆ Medium     | 2023          |
| Bloomberg  | ★★☆☆☆ Low        | 2022          |

> ⚠️ Frequency data is crowd-sourced from LeetCode Discuss, Glassdoor, and
> community interview reports. Treat as a signal, not a guarantee.

---

## DSA Concepts Used

- **Binary Search Tree** — the compact codec exploits the BST invariant (left < node < right) so a bare preorder value list is enough to reconstruct the tree, no null markers → see [`/dsa/binary_search_tree.md`](/dsa/binary_search_tree.md)
- **Tree Traversal (preorder DFS)** — both serialization and reconstruction are preorder walks; preorder emits the root before its subtrees, which is exactly what bound-based rebuild needs → see [`/dsa/tree_traversal.md`](/dsa/tree_traversal.md)
- **Design of Data Structures** — the deliverable is a `Codec` type with `serialize`/`deserialize` methods forming a lossless round-trip → see [`/dsa/design_data_structures.md`](/dsa/design_data_structures.md)

---

## Approaches Overview

| # | Approach | Time | Space | When to use |
|---|----------|------|-------|-------------|
| 1 | Preorder with Null Markers | O(n) | O(n) | General binary trees; simplest, but ~2× larger payload |
| 2 | Preorder-Only, BST-Bounded (Optimal) | O(n) | O(n) | BSTs; drops markers → "as compact as possible" |

---

## Approach 1 — Preorder with Null Markers

### Intuition

This is the general-purpose codec that ignores the BST property. A preorder value list alone is ambiguous for an arbitrary binary tree, but writing an explicit sentinel (`#`) for every `nil` child removes the ambiguity: the shape is fully captured, and consuming the same preorder stream rebuilds exactly one tree. It works for *any* binary tree — which is why it is the baseline before we specialise to BSTs. Cost: every leaf spends two `#` tokens, so the string is roughly twice as long as it needs to be.

### Algorithm

1. **serialize:** preorder DFS. Emit the node value, or `#` when the node is `nil`. Space-join.
2. **deserialize:** split into tokens; read left-to-right. `#` → `nil`; otherwise create the node, then recursively build its left subtree, then its right subtree (preorder order).

### Complexity

- **Time:** O(n) — serialize and deserialize each touch every node and marker once.
- **Space:** O(n) — the encoded string, plus O(h) recursion depth.

### Code

```go
type PreorderNullCodec struct{}

// serialize writes a preorder traversal with "#" sentinels for nil children.
func (PreorderNullCodec) serialize(root *TreeNode) string {
	var sb strings.Builder
	var pre func(node *TreeNode)
	pre = func(node *TreeNode) {
		if node == nil {
			sb.WriteString("# ") // sentinel marks an absent child
			return
		}
		sb.WriteString(strconv.Itoa(node.Val)) // node value first (preorder)
		sb.WriteByte(' ')
		pre(node.Left)  // then entire left subtree
		pre(node.Right) // then entire right subtree
	}
	pre(root)
	return strings.TrimSpace(sb.String())
}

// deserialize rebuilds the tree by consuming the preorder token stream.
func (PreorderNullCodec) deserialize(data string) *TreeNode {
	if data == "" {
		return nil
	}
	tokens := strings.Fields(data) // split on whitespace into value/"#" tokens
	pos := 0                       // cursor into tokens, advanced as we consume
	var build func() *TreeNode
	build = func() *TreeNode {
		tok := tokens[pos] // current token dictates node vs nil
		pos++
		if tok == "#" {
			return nil // sentinel → no node here
		}
		val, _ := strconv.Atoi(tok)
		node := &TreeNode{Val: val}
		node.Left = build()  // preorder: left subtree consumed next
		node.Right = build() // then right subtree
		return node
	}
	return build()
}
```

### Dry Run

Example 1: `root = [2,1,3]` (root 2, left 1, right 3). `serialize` produces `"2 1 # # 3 # #"`. Now `deserialize` consumes tokens:

| pos | token | action | tree so far |
|-----|-------|--------|-------------|
| 0 | `2` | make node 2; build its left | 2 |
| 1 | `1` | make node 1 (left of 2); build its left | 2←1 |
| 2 | `#` | node 1's left = nil | 2←1 |
| 3 | `#` | node 1's right = nil → 1 done | 2←1 |
| 4 | `3` | make node 3 (right of 2); build its left | 2←1, 2→3 |
| 5 | `#` | node 3's left = nil | — |
| 6 | `#` | node 3's right = nil → 3 done | — |

Reconstructed tree level-order: `[2,1,3]`. ✔

---

## Approach 2 — Preorder-Only, BST-Bounded (Optimal)

### Intuition

For a **BST**, the preorder sequence *by itself* determines the tree — no sentinels required — because the search-tree ordering tells us exactly where each value belongs. The first value is the root. Among the remaining preorder values, a contiguous prefix is smaller than the root (the left subtree) and the rest is larger (the right subtree). Reconstruct recursively with an allowed value window `(lower, upper)`: consume the next value only while it fits the current node's range. Because you never look at a value twice and the bounds are O(1) checks, rebuild is O(n) — and the payload has no markers, satisfying "as compact as possible."

### Algorithm

1. **serialize:** plain preorder DFS, values only, space-joined.
2. **deserialize:** keep a cursor over the parsed values and call `build(lower, upper)`:
   - If the stream is exhausted, or the next value falls outside `(lower, upper)`, return `nil` (it belongs to a different subtree).
   - Otherwise consume it as this subtree's root; recurse left with `upper = val`, recurse right with `lower = val`.
   - Bound the root call with `±∞`.

### Complexity

- **Time:** O(n) — every value is consumed exactly once; each range test is O(1).
- **Space:** O(n) — the value string, plus O(h) recursion depth.

### Code

```go
type BSTPreorderCodec struct{}

// serialize writes just the preorder values — BST ordering makes markers unnecessary.
func (BSTPreorderCodec) serialize(root *TreeNode) string {
	var vals []string
	var pre func(node *TreeNode)
	pre = func(node *TreeNode) {
		if node == nil {
			return // no marker: absence is inferred from value bounds later
		}
		vals = append(vals, strconv.Itoa(node.Val)) // preorder value
		pre(node.Left)
		pre(node.Right)
	}
	pre(root)
	return strings.Join(vals, " ")
}

// deserialize rebuilds the BST from the preorder values using value-range bounds.
func (BSTPreorderCodec) deserialize(data string) *TreeNode {
	if data == "" {
		return nil
	}
	fields := strings.Fields(data)
	vals := make([]int, len(fields))
	for i, f := range fields {
		vals[i], _ = strconv.Atoi(f) // parse preorder values once
	}
	pos := 0 // cursor into vals; only moves when a value is placed

	// build constructs the subtree whose node values must fall in (lower, upper).
	var build func(lower, upper int) *TreeNode
	build = func(lower, upper int) *TreeNode {
		if pos == len(vals) {
			return nil // stream exhausted
		}
		v := vals[pos] // peek the next preorder value
		if v < lower || v > upper {
			return nil // out of this node's allowed window → belongs elsewhere
		}
		pos++ // consume v as the current subtree's root
		node := &TreeNode{Val: v}
		node.Left = build(lower, v)  // left subtree values must be < v
		node.Right = build(v, upper) // right subtree values must be > v
		return node
	}
	return build(-1<<62, 1<<62)
}
```

### Dry Run

Example 1: `root = [2,1,3]`. `serialize` → `"2 1 3"`; parsed `vals = [2,1,3]`, `pos = 0`. Call `build(-∞, +∞)`:

| call | (lower, upper) | pos | v = vals[pos] | in range? | action |
|------|----------------|-----|---------------|-----------|--------|
| build(-∞,+∞) | (-∞, +∞) | 0 | 2 | yes | node 2, pos→1; recurse left (…,2) then right (2,…) |
| build(-∞,2) | (-∞, 2) | 1 | 1 | yes | node 1 (left of 2), pos→2; recurse left (-∞,1), right (1,2) |
| build(-∞,1) | (-∞, 1) | 2 | 3 | 3 > 1 → no | return nil (1's left) |
| build(1,2) | (1, 2) | 2 | 3 | 3 > 2 → no | return nil (1's right) → node 1 done |
| build(2,+∞) | (2, +∞) | 2 | 3 | yes | node 3 (right of 2), pos→3; recurse left (2,3), right (3,+∞) |
| build(2,3) | (2, 3) | 3 | — | pos==len | return nil (3's left) |
| build(3,+∞) | (3, +∞) | 3 | — | pos==len | return nil (3's right) → done |

Reconstructed level-order: `[2,1,3]`. ✔ — and it consumed 3 values with zero markers.

---

## Key Takeaways

- **A BST is uniquely determined by its preorder traversal alone.** The ordering invariant lets you infer subtree boundaries, so null markers are pure waste for a BST.
- **Bounds-based reconstruction is the workhorse.** Passing an `(lower, upper)` window down the recursion — consuming a value only if it fits — is the same technique behind LC #1008 (BST from preorder) and BST validity checks (LC #98).
- **Preorder gives you the root first**, which is exactly what top-down construction needs; that is why preorder (not inorder) is the natural serialization order here. Inorder alone would be ambiguous (it is just the sorted values).
- **Codec design = lossless round-trip.** Test `deserialize(serialize(t))` reproduces the tree structurally, including for the empty tree.
- If the tree were *not* guaranteed to be a BST, fall back to Approach 1 (markers) or an inorder+preorder pair — the compact trick relies on the ordering.

---

## Related Problems

- LeetCode #297 — Serialize and Deserialize Binary Tree (general tree; markers required)
- LeetCode #1008 — Construct Binary Search Tree from Preorder Traversal (the reconstruction half)
- LeetCode #98 — Validate Binary Search Tree (same `(lower, upper)` bounds idea)
- LeetCode #105 — Construct Binary Tree from Preorder and Inorder Traversal (non-BST reconstruction)
- LeetCode #428 — Serialize and Deserialize N-ary Tree
