package main

import "fmt"

// ── Approach 1: Memoized Recursion ────────────────────────────────────────────
//
// numDecodings solves Decode Ways using top-down DP.
//
// Intuition:
//   s[i:] can be decoded by taking 1 char (if '1'-'9') or 2 chars (if '10'-'26').
//   Let dp(i) = number of ways to decode s[i:].
//
// Time:  O(n)
// Space: O(n)
func numDecodings(s string) int {
	memo := make(map[int]int)
	var dp func(i int) int
	dp = func(i int) int {
		if i == len(s) {
			return 1 // successfully decoded all characters
		}
		if s[i] == '0' {
			return 0 // leading zero, can't decode
		}
		if v, ok := memo[i]; ok {
			return v
		}
		// take 1 digit
		result := dp(i + 1)
		// take 2 digits if valid (10..26)
		if i+1 < len(s) {
			twoDigit := (int(s[i]-'0'))*10 + int(s[i+1]-'0')
			if twoDigit >= 10 && twoDigit <= 26 {
				result += dp(i + 2)
			}
		}
		memo[i] = result
		return result
	}
	return dp(0)
}

// ── Approach 2: Bottom-Up DP ──────────────────────────────────────────────────
//
// numDecodingsDP solves Decode Ways using bottom-up DP.
//
// Intuition:
//   dp[i] = number of ways to decode s[0:i].
//   dp[0] = 1 (empty prefix — one way: do nothing).
//   dp[1] = 1 if s[0] != '0', else 0.
//   dp[i] = dp[i-1] if s[i-1] != '0'     (take 1 digit)
//          + dp[i-2] if s[i-2:i] in 10..26 (take 2 digits)
//
// Time:  O(n)
// Space: O(n), reducible to O(1).
func numDecodingsDP(s string) int {
	n := len(s)
	dp := make([]int, n+1)
	dp[0] = 1 // empty string has 1 decoding
	if s[0] == '0' {
		dp[1] = 0
	} else {
		dp[1] = 1
	}
	for i := 2; i <= n; i++ {
		// 1-digit decode
		if s[i-1] != '0' {
			dp[i] += dp[i-1]
		}
		// 2-digit decode
		twoDigit := (int(s[i-2]-'0'))*10 + int(s[i-1]-'0')
		if twoDigit >= 10 && twoDigit <= 26 {
			dp[i] += dp[i-2]
		}
	}
	return dp[n]
}

// ── Approach 3: O(1) Space DP ─────────────────────────────────────────────────
//
// numDecodingsO1 solves Decode Ways with O(1) space by keeping only two
// previous dp values (prev2 and prev1).
//
// Time:  O(n)
// Space: O(1)
func numDecodingsO1(s string) int {
	n := len(s)
	prev2 := 1 // dp[i-2]
	var prev1 int
	if s[0] == '0' {
		prev1 = 0
	} else {
		prev1 = 1
	}

	for i := 2; i <= n; i++ {
		curr := 0
		if s[i-1] != '0' {
			curr += prev1
		}
		twoDigit := (int(s[i-2]-'0'))*10 + int(s[i-1]-'0')
		if twoDigit >= 10 && twoDigit <= 26 {
			curr += prev2
		}
		prev2 = prev1
		prev1 = curr
	}
	return prev1
}

func main() {
	fmt.Println("=== Approach 1: Memoized Recursion ===")
	fmt.Printf("s=%q  got=%d  expected 2\n", "12", numDecodings("12"))
	fmt.Printf("s=%q  got=%d  expected 3\n", "226", numDecodings("226"))
	fmt.Printf("s=%q  got=%d  expected 0\n", "06", numDecodings("06"))

	fmt.Println("=== Approach 2: Bottom-Up DP ===")
	fmt.Printf("s=%q  got=%d  expected 2\n", "12", numDecodingsDP("12"))
	fmt.Printf("s=%q  got=%d  expected 3\n", "226", numDecodingsDP("226"))
	fmt.Printf("s=%q  got=%d  expected 0\n", "06", numDecodingsDP("06"))
	fmt.Printf("s=%q  got=%d  expected 1\n", "10", numDecodingsDP("10"))

	fmt.Println("=== Approach 3: O(1) Space ===")
	fmt.Printf("s=%q  got=%d  expected 2\n", "12", numDecodingsO1("12"))
	fmt.Printf("s=%q  got=%d  expected 3\n", "226", numDecodingsO1("226"))
	fmt.Printf("s=%q  got=%d  expected 0\n", "06", numDecodingsO1("06"))
	fmt.Printf("s=%q  got=%d  expected 1\n", "10", numDecodingsO1("10"))
}
