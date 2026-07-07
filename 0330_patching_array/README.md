# 0330 — Patching Array

> LeetCode #330 · Difficulty: Hard
> **Categories:** Array, Greedy

---

## Problem Statement

Given a sorted integer array `nums` and an integer `n`, add/patch elements to the array such that any number in the range `[1, n]` inclusive can be formed by the sum of some elements in the array.

Return *the minimum number of patches required*.

**Example 1:**

```
Input: nums = [1,3], n = 6
Output: 1
Explanation:
Combinations of nums are [1], [3], [1,3], which form possible sums of: 1, 3, 4.
Now if we add/patch 2 to nums, the combinations are: [1], [2], [3], [1,3], [2,3], [1,2,3].
Possible sums are 1, 2, 3, 4, 5, 6, which now covers the range [1, 6].
So we only need 1 patch.
```

**Example 2:**

```
Input: nums = [1,5,10], n = 20
Output: 2
Explanation: The two patches can be [2, 4].
```

**Example 3:**

```
Input: nums = [1,2,2], n = 5
Output: 0
```

**Constraints:**

- `1 <= nums.length <= 1000`
- `1 <= nums[i] <= 10^4`
- `nums` is sorted in ascending order.
- `1 <= n <= 2^31 - 1`

---

## Company Frequency

| Company    | Frequency        | Last Reported |
|------------|------------------|---------------|
| Google     | ★★★★☆ High       | 2024          |
| Amazon     | ★★★☆☆ Medium     | 2024          |
| Microsoft  | ★★★☆☆ Medium     | 2023          |
| Meta       | ★★☆☆☆ Low        | 2023          |
| Apple      | ★★☆☆☆ Low        | 2022          |
| Bloomberg  | ★★☆☆☆ Low        | 2022          |

> ⚠️ Frequency data is crowd-sourced from LeetCode Discuss, Glassdoor, and
> community interview reports. Treat as a signal, not a guarantee.

---

## DSA Concepts Used

- **Greedy** — at each gap we make the locally-best choice (patch the smallest currently-unreachable value `miss`), and this local choice is provably globally optimal because no larger value can fill the hole at `miss` and no smaller value extends coverage further → see [`/dsa/greedy.md`](/dsa/greedy.md)
- **Number Theory / Reachable Subset Sums** — the core lemma is: if every integer in `[1, k]` is a reachable subset sum, then adding any `x <= k+1` extends reachability contiguously to `[1, k+x]`; patching `k+1` doubles the range to `[1, 2k+1]`. This "coverage doubling" is the number-theoretic engine of the algorithm → see [`/dsa/math_number_theory.md`](/dsa/math_number_theory.md)

---

## Approaches Overview

| # | Approach | Time | Space | When to use |
|---|----------|------|-------|-------------|
| 1 | Greedy Patching (Optimal) | O(m + log n) | O(1) | The canonical answer; `miss` tracks the smallest unreachable value directly |
| 2 | Greedy Verbose (Explicit Coverage Bound) | O(m + log n) | O(1) | Same greedy re-expressed with a `covered` upper bound; clearest invariant, used to cross-check |

(m = `len(nums)`. Both are the same optimal greedy, structured differently for clarity and mutual verification — the optimal solution to this problem is uniquely greedy.)

---

## Approach 1 — Greedy Patching (Optimal)

### Intuition

Track `miss` = the smallest positive integer we cannot yet form as a subset sum. The invariant is that **everything in `[1, miss-1]` is already reachable**. Now consider the next available number `nums[i]`:

- If `nums[i] <= miss`, it slots onto the reachable prefix with no gap: reachability extends from `[1, miss-1]` to `[1, miss-1+nums[i]]`, so `miss` grows to `miss + nums[i]` — for free, no patch.
- If `nums[i] > miss` (or we have run out of numbers), nothing available can bridge the hole at `miss`. The optimal move is to patch `miss` itself. Adding `miss` extends reachability to `[1, 2*miss-1]` — the largest jump any single patch can make — and costs one patch.

Why patching `miss` is optimal: any patch value `> miss` leaves the hole at `miss` unfilled; any value `< miss` is already reachable and extends coverage less. So `miss` is simultaneously the only value that fixes the current gap and the one that reaches furthest.

### Algorithm

1. Initialise `miss = 1`, `i = 0`, `patches = 0`.
2. While `miss <= n`:
   1. If `i < len(nums)` and `nums[i] <= miss`: set `miss += nums[i]`; `i++` (free extension).
   2. Else: set `miss += miss` (double coverage by patching `miss`); `patches++`.
3. Return `patches`.

### Complexity

- **Time:** O(m + log n) — each of the `m = len(nums)` numbers is consumed at most once, and every patch at least doubles `miss`, so at most about `log2(n)` patches occur.
- **Space:** O(1) — just `miss`, `i`, and `patches`.

### Code

```go
func greedyPatching(nums []int, n int) int {
	// miss is int64 because doubling can push it past 2^31 before the loop's
	// `miss <= n` guard stops it (n itself can be up to 2^31-1). A 32-bit int
	// would overflow on that last doubling; int64 keeps it exact.
	var miss int64 = 1
	i := 0       // index into nums
	patches := 0 // count of numbers we had to add

	for miss <= int64(n) { // continue until every value in [1, n] is reachable
		if i < len(nums) && int64(nums[i]) <= miss {
			// nums[i] plugs into the current coverage without leaving a gap:
			// extend reachability by nums[i] at no cost.
			miss += int64(nums[i])
			i++
		} else {
			// Gap at `miss`: patch it. Adding `miss` doubles the reachable
			// range (best possible jump) and consumes exactly one patch.
			miss += miss // == miss *= 2; may exceed 2^31, hence int64
			patches++
		}
	}
	return patches
}
```

