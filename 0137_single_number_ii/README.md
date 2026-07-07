# 0137 — Single Number II

> LeetCode #137 · Difficulty: Medium
> **Categories:** Array, Bit Manipulation

---

## Problem Statement

Given an integer array `nums` where every element appears **three times** except for one, which appears **exactly once**. Find the single element and return it.

You must implement a solution with a linear runtime complexity and use only constant extra space.

**Example 1:**
```
Input: nums = [2,2,3,2]
Output: 3
```

**Example 2:**
```
Input: nums = [0,1,0,1,0,1,99]
Output: 99
```

**Constraints:**
- `1 <= nums.length <= 3 * 10^4`
- `-2^31 <= nums[i] <= 2^31 - 1`
- Each element in `nums` appears exactly **three times** except for one element which appears **once**.

---

## Company Frequency

| Company   | Frequency       | Last Reported |
|-----------|-----------------|---------------|
| Google    | ★★★★☆ High      | 2024          |
| Amazon    | ★★★★☆ High      | 2024          |
| Microsoft | ★★★☆☆ Medium    | 2024          |
| Apple     | ★★★☆☆ Medium    | 2023          |
| Facebook  | ★★★☆☆ Medium    | 2023          |

> ⚠️ Frequency data is crowd-sourced from LeetCode Discuss, Glassdoor, and community interview reports. Treat as a signal, not a guarantee.

---

## DSA Concepts Used

- **Bit Manipulation** — per-bit counting mod 3, and the ones/twos bitmask state machine → see [`/dsa/bit_manipulation.md`](/dsa/bit_manipulation.md)
- **Hash Map** — the generic frequency-count fallback → see [`/dsa/hash_map.md`](/dsa/hash_map.md)

---

## Approaches Overview

| # | Approach                    | Time     | Space | When to use                                            |
|---|-----------------------------|----------|-------|--------------------------------------------------------|
| 1 | Brute Force                 | O(n²)    | O(1)  | Explanation baseline only                              |
| 2 | Hash Map                    | O(n)     | O(n)  | Interview warm-up; violates the O(1)-space follow-up   |
| 3 | Math (Set Sum)              | O(n)     | O(n)  | Elegant one-liner idea; still O(n) space               |
| 4 | Bit Counting                | O(32·n)  | O(1)  | Meets constraints; generalizes to "appears k times"    |
| 5 | Ones/Twos Bitmask DFA (Optimal) | O(n) | O(1)  | The canonical optimal answer; single pass, two registers |

---

## Approach 1 — Brute Force

### Intuition
Directly test each element: the single number is the one whose total occurrence count is 1 rather than 3.

### Algorithm
1. For each index `i`:
   1. Count `j` with `nums[j] == nums[i]` over the whole array.
   2. If the count is `1`, return `nums[i]`.

### Complexity
- **Time:** O(n²) — n candidates × O(n) counting scan each.
- **Space:** O(1) — loop counters only.

### Code
```go
func bruteForce(nums []int) int {
	for i := 0; i < len(nums); i++ { // candidate element
		count := 0
		for j := 0; j < len(nums); j++ { // count occurrences everywhere
			if nums[j] == nums[i] {
				count++
			}
		}
		if count == 1 { // appears once → the single number
			return nums[i]
		}
	}
	return -1 // unreachable per problem guarantee
}
```

### Dry Run
`nums = [2,2,3,2]` (Example 1):

| i | nums[i] | matches at indices | count | action        |
|---|---------|--------------------|-------|---------------|
| 0 | 2       | 0, 1, 3            | 3     | keep looking  |
| 1 | 2       | 0, 1, 3            | 3     | keep looking  |
| 2 | 3       | 2                  | 1     | return **3** ✅ |

---

## Approach 2 — Hash Map

### Intuition
Count frequencies in a map, then report the value whose count is 1. Works for any repetition factor, at the cost of O(n) memory.

### Algorithm
1. One pass: `freq[num]++`.
2. Scan the map; return the key with count `1`.

### Complexity
- **Time:** O(n) — two linear passes with amortized O(1) map operations.
- **Space:** O(n) — up to ⌈n/3⌉+1 distinct keys stored.

