package main

import "fmt"

// ── Approach 1: Full Generation with a Group Queue ───────────────────────────
//
// generateWithQueue builds the magical string one full character at a time by
// treating the already-built prefix as a queue of "how long is the next run".
//
// Intuition:
//
//	The magical string is self-describing: reading its own characters left to
//	right tells you the length of each successive run of identical digits.
//	s[0]=1 means the first run ("1") is length 1. s[1]=2 means the second run
//	("22") is length 2. s[2]=2 means the third run ("11") is length 2, and so
//	on. So a "read head" walking the string we are still building dictates how
//	many copies of the next digit (which strictly alternates 1,2,1,2,…) to
//	append. Because a run is at most length 2, the read head never overtakes
//	the write head, and the string grows without ever stalling.
//
// Algorithm:
//  1. Seed the string with the known prefix "122" (bootstrapping: the first
//     digit cannot be derived from an empty string).
//  2. Keep an index `head` that starts at 2 (s[2] is the first run length we
//     have not yet consumed) and a `next` digit that starts at 1.
//  3. While the string is shorter than n, append s[head] copies of `next`,
//     flip `next` between 1 and 2, and advance `head`.
//  4. Count the 1's among the first n characters.
//
// Time:  O(n) — every append writes at most 2 characters and each character is
//
//	produced exactly once; the final count is a single O(n) pass.
//
// Space: O(n) — the generated string of length ≈ n.
func generateWithQueue(n int) int {
	if n == 0 {
		return 0 // no characters → no 1's
	}
	if n <= 3 {
		// The bootstrap prefix "122" already covers n = 1,2,3. Its 1's:
		// "1"→1, "12"→1, "122"→1. So the count is always 1 here.
		return 1
	}

	s := []int{1, 2, 2} // known self-consistent seed of the magical string
	head := 2           // s[2]=2 is the next run length we have not applied yet
	next := 1           // the digit that the next run will consist of (alternates)

	// Grow until we have at least n characters; a run of length s[head]
	// appends that many copies of `next`.
	for len(s) < n {
		runLen := s[head] // how many identical digits the next group holds (1 or 2)
		for i := 0; i < runLen; i++ {
			s = append(s, next) // emit one digit of the current run
		}
		next ^= 3 // flip 1<->2 (1^3=2, 2^3=1) so runs alternate digit
		head++    // consume the next run-length instruction
	}

	ones := 0
	for i := 0; i < n; i++ { // count 1's within the first n characters only
		if s[i] == 1 {
			ones++
		}
	}
	return ones
}

// ── Approach 2: In-Place Two-Pointer, Count While Building (Optimal) ──────────
//
// twoPointers builds the string in the same self-referential way but folds the
// 1-count into the generation loop, so no separate final pass is needed. It is
// the canonical, most memory-frugal formulation.
//
// Intuition:
//
//	Identical mechanics to Approach 1 — a slow read pointer `i` supplies run
//	lengths for a fast write pointer that appends the alternating digit — but
//	we increment a running `ones` counter the moment we append a 1, and we
//	stop the instant the string reaches length n. This removes the second
//	sweep and the need to remember anything past `ones`.
//
// Algorithm:
//  1. Handle n ≤ 3 directly (prefix "122" always contributes exactly one 1).
//  2. Seed s = [1,2,2] with ones = 1 (the single 1 in the seed).
//  3. Read pointer i starts at 2; alternating digit starts at 1.
//  4. While len(s) < n: append s[i] copies of the current digit, and for every
//     appended 1 that lands at an index < n, bump ones. Flip the digit, i++.
//  5. Return ones.
//
// Time:  O(n) — one character produced per unit of work, counted as we go.
// Space: O(n) — the string buffer; no auxiliary structures.
func twoPointers(n int) int {
	if n == 0 {
		return 0 // empty prefix has zero 1's
	}
	if n <= 3 {
		return 1 // "1","12","122" each contain exactly one 1
	}

	s := []int{1, 2, 2} // magical-string seed
	ones := 1           // the seed contributes exactly one '1' (at index 0)
	i := 2              // read pointer: s[2] is the next run length to apply
	digit := 1          // digit of the run being appended right now (alternates)

	for len(s) < n {
		for c := 0; c < s[i]; c++ { // append s[i] copies of `digit`
			if len(s) >= n {
				break // never write past position n-1
			}
			s = append(s, digit)
			if digit == 1 { // count 1's inline, only for indices < n
				ones++
			}
		}
		digit ^= 3 // 1<->2 alternation for the next run
		i++        // advance the read pointer to the next run length
	}
	return ones
}

func main() {
	fmt.Println("=== Approach 1: Full Generation with a Group Queue ===")
	fmt.Printf("n=6   got=%d  expected 3\n", generateWithQueue(6))  // "122112" → three 1's
	fmt.Printf("n=1   got=%d  expected 1\n", generateWithQueue(1))  // "1"      → one 1
	fmt.Printf("n=0   got=%d  expected 0\n", generateWithQueue(0))  // ""       → zero 1's
	fmt.Printf("n=10  got=%d  expected 5\n", generateWithQueue(10)) // "1221121221" → five 1's

	fmt.Println("=== Approach 2: In-Place Two-Pointer, Count While Building (Optimal) ===")
	fmt.Printf("n=6   got=%d  expected 3\n", twoPointers(6))  // three 1's
	fmt.Printf("n=1   got=%d  expected 1\n", twoPointers(1))  // one 1
	fmt.Printf("n=0   got=%d  expected 0\n", twoPointers(0))  // zero 1's
	fmt.Printf("n=10  got=%d  expected 5\n", twoPointers(10)) // five 1's
}
