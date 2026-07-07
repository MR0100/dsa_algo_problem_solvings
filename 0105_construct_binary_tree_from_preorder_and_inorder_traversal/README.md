# 0105 — Construct Binary Tree from Preorder and Inorder Traversal

> LeetCode #105 · Difficulty: Medium
> **Categories:** Array, Hash Table, Divide and Conquer, Tree, Binary Tree

---

## Problem Statement

Given two integer arrays `preorder` and `inorder` where `preorder` is the preorder traversal of a binary tree and `inorder` is the inorder traversal of the same tree, construct and return the binary tree.

**Example 1:**
```
Input: preorder = [3,9,20,15,7], inorder = [9,3,15,20,7]
Output: [3,9,20,null,null,15,7]
```

**Example 2:**
```
Input: preorder = [-1], inorder = [-1]
Output: [-1]
```

**Constraints:**
- `1 <= preorder.length <= 3000`
- `inorder.length == preorder.length`
- `-3000 <= preorder[i], inorder[i] <= 3000`
- `preorder` and `inorder` consist of unique values.
- Each value of `inorder` also appears in `preorder`.
- `preorder` is guaranteed to be the preorder traversal.
- `inorder` is guaranteed to be the inorder traversal.

---

## Company Frequency

| Company   | Frequency       | Last Reported |
|-----------|-----------------|---------------|
| Amazon    | ★★★★★ Very High | 2024          |
| Google    | ★★★★★ Very High | 2024          |
| Microsoft | ★★★★☆ High      | 2024          |
| Bloomberg | ★★★★☆ High      | 2023          |
| Facebook  | ★★★☆☆ Medium    | 2024          |

> ⚠️ Frequency data is crowd-sourced from LeetCode Discuss, Glassdoor, and community interview reports. Treat as a signal, not a guarantee.

---

## DSA Concepts Used

- **Divide and Conquer** — split problem into left/right subtrees
- **Hash Map** — O(1) lookup of root index in inorder array
- **Tree Construction** — linking nodes via recursion

---

## Approaches Overview

| # | Approach                  | Time | Space | When to use               |
|---|---------------------------|------|-------|---------------------------|
| 1 | Recursive + HashMap       | O(n) | O(n)  | Optimal; most common      |
| 2 | Iterative Stack           | O(n) | O(h)  | Avoids recursion          |

---

## Approach 1 — Recursive with HashMap

### Intuition
Key observations:
- `preorder[0]` is always the root of the current subtree.
- In `inorder`, elements to the left of the root form the left subtree; elements to the right form the right subtree.
- `leftSize = rootInorderIdx - inStart` tells us how many elements are in the left subtree, and therefore how many preorder elements to allocate to it.

Build a hashmap of `value → inorder index` upfront to avoid O(n) scans.

### Algorithm
1. Build `inMap`: `value → index` for all of `inorder`.
2. Define `build(preStart, preEnd, inStart, inEnd)`:
   - If `preStart > preEnd`, return nil.
   - `rootVal = preorder[preStart]`, create root node.
   - `rootIdx = inMap[rootVal]`, `leftSize = rootIdx - inStart`.
   - `root.Left = build(preStart+1, preStart+leftSize, inStart, rootIdx-1)`
   - `root.Right = build(preStart+leftSize+1, preEnd, rootIdx+1, inEnd)`
3. Return `build(0, n-1, 0, n-1)`.

### Complexity
- **Time:** O(n) — each node visited once; hashmap lookups O(1).
- **Space:** O(n) — hashmap of size n; recursion stack O(h).

### Code
```go
func buildTree(preorder []int, inorder []int) *TreeNode {
    inMap := make(map[int]int, len(inorder))
    for i, v := range inorder { inMap[v] = i }

    var build func(preStart, preEnd, inStart, inEnd int) *TreeNode
    build = func(preStart, preEnd, inStart, inEnd int) *TreeNode {
        if preStart > preEnd { return nil }
        rootVal := preorder[preStart]
        root := &TreeNode{Val: rootVal}
        rootIdx := inMap[rootVal]
        leftSize := rootIdx - inStart
        root.Left  = build(preStart+1, preStart+leftSize, inStart, rootIdx-1)
        root.Right = build(preStart+leftSize+1, preEnd, rootIdx+1, inEnd)
        return root
    }
    return build(0, len(preorder)-1, 0, len(inorder)-1)
}
```

### Dry Run
`preorder = [3,9,20,15,7]`, `inorder = [9,3,15,20,7]`

