# 0385 — Mini Parser

> LeetCode #385 · Difficulty: Medium
> **Categories:** String, Stack, Depth-First Search, Recursion, Design

---

## Problem Statement

Given a string `s` represents the serialization of a nested list, implement a parser to deserialize it and return *the deserialized* `NestedInteger`.

Each element is either an integer or a list whose elements may also be integers or other lists.

The `NestedInteger` interface is:

- `NestedInteger()` Initializes an empty nested list.
- `boolean isInteger()` Returns `true` if this `NestedInteger` holds a single integer, rather than a nested list.
- `int getInteger()` Returns the single integer that this `NestedInteger` holds, if it holds a single integer. Returns `null` if this `NestedInteger` holds a nested list.
- `void setInteger(int value)` Sets this `NestedInteger` to hold a single integer equal to `value`.
- `void add(NestedInteger ni)` Sets this `NestedInteger` to hold a nested list and adds a nested integer `ni` to it.
- `List<NestedInteger> getList()` Returns the nested list that this `NestedInteger` holds, if it holds a nested list. Returns an empty list if this `NestedInteger` holds a single integer.

**Example 1:**
```
Input: s = "324"
Output: 324
Explanation: You should return a NestedInteger object which contains a single integer 324.
```

**Example 2:**
```
Input: s = "[123,[456,[789]]]"
Output: [123,[456,[789]]]
Explanation: Return a NestedInteger object containing a nested list with 2 elements:
1. An integer containing value 123.
2. A nested list containing two elements:
    i.  An integer containing value 456.
    ii. A nested list with one element:
         a. An integer containing value 789
```

**Constraints:**
- `1 <= s.length <= 5 * 10⁴`
- `s` consists of digits, square brackets `"[]"`, negative sign `'-'`, and commas `','`.
- `s` is the serialization of valid `NestedInteger`.
- All the values in the input are in the range `[-10⁶, 10⁶]`.

---

## Company Frequency

| Company    | Frequency        | Last Reported |
|------------|------------------|---------------|
| Google     | ★★★☆☆ Medium     | 2024          |
| Amazon     | ★★★☆☆ Medium     | 2023          |
| Facebook   | ★★☆☆☆ Low        | 2023          |
| Microsoft  | ★★☆☆☆ Low        | 2023          |
| Bloomberg  | ★★☆☆☆ Low        | 2023          |

> ⚠️ Frequency data is crowd-sourced from LeetCode Discuss, Glassdoor, and
> community interview reports. Treat as a signal, not a guarantee.

---

## DSA Concepts Used

- **Stack** — brackets nest, so an explicit stack of open lists tracks the current nesting context as we scan → see [`/dsa/stack.md`](/dsa/stack.md)
- **Recursion / DFS (Recursive Descent)** — the grammar `element = int | '[' … ']'` maps directly to a recursive parser where the call stack replaces the explicit stack → see [`/dsa/graph_bfs_dfs.md`](/dsa/graph_bfs_dfs.md)
- **String Parsing / Tokenization** — single left-to-right pass, handling signed integer tokens, `[`, `]`, `,` → see [`/dsa/string_algorithms.md`](/dsa/string_algorithms.md)
- **Design / API Implementation** — build the `NestedInteger` type and use its `Add` / `SetInteger` operations → see [`/dsa/design_data_structures.md`](/dsa/design_data_structures.md)

---

## Approaches Overview

Let L = len(s), D = maximum nesting depth.

| # | Approach | Time | Space | When to use |
|---|----------|------|-------|-------------|
| 1 | Explicit Stack | O(L) | O(D) | Iterative preference; avoids recursion-depth limits |
| 2 | Recursive Descent (Optimal) | O(L) | O(D) | Cleanest; grammar-driven, least bookkeeping |

---

## Approach 1 — Explicit Stack

### Intuition
Brackets nest, so the natural tool is a stack. If the string does not start with `[`, it is a bare integer — return it directly. Otherwise each `[` opens a new list pushed on the stack; a number token becomes a `NestedInteger` added to the top list; each `]` closes the top list and nests it into the list below. The single remaining stack item is the answer.

### Algorithm
1. If `s[0] != '['`, return `NewInt(atoi(s))`.
2. Walk `s`, maintaining a stack of list `NestedInteger`s and a pending-number range.
3. On `[`: push a fresh empty list.
4. On `,` or `]`: flush any buffered number as an integer into the top list; on `]` additionally pop the top list and `Add` it into the new top (unless it is the root).
5. On a digit or `-`: extend the current number token.
6. Return the single remaining list.

### Complexity
- **Time:** O(L) — one pass; each `atoi` is amortized over the digits it consumes.
- **Space:** O(D) — stack depth equals maximum nesting; O(L) worst case.

