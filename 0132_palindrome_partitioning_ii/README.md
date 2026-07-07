# 0132 — Palindrome Partitioning II

> LeetCode #132 · Difficulty: Hard
> **Categories:** String, Dynamic Programming

---

## Problem Statement

Given a string `s`, partition `s` such that every substring of the partition is a **palindrome**.

Return *the **minimum** cuts needed for a palindrome partitioning of* `s`.

**Example 1:**
```
Input: s = "aab"
Output: 1
Explanation: The palindrome partitioning ["aa","b"] could be produced using 1 cut.
```

**Example 2:**
```
Input: s = "a"
Output: 0
```

**Example 3:**
```
Input: s = "ab"
Output: 1
```

**Constraints:**
- `1 <= s.length <= 2000`
- `s` consists of lowercase English letters only.

---

## Company Frequency

| Company   | Frequency        | Last Reported |
|-----------|------------------|---------------|
| Google    | ★★★★☆ High       | 2024          |
| Amazon    | ★★★★☆ High       | 2024          |
| Microsoft | ★★★☆☆ Medium     | 2023          |
| Apple     | ★★☆☆☆ Low        | 2023          |
| Bloomberg | ★★☆☆☆ Low        | 2023          |

> ⚠️ Frequency data is crowd-sourced from LeetCode Discuss, Glassdoor, and
> community interview reports. Treat as a signal, not a guarantee.

---

## DSA Concepts Used

- **Dynamic Programming (1D)** — `cuts[i]` = minimum cuts for a prefix; each state built from smaller prefixes → see [`/dsa/dynamic_programming_1d.md`](/dsa/dynamic_programming_1d.md)
- **Dynamic Programming (2D)** — the `isPal[i][j]` palindrome table with recurrence `s[i]==s[j] && isPal[i+1][j-1]` → see [`/dsa/dynamic_programming_2d.md`](/dsa/dynamic_programming_2d.md)
- **Two Pointers** — palindrome verification and center expansion both move a symmetric pair of indices → see [`/dsa/two_pointers.md`](/dsa/two_pointers.md)

---

## Approaches Overview

| # | Approach | Time | Space | When to use |
|---|----------|------|-------|-------------|
| 1 | Brute-Force Recursion | O(n · 2ⁿ) | O(n) | Only to explain the search space; times out for n ≥ ~25 |
| 2 | Top-Down DP (Memoization) | O(n²) | O(n²) | Natural refinement of the recursion; easy to derive |
| 3 | Bottom-Up DP | O(n²) | O(n²) | Standard iterative interview answer; no recursion |
| 4 | Expand Around Center (Optimal) | O(n²) | O(n) | Best space; n = 2000 makes the n×n table (4M bools) avoidable |

---

## Approach 1 — Brute-Force Recursion

### Intuition
If the whole string is already a palindrome, the answer is 0 cuts. Otherwise, the first cut must split off some palindromic prefix, and the total cost is that 1 cut plus the minimum cuts for the remaining suffix. Trying every palindromic prefix and recursing explores every possible partition — the minimum over all of them is the answer. This is exactly LeetCode 131's enumeration, collapsed to a single number.

### Algorithm
1. Define `solve(t)`:
   1. If `t` is a palindrome (two-pointer scan), return 0.
   2. Initialize `best = len(t) - 1` (cutting into single characters always works).
   3. For every split point `i` in `[1, len(t)-1]`: if `t[:i]` is a palindrome, compute `1 + solve(t[i:])` and keep the minimum.
   4. Return `best`.
2. Answer is `solve(s)`.

### Complexity
- **Time:** O(n · 2ⁿ) — each of the n−1 gaps is independently cut or not, giving up to 2ⁿ⁻¹ partitions; each prefix palindrome check costs O(n).
- **Space:** O(n) — recursion depth; Go substrings share the original backing array, so no copies are made.

### Code
```go
func bruteForce(s string) int {
	// isPalStr checks a whole string with two converging pointers.
	isPalStr := func(t string) bool {
		for i, j := 0, len(t)-1; i < j; i, j = i+1, j-1 {
			if t[i] != t[j] { // mismatch → not a palindrome
				return false
			}
		}
		return true
	}

	var solve func(t string) int
	solve = func(t string) int {
		if isPalStr(t) {
			return 0 // whole remainder is one palindrome: no cut needed
		}
		best := len(t) - 1 // cutting into single characters always works
		for i := 1; i < len(t); i++ {
			if isPalStr(t[:i]) { // prefix must be a palindrome to cut here
				if c := 1 + solve(t[i:]); c < best {
					best = c // found a partition with fewer cuts
				}
			}
		}
		return best
	}

	return solve(s)
}
```

