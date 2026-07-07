package main

import (
	"fmt"
	"sort"
	"strconv"
)

// ── Approach 1: Brute Force (Generate + Sort Lexicographically) ──────────────
//
// bruteForce solves K-th Smallest in Lexicographical Order by materialising all
// numbers 1..n, sorting them as strings, and indexing the (k-1)-th.
//
// Intuition:
//
//	"Lexicographical order" is just string order. So write out 1..n, sort by
//	their decimal spelling, and pick element k. This is the definition made
//	literal. It is only viable for small n — n can be 10^9, where an O(n)
//	array of strings is far too big — but it is the ground truth the optimal
//	method must reproduce.
//
// Algorithm:
//  1. Build the slice [1, 2, ..., n].
//  2. Sort it with a comparator that compares decimal string representations.
//  3. Return the element at index k-1.
//
// Time:  O(n log n · L) — sorting n items, each string compare up to L digits.
// Space: O(n) — the materialised list (and its string keys during compare).
func bruteForce(n int, k int) int {
	nums := make([]int, n) // all candidates 1..n
	for i := 0; i < n; i++ {
		nums[i] = i + 1
	}
	// Compare by decimal spelling so "10" < "2" (lexicographic, not numeric).
	sort.Slice(nums, func(a, b int) bool {
		return strconv.Itoa(nums[a]) < strconv.Itoa(nums[b])
	})
	return nums[k-1] // k is 1-indexed
}

// ── Approach 2: Prefix-Tree Step Counting (Optimal) ──────────────────────────
//
// prefixTreeCount solves K-th Smallest in Lexicographical Order by walking a
// 10-ary "denary trie" of the numbers 1..n in pre-order, counting how many
// numbers live under each prefix so it can skip whole subtrees instead of
// enumerating them.
//
// Intuition:
//
//	Arrange 1..n as a tree: root has children 1..9, and every node `p` has
//	children `p*10 .. p*10+9` (those that are <= n). A pre-order DFS of this
//	tree visits numbers in exactly lexicographical order. We don't build it —
//	we COUNT. For a current prefix, countUnder(prefix) tells how many numbers in
//	[1,n] start with it. Starting at prefix 1 with a budget of k-1 steps to
//	take: if the whole subtree under `prefix` has <= remaining numbers, skip it
//	all (advance to the next sibling `prefix+1`, subtract the count); otherwise
//	the answer is inside, so descend (prefix *= 10, spend one step). When the
//	budget hits 0, `prefix` is the k-th number.
//
// Algorithm:
//  1. curr = 1; k-- (we are already standing on the 1st number, prefix "1").
//  2. While k > 0:
//     a. cnt = countUnder(curr, n) — numbers in [1,n] with prefix curr.
//     b. If cnt <= k: this whole subtree precedes the target — curr++ (next
//     sibling), k -= cnt.
//     c. Else: target is inside — curr *= 10 (first child), k--.
//  3. Return curr.
//
// Time:  O((log n)^2) — the outer walk takes O(log n) big steps; each countUnder
//
//	does O(log n) work. Independent of k's magnitude.
//
// Space: O(1) — only integer counters.
func prefixTreeCount(n int, k int) int {
	curr := 1 // start at the smallest lexicographic number, prefix "1"
	k--       // standing on the 1st number already; k more steps to walk
	for k > 0 {
		cnt := countUnder(curr, n) // how many numbers in [1,n] begin with `curr`
		if cnt <= k {
			// The entire subtree under curr comes before the target; hop to the
			// next sibling and account for all cnt numbers we just skipped.
			curr++
			k -= cnt
		} else {
			// The target lives within curr's subtree; step down to its first
			// child (curr*10), spending one step to land on that number.
			curr *= 10
			k--
		}
	}
	return curr
}

// countUnder returns how many integers in [1, n] have `prefix` as a leading
// prefix in the denary trie (i.e. prefix, prefix0..prefix9, prefix00.., ...).
//
// It sums, level by level, the size of the intersection of [prefix, next) with
// [1, n], where at each deeper level the range [prefix, next) is multiplied by
// 10. `next` is the first number just past the prefix's range on that level.
func countUnder(prefix int, n int) int {
	count := 0
	cur := prefix      // left edge of the prefix's range at the current level
	next := prefix + 1 // right edge (exclusive) at the current level
	for cur <= n {
		// On this level the prefix covers [cur, next); clamp the right edge to
		// n+1 so we never count numbers greater than n.
		hi := next
		if n+1 < hi {
			hi = n + 1
		}
		count += hi - cur // numbers of this length sharing the prefix
		cur *= 10         // descend one level: ranges widen by a factor of 10
		next *= 10
	}
	return count
}

func main() {
	fmt.Println("=== Approach 1: Brute Force (Generate + Sort Lexicographically) ===")
	fmt.Println(bruteForce(13, 2)) // expected 10
	fmt.Println(bruteForce(1, 1))  // expected 1
	fmt.Println(bruteForce(10, 3)) // expected 2 (order: 1,10,2,3,4,5,6,7,8,9)

	fmt.Println("=== Approach 2: Prefix-Tree Step Counting (Optimal) ===")
	fmt.Println(prefixTreeCount(13, 2))         // expected 10
	fmt.Println(prefixTreeCount(1, 1))          // expected 1
	fmt.Println(prefixTreeCount(10, 3))         // expected 2
	fmt.Println(prefixTreeCount(1000000000, 1)) // expected 1 (huge n, still O(log^2 n))
}
