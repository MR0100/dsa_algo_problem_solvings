# 0161 — One Edit Distance

> LeetCode #161 · Difficulty: Medium (Premium)
> **Categories:** Two Pointers, String, Dynamic Programming

---

## Problem Statement

Given two strings `s` and `t`, return `true` if they are both one edit distance apart, otherwise return `false`.

A string `s` is said to be one distance apart from a string `t` if you can:

- Insert **exactly one** character into `s` to get `t`.
- Delete **exactly one** character from `s` to get `t`.
- Replace **exactly one** character of `s` with **a different character** to get `t`.

**Example 1:**
```
Input: s = "ab", t = "acb"
Output: true
Explanation: We can insert 'c' into s to get t.
```

**Example 2:**
```
Input: s = "", t = ""
Output: false
Explanation: We cannot get t from s by only one step.
```

**Constraints:**
- `0 <= s.length, t.length <= 10^4`
- `s` and `t` consist of lowercase letters, uppercase letters, and digits.

---

## Company Frequency

| Company   | Frequency       | Last Reported |
|-----------|-----------------|---------------|
| Meta      | ★★★★★ Very High | 2024          |
| Uber      | ★★★★☆ High      | 2024          |
| Microsoft | ★★★☆☆ Medium    | 2024          |
| Amazon    | ★★★☆☆ Medium    | 2023          |
| Google    | ★★★☆☆ Medium    | 2023          |
| Snap      | ★★☆☆☆ Low       | 2023          |

> ⚠️ Frequency data is crowd-sourced from LeetCode Discuss, Glassdoor, and community interview reports. Treat as a signal, not a guarantee.

---

## DSA Concepts Used

- **Two Pointers** — walking both strings simultaneously and "spending" a single allowed divergence detects the one edit in one pass → see [`/dsa/two_pointers.md`](/dsa/two_pointers.md)
- **Dynamic Programming (2D)** — the brute force computes the full Levenshtein edit-distance table and checks whether it equals 1 → see [`/dsa/dynamic_programming_2d.md`](/dsa/dynamic_programming_2d.md)
- **String Algorithms** — prefix/suffix structure of strings that differ by one edit (common prefix + one repair + common suffix) → see [`/dsa/string_algorithms.md`](/dsa/string_algorithms.md)

---

## Approaches Overview

| # | Approach                             | Time      | Space   | When to use                                                    |
|---|--------------------------------------|-----------|---------|----------------------------------------------------------------|
| 1 | Brute Force (Full Edit-Distance DP)  | O(m·n)    | O(m·n)  | Never for this problem; shows the link to LeetCode #72         |
| 2 | First Mismatch + Suffix Comparison   | O(m + n)  | O(1)    | Clean and easy to reason about; great whiteboard version       |
| 3 | One-Pass Two Pointers (Optimal)      | O(min(m,n)) | O(1)  | Always — single pass, constant space, no slicing needed        |

---

## Approach 1 — Brute Force (Full Edit-Distance DP)

