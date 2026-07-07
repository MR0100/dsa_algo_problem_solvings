# 0139 — Word Break

> LeetCode #139 · Difficulty: Medium
> **Categories:** Array, Hash Table, String, Dynamic Programming, Trie, Memoization

---

## Problem Statement

Given a string `s` and a dictionary of strings `wordDict`, return `true` if `s` can be segmented into a space-separated sequence of one or more dictionary words.

**Note** that the same word in the dictionary may be reused multiple times in the segmentation.

**Example 1:**
```
Input: s = "leetcode", wordDict = ["leet","code"]
Output: true
Explanation: Return true because "leetcode" can be segmented as "leet code".
```

**Example 2:**
```
Input: s = "applepenapple", wordDict = ["apple","pen"]
Output: true
Explanation: Return true because "applepenapple" can be segmented as "apple pen apple".
Note that you are allowed to reuse a dictionary word.
```

**Example 3:**
```
Input: s = "catsandog", wordDict = ["cats","dog","sand","and","cat"]
Output: false
```

**Constraints:**
- `1 <= s.length <= 300`
- `1 <= wordDict.length <= 1000`
- `1 <= wordDict[i].length <= 20`
- `s` and `wordDict[i]` consist of only lowercase English letters.
- All the strings of `wordDict` are **unique**.

---

## Company Frequency

| Company    | Frequency       | Last Reported |
|------------|-----------------|---------------|
| Amazon     | ★★★★★ Very High | 2024          |
| Facebook   | ★★★★★ Very High | 2024          |
| Google     | ★★★★☆ High      | 2024          |
| Microsoft  | ★★★★☆ High      | 2024          |
| Apple      | ★★★☆☆ Medium    | 2023          |
| Bloomberg  | ★★★☆☆ Medium    | 2023          |

> ⚠️ Frequency data is crowd-sourced from LeetCode Discuss, Glassdoor, and community interview reports. Treat as a signal, not a guarantee.

---

## DSA Concepts Used

- **1-D Dynamic Programming** — `dp[i]` = "prefix of length i is segmentable"; overlapping suffix subproblems → see [`/dsa/dynamic_programming_1d.md`](/dsa/dynamic_programming_1d.md)
- **Hash Map / Hash Set** — O(1) dictionary membership tests → see [`/dsa/hash_map.md`](/dsa/hash_map.md)
- **BFS on an implicit graph** — split indices are nodes, words are edges; segmentability = reachability → see [`/dsa/graph_bfs_dfs.md`](/dsa/graph_bfs_dfs.md)
- **Trie** — prefix-walking replaces repeated substring hashing → see [`/dsa/trie.md`](/dsa/trie.md)

---

## Approaches Overview

| # | Approach                  | Time         | Space | When to use                                          |
|---|---------------------------|--------------|-------|------------------------------------------------------|
| 1 | Brute Force (Recursion)   | O(2ⁿ·n)      | O(n)  | Never; shows why memoization is needed               |
| 2 | DP Top-Down (Memoization) | O(n²·L)      | O(n)  | Natural evolution of the recursion; easy to derive   |
| 3 | DP Bottom-Up (Optimal)    | O(n²·L)      | O(n)  | The standard interview answer; no recursion overhead |
| 4 | BFS over Indices          | O(n²·L)      | O(n)  | Alternative framing; nice reachability story          |
| 5 | Trie + DP                 | O(n·M)       | O(W)  | Huge dictionaries / long strings; avoids re-hashing   |

*(n = len(s), L = max word length for substring hashing, M = longest word, W = total characters across dictionary words.)*

---

## Approach 1 — Brute Force (Recursion)

### Intuition
Segmenting `s` means choosing the *first* word, then segmenting the rest. So try every possible prefix `s[start:end]`: if it is a dictionary word and the remaining suffix can also be segmented, we're done. This is a full search over all 2^(n−1) ways to place cut points, and identical suffixes get re-solved over and over.

### Algorithm
1. Put `wordDict` into a hash set.
2. Define `canBreak(start)`:
   1. If `start == len(s)` return `true` (nothing left to match).
   2. For each `end` from `start+1` to `len(s)`: if `s[start:end]` is in the set **and** `canBreak(end)`, return `true`.
   3. Return `false`.
