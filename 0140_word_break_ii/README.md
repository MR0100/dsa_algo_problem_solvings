# 0140 — Word Break II

> LeetCode #140 · Difficulty: Hard
> **Categories:** Array, Hash Table, String, Dynamic Programming, Backtracking, Trie, Memoization

---

## Problem Statement

Given a string `s` and a dictionary of strings `wordDict`, add spaces in `s` to construct a sentence where each word is a valid dictionary word. Return all such possible sentences in **any order**.

**Note** that the same word in the dictionary may be reused multiple times in the segmentation.

**Example 1:**
```
Input: s = "catsanddog", wordDict = ["cat","cats","and","sand","dog"]
Output: ["cats and dog","cat sand dog"]
```

**Example 2:**
```
Input: s = "pineapplepenapple", wordDict = ["apple","pen","applepen","pine","pineapple"]
Output: ["pine apple pen apple","pineapple pen apple","pine applepen apple"]
Explanation: Note that you are allowed to reuse a dictionary word.
```

**Example 3:**
```
Input: s = "catsandog", wordDict = ["cats","dog","sand","and","cat"]
Output: []
```

**Constraints:**
- `1 <= s.length <= 20`
- `1 <= wordDict.length <= 1000`
- `1 <= wordDict[i].length <= 10`
- `s` and `wordDict[i]` consist of only lowercase English letters.
- All the strings of `wordDict` are **unique**.
- Input is generated in a way that the length of the answer doesn't exceed 10⁵.

---

## Company Frequency

| Company    | Frequency       | Last Reported |
|------------|-----------------|---------------|
| Amazon     | ★★★★☆ High      | 2024          |
| Google     | ★★★★☆ High      | 2024          |
| Facebook   | ★★★★☆ High      | 2024          |
| Microsoft  | ★★★☆☆ Medium    | 2023          |
| Uber       | ★★★☆☆ Medium    | 2023          |
| Snap       | ★★☆☆☆ Low       | 2023          |

> ⚠️ Frequency data is crowd-sourced from LeetCode Discuss, Glassdoor, and community interview reports. Treat as a signal, not a guarantee.

---

## DSA Concepts Used

- **Backtracking** — choose a word, recurse, un-choose; the standard enumeration engine → see [`/dsa/backtracking.md`](/dsa/backtracking.md)
- **Dynamic Programming / Memoization** — suffix→sentences caching turns repeated subtrees into lookups → see [`/dsa/dynamic_programming_1d.md`](/dsa/dynamic_programming_1d.md)
- **Hash Map / Hash Set** — O(1) dictionary membership → see [`/dsa/hash_map.md`](/dsa/hash_map.md)

---

## Approaches Overview

| # | Approach                            | Time            | Space           | When to use                                            |
|---|-------------------------------------|-----------------|-----------------|--------------------------------------------------------|
| 1 | Backtracking (Brute Force)          | O(2ⁿ·n)         | O(n) + output   | n ≤ 20 (as here); simplest correct enumeration          |
| 2 | DP Top-Down (Memoized DFS)          | O(2ⁿ·n) worst   | O(2ⁿ·n)         | Repeated suffixes; the classic interview answer         |
| 3 | DP Bottom-Up (Prefix Sentence Table)| O(2ⁿ·n) worst   | O(2ⁿ·n)         | Shows #139→#140 evolution; feasibility pruning built in |

*(The output itself can be exponential in n, so every correct algorithm is Ω(answer size); the notes distinguish how much redundant work each avoids.)*

---

## Approach 1 — Backtracking (Brute Force)

### Intuition
A sentence is a sequence of choices: "which dictionary word comes next?". Depth-first search over those choices enumerates every sentence exactly once. At position `start`, try each `end` such that `s[start:end]` is a word: push the word, recurse from `end`, then pop it so the next candidate starts from a clean path. Hitting `start == len(s)` means the current path spells the entire string — record it.

Unlike #139 we cannot stop at the first success; we must exhaust the whole tree, which is why the failure case (`"catsandog"`) re-explores the doomed `"og"` tail from both the `cats and` and `cat sand` branches.

### Algorithm
1. Load `wordDict` into a hash set.
2. `dfs(start)`:
   1. If `start == len(s)`: join `path` with spaces, append to `results`, return.
   2. For `end` from `start+1` to `len(s)`:
      1. If `s[start:end]` is a dictionary word: append it to `path`, call `dfs(end)`, then remove it (backtrack).
3. Call `dfs(0)`; return `results`.

### Complexity
- **Time:** O(2ⁿ·n) — up to 2^(n−1) segmentations (e.g. `s="aaaa…"`, dict `{"a","aa",…}`), each sentence assembled in O(n).
- **Space:** O(n) auxiliary — recursion depth plus the path buffer; output storage excluded by convention.

