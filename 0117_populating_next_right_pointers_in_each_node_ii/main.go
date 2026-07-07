package main

import "fmt"

// Node is a binary tree node with a next pointer.
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
		// find leftmost child using next pointers
		next := root
		for next != nil {
			if next.Left != nil { root = next.Left; break }
			if next.Right != nil { root = next.Right; break }
			next = next.Next
		}
		if next == nil { break }
	}
	return result
}

// ── Approach 1: BFS Level Order ───────────────────────────────────────────────
//
// connect solves Populating Next Right Pointers II (arbitrary binary tree) using BFS.
//
// Intuition:
//   Same as #116 BFS approach — works for any binary tree, not just perfect ones.
//
// Time:  O(n)
// Space: O(w)
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
				node.Next = queue[0]
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

// ── Approach 2: O(1) Space with Dummy Head ───────────────────────────────────
//
// connectO1 solves Populating Next Right Pointers II using O(1) space.
//
// Intuition:
//   Use the current level's Next pointers (already set) to traverse it.
//   For each node, wire its children to a dummy-headed "next level" linked list.
//   A dummy head simplifies "first child of level" tracking.
//
// Algorithm:
//   For each level:
//     Use `curr` to scan across the level via Next pointers.
//     Use a `dummy` node and `tail` to build the next level's linked list.
//     Advance `curr = curr.Next` to next node in current level.
//     When done, advance to `dummy.Next` (first node of next level).
//
// Time:  O(n)
// Space: O(1)
func connectO1(root *Node) *Node {
	curr := root

	for curr != nil {
		dummy := &Node{} // head of next level's linked list
		tail := dummy

		for curr != nil {
			if curr.Left != nil {
				tail.Next = curr.Left
				tail = tail.Next
			}
			if curr.Right != nil {
				tail.Next = curr.Right
				tail = tail.Next
			}
			curr = curr.Next // advance across current level
		}
		curr = dummy.Next // descend to next level
	}
	return root
}

func main() {
	// [1,2,3,4,5,null,7]
	build := func() *Node {
		return &Node{Val: 1,
			Left: &Node{Val: 2,
				Left: &Node{Val: 4}, Right: &Node{Val: 5}},
			Right: &Node{Val: 3,
				Right: &Node{Val: 7}},
		}
	}

	fmt.Println("=== Approach 1: BFS ===")
	t1 := build()
	connect(t1)
	fmt.Printf("levels=%v  expected [[1] [2 3] [4 5 7]]\n", printLevels(t1))

	fmt.Println("=== Approach 2: O(1) Space + Dummy Head ===")
	t2 := build()
	connectO1(t2)
	fmt.Printf("levels=%v  expected [[1] [2 3] [4 5 7]]\n", printLevels(t2))
}
