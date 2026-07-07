package main

import "fmt"

// ── Approach 1: Two Pointers with Map (Optimal) ──────────────────────────────
//
// twoPointers decides whether num is strobogrammatic — i.e. it looks the same
// after being rotated 180 degrees.
//
// Intuition:
//
//	A number reads the same upside-down only if, reading it from BOTH ends
//	simultaneously, each outer pair rotates into the other. Rotating 180°
//	reverses the digit order, so the first digit must be the rotation of the
//	last digit, the second the rotation of the second-to-last, and so on.
//	Only these rotations are legal: 0→0, 1→1, 8→8, 6→9, 9→6. Any digit
//	outside {0,1,6,8,9} instantly disqualifies the number (2,3,4,5,7 have no
//	valid upside-down form).
//
// Algorithm:
//  1. Build a map rotate: digit -> the digit it becomes when flipped.
//  2. Set left = 0, right = len(num)-1.
//  3. While left <= right:
//     a. If num[left] is not a rotatable digit, return false.
//     b. If rotate[num[left]] != num[right], return false (pair mismatch).
//     c. Move left++ and right--.
//  4. If every pair matched, return true.
//
// Time:  O(n) — one pass over half the string.
// Space: O(1) — the rotate map has a fixed 5 entries.
func twoPointers(num string) bool {
	// The only digits that have a valid 180° rotation, mapped to their image.
	rotate := map[byte]byte{
		'0': '0',
		'1': '1',
		'8': '8',
		'6': '9', // 6 upside-down looks like 9
		'9': '6', // 9 upside-down looks like 6
	}

	left, right := 0, len(num)-1 // converge from both ends
	for left <= right {
		// The left digit must itself be rotatable; if not, fail immediately.
		mirror, ok := rotate[num[left]]
		if !ok {
			return false
		}
		// The rotation of the left digit must equal the right digit, because
		// flipping reverses order and swaps the two positions.
		if mirror != num[right] {
			return false
		}
		left++  // advance the left pointer inward
		right-- // advance the right pointer inward
	}
	return true // every outer pair rotated into its partner
}

// ── Approach 2: Build Rotated String and Compare ─────────────────────────────
//
// buildAndCompare constructs the full 180°-rotated version of num and checks it
// against the original.
//
// Intuition:
//
//	Rotating a number 180° means: rotate each digit individually, then reverse
//	the whole sequence (because the last digit ends up first). If any digit is
//	not rotatable the number cannot be strobogrammatic. Building that rotated
//	string explicitly and comparing to the input is the most literal reading of
//	the definition — clearer, though it allocates O(n) extra space.
//
// Algorithm:
//  1. Walk num from the LAST digit to the FIRST.
//  2. For each digit, look up its rotation; if it has none, return false.
//  3. Append the rotation to a builder (last digit first ⇒ builds the reversal).
//  4. Compare the built string to num; equal ⇒ strobogrammatic.
//
// Time:  O(n) — one reverse pass to build + one comparison.
// Space: O(n) — the rotated string buffer.
func buildAndCompare(num string) bool {
	rotate := map[byte]byte{
		'0': '0', '1': '1', '8': '8', '6': '9', '9': '6',
	}

	rotated := make([]byte, 0, len(num)) // buffer for the flipped number
	// Iterate from the end so appending naturally reverses the order.
	for i := len(num) - 1; i >= 0; i-- {
		mirror, ok := rotate[num[i]] // rotation of this digit
		if !ok {
			return false // a non-rotatable digit ⇒ never strobogrammatic
		}
		rotated = append(rotated, mirror) // append rotated digit (reversed order)
	}
	// The number is strobogrammatic iff its 180° rotation equals itself.
	return string(rotated) == num
}

func main() {
	fmt.Println("=== Approach 1: Two Pointers with Map (Optimal) ===")
	fmt.Printf("num=\"69\"   got=%v  expected true\n", twoPointers("69"))   // expected true
	fmt.Printf("num=\"88\"   got=%v  expected true\n", twoPointers("88"))   // expected true
	fmt.Printf("num=\"962\"  got=%v  expected false\n", twoPointers("962")) // expected false
	fmt.Printf("num=\"1\"    got=%v  expected true\n", twoPointers("1"))    // expected true

	fmt.Println("=== Approach 2: Build Rotated String and Compare ===")
	fmt.Printf("num=\"69\"   got=%v  expected true\n", buildAndCompare("69"))   // expected true
	fmt.Printf("num=\"88\"   got=%v  expected true\n", buildAndCompare("88"))   // expected true
	fmt.Printf("num=\"962\"  got=%v  expected false\n", buildAndCompare("962")) // expected false
	fmt.Printf("num=\"1\"    got=%v  expected true\n", buildAndCompare("1"))    // expected true
}
