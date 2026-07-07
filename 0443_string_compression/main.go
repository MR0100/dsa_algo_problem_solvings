package main

import (
	"fmt"
	"strconv"
)

// ── Approach 1: Extra Buffer (Build Then Copy Back) ──────────────────────────
//
// extraBuffer solves String Compression by building the compressed run-length
// string in a separate slice, then copying it back into chars.
//
// Intuition:
//
//	Compression is a run-length encoding: for each maximal run of one character,
//	emit the character, and (if the run is longer than 1) its length's digits.
//	Building into a fresh buffer keeps the logic obvious; we copy back to satisfy
//	the "store in chars" requirement (though it uses O(n) extra space).
//
// Algorithm:
//  1. Scan runs: for group start i, extend j while chars[j] == chars[i].
//  2. Append chars[i]; if run length > 1, append each digit of the length.
//  3. Copy the built buffer into chars and return its length.
//
// Time:  O(n) — single scan plus a copy.
// Space: O(n) — the temporary buffer (does NOT meet the O(1) follow-up).
func extraBuffer(chars []byte) int {
	out := make([]byte, 0, len(chars)) // compressed result
	i := 0
	for i < len(chars) {
		ch := chars[i] // the character of this run
		j := i
		// Extend the run while the same character repeats.
		for j < len(chars) && chars[j] == ch {
			j++
		}
		count := j - i        // length of the run
		out = append(out, ch) // always emit the character
		if count > 1 {
			// Emit the count's decimal digits (may be multiple, e.g. "12").
			out = append(out, []byte(strconv.Itoa(count))...)
		}
		i = j // jump to the next run
	}
	copy(chars, out) // write compressed content back into chars in place
	return len(out)  // new logical length
}

// ── Approach 2: Two Pointers In-Place (Optimal) ──────────────────────────────
//
// twoPointers solves String Compression in O(1) extra space using a read pointer
// and a write pointer over the same array.
//
// Intuition:
//
//	A run compresses to at most its own length (1 char + up to a few digits),
//	so the WRITE pointer can never overtake the READ pointer — it is always safe
//	to overwrite chars in place. Read groups the current run, then write emits the
//	character and, if the run > 1, the digits of the count.
//
// Algorithm:
//  1. read = 0, write = 0.
//  2. While read < n: mark ch = chars[read]; advance read over the whole run,
//     counting its length.
//  3. chars[write++] = ch. If count > 1, convert count to a digit string and
//     write each digit at chars[write++].
//  4. Return write.
//
// Time:  O(n) — read visits each character once; write emits ≤ read total.
// Space: O(1) — only pointers and a tiny digit buffer (bounded by log10(n) ≤ 6).
func twoPointers(chars []byte) int {
	write := 0 // next slot to write compressed output
	read := 0  // scanning position over the input runs
	for read < len(chars) {
		ch := chars[read] // character starting this run
		count := 0
		// Consume the entire run of ch, counting its length.
		for read < len(chars) && chars[read] == ch {
			read++
			count++
		}
		chars[write] = ch // write the run's character
		write++
		if count > 1 {
			// Write the count's digits in order (in-place, left to right).
			for _, d := range strconv.Itoa(count) {
				chars[write] = byte(d) // rune digit fits in a byte ('0'..'9')
				write++
			}
		}
	}
	return write // length of the compressed array
}

// stringify renders the first `n` bytes of chars for readable expected-output
// comparison in main().
func stringify(chars []byte, n int) string {
	return string(chars[:n])
}

func main() {
	fmt.Println("=== Approach 1: Extra Buffer (Build Then Copy Back) ===")
	c1 := []byte{'a', 'a', 'b', 'b', 'c', 'c', 'c'}
	n1 := extraBuffer(c1)
	fmt.Printf("[a a b b c c c]  len=%d chars=%q  expected len=6 chars=\"a2b2c3\"\n", n1, stringify(c1, n1))

	c2 := []byte{'a', 'b', 'b', 'b', 'b', 'b', 'b', 'b', 'b', 'b', 'b', 'b', 'b'}
	n2 := extraBuffer(c2)
	fmt.Printf("[a b*12]         len=%d chars=%q  expected len=4 chars=\"ab12\"\n", n2, stringify(c2, n2))

	c3 := []byte{'a'}
	n3 := extraBuffer(c3)
	fmt.Printf("[a]              len=%d chars=%q  expected len=1 chars=\"a\"\n", n3, stringify(c3, n3))

	fmt.Println("=== Approach 2: Two Pointers In-Place (Optimal) ===")
	d1 := []byte{'a', 'a', 'b', 'b', 'c', 'c', 'c'}
	m1 := twoPointers(d1)
	fmt.Printf("[a a b b c c c]  len=%d chars=%q  expected len=6 chars=\"a2b2c3\"\n", m1, stringify(d1, m1))

	d2 := []byte{'a', 'b', 'b', 'b', 'b', 'b', 'b', 'b', 'b', 'b', 'b', 'b', 'b'}
	m2 := twoPointers(d2)
	fmt.Printf("[a b*12]         len=%d chars=%q  expected len=4 chars=\"ab12\"\n", m2, stringify(d2, m2))

	d3 := []byte{'a'}
	m3 := twoPointers(d3)
	fmt.Printf("[a]              len=%d chars=%q  expected len=1 chars=\"a\"\n", m3, stringify(d3, m3))
}
