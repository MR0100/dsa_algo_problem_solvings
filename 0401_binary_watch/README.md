# 0401 — Binary Watch

> LeetCode #401 · Difficulty: Easy
> **Categories:** Backtracking, Bit Manipulation

---

## Problem Statement

A binary watch has 4 LEDs on the top which represent the hours (0-11), and 6 LEDs on the bottom which represent the minutes (0-59). Each LED represents a zero or one, with the least significant bit on the right.

- For example, the below binary watch reads `"4:51"`.

Given an integer `turnedOn` which represents the number of LEDs that are currently on (ignoring the PM), return *all possible times the watch could represent*. You may return the answer in **any order**.

The hour must not contain a leading zero.

- For example, `"01:00"` is not valid. It should be `"1:00"`.

The minute must consist of two digits and may contain a leading zero.

- For example, `"10:2"` is not valid. It should be `"10:02"`.

**Example 1:**

```
Input: turnedOn = 1
Output: ["0:01","0:02","0:04","0:08","0:16","0:32","1:00","2:00","4:00","8:00"]
```

**Example 2:**

```
Input: turnedOn = 9
Output: []
```

**Constraints:**

- `0 <= turnedOn <= 10`

---

## Company Frequency

| Company    | Frequency        | Last Reported |
|------------|------------------|---------------|
| Google     | ★★★☆☆ Medium     | 2023          |
| Amazon     | ★★★☆☆ Medium     | 2023          |
| Microsoft  | ★★☆☆☆ Low        | 2022          |
| Apple      | ★★☆☆☆ Low        | 2022          |
| Bloomberg  | ★★☆☆☆ Low        | 2021          |

> ⚠️ Frequency data is crowd-sourced from LeetCode Discuss, Glassdoor, and
> community interview reports. Treat as a signal, not a guarantee.

---

## DSA Concepts Used

- **Bit Manipulation** — a lit-LED count is exactly `popcount(hour) + popcount(minute)`; the whole problem is "which numbers have a given number of set bits" → see [`/dsa/bit_manipulation.md`](/dsa/bit_manipulation.md)
- **Backtracking** — the optimal framing chooses which of the 10 LEDs to switch on, pruning any partial choice whose hour exceeds 11 or minute exceeds 59 → see [`/dsa/backtracking.md`](/dsa/backtracking.md)

---

## Approaches Overview

| # | Approach | Time | Space | When to use |
|---|----------|------|-------|-------------|
| 1 | Brute Force (scan 12×60) | O(720) = O(1) | O(1) | Simplest correct answer; the tiny fixed search space makes it optimal in practice |
| 2 | Split the Budget | O(720) = O(1) | O(1) | Generalises to watches with different LED counts; frames it as distributing set bits |
| 3 | Backtracking over 10 LEDs | O(C(10,k)) pruned | O(k) | The "intended" combinatorial answer; demonstrates subset generation with pruning |

---

## Approach 1 — Brute Force

### Intuition

A binary watch just displays a number in binary, so the number of lit LEDs for a time is `popcount(hour) + popcount(minute)`. There are only `12 × 60 = 720` valid times, so rather than deciding which LEDs to light, we can test every time and keep those whose total bit count equals `turnedOn`. The tiny, fixed search space means this is not merely acceptable — it is effectively optimal.

### Algorithm

1. Loop `h` from `0` to `11` (the 4 hour LEDs cover this range).
2. Loop `m` from `0` to `59` (the 6 minute LEDs cover this range).
3. If `popcount(h) + popcount(m) == turnedOn`, format as `"h:mm"` (no leading zero on the hour, two digits on the minute) and collect it.
4. Return all collected strings.

### Complexity

- **Time:** O(720) = O(1) — a fixed 720 iterations no matter what `turnedOn` is; each does a constant-time popcount.
- **Space:** O(1) beyond the output list, which holds at most a few dozen strings.

### Code

```go
func bruteForce(turnedOn int) []string {
	var result []string
	for h := 0; h < 12; h++ { // 4 hour LEDs cover 0..11
		for m := 0; m < 60; m++ { // 6 minute LEDs cover 0..59
			// bits.OnesCount counts the set bits (lit LEDs) of the value.
			if bits.OnesCount(uint(h))+bits.OnesCount(uint(m)) == turnedOn {
				// %d gives the hour with no leading zero; %02d pads the
				// minute to exactly two digits (e.g. 2 -> "02").
				result = append(result, fmt.Sprintf("%d:%02d", h, m))
			}
		}
	}
	return result
}
```

