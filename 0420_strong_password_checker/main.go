package main

import "fmt"

// Problem: return the MINIMUM number of steps (insert / delete / replace one
// character) to make `password` strong. Strong means:
//   (a) length in [6, 20],
//   (b) contains a lowercase, an uppercase, and a digit,
//   (c) no three identical characters in a row.

// ── Approach 1: Brute Force (BFS Over Edit States) ───────────────────────────
//
// bruteForceBFS solves Strong Password Checker by breadth-first searching the
// space of strings reachable by single edits, returning the depth (number of
// edits) at which the first strong password appears.
//
// Intuition:
//
//	"Minimum number of single-character edits to satisfy a predicate" is a
//	shortest-path problem: nodes are strings, edges are one insert/delete/
//	replace, and BFS finds the fewest edits to reach ANY strong string. This is
//	exponential and only viable for very short inputs, but it is an unimpeachable
//	oracle to validate the greedy solution against.
//
// Algorithm:
//  1. If s is already strong, return 0.
//  2. BFS level by level. From each string generate all neighbours via one
//     delete, one replace (to a char of each needed class — we allow the full
//     small alphabet here), and one insert.
//  3. The first level containing a strong string is the answer.
//
// Time:  exponential in |s| — only used for tiny strings (length ≤ ~6).
// Space: exponential — the visited set / frontier.
func bruteForceBFS(password string) int {
	if isStrong(password) {
		return 0
	}
	// Small alphabet that still covers all three required classes plus a symbol.
	alphabet := []byte("aB3!")

	seen := map[string]bool{password: true}
	frontier := []string{password}
	steps := 0
	for len(frontier) > 0 {
		steps++
		var next []string
		for _, cur := range frontier {
			// 1) Deletions.
			for i := 0; i < len(cur); i++ {
				cand := cur[:i] + cur[i+1:]
				if !seen[cand] {
					if isStrong(cand) {
						return steps
					}
					seen[cand] = true
					next = append(next, cand)
				}
			}
			// 2) Replacements.
			for i := 0; i < len(cur); i++ {
				for _, ch := range alphabet {
					if cur[i] == ch {
						continue
					}
					cand := cur[:i] + string(ch) + cur[i+1:]
					if !seen[cand] {
						if isStrong(cand) {
							return steps
						}
						seen[cand] = true
						next = append(next, cand)
					}
				}
			}
			// 3) Insertions (positions 0..len).
			for i := 0; i <= len(cur); i++ {
				for _, ch := range alphabet {
					cand := cur[:i] + string(ch) + cur[i:]
					if !seen[cand] {
						if isStrong(cand) {
							return steps
						}
						seen[cand] = true
						next = append(next, cand)
					}
				}
			}
		}
		frontier = next
	}
	return steps // unreachable for valid inputs
}

// isStrong reports whether s already satisfies all three strength rules.
func isStrong(s string) bool {
	if len(s) < 6 || len(s) > 20 {
		return false
	}
	var lower, upper, digit bool
	for i := 0; i < len(s); i++ {
		c := s[i]
		switch {
		case c >= 'a' && c <= 'z':
			lower = true
		case c >= 'A' && c <= 'Z':
			upper = true
		case c >= '0' && c <= '9':
			digit = true
		}
	}
	if !(lower && upper && digit) {
		return false
	}
	// No run of three identical characters.
	for i := 2; i < len(s); i++ {
		if s[i] == s[i-1] && s[i-1] == s[i-2] {
			return false
		}
	}
	return true
}

