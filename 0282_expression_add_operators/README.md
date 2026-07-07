# 0282 — Expression Add Operators

> LeetCode #282 · Difficulty: Hard
> **Categories:** Backtracking, Math, String, Recursion

---

## Problem Statement

Given a string `num` that contains only digits and an integer `target`, return **all possibilities** to insert the binary operators `'+'`, `'-'`, and/or `'*'` between the digits of `num` so that the resultant expression evaluates to the `target` value.

Note that operands in the returned expressions **should not** contain leading zeros.

**Example 1:**
```
Input: num = "123", target = 6
Output: ["1*2*3","1+2+3"]
Explanation: Both "1*2*3" and "1+2+3" evaluate to 6.
```

**Example 2:**
```
Input: num = "232", target = 8
Output: ["2*3+2","2+3*2"]
Explanation: Both "2*3+2" and "2+3*2" evaluate to 8.
```

**Example 3:**
```
Input: num = "3456237490", target = 9191
Output: []
Explanation: There are no expressions that can be created from "3456237490" to evaluate to 9191.
```

**Constraints:**
- `1 <= num.length <= 10`
- `num` consists of only digits.
- `-2³¹ <= target <= 2³¹ - 1`

---

## Company Frequency

| Company    | Frequency        | Last Reported |
|------------|------------------|---------------|
| Google     | ★★★★★ Very High  | 2024          |
| Meta       | ★★★★☆ High       | 2024          |
| Amazon     | ★★★☆☆ Medium     | 2023          |
| Microsoft  | ★★★☆☆ Medium     | 2023          |
| Apple      | ★★☆☆☆ Low        | 2023          |
| Bloomberg  | ★★☆☆☆ Low        | 2022          |

> ⚠️ Frequency data is crowd-sourced from LeetCode Discuss, Glassdoor, and
> community interview reports. Treat as a signal, not a guarantee.

---

## DSA Concepts Used
- **Backtracking** — DFS over operator choices at every gap, undoing state as we return → see [`/dsa/backtracking.md`](/dsa/backtracking.md)
- **String partitioning / building** — split digits into multi-digit operands and assemble candidate expressions → see [`/dsa/string_algorithms.md`](/dsa/string_algorithms.md)
- **Operator-precedence evaluation** — track a running term so `*` binds tighter than `+`/`-` → see [`/dsa/math_number_theory.md`](/dsa/math_number_theory.md)

---

## Approaches Overview
| # | Approach | Time | Space | When to use |
|---|----------|------|-------|-------------|
| 1 | Brute Force (enumerate assignments) | O(4ⁿ·n) | O(n) | Baseline; builds then re-evaluates each string |
| 2 | Backtracking with carried value (Optimal) | O(4ⁿ·n) | O(n) | Standard; evaluates incrementally, prunes early |

---

## Approach 1 — Brute Force

### Intuition
Between `n` digits there are `n-1` gaps. Each gap gets one of four choices: `+`, `-`, `*`, or "nothing" (which glues digits into a longer number). That's `4^(n-1)` candidate strings. Build every one, reject those with leading-zero operands, evaluate honouring `*` precedence, and keep the ones equal to `target`.

### Algorithm
1. Enumerate every base-4 mask over the `n-1` gaps.
2. For each mask, insert the chosen operator (or nothing) between digits, tracking numeric tokens.
3. Reject expressions whose any operand has a leading zero (e.g. `"05"`).
4. Evaluate with precedence (`*` before `+`/`-`); if it equals `target`, record the expression.

### Complexity
- **Time:** O(4ⁿ · n) — `4^(n-1)` expressions, each parsed and evaluated in O(n).
- **Space:** O(n) for token/expression building, plus the result list.

### Code
```go
func bruteForce(num string, target int) []string {
	n := len(num)
	res := []string{}
	if n == 0 {
		return res
	}
	ops := []string{"", "+", "-", "*"}
	gaps := n - 1

	total := 1
	for i := 0; i < gaps; i++ {
		total *= 4
	}

	for mask := 0; mask < total; mask++ {
		expr := string(num[0])
		tokens := []string{string(num[0])}
		valid := true
		m := mask
		for g := 0; g < gaps; g++ {
			op := ops[m%4]
			m /= 4
			expr += op + string(num[g+1])
			if op == "" {
				tokens[len(tokens)-1] += string(num[g+1])
			} else {
				tokens = append(tokens, op, string(num[g+1]))
			}
		}
		for i := 0; i < len(tokens); i += 2 {
			t := tokens[i]
			if len(t) > 1 && t[0] == '0' {
				valid = false
				break
			}
		}
		if valid && evaluate(tokens) == target {
			res = append(res, expr)
		}
	}
	sort.Strings(res)
	return res
}

func evaluate(tokens []string) int {
	result := 0
	term, _ := strconv.Atoi(tokens[0])
	for i := 1; i < len(tokens); i += 2 {
		op := tokens[i]
		v, _ := strconv.Atoi(tokens[i+1])
		switch op {
		case "+":
			result += term
			term = v
		case "-":
			result += term
			term = -v
		case "*":
			term *= v
		}
	}
	return result + term
}
```

