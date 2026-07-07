# 0006 — Zigzag Conversion

> LeetCode #6 · Difficulty: Medium
> **Categories:** String, Simulation, Math

---

## Problem Statement

The string `"PAYPALISHIRING"` is written in a zigzag pattern on a given number of rows like this:

```
P   A   H   N
A P L S I I G
Y   I   R
```

And then read line by line: `"PAHNAPLSIIGYIR"`.

Write the code that will take a string and make this conversion given a number of rows.

**Example 1**
```
Input:  s = "PAYPALISHIRING", numRows = 3
Output: "PAHNAPLSIIGYIR"
```

**Example 2**
```
Input:  s = "PAYPALISHIRING", numRows = 4
Output: "PINALSIGYAHRPI"
Explanation:
P     I    N
A   L S  I G
Y A   H R
P     I
```

**Example 3**
```
Input:  s = "A", numRows = 1
Output: "A"
```

**Constraints**
- `1 <= s.length <= 1000`
- `s` consists of English letters (lower-case and upper-case), `','` and `'.'`.
- `1 <= numRows <= 1000`

---

## Company Frequency

| Company   | Frequency       | Last Reported |
|-----------|-----------------|---------------|
| Amazon    | ★★★★☆ High      | 2024          |
| Microsoft | ★★★★☆ High      | 2024          |
| Google    | ★★★☆☆ Medium    | 2023          |
| Bloomberg | ★★★☆☆ Medium    | 2023          |
| Apple     | ★★☆☆☆ Low       | 2022          |
| Adobe     | ★★☆☆☆ Low       | 2022          |

> ⚠️ Frequency data is crowd-sourced from LeetCode Discuss, Glassdoor, and community interview reports. Treat as a signal, not a guarantee.

---

## DSA Concepts Used

- **String / Simulation** — Approach 1 simulates the physical writing of characters into rows by tracking a current row index and a bouncing direction.
- **Math / Cycle Pattern** — Approach 2 observes that the zigzag repeats with period `cycleLen = 2*(numRows-1)` and derives the output positions analytically.

---

## Approaches Overview

| # | Approach | Time | Space | When to use |
|---|----------|------|-------|-------------|
| 1 | Simulate with Row Buffers | O(n) | O(n) | Easy to code; great for interviews |
| 2 | Math — Direct Formula ✅ | O(n) | O(n) | Faster constant; no extra data structures |

Both approaches are O(n) time and O(n) space (the output string is always O(n)). The math formula has a smaller constant and is the "slick" interview answer once you see the cycle pattern.

---

## Approach 1 — Simulate with Row Buffers

### Intuition
Physically simulate writing the characters: maintain a "current row" integer and a "going down" boolean. Each character is appended to the buffer for the current row. When `curRow` reaches 0 or `numRows-1`, flip direction.

### Algorithm
1. Create `numRows` string builders.
2. For each character `ch` in `s`:
   - Append to `rows[curRow]`.
   - If `curRow == 0` or `curRow == numRows-1`: flip `goingDown`.
   - Advance `curRow` by +1 or -1.
3. Concatenate all builders in order.

### Complexity
- **Time:** O(n) — one pass over `s` + O(n) concatenation.
- **Space:** O(n) — row buffers total n characters.

### Code
```go
func simulate(s string, numRows int) string {
    if numRows == 1 || numRows >= len(s) { return s }
    rows := make([]strings.Builder, numRows)
    curRow, goingDown := 0, false
    for _, ch := range s {
        rows[curRow].WriteRune(ch)
        if curRow == 0 || curRow == numRows-1 { goingDown = !goingDown }
        if goingDown { curRow++ } else { curRow-- }
    }
    var result strings.Builder
    for _, row := range rows { result.WriteString(row.String()) }
    return result.String()
}
```

### Dry Run — `s = "PAYPALISHIRING"`, `numRows = 3`
```
Row 0: P  A  H  N           (indices 0,4,8,12)
Row 1: A  P  L  S  I  I  G  (indices 1,3,5,7,9,11,13)
Row 2: Y  I  R              (indices 2,6,10)

Direction trace:
  P(0↓) A(1↓) Y(2↑) P(1↑) A(0↓) L(1↓) I(2↑) S(1↑) H(0↓) I(1↓) R(2↑) I(1↑) N(0↓) G(1↓)

Concat rows: "PAHN" + "APLSIIG" + "YIR" = "PAHNAPLSIIGYIR" ✓
```

