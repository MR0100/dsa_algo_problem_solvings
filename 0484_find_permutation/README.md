# 0484 — Find Permutation

> LeetCode #484 · Difficulty: Medium
> **Categories:** Array, String, Stack, Greedy

---

## Problem Statement

A permutation `perm` of `n` integers of all the integers in the range `[1, n]` can be represented as a string `s` of length `n - 1` where:

- `s[i] == 'I'` if `perm[i] < perm[i + 1]`, and
- `s[i] == 'D'` if `perm[i] > perm[i + 1]`.

Given a string `s`, reconstruct the lexicographically smallest permutation `perm` and return it.

**Example 1:**

```
Input: s = "I"
Output: [1,2]
Explanation: [1,2] is the only legal permutation that can represented by s, where the number 1 and 2 construct an increasing relationship.
```

**Example 2:**

```
Input: s = "DI"
Output: [2,1,3]
Explanation: Both [2,1,3] and [3,1,2] can be represented as "DI", but since we want to find the smallest lexicographical permutation, you should return [2,1,3].
```

**Constraints:**

- `1 <= s.length <= 10^5`
- `s[i]` is either `'I'` or `'D'`.

---

## Company Frequency

| Company    | Frequency        | Last Reported |
|------------|------------------|---------------|
| Google     | ★★★☆☆ Medium     | 2023          |
| Amazon     | ★★☆☆☆ Low        | 2022          |
| Facebook   | ★★☆☆☆ Low        | 2022          |
| Microsoft  | ★☆☆☆☆ Rare       | 2021          |

> ⚠️ Frequency data is crowd-sourced from LeetCode Discuss, Glassdoor, and
> community interview reports. Treat as a signal, not a guarantee.

---

## DSA Concepts Used

- **Greedy (minimal local disturbance)** — the identity `[1..n]` is already lex-smallest; each maximal `'D'` run is fixed by reversing exactly its ascending window, the least change that satisfies the constraint. Choosing that local minimum everywhere yields the global minimum → see [`/dsa/greedy.md`](/dsa/greedy.md)
- **Stack (LIFO reversal)** — pushing `1..n` and flushing on each `'I'` uses a stack's natural order-reversal to realise the descending blocks in a single pass → see [`/dsa/stack.md`](/dsa/stack.md)
- **Array in-place reversal** — the D-run fix is a windowed reverse over the working array → see [`/dsa/arrays.md`](/dsa/arrays.md)

---

## Approaches Overview

| # | Approach | Time | Space | When to use |
|---|----------|------|-------|-------------|
| 1 | Brute Force via next_permutation Scan | O(n! · n) | O(n) | Tiny `n` / correctness oracle only |
| 2 | Greedy Reverse of D-Runs (Optimal) | O(n) | O(n) | Cleanest optimal; reverse each D-block over the identity |
| 3 | Stack-Based Emit-on-I (Optimal) | O(n) | O(n) | Single pass; descending blocks fall out of LIFO pops |

> Both optimal approaches encode the same idea — descending numbers within each `'D'` block — either by reversing a window (2) or by delaying output through a stack (3).

---

## Approach 1 — Brute Force via next_permutation Scan

### Intuition

We want the smallest permutation satisfying the pattern, so enumerate permutations of `[1..n]` in lexicographic order (identity first) and return the first that fits `s`. Correct by definition, and handy as an oracle to validate the greedy result for small `n` — but factorial time, so only viable for tiny inputs.

### Algorithm

1. `perm = [1,2,…,n]` (the lexicographically first permutation).
2. Loop: if `perm` satisfies every I/D constraint, return it; otherwise advance `perm` to its next lexicographic permutation.
3. Stop when permutations are exhausted (never needed for valid input — the identity already satisfies an all-`'I'` string).

### Complexity

- **Time:** O(n! · n) worst case — up to `n!` permutations, each validated in `O(n)`.
- **Space:** O(n) — the working permutation.

### Code

```go
func bruteForce(s string) []int {
	n := len(s) + 1
	perm := make([]int, n)
	for i := range perm {
		perm[i] = i + 1 // identity = lexicographically first permutation
	}

	for {
		if satisfies(perm, s) { // first fit in lex order is the smallest
			return perm
		}
		if !nextPermutation(perm) { // exhausted all permutations
			return perm // (unreachable for valid input; identity handles all-'I')
		}
	}
}

func satisfies(perm []int, s string) bool {
	for i := 0; i < len(s); i++ {
		if s[i] == 'I' && !(perm[i] < perm[i+1]) {
			return false // 'I' demands an ascent here
		}
		if s[i] == 'D' && !(perm[i] > perm[i+1]) {
			return false // 'D' demands a descent here
		}
	}
	return true
}

func nextPermutation(perm []int) bool {
	n := len(perm)
	i := n - 2
	for i >= 0 && perm[i] >= perm[i+1] { // find last ascent perm[i] < perm[i+1]
		i--
	}
	if i < 0 {
		return false // wholly descending → no greater permutation
	}
	j := n - 1
	for perm[j] <= perm[i] { // find rightmost element greater than perm[i]
		j--
	}
	perm[i], perm[j] = perm[j], perm[i] // swap pivot with its successor
	reverseInts(perm, i+1, n-1)         // reverse the suffix to its smallest order
	return true
}
```

### Dry Run

Example 1: `s = "I"`, `n = 2`.

| Step | perm | satisfies "I" (perm[0] < perm[1])? | action |
|------|------|-------------------------------------|--------|
| 1 | `[1,2]` | 1 < 2 → yes | **return [1,2]** |

Result: `[1,2]` ✔ (found immediately — identity fits an all-`'I'` pattern).

