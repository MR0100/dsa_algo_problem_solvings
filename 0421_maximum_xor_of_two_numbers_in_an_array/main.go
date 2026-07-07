package main

import "fmt"

// ── Approach 1: Brute Force ──────────────────────────────────────────────────
//
// bruteForce solves Maximum XOR of Two Numbers in an Array by trying every
// unordered pair (i, j) and keeping the largest XOR.
//
// Intuition:
//
//	The problem asks for the maximum of nums[i] ^ nums[j] over all pairs. The
//	most direct thing possible is to compute that XOR for every pair and take
//	the max. It is obviously correct and a good baseline to measure the smarter
//	solutions against — the only thing it wastes is the structure of binary
//	numbers (it never exploits that high bits dominate the result).
//
// Algorithm:
//  1. best = 0.
//  2. For each i, for each j > i: candidate = nums[i] ^ nums[j].
//  3. best = max(best, candidate).
//  4. Return best.
//
// Time:  O(n²) — every pair is examined once.
// Space: O(1) — only the running maximum is stored.
func bruteForce(nums []int) int {
	best := 0 // XOR of two non-negative ints is >= 0, so 0 is a safe floor
	for i := 0; i < len(nums); i++ {
		for j := i + 1; j < len(nums); j++ {
			x := nums[i] ^ nums[j] // XOR of this specific pair
			if x > best {
				best = x // remember the largest XOR seen so far
			}
		}
	}
	return best
}

// ── Approach 2: Greedy Prefix Hash Set ───────────────────────────────────────
//
// greedyHashSet solves Maximum XOR of Two Numbers in an Array by building the
// answer one bit at a time, from the most significant bit down, using a hash
// set of number prefixes to test whether a "1" is achievable at each bit.
//
// Intuition:
//
//	Because a higher bit outweighs all lower bits combined, we greedily try to
//	set the answer's bits to 1 from the top down. Suppose we have already fixed
//	the top bits of the answer to some value `answer`. We hope to also set the
//	next bit, i.e. reach `candidate = answer | (1<<bit)`. Using the identity
//	a ^ b == candidate  ⇔  a ^ candidate == b, if we put the top-bit prefixes
//	of every number into a set, a bit is achievable iff some prefix p in the
//	set has (p ^ candidate) also in the set. If yes, keep the bit; otherwise
//	leave it 0 and move on.
//
// Algorithm:
//  1. answer = 0.
//  2. For bit from 30 down to 0:
//     a. candidate = answer with this bit turned on.
//     b. Build a set of every number shifted right by `bit` (its prefix of the
//     top bits down to `bit`).
//     c. If any prefix p has (p ^ candidate) in the set, this bit is reachable:
//     set answer = candidate.
//  3. Return answer.
//
// Time:  O(31·n) = O(n) — for each of 31 bit positions we scan all n numbers.
// Space: O(n) — the prefix set holds up to n entries per bit.
func greedyHashSet(nums []int) int {
	answer := 0 // the maximum XOR built so far (top bits fixed)
	// Values are < 2^31, so bit 30 is the highest meaningful bit.
	for bit := 30; bit >= 0; bit-- {
		answer <<= 1            // shift answer left to make room for the new bit (it is 0 for now)
		candidate := answer | 1 // hypothesis: the new bit can also be a 1

		// Collect each number's prefix = its top bits down to `bit`.
		prefixes := make(map[int]bool, len(nums))
		for _, n := range nums {
			prefixes[n>>bit] = true // drop the low `bit` bits we haven't decided yet
		}

		// candidate is reachable iff two prefixes XOR to it: p ^ q == candidate.
		for p := range prefixes {
			// If (p ^ candidate) is also a stored prefix q, then p ^ q == candidate.
			if prefixes[p^candidate] {
				answer = candidate // lock in the 1 at this bit
				break              // one witnessing pair is enough
			}
		}
	}
	return answer
}

// ── Approach 3: Binary Trie (Optimal) ────────────────────────────────────────
//
// binaryTrie solves Maximum XOR of Two Numbers in an Array by inserting every
// number's 31-bit representation into a binary trie, then for each number
// walking the trie greedily choosing the opposite bit whenever possible to
// maximise its XOR partner.
//
// Intuition:
//
//	Store numbers as root-to-leaf paths of bits (bit 30 first). To maximise
//	x ^ partner for a fixed x, at every bit we want partner's bit to differ
//	from x's bit — that puts a 1 in the result at that (high) position. The
//	trie lets us check, in O(1) per bit, whether a partner with the desired
//	opposite bit exists; if it does we follow it, otherwise we are forced down
//	the same-bit branch (contributing a 0 there). Doing this for every x and
//	taking the max gives the global answer.
//
// Algorithm:
//  1. Insert every number into a binary trie keyed by bits 30..0.
//  2. For each number x, walk the trie from the root: at each bit prefer the
//     child for the opposite bit (adds 1<<bit to this pair's XOR); if it is
//     missing, take the same-bit child.
//  3. Track the best XOR obtained over all x.
//
// Time:  O(31·n) = O(n) — each insert and each query is 31 steps.
// Space: O(31·n) — up to 31 trie nodes per inserted number.
func binaryTrie(nums []int) int {
	const highBit = 30 // numbers < 2^31 → top meaningful bit is 30

	root := &trieNode{} // empty trie

	// Insert every number as a 31-bit path (MSB first).
	for _, n := range nums {
		node := root
		for b := highBit; b >= 0; b-- {
			bit := (n >> b) & 1 // the bit of n at position b
			if node.children[bit] == nil {
				node.children[bit] = &trieNode{} // create the branch on first use
			}
			node = node.children[bit] // descend
		}
	}

	best := 0
	// For each number, find its best XOR partner via a greedy walk.
	for _, n := range nums {
		node := root
		cur := 0 // XOR accumulated for this particular n
		for b := highBit; b >= 0; b-- {
			bit := (n >> b) & 1 // n's bit at position b
			opp := bit ^ 1      // the opposite bit — differing bits give a 1 in XOR
			if node.children[opp] != nil {
				cur |= 1 << b             // a partner differs here: set this bit of the XOR
				node = node.children[opp] // and follow that (better) branch
			} else {
				node = node.children[bit] // forced to match: this bit contributes 0
			}
		}
		if cur > best {
			best = cur // keep the largest XOR across all numbers
		}
	}
	return best
}

// trieNode is one node of the binary trie; children[0] / children[1] are the
// subtrees for a 0-bit and a 1-bit respectively.
type trieNode struct {
	children [2]*trieNode
}

func main() {
	fmt.Println("=== Approach 1: Brute Force ===")
	fmt.Println(bruteForce([]int{3, 10, 5, 25, 2, 8}))                             // expected 28
	fmt.Println(bruteForce([]int{14, 70, 53, 83, 49, 91, 36, 80, 92, 51, 66, 70})) // expected 127

	fmt.Println("=== Approach 2: Greedy Prefix Hash Set ===")
	fmt.Println(greedyHashSet([]int{3, 10, 5, 25, 2, 8}))                             // expected 28
	fmt.Println(greedyHashSet([]int{14, 70, 53, 83, 49, 91, 36, 80, 92, 51, 66, 70})) // expected 127

	fmt.Println("=== Approach 3: Binary Trie (Optimal) ===")
	fmt.Println(binaryTrie([]int{3, 10, 5, 25, 2, 8}))                             // expected 28
	fmt.Println(binaryTrie([]int{14, 70, 53, 83, 49, 91, 36, 80, 92, 51, 66, 70})) // expected 127
}
