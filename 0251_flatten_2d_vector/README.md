# 0251 — Flatten 2D Vector

> LeetCode #251 · Difficulty: Medium
> **Categories:** Design, Array, Two Pointers, Iterator

---

## Problem Statement

Design an iterator to flatten a 2D vector. It should support the `next` and `hasNext` operations.

Implement the `Vector2D` class:

- `Vector2D(int[][] vec)` initializes the object with the 2D vector `vec`.
- `next()` returns the next element from the 2D vector and moves the pointer one step forward. You may assume that all the calls to `next` are valid.
- `hasNext()` returns `true` if there are still some elements in the vector, and `false` otherwise.

**Example 1:**

```
Input
["Vector2D", "next", "next", "next", "hasNext", "hasNext", "next", "hasNext"]
[[[[1, 2], [3], [4]]], [], [], [], [], [], [], []]
Output
[null, 1, 2, 3, true, true, 4, false]

Explanation
Vector2D vector2D = new Vector2D([[1, 2], [3], [4]]);
vector2D.next();    // return 1
vector2D.next();    // return 2
vector2D.next();    // return 3
vector2D.hasNext(); // return True
vector2D.hasNext(); // return True
vector2D.next();    // return 4
vector2D.hasNext(); // return False
```

**Constraints:**

- `0 <= vec.length <= 200`
- `0 <= vec[i].length <= 500`
- `-500 <= vec[i][j] <= 500`
- At most `10^5` calls will be made to `next` and `hasNext`.

**Follow-up:** As an added challenge, try to code it using only iterators in C++ or iterators in Java.

---

## Company Frequency

| Company   | Frequency        | Last Reported |
|-----------|------------------|---------------|
| Google    | ★★★★☆ High       | 2023          |
| Amazon    | ★★★☆☆ Medium     | 2023          |
| Facebook  | ★★★☆☆ Medium     | 2022          |
| Microsoft | ★★☆☆☆ Low        | 2022          |

> ⚠️ Frequency data is crowd-sourced from LeetCode Discuss, Glassdoor, and community interview reports. Treat as a signal, not a guarantee.

---

## DSA Concepts Used

- **Iterator / stateful design** — the class exposes `next`/`hasNext` and keeps traversal state between calls → see [`/dsa/design_data_structures.md`](/dsa/design_data_structures.md)
- **Two pointers** — the optimal iterator advances an (row, col) pair lazily instead of copying → see [`/dsa/two_pointers.md`](/dsa/two_pointers.md)

---

## Approaches Overview

| # | Approach | Time | Space | When to use |
|---|----------|------|-------|-------------|
| 1 | Eager Flatten | Ctor O(N), ops O(1) | O(N) | Simple, when memory is cheap and the data is static |
| 2 | Lazy Two-Pointer (Optimal) | ops amortised O(1) | O(1) | When you must not copy (huge/streaming input), handles empty rows |

---

## Approach 1 — Eager Flatten

### Intuition
The easiest way to iterate a 2D structure as if it were 1D is to actually build the 1D version once in the constructor, then walk a cursor.

### Algorithm
1. In the constructor, loop over each row and append every element into a flat slice `data`.
2. Keep an index `pos` starting at 0.
3. `HasNext` = `pos < len(data)`.
4. `Next` returns `data[pos]` and increments `pos`.

### Complexity
- **Time:** Constructor O(N) where N is total elements; `Next`/`HasNext` O(1) — pure index math.
- **Space:** O(N) — the full flattened copy.

### Code
```go
type EagerVector2D struct {
	data []int
	pos  int
}

func NewEagerVector2D(vec [][]int) *EagerVector2D {
	flat := []int{}
	for _, row := range vec {
		flat = append(flat, row...)
	}
	return &EagerVector2D{data: flat, pos: 0}
}

func (v *EagerVector2D) Next() int {
	val := v.data[v.pos]
	v.pos++
	return val
}

func (v *EagerVector2D) HasNext() bool {
	return v.pos < len(v.data)
}
```

