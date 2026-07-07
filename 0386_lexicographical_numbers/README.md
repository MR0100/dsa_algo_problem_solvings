# 0386 — Lexicographical Numbers

> LeetCode #386 · Difficulty: Medium
> **Categories:** Depth-First Search, Trie, Math

---

## Problem Statement

Given an integer `n`, return all the numbers in the range `[1, n]` sorted in lexicographical order.

You must write an algorithm that runs in `O(n)` time and uses `O(1)` extra space (the returned list is not counted as extra space).

**Example 1:**

```
Input: n = 13
Output: [1,10,11,12,13,2,3,4,5,6,7,8,9]
```

**Example 2:**

```
Input: n = 2
Output: [1,2]
```

**Constraints:**

- `1 <= n <= 5 * 10^4`

---

## Company Frequency

| Company    | Frequency        | Last Reported |
|------------|------------------|---------------|
| ByteDance  | ★★★★☆ High       | 2024          |
| Amazon     | ★★★☆☆ Medium     | 2024          |
| Google     | ★★★☆☆ Medium     | 2023          |
| Microsoft  | ★★☆☆☆ Low        | 2023          |
| Apple      | ★★☆☆☆ Low        | 2022          |

> ⚠️ Frequency data is crowd-sourced from LeetCode Discuss, Glassdoor, and
> community interview reports. Treat as a signal, not a guarantee.

---

## DSA Concepts Used

- **Depth-First Search** — the numbers form a 10-ary "digit trie" (root children `1..9`, each node `x` has children `10x..10x+9`); a *preorder* DFS visits them in exactly lexicographic order → see [`/dsa/graph_bfs_dfs.md`](/dsa/graph_bfs_dfs.md)
- **Trie** — the mental model is a prefix tree over decimal strings; dictionary order = preorder over that tree → see [`/dsa/trie.md`](/dsa/trie.md)
- **Math / Number Theory** — the O(1)-space version computes the lexicographic *successor* of a number using `*10`, `/10`, and last-digit arithmetic → see [`/dsa/math_number_theory.md`](/dsa/math_number_theory.md)

---

## Approaches Overview

| # | Approach | Time | Space | When to use |
|---|----------|------|-------|-------------|
| 1 | Brute Force (String Sort) | O(n · L · log n) | O(n · L) | Baseline; ignores the O(n)/O(1) requirement |
| 2 | Preorder DFS | O(n) | O(log n) recursion | Cleanest to reason about; the trie insight |
| 3 | Iterative Successor (Optimal) | O(n) | O(1) | Meets the stated O(n) time / O(1) space bound |

---

## Approach 1 — Brute Force (String Sort)

### Intuition

"Lexicographical order" is just dictionary order of the decimal strings. The most direct thing that works: generate `1..n`, sort the values by their string form, done. It violates the required `O(n)` time / `O(1)` space but is a correct reference.

### Algorithm

1. Build `nums = [1, 2, ..., n]`.
2. Sort `nums` using a comparator that compares `strconv.Itoa(a)` with `strconv.Itoa(b)` (string comparison, not numeric).
3. Return `nums`.

### Complexity

- **Time:** O(n · L · log n) — `O(n log n)` comparisons, each comparing two strings of length up to `L = O(log n)`.
- **Space:** O(n · L) — the n values plus the transient strings created during comparison.

### Code

```go
func bruteForce(n int) []int {
	nums := make([]int, n) // will hold 1..n
	for i := 0; i < n; i++ {
		nums[i] = i + 1 // fill with the values 1..n
	}
	// Sort by the DECIMAL STRING of each value, i.e. dictionary order:
	// "10" < "2" because '1' < '2' at the first character.
	sort.Slice(nums, func(a, b int) bool {
		return strconv.Itoa(nums[a]) < strconv.Itoa(nums[b])
	})
	return nums
}
```

### Dry Run

Input `n = 13`. Initial `nums = [1,2,3,4,5,6,7,8,9,10,11,12,13]`.

| Step | Comparison (as strings) | Result |
|------|-------------------------|--------|
| — | `"1"` vs `"2"` | `"1" < "2"` → 1 before 2 |
| — | `"10"` vs `"2"` | `'1' < '2'` → `"10"` before `"2"` |
| — | `"10"` vs `"11"` | `'0' < '1'` → `"10"` before `"11"` |
| — | `"13"` vs `"2"` | `'1' < '2'` → `"13"` before `"2"` |
| final | full sort | `[1,10,11,12,13,2,3,4,5,6,7,8,9]` |

---

## Approach 2 — Preorder DFS

### Intuition

Numbers are paths in a tree keyed by digits. From `1` you descend to `10, 11, ..., 19` (append a digit) *before* moving to the sibling `2`. "Go deep appending `0..9`, then move to the next sibling" is exactly dictionary order — which is what a depth-first **preorder** traversal produces. Prune any branch whose value exceeds `n`.

### Algorithm

1. For each root digit `d` in `1..9`, call `visit(d)`.
2. `visit(cur)`:
   1. If `cur > n`, return (prune this whole subtree).
   2. Append `cur` (preorder: record before descending).
   3. For `next` in `0..9`, recurse `visit(cur*10 + next)`.

### Complexity

