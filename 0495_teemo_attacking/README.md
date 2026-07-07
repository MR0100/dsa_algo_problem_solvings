# 0495 — Teemo Attacking

> LeetCode #495 · Difficulty: Easy
> **Categories:** Array, Simulation, Intervals

---

## Problem Statement

Our hero Teemo is attacking an enemy Ashe with poison attacks! When Teemo attacks Ashe, Ashe gets poisoned for a exactly `duration` seconds. More formally, an attack at second `t` will mean Ashe is poisoned during the **inclusive** time interval `[t, t + duration - 1]`. If Teemo attacks again **before** the poison effect ends, the timer for it is **reset**, and the poison effect will end `duration` seconds after the new attack.

You are given a **non-decreasing** integer array `timeSeries`, where `timeSeries[i]` denotes that Teemo attacks Ashe at second `timeSeries[i]`, and an integer `duration`.

Return *the **total** number of seconds that Ashe is poisoned*.

**Example 1:**

```
Input: timeSeries = [1,4], duration = 2
Output: 4
Explanation: Teemo's attacks on Ashe go as follows:
- At second 1, Teemo attacks, and Ashe is poisoned for seconds 1 and 2.
- At second 4, Teemo attacks, and Ashe is poisoned for seconds 4 and 5.
Ashe is poisoned for seconds 1, 2, 4, and 5, which is 4 seconds in total.
```

**Example 2:**

```
Input: timeSeries = [1,2], duration = 2
Output: 3
Explanation: Teemo's attacks on Ashe go as follows:
- At second 1, Teemo attacks, and Ashe is poisoned for seconds 1 and 2.
- At second 2 however, Teemo attacks again and resets the poison timer. Ashe is poisoned for seconds 2 and 3.
Ashe is poisoned for seconds 1, 2, and 3, which is 3 seconds in total.
```

**Constraints:**

- `1 <= timeSeries.length <= 10^4`
- `0 <= timeSeries[i], duration <= 10^7`
- `timeSeries` is sorted in **non-decreasing** order.

---

## Company Frequency

| Company    | Frequency        | Last Reported |
|------------|------------------|---------------|
| Amazon     | ★★★☆☆ Medium     | 2023          |
| Google     | ★★☆☆☆ Low        | 2023          |
| Microsoft  | ★★☆☆☆ Low        | 2022          |
| Bloomberg  | ★★☆☆☆ Low        | 2022          |
| Apple      | ★☆☆☆☆ Rare       | 2021          |

> ⚠️ Frequency data is crowd-sourced from LeetCode Discuss, Glassdoor, and
> community interview reports. Treat as a signal, not a guarantee.

---

## DSA Concepts Used

- **Interval union / merge** — each attack is the interval `[t, t+duration-1]`; the answer is the measure of their union, and because `timeSeries` is sorted, adjacent intervals either overlap (extend) or are disjoint (start a new block) → see [`/dsa/intervals.md`](/dsa/intervals.md)
- **Array single-pass scan** — the optimal solution is one linear sweep over adjacent gaps, `Σ min(gap, duration)` plus a final full duration → see [`/dsa/arrays.md`](/dsa/arrays.md)

---

## Approaches Overview

| # | Approach | Time | Space | When to use |
|---|----------|------|-------|-------------|
| 1 | Timeline Simulation | O(n · duration) | O(span) | Intuition builder; TLE/MLE when duration ~ 10⁷ |
| 2 | Single-Pass Gap Sum (Optimal) | O(n) | O(1) | The intended answer; shortest and fastest |
| 3 | Interval Union Merge | O(n) | O(1) | Same idea framed as generic interval merging |

---

## Approach 1 — Timeline Simulation

### Intuition

Take the problem literally: lay out an explicit second-by-second timeline covering the whole poisoned window (from the first attack to the last attack + `duration`). For each attack, paint its `duration` seconds as poisoned. Overlaps take care of themselves because painting an already-poisoned second changes nothing. The answer is the count of painted seconds. Correct and obvious, but the timeline can be enormous (`duration` up to `10⁷`).

