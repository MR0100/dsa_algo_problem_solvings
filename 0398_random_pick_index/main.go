package main

import (
	"fmt"
	"math/rand"
)

// ── Approach 1: Hash Map of Indices (Precompute) ─────────────────────────────
//
// SolutionHashMap solves Random Pick Index by precomputing, for every value,
// the full list of indices at which it appears, then picking one uniformly.
//
// Intuition:
//
//	pick(target) must return a uniformly random index among all positions equal
//	to target. If we already have the slice of those positions, a single call
//	to rand.Intn over its length gives a uniform pick in O(1).
//
// Algorithm (constructor):
//  1. Build map[value] -> []index by scanning nums once.
//
// Algorithm (pick):
//  1. Look up the index slice for target.
//  2. Return slice[rand.Intn(len(slice))].
//
// Time:  Constructor O(n); pick O(1).
// Space: O(n) — stores every index grouped by value.
type SolutionHashMap struct {
	idx map[int][]int // value -> all indices where it occurs
}

// NewSolutionHashMap builds the value->indices map from nums.
func NewSolutionHashMap(nums []int) *SolutionHashMap {
	m := make(map[int][]int, len(nums))
	for i, v := range nums {
		m[v] = append(m[v], i) // group indices by their value
	}
	return &SolutionHashMap{idx: m}
}

// Pick returns a uniformly random index i with nums[i] == target.
func (s *SolutionHashMap) Pick(target int) int {
	positions := s.idx[target]                  // all indices holding target
	return positions[rand.Intn(len(positions))] // uniform choice among them
}

// ── Approach 2: Reservoir Sampling (Optimal Space) ───────────────────────────
//
// SolutionReservoir solves Random Pick Index with O(1) extra space by using
// reservoir sampling of size 1 over a single scan per pick.
//
// Intuition:
//
//	Keep only the raw array. On pick(target), stream over the array counting
//	matches; the c-th match replaces the current choice with probability 1/c.
//	After the scan, each match has been kept with probability exactly 1/count,
//	i.e. uniform — the classic size-1 reservoir argument.
//
// Algorithm (pick):
//  1. count = 0, result = -1.
//  2. For each i with nums[i] == target:
//     a. count++.
//     b. With probability 1/count (rand.Intn(count)==0), set result = i.
//  3. Return result.
//
// Time:  Constructor O(1); pick O(n) (one scan).
// Space: O(1) extra — only the original slice reference is kept.
type SolutionReservoir struct {
	nums []int // reference to the input; no per-value preprocessing
}

// NewSolutionReservoir stores the array reference; no preprocessing.
func NewSolutionReservoir(nums []int) *SolutionReservoir {
	return &SolutionReservoir{nums: nums}
}

// Pick reservoir-samples one index equal to target uniformly in a single scan.
func (s *SolutionReservoir) Pick(target int) int {
	count := 0   // how many matches seen so far
	result := -1 // currently chosen index
	for i, v := range s.nums {
		if v != target {
			continue // ignore non-matching positions
		}
		count++
		// The c-th match wins the "seat" with probability 1/c, keeping the
		// distribution uniform over all matches seen so far.
		if rand.Intn(count) == 0 {
			result = i
		}
	}
	return result
}

func main() {
	// Seed for deterministic, reproducible demo output.
	rand.Seed(1)

	// Official example: nums=[1,2,3,3,3], pick(3) must be one of {2,3,4},
	// pick(1) must be 0. We show each returns a VALID index for its target.
	nums := []int{1, 2, 3, 3, 3}

	fmt.Println("=== Approach 1: Hash Map of Indices ===")
	hm := NewSolutionHashMap(nums)
	p3 := hm.Pick(3)
	fmt.Printf("pick(3) -> %d  (valid if in {2,3,4}: %v)\n", p3, p3 == 2 || p3 == 3 || p3 == 4) // expected true
	p1 := hm.Pick(1)
	fmt.Printf("pick(1) -> %d  (valid if == 0: %v)\n", p1, p1 == 0) // expected true

	fmt.Println("=== Approach 2: Reservoir Sampling (Optimal Space) ===")
	rs := NewSolutionReservoir(nums)
	r3 := rs.Pick(3)
	fmt.Printf("pick(3) -> %d  (valid if in {2,3,4}: %v)\n", r3, r3 == 2 || r3 == 3 || r3 == 4) // expected true
	r1 := rs.Pick(1)
	fmt.Printf("pick(1) -> %d  (valid if == 0: %v)\n", r1, r1 == 0) // expected true

	// Distribution check: over many picks of 3, each of {2,3,4} should appear
	// roughly a third of the time — evidence the sampling is uniform.
	fmt.Println("=== Uniformity check (Reservoir, 30000 picks of 3) ===")
	counts := map[int]int{}
	for i := 0; i < 30000; i++ {
		counts[rs.Pick(3)]++
	}
	fmt.Printf("index 2: %d, index 3: %d, index 4: %d  (each ~10000)\n", counts[2], counts[3], counts[4])
}
