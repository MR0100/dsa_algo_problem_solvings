# 0458 — Poor Pigs

> LeetCode #458 · Difficulty: Hard
> **Categories:** Math, Dynamic Programming, Combinatorics

---

## Problem Statement

There are `buckets` buckets of liquid, where **exactly one** of the buckets is poisonous. To figure out which one is poisonous, you feed some number of (poor) pigs the liquid to see whether they will die or not. Unfortunately, you only have `minutesToTest` minutes to determine which bucket is poisonous.

You can feed the pigs according to these steps:

1. Choose some live pigs to feed.
2. For each pig, choose which buckets to feed it. The pig will consume all the chosen buckets simultaneously and will take no time. Each pig can feed from any number of buckets, and each bucket can be fed from by any number of pigs.
3. Wait for `minutesToDie` minutes. You may **not** feed any other pigs during this time.
4. After `minutesToDie` minutes have passed, any pigs that have been fed the poisonous bucket will die, and all others will survive.
5. Repeat this process until you run out of time.

Given `buckets`, `minutesToDie`, and `minutesToTest`, return *the **minimum** number of pigs needed to figure out which bucket is poisonous within the allotted time*.

**Example 1:**

```
Input: buckets = 1000, minutesToDie = 15, minutesToTest = 60
Output: 5
```

**Example 2:**

```
Input: buckets = 4, minutesToDie = 15, minutesToTest = 15
Output: 2
```

**Example 3:**

```
Input: buckets = 4, minutesToDie = 15, minutesToTest = 30
Output: 2
```

**Constraints:**

- `1 <= buckets <= 1000`
- `1 <= minutesToDie <= minutesToTest <= 100`

---

## Company Frequency

| Company    | Frequency        | Last Reported |
|------------|------------------|---------------|
| Google     | ★★★☆☆ Medium     | 2023          |
| Amazon     | ★★☆☆☆ Low        | 2023          |
| Microsoft  | ★★☆☆☆ Low        | 2022          |
| Apple      | ★★☆☆☆ Low        | 2022          |
| Two Sigma  | ★★☆☆☆ Low        | 2021          |

> ⚠️ Frequency data is crowd-sourced from LeetCode Discuss, Glassdoor, and
> community interview reports. Treat as a signal, not a guarantee.

---

## DSA Concepts Used

- **Math / Number Theory** — the answer is a pure information-theoretic count: each pig is a base-`(rounds+1)` digit, and we need the smallest `p` with `(rounds+1)^p ≥ buckets`, i.e. `ceil(log_base(buckets))` → see [`/dsa/math_number_theory.md`](/dsa/math_number_theory.md)
- **Combinatorics / Counting** — each pig is one independent radix digit, so `states^pigs` distinguishable outcomes; finding the fewest pigs is a pure counting / information-theory argument → see [`/dsa/combinatorics.md`](/dsa/combinatorics.md)

---

## Approaches Overview

| # | Approach | Time | Space | When to use |
|---|----------|------|-------|-------------|
| 1 | Iterative state counting | O(log₍base₎ buckets) | O(1) | Loop that avoids floating point; the safe interview answer |
| 2 | Closed-form logarithm (Optimal) | O(1) | O(1) | One-line answer; needs an epsilon to dodge FP rounding |

---

## Approach 1 — Iterative Counting of States

### Intuition

Think of each pig as an independent sensor. Over the whole experiment there are `rounds = minutesToTest / minutesToDie` feeding rounds. A single pig ends in one of `rounds + 1` distinguishable states: it dies in round 1, in round 2, …, in round `rounds`, or it survives everything. Call `base = rounds + 1`. With `p` pigs, the joint outcome is a vector of `p` states, giving `base^p` distinct outcome patterns — and each pattern can be pre-assigned to a different bucket (bucket number written in base-`base`, one digit per pig, telling us which rounds to feed that pig). So we need the smallest `p` with `base^p ≥ buckets`. Compute it by multiplying `reachable` by `base` one pig at a time.

### Algorithm

1. `base = minutesToTest / minutesToDie + 1`.
2. `pigs = 0`, `reachable = 1` (with 0 pigs we can distinguish exactly 1 bucket).
3. While `reachable < buckets`: `reachable *= base`, `pigs++`.
4. Return `pigs`.

### Complexity

- **Time:** O(log₍base₎ buckets) — one multiplication per pig; the pig count grows logarithmically in `buckets`.
- **Space:** O(1) — two integer counters.

### Code

