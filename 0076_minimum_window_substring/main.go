package main

import "fmt"

// ── Approach 1: Brute Force ───────────────────────────────────────────────────
//
// bruteForce solves Minimum Window Substring by checking all substrings.
//
// Intuition:
//   Try every substring s[i:j]. For each, check if it contains all characters
//   of t. Track the shortest valid one.
//
// Time:  O(n² × |t|) — n² substrings, each checked in O(|t|).
// Space: O(|Σ|) — frequency map.
func bruteForce(s string, t string) string {
	if len(s) == 0 || len(t) == 0 {
		return ""
	}
	need := make(map[byte]int)
	for i := 0; i < len(t); i++ {
		need[t[i]]++
	}
	best := ""
	for i := 0; i < len(s); i++ {
		have := make(map[byte]int)
		for j := i; j < len(s); j++ {
			have[s[j]]++
			// check if all chars of t are covered
			valid := true
			for ch, cnt := range need {
				if have[ch] < cnt {
					valid = false
					break
				}
			}
			if valid {
				if best == "" || j-i+1 < len(best) {
					best = s[i : j+1]
				}
				break // shorter extension won't help
			}
		}
	}
	return best
}

// ── Approach 2: Sliding Window (Optimal) ─────────────────────────────────────
//
// slidingWindow solves Minimum Window Substring using a two-pointer sliding window.
//
// Intuition:
//   Expand the right pointer to include characters. Once all characters of t
//   are covered (formed == required), try shrinking from the left to find the
//   minimum window. When shrinking makes the window invalid, expand right again.
//
//   Track `formed` (count of unique chars in t that are satisfied) vs
//   `required` (total unique chars in t).
//
// Algorithm:
//   need = freq map of t; required = len(unique chars in t)
//   l=0, formed=0; windowCnt={}
//   for r in 0..n-1:
//     add s[r] to windowCnt
//     if windowCnt[s[r]] == need[s[r]]: formed++
//     while formed == required:
//       update best window
//       remove s[l] from windowCnt; if drops below need: formed--; l++
//
// Time:  O(|s| + |t|) — each char added/removed from window at most once.
// Space: O(|Σ|) — frequency maps.
func slidingWindow(s string, t string) string {
	if len(s) == 0 || len(t) == 0 {
		return ""
	}
	need := make(map[byte]int)
	for i := 0; i < len(t); i++ {
		need[t[i]]++
	}
	required := len(need) // number of unique chars in t

	windowCnt := make(map[byte]int)
	formed := 0 // unique chars in window that meet their required frequency

	l := 0
	bestLen := -1
	bestL, bestR := 0, 0

	for r := 0; r < len(s); r++ {
		ch := s[r]
		windowCnt[ch]++
		if need[ch] > 0 && windowCnt[ch] == need[ch] {
			formed++ // this char's requirement is now satisfied
		}

		// try to shrink from left
		for formed == required {
			if bestLen == -1 || r-l+1 < bestLen {
				bestLen = r - l + 1
				bestL, bestR = l, r
			}
			lch := s[l]
			windowCnt[lch]--
			if need[lch] > 0 && windowCnt[lch] < need[lch] {
				formed-- // this char's requirement is no longer satisfied
			}
			l++
		}
	}

	if bestLen == -1 {
		return ""
	}
	return s[bestL : bestR+1]
}

func main() {
	fmt.Println("=== Approach 1: Brute Force ===")
	fmt.Printf("s=%q t=%q  got=%q  expected %q\n", "ADOBECODEBANC", "ABC", bruteForce("ADOBECODEBANC", "ABC"), "BANC")
	fmt.Printf("s=%q t=%q  got=%q  expected %q\n", "a", "a", bruteForce("a", "a"), "a")
	fmt.Printf("s=%q t=%q  got=%q  expected %q\n", "a", "aa", bruteForce("a", "aa"), "")

	fmt.Println("=== Approach 2: Sliding Window ===")
	fmt.Printf("s=%q t=%q  got=%q  expected %q\n", "ADOBECODEBANC", "ABC", slidingWindow("ADOBECODEBANC", "ABC"), "BANC")
	fmt.Printf("s=%q t=%q  got=%q  expected %q\n", "a", "a", slidingWindow("a", "a"), "a")
	fmt.Printf("s=%q t=%q  got=%q  expected %q\n", "a", "aa", slidingWindow("a", "aa"), "")
	fmt.Printf("s=%q t=%q  got=%q  expected %q\n", "aa", "aa", slidingWindow("aa", "aa"), "aa")
}