---

## Approach 2 — Math Formula (Recommended ✅)

### Intuition
The zigzag repeats every `cycleLen = 2*(numRows-1)` characters. For row `r` and cycle starting at index `j` (j = 0, cycleLen, 2·cycleLen, ...):
- **Down-stroke character:** `s[j + r]`
- **Up-stroke character:** `s[j + cycleLen - r]` (only for interior rows)

By iterating rows then cycles, we read characters in the correct output order.

### Algorithm
1. `cycleLen = 2 * (numRows - 1)`.
2. For each `row` from 0 to numRows-1:
   - For each cycle start `j = 0, cycleLen, 2*cycleLen, ...`:
     - Append `s[j+row]` (down-stroke).
     - If interior row and `j+cycleLen-row < n`: append `s[j+cycleLen-row]` (up-stroke).

### Complexity
- **Time:** O(n) — every character output exactly once.
- **Space:** O(n) — result string.

### Code
```go
func mathFormula(s string, numRows int) string {
    if numRows == 1 || numRows >= len(s) { return s }
    n, cycleLen := len(s), 2*(numRows-1)
    var result strings.Builder
    for row := 0; row < numRows; row++ {
        for j := 0; j+row < n; j += cycleLen {
            result.WriteByte(s[j+row])
            upIdx := j + cycleLen - row
            if row != 0 && row != numRows-1 && upIdx < n {
                result.WriteByte(s[upIdx])
            }
        }
    }
    return result.String()
}
```

### Dry Run — `s = "PAYPALISHIRING"`, `numRows = 4`, `cycleLen = 6`
```
Indices: 0  1  2  3  4  5  6  7  8  9  10 11 12 13
Chars:   P  A  Y  P  A  L  I  S  H  I  R  I  N  G

Row 0: j=0 → s[0]='P'; j=6 → s[6]='I'; j=12 → s[12]='N'         → "PIN"
Row 1: j=0 → down s[1]='A', up s[5]='L'
       j=6 → down s[7]='S', up s[11]='I'                           → "ALSI"
Row 2: j=0 → down s[2]='Y', up s[4]='A'
       j=6 → down s[8]='H', up s[10]='R'                           → "YAHR"
Row 3: j=0 → s[3]='P'; j=6 → s[9]='I'                             → "PI"

Concat: "PIN"+"ALSI"+"YAHR"+"PI" = "PINALSIGYAHRPI" ✓
```

---

## Key Takeaways

- **CycleLen = 2*(numRows-1)** — this is the fundamental period of the zigzag pattern. Memorise this.
- **Top and bottom rows have one character per cycle; interior rows have two** — the up-stroke exists only for `0 < row < numRows-1`.
- **Edge cases:** `numRows == 1` (no zigzag) and `numRows >= len(s)` (each character on its own row, output = input). Handle both with an early return.
- **Both approaches are O(n)** — prefer the simulation in interviews for clarity; mention the math formula as an optimisation for the constant factor.

---

## Implementation (Go)

See [main.go](main.go).

### Verification
```
s="PAYPALISHIRING"  numRows=3  expect="PAHNAPLSIIGYIR"
  Approach 1: Simulate (row buffers) O(n) T | O(n) S      → "PAHNAPLSIIGYIR"
  Approach 2: Math formula         ✅ O(n) T | O(n) S      → "PAHNAPLSIIGYIR"

s="PAYPALISHIRING"  numRows=4  expect="PINALSIGYAHRPI"
  Approach 1: Simulate (row buffers) O(n) T | O(n) S      → "PINALSIGYAHRPI"
  Approach 2: Math formula         ✅ O(n) T | O(n) S      → "PINALSIGYAHRPI"
```

---

## Related Problems

- LeetCode #5 — Longest Palindromic Substring (also involves mapping string positions)
- LeetCode #944 — Delete Columns to Make Sorted (column reading of a 2-D string grid)