// ── Approach 2: Greedy (Optimal) ─────────────────────────────────────────────
//
// greedy solves Strong Password Checker in O(n) by splitting on the length
// regime and carefully overlapping the fixes for "missing character types" and
// "three-in-a-row runs".
//
// Intuition:
//
//	Three independent needs interact:
//	  • missing = how many of {lower, upper, digit} are absent (0..3),
//	  • runs of length ≥ 3, each needing ⌊len/3⌋ replacements to break,
//	  • length: too short (<6), fine (6..20), or too long (>20).
//	Replacements and insertions can each simultaneously add a missing type AND
//	break a run, so we always overlap those. The three regimes:
//	  1) len < 6: we only insert (never delete). Each insert can both extend
//	     length and split a run/add a type, so the answer is
//	     max(6 - len, missing).
//	  2) 6 ≤ len ≤ 20: no length change needed. Replacements fix runs, and a
//	     replacement can double as adding a missing type, so the answer is
//	     max(replaceForRuns, missing).
//	  3) len > 20: we MUST delete (len - 20) characters. Deletions can shorten
//	     runs and thereby reduce the replacements they need — and crucially,
//	     deleting to hit run lengths ≡ 0, 1, 2 (mod 3) saves replacements at
//	     different efficiencies. Spend deletions greedily on runs of length
//	     ≡ 0 (mod 3) first (1 deletion saves 1 replacement), then ≡ 1 (mod 3)
//	     (2 deletions save 1), then the rest. Final answer =
//	     deletions + max(replaceRemaining, missing).
//
// Algorithm:
//  1. Compute `missing` from a single scan for the three classes.
//  2. Scan runs of equal characters; for len < 6 or in-range, sum ⌊len/3⌋
//     replacements. For len > 20, bucket each run's length mod 3 so deletions
//     can be applied where they save the most replacements.
//  3. Combine per the regime formulas above.
//
// Time:  O(n) — a couple of linear scans.
// Space: O(1) — a few counters (O(1) buckets for the over-length case).
func greedy(password string) int {
	n := len(password)

	// --- missing character types ---
	var lower, upper, digit bool
	for i := 0; i < n; i++ {
		c := password[i]
		switch {
		case c >= 'a' && c <= 'z':
			lower = true
		case c >= 'A' && c <= 'Z':
			upper = true
		case c >= '0' && c <= '9':
			digit = true
		}
	}
	missing := 0
	if !lower {
		missing++
	}
	if !upper {
		missing++
	}
	if !digit {
		missing++
	}

	// --- collect run lengths of ≥3 repeats ---
	// replace = total replacements needed to break all runs = Σ ⌊len/3⌋.
	// For the over-length regime we also bucket runs by len % 3.
	replace := 0
	// oneMod[r] = number of runs whose length ≡ r (mod 3), among runs of length ≥ 3.
	var buckets [3]int
	i := 0
	for i < n {
		j := i
		for j < n && password[j] == password[i] {
			j++ // extend the run of identical characters
		}
		runLen := j - i
		if runLen >= 3 {
			replace += runLen / 3 // ⌊len/3⌋ replacements break this run
			buckets[runLen%3]++   // remember its residue for deletion targeting
		}
		i = j
	}

	if n < 6 {
		// Only insertions. Each insertion can add length and (if aimed well)
		// add a missing type or split a run. So we need at least (6-n) inserts
		// for length and at least `missing` inserts for types; one insert can
		// serve both, hence the max.
		return max(6-n, missing)
	}

	if n <= 20 {
		// No length change: fix runs with replacements, and reuse replacements
		// to add missing types. So max(replace, missing).
		return max(replace, missing)
	}

	// n > 20: we must delete exactly `over` characters.
	over := n - 20
	deletions := over // every over-length char must go

	// Spend deletions to reduce the replacements the runs still require.
	// A run of length L needs ⌊L/3⌋ replacements. Deleting characters can drop
	// L into a lower ⌊L/3⌋ bracket:
	//   • L ≡ 0 (mod 3): deleting 1 char saves 1 replacement (best value).
	//   • L ≡ 1 (mod 3): deleting 2 chars saves 1 replacement.
	//   • L ≡ 2 (mod 3): deleting 3 chars saves 1 replacement.
	// Apply the cheap savings first.

	// Pass 1: runs with len % 3 == 0 — 1 deletion each saves 1 replacement.
	if over > 0 {
		use := min(buckets[0], over)
		over -= use
		replace -= use // each such deletion removes one required replacement
	}
	// Pass 2: runs with len % 3 == 1 — 2 deletions each save 1 replacement.
	if over > 0 {
		use := min(buckets[1]*2, over)
		over -= use
		replace -= use / 2
	}
	// Pass 3: any remaining deletions — every 3 deletions save 1 replacement
	// (covers len % 3 == 2 and leftovers after chewing through longer runs).
	if over > 0 {
		replace -= over / 3
	}
	if replace < 0 {
		replace = 0 // can't need negative replacements
	}

	// After deletions, remaining run-replacements can still double as adding a
	// missing type, so the non-deletion cost is max(replace, missing).
	return deletions + max(replace, missing)
}

// max returns the larger of two ints.
func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

// min returns the smaller of two ints.
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func main() {
	fmt.Println("=== Approach 1: Brute Force (BFS over edits) ===")
	fmt.Println(bruteForceBFS("a"))   // expected 5
	fmt.Println(bruteForceBFS("aA1")) // expected 3
	fmt.Println(bruteForceBFS("aaa")) // expected 3 (need length+upper+digit; inserts/replaces overlap)

	fmt.Println("=== Approach 2: Greedy (Optimal) ===")
	fmt.Println(greedy("a"))                        // expected 5
	fmt.Println(greedy("aA1"))                      // expected 3
	fmt.Println(greedy("1337C0d3"))                 // expected 0
	fmt.Println(greedy("aaa"))                      // expected 3
	fmt.Println(greedy("aaaaa"))                    // expected 2  (len 5<6: max(6-5, missing=2)=max(1,2)=2)
	fmt.Println(greedy("aaaabbaaaabbaaaabbaaaabb")) // expected 6  (len 24: 4 deletions + max(replace=2, missing=2))
	fmt.Println(greedy("ABABABABABABABABABAB1"))    // expected 2  (len 21>20: 1 delete + missing lower =1 → max)
}
