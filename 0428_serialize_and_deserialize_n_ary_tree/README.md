# 0428 — Serialize and Deserialize N-ary Tree

> LeetCode #428 · Difficulty: Hard
> **Categories:** Tree, Depth-First Search, Breadth-First Search, Design, String

---

## Problem Statement

Serialization is the process of converting a data structure or object into a sequence of bits so that it can be stored in a file or memory buffer, or transmitted across a network connection link to be reconstructed later in the same or another computer environment.

Design an algorithm to serialize and deserialize an **N-ary tree**. An N-ary tree is a rooted tree in which each node has no more than N children. There is no restriction on how your serialization/deserialization algorithm should work. You just need to ensure that an N-ary tree can be serialized to a string and this string can be deserialized to the original tree structure.

For example, you may serialize the following `3-ary` tree
```
        1
      / | \
     3  2  4
    / \
   5   6
```
as `[1 [3[5 6] 2 4]]`. Note that this is just an example, you do not necessarily need to follow this format.

Or you can follow LeetCode's level order traversal serialization format, where each group of children is separated by the `null` value:
```
Input: root = [1,null,2,3,4,5,null,null,6,7,null,8,null,9,10,null,null,11,null,12,null,13,null,null,14]
```

**Constraints:**
- The height of the n-ary tree is less than or equal to `1000`.
- The total number of nodes is between `[0, 10⁴]`.
- Do not use class member/global/static variables to store states. Your serialize and deserialize algorithms should be stateless.

---

## Company Frequency

| Company    | Frequency        | Last Reported |
|------------|------------------|---------------|
| Amazon     | ★★★☆☆ Medium     | 2023          |
| Google     | ★★★☆☆ Medium     | 2023          |
| Facebook   | ★★☆☆☆ Low        | 2023          |
| Microsoft  | ★★☆☆☆ Low        | 2022          |
| LinkedIn   | ★★☆☆☆ Low        | 2022          |

> ⚠️ Frequency data is crowd-sourced from LeetCode Discuss, Glassdoor, and
> community interview reports. Treat as a signal, not a guarantee.

---

## DSA Concepts Used

- **Tree Traversal (preorder DFS)** — both codecs write and read the tree in a single preorder walk → see [`/dsa/tree_traversal.md`](/dsa/tree_traversal.md)
- **Design / stateless codec** — the API is a `serialize`/`deserialize` pair judged on round-trip fidelity; state must live in the string, not in fields → see [`/dsa/design_data_structures.md`](/dsa/design_data_structures.md)
- **String parsing / tokenization** — deserialization tokenizes the encoded string (comma split, or a hand-rolled bracket lexer) → see [`/dsa/string_algorithms.md`](/dsa/string_algorithms.md)

---

## Approaches Overview

Let n = number of nodes, h = tree height.

| # | Approach | Time (ser/deser) | Space | When to use |
|---|----------|------------------|-------|-------------|
| 1 | Preorder with explicit child counts | O(n) / O(n) | O(n) + O(h) | Simplest to get correct; count removes all ambiguity, no sentinel needed |
| 2 | Bracketed / parenthesized encoding | O(n) / O(n) | O(n) + O(h) | Human-readable (S-expression style); structure carried by matching brackets |

---

## Approach 1 — Preorder with Explicit Child Counts

### Intuition
Serializing a *binary* tree is easy because each node has a fixed slot for left and right; you can mark absent children with a sentinel. An N-ary node has an **unbounded** child list, so a reader cannot tell where one node's children stop. The clean fix: write the **child count** immediately after each value. A preorder stream `val,count,...` is then unambiguous — read a value and its `k`, then recursively read exactly `k` children.

### Algorithm
**serialize:**
1. Preorder walk; for each node emit `Val` then `len(Children)`.
2. Recurse into every child in order.
3. Join tokens with commas.

**deserialize:**
1. Split on commas; keep a shared cursor `pos`.
2. `build()`: read `Val` and `count` (advance `pos` by 2), make the node, then call `build()` exactly `count` times for its children.

### Complexity
- **Time:** O(n) each way — every node is emitted once and consumed once.
- **Space:** O(n) for the string plus O(h) recursion depth.

