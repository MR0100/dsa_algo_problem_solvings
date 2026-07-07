package main

import (
	"fmt"
	"math/rand"
)

// randomizedCollection is the common API every approach implements, so main()
// can drive the same official operation sequence through all of them.
type randomizedCollection interface {
	Insert(val int) bool // true if val was NOT already present
	Remove(val int) bool // true if val was present (one copy removed)
	GetRandom() int      // uniformly random element (by multiplicity)
}

// ── Approach 1: Brute Force (Plain Slice) ────────────────────────────────────
//
// SliceCollection stores every element (including duplicates) in a flat slice.
//
// Intuition:
//
//	The simplest multiset is just a list. Insert appends. GetRandom indexes a
//	random position — automatically weighted by multiplicity. Remove is the
//	slow part: find any occurrence by linear scan and delete it. This is
//	correct but Remove is O(n), which the optimal version fixes with an index
//	map + swap-with-last.
//
// Algorithm:
//
//	Insert:    remember whether val already appears (scan), then append.
//	Remove:    linear-scan for val; if found, delete that index (shift) and
//	           report true; else false.
//	GetRandom: return elems[rand.Intn(len)].
//
// Time:  Insert O(n) (presence scan), Remove O(n), GetRandom O(1).
// Space: O(n) — the element list.
type SliceCollection struct {
	elems []int
	rng   *rand.Rand
}

// NewSliceCollection builds an empty slice-backed collection.
func NewSliceCollection(rng *rand.Rand) *SliceCollection {
	return &SliceCollection{elems: []int{}, rng: rng}
}

// Insert appends val and reports whether it was newly seen.
func (c *SliceCollection) Insert(val int) bool {
	existed := false
	for _, v := range c.elems {
		if v == val {
			existed = true // val already present → return false later
			break
		}
	}
	c.elems = append(c.elems, val) // always store the new copy
	return !existed
}

// Remove deletes one occurrence of val, reporting whether any existed.
func (c *SliceCollection) Remove(val int) bool {
	for i, v := range c.elems {
		if v == val {
			// delete index i by shifting the tail left
			c.elems = append(c.elems[:i], c.elems[i+1:]...)
			return true
		}
	}
	return false // val not found
}

// GetRandom returns a uniformly random element (weighted by multiplicity).
func (c *SliceCollection) GetRandom() int {
	return c.elems[c.rng.Intn(len(c.elems))]
}

// ── Approach 2: Slice + Map of Index Sets (Optimal, O(1) average) ────────────
//
// IndexMapCollection achieves average O(1) Insert, Remove, and GetRandom even
// with duplicates, by pairing a flat slice with a map value → set of the
// indices where that value currently lives.
//
// Intuition:
//
//	GetRandom needs a contiguous array to index; Remove needs O(1) deletion.
//	The trick: keep all elements in a slice `nums`; keep `idx[val]` = the set
//	of positions holding val. To remove one copy of val, take any of its
//	positions, overwrite it with the LAST element of the slice, fix that moved
//	element's index bookkeeping, then shrink the slice. Swap-with-last makes
//	deletion O(1); the index sets keep everything consistent under duplicates.
//
// Algorithm:
//
//	Insert(val):
//	  append val to nums; add its new index (len-1) to idx[val];
//	  return true iff idx[val] had size 1 after (val was previously absent).
//	Remove(val):
//	  if idx[val] empty → false.
//	  pick some index i in idx[val]; let last = len(nums)-1 and lastVal=nums[last].
//	  move lastVal into position i (nums[i]=lastVal); update idx sets:
//	    remove i from idx[val]; if i != last, remove `last` from idx[lastVal]
//	    and add i to idx[lastVal]. shrink nums by one. return true.
//	GetRandom: nums[rand.Intn(len(nums))].
//
// Time:  Insert O(1) avg, Remove O(1) avg, GetRandom O(1).
// Space: O(n) — the slice plus the index sets.
type IndexMapCollection struct {
	nums []int                    // flat storage, contiguous for O(1) random indexing
	idx  map[int]map[int]struct{} // val → set of indices in nums holding val
	rng  *rand.Rand
}

// NewIndexMapCollection builds the empty optimal collection.
func NewIndexMapCollection(rng *rand.Rand) *IndexMapCollection {
	return &IndexMapCollection{
		nums: []int{},
		idx:  map[int]map[int]struct{}{},
		rng:  rng,
	}
}

