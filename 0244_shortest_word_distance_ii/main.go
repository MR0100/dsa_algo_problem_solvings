package main

import "fmt"

// Shortest Word Distance II — DESIGN problem.
//
// Build a data structure once from wordsDict, then answer many shortest(word1,
// word2) queries efficiently. The key idea is to precompute, for every word,
// the sorted list of indices where it appears, so each query merges two sorted
// index lists instead of rescanning the whole array.

// ── Approach 1: Precompute Index Map + Merge Two Sorted Lists (Optimal) ───────
//
// WordDistance is the design type. Its constructor records every word's sorted
// occurrence indices; each shortest() query walks the two relevant lists with
// a two-pointer merge.
//
// Intuition:
//
//	Occurrences of each word appear in increasing index order as we scan left
//	to right, so wordIndex[word] is already sorted. To find the closest pair
//	between two sorted lists, advance whichever pointer points to the smaller
//	index — that is the only move that can shrink the gap — and track the
//	minimum |a - b| along the way.
//
// Algorithm (Constructor):
//  1. For each position i, append i to wordIndex[words[i]].
//
// Algorithm (shortest):
//  1. Let l1 = wordIndex[word1], l2 = wordIndex[word2], p1 = p2 = 0.
//  2. While both pointers in range: update min with |l1[p1] - l2[p2]|, then
//     advance the pointer at the smaller value.
//  3. Return min.
//
// Time:  Constructor O(n); each query O(a + b) where a,b are the two words'
//
//	occurrence counts (≤ O(n)).
//
// Space: O(n) — every index stored once across the map.
type WordDistance struct {
	wordIndex map[string][]int // word → sorted list of positions
}

// Constructor builds the index map from wordsDict in one pass.
func Constructor(wordsDict []string) WordDistance {
	idx := make(map[string][]int, len(wordsDict))
	for i, w := range wordsDict {
		idx[w] = append(idx[w], i) // positions accumulate in increasing order
	}
	return WordDistance{wordIndex: idx}
}

// shortest returns the minimum distance between word1 and word2 by merging
// their two sorted index lists.
func (wd *WordDistance) shortest(word1 string, word2 string) int {
	l1 := wordIndex(wd, word1) // sorted indices of word1
	l2 := wordIndex(wd, word2) // sorted indices of word2

	p1, p2 := 0, 0
	const maxInt = int(^uint(0) >> 1) // safe large upper bound for the running min
	min := maxInt
	// Merge: compare current heads, keep the gap, move the smaller forward.
	for p1 < len(l1) && p2 < len(l2) {
		a, b := l1[p1], l2[p2]
		if d := abs(a - b); d < min {
			min = d // closer pair
		}
		if a < b {
			p1++ // advancing the smaller index is the only way to reduce the gap
		} else {
			p2++
		}
	}
	return min
}

// ── Approach 2: Precompute Index Map + Linear Rescan Per Query (Baseline) ─────
//
// shortestRescan answers a query by re-running the one-pass two-pointer scan of
// problem 243 over the original words array. Kept as a baseline that shows why
// the merge version wins when queries are frequent.
//
// Intuition:
//
//	Without stored index lists, each query must sweep the whole array tracking
//	the last-seen index of each word (the LeetCode 243 technique).
//
// Algorithm:
//  1. i1 = i2 = -1.
//  2. Scan words: on word1 set i1 and try i1-i2; on word2 set i2 and try i2-i1.
//  3. Return the minimum.
//
// Time:  O(n) per query regardless of how rare the words are.
// Space: O(n) to keep the original array around.
type WordDistanceRescan struct {
	words []string
}

func NewRescan(wordsDict []string) WordDistanceRescan {
	return WordDistanceRescan{words: wordsDict}
}

func (wd *WordDistanceRescan) shortest(word1 string, word2 string) int {
	i1, i2 := -1, -1
	min := len(wd.words)
	for k, w := range wd.words {
		switch w {
		case word1:
			i1 = k
			if i2 != -1 && i1-i2 < min {
				min = i1 - i2
			}
		case word2:
			i2 = k
			if i1 != -1 && i2-i1 < min {
				min = i2 - i1
			}
		}
	}
	return min
}

// wordIndex fetches a word's stored index list.
func wordIndex(wd *WordDistance, w string) []int { return wd.wordIndex[w] }

// abs returns the absolute value of x.
func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func main() {
	words := []string{"practice", "makes", "perfect", "coding", "makes"}

	fmt.Println("=== Approach 1: Index Map + Merge (Optimal) ===")
	wd := Constructor(words)
	fmt.Println(wd.shortest("coding", "practice")) // expected 3
	fmt.Println(wd.shortest("makes", "coding"))    // expected 1

	fmt.Println("=== Approach 2: Index Map placeholder + Rescan Per Query ===")
	wr := NewRescan(words)
	fmt.Println(wr.shortest("coding", "practice")) // expected 3
	fmt.Println(wr.shortest("makes", "coding"))    // expected 1
}
