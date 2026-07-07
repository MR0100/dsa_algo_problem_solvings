package main

import "fmt"

// ── Approach 1: Brute Force ───────────────────────────────────────────────────
//
// bruteForce solves Trapping Rain Water by computing, for each position i,
// the minimum of the maximum height to its left and right.
//
// Intuition: Water trapped at column i = min(maxLeft[i], maxRight[i]) - height[i].
// Iterate over all positions, and for each compute max left and max right.
//
// Time:  O(n²) — O(n) positions × O(n) for max computation per position
// Space: O(1)
func bruteForce(height []int) int {
	n := len(height)
	total := 0
	for i := 0; i < n; i++ {
		maxL, maxR := 0, 0
		for l := 0; l <= i; l++ {
			if height[l] > maxL {
				maxL = height[l]
			}
		}
		for r := i; r < n; r++ {
			if height[r] > maxR {
				maxR = height[r]
			}
		}
		// water above this column = min(maxL, maxR) - height[i]
		minWall := maxL
		if maxR < minWall {
			minWall = maxR
		}
		total += minWall - height[i]
	}
	return total
}

// ── Approach 2: Precomputed Left/Right Max Arrays ─────────────────────────────
//
// precomputed solves Trapping Rain Water using two extra arrays.
//
// Intuition: Precompute maxLeft[i] and maxRight[i] in O(n) each, then compute
// the answer in a final O(n) pass.
//
// Time:  O(n)
// Space: O(n)
func precomputed(height []int) int {
	n := len(height)
	if n == 0 {
		return 0
	}
	maxL := make([]int, n) // maxL[i] = max height in height[0..i]
	maxR := make([]int, n) // maxR[i] = max height in height[i..n-1]

	maxL[0] = height[0]
	for i := 1; i < n; i++ {
		if height[i] > maxL[i-1] {
			maxL[i] = height[i]
		} else {
			maxL[i] = maxL[i-1]
		}
	}
	maxR[n-1] = height[n-1]
	for i := n - 2; i >= 0; i-- {
		if height[i] > maxR[i+1] {
			maxR[i] = height[i]
		} else {
			maxR[i] = maxR[i+1]
		}
	}

	total := 0
	for i := 0; i < n; i++ {
		minWall := maxL[i]
		if maxR[i] < minWall {
			minWall = maxR[i]
		}
		total += minWall - height[i]
	}
	return total
}

// ── Approach 3: Two Pointers (Optimal) ───────────────────────────────────────
//
// twoPointers solves Trapping Rain Water in O(n) time and O(1) space.
//
// Intuition: Maintain left and right pointers. Track maxLeft (max height seen
// from the left so far) and maxRight (max height seen from the right so far).
//
// Key insight: if maxLeft <= maxRight, the water level at `left` is determined
// by maxLeft (the shorter wall). We can safely compute the water at `left` and
// move left++. Symmetrically for the right side.
//
// Why this works: when maxLeft <= maxRight, even though we don't know the full
// right side, we know maxRight is at least maxLeft (already seen), so the water
// level at `left` is exactly maxLeft - height[left].
//
// Algorithm:
//  left=0, right=n-1, maxL=0, maxR=0, total=0
//  while left < right:
//    if height[left] <= height[right]:
//      if height[left] >= maxL: maxL = height[left]
//      else: total += maxL - height[left]
//      left++
//    else:
//      if height[right] >= maxR: maxR = height[right]
//      else: total += maxR - height[right]
//      right--
//
// Time:  O(n)
// Space: O(1)
func twoPointers(height []int) int {
	left, right := 0, len(height)-1
	maxL, maxR := 0, 0
	total := 0
	for left < right {
		if height[left] <= height[right] {
			if height[left] >= maxL {
				maxL = height[left] // new max on left side
			} else {
				total += maxL - height[left] // water fills up to maxL
			}
			left++
		} else {
			if height[right] >= maxR {
				maxR = height[right] // new max on right side
			} else {
				total += maxR - height[right] // water fills up to maxR
			}
			right--
		}
	}
	return total
}

// ── Approach 4: Stack ─────────────────────────────────────────────────────────
//
// stackApproach solves Trapping Rain Water using a monotonic decreasing stack.
//
// Intuition: Maintain a stack of indices in decreasing height order. When we
// encounter a bar taller than the stack top, a valley (water pocket) forms.
// Pop the valley bottom; the water width is (current - stack.top - 1) and
// height is min(height[current], height[stack.top]) - height[valley].
//
// Time:  O(n)
// Space: O(n)
func stackApproach(height []int) int {
	stack := []int{} // stack of indices
	total := 0
	for i := 0; i < len(height); i++ {
		for len(stack) > 0 && height[i] > height[stack[len(stack)-1]] {
			valley := stack[len(stack)-1]
			stack = stack[:len(stack)-1] // pop
			if len(stack) == 0 {
				break // no left wall
			}
			leftWall := stack[len(stack)-1]
			width := i - leftWall - 1
			wallH := height[leftWall]
			if height[i] < wallH {
				wallH = height[i]
			}
			total += (wallH - height[valley]) * width
		}
		stack = append(stack, i) // push current index
	}
	return total
}

func main() {
	cases := []struct {
		h    []int
		want int
	}{
		{[]int{0, 1, 0, 2, 1, 0, 1, 3, 2, 1, 2, 1}, 6},
		{[]int{4, 2, 0, 3, 2, 5}, 9},
		{[]int{3, 0, 2, 0, 4}, 7},
		{[]int{1, 0, 1}, 1},
	}

	fmt.Println("=== Approach 1: Brute Force ===")
	for _, c := range cases {
		fmt.Printf("height=%v => %d  expected %d\n", c.h, bruteForce(c.h), c.want)
	}

	fmt.Println("\n=== Approach 2: Precomputed Arrays ===")
	for _, c := range cases {
		fmt.Printf("height=%v => %d  expected %d\n", c.h, precomputed(c.h), c.want)
	}

	fmt.Println("\n=== Approach 3: Two Pointers (Optimal) ===")
	for _, c := range cases {
		fmt.Printf("height=%v => %d  expected %d\n", c.h, twoPointers(c.h), c.want)
	}

	fmt.Println("\n=== Approach 4: Stack ===")
	for _, c := range cases {
		fmt.Printf("height=%v => %d  expected %d\n", c.h, stackApproach(c.h), c.want)
	}
}
