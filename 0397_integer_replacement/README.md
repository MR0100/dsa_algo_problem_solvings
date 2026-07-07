# 0397 — Integer Replacement

> LeetCode #397 · Difficulty: Medium
> **Categories:** Greedy, Bit Manipulation, Memoization, Dynamic Programming

---

## Problem Statement

Given a positive integer `n`, you can apply one of the following operations:

1. If `n` is even, replace `n` with `n / 2`.
2. If `n` is odd, replace `n` with either `n + 1` or `n - 1`.

Return *the minimum number of operations needed for* `n` *to become* `1`.

**Example 1:**

```
Input: n = 8
Output: 3
Explanation: 8 -> 4 -> 2 -> 1
```

**Example 2:**

```
Input: n = 7
Output: 4
Explanation: 7 -> 8 -> 4 -> 2 -> 1
or 7 -> 6 -> 3 -> 2 -> 1
```

**Example 3:**

```
Input: n = 4
Output: 2
```

**Constraints:**

- `1 <= n <= 2^31 - 1`

---

## Company Frequency

| Company    | Frequency        | Last Reported |
|------------|------------------|---------------|
| Google     | ★★★☆☆ Medium     | 2024          |
| Amazon     | ★★★☆☆ Medium     | 2023          |
| Microsoft  | ★★☆☆☆ Low        | 2023          |
| Bloomberg  | ★★☆☆☆ Low        | 2022          |
| Apple      | ★☆☆☆☆ Rare       | 2022          |

> ⚠️ Frequency data is crowd-sourced from LeetCode Discuss, Glassdoor, and
> community interview reports. Treat as a signal, not a guarantee.

---

## DSA Concepts Used

- **Bit Manipulation** — the optimal solution reasons over the two lowest bits to decide whether `+1` (carry-away trailing 1s) or `-1` (clear a lone 1) reaches `1` fastest → see [`/dsa/bit_manipulation.md`](/dsa/bit_manipulation.md)
- **Greedy** — at every odd step we greedily pick the move that eliminates the most low set bits; a local optimum is provably global here → see [`/dsa/greedy.md`](/dsa/greedy.md)
- **Recursion / Top-Down search** — the branch-both-ways recursion is the natural baseline before the greedy insight → see [`/dsa/dynamic_programming_1d.md`](/dsa/dynamic_programming_1d.md)

---

## Approaches Overview

| # | Approach | Time | Space | When to use |
|---|----------|------|-------|-------------|
| 1 | Recursion (Top-Down) | O(2^log n) worst branch | O(log n) | Clear, correct baseline; each branch still halves quickly |
| 2 | Greedy Bit Manipulation (Optimal) | O(log n) | O(1) | Always — one linear-in-bits pass, overflow-safe with `uint` |

---

## Approach 1 — Recursion (Top-Down)

### Intuition

The even case is forced: halving is the only legal move. The odd case is a fork — `+1` or `-1` — and we can't tell locally which leads to fewer steps, so try both and take the minimum. Every path shrinks the number (a halving always follows an odd step), so the recursion terminates at `n == 1`.

The one gotcha: `n` can be `2^31 − 1`, and computing `n + 1` overflows a 32-bit `int`. Working in `uint` sidesteps that.

### Algorithm

1. Base case: `n == 1` → `0` operations.
2. If `n` is even → `1 + solve(n/2)`.
3. If `n` is odd → `1 + min(solve(n+1), solve(n-1))`.

### Complexity

- **Time:** O(2^log n) in the worst odd-heavy branch (no memoization); in practice far less because even steps don't branch and each branch halves.
- **Space:** O(log n) — recursion-stack depth proportional to the bit length.

### Code

```go
func recursion(n int) int {
	var solve func(x uint) int
	solve = func(x uint) int {
		if x == 1 {
			return 0 // reached the target, no more operations
		}
		if x%2 == 0 {
			return 1 + solve(x/2) // even: the only legal move is to halve
		}
		// odd: try both neighbours (each becomes even) and take the cheaper.
		return 1 + min(solve(x+1), solve(x-1))
	}
	return solve(uint(n))
}
```

### Dry Run

Example 2: `n = 7`.

