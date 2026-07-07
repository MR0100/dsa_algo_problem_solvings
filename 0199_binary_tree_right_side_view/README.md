# 0199 — Binary Tree Right Side View

> LeetCode #199 · Difficulty: Medium
> **Categories:** Tree, Depth-First Search, Breadth-First Search, Binary Tree

---

## Problem Statement

Given the `root` of a binary tree, imagine yourself standing on the **right side** of it, return *the values of the nodes you can see ordered from top to bottom*.

**Example 1:**

```
Input: root = [1,2,3,null,5,null,4]
Output: [1,3,4]
```

Explanation:

```
        1            <---   you see 1
       / \
      2   3          <---   you see 3 (2 is hidden behind it)
       \   \
        5   4        <---   you see 4 (5 is hidden behind it)
```

**Example 2:**

```
Input: root = [1,2,3,4,null,null,null,5]
Output: [1,3,4,5]
```

Explanation:

```
          1          <---   you see 1
         / \
        2   3        <---   you see 3
       /
      4              <---   you see 4 (nothing to its right on this level)
     /
    5                <---   you see 5
```

**Example 3:**

```
Input: root = [1,null,3]
Output: [1,3]
```

**Example 4:**

```
Input: root = []
Output: []
```

**Constraints:**

- The number of nodes in the tree is in the range `[0, 100]`.
- `-100 <= Node.val <= 100`

---

## Company Frequency

| Company    | Frequency        | Last Reported |
|------------|------------------|---------------|
| Amazon     | ★★★★★ Very High  | 2024          |
| Meta       | ★★★★★ Very High  | 2024          |
| Microsoft  | ★★★★☆ High       | 2024          |
| Google     | ★★★★☆ High       | 2023          |
| Bloomberg  | ★★★☆☆ Medium     | 2023          |
| Apple      | ★★★☆☆ Medium     | 2023          |
| ByteDance  | ★★☆☆☆ Low        | 2022          |

> ⚠️ Frequency data is crowd-sourced from LeetCode Discuss, Glassdoor, and
> community interview reports. Treat as a signal, not a guarantee.

---

## DSA Concepts Used

- **Tree Traversal** — the answer is one node per depth; every solution is a
  disciplined walk over the tree that decides which node "wins" each level →
  see [`/dsa/tree_traversal.md`](/dsa/tree_traversal.md)
- **Breadth-First Search (Level Order)** — Approach 1 processes the tree one
  full level at a time and keeps the last node of each level → see
  [`/dsa/graph_bfs_dfs.md`](/dsa/graph_bfs_dfs.md)
- **Depth-First Search** — Approaches 2 and 3 recurse depth-by-depth, using
  child-visit order to guarantee the rightmost node is the one recorded → see
  [`/dsa/graph_bfs_dfs.md`](/dsa/graph_bfs_dfs.md)
- **Queue** — the BFS approach uses a FIFO queue as its level buffer → see
  [`/dsa/queue_deque.md`](/dsa/queue_deque.md)

---

## Approaches Overview

| # | Approach | Time | Space | When to use |
|---|----------|------|-------|-------------|
| 1 | BFS Level Order | O(n) | O(w) | The most intuitive answer; "rightmost per level" maps directly onto level-order traversal. Also generalises to left-view / averages / max-per-level |
| 2 | DFS Left-First (Overwrite) | O(n) | O(h) | Recursive one-liner; natural if you already traverse left→right and want to overwrite each depth's slot |
| 3 | DFS Right-First (Optimal) | O(n) | O(h) | Cleanest recursion: record only on first arrival at a depth; no overwriting, no queue |

*(n = number of nodes, w = max tree width, h = tree height.)*

---

## Approach 1 — BFS Level Order

### Intuition

Standing on the right, you see exactly **one node per level**: the rightmost one. Breadth-first search naturally visits the tree level by level, so if you can tell when you are about to dequeue the *last* node of the current level, that node is precisely what is visible. The trick is to snapshot the queue's size at the start of each level — that count is the level's width, so the `size`-th node dequeued is its rightmost.

### Algorithm

