# 0247 — Strobogrammatic Number II

> LeetCode #247 · Difficulty: Medium
> **Categories:** Array, String, Recursion

---

## Problem Statement

Given an integer `n`, return all the **strobogrammatic numbers** that are of length `n`. You may return the answer in **any order**.

A **strobogrammatic number** is a number that looks the same when rotated `180` degrees (looked at upside down).

**Example 1:**

```
Input: n = 2
Output: ["11","69","88","96"]
```

**Example 2:**

```
Input: n = 1
Output: ["0","1","8"]
```

**Constraints:**

- `1 <= n <= 14`

---

## Company Frequency

| Company    | Frequency        | Last Reported |
|------------|------------------|---------------|
| Google     | ★★★★☆ High       | 2024          |
| Facebook   | ★★★☆☆ Medium     | 2023          |
| Uber       | ★★☆☆☆ Low        | 2022          |
| Microsoft  | ★★☆☆☆ Low        | 2022          |

> ⚠️ Frequency data is crowd-sourced from LeetCode Discuss, Glassdoor, and
> community interview reports. Treat as a signal, not a guarantee.

---

## DSA Concepts Used

- **Backtracking / Recursion** — build the number from the outermost pair inward, delegating the inner core to a recursive subproblem of length `n-2` → see [`/dsa/backtracking.md`](/dsa/backtracking.md)
- **String Algorithms** — every result is a mirror pair wrapped around a shorter strobogrammatic core → see [`/dsa/string_algorithms.md`](/dsa/string_algorithms.md)

---

## Approaches Overview

| # | Approach | Time | Space | When to use |
|---|----------|------|-------|-------------|
| 1 | Recursive Build (Optimal) | O(5^(n/2)) | O(5^(n/2)·n) | Cleanest to reason about; natural n → n-2 recursion |
| 2 | Iterative Layer-by-Layer | O(5^(n/2)) | O(5^(n/2)·n) | Avoids recursion; grows from the center outward |

---

## Approach 1 — Recursive Build (Optimal)

### Intuition

A strobogrammatic number is symmetric under 180° rotation. Build it from the outside in: the outermost pair of characters must be a mirror pair `(0,0),(1,1),(8,8),(6,9),(9,6)`, wrapped around a strobogrammatic **core** of length `n-2`. That core is exactly the same subproblem two sizes smaller. Base cases: length 0 gives `[""]`, length 1 gives `["0","1","8"]`. The only special rule: the **outermost** pair may not start with `0` (leading zero), but inner pairs may.

### Algorithm

1. `helper(n, total)` returns all strobogrammatic strings of length `n`; `total` is the original requested length so we can detect the outer layer.
2. Base: `n==0 → [""]`; `n==1 → ["0","1","8"]`.
3. Recurse: `inner = helper(n-2, total)`.
4. For each `core` in `inner`, and each mirror pair `(a,b)`:
   - If `n == total` and `a == '0'`, skip (leading zero).
   - Else append `a + core + b`.
5. Return the collected list.

### Complexity

- **Time:** O(5^(n/2)) — about five choices per outer pair over `n/2` pairs.
- **Space:** O(5^(n/2)·n) — all results plus recursion depth `n/2`.

### Code

```go
var strobPairs = [][2]byte{
	{'0', '0'},
	{'1', '1'},
	{'6', '9'},
	{'8', '8'},
	{'9', '6'},
}

func strobHelper(n, total int) []string {
	if n == 0 {
		return []string{""}
	}
	if n == 1 {
		return []string{"0", "1", "8"}
	}

	inner := strobHelper(n-2, total)
	result := []string{}

	for _, core := range inner {
		for _, p := range strobPairs {
			if n == total && p[0] == '0' {
				continue
			}
			result = append(result, string(p[0])+core+string(p[1]))
		}
	}
	return result
}

func recursiveBuild(n int) []string {
	return strobHelper(n, n)
}
```

### Dry Run

Trace `recursiveBuild(2)` → `strobHelper(2, 2)`:

| Step | Action | State |
|------|--------|-------|
| 1 | `n=2, total=2`, recurse `strobHelper(0,2)` | inner = `[""]` |
| 2 | core=`""`, pair `(0,0)` | `n==total && '0'` → **skip** |
| 3 | core=`""`, pair `(1,1)` | append `"11"` |
| 4 | core=`""`, pair `(6,9)` | append `"69"` |
| 5 | core=`""`, pair `(8,8)` | append `"88"` |
| 6 | core=`""`, pair `(9,6)` | append `"96"` |

Result = `["11","69","88","96"]`. ✓

---

## Approach 2 — Iterative Layer-by-Layer Build

### Intuition

Same construction, no recursion: start from the innermost layer (`""` for even `n`, `["0","1","8"]` for odd `n`) and repeatedly wrap every current string in each mirror pair, adding two characters per round, until the strings reach length `n`. On the final (outermost) round, forbid the `(0,0)` pair to avoid a leading zero.

### Algorithm

1. Seed `current = [""]` if `n` even, else `["0","1","8"]`; set `length` accordingly (0 or 1).
2. While `length < n`: `length += 2`; wrap each string in every pair; if `length == n`, skip pairs starting with `0`.
3. Return `current`.

### Complexity

- **Time:** O(5^(n/2)).
- **Space:** O(5^(n/2)·n).

### Code

```go
func iterativeBuild(n int) []string {
	var current []string
	length := 0
	if n%2 == 0 {
		current = []string{""}
	} else {
		current = []string{"0", "1", "8"}
		length = 1
	}

	for length < n {
		length += 2
		next := []string{}
		for _, core := range current {
			for _, p := range strobPairs {
				if length == n && p[0] == '0' {
					continue
				}
				next = append(next, string(p[0])+core+string(p[1]))
			}
		}
		current = next
	}
	return current
}
```

### Dry Run

Trace `iterativeBuild(2)`:

| Step | length | current before | Action | current after |
|------|--------|----------------|--------|---------------|
| 1 | 0 | `[""]` | enter loop, `length=2` | — |
| 2 | 2 | `[""]` | wrap `""`; skip `(0,0)`; add `11,69,88,96` | `["11","69","88","96"]` |
| 3 | 2 | — | `length == n`, loop ends | `["11","69","88","96"]` |

Result = `["11","69","88","96"]`. ✓

---

## Key Takeaways

- Strobogrammatic strings are built **outside-in** by wrapping a shorter core in a mirror pair — a textbook `n → n-2` recursion.
- The leading-zero rule applies **only to the outermost pair**; inner `(0,0)` is fine, which is why the `total` (or `length == n`) guard matters.
- The center of an odd-length number must be self-symmetric: one of `0,1,8`.

---

## Related Problems

- LeetCode #246 — Strobogrammatic Number (validate a single number)
- LeetCode #248 — Strobogrammatic Number III (count in a range)
- LeetCode #22 — Generate Parentheses (recursive constructive enumeration)
- LeetCode #17 — Letter Combinations of a Phone Number (build-up recursion)
