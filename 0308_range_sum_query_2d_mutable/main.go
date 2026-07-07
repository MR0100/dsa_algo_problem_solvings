package main

import "fmt"

// numMatrixADT is the shared contract each implementation satisfies so main()
// can drive the same official operation sequence through every approach.
type numMatrixADT interface {
	Update(row, col, val int)                 // set matrix[row][col] = val
	SumRegion(row1, col1, row2, col2 int) int // sum of the sub-rectangle
}

// ── Approach 1: Brute Force (mutable matrix, linear region sum) ───────────────
//
// BruteForce solves Range Sum Query 2D - Mutable by keeping the raw matrix and
// summing the requested rectangle on every query.
//
// Intuition:
//
//	The simplest structure: store the grid. Update is one assignment; SumRegion
//	loops over every cell in the rectangle. It is the correctness baseline and
//	shows the cost the smarter structures remove.
//
// Algorithm:
//
//	Update:    matrix[row][col] = val — O(1).
//	SumRegion: double loop over rows row1..row2, cols col1..col2 — O(mn).
//
// Time:  Update O(1), SumRegion O(m·n) (rectangle area).
// Space: O(m·n) — the stored matrix.
type BruteForce struct {
	mat [][]int // live copy of the grid
}

// NewBruteForce copies the input matrix.
func NewBruteForce(matrix [][]int) *BruteForce {
	m := len(matrix)
	cp := make([][]int, m)
	for i := range matrix {
		cp[i] = append([]int(nil), matrix[i]...) // deep-copy each row
	}
	return &BruteForce{mat: cp}
}

// Update assigns a new value in O(1).
func (b *BruteForce) Update(row, col, val int) {
	b.mat[row][col] = val
}

// SumRegion sums the inclusive rectangle in O(area).
func (b *BruteForce) SumRegion(row1, col1, row2, col2 int) int {
	sum := 0
	for r := row1; r <= row2; r++ {
		for c := col1; c <= col2; c++ {
			sum += b.mat[r][c] // accumulate each cell in the rectangle
		}
	}
	return sum
}

// ── Approach 2: Per-Row Prefix Sums ──────────────────────────────────────────
//
// RowPrefix solves Range Sum Query 2D - Mutable by keeping a prefix-sum array
// for each row, so a region sum scans only the rows (not every cell).
//
// Intuition:
//
//	Balance update against query. Each row keeps prefix sums so a single row's
//	segment [col1, col2] is answered in O(1). A region then loops the rows
//	row1..row2, adding each row's segment: O(rows). An update rewrites the
//	suffix of one row's prefixes: O(cols). Good when the grid is wide but short,
//	or updates and queries are balanced.
//
// Algorithm:
//
//	Build:     rowPre[r][c+1] = rowPre[r][c] + matrix[r][c].
//	Update:    set matrix[r][c]=val; recompute rowPre[r][c+1..n].
//	SumRegion: for r in row1..row2, add rowPre[r][col2+1] - rowPre[r][col1].
//
// Time:  Update O(n), SumRegion O(m) where m = rows in the region, n = cols.
// Space: O(m·n) — one prefix array per row plus the matrix.
type RowPrefix struct {
	mat    [][]int // current values, needed to rebuild a row after Update
	rowPre [][]int // rowPre[r] has length cols+1; rowPre[r][k] = sum of first k cells
	cols   int
}

// NewRowPrefix builds per-row prefix sums.
func NewRowPrefix(matrix [][]int) *RowPrefix {
	m := len(matrix)
	n := 0
	if m > 0 {
		n = len(matrix[0])
	}
	mat := make([][]int, m)
	rowPre := make([][]int, m)
	for r := 0; r < m; r++ {
		mat[r] = append([]int(nil), matrix[r]...)
		rowPre[r] = make([]int, n+1)
		for c := 0; c < n; c++ {
			rowPre[r][c+1] = rowPre[r][c] + matrix[r][c] // running row prefix
		}
	}
	return &RowPrefix{mat: mat, rowPre: rowPre, cols: n}
}

// Update rewrites the affected row's prefix suffix.
func (p *RowPrefix) Update(row, col, val int) {
	p.mat[row][col] = val // record the new value
	// Rebuild prefixes from the changed column to the end of the row.
	for c := col; c < p.cols; c++ {
		p.rowPre[row][c+1] = p.rowPre[row][c] + p.mat[row][c]
	}
}

// SumRegion adds each row's segment sum across the region's rows.
func (p *RowPrefix) SumRegion(row1, col1, row2, col2 int) int {
	sum := 0
	for r := row1; r <= row2; r++ {
		// O(1) segment sum for this row via its prefix array.
		sum += p.rowPre[r][col2+1] - p.rowPre[r][col1]
	}
	return sum
}