### Code
```go
func backtracking(s string, wordDict []string) []string {
	words := make(map[string]bool, len(wordDict)) // O(1) membership test
	for _, w := range wordDict {
		words[w] = true
	}
	var results []string
	var path []string // words chosen so far on the current root-to-here branch
	var dfs func(start int)
	dfs = func(start int) {
		if start == len(s) {
			results = append(results, strings.Join(path, " ")) // full sentence found
			return
		}
		for end := start + 1; end <= len(s); end++ {
			word := s[start:end]
			if !words[word] {
				continue // this slice is not a dictionary word
			}
			path = append(path, word) // choose
			dfs(end)                  // explore the suffix
			path = path[:len(path)-1] // un-choose (backtrack)
		}
	}
	dfs(0)
	return results
}
```

### Dry Run
`s = "catsanddog"`, `wordDict = ["cat","cats","and","sand","dog"]` (Example 1):

| step | call        | word tried | path after choose        | outcome                          |
|------|-------------|------------|--------------------------|----------------------------------|
| 1    | dfs(0)      | "cat"      | [cat]                    | recurse dfs(3)                   |
| 2    | dfs(3)      | "sand"     | [cat, sand]              | recurse dfs(7)                   |
| 3    | dfs(7)      | "dog"      | [cat, sand, dog]         | recurse dfs(10)                  |
| 4    | dfs(10)     | —          | —                        | record **"cat sand dog"**        |
| 5    | backtrack   | —          | [cat, sand] → [cat]      | dfs(7), dfs(3) exhausted         |
| 6    | dfs(0)      | "cats"     | [cats]                   | recurse dfs(4)                   |
| 7    | dfs(4)      | "and"      | [cats, and]              | recurse dfs(7)                   |
| 8    | dfs(7)      | "dog"      | [cats, and, dog]         | recurse dfs(10)                  |
| 9    | dfs(10)     | —          | —                        | record **"cats and dog"**        |
| 10   | unwind all  | —          | []                       | return both sentences ✅          |

---

## Approach 2 — DP Top-Down (Memoized DFS)

### Intuition
The backtracker revisits the same suffix from different prefixes (step 3 and step 8 above both solve `dfs(7)` = "dog"). The full *set of sentences* for a suffix is a pure function of its start index, so cache it: `sentences(start)` returns every segmentation of `s[start:]`, computed once. Prefixes then combine as `word + " " + tail` for every cached tail. The base case returns `[""]` — one way to segment the empty tail — which makes the gluing rule uniform.

### Algorithm
1. Load the word set; create `memo: map[int][]string`.
2. `sentences(start)`:
   1. Memo hit → return cached slice.
   2. If `start == len(s)` → return `[""]`.
   3. For each `end` with `s[start:end]` in the set: for every `tail` in `sentences(end)`, append `word` (if `tail == ""`) or `word + " " + tail`.
   4. Store in `memo[start]`; return.
3. Return `sentences(0)`.

### Complexity
- **Time:** O(2ⁿ·n) worst case — output-bound; but each suffix's sentence set is computed once, so shared subtrees (like the "dog" tail) cost O(1) on re-visits.
- **Space:** O(2ⁿ·n) — the memo can hold every sentence of every suffix (this is the trade-off versus Approach 1).

### Code
```go
func dpTopDown(s string, wordDict []string) []string {
	words := make(map[string]bool, len(wordDict))
	for _, w := range wordDict {
		words[w] = true
	}
	memo := make(map[int][]string) // start index → all sentences for s[start:]
	var sentences func(start int) []string
	sentences = func(start int) []string {
		if res, ok := memo[start]; ok {
			return res // suffix already solved
		}
		if start == len(s) {
			return []string{""} // one way to break the empty tail: no words at all
		}
		var res []string
		for end := start + 1; end <= len(s); end++ {
			word := s[start:end]
			if !words[word] {
				continue // not a word, no branch here
			}
			for _, tail := range sentences(end) {
				if tail == "" {
					res = append(res, word) // word ends the sentence exactly
				} else {
					res = append(res, word+" "+tail) // glue word before the tail
				}
			}
		}
		memo[start] = res // cache for any other path that reaches this suffix
		return res
	}
	return sentences(0)
}
```

### Dry Run
`s = "catsanddog"` (Example 1). Resolution order of the suffix subproblems:

| suffix solved     | via word(s)       | tails used            | memo[start] result          |
|-------------------|-------------------|-----------------------|-----------------------------|
| sentences(10)     | base case         | —                     | [""]                        |
| sentences(7) "dog"| "dog" + tail ""   | [""]                  | ["dog"]                     |
| sentences(3) "sanddog" | "sand" + tails of 7 | ["dog"]         | ["sand dog"]                |
| sentences(4) "anddog"  | "and" + tails of 7  | ["dog"] (memo hit!) | ["and dog"]              |
| sentences(0)      | "cat" + tails of 3; "cats" + tails of 4 | ["sand dog"], ["and dog"] | ["cat sand dog", "cats and dog"] ✅ |

The memo hit on `sentences(7)` is exactly the work Approach 1 duplicated.

---

## Approach 3 — DP Bottom-Up (Prefix Sentence Table)

