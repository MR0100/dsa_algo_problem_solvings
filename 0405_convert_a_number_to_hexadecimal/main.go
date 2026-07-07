package main

import (
	"fmt"
	"strings"
)

// Given a 32-bit signed integer num, return its hexadecimal string. Negative
// numbers use two's complement (so -1 -> "ffffffff"). Output uses lowercase
// hex digits, no leading zeros (except num == 0 -> "0"). Built-in base
// conversions (e.g. strconv.FormatInt with radix 16) are not allowed.

const hexDigits = "0123456789abcdef" // maps a 0..15 nibble to its hex character

// ── Approach 1: Bit Masking Nibble-by-Nibble (Optimal) ───────────────────────
//
// bitMasking reads the number four bits at a time from the least-significant
// end, converting each 4-bit group (nibble) to one hex digit.
//
// Intuition:
//
//	Hex is base 16, so each hex digit encodes exactly 4 bits. Reinterpret num
//	as an unsigned 32-bit value (this IS two's complement for negatives), then
//	repeatedly take the low 4 bits with `& 0xf` to get a 0..15 nibble, map it to
//	a hex char, and shift right by 4 to expose the next nibble. Digits come out
//	least-significant first, so reverse (or prepend) at the end.
//
// Algorithm:
//  1. If num == 0, return "0".
//  2. n = uint32(num) — negatives become their two's-complement bit pattern.
//  3. While n != 0: digit = n & 0xf; prepend hexDigits[digit]; n >>= 4.
//  4. Return the assembled string.
//
// Time:  O(8) = O(1) — at most 8 nibbles in a 32-bit number.
// Space: O(8) = O(1) — the 8-character output buffer.
func bitMasking(num int) string {
	if num == 0 {
		return "0" // the only case that legitimately prints a single zero
	}
	// Reinterpreting the signed int as uint32 yields the two's-complement bits,
	// e.g. -1 -> 0xFFFFFFFF, which is exactly what the problem wants.
	n := uint32(num)

	var sb strings.Builder
	// Collect nibbles from least-significant to most-significant.
	var digits []byte
	for n != 0 {
		nibble := n & 0xf                          // isolate the low 4 bits (0..15)
		digits = append(digits, hexDigits[nibble]) // its hex character
		n >>= 4                                    // drop those 4 bits
	}
	// digits are reversed (low first); write them back high-to-low.
	for i := len(digits) - 1; i >= 0; i-- {
		sb.WriteByte(digits[i])
	}
	return sb.String()
}

// ── Approach 2: Fixed 8-Nibble Scan from the Top, Skipping Leading Zeros ──────
//
// topDownFixed inspects all 8 nibbles from most-significant to least, emitting
// characters once the first non-zero nibble is seen (to drop leading zeros).
//
// Intuition:
//
//	A 32-bit value is exactly 8 nibbles. Walk them high to low by shifting right
//	28, 24, ..., 0 bits and masking 0xf. Skip zero nibbles until the first
//	significant one appears; from then on append every nibble (interior zeros
//	must stay). This produces the digits already in the correct order — no
//	reversal needed.
//
// Algorithm:
//  1. If num == 0, return "0".
//  2. n = uint32(num). leading = true.
//  3. For shift = 28 down to 0 step 4: nibble = (n >> shift) & 0xf.
//     - If nibble != 0, leading = false.
//     - If not leading, append hexDigits[nibble].
//  4. Return the string (guaranteed non-empty since num != 0).
//
// Time:  O(8) = O(1) — a fixed 8 iterations.
// Space: O(8) = O(1) — output buffer.
func topDownFixed(num int) string {
	if num == 0 {
		return "0"
	}
	n := uint32(num)
	var sb strings.Builder
	leading := true // still skipping high-order zero nibbles?
	for shift := 28; shift >= 0; shift -= 4 {
		nibble := (n >> uint(shift)) & 0xf // the nibble at this position
		if nibble != 0 {
			leading = false // first significant nibble reached
		}
		if !leading {
			sb.WriteByte(hexDigits[nibble]) // emit real digits and interior zeros
		}
	}
	return sb.String()
}

// ── Approach 3: Repeated Division / Modulo on the Unsigned Value ──────────────
//
// divisionMod builds the hex string by taking n % 16 for each digit and n /= 16
// to advance, operating on the unsigned reinterpretation so negatives work.
//
// Intuition:
//
//	Base conversion by hand: the last hex digit is value mod 16, then divide by
//	16 and repeat. Doing this on uint32(num) handles negatives via two's
//	complement automatically. Digits emerge least-significant first, so reverse.
//	This avoids bitwise ops entirely, showing the arithmetic view of the same
//	computation (÷16 ≡ >>4, %16 ≡ &0xf).
//
// Algorithm:
//  1. If num == 0, return "0".
//  2. n = uint32(num). While n > 0: prepend hexDigits[n % 16]; n /= 16.
//  3. Return the assembled string.
//
// Time:  O(8) = O(1) — at most 8 divisions.
// Space: O(8) = O(1) — output buffer.
func divisionMod(num int) string {
	if num == 0 {
		return "0"
	}
	n := uint32(num) // two's-complement bit pattern viewed as an unsigned value

	var digits []byte
	for n > 0 {
		digits = append(digits, hexDigits[n%16]) // low-order hex digit
		n /= 16                                  // shift right one hex place
	}
	// Reverse: digits were produced least-significant first.
	var sb strings.Builder
	for i := len(digits) - 1; i >= 0; i-- {
		sb.WriteByte(digits[i])
	}
	return sb.String()
}

func main() {
	fmt.Println("=== Approach 1: Bit Masking (nibbles) ===")
	fmt.Printf("num=26 -> %q\n", bitMasking(26)) // expected "1a"
	fmt.Printf("num=-1 -> %q\n", bitMasking(-1)) // expected "ffffffff"
	fmt.Printf("num=0  -> %q\n", bitMasking(0))  // expected "0"

	fmt.Println("=== Approach 2: Top-Down Fixed 8-Nibble Scan ===")
	fmt.Printf("num=26 -> %q\n", topDownFixed(26)) // expected "1a"
	fmt.Printf("num=-1 -> %q\n", topDownFixed(-1)) // expected "ffffffff"
	fmt.Printf("num=0  -> %q\n", topDownFixed(0))  // expected "0"

	fmt.Println("=== Approach 3: Division / Modulo ===")
	fmt.Printf("num=26 -> %q\n", divisionMod(26)) // expected "1a"
	fmt.Printf("num=-1 -> %q\n", divisionMod(-1)) // expected "ffffffff"
	fmt.Printf("num=0  -> %q\n", divisionMod(0))  // expected "0"
}
