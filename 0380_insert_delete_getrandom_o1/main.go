package main

import (
	"fmt"
	"math/rand"
)

// randomSetADT is the common interface both implementations satisfy so main()
// can drive the same official operation sequence through each.
type randomSetADT interface {
	Insert(val int) bool
	Remove(val int) bool
	GetRandom() int
}

// ── Approach 1: Slice Only (Brute Force) ─────────────────────────────────────
//
// SliceSet implements the set with just a slice. Membership tests and removals
// scan the slice linearly, so Insert/Remove are O(n) — only GetRandom is O(1).
//
// Intuition:
//
//	The dumbest working model: keep every element in a slice. Insert appends
//	after a linear "already present?" scan; Remove finds the element by scanning
//	then swap-removes it; GetRandom indexes a random slot. It shows exactly what
//	the hash map in Approach 2 buys us: O(1) membership and locate.
//
// Algorithm:
//
//	Insert:    scan for val; if present return false; else append, return true.
//	Remove:    scan for val; if absent return false; else swap with last, shrink.
//	GetRandom: return vals[rand index].
//
// Time:  Insert O(n), Remove O(n), GetRandom O(1).
// Space: O(n).
type SliceSet struct {
	vals []int      // every element currently in the set
	rng  *rand.Rand // deterministic source for reproducible GetRandom
}

// NewSliceSet builds an empty set with a seeded RNG.
func NewSliceSet() *SliceSet {
	return &SliceSet{rng: rand.New(rand.NewSource(1))}
}

// Insert adds val if absent; returns whether it was newly inserted.
func (s *SliceSet) Insert(val int) bool {
	for _, v := range s.vals {
		if v == val {
			return false // already present
		}
	}
	s.vals = append(s.vals, val) // new element
	return true
}

// Remove deletes val if present; returns whether it was there.
func (s *SliceSet) Remove(val int) bool {
	for i, v := range s.vals {
		if v == val {
			last := len(s.vals) - 1
			s.vals[i] = s.vals[last] // overwrite with the last element
			s.vals = s.vals[:last]   // drop the (now duplicated) tail
			return true
		}
	}
	return false // not found
}

// GetRandom returns a uniformly random current element.
func (s *SliceSet) GetRandom() int {
	return s.vals[s.rng.Intn(len(s.vals))]
}

// ── Approach 2: Hash Map + Slice / Swap-Remove (Optimal) ─────────────────────
//
// RandomizedSet gives true O(1) average Insert, Remove, and GetRandom by pairing
// a slice (for O(1) random indexing) with a map val→index (for O(1) locate).
//
// Intuition:
//
//	GetRandom needs contiguous storage → a slice indexed by rand. Insert/Remove
//	need O(1) membership → a hash map from value to its slice index. The trick
//	for O(1) Remove without leaving a hole: swap the target with the LAST slice
//	element, fix that element's index in the map, then pop the tail. Order is
//	irrelevant, so the swap is free.
//
// Algorithm:
//
//	Insert:    if val in map return false; append to slice, record index, true.
//	Remove:    if val not in map return false; look up its index i, move the last
//	           element into slot i (update its map index), pop tail, delete val.
//	GetRandom: slice[rand index].
//
// Time:  Insert O(1) avg, Remove O(1) avg, GetRandom O(1).
// Space: O(n) — slice + map.
type RandomizedSet struct {
	vals []int       // contiguous storage for O(1) random access
	idx  map[int]int // value → its position in vals
	rng  *rand.Rand  // deterministic source for reproducible GetRandom
}

// NewRandomizedSet builds an empty set with a seeded RNG.
func NewRandomizedSet() *RandomizedSet {
	return &RandomizedSet{
		idx: make(map[int]int),
		rng: rand.New(rand.NewSource(1)),
	}
}

// Insert adds val if absent; returns whether it was newly inserted.
func (s *RandomizedSet) Insert(val int) bool {
	if _, ok := s.idx[val]; ok {
		return false // already present
	}
	s.idx[val] = len(s.vals)     // its index is the current tail position
	s.vals = append(s.vals, val) // store it contiguously
	return true
}

// Remove deletes val if present; returns whether it was there.
func (s *RandomizedSet) Remove(val int) bool {
	i, ok := s.idx[val]
	if !ok {
		return false // not present
	}
	last := len(s.vals) - 1
	lastVal := s.vals[last]
	s.vals[i] = lastVal    // move the last element into the freed slot
	s.idx[lastVal] = i     // fix that element's recorded index
	s.vals = s.vals[:last] // pop the tail (now a duplicate)
	delete(s.idx, val)     // forget the removed value
	return true
}

// GetRandom returns a uniformly random current element.
func (s *RandomizedSet) GetRandom() int {
	return s.vals[s.rng.Intn(len(s.vals))]
}

// runExample drives the official LeetCode operation sequence through one
// implementation. GetRandom is non-deterministic in value, so we verify it
// returns a currently-present element rather than a fixed literal.
//
// Ops:  ["RandomizedSet","insert","remove","insert","getRandom","remove","insert","getRandom"]
// Args: [[],[1],[2],[2],[],[1],[2],[]]
func runExample(newSet func() randomSetADT, present map[int]bool) string {
	s := newSet()
	out := []string{"null"} // constructor

	out = append(out, fmt.Sprintf("%t", s.Insert(1))) // true  — 1 was absent
	present[1] = true

	out = append(out, fmt.Sprintf("%t", s.Remove(2))) // false — 2 not present

	out = append(out, fmt.Sprintf("%t", s.Insert(2))) // true  — {1,2}
	present[2] = true

	r1 := s.GetRandom() // must be 1 or 2
	out = append(out, fmt.Sprintf("%t", present[r1]))

	out = append(out, fmt.Sprintf("%t", s.Remove(1))) // true  — {2}
	delete(present, 1)

	out = append(out, fmt.Sprintf("%t", s.Insert(2))) // false — 2 already there

	r2 := s.GetRandom() // set is {2}, must be 2
	out = append(out, fmt.Sprintf("%t", present[r2] && r2 == 2))

	res := "[" + join(out) + "]"
	return res
}

// join renders a []string as comma-separated values (avoids importing strings
// alongside the tiny formatting need here).
func join(xs []string) string {
	out := ""
	for i, x := range xs {
		if i > 0 {
			out += ", "
		}
		out += x
	}
	return out
}

func main() {
	fmt.Println("=== Approach 1: Slice Only (Brute Force) ===")
	fmt.Println(runExample(func() randomSetADT { return NewSliceSet() }, map[int]bool{})) // [null, true, false, true, true, true, false, true]

	fmt.Println("=== Approach 2: Hash Map + Slice / Swap-Remove (Optimal) ===")
	fmt.Println(runExample(func() randomSetADT { return NewRandomizedSet() }, map[int]bool{})) // [null, true, false, true, true, true, false, true]
}
