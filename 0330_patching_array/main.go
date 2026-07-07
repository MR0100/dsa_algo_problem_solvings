package main

import "fmt"

// ── Approach 1: Greedy Patching (Optimal) ────────────────────────────────────
//
// greedyPatching solves Patching Array with the canonical greedy that tracks
// the smallest sum in [1, n] not yet reachable.
//
// Intuition:
//
//	Let `miss` be the smallest positive integer that we cannot yet form as a
//	subset sum. Invariant: everything in [1, miss-1] is already reachable.
//	Look at the next available number nums[i]:
//	  • If nums[i] <= miss, adding it extends reachability from [1, miss-1] to
//	    [1, miss-1+nums[i]], so `miss` grows to miss+nums[i] for free (no patch).
//	  • If nums[i] > miss (or nums are exhausted), no existing number can bridge
//	    the gap at `miss`. The greedy-optimal patch is `miss` itself: it doubles
//	    coverage to [1, 2*miss-1], the largest possible jump, and costs one patch.
//	Repeat until miss > n. Patching `miss` is provably optimal because no value
//	larger than `miss` could fill the hole at `miss`, and no smaller value
//	extends coverage further.
//
// Algorithm:
//  1. miss = 1, i = 0, patches = 0.
//  2. While miss <= n:
//     a. If i < len && nums[i] <= miss: miss += nums[i]; i++ (free extension).
//     b. Else: miss += miss (double coverage by patching `miss`); patches++.
//  3. Return patches.
//
// Time:  O(m + log n) — each num is consumed once (m = len(nums)); each patch
//
//	at least doubles `miss`, so at most ~log2(n) patches occur.
//
// Space: O(1) — a few scalars.
func greedyPatching(nums []int, n int) int {
	// miss is int64 because doubling can push it past 2^31 before the loop's
	// `miss <= n` guard stops it (n itself can be up to 2^31-1). A 32-bit int
	// would overflow on that last doubling; int64 keeps it exact.
	var miss int64 = 1
	i := 0       // index into nums
	patches := 0 // count of numbers we had to add

	for miss <= int64(n) { // continue until every value in [1, n] is reachable
		if i < len(nums) && int64(nums[i]) <= miss {
			// nums[i] plugs into the current coverage without leaving a gap:
			// extend reachability by nums[i] at no cost.
			miss += int64(nums[i])
			i++
		} else {
			// Gap at `miss`: patch it. Adding `miss` doubles the reachable
			// range (best possible jump) and consumes exactly one patch.
			miss += miss // == miss *= 2; may exceed 2^31, hence int64
			patches++
		}
	}
	return patches
}

// ── Approach 2: Greedy Verbose (Explicit Coverage Bound) ─────────────────────
//
// greedyVerbose solves Patching Array with the same greedy, but keeps an
// explicit `covered` variable (the largest reachable sum) so the invariant is
// spelled out. It is a structurally different re-implementation used to
// cross-check greedyPatching.
//
// Intuition:
//
//	Maintain `covered` = the largest value such that every integer in
//	[0, covered] is a reachable subset sum (start covered = 0: only the empty
//	sum 0 is reachable). The first unreachable value is therefore covered+1,
//	which is exactly `miss` from Approach 1. Decision at each step:
//	  • If the next number nums[i] <= covered+1, it fits snugly onto the
//	    covered prefix, so covered grows to covered + nums[i] (no patch).
//	  • Otherwise patch covered+1: the new reachable range becomes
//	    [0, covered + (covered+1)] = [0, 2*covered+1], one patch spent.
//	Stop once covered >= n, meaning [1, n] is fully reachable.
//
// Algorithm:
//  1. covered = 0, i = 0, patches = 0.
//  2. While covered < n:
//     a. next = covered + 1 (smallest currently-unreachable value).
//     b. If i < len && nums[i] <= next: covered += nums[i]; i++.
//     c. Else: covered += next (patch `next`); patches++.
//  3. Return patches.
//
// Time:  O(m + log n) — same accounting as Approach 1.
// Space: O(1) — scalars only.
func greedyVerbose(nums []int, n int) int {
	// covered = upper bound of the contiguous reachable range [0, covered].
	// int64 for the same overflow reason: covered can approach 2^31, and adding
	// `next` (~covered+1) nearly doubles it past the 32-bit limit.
	var covered int64 = 0
	i := 0
	patches := 0

	for covered < int64(n) { // stop once [1, n] is fully covered
		next := covered + 1 // smallest value not yet reachable == "miss"
		if i < len(nums) && int64(nums[i]) <= next {
			// nums[i] attaches to the reachable prefix with no gap.
			covered += int64(nums[i])
			i++
		} else {
			// Patch `next`; reachable range extends to [0, covered+next].
			covered += next // ~doubles covered; may exceed 2^31, hence int64
			patches++
		}
	}
	return patches
}

func main() {
	fmt.Println("=== Approach 1: Greedy Patching (Optimal) ===")
	fmt.Println(greedyPatching([]int{1, 3}, 6))                  // expected 1
	fmt.Println(greedyPatching([]int{1, 5, 10}, 20))             // expected 2
	fmt.Println(greedyPatching([]int{1, 2, 2}, 5))               // expected 0
	fmt.Println(greedyPatching([]int{2}, 1))                     // expected 1
	fmt.Println(greedyPatching([]int{1, 2, 31, 33}, 2147483647)) // expected 28 (large-n stress)

	fmt.Println("=== Approach 2: Greedy Verbose (Explicit Coverage Bound) ===")
	fmt.Println(greedyVerbose([]int{1, 3}, 6))                  // expected 1
	fmt.Println(greedyVerbose([]int{1, 5, 10}, 20))             // expected 2
	fmt.Println(greedyVerbose([]int{1, 2, 2}, 5))               // expected 0
	fmt.Println(greedyVerbose([]int{2}, 1))                     // expected 1
	fmt.Println(greedyVerbose([]int{1, 2, 31, 33}, 2147483647)) // expected 28 (must match Approach 1)
}
