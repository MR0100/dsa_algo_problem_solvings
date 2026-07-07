# 0133 — Clone Graph

> LeetCode #133 · Difficulty: Medium
> **Categories:** Hash Table, Depth-First Search, Breadth-First Search, Graph

---

## Problem Statement

Given a reference of a node in a **connected** undirected graph, return a **deep copy** (clone) of the graph.

Each node in the graph contains a value (`int`) and a list (`List[Node]`) of its neighbors.

```
class Node {
    public int val;
    public List<Node> neighbors;
}
```

**Test case format:**

For simplicity, each node's value is the same as the node's index (1-indexed). For example, the first node with `val == 1`, the second node with `val == 2`, and so on. The graph is represented in the test case using an adjacency list.

An **adjacency list** is a collection of unordered lists used to represent a finite graph. Each list describes the set of neighbors of a node in the graph.

The given node will always be the first node with `val = 1`. You must return the **copy of the given node** as a reference to the cloned graph.

**Example 1:**
```
Input: adjList = [[2,4],[1,3],[2,4],[1,3]]
Output: [[2,4],[1,3],[2,4],[1,3]]
Explanation: There are 4 nodes in the graph.
1st node (val = 1)'s neighbors are 2nd node (val = 2) and 4th node (val = 4).
2nd node (val = 2)'s neighbors are 1st node (val = 1) and 3rd node (val = 3).
3rd node (val = 3)'s neighbors are 2nd node (val = 2) and 4th node (val = 4).
4th node (val = 4)'s neighbors are 1st node (val = 1) and 3rd node (val = 3).
```

**Example 2:**
```
Input: adjList = [[]]
Output: [[]]
Explanation: Note that the input contains one empty list. The graph consists of
only one node with val = 1 and it does not have any neighbors.
```

**Example 3:**
```
Input: adjList = []
Output: []
Explanation: This is an empty graph, it does not have any nodes.
```

**Constraints:**
- The number of nodes in the graph is in the range `[0, 100]`.
- `1 <= Node.val <= 100`
- `Node.val` is unique for each node.
- There are no repeated edges and no self-loops in the graph.
- The Graph is connected and all nodes can be visited starting from the given node.

---

## Company Frequency

| Company   | Frequency        | Last Reported |
|-----------|------------------|---------------|
| Meta      | ★★★★★ Very High  | 2024          |
| Amazon    | ★★★★☆ High       | 2024          |
| Google    | ★★★★☆ High       | 2024          |
| Microsoft | ★★★☆☆ Medium     | 2023          |
| Uber      | ★★★☆☆ Medium     | 2023          |

> ⚠️ Frequency data is crowd-sourced from LeetCode Discuss, Glassdoor, and
> community interview reports. Treat as a signal, not a guarantee.

---

## DSA Concepts Used

- **Graph BFS/DFS** — cloning requires visiting every node and edge exactly once → see [`/dsa/graph_bfs_dfs.md`](/dsa/graph_bfs_dfs.md)
- **Hash Map** — `original → clone` map is simultaneously the visited set and the pointer-translation table → see [`/dsa/hash_map.md`](/dsa/hash_map.md)
- **Queue / Deque** — the iterative BFS frontier → see [`/dsa/queue_deque.md`](/dsa/queue_deque.md)

---

## Approaches Overview

| # | Approach | Time | Space | When to use |
|---|----------|------|-------|-------------|
| 1 | Two-Pass Copy (Brute Force) | O(V + E) | O(V) | Conceptually simplest; separates node creation from edge wiring |
| 2 | DFS + Hash Map | O(V + E) | O(V) | Shortest code; the classic interview answer |
| 3 | BFS + Hash Map (Optimal, iterative) | O(V + E) | O(V) | No recursion — safe for deep/large graphs; single pass |

*(All three are asymptotically identical — the "optimal" label goes to the single-pass iterative version, which avoids both the second pass and stack-overflow risk.)*

