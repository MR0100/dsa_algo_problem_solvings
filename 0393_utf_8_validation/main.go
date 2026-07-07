package main

import "fmt"

// ── Approach 1: Count Continuation Bytes via Bit Masks (Optimal) ──────────────
//
// validUtf8 checks whether the integer array (only the low 8 bits of each entry
// matter) forms a valid UTF-8 encoding.
//
// UTF-8 rules:
//   - 1-byte char: 0xxxxxxx
//   - n-byte char (n=2..4): the first byte starts with n leading 1s then a 0,
//     and is followed by (n-1) continuation bytes each starting with 10xxxxxx.
//
// Intuition:
//
//	Walk the bytes left to right. When we are NOT mid-character, inspect the
//	leading bits of the current byte to decide how many continuation bytes must
//	follow (0, 1, 2, or 3). Then verify exactly that many following bytes each
//	begin with the bit pattern 10. Any deviation — a bad leader (like 10xxxxxx
//	when a fresh char is expected, or 11111xxx which never starts a char), too
//	few continuation bytes, or a continuation byte that doesn't start with 10 —
//	makes the whole array invalid.
//
// Algorithm:
//  1. remaining = 0 (continuation bytes still expected).
//  2. For each byte b (low 8 bits):
//     - If remaining == 0 (start of a new character):
//     • 0xxxxxxx (b < 0x80) → 1-byte char, remaining stays 0.
//     • 110xxxxx (b>>5 == 0b110) → remaining = 1.
//     • 1110xxxx (b>>4 == 0b1110) → remaining = 2.
//     • 11110xxx (b>>3 == 0b11110) → remaining = 3.
//     • anything else → invalid.
//     - Else (mid-character): b must be 10xxxxxx (b>>6 == 0b10); remaining--.
//  3. Valid iff remaining == 0 at the end (no truncated char).
//
// Time:  O(n) — one pass over the bytes.
// Space: O(1) — a single counter.
func validUtf8(data []int) bool {
	remaining := 0 // continuation bytes still expected for the current character
	for _, num := range data {
		b := num & 0xFF // keep only the least-significant 8 bits
		if remaining == 0 {
			// Start of a new character: classify by leading bits.
			switch {
			case b>>7 == 0b0: // 0xxxxxxx — 1-byte ASCII char
				remaining = 0
			case b>>5 == 0b110: // 110xxxxx — 2-byte leader
				remaining = 1
			case b>>4 == 0b1110: // 1110xxxx — 3-byte leader
				remaining = 2
			case b>>3 == 0b11110: // 11110xxx — 4-byte leader
				remaining = 3
			default:
				return false // 10xxxxxx as a starter, or 11111xxx: invalid leader
			}
		} else {
			// Mid-character: must be a 10xxxxxx continuation byte.
			if b>>6 != 0b10 {
				return false
			}
			remaining-- // one expected continuation byte consumed
		}
	}
	// A leftover expectation means the last character was truncated.
	return remaining == 0
}

// ── Approach 2: Bit-Pattern String Simulation (Readable) ─────────────────────
//
// validUtf8Strings validates by turning each byte into its 8-char binary string
// and scanning leading '1's. Slower and heavier than the mask version, but the
// logic reads almost like the spec, which is handy for explaining/debugging.
//
// Intuition:
//
//	The number of leading 1 bits of a leader byte tells you the character
//	length: 0 leading ones → 1-byte; k leading ones (k in 2..4) → k-byte char
//	needing k-1 continuation bytes. Continuation bytes have exactly the prefix
//	"10". We render bytes as strings to count leading ones directly.
//
// Algorithm:
//  1. For each byte, format its low 8 bits as an 8-character binary string.
//  2. Count leading '1's = ones.
//  3. If ones == 0 → 1-byte char (fine). If ones == 1 → stray continuation → invalid.
//     If ones in 2..4 → this leader needs ones-1 continuation bytes: check each
//     next byte begins with "10". If ones > 4 → invalid.
//  4. Advance the index past the whole character; if there aren't enough bytes
//     left, invalid.
//
// Time:  O(n) — each byte formatted/scanned a constant number of times.
// Space: O(1) extra beyond the fixed-size 8-char strings.
func validUtf8Strings(data []int) bool {
	n := len(data)
	i := 0
	for i < n {
		bits := byteToBits(data[i]) // 8-char binary string, e.g. "11000101"
		ones := leadingOnes(bits)   // number of leading '1' characters
		switch {
		case ones == 0:
			i++ // 1-byte char, move on
		case ones == 1 || ones > 4:
			return false // "10xxxxxx" as a leader, or too many leading ones
		default:
			// ones is 2..4: need ones-1 continuation bytes starting with "10".
			need := ones - 1
			if i+need >= n {
				return false // not enough bytes remain for this character
			}
			for k := 1; k <= need; k++ {
				cont := byteToBits(data[i+k])
				if cont[0] != '1' || cont[1] != '0' {
					return false // continuation byte must be 10xxxxxx
				}
			}
			i += 1 + need // skip leader + its continuation bytes
		}
	}
	return true
}

// byteToBits renders the low 8 bits of num as an 8-character binary string.
func byteToBits(num int) string {
	b := num & 0xFF
	out := make([]byte, 8)
	for k := 7; k >= 0; k-- {
		if b&1 == 1 {
			out[k] = '1'
		} else {
			out[k] = '0'
		}
		b >>= 1
	}
	return string(out)
}

// leadingOnes counts how many '1' characters appear at the start of s.
func leadingOnes(s string) int {
	count := 0
	for count < len(s) && s[count] == '1' {
		count++
	}
	return count
}

func main() {
	fmt.Println("=== Approach 1: Count Continuation Bytes (Bit Masks) ===")
	fmt.Printf("data=[197,130,1]: got=%v  expected true\n", validUtf8([]int{197, 130, 1}))              // expected true
	fmt.Printf("data=[235,140,4]: got=%v  expected false\n", validUtf8([]int{235, 140, 4}))             // expected false
	fmt.Printf("data=[240,162,138,147]: got=%v  expected true\n", validUtf8([]int{240, 162, 138, 147})) // expected true
	fmt.Printf("data=[255]: got=%v  expected false\n", validUtf8([]int{255}))                           // expected false

	fmt.Println("=== Approach 2: Bit-Pattern String Simulation ===")
	fmt.Printf("data=[197,130,1]: got=%v  expected true\n", validUtf8Strings([]int{197, 130, 1}))              // expected true
	fmt.Printf("data=[235,140,4]: got=%v  expected false\n", validUtf8Strings([]int{235, 140, 4}))             // expected false
	fmt.Printf("data=[240,162,138,147]: got=%v  expected true\n", validUtf8Strings([]int{240, 162, 138, 147})) // expected true
	fmt.Printf("data=[255]: got=%v  expected false\n", validUtf8Strings([]int{255}))                           // expected false
}
