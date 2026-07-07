# 0207 — Course Schedule

> LeetCode #207 · Difficulty: Medium
> **Categories:** Graph, Topological Sort, Depth-First Search, Breadth-First Search

---

## Problem Statement

There are a total of `numCourses` courses you have to take, labeled from `0` to `numCourses - 1`. You are given an array `prerequisites` where `prerequisites[i] = [aᵢ, bᵢ]` indicates that you **must** take course `bᵢ` first if you want to take course `aᵢ`.

- For example, the pair `[0, 1]` indicates that to take course `0` you have to first take course `1`.

Return `true` if you can finish all courses. Otherwise, return `false`.

**Example 1:**
```
Input: numCourses = 2, prerequisites = [[1,0]]
Output: true
Explanation: There are a total of 2 courses to take.
To take course 1 you should have finished course 0. So it is possible.
```

**Example 2:**
```
Input: numCourses = 2, prerequisites = [[1,0],[0,1]]
Output: false
Explanation: There are a total of 2 courses to take.
To take course 1 you should have finished course 0, and to take course 0 you should also have finished course 1. So it is impossible.
```

**Constraints:**
- `1 <= numCourses <= 2000`
- `0 <= prerequisites.length <= 5000`
- `prerequisites[i].length == 2`
- `0 <= aᵢ, bᵢ < numCourses`
- All the pairs `prerequisites[i]` are **unique**.

---

## Company Frequency

| Company    | Frequency        | Last Reported |
|------------|------------------|---------------|
| Amazon     | ★★★★★ Very High  | 2024          |
| Google     | ★★★★☆ High       | 2024          |
| Meta       | ★★★★☆ High       | 2024          |
| Microsoft  | ★★★☆☆ Medium     | 2024          |
| TikTok     | ★★★☆☆ Medium     | 2024          |
| Uber       | ★★★☆☆ Medium     | 2023          |
| Apple      | ★★☆☆☆ Low        | 2023          |
| Oracle     | ★★☆☆☆ Low        | 2023          |

> ⚠️ Frequency data is crowd-sourced from LeetCode Discuss, Glassdoor, and
> community interview reports. Treat as a signal, not a guarantee.

---

## DSA Concepts Used

- **Graph (Directed)** — courses are nodes; each pair `[a, b]` is a directed edge `b → a` ("finish b unlocks a"); the question reduces to a structural property of this graph → see [`/dsa/graph_bfs_dfs.md`](/dsa/graph_bfs_dfs.md)
- **Topological Sort** — "can all courses be finished?" is exactly "does a topological order exist?", which holds iff the graph is a DAG → see [`/dsa/topological_sort.md`](/dsa/topological_sort.md)
- **Cycle Detection (DFS 3-colour)** — a back edge to a node on the current DFS path proves a cycle → see [`/dsa/graph_bfs_dfs.md`](/dsa/graph_bfs_dfs.md)
- **Queue** — Kahn's algorithm processes zero-in-degree nodes in FIFO order → see [`/dsa/queue_deque.md`](/dsa/queue_deque.md)

---

## Approaches Overview

| # | Approach | Time | Space | When to use |
|---|----------|------|-------|-------------|
| 1 | Brute Force (Repeated Elimination) | O(V·(V+E)) | O(V) | Baseline; shows *why* cycles block progress, fine for tiny inputs |
| 2 | DFS Cycle Detection (3-Colour) | O(V+E) | O(V+E) | When you only need yes/no; natural if you think recursively |
| 3 | Kahn's / BFS Topological Sort (Optimal) | O(V+E) | O(V+E) | Interview default; iterative, no stack-overflow risk, extends directly to #210 |

---

## Approach 1 — Brute Force (Repeated Elimination)

### Intuition
Simulate a student with no planning skills. Each "semester", scan the entire catalogue and enrol in every course whose prerequisites are already complete. Repeat semesters until either everything is taken (schedule feasible) or a whole semester passes with zero new enrolments. In the stuck case every remaining course waits on another remaining course — which is precisely a prerequisite cycle — so the answer is `false`. This is Kahn's algorithm without the bookkeeping: instead of tracking in-degrees, it re-derives readiness by brute scanning.

