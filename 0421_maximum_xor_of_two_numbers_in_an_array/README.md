# 0421 — Maximum XOR of Two Numbers in an Array

> LeetCode #421 · Difficulty: Medium
> **Categories:** Bit Manipulation, Trie, Hash Table, Greedy

---

## Problem Statement

Given an integer array `nums`, return *the maximum result of* `nums[i] XOR nums[j]`, where `0 <= i <= j < n`.

**Example 1:**

```
Input: nums = [3,10,5,25,2,8]
Output: 28
Explanation: The maximum result is 5 XOR 25 = 28.
```

**Example 2:**

```
Input: nums = [14,70,53,83,49,91,36,80,92,51,66,70]
Output: 127
```

**Constraints:**

- `1 <= nums.length <= 2 * 10^5`
- `0 <= nums[i] <= 2^31 - 1`

**Follow-up:** Could you do this in `O(n)` runtime?

---

## Company Frequency

| Company    | Frequency        | Last Reported |
|------------|------------------|---------------|
| Amazon     | ★★★★☆ High       | 2024          |
| Google     | ★★★★☆ High       | 2024          |
| Meta       | ★★★☆☆ Medium     | 2023          |
| Microsoft  | ★★★☆☆ Medium     | 2023          |
| Bloomberg  | ★★☆☆☆ Low        | 2023          |
| Apple      | ★★☆☆☆ Low        | 2022          |

> ⚠️ Frequency data is crowd-sourced from LeetCode Discuss, Glassdoor, and
> community interview reports. Treat as a signal, not a guarantee.

---

## DSA Concepts Used

- **Trie (bitwise prefix tree)** — storing each number as a 31-bit root-to-leaf path lets us, for a fixed `x`, greedily pick the opposite bit at every level and find the XOR-maximising partner in O(31) → see [`/dsa/trie.md`](/dsa/trie.md)
- **Bit Manipulation** — the answer is built most-significant-bit first because a high bit outweighs every lower bit combined; XOR is 1 exactly where two bits differ → see [`/dsa/bit_manipulation.md`](/dsa/bit_manipulation.md)
- **Hash Table** — the greedy prefix approach tests "can this bit be 1?" using the identity `a ^ b = c ⇔ a ^ c = b` against a set of number prefixes → see [`/dsa/hash_map.md`](/dsa/hash_map.md)
- **Greedy** — committing to a 1 at the highest reachable bit is never regretted; lower bits can never make up for a missed high bit → see [`/dsa/greedy.md`](/dsa/greedy.md)

---

## Approaches Overview

| # | Approach | Time | Space | When to use |
|---|----------|------|-------|-------------|
| 1 | Brute Force | O(n²) | O(1) | Tiny inputs or a correctness oracle; `n = 2·10⁵` gives 4·10¹⁰ pairs — TLE |
| 2 | Greedy Prefix Hash Set | O(31·n) = O(n) | O(n) | Clean O(n) with no explicit tree; great when you like the bit identity trick |
| 3 | Binary Trie (Optimal) | O(31·n) = O(n) | O(31·n) | The canonical answer; reusable pattern for max/min-XOR-partner problems |

---

## Approach 1 — Brute Force

### Intuition

The task is literally "maximum of `nums[i] ^ nums[j]` over all pairs". Compute that XOR for every pair and keep the biggest. It ignores everything special about binary numbers, so it is only a baseline — but a perfect correctness oracle for the clever approaches.

### Algorithm

1. Initialise `best = 0` (XOR of two non-negative integers is always ≥ 0).
2. For every `i`, and every `j > i`, compute `x = nums[i] ^ nums[j]`.
3. If `x > best`, update `best = x`.
4. Return `best`.

### Complexity

- **Time:** O(n²) — there are ~n²/2 unordered pairs, each XOR is O(1).
- **Space:** O(1) — only the running maximum is kept.

### Code

```go
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
```

### Dry Run

Input `nums = [3,10,5,25,2,8]`. Only the pairs that raise `best` are shown.

| Step | Pair (i,j) | nums[i] ^ nums[j] | best after |
|------|------------|-------------------|------------|
| init | —          | —                 | 0          |
| 1    | (0,1) 3^10 | 9                 | 9          |
| 2    | (0,3) 3^25 | 26                | 26         |
| 3    | (2,3) 5^25 | **28**            | **28**     |
| …    | remaining pairs | all ≤ 28     | 28         |

Return `28`.

---

## Approach 2 — Greedy Prefix Hash Set

### Intuition

Build the answer bit by bit from the most significant bit down, always trying to place a `1`. If the top bits of the answer are already fixed to `answer`, we hope to also set the next bit, giving `candidate = (answer<<1) | 1`. Using `a ^ b = candidate  ⇔  a ^ candidate = b`, we put every number's **prefix** (its top bits down to the current one) into a set; the bit is reachable iff some prefix `p` in the set has `p ^ candidate` also in the set (that pair XORs to `candidate`). Keep the `1` when reachable, otherwise leave `0`.

### Algorithm

