package main

import "fmt"

// ── Approach 1: Flatten Eagerly in the Constructor ───────────────────────────
//
// EagerVector2D solves Flatten 2D Vector by copying every element into a single
// flat slice up front, then walking a cursor through it.
//
// Intuition:
//
//	The simplest way to iterate a 2D structure "as if" it were 1D is to actually
//	build the 1D version once. After that Next/HasNext are trivial index moves.
//
// Algorithm:
//  1. In the constructor, loop over each row and append every element to data.
//  2. Keep an index pos starting at 0.
//  3. HasNext = pos < len(data).
//  4. Next returns data[pos] and advances pos.
//
// Time:  Constructor O(N) where N is the total number of elements; Next/HasNext O(1).
// Space: O(N) — the flattened copy.
type EagerVector2D struct {
	data []int // every element copied out in row-major order
	pos  int   // index of the next element to return
}

// NewEagerVector2D flattens vec immediately into one slice.
func NewEagerVector2D(vec [][]int) *EagerVector2D {
	flat := []int{}           // will hold all elements in order
	for _, row := range vec { // walk each inner slice
		flat = append(flat, row...) // copy the whole row in
	}
	return &EagerVector2D{data: flat, pos: 0}
}

// Next returns the current element and advances the cursor.
func (v *EagerVector2D) Next() int {
	val := v.data[v.pos] // read element at the cursor
	v.pos++              // move cursor forward
	return val
}

// HasNext reports whether any element remains.
func (v *EagerVector2D) HasNext() bool {
	return v.pos < len(v.data) // still inside the flattened slice?
}

// ── Approach 2: Lazy Two-Pointer Iterator (Optimal) ──────────────────────────
//
// LazyVector2D solves Flatten 2D Vector without copying: it keeps two indices
// (outer row, inner column) and advances them on demand.
//
// Intuition:
//
//	Copying wastes O(N) memory. Instead track a "position" as a (row, col) pair.
//	The key helper advance() skips over any empty rows and rows we've exhausted,
//	so that whenever it finishes, (row, col) either points at a real element or
//	past the end. Calling advance() before every check keeps the state valid even
//	when the input contains empty inner lists like [] .
//
// Algorithm:
//  1. Store the original vec, row = 0, col = 0.
//  2. advance(): while row < len(vec) and col == len(vec[row]) (current row
//     used up or empty), move to the next row and reset col to 0.
//  3. HasNext: call advance(), then return row < len(vec).
//  4. Next: call HasNext() to normalise state, read vec[row][col], then col++.
//
// Time:  Next/HasNext amortised O(1) — each row/col index is advanced at most once total.
// Space: O(1) — only two integer cursors, no copy of the data.
type LazyVector2D struct {
	vec [][]int // reference to the original (not copied)
	row int     // current outer index
	col int     // current inner index within vec[row]
}

// NewLazyVector2D just stores the reference; no flattening work is done here.
func NewLazyVector2D(vec [][]int) *LazyVector2D {
	return &LazyVector2D{vec: vec, row: 0, col: 0}
}

// advance skips exhausted or empty rows so (row,col) points at a live element
// or one-past-the-end. Safe to call repeatedly.
func (v *LazyVector2D) advance() {
	// While we're on a real row but have consumed all of its columns
	// (col == len, which also catches empty rows where len == 0), step down.
	for v.row < len(v.vec) && v.col == len(v.vec[v.row]) {
		v.row++   // move to the next row
		v.col = 0 // restart column scan at its beginning
	}
}

// HasNext normalises the cursor then checks whether a row still remains.
func (v *LazyVector2D) HasNext() bool {
	v.advance()               // ensure cursor sits on a valid element (or the end)
	return v.row < len(v.vec) // any row left means an element is available
}

// Next returns the pointed-to element and steps the column cursor forward.
func (v *LazyVector2D) Next() int {
	v.HasNext()                // guarantee the cursor is on a valid element
	val := v.vec[v.row][v.col] // read it
	v.col++                    // advance within the row for next time
	return val
}

func main() {
	// Official example:
	// Vector2D(vec = [[1,2],[3],[4]])
	// next()    -> 1
	// next()    -> 2
	// next()    -> 3
	// hasNext() -> true
	// hasNext() -> true
	// next()    -> 4
	// hasNext() -> false
	fmt.Println("=== Approach 1: Eager Flatten ===")
	e := NewEagerVector2D([][]int{{1, 2}, {3}, {4}})
	fmt.Println(e.Next())    // expected 1
	fmt.Println(e.Next())    // expected 2
	fmt.Println(e.Next())    // expected 3
	fmt.Println(e.HasNext()) // expected true
	fmt.Println(e.HasNext()) // expected true
	fmt.Println(e.Next())    // expected 4
	fmt.Println(e.HasNext()) // expected false

	fmt.Println("=== Approach 2: Lazy Two-Pointer ===")
	l := NewLazyVector2D([][]int{{1, 2}, {3}, {4}})
	fmt.Println(l.Next())    // expected 1
	fmt.Println(l.Next())    // expected 2
	fmt.Println(l.Next())    // expected 3
	fmt.Println(l.HasNext()) // expected true
	fmt.Println(l.HasNext()) // expected true
	fmt.Println(l.Next())    // expected 4
	fmt.Println(l.HasNext()) // expected false

	// Edge case: empty inner lists must be skipped transparently.
	fmt.Println("=== Approach 2: Empty rows edge case ===")
	l2 := NewLazyVector2D([][]int{{}, {}, {1}, {}, {2, 3}, {}})
	out := []int{}
	for l2.HasNext() {
		out = append(out, l2.Next())
	}
	fmt.Println(out) // expected [1 2 3]
}
