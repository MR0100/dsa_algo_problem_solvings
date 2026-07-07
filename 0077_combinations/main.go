package main

import "fmt"

// ── Approach 1: Backtracking ──────────────────────────────────────────────────
//
// backtracking solves Combinations by building k-combinations of [1..n].
//
// Intuition:
//   Choose numbers in increasing order to avoid duplicates. At each step,
//   choose a number from [start..n] and recurse. When the path has k elements,
//   record it.
//
//   Pruning: there's no point starting at `start` if remaining numbers (n-start+1)
//   are fewer than the numbers still needed (k - len(path)). This prunes many
//   branches early.
//
// Algorithm:
//   bt(start, path):
//     if len(path)==k: record; return
//     for i=start to n - (k-len(path)) + 1:  // pruning upper bound
//       bt(i+1, path+[i])
//
// Time:  O(C(n,k) × k) — C(n,k) combinations, each copied in O(k).
// Space: O(k) — recursion depth k.
func backtracking(n, k int) [][]int {
	var result [][]int
	var bt func(start int, path []int)
	bt = func(start int, path []int) {
		if len(path) == k {
			tmp := make([]int, k)
			copy(tmp, path)
			result = append(result, tmp)
			return
		}
		// pruning: need (k - len(path)) more elements; they come from [start..n]
		// so start must be <= n - (k - len(path)) + 1
		limit := n - (k - len(path)) + 1
		for i := start; i <= limit; i++ {
			bt(i+1, append(path, i))
		}
	}
	bt(1, nil)
	return result
}

// ── Approach 2: Iterative (Lexicographic) ────────────────────────────────────
//
// iterative solves Combinations iteratively by simulating the recursive
// backtracking using a pointer into the current combination.
//
// Intuition:
//   Start with [1,2,...,k]. At each step:
//   - Output the current combination.
//   - Find the rightmost element that can be incremented (nums[i] < n-k+1+i).
//   - Increment it and fill the rest sequentially.
//
// Time:  O(C(n,k) × k)
// Space: O(k) — current combination.
func iterative(n, k int) [][]int {
	var result [][]int
	nums := make([]int, k)
	for i := range nums {
		nums[i] = i + 1 // start with [1,2,...,k]
	}

	for {
		// record current combination
		tmp := make([]int, k)
		copy(tmp, nums)
		result = append(result, tmp)

		// find rightmost element that can be incremented
		i := k - 1
		for i >= 0 && nums[i] == n-k+i+1 {
			i--
		}
		if i < 0 {
			break // all elements at maximum; done
		}
		// increment nums[i] and fill rest sequentially
		nums[i]++
		for j := i + 1; j < k; j++ {
			nums[j] = nums[j-1] + 1
		}
	}
	return result
}

func main() {
	fmt.Println("=== Approach 1: Backtracking ===")
	r1 := backtracking(4, 2)
	fmt.Printf("n=4 k=2  count=%d  expected 6\n", len(r1))
	fmt.Println("combinations:", r1)

	r2 := backtracking(1, 1)
	fmt.Printf("n=1 k=1  count=%d  expected 1\n", len(r2))
	fmt.Println("combinations:", r2)

	fmt.Println("=== Approach 2: Iterative ===")
	r3 := iterative(4, 2)
	fmt.Printf("n=4 k=2  count=%d  expected 6\n", len(r3))
	fmt.Println("combinations:", r3)

	r4 := iterative(1, 1)
	fmt.Printf("n=1 k=1  count=%d  expected 1\n", len(r4))
	fmt.Println("combinations:", r4)
}