### Algorithm
1. Create `taken[0..numCourses-1]`, all `false`, and `count = 0`.
2. Loop over "sweeps":
   1. For every course `c` not yet taken, scan **all** prerequisite pairs; `c` is ready iff every pair `[c, b]` has `taken[b] == true`.
   2. If ready, set `taken[c] = true`, increment `count`, and mark that this sweep made progress.
3. Stop when a sweep makes no progress.
4. Return `count == numCourses`.

### Complexity
- **Time:** O(V·(V+E)) — at most V sweeps (each successful sweep takes ≥ 1 course), and every sweep scans V courses × the full E-pair list.
- **Space:** O(V) — only the `taken` array.

### Code
```go
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
```

### Dry Run (Example 1: numCourses = 2, prerequisites = [[1,0]])

| Sweep | Course examined | Pairs constraining it | Ready? | `taken` after | `count` | `progress` |
|-------|-----------------|------------------------|--------|---------------|---------|------------|
| 1 | 0 | none ([1,0] constrains course 1, not 0) | yes | [true, false] | 1 | true |
| 1 | 1 | [1,0] → needs course 0, taken ✓ | yes | [true, true] | 2 | true |
| 2 | — (both taken, nothing examined) | — | — | [true, true] | 2 | false → stop |

`count (2) == numCourses (2)` → return **true** ✓

(Example 2 for contrast: sweep 1 finds course 0 waiting on 1 and course 1 waiting on 0 — no progress, loop stops, `count = 0 ≠ 2` → `false`.)

---

## Approach 2 — DFS Cycle Detection (3-Colour)

### Intuition
Build the graph with an edge `b → a` for every pair `[a, b]`. All courses are finishable **iff this directed graph has no cycle** (a cycle means "a before b before … before a", which is unsatisfiable; without cycles, any topological order works). DFS detects directed cycles with three node colours: **white** = untouched, **grey** = on the current recursion path, **black** = fully explored and proven safe. Reaching a *grey* node again means the current path looped back onto itself → cycle. Reaching a *black* node is harmless — everything beneath it was already checked.

### Algorithm
1. Build adjacency list: for each `[a, b]`, append `a` to `adj[b]`.
2. Initialise `state[i] = white` for all courses.
3. Define `dfs(u)` → "cycle reachable from u?":
   1. Mark `state[u] = grey`.
   2. For each neighbour `v`: if `v` is grey → return true; if `v` is white and `dfs(v)` → return true.
   3. Mark `state[u] = black`; return false.
4. For every white course `c`, run `dfs(c)`; if any returns true → return `false`.
5. Otherwise return `true`.

### Complexity
- **Time:** O(V+E) — every node is DFS-entered once (white→grey→black is one-way), every edge examined once.
- **Space:** O(V+E) — adjacency list + state array + recursion stack (up to V deep).

### Code
```go
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
```

### Dry Run (Example 1: numCourses = 2, prerequisites = [[1,0]])

Graph: `adj[0] = [1]`, `adj[1] = []`. States shown as [course0, course1].

| Step | Action | `state` | Result |
|------|--------|---------|--------|
| 1 | Outer loop: course 0 is white → `dfs(0)` | [grey, white] | — |
| 2 | `dfs(0)` scans neighbour 1: white → recurse `dfs(1)` | [grey, grey] | — |
| 3 | `dfs(1)` has no neighbours → mark 1 black, return false | [grey, black] | no cycle below 1 |
| 4 | back in `dfs(0)`: neighbours done → mark 0 black, return false | [black, black] | no cycle below 0 |
| 5 | Outer loop: course 1 already black → skip | [black, black] | — |

No DFS reported a cycle → return **true** ✓

(Example 2: `dfs(0)` goes 0(grey) → 1(grey) → neighbour 0 is **grey** → cycle → `false`.)

---

## Approach 3 — Kahn's Algorithm / BFS Topological Sort (Optimal)

