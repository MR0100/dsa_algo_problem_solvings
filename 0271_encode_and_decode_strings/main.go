package main

import (
	"fmt"
	"strconv"
	"strings"
)

// Encode and Decode Strings (LeetCode #271)
//
// Design an algorithm to encode a LIST of strings into a SINGLE string, which
// is then sent over the network and decoded back into the original list.
//
// The core difficulty: strings may contain ANY character — including whatever
// delimiter we might naively pick (commas, spaces, newlines, etc.). A robust
// codec must survive an adversary who deliberately embeds our delimiter inside
// the payload. Below are three encoding schemes, from a fragile naive one to
// two production-safe ones.

// ── Approach 1: Length-Prefix (Chunked Transfer / Optimal) ───────────────────
//
// lengthPrefixEncode / lengthPrefixDecode encode each string as
// "<length>#<payload>" concatenated together, mirroring HTTP chunked transfer
// encoding.
//
// Intuition:
//
//	A delimiter is only ambiguous if the payload can contain it. Instead of
//	searching for a delimiter INSIDE the payload, we announce the payload's
//	exact byte length up front. The decoder reads the number, reads the '#'
//	separator, then blindly consumes exactly that many bytes — so the payload
//	can contain '#', digits, or anything else with zero ambiguity.
//
// Algorithm (encode):
//  1. For each string s, append len(s), then '#', then s itself.
//  2. Concatenate all such chunks.
//
// Algorithm (decode):
//  1. i = 0. While i < len(encoded):
//  2. Read digits from i until the '#' → that number is the chunk length L.
//  3. The payload is the L bytes right after the '#'.
//  4. Append it, then advance i past the payload; repeat.
//
// Time:  O(N) encode and decode, N = total bytes across all strings.
// Space: O(N) for the output.
func lengthPrefixEncode(strs []string) string {
	var b strings.Builder
	for _, s := range strs {
		// len(s) is the BYTE length (Go strings are byte slices); this is what
		// the decoder will consume, so byte length — not rune count — is correct.
		b.WriteString(strconv.Itoa(len(s))) // announce payload size
		b.WriteByte('#')                    // separator between size and payload
		b.WriteString(s)                    // the raw payload, unescaped
	}
	return b.String()
}

func lengthPrefixDecode(encoded string) []string {
	res := []string{}
	i := 0
	for i < len(encoded) {
		// Scan forward to the '#' that terminates the length header.
		j := i
		for encoded[j] != '#' {
			j++
		}
		// encoded[i:j] is the ASCII length; parse it.
		length, _ := strconv.Atoi(encoded[i:j])
		// Payload starts right after '#' (at j+1) and is exactly `length` bytes.
		start := j + 1
		res = append(res, encoded[start:start+length])
		// Jump past this whole chunk to the next length header.
		i = start + length
	}
	return res
}

// ── Approach 2: Escaping a Delimiter ─────────────────────────────────────────
//
// escapeEncode / escapeDecode pick a delimiter and make it safe by escaping any
// occurrence of it (and of the escape char) inside the payload.
//
// Intuition:
//
//	Keep a human-readable delimiter (say ':') between strings, but first
//	neutralise every ':' and every escape marker inside each payload. We map
//	each payload character to a two-char safe encoding when needed, so the
//	only *unescaped* delimiter in the stream marks a true boundary.
//
//	Concretely: escape ':' as "::" ... but that collides with the delimiter
//	itself. The classic trick is a distinct escape rule: encode the char '#'
//	as "#h" and use "#:" as the real separator. Then a lone "#:" can only be a
//	boundary, because any literal '#' in the payload became "#h".
//
// Algorithm (encode): for each s, replace every '#' with "#h", then append the
// escaped string followed by the sentinel "#:".
//
// Algorithm (decode): split the stream on the sentinel "#:" (dropping the
// trailing empty piece), then un-escape each piece by turning "#h" back to '#'.
//
// Time:  O(N). Space: O(N).
func escapeEncode(strs []string) string {
	var b strings.Builder
	for _, s := range strs {
		// Escape the escape-introducer '#' so no literal '#' can be mistaken
		// for the start of our sentinel.
		esc := strings.ReplaceAll(s, "#", "#h")
		b.WriteString(esc)  // safe payload
		b.WriteString("#:") // sentinel: a '#' followed by ':' — never appears
		//                     inside an escaped payload (every payload '#' is "#h")
	}
	return b.String()
}

