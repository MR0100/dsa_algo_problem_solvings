# 0190 — Reverse Bits

> LeetCode #190 · Difficulty: Easy
> **Categories:** Bit Manipulation, Divide and Conquer

---

## Problem Statement

Reverse bits of a given 32 bits unsigned integer.

**Note:**

- Note that in some languages, such as Java, there is no unsigned integer type. In this case, both input and output will be given as a signed integer type. They should not affect your implementation, as the integer's internal binary representation is the same, whether it is signed or unsigned.
- In Java, the compiler represents the signed integers using [2's complement notation](https://en.wikipedia.org/wiki/Two%27s_complement). Therefore, in **Example 2** below, the input represents the signed integer `-3` and the output represents the signed integer `-1073741825`.

**Example 1:**

```
Input: n = 00000010100101000001111010011100
Output:    964176192 (00111001011110000101001010000000)
Explanation: The input binary string 00000010100101000001111010011100 represents the unsigned integer 43261596, so return 964176192 which its binary representation is 00111001011110000101001010000000.
```

**Example 2:**

```
Input: n = 11111111111111111111111111111101
Output:   3221225471 (10111111111111111111111111111101)
Explanation: The input binary string 11111111111111111111111111111101 represents the unsigned integer 4294967293, so return 3221225471 which its binary representation is 10111111111111111111111111111101.
```

**Constraints:**

- The input must be a **binary string** of length `32`.

**Follow up:** If this function is called many times, how would you optimize it?

---

## Company Frequency

| Company    | Frequency        | Last Reported |
|------------|------------------|---------------|
| Apple      | ★★★★☆ High       | 2024          |
| Amazon     | ★★★☆☆ Medium     | 2024          |
| Microsoft  | ★★★☆☆ Medium     | 2023          |
| Google     | ★★★☆☆ Medium     | 2023          |
| Adobe      | ★★☆☆☆ Low        | 2023          |
| Bloomberg  | ★★☆☆☆ Low        | 2022          |

> ⚠️ Frequency data is crowd-sourced from LeetCode Discuss, Glassdoor, and
> community interview reports. Treat as a signal, not a guarantee.

---

## DSA Concepts Used

- **Bit Manipulation** — the task is a pure bit-level transform: peel bits off the input with `& 1` / `>>`, push them onto the result with `<<` / `|`, and use fixed masks (`0xaaaaaaaa`, `0x55555555`, …) to move whole blocks of bits in parallel → see [`/dsa/bit_manipulation.md`](/dsa/bit_manipulation.md)
- **Divide and Conquer** — the optimal approach reverses a 32-bit word by swapping its two 16-bit halves, then the bytes inside each half, then nibbles, pairs, and single bits — `log₂(32) = 5` parallel swap levels instead of a 32-step loop → see [`/dsa/divide_and_conquer.md`](/dsa/divide_and_conquer.md)

---

## Approaches Overview

| # | Approach | Time | Space | When to use |
|---|----------|------|-------|-------------|
| 1 | Brute Force (Binary String Round-Trip) | O(32) = O(1) | O(32) = O(1) | Correctness baseline; readable but allocates strings |
| 2 | Bit by Bit | O(32) = O(1) | O(1) | The standard interview answer; no allocations, easy to prove |
| 3 | Byte Lookup Table | O(1) per call | O(256) = O(1) | The **follow-up** answer — many calls amortise the 256-entry table |
| 4 | Divide and Conquer Mask-Swap (Optimal) | O(1) | O(1) | Branch-free, loop-free; 5 constant operations reverse the word |

---

## Approach 1 — Brute Force (Binary String Round-Trip)

### Intuition

Take the problem statement literally. It shows the input as a 32-character binary string, so materialise exactly that string, mirror its characters, and parse the mirror back to an integer. No bit-twiddling insight is needed — this is a plain string reversal that serves as a correctness oracle for the sharper bit-level approaches.

### Algorithm

1. Format `num` as a zero-padded 32-character binary string with `"%032b"`.
2. Build a mirrored buffer where `reversed[i] = bits[31-i]`.
3. Parse the mirrored 32-char string back with base-2 `ParseUint` and return it.

### Complexity

- **Time:** O(32) = O(1) — one fixed format pass, one 32-step mirror, one 32-char parse.
- **Space:** O(32) = O(1) — two fixed 32-byte buffers; these heap-allocate, unlike Approaches 2–4.

### Code

```go
func bruteForce(num uint32) uint32 {
	bits := fmt.Sprintf("%032b", num) // zero-padded so all 32 positions are explicit
	reversed := make([]byte, 32)
	for i := 0; i < 32; i++ {
		reversed[i] = bits[31-i] // mirror character positions around the middle
	}
	// error can't occur: the string is exactly 32 chars of '0'/'1'
	v, _ := strconv.ParseUint(string(reversed), 2, 32)
	return uint32(v)
}
```

### Dry Run

Example 1: `num = 43261596` → `bits = "00000010100101000001111010011100"` (index 0 = leftmost / most-significant char).

| i | source index `31-i` | `bits[31-i]` | reversed builds up |
|---|---------------------|--------------|--------------------|
| 0 | 31 | `0` | `0` |
| 1 | 30 | `0` | `00` |
| 2 | 29 | `1` | `001` |
| 3 | 28 | `1` | `0011` |
| … | … | … | … |
| 29 | 2 | `0` | `0011100101111000010100101` … |
| 30 | 1 | `0` | `…000000` |
| 31 | 0 | `0` | `00111001011110000101001010000000` |

Parse `"00111001011110000101001010000000"` in base 2 → `964176192`. Result: `964176192` ✔

---

## Approach 2 — Bit by Bit

### Intuition

Reversing is "last in, first out" — exactly what repeated shifting gives. Pop the lowest bit off `num` and push it onto the low end of `result`; the first-popped (lowest input) bit gets shoved up as later bits arrive, so it lands highest. After 32 pops the whole word is mirrored, like reversing a list by pushing each element onto a stack.

### Algorithm

1. Initialise `result = 0`.
2. Repeat 32 times: shift `result` left by one, OR in `num`'s lowest bit (`num & 1`), then shift `num` right by one.
3. Return `result`.

### Complexity

- **Time:** O(32) = O(1) — one constant-work iteration per bit.
- **Space:** O(1) — a single accumulator, no allocations.

### Code

```go
func bitByBit(num uint32) uint32 {
	var result uint32
	for i := 0; i < 32; i++ {
		result <<= 1      // make room at the bottom for the next bit
		result |= num & 1 // take num's current lowest bit
		num >>= 1         // consume it
	}
	return result
}
```

### Dry Run

Example 1: `num = 43261596` = `00000010100101000001111010011100`. Tracing the first few and last iterations (bit taken = `num & 1` before the shift):

| i | `num & 1` | `result` before `<<` | `result` after `<<1` then `\| bit` | `num` after `>>1` (low bits) |
|---|-----------|----------------------|------------------------------------|-------------------------------|
| 0 | 0 | `0` | `0` | `…10011100` → `…1001110` |
| 1 | 0 | `0` | `00` | `…100111` |
| 2 | 1 | `00` | `001` | `…10011` |
| 3 | 1 | `001` | `0011` | `…1001` |
| 4 | 1 | `0011` | `00111` | `…100` |
| … | … | … | … | … |
| 31 | 0 | `0011100101111000010100101000000` | `00111001011110000101001010000000` | `0` |

After 32 iterations `result = 00111001011110000101001010000000` = `964176192` ✔

---

## Approach 3 — Byte Lookup Table (Follow-Up: Called Many Times)

### Intuition

This answers the **follow-up**: if the function is called many times, stop re-deriving per-bit work on every call. Reversing a 32-bit word equals reversing each of its 4 bytes **and** reversing the order of those bytes — so the lowest input byte, bit-reversed, becomes the highest output byte. Precompute the 8-bit reversal of every possible byte once (a 256-entry table), then each query is just 4 table lookups and 3 ORs.

### Algorithm

1. One-time: build `revByte[b]` = the 8-bit reversal of `b` for every `b` in `0..255`.
2. Split `num` into its 4 bytes and look each up in the table.
3. Reassemble in mirrored byte order: low byte → bits 24–31, 2nd byte → bits 16–23, 3rd byte → bits 8–15, high byte → bits 0–7.

### Complexity

- **Time:** O(1) per call — 4 lookups and 3 ORs, after a one-time 256×8-step table build.
- **Space:** O(256) = O(1) — the shared lookup table, amortised across all calls.

### Code

```go
func byteTable(num uint32) uint32 {
	return revByte[num&0xff]<<24 | // lowest byte, reversed, becomes highest
		revByte[(num>>8)&0xff]<<16 | // 2nd byte → 3rd slot
		revByte[(num>>16)&0xff]<<8 | // 3rd byte → 2nd slot
		revByte[num>>24] // highest byte, reversed, becomes lowest
}

// revByte[b] holds b with its 8 bits reversed; built once at program start.
var revByte = buildRevByteTable()

func buildRevByteTable() [256]uint32 {
	var table [256]uint32
	for b := 0; b < 256; b++ {
		r := 0
		for i := 0; i < 8; i++ {
			r = r<<1 | b>>i&1 // push b's i-th bit onto r (same LIFO idea as bitByBit)
		}
		table[b] = uint32(r)
	}
	return table
}
```

### Dry Run

Example 1: `num = 43261596` = `00000010 10010100 00011110 10011100` (bytes, high → low).

| byte position | raw byte | `revByte[byte]` (8-bit reversed) | shifted into slot |
|---------------|----------|----------------------------------|-------------------|
| low = `num & 0xff` | `10011100` (156) | `00111001` (57) | `<< 24` → `00111001 00000000 00000000 00000000` |
| `(num>>8) & 0xff` | `00011110` (30) | `01111000` (120) | `<< 16` → `00000000 01111000 00000000 00000000` |
| `(num>>16) & 0xff` | `10010100` (148) | `00101001` (41) | `<< 8`  → `00000000 00000000 00101001 00000000` |
| high = `num >> 24` | `00000010` (2) | `01000000` (64) | `<< 0`  → `00000000 00000000 00000000 01000000` |

OR the four rows together:
`00111001 01111000 00101001 01000000` = `964176192` ✔

---

## Approach 4 — Divide and Conquer Mask-Swap (Optimal)

### Intuition

A reversed word is: its two 16-bit halves swapped, with each half itself reversed. Apply that definition recursively — swap the halves, then swap the bytes inside each half, then the nibbles inside each byte, then bit-pairs, then adjacent bits. Each level swaps **all** blocks of one size *in parallel* with constant masks (`0xff00ff00` selects every high byte, `0xf0f0f0f0` every high nibble, `0xcccccccc` every high pair, `0xaaaaaaaa` every odd bit). So `log₂(32) = 5` operations reverse the entire word — no loop at all.

### Algorithm

1. Swap the two 16-bit halves.
2. Swap adjacent bytes within each half.
3. Swap adjacent nibbles within each byte.
4. Swap adjacent 2-bit pairs within each nibble.
5. Swap adjacent single bits within each pair.

### Complexity

- **Time:** O(1) — exactly 5 shift/mask/OR lines, branch-free and loop-free.
- **Space:** O(1) — the value is transformed in place in a register.

### Code

```go
func divideAndConquer(num uint32) uint32 {
	num = (num >> 16) | (num << 16)                             // 1) swap the two 16-bit halves
	num = ((num & 0xff00ff00) >> 8) | ((num & 0x00ff00ff) << 8) // 2) swap bytes inside each half
	num = ((num & 0xf0f0f0f0) >> 4) | ((num & 0x0f0f0f0f) << 4) // 3) swap nibbles inside each byte
	num = ((num & 0xcccccccc) >> 2) | ((num & 0x33333333) << 2) // 4) swap 2-bit pairs inside each nibble
	num = ((num & 0xaaaaaaaa) >> 1) | ((num & 0x55555555) << 1) // 5) swap adjacent bits inside each pair
	return num
}
```

### Dry Run

Example 1: `num = 00000010 10010100 00011110 10011100` (43261596). Grouping the 32 bits into 8 hex nibbles: `0x0294_1E9C`.

| Step | Operation | `num` after (binary, grouped by byte) | hex |
|------|-----------|----------------------------------------|-----|
| start | — | `00000010 10010100 00011110 10011100` | `0294 1E9C` |
| 1 | swap 16-bit halves | `00011110 10011100 00000010 10010100` | `1E9C 0294` |
| 2 | swap bytes in each half | `10011100 00011110 10010100 00000010` | `9C1E 9402` |
| 3 | swap nibbles in each byte | `11001001 11100001 01001001 00100000` | `C9E1 4920` |
| 4 | swap 2-bit pairs in each nibble | `00110110 10110100 00010110 10000000` | `36B4 1680` |
| 5 | swap adjacent bits in each pair | `00111001 01111000 00101001 01000000` | `3978 2940` |

Final `num = 00111001 01111000 00101001 01000000` = `0x39782940` = `964176192` ✔

---

## Key Takeaways

- **Reversing bits = LIFO.** Peel the lowest bit off the input (`& 1`, `>> 1`) and push it onto the low end of the result (`<< 1`, `| bit`); the earliest-popped bit rises to the top. This "shift-out / shift-in" loop is the bread-and-butter bit reversal.
- **The follow-up wants precomputation.** "Called many times" is a signal to build a lookup table: reverse every byte once (256 entries) and answer each query with a handful of lookups. Reversing a word = reverse each chunk **and** reverse the chunk order.
- **Divide-and-conquer on bits uses parallel masks.** Swapping halves → bytes → nibbles → pairs → bits with fixed masks reverses a 32-bit word in `log₂(32) = 5` branch-free steps. Memorise the mask ladder `0xff00ff00 / 0xf0f0f0f0 / 0xcccccccc / 0xaaaaaaaa` (and their complements) — it recurs in bit-reversal, popcount, and byte-swap tricks.
- **Zero-pad when printing bits.** `%032b` makes all 32 positions explicit; forgetting the width drops leading zeros and silently corrupts a "reverse the 32 bits" task.

---

## Related Problems

- LeetCode #191 — Number of 1 Bits (per-bit peeling with `& 1` / `>> 1`)
- LeetCode #201 — Bitwise AND of Numbers Range (common-prefix bit reasoning)
- LeetCode #338 — Counting Bits (bit DP built on the same shift idioms)
- LeetCode #7 — Reverse Integer (decimal-digit analogue of this reversal)
- LeetCode #371 — Sum of Two Integers (pure bit-manipulation arithmetic)
