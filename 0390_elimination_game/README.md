# 0390 — Elimination Game

> LeetCode #390 · Difficulty: Medium
> **Categories:** Math, Simulation, Recursion

---

## Problem Statement

You have a list `arr` of all integers in the range `[1, n]` sorted in a strictly increasing order. Apply the following algorithm on `arr`:

- Starting from left to right, remove the first number and every other number afterward until you reach the end of the list.
- Repeat the previous step again, but this time from right to left, remove the rightmost number and every other number from the remaining numbers.
- Keep repeating the steps again, alternating left to right and right to left, until a single number remains.

Given the integer `n`, return the last number that remains in `arr`.

**Example 1:**

```
Input: n = 9
Output: 6
Explanation:
arr = [1, 2, 3, 4, 5, 6, 7, 8, 9]
arr = [2, 4, 6, 8]
arr = [2, 6]
arr = [6]
```

**Example 2:**

```
Input: n = 1
Output: 1
```

**Constraints:**

- `1 <= n <= 10^9`

---

## Company Frequency

| Company    | Frequency        | Last Reported |
|------------|------------------|---------------|
| Amazon     | ★★★☆☆ Medium     | 2023          |
| Google     | ★★★☆☆ Medium     | 2023          |
| Adobe      | ★★☆☆☆ Low        | 2022          |
| Microsoft  | ★★☆☆☆ Low        | 2022          |
| Bloomberg  | ★☆☆☆☆ Rare       | 2021          |

> ⚠️ Frequency data is crowd-sourced from LeetCode Discuss, Glassdoor, and
> community interview reports. Treat as a signal, not a guarantee.

---

## DSA Concepts Used

- **Math / Pattern Reasoning** — instead of storing the list, track only the value of the leftmost survivor (`head`), the gap (`step`) between survivors, and the count; derive a closed recurrence for how `head` moves → see [`/dsa/math_number_theory.md`](/dsa/math_number_theory.md)
- **Simulation** — the brute-force baseline directly enacts the described elimination procedure → see [`/dsa/string_algorithms.md`](/dsa/string_algorithms.md)

---

## Approaches Overview

| # | Approach | Time | Space | When to use |
|---|----------|------|-------|-------------|
| 1 | Brute Force Simulation | O(n) | O(n) | Small n; verifying the pattern |
| 2 | Track Head Pointer (Optimal) | O(log n) | O(1) | n up to 10⁹ — the intended solution |

---

## Approach 1 — Brute Force Simulation

### Intuition

Do exactly what the problem says: hold the surviving numbers in a slice and, each pass, keep every other element — flipping the sweep direction each round. On a left→right pass the first element is always removed (keep indices `1,3,5,…`). On a right→left pass the pattern mirrors and depends on the current length's parity. When one element remains, it is the answer. Correct, but O(n) memory — impossible for `n = 10⁹`.

### Algorithm

1. `nums = [1..n]`; `leftToRight = true`.
2. While `len(nums) > 1`:
   1. If `leftToRight`: keep indices `1,3,5,…`.
   2. Else: keep starting at index `0` when length is even, or index `1` when odd, then every other.
   3. Flip `leftToRight`.
3. Return `nums[0]`.

### Complexity

- **Time:** O(n) — the list halves every pass, so total work is `n + n/2 + … ≈ 2n`.
- **Space:** O(n) — the full list is materialised.

### Code

```go
func bruteForce(n int) int {
	nums := make([]int, n) // the surviving numbers, initially 1..n
	for i := 0; i < n; i++ {
		nums[i] = i + 1
	}
	leftToRight := true // direction of the current pass

	for len(nums) > 1 {
		next := make([]int, 0, len(nums)/2) // survivors of this pass
		if leftToRight {
			// Left→right: always remove the first, so keep indices 1,3,5,...
			for i := 1; i < len(nums); i += 2 {
				next = append(next, nums[i])
			}
		} else {
			// Right→left: mirror image. If length is even, we keep 0,2,4,...;
			// if odd, we keep 1,3,5,... (the leftmost is removed).
			start := 0
			if len(nums)%2 == 1 {
				start = 1
			}
			for i := start; i < len(nums); i += 2 {
				next = append(next, nums[i])
			}
		}
		nums = next             // advance to the survivors
		leftToRight = !leftToRight // alternate the sweep direction
	}
	return nums[0]
}
```

