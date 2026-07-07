package main

import "fmt"

// Node is a perfect binary tree node with a next pointer.
type Node struct {
	Val   int
	Left  *Node
	Right *Node
	Next  *Node
}

func printLevels(root *Node) [][]int {
	var result [][]int
	for root != nil {
		var level []int
		curr := root
		for curr != nil {
			level = append(level, curr.Val)
			curr = curr.Next
		}
		result = append(result, level)
		root = root.Left
	}
	return result
}

// ── Approach 1: BFS Level Order ───────────────────────────────────────────────
//
// connect solves Populating Next Right Pointers in Each Node using BFS.
//
// Intuition:
//   Process level by level. Within each level, wire node.Next = next node in queue.
//   Last node in each level gets Next = nil (default).
//
// Time:  O(n)
// Space: O(w) — queue holds at most one level.
func connect(root *Node) *Node {
	if root == nil {
		return nil
	}
	queue := []*Node{root}

	for len(queue) > 0 {
		levelSize := len(queue)
		for i := 0; i < levelSize; i++ {
			node := queue[0]
			queue = queue[1:]
			if i < levelSize-1 {
				node.Next = queue[0] // point to next in same level
			}
			if node.Left != nil {
				queue = append(queue, node.Left)
			}
			if node.Right != nil {
				queue = append(queue, node.Right)
			}
		}
	}
	return root
}

// ── Approach 2: O(1) Space — Exploit Perfect Tree Structure ──────────────────
//
// connectO1 solves Populating Next Right Pointers using the already-set Next
// pointers of the previous level to traverse without a queue.
//
// Intuition:
//   For a perfect binary tree, at each level we can use the Next pointers
//   (already set) to traverse all nodes and wire their children's Next pointers.
//   Two cases:
//   - Same parent: node.Left.Next = node.Right
//   - Different parent: node.Right.Next = node.Next.Left (if node.Next exists)
//
// Time:  O(n)
// Space: O(1)
func connectO1(root *Node) *Node {
	if root == nil {
		return nil
	}
	leftmost := root // start of each level

	for leftmost.Left != nil { // stop at leaf level
		curr := leftmost
		for curr != nil {
			curr.Left.Next = curr.Right // same-parent connection
			if curr.Next != nil {
				curr.Right.Next = curr.Next.Left // cross-parent connection
			}
			curr = curr.Next // traverse level using already-set Next pointers
		}
		leftmost = leftmost.Left // move to next level
	}
	return root
}

func main() {
	build := func() *Node {
		return &Node{Val: 1,
			Left: &Node{Val: 2,
				Left: &Node{Val: 4}, Right: &Node{Val: 5}},
			Right: &Node{Val: 3,
				Left: &Node{Val: 6}, Right: &Node{Val: 7}},
		}
	}

	fmt.Println("=== Approach 1: BFS ===")
	t1 := build()
	connect(t1)
	fmt.Printf("levels=%v  expected [[1] [2 3] [4 5 6 7]]\n", printLevels(t1))

	fmt.Println("=== Approach 2: O(1) Space ===")
	t2 := build()
	connectO1(t2)
	fmt.Printf("levels=%v  expected [[1] [2 3] [4 5 6 7]]\n", printLevels(t2))

	// nil case
	fmt.Printf("nil case: %v\n", connect(nil))
}
