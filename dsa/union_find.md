# Union Find (Disjoint Set Union / DSU)

> **Also known as:** Disjoint Set Union (DSU), Merge-Find Set
> **Core operations:** `Find(x)` — which group is `x` in? · `Union(x, y)` — merge the groups of `x` and `y`.
> **Complexity with both optimizations:** near-O(1) amortized per operation — formally O(α(n)), where α is the inverse Ackermann function (α(n) ≤ 4 for any realistic n).

---

## What it is

Union Find is a data structure that maintains a collection of **disjoint (non-overlapping) sets** and supports two operations extremely fast:

1. **Find(x)** — return a canonical *representative* (the "root") of the set containing `x`. Two elements are in the same set iff their roots are equal.
2. **Union(x, y)** — merge the two sets containing `x` and `y` into one.

Internally, each set is stored as a **tree**: every element points to a parent, and the root points to itself. `Find` walks up the parent chain to the root; `Union` attaches one root under the other.

Two optimizations make this blazingly fast:

- **Path compression** — during `Find`, re-point every node visited directly at the root, flattening the tree for all future queries.
- **Union by rank (or by size)** — during `Union`, always attach the *shorter/smaller* tree under the *taller/larger* one, keeping trees shallow.

Either optimization alone gives O(log n) amortized; **both together give O(α(n))** — effectively constant.

### Mental model

Think of it as a **dynamic connectivity** tracker. You start with `n` isolated islands. Each `Union` builds a bridge. At any moment you can ask, in constant time, "are these two islands connected (possibly through many bridges)?" and "how many separate island-groups remain?"

What it **cannot** do (know these limits for interviews):

- **No un-union / deletion** — merges are irreversible. If a problem removes edges over time, the standard trick is to process operations **in reverse** (deletions become insertions).
- **No path retrieval** — it answers *whether* two nodes are connected, not *how*. For actual paths, use BFS/DFS.
- **Undirected connectivity only** — it models symmetric relations. Directed reachability needs other tools (topological sort, SCC/Tarjan).

---

## How to recognise it — signals in the problem statement

Reach for Union Find when you see:

| Signal | Example phrasing |
|--------|------------------|
| **Counting connected components** | "How many provinces / islands / networks / friend circles are there?" |
| **Grouping by an equivalence relation** | "Merge accounts that share an email", "words that are synonyms", "similar string groups" — any relation that is reflexive, symmetric, and transitive |
| **Dynamic / incremental connectivity** | "After each edge is added, report the number of components" — edges arrive over time and you must answer queries between insertions (BFS/DFS would need a full re-scan each time) |
| **Cycle detection in an undirected graph** | "Find the redundant edge", "can this be a valid tree?" — a `Union` of two nodes that already share a root means the new edge closes a cycle |
| **"Are X and Y connected?" queries** | Many pairwise connectivity queries after (or interleaved with) merges |
| **Kruskal's Minimum Spanning Tree** | Sort edges by weight, add each edge unless it creates a cycle — the cycle check *is* Union Find |
| **Grid problems as an alternative to DFS/BFS** | Number of islands, surrounded regions — flatten cell `(r, c)` to index `r*cols + c` and union adjacent same-value cells. Especially good when cells are **added online** (e.g. "Number of Islands II") |
| **Offline query tricks** | Sort queries and edges by threshold, add edges as the threshold grows, answer each query with `Find` (e.g. "paths with edge weights < limit") |
| **Percolation / union-with-virtual-node** | Introduce a virtual "super node" and union all boundary cells to it (e.g. all border-connected `O`s in Surrounded Regions, top/bottom rows in water-flow problems) |

**When *not* to use it:** single static graph, single connectivity question → plain DFS/BFS is simpler and strictly O(V+E). Need shortest paths or the path itself → BFS/Dijkstra. Directed edges → not DSU territory.

---

## General template (Go)

Idiomatic, interview-ready struct with path compression + union by rank, plus the two bookkeeping fields you almost always want (`count` of components, sometimes `size` of each set):

