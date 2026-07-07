# 0440 — K-th Smallest in Lexicographical Order

> LeetCode #440 · Difficulty: Hard
> **Categories:** Trie, Greedy, Math

---

## Problem Statement

Given two integers `n` and `k`, return *the* `kth` *lexicographically smallest integer in the range* `[1, n]`.

**Example 1:**

```
Input: n = 13, k = 2
Output: 10
Explanation: The lexicographical order is [1, 10, 11, 12, 13, 2, 3, 4, 5, 6, 7, 8, 9], so the second smallest number is 10.
```

**Example 2:**

```
Input: n = 1, k = 1
Output: 1
```

**Constraints:**

- `1 <= k <= n <= 10^9`

---

## Company Frequency

| Company    | Frequency        | Last Reported |
|------------|------------------|---------------|
| ByteDance  | ★★★☆☆ Medium     | 2024          |
| Google     | ★★★☆☆ Medium     | 2023          |
| Amazon     | ★★☆☆☆ Low        | 2023          |
| Microsoft  | ★★☆☆☆ Low        | 2022          |
| Apple      | ★★☆☆☆ Low        | 2022          |

> ⚠️ Frequency data is crowd-sourced from LeetCode Discuss, Glassdoor, and
> community interview reports. Treat as a signal, not a guarantee.

---

## DSA Concepts Used

- **Trie (denary prefix tree)** — the numbers `1..n` form a 10-ary tree where node `p` has children `p·10 .. p·10+9`; a pre-order walk of it yields lexicographic order, and counting nodes under a prefix lets us skip whole subtrees → see [`/dsa/trie.md`](/dsa/trie.md)
- **Greedy descent** — at each step we greedily either skip an entire sibling subtree or descend into the current one, guided by the count, never backtracking → see [`/dsa/greedy.md`](/dsa/greedy.md)
- **Math / counting** — `countUnder(prefix)` is a closed-form level-by-level range count (`min(next, n+1) − cur`), not an enumeration → see [`/dsa/math_number_theory.md`](/dsa/math_number_theory.md)

---

## Approaches Overview

| # | Approach | Time | Space | When to use |
|---|----------|------|-------|-------------|
| 1 | Brute Force (Generate + Sort) | O(n log n · L) | O(n) | Only small n; the literal definition and correctness oracle |
| 2 | Prefix-Tree Step Counting (Optimal) | O((log n)²) | O(1) | The required solution; handles n = 10⁹ instantly |

`L` = number of digits.

---

## Approach 1 — Brute Force (Generate + Sort)

### Intuition

Lexicographic order *is* string order. Materialise every number `1..n`, sort by decimal spelling (so `"10"` precedes `"2"`), and read off element `k`. It is the definition made literal — perfectly correct, but it allocates and sorts up to `n = 10⁹` items, so it only survives small inputs. Its value here is as the oracle that Approach 2 must match.

### Algorithm

1. Build `nums = [1, 2, …, n]`.
2. Sort with comparator `strconv.Itoa(a) < strconv.Itoa(b)` (string comparison).
3. Return `nums[k-1]` (`k` is 1-indexed).

### Complexity

- **Time:** O(n log n · L) — `n log n` comparisons, each comparing up to `L` digits.
- **Space:** O(n) — the materialised slice plus transient string keys.

### Code

```go
func bruteForce(n int, k int) int {
	nums := make([]int, n) // all candidates 1..n
	for i := 0; i < n; i++ {
		nums[i] = i + 1
	}
	// Compare by decimal spelling so "10" < "2" (lexicographic, not numeric).
	sort.Slice(nums, func(a, b int) bool {
		return strconv.Itoa(nums[a]) < strconv.Itoa(nums[b])
	})
	return nums[k-1] // k is 1-indexed
}
```

### Dry Run

Example 1: `n = 13, k = 2`.

1. `nums = [1,2,3,4,5,6,7,8,9,10,11,12,13]`.
2. Sort by string → `["1","10","11","12","13","2","3","4","5","6","7","8","9"]` → `[1,10,11,12,13,2,3,4,5,6,7,8,9]`.
3. `nums[k-1] = nums[1] = 10`.

Result: `10` ✔

---

## Approach 2 — Prefix-Tree Step Counting (Optimal)

### Intuition

Picture `1..n` as a **denary trie**: the root's children are `1..9`, and every node `p` has children `p·10 … p·10+9` that do not exceed `n`. A **pre-order** traversal of this tree emits numbers in exactly lexicographic order (`1`, then `10`, `11`, … before `2`). We must not build or even walk every node — instead we **count**. `countUnder(prefix)` returns how many numbers in `[1, n]` begin with `prefix`. Now walk greedily from `prefix = 1` with a remaining budget of `k−1` steps:

