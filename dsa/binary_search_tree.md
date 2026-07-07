# Binary Search Tree (BST)

> **Category:** Data Structure / Tree
> **Prerequisites:** [`binary tree traversal`](/dsa/), recursion, [`binary_search.md`](/dsa/binary_search.md) (same "halve the search space" idea, applied to a tree)

---

## 1. What it is

A **Binary Search Tree** is a binary tree with an **ordering invariant** on every node:

```
For EVERY node n:
    all keys in n.Left  <  n.Val  <  all keys in n.Right
```

Two things interviewers love to test about that definition:

1. **The invariant is global, not local.** It is *not* enough that
   `left child < node < right child`. **Every** node in the entire left
   subtree must be smaller, and **every** node in the entire right subtree
   must be larger. This is the #1 source of wrong `isValidBST` solutions.

   ```
        5
       / \
      3   8
         / \
        4   9     ← 4 < 8 (locally fine) but 4 < 5 → NOT a BST
   ```

2. **Inorder traversal of a BST yields keys in strictly increasing order.**
   (Left → Node → Right visits smaller keys, then the key, then larger keys.)
   This single property is the engine behind most BST problems: validation,
   recovery, kth-smallest, converting to/from sorted sequences, successor
   queries, two-sum on a BST, and more.

### Core operations and their costs

| Operation        | Balanced BST | Degenerate (linked-list-shaped) BST |
|------------------|--------------|-------------------------------------|
| Search           | O(log n)     | O(n)                                |
| Insert           | O(log n)     | O(n)                                |
| Delete           | O(log n)     | O(n)                                |
| Min / Max        | O(log n)     | O(n)                                |
| Inorder (all)    | O(n)         | O(n)                                |
| Predecessor/Succ | O(log n)     | O(n)                                |

`h` (height) drives everything: all point operations are **O(h)**. Inserting
sorted data into a naive BST produces a chain with `h = n` — which is exactly
why self-balancing variants (AVL, Red-Black) and the "build from sorted array
by picking the middle" trick (LeetCode #108) exist.

### BST vs. hash map vs. sorted array

| Need                                   | Best structure |
|----------------------------------------|----------------|
| Exact lookup only                      | Hash map — O(1) avg |
| Lookup **and** ordered iteration       | BST — O(h) lookup, O(n) sorted walk |
| Range queries (`count keys in [a,b]`)  | BST — prune subtrees outside range |
| Predecessor / successor / floor / ceil | BST — O(h) |
| Static data, no inserts                | Sorted array + binary search |

A BST is "a sorted array that supports O(log n) insert/delete" — reach for it
whenever a problem needs **order-aware queries on a changing set**.

---

## 2. How to recognise a BST problem (signals in the statement)

- The phrase **"binary search tree"** or **"BST"** appears — obviously — but
  then ask: *which invariant am I supposed to exploit?* If your solution would
  work on any binary tree, you're probably missing the intended O(h) or
  inorder-based trick.
- **"Sorted"** + **"tree"** in the same problem: sorted array/list → BST
  (#108, #109), BST → sorted output, "two nodes were swapped" in a BST (#99 —
  i.e. two elements out of place in a sorted sequence).
- **kth smallest / largest** element in a tree → inorder walk, stop at k.
- **Range** conditions: "count/sum nodes with values in [lo, hi]" → recurse
  and **prune**: if `node.Val < lo`, the whole left subtree is out.
- **Successor / predecessor / closest value** → walk root-to-leaf, binary
  search style, recording the best candidate on the way down.