```go
// DSU is a disjoint-set-union over elements 0..n-1.
type DSU struct {
    parent []int // parent[i] = parent of i; root iff parent[i] == i
    rank   []int // rank[i] = upper bound on tree height rooted at i
    size   []int // size[root] = number of elements in that set
    count  int   // number of disjoint sets currently
}

// NewDSU creates n singleton sets: {0}, {1}, ..., {n-1}.
func NewDSU(n int) *DSU {
    d := &DSU{
        parent: make([]int, n),
        rank:   make([]int, n),
        size:   make([]int, n),
        count:  n,
    }
    for i := 0; i < n; i++ {
        d.parent[i] = i // each element starts as its own root
        d.size[i] = 1   // each set has one element
    }
    return d
}

// Find returns the root of x's set, compressing the path as it goes.
func (d *DSU) Find(x int) int {
    // Walk up until we hit the root (a node that is its own parent).
    if d.parent[x] != x {
        // Path compression: point x directly at the ultimate root.
        // The recursion returns the root, and every node on the way
        // up gets re-parented to it — flattening the whole chain.
        d.parent[x] = d.Find(d.parent[x])
    }
    return d.parent[x]
}

// Union merges the sets containing x and y.
// Returns false if they were already in the same set (useful for cycle detection).
func (d *DSU) Union(x, y int) bool {
    rx, ry := d.Find(x), d.Find(y)
    if rx == ry {
        return false // already connected — adding this edge would form a cycle
    }
    // Union by rank: attach the shorter tree under the taller one
    // so the resulting tree's height doesn't grow unnecessarily.
    if d.rank[rx] < d.rank[ry] {
        rx, ry = ry, rx // ensure rx is the taller (or equal) tree
    }
    d.parent[ry] = rx        // shorter root ry now hangs under rx
    d.size[rx] += d.size[ry] // rx's set absorbed ry's members
    if d.rank[rx] == d.rank[ry] {
        d.rank[rx]++ // equal heights: the merged tree grows by one level
    }
    d.count-- // two sets became one
    return true
}

// Connected reports whether x and y are in the same set.
func (d *DSU) Connected(x, y int) bool {
    return d.Find(x) == d.Find(y)
}
```

### Commented pseudocode (the skeleton to memorise)

```text
init(n):
    parent[i] = i  for all i        # everyone is their own root
    count = n                        # n singleton components

find(x):
    while parent[x] != x:            # climb to the root
        parent[x] = parent[parent[x]]  # (path halving — iterative compression)
        x = parent[x]
    return x

union(x, y):
    rx, ry = find(x), find(y)
    if rx == ry: return false        # same set → edge would close a cycle
    attach smaller-rank root under larger-rank root
    count -= 1
    return true
```

### Common variants

**Iterative Find with path halving** (avoids recursion; every node ends up pointing at its grandparent — same asymptotics):

```go
func (d *DSU) Find(x int) int {
    for d.parent[x] != x {
        d.parent[x] = d.parent[d.parent[x]] // skip a level: point at grandparent
        x = d.parent[x]
    }
    return x
}
```

**Union by size** (interchangeable with rank; handy when the problem asks for the largest component):

```go
if d.size[rx] < d.size[ry] { rx, ry = ry, rx }
d.parent[ry] = rx
d.size[rx] += d.size[ry]
```

**Map-based DSU** for non-integer / sparse elements (e.g. email strings, coordinates in a huge grid):

```go
parent := map[string]string{}
var find func(x string) string
find = func(x string) string {
    if _, ok := parent[x]; !ok {
        parent[x] = x // lazily register new elements as singletons
    }
    if parent[x] != x {
        parent[x] = find(parent[x])
    }
    return parent[x]
}
```

**Grid flattening** — treat cell `(r, c)` of an `rows × cols` grid as element `r*cols + c`; add one extra element `rows*cols` as a **virtual node** when you need a "connected to the border / to the top / to water" super-group.

---

## Worked example — step-by-step trace

Problem: given `n = 6` nodes (0..5) and edges `[(0,1), (1,2), (3,4), (2,0), (4,5)]`, how many connected components remain, and does any edge close a cycle?

Initial state — six singletons:

```
index:  0  1  2  3  4  5
parent: 0  1  2  3  4  5
rank:   0  0  0  0  0  0        count = 6
```

**1. Union(0, 1)** — `Find(0)=0`, `Find(1)=1`. Different roots, equal ranks (0,0): attach 1 under 0, bump rank[0] to 1.

```
parent: 0  0  2  3  4  5
rank:   1  0  0  0  0  0        count = 5      sets: {0,1} {2} {3} {4} {5}
```

**2. Union(1, 2)** — `Find(1)`: parent[1]=0, 0 is root → root 0. `Find(2)=2`. rank[0]=1 > rank[2]=0: attach 2 under 0.

```
parent: 0  0  0  3  4  5
rank:   1  0  0  0  0  0        count = 4      sets: {0,1,2} {3} {4} {5}
```

**3. Union(3, 4)** — roots 3 and 4, equal ranks: attach 4 under 3, rank[3]→1.

```
parent: 0  0  0  3  3  5
rank:   1  0  0  1  0  0        count = 3      sets: {0,1,2} {3,4} {5}
```

