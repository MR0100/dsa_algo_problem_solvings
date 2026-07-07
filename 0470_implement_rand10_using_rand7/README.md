# 0470 — Implement Rand10() Using Rand7()

> LeetCode #470 · Difficulty: Medium
> **Categories:** Math, Randomized, Rejection Sampling, Probability and Statistics

---

## Problem Statement

Given the **API** `rand7()` that generates a uniform random integer in the range `[1, 7]`, write a function `rand10()` that generates a uniform random integer in the range `[1, 10]`. You can only call the API `rand7()`, and you shouldn't call any other API. Please **do not** use a language's built-in random API.

Each test case will have one **internal** argument `n`, the number of times that your implemented function `rand10()` will be called while testing. Note that this is **not an argument** passed to `rand10()`.

**Example 1:**

```
Input: n = 1
Output: [2]
```

**Example 2:**

```
Input: n = 2
Output: [2,8]
```

**Example 3:**

```
Input: n = 3
Output: [3,8,10]
```

**Constraints:**

- `1 <= n <= 10^5`

**Follow up:**

- What is the [expected value](https://en.wikipedia.org/wiki/Expected_value) for the number of calls to `rand7()` function?
- Could you minimize the number of calls to `rand7()`?

> Note: the outputs in the examples above are just *sample* results of a random function; there is no fixed correct sequence. A solution is correct iff every output lies in `[1, 10]` and the distribution is uniform.

---

## Company Frequency

| Company    | Frequency        | Last Reported |
|------------|------------------|---------------|
| Google     | ★★★★☆ High       | 2024          |
| ByteDance  | ★★★★☆ High       | 2024          |
| Bloomberg  | ★★★☆☆ Medium     | 2023          |
| LinkedIn   | ★★★☆☆ Medium     | 2023          |
| Microsoft  | ★★★☆☆ Medium     | 2023          |
| Amazon     | ★★☆☆☆ Low        | 2022          |

> ⚠️ Frequency data is crowd-sourced from LeetCode Discuss, Glassdoor, and
> community interview reports. Treat as a signal, not a guarantee.

---

## DSA Concepts Used

- **Rejection sampling / uniform generation via modular reduction** — build a larger uniform range from a smaller one (7 → 49), keep the largest multiple of the target (40), reject the rest to stay exactly uniform, then reduce mod 10; the correctness argument is pure modular arithmetic and probability → see [`/dsa/math_number_theory.md`](/dsa/math_number_theory.md)

---

## Approaches Overview

| # | Approach | Expected rand7() calls | Space | When to use |
|---|----------|------------------------|-------|-------------|
| 1 | Rejection Sampling on a 7×7 Grid (Optimal, standard) | ≈ 2.45 | O(1) | The canonical answer; simplest to prove uniform |
| 2 | Rejection Sampling that Reuses Rejects | ≈ 2.19 | O(1) | The follow-up "minimise calls"; recycles rejected entropy |

Both are **expected O(1)** time (each trial succeeds with high probability); neither has a hard worst-case bound because, in theory, rejection can repeat — but the probability decays geometrically.

---

## Approach 1 — Rejection Sampling on a 7×7 Grid (Optimal, standard)

### Intuition

One `rand7()` gives 7 equally likely values; two independent calls give a uniform point on a `7 × 7 = 49`-cell grid, i.e. a uniform integer in `[1, 49]`. `40` is the largest multiple of `10` that is `≤ 49`, so keep only values in `[1, 40]` (each of the 10 outputs then has exactly 4 pre-images) and **reject** `41..49`, retrying. Rejection is essential: `49` is not divisible by `10`, and folding the leftover 9 values back into `[1,10]` would over-represent some outcomes and break uniformity.

### Algorithm

1. `row = rand7()` (1..7), `col = rand7()` (1..7).
2. `idx = (row−1)·7 + col` → uniform integer in `[1, 49]`.
3. If `idx > 40`, discard and go back to step 1 (rejection).
4. Return `(idx−1) % 10 + 1` → uniform integer in `[1, 10]`.

### Complexity

- **Time:** Expected O(1) — each trial succeeds with probability `40/49 ≈ 0.816`; expected trials `= 49/40 ≈ 1.225`, i.e. `≈ 2.45` `rand7()` calls (this answers the first follow-up).
- **Space:** O(1) — a few local integers.

### Code

```go
func rejectionSampling() int {
	for {
		row := rand7()             // 1..7 — chooses the grid row
		col := rand7()             // 1..7 — chooses the grid column
		idx := (row-1)*7 + col     // 1..49, uniform over the 7×7 grid
		if idx <= 40 {             // keep only the first 40 (a multiple of 10)
			return (idx-1)%10 + 1  // fold 40 outcomes → [1,10], 4 each, unbiased
		}
		// idx in 41..49 → reject and retry so the distribution stays exact
	}
}
```

### Dry Run

Suppose the first two `rand7()` calls return `row = 6`, `col = 3`:

| Step | Computation | Value |
|------|-------------|-------|
| 1 | `row = rand7()` | 6 |
| 1 | `col = rand7()` | 3 |
| 2 | `idx = (6−1)·7 + 3 = 35 + 3` | 38 |
| 3 | `38 ≤ 40`? | yes → accept |
| 4 | `(38−1) % 10 + 1 = 37 % 10 + 1 = 7 + 1` | **8** |

Returns `8`. Had `idx` been, say, `44` (e.g. `row = 7, col = 2 → 44`), step 3 would reject it and the loop would draw two fresh values. Because each accepted `idx ∈ [1,40]` is equally likely and the 40 values split evenly into 10 residue classes, every output `1..10` occurs with probability exactly `1/10`.

---

## Approach 2 — Rejection Sampling that Reuses Rejects (Fewer rand7() Calls)

### Intuition

When `idx ∈ [41, 49]` we reject in Approach 1 — but that rejected `idx` still hides a **uniform** integer in `[1, 9]` (`idx − 40`). Don't waste it: combine it with a fresh `rand7()` to form a uniform `[1, 63]`, keep `[1, 60]` (→ `[1,10]`), and if that also lands in the reject zone `[61,63]` you still hold a uniform `[1, 3]`, which combines with one more `rand7()` into `[1, 21]`; keep `[1, 20]`. Recycling the leftover entropy at each stage means fewer cold restarts, lowering the expected `rand7()` count.

### Algorithm

1. `a = rand7()`, `b = rand7()`; `idx = (a−1)·7 + b` (1..49). If `idx ≤ 40`, return `(idx−1)%10 + 1`.
2. Else `rem = idx − 40` (uniform 1..9); `c = rand7()`; `idx = (rem−1)·7 + c` (1..63). If `idx ≤ 60`, return `(idx−1)%10 + 1`.
3. Else `rem = idx − 60` (uniform 1..3); `d = rand7()`; `idx = (rem−1)·7 + d` (1..21). If `idx ≤ 20`, return `(idx−1)%10 + 1`.
4. Else (only value 21 remains) loop back to step 1.

### Complexity

- **Time:** Expected `≈ 2.19` `rand7()` calls per `rand10()` — lower than Approach 1's `≈ 2.45` because rejected randomness is reused (this answers the second follow-up). Still expected O(1).
- **Space:** O(1).

### Code

```go
func reuseRejects() int {
	for {
		a := rand7()
		b := rand7()
		idx := (a-1)*7 + b // 1..49
		if idx <= 40 {
			return (idx-1)%10 + 1
		}
		// Salvage the uniform [1,9] hiding in the rejected 41..49.
		rem := idx - 40    // 1..9
		c := rand7()
		idx = (rem-1)*7 + c // 1..63
		if idx <= 60 {
			return (idx-1)%10 + 1
		}
		// Salvage the uniform [1,3] hiding in the rejected 61..63.
		rem = idx - 60     // 1..3
		d := rand7()
		idx = (rem-1)*7 + d // 1..21
		if idx <= 20 {
			return (idx-1)%10 + 1
		}
		// Only 1 value left (21); nothing to salvage — loop and start over.
	}
}
```

### Dry Run

Suppose `a = 7, b = 6`, then `c = 4`:

| Step | Computation | Value | Decision |
|------|-------------|-------|----------|
| 1 | `idx = (7−1)·7 + 6 = 42 + 6` | 48 | `48 > 40` → reject, salvage |
| 2 | `rem = 48 − 40` | 8 | uniform in [1,9] |
| 2 | `c = rand7()` | 4 | — |
| 2 | `idx = (8−1)·7 + 4 = 49 + 4` | 53 | `53 ≤ 60` → accept |
| 2 | `(53−1) % 10 + 1 = 52 % 10 + 1 = 2 + 1` | **3** | return 3 |

Returns `3` using **3** `rand7()` calls, whereas Approach 1 would have discarded the `48` and restarted, spending more calls on average. Every accepted value across all three stages is still uniform over a range whose kept prefix is a multiple of 10, so the output stays exactly uniform on `[1,10]`.

---

## Key Takeaways

- **Build up, then reject down.** To turn `rand(a)` into `rand(b)` with `b > a`, stack independent draws to reach a range `R ≥ b`, keep the largest multiple of `b` that is `≤ R`, reject the tail, and reduce mod `b`. The rejection is what preserves exact uniformity.
- **Never fold the remainder back in.** Because `49 % 10 ≠ 0`, reusing `41..49` as-is would bias the result. Any reuse must first re-expand the leftover into a fresh, properly-sized uniform range (Approach 2).
- **Expected-value analysis of rejection:** with success probability `p`, the expected number of trials is `1/p` (geometric distribution). Here `p = 40/49`, giving `≈ 1.225` trials × 2 calls `≈ 2.45` `rand7()` calls.
- **Minimising calls = recycling entropy.** The rejected value is not "no information" — it is a uniform sample over the reject range, worth reusing. This is the heart of the follow-up.
- **Testing a randomized function:** you cannot assert a fixed output. Assert the *invariants* instead — range membership (`1..10`) and approximate uniformity over a large sample (bucket counts within a tolerance of `N/10`), which is exactly what `main()` does.

---

## Related Problems

- LeetCode #478 — Generate Random Point in a Circle (rejection sampling in 2-D)
- LeetCode #497 — Random Point in Non-overlapping Rectangles (weighted uniform sampling)
- LeetCode #528 — Random Pick with Weight (prefix sums + binary search over a distribution)
- LeetCode #382 — Linked List Random Node (reservoir sampling — a different uniform trick)
- LeetCode #384 — Shuffle an Array (Fisher–Yates; uniform permutations)
