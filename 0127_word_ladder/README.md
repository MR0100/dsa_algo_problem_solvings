# 0127 — Word Ladder

> LeetCode #127 · Difficulty: Hard
> **Categories:** Hash Table, String, Breadth-First Search

---

## Problem Statement

A **transformation sequence** from word `beginWord` to word `endWord` using a dictionary `wordList` is a sequence of words `beginWord -> s1 -> s2 -> ... -> sk` such that every adjacent pair differs by a single letter and every `si` is in `wordList`.

Given `beginWord`, `endWord`, and `wordList`, return the **number of words** in the shortest transformation sequence, or `0` if no sequence exists.

**Example 1:**
```
Input: beginWord = "hit", endWord = "cog", wordList = ["hot","dot","dog","lot","log","cog"]
Output: 5
Explanation: hit → hot → dot → dog → cog, length 5.
```

**Example 2:**
```
Input: beginWord = "hit", endWord = "cog", wordList = ["hot","dot","dog","lot","log"]
Output: 0
```

**Constraints:**
- `1 <= beginWord.length <= 10`
- `endWord.length == beginWord.length`
- `1 <= wordList.length <= 5000`
- All words consist of lowercase English letters.

---

## Company Frequency

| Company   | Frequency       | Last Reported |
|-----------|-----------------|---------------|
| Amazon    | ★★★★★ Very High | 2024          |
| Google    | ★★★★☆ High      | 2024          |
| Microsoft | ★★★★☆ High      | 2024          |
| Facebook  | ★★★☆☆ Medium    | 2024          |

> ⚠️ Frequency data is crowd-sourced from LeetCode Discuss, Glassdoor, and community interview reports. Treat as a signal, not a guarantee.

---

## DSA Concepts Used

- **BFS on implicit graph** — words are nodes, single-letter-difference = edge
- **Bidirectional BFS** — expand from both ends simultaneously for faster convergence

---

## Approaches Overview

| # | Approach            | Time          | Space   | When to use              |
|---|---------------------|---------------|---------|--------------------------|
| 1 | BFS                 | O(N·26·L)     | O(N·L)  | Standard                 |
| 2 | Bidirectional BFS   | O(N·26·L/2)   | O(N·L)  | Larger inputs            |

---

## Approach 1 — BFS

### Intuition
Model as a graph where each word is a node and there's an edge between words differing by one letter. BFS finds the shortest path from `beginWord` to `endWord`.

For each word in the queue, try replacing each character with all 26 letters and check if the result is in `wordSet`.

### Algorithm
1. Build `wordSet`.
2. BFS: queue starts with `[beginWord]`, `steps=1`.
3. Each level: `steps++`. For each word, try all substitutions. Return `steps` if `endWord` found. Mark visited.

### Complexity
- **Time:** O(N · 26 · L) — N words × L positions × 26 chars.
- **Space:** O(N · L)

### Code
```go
func ladderLength(beginWord, endWord string, wordList []string) int {
    wordSet := make(map[string]bool)
    for _, w := range wordList { wordSet[w] = true }
    if !wordSet[endWord] { return 0 }
    queue := []string{beginWord}
    visited := map[string]bool{beginWord: true}
    steps := 1
    for len(queue) > 0 {
        levelSize := len(queue); steps++
        for i := 0; i < levelSize; i++ {
            word := queue[0]; queue = queue[1:]
            bs := []byte(word)
            for pos := 0; pos < len(bs); pos++ {
                orig := bs[pos]
                for c := byte('a'); c <= byte('z'); c++ {
                    if c == orig { continue }
                    bs[pos] = c; next := string(bs)
                    if next == endWord { return steps }
                    if wordSet[next] && !visited[next] { visited[next]=true; queue=append(queue,next) }
                    bs[pos] = orig
                }
            }
        }
    }
    return 0
}
```

### Dry Run
`hit → cog`:

| level | queue        | steps |
|-------|-------------|-------|
| 1     | [hit]       | 2     |
| 2     | [hot]       | 3     |
| 3     | [dot,lot]   | 4     |
| 4     | [dog,log]   | 5     |
| expand dog: next=cog → return 5 |

---

## Approach 2 — Bidirectional BFS

### Intuition
Expand two frontiers simultaneously — one from `beginWord`, one from `endWord`. When they intersect, we found the shortest path. Always expand the smaller frontier (reduces branching).

### Complexity
- **Time:** O(N · 26 · L) but frontier sizes are smaller → faster in practice.
- **Space:** O(N · L)

### Code
```go
// ladderLengthBiBFS solves Word Ladder using bidirectional BFS.
//
// Time:  O(N * 26 * L) but with a smaller constant.
// Space: O(N * L)
func ladderLengthBiBFS(beginWord string, endWord string, wordList []string) int {
	wordSet := make(map[string]bool)
	for _, w := range wordList {
		wordSet[w] = true
	}
	if !wordSet[endWord] {
		return 0
	}

	frontBegin := map[string]bool{beginWord: true}
	frontEnd := map[string]bool{endWord: true}
	steps := 1

	for len(frontBegin) > 0 && len(frontEnd) > 0 {
		// always expand the smaller frontier
		if len(frontBegin) > len(frontEnd) {
			frontBegin, frontEnd = frontEnd, frontBegin
		}
		steps++
		nextFront := make(map[string]bool)
		for word := range frontBegin {
			bs := []byte(word)
			for pos := 0; pos < len(bs); pos++ {
				orig := bs[pos]
				for c := byte('a'); c <= byte('z'); c++ {
					if c == orig {
						continue
					}
					bs[pos] = c
					next := string(bs)
					if frontEnd[next] {
						return steps // frontiers meet
					}
					if wordSet[next] {
						wordSet[next] = false // mark visited by deleting
						nextFront[next] = true
					}
					bs[pos] = orig
				}
			}
		}
		frontBegin = nextFront
	}
	return 0
}
```

### Dry Run
`hit → cog`:
- frontBegin={hit}, frontEnd={cog}, steps=1.
- Expand begin (smaller): next={hot}. steps=2.
- Expand end: next={dog,log}. steps=3.
- Expand begin {hot}: next={dot,lot}. steps=4.
- Expand end {dog,log}: try dot→ in frontBegin? No. Try lot→ in frontBegin? Yes! return 5.

---

## Key Takeaways
- BFS on an implicit graph: don't build the graph explicitly, generate neighbors on the fly.
- Delete from wordSet to mark visited (faster than a separate visited set).
- Bidirectional BFS halves the search depth → much faster for large word lists.

---

## Related Problems
- LeetCode #126 — Word Ladder II (all shortest paths)
- LeetCode #433 — Minimum Genetic Mutation (same pattern)
