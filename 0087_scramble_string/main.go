package main

import "fmt"

// ── Approach 1: Memoized Recursion (Top-Down DP) ──────────────────────────────
//
// isScramble solves Scramble String using memoized recursion.
//
// Intuition:
//   s1 is a scramble of s2 iff we can split s1 at some index k and either:
//   (a) [s1[:k] is scramble of s2[:k]] AND [s1[k:] is scramble of s2[k:]]  (no swap), OR
//   (b) [s1[:k] is scramble of s2[n-k:]] AND [s1[k:] is scramble of s2[:n-k]]  (swap)
//
//   Base case: s1 == s2 (identical strings are trivially scrambles).
//   Early exit: if sorted characters differ, they can't be scrambles.
//
// Time:  O(n^4) — O(n^3) states (two indices + length), O(n) per state to check splits.
// Space: O(n^3) — memoization map.
func isScramble(s1 string, s2 string) bool {
	memo := make(map[string]bool)

	var dp func(a, b string) bool
	dp = func(a, b string) bool {
		if a == b {
			return true
		}
		if len(a) != len(b) {
			return false
		}
		key := a + "#" + b
		if v, ok := memo[key]; ok {
			return v
		}

		// pruning: if character frequencies differ, can't be scrambles
		freq := [26]int{}
		for i := 0; i < len(a); i++ {
			freq[a[i]-'a']++
			freq[b[i]-'a']--
		}
		for _, f := range freq {
			if f != 0 {
				memo[key] = false
				return false
			}
		}

		n := len(a)
		for k := 1; k < n; k++ {
			// no swap: a[:k] ~ b[:k] && a[k:] ~ b[k:]
			if dp(a[:k], b[:k]) && dp(a[k:], b[k:]) {
				memo[key] = true
				return true
			}
			// swap: a[:k] ~ b[n-k:] && a[k:] ~ b[:n-k]
			if dp(a[:k], b[n-k:]) && dp(a[k:], b[:n-k]) {
				memo[key] = true
				return true
			}
		}
		memo[key] = false
		return false
	}

	return dp(s1, s2)
}

// ── Approach 2: 3D DP (Bottom-Up) ─────────────────────────────────────────────
//
// isScrambleDP solves Scramble String using 3D bottom-up DP.
//
// Intuition:
//   dp[len][i][j] = true if s1[i:i+len] is a scramble of s2[j:j+len].
//   Fill from smaller lengths to larger.
//
//   dp[len][i][j] = OR over k=1..len-1 of:
//     (dp[k][i][j] && dp[len-k][i+k][j+k])       no-swap
//     (dp[k][i][j+len-k] && dp[len-k][i+k][j])   swap
//
// Time:  O(n^4) — O(n^3) states, O(n) transitions per state.
// Space: O(n^3)
func isScrambleDP(s1 string, s2 string) bool {
	n := len(s1)
	if n != len(s2) {
		return false
	}
	// dp[length][i][j]: s1[i:i+length] is scramble of s2[j:j+length]
	dp := make([][][]bool, n+1)
	for l := 0; l <= n; l++ {
		dp[l] = make([][]bool, n)
		for i := range dp[l] {
			dp[l][i] = make([]bool, n)
		}
	}
	// base case: length 1
	for i := 0; i < n; i++ {
		for j := 0; j < n; j++ {
			dp[1][i][j] = s1[i] == s2[j]
		}
	}
	// fill for length 2..n
	for length := 2; length <= n; length++ {
		for i := 0; i <= n-length; i++ {
			for j := 0; j <= n-length; j++ {
				for k := 1; k < length; k++ {
					// no-swap
					if dp[k][i][j] && dp[length-k][i+k][j+k] {
						dp[length][i][j] = true
						break
					}
					// swap
					if dp[k][i][j+length-k] && dp[length-k][i+k][j] {
						dp[length][i][j] = true
						break
					}
				}
			}
		}
	}
	return dp[n][0][0]
}

func main() {
	fmt.Println("=== Approach 1: Memoized Recursion ===")
	fmt.Printf("s1=%q s2=%q  got=%v  expected true\n", "great", "rgeat", isScramble("great", "rgeat"))
	fmt.Printf("s1=%q s2=%q  got=%v  expected false\n", "abcde", "caebd", isScramble("abcde", "caebd"))
	fmt.Printf("s1=%q s2=%q  got=%v  expected true\n", "a", "a", isScramble("a", "a"))

	fmt.Println("=== Approach 2: 3D Bottom-Up DP ===")
	fmt.Printf("s1=%q s2=%q  got=%v  expected true\n", "great", "rgeat", isScrambleDP("great", "rgeat"))
	fmt.Printf("s1=%q s2=%q  got=%v  expected false\n", "abcde", "caebd", isScrambleDP("abcde", "caebd"))
	fmt.Printf("s1=%q s2=%q  got=%v  expected true\n", "a", "a", isScrambleDP("a", "a"))
}
