package main

import "fmt"

// ── Approach 1: Brute Force (Repeated Sweeps) ────────────────────────────────
//
// bruteForce solves Candy by giving everyone one candy and repeatedly fixing
// violations until a full sweep changes nothing.
//
// Intuition: the two rules ("higher rating than left neighbor → more candy
// than left" and symmetrically for right) are local constraints. Start from
// the minimum possible allocation (all ones) and keep bumping any child that
// violates a constraint to the smallest value that fixes it (neighbor + 1).
// Each bump only ever raises values by the minimum necessary, so when no
// violations remain the allocation is both valid and minimal.
//
// Algorithm:
//  1. candies[i] = 1 for all i.
//  2. Loop: for every i, if ratings[i] > ratings[i-1] and
//     candies[i] <= candies[i-1], set candies[i] = candies[i-1]+1;
//     likewise against the right neighbor. Repeat while anything changed.
//  3. Return the sum.
//
// Time:  O(n²) — worst case ~n sweeps (long increasing slope), each O(n).
// Space: O(n) — the candies array.
func bruteForce(ratings []int) int {
	n := len(ratings)
	candies := make([]int, n)
	for i := range candies {
		candies[i] = 1 // everyone gets at least one candy
	}

	for changed := true; changed; {
		changed = false // assume this sweep is clean until proven otherwise
		for i := 0; i < n; i++ {
			// left rule: strictly higher rating than left neighbor
			if i > 0 && ratings[i] > ratings[i-1] && candies[i] <= candies[i-1] {
				candies[i] = candies[i-1] + 1 // minimal fix for the violation
				changed = true
			}
			// right rule: strictly higher rating than right neighbor
			if i < n-1 && ratings[i] > ratings[i+1] && candies[i] <= candies[i+1] {
				candies[i] = candies[i+1] + 1 // minimal fix for the violation
				changed = true
			}
		}
	}

	total := 0
	for _, c := range candies {
		total += c // sum the final minimal allocation
	}
	return total
}

// ── Approach 2: Two-Pass Arrays ──────────────────────────────────────────────
//
// twoPassArrays solves Candy with one left-to-right pass and one
// right-to-left pass, taking the max of the two requirements per child.
//
// Intuition: the constraint couples both neighbors, which is awkward — but it
// splits cleanly into two one-sided constraints. A left pass alone can
// satisfy "higher than left neighbor → more candy": walk rightward and set
// candies[i] = candies[i-1]+1 on every ascent. A right pass alone satisfies
// the mirrored rule. A child on a peak needs to satisfy BOTH sides, and
// taking max(left[i], right[i]) does exactly that without breaking either
// pass's guarantee (each pass's value is the minimum for its own side).
//
// Algorithm:
//  1. left[i]: scan i = 1..n-1, left[i] = left[i-1]+1 if ratings ascend,
//     else 1.
//  2. right[i]: scan i = n-2..0, right[i] = right[i+1]+1 if ratings ascend
//     leftward, else 1.
//  3. Sum max(left[i], right[i]).
//
// Time:  O(n) — three linear passes.
// Space: O(n) — the two requirement arrays.
func twoPassArrays(ratings []int) int {
	n := len(ratings)

	left := make([]int, n)  // candies needed looking only at left neighbors
	right := make([]int, n) // candies needed looking only at right neighbors
	for i := range left {
		left[i] = 1 // base: one candy each
		right[i] = 1
	}

	// left-to-right: enforce "ascent from the left means +1 candy"
	for i := 1; i < n; i++ {
		if ratings[i] > ratings[i-1] {
			left[i] = left[i-1] + 1 // strictly more than the left neighbor
		}
	}

	// right-to-left: enforce "ascent from the right means +1 candy"
	for i := n - 2; i >= 0; i-- {
		if ratings[i] > ratings[i+1] {
			right[i] = right[i+1] + 1 // strictly more than the right neighbor
		}
	}

	total := 0
	for i := 0; i < n; i++ {
		// each child must satisfy both one-sided requirements simultaneously
		if left[i] > right[i] {
			total += left[i]
		} else {
			total += right[i]
		}
	}
	return total
}

// ── Approach 3: Slope Counting (Optimal, O(1) space) ─────────────────────────
//
// slopeCounting solves Candy in one pass with O(1) extra space by tracking
// the lengths of the current ascending and descending runs.
//
// Intuition: in an optimal allocation the candy counts along an ascending
// run of ratings look like 1,2,3,...,and along a descending run they mirror
// to ...,3,2,1. So the total only depends on run LENGTHS, not on the actual
// values — no array needed. The only subtlety is the peak between an ascent
// of length `up` and a descent of length `down`: the peak child needs
// max(up, down)+1 candies. We optimistically give the peak up+1 during the
// ascent; while descending, once the descent grows strictly longer than the
// recorded peak height, each further step must also raise the peak by one —
// accounted lazily by adding 1 extra per step (via not subtracting).
//
// Algorithm:
//  1. total = 1 (first child), up = down = peak = 0.
//  2. For each i from 1: compare ratings[i] with ratings[i-1]:
//     - ascent: up++, down = 0, peak = up; total += up + 1.
//     - flat:   up = down = peak = 0; total += 1.
//     - descent: down++, up = 0; total += down + 1; if peak >= down,
//     total-- (the peak already covers this descent length).
//  3. Return total.
//
// Time:  O(n) — single pass, O(1) work per child.
// Space: O(1) — four integer counters, no arrays.
func slopeCounting(ratings []int) int {
	n := len(ratings)
	if n <= 1 {
		return n // 0 children → 0 candies; 1 child → 1 candy
	}

	total := 1 // the first child always gets 1 candy to start
	up, down, peak := 0, 0, 0

	for i := 1; i < n; i++ {
		switch {
		case ratings[i] > ratings[i-1]:
			up++      // ascending run grows
			down = 0  // any previous descent is over
			peak = up // remember the height of the (current) peak
			// child i sits `up` steps above the run's start → needs up+1
			total += up + 1
		case ratings[i] == ratings[i-1]:
			// equal ratings carry no constraint: reset all runs,
			// this child can drop back to a single candy
			up, down, peak = 0, 0, 0
			total += 1
		default: // ratings[i] < ratings[i-1]
			down++ // descending run grows
			up = 0 // any previous ascent is over
			// tentatively this child starts a mirrored 1,2,...,down chain:
			// every earlier child in the descent shifts up by one → +down,
			// plus 1 candy for this child itself
			total += down + 1
			if peak >= down {
				// the peak (given up+1 earlier) is still strictly higher
				// than the descent chain, so it need not grow: take back
				// the one candy we just over-counted for it
				total--
			}
			// if peak < down, the peak must rise with the chain: the +1
			// stays, effectively raising the peak by one this step
		}
	}

	return total
}

func main() {
	// Example 1
	r1 := []int{1, 0, 2}
	// Example 2
	r2 := []int{1, 2, 2}

	fmt.Println("=== Approach 1: Brute Force (Repeated Sweeps) ===")
	fmt.Println(bruteForce(r1)) // 5
	fmt.Println(bruteForce(r2)) // 4

	fmt.Println("=== Approach 2: Two-Pass Arrays ===")
	fmt.Println(twoPassArrays(r1)) // 5
	fmt.Println(twoPassArrays(r2)) // 4

	fmt.Println("=== Approach 3: Slope Counting (Optimal, O(1) space) ===")
	fmt.Println(slopeCounting(r1)) // 5
	fmt.Println(slopeCounting(r2)) // 4
}
