# 0202 — Happy Number

> LeetCode #202 · Difficulty: Easy
> **Categories:** Hash Table, Math, Two Pointers

---

## Problem Statement

Write an algorithm to determine if a number `n` is happy.

A **happy number** is a number defined by the following process:

- Starting with any positive integer, replace the number by the sum of the squares of its digits.
- Repeat the process until the number equals 1 (where it will stay), or it **loops endlessly in a cycle** which does not include 1.
- Those numbers for which this process **ends in 1** are happy.

Return `true` *if* `n` *is a happy number, and* `false` *if not*.

**Example 1:**

```
Input: n = 19
Output: true
Explanation:
1² + 9² = 82
8² + 2² = 68
6² + 8² = 100
1² + 0² + 0² = 1
```

**Example 2:**

```
Input: n = 2
Output: false
```

**Constraints:**

- `1 <= n <= 2^31 - 1`

---

## Company Frequency

| Company    | Frequency        | Last Reported |
|------------|------------------|---------------|
| Google     | ★★★★☆ High       | 2024          |
| Amazon     | ★★★★☆ High       | 2024          |
| Apple      | ★★★☆☆ Medium     | 2024          |
| Microsoft  | ★★★☆☆ Medium     | 2023          |
| Adobe      | ★★☆☆☆ Low        | 2023          |
| JPMorgan   | ★★☆☆☆ Low        | 2023          |

> ⚠️ Frequency data is crowd-sourced from LeetCode Discuss, Glassdoor, and
> community interview reports. Treat as a signal, not a guarantee.

---

## DSA Concepts Used

