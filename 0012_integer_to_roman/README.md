# 0012 — Integer to Roman

> LeetCode #12 · Difficulty: Medium
> **Categories:** Hash Table, Math, String

---

## Problem Statement

Seven different symbols represent Roman numerals with the following values:

| Symbol | Value |
|--------|-------|
| I | 1 |
| V | 5 |
| X | 10 |
| L | 50 |
| C | 100 |
| D | 500 |
| M | 1000 |

Roman numerals are formed by appending the conversions of decimal place values from highest to lowest. Converting a decimal place value into a Roman numeral has the following rules:

- If the value does not start with 4 or 9, select the symbol of the maximum value that can be subtracted from the input, append that symbol to the result, subtract its value, and convert the remainder to a Roman numeral.
- If the value starts with 4 or 9 use the **subtractive form** representing one symbol subtracted from the following symbol: IV (4), IX (9), XL (40), XC (90), CD (400), CM (900).
- Only powers of 10 (I, X, C, M) can be appended multiple times; V, L, D can appear only once.

Given an integer, convert it to a Roman numeral.

**Example 1**
```
Input:  num = 3749
Output: "MMMDCCXLIX"
Explanation: 3000 = MMM, 700 = DCC, 40 = XL, 9 = IX
```

**Example 2**
```
Input:  num = 58
Output: "LVIII"
Explanation: 50 = L, 8 = VIII
```

**Example 3**
```
Input:  num = 1994
Output: "MCMXCIV"
Explanation: 1000 = M, 900 = CM, 90 = XC, 4 = IV
```

**Constraints**
- `1 <= num <= 3999`

---

## Company Frequency

| Company   | Frequency       | Last Reported |
|-----------|-----------------|---------------|
| Amazon    | ★★★★☆ High      | 2024          |
| Google    | ★★★★☆ High      | 2024          |
| Microsoft | ★★★☆☆ Medium    | 2024          |
| Bloomberg | ★★★☆☆ Medium    | 2023          |
| Apple     | ★★★☆☆ Medium    | 2023          |
| Meta      | ★★☆☆☆ Low       | 2022          |
| Adobe     | ★★☆☆☆ Low       | 2022          |

> ⚠️ Frequency data is crowd-sourced from LeetCode Discuss, Glassdoor, and community interview reports. Treat as a signal, not a guarantee.

---

## DSA Concepts Used

- **Greedy** — Approach 1 always picks the largest symbol that fits, which is provably optimal for Roman numerals.
- **Math / Digit Decomposition** — Approach 2 separates the number into decimal digit positions and maps each digit independently using pre-built lookup arrays.

---

## Approaches Overview

| # | Approach | Time | Space | When to use |
|---|----------|------|-------|-------------|
| 1 | Greedy — Value-Symbol Table | O(1) | O(1) | Simple to derive; great for interviews |
| 2 | Digit-by-Digit Table Lookup ✅ | O(1) | O(1) | Fastest; most elegant; zero arithmetic at runtime |

Both are O(1) because `num ≤ 3999` bounds the input. The output length is bounded by `~15` characters.

---

## Approach 1 — Greedy with Value-Symbol Table

### Intuition
Build a table of all 13 Roman values (7 additive + 6 subtractive) in descending order. Repeatedly find the largest value that fits into `num`, append its symbol, and subtract it. The greedy choice is always optimal because Roman numeral construction is defined this way.

The 6 subtractive cases to include in the table:
```
CM=900, CD=400, XC=90, XL=40, IX=9, IV=4
```

### Algorithm
1. `vals = [1000, 900, 500, 400, 100, 90, 50, 40, 10, 9, 5, 4, 1]`
2. For each `(val, sym)` pair: while `num >= val`, append `sym`, subtract `val`.

### Complexity
- **Time:** O(1) — fixed-size table; bounded iterations.
- **Space:** O(1).

### Code
```go
func greedyTable(num int) string {
    vals := []int{1000, 900, 500, 400, 100, 90, 50, 40, 10, 9, 5, 4, 1}
    syms := []string{"M","CM","D","CD","C","XC","L","XL","X","IX","V","IV","I"}
    result := ""
    for i, v := range vals {
        for num >= v { result += syms[i]; num -= v }
    }
    return result
}
```

### Dry Run — `num = 1994`
```
1994 >= 1000 → result="M",   num=994
 994 >= 900  → result="MCM", num=94
  94 >= 90   → result="MCMXC", num=4
   4 >= 4    → result="MCMXCIV", num=0
Return "MCMXCIV" ✓
```

---

## Approach 2 — Digit-by-Digit Table Lookup (Recommended ✅)

### Intuition
For any number 1–3999, each decimal digit position (thousands, hundreds, tens, ones) contributes an independent Roman numeral fragment. Pre-build four lookup arrays, one per position, indexed by the digit value 0–9. Extract each digit and concatenate the four fragments.

This avoids any loops at runtime — it's four array accesses and a string concatenation.

### Algorithm
1. Pre-build `thousands[0..3]`, `hundreds[0..9]`, `tens[0..9]`, `ones[0..9]`.
2. Return `thousands[num/1000] + hundreds[(num%1000)/100] + tens[(num%100)/10] + ones[num%10]`.

### Complexity
- **Time:** O(1) — four lookups.
- **Space:** O(1) — the four fixed arrays.

### Code
```go
func digitByDigit(num int) string {
    thousands := []string{"","M","MM","MMM"}
    hundreds  := []string{"","C","CC","CCC","CD","D","DC","DCC","DCCC","CM"}
    tens      := []string{"","X","XX","XXX","XL","L","LX","LXX","LXXX","XC"}
    ones      := []string{"","I","II","III","IV","V","VI","VII","VIII","IX"}
    return thousands[num/1000] + hundreds[(num%1000)/100] + tens[(num%100)/10] + ones[num%10]
}
```

### Dry Run — `num = 3749`
```
num/1000      = 3 → thousands[3] = "MMM"
(num%1000)/100 = 7 → hundreds[7]  = "DCC"
(num%100)/10   = 4 → tens[4]      = "XL"
num%10         = 9 → ones[9]      = "IX"
Concat: "MMM" + "DCC" + "XL" + "IX" = "MMMDCCXLIX" ✓
```

---

## Key Takeaways

- **Include the 6 subtractive cases in the greedy table** — 4, 9, 40, 90, 400, 900. Forgetting any one of them is the most common interview mistake.
- **Digit-by-digit is the cleanest** — once you have the four lookup tables, the entire function is one `return` statement. Interviewers love the elegance.
- **`num ≤ 3999` simplifies things** — there is no M equivalent at 4000+ (MMMM would not be valid), so thousands only go up to 3.

---

## Implementation (Go)

See [main.go](main.go).

### Verification
```
3    → "III"       ✓
58   → "LVIII"     ✓
1994 → "MCMXCIV"   ✓
3999 → "MMMCMXCIX" ✓
```

---

## Related Problems

- LeetCode #13 — Roman to Integer (the reverse direction)
- LeetCode #273 — Integer to English Words (similar digit-decomposition pattern)