- If the whole subtree under `prefix` holds `cnt ≤ remaining` numbers, the target is **not** inside — skip the entire subtree, move to the next sibling (`prefix+1`), and subtract `cnt`.
- Otherwise the target **is** inside — descend to the first child (`prefix·10`), spending one step.

When the budget reaches 0, `prefix` is the answer. Each "skip a subtree" collapses up to billions of numbers into one arithmetic step.

### Algorithm

1. `curr = 1`; `k--` (we already stand on the 1st lexicographic number, `"1"`).
2. While `k > 0`:
   - `cnt = countUnder(curr, n)`.
   - If `cnt <= k`: `curr++`, `k -= cnt` (skip this subtree, hop to next sibling).
   - Else: `curr *= 10`, `k--` (descend into this subtree).
3. Return `curr`.

`countUnder(prefix, n)` sums, level by level, the count of numbers of each length sharing `prefix`: start with range `[cur, next) = [prefix, prefix+1)`, add `min(next, n+1) − cur`, then multiply both edges by 10 to widen to the next level, until `cur > n`.

### Complexity

- **Time:** O((log n)²) — the outer greedy walk runs O(log n) big steps; each `countUnder` is O(log n). Completely independent of `k`'s size.
- **Space:** O(1) — only integer counters.

### Code

```go
func prefixTreeCount(n int, k int) int {
	curr := 1 // start at the smallest lexicographic number, prefix "1"
	k--       // standing on the 1st number already; k more steps to walk
	for k > 0 {
		cnt := countUnder(curr, n) // how many numbers in [1,n] begin with `curr`
		if cnt <= k {
			// The entire subtree under curr comes before the target; hop to the
			// next sibling and account for all cnt numbers we just skipped.
			curr++
			k -= cnt
		} else {
			// The target lives within curr's subtree; step down to its first
			// child (curr*10), spending one step to land on that number.
			curr *= 10
			k--
		}
	}
	return curr
}

func countUnder(prefix int, n int) int {
	count := 0
	cur := prefix      // left edge of the prefix's range at the current level
	next := prefix + 1 // right edge (exclusive) at the current level
	for cur <= n {
		// On this level the prefix covers [cur, next); clamp the right edge to
		// n+1 so we never count numbers greater than n.
		hi := next
		if n+1 < hi {
			hi = n + 1
		}
		count += hi - cur // numbers of this length sharing the prefix
		cur *= 10         // descend one level: ranges widen by a factor of 10
		next *= 10
	}
	return count
}
```

### Dry Run

Example 1: `n = 13, k = 2`. Start `curr = 1`, then `k-- → k = 1`.

First, `countUnder` values we will need (for `n = 13`):
- `countUnder(1)`: level `[1,2)` → 1; level `[10,20)` clamped to `[10,14)` → 4; next level `cur=100 > 13` stop. Total **5** (numbers 1,10,11,12,13).
- `countUnder(2)`: level `[2,3)` → 1; level `[20,30)` clamped to `[20,14)` → 0. Total **1** (just 2).

| iteration | curr | k (remaining) | cnt = countUnder(curr) | cnt ≤ k? | action | curr' | k' |
|-----------|------|---------------|------------------------|----------|--------|-------|----|
| 1 | 1 | 1 | 5 | no (5 > 1) | descend: `curr *= 10`, `k--` | 10 | 0 |

`k == 0` → stop. Return `curr = 10`.

Result: `10` ✔ — the walk descended from prefix `1` into `10` in a single step because the target (2nd number) lives inside the `1`-subtree, which contains 5 numbers.

---

## Key Takeaways

- **"K-th in lexicographic order" ⇒ pre-order over a denary trie.** Numbers `1..n` are a 10-ary tree (`p → p·10 … p·10+9`); their string order is that tree's pre-order.
- **Count subtrees, don't enumerate them.** `countUnder(prefix)` is a level-by-level clamped range sum, O(log n). It lets each step either skip a whole subtree (`curr++`, `k -= cnt`) or descend into it (`curr *= 10`, `k--`).
- **Two moves, one invariant.** `cnt ≤ k` means "target is past this subtree, skip it"; `cnt > k` means "target is inside, go deeper". The budget `k` monotonically shrinks to 0.
- **Runtime is independent of k and n's magnitude** — O((log n)²). The same idea powers "K-th smallest" queries over implicit ordered structures (e.g. #1439, #668, #719).
- **Clamp with `min(next, n+1)`** so numbers beyond `n` are never counted — the single easiest place to get an off-by-one.

---

## Related Problems

- LeetCode #386 — Lexicographical Numbers (the pre-order walk itself, no k-th)
- LeetCode #668 — Kth Smallest Number in Multiplication Table (count-and-descend / binary search on answer)
- LeetCode #719 — Find K-th Smallest Pair Distance (count ≤ x, binary search)
- LeetCode #1439 — Kth Smallest Sum in a Sorted Matrix (k-th over implicit order)
- LeetCode #440-style prefix counting reused in trie-of-numbers problems
