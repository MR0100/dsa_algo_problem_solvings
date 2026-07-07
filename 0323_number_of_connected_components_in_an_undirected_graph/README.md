# 0323 — Number of Connected Components in an Undirected Graph

> LeetCode #323 · Difficulty: Medium (Premium)
> **Categories:** Union Find, Graph, DFS, BFS

---

## Problem Statement

You have a graph of `n` nodes. You are given an integer `n` and an array `edges`
where `edges[i] = [ai, bi]` indicates that there is an edge between `ai` and `bi`
in the graph.

Return the number of connected components in the graph.

**Example 1:**

```
Input:  n = 5, edges = [[0,1],[1,2],[3,4]]
Output: 2
```

Nodes `{0,1,2}` form one component and `{3,4}` form another.

**Example 2:**

```
Input:  n = 5, edges = [[0,1],[1,2],[2,3],[3,4]]
Output: 1
```

All five nodes are linked into a single chain, so there is one component.

**Constraints:**

- `1 <= n <= 2000`
- `1 <= edges.length <= 5000`
- `edges[i].length == 2`
- `0 <= ai <= bi < n`
- `ai != bi`
- There are no repeated edges.

---

## Company Frequency

| Company   | Frequency        | Last Reported |
|-----------|------------------|---------------|
| Google    | ★★★★☆ High       | 2024          |
| Amazon    | ★★★☆☆ Medium     | 2023          |
| Facebook  | ★★★☆☆ Medium     | 2023          |
| Twitter   | ★★☆☆☆ Low        | 2022          |

> ⚠️ Frequency data is crowd-sourced from LeetCode Discuss, Glassdoor, and
> community interview reports. Treat as a signal, not a guarantee.

---

## DSA Concepts Used

- **Union-Find (DSU)** — merge edge endpoints, decrement a component counter on
  each real union → see [`/dsa/union_find.md`](/dsa/union_find.md)
- **Graph DFS/BFS** — flood each component from an unvisited seed node → see
  [`/dsa/graph_bfs_dfs.md`](/dsa/graph_bfs_dfs.md)

---

## Approaches Overview

| # | Approach | Time | Space | When to use |
|---|----------|------|-------|-------------|
| 1 | Union-Find (Optimal) | O(n + e·α(n)) | O(n) | Dynamic merges / online connectivity |
| 2 | DFS | O(n + e) | O(n + e) | Simple static graph, recursion OK |
| 3 | BFS | O(n + e) | O(n + e) | Same, but avoids deep recursion |

(`e = len(edges)`, `α` = inverse Ackermann)

---

## Approach 1 — Union-Find (Disjoint Set Union) (Optimal)

### Intuition
Start with `n` components (each node alone). Every edge that connects two nodes
in *different* sets fuses them into one, dropping the count by one. Edges inside
an already-joined set change nothing. Union-by-rank + path compression make
`find` effectively O(1).

### Algorithm
1. `parent[i] = i`, `rank[i] = 0`, `count = n`.
2. For each edge `(a,b)`: if `find(a) != find(b)`, union them and `count--`.
3. Return `count`.

### Complexity
- **Time:** O(n + e·α(n)) — near-linear; α(n) is effectively constant.
- **Space:** O(n) — parent and rank arrays.

### Code
```go
func unionFind(n int, edges [][]int) int {
	parent := make([]int, n)
	rank := make([]int, n)
	for i := range parent {
		parent[i] = i
	}
	var find func(x int) int
	find = func(x int) int {
		for parent[x] != x {
			parent[x] = parent[parent[x]]
			x = parent[x]
		}
		return x
	}
	count := n
	for _, e := range edges {
		ra, rb := find(e[0]), find(e[1])
		if ra == rb {
			continue
		}
		if rank[ra] < rank[rb] {
			ra, rb = rb, ra
		}
		parent[rb] = ra
		if rank[ra] == rank[rb] {
			rank[ra]++
		}
		count--
	}
	return count
}
```

### Dry Run
Example 1: `n = 5`, `edges = [[0,1],[1,2],[3,4]]`. Start `count = 5`,
`parent = [0,1,2,3,4]`.

| edge  | find(a), find(b) | same? | action                 | count |
|-------|------------------|-------|------------------------|-------|
| (0,1) | 0, 1             | no    | parent[1]=0            | 4 |
| (1,2) | 0, 2             | no    | parent[2]=0            | 3 |
| (3,4) | 3, 4             | no    | parent[4]=3            | 2 |

