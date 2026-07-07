# 0191 — Number of 1 Bits

> LeetCode #191 · Difficulty: Easy
> **Categories:** Bit Manipulation, Divide and Conquer

---

## Problem Statement

Given a positive integer `n`, write a function that returns the number of set bits in its binary representation (also known as the [Hamming weight](https://en.wikipedia.org/wiki/Hamming_weight)).

**Example 1:**

```
Input: n = 11
Output: 3
Explanation: The input binary string 1011 has a total of three set bits.
```

**Example 2:**

```
Input: n = 128
Output: 1
Explanation: The input binary string 10000000 has a total of one set bit.
```

**Example 3:**

```
Input: n = 2147483645
Output: 30
Explanation: The input binary string 1111111111111111111111111111101 has a total of thirty set bits.
```

**Constraints:**

- `1 <= n <= 2^31 - 1`

**Follow up:** If this function is called many times, how would you optimize it?

---

## Company Frequency

| Company   | Frequency        | Last Reported |
|-----------|------------------|---------------|
| Apple     | ★★★★☆ High       | 2024          |
| Microsoft | ★★★★☆ High       | 2024          |
| Amazon    | ★★★☆☆ Medium     | 2024          |
| Google    | ★★★☆☆ Medium     | 2023          |
| Qualcomm  | ★★★★☆ High       | 2024          |
| Samsung   | ★★★☆☆ Medium     | 2023          |
| Meta      | ★★☆☆☆ Low        | 2023          |
| Box       | ★★☆☆☆ Low        | 2022          |

> ⚠️ Frequency data is crowd-sourced from LeetCode Discuss, Glassdoor, and
> community interview reports. Treat as a signal, not a guarantee.

---

## DSA Concepts Used

- **Bit Manipulation** — the whole problem is reading and clearing individual bits; masks, shifts, and the `n & (n-1)` lowest-set-bit trick are the core tools → see [`/dsa/bit_manipulation.md`](/dsa/bit_manipulation.md)
- **Divide and Conquer** — the SWAR approach sums bit-pairs, then nibbles, then bytes: it halves the number of "counters" each round, conquering all 32 bits in log₂32 = 5 merge steps → see [`/dsa/divide_and_conquer.md`](/dsa/divide_and_conquer.md)

---

## Approaches Overview

| # | Approach | Time | Space | When to use |
|---|----------|------|-------|-------------|
| 1 | Brute Force (bit-by-bit mask scan) | O(32) = O(1) | O(1) | First idea; always 32 iterations regardless of input |
| 2 | Brian Kernighan `n & (n-1)` | O(k), k = set bits | O(1) | The expected interview answer; loops only k ≤ 32 times |
| 3 | Lookup Table (8-bit chunks) | O(1) — 4 lookups | O(256) | Follow-up: function called millions of times |
| 4 | Parallel Bit Count (SWAR) | O(1) — ~12 ops, branch-free | O(1) | Hot loops without POPCNT hardware; shows deep bit fluency |
| 5 | Built-in `bits.OnesCount32` (Optimal) | O(1) — 1 CPU instruction | O(1) | Production Go code; compiles to hardware POPCNT |

---

## Approach 1 — Brute Force (Bit-by-Bit Mask Scan)

### Intuition
A `uint32` has exactly 32 bit slots. Peek into each slot with a mask that has only that one bit set and ask "is this bit a 1?" — the number of yes-answers is the Hamming weight. No cleverness, just inspect everything.

### Algorithm
1. Initialise `count = 0`.
2. For every position `i` in `0..31`, build the mask `1 << i` (only bit `i` set).
3. AND the mask with `n`; a non-zero result means bit `i` of `n` is 1 → `count++`.
4. Return `count`.

### Complexity
- **Time:** O(32) = O(1) — the loop always runs exactly 32 times, independent of the value of `n`.
- **Space:** O(1) — only a counter and a mask variable.

### Code
```go
func bruteForce(n uint32) int {
	count := 0 // running total of set bits found so far
	for i := 0; i < 32; i++ { // visit every bit position 0..31
		mask := uint32(1) << i // mask with only bit i set, e.g. i=3 → 0b1000
		if n&mask != 0 {       // AND isolates bit i; non-zero ⇒ that bit is 1
			count++ // bit i is set → record it
		}
	}
	return count
}
```

### Dry Run
Example 1: `n = 11` = `0b1011`.

| i | mask (binary) | n & mask | bit set? | count after |
|---|---------------|----------|----------|-------------|
| 0 | `0001` | `0001` | yes | 1 |
| 1 | `0010` | `0010` | yes | 2 |
| 2 | `0100` | `0000` | no  | 2 |
| 3 | `1000` | `1000` | yes | 3 |
| 4–31 | `10000` … `1<<31` | `0` every time | no | 3 |

Loop ends after i = 31 → return **3**. ✓ (matches expected output 3)

---

## Approach 2 — Brian Kernighan's Trick

### Intuition
Subtracting 1 from `n` flips the lowest set bit to 0 and turns every bit below it to 1 (e.g. `0b10100 − 1 = 0b10011`). AND-ing `n` with `n−1` therefore wipes out *exactly* the lowest set bit and nothing else. Counting how many wipes it takes to reach 0 counts the set bits — and the loop never wastes time on zero bits.

### Algorithm
1. Initialise `count = 0`.
2. While `n != 0`: replace `n` with `n & (n-1)` (drops the lowest set bit), then `count++`.
3. Return `count`.

### Complexity
- **Time:** O(k) where k = number of set bits (≤ 32) — one iteration per set bit, so sparse numbers like `128` finish in a single iteration.
- **Space:** O(1) — in-place bit arithmetic.

### Code
```go
func brianKernighan(n uint32) int {
	count := 0
	for n != 0 { // one iteration per set bit, not per bit position
		n &= n - 1 // n-1 flips the lowest set bit and all zeros below it;
		//            AND-ing erases exactly that lowest set bit from n
		count++ // one set bit removed ⇒ one set bit counted
	}
	return count
}
```

### Dry Run
Example 1: `n = 11` = `0b1011`.

| iteration | n (binary) | n−1 (binary) | n & (n−1) | count after |
|-----------|------------|--------------|-----------|-------------|
| 1 | `1011` | `1010` | `1010` | 1 |
| 2 | `1010` | `1001` | `1000` | 2 |
| 3 | `1000` | `0111` | `0000` | 3 |
| — | `0000` → loop exits | | | 3 |

Return **3**. ✓ Exactly 3 iterations for 3 set bits — compare with 32 iterations in Approach 1.

---

## Approach 3 — Lookup Table (8-Bit Chunks)

### Intuition
This answers the follow-up *"what if the function is called many times?"*: pay once to precompute the popcount of all 256 possible byte values, then every future query is just 4 table reads and 3 additions — no per-bit work ever again. The table itself is built with the recurrence `popcount(i) = popcount(i >> 1) + (i & 1)`: dropping the last bit of `i` gives a smaller value whose answer is already in the table.

### Algorithm
1. (Once, at start-up) fill `table8[b] = popcount(b)` for every `b` in `0..255` via the recurrence above.
2. Slice `n` into 4 bytes using shifts of 0, 8, 16, 24 and masking with `0xFF`.
3. Return the sum of the 4 table entries.

### Complexity
- **Time:** O(1) — exactly 4 array lookups and 3 additions per call, after a one-time O(256) precomputation amortised over all calls.
- **Space:** O(256) = O(1) — a fixed 256-entry table shared by every call.

### Code
```go
// table8 caches the popcount of every possible byte value 0..255.
var table8 [256]int

func init() {
	for i := 1; i < 256; i++ {
		table8[i] = table8[i>>1] + (i & 1) // reuse the answer for i/2, add i's last bit
	}
}

func lookupTable(n uint32) int {
	return table8[n&0xFF] + // byte 0: bits 0..7
		table8[(n>>8)&0xFF] + // byte 1: bits 8..15
		table8[(n>>16)&0xFF] + // byte 2: bits 16..23
		table8[(n>>24)&0xFF] // byte 3: bits 24..31
}
```

### Dry Run
Example 1: `n = 11` = `0x0000000B`.

| step | expression | byte value | table8[value] | running sum |
|------|------------|------------|---------------|-------------|
| 1 | `n & 0xFF` | 11 (`0b1011`) | 3 | 3 |
| 2 | `(n>>8) & 0xFF` | 0 | 0 | 3 |
| 3 | `(n>>16) & 0xFF` | 0 | 0 | 3 |
| 4 | `(n>>24) & 0xFF` | 0 | 0 | 3 |

(Table sanity: `table8[11] = table8[5] + 1 = (table8[2] + 1) + 1 = ((table8[1] + 0) + 1) + 1 = 3`.)

Return **3**. ✓

---

## Approach 4 — Parallel Bit Count (SWAR)

### Intuition
SWAR ("SIMD Within A Register") treats the 32-bit word as 16 two-bit counters, then 8 four-bit counters, then 4 byte counters. Each masked add merges neighbouring counters *in parallel* with a single machine add — a divide-and-conquer on the bits that finishes in log₂32 = 5 conceptual rounds, with zero loops and zero branches. This is essentially the software fallback behind hardware POPCNT.

### Algorithm
1. `n - ((n>>1) & 0x55555555)` — every 2-bit field now holds the popcount of the 2 bits it replaced (identity: a 2-bit value `ab` satisfies `popcount(ab) = ab − a`).
2. `(n & 0x33333333) + ((n>>2) & 0x33333333)` — adjacent 2-bit counts merge into 4-bit field sums.
3. `(n + (n>>4)) & 0x0F0F0F0F` — adjacent nibbles merge into per-byte sums (each ≤ 8, so no overflow).
4. `(n * 0x01010101) >> 24` — the multiply adds all four bytes into the top byte; shift it down for the total.

### Complexity
- **Time:** O(1) — a fixed straight-line sequence of ~12 arithmetic operations, branch-free (great for pipelined CPUs).
- **Space:** O(1) — pure register arithmetic, no memory.

### Code
```go
func parallelCount(n uint32) int {
	n = n - ((n >> 1) & 0x55555555) // each 2-bit field := popcount of its 2 bits
	n = (n & 0x33333333) + ((n >> 2) & 0x33333333) // each 4-bit field := sum of its two 2-bit fields
	n = (n + (n >> 4)) & 0x0F0F0F0F // each byte := sum of its two nibbles (≤ 8, no overflow)
	return int((n * 0x01010101) >> 24) // top byte accumulates b0+b1+b2+b3; shift it down
}
```

### Dry Run
Example 1: `n = 11` = `0b1011` (high bits all 0, shown as 8 bits for clarity).

| step | operation | value (binary) | meaning |
|------|-----------|----------------|---------|
| 0 | start | `0000 1011` | pairs: `10`,`11` |
| 1 | `n>>1 = 0000 0101`; `& 0x55.. = 0000 0101`; `n − that` | `0000 0110` | pair counts: `01`(=1 for `10`), `10`(=2 for `11`) |
| 2 | `n & 0x33.. = 0000 0010`; `(n>>2) & 0x33.. = 0000 0001`; add | `0000 0011` | nibble count: 3 set bits in low nibble |
| 3 | `n + (n>>4) = 0000 0011`; `& 0x0F.. ` | `0000 0011` | byte 0 holds 3, bytes 1–3 hold 0 |
| 4 | `n * 0x01010101 = 0x03030303`; `>> 24` | `3` | top byte = 3+0+0+0 |

Return **3**. ✓

---

## Approach 5 — Built-in Popcount (Optimal)

### Intuition
Population count is so common (cryptography, chess engines, error-correcting codes, bitset cardinality) that CPUs implement it in silicon. Go exposes it as `math/bits.OnesCount32`, which the compiler lowers to the hardware `POPCNT` instruction on amd64/arm64. In production this is the correct answer; in an interview, mention it *after* deriving Approaches 2–4 by hand.

### Algorithm
1. Return `bits.OnesCount32(n)`.

### Complexity
- **Time:** O(1) — a single CPU instruction where supported (SWAR fallback identical to Approach 4 otherwise).
- **Space:** O(1) — no allocations.

### Code
```go
func builtinPopcount(n uint32) int {
	return bits.OnesCount32(n) // compiler intrinsic → hardware POPCNT
}
```

### Dry Run
Example 1: `n = 11`.

| step | state | result |
|------|-------|--------|
| 1 | `bits.OnesCount32(0b1011)` executes as one `POPCNT` instruction | 3 |

Return **3**. ✓ (Internally the fallback path is byte-for-byte the SWAR trace of Approach 4.)

---

## Key Takeaways

- **`n & (n-1)` clears the lowest set bit** — the single most reusable bit trick; it also powers LC 231 *Power of Two* (`n & (n-1) == 0`) and LC 338 *Counting Bits* (`dp[i] = dp[i&(i-1)] + 1`).
- **`n & (-n)` isolates the lowest set bit** — the sibling trick (basis of Fenwick trees); know both and which does which.
- **Loop over set bits, not positions**: Brian Kernighan turns O(bits-in-word) into O(bits-that-are-1) — the general pattern of "iterate only over what exists".
- **Follow-up "called many times" ⇒ precompute**: a 256-entry byte table converts per-call work into 4 O(1) lookups; the classic time–space trade.
- **SWAR halving** (pairs → nibbles → bytes) is divide-and-conquer without recursion; the masks `0x55555555`, `0x33333333`, `0x0F0F0F0F` are worth recognising on sight.
- In Go, reach for **`math/bits`** (`OnesCount32`, `TrailingZeros32`, `LeadingZeros32`, …) — intrinsics that compile to single instructions.

---

## Related Problems

- LeetCode #190 — Reverse Bits (same 32-bit loop/mask machinery)
- LeetCode #231 — Power of Two (`n & (n-1)` in one shot)
- LeetCode #338 — Counting Bits (popcount of every number 0..n via DP on `n & (n-1)`)
- LeetCode #461 — Hamming Distance (popcount of `x XOR y`)
- LeetCode #268 — Missing Number (XOR bit trick family)
- LeetCode #693 — Binary Number with Alternating Bits (bit pattern inspection)