### Code
```go
type Codec1 struct{}

// serialize turns the tree into "v0,c0,v1,c1,..." preorder token stream.
func (Codec1) serialize(root *Node) string {
	if root == nil {
		return "" // empty tree ⇒ empty string
	}
	tokens := []string{}
	var pre func(n *Node)
	pre = func(n *Node) {
		// Each node contributes its value AND its child count so the reader
		// knows how many children to pull next.
		tokens = append(tokens, strconv.Itoa(n.Val), strconv.Itoa(len(n.Children)))
		for _, c := range n.Children {
			pre(c)
		}
	}
	pre(root)
	return strings.Join(tokens, ",")
}

// deserialize rebuilds the tree from the "v,c,v,c,..." token stream.
func (Codec1) deserialize(data string) *Node {
	if data == "" {
		return nil
	}
	tokens := strings.Split(data, ",")
	pos := 0 // shared cursor into tokens
	var build func() *Node
	build = func() *Node {
		val, _ := strconv.Atoi(tokens[pos])     // current node's value
		count, _ := strconv.Atoi(tokens[pos+1]) // its child count
		pos += 2                                // consume the (val,count) pair
		node := &Node{Val: val, Children: []*Node{}}
		for i := 0; i < count; i++ {
			node.Children = append(node.Children, build()) // pull exactly `count` children
		}
		return node
	}
	return build()
}
```

### Dry Run (Example tree `1 → (3,2,4)`, `3 → (5,6)`)

**serialize** (preorder, emitting `val,count`):

| Visit | Emit | tokens so far |
|-------|------|---------------|
| 1 (3 children) | `1,3` | `1,3` |
| 3 (2 children) | `3,2` | `1,3,3,2` |
| 5 (0) | `5,0` | `1,3,3,2,5,0` |
| 6 (0) | `6,0` | `1,3,3,2,5,0,6,0` |
| 2 (0) | `2,0` | `…,2,0` |
| 4 (0) | `4,0` | `1,3,3,2,5,0,6,0,2,0,4,0` |

**deserialize** `1,3,3,2,5,0,6,0,2,0,4,0`:

| `pos` | Read (val,count) | Action |
|-------|------------------|--------|
| 0 | (1,3) | node 1, will read 3 children |
| 2 | (3,2) | node 3, will read 2 children |
| 4 | (5,0) | leaf 5 → child of 3 |
| 6 | (6,0) | leaf 6 → child of 3; node 3 done |
| 8 | (2,0) | leaf 2 → child of 1 |
| 10 | (4,0) | leaf 4 → child of 1; node 1 done |

Rebuilt level order: `[[1],[3,2,4],[5,6]]` ✓

---

## Approach 2 — Bracketed / Parenthesized Encoding

### Intuition
Another way to delimit an unbounded child list is to **bracket** it, exactly like JSON or an S-expression: a node prints its value, and if it has children, `[` … children … `]`. The matching `]` marks where this node's children end, so no counts are needed — the structure lives in the nesting. Deserialization is then a tiny recursive-descent parser: read a number (new node); if the next token is `[`, recurse to collect children until the matching `]`.

### Algorithm
**serialize:**
1. Emit the node's value.
2. If it has children, emit `[`, recurse into each child (space-separated), emit `]`.

**deserialize:**
1. Tokenize into numbers, `[`, `]`.
2. `build()`: read a number → new node; if the next token is `[`, consume it, loop `build()` until `]`, consume `]`.

### Complexity
- **Time:** O(n) each way — one pass to emit, one pass to parse.
- **Space:** O(n) string plus O(h) recursion depth.

