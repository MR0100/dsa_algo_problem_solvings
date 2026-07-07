package main

import "fmt"

// Node is the LeetCode graph node: a value plus a list of neighbor pointers.
// The graph is undirected and connected; Val is unique and equals the
// 1-indexed position in the adjacency-list representation.
type Node struct {
	Val       int
	Neighbors []*Node
}

// ── Approach 1: Two-Pass Copy (Brute Force) ──────────────────────────────────
//
// twoPassCopy solves Clone Graph by first discovering every node and creating
// a bare clone for each, then wiring all neighbor lists in a second pass.
//
// Intuition: the difficulty of cloning a graph is that neighbors may not have
// been cloned yet when we need to point at them. The bluntest fix: separate
// the two concerns completely. Pass 1 traverses the graph and creates every
// clone node (values only, no edges). Pass 2 walks the original nodes again;
// now every clone already exists, so each edge can be copied directly.
//
// Algorithm:
//  1. BFS from the start node, collecting every original node and creating
//     clones[orig] = &Node{Val: orig.Val} with an empty neighbor list.
//  2. For every original node, for every neighbor, append the neighbor's
//     clone to the clone's neighbor list (preserving order).
//  3. Return clones[start].
//
// Time:  O(V + E) — each pass touches every node and edge once.
// Space: O(V) — the clone map, the visit queue, and the discovered-node list.
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

// ── Approach 2: DFS + Hash Map ───────────────────────────────────────────────
//
// dfsClone solves Clone Graph with a single recursive depth-first traversal,
// using a hash map both as the visited set and as the original→clone mapping.
//
// Intuition: to clone a node, clone its value, then clone each neighbor
// recursively and link to the results. Cycles (this is an undirected graph,
// so even one edge is a 2-cycle) would recurse forever — unless we memoize:
// the map returns the already-created clone the moment any node is revisited,
// which simultaneously breaks cycles and shares clones between edges.
//
// Algorithm:
//  1. If node is nil, return nil.
//  2. If node is already in the map, return its clone (visited check).
//  3. Create the clone, register it in the map BEFORE recursing (so cycles
//     back to this node find it), then recurse on every neighbor and append
//     each returned clone pointer.
//
// Time:  O(V + E) — every node cloned once, every edge walked once per side.
// Space: O(V) — the map plus recursion stack (O(V) deep for a path graph).
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

// ── Approach 3: BFS + Hash Map (Optimal, iterative) ──────────────────────────
//
// bfsClone solves Clone Graph with an iterative breadth-first traversal,
// cloning nodes level by level and wiring edges as they are crossed.
//
// Intuition: same map trick as DFS, but with an explicit queue instead of the
// call stack — immune to stack overflow on deep graphs. Each dequeued node's
// edges are processed once: an unseen neighbor gets a clone and a queue slot;
// either way the clone edge cur'→nb' is appended immediately.
//
// Algorithm:
//  1. If node is nil, return nil.
//  2. Create the entry clone and register it; enqueue the original entry.
//  3. While the queue is non-empty: pop cur; for each neighbor nb:
//     a. If nb has no clone yet, create it and enqueue nb.
//     b. Append nb's clone to cur's clone's neighbor list.
//  4. Return the entry clone.
//
// Time:  O(V + E) — each node enqueued once, each edge relaxed once per side.
// Space: O(V) — the map and the queue.
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

// ── Test helpers (build / serialize the LeetCode adjacency-list format) ──────

// buildGraph constructs a graph from LeetCode's adjacency list, where
// adj[i] lists the (1-indexed) values adjacent to node i+1. Returns node 1.
func buildGraph(adj [][]int) *Node {
	if len(adj) == 0 {
		return nil // empty adjacency list → empty graph
	}
	nodes := make([]*Node, len(adj))
	for i := range adj {
		nodes[i] = &Node{Val: i + 1} // create all nodes first
	}
	for i, nbs := range adj {
		for _, v := range nbs {
			nodes[i].Neighbors = append(nodes[i].Neighbors, nodes[v-1]) // wire edges
		}
	}
	return nodes[0]
}

