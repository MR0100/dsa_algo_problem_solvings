# 0280 — Wiggle Sort

> LeetCode #280 · Difficulty: Medium
> **Categories:** Array, Greedy, Sorting

---

## Problem Statement

Given an integer array `nums`, reorder it such that
`nums[0] <= nums[1] >= nums[2] <= nums[3]...`.

You may assume the input array always has a valid answer.

**Example 1:**

```
Input: nums = [3,5,2,1,6,4]
Output: [3,5,1,6,2,4]
Explanation: [1,6,2,5,3,4] is also accepted.
```

**Example 2:**

```
Input: nums = [6,6,5,6,3,8]
Output: [6,6,5,6,3,8]
```

**Constraints:**

- `1 <= nums.length <= 5 * 10^4`
- `0 <= nums[i] <= 10^4`
- It is guaranteed that there will be an answer for the given input `nums`.

**Follow up:** Could you solve the problem in `O(n)` time complexity?

> Note: any output satisfying `nums[0] <= nums[1] >= nums[2] <= ...` is accepted;
> the examples show one valid arrangement each.

---

## Company Frequency

| Company   | Frequency        | Last Reported |
|-----------|------------------|---------------|
| Google    | ★★★☆☆ Medium     | 2023          |
| Facebook  | ★★★☆☆ Medium     | 2022          |
| Amazon    | ★★☆☆☆ Low        | 2022          |
| Microsoft | ★★☆☆☆ Low        | 2022          |

> ⚠️ Frequency data is crowd-sourced from LeetCode Discuss, Glassdoor, and
> community interview reports. Treat as a signal, not a guarantee.

---

## DSA Concepts Used

- **Greedy** — fix each adjacent pair locally with a single swap; a local fix
  never breaks the relation to its left → see [`/dsa/greedy.md`](/dsa/greedy.md)
- **Sorting** — one straightforward solution sorts then swaps pairs → see [`/dsa/sorting.md`](/dsa/sorting.md)

---

## Approaches Overview

| # | Approach | Time | Space | When to use |
|---|----------|------|-------|-------------|
| 1 | Sort then Swap Pairs | O(n log n) | O(1) | Simple to reason about; not optimal |
| 2 | One-Pass Greedy Swap (Optimal) | O(n) | O(1) | Meets the O(n) follow-up |

---

## Approach 1 — Sort then Swap Pairs

### Intuition
After sorting ascending, the array is monotone. Swapping the adjacent pairs
`(1,2), (3,4), …` pushes the larger element of each pair into the odd "peak"
positions, producing the `<= >= <= >=` pattern.

### Algorithm
1. Sort `nums` ascending.
2. For `i = 1; i < n-1; i += 2`: swap `nums[i]` and `nums[i+1]`.

### Complexity
- **Time:** O(n log n) — dominated by the sort.
- **Space:** O(1) extra — in-place swaps (Go's `sort.Ints` sorts in place).

### Code
```go
func sortAndSwap(nums []int) {
	sort.Ints(nums) // ascending order first
	for i := 1; i < len(nums)-1; i += 2 {
		// swap each peak position with its successor to build the wiggle
		nums[i], nums[i+1] = nums[i+1], nums[i]
	}
}
```

### Dry Run
`nums = [3,5,2,1,6,4]`:

| Step | Array |
|------|-------|
| after sort | `[1,2,3,4,5,6]` |
| swap i=1 (2,3) | `[1,3,2,4,5,6]` |
| swap i=3 (4,5) | `[1,3,2,5,4,6]` |
| i=5 → `i < n-1`? 5<5 false, stop | `[1,3,2,5,4,6]` |

Result `[1,3,2,5,4,6]`: `1<=3>=2<=5>=4<=6` ✓ (a valid wiggle).

---

## Approach 2 — One-Pass Greedy Swap (Optimal)

### Intuition
Walk left to right. At an **even** index `i` we want a valley
(`nums[i] <= nums[i+1]`); at an **odd** index `i` we want a peak
(`nums[i] >= nums[i+1]`). Whenever the current pair violates the desired
relation, swap the two elements. Crucially, this local swap can only make the
element at `i` more extreme in the right direction — it never breaks the
already-satisfied relation with `i-1` — so a single pass fixes the whole array.

### Algorithm
1. For `i = 0..n-2`:
   - if `i` is even and `nums[i] > nums[i+1]`: swap.
   - if `i` is odd and `nums[i] < nums[i+1]`: swap.

### Complexity
- **Time:** O(n) — a single pass.
- **Space:** O(1) — in place.

### Code
```go
func greedy(nums []int) {
	for i := 0; i < len(nums)-1; i++ {
		if i%2 == 0 {
			// even index should be a "valley": nums[i] <= nums[i+1]
			if nums[i] > nums[i+1] {
				nums[i], nums[i+1] = nums[i+1], nums[i]
			}
		} else {
			// odd index should be a "peak": nums[i] >= nums[i+1]
			if nums[i] < nums[i+1] {
				nums[i], nums[i+1] = nums[i+1], nums[i]
			}
		}
	}
}
```

### Dry Run
`nums = [3,5,2,1,6,4]`:

| i | parity | want | pair (nums[i],nums[i+1]) | swap? | array after |
|---|--------|------|--------------------------|-------|-------------|
| 0 | even | `<=` | (3,5) → 3<=5 | no | `[3,5,2,1,6,4]` |
| 1 | odd  | `>=` | (5,2) → 5>=2 | no | `[3,5,2,1,6,4]` |
| 2 | even | `<=` | (2,1) → 2>1 | yes | `[3,5,1,2,6,4]` |
| 3 | odd  | `>=` | (2,6) → 2<6 | yes | `[3,5,1,6,2,4]` |
| 4 | even | `<=` | (2,4) → 2<=4 | no | `[3,5,1,6,2,4]` |

Result `[3,5,1,6,2,4]`: `3<=5>=1<=6>=2<=4` ✓.

---

## Key Takeaways

- The greedy insight — "fix each pair locally, a valid swap never breaks the
  left neighbor" — turns an O(n log n) sort into an O(n) one-pass solution.
- Parity of the index encodes the desired relation: even = valley, odd = peak.
- Unlike Wiggle Sort II (#324), duplicates at the boundary are fine here because
  the relations are non-strict (`<=`, `>=`), so a simple in-place swap works.

---

## Related Problems

- LeetCode #324 — Wiggle Sort II (strict `<`/`>`, needs median + interleaving)
- LeetCode #75 — Sort Colors (in-place one-pass rearrangement)
- LeetCode #215 — Kth Largest Element (median/selection, used by Wiggle Sort II)
