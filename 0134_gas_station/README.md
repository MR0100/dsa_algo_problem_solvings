# 0134 — Gas Station

> LeetCode #134 · Difficulty: Medium
> **Categories:** Array, Greedy

---

## Problem Statement

There are `n` gas stations along a circular route, where the amount of gas at the `iᵗʰ` station is `gas[i]`.

You have a car with an unlimited gas tank and it costs `cost[i]` of gas to travel from the `iᵗʰ` station to its next `(i + 1)ᵗʰ` station. You begin the journey with an empty tank at one of the gas stations.

Given two integer arrays `gas` and `cost`, return *the starting gas station's index if you can travel around the circuit once in the clockwise direction, otherwise return* `-1`. If there exists a solution, it is **guaranteed** to be **unique**.

**Example 1:**
```
Input: gas = [1,2,3,4,5], cost = [3,4,5,1,2]
Output: 3
Explanation:
Start at station 3 (index 3) and fill up with 4 unit of gas. Your tank = 0 + 4 = 4
Travel to station 4. Your tank = 4 - 1 + 5 = 8
Travel to station 0. Your tank = 8 - 2 + 1 = 7
Travel to station 1. Your tank = 7 - 3 + 2 = 6
Travel to station 2. Your tank = 6 - 4 + 3 = 5
Travel to station 3. The cost is 5. Your gas is just enough to travel back to station 3.
Therefore, return 3 as the starting index.
```

**Example 2:**
```
Input: gas = [2,3,4], cost = [3,4,3]
Output: -1
Explanation:
You can't start at station 0 or 1, as there is not enough gas to travel to the next station.
Let's start at station 2 and fill up with 4 unit of gas. Your tank = 0 + 4 = 4
Travel to station 0. Your tank = 4 - 3 + 2 = 3
Travel to station 1. Your tank = 3 - 3 + 3 = 3
You cannot travel back to station 2, as it requires 4 unit of gas but you only have 3.
Therefore, you can't travel around the circuit once no matter where you start.
```

**Constraints:**
- `n == gas.length == cost.length`
- `1 <= n <= 10⁵`
- `0 <= gas[i], cost[i] <= 10⁴`

---

## Company Frequency

| Company   | Frequency        | Last Reported |
|-----------|------------------|---------------|
| Amazon    | ★★★★★ Very High  | 2024          |
| Microsoft | ★★★★☆ High       | 2024          |
| Google    | ★★★☆☆ Medium     | 2024          |
| Bloomberg | ★★★☆☆ Medium     | 2023          |
| Apple     | ★★☆☆☆ Low        | 2023          |

> ⚠️ Frequency data is crowd-sourced from LeetCode Discuss, Glassdoor, and
> community interview reports. Treat as a signal, not a guarantee.

---

## DSA Concepts Used

- **Greedy** — a local failure ("tank went negative") eliminates a whole range of candidate starts at once, so one pass suffices → see [`/dsa/greedy.md`](/dsa/greedy.md)
- **Prefix Sum** — the running surplus of `gas[i] - cost[i]` is a prefix sum over a circular array; its global minimum pinpoints the answer → see [`/dsa/prefix_sum.md`](/dsa/prefix_sum.md)

---

## Approaches Overview

| # | Approach | Time | Space | When to use |
|---|----------|------|-------|-------------|
| 1 | Brute Force | O(n²) | O(1) | Baseline; too slow for n = 10⁵ |
| 2 | Prefix-Sum Minimum | O(n) | O(1) | Elegant "deepest valley" argument; easy to prove |
| 3 | Greedy One-Pass (Optimal) | O(n) | O(1) | Standard interview answer; single scan, one theorem |

---

## Approach 1 — Brute Force

### Intuition
There are only n possible starting stations, so try them all. From each candidate start, simulate the drive: at every station add its gas, then pay the cost to reach the next one. If the tank ever dips below zero, that start fails. Since the problem guarantees a unique answer when one exists, the first start that survives all n hops is the answer.

### Algorithm
1. For each `start` in `[0, n-1]`:
   1. Set `tank = 0`.
   2. For `step` from 0 to n−1: let `i = (start + step) % n`; update `tank += gas[i] - cost[i]`.
   3. If `tank < 0` at any point, abandon this start and try the next.
   4. If all n hops succeed, return `start`.
