package main

import "fmt"

// ── Approach 1: Brute Force ───────────────────────────────────────────────────
//
// bruteForce solves Longest Valid Parentheses by checking every substring.
//
// Intuition: Try every possible start and end position; validate with a
// balance counter; track the maximum valid length.
//
// Algorithm:
//  1. For every pair (i, j) where j-i+1 is even: validate s[i..j].
//  2. Validation: balance starts at 0; increment on '(', decrement on ')'.
//     If balance < 0 at any point, invalid. Valid iff balance == 0 at end.
//  3. Track max valid length.
//
// Time:  O(n³) — O(n²) substrings × O(n) validation
// Space: O(1)
func bruteForce(s string) int {
	n := len(s)
	best := 0
	for i := 0; i < n; i++ {
		balance := 0
		for j := i; j < n; j++ {
			if s[j] == '(' {
				balance++
			} else {
				balance--
			}
			if balance < 0 {
				break // can never recover
			}
			if balance == 0 {
				length := j - i + 1
				if length > best {
					best = length
				}
			}
		}
	}
	return best
}

// ── Approach 2: Stack ─────────────────────────────────────────────────────────
//
// stackApproach solves Longest Valid Parentheses using a stack of indices.
//
// Intuition: Push indices onto a stack. The stack always holds the index of
// the last "unmatched" character. When we find a matching pair, pop the top;
// the length of the current valid substring is i - stack[top].
//
// Algorithm:
//  1. Push -1 onto the stack as a sentinel (the "base" before any valid string).
//  2. For each i:
//     if s[i]=='(': push i.
//     if s[i]==')':
//       pop the top.
//       if stack is empty: push i (new base).
//       else: current valid length = i - stack[top].
//  3. Return max length seen.
//
// Time:  O(n)
// Space: O(n)
func stackApproach(s string) int {
	stack := []int{-1} // sentinel: base before valid substrings
	best := 0
	for i := 0; i < len(s); i++ {
		if s[i] == '(' {
			stack = append(stack, i) // push opening bracket index
		} else {
			// pop the matching '(' (or the base sentinel)
			stack = stack[:len(stack)-1]
			if len(stack) == 0 {
				stack = append(stack, i) // this ')' is unmatched; becomes new base
			} else {
				length := i - stack[len(stack)-1] // distance from base to current index
				if length > best {
					best = length
				}
			}
		}
	}
	return best
}

// ── Approach 3: Dynamic Programming ──────────────────────────────────────────
//
// dpApproach solves Longest Valid Parentheses using a DP array.
//
// Intuition: dp[i] = length of the longest valid substring ending at index i.
//   - If s[i] == '(': dp[i] = 0 (no valid string ends with '(').
//   - If s[i] == ')':
//       let j = i - dp[i-1] - 1 (the index just before the current valid suffix).
//       if s[j] == '(': dp[i] = dp[i-1] + 2 + dp[j-1]
//         (the current pair + the valid suffix before the pair + any valid
//         substring that ended just before j).
//
// Time:  O(n)
// Space: O(n)
func dpApproach(s string) int {
	n := len(s)
	dp := make([]int, n) // dp[i] = longest valid substring ending at i
	best := 0
	for i := 1; i < n; i++ {
		if s[i] == ')' {
			j := i - dp[i-1] - 1 // index of the potential matching '('
			if j >= 0 && s[j] == '(' {
				dp[i] = dp[i-1] + 2 // match the pair
				if j > 0 {
					dp[i] += dp[j-1] // add valid substring before the '('
				}
			}
		}
		if dp[i] > best {
			best = dp[i]
		}
	}
	return best
}

// ── Approach 4: Two Passes (Left-Right Counters) — Optimal ───────────────────
//
// twoPass solves Longest Valid Parentheses in O(n) time and O(1) space.
//
// Intuition: Scan left→right keeping (open, close) counters.
//   - Increment open on '(', close on ')'.
//   - When open == close: record 2*close as a candidate.
//   - When close > open: reset both to 0 (this ')' can never be matched).
//   Then scan right→left with the same logic (swapping open/close roles) to
//   catch valid strings where open never equals close in the forward pass
//   (e.g., "((()" has valid suffix "()" but forward pass never resets).
//
// Time:  O(n)
// Space: O(1)
func twoPass(s string) int {
	best := 0
	// left → right pass
	open, close := 0, 0
	for _, ch := range s {
		if ch == '(' {
			open++
		} else {
			close++
		}
		if open == close {
			length := 2 * close
			if length > best {
				best = length
			}
		} else if close > open {
			open, close = 0, 0 // unrecoverable excess ')'
		}
	}
	// right → left pass
	open, close = 0, 0
	for i := len(s) - 1; i >= 0; i-- {
		if s[i] == '(' {
			open++
		} else {
			close++
		}
		if open == close {
			length := 2 * open
			if length > best {
				best = length
			}
		} else if open > close {
			open, close = 0, 0 // unrecoverable excess '('
		}
	}
	return best
}

func main() {
	cases := []struct {
		s    string
		want int
	}{
		{"(()", 2},
		{")()())", 4},
		{"", 0},
		{"()()", 4},
		{"(())", 4},
		{"()(()", 2},
	}

	fmt.Println("=== Approach 1: Brute Force ===")
	for _, c := range cases {
		fmt.Printf("s=%-12q  got=%d  expected=%d\n", c.s, bruteForce(c.s), c.want)
	}

	fmt.Println("\n=== Approach 2: Stack ===")
	for _, c := range cases {
		fmt.Printf("s=%-12q  got=%d  expected=%d\n", c.s, stackApproach(c.s), c.want)
	}

	fmt.Println("\n=== Approach 3: DP ===")
	for _, c := range cases {
		fmt.Printf("s=%-12q  got=%d  expected=%d\n", c.s, dpApproach(c.s), c.want)
	}

	fmt.Println("\n=== Approach 4: Two Pass (Optimal) ===")
	for _, c := range cases {
		fmt.Printf("s=%-12q  got=%d  expected=%d\n", c.s, twoPass(c.s), c.want)
	}
}