### Intuition
A course with **in-degree 0** has no unmet prerequisites — it can be taken immediately. Taking it conceptually deletes its outgoing edges, which may drop other courses to in-degree 0, making them takeable next. Process such courses with a queue and count how many get taken. Nodes on a cycle can never reach in-degree 0 (each cycle member permanently waits on the previous one), so `processed == numCourses` is exactly the acyclicity test. Bonus: the order in which nodes leave the queue *is* a valid course order — which is why this approach extends verbatim to #210.

### Algorithm
1. Build `adj[b] ← a` and `indegree[a]++` for every pair `[a, b]`.
2. Enqueue every course whose `indegree` is 0.
3. While the queue is non-empty:
   1. Dequeue `u`; increment `processed`.
   2. For each neighbour `v` of `u`: decrement `indegree[v]`; if it reaches 0, enqueue `v`.
4. Return `processed == numCourses`.

### Complexity
- **Time:** O(V+E) — each course enters the queue at most once; each edge decrements one in-degree exactly once.
- **Space:** O(V+E) — adjacency list (E entries), in-degree array and queue (V entries).

### Code
```go
func bfsKahn(numCourses int, prerequisites [][]int) bool {
	adj := make([][]int, numCourses)     // adj[b] = courses unlocked by b
	indegree := make([]int, numCourses)  // indegree[a] = prerequisites of a still unmet
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
```

### Dry Run (Example 1: numCourses = 2, prerequisites = [[1,0]])

Setup: `adj[0] = [1]`, `adj[1] = []`, `indegree = [0, 1]`. Course 0 has in-degree 0 → initial `queue = [0]`.

| Step | Dequeued `u` | `processed` | Edge relaxed | `indegree` after | Queue after |
|------|--------------|-------------|--------------|------------------|-------------|
| init | — | 0 | — | [0, 1] | [0] |
| 1 | 0 | 1 | 0→1: indegree[1] 1→0 → enqueue 1 | [0, 0] | [1] |
| 2 | 1 | 2 | (no outgoing edges) | [0, 0] | [] |

Queue empty; `processed (2) == numCourses (2)` → return **true** ✓

(Example 2: `indegree = [1, 1]` — nothing ever enters the queue, `processed = 0 ≠ 2` → `false`.)

---

## Key Takeaways

- **Recognise the reduction:** "can all tasks with dependencies be completed?" = "is the dependency graph a DAG?" = "does a topological order exist?". This phrasing appears with courses, build systems, spreadsheet formulas, and recipe crafting.
- **Watch the edge direction.** `[a, b]` = "b before a" → edge `b → a`. Flipping it silently still *detects cycles* correctly (a cycle reversed is still a cycle), which hides the bug until you need the actual order in #210 — build it right from the start.
- **Kahn's invariant:** a node enters the queue exactly when its last unmet dependency completes; cycle members never do. `processed == V` doubles as the acyclicity check for free.
- **3-colour DFS:** grey-hit = cycle is the directed-graph cycle test; a plain visited boolean is *not* enough for directed graphs (black nodes must be revisitable without alarm).
- **DFS vs BFS choice:** DFS is shorter to write recursively but risks deep recursion (V up to 2000 here is fine); Kahn's is iterative, cycle-proof by counting, and yields the ordering itself.
- The brute-force elimination is Kahn's algorithm minus the in-degree bookkeeping — a good way to *derive* Kahn's in an interview if you blank on it.

---

## Related Problems

- LeetCode #210 — Course Schedule II (same graph; return the actual order)
- LeetCode #269 — Alien Dictionary (build the graph from words, then topo-sort)
- LeetCode #1462 — Course Schedule IV (prerequisite reachability queries)
- LeetCode #630 — Course Schedule III (greedy + heap, different flavour)
- LeetCode #802 — Find Eventual Safe States (reverse-graph topological peel / DFS colours)
- LeetCode #444 — Sequence Reconstruction (uniqueness of a topological order)
- LeetCode #2115 — Find All Possible Recipes from Given Supplies (Kahn's on ingredients)
