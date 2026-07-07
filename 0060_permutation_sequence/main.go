package main

import (
	"fmt"
	"strconv"
	"strings"
)

// ── Approach 1: Generate All Permutations ────────────────────────────────────
//
// generateAll solves Permutation Sequence by generating all permutations in
// lexicographic order and returning the k-th one.
//
// Intuition:
//   Generate all n! permutations via backtracking, collect them in order, return k-th.
//   Very slow for large n (n=9 → 362880 permutations).
//
// Time:  O(n! × n) — n! permutations, each of length n.
// Space: O(n! × n)
func generateAll(n, k int) string {
	nums := make([]int, n)
	for i := range nums {
		nums[i] = i + 1
	}
	var result []string
	visited := make([]bool, n)

	var bt func(path []int)
	bt = func(path []int) {
		if len(path) == n {
			parts := make([]string, n)
			for i, v := range path {
				parts[i] = strconv.Itoa(v)
			}
			result = append(result, strings.Join(parts, ""))
			return
		}
		for i := 0; i < n; i++ {
			if !visited[i] {
				visited[i] = true
				bt(append(path, nums[i]))
				visited[i] = false
			}
		}
	}

	bt(nil)
	return result[k-1]
}

// ── Approach 2: Factorial Number System (Optimal) ────────────────────────────
//
// factorialSystem solves Permutation Sequence directly in O(n²) time.
//
// Intuition:
//   The n! permutations of n digits are divided into n groups of (n-1)! each,
//   one per leading digit. The group is determined by (k-1) / (n-1)!.
//   After choosing the leading digit, recurse on the remaining n-1 digits with
//   the updated k.
//
//   This is the "factorial number system" representation:
//   k-1 = d[0]*(n-1)! + d[1]*(n-2)! + ... + d[n-1]*0!
//   d[i] is the index into the remaining digits at step i.
//
// Algorithm:
//   digits = [1,2,...,n]; fact = [0!,1!,...,(n-1)!]; k -= 1 (0-indexed)
//   for i = n downto 1:
//     idx = k / fact[i-1]
//     result += digits[idx]
//     digits = remove digits[idx]
//     k %= fact[i-1]
//
// Time:  O(n²) — n steps; each removal from digits list is O(n).
// Space: O(n)  — digits list.
func factorialSystem(n, k int) string {
	// precompute factorials 0! to (n-1)!
	fact := make([]int, n)
	fact[0] = 1
	for i := 1; i < n; i++ {
		fact[i] = fact[i-1] * i
	}

	// available digits
	digits := make([]int, n)
	for i := range digits {
		digits[i] = i + 1
	}

	k-- // convert to 0-indexed
	var sb strings.Builder

	for i := n; i >= 1; i-- {
		idx := k / fact[i-1]   // which digit to pick
		sb.WriteString(strconv.Itoa(digits[idx]))
		digits = append(digits[:idx], digits[idx+1:]...) // remove chosen digit
		k %= fact[i-1]         // remainder for next step
	}

	return sb.String()
}

func main() {
	fmt.Println("=== Approach 1: Generate All Permutations ===")
	fmt.Printf("n=3 k=3  got=%s  expected 213\n", generateAll(3, 3))
	fmt.Printf("n=4 k=9  got=%s  expected 2314\n", generateAll(4, 9))
	fmt.Printf("n=3 k=1  got=%s  expected 123\n", generateAll(3, 1))

	fmt.Println("=== Approach 2: Factorial Number System (Optimal) ===")
	fmt.Printf("n=3 k=3  got=%s  expected 213\n", factorialSystem(3, 3))
	fmt.Printf("n=4 k=9  got=%s  expected 2314\n", factorialSystem(4, 9))
	fmt.Printf("n=3 k=1  got=%s  expected 123\n", factorialSystem(3, 1))
}