### Algorithm

1. If `timeSeries` is empty or `duration == 0`, return `0`.
2. Offset every attack by `base = timeSeries[0]` so the array starts at index `0`; size it to `timeSeries[last] − base + duration`.
3. For each attack, set `poisoned[start .. start+duration-1] = true`.
4. Count and return the number of `true` entries.

### Complexity

- **Time:** O(n · duration) — each of `n` attacks paints up to `duration` seconds.
- **Space:** O(span) — one boolean per second of the modeled window.

### Code

```go
func simulate(timeSeries []int, duration int) int {
	if len(timeSeries) == 0 || duration == 0 {
		return 0
	}
	base := timeSeries[0]                                   // shift so the array is 0-indexed
	span := timeSeries[len(timeSeries)-1] - base + duration // total seconds to model
	poisoned := make([]bool, span)                          // poisoned[k] == true if second base+k is poisoned
	for _, t := range timeSeries {
		start := t - base // this attack's start, relative to base
		for k := start; k < start+duration; k++ {
			poisoned[k] = true // paint each second of [t, t+duration-1]
		}
	}
	count := 0
	for _, p := range poisoned {
		if p {
			count++ // tally the painted (poisoned) seconds
		}
	}
	return count
}
```

### Dry Run

Example 1: `timeSeries = [1,4]`, `duration = 2`. `base = 1`, `span = 4 − 1 + 2 = 5` (models seconds 1..5 → indices 0..4).

| Attack t | start = t−base | paints indices | poisoned array (idx 0..4) |
|----------|----------------|----------------|---------------------------|
| 1 | 0 | 0,1 | `[T,T,F,F,F]` |
| 4 | 3 | 3,4 | `[T,T,F,T,T]` |

Count of `true` = `4` ✔ (seconds 1,2,4,5).

---

## Approach 2 — Single-Pass Gap Sum (Optimal)

### Intuition

Look at two consecutive attacks `i` and `i+1`. Attack `i` *would* poison a full `duration` seconds, but it may be cut short if the next attack arrives first. The **gap** `timeSeries[i+1] − timeSeries[i]` is exactly how many fresh seconds attack `i` contributes before its timer is reset:

- If `gap >= duration`, attack `i` finished uninterrupted → it added `duration` fresh seconds.
- If `gap < duration`, attack `i+1` reset the timer early → attack `i` only added `gap` fresh seconds (the overlap is credited to attack `i+1`).

So each adjacent pair contributes `min(gap, duration)`, and the **last** attack always contributes a full `duration` (nothing comes after to reset it). One pass, constant space.

### Algorithm

1. If empty or `duration == 0`, return `0`.
2. `total = 0`; for `i` from `0` to `n-2`: `total += min(timeSeries[i+1] − timeSeries[i], duration)`.
3. `total += duration` for the final, uninterrupted attack.
4. Return `total`.

### Complexity

- **Time:** O(n) — a single sweep over adjacent pairs.
- **Space:** O(1) — one accumulator.

### Code

```go
func gapSum(timeSeries []int, duration int) int {
	if len(timeSeries) == 0 || duration == 0 {
		return 0
	}
	total := 0
	for i := 0; i+1 < len(timeSeries); i++ {
		gap := timeSeries[i+1] - timeSeries[i] // seconds until the next reset
		if gap < duration {
			total += gap // interrupted early: only `gap` fresh seconds counted
		} else {
			total += duration // uninterrupted: full poison duration
		}
	}
	total += duration // the final attack always runs its full duration
	return total
}
```

### Dry Run

Example 2: `timeSeries = [1,2]`, `duration = 2`.

| i | gap = ts[i+1] − ts[i] | min(gap, duration) | total |
|---|-----------------------|--------------------|-------|
| 0 | 2 − 1 = 1 | min(1, 2) = 1 | 1 |
| — | (loop ends) | add final duration = 2 | 1 + 2 = **3** |

