# 0422 — Valid Word Square

> LeetCode #422 · Difficulty: Easy · 🔒 Premium
> **Categories:** Array, Matrix, String

---

## Problem Statement

Given an array of strings `words`, return `true` *if it forms a valid word square*.

A sequence of strings forms a valid **word square** if the `k`th row and column read the same string, where `0 <= k < max(numRows, numColumns)`.

**Example 1:**

```
Input: words = ["abcd","bnrt","crmy","dtye"]
Output: true
Explanation:
The 1st row and column both read "abcd".
The 2nd row and column both read "bnrt".
The 3rd row and column both read "crmy".
The 4th row and column both read "dtye".
Therefore, it is a valid word square.
```

**Example 2:**

```
Input: words = ["abcd","bnrt","crm","dt"]
Output: true
Explanation:
The 1st row and column both read "abcd".
The 2nd row and column both read "bnrt".
The 3rd row and column both read "crm".
The 4th row and column both read "dt".
Therefore, it is a valid word square.
```

**Example 3:**

```
Input: words = ["ball","area","read","lady"]
Output: false
Explanation:
The 3rd row reads "read" while the 3rd column reads "lead".
Therefore, it is NOT a valid word square.
```

**Constraints:**

- `1 <= words.length <= 500`
- `1 <= words[i].length <= 500`
- `words[i]` consists of only lowercase English letters.

---

## Company Frequency

| Company    | Frequency        | Last Reported |
|------------|------------------|---------------|
| Google     | ★★★☆☆ Medium     | 2023          |
| Amazon     | ★★☆☆☆ Low        | 2023          |
| Microsoft  | ★★☆☆☆ Low        | 2022          |
| Bloomberg  | ★★☆☆☆ Low        | 2022          |
| Facebook   | ★★☆☆☆ Low        | 2021          |

> ⚠️ Frequency data is crowd-sourced from LeetCode Discuss, Glassdoor, and
> community interview reports. Treat as a signal, not a guarantee.

---

## DSA Concepts Used

- **Matrix / grid transpose** — a valid word square is precisely a grid that is symmetric across its main diagonal (`grid[i][j] == grid[j][i]`); the whole task is a guarded transpose check → see [`/dsa/matrix_traversal.md`](/dsa/matrix_traversal.md)
- **String indexing with bounds care** — rows may be *ragged* (different lengths), so every character access must be guarded; the crux of the problem is handling missing cells correctly → see [`/dsa/string_algorithms.md`](/dsa/string_algorithms.md)

---

## Approaches Overview

| # | Approach | Time | Space | When to use |
|---|----------|------|-------|-------------|
| 1 | Build Columns and Compare | O(n·L) | O(L) | Most literal reading; easiest to explain the "row k = column k" rule |
| 2 | Symmetric Index Check (Optimal) | O(n·L) | O(1) | No string building; reframes the problem as diagonal symmetry |

`n` = number of rows, `L` = maximum word length.

---

## Approach 1 — Build Columns and Compare

### Intuition

The definition says row `k` must equal column `k`. So literally build column `k` — the `k`-th character of every row long enough to have one — and compare it to `words[k]`. The subtlety is ragged rows: a row only contributes to column `k` when it actually has a `k`-th character, and columns of a genuine square are contiguous from the top, so we stop at the first row that is too short.

### Algorithm

1. Let `n = len(words)`.
2. For each `k` in `0..n-1`:
   1. Build column `k`: walk rows top-to-bottom; while row `r` has a `k`-th character, append `words[r][k]`; stop at the first row that is too short.
   2. If the assembled column string `!= words[k]`, return `false`.
3. Return `true`.

### Complexity

- **Time:** O(n·L) — every character in the grid is read once while assembling columns.
- **Space:** O(L) — one column string at a time (its length is at most `n`, bounded by `L` for a square).

### Code

```go
func buildColumns(words []string) bool {
	n := len(words) // number of rows; a square has at most n columns
	for k := 0; k < n; k++ {
		var col strings.Builder // column k, assembled character by character
		for r := 0; r < n; r++ {
			// Row r contributes to column k only if it is long enough to have
			// a k-th character; a ragged (short) row simply stops early.
			if k < len(words[r]) {
				col.WriteByte(words[r][k]) // the k-th char of row r sits in column k, row r
			} else {
				break // rows are read top-to-bottom; once one is too short, a full
				// square would require the column to end here — but we compare the
				// assembled prefix directly to row k below, which catches mismatches.
			}
		}
		if col.String() != words[k] { // column k must read exactly like row k
			return false
		}
	}
	return true
}
```

