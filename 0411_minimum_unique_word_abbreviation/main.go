package main

import (
	"fmt"
	"strconv"
	"strings"
)

// Minimum Unique Word Abbreviation
//
// Given a `target` string and a `dictionary`, find an abbreviation of `target`
// with the smallest possible *display* length such that it does not collide with
// (is not also a valid abbreviation of) any word in the dictionary.
//
// Core observations shared by every approach below:
//
//  1. Only dictionary words with the SAME length as target can ever collide —
//     an abbreviation encodes the exact length, so different-length words are
//     automatically safe and can be dropped.
//  2. An abbreviation is fully described by WHICH positions of target we keep as
//     literal letters (the maximal gaps between kept letters become numbers).
//     Encode a "kept-positions" choice as a bitmask `cand` over target's indices.
//  3. For each same-length dictionary word build a diff-mask `diff` whose bit i is
//     1 exactly where that word differs from target. A choice `cand` DISTINGUISHES
//     target from that word iff `cand & diff != 0` (we kept at least one position
//     where they differ). The abbreviation is valid iff this holds for EVERY word.
//  4. Among valid choices we minimise the abbreviation's display length
//     (letters counted 1 each, every maximal run of skipped positions counts as
//     one number token), NOT the number of kept letters.

// abbrevLen returns the display length of the abbreviation of a length-n word
// that keeps exactly the positions whose bit is set in `mask`.
//
// A kept position contributes 1 (the letter). A maximal run of unkept positions
// contributes 1 (a single number like "12"). We sweep left to right and count a
// number token once per run of consecutive skipped positions.
func abbrevLen(mask, n int) int {
	length := 0
	i := 0
	for i < n {
		if mask&(1<<i) != 0 { // this position is kept as a literal letter
			length++ // the letter costs 1
			i++
			continue
		}
		// A run of skipped positions collapses into ONE number token.
		length++ // the number token costs 1 (regardless of how many digits)
		for i < n && mask&(1<<i) == 0 {
			i++ // consume the whole run of skipped positions
		}
	}
	return length
}

// buildAbbrev renders the actual abbreviation string for a kept-positions mask,
// e.g. target="apple", mask keeping only index 1 → "1p3".
func buildAbbrev(target string, mask int) string {
	var sb strings.Builder
	n := len(target)
	i := 0
	for i < n {
		if mask&(1<<i) != 0 { // kept letter
			sb.WriteByte(target[i])
			i++
			continue
		}
		count := 0 // length of the skipped run → becomes the number
		for i < n && mask&(1<<i) == 0 {
			count++
			i++
		}
		sb.WriteString(strconv.Itoa(count)) // one number token for the whole run
	}
	return sb.String()
}

// buildDiffMasks filters the dictionary to same-length words and returns, for
// each, the bitmask of positions where it differs from target.
func buildDiffMasks(target string, dictionary []string) []int {
	n := len(target)
	diffs := make([]int, 0, len(dictionary))
	for _, w := range dictionary {
		if len(w) != n {
			continue // different length can never collide with our abbreviation
		}
		mask := 0
		for i := 0; i < n; i++ {
			if w[i] != target[i] {
				mask |= 1 << i // mark this differing position
			}
		}
		// Note: an identical same-length word (mask==0) can never be
		// distinguished; per constraints target is not required to be absent,
		// but if such a word existed no abbreviation could ever be unique. We
		// keep it so validity naturally fails, matching the problem's intent.
		diffs = append(diffs, mask)
	}
	return diffs
}

// ── Approach 1: Brute Force Over All Kept-Masks ──────────────────────────────
//
// bruteForce solves Minimum Unique Word Abbreviation by trying EVERY subset of
// kept positions (all 2^n bitmasks), keeping those that distinguish target from
// all same-length dictionary words, and returning the one with the smallest
// display length.
//
// Intuition:
//
//	A valid abbreviation ⇔ a kept-positions mask `cand` such that cand & diff != 0
//	for every same-length dictionary diff-mask. Enumerate all 2^n masks, test that
//	predicate, and track the minimum abbrevLen. The constraint log2(n)+m ≤ 20 (with
//	m ≤ 21) keeps 2^m within reach.
//
// Algorithm:
//  1. Build diff-masks for same-length words. If none, "n" (all abbreviated) is optimal.
//  2. For cand in 0 .. 2^n - 1: valid iff cand & diff != 0 for all diffs.
//  3. Track cand with minimum abbrevLen; render it with buildAbbrev.
//
// Time:  O(2^n · D) — 2^n masks, each checked against D same-length words.
// Space: O(D) — the diff-mask list.
func bruteForce(target string, dictionary []string) string {
	n := len(target)
	diffs := buildDiffMasks(target, dictionary)

	// No same-length word to avoid → abbreviate everything to just the number n.
	if len(diffs) == 0 {
		return strconv.Itoa(n)
	}

	bestMask := (1 << n) - 1 // fallback: keep all letters (always valid, length n)
	bestLen := n             // its display length is exactly n
	for cand := 0; cand < (1 << n); cand++ {
		valid := true
		for _, d := range diffs {
			if cand&d == 0 { // fails to distinguish target from this word
				valid = false
				break
			}
		}
		if !valid {
			continue
		}
		if l := abbrevLen(cand, n); l < bestLen { // strictly shorter → adopt
			bestLen = l
			bestMask = cand
		}
	}
	return buildAbbrev(target, bestMask)
}

