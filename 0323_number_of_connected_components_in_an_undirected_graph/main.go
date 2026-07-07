package main

import "fmt"

// ── Approach 1: Union-Find (Disjoint Set Union) (Optimal) ────────────────────
//
// unionFind solves Number of Connected Components by starting with n separate
// nodes and merging the two endpoints of every edge. Each successful merge joins
// two components into one, so the component count drops by one.
//
// Intuition:
//
//	Begin with n components (each node alone). Every edge that links two nodes
//	from *different* sets fuses those sets. Track a running count that starts at
//	n and decrements on each real merge; edges within an already-joined set do
//	nothing. Union by rank + path compression keep find() near O(1).
//
// Algorithm:
//  1. parent[i] = i, rank[i] = 0, count = n.
//  2. For each edge (a,b): if find(a) != find(b), union them and count--.
//  3. Return count.
//
// Time:  O(n + e·α(n)) — α is the near-constant inverse Ackermann.
// Space: O(n) — parent and rank arrays.
func unionFind(n int, edges [][]int) int {
	parent := make([]int, n) // parent[i] = representative of i's set
	rank := make([]int, n)   // rank[i] = tree-height upper bound for union-by-rank
	for i := range parent {
		parent[i] = i // each node is initially its own root
	}
	// find returns the root of x with path compression.
	var find func(x int) int
	find = func(x int) int {
		for parent[x] != x { // walk up to the root
			parent[x] = parent[parent[x]] // path halving: point to grandparent
			x = parent[x]
		}
		return x
	}
	count := n // start with n isolated components
	for _, e := range edges {
		ra, rb := find(e[0]), find(e[1]) // roots of the two endpoints
		if ra == rb {
			continue // already in the same component: no merge
		}
		// Union by rank: attach the shorter tree under the taller one.
		if rank[ra] < rank[rb] {
			ra, rb = rb, ra // ensure ra is the taller (or equal) root
		}
		parent[rb] = ra // hang rb's tree under ra
		if rank[ra] == rank[rb] {
			rank[ra]++ // equal ranks: resulting tree grows by one
		}
		count-- // two components became one
	}
	return count
}

// ── Approach 2: DFS over Adjacency List ──────────────────────────────────────
//
// dfsCount solves Number of Connected Components by building an adjacency list
// and launching a DFS from every not-yet-visited node; each launch discovers one
// whole component.
//
// Intuition:
//
//	A component is a maximal set of nodes reachable from any of its members. Loop
//	over all nodes; the first time we meet an unvisited node it seeds a new
//	component, and DFS floods everything reachable from it so those nodes are not
//	counted again.
//
// Algorithm:
//  1. Build adjacency list from edges (undirected → add both directions).
//  2. visited[] all false; count = 0.
//  3. For each node i: if unvisited, count++ and DFS-mark its whole component.
//  4. Return count.
//
// Time:  O(n + e) — build the list once, then visit every node/edge once.
// Space: O(n + e) — adjacency list plus visited/recursion.
func dfsCount(n int, edges [][]int) int {
	adj := make([][]int, n) // adj[i] = neighbours of node i
	for _, e := range edges {
		adj[e[0]] = append(adj[e[0]], e[1]) // undirected: record both ends
		adj[e[1]] = append(adj[e[1]], e[0])
	}
	visited := make([]bool, n) // visited[i] = node i already explored
	var dfs func(u int)
	dfs = func(u int) {
		visited[u] = true // mark on entry
		for _, v := range adj[u] {
			if !visited[v] {
				dfs(v) // flood into unvisited neighbours
			}
		}
	}
	count := 0
	for i := 0; i < n; i++ {
		if !visited[i] { // a fresh, undiscovered node
			count++ // seeds a new component
			dfs(i)  // consume the entire component
		}
	}
	return count
}

// ── Approach 3: BFS over Adjacency List ──────────────────────────────────────
//
// bfsCount solves Number of Connected Components the same way as DFS but explores
// each component with a queue instead of recursion.
//
// Intuition:
//
//	Identical counting logic: each unvisited node starts a new component; BFS
//	then drains everything reachable from it level by level. Handy when recursion
//	depth could blow the stack on long chains.
//
// Algorithm:
//  1. Build adjacency list.
//  2. For each unvisited node: count++, seed a queue, and BFS-mark reachable nodes.
//  3. Return count.
//
// Time:  O(n + e) — every node and edge processed once.
// Space: O(n + e) — adjacency list plus the queue.
func bfsCount(n int, edges [][]int) int {
	adj := make([][]int, n) // adjacency list
	for _, e := range edges {
		adj[e[0]] = append(adj[e[0]], e[1])
		adj[e[1]] = append(adj[e[1]], e[0])
	}
	visited := make([]bool, n)
	count := 0
	for i := 0; i < n; i++ {
		if visited[i] {
			continue // already part of a counted component
		}
		count++           // new component discovered
		queue := []int{i} // BFS frontier
		visited[i] = true
		for len(queue) > 0 {
			u := queue[0] // dequeue front
			queue = queue[1:]
			for _, v := range adj[u] {
				if !visited[v] {
					visited[v] = true // mark before enqueue (avoids dupes)
					queue = append(queue, v)
				}
			}
		}
	}
	return count
}

func main() {
	fmt.Println("=== Approach 1: Union-Find (Optimal) ===")
	fmt.Println(unionFind(5, [][]int{{0, 1}, {1, 2}, {3, 4}}))         // expected 2
	fmt.Println(unionFind(5, [][]int{{0, 1}, {1, 2}, {2, 3}, {3, 4}})) // expected 1

	fmt.Println("=== Approach 2: DFS ===")
	fmt.Println(dfsCount(5, [][]int{{0, 1}, {1, 2}, {3, 4}}))         // expected 2
	fmt.Println(dfsCount(5, [][]int{{0, 1}, {1, 2}, {2, 3}, {3, 4}})) // expected 1

	fmt.Println("=== Approach 3: BFS ===")
	fmt.Println(bfsCount(5, [][]int{{0, 1}, {1, 2}, {3, 4}}))         // expected 2
	fmt.Println(bfsCount(5, [][]int{{0, 1}, {1, 2}, {2, 3}, {3, 4}})) // expected 1
}
