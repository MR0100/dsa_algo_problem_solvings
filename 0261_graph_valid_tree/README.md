# 0261 — Graph Valid Tree

> LeetCode #261 · Difficulty: Medium
> **Categories:** Graph, Union Find, Depth-First Search, Breadth-First Search

---

## Problem Statement

You have a graph of `n` nodes labeled from `0` to `n - 1`. You are given an integer `n` and a list of `edges` where `edges[i] = [aᵢ, bᵢ]` indicates that there is an undirected edge between nodes `aᵢ` and `bᵢ` in the graph.

Return `true` if the edges of the given graph make up a valid tree, and `false` otherwise.

**Example 1:**

```
Input: n = 5, edges = [[0,1],[0,2],[0,3],[1,4]]
Output: true
```

**Example 2:**

```
Input: n = 5, edges = [[0,1],[1,2],[2,3],[1,3],[1,4]]
Output: false
```

**Constraints:**

- `1 <= n <= 2000`
- `0 <= edges.length <= 5000`
- `edges[i].length == 2`
- `0 <= aᵢ, bᵢ < n`
- `aᵢ != bᵢ`
- There are no self-loops or repeated edges.

---

## Company Frequency

| Company   | Frequency        | Last Reported |
|-----------|------------------|---------------|
| Google    | ★★★★☆ High       | 2024          |
| Meta      | ★★★★☆ High       | 2024          |
| Amazon    | ★★★☆☆ Medium     | 2023          |
| Microsoft | ★★★☆☆ Medium     | 2023          |
| Bloomberg | ★★☆☆☆ Low        | 2022          |

> ⚠️ Frequency data is crowd-sourced from LeetCode Discuss, Glassdoor, and community interview reports. Treat as a signal, not a guarantee.

---

## DSA Concepts Used

- **Union-Find (Disjoint Set Union)** — merge edge endpoints; a repeated root means a cycle → see [`/dsa/union_find.md`](/dsa/union_find.md)
- **Graph BFS/DFS** — traverse from a source to test connectivity → see [`/dsa/graph_bfs_dfs.md`](/dsa/graph_bfs_dfs.md)
- **Tree edge-count invariant** — a tree on `n` nodes has exactly `n-1` edges → see [`/dsa/math_number_theory.md`](/dsa/math_number_theory.md)

---

## Approaches Overview

| # | Approach | Time | Space | When to use |
|---|----------|------|-------|-------------|
| 1 | DFS (edge count + reachability) | O(n + e) | O(n + e) | Recursive, concise; fine for moderate depth |
| 2 | BFS (edge count + reachability) | O(n + e) | O(n + e) | Avoids deep recursion on long chains |
| 3 | Union-Find (Optimal) | O(n + e·α(n)) | O(n) | Cycle detection while streaming edges |

---

## Approach 1 — DFS (edge count + reachability)

### Intuition
A tree on `n` nodes has exactly `n-1` edges. If the edge count is `n-1` **and** the graph is connected, it must also be acyclic — so we only need to confirm connectivity. One DFS from node `0` that reaches all `n` nodes proves the graph is a single connected component.

### Algorithm
1. If `len(edges) != n-1`, return `false` (too few → disconnected; too many → has a cycle).
2. Build an adjacency list from the undirected edges.
3. DFS from node `0`, marking every reached node as visited.
4. Return `true` iff every node was visited.

### Complexity
- **Time:** O(n + e) — build adjacency once, visit each node/edge once.
- **Space:** O(n + e) — adjacency list, visited array, recursion stack.

### Code
```go
func dfsTree(n int, edges [][]int) bool {
	if len(edges) != n-1 { // a tree has exactly n-1 edges; any other count fails
		return false
	}
	adj := make([][]int, n) // adjacency list: adj[u] = neighbours of u
	for _, e := range edges {
		u, v := e[0], e[1]
		adj[u] = append(adj[u], v) // undirected ⇒ store both directions
		adj[v] = append(adj[v], u)
	}
	visited := make([]bool, n) // visited[i] true once DFS reaches node i
	var dfs func(node int)
	dfs = func(node int) {
		visited[node] = true // mark current node reached
		for _, nb := range adj[node] {
			if !visited[nb] { // recurse only into unvisited neighbours
				dfs(nb)
			}
		}
	}
	dfs(0) // explore the single connected component containing node 0
	for _, seen := range visited {
		if !seen { // any unvisited node ⇒ graph is disconnected ⇒ not a tree
			return false
		}
	}
	return true // n-1 edges AND connected ⇒ valid tree
}
```

### Dry Run
Example 1: `n = 5`, `edges = [[0,1],[0,2],[0,3],[1,4]]`. Edge count `4 == 5-1` ✓.
Adjacency: `0→[1,2,3]`, `1→[0,4]`, `2→[0]`, `3→[0]`, `4→[1]`.

| Step | Call        | Action                        | visited            |
|------|-------------|-------------------------------|--------------------|
| 1    | dfs(0)      | mark 0; neighbours 1,2,3      | {0}                |
| 2    | dfs(1)      | mark 1; neighbours 0(seen),4  | {0,1}              |
| 3    | dfs(4)      | mark 4; neighbour 1(seen)     | {0,1,4}            |
| 4    | dfs(2)      | mark 2; neighbour 0(seen)     | {0,1,4,2}          |
| 5    | dfs(3)      | mark 3; neighbour 0(seen)     | {0,1,4,2,3}        |

All 5 nodes visited → return `true`.

---

## Approach 2 — BFS (edge count + reachability)

