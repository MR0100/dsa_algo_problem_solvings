package main

import "fmt"

// ── Approach 1: Brute Force (All Pairs) ──────────────────────────────────────
//
// bruteForce solves Shortest Word Distance by comparing every occurrence of
// word1 with every occurrence of word2.
//
// Intuition:
//
//	The distance is |i - j| for some index i holding word1 and some index j
//	holding word2. Try all such pairs and keep the smallest gap.
//
// Algorithm:
//  1. min = large sentinel.
//  2. For each i where words[i] == word1:
//     For each j where words[j] == word2: min = min(min, |i - j|).
//  3. Return min.
//
// Time:  O(n²) worst case — nested scan over all matching pairs.
// Space: O(1) — only a running minimum.
func bruteForce(words []string, word1 string, word2 string) int {
	min := len(words) // any real distance is < len(words); safe upper bound
	for i := 0; i < len(words); i++ {
		if words[i] != word1 {
			continue // outer loop only stops on word1 occurrences
		}
		for j := 0; j < len(words); j++ {
			if words[j] == word2 { // inner loop only on word2 occurrences
				if d := abs(i - j); d < min {
					min = d // tighter pair found
				}
			}
		}
	}
	return min
}

// ── Approach 2: One-Pass Two Pointers (Optimal) ──────────────────────────────
//
// twoPointers solves Shortest Word Distance in a single pass by tracking the
// most recent index of each word.
//
// Intuition:
//
//	The closest occurrence of word2 to a given word1 is the most recent one
//	seen. So scan left to right; whenever we hit either word, update its last
//	position, and if we already have a position for the other word, the gap
//	between them is a candidate answer. Because we always use the latest index
//	of each, no closer pair is ever missed.
//
// Algorithm:
//  1. i1 = i2 = -1 (last seen indices).
//  2. For k, w in words:
//     - if w == word1: i1 = k; if i2 != -1, update min with i1 - i2.
//     - if w == word2: i2 = k; if i1 != -1, update min with i2 - i1.
//  3. Return min.
//
// Time:  O(n) — one linear scan.
// Space: O(1) — two index variables and a minimum.
func twoPointers(words []string, word1 string, word2 string) int {
	i1, i2 := -1, -1  // last positions of word1 and word2 (-1 = not yet seen)
	min := len(words) // safe upper bound on any distance
	for k, w := range words {
		switch w {
		case word1:
			i1 = k // record newest word1 position
			if i2 != -1 && i1-i2 < min {
				min = i1 - i2 // pair with the latest word2
			}
		case word2:
			i2 = k // record newest word2 position
			if i1 != -1 && i2-i1 < min {
				min = i2 - i1 // pair with the latest word1
			}
		}
	}
	return min
}

// abs returns the absolute value of x.
func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func main() {
	words := []string{"practice", "makes", "perfect", "coding", "makes"}

	fmt.Println("=== Approach 1: Brute Force ===")
	fmt.Println(bruteForce(words, "coding", "practice")) // expected 3
	fmt.Println(bruteForce(words, "makes", "coding"))    // expected 1

	fmt.Println("=== Approach 2: One-Pass Two Pointers (Optimal) ===")
	fmt.Println(twoPointers(words, "coding", "practice")) // expected 3
	fmt.Println(twoPointers(words, "makes", "coding"))    // expected 1
}
