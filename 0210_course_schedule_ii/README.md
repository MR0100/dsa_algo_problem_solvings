# 0210 — Course Schedule II

> LeetCode #210 · Difficulty: Medium
> **Categories:** Graph, Topological Sort, BFS, DFS, Depth-First Search, Breadth-First Search

---

## Problem Statement

There are a total of `numCourses` courses you have to take, labeled from `0` to `numCourses - 1`. You are given an array `prerequisites` where `prerequisites[i] = [ai, bi]` indicates that you **must** take course `bi` first if you want to take course `ai`.

- For example, the pair `[0, 1]`, indicates that to take course `0` you have to first take course `1`.

Return *the ordering of courses you should take to finish all courses*. If there are many valid answers, return **any** of them. If it is impossible to finish all courses, return **an empty array**.

**Example 1:**

```
Input: numCourses = 2, prerequisites = [[1,0]]
Output: [0,1]
Explanation: There are a total of 2 courses to take. To take course 1 you should have finished course 0. So the correct course order is [0,1].
```

**Example 2:**

```
Input: numCourses = 4, prerequisites = [[1,0],[2,0],[3,1],[3,2]]
Output: [0,2,1,3]
Explanation: There are a total of 4 courses to take. To take course 3 you should have finished both courses 1 and 2. Both courses 1 and 2 should be taken after you finished course 0.
So one correct course order is [0,1,2,3]. Another correct ordering is [0,2,1,3].
```

**Example 3:**

```
Input: numCourses = 1, prerequisites = []
Output: [0]
```

**Constraints:**

- `1 <= numCourses <= 2000`
- `0 <= prerequisites.length <= numCourses * (numCourses - 1)`
- `prerequisites[i].length == 2`
- `0 <= ai, bi < numCourses`
- `ai != bi`
- All the pairs `[ai, bi]` are **distinct**.

---

## Company Frequency

| Company    | Frequency        | Last Reported |
|------------|------------------|---------------|
| Amazon     | ★★★★★ Very High  | 2024          |
| Google     | ★★★★☆ High       | 2024          |
| Meta       | ★★★★☆ High       | 2024          |
| Microsoft  | ★★★☆☆ Medium     | 2024          |
| Bloomberg  | ★★★☆☆ Medium     | 2023          |
| Apple      | ★★☆☆☆ Low        | 2023          |
| Uber       | ★★☆☆☆ Low        | 2022          |

> ⚠️ Frequency data is crowd-sourced from LeetCode Discuss, Glassdoor, and
> community interview reports. Treat as a signal, not a guarantee.

---

## DSA Concepts Used

- **Topological Sort** — the whole problem *is* topological sorting: emit vertices of a DAG so every prerequisite precedes the course it unlocks; a cycle means no valid order exists → see [`/dsa/topological_sort.md`](/dsa/topological_sort.md)
- **Graph BFS / DFS** — the prerequisite pairs form a directed graph; both Kahn's BFS (in-degree peeling) and DFS postorder traverse it to produce the ordering and detect cycles → see [`/dsa/graph_bfs_dfs.md`](/dsa/graph_bfs_dfs.md)
- **Queue (FIFO)** — Kahn's algorithm feeds every zero-in-degree course through a queue; the dequeue order *is* the answer → see [`/dsa/queue_deque.md`](/dsa/queue_deque.md)

---

## Approaches Overview

| # | Approach | Time | Space | When to use |
|---|----------|------|-------|-------------|
| 1 | Brute Force (Repeated Elimination) | O(V · (V + E)) | O(V) | Baseline that mirrors the intuition ("keep taking whatever is available"); fine for tiny inputs, quadratic sweeps otherwise |
| 2 | DFS Postorder + Reverse | O(V + E) | O(V + E) | When you prefer recursion; postorder + reverse is the classic DFS topo-sort with 3-colour cycle detection |
| 3 | Kahn's Algorithm / BFS Topological Sort (Optimal) | O(V + E) | O(V + E) | The standard interview answer; iterative, no recursion depth risk, order falls out of the queue directly |

Here `V = numCourses` and `E = len(prerequisites)`.

---

## Approach 1 — Brute Force (Repeated Elimination)

### Intuition