1. If `root` is `nil`, return an empty view (nothing is visible).
2. Push `root` into a queue.
3. While the queue is non-empty:
   1. Let `size = len(queue)` — the number of nodes on the current level.
   2. Dequeue exactly `size` nodes. For each dequeued node, enqueue its non-nil `Left` then `Right` child.
   3. When dequeuing the `size`-th (index `size-1`) node, append its value to the view — it is the rightmost node of this level.
4. Return the view.

### Complexity

- **Time:** O(n) — every node is enqueued once and dequeued once.
- **Space:** O(w) — the queue holds at most one full level; the widest level can be up to ~n/2 nodes in a complete tree.

### Code

```go
func bfsLevelOrder(root *TreeNode) []int {
	view := []int{}
	if root == nil {
		return view // empty tree: nothing is visible
	}
	queue := []*TreeNode{root}
	for len(queue) > 0 {
		size := len(queue) // number of nodes on the current level
		for i := 0; i < size; i++ {
			node := queue[0]
			queue = queue[1:]
			if i == size-1 {
				view = append(view, node.Val) // last node of the level = rightmost
			}
			if node.Left != nil {
				queue = append(queue, node.Left)
			}
			if node.Right != nil {
				queue = append(queue, node.Right)
			}
		}
	}
	return view
}
```

### Dry Run

Example 1: `root = [1,2,3,null,5,null,4]`.

```
        1
       / \
      2   3
       \   \
        5   4
```

| Level | queue at level start | size | dequeue order (i) | children enqueued | view after level |
|-------|----------------------|------|-------------------|-------------------|------------------|
| 0 | `[1]` | 1 | i=0 → **1** (last) | 2, 3 | `[1]` |
| 1 | `[2, 3]` | 2 | i=0 → 2 (enq 5); i=1 → **3** (last, enq 4) | 5, 4 | `[1, 3]` |
| 2 | `[5, 4]` | 2 | i=0 → 5; i=1 → **4** (last) | — | `[1, 3, 4]` |

Queue empty → return `[1, 3, 4]` ✔

---

## Approach 2 — DFS Left-First (Overwrite)

### Intuition

Do an ordinary left-to-right DFS (root → left → right), but track the current `depth`. In this order, the **last** node visited at any given depth is the rightmost node of that depth — because everything to its left was visited earlier. So instead of trying to detect "the rightmost node," just keep overwriting `view[depth]` on every visit and let the final writer win. The first time you reach a new depth, you append; after that, you overwrite.

### Algorithm

1. Maintain a `view` slice indexed by depth.
2. DFS from the root with `depth = 0`, visiting **left before right**.
3. On visiting a node:
   - If `depth == len(view)`, this is the first node ever seen at this depth → append its value.
   - Otherwise, overwrite `view[depth]` with this node's value (a node further right is winning).
4. Recurse into `Left`, then `Right`.

### Complexity

- **Time:** O(n) — each node is visited exactly once.
- **Space:** O(h) — the recursion stack depth equals the tree height (O(n) for a fully skewed tree).

### Code

```go
func dfsLeftFirst(root *TreeNode) []int {
	view := []int{}
	var dfs func(node *TreeNode, depth int)
	dfs = func(node *TreeNode, depth int) {
		if node == nil {
			return
		}
		if depth == len(view) {
			view = append(view, node.Val) // first node seen at this depth
		} else {
			view[depth] = node.Val // a node further right overwrites the slot
		}
		dfs(node.Left, depth+1) // left first: the rightmost node writes last
		dfs(node.Right, depth+1)
	}
	dfs(root, 0)
	return view
}
```

### Dry Run

Example 1: `root = [1,2,3,null,5,null,4]`.

```
        1
       / \
      2   3
       \   \
        5   4
```

Visit order is root → left → right. Each row shows the node visited, its depth, whether it appends (new depth) or overwrites, and the resulting `view`.

| Step | node | depth | len(view) | action | view after |
|------|------|-------|-----------|--------|------------|
| 1 | 1 | 0 | 0 | append (new depth) | `[1]` |
| 2 | 2 | 1 | 1 | append (new depth) | `[1, 2]` |
| 3 | 5 | 2 | 2 | append (new depth) | `[1, 2, 5]` |
| 4 | 3 | 1 | 3 | overwrite view[1] | `[1, 3, 5]` |
| 5 | 4 | 2 | 3 | overwrite view[2] | `[1, 3, 4]` |

