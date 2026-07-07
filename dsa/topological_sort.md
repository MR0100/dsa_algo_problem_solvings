# Topological Sort

> **Category:** Graphs · Directed Acyclic Graphs (DAGs)
> **Prerequisites:** [`graph representation`](/dsa/graph.md) (adjacency lists), BFS, DFS, [`queue`](/dsa/queue_deque.md), [`stack`](/dsa/stack.md)

---

## What it is

A **topological sort** (topo sort / topological ordering) of a **directed graph** is a
linear ordering of its vertices such that for every directed edge `u → v`,
vertex `u` appears **before** vertex `v` in the ordering.

Think of it as flattening a dependency graph into a valid execution order:
if task `u` must happen before task `v`, the ordering respects that.

Key facts:

- A topological ordering exists **if and only if** the graph is a **DAG**
  (Directed Acyclic Graph). A cycle makes ordering impossible — `a → b → a`
  means `a` must come both before and after `b`.
- A DAG can have **many** valid topological orderings. Any of them is a
  correct answer unless the problem asks for a specific one (e.g.
  lexicographically smallest → use a min-heap instead of a plain queue).
- Topological sort doubles as a **cycle detector**: if the algorithm cannot
  place every vertex, the graph has a cycle. Many "is this schedule possible?"
  problems are exactly this check.
- On a topologically sorted DAG you can run **DP over the graph** (longest
  path, counting paths, propagating values) in a single linear pass, because
  every node is processed only after all of its prerequisites.

---

## How to recognise it — signals in the problem statement

Reach for topological sort when you see:

| Signal | Example phrasing |
|--------|------------------|
| **Prerequisites / dependencies** | "course `b` must be taken before course `a`", "build `x` depends on `y`" |
| **Ordering tasks under constraints** | "return any valid order to finish all tasks" |
| **"Is it possible to finish?"** | "return `true` if you can finish all courses" → cycle detection on a directed graph |
| **Deriving an unknown order from pairwise clues** | "given sorted words from an alien dictionary, recover the alphabet" |
| **Compilation / build / pipeline ordering** | "in what order should packages be installed?" |
| **Layer-by-layer peeling of a graph** | "repeatedly remove nodes with no outgoing/incoming edges" (e.g. finding tree centres, safe states) |
| **Longest chain through a DAG** | "longest path where each step must follow the previous" — topo order + DP |
| **Counting semesters / rounds** | "minimum number of semesters to take all courses" — BFS topo sort, counting levels |

Rules of thumb:

- **Directed** relationships + a need for a **global order** ⇒ topological sort.
- If the graph is undirected, topo sort does not apply (use union-find, BFS, DFS).
- If the answer only needs "possible or not", you still run a full topo sort —
  the possibility check *is* whether all `V` vertices get ordered.

---

## Template 1 — Kahn's algorithm (BFS, in-degree)

The iterative, queue-based approach. Usually the one to reach for in
interviews: no recursion, cycle detection falls out for free, and it
naturally supports "process in rounds/levels" variants.

**Pseudocode:**

```
1. Build adjacency list adj[u] = list of v for every edge u → v.
2. Compute inDegree[v] = number of edges pointing INTO v.
3. Push every vertex with inDegree == 0 into a queue.        // no prerequisites
4. While queue not empty:
     u = pop front
     append u to order
     for each neighbour v of u:
         inDegree[v]--                                        // u's prerequisite satisfied
         if inDegree[v] == 0: push v                          // all prereqs of v done
5. If len(order) < V ⇒ cycle (some vertices never reached 0). Otherwise order is the answer.
```

**Go:**

```go
// topoSortKahn returns a topological ordering of a directed graph with n
// vertices (0..n-1) and edge list edges, where edges[i] = [u, v] means u → v
// (u must come before v). Returns nil if the graph contains a cycle.
//
// Time:  O(V + E) — each vertex enqueued/dequeued once, each edge relaxed once.
// Space: O(V + E) — adjacency list, in-degree array, queue, output.
func topoSortKahn(n int, edges [][]int) []int {
    adj := make([][]int, n)   // adjacency list: adj[u] = all v with edge u → v
    inDegree := make([]int, n) // inDegree[v] = number of unmet prerequisites of v

    for _, e := range edges {
        u, v := e[0], e[1]
        adj[u] = append(adj[u], v) // record dependency u → v
        inDegree[v]++              // v gains one more prerequisite
    }

    // Seed the queue with every vertex that has no prerequisites.
    queue := []int{}
    for v := 0; v < n; v++ {
        if inDegree[v] == 0 {
            queue = append(queue, v)
        }
    }

    order := make([]int, 0, n) // the topological ordering we build
    for len(queue) > 0 {
        u := queue[0]     // pop from the front (FIFO)
        queue = queue[1:]
        order = append(order, u) // u has no remaining prerequisites → safe to place

        for _, v := range adj[u] {
            inDegree[v]-- // prerequisite u is now satisfied for v
            if inDegree[v] == 0 {
                queue = append(queue, v) // all of v's prerequisites done
            }
        }
    }

    if len(order) < n { // some vertices never reached in-degree 0
        return nil // cycle detected — no valid ordering exists
    }
    return order
}
```

