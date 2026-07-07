# 0269 — Alien Dictionary

> LeetCode #269 · Difficulty: Hard
> **Categories:** Array, String, Depth-First Search, Breadth-First Search, Graph, Topological Sort

---

## Problem Statement

There is a new alien language that uses the English alphabet. However, the order
among the letters is unknown to you.

You are given a list of strings `words` from the alien language's dictionary. Now
it is claimed that the strings in `words` are **sorted lexicographically** by the
rules of this new language.

If this claim is incorrect, and the given arrangement of strings in `words`
cannot correspond to any order of letters, return `""`.

Otherwise, return a string of the unique letters in the new alien language sorted
in **lexicographically increasing order** by the new language's rules. If there
are multiple solutions, return **any of them**.

**Example 1:**
```
Input: words = ["wrt","wrf","er","ett","rftt"]
Output: "wertf"
```

**Example 2:**
```
Input: words = ["z","x"]
Output: "zx"
```

**Example 3:**
```
Input: words = ["z","x","z"]
Output: ""
```
Explanation: The order is invalid, so return "".

**Constraints:**
- `1 <= words.length <= 100`
- `1 <= words[i].length <= 100`
- `words[i]` consists of only lowercase English letters.

**Note on invalid prefixes:** If `words[i]` is longer than `words[i+1]` but is a
prefix of it (e.g. `["abc","ab"]`), the ordering is impossible → return `""`.

---

## Company Frequency

| Company    | Frequency        | Last Reported |
|------------|------------------|---------------|
| Amazon     | ★★★★★ Very High  | 2024          |
| Google     | ★★★★☆ High       | 2024          |
| Facebook   | ★★★★☆ High       | 2024          |
| Airbnb     | ★★★☆☆ Medium     | 2023          |
| Microsoft  | ★★★☆☆ Medium     | 2023          |

> ⚠️ Frequency data is crowd-sourced from LeetCode Discuss, Glassdoor, and
> community interview reports. Treat as a signal, not a guarantee.

---

## DSA Concepts Used
- **Topological Sort** — order letters subject to precedence constraints → see [`/dsa/topological_sort.md`](/dsa/topological_sort.md)
- **Graph BFS/DFS** — build and traverse the letter precedence graph → see [`/dsa/graph_bfs_dfs.md`](/dsa/graph_bfs_dfs.md)
- **Queue / Deque** — Kahn's algorithm uses a ready-queue of in-degree-0 nodes → see [`/dsa/queue_deque.md`](/dsa/queue_deque.md)
- **Hash Map / Set** — adjacency lists and in-degree bookkeeping → see [`/dsa/hash_map.md`](/dsa/hash_map.md)

---

## Approaches Overview

| # | Approach | Time | Space | When to use |
|---|----------|------|-------|-------------|
| 1 | BFS Topological Sort (Kahn) | O(C + V + E) | O(V + E) | Easy cycle detection via leftover count |
| 2 | DFS Topological Sort (Optimal) | O(C + V + E) | O(V + E) | Reverse post-order + colour cycle detection |

> `C` = total characters across all words; `V` ≤ 26 letters; `E` ≤ 26² edges.

---

## Approach 1 — BFS Topological Sort (Kahn's Algorithm)

### Intuition
The dictionary is sorted by an unknown alphabet. Comparing two **adjacent**
words, the first position where they differ reveals one ordering rule: the
earlier word's char precedes the later word's char. Collect all such rules into
a directed graph over letters; any topological order of that graph is a valid
alien alphabet. A cycle means the constraints contradict → `""`.

### Algorithm
1. Seed every appearing letter with in-degree 0 (so isolated letters still
   appear in the output).
2. For each adjacent pair, find the first differing chars `c1,c2` and add edge
   `c1 → c2` (once). If `word1` is longer than `word2` but is its prefix, the
   arrangement is invalid → return `""`.
3. Kahn: enqueue all in-degree-0 letters (kept sorted for a deterministic,
   lexicographically-smallest order); repeatedly pop a letter, append it, and
   decrement neighbours, enqueuing any that hit 0.
