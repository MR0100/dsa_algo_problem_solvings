package main

import "fmt"

// ── Approach 1: Brute Force (Recursive Subsequence Enumeration) ───────────────
//
// bruteForce solves Arithmetic Slices II - Subsequence by explicitly extending
// every arithmetic subsequence one element at a time via recursion.
//
// Intuition:
//
//	An arithmetic subsequence is fully described by (its last index, its common
//	difference). Anchor a pair (i, j) with i < j as the first two terms; that
//	fixes the difference d = nums[j] - nums[i]. From index j we then try to
//	append any later index k whose value equals nums[j] + d, and recurse. Every
//	time we successfully append a THIRD-or-later element we have completed a
//	valid arithmetic subsequence (length >= 3), so we count it. This walks the
//	entire space of arithmetic subsequences, hence it is correct but exponential.
//
// Algorithm:
//  1. For every ordered starting pair (i, j), i < j, let d = nums[j]-nums[i].
//  2. extend(j, d): scan k > j; if nums[k]-nums[j] == d, this append makes a
//     subsequence of length >= 3 → count++, then recurse extend(k, d).
//  3. Sum all counts.
//
// Time:  O(2^n) worst case — e.g. all-equal arrays branch into every subset of
//
//	size >= 3; each arithmetic subsequence is generated once.
//
// Space: O(n) — recursion depth is bounded by the array length.
//
// Note: int64 differences avoid overflow because nums[i] spans the full int32
// range and a difference of two int32 values can exceed int32.
func bruteForce(nums []int) int {
	n := len(nums)
	count := 0 // total valid arithmetic subsequences (length >= 3)

	// extend continues an arithmetic subsequence whose last element sits at
	// index last and whose common difference is d. Any successful append here
	// is the 3rd (or later) element, so it forms a valid slice.
	var extend func(last int, d int64)
	extend = func(last int, d int64) {
		for k := last + 1; k < n; k++ {
			// Does appending nums[k] preserve the common difference d?
			if int64(nums[k])-int64(nums[last]) == d {
				count++      // length >= 3 reached → a valid subsequence
				extend(k, d) // try to grow it further from k
			}
		}
	}

	// Every pair (i, j) seeds a difference; recursion supplies the >=3rd term.
	for i := 0; i < n; i++ {
		for j := i + 1; j < n; j++ {
			d := int64(nums[j]) - int64(nums[i]) // fixed common difference
			extend(j, d)                         // count all extensions
		}
	}
	return count
}

// ── Approach 2: Dynamic Programming with Hash Maps (Optimal) ──────────────────
//
// dpHashMap solves Arithmetic Slices II - Subsequence by counting, for each
// index i and each difference d, how many WEAK arithmetic subsequences of
// length >= 2 end at i with difference d.
//
// Intuition:
//
//	Define dp[i][d] = number of arithmetic subsequences (allowing length 2, the
//	"weak" ones) that END at index i with common difference d. When we look at a
//	pair (i, j) with j < i and difference d = nums[i]-nums[j], every weak
//	subsequence ending at j with the SAME d can be extended by nums[i]. Each of
//	those had length >= 2, so appending nums[i] makes length >= 3 — a real
//	answer. Thus we add dp[j][d] to the global answer, and we also grow
//	dp[i][d] by dp[j][d] + 1 (the +1 is the brand-new length-2 pair (j, i)).
//	Summing the "promotions" over all pairs counts every arithmetic subsequence
//	exactly once, because a subsequence of length L is counted at the moment its
//	last element is appended, from the state stored at its second-to-last.
//
// Algorithm:
//  1. dp[i] is a map[int64]int: difference → count of weak subsequences ending
//     at i with that difference.
//  2. For each i, for each j < i: d = nums[i]-nums[j];
//     cnt = dp[j][d] (weak subsequences ending at j with diff d);
//     ans += cnt  (each becomes a valid length>=3 subsequence via nums[i]);
//     dp[i][d] += cnt + 1  (extend them, plus the new pair (j,i)).
//  3. Return ans.
//
// Time:  O(n^2) — every ordered pair (j, i) is processed once with O(1)
//
//	amortised map work.
//
// Space: O(n^2) — up to O(n) distinct differences stored per index.
//
// Note: differences are int64 to survive int32-range subtraction; the map key
// type is int64 for the same reason.
func dpHashMap(nums []int) int {
	n := len(nums)
	// dp[i][d] = # of weak (len>=2) arithmetic subsequences ending at i, diff d.
	dp := make([]map[int64]int, n)
	for i := range dp {
		dp[i] = make(map[int64]int)
	}

	ans := 0 // count of STRONG (len>=3) arithmetic subsequences
	for i := 0; i < n; i++ {
		for j := 0; j < i; j++ {
			// Common difference contributed by ending pair (j, i).
			d := int64(nums[i]) - int64(nums[j])
			// Weak subsequences ending at j with this exact difference:
			// each already has length >= 2, so nums[i] promotes it to >= 3.
			cnt := dp[j][d]
			ans += cnt // these are genuine arithmetic subsequences now
			// Extend those weak ones to end at i, and add 1 for the fresh
			// length-2 pair (j, i) that starts a new arithmetic run.
			dp[i][d] += cnt + 1
		}
	}
	return ans
}

func main() {
	fmt.Println("=== Approach 1: Brute Force (Recursive Enumeration) ===")
	fmt.Printf("nums=[2,4,6,8,10]  got=%d  expected 7\n", bruteForce([]int{2, 4, 6, 8, 10}))
	fmt.Printf("nums=[7,7,7,7,7]   got=%d  expected 16\n", bruteForce([]int{7, 7, 7, 7, 7}))
	fmt.Printf("nums=[2,2,3,4]     got=%d  expected 2\n", bruteForce([]int{2, 2, 3, 4})) // [2,3,4] via either 2

	fmt.Println("=== Approach 2: DP with Hash Maps (Optimal) ===")
	fmt.Printf("nums=[2,4,6,8,10]  got=%d  expected 7\n", dpHashMap([]int{2, 4, 6, 8, 10}))
	fmt.Printf("nums=[7,7,7,7,7]   got=%d  expected 16\n", dpHashMap([]int{7, 7, 7, 7, 7}))
	fmt.Printf("nums=[2,2,3,4]     got=%d  expected 2\n", dpHashMap([]int{2, 2, 3, 4}))
}
