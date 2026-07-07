# 0394 — Decode String

> LeetCode #394 · Difficulty: Medium
> **Categories:** Stack, String, Recursion

---

## Problem Statement

Given an encoded string, return its decoded string.

The encoding rule is: `k[encoded_string]`, where the `encoded_string` inside the square
brackets is being repeated exactly `k` times. Note that `k` is guaranteed to be a positive
integer.

You may assume that the input string is always valid; there are no extra white spaces,
square brackets are well-formed, etc. Furthermore, you may assume that the original data
does not contain any digits and that digits are only for those repeat numbers, `k`. For
example, there will not be input like `3a` or `2[4]`.

The test cases are generated so that the length of the output will never exceed `10^5`.

**Example 1:**

```
Input: s = "3[a]2[bc]"
Output: "aaabcbc"
```

**Example 2:**

```
Input: s = "3[a2[c]]"
Output: "accaccacc"
```

**Example 3:**

```
Input: s = "2[abc]3[cd]ef"
Output: "abcabccdcdcdef"
```

**Constraints:**

- `1 <= s.length <= 30`
- `s` consists of lowercase English letters, digits, and square brackets `'[]'`.
- `s` is guaranteed to be a valid input.
- All the integers in `s` are in the range `[1, 300]`.

---

## Company Frequency

| Company   | Frequency         | Last Reported |
|-----------|-------------------|---------------|
| Google    | ★★★★★ Very High   | 2024          |
| Amazon    | ★★★★☆ High        | 2024          |
| Microsoft | ★★★★☆ High        | 2024          |
| Facebook  | ★★★☆☆ Medium      | 2023          |
| Bloomberg | ★★★☆☆ Medium      | 2023          |

> ⚠️ Frequency data is crowd-sourced from LeetCode Discuss, Glassdoor, and
> community interview reports. Treat as a signal, not a guarantee.

---

## DSA Concepts Used

- **Stack** — push outer context on `[`, pop and combine on `]` for nesting → see [`/dsa/stack.md`](/dsa/stack.md)
- **Recursion / recursive descent** — a grammar-driven parser mirroring `k[ ... ]` → see [`/dsa/divide_and_conquer.md`](/dsa/divide_and_conquer.md)
- **String building** — repeat and concatenate decoded fragments → see [`/dsa/string_algorithms.md`](/dsa/string_algorithms.md)

---

## Approaches Overview

| # | Approach | Time | Space | When to use |
|---|----------|------|-------|-------------|
| 1 | Two Stacks (iterative) | O(N) | O(N) | Avoids recursion; classic interview answer |
| 2 | Recursive Descent | O(N) | O(N) | Cleanest mapping to the grammar |

(N = total length of the decoded output.)

---

## Approach 1 — Two Stacks

### Intuition

Nesting screams stack. Hold the string built at the current level (`cur`) and the count
being parsed (`num`). On `[` we descend: push `num` and `cur`, then reset both. On `]` we
finish an inner level: pop the repeat count `k` and the outer prefix, and set
`cur = prefix + cur*k`.

### Algorithm

1. `countStack`, `stringStack`, `cur = ""`, `num = 0`.
2. Scan each char: digit → `num = num*10 + d`; `[` → push `num`, push `cur`, reset both;
   `]` → pop `k`, pop `prev`, `cur = prev + repeat(cur, k)`; letter → `cur += ch`.
3. Return `cur`.

### Complexity

- **Time:** O(N) — each output character is produced once.
- **Space:** O(N) — stacks plus the building string.

### Code

```go
func twoStacks(s string) string {
	countStack := []int{}
	stringStack := []string{}
	cur := ""
	num := 0

	for i := 0; i < len(s); i++ {
		ch := s[i]
		switch {
		case ch >= '0' && ch <= '9':
			num = num*10 + int(ch-'0')
		case ch == '[':
			countStack = append(countStack, num)
			stringStack = append(stringStack, cur)
			num = 0
			cur = ""
		case ch == ']':
			k := countStack[len(countStack)-1]
			countStack = countStack[:len(countStack)-1]
			prev := stringStack[len(stringStack)-1]
			stringStack = stringStack[:len(stringStack)-1]
			cur = prev + strings.Repeat(cur, k)
		default:
			cur += string(ch)
		}
	}
	return cur
}
```

