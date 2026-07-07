package main

import "fmt"

// NestedInteger is the interface LeetCode provides. Each element is either a
// single integer OR a (possibly empty, possibly nested) list of NestedIntegers.
// Here we model it with a concrete struct so we can build the test inputs.
type NestedInteger struct {
	isInt bool             // true if this holds a single integer
	num   int              // valid only when isInt == true
	list  []*NestedInteger // valid only when isInt == false
}

// IsInteger reports whether this NestedInteger is a single integer (not a list).
func (n *NestedInteger) IsInteger() bool { return n.isInt }

// GetInteger returns the single integer this holds (undefined if it is a list).
func (n *NestedInteger) GetInteger() int { return n.num }

// GetList returns the nested list this holds (undefined if it is an integer).
func (n *NestedInteger) GetList() []*NestedInteger { return n.list }

// helpers to build test data
func leaf(v int) *NestedInteger                { return &NestedInteger{isInt: true, num: v} }
func nest(xs ...*NestedInteger) *NestedInteger { return &NestedInteger{isInt: false, list: xs} }

// ── Approach 1: Eager Flatten in Constructor (DFS Pre-flatten) ───────────────
//
// NestedIteratorEager solves Flatten Nested List Iterator by fully flattening
// the nested structure into a flat []int the moment the iterator is built, then
// walking that slice with a cursor.
//
// Intuition:
//
//	The traversal order is a simple depth-first, left-to-right walk of the
//	nested tree. If we do that whole walk once up front and record every
//	integer we see into a flat slice, Next/HasNext become trivial index
//	operations.
//
// Algorithm:
//  1. In the constructor, recursively DFS the nested list. For a leaf integer
//     append it; for a sub-list recurse into it.
//  2. HasNext: cursor < len(flat).
//  3. Next: return flat[cursor] and advance the cursor.
//
// Time:  O(N) total to build (N = count of integers) + O(1) per Next/HasNext.
// Space: O(N) — the flattened slice holds every integer.
type NestedIteratorEager struct {
	flat []int // every integer in DFS order
	pos  int   // index of the next integer to return
}

// ConstructorEager flattens the whole structure once, up front.
func ConstructorEager(nestedList []*NestedInteger) *NestedIteratorEager {
	it := &NestedIteratorEager{}
	var dfs func(list []*NestedInteger)
	dfs = func(list []*NestedInteger) {
		for _, ni := range list { // left-to-right
			if ni.IsInteger() {
				it.flat = append(it.flat, ni.GetInteger()) // record the leaf
			} else {
				dfs(ni.GetList()) // descend into the sub-list
			}
		}
	}
	dfs(nestedList)
	return it
}

// HasNext reports whether any integer remains.
func (it *NestedIteratorEager) HasNext() bool { return it.pos < len(it.flat) }

// Next returns the next integer and advances the cursor.
func (it *NestedIteratorEager) Next() int {
	v := it.flat[it.pos]
	it.pos++
	return v
}

// ── Approach 2: Lazy Stack (Optimal for memory) ──────────────────────────────
//
// NestedIteratorStack solves Flatten Nested List Iterator lazily: it never
// pre-flattens. It keeps an explicit stack of pending NestedIntegers and only
// unpacks lists on demand, so memory is proportional to nesting depth, not to
// the number of integers.
//
// Intuition:
//
//	Recursion uses a call stack; we can make that stack explicit. We push the
//	top-level list reversed so the first element is on top. Before each read
//	we "prime" the stack: while the top is a list, pop it and push its
//	children (reversed) so a leaf integer floats to the top.
//
// Algorithm:
//  1. Constructor: push all top-level items onto the stack in reverse order.
//  2. HasNext: while stack non-empty and top is a list, pop it and push its
//     children reversed. Return true iff the stack is non-empty afterwards
//     (top is now guaranteed to be an integer).
//  3. Next: call HasNext to prime, then pop and return the top integer.
//
// Time:  O(1) amortised per element (each NestedInteger is pushed/popped once).
// Space: O(D + top-level width) — proportional to nesting depth, not total N.
type NestedIteratorStack struct {
	stack []*NestedInteger // top of stack is the last element
}

// ConstructorStack seeds the stack with the top-level list, reversed so the
// first logical element sits on top of the stack.
func ConstructorStack(nestedList []*NestedInteger) *NestedIteratorStack {
	it := &NestedIteratorStack{}
	for i := len(nestedList) - 1; i >= 0; i-- { // reverse push
		it.stack = append(it.stack, nestedList[i])
	}
	return it
}

// HasNext primes the stack so its top is an integer, then reports emptiness.
func (it *NestedIteratorStack) HasNext() bool {
	for len(it.stack) > 0 {
		top := it.stack[len(it.stack)-1] // peek
		if top.IsInteger() {
			return true // a leaf is ready to be returned
		}
		it.stack = it.stack[:len(it.stack)-1] // pop the list
		list := top.GetList()
		for i := len(list) - 1; i >= 0; i-- { // push children reversed
			it.stack = append(it.stack, list[i])
		}
	}
	return false
}

// Next assumes HasNext primed the stack; pops and returns the top integer.
func (it *NestedIteratorStack) Next() int {
	it.HasNext()                          // ensure top is an integer
	top := it.stack[len(it.stack)-1]      // peek
	it.stack = it.stack[:len(it.stack)-1] // pop
	return top.GetInteger()
}

// drain runs the standard "while it.HasNext(): res = append(res, it.Next())"
// operation sequence LeetCode uses to validate the iterator.
func drainEager(nestedList []*NestedInteger) []int {
	it := ConstructorEager(nestedList)
	res := []int{}
	for it.HasNext() {
		res = append(res, it.Next())
	}
	return res
}

func drainStack(nestedList []*NestedInteger) []int {
	it := ConstructorStack(nestedList)
	res := []int{}
	for it.HasNext() {
		res = append(res, it.Next())
	}
	return res
}

func main() {
	// Example 1: [[1,1],2,[1,1]] -> [1,1,2,1,1]
	ex1 := []*NestedInteger{nest(leaf(1), leaf(1)), leaf(2), nest(leaf(1), leaf(1))}
	// Example 2: [1,[4,[6]]] -> [1,4,6]
	ex2 := []*NestedInteger{leaf(1), nest(leaf(4), nest(leaf(6)))}

	fmt.Println("=== Approach 1: Eager Flatten in Constructor ===")
	fmt.Println(drainEager(ex1)) // expected [1 1 2 1 1]
	fmt.Println(drainEager(ex2)) // expected [1 4 6]

	fmt.Println("=== Approach 2: Lazy Stack (Optimal) ===")
	fmt.Println(drainStack(ex1)) // expected [1 1 2 1 1]
	fmt.Println(drainStack(ex2)) // expected [1 4 6]
}
