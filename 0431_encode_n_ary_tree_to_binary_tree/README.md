# 0431 — Encode N-ary Tree to Binary Tree

> LeetCode #431 · Difficulty: Hard
> **Categories:** Tree, Depth-First Search, Breadth-First Search, Design, Binary Tree

---

## Problem Statement

Design an algorithm to encode an **N-ary tree** into a **binary tree** and decode the binary tree to get the original N-ary tree. An N-ary tree is a rooted tree in which each node has no more than N children. Similarly, a binary tree is a rooted tree in which each node has no more than 2 children. There is no restriction on how your encode/decode algorithm should work. You just need to ensure that an N-ary tree can be encoded to a binary tree and this binary tree can be decoded to the original N-ary tree structure.

*Nary-Tree* input serialization is represented in their level order traversal, each group of children is separated by the `null` value (see the example below).

For example, you may encode the following `3-ary` tree to a binary tree in this way:

```
Input: root = [1,null,3,2,4,null,5,6]

N-ary tree:
          1
       /  |  \
      3   2   4
     / \
    5   6
```

Note that the above is just an example which _might or might not_ work. You do not necessarily need to follow this format, so please be creative and come up with different approaches yourself.

**Constraints:**

- The number of nodes in the tree is in the range `[0, 10^4]`.
- `0 <= Node.val <= 10^4`
- The height of the n-ary tree is less than or equal to `1000`.
- Do not use class member/global/static variables to store states. Your encode and decode algorithms should be stateless.

---

## Company Frequency

| Company    | Frequency        | Last Reported |
|------------|------------------|---------------|
| Microsoft  | ★★★☆☆ Medium     | 2023          |
| Amazon     | ★★★☆☆ Medium     | 2023          |
| Google     | ★★☆☆☆ Low        | 2022          |
| Facebook   | ★★☆☆☆ Low        | 2022          |
| Uber       | ★★☆☆☆ Low        | 2021          |

> ⚠️ Frequency data is crowd-sourced from LeetCode Discuss, Glassdoor, and
> community interview reports. Treat as a signal, not a guarantee.

---

## DSA Concepts Used

- **Tree Traversal** — both encode and decode are single-pass tree walks (pre-order DFS for the sibling-chain mapping, BFS/level-order for the serialization variant); recognising that the transform is just "visit every node once" is the whole game → see [`/dsa/tree_traversal.md`](/dsa/tree_traversal.md)
- **Graph BFS/DFS** — the level-order serialization approach is a textbook BFS over the tree using an explicit queue, and the DFS recursion in the LCRS mapping is the same frontier idea depth-first → see [`/dsa/graph_bfs_dfs.md`](/dsa/graph_bfs_dfs.md)
- **Design / Serialization** — this is a "design a reversible codec" problem: pick a representation that stores enough information to reconstruct the original exactly (a bijection) → see [`/dsa/design_data_structures.md`](/dsa/design_data_structures.md)

---

## Approaches Overview

| # | Approach | Time | Space | When to use |
|---|----------|------|-------|-------------|
| 1 | Left-Child / Right-Sibling (Optimal) | O(V) encode + O(V) decode | O(H) recursion + O(V) tree | The canonical answer — a clean structural bijection, no strings, no parsing |
| 2 | BFS Level-Order String Serialization | O(V) encode + O(V) decode | O(V) | Shows the problem is "serialise then deserialise"; useful when you already have a serializer |

`V` = number of nodes, `H` = tree height.

---

## Approach 1 — Left-Child / Right-Sibling (Optimal)

### Intuition

A binary node has exactly two pointers; an N-ary node needs to remember *its first child* and *its next sibling*. Those are two things — so map them onto the two binary pointers:

- binary **`Left`** = the N-ary node's **first child**
- binary **`Right`** = the N-ary node's **next sibling**

All children of one N-ary node become a `Right`-linked chain that hangs off the parent's `Left` pointer. This "left-child, right-sibling" (LCRS) representation is a classic, and it is a **bijection**: no information is lost, so decoding is exact. To decode, `Left` gives you the first child; then you walk the `Right` chain to collect every remaining child of that node.

### Algorithm

1. **encode(nary):**
   1. If `nary` is `nil`, return `nil`.
   2. Create binary node `b` with `b.Val = nary.Val`.
   3. If `nary` has children: set `b.Left = encode(children[0])`.
   4. Walk `cur = b.Left`; for each remaining child `i ≥ 1`, set `cur.Right = encode(children[i])` and advance `cur = cur.Right`.
   5. Return `b`.
