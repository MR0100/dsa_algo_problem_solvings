package main

import "fmt"

// ── Approach 1: Brute Force ───────────────────────────────────────────────────
//
// bruteForce solves Substring with Concatenation of All Words by checking
// every possible starting position of length total = len(words) * wordLen.
//
// Intuition: A valid window is exactly len(words)*wordLen characters long.
// Slide a window of that size across s; for each window, extract the words
// and check that they match the multiset of the given words array.
//
// Algorithm:
//  1. Build frequency map need from words.
//  2. For each start i in [0, len(s)-total]:
//     extract words from s[i..i+total-1] in chunks of wordLen.
//     build a "have" map; if any word is unknown or over-count, break.
//     if have == need: record i.
//  3. Return all valid starts.
//
// Time:  O(n * k * wordLen) where n=len(s), k=len(words)
// Space: O(k) — frequency maps
func bruteForce(s string, words []string) []int {
	if len(s) == 0 || len(words) == 0 {
		return nil
	}
	wordLen := len(words[0])
	k := len(words)
	total := wordLen * k

	// build required frequency map
	need := make(map[string]int)
	for _, w := range words {
		need[w]++
	}

	var result []int
	for i := 0; i <= len(s)-total; i++ {
		have := make(map[string]int)
		j := 0
		for j < k {
			word := s[i+j*wordLen : i+j*wordLen+wordLen] // extract j-th word in window
			if need[word] == 0 {                           // word not in words list
				break
			}
			have[word]++
			if have[word] > need[word] { // word appears too many times
				break
			}
			j++
		}
		if j == k { // all k words matched
			result = append(result, i)
		}
	}
	return result
}

// ── Approach 2: Sliding Window with Frequency Map (Optimal) ──────────────────
//
// slidingWindow solves Substring with Concatenation of All Words using
// k sliding windows (one per starting offset mod wordLen).
//
// Intuition: Words are fixed-length. There are wordLen distinct "grids" to
// consider — starting at offset 0, 1, …, wordLen-1. For each grid, use a
// sliding window at word granularity: maintain a "have" map and a count of
// matched words. When a word is unknown, restart. When a word is over-count,
// shrink from the left until the count is correct. When matched == k, record.
//
// Algorithm:
//  1. For offset = 0 to wordLen-1:
//     left = offset, matched = 0, have = {}.
//     i steps from offset to len(s)-wordLen in steps of wordLen:
//       word = s[i..i+wordLen-1].
//       if word not in need: reset left=i+wordLen, have={}, matched=0; continue.
//       have[word]++; if have[word] <= need[word]: matched++.
//       while have[word] > need[word]: shrink from left.
//       if matched == k: record left; shrink left by one word.
//
// Time:  O(n) — each character examined twice per offset pass; total O(wordLen * n/wordLen) = O(n)
// Space: O(k) — have map
func slidingWindow(s string, words []string) []int {
	if len(s) == 0 || len(words) == 0 {
		return nil
	}
	wordLen := len(words[0])
	k := len(words)
	total := wordLen * k
	n := len(s)

	if n < total {
		return nil
	}

	need := make(map[string]int)
	for _, w := range words {
		need[w]++
	}

	var result []int

	// run a separate sliding window for each starting offset
	for offset := 0; offset < wordLen; offset++ {
		have := make(map[string]int)
		matched := 0  // number of words currently matched with correct count
		left := offset // left boundary of the current window (word-aligned)

		for i := offset; i <= n-wordLen; i += wordLen {
			word := s[i : i+wordLen]

			if _, ok := need[word]; !ok {
				// unknown word: reset window entirely
				have = make(map[string]int)
				matched = 0
				left = i + wordLen
				continue
			}

			have[word]++
			if have[word] <= need[word] {
				matched++ // this copy contributes to a valid window
			}

			// shrink from left if word is over-count
			for have[word] > need[word] {
				leftWord := s[left : left+wordLen]
				have[leftWord]--
				if have[leftWord] < need[leftWord] {
					matched-- // we lost a contributing copy
				}
				left += wordLen
			}

			// window contains exactly k matched words
			if matched == k {
				result = append(result, left)
				// slide left forward by one word for next iteration
				leftWord := s[left : left+wordLen]
				have[leftWord]--
				if have[leftWord] < need[leftWord] {
					matched--
				}
				left += wordLen
			}
		}
	}
	return result
}

func main() {
	fmt.Println("=== Approach 1: Brute Force ===")
	fmt.Printf("s=barfoothefoobarman words=[foo,bar]  => %v  expected [0 9]\n",
		bruteForce("barfoothefoobarman", []string{"foo", "bar"}))
	fmt.Printf("s=wordgoodgoodgoodbestword words=[word,good,best,word]  => %v  expected []\n",
		bruteForce("wordgoodgoodgoodbestword", []string{"word", "good", "best", "word"}))
	fmt.Printf("s=barfoofoobarthefoobarman words=[bar,foo,the]  => %v  expected [6 9 12]\n",
		bruteForce("barfoofoobarthefoobarman", []string{"bar", "foo", "the"}))

	fmt.Println("\n=== Approach 2: Sliding Window (Optimal) ===")
	fmt.Printf("s=barfoothefoobarman words=[foo,bar]  => %v  expected [0 9]\n",
		slidingWindow("barfoothefoobarman", []string{"foo", "bar"}))
	fmt.Printf("s=wordgoodgoodgoodbestword words=[word,good,best,word]  => %v  expected []\n",
		slidingWindow("wordgoodgoodgoodbestword", []string{"word", "good", "best", "word"}))
	fmt.Printf("s=barfoofoobarthefoobarman words=[bar,foo,the]  => %v  expected [6 9 12]\n",
		slidingWindow("barfoofoobarthefoobarman", []string{"bar", "foo", "the"}))
}
