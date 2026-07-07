# 0399 — Evaluate Division

> LeetCode #399 · Difficulty: Medium
> **Categories:** Graph, DFS, BFS, Union-Find, Shortest Path

---

## Problem Statement

You are given an array of variable pairs `equations` and an array of real numbers `values`, where `equations[i] = [Ai, Bi]` and `values[i]` represent the equation `Ai / Bi = values[i]`. Each `Ai` or `Bi` is a string that represents a single variable.

You are also given some `queries`, where `queries[j] = [Cj, Dj]` represents the `j`th query where you must find the answer for `Cj / Dj = ?`.

Return *the answers to all queries*. If a single answer cannot be determined, return `-1.0`.

**Note:** The input is always valid. You may assume that evaluating the queries will not result in division by zero and that there is no contradiction.

**Note:** The variables that do not occur in the list of equations are undefined, so the answer cannot be determined for them.

**Example 1:**

```
Input: equations = [["a","b"],["b","c"]], values = [2.0,3.0], queries = [["a","c"],["b","a"],["a","e"],["a","a"],["x","x"]]
Output: [6.00000,0.50000,-1.00000,1.00000,-1.00000]
Explanation:
Given: a / b = 2.0, b / c = 3.0
queries are: a / c = ?, b / a = ?, a / e = ?, a / a = ?, x / x = ?
return: [6.0, 0.5, -1.0, 1.0, -1.0 ]
note: x is undefined => -1.0
```

**Example 2:**

```
Input: equations = [["a","b"],["b","c"],["bc","cd"]], values = [1.5,2.5,5.0], queries = [["a","c"],["c","b"],["bc","cd"],["cd","bc"]]
Output: [3.75000,0.40000,5.00000,0.20000]
```

**Example 3:**

```
Input: equations = [["a","b"]], values = [0.5], queries = [["a","b"],["b","a"],["a","c"],["x","y"]]
Output: [0.50000,2.00000,-1.00000,-1.00000]
```

**Constraints:**

- `1 <= equations.length <= 20`
- `equations[i].length == 2`
- `1 <= Ai.length, Bi.length <= 5`
- `values.length == equations.length`
- `0.0 < values[i] <= 20.0`
- `1 <= queries.length <= 20`
- `queries[i].length == 2`
- `1 <= Cj.length, Dj.length <= 5`
- `Ai, Bi, Cj, Dj` consist of lower case English letters and digits.

---

## Company Frequency

| Company    | Frequency        | Last Reported |
|------------|------------------|---------------|
| Google     | ★★★★★ Very High  | 2024          |
| Facebook   | ★★★★☆ High       | 2024          |
| Amazon     | ★★★★☆ High       | 2024          |
| Microsoft  | ★★★☆☆ Medium     | 2023          |
| Uber       | ★★★☆☆ Medium     | 2023          |
| Bloomberg  | ★★☆☆☆ Low        | 2022          |

> ⚠️ Frequency data is crowd-sourced from LeetCode Discuss, Glassdoor, and
> community interview reports. Treat as a signal, not a guarantee.

---

## DSA Concepts Used

- **Graph modelling + DFS** — each equation `a/b = v` is a pair of weighted directed edges; a query is the product of edge weights along a path, found by depth-first search → see [`/dsa/graph_bfs_dfs.md`](/dsa/graph_bfs_dfs.md)
- **Weighted Union-Find** — group variables by connected component and store each node's ratio to its root, so a query becomes a division of two ratios in near-constant time → see [`/dsa/union_find.md`](/dsa/union_find.md)
- **Hash Map** — variables are strings, so adjacency lists / parent pointers are keyed by string in maps → see [`/dsa/hash_map.md`](/dsa/hash_map.md)

---

## Approaches Overview

| # | Approach | Time | Space | When to use |
|---|----------|------|-------|-------------|
| 1 | Graph DFS | O(Q·(V+E)) | O(V+E) | Intuitive; fine for the small constraints (≤20 edges) |
| 2 | Weighted Union-Find (Optimal) | ~O((E+Q)·α) | O(V) | Many queries / large graphs; amortised near-constant per op |

---

## Approach 1 — Graph DFS

### Intuition

