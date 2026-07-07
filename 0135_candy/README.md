# 0135 — Candy

> LeetCode #135 · Difficulty: Hard
> **Categories:** Array, Greedy

---

## Problem Statement

There are `n` children standing in a line. Each child is assigned a rating value given in the integer array `ratings`.

You are giving candies to these children subjected to the following requirements:

- Each child must have at least one candy.
- Children with a higher rating get more candies than their neighbors.

Return *the minimum number of candies you need to have to distribute the candies to the children*.

**Example 1:**
```
Input: ratings = [1,0,2]
Output: 5
Explanation: You can allocate to the first, second and third child with 2, 1, 2 candies respectively.
```

**Example 2:**
```
Input: ratings = [1,2,2]
Output: 4
Explanation: You can allocate to the first, second and third child with 1, 2, 1 candies respectively.
The third child gets 1 candy because it satisfies the above two conditions.
```

**Constraints:**
- `n == ratings.length`
- `1 <= n <= 2 * 10⁴`
- `0 <= ratings[i] <= 2 * 10⁴`

---

## Company Frequency

| Company   | Frequency        | Last Reported |
|-----------|------------------|---------------|
| Amazon    | ★★★★☆ High       | 2024          |
| Google    | ★★★★☆ High       | 2024          |
| Microsoft | ★★★☆☆ Medium     | 2024          |
| Flipkart  | ★★★☆☆ Medium     | 2023          |
| Adobe     | ★★☆☆☆ Low        | 2023          |

> ⚠️ Frequency data is crowd-sourced from LeetCode Discuss, Glassdoor, and
> community interview reports. Treat as a signal, not a guarantee.

---

## DSA Concepts Used

- **Greedy** — give each child the minimum that satisfies its constraints; decomposing a two-sided constraint into two one-sided greedy passes → see [`/dsa/greedy.md`](/dsa/greedy.md)

---

## Approaches Overview

| # | Approach | Time | Space | When to use |
|---|----------|------|-------|-------------|
| 1 | Brute Force (Repeated Sweeps) | O(n²) | O(n) | Baseline; shows why minimal local fixes converge |
| 2 | Two-Pass Arrays | O(n) | O(n) | The standard interview answer; clean correctness proof |
| 3 | Slope Counting (Optimal, O(1) space) | O(n) | O(1) | Follow-up "can you do O(1) space?"; single pass |

---

## Approach 1 — Brute Force (Repeated Sweeps)

### Intuition
The rules are purely local: a child only ever compares with its two neighbors. Start everyone at the legal minimum (1 candy each) and repeatedly repair violations: whenever a child has a strictly higher rating than a neighbor but not strictly more candy, raise that child to `neighbor + 1` — the *smallest* fix. Since every adjustment is forced (any valid allocation must be at least that high) and only ever increases values, the process converges to the unique minimal valid allocation.

### Algorithm
1. Initialize `candies[i] = 1` for all i.
2. Repeat until a full sweep makes no change; in each sweep, for every `i`:
   1. If `i > 0` and `ratings[i] > ratings[i-1]` and `candies[i] <= candies[i-1]`, set `candies[i] = candies[i-1] + 1`.
   2. If `i < n-1` and `ratings[i] > ratings[i+1]` and `candies[i] <= candies[i+1]`, set `candies[i] = candies[i+1] + 1`.
3. Return the sum of `candies`.

### Complexity
- **Time:** O(n²) — each sweep is O(n), and a long monotone slope can require up to ~n sweeps for the `+1`s to propagate to its far end.
- **Space:** O(n) — the `candies` array.

### Code
```go
func bruteForce(ratings []int) int {
	n := len(ratings)
	candies := make([]int, n)
	for i := range candies {
		candies[i] = 1 // everyone gets at least one candy
	}

	for changed := true; changed; {
		changed = false // assume this sweep is clean until proven otherwise
		for i := 0; i < n; i++ {
			// left rule: strictly higher rating than left neighbor
			if i > 0 && ratings[i] > ratings[i-1] && candies[i] <= candies[i-1] {
				candies[i] = candies[i-1] + 1 // minimal fix for the violation
				changed = true
			}
			// right rule: strictly higher rating than right neighbor
			if i < n-1 && ratings[i] > ratings[i+1] && candies[i] <= candies[i+1] {
				candies[i] = candies[i+1] + 1 // minimal fix for the violation
				changed = true
			}
		}
	}

	total := 0
	for _, c := range candies {
		total += c // sum the final minimal allocation
	}
	return total
}
```

### Dry Run — Example 1: `ratings = [1,0,2]`