2. Return `-1`.

### Complexity
- **Time:** O(n²) — n candidate starts, each simulated for up to n hops (≈10¹⁰ operations at n = 10⁵ — TLE on LeetCode).
- **Space:** O(1) — a single tank counter.

### Code
```go
func bruteForce(gas []int, cost []int) int {
	n := len(gas)
	for start := 0; start < n; start++ {
		tank := 0 // fuel in the tank when leaving each station
		ok := true
		for step := 0; step < n; step++ {
			i := (start + step) % n  // wrap around the circular route
			tank += gas[i] - cost[i] // fill up at i, pay to reach i+1
			if tank < 0 {            // ran dry before the next station
				ok = false
				break // this start is infeasible; try the next one
			}
		}
		if ok {
			return start // completed the full circle
		}
	}
	return -1 // no starting station works
}
```

### Dry Run — Example 1: `gas = [1,2,3,4,5]`, `cost = [3,4,5,1,2]`

(`diff = gas - cost = [-2, -2, -2, 3, 3]`)

| start | Hops simulated (tank after each) | Outcome |
|-------|----------------------------------|---------|
| 0 | i=0: 0 + (−2) = **−2** | fail immediately |
| 1 | i=1: 0 + (−2) = **−2** | fail immediately |
| 2 | i=2: 0 + (−2) = **−2** | fail immediately |
| 3 | i=3: 3 → i=4: 6 → i=0: 4 → i=1: 2 → i=2: 0 | **all 5 hops ok → return 3** |

Final answer: `3`. ✅

---

## Approach 2 — Prefix-Sum Minimum

### Intuition
Define `diff[i] = gas[i] - cost[i]` and walk the circle once from station 0, tracking the running sum (a prefix sum). Two observations:

1. **Feasibility** depends only on the total: if `Σ diff < 0` the circle consumes more than it provides — impossible from anywhere. If `Σ diff ≥ 0`, a start must exist.
2. The running sum traces a "elevation profile" of the trip. Its **global minimum** is the deepest valley — the point where the journey so far has been most fuel-starved. Starting **just after the valley** reorders the trip so that the draining stretch is traversed *last*, when the accumulated surplus from the rest of the circle is at its largest. The tank therefore never goes negative.

### Algorithm
1. Initialize `total = 0`, `minSum = 0`, `minIndex = -1`.
2. For `i` from 0 to n−1: `total += gas[i] - cost[i]`; if `total < minSum`, set `minSum = total`, `minIndex = i`.
3. If `total < 0`, return `-1`.
4. Return `(minIndex + 1) % n` (the station right after the deepest valley; wraps to 0 when the minimum is at the last station).

### Complexity
- **Time:** O(n) — one pass over the stations.
- **Space:** O(1) — three scalar accumulators.

### Code
```go
func prefixMinimum(gas []int, cost []int) int {
	n := len(gas)
	total := 0     // running sum of gas[i]-cost[i] over the whole circle
	minSum := 0    // lowest prefix sum seen so far
	minIndex := -1 // index at which the lowest prefix sum occurs

	for i := 0; i < n; i++ {
		total += gas[i] - cost[i] // surplus after leaving station i
		if total < minSum {
			minSum = total // deeper valley found
			minIndex = i   // valley bottoms out after station i
		}
	}

	if total < 0 {
		return -1 // circle consumes more than it provides: impossible
	}
	// start just past the deepest valley; wraps to 0 when minIndex == n-1
	return (minIndex + 1) % n
}
```

### Dry Run — Example 1: `gas = [1,2,3,4,5]`, `cost = [3,4,5,1,2]`

| i | `diff[i]` | `total` (prefix sum) | `total < minSum`? | `minSum` | `minIndex` |
|---|-----------|----------------------|-------------------|----------|------------|
| 0 | −2 | −2 | yes (−2 < 0) | −2 | 0 |
| 1 | −2 | −4 | yes | −4 | 1 |
| 2 | −2 | −6 | yes | −6 | 2 |
| 3 | +3 | −3 | no | −6 | 2 |
| 4 | +3 | 0 | no | −6 | 2 |

`total = 0 ≥ 0` → answer `(minIndex + 1) % n = (2 + 1) % 5 = 3`. ✅
(The valley bottoms out after station 2; starting at 3 puts the −2,−2,−2 stretch last, cushioned by +3,+3.)

