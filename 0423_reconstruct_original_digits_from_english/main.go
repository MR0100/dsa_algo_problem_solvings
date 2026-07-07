package main

import (
	"fmt"
	"sort"
	"strings"
)

// digitWords[d] is the English spelling of digit d.
var digitWords = [10]string{"zero", "one", "two", "three", "four", "five", "six", "seven", "eight", "nine"}

// ── Approach 1: Backtracking (Brute Force) ───────────────────────────────────
//
// backtracking solves Reconstruct Original Digits by trying to peel spelled-out
// digit words out of the multiset of letters in every order until the multiset
// is empty, then returning the digits (sorted) of the first full decomposition.
//
// Intuition:
//
//	The string is a shuffle of some digit words concatenated. If we model the
//	string as a letter-count multiset, "peeling" a digit word means subtracting
//	its letters (only allowed when they are all available). A full solution is a
//	sequence of peels that empties the multiset. Depth-first search over which
//	digit to peel next finds one such sequence; because the answer must be
//	returned in ascending digit order, we sort the collected digits at the end.
//
// Algorithm:
//  1. Build count[26] from s.
//  2. DFS: if all counts are zero, success — record current digit list.
//  3. Otherwise try each digit 0..9 whose word can be subtracted; subtract,
//     recurse, and undo (add back) on return.
//  4. On the first success, sort the recorded digits and join into a string.
//
// Time:  Exponential in the worst case — the search tree branches over 10
//
//	digits at each level; fine for correctness/tiny inputs, not for 10^5.
//
// Space: O(depth) recursion = O(number of digits) plus the O(1) counts.
func backtracking(s string) string {
	var count [26]int
	for i := 0; i < len(s); i++ {
		count[s[i]-'a']++ // tally each letter of the shuffled input
	}

	var result []int    // digits of the first successful decomposition
	var chosen []int    // digits peeled on the current DFS path
	var dfs func() bool // returns true once the multiset is emptied

	dfs = func() bool {
		empty := true
		for _, c := range count {
			if c != 0 {
				empty = false // still letters left to consume
				break
			}
		}
		if empty {
			result = append(result, chosen...) // copy the winning path out
			return true
		}
		for d := 0; d <= 9; d++ {
			if canSubtract(&count, digitWords[d]) { // is digit d's word fully available?
				subtract(&count, digitWords[d], -1) // remove its letters
				chosen = append(chosen, d)
				if dfs() { // recurse on the smaller multiset
					return true // stop at the first complete decomposition
				}
				chosen = chosen[:len(chosen)-1]     // undo the choice
				subtract(&count, digitWords[d], +1) // add the letters back
			}
		}
		return false // no digit word fits — dead end on this path
	}

	dfs()
	sort.Ints(result) // answer must be in ascending digit order
	var sb strings.Builder
	for _, d := range result {
		sb.WriteByte(byte('0' + d)) // render each digit as a character
	}
	return sb.String()
}

// canSubtract reports whether every letter of word is currently available.
func canSubtract(count *[26]int, word string) bool {
	var need [26]int
	for i := 0; i < len(word); i++ {
		need[word[i]-'a']++ // how many of each letter this word requires
	}
	for i := 0; i < 26; i++ {
		if count[i] < need[i] {
			return false // not enough of some letter to spell this word
		}
	}
	return true
}

// subtract adds sign*occurrences of each letter of word into count.
func subtract(count *[26]int, word string, sign int) {
	for i := 0; i < len(word); i++ {
		count[word[i]-'a'] += sign // sign = -1 removes the word, +1 restores it
	}
}