### Dry Run

Input `n = 9`:

| pass | direction | nums before | survivors kept | nums after |
|------|-----------|-------------|----------------|------------|
| 1 | L→R | `[1..9]` (len 9) | indices 1,3,5,7 | `[2,4,6,8]` |
| 2 | R→L | `[2,4,6,8]` (len 4, even) | indices 0,2 | `[2,6]` |
| 3 | L→R | `[2,6]` (len 2) | index 1 | `[6]` |
| — | — | len 1 | stop | **6** |

Answer: `6`.

---

## Approach 2 — Track Head Pointer (Optimal)

### Intuition

The final answer is wherever the **head** (leftmost survivor) ends up. We never need the list — only `head`, the `step` between consecutive survivors, and how many `remaining` numbers there are. The head moves under two conditions:

- On a **left→right** pass, the head is *always* removed, so `head += step`.
- On a **right→left** pass, the head moves *only if* the count is odd (then the leftmost is also swept away); if the count is even, the head survives untouched.

Each pass halves `remaining` and doubles `step`. When `remaining == 1`, `head` is the answer.

### Algorithm

1. `head = 1`, `step = 1`, `remaining = n`, `leftToRight = true`.
2. While `remaining > 1`:
   1. If `leftToRight` **or** `remaining` is odd: `head += step`.
   2. `remaining /= 2`; `step *= 2`; flip `leftToRight`.
3. Return `head`.

### Complexity

- **Time:** O(log n) — `remaining` halves each iteration.
- **Space:** O(1) — a few scalars.

### Code

```go
func headPointer(n int) int {
	head := 1          // value of the leftmost surviving number
	step := 1          // gap between consecutive survivors
	remaining := n     // how many numbers are still in play
	leftToRight := true

	for remaining > 1 {
		// The head is eliminated (so it must move forward by one step) when:
		//  - we sweep left→right (head is always first to go), OR
		//  - we sweep right→left but the count is odd (leftmost also removed).
		if leftToRight || remaining%2 == 1 {
			head += step
		}
		remaining /= 2       // half the numbers are eliminated this pass
		step *= 2            // survivors are now twice as far apart
		leftToRight = !leftToRight
	}
	return head
}
```

### Dry Run

Input `n = 9`:

| iter | leftToRight | remaining | odd? | head moves? | head after | step after | remaining after |
|------|-------------|-----------|------|-------------|------------|------------|-----------------|
| 1 | true | 9 | yes | yes (L→R) | 1+1 = 2 | 2 | 4 |
| 2 | false | 4 | no | no (R→L, even) | 2 | 4 | 2 |
| 3 | true | 2 | no | yes (L→R) | 2+4 = 6 | 8 | 1 |
| — | — | 1 | — | loop ends | **6** | — | — |

Answer: `6`.

---

## Key Takeaways

- **Don't simulate a huge list — track the invariant.** Here the invariant is `(head, step, remaining)`; the survivors are always an arithmetic progression `head, head+step, head+2·step, …`.
- **Head only ever moves forward.** The leftmost survivor is removed on every L→R pass, and on R→L passes exactly when the count is odd. That single rule collapses the whole simulation to O(log n).
- **Josephus-flavoured problems** reward finding the closed-form movement of a single tracked position rather than materialising all eliminations.

---

## Related Problems

- LeetCode #1823 — Find the Winner of the Circular Game (Josephus problem)
- LeetCode #390 is a linear cousin of the circular Josephus elimination
- LeetCode #294 / #292 — Nim-style elimination reasoning
