# 0365 — Water and Jug Problem

> LeetCode #365 · Difficulty: Medium
> **Categories:** Math, Number Theory, Depth-First Search, Breadth-First Search

---

## Problem Statement

You are given two jugs with capacities `x` liters and `y` liters. You have an infinite water supply. Return whether the total amount of exactly `target` liters can be measured using the two jugs by performing the following operations:

- Fill either jug completely with water.
- Completely empty either jug.
- Pour water from one jug into another until the receiving jug is full, or the transferring jug is empty.

**Example 1:**

```
Input: x = 3, y = 5, target = 4
Output: true
Explanation:
Follow these steps to reach a total of 4 liters:
1. Fill the 5-liter jug (0, 5).
2. Pour from the 5-liter jug into the 3-liter jug, leaving 2 liters (3, 2).
3. Empty the 3-liter jug (0, 2).
4. Transfer the 2 liters from the 5-liter jug to the 3-liter jug (2, 0).
5. Fill the 5-liter jug again (2, 5).
6. Pour from the 5-liter jug into the 3-liter jug until full. The 5-liter jug is left with 4 liters (3, 4).
7. Empty the 3-liter jug. Now, you have exactly 4 liters in the 5-liter jug (0, 4).
Reference: The Die Hard example.
```

**Example 2:**

```
Input: x = 2, y = 6, target = 5
Output: false
```

**Example 3:**

```
Input: x = 1, y = 2, target = 3
Output: true
```

**Constraints:**

- `1 <= x, y, target <= 10^3`

---

## Company Frequency

| Company    | Frequency        | Last Reported |
|------------|------------------|---------------|
| Amazon     | ★★★☆☆ Medium     | 2023          |
| Microsoft  | ★★★☆☆ Medium     | 2023          |
| Google     | ★★☆☆☆ Low        | 2022          |
| Bloomberg  | ★★☆☆☆ Low        | 2022          |
| Apple      | ★★☆☆☆ Low        | 2021          |

> ⚠️ Frequency data is crowd-sourced from LeetCode Discuss, Glassdoor, and
> community interview reports. Treat as a signal, not a guarantee.

---

## DSA Concepts Used