// ── Approach 2: Unique-Letter Counting (Optimal) ─────────────────────────────
//
// uniqueLetterCounting solves Reconstruct Original Digits in one pass by
// exploiting that certain digit words own a letter no other digit word has,
// which pins down their counts directly; the rest fall out by subtraction.
//
// Intuition:
//
//	Some spelled digits have a *unique* letter across all ten words:
//	  'z' → only "zero"  (0)      'w' → only "two"   (2)
//	  'u' → only "four"  (4)      'x' → only "six"   (6)
//	  'g' → only "eight" (8)
//	So count['z'] is exactly the number of 0s, etc. Once those five are known,
//	other letters become unique *relative to the remaining* digits:
//	  'o' appears in zero, two, four (all counted) and in "one" → count['o']
//	      minus zeros/twos/fours gives the ones (1).
//	  'h' appears in "three" and "eight" → minus eights gives threes (3).
//	  'f' appears in "four" and "five" → minus fours gives fives (5).
//	  's' appears in "six" and "seven" → minus sixes gives sevens (7).
//	  'i' appears in five, six, eight, nine → minus (5,6,8) gives nines (9).
//	Process the digits in this dependency order; each is determined exactly.
//
// Algorithm:
//  1. Build letter counts of s.
//  2. cnt[0]=z, cnt[2]=w, cnt[4]=u, cnt[6]=x, cnt[8]=g (unique letters).
//  3. cnt[3]=h−cnt[8]; cnt[5]=f−cnt[4]; cnt[7]=s−cnt[6].
//  4. cnt[1]=o−cnt[0]−cnt[2]−cnt[4]; cnt[9]=i−cnt[5]−cnt[6]−cnt[8].
//  5. Emit digit d exactly cnt[d] times, ascending.
//
// Time:  O(n) — one pass to count letters, then O(1) arithmetic and O(n) output.
// Space: O(1) — fixed-size counts (output aside).
func uniqueLetterCounting(s string) string {
	var c [26]int
	for i := 0; i < len(s); i++ {
		c[s[i]-'a']++ // frequency of each letter in the shuffled string
	}
	// helper to read a letter's remaining frequency by its rune
	at := func(ch byte) int { return c[ch-'a'] }

	var cnt [10]int
	// Digits pinned by a letter unique to their word.
	cnt[0] = at('z') // zero: only word containing 'z'
	cnt[2] = at('w') // two:  only word containing 'w'
	cnt[4] = at('u') // four: only word containing 'u'
	cnt[6] = at('x') // six:  only word containing 'x'
	cnt[8] = at('g') // eight:only word containing 'g'

	// Digits whose defining letter is shared only with an already-known digit.
	cnt[3] = at('h') - cnt[8] // 'h' in three & eight
	cnt[5] = at('f') - cnt[4] // 'f' in five & four
	cnt[7] = at('s') - cnt[6] // 's' in seven & six

	// Digits determined after removing the contributions counted above.
	cnt[1] = at('o') - cnt[0] - cnt[2] - cnt[4] // 'o' in one, zero, two, four
	cnt[9] = at('i') - cnt[5] - cnt[6] - cnt[8] // 'i' in nine, five, six, eight

	var sb strings.Builder
	for d := 0; d <= 9; d++ {
		for k := 0; k < cnt[d]; k++ {
			sb.WriteByte(byte('0' + d)) // append digit d, cnt[d] times, in order
		}
	}
	return sb.String()
}

// ── Approach 3: Counting via strconv sanity double-check ─────────────────────
//
// countingVerified is the same O(n) counting idea, but it additionally rebuilds
// the letter multiset from its answer and asserts it matches the input, showing
// how to self-check the reconstruction. It exists to demonstrate verification;
// the returned string is identical to uniqueLetterCounting.
//
// Intuition:
//
//	The counting method is provably exact, but in an interview you can cheaply
//	*prove* your output: re-spell every emitted digit, tally the letters, and
//	confirm the tally equals the original. This catches arithmetic slips.
//
// Algorithm:
//  1. Compute cnt[0..9] exactly as in Approach 2.
//  2. Rebuild the expected letter counts by summing each chosen digit's word.
//  3. If it differs from the input's counts, that would signal a bug (never
//     happens for valid input); otherwise emit the digits.
//
// Time:  O(n) — counting plus a linear reconstruction check.
// Space: O(1) — fixed letter tables.
func countingVerified(s string) string {
	var c [26]int
	for i := 0; i < len(s); i++ {
		c[s[i]-'a']++
	}
	at := func(ch byte) int { return c[ch-'a'] }

	var cnt [10]int
	cnt[0], cnt[2], cnt[4], cnt[6], cnt[8] = at('z'), at('w'), at('u'), at('x'), at('g')
	cnt[3] = at('h') - cnt[8]
	cnt[5] = at('f') - cnt[4]
	cnt[7] = at('s') - cnt[6]
	cnt[1] = at('o') - cnt[0] - cnt[2] - cnt[4]
	cnt[9] = at('i') - cnt[5] - cnt[6] - cnt[8]

	// Self-check: re-spell the digits and confirm the letters add back to s.
	var rebuilt [26]int
	for d := 0; d <= 9; d++ {
		for k := 0; k < cnt[d]; k++ {
			for i := 0; i < len(digitWords[d]); i++ {
				rebuilt[digitWords[d][i]-'a']++
			}
		}
	}
	if rebuilt != c {
		return "" // would indicate a logic error for valid inputs; never triggers
	}

	var sb strings.Builder
	for d := 0; d <= 9; d++ {
		for k := 0; k < cnt[d]; k++ {
			sb.WriteByte(byte('0' + d))
		}
	}
	return sb.String()
}

func main() {
	fmt.Println("=== Approach 1: Backtracking (Brute Force) ===")
	fmt.Println(backtracking("owoztneoer")) // expected 012
	fmt.Println(backtracking("fviefuro"))   // expected 45

	fmt.Println("=== Approach 2: Unique-Letter Counting (Optimal) ===")
	fmt.Println(uniqueLetterCounting("owoztneoer")) // expected 012
	fmt.Println(uniqueLetterCounting("fviefuro"))   // expected 45

	fmt.Println("=== Approach 3: Counting with Self-Check ===")
	fmt.Println(countingVerified("owoztneoer")) // expected 012
	fmt.Println(countingVerified("fviefuro"))   // expected 45
}
