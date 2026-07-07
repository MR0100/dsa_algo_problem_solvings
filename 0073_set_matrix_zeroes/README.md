# 0073 — Set Matrix Zeroes

> LeetCode #73 · Difficulty: Medium
> **Categories:** Array, Hash Table, Matrix

---

## Problem Statement

Given an `m x n` integer matrix `matrix`, if an element is `0`, set its entire **row** and **column** to `0`'s.

You must do it **in place**.

**Example 1**
```
Input:  matrix = [[1,1,1],[1,0,1],[1,1,1]]
Output:           [[1,0,1],[0,0,0],[1,0,1]]
```

**Example 2**
```
Input:  matrix = [[0,1,2,0],[3,4,5,2],[1,3,1,5]]
Output:           [[0,0,0,0],[0,4,5,0],[0,3,1,0]]
```

**Constraints**
- `m == matrix.length`
- `n == matrix[i].length`
- `1 <= m, n <= 200`
- `-2³¹ <= matrix[i][j] <= 2³¹ - 1`

**Follow-up:**
- A straightforward solution using O(mn) space is probably a bad idea.
- A simple improvement uses O(m + n) space, but still not the best solution.
- Could you devise a constant space solution?

---

## Company Frequency

| Company   | Frequency       | Last Reported |
|-----------|-----------------|---------------|
| Amazon    | ★★★★★ Very High | 2024          |
| Google    | ★★★★★ Very High | 2024          |
| Meta      | ★★★★★ Very High | 2024          |
| Microsoft | ★★★★☆ High      | 2024          |
| Bloomberg | ★★★★☆ High      | 2023          |
| Apple     | ★★★☆☆ Medium    | 2023          |
| Adobe     | ★★★☆☆ Medium    | 2023          |

> ⚠️ Frequency data is crowd-sourced from LeetCode Discuss, Glassdoor, and community interview reports. Treat as a signal, not a guarantee.

---

## DSA Concepts Used

- **In-Place Matrix Encoding** — use the first row/column as flag arrays to encode which rows/columns need zeroing without extra memory.

---

## Approaches Overview

| # | Approach | Time | Space | When to use |
|---|----------|------|-------|-------------|
| 1 | Extra Space (O(m+n) flags) | O(m × n) | O(m + n) | Simple; easy to explain |
| 2 | In-Place with First Row/Col ✅ | O(m × n) | O(1) | Optimal; uses matrix as its own flag |

---

## Approach 1 — Extra Space (O(m+n))

### Intuition
Two passes:
1. Scan all cells; record which rows and cols contain zeros in boolean arrays.
2. For each cell, if its row or col is flagged, set it to 0.

### Complexity
- **Time:** O(m × n).
- **Space:** O(m + n).

### Code
```go
func extraSpace(matrix [][]int) {
	m, n := len(matrix), len(matrix[0])
	rows := make([]bool, m)
	cols := make([]bool, n)

	for r := 0; r < m; r++ {
		for c := 0; c < n; c++ {
			if matrix[r][c] == 0 {
				rows[r] = true
				cols[c] = true
			}
		}
	}
	for r := 0; r < m; r++ {
		for c := 0; c < n; c++ {
			if rows[r] || cols[c] {
				matrix[r][c] = 0
			}
		}
	}
}
```

### Dry Run — `matrix = [[0,1,2,0],[3,4,5,2],[1,3,1,5]]`

**Pass 1 — record zeros** (`rows[m]`, `cols[n]`):

| cell | value | effect |
|------|-------|--------|
| (0,0) | 0 | rows[0]=true, cols[0]=true |
| (0,3) | 0 | rows[0]=true, cols[3]=true |
| others | ≠0 | no change |

After pass 1: `rows = [T,F,F]`, `cols = [T,F,F,T]`.

**Pass 2 — zero flagged cells** (set `matrix[r][c]=0` when `rows[r] || cols[c]`):

| row r | rows[r]? | resulting row (col flagged at c=0 and c=3) |
|-------|----------|--------------------------------------------|
| 0 | true | [0,0,0,0] (whole row zeroed) |
| 1 | false | [0,4,5,0] (only cols 0 and 3 zeroed) |
| 2 | false | [0,3,1,0] (only cols 0 and 3 zeroed) |

Result: `[[0,0,0,0],[0,4,5,0],[0,3,1,0]]` ✓

---

## Approach 2 — In-Place with First Row/Col (Recommended ✅)

