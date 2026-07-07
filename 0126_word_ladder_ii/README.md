# 0126 — Word Ladder II

> LeetCode #126 · Difficulty: Hard
> **Categories:** Hash Table, String, Backtracking, Breadth-First Search

---

## Problem Statement

A **transformation sequence** from word `beginWord` to word `endWord` using a dictionary `wordList` is a sequence of words `beginWord -> s1 -> s2 -> ... -> sk` such that:
- Every adjacent pair of words differs by a single letter.
- Every `si` for `1 <= i <= k` is in `wordList`. Note that `beginWord` does not need to be in `wordList`.
- `sk == endWord`

Given `beginWord`, `endWord`, and `wordList`, return all the **shortest transformation sequences** from `beginWord` to `endWord`, or an empty list if no such sequence exists.

**Example 1:**
```
Input: beginWord = "hit", endWord = "cog",
       wordList = ["hot","dot","dog","lot","log","cog"]
Output: [["hit","hot","dot","dog","cog"],["hit","hot","lot","log","cog"]]
```

**Example 2:**
```
Input: beginWord = "hit", endWord = "cog",
       wordList = ["hot","dot","dog","lot","log"]
Output: []
```

**Constraints:**
- `1 <= beginWord.length <= 5`
- `endWord.length == beginWord.length`
- `1 <= wordList.length <= 500`
- All words have the same length and consist of lowercase English letters.

---

## Company Frequency

| Company   | Frequency    | Last Reported |
|-----------|--------------|---------------|
| Amazon    | ★★★★☆ High  | 2024          |
| Google    | ★★★★☆ High  | 2024          |
| Microsoft | ★★★☆☆ Medium | 2023          |

> ⚠️ Frequency data is crowd-sourced from LeetCode Discuss, Glassdoor, and community interview reports. Treat as a signal, not a guarantee.

---

## DSA Concepts Used

- **BFS** — find shortest path length and build parent map layer by layer
- **Backtracking** — DFS from endWord to beginWord using parent map

---

## Approaches Overview

| # | Approach              | Time                 | Space   | When to use |
|---|-----------------------|----------------------|---------|-------------|
| 1 | BFS + DFS Backtrack   | O(N·26·L + paths)    | O(N·L)  | Standard    |

---

## Approach 1 — BFS + Backtracking

### Intuition
Two-phase approach:

**Phase 1 (BFS):** Build a `parents` map: `parents[word]` = list of words that can reach `word` in BFS layer `k-1` (when `word` is first discovered at layer `k`). Stop BFS as soon as `endWord` is reached.

**Phase 2 (DFS):** Backtrack from `endWord` to `beginWord` using the `parents` map. Reverse each collected path before recording.

Key: we mark entire BFS layers as visited before expanding them, to correctly track shortest paths (multiple parents at the same BFS depth are all valid).

### Algorithm
1. Build `wordSet`.
2. BFS with `currentLayer` set. For each word in current layer, try all 26-letter substitutions. If neighbor in wordSet and not yet visited: add to nextLayer, record parent.
3. Stop when endWord appears in nextLayer.
4. DFS backtrack: path starts at [endWord], recurse over parents.

### Complexity
- **Time:** O(N · 26 · L) for BFS + O(paths · L) for backtracking.
- **Space:** O(N · L) — parents map.

### Code
See `main.go` — `findLadders`.

### Dry Run
`hit → cog`, wordList = [hot,dot,dog,lot,log,cog]:

BFS layers:
- Layer 1: {hit}
- Layer 2: {hot}. parents[hot]=[hit].
- Layer 3: {dot,lot}. parents[dot]=[hot], parents[lot]=[hot].
- Layer 4: {dog,log}. parents[dog]=[dot], parents[log]=[lot].
- Layer 5: {cog}. parents[cog]=[dog,log]. Found!

DFS from cog:
- cog←dog←dot←hot←hit → [hit,hot,dot,dog,cog]
- cog←log←lot←hot←hit → [hit,hot,lot,log,cog]

---

## Key Takeaways
- Mark entire BFS layer as visited *before* expanding it — otherwise the same word can appear at two positions in the same layer, losing valid parents.
- Phase separation: BFS for distance, DFS for path reconstruction.
- `parents` map (word → list of predecessors) is the bridge between phases.

---

## Related Problems
- LeetCode #127 — Word Ladder (only shortest length, no path enumeration)
