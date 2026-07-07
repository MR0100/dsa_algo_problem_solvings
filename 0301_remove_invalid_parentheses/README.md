# 0301 — Remove Invalid Parentheses

> LeetCode #301 · Difficulty: Hard
> **Categories:** String, Backtracking, Breadth-First Search

---

## Problem Statement

Given a string `s` that contains parentheses and letters, remove the minimum number of invalid parentheses to make the input string valid.

Return _all the possible results_. You may return the answer in **any order**.

**Example 1:**

```
Input:  s = "()())()"
Output: ["(())()","()()()"]
```

**Example 2:**

```
Input:  s = "(a)())()"
Output: ["(a())()","(a)()()"]
```

**Example 3:**

```
Input:  s = ")("
Output: [""]
```

> Note: for `")("` the only valid minimal result is the empty string `""`; the outputs above are shown after sorting, so an empty slice prints as `[]`.

**Constraints:**

- `1 <= s.length <= 25`
- `s` consists of lowercase English letters and parentheses `'('` and `')'`.
- There will be at most `20` parentheses in `s`.

---

## Company Frequency

| Company   | Frequency        | Last Reported |
|-----------|------------------|---------------|
| Facebook  | ★★★★★ Very High  | 2024          |
| Google    | ★★★★☆ High       | 2023          |
| Amazon    | ★★★☆☆ Medium     | 2023          |
| Microsoft | ★★☆☆☆ Low        | 2022          |

> ⚠️ Frequency data is crowd-sourced from LeetCode Discuss, Glassdoor, and community interview reports. Treat as a signal, not a guarantee.

---

## DSA Concepts Used

- **Breadth-First Search** — treat "number of removals" as depth; the first level with a valid string is the minimal-removal answer → see [`/dsa/graph_bfs_dfs.md`](/dsa/graph_bfs_dfs.md)
- **Backtracking** — enumerate exactly the deletions that spend the mandatory removal quotas, pruning aggressively → see [`/dsa/backtracking.md`](/dsa/backtracking.md)
- **Stack / balance counter** — validate parentheses with a single running open-minus-close count → see [`/dsa/stack.md`](/dsa/stack.md)
- **String manipulation** — build candidate strings by index-wise deletion → see [`/dsa/string_algorithms.md`](/dsa/string_algorithms.md)

---

## Approaches Overview

| # | Approach | Time | Space | When to use |
|---|----------|------|-------|-------------|
| 1 | BFS level-by-level | O(2ⁿ·n) | O(2ⁿ·n) | Intuitive "minimum = shallowest level"; easy to reason about correctness |
| 2 | DFS with precomputed removal counts (Optimal) | O(2ⁿ·n) | O(n) | Best in practice — exact quotas prune the search and cut memory to O(n) |

---

## Approach 1 — BFS Level-by-Level

### Intuition
The minimum number of removals is a shortest-path quantity. Picture a tree whose root is the original string and whose children delete one more character. BFS reaches the shallowest valid strings first, so the first level that contains any valid string holds the entire minimal-removal answer set.

### Algorithm
1. Seed a queue with the original string and a visited set to avoid duplicate states.
2. Process the queue one level at a time. For each string on the level, if it is valid, mark `found = true` and collect it.
3. Once `found` is true, stop expanding — deeper levels only remove more characters.
4. If no string on the level was valid, generate all children (delete each single parenthesis once), enqueue unseen ones, and descend to the next level.
5. Return the collected valid strings.

### Complexity
- **Time:** O(2ⁿ · n) — up to 2ⁿ distinct subsequences may be visited, each validated in O(n).
- **Space:** O(2ⁿ · n) — the visited set and queue can hold exponentially many strings.

### Code
```go
func bfs(s string) []string {
	result := []string{}          // valid strings found on the minimal level
	visited := map[string]bool{}  // strings already enqueued (dedupe)
	queue := []string{s}          // BFS frontier, starts with the whole string
	visited[s] = true             // mark the root as seen
	found := false                // becomes true once a valid string appears

	for len(queue) > 0 { // process the tree level by level
		next := []string{} // children to explore at the next depth
		for _, cur := range queue {
			if isValidParen(cur) { // this string needs no more removals
				result = append(result, cur)
				found = true // remember that this whole level is the answer level
			}
			if found {
				continue // once found on this level, never expand deeper
			}
			// Generate all strings with exactly one more character removed.
			for i := 0; i < len(cur); i++ {
				c := cur[i]
				if c != '(' && c != ')' {
					continue // only removing parentheses can help validity
				}
				child := cur[:i] + cur[i+1:] // delete character i
				if !visited[child] {
					visited[child] = true // avoid re-expanding duplicates
					next = append(next, child)
				}
			}
		}
		if found {
			break // minimal level fully collected — deeper levels remove more
		}
		queue = next // descend one level (one additional removal)
	}
	sort.Strings(result) // deterministic order for stable test output
	return result
}
```

### Dry Run
Input `s = "()())()"`.