- **"How many structurally unique BSTs..."** → Catalan numbers / DP over root
  choices (#95, #96): every value can be the root; smaller values form the
  left subtree, larger values the right.
- **Validate / recover / fix** a BST → inorder must be strictly increasing;
  any violation shows up as an "inversion" `prev.Val >= curr.Val`.
- **Design an ordered map / interval set** (my calendar, snapshot of ranges)
  → in Go, which has no built-in TreeMap, you either implement a BST or keep
  a sorted slice.

---

## 3. General templates (Go)

The node type used throughout (matches LeetCode's `TreeNode`):

```go
type TreeNode struct {
    Val   int
    Left  *TreeNode
    Right *TreeNode
}
```

### 3.1 Search — the fundamental O(h) descent

```go
// searchBST returns the node with value target, or nil.
//
// Pseudocode:
//   while node != nil:
//     if target == node.Val: found
//     if target <  node.Val: go left   (everything right is even bigger)
//     else:                  go right  (everything left is even smaller)
//
// Time: O(h)   Space: O(1) — iterative, no stack.
func searchBST(root *TreeNode, target int) *TreeNode {
    node := root
    for node != nil {
        switch {
        case target == node.Val:
            return node // exact hit
        case target < node.Val:
            node = node.Left // target must live in the left subtree
        default:
            node = node.Right // target must live in the right subtree
        }
    }
    return nil // fell off a leaf — not present
}
```

### 3.2 Insert — descend, attach at the nil where the key belongs

```go
// insertBST inserts val and returns the (possibly new) root.
//
// Time: O(h)   Space: O(h) recursion (O(1) if written iteratively).
func insertBST(root *TreeNode, val int) *TreeNode {
    if root == nil {
        return &TreeNode{Val: val} // empty spot — this is where val lives
    }
    if val < root.Val {
        root.Left = insertBST(root.Left, val) // reattach updated subtree
    } else {
        root.Right = insertBST(root.Right, val)
    }
    return root
}
```

### 3.3 Delete — the three-case classic

```go
// deleteBST removes key and returns the new subtree root.
//
// Cases once the key node is found:
//   1. Leaf            → return nil.
//   2. One child       → splice: return the child.
//   3. Two children    → replace node's value with its inorder successor
//                        (min of right subtree), then delete that successor
//                        from the right subtree (it has at most one child).
//
// Time: O(h)   Space: O(h) recursion.
func deleteBST(root *TreeNode, key int) *TreeNode {
    if root == nil {
        return nil // key not in tree
    }
    switch {
    case key < root.Val:
        root.Left = deleteBST(root.Left, key) // key is left of here
    case key > root.Val:
        root.Right = deleteBST(root.Right, key) // key is right of here
    default: // found the node to delete
        if root.Left == nil {
            return root.Right // covers leaf (nil) and right-only child
        }
        if root.Right == nil {
            return root.Left // left-only child
        }
        // Two children: find inorder successor = leftmost node of right subtree.
        succ := root.Right
        for succ.Left != nil {
            succ = succ.Left
        }
        root.Val = succ.Val                          // adopt successor's value
        root.Right = deleteBST(root.Right, succ.Val) // remove the duplicate
    }
    return root
}
```

### 3.4 Validation — min/max bounds passed down

```go
// isValidBST checks the GLOBAL invariant by narrowing an allowed
// (lo, hi) window as we descend. Root may be anything; going left
// tightens hi, going right tightens lo.
//
// Time: O(n)   Space: O(h).
func isValidBST(root *TreeNode) bool {
    var check func(node *TreeNode, lo, hi *int) bool
    check = func(node *TreeNode, lo, hi *int) bool {
        if node == nil {
            return true // empty subtree is trivially valid
        }
        if lo != nil && node.Val <= *lo { // must be strictly > lower bound
            return false
        }
        if hi != nil && node.Val >= *hi { // must be strictly < upper bound
            return false
        }
        // Left subtree: same lo, hi becomes this node's value.
        // Right subtree: lo becomes this node's value, same hi.
        return check(node.Left, lo, &node.Val) && check(node.Right, &node.Val, hi)
    }
    return check(root, nil, nil) // nil pointers = unbounded (avoids MinInt/MaxInt traps)
}
```

### 3.5 Inorder traversal with a `prev` pointer — the BST workhorse

Validation, recovery (#99), kth smallest, min absolute difference — all are
this template with a different body at the "visit" step:

```go
// inorderWithPrev walks nodes in sorted order, comparing each to the
// previous one. In a valid BST, prev.Val < curr.Val at every step.
//
// Time: O(n)   Space: O(h) for the explicit stack.
func inorderWithPrev(root *TreeNode, visit func(prev, curr *TreeNode)) {
    stack := []*TreeNode{}
    var prev *TreeNode
    curr := root
    for curr != nil || len(stack) > 0 {
        for curr != nil { // dive left as far as possible
            stack = append(stack, curr)
            curr = curr.Left
        }
        curr = stack[len(stack)-1] // smallest unvisited node
        stack = stack[:len(stack)-1]

        visit(prev, curr) // ← problem-specific logic goes here

        prev = curr       // remember for the next comparison
        curr = curr.Right // then handle the right subtree
    }
}
```

### 3.6 Build a balanced BST from a sorted array (#108)

```go
// sortedArrayToBST picks the middle element as root so both halves have
// (nearly) equal size → height O(log n). Recurse on each half.
//
// Time: O(n)   Space: O(log n) recursion.
func sortedArrayToBST(nums []int) *TreeNode {
    if len(nums) == 0 {
        return nil // empty slice → empty subtree
    }
    mid := len(nums) / 2 // middle keeps the tree height-balanced
    return &TreeNode{
        Val:   nums[mid],
        Left:  sortedArrayToBST(nums[:mid]),   // smaller half → left
        Right: sortedArrayToBST(nums[mid+1:]), // larger half → right
    }
}
```

### 3.7 Range pruning — count/sum keys in [lo, hi]

```go
// rangeSumBST sums values in [lo, hi], skipping entire subtrees that
// cannot contain in-range keys.
//
// Time: O(h + k) where k = matches   Space: O(h).
func rangeSumBST(root *TreeNode, lo, hi int) int {
    if root == nil {
        return 0
    }
    if root.Val < lo { // whole left subtree < root.Val < lo → skip it
        return rangeSumBST(root.Right, lo, hi)
    }
    if root.Val > hi { // whole right subtree > root.Val > hi → skip it
        return rangeSumBST(root.Left, lo, hi)
    }
    // Root in range: count it, and both sides may still contain matches.
    return root.Val + rangeSumBST(root.Left, lo, hi) + rangeSumBST(root.Right, lo, hi)
}
```

---

## 4. Worked example — validate a BST (LeetCode #98), traced step by step

Input tree (the classic trap case):

```
        5
       / \
      4   6
         / \
        3   7
```

Locally every parent/child pair looks fine on the right spine (3 < 6, 6 < 7),
but `3` sits in the **right** subtree of `5` while being smaller than 5.

### Trace of the bounds template (`check(node, lo, hi)`)

`nil` bound = unbounded. Notation: `(lo, hi)`.

| Step | Node | Window (lo, hi) | Check `lo < val < hi`      | Action |
|------|------|-----------------|-----------------------------|--------|
| 1    | 5    | (−∞, +∞)        | ok                          | recurse left with (−∞, 5), right with (5, +∞) |
| 2    | 4    | (−∞, 5)         | 4 < 5 ok                    | children are nil → subtree valid |
| 3    | nil  | (−∞, 4)         | —                           | return true |
| 4    | nil  | (4, 5)          | —                           | return true |
| 5    | 6    | (5, +∞)         | 6 > 5 ok                    | recurse left with (5, 6) |
| 6    | **3**| **(5, 6)**      | **3 ≤ 5 → violates lo**     | return **false** — bubbles up |

Result: `false`. Note the failure is caught by the *inherited* lower bound 5
(from the great-grandparent), which a "compare with parent only" solution
never sees.

### Same tree via the inorder-with-prev template

Inorder visit order: `4, 5, 3, 6, 7`.

| Visit | prev | curr | prev.Val < curr.Val? |
|-------|------|------|----------------------|
| 1     | nil  | 4    | (no prev) ok         |
| 2     | 4    | 5    | 4 < 5 ok             |
| 3     | 5    | 3    | **5 < 3 false → not a BST** |

The inorder sequence `4, 5, 3, 6, 7` isn't sorted — the dip at `3` is the
violation. (In #99 Recover BST, the same dips tell you *which two nodes to
swap back*.)

---

## 5. Common pitfalls and how to avoid them

1. **Checking only parent vs. children.** The invariant is subtree-wide.
   Fix: pass down `(lo, hi)` bounds, or verify inorder is strictly increasing.
2. **Using `math.MinInt`/`math.MaxInt` as sentinels** when node values can
   themselves equal those extremes. Fix: use `*int` (nil = unbounded), as in
   template 3.4, or use the `prev`-pointer inorder check.
3. **`<=` vs `<` — duplicates.** LeetCode #98 requires *strictly* less/greater;
   a node equal to an ancestor bound is invalid. Other problems allow
   duplicates (usually "duplicates go right"). Read the constraints, then be
   consistent in both the bound check and the insert direction.
4. **Forgetting to reattach subtrees after recursive insert/delete**
   (`root.Left = insertBST(root.Left, val)` — the assignment matters; the
   recursive call may return a *new* subtree root).
5. **Two-children delete done wrong:** copying the successor's value but
   deleting the wrong node, or picking the successor as "right child" instead
   of "leftmost node of the right subtree". Also fine: use the inorder
   *predecessor* (rightmost of left subtree) — just pick one and delete from
   the correct side.
6. **Assuming O(log n) on a possibly unbalanced tree.** All O(h) claims become
   O(n) on skewed input. If the problem gives sorted input and asks for a BST,
   it wants the balanced middle-element construction (#108/#109), not repeated
   inserts (which would be O(n²) and produce a chain).
7. **Full inorder when early exit suffices.** For kth smallest / first
   violation, use the iterative stack version and `return` as soon as the
   answer is known — recursion makes early termination awkward in Go.
8. **Recursion depth on skewed trees.** h can be n (up to 10⁴–10⁵ on
   LeetCode). Go goroutine stacks grow dynamically so this rarely crashes,
   but mention O(h) stack space in complexity analysis, and know the
   iterative (explicit stack) and Morris-traversal (O(1) space) alternatives.
9. **Recomputing subtree properties per node** (e.g. re-walking subtrees to
   find min/max under every node → O(n²)). Fix: return the needed info
   (min, max, valid?) *up* from recursion, or pass bounds *down*.
10. **Sorted-list → BST via index conversion only.** #109's follow-up wants
    the O(n) inorder-simulation trick (build left subtree first, consume list
    pointer at the "visit" step) instead of converting to an array.

---

## 6. Problems in this repo

Problems 0001–0130 exist today; BST problems in the 0131–0400 range (e.g.
#173, #230, #235) will be linked in a later pass as those folders land.

| Problem | What it exercises |
|---------|-------------------|
| [0094 — Binary Tree Inorder Traversal](../0094_binary_tree_inorder_traversal/README.md) | The traversal that makes BSTs useful: inorder = sorted order; recursive, iterative-stack, and Morris variants |
| [0095 — Unique Binary Search Trees II](../0095_unique_binary_search_trees_ii/README.md) | Enumerate all structurally unique BSTs: every value as root, cartesian product of left/right subtree sets |
| [0096 — Unique Binary Search Trees](../0096_unique_binary_search_trees/README.md) | Count unique BSTs — Catalan numbers / DP over root choices |
| [0098 — Validate Binary Search Tree](../0098_validate_binary_search_tree/README.md) | The global invariant: min/max bounds vs. inorder-with-prev validation |
| [0099 — Recover Binary Search Tree](../0099_recover_binary_search_tree/README.md) | Two swapped nodes = inversions in the inorder sequence; O(1)-space Morris follow-up |
| [0108 — Convert Sorted Array to BST](../0108_convert_sorted_array_to_binary_search_tree/README.md) | Balanced construction: middle element as root, recurse on halves |
| [0109 — Convert Sorted List to BST](../0109_convert_sorted_list_to_binary_search_tree/README.md) | Same idea on a linked list: slow/fast middle, or O(n) inorder simulation |

---

## See also

- [`binary_search.md`](/dsa/binary_search.md) — the array-flavoured sibling of the BST descent
- [`stack.md`](/dsa/stack.md) — explicit-stack iterative traversals
- [`linked_list.md`](/dsa/linked_list.md) — sorted list → BST conversions