### Dry Run

Example 1: `turnedOn = 1`. We need `popcount(h) + popcount(m) == 1`, i.e. exactly one lit LED total.

| h | popcount(h) | which m qualify (popcount(m) == 1−popcount(h)) | emitted |
|---|-------------|-----------------------------------------------|---------|
| 0 | 0 | m with popcount 1: 1,2,4,8,16,32 | `0:01 0:02 0:04 0:08 0:16 0:32` |
| 1 | 1 | m with popcount 0: 0 | `1:00` |
| 2 | 1 | m with popcount 0: 0 | `2:00` |
| 3 | 2 | need popcount(m) = −1 → none | — |
| 4 | 1 | m = 0 | `4:00` |
| 5..7 | 2 | none | — |
| 8 | 1 | m = 0 | `8:00` |
| 9..11 | 2 | none | — |

Result: `[0:01 0:02 0:04 0:08 0:16 0:32 1:00 2:00 4:00 8:00]` ✔

---

## Approach 2 — Split the Budget

### Intuition

Instead of scanning raw times, decide **how many** of the `turnedOn` lit LEDs belong to the hour. If the hour uses `hb` LEDs then the minute must use `turnedOn − hb`. For each split, gather the hours with exactly `hb` set bits and the minutes with exactly `turnedOn − hb` set bits, and take their cross product. This makes the "distribute set bits between two fields" structure explicit and adapts cleanly if the LED counts ever change.

### Algorithm

1. For `hb` from `0` to `min(turnedOn, 4)` (at most 4 hour LEDs exist):
   1. Let `mb = turnedOn − hb`; skip if `mb < 0` or `mb > 6` (only 6 minute LEDs).
   2. Collect every hour `h` in `0..11` with `popcount(h) == hb`.
   3. Collect every minute `m` in `0..59` with `popcount(m) == mb`.
   4. Emit `"h:mm"` for every (hour, minute) pair.
2. Return the collected list.

### Complexity

- **Time:** O(720) = O(1) — the inner scans over `0..11` and `0..59` are fixed; summed over the ≤5 splits it is still constant.
- **Space:** O(1) beyond the output (small temporary lists of hours/minutes).

### Code

```go
func splitBudget(turnedOn int) []string {
	var result []string
	// hb = number of lit LEDs assigned to the hour. At most 4 hour LEDs exist,
	// and it obviously cannot exceed the total budget.
	for hb := 0; hb <= turnedOn && hb <= 4; hb++ {
		mb := turnedOn - hb // the rest of the budget goes to the minute
		if mb < 0 || mb > 6 {
			continue // impossible: only 6 minute LEDs are available
		}
		// Collect all hours whose popcount matches the hour budget.
		var hours []int
		for h := 0; h < 12; h++ {
			if bits.OnesCount(uint(h)) == hb {
				hours = append(hours, h)
			}
		}
		// Collect all minutes whose popcount matches the minute budget.
		var mins []int
		for m := 0; m < 60; m++ {
			if bits.OnesCount(uint(m)) == mb {
				mins = append(mins, m)
			}
		}
		// Cross product: every valid hour with every valid minute.
		for _, h := range hours {
			for _, m := range mins {
				result = append(result, fmt.Sprintf("%d:%02d", h, m))
			}
		}
	}
	return result
}
```

### Dry Run

Example 1: `turnedOn = 1`.

| hb | mb = 1−hb | hours (popcount hb) | minutes (popcount mb) | emitted pairs |
|----|-----------|---------------------|-----------------------|---------------|
| 0 | 1 | {0} | {1,2,4,8,16,32} | `0:01 0:02 0:04 0:08 0:16 0:32` |
| 1 | 0 | {1,2,4,8} | {0} | `1:00 2:00 4:00 8:00` |

`hb` cannot reach 2 (that would need `mb = −1`). Result: `[0:01 0:02 0:04 0:08 0:16 0:32 1:00 2:00 4:00 8:00]` ✔

---

## Approach 3 — Backtracking

### Intuition

Model the watch as 10 LEDs with fixed weights: hour LEDs contribute `{1,2,4,8}` and minute LEDs contribute `{1,2,4,8,16,32}`. Choose exactly `turnedOn` LEDs to switch on; each choice adds its weight to either the running hour or minute. When we have placed `turnedOn` LEDs we have a concrete `(hour, minute)` to record. Crucially, the moment a partial choice makes `hour > 11` or `minute > 59` we abandon that branch — this pruning is what turns a blind subset enumeration into a targeted one.

