# 0131 — Palindrome Partitioning

> LeetCode #131 · Difficulty: Medium
> **Categories:** String, Dynamic Programming, Backtracking

---

## Problem Statement

Given a string `s`, partition `s` such that every substring of the partition is a **palindrome**. Return *all possible palindrome partitionings of* `s`.

**Example 1:**
```
Input: s = "aab"
Output: [["a","a","b"],["aa","b"]]
```

**Example 2:**
```
Input: s = "a"
Output: [["a"]]
```

**Constraints:**
- `1 <= s.length <= 16`
- `s` contains only lowercase English letters.

---

## Company Frequency

| Company   | Frequency        | Last Reported |
|-----------|------------------|---------------|
| Amazon    | ★★★★★ Very High  | 2024          |
| Microsoft | ★★★★☆ High       | 2024          |
| Google    | ★★★★☆ High       | 2024          |
| Meta      | ★★★☆☆ Medium     | 2023          |
| Bloomberg | ★★★☆☆ Medium     | 2023          |

> ⚠️ Frequency data is crowd-sourced from LeetCode Discuss, Glassdoor, and
> community interview reports. Treat as a signal, not a guarantee.

---

## DSA Concepts Used

- **Backtracking** — enumerate all partitions by choosing a palindromic first piece, recursing on the rest, and undoing the choice → see [`/dsa/backtracking.md`](/dsa/backtracking.md)
- **Dynamic Programming (2D)** — `isPal[i][j]` table built from smaller substrings makes every palindrome query O(1) → see [`/dsa/dynamic_programming_2d.md`](/dsa/dynamic_programming_2d.md)
- **Two Pointers** — the on-the-fly palindrome check converges two indices toward the middle → see [`/dsa/two_pointers.md`](/dsa/two_pointers.md)

---

## Approaches Overview

| # | Approach | Time | Space | When to use |
|---|----------|------|-------|-------------|
| 1 | Brute-Force Backtracking | O(n · 2ⁿ) | O(n) aux | Simplest to write; fine for n ≤ 16 |
| 2 | Backtracking + DP Palindrome Table | O(n · 2ⁿ), O(1) per check | O(n²) | Standard interview answer; removes redundant scans |
| 3 | Memoized Suffix Partitions (Optimal reuse) | O(n · 2ⁿ) | O(n · 2ⁿ) | When suffix results are reused / to show top-down DP thinking |

*(The output itself can contain 2ⁿ⁻¹ partitions, so no algorithm can beat exponential time — "optimal" here means eliminating all redundant work.)*

---

## Approach 1 — Brute-Force Backtracking

### Intuition
A partition of `s` is a sequence of cut positions. Walk the string left to right: at index `start`, try every candidate first piece `s[start..end]`. If that piece is a palindrome, commit to it and recursively partition the remainder `s[end+1:]`. When `start` reaches the end, the pieces chosen along the way form one complete valid partition. Undoing each choice after the recursion (backtracking) lets one `path` slice explore every branch.

### Algorithm
1. Initialize `result = []`, `path = []`.
2. Define `dfs(start)`:
   1. If `start == len(s)`, append a **copy** of `path` to `result` and return.
   2. For `end` from `start` to `len(s)-1`:
      1. Check `s[start..end]` with two pointers (`i` from the left, `j` from the right, compare and converge).
      2. If it is a palindrome: append `s[start:end+1]` to `path`, call `dfs(end+1)`, then pop the last element of `path`.
3. Call `dfs(0)` and return `result`.

### Complexity
- **Time:** O(n · 2ⁿ) — there are up to 2ⁿ⁻¹ ways to cut a string of length n (each of the n−1 gaps is cut or not); each explored piece pays an O(n) palindrome scan and each recorded partition costs O(n) to copy.
- **Space:** O(n) auxiliary — recursion depth ≤ n and the shared `path` slice (the exponential `result` is required output, not counted).

