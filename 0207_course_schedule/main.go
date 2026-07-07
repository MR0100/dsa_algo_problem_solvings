package main

import "fmt"

// ── Approach 1: Brute Force (Repeated Elimination) ───────────────────────────
//
// bruteForce solves Course Schedule by repeatedly sweeping all courses and
// "taking" any course whose prerequisites have all been taken, until either
// every course is taken or a full sweep takes nothing.
//
// Intuition:
//
//	Simulate a student with no planning skills: each semester, scan the whole
//	catalogue and enrol in every course whose prerequisites are already done.
//	If all courses eventually get taken, the schedule is feasible. If one full
//	scan takes nothing while courses remain, those remaining courses all wait
//	on each other — a prerequisite cycle — so finishing is impossible.
//
// Algorithm:
//  1. Keep a taken[i] boolean per course and a count of courses taken.
//  2. Loop: sweep every course i not yet taken; if every prerequisite pair
//     [i, b] has taken[b] == true, mark taken[i] and note progress.
//  3. Stop sweeping when a pass makes no progress.
//  4. Return true iff the number taken equals numCourses.
//
// Time:  O(V · (V + E)) — up to V sweeps (each sweep takes ≥1 course or
//
//	stops), each sweep scans V courses and all E prerequisite pairs.
//
// Space: O(V) — the taken array.
func bruteForce(numCourses int, prerequisites [][]int) bool {
	taken := make([]bool, numCourses) // taken[i] = course i already completed
	count := 0                        // how many courses completed so far

	for {
		progress := false // did this sweep manage to take any course?
		for c := 0; c < numCourses; c++ {
			if taken[c] {
				continue // already completed — skip
			}
			// Check every prerequisite pair that constrains course c.
			ready := true
			for _, p := range prerequisites {
				// p = [a, b] means "to take a you must first take b".
				if p[0] == c && !taken[p[1]] {
					ready = false // some prerequisite of c is still pending
					break
				}
			}
			if ready {
				taken[c] = true // enrol: all prerequisites satisfied
				count++
				progress = true // this sweep achieved something
			}
		}
		if !progress {
			break // stuck: nothing new can be taken — stop simulating
		}
	}
	// Feasible iff the simulation completed every course.
	return count == numCourses
}

// ── Approach 2: DFS Cycle Detection (3-Colour) ───────────────────────────────
//
// dfsCycleDetection solves Course Schedule by depth-first searching the
// prerequisite graph and reporting whether any back edge (cycle) exists.
//
// Intuition:
//
//	Model courses as graph nodes with an edge b → a for each pair [a, b]
//	("b unlocks a"). All courses are finishable iff this directed graph has
//	no cycle. DFS finds cycles with three node states: white (unvisited),
//	grey (in the current recursion path), black (fully explored). Meeting a
//	grey node again means the path looped back onto itself — a cycle.
//
// Algorithm:
//  1. Build adjacency list adj[b] = all courses that list b as prerequisite.
//  2. state[i] ∈ {0 white, 1 grey, 2 black}, all initially white.
//  3. DFS(u): mark u grey; for each neighbour v: grey v → cycle; white v →
//     recurse (propagate cycle upwards). Afterwards mark u black.
//  4. Run DFS from every white node; return true iff no DFS found a cycle.
//
// Time:  O(V + E) — each node and edge is processed exactly once.
// Space: O(V + E) — adjacency list, state array, and recursion stack.
func dfsCycleDetection(numCourses int, prerequisites [][]int) bool {
	// adj[b] lists every course unlocked by b (edge b → a).
	adj := make([][]int, numCourses)
	for _, p := range prerequisites {
		adj[p[1]] = append(adj[p[1]], p[0]) // b → a: finish b before a
	}

	const (
		white = 0 // never visited
		grey  = 1 // on the current DFS path (still being explored)
		black = 2 // completely explored — provably cycle-free below it
	)
	state := make([]int, numCourses) // all start white (zero value)

	// dfs returns true if a cycle is reachable from node u.
	var dfs func(u int) bool
	dfs = func(u int) bool {
		state[u] = grey // u joins the active path
		for _, v := range adj[u] {
			if state[v] == grey {
				return true // back edge to the active path → cycle
			}
			if state[v] == white && dfs(v) {
				return true // a deeper call found a cycle — bubble it up
			}
			// black neighbours are already proven safe — skip them
		}
		state[u] = black // u fully explored, nothing below it cycles
		return false
	}

	// The graph may be disconnected — start a DFS from every unvisited node.
	for c := 0; c < numCourses; c++ {
		if state[c] == white && dfs(c) {
			return false // any cycle makes the schedule impossible
		}
	}
	return true // no cycle anywhere → every course can be finished
}

