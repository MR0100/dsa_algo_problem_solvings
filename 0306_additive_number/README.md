# 0306 — Additive Number

> LeetCode #306 · Difficulty: Medium
> **Categories:** String, Backtracking, Big Number

---

## Problem Statement

An **additive number** is a string whose digits can form an **additive sequence**.

A valid additive sequence should contain **at least** three numbers. Except for the first two numbers, each subsequent number in the sequence must be the sum of the preceding two.

Given a string containing only digits, return `true` if it is an **additive number** or `false` otherwise.

**Note:** Numbers in the additive sequence **cannot** have leading zeros, so sequence `1, 2, 03` or `1, 02, 3` is invalid.

**Example 1:**

```
Input: "112358"
Output: true
Explanation:
The digits can form an additive sequence: 1, 1, 2, 3, 5, 8.
1 + 1 = 2, 1 + 2 = 3, 2 + 3 = 5, 3 + 5 = 8
```

**Example 2:**

```
Input: "199100199"
Output: true
Explanation:
The additive sequence is: 1, 99, 100, 199.
1 + 99 = 100, 99 + 100 = 199
```

**Example 3:**

```
Input: "1203"
Output: false
Explanation: There is no valid additive sequence.
```

**Constraints:**

- `1 <= num.length <= 35`
- `num` consists only of digits.

**Follow up:** How would you handle overflow for very large input integers?

---

## Company Frequency

| Company   | Frequency        | Last Reported |
|-----------|------------------|---------------|
| Epic      | ★★★☆☆ Medium     | 2023          |
| Amazon    | ★★☆☆☆ Low        | 2023          |
| Google    | ★★☆☆☆ Low        | 2022          |
| Microsoft | ★★☆☆☆ Low        | 2022          |

> ⚠️ Frequency data is crowd-sourced from LeetCode Discuss, Glassdoor, and community interview reports. Treat as a signal, not a guarantee.

---

## DSA Concepts Used

- **Backtracking / Exhaustive Split** — the sequence is pinned down by its first two numbers, so we enumerate every prefix pair and verify → see [`/dsa/backtracking.md`](/dsa/backtracking.md)
- **String Algorithms** — leading-zero validation and matching forced sums against string suffixes → see [`/dsa/string_algorithms.md`](/dsa/string_algorithms.md)
- **Big-Number Arithmetic** — up to 35 digits overflows int64; add as decimal strings or with `math/big` → see [`/dsa/math_number_theory.md`](/dsa/math_number_theory.md)

---

## Approaches Overview

| # | Approach | Time | Space | When to use |
|---|----------|------|-------|-------------|
| 1 | Brute Force Backtracking (string add) | O(n³) | O(n) | Interview default; overflow-safe by hand |
| 2 | Iterative Split + `big.Int` (Optimal) | O(n³) | O(n) | Cleanest arbitrary-precision version |

---

## Approach 1 — Brute Force Backtracking

### Intuition

An additive sequence is completely determined by its first two numbers: once `num1` and `num2` are fixed, `num3 = num1 + num2` is forced, then `num2 + num3`, and so on. So we only enumerate the two prefixes and check whether the rest of the string is the forced continuation. Because numbers can reach 35 digits, we add them as decimal strings to stay overflow-free.

### Algorithm

1. Loop `i` over the end (exclusive) of the first number `num[0:i]`.
2. Loop `j` over the end (exclusive) of the second number `num[i:j]`.
3. Reject any candidate with a leading zero (length > 1 starting with `'0'`).
4. Verify the continuation: repeatedly compute `sum = num1 + num2` (string addition), require the remainder to start with `sum`, then slide the window `(num2, sum)`.
5. Return `true` on the first split that consumes the whole string.

### Complexity

- **Time:** O(n³) — O(n²) prefix pairs, each verified in O(n) with O(n) string additions along the way.
- **Space:** O(n) — recursion depth and addition buffers proportional to the string length.

### Code

