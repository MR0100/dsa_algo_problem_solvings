package main

import (
	"container/heap"
	"fmt"
)

// ── Approach 1: Greedy with Max-Heap ─────────────────────────────────────────
//
// maxHeapGreedy solves Rearrange String k Distance Apart by always placing the
// most-frequent still-available character next, using a cooldown queue to keep
// each character on hold for k positions.
//
// Intuition:
//
//	To spread characters at least k apart, greedily emit the character with the
//	highest remaining count (it is the scarcest resource — hardest to place).
//	After emitting a character, it must "cool down" for k-1 more slots before it
//	may be used again, so park it in a FIFO cooldown queue of size k. When the
//	queue reaches size k, the front becomes eligible and, if it still has count,
//	returns to the heap. If at any point the heap is empty but characters remain
//	on cooldown, no valid arrangement exists.
//
// Algorithm:
//  1. Count frequencies; push (count, char) into a max-heap by count.
//  2. While the heap is non-empty: pop the top, append its char, decrement its
//     count, and push it into the cooldown queue.
//  3. Once the cooldown queue has k entries, pop its front; if that char still
//     has count > 0, push it back into the heap.
//  4. If the result length equals the input length, return it; else "".
//
// Time:  O(n log a) — n emissions, each a heap op over at most a distinct chars.
// Space: O(a) — heap + cooldown queue over the alphabet.
func maxHeapGreedy(s string, k int) string {
	if k <= 1 {
		return s // any arrangement (or the string itself) already satisfies k<=1
	}
	freq := map[byte]int{}
	for i := 0; i < len(s); i++ {
		freq[s[i]]++ // tally each character
	}

	h := &charHeap{}
	for c, f := range freq {
		heap.Push(h, charCount{c, f}) // seed heap with every distinct char
	}

	type cooling struct {
		char  byte
		count int
	}
	queue := []cooling{} // FIFO of characters waiting out their cooldown
	result := make([]byte, 0, len(s))

	for h.Len() > 0 {
		top := heap.Pop(h).(charCount) // most frequent available char
		result = append(result, top.char)
		top.count-- // one occurrence consumed
		// Park it on cooldown regardless of remaining count (we filter later).
		queue = append(queue, cooling{top.char, top.count})

		// Only after k placements does the oldest char become eligible again.
		if len(queue) >= k {
			front := queue[0]
			queue = queue[1:]
			if front.count > 0 {
				heap.Push(h, charCount{front.char, front.count})
			}
		}
	}

	if len(result) != len(s) {
		return "" // ran out of eligible chars before placing all — impossible
	}
	return string(result)
}

// charCount pairs a character with its remaining count for the heap.
type charCount struct {
	char  byte
	count int
}

// charHeap is a max-heap ordered by count (most frequent on top).
type charHeap []charCount

func (h charHeap) Len() int            { return len(h) }
func (h charHeap) Less(i, j int) bool  { return h[i].count > h[j].count } // max-heap
func (h charHeap) Swap(i, j int)       { h[i], h[j] = h[j], h[i] }
func (h *charHeap) Push(x interface{}) { *h = append(*h, x.(charCount)) }
func (h *charHeap) Pop() interface{} {
	old := *h
	n := len(old)
	item := old[n-1]
	*h = old[:n-1]
	return item
}