Traversal ends → return `[1, 3, 4]` ✔ — node 3 overwrote 2, and node 4 overwrote 5.

---

## Approach 3 — DFS Right-First (Optimal)

### Intuition

Same DFS idea, but flip the child order to root → **right** → left. Now the **first** node reached at any depth is the rightmost node of that depth — every node to its left is visited strictly later. So you only ever record on *first arrival* at a depth (`depth == len(view)`), and simply ignore every later visit to that depth. No overwriting, no queue — each answer is nailed at the earliest possible moment.

### Algorithm

1. Maintain a `view` slice indexed by depth.
2. DFS from the root with `depth = 0`, visiting **right before left**.
3. On visiting a node: if `depth == len(view)`, this is the first (hence rightmost) node seen at this depth → append its value. Otherwise do nothing.
4. Recurse into `Right`, then `Left`.

### Complexity

- **Time:** O(n) — each node is visited exactly once.
- **Space:** O(h) — recursion stack only, equal to the tree height; no auxiliary queue.

### Code

```go
func dfsRightFirst(root *TreeNode) []int {
	view := []int{}
	var dfs func(node *TreeNode, depth int)
	dfs = func(node *TreeNode, depth int) {
		if node == nil {
			return
		}
		if depth == len(view) {
			view = append(view, node.Val) // first arrival at this depth = rightmost
		}
		dfs(node.Right, depth+1) // right first, so the rightmost node wins each depth
		dfs(node.Left, depth+1)
	}
	dfs(root, 0)
	return view
}
```

### Dry Run

Example 1: `root = [1,2,3,null,5,null,4]`.

```
        1
       / \
      2   3
       \   \
        5   4
```

Visit order is root → right → left. We append only when `depth == len(view)` (first arrival); later visits at an already-recorded depth are skipped.

| Step | node | depth | len(view) | first at depth? | view after |
|------|------|-------|-----------|-----------------|------------|
| 1 | 1 | 0 | 0 | yes → append | `[1]` |
| 2 | 3 | 1 | 1 | yes → append | `[1, 3]` |
| 3 | 4 | 2 | 2 | yes → append | `[1, 3, 4]` |
| 4 | 2 | 1 | 3 | no (depth 1 filled) → skip | `[1, 3, 4]` |
| 5 | 5 | 2 | 3 | no (depth 2 filled) → skip | `[1, 3, 4]` |

Traversal ends → return `[1, 3, 4]` ✔ — each depth was captured on its first (rightmost) visit; nodes 2 and 5 were correctly skipped.

---

## Key Takeaways

- **"One value per level" = level-indexed traversal.** Whether you go BFS or DFS, the core move is grouping nodes by depth and picking one winner per depth. This pattern also solves left-side view, per-level averages/max/min, and level-order zigzag with tiny tweaks.
- **In BFS, snapshot `len(queue)` before draining a level** — that count *is* the level width, letting you identify the first node (left view) or last node (right view) without extra bookkeeping.
- **Child-visit order encodes priority in DFS.** Left-first + overwrite and right-first + first-write are two sides of the same coin; right-first is cleaner because it never rewrites a slot.
- **`depth == len(view)` is the idiomatic "first time this deep" test** in Go — it doubles as both the append condition and the "have we recorded this depth yet?" check.
- Handle the **empty tree** (`root == nil`) up front — it must return `[]`, not `nil`-panic.

---

## Related Problems

- LeetCode #102 — Binary Tree Level Order Traversal (BFS per-level buffering)
- LeetCode #107 — Binary Tree Level Order Traversal II (level order, bottom-up)
- LeetCode #103 — Binary Tree Zigzag Level Order Traversal (level order with alternating direction)
- LeetCode #513 — Find Bottom Left Tree Value (left-view / last level's leftmost)
- LeetCode #637 — Average of Levels in Binary Tree (one aggregate per level)
- LeetCode #116 — Populating Next Right Pointers in Each Node (level-linked traversal)