2. **decode(bin):**
   1. If `bin` is `nil`, return `nil`.
   2. Create N-ary node `n` with `n.Val = bin.Val`.
   3. Walk `child = bin.Left` along `Right` pointers; for each, append `decode(child)` to `n.Children`.
   4. Return `n`.

### Complexity

- **Time:** O(V) — encode touches each N-ary node once (creating one binary node); decode touches each binary node once. Both are linear in the node count.
- **Space:** O(H) for the recursion stack (H = height, ≤ 1000 per constraints) plus O(V) for the produced tree.

### Code

```go
type LCRSCodec struct{}

// encode builds the binary tree using the left-child/right-sibling mapping.
func (LCRSCodec) encode(root *Node) *TreeNode {
	if root == nil {
		return nil // empty N-ary tree ↔ empty binary tree
	}
	b := &TreeNode{Val: root.Val} // binary node carrying the same value
	if len(root.Children) > 0 {
		// The FIRST child hangs off the Left pointer.
		b.Left = (LCRSCodec{}).encode(root.Children[0])
		cur := b.Left // cur walks the right-linked sibling chain
		// Every subsequent child is chained via Right (sibling links).
		for i := 1; i < len(root.Children); i++ {
			cur.Right = (LCRSCodec{}).encode(root.Children[i])
			cur = cur.Right // advance to the newly attached sibling
		}
	}
	return b
}

// decode rebuilds the N-ary tree from the left-child/right-sibling binary tree.
func (LCRSCodec) decode(root *TreeNode) *Node {
	if root == nil {
		return nil // empty binary tree ↔ empty N-ary tree
	}
	n := &Node{Val: root.Val, Children: []*Node{}} // N-ary node, same value
	child := root.Left                             // Left points at the first child
	// Walk the sibling chain: each Right hop is the next child of `n`.
	for child != nil {
		n.Children = append(n.Children, (LCRSCodec{}).decode(child))
		child = child.Right // move to the next sibling in the chain
	}
	return n
}
```

### Dry Run

Encoding the example N-ary tree `1[3[5,6],2,4]` (root `1` has children `3,2,4`; node `3` has children `5,6`).

| Step | N-ary node | Action | Binary result |
|------|------------|--------|---------------|
| 1 | `1` | make binary `1`; first child `3` → `1.Left` | `1.Left = 3` |
| 2 | `1`'s siblings | child `2` → `3.Right`; child `4` → `2.Right` | `3.Right = 2`, `2.Right = 4` |
| 3 | `3` | make binary `3`; first child `5` → `3.Left` | `3.Left = 5` |
| 4 | `3`'s siblings | child `6` → `5.Right` | `5.Right = 6` |
| 5 | `2`, `4`, `5`, `6` | no children → no `Left` | leaves |

Encoded binary tree (Left = down, Right = across):

```
        1
       /
      3 → 2 → 4
     /
    5 → 6
```

**Decoding** reverses it: at binary `1`, `Left=3` starts the child chain; walk `3 → 2 → 4` via `Right` giving children `[3,2,4]`. Recurse at `3`: `Left=5`, walk `5 → 6` giving `[5,6]`. Result: `1[3[5,6],2,4]` — identical to the input. ✔

---

## Approach 2 — BFS Level-Order String Serialization

### Intuition

"Encode into a binary tree" really only demands that the binary tree carry **enough information** to rebuild the original. A serialised token stream already does that — so we can smuggle the stream through the binary tree by chaining nodes via `Left`, one token per node. We serialise the N-ary tree in **BFS/level order** as a self-describing sequence of `(value, childCount)` pairs. Because BFS emits a parent before its children, a queue during decode always creates a parent before it needs to attach its kids. This is deliberately *not* the elegant structural trick — it exists to show the problem reduces to "serialise then deserialise", and to make the LCRS mapping's cleanliness obvious by contrast.

### Algorithm

1. **encode:** BFS the N-ary tree with a queue; for each node emit two tokens: its value and its child count, then enqueue its children. Thread the token list through `Left` pointers of fresh binary nodes (each node holds one token in `Val`).
2. **decode:** Read the tokens back off the `Left` chain into a flat list. Pop the root's `(val, count)`, create it, and enqueue it with `remaining = count`. Repeatedly: if the front parent still needs children, pop the next `(val, count)` as one child, attach it, decrement the parent's `remaining`, and enqueue the child (it may have its own children).

### Complexity

- **Time:** O(V) — each node yields O(1) tokens on encode; decode consumes each token pair once.
- **Space:** O(V) — the token chain plus the BFS queues.

### Code

