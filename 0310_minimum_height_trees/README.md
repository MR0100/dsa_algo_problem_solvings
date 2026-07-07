# 0310 — Minimum Height Trees

> LeetCode #310 · Difficulty: Medium
> **Categories:** Graph, BFS, Topological Sort, Tree

---

## Problem Statement

A tree is an undirected graph in which any two vertices are connected by *exactly* one path. In other words, any connected graph without simple cycles is a tree.

Given a tree of `n` nodes labelled from `0` to `n - 1`, and an array of `n - 1` `edges` where `edges[i] = [aᵢ, bᵢ]` indicates that there is an undirected edge between the two nodes `aᵢ` and `bᵢ` in the tree, you can choose any node of the tree as the root. When you select a node `x` as the root, the result tree has height `h`. Among all possible rooted trees, those with minimum height (i.e. `min(h)`) are called **minimum height trees** (MHTs).

Return *a list of all **MHTs'** root labels*. You can return the answer in **any order**.

The **height** of a rooted tree is the number of edges on the longest downward path between the root and a leaf.

**Example 1:**

```
Input: n = 4, edges = [[1,0],[1,2],[1,3]]
Output: [1]
Explanation: As shown, the height of the tree is 1 when the root is the node with label 1 which is the only MHT.
```

**Example 2:**

```
Input: n = 6, edges = [[3,0],[3,1],[3,2],[3,4],[5,4]]
Output: [3,4]
```

**Constraints:**

- `1 <= n <= 2 * 10^4`
- `edges.length == n - 1`
- `0 <= aᵢ, bᵢ < n`
- `aᵢ != bᵢ`
- All the pairs `(aᵢ, bᵢ)` are distinct.
- The given input is **guaranteed** to be a tree and there will be **no** repeated edges.

---

## Company Frequency

| Company   | Frequency        | Last Reported |
|-----------|------------------|---------------|
| Amazon    | ★★★☆☆ Medium     | 2024          |
| Google    | ★★★☆☆ Medium     | 2023          |
| Microsoft | ★★☆☆☆ Low        | 2023          |
| Uber      | ★★☆☆☆ Low        | 2022          |

> ⚠️ Frequency data is crowd-sourced from LeetCode Discuss, Glassdoor, and community interview reports. Treat as a signal, not a guarantee.

---

## DSA Concepts Used

- **Graph BFS** — measuring a rooting's height as the farthest BFS layer; brute-force eccentricity → see [`/dsa/graph_bfs_dfs.md`](/dsa/graph_bfs_dfs.md)
- **Topological Leaf Peeling** — Kahn-style removal of degree-1 nodes layer by layer until the centers remain → see [`/dsa/topological_sort.md`](/dsa/topological_sort.md)
- **Tree Properties** — a tree has at most two centers, which lie in the middle of its diameter → see [`/dsa/tree_traversal.md`](/dsa/tree_traversal.md)

---

## Approaches Overview

| # | Approach | Time | Space | When to use |
|---|----------|------|-------|-------------|
| 1 | Brute Force (BFS per node) | O(n²) | O(n) | Small n; conceptual clarity |
| 2 | Topological Leaf Peeling (Optimal) | O(n) | O(n) | Any n up to 2·10⁴ |

---

## Approach 1 — Brute Force (BFS from every node)

### Intuition

The height when rooted at node `v` equals `v`'s **eccentricity** — the distance to the farthest node. The MHT roots are the nodes with minimum eccentricity (the graph centers). Measure every node's eccentricity by BFS and keep the minima.

### Algorithm

1. Build an adjacency list.
2. For each node `r`: BFS from `r`, recording the maximum depth reached = its rooting height.
3. Track the minimum height; collect all `r` that achieve it (ties included).

### Complexity

- **Time:** O(n²) — `n` BFS traversals, each O(n) on a tree with `n-1` edges.
- **Space:** O(n) — adjacency list plus BFS queue and visited array.

### Code