// ── Approach 3: Kahn's Algorithm / BFS Topological Sort (Optimal) ────────────
//
// bfsKahn solves Course Schedule by peeling off zero-in-degree courses layer
// by layer; all courses get peeled iff the graph is acyclic.
//
// Intuition:
//
//	A course with in-degree 0 has no pending prerequisites — take it now.
//	Taking it "removes" its outgoing edges, possibly dropping other courses
//	to in-degree 0. Repeat with a queue. Nodes on a cycle can never reach
//	in-degree 0 (each waits on another cycle member), so counting how many
//	courses got processed tells us whether a full ordering exists.
//
// Algorithm:
//  1. Build adjacency list (edge b → a per pair [a, b]) and indegree[a]++.
//  2. Enqueue every course with in-degree 0.
//  3. Pop u: count it processed; decrement in-degree of each neighbour v,
//     enqueueing v when its in-degree hits 0.
//  4. Return processed == numCourses.
//
// Time:  O(V + E) — each course enqueued once, each edge relaxed once.
// Space: O(V + E) — adjacency list, in-degree array, queue.
func bfsKahn(numCourses int, prerequisites [][]int) bool {
	adj := make([][]int, numCourses)    // adj[b] = courses unlocked by b
	indegree := make([]int, numCourses) // indegree[a] = prerequisites of a still unmet
	for _, p := range prerequisites {
		adj[p[1]] = append(adj[p[1]], p[0]) // edge b → a
		indegree[p[0]]++                    // a waits on one more course
	}

	// Seed the queue with all courses takeable right now (no prerequisites).
	queue := []int{}
	for c := 0; c < numCourses; c++ {
		if indegree[c] == 0 {
			queue = append(queue, c)
		}
	}

	processed := 0 // number of courses successfully ordered
	for len(queue) > 0 {
		u := queue[0]     // take the front course
		queue = queue[1:] // dequeue it
		processed++       // u is now "taken"
		for _, v := range adj[u] {
			indegree[v]-- // u is done, so v has one fewer unmet prerequisite
			if indegree[v] == 0 {
				queue = append(queue, v) // all of v's prerequisites met → takeable
			}
		}
	}

	// Courses stuck on a cycle never reach in-degree 0 and are never counted.
	return processed == numCourses
}

func main() {
	fmt.Println("=== Approach 1: Brute Force (Repeated Elimination) ===")
	fmt.Println(bruteForce(2, [][]int{{1, 0}}))         // true
	fmt.Println(bruteForce(2, [][]int{{1, 0}, {0, 1}})) // false

	fmt.Println("=== Approach 2: DFS Cycle Detection (3-Colour) ===")
	fmt.Println(dfsCycleDetection(2, [][]int{{1, 0}}))         // true
	fmt.Println(dfsCycleDetection(2, [][]int{{1, 0}, {0, 1}})) // false

	fmt.Println("=== Approach 3: Kahn's Algorithm / BFS Topological Sort (Optimal) ===")
	fmt.Println(bfsKahn(2, [][]int{{1, 0}}))         // true
	fmt.Println(bfsKahn(2, [][]int{{1, 0}, {0, 1}})) // false
}
