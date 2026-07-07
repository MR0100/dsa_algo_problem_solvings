# Tree DP (Dynamic Programming on Trees)

> **Core idea:** run a single **post-order DFS**. Each node computes its answer from the *already-computed* answers of its children, returning a small **state tuple** up the recursion. The tree's global answer is either the root's state or a running maximum updated at every node.
> **Complexity:** O(n) time, O(h) stack space — every node and edge is touched exactly once.

---

## What it is

Tree DP is dynamic programming where the **DAG of subproblems is the tree itself**. A subproblem is "the best answer for the subtree rooted at node `v`". Because a child's subtree is strictly smaller and fully contained in its parent's, there is a natural evaluation order: **children before parents**, i.e. **post-order** traversal.

Each recursive call returns a **fixed-size state** — one number, or a tuple of a few numbers — that is *sufficient* for the parent to compute its own state without looking further down. The art of tree DP is choosing that state so it is:

- **small** (O(1) values, independent of subtree size), and
- **composable** — the parent can combine its children's states in O(1) (or O(children)).

Typical state shapes:

| Problem | State returned per node | Combine rule |
|---------|-------------------------|--------------|
| House Robber III (#337) | `(rob, notRob)` — best sum if we **take** this node vs **skip** it | `rob = val + Σ child.notRob`; `notRob = Σ max(child.rob, child.notRob)` |
| Max Path Sum (#124) | `down` — best *straight* path from this node **downward** into one child | node's best *through-path* = `val + max(0,left.down) + max(0,right.down)`; return `val + max(0, max(left.down, right.down))` |
| Diameter / height (#543-style, #104) | `height` — longest downward chain | update global `diameter = max(diameter, leftH + rightH)`; return `1 + max(leftH, rightH)` |
| Count Univalue Subtrees (#250) | `bool` — is this subtree all one value? | true iff both children are univalue *and* equal this node's value (respecting nils) |
| Find Leaves (#366) | `height` (distance to farthest leaf) | bucket node into `answer[height]`; return `1 + max(leftH, rightH)` |
| Balanced check (#110) | `height`, with a sentinel `-1` meaning "unbalanced below" | if `|leftH-rightH|>1` propagate `-1` |

### The two ways the answer surfaces

This is the single most important distinction in tree DP:

1. **Answer = root's returned state.** The value you want *is* what bubbles up to the top (House Robber III: `max(root.rob, root.notRob)`; univalue count if you thread a counter).
2. **Answer = a global maximum updated at every node**, while the function *returns something different*. This happens when the optimal structure **cannot be extended by the parent**. In Max Path Sum and Diameter, a path that bends through a node (uses *both* children) is a complete answer at that node — the parent can't use it, because a parent can only attach to a single straight chain. So we **record** the bend into a global, but **return** only the extendable straight chain.

Mixing these up ("returning the through-path") is the #1 tree-DP bug.

---

## When to recognise it

Reach for tree DP when the input is a tree (or forest) and the question asks for an optimum, count, or property that depends on **combining subtree results**:

| Signal in the problem | Why tree DP fits |
|-----------------------|------------------|
| "Binary tree" + "maximum / minimum / count of ..." | a global optimum decomposes over subtrees → post-order combine |
| **Adjacency constraint between a node and its children** ("can't rob two directly-linked houses", "no two adjacent") | classic take/skip state `(include, exclude)` |
| "Path", "diameter", "longest chain" in a tree | height-returning DFS + global max over `leftH + rightH` |
| "Is every subtree / how many subtrees satisfy property P?" | return a boolean (or richer info) and count on the way up |
| "Collect / group nodes by their distance to the nearest leaf" | leaf-height DFS, bucket by returned height (#366) |
| "Is the tree balanced / is this subtree a valid BST?" | return height (or `(min,max,size,isValid)`) with a sentinel for failure |
| Minimum Height Trees (#310) | the *centroid* view — trees have ≤2 centroids; peel leaves layer by layer (a BFS "reverse tree DP") |

**When it is *not* tree DP:** simple traversals that only *emit* nodes (inorder print, level-order) carry no combined state. Shortest path between two arbitrary nodes in a *graph* is BFS/Dijkstra, not tree DP. If the recurrence needs information from *ancestors or siblings*, a single post-order pass is not enough — you need a second (pre-order / re-rooting) pass (see "rerooting" pitfall).

---

## Contrast with array/linear DP

Tree DP is the same principle as 1-D DP — *optimal substructure + overlapping-free evaluation in dependency order* — but the **shape of the dependency graph** differs, and that changes the mechanics:

| | Linear DP (e.g. House Robber #198) | Tree DP (e.g. House Robber III #337) |
|---|---|---|
| Subproblem space | `dp[i]` over indices `0..n-1` | one subproblem per **node** |
| Dependency order | left→right index order | **post-order** (children first) |
| Transition reads | `dp[i-1]`, `dp[i-2]` (fixed neighbours) | states of **all children** |
| Storage | a `dp[]` array (or two rolling scalars) | the **call stack** (return value *is* the memo) |
| "Take / skip" recurrence | `take=nums[i]+dp[i-2]`, `skip=dp[i-1]` | `rob=val+Σchild.notRob`, `notRob=Σmax(child.rob,child.notRob)` |
| Answer location | `dp[n-1]` (last cell) | root's state, or a global updated everywhere |

The linear "rob a *row* of houses" and the tree "rob a *tree* of houses" are literally the same DP; only the neighbour set changes from `{i-1, i-2}` to `{children}`. Recognising this equivalence is the payoff of studying both. Note there is usually **no explicit memo table** in tree DP: because every node has exactly one parent, each subproblem is visited once, so plain recursion already runs in O(n) — the return value plays the role the `dp[]` array plays in linear DP.

---

## General template / pseudocode

Assume the standard LeetCode node:

```go
type TreeNode struct {
    Val   int
    Left  *TreeNode
    Right *TreeNode
}
```

### Template A — return a state tuple; answer is the root's state

House Robber III: each node returns `(rob, notRob)`.

```go
// rob returns the maximum money obtainable from the tree rooted at `root`.
//
// State per node: a pair (withNode, withoutNode)
//   withNode    = best sum that INCLUDES this node's value
//                 → children must then be skipped
//   withoutNode = best sum that EXCLUDES this node
//                 → each child is free to be taken or not (whichever is larger)
//
// Post-order: we need both children's pairs before computing this node's pair.
func rob(root *TreeNode) int {
    withRoot, withoutRoot := dfs(root)
    return max(withRoot, withoutRoot) // root is free: take the better option
}

func dfs(node *TreeNode) (withNode, withoutNode int) {
    if node == nil {
        return 0, 0 // empty subtree contributes nothing under either choice
    }
    lWith, lWithout := dfs(node.Left)  // solve children first (post-order)
    rWith, rWithout := dfs(node.Right)

    // Take this node ⇒ we are forbidden to take either child.
    withNode = node.Val + lWithout + rWithout
    // Skip this node ⇒ each child independently picks its own best.
    withoutNode = max(lWith, lWithout) + max(rWith, rWithout)
    return withNode, withoutNode
}

func max(a, b int) int { if a > b { return a }; return b }
```

### Template B — return an extendable chain; answer is a global max (bend recorded here)

Binary Tree Maximum Path Sum, and the diameter/#543 pattern share this skeleton.

```go
// maxPathSum: a path may bend through a node using BOTH children, but such a
// bent path CANNOT be extended by the parent — so we record it in `best` and
// return only the straight downward gain the parent can actually attach to.
func maxPathSum(root *TreeNode) int {
    best := math.MinInt // paths can be all-negative, so no 0 floor here
    var gain func(node *TreeNode) int
    gain = func(node *TreeNode) int {
        if node == nil {
            return 0
        }
        // Clamp negative child contributions to 0 = "don't take that side".
        left := max(0, gain(node.Left))
        right := max(0, gain(node.Right))

        // Candidate BEST path that peaks (bends) at this node.
        best = max(best, node.Val+left+right)

        // Return the straight chain (node + ONE side) the parent can extend.
        return node.Val + max(left, right)
    }
    gain(root)
    return best
}
```

Diameter (#543) is Template B with `node.Val` dropped and edges counted:

```go
func diameterOfBinaryTree(root *TreeNode) int {
    diameter := 0
    var height func(node *TreeNode) int // longest downward chain in EDGES... here in nodes
    height = func(node *TreeNode) int {
        if node == nil {
            return 0
        }
        lh, rh := height(node.Left), height(node.Right)
        diameter = max(diameter, lh+rh) // path through node = left depth + right depth (edges)
        return 1 + max(lh, rh)          // return height (nodes) for the parent
    }
    height(root)
    return diameter
}
```

### Template C — return a boolean/rich struct, count/validate on the way up

Count Univalue Subtrees (#250): return whether the subtree is univalue; increment a counter.

```go
func countUnivalSubtrees(root *TreeNode) int {
    count := 0
    var isUni func(node *TreeNode) bool
    isUni = func(node *TreeNode) bool {
        if node == nil {
            return true // vacuously univalue; lets leaves' nil children pass
        }
        left := isUni(node.Left)   // evaluate BOTH children (no short-circuit!)
        right := isUni(node.Right) // — we must recurse to count deeper subtrees
        if !left || !right {
            return false
        }
        // Children are univalue; now their values must match this node's.
        if node.Left != nil && node.Left.Val != node.Val {
            return false
        }
        if node.Right != nil && node.Right.Val != node.Val {
            return false
        }
        count++ // this whole subtree is univalue
        return true
    }
    isUni(root)
    return count
}
```

### Template D — leaf-height bucketing (Find Leaves, #366)

A node's "layer" is its distance to its farthest leaf; equal-height nodes are removed together.

```go
func findLeaves(root *TreeNode) [][]int {
    var res [][]int
    var height func(node *TreeNode) int
    height = func(node *TreeNode) int {
        if node == nil {
            return -1 // so a leaf gets height 0
        }
        h := 1 + max(height(node.Left), height(node.Right)) // leaf-distance
        if h == len(res) {
            res = append(res, []int{}) // first node at this new layer
        }
        res[h] = append(res[h], node.Val) // bucket by distance-to-leaf
        return h
    }
    height(root)
    return res
}
```

### The odd one out — Minimum Height Trees (#310): peel leaves (topological "reverse DP")

Not a rooted post-order DFS, but the same subtree-combining spirit run *inward*. A tree has **at most two centroids**; they are what remain after repeatedly stripping all current leaves (like Kahn's topological sort on an undirected tree).

```go
func findMinHeightTrees(n int, edges [][]int) []int {
    if n == 1 {
        return []int{0}
    }
    adj := make([]map[int]bool, n)
    for i := range adj {
        adj[i] = map[int]bool{}
    }
    for _, e := range edges {
        adj[e[0]][e[1]] = true
        adj[e[1]][e[0]] = true
    }
    leaves := []int{}
    for i := 0; i < n; i++ {
        if len(adj[i]) == 1 { // degree-1 nodes are leaves
            leaves = append(leaves, i)
        }
    }
    remaining := n
    for remaining > 2 { // stop when ≤2 centroids remain
        remaining -= len(leaves)
        next := []int{}
        for _, leaf := range leaves {
            for nb := range adj[leaf] { // each leaf has exactly one neighbour
                delete(adj[nb], leaf)
                if len(adj[nb]) == 1 { // neighbour just became a leaf
                    next = append(next, nb)
                }
            }
        }
        leaves = next
    }
    return leaves // the 1 or 2 centroids
}
```

---

## Worked example

Take House Robber III on this tree (the LeetCode Example 1). Values shown; the answer is 7 (rob the two 3s):

```
        3
       / \
      2   3
       \   \
        3   1
```

Post-order visits leaves first, then parents. We annotate each node with its returned pair **(withNode, withoutNode)**.

| Step | Node (value) | Left pair | Right pair | `withNode = val + lWithout + rWithout` | `withoutNode = max(lPair) + max(rPair)` | Returns |
|------|--------------|-----------|------------|----------------------------------------|-----------------------------------------|---------|
| 1 | left-`2`'s child `3` (leaf) | (0,0) | (0,0) | 3 + 0 + 0 = **3** | 0 + 0 = **0** | (3, 0) |
| 2 | right-`3`'s child `1` (leaf) | (0,0) | (0,0) | 1 + 0 + 0 = **1** | 0 + 0 = **0** | (1, 0) |
| 3 | left child `2` | (3,0) *(from step 1, as its right child)* | (0,0) | 2 + 0 + 0 = **2** | max(3,0) + max(0,0) = **3** | (2, 3) |
| 4 | right child `3` | (0,0) | (1,0) *(from step 2)* | 3 + 0 + 0 = **3** | max(0,0) + max(1,0) = **1** | (3, 1) |
| 5 | root `3` | (2,3) *(step 3)* | (3,1) *(step 4)* | 3 + 3 + 1 = **7** | max(2,3) + max(3,1) = 3 + 3 = **6** | (7, 6) |

Root returns `(7, 6)`; the answer is `max(7, 6) = 7`.

Trace intuition: at the root, **taking** the root (value 3) forces us to skip both children, but each child's `withoutNode` already banked its own best grandchild — `3` (the left grandchild) and `1` (the right grandchild) — giving `3 + 3 + 1 = 7`. **Skipping** the root lets each child choose freely but they can't beat 6. So we rob the root plus the two leaf `3`/`1`... wait — the 7 comes from root(3) + leftGrandchild(3) + rightGrandchild(1) = 7, exactly the non-adjacent set the DP found without ever enumerating subsets.

---

## Complexity

| Aspect | Cost | Reason |
|--------|------|--------|
| Time | **O(n)** | one visit per node; each combine step is O(1) for a binary tree (O(deg) for general trees, summing to O(n) over all edges) |
| Space (call stack) | **O(h)** | recursion depth = tree height `h`; `O(log n)` if balanced, `O(n)` if degenerate/skewed |
| Auxiliary state | **O(1)** per node | the returned tuple is fixed-size; no per-node memo table is needed because each subproblem is reached exactly once |

There is **no memoization table** in classic rooted tree DP — the single-parent property already guarantees each node's subproblem is solved once. (Memoization *does* reappear when the "tree" is actually a DAG, or when you re-root and want to cache, or in "DP on subtrees indexed by (node, extra-parameter)".)

---

## Common pitfalls

1. **Returning the bent/through path instead of the extendable chain (Template B).** In Max Path Sum / Diameter you must *record* `val + left + right` into a global but *return* `val + max(left, right)`. Returning the through-path lets a parent attach to a path that already used both children — impossible — inflating the answer.

2. **Using `0` as the initial best when values can be negative (#124).** Max Path Sum inputs can be entirely negative; seeding `best := 0` wrongly returns 0 for `[-3]`. Seed `best := math.MinInt`. The `max(0, childGain)` clamp is separate — it means "drop a negative branch", which is legitimate because a path may consist of a single node.

3. **Short-circuiting the recursion when you still need to count/visit deeper (#250).** Writing `return isUni(left) && isUni(right) && ...` lets Go's `&&` skip the right subtree once the left is false — but then you miss counting univalue subtrees on the skipped side. Evaluate both children into variables *first*, then combine.

4. **Confusing the two "answer surfaces."** Decide up front: is the answer the **root's return value** (House Robber III, balanced-check) or a **global max updated at every node** (diameter, max path sum)? Trying to return the global answer up the stack, or forgetting to update the global, are dual failure modes.

5. **Nil-child handling / base case value.** The base case's returned state must be the *identity* for your combine: `(0,0)` for rob/notRob sums, `0` (or `-1`) for heights depending on whether you count nodes or edges, `true` for "is univalue". An off-by-one here (e.g. returning `0` vs `-1` for a nil in leaf-height code) shifts every layer.

6. **Height in edges vs nodes.** Diameter is usually asked in **edges** (path through node = `leftDepth + rightDepth`), while many height problems count **nodes** (`1 + max(lh, rh)`). Pin down which the problem wants; the off-by-one silently gives answers that are one too big or small.

7. **Needing ancestor/sibling info → one pass is not enough (rerooting).** "For *every* node, the answer if the tree were rooted there" (sum of distances to all nodes, etc.) can't be done in a single post-order pass — you need a second **pre-order re-rooting** pass that pushes parent-side information down. Recognise when your recurrence looks *upward*, not just downward.

8. **Stack overflow on skewed trees.** Recursion depth is O(h); a 10⁵-node right-skewed tree can exhaust the default stack in some languages. Go grows goroutine stacks, so it usually survives, but an explicit stack (iterative post-order) is the defensive choice for pathological depth.

9. **Minimum Height Trees is *not* "root each node and DFS".** The brute force (BFS from every node, take min eccentricity) is O(n²) and TLEs. The intended solution is the leaf-peeling centroid trick — remember trees have **≤2** centroids.

---

## Problems in this repo that use it

- [0337 — House Robber III](/0337_house_robber_iii/README.md) — take/skip state `(rob, notRob)`; the tree analogue of linear House Robber.
- [0124 — Binary Tree Maximum Path Sum](/0124_binary_tree_maximum_path_sum/README.md) — Template B: record the bent path in a global, return the straight chain.
- [0310 — Minimum Height Trees](/0310_minimum_height_trees/README.md) — leaf-peeling to the ≤2 centroids (reverse/topological tree DP).
- [0250 — Count Univalue Subtrees](/0250_count_univalue_subtrees/README.md) — return a boolean, count univalue subtrees on the way up.
- [0366 — Find Leaves of Binary Tree](/0366_find_leaves_of_binary_tree/README.md) — bucket nodes by distance-to-leaf using a height-returning DFS.
- [0110 — Balanced Binary Tree](/0110_balanced_binary_tree/README.md) — height-returning DFS with a `-1` sentinel that short-circuits "unbalanced below".
- [0333 — Largest BST Subtree](/0333_largest_bst_subtree/README.md) — return a rich tuple `(min, max, size, isBST)`; the archetype of "carry more than one value up".
- [0298 — Binary Tree Longest Consecutive Sequence](/0298_binary_tree_longest_consecutive_sequence/README.md) — DFS returning the current consecutive run length, updating a global maximum.

### Linear-DP cousins (study the contrast)

- [0198 — House Robber](/0198_house_robber/README.md) — the same take/skip recurrence over an array (`{i-1, i-2}` neighbours instead of children).
- [0213 — House Robber II](/0213_house_robber_ii/README.md) — circular array variant; still linear DP, run twice.

### Related classics to know (may not be in repo)

- LeetCode #543 — Diameter of Binary Tree (the canonical Template B height/diameter problem).
- LeetCode #104 — Maximum Depth of Binary Tree (the simplest height-returning DFS).
- LeetCode #687 — Longest Univalue Path · #563 — Binary Tree Tilt (both Template B with a twist).
- LeetCode #834 — Sum of Distances in Tree (the rerooting / two-pass classic).
