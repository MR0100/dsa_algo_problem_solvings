package main

import "fmt"

// ── Approach 1: BFS ───────────────────────────────────────────────────────────
//
// ladderLength solves Word Ladder using BFS.
//
// Intuition:
//   BFS from beginWord, treating word transformations as edges in a graph.
//   At each step, try all 26-letter substitutions at every position. If a
//   resulting word is in the wordSet and not yet visited, enqueue it.
//   Return the level count when endWord is found.
//
// Time:  O(N * 26 * L) — N words, L letters each.
// Space: O(N * L)
func ladderLength(beginWord string, endWord string, wordList []string) int {
	wordSet := make(map[string]bool)
	for _, w := range wordList {
		wordSet[w] = true
	}
	if !wordSet[endWord] {
		return 0
	}

	queue := []string{beginWord}
	visited := map[string]bool{beginWord: true}
	steps := 1

	for len(queue) > 0 {
		levelSize := len(queue)
		steps++
		for i := 0; i < levelSize; i++ {
			word := queue[0]
			queue = queue[1:]
			bs := []byte(word)
			for pos := 0; pos < len(bs); pos++ {
				orig := bs[pos]
				for c := byte('a'); c <= byte('z'); c++ {
					if c == orig {
						continue
					}
					bs[pos] = c
					next := string(bs)
					if next == endWord {
						return steps
					}
					if wordSet[next] && !visited[next] {
						visited[next] = true
						queue = append(queue, next)
					}
					bs[pos] = orig
				}
			}
		}
	}
	return 0
}

// ── Approach 2: Bidirectional BFS ─────────────────────────────────────────────
//
// ladderLengthBiBFS solves Word Ladder using bidirectional BFS.
//
// Intuition:
//   Expand BFS from both beginWord and endWord simultaneously.
//   Always expand the smaller frontier. When frontiers meet, the path is found.
//   Reduces search space from O(b^d) to O(b^(d/2)) where b=branching factor.
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

func main() {
	fmt.Println("=== Approach 1: BFS ===")
	fmt.Printf("hit→cog  got=%d  expected 5\n", ladderLength("hit", "cog", []string{"hot", "dot", "dog", "lot", "log", "cog"}))
	fmt.Printf("hit→cog(no cog)  got=%d  expected 0\n", ladderLength("hit", "cog", []string{"hot", "dot", "dog", "lot", "log"}))

	fmt.Println("=== Approach 2: Bidirectional BFS ===")
	fmt.Printf("hit→cog  got=%d  expected 5\n", ladderLengthBiBFS("hit", "cog", []string{"hot", "dot", "dog", "lot", "log", "cog"}))
	fmt.Printf("hit→cog(no cog)  got=%d  expected 0\n", ladderLengthBiBFS("hit", "cog", []string{"hot", "dot", "dog", "lot", "log"}))
}