### Dry Run — `s = "aab"`

| Step | Call | Palindrome check | Action | Returned value |
|------|------|------------------|--------|----------------|
| 1 | solve("aab") | "aab"? no (`a≠b`) | try prefixes; `best = 2` | pending |
| 2 | — prefix `"a"` | "a"? yes | recurse solve("ab") | pending |
| 3 | solve("ab") | "ab"? no | try prefixes; `best = 1` | pending |
| 4 | — prefix `"a"` | "a"? yes | recurse solve("b") | pending |
| 5 | solve("b") | "b"? yes | palindrome base case | **0** |
| 6 | back in solve("ab") | — | candidate `1 + 0 = 1`; `best = 1` | **1** |
| 7 | back in solve("aab") | — | candidate `1 + 1 = 2`; `best = 2` | pending |
| 8 | — prefix `"aa"` | "aa"? yes | recurse solve("b") → 0 | pending |
| 9 | back in solve("aab") | — | candidate `1 + 0 = 1`; `best = 1` | **1** |

Final answer: `1` (partition `["aa","b"]`). ✅

---

## Approach 2 — Top-Down DP (Memoization)

### Intuition
The brute force recomputes the same suffixes exponentially often, but the minimum cuts for `s[start:]` do not depend on anything before `start` — there are only n distinct subproblems. Two memoizations fix everything: (a) an O(n²) precomputed `isPal[i][j]` table replaces repeated palindrome scans, and (b) a `memo[start]` array caches each suffix's answer.

### Algorithm
1. Build `isPal[i][j]` bottom-up: single characters true; for longer windows, `isPal[i][j] = s[i]==s[j] && (length==2 || isPal[i+1][j-1])`, filled by increasing length.
2. Define `solve(start)`:
   1. If `isPal[start][n-1]`, return 0 (suffix is one palindrome).
   2. If memoized, return `memo[start]`.
   3. For each `end` in `[start, n-2]` with `isPal[start][end]`: candidate `1 + solve(end+1)`; keep the minimum.
   4. Memoize and return.
3. Answer is `solve(0)`.

### Complexity
- **Time:** O(n²) — n suffix states, each trying at most n first pieces with O(1) palindrome lookups; table construction is also O(n²).
- **Space:** O(n²) — the palindrome table dominates; memo and recursion stack are O(n).

### Code
```go
func dpTopDown(s string) int {
	n := len(s)

	// isPal[i][j] == true iff s[i..j] is a palindrome (bottom-up by length).
	isPal := make([][]bool, n)
	for i := range isPal {
		isPal[i] = make([]bool, n)
		isPal[i][i] = true // single characters are palindromes
	}
	for length := 2; length <= n; length++ {
		for i := 0; i+length-1 < n; i++ {
			j := i + length - 1
			if s[i] == s[j] {
				// ends match and the inside is a palindrome (or empty)
				isPal[i][j] = length == 2 || isPal[i+1][j-1]
			}
		}
	}

	memo := make([]int, n)
	for i := range memo {
		memo[i] = -1 // -1 marks "not computed yet"
	}

	var solve func(start int) int
	solve = func(start int) int {
		if isPal[start][n-1] {
			return 0 // whole suffix is a palindrome: zero cuts
		}
		if memo[start] != -1 {
			return memo[start] // reuse previously solved suffix
		}
		best := n - 1 - start // upper bound: cut every position
		for end := start; end < n-1; end++ {
			if isPal[start][end] { // first piece s[start..end] is valid
				if c := 1 + solve(end+1); c < best {
					best = c
				}
			}
		}
		memo[start] = best
		return best
	}

	return solve(0)
}
```

### Dry Run — `s = "aab"`

Palindrome table (true cells): `[0,0] [1,1] [2,2] [0,1]("aa")`; `[1,2]("ab")` and `[0,2]("aab")` are false.

| Step | Call | `isPal[start][2]`? | Transition tried | Result | `memo` |
|------|------|--------------------|------------------|--------|--------|
| 1 | solve(0) | `isPal[0][2]` false | end=0 (`"a"` ok) → solve(1) | pending | `[-1,-1,-1]` |
| 2 | solve(1) | `isPal[1][2]` false | end=1 (`"a"` ok) → solve(2) | pending | `[-1,-1,-1]` |
| 3 | solve(2) | `isPal[2][2]` **true** | base case | **0** | `[-1,-1,-1]` |
| 4 | solve(1) resumes | — | candidate `1+0=1`; no more ends | **1** | `[-1,1,-1]` |
| 5 | solve(0) resumes | — | candidate `1+1=2`; best=2 | pending | `[-1,1,-1]` |
| 6 | solve(0), end=1 | `isPal[0][1]` **true** (`"aa"`) | → solve(2) = 0 (base case) | candidate `1+0=1` | `[-1,1,-1]` |
| 7 | solve(0) finishes | — | best = min(2, 1) = 1 | **1** | `[1,1,-1]` |

