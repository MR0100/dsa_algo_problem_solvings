package main

import "fmt"

// NestedInteger models LeetCode's NestedInteger interface: each element is
// either a single integer OR a list of NestedInteger. We tag which one it is.
type NestedInteger struct {
	isInt bool             // true → this holds a plain integer
	value int              // valid when isInt is true
	list  []*NestedInteger // valid when isInt is false (a nested list)
}

// NewInt builds a leaf integer node.
func NewInt(v int) *NestedInteger { return &NestedInteger{isInt: true, value: v} }

// NewList builds a nested-list node from children.
func NewList(children ...*NestedInteger) *NestedInteger {
	return &NestedInteger{isInt: false, list: children}
}

// IsInteger reports whether this node is a single integer.
func (ni *NestedInteger) IsInteger() bool { return ni.isInt }

// GetInteger returns the integer value (valid only if IsInteger()).
func (ni *NestedInteger) GetInteger() int { return ni.value }

// GetList returns the child list (valid only if !IsInteger()).
func (ni *NestedInteger) GetList() []*NestedInteger { return ni.list }

// ── Approach 1: Two-Pass DFS (find max depth, then weight bottom-up) ──────────
//
// twoPassDFS solves Nested List Weight Sum II by first finding the maximum depth
// of the structure, then summing each integer times (maxDepth - depth + 1).
//
// Intuition:
//
//	Unlike problem 339, weight grows toward the ROOT: leaves have weight 1,
//	the shallowest integers have the largest weight. If maxDepth is the deepest
//	nesting level, an integer at depth d contributes value*(maxDepth - d + 1).
//	So first learn maxDepth, then do a weighted DFS.
//
// Algorithm:
//  1. DFS to compute maxDepth over the whole structure.
//  2. DFS again; for an integer at depth d add value*(maxDepth - d + 1).
//  3. Return the accumulated sum.
//
// Time:  O(N) — N = total number of nodes (ints + lists); two linear passes.
// Space: O(D) — recursion stack, D = maximum nesting depth.
func twoPassDFS(nestedList []*NestedInteger) int {
	maxDepth := findMaxDepth(nestedList, 1) // depth of top-level items is 1
	return weightedSum(nestedList, 1, maxDepth)
}

// findMaxDepth returns the deepest nesting level reached from this list.
func findMaxDepth(list []*NestedInteger, depth int) int {
	best := depth // at least this deep just by being here
	for _, ni := range list {
		if !ni.IsInteger() {
			// Recurse one level deeper into the sublist.
			if d := findMaxDepth(ni.GetList(), depth+1); d > best {
				best = d
			}
		}
	}
	return best
}

// weightedSum adds value*(maxDepth-depth+1) for every integer beneath `list`.
func weightedSum(list []*NestedInteger, depth, maxDepth int) int {
	sum := 0
	for _, ni := range list {
		if ni.IsInteger() {
			// Weight is inverted: shallow items weigh more, leaves weigh 1.
			sum += ni.GetInteger() * (maxDepth - depth + 1)
		} else {
			sum += weightedSum(ni.GetList(), depth+1, maxDepth)
		}
	}
	return sum
}

// ── Approach 2: One-Pass BFS (accumulate unweighted level sums) (Optimal) ─────
//
// onePassBFS solves it in a single traversal using the identity: summing the
// running prefix of level sums equals summing each integer by its inverted
// weight — no need to know maxDepth in advance.
//
// Intuition:
//
//	Let levelSum be the running total of integer values seen down to the
//	current BFS level, and total the running answer. After processing each
//	level, do total += levelSum. An integer first entering at level ℓ then
//	gets added once for level ℓ, once for ℓ+1, ... down to the last level —
//	i.e. exactly (maxDepth - ℓ + 1) times, which is precisely its weight.
//	This computes the inverted-weight sum without a separate depth pass.
//
// Algorithm:
//  1. Start a queue with the top-level items; levelSum = 0, total = 0.
//  2. For each level: add every integer at this level into levelSum, enqueue
//     children of every list. Then total += levelSum.
//  3. When the queue empties, return total.
//
// Time:  O(N) — every node enqueued and dequeued once.
// Space: O(W) — W = maximum number of nodes on one level (queue width).
func onePassBFS(nestedList []*NestedInteger) int {
	queue := append([]*NestedInteger{}, nestedList...) // level 1 frontier
	levelSum, total := 0, 0
	for len(queue) > 0 {
		var next []*NestedInteger // nodes forming the next level
		for _, ni := range queue {
			if ni.IsInteger() {
				levelSum += ni.GetInteger() // carried into every deeper level
			} else {
				next = append(next, ni.GetList()...) // descend one level
			}
		}
		// Adding levelSum once per remaining level accumulates each integer's
		// inverted weight (maxDepth - itsDepth + 1) automatically.
		total += levelSum
		queue = next
	}
	return total
}

func main() {
	// Example 1: [[1,1],2,[1,1]] → 8. Four 1's at depth 2 (weight 1 each) and
	// one 2 at depth 1 (weight 2): 4*1 + 2*2 = 4 + 4 = 8.
	ex1 := []*NestedInteger{
		NewList(NewInt(1), NewInt(1)),
		NewInt(2),
		NewList(NewInt(1), NewInt(1)),
	}
	// Example 2: [1,[4,[6]]] → 17. 1 at depth 1 (weight 3), 4 at depth 2
	// (weight 2), 6 at depth 3 (weight 1): 1*3 + 4*2 + 6*1 = 3 + 8 + 6 = 17.
	ex2 := []*NestedInteger{
		NewInt(1),
		NewList(NewInt(4), NewList(NewInt(6))),
	}

	fmt.Println("=== Approach 1: Two-Pass DFS ===")
	fmt.Println(twoPassDFS(ex1)) // expected 8
	fmt.Println(twoPassDFS(ex2)) // expected 17

	fmt.Println("=== Approach 2: One-Pass BFS (Optimal) ===")
	fmt.Println(onePassBFS(ex1)) // expected 8
	fmt.Println(onePassBFS(ex2)) // expected 17
}