Final `count = 2`. Sets: `{0,1,2}` and `{3,4}`. Output `2`.

---

## Approach 2 — DFS over Adjacency List

### Intuition
A component is a maximal set of mutually reachable nodes. Loop over all nodes;
the first unvisited node seeds a new component, and DFS floods everything
reachable so those nodes are not re-counted.

### Algorithm
1. Build an undirected adjacency list from `edges`.
2. `visited[]` all false; `count = 0`.
3. For each node `i`: if unvisited, `count++` and DFS-mark its whole component.
4. Return `count`.

### Complexity
- **Time:** O(n + e) — build the list once, visit each node/edge once.
- **Space:** O(n + e) — adjacency list plus visited/recursion.

### Code
```go
func dfsCount(n int, edges [][]int) int {
	adj := make([][]int, n)
	for _, e := range edges {
		adj[e[0]] = append(adj[e[0]], e[1])
		adj[e[1]] = append(adj[e[1]], e[0])
	}
	visited := make([]bool, n)
	var dfs func(u int)
	dfs = func(u int) {
		visited[u] = true
		for _, v := range adj[u] {
			if !visited[v] {
				dfs(v)
			}
		}
	}
	count := 0
	for i := 0; i < n; i++ {
		if !visited[i] {
			count++
			dfs(i)
		}
	}
	return count
}
```

### Dry Run
Example 1: `n = 5`, `edges = [[0,1],[1,2],[3,4]]`.
Adjacency: `0:[1] 1:[0,2] 2:[1] 3:[4] 4:[3]`.

| i | visited? | action              | visited set after |
|---|----------|---------------------|-------------------|
| 0 | no       | count=1, dfs(0)→1→2 | {0,1,2} |
| 1 | yes      | skip                | {0,1,2} |
| 2 | yes      | skip                | {0,1,2} |
| 3 | no       | count=2, dfs(3)→4   | {0,1,2,3,4} |
| 4 | yes      | skip                | all |

Output `2`.

---

## Approach 3 — BFS over Adjacency List

### Intuition
Same counting logic as DFS, but each component is drained with a queue instead of
recursion — useful when a long chain (up to 2000 nodes here) could stress the
recursion stack.

### Algorithm
1. Build adjacency list.
2. For each unvisited node: `count++`, seed a queue, BFS-mark reachable nodes.
3. Return `count`.

### Complexity
- **Time:** O(n + e) — every node and edge processed once.
- **Space:** O(n + e) — adjacency list plus the queue.

### Code
```go
func bfsCount(n int, edges [][]int) int {
	adj := make([][]int, n)
	for _, e := range edges {
		adj[e[0]] = append(adj[e[0]], e[1])
		adj[e[1]] = append(adj[e[1]], e[0])
	}
	visited := make([]bool, n)
	count := 0
	for i := 0; i < n; i++ {
		if visited[i] {
			continue
		}
		count++
		queue := []int{i}
		visited[i] = true
		for len(queue) > 0 {
			u := queue[0]
			queue = queue[1:]
			for _, v := range adj[u] {
				if !visited[v] {
					visited[v] = true
					queue = append(queue, v)
				}
			}
		}
	}
	return count
}
```

### Dry Run
Example 1: `n = 5`, `edges = [[0,1],[1,2],[3,4]]`.

| i | seed | queue evolution         | visited after   | count |
|---|------|-------------------------|-----------------|-------|
| 0 | 0    | [0]→[1]→[2]→[]           | {0,1,2}         | 1 |
| 1 | —    | visited, skip           | {0,1,2}         | 1 |
| 2 | —    | visited, skip           | {0,1,2}         | 1 |
| 3 | 3    | [3]→[4]→[]              | {0,1,2,3,4}     | 2 |
| 4 | —    | visited, skip           | all             | 2 |

Output `2`.

---

## Key Takeaways
- Counting connected components = **start at n and subtract one per real union**,
  or **launch one flood per unvisited seed**.
- Union-Find shines when edges arrive dynamically (online connectivity); DFS/BFS
  are simplest for a static, fully-given graph.
- Always mark BFS nodes visited **on enqueue**, not on dequeue, to avoid pushing
  the same node twice.

---

## Related Problems
- LeetCode #547 — Number of Provinces (same count, adjacency-matrix input)
- LeetCode #200 — Number of Islands (component count on a grid)
- LeetCode #261 — Graph Valid Tree (Union-Find: n-1 edges + connected)
- LeetCode #684 — Redundant Connection (Union-Find cycle detection)