### Dry Run

`s = "3[a]2[bc]"`:

| char | num | cur | countStack | stringStack | note |
|------|-----|-----|------------|-------------|------|
| 3 | 3 | "" | [] | [] | build number |
| [ | 0 | "" | [3] | [""] | descend |
| a | 0 | "a" | [3] | [""] | letter |
| ] | 0 | "aaa" | [] | [] | pop k=3: ""+"a"*3 |
| 2 | 2 | "aaa" | [] | [] | build number |
| [ | 0 | "" | [2] | ["aaa"] | descend |
| b | 0 | "b" | [2] | ["aaa"] | letter |
| c | 0 | "bc" | [2] | ["aaa"] | letter |
| ] | 0 | "aaabcbc" | [] | [] | pop k=2: "aaa"+"bc"*2 |

Result **`"aaabcbc"`**.

---

## Approach 2 — Recursive Descent

### Intuition

The grammar is a sequence of `letters` or `number '[' sequence ']'`. A function that
decodes from a shared cursor and returns at `]` naturally handles nesting: on a digit read
the count, consume `[`, recurse for the inside, consume `]`, then repeat the inner result
`k` times.

### Algorithm

1. Shared index `i` into `s`.
2. `decode()` loops while `i < len` and `s[i] != ']'`: letters append; a digit run parses
   `k`, skips `[`, recurses, skips `]`, appends inner repeated `k` times.
3. The top-level call returns the whole decoding.

### Complexity

- **Time:** O(N) — each output char written once.
- **Space:** O(N + depth) — output builder plus recursion stack (nesting depth).

### Code

```go
func recursiveDescent(s string) string {
	i := 0
	var decode func() string
	decode = func() string {
		var sb strings.Builder
		for i < len(s) && s[i] != ']' {
			ch := s[i]
			if ch >= '0' && ch <= '9' {
				k := 0
				for i < len(s) && s[i] >= '0' && s[i] <= '9' {
					k = k*10 + int(s[i]-'0')
					i++
				}
				i++
				inner := decode()
				i++
				for r := 0; r < k; r++ {
					sb.WriteString(inner)
				}
			} else {
				sb.WriteByte(ch)
				i++
			}
		}
		return sb.String()
	}
	return decode()
}
```

### Dry Run

`s = "3[a]2[bc]"` — top-level `decode()`:

| i | char | action | sb so far |
|---|------|--------|-----------|
| 0 | 3 | read k=3; i→1; skip `[` i→2; recurse | "" |
| 2 | a | inner decode returns "a"; i→3 skip `]` i→4 | append "a"×3 → "aaa" |
| 4 | 2 | read k=2; i→5; skip `[` i→6; recurse | "aaa" |
| 6 | b,c | inner decode returns "bc"; i→8 skip `]` i→9 | append "bc"×2 → "aaabcbc" |

Result **`"aaabcbc"`**.

---

## Key Takeaways

- **Two-stack pattern** for bracketed/nested structures: one stack for the pending
  multiplier (or operator), one for the string (or value) built before the bracket.
- On `[` push-and-reset; on `]` pop-and-combine — the mirror-image bracket handling is the
  reusable trick (see also basic calculators).
- Recursive descent gives the shortest code when you keep a **shared cursor** the recursion
  advances.
- Multi-digit counts: accumulate `num = num*10 + digit`.

---

## Related Problems

- LeetCode #726 — Number of Atoms (nested multipliers, same stack pattern)
- LeetCode #224 — Basic Calculator (stack for nested parentheses/signs)
- LeetCode #227 — Basic Calculator II (operator/operand stacks)
- LeetCode #1190 — Reverse Substrings Between Each Pair of Parentheses (stack of strings)