| Call | x | branch | contributes |
|------|---|--------|-------------|
| solve(7) | 7 odd | min(solve(8), solve(6)) | 1 + min(...) |
| solve(8) | 8 even | solve(4) | 1 + ... |
| solve(4) | 4 even | solve(2) | 1 + ... |
| solve(2) | 2 even | solve(1) | 1 + ... |
| solve(1) | 1 | base | 0 |
| ⇒ solve(8) | | 8→4→2→1 | 3 |
| solve(6) | 6 even | solve(3) | 1 + ... |
| solve(3) | 3 odd | min(solve(4), solve(2)) → solve(2) cheaper | 1 + 1 |
| ⇒ solve(6) | | 6→3→2→1 | 3 |

`solve(7) = 1 + min(3, 3) = 4`. Result: `4` ✔ (both `7→8→4→2→1` and `7→6→3→2→1` cost 4).

---

## Approach 2 — Greedy Bit Manipulation (Optimal)

### Intuition

Halving strips a trailing `0`, so we want to *reach* numbers with many trailing zeros. For odd `n`, the choice `+1` vs `-1` should clear the most low set bits:

- **`n == 3`** is special: `3 → 2 → 1` (subtract) beats `3 → 4 → 2 → 1` (add). Always subtract.
- **Second-lowest bit is 1** (n ends in `…11`): adding 1 triggers a carry that clears a whole run of trailing 1s → do `n++`.
- **Second-lowest bit is 0** (n ends in `…01`): subtracting 1 clears the single low 1 without disturbing higher bits → do `n--`.

Each step is followed by a halving, so we march down to 1 in a number of steps linear in the bit length. `uint` again guards the `2^31 − 1` overflow.

### Algorithm

1. `count = 0`, work in `uint`.
2. While `n != 1`:
   1. If `n` even → `n >>= 1`.
   2. Else if `n == 3` **or** bit 1 is 0 (`n & 2 == 0`) → `n--`.
   3. Else → `n++`.
   4. `count++`.
3. Return `count`.

### Complexity

- **Time:** O(log n) — each iteration reduces the number's magnitude by at least one bit.
- **Space:** O(1) — a single counter, no recursion.

### Code

```go
func greedyBits(n int) int {
	x := uint(n)
	count := 0
	for x != 1 {
		if x%2 == 0 {
			x >>= 1 // even: halve (drop the trailing 0)
		} else if x == 3 || x&2 == 0 {
			// x == 3: subtract to avoid the longer 3->4->2->1 path.
			// x&2 == 0 (ends in 01): subtracting clears the single low 1 bit.
			x--
		} else {
			// ends in 11: adding 1 carries and clears a run of trailing 1s.
			x++
		}
		count++ // every branch above performed one operation
	}
	return count
}
```

### Dry Run

Example 2: `n = 7` (binary `111`).

| Step | x (bin) | condition | action | x after | count |
|------|---------|-----------|--------|---------|-------|
| 1 | `111` (7) | odd, not 3, bit1=1 (`…11`) | x++ | `1000` (8) | 1 |
| 2 | `1000` (8) | even | x >>= 1 | `100` (4) | 2 |
| 3 | `100` (4) | even | x >>= 1 | `10` (2) | 3 |
| 4 | `10` (2) | even | x >>= 1 | `1` (1) | 4 |

Loop ends (`x == 1`). Result: `4` ✔ — the greedy `+1` on `111` cleared both trailing 1s at once.

---

## Key Takeaways

- **Trailing-bit reasoning** decides the optimal odd move: `…11` → add (carry clears a run), `…01` → subtract (removes the lone 1). This "make more trailing zeros" heuristic generalises to many halving puzzles.
- **`n == 3` is the lone exception** to the rule — always subtract. Memorise it; it's the classic edge that breaks a naive greedy.
- **Always use `uint`** (or `int64`) when the input can be `2^31 − 1` and you compute `n + 1` — otherwise you overflow to a negative number and loop forever/incorrectly.
- Branch-and-take-min recursion is a fine baseline, but recognising the bit pattern removes the branching entirely for an O(log n), O(1) solution.

---

## Related Problems

- LeetCode #191 — Number of 1 Bits (`x & (x-1)` bit-clearing)
- LeetCode #231 — Power of Two (trailing-bit structure)
- LeetCode #201 — Bitwise AND of Numbers Range (common-prefix bit reasoning)
- LeetCode #1611 — Minimum One Bit Operations (bit-by-bit reduction to a target)
