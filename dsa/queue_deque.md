# Queue & Deque

> **Pattern family:** FIFO Queue · Double-Ended Queue (Deque) · Monotonic Deque
> **Core use cases:** BFS / level-order traversal, sliding-window maximum/minimum,
> order-preserving buffering, 0-1 BFS.

---

## 1. What the concept is

### Queue (FIFO)

A **queue** is a linear collection with two ends and one rule: elements are
**enqueued at the back** and **dequeued from the front** — *First In, First
Out*. The element that has waited longest is always served next.

```
enqueue →  [ back | ... | front ]  → dequeue
```

| Operation | Meaning                 | Cost (amortised) |
|-----------|-------------------------|------------------|
| Enqueue   | insert at back          | O(1)             |
| Dequeue   | remove from front       | O(1)             |
| Peek      | read front, don't remove| O(1)             |
| Empty?    | length check            | O(1)             |

### Deque (Double-Ended Queue)

A **deque** generalises the queue: you may **push and pop at *both* ends** in
O(1). It can therefore behave as a queue (push back / pop front), a stack
(push back / pop back), or something strictly more powerful — which is exactly
what the *monotonic deque* pattern exploits.

```
push front → [ front | ... | back ] ← push back
pop  front ← [ front | ... | back ] → pop  back
```

### Monotonic deque (the interview superstar)

A deque whose stored values (or the values at the stored *indices*) are kept
in strictly increasing or decreasing order. Before pushing a new element at
the back, you pop every back element that violates the order. Result: the
**front of the deque is always the max (or min) of the current window** — the
key to O(n) sliding-window-maximum, and to many DP-with-window-transition
problems.

Each element is pushed once and popped at most once → **total O(n)** across
the whole scan, i.e. **amortised O(1)** per step.

---

## 2. How to recognise a queue/deque problem

Signals in the problem statement:

