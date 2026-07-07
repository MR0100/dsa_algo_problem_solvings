package main

import "fmt"

// ── Approach 1: Brute Force (Repeated Elimination) ───────────────────────────
//
// bruteForce solves Course Schedule II by repeatedly sweeping the catalogue
// and taking, in index order, every course whose prerequisites are already
// taken — recording the order in which courses get taken.
//
// Intuition:
//
//	Same simulation as #207's brute force, but now we write down the order.
//	Each sweep enrols every course that has become available; courses taken
//	earlier in the same sweep immediately unlock later-indexed ones. If a
//	full sweep enrols nobody while courses remain, the leftovers form a
//	cycle and no ordering exists — return the empty slice.
//
// Algorithm:
//  1. taken[i] tracks completion; order accumulates the schedule.
//  2. Sweep courses 0..n-1: any untaken course whose every pair [c, b] has
//     taken[b] == true is appended to order and marked taken.
//  3. Repeat sweeps until one makes no progress.
//  4. Return order if it contains all courses, else an empty slice.
//
// Time:  O(V · (V + E)) — up to V sweeps, each scanning V courses × E pairs.
// Space: O(V) — the taken array and the output order.
func bruteForce(numCourses int, prerequisites [][]int) []int {
	taken := make([]bool, numCourses)   // taken[i] = course i already completed
	order := make([]int, 0, numCourses) // the schedule being built

	for {
		progress := false // does this sweep enrol anyone?
		for c := 0; c < numCourses; c++ {
			if taken[c] {
				continue // course already scheduled
			}
			// Course c is ready iff every pair [c, b] has b taken.
			ready := true
			for _, p := range prerequisites {
				// p = [a, b]: to take course a you must first take course b.
				if p[0] == c && !taken[p[1]] {
					ready = false // some prerequisite of c still missing
					break
				}
			}
			if ready {
				taken[c] = true          // enrol course c now
				order = append(order, c) // record its position in the schedule
				progress = true
			}
		}
		if !progress {
			break // no course could be added — either done or stuck on a cycle
		}
	}

	if len(order) != numCourses {
		return []int{} // cycle: some courses can never be scheduled
	}
	return order
}

// ── Approach 2: DFS Postorder + Reverse ──────────────────────────────────────
//
// dfsPostorder solves Course Schedule II by depth-first searching the
// prerequisite graph, emitting each course after all courses depending on it,
// then reversing that postorder to get a valid schedule.
//
// Intuition:
//
//	With edges b → a ("b unlocks a"), DFS finishes a node only after every
//	node reachable from it is finished. So the DFS finish (postorder) lists
//	each course AFTER all its dependents — the exact reverse of a valid
//	schedule. Reverse the postorder and prerequisites come first. Cycles are
//	caught with the same 3-colour states as #207: revisiting a grey node
//	means the current path bit its own tail.
//
// Algorithm:
//  1. Build adjacency list adj[b] ← a for each pair [a, b].
//  2. state[i] ∈ {white, grey, black}; post = empty list.
//  3. dfs(u): mark grey; recurse into white neighbours (grey neighbour →
//     cycle); mark black and append u to post.
//  4. Run dfs from every white node; on any cycle return the empty slice.
//  5. Reverse post in place and return it.
//
// Time:  O(V + E) — each node entered once, each edge crossed once.
// Space: O(V + E) — adjacency list, states, postorder list, recursion stack.
func dfsPostorder(numCourses int, prerequisites [][]int) []int {
	// adj[b] = all courses that need b done first (edge b → a).
	adj := make([][]int, numCourses)
	for _, p := range prerequisites {
		adj[p[1]] = append(adj[p[1]], p[0]) // finish b before a
	}

	const (
		white = 0 // untouched
		grey  = 1 // on the current DFS path
		black = 2 // fully explored (already in post)
	)
	state := make([]int, numCourses)
	post := make([]int, 0, numCourses) // DFS finish order (reverse schedule)

	// dfs returns false if it finds a cycle reachable from u.
	var dfs func(u int) bool
	dfs = func(u int) bool {
		state[u] = grey // u is now on the active path
		for _, v := range adj[u] {
			if state[v] == grey {
				return false // back edge: the active path loops → cycle
			}
			if state[v] == white && !dfs(v) {
				return false // cycle found deeper down — abort everything
			}
		}
		state[u] = black       // all dependents of u fully processed
		post = append(post, u) // u finishes AFTER everyone who needs it
		return true
	}

	// Cover disconnected components: try every course as a DFS root.
	for c := 0; c < numCourses; c++ {
		if state[c] == white && !dfs(c) {
			return []int{} // impossible schedule
		}
	}

	// post is "dependents first"; a schedule needs prerequisites first.
	for i, j := 0, len(post)-1; i < j; i, j = i+1, j-1 {
		post[i], post[j] = post[j], post[i] // in-place reversal
	}
	return post
}