---

## Approach 3 — Greedy One-Pass (Optimal)

### Intuition
Two facts make a single greedy scan correct:

1. **Global feasibility:** if `Σ (gas[i] − cost[i]) < 0`, no start can work (total demand exceeds total supply). Conversely if the total is ≥ 0, a valid start is guaranteed to exist.
2. **Failure elimination:** suppose we start at `s` and the tank first goes negative after taking station `i`'s hop. Then **no station in `(s, i]` can be a valid start either.** Why: for any `m` in that range, the drive from `s` arrived at `m` with `tank ≥ 0`; starting at `m` instead means arriving at `i` with *less or equal* fuel (it starts from 0 instead of a non-negative carry-over), so it fails at `i` too. Therefore the next candidate worth trying is `i + 1` — we can jump the start forward and never re-examine a station.

Each station is visited exactly once, giving O(n).

### Algorithm
1. Initialize `total = 0`, `tank = 0`, `start = 0`.
2. For `i` from 0 to n−1:
   1. `d = gas[i] - cost[i]`; add `d` to both `total` and `tank`.
   2. If `tank < 0`: set `start = i + 1` and reset `tank = 0` (every start in `(old start, i]` is eliminated by fact 2).
3. If `total < 0`, return `-1`; otherwise return `start` (unique by problem statement; fact 1 guarantees it completes the circle).

### Complexity
- **Time:** O(n) — one pass; each index processed once, all updates O(1).
- **Space:** O(1) — three scalars.

### Code
```go
func greedyApproach(gas []int, cost []int) int {
	total := 0 // net surplus over the entire circle (feasibility test)
	tank := 0  // fuel since the current candidate start
	start := 0 // current candidate starting station

	for i := 0; i < len(gas); i++ {
		d := gas[i] - cost[i] // net fuel gained by visiting station i
		total += d
		tank += d
		if tank < 0 {
			// candidate start (and everything between it and i) is doomed:
			// restart the attempt from the next station with an empty tank
			start = i + 1
			tank = 0
		}
	}

	if total < 0 {
		return -1 // whole circle is a net loss: no start can work
	}
	return start // guaranteed unique valid start
}
```

### Dry Run — Example 1: `gas = [1,2,3,4,5]`, `cost = [3,4,5,1,2]`

| i | `d = gas[i]−cost[i]` | `total` | `tank` after `+= d` | `tank < 0`? | `start` | `tank` after reset |
|---|----------------------|---------|---------------------|-------------|---------|--------------------|
| 0 | −2 | −2 | −2 | yes | 1 | 0 |
| 1 | −2 | −4 | −2 | yes | 2 | 0 |
| 2 | −2 | −6 | −2 | yes | 3 | 0 |
| 3 | +3 | −3 | 3 | no | 3 | 3 |
| 4 | +3 | 0 | 6 | no | 3 | 6 |

End of loop: `total = 0 ≥ 0` → return `start = 3`. ✅

---

## Key Takeaways

- **`Σ gas ≥ Σ cost` ⟺ a valid start exists** — separate the *feasibility* question (a global sum) from the *location* question (a scan). Many circular-route problems decompose the same way.
- **Greedy elimination lemma:** if starting at `s` the tank first fails at `i`, every start in `(s, i]` also fails — because each was reached with non-negative fuel, so it begins with no advantage. Skipping the candidate to `i+1` turns O(n²) into O(n).
- **"Start after the global minimum of the prefix sums"** is the prefix-sum view of the same fact — the deepest valley of the surplus profile must be crossed last.
- Reducing two arrays to one (`diff[i] = gas[i] − cost[i]`) is a common first simplification step.
- Both O(n) solutions need **no wrap-around second pass** — the feasibility total plus the elimination argument covers the circularity.

---

## Related Problems

- LeetCode #135 — Candy (adjacent greedy constraint, two-directional reasoning)
- LeetCode #55 — Jump Game (greedy reachability with a running budget)
- LeetCode #45 — Jump Game II (greedy interval extension)
- LeetCode #53 — Maximum Subarray (Kadane's: same "reset when the running sum hurts" instinct)
- LeetCode #918 — Maximum Sum Circular Subarray (prefix sums on a circular array)
- LeetCode #871 — Minimum Number of Refueling Stops (fuel feasibility, heap-greedy)
