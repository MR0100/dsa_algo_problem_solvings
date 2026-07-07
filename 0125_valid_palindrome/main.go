package main

import (
	"fmt"
	"strings"
	"unicode"
)

// ── Approach 1: Clean String then Two Pointers ────────────────────────────────
//
// isPalindrome solves Valid Palindrome by filtering then checking.
//
// Intuition:
//   Keep only alphanumeric characters, lowercase. Use two-pointer palindrome check.
//
// Time:  O(n)
// Space: O(n) — filtered string.
func isPalindrome(s string) bool {
	var sb strings.Builder
	for _, ch := range s {
		if unicode.IsLetter(ch) || unicode.IsDigit(ch) {
			sb.WriteRune(unicode.ToLower(ch))
		}
	}
	filtered := sb.String()
	for i, j := 0, len(filtered)-1; i < j; i, j = i+1, j-1 {
		if filtered[i] != filtered[j] {
			return false
		}
	}
	return true
}

// ── Approach 2: Two Pointers In-Place (O(1) Space) ───────────────────────────
//
// isPalindromeO1 solves Valid Palindrome without extra string allocation.
//
// Intuition:
//   Use two pointers l, r on the original string. Skip non-alphanumeric characters.
//   Compare lowercase versions of the characters they point to.
//
// Time:  O(n)
// Space: O(1)
func isPalindromeO1(s string) bool {
	l, r := 0, len(s)-1

	for l < r {
		// skip non-alphanumeric from left
		for l < r && !isAlphanumeric(s[l]) {
			l++
		}
		// skip non-alphanumeric from right
		for l < r && !isAlphanumeric(s[r]) {
			r--
		}
		if toLower(s[l]) != toLower(s[r]) {
			return false
		}
		l++
		r--
	}
	return true
}

func isAlphanumeric(b byte) bool {
	return (b >= 'a' && b <= 'z') || (b >= 'A' && b <= 'Z') || (b >= '0' && b <= '9')
}

func toLower(b byte) byte {
	if b >= 'A' && b <= 'Z' {
		return b + 32
	}
	return b
}

func main() {
	fmt.Println("=== Approach 1: Filter + Two Pointers ===")
	fmt.Printf(`s="A man, a plan, a canal: Panama"  got=%v  expected true`+"\n", isPalindrome("A man, a plan, a canal: Panama"))
	fmt.Printf(`s="race a car"  got=%v  expected false`+"\n", isPalindrome("race a car"))
	fmt.Printf(`s=" "  got=%v  expected true`+"\n", isPalindrome(" "))

	fmt.Println("=== Approach 2: In-Place Two Pointers ===")
	fmt.Printf(`s="A man, a plan, a canal: Panama"  got=%v  expected true`+"\n", isPalindromeO1("A man, a plan, a canal: Panama"))
	fmt.Printf(`s="race a car"  got=%v  expected false`+"\n", isPalindromeO1("race a car"))
	fmt.Printf(`s=" "  got=%v  expected true`+"\n", isPalindromeO1(" "))
}
