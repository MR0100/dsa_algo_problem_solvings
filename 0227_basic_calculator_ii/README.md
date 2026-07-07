# 0227 — Basic Calculator II

> LeetCode #227 · Difficulty: Medium
> **Categories:** Math, String, Stack

---

## Problem Statement

Given a string `s` which represents an expression, *evaluate this expression and return its value*.

The integer division should truncate toward zero.

You may assume that the given expression is always valid. All intermediate results will be in the range of `[-2³¹, 2³¹ - 1]`.

**Note:** You are not allowed to use any built-in function which evaluates strings as mathematical expressions, such as `eval()`.

**Example 1:**
```
Input: s = "3+2*2"
Output: 7
```

**Example 2:**
```
Input: s = " 3/2 "
Output: 1
```

**Example 3:**
```
Input: s = " 3+5 / 2 "
Output: 5
```

**Constraints:**
- `1 <= s.length <= 3 * 10⁵`
- `s` consists of integers and operators `('+', '-', '*', '/')` separated by some number of spaces.
- `s` represents a valid expression.
- All the integers in the expression are non-negative integers in the range `[0, 2³¹ - 1]`.
- The answer is guaranteed to fit in a 32-bit integer.

---

## Company Frequency

| Company    | Frequency        | Last Reported |
|------------|------------------|---------------|
| Amazon     | ★★★★★ Very High  | 2024          |
| Google     | ★★★★☆ High       | 2024          |
| Microsoft  | ★★★★☆ High       | 2024          |
| Meta       | ★★★☆☆ Medium     | 2023          |
| Bloomberg  | ★★★☆☆ Medium     | 2024          |
| Apple      | ★★☆☆☆ Low        | 2023          |

> ⚠️ Frequency data is crowd-sourced from LeetCode Discuss, Glassdoor, and
> community interview reports. Treat as a signal, not a guarantee.

---

## DSA Concepts Used

- **Stack** — deferring additive terms while resolving `*`/`/` immediately is a textbook stack use, respecting operator precedence in one pass → see [`/dsa/stack.md`](/dsa/stack.md)
- **String Parsing** — building multi-digit numbers character by character and skipping spaces → see [`/dsa/string_algorithms.md`](/dsa/string_algorithms.md)
- **Math (operator precedence)** — `*` and `/` bind tighter than `+` and `-`; truncating division toward zero → see [`/dsa/math_number_theory.md`](/dsa/math_number_theory.md)

---

## Approaches Overview

| # | Approach | Time | Space | When to use |
|---|----------|------|-------|-------------|
| 1 | Stack | O(n) | O(n) | Clear, easy to reason about precedence |
| 2 | O(1) Space (accumulators) | O(n) | O(1) | Optimal; only the stack top ever matters |

---

## Approach 1 — Stack

### Intuition
Multiplication and division bind tighter than addition and subtraction. Scan
left to right and remember the operator that *preceded* the current number.
Precedence is then resolved with a stack: on `+` push the number, on `-` push
its negation, and on `*`/`/` pop the top and push the combined value. Whatever
remains on the stack are independent additive terms — sum them for the answer.

### Algorithm
1. Track `prevOp` (operator before the current number), initialized to `'+'`.
2. Walk each character, building a multi-digit number `num`.
3. When an operator is reached, or at the end of the string, act on `prevOp`:
   - `'+'` → push `num`
   - `'-'` → push `-num`
   - `'*'` → replace stack top with `top * num`
   - `'/'` → replace stack top with `top / num` (Go truncates toward zero)
4. Reset `num` to 0 and set `prevOp` to the current operator.
5. Sum all stack entries.

### Complexity
- **Time:** O(n) — a single left-to-right pass; each number is pushed/popped O(1) times.
- **Space:** O(n) — the stack can hold one entry per additive term in the worst case.

### Code
```go
func stackEval(s string) int {
	stack := []int{} // holds resolved additive terms (already signed / multiplied)
	num := 0         // the integer currently being parsed
	prevOp := byte('+')

	for i := 0; i < len(s); i++ {
		c := s[i]
		if c >= '0' && c <= '9' {
			num = num*10 + int(c-'0') // extend the multi-digit number
		}
		// Act at an operator OR at the final character (flush the last number).
		if (c != ' ' && c < '0') || i == len(s)-1 {
			switch prevOp {
			case '+':
				stack = append(stack, num) // additive term, keep sign positive
			case '-':
				stack = append(stack, -num) // subtraction = adding a negative
			case '*':
				top := stack[len(stack)-1]      // multiply into the pending term
				stack[len(stack)-1] = top * num // replace top with product
			case '/':
				top := stack[len(stack)-1]      // integer-divide the pending term
				stack[len(stack)-1] = top / num // Go truncates toward zero, as required
			}
			prevOp = c // this operator governs the NEXT number
			num = 0    // reset the number accumulator
		}
	}

	sum := 0
	for _, v := range stack { // remaining entries are independent additive terms
		sum += v
	}
	return sum
}
```