### Code
```go
type Codec2 struct{}

// serialize produces "val" or "val[child child ...]" recursively.
func (Codec2) serialize(root *Node) string {
	if root == nil {
		return ""
	}
	var sb strings.Builder
	var enc func(n *Node)
	enc = func(n *Node) {
		sb.WriteString(strconv.Itoa(n.Val)) // node value
		if len(n.Children) > 0 {
			sb.WriteByte('[') // open the child group
			for i, c := range n.Children {
				if i > 0 {
					sb.WriteByte(' ') // separate siblings
				}
				enc(c)
			}
			sb.WriteByte(']') // close the child group
		}
	}
	enc(root)
	return sb.String()
}

// deserialize parses the bracketed string back into a tree.
func (Codec2) deserialize(data string) *Node {
	if data == "" {
		return nil
	}
	tokens := tokenizeBrackets(data)
	pos := 0
	var build func() *Node
	build = func() *Node {
		val, _ := strconv.Atoi(tokens[pos]) // a number token starts a node
		pos++
		node := &Node{Val: val, Children: []*Node{}}
		if pos < len(tokens) && tokens[pos] == "[" {
			pos++ // consume "["
			for tokens[pos] != "]" {
				node.Children = append(node.Children, build()) // read one child
			}
			pos++ // consume the matching "]"
		}
		return node
	}
	return build()
}

// tokenizeBrackets splits a bracketed string into number / "[" / "]" tokens.
func tokenizeBrackets(s string) []string {
	tokens := []string{}
	i := 0
	for i < len(s) {
		switch s[i] {
		case '[', ']':
			tokens = append(tokens, string(s[i])) // structural token
			i++
		case ' ':
			i++ // skip separators
		default:
			j := i
			for j < len(s) && s[j] != '[' && s[j] != ']' && s[j] != ' ' {
				j++ // extend over the whole number (supports multi-digit / negatives)
			}
			tokens = append(tokens, s[i:j])
			i = j
		}
	}
	return tokens
}
```

### Dry Run (same example tree)

**serialize** produces `1[3[5 6] 2 4]`:
- `1` then `[` (has children)
- `3` then `[5 6]` (3's children), space, `2`, space, `4`
- close `]` → `1[3[5 6] 2 4]`

**deserialize** `1[3[5 6] 2 4]` → tokens `1 [ 3 [ 5 6 ] 2 4 ]`:

| `pos` | Token | Action | Stack of open nodes |
|-------|-------|--------|---------------------|
| 0 | `1` | new node 1 | [1] |
| 1 | `[` | open 1's children | [1] |
| 2 | `3` | new node 3 (child of 1) | [1,3] |
| 3 | `[` | open 3's children | [1,3] |
| 4 | `5` | leaf 5 (child of 3) | [1,3] |
| 5 | `6` | leaf 6 (child of 3) | [1,3] |
| 6 | `]` | close 3's children; 3 done | [1] |
| 7 | `2` | leaf 2 (child of 1) | [1] |
| 8 | `4` | leaf 4 (child of 1) | [1] |
| 9 | `]` | close 1's children; 1 done | [] |

Rebuilt level order: `[[1],[3,2,4],[5,6]]` ✓

---

## Key Takeaways

- **The N-ary difficulty is "where do the children end?"** Two clean answers: write the **count** (Approach 1) or write **delimiters** (Approach 2). Both make a preorder stream unambiguous; pick whichever the interviewer prefers.
- **Serialize and deserialize must share the exact grammar** — the reader mirrors the writer. Keep the two in lockstep and test with a round-trip (`deserialize(serialize(t))` reproduces `t`), which is the only correctness check that matters.
- **Stateless requirement:** carry the parse cursor as a *local* (`pos`) threaded through the recursion, never a struct field — otherwise two concurrent codecs would corrupt each other, and the problem explicitly forbids it.
- Preorder + count generalises to any tree/graph shape where a node's out-degree varies; it is the same idea behind length-prefixed framing in binary protocols.
- Handle the **empty tree** (`nil` ⇄ `""`) and **single node** explicitly — the two edge cases most likely to break a codec.

---

## Related Problems
- LeetCode #297 — Serialize and Deserialize Binary Tree (the binary-tree original)
- LeetCode #449 — Serialize and Deserialize BST (BST-specialised, more compact)
- LeetCode #431 — Encode N-ary Tree to Binary Tree (a related N-ary encoding trick)
- LeetCode #536 — Construct Binary Tree from String (bracketed-format parsing)
- LeetCode #606 — Construct String from Binary Tree (the serialize half, bracketed)
- LeetCode #429 — N-ary Tree Level Order Traversal (BFS over the same node type)