### Code
```go
func hashMap(nums []int) int {
	freq := make(map[int]int, len(nums)) // value → count
	for _, num := range nums {
		freq[num]++ // tally
	}
	for num, count := range freq {
		if count == 1 { // unique element
			return num
		}
	}
	return -1 // unreachable per problem guarantee
}
```

### Dry Run
`nums = [2,2,3,2]` (Example 1):

| step | num | freq after   |
|------|-----|--------------|
| 1    | 2   | {2:1}        |
| 2    | 2   | {2:2}        |
| 3    | 3   | {2:2, 3:1}   |
| 4    | 2   | {2:3, 3:1}   |
| scan | —   | 2→3 (skip), 3→1 → return **3** ✅ |

---

## Approach 3 — Math (Set Sum)

### Intuition
If every distinct value appeared exactly three times, the array total would be `3·sum(distinct)`. The single number shows up once instead of three times, so the total is short exactly **two** copies of it:

`3·sum(set) − sum(all) = 2·single  →  single = (3·sum(set) − sum(all)) / 2`

### Algorithm
1. One pass: accumulate `sumAll`; the first time each value is seen, add it to `sumSet` (tracked with a set).
2. Return `(3*sumSet − sumAll) / 2`.

### Complexity
- **Time:** O(n) — one pass over the array plus O(1) set ops.
- **Space:** O(n) — the distinct-value set.

### Code
```go
func mathSum(nums []int) int {
	seen := make(map[int]bool, len(nums)) // distinct values
	sumAll, sumSet := 0, 0
	for _, num := range nums {
		sumAll += num // total including triplicates
		if !seen[num] {
			seen[num] = true
			sumSet += num // each distinct value once
		}
	}
	return (3*sumSet - sumAll) / 2 // missing two copies of the single number
}
```

### Dry Run
`nums = [2,2,3,2]` (Example 1):

| step | num | seen after | sumAll | sumSet |
|------|-----|------------|--------|--------|
| 1    | 2   | {2}        | 2      | 2      |
| 2    | 2   | {2}        | 4      | 2      |
| 3    | 3   | {2,3}      | 7      | 5      |
| 4    | 2   | {2,3}      | 9      | 5      |

Result: `(3·5 − 9) / 2 = 6/2 = 3` → **3** ✅

---

## Approach 4 — Bit Counting

### Intuition
Treat every bit position independently. Across the array, a value appearing three times contributes either 0 or 3 to the "set count" of a position — always a multiple of 3. Therefore `count % 3` at each position equals the single number's bit there. Doing the arithmetic in `int32` makes bit 31 land on the sign bit, so negative answers reconstruct correctly under two's complement.

### Algorithm
1. For `bit = 0 .. 31`:
   1. Count how many `nums[i]` have that bit set (`(int32(num)>>bit)&1`).
   2. If `count % 3 != 0`, set that bit in the `int32` result.
2. Convert the `int32` result to `int` and return (sign preserved).

### Complexity
- **Time:** O(32·n) = O(n) — 32 linear passes (constant factor).
- **Space:** O(1) — one 32-bit accumulator and counters.

### Code
```go
func bitCount(nums []int) int {
	var result int32 // assemble the answer bit by bit (int32 keeps the sign bit honest)
	for bit := 0; bit < 32; bit++ {
		count := 0
		for _, num := range nums {
			if (int32(num)>>bit)&1 == 1 { // is this bit set in num?
				count++
			}
		}
		if count%3 != 0 { // triples contribute multiples of 3; remainder = single's bit
			result |= int32(1) << bit // plant the bit (bit 31 wraps to the sign bit correctly)
		}
	}
	return int(result) // int32 → int preserves the two's-complement value
}
```

### Dry Run
`nums = [2,2,3,2]` (Example 1). Binary: 2 = `10`, 3 = `11`.

| bit | set in            | count | count % 3 | result bit |
|-----|-------------------|-------|-----------|------------|
| 0   | 3                 | 1     | 1         | 1          |
| 1   | 2, 2, 3, 2        | 4     | 1         | 1          |
| 2–31| none              | 0     | 0         | 0          |