`inMap = {9:0, 3:1, 15:2, 20:3, 7:4}`

| Call                          | rootVal | rootIdx | leftSize | Left call       | Right call     |
|-------------------------------|---------|---------|----------|-----------------|----------------|
| build(0,4,0,4)                | 3       | 1       | 1        | build(1,1,0,0)  | build(2,4,2,4) |
| build(1,1,0,0)                | 9       | 0       | 0        | build(1,0,0,-1)=nil | build(1,1,1,0)=nil |
| build(2,4,2,4)                | 20      | 3       | 1        | build(3,3,2,2)  | build(4,4,4,4) |
| build(3,3,2,2)                | 15      | 2       | 0        | nil             | nil            |
| build(4,4,4,4)                | 7       | 4       | 0        | nil             | nil            |

Tree built: root=3, left=9, right=20, 20.left=15, 20.right=7.

---

## Approach 2 — Iterative Stack

### Intuition
Walk preorder left-to-right. Maintain a stack of nodes. For each preorder element:
- If the top of stack doesn't match `inorder[inIdx]`, the current node is a left child of the stack top.
- If it does match, we've finished the left subtree. Pop until mismatch — the last popped node needs the current node as its right child.

### Algorithm
1. Create root from `preorder[0]`, push onto stack. `inIdx = 0`.
2. For `i = 1..n-1`:
   - Create `node = preorder[i]`.
   - If `stack.top.Val != inorder[inIdx]`: set `stack.top.Left = node`.
   - Else: pop while match, set `lastPopped.Right = node`.
   - Push `node`.
3. Return root.

### Complexity
- **Time:** O(n) — each node pushed/popped once.
- **Space:** O(h) — stack size.

### Code
```go
func buildTreeIterative(preorder []int, inorder []int) *TreeNode {
    if len(preorder) == 0 { return nil }
    root := &TreeNode{Val: preorder[0]}
    stack := []*TreeNode{root}
    inIdx := 0
    for i := 1; i < len(preorder); i++ {
        node := &TreeNode{Val: preorder[i]}
        if stack[len(stack)-1].Val != inorder[inIdx] {
            stack[len(stack)-1].Left = node
        } else {
            var parent *TreeNode
            for len(stack) > 0 && stack[len(stack)-1].Val == inorder[inIdx] {
                parent = stack[len(stack)-1]
                stack = stack[:len(stack)-1]; inIdx++
            }
            parent.Right = node
        }
        stack = append(stack, node)
    }
    return root
}
```

### Dry Run
`preorder = [3,9,20,15,7]`, `inorder = [9,3,15,20,7]`, `inIdx=0`

| i | preorder[i] | stack top | inorder[inIdx] | Action                  |
|---|-------------|-----------|----------------|-------------------------|
| 1 | 9           | 3         | 9              | top(3)≠9? No. Pop 3,inIdx=1; but wait inorder[0]=9≠3.Val? Actually 3≠9→left. |

Trace corrected:
- `stack=[3]`, `inIdx=0`, `inorder[0]=9`. `stack.top(3).Val=3 != 9` → `3.Left=9`. Push 9. `stack=[3,9]`.
- `i=2`, node=20. `stack.top=9`, `inorder[0]=9`. `9==9` → pop: parent=9,inIdx=1. `inorder[1]=3`, `stack.top=3==3` → pop: parent=3,inIdx=2. `stack=[]`. `parent(3).Right=20`. Push 20. `stack=[20]`.
- `i=3`, node=15. `stack.top=20`, `inorder[2]=15`. `20≠15` → `20.Left=15`. Push 15.
- `i=4`, node=7. `stack.top=15`, `inorder[2]=15`. `15==15` → pop: parent=15,inIdx=3. `inorder[3]=20`, `stack.top=20==20` → pop: parent=20,inIdx=4. `parent(20).Right=7`.

Result: 3(left=9, right=20(left=15,right=7)) ✓

---

## Key Takeaways
- `preorder[0]` = root; inorder position of root divides left/right subtrees.
- `leftSize = rootIdx - inStart` → number of preorder elements for left subtree.
- Always build an inorder hashmap to avoid O(n) scans, making the solution O(n) not O(n²).
- Iterative approach uses a stack and an inorder pointer — tricky but O(h) space without recursion.

---

## Related Problems
- LeetCode #106 — Construct Binary Tree from Inorder and Postorder Traversal
- LeetCode #889 — Construct Binary Tree from Preorder and Postorder Traversal
- LeetCode #108 — Convert Sorted Array to Binary Search Tree