func escapeDecode(encoded string) []string {
	if encoded == "" {
		return []string{}
	}
	// Every original string ended with "#:"; splitting on it yields the pieces
	// plus one trailing "" (after the final sentinel), which we drop.
	parts := strings.Split(encoded, "#:")
	parts = parts[:len(parts)-1] // remove the trailing empty element
	res := make([]string, len(parts))
	for i, p := range parts {
		// Reverse the escaping: "#h" → "#".
		res[i] = strings.ReplaceAll(p, "#h", "#")
	}
	return res
}

// ── Approach 3: Non-ASCII Sentinel (Naive/Fragile) ───────────────────────────
//
// sentinelEncode / sentinelDecode join with a rare Unicode sentinel and split
// on it. Simple, but ONLY correct if payloads never contain that rune — so it
// is a fragile approach shown for contrast, not for production.
//
// Intuition:
//
//	If we assume input strings never contain some exotic character, we can
//	just Join/Split on it. This is what most "just use a delimiter" answers
//	do; it fails the moment an adversary includes the sentinel.
//
// Algorithm (encode): join the list with the sentinel rune.
// Algorithm (decode): split the string on the sentinel rune.
//
// Time:  O(N). Space: O(N).
func sentinelEncode(strs []string) string {
	// U+241F ("SYMBOL FOR UNIT SEPARATOR") stands in for a byte unlikely to
	// occur in normal text. This is the fragile assumption.
	return strings.Join(strs, "␟")
}

func sentinelDecode(encoded string) []string {
	if encoded == "" {
		// Join of a single "" and a truly empty list both yield "" — we can't
		// distinguish, so we standardise on the single-empty-string case being
		// unreachable here and return [] for the empty stream.
		return []string{}
	}
	return strings.Split(encoded, "␟")
}

// equalSlices is a tiny helper so main() can print a clean PASS/FAIL.
func equalSlices(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

func main() {
	// Official examples.
	ex1 := []string{"Hello", "World"}
	ex2 := []string{""}
	// Adversarial example: payloads that embed our delimiters/escape chars.
	ex3 := []string{"we", "say", ":", "#", "yes", "3#abc"}

	fmt.Println("=== Approach 1: Length-Prefix ===")
	fmt.Println(equalSlices(lengthPrefixDecode(lengthPrefixEncode(ex1)), ex1)) // expected true
	fmt.Println(equalSlices(lengthPrefixDecode(lengthPrefixEncode(ex2)), ex2)) // expected true
	fmt.Println(equalSlices(lengthPrefixDecode(lengthPrefixEncode(ex3)), ex3)) // expected true

	fmt.Println("=== Approach 2: Escaping a Delimiter ===")
	fmt.Println(equalSlices(escapeDecode(escapeEncode(ex1)), ex1)) // expected true
	fmt.Println(equalSlices(escapeDecode(escapeEncode(ex2)), ex2)) // expected true
	fmt.Println(equalSlices(escapeDecode(escapeEncode(ex3)), ex3)) // expected true

	fmt.Println("=== Approach 3: Non-ASCII Sentinel (Fragile) ===")
	fmt.Println(equalSlices(sentinelDecode(sentinelEncode(ex1)), ex1)) // expected true
	fmt.Println(equalSlices(sentinelDecode(sentinelEncode(ex3)), ex3)) // expected true
}
