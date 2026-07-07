# 0052 — N-Queens II

> LeetCode #52 · Difficulty: Hard
> **Categories:** Backtracking

---

## Problem Statement

Given an integer `n`, return the number of distinct solutions to the **n-queens puzzle**.

**Example 1**
```
Input:  n = 4
Output: 2
```

**Example 2**
```
Input:  n = 1
Output: 1
```

**Constraints**
- `1 <= n <= 9`

---

## Company Frequency

| Company   | Frequency       | Last Reported |
|-----------|-----------------|---------------|
| Amazon    | ★★★★☆ High      | 2024          |
| Google    | ★★★★☆ High      | 2024          |
| Microsoft | ★★★★☆ High      | 2024          |
| Meta      | ★★★☆☆ Medium    | 2023          |
| Bloomberg | ★★★☆☆ Medium    | 2023          |
| Apple     | ★★☆☆☆ Low       | 2022          |

> ⚠️ Frequency data is crowd-sourced from LeetCode Discuss, Glassdoor, and community interview reports. Treat as a signal, not a guarantee.

---

## DSA Concepts Used

- **Backtracking** — same as #51; skip board construction, just count.
- **Bitmask Backtracking** — encode attacked columns/diagonals as integer bitmasks for O(1) set operations and cache-friendly state passing.

---

## Approaches Overview

| # | Approach | Time | Space | When to use |
|---|----------|------|-------|-------------|
| 1 | Backtracking (Bool Arrays) | O(n!) | O(n) | Clear; same as #51 |
| 2 | Bitmask Backtracking ✅ | O(n!) | O(n) | Fastest in practice; O(1) bit ops, compact state |

---

## Approach 1 — Backtracking (Count Only)

### Intuition
Identical to #51 except we don't maintain or store the board — just increment `count` when `row == n`.

### Complexity
- **Time:** O(n!).
- **Space:** O(n) — `cols`/`diag`/`anti` arrays + recursion stack.

### Dry Run — `n = 4`
```
Total paths explored: starts at row 0, 4 column choices,
pruned by conflict checks → 2 complete placements counted.
```

---

## Approach 2 — Bitmask Backtracking (Recommended ✅)

### Intuition
Represent attacked columns/left-diagonals/right-diagonals as integer bitmasks:
- `cols` — bit `c` is set if column `c` is attacked.
- `leftDiag` — diagonals propagate **right** as we move to the next row, so shift right (`>>1`) each level.
- `rightDiag` — anti-diagonals propagate **left**, so shift left (`<<1`) each level.

`avail = allMask & ^(cols | leftDiag | rightDiag)` gives columns still free in this row.

Pick each candidate with `bit = avail & (-avail)` (lowest set bit), then `avail &= avail-1` to remove it.

### Why diagonals shift
When we descend from row `r` to row `r+1`:
- A left-to-right diagonal (`\`) column increases by 1 → shift left diag bits **right** (`>>1`) to realign.
- A right-to-left diagonal (`/`) column decreases by 1 → shift right diag bits **left** (`<<1`).

After `n` rows these bits have shifted out of the `n`-bit window, naturally removing stale attacks.

### Algorithm
```
bt(row, cols, leftDiag, rightDiag):
  if row == n: count++; return
  avail = allMask & NOT(cols | leftDiag | rightDiag)
  while avail != 0:
    bit = avail & (-avail)       // lowest set bit = one column
    avail &= avail - 1           // clear that bit
    bt(row+1, cols|bit, (leftDiag|bit)>>1, (rightDiag|bit)<<1)
```

### Complexity
- **Time:** O(n!) — same number of nodes as backtracking.
- **Space:** O(n) — recursion stack; no arrays needed.

### Code
```go
func bitmask(n int) int {
    count := 0; allMask := (1 << n) - 1
    var bt func(row, cols, leftDiag, rightDiag int)
    bt = func(row, cols, leftDiag, rightDiag int) {
        if row == n { count++; return }
        avail := allMask & ^(cols | leftDiag | rightDiag)
        for avail != 0 {
            bit := avail & (-avail)
            avail &= avail - 1
            bt(row+1, cols|bit, (leftDiag|bit)>>1, (rightDiag|bit)<<1)
        }
    }
    bt(0, 0, 0, 0); return count
}
```

### Dry Run — `n = 4`, `allMask = 0b1111 = 15`
```
bt(0, cols=0000, left=0000, right=0000):
  avail = 1111 & ~0000 = 1111  (all 4 cols free)
  bit=0001 (col 0): bt(1, 0001, 0000>>1=0000, 0001<<1=0010):
    avail = 1111 & ~(0001|0000|0010) = 1111 & ~0011 = 1100  (cols 2,3)
    bit=0100 (col 2): bt(2, 0101, 0010>>1=0001, 0110<<1=1100→too wide, masked to 1100):
      ...eventually leads to deadend (no valid column in row 3)
    bit=1000 (col 3): bt(2, ...): → leads to solution [col3,col1,col0,...] No
  bit=0010 (col 1): bt(1, 0010, 0001, 0100):
    avail = 1111 & ~(0010|0001|0100) = 1111 & ~0111 = 1000  (col 3 only)
    bit=1000 (col 3): bt(2, 1010, 0101, 0000):
      avail = 1111 & ~(1010|0101|0000) = 1111 & ~1111 = 0000  → dead end
  ...
  bit=0100 (col 2): bt(1, 0100, 0010, 1000):
    ...→ eventually finds 2 solutions
```
Result: count = 2 ✓

---

## Key Takeaways

- **`avail & (-avail)` isolates the lowest set bit** — a classic bitmask trick; `-x` in two's complement flips all bits above the lowest set bit.
- **`avail &= avail-1` clears the lowest set bit** — another standard bitmask pattern.
- **Diagonal shift encodes propagation for free** — no need for an offset-indexed diagonal array; the shifts naturally model diagonal movement.
- **No arrays, no hashing** — the entire state fits in three integers. This makes bitmask backtracking ~3–5× faster than array-based for n ≤ 16.
- **Known answer for n=8 is 92** — a good sanity check.

---

## Related Problems

- LeetCode #51 — N-Queens (same backtracking; also reconstruct boards)
- LeetCode #37 — Sudoku Solver (backtracking with constraint sets)
