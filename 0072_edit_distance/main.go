package main

import "fmt"

// ── Approach 1: Recursion with Memoization ────────────────────────────────────
//
// memoization solves Edit Distance using top-down DP.
//
// Intuition:
//   dp(i,j) = min edit ops to convert word1[i:] to word2[j:].
//   Base: dp(i, len2) = len2 - i is wrong; actually:
//         dp(i, len2) = len1-i (delete remaining chars of word1)
//         dp(len1, j) = len2-j (insert remaining chars of word2)
//   Recurrence:
//     if word1[i]==word2[j]: dp(i,j) = dp(i+1,j+1)
//     else: dp(i,j) = 1 + min(dp(i+1,j),   // delete word1[i]
//                              dp(i,j+1),   // insert word2[j]
//                              dp(i+1,j+1)) // replace word1[i] with word2[j]
//
// Time:  O(m × n)
// Space: O(m × n)
func memoization(word1 string, word2 string) int {
	m, n := len(word1), len(word2)
	memo := make([][]int, m+1)
	for i := range memo {
		memo[i] = make([]int, n+1)
		for j := range memo[i] {
			memo[i][j] = -1
		}
	}

	var dp func(i, j int) int
	dp = func(i, j int) int {
		if i == m {
			return n - j // insert remaining chars of word2
		}
		if j == n {
			return m - i // delete remaining chars of word1
		}
		if memo[i][j] != -1 {
			return memo[i][j]
		}
		if word1[i] == word2[j] {
			memo[i][j] = dp(i+1, j+1)
		} else {
			del := dp(i+1, j)   // delete word1[i]
			ins := dp(i, j+1)   // insert word2[j] before word1[i]
			rep := dp(i+1, j+1) // replace word1[i] with word2[j]
			best := del
			if ins < best {
				best = ins
			}
			if rep < best {
				best = rep
			}
			memo[i][j] = 1 + best
		}
		return memo[i][j]
	}

	return dp(0, 0)
}

// ── Approach 2: DP Bottom-Up (2D Table) ──────────────────────────────────────
//
// dpBottomUp solves Edit Distance using the classic DP table.
//
// Intuition:
//   dp[i][j] = min ops to convert word1[0..i-1] to word2[0..j-1].
//   dp[i][0] = i (delete all chars of word1[:i])
//   dp[0][j] = j (insert all chars of word2[:j])
//   dp[i][j] = dp[i-1][j-1] if word1[i-1]==word2[j-1]
//            = 1 + min(dp[i-1][j], dp[i][j-1], dp[i-1][j-1]) otherwise
//
// Time:  O(m × n)
// Space: O(m × n)
func dpBottomUp(word1, word2 string) int {
	m, n := len(word1), len(word2)
	dp := make([][]int, m+1)
	for i := range dp {
		dp[i] = make([]int, n+1)
		dp[i][0] = i // cost to delete i chars
	}
	for j := 0; j <= n; j++ {
		dp[0][j] = j // cost to insert j chars
	}
	for i := 1; i <= m; i++ {
		for j := 1; j <= n; j++ {
			if word1[i-1] == word2[j-1] {
				dp[i][j] = dp[i-1][j-1] // no op needed
			} else {
				best := dp[i-1][j] // delete
				if dp[i][j-1] < best {
					best = dp[i][j-1] // insert
				}
				if dp[i-1][j-1] < best {
					best = dp[i-1][j-1] // replace
				}
				dp[i][j] = 1 + best
			}
		}
	}
	return dp[m][n]
}

// ── Approach 3: DP Rolling Row (O(n) Space) ───────────────────────────────────
//
// dpRolling solves Edit Distance with O(n) space by maintaining only two rows.
//
// Intuition:
//   dp[i][j] depends on dp[i-1][j-1], dp[i-1][j], dp[i][j-1].
//   Keep prev (previous row) and curr (current row), updating curr in-place.
//
// Time:  O(m × n)
// Space: O(n)
func dpRolling(word1, word2 string) int {
	m, n := len(word1), len(word2)
	prev := make([]int, n+1)
	for j := 0; j <= n; j++ {
		prev[j] = j // dp[0][j] = j
	}
	for i := 1; i <= m; i++ {
		curr := make([]int, n+1)
		curr[0] = i // dp[i][0] = i
		for j := 1; j <= n; j++ {
			if word1[i-1] == word2[j-1] {
				curr[j] = prev[j-1]
			} else {
				best := prev[j] // delete
				if curr[j-1] < best {
					best = curr[j-1] // insert
				}
				if prev[j-1] < best {
					best = prev[j-1] // replace
				}
				curr[j] = 1 + best
			}
		}
		prev = curr
	}
	return prev[n]
}

func main() {
	fmt.Println("=== Approach 1: Memoization ===")
	fmt.Printf("horse→ros  got=%d  expected 3\n", memoization("horse", "ros"))
	fmt.Printf("intention→execution  got=%d  expected 5\n", memoization("intention", "execution"))
	fmt.Printf("\"\"→\"\"  got=%d  expected 0\n", memoization("", ""))
	fmt.Printf("a→\"\"  got=%d  expected 1\n", memoization("a", ""))

	fmt.Println("=== Approach 2: DP 2D Table ===")
	fmt.Printf("horse→ros  got=%d  expected 3\n", dpBottomUp("horse", "ros"))
	fmt.Printf("intention→execution  got=%d  expected 5\n", dpBottomUp("intention", "execution"))
	fmt.Printf("\"\"→\"\"  got=%d  expected 0\n", dpBottomUp("", ""))
	fmt.Printf("a→\"\"  got=%d  expected 1\n", dpBottomUp("a", ""))

	fmt.Println("=== Approach 3: DP Rolling Row ===")
	fmt.Printf("horse→ros  got=%d  expected 3\n", dpRolling("horse", "ros"))
	fmt.Printf("intention→execution  got=%d  expected 5\n", dpRolling("intention", "execution"))
	fmt.Printf("\"\"→\"\"  got=%d  expected 0\n", dpRolling("", ""))
	fmt.Printf("a→\"\"  got=%d  expected 1\n", dpRolling("a", ""))
}
