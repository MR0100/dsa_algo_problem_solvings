# Dijkstra's Algorithm (Single-Source Shortest Path)

> **Solves:** shortest path from one source to every node in a graph with **non-negative** edge weights.
> **Engine:** a min-heap (priority queue) that always expands the *closest un-settled* node next.
> **Complexity (binary heap):** O((V + E) log V) time, O(V) space.

---

## What it is

Dijkstra's algorithm computes the shortest distance from a single **source** vertex to all other vertices of a weighted graph, provided **every edge weight is ≥ 0**. It is a **greedy** algorithm built on one invariant:

> The first time a node is popped from the priority queue, its recorded distance is already the true shortest distance — it is *settled* and never needs revisiting.

That invariant is exactly why non-negative weights are required (see below). You maintain a `dist[]` array (best-known distance to each node, initialised to +∞ except `dist[source] = 0`) and a min-heap of `(distance, node)` pairs. Repeatedly pop the smallest, and for each outgoing edge `u → v` of weight `w`, **relax** it:

```
if dist[u] + w < dist[v] {
    dist[v] = dist[u] + w
    push (dist[v], v)
}
```

### Where it sits among shortest-path tools

| Algorithm | Edge weights | Time | Use when |
|-----------|--------------|------|----------|
| **BFS** | all equal (unweighted) | O(V + E) | every edge costs the same — a queue *is* the priority queue |
| **Dijkstra** | non-negative | O((V+E) log V) | weighted, no negative edges — the common case |
| **Bellman-Ford** | any (incl. negative) | O(V·E) | negative edges present; also detects negative cycles |
| **0-1 BFS** | only 0 or 1 | O(V + E) | weights are just 0/1 — a deque replaces the heap |
| **Floyd-Warshall** | any (no neg cycle) | O(V³) | *all-pairs* shortest paths on a small dense graph |

**Mental model:** think of distances as water flooding outward from the source. The heap always releases water into the nearest dry node first; once a node is wet (popped), its fill-level is final.

### What Dijkstra cannot do