### Code
```go
func bruteForceBacktracking(s string) [][]string {
	n := len(s)
	result := [][]string{} // all valid partitions collected here
	path := []string{}     // current partial partition being built

	// isPal checks s[i..j] inclusive with two converging pointers.
	isPal := func(i, j int) bool {
		for i < j {
			if s[i] != s[j] { // mismatch means not a palindrome
				return false
			}
			i++ // move both ends toward the middle
			j--
		}
		return true
	}

	var dfs func(start int)
	dfs = func(start int) {
		if start == n {
			// consumed the whole string: snapshot the path (copy, because
			// path's backing array keeps mutating during backtracking)
			part := make([]string, len(path))
			copy(part, path)
			result = append(result, part)
			return
		}
		for end := start; end < n; end++ {
			if isPal(start, end) { // only recurse on palindromic prefixes
				path = append(path, s[start:end+1]) // choose this piece
				dfs(end + 1)                        // partition the remainder
				path = path[:len(path)-1]           // undo the choice (backtrack)
			}
		}
	}

	dfs(0)
	return result
}
```

### Dry Run — `s = "aab"`

| Step | Call | end | Piece `s[start..end]` | Palindrome? | `path` after action | `result` |
|------|------|-----|----------------------|-------------|---------------------|----------|
| 1 | dfs(0) | 0 | `"a"` | yes | `[a]` | `[]` |
| 2 | dfs(1) | 1 | `"a"` | yes | `[a a]` | `[]` |
| 3 | dfs(2) | 2 | `"b"` | yes | `[a a b]` | `[]` |
| 4 | dfs(3) | — | start==3 → record | — | `[a a b]` | `[[a a b]]` |
| 5 | back in dfs(2) | — | pop `b`; loop ends | — | `[a a]` | `[[a a b]]` |
| 6 | back in dfs(1) | 2 | `"ab"` | **no** | `[a a]` → pop `a` → `[a]` | `[[a a b]]` |
| 7 | back in dfs(0) | 1 | `"aa"` | yes | `[aa]` | `[[a a b]]` |
| 8 | dfs(2) | 2 | `"b"` | yes | `[aa b]` | `[[a a b]]` |
| 9 | dfs(3) | — | start==3 → record | — | `[aa b]` | `[[a a b] [aa b]]` |
| 10 | unwind | 2 | `"aab"` | **no** | `[]` | `[[a a b] [aa b]]` |

Final answer: `[["a","a","b"],["aa","b"]]`. ✅

---

## Approach 2 — Backtracking + DP Palindrome Table

### Intuition
The brute force re-scans the same substrings many times across branches. Palindromicity has optimal substructure: `s[i..j]` is a palindrome **iff** `s[i] == s[j]` and the inner substring `s[i+1..j-1]` is a palindrome. Filling a table bottom-up over increasing substring lengths answers *every possible* palindrome query in O(n²) total preprocessing, turning each check inside the search into an O(1) lookup.

### Algorithm
1. Allocate `isPal[n][n]`; set `isPal[i][i] = true` (length‑1 substrings).
2. For `length` from 2 to n, for each window `[i, j=i+length-1]`:
   1. If `s[i] == s[j]` and (`length == 2` or `isPal[i+1][j-1]`), set `isPal[i][j] = true`.
3. Run the same backtracking as Approach 1, replacing the two-pointer scan with the lookup `isPal[start][end]`.

### Complexity
- **Time:** O(n · 2ⁿ) overall — the exponential number of partitions still dominates, but each palindrome check drops from O(n) to O(1); the table itself is built in O(n²).
- **Space:** O(n²) — the boolean table; plus O(n) recursion depth.

