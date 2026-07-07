# 0415 — Add Strings

> LeetCode #415 · Difficulty: Easy
> **Categories:** Math, String, Simulation, Two Pointers

---

## Problem Statement

Given two non-negative integers, `num1` and `num2` represented as string, return *the sum of `num1` and `num2` as a string*.

You must solve the problem **without using any built-in library for handling large integers** (such as `BigInteger`), and **without converting the inputs to integers directly**.

**Example 1:**

```
Input: num1 = "11", num2 = "123"
Output: "134"
```

**Example 2:**

```
Input: num1 = "456", num2 = "77"
Output: "533"
```

**Example 3:**

```
Input: num1 = "0", num2 = "0"
Output: "0"
```

**Constraints:**

- `1 <= num1.length, num2.length <= 10^4`
- `num1` and `num2` consist of only digits.
- `num1` and `num2` don't have any leading zeros except for the zero itself.

---

## Company Frequency

| Company    | Frequency        | Last Reported |
|------------|------------------|---------------|
| Google     | ★★★★☆ High       | 2024          |
| Amazon     | ★★★☆☆ Medium     | 2024          |
| Facebook   | ★★★☆☆ Medium     | 2023          |
| Microsoft  | ★★☆☆☆ Low        | 2023          |
| Airbnb     | ★★☆☆☆ Low        | 2022          |
| Bloomberg  | ★★☆☆☆ Low        | 2022          |

> ⚠️ Frequency data is crowd-sourced from LeetCode Discuss, Glassdoor, and
> community interview reports. Treat as a signal, not a guarantee.

---

## DSA Concepts Used

- **String manipulation (digit arithmetic)** — the number lives in a string and must be added a character at a time, converting `'0'..'9'` to/from integer values via ASCII offset → see [`/dsa/string_algorithms.md`](/dsa/string_algorithms.md)
- **Two pointers (converging from both ends)** — two indices walk each string from its least-significant digit toward the most-significant, aligning columns → see [`/dsa/two_pointers.md`](/dsa/two_pointers.md)
- **Elementary number theory (positional addition & carry)** — column sum splits into output digit `sum%10` and carry `sum/10`, the base-10 place-value algorithm → see [`/dsa/math_number_theory.md`](/dsa/math_number_theory.md)

---

## Approaches Overview

| # | Approach | Time | Space | When to use |
|---|----------|------|-------|-------------|
| 1 | Two Pointers + Carry (Optimal) | O(max(m, n)) | O(max(m, n)) | The canonical answer; build reversed, flip once |
| 2 | Pre-Sized Buffer (No Reverse) | O(max(m, n)) | O(max(m, n)) | Avoids the reverse by writing back-to-front into a fixed buffer |

> `m = len(num1)`, `n = len(num2)`. Output space is unavoidable — the sum has up to `max(m,n)+1` digits.

---

## Approach 1 — Two Pointers + Carry (Grade-School Addition)

### Intuition

Add the two numbers exactly the way you would on paper: line them up at the right, and for each column sum the two current digits plus the carry from the previous column. The digit written is `sum % 10`; the carry forward is `sum / 10`. When one number is shorter, its missing high-order digits count as `0`. After the last column, any leftover carry becomes a new leading digit (e.g. `99 + 1 = 100`). Digits come out least-significant first, so reverse the buffer once at the end. ASCII makes the digit conversions trivial: `c - '0'` gives the value, `value + '0'` gives the character.

### Algorithm

1. Set `i = len(num1)-1`, `j = len(num2)-1`, `carry = 0`.
2. While `i >= 0` **or** `j >= 0` **or** `carry > 0`:
   - `d1 = num1[i]-'0'` if `i >= 0` else `0`; decrement `i`.
   - `d2 = num2[j]-'0'` if `j >= 0` else `0`; decrement `j`.
   - `sum = d1 + d2 + carry`; append `(sum%10)+'0'`; `carry = sum/10`.
3. Reverse the appended digits and return them as a string.

### Complexity

- **Time:** O(max(m, n)) — one iteration per output digit.
- **Space:** O(max(m, n)) — the result buffer (required for the answer).

### Code

```go
func twoPointers(num1 string, num2 string) string {
	i, j := len(num1)-1, len(num2)-1 // start at the least-significant digit of each
	carry := 0                       // running carry into the current column
	var sb strings.Builder           // collects result digits in REVERSE order

	// Continue while either number has digits left OR a carry still needs placing.
	for i >= 0 || j >= 0 || carry > 0 {
		sum := carry // begin the column with the incoming carry
		if i >= 0 {
			sum += int(num1[i] - '0') // ASCII digit → integer value
			i--
		}
		if j >= 0 {
			sum += int(num2[j] - '0')
			j--
		}
		sb.WriteByte(byte(sum%10) + '0') // low digit of the column → output (as ASCII)
		carry = sum / 10                 // high part carries to the next column
	}

	return reverseString(sb.String()) // digits were emitted LSB-first; flip them
}
```