- **"Level by level" / "level order" / "layer" / "depth of a tree"** — BFS
  with a plain FIFO queue. (#102, #103, #104, #107, #111)
- **"Shortest path in an unweighted graph"** or "minimum number of steps /
  transformations / moves" — BFS explores states in increasing distance
  order, so the *first* time you reach the target is optimal. (#127 Word
  Ladder, #130's BFS variant, knight moves, sliding-puzzle problems)
- **"Sliding window maximum / minimum"** or "max of every window of size k"
  — monotonic deque. (#239)
- **"Process in the order received"** / task scheduling / recent-calls
  counters (#933) — plain queue as an order-preserving buffer.
- **"Design a queue / circular queue / stack using queues"** (#232, #225,
  #622, #641) — direct implementation questions.
- **DP where `dp[i]` depends on `max(dp[i-k..i-1])`** — monotonic deque
  turns an O(n·k) DP into O(n). (#1696 Jump Game VI, #862)
- **Alternating left-to-right / right-to-left traversal ("zigzag")** — a
  deque lets you consume from either end. (#103)
- **Edge weights only 0 or 1, shortest path** — 0-1 BFS: push weight-0
  neighbours to the *front* of a deque, weight-1 neighbours to the back.

Rule of thumb: **stack = "nearest thing looking backwards"; queue =
"explore in arrival/distance order"; deque = "I need the best of a moving
window" or "I need both ends."**

---

## 3. Go templates

Go has no built-in queue/deque type. In interviews (and this repo) the idiom
is a **slice**.

### 3.1 Plain FIFO queue (BFS)

```go
// Generic BFS skeleton over an implicit graph.
//
// Pseudocode:
//   queue ← {start}; visited ← {start}
//   while queue not empty:
//       node ← pop front
//       if node is goal: return distance
//       for each unvisited neighbour: mark visited, push back
func bfs(start Node) int {
    queue := []Node{start}          // slice used as queue; front = index 0
    visited := map[Node]bool{start: true}
    steps := 0                      // distance of the current level from start

    for len(queue) > 0 {
        levelSize := len(queue)     // freeze size: process exactly one level
        for i := 0; i < levelSize; i++ {
            node := queue[0]        // peek front
            queue = queue[1:]       // dequeue (O(1) amortised; see pitfalls)

            if isGoal(node) {
                return steps        // first arrival = shortest, because BFS
            }
            for _, next := range neighbours(node) {
                if !visited[next] { // mark BEFORE enqueueing, not after popping
                    visited[next] = true
                    queue = append(queue, next) // enqueue at back
                }
            }
        }
        steps++                     // finished a whole level → one step farther
    }
    return -1                       // goal unreachable
}
```

The `levelSize` freeze is the **level-order trick**: everything currently in
the queue is exactly one BFS level, so draining `levelSize` items groups
nodes by depth (used in #102/#103/#107/#111/#116/#117).

### 3.2 Deque on a slice

```go
// A slice supports all four deque operations; only popFront needs care.
dq := []int{}

dq = append(dq, x)            // push back
dq = append([]int{x}, dq...)  // push front — O(n)! see 3.3 for the O(1) fix
back := dq[len(dq)-1]         // peek back
dq = dq[:len(dq)-1]           // pop back
front := dq[0]                // peek front
dq = dq[1:]                   // pop front (re-slice; memory freed when slice dies)
```

If you genuinely need O(1) push-front (rare — mostly 0-1 BFS), use
`container/list` or a ring buffer with head/tail indices; for the monotonic
pattern below you never push front, so a slice is perfect.

### 3.3 Monotonic deque — sliding window maximum

```go
// slidingWindowMax returns the max of every window of size k in O(n).
//
// Invariant: dq holds INDICES whose values are strictly decreasing,
// so nums[dq[0]] is always the max of the current window.
//
// Pseudocode per index i:
//   1. pop front while it has slid out of the window (dq[0] <= i-k)
//   2. pop back  while nums[back] <= nums[i]   (they can never be max again)
//   3. push i at back
//   4. once i >= k-1, record nums[dq[0]]
func slidingWindowMax(nums []int, k int) []int {
    dq := []int{}                       // deque of indices, values decreasing
    res := make([]int, 0, len(nums)-k+1)

    for i, v := range nums {
        // (1) evict indices that fell out of the window on the LEFT
        if len(dq) > 0 && dq[0] <= i-k {
            dq = dq[1:]                 // pop front
        }
        // (2) evict smaller-or-equal values from the RIGHT:
        //     v arrives later AND is >=, so those can never be a window max
        for len(dq) > 0 && nums[dq[len(dq)-1]] <= v {
            dq = dq[:len(dq)-1]         // pop back
        }
        // (3) current index joins the candidate list
        dq = append(dq, i)
        // (4) window is complete from i = k-1 onward; front is the max
        if i >= k-1 {
            res = append(res, nums[dq[0]])
        }
    }
    return res
}
```

For a sliding-window **minimum**, flip the comparison in step (2) to
`nums[back] >= v` (keep values increasing).

### 3.4 Queue of tree nodes (level order)

```go
func levelOrder(root *TreeNode) [][]int {
    if root == nil {
        return nil                       // guard: never enqueue nil roots
    }
    res := [][]int{}
    queue := []*TreeNode{root}
    for len(queue) > 0 {
        n := len(queue)                  // nodes on this level
        level := make([]int, 0, n)
        for i := 0; i < n; i++ {
            node := queue[0]
            queue = queue[1:]
            level = append(level, node.Val)
            if node.Left != nil {        // enqueue children for next level
                queue = append(queue, node.Left)
            }
            if node.Right != nil {
                queue = append(queue, node.Right)
            }
        }
        res = append(res, level)
    }
    return res
}
```

---

## 4. Worked example — Sliding Window Maximum (LeetCode #239)

`nums = [1, 3, -1, -3, 5, 3, 6, 7]`, `k = 3`. Expected: `[3, 3, 5, 5, 6, 7]`.

Deque stores **indices**; shown as `idx(val)`. Front is leftmost.

| i | v  | (1) evict front? | (2) pop back while `val ≤ v` | (3) deque after push | (4) emit |
|---|----|------------------|------------------------------|----------------------|----------|
| 0 | 1  | no (empty)       | —                            | `0(1)`               | — (i<2)  |
| 1 | 3  | no               | pop `0(1)` (1 ≤ 3)           | `1(3)`               | — (i<2)  |
| 2 | -1 | no               | — (3 > -1)                   | `1(3) 2(-1)`         | **3**    |
| 3 | -3 | `dq[0]=1 > 3-3=0` → no | — (-1 > -3)            | `1(3) 2(-1) 3(-3)`   | **3**    |
| 4 | 5  | `dq[0]=1 ≤ 4-3=1` → pop `1(3)` | pop `3(-3)`, pop `2(-1)` (both ≤ 5) | `4(5)` | **5** |
| 5 | 3  | no               | — (5 > 3)                    | `4(5) 5(3)`          | **5**    |
| 6 | 6  | no               | pop `5(3)`, pop `4(5)`       | `6(6)`               | **6**    |
| 7 | 7  | no               | pop `6(6)`                   | `7(7)`               | **7**    |

Result: `[3, 3, 5, 5, 6, 7]`. ✔

Why it works: at every step the deque is a **decreasing staircase of "still
possibly the max" candidates**. Anything smaller than a newer element is
dominated (the newer element is bigger *and* stays in the window longer), so
popping it is safe. Each of the 8 indices was pushed once and popped at most
once → 16 deque operations total → O(n).

---

## 5. Common pitfalls (and fixes)

1. **`queue = queue[1:]` and memory.** Re-slicing never shrinks the backing
   array, so a long-lived queue can pin memory for already-dequeued elements
   (and for pointer elements, keep them from GC). Fine for typical LeetCode
   input sizes; for pointer queues you can `queue[0] = nil` before
   re-slicing, or use head/tail indices into one buffer.
2. **O(n) push-front with `append([]T{x}, dq...)`.** That copies the whole
   slice. If an algorithm truly needs push-front (0-1 BFS), use
   `container/list` or a ring buffer — don't accidentally turn O(n) into O(n²).
3. **Marking visited at dequeue time instead of enqueue time.** In BFS, if
   you mark a node visited only when you *pop* it, the same node can be
   enqueued many times first → blow-up on dense graphs (classic Word Ladder
   TLE). Mark **when you enqueue**.
4. **Forgetting to freeze `levelSize`.** Writing
   `for i := 0; i < len(queue); i++` while also appending children reads a
   *growing* length and merges levels. Snapshot `n := len(queue)` first.
5. **Storing values instead of indices in a monotonic deque.** With values
   alone you can't tell when the front has slid out of the window. Store
   **indices**; evict front when `dq[0] <= i-k`.
6. **Wrong strictness in the back-pop (`<` vs `<=`).** For window *maximum*,
   pop while `nums[back] <= v` — keeping equal older values wastes space and
   in count/index variants gives wrong answers. Decide duplicate handling
   deliberately.
7. **Evicting the front with a `for` loop out of habit.** At most one index
   leaves the window per step, so `if` suffices — a `while` is harmless but
   signals you don't know the invariant.
8. **Confusing deque with priority queue.** A monotonic deque gives the
   window max in O(1) *only because eviction is age-ordered*. If elements
   leave in arbitrary order, you need a heap (see `/dsa/heap.md`).
9. **Nil children / nil root enqueued in tree BFS.** Either guard before
   `append` or nil-check after popping — pick one convention and be
   consistent, or you'll dereference nil.
10. **Using BFS where DFS is simpler (or vice versa).** BFS guarantees
    *shortest* in unweighted graphs; DFS doesn't. If the problem asks for
    *any* path/existence, DFS is usually less code; if it asks for *minimum
    steps*, BFS is mandatory.

---

## 6. Problems in this repo

Plain FIFO queue — BFS / level order:

- [0017 — Letter Combinations of a Phone Number](../0017_letter_combinations_of_a_phone_number/README.md) — BFS queue expansion of partial combinations, level per digit
- [0101 — Symmetric Tree](../0101_symmetric_tree/README.md) — iterative check with a queue of mirrored node pairs
- [0102 — Binary Tree Level Order Traversal](../0102_binary_tree_level_order_traversal/README.md) — the canonical BFS-with-levelSize template
- [0103 — Binary Tree Zigzag Level Order Traversal](../0103_binary_tree_zigzag_level_order_traversal/README.md) — level order with alternating direction (deque-flavoured)
- [0104 — Maximum Depth of Binary Tree](../0104_maximum_depth_of_binary_tree/README.md) — BFS level counting
- [0107 — Binary Tree Level Order Traversal II](../0107_binary_tree_level_order_traversal_ii/README.md) — level order, bottom-up output
- [0111 — Minimum Depth of Binary Tree](../0111_minimum_depth_of_binary_tree/README.md) — BFS wins here: first leaf found is the shallowest
- [0116 — Populating Next Right Pointers in Each Node](../0116_populating_next_right_pointers_in_each_node/README.md) — level order to wire `next` pointers
- [0117 — Populating Next Right Pointers in Each Node II](../0117_populating_next_right_pointers_in_each_node_ii/README.md) — same, on an arbitrary binary tree
- [0127 — Word Ladder](../0127_word_ladder/README.md) — shortest transformation via BFS over word states; mark-visited-on-enqueue matters
- [0130 — Surrounded Regions](../0130_surrounded_regions/README.md) — BFS flood fill from the border

Monotonic deque: no solved problem below #131 uses it yet; #239 (Sliding
Window Maximum), #862, and #1696 are the flagship problems and will be
linked when solved. (Problems 0131+ are being added; a later pass will
extend this list.)

Related references: [`/dsa/sliding_window.md`](/dsa/sliding_window.md) (window
mechanics), [`/dsa/two_pointers.md`](/dsa/two_pointers.md).