Result bits `...011` = **3** ✅

---

## Approach 5 — Ones/Twos Bitmask DFA (Optimal)

### Intuition
Approach 4 counts mod 3 with an integer per bit; we can instead store the mod-3 counter **in two bits spread across two registers** and update all 32 lanes simultaneously. `ones` holds bits currently seen once; `twos` holds bits seen twice. A bit's state machine is `00 → 01 → 10 → 00` (never `11`). When a bit completes three sightings it is wiped from both registers; after the full pass, `ones` is exactly the number that appeared once.

The update `ones = (ones ^ num) &^ twos` means: toggle the incoming bits into `ones`, but suppress any bit already at state "two" (that sighting is its third, so it must go to zero, not back to one). The symmetric update for `twos` uses the freshly updated `ones` as the suppressor.

### Algorithm
1. Initialize `ones = 0`, `twos = 0`.
2. For each `num`:
   1. `ones = (ones ^ num) &^ twos`
   2. `twos = (twos ^ num) &^ ones`
3. Return `ones`.

### Complexity
- **Time:** O(n) — a single pass with two bitwise operations per element.
- **Space:** O(1) — two integer registers.

### Code
```go
func onesTwos(nums []int) int {
	ones, twos := 0, 0
	for _, num := range nums {
		ones = (ones ^ num) &^ twos // advance state 00→01, 10→(blocked, cleared below)
		twos = (twos ^ num) &^ ones // advance state 01→10, 10→00
	}
	return ones // bits seen exactly once (mod 3) = the single number
}
```

### Dry Run
`nums = [2,2,3,2]` (Example 1). Values in binary: 2 = `10`, 3 = `11`.

| step | num | ones before | twos before | ones after = (ones⊕num)&^twos | twos after = (twos⊕num)&^ones |
|------|-----|-------------|-------------|-------------------------------|-------------------------------|
| 1    | 10  | 00          | 00          | (00⊕10) &^ 00 = **10**        | (00⊕10) &^ 10 = **00**        |
| 2    | 10  | 10          | 00          | (10⊕10) &^ 00 = **00**        | (00⊕10) &^ 00 = **10**        |
| 3    | 11  | 00          | 10          | (00⊕11) &^ 10 = **01**        | (10⊕11) &^ 01 = **00**        |
| 4    | 10  | 01          | 00          | (01⊕10) &^ 00 = **11**        | (00⊕10) &^ 11 = **00**        |

Final `ones = 11₂ = 3` → **3** ✅

Sanity check — remember the DFA counts **bit** occurrences mod 3, not element occurrences:

| bit | seen in                          | total sightings | mod 3 | answer's bit |
|-----|----------------------------------|-----------------|-------|--------------|
| 0   | the `3` (step 3)                 | 1               | 1     | 1            |
| 1   | all three `2`s **and** the `3`   | 4               | 1     | 1            |

Every bit of every tripled element contributes a multiple of 3, so each lane's mod-3 residue is exactly the single number's bit — bits `11₂` = **3** ✅

---

## Key Takeaways

- The ones/twos trick is a **parallel mod-3 counter**: two bitmask registers encode states 00/01/10 for all 32 bit lanes at once. Generalize to "appears k times" with ⌈log₂k⌉ registers.
- Per-bit counting mod k is the easy-to-derive fallback: bits from k-fold repeats always sum to multiples of k, so the remainder is the loner's bit.
- **Sign-bit care in Go:** ints are 64-bit, so assemble the 32-bit answer in an `int32` (or sign-extend manually) or negative results come out wrong.
- The set-sum identity `(k·sum(set) − sum(all)) / (k−1)` solves the whole family in one line, at O(n) space.
- XOR alone (LeetCode #136) is the special case k = 2 — its "registers" collapse to a single `ones`.

---

## Related Problems

- LeetCode #136 — Single Number (every element appears twice)
- LeetCode #260 — Single Number III (two elements appear once)
- LeetCode #645 — Set Mismatch (counting/XOR hybrid)
- LeetCode #1009 — Complement of Base 10 Integer (bit reconstruction practice)