### Code
```go
func dpTableBacktracking(s string) [][]string {
	n := len(s)

	// isPal[i][j] == true iff s[i..j] is a palindrome.
	isPal := make([][]bool, n)
	for i := range isPal {
		isPal[i] = make([]bool, n)
		isPal[i][i] = true // every single character is a palindrome
	}
	// fill by increasing substring length so smaller answers already exist
	for length := 2; length <= n; length++ {
		for i := 0; i+length-1 < n; i++ {
			j := i + length - 1 // right end of the window
			if s[i] == s[j] {
				// inner part must be a palindrome too (or be empty, length 2)
				isPal[i][j] = length == 2 || isPal[i+1][j-1]
			}
		}
	}

	result := [][]string{}
	path := []string{}

	var dfs func(start int)
	dfs = func(start int) {
		if start == n {
			part := make([]string, len(path)) // snapshot the finished partition
			copy(part, path)
			result = append(result, part)
			return
		}
		for end := start; end < n; end++ {
			if isPal[start][end] { // O(1) table lookup instead of a scan
				path = append(path, s[start:end+1])
				dfs(end + 1)
				path = path[:len(path)-1] // backtrack
			}
		}
	}

	dfs(0)
	return result
}
```

### Dry Run — `s = "aab"`

Table construction (`n = 3`):

| Window | `s[i]` vs `s[j]` | Inner check | `isPal[i][j]` |
|--------|------------------|-------------|----------------|
| [0,0] `"a"` | — | length 1 | **true** |
| [1,1] `"a"` | — | length 1 | **true** |
| [2,2] `"b"` | — | length 1 | **true** |
| [0,1] `"aa"` | `a == a` | length 2 | **true** |
| [1,2] `"ab"` | `a != b` | — | false |
| [0,2] `"aab"` | `a != b` | — | false |

Search phase (lookups only, no rescanning):

| Step | Call | Lookup | Result | `path` | `result` |
|------|------|--------|--------|--------|----------|
| 1 | dfs(0) | `isPal[0][0]` | true | `[a]` | `[]` |
| 2 | dfs(1) | `isPal[1][1]` | true | `[a a]` | `[]` |
| 3 | dfs(2) | `isPal[2][2]` | true | `[a a b]` | `[]` |
| 4 | dfs(3) | base case | record | `[a a b]` | `[[a a b]]` |
| 5 | dfs(1) | `isPal[1][2]` | false | `[a]` | `[[a a b]]` |
| 6 | dfs(0) | `isPal[0][1]` | true | `[aa]` | `[[a a b]]` |
| 7 | dfs(2) | `isPal[2][2]` | true | `[aa b]` | `[[a a b]]` |
| 8 | dfs(3) | base case | record | `[aa b]` | `[[a a b] [aa b]]` |
| 9 | dfs(0) | `isPal[0][2]` | false | `[]` | `[[a a b] [aa b]]` |

Final answer: `[["a","a","b"],["aa","b"]]`. ✅

---

## Approach 3 — Memoized Suffix Partitions (Optimal reuse)

### Intuition
`dfs(start)` gets re-entered with the same `start` along many different prefixes, yet the set of partitions of the suffix `s[start:]` never depends on how we arrived. That is a textbook overlapping-subproblem: cache `partitions(start)` once and reuse it. The recurrence is `partitions(start) = { [piece] ++ rest : piece = s[start..end] is a palindrome, rest ∈ partitions(end+1) }`, with base case `partitions(n) = [[]]` (one empty partition).

### Algorithm
1. Build the O(n²) palindrome table exactly as in Approach 2.
2. Define `solve(start)`:
   1. Base case: `start == n` → return `[[]]`.
   2. If `memo[start]` exists, return it.
   3. For each `end` with `isPal[start][end]`: take `piece = s[start:end+1]`, and for every `rest` in `solve(end+1)` emit a fresh slice `piece ++ rest`.
   4. Store the list in `memo[start]` and return it.
3. Answer is `solve(0)`.

