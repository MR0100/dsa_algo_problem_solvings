package main

import (
	"fmt"
	"sort"
)

// ── Approach 1: Brute Force (Unsorted Slice, Check All Pairs) ────────────────
//
// TwoSumBrute solves Two Sum III with an unsorted slice: add appends,
// find scans every pair.
//
// Intuition:
//
//	The laziest data structure: store the stream as-is. add is a trivial
//	append; find pays the full price by testing all O(n^2) pairs. Correct,
//	but find is far too slow once thousands of numbers accumulate.
//
// Algorithm:
//
//	add : append number to the slice.
//	find: double loop over all i < j; return true if nums[i]+nums[j] == value.
//
// Time:  add O(1); find O(n^2) — all pairs in the worst case.
// Space: O(n) — the slice of added numbers.
type TwoSumBrute struct {
	nums []int // every added number, in arrival order
}

// NewTwoSumBrute initializes the object with an empty container.
func NewTwoSumBrute() *TwoSumBrute {
	return &TwoSumBrute{nums: []int{}}
}

// Add stores number in the data structure. Time O(1).
func (t *TwoSumBrute) Add(number int) {
	t.nums = append(t.nums, number) // just remember it
}

// Find reports whether any pair of stored numbers sums to value. Time O(n^2).
func (t *TwoSumBrute) Find(value int) bool {
	for i := 0; i < len(t.nums)-1; i++ {
		for j := i + 1; j < len(t.nums); j++ {
			// Two distinct stored elements (indices differ) forming the sum.
			if t.nums[i]+t.nums[j] == value {
				return true
			}
		}
	}
	return false // no pair matched
}

// ── Approach 2: Sorted Slice + Two Pointers ──────────────────────────────────
//
// TwoSumSorted solves Two Sum III by keeping the numbers sorted: add does a
// binary-search insertion, find runs the classic converging two-pointer scan.
//
// Intuition:
//
//	Problem #167 taught us that a sorted array answers pair-sum queries in
//	O(n) with two pointers and no extra memory. The cost moves to add, which
//	must keep the slice sorted — binary search finds the spot in O(log n) but
//	the insertion shift is O(n).
//
// Algorithm:
//
//	add : binary-search the insertion index, then splice the number in.
//	find: left/right pointers converge; sum < value → left++, sum > value →
//	      right--, equal → true.
//
// Time:  add O(n) — O(log n) search + O(n) shift; find O(n) — pointer sweep.
// Space: O(n) — the sorted slice.
type TwoSumSorted struct {
	nums []int // all added numbers, maintained in non-decreasing order
}

// NewTwoSumSorted initializes the object with an empty container.
func NewTwoSumSorted() *TwoSumSorted {
	return &TwoSumSorted{nums: []int{}}
}

// Add inserts number keeping the slice sorted. Time O(n).
func (t *TwoSumSorted) Add(number int) {
	// sort.SearchInts = lower bound: first index whose element >= number.
	i := sort.SearchInts(t.nums, number)
	t.nums = append(t.nums, 0)     // grow by one slot
	copy(t.nums[i+1:], t.nums[i:]) // shift the tail right to open index i
	t.nums[i] = number             // drop the new number into place
}

// Find reports whether any pair sums to value via two pointers. Time O(n).
func (t *TwoSumSorted) Find(value int) bool {
	left, right := 0, len(t.nums)-1
	for left < right {
		sum := t.nums[left] + t.nums[right]
		switch {
		case sum == value:
			return true // found a pair of distinct elements
		case sum < value:
			left++ // need a bigger sum → advance the small end
		default:
			right-- // need a smaller sum → retreat the large end
		}
	}
	return false // pointers met without hitting the target
}

// ── Approach 3: Hash Map of Counts (Optimal Balance) ─────────────────────────
//
// TwoSumHashMap solves Two Sum III with a frequency map: add is O(1), find
// walks the distinct values looking for complements.
//
// Intuition:
//
//	Store how many times each value was added. For a query, every distinct
//	value x needs its complement value-x to exist; when the complement is x
//	itself, x must have been added at least twice. Iterating distinct values
//	(not all additions) also makes duplicate-heavy streams cheap. This is
//	the expected interview answer: O(1) ingestion, O(n) query, O(n) space.
//
// Algorithm:
//
//	add : counts[number]++.
//	find: for each key x, need = value - x; true if need exists and
//	      (need != x, or counts[x] >= 2).
//
// Time:  add O(1); find O(n) — n = number of distinct values.
// Space: O(n) — one map entry per distinct value.
type TwoSumHashMap struct {
	counts map[int]int // value → how many times it was added
}

// NewTwoSumHashMap initializes the object with an empty container.
func NewTwoSumHashMap() *TwoSumHashMap {
	return &TwoSumHashMap{counts: map[int]int{}}
}

// Add records one more occurrence of number. Time O(1).
func (t *TwoSumHashMap) Add(number int) {
	t.counts[number]++
}