**4. Union(2, 0)** — `Find(2)`: parent[2]=0 → root 0. `Find(0)=0`. **Same root → return false: edge (2,0) closes a cycle.** Nothing changes; count stays 3.

**5. Union(4, 5)** — `Find(4)`: parent[4]=3 → root 3. `Find(5)=5`. rank[3]=1 > rank[5]=0: attach 5 under 3.

```
parent: 0  0  0  3  3  3
rank:   1  0  0  1  0  0        count = 2      sets: {0,1,2} {3,4,5}
```

**Result:** `count = 2` components; edge `(2,0)` was the redundant (cycle-forming) edge.

**Path compression in action:** suppose instead step 2 had been `Union(2, 1)` on a chain `2→1→0`. Then a later `Find(2)` walks 2→1→0 and rewrites `parent[2] = 0` on the way back — the next `Find(2)` is a single hop. Over many operations, trees stay almost completely flat; that is where the O(α(n)) bound comes from.

---

## Common pitfalls and how to avoid them

1. **Comparing elements instead of roots.** `x == y`, `parent[x] == parent[y]`, or `union(parent[x], parent[y])` are all wrong — only `Find(x) == Find(y)` is a valid same-set test, and `Union` must link *roots*, not arbitrary nodes. Always call `Find` first.

2. **Forgetting path compression (or rank), then TLE.** A plain quick-union degenerates into a linked list — `Find` becomes O(n) and the whole solution O(n²). With n up to 10⁵–10⁶ this times out. Both optimizations are ~3 extra lines; always include them.

3. **Updating rank incorrectly.** Rank only increases when merging two trees of **equal** rank, and only on the surviving root. Bumping rank on every union, or updating the absorbed root, silently degrades to O(log n) or worse.

4. **Reading `size[x]` of a non-root.** After merges, `size` is only maintained on roots. The size of x's component is `size[Find(x)]`, never `size[x]`.

5. **Stale `count` / recounting components wrong.** Decrement `count` **only when the union actually merges two different sets** (i.e. inside the `rx != ry` branch). Alternatively, count roots at the end: `parent[i] == i`. Counting `Find(i) == i` over a map-DSU without registering all elements also misses isolated nodes.

6. **Off-by-one on node labels.** LeetCode graphs are sometimes 1-indexed (`n` nodes labeled 1..n). Either allocate `n+1` slots or subtract 1 everywhere — mixing the two causes index-out-of-range or, worse, silently unions node 0 into everything.

7. **Recursion depth on huge inputs.** Recursive `Find` on a pathological chain of 10⁶ nodes can blow the stack (Go's default goroutine stack grows, but other languages crash). The iterative path-halving version is the safe default.

8. **Using DSU where it can't work.** Edge *deletions*, directed reachability, or "print the actual path" — DSU handles none of these directly. For deletions, process operations in reverse; for the others, use BFS/DFS/SCC.

9. **Grid problems: unioning across the wrong neighbours.** When flattening a grid, only union a cell with its **already-valid** neighbours (usually just up and left during a row-major scan — right and down get covered when those cells are visited), and double-check the flatten formula `r*cols + c` (not `r*rows + c`).

10. **Mutating during iteration in map-based DSU.** Lazily inserting into the parent map inside `find` while ranging over that same map in Go is undefined-order trouble. Register all keys first, or collect keys before iterating.

---

## Problems in this repo

Problems currently in the repo whose README lists Union Find as a category (both are solved with DFS/BFS or hashing as the primary approach; Union Find is the alternative pattern):

- [0128 — Longest Consecutive Sequence](/0128_longest_consecutive_sequence/README.md) — union `x` with `x+1` when both exist; the largest set size is the answer (hash-set scan is the usual optimal, DSU is the equivalence-class view).
- [0130 — Surrounded Regions](/0130_surrounded_regions/README.md) — classic virtual-node trick: union every border-connected `O` to a virtual "safe" node; any `O` not connected to it gets flipped to `X`.

> Problems 0131–0400 are being written concurrently; a later pass will add
> further Union Find problems (e.g. #200 Number of Islands, #261 Graph Valid
> Tree, #305 Number of Islands II, #323 Number of Connected Components,
> #399 Evaluate Division) as they land in the repo.

### Related classics to know (not yet in repo)

- LeetCode #200 — Number of Islands · #305 — Number of Islands II (DSU shines: online cell additions)
- LeetCode #547 — Number of Provinces (the canonical "count components" DSU problem)
- LeetCode #684 — Redundant Connection (cycle detection via failed union)
- LeetCode #721 — Accounts Merge (map-based DSU over strings)
- LeetCode #1584 — Min Cost to Connect All Points (Kruskal's MST)
- LeetCode #990 — Satisfiability of Equality Equations (equivalence classes)
