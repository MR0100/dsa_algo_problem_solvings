# 0402 — Remove K Digits

> LeetCode #402 · Difficulty: Medium
> **Categories:** Greedy, Monotonic Stack, String

---

## Problem Statement

Given string `num` representing a non-negative integer `num`, and an integer `k`, return *the smallest possible integer after removing* `k` *digits from* `num`.

**Example 1:**

```
Input: num = "1432219", k = 3
Output: "1219"
Explanation: Remove the three digits 4, 3, and 2 to form the new number 1219 which is the smallest.
```

**Example 2:**

```
Input: num = "10200", k = 1
Output: "200"
Explanation: Remove the leading 1 and the number is 200. Note that the output must not contain leading zeroes.
```

**Example 3:**

```
Input: num = "10", k = 2
Output: "0"
Explanation: Remove all the digits from the number and it is left with nothing which is 0.
```

**Constraints:**

- `1 <= k <= num.length <= 10^5`
- `num` consists of only digits.
- `num` does not have any leading zeros except for the zero itself.

---

## Company Frequency

| Company    | Frequency        | Last Reported |
|------------|------------------|---------------|
| Amazon     | ★★★★☆ High       | 2024          |
| Google     | ★★★★☆ High       | 2024          |
| Microsoft  | ★★★☆☆ Medium     | 2023          |
| Bloomberg  | ★★★☆☆ Medium     | 2023          |
| Meta       | ★★☆☆☆ Low        | 2023          |
| Uber       | ★★☆☆☆ Low        | 2022          |

> ⚠️ Frequency data is crowd-sourced from LeetCode Discuss, Glassdoor, and
> community interview reports. Treat as a signal, not a guarantee.

---

## DSA Concepts Used

- **Monotonic Stack** — the kept digits must stay non-decreasing from most- to least-significant; a stack that pops any larger predecessor when a smaller digit arrives builds exactly that → see [`/dsa/monotonic_stack.md`](/dsa/monotonic_stack.md)
- **Greedy** — at each step it is always locally optimal to remove the first digit that is larger than its successor, because fixing a more-significant place dominates every lower place → see [`/dsa/greedy.md`](/dsa/greedy.md)
- **String Algorithms** — the input and output are decimal strings; leading-zero trimming and length bookkeeping are string operations → see [`/dsa/string_algorithms.md`](/dsa/string_algorithms.md)

---

## Approaches Overview

| # | Approach | Time | Space | When to use |
|---|----------|------|-------|-------------|
| 1 | Brute Force (delete first descent, k times) | O(k·n) | O(n) | Easy to reason about correctness; too slow for n up to 10⁵ with large k |
| 2 | Monotonic Increasing Stack | O(n) | O(n) | The canonical optimal solution; one linear pass |
| 3 | Preallocated Array Stack | O(n) | O(n) | Same algorithm, allocation-free with an explicit top pointer |

---

## Approach 1 — Brute Force

### Intuition

To make a number smaller by deleting a single digit, find the first **descent** — the first index `i` where `num[i] > num[i+1]` — and remove `num[i]`. That peak sits in a more significant place than anything after it, so eliminating it lowers the number more than any later deletion could. If there is no descent, the digits are non-decreasing and the least valuable digit to keep is the last one, so drop the tail. Repeat this `k` times.

### Algorithm

1. Repeat `k` times:
   1. Scan for the first index `i` with `num[i] > num[i+1]`.
   2. Remove `num[i]`. If no such `i` exists (string is non-decreasing), remove the final character instead.
2. Strip leading zeros from the result; return `"0"` if nothing remains.

### Complexity

- **Time:** O(k·n) — each of the `k` deletions scans the up-to-`n`-length string and splices it.
- **Space:** O(n) — the mutable working copy of the digits.

### Code

```go
func bruteForce(num string, k int) string {
	b := []byte(num) // mutable copy so we can splice out characters
	for ; k > 0; k-- {
		i := 0
		// Walk to the first spot where a digit exceeds its successor.
		for i < len(b)-1 && b[i] <= b[i+1] {
			i++
		}
		// Delete b[i]: if a descent was found, i is that peak; otherwise the
		// string is non-decreasing and i == len(b)-1 (the last, largest digit).
		b = append(b[:i], b[i+1:]...)
	}
	return normalize(string(b))
}
```

### Dry Run

Example 1: `num = "1432219", k = 3`.

| Pass | current | first descent i (num[i] > num[i+1]) | removed | after |
|------|---------|-------------------------------------|---------|-------|
| 1 | `1432219` | i=1 (`4` > `3`) | `4` | `132219` |
| 2 | `132219` | i=1 (`3` > `2`) | `3` | `12219` |
| 3 | `12219` | i=2 (`2` > `1`) | `2` | `1219` |

No leading zeros to strip. Result: `"1219"` ✔

---

## Approach 2 — Monotonic Stack

### Intuition

We want the kept digits to read as small as possible, which means they should be **non-decreasing** from the most-significant place onward: a smaller digit in an earlier place beats anything later. Process digits left to right on a stack. When a new digit `d` arrives, any kept digit larger than `d` is a bad occupant of its place — pop it (spending one removal) so `d` can take a more significant slot. Stop popping when removals run out. If removals remain after the whole string, the stack is non-decreasing, so its largest digits are at the tail — drop the last `k`.

### Algorithm

