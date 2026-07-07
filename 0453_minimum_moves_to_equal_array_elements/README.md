# 0453 — Minimum Moves to Equal Array Elements

> LeetCode #453 · Difficulty: Medium
> **Categories:** Array, Math

---

## Problem Statement

Given an integer array `nums` of size `n`, return *the minimum number of moves required to make all array elements equal*.

In one move, you can increment `n - 1` elements of the array by `1`.

**Example 1:**

```
Input: nums = [1,2,3]
Output: 3
Explanation: Only three moves are needed (remember each move increments two elements):
[1,2,3]  =>  [2,3,3]  =>  [3,4,3]  =>  [4,4,4]
```

**Example 2:**

```
Input: nums = [1,1,1]
Output: 0
```

**Constraints:**

- `n == nums.length`
- `1 <= nums.length <= 10^5`
- `-10^9 <= nums[i] <= 10^9`
- The answer is guaranteed to fit in a **32-bit** integer.

---

## Company Frequency

| Company    | Frequency        | Last Reported |
|------------|------------------|---------------|
| Amazon     | ★★★☆☆ Medium     | 2023          |
| Google     | ★★☆☆☆ Low        | 2023          |
| Microsoft  | ★★☆☆☆ Low        | 2022          |
| Adobe      | ★★☆☆☆ Low        | 2022          |
| Bloomberg  | ★☆☆☆☆ Rare       | 2021          |

> ⚠️ Frequency data is crowd-sourced from LeetCode Discuss, Glassdoor, and
> community interview reports. Treat as a signal, not a guarantee.

---

## DSA Concepts Used

- **Math / Number Theory (invariant reframing)** — the trick is realising that "add 1 to `n-1` elements" is *relatively* identical to "subtract 1 from 1 element": both shrink exactly one pairwise gap by 1. That reframing collapses the whole problem to `sum - n * min`, computed in one pass → see [`/dsa/math_number_theory.md`](/dsa/math_number_theory.md)
- **Array (single-pass aggregation)** — the optimal solution is a textbook running-sum-and-minimum scan over the array → see [`/dsa/arrays.md`](/dsa/arrays.md)

---

## Approaches Overview

| # | Approach | Time | Space | When to use |
|---|----------|------|-------|-------------|
| 1 | Brute Force Simulation | O(moves · n), moves up to ~10⁹ | O(1) | Only to build intuition; TLE on large gaps |
| 2 | Math — Increment ≡ Decrement (Optimal) | O(n) | O(1) | The answer; `sum - n*min` in one pass |
| 3 | Sort + Sum of Differences From Min | O(n log n) | O(n) | Same math after sorting so `min = nums[0]`; a common first instinct |

---

## Approach 1 — Brute Force Simulation

### Intuition

Follow the statement literally. Incrementing `n-1` elements by 1 is the same as *not* incrementing exactly one element — and the sensible one to leave out is the current maximum (raising it too would only widen the gap). So each move keeps the max fixed while everyone else rises by 1. Repeat until `min == max`. This is faithful but its running time is proportional to the *answer*, which can reach ~10⁹, so it is a teaching baseline, not a submission.

### Algorithm

1. Set `moves = 0`.
2. Repeat:
   - Scan for the minimum, the maximum, and an index of the maximum.
   - If `min == max`, stop and return `moves`.
   - Otherwise increment every element except that one maximum, and `moves++`.

### Complexity

- **Time:** O(moves · n) — each move costs an O(n) scan, and there can be ~10⁹ moves. TLE on large value gaps.
- **Space:** O(1) — increments happen in place.

### Code

```go
func bruteForce(nums []int) int {
	moves := 0
	for {
		// Find current min and the index of a maximum.
		minVal, maxVal, maxIdx := nums[0], nums[0], 0
		for i, v := range nums {
			if v < minVal {
				minVal = v
			}
			if v > maxVal {
				maxVal = v
				maxIdx = i
			}
		}
		if minVal == maxVal {
			return moves // all equal → done
		}
		// One move: increment every element except one chosen maximum.
		for i := range nums {
			if i != maxIdx {
				nums[i]++ // n-1 elements go up by 1
			}
		}
		moves++
	}
}
```

### Dry Run

Example 1: `nums = [1,2,3]`.

| Move | nums before | min | max (idx) | increment all but idx | nums after | moves |
|------|-------------|-----|-----------|-----------------------|------------|-------|
| 1 | `[1,2,3]` | 1 | 3 (i=2) | +1 to indices 0,1 | `[2,3,3]` | 1 |
| 2 | `[2,3,3]` | 2 | 3 (i=1) | +1 to indices 0,2 | `[3,4,3]` | 2 |
| 3 | `[3,4,3]` | 3 | 4 (i=1) | +1 to indices 0,2 | `[4,4,4]` | 3 |
| — | `[4,4,4]` | 4 | 4 | min == max → stop | — | 3 |