- **Time:** O(n) — each value in `[1..n]` is appended once; pruned calls are a constant factor per emitted node.
- **Space:** O(log n) — recursion depth = number of digits of `n`. The output slice is required output, not auxiliary.

### Code

```go
func dfs(n int) []int {
	result := make([]int, 0, n) // preallocate for the n emitted values

	// visit performs the preorder walk rooted at cur.
	var visit func(cur int)
	visit = func(cur int) {
		if cur > n { // cur (and all its 10x children) exceed n → prune
			return
		}
		result = append(result, cur) // preorder: record before descending
		// Children of cur are cur*10 + 0 .. cur*10 + 9, in ascending digit order.
		for next := 0; next <= 9; next++ {
			visit(cur*10 + next) // append the next digit and go deeper
		}
	}

	for d := 1; d <= 9; d++ { // the tree has 9 roots: 1..9 (no leading zero)
		visit(d)
	}
	return result
}
```

### Dry Run

Input `n = 13`. Trace the preorder walk (only relevant calls shown):

| Call | `cur` | Action |
|------|-------|--------|
| root | 1 | append 1 → `[1]` |
| child | 10 | append 10 → `[1,10]` |
| 10's children | 100.. | `100 > 13` → prune all |
| sibling | 11 | append → `[1,10,11]` |
| sibling | 12 | append → `[1,10,11,12]` |
| sibling | 13 | append → `[1,10,11,12,13]` |
| 14 | 14 | `14 > 13` → prune |
| root | 2 | append 2 → `[...,2]`; `20 > 13` prune |
| roots | 3..9 | append each; children all `>13` |

Final: `[1,10,11,12,13,2,3,4,5,6,7,8,9]`.

---

## Approach 3 — Iterative Successor (Optimal)

### Intuition

We can reproduce the DFS preorder with no stack by computing the lexicographic **successor** of the current number using only arithmetic:

1. **Go deeper** — `cur*10` (append a `0`) if it stays `≤ n`.
2. **Go right / climb** — otherwise, while the last digit is `9` (no right sibling) or `cur+1` would exceed `n`, climb up (`cur /= 10`), then `cur++`.

This is preorder traversal in `O(1)` memory, meeting the problem's requirement.

### Algorithm

1. `cur = 1`. Repeat `n` times:
   1. Append `cur`.
   2. If `cur*10 <= n`: `cur *= 10`.
   3. Else: while `cur%10 == 9 || cur+1 > n`: `cur /= 10`; then `cur++`.

### Complexity

- **Time:** O(n) — n iterations; the inner climb removes digits that were previously added, so it is amortised O(1) per step.
- **Space:** O(1) — only the scalar `cur` (output slice is required output).

### Code

```go
func iterativeNext(n int) []int {
	result := make([]int, 0, n) // n values will be produced
	cur := 1                    // lexicographically smallest positive number
	for i := 0; i < n; i++ {
		result = append(result, cur) // emit current number
		if cur*10 <= n {             // rule 1: can we append a '0' and stay ≤ n?
			cur *= 10 // go deeper: e.g. 1 -> 10
		} else {
			// rule 2: cannot go deeper; move right, climbing when blocked.
			// Climb while the last digit is 9 (no right sibling) OR
			// incrementing would exceed n (sibling out of range).
			for cur%10 == 9 || cur+1 > n {
				cur /= 10 // step up to the parent
			}
			cur++ // move to the next sibling at this (possibly higher) level
		}
	}
	return result
}
```

### Dry Run

Input `n = 13`:

| i | emit | `cur*10 ≤ 13`? | successor computation | new `cur` |
|---|------|----------------|-----------------------|-----------|
| 0 | 1 | `10 ≤ 13` yes | `cur *= 10` | 10 |
| 1 | 10 | `100 ≤ 13` no | last digit 0, `11 ≤ 13` → no climb; `cur++` | 11 |
| 2 | 11 | no | `12 ≤ 13` → `cur++` | 12 |
| 3 | 12 | no | `13 ≤ 13` → `cur++` | 13 |
| 4 | 13 | no | `13%10=3`, but `14 > 13` → climb `13/10=1`; then `cur++` | 2 |
| 5 | 2 | `20 ≤ 13` no | `3 ≤ 13` → `cur++` | 3 |
| 6–12 | 3..9 | no | `cur++` each; at 9 loop ends | 4..10 |

Final: `[1,10,11,12,13,2,3,4,5,6,7,8,9]`.

---

## Key Takeaways

- **Dictionary order = preorder over the digit trie.** Whenever a problem says "lexicographic order of numbers," think of the 10-ary tree (`x → 10x..10x+9`) and a preorder DFS.
- **Successor arithmetic replaces the stack.** `*10` = descend, `/10` = ascend, last-digit checks = sibling logic. This turns an O(depth)-space DFS into an O(1)-space iteration.
- **Pruning by value** (`cur > n`) keeps the traversal O(n): we never descend into a subtree whose smallest member already exceeds `n`.

---

## Related Problems

- LeetCode #440 — K-th Smallest in Lexicographical Order (count subtree sizes on the same digit trie)
- LeetCode #386 pattern also appears in Trie preorder enumeration problems
- LeetCode #1291 — Sequential Digits (enumerate numbers by digit structure)