This is the literal simulation of "keep enrolling in whatever you can". Sweep the course catalogue in index order; any course whose every prerequisite is already taken gets taken now and appended to the schedule. Because a course taken earlier in a sweep can unlock a later-indexed one in the *same* sweep, one pass can enrol several courses. Repeat sweeps until a full pass enrols nobody. If courses still remain untaken at that point, they mutually block each other — a cycle — and no ordering exists, so we return the empty slice.

### Algorithm

1. `taken[i]` tracks completion; `order` accumulates the schedule.
2. Sweep courses `0..n-1`. A course `c` is ready iff every pair `[c, b]` has `taken[b] == true`. Append each ready course to `order` and mark it taken.
3. Track a `progress` flag: if a whole sweep enrols nobody, stop.
4. Return `order` if it holds all `numCourses` courses, else return `[]`.

### Complexity

- **Time:** O(V · (V + E)) — up to V sweeps (each sweep may unlock only one course), and every sweep scans all V courses, checking each against all E prerequisite pairs.
- **Space:** O(V) — the `taken` boolean array plus the output `order`.

### Code

```go
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
```

### Dry Run

Example 2: `numCourses = 4`, `prerequisites = [[1,0],[2,0],[3,1],[3,2]]`.

Prerequisites read as: `1→0`, `2→0`, `3→1`, `3→2` (course on left needs course on right first).

| Sweep | c | taken[] before | ready? (missing prereq) | Action | order after |
|-------|---|----------------|-------------------------|--------|-------------|
| 1 | 0 | `[F F F F]` | yes (no prereqs) | take 0 | `[0]` |
| 1 | 1 | `[T F F F]` | yes (needs 0 ✔) | take 1 | `[0 1]` |
| 1 | 2 | `[T T F F]` | yes (needs 0 ✔) | take 2 | `[0 1 2]` |
| 1 | 3 | `[T T T F]` | yes (needs 1 ✔, 2 ✔) | take 3 | `[0 1 2 3]` |
| 2 | — | `[T T T T]` | nothing untaken | no progress → break | `[0 1 2 3]` |

`len(order) == 4 == numCourses`. Result: `[0 1 2 3]` ✔ (a valid topological order).

---

## Approach 2 — DFS Postorder + Reverse

### Intuition

Model the graph with edges `b → a` meaning "`b` unlocks `a`". Depth-first search finishes a node only after every node reachable from it is finished. So the DFS **finish order** (postorder) lists each course *after* all courses that depend on it — which is the exact reverse of a valid schedule. Emit postorder, then reverse it, and prerequisites land before their dependents. Cycles are caught with the 3-colour scheme: hitting a **grey** (on-path) node means the current path looped back on itself, so no ordering exists.

### Algorithm

1. Build adjacency list `adj[b] ← a` for each pair `[a, b]`.
2. Track `state[i] ∈ {white, grey, black}` and a `post` list.
3. `dfs(u)`: mark `u` grey; recurse into white neighbours; a grey neighbour is a cycle → abort. Then mark `u` black and append it to `post`.
4. Run `dfs` from every white node to cover disconnected components; any cycle returns `[]`.
5. Reverse `post` in place and return it.

### Complexity

- **Time:** O(V + E) — every node is entered once and every edge crossed once.
- **Space:** O(V + E) — adjacency list, state array, postorder list, and recursion stack.

### Code

```go
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
```

### Dry Run

Example 2: `numCourses = 4`, `prerequisites = [[1,0],[2,0],[3,1],[3,2]]`.

Adjacency (`adj[b]` = courses unlocked by `b`): `adj[0]=[1,2]`, `adj[1]=[3]`, `adj[2]=[3]`, `adj[3]=[]`.

DFS roots are tried in order `0,1,2,3`; only course `0` is white when we start (the rest get coloured during its recursion).

| Step | Call | state changes | post after | Notes |
|------|------|---------------|------------|-------|
| 1 | `dfs(0)` | 0→grey | `[]` | visit adj[0] = [1,2] |
| 2 | `dfs(1)` | 1→grey | `[]` | visit adj[1] = [3] |
| 3 | `dfs(3)` | 3→grey | `[]` | adj[3] empty |
| 4 | finish 3 | 3→black | `[3]` | append 3 |
| 5 | finish 1 | 1→black | `[3 1]` | back in dfs(1); append 1 |
| 6 | `dfs(2)` | 2→grey | `[3 1]` | back in dfs(0), visit 2; adj[2]=[3], 3 is black → skip |
| 7 | finish 2 | 2→black | `[3 1 2]` | append 2 |
| 8 | finish 0 | 0→black | `[3 1 2 0]` | append 0 |