### Algorithm

1. Fix LED weights: indices `0..3` → `1,2,4,8` (hours), indices `4..9` → `1,2,4,8,16,32` (minutes).
2. `dfs(index, remaining, hour, minute)`:
   1. If `hour > 11` or `minute > 59`: prune (return).
   2. If `remaining == 0`: record `"hour:minute"` and return.
   3. For each LED `i` from `index` to `9`: light it (add its weight to hour if `i < 4`, else to minute) and recurse with `remaining − 1` and start `i + 1`.
3. Sort the results for a deterministic output.

### Complexity

- **Time:** O(C(10, turnedOn)) subsets in the worst case, but pruning discards most branches; since there are only 10 LEDs this is effectively O(1).
- **Space:** O(turnedOn) recursion depth plus the output list.

### Code

```go
func backtracking(turnedOn int) []string {
	// Weight of each of the 10 LEDs. Indices 0..3 are hour bits, 4..9 minute.
	weights := []int{1, 2, 4, 8, 1, 2, 4, 8, 16, 32}

	var result []string
	var dfs func(index, remaining, hour, minute int)
	dfs = func(index, remaining, hour, minute int) {
		if hour > 11 || minute > 59 {
			return // invalid clock face — abandon this branch immediately
		}
		if remaining == 0 {
			// Exactly turnedOn LEDs are lit; emit the formatted time.
			result = append(result, fmt.Sprintf("%d:%02d", hour, minute))
			return
		}
		// Choose the next LED to light among the remaining LEDs (index..9),
		// which guarantees each subset of LEDs is generated exactly once.
		for i := index; i < 10; i++ {
			if i < 4 {
				// Hour LED: add its weight to the hour half.
				dfs(i+1, remaining-1, hour+weights[i], minute)
			} else {
				// Minute LED: add its weight to the minute half.
				dfs(i+1, remaining-1, hour, minute+weights[i])
			}
		}
	}
	dfs(0, turnedOn, 0, 0)

	// The recursion visits LEDs in a fixed order, so results are already grouped,
	// but sorting makes the output stable and easy to compare in tests.
	sort.Strings(result)
	return result
}
```

### Dry Run

Example 1: `turnedOn = 1`. We start `dfs(0, 1, 0, 0)` and pick exactly one LED. Because `remaining` starts at 1, every recursive branch lights a single LED then hits the `remaining == 0` base case.

| LED i picked | row | weight | resulting (hour, minute) | valid? | recorded |
|--------------|-----|--------|--------------------------|--------|----------|
| 0 | hour | 1 | (1, 0) | yes | `1:00` |
| 1 | hour | 2 | (2, 0) | yes | `2:00` |
| 2 | hour | 4 | (4, 0) | yes | `4:00` |
| 3 | hour | 8 | (8, 0) | yes | `8:00` |
| 4 | minute | 1 | (0, 1) | yes | `0:01` |
| 5 | minute | 2 | (0, 2) | yes | `0:02` |
| 6 | minute | 4 | (0, 4) | yes | `0:04` |
| 7 | minute | 8 | (0, 8) | yes | `0:08` |
| 8 | minute | 16 | (0, 16) | yes | `0:16` |
| 9 | minute | 32 | (0, 32) | yes | `0:32` |

After `sort.Strings`: `[0:01 0:02 0:04 0:08 0:16 0:32 1:00 2:00 4:00 8:00]` ✔

---

## Key Takeaways

- **Tiny fixed search spaces beat clever algorithms.** With only 720 times, a direct scan is simplest *and* effectively optimal — reach for combinatorics only when the space actually blows up.
- **"Number of lit LEDs" is popcount.** Recognising that a display of binary numbers has bit-count semantics collapses the problem to "how many integers in a range have exactly k set bits".
- **Backtracking = choose + prune.** Generate each subset once by only ever picking indices `≥ start`, and cut invalid branches (`hour > 11`, `minute > 59`) as early as possible.
- **Formatting matters in interviews.** `%d:%02d` cleanly enforces "hour without leading zero, minute always two digits" — a common source of wrong-answer verdicts here.

---

## Related Problems

- LeetCode #191 — Number of 1 Bits (popcount, the core primitive)
- LeetCode #78 — Subsets (subset generation by backtracking)
- LeetCode #77 — Combinations (choose k of n with pruning)
- LeetCode #17 — Letter Combinations of a Phone Number (cross-product enumeration)
- LeetCode #338 — Counting Bits (bit counts across a range)