Final answer: `1`. ✅

---

## Approach 3 — Bottom-Up DP

### Intuition
Flip the direction: instead of asking "how do I cut this suffix?", build answers for growing **prefixes**. Let `cuts[i]` be the minimum cuts for `s[0..i]`. In an optimal partition of that prefix, the *last* piece is some palindrome `s[j..i]`; everything before it is the prefix `s[0..j-1]`, already solved optimally. So `cuts[i] = min over valid j of cuts[j-1] + 1`, and `cuts[i] = 0` outright when `s[0..i]` is itself a palindrome.

### Algorithm
1. Build the `isPal` table exactly as in Approach 2.
2. For `i` from 0 to n−1:
   1. If `isPal[0][i]`, set `cuts[i] = 0` and continue.
   2. Else initialize `cuts[i] = i` (all single characters) and for each `j` in `[1, i]`: if `isPal[j][i]` and `cuts[j-1] + 1 < cuts[i]`, update `cuts[i]`.
3. Return `cuts[n-1]`.

### Complexity
- **Time:** O(n²) — the (i, j) double loop and the table fill are both quadratic.
- **Space:** O(n²) — the palindrome table; the `cuts` array itself is O(n).

### Code
```go
func dpBottomUp(s string) int {
	n := len(s)

	// palindrome table, same construction as dpTopDown
	isPal := make([][]bool, n)
	for i := range isPal {
		isPal[i] = make([]bool, n)
		isPal[i][i] = true
	}
	for length := 2; length <= n; length++ {
		for i := 0; i+length-1 < n; i++ {
			j := i + length - 1
			if s[i] == s[j] {
				isPal[i][j] = length == 2 || isPal[i+1][j-1]
			}
		}
	}

	cuts := make([]int, n) // cuts[i] = min cuts for prefix s[0..i]
	for i := 0; i < n; i++ {
		if isPal[0][i] {
			cuts[i] = 0 // whole prefix is one palindrome
			continue
		}
		cuts[i] = i // worst case: i cuts → i+1 single characters
		for j := 1; j <= i; j++ {
			// last piece s[j..i] must be a palindrome; prefix s[0..j-1]
			// is already optimally solved in cuts[j-1]
			if isPal[j][i] && cuts[j-1]+1 < cuts[i] {
				cuts[i] = cuts[j-1] + 1
			}
		}
	}

	return cuts[n-1]
}
```

### Dry Run — `s = "aab"`

| i | Prefix | `isPal[0][i]`? | j tried | Last piece `s[j..i]` | Palindrome? | Candidate `cuts[j-1]+1` | `cuts` after |
|---|--------|----------------|---------|----------------------|-------------|--------------------------|--------------|
| 0 | `"a"` | true | — | — | — | — | `[0, _, _]` |
| 1 | `"aa"` | true | — | — | — | — | `[0, 0, _]` |
| 2 | `"aab"` | false | init | — | — | `cuts[2] = 2` | `[0, 0, 2]` |
| 2 | | | j=1 | `"ab"` | no | skip | `[0, 0, 2]` |
| 2 | | | j=2 | `"b"` | yes | `cuts[1]+1 = 1` | `[0, 0, 1]` |

Final answer: `cuts[2] = 1`. ✅

---

## Approach 4 — Expand Around Center (Optimal)

### Intuition
For n = 2000 the n×n table means 4,000,000 booleans — correct, but avoidable. Every palindrome has a center: a character (odd length) or a gap between characters (even length), i.e. 2n−1 centers total. Expanding outward from each center enumerates **every palindromic substring exactly once, in order of discovery**. Each time a palindrome `s[j..k]` is certified, treat it as the *last piece* of a partition of the prefix ending at `k` and relax: `cuts[k+1] = min(cuts[k+1], cuts[j] + 1)`. The DP array is indexed by prefix *length*, with sentinel `cuts[0] = -1` so a palindrome spanning the whole prefix yields `-1 + 1 = 0` cuts.

### Algorithm
1. Allocate `cuts[0..n]` with `cuts[i] = i - 1` (worst case; `cuts[0] = -1` sentinel).
2. For each center `c` in `[0, n-1]`:
   1. Expand `(l, r) = (c, c)` — odd-length palindromes — while `l >= 0 && r < n && s[l] == s[r]`: relax `cuts[r+1]` with `cuts[l] + 1`, then `l--, r++`.
   2. Expand `(l, r) = (c, c+1)` — even-length palindromes — same loop.
