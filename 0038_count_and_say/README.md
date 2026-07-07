# 0038 — Count and Say

> LeetCode #38 · Difficulty: Medium
> **Categories:** String

---

## Problem Statement

The **count-and-say** sequence is a sequence of digit strings defined by the recursive formula:

- `countAndSay(1) = "1"`
- `countAndSay(n)` is the run-length encoding of `countAndSay(n - 1)`.

Run-length encoding (RLE) is a string compression method that works by replacing consecutive identical characters with the concatenation of the character and the marking number of occurrences of this character: e.g., to compress the string `"3322251"`, we replace `"33"` with `"23"`, replace `"222"` with `"32"`, replace `"5"` with `"15"`, and replace `"1"` with `"11"`. Thus the compressed string becomes `"23321511"`.

Return the `n`th term of the count-and-say sequence.

**Example 1**
```
Input:  n = 1
Output: "1"
Explanation:
  This is the base case.
```

**Example 2**
```
Input:  n = 4
Output: "1211"
Explanation:
  countAndSay(1) = "1"
  countAndSay(2) = say "1"         = one 1         = "11"
  countAndSay(3) = say "11"        = two 1s         = "21"
  countAndSay(4) = say "21"        = one 2, one 1   = "1211"
```

**Constraints**
- `1 <= n <= 30`

---

## Company Frequency

| Company   | Frequency       | Last Reported |
|-----------|-----------------|---------------|
| Amazon    | ★★★★☆ High      | 2024          |
| Google    | ★★★★☆ High      | 2024          |
| Microsoft | ★★★☆☆ Medium    | 2024          |
| Meta      | ★★★☆☆ Medium    | 2023          |
| Bloomberg | ★★★☆☆ Medium    | 2023          |
| Apple     | ★★☆☆☆ Low       | 2022          |

> ⚠️ Frequency data is crowd-sourced from LeetCode Discuss, Glassdoor, and community interview reports. Treat as a signal, not a guarantee.

---

## DSA Concepts Used

- **Run-Length Encoding (RLE)** — the core operation: walk a string, count consecutive identical characters, emit `count + char`.
- **String Building** — use `strings.Builder` to avoid O(n²) string concatenation in loops.

---

## Approaches Overview

| # | Approach | Time | Space | When to use |
|---|----------|------|-------|-------------|
| 1 | Iterative ✅ | O(n × L) | O(L) | Preferred; no recursion stack |
| 2 | Recursive | O(n × L) | O(n × L) | Elegant; slight stack overhead |

L = length of the current term (grows by Conway's constant ≈ 1.30 per step).

---

## Approach 1 — Iterative Simulation (Recommended ✅)

### Intuition
Start with `"1"`. Repeatedly apply run-length encoding `n-1` times.

The RLE helper scans the current string with a walk pointer, counts consecutive equal characters, emits `count + char`, then advances by the run length.

### Algorithm
```
result = "1"
for i = 2 to n:
  result = rle(result)
return result

rle(s):
  sb = ""
  j = 0
  while j < len(s):
    ch = s[j]; count = 1
    while j+count < len(s) and s[j+count] == ch: count++
    sb += string(count) + string(ch)
    j += count
  return sb
```

### Complexity
- **Time:** O(n × L) — n-1 RLE applications; each processes a string of length L.
- **Space:** O(L) — the current sequence string; old strings are garbage-collected.

### Code
```go
func iterative(n int) string {
    result := "1"
    for i := 2; i <= n; i++ { result = rle(result) }
    return result
}
func rle(s string) string {
    var sb strings.Builder
    for j := 0; j < len(s); {
        ch := s[j]; count := 1
        for j+count < len(s) && s[j+count] == ch { count++ }
        sb.WriteByte(byte('0' + count)); sb.WriteByte(ch)
        j += count
    }
    return sb.String()
}
```

### Dry Run — building `countAndSay(5)`
```
n=1: "1"
n=2: rle("1") → one 1 → "11"
n=3: rle("11") → two 1s → "21"
n=4: rle("21") → one 2, one 1 → "1211"
n=5: rle("1211") → one 1, one 2, two 1s → "111221"
✓
```

---

## Approach 2 — Recursive

### Intuition
`countAndSay(n) = rle(countAndSay(n-1))` with base case `countAndSay(1) = "1"`.

### Complexity
- **Time:** O(n × L).
- **Space:** O(n × L) — the call stack holds n strings simultaneously.

### Code
```go
func recursive(n int) string {
    if n == 1 {
        return "1"
    }
    return rle(recursive(n - 1))
}

// rle returns the run-length encoding of s.
func rle(s string) string {
    var sb strings.Builder
    j := 0
    for j < len(s) {
        ch := s[j]
        count := 1
        for j+count < len(s) && s[j+count] == ch { // count the run
            count++
        }
        sb.WriteByte(byte('0' + count)) // write count
        sb.WriteByte(ch)                // write digit
        j += count
    }
    return sb.String()
}
```

### Dry Run — `recursive(5)`
Calls unwind to the base case `n==1`, then each level applies `rle` to the value returned from below.

**Descent (recursion goes down to the base case):**

| call | n | action |
|------|---|--------|
| recursive(5) | 5 | needs rle(recursive(4)) |
| recursive(4) | 4 | needs rle(recursive(3)) |
| recursive(3) | 3 | needs rle(recursive(2)) |
| recursive(2) | 2 | needs rle(recursive(1)) |
| recursive(1) | 1 | base case → return `"1"` |

**Return (each level applies rle to the value from below):**

| returning call | rle(input) | says | result |
|----------------|-----------|------|--------|
| recursive(2) | rle("1")    | one 1              | `"11"` |
| recursive(3) | rle("11")   | two 1s             | `"21"` |
| recursive(4) | rle("21")   | one 2, one 1       | `"1211"` |
| recursive(5) | rle("1211") | one 1, one 2, two 1s | `"111221"` |

Final answer: `"111221"` ✓

---

## Key Takeaways

- **`strings.Builder` vs `+=`** — in Go (and Java/Python), naive string concatenation in a loop is O(L²) because strings are immutable and each `+=` allocates a new string. `strings.Builder` amortises to O(L).
- **Run length is always a single digit** — for the count-and-say sequence, runs never exceed 3 characters, so the count is always 1, 2, or 3. This is a property of the sequence, not a general truth about RLE.
- **Conway's constant** — the length of the n-th term grows by a factor approaching 1.3035... per step (Conway's constant). For n=30, the term has ~5 million characters; feasible but large.

---

## Related Problems

- LeetCode #271 — Encode and Decode Strings (RLE in a real-world context)
- LeetCode #443 — String Compression (RLE compression in-place)