- **Hash Map / Hash Set** — remember every value the trajectory produces; a revisit proves an endless cycle → see [`/dsa/hash_map.md`](/dsa/hash_map.md)
- **Two Pointers (Floyd's tortoise & hare)** — the sequence `n → step(n) → …` is an implicit linked list; detect its cycle in O(1) space exactly like LeetCode #141 → see [`/dsa/two_pointers.md`](/dsa/two_pointers.md)
- **Math / Number Theory (digit manipulation)** — the step function pops digits with `% 10` / `/ 10`, and a size argument proves every trajectory collapses below 1000, making cycles inevitable → see [`/dsa/math_number_theory.md`](/dsa/math_number_theory.md)

---

## Approaches Overview

| # | Approach | Time | Space | When to use |
|---|----------|------|-------|-------------|
| 1 | Hash Set Cycle Detection | O(log n) | O(log n) | Default first answer; simplest to write and to prove |
| 2 | Floyd's Cycle Detection (Optimal) | O(log n) | O(1) | When the interviewer asks to remove the extra memory |
| 3 | Hardcoded Cycle (Math Fact) | O(log n) | O(1) | Fastest constant factor; requires citing the known 4-cycle |

*Why O(log n) time:* one `digitSquareSum` costs O(log₁₀ n) digit operations, and the number of steps before the walk is trapped in the small constant-size zone (< 1000) is bounded — a 10-digit number maps below 811, so after the first couple of steps everything is constant work.

---

## Approach 1 — Hash Set Cycle Detection

### Intuition

The process is a deterministic walk: from `n` there is exactly one next value, `digitSquareSum(n)`. First, the walk cannot escape to infinity — a number with `d ≥ 4` digits maps to at most `81·d`, which has fewer digits (e.g. any 10-digit number maps to ≤ 810), so every trajectory quickly falls into the finite zone `[1, 999]` and stays there. A deterministic walk on a finite set has only two possible fates: it reaches the fixed point `1` (since `1² = 1`, it stays), or it revisits some value and therefore repeats the same loop forever. A hash set of visited values distinguishes the two fates directly.

### Algorithm

1. Create an empty set `seen`.
2. While `n != 1` **and** `n ∉ seen`:
   1. Insert `n` into `seen`.
   2. Replace `n` with `digitSquareSum(n)` (pop digits with `% 10`, square, accumulate, shrink with `/ 10`).
3. The loop ends at `n == 1` (happy) or on a repeat (unhappy). Return `n == 1`.

### Complexity

- **Time:** O(log n) — O(log₁₀ n) digit work for the first step; afterwards the walk lives among numbers < 1000, of which only a bounded constant number can be visited before a repeat.
- **Space:** O(log n) — the set holds the visited chain, whose length is bounded by the same argument.

### Code

```go
func hashSet(n int) bool {
	seen := map[int]bool{} // every value the trajectory has produced
	for n != 1 && !seen[n] {
		seen[n] = true        // record before stepping, so a revisit is caught
		n = digitSquareSum(n) // advance the walk one step
	}
	return n == 1 // loop left either at 1 (happy) or on a repeat (cycle)
}

// digitSquareSum returns the sum of the squares of the decimal digits of n.
func digitSquareSum(n int) int {
	sum := 0
	for n > 0 {
		digit := n % 10      // pop the lowest decimal digit
		sum += digit * digit // accumulate its square
		n /= 10              // shrink the number by one digit
	}
	return sum
}
```

### Dry Run

Example 1: `n = 19`.

| Step | n | n == 1? | n in seen? | seen after insert | digitSquareSum |
|------|-----|---------|------------|----------------------------|----------------|
| 1 | 19 | no | no | {19} | 1² + 9² = 82 |
| 2 | 82 | no | no | {19, 82} | 8² + 2² = 68 |
| 3 | 68 | no | no | {19, 82, 68} | 6² + 8² = 100 |
| 4 | 100 | no | no | {19, 82, 68, 100} | 1² + 0² + 0² = 1 |
| 5 | 1 | yes | — | — | exit loop |

Return `1 == 1` → `true` ✔ (For `n = 2` the walk reaches 4 → 16 → 37 → 58 → 89 → 145 → 42 → 20 → **4 again**, the set lookup fires, and the function returns `false`.)

---

## Approach 2 — Floyd's Cycle Detection (Two Pointers, Optimal)

### Intuition

Reinterpret the value sequence as an **implicit linked list**: node = current value, `Next` = `digitSquareSum`. "Unhappy" means this list flows into a cycle that does not contain 1; "happy" means it flows into the self-loop at 1. That is precisely LeetCode #141 — so use its O(1)-space tool. Advance `slow` one step and `fast` two steps per iteration: if there is a cycle, `fast` laps `slow` and they meet inside it; if the number is happy, `fast` hits the absorbing state 1 first and we stop. No memory of the path is needed because the *relative speed* of the pointers does the remembering.

### Algorithm

1. `slow = n`, `fast = digitSquareSum(n)` (fast starts one step ahead so the loop condition works from the start).
2. While `fast != 1` and `slow != fast`:
   1. `slow = digitSquareSum(slow)` — one step.
   2. `fast = digitSquareSum(digitSquareSum(fast))` — two steps.
3. Return `fast == 1`: reaching 1 means happy; meeting anywhere else proves a 1-free cycle.

### Complexity

- **Time:** O(log n) — same trajectory bound as Approach 1; once both pointers are inside the constant-size cycle, they meet within one lap (< 20 steps in practice).
- **Space:** O(1) — exactly two integer variables; this is the answer to "can you do it without the set?".

### Code

```go
func floydCycleDetection(n int) bool {
	slow := n                 // moves one step at a time
	fast := digitSquareSum(n) // starts one step ahead, moves two at a time
	for fast != 1 && slow != fast {
		slow = digitSquareSum(slow)                 // tortoise: one step
		fast = digitSquareSum(digitSquareSum(fast)) // hare: two steps
	}
	// Either fast reached the fixed point 1 (happy), or the pointers met
	// inside a cycle that never contains 1 (unhappy).
	return fast == 1
}
```

### Dry Run

Example 1: `n = 19`. Chain: 19 → 82 → 68 → 100 → 1 → 1 → …

| Step | slow | fast | fast == 1? | slow == fast? | Action |
|------|------|------|------------|---------------|--------|
| 0 | 19 | 82 | no | no | enter loop |
| 1 | 82 | step(step(82)) = step(68) = 100 | no | no | continue |
| 2 | 68 | step(step(100)) = step(1) = 1 | **yes** | — | exit loop |

Return `fast == 1` → `true` ✔ (For `n = 2` the pointers meet at 42 inside the 4-cycle: slow 2→4→16→37→58→89→145→**42**, fast 4→37→89→42→4→37→89→**42** — met, and 42 ≠ 1 → `false`.)

---

## Approach 3 — Hardcoded Cycle (Math Fact)

### Intuition

The finite-zone argument doesn't just say "some cycle exists" — the zone is small enough to check exhaustively, once, offline. Doing so reveals a remarkable fact: **the only cycle other than the fixed point 1 is** `4 → 16 → 37 → 58 → 89 → 145 → 42 → 20 → 4`. Therefore every unhappy trajectory must pass through 4, and every happy one ends at 1. At runtime, iterate the step function until the value is `1` or `4` — no set, no second pointer, just two sentinel comparisons.

### Algorithm

1. While `n != 1` and `n != 4`: replace `n` with `digitSquareSum(n)`.
2. Return `n == 1` (`n == 4` means the trajectory just entered the unique unhappy cycle).

### Complexity

- **Time:** O(log n) — bounded number of steps to fall into `{1} ∪ {4-cycle}`, each step costing O(log n) digit work.
- **Space:** O(1) — one integer; even less bookkeeping than Floyd.

### Code

```go
func hardcodedCycle(n int) bool {
	// 4 is the sentinel: every unhappy trajectory provably passes through it.
	for n != 1 && n != 4 {
		n = digitSquareSum(n)
	}
	return n == 1
}
```

### Dry Run

Example 1: `n = 19`.

| Step | n | n == 1? | n == 4? | digitSquareSum |
|------|-----|---------|---------|----------------|
| 1 | 19 | no | no | 82 |
| 2 | 82 | no | no | 68 |
| 3 | 68 | no | no | 100 |
| 4 | 100 | no | no | 1 |
| 5 | 1 | yes | — | exit loop |

Return `true` ✔ (For `n = 2`: 2 → **4** — the sentinel fires after one step → `false`.)

---

## Key Takeaways

- **"Iterated function + does it loop?" = cycle detection.** Any deterministic `x → f(x)` process is an implicit linked list; the hash-set and Floyd tools from #141/#142 transfer verbatim. Recognising this reframing is the whole problem.
- **Bound the state space first.** The digit-square map sends every d-digit number to ≤ 81·d, so trajectories are trapped below 1000 — that single inequality guarantees termination of all three loops.
- **Floyd's tortoise & hare** removes the O(chain) memory whenever you can afford to re-run the step function; the meeting point inside the cycle needs no bookkeeping.
- The **4-cycle fact** (`4→16→37→58→89→145→42→20→4` is the only non-trivial cycle) is worth memorising: it turns the problem into a two-sentinel loop and is a favourite interview follow-up.
- Digit extraction idiom: `d := n % 10; n /= 10` — O(number of digits) = O(log₁₀ n), the same skeleton as #7 Reverse Integer and #9 Palindrome Number.

---

## Related Problems

- LeetCode #141 — Linked List Cycle (the same Floyd detection on a real list)
- LeetCode #142 — Linked List Cycle II (locate the cycle entry point)
- LeetCode #258 — Add Digits (iterated digit map with a closed-form answer)
- LeetCode #263 — Ugly Number (loop on a number-theoretic property)
- LeetCode #1945 — Sum of Digits of String After Convert (iterated digit-sum process)