Roots 1,2,3 are already black — no more work. Reverse `post = [3 1 2 0]` → `[0 2 1 3]`. Result: `[0 2 1 3]` ✔ (a valid topological order, matching the alternate ordering in Example 2).

---

## Approach 3 — Kahn's Algorithm / BFS Topological Sort (Optimal)

### Intuition

A course can be taken exactly when all its prerequisites are done — i.e. when its **in-degree** (count of unmet prerequisites) hits zero. Start by queuing every zero-in-degree course. Each time you dequeue a course, that's the next safe course to take, so append it to the schedule; then relax its outgoing edges, decrementing each neighbour's in-degree and queuing any that just reached zero. The dequeue order **is** a valid topological order. If a cycle exists, its members never reach in-degree 0, so fewer than `numCourses` courses dequeue — that's how we detect impossibility.

### Algorithm

1. Build `adj[b] ← a` and `indegree[a]++` for every pair `[a, b]`.
2. Enqueue all courses with `indegree == 0`.
3. Dequeue `u`, append it to `order`; for each neighbour `v`, do `indegree[v]--` and enqueue `v` if it reaches 0.
4. Return `order` if `len(order) == numCourses`, else `[]`.

### Complexity

- **Time:** O(V + E) — each course is enqueued and dequeued once; each edge is relaxed once.
- **Space:** O(V + E) — adjacency list, in-degree array, queue, and output.

### Code

```go
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
```

### Dry Run

Example 2: `numCourses = 4`, `prerequisites = [[1,0],[2,0],[3,1],[3,2]]`.

Build: `adj[0]=[1,2]`, `adj[1]=[3]`, `adj[2]=[3]`, `adj[3]=[]`; `indegree = [0, 1, 1, 2]` (course 0 free, 1 needs 0, 2 needs 0, 3 needs 1 and 2).

Initial queue = `[0]` (only course 0 has in-degree 0).

| Step | queue before | dequeue u | order after | relax neighbours (indegree after) | queue after |
|------|--------------|-----------|-------------|-----------------------------------|-------------|
| 1 | `[0]` | 0 | `[0]` | 1: 1→0 (enqueue), 2: 1→0 (enqueue) | `[1 2]` |
| 2 | `[1 2]` | 1 | `[0 1]` | 3: 2→1 | `[2]` |
| 3 | `[2]` | 2 | `[0 1 2]` | 3: 1→0 (enqueue) | `[3]` |
| 4 | `[3]` | 3 | `[0 1 2 3]` | (no neighbours) | `[]` |

Queue empty; `len(order) == 4 == numCourses`. Result: `[0 1 2 3]` ✔.

---

## Key Takeaways

- **"Find an order that respects dependencies" = topological sort of a DAG.** The moment prerequisites appear, reach for topo-sort. Course Schedule II is the canonical "return the order" version; #207 is the "does an order exist?" version — same graph, one returns the list, the other a bool.
- **Kahn's BFS gives the answer for free.** The order in which zero-in-degree nodes leave the queue is already a valid topological order — no post-processing. Cycle detection is "did fewer than V nodes get processed?"
- **DFS postorder needs a reverse.** DFS finishes prerequisites *last* along `b→a` edges, so the finish order is the reverse schedule. Append on finish, then reverse. Use 3-colour (white/grey/black) states so a back edge (grey) flags a cycle mid-traversal.
- **Edge direction is a choice — commit to it.** Here `b → a` ("b unlocks a") makes in-degree = unmet prerequisites, which is what both Kahn and the DFS reversal rely on. Flipping the direction flips the whole algorithm; write the mapping down before coding.
- **Any valid order is accepted.** Different approaches (and different queue/adjacency orders) yield different-but-correct schedules, e.g. `[0 1 2 3]` vs `[0 2 1 3]`. Don't over-index on matching one specific output.

---

## Related Problems

- LeetCode #207 — Course Schedule (same graph; return whether an ordering exists instead of the ordering itself)
- LeetCode #269 — Alien Dictionary (topological sort over letter-ordering constraints)
- LeetCode #310 — Minimum Height Trees (peel-leaves BFS, a Kahn's-style in-degree variant)
- LeetCode #444 — Sequence Reconstruction (verify a unique topological order via Kahn's)
- LeetCode #1136 — Parallel Courses (topological sort measuring the number of semesters/levels)