4. If the result covers every letter, return it; otherwise a cycle exists → `""`.

### Complexity
- **Time:** O(C + V + E) — building the graph reads all characters; the sort
  visits every node and edge.
- **Space:** O(V + E) — adjacency + in-degree over ≤ 26 letters.

### Code
```go
func bfsTopoSort(words []string) string {
	adj := make(map[byte]map[byte]struct{}) // c -> set of letters after c
	indeg := make(map[byte]int)             // c -> number of prerequisites

	// Seed all appearing letters with in-degree 0.
	for _, w := range words {
		for i := 0; i < len(w); i++ {
			if _, ok := adj[w[i]]; !ok {
				adj[w[i]] = make(map[byte]struct{})
				indeg[w[i]] = 0
			}
		}
	}

	// Derive ordering rules from adjacent word pairs.
	for i := 0; i+1 < len(words); i++ {
		w1, w2 := words[i], words[i+1]
		minLen := len(w1)
		if len(w2) < minLen {
			minLen = len(w2)
		}
		j := 0
		for j < minLen && w1[j] == w2[j] {
			j++ // skip the common prefix
		}
		if j == minLen {
			// No differing char within the shared length. Invalid only if the
			// first word is the longer one ("abc" cannot precede "ab").
			if len(w1) > len(w2) {
				return ""
			}
			continue
		}
		c1, c2 := w1[j], w2[j] // first difference gives rule c1 -> c2
		if _, exists := adj[c1][c2]; !exists {
			adj[c1][c2] = struct{}{}
			indeg[c2]++ // c2 gains a prerequisite
		}
	}

	// Kahn's algorithm. We keep the ready set sorted so the output is
	// deterministic (lexicographically smallest valid order).
	ready := make([]byte, 0)
	for c, d := range indeg {
		if d == 0 {
			ready = append(ready, c) // no prerequisites -> ready
		}
	}
	sort.Slice(ready, func(i, j int) bool { return ready[i] < ready[j] })
	result := make([]byte, 0, len(indeg))
	for len(ready) > 0 {
		c := ready[0] // smallest available letter
		ready = ready[1:]
		result = append(result, c) // commit this letter
		freed := make([]byte, 0)
		for next := range adj[c] {
			indeg[next]--
			if indeg[next] == 0 {
				freed = append(freed, next) // all its prereqs placed
			}
		}
		sort.Slice(freed, func(i, j int) bool { return freed[i] < freed[j] })
		ready = append(ready, freed...)
		sort.Slice(ready, func(i, j int) bool { return ready[i] < ready[j] })
	}
	if len(result) < len(indeg) {
		return "" // some letters never freed -> cycle
	}
	return string(result)
}
```

### Dry Run
Input `["wrt","wrf","er","ett","rftt"]`. Derived edges from adjacent pairs:

| Pair            | First diff | Edge   |
|-----------------|------------|--------|
| wrt vs wrf      | t,f        | t → f  |
| wrf vs er       | w,e        | w → e  |
| er vs ett       | r,t        | r → t  |
| ett vs rftt     | e,r        | e → r  |

In-degrees: `w:0, e:1, r:1, t:1, f:1`. Kahn:

| Step | ready (sorted) | pop | append | frees      | result |
|------|----------------|-----|--------|------------|--------|
| 1    | [w]            | w   | w      | e (→0)     | w      |
| 2    | [e]            | e   | e      | r (→0)     | we     |
| 3    | [r]            | r   | r      | t (→0)     | wer    |
| 4    | [t]            | t   | t      | f (→0)     | wert   |
| 5    | [f]            | f   | f      | —          | wertf  |

All 5 letters placed → return **"wertf"**. ✅

---

## Approach 2 — DFS Topological Sort (Optimal)

### Intuition
A topological order is the **reverse of DFS finish times**. Colour nodes white
(unseen) / grey (on the current path) / black (done). Meeting a grey node again
is a back edge → cycle → `""`. Otherwise append each node after its descendants
and reverse the list at the end.

### Algorithm
1. Build the same precedence graph and letter set from adjacent word pairs
   (handling the invalid-prefix case exactly as in Approach 1).
