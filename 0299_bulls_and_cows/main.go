package main

import (
	"fmt"
	"strconv"
)

// Bulls and Cows
//
// You play a guessing game with a secret number. Given the secret and a guess
// (equal-length numeric strings), return a hint "xAyB" where:
//   - x (bulls)  = digits that match in both value AND position.
//   - y (cows)   = digits present in the secret but at a wrong position (each
//                  digit can only be counted once, honouring multiplicities).

// ── Approach 1: Two-Pass with Digit Counts ───────────────────────────────────
//
// twoPassCount first counts exact-position matches (bulls), then uses frequency
// tables of the leftover digits to count cows.
//
// Intuition:
//
//	Bulls are trivial: same digit, same index. For cows, ignore all bull
//	positions and tally how many of each digit remain in the secret and in the
//	guess separately. For each digit 0..9, the number of cows it can form is the
//	minimum of its leftover counts on both sides.
//
// Algorithm:
//  1. Loop indices: if secret[i] == guess[i], bulls++; else increment
//     secretCount[secret[i]] and guessCount[guess[i]].
//  2. cows = sum over d in 0..9 of min(secretCount[d], guessCount[d]).
//  3. Return "{bulls}A{cows}B".
//
// Time:  O(n) — one pass plus a fixed 10-digit merge.
// Space: O(1) — two size-10 arrays.
func twoPassCount(secret string, guess string) string {
	var secretCount, guessCount [10]int
	bulls := 0

	for i := 0; i < len(secret); i++ {
		if secret[i] == guess[i] {
			bulls++ // exact match in value and position
		} else {
			// tally only NON-bull digits for later cow matching
			secretCount[secret[i]-'0']++
			guessCount[guess[i]-'0']++
		}
	}

	cows := 0
	for d := 0; d < 10; d++ {
		// a digit forms cows up to the smaller of the two leftover counts
		cows += min(secretCount[d], guessCount[d])
	}

	return strconv.Itoa(bulls) + "A" + strconv.Itoa(cows) + "B"
}

// ── Approach 2: Single-Pass with a Signed Balance Array (Optimal) ────────────
//
// singlePassBalance computes bulls and cows in ONE loop using a single count
// array whose sign indicates which side currently has a surplus of a digit.
//
// Intuition:
//
//	Keep one array `count`. When we see a non-bull secret digit we do count[d]++;
//	if it was already negative, the guess had earlier "requested" this digit, so
//	a cow is formed. Symmetrically, a non-bull guess digit does count[d]--; if it
//	was positive, the secret had a surplus, forming a cow. Each pending surplus
//	on one side is matched at most once by the other side.
//
// Algorithm:
//  1. For each i:
//     - if secret[i] == guess[i]: bulls++, continue.
//     - s = secret[i]: if count[s] < 0 -> cows++ (guess was waiting); count[s]++.
//     - g = guess[i]:  if count[g] > 0 -> cows++ (secret was waiting); count[g]--.
//  2. Return "{bulls}A{cows}B".
//
// Time:  O(n) — a single pass.
// Space: O(1) — one size-10 array.
func singlePassBalance(secret string, guess string) string {
	var count [10]int // >0: secret surplus, <0: guess surplus, for that digit
	bulls, cows := 0, 0

	for i := 0; i < len(secret); i++ {
		if secret[i] == guess[i] {
			bulls++ // exact positional match
			continue
		}
		s := secret[i] - '0'
		if count[s] < 0 { // guess had already seen this digit unmatched
			cows++ // pair them up now
		}
		count[s]++ // secret now holds a surplus (or cancels a guess deficit)

		g := guess[i] - '0'
		if count[g] > 0 { // secret had a surplus of this digit
			cows++ // pair them up now
		}
		count[g]-- // guess now holds a surplus (or cancels a secret surplus)
	}

	return strconv.Itoa(bulls) + "A" + strconv.Itoa(cows) + "B"
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func main() {
	fmt.Println("=== Approach 1: Two-Pass with Digit Counts ===")
	fmt.Println(twoPassCount("1807", "7810")) // expected 1A3B
	fmt.Println(twoPassCount("1123", "0111")) // expected 1A1B
	fmt.Println(twoPassCount("1234", "1234")) // expected 4A0B
	fmt.Println(twoPassCount("1234", "5678")) // expected 0A0B

	fmt.Println("=== Approach 2: Single-Pass Balance (Optimal) ===")
	fmt.Println(singlePassBalance("1807", "7810")) // expected 1A3B
	fmt.Println(singlePassBalance("1123", "0111")) // expected 1A1B
	fmt.Println(singlePassBalance("1234", "1234")) // expected 4A0B
	fmt.Println(singlePassBalance("1234", "5678")) // expected 0A0B
}