| Sweep | i | Check | Violation? | Action | `candies` |
|-------|---|-------|------------|--------|-----------|
| — | — | initialize | — | all ones | `[1, 1, 1]` |
| 1 | 0 | right: `1 > 0`, `candies[0](1) <= candies[1](1)` | yes | `candies[0] = 2` | `[2, 1, 1]` |
| 1 | 1 | rating 0 lower than both neighbors | no | — | `[2, 1, 1]` |
| 1 | 2 | left: `2 > 0`, `candies[2](1) <= candies[1](1)` | yes | `candies[2] = 2` | `[2, 1, 2]` |
| 2 | 0..2 | re-check all | none | sweep clean → stop | `[2, 1, 2]` |

Total = 2 + 1 + 2 = `5`. ✅

---

## Approach 2 — Two-Pass Arrays

### Intuition
The neighbor rule couples each child to *both* sides at once — hard to satisfy in one sweep. But it decomposes perfectly into two independent one-sided rules:

- **Left rule:** if `ratings[i] > ratings[i-1]`, child i needs more than child i−1. A single left-to-right pass computes the minimal `left[]` satisfying only this.
- **Right rule:** if `ratings[i] > ratings[i+1]`, child i needs more than child i+1. A right-to-left pass computes the minimal `right[]`.

A child (e.g. on a peak) must satisfy both simultaneously, and `max(left[i], right[i])` does: it meets each side's requirement because it is ≥ each pass's minimum, and it is minimal because any valid allocation must be ≥ both one-sided minimums pointwise.

### Algorithm
1. Initialize `left[i] = right[i] = 1` for all i.
2. Left pass: for `i` from 1 to n−1, if `ratings[i] > ratings[i-1]`, set `left[i] = left[i-1] + 1`.
3. Right pass: for `i` from n−2 down to 0, if `ratings[i] > ratings[i+1]`, set `right[i] = right[i+1] + 1`.
4. Return `Σ max(left[i], right[i])`.

### Complexity
- **Time:** O(n) — three linear passes (left, right, sum).
- **Space:** O(n) — two auxiliary requirement arrays.

### Code
```go
func twoPassArrays(ratings []int) int {
	n := len(ratings)

	left := make([]int, n)  // candies needed looking only at left neighbors
	right := make([]int, n) // candies needed looking only at right neighbors
	for i := range left {
		left[i] = 1 // base: one candy each
		right[i] = 1
	}

	// left-to-right: enforce "ascent from the left means +1 candy"
	for i := 1; i < n; i++ {
		if ratings[i] > ratings[i-1] {
			left[i] = left[i-1] + 1 // strictly more than the left neighbor
		}
	}

	// right-to-left: enforce "ascent from the right means +1 candy"
	for i := n - 2; i >= 0; i-- {
		if ratings[i] > ratings[i+1] {
			right[i] = right[i+1] + 1 // strictly more than the right neighbor
		}
	}

	total := 0
	for i := 0; i < n; i++ {
		// each child must satisfy both one-sided requirements simultaneously
		if left[i] > right[i] {
			total += left[i]
		} else {
			total += right[i]
		}
	}
	return total
}
```

### Dry Run — Example 1: `ratings = [1,0,2]`

Left pass:

| i | Compare | `left` |
|---|---------|--------|
| — | init | `[1, 1, 1]` |
| 1 | `ratings[1](0) > ratings[0](1)`? no | `[1, 1, 1]` |
| 2 | `ratings[2](2) > ratings[1](0)`? yes → `left[1]+1 = 2` | `[1, 1, 2]` |

Right pass:

| i | Compare | `right` |
|---|---------|---------|
| — | init | `[1, 1, 1]` |
| 1 | `ratings[1](0) > ratings[2](2)`? no | `[1, 1, 1]` |
| 0 | `ratings[0](1) > ratings[1](0)`? yes → `right[1]+1 = 2` | `[2, 1, 1]` |

Combine:

| i | `left[i]` | `right[i]` | `max` | running total |
|---|-----------|------------|-------|---------------|
| 0 | 1 | 2 | 2 | 2 |
| 1 | 1 | 1 | 1 | 3 |
| 2 | 2 | 1 | 2 | 5 |

Total = `5`. ✅

---

## Approach 3 — Slope Counting (Optimal, O(1) space)

### Intuition
Look at what the optimal allocation actually is: along a strictly **ascending** run of ratings, candies go 1, 2, 3, …; along a strictly **descending** run they mirror to …, 3, 2, 1; **equal** neighbors carry no constraint, so the count resets to 1. The total therefore depends only on the *lengths* of the monotone runs — no per-child array is needed.

