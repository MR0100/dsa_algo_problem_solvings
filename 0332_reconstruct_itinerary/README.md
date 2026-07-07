# 0332 — Reconstruct Itinerary

> LeetCode #332 · Difficulty: Hard
> **Categories:** Graph, Depth-First Search, Eulerian Path, Backtracking

---

## Problem Statement

You are given a list of airline `tickets` where `tickets[i] = [from_i, to_i]` represent the departure and the arrival airports of one flight. Reconstruct the itinerary in order and return it.

All of the tickets belong to a man who departs from `"JFK"`, thus, the itinerary must begin with `"JFK"`. If there are multiple valid itineraries, you should return the itinerary that has the smallest lexical order when read as a single string.

- For example, the itinerary `["JFK", "LGA"]` has a smaller lexical order than `["JFK", "LGB"]`.

You may assume all tickets form at least one valid itinerary. You must use all the tickets once and only once.

**Example 1:**

```
Input: tickets = [["MUC","LHR"],["JFK","MUC"],["SFO","SJC"],["LHR","SFO"]]
Output: ["JFK","MUC","LHR","SFO","SJC"]
```

**Example 2:**

```
Input: tickets = [["JFK","SFO"],["JFK","ATL"],["SFO","ATL"],["ATL","JFK"],["ATL","SFO"]]
Output: ["JFK","ATL","JFK","SFO","ATL","SFO"]
Explanation: Another possible reconstruction is ["JFK","SFO","ATL","JFK","ATL","SFO"] but it is larger in lexical order.
```

**Constraints:**

- `1 <= tickets.length <= 300`
- `tickets[i].length == 2`
- `from_i.length == 3`
- `to_i.length == 3`
- `from_i` and `to_i` consist of uppercase English letters.
- `from_i != to_i`

---

## Company Frequency

| Company    | Frequency        | Last Reported |
|------------|------------------|---------------|
| Google     | ★★★★☆ High       | 2024          |
| Amazon     | ★★★★☆ High       | 2024          |
| Meta       | ★★★☆☆ Medium     | 2023          |
| Microsoft  | ★★★☆☆ Medium     | 2023          |
| Bloomberg  | ★★☆☆☆ Low        | 2023          |
| Uber       | ★★☆☆☆ Low        | 2022          |

> ⚠️ Frequency data is crowd-sourced from LeetCode Discuss, Glassdoor, and
> community interview reports. Treat as a signal, not a guarantee.

---

## DSA Concepts Used