// Find reports whether two stored occurrences sum to value. Time O(n).
func (t *TwoSumHashMap) Find(value int) bool {
	for x := range t.counts {
		need := value - x // complement that would complete the pair
		if need == x {
			// Pairing x with itself needs two separate additions of x.
			if t.counts[x] >= 2 {
				return true
			}
		} else if t.counts[need] > 0 {
			// Distinct complement present at least once.
			return true
		}
	}
	return false // no value has a usable complement
}

// ── Approach 4: Precomputed Pair Sums (Fast Find) ────────────────────────────
//
// TwoSumPairSums solves Two Sum III by materialising every achievable pair
// sum at add time, making find a single set lookup.
//
// Intuition:
//
//	Flip the trade-off of Approach 3: if find dominates the workload, pay at
//	ingestion instead. Each new number forms a pair with every number already
//	stored, so add records all those sums in a set; find is then O(1). The
//	price is O(n) per add and O(n^2) worst-case memory for the sums set.
//
// Algorithm:
//
//	add : for every previously stored number x, insert x+number into the sum
//	      set; then store number.
//	find: return whether value is in the sum set.
//
// Time:  add O(n) — pairs with all existing numbers; find O(1) — set lookup.
// Space: O(n^2) — up to one entry per pair of additions.
type TwoSumPairSums struct {
	nums []int        // all added numbers (duplicates kept: they form real pairs)
	sums map[int]bool // every sum achievable by two distinct additions
}

// NewTwoSumPairSums initializes the object with an empty container.
func NewTwoSumPairSums() *TwoSumPairSums {
	return &TwoSumPairSums{nums: []int{}, sums: map[int]bool{}}
}

// Add records number and all new pair sums it creates. Time O(n).
func (t *TwoSumPairSums) Add(number int) {
	// Every existing element pairs with the newcomer exactly once.
	for _, x := range t.nums {
		t.sums[x+number] = true
	}
	t.nums = append(t.nums, number)
}

// Find reports whether value is an achievable pair sum. Time O(1).
func (t *TwoSumPairSums) Find(value int) bool {
	return t.sums[value] // set membership answers the query directly
}

func main() {
	fmt.Println("=== Approach 1: Brute Force (Unsorted Slice, Check All Pairs) ===")
	b := NewTwoSumBrute()
	b.Add(1)                                              // [] --> [1]
	b.Add(3)                                              // [1] --> [1,3]
	b.Add(5)                                              // [1,3] --> [1,3,5]
	fmt.Printf("find(4)=%v  expected true\n", b.Find(4))  // 1 + 3 = 4
	fmt.Printf("find(7)=%v  expected false\n", b.Find(7)) // no two integers sum to 7
	bDup := NewTwoSumBrute()
	bDup.Add(0)
	fmt.Printf("dup: find(0)=%v  expected false\n", bDup.Find(0)) // single 0 can't pair with itself
	bDup.Add(0)
	fmt.Printf("dup: find(0)=%v  expected true\n", bDup.Find(0)) // 0 + 0 = 0 after second add

	fmt.Println("=== Approach 2: Sorted Slice + Two Pointers ===")
	s := NewTwoSumSorted()
	s.Add(1)
	s.Add(3)
	s.Add(5)
	fmt.Printf("find(4)=%v  expected true\n", s.Find(4))
	fmt.Printf("find(7)=%v  expected false\n", s.Find(7))
	sDup := NewTwoSumSorted()
	sDup.Add(0)
	fmt.Printf("dup: find(0)=%v  expected false\n", sDup.Find(0))
	sDup.Add(0)
	fmt.Printf("dup: find(0)=%v  expected true\n", sDup.Find(0))

	fmt.Println("=== Approach 3: Hash Map of Counts (Optimal Balance) ===")
	h := NewTwoSumHashMap()
	h.Add(1)
	h.Add(3)
	h.Add(5)
	fmt.Printf("find(4)=%v  expected true\n", h.Find(4))
	fmt.Printf("find(7)=%v  expected false\n", h.Find(7))
	hDup := NewTwoSumHashMap()
	hDup.Add(0)
	fmt.Printf("dup: find(0)=%v  expected false\n", hDup.Find(0))
	hDup.Add(0)
	fmt.Printf("dup: find(0)=%v  expected true\n", hDup.Find(0))

	fmt.Println("=== Approach 4: Precomputed Pair Sums (Fast Find) ===")
	p := NewTwoSumPairSums()
	p.Add(1)
	p.Add(3)
	p.Add(5)
	fmt.Printf("find(4)=%v  expected true\n", p.Find(4))
	fmt.Printf("find(7)=%v  expected false\n", p.Find(7))
	pDup := NewTwoSumPairSums()
	pDup.Add(0)
	fmt.Printf("dup: find(0)=%v  expected false\n", pDup.Find(0))
	pDup.Add(0)
	fmt.Printf("dup: find(0)=%v  expected true\n", pDup.Find(0))
}
