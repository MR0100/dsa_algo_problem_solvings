# 0138 — Copy List with Random Pointer

> LeetCode #138 · Difficulty: Medium
> **Categories:** Hash Table, Linked List

---

## Problem Statement

A linked list of length `n` is given such that each node contains an additional random pointer, which could point to any node in the list, or `null`.

Construct a **deep copy** of the list. The deep copy should consist of exactly `n` **brand new** nodes, where each new node has its value set to the value of its corresponding original node. Both the `next` and `random` pointer of the new nodes should point to new nodes in the copied list such that the pointers in the original list and copied list represent the same list state. **None of the pointers in the new list should point to nodes in the original list.**

For example, if there are two nodes `X` and `Y` in the original list, where `X.random --> Y`, then for the corresponding two nodes `x` and `y` in the copied list, `x.random --> y`.

Return the head of the copied linked list.

The linked list is represented in the input/output as a list of `n` nodes. Each node is represented as a pair of `[val, random_index]` where:

- `val`: an integer representing `Node.val`
- `random_index`: the index of the node (range from `0` to `n-1`) that the `random` pointer points to, or `null` if it does not point to any node.

Your code will **only** be given the `head` of the original linked list.

**Example 1:**
```
Input: head = [[7,null],[13,0],[11,4],[10,2],[1,0]]
Output: [[7,null],[13,0],[11,4],[10,2],[1,0]]
```

**Example 2:**
```
Input: head = [[1,1],[2,1]]
Output: [[1,1],[2,1]]
```

**Example 3:**
```
Input: head = [[3,null],[3,0],[3,null]]
Output: [[3,null],[3,0],[3,null]]
```

**Constraints:**
- `0 <= n <= 1000`
- `-10^4 <= Node.val <= 10^4`
- `Node.random` is `null` or is pointing to some node in the linked list.

---

## Company Frequency

| Company    | Frequency       | Last Reported |
|------------|-----------------|---------------|
| Amazon     | ★★★★★ Very High | 2024          |
| Facebook   | ★★★★★ Very High | 2024          |
| Microsoft  | ★★★★☆ High      | 2024          |
| Bloomberg  | ★★★★☆ High      | 2024          |
| Google     | ★★★☆☆ Medium    | 2023          |
| Oracle     | ★★★☆☆ Medium    | 2023          |

> ⚠️ Frequency data is crowd-sourced from LeetCode Discuss, Glassdoor, and community interview reports. Treat as a signal, not a guarantee.

---

## DSA Concepts Used

- **Linked List** — pointer surgery: weaving and unweaving interleaved nodes → see [`/dsa/linked_list.md`](/dsa/linked_list.md)
- **Hash Map** — the original→copy identity map is the core of the straightforward solutions → see [`/dsa/hash_map.md`](/dsa/hash_map.md)
- **Graph DFS (clone-graph pattern)** — random pointers make the list a graph; memoized DFS clones it → see [`/dsa/graph_bfs_dfs.md`](/dsa/graph_bfs_dfs.md)

---

## Approaches Overview

| # | Approach                 | Time | Space | When to use                                              |
|---|--------------------------|------|-------|----------------------------------------------------------|
| 1 | Hash Map (Two Pass)      | O(n) | O(n)  | Default interview answer; clearest to explain            |
| 2 | Recursion + Memoization  | O(n) | O(n)  | When you want to show the clone-graph generalization     |
| 3 | Interleaving (Optimal)   | O(n) | O(1)  | Follow-up "can you do O(1) extra space?"                 |

---

## Approach 1 — Hash Map (Two Pass)

### Intuition
A naive one-pass clone breaks because `random` may point to a node whose copy doesn't exist yet (a forward jump). The fix is to decouple *node creation* from *pointer wiring*: create every copy first, remembering which copy belongs to which original in a map, then wire all `next`/`random` pointers by translating each original pointer through the map.