// Insert appends val and records its index; true iff val was newly introduced.
func (c *IndexMapCollection) Insert(val int) bool {
	if c.idx[val] == nil {
		c.idx[val] = map[int]struct{}{} // first time we ever see val
	}
	c.nums = append(c.nums, val)           // store the new copy at the end
	c.idx[val][len(c.nums)-1] = struct{}{} // record its position
	return len(c.idx[val]) == 1            // size 1 ⇒ val was absent before
}

// Remove deletes one occurrence of val via swap-with-last; false if none exist.
func (c *IndexMapCollection) Remove(val int) bool {
	positions := c.idx[val]
	if len(positions) == 0 {
		return false // no copy of val present
	}
	// grab an arbitrary index i holding val
	var i int
	for p := range positions {
		i = p
		break
	}
	last := len(c.nums) - 1
	lastVal := c.nums[last]

	c.nums[i] = lastVal   // overwrite the removed slot with the last element
	delete(c.idx[val], i) // i no longer holds val

	if i != last { // the moved element lands at a new index i
		delete(c.idx[lastVal], last)   // it is no longer at `last`
		c.idx[lastVal][i] = struct{}{} // it is now at i
	}

	c.nums = c.nums[:last] // drop the (now duplicated) tail slot
	if len(c.idx[val]) == 0 {
		delete(c.idx, val) // keep the map clean when a value is fully gone
	}
	return true
}

// GetRandom returns a uniformly random element (weighted by multiplicity).
func (c *IndexMapCollection) GetRandom() int {
	return c.nums[c.rng.Intn(len(c.nums))]
}

// runExample drives the official operation sequence through one implementation
// and returns the output list in LeetCode's format.
//
// Ops:  ["RandomizedCollection","insert","insert","insert","getRandom","remove","getRandom"]
// Args: [[],[1],[1],[2],[],[1],[]]
// Out:  [null,true,false,true,<1 or 2>,true,<1 or 2>]
func runExample(newColl func() randomizedCollection) []string {
	c := newColl()
	out := []string{"null"}
	out = append(out, fmt.Sprintf("%t", c.Insert(1))) // true  (collection now [1])
	out = append(out, fmt.Sprintf("%t", c.Insert(1))) // false (duplicate; now [1,1])
	out = append(out, fmt.Sprintf("%t", c.Insert(2))) // true  (now [1,1,2])
	r1 := c.GetRandom()                               // 1 with prob 2/3, 2 with prob 1/3
	out = append(out, fmt.Sprintf("in{1,2}:%t", r1 == 1 || r1 == 2))
	out = append(out, fmt.Sprintf("%t", c.Remove(1))) // true  (remove one 1; now [1,2])
	r2 := c.GetRandom()                               // 1 or 2 with prob 1/2 each
	out = append(out, fmt.Sprintf("in{1,2}:%t", r2 == 1 || r2 == 2))
	return out
}

// countRandom samples GetRandom to show multiplicity-weighted uniformity.
func countRandom(get func() int, trials int) map[int]int {
	hist := map[int]int{}
	for t := 0; t < trials; t++ {
		hist[get()]++
	}
	return hist
}

func main() {
	fmt.Println("=== Approach 1: Brute Force (Plain Slice) ===")
	fmt.Println(runExample(func() randomizedCollection {
		return NewSliceCollection(rand.New(rand.NewSource(1)))
	})) // expected [null true false true in{1,2}:true true in{1,2}:true]

	fmt.Println("=== Approach 2: Slice + Index Map (Optimal) ===")
	fmt.Println(runExample(func() randomizedCollection {
		return NewIndexMapCollection(rand.New(rand.NewSource(1)))
	})) // expected [null true false true in{1,2}:true true in{1,2}:true]

	// Distribution check: [1,1,2] ⇒ getRandom ~ 2/3 ones, 1/3 twos.
	fmt.Println("=== Distribution: {1,1,2} over 30000 draws ===")
	c := NewIndexMapCollection(rand.New(rand.NewSource(1)))
	c.Insert(1)
	c.Insert(1)
	c.Insert(2)
	fmt.Println(countRandom(c.GetRandom, 30000)) // expected ~ map[1:20000 2:10000]
}
