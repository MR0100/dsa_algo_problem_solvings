# 0241 — Different Ways to Add Parentheses

> LeetCode #241 · Difficulty: Medium
> **Categories:** Divide and Conquer, Recursion, Memoization, Math, String

---

## Problem Statement

Given a string `expression` of numbers and operators, return *all possible results from computing all the different possible ways to group numbers and operators*. You may return the answer in **any order**.

The test cases are generated such that the output values fit in a 32-bit integer and the number of different results does not exceed `10^4`.

**Example 1:**

```
Input: expression = "2-1-1"
Output: [0,2]
Explanation:
((2-1)-1) = 0
(2-(1-1)) = 2
```

**Example 2:**

```
Input: expression = "2*3-4*5"
Output: [-34,-14,-10,-10,10]
Explanation:
(2*(3-(4*5))) = -34
((2*3)-(4*5)) = -14
((2*(3-4))*5) = -10
(2*((3-4)*5)) = -10
(((2*3)-4)*5) = 10
```

**Constraints:**

- `1 <= expression.length <= 20`
- `expression` consists of digits and the operator `'+'`, `'-'`, and `'*'`.
- All the integer values in the input expression are in the range `[0, 99]`.

---

## Company Frequency

| Company    | Frequency        | Last Reported |
|------------|------------------|---------------|
| Amazon     | ★★★★☆ High       | 2024          |
| Google     | ★★★★☆ High       | 2024          |
| Microsoft  | ★★★☆☆ Medium     | 2023          |
| Meta       | ★★★☆☆ Medium     | 2023          |
| Bloomberg  | ★★☆☆☆ Low        | 2023          |

> ⚠️ Frequency data is crowd-sourced from LeetCode Discuss, Glassdoor, and
> community interview reports. Treat as a signal, not a guarantee.

---

## DSA Concepts Used

- **Divide and Conquer** — each operator is a candidate "last operation"; fixing it splits the expression into an independent left and right sub-problem whose results combine → see [`/dsa/divide_and_conquer.md`](/dsa/divide_and_conquer.md)
- **Memoization / Dynamic Programming** — identical substrings recur across many splits; caching each substring's result set removes exponential re-computation → see [`/dsa/dynamic_programming_1d.md`](/dsa/dynamic_programming_1d.md)
- **String Parsing** — walking the expression to locate operator positions and slice operands → see [`/dsa/string_algorithms.md`](/dsa/string_algorithms.md)

---

## Approaches Overview

| # | Approach | Time | Space | When to use |
|---|----------|------|-------|-------------|
| 1 | Divide and Conquer | O(Catalan(n)) | O(Catalan(n)) | Clean baseline; fine for the tiny input bound |
| 2 | Divide and Conquer + Memoization (Optimal) | O(Catalan(n)) | O(distinct substrings × results) | Avoids re-solving shared sub-expressions |

(n = number of operators.)

---

## Approach 1 — Divide and Conquer

### Intuition
Every valid parenthesization designates exactly one operator as the *last* one evaluated. If we fix that operator, everything to its left is one fully-parenthesizable sub-expression and everything to its right is another. Recursively enumerate every result of each side, then apply the fixed operator to every left/right pair.

### Algorithm
1. If the string has no operator, it is a single integer — return `[that number]`.
2. For each index `i` where `expression[i]` is `+`, `-`, or `*`:
   1. Recurse on `expression[:i]` to get the list of left values.
   2. Recurse on `expression[i+1:]` to get the list of right values.
   3. For every pair `(l, r)`, append `l op r` to the results.
3. Return all accumulated results.

### Complexity
- **Time:** O(Catalan(n)) — the count of distinct parenthesizations of `n` operators is the nth Catalan number; each is produced exactly once.
- **Space:** O(Catalan(n)) for the output list, plus O(n) recursion depth.

### Code
```go
func divideAndConquer(expression string) []int {
	if isNumber(expression) {
		n, _ := strconv.Atoi(expression)
		return []int{n}
	}

	var results []int
	for i := 0; i < len(expression); i++ {
		c := expression[i]
		if c == '+' || c == '-' || c == '*' {
			left := divideAndConquer(expression[:i])
			right := divideAndConquer(expression[i+1:])
			for _, l := range left {
				for _, r := range right {
					results = append(results, apply(l, r, c))
				}
			}
		}
	}
	return results
}
```