Result: `3` ✔.

---

## Approach 2 — Math — Increment ≡ Decrement (Optimal)

### Intuition

Only the **differences** between elements matter for "make them all equal", and the actual target value is irrelevant. Adding 1 to `n-1` elements raises the whole array's baseline but shrinks the gap between the left-out element and the rest by 1 — exactly what subtracting 1 from that one element would do to the gaps. So flip the picture: instead of pushing everyone *up* to the maximum, pull everyone *down* to the minimum, one unit per move. Element `nums[i]` needs `nums[i] - min` such moves, and summing gives

```
moves = Σ (nums[i] - min) = sum(nums) - n * min(nums).
```

### Algorithm

1. In one pass, accumulate `sum` and track `min`.
2. Return `sum - n * min`.

### Complexity

- **Time:** O(n) — a single pass.
- **Space:** O(1) — two scalars.

### Code

```go
func mathDecrement(nums []int) int {
	sum := 0
	minVal := nums[0]
	for _, v := range nums {
		sum += v // running total of all elements
		if v < minVal {
			minVal = v // smallest element = the level everyone descends to
		}
	}
	// Each element must drop to minVal; the drops summed = total moves.
	return sum - len(nums)*minVal
}
```

### Dry Run

Example 1: `nums = [1,2,3]`.

| Step | v | sum after | minVal after |
|------|---|-----------|--------------|
| 1 | 1 | 1 | 1 |
| 2 | 2 | 3 | 1 |
| 3 | 3 | 6 | 1 |

Final: `sum = 6`, `n = 3`, `min = 1` → `6 - 3*1 = 3` ✔.

---

## Approach 3 — Sort + Sum of Differences From Min

### Intuition

The identical `Σ (nums[i] - min)` sum, but derived after sorting so the minimum is trivially `nums[0]`. Sorting adds a log factor and buys nothing over Approach 2, yet it is a natural first instinct: "line them up, then measure how far each sits above the smallest." Summing those distances is the answer.

### Algorithm

1. Copy `nums` (to avoid mutating the caller) and sort ascending; now `nums[0]` is the minimum.
2. Accumulate `moves = Σ (nums[i] - nums[0])`.
3. Return `moves`.

### Complexity

- **Time:** O(n log n) — the sort dominates.
- **Space:** O(n) — a copy is sorted so the caller's slice is untouched.

### Code

```go
func sortAndSum(nums []int) int {
	// Work on a copy so we don't disturb the caller's slice ordering.
	cp := make([]int, len(nums))
	copy(cp, nums)
	sort.Ints(cp) // ascending; cp[0] is now the minimum element

	moves := 0
	for _, v := range cp {
		moves += v - cp[0] // distance of each element down to the smallest
	}
	return moves
}
```

### Dry Run

Example 1: `nums = [1,2,3]`.

| Step | Action | State |
|------|--------|-------|
| 1 | Copy + sort | `cp = [1,2,3]`, `cp[0] = 1` |
| 2 | add `1 - 1` | `moves = 0` |
| 2 | add `2 - 1` | `moves = 1` |
| 2 | add `3 - 1` | `moves = 3` |

Result: `3` ✔ (identical to `sum - n*min`).

---

## Key Takeaways

- **Reframe operations by their invariant.** "Add 1 to n−1 elements" feels complicated; recognising it as "subtract 1 from 1 element" (same effect on all pairwise gaps) turns an O(answer) simulation into an O(n) formula.
- **When the target value is free, only differences matter.** You are never told what final value to reach, so pick the most convenient one (the minimum) and measure distances to it.
- **`sum - n * min`** is the whole solution; watch that intermediate `sum` can be large (up to ~10¹⁴ for the max input), so use 64-bit — the *answer* fits in 32 bits but the running sum need not on the way there.
- Sorting is a tempting but unnecessary detour here; whenever you only need the minimum and a sum, a single linear scan beats an O(n log n) sort.

---

## Related Problems

- LeetCode #462 — Minimum Moves to Equal Array Elements II (move to the *median*, ±1 on one element)
- LeetCode #2033 — Minimum Operations to Make a Uni-Value Grid (same distance-to-central-value idea with a step `x`)
- LeetCode #2244 — Minimum Rounds to Complete All Tasks (counting-based minimum operations)
