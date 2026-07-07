# 0106 — Construct Binary Tree from Inorder and Postorder Traversal

> LeetCode #106 · Difficulty: Medium
> **Categories:** Array, Hash Table, Divide and Conquer, Tree, Binary Tree

---

## Problem Statement

Given two integer arrays `inorder` and `postorder` where `inorder` is the inorder traversal of a binary tree and `postorder` is the postorder traversal of the same tree, construct and return the binary tree.

**Example 1:**
```
Input: inorder = [9,3,15,20,7], postorder = [9,15,7,20,3]
Output: [3,9,20,null,null,15,7]
```

**Example 2:**
```
Input: inorder = [-1], postorder = [-1]
Output: [-1]
```

**Constraints:**
- `1 <= inorder.length <= 3000`
- `inorder.length == postorder.length`
- `-3000 <= inorder[i], postorder[i] <= 3000`
- All values are unique.
- `inorder` is the inorder traversal of the tree.
- `postorder` is the postorder traversal of the tree.

---

## Company Frequency

| Company   | Frequency    | Last Reported |
|-----------|--------------|---------------|
| Amazon    | ★★★★☆ High  | 2024          |
| Google    | ★★★★☆ High  | 2024          |
| Microsoft | ★★★☆☆ Medium | 2023          |

> ⚠️ Frequency data is crowd-sourced from LeetCode Discuss, Glassdoor, and community interview reports. Treat as a signal, not a guarantee.

---

## DSA Concepts Used

- **Divide and Conquer** — split into left/right subtrees using root's inorder position
- **Hash Map** — O(1) lookup of root index in inorder array
- **Tree Construction** — mirror of #105 but using postorder's last element as root

---

## Approaches Overview

| # | Approach                  | Time | Space | When to use          |
|---|---------------------------|------|-------|----------------------|
| 1 | Recursive + HashMap       | O(n) | O(n)  | Optimal; standard    |
| 2 | Iterative Stack (reversed)| O(n) | O(h)  | Avoids recursion     |

---

## Approach 1 — Recursive with HashMap

### Intuition
Mirror of #105. Key differences:
- **Root** = `postorder[last]` (not `preorder[0]`).
- **Right** subtree uses postorder elements before the root; **left** before that.
- `rightSize = inEnd - rootIdx` elements go to the right subtree.

### Algorithm
1. Build `inMap[val] = index` for all inorder values.
2. `build(postStart, postEnd, inStart, inEnd)`:
   - Root = `postorder[postEnd]`.
   - `rootIdx = inMap[rootVal]`, `rightSize = inEnd - rootIdx`.
   - `root.Right = build(postEnd-rightSize, postEnd-1, rootIdx+1, inEnd)`.
   - `root.Left  = build(postStart, postEnd-rightSize-1, inStart, rootIdx-1)`.

### Complexity
- **Time:** O(n)
- **Space:** O(n) — hashmap; O(h) recursion.

### Code
```go
func buildTree(inorder []int, postorder []int) *TreeNode {
    inMap := make(map[int]int, len(inorder))
    for i, v := range inorder { inMap[v] = i }

    var build func(postStart, postEnd, inStart, inEnd int) *TreeNode
    build = func(postStart, postEnd, inStart, inEnd int) *TreeNode {
        if postStart > postEnd { return nil }
        rootVal := postorder[postEnd]
        root := &TreeNode{Val: rootVal}
        rootIdx := inMap[rootVal]
        rightSize := inEnd - rootIdx
        root.Right = build(postEnd-rightSize, postEnd-1, rootIdx+1, inEnd)
        root.Left  = build(postStart, postEnd-rightSize-1, inStart, rootIdx-1)
        return root
    }
    return build(0, len(postorder)-1, 0, len(inorder)-1)
}
```

### Dry Run
`inorder=[9,3,15,20,7]`, `postorder=[9,15,7,20,3]`, `inMap={9:0,3:1,15:2,20:3,7:4}`

| Call              | rootVal | rootIdx | rightSize | Right call         | Left call          |
|-------------------|---------|---------|-----------|--------------------|--------------------|
| build(0,4,0,4)    | 3       | 1       | 3         | build(2,4,2,4)     | build(0,0,0,0)     |
| build(0,0,0,0)    | 9       | 0       | 0         | nil                | nil                |
| build(2,4,2,4)    | 20      | 3       | 1         | build(4,3,4,4)=nil | build(2,2,2,2)     |

Wait — `build(4,3,…)` has `postStart>postEnd` → nil. `build(2,2,2,2)` → root=15.

Tree: 3(left=9, right=20(left=15, right=7)) ✓

---

## Approach 2 — Iterative Stack (Reverse Postorder)

### Intuition
Reversed postorder = root → right → left. Apply the same stack trick as #105 iterative, but:
- Walk `postorder` from right to left.
- Walk `inorder` from right to left.
- When match found, attach as **left** child (mirror of right in #105).

### Complexity
- **Time:** O(n)
- **Space:** O(h)

### Code
```go
func buildTreeIterative(inorder []int, postorder []int) *TreeNode {
    n := len(postorder)
    if n == 0 { return nil }
    root := &TreeNode{Val: postorder[n-1]}
    stack := []*TreeNode{root}
    inIdx := n - 1
    for i := n - 2; i >= 0; i-- {
        node := &TreeNode{Val: postorder[i]}
        if stack[len(stack)-1].Val != inorder[inIdx] {
            stack[len(stack)-1].Right = node
        } else {
            var parent *TreeNode
            for len(stack) > 0 && stack[len(stack)-1].Val == inorder[inIdx] {
                parent = stack[len(stack)-1]
                stack = stack[:len(stack)-1]; inIdx--
            }
            parent.Left = node
        }
        stack = append(stack, node)
    }
    return root
}
```

### Dry Run
`postorder=[9,15,7,20,3]`, `inorder=[9,3,15,20,7]` (walk both from right).

Reversed postorder walk: 3,20,7,15,9.
- 3: root, stack=[3], inIdx=4.
- 20: top=3, inorder[4]=7. 3≠7 → 3.Right=20. stack=[3,20].
- 7: top=20, inorder[4]=7. 20≠7 → 20.Right=7. stack=[3,20,7].
- 15: top=7, inorder[4]=7. match: pop 7,inIdx=3. inorder[3]=20, top=20, match: pop 20,inIdx=2. inorder[2]=15, top=3≠15 → parent=20. 20.Left=15. stack=[3,15].
- 9: top=15, inorder[2]=15. match: pop 15,inIdx=1. inorder[1]=3, top=3, match: pop 3,inIdx=0. stack=[]. parent=3. 3.Left=9.

Tree: 3(left=9, right=20(left=15,right=7)) ✓

---

## Key Takeaways
- Postorder root = last element; preorder root = first element.
- `rightSize = inEnd - rootIdx` (not `leftSize`) since right elements come last in postorder.
- Iterative trick: reversed postorder is root→right→left, exactly mirroring #105's preorder=root→left→right.

---

## Related Problems
- LeetCode #105 — Construct Binary Tree from Preorder and Inorder Traversal
- LeetCode #889 — Construct Binary Tree from Preorder and Postorder Traversal
