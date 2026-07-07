# 0022 — Generate Parentheses

> LeetCode #22 · Difficulty: Medium
> **Categories:** String, Dynamic Programming, Backtracking

---

## Problem Statement

Given `n` pairs of parentheses, write a function to *generate all combinations of well-formed parentheses*.

**Example 1**
```
Input:  n = 3
Output: ["((()))","(()())","(())()","()(())","()()()"]
```

**Example 2**
```
Input:  n = 1
Output: ["()"]
```

**Constraints**
- `1 <= n <= 8`

---

## Company Frequency

| Company   | Frequency       | Last Reported |
|-----------|-----------------|---------------|
| Amazon    | ★★★★★ Very High | 2024          |
| Google    | ★★★★★ Very High | 2024          |
| Meta      | ★★★★★ Very High | 2024          |
| Microsoft | ★★★★☆ High      | 2024          |
| Bloomberg | ★★★★☆ High      | 2024          |
| Apple     | ★★★☆☆ Medium    | 2023          |
| Adobe     | ★★★☆☆ Medium    | 2023          |
| Uber      | ★★★☆☆ Medium    | 2023          |

> ⚠️ Frequency data is crowd-sourced from LeetCode Discuss, Glassdoor, and community interview reports. Treat as a signal, not a guarantee.

---

## DSA Concepts Used

- **Backtracking** — Approach 2 prunes the search tree by only placing `(` when opens remain and `)` when closes < opens.
- **Dynamic Programming / Catalan Numbers** — Approach 3 builds `dp[k]` from `dp[0..k-1]` using the recursive decomposition: every valid string of k pairs is `(dp[i])dp[j]` where `i+j = k-1`.
- **Combinatorics** — the answer count is the n-th Catalan number: C(n) = C(2n,n)/(n+1).

---

## Approaches Overview

| # | Approach | Time | Space | When to use |
|---|----------|------|-------|-------------|
| 1 | Brute Force | O(2^2n · n) | O(2^2n · n) | Never; exponentially many invalid strings generated |
| 2 | Backtracking ✅ | O(4^n / √n) | O(n) | The canonical interview answer |
| 3 | DP (Catalan decomposition) | O(4^n / √n) | O(4^n / √n) | Demonstrates the Catalan structure; iterative |

---

## Approach 1 — Brute Force

### Intuition
Generate every string of length 2n using `(` and `)`, then filter the valid ones with a balance counter.

### Algorithm
- Enumerate all 2^(2n) binary strings (0 → `(`, 1 → `)`).
- Validate: balance never goes negative and ends at 0.

### Complexity
- **Time:** O(2^(2n) · n) — 4^n strings, each validated in O(n).
- **Space:** O(4^n · n).

### Code
```go
// bruteForce generates every string of length 2n using '(' and ')', then
// filters out invalid ones.
//
// Time:  O(2^(2n) · n) — generate 2^(2n) strings, each validated in O(n).
// Space: O(2^(2n) · n) — all generated strings.
func bruteForce(n int) []string {
	var result []string
	generate(make([]byte, 2*n), 0, n, &result)
	return result
}

func generate(current []byte, pos, n int, result *[]string) {
	if pos == 2*n {
		if isValid(current) {
			*result = append(*result, string(current))
		}
		return
	}
	current[pos] = '('
	generate(current, pos+1, n, result)
	current[pos] = ')'
	generate(current, pos+1, n, result)
}

func isValid(s []byte) bool {
	balance := 0
	for _, ch := range s {
		if ch == '(' {
			balance++
		} else {
			balance--
		}
		if balance < 0 {
			return false
		}
	}
	return balance == 0
}
```

### Dry Run — `n = 2`
`generate` fills 2n = 4 positions, trying `(` then `)` at each, and `isValid` filters the 2^4 = 16 leaves with a balance counter (`balance` never goes < 0, must end at 0). Representative leaves:

| String | balance trace | Valid? |
|--------|---------------|--------|
| `(())` | 1, 2, 1, 0 | ✓ record |
| `()()` | 1, 0, 1, 0 | ✓ record |
| `(()(` | 1, 2, 1, 2 → ends 2 ≠ 0 | ✗ |
| `()((` | 1, 0, 1, 2 → ends 2 ≠ 0 | ✗ |
| `)(()` | −1 on first char | ✗ |
| `))((` | −1 on first char | ✗ |

Only two of the sixteen candidates pass. **Result:** `["(())", "()()"]` ✓

---

## Approach 2 — Backtracking (Recommended ✅)

### Intuition
At each position we have at most two choices, but we only branch when that choice can still lead to a valid string:
- Add `(` if `open < n` — we still have opens to use.
- Add `)` if `close < open` — there is an unmatched `(` to close.

This prunes the tree to exactly the valid strings — no filtering needed.