3. Answer is `canBreak(0)`.

### Complexity
- **Time:** O(2ⁿ·n) — worst case explores every subset of cut positions, each step doing O(n)-ish substring hashing.
- **Space:** O(n) — recursion depth, plus the word set.

### Code
```go
func bruteForce(s string, wordDict []string) bool {
	words := make(map[string]bool, len(wordDict)) // O(1) word lookup
	for _, w := range wordDict {
		words[w] = true
	}
	var canBreak func(start int) bool
	canBreak = func(start int) bool {
		if start == len(s) {
			return true // consumed the whole string successfully
		}
		for end := start + 1; end <= len(s); end++ {
			// try s[start:end] as the next word, then recurse on the suffix
			if words[s[start:end]] && canBreak(end) {
				return true
			}
		}
		return false // no word fits the front of s[start:]
	}
	return canBreak(0)
}
```

### Dry Run
`s = "leetcode"`, `wordDict = ["leet","code"]` (Example 1):

| call         | end tried | s[start:end] | in dict? | recursion result        |
|--------------|-----------|--------------|----------|-------------------------|
| canBreak(0)  | 1..3      | "l","le","lee" | no     | keep scanning           |
| canBreak(0)  | 4         | "leet"       | yes      | recurse canBreak(4)     |
| canBreak(4)  | 5..7      | "c","co","cod" | no     | keep scanning           |
| canBreak(4)  | 8         | "code"       | yes      | recurse canBreak(8)     |
| canBreak(8)  | —         | —            | —        | start == 8 → **true**   |
| canBreak(4)  | —         | —            | —        | returns **true**        |
| canBreak(0)  | —         | —            | —        | returns **true** ✅      |

---

## Approach 2 — DP Top-Down (Memoization)

### Intuition
`canBreak(start)` depends only on `start` — there are just `n+1` distinct subproblems, yet the brute force may solve each exponentially many times (classic on inputs like `"aaaa...b"`). Cache each suffix's verdict in a memo and every subproblem is computed once.

### Algorithm
1. Build the word set; create `memo` (start → bool).
2. `canBreak(start)`:
   1. `start == len(s)` → `true`.
   2. Memo hit → return the cached value.
   3. Scan `end` as before; on success cache `memo[start] = true` and return.
   4. Cache `memo[start] = false` and return.
3. Answer is `canBreak(0)`.

### Complexity
- **Time:** O(n²·L) — n suffixes × n split points, each split doing an O(L) substring hash (bounded above by O(n³)).
- **Space:** O(n) — memo plus recursion stack.

### Code
```go
func dpTopDown(s string, wordDict []string) bool {
	words := make(map[string]bool, len(wordDict))
	for _, w := range wordDict {
		words[w] = true
	}
	memo := make(map[int]bool, len(s)) // start index → can s[start:] be broken?
	var canBreak func(start int) bool
	canBreak = func(start int) bool {
		if start == len(s) {
			return true // empty suffix always breaks
		}
		if res, ok := memo[start]; ok {
			return res // already solved this suffix
		}
		for end := start + 1; end <= len(s); end++ {
			if words[s[start:end]] && canBreak(end) {
				memo[start] = true // cache success
				return true
			}
		}
		memo[start] = false // cache failure so we never redo this suffix
		return false
	}
	return canBreak(0)
}
```

### Dry Run
`s = "leetcode"` (Example 1). Same call tree as Approach 1, but now every resolved suffix lands in the memo:

| event                    | memo state after            |
|--------------------------|-----------------------------|
| canBreak(8) → true       | {8:true}*                   |
| canBreak(4) → true       | {8:true*, 4:true}           |
| canBreak(0) → true       | {8:true*, 4:true, 0:true}   |

\* index 8 is the `start == len(s)` base case, returned before touching the memo — shown for completeness. On failure-heavy inputs (Example 3, `"catsandog"`), `memo[5..8] = false` entries are what kill the exponential blow-up: `canBreak(7)` ("og") is computed once even though both the "cats·and" and "cat·sand" branches reach it.

---

## Approach 3 — DP Bottom-Up (Optimal)