---

## Approach 1 — Two-Pass Copy (Brute Force)

### Intuition
The core difficulty of cloning a graph is a chicken-and-egg problem: when copying node u's edge to v, the clone of v may not exist yet. The bluntest possible fix is to decouple the two concerns entirely. **Pass 1** traverses the graph and creates a bare clone (value only, no edges) for every node. **Pass 2** walks the original nodes again — now every clone is guaranteed to exist, so copying each edge is a mechanical map lookup.

### Algorithm
1. If the input node is `nil`, return `nil`.
2. **Pass 1 (discover + create):** BFS from the start node. When a node is first seen, store `clones[orig] = &Node{Val: orig.Val}` and enqueue it. Also record every discovered node in a list `order`.
3. **Pass 2 (wire edges):** for every original node in `order`, for each of its neighbors `nb`, append `clones[nb]` to `clones[orig].Neighbors` (preserving neighbor order).
4. Return `clones[start]`.

### Complexity
- **Time:** O(V + E) — pass 1 visits each node/edge once; pass 2 touches each edge once more (still linear).
- **Space:** O(V) — the clone map, the queue, and the discovered-node list.

### Code
```go
func twoPassCopy(node *Node) *Node {
	if node == nil {
		return nil // empty graph clones to an empty graph
	}

	clones := map[*Node]*Node{} // original node → its clone
	order := []*Node{}          // every discovered original node

	// Pass 1: BFS to discover all nodes and create edge-less clones.
	queue := []*Node{node}
	clones[node] = &Node{Val: node.Val} // clone the entry node up front
	for len(queue) > 0 {
		cur := queue[0]
		queue = queue[1:] // pop front
		order = append(order, cur)
		for _, nb := range cur.Neighbors {
			if _, seen := clones[nb]; !seen {
				clones[nb] = &Node{Val: nb.Val} // create clone on first sight
				queue = append(queue, nb)       // and schedule its exploration
			}
		}
	}

	// Pass 2: every clone exists now, so edges can be copied mechanically.
	for _, orig := range order {
		c := clones[orig]
		for _, nb := range orig.Neighbors {
			// point the clone at the *clone* of each neighbor, same order
			c.Neighbors = append(c.Neighbors, clones[nb])
		}
	}

	return clones[node]
}
```

### Dry Run — Example 1: `adjList = [[2,4],[1,3],[2,4],[1,3]]`

Pass 1 (node discovery; `n'` denotes the clone of node n):

| Step | Dequeued | Neighbors seen | New clones created | Queue after | `order` |
|------|----------|----------------|--------------------|-------------|---------|
| 0 | — (init) | — | `1'` | `[1]` | `[]` |
| 1 | 1 | 2, 4 | `2'`, `4'` | `[2, 4]` | `[1]` |
| 2 | 2 | 1 (seen), 3 | `3'` | `[4, 3]` | `[1, 2]` |
| 3 | 4 | 1 (seen), 3 (seen) | — | `[3]` | `[1, 2, 4]` |
| 4 | 3 | 2 (seen), 4 (seen) | — | `[]` | `[1, 2, 4, 3]` |

Pass 2 (edge wiring):

| Original | Edges copied to its clone |
|----------|---------------------------|
| 1 | `1'.Neighbors = [2', 4']` |
| 2 | `2'.Neighbors = [1', 3']` |
| 4 | `4'.Neighbors = [1', 3']` |
| 3 | `3'.Neighbors = [2', 4']` |

Serialized output: `[[2,4],[1,3],[2,4],[1,3]]` — identical structure, all-new pointers. ✅

---

## Approach 2 — DFS + Hash Map

