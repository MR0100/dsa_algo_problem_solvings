# 0254 — Factor Combinations

> LeetCode #254 · Difficulty: Medium
> **Categories:** Backtracking, Recursion, Math

---

## Problem Statement

Numbers can be regarded as the product of their factors.

- For example, `8 = 2 x 2 x 2 = 2 x 4`.

Given an integer `n`, return all possible combinations of its factors. You may return the answer in any order.

**Note** that the factors should be in the range `[2, n - 1]`.

**Example 1:**

```
Input: n = 1
Output: []
```

**Example 2:**

```
Input: n = 12
Output: [[2,6],[2,2,3],[3,4]]
```

**Example 3:**

```
Input: n = 37
Output: []
```

**Constraints:**

- `1 <= n <= 10^7`

---

## Company Frequency

| Company   | Frequency        | Last Reported |
|-----------|------------------|---------------|
| Google    | ★★★★☆ High       | 2023          |
| Amazon    | ★★★☆☆ Medium     | 2022          |
| Uber      | ★★☆☆☆ Low        | 2022          |
| LinkedIn  | ★★☆☆☆ Low        | 2021          |

> ⚠️ Frequency data is crowd-sourced from LeetCode Discuss, Glassdoor, and community interview reports. Treat as a signal, not a guarantee.

---

## DSA Concepts Used

- **Backtracking** — explore factorizations by choosing factors in order, recursing on the quotient → see [`/dsa/backtracking.md`](/dsa/backtracking.md)
- **Divide & conquer** — the answer for `n` is built from factorizations of its quotients → see [`/dsa/divide_and_conquer.md`](/dsa/divide_and_conquer.md)
- **Number theory (divisors)** — we only iterate candidate factors up to √n → see [`/dsa/math_number_theory.md`](/dsa/math_number_theory.md)

---

## Approaches Overview

| # | Approach | Time | Space | When to use |
|---|----------|------|-------|-------------|
| 1 | Backtracking over Divisors | O(#factorizations · depth) | O(log n) recursion | Standard in-place path building |
| 2 | DFS Collect (Alt Form) | O(#factorizations · depth) | O(log n) + output | When you prefer returning sub-results explicitly |

---

## Approach 1 — Backtracking over Divisors

### Intuition
Every non-trivial factorization of `n` is a multiset of factors ≥ 2 whose product is `n`. To avoid duplicates like `[2,6]` vs `[6,2]`, force factors into non-decreasing order by only trying candidates ≥ the last one chosen. For a remaining value, each divisor `f` (with `f*f <= remaining`) yields a complete factorization `path + [f, remaining/f]`, and we recurse on the quotient to break it down further.

### Algorithm
1. `dfs(remaining, start, path)`: iterate `f = start; f*f <= remaining; f++`.
2. If `remaining % f == 0`: record `path + [f, remaining/f]` as one answer, then recurse `dfs(remaining/f, f, path+[f])`.
3. Start with `dfs(n, 2, [])`. Factors begin at 2; capping the loop at √remaining avoids the trivial `[n]` and duplicate splits.

### Complexity
- **Time:** O(number of factorizations · depth) — bounded by the divisor tree of `n`; each valid path is emitted once.
- **Space:** O(log n) recursion depth (each division shrinks the value), plus the output list.

### Code
```go
func backtracking(n int) [][]int {
	result := [][]int{}
	var dfs func(remaining, start int, path []int)
	dfs = func(remaining, start int, path []int) {
		for f := start; f*f <= remaining; f++ {
			if remaining%f == 0 {
				comb := make([]int, len(path))
				copy(comb, path)
				comb = append(comb, f, remaining/f)
				result = append(result, comb)

				dfs(remaining/f, f, append(path, f))
			}
		}
	}
	dfs(n, 2, []int{})
	return result
}
```

### Dry Run
Input `n = 12`. Start `dfs(12, 2, [])`.

| call                 | f | f*f ≤ rem? | 12%f | emit             | recurse            |
|----------------------|---|-----------|------|------------------|--------------------|
| dfs(12,2,[])         | 2 | 4≤12 yes  | 0    | [2,6]            | dfs(6,2,[2])       |
| dfs(6,2,[2])         | 2 | 4≤6 yes   | 0    | [2,2,3]          | dfs(3,2,[2,2])     |
| dfs(3,2,[2,2])       | 2 | 4≤3 no    | —    | (loop ends)      | —                  |
| dfs(12,2,[]) cont.   | 3 | 9≤12 yes  | 0    | [3,4]            | dfs(4,3,[3])       |
| dfs(4,3,[3])         | 3 | 9≤4 no    | —    | (loop ends)      | —                  |
| dfs(12,2,[]) cont.   | 4 | 16≤12 no  | —    | loop ends        | —                  |

Result: `[[2,6],[2,2,3],[3,4]]`.

---

## Approach 2 — DFS Collect (Alt Form)

### Intuition
Same recursion, framed as "combine the current factor with sub-factorizations of the quotient". `factorize(value, start)` returns all factorizations of `value` using factors ≥ `start`: for each divisor `f`, emit the simple split `[f, value/f]`, and prepend `f` to every deeper factorization of `value/f`.

### Algorithm
1. `helper(value, start)` returns `[][]int`.
2. For `f = start; f*f <= value; f++` with `value % f == 0`: add `[f, value/f]`, then for each `sub` in `helper(value/f, f)` add `[f] ++ sub`.
3. Return `helper(n, 2)`.

### Complexity
- **Time:** O(number of factorizations · depth), same class as Approach 1.
- **Space:** O(log n) recursion depth plus the produced result lists.

### Code
```go
func dfsCollect(n int) [][]int {
	var helper func(value, start int) [][]int
	helper = func(value, start int) [][]int {
		res := [][]int{}
		for f := start; f*f <= value; f++ {
			if value%f == 0 {
				co := value / f
				res = append(res, []int{f, co})
				for _, sub := range helper(co, f) {
					combined := append([]int{f}, sub...)
					res = append(res, combined)
				}
			}
		}
		return res
	}
	return helper(n, 2)
}
```

### Dry Run
Input `n = 12`. Call `helper(12, 2)`.

| f | 12%f | co  | direct add | helper(co,f)        | prepended     |
|---|------|-----|-----------|---------------------|---------------|
| 2 | 0    | 6   | [2,6]     | helper(6,2)=[[2,3]] | [2,2,3]       |
| 3 | 0    | 4   | [3,4]     | helper(4,3)=[]      | (none)        |
| 4 | 16>12 stop |  |          |                     |               |

`helper(6,2)`: f=2 → [2,3], helper(3,2)=[] → returns `[[2,3]]`.

Result order: `[[2,6],[2,2,3],[3,4]]` — matches Approach 1.

---

## Key Takeaways
- Enforcing non-decreasing factor order (candidates ≥ last chosen) is the standard trick to emit each multiset exactly once — no duplicate `[2,6]`/`[6,2]`.
- Looping factors only up to `√remaining` both bounds the work and naturally excludes the trivial `[n]` factorization.
- Copy the path when recording a result (`make` + `copy`); reusing the same backing slice across branches causes classic aliasing bugs in Go backtracking.

---

## Related Problems
- LeetCode #39 — Combination Sum (choose with reuse, non-decreasing order)
- LeetCode #40 — Combination Sum II (combinations without duplicates)
- LeetCode #46 — Permutations (backtracking template)
- LeetCode #78 — Subsets (enumerate all subsets)