### Algorithm
```
btHelper(n, open, close, path, result):
  if len(path) == 2n:
    result.append(path)
    return
  if open < n:
    btHelper(n, open+1, close, path+'(', result)
  if close < open:
    btHelper(n, open, close+1, path+')', result)
```

### Complexity
- **Time:** O(4^n / √n) — the n-th Catalan number; the exact count of valid strings times O(n) to copy each.
- **Space:** O(n) — recursion depth 2n + the path buffer.

### Code
```go
func backtracking(n int) []string {
    var result []string
    var bt func(open, close int, path []byte)
    bt = func(open, close int, path []byte) {
        if len(path) == 2*n {
            result = append(result, string(path))
            return
        }
        if open < n  { bt(open+1, close, append(path, '(')) }
        if close < open { bt(open, close+1, append(path, ')')) }
    }
    bt(0, 0, []byte{})
    return result
}
```

### Dry Run — `n = 2`
```
bt(0,0,""):
  open<2 → bt(1,0,"("):
    open<2 → bt(2,0,"(("):
      close<open → bt(2,1,"(()"):
        close<open → bt(2,2,"(())") → record "(())"
    close<open → bt(1,1,"()"):
      open<2 → bt(2,1,"()("):
        close<open → bt(2,2,"()()") → record "()()"
Result: ["(())", "()()"] ✓
```

---

## Approach 3 — Dynamic Programming (Catalan Decomposition)

### Intuition
Every valid parenthesisation of `k` pairs can be uniquely written as:
```
  ( [inner] ) [outer]
```
where `inner` is a valid parenthesisation of `i` pairs and `outer` is a valid parenthesisation of `j = k-1-i` pairs, for `i` from 0 to k-1. This gives the Catalan recurrence: `C(k) = Σ C(i) · C(k-1-i)`.

Build iteratively: `dp[0] = [""]`, `dp[k] = ["("+inner+")"+outer for i in 0..k-1 for inner in dp[i] for outer in dp[k-1-i]]`.

### Complexity
- **Time:** O(4^n / √n) — same total number of strings.
- **Space:** O(4^n / √n) — the full dp table.

### Code
```go
// dpApproach builds dp[k] (all valid strings of k pairs) from dp[0..k-1].
//
// Time:  O(4^n / sqrt(n)) — same count of valid strings to produce.
// Space: O(4^n / sqrt(n)) — all generated strings stored in dp table.
func dpApproach(n int) []string {
	// dp[k] = all valid parenthesisations of k pairs.
	dp := make([][]string, n+1)
	dp[0] = []string{""}

	for k := 1; k <= n; k++ {
		var combinations []string
		for i := 0; i < k; i++ {
			j := k - 1 - i
			for _, inner := range dp[i] {
				for _, outer := range dp[j] {
					combinations = append(combinations, "("+inner+")"+outer)
				}
			}
		}
		dp[k] = combinations
	}
	return dp[n]
}
```

### Dry Run — `n = 2`
Build `dp[0..2]` bottom-up using `dp[k] = "(" + dp[i] + ")" + dp[k-1-i]` for `i` in `0..k-1`:

| k | i (j = k−1−i) | inner (dp[i]) × outer (dp[j]) | strings produced | dp[k] |
|---|---------------|-------------------------------|------------------|-------|
| 0 | — | — | — | `[""]` |
| 1 | i=0, j=0 | `""` × `""` | `"(" + "" + ")" + "" = "()"` | `["()"]` |
| 2 | i=0, j=1 | `""` × `"()"` | `"()" + "()" → "()()"` | — |
| 2 | i=1, j=0 | `"()"` × `""` | `"(()" + ")" → "(())"` | — |

Concatenating the k=2 rows gives `dp[2] = ["()()", "(())"]`. **Result:** `["()()", "(())"]` ✓ (same set as the other approaches, different order).

---

## Key Takeaways

- **The pruning condition** — `open < n` for `(` and `close < open` for `)`. These two guards are the entire validity logic. No post-processing filter needed.
- **Catalan number = answer count** — C(n) = (2n choose n) / (n+1). For n=4: 14, n=5: 42, n=6: 132. Memorise a few for interviews.
- **This is the canonical backtracking problem** — master the `open/close counter` pattern here; it appears in many variants (remove invalid parentheses #301, longest valid parentheses #32).
- **DP decomposition insight** — `(inner)outer` shows that the first character is always `(` and its matching `)` splits the string into two independent valid parts. This decomposition appears again in #10 (regex DP) and #96 (unique BSTs).

---

## Related Problems

- LeetCode #20 — Valid Parentheses (validation, not generation)
- LeetCode #32 — Longest Valid Parentheses (longest valid substring)
- LeetCode #301 — Remove Invalid Parentheses (backtracking with BFS)
- LeetCode #96 — Unique Binary Search Trees (Catalan number recurrence)
