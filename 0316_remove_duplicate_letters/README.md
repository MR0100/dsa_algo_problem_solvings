# 0316 — Remove Duplicate Letters

> LeetCode #316 · Difficulty: Medium
> **Categories:** String, Stack, Greedy, Monotonic Stack

---

## Problem Statement

Given a string `s`, remove duplicate letters so that every letter appears once
and only once. You must make sure your result is the **smallest in
lexicographical order** among all possible results.

**Example 1:**

```
Input:  s = "bcabc"
Output: "abc"
```

**Example 2:**

```
Input:  s = "cbacdcbc"
Output: "acdb"
```

Explanation: The distinct letters are `a, b, c, d`. Among all orderings that
keep each letter once and respect the relative feasibility of the original
string, `"acdb"` is the lexicographically smallest.

**Constraints:**

- `1 <= s.length <= 10^4`
- `s` consists of lowercase English letters.

**Note:** This problem is the same as LeetCode #1081 — Smallest Subsequence of
Distinct Characters.

---

## Company Frequency

| Company   | Frequency        | Last Reported |
|-----------|------------------|---------------|
| Google    | ★★★★☆ High       | 2024          |
| Amazon    | ★★★★☆ High       | 2024          |
| Meta      | ★★★☆☆ Medium     | 2023          |
| Bloomberg | ★★★☆☆ Medium     | 2023          |
| Microsoft | ★★☆☆☆ Low        | 2023          |

> ⚠️ Frequency data is crowd-sourced from LeetCode Discuss, Glassdoor, and
> community interview reports. Treat as a signal, not a guarantee.

---

## DSA Concepts Used

- **Monotonic Stack** — maintain a near-increasing result buffer, popping a
  larger tail letter when a smaller one arrives and the tail recurs later →
  see [`/dsa/monotonic_stack.md`](/dsa/monotonic_stack.md)
- **Greedy** — commit to the smallest feasible leading letter at each step →
  see [`/dsa/greedy.md`](/dsa/greedy.md)
- **Stack** — the buffer itself is a LIFO structure →
  see [`/dsa/stack.md`](/dsa/stack.md)
- **String Algorithms** — last-occurrence indexing and character bookkeeping →
  see [`/dsa/string_algorithms.md`](/dsa/string_algorithms.md)

---

## Approaches Overview

| # | Approach | Time | Space | When to use |
|---|----------|------|-------|-------------|
| 1 | Brute Force Recursion | O(26·n) | O(n) | Teaching the greedy "smallest first" insight |
| 2 | Greedy with Counts | O(n) | O(1) | When you prefer explicit remaining-count logic |
| 3 | Monotonic Stack + Last Index (Optimal) | O(n) | O(1) | Interview default — cleanest linear solution |

---

## Approach 1 — Brute Force Recursion

### Intuition
The first letter of the answer is the smallest letter we can pick such that all
other distinct letters still appear after our pick. If we look at the last
occurrence of each letter, the smallest of those last indices (`pos`) is a hard
right boundary: our first output letter must be chosen at or before `pos`, or
the letter that ends at `pos` would be lost. Within `s[0..pos]` we greedily take
the smallest letter, then recurse.

### Algorithm
1. If `s` is empty, return `""`.
2. Compute `last[c]` = final index of every letter present.
3. `pos` = the minimum value among all `last[c]`.
4. In window `s[0..pos]`, pick the smallest character; let `bestIdx` be its
   first index.
5. Emit that character; recurse on `s[bestIdx+1:]` with every copy of that
   character removed.

### Complexity
- **Time:** O(26·n) — at most 26 recursion levels (one per distinct letter),
  each doing an O(n) scan.
- **Space:** O(n) — recursion depth plus filtered substrings.

### Code
```go
func bruteForce(s string) string {
	if len(s) == 0 {
		return ""
	}
	last := map[byte]int{}
	for i := 0; i < len(s); i++ {
		last[s[i]] = i
	}
	pos := len(s)
	for _, idx := range last {
		if idx < pos {
			pos = idx
		}
	}
	best := byte('z' + 1)
	bestIdx := 0
	for i := 0; i <= pos; i++ {
		if s[i] < best {
			best = s[i]
			bestIdx = i
		}
	}
	rest := make([]byte, 0, len(s))
	for i := bestIdx + 1; i < len(s); i++ {
		if s[i] != best {
			rest = append(rest, s[i])
		}
	}
	return string(best) + bruteForce(string(rest))
}
```

### Dry Run
Input `s = "bcabc"`.

| Call | s              | last (b,c,a) | pos | window s[0..pos] | smallest | bestIdx | emit | rest (drop 'best') |
|------|----------------|--------------|-----|------------------|----------|---------|------|--------------------|
| 1    | `bcabc`        | b→3,c→4,a→2  | 2   | `bca`            | `a`      | 2       | `a`  | `bc`               |
| 2    | `bc`           | b→0,c→1      | 0   | `b`              | `b`      | 0       | `b`  | `c`                |
| 3    | `c`            | c→0          | 0   | `c`              | `c`      | 0       | `c`  | ``                 |
| 4    | ``             | —            | —   | —                | —        | —       | —    | return `""`        |

Result: `"a" + "b" + "c" = "abc"`. ✓

---

## Approach 2 — Greedy with Counts

