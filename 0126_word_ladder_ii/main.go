package main

import "fmt"

// ── Approach 1: BFS + Backtracking ────────────────────────────────────────────
//
// findLadders solves Word Ladder II using BFS to find shortest path length
// and then DFS/backtracking to collect all shortest paths.
//
// Intuition:
//   BFS from beginWord builds a layer-by-layer graph. Record for each word
//   which words in the *previous* BFS layer are its parents. Stop BFS when
//   endWord is reached. Then DFS from endWord back to beginWord using the
//   parent map to reconstruct all shortest paths.
//
// Algorithm:
//   1. BFS: at each step expand all words reachable from the current layer.
//      For each neighbor, if not yet visited, add to next layer and record parent.
//      Stop when endWord appears in the next layer.
//   2. DFS: backtrack from endWord → beginWord using the parent map.
//      Reverse each path before appending to result.
//
// Time:  O(N * 26 * L + paths) where N=word list size, L=word length.
// Space: O(N * L) — BFS layers + parent map.
func findLadders(beginWord string, endWord string, wordList []string) [][]string {
	wordSet := make(map[string]bool)
	for _, w := range wordList {
		wordSet[w] = true
	}
	if !wordSet[endWord] {
		return nil
	}

	// parents[word] = set of words that can transform into word at this BFS level
	parents := make(map[string][]string)
	currentLayer := map[string]bool{beginWord: true}
	visited := map[string]bool{beginWord: true}
	found := false

	for len(currentLayer) > 0 && !found {
		nextLayer := make(map[string]bool)
		// mark current layer as visited before processing (avoid going back)
		for w := range currentLayer {
			visited[w] = true
		}

		for word := range currentLayer {
			bs := []byte(word)
			for i := 0; i < len(bs); i++ {
				orig := bs[i]
				for c := byte('a'); c <= byte('z'); c++ {
					if c == orig {
						continue
					}
					bs[i] = c
					next := string(bs)
					if wordSet[next] && !visited[next] {
						nextLayer[next] = true
						parents[next] = append(parents[next], word)
					}
					bs[i] = orig
				}
			}
		}

		if nextLayer[endWord] {
			found = true
		}
		currentLayer = nextLayer
	}

	if !found {
		return nil
	}

	// DFS backtrack from endWord to beginWord
	var result [][]string
	path := []string{endWord}

	var dfs func(word string)
	dfs = func(word string) {
		if word == beginWord {
			cp := make([]string, len(path))
			copy(cp, path)
			// path was built end→begin, reverse it
			for i, j := 0, len(cp)-1; i < j; i, j = i+1, j-1 {
				cp[i], cp[j] = cp[j], cp[i]
			}
			result = append(result, cp)
			return
		}
		for _, parent := range parents[word] {
			path = append(path, parent)
			dfs(parent)
			path = path[:len(path)-1]
		}
	}
	dfs(endWord)
	return result
}

func main() {
	fmt.Println("=== Approach 1: BFS + Backtracking ===")
	r1 := findLadders("hit", "cog", []string{"hot", "dot", "dog", "lot", "log", "cog"})
	fmt.Printf("beginWord=hit endWord=cog  got=%v\n  expected [[hit hot dot dog cog] [hit hot lot log cog]]\n", r1)

	r2 := findLadders("hit", "cog", []string{"hot", "dot", "dog", "lot", "log"})
	fmt.Printf("beginWord=hit endWord=cog (no cog in list)  got=%v  expected []\n", r2)
}