### Intuition
"One edit distance apart" is literally the statement "the Levenshtein edit distance between `s` and `t` equals exactly 1". Edit distance (LeetCode #72) is a classic 2-D DP, so the most direct — if massively over-powered — solution is to fill the whole table and test `dp[m][n] == 1`. Equal strings have distance 0, which correctly maps to `false`: zero edits is not one edit.

### Algorithm
1. Let `m = len(s)`, `n = len(t)`. Allocate `dp` of size `(m+1) × (n+1)` where `dp[i][j]` = edit distance between the prefixes `s[:i]` and `t[:j]`.
2. Base cases: `dp[i][0] = i` (delete all `i` characters) and `dp[0][j] = j` (insert all `j` characters).
3. For every `i` from 1 to `m` and `j` from 1 to `n`:
   1. If `s[i-1] == t[j-1]`, no edit is needed at this position: `dp[i][j] = dp[i-1][j-1]`.
   2. Otherwise `dp[i][j] = 1 + min(dp[i-1][j], dp[i][j-1], dp[i-1][j-1])` — the cheapest of delete, insert, replace.
4. Return `dp[m][n] == 1`.

### Complexity
- **Time:** O(m·n) — every one of the (m+1)·(n+1) table cells is computed once with O(1) work.
- **Space:** O(m·n) — the full table is stored (a rolling row would cut this to O(n), but the point of this approach is clarity).

### Code
```go
func bruteForce(s, t string) bool {
	m, n := len(s), len(t)
	dp := make([][]int, m+1) // dp[i][j] = edit distance of s[:i] vs t[:j]
	for i := range dp {
		dp[i] = make([]int, n+1)
		dp[i][0] = i // turning s[:i] into "" needs i deletions
	}
	for j := 0; j <= n; j++ {
		dp[0][j] = j // turning "" into t[:j] needs j insertions
	}
	for i := 1; i <= m; i++ {
		for j := 1; j <= n; j++ {
			if s[i-1] == t[j-1] {
				dp[i][j] = dp[i-1][j-1] // chars match → no edit needed here
			} else {
				del := dp[i-1][j]   // delete s[i-1]
				ins := dp[i][j-1]   // insert t[j-1]
				rep := dp[i-1][j-1] // replace s[i-1] with t[j-1]
				dp[i][j] = 1 + min(min(del, ins), rep)
			}
		}
	}
	return dp[m][n] == 1 // exactly one edit — 0 (equal) must return false
}
```

### Dry Run
`s = "ab"`, `t = "acb"` (Example 1). Rows are prefixes of `s`, columns prefixes of `t`:

| dp        | "" (j=0) | "a" (j=1) | "ac" (j=2) | "acb" (j=3) | how the row is filled                                        |
|-----------|----------|-----------|------------|-------------|--------------------------------------------------------------|
| "" (i=0)  | 0        | 1         | 2          | 3           | base case: j insertions                                      |
| "a" (i=1) | 1        | **0** ('a'='a' → diag 0) | **1** ('a'≠'c' → 1+min(2,0,1)) | **2** ('a'≠'b' → 1+min(3,1,2)) | mix of match and 1+min |
| "ab" (i=2)| 2        | **1** ('b'≠'a' → 1+min(0,2,1)) | **1** ('b'≠'c' → 1+min(1,1,0)) | **1** ('b'='b' → diag 1) | final cell = 1 |

`dp[2][3] = 1` → return **true** ✅ (the single edit is inserting `'c'`).

For `s = ""`, `t = ""` (Example 2): `dp[0][0] = 0`, and `0 != 1` → **false** ✅.

---

## Approach 2 — First Mismatch + Suffix Comparison

### Intuition
Two strings that are one edit apart must look like: *common prefix* + *the one repaired position* + *identical remainder*. So find the first index where they disagree. What happens there depends only on the lengths:

- equal lengths → the mismatch must be a **replacement**, so the suffixes *after* it must match exactly;
- `s` shorter by one → the mismatched character of `t` was **inserted**, so `s`'s suffix from the mismatch must equal `t`'s suffix after it;
- `s` longer by one → symmetric **deletion** case.

If no mismatch exists in the overlap, the strings can only be one edit apart if one of them has exactly one extra trailing character.

### Algorithm
1. If `|m − n| > 1`, return `false` — a single insert/delete changes length by at most 1.
2. Walk `i` from 0 over the shorter length; stop at the first `i` with `s[i] != t[i]`.
3. At that first mismatch return:
   1. `s[i+1:] == t[i+1:]` when `m == n` (replace `s[i]`);
   2. `s[i:] == t[i+1:]` when `m < n` (insert `t[i]` into `s`);
   3. `s[i+1:] == t[i:]` when `m > n` (delete `s[i]`).
4. If the whole overlap matched, return `|m − n| == 1` — equal strings (diff 0) are **not** one edit apart.

### Complexity
- **Time:** O(m + n) — one scan to the first mismatch plus one suffix comparison; each character is examined at most twice.
- **Space:** O(1) — Go string slicing (`s[i:]`) creates a view over the same bytes, not a copy.

### Code
```go
func suffixCompare(s, t string) bool {
	m, n := len(s), len(t)
	diff := m - n
	if diff < 0 {
		diff = -diff // absolute length difference
	}
	if diff > 1 {
		return false // would need at least two inserts/deletes
	}
	shorter := m
	if n < m {
		shorter = n // only the overlapping prefix can be compared index-wise
	}
	for i := 0; i < shorter; i++ {
		if s[i] != t[i] { // first mismatch — decide which single edit fixes it
			switch {
			case m == n:
				return s[i+1:] == t[i+1:] // replace s[i] with t[i]
			case m < n:
				return s[i:] == t[i+1:] // insert t[i] into s at position i
			default:
				return s[i+1:] == t[i:] // delete s[i] from s
			}
		}
	}
	// The overlap matched entirely: strings are equal (diff == 0 → false,
	// because zero edits is not one edit) or one has a single extra tail char.
	return diff == 1
}
```

### Dry Run
`s = "ab"`, `t = "acb"` (Example 1):

| step | variable state                                | action                                                     |
|------|-----------------------------------------------|------------------------------------------------------------|
| 1    | `m=2, n=3, diff=1`                            | `diff ≤ 1` → continue                                       |
| 2    | `shorter = 2`                                 | compare the overlapping prefix index by index               |
| 3    | `i=0`: `s[0]='a'`, `t[0]='a'`                 | equal → keep walking                                        |
| 4    | `i=1`: `s[1]='b'`, `t[1]='c'`                 | first mismatch; `m < n` → insert case                       |
| 5    | compare `s[1:]="b"` with `t[2:]="b"`          | `"b" == "b"` → return **true** ✅                            |

For `s = ""`, `t = ""` (Example 2): loop body never runs (`shorter = 0`), `diff = 0` → return **false** ✅.

---

## Approach 3 — One-Pass Two Pointers (Optimal)

### Intuition
Approach 2 still re-compares a suffix after the mismatch. Instead, walk pointers `i` (shorter string) and `j` (longer string) together and treat the single allowed edit as a budget: the first time the characters differ, spend it — skip just the longer string's character when lengths differ (insert/delete), or skip both when lengths match (replace). Any *second* difference proves the distance is ≥ 2. Every character is looked at exactly once.

### Algorithm
1. Swap if needed so `s` is the shorter (or equal-length) string; if `n − m > 1`, return `false`.
2. Set `i = j = 0` and `usedEdit = false`.
3. While `i < m` and `j < n`:
   1. If `s[i] == t[j]`, advance both pointers and continue.
   2. Otherwise, if `usedEdit` is already `true`, return `false` (second mismatch).
   3. Set `usedEdit = true`; advance `j` always; advance `i` too only when `m == n` (replacement).
4. After the loop return `usedEdit || n−m == 1` — either the edit was spent and the tails matched, or the strings matched fully and `t` has exactly one leftover trailing character.

### Complexity
- **Time:** O(min(m, n)) — both pointers only move forward; the loop runs at most `min(m, n)` + 1 mismatch steps.
- **Space:** O(1) — two indices and one boolean flag.

### Code
```go
func twoPointers(s, t string) bool {
	m, n := len(s), len(t)
	if m > n {
		s, t = t, s // ensure s is the shorter string
		m, n = n, m
	}
	if n-m > 1 {
		return false // one insert/delete can bridge a gap of at most 1
	}
	i, j := 0, 0
	usedEdit := false // whether the single allowed edit has been consumed
	for i < m && j < n {
		if s[i] == t[j] { // characters agree → advance both pointers
			i++
			j++
			continue
		}
		if usedEdit {
			return false // second mismatch → at least two edits required
		}
		usedEdit = true
		if m == n {
			i++ // equal lengths → this mismatch must be a replacement
		}
		j++ // unequal lengths → skip t[j] (delete from t / insert into s)
	}
	// Either the edit was already spent (tails matched afterwards), or the
	// strings matched fully and t has exactly one leftover character.
	return usedEdit || n-m == 1
}
```

### Dry Run
`s = "ab"`, `t = "acb"` (Example 1); `s` is already the shorter string, `n − m = 1`:

| step | i | j | s[i] vs t[j] | usedEdit | action                                            |
|------|---|---|--------------|----------|---------------------------------------------------|
| 1    | 0 | 0 | 'a' vs 'a'   | false    | match → `i=1, j=1`                                |
| 2    | 1 | 1 | 'b' vs 'c'   | false→**true** | first mismatch, `m≠n` → skip `t[1]`: `j=2`  |
| 3    | 1 | 2 | 'b' vs 'b'   | true     | match → `i=2, j=3`                                |
| 4    | 2 | 3 | loop ends    | true     | return `usedEdit=true` → **true** ✅               |

For `s = ""`, `t = ""` (Example 2): `m=n=0`, loop never runs, `usedEdit=false`, `n−m=0` → `false || false` = **false** ✅.

---

## Key Takeaways

- **"Exactly one edit" ≠ "at most one edit"** — equal strings must return `false`. This edge case (Example 2 is literally two empty strings) is the most common interview slip.
- Strings one edit apart always decompose as *common prefix + single repair + identical suffix* — so the first mismatch tells you everything, and lengths alone decide whether the repair is a replace (`m == n`) or an insert/delete (`|m − n| == 1`).
- The **edit budget** pattern (a `usedEdit` flag spent on the first divergence, instant `false` on the second) generalises to problems like Valid Palindrome II (#680) where you may skip at most one character.
- Normalising by swapping so `s` is always the shorter string halves the number of cases — a recurring trick in two-string problems.
- Recognise reductions: this is edit distance (#72) with the answer capped at 1, which collapses an O(m·n) DP to a linear scan. When a DP's output is restricted to a tiny set, look for a direct combinatorial argument.

---

## Related Problems

- LeetCode #72 — Edit Distance (the general DP this problem specialises)
- LeetCode #583 — Delete Operation for Two Strings (edit distance with deletes only)
- LeetCode #712 — Minimum ASCII Delete Sum for Two Strings (weighted edit-distance variant)
- LeetCode #680 — Valid Palindrome II (same "spend one edit budget" two-pointer pattern)
- LeetCode #392 — Is Subsequence (two-pointer simultaneous walk over two strings)