### Dry Run

Input `words = ["abcd","bnrt","crmy","dtye"]`, `n = 4`.

| k | column built (row 0..3 char k) | words[k] | equal? |
|---|--------------------------------|----------|--------|
| 0 | a,b,c,d → `"abcd"`             | `"abcd"` | yes    |
| 1 | b,n,r,t → `"bnrt"`             | `"bnrt"` | yes    |
| 2 | c,r,m,y → `"crmy"`             | `"crmy"` | yes    |
| 3 | d,t,y,e → `"dtye"`             | `"dtye"` | yes    |

All columns match their rows → return `true`.

---

## Approach 2 — Symmetric Index Check (Optimal)

### Intuition

"Row `k` equals column `k` for all `k`" is exactly the statement that the grid is symmetric about its main diagonal: cell `(i, j)` must equal cell `(j, i)`. With ragged rows we must be strict about *existence*: if `(i, j)` is a real character then `(j, i)` must also be real and equal. In particular, if row `i` has a `j`-th character but there is no row `j` at all, the mirror cell is missing and the square is invalid. This checks each character once against its mirror — no strings built.

### Algorithm

1. For each row `i`, and each column `j` with `j < len(words[i])`:
   1. If `j >= n` (no `j`-th row exists) → return `false`.
   2. If `i >= len(words[j])` (row `j` has no `i`-th character, mirror missing) → return `false`.
   3. If `words[i][j] != words[j][i]` (asymmetric) → return `false`.
2. Return `true`.

### Complexity

- **Time:** O(n·L) — each existing character is compared once with its transpose.
- **Space:** O(1) — pure index arithmetic, no auxiliary structures.

### Code

```go
func symmetricCheck(words []string) bool {
	n := len(words) // number of rows == max possible number of columns
	for i := 0; i < n; i++ {
		for j := 0; j < len(words[i]); j++ {
			// (i, j) exists. Its mirror is (j, i). For validity that mirror
			// must also exist and match.

			// There must be a j-th row at all; otherwise column i extends past
			// the number of rows while row i still has a character — invalid.
			if j >= n {
				return false
			}
			// Row j must have an i-th character, else the mirror cell is absent.
			if i >= len(words[j]) {
				return false
			}
			// The defining symmetry: cell (i,j) equals cell (j,i).
			if words[i][j] != words[j][i] {
				return false
			}
		}
	}
	return true
}
```

### Dry Run

Input `words = ["ball","area","read","lady"]`, `n = 4` (Example 3, expected `false`). We scan cells until the first violation.

| (i, j) | words[i][j] | mirror (j, i) | words[j][i] | match? |
|--------|-------------|---------------|-------------|--------|
| (0,0)  | `b`         | (0,0)         | `b`         | yes    |
| (0,1)  | `a`         | (1,0) area[0] | `a`         | yes    |
| (0,2)  | `l`         | (2,0) read[0] | `r`         | **no** → return `false` |

The cell `(0,2)='l'` disagrees with its mirror `(2,0)='r'` (row 0 is `ball`, but column 0 reads `b,a,r,l`), so the square is invalid.

---

## Key Takeaways

- **Word square = diagonally symmetric grid.** The cleanest mental model is `grid[i][j] == grid[j][i]`; every valid-word-square problem reduces to a transpose comparison.
- **Ragged rows are the whole difficulty.** Guard every character access: a mirror cell can be missing either because the target row is too short *or* because that row does not exist at all. Both mean "not a square".
- Reframing "compare row to column" as "compare cell to its mirror" removes all string allocation and drops space to O(1).

---

## Related Problems

- LeetCode #425 — Word Squares (build all squares from a word list — the constructive counterpart)
- LeetCode #48 — Rotate Image (in-place transpose + reverse)
- LeetCode #766 — Toeplitz Matrix (another diagonal-based grid check)
- LeetCode #867 — Transpose Matrix (the operation this problem tests against)