### Dry Run
`s = "3+2*2"`:

| i | char | num | Triggered? | prevOp handled | Stack after | new prevOp |
|---|------|-----|------------|----------------|-------------|------------|
| 0 | 3    | 3   | no         | —              | []          | +          |
| 1 | +    | 3   | yes        | + → push 3     | [3]         | +          |
| 2 | 2    | 2   | no         | —              | [3]         | +          |
| 3 | *    | 2   | yes        | + → push 2     | [3, 2]      | *          |
| 4 | 2    | 2   | yes (end)  | * → 2*2=4      | [3, 4]      | *          |

Sum = 3 + 4 = **7**. ✅

---

## Approach 2 — O(1) Space (Running Accumulators)

### Intuition
The stack only ever needs its *top* element for `*` and `/`; everything below
is simply being summed. Replace the stack with two ints: `result` (sum of all
finished terms) and `lastNum` (the most recent term, still open to `*`/`/`).
When `+`/`-` arrives, the previous term is finalized, so fold `lastNum` into
`result`. For `*`/`/`, combine `lastNum` with the new number in place.

### Algorithm
1. Keep `result = 0`, `lastNum = 0`, `prevOp = '+'`.
2. Parse each number; at each operator or end of string, act on `prevOp`:
   - `'+'` → `result += lastNum`, then `lastNum = num`
   - `'-'` → `result += lastNum`, then `lastNum = -num`
   - `'*'` → `lastNum *= num`
   - `'/'` → `lastNum /= num`
3. After the loop, `result += lastNum` and return `result`.

### Complexity
- **Time:** O(n) — single pass over the string.
- **Space:** O(1) — three integer accumulators, no stack.

### Code
```go
func constantSpaceEval(s string) int {
	result := 0  // sum of all fully-committed additive terms
	lastNum := 0 // the current term, still open to * or /
	num := 0     // integer being parsed
	prevOp := byte('+')

	for i := 0; i < len(s); i++ {
		c := s[i]
		if c >= '0' && c <= '9' {
			num = num*10 + int(c-'0') // build the multi-digit number
		}
		if (c != ' ' && c < '0') || i == len(s)-1 {
			switch prevOp {
			case '+':
				result += lastNum // commit previous term
				lastNum = num     // new open term is +num
			case '-':
				result += lastNum // commit previous term
				lastNum = -num    // new open term is -num
			case '*':
				lastNum *= num // fold multiplication into the open term
			case '/':
				lastNum /= num // fold division into the open term (truncates toward zero)
			}
			prevOp = c
			num = 0
		}
	}
	result += lastNum // commit the final open term
	return result
}
```

### Dry Run
`s = "3+2*2"`:

| i | char | num | Triggered? | prevOp handled            | result | lastNum | new prevOp |
|---|------|-----|------------|---------------------------|--------|---------|------------|
| 0 | 3    | 3   | no         | —                         | 0      | 0       | +          |
| 1 | +    | 3   | yes        | + → result+=0, lastNum=3  | 0      | 3       | +          |
| 2 | 2    | 2   | no         | —                         | 0      | 3       | +          |
| 3 | *    | 2   | yes        | + → result+=3, lastNum=2  | 3      | 2       | *          |
| 4 | 2    | 2   | yes (end)  | * → lastNum=2*2           | 3      | 4       | *          |

Final: `result += lastNum` → 3 + 4 = **7**. ✅

---

## Key Takeaways
- Precedence with only `+ - * /` (no parentheses) is handled by deferring
  additive terms and eagerly resolving multiplicative ones against the last term.
- Recognizing that a stack is only ever touched at its top collapses O(n) space
  into O(1) with a `lastNum` accumulator — a broadly reusable optimization.
- Trigger the "flush" both on an operator *and* on the final character so the
  last number is never dropped.
- Go's integer division already truncates toward zero, matching the spec.

---

## Related Problems
- LeetCode #224 — Basic Calculator (adds parentheses, `+`/`-` only)
- LeetCode #772 — Basic Calculator III (parentheses + all four operators)
- LeetCode #150 — Evaluate Reverse Polish Notation (stack evaluation)
- LeetCode #394 — Decode String (stack-based parsing)