3. Return `cuts[n]`.

### Complexity
- **Time:** O(n²) — 2n−1 centers, each expanding at most n/2 steps; every step does O(1) work.
- **Space:** O(n) — only the `cuts` array; the palindrome table is never materialized.

### Code
```go
func expandAroundCenter(s string) int {
	n := len(s)

	// cuts[i] = min cuts for prefix s[:i]; cuts[0] = -1 is a sentinel so
	// that a palindrome covering the whole prefix gives cuts = -1 + 1 = 0.
	cuts := make([]int, n+1)
	for i := 0; i <= n; i++ {
		cuts[i] = i - 1 // worst case: prefix of length i needs i-1 cuts
	}

	// expand grows a palindrome outward from (l, r) while ends match,
	// relaxing the cuts array for every palindrome it certifies.
	expand := func(l, r int) {
		for l >= 0 && r < n && s[l] == s[r] {
			// s[l..r] is a palindrome: use it as the final piece of the
			// prefix s[:r+1]; everything before it costs cuts[l], +1 cut.
			if cuts[l]+1 < cuts[r+1] {
				cuts[r+1] = cuts[l] + 1
			}
			l-- // widen the window symmetrically
			r++
		}
	}

	for c := 0; c < n; c++ {
		expand(c, c)   // odd-length palindromes centered at character c
		expand(c, c+1) // even-length palindromes centered at gap (c, c+1)
	}

	return cuts[n]
}
```

### Dry Run — `s = "aab"`

Initial `cuts = [-1, 0, 1, 2]` (index = prefix length).

| Step | Center | (l, r) | `s[l..r]` | Palindrome? | Relaxation | `cuts` after |
|------|--------|--------|-----------|-------------|------------|--------------|
| 1 | c=0 odd | (0,0) | `"a"` | yes | `cuts[1] = min(0, cuts[0]+1=0)` | `[-1, 0, 1, 2]` |
| 2 | c=0 odd | (-1,1) | — | l<0, stop | — | `[-1, 0, 1, 2]` |
| 3 | c=0 even | (0,1) | `"aa"` | yes | `cuts[2] = min(1, cuts[0]+1=0)` → **0** | `[-1, 0, 0, 2]` |
| 4 | c=0 even | (-1,2) | — | l<0, stop | — | `[-1, 0, 0, 2]` |
| 5 | c=1 odd | (1,1) | `"a"` | yes | `cuts[2] = min(0, cuts[1]+1=1)` → keep 0 | `[-1, 0, 0, 2]` |
| 6 | c=1 odd | (0,2) | `"aab"` | `a≠b`, stop | — | `[-1, 0, 0, 2]` |
| 7 | c=1 even | (1,2) | `"ab"` | `a≠b`, stop | — | `[-1, 0, 0, 2]` |
| 8 | c=2 odd | (2,2) | `"b"` | yes | `cuts[3] = min(2, cuts[2]+1=1)` → **1** | `[-1, 0, 0, 1]` |
| 9 | c=2 even | (2,3) | — | r=3 out of range | — | `[-1, 0, 0, 1]` |

Final answer: `cuts[3] = 1`. ✅

---

## Key Takeaways

- **Counting/minimizing over partitions ⇒ DP, enumerating them ⇒ backtracking.** LeetCode 131 vs 132 is the canonical pair: same cut structure, but asking for a *number* collapses the exponential search into O(n²) states.
- **`cuts[i] = min(cuts[j-1] + 1)` over palindromic last pieces `s[j..i]`** — "fix the last piece" is the standard way to linearize a partition DP.
- **Sentinel `cuts[0] = -1`** elegantly makes a whole-prefix palindrome cost 0 without a special case.
- **Center expansion replaces the O(n²)-space palindrome table** whenever you only need palindromes *streamed* rather than *queried at random* — an O(n²) → O(n) space trick worth remembering.
- The `isPal` bottom-up table (`s[i]==s[j] && isPal[i+1][j-1]`, by increasing length) is a reusable component across #5, #131, #132, #647.

---

## Related Problems

- LeetCode #131 — Palindrome Partitioning (enumerate all partitions; same palindrome table)
- LeetCode #5 — Longest Palindromic Substring (center expansion / DP table)
- LeetCode #647 — Palindromic Substrings (count palindromes via the same machinery)
- LeetCode #139 — Word Break (same prefix-DP shape with a dictionary predicate)
- LeetCode #1278 — Palindrome Partitioning III (min changes to make k palindromic pieces)
- LeetCode #1745 — Palindrome Partitioning IV (split into exactly three palindromes)
