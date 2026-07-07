# Tree Traversal

> Visiting every node of a tree exactly once, in a well-defined order.
> The single most reusable skill on LeetCode trees: almost every tree problem
> is "pick the right traversal + do something at each node".

---

## What it is

A **tree** is a connected acyclic structure with a root; each node has zero or
more children (binary tree: at most `Left` and `Right`). A **traversal** is a
systematic walk that visits every node exactly once. Because a tree has no
cycles, traversal needs no `visited` set (unlike graphs) — the parent→child
direction alone guarantees termination.

The two families:

| Family | Data structure | Orders |
|--------|----------------|--------|
| **DFS** (depth-first) | Recursion / explicit stack | Preorder, Inorder, Postorder |
| **BFS** (breadth-first) | Queue | Level order (and variants: zigzag, bottom-up, right view) |

The three DFS orders differ only in **when the node itself is processed**
relative to its subtrees:

- **Preorder** — `Node, Left, Right` → process a node *before* its children.
  "Top-down": pass information from parent to child (prefix sums, path so far,
  copying/serialising a tree, building from preorder arrays).
- **Inorder** — `Left, Node, Right` → for a **BST this yields sorted order**.
  Anything about BST ordering (validate, k-th smallest, recover swapped nodes)
  screams inorder.
- **Postorder** — `Left, Right, Node` → process a node *after* its children.
  "Bottom-up": compute a value for a node from its children's answers
  (height, balance, subtree sums, max path sum, deleting/freeing a tree).

**Level order (BFS)** — visit depth 0, then depth 1, ... Anything phrased
"per level", "nearest/shallowest", "left/right view", "connect nodes on the
same level" wants BFS (or DFS carrying a `depth` parameter).

### How to recognise a tree-traversal problem

Signals in the problem statement:

- The input is a `TreeNode` root — 95 % of the time the answer is "traverse it".
- "Return the ... **traversal** of its nodes' values" — literal.
- "**Depth** / height / balanced / diameter" — postorder (children first).
- "**Level by level** / zigzag / bottom-up / right side view / minimum depth" — BFS.
- "**Root-to-leaf path** / path sum" — preorder DFS carrying running state down.
- "Is this a valid **BST** / k-th smallest in BST / two nodes swapped" — inorder.
- "Build a tree from **preorder+inorder** (or inorder+postorder) arrays" —
  exploit what each order tells you: pre/post gives the root, inorder splits
  left/right subtrees.
- "**Same tree / symmetric / mirror**" — traverse two trees (or two halves) in
  lock-step.
- "Flatten / serialise / clone a tree" — pick the order that makes the output
  format fall out naturally (usually preorder).

---

## Templates (Go)

Node definition used throughout:

```go
type TreeNode struct {
    Val   int
    Left  *TreeNode
    Right *TreeNode
}
```

### 1. Recursive DFS — the universal skeleton

```go
// dfs visits every node. Moving the "process(node)" line changes the order.
func dfs(node *TreeNode) {
    if node == nil {        // base case: empty subtree — always check first
        return
    }
    // process(node)        // ← HERE = preorder  (Node, Left, Right)
    dfs(node.Left)
    // process(node)        // ← HERE = inorder   (Left, Node, Right)
    dfs(node.Right)
    // process(node)        // ← HERE = postorder (Left, Right, Node)
}
```

Two common shapes of the recursion:

```go
// Top-down (preorder-style): parent passes state DOWN via parameters.
// e.g. Path Sum: "does some root-to-leaf path add up to target?"
func topDown(node *TreeNode, runningSum int, target int) bool {
    if node == nil {
        return false
    }
    runningSum += node.Val                       // absorb current node
    if node.Left == nil && node.Right == nil {   // leaf: decide here
        return runningSum == target
    }
    return topDown(node.Left, runningSum, target) ||
        topDown(node.Right, runningSum, target)
}

// Bottom-up (postorder-style): children RETURN values, parent combines them.
// e.g. Maximum Depth.
func bottomUp(node *TreeNode) int {
    if node == nil {
        return 0                       // empty tree has depth 0
    }
    left := bottomUp(node.Left)        // solve left subtree first
    right := bottomUp(node.Right)      // then right subtree
    return 1 + max(left, right)        // combine: this node adds one level
}
```

