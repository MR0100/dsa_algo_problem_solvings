# 0439 — Ternary Expression Parser

> LeetCode #439 · Difficulty: Medium
> **Categories:** String, Stack, Recursion

---

## Problem Statement

Given a string `expression` representing arbitrarily nested ternary expressions, evaluate the expression, and return *the result of it*.

You can always assume that the given expression is valid and only contains digits, `'?'`, `':'`, `'T'`, and `'F'`, where `'T'` is true and `'F'` is false. All the numbers in the expression are **one-digit** numbers (i.e., in the range `[0, 9]`).

The conditional expressions group **right-to-left** (as usual in most languages), and the result of the expression will always evaluate to either a digit, `'T'` or `'F'`.

**Example 1:**

```
Input: expression = "T?2:3"
Output: "2"
Explanation: If true, then result is 2; otherwise result is 3.
```

**Example 2:**

```
Input: expression = "F?1:T?4:5"
Output: "4"
Explanation: The conditional expressions group right-to-left. Using parenthesis, it is read/evaluated as:
"(F ? 1 : (T ? 4 : 5))"   -->   "(F ? 1 : 4)"   -->   "4"
or   "(F ? 1 : (T ? 4 : 5))"   -->   "(T ? 4 : 5)"   -->   "4"
```

**Example 3:**

```
Input: expression = "T?T?F:5:3"
Output: "F"
Explanation: The conditional expressions group right-to-left. Using parenthesis, it is read/evaluated as:
"(T ? (T ? F : 5) : 3)"   -->   "(T ? F : 3)"   -->   "F"
"(T ? (T ? F : 5) : 3)"   -->   "(T ? F : 5)"   -->   "F"
```

**Constraints:**

- `5 <= expression.length <= 10^4`
- `expression` consists of digits, `'T'`, `'F'`, `'?'`, and `':'`.
- It is **guaranteed** that `expression` is a valid ternary expression and that each number is a **one-digit number**.

---

## Company Frequency

| Company    | Frequency        | Last Reported |
|------------|------------------|---------------|
| Facebook   | ★★★☆☆ Medium     | 2023          |
| Snapchat   | ★★★☆☆ Medium     | 2022          |
| Amazon     | ★★☆☆☆ Low        | 2022          |
| Google     | ★★☆☆☆ Low        | 2022          |
| Microsoft  | ★★☆☆☆ Low        | 2021          |

> ⚠️ Frequency data is crowd-sourced from LeetCode Discuss, Glassdoor, and
> community interview reports. Treat as a signal, not a guarantee.

---

## DSA Concepts Used

- **Stack** — scanning right-to-left and pushing tokens lets each `?` collapse the two operands already sitting on top; the innermost (rightmost) ternary resolves first, matching the right-associative grouping → see [`/dsa/stack.md`](/dsa/stack.md)
- **String Parsing / Recursive Descent** — the expression follows a tiny grammar `expr = value | value '?' expr ':' expr`; parsing it with a moving cursor is textbook recursive descent → see [`/dsa/string_algorithms.md`](/dsa/string_algorithms.md)

---

## Approaches Overview

| # | Approach | Time | Space | When to use |
|---|----------|------|-------|-------------|
| 1 | Right-to-Left Stack | O(n) | O(n) | Cleanest to reason about; leans on right-associativity directly |
| 2 | Recursive Descent | O(n) | O(n) recursion | Most readable as a grammar; natural if you think in parsers |
| 3 | Iterative Forward Skip (Optimal) | O(n) | O(1) | No stack, no recursion — keep taken branch, skip the other |

---

## Approach 1 — Right-to-Left Stack

### Intuition

Ternaries group **right-to-left**, so the rightmost `?` is the innermost complete expression. Walk the string from the end pushing characters onto a stack. When you hit a `?`, the two operands it needs are already resolved and sitting on top: pop the true-branch, pop the `:` separator, pop the false-branch. The character immediately to the left of the `?` is the condition — push back the winning branch. Continue leftward; by index 0 the stack holds exactly one value.

### Algorithm

1. Iterate `i` from `len-1` down to `0`.
2. If `s[i] == '?'`: pop `trueBranch`, pop `':'`, pop `falseBranch`; read condition `s[i-1]`; push `trueBranch` if condition is `'T'` else `falseBranch`; decrement `i` once more to consume the condition.
3. Otherwise push `s[i]`.
4. Return the single remaining stack element.