```go
func bruteForce(num string) bool {
	n := len(num)
	for i := 1; i <= n-2; i++ {
		if num[0] == '0' && i > 1 {
			break
		}
		for j := i + 1; j <= n-1; j++ {
			if num[i] == '0' && j-i > 1 {
				break
			}
			if isValid(num[0:i], num[i:j], num[j:]) {
				return true
			}
		}
	}
	return false
}

func isValid(num1, num2, rest string) bool {
	if len(rest) == 0 {
		return true
	}
	sum := addStrings(num1, num2)
	if len(rest) < len(sum) || rest[:len(sum)] != sum {
		return false
	}
	return isValid(num2, sum, rest[len(sum):])
}

func addStrings(a, b string) string {
	i, j := len(a)-1, len(b)-1
	carry := 0
	res := []byte{}
	for i >= 0 || j >= 0 || carry > 0 {
		sum := carry
		if i >= 0 {
			sum += int(a[i] - '0')
			i--
		}
		if j >= 0 {
			sum += int(b[j] - '0')
			j--
		}
		carry = sum / 10
		res = append(res, byte(sum%10)+'0')
	}
	for l, r := 0, len(res)-1; l < r; l, r = l+1, r-1 {
		res[l], res[r] = res[r], res[l]
	}
	return string(res)
}
```

### Dry Run

Input `num = "112358"` (n = 6).

| i | num1 | j | num2 | Forced continuation of `num[j:]` | Result |
|---|------|---|------|----------------------------------|--------|
| 1 | "1"  | 2 | "1"  | rest="2358": 1+1=2 ✓ (rest "358"), 1+2=3 ✓ (rest "58"), 2+3=5 ✓ (rest "8"), 3+5=8 ✓ (rest "") | **true** |

Since the first split `(1, 1)` validates the whole string, `bruteForce` returns `true` immediately.

---

## Approach 2 — Iterative Split + big.Int (Optimal)

### Intuition

Same search space — every choice of the first two numbers — but the verification is a tight loop, and `math/big` handles arbitrary-precision addition cleanly. This is the version to present after justifying that int64 can overflow (up to 35 digits).

### Algorithm

1. Enumerate end index `i` of `num1` and end index `j` of `num2`, with the same leading-zero pruning.
2. Parse `num1`, `num2` as `big.Int`, set `pos = j`.
3. While `pos < n`: compute `next = num1 + num2`, render it; the suffix from `pos` must start with it, else break; advance `pos`, shift the pair `(num2, next)`.
4. If `pos == n` exactly, the whole string was consumed → return `true`.

### Complexity

- **Time:** O(n³) — O(n²) splits, each validated in O(n) big additions.
- **Space:** O(n) — `big.Int` operands proportional to the number length.

### Code

```go
func bigIntSplit(num string) bool {
	n := len(num)
	for i := 1; i <= n-2; i++ {
		if num[0] == '0' && i > 1 {
			break
		}
		first := new(big.Int)
		first.SetString(num[0:i], 10)
		for j := i + 1; j <= n-1; j++ {
			if num[i] == '0' && j-i > 1 {
				break
			}
			second := new(big.Int)
			second.SetString(num[i:j], 10)
			if validateFrom(num, j, new(big.Int).Set(first), second) {
				return true
			}
		}
	}
	return false
}

func validateFrom(num string, pos int, a, b *big.Int) bool {
	n := len(num)
	for pos < n {
		sum := new(big.Int).Add(a, b)
		s := sum.String()
		if pos+len(s) > n || num[pos:pos+len(s)] != s {
			return false
		}
		pos += len(s)
		a, b = b, sum
	}
	return pos == n
}
```

### Dry Run

Input `num = "199100199"` (n = 9).

| i | num1 | j | num2 | pos | next = a+b | suffix starts with? | new pos |
|---|------|---|------|-----|-----------|---------------------|---------|
| 1 | 1    | 3 | 99   | 3   | 1+99=100  | "100199" → yes      | 6       |
|   |      |   |      | 6   | 99+100=199| "199" → yes         | 9       |
|   |      |   |      | 9   | pos == n  | loop ends           | **true**|

`validateFrom` returns `true`, so `bigIntSplit("199100199")` is `true`.

---

## Key Takeaways

- **Fix the seed, force the rest:** many "sequence exists?" problems collapse to enumerating the first one or two terms; everything after is deterministic.
- **Two nested prefix loops** give O(n²) candidate starts — a recurring pattern for string-partition problems.
- **Leading-zero rule:** a multi-digit number may not start with `'0'`; prune early (`break`) when the first char of a candidate is `'0'` and it has length > 1.
- **Overflow answer for the follow-up:** add as decimal strings or use big integers; never rely on int64 when inputs can be 35 digits.

---

## Related Problems

- LeetCode #842 — Split Array into Fibonacci Sequence (same forced-continuation split, returns the sequence)
- LeetCode #415 — Add Strings (the string-addition helper reused here)
- LeetCode #93 — Restore IP Addresses (enumerate valid prefix partitions with leading-zero rules)
- LeetCode #131 — Palindrome Partitioning (backtracking over string splits)
