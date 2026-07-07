# 0117 — Populating Next Right Pointers in Each Node II

> LeetCode #117 · Difficulty: Medium
> **Categories:** Linked List, Tree, Depth-First Search, Breadth-First Search, Binary Tree

---

## Problem Statement

Given a binary tree (not necessarily perfect), populate each `next` pointer to point to its next right node. If there is no next right node, set it to `NULL`.

**Example 1:**
```
Input: root = [1,2,3,4,5,null,7]
Output: [1,#,2,3,#,4,5,7,#]
```

**Example 2:**
```
Input: root = []
Output: []
```

**Constraints:**
- The number of nodes in the tree is in the range `[0, 6000]`.
- `-100 <= node.val <= 100`

**Follow-up:** Use only O(1) extra space.

---

## Company Frequency

| Company   | Frequency    | Last Reported |
|-----------|--------------|---------------|
| Amazon    | ★★★★☆ High  | 2024          |
| Microsoft | ★★★★☆ High  | 2024          |
| Google    | ★★★☆☆ Medium | 2023          |

> ⚠️ Frequency data is crowd-sourced from LeetCode Discuss, Glassdoor, and community interview reports. Treat as a signal, not a guarantee.

---

## DSA Concepts Used

- **BFS** — handles any tree structure with a queue
- **Dummy head linked list** — simplifies tracking first child of next level

---

## Approaches Overview

| # | Approach               | Time | Space | When to use             |
|---|------------------------|------|-------|-------------------------|
| 1 | BFS Queue              | O(n) | O(w)  | General; simple         |
| 2 | O(1) Space + Dummy Head| O(n) | O(1)  | Follow-up requirement   |

---

## Approach 1 — BFS Level Order

### Intuition
Identical to #116 BFS — works for any tree, not just perfect ones.

### Complexity
- **Time:** O(n)
- **Space:** O(w)

### Code
```go
func connect(root *Node) *Node {
    if root == nil { return nil }
    queue := []*Node{root}
    for len(queue) > 0 {
        levelSize := len(queue)
        for i := 0; i < levelSize; i++ {
            node := queue[0]; queue = queue[1:]
            if i < levelSize-1 { node.Next = queue[0] }
            if node.Left  != nil { queue = append(queue, node.Left)  }
            if node.Right != nil { queue = append(queue, node.Right) }
        }
    }
    return root
}
```

### Dry Run
Tree `[1,2,3,4,5,null,7]`, Level [2,3]:
- node=2: `2.Next=3`. Enqueue 4, 5.
- node=3: last. Enqueue 7 (no 3.Left).

Level [4,5,7]:
- node=4: `4.Next=5`.
- node=5: `5.Next=7`.
- node=7: last.

---

## Approach 2 — O(1) Space with Dummy Head

### Intuition
Can't use the perfect-tree trick (#116) since nodes may not have both children. Instead:
- Traverse the current level via `curr` and its `Next` pointers.
- Build the next level's linked list using a `dummy` head and `tail` pointer.
- After processing the level, descend: `curr = dummy.Next`.

The dummy head avoids a special case for "first child of level."

### Algorithm
```
curr = root
while curr != nil:
    dummy = new Node()
    tail = dummy
    while curr != nil:
        if curr.Left: tail.Next = curr.Left; tail = tail.Next
        if curr.Right: tail.Next = curr.Right; tail = tail.Next
        curr = curr.Next
    curr = dummy.Next
```

### Complexity
- **Time:** O(n)
- **Space:** O(1)

### Code
```go
func connectO1(root *Node) *Node {
    curr := root
    for curr != nil {
        dummy := &Node{}; tail := dummy
        for curr != nil {
            if curr.Left  != nil { tail.Next = curr.Left;  tail = tail.Next }
            if curr.Right != nil { tail.Next = curr.Right; tail = tail.Next }
            curr = curr.Next
        }
        curr = dummy.Next
    }
    return root
}
```

### Dry Run
Level [2→3], curr=2:
- curr=2: tail→4, tail→5. curr=2.Next=3.
- curr=3: skip left (nil), tail→7. curr=3.Next=nil.
- Exit inner loop. dummy.Next=4. 4.Next=5. 5.Next=7.
- curr = 4.

Level [4→5→7], curr=4:
- No children → tail stays at dummy. curr=dummy.Next=nil. Done.

---

## Key Takeaways
- Dummy head pattern makes "append to next level list" uniform — no nil check for first node.
- The outer while loop iterates one level per iteration using `curr` which traverses via `Next`.
- Works for any binary tree (arbitrary children), unlike #116's O(1) method.

---

## Related Problems
- LeetCode #116 — Populating Next Right Pointers (perfect binary tree only)
- LeetCode #102 — Binary Tree Level Order Traversal