// ── Approach 2: Greedy Round-Robin (Sort by Count) ───────────────────────────
//
// roundRobinGreedy solves Rearrange String k Distance Apart by an explicit
// slot-filling schedule: place the most frequent char first into positions
// 0, k, 2k, …, then the next char continues filling remaining slots.
//
// Intuition:
//
//	Feasibility hinges on the most frequent char m with count f: it needs
//	(f-1) gaps of size >= k, i.e. at least (f-1)*k + 1 total positions. If
//	n < (f-1)*k + 1 it is impossible. Otherwise, sort characters by descending
//	count and drop them into an array by columns: fill index 0, k, 2k, … first
//	(one char after another), wrapping to index 1, k+1, … This guarantees two
//	equal characters land at least k apart, because we never place a character
//	twice within the same "column stride" until k slots have passed.
//
// Algorithm:
//  1. Count and sort chars by descending frequency.
//  2. Feasibility: if (maxFreq-1)*k + (# chars with maxFreq) > n ⇒ "".
//     (Equivalently maxFreq-1 full rows of width k plus the last partial row.)
//  3. Walk an index that jumps by k each placement, wrapping to the next start
//     column (0,1,2,…) when it overflows; assign chars in frequency order.
//  4. Join the filled array.
//
// Time:  O(n + a log a) — counting, sort of the small alphabet, linear fill.
// Space: O(n) — the output buffer.
func roundRobinGreedy(s string, k int) string {
	if k <= 1 {
		return s
	}
	freq := map[byte]int{}
	for i := 0; i < len(s); i++ {
		freq[s[i]]++
	}
	// Collect (char, count) and sort by count descending (simple insertion).
	pairs := make([]charCount, 0, len(freq))
	for c, f := range freq {
		pairs = append(pairs, charCount{c, f})
	}
	for i := 1; i < len(pairs); i++ { // insertion sort by count desc
		for j := i; j > 0 && pairs[j].count > pairs[j-1].count; j-- {
			pairs[j], pairs[j-1] = pairs[j-1], pairs[j]
		}
	}

	n := len(s)
	// Feasibility: the most frequent char needs (maxFreq-1)*k + 1 slots minimum.
	maxFreq := pairs[0].count
	if (maxFreq-1)*k+1 > n {
		return "" // cannot spread the top char far enough
	}

	res := make([]byte, n)
	idx := 0 // current slot to fill; jumps by k each time
	for _, p := range pairs {
		for c := 0; c < p.count; c++ {
			res[idx] = p.char // place one occurrence
			idx += k          // next occurrence of this (and the next) char is k away
			if idx >= n {
				idx = idx%k + 1 // wrap to the next start column (0→1→2…)
			}
		}
	}
	return string(res)
}

func main() {
	fmt.Println("=== Approach 1: Greedy Max-Heap ===")
	fmt.Printf("s=\"aabbcc\", k=3   got=%q  expected valid (e.g. \"abcabc\")\n", maxHeapGreedy("aabbcc", 3))
	fmt.Printf("s=\"aaabc\",  k=3   got=%q  expected \"\"\n", maxHeapGreedy("aaabc", 3))
	fmt.Printf("s=\"aaadbbcc\", k=2 got=%q  expected valid (e.g. \"abacabcd\")\n", maxHeapGreedy("aaadbbcc", 2))
	fmt.Printf("valid(\"aabbcc\",3,heap)  = %t  expected true\n", isValid(maxHeapGreedy("aabbcc", 3), "aabbcc", 3))
	fmt.Printf("valid(\"aaadbbcc\",2,heap)= %t  expected true\n", isValid(maxHeapGreedy("aaadbbcc", 2), "aaadbbcc", 2))

	fmt.Println("=== Approach 2: Greedy Round-Robin ===")
	fmt.Printf("s=\"aabbcc\", k=3   got=%q  expected valid (e.g. \"abcabc\")\n", roundRobinGreedy("aabbcc", 3))
	fmt.Printf("s=\"aaabc\",  k=3   got=%q  expected \"\"\n", roundRobinGreedy("aaabc", 3))
	fmt.Printf("s=\"aaadbbcc\", k=2 got=%q  expected valid (e.g. \"abacabcd\")\n", roundRobinGreedy("aaadbbcc", 2))
	fmt.Printf("valid(\"aabbcc\",3,rr)    = %t  expected true\n", isValid(roundRobinGreedy("aabbcc", 3), "aabbcc", 3))
	fmt.Printf("valid(\"aaadbbcc\",2,rr)  = %t  expected true\n", isValid(roundRobinGreedy("aaadbbcc", 2), "aaadbbcc", 2))
}

// isValid checks that `out` is a permutation of `s` in which equal characters
// are at least k apart. Used only to verify the greedy outputs in main().
func isValid(out, s string, k int) bool {
	if out == "" {
		return false
	}
	if len(out) != len(s) {
		return false
	}
	// Same multiset of characters?
	cnt := map[byte]int{}
	for i := 0; i < len(s); i++ {
		cnt[s[i]]++
	}
	for i := 0; i < len(out); i++ {
		cnt[out[i]]--
	}
	for _, v := range cnt {
		if v != 0 {
			return false
		}
	}
	// Each character's occurrences at least k apart?
	last := map[byte]int{}
	for i := 0; i < len(out); i++ {
		if p, ok := last[out[i]]; ok && i-p < k {
			return false
		}
		last[out[i]] = i
	}
	return true
}