// ── Approach 3: Kahn's Algorithm / BFS Topological Sort (Optimal) ────────────
//
// bfsKahn solves Course Schedule II by dequeuing zero-in-degree courses; the
// dequeue order IS a topological order of the prerequisite graph.
//
// Intuition:
//
//	Identical machinery to #207, except the point here is the by-product:
//	every course leaves the queue only once its final prerequisite has left
//	before it, so appending courses as they dequeue directly builds a valid
//	schedule. If some courses never reach in-degree 0 (a cycle starves
//	them), fewer than numCourses dequeue and we report impossibility.
//
// Algorithm:
//  1. Build adj[b] ← a and indegree[a]++ per pair [a, b].
//  2. Enqueue all courses with in-degree 0.
//  3. Dequeue u → append to order; decrement each neighbour's in-degree,
//     enqueueing any that reach 0.
//  4. Return order if len(order) == numCourses, else the empty slice.
//
// Time:  O(V + E) — every course enqueued/dequeued once, every edge relaxed once.
// Space: O(V + E) — adjacency list, in-degree array, queue, output.
func bfsKahn(numCourses int, prerequisites [][]int) []int {
	adj := make([][]int, numCourses)    // adj[b] = courses unlocked by b
	indegree := make([]int, numCourses) // unmet-prerequisite count per course
	for _, p := range prerequisites {
		adj[p[1]] = append(adj[p[1]], p[0]) // edge b → a
		indegree[p[0]]++                    // course a waits on one more
	}

	// Start with everything takeable immediately (no prerequisites).
	queue := []int{}
	for c := 0; c < numCourses; c++ {
		if indegree[c] == 0 {
			queue = append(queue, c)
		}
	}

	order := make([]int, 0, numCourses) // schedule in dequeue order
	for len(queue) > 0 {
		u := queue[0]            // next course with all prerequisites done
		queue = queue[1:]        // pop the front
		order = append(order, u) // taking u right now is safe
		for _, v := range adj[u] {
			indegree[v]-- // u done → v loses one unmet prerequisite
			if indegree[v] == 0 {
				queue = append(queue, v) // v just became takeable
			}
		}
	}

	if len(order) != numCourses {
		return []int{} // cycle members never dequeued — no valid order
	}
	return order
}

func main() {
	// Note: any valid topological order is accepted by LeetCode. The inline
	// expectations below are the deterministic outputs of *these* functions.

	fmt.Println("=== Approach 1: Brute Force (Repeated Elimination) ===")
	fmt.Println(bruteForce(2, [][]int{{1, 0}}))                         // [0 1]
	fmt.Println(bruteForce(4, [][]int{{1, 0}, {2, 0}, {3, 1}, {3, 2}})) // [0 1 2 3]
	fmt.Println(bruteForce(1, [][]int{}))                               // [0]

	fmt.Println("=== Approach 2: DFS Postorder + Reverse ===")
	fmt.Println(dfsPostorder(2, [][]int{{1, 0}}))                         // [0 1]
	fmt.Println(dfsPostorder(4, [][]int{{1, 0}, {2, 0}, {3, 1}, {3, 2}})) // [0 2 1 3] (also valid)
	fmt.Println(dfsPostorder(1, [][]int{}))                               // [0]

	fmt.Println("=== Approach 3: Kahn's Algorithm / BFS Topological Sort (Optimal) ===")
	fmt.Println(bfsKahn(2, [][]int{{1, 0}}))                         // [0 1]
	fmt.Println(bfsKahn(4, [][]int{{1, 0}, {2, 0}, {3, 1}, {3, 2}})) // [0 1 2 3]
	fmt.Println(bfsKahn(1, [][]int{}))                               // [0]
}