### Complexity

- **Time:** O(n) — each character is pushed and popped a constant number of times.
- **Space:** O(n) — the stack, worst case proportional to the input.

### Code

```go
func stackRightToLeft(expression string) string {
	stack := []byte{} // holds resolved single characters (values / ':' markers)
	for i := len(expression) - 1; i >= 0; i-- {
		c := expression[i]
		if c == '?' {
			trueBranch := stack[len(stack)-1]  // top = value if condition true
			stack = stack[:len(stack)-1]       // pop true-branch
			stack = stack[:len(stack)-1]       // pop the ':' separator
			falseBranch := stack[len(stack)-1] // next = value if condition false
			stack = stack[:len(stack)-1]       // pop false-branch
			cond := expression[i-1]            // the condition sits just left of '?'
			if cond == 'T' {
				stack = append(stack, trueBranch) // condition true → keep true-branch
			} else {
				stack = append(stack, falseBranch) // condition false → keep false-branch
			}
			i-- // we consumed the condition character too; skip it
		} else {
			stack = append(stack, c) // digit, 'T', 'F', or ':' — defer resolution
		}
	}
	return string(stack[0]) // exactly one value remains
}
```

### Dry Run

Example 1: `expression = "T?2:3"`. Indices: 0=`T`, 1=`?`, 2=`2`, 3=`:`, 4=`3`. Scan right-to-left.

| i | char | action | stack (top → right) |
|---|------|--------|---------------------|
| 4 | `3` | push | `[3]` |
| 3 | `:` | push | `[3, :]` |
| 2 | `2` | push | `[3, :, 2]` |
| 1 | `?` | pop true=`2`, pop `:`, pop false=`3`; cond=`s[0]='T'` → push `2`; i-- | `[2]` |
| 0 | (consumed as condition) | — | `[2]` |

Return `"2"` ✔

---

## Approach 2 — Recursive Descent

### Intuition

The expression obeys the grammar `expr = value | value '?' expr ':' expr`. Read the leading `value`; if the next character is `?`, it is a ternary — recursively parse the true-branch and the false-branch. Crucially, we parse **both** branches even though only one value is kept, because parsing advances the shared cursor to the correct resume point for the caller. Right-associativity falls out for free: the recursion for the true-branch consumes its own nested ternary before the `:` is reached.

### Algorithm

1. Maintain a shared cursor `pos`.
2. `parse()`: read `cond = s[pos]`, advance. If the next char is not `'?'`, return `cond` (a bare value).
3. Otherwise skip `'?'`; `trueVal = parse()`; skip `':'`; `falseVal = parse()`; return `trueVal` if `cond == 'T'` else `falseVal`.

### Complexity

- **Time:** O(n) — the cursor advances through each character exactly once.
- **Space:** O(n) — recursion stack up to the expression's nesting depth.

### Code

```go
func recursiveDescent(expression string) string {
	pos := 0 // shared cursor into expression

	var parse func() string
	parse = func() string {
		cond := expression[pos] // first atom: a value or a condition (T/F)
		pos++                   // move past it
		// A bare value (no trailing '?') is a complete expression.
		if pos >= len(expression) || expression[pos] != '?' {
			return string(cond)
		}
		pos++               // skip '?'
		trueVal := parse()  // recursively evaluate the true-branch
		pos++               // skip ':'
		falseVal := parse() // recursively evaluate the false-branch
		if cond == 'T' {    // choose based on the condition
			return trueVal
		}
		return falseVal
	}

	return parse()
}
```

### Dry Run

Example 1: `expression = "T?2:3"`.

| call depth | pos on entry | cond | next char | action | returns |
|------------|--------------|------|-----------|--------|---------|
| parse #1 | 0 | `T` | `s[1]='?'` → ternary | skip `?` (pos→2); call #2 | (waits) |
| parse #2 | 2 | `2` | `s[3]=':'` ≠ `?` | bare value | `"2"` |
| back in #1 | pos=3 | — | skip `:` (pos→4); call #3 | (waits) |
| parse #3 | 4 | `3` | pos=5 ≥ len | bare value | `"3"` |
| back in #1 | — | `T` | cond=='T' → keep trueVal | `"2"` |