1. For each digit `d` in `num`:
   1. While `k > 0` and the stack is non-empty and its top `> d`: pop and decrement `k`.
   2. Push `d`.
2. If `k > 0` still remains, remove the last `k` digits (the tail is largest).
3. Strip leading zeros; return `"0"` if empty.

### Complexity

- **Time:** O(n) — each digit is pushed exactly once and popped at most once.
- **Space:** O(n) — the stack of kept digits.

### Code

```go
func monotonicStack(num string, k int) string {
	stack := make([]byte, 0, len(num)) // kept digits, maintained non-decreasing
	for i := 0; i < len(num); i++ {
		d := num[i]
		// Pop any kept digit strictly greater than the incoming one, as long as
		// we still have removals: the new smaller digit improves that place.
		for k > 0 && len(stack) > 0 && stack[len(stack)-1] > d {
			stack = stack[:len(stack)-1]
			k--
		}
		stack = append(stack, d) // the incoming digit is now kept
	}
	// Any leftover removals: the remaining digits are non-decreasing, so the
	// biggest ones are at the end — chop them off.
	stack = stack[:len(stack)-k]
	return normalize(string(stack))
}
```

### Dry Run

Example 1: `num = "1432219", k = 3`. Stack shown left (bottom) to right (top).

| digit | pops (top > digit while k>0) | k after | stack after push |
|-------|------------------------------|---------|------------------|
| `1` | none | 3 | `1` |
| `4` | none (`1` ≤ `4`) | 3 | `1 4` |
| `3` | pop `4` | 2 | `1 3` |
| `2` | pop `3` | 1 | `1 2` |
| `2` | none (`2` ≤ `2`) | 1 | `1 2 2` |
| `1` | pop `2` | 0 | `1 2 1` |
| `9` | none (k = 0) | 0 | `1 2 1 9` |

`k = 0`, so no tail trim. Stack = `1219`, no leading zeros. Result: `"1219"` ✔

---

## Approach 3 — Array Stack

### Intuition

Exactly the same greedy as Approach 2 — maintain a non-decreasing run of kept digits, popping a larger predecessor whenever a smaller digit arrives and removals remain. The only change is representation: use a preallocated `[]byte` with a manual `top` index instead of growing/shrinking a slice. This avoids re-slicing overhead and makes the O(n) push/pop accounting explicit, a form many interviewers like to see.

### Algorithm

1. Allocate a result buffer of length `n`; set `top = 0`.
2. For each digit `d`: while `top > 0` and `k > 0` and `buf[top-1] > d`, do `top--`, `k--`. Then `buf[top] = d`, `top++`.
3. Kept length = `top − k` (leftover removals trim the largest tail digits).
4. Strip leading zeros over `buf[:keptLen]`; return `"0"` if empty.

### Complexity

- **Time:** O(n) — one push and at most one pop per digit.
- **Space:** O(n) — the result buffer.

### Code

```go
func arrayStack(num string, k int) string {
	buf := make([]byte, len(num)) // preallocated stack storage
	top := 0                      // index one past the last kept digit
	for i := 0; i < len(num); i++ {
		d := num[i]
		// Discard larger kept digits while budget allows.
		for top > 0 && k > 0 && buf[top-1] > d {
			top--
			k--
		}
		buf[top] = d // push the incoming digit
		top++
	}
	keptLen := top - k // leftover removals trim the (largest) tail digits
	return normalize(string(buf[:keptLen]))
}
```

### Dry Run

Example 1: `num = "1432219", k = 3`. `buf` shown as its live prefix `buf[:top]`.

| digit | pops while buf[top-1] > digit and k>0 | k after | top after | buf[:top] |
|-------|----------------------------------------|---------|-----------|-----------|
| `1` | none | 3 | 1 | `1` |
| `4` | none | 3 | 2 | `1 4` |
| `3` | pop `4` (top→1) | 2 | 2 | `1 3` |
| `2` | pop `3` (top→1) | 1 | 2 | `1 2` |
| `2` | none | 1 | 3 | `1 2 2` |
| `1` | pop `2` (top→2) | 0 | 3 | `1 2 1` |
| `9` | none (k=0) | 0 | 4 | `1 2 1 9` |

`keptLen = top − k = 4 − 0 = 4` → `1219`. Result: `"1219"` ✔

---

## Key Takeaways

- **"Smallest/largest number after removals" ⇒ monotonic stack.** Keeping digits monotonic while spending a removal budget is the signature pattern; the direction (increasing vs decreasing) depends on whether you minimise or maximise.
- **Leftover budget goes to the tail.** If the string is already monotonic you never popped; unused removals must strip the least-significant (largest) digits from the end.
- **Handle leading zeros and the empty case last.** Trim leading `0`s once at the end, and collapse an empty result to `"0"` — both examples 2 and 3 exist purely to test this.
- **Greedy proof by exchange:** removing the first descent is safe because any other first removal leaves a larger digit in a more significant place, which can never be compensated later.

---

## Related Problems

- LeetCode #316 — Remove Duplicate Letters (monotonic stack with a keep budget)
- LeetCode #1673 — Find the Most Competitive Subsequence (same stack, fixed keep count)
- LeetCode #321 — Create Maximum Number (merge + remove-k generalisation)
- LeetCode #84 — Largest Rectangle in Histogram (monotonic stack core pattern)
- LeetCode #739 — Daily Temperatures (monotonic stack warm-up)