### Intuition
To clone a node: copy its value, then clone each neighbor recursively and link to the results. The catch is cycles — in an undirected graph even a single edge is a 2-cycle (1→2→1→2…), so naive recursion never terminates. The fix is memoization: keep a map `original → clone`, and **register the clone before recursing into neighbors**. Any cycle that comes back to a node finds the clone already in the map and returns it immediately — the map is simultaneously the visited set, the cycle breaker, and the mechanism that makes two edges into the same node share one clone.

### Algorithm
1. If `node` is `nil`, return `nil`.
2. If `clones[node]` exists, return it (already cloned — visited check).
3. Create `c = &Node{Val: node.Val}` and set `clones[node] = c` **before** any recursion.
4. For each neighbor `nb` of `node`, append `dfs(nb)` to `c.Neighbors`.
5. Return `c`.

### Complexity
- **Time:** O(V + E) — each node is cloned exactly once; each undirected edge is traversed once from each side.
- **Space:** O(V) — the map plus the recursion stack, which can reach depth V on a path-shaped graph.

### Code
```go
func dfsClone(node *Node) *Node {
	clones := map[*Node]*Node{} // original → clone; doubles as visited set

	var dfs func(cur *Node) *Node
	dfs = func(cur *Node) *Node {
		if cur == nil {
			return nil // empty graph
		}
		if c, ok := clones[cur]; ok {
			return c // already cloned: reuse (breaks cycles)
		}
		c := &Node{Val: cur.Val}
		clones[cur] = c // register BEFORE recursing, or cycles loop forever
		for _, nb := range cur.Neighbors {
			// recursive call returns the (possibly shared) clone of nb
			c.Neighbors = append(c.Neighbors, dfs(nb))
		}
		return c
	}

	return dfs(node)
}
```

### Dry Run — Example 1: `adjList = [[2,4],[1,3],[2,4],[1,3]]`

| Step | Call | Map lookup | Action | `clones` map keys | Clone edges completed |
|------|------|------------|--------|-------------------|-----------------------|
| 1 | dfs(1) | miss | create `1'`, register, recurse into 2 | {1} | — |
| 2 | dfs(2) | miss | create `2'`, register, recurse into 1 | {1, 2} | — |
| 3 | dfs(1) | **hit** | return `1'` (cycle broken) | {1, 2} | `2' → [1']` partial |
| 4 | dfs(3) | miss | create `3'`, register, recurse into 2 | {1, 2, 3} | — |
| 5 | dfs(2) | **hit** | return `2'` | {1, 2, 3} | `3' → [2']` partial |
| 6 | dfs(4) | miss | create `4'`, register, recurse into 1 | {1, 2, 3, 4} | — |
| 7 | dfs(1) | **hit** | return `1'` | {1, 2, 3, 4} | `4' → [1']` partial |
| 8 | dfs(3) | **hit** | return `3'`; `4'` done → `[1', 3']` | {1, 2, 3, 4} | `4'` complete |
| 9 | unwind | — | `3' = [2', 4']`, `2' = [1', 3']` complete | {1, 2, 3, 4} | `3'`, `2'` complete |
| 10 | back in dfs(1) | **hit** on 4 | `1' = [2', 4']` complete; return `1'` | {1, 2, 3, 4} | all complete |

Serialized output: `[[2,4],[1,3],[2,4],[1,3]]`. ✅

---

## Approach 3 — BFS + Hash Map (Optimal, iterative)

### Intuition
Same map idea as DFS, but drive the traversal with an explicit queue instead of the call stack — immune to stack overflow on deep graphs and does everything in **one** pass (unlike Approach 1). When dequeuing a node, look at each neighbor: if it has no clone yet, create one and enqueue the neighbor; either way, the neighbor's clone now exists, so the cloned edge can be appended immediately.

### Algorithm
1. If `node` is `nil`, return `nil`.
2. Seed `clones = {node: &Node{Val: node.Val}}` and `queue = [node]`.
3. While the queue is non-empty:
   1. Dequeue `cur`.
   2. For each neighbor `nb` of `cur`:
      1. If `nb` has no clone, create `clones[nb]` and enqueue `nb`.
      2. Append `clones[nb]` to `clones[cur].Neighbors`.