### Dry Run
`num = "123"`, `target = 6`. Gaps = 2, so 16 masks. Showing the two that hit target:

| mask (gap1,gap2 ops) | expression | tokens | evaluate | == 6? |
|----------------------|------------|--------|----------|-------|
| `+`,`+` | `1+2+3` | [1,+,2,+,3] | 1→term1; +2→result1,term2; +3→result3,term3 ⇒ 6 | ✅ |
| `*`,`*` | `1*2*3` | [1,*,2,*,3] | term1; *2→term2; *3→term6 ⇒ 0+6 | ✅ |
| `+`,`*` | `1+2*3` | [1,+,2,*,3] | 1→result1,term2; *3→term6 ⇒ 1+6=7 | ❌ |

Sorted result: `["1*2*3","1+2+3"]`.

---

## Approach 2 — Backtracking (Optimal)

### Intuition
The tricky part is `*` precedence. Instead of re-parsing strings, carry two values down the recursion: `total` (value so far) and `prev` (the signed value of the **last operand**). Then:
- `+v` → `total += v`, `prev = v`.
- `-v` → `total -= v`, `prev = -v`.
- `*v` → we must undo the last addition and multiply into it: `total = total - prev + prev*v`, `prev = prev*v`.

Each leaf is evaluated in O(1), and leading-zero operands are pruned the moment they'd form.

### Algorithm
1. `dfs(pos, expr, total, prev)`:
2. If `pos == n`: if `total == target`, record `expr`; return.
3. For `end = pos … n-1`: take operand `num[pos..end]`; if it has a leading zero (`num[pos]=='0'` and `end>pos`) break.
4. If `pos == 0`: seed the first operand with no operator → `dfs(end+1, str, cur, cur)`.
5. Else branch three ways: `+`, `-`, `*` with the updates above.

### Complexity
- **Time:** O(4ⁿ · n) worst case — branching up to 4 per gap, O(n) to append the string at each step.
- **Space:** O(n) recursion depth (plus output storage).

### Code
```go
func backtracking(num string, target int) []string {
	res := []string{}
	n := len(num)
	if n == 0 {
		return res
	}

	var dfs func(pos int, expr string, total, prev int)
	dfs = func(pos int, expr string, total, prev int) {
		if pos == n {
			if total == target {
				res = append(res, expr)
			}
			return
		}
		for end := pos; end < n; end++ {
			if end > pos && num[pos] == '0' {
				break
			}
			cur, _ := strconv.Atoi(num[pos : end+1])
			str := num[pos : end+1]
			if pos == 0 {
				dfs(end+1, str, cur, cur)
			} else {
				dfs(end+1, expr+"+"+str, total+cur, cur)
				dfs(end+1, expr+"-"+str, total-cur, -cur)
				dfs(end+1, expr+"*"+str, total-prev+prev*cur, prev*cur)
			}
		}
	}
	dfs(0, "", 0, 0)
	sort.Strings(res)
	return res
}
```

### Dry Run
`num = "123"`, `target = 6`. Trace the two accepting paths (operands taken one digit at a time):

| step | expr | total | prev | note |
|------|------|-------|------|------|
| seed | `1` | 1 | 1 | first operand |
| `+2` | `1+2` | 3 | 2 | total 1+2 |
| `+3` | `1+2+3` | 6 | 3 | leaf: total==6 ✅ record |
| `*2` (from `1`) | `1*2` | `1-1+1*2 = 2` | 2 | undo prev 1, apply ×2 |
| `*3` | `1*2*3` | `2-2+2*3 = 6` | 6 | leaf: total==6 ✅ record |

Other branches (`1+2*3=7`, `1-…`, etc.) miss. Sorted result: `["1*2*3","1+2+3"]`.

---

## Key Takeaways
- **Carry the running value, not the string.** Tracking `(total, prev)` turns multiplication precedence into an O(1) arithmetic fix `total - prev + prev*cur` — the signature trick of this problem.
- **Prune leading zeros at the source:** once `num[pos]=='0'` you may only take the single digit `"0"` as an operand, so `break` after `end == pos`.
- Enumerating operators is a **4-way branch per gap** (`+`, `-`, `*`, concatenate), giving the `4^(n-1)` bound.

---

## Related Problems
- LeetCode #494 — Target Sum (assign +/- to reach a target)
- LeetCode #227 — Basic Calculator II (evaluate expression with `*`/`/` precedence)
- LeetCode #241 — Different Ways to Add Parentheses (all parenthesizations)
- LeetCode #679 — 24 Game (insert operators/parentheses to hit 24)