Treat variables as nodes. `a/b = 2` means an edge `a → b` with weight `2` and, because division inverts, an edge `b → a` with weight `1/2`. Chaining edges multiplies weights: `a/b = 2` and `b/c = 3` give `a/c = 2·3 = 6` along the path `a → b → c`. So a query `x/y` is the product of weights on any path from `x` to `y`; DFS finds one. Unknown variable or no path ⇒ `-1`.

### Algorithm

1. Build adjacency: for each equation set `graph[a][b] = v` and `graph[b][a] = 1/v`.
2. For each query `(x, y)`:
   1. If `x` or `y` is absent from the graph → `-1`.
   2. DFS from `x` toward `y`, carrying the accumulated product; return it on arrival. Use a `seen` set to avoid cycles.
   3. If `y` is unreachable → `-1`.

### Complexity

- **Time:** O(Q·(V+E)) — each query can traverse the whole graph once.
- **Space:** O(V+E) — adjacency map, plus O(V) recursion/visited per query.

### Code

```go
func graphDFS(equations [][]string, values []float64, queries [][]string) []float64 {
	graph := map[string]map[string]float64{}
	addEdge := func(u, v string, w float64) {
		if graph[u] == nil {
			graph[u] = map[string]float64{}
		}
		graph[u][v] = w
	}
	for i, eq := range equations {
		a, b := eq[0], eq[1]
		addEdge(a, b, values[i])     // a/b = v
		addEdge(b, a, 1.0/values[i]) // b/a = 1/v
	}

	var dfs func(cur, target string, acc float64, seen map[string]bool) (float64, bool)
	dfs = func(cur, target string, acc float64, seen map[string]bool) (float64, bool) {
		if cur == target {
			return acc, true // arrived; acc is the product cur/... /target
		}
		seen[cur] = true // avoid revisiting (cycles)
		for next, w := range graph[cur] {
			if seen[next] {
				continue
			}
			if val, ok := dfs(next, target, acc*w, seen); ok {
				return val, true // propagate the first successful path
			}
		}
		return 0, false
	}

	ans := make([]float64, len(queries))
	for i, q := range queries {
		x, y := q[0], q[1]
		if graph[x] == nil || graph[y] == nil {
			ans[i] = -1.0 // unknown variable → undefined
			continue
		}
		if val, ok := dfs(x, y, 1.0, map[string]bool{}); ok {
			ans[i] = val
		} else {
			ans[i] = -1.0 // no connecting path
		}
	}
	return ans
}
```

### Dry Run

Example 1: `a/b = 2`, `b/c = 3`. Graph edges: `a→b=2, b→a=0.5, b→c=3, c→b=1/3`.

| Query | DFS trace | result |
|-------|-----------|--------|
| a/c | a→b (acc 2) → c (acc 2·3=6) | **6** |
| b/a | b→a (acc 0.5) | **0.5** |
| a/e | `e` not in graph | **-1** |
| a/a | cur==target immediately, acc=1 | **1** |
| x/x | `x` not in graph | **-1** |

Result: `[6, 0.5, -1, 1, -1]` ✔

---

## Approach 2 — Weighted Union-Find (Optimal)

### Intuition

If every variable knows its value **relative to the root** of its group, then for `x` and `y` in the same group, `x/y = (x/root)/(y/root)`. Maintain Disjoint-Set-Union where `weight[node] = value(node)/value(parent)`. `find(x)` walks to the root, compressing the path and multiplying weights to yield `x/root`. `union(a,b,v)` merges two groups by attaching one root under the other with the ratio that makes `a/b = v` hold.

To attach root `ra` under root `rb`: from `a/ra = wa`, `b/rb = wb`, `a/b = v`, we get `ra/rb = v·wb/wa`, which becomes `weight[ra]`.

### Algorithm

1. For each equation `a/b = v`: `union(a, b, v)` (creating fresh singletons as needed).
2. `find(x)` returns `(root, x/root)`, compressing the path along the way.
3. Query `x/y`: if either is unknown or roots differ → `-1`; else return `(x/root)/(y/root)`.

### Complexity

- **Time:** ~O((E + Q)·α(V)) with path compression — near-constant amortised per operation.
- **Space:** O(V) — a `parent` map and a `weight` map.