### 2. Iterative preorder — explicit stack

```go
func preorderIterative(root *TreeNode) []int {
    result := []int{}
    if root == nil {
        return result
    }
    stack := []*TreeNode{root}         // stack replaces the call stack
    for len(stack) > 0 {
        node := stack[len(stack)-1]    // pop the top
        stack = stack[:len(stack)-1]
        result = append(result, node.Val) // visit BEFORE children = preorder
        if node.Right != nil {         // push RIGHT first...
            stack = append(stack, node.Right)
        }
        if node.Left != nil {          // ...so LEFT pops (and is visited) first
            stack = append(stack, node.Left)
        }
    }
    return result
}
```

### 3. Iterative inorder — "slide left, pop, go right"

The workhorse for BST problems (and the basis of a BST iterator):

```go
func inorderIterative(root *TreeNode) []int {
    result := []int{}
    stack := []*TreeNode{}
    curr := root
    for curr != nil || len(stack) > 0 {
        for curr != nil {              // 1) slide as far left as possible,
            stack = append(stack, curr) //    stacking every node on the way
            curr = curr.Left
        }
        curr = stack[len(stack)-1]     // 2) pop the deepest unvisited node
        stack = stack[:len(stack)-1]
        result = append(result, curr.Val) // 3) visit it (leftmost remaining)
        curr = curr.Right              // 4) then explore its right subtree
    }
    return result
}
```

### 4. Iterative postorder — trick: reversed "Node, Right, Left"

```go
func postorderIterative(root *TreeNode) []int {
    result := []int{}
    if root == nil {
        return result
    }
    stack := []*TreeNode{root}
    for len(stack) > 0 {
        node := stack[len(stack)-1]
        stack = stack[:len(stack)-1]
        result = append(result, node.Val)  // build Node,Right,Left...
        if node.Left != nil {
            stack = append(stack, node.Left)
        }
        if node.Right != nil {
            stack = append(stack, node.Right)
        }
    }
    // ...then reverse: Left,Right,Node = postorder
    for i, j := 0, len(result)-1; i < j; i, j = i+1, j-1 {
        result[i], result[j] = result[j], result[i]
    }
    return result
}
```

### 5. Level order (BFS) — queue with level-size snapshot

```go
func levelOrder(root *TreeNode) [][]int {
    result := [][]int{}
    if root == nil {
        return result
    }
    queue := []*TreeNode{root}
    for len(queue) > 0 {
        size := len(queue)             // SNAPSHOT: nodes currently in queue
        level := make([]int, 0, size)  //           = exactly one level
        for i := 0; i < size; i++ {    // drain only this level
            node := queue[0]
            queue = queue[1:]          // dequeue front
            level = append(level, node.Val)
            if node.Left != nil {      // children join the NEXT level
                queue = append(queue, node.Left)
            }
            if node.Right != nil {
                queue = append(queue, node.Right)
            }
        }
        result = append(result, level)
    }
    return result
}
```

Variants fall out of this loop: reverse `result` at the end (bottom-up),
reverse `level` on odd depths (zigzag), keep only `level[size-1]` (right view),
return on the first leaf found (minimum depth — BFS stops early, DFS cannot).

### 6. Morris traversal — O(1) space inorder

Threads the tree temporarily instead of using a stack. Interview flex; also
the trick behind O(1)-space solutions to Recover BST and Flatten.

```go
func morrisInorder(root *TreeNode) []int {
    result := []int{}
    curr := root
    for curr != nil {
        if curr.Left == nil {
            result = append(result, curr.Val) // no left subtree: visit, go right
            curr = curr.Right
        } else {
            pred := curr.Left                  // find inorder predecessor:
            for pred.Right != nil && pred.Right != curr {
                pred = pred.Right              // rightmost node of left subtree
            }
            if pred.Right == nil {
                pred.Right = curr              // create thread back to curr
                curr = curr.Left               // and dive left
            } else {
                pred.Right = nil               // second arrival: remove thread
                result = append(result, curr.Val) // left subtree done — visit
                curr = curr.Right
            }
        }
    }
    return result
}
```