4. Return `clones[node]`.

### Complexity
- **Time:** O(V + E) — every node enqueued exactly once, every edge processed once from each endpoint.
- **Space:** O(V) — the map plus the queue (bounded by the widest BFS level, ≤ V).

### Code
```go
func bfsClone(node *Node) *Node {
	if node == nil {
		return nil // empty graph
	}

	clones := map[*Node]*Node{node: {Val: node.Val}} // seed with entry clone
	queue := []*Node{node}                           // BFS frontier (originals)

	for len(queue) > 0 {
		cur := queue[0]
		queue = queue[1:] // pop front
		for _, nb := range cur.Neighbors {
			if _, seen := clones[nb]; !seen {
				clones[nb] = &Node{Val: nb.Val} // first visit: create clone
				queue = append(queue, nb)       // explore it later
			}
			// wire the cloned edge now; nb's clone certainly exists
			clones[cur].Neighbors = append(clones[cur].Neighbors, clones[nb])
		}
	}

	return clones[node]
}
```

### Dry Run — Example 1: `adjList = [[2,4],[1,3],[2,4],[1,3]]`

| Step | Dequeued | Neighbor | Clone exists? | Action | Queue after | Clone edges after |
|------|----------|----------|---------------|--------|-------------|-------------------|
| 0 | — (init) | — | — | seed `1'` | `[1]` | — |
| 1 | 1 | 2 | no | create `2'`, enqueue 2, wire `1'→2'` | `[2]` | `1'=[2']` |
| 2 | 1 | 4 | no | create `4'`, enqueue 4, wire `1'→4'` | `[2, 4]` | `1'=[2',4']` |
| 3 | 2 | 1 | **yes** | wire `2'→1'` only | `[4]` | `2'=[1']` |
| 4 | 2 | 3 | no | create `3'`, enqueue 3, wire `2'→3'` | `[4, 3]` | `2'=[1',3']` |
| 5 | 4 | 1 | **yes** | wire `4'→1'` | `[3]` | `4'=[1']` |
| 6 | 4 | 3 | **yes** | wire `4'→3'` | `[3]` | `4'=[1',3']` |
| 7 | 3 | 2 | **yes** | wire `3'→2'` | `[]` | `3'=[2']` |
| 8 | 3 | 4 | **yes** | wire `3'→4'` | `[]` | `3'=[2',4']` |

Queue empty → return `1'`. Serialized output: `[[2,4],[1,3],[2,4],[1,3]]`. ✅

---

## Key Takeaways

- **The `original → clone` hash map is the whole trick:** one structure serves as visited set, cycle breaker, and pointer-translation table. The same pattern deep-copies any pointer structure with cycles (LeetCode #138).
- **Register the clone *before* recursing/enqueuing neighbors** — inserting it after the recursive calls is the classic infinite-loop bug on cyclic graphs.
- **DFS vs BFS is a stack-vs-queue choice:** identical complexity; pick BFS iterative when the graph can be deep (recursion depth O(V)), DFS recursive for the shortest code.
- **Deep copy means zero shared pointers** — values equal, structure equal, but every `*Node` fresh. Verifying this (as `isDeepCopy` in `main.go` does) is a good habit; returning the input graph passes a naive value comparison.
- Edge cases worth stating in an interview: `nil` input (empty graph) and a single node with no neighbors.

---

## Related Problems

- LeetCode #138 — Copy List with Random Pointer (same clone-map pattern on a linked list)
- LeetCode #1485 — Clone Binary Tree With Random Pointer (same pattern on a tree)
- LeetCode #1490 — Clone N-ary Tree (simpler: no cycles, map optional)
- LeetCode #200 — Number of Islands (graph traversal with a visited set)
- LeetCode #547 — Number of Provinces (BFS/DFS over an adjacency structure)