### Intuition
Flip the recursion into a table over prefixes: `dp[i]` answers "can `s[:i]` be segmented?". The empty prefix trivially can. A longer prefix `s[:i]` can be segmented iff there is some cut `j < i` where the left part `s[:j]` is already known-good and the right part `s[j:i]` is a dictionary word. Filling `i` from 1 to n guarantees every `dp[j]` we consult is final.

### Algorithm
1. Build the word set. Allocate `dp[0..n]`, set `dp[0] = true`.
2. For `i = 1..n`:
   1. For `j = 0..i-1`: if `dp[j] && words[s[j:i]]`, set `dp[i] = true` and break.
3. Return `dp[n]`.

### Complexity
- **Time:** O(n²·L) — all (j, i) pairs, each with an O(L) substring hash. (Bounding `i−j` by the max word length of 20 makes it effectively O(n·20·L).)
- **Space:** O(n) — the boolean table.

### Code
```go
func dpBottomUp(s string, wordDict []string) bool {
	words := make(map[string]bool, len(wordDict))
	for _, w := range wordDict {
		words[w] = true
	}
	n := len(s)
	dp := make([]bool, n+1) // dp[i] ⇔ s[:i] is segmentable
	dp[0] = true            // empty prefix is trivially segmentable
	for i := 1; i <= n; i++ {
		for j := 0; j < i; j++ {
			// s[:i] works if s[:j] works and the gap s[j:i] is a word
			if dp[j] && words[s[j:i]] {
				dp[i] = true
				break // one witness is enough
			}
		}
	}
	return dp[n]
}
```

### Dry Run
`s = "leetcode"`, `wordDict = ["leet","code"]` (Example 1), n = 8:

| i | witness j found       | s[j:i]  | dp[i] | dp array (indices 0..8)                       |
|---|-----------------------|---------|-------|-----------------------------------------------|
| 1 | none                  | —       | F     | T F                                            |
| 2 | none                  | —       | F     | T F F                                          |
| 3 | none                  | —       | F     | T F F F                                        |
| 4 | j=0 (dp[0]=T)         | "leet"  | **T** | T F F F T                                      |
| 5 | none                  | —       | F     | T F F F T F                                    |
| 6 | none                  | —       | F     | T F F F T F F                                  |
| 7 | none                  | —       | F     | T F F F T F F F                                |
| 8 | j=4 (dp[4]=T)         | "code"  | **T** | T F F F T F F F T                              |

`dp[8] = true` → **true** ✅

---

## Approach 4 — BFS over Indices

### Intuition
Model split positions `0..n` as graph nodes with an edge `start → end` whenever `s[start:end]` is a dictionary word. "Can `s` be segmented?" becomes "is node `n` reachable from node 0?" — plain BFS. The `visited` set plays exactly the role of the DP memo: each index is expanded once.

### Algorithm
1. Build the word set; queue `{0}`; mark 0 visited.
2. Pop `start`; for each `end` in `(start, n]` with `s[start:end]` in the set:
   1. If `end == n`, return `true`.
   2. If `end` unvisited, mark and enqueue it.
3. Queue exhausted → return `false`.

### Complexity
- **Time:** O(n²·L) — each of the n indices is expanded once over up to n edges, each edge test hashing an O(L) substring.
- **Space:** O(n) — queue plus visited array.

### Code
```go
func bfsApproach(s string, wordDict []string) bool {
	words := make(map[string]bool, len(wordDict))
	for _, w := range wordDict {
		words[w] = true
	}
	n := len(s)
	if n == 0 {
		return true // empty string trivially breaks
	}
	visited := make([]bool, n+1) // indices already expanded (avoid rework)
	queue := []int{0}            // frontier of reachable split points
	visited[0] = true
	for len(queue) > 0 {
		start := queue[0]
		queue = queue[1:] // pop front
		for end := start + 1; end <= n; end++ {
			if !words[s[start:end]] {
				continue // no edge start→end
			}
			if end == n {
				return true // reached the end of the string
			}
			if !visited[end] {
				visited[end] = true // mark before enqueue to prevent duplicates
				queue = append(queue, end)
			}
		}
	}
	return false // end index unreachable
}
```

### Dry Run
`s = "leetcode"` (Example 1):

