# Graph BFS & DFS

> **Pattern family:** Breadth-First Search · Depth-First Search · Grid Traversal ·
> Implicit Graphs · Connected Components · Flood Fill · Multi-Source BFS ·
> Bidirectional BFS
> **Core use cases:** shortest path in unweighted graphs, connectivity / component
> counting, region marking on grids, cycle detection, reachability, level-by-level
> processing, state-space search.

---

## 1. What the concept is

A **graph** is a set of nodes (vertices) plus edges connecting them. Many problems
never say the word "graph" — the graph is *implicit*: grid cells connected to their
4 neighbours, words connected when they differ by one letter, states connected by
one legal move. Recognising the hidden graph is usually the whole battle.

Two fundamental ways to explore a graph from a starting node:

### BFS — Breadth-First Search

Explore in **expanding rings**: first all nodes at distance 1, then distance 2, …
Uses a **FIFO queue**. Because nodes are dequeued in non-decreasing distance
order, **the first time BFS reaches a node is via a shortest path (in edge count)**.

```
        start
       /  |  \
      A   B   C      ← level 1 (dist 1) — all visited before…
     / \      |
    D   E    F       ← level 2 (dist 2)
```

### DFS — Depth-First Search

Explore **as deep as possible** along one branch before backtracking. Uses
**recursion** (implicit call stack) or an explicit stack. Natural for
"explore/consume an entire region", path enumeration, cycle detection, and
anything where you need to *finish* a subtree before moving on (topological order,
post-order computations).

### BFS vs DFS — decision table

| Question the problem asks                              | Use     | Why |
|--------------------------------------------------------|---------|-----|
| *Shortest* path / *minimum* steps (unweighted)         | **BFS** | first arrival = shortest |
| Process "level by level" / by distance / by generation | **BFS** | queue naturally layers |
| Does a path exist? / mark a whole region               | Either  | DFS is less code; BFS avoids stack overflow |
| Count connected components / flood fill                | Either  | run from every unvisited node |
| Enumerate *all* paths, backtracking with undo          | **DFS** | recursion mirrors the path |
| Cycle detection, topological sort, post-order values   | **DFS** | needs finish-time information |
| Very deep graph (recursion depth ~10⁵+)                | **BFS** or iterative DFS | avoid stack overflow |
| Shortest path with *weighted* edges                    | Neither | Dijkstra / Bellman-Ford (see `/dsa/heap_priority_queue.md`) |

Both run in **O(V + E)** time and **O(V)** space (visited set + queue/stack).
On an m×n grid: V = m·n, E ≤ 4·m·n → **O(m·n)**.

---

## 2. How to recognise a BFS/DFS problem

Signals in the problem statement:

- **"Shortest / minimum number of steps / transformations / moves"** on an
  unweighted structure → BFS. (#127 Word Ladder: minimum transformation
  sequence length.)
