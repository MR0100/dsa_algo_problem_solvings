package main

import "fmt"

// ── Approach 1: Brute Force ───────────────────────────────────────────────────
//
// bruteForce solves Largest Rectangle in Histogram by trying every pair (i, j)
// as the left and right boundary of a rectangle.
//
// Intuition:
//   For each pair of bars (i, j), the rectangle extending from i to j has
//   height = min(heights[i..j]). Track the maximum such area.
//
// Time:  O(n²)
// Space: O(1)
func bruteForce(heights []int) int {
	n := len(heights)
	maxArea := 0
	for i := 0; i < n; i++ {
		minH := heights[i]
		for j := i; j < n; j++ {
			if heights[j] < minH {
				minH = heights[j]
			}
			area := minH * (j - i + 1)
			if area > maxArea {
				maxArea = area
			}
		}
	}
	return maxArea
}

// ── Approach 2: Stack (Monotonic Stack, Optimal) ──────────────────────────────
//
// monoStack solves Largest Rectangle in Histogram using a monotonic
// increasing stack.
//
// Intuition:
//   For each bar, find the first bar to its left that is shorter (left boundary)
//   and the first bar to its right that is shorter (right boundary). The
//   rectangle using this bar as the shortest bar spans from left to right.
//
//   A monotonic increasing stack maintains bars in order of increasing height.
//   When we encounter a bar shorter than the top of the stack, the top bar's
//   right boundary has been found. Its left boundary is the new top of the stack.
//
//   Sentinel: append a 0 at the end to flush all remaining bars from the stack.
//
// Algorithm:
//   stack = [] (stores indices); append 0 to heights
//   for i, h in heights:
//     while stack not empty && heights[stack.top] > h:
//       height = heights[stack.pop]
//       width = i if stack empty else i - stack.top - 1
//       maxArea = max(maxArea, height * width)
//     stack.push(i)
//
// Time:  O(n) — each bar pushed and popped at most once.
// Space: O(n) — stack.
func monoStack(heights []int) int {
	// append sentinel 0 to flush all remaining elements
	heights = append(heights, 0)
	stack := []int{} // indices into heights
	maxArea := 0

	for i, h := range heights {
		for len(stack) > 0 && heights[stack[len(stack)-1]] > h {
			// pop the top; this is the shortest bar in a span
			topIdx := stack[len(stack)-1]
			stack = stack[:len(stack)-1]
			height := heights[topIdx]

			var width int
			if len(stack) == 0 {
				width = i // extends all the way to the left
			} else {
				width = i - stack[len(stack)-1] - 1
			}
			area := height * width
			if area > maxArea {
				maxArea = area
			}
		}
		stack = append(stack, i)
	}
	return maxArea
}

func main() {
	fmt.Println("=== Approach 1: Brute Force ===")
	fmt.Printf("heights=[2,1,5,6,2,3]  got=%d  expected 10\n", bruteForce([]int{2, 1, 5, 6, 2, 3}))
	fmt.Printf("heights=[2,4]  got=%d  expected 4\n", bruteForce([]int{2, 4}))
	fmt.Printf("heights=[1]  got=%d  expected 1\n", bruteForce([]int{1}))

	fmt.Println("=== Approach 2: Monotonic Stack ===")
	fmt.Printf("heights=[2,1,5,6,2,3]  got=%d  expected 10\n", monoStack([]int{2, 1, 5, 6, 2, 3}))
	fmt.Printf("heights=[2,4]  got=%d  expected 4\n", monoStack([]int{2, 4}))
	fmt.Printf("heights=[1]  got=%d  expected 1\n", monoStack([]int{1}))
	fmt.Printf("heights=[6,2,5,4,5,1,6]  got=%d  expected 12\n", monoStack([]int{6, 2, 5, 4, 5, 1, 6}))
}
