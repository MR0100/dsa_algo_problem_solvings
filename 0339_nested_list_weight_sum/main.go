package main

import "fmt"

// NestedInteger models LeetCode's opaque NestedInteger interface: each element
// is EITHER a single integer OR a list of NestedIntegers. We model it as a
// struct where isInt distinguishes the two cases.
type NestedInteger struct {
	isInt bool             // true → this holds a single integer
	value int              // valid only when isInt == true
	list  []*NestedInteger // valid only when isInt == false
}

// IsInteger reports whether this NestedInteger holds a single integer.
func (n *NestedInteger) IsInteger() bool { return n.isInt }

// GetInteger returns the single integer (only meaningful if IsInteger()).
func (n *NestedInteger) GetInteger() int { return n.value }

// GetList returns the nested list (only meaningful if !IsInteger()).
func (n *NestedInteger) GetList() []*NestedInteger { return n.list }

// leaf builds a NestedInteger wrapping a single integer.
func leaf(v int) *NestedInteger { return &NestedInteger{isInt: true, value: v} }

// nest builds a NestedInteger wrapping a list.
func nest(items ...*NestedInteger) *NestedInteger {
	return &NestedInteger{isInt: false, list: items}
}

// ── Approach 1: DFS (recursive, depth carried down) ──────────────────────────
//
// dfs solves Nested List Weight Sum by recursing into the structure, passing
// the current depth down and accumulating value*depth at each integer.
//
// Intuition:
//
//	Weight = depth (root list is depth 1). Walk the tree of nested integers.
//	When you meet a plain integer at depth d, it contributes value*d. When you
//	meet a list, recurse into it at depth d+1. Summing all contributions gives
//	the answer.
//
// Algorithm:
//  1. helper(list, depth): for each element, if integer add value*depth,
//     else recurse with depth+1.
//  2. Start with helper(input, 1).
//
// Time:  O(N) — N = total number of integers and lists visited once.
// Space: O(D) — recursion stack, D = maximum nesting depth.
func dfs(nestedList []*NestedInteger) int {
	var helper func(list []*NestedInteger, depth int) int
	helper = func(list []*NestedInteger, depth int) int {
		sum := 0
		for _, ni := range list {
			if ni.IsInteger() {
				sum += ni.GetInteger() * depth // leaf contributes value×depth
			} else {
				sum += helper(ni.GetList(), depth+1) // descend: deeper by one level
			}
		}
		return sum
	}
	return helper(nestedList, 1) // top-level list is depth 1
}

// ── Approach 2: BFS (level by level) ─────────────────────────────────────────
//
// bfs solves it by processing the structure level by level with a queue,
// tracking the current depth explicitly.
//
// Intuition:
//
//	Instead of recursion, hold a queue of NestedIntegers for the current level.
//	At depth d, sum the integers and enqueue the contents of any lists for the
//	next level (depth d+1). Repeat until the queue empties.
//
// Algorithm:
//  1. queue = top-level items, depth = 1.
//  2. While queue non-empty: pop all current-level items; integers add
//     value*depth, lists enqueue their children; then depth++.
//
// Time:  O(N) — each element enqueued/dequeued once.
// Space: O(W) — W = maximum number of elements on a single level.
func bfs(nestedList []*NestedInteger) int {
	queue := nestedList // items pending at the current depth
	depth := 1          // top level is depth 1
	sum := 0
	for len(queue) > 0 {
		var next []*NestedInteger // items collected for the following depth
		for _, ni := range queue {
			if ni.IsInteger() {
				sum += ni.GetInteger() * depth // weigh this integer by its depth
			} else {
				next = append(next, ni.GetList()...) // its children live one level deeper
			}
		}
		queue = next // advance to the next level
		depth++      // and increase the weight
	}
	return sum
}

// ── Approach 3: DFS Explicit Stack (Optimal, no recursion) ───────────────────
//
// iterativeStack solves it with an explicit stack of (NestedInteger, depth)
// frames, avoiding recursion entirely.
//
// Intuition:
//
//	Simulate the DFS with a manual stack. Each frame pairs an element with its
//	depth. Pop a frame: if it's an integer, add value*depth; if it's a list,
//	push each child paired with depth+1. Purely iterative, same O(N) work.
//
// Algorithm:
//  1. Push all top-level items with depth 1.
//  2. Pop frame; integer → add value*depth; list → push children at depth+1.
//  3. Continue until the stack empties.
//
// Time:  O(N) — every element pushed and popped once.
// Space: O(N) worst case — the stack can hold all elements of a wide level.
func iterativeStack(nestedList []*NestedInteger) int {
	// frame pairs an element with the depth at which it sits.
	type frame struct {
		ni    *NestedInteger
		depth int
	}
	stack := make([]frame, 0, len(nestedList))
	for _, ni := range nestedList {
		stack = append(stack, frame{ni, 1}) // seed with top-level, depth 1
	}
	sum := 0
	for len(stack) > 0 {
		f := stack[len(stack)-1] // pop the top frame
		stack = stack[:len(stack)-1]
		if f.ni.IsInteger() {
			sum += f.ni.GetInteger() * f.depth // weigh by its depth
		} else {
			for _, child := range f.ni.GetList() {
				stack = append(stack, frame{child, f.depth + 1}) // children go deeper
			}
		}
	}
	return sum
}

func main() {
	// Example 1: [[1,1],2,[1,1]] → 10
	//   (1+1)*2  + 2*1 + (1+1)*2 = 4 + 2 + 4 = 10
	ex1 := []*NestedInteger{
		nest(leaf(1), leaf(1)),
		leaf(2),
		nest(leaf(1), leaf(1)),
	}
	// Example 2: [1,[4,[6]]] → 27
	//   1*1 + 4*2 + 6*3 = 1 + 8 + 18 = 27
	ex2 := []*NestedInteger{
		leaf(1),
		nest(leaf(4), nest(leaf(6))),
	}

	fmt.Println("=== Approach 1: DFS (recursive) ===")
	fmt.Println(dfs(ex1)) // 10
	fmt.Println(dfs(ex2)) // 27

	fmt.Println("=== Approach 2: BFS (level by level) ===")
	fmt.Println(bfs(ex1)) // 10
	fmt.Println(bfs(ex2)) // 27

	fmt.Println("=== Approach 3: DFS Explicit Stack (Optimal) ===")
	fmt.Println(iterativeStack(ex1)) // 10
	fmt.Println(iterativeStack(ex2)) // 27
}
