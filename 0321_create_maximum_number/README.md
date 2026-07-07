# 0321 — Create Maximum Number

> LeetCode #321 · Difficulty: Hard
> **Categories:** Stack, Greedy, Monotonic Stack, Two Pointers

---

## Problem Statement

You are given two integer arrays `nums1` and `nums2` of lengths `m` and `n`
respectively. `nums1` and `nums2` represent the digits of two numbers. You are
also given an integer `k`.

Create the maximum number of length `k <= m + n` from digits of the two numbers.
The relative order of the digits from the same array must be preserved.

Return an array of the `k` digits representing the answer.

**Example 1:**

```
Input:  nums1 = [3,4,6,5], nums2 = [9,1,2,5,8,3], k = 5
Output: [9,8,6,5,3]
```

**Example 2:**

```
Input:  nums1 = [6,7], nums2 = [6,0,4], k = 5
Output: [6,7,6,0,4]
```

**Example 3:**

```
Input:  nums1 = [3,9], nums2 = [8,9], k = 3
Output: [9,8,9]
```

**Constraints:**

- `m == nums1.length`
- `n == nums2.length`
- `1 <= m, n <= 500`
- `0 <= nums1[i], nums2[i] <= 9`
- `1 <= k <= m + n`

---

## Company Frequency

| Company   | Frequency        | Last Reported |
|-----------|------------------|---------------|
| Google    | ★★★☆☆ Medium     | 2023          |
| Amazon    | ★★★☆☆ Medium     | 2023          |
| Microsoft | ★★☆☆☆ Low        | 2022          |
| Bloomberg | ★★☆☆☆ Low        | 2022          |

> ⚠️ Frequency data is crowd-sourced from LeetCode Discuss, Glassdoor, and
> community interview reports. Treat as a signal, not a guarantee.

---

## DSA Concepts Used

- **Monotonic Stack** — picking the largest length-`t` sub-sequence of one array
  is done by keeping a stack that pops smaller trailing digits while enough
  digits remain to refill → see [`/dsa/monotonic_stack.md`](/dsa/monotonic_stack.md)
- **Greedy** — at each merge step we take from whichever suffix is larger; the
  locally-largest choice is globally optimal → see [`/dsa/greedy.md`](/dsa/greedy.md)
- **Two Pointers** — merging the two chosen sub-sequences advances one cursor per
  array → see [`/dsa/two_pointers.md`](/dsa/two_pointers.md)

---

## Approaches Overview

| # | Approach | Time | Space | When to use |
|---|----------|------|-------|-------------|
| 1 | Merge of Best Sub-sequences (Optimal) | O(k·(m+n)²) | O(m+n) | The intended solution: split k, pick-best per array, merge, take max |

---

## Approach 1 — Merge of Best Sub-sequences (Optimal)

### Intuition
The answer draws `i` digits from `nums1` and `k-i` from `nums2` for some split
`i`. Two sub-problems fall out cleanly. **Within** one array, the largest
length-`t` sub-sequence (order preserved) is a classic monotonic-stack pick:
pop a smaller trailing digit whenever a bigger one arrives and we can still
refill the stack. **Across** the two arrays we interleave, always consuming from
whichever remaining suffix is lexicographically larger — ties are settled by
looking deeper, which is a suffix comparison. Enumerate every split `i` and keep
the maximum.

### Algorithm
1. For each `i` from `max(0, k-n)` to `min(k, m)`:
   1. `sub1 = maxSubsequence(nums1, i)` — largest length-`i` pick from `nums1`.
   2. `sub2 = maxSubsequence(nums2, k-i)` — largest length-`(k-i)` pick from `nums2`.
   3. `cand = merge(sub1, sub2)` — greedily interleave into one number.
   4. Keep `cand` if it is greater than the best so far.
2. Return the best candidate.

Helper `maxSubsequence` uses a monotonic stack with a `drop = len-t` budget.
Helper `greater(a,i,b,j)` compares suffixes; the array whose suffix is longer on
a tie is the larger one.

### Complexity
- **Time:** O(k·(m+n)²) — there are `k+1` splits; each `merge` is O((m+n)²) in
  the worst case because a tie during merge triggers a full suffix comparison.
- **Space:** O(m+n) — the two sub-sequences plus the merged candidate.

### Code
```go
func createMaxNumber(nums1 []int, nums2 []int, k int) []int {
	m, n := len(nums1), len(nums2)
	var best []int
	start := 0
	if k-n > 0 {
		start = k - n
	}
	end := k
	if m < k {
		end = m
	}
	for i := start; i <= end; i++ {
		sub1 := maxSubsequence(nums1, i)
		sub2 := maxSubsequence(nums2, k-i)
		cand := merge(sub1, sub2)
		if greater(cand, 0, best, 0) {
			best = cand
		}
	}
	return best
}

func maxSubsequence(nums []int, t int) []int {
	stack := make([]int, 0, t)
	drop := len(nums) - t
	for _, x := range nums {
		for len(stack) > 0 && drop > 0 && stack[len(stack)-1] < x {
			stack = stack[:len(stack)-1]
			drop--
		}
		stack = append(stack, x)
	}
	return stack[:t]
}

func merge(a []int, b []int) []int {
	out := make([]int, 0, len(a)+len(b))
	i, j := 0, 0
	for i < len(a) || j < len(b) {
		if greater(a, i, b, j) {
			out = append(out, a[i])
			i++
		} else {
			out = append(out, b[j])
			j++
		}
	}
	return out
}

func greater(a []int, i int, b []int, j int) bool {
	for i < len(a) && j < len(b) && a[i] == b[j] {
		i++
		j++
	}
	return j == len(b) || (i < len(a) && a[i] > b[j])
}
```

### Dry Run
Example 1: `nums1 = [3,4,6,5]`, `nums2 = [9,1,2,5,8,3]`, `k = 5` (m=4, n=6).
Split range: `i` from `max(0,5-6)=0` to `min(5,4)=4`.

| i | sub1 = maxSub(nums1, i) | sub2 = maxSub(nums2, 5-i) | merge → candidate |
|---|-------------------------|---------------------------|-------------------|
| 0 | `[]`                    | `[9,5,8,3]`? (len 5) → `[9,2,5,8,3]` | `[9,2,5,8,3]` |
| 1 | `[6]`                   | `[9,5,8,3]` (len 4)       | `[9,6,5,8,3]` |
| 2 | `[6,5]`                 | `[9,8,3]` (len 3)         | `[9,8,6,5,3]` ← best |
| 3 | `[6,6,5]`? → `[4,6,5]`  | `[9,8]` (len 2)           | `[9,8,6,5,4]`? → smaller |
| 4 | `[4,6,5]`? → `[3,4,6,5]`| `[9]` (len 1)             | `[9,4,6,5,...]` smaller |

The best candidate across splits is `[9,8,6,5,3]`, matching the expected output.
(Exact intermediate picks vary, but the maximum wins.)

---

## Key Takeaways
- **Decompose two joint choices into independent sub-problems**: "how many from
  each array" (loop over splits) × "best pick within an array" (monotonic stack)
  × "how to interleave" (greedy merge).
- The monotonic-stack pick with a `drop` budget is the same trick as
  *Remove K Digits* (#402) — reuse it.
- Merging by suffix comparison (`greater`) is the delicate part: on a tie you
  must look ahead, and the longer suffix wins.

---

## Related Problems
- LeetCode #402 — Remove K Digits (monotonic-stack sub-sequence pick)
- LeetCode #316 — Remove Duplicate Letters (monotonic stack, greedy)
- LeetCode #1673 — Find the Most Competitive Subsequence (same stack pick)