### Code

```go
type dsu struct {
	parent map[string]string  // node -> parent
	weight map[string]float64 // node -> value(node)/value(parent)
}

func newDSU() *dsu {
	return &dsu{parent: map[string]string{}, weight: map[string]float64{}}
}

func (d *dsu) add(x string) {
	if _, ok := d.parent[x]; !ok {
		d.parent[x] = x
		d.weight[x] = 1.0
	}
}

func (d *dsu) find(x string) (string, float64) {
	if d.parent[x] == x {
		return x, 1.0
	}
	root, w := d.find(d.parent[x]) // recurse to get parent/root ratio
	d.parent[x] = root             // path compression: point x straight at root
	d.weight[x] *= w               // x/root = (x/parent)*(parent/root)
	return root, d.weight[x]
}

func (d *dsu) union(a, b string, value float64) {
	d.add(a)
	d.add(b)
	ra, wa := d.find(a) // ra = root of a, wa = a/ra
	rb, wb := d.find(b) // rb = root of b, wb = b/rb
	if ra == rb {
		return // already related; assume input is consistent
	}
	// a/b = value, a/ra = wa, b/rb = wb  =>  ra/rb = value*wb/wa.
	d.parent[ra] = rb
	d.weight[ra] = value * wb / wa
}

func weightedUnionFind(equations [][]string, values []float64, queries [][]string) []float64 {
	d := newDSU()
	for i, eq := range equations {
		d.union(eq[0], eq[1], values[i])
	}
	ans := make([]float64, len(queries))
	for i, q := range queries {
		x, y := q[0], q[1]
		_, okX := d.parent[x]
		_, okY := d.parent[y]
		if !okX || !okY {
			ans[i] = -1.0
			continue
		}
		rx, wx := d.find(x) // x/root
		ry, wy := d.find(y) // y/root
		if rx != ry {
			ans[i] = -1.0
		} else {
			ans[i] = wx / wy // (x/root)/(y/root) = x/y
		}
	}
	return ans
}
```

### Dry Run

Example 1: equations `a/b = 2`, `b/c = 3`.

| Step | action | resulting state |
|------|--------|-----------------|
| union(a,b,2) | add a,b; ra=a,rb=b; attach a under b, weight[a]=2·1/1=2 | parent{a:b,b:b}, weight{a:2,b:1} |
| union(b,c,3) | add c; find(b)=(b,1), find(c)=(c,1); attach b under c, weight[b]=3 | parent{a:b,b:c,c:c}, weight{a:2,b:3,c:1} |

Queries (each `find` compresses and accumulates):

| Query | find(x) → (root, ratio) | find(y) → (root, ratio) | answer |
|-------|-------------------------|-------------------------|--------|
| a/c | a: a→b→c ⇒ (c, 2·3=6) | c: (c, 1) | 6/1 = **6** |
| b/a | b: (c, 3) | a: (c, 6) | 3/6 = **0.5** |
| a/e | e unknown | — | **-1** |
| a/a | (c, 6) | (c, 6) | 6/6 = **1** |
| x/x | x unknown | — | **-1** |

Result: `[6, 0.5, -1, 1, -1]` ✔

---

## Key Takeaways

- **Ratios chain multiplicatively** — modelling `a/b = v` as a weighted graph turns "evaluate a/c" into "multiply weights along a path", the core insight for both approaches.
- **Add both directions** (`a→b = v` and `b→a = 1/v`) so queries can be answered in either orientation.
- **Weighted Union-Find** stores each node's ratio to its root; `find` compresses paths *and* folds the accumulated ratio, giving near-constant queries — the scalable answer when queries or the graph are large.
- Undefined variables (never seen in any equation) must return `-1`, and so must cross-component queries. Check membership and root-equality explicitly.

---

## Related Problems

- LeetCode #547 — Number of Provinces (Union-Find connectivity)
- LeetCode #684 — Redundant Connection (Union-Find cycle detection)
- LeetCode #990 — Satisfiability of Equality Equations (Union-Find over `==`/`!=`)
- LeetCode #785 — Is Graph Bipartite (graph traversal with per-node labels)
- LeetCode #1631 — Path With Minimum Effort (weighted-graph pathfinding)
