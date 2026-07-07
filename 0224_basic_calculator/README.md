# 0224 — Basic Calculator

> LeetCode #224 · Difficulty: Hard
> **Categories:** Math, String, Stack, Recursion

---

## Problem Statement

Given a string `s` representing a valid expression, implement a basic calculator to evaluate it, and return *the result of the evaluation*.

**Note:** You are **not** allowed to use any built-in function which evaluates strings as mathematical expressions, such as `eval()`.

**Example 1:**

```
Input: s = "1 + 1"
Output: 2
```

**Example 2:**

```
Input: s = " 2-1 + 2 "
Output: 3
```

**Example 3:**

```
Input: s = "(1+(4+5+2)-3)+(6+8)"
Output: 23
```

**Constraints:**

- `1 <= s.length <= 3 * 10^5`
- `s` consists of digits, `'+'`, `'-'`, `'('`, `')'`, and `' '`.
- `s` represents a valid expression.
- `'+'` is **not** used as a unary operation (i.e., `"+1"` and `"+(2 + 3)"` is invalid).
- `'-'` **could** be used as a unary operation (i.e., `"-1"` and `"-(2 + 3)"` is valid).
- There will be no two consecutive operators in the input.
- Every number and running calculation will fit in a signed 32-bit integer.

---

## Company Frequency

| Company    | Frequency        | Last Reported |
|------------|------------------|---------------|
| Google     | ★★★★★ Very High  | 2024          |
| Amazon     | ★★★★☆ High       | 2024          |
| Microsoft  | ★★★★☆ High       | 2024          |
| Meta       | ★★★☆☆ Medium     | 2023          |
| Bloomberg  | ★★★☆☆ Medium     | 2023          |
| Apple      | ★★☆☆☆ Low        | 2023          |

> ⚠️ Frequency data is crowd-sourced from LeetCode Discuss, Glassdoor, and
> community interview reports. Treat as a signal, not a guarantee.

---

## DSA Concepts Used

- **Stack** — an explicit stack saves the enclosing `(result, sign)` context when a `(` opens a sub-expression and restores it on `)` → see [`/dsa/stack.md`](/dsa/stack.md)
- **String Parsing** — single left-to-right scan that builds multi-digit numbers and reacts to operators/parentheses → see [`/dsa/string_algorithms.md`](/dsa/string_algorithms.md)

---

## Approaches Overview

| # | Approach | Time | Space | When to use |
|---|----------|------|-------|-------------|
| 1 | Stack of (result, sign) Contexts | O(n) | O(n) | The canonical iterative solution; safe against deep nesting |
| 2 | Recursive Descent | O(n) | O(d) | Cleaner to read; `d` = nesting depth (uses the call stack instead) |

---

## Approach 1 — Stack of (result, sign) Contexts

### Intuition
There is no `*` or `/`, so the only precedence comes from parentheses. Scan left to right maintaining a running `result` and the `sign` for the next number, folding numbers in with `result += sign * number`. A `(` starts a fresh sub-expression whose value must later be combined with the sign that preceded it — so **push** the current `(result, sign)` and reset. A `)` finishes the group: pop the saved context and combine `savedResult + savedSign * innerResult`. Unary minus works for free: `-(2+3)` sees `sign = −1` pushed before the group.

### Algorithm
1. Init `result=0`, `number=0`, `sign=+1`, empty `stack`.
2. For each character:
   - digit → `number = number*10 + digit`.
   - `+` → `result += sign*number`; reset `number`; `sign = +1`.
   - `-` → `result += sign*number`; reset `number`; `sign = −1`.
   - `(` → commit pending number; push `result` then `sign`; reset `result=0, sign=+1`.
   - `)` → commit number; pop `savedSign`, `savedResult`; `result = savedResult + savedSign*result`.
   - space → skip.
3. Commit any trailing number; return `result`.

### Complexity
- **Time:** O(n) — one linear scan.
- **Space:** O(n) — stack grows with parenthesis nesting (worst case O(n)).

### Code
```go
func stackCalculator(s string) int {
	result := 0       // running total of the current (sub)expression
	number := 0       // integer currently being parsed
	sign := 1         // sign applied to the next number (+1 or −1)
	stack := []int{}  // holds (result, sign) pairs of enclosing contexts
	hasDigit := false // whether `number` currently holds parsed digits

	for i := 0; i < len(s); i++ {
		c := s[i]
		switch {
		case c >= '0' && c <= '9':
			number = number*10 + int(c-'0') // extend the multi-digit number
			hasDigit = true
		case c == '+':
			result += sign * number // commit the number just parsed
			number, hasDigit = 0, false
			sign = 1 // next number is added
		case c == '-':
			result += sign * number
			number, hasDigit = 0, false
			sign = -1 // next number is subtracted
		case c == '(':
			// entering a group: save context, then start fresh inside
			stack = append(stack, result, sign)
			result, sign = 0, 1
			number, hasDigit = 0, false
		case c == ')':
			result += sign * number // commit the last number inside the group
			number, hasDigit = 0, false
			savedSign := stack[len(stack)-1]
			savedResult := stack[len(stack)-2]
			stack = stack[:len(stack)-2]
			result = savedResult + savedSign*result // fold inner value back
		}
		// spaces fall through and are ignored
	}
	if hasDigit {
		result += sign * number // commit any trailing number
	}
	return result
}
```

### Dry Run
Example 3: `s = "(1+(4+5+2)-3)+(6+8)"`. (`R`=result, `sgn`=sign, `num`=number, `st`=stack)

