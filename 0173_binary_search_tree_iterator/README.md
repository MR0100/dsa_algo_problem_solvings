# 0173 — Binary Search Tree Iterator

> LeetCode #173 · Difficulty: Medium
> **Categories:** Stack, Tree, Design, Binary Search Tree, Binary Tree, Iterator

---

## Problem Statement

Implement the `BSTIterator` class that represents an iterator over the [**in-order traversal**](https://en.wikipedia.org/wiki/Tree_traversal#In-order_(LNR)) of a binary search tree (BST):

- `BSTIterator(TreeNode root)` Initializes an object of the `BSTIterator` class. The `root` of the BST is given as part of the constructor. The pointer should be initialized to a non-existent number smaller than any element in the BST.
- `boolean hasNext()` Returns `true` if there exists a number in the traversal to the right of the pointer, otherwise returns `false`.
- `int next()` Moves the pointer to the right, then returns the number at the pointer.

Notice that by initializing the pointer to a non-existent smallest number, the first call to `next()` will return the smallest element in the BST.

You may assume that `next()` calls will always be valid. That is, there will be at least a next number in the in-order traversal when `next()` is called.

**Example 1:**

```
        7
       / \
      3   15
          / \
         9   20

Input
["BSTIterator", "next", "next", "hasNext", "next", "hasNext", "next", "hasNext", "next", "hasNext"]
[[[7, 3, 15, null, null, 9, 20]], [], [], [], [], [], [], [], [], []]
Output
[null, 3, 7, true, 9, true, 15, true, 20, false]

Explanation
BSTIterator bSTIterator = new BSTIterator([7, 3, 15, null, null, 9, 20]);
bSTIterator.next();    // return 3
bSTIterator.next();    // return 7
bSTIterator.hasNext(); // return True
bSTIterator.next();    // return 9
bSTIterator.hasNext(); // return True
bSTIterator.next();    // return 15
bSTIterator.hasNext(); // return True
bSTIterator.next();    // return 20
bSTIterator.hasNext(); // return False
```

**Constraints:**

- The number of nodes in the tree is in the range `[1, 10^5]`.
- `0 <= Node.val <= 10^6`
- At most `10^5` calls will be made to `hasNext`, and `next`.

**Follow-up:** Could you implement `next()` and `hasNext()` to run in average `O(1)` time and use `O(h)` memory, where `h` is the height of the tree?

---

## Company Frequency

| Company    | Frequency        | Last Reported |
|------------|------------------|---------------|
| Meta       | ★★★★★ Very High  | 2024          |
| Microsoft  | ★★★★☆ High       | 2024          |
| Amazon     | ★★★★☆ High       | 2024          |
| Google     | ★★★☆☆ Medium     | 2024          |
| LinkedIn   | ★★★☆☆ Medium     | 2023          |
| Bloomberg  | ★★☆☆☆ Low        | 2023          |

> ⚠️ Frequency data is crowd-sourced from LeetCode Discuss, Glassdoor, and
> community interview reports. Treat as a signal, not a guarantee.

---

## DSA Concepts Used

- **Binary Search Tree** — inorder traversal of a BST visits values in ascending sorted order; the iterator streams that order → see [`/dsa/binary_search_tree.md`](/dsa/binary_search_tree.md)
- **Tree Traversal (inorder)** — all three approaches are inorder traversal, run eagerly, paused with a stack, or threaded through the tree → see [`/dsa/tree_traversal.md`](/dsa/tree_traversal.md)
- **Stack** — the optimal solution replaces the recursion call stack with an explicit stack so the traversal can be frozen between calls → see [`/dsa/stack.md`](/dsa/stack.md)
- **Design / Data Structure APIs** — amortized-cost reasoning across a sequence of operations, iterator invariants → see [`/dsa/design_data_structures.md`](/dsa/design_data_structures.md)

---

## Approaches Overview

| # | Approach | Time | Space | When to use |
|---|----------|------|-------|-------------|
| 1 | Brute Force (flatten to array) | ctor O(n); `next`/`hasNext` O(1) | O(n) | Simplest correct code; fine when n is small or all values get consumed anyway |
| 2 | Controlled Recursion with a Stack (Optimal) | ctor O(h); `next` amortized O(1); `hasNext` O(1) | O(h) | Always — exactly meets the follow-up (O(h) memory, average O(1) calls) |
| 3 | Morris Threading | `next` amortized O(1); `hasNext` O(1) | O(1) | When memory is critical *and* temporarily mutating the tree is acceptable |

---

## Approach 1 — Brute Force (Flatten to Array)

### Intuition

Inorder traversal of a BST yields its values in sorted order. The laziest correct iterator: run the *entire* traversal in the constructor, store the values in a slice, and let `next()` / `hasNext()` be array-index operations. It is perfectly O(1) per call — but it front-loads O(n) work and pins O(n) memory even if the caller reads only two values, which is precisely what the follow-up's O(h) bound is designed to rule out.

### Algorithm

1. **Constructor:** run a recursive inorder traversal (`left → node → right`), appending every value to `vals`; set `idx = 0`.
2. **next():** return `vals[idx]`, then increment `idx`.
3. **hasNext():** return `idx < len(vals)`.

### Complexity

- **Time:** constructor O(n) — visits every node once; `next()` and `hasNext()` O(1) — pure index arithmetic.
- **Space:** O(n) — the flattened array of all n values (plus O(h) recursion stack during construction).

### Code

```go
type FlattenIterator struct {
	vals []int // full inorder sequence, i.e. sorted BST values
	idx  int   // position of the next unreturned value
}

func NewFlattenIterator(root *TreeNode) *FlattenIterator {
	it := &FlattenIterator{}
	var inorder func(*TreeNode)
	inorder = func(n *TreeNode) {
		if n == nil {
			return
		}
		inorder(n.Left)                  // everything smaller first
		it.vals = append(it.vals, n.Val) // then the node itself
		inorder(n.Right)                 // then everything larger
	}
	inorder(root)
	return it
}

func (it *FlattenIterator) Next() int {
	v := it.vals[it.idx] // problem guarantees Next() is only called when valid
	it.idx++             // advance the pointer past the returned value
	return v
}

func (it *FlattenIterator) HasNext() bool { return it.idx < len(it.vals) }
```

### Dry Run

Example 1: tree `[7,3,15,null,null,9,20]`. Constructor flattens inorder → `vals = [3,7,9,15,20]`, `idx = 0`.

| Step | Operation | idx before | Returned | idx after |
|------|-----------|------------|----------|-----------|
| 1 | `BSTIterator(...)` | — | null | 0 |
| 2 | `next()` | 0 | vals[0] = **3** | 1 |
| 3 | `next()` | 1 | vals[1] = **7** | 2 |
| 4 | `hasNext()` | 2 | 2 < 5 → **true** | 2 |
| 5 | `next()` | 2 | vals[2] = **9** | 3 |
| 6 | `hasNext()` | 3 | 3 < 5 → **true** | 3 |
| 7 | `next()` | 3 | vals[3] = **15** | 4 |
| 8 | `hasNext()` | 4 | 4 < 5 → **true** | 4 |
| 9 | `next()` | 4 | vals[4] = **20** | 5 |
| 10 | `hasNext()` | 5 | 5 < 5 → **false** | 5 |

Output `[null,3,7,true,9,true,15,true,20,false]` ✔

---

## Approach 2 — Controlled Recursion with a Stack (Optimal)

### Intuition

Recursive inorder is "go left as deep as possible, visit, then recurse right". Replace the call stack with our own explicit stack and the traversal can be **frozen between `next()` calls**. The invariant that makes it work: *the stack always holds exactly the nodes whose left subtrees are fully consumed but which are themselves unvisited* — stacked along a left-descending path, so the top of the stack is always the next-smallest value. Popping a node "unlocks" its right subtree; pushing that subtree's left spine restores the invariant. Each node is pushed once and popped once over the whole iteration, so `next()` is O(1) amortized even though a single call can push O(h) nodes.

### Algorithm

1. **Constructor:** push the left spine of `root` (root, root.Left, root.Left.Left, …) onto the stack.
2. **next():**
   1. Pop the top node — it is the next-smallest; call it `top`.
   2. Push the left spine of `top.Right` (if any).
   3. Return `top.Val`.
3. **hasNext():** return `len(stack) > 0`.

### Complexity

- **Time:** constructor O(h); `next()` **amortized O(1)** — across the entire iteration every node is pushed exactly once and popped exactly once (a single call is worst-case O(h)); `hasNext()` O(1).
- **Space:** O(h) — the stack holds at most one root-to-leaf path (h = log n balanced, n skewed).

### Code

```go
type StackIterator struct {
	stack []*TreeNode // partially-frozen inorder traversal; top = next value
}

func NewStackIterator(root *TreeNode) *StackIterator {
	it := &StackIterator{}
	it.pushLeftSpine(root) // prime the stack so the minimum is on top
	return it
}

// pushLeftSpine pushes node and all its left descendants onto the stack.
func (it *StackIterator) pushLeftSpine(node *TreeNode) {
	for node != nil {
		it.stack = append(it.stack, node) // owe this node's value later
		node = node.Left                  // its left subtree comes first
	}
}

func (it *StackIterator) Next() int {
	top := it.stack[len(it.stack)-1]      // smallest unvisited node
	it.stack = it.stack[:len(it.stack)-1] // pop it — we are visiting it now
	it.pushLeftSpine(top.Right)           // its right subtree is next in line
	return top.Val
}

func (it *StackIterator) HasNext() bool { return len(it.stack) > 0 }
```

### Dry Run

Example 1: tree `[7,3,15,null,null,9,20]` (stack shown bottom→top).

| Step | Operation | Action | Stack after | Returned |
|------|-----------|--------|-------------|----------|
| 1 | `BSTIterator(...)` | push left spine of 7: push 7, push 3 | [7, 3] | null |
| 2 | `next()` | pop 3; 3.Right = nil → nothing pushed | [7] | **3** |
| 3 | `next()` | pop 7; push left spine of 15: push 15, push 9 | [15, 9] | **7** |
| 4 | `hasNext()` | stack non-empty | [15, 9] | **true** |
| 5 | `next()` | pop 9; 9.Right = nil | [15] | **9** |
| 6 | `hasNext()` | stack non-empty | [15] | **true** |
| 7 | `next()` | pop 15; push left spine of 20: push 20 | [20] | **15** |
| 8 | `hasNext()` | stack non-empty | [20] | **true** |
| 9 | `next()` | pop 20; 20.Right = nil | [] | **20** |
| 10 | `hasNext()` | stack empty | [] | **false** |

Output `[null,3,7,true,9,true,15,true,20,false]` ✔

---

## Approach 3 — Morris Threading (O(1) Extra Space)

### Intuition

The stack exists only so we can climb back up after finishing a left subtree. **Morris traversal stores that return path inside the tree itself**: before descending into a left subtree, find the current node's inorder predecessor (the rightmost node of the left subtree) and point its nil right pointer back at the current node — a *thread*. When the traversal later walks off the predecessor's right pointer, it lands exactly on the node where it must resume; arriving there a second time proves the left subtree is finished, so the thread is removed (restoring the tree) and the node is emitted. The only state between `next()` calls is one pointer — O(1) memory. Trade-off: the tree is temporarily mutated mid-iteration (not acceptable for concurrent readers), and is fully restored only once iteration runs to completion.

### Algorithm

One `next()` call — loop on `cur`:

1. If `cur.Left == nil`: `cur` is the next value — save `cur.Val`, move `cur = cur.Right` (real child *or* thread to the successor), return the saved value.
2. Otherwise find `pred` = rightmost node of `cur.Left`, stopping early if `pred.Right` already threads to `cur`.
3. If `pred.Right == nil` (first arrival): lay the thread `pred.Right = cur`; descend `cur = cur.Left`.
4. Else (second arrival — left subtree consumed): remove the thread (`pred.Right = nil`), save `cur.Val`, move `cur = cur.Right`, return the saved value.

`hasNext()`: `cur != nil`.

### Complexity

- **Time:** `next()` amortized O(1) — over the full iteration every edge is traversed at most twice (once laying a thread, once removing it), ≈ 4n pointer moves for n values; `hasNext()` O(1).
- **Space:** O(1) — a single resume pointer; the traversal bookkeeping lives in the tree's own nil right pointers (tree mutated during iteration, restored by completion).

### Code

```go
type MorrisIterator struct {
	cur *TreeNode // where the paused traversal will resume
}

func NewMorrisIterator(root *TreeNode) *MorrisIterator {
	return &MorrisIterator{cur: root}
}

func (it *MorrisIterator) Next() int {
	for it.cur != nil {
		if it.cur.Left == nil {
			// No left subtree → cur itself is the next inorder value.
			v := it.cur.Val
			it.cur = it.cur.Right // real child or thread to the successor
			return v
		}
		// Find cur's inorder predecessor: rightmost node in the left subtree.
		pred := it.cur.Left
		for pred.Right != nil && pred.Right != it.cur {
			pred = pred.Right
		}
		if pred.Right == nil {
			pred.Right = it.cur  // lay a thread so we can return to cur later
			it.cur = it.cur.Left // now safe to descend into the left subtree
		} else {
			pred.Right = nil      // second arrival: left subtree finished — unthread
			v := it.cur.Val       // cur is the next inorder value
			it.cur = it.cur.Right // continue with the right subtree
			return v
		}
	}
	return null // unreachable: problem guarantees Next() calls are valid
}

func (it *MorrisIterator) HasNext() bool { return it.cur != nil }
```

### Dry Run

Example 1: tree `[7,3,15,null,null,9,20]`, `cur = 7` after construction.

| Step | Operation | cur | Action | Tree threads | Returned |
|------|-----------|-----|--------|--------------|----------|
| 1 | `BSTIterator(...)` | 7 | store root pointer | none | null |
| 2 | `next()` | 7 | 7 has left; pred(7) = 3, 3.Right = nil → thread 3→7, descend | 3.Right⇢7 | — |
|   |          | 3 | 3 has no left → emit 3, follow 3.Right (the thread) | 3.Right⇢7 | **3** |
| 3 | `next()` | 7 | 7 has left; pred(7) = 3 but 3.Right == 7 → thread found: unthread, emit 7, go right | none | **7** |
| 4 | `hasNext()` | 15 | cur ≠ nil | none | **true** |
| 5 | `next()` | 15 | 15 has left; pred(15) = 9, 9.Right = nil → thread 9→15, descend | 9.Right⇢15 | — |
|   |          | 9 | 9 has no left → emit 9, follow thread | 9.Right⇢15 | **9** |
| 6 | `hasNext()` | 15 | cur ≠ nil | 9.Right⇢15 | **true** |
| 7 | `next()` | 15 | pred(15) = 9, 9.Right == 15 → unthread, emit 15, go right | none | **15** |
| 8 | `hasNext()` | 20 | cur ≠ nil | none | **true** |
| 9 | `next()` | 20 | 20 has no left → emit 20, cur = 20.Right = nil | none | **20** |
| 10 | `hasNext()` | nil | cur == nil | none | **false** |

Output `[null,3,7,true,9,true,15,true,20,false]` ✔ — and every thread laid was later removed, so the tree ends unmodified.

---

## Key Takeaways

- **"Iterator over a traversal" = pause a recursion.** The reusable trick is replacing the call stack with an explicit stack whose invariant is *top = next element*; `pushLeftSpine` is the whole secret and reappears in #285 (Inorder Successor), #230 (Kth Smallest), and two-iterator problems like #272.
- **Amortized analysis is the expected talking point:** one `next()` may push h nodes, but n calls do exactly n pushes + n pops total → O(1) average. Say "amortized", not "always O(1)".
- Stack top is always the minimum unvisited node — this iterator is effectively a lazy sorted stream, which is why "merge two BSTs" style problems build directly on it.
- **Morris threading** converts O(h) stack into O(1) pointers by borrowing nil right pointers as return addresses; know its caveats (mutates the tree mid-iteration, unsafe with concurrent readers, restored only at completion).
- The flatten approach is not "wrong" — it is the right call when all elements will be consumed anyway and h could be as bad as n (skewed tree); state the trade-off instead of dismissing it.

---

## Related Problems

- LeetCode #94 — Binary Tree Inorder Traversal (the underlying traversal, all three flavors)
- LeetCode #230 — Kth Smallest Element in a BST (stop this iterator after k `next()` calls)
- LeetCode #285 — Inorder Successor in BST (single `next()` from an arbitrary node)
- LeetCode #251 — Flatten 2D Vector (same design-an-iterator pattern)
- LeetCode #284 — Peeking Iterator (iterator adapter design)
- LeetCode #341 — Flatten Nested List Iterator (stack-paused traversal of a nested structure)