// toAdjList serializes a graph back to the adjacency-list format by BFS,
// indexing nodes by Val (guaranteed unique, 1..n).
func toAdjList(node *Node) [][]int {
	if node == nil {
		return [][]int{} // empty graph serializes to []
	}
	byVal := map[int]*Node{} // Val → node, for stable ordered output
	visited := map[*Node]bool{node: true}
	queue := []*Node{node}
	maxVal := node.Val
	for len(queue) > 0 {
		cur := queue[0]
		queue = queue[1:]
		byVal[cur.Val] = cur
		if cur.Val > maxVal {
			maxVal = cur.Val // track highest value = node count
		}
		for _, nb := range cur.Neighbors {
			if !visited[nb] {
				visited[nb] = true
				queue = append(queue, nb)
			}
		}
	}
	adj := make([][]int, maxVal)
	for v := 1; v <= maxVal; v++ {
		adj[v-1] = []int{} // ensure [] rather than nil for isolated nodes
		for _, nb := range byVal[v].Neighbors {
			adj[v-1] = append(adj[v-1], nb.Val)
		}
	}
	return adj
}

// isDeepCopy verifies that no node pointer is shared between the two graphs.
func isDeepCopy(orig, clone *Node) bool {
	if orig == nil || clone == nil {
		return orig == nil && clone == nil // both empty is a valid deep copy
	}
	seen := map[*Node]bool{}
	var collect func(n *Node) // gather every original node pointer
	collect = func(n *Node) {
		if seen[n] {
			return
		}
		seen[n] = true
		for _, nb := range n.Neighbors {
			collect(nb)
		}
	}
	collect(orig)
	ok := true
	visited := map[*Node]bool{}
	var check func(n *Node) // no clone pointer may appear in the original set
	check = func(n *Node) {
		if visited[n] {
			return
		}
		visited[n] = true
		if seen[n] {
			ok = false // shared pointer → shallow copy, fail
			return
		}
		for _, nb := range n.Neighbors {
			check(nb)
		}
	}
	check(clone)
	return ok
}

func main() {
	// Example 1: 4-node cycle 1-2-3-4-1
	adj1 := [][]int{{2, 4}, {1, 3}, {2, 4}, {1, 3}}
	// Example 2: single node with no neighbors
	adj2 := [][]int{{}}
	// Example 3: empty graph
	adj3 := [][]int{}

	g1, g2, g3 := buildGraph(adj1), buildGraph(adj2), buildGraph(adj3)

	fmt.Println("=== Approach 1: Two-Pass Copy (Brute Force) ===")
	c1 := twoPassCopy(g1)
	fmt.Println(toAdjList(c1), "deep copy:", isDeepCopy(g1, c1)) // [[2 4] [1 3] [2 4] [1 3]] deep copy: true
	c2 := twoPassCopy(g2)
	fmt.Println(toAdjList(c2), "deep copy:", isDeepCopy(g2, c2)) // [[]] deep copy: true
	c3 := twoPassCopy(g3)
	fmt.Println(toAdjList(c3), "deep copy:", isDeepCopy(g3, c3)) // [] deep copy: true

	fmt.Println("=== Approach 2: DFS + Hash Map ===")
	d1 := dfsClone(g1)
	fmt.Println(toAdjList(d1), "deep copy:", isDeepCopy(g1, d1)) // [[2 4] [1 3] [2 4] [1 3]] deep copy: true
	d2 := dfsClone(g2)
	fmt.Println(toAdjList(d2), "deep copy:", isDeepCopy(g2, d2)) // [[]] deep copy: true
	d3 := dfsClone(g3)
	fmt.Println(toAdjList(d3), "deep copy:", isDeepCopy(g3, d3)) // [] deep copy: true

	fmt.Println("=== Approach 3: BFS + Hash Map (Optimal, iterative) ===")
	b1 := bfsClone(g1)
	fmt.Println(toAdjList(b1), "deep copy:", isDeepCopy(g1, b1)) // [[2 4] [1 3] [2 4] [1 3]] deep copy: true
	b2 := bfsClone(g2)
	fmt.Println(toAdjList(b2), "deep copy:", isDeepCopy(g2, b2)) // [[]] deep copy: true
	b3 := bfsClone(g3)
	fmt.Println(toAdjList(b3), "deep copy:", isDeepCopy(g3, b3)) // [] deep copy: true
}