### Complexity summary

| Template | Time | Space |
|----------|------|-------|
| Recursive DFS | O(n) | O(h) call stack — O(log n) balanced, O(n) skewed |
| Iterative DFS | O(n) | O(h) explicit stack |
| BFS | O(n) | O(w) queue — up to O(n/2) ≈ O(n) for the widest level |
| Morris | O(n) — each edge walked ≤ 3× | **O(1)** |

---

## Worked example — inorder traversal, step by step

Tree (LeetCode #94, Example 1: `root = [1,null,2,3]`):

```
1
 \
  2
 /
3
```

Expected inorder (`Left, Node, Right`): `[1, 3, 2]`.

Trace of the **iterative** template (`stack`, `curr`, `result`):

| Step | Action | stack | curr | result |
|------|--------|-------|------|--------|
| 1 | start | `[]` | `1` | `[]` |
| 2 | slide left: push 1; 1 has no Left | `[1]` | `nil` | `[]` |
| 3 | pop 1, **visit** | `[]` | `1` | `[1]` |
| 4 | go right of 1 | `[]` | `2` | `[1]` |
| 5 | slide left: push 2, move to 3; push 3; 3 has no Left | `[2,3]` | `nil` | `[1]` |
| 6 | pop 3, **visit** | `[2]` | `3` | `[1,3]` |
| 7 | go right of 3 → nil | `[2]` | `nil` | `[1,3]` |
| 8 | pop 2, **visit** | `[]` | `2` | `[1,3,2]` |
| 9 | go right of 2 → nil; stack empty, curr nil → stop | `[]` | `nil` | `[1,3,2]` |

Same tree, recursive view (call nesting shows the order falling out naturally):

```
inorder(1)
├── inorder(nil)        left of 1 — returns
├── visit 1             → [1]
└── inorder(2)
    ├── inorder(3)
    │   ├── inorder(nil)
    │   ├── visit 3     → [1,3]
    │   └── inorder(nil)
    ├── visit 2         → [1,3,2]
    └── inorder(nil)
```

For contrast, on the same tree: preorder = `[1,2,3]`, postorder = `[3,2,1]`,
level order = `[[1],[2],[3]]`.

---

## Common pitfalls

1. **Forgetting the `nil` base case** — first line of every recursive helper
   must be `if node == nil { return ... }`. In Go this is a panic
   (`nil pointer dereference`), not a silent bug.
2. **Wrong stack push order in iterative preorder** — a stack is LIFO, so push
   **Right before Left** to visit Left first. Getting this backwards mirrors
   the whole traversal.
3. **BFS without the level-size snapshot** — if you don't capture
   `size := len(queue)` before the inner loop, children enqueued mid-loop
   bleed into the current level and level boundaries are lost.
4. **Validating a BST by checking only `node.Left.Val < node.Val < node.Right.Val`**
   — the constraint is on the *whole* subtree, not just direct children.
   Either pass down `(min, max)` bounds, or do an inorder walk and check it is
   strictly increasing.
5. **Confusing depth conventions** — "max depth counts nodes" (empty tree = 0)
   vs. counting edges. And **minimum depth is not `1 + min(left, right)`**
   when one child is nil: a node with a single child is not a leaf, so you
   must recurse into the existing side only.
6. **O(h) space claims that ignore skewed trees** — recursion is O(h), and
   h = n for a linked-list-shaped tree. Constraints allowing n = 10⁴+ skewed
   nodes can blow the stack in some languages; Go's growable goroutine stacks
   are forgiving, but say "O(h), worst-case O(n)" in interviews.
7. **Mutating shared state across recursive branches** — appending a slice
   that later gets mutated (e.g. collecting root-to-leaf paths) requires
   copying the path before storing it (`append([]int{}, path...)`); slices
   share backing arrays. Backtrack (`path = path[:len(path)-1]`) after
   recursing if you reuse one buffer.
8. **Recomputing subtree answers** — a naive Balanced Binary Tree that calls
   `height()` inside `isBalanced()` is O(n²). Postorder returns height and
   balance in a single pass: compute once at the child, reuse at the parent.
9. **Bottom-up value vs. global answer** — for Diameter / Max Path Sum, what
   a node *returns to its parent* (best single downward arm) differs from
   what it *contributes to the answer* (left arm + node + right arm). Keep a
   separate global/closure variable for the answer.
10. **Morris traversal left dirty** — if you exit early (found what you were
    looking for) you may leave threads (`pred.Right = curr`) in the tree.
    Either finish the walk or restore threads before returning.

---

## Problems in this repo

DFS orders and fundamentals:

- [0094 — Binary Tree Inorder Traversal](../0094_binary_tree_inorder_traversal/README.md) — the canonical inorder problem: recursive, iterative, Morris
- [0098 — Validate Binary Search Tree](../0098_validate_binary_search_tree/README.md) — inorder yields sorted order ⇔ valid BST
- [0099 — Recover Binary Search Tree](../0099_recover_binary_search_tree/README.md) — inorder walk to find the two swapped nodes
- [0100 — Same Tree](../0100_same_tree/README.md) — lock-step DFS over two trees
- [0101 — Symmetric Tree](../0101_symmetric_tree/README.md) — lock-step DFS over mirrored halves
- [0104 — Maximum Depth of Binary Tree](../0104_maximum_depth_of_binary_tree/README.md) — the simplest bottom-up postorder
- [0110 — Balanced Binary Tree](../0110_balanced_binary_tree/README.md) — postorder returning height + balance in one pass
- [0111 — Minimum Depth of Binary Tree](../0111_minimum_depth_of_binary_tree/README.md) — leaf-aware depth; BFS early exit
- [0112 — Path Sum](../0112_path_sum/README.md) — top-down preorder with running sum
- [0113 — Path Sum II](../0113_path_sum_ii/README.md) — root-to-leaf paths with backtracking (copy-before-store pitfall)
- [0124 — Binary Tree Maximum Path Sum](../0124_binary_tree_maximum_path_sum/README.md) — postorder; return-arm vs. global-answer distinction
- [0129 — Sum Root to Leaf Numbers](../0129_sum_root_to_leaf_numbers/README.md) — top-down preorder accumulating a number

Level order (BFS) family:

- [0102 — Binary Tree Level Order Traversal](../0102_binary_tree_level_order_traversal/README.md) — the canonical BFS template
- [0103 — Binary Tree Zigzag Level Order Traversal](../0103_binary_tree_zigzag_level_order_traversal/README.md) — alternate level direction
- [0107 — Binary Tree Level Order Traversal II](../0107_binary_tree_level_order_traversal_ii/README.md) — bottom-up levels
- [0116 — Populating Next Right Pointers in Each Node](../0116_populating_next_right_pointers_in_each_node/README.md) — level links; O(1)-space level walk
- [0117 — Populating Next Right Pointers in Each Node II](../0117_populating_next_right_pointers_in_each_node_ii/README.md) — same on an arbitrary tree

Building and reshaping trees via traversal orders:

- [0095 — Unique Binary Search Trees II](../0095_unique_binary_search_trees_ii/README.md) — recursively constructing all BSTs over a range
- [0105 — Construct Binary Tree from Preorder and Inorder Traversal](../0105_construct_binary_tree_from_preorder_and_inorder_traversal/README.md) — preorder gives roots, inorder splits subtrees
- [0106 — Construct Binary Tree from Inorder and Postorder Traversal](../0106_construct_binary_tree_from_inorder_and_postorder_traversal/README.md) — same idea, roots from the back
- [0108 — Convert Sorted Array to Binary Search Tree](../0108_convert_sorted_array_to_binary_search_tree/README.md) — build order = reverse of inorder flattening
- [0109 — Convert Sorted List to Binary Search Tree](../0109_convert_sorted_list_to_binary_search_tree/README.md) — inorder simulation over a linked list
- [0114 — Flatten Binary Tree to Linked List](../0114_flatten_binary_tree_to_linked_list/README.md) — flatten to preorder; reverse-postorder / Morris-style O(1) variants

> Problems 0131+ will be linked in a later pass as they are added to the repo.
