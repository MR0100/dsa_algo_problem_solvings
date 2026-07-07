package main

import "fmt"

// TreeNode is a binary tree node.
type TreeNode struct {
	Val   int
	Left  *TreeNode
	Right *TreeNode
}

// ── Approach 1: DFS with Running Number ──────────────────────────────────────
//
// sumNumbers solves Sum Root to Leaf Numbers using DFS.
//
// Intuition:
//   As we descend, maintain the running number formed so far:
//   running = running*10 + node.Val.
//   At a leaf, add running to the total.
//
// Time:  O(n)
// Space: O(h)
func sumNumbers(root *TreeNode) int {
	var dfs func(node *TreeNode, running int) int
	dfs = func(node *TreeNode, running int) int {
		if node == nil {
			return 0
		}
		running = running*10 + node.Val
		if node.Left == nil && node.Right == nil {
			return running // leaf: full number formed
		}
		return dfs(node.Left, running) + dfs(node.Right, running)
	}
	return dfs(root, 0)
}

// ── Approach 2: Iterative DFS (Stack) ────────────────────────────────────────
//
// sumNumbersIterative solves Sum Root to Leaf Numbers iteratively.
//
// Intuition:
//   Stack stores (node, running) pairs. Same logic as recursive: multiply by 10
//   and add digit. At leaf, add to total.
//
// Time:  O(n)
// Space: O(h)
func sumNumbersIterative(root *TreeNode) int {
	if root == nil {
		return 0
	}
	type item struct {
		node    *TreeNode
		running int
	}
	stack := []item{{root, 0}}
	total := 0

	for len(stack) > 0 {
		curr := stack[len(stack)-1]
		stack = stack[:len(stack)-1]
		running := curr.running*10 + curr.node.Val
		if curr.node.Left == nil && curr.node.Right == nil {
			total += running
			continue
		}
		if curr.node.Left != nil {
			stack = append(stack, item{curr.node.Left, running})
		}
		if curr.node.Right != nil {
			stack = append(stack, item{curr.node.Right, running})
		}
	}
	return total
}

func main() {
	// [1,2,3] → 12 + 13 = 25
	t1 := &TreeNode{Val: 1, Left: &TreeNode{Val: 2}, Right: &TreeNode{Val: 3}}
	// [4,9,0,5,1] → 495 + 491 + 40 = 1026
	t2 := &TreeNode{Val: 4,
		Left: &TreeNode{Val: 9,
			Left: &TreeNode{Val: 5}, Right: &TreeNode{Val: 1}},
		Right: &TreeNode{Val: 0},
	}

	fmt.Println("=== Approach 1: Recursive DFS ===")
	fmt.Printf("tree=[1,2,3]  got=%d  expected 25\n", sumNumbers(t1))
	fmt.Printf("tree=[4,9,0,5,1]  got=%d  expected 1026\n", sumNumbers(t2))

	t3 := &TreeNode{Val: 1, Left: &TreeNode{Val: 2}, Right: &TreeNode{Val: 3}}
	t4 := &TreeNode{Val: 4,
		Left: &TreeNode{Val: 9,
			Left: &TreeNode{Val: 5}, Right: &TreeNode{Val: 1}},
		Right: &TreeNode{Val: 0},
	}

	fmt.Println("=== Approach 2: Iterative DFS ===")
	fmt.Printf("tree=[1,2,3]  got=%d  expected 25\n", sumNumbersIterative(t3))
	fmt.Printf("tree=[4,9,0,5,1]  got=%d  expected 1026\n", sumNumbersIterative(t4))
}