2. DFS each unvisited letter (in sorted order for determinism); a grey re-visit
   signals a cycle. Post-append blacks.
3. Reverse the post-order to obtain the alphabet.

### Complexity
- **Time:** O(C + V + E) — build + traversal, `V` ≤ 26, `E` ≤ 26².
- **Space:** O(V + E) — graph + recursion stack.

### Code
```go
func dfsTopoSort(words []string) string {
	adj := make(map[byte]map[byte]struct{})
	letters := make(map[byte]struct{})
	for _, w := range words {
		for i := 0; i < len(w); i++ {
			letters[w[i]] = struct{}{}
			if _, ok := adj[w[i]]; !ok {
				adj[w[i]] = make(map[byte]struct{})
			}
		}
	}
	for i := 0; i+1 < len(words); i++ {
		w1, w2 := words[i], words[i+1]
		minLen := len(w1)
		if len(w2) < minLen {
			minLen = len(w2)
		}
		j := 0
		for j < minLen && w1[j] == w2[j] {
			j++
		}
		if j == minLen {
			if len(w1) > len(w2) {
				return "" // invalid prefix ordering
			}
			continue
		}
		adj[w1[j]][w2[j]] = struct{}{}
	}

	const (
		white = 0 // unvisited
		grey  = 1 // on current DFS path
		black = 2 // fully processed
	)
	color := make(map[byte]int)
	order := make([]byte, 0, len(letters))
	// sortedKeys returns a map's byte keys in ascending order for determinism.
	sortedKeys := func(m map[byte]struct{}) []byte {
		ks := make([]byte, 0, len(m))
		for k := range m {
			ks = append(ks, k)
		}
		sort.Slice(ks, func(i, j int) bool { return ks[i] < ks[j] })
		return ks
	}
	var dfs func(c byte) bool
	dfs = func(c byte) bool {
		color[c] = grey // enter the recursion stack
		for _, next := range sortedKeys(adj[c]) {
			if color[next] == grey {
				return false // back edge -> cycle
			}
			if color[next] == white && !dfs(next) {
				return false // cycle found deeper
			}
		}
		color[c] = black         // finished
		order = append(order, c) // post-order append
		return true
	}
	for _, c := range sortedKeys(letters) {
		if color[c] == white {
			if !dfs(c) {
				return "" // cyclic constraints
			}
		}
	}
	// order is reverse topological; flip it.
	for i, j := 0, len(order)-1; i < j; i, j = i+1, j-1 {
		order[i], order[j] = order[j], order[i]
	}
	return string(order)
}
```

### Dry Run
Same graph: edges `t→f, w→e, r→t, e→r`. DFS visits sorted letters; take the
chain from `w`:

| Call        | Action                              | order (post) |
|-------------|-------------------------------------|--------------|
| dfs(e)      | e→r→t→f, bottom-out at f            | [f]          |
| finish t    | append t                            | [f,t]        |
| finish r    | append r                            | [f,t,r]      |
| finish e    | append e                            | [f,t,r,e]    |
| dfs(w)      | w→e (black), finish w               | [f,t,r,e,w]  |

Post-order `[f,t,r,e,w]`; reverse → **"wertf"**. ✅ (No grey re-visit ⇒ no cycle.)

---

## Key Takeaways
- Adjacent-pair comparison yields **exactly one** ordering rule per pair (the first differing char); characters after it carry no information.
- The **invalid-prefix** edge case (`["abc","ab"]`) must return `""` — easy to miss.
- Seed **all** appearing letters into the graph so isolated letters are not dropped from the answer.
- Kahn (BFS) detects cycles by a short output; DFS detects them via a grey (on-stack) re-visit.
- Sorting the ready-set / neighbours makes the output deterministic and lexicographically smallest.

---

## Related Problems
- LeetCode #207 — Course Schedule (cycle detection via topological sort)
- LeetCode #210 — Course Schedule II (return a valid order)
- LeetCode #310 — Minimum Height Trees (peeling by degree, BFS layering)
- LeetCode #444 — Sequence Reconstruction (unique topological order check)
