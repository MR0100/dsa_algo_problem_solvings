package main

import "fmt"

// ── Approach 1: Backtracking with One Map (Brute Force) ───────────────────────
//
// backtrackOneMap solves Word Pattern II by trying every way to cut the string
// into pieces matching the pattern, mapping each pattern letter to a substring,
// and scanning the map's values to reject two letters mapping to the same word.
//
// Intuition:
//
//	We need a BIJECTION between pattern characters and non-empty substrings such
//	that concatenating the substring for each pattern char (in order) rebuilds
//	s. Since substring lengths are unknown, we search: for the current pattern
//	char, if it's already mapped, the next slice of s must equal that word;
//	otherwise try every possible prefix of the remaining string as its word.
//	A single char→word map enforces "same letter → same word"; to also enforce
//	"different letters → different words" we linearly check no other key already
//	holds this candidate word.
//
// Algorithm:
//
//	dfs(pi, si):
//	  if pi == len(pattern) and si == len(s): success.
//	  c = pattern[pi].
//	  if c mapped to w: s must start with w at si; recurse (pi+1, si+len(w)).
//	  else: for each end from si+1..len(s): w = s[si:end];
//	        skip if w already used by another char; map c→w, recurse; unmap.
//
// Time:  O(n^m) worst case — m pattern chars, each tries O(n) split lengths,
//
//	with an O(m) value scan per assignment.
//
// Space: O(m) map + recursion depth O(m).
func backtrackOneMap(pattern string, s string) bool {
	charToWord := map[byte]string{} // pattern letter → assigned substring
	var dfs func(pi, si int) bool
	dfs = func(pi, si int) bool {
		if pi == len(pattern) { // consumed the whole pattern...
			return si == len(s) // ...succeed only if we also consumed all of s
		}
		c := pattern[pi]
		if w, ok := charToWord[c]; ok {
			// c is already bound: the next slice of s MUST equal its word.
			end := si + len(w)
			if end > len(s) || s[si:end] != w {
				return false // mismatch → this branch is dead
			}
			return dfs(pi+1, end) // consume w and move on
		}
		// c is unbound: try every non-empty prefix of the remaining string.
		for end := si + 1; end <= len(s); end++ {
			w := s[si:end]
			// Enforce injectivity: no other pattern char may already own this word.
			used := false
			for _, existing := range charToWord {
				if existing == w {
					used = true
					break
				}
			}
			if used {
				continue // would break the bijection
			}
			charToWord[c] = w // tentatively assign
			if dfs(pi+1, end) {
				return true
			}
			delete(charToWord, c) // backtrack: undo the assignment
		}
		return false
	}
	return dfs(0, 0)
}

// ── Approach 2: Backtracking with Two Maps (Optimal) ──────────────────────────
//
// backtrackTwoMaps solves Word Pattern II with the same search but replaces the
// O(m) "is this word already used" scan with a second word→char map, making the
// injectivity check O(1).
//
// Intuition:
//
//	A bijection needs BOTH directions enforced: char→word AND word→char. Keeping
//	the inverse map lets us reject in O(1) both "letter already maps elsewhere"
//	and "word already claimed by another letter", instead of scanning all values.
//
// Algorithm:
//
//	dfs(pi, si):
//	  if pi == len(pattern): return si == len(s).
//	  c = pattern[pi].
//	  if c mapped to w: verify s[si:si+len(w)] == w; recurse.
//	  else: for each candidate word w = s[si:end]:
//	        if w already claimed (in wordToChar): skip.
//	        set both maps, recurse, then unset both on backtrack.
//
// Time:  O(n^m) split search, but O(1) per bijection check.
// Space: O(m) for the two maps + O(m) recursion.
func backtrackTwoMaps(pattern string, s string) bool {
	charToWord := map[byte]string{} // forward map: letter → word
	wordToChar := map[string]byte{} // inverse map: word → letter
	var dfs func(pi, si int) bool
	dfs = func(pi, si int) bool {
		if pi == len(pattern) {
			return si == len(s) // both fully consumed ⇒ valid matching
		}
		c := pattern[pi]
		if w, ok := charToWord[c]; ok {
			// Bound letter: the upcoming slice must equal its word exactly.
			end := si + len(w)
			if end > len(s) || s[si:end] != w {
				return false
			}
			return dfs(pi+1, end)
		}
		for end := si + 1; end <= len(s); end++ {
			w := s[si:end]
			if _, taken := wordToChar[w]; taken {
				continue // another letter already owns this word → injectivity broken
			}
			charToWord[c] = w // bind both directions
			wordToChar[w] = c
			if dfs(pi+1, end) {
				return true
			}
			delete(charToWord, c) // backtrack both maps together
			delete(wordToChar, w)
		}
		return false
	}
	return dfs(0, 0)
}

func main() {
	fmt.Println("=== Approach 1: Backtracking (One Map) ===")
	fmt.Println(backtrackOneMap("abab", "redblueredblue")) // expected true
	fmt.Println(backtrackOneMap("aaaa", "asdasdasdasd"))   // expected true
	fmt.Println(backtrackOneMap("aabb", "xyzabcxzyabc"))   // expected false

	fmt.Println("=== Approach 2: Backtracking (Two Maps, Optimal) ===")
	fmt.Println(backtrackTwoMaps("abab", "redblueredblue")) // expected true
	fmt.Println(backtrackTwoMaps("aaaa", "asdasdasdasd"))   // expected true
	fmt.Println(backtrackTwoMaps("aabb", "xyzabcxzyabc"))   // expected false
}
