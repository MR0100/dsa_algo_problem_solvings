package main

import "fmt"

// seqLen is the fixed window size the problem asks about: 10-letter sequences.
const seqLen = 10

// ── Approach 1: Brute Force ──────────────────────────────────────────────────
//
// bruteForce solves Repeated DNA Sequences by comparing every 10-letter window
// against every other window with nested loops.
//
// Intuition:
//
//	A sequence is an answer iff it starts at two different indices. So for
//	each window, scan the rest of the string for a second occurrence. To
//	avoid emitting duplicates, a window is only *reported* by its first
//	occurrence: if the same text already appeared at an earlier index, the
//	current window is skipped because that earlier pass already handled it.
//
// Algorithm:
//  1. For each start index i (0 .. n-10):
//     a. If s[i:i+10] also occurs at some j < i, skip (already reported).
//     b. Otherwise, if s[i:i+10] occurs again at some j > i, append it.
//  2. Return the collected sequences (ordered by first occurrence).
//
// Time:  O(n^2 · L) — for each of the n windows we scan up to n others, each comparison costing L = 10.
// Space: O(1) extra — ignoring the output slice; comparisons reuse string views.
func bruteForce(s string) []string {
	result := []string{} // answers, ordered by first occurrence
	n := len(s)
	for i := 0; i+seqLen <= n; i++ {
		window := s[i : i+seqLen] // candidate 10-letter sequence
		// (a) reported already? — check all earlier windows for the same text
		seenBefore := false
		for j := 0; j < i; j++ {
			if s[j:j+seqLen] == window {
				seenBefore = true // its first occurrence handled the reporting
				break
			}
		}
		if seenBefore {
			continue
		}
		// (b) does it repeat later? — one later match is enough to qualify
		for j := i + 1; j+seqLen <= n; j++ {
			if s[j:j+seqLen] == window {
				result = append(result, window) // first occurrence reports it
				break
			}
		}
	}
	return result
}

// ── Approach 2: Hash Map Counting ────────────────────────────────────────────
//
// hashMap solves Repeated DNA Sequences by counting every 10-letter substring
// in a map and reporting each one the moment its count reaches two.
//
// Intuition:
//
//	Slide a fixed window of 10 across the string and tally each snapshot in a
//	substring → count map. Appending exactly when a count *becomes* 2 (not >2)
//	reports every repeated sequence once, in deterministic left-to-right
//	detection order — no separate dedup set needed.
//
// Algorithm:
//  1. For every start index i, take sub = s[i:i+10] and increment counts[sub].
//  2. If counts[sub] == 2, append sub to the result (first moment it repeats).
//  3. Return the result.
//
// Time:  O(n · L) — n windows; hashing/slicing each 10-byte window costs L = 10.
// Space: O(n · L) — the map may hold a distinct 10-byte key per window.
func hashMap(s string) []string {
	result := []string{}          // sequences confirmed repeated, in detection order
	counts := map[string]int{}    // 10-letter window → occurrences so far
	for i := 0; i+seqLen <= len(s); i++ {
		sub := s[i : i+seqLen] // current window snapshot
		counts[sub]++
		if counts[sub] == 2 { // exactly the moment it becomes repeated
			result = append(result, sub) // == 2 (not >= 2) guarantees a single report
		}
	}
	return result
}

// ── Approach 3: Rolling Hash / Bit Manipulation (Optimal) ────────────────────
//
// rollingHash solves Repeated DNA Sequences by packing each window into a
// 20-bit integer fingerprint maintained in O(1) per step.
//
// Intuition:
//
//	The alphabet has only 4 letters, so each letter fits in 2 bits and a
//	10-letter window fits in exactly 20 bits — a perfect, collision-free
//	hash. Sliding the window right is just "shift left by 2, OR the new
//	letter, mask to 20 bits", so each step costs O(1) instead of re-hashing
//	10 characters, and set keys are cheap ints instead of 10-byte strings.
//
// Algorithm:
//  1. Map A→0, C→1, G→2, T→3 (2 bits each).
//  2. Maintain window = the last 10 letters encoded in the low 20 bits:
//     window = ((window << 2) | code) & 0xFFFFF.
//  3. Once i >= 9 the window is full: if its fingerprint is in seen but not
//     yet in added, append s[i-9:i+1]; then mark the fingerprint seen.
//  4. Return the result.
//
// Time:  O(n) — one pass; each step is O(1) shift/mask/set work (no 10-char re-hash).
// Space: O(n) — the two integer sets hold at most one 20-bit key per window.
func rollingHash(s string) []string {
	result := []string{}
	if len(s) < seqLen {
		return result // no 10-letter window even exists
	}
	// 2-bit code per nucleotide — 10 letters pack into 20 bits losslessly
	code := map[byte]uint32{'A': 0, 'C': 1, 'G': 2, 'T': 3}
	const mask = 1<<(2*seqLen) - 1 // 0xFFFFF: keep only the newest 10 letters
	var window uint32              // rolling 20-bit fingerprint of the last 10 letters
	seen := map[uint32]bool{}      // fingerprints observed at least once
	added := map[uint32]bool{}     // fingerprints already appended to result
	for i := 0; i < len(s); i++ {
		// push the new letter into the low bits, drop the letter that fell out
		window = (window<<2 | code[s[i]]) & mask
		if i >= seqLen-1 { // window spans a full 10 letters ending at i
			if seen[window] && !added[window] {
				result = append(result, s[i-seqLen+1:i+1]) // second sighting → report once
				added[window] = true                       // never report this text again
			}
			seen[window] = true // record this window for future sightings
		}
	}
	return result
}

func main() {
	// Example 1: s = "AAAAACCCCCAAAAACCCCCCAAAAAGGGTTT"
	// Example 2: s = "AAAAAAAAAAAAA"
	ex1 := "AAAAACCCCCAAAAACCCCCCAAAAAGGGTTT"
	ex2 := "AAAAAAAAAAAAA"

	fmt.Println("=== Approach 1: Brute Force ===")
	fmt.Println(bruteForce(ex1)) // expected: [AAAAACCCCC CCCCCAAAAA]
	fmt.Println(bruteForce(ex2)) // expected: [AAAAAAAAAA]

	fmt.Println("=== Approach 2: Hash Map Counting ===")
	fmt.Println(hashMap(ex1)) // expected: [AAAAACCCCC CCCCCAAAAA]
	fmt.Println(hashMap(ex2)) // expected: [AAAAAAAAAA]

	fmt.Println("=== Approach 3: Rolling Hash / Bit Manipulation (Optimal) ===")
	fmt.Println(rollingHash(ex1)) // expected: [AAAAACCCCC CCCCCAAAAA]
	fmt.Println(rollingHash(ex2)) // expected: [AAAAAAAAAA]
}