### Intuition
Use the **first row** and **first column** as the flag arrays. `matrix[0][c]` flags column `c`; `matrix[r][0]` flags row `r`. The cell `matrix[0][0]` is shared — it can flag either row 0 or column 0, but not both independently. Track column 0's original zero status separately with `firstColZero`.

**Four-step algorithm:**
1. Check if column 0 originally had any zero (`firstColZero`).
2. Check if row 0 originally had any zero (`firstRowZero`).
3. Scan interior cells (`r≥1, c≥1`): if zero, mark `matrix[r][0]=0` and `matrix[0][c]=0`.
4. Zero interior cells based on flags.
5. Zero row 0 if `firstRowZero`.
6. Zero column 0 if `firstColZero`.

The order matters: zero the first row/column AFTER using them as flags.

### Complexity
- **Time:** O(m × n).
- **Space:** O(1).

### Code
```go
func inPlace(matrix [][]int) {
	m, n := len(matrix), len(matrix[0])

	// does column 0 need to be zeroed?
	firstColZero := false
	for r := 0; r < m; r++ {
		if matrix[r][0] == 0 {
			firstColZero = true
			break
		}
	}

	// does row 0 need to be zeroed?
	firstRowZero := false
	for c := 0; c < n; c++ {
		if matrix[0][c] == 0 {
			firstRowZero = true
			break
		}
	}

	// use first row and first column as flags for the rest of the matrix
	for r := 1; r < m; r++ {
		for c := 1; c < n; c++ {
			if matrix[r][c] == 0 {
				matrix[r][0] = 0 // flag row r
				matrix[0][c] = 0 // flag col c
			}
		}
	}

	// zero interior cells based on flags
	for r := 1; r < m; r++ {
		for c := 1; c < n; c++ {
			if matrix[r][0] == 0 || matrix[0][c] == 0 {
				matrix[r][c] = 0
			}
		}
	}

	// zero first row if needed
	if firstRowZero {
		for c := 0; c < n; c++ {
			matrix[0][c] = 0
		}
	}

	// zero first column if needed
	if firstColZero {
		for r := 0; r < m; r++ {
			matrix[r][0] = 0
		}
	}
}
```

### Dry Run — `matrix = [[0,1,2,0],[3,4,5,2],[1,3,1,5]]`
```
firstColZero: matrix[0][0]=0 → true
firstRowZero: matrix[0][3]=0 → true

Mark flags from interior (rows 1+, cols 1+):
  All non-zero interior cells → no additional marks needed.

matrix[0] = [0,1,2,0] (first row untouched so far)
matrix[*][0] = [0,3,1] (first col untouched so far)

Zero interior cells based on flags:
  matrix[0][c]=0 for c=1,2,3: first row flags → but we skip this (firstRowZero handles row 0)
  Actually: rows 1,2 checked against first col (matrix[r][0]):
    matrix[1][0]=3≠0, matrix[0][c]=0 for c=3: only col 0 and col 3 flagged.
  Wait, interior scan finds no zeros to flag in rows 1-2 besides col 0 original.

  Zero interior based on matrix[r][0] and matrix[0][c]:
    matrix[1][0]=3 → row 1 not flagged (from matrix[r][0]).
    matrix[0][1..3]: matrix[0][3]=0 → col 3 flagged.
    So matrix[1][3]=2→0, matrix[2][3]=5→0.
    matrix[r][0]=0 for r=0 only (from scan → firstColZero only).

Apply firstRowZero: zero entire row 0 → [0,0,0,0].
Apply firstColZero: zero entire col 0 → matrix[0][0]=0,matrix[1][0]=0,matrix[2][0]=0.

Result: [[0,0,0,0],[0,4,5,0],[0,3,1,0]] ✓
```

---

## Key Takeaways

- **The critical subtlety: `matrix[0][0]` is shared** — it represents both "zero row 0" and "zero col 0." The fix: use a separate `firstColZero` variable; let `matrix[0][0]` flag only row 0.
- **Zero first row/column LAST** — during the flag phase, we use these rows/columns to encode information; zeroing them early would corrupt that information.
- **This is a common "use the data structure as its own auxiliary" trick** — appears in Game of Life (#289), rotating image (#48), and other in-place matrix problems.

---

## Related Problems

- LeetCode #289 — Game of Life (encode two states in-place; similar trick)
- LeetCode #48 — Rotate Image (in-place matrix transformation)