---

## Approach 2 — Greedy Reverse of D-Runs (Optimal, array walk)

### Intuition

Start from `[1,2,…,n]`, which is already the lexicographically smallest sequence and satisfies every `'I'`. Wherever `s` holds a maximal run of `'D'` at positions `i..j`, the array indices `i..j+1` must strictly decrease. Those slots currently hold the consecutive ascending numbers `i+1..j+2`; **reversing exactly that window** turns them descending while leaving every element outside untouched — the minimal change, hence still lex-smallest. Isolated `'I'` positions need nothing. Sweeping left to right and reversing each `'D'` block gives the global answer.

### Algorithm

1. `perm = [1,2,…,n]`.
2. `i = 0`. While `i < len(s)`: if `s[i] == 'D'`, extend to the run's end `j` (last consecutive `'D'`), reverse `perm[i .. j]` (indices `i` through `j`, i.e. `j+1` elements), then set `i = j`; otherwise `i++`.
3. Return `perm`.

### Complexity

- **Time:** O(n) — the `'D'`-runs are disjoint, so each element is reversed at most once.
- **Space:** O(n) — the output permutation (`O(1)` extra beyond it).

### Code

```go
func greedyReverseRuns(s string) []int {
	n := len(s) + 1
	perm := make([]int, n)
	for i := range perm {
		perm[i] = i + 1 // baseline ascending permutation
	}

	i := 0
	for i < len(s) {
		if s[i] == 'D' {
			j := i
			for j < len(s) && s[j] == 'D' { // extend over the whole D-run
				j++
			}
			// s[i..j-1] are 'D' ⇒ indices i..j must descend; reverse that window.
			reverseInts(perm, i, j)
			i = j // continue past the run (position j was the boundary)
		} else {
			i++ // 'I' already satisfied by the ascending baseline
		}
	}
	return perm
}
```

### Dry Run

Example 2: `s = "DI"`, `n = 3`, baseline `perm = [1,2,3]`.

| i | s[i] | action | run end j | window reversed | perm after | next i |
|---|------|--------|-----------|-----------------|------------|--------|
| 0 | `D` | start D-run | j: s[0]=D→1, s[1]=I stop ⇒ j=1 | reverse perm[0..1] `[1,2]→[2,1]` | `[2,1,3]` | 1 |
| 1 | `I` | nothing (baseline ok) | — | — | `[2,1,3]` | 2 |
| — | end | — | — | — | `[2,1,3]` | — |

Result: `[2,1,3]` ✔

---

## Approach 3 — Stack-Based Emit-on-I (Optimal, one pass)

### Intuition

Push `1,2,3,…` onto a stack. A stack reverses order on pop, so any numbers held between flushes emerge **descending** — precisely what a `'D'` run needs. At every `'I'` boundary (and after pushing the final number) we pop everything queued so far into the output. This delayed LIFO flush reproduces the reverse-D-run result in a single left-to-right pass, with no explicit reversal call.

### Algorithm

1. For `pos = 0 … n-1`: push `pos+1` onto the stack.
2. If `pos == n-1` **or** `s[pos] == 'I'`: pop the entire stack into the output (the queued numbers form one descending block).
3. The output, filled in pop order, is the lex-smallest permutation.

### Complexity

- **Time:** O(n) — each number is pushed once and popped once.
- **Space:** O(n) — stack plus output.

### Code

```go
func stackEmit(s string) []int {
	n := len(s) + 1
	out := make([]int, 0, n)
	stack := make([]int, 0, n)

	for pos := 0; pos < n; pos++ {
		stack = append(stack, pos+1) // queue the next natural number
		// Flush at an 'I' boundary or at the very end: everything queued since
		// the last flush belongs to one descending block, reversed by the pops.
		if pos == n-1 || s[pos] == 'I' {
			for len(stack) > 0 {
				out = append(out, stack[len(stack)-1]) // pop top (LIFO ⇒ reversed)
				stack = stack[:len(stack)-1]
			}
		}
	}
	return out
}
```

### Dry Run

Example 2: `s = "DI"`, `n = 3`.

| pos | push | s[pos] | flush? (pos==n-1 or 'I') | pops → out | stack after |
|-----|------|--------|--------------------------|------------|-------------|
| 0 | 1 | `D` | no | — | `[1]` |
| 1 | 2 | `I` | yes | pop 2, pop 1 → `[2,1]` | `[]` |
| 2 | 3 | (end) | yes (pos==2) | pop 3 → `[2,1,3]` | `[]` |

Result: `[2,1,3]` ✔ — LIFO turned the pending `[1,2]` into the descent `[2,1]`.

---

## Key Takeaways

- **Identity is the lex-smallest baseline.** For "smallest permutation matching a comparison pattern", start from `[1..n]` and disturb it as little as possible.
- **A `'D'` run ⇔ reverse a contiguous window.** Consecutive decreases over indices `i..j` are satisfied by reversing exactly that ascending block — nothing outside changes, preserving minimality.
- **Stacks encode "reverse on demand."** Pushing in natural order and popping at boundaries reverses each pending block for free — a one-pass alternative to explicit window reversal.
- Both optimal forms are `O(n)`; pick the reverse-runs version for clarity, the stack version when you like streaming/one-pass phrasing.

---

## Related Problems

- LeetCode #942 — DI String Match (construct any permutation from an I/D string)
- LeetCode #46 — Permutations (enumerate all permutations)
- LeetCode #31 — Next Permutation (the lexicographic-successor building block)
- LeetCode #316 — Remove Duplicate Letters (greedy + stack for lex-smallest output)