```go
func bruteForce(n int, edges [][]int) []int {
	if n == 1 {
		return []int{0}
	}
	adj := buildAdj(n, edges)
	minHeight := n
	best := []int{}
	for r := 0; r < n; r++ {
		h := bfsHeight(r, adj, n)
		if h < minHeight {
			minHeight = h
			best = []int{r}
		} else if h == minHeight {
			best = append(best, r)
		}
	}
	return best
}

func bfsHeight(start int, adj [][]int, n int) int {
	visited := make([]bool, n)
	visited[start] = true
	queue := []int{start}
	height := 0
	for len(queue) > 0 {
		next := []int{}
		for _, u := range queue {
			for _, v := range adj[u] {
				if !visited[v] {
					visited[v] = true
					next = append(next, v)
				}
			}
		}
		if len(next) > 0 {
			height++
		}
		queue = next
	}
	return height
}
```

### Dry Run

`n = 4`, `edges = [[1,0],[1,2],[1,3]]` (a star centered at 1):

| Root r | BFS layers | Height |
|--------|------------|--------|
| 0 | {0} → {1} → {2,3} | 2 |
| 1 | {1} → {0,2,3} | 1 |
| 2 | {2} → {1} → {0,3} | 2 |
| 3 | {3} → {1} → {0,2} | 2 |

Minimum height is 1, achieved only by node 1 → answer `[1]`.

---

## Approach 2 — Topological Leaf Peeling (Optimal)

### Intuition

The MHT roots are the tree's **centroids**, sitting in the middle of its longest path (diameter). Peeling all leaves layer by layer shrinks the tree inward from both ends of every path at once. When `≤ 2` nodes remain, they are the middle of the diameter — the centers. A tree always has exactly 1 or 2 centers.

### Algorithm

1. If `n ≤ 2`, every node is a center — return `0..n-1`.
2. Build adjacency and a degree array; collect all degree-1 nodes as the first leaf layer.
3. While more than 2 nodes remain: peel the current leaves, decrement each neighbor's degree, and any neighbor dropping to degree 1 becomes a leaf for the next round.
4. The surviving 1 or 2 nodes are the answer.

### Complexity

- **Time:** O(n) — each node and edge is processed once.
- **Space:** O(n) — adjacency list, degree array, and leaf queues.

### Code

```go
func leafPeeling(n int, edges [][]int) []int {
	if n <= 2 {
		res := make([]int, n)
		for i := range res {
			res[i] = i
		}
		return res
	}
	adj := buildAdj(n, edges)
	degree := make([]int, n)
	for i := 0; i < n; i++ {
		degree[i] = len(adj[i])
	}
	leaves := []int{}
	for i := 0; i < n; i++ {
		if degree[i] == 1 {
			leaves = append(leaves, i)
		}
	}
	remaining := n
	for remaining > 2 {
		remaining -= len(leaves)
		next := []int{}
		for _, leaf := range leaves {
			for _, nb := range adj[leaf] {
				degree[nb]--
				if degree[nb] == 1 {
					next = append(next, nb)
				}
			}
		}
		leaves = next
	}
	return leaves
}
```

### Dry Run

`n = 4`, `edges = [[1,0],[1,2],[1,3]]`:

| Step | degree | leaves | remaining |
|------|--------|--------|-----------|
| init | [1,3,1,1] | [0,2,3] | 4 |
| peel [0,2,3] | node 1: 3→2→1→0 | next = [1] (when degree hits 1) | 4 − 3 = 1 |

`remaining = 1 ≤ 2`, loop stops. Survivor is `[1]` → answer `[1]`.

(For `n = 6` example, two rounds peel `[0,1,2,5]` then leave `remaining = 2`, giving centers `[3,4]`.)

---

## Key Takeaways

- **MHT roots = tree centers**, which are 1 or 2 nodes in the middle of the diameter — never more than two.
- **Leaf peeling (topological trimming)** is the linear-time centroid finder: strip degree-1 nodes layer by layer until `≤ 2` remain. It is Kahn's algorithm specialized to trees.
- **Handle `n ≤ 2` up front** — the loop assumes at least 3 nodes so it can shrink to a 1- or 2-node core.
- Think in terms of **eccentricity/centrality** when a problem asks for the best root of a tree; the brute force computes it directly, the peeling finds the minimum without measuring anyone.

---

## Related Problems

- LeetCode #207 — Course Schedule (Kahn's topological sort / leaf-like peeling by in-degree)
- LeetCode #543 — Diameter of Binary Tree (the diameter whose middle these centers are)
- LeetCode #863 — All Nodes Distance K in Binary Tree (BFS over a tree treated as a graph)
- LeetCode #1245 — Tree Diameter (find the longest path in a general tree)