### Intuition
Identical reasoning to DFS — `n-1` edges plus "all nodes reachable from node 0" implies a valid tree — but connectivity is checked iteratively with a queue, avoiding recursion depth on long path-like inputs.

### Algorithm
1. If `len(edges) != n-1`, return `false`.
2. Build the adjacency list.
3. BFS from node `0` using a queue; mark and count each newly discovered node.
4. Return `true` iff the visited count equals `n`.

### Complexity
- **Time:** O(n + e) — each node enqueued once, each edge relaxed once.
- **Space:** O(n + e) — adjacency list plus queue and visited array.

### Code
```go
func bfsTree(n int, edges [][]int) bool {
	if len(edges) != n-1 { // tree edge-count invariant
		return false
	}
	adj := make([][]int, n)
	for _, e := range edges {
		u, v := e[0], e[1]
		adj[u] = append(adj[u], v)
		adj[v] = append(adj[v], u)
	}
	visited := make([]bool, n)
	queue := []int{0} // start BFS from node 0
	visited[0] = true
	count := 1 // number of distinct nodes reached so far
	for len(queue) > 0 {
		node := queue[0]  // pop front
		queue = queue[1:] // dequeue
		for _, nb := range adj[node] {
			if !visited[nb] { // first time we see this neighbour
				visited[nb] = true
				count++
				queue = append(queue, nb) // enqueue for later expansion
			}
		}
	}
	return count == n // reached every node ⇒ connected ⇒ valid tree
}
```

### Dry Run
Example 1: `n = 5`, edges as above. Edge count `4 == 4` ✓.

| Step | queue (front→) | node popped | new neighbours added | visited count |
|------|----------------|-------------|----------------------|---------------|
| init | [0]            | —           | —                    | 1             |
| 1    | [1,2,3]        | 0           | 1,2,3                | 4             |
| 2    | [2,3,4]        | 1           | 4                    | 5             |
| 3    | [3,4]          | 2           | —                    | 5             |
| 4    | [4]            | 3           | —                    | 5             |
| 5    | []             | 4           | —                    | 5             |

`count == 5 == n` → return `true`.

---

## Approach 3 — Union-Find (Optimal)

### Intuition
Merge the two endpoints of each edge into the same set. If an edge's endpoints are **already** in the same set, adding it closes a cycle → not a tree. If no cycle is ever formed, the graph is a forest; it is a single tree exactly when one component remains at the end.

### Algorithm
1. Initialise `parent[i] = i` for all nodes and `count = n` components.
2. For each edge `(u, v)`: find both roots. If equal → cycle → return `false`. Otherwise union them and decrement `count`.
3. Return `true` iff `count == 1` (everything merged into one tree).

### Complexity
- **Time:** O(n + e·α(n)) — near-linear with path compression + union by rank.
- **Space:** O(n) — parent and rank arrays.

### Code
```go
func unionFind(n int, edges [][]int) bool {
	parent := make([]int, n) // parent[i] = representative of i's set
	rank := make([]int, n)   // rank[i] ≈ tree height, for union by rank
	for i := range parent {
		parent[i] = i // every node starts as its own root
	}
	var find func(x int) int
	find = func(x int) int {
		for parent[x] != x { // walk up to the set's root
			parent[x] = parent[parent[x]] // path compression (halving)
			x = parent[x]
		}
		return x
	}
	count := n // number of disjoint components, starts at n singletons
	for _, e := range edges {
		ru, rv := find(e[0]), find(e[1])
		if ru == rv { // endpoints already connected ⇒ this edge closes a cycle
			return false
		}
		// union by rank: attach the shorter tree under the taller one
		if rank[ru] < rank[rv] {
			ru, rv = rv, ru
		}
		parent[rv] = ru
		if rank[ru] == rank[rv] {
			rank[ru]++
		}
		count-- // two components merged into one
	}
	return count == 1 // one component left ⇒ connected and acyclic ⇒ valid tree
}
```

### Dry Run
Example 2: `n = 5`, `edges = [[0,1],[1,2],[2,3],[1,3],[1,4]]`. Start `parent = [0,1,2,3,4]`, `count = 5`.

| Edge  | find(u) | find(v) | Same set? | Action                | count |
|-------|---------|---------|-----------|-----------------------|-------|
| (0,1) | 0       | 1       | no        | union → parent[1]=0   | 4     |
| (1,2) | 0       | 2       | no        | union → parent[2]=0   | 3     |
| (2,3) | 0       | 3       | no        | union → parent[3]=0   | 2     |
| (1,3) | 0       | 0       | **yes**   | cycle → return false  | —     |

Edge `(1,3)` finds both roots equal → returns `false`, correctly rejecting the cyclic graph.

---

## Key Takeaways

- **Tree = connected + acyclic + exactly `n-1` edges.** Any two of these three imply the third, so the edge-count shortcut (`len(edges) == n-1`) collapses the problem to a single connectivity check.
- **Union-Find detects cycles online:** two endpoints sharing a root before union means a cycle. This generalises to redundant-connection and MST (Kruskal) problems.
- Both DFS and BFS need an explicit visited array here because the graph is undirected (avoid walking back along the edge you came from — done implicitly by the visited check).

---

## Related Problems

- LeetCode #323 — Number of Connected Components in an Undirected Graph (union-find / DFS component counting)
- LeetCode #684 — Redundant Connection (union-find cycle detection)
- LeetCode #547 — Number of Provinces (connectivity via union-find)
- LeetCode #200 — Number of Islands (grid connectivity via DFS/BFS)