- **Number Theory (Bézout's identity, GCD)** — the target is measurable iff it is a multiple of `gcd(x,y)` and `≤ x+y`; Euclid's algorithm computes the gcd → see [`/dsa/math_number_theory.md`](/dsa/math_number_theory.md)
- **Breadth-First Search over a state graph** — each `(a,b)` jug configuration is a node; the six operations are edges → see [`/dsa/graph_bfs_dfs.md`](/dsa/graph_bfs_dfs.md)
- **Hashing / visited set** — states are deduplicated with a hash set so BFS terminates → see [`/dsa/hash_map.md`](/dsa/hash_map.md)

---

## Approaches Overview

| # | Approach | Time | Space | When to use |
|---|----------|------|-------|-------------|
| 1 | BFS over Reachable States | O(x·y) | O(x·y) | Intuitive; also recovers the sequence of moves if needed |
| 2 | Bézout / GCD (Optimal) | O(log·min(x,y)) | O(1) | The intended answer; a two-line number-theory check |

---

## Approach 1 — BFS over Reachable States

### Intuition

A state is the pair `(a,b)` of how much water each jug holds. From any state there are six legal moves: fill X, fill Y, empty X, empty Y, pour X→Y, pour Y→X. This defines a graph; BFS explores all reachable states until one has combined water equal to `target`. The state space is bounded by `(x+1)(y+1)`, so the search terminates.

### Algorithm

1. Start at `(0,0)`; mark visited. Handle `target == 0` (true) and `target > x+y` (false) up front.
2. Dequeue `(a,b)`; if `a+b == target`, return true.
3. Generate the six successor states, enqueue unvisited ones.
4. If the queue drains without success, return false.

### Complexity

- **Time:** O(x·y) — at most `(x+1)(y+1)` states, each expanding O(1) moves.
- **Space:** O(x·y) — visited set plus queue.

### Code

```go
func bfsStates(x, y, target int) bool {
	type state struct{ a, b int }
	start := state{0, 0}
	if target == 0 {
		return true // both jugs empty already measures 0
	}
	if target > x+y {
		return false // can never hold more than both jugs combined
	}
	visited := map[state]bool{start: true}
	queue := []state{start}
	for len(queue) > 0 {
		cur := queue[0]
		queue = queue[1:]
		a, b := cur.a, cur.b
		if a+b == target { // combined water in the two jugs equals target
			return true
		}
		// Six legal moves from (a,b):
		nexts := []state{
			{x, b},                     // fill jug X
			{a, y},                     // fill jug Y
			{0, b},                     // empty jug X
			{a, 0},                     // empty jug Y
			{a - min(a, y-b), b + min(a, y-b)}, // pour X → Y
			{a + min(b, x-a), b - min(b, x-a)}, // pour Y → X
		}
		for _, nx := range nexts {
			if !visited[nx] {
				visited[nx] = true
				queue = append(queue, nx)
			}
		}
	}
	return false
}
```

### Dry Run

Example 1: `x=3, y=5, target=4`. A shortest path BFS discovers (states as `(X,Y)`):

| Step | state `(a,b)` | a+b | == 4? | notable successor generated |
|------|---------------|-----|-------|-----------------------------|
| 0 | (0,0) | 0 | no | fill Y → (0,5) |
| … | (0,5) | 5 | no | pour Y→X → (3,2) |
| … | (3,2) | 5 | no | empty X → (0,2) |
| … | (0,2) | 2 | no | pour Y→X → (2,0), fill Y → (2,5) |
| … | (2,5) | 7 | no | pour Y→X → (3,4) |
| … | (3,4) | 7 | no | empty X → (0,4) |
| ✔ | (0,4) | 4 | **yes** | return true |

Reaches combined water 4 → `true` ✔

---

## Approach 2 — Bézout / GCD (Optimal)

### Intuition

Every operation changes the total water by a multiple of `x` or `y` (a fill adds a jug's capacity, an empty subtracts it, a pour just moves water between jugs). So every reachable total is `a·x + b·y` for integers `a,b`. Bézout's identity: the set of all such integer combinations is exactly the multiples of `gcd(x,y)`. Therefore `target` is achievable iff it is a multiple of `gcd(x,y)` — bounded by the physical cap `target ≤ x+y`.

### Algorithm

1. If `target == 0`, return true.
2. If `target > x+y`, return false.
3. Return `target % gcd(x,y) == 0`.

### Complexity

- **Time:** O(log·min(x,y)) — Euclid's algorithm for the gcd.
- **Space:** O(1) — iterative gcd, constant scalars.

### Code

```go
func gcdBezout(x, y, target int) bool {
	if target == 0 {
		return true // measuring nothing is always possible
	}
	if target > x+y {
		return false // cannot exceed the combined capacity of both jugs
	}
	// Bézout: a*x + b*y spans exactly the multiples of gcd(x,y).
	return target%gcd(x, y) == 0
}

func gcd(a, b int) int {
	for b != 0 {
		a, b = b, a%b // Euclid: replace (a,b) with (b, a mod b)
	}
	return a
}
```

### Dry Run

Example 1: `x=3, y=5, target=4`.

| Check | value |
|-------|-------|
| target == 0? | no |
| target > x+y (4 > 8)? | no |
| gcd(3,5): (3,5)→(5,3)→(3,2)→(2,1)→(1,0) | **1** |
| target % gcd == 4 % 1 | 0 → divisible |

Returns `true` ✔

Cross-check Example 2 `x=2,y=6,target=5`: `gcd(2,6)=2`, `5 % 2 = 1 ≠ 0` → `false` ✔.
Example 3 `x=1,y=2,target=3`: `3 ≤ 1+2` and `gcd(1,2)=1`, `3 % 1 = 0` → `true` ✔.

---

## Key Takeaways

- **Reachable totals from jug operations = integer combinations `a·x + b·y` = multiples of `gcd(x,y)`** (Bézout). This collapses a search problem into a one-line divisibility test.
- **Always apply the physical bound `target ≤ x+y`** — divisibility alone allows impossible totals larger than both jugs combined.
- **Euclid's gcd** `gcd(a,b) = gcd(b, a mod b)` is the reusable primitive; `gcd(n,0)=n` handles a zero-capacity jug gracefully.
- When a puzzle has a small bounded state space, **BFS is a safe fallback** that also yields the actual move sequence — but recognizing the number-theoretic invariant gives an exponentially faster answer.

---

## Related Problems

- LeetCode #914 — X of a Kind in a Deck of Cards (gcd over counts)
- LeetCode #1071 — Greatest Common Divisor of Strings (gcd structure)
- LeetCode #372 — Super Pow (number-theory / modular arithmetic)
- LeetCode #1201 — Ugly Number III (gcd/lcm inclusion-exclusion)