| Level (removals) | Queue sample | Any valid? | Action |
|---|---|---|---|
| 0 | `"()())()"` | No (extra `)` at index 4) | Expand: delete each paren once |
| 1 | `"))())()"`, ... , `"()()()"`, `"(())()"`, ... | **Yes** — `"()()()"` and `"(())()"` are valid | Collect all valid on this level, stop |

Result after sort: `[(())() ()()()]`. Both required exactly one removal, so this level is minimal.

---

## Approach 2 — DFS with Precomputed Removal Counts (Optimal)

### Intuition
A single left-to-right pass pins down the minimum removals exactly: every `)` that has no available `(` before it must be removed (`rightRem`), and any `(` still unmatched at the end must be removed (`leftRem`). Knowing these quotas, we only build strings that delete precisely that many parentheses, and by skipping the "delete" branch when the quota is exhausted (and validating balance at the end) we generate only minimal, distinct valid strings.

### Algorithm
1. Pass once over `s` to compute `leftRem` (unmatched `(`) and `rightRem` (unmatched `)`).
2. DFS over indices, carrying the string built so far, the current `open` balance, and remaining quotas `lRem`, `rRem`.
3. Prune whenever `open`, `lRem`, or `rRem` goes negative.
4. At each paren you may **delete** it (spending the matching quota) or **keep** it (updating `open`). Letters are always kept.
5. At the end, record the built string when both quotas are zero and `open` is zero (balanced).

### Complexity
- **Time:** O(2ⁿ · n) worst case, but the exact quotas prune the tree heavily in practice.
- **Space:** O(n) recursion depth, plus the output set.

### Code
```go
func dfsBacktrack(s string) []string {
	// One pass to compute the mandatory removal quotas.
	leftRem, rightRem := 0, 0
	for i := 0; i < len(s); i++ {
		switch s[i] {
		case '(':
			leftRem++ // tentatively an unmatched '('
		case ')':
			if leftRem > 0 {
				leftRem-- // this ')' matches a pending '('
			} else {
				rightRem++ // ')' with nothing to match → must remove
			}
		}
	}

	resultSet := map[string]bool{} // dedupe defensively
	var dfs func(idx int, built string, open, lRem, rRem int)
	dfs = func(idx int, built string, open, lRem, rRem int) {
		if lRem < 0 || rRem < 0 || open < 0 {
			return // pruned: removed too many or unbalanced prefix
		}
		if idx == len(s) {
			if lRem == 0 && rRem == 0 && open == 0 {
				resultSet[built] = true // exact quotas used, balanced → valid
			}
			return
		}
		c := s[idx]
		// Option A: delete a parenthesis, spending the matching quota.
		if c == '(' && lRem > 0 {
			dfs(idx+1, built, open, lRem-1, rRem)
		} else if c == ')' && rRem > 0 {
			dfs(idx+1, built, open, lRem, rRem-1)
		}
		// Option B: keep the current character.
		switch c {
		case '(':
			dfs(idx+1, built+"(", open+1, lRem, rRem) // one more open
		case ')':
			dfs(idx+1, built+")", open-1, lRem, rRem) // closes an open
		default:
			dfs(idx+1, built+string(c), open, lRem, rRem) // letter: passthrough
		}
	}
	dfs(0, "", 0, leftRem, rightRem)

	result := make([]string, 0, len(resultSet))
	for str := range resultSet {
		result = append(result, str)
	}
	sort.Strings(result) // deterministic order for test output
	return result
}
```

### Dry Run
Input `s = "()())()"`.

1. **Quota pass:** `(` → leftRem 1; `)` matches → leftRem 0; `(` → leftRem 1; `)` matches → leftRem 0; `)` no open → rightRem 1; `(` → leftRem 1; `)` matches → leftRem 0. Final: `leftRem = 0`, `rightRem = 1`.
2. **DFS** must delete exactly one `)`. The deletable `)` positions that yield a balanced string are index 4 (giving `"()()()"`... depending on branch) and the second `)`, producing the two distinct valid strings.

| built (end state) | lRem | rRem | open | recorded? |
|---|---|---|---|---|
| `()()()`  | 0 | 0 | 0 | ✅ |
| `(())()`  | 0 | 0 | 0 | ✅ |
| `())()` (kept extra `)`) | 0 | 1 | – | ❌ quota unused |

Result after sort: `[(())() ()()()]`.

---

## Key Takeaways

- "Minimum number of X operations" often maps to **BFS depth** — the first level with a solution is minimal by construction.
- Computing exact **removal quotas** in one pass turns an open-ended search into a bounded one, enabling strong pruning.
- A single **balance counter** (`open − close`, never going negative) validates parentheses in O(n) O(1).
- Skipping consecutive duplicate deletions (or using a set) keeps results **distinct** without post-filtering.

---

## Related Problems

- LeetCode #20 — Valid Parentheses (the validity check used here)
- LeetCode #22 — Generate Parentheses (enumerate valid strings via backtracking)
- LeetCode #32 — Longest Valid Parentheses (balance-counter DP)
- LeetCode #921 — Minimum Add to Make Parentheses Valid (quota counting only)
- LeetCode #1249 — Minimum Remove to Make Valid Parentheses (single valid result)
