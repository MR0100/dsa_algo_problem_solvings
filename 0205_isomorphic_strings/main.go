package main

import "fmt"

// ── Approach 1: Brute Force (Pairwise Consistency) ───────────────────────────
//
// bruteForce solves Isomorphic Strings by checking every pair of positions
// for a mapping contradiction.
//
// Intuition:
//
//	A character mapping exists and is one-to-one exactly when the two strings
//	have the same "equality pattern": for every pair of positions (j, i),
//	s[j] == s[i] must hold precisely when t[j] == t[i]. If s repeats where t
//	doesn't, one s-char would need two images; if t repeats where s doesn't,
//	two s-chars would collide on one image. Check all pairs directly.
//
// Algorithm:
//  1. For each i from 0..n−1, for each j < i:
//     if (s[i] == s[j]) != (t[i] == t[j]), return false.
//  2. No contradiction found → return true.
//
// Time:  O(n²) — all ~n²/2 position pairs are compared.
// Space: O(1) — no auxiliary structures at all.
func bruteForce(s string, t string) bool {
	for i := 0; i < len(s); i++ {
		for j := 0; j < i; j++ {
			// The equality patterns must agree at every pair of positions:
			// a mismatch in either direction breaks the bijection.
			if (s[i] == s[j]) != (t[i] == t[j]) {
				return false
			}
		}
	}
	return true // identical equality structure → an isomorphism exists
}

// ── Approach 2: Two Hash Maps ────────────────────────────────────────────────
//
// twoHashMaps solves Isomorphic Strings by building the forward (s→t) and
// backward (t→s) character mappings simultaneously.
//
// Intuition:
//
//	Isomorphism is a consistent AND injective replacement. One map s→t
//	enforces consistency ("every 'a' becomes the same thing"); it cannot
//	catch two different s-chars claiming the same t-char (e.g. s="badc",
//	t="baba": d→b collides with b→b). A second map t→s enforces that
//	injectivity. Walk once, growing both maps and failing on any conflict.
//
// Algorithm:
//  1. sToT, tToS = empty maps.
//  2. For each position i with bytes a = s[i], b = t[i]:
//     - if sToT[a] exists and differs from b → false;
//     - if tToS[b] exists and differs from a → false;
//     - record sToT[a] = b, tToS[b] = a.
//  3. Return true if the whole string passes.
//
// Time:  O(n) — one pass with O(1) average map operations.
// Space: O(k) — k = alphabet size (≤ 256 distinct byte keys per map).
func twoHashMaps(s string, t string) bool {
	sToT := map[byte]byte{} // established forward mapping s-char → t-char
	tToS := map[byte]byte{} // established backward mapping t-char → s-char
	for i := 0; i < len(s); i++ {
		a, b := s[i], t[i]
		// Forward check: a must always map to the same image.
		if mapped, ok := sToT[a]; ok && mapped != b {
			return false // a already maps elsewhere — inconsistent
		}
		// Backward check: b must not already be claimed by another s-char.
		if mapped, ok := tToS[b]; ok && mapped != a {
			return false // two s-chars would share the image b — not injective
		}
		sToT[a] = b // (re-)recording an identical pair is harmless
		tToS[b] = a
	}
	return true
}

// ── Approach 3: First-Occurrence Encoding (Optimal) ──────────────────────────
//
// firstOccurrenceEncoding solves Isomorphic Strings by comparing where each
// character was last seen, using two fixed-size arrays instead of maps.
//
// Intuition:
//
//	Replace every character by the index at which it was last seen: two
//	strings are isomorphic iff these index sequences coincide (both "egg"
//	and "add" encode to 0,1,2-with-repeat-structure). It suffices to check,
//	at every position, that s[i] and t[i] were last seen at the SAME earlier
//	position — storing i+1 (so the zero value means "never seen") in two
//	256-slot arrays makes each check two array reads.
//
// Algorithm:
//  1. lastSeenS, lastSeenT = [256]int zero arrays.
//  2. For each i: if lastSeenS[s[i]] != lastSeenT[t[i]], return false
//     (one char is new while the other isn't, or they recur from
//     different positions).
//  3. Store i+1 into both entries; return true at the end.
//
// Time:  O(n) — one pass, two array reads + two writes per position.
// Space: O(1) — two fixed 256-entry arrays, independent of n.
func firstOccurrenceEncoding(s string, t string) bool {
	var lastSeenS, lastSeenT [256]int // last position + 1; 0 = never seen
	for i := 0; i < len(s); i++ {
		// Both characters must have the same "history": either both brand
		// new (0 == 0) or both last seen at the same index.
		if lastSeenS[s[i]] != lastSeenT[t[i]] {
			return false
		}
		lastSeenS[s[i]] = i + 1 // +1 shift keeps index 0 distinguishable
		lastSeenT[t[i]] = i + 1
	}
	return true
}

func main() {
	fmt.Println("=== Approach 1: Brute Force (Pairwise Consistency) ===")
	fmt.Printf("s=\"egg\", t=\"add\"      got=%t  expected true\n", bruteForce("egg", "add"))
	fmt.Printf("s=\"foo\", t=\"bar\"      got=%t  expected false\n", bruteForce("foo", "bar"))
	fmt.Printf("s=\"paper\", t=\"title\"  got=%t  expected true\n", bruteForce("paper", "title"))
	fmt.Printf("s=\"badc\", t=\"baba\"    got=%t  expected false\n", bruteForce("badc", "baba")) // injectivity edge

	fmt.Println("=== Approach 2: Two Hash Maps ===")
	fmt.Printf("s=\"egg\", t=\"add\"      got=%t  expected true\n", twoHashMaps("egg", "add"))
	fmt.Printf("s=\"foo\", t=\"bar\"      got=%t  expected false\n", twoHashMaps("foo", "bar"))
	fmt.Printf("s=\"paper\", t=\"title\"  got=%t  expected true\n", twoHashMaps("paper", "title"))
	fmt.Printf("s=\"badc\", t=\"baba\"    got=%t  expected false\n", twoHashMaps("badc", "baba"))

	fmt.Println("=== Approach 3: First-Occurrence Encoding (Optimal) ===")
	fmt.Printf("s=\"egg\", t=\"add\"      got=%t  expected true\n", firstOccurrenceEncoding("egg", "add"))
	fmt.Printf("s=\"foo\", t=\"bar\"      got=%t  expected false\n", firstOccurrenceEncoding("foo", "bar"))
	fmt.Printf("s=\"paper\", t=\"title\"  got=%t  expected true\n", firstOccurrenceEncoding("paper", "title"))
	fmt.Printf("s=\"badc\", t=\"baba\"    got=%t  expected false\n", firstOccurrenceEncoding("badc", "baba"))
}