### Intuition
Build the answer left to right in a buffer. Before adding character `c`, pop any
trailing buffer letter that is **larger** than `c` as long as it still occurs
later (remaining count > 0): demoting it after `c` produces a smaller string,
and a later copy will re-supply it. Never add a letter already in the buffer.

### Algorithm
1. `count[c]` = total occurrences of each letter.
2. `inResult[c]` = whether `c` is currently in the buffer.
3. For each char `c`: decrement `count[c]`.
   - If `c` already in buffer, continue.
   - While buffer top `> c` and `count[top] > 0`: pop it, clear its flag.
   - Push `c`, set its flag.
4. Return the buffer as a string.

### Complexity
- **Time:** O(n) — every character is pushed and popped at most once.
- **Space:** O(1) — two fixed 26-size arrays and a buffer of size ≤ 26.

### Code
```go
func greedyCounts(s string) string {
	var count [26]int
	var inResult [26]bool
	for i := 0; i < len(s); i++ {
		count[s[i]-'a']++
	}
	stack := make([]byte, 0, 26)
	for i := 0; i < len(s); i++ {
		c := s[i]
		count[c-'a']--
		if inResult[c-'a'] {
			continue
		}
		for len(stack) > 0 {
			top := stack[len(stack)-1]
			if top > c && count[top-'a'] > 0 {
				stack = stack[:len(stack)-1]
				inResult[top-'a'] = false
			} else {
				break
			}
		}
		stack = append(stack, c)
		inResult[c-'a'] = true
	}
	return string(stack)
}
```

### Dry Run
Input `s = "bcabc"`. Initial counts: b=2, c=2, a=1.

| i | c | count after -- | stack before | action | stack after |
|---|---|----------------|--------------|--------|-------------|
| 0 | b | b=1            | ``           | push b | `b`         |
| 1 | c | c=1            | `b`          | b<c, push c | `bc`   |
| 2 | a | a=0            | `bc`         | c>a & count[c]=1>0 pop c; b>a & count[b]=1>0 pop b; push a | `a` |
| 3 | b | b=0            | `a`          | a<b, push b | `ab`   |
| 4 | c | c=0            | `ab`         | b<c, push c | `abc`  |

Result: `"abc"`. ✓

---

## Approach 3 — Monotonic Stack + Last Index (Optimal)

### Intuition
Identical greedy logic to Approach 2, but rather than tracking remaining counts
we precompute `last[c]`, the final index of each letter. We may pop a larger
buffer top when a smaller letter arrives at index `i` precisely when the top
still appears later, i.e. `last[top] > i`. This avoids mutating counts and reads
cleanly as a monotonic-stack sweep.

### Algorithm
1. `last[c]` = final index of each letter.
2. `seen[c]` = buffer membership.
3. For `i, c` in `s`:
   - If `seen[c]`, skip.
   - While top `> c` and `last[top] > i`: pop, clear `seen`.
   - Push `c`, set `seen`.
4. Return the buffer.

### Complexity
- **Time:** O(n) — single pass, amortized O(1) per stack operation.
- **Space:** O(1) — fixed 26-size arrays plus a bounded buffer.

### Code
```go
func monotonicStack(s string) string {
	var last [26]int
	var seen [26]bool
	for i := 0; i < len(s); i++ {
		last[s[i]-'a'] = i
	}
	stack := make([]byte, 0, 26)
	for i := 0; i < len(s); i++ {
		c := s[i]
		if seen[c-'a'] {
			continue
		}
		for len(stack) > 0 {
			top := stack[len(stack)-1]
			if top > c && last[top-'a'] > i {
				stack = stack[:len(stack)-1]
				seen[top-'a'] = false
			} else {
				break
			}
		}
		stack = append(stack, c)
		seen[c-'a'] = true
	}
	return string(stack)
}
```

### Dry Run
Input `s = "bcabc"`. Last indices: b→3, c→4, a→2.

| i | c | seen[c]? | stack before | pops (top>c & last[top]>i) | stack after |
|---|---|----------|--------------|----------------------------|-------------|
| 0 | b | no       | ``           | none                       | `b`         |
| 1 | c | no       | `b`          | b<c, no pop                | `bc`        |
| 2 | a | no       | `bc`         | c>a & last[c]=4>2 pop c; b>a & last[b]=3>2 pop b | `a` |
| 3 | b | no       | `a`          | a<b, no pop                | `ab`        |
| 4 | c | no       | `ab`         | b<c, no pop                | `abc`       |

Result: `"abc"`. ✓

---

## Key Takeaways

- **"Smallest lexicographic subsequence keeping each element once"** is a
  canonical monotonic-stack pattern: sweep, and pop a bigger tail whenever a
  smaller element arrives that can safely demote it (it recurs later).
- The **`seen`/`inResult` membership guard** is essential — without it you would
  re-add letters and break the "each letter once" rule.
- **Last-occurrence index vs remaining count** are interchangeable ways to
  answer "does this letter appear again later?" — both give O(n)/O(1).
- The right boundary insight (`pos` = min of last indices) is the seed idea that
  makes the greedy provably correct.

---

## Related Problems

- LeetCode #1081 — Smallest Subsequence of Distinct Characters (identical)
- LeetCode #402 — Remove K Digits (monotonic stack, smallest number)
- LeetCode #321 — Create Maximum Number (greedy + monotonic stack)
- LeetCode #1673 — Find the Most Competitive Subsequence (monotonic stack)