### Algorithm
1. If `head` is nil, return nil.
2. **Pass 1:** walk the list; for every original node `curr`, create `&Node{Val: curr.Val}` and store `oldToNew[curr] = copy`.
3. **Pass 2:** walk again; for each `curr` set:
   - `oldToNew[curr].Next = oldToNew[curr.Next]`
   - `oldToNew[curr].Random = oldToNew[curr.Random]`
   (a nil original pointer looks up to Go's zero value nil — the tail and null randoms wire themselves.)
4. Return `oldToNew[head]`.

### Complexity
- **Time:** O(n) — two passes; each map get/put is amortized O(1).
- **Space:** O(n) — the map keeps one entry per node (beyond the mandatory output list).

### Code
```go
func hashMapTwoPass(head *Node) *Node {
	if head == nil {
		return nil // empty list clones to empty list
	}
	oldToNew := make(map[*Node]*Node) // original node → its copy

	// Pass 1: create all copies so every target of next/random already exists.
	for curr := head; curr != nil; curr = curr.Next {
		oldToNew[curr] = &Node{Val: curr.Val} // copy carries only the value for now
	}

	// Pass 2: wire pointers by translating originals through the map.
	for curr := head; curr != nil; curr = curr.Next {
		copyNode := oldToNew[curr]
		copyNode.Next = oldToNew[curr.Next]     // nil key → nil value, handles tail
		copyNode.Random = oldToNew[curr.Random] // nil key → nil value, handles null random
	}

	return oldToNew[head] // the copy of the head is the new list
}
```

### Dry Run
`head = [[7,null],[13,0],[11,4],[10,2],[1,0]]` (Example 1). Call the original nodes `A(7) B(13) C(11) D(10) E(1)` and their copies `A'..E'`.

Pass 1 (creation):

| step | node | map after                                  |
|------|------|--------------------------------------------|
| 1    | A    | {A→A'}                                     |
| 2    | B    | {A→A', B→B'}                               |
| 3    | C    | {A→A', B→B', C→C'}                         |
| 4    | D    | {A→A', B→B', C→C', D→D'}                   |
| 5    | E    | {A→A', B→B', C→C', D→D', E→E'}             |

Pass 2 (wiring):

| step | node | orig.Next | orig.Random | copy.Next set to | copy.Random set to |
|------|------|-----------|-------------|------------------|--------------------|
| 1    | A    | B         | nil         | B'               | nil                |
| 2    | B    | C         | A           | C'               | A'                 |
| 3    | C    | D         | E           | D'               | E'                 |
| 4    | D    | E         | C           | E'               | C'                 |
| 5    | E    | nil       | A           | nil              | A'                 |

Result `A'→B'→C'→D'→E'` serializes to `[[7,null],[13,0],[11,4],[10,2],[1,0]]` ✅

---

## Approach 2 — Recursion + Memoization

### Intuition
Each node has two outgoing edges (`next` and `random`), so the structure is really a directed graph — and this problem is *Clone Graph* in disguise. Deep-cloning a graph is memoized DFS: the memo map guarantees each original node produces exactly one copy, and inserting into the memo **before** recursing breaks the cycles that `random` (or `random`-induced loops) can create.

### Algorithm
1. Define `clone(node)`:
   1. If `node` is nil, return nil.
   2. If `node` is already in `memo`, return the stored copy.
   3. Create `copy = &Node{Val: node.Val}` and set `memo[node] = copy` **now** (before recursion).
   4. `copy.Next = clone(node.Next)`; `copy.Random = clone(node.Random)`.
   5. Return `copy`.
2. Return `clone(head)`.

### Complexity
- **Time:** O(n) — each node is cloned once; repeat visits are O(1) memo hits.
- **Space:** O(n) — memo map plus recursion stack up to n frames (the `next` chain).

### Code
```go
func recursiveMemo(head *Node) *Node {
	memo := make(map[*Node]*Node) // original → copy, shared across the recursion
	var clone func(node *Node) *Node
	clone = func(node *Node) *Node {
		if node == nil {
			return nil // base case: null pointer clones to null
		}
		if copyNode, ok := memo[node]; ok {
			return copyNode // already cloned (revisited via random or a cycle)
		}
		copyNode := &Node{Val: node.Val}
		memo[node] = copyNode           // memoize BEFORE recursing to break cycles
		copyNode.Next = clone(node.Next)     // deep-clone the next edge
		copyNode.Random = clone(node.Random) // deep-clone the random edge
		return copyNode
	}
	return clone(head)
}
```

### Dry Run
`head = [[7,null],[13,0],[11,4],[10,2],[1,0]]` (Example 1), nodes `A(7) B(13) C(11) D(10) E(1)`:

| call            | memo hit? | action                                              | memo after           |
|-----------------|-----------|-----------------------------------------------------|----------------------|
| clone(A)        | no        | create A', recurse Next→B                           | {A}                  |
| clone(B)        | no        | create B', recurse Next→C                           | {A,B}                |
| clone(C)        | no        | create C', recurse Next→D                           | {A,B,C}              |
| clone(D)        | no        | create D', recurse Next→E                           | {A,B,C,D}            |
| clone(E)        | no        | create E', Next=clone(nil)=nil, Random=clone(A)     | {A,B,C,D,E}          |
| clone(A) via E  | **yes**   | return existing A' (cycle safely broken)            | unchanged            |
| unwind D        | —         | D'.Random = clone(C) → memo hit, C'                 | unchanged            |
| unwind C        | —         | C'.Random = clone(E) → memo hit, E'                 | unchanged            |
| unwind B        | —         | B'.Random = clone(A) → memo hit, A'                 | unchanged            |
| unwind A        | —         | A'.Random = clone(nil) = nil                        | unchanged            |

Output: `[[7,null],[13,0],[11,4],[10,2],[1,0]]` ✅

---

## Approach 3 — Interleaving (Optimal, O(1) extra space)

### Intuition
The hash map answers exactly one question: *"where is the copy of this original node?"* We can encode that answer **in the list structure itself**: insert each copy immediately after its original, producing `A→A'→B→B'→C→C'…`. Now the copy of any node `X` is simply `X.Next`, so `X'.Random = X.Random.Next` — no map at all. A final pass unweaves the two lists and restores the original.

### Algorithm
1. **Weave:** for each original `curr`, insert `curr' = &Node{Val: curr.Val}` between `curr` and `curr.Next`; advance `curr = curr.Next.Next`.
2. **Assign randoms:** for each original `curr` (stepping by 2): if `curr.Random != nil`, set `curr.Next.Random = curr.Random.Next`.
3. **Unweave:** `newHead = head.Next`; for each original `curr`: `curr.Next = curr.Next.Next` (restore original), and link its copy to the following copy. Return `newHead`.

### Complexity
- **Time:** O(n) — three linear passes.
- **Space:** O(1) — a constant number of pointers; the only allocations are the required output nodes.

### Code
```go
func interleaving(head *Node) *Node {
	if head == nil {
		return nil // nothing to copy
	}

	// Step 1: weave copies into the original list: A→A'→B→B'→...
	for curr := head; curr != nil; curr = curr.Next.Next {
		copyNode := &Node{Val: curr.Val}
		copyNode.Next = curr.Next // copy points at the rest of the list
		curr.Next = copyNode      // original points at its copy
	}

	// Step 2: set random pointers on the copies using the weave invariant.
	for curr := head; curr != nil; curr = curr.Next.Next {
		if curr.Random != nil {
			curr.Next.Random = curr.Random.Next // copy of X.Random is X.Random.Next
		}
	}

	// Step 3: unweave — detach the copies and restore the original list.
	newHead := head.Next
	for curr := head; curr != nil; {
		copyNode := curr.Next
		curr.Next = copyNode.Next // restore original's next
		if copyNode.Next != nil {
			copyNode.Next = copyNode.Next.Next // link copy to the next copy
		}
		curr = curr.Next // move to the next original node
	}

	return newHead
}
```

### Dry Run
`head = [[7,null],[13,0],[11,4],[10,2],[1,0]]` (Example 1), nodes `A(7) B(13) C(11) D(10) E(1)` with `A.rnd=nil, B.rnd=A, C.rnd=E, D.rnd=C, E.rnd=A`.

**Step 1 — weave:**

| iteration | insert | list state after                                  |
|-----------|--------|---------------------------------------------------|
| 1         | A'     | A→A'→B→C→D→E                                      |
| 2         | B'     | A→A'→B→B'→C→D→E                                   |
| 3         | C'     | A→A'→B→B'→C→C'→D→E                                |
| 4         | D'     | A→A'→B→B'→C→C'→D→D'→E                             |
| 5         | E'     | A→A'→B→B'→C→C'→D→D'→E→E'                          |

**Step 2 — randoms** (`X'.Random = X.Random.Next`):

| original | X.Random | X.Random.Next | assignment       |
|----------|----------|---------------|------------------|
| A        | nil      | —             | A'.Random = nil  |
| B        | A        | A'            | B'.Random = A'   |
| C        | E        | E'            | C'.Random = E'   |
| D        | C        | C'            | D'.Random = C'   |
| E        | A        | A'            | E'.Random = A'   |

**Step 3 — unweave:**

| iteration | original restored | copy chained  |
|-----------|-------------------|---------------|
| 1         | A.Next = B        | A'.Next = B'  |
| 2         | B.Next = C        | B'.Next = C'  |
| 3         | C.Next = D        | C'.Next = D'  |
| 4         | D.Next = E        | D'.Next = E'  |
| 5         | E.Next = nil      | E'.Next = nil |

Return `A'` → `[[7,null],[13,0],[11,4],[10,2],[1,0]]`, original list intact ✅

---

## Key Takeaways

- **Identity map (old→new) is the universal deep-copy pattern** — same trick clones graphs (#133), trees with random pointers (#1485), and any pointer-rich structure.
- Two-phase "create all, then wire all" sidesteps forward references entirely.
- When memoizing recursive clones, **insert into the memo before recursing** or cycles recurse forever.
- The interleaving trick shows a recurring space optimization: **encode an auxiliary mapping inside the data structure itself** (copy of X lives at `X.Next`), then undo the mutation before returning.
- Always restore the input list in the O(1) approach — interviewers check that the original is unmodified.

---

## Related Problems

- LeetCode #133 — Clone Graph (same memoized-DFS deep-copy pattern)
- LeetCode #1485 — Clone Binary Tree With Random Pointer (tree version)
- LeetCode #1490 — Clone N-ary Tree (simpler clone)
- LeetCode #141 — Linked List Cycle (pointer-weaving intuition)
- LeetCode #430 — Flatten a Multilevel Doubly Linked List (multi-pointer list surgery)