| step | popped start | word edges found     | end reached | queue after | visited          |
|------|--------------|----------------------|-------------|-------------|------------------|
| 1    | 0            | "leet" → 4           | no          | [4]         | {0, 4}           |
| 2    | 4            | "code" → 8           | 8 == n      | —           | return **true** ✅ |

---

## Approach 5 — Trie + DP

### Intuition
The DP's inner loop hashes each candidate substring `s[j:i]` from scratch. A trie of the dictionary shares that work: starting at breakable index `j`, walk characters `s[j], s[j+1], …` down the trie **once**; every trie node flagged `isWord` along the way marks a valid new breakable index `i+1`. Matching all words that start at `j` costs one walk of at most `maxWordLen` steps — no substring allocation or hashing at all.

### Algorithm
1. Insert every dictionary word into a trie; mark terminal nodes `isWord`.
2. `dp[0] = true`. For each `j` from 0 to n with `dp[j]` true:
   1. Walk the trie along `s[j:]` character by character.
   2. If the walk falls off the trie, stop.
   3. Whenever the current trie node has `isWord`, set `dp[i+1] = true`.
3. Return `dp[n]`.

### Complexity
- **Time:** O(n·M + W) — building the trie costs W (total dictionary characters); each of the n start indices walks at most M (longest word) trie steps.
- **Space:** O(W + n) — trie nodes plus the dp table.

### Code
```go
func trieDP(s string, wordDict []string) bool {
	type trieNode struct {
		children map[byte]*trieNode
		isWord   bool
	}
	newNode := func() *trieNode { return &trieNode{children: map[byte]*trieNode{}} }

	// Build the trie of dictionary words.
	root := newNode()
	for _, w := range wordDict {
		curr := root
		for i := 0; i < len(w); i++ {
			c := w[i]
			if curr.children[c] == nil {
				curr.children[c] = newNode() // create the path lazily
			}
			curr = curr.children[c]
		}
		curr.isWord = true // mark the end of a full word
	}

	n := len(s)
	dp := make([]bool, n+1) // dp[i] ⇔ s[:i] is segmentable
	dp[0] = true
	for j := 0; j <= n; j++ {
		if !dp[j] {
			continue // can't start a word from an unreachable index
		}
		curr := root
		for i := j; i < n; i++ {
			curr = curr.children[s[i]] // extend the match by one character
			if curr == nil {
				break // no dictionary word continues this way
			}
			if curr.isWord {
				dp[i+1] = true // s[j:i+1] is a word extending a breakable prefix
			}
		}
	}
	return dp[n]
}
```

### Dry Run
`s = "leetcode"` (Example 1). Trie contains paths `l-e-e-t✓` and `c-o-d-e✓`:

| j (dp[j]=T) | trie walk over s[j:]                          | isWord hits | dp updates   |
|-------------|-----------------------------------------------|-------------|--------------|
| 0           | l → e → e → t✓, then 'c' has no child → stop  | at "leet"   | dp[4] = true |
| 4           | c → o → d → e✓, end of string → stop          | at "code"   | dp[8] = true |
| 8           | j == n, walk is empty                         | —           | —            |

`dp[8] = true` → **true** ✅

---

## Key Takeaways

- "Can this string be split into pieces from a set?" is the archetype of **prefix DP**: `dp[i]` over prefix lengths with `dp[0] = true`.
- Memoized recursion, bottom-up tables, and BFS-with-visited are three costumes for the **same O(n²) state graph** — recognize the equivalence and pick the one you narrate best.
- Practical constant-factor win: bound the inner loop by the **maximum word length** (`i − j ≤ 20` here) — turns O(n²) pair scanning into O(n·20).
- A **trie** shines when the same text positions are matched against many words: one walk from each start index replaces per-word hashing.
- This decision problem is the feasibility filter for LeetCode #140 (Word Break II), where you must actually enumerate the sentences.

---

## Related Problems

- LeetCode #140 — Word Break II (enumerate all segmentations)
- LeetCode #472 — Concatenated Words (word break applied to each dictionary word)
- LeetCode #91 — Decode Ways (same prefix-DP shape, counting instead of boolean)
- LeetCode #279 — Perfect Squares (dp[i] over "remaining amount" with a fixed candidate set)
- LeetCode #1043 — Partition Array for Maximum Sum (prefix DP with bounded segment length)