### Dry Run
Trace `"2-1-1"`:

| Call | Split operator | Left results | Right results | Combined |
|------|----------------|--------------|---------------|----------|
| `"2-1-1"` | `-` at index 1 | `solve("2")=[2]` | `solve("1-1")` | see below |
| `"1-1"` | `-` at index 1 | `[1]` | `[1]` | `1-1 = [0]` |
| back in `"2-1-1"` split@1 | | `[2]` | `[0]` | `2-0 = 2` |
| `"2-1-1"` | `-` at index 3 | `solve("2-1")` | `solve("1")=[1]` | see below |
| `"2-1"` | `-` at index 1 | `[2]` | `[1]` | `2-1 = [1]` |
| back in `"2-1-1"` split@3 | | `[1]` | `[1]` | `1-1 = 0` |

Collected results = `{2, 0}` → sorted `[0, 2]`. ✓

---

## Approach 2 — Divide and Conquer + Memoization (Optimal)

### Intuition
Approach 1 re-solves the same substring many times — `"2*3"` inside `"2*3-4*5"` is recomputed for several outer splits. Since a substring's result set depends only on the substring, cache it in a map keyed by the substring text and reuse it.

### Algorithm
1. Keep a `memo` map from substring → its list of possible values.
2. On entry, if the substring is cached, return the cached slice.
3. Otherwise apply the same split-and-combine logic as Approach 1.
4. Store the computed slice in `memo` before returning.

### Complexity
- **Time:** O(Catalan(n)) results, but each distinct substring is expanded only once, eliminating the repeated exponential work of shared sub-expressions.
- **Space:** O(number of distinct substrings × results per substring) for the cache, plus recursion depth.

### Code
```go
func memoized(expression string) []int {
	memo := make(map[string][]int)
	var solve func(expr string) []int
	solve = func(expr string) []int {
		if v, ok := memo[expr]; ok {
			return v
		}
		if isNumber(expr) {
			n, _ := strconv.Atoi(expr)
			memo[expr] = []int{n}
			return memo[expr]
		}
		var results []int
		for i := 0; i < len(expr); i++ {
			c := expr[i]
			if c == '+' || c == '-' || c == '*' {
				left := solve(expr[:i])
				right := solve(expr[i+1:])
				for _, l := range left {
					for _, r := range right {
						results = append(results, apply(l, r, c))
					}
				}
			}
		}
		memo[expr] = results
		return results
	}
	return solve(expression)
}
```

### Dry Run
Trace `"2-1-1"` (showing cache hits):

| Step | Call | Cache state before | Action |
|------|------|--------------------|--------|
| 1 | `solve("2-1-1")` | empty | split@1 and split@3 |
| 2 | `solve("2")` | — | miss → store `"2":[2]` |
| 3 | `solve("1-1")` | `{"2"}` | miss → recurse |
| 4 | `solve("1")` | `{"2"}` | miss → store `"1":[1]` |
| 5 | `"1-1"` combines `[1],[1]` | — | store `"1-1":[0]` |
| 6 | back in split@1: `2-0=2` | | |
| 7 | `solve("2-1")` | `{"2","1","1-1"}` | miss → `"2":[2]` **hit**, `"1":[1]` **hit** → store `"2-1":[1]` |
| 8 | `solve("1")` (split@3 right) | | **hit** returns `[1]` |
| 9 | back in split@3: `1-1=0` | | |

Results `{2, 0}` → `[0, 2]`. Note steps 7–8 reuse cached `"2"` and `"1"`. ✓

---

## Key Takeaways
- **"Pick the last operation" is the divide-and-conquer key** for expression enumeration problems — fixing the top-level operator decomposes the string into independent halves.
- The number of parenthesizations of `n` binary operators is the **Catalan number** `C(n)`; that is the intrinsic size of the answer.
- **Memoize on the substring**, not on indices, because the value set is a pure function of the substring text — a clean cache key.
- Combining two result lists is a **Cartesian product** with the operator applied to each pair.

---

## Related Problems
- LeetCode #95 — Unique Binary Search Trees II (same split-at-root divide and conquer)
- LeetCode #96 — Unique Binary Search Trees (Catalan counting)
- LeetCode #22 — Generate Parentheses (enumerating parenthesizations)
- LeetCode #282 — Expression Add Operators (inserting operators into a digit string)