### Dry Run
Input `vec = [[1,2],[3],[4]]`. Constructor flattens to `data = [1,2,3,4]`, `pos = 0`.

| Call      | pos before | Action                | Returns |
|-----------|-----------|-----------------------|---------|
| Next()    | 0         | read data[0]=1, pos→1 | 1       |
| Next()    | 1         | read data[1]=2, pos→2 | 2       |
| Next()    | 2         | read data[2]=3, pos→3 | 3       |
| HasNext() | 3         | 3 < 4                 | true    |
| HasNext() | 3         | 3 < 4                 | true    |
| Next()    | 3         | read data[3]=4, pos→4 | 4       |
| HasNext() | 4         | 4 < 4 is false        | false   |

---

## Approach 2 — Lazy Two-Pointer (Optimal)

### Intuition
Copying wastes O(N) memory. Instead track a position as a `(row, col)` pair over the original vector. A helper `advance()` skips empty or exhausted rows so the cursor always points at a real element or past the end — this transparently handles empty inner lists like `[]`.

### Algorithm
1. Store `vec`, `row = 0`, `col = 0`.
2. `advance()`: while `row < len(vec)` and `col == len(vec[row])` (current row used up or empty), do `row++`, `col = 0`.
3. `HasNext`: call `advance()`, then return `row < len(vec)`.
4. `Next`: call `HasNext()` to normalise state, read `vec[row][col]`, then `col++`.

### Complexity
- **Time:** `Next`/`HasNext` amortised O(1) — each row and column index is advanced at most once across the whole traversal.
- **Space:** O(1) — only two integer cursors; the input is never copied.

### Code
```go
type LazyVector2D struct {
	vec [][]int
	row int
	col int
}

func NewLazyVector2D(vec [][]int) *LazyVector2D {
	return &LazyVector2D{vec: vec, row: 0, col: 0}
}

func (v *LazyVector2D) advance() {
	for v.row < len(v.vec) && v.col == len(v.vec[v.row]) {
		v.row++
		v.col = 0
	}
}

func (v *LazyVector2D) HasNext() bool {
	v.advance()
	return v.row < len(v.vec)
}

func (v *LazyVector2D) Next() int {
	v.HasNext()
	val := v.vec[v.row][v.col]
	v.col++
	return val
}
```

### Dry Run
Input `vec = [[1,2],[3],[4]]`, start `row=0, col=0`.

| Call      | advance() result | read           | after   | Returns |
|-----------|------------------|----------------|---------|---------|
| Next()    | row=0,col=0 (0<2)| vec[0][0]=1    | col→1   | 1       |
| Next()    | row=0,col=1 (1<2)| vec[0][1]=2    | col→2   | 2       |
| Next()    | col==len → row=1,col=0 | vec[1][0]=3 | col→1 | 3       |
| HasNext() | col==len → row=2,col=0; 2<3 | —    | —       | true    |
| HasNext() | already valid; 2<3 | —            | —       | true    |
| Next()    | row=2,col=0 (0<1)| vec[2][0]=4    | col→1   | 4       |
| HasNext() | col==len → row=3; 3<3 false | —    | —       | false   |

---

## Key Takeaways
- An iterator over nested structure is easiest to reason about as a `(row, col)` cursor plus a normalising `advance()` that hops over empty/exhausted rows.
- Calling `HasNext()` from inside `Next()` centralises the "make the cursor valid" logic in one place — no duplicated skipping code.
- Eager flatten trades memory for simplicity; lazy iteration is the right answer when data is large or streamed.

---

## Related Problems
- LeetCode #341 — Flatten Nested List Iterator (nested structure, iterator design)
- LeetCode #284 — Peeking Iterator (wrapping an iterator with lookahead)
- LeetCode #173 — Binary Search Tree Iterator (lazy in-order traversal state)
- LeetCode #281 — Zigzag Iterator (multi-list interleaving cursor)
