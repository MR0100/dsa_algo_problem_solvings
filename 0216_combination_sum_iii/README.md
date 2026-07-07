# 0216 — Combination Sum III

> LeetCode #216 · Difficulty: Medium
> **Categories:** Array, Backtracking

---

## Problem Statement

Find all valid combinations of `k` numbers that sum up to `n` such that the following conditions are true:

- Only numbers `1` through `9` are used.
- Each number is used **at most once**.

Return *a list of all possible valid combinations*. The list must not contain the same combination twice, and the combinations may be returned in any order.

**Example 1:**

```
Input: k = 3, n = 7
Output: [[1,2,4]]
Explanation:
1 + 2 + 4 = 7
There are no other valid combinations.
```

**Example 2:**

```
Input: k = 3, n = 9
Output: [[1,2,6],[1,3,5],[2,3,4]]
Explanation:
1 + 2 + 6 = 9
1 + 3 + 5 = 9
2 + 3 + 4 = 9
There are no other valid combinations.
```

**Example 3:**

```
Input: k = 4, n = 1
Output: []
Explanation: There are no valid combinations.
Using 4 different numbers in the range [1,9], the smallest sum we can get is 1+2+3+4 = 10 and since 10 > 1, there are no valid combinations.
```

**Constraints:**

- `2 <= k <= 9`
- `1 <= n <= 60`

---

## Company Frequency

| Company    | Frequency        | Last Reported |
|------------|------------------|---------------|
| Amazon     | ★★★☆☆ Medium     | 2024          |
| Google     | ★★★☆☆ Medium     | 2023          |
| Microsoft  | ★★☆☆☆ Low        | 2023          |
| Facebook   | ★★☆☆☆ Low        | 2022          |
| Adobe      | ★★☆☆☆ Low        | 2022          |

> ⚠️ Frequency data is crowd-sourced from LeetCode Discuss, Glassdoor, and
> community interview reports. Treat as a signal, not a guarantee.

---

## DSA Concepts Used

- **Backtracking** — the canonical "choose / explore / un-choose" template over a fixed candidate set, using a `start` index to enforce sorted, duplicate-free combinations and pruning to cut dead branches → see [`/dsa/backtracking.md`](/dsa/backtracking.md)
- **Bit Manipulation** — the alternate approach enumerates all `2^9` subsets of the 9 digits as bitmasks, one bit per digit → see [`/dsa/bit_manipulation.md`](/dsa/bit_manipulation.md)

---

## Approaches Overview

| # | Approach | Time | Space | When to use |
|---|----------|------|-------|-------------|
| 1 | Backtracking (Optimal) | O(C(9,k)·k) | O(k) | The standard, extensible answer; generalizes to larger candidate sets |
| 2 | Bitmask Enumeration | O(2⁹·9) | O(k) | Only works because the universe is tiny (9 fixed digits); neat and branch-free |

---

## Approach 1 — Backtracking (Optimal)

### Intuition

We need exactly `k` **distinct** digits from `1..9` summing to `n`, order irrelevant. Walking the digits left-to-right and always forcing the next chosen digit to be strictly larger than the current one means every combination is produced exactly once, already sorted — no duplicate filtering needed. Along the way we prune hard: if the current digit already exceeds the remaining sum, every larger digit does too, so we stop scanning; and once we have picked `k` digits we check the sum and return.

### Algorithm

1. DFS state: `start` (smallest digit still allowed), `count` (digits still to pick), `sum` (target still to reach), plus the shared `path`.
2. If `count == 0`: success iff `sum == 0` → record a **copy** of `path`; return either way.
3. Loop `d` from `start` to `9`. If `d > sum`, `break` (sorted digits overshoot).
4. Choose `d` (append to `path`), recurse with `(d+1, count-1, sum-d)`, then pop `d` off `path`.
5. Start the recursion at `dfs(1, k, n)`.

### Complexity

- **Time:** O(C(9,k)·k) — at most C(9,k) valid leaves, each copied in O(k); the pruned search tree over 9 digits is small.
- **Space:** O(k) — recursion depth is bounded by `k`, and `path` holds at most `k` digits (output array excluded).

### Code

```go
func backtracking(k int, n int) [][]int {
	var result [][]int      // collected valid combinations
	path := make([]int, 0, k) // current partial combination being built

	// dfs tries digits from `start` upward, needing `count` more digits that
	// together sum to `sum`.
	var dfs func(start, count, sum int)
	dfs = func(start, count, sum int) {
		if count == 0 { // we have chosen exactly k digits
			if sum == 0 { // and they hit the target sum
				combo := make([]int, len(path)) // copy: path is mutated after return
				copy(combo, path)
				result = append(result, combo)
			}
			return // either recorded or over-shot the sum with k digits — stop
		}
		// Try each candidate digit in strictly increasing order.
		for d := start; d <= 9; d++ {
			if d > sum { // digits are sorted; d and everything after overshoots
				break
			}
			path = append(path, d)       // choose d
			dfs(d+1, count-1, sum-d)     // recurse: next digit must exceed d
			path = path[:len(path)-1]    // un-choose d (backtrack)
		}
	}

	dfs(1, k, n) // start from digit 1, needing k digits summing to n
	return result
}
```