- **Graph DFS** — the tickets form a directed multigraph; both solutions walk it with depth-first search consuming edges → see [`/dsa/graph_bfs_dfs.md`](/dsa/graph_bfs_dfs.md)
- **Backtracking** — Approach 1 tries the smallest edge and undoes it when the tail can't consume all remaining tickets → see [`/dsa/backtracking.md`](/dsa/backtracking.md)
- **Sorting** — visiting destinations in lexical order (sorted adjacency lists) is what makes the result the smallest valid itinerary → see [`/dsa/sorting.md`](/dsa/sorting.md)
- **Eulerian Path (Hierholzer's Algorithm)** — the optimal solution recognizes the itinerary as an Eulerian path and builds it in post-order. No dedicated file exists; closest is graph DFS above.

---

## Approaches Overview

| # | Approach | Time | Space | When to use |
|---|----------|------|-------|-------------|
| 1 | Backtracking (DFS + Undo) | O(E^d) worst | O(E) | When you don't recall Hierholzer; simple to reason about, but can blow up |
| 2 | Hierholzer's Algorithm (Optimal) | O(E log E) | O(E) | The intended solution — linear-ish, no backtracking, provably smallest |

*(E = number of tickets, d = max out-degree of a node.)*

---

## Approach 1 — Backtracking (DFS + Undo)

### Intuition

We need a path that uses every ticket exactly once (an **Eulerian path**), and among all such paths the lexicographically smallest. A pure greedy "always fly to the smallest next airport" can dead-end with tickets unused — imagine flying into a leaf that has no way back. So we treat it as search: sort each airport's destinations, try the smallest unused ticket, recurse, and if the remaining flights can't all be consumed, *undo* that ticket and try the next. Because we try destinations in ascending order, the first complete itinerary we finish is the smallest.

### Algorithm

1. Build `adj[from]` = the sorted list of destinations, with a parallel `used[]` flag array per source.
2. Start `route = ["JFK"]`. Define `dfs()`:
   - If `len(route) == tickets+1`, all tickets are used → return `true`.
   - From the current airport, for each unused edge in lexical order: mark it used, append the destination, recurse. If the recursion returns `true`, propagate `true`. Otherwise unmark and pop (backtrack).
3. Run `dfs()`; `route` now holds the answer.

### Complexity

- **Time:** O(E^d) worst case — at each of E steps we may branch over up to d outgoing edges and backtrack; exponential but fine for E ≤ 300 in practice due to the tree structure.
- **Space:** O(E) — the route, the `used` flags, and recursion depth.

### Code

```go
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
	total := len(tickets) + 1     // number of airports in a full itinerary
	route := []string{"JFK"}      // every itinerary starts at JFK
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
```

### Dry Run

Example 1: `tickets = [[MUC,LHR],[JFK,MUC],[SFO,SJC],[LHR,SFO]]`, `total = 5`.

Sorted adjacency: `JFK→[MUC]`, `MUC→[LHR]`, `LHR→[SFO]`, `SFO→[SJC]`.

| Step | route | current | edge tried | result |
|------|-------|---------|-----------|--------|
| 1 | `[JFK]` | JFK | MUC | fly, recurse |
| 2 | `[JFK,MUC]` | MUC | LHR | fly, recurse |
| 3 | `[JFK,MUC,LHR]` | LHR | SFO | fly, recurse |
| 4 | `[JFK,MUC,LHR,SFO]` | SFO | SJC | fly, recurse |
| 5 | `[JFK,MUC,LHR,SFO,SJC]` | — | — | `len==5` → return true |

No backtracking needed here (linear chain). Result: `[JFK,MUC,LHR,SFO,SJC]` ✔

---

## Approach 2 — Hierholzer's Algorithm (Optimal)

### Intuition

The problem *is* "find an Eulerian path from JFK", and one is guaranteed to exist. Hierholzer's insight: greedily walk edges (always the smallest destination) until you get **stuck** at a node with no outgoing edges left. That stuck node must be the itinerary's final airport, so record it **first**. Recording each node in *post-order* — only after all its edges are exhausted — and then reversing the list yields a valid Eulerian path. Because we always consume the lexically smallest edge available, the reversed post-order is the smallest valid itinerary. No undo is ever required: any edge we "skip past" gets stitched in as a detour earlier in the reversal.

### Algorithm

1. Build `adj[from]` = sorted list of destinations, used as a queue (front = smallest).
2. Define `dfs(node)`: while `node` still has edges, pop its smallest destination and `dfs` into it; once `node` has no edges left, append `node` to `route` (post-order).
3. Call `dfs("JFK")`, then reverse `route`.

### Complexity

- **Time:** O(E log E) — sorting the adjacency lists dominates; the traversal itself visits each edge once, O(E).
- **Space:** O(E) — adjacency lists, recursion stack, and the output route.

### Code

```go
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
			next := adj[node][0]     // smallest available destination
			adj[node] = adj[node][1:] // remove that edge (used exactly once)
			dfs(next)                // walk deeper before recording node
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
```

### Dry Run

Example 1: `tickets = [[MUC,LHR],[JFK,MUC],[SFO,SJC],[LHR,SFO]]`.

Adjacency: `JFK→[MUC]`, `MUC→[LHR]`, `LHR→[SFO]`, `SFO→[SJC]`.

| Call | Node | Edges left | Action |
|------|------|-----------|--------|
| dfs(JFK) | JFK | [MUC] | take MUC → dfs(MUC) |
| dfs(MUC) | MUC | [LHR] | take LHR → dfs(LHR) |
| dfs(LHR) | LHR | [SFO] | take SFO → dfs(SFO) |
| dfs(SFO) | SFO | [SJC] | take SJC → dfs(SJC) |
| dfs(SJC) | SJC | [] | append SJC → route=`[SJC]` |
| ↑ back in SFO | SFO | [] | append SFO → route=`[SJC,SFO]` |
| ↑ back in LHR | LHR | [] | append LHR → route=`[SJC,SFO,LHR]` |
| ↑ back in MUC | MUC | [] | append MUC → route=`[SJC,SFO,LHR,MUC]` |
| ↑ back in JFK | JFK | [] | append JFK → route=`[SJC,SFO,LHR,MUC,JFK]` |

Reverse: `[JFK,MUC,LHR,SFO,SJC]`. Result ✔

---

## Key Takeaways

- **Recognize the Eulerian-path signature:** "use every edge exactly once" + "a valid one is guaranteed" ⇒ Hierholzer's algorithm. Post-order append then reverse is the whole trick.
- **Lexical-smallest comes for free** if adjacency lists are sorted and you always consume the front — greedy edge choice does not need backtracking under Hierholzer because skipped detours re-attach during the reversal.
- **Why naive greedy fails but Hierholzer doesn't:** greedy commits to a full-path prefix and can strand tickets; Hierholzer only commits at the moment a node is *exhausted*, so a premature dead-end simply becomes the tail of the route.
- Backtracking is the safe fallback when you can't recall Hierholzer, but its worst case is exponential — know the linear algorithm for Hard interviews.

---

## Related Problems

- LeetCode #753 — Cracking the Safe (Eulerian path on a de Bruijn graph)
- LeetCode #2097 — Valid Arrangement of Pairs (general Eulerian path reconstruction)
- LeetCode #207 — Course Schedule (directed-graph DFS / cycle reasoning)
- LeetCode #797 — All Paths From Source to Target (DFS path enumeration)
- LeetCode #1743 — Restore the Array From Adjacent Pairs (reconstruct a sequence from edges)
