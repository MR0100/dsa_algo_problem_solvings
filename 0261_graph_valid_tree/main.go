package main

import "fmt"

// LeetCode 261 — Graph Valid Tree.
//
// Given n nodes labelled 0..n-1 and a list of undirected edges, decide whether
// the graph forms a VALID TREE. A graph is a tree iff it is:
//   - fully connected  (every node reachable from any other), AND
//   - acyclic          (contains no cycle).
//
// A key counting fact: a tree on n nodes has EXACTLY n-1 edges. So the quick
// necessary check `len(edges) == n-1` plus "connected" already implies acyclic
// (n-1 edges + connected ⇒ no cycle), and vice-versa.

// ── Approach 1: DFS (edge count + reachability) ──────────────────────────────
//
// dfsTree solves Graph Valid Tree by first checking the edge count, then doing
// a single DFS to confirm every node is reachable from node 0.
//
// Intuition:
//
//	A tree on n nodes has exactly n-1 edges. If we have n-1 edges and the graph
//	is connected, it cannot contain a cycle (a cycle would need an extra edge
//	somewhere else to still reach all nodes, exceeding n-1). So: verify
//	len(edges) == n-1, then verify one DFS from node 0 visits all n nodes.
//
// Algorithm:
//  1. If len(edges) != n-1, return false immediately (too few ⇒ disconnected,
//     too many ⇒ must contain a cycle).
//  2. Build an adjacency list.
//  3. DFS from node 0, marking visited nodes.
//  4. Return true iff every node was visited (fully connected).
//
// Time:  O(n + e) — build adjacency (O(e)) and DFS every node/edge once.
// Space: O(n + e) — adjacency list plus visited set and recursion stack.
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

// ── Approach 2: BFS (edge count + reachability) ──────────────────────────────
//
// bfsTree solves Graph Valid Tree with the same edge-count check but explores
// connectivity iteratively with a queue instead of recursion.
//
// Intuition:
//
//	Same reasoning as DFS: n-1 edges plus "all reachable from node 0" ⇒ tree.
//	BFS avoids deep recursion, which matters for long path-like graphs.
//
// Algorithm:
//  1. If len(edges) != n-1, return false.
//  2. Build adjacency list.
//  3. BFS from node 0 using a queue, counting distinct visited nodes.
//  4. Return true iff the visited count equals n.
//
// Time:  O(n + e) — each node enqueued once, each edge relaxed once.
// Space: O(n + e) — adjacency list plus the queue and visited set.
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

// ── Approach 3: Union-Find (Optimal) ─────────────────────────────────────────
//
// unionFind solves Graph Valid Tree by unioning the endpoints of every edge; a
// cycle is detected the moment an edge connects two nodes already in the same
// set, and connectivity falls out of the final component count.
//
// Intuition:
//
//	Process edges one by one, merging the two endpoints' sets. If both endpoints
//	are ALREADY in the same set, adding this edge creates a cycle ⇒ not a tree.
//	After processing all edges with no cycle, the graph is a forest; it is a
//	single tree iff exactly one component remains (equivalently n-1 unions
//	succeeded).
//
// Algorithm:
//  1. Initialise parent[i] = i (each node its own set) and count = n.
//  2. For each edge (u, v): find roots; if equal ⇒ cycle ⇒ return false;
//     else union them and decrement the component count.
//  3. Return true iff count == 1 (all nodes merged into one tree).
//
// Time:  O(n + e·α(n)) — near-linear; α is the inverse Ackermann function.
// Space: O(n) — the parent (and rank) arrays.
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

func main() {
	// Example 1: n = 5, edges = [[0,1],[0,2],[0,3],[1,4]] ⇒ true (a valid tree).
	n1 := 5
	edges1 := [][]int{{0, 1}, {0, 2}, {0, 3}, {1, 4}}

	// Example 2: n = 5, edges = [[0,1],[1,2],[2,3],[1,3],[1,4]] ⇒ false (cycle 1-2-3-1).
	n2 := 5
	edges2 := [][]int{{0, 1}, {1, 2}, {2, 3}, {1, 3}, {1, 4}}

	fmt.Println("=== Approach 1: DFS ===")
	fmt.Println(dfsTree(n1, edges1)) // expected true
	fmt.Println(dfsTree(n2, edges2)) // expected false

	fmt.Println("=== Approach 2: BFS ===")
	fmt.Println(bfsTree(n1, edges1)) // expected true
	fmt.Println(bfsTree(n2, edges2)) // expected false

	fmt.Println("=== Approach 3: Union-Find (Optimal) ===")
	fmt.Println(unionFind(n1, edges1)) // expected true
	fmt.Println(unionFind(n2, edges2)) // expected false
}