### Dry Run

Example 1: `nums = [1, 3]`, `n = 6`.

| Step | miss (before) | i | nums[i] | nums[i] <= miss? | Action | miss (after) | patches |
|------|---------------|---|---------|------------------|--------|--------------|---------|
| 0 | 1 | 0 | 1 | yes (1 <= 1) | consume nums[0]=1, i→1 | 2 | 0 |
| 1 | 2 | 1 | 3 | no (3 > 2) | patch 2 (double) | 4 | 1 |
| 2 | 4 | 1 | 3 | yes (3 <= 4) | consume nums[1]=3, i→2 | 7 | 1 |
| 3 | 7 | 2 | — | i out of range | `miss=7 > n=6`, exit loop | 7 | 1 |

`miss = 7 > 6`, loop ends. Result: **1** ✔ — exactly the single patch (value `2`) described in the problem.

---

## Approach 2 — Greedy Verbose (Explicit Coverage Bound)

### Intuition

Identical greedy, but instead of tracking the first *unreachable* value we track the largest *reachable* one. Let `covered` = the largest value such that every integer in `[0, covered]` is a reachable subset sum; start at `covered = 0` (only the empty sum `0` is reachable). Then `covered + 1` is the first unreachable value — exactly `miss` from Approach 1.

At each step, look at the smallest unreachable value `next = covered + 1`:

- If `nums[i] <= next`, the number attaches to the covered prefix with no gap, so `covered` grows to `covered + nums[i]` (no patch).
- Otherwise patch `next`: the reachable range extends to `[0, covered + next] = [0, 2*covered + 1]`, spending one patch.

Stop once `covered >= n`, i.e. `[1, n]` is fully covered. Making the `covered` bound explicit spells the loop invariant out in the open, which is why it is handy as an independent cross-check of Approach 1.

### Algorithm

1. Initialise `covered = 0`, `i = 0`, `patches = 0`.
2. While `covered < n`:
   1. Let `next = covered + 1` (smallest currently-unreachable value).
   2. If `i < len(nums)` and `nums[i] <= next`: set `covered += nums[i]`; `i++`.
   3. Else: set `covered += next` (patch `next`); `patches++`.
3. Return `patches`.

### Complexity

- **Time:** O(m + log n) — same accounting as Approach 1: each number consumed once, each patch roughly doubles `covered`.
- **Space:** O(1) — `covered`, `i`, `patches`.

### Code

```go
func greedyVerbose(nums []int, n int) int {
	// covered = upper bound of the contiguous reachable range [0, covered].
	// int64 for the same overflow reason: covered can approach 2^31, and adding
	// `next` (~covered+1) nearly doubles it past the 32-bit limit.
	var covered int64 = 0
	i := 0
	patches := 0

	for covered < int64(n) { // stop once [1, n] is fully covered
		next := covered + 1 // smallest value not yet reachable == "miss"
		if i < len(nums) && int64(nums[i]) <= next {
			// nums[i] attaches to the reachable prefix with no gap.
			covered += int64(nums[i])
			i++
		} else {
			// Patch `next`; reachable range extends to [0, covered+next].
			covered += next // ~doubles covered; may exceed 2^31, hence int64
			patches++
		}
	}
	return patches
}
```

### Dry Run

Example 1: `nums = [1, 3]`, `n = 6`.

| Step | covered (before) | next = covered+1 | i | nums[i] | nums[i] <= next? | Action | covered (after) | patches |
|------|------------------|------------------|---|---------|------------------|--------|-----------------|---------|
| 0 | 0 | 1 | 0 | 1 | yes (1 <= 1) | consume nums[0]=1, i→1 | 1 | 0 |
| 1 | 1 | 2 | 1 | 3 | no (3 > 2) | patch 2 | 3 | 1 |
| 2 | 3 | 4 | 1 | 3 | yes (3 <= 4) | consume nums[1]=3, i→2 | 6 | 1 |
| 3 | 6 | — | 2 | — | — | `covered=6 >= n=6`, exit loop | 6 | 1 |

`covered = 6 >= 6`, loop ends. Result: **1** ✔ — agrees with Approach 1 and the problem's stated answer.

---

## Key Takeaways

- **Coverage-doubling invariant.** If `[1, k]` is fully reachable and you add the value `k+1`, reachability jumps to `[1, 2k+1]` — the biggest single-step expansion possible. Patching the smallest gap `miss` (= `covered+1`) is therefore both necessary (only it fills the current hole) and optimal (nothing reaches further). This is the whole problem.
- **Greedy correctness by an exchange argument.** No value `> miss` can cover `miss`; no value `< miss` extends coverage as far. So the greedy choice can never be beaten by a different first patch, and the argument recurses.
- **Consume vs. patch.** A sorted `nums[i] <= miss` is a free extension — always take it before patching. Because `nums` is sorted, a single forward index suffices; no re-scanning.
- **Watch the overflow.** `miss` (or `covered`) can approach `2^31` and then double on the last patch, overshooting a 32-bit `int`. Use `int64` (or `uint64`) for the running bound; the *answer* (patch count ≤ ~31) still fits comfortably in `int`.
- **Range width is irrelevant to cost.** `n = 2^31 - 1` costs only ~`log2(n)` ≈ 30 patches, not `O(n)` work — the loop count is governed by doublings, not by `n`.

---

## Related Problems

- LeetCode #45 — Jump Game II (greedy reach-extension, minimise steps)
- LeetCode #55 — Jump Game (greedy reachability frontier)
- LeetCode #134 — Gas Station (greedy running-deficit argument)
- LeetCode #763 — Partition Labels (greedy interval-extension over a frontier)