### Dry Run

Example 1: `num1 = "11", num2 = "123"`.

| Step | i | j | num1[i] | num2[j] | carry in | sum | write (sum%10) | carry out | builder (reversed) |
|------|---|---|---------|---------|----------|-----|----------------|-----------|--------------------|
| 1 | 1 | 2 | `1` | `3` | 0 | 4 | `4` | 0 | `4` |
| 2 | 0 | 1 | `1` | `2` | 0 | 3 | `3` | 0 | `43` |
| 3 | -1 | 0 | — (0) | `1` | 0 | 1 | `1` | 0 | `431` |

Loop ends (`i,j < 0`, `carry = 0`). Reverse `"431"` → `"134"` ✔

---

## Approach 2 — Pre-Sized Buffer (No Reverse)

### Intuition

The same column-by-column addition, but skip the final reverse. The sum of an `m`-digit and an `n`-digit number has at most `max(m,n)+1` digits, so allocate exactly that many bytes and fill them from the **rightmost** cell leftward — which is the natural direction addition already runs in. Each column's digit lands directly in its final position. If the top carry turned out to be `0`, the single leading cell was never written, so we return the slice starting one past the last write position, trimming it.

### Algorithm

1. `size = max(m, n) + 1`; `buf = make([]byte, size)`; `pos = size-1`.
2. Add columns exactly as in Approach 1, but write `buf[pos] = (sum%10)+'0'` and decrement `pos` each step.
3. Return `string(buf[pos+1:])` — the filled region, already most-significant-first.

### Complexity

- **Time:** O(max(m, n)) — single pass.
- **Space:** O(max(m, n)) — one pre-allocated buffer, no separate reversal copy.

### Code

```go
func preSizedBuffer(num1 string, num2 string) string {
	m, n := len(num1), len(num2)
	size := maxInt(m, n) + 1  // +1 leaves room for a possible final carry
	buf := make([]byte, size) // result digits, written back-to-front
	pos := size - 1           // next cell to write (rightmost first)

	i, j, carry := m-1, n-1, 0
	for i >= 0 || j >= 0 || carry > 0 {
		sum := carry
		if i >= 0 {
			sum += int(num1[i] - '0')
			i--
		}
		if j >= 0 {
			sum += int(num2[j] - '0')
			j--
		}
		buf[pos] = byte(sum%10) + '0' // place digit directly in its final position
		pos--                         // move one cell to the left
		carry = sum / 10
	}
	// buf[pos+1:] is the filled, correctly-ordered result (leading cell, if the
	// top carry was 0, is simply left out by starting the slice at pos+1).
	return string(buf[pos+1:])
}
```

### Dry Run

Example 1: `num1 = "11", num2 = "123"`. `size = max(2,3)+1 = 4`, `buf = [_,_,_,_]`, `pos = 3`.

| Step | i | j | sum | buf[pos] written | pos after | buf state |
|------|---|---|-----|------------------|-----------|-----------|
| 1 | 1 | 2 | 4 | `buf[3]='4'` | 2 | `[_,_,_,4]` |
| 2 | 0 | 1 | 3 | `buf[2]='3'` | 1 | `[_,_,3,4]` |
| 3 | -1 | 0 | 1 | `buf[1]='1'` | 0 | `[_,1,3,4]` |

Loop ends. Return `buf[pos+1:] = buf[1:] = "134"` ✔ (the unused `buf[0]` is trimmed).

---

## Key Takeaways

- **ASCII digit arithmetic:** `c - '0'` converts a digit character to its value and `v + '0'` converts back — the workhorse of string-number problems, no `strconv` needed.
- **The `carry > 0` loop condition is essential.** Keeping the loop alive while a carry remains handles the "grows a new leading digit" case (`"99" + "1" = "100"`) with no special-casing.
- **Pad the shorter number with implicit zeros** by guarding each index with `i >= 0` instead of pre-padding the strings.
- **Build-reversed-then-flip vs. pre-sized-back-to-front** are the two standard ways to emit a number whose digits are produced least-significant-first; the second avoids the reverse at the cost of one length computation.
- This is the linear-string template behind Add Two Numbers (linked list), Multiply Strings, and Add Binary — only the base and the container change.

---

## Related Problems

- LeetCode #2 — Add Two Numbers (same carry logic on a linked list)
- LeetCode #67 — Add Binary (identical algorithm in base 2)
- LeetCode #43 — Multiply Strings (grade-school multiplication built on this addition)
- LeetCode #989 — Add to Array-Form of Integer (digit-array variant)
- LeetCode #66 — Plus One (carry propagation over a digit array)
