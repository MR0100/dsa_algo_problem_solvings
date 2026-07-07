# 0066 — Plus One

> LeetCode #66 · Difficulty: Easy
> **Categories:** Array, Math

---

## Problem Statement

You are given a **large integer** represented as an integer array `digits`, where each `digits[i]` is the `i`th digit of the integer. The digits are ordered from most significant to least significant in left-to-right order. The large integer does not contain any leading `0`'s.

Increment the large integer by one and return the resulting array of digits.

**Example 1**
```
Input:  digits = [1,2,3]
Output: [1,2,4]
Explanation: The array represents the integer 123. Incrementing by one gives 123 + 1 = 124.
```

**Example 2**
```
Input:  digits = [4,3,2,1]
Output: [4,3,2,2]
```

**Example 3**
```
Input:  digits = [9]
Output: [1,0]
```

**Constraints**
- `1 <= digits.length <= 100`
- `0 <= digits[i] <= 9`
- `digits` does not contain any leading `0`'s.

---

## Company Frequency

| Company   | Frequency       | Last Reported |
|-----------|-----------------|---------------|
| Google    | ★★★★★ Very High | 2024          |
| Amazon    | ★★★★☆ High      | 2024          |
| Microsoft | ★★★★☆ High      | 2024          |
| Meta      | ★★★☆☆ Medium    | 2023          |
| Bloomberg | ★★★☆☆ Medium    | 2023          |
| Apple     | ★★☆☆☆ Low       | 2022          |

> ⚠️ Frequency data is crowd-sourced from LeetCode Discuss, Glassdoor, and community interview reports. Treat as a signal, not a guarantee.

---

## DSA Concepts Used

- **Carry Propagation** — incrementing a digit and propagating the carry leftward when a digit is 9.

---

## Approaches Overview

| # | Approach | Time | Space | When to use |
|---|----------|------|-------|-------------|
| 1 | Scan from Right ✅ | O(n) worst / O(1) amortised | O(n) worst / O(1) typical | The only approach needed |

---

## Approach 1 — Scan from Right (Recommended ✅)

### Intuition
Walk from the least significant digit (rightmost). If the digit is less than 9, increment it and return immediately — no carry needed. If it's 9, set it to 0 and carry propagates left. If we exhaust all digits (all were 9), prepend a `1`.

### Algorithm
```
for i = len-1 downto 0:
  if digits[i] < 9: digits[i]++; return digits
  digits[i] = 0        // 9+1=10; write 0, carry 1
// all digits were 9: e.g. [9,9,9] → [1,0,0,0]
prepend 1
```

### Complexity
- **Time:** O(n) worst case (e.g., `[9,9,...,9]`). O(1) amortised — most numbers end in non-9.
- **Space:** O(n) worst case (extra element when all 9s); O(1) otherwise.

### Code
```go
// plusOne solves Plus One by simulating addition from the least significant digit.
//
// Time:  O(n) worst case; O(1) amortised.
// Space: O(n) worst case (new slice when all 9s); O(1) otherwise.
func plusOne(digits []int) []int {
	for i := len(digits) - 1; i >= 0; i-- {
		if digits[i] < 9 {
			digits[i]++
			return digits
		}
		digits[i] = 0 // carry: 9+1=10, write 0
	}
	// all digits were 9 (e.g., [9,9] → [1,0,0])
	return append([]int{1}, digits...)
}
```

### Dry Run — `digits = [1,9,9]`
```
i=2: digits[2]=9 → digits[2]=0. carry.
i=1: digits[1]=9 → digits[1]=0. carry.
i=0: digits[0]=1 < 9 → digits[0]=2. return [2,0,0] ✓
```

### Dry Run — `digits = [9,9,9]`
```
i=2: 9 → 0. i=1: 9 → 0. i=0: 9 → 0. Loop ends.
Prepend 1 → [1,0,0,0] ✓
```

---

## Key Takeaways

- **The only tricky case is all-9s** — all other cases terminate early. Write the `prepend 1` path clearly.
- **In-place modification** — we modify the original slice (no new slice needed in the typical case). Only allocate a new slice for the all-9s edge case.
- **This generalises to arbitrary addition** — see #67 (Add Binary) and #43 (Multiply Strings) for related carry-propagation patterns.

---

## Related Problems

- LeetCode #67 — Add Binary (similar carry propagation in base 2)
- LeetCode #43 — Multiply Strings (carry-based grade-school multiplication)
- LeetCode #369 — Plus One Linked List (same operation on a linked list)
