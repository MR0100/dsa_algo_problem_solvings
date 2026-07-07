package main

import (
	"fmt"
	"math/rand"
	"sort"
)

// ── Approach 1: Brute Force (Draw-from-a-bag) ────────────────────────────────
//
// BruteForceSolution shuffles by repeatedly picking a random element out of a
// "bag" (remaining unused elements) and appending it to the result.
//
// Intuition:
//
//	A fair shuffle is exactly "draw items one at a time, uniformly, without
//	replacement". Model the bag as a mutable copy of the array. Each round,
//	pick a random remaining index, take that element, and physically remove
//	it (swap-with-last + shrink) so it can't be drawn again. This is provably
//	uniform but removing from an arbitrary index is the wasteful part the
//	optimal version fixes.
//
// Algorithm:
//  1. Copy the original into a "bag" slice.
//  2. While the bag is non-empty: pick k = rand.Intn(len(bag)); append bag[k]
//     to the result; remove bag[k] by overwriting it with the last element and
//     shrinking the bag.
//  3. Return the result. Reset() returns a copy of the original.
//
// Time:  reset O(n); shuffle O(n) — each of n draws is O(1) with swap-remove.
// Space: O(n) — the original copy plus the bag/result.
type BruteForceSolution struct {
	original []int
	rng      *rand.Rand
}

// NewBruteForceSolution stores an immutable copy of nums.
func NewBruteForceSolution(nums []int, rng *rand.Rand) *BruteForceSolution {
	orig := make([]int, len(nums))
	copy(orig, nums) // keep a pristine copy for Reset
	return &BruteForceSolution{original: orig, rng: rng}
}

// Reset returns the array to its original configuration.
func (s *BruteForceSolution) Reset() []int {
	out := make([]int, len(s.original))
	copy(out, s.original) // never expose the internal slice
	return out
}

// Shuffle returns a uniformly random permutation via draw-from-bag.
func (s *BruteForceSolution) Shuffle() []int {
	bag := make([]int, len(s.original))
	copy(bag, s.original) // mutable working copy
	res := make([]int, 0, len(bag))
	for len(bag) > 0 {
		k := s.rng.Intn(len(bag)) // pick a random remaining element
		res = append(res, bag[k]) // draw it
		bag[k] = bag[len(bag)-1]  // swap the hole with the last element
		bag = bag[:len(bag)-1]    // shrink: that element is now used
	}
	return res
}

// ── Approach 2: Fisher–Yates Shuffle (Optimal, in-place) ─────────────────────
//
// FisherYatesSolution shuffles the array in place using the Fisher–Yates
// (Knuth) algorithm: every one of the n! permutations is equally likely.
//
// Intuition:
//
//	Walk i from the last index down to 1. At each step pick j uniformly in
//	[0, i] and swap arr[i], arr[j]. This "locks in" position i with a
//	uniformly chosen element from the still-unfixed prefix. Because index i is
//	filled from i+1 equally likely candidates and later choices never disturb
//	it, all n! orderings occur with probability 1/n!. It is in-place, O(n),
//	and the canonical fair shuffle.
//
// Algorithm:
//  1. Copy the current array so shuffle doesn't destroy the original.
//  2. For i from n-1 down to 1: j = rand.Intn(i+1); swap arr[i], arr[j].
//  3. Return arr. Reset() returns a copy of the original.
//
// Time:  reset O(n); shuffle O(n) — one pass, O(1) per step.
// Space: O(n) — the returned shuffled copy (in-place over that copy).
type FisherYatesSolution struct {
	original []int
	rng      *rand.Rand
}

// NewFisherYatesSolution stores an immutable copy of nums.
func NewFisherYatesSolution(nums []int, rng *rand.Rand) *FisherYatesSolution {
	orig := make([]int, len(nums))
	copy(orig, nums)
	return &FisherYatesSolution{original: orig, rng: rng}
}

// Reset returns the original configuration (fresh copy).
func (s *FisherYatesSolution) Reset() []int {
	out := make([]int, len(s.original))
	copy(out, s.original)
	return out
}

// Shuffle returns a uniformly random permutation via Fisher–Yates.
func (s *FisherYatesSolution) Shuffle() []int {
	arr := make([]int, len(s.original))
	copy(arr, s.original) // shuffle a copy, keep original intact
	for i := len(arr) - 1; i >= 1; i-- {
		j := s.rng.Intn(i + 1)          // uniform in [0, i]
		arr[i], arr[j] = arr[j], arr[i] // fix position i with a random unfixed element
	}
	return arr
}

// permKey turns a slice into a comparable string key for the histogram.
func permKey(a []int) string { return fmt.Sprint(a) }

// isUniform runs many shuffles and checks every permutation appears roughly
// equally often (within tolerance). Returns true if plausibly uniform.
func isUniform(shuffle func() []int, n, trials int, tol float64) bool {
	hist := map[string]int{}
	for t := 0; t < trials; t++ {
		hist[permKey(shuffle())]++
	}
	// count number of distinct permutations produced
	fact := 1
	for i := 2; i <= n; i++ {
		fact *= i
	}
	if len(hist) != fact {
		return false // some permutation never appeared
	}
	expected := float64(trials) / float64(fact)
	for _, c := range hist {
		if float64(c) < expected*(1-tol) || float64(c) > expected*(1+tol) {
			return false
		}
	}
	return true
}

func main() {
	// Official example: nums = [1,2,3]; shuffle returns a random permutation,
	// reset returns [1,2,3].
	nums := []int{1, 2, 3}

	fmt.Println("=== Approach 1: Brute Force (Draw-from-a-bag) ===")
	bf := NewBruteForceSolution(nums, rand.New(rand.NewSource(1)))
	sh := bf.Shuffle()
	sorted := append([]int(nil), sh...)
	sort.Ints(sorted)
	fmt.Println(fmt.Sprint(sorted) == "[1 2 3]")       // expected true (shuffle is a permutation of [1,2,3])
	fmt.Println(bf.Reset())                            // expected [1 2 3]
	fmt.Println(isUniform(bf.Shuffle, 3, 60000, 0.10)) // expected true (all 6 perms ~equally likely)

	fmt.Println("=== Approach 2: Fisher–Yates (Optimal) ===")
	fy := NewFisherYatesSolution(nums, rand.New(rand.NewSource(1)))
	sh2 := fy.Shuffle()
	sorted2 := append([]int(nil), sh2...)
	sort.Ints(sorted2)
	fmt.Println(fmt.Sprint(sorted2) == "[1 2 3]")      // expected true (shuffle is a permutation of [1,2,3])
	fmt.Println(fy.Reset())                            // expected [1 2 3]
	fmt.Println(isUniform(fy.Shuffle, 3, 60000, 0.10)) // expected true (all 6 perms ~equally likely)
}
