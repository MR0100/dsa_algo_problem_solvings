# 0374 — Guess Number Higher or Lower

> LeetCode #374 · Difficulty: Easy
> **Categories:** Binary Search, Interactive

---

## Problem Statement

We are playing the Guess Game. The game is as follows:

I pick a number from `1` to `n`. You have to guess which number I picked.

Every time you guess wrong, I will tell you whether the number I picked is higher or lower than your guess.

You call a pre-defined API `int guess(int num)`, which returns three possible results:

- `-1`: Your guess is higher than the number I picked (i.e. `num > pick`).
- `1`: Your guess is lower than the number I picked (i.e. `num < pick`).
- `0`: Your guess is equal to the number I picked (i.e. `num == pick`).

Return *the number that I picked.*

**Example 1:**

```
Input: n = 10, pick = 6
Output: 6
```

**Example 2:**

```
Input: n = 1, pick = 1
Output: 1
```

**Example 3:**

```
Input: n = 2, pick = 1
Output: 1
```

**Constraints:**

- `1 <= n <= 2^31 - 1`
- `1 <= pick <= n`

---

## Company Frequency

| Company    | Frequency        | Last Reported |
|------------|------------------|---------------|
| Google     | ★★★★☆ High       | 2024          |
| Amazon     | ★★★☆☆ Medium     | 2023          |
| Apple      | ★★★☆☆ Medium     | 2023          |
| Microsoft  | ★★☆☆☆ Low        | 2022          |

> ⚠️ Frequency data is crowd-sourced from LeetCode Discuss, Glassdoor, and community interview reports. Treat as a signal, not a guarantee.

---

## DSA Concepts Used

- **Binary Search** — the higher/lower response is a monotone predicate that halves the search space each guess → see [`/dsa/binary_search.md`](/dsa/binary_search.md)
- **Overflow-safe midpoint** — use `lo + (hi-lo)/2` because `n` can be near `2^31 - 1` → see [`/dsa/math_number_theory.md`](/dsa/math_number_theory.md)

---

## Approaches Overview

| # | Approach | Time | Space | When to use |
|---|----------|------|-------|-------------|
| 1 | Linear scan | O(n) | O(1) | Baseline; ignores the hint |
| 2 | Binary search | O(log n) | O(1) | Optimal; uses the higher/lower feedback |

---

## Approach 1 — Linear Scan

### Intuition
`guess(num) == 0` uniquely identifies the picked number, so walking upward from 1 must eventually hit it. This ignores the ordering information the API gives us and is only a baseline.

### Algorithm
1. For `num = 1..n`: if `guess(num) == 0` return `num`.
2. (Unreachable, per constraints.) Return `-1`.

### Complexity
- **Time:** O(n) — up to `n` guesses.
- **Space:** O(1).

### Code
```go
func linearScan(n int) int {
	for num := 1; num <= n; num++ { // try every candidate in order
		if guess(num) == 0 { // found the picked number
			return num
		}
	}
	return -1 // per constraints the pick is always in [1, n]
}
```

### Dry Run
`n = 10`, `pick = 6`:

| num | guess(num) | action |
|-----|------------|--------|
| 1 | 1 (too low) | continue |
| 2 | 1 | continue |
| 3 | 1 | continue |
| 4 | 1 | continue |
| 5 | 1 | continue |
| 6 | 0 | **return 6** |

Return `6`. ✓

---

## Approach 2 — Binary Search (Optimal)

### Intuition
The response tells us which half of `[lo, hi]` contains the pick, so this is a textbook binary search over a monotone predicate. Halving the range each guess reaches the answer in O(log n). Compute the midpoint as `lo + (hi-lo)/2` to avoid integer overflow when `n` is near `2^31 - 1`.

### Algorithm
1. `lo = 1`, `hi = n`.
2. While `lo <= hi`:
   1. `mid = lo + (hi-lo)/2`.
   2. `r = guess(mid)`. If `r == 0` return `mid`.
   3. If `r < 0` (`mid` too high) → `hi = mid - 1`; else (`mid` too low) → `lo = mid + 1`.
3. (Unreachable.) Return `-1`.

### Complexity
- **Time:** O(log n) — the range halves each step.
- **Space:** O(1).

### Code
```go
func binarySearch(n int) int {
	lo, hi := 1, n
	for lo <= hi {
		mid := lo + (hi-lo)/2 // overflow-safe midpoint
		switch guess(mid) {
		case 0: // exact hit
			return mid
		case -1: // mid is higher than pick → search lower half
			hi = mid - 1
		default: // guess returned 1: mid is lower → search upper half
			lo = mid + 1
		}
	}
	return -1
}
```

### Dry Run
`n = 10`, `pick = 6`:

| Step | lo | hi | mid = lo+(hi-lo)/2 | guess(mid) | action |
|------|----|----|--------------------|------------|--------|
| 1 | 1 | 10 | 5 | 1 (mid too low) | lo = 6 |
| 2 | 6 | 10 | 8 | -1 (mid too high) | hi = 7 |
| 3 | 6 | 7 | 6 | 0 (match) | **return 6** |

Return `6`. ✓

---

## Key Takeaways
- **Interactive higher/lower search is just binary search**: the API's response is the branch selector.
- **Always use `lo + (hi-lo)/2`** for the midpoint when bounds can approach `INT_MAX` — `(lo+hi)/2` overflows.
- Watch the sign convention: here `guess` returns `-1` when your guess is **too high** (opposite of some intuitions), so map `-1 → hi = mid-1`.

---

## Related Problems
- LeetCode #278 — First Bad Version (same interactive binary search shape)
- LeetCode #704 — Binary Search (canonical form)
- LeetCode #375 — Guess Number Higher or Lower II (minimax DP variant)
- LeetCode #35 — Search Insert Position