**Variants of Template 1:**

- **Level-by-level (minimum rounds/semesters):** process the queue one whole
  level at a time (like BFS level-order); the number of levels is the minimum
  number of parallel rounds.
- **Lexicographically smallest order:** replace the FIFO queue with a
  **min-heap**; always pop the smallest available vertex. Cost becomes
  `O((V + E) log V)`.

---

## Template 2 — DFS with post-order (colours / visited states)

Run DFS; a vertex is *finished* only after all vertices it points to are
finished. Appending each vertex on finish and **reversing** gives a
topological order. Cycle detection needs a 3-state marking, because a plain
boolean `visited` cannot distinguish "currently on the recursion stack"
(back edge ⇒ cycle) from "fully processed earlier" (harmless cross edge).

**Pseudocode:**

```
state[v] ∈ {white=unvisited, gray=in progress, black=done}

dfs(u):
    state[u] = gray                       // u is on the current recursion path
    for each v in adj[u]:
        if state[v] == gray: CYCLE        // back edge to an ancestor
        if state[v] == white: dfs(v)
    state[u] = black                      // everything reachable from u is done
    append u to order                     // post-order position

run dfs from every white vertex; reverse(order) is the topological sort.
```

**Go:**

```go
// topoSortDFS returns a topological ordering using DFS post-order, or nil if
// the graph contains a cycle.
//
// Time:  O(V + E) — every vertex and edge visited once.
// Space: O(V + E) — adjacency list, state array, recursion stack (up to O(V)).
func topoSortDFS(n int, edges [][]int) []int {
    adj := make([][]int, n)
    for _, e := range edges {
        adj[e[0]] = append(adj[e[0]], e[1]) // edge u → v
    }

    const (
        white = 0 // unvisited
        gray  = 1 // on the current DFS path (in progress)
        black = 2 // fully processed
    )
    state := make([]int, n)
    order := make([]int, 0, n) // vertices in post-order (finish time)
    hasCycle := false

    var dfs func(u int)
    dfs = func(u int) {
        state[u] = gray // mark: u is on the active recursion path
        for _, v := range adj[u] {
            if state[v] == gray { // back edge → v is an ancestor of u → cycle
                hasCycle = true
                return
            }
            if state[v] == white {
                dfs(v)
                if hasCycle {
                    return // abort early once a cycle is found
                }
            }
            // state[v] == black: already finished, safe to ignore
        }
        state[u] = black          // all descendants of u are placed
        order = append(order, u)  // u finishes now → post-order append
    }

    for v := 0; v < n; v++ {
        if state[v] == white {
            dfs(v)
            if hasCycle {
                return nil
            }
        }
    }

    // Post-order lists dependencies before dependents in REVERSE, so flip it.
    for i, j := 0, len(order)-1; i < j; i, j = i+1, j-1 {
        order[i], order[j] = order[j], order[i]
    }
    return order
}
```

**When to prefer which:**

| | Kahn (BFS) | DFS post-order |
|---|---|---|
| Cycle detection | free (`len(order) < n`) | needs 3-colour states |
| "Minimum rounds / levels" | natural (level BFS) | awkward |
| Lexicographic order | swap queue → min-heap | awkward |
| Combine with DP on DAG | works | very natural (memoised DFS) |
| Recursion depth risk | none (iterative) | O(V) stack — deep chains can overflow |

---

## Worked example — full Kahn trace

Course-schedule style input: `n = 6` vertices, edges (u → v means "u before v"):

```
5 → 2,  5 → 0,  4 → 0,  4 → 1,  2 → 3,  3 → 1
```

**Setup:**

| Vertex | adj (points to) | inDegree |
|--------|-----------------|----------|
| 0 | — | 2 (from 5, 4) |
| 1 | — | 2 (from 4, 3) |
| 2 | 3 | 1 (from 5) |
| 3 | 1 | 1 (from 2) |
| 4 | 0, 1 | 0 |
| 5 | 2, 0 | 0 |

Initial queue: `[4, 5]` (in-degree 0). `order = []`.

**Trace:**

