package main

import "fmt"

// ── Approach 1: DP (Catalan Number) ──────────────────────────────────────────
//
// numTrees solves Unique Binary Search Trees using DP.
//
// Intuition:
//   Let dp[n] = number of structurally unique BSTs with n nodes.
//   For each root i (1..n), left subtree has (i-1) nodes and right has (n-i) nodes.
//   dp[n] = sum over i=1..n of dp[i-1] * dp[n-i].
//   dp[0] = 1 (empty tree), dp[1] = 1.
//   This produces the Catalan numbers: 1, 1, 2, 5, 14, 42, 132, ...
//
// Time:  O(n²)
// Space: O(n)
func numTrees(n int) int {
	dp := make([]int, n+1)
	dp[0] = 1
	dp[1] = 1
	for i := 2; i <= n; i++ {
		for j := 1; j <= i; j++ {
			dp[i] += dp[j-1] * dp[i-j] // j is root: left has j-1 nodes, right has i-j nodes
		}
	}
	return dp[n]
}

// ── Approach 2: Catalan Number Formula ───────────────────────────────────────
//
// numTreesFormula solves Unique Binary Search Trees using the closed-form
// Catalan number formula: C(n) = C(2n, n) / (n+1).
//
// Compute iteratively to avoid overflow: C(2n,n)/(n+1) = product over i=0..n-1 of (n+1+i)/(i+1).
//
// Time:  O(n)
// Space: O(1)
func numTreesFormula(n int) int {
	// C(2n, n) / (n+1) computed iteratively
	// C_n = (2n)! / (n! * (n+1)!)
	result := 1
	for i := 0; i < n; i++ {
		// multiply by (n+1+i) and divide by (i+1)
		// division is exact at each step (Catalan numbers are always integers)
		result = result * (n + 1 + i) / (i + 1)
	}
	return result / (n + 1)
}

func main() {
	fmt.Println("=== Approach 1: DP ===")
	fmt.Printf("n=3  got=%d  expected 5\n", numTrees(3))
	fmt.Printf("n=1  got=%d  expected 1\n", numTrees(1))
	fmt.Printf("n=5  got=%d  expected 42\n", numTrees(5))
	fmt.Printf("n=19  got=%d  expected 1767263190\n", numTrees(19))

	fmt.Println("=== Approach 2: Catalan Formula ===")
	fmt.Printf("n=3  got=%d  expected 5\n", numTreesFormula(3))
	fmt.Printf("n=1  got=%d  expected 1\n", numTreesFormula(1))
	fmt.Printf("n=5  got=%d  expected 42\n", numTreesFormula(5))
}