| char | action | R | sgn | num | st |
|------|--------|---|-----|-----|----|
| `(` | push (0,+1), reset | 0 | +1 | 0 | [0,1] |
| `1` | build num | 0 | +1 | 1 | [0,1] |
| `+` | R+=+1·1 | 1 | +1 | 0 | [0,1] |
| `(` | push (1,+1), reset | 0 | +1 | 0 | [0,1,1,1] |
| `4` `+` | R+=4 | 4 | +1 | 0 | [0,1,1,1] |
| `5` `+` | R+=5 | 9 | +1 | 0 | [0,1,1,1] |
| `2` | build | 9 | +1 | 2 | [0,1,1,1] |
| `)` | R+=2=11; pop sgn=+1,R'=1 → 1+1·11 | 11 | +1 | 0 | [0,1] |
| `-` | (num 0) sign=−1 | 11 | −1 | 0 | [0,1] |
| `3` | build | 11 | −1 | 3 | [0,1] |
| `)` | R+=−1·3=8; pop sgn=+1,R'=0 → 0+1·8 | 8 | +1 | 0 | [] |
| `+` | (num 0) sign=+1 | 8 | +1 | 0 | [] |
| `(` | push (8,+1), reset | 0 | +1 | 0 | [8,1] |
| `6` `+` | R+=6 | 6 | +1 | 0 | [8,1] |
| `8` | build | 6 | +1 | 8 | [8,1] |
| `)` | R+=8=14; pop sgn=+1,R'=8 → 8+1·14 | 22 | +1 | 0 | [] |
| end | commit num 0 | **22** | | | |

Wait — the trace lands on 22? Re-checking the first `)`: after the inner group `(4+5+2)` we get 11, folded with saved `(1,+1)` → `1 + 1·11 = 12`, not 11. Correcting: at that pop `savedResult=1`, so `R = 1 + 1·11 = 12`. Then `-3` → `R += −1·3 = 12−3 = 9`; the outer `)` pops saved `(0,+1)` → `0 + 1·9 = 9`. Then `+(6+8)=14` → `9 + 14 = 23`. Final answer **23**. ✔ (The table above under-counted the saved `1`; the corrected fold gives 23.)

---

## Approach 2 — Recursive Descent

### Intuition
Every `(` opens a self-contained sub-expression. Instead of an explicit stack, let recursion hold the enclosing context: when the scanner meets `(`, it recursively evaluates the inside, receives the group's value, and treats that value as if it were a plain number. Each `)` returns the level's value and the position just past it, so the caller resumes cleanly.

### Algorithm
1. `eval(s, i)` keeps local `result`, `number`, `sign`, starting at index `i`.
2. digit → build `number`; `+`/`-` → commit and set `sign`; space → skip.
3. `(` → `(inner, i) = eval(s, i+1)`; set `number = inner`.
4. `)` → commit `number`, return `(result, i+1)`.
5. At end of string, commit trailing `number` and return.

### Complexity
- **Time:** O(n) — each character consumed exactly once across all frames.
- **Space:** O(d) — recursion depth equals maximum parenthesis nesting.

### Code
```go
func recursiveCalculator(s string) int {
	val, _ := eval(s, 0)
	return val
}

func eval(s string, i int) (int, int) {
	result := 0 // running total for this parenthesis level
	number := 0 // integer being parsed
	sign := 1   // sign for the next number
	for i < len(s) {
		c := s[i]
		switch {
		case c >= '0' && c <= '9':
			number = number*10 + int(c-'0')
			i++
		case c == '+':
			result += sign * number
			number, sign = 0, 1
			i++
		case c == '-':
			result += sign * number
			number, sign = 0, -1
			i++
		case c == '(':
			var inner int
			inner, i = eval(s, i+1) // evaluate nested group; value acts as a number
			number = inner
		case c == ')':
			result += sign * number // commit last number of this level
			return result, i + 1    // hand back position after ')'
		default: // space
			i++
		}
	}
	result += sign * number // commit trailing number at top level
	return result, i
}
```

### Dry Run
Example 3: `"(1+(4+5+2)-3)+(6+8)"`.

| Frame | Sub-expression | Steps | Returns value |
|-------|----------------|-------|---------------|
| eval #3 | `4+5+2` | 4, +5, +2 → 11 | 11 |
| eval #2 | `1+(…)-3` | 1, +[11]=12, −3 | 9 |
| eval #4 | `6+8` | 6, +8 | 14 |
| eval #1 (top) | `(…)+(…)` | number=9 from #2, +14 from #4 | **23** |

Answer `23`. ✔

---

## Key Takeaways
- With only `+`/`-` and parentheses, the entire evaluation is `result += sign * number` — no operator precedence stack needed, just a **sign** and a way to save/restore context at parentheses.
- Pushing `(result, sign)` on `(` and folding back `savedResult + savedSign*inner` on `)` is the reusable "evaluate with parentheses" pattern.
- Unary minus needs no special case: the `sign` carried into a `(` is exactly the sign applied to the whole group.
- Recursion and an explicit stack are duals here; recursion uses O(depth) call-stack space, the explicit stack makes that space visible and avoids deep-recursion limits.

---

## Related Problems
- LeetCode #227 — Basic Calculator II (adds `*` and `/`, no parentheses — needs precedence)
- LeetCode #772 — Basic Calculator III (`+ - * /` **and** parentheses — combines both)
- LeetCode #150 — Evaluate Reverse Polish Notation (postfix, pure stack)
- LeetCode #394 — Decode String (same push/pop-context-on-brackets pattern)
