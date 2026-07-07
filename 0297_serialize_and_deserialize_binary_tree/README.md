# 0297 — Serialize and Deserialize Binary Tree

> LeetCode #297 · Difficulty: Hard
> **Categories:** Tree, Depth-First Search, Breadth-First Search, Design, String, Binary Tree

---

## Problem Statement

Serialization is the process of converting a data structure or object into a sequence of bits so that it can be stored in a file or memory buffer, or transmitted across a network connection link to be reconstructed later in the same or another computer environment.

Design an algorithm to serialize and deserialize a binary tree. There is no restriction on how your serialization/deserialization algorithm should work. You just need to ensure that a binary tree can be serialized to a string and this string can be deserialized to the original tree structure.

**Clarification:** The input/output format is the same as [how LeetCode serializes a binary tree](https://support.leetcode.com/hc/en-us/articles/32442719377939-How-to-create-test-cases-on-LeetCode#h_01J5EGREAW3NAEJ14XC07GRW1A). You do not necessarily need to follow this format, so please be creative and come up with different approaches yourself.

**Example 1:**

```
Input: root = [1,2,3,null,null,4,5]
Output: [1,2,3,null,null,4,5]
```

**Example 2:**

```
Input: root = []
Output: []
```

**Constraints:**

- The number of nodes in the tree is in the range `[0, 10^4]`.
- `-1000 <= Node.val <= 1000`

---

## Company Frequency

| Company    | Frequency        | Last Reported |
|------------|------------------|---------------|
| Amazon     | ★★★★★ Very High  | 2024          |
| Google     | ★★★★★ Very High  | 2024          |
| Meta       | ★★★★☆ High       | 2024          |
| Microsoft  | ★★★★☆ High       | 2023          |
| Bloomberg  | ★★★☆☆ Medium     | 2023          |
| LinkedIn   | ★★★☆☆ Medium     | 2023          |

> ⚠️ Frequency data is crowd-sourced from LeetCode Discuss, Glassdoor, and community interview reports. Treat as a signal, not a guarantee.

---

## DSA Concepts Used

- **Depth-First Search (preorder)** — a preorder walk with null markers uniquely encodes a tree → see [`/dsa/tree_traversal.md`](/dsa/tree_traversal.md)
- **Breadth-First Search (level order)** — queue-based encode/decode mirroring LeetCode's display format → see [`/dsa/graph_bfs_dfs.md`](/dsa/graph_bfs_dfs.md)
- **Design / Codec pattern** — pairing a `serialize`/`deserialize` method set → see [`/dsa/design_data_structures.md`](/dsa/design_data_structures.md)
- **String parsing** — tokenizing on commas, handling markers → see [`/dsa/string_algorithms.md`](/dsa/string_algorithms.md)

---

## Approaches Overview

| # | Approach | Time (ser+deser) | Space | When to use |
|---|----------|------------------|-------|-------------|
| 1 | Preorder DFS Codec | O(n) | O(n) | Simplest to reason about; default |
| 2 | BFS Level-Order Codec | O(n) | O(n) | Matches LeetCode array format |

---

## Approach 1 — Preorder DFS Codec (Optimal)

### Intuition
A preorder traversal that **also records nil children** as a sentinel (`#`) is enough to reconstruct the tree. While reading back, you always know whether the current slot is a real node or an absent child, so a single traversal suffices in both directions.

### Algorithm
1. **Serialize:** visit root; append `#` if nil else its value; recurse left then right; join with commas.
2. **Deserialize:** split into tokens, keep an index; read one token — if `#` return nil, else create a node and recursively build left then right (preorder consumes tokens in the same order they were written).

### Complexity
- **Time:** O(n) to serialize + O(n) to deserialize — each node/marker touched once.
- **Space:** O(n) — the output string plus O(h) recursion depth.

### Code
```go
type PreorderCodec struct{}

func (PreorderCodec) serialize(root *TreeNode) string {
	var sb strings.Builder
	var dfs func(node *TreeNode)
	dfs = func(node *TreeNode) {
		if node == nil {
			sb.WriteString("#,")
			return
		}
		sb.WriteString(strconv.Itoa(node.Val))
		sb.WriteByte(',')
		dfs(node.Left)
		dfs(node.Right)
	}
	dfs(root)
	return sb.String()
}

func (PreorderCodec) deserialize(data string) *TreeNode {
	tokens := strings.Split(data, ",")
	idx := 0
	var build func() *TreeNode
	build = func() *TreeNode {
		tok := tokens[idx]
		idx++
		if tok == "#" {
			return nil
		}
		val, _ := strconv.Atoi(tok)
		node := &TreeNode{Val: val}
		node.Left = build()
		node.Right = build()
		return node
	}
	return build()
}
```

### Dry Run
Serialize Example 1 (`1,2,3` with `4,5` under `3`):

| visit | emit | running string |
|-------|------|----------------|
| 1     | `1`  | `1,` |
| 2     | `2`  | `1,2,` |
| 2.L   | `#`  | `1,2,#,` |
| 2.R   | `#`  | `1,2,#,#,` |
| 3     | `3`  | `1,2,#,#,3,` |
| 4     | `4`  | `1,2,#,#,3,4,` |
| 4.L,4.R | `#,#` | `1,2,#,#,3,4,#,#,` |
| 5     | `5`  | `...,5,` |
| 5.L,5.R | `#,#` | `1,2,#,#,3,4,#,#,5,#,#,` |

Deserialize reads `1` (root) → build left starts, reads `2` → its children `#,#` (both nil) → back up, build right reads `3` → left `4` (`#,#`), right `5` (`#,#`). Reconstructed tree matches the original. ✔

---

## Approach 2 — BFS Level-Order Codec

### Intuition
A breadth-first sweep records nodes level by level; writing `#` for each missing child preserves shape. On the way back, a queue re-links children in the exact order they were written — this is the same format LeetCode uses to display `[1,2,3,null,null,4,5]`.

### Algorithm
1. **Serialize:** push root; pop a node — if nil append `#`, else append value and push both children; repeat until the queue empties.
2. **Deserialize:** first token is the root, enqueue it; pop a parent — the next two tokens are its left/right children; create the non-`#` ones, attach, and enqueue them.

### Complexity
- **Time:** O(n) to serialize + O(n) to deserialize.
- **Space:** O(n) — queue plus output.

### Code
```go
type BFSCodec struct{}

func (BFSCodec) serialize(root *TreeNode) string {
	if root == nil {
		return ""
	}
	var out []string
	queue := []*TreeNode{root}
	for len(queue) > 0 {
		node := queue[0]
		queue = queue[1:]
		if node == nil {
			out = append(out, "#")
			continue
		}
		out = append(out, strconv.Itoa(node.Val))
		queue = append(queue, node.Left)
		queue = append(queue, node.Right)
	}
	return strings.Join(out, ",")
}

func (BFSCodec) deserialize(data string) *TreeNode {
	if data == "" {
		return nil
	}
	tokens := strings.Split(data, ",")
	root := &TreeNode{Val: mustAtoi(tokens[0])}
	queue := []*TreeNode{root}
	i := 1
	for len(queue) > 0 && i < len(tokens) {
		parent := queue[0]
		queue = queue[1:]

		if tokens[i] != "#" {
			parent.Left = &TreeNode{Val: mustAtoi(tokens[i])}
			queue = append(queue, parent.Left)
		}
		i++
		if i < len(tokens) && tokens[i] != "#" {
			parent.Right = &TreeNode{Val: mustAtoi(tokens[i])}
			queue = append(queue, parent.Right)
		}
		i++
	}
	return root
}
```

### Dry Run
Serialize Example 1:

| pop | emit | enqueue |
|-----|------|---------|
| 1   | `1`  | 2, 3 |
| 2   | `2`  | nil, nil |
| 3   | `3`  | 4, 5 |
| nil | `#`  | — |
| nil | `#`  | — |
| 4   | `4`  | nil, nil |
| 5   | `5`  | nil, nil |
| nil×4 | `#,#,#,#` | — |

Result: `1,2,3,#,#,4,5,#,#,#,#`.

Deserialize: root `1`; parent `1` → children `2`,`3` (enqueue); parent `2` → children `#`,`#` (none); parent `3` → children `4`,`5` (enqueue); parents `4`,`5` → all `#`. Tree matches. ✔

---

## Key Takeaways
- **Null markers are the trick**: a traversal that records absent children is self-delimiting and needs only ONE traversal order to reconstruct (unlike preorder+inorder for value-unique trees).
- **Preorder DFS** reconstructs recursively with a shared index — left builds fully before right because both writer and reader agree on order.
- **BFS level order** naturally reproduces LeetCode's `[...]` display and re-links via a queue of parents.
- Duplicate values are fine here because shape (via markers), not value uniqueness, drives reconstruction.

---

## Related Problems
- LeetCode #449 — Serialize and Deserialize BST (BST lets you drop null markers)
- LeetCode #428 — Serialize and Deserialize N-ary Tree
- LeetCode #536 — Construct Binary Tree from String
- LeetCode #105 — Construct Binary Tree from Preorder and Inorder Traversal