// ── Approach 2: Backtracking Abbreviation Builder ────────────────────────────
//
// backtracking solves Minimum Unique Word Abbreviation by recursively deciding,
// position by position, whether to abbreviate a run of characters or keep the
// next letter, pruning branches once they cannot beat the best answer so far.
//
// Intuition:
//
//	Build the abbreviation left to right. At each position we either (a) keep the
//	current letter, or (b) skip a run of k >= 1 characters as one number. Track the
//	kept-positions mask so we can validate against the diff-masks at the end. Two
//	prunes keep it fast: stop a branch whose running token count already reaches
//	the best complete answer, and only accept a completed abbreviation if its mask
//	distinguishes target from every dictionary word.
//
// Algorithm:
//  1. Build diff-masks; if none, return "n".
//  2. Recurse over (pos, tokensSoFar, keptMask). At pos:
//     - skip a run of length 1..(n-pos): +1 token, advance pos.
//     - keep the letter at pos: +1 token, set bit, advance pos+1.
//  3. At pos == n, if keptMask beats all diffs and token count is smaller, record it.
//
// Time:  O(2^n · D) worst case (same search space as brute force, but pruned).
// Space: O(n) recursion depth.
func backtracking(target string, dictionary []string) string {
	n := len(target)
	diffs := buildDiffMasks(target, dictionary)
	if len(diffs) == 0 {
		return strconv.Itoa(n)
	}

	bestLen := n + 1 // no valid abbreviation found yet (n letters is the ceiling)
	bestMask := (1 << n) - 1

	// distinguishes reports whether keeping `mask` separates target from every word.
	distinguishes := func(mask int) bool {
		for _, d := range diffs {
			if mask&d == 0 {
				return false
			}
		}
		return true
	}

	var dfs func(pos, tokens, keptMask int)
	dfs = func(pos, tokens, keptMask int) {
		if tokens >= bestLen { // cannot improve on the best complete answer → prune
			return
		}
		if pos == n {
			// Completed an abbreviation of display length `tokens`.
			if distinguishes(keptMask) && tokens < bestLen {
				bestLen = tokens
				bestMask = keptMask
			}
			return
		}
		// Option A: abbreviate a run of length runLen as a single number token.
		for runLen := 1; pos+runLen <= n; runLen++ {
			dfs(pos+runLen, tokens+1, keptMask) // one token for the whole run
		}
		// Option B: keep the letter at pos as a literal (also one token).
		dfs(pos+1, tokens+1, keptMask|(1<<pos))
	}
	dfs(0, 0, 0)

	return buildAbbrev(target, bestMask)
}

func main() {
	// Example 1: target = "apple", dictionary = ["blade"] → "a4"
	// ("5" and "4e" both collide with "blade"; "a4" keeps position 0 where they
	//  differ — 'a' vs 'b' — so it is unique, display length 2.)
	fmt.Println("=== Approach 1: Brute Force Over All Kept-Masks ===")
	fmt.Println(bruteForce("apple", []string{"blade"})) // expected a4

	// Example 2: target = "apple", dictionary = ["plain","amber","blade"] → "1p3"
	// (Length 3; other valid minima exist: ap3, a3e, 2p2, 3le, 3l1. We print the
	//  smallest kept-mask that achieves the minimum length, which is "1p3".)
	fmt.Println(bruteForce("apple", []string{"plain", "amber", "blade"})) // expected 1p3

	// No same-length word → abbreviate everything.
	fmt.Println(bruteForce("apple", []string{"banana"})) // expected 5

	fmt.Println("=== Approach 2: Backtracking Abbreviation Builder ===")
	fmt.Println(backtracking("apple", []string{"blade"}))                   // expected a4
	fmt.Println(backtracking("apple", []string{"plain", "amber", "blade"})) // expected 1p3
	fmt.Println(backtracking("apple", []string{"banana"}))                  // expected 5
}
