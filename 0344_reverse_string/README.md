# 0344 — Reverse String

> LeetCode #344 · Difficulty: Easy
> **Categories:** Two Pointers, String, Recursion

---

## Problem Statement

Write a function that reverses a string. The input string is given as an array of characters `s`.

You must do this by modifying the input array **in-place** with `O(1)` extra memory.

**Example 1:**

```
Input: s = ["h","e","l","l","o"]
Output: ["o","l","l","e","h"]
```

**Example 2:**

```
Input: s = ["H","a","n","n","a","h"]
Output: ["h","a","n","n","a","H"]
```

**Constraints:**

- `1 <= s.length <= 10^5`
- `s[i]` is a printable ascii character.

---

## Company Frequency

| Company    | Frequency        | Last Reported |
|------------|------------------|---------------|
| Amazon     | ★★★★☆ High       | 2024          |
| Microsoft  | ★★★☆☆ Medium     | 2023          |
| Google     | ★★★☆☆ Medium     | 2023          |
| Apple      | ★★☆☆☆ Low        | 2023          |
| Facebook   | ★★☆☆☆ Low        | 2022          |

> ⚠️ Frequency data is crowd-sourced from LeetCode Discuss, Glassdoor, and
> community interview reports. Treat as a signal, not a guarantee.

---

## DSA Concepts Used

- **Two Pointers** — swap the element at each end and step both pointers inward until they meet; the canonical in-place reversal → see [`/dsa/two_pointers.md`](/dsa/two_pointers.md)
- **String** — reversal is a fundamental string/array manipulation → see [`/dsa/string_algorithms.md`](/dsa/string_algorithms.md)
- **Recursion** — the reversal can also be expressed as "swap ends, recurse on the middle" → see [`/dsa/divide_and_conquer.md`](/dsa/divide_and_conquer.md)

---

## Approaches Overview

| # | Approach | Time | Space | When to use |
|---|----------|------|-------|-------------|
| 1 | Two Pointers (Optimal) | O(n) | O(1) | The intended answer; true in-place, constant space |
| 2 | Recursion | O(n) | O(n) | Illustrates the recursive framing; uses call-stack space |

---

## Approach 1 — Two Pointers (Optimal)

### Intuition
The first char must end at the last position, the second at the second-last, and so on. Walk one pointer from the left, one from the right; swap the pair they point at and step both inward. Stop when they cross the middle.

### Algorithm
1. `left = 0`, `right = len(s)-1`.
2. While `left < right`: swap `s[left]` and `s[right]`; `left++`; `right--`.

### Complexity
- **Time:** O(n) — each element is visited once.
- **Space:** O(1) — swaps are in place.

### Code
```go
func twoPointers(s []byte) {
	left, right := 0, len(s)-1
	for left < right {
		s[left], s[right] = s[right], s[left]
		left++
		right--
	}
}
```

### Dry Run
Input `s = ["h","e","l","l","o"]`:

| Step | left | right | swap | s |
|------|------|-------|------|---|
| 0 | 0 | 4 | h↔o | `o e l l h` |
| 1 | 1 | 3 | e↔l | `o l l e h` |
| 2 | 2 | 2 | left == right, stop | `o l l e h` |

Output `["o","l","l","e","h"]`.

---

## Approach 2 — Recursion

### Intuition
`reverse(left, right)` = swap the two ends, then `reverse(left+1, right-1)`. The recursion bottoms out when the window has 0 or 1 element (`left >= right`).

### Algorithm
1. `helper(left, right)`: if `left >= right`, return.
2. Swap `s[left]` and `s[right]`.
3. Recurse `helper(left+1, right-1)`.

### Complexity
- **Time:** O(n) — n/2 swaps.
- **Space:** O(n) — recursion stack of depth n/2 (so *not* O(1) extra; it violates the follow-up constraint but is instructive).

### Code
```go
func recursion(s []byte) {
	var helper func(left, right int)
	helper = func(left, right int) {
		if left >= right {
			return
		}
		s[left], s[right] = s[right], s[left]
		helper(left+1, right-1)
	}
	helper(0, len(s)-1)
}
```

### Dry Run
Input `s = ["h","e","l","l","o"]`:

| Call | left | right | Action | s after |
|------|------|-------|--------|---------|
| helper(0,4) | 0 | 4 | swap h↔o | `o e l l h` |
| helper(1,3) | 1 | 3 | swap e↔l | `o l l e h` |
| helper(2,2) | 2 | 2 | left >= right → return | `o l l e h` |

Output `["o","l","l","e","h"]`.

---

## Key Takeaways

- In-place reversal is the archetypal two-pointer, opposite-ends pattern — memorise the `left < right` swap loop.
- Go tuple assignment `s[left], s[right] = s[right], s[left]` swaps without a temp variable.
- Recursion gives the same result but costs O(n) stack space, so it does not satisfy an O(1)-extra-space requirement — prefer the iterative version in interviews.
- Strings in Go are immutable; convert to `[]byte` (or `[]rune` for Unicode) to mutate in place.

---

## Related Problems

- LeetCode #345 — Reverse Vowels of a String (two pointers with a skip condition)
- LeetCode #541 — Reverse String II (reverse in blocks)
- LeetCode #151 — Reverse Words in a String (reverse whole, then each word)
- LeetCode #7 — Reverse Integer (reversal on digits)
- LeetCode #234 — Palindrome Linked List (two pointers / reversal on a list)