| Step | Pop | Action on neighbours | inDegree after | Queue after | order |
|------|-----|----------------------|----------------|-------------|-------|
| 1 | 4 | 0: 2→1; 1: 2→1 | 0:1, 1:1, 2:1, 3:1 | [5] | [4] |
| 2 | 5 | 2: 1→**0** push; 0: 1→**0** push | 0:0, 1:1, 2:0, 3:1 | [2, 0] | [4, 5] |
| 3 | 2 | 3: 1→**0** push | 1:1, 3:0 | [0, 3] | [4, 5, 2] |
| 4 | 0 | (no neighbours) | — | [3] | [4, 5, 2, 0] |
| 5 | 3 | 1: 1→**0** push | 1:0 | [1] | [4, 5, 2, 0, 3] |
| 6 | 1 | (no neighbours) | — | [] | [4, 5, 2, 0, 3, 1] |

`len(order) == 6 == n` → no cycle. Answer: `[4, 5, 2, 0, 3, 1]`.
Check any edge, e.g. `3 → 1`: 3 appears at index 4, 1 at index 5. ✓
(Other orders like `[5, 4, 2, 0, 3, 1]` would be equally valid.)

**Cycle case:** add edge `1 → 5`. Now 5 has in-degree 1, the initial queue is
`[4]`, and after processing 4 nothing else ever reaches in-degree 0 (the cycle
`5 → 2 → 3 → 1 → 5` keeps every member at in-degree ≥ 1). `order = [4]`,
`len(order) = 1 < 6` → cycle reported.

---

## Common pitfalls and how to avoid them

1. **Edge direction flipped.** LeetCode prerequisite pairs are often given as
   `[a, b]` meaning "take `b` before `a`", i.e. edge `b → a` — the reverse of
   the pair's reading order. Getting this backwards produces the reverse
   ordering (or wrong cycle results on non-symmetric inputs). *Fix:* write a
   one-line comment stating your edge convention before building `adj`, and
   double-check against Example 1.

2. **Forgetting the cycle check.** Kahn's loop terminates quietly on a cyclic
   graph — it just outputs fewer than `V` vertices. Always compare
   `len(order)` with `V` at the end; returning a partial order is a silent bug.

3. **Using a boolean `visited` in the DFS version.** You must distinguish
   "on the current path" (gray → back edge → cycle) from "finished earlier"
   (black → fine). Two states can't express this; a diamond `a→b, a→c, b→d, c→d`
   gets falsely flagged as a cycle, or a real cycle gets missed. Use 3 states.

4. **Forgetting to reverse the DFS post-order.** Post-order finishes sinks
   first, so the raw list is the topological order backwards. Either reverse
   at the end or prepend (reversing a slice is cheaper in Go).

5. **Seeding the queue with only one zero-in-degree vertex.** A DAG can have
   many sources, and the graph may be disconnected. Scan **all** `n` vertices
   for `inDegree == 0`, and in the DFS version launch `dfs` from every white
   vertex — never just from vertex 0.

6. **Counting in-degree from the adjacency list of the same vertex.**
   `inDegree[v]` counts edges *into* `v`; it is incremented when scanning
   edges, not derived from `len(adj[v])` (that's the *out*-degree).

7. **Duplicate edges inflating in-degree.** If the input can contain the same
   prerequisite pair twice, either deduplicate edges or make sure you
   decrement once per stored edge (consistent build + relax is fine;
   dedup-on-build but not on-relax is not).

8. **Assuming the answer is unique.** Judges accept any valid ordering, but if
   the problem demands a specific one (smallest lexicographic, must match a
   given sequence), a plain queue is not enough — use a min-heap, or verify
   the given sequence is *the only* order (queue size must stay ≤ 1, e.g.
   LeetCode #444 Sequence Reconstruction).

9. **Recursion depth on huge inputs.** Go goroutine stacks grow dynamically,
   so DFS to depth 10⁵ is usually fine — but on other platforms/languages a
   long chain graph overflows the stack. When in doubt, use Kahn's (fully
   iterative).

10. **Applying topo sort to an undirected graph.** In-degree/ordering only
    makes sense for directed edges. Undirected "dependency" problems are
    usually connectivity (union-find) or leaf-peeling (e.g. Minimum Height
    Trees uses a Kahn-*like* peel of degree-1 nodes, but it is not a
    topological sort).

---

## Problems in this repo

*No problems in this repo (currently 0001–0130) use topological sort yet.*
Classic LeetCode problems that will link here once added:

- LeetCode #207 — Course Schedule (cycle detection via topo sort)
- LeetCode #210 — Course Schedule II (return the ordering)
- LeetCode #269 — Alien Dictionary (build graph from pairwise clues, then topo sort)
- LeetCode #310 — Minimum Height Trees (Kahn-style leaf peeling on an undirected tree)
- LeetCode #329 — Longest Increasing Path in a Matrix (DP over an implicit DAG)

<!-- Later pass: add relative links like [0207_course_schedule](../0207_course_schedule/README.md) as those folders land. -->
