package main

import "fmt"

// ── Approach 1: DFS from Border ───────────────────────────────────────────────
//
// solve solves Surrounded Regions by first marking all 'O's connected to the
// border (cannot be captured), then flipping the rest.
//
// Intuition:
//   An 'O' is NOT captured if it is connected to a border 'O'. So:
//   1. DFS from every border 'O', marking safe cells as 'S'.
//   2. Flip all remaining 'O' → 'X' (captured), 'S' → 'O' (restore safe).
//
// Time:  O(m*n)
// Space: O(m*n) — DFS recursion stack.
func solve(board [][]byte) {
	if len(board) == 0 {
		return
	}
	m, n := len(board), len(board[0])

	var dfs func(r, c int)
	dfs = func(r, c int) {
		if r < 0 || r >= m || c < 0 || c >= n || board[r][c] != 'O' {
			return
		}
		board[r][c] = 'S' // mark safe
		dfs(r+1, c); dfs(r-1, c); dfs(r, c+1); dfs(r, c-1)
	}

	// mark border-connected 'O's
	for c := 0; c < n; c++ {
		dfs(0, c); dfs(m-1, c)
	}
	for r := 0; r < m; r++ {
		dfs(r, 0); dfs(r, n-1)
	}

	// flip
	for r := 0; r < m; r++ {
		for c := 0; c < n; c++ {
			if board[r][c] == 'O' {
				board[r][c] = 'X' // captured
			} else if board[r][c] == 'S' {
				board[r][c] = 'O' // restore safe
			}
		}
	}
}

// ── Approach 2: BFS from Border ───────────────────────────────────────────────
//
// solveBFS solves Surrounded Regions using BFS instead of DFS.
//
// Time:  O(m*n)
// Space: O(m*n) — BFS queue.
func solveBFS(board [][]byte) {
	if len(board) == 0 {
		return
	}
	m, n := len(board), len(board[0])
	type pt struct{ r, c int }
	dirs := []pt{{1, 0}, {-1, 0}, {0, 1}, {0, -1}}

	var queue []pt
	enqueue := func(r, c int) {
		if r >= 0 && r < m && c >= 0 && c < n && board[r][c] == 'O' {
			board[r][c] = 'S'
			queue = append(queue, pt{r, c})
		}
	}

	// seed from border
	for c := 0; c < n; c++ { enqueue(0, c); enqueue(m-1, c) }
	for r := 0; r < m; r++ { enqueue(r, 0); enqueue(r, n-1) }

	for len(queue) > 0 {
		curr := queue[0]; queue = queue[1:]
		for _, d := range dirs {
			enqueue(curr.r+d.r, curr.c+d.c)
		}
	}

	for r := 0; r < m; r++ {
		for c := 0; c < n; c++ {
			if board[r][c] == 'O' { board[r][c] = 'X' } else if board[r][c] == 'S' { board[r][c] = 'O' }
		}
	}
}

func main() {
	b1 := [][]byte{
		{'X', 'X', 'X', 'X'},
		{'X', 'O', 'O', 'X'},
		{'X', 'X', 'O', 'X'},
		{'X', 'O', 'X', 'X'},
	}
	fmt.Println("=== Approach 1: DFS from Border ===")
	solve(b1)
	for _, row := range b1 {
		fmt.Println(string(row))
	}
	fmt.Println("expected:")
	fmt.Println("XXXX")
	fmt.Println("XXXX")
	fmt.Println("XXXX")
	fmt.Println("XOXX")

	b2 := [][]byte{
		{'X', 'X', 'X', 'X'},
		{'X', 'O', 'O', 'X'},
		{'X', 'X', 'O', 'X'},
		{'X', 'O', 'X', 'X'},
	}
	fmt.Println("=== Approach 2: BFS from Border ===")
	solveBFS(b2)
	for _, row := range b2 {
		fmt.Println(string(row))
	}
}
