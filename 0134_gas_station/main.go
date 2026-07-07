package main

import "fmt"

// ── Approach 1: Brute Force ──────────────────────────────────────────────────
//
// bruteForce solves Gas Station by simulating the full circuit from every
// possible starting station.
//
// Intuition: there are only n candidate starts. For each one, drive around
// the circle station by station, adding gas and paying cost; if the tank
// ever goes negative the start fails. The first start that completes the
// loop is the answer (the problem guarantees uniqueness if one exists).
//
// Algorithm:
//  1. For each start in [0, n-1]:
//     a. tank = 0.
//     b. For step = 0..n-1: i = (start+step) % n; tank += gas[i] - cost[i];
//     if tank < 0, abandon this start.
//     c. If all n steps succeed, return start.
//  2. Return -1 if every start fails.
//
// Time:  O(n²) — n starts × n-step simulation each.
// Space: O(1) — only the tank counter.
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

// ── Approach 2: Prefix-Sum Minimum ───────────────────────────────────────────
//
// prefixMinimum solves Gas Station by locating the global minimum of the
// running surplus and starting right after it.
//
// Intuition: let diff[i] = gas[i] - cost[i] and walk the circle once from
// station 0, tracking the running sum. The point where the running surplus
// is at its lowest is the "deepest valley" of the trip — every station up to
// and including it drains the tank. Starting at the station just AFTER the
// valley means the valley's deficit is paid last, when the accumulated
// surplus is largest, so the tank never dips below zero. Feasibility overall
// only depends on total(gas) >= total(cost).
//
// Algorithm:
//  1. Scan i = 0..n-1 accumulating total += gas[i]-cost[i].
//  2. Track minSum (lowest prefix value) and minIndex (where it occurs).
//  3. If total < 0, return -1; else return (minIndex + 1) % n.
//
// Time:  O(n) — single pass.
// Space: O(1) — three scalars.
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

// ── Approach 3: Greedy One-Pass (Optimal) ────────────────────────────────────
//
// greedyApproach solves Gas Station in one pass by discarding any start whose
// running tank goes negative and restarting from the next station.
//
// Intuition: two facts make greedy correct.
//  1. If the total surplus of the whole circle is negative, no start works.
//  2. If starting at s the tank first goes negative when leaving station i,
//     then NO station in (s..i] can be a valid start either: every such
//     station was reached with tank >= 0, so starting there (tank = 0, i.e.
//     with less or equal fuel) fails at i too. The next candidate is i+1.
//
// So one scan suffices: each failure jumps the candidate past the failure
// point, and each station is examined exactly once.
//
// Algorithm:
//  1. total = 0, tank = 0, start = 0.
//  2. For i = 0..n-1: d = gas[i]-cost[i]; total += d; tank += d.
//     If tank < 0: start = i+1; tank = 0.
//  3. Return start if total >= 0, else -1.
//
// Time:  O(n) — single pass, constant work per station.
// Space: O(1) — three scalars.
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

func main() {
	// Example 1
	gas1, cost1 := []int{1, 2, 3, 4, 5}, []int{3, 4, 5, 1, 2}
	// Example 2
	gas2, cost2 := []int{2, 3, 4}, []int{3, 4, 3}

	fmt.Println("=== Approach 1: Brute Force ===")
	fmt.Println(bruteForce(gas1, cost1)) // 3
	fmt.Println(bruteForce(gas2, cost2)) // -1

	fmt.Println("=== Approach 2: Prefix-Sum Minimum ===")
	fmt.Println(prefixMinimum(gas1, cost1)) // 3
	fmt.Println(prefixMinimum(gas2, cost2)) // -1

	fmt.Println("=== Approach 3: Greedy One-Pass (Optimal) ===")
	fmt.Println(greedyApproach(gas1, cost1)) // 3
	fmt.Println(greedyApproach(gas2, cost2)) // -1
}