The only interaction is at a **peak** between an ascent of length `up` and a descent of length `down`: the peak child must top both sides, needing `max(up, down) + 1`. The trick: award the peak `up + 1` optimistically during the ascent. While descending, each new step adds `down + 1` candies (1 for the new child, plus 1 to each of the `down − 1` earlier descent children whose mirrored chain shifted up, plus 1 tentatively for the peak). If the recorded `peak` height still exceeds the descent length (`peak >= down`), the peak did not actually need raising — subtract the 1 back. Once `down` exceeds `peak`, the +1 stays each step, lazily growing the peak with the chain.

### Algorithm
1. If `n <= 1`, return `n`. Set `total = 1`, `up = down = peak = 0`.
2. For `i` from 1 to n−1, compare `ratings[i]` to `ratings[i-1]`:
   1. **Greater:** `up++`, `down = 0`, `peak = up`, `total += up + 1`.
   2. **Equal:** `up = down = peak = 0`, `total += 1`.
   3. **Less:** `down++`, `up = 0`, `total += down + 1`; if `peak >= down`, `total--`.
3. Return `total`.

### Complexity
- **Time:** O(n) — one pass, constant work per child.
- **Space:** O(1) — four integer counters; no arrays at all.

### Code
```go
func slopeCounting(ratings []int) int {
	n := len(ratings)
	if n <= 1 {
		return n // 0 children → 0 candies; 1 child → 1 candy
	}

	total := 1 // the first child always gets 1 candy to start
	up, down, peak := 0, 0, 0

	for i := 1; i < n; i++ {
		switch {
		case ratings[i] > ratings[i-1]:
			up++      // ascending run grows
			down = 0  // any previous descent is over
			peak = up // remember the height of the (current) peak
			// child i sits `up` steps above the run's start → needs up+1
			total += up + 1
		case ratings[i] == ratings[i-1]:
			// equal ratings carry no constraint: reset all runs,
			// this child can drop back to a single candy
			up, down, peak = 0, 0, 0
			total += 1
		default: // ratings[i] < ratings[i-1]
			down++ // descending run grows
			up = 0 // any previous ascent is over
			// tentatively this child starts a mirrored 1,2,...,down chain:
			// every earlier child in the descent shifts up by one → +down,
			// plus 1 candy for this child itself
			total += down + 1
			if peak >= down {
				// the peak (given up+1 earlier) is still strictly higher
				// than the descent chain, so it need not grow: take back
				// the one candy we just over-counted for it
				total--
			}
			// if peak < down, the peak must rise with the chain: the +1
			// stays, effectively raising the peak by one this step
		}
	}

	return total
}
```

### Dry Run — Example 1: `ratings = [1,0,2]`

| i | Comparison | Branch | `up` | `down` | `peak` | `total` update | `total` |
|---|------------|--------|------|--------|--------|----------------|---------|
| — | init | — | 0 | 0 | 0 | first child gets 1 | 1 |
| 1 | `0 < 1` | descent | 0 | 1 | 0 | `+ (down+1) = +2`; `peak(0) >= down(1)`? no → keep | 3 |
| 2 | `2 > 0` | ascent | 1 | 0 | 1 | `+ (up+1) = +2` | 5 |

Return `5` — implied allocation `[2, 1, 2]`. ✅

(Example 2, `[1,2,2]`: i=1 ascent → total 1+2 = 3; i=2 equal → reset, +1 → 4. ✅)

---

## Key Takeaways

- **Two-sided neighbor constraints split into two one-sided passes** — solve left-to-right and right-to-left independently, then combine with `max`. This left/right decomposition is a reusable pattern (Trapping Rain Water, Product of Array Except Self).
- **`max` of pointwise minimums stays minimal:** each pass computes the least value satisfying its own rule; any valid allocation dominates both, so their max is optimal.
- **Equal ratings carry no constraint** — the classic trap: `[1,2,2]` lets the third child drop back to 1 candy.
- **Monotone-run bookkeeping kills the arrays:** when the answer along a slope is forced to be 1,2,3,…, only run lengths matter — track `up`, `down`, `peak` and settle the peak's height lazily (`peak >= down` check) for O(1) space.
- Brute force converging via "minimal local repairs only ever increase values" is a useful mental model for why the greedy answers are forced lower bounds.

---

## Related Problems

- LeetCode #134 — Gas Station (adjacent greedy on a circle, same chapter of greedy)
- LeetCode #42 — Trapping Rain Water (left/right pass decomposition with max)
- LeetCode #238 — Product of Array Except Self (prefix and suffix passes combined)
- LeetCode #2193 — Minimum Number of Moves to Make Palindrome (greedy with positional constraints)
- LeetCode #2242 — Maximum Score of a Node Sequence (neighbor-constrained optimization)