### Code
```go
func stackParse(s string) *NestedInteger {
	if s[0] != '[' { // a bare integer like "324" or "-7"
		v, _ := strconv.Atoi(s)
		return NewInt(v)
	}

	stack := []*NestedInteger{} // open lists, innermost on top
	numStart := -1              // start index of current number token, -1 = none

	flushNumber := func(end int) {
		if numStart != -1 {
			v, _ := strconv.Atoi(s[numStart:end])
			top := stack[len(stack)-1]
			top.Add(NewInt(v)) // nest the integer into the current list
			numStart = -1
		}
	}

	for i := 0; i < len(s); i++ {
		switch c := s[i]; {
		case c == '[':
			stack = append(stack, NewList()) // open a new list
		case c == ',':
			flushNumber(i) // a number token (if any) just ended
		case c == ']':
			flushNumber(i) // close out any trailing number in this list
			if len(stack) > 1 {
				top := stack[len(stack)-1]
				stack = stack[:len(stack)-1]  // pop the finished list
				stack[len(stack)-1].Add(top)  // nest it into its parent
			}
		default: // digit or leading '-'
			if numStart == -1 {
				numStart = i // begin a new number token
			}
		}
	}
	return stack[0] // the root list
}
```

### Dry Run
`s = "324"` → `s[0] != '['`, so return `NewInt(324)`. Output `324`.

Now the richer trace for `s = "[123,[456,[789]]]"`:

| i | char | action | stack (innermost last) |
|---|------|--------|------------------------|
| 0 | `[` | push list A | `[A]` |
| 1-3 | `123` | build number token | `[A]` |
| 4 | `,` | flush 123 → A.Add(123) | `[A(123)]` |
| 5 | `[` | push list B | `[A(123), B]` |
| 6-8 | `456` | number token | `[A, B]` |
| 9 | `,` | flush 456 → B.Add(456) | `[A, B(456)]` |
| 10 | `[` | push list C | `[A, B, C]` |
| 11-13 | `789` | number token | `[A, B, C]` |
| 14 | `]` | flush 789 → C.Add(789); pop C; B.Add(C) | `[A, B(456,[789])]` |
| 15 | `]` | pop B; A.Add(B) | `[A(123,[456,[789]])]` |
| 16 | `]` | len(stack)==1, nothing to pop | `[A]` |
| end | — | return A | `[123,[456,[789]]]` |

---

## Approach 2 — Recursive Descent (Optimal)

### Intuition
The data is a grammar, so parse it recursively. A shared cursor `i` walks the string. `parseValue` inspects `s[i]`: if not `[`, read a (possibly negative) integer token; otherwise consume `[`, then repeatedly `parseValue` each comma-separated child until `]`. Recursion naturally matches nesting depth — the call stack *is* the stack.

### Algorithm
1. Keep an index `i` shared across recursive calls (closure variable).
2. `parseValue`: if `s[i] != '['`, scan an optional `-` then digits, wrap in `NewInt`.
3. Else consume `[`; make an empty list; while `s[i] != ']'`, `parseValue` a child, `Add` it, skip a `,` if present; consume `]`; return the list.
4. Call `parseValue()` once from index 0.

### Complexity
- **Time:** O(L) — every character consumed exactly once across the recursion.
- **Space:** O(D) — recursion depth = maximum nesting; O(L) worst case.

### Code
```go
func recursiveParse(s string) *NestedInteger {
	i := 0 // shared cursor into s

	var parseValue func() *NestedInteger
	parseValue = func() *NestedInteger {
		if s[i] != '[' { // integer token: optional '-' then digits
			start := i
			if s[i] == '-' {
				i++
			}
			for i < len(s) && s[i] >= '0' && s[i] <= '9' {
				i++
			}
			v, _ := strconv.Atoi(s[start:i])
			return NewInt(v)
		}

		i++ // consume '['
		lst := NewList()
		for s[i] != ']' { // parse children until the matching ']'
			lst.Add(parseValue()) // recurse for each element
			if s[i] == ',' {
				i++ // skip the separator
			}
		}
		i++ // consume ']'
		return lst
	}

	return parseValue()
}
```

### Dry Run
`s = "324"`.

| call | i at entry | s[i] | branch | action | return |
|------|-----------|------|--------|--------|--------|
| parseValue | 0 | `3` | integer | scan `324`, i→3 | `NewInt(324)` |

Output `324`.

Nested trace for `s = "[123,[456,[789]]]"`: outer call sees `[`, consumes it, then loops — `parseValue` reads `123`, skips `,`, recurses into `[456,[789]]` which itself reads `456`, skips `,`, recurses into `[789]` reading `789`; each inner list returns and is `Add`ed to its parent, yielding `[123,[456,[789]]]`.

---

## Key Takeaways
- **Nested-bracket parsing** has two canonical shapes: an explicit stack of open containers, or recursive descent driven by the grammar. They are duals — the recursion's call stack mirrors the explicit stack.
- Handle the **bare-integer** base case first (`s[0] != '['`), and remember integers can be **negative** (`-`).
- A shared **cursor variable** (closure or pointer) is the clean way to thread parse position through recursion.
- Building the `String()` method on `NestedInteger` lets you round-trip and verify the parse against the input.

---

## Related Problems
- LeetCode #341 — Flatten Nested List Iterator (consume a NestedInteger)
- LeetCode #339 — Nested List Weight Sum (DFS over NestedInteger)
- LeetCode #394 — Decode String (stack-based nested parsing)
- LeetCode #726 — Number of Atoms (recursive/stack parsing of nested formulas)