1. `answer = 0`.
2. For `bit` from 30 down to 0:
   1. Shift `answer` left by one (new bit starts as 0), form `candidate = answer | 1`.
   2. Build a set of `n >> bit` for every `n` (each number's top-bits prefix).
   3. If any prefix `p` has `p ^ candidate` in the set, set `answer = candidate`.
3. Return `answer`.

### Complexity

- **Time:** O(31·n) = O(n) — 31 bit positions, each scanning all n numbers once.
- **Space:** O(n) — the prefix set holds up to n keys.

### Code

```go
func greedyHashSet(nums []int) int {
	answer := 0 // the maximum XOR built so far (top bits fixed)
	// Values are < 2^31, so bit 30 is the highest meaningful bit.
	for bit := 30; bit >= 0; bit-- {
		answer <<= 1 // shift answer left to make room for the new bit (it is 0 for now)
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
```

### Dry Run

Input `nums = [3,10,5,25,2,8]`. All numbers fit in 5 bits, so bits 30..5 add nothing (`answer` stays 0). We trace bits 4..0. `answer` is shown in 5-bit binary; each step first shifts `answer` left, then tries `candidate = answer | 1`.

Binary values: `3=00011, 10=01010, 5=00101, 25=11001, 2=00010, 8=01000`.

| bit | candidate | prefixes `{n>>bit}` (sorted) | pair XORing to candidate? | answer after |
|-----|-----------|------------------------------|---------------------------|--------------|
| 4   | `00001`   | {0, 1}                       | `0 ^ 1 = 1` → **yes**     | `00001` (1)  |
| 3   | `00011`   | {0, 1, 3}                    | `0 ^ 3 = 3` → **yes**     | `00011` (3)  |
| 2   | `00111`   | {0, 1, 2, 6}                 | `1 ^ 6 = 7` → **yes**     | `00111` (7)  |
| 1   | `01111`   | {1, 2, 4, 5, 12}             | none XOR to 15 → **no**   | `01110` (14) |
| 0   | `11101`   | {2, 3, 5, 8, 10, 25}         | `5 ^ 25 = 28`; but candidate=29 unreachable → **no** | `11100` (28) |

Key subtlety: when a bit is **unreachable**, `answer` keeps the value it had *after the left shift* (new bit = 0). So an unreachable bit 1 leaves `answer = 01110` (14), not 7. The final `answer = 11100 = 28`, realised by the pair `5 ^ 25`.

> The greedy hash-set method and the trie method encode the same greedy bit-by-bit idea; the trie dry run below traces the winning pair `5 ^ 25 = 28` directly.

---

## Approach 3 — Binary Trie (Optimal)

### Intuition

Insert each number as a path of 31 bits (bit 30 first) into a binary trie. To maximise `x ^ partner` for a fixed `x`, at each bit we want the partner's bit to **differ** from `x`'s bit, because differing bits put a `1` in the XOR — and a `1` at a higher position dominates everything below. The trie answers "does a partner with the opposite bit exist here?" in O(1); if yes we follow it (banking a `1`), else we are forced onto the matching branch (a `0`). Run this for every `x` and keep the maximum.

### Algorithm

1. Insert every number into a binary trie keyed by bits 30..0 (MSB first).
2. For each number `x`, walk from the root: at each bit prefer the child for the **opposite** bit — that adds `1<<bit` to this pair's XOR — and fall back to the same-bit child when the opposite is missing.
3. Track the maximum XOR obtained over all `x`.

### Complexity

- **Time:** O(31·n) = O(n) — each insertion and each query touches 31 nodes.
- **Space:** O(31·n) — at most 31 new nodes per inserted number.

### Code

```go
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
				cur |= 1 << b          // a partner differs here: set this bit of the XOR
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

type trieNode struct {
	children [2]*trieNode
}
```

### Dry Run

Input `nums = [3,10,5,25,2,8]`. All fit in 5 bits, so we show bits 4..0. In binary:
`3=00011, 10=01010, 5=00101, 25=11001, 2=00010, 8=01000`. All six are inserted.

Query `x = 5 = 00101` and walk the trie choosing the opposite bit when available:

| bit b | x bit | want opp | opposite branch present? | cur (running XOR) |
|-------|-------|----------|--------------------------|-------------------|
| 4     | 0     | 1        | yes (25 has bit4=1)      | `10000` = 16      |
| 3     | 0     | 1        | yes (25 has bit3=1)      | `11000` = 24      |
| 2     | 1     | 0        | yes (25 has bit2=0)      | `11100` = 28      |
| 1     | 0     | 1        | no partner on this path with bit1=1 → take same | `11100` = 28 |
| 0     | 1     | 0        | matches 25's bit0=1 → forced same | `11100` = 28 |

Query for `x = 5` yields `28` (the partner traced is `25`). No other `x` beats it, so `best = 28`.

---

## Key Takeaways

- **Max-XOR-partner is a trie problem.** Whenever you must maximise (or minimise) `x ^ y` over a set, a binary trie with a greedy opposite-bit walk gives O(n·B) where B is the bit width.
- **High bits dominate.** Any XOR/OR/AND extremal problem is solved most-significant-bit first: locking a `1` at a higher bit beats any configuration of the lower bits.
- **The identity `a ^ b = c ⇔ a ^ c = b`** turns "does a pair XOR to `c`?" into a single hash-set membership test — the engine of the greedy prefix approach.
- Two very different-looking O(n) methods (prefix hash set vs. binary trie) implement the same greedy bit-by-bit idea.

---

## Related Problems

- LeetCode #1707 — Maximum XOR With an Element From Array (trie + offline queries)
- LeetCode #1938 — Maximum Genetic Difference Query (binary trie on a tree)
- LeetCode #208 — Implement Trie (Prefix Tree) (the underlying structure)
- LeetCode #211 — Design Add and Search Words Data Structure (trie variant)
- LeetCode #1803 — Count Pairs With XOR in a Range (bitwise trie counting)