Result `3` ✔ — the first attack contributed only `1` fresh second (interrupted at second 2), the second contributed its full `2`.

---

## Approach 3 — Interval Union Merge

### Intuition

Each attack is the interval `[t, t + duration − 1]`; the poisoned total is the **measure of the union** of these intervals. Since `timeSeries` is sorted, sweep left→right keeping the current merged block `[start, end]`. If the next attack starts at or before `end` it overlaps — extend `end`. Otherwise the block is complete: add its length `end − start + 1`, and open a new block at this attack. This is the textbook "merge intervals, then sum lengths" routine, and it would still work on unsorted input after a sort.

### Algorithm

1. If empty or `duration == 0`, return `0`.
2. `start = timeSeries[0]`, `end = start + duration − 1`, `total = 0`.
3. For each later attack `t` with interval `[t, t+duration−1]`:
   - if `t <= end` (overlap/touch), `end = max(end, t+duration−1)`;
   - else close the block: `total += end − start + 1`; reset `start, end` to this attack's interval.
4. Add the final block's length; return `total`.

### Complexity

- **Time:** O(n) with the given sorted input (O(n log n) if a sort were required).
- **Space:** O(1).

### Code

```go
func mergeIntervals(timeSeries []int, duration int) int {
	if len(timeSeries) == 0 || duration == 0 {
		return 0
	}
	total := 0
	start := timeSeries[0]      // current merged block's start second
	end := start + duration - 1 // current merged block's end second (inclusive)
	for i := 1; i < len(timeSeries); i++ {
		t := timeSeries[i]
		newEnd := t + duration - 1 // interval contributed by this attack
		if t <= end {              // overlaps or touches the current block
			if newEnd > end {
				end = newEnd // extend the block to cover this attack
			}
		} else {
			total += end - start + 1 // close the finished block, tally its length
			start, end = t, newEnd   // begin a fresh block at this attack
		}
	}
	total += end - start + 1 // add the last open block
	return total
}
```

### Dry Run

Example 1: `timeSeries = [1,4]`, `duration = 2`. Intervals: attack 1 → `[1,2]`, attack 4 → `[4,5]`.

| i | t | newEnd = t+d−1 | current [start, end] | t <= end? | action | total |
|---|---|----------------|----------------------|-----------|--------|-------|
| init | — | — | [1, 2] | — | open first block | 0 |
| 1 | 4 | 5 | [1, 2] | 4 ≤ 2? no | close [1,2]: +2; open [4,5] | 2 |
| end | — | — | [4, 5] | — | add final: +(5−4+1)=2 | 2 + 2 = **4** |

Result `4` ✔ — two disjoint blocks of length 2 each.

---

## Key Takeaways

- **"Total time under repeated, possibly-overlapping effects" = union of intervals.** When the events are sorted, you never need to build the timeline — compare each adjacent pair.
- The optimal one-liner insight: **each attack adds `min(gap_to_next, duration)`**, and the last adds a full `duration`. This "min of gap and effect length" is the reusable trick.
- **Inclusive intervals shift the arithmetic by one:** `[t, t+duration-1]` has length `duration`, and a merged block `[start,end]` has length `end−start+1`. Off-by-one here is the classic bug.
- Guard the degenerate `duration == 0` (poison lasts zero seconds → total `0`), which the constraints allow (`0 <= duration`).
- Avoid the simulation for large `duration` (up to `10⁷`): it is O(n·duration) time and O(span) memory — both blow up while the O(n)/O(1) sweep does not.

---

## Related Problems

- LeetCode #56 — Merge Intervals (the general merge routine)
- LeetCode #57 — Insert Interval (interval overlap handling)
- LeetCode #253 — Meeting Rooms II (interval overlap counting)
- LeetCode #452 — Minimum Number of Arrows to Burst Balloons (greedy over sorted intervals)
- LeetCode #435 — Non-overlapping Intervals (interval sweeping)