Return `"2"` ✔

---

## Approach 3 — Iterative Forward Skip (Optimal)

### Intuition

Neither a stack nor recursion is strictly required. Move a single pointer forward, always keeping the **taken** branch and fast-forwarding past the **untaken** one. Read the leading condition: if `'T'`, step straight into the true-branch (right after `"T?"`); if `'F'`, discard the whole true-branch and jump to just after the `:` that matches this `?`. Locating that matching `:` means skipping any nested ternaries, which a depth counter handles: `+1` on `?`, `−1` on `:`, and the matching separator is the `:` seen at depth 0. The pointer always lands on the head of the next value we still care about, so the value it finally rests on is the answer — all in O(1) extra space.

### Algorithm

1. `i = 0`.
2. If the character after `s[i]` is not `'?'`, then `s[i]` is the final value — return it.
3. Otherwise `s[i]` is a condition:
   - `'T'` → set `i += 2` (descend into the true-branch).
   - `'F'` → set `i += 2`, then scan forward with a depth counter to the matching `:`, then step past it onto the false-branch.
4. Repeat.

### Complexity

- **Time:** O(n) — the pointer only ever moves forward.
- **Space:** O(1) — two integer counters; no stack, no recursion.

### Code

```go
func iterativeForwardSkip(expression string) string {
	i := 0
	for {
		// A value with no following '?' is a complete (sub-)expression's result.
		if i+1 >= len(expression) || expression[i+1] != '?' {
			return string(expression[i])
		}
		if expression[i] == 'T' {
			// Condition true → descend into the true-branch (right after "T?").
			i += 2
		} else {
			// Condition false → skip the ENTIRE true-branch, land on false-branch.
			// Start just past "F?" and walk to the ':' that matches this '?'.
			i += 2       // move past "F?" to the first char of the true-branch
			depth := 0   // nesting level relative to this branch
			for depth > 0 || expression[i] != ':' {
				switch expression[i] {
				case '?':
					depth++ // entering a nested ternary
				case ':':
					depth-- // leaving a nested ternary
				}
				i++
			}
			i++ // step over the matching ':' onto the false-branch head
		}
	}
}
```

### Dry Run

Example 1: `expression = "T?2:3"`.

| i | s[i] | s[i+1] | decision | i after |
|---|------|--------|----------|---------|
| 0 | `T` | `?` | condition true → `i += 2` | 2 |
| 2 | `2` | `:` (≠ `?`) | bare value → return `"2"` | — |

Return `"2"` ✔

Nested check — Example 3 `"T?T?F:5:3"` (0=`T`,1=`?`,2=`T`,3=`?`,4=`F`,5=`:`,6=`5`,7=`:`,8=`3`):

| i | s[i] | s[i+1] | decision | i after |
|---|------|--------|----------|---------|
| 0 | `T` | `?` | true → `i += 2` | 2 |
| 2 | `T` | `?` | true → `i += 2` | 4 |
| 4 | `F` | `:` (≠ `?`) | bare value → return `"F"` | — |

Return `"F"` ✔ (a false condition would instead have skipped its true-branch via the depth counter).

---

## Key Takeaways

- **Right-associative parsing loves a right-to-left stack.** When the innermost expression sits at the *end*, scanning backwards makes each operator find its operands already resolved on top of the stack.
- **Recursive descent mirrors the grammar.** For `expr = value | value '?' expr ':' expr`, the parser writes itself — and a shared mutable cursor keeps both branches advancing the same position.
- **You can trade the stack for a forward skip** when you only need one value out: keep the taken branch, and skip the untaken branch with a balanced depth counter (`?` opens, `:` closes). O(1) memory.
- **Single-digit values** and a *guaranteed-valid* expression remove all tokenising pain — every operand is exactly one character, so no number assembly or error handling is needed.

---

## Related Problems

- LeetCode #227 — Basic Calculator II (stack-based expression evaluation)
- LeetCode #224 — Basic Calculator (parentheses + recursion/stack)
- LeetCode #772 — Basic Calculator III (nested operators and parentheses)
- LeetCode #150 — Evaluate Reverse Polish Notation (operand stack)
- LeetCode #394 — Decode String (right-to-left / stack nesting)