- **Negative edge weights** break the settle-on-first-pop invariant — a cheaper path may arrive *after* the node was already finalised. Use Bellman-Ford (or Johnson's for all-pairs).
- **Longest path / negative cycles** — outside its model entirely.

---

## When to recognise it — signals in the problem statement

| Signal | Why Dijkstra |
|--------|--------------|
| "shortest / cheapest / minimum-cost path" on a graph with **weighted, non-negative** edges | textbook trigger — BFS would be wrong because edges differ in cost |
| "minimum time for a signal / package / robot to reach all nodes" | single-source, weights = latency/cost → Dijkstra; answer is `max(dist)` for "reach *all*" |
| a grid where **moving costs vary** (roll distance, elevation, effort) | flatten `(r,c)` → node; edge weight = the move's cost |
| "shortest path, and among ties the **lexicographically smallest** route" | Dijkstra with an augmented key `(distance, path-string)` — the heap breaks ties for free |
| "minimise the **maximum** edge / effort along a path" (bottleneck) | the *minimax* flavour: replace `dist[u]+w` with `max(dist[u], w)` — same heap machinery |
| "water pours over the **lowest surrounding wall**", "path of least resistance" | minimal-enclosing / bottleneck-shortest — a boundary-seeded Dijkstra |
| "cheapest flight with **at most K stops**" | Dijkstra with an extra state dimension (`stops`), or Bellman-Ford bounded to K+1 rounds |

**When *not* to use it:** unweighted graph → plain BFS is simpler and faster. Any negative edge → Bellman-Ford. All-pairs on a tiny graph → Floyd-Warshall.

---

## General template (Go, `container/heap`)

Go has no built-in priority queue; you implement `heap.Interface` (five methods) over a slice. This is the exact idiom used by the Dijkstra solutions in this repo.

```go
import "container/heap"

// item is one entry in the priority queue: "node is reachable in dist".
type item struct {
    node int
    dist int
}

// pq is a min-heap of items ordered by dist (smallest dist = highest priority).
type pq []item

func (p pq) Len() int            { return len(p) }
func (p pq) Less(i, j int) bool  { return p[i].dist < p[j].dist } // min-heap on dist
func (p pq) Swap(i, j int)       { p[i], p[j] = p[j], p[i] }
func (p *pq) Push(x interface{}) { *p = append(*p, x.(item)) }
func (p *pq) Pop() interface{} {
    old := *p
    n := len(old)
    it := old[n-1]
    *p = old[:n-1]
    return it
}

// dijkstra returns the shortest distance from src to every node.
// graph is an adjacency list: graph[u] = list of (neighbour, weight) edges.
// All weights must be >= 0.
func dijkstra(n, src int, graph [][]item) []int {
    dist := make([]int, n)
    for i := range dist {
        dist[i] = 1 << 62 // +infinity sentinel
    }
    dist[src] = 0

    h := &pq{{node: src, dist: 0}}
    heap.Init(h)

    for h.Len() > 0 {
        cur := heap.Pop(h).(item) // the closest node not yet settled

        // Lazy deletion: if we already found a better distance for this node,
        // this popped entry is stale — skip it. (We never decrease-key; we
        // just push duplicates and discard the outdated pops.)
        if cur.dist > dist[cur.node] {
            continue
        }

        // cur.node is now SETTLED: dist[cur.node] is final.
        for _, e := range graph[cur.node] {
            nd := cur.dist + e.dist // candidate distance to neighbour e.node
            if nd < dist[e.node] {  // relaxation: found a shorter route
                dist[e.node] = nd
                heap.Push(h, item{node: e.node, dist: nd})
            }
        }
    }
    return dist
}
```

### The `visited` / settled variant

Instead of the `cur.dist > dist[cur.node]` staleness check you can keep an explicit `settled []bool` and `continue` when `settled[cur.node]` is already true, setting it to true right after the pop. Both are correct; the lazy-deletion form above needs no extra array.

```go
if settled[cur.node] { continue }
settled[cur.node] = true
```

### Lexicographic / tie-break variant (used by #499 The Maze III)

When two paths have the **same** distance and the problem wants the lexicographically smallest route, augment the heap key so ties fall back to the path string. `Less` compares `dist` first, then `path`:

```go
type state struct {
    dist int
    path string // instructions accumulated so far, e.g. "dl"
    r, c int
}

func (p pq) Less(i, j int) bool {
    if p[i].dist != p[j].dist {
        return p[i].dist < p[j].dist // primary key: shortest distance
    }
    return p[i].path < p[j].path     // tie-break: lexicographically smaller route
}
```

Because the heap yields the `(dist, path)`-minimum, the first time you settle the target the accumulated `path` is guaranteed both shortest and lexicographically minimal. (Iterate the four directions in sorted order `d < l < r < u` so equal-length alternatives are generated deterministically.)

### Bottleneck / minimal-enclosing variant (used by #407 Trapping Rain Water II)

Some problems don't sum edge weights — they want to **minimise the maximum** value crossed along the path (a *minimax* / bottleneck shortest path). Only the relaxation step changes: `+` becomes `max`.

```go
// Instead of nd := cur.dist + e.dist  (sum of costs)
nd := cur.dist
if e.dist > nd {
    nd = e.dist // path cost = the highest wall we had to clear
}
if nd < dist[e.node] { ... } // still "keep the cheapest such maximum"
```

For Trapping Rain Water II this becomes a **boundary-seeded** Dijkstra: push every border cell into the min-heap keyed on its height, then always pour over the *lowest wall on the frontier*. When you expand into a lower inner cell, water rises to the current wall height, and the trapped volume is `wall − cellHeight`; the effective height propagated inward is `max(wall, cellHeight)`. The min-heap guaranteeing "lowest wall first" is exactly Dijkstra's greedy choice.

---

## Worked example — step-by-step trace

Graph with 5 nodes `0..4`, directed weighted edges, source `0`:

```
0 → 1 (10)   0 → 3 (5)
1 → 2 (1)    1 → 3 (2)
3 → 1 (3)    3 → 2 (9)   3 → 4 (2)
2 → 4 (4)    4 → 2 (6)
```

Initialise: `dist = [0, ∞, ∞, ∞, ∞]`, heap = `{(0,0)}`.

| Step | Pop `(d,u)` | Settled? | Relaxations | `dist` after | Heap after (by dist) |
|------|-------------|----------|-------------|--------------|----------------------|
| 1 | (0, **0**) | yes | 1: 0+10=10 ✓ · 3: 0+5=5 ✓ | `[0,10,∞,5,∞]` | (5,3),(10,1) |
| 2 | (5, **3**) | yes | 1: 5+3=8<10 ✓ · 2: 5+9=14 ✓ · 4: 5+2=7 ✓ | `[0,8,14,5,7]` | (7,4),(8,1),(10,1*),(14,2) |
| 3 | (7, **4**) | yes | 2: 7+6=13<14 ✓ | `[0,8,13,5,7]` | (8,1),(10,1*),(13,2),(14,2*) |
| 4 | (8, **1**) | yes | 2: 8+1=9<13 ✓ · 3: 8+2=10, not < 5 ✗ | `[0,8,9,5,7]` | (9,2),(10,1*),(13,2*),(14,2*) |
| 5 | (9, **2**) | yes | 4: 9+4=13, not < 7 ✗ | `[0,8,9,5,7]` | (10,1*),(13,2*),(14,2*) |
| 6 | (10, **1***) | **stale** (10 > dist[1]=8) → skip | — | unchanged | (13,2*),(14,2*) |
| 7 | (13, **2***) | **stale** (13 > dist[2]=9) → skip | — | unchanged | (14,2*) |
| 8 | (14, **2***) | **stale** → skip | — | unchanged | empty |

Entries marked `*` are the outdated duplicates we pushed before finding a shorter distance; each is discarded in O(1) at pop time by the staleness check — this is the "lazy decrease-key" that keeps the code simple.

**Final shortest distances from 0:** `dist = [0, 8, 9, 5, 7]`.
Note node 1's distance improved from 10 (direct) to 8 (via 3) *before* it was ever settled — that is the whole point of relaxation, and it only works because no edge is negative.

---

## Complexity

Let V = vertices, E = edges.

| Operation | Cost | Why |
|-----------|------|-----|
| Each edge relaxed | O(1), pushes ≤1 heap entry | scanned once per settle of its tail |
| Heap push / pop | O(log V) | binary heap; at most O(E) entries live |
| **Total time** | **O((V + E) log V)** | ≈ O(E log V) on a connected graph |
| Space | O(V + E) | adjacency list + `dist[]` + heap (≤ E entries) |

Notes:
- With a **Fibonacci heap** the bound drops to O(E + V log V) — theoretically better, but rarely worth the constant factors and complexity on interview-scale inputs.
- For a **dense** graph (E ≈ V²) a plain O(V²) array-scan Dijkstra (no heap) can be competitive.
- The lazy-deletion approach may hold up to O(E) heap entries; that is fine and standard.

---

## Common pitfalls

1. **Using Dijkstra with negative edges.** The settle-on-first-pop invariant fails — a node can be finalised before a cheaper (negative-edge) path reaches it. Symptom: silently wrong distances. Fix: detect negative weights and switch to Bellman-Ford.

2. **Forgetting the staleness / visited check.** Without `if cur.dist > dist[cur.node] { continue }` (or a `settled[]` guard) you re-expand nodes through outdated entries. Usually still correct but wastefully slow; with certain augmented states it can produce wrong results.

3. **Confusing "popped" with "reachable".** A node's distance is only *final* when it is **popped**, not when it is first pushed. Reading `dist[target]` and returning early *before* target is popped can give a non-minimal value if you also short-circuit relaxations.

4. **Trying to `decrease-key` in Go.** `container/heap` has no decrease-key. Don't hunt an entry to mutate it — just push a new `(smaller-dist, node)` and let the stale one be skipped on pop.

5. **Wrong `Less` direction.** `container/heap` is a **min-heap** by default *only if* `Less` returns "i has higher priority than j" = "i.dist < j.dist". Flip it and you build a max-heap (longest paths). Double-check the comparator.

6. **Integer overflow on the ∞ sentinel.** If `dist[u]` is `math.MaxInt` and you compute `dist[u] + w`, it overflows and may become negative, corrupting relaxations. Use a sentinel like `1 << 62` (leaves head-room) and/or skip relaxation when `dist[u]` is still ∞.

7. **Bottleneck vs additive mix-up.** For minimax problems the relaxation must be `max(dist[u], w)`, not `dist[u] + w`. Using the additive form answers a different question (total cost, not worst edge).

8. **Lexicographic ties done wrong.** Comparing only distance and then hoping the traversal order is "naturally" sorted is fragile. Put the tie-break *inside* the heap key (`(dist, path)`) and generate neighbours in sorted order so the first settle of the target is provably minimal.

9. **1-indexed node labels.** LeetCode graphs (e.g. network-delay style) are often labeled `1..n`. Allocate `n+1` slots or subtract 1 consistently.

10. **Grid flatten formula.** When mapping a cell to a node use `r*cols + c` (not `r*rows + c`), and only add edges to in-bounds, non-wall neighbours.

---

## Problems in this repo that use Dijkstra

- [0499 — The Maze III](/0499_the_maze_iii/README.md) — Dijkstra on the *roll-graph*: each "stop against a wall" is a node, edge weight = empty cells rolled; the heap key is `(distance, path-string)` so the first settle of the hole yields the shortest **and lexicographically minimum** instruction string.
- [0407 — Trapping Rain Water II](/0407_trapping_rain_water_ii/README.md) — boundary-seeded min-heap Dijkstra of the **bottleneck / minimal-enclosing** flavour: always pour over the lowest wall on the frontier; trapped water at an inner cell = `frontierWall − cellHeight`, and the propagated height is `max(frontierWall, cellHeight)`.

Closely related grid shortest-path problems also in the repo (BFS / Dijkstra-adjacent):

- [0286 — Walls and Gates](/0286_walls_and_gates/README.md) — multi-source BFS (unit weights): the unweighted cousin where a plain queue replaces the heap.
- [0490 — The Maze](/0490_the_maze/README.md) — reachability on the same roll-graph as #499, but unweighted (BFS/DFS suffices because we only ask *can* it stop at the target).

### Related classics to know (not yet in repo)

- LeetCode #743 — Network Delay Time (the canonical single-source Dijkstra; answer is `max(dist)`)
- LeetCode #787 — Cheapest Flights Within K Stops (Dijkstra with an extra `stops` state dimension, or bounded Bellman-Ford)
- LeetCode #505 — The Maze II (weighted roll-graph like #499, but return the shortest *distance* — no lexicographic tie-break)
- LeetCode #1631 — Path With Minimum Effort (pure **bottleneck** Dijkstra: minimise the maximum absolute height difference along the path)
- LeetCode #1514 — Path with Maximum Probability (Dijkstra on a **max-heap** with multiplicative edge relaxation)
- LeetCode #778 — Swim in Rising Water (bottleneck Dijkstra / minimax on a grid)
