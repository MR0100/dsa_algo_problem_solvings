# 0246 — Strobogrammatic Number

> LeetCode #246 · Difficulty: Easy
> **Categories:** Hash Table, Two Pointers, String, Math

---

## Problem Statement

Given a string `num` which represents an integer, return `true` *if* `num` *is a **strobogrammatic number***.

A **strobogrammatic number** is a number that looks the same when rotated `180` degrees (looked at upside down).

**Example 1:**

```
Input: num = "69"
Output: true
```

**Example 2:**

```
Input: num = "88"
Output: true
```

**Example 3:**

```
Input: num = "962"
Output: false
```

**Constraints:**

- `1 <= num.length <= 50`
- `num` consists of only digits.
- `num` does not contain any leading zeros except for zero itself.

---

## Company Frequency

| Company    | Frequency        | Last Reported |
|------------|------------------|---------------|
| Google     | ★★★★☆ High       | 2024          |
| Facebook   | ★★★☆☆ Medium     | 2023          |
| Microsoft  | ★★☆☆☆ Low        | 2023          |
| Amazon     | ★★☆☆☆ Low        | 2022          |

> ⚠️ Frequency data is crowd-sourced from LeetCode Discuss, Glassdoor, and
> community interview reports. Treat as a signal, not a guarantee.

---

## DSA Concepts Used

- **Two Pointers** — converge from both ends, pairing each digit with its mirror partner across the center → see [`/dsa/two_pointers.md`](/dsa/two_pointers.md)
- **Hash Map** — a fixed lookup table maps each rotatable digit to its 180° image (`0↔0, 1↔1, 8↔8, 6↔9`) → see [`/dsa/hash_map.md`](/dsa/hash_map.md)
- **String Algorithms** — the check is fundamentally about a string being equal to its rotated reverse → see [`/dsa/string_algorithms.md`](/dsa/string_algorithms.md)

---

## Approaches Overview

| # | Approach | Time | Space | When to use |
|---|----------|------|-------|-------------|
| 1 | Two Pointers with Map (Optimal) | O(n) | O(1) | Best: no extra allocation, early exit on first mismatch |
| 2 | Build Rotated String and Compare | O(n) | O(n) | Most literal reading of the definition; easy to explain |

---

## Approach 1 — Two Pointers with Map (Optimal)

### Intuition

Rotating a number 180° reverses its digit order AND flips each digit. So the number is unchanged iff, pairing the `i`-th digit from the left with the `i`-th from the right, each left digit rotates exactly into its partner. Only `{0,1,6,8,9}` have a valid upside-down form; `6↔9` swap, the rest map to themselves. Any other digit disqualifies the number.

### Algorithm

1. Build a map `rotate`: `0→0, 1→1, 8→8, 6→9, 9→6`.
2. Set `left = 0`, `right = len(num)-1`.
3. While `left <= right`:
   1. If `num[left]` is not in `rotate`, return `false`.
   2. If `rotate[num[left]] != num[right]`, return `false`.
   3. `left++`, `right--`.
4. Return `true`.

### Complexity

- **Time:** O(n) — each digit is visited at most once from one side.
- **Space:** O(1) — the map has a constant 5 entries.

### Code

```go
func twoPointers(num string) bool {
	rotate := map[byte]byte{
		'0': '0',
		'1': '1',
		'8': '8',
		'6': '9',
		'9': '6',
	}

	left, right := 0, len(num)-1
	for left <= right {
		mirror, ok := rotate[num[left]]
		if !ok {
			return false
		}
		if mirror != num[right] {
			return false
		}
		left++
		right--
	}
	return true
}
```

### Dry Run

Trace `num = "69"` (`left=0, right=1`):

| Step | left | right | num[left] | rotate[num[left]] | num[right] | Match? |
|------|------|-------|-----------|-------------------|------------|--------|
| 1    | 0    | 1     | `6`       | `9`               | `9`        | yes → left=1, right=0 |
| 2    | 1    | 0     | loop ends (left > right) | — | — | — |

Every pair matched → return `true`. ✓

---

## Approach 2 — Build Rotated String and Compare

### Intuition

The most direct reading of "looks the same rotated 180°": literally construct the rotated number and compare it to the original. Building it means rotating each digit and reversing the order (last digit becomes first). If any digit is not rotatable, the number cannot be strobogrammatic.

### Algorithm

1. Walk `num` from the last index to the first.
2. For each digit, look up its rotation; if none, return `false`.
3. Append the rotation to a buffer (walking backwards builds the reversal).
4. Return whether the built string equals `num`.

### Complexity

- **Time:** O(n) — one reverse pass to build plus one comparison.
- **Space:** O(n) — the rotated-string buffer.

### Code

```go
func buildAndCompare(num string) bool {
	rotate := map[byte]byte{
		'0': '0', '1': '1', '8': '8', '6': '9', '9': '6',
	}

	rotated := make([]byte, 0, len(num))
	for i := len(num) - 1; i >= 0; i-- {
		mirror, ok := rotate[num[i]]
		if !ok {
			return false
		}
		rotated = append(rotated, mirror)
	}
	return string(rotated) == num
}
```

### Dry Run

Trace `num = "69"`:

| Step | i | num[i] | rotate[num[i]] | rotated so far |
|------|---|--------|----------------|----------------|
| 1    | 1 | `9`    | `6`            | `6`            |
| 2    | 0 | `6`    | `9`            | `69`           |

`rotated = "69"` equals `num = "69"` → return `true`. ✓

---

## Key Takeaways

- Strobogrammatic = the string equals its **rotated reverse**; only `{0,1,6,8,9}` are rotatable and `6↔9` is the only cross pair.
- Two-pointer converging with a fixed lookup table gives O(1) space and early exit — the canonical answer.
- When the center digit exists (odd length), the loop's `left == right` case forces it to be self-rotatable (`0,1,8`), which the same map check handles automatically.

---

## Related Problems

- LeetCode #247 — Strobogrammatic Number II (generate all of a given length)
- LeetCode #248 — Strobogrammatic Number III (count in a range)
- LeetCode #9 — Palindrome Number (self-equal under simple reversal)
- LeetCode #125 — Valid Palindrome (two-pointer symmetry check)