// ── Approach 3: 2D Fenwick Tree / Binary Indexed Tree (Optimal) ──────────────
//
// FenwickTree2D solves Range Sum Query 2D - Mutable with a 2D Binary Indexed
// Tree: point update and prefix-rectangle sum both in O(log m · log n).
//
// Intuition:
//
//	Extend the 1D Fenwick idea to two dimensions. tree[i][j] stores a partial
//	sum over a rectangle whose extents are the lowbits of i and j. A prefix sum
//	over the rectangle (0,0)..(r,c) is a nested lowbit descent; a point update
//	is a nested lowbit ascent. A general region uses inclusion–exclusion:
//	  sum(r1,c1,r2,c2) = P(r2,c2) - P(r1-1,c2) - P(r2,c1-1) + P(r1-1,c1-1).
//	Since updates set (not add) a value, we cache the grid and push the delta.
//
// Algorithm:
//
//	Build:     start all zeros; add each initial value via add().
//	add(r,c,d):for i=r..m by lowbit: for j=c..n by lowbit: tree[i][j]+=d.
//	query(r,c):for i=r..1 by lowbit: for j=c..1 by lowbit: s+=tree[i][j].
//	Update:    d=val-nums[r][c]; nums[r][c]=val; add(r+1,c+1,d) (1-indexed).
//	SumRegion: inclusion–exclusion of four prefix queries.
//
// Time:  Update O(log m · log n), SumRegion O(log m · log n), Build O(mn log m log n).
// Space: O(m·n) — the tree plus the value cache.
type FenwickTree2D struct {
	m, n int
	tree [][]int // 1-indexed, (m+1) x (n+1)
	nums [][]int // current values for delta computation on Update
}

// NewFenwickTree2D builds the 2D BIT from the initial matrix.
func NewFenwickTree2D(matrix [][]int) *FenwickTree2D {
	m := len(matrix)
	n := 0
	if m > 0 {
		n = len(matrix[0])
	}
	f := &FenwickTree2D{m: m, n: n}
	f.tree = make([][]int, m+1)
	for i := range f.tree {
		f.tree[i] = make([]int, n+1) // index 0 rows/cols unused
	}
	f.nums = make([][]int, m)
	for r := 0; r < m; r++ {
		f.nums[r] = make([]int, n)
		for c := 0; c < n; c++ {
			f.nums[r][c] = 0             // start empty so Update computes clean deltas
			f.Update(r, c, matrix[r][c]) // seed each value through the delta path
		}
	}
	return f
}

// add pushes delta at 1-indexed (r, c), climbing both dimensions by lowbit.
func (f *FenwickTree2D) add(r, c, delta int) {
	for i := r; i <= f.m; i += i & (-i) {
		for j := c; j <= f.n; j += j & (-j) {
			f.tree[i][j] += delta
		}
	}
}

// query returns the prefix-rectangle sum over 1-indexed (1,1)..(r,c).
func (f *FenwickTree2D) query(r, c int) int {
	s := 0
	for i := r; i > 0; i -= i & (-i) {
		for j := c; j > 0; j -= j & (-j) {
			s += f.tree[i][j]
		}
	}
	return s
}

// Update sets matrix[row][col] = val by pushing the difference into the tree.
func (f *FenwickTree2D) Update(row, col, val int) {
	delta := val - f.nums[row][col] // change from current value
	f.nums[row][col] = val          // record the new value
	f.add(row+1, col+1, delta)      // 1-indexed add
}

// SumRegion answers the rectangle via inclusion–exclusion of four prefixes.
func (f *FenwickTree2D) SumRegion(row1, col1, row2, col2 int) int {
	// All queries use 1-indexed inclusive corners.
	return f.query(row2+1, col2+1) -
		f.query(row1, col2+1) -
		f.query(row2+1, col1) +
		f.query(row1, col1)
}

// runExample drives the official operation sequence through one implementation
// and returns the outputs in LeetCode's format (null for void ops).
//
// Ops:  ["NumMatrix","sumRegion","update","sumRegion"]
// Args: [[[[3,0,1,4,2],[5,6,3,2,1],[1,2,0,1,5],[4,1,0,1,7],[1,0,3,0,5]]],
//
//	[2,1,4,3],[3,2,2],[2,1,4,3]]
func runExample(build func(matrix [][]int) numMatrixADT) []interface{} {
	matrix := [][]int{
		{3, 0, 1, 4, 2},
		{5, 6, 3, 2, 1},
		{1, 2, 0, 1, 5},
		{4, 1, 0, 1, 7},
		{1, 0, 3, 0, 5},
	}
	nm := build(matrix)                         // constructor → null
	out := []interface{}{nil}                   //
	out = append(out, nm.SumRegion(2, 1, 4, 3)) // → 8
	nm.Update(3, 2, 2)                          // → null
	out = append(out, nil)                      //
	out = append(out, nm.SumRegion(2, 1, 4, 3)) // → 10
	return out
}

func main() {
	fmt.Println("=== Approach 1: Brute Force ===")
	fmt.Println(runExample(func(m [][]int) numMatrixADT { return NewBruteForce(m) })) // [<nil> 8 <nil> 10]

	fmt.Println("=== Approach 2: Per-Row Prefix Sums ===")
	fmt.Println(runExample(func(m [][]int) numMatrixADT { return NewRowPrefix(m) })) // [<nil> 8 <nil> 10]

	fmt.Println("=== Approach 3: 2D Fenwick Tree (Optimal) ===")
	fmt.Println(runExample(func(m [][]int) numMatrixADT { return NewFenwickTree2D(m) })) // [<nil> 8 <nil> 10]
}