### Intuition
Upgrade #139's boolean table to a table of sentence lists: `dp[i]` holds **all** sentences spelling the prefix `s[:i]`. Each `dp[i]` is formed by extending every sentence of a smaller `dp[j]` with the word `s[j:i]`. Two practical guards keep it fast: run the cheap boolean DP first and (a) bail out immediately when the whole string is unbreakable, (b) skip building sentences for any prefix that the boolean table marks dead — no sentence built along an infeasible prefix can ever reach `dp[n]`.

### Algorithm
1. Compute `canBreak[0..n]` exactly as in LeetCode #139.
2. If `!canBreak[n]`, return `[]` (Example 3 exits here without building anything).
3. `dp[0] = [""]`. For `i = 1..n` with `canBreak[i]`:
   1. For each `j < i` with `dp[j]` non-empty and `s[j:i]` in the word set:
      1. For every `sentence` in `dp[j]`: append `word` if `sentence == ""` else `sentence + " " + word` to `dp[i]`.
4. Return `dp[n]`.

### Complexity
- **Time:** O(2ⁿ·n) worst case — the sentence lists themselves can be exponential; the feasibility pre-pass keeps unbreakable inputs at O(n²·L) with zero sentence construction.
- **Space:** O(2ⁿ·n) — sentence lists for every feasible prefix (strictly more retained than top-down, which can drop lists once consumed).

### Code
```go
func dpBottomUp(s string, wordDict []string) []string {
	words := make(map[string]bool, len(wordDict))
	for _, w := range wordDict {
		words[w] = true
	}
	n := len(s)

	// Phase 1: cheap boolean feasibility (exactly LeetCode #139).
	canBreak := make([]bool, n+1)
	canBreak[0] = true
	for i := 1; i <= n; i++ {
		for j := 0; j < i; j++ {
			if canBreak[j] && words[s[j:i]] {
				canBreak[i] = true
				break // one witness suffices
			}
		}
	}
	if !canBreak[n] {
		return []string{} // no segmentation exists; skip the expensive phase
	}

	// Phase 2: build actual sentences only along feasible prefixes.
	dp := make([][]string, n+1) // dp[i] = all sentences spelling s[:i]
	dp[0] = []string{""}        // empty prefix = one empty sentence
	for i := 1; i <= n; i++ {
		if !canBreak[i] {
			continue // dead prefix: building sentences here is wasted work
		}
		for j := 0; j < i; j++ {
			if len(dp[j]) == 0 || !words[s[j:i]] {
				continue // either s[:j] unbuildable or gap isn't a word
			}
			word := s[j:i]
			for _, sentence := range dp[j] {
				if sentence == "" {
					dp[i] = append(dp[i], word) // first word needs no leading space
				} else {
					dp[i] = append(dp[i], sentence+" "+word)
				}
			}
		}
	}
	return dp[n]
}
```

### Dry Run
`s = "catsanddog"` (Example 1), n = 10. Phase 1 gives `canBreak` true at indices 0, 3, 4, 7, 10.

Phase 2 (only feasible `i` shown):

| i  | contributing j | word s[j:i] | dp[j] sentences        | dp[i] after                                   |
|----|----------------|-------------|------------------------|-----------------------------------------------|
| 3  | 0              | "cat"       | [""]                   | ["cat"]                                       |
| 4  | 0              | "cats"      | [""]                   | ["cats"]                                      |
| 7  | 3              | "sand"      | ["cat"]                | ["cat sand"]                                  |
| 7  | 4              | "and"       | ["cats"]               | ["cat sand", "cats and"]                      |
| 10 | 7              | "dog"       | ["cat sand","cats and"]| ["cat sand dog", "cats and dog"] ✅            |

`dp[10]` is the answer: `["cat sand dog","cats and dog"]` (any order accepted).

---

## Key Takeaways

- **#139 is the feasibility oracle for #140** — always run the cheap boolean DP first; it converts the killer no-solution cases (like `"aaaa…b"`) from exponential to quadratic.
- Enumeration problems ("return **all** …") are backtracking at heart; memoization only pays when subproblems (suffixes) are genuinely shared, and the cache must then store *collections*, not booleans.
- The memo base case `[""]` for an empty tail is the clean trick that makes sentence gluing uniform (`word` vs `word + " " + tail`).
- Output-sensitive complexity: when the answer itself can be exponential, quote costs as "O(answer) plus overhead" and focus on eliminating *redundant* work.
- Choose-explore-unchoose (`append` → recurse → re-slice) is the canonical Go backtracking idiom; re-slicing `path[:len(path)-1]` is O(1).

---

## Related Problems

- LeetCode #139 — Word Break (the boolean feasibility version)
- LeetCode #472 — Concatenated Words (word break across the dictionary itself)
- LeetCode #93 — Restore IP Addresses (same enumerate-all-splits backtracking)
- LeetCode #131 — Palindrome Partitioning (enumerate splits with a different predicate)
- LeetCode #282 — Expression Add Operators (path-building DFS with pruning)
