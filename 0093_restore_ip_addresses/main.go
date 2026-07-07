package main

import (
	"fmt"
	"strconv"
)

// ── Approach 1: Backtracking ──────────────────────────────────────────────────
//
// restoreIpAddresses solves Restore IP Addresses by trying all valid splits
// of the string into exactly 4 octets.
//
// Intuition:
//   An IP address has exactly 4 parts, each 1-3 digits, with value 0-255.
//   Leading zeros are only valid for the octet "0" itself (not "01" or "00").
//   Use backtracking: at each step, try taking 1, 2, or 3 characters as the
//   next octet. When we have 4 octets and have consumed all characters, record.
//
// Time:  O(1) — bounded by the fixed max branches (3 choices × 4 levels = 81).
// Space: O(1) — recursion depth is exactly 4.
func restoreIpAddresses(s string) []string {
	var result []string
	var bt func(start, parts int, current string)
	bt = func(start, parts int, current string) {
		if parts == 4 && start == len(s) {
			result = append(result, current[:len(current)-1]) // remove trailing dot
			return
		}
		if parts == 4 || start == len(s) {
			return
		}
		// try 1, 2, 3 digits for the next octet
		for length := 1; length <= 3; length++ {
			if start+length > len(s) {
				break
			}
			segment := s[start : start+length]
			// no leading zeros (except "0" itself)
			if length > 1 && segment[0] == '0' {
				break
			}
			// value must be <= 255
			val, _ := strconv.Atoi(segment)
			if val > 255 {
				break
			}
			bt(start+length, parts+1, current+segment+".")
		}
	}
	bt(0, 0, "")
	return result
}

// ── Approach 2: Three Nested Loops ───────────────────────────────────────────
//
// restoreIpAddressesIter solves Restore IP Addresses using three nested loops
// to place 3 dots and check all splits.
//
// Intuition:
//   Place dots at positions i, j, k in the string s. The four segments are
//   s[0:i], s[i:j], s[j:k], s[k:]. Check each segment is a valid octet.
//
// Time:  O(1) — O(n³) loops but n≤12, so bounded.
// Space: O(1)
func restoreIpAddressesIter(s string) []string {
	var result []string
	n := len(s)

	isValid := func(seg string) bool {
		if len(seg) == 0 || len(seg) > 3 {
			return false
		}
		if len(seg) > 1 && seg[0] == '0' {
			return false
		}
		v, _ := strconv.Atoi(seg)
		return v <= 255
	}

	for i := 1; i <= 3 && i < n; i++ {
		for j := i + 1; j <= i+3 && j < n; j++ {
			for k := j + 1; k <= j+3 && k < n; k++ {
				a, b, c, d := s[:i], s[i:j], s[j:k], s[k:]
				if isValid(a) && isValid(b) && isValid(c) && isValid(d) {
					result = append(result, a+"."+b+"."+c+"."+d)
				}
			}
		}
	}
	return result
}

func main() {
	fmt.Println("=== Approach 1: Backtracking ===")
	fmt.Printf("s=%q  got=%v\n", "25525511135", restoreIpAddresses("25525511135"))
	fmt.Printf("s=%q  got=%v\n", "0000", restoreIpAddresses("0000"))
	fmt.Printf("s=%q  got=%v\n", "101023", restoreIpAddresses("101023"))

	fmt.Println("=== Approach 2: Three Nested Loops ===")
	fmt.Printf("s=%q  got=%v\n", "25525511135", restoreIpAddressesIter("25525511135"))
	fmt.Printf("s=%q  got=%v\n", "0000", restoreIpAddressesIter("0000"))
	fmt.Printf("s=%q  got=%v\n", "101023", restoreIpAddressesIter("101023"))
}
