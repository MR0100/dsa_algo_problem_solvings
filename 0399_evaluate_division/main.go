package main

import "fmt"

// ── Approach 1: Graph DFS ────────────────────────────────────────────────────
//
// graphDFS solves Evaluate Division by modelling variables as graph nodes and
// each equation a/b = v as two weighted directed edges a->b (v) and b->a (1/v),
// then answering each query by DFS-multiplying edge weights along a path.
//
// Intuition:
//
//	a/b = 2 and b/c = 3 imply a/c = 6 by multiplying along the chain a->b->c.
//	Build a directed weighted graph; a query x/y is the product of edge weights
//	on any path from x to y. If x or y is unknown, or no path exists, answer -1.
//
// Algorithm:
//  1. Build adjacency: graph[a][b] = v, graph[b][a] = 1/v for each equation.
//  2. For each query (x, y):
//     a. If x or y not in graph → -1.
//     b. DFS from x toward y, carrying the accumulated product; return it on
//     arrival, using a visited set to avoid cycles.
//     c. If unreachable → -1.
//
// Time:  O(Q · (V + E)) — each query may traverse the whole graph.
// Space: O(V + E) — adjacency map plus recursion/visited per query.
func graphDFS(equations [][]string, values []float64, queries [][]string) []float64 {
	// graph[u][v] = weight of edge u->v (i.e. u/v).
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

	// dfs multiplies weights along a path from cur to target.
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

// ── Approach 2: Weighted Union-Find (Optimal) ────────────────────────────────
//
// weightedUnionFind solves Evaluate Division by grouping connected variables
// under a common root, storing each node's ratio to its parent so that a query
// reduces to comparing two nodes' ratios to their shared root.
//
// Intuition:
//
//	If every variable knows its value relative to the root of its group, then
//	for x and y in the same group, x/y = (x/root)/(y/root). Union-Find with a
//	weight[node] = node/parent lets us find(x) while compressing the path and
//	accumulating the ratio to the root; union merges two groups by wiring one
//	root under the other with the correct connecting weight.
//
// Algorithm:
//  1. For each equation a/b = v: union(a, b, v).
//  2. find(x) returns (root, ratio x/root), compressing along the way.
//  3. Query x/y: if unknown or different roots → -1; else ratioX / ratioY.
//
// Time:  ~O((E + Q)·α) with path compression — near-constant per op.
// Space: O(V) — parent and weight maps.
type dsu struct {
	parent map[string]string  // node -> parent
	weight map[string]float64 // node -> value(node)/value(parent)
}

func newDSU() *dsu {
	return &dsu{parent: map[string]string{}, weight: map[string]float64{}}
}

// add registers a fresh variable as its own root with ratio 1.
func (d *dsu) add(x string) {
	if _, ok := d.parent[x]; !ok {
		d.parent[x] = x
		d.weight[x] = 1.0
	}
}

// find returns the root of x and the ratio x/root, compressing the path.
func (d *dsu) find(x string) (string, float64) {
	if d.parent[x] == x {
		return x, 1.0
	}
	root, w := d.find(d.parent[x]) // recurse to get parent/root ratio
	d.parent[x] = root             // path compression: point x straight at root
	d.weight[x] *= w               // x/root = (x/parent)*(parent/root)
	return root, d.weight[x]
}

// union records a/b = value by merging their two groups.
func (d *dsu) union(a, b string, value float64) {
	d.add(a)
	d.add(b)
	ra, wa := d.find(a) // ra = root of a, wa = a/ra
	rb, wb := d.find(b) // rb = root of b, wb = b/rb
	if ra == rb {
		return // already related; assume input is consistent
	}
	// Attach ra under rb. We need weight[ra] = ra/rb.
	// a/b = value, a/ra = wa, b/rb = wb  =>  ra/rb = value*wb/wa.
	d.parent[ra] = rb
	d.weight[ra] = value * wb / wa
}

func weightedUnionFind(equations [][]string, values []float64, queries [][]string) []float64 {
	d := newDSU()
	for i, eq := range equations {
		d.union(eq[0], eq[1], values[i]) // build the disjoint sets with ratios
	}
	ans := make([]float64, len(queries))
	for i, q := range queries {
		x, y := q[0], q[1]
		_, okX := d.parent[x]
		_, okY := d.parent[y]
		if !okX || !okY {
			ans[i] = -1.0 // unknown variable
			continue
		}
		rx, wx := d.find(x) // x/root
		ry, wy := d.find(y) // y/root
		if rx != ry {
			ans[i] = -1.0 // different components → no relation
		} else {
			ans[i] = wx / wy // (x/root)/(y/root) = x/y
		}
	}
	return ans
}

func main() {
	// Example 1
	eq1 := [][]string{{"a", "b"}, {"b", "c"}}
	val1 := []float64{2.0, 3.0}
	q1 := [][]string{{"a", "c"}, {"b", "a"}, {"a", "e"}, {"a", "a"}, {"x", "x"}}
	// expected: [6.0, 0.5, -1.0, 1.0, -1.0]

	// Example 2
	eq2 := [][]string{{"a", "b"}, {"b", "c"}, {"bc", "cd"}}
	val2 := []float64{1.5, 2.5, 5.0}
	q2 := [][]string{{"a", "c"}, {"c", "b"}, {"bc", "cd"}, {"cd", "bc"}}
	// expected: [3.75, 0.4, 5.0, 0.2]

	// Example 3
	eq3 := [][]string{{"a", "b"}}
	val3 := []float64{0.5}
	q3 := [][]string{{"a", "b"}, {"b", "a"}, {"a", "c"}, {"x", "y"}}
	// expected: [0.5, 2.0, -1.0, -1.0]

	fmt.Println("=== Approach 1: Graph DFS ===")
	fmt.Printf("ex1 got=%v  expected [6 0.5 -1 1 -1]\n", graphDFS(eq1, val1, q1))
	fmt.Printf("ex2 got=%v  expected [3.75 0.4 5 0.2]\n", graphDFS(eq2, val2, q2))
	fmt.Printf("ex3 got=%v  expected [0.5 2 -1 -1]\n", graphDFS(eq3, val3, q3))

	fmt.Println("=== Approach 2: Weighted Union-Find (Optimal) ===")
	fmt.Printf("ex1 got=%v  expected [6 0.5 -1 1 -1]\n", weightedUnionFind(eq1, val1, q1))
	fmt.Printf("ex2 got=%v  expected [3.75 0.4 5 0.2]\n", weightedUnionFind(eq2, val2, q2))
	fmt.Printf("ex3 got=%v  expected [0.5 2 -1 -1]\n", weightedUnionFind(eq3, val3, q3))
}
