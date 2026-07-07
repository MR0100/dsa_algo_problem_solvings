package main

import "fmt"

// ── Approach 1: Brute Force (BFS height from every node) ──────────────────────
//
// bruteForce solves Minimum Height Trees by rooting the tree at each node in
// turn, computing that rooting's height with a BFS, and keeping the roots that
// achieve the global minimum height.
//
// Intuition:
//
//	The height when rooted at node v is the eccentricity of v — the distance to
//	the farthest node. The answer is exactly the set of nodes with minimum
//	eccentricity (the graph "centers"). The most direct way to find them is to
//	measure every node's eccentricity by BFS and take the minima.
//
// Algorithm:
//
//  1. Build an adjacency list.
//  2. For each node r: BFS from r, record the maximum depth reached = height.
//  3. Track the minimum height seen; collect all r achieving it.
//
// Time:  O(n^2) — n BFS traversals, each O(n) on a tree (n-1 edges).
// Space: O(n) — adjacency list + BFS queue/visited.
func bruteForce(n int, edges [][]int) []int {
	if n == 1 {
		return []int{0} // single node: height 0, it is the only center
	}
	adj := buildAdj(n, edges)
	minHeight := n // upper bound
	best := []int{}
	for r := 0; r < n; r++ {
		h := bfsHeight(r, adj, n) // eccentricity of r
		if h < minHeight {
			minHeight = h
			best = []int{r} // new best → reset the list
		} else if h == minHeight {
			best = append(best, r) // ties join the answer set
		}
	}
	return best
}

// bfsHeight returns the height of the tree when rooted at start (max depth).
func bfsHeight(start int, adj [][]int, n int) int {
	visited := make([]bool, n)
	visited[start] = true
	queue := []int{start}
	height := 0
	for len(queue) > 0 {
		next := []int{} // nodes of the next BFS layer
		for _, u := range queue {
			for _, v := range adj[u] {
				if !visited[v] {
					visited[v] = true
					next = append(next, v)
				}
			}
		}
		if len(next) > 0 {
			height++ // completed one more level below the root
		}
		queue = next
	}
	return height
}

// ── Approach 2: Topological Peeling of Leaves (Optimal) ──────────────────────
//
// leafPeeling solves Minimum Height Trees by repeatedly stripping all current
// leaves until 1 or 2 nodes remain — those remaining nodes are the centers.
//
// Intuition:
//
//	The minimum-height roots are the centroids of the tree, found on the longest
//	path (diameter). Peeling leaves layer by layer shrinks the tree inward from
//	both ends of every path simultaneously. When ≤ 2 nodes are left, they are the
//	middle of the diameter — the centers. A tree always has 1 or 2 centers.
//
// Algorithm:
//
//  1. Handle n ≤ 2 directly (all nodes are centers).
//  2. Build adjacency and a degree count; collect degree-1 nodes as leaves.
//  3. While more than 2 nodes remain: remove the current leaf layer, decrement
//     neighbors' degrees, and any neighbor dropping to degree 1 becomes a new
//     leaf for the next round.
//  4. The surviving 1 or 2 nodes are the answer.
//
// Time:  O(n) — every node and edge is processed once.
// Space: O(n) — adjacency list, degree array, leaf queues.
func leafPeeling(n int, edges [][]int) []int {
	if n <= 2 {
		res := make([]int, n) // 0..n-1 are all centers when n is 1 or 2
		for i := range res {
			res[i] = i
		}
		return res
	}
	adj := buildAdj(n, edges)
	degree := make([]int, n)
	for i := 0; i < n; i++ {
		degree[i] = len(adj[i]) // initial degree of each node
	}
	// Seed the first layer with all leaves (degree 1).
	leaves := []int{}
	for i := 0; i < n; i++ {
		if degree[i] == 1 {
			leaves = append(leaves, i)
		}
	}
	remaining := n
	for remaining > 2 {
		remaining -= len(leaves) // these leaves get peeled off
		next := []int{}          // next layer of leaves
		for _, leaf := range leaves {
			for _, nb := range adj[leaf] {
				degree[nb]-- // this edge is gone
				if degree[nb] == 1 {
					next = append(next, nb) // neighbor is now a leaf
				}
			}
		}
		leaves = next // advance one layer inward
	}
	return leaves // the 1 or 2 survivors are the centers
}

// buildAdj constructs an undirected adjacency list from the edge list.
func buildAdj(n int, edges [][]int) [][]int {
	adj := make([][]int, n)
	for _, e := range edges {
		u, v := e[0], e[1]
		adj[u] = append(adj[u], v) // undirected → add both directions
		adj[v] = append(adj[v], u)
	}
	return adj
}

func main() {
	fmt.Println("=== Approach 1: Brute Force (BFS from every node) ===")
	fmt.Println(bruteForce(4, [][]int{{1, 0}, {1, 2}, {1, 3}}))                 // expected [1]
	fmt.Println(bruteForce(6, [][]int{{3, 0}, {3, 1}, {3, 2}, {3, 4}, {5, 4}})) // expected [3 4]
	fmt.Println(bruteForce(1, [][]int{}))                                       // expected [0]
	fmt.Println(bruteForce(2, [][]int{{0, 1}}))                                 // expected [0 1]

	fmt.Println("=== Approach 2: Topological Leaf Peeling (Optimal) ===")
	fmt.Println(leafPeeling(4, [][]int{{1, 0}, {1, 2}, {1, 3}}))                 // expected [1]
	fmt.Println(leafPeeling(6, [][]int{{3, 0}, {3, 1}, {3, 2}, {3, 4}, {5, 4}})) // expected [3 4]
	fmt.Println(leafPeeling(1, [][]int{}))                                       // expected [0]
	fmt.Println(leafPeeling(2, [][]int{{0, 1}}))                                 // expected [0 1]
}