### Dry Run

Example 1: `k = 3, n = 7`. State shown as `(start, count, sum)` with `path`.

| Step | Call | path | Action |
|------|------|------|--------|
| 1 | dfs(1,3,7) | [] | try d=1 → choose |
| 2 | dfs(2,2,6) | [1] | try d=2 → choose |
| 3 | dfs(3,1,4) | [1,2] | try d=3 → choose |
| 4 | dfs(4,0,1) | [1,2,3] | count=0 but sum=1≠0 → dead; pop 3 |
| 5 | dfs(3,1,4) | [1,2] | try d=4 → choose |
| 6 | dfs(5,0,0) | [1,2,4] | count=0 and sum=0 → **record [1,2,4]**; pop 4 |
| 7 | dfs(3,1,4) | [1,2] | try d=5..: 5>sum? no, but dfs(6,0,-1)… actually d=5→sum-5=-1, recorded none; d>sum breaks at d=5>4 → break; pop 2 |
| 8 | dfs(2,2,6) | [1] | try d=3 → dfs(4,1,3): d=4>3 break → dead … continues, none succeed |
| … | | | all remaining branches prune out |

Result: `[[1,2,4]]` ✔

---

## Approach 2 — Bitmask Enumeration

### Intuition

The candidate universe is fixed and tiny: the nine digits `1..9`. There are only `2^9 = 512` subsets, so we can just enumerate every subset with a 9-bit integer mask — bit `i` set means "digit `i+1` is in the subset". For each mask we compute its size (popcount) and the sum of its digits, and keep it exactly when size equals `k` and sum equals `n`. No recursion, no pruning logic — brute force is affordable here.

### Algorithm

1. For `mask` from `0` to `511`:
2.   Scan bits `i = 0..8`; whenever bit `i` is set, add digit `i+1` to a running `sum` and to a temp `combo`.
3.   If `len(combo) == k` and `sum == n`, append `combo` to the result.
4. Return the result (each combo internally sorted; combos come out in mask order).

### Complexity

- **Time:** O(2⁹·9) ≈ 4608 constant operations — independent of `k` and `n`.
- **Space:** O(k) scratch per mask for `combo` (output array excluded).

### Code

```go
func bitmaskEnumeration(k int, n int) [][]int {
	var result [][]int
	// 1<<9 == 512 subsets of the digit set {1..9}.
	for mask := 0; mask < (1 << 9); mask++ {
		sum := 0                  // running sum of chosen digits
		combo := make([]int, 0, k) // chosen digits for this mask
		for i := 0; i < 9; i++ {  // inspect each of the 9 candidate digits
			if mask&(1<<i) != 0 { // bit i set → digit (i+1) is included
				digit := i + 1
				sum += digit
				combo = append(combo, digit)
			}
		}
		// Keep only subsets of the right size that hit the target sum.
		if len(combo) == k && sum == n {
			result = append(result, combo)
		}
	}
	return result
}
```

### Dry Run

Example 1: `k = 3, n = 7`. A few representative masks (bit `i` ↔ digit `i+1`):

| mask (binary, bits 8..0) | digits chosen | size | sum | keep? |
|--------------------------|---------------|------|-----|-------|
| `000000111` | 1,2,3 | 3 | 6 | no (sum 6≠7) |
| `000001011` | 1,2,4 | 3 | 7 | **yes** → [1,2,4] |
| `000010011` | 1,2,5 | 3 | 8 | no |
| `000000011` | 1,2 | 2 | 3 | no (size 2≠3) |
| `000001101` | 1,3,4 | 3 | 8 | no |

Scanning all 512 masks, only `000001011` satisfies both conditions. Result: `[[1,2,4]]` ✔

---

## Key Takeaways

- **`start` index = duplicate-free sorted combinations.** Passing `d+1` into the recursive call is the one-line trick that both forbids reusing a number and guarantees each combination is emitted once, already sorted.
- **Prune with sorted candidates.** Because digits are visited in increasing order, `d > remaining sum` lets you `break` the whole remaining loop, not just `continue`.
- **Always copy the path before recording.** `path` is a shared, mutated buffer; appending it directly would store an alias that later backtracking corrupts.
- **When the universe is tiny (≤ ~20 elements), bitmask enumeration is a legitimate, bug-resistant alternative** — no recursion, trivial to reason about, and here it is effectively O(1).

---

## Related Problems

- LeetCode #39 — Combination Sum (reuse allowed, unbounded candidates)
- LeetCode #40 — Combination Sum II (candidates with duplicates, each used once)
- LeetCode #77 — Combinations (choose k of n, no sum constraint)
- LeetCode #78 — Subsets (all subsets — same bitmask/backtracking skeleton)
- LeetCode #46 — Permutations (order matters — no `start` index)
