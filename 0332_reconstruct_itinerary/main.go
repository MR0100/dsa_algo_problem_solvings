package main

import (
	"fmt"
	"sort"
)

// ── Approach 1: Backtracking (DFS + Undo) ────────────────────────────────────
//
// backtracking reconstructs the itinerary by trying destinations in
// lexical order and undoing whenever a branch fails to use all tickets.
//
// Intuition:
//
//	We must use every ticket exactly once (an Eulerian path) and, among all
//	such paths, pick the lexicographically smallest. Greedily always taking
//	the smallest available next airport can strand us (a ticket left unused),
//	so we backtrack: pick the smallest unused edge, recurse, and if the tail
//	cannot consume all remaining tickets, put the edge back and try the next.
//
// Algorithm:
//  1. Build adj[from] = sorted list of destinations; keep a used[] flag per edge.
//  2. DFS(node): if all tickets used, done. Else for each unused edge from node
//     in lexical order: mark used, append dest, recurse; if it completed, return
//     true; otherwise unmark and pop (undo).
//
// Time:  O(E^d) worst case (d = max out-degree) — exponential backtracking.
// Space: O(E) for the route/recursion.
func backtracking(tickets [][]string) []string {
	adj := map[string][]string{} // from → sorted destinations
	for _, t := range tickets {
		adj[t[0]] = append(adj[t[0]], t[1])
	}
	for k := range adj {
		sort.Strings(adj[k]) // lexical order so we try smallest first
	}
	used := map[string][]bool{} // parallel "edge already flown" flags
	for k, v := range adj {
		used[k] = make([]bool, len(v))
	}
	total := len(tickets) + 1 // number of airports in a full itinerary
	route := []string{"JFK"}  // every itinerary starts at JFK
	var dfs func() bool
	dfs = func() bool {
		if len(route) == total {
			return true // used all tickets → valid complete route
		}
		cur := route[len(route)-1] // where we are now
		for i, dest := range adj[cur] {
			if used[cur][i] {
				continue // this ticket already flown on this path
			}
			used[cur][i] = true         // fly it
			route = append(route, dest) // extend the route
			if dfs() {
				return true // downstream completed the itinerary
			}
			used[cur][i] = false         // undo: mark ticket unused again
			route = route[:len(route)-1] // undo: drop the airport
		}
		return false // no destination from here completes the route
	}
	dfs()
	return route
}

// ── Approach 2: Hierholzer's Algorithm (Optimal) ─────────────────────────────
//
// hierholzer reconstructs the itinerary as an Eulerian path using
// Hierholzer's post-order edge-consumption, no backtracking needed.
//
// Intuition:
//
//	The graph is guaranteed to have an Eulerian path from JFK. Hierholzer:
//	greedily walk edges (always the smallest destination) until stuck; the
//	airport where you get stuck is the true end and is added to the route
//	first. Recording nodes in POST-ORDER (after exhausting their edges) and
//	reversing yields a valid Eulerian path — and because we always take the
//	lexically smallest edge, that path is the smallest one.
//
// Algorithm:
//  1. Build adj[from] = min-heap-like sorted destinations (consumed front to back).
//  2. DFS(node): while node has unused edges, take the smallest and recurse;
//     when node has no edges left, prepend it to the answer (post-order).
//  3. Reverse the post-order list → the itinerary.
//
// Time:  O(E log E) — sorting each adjacency list dominates.
// Space: O(E) — adjacency lists + recursion + output.
func hierholzer(tickets [][]string) []string {
	adj := map[string][]string{} // from → sorted destinations (a queue)
	for _, t := range tickets {
		adj[t[0]] = append(adj[t[0]], t[1])
	}
	for k := range adj {
		sort.Strings(adj[k]) // smallest destination first
	}
	route := []string{} // built in reverse (post-order)
	var dfs func(node string)
	dfs = func(node string) {
		// Consume edges from this node until none remain.
		for len(adj[node]) > 0 {
			next := adj[node][0]      // smallest available destination
			adj[node] = adj[node][1:] // remove that edge (used exactly once)
			dfs(next)                 // walk deeper before recording node
		}
		route = append(route, node) // post-order: add after edges exhausted
	}
	dfs("JFK")
	// route is in reverse Eulerian order; reverse it in place.
	for i, j := 0, len(route)-1; i < j; i, j = i+1, j-1 {
		route[i], route[j] = route[j], route[i]
	}
	return route
}

func main() {
	ex1 := [][]string{{"MUC", "LHR"}, {"JFK", "MUC"}, {"SFO", "SJC"}, {"LHR", "SFO"}}
	ex2 := [][]string{{"JFK", "SFO"}, {"JFK", "ATL"}, {"SFO", "ATL"}, {"ATL", "JFK"}, {"ATL", "SFO"}}

	fmt.Println("=== Approach 1: Backtracking (DFS + Undo) ===")
	fmt.Println(backtracking(ex1)) // expected [JFK MUC LHR SFO SJC]
	fmt.Println(backtracking(ex2)) // expected [JFK ATL JFK SFO ATL SFO]

	fmt.Println("=== Approach 2: Hierholzer's Algorithm (Optimal) ===")
	fmt.Println(hierholzer(ex1)) // expected [JFK MUC LHR SFO SJC]
	fmt.Println(hierholzer(ex2)) // expected [JFK ATL JFK SFO ATL SFO]
}