### Complexity
- **Time:** O(n · 2ⁿ) — bounded by output size; each suffix's partition list is computed exactly once and shared, so no branch is ever re-explored.
- **Space:** O(n · 2ⁿ) — the memo holds complete partition lists for every suffix (this is the trade-off versus Approach 2's O(n²)).

### Code
```go
func memoizedSuffixes(s string) [][]string {
	n := len(s)

	// palindrome table, identical construction to Approach 2
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

	memo := make(map[int][][]string) // start index → all partitions of s[start:]

	var solve func(start int) [][]string
	solve = func(start int) [][]string {
		if start == n {
			// exactly one way to partition the empty suffix: the empty list
			return [][]string{{}}
		}
		if cached, ok := memo[start]; ok {
			return cached // suffix already fully solved
		}
		res := [][]string{}
		for end := start; end < n; end++ {
			if !isPal[start][end] {
				continue // first piece must be a palindrome
			}
			piece := s[start : end+1]
			for _, rest := range solve(end + 1) {
				// build a fresh slice: piece followed by the cached tail
				part := make([]string, 0, len(rest)+1)
				part = append(part, piece)
				part = append(part, rest...)
				res = append(res, part)
			}
		}
		memo[start] = res // cache before returning
		return res
	}

	return solve(0)
}
```

### Dry Run — `s = "aab"`

| Step | Call | Palindromic first pieces | Sub-result used | Value computed | `memo` after |
|------|------|--------------------------|-----------------|----------------|--------------|
| 1 | solve(0) | `"a"` [0,0], `"aa"` [0,1] | needs solve(1), solve(2) | — (pending) | `{}` |
| 2 | solve(1) | `"a"` [1,1] (`"ab"` fails) | needs solve(2) | — (pending) | `{}` |
| 3 | solve(2) | `"b"` [2,2] | needs solve(3) | — (pending) | `{}` |
| 4 | solve(3) | base case | — | `[[]]` | `{}` |
| 5 | solve(2) resumes | `"b"` ++ `[]` | solve(3) | `[[b]]` | `{2: [[b]]}` |
| 6 | solve(1) resumes | `"a"` ++ `[b]` | solve(2) | `[[a b]]` | `{2, 1: [[a b]]}` |
| 7 | solve(0), piece `"a"` | `"a"` ++ `[a b]` | solve(1) | `[[a a b]]` so far | `{2, 1}` |
| 8 | solve(0), piece `"aa"` | `"aa"` ++ `[b]` | **memo hit** solve(2) | `[[a a b] [aa b]]` | `{2, 1, 0: [[a a b] [aa b]]}` |

Note step 8: `solve(2)` is answered from the memo — the reuse the backtracking versions never get.

Final answer: `[["a","a","b"],["aa","b"]]`. ✅

---

## Key Takeaways

- **"All partitions" ⇒ backtracking skeleton:** choose a first piece, recurse on the remainder, undo. The palindrome condition is just a pruning filter on the choice.
- **Precompute predicates you query repeatedly:** the `isPal[i][j]` DP table (`s[i]==s[j] && isPal[i+1][j-1]`, filled by increasing length) is a reusable building block — it also powers LeetCode 132, 5, and 647.
- **Always copy the path before recording it** — the shared slice keeps mutating; appending it directly is a classic aliasing bug in Go backtracking.
- **Output-sensitive lower bound:** with up to 2ⁿ⁻¹ partitions, exponential time is unavoidable; optimization means removing redundant *checks* (DP table) or redundant *exploration* (suffix memoization), not beating 2ⁿ.
- **Memoize on the suffix index** when the answer for "the rest of the string" is independent of the prefix — a general top-down DP trigger.

---

## Related Problems

- LeetCode #132 — Palindrome Partitioning II (same palindrome table, min-cut DP instead of enumeration)
- LeetCode #5 — Longest Palindromic Substring (same `isPal` expansion/DP machinery)
- LeetCode #647 — Palindromic Substrings (count entries of the same table)
- LeetCode #93 — Restore IP Addresses (identical "cut the string with a validity predicate" backtracking)
- LeetCode #139 — Word Break (partition with a dictionary predicate instead of palindromes)
- LeetCode #2472 — Maximum Number of Non-overlapping Palindrome Substrings (palindrome table + DP)