```go
type SerializeCodec struct{}

// encode serialises the N-ary tree (BFS) and threads the integer tokens through
// a Left-linked chain of binary nodes.
func (SerializeCodec) encode(root *Node) *TreeNode {
	if root == nil {
		return nil
	}
	tokens := []int{}      // flat stream: val, count, val, count, ...
	queue := []*Node{root} // BFS frontier over the N-ary tree
	for len(queue) > 0 {
		node := queue[0]
		queue = queue[1:]                                     // dequeue
		tokens = append(tokens, node.Val, len(node.Children)) // emit (val,count)
		queue = append(queue, node.Children...)               // children explored next
	}
	dummy := &TreeNode{}
	cur := dummy
	for _, t := range tokens {
		cur.Left = &TreeNode{Val: t} // store token in a binary node's value
		cur = cur.Left               // extend the chain downward via Left
	}
	return dummy.Left // first real token node is the encoded root
}

// decode reads the token chain back and BFS-rebuilds the N-ary tree.
func (SerializeCodec) decode(root *TreeNode) *Node {
	if root == nil {
		return nil
	}
	tokens := []int{}
	for cur := root; cur != nil; cur = cur.Left {
		tokens = append(tokens, cur.Val)
	}
	idx := 0
	rootVal, rootCnt := tokens[idx], tokens[idx+1]
	idx += 2
	nRoot := &Node{Val: rootVal, Children: []*Node{}}
	type item struct {
		node      *Node
		remaining int // children still to attach to this node
	}
	queue := []item{{nRoot, rootCnt}}
	for len(queue) > 0 && idx < len(tokens) {
		parent := &queue[0]
		if parent.remaining == 0 {
			queue = queue[1:] // this parent is satisfied; move on
			continue
		}
		val, cnt := tokens[idx], tokens[idx+1]
		idx += 2
		child := &Node{Val: val, Children: []*Node{}}
		parent.node.Children = append(parent.node.Children, child)
		parent.remaining--                       // one fewer child to place
		queue = append(queue, item{child, cnt}) // the child may have its own kids
	}
	return nRoot
}
```

### Dry Run

Encoding `1[3[5,6],2,4]`. BFS visit order: `1, 3, 2, 4, 5, 6`.

| Dequeued node | children count | tokens appended | queue after |
|---------------|----------------|-----------------|-------------|
| `1` | 3 | `1, 3` | `[3, 2, 4]` |
| `3` | 2 | `3, 2` | `[2, 4, 5, 6]` |
| `2` | 0 | `2, 0` | `[4, 5, 6]` |
| `4` | 0 | `4, 0` | `[5, 6]` |
| `5` | 0 | `5, 0` | `[6]` |
| `6` | 0 | `6, 0` | `[]` |

Token stream: `[1,3, 3,2, 2,0, 4,0, 5,0, 6,0]`, threaded through `Left` pointers.

**Decode:** pop root `(1,3)` → node `1` needs 3 children. Pop `(3,2)` → child of `1`, needs 2. Pop `(2,0)` → child of `1`. Pop `(4,0)` → child of `1` (now satisfied). Pop `(5,0)`, `(6,0)` → children of `3`. Reconstructed: `1[3[5,6],2,4]`. ✔

---

## Key Takeaways

- **"Encode X to Y" = design a reversible codec.** Any representation works as long as encode∘decode is the identity — pick the one that is easiest to prove correct. Here two very different encodings (structural pointers vs. a serialised stream) both satisfy the contract.
- **Left-child / right-sibling** is *the* standard way to map a multi-way tree onto a binary tree: `Left` = first child, `Right` = next sibling. Memorise it — it also shows up in the Fibonacci/binomial-heap literature and in "flatten a multilevel list" problems.
- The mapping turns *sibling breadth* into *right-going depth*: N-ary children become a right-linked list, which is why the binary tree can be tall and skinny even when the N-ary tree is shallow and wide.
- **Statelessness matters:** the constraint forbids storing the tree in a field, so recursion (its own stack) or a locally-built queue is the right tool — never a struct-level cache.

---

## Related Problems

- LeetCode #428 — Serialize and Deserialize N-ary Tree (direct serialization codec)
- LeetCode #297 — Serialize and Deserialize Binary Tree (same "reversible codec" design)
- LeetCode #449 — Serialize and Deserialize BST (BST-specialised codec)
- LeetCode #589 — N-ary Tree Preorder Traversal (the DFS this encode is built on)
- LeetCode #429 — N-ary Tree Level Order Traversal (the BFS the serialization variant uses)
- LeetCode #114 — Flatten Binary Tree to Linked List (same "reuse the two pointers" idea)
