package main

import "fmt"

// ── Approach 1: BFS over Reachable States ────────────────────────────────────
//
// bfsStates solves the Water and Jug Problem by exploring every reachable
// (amountX, amountY) state via the six legal operations until some state has
// total water equal to target.
//
// Intuition:
//
//	A "state" is how much water each jug currently holds. From any state we
//	can: fill either jug, empty either jug, or pour one jug into the other
//	until the source is empty or the destination is full — six moves. Search
//	the state graph (BFS) for a state whose combined water equals target.
//	The state space is bounded by (x+1)*(y+1), so it terminates.
//
// Algorithm:
//  1. Start from (0,0); mark visited.
//  2. From (a,b) generate the six successor states.
//  3. If any state has a+b == target (or a==target / b==target), return true.
//  4. If the queue empties without success, return false.
//
// Time:  O(x*y) — at most (x+1)*(y+1) distinct states, each with O(1) moves.
// Space: O(x*y) — the visited set and BFS queue.
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
			{x, b},                             // fill jug X
			{a, y},                             // fill jug Y
			{0, b},                             // empty jug X
			{a, 0},                             // empty jug Y
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

// min returns the smaller of two ints.
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// ── Approach 2: Bézout / GCD Number Theory (Optimal) ─────────────────────────
//
// gcdBezout solves the Water and Jug Problem in O(log(min(x,y))) using number
// theory: target is measurable iff target ≤ x+y and target is a multiple of
// gcd(x, y).
//
// Intuition:
//
//	Every operation changes the total water by a multiple of x or y (fill or
//	empty a jug, or pour, which is conservative). So any reachable total is of
//	the form a*x + b*y for integers a,b. Bézout's identity says the set of all
//	such integer combinations is exactly the multiples of gcd(x,y). Hence the
//	target is achievable exactly when it is a multiple of gcd(x,y), subject to
//	the physical cap target ≤ x+y (you cannot hold more than both jugs).
//
// Algorithm:
//  1. If target == 0, return true (trivially measured).
//  2. If target > x+y, return false (exceeds total capacity).
//  3. Return target % gcd(x, y) == 0.
//
// Time:  O(log(min(x,y))) — Euclid's algorithm.
// Space: O(1) — iterative gcd (constant scalars).
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

// gcd computes the greatest common divisor via the iterative Euclid algorithm.
// gcd(n, 0) = n handles the case where one jug has capacity 0.
func gcd(a, b int) int {
	for b != 0 {
		a, b = b, a%b // Euclid: replace (a,b) with (b, a mod b)
	}
	return a
}

func main() {
	// Example 1: x=3, y=5, target=4 → true.
	// Example 2: x=2, y=6, target=5 → false (gcd 2 does not divide 5).
	// Example 3: x=1, y=2, target=3 → true (3 = 1+2, gcd 1 divides 3).

	fmt.Println("=== Approach 1: BFS over Reachable States ===")
	fmt.Println(bfsStates(3, 5, 4)) // expected true
	fmt.Println(bfsStates(2, 6, 5)) // expected false
	fmt.Println(bfsStates(1, 2, 3)) // expected true

	fmt.Println("=== Approach 2: Bézout / GCD (Optimal) ===")
	fmt.Println(gcdBezout(3, 5, 4)) // expected true
	fmt.Println(gcdBezout(2, 6, 5)) // expected false
	fmt.Println(gcdBezout(1, 2, 3)) // expected true
}