- **"Grid" / "matrix" of cells with regions, islands, or spreading** — implicit
  graph where each cell connects to 4 (or 8) neighbours. (#130 Surrounded
  Regions, #79 Word Search.)
- **"Connected" / "adjacent" / "region" / "island" / "province"** — connected
  components; run BFS/DFS from every unvisited node and count/mark.
- **"All paths" / "all sequences" / "every way to…"** → DFS with backtracking
  (see `/dsa/backtracking.md`). (#126 Word Ladder II: *all* shortest sequences.)
- **A transformation between states** ("change one letter", "flip one switch",
  "rotate one wheel") — nodes are states, edges are single moves → BFS over the
  state space.
- **"Spreads simultaneously from multiple sources"** (rotting oranges, fire,
  walls-and-gates) → **multi-source BFS**: seed the queue with *all* sources at
  distance 0.
- **"Level order" / "depth" of a tree** — BFS with a queue (trees are just
  acyclic graphs; no visited set needed). (#102, #104, #107, #111.)
- **"Can you finish / is there a valid ordering"** with prerequisites →
  cycle detection / topological sort via DFS colours or Kahn's BFS.

Red herring check: if edges have **weights**, plain BFS/DFS won't give shortest
paths — reach for Dijkstra instead.

---

## 3. General templates (Go)

### 3.1 BFS on an adjacency-list graph

```go
// bfs returns the minimum edge-distance from start to every reachable node.
//
// Time:  O(V + E) — each node enqueued once, each edge scanned once.
// Space: O(V)     — queue + dist/visited map.
func bfs(adj map[int][]int, start int) map[int]int {
    dist := map[int]int{start: 0} // doubles as the visited set
    queue := []int{start}         // FIFO queue, seeded with the source

    for len(queue) > 0 {
        node := queue[0]   // dequeue from the front
        queue = queue[1:]

        for _, next := range adj[node] { // scan every neighbour
            if _, seen := dist[next]; seen {
                continue // CRITICAL: mark/check visited at enqueue time
            }
            dist[next] = dist[node] + 1 // first arrival = shortest distance
            queue = append(queue, next) // explore it in a later ring
        }
    }
    return dist
}
```

### 3.2 BFS with explicit levels (when you need the layer number)

```go
// Process the queue one full level at a time — the standard trick for
// "minimum steps" answers and level-order traversal.
level := 0
queue := []int{start}
visited := map[int]bool{start: true}

for len(queue) > 0 {
    size := len(queue)          // freeze the size: exactly this level's nodes
    for i := 0; i < size; i++ { // consume only the frozen prefix
        node := queue[0]
        queue = queue[1:]

        // ... process node; if node is the target, return level ...

        for _, next := range neighbors(node) {
            if !visited[next] {
                visited[next] = true // mark when enqueuing, not when popping
                queue = append(queue, next)
            }
        }
    }
    level++ // everything enqueued during this pass belongs to the next ring
}
```

### 3.3 Recursive DFS on a grid (flood fill / region marking)

```go
// dfs marks every cell reachable from (r, c) through same-region cells.
//
// Time:  O(m·n) total across all calls — each cell visited once.
// Space: O(m·n) worst-case recursion depth (a snake-shaped region).
func dfs(grid [][]byte, r, c int, target, mark byte) {
    // bounds + "is this cell part of the region?" — the single guard clause
    if r < 0 || r >= len(grid) || c < 0 || c >= len(grid[0]) || grid[r][c] != target {
        return
    }
    grid[r][c] = mark // mark BEFORE recursing — this is the visited check

    dfs(grid, r+1, c, target, mark) // down
    dfs(grid, r-1, c, target, mark) // up
    dfs(grid, r, c+1, target, mark) // right
    dfs(grid, r, c-1, target, mark) // left
}
```

Direction-vector variant (scales to 8 directions / knight moves):

```go
var dirs = [4][2]int{{1, 0}, {-1, 0}, {0, 1}, {0, -1}}
for _, d := range dirs {
    dfs(grid, r+d[0], c+d[1], target, mark)
}
```

### 3.4 Iterative DFS (explicit stack — same code shape as BFS)

```go
// Identical to BFS except LIFO instead of FIFO: pop from the BACK.
stack := []int{start}
visited := map[int]bool{start: true}

for len(stack) > 0 {
    node := stack[len(stack)-1]      // peek back
    stack = stack[:len(stack)-1]     // pop back  ← the only line that differs

    for _, next := range neighbors(node) {
        if !visited[next] {
            visited[next] = true
            stack = append(stack, next)
        }
    }
}
```

### 3.5 Multi-source BFS (pseudocode)

```
seed queue with EVERY source node, all at distance 0
mark all sources visited
run normal BFS
→ dist[v] = distance from v to the NEAREST source
```

### 3.6 Bidirectional BFS (pseudocode)

```
beginSet ← {start};  endSet ← {target}
while both sets non-empty:
    always expand the SMALLER set by one level
    if the new frontier intersects the other set → paths met, return steps
→ cuts O(b^d) work to O(b^(d/2)) on high-branching graphs   (used in #127)
```

### 3.7 Counting connected components (pseudocode)

```
components ← 0
for every node v:
    if v not visited:
        components++
        BFS/DFS from v, marking everything reachable
```

---

## 4. Worked example — Surrounded Regions (#130)

Board (`X` = wall, `O` = open):

```
X X X X          row 0
X O O X          row 1
X X O X          row 2
X O X X          row 3
```

Goal: flip every `O` region *not* touching the border to `X`.

**Inverse-thinking trick:** finding "surrounded" regions directly is awkward;
instead DFS from every **border** `O`, mark all connected `O`s as safe (`S`),
then flip everything else.

Step-by-step trace:

| Step | Action | Board state / notes |
|------|--------|---------------------|
| 1 | Scan border cells: row 0, row 3, col 0, col 3 for `O` | Only `(3,1)` is a border `O` |
| 2 | `dfs(3,1)`: cell is `O` → mark `S` | `(3,1) = S` |
| 3 | Recurse down `(4,1)` | out of bounds → return |
| 4 | Recurse up `(2,1)` | `X`, not target → return |
| 5 | Recurse right `(3,2)` | `X` → return |
| 6 | Recurse left `(3,0)` | `X` → return; `dfs(3,1)` finished |
| 7 | Border scan complete | Interior region `{(1,1),(1,2),(2,2)}` never reached — it is surrounded |
| 8 | Final sweep over all cells | `O` → `X` (captured); `S` → `O` (restore safe) |

Result:

```
X X X X
X X X X
X X X X
X O X X      ← the border-connected O survives
```

Every cell is visited a constant number of times → **O(m·n)** time,
**O(m·n)** worst-case stack depth, and the in-place `S` marker doubles as the
visited set → no extra structure needed.

For a BFS worked example on an *implicit* graph, see the dry runs in
[`../0127_word_ladder/README.md`](../0127_word_ladder/README.md): each word is a
node, an edge joins words differing in one letter, and level-by-level BFS makes
the first arrival at `endWord` provably the shortest transformation sequence.

---

## 5. Common pitfalls and how to avoid them

1. **Marking visited at pop time instead of enqueue time (BFS).** The same node
   gets enqueued many times before it is first popped → blow-up to O(V²) or
   worse, and wrong level counts. **Always mark when you `append` to the queue.**
2. **Forgetting the visited set entirely.** On any graph with a cycle (or a grid,
   where `(r,c) → (r+1,c) → (r,c)` is a 2-cycle), the search loops forever.
   Trees are the *only* place you can skip it.
3. **Recursion depth on big grids.** A 10⁶-cell single region means ~10⁶ deep
   recursion. Go's goroutine stacks grow dynamically so this usually survives,
   but in interviews call it out and offer the iterative/BFS version.
4. **Confusing DFS with shortest path.** DFS finds *a* path, not the shortest.
   If the problem says "minimum", DFS + pruning is almost always the wrong tool —
   use BFS.
5. **Mutating shared path state in DFS without undoing it.** When enumerating
   paths (backtracking), restore the cell/slice after the recursive call —
   e.g. #79 Word Search marks a cell `'#'` before recursing and restores the
   letter afterwards. Also beware appending a shared slice into results without
   copying it.
6. **Level counting off-by-one.** Decide up front: does the start count as
   level 0 or 1? (#127 counts the start word, so the answer is `level` starting
   at 1.) Freeze `size := len(queue)` *before* the inner loop — reading
   `len(queue)` live mixes levels.
7. **Grid bounds checks scattered everywhere.** Do bounds + target check as the
   *first* guard clause inside the recursive call (template 3.3) instead of
   before each of the four calls — one place to get right instead of four.
8. **Building the graph explicitly when it's cheaper implicitly.** In #127,
   pre-computing all word pairs is O(N²·L); generating the 26·L one-letter
   mutations of the current word on the fly is O(N·26·L). Ask "can I compute
   neighbours on demand?" before materialising an adjacency list.
9. **Popping from the front of a Go slice and expecting O(1).** `queue[1:]`
   is O(1) but the backing array is never freed while the slice lives; for huge
   BFS frontiers use an index pointer (`head++`) or `container/list`.
10. **Using plain BFS on weighted edges.** First arrival ≠ shortest once weights
    differ. Weighted → Dijkstra; weights ∈ {0,1} → 0-1 BFS with a deque
    (see `/dsa/queue_deque.md`).

---

## 6. Problems in this repo

Graph / grid BFS & DFS:

- [0079 — Word Search](../0079_word_search/README.md) — DFS on a grid with
  backtracking: mark cell, recurse in 4 directions, unmark on return.
- [0126 — Word Ladder II](../0126_word_ladder_ii/README.md) — BFS to layer the
  implicit word graph by distance, then DFS along parent links to enumerate
  *all* shortest transformation sequences.
- [0127 — Word Ladder](../0127_word_ladder/README.md) — level-by-level BFS on an
  implicit graph (one-letter mutations as edges); bidirectional BFS optimisation.
- [0130 — Surrounded Regions](../0130_surrounded_regions/README.md) — inverse
  thinking: DFS/BFS flood fill from border `O`s marks the safe cells; everything
  unmarked is captured.

Tree BFS/DFS (same machinery, no visited set needed — trees have no cycles):

- [0102 — Binary Tree Level Order Traversal](../0102_binary_tree_level_order_traversal/README.md) — the canonical level-by-level BFS template (3.2).
- [0103 — Binary Tree Zigzag Level Order Traversal](../0103_binary_tree_zigzag_level_order_traversal/README.md) — level BFS with alternating direction.
- [0104 — Maximum Depth of Binary Tree](../0104_maximum_depth_of_binary_tree/README.md) — DFS depth vs BFS level count.
- [0107 — Binary Tree Level Order Traversal II](../0107_binary_tree_level_order_traversal_ii/README.md) — level BFS, reversed output.
- [0111 — Minimum Depth of Binary Tree](../0111_minimum_depth_of_binary_tree/README.md) — BFS shines: first leaf found = minimum depth, no need to explore deeper.
- [0112 — Path Sum](../0112_path_sum/README.md) / [0113 — Path Sum II](../0113_path_sum_ii/README.md) — root-to-leaf DFS; II adds backtracking to collect paths.
- [0129 — Sum Root to Leaf Numbers](../0129_sum_root_to_leaf_numbers/README.md) — DFS carrying accumulated state down the recursion.
- [0116](../0116_populating_next_right_pointers_in_each_node/README.md) / [0117 — Populating Next Right Pointers](../0117_populating_next_right_pointers_in_each_node_ii/README.md) — level-order connection with O(1)-space refinement.

> Problems #0131+ are being added; a later pass will extend this list
> (e.g. #133 Clone Graph, #200 Number of Islands, #207 Course Schedule).

Related references: [`/dsa/queue_deque.md`](/dsa/queue_deque.md) (the BFS queue,
0-1 BFS) · [`/dsa/stack.md`](/dsa/stack.md) (iterative DFS) ·
[`/dsa/heap_priority_queue.md`](/dsa/heap_priority_queue.md) (Dijkstra, when
edges gain weights).