```go
func iterativeStates(buckets int, minutesToDie int, minutesToTest int) int {
	base := minutesToTest/minutesToDie + 1 // distinguishable states per pig
	pigs := 0                              // pigs used so far
	reachable := 1                         // buckets distinguishable = base^pigs
	// Grow the reach one pig at a time until it covers every bucket.
	for reachable < buckets {
		reachable *= base // adding one pig multiplies the outcome space by `base`
		pigs++            // account for that pig
	}
	return pigs
}
```

### Dry Run

Example 1: `buckets = 1000, minutesToDie = 15, minutesToTest = 60`.
`base = 60/15 + 1 = 4 + 1 = 5`. Start `pigs = 0`, `reachable = 1`.

| Step | reachable < 1000? | reachable *= 5 | pigs after |
|------|-------------------|----------------|------------|
| 1 | 1 < 1000 yes | 5 | 1 |
| 2 | 5 < 1000 yes | 25 | 2 |
| 3 | 25 < 1000 yes | 125 | 3 |
| 4 | 125 < 1000 yes | 625 | 4 |
| 5 | 625 < 1000 yes | 3125 | 5 |
| 6 | 3125 < 1000 no | — | stop |

`5^5 = 3125 ≥ 1000` first happens at 5 pigs → return `5` ✔

---

## Approach 2 — Closed-Form Logarithm (Optimal)

### Intuition

`base^p ≥ buckets` is equivalent to `p ≥ log_base(buckets) = ln(buckets) / ln(base)`. The minimum integer `p` is the ceiling of that ratio. This is the loop collapsed into a single expression. The only trap is floating-point: an exact power such as `1000 = 10^3` may compute as `2.9999999…`, which would ceil to `4` instead of `3`. Subtracting a tiny epsilon (`1e-9`) before the ceiling absorbs that noise. The `buckets == 1` case is handled explicitly (0 pigs, and `log(1) = 0`).

### Algorithm

1. If `buckets == 1`, return `0`.
2. `base = minutesToTest / minutesToDie + 1`.
3. `ratio = ln(buckets) / ln(base)`.
4. Return `ceil(ratio − 1e-9)`.

### Complexity

- **Time:** O(1) — two logarithms and one ceiling.
- **Space:** O(1).

### Code

```go
func logClosedForm(buckets int, minutesToDie int, minutesToTest int) int {
	if buckets == 1 {
		return 0 // only one bucket → it is trivially the poisonous one, no test needed
	}
	base := float64(minutesToTest/minutesToDie + 1) // states per pig, as float
	// p = ceil(log_base(buckets)); the epsilon absorbs floating-point noise so
	// exact powers (e.g. 1000 = 10^3) don't spuriously round up.
	ratio := math.Log(float64(buckets)) / math.Log(base)
	return int(math.Ceil(ratio - 1e-9))
}
```

### Dry Run

Example 1: `buckets = 1000, minutesToDie = 15, minutesToTest = 60`.

| Step | value |
|------|-------|
| buckets == 1? | no |
| base | 60/15 + 1 = 5 |
| ln(1000) | ≈ 6.907755 |
| ln(5) | ≈ 1.609438 |
| ratio = ln(1000)/ln(5) | ≈ 4.29203 |
| ceil(4.29203 − 1e-9) | 5 |

Return `5` ✔ (matches Approach 1; the epsilon is irrelevant here since the ratio isn't an integer, but it protects cases like `buckets = 125, base = 5` where the exact `3.0` must not round to `4`).

---

## Key Takeaways

- **Reframe "minimum resource" as "encode enough states."** Each pig contributes `rounds + 1` independent states; `p` pigs give `(rounds+1)^p` outcomes. The answer is the smallest exponent covering all buckets — a base-`(rounds+1)` digit-counting argument.
- **`rounds = minutesToTest / minutesToDie`** — the "+1" (surviving all rounds) is the state that trips people up. Miss it and you overshoot the pig count.
- **`ceil(log_b(N))` = smallest p with `b^p ≥ N`.** A reusable identity; pair it with an epsilon whenever floating-point logs feed an integer ceiling.
- **The iterative form is FP-free and interview-safe.** Prefer multiplying integers when correctness near exact powers matters; the closed form is elegant but must be guarded.

---

## Related Problems

- LeetCode #375 — Guess Number Higher or Lower II (information / decision cost)
- LeetCode #887 — Super Egg Drop (minimise trials, DP on states — same "states per trial" spirit)
- LeetCode #464 — Can I Win (game-state reasoning)
- LeetCode #390 — Elimination Game (mathematical reduction over rounds)
- LeetCode #843 — Guess the Word (information-theoretic guessing)
