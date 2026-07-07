# 0150 — Evaluate Reverse Polish Notation

> LeetCode #150 · Difficulty: Medium
> **Categories:** Array, Math, Stack

---

## Problem Statement

You are given an array of strings `tokens` that represents an arithmetic expression in a [Reverse Polish Notation](http://en.wikipedia.org/wiki/Reverse_Polish_notation).

Evaluate the expression. Return an integer that represents the value of the expression.

**Note** that:

- The valid operators are `'+'`, `'-'`, `'*'`, and `'/'`.
- Each operand may be an integer or another expression.
- The division between two integers always **truncates toward zero**.
- There will not be any division by zero.
- The input represents a valid arithmetic expression in a reverse polish notation.
- The answer and all the intermediate calculations can be represented in a **32-bit** integer.

**Example 1:**
```
Input: tokens = ["2","1","+","3","*"]
Output: 9
Explanation: ((2 + 1) * 3) = 9
```

**Example 2:**
```
Input: tokens = ["4","13","5","/","+"]
Output: 6
Explanation: (4 + (13 / 5)) = 6
```

**Example 3:**
```
Input: tokens = ["10","6","9","3","+","-11","*","/","*","17","+","5","+"]
Output: 22
Explanation: ((10 * (6 / ((9 + 3) * -11))) + 17) + 5
= ((10 * (6 / (12 * -11))) + 17) + 5
= ((10 * (6 / -132)) + 17) + 5
= ((10 * 0) + 17) + 5
= (0 + 17) + 5
= 17 + 5
= 22
```

**Constraints:**
- `1 <= tokens.length <= 10^4`
- `tokens[i]` is either an operator: `"+"`, `"-"`, `"*"`, or `"/"`, or an integer in the range `[-200, 200]`.

---

## Company Frequency

| Company    | Frequency        | Last Reported |
|------------|------------------|---------------|
| Amazon     | ★★★★☆ High       | 2024          |
| LinkedIn   | ★★★★☆ High       | 2024          |
| Facebook   | ★★★★☆ High       | 2024          |
| Microsoft  | ★★★☆☆ Medium     | 2024          |
| Google     | ★★★☆☆ Medium     | 2023          |
| Uber       | ★★☆☆☆ Low        | 2023          |

> ⚠️ Frequency data is crowd-sourced from LeetCode Discuss, Glassdoor, and
> community interview reports. Treat as a signal, not a guarantee.

---

## DSA Concepts Used

- **Stack** — postfix notation is *defined* for stack evaluation: operands wait on the stack until their operator arrives → see [`/dsa/stack.md`](/dsa/stack.md)
- **Expression trees / recursion** — the last RPN token is the root of the expression tree; evaluating right-to-left is a post-order-in-reverse walk → see [`/dsa/tree_traversal.md`](/dsa/tree_traversal.md)
- **Integer division semantics** — "truncate toward zero" (Go's native `/` behaviour) vs floor division; a classic cross-language trap → see [`/dsa/math_number_theory.md`](/dsa/math_number_theory.md)

---

## Approaches Overview

| # | Approach | Time | Space | When to use |
|---|----------|------|-------|-------------|
| 1 | Brute Force (Repeated Scan and Reduce) | O(n²) | O(n) | Mirrors hand-evaluation; fine for tiny inputs only |
| 2 | Recursion from the Right | O(n) | O(n) stack | Elegant; shows the expression-tree view |
| 3 | Stack (Optimal) | O(n) | O(n) | The intended answer — one pass, no recursion |

---

## Approach 1 — Brute Force (Repeated Scan and Reduce)

### Intuition
Evaluate the way you would on paper: find an operator that directly follows two plain numbers — such a triple `number number op` is always evaluable — replace those three tokens with the computed value, and repeat. Each reduction shrinks the expression by two tokens but keeps it valid RPN, so after ~n/2 rounds a single number remains: the answer.

### Algorithm
1. Copy `tokens` into a working slice (we splice it in place).
2. While more than one token remains:
   1. Scan left→right for the first index `i` where `work[i]` is an operator and `work[i-2]`, `work[i-1]` are both numbers.
   2. Compute `apply(work[i-2], work[i-1], work[i])`.
   3. Overwrite `work[i-2]` with the result string and splice out positions `i-1` and `i`.
3. Parse the single remaining token and return it.

### Complexity
- **Time:** O(n²) — ~n/2 reductions, each rescanning and shifting up to O(n) tokens.
- **Space:** O(n) — the mutable working copy.

### Code
```go
func bruteForceReduce(tokens []string) int {
	work := make([]string, len(tokens))
	copy(work, tokens) // never mutate the caller's slice

	for len(work) > 1 {
		for i := 2; i < len(work); i++ {
			if isOperator(work[i]) && !isOperator(work[i-2]) && !isOperator(work[i-1]) {
				a, _ := strconv.Atoi(work[i-2])
				b, _ := strconv.Atoi(work[i-1])
				result := apply(a, b, work[i])
				work[i-2] = strconv.Itoa(result)
				work = append(work[:i-1], work[i+1:]...) // splice 3 tokens → 1
				break
			}
		}
	}
	value, _ := strconv.Atoi(work[0])
	return value
}
```

### Dry Run
Example 1: `tokens = ["2","1","+","3","*"]`.

| Round | Working expression        | First reducible triple found | Reduction    | Expression after   |
|-------|---------------------------|-------------------------------|--------------|--------------------|
| 1     | `["2","1","+","3","*"]`   | `2 1 +` at i = 2              | 2 + 1 = 3    | `["3","3","*"]`    |
| 2     | `["3","3","*"]`           | `3 3 *` at i = 2              | 3 * 3 = 9    | `["9"]`            |
| —     | `["9"]`                   | single token → stop           | —            | returns **9**      |

Output: `9` ✓

---

## Approach 2 — Recursion from the Right

### Intuition
An RPN expression is a **post-order serialization of an expression tree** — so its *last* token is the tree's root. If that token is an operator, everything before it is: (left subtree tokens)(right subtree tokens). Reading backwards with a cursor, the right operand's subtree comes first, then the left's. So: consume the token at the cursor; if it's an operator, recursively evaluate the **right** operand, then the **left**, and combine. Each call consumes exactly its subtree's tokens — no boundary computation needed.

### Algorithm
1. `idx := len(tokens) - 1` — shared cursor, moves right→left.
2. `eval()`:
   1. `token := tokens[idx]`; `idx--`.
   2. If `token` is an operator: `right := eval()`, then `left := eval()`, return `apply(left, right, token)`.
   3. Otherwise return the parsed integer (a leaf).
3. The top-level `eval()` returns the whole expression's value.

### Complexity
- **Time:** O(n) — the cursor visits each token exactly once.
- **Space:** O(n) — recursion depth = expression-tree height; worst case (e.g. `1 1 + 1 + 1 + …`) is O(n).

### Code
```go
func recursiveEval(tokens []string) int {
	idx := len(tokens) - 1 // shared cursor, consumed right-to-left
	var eval func() int
	eval = func() int {
		token := tokens[idx]
		idx--
		if isOperator(token) {
			right := eval() // right subtree sits immediately before the operator
			left := eval()
			return apply(left, right, token)
		}
		value, _ := strconv.Atoi(token)
		return value
	}
	return eval()
}
```

### Dry Run
Example 1: `tokens = ["2","1","+","3","*"]` (indices 0–4).

| Step | Call depth | Token consumed (idx) | Action                       | Returns |
|------|------------|----------------------|------------------------------|---------|
| 1    | eval₀      | `*` (4)              | operator → need right, left  | …       |
| 2    | eval₁ (right of `*`) | `3` (3)    | number leaf                  | 3       |
| 3    | eval₁ (left of `*`)  | `+` (2)    | operator → need right, left  | …       |
| 4    | eval₂ (right of `+`) | `1` (1)    | number leaf                  | 1       |
| 5    | eval₂ (left of `+`)  | `2` (0)    | number leaf                  | 2       |
| 6    | eval₁ resumes        | —          | apply(2, 1, `+`) = 3         | 3       |
| 7    | eval₀ resumes        | —          | apply(3, 3, `*`) = 9         | **9**   |

Output: `9` ✓

---

## Approach 3 — Stack (Optimal)

### Intuition
RPN exists *because of* stacks — it's the machine format of HP calculators and the JVM's bytecode operand model. Scan left to right: a number is pushed; an operator pops the top two values, applies itself, pushes the result. Postfix order guarantees an operator's two operands are always the two most recently produced values, i.e. exactly the top of the stack. One subtlety: the **top of the stack is the RIGHT operand** — popping in the wrong order breaks `-` and `/`.

### Algorithm
1. Create an int stack (a Go slice with capacity `len(tokens)`).
2. For each token:
   1. Operator → `b := pop()` (right), `a := pop()` (left), `push(apply(a, b, op))`.
   2. Number → parse and push.
3. Return the single value remaining on the stack.

### Complexity
- **Time:** O(n) — each token is processed once with O(1) stack work.
- **Space:** O(n) — worst-case stack height when all numbers precede the operators.

### Code
```go
func stackApproach(tokens []string) int {
	stack := make([]int, 0, len(tokens))
	for _, token := range tokens {
		if isOperator(token) {
			b := stack[len(stack)-1] // top = RIGHT operand
			a := stack[len(stack)-2] // beneath = LEFT operand
			stack = stack[:len(stack)-2]
			stack = append(stack, apply(a, b, token))
		} else {
			value, _ := strconv.Atoi(token)
			stack = append(stack, value)
		}
	}
	return stack[0]
}

func apply(a, b int, op string) int {
	switch op {
	case "+":
		return a + b
	case "-":
		return a - b
	case "*":
		return a * b
	default: // "/"
		return a / b // Go truncates toward zero, e.g. 6 / -132 = 0
	}
}
```

### Dry Run
Example 1: `tokens = ["2","1","+","3","*"]`.

| Step | Token | Action                             | Stack (bottom → top) |
|------|-------|------------------------------------|-----------------------|
| 1    | `2`   | push 2                             | `[2]`                 |
| 2    | `1`   | push 1                             | `[2, 1]`              |
| 3    | `+`   | pop b=1, a=2 → push 2+1 = 3        | `[3]`                 |
| 4    | `3`   | push 3                             | `[3, 3]`              |
| 5    | `*`   | pop b=3, a=3 → push 3*3 = 9        | `[9]`                 |

Return `stack[0]` → Output: `9` ✓

(Example 3's key moment, for the division trap: stack `[10, 6, -132]` hits `/` → `6 / -132` truncates toward zero to `0`, not floor to `-1`.)

---

## Key Takeaways

- **Postfix ⇒ stack.** Whenever operands appear before their operator (RPN, post-order anything), a stack evaluates it in one pass. This generalizes to LeetCode #224/#227 (infix calculators = shunting-yard + this evaluator).
- **Pop order matters:** top of stack is the **right** operand. `a - b` with swapped pops is the most common bug in this problem.
- **Know your division:** the problem demands truncation toward zero. Go and C/Java/Rust truncate natively; Python's `//` floors (`6 // -132 == -1` — wrong here). State this in interviews.
- **Operator vs negative number:** match tokens against the four operator strings — `"-11"` fails that match and falls through to number parsing. Checking "starts with `-`" is a bug.
- RPN is a **post-order serialization of an expression tree** — recursion from the right end rebuilds/evaluates the tree without any parentheses or precedence rules. That's the whole point of the notation: precedence is encoded in position.

---

## Related Problems

- LeetCode #224 — Basic Calculator (infix with parentheses; stack of signs/results)
- LeetCode #227 — Basic Calculator II (infix with precedence; converts naturally via stacks)
- LeetCode #772 — Basic Calculator III (infix, parentheses + precedence combined)
- LeetCode #1006 — Clumsy Factorial (stack evaluation with fixed operator cycle)
- LeetCode #682 — Baseball Game (push/pop token processing on a stack)
