package main

import "fmt"

// isVowel reports whether b is a vowel (case-insensitive). Vowels are a,e,i,o,u.
func isVowel(b byte) bool {
	switch b {
	case 'a', 'e', 'i', 'o', 'u', 'A', 'E', 'I', 'O', 'U':
		return true
	}
	return false
}

// ── Approach 1: Two Pointers (Optimal) ───────────────────────────────────────
//
// twoPointers reverses only the vowels of s in place, leaving consonants fixed.
//
// Intuition:
//
//	Reversing "only the vowels" means the sequence of vowel characters is
//	reversed while every non-vowel stays at its position. Use one pointer from
//	each end: advance each until it lands on a vowel, then swap those two
//	vowels and move both inward. Consonants are simply skipped.
//
// Algorithm:
//  1. Convert to a byte slice; left=0, right=len-1.
//  2. Advance left past non-vowels; advance right (leftward) past non-vowels.
//  3. When both point at vowels and left < right, swap them; step both inward.
//  4. Repeat until left >= right.
//
// Time:  O(n) — each index visited once. Space: O(n) for the byte slice (O(1)
//
//	extra beyond the required mutable buffer).
func twoPointers(s string) string {
	b := []byte(s) // strings are immutable in Go; work on a byte slice
	left, right := 0, len(b)-1
	for left < right {
		for left < right && !isVowel(b[left]) { // skip non-vowels from the left
			left++
		}
		for left < right && !isVowel(b[right]) { // skip non-vowels from the right
			right--
		}
		if left < right {
			b[left], b[right] = b[right], b[left] // swap the two vowels
			left++
			right--
		}
	}
	return string(b)
}

// ── Approach 2: Collect Indices then Reverse (Two-Pass) ──────────────────────
//
// collectIndices records all vowel positions, then rewrites them with the vowel
// characters in reverse order.
//
// Intuition:
//
//	First find WHERE the vowels are (their indices, left to right). The vowels
//	must end up in reverse order at those same slots, so pair the k-th index
//	from the front with the k-th vowel from the back.
//
// Algorithm:
//  1. Pass 1: collect indices idx[] where b[i] is a vowel.
//  2. Pass 2: for k in 0..len(idx)-1, place b[idx[len-1-k]]'s original vowel at
//     idx[k]. Concretely, swap using two pointers over idx.
//
// Time:  O(n) two passes. Space: O(V) for the index list (V = vowel count).
func collectIndices(s string) string {
	b := []byte(s)
	idx := []int{}
	for i := 0; i < len(b); i++ { // pass 1: where are the vowels?
		if isVowel(b[i]) {
			idx = append(idx, i)
		}
	}
	// pass 2: swap symmetric pairs of vowel positions.
	for i, j := 0, len(idx)-1; i < j; i, j = i+1, j-1 {
		b[idx[i]], b[idx[j]] = b[idx[j]], b[idx[i]]
	}
	return string(b)
}

func main() {
	fmt.Println("=== Approach 1: Two Pointers (Optimal) ===")
	fmt.Println(twoPointers("IceCreAm")) // expected AceCreIm
	fmt.Println(twoPointers("leetcode")) // expected leotcede

	fmt.Println("=== Approach 2: Collect Indices then Reverse ===")
	fmt.Println(collectIndices("IceCreAm")) // expected AceCreIm
	fmt.Println(collectIndices("leetcode")) // expected leotcede
}
